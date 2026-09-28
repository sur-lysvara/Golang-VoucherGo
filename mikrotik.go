package main

import (
	"bytes"
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"log"
)

func (a *App) routers(w http.ResponseWriter, r *http.Request) {
	// ADMIN_ROUTER_SEPARATE_PAGES_V1
	role, scopedRouterID := a.currentAdminRole(r)
	isSuper := role == "superadmin"

	var rows *sql.Rows
	var err error

	if role == "router" && scopedRouterID > 0 {
		rows, err = a.db.Query(`SELECT id,name,slug,COALESCE(hotspot_login_url,''),host,port,COALESCE(api_mode,'auto'),username,password,use_ssl,insecure_ssl,enabled,COALESCE(last_test_status,''),COALESCE(last_test_message,''),COALESCE(DATE_FORMAT(last_test_at,'%Y-%m-%d %H:%i:%s'),'') FROM routers WHERE id=? ORDER BY id DESC`, scopedRouterID)
	} else {
		rows, err = a.db.Query(`SELECT id,name,slug,COALESCE(hotspot_login_url,''),host,port,COALESCE(api_mode,'auto'),username,password,use_ssl,insecure_ssl,enabled,COALESCE(last_test_status,''),COALESCE(last_test_message,''),COALESCE(DATE_FORMAT(last_test_at,'%Y-%m-%d %H:%i:%s'),'') FROM routers ORDER BY id DESC`)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var list []Router

	for rows.Next() {
		var x Router
		if err := rows.Scan(&x.ID, &x.Name, &x.Slug, &x.HotspotLoginURL, &x.Host, &x.Port, &x.APIMode, &x.Username, &x.Password, &x.UseSSL, &x.InsecureSSL, &x.Enabled, &x.LastTestStatus, &x.LastTestMessage, &x.LastTestAt); err != nil {
			log.Printf("router scan error: %v", err)
			continue
		}
		x.LastTestOK = x.LastTestStatus == "ok"
		x.LastTestFail = x.LastTestStatus == "fail"
		list = append(list, x)
	}

	createMode := isSuper && r.URL.Query().Get("new") == "1"
	editMode := false

	edit := Router{
		Port:        443,
		APIMode:     "auto",
		UseSSL:      true,
		InsecureSSL: true,
		Enabled:     true,
	}

	if isSuper {
		if id := r.URL.Query().Get("edit"); id != "" {
			editID, _ := strconv.ParseInt(id, 10, 64)
			if editID <= 0 {
				http.Error(w, "Router tidak valid", http.StatusBadRequest)
				return
			}

			err := a.db.QueryRow(`SELECT id,name,slug,COALESCE(hotspot_login_url,''),host,port,COALESCE(api_mode,'auto'),username,password,use_ssl,insecure_ssl,enabled,COALESCE(last_test_status,''),COALESCE(last_test_message,''),COALESCE(DATE_FORMAT(last_test_at,'%Y-%m-%d %H:%i:%s'),'') FROM routers WHERE id=?`, editID).
				Scan(&edit.ID, &edit.Name, &edit.Slug, &edit.HotspotLoginURL, &edit.Host, &edit.Port, &edit.APIMode, &edit.Username, &edit.Password, &edit.UseSSL, &edit.InsecureSSL, &edit.Enabled, &edit.LastTestStatus, &edit.LastTestMessage, &edit.LastTestAt)
			if err != nil {
				http.Error(w, "Router tidak ditemukan", http.StatusNotFound)
				return
			}

			edit.LastTestOK = edit.LastTestStatus == "ok"
			edit.LastTestFail = edit.LastTestStatus == "fail"
			editMode = true
		}
	}

	testStatus := strings.TrimSpace(r.URL.Query().Get("router_test"))
	testMessage := strings.TrimSpace(r.URL.Query().Get("router_msg"))
	testNotice := testStatus == "notice"
	testOK := testStatus == "ok" || testNotice

	body := renderHTML(routerAdminPageTemplate, map[string]any{
		"List":        list,
		"Edit":        edit,
		"APIModes":    routerAPIModeOptions(edit.APIMode),
		"TestMessage": testMessage,
		"TestOK":      testOK,
		"TestNotice":  testNotice,
		"IsSuper":     isSuper,
		"CreateMode":  createMode,
		"EditMode":    editMode,
	})

	page(w, "Router", body)
}

func (a *App) routerSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/routers", http.StatusSeeOther)
		return
	}

	role, _ := a.currentAdminRole(r)
	if role != "superadmin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))

	port, err := strconv.Atoi(strings.TrimSpace(r.FormValue("port")))
	if err != nil || port <= 0 {
		port = 443
	}

	useSSL := r.FormValue("use_ssl") == "on"
	insecure := r.FormValue("insecure_ssl") == "on"
	enabled := r.FormValue("enabled") == "on"
	apiMode := normalizeRouterAPIMode(r.FormValue("api_mode"))

	if id == "0" || id == "" {
		_, err := a.db.Exec(
			`INSERT INTO routers(name,slug,hotspot_login_url,host,port,api_mode,username,password,use_ssl,insecure_ssl,enabled) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
			r.FormValue("name"),
			slugify(r.FormValue("slug")),
			strings.TrimSpace(r.FormValue("hotspot_login_url")),
			r.FormValue("host"),
			port,
			apiMode,
			r.FormValue("username"),
			r.FormValue("password"),
			useSSL,
			insecure,
			enabled,
		)
		if err != nil {
			log.Printf("router insert error: %v", err)
			http.Error(w, "Gagal menyimpan router", http.StatusInternalServerError)
			return
		}
	} else {
		routerID, err := strconv.ParseInt(id, 10, 64)
		if err != nil || routerID <= 0 {
			http.Error(w, "Router tidak valid", http.StatusBadRequest)
			return
		}

		password := r.FormValue("password")
		var res sql.Result

		if password != "" {
			res, err = a.db.Exec(
				`UPDATE routers SET name=?,slug=?,hotspot_login_url=?,host=?,port=?,api_mode=?,username=?,password=?,use_ssl=?,insecure_ssl=?,enabled=? WHERE id=?`,
				r.FormValue("name"),
				slugify(r.FormValue("slug")),
				strings.TrimSpace(r.FormValue("hotspot_login_url")),
				r.FormValue("host"),
				port,
				apiMode,
				r.FormValue("username"),
				password,
				useSSL,
				insecure,
				enabled,
				routerID,
			)
		} else {
			res, err = a.db.Exec(
				`UPDATE routers SET name=?,slug=?,hotspot_login_url=?,host=?,port=?,api_mode=?,username=?,use_ssl=?,insecure_ssl=?,enabled=? WHERE id=?`,
				r.FormValue("name"),
				slugify(r.FormValue("slug")),
				strings.TrimSpace(r.FormValue("hotspot_login_url")),
				r.FormValue("host"),
				port,
				apiMode,
				r.FormValue("username"),
				useSSL,
				insecure,
				enabled,
				routerID,
			)
		}

		if err != nil {
			log.Printf("router update error id=%d: %v", routerID, err)
			http.Error(w, "Gagal menyimpan router", http.StatusInternalServerError)
			return
		}

		if affected, err := res.RowsAffected(); err == nil && affected == 0 {
			var exists int
			if err := a.db.QueryRow(`SELECT COUNT(*) FROM routers WHERE id=?`, routerID).Scan(&exists); err != nil {
				log.Printf("router existence check error id=%d: %v", routerID, err)
				http.Error(w, "Gagal menyimpan router", http.StatusInternalServerError)
				return
			}
			if exists == 0 {
				http.Error(w, "Router tidak ditemukan", http.StatusNotFound)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/admin/routers?edit=%d&router_test=notice&router_msg=Data+router+masih+sama%%2C+jadi+tidak+ada+perubahan+yang+perlu+disimpan.", routerID), http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/admin/routers", http.StatusSeeOther)
}

func (a *App) routerDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/routers", http.StatusSeeOther)
		return
	}

	role, _ := a.currentAdminRole(r)
	if role != "superadmin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	routerID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("id")), 10, 64)
	if err != nil || routerID <= 0 {
		http.Error(w, "Router tidak valid", http.StatusBadRequest)
		return
	}

	res, err := a.db.Exec(`DELETE FROM routers WHERE id=?`, routerID)
	if err != nil {
		log.Printf("router delete error id=%d: %v", routerID, err)
		http.Error(w, "Gagal menghapus router", http.StatusInternalServerError)
		return
	}

	if affected, err := res.RowsAffected(); err == nil && affected == 0 {
		http.Error(w, "Router tidak ditemukan", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/admin/routers", http.StatusSeeOther)
}

func (a *App) publicRouter(w http.ResponseWriter, r *http.Request) {
	slug := strings.Trim(r.URL.Path, "/")
	if strings.HasPrefix(r.URL.Path, "/r/") {
		slug = strings.TrimPrefix(r.URL.Path, "/r/")
		slug = strings.Trim(slug, "/")
	}

	var ro Router

	err := a.db.QueryRow(`SELECT id,name,slug FROM routers WHERE slug=? AND enabled=1`, slug).Scan(&ro.ID, &ro.Name, &ro.Slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if r.URL.Query().Get("new_order") == "1" {
		clearPendingInvoiceCookie(w, r)
	} else if target, clear := a.pendingInvoiceRedirectForRouter(r, ro.Slug); target != "" {
		http.Redirect(w, r, target, http.StatusSeeOther)
		return
	} else if clear {
		clearPendingInvoiceCookie(w, r)
	}

	rows, err := a.db.Query(`SELECT id,name,price,profile,limit_uptime FROM packages WHERE router_id=? AND active=1 ORDER BY price`, ro.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var pkgs []Package

	for rows.Next() {
		var p Package
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Profile, &p.LimitUptime); err != nil {
			log.Printf("public package scan error: %v", err)
			continue
		}
		pkgs = append(pkgs, p)
	}

	publicSubtitle := strings.TrimSpace(a.setting("public_subtitle"))
	if publicSubtitle == "" {
		publicSubtitle = "Beli voucher hotspot, bayar via QRIS."
	}

	body := renderHTML(publicRouterPageTemplate, map[string]any{
		"Router":         ro,
		"Packages":       pkgs,
		"BuyExtra":       publicBuyExtraQuery(r, ro.Slug),
		"PublicSubtitle": publicSubtitle,
	})

	publicPage(w, "Selamat Datang di "+ro.Name, body) // PUBLIC_ROUTER_TITLE_WITH_NAME_V1
}

func (a *App) hotspotLoginRedirect(w http.ResponseWriter, r *http.Request) {
	routerSlug := strings.TrimSpace(r.URL.Query().Get("router"))

	u := ""
	if routerSlug != "" {
		_ = a.db.QueryRow(`SELECT COALESCE(hotspot_login_url,'') FROM routers WHERE slug=? AND enabled=1`, routerSlug).Scan(&u)
	}

	// Fallback untuk kompatibilitas setting lama.
	if strings.TrimSpace(u) == "" {
		u = strings.TrimSpace(a.setting("hotspot_login_url"))
	}

	if strings.TrimSpace(u) == "" {
		u = "http://login.kucing.net"
	}

	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

	http.Redirect(w, r, u, http.StatusFound)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")

	var b strings.Builder

	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}

	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "router" + randCode(4)
	}

	return out
}

func (a *App) provisionVoucher(txID int64) error {
	var ro Router
	var p Package
	var orderID, user, pass string

	err := a.db.QueryRow(`SELECT
			r.id,
			r.name,
			r.slug,
			COALESCE(r.host,''),
			r.port,
			COALESCE(r.api_mode,'auto'),
			COALESCE(r.username,''),
			COALESCE(r.password,''),
			r.use_ssl,
			r.insecure_ssl,
			r.enabled,
			p.id,
			p.name,
			p.price,
			COALESCE(p.profile,''),
			COALESCE(p.limit_uptime,''),
			p.shared_users,
			p.active,
			t.order_id,
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,'')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		JOIN packages p ON p.id=t.package_id
		WHERE t.id=?`, txID).
		Scan(
			&ro.ID,
			&ro.Name,
			&ro.Slug,
			&ro.Host,
			&ro.Port,
			&ro.APIMode,
			&ro.Username,
			&ro.Password,
			&ro.UseSSL,
			&ro.InsecureSSL,
			&ro.Enabled,
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Profile,
			&p.LimitUptime,
			&p.SharedUsers,
			&p.Active,
			&orderID,
			&user,
			&pass,
		)

	if err != nil {
		return err
	}

	if ro.Host == "" || ro.Username == "" || p.Profile == "" {
		_, _ = a.db.Exec(
			`UPDATE transactions SET provision_status='skipped', provision_error='router/profile belum lengkap' WHERE id=?`,
			txID,
		)
		return nil
	}

	mode := selectedMTAPIMode(ro.APIMode)

	var provisionErr error
	switch mode {
	case "rest":
		provisionErr = provisionVoucherREST(ro, p, orderID, user, pass)
	case "classic":
		provisionErr = provisionVoucherClassic(ro, p, orderID, user, pass, false, true)
	case "classic_ssl":
		provisionErr = provisionVoucherClassic(ro, p, orderID, user, pass, true, true)
	default:
		restErr := provisionVoucherREST(ro, p, orderID, user, pass)
		if restErr == nil {
			provisionErr = nil
			break
		}

		classicErr := provisionVoucherClassic(ro, p, orderID, user, pass, false, false)
		if classicErr == nil {
			provisionErr = nil
			break
		}

		provisionErr = fmt.Errorf("REST gagal: %v | Classic API gagal: %v", restErr, classicErr)
	}

	if provisionErr != nil {
		_, _ = a.db.Exec(`UPDATE transactions SET provision_status='failed', provision_error=? WHERE id=?`, provisionErr.Error(), txID)
		return provisionErr
	}

	_, _ = a.db.Exec(`UPDATE transactions SET provision_status='success', provision_error='' WHERE id=?`, txID)
	return nil
}

func provisionVoucherREST(ro Router, p Package, orderID, user, pass string) error {
	scheme := "http"
	if ro.UseSSL {
		scheme = "https"
	}

	u := fmt.Sprintf("%s://%s:%d/rest/ip/hotspot/user", scheme, ro.Host, ro.Port)

	body := map[string]string{
		"name":     user,
		"password": pass,
		"profile":  p.Profile,
		"comment":  voucherGoComment(orderID),
	}

	if p.LimitUptime != "" {
		body["limit-uptime"] = p.LimitUptime
	}

	b, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, u, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(ro.Username, ro.Password)

	resp, err := mtRestClientForRouter(ro).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rb, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		{
			// VOUCHERGO_MIKROTIK_DUPLICATE_USER_IDEMPOTENT_V1_CALL
			mtBodyText := string(rb)
			if isMikroTikDuplicateUserText(mtBodyText) {
				return nil
			}
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, mtBodyText)
		}
	}

	return nil
}

func (a *App) newVoucherCode() string {
	for i := 0; i < 50; i++ {
		code := randCode(4)

		var n int
		err := a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE voucher_username=?`, code).Scan(&n)
		if err == nil && n == 0 {
			return code
		}
	}

	return randCode(4)
}

func randCode(n int) string {
	alphabet := "abcdefhijkmnpqrstuvwxyz2345678"

	b := make([]byte, n)

	if _, err := crand.Read(b); err != nil {
		x := strconv.FormatInt(time.Now().UnixNano(), 36)
		if len(x) >= n {
			return x[:n]
		}
		return x
	}

	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}

	return string(b)
}

func voucherGoComment(orderID string) string {
	clean := strings.TrimSpace(orderID)

	var b strings.Builder
	for _, r := range clean {
		if r >= 'a' && r <= 'z' {
			b.WriteRune(r - 32)
			continue
		}
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}

	suffix := b.String()
	if len(suffix) > 3 {
		suffix = suffix[len(suffix)-3:]
	}
	if suffix == "" {
		suffix = strings.ToUpper(randCode(3))
	}

	return fmt.Sprintf("vc-vku-%s-VoucherGo%s", time.Now().Format("01.02.06"), suffix)
}
