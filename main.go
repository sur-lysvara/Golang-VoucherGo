package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	db         *sql.DB
	sessionKey []byte
}

func newMTRestClient(insecure bool) *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       120 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 12 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     false,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
				MinVersion:         tls.VersionTLS12,
				MaxVersion:         tls.VersionTLS12,
			},
		},
	}
}

var (
	mtRestClientVerified = newMTRestClient(false)
	mtRestClientInsecure = newMTRestClient(true)
)

func mtRestClientForRouter(ro Router) *http.Client {
	if ro.UseSSL && ro.InsecureSSL {
		return mtRestClientInsecure
	}

	return mtRestClientVerified
}

func main() {
	dsn := mustEnv("DB_DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	app := &App{
		db:         db,
		sessionKey: []byte(mustEnv("SESSION_KEY")),
	}
	go app.startUnusedVoucherReminderWorker()

	if err := app.migrate(); err != nil {
		log.Fatal(err)
	}

	if err := app.seedAdmin(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cek-voucher", app.publicMaintenance(app.publicVoucherCheck)) // VOUCHERGO_PUBLIC_CEK_VOUCHER_V1

	app.registerRoutes(mux)

	addr := env("APP_ADDR", "127.0.0.1:8097")
	log.Println("VoucherGo app running at", addr)
	log.Fatal(http.ListenAndServe(addr, securityHeaders(mux)))
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

// page renders admin pages with admin navigation.
// Re-added as a global function after previous layout refactor broke admin build.
