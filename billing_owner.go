package main

import (
	"bufio"
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// OWNER_BILLING_PAY_PHASE3A_V1
func readClientEnvFileBillingV1(slug string) map[string]string {
	out := map[string]string{}

	slug = strings.TrimSpace(slug)
	if slug == "" {
		return out
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,60}$`)
	if !re.MatchString(slug) {
		return out
	}

	candidates := []string{
		filepath.Join("/etc", "tuku-"+slug+".env"),
		filepath.Join("/var/www/vcr-clients", slug, "app", ".env"),
		filepath.Join("/var/www/vcr-clients", slug, ".env"),
	}

	var f *os.File
	for _, path := range candidates {
		fh, err := os.Open(path)
		if err == nil {
			f = fh
			out["_ENV_FILE"] = path
			break
		}
	}

	if f == nil {
		out["_ERROR"] = "env client tidak ditemukan. Dicek: " + strings.Join(candidates, ", ")
		return out
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		out[k] = v
	}

	return out
}

func (a *App) ensureClientBillingInvoicesV1() error {
	_, err := a.db.Exec(`
CREATE TABLE IF NOT EXISTS client_billing_invoices (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	order_id VARCHAR(100) NOT NULL UNIQUE,
	client_slug VARCHAR(80) NOT NULL,
	amount BIGINT NOT NULL,
	months INT NOT NULL DEFAULT 1,
	status ENUM('pending','paid','failed') NOT NULL DEFAULT 'pending',
	payment_method VARCHAR(80) NULL,
	raw_webhook MEDIUMTEXT NULL,
	paid_at DATETIME NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	INDEX(client_slug),
	INDEX(status),
	INDEX(created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return err
}

func billingAmountFromTextV1(v string) int64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}

	var b strings.Builder
	for _, r := range v {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	n, _ := strconv.ParseInt(b.String(), 10, 64)
	return n
}

func ownerBillingBaseURLV1(r *http.Request, configured string) string {
	configured = strings.TrimRight(strings.TrimSpace(configured), "/")
	if configured != "" {
		return configured
	}

	scheme := "https"
	if r.TLS == nil && strings.HasPrefix(r.Host, "127.0.0.1") {
		scheme = "http"
	}

	return scheme + "://" + r.Host
}

func billingStatusLabelV1(status string) (label, class, note string) {
	status = strings.ToLower(strings.TrimSpace(status))

	switch status {
	case "warning":
		return "Mendekati Jatuh Tempo", "warn", "Masa aktif aplikasi akan segera habis. Silakan lakukan perpanjangan sebelum jatuh tempo."
	case "grace":
		return "Masa Tenggang", "warn", "Aplikasi masih bisa digunakan sementara. Segera lakukan pembayaran agar tidak terkena suspend."
	case "expired":
		return "Expired", "bad", "Masa aktif aplikasi sudah habis. Lakukan pembayaran untuk menghindari suspend."
	case "suspended":
		return "Suspended", "bad", "Aplikasi sedang ditangguhkan. Silakan lakukan pembayaran untuk mengaktifkan kembali layanan."
	case "active", "":
		return "Active", "ok", "Layanan aplikasi aktif. Semua fitur dapat digunakan normal."
	default:
		return status, "warn", "Status billing membutuhkan pengecekan dari admin owner."
	}
}

func (a *App) ownerBillingCreateInvoiceV1(w http.ResponseWriter, r *http.Request, slug string, cfg map[string]string) {
	if err := a.ensureClientBillingInvoicesV1(); err != nil {
		http.Error(w, "gagal menyiapkan data pembayaran: "+err.Error(), http.StatusInternalServerError)
		return
	}

	project := strings.TrimSpace(a.setting("pakasir_slug"))
	apiKey := strings.TrimSpace(a.setting("pakasir_api_key"))
	if project == "" || apiKey == "" {
		http.Error(w, "Pengaturan pembayaran Pakasir belum lengkap.", http.StatusBadRequest)
		return
	}

	amount := billingAmountFromTextV1(cfg["CLIENT_BILLING_PRICE"])
	if amount <= 0 {
		http.Error(w, "Nominal sewa client belum valid.", http.StatusBadRequest)
		return
	}

	months := 1
	if m, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("months"))); m > 0 && m <= 24 {
		months = m
	}
	total := amount * int64(months)

	orderID := fmt.Sprintf("BILL-%s-%s", slug, time.Now().Format("20060102150405"))

	res, err := a.db.Exec(
		`INSERT INTO client_billing_invoices(order_id,client_slug,amount,months,status) VALUES(?,?,?,?, 'pending')`,
		orderID,
		slug,
		total,
		months,
	)
	if err != nil {
		http.Error(w, "gagal membuat pembayaran: "+err.Error(), http.StatusInternalServerError)
		return
	}

	invoiceID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "gagal membaca ID pembayaran: "+err.Error(), http.StatusInternalServerError)
		return
	}

	qris, err := a.createPakasirQRIS(orderID, total)
	if err != nil {
		_, _ = a.db.Exec(
			`UPDATE client_billing_invoices SET status='failed' WHERE id=?`,
			invoiceID,
		)
		http.Error(w, "gagal membuat transaksi Pakasir: "+err.Error(), http.StatusBadGateway)
		return
	}

	_, err = a.db.Exec(`
UPDATE client_billing_invoices
SET pakasir_txn_id=?,
    payment_method=?,
    qris_string=?,
    qris_total_payment=?,
    qris_expired_at=?
WHERE id=?`,
		qris.TxnID,
		qris.PaymentMethod,
		qris.QRISString,
		qris.TotalPayment,
		qris.ExpiredAt,
		invoiceID,
	)
	if err != nil {
		http.Error(w, "gagal menyimpan data QRIS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/billing/pay?client="+url.QueryEscape(slug)+"&order_id="+url.QueryEscape(orderID),
		http.StatusSeeOther,
	)
}

func formatRupiahBillingV1(n int64) string {
	if n <= 0 {
		return "Rp0"
	}

	raw := strconv.FormatInt(n, 10)
	var out []byte
	c := 0
	for i := len(raw) - 1; i >= 0; i-- {
		if c == 3 {
			out = append(out, '.')
			c = 0
		}
		out = append(out, raw[i])
		c++
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return "Rp" + string(out)
}

func (a *App) ownerBillingRecentInvoicesV1(slug string) []map[string]any {
	rows := []map[string]any{}

	slug = strings.TrimSpace(slug)
	if slug == "" {
		return rows
	}

	_ = a.ensureClientBillingInvoicesV1()

	q, err := a.db.Query(`
SELECT order_id, amount, months, status, COALESCE(payment_method,''), COALESCE(DATE_FORMAT(paid_at,'%Y-%m-%d %H:%i'),''), DATE_FORMAT(created_at,'%Y-%m-%d %H:%i')
FROM client_billing_invoices
WHERE client_slug=?
ORDER BY id DESC
LIMIT 10`, slug)
	if err != nil {
		return rows
	}
	defer q.Close()

	for q.Next() {
		var orderID, status, paymentMethod, paidAt, createdAt string
		var amount int64
		var months int

		if q.Scan(&orderID, &amount, &months, &status, &paymentMethod, &paidAt, &createdAt) != nil {
			continue
		}

		statusClass := "warn"
		statusLabel := status
		switch strings.ToLower(status) {
		case "paid":
			statusClass = "ok"
			statusLabel = "paid"
		case "failed":
			statusClass = "bad"
			statusLabel = "failed"
		case "pending":
			statusClass = "warn"
			statusLabel = "pending"
		}

		rows = append(rows, map[string]any{
			"OrderID":       orderID,
			"Amount":        amount,
			"AmountText":    formatRupiahBillingV1(amount),
			"Months":        months,
			"Status":        status,
			"StatusLabel":   statusLabel,
			"StatusClass":   statusClass,
			"PaymentMethod": paymentMethod,
			"PaidAt":        paidAt,
			"CreatedAt":     createdAt,
		})
	}

	return rows
}

func (a *App) ownerBillingInvoiceStatusV1(orderID string) map[string]any {
	out := map[string]any{
		"HasInvoice": false,
	}

	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return out
	}

	_ = a.ensureClientBillingInvoicesV1()

	var clientSlug, status, paymentMethod string
	var pakasirTxnID, qrisString, qrisExpiredAt string
	var amount, qrisTotalPayment int64
	var months int
	var paidAt string

	err := a.db.QueryRow(`
SELECT client_slug,
       amount,
       months,
       status,
       COALESCE(payment_method,''),
       COALESCE(pakasir_txn_id,''),
       COALESCE(qris_string,''),
       COALESCE(qris_total_payment,0),
       COALESCE(qris_expired_at,''),
       COALESCE(DATE_FORMAT(paid_at,'%Y-%m-%d %H:%i'),'')
FROM client_billing_invoices
WHERE order_id=?`, orderID).Scan(
		&clientSlug,
		&amount,
		&months,
		&status,
		&paymentMethod,
		&pakasirTxnID,
		&qrisString,
		&qrisTotalPayment,
		&qrisExpiredAt,
		&paidAt,
	)
	if err != nil {
		out["HasInvoice"] = false
		return out
	}

	out["HasInvoice"] = true
	out["OrderID"] = orderID
	out["ClientSlug"] = clientSlug
	out["Amount"] = amount
	out["Months"] = months
	out["Status"] = status
	out["PaymentMethod"] = paymentMethod
	out["PakasirTxnID"] = pakasirTxnID
	out["QRISString"] = qrisString
	out["QRISTotalPayment"] = qrisTotalPayment
	out["QRISExpiredAt"] = qrisExpiredAt
	out["PaidAt"] = paidAt
	return out
}

func (a *App) billingQRISPNG(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/billing-qris/")
	orderID = strings.TrimSuffix(orderID, ".png")
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		http.NotFound(w, r)
		return
	}

	var qris string
	err := a.db.QueryRow(`
SELECT COALESCE(qris_string,'')
FROM client_billing_invoices
WHERE order_id=?`, orderID).Scan(&qris)
	if err != nil || strings.TrimSpace(qris) == "" {
		http.NotFound(w, r)
		return
	}

	png, err := qrcode.Encode(qris, qrcode.Medium, 640)
	if err != nil {
		http.Error(w, "failed generate billing qris", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(png)
}

func (a *App) ownerBillingPay(w http.ResponseWriter, r *http.Request) {
	if env("OWNER_PANEL", "0") != "1" {
		http.NotFound(w, r)
		return
	}

	slug := strings.TrimSpace(r.URL.Query().Get("client"))
	if slug == "" {
		http.Error(w, "client kosong", http.StatusBadRequest)
		return
	}

	cfg := readClientEnvFileBillingV1(slug)
	if len(cfg) == 0 || cfg["_ERROR"] != "" {
		msg := cfg["_ERROR"]
		if msg == "" {
			msg = "client tidak ditemukan"
		}
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		a.ownerBillingCreateInvoiceV1(w, r, slug, cfg)
		return
	}

	status := strings.TrimSpace(cfg["CLIENT_BILLING_STATUS"])
	if status == "" {
		status = "active"
	}

	dueAt := strings.TrimSpace(cfg["CLIENT_BILLING_DUE_AT"])
	if dueAt == "" {
		dueAt = "-"
	}

	graceUntil := strings.TrimSpace(cfg["CLIENT_BILLING_GRACE_UNTIL"])
	if graceUntil == "" {
		graceUntil = "-"
	}

	price := strings.TrimSpace(cfg["CLIENT_BILLING_PRICE"])
	if price == "" {
		price = "-"
	}

	clientURL := strings.TrimSpace(cfg["PUBLIC_BASE_URL"])
	if clientURL == "" {
		clientURL = strings.TrimSpace(cfg["APP_PUBLIC_URL"])
	}
	if clientURL == "" {
		clientURL = "https://" + slug + ".vouchergo.biz.id"
	}

	statusLabel, statusClass, statusNote := billingStatusLabelV1(status)
	invoice := a.ownerBillingInvoiceStatusV1(r.URL.Query().Get("order_id"))

	body := renderHTML(`
<div class="owner-billing-pay-v1">
	<div class="owner-billing-pay-hero-v1">
		<div>
			<div class="owner-billing-kicker-v1">VoucherGo Billing</div>
			<h1>Pembayaran Sewa Aplikasi</h1>
			<p>Halaman pembayaran dan perpanjangan masa aktif aplikasi.</p>
		</div>
		<div class="owner-billing-status-v1 {{.StatusClass}}">{{.StatusLabel}}</div>
	</div>

	<div class="owner-billing-pay-grid-v1">
		<div class="card">
			<b>Client</b>
			<div class="price">{{.Slug}}</div>
			<p class="muted">{{.ClientURL}}</p>
		</div>
		<div class="card">
			<b>Jatuh Tempo</b>
			<div class="price">{{.DueAt}}</div>
			<p class="muted">Masa tenggang: {{.GraceUntil}}</p>
		</div>
		<div class="card">
			<b>Biaya Sewa</b>
			<div class="price">{{.Price}}</div>
			<p class="muted">Biaya sewa aplikasi per periode.</p>
		</div>
	</div>

	{{if .Invoice.HasInvoice}}
	<div class="card owner-billing-pay-card-v1">
		<h2>Status Pembayaran</h2>
		<p>Order: <b>{{.Invoice.OrderID}}</b></p>
		<p>Status: <b>{{.Invoice.Status}}</b>{{if .Invoice.PaidAt}} · Dibayar: {{.Invoice.PaidAt}}{{end}}</p>

		{{if and (eq .Invoice.Status "pending") .Invoice.QRISString}}
		<div style="margin-top:20px; text-align:center; padding:20px; border:1px solid #e2e8f0; border-radius:18px; background:#f8fafc;">
			<h3 style="margin:0 0 6px;">Selesaikan Pembayaran</h3>
			<p class="muted" style="margin:0 0 16px;">Scan QRIS berikut melalui m-Banking atau e-Wallet.</p>

			<a href="/billing-qris/{{.Invoice.OrderID}}.png" target="_blank" rel="noopener" style="display:inline-block; padding:10px; background:#fff; border:1px solid #e2e8f0; border-radius:14px;">
				<img src="/billing-qris/{{.Invoice.OrderID}}.png" alt="QRIS Pembayaran" style="width:280px; max-width:100%; display:block; border-radius:10px;">
			</a>

			{{if .Invoice.QRISTotalPayment}}
			<div style="margin-top:14px; font-size:22px; font-weight:800;">
				{{formatRupiahBillingV1 .Invoice.QRISTotalPayment}}
			</div>
			{{end}}

			<div style="margin-top:8px; font-size:13px; color:#64748b;">
				Order: {{.Invoice.OrderID}}
			</div>

			{{if .Invoice.QRISExpiredAt}}
			<div style="margin-top:4px; font-size:13px; color:#64748b;">
				Berlaku sampai: {{.Invoice.QRISExpiredAt}}
			</div>
			{{end}}
		</div>
		{{end}}
	</div>
	{{end}}

	<div class="card owner-billing-history-v1">
		<div class="owner-billing-history-head-v1">
			<div>
				<h2>Riwayat Pembayaran</h2>
				<p class="muted">Menampilkan 10 pembayaran terakhir untuk client ini.</p>
			</div>
		</div>

		{{if .Invoices}}
		<div class="owner-billing-history-list-v1">
			{{range .Invoices}}
			<div class="owner-billing-history-item-v1">
				<div class="owner-billing-history-main-v1">
					<b>{{.OrderID}}</b>
					<span>Dibuat: {{.CreatedAt}}</span>
					{{if .PaidAt}}<span>Dibayar: {{.PaidAt}}</span>{{end}}
				</div>
				<div class="owner-billing-history-meta-v1">
					<strong>{{.AmountText}}</strong>
					<small>{{.Months}} bulan</small>
				</div>
				<div class="owner-billing-history-action-v1">
					<span class="owner-billing-mini-status-v1 {{.StatusClass}}">{{.StatusLabel}}</span>
					<a class="btn light" href="/billing/pay?client={{$.Slug}}&order_id={{.OrderID}}">Detail</a>
				</div>
			</div>
			{{end}}
		</div>
		{{else}}
		<div class="owner-billing-empty-v1">Belum ada riwayat pembayaran untuk client ini.</div>
		{{end}}
	</div>

	<div class="card owner-billing-pay-card-v1">
		<h2>Instruksi Pembayaran</h2>
		<p>{{.StatusNote}}</p>

		<form method="post" class="owner-billing-form-v3a">
			<input type="hidden" name="months" value="1">
			<button class="btn pri" type="submit">Bayar Sekarang</button>
		</form>

		<div class="owner-billing-alert-v1">
			<b>Info:</b> Setelah pembayaran berhasil, masa aktif client akan diperpanjang otomatis.
		</div>

		<div class="actions">
			<a class="btn light" href="/admin/clients">Buka Client Manager</a>
			<a class="btn light" href="{{.ClientURL}}" target="_blank" rel="noopener">Buka Web Client</a>
		</div>
	</div>
</div>
`, map[string]any{
		"Slug":        slug,
		"Status":      status,
		"StatusLabel": statusLabel,
		"StatusClass": statusClass,
		"StatusNote":  statusNote,
		"DueAt":       dueAt,
		"GraceUntil":  graceUntil,
		"Price":       price,
		"ClientURL":   clientURL,
		"Invoice":     invoice,
		"Invoices":    a.ownerBillingRecentInvoicesV1(slug),
	})

	publicPage(w, "VoucherGo Billing", body)
}

func (a *App) handleClientBillingPakasirWebhookV2(h PakasirWebhookV2, raw []byte) bool {
	if err := a.ensureClientBillingInvoicesV1(); err != nil {
		return false
	}

	var id int64
	var amount int64
	var clientSlug string
	var months int
	var oldStatus string
	var pakasirTxnID string
	var paymentMethod string

	err := a.db.QueryRow(`
SELECT id,
       amount,
       client_slug,
       months,
       status,
       COALESCE(pakasir_txn_id,''),
       COALESCE(payment_method,'')
FROM client_billing_invoices
WHERE order_id=?`,
		h.OrderID,
	).Scan(
		&id,
		&amount,
		&clientSlug,
		&months,
		&oldStatus,
		&pakasirTxnID,
		&paymentMethod,
	)
	if err != nil {
		return false
	}

	if pakasirTxnID == "" ||
		pakasirTxnID != h.TxnID ||
		amount != h.Amount {
		_, _ = a.db.Exec(
			`UPDATE client_billing_invoices SET raw_webhook=? WHERE id=?`,
			string(raw),
			id,
		)
		return false
	}

	if oldStatus == "paid" {
		return true
	}

	if !a.verifyPakasirV2(h.TxnID, h.OrderID, amount) {
		_, _ = a.db.Exec(
			`UPDATE client_billing_invoices SET raw_webhook=? WHERE id=?`,
			string(raw),
			id,
		)
		return false
	}

	res, err := a.db.Exec(`
UPDATE client_billing_invoices
SET status='paid',
    payment_method=?,
    raw_webhook=?,
    paid_at=IFNULL(paid_at,NOW())
WHERE id=? AND status<>'paid'`,
		paymentMethod,
		string(raw),
		id,
	)
	if err != nil {
		return false
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return true
	}

	// OWNER_BILLING_AUTO_EXTEND_PHASE3B_V1
	a.applyClientBillingPaidV1(clientSlug, months)

	return true
}

func (a *App) applyClientBillingPaidV1(clientSlug string, months int) {
	clientSlug = strings.TrimSpace(clientSlug)
	if clientSlug == "" {
		return
	}
	if months <= 0 {
		months = 1
	}

	cmd := exec.Command("sudo", "/usr/local/sbin/vcr-client-billing.sh", "paid", clientSlug, strconv.Itoa(months))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("client billing apply failed slug=%s months=%d err=%v output=%s", clientSlug, months, err, strings.TrimSpace(string(out)))
		return
	}

	log.Printf("client billing applied slug=%s months=%d output=%s", clientSlug, months, strings.TrimSpace(string(out)))
}
