package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"encoding/base64"
)

func (a *App) healthAdmin(w http.ResponseWriter, r *http.Request) {
	// ADMIN_HEALTH_PAGE_V1
	// ADMIN_HEALTH_COMPACT_V1
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	dbStatus := "OK"
	dbError := ""

	if err := a.db.PingContext(ctx); err != nil {
		dbStatus = "ERROR"
		dbError = err.Error()
	}

	var routerCount, packageCount, txCount, paidCount, pendingCount int
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM routers`).Scan(&routerCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM packages`).Scan(&packageCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions`).Scan(&txCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status='paid'`).Scan(&paidCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status='pending'`).Scan(&pendingCount)

	waStatus := "ERROR"
	waProvider := "-"
	waConnection := "-"
	waConnected := false
	waRaw := ""
	waHTTP := 0

	code, bodyText, err := a.localWACall(ctx, http.MethodGet, "/status", nil)
	waHTTP = code
	waRaw = maskHealthWARaw(bodyText)

	if err != nil {
		waRaw = err.Error()
	} else {
		var st map[string]any
		if json.Unmarshal([]byte(bodyText), &st) == nil {
			if v, ok := st["provider"].(string); ok {
				waProvider = v
			}
			if v, ok := st["connection"].(string); ok {
				waConnection = v
			}
			if v, ok := st["connected"].(bool); ok {
				waConnected = v
			}

			if waConnected {
				waStatus = "OK"
			} else {
				waStatus = "DISCONNECTED"
			}
		}
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	body := renderHTML(`
<div class="health-page">
	<div class="card health-hero">
		<div>
			<h2>Health Status</h2>
			<p class="muted">Ringkasan kondisi aplikasi, database, WhatsApp, dan transaksi.</p>
		</div>
		<div class="health-pill {{if eq .WAStatus "OK"}}ok{{else}}bad{{end}}">
			WA {{.WAStatus}}
		</div>
	</div>

	<div class="grid health-grid health-top-grid">
		<div class="card health-card">
			<b>Database</b>
			<div class="price">{{.DBStatus}}</div>
			{{if .DBError}}<p class="muted">{{.DBError}}</p>{{end}}
		</div>

		<div class="card health-card">
			<b>WhatsApp</b>
			<div class="price">{{.WAStatus}}</div>
			<p class="muted">{{.WAProvider}} · {{.WAConnection}} · HTTP {{.WAHTTP}}</p>
		</div>

		<div class="card health-card">
			<b>Go Memory</b>
			<div class="price">{{.GoAlloc}}</div>
			<p class="muted">Sys: {{.GoSys}}</p>
		</div>
	</div>

	<div class="grid health-grid health-stat-grid">
		<div class="card health-mini"><b>Router</b><div class="price">{{.RouterCount}}</div></div>
		<div class="card health-mini"><b>Paket</b><div class="price">{{.PackageCount}}</div></div>
		<div class="card health-mini"><b>Transaksi</b><div class="price">{{.TxCount}}</div></div>
		<div class="card health-mini"><b>Paid</b><div class="price">{{.PaidCount}}</div></div>
		<div class="card health-mini"><b>Pending</b><div class="price">{{.PendingCount}}</div></div>
	</div>

	<div class="card health-settings">
		<h2>Setting WhatsApp</h2>
		<div class="grid health-setting-grid">
			<div class="field">
				<label>Provider</label>
				<input readonly value="{{.WASettingProvider}}">
			</div>
			<div class="field">
				<label>Local URL</label>
				<input readonly value="{{.WASettingURL}}">
			</div>
			<div class="field">
				<label>Wanesia Enabled</label>
				<input readonly value="{{.WanesiaEnabled}}">
			</div>
		</div>
	</div>

	<details class="card health-raw-card">
		<summary>Raw WhatsApp Status</summary>
		<pre class="wa-raw">{{.WARaw}}</pre>
	</details>
</div>`, map[string]any{
		"DBStatus":          dbStatus,
		"DBError":           dbError,
		"WAStatus":          waStatus,
		"WAProvider":        waProvider,
		"WAConnection":      waConnection,
		"WAHTTP":            waHTTP,
		"WARaw":             waRaw,
		"GoAlloc":           formatHealthBytes(mem.Alloc),
		"GoSys":             formatHealthBytes(mem.Sys),
		"RouterCount":       routerCount,
		"PackageCount":      packageCount,
		"TxCount":           txCount,
		"PaidCount":         paidCount,
		"PendingCount":      pendingCount,
		"WASettingProvider": a.setting("wa_provider"),
		"WASettingURL":      a.setting("wa_local_url"),
		"WanesiaEnabled":    a.setting("wanesia_enabled"),
	})

	page(w, "Health", body)
}

// ADMIN_WA_QR_PROXY_V1

// WA_GUI_STATUS_JSON_V1
func (a *App) whatsappStatusJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	code, bodyText, err := a.localWACall(r.Context(), http.MethodGet, "/status", nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"ok":false,"connected":false,"connection":"error","state":"error","error":"gagal ambil status WhatsApp lokal"}`))
		return
	}

	if code < 200 || code >= 300 {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"ok":false,"connected":false,"connection":"error","state":"error","error":"status WhatsApp lokal tidak normal"}`))
		return
	}

	_, _ = w.Write([]byte(bodyText))
}

func (a *App) whatsappQRProxy(w http.ResponseWriter, r *http.Request) {
	// VOUCHERGO_WEBJS_QR_PROXY_DATAURL_V1
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	code, body, err := a.localWACall(ctx, http.MethodGet, "/qr", nil)
	if err != nil {
		http.Error(w, "gagal ambil QR WhatsApp: "+err.Error(), http.StatusBadGateway)
		return
	}
	if code < 200 || code >= 300 {
		http.Error(w, fmt.Sprintf("gagal ambil QR WhatsApp: HTTP %d", code), http.StatusBadGateway)
		return
	}

	// Provider baru seperti whatsapp-web.js mengembalikan JSON:
	// {"qr":"data:image/png;base64,...","qr_data_url":"data:image/png;base64,..."}
	var qrResp map[string]any
	if json.Unmarshal([]byte(body), &qrResp) == nil {
		for _, key := range []string{"qr", "qr_data_url", "qrDataUrl", "data_url", "image"} {
			if raw, ok := qrResp[key].(string); ok && raw != "" {
				if strings.HasPrefix(raw, "data:image/") {
					comma := strings.Index(raw, ",")
					if comma > 0 {
						img, decErr := base64.StdEncoding.DecodeString(raw[comma+1:])
						if decErr == nil && len(img) > 0 {
							w.Header().Set("Content-Type", "image/png")
							w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
							_, _ = w.Write(img)
							return
						}
					}
				}
			}
		}
	}

	// Fallback untuk provider lama yang langsung balikin PNG bytes.
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	_, _ = w.Write([]byte(body))
}

func (a *App) whatsappAdmin(w http.ResponseWriter, r *http.Request) {
	// ADMIN_WHATSAPP_MENU_V1
	// ADMIN_WHATSAPP_RESPONSIVE_V1
	// ADMIN_DEMO_DISABLE_WA_UI_V1
	msg := strings.TrimSpace(r.URL.Query().Get("msg"))

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		action := r.FormValue("action")
		ctx, cancel := context.WithTimeout(r.Context(), 35*time.Second)
		defer cancel()

		switch action {
		case "demo_wa_toggle":
			pin := strings.TrimSpace(r.FormValue("pin"))
			wantPin := strings.TrimSpace(env("DEMO_CONTROL_PIN", ""))

			if wantPin == "" {
				msg = "DEMO_CONTROL_PIN belum diset di env demo."
				break
			}

			if pin != wantPin {
				msg = "PIN kontrol demo salah."
				break
			}

			enabled := "0"
			if r.FormValue("enabled") == "1" {
				enabled = "1"
			}

			_, err := a.db.Exec(`INSERT INTO settings(k,v) VALUES('demo_wa_enabled',?) ON DUPLICATE KEY UPDATE v=?`, enabled, enabled)
			if err != nil {
				msg = "Gagal update mode WA demo: " + err.Error()
			} else if enabled == "1" {
				msg = "Private trial WA demo aktif. QR dan test kirim WA ditampilkan."
			} else {
				msg = "Public demo aman aktif. QR dan test kirim WA disembunyikan."
			}

		case "reconnect":
			code, body, err := a.localWACall(ctx, http.MethodPost, "/reconnect", map[string]string{})
			if err != nil {
				msg = "Reconnect gagal: " + err.Error()
			} else {
				msg = fmt.Sprintf("Reconnect response HTTP %d: %s", code, truncateLocalWAString(body, 300))
			}

		case "logout":
			code, body, err := a.localWACall(ctx, http.MethodPost, "/logout", map[string]string{})
			if err != nil {
				msg = "Logout gagal: " + err.Error()
			} else {
				msg = fmt.Sprintf("Logout response HTTP %d: %s", code, truncateLocalWAString(body, 300))
			}

		case "test":
			rawNumber := strings.TrimSpace(r.FormValue("number"))
			number := ""
			for _, ch := range rawNumber {
				if ch >= '0' && ch <= '9' {
					number += string(ch)
				}
			}
			if strings.HasPrefix(number, "08") {
				number = "62" + strings.TrimPrefix(number, "0")
			} else if strings.HasPrefix(number, "8") {
				number = "62" + number
			}

			message := strings.TrimSpace(r.FormValue("message"))
			if message == "" {
				message = "Tes WhatsApp VoucherGo berhasil."
			}

			if number == "" || len(number) < 10 {
				msg = "Nomor WhatsApp tidak valid. Contoh: 6285775898626"
				break
			}

			code, body, err := a.localWACall(ctx, http.MethodPost, "/send", map[string]string{
				"to":      number,
				"number":  number,
				"phone":   number,
				"message": message,
				"text":    message,
			})

			if err != nil {
				msg = "Test kirim gagal: " + err.Error()
			} else if code >= 200 && code < 300 {
				msg = fmt.Sprintf("Test kirim berhasil ke %s. Response HTTP %d: %s", number, code, truncateLocalWAString(body, 500))
			} else {
				msg = fmt.Sprintf("Test kirim gagal HTTP %d: %s", code, truncateLocalWAString(body, 500))
			}
		}
	}

	// ADMIN_WA_POST_REDIRECT_V1
	// Hindari browser mengulang POST saat halaman di-refresh.
	if r.Method == http.MethodPost {
		redirectTo := "/admin/whatsapp"
		if strings.TrimSpace(msg) != "" {
			redirectTo += "?msg=" + url.QueryEscape(msg)
		}
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	statusCode, statusBody, statusErr := a.localWACall(ctx, http.MethodGet, "/status", nil)

	statusText := ""
	provider := "-"
	connection := "-"
	connected := false
	me := "-"

	if statusErr != nil {
		statusText = statusErr.Error()
	} else {
		statusText = statusBody

		var st map[string]any
		if json.Unmarshal([]byte(statusBody), &st) == nil {
			if v, ok := st["provider"].(string); ok && v != "" {
				provider = v
			}
			if v, ok := st["connection"].(string); ok && v != "" {
				connection = v
			}
			if v, ok := st["connected"].(bool); ok {
				connected = v
			}
			if v, ok := st["me"]; ok {
				b, _ := json.Marshal(v)
				if len(b) > 0 {
					me = string(b)
				}
			}
		}
	}

	// ADMIN_WHATSAPP_DEMO_MASK_V1
	demoMode := env("DEMO_MODE", "0") == "1"
	demoWAEnabled := env("DEMO_WA_ENABLED", "0") == "1"
	if demoMode {
		var demoWASetting string
		_ = a.db.QueryRow(`SELECT COALESCE(v,'') FROM settings WHERE k='demo_wa_enabled'`).Scan(&demoWASetting)
		demoWASetting = strings.TrimSpace(demoWASetting)
		if demoWASetting == "1" {
			demoWAEnabled = true
		} else if demoWASetting == "0" {
			demoWAEnabled = false
		}

		if strings.TrimSpace(me) != "" && me != "-" && me != "null" {
			me = "[masked in demo mode]"
		}
		statusText = maskHealthWARaw(statusText)
	}

	qrData := ""
	if strings.TrimSpace(strings.ToLower(me)) == "null" || strings.TrimSpace(me) == "" {
		me = "-"
	}

	qrCode, qrBody, qrErr := a.localWACall(ctx, http.MethodGet, "/qr", nil)
	_ = qrData // ADMIN_WA_QR_PROXY_UNUSED_LEGACY_V1
	qrProxyURL := ""
	if qrErr == nil && qrCode == http.StatusOK {
		qrProxyURL = fmt.Sprintf("/admin/whatsapp/qr.png?t=%d", time.Now().Unix())
	}
	if qrErr == nil && qrCode >= 200 && qrCode < 300 {
		var qr map[string]any
		if json.Unmarshal([]byte(qrBody), &qr) == nil {
			if v, ok := qr["data_url"].(string); ok {
				qrData = v
			}
		}
	}

	// WA_MESSAGE_STATS_CARD_V1
	var waSentCount, waPendingCount, waFailedCount int
	_ = a.db.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN COALESCE(wa_status,'pending')='sent' THEN 1 ELSE 0 END),0),
			COALESCE(SUM(CASE WHEN COALESCE(wa_status,'pending') IN ('pending','sending','') THEN 1 ELSE 0 END),0),
			COALESCE(SUM(CASE WHEN COALESCE(wa_status,'pending') IN ('failed','skipped') THEN 1 ELSE 0 END),0)
		FROM transactions
	`).Scan(&waSentCount, &waPendingCount, &waFailedCount)

	waTotalCount := waSentCount + waPendingCount + waFailedCount
	waSentDeg := "0"
	waPendingEndDeg := "0"
	if waTotalCount > 0 {
		sentDeg := float64(waSentCount) / float64(waTotalCount) * 360
		pendingDeg := float64(waPendingCount) / float64(waTotalCount) * 360
		waSentDeg = fmt.Sprintf("%.2f", sentDeg)
		waPendingEndDeg = fmt.Sprintf("%.2f", sentDeg+pendingDeg)
	}

	body := renderHTML(`
<div class="wa-page">
	<div class="card wa-hero">
		<div>
			<span class="wa-mobile-compact-v1" hidden></span>

<style>
/* WA_MOBILE_COMPACT_FINAL_V3 */
@media (max-width: 720px){
  body:has(.wa-mobile-compact-v1) .wa-test-form{
    grid-template-columns:1fr!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form .field{
    width:100%!important;
    min-width:0!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form input,
  body:has(.wa-mobile-compact-v1) .wa-test-form button{
    width:100%!important;
    box-sizing:border-box!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form input[name="message"]{
    font-size:12px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-device-card{
    display:none!important;
  }

  body:has(.wa-mobile-compact-v1) .grid{
    margin-bottom:6px!important;
  }

  body:has(.wa-mobile-compact-v1) .grid > .card{
    min-height:62px!important;
  }

  body:has(.wa-mobile-compact-v1) .notice{
    margin:7px 0!important;
  }
}
</style>

<style>
/* WA_MOBILE_COMPACT_V1 */
@media (max-width: 720px){
  body:has(.wa-mobile-compact-v1) .card{
    padding:11px!important;
    border-radius:14px!important;
    margin-bottom:9px!important;
  }

  body:has(.wa-mobile-compact-v1) .grid{
    display:grid!important;
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:8px!important;
    margin-bottom:8px!important;
  }

  body:has(.wa-mobile-compact-v1) .grid > .card{
    min-width:0!important;
  }

  body:has(.wa-mobile-compact-v1) .grid > .card:nth-child(3){
    grid-column:1/-1!important;
  }

  body:has(.wa-mobile-compact-v1) .notice,
  body:has(.wa-mobile-compact-v1) .wa-qr-note{
    padding:10px 12px!important;
    margin:8px 0!important;
    border-radius:14px!important;
    font-size:13px!important;
    line-height:1.35!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-action-card,
  body:has(.wa-mobile-compact-v1) .wa-test-card,
  body:has(.wa-mobile-compact-v1) .wa-qr-box{
    padding:11px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-actions{
    display:grid!important;
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:8px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-actions form,
  body:has(.wa-mobile-compact-v1) .wa-actions button{
    width:100%!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-actions button{
    min-height:38px!important;
    padding:9px 8px!important;
    font-size:12px!important;
    line-height:1.15!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form{
    display:grid!important;
    grid-template-columns:1fr!important;
    gap:8px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form .field{
    margin:0!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form input{
    min-height:38px!important;
    padding:9px 10px!important;
    font-size:13px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-form button{
    width:100%!important;
    min-height:38px!important;
    padding:9px 10px!important;
    font-size:13px!important;
  }

  body:has(.wa-mobile-compact-v1) h2{
    font-size:16px!important;
    margin-bottom:6px!important;
  }

  body:has(.wa-mobile-compact-v1) label{
    font-size:11px!important;
    margin-bottom:5px!important;
  }

  body:has(.wa-mobile-compact-v1) .price,
  body:has(.wa-mobile-compact-v1) b{
    font-size:14px!important;
  }

  body:has(.wa-mobile-compact-v1) .badge{
    padding:5px 8px!important;
    font-size:11px!important;
  }

  body:has(.wa-mobile-compact-v1) details,
  body:has(.wa-mobile-compact-v1) pre,
  body:has(.wa-mobile-compact-v1) .wa-raw{
    font-size:11px!important;
    max-height:150px!important;
    overflow:auto!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-qr-img{
    max-width:220px!important;
    width:100%!important;
    margin:auto!important;
    display:block!important;
  }
}

@media (max-width: 380px){
  body:has(.wa-mobile-compact-v1) .card{
    padding:10px!important;
  }

  body:has(.wa-mobile-compact-v1) .grid{
    gap:7px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-actions{
    grid-template-columns:1fr!important;
  }
}
</style>

<style>
/* WA_MOBILE_COMPACT_V2 */
@media (max-width: 720px){
  body:has(.wa-mobile-compact-v1) .grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }

  body:has(.wa-mobile-compact-v1) .grid > .card:nth-child(3){
    grid-column:auto!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-device-card{
    grid-column:1/-1!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-device-card small{
    display:block!important;
    max-height:42px!important;
    overflow:auto!important;
    font-size:10px!important;
    line-height:1.25!important;
    word-break:break-all!important;
  }

  body:has(.wa-mobile-compact-v1) .health-pill,
  body:has(.wa-mobile-compact-v1) .badge{
    transform:scale(.92);
    transform-origin:left center;
  }

  body:has(.wa-mobile-compact-v1) .notice:first-of-type{
    font-size:12px!important;
    padding:8px 10px!important;
    max-height:46px!important;
    overflow:auto!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-action-card h2,
  body:has(.wa-mobile-compact-v1) .wa-test-card h2{
    margin-bottom:8px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-action-card{
    margin-top:8px!important;
  }

  body:has(.wa-mobile-compact-v1) .wa-test-card{
    margin-top:8px!important;
  }
}
</style>

			<h2>WhatsApp Gateway</h2>
			<p class="muted">{{if and .Provider (ne .Provider "-")}}Gateway lokal aktif: {{.Provider}}{{else}}Gateway WhatsApp lokal belum dikonfigurasi{{end}}</p>
		</div>
		<div class="wa-hero-pill {{if .Connected}}ok{{else}}bad{{end}}">
			{{if .Connected}}ONLINE{{else}}OFFLINE{{end}}
		</div>
	</div>

	{{if .Msg}}
	<div class="notice wa-message">{{.Msg}}</div>
	{{end}}

	<div class="grid wa-status-grid">
		<div class="card wa-stat">
			<b>Provider</b>
			<span>{{.Provider}}</span>
		</div>
		<div class="card wa-stat">
			<b>Koneksi</b>
			<span class="badge {{if .Connected}}sent{{else}}failed{{end}}">{{.Connection}}</span>
		</div>
		<div class="card wa-stat">
			<b>Status</b>
			<span>{{if .Connected}}Connected{{else}}QR Required{{end}}</span>
		</div>
	</div>

	
	<style>
	/* WA_MESSAGE_STATS_CARD_V1 */
	body:has(.wa-message-stats-v1) .wa-message-stats-v1{
		display:grid;
		grid-template-columns:minmax(260px,1fr) minmax(260px,.75fr);
		gap:18px;
		align-items:center;
		margin-top:14px;
	}
	body:has(.wa-message-stats-v1) .wa-message-chart-wrap-v1{
		display:flex;
		align-items:center;
		justify-content:center;
		min-height:260px;
	}
	body:has(.wa-message-stats-v1) .wa-message-donut-v1{
		width:230px;
		height:230px;
		border-radius:999px;
		background:conic-gradient(#22cfa2 0deg {{.WASentDeg}}deg, #ff8a0a {{.WASentDeg}}deg {{.WAPendingEndDeg}}deg, #ef4b4b {{.WAPendingEndDeg}}deg 360deg);
		position:relative;
		box-shadow:0 18px 45px rgba(0,0,0,.24);
	}
	body:has(.wa-message-stats-v1) .wa-message-donut-v1:after{
		content:"";
		position:absolute;
		inset:64px;
		border-radius:999px;
		background:var(--card, #20142d);
		box-shadow:inset 0 0 0 1px rgba(255,255,255,.07);
	}
	body:has(.wa-message-stats-v1) .wa-message-legend-v1{
		display:flex;
		justify-content:center;
		gap:14px;
		flex-wrap:wrap;
		margin-bottom:18px;
		font-size:13px;
		color:var(--muted, #94a3b8);
	}
	body:has(.wa-message-stats-v1) .wa-message-legend-v1 span{
		display:inline-flex;
		align-items:center;
		gap:7px;
	}
	body:has(.wa-message-stats-v1) .wa-dot-v1{
		width:13px;
		height:13px;
		border-radius:999px;
		display:inline-block;
	}
	body:has(.wa-message-stats-v1) .wa-message-side-v1{
		display:grid;
		gap:12px;
	}
	body:has(.wa-message-stats-v1) .wa-message-side-item-v1{
		display:grid;
		grid-template-columns:auto 1fr;
		gap:14px;
		align-items:center;
		padding:16px 18px;
		border-radius:18px;
		background:rgba(255,255,255,.055);
		border:1px solid rgba(255,255,255,.075);
	}
	body:has(.wa-message-stats-v1) .wa-message-side-item-v1 .num{
		font-size:28px;
		font-weight:900;
		line-height:1;
	}
	body:has(.wa-message-stats-v1) .wa-message-side-item-v1 .label{
		color:var(--muted, #94a3b8);
		margin-top:4px;
	}
	@media(max-width:720px){
		body:has(.wa-message-stats-v1) .wa-message-stats-v1{
			grid-template-columns:1fr;
			gap:12px;
			margin-top:10px;
		}
		body:has(.wa-message-stats-v1) .wa-message-chart-wrap-v1{
			min-height:190px;
		}
		body:has(.wa-message-stats-v1) .wa-message-donut-v1{
			width:170px;
			height:170px;
		}
		body:has(.wa-message-stats-v1) .wa-message-donut-v1:after{
			inset:48px;
		}
		body:has(.wa-message-stats-v1) .wa-message-side-v1{
			grid-template-columns:1fr;
		}
		body:has(.wa-message-stats-v1) .wa-message-side-item-v1{
			padding:12px 14px;
		}
		body:has(.wa-message-stats-v1) .wa-message-side-item-v1 .num{
			font-size:22px;
		}
	}
	</style>

	<div class="card wa-message-stats-v1">
		<div>
			<h2>Pesan Terkirim</h2>
			<div class="wa-message-legend-v1">
				<span><i class="wa-dot-v1" style="background:#22cfa2"></i>Terkirim</span>
				<span><i class="wa-dot-v1" style="background:#ff8a0a"></i>Pending</span>
				<span><i class="wa-dot-v1" style="background:#ef4b4b"></i>Gagal</span>
			</div>
			<div class="wa-message-chart-wrap-v1">
				<div class="wa-message-donut-v1" aria-label="Statistik pesan WhatsApp"></div>
			</div>
		</div>
		<div class="wa-message-side-v1">
			<div class="wa-message-side-item-v1">
				<i class="wa-dot-v1" style="background:#22cfa2"></i>
				<div><div class="num">{{.WASentCount}}</div><div class="label">Terkirim</div></div>
			</div>
			<div class="wa-message-side-item-v1">
				<i class="wa-dot-v1" style="background:#ff8a0a"></i>
				<div><div class="num">{{.WAPendingCount}}</div><div class="label">Pending</div></div>
			</div>
			<div class="wa-message-side-item-v1">
				<i class="wa-dot-v1" style="background:#ef4b4b"></i>
				<div><div class="num">{{.WAFailedCount}}</div><div class="label">Gagal</div></div>
			</div>
		</div>
	</div>

<div class="card wa-device-card">
		<b>Device</b>
		<small>{{.Me}}</small>
	</div>

	{{if .DemoMode}}
	<div class="card wa-demo-control-card">
		<h2>Mode WhatsApp Demo</h2>
		<p class="muted">{{if .DemoWAEnabled}}Private trial aktif: QR dan test WA bisa dipakai.{{else}}Public demo aman: QR dan test WA terkunci.{{end}}</p>

		<form method="post" class="wa-demo-control-form">
			<input type="hidden" name="action" value="demo_wa_toggle">
			<input type="hidden" name="enabled" value="{{if .DemoWAEnabled}}0{{else}}1{{end}}">
			<div class="field">
				<label>PIN Kontrol Demo</label>
				<input name="pin" type="password" placeholder="masukkan PIN owner" required>
			</div>
			<div class="field">
				<label>&nbsp;</label>
				<button class="btn {{if .DemoWAEnabled}}danger{{else}}pri{{end}}" type="submit">
					{{if .DemoWAEnabled}}Kunci WA Demo{{else}}Aktifkan Private Trial WA{{end}}
				</button>
			</div>
		</form>
	</div>
	{{end}}

	{{if and .DemoMode (not .DemoWAEnabled)}}
	<div class="notice wa-demo-safe-note">
		<b>Demo Mode aktif.</b><br>
		QR WhatsApp, logout, reconnect, dan test kirim WA disembunyikan agar aman saat demo ke client.
	</div>
	{{else if and .QRData (or (not .DemoMode) .DemoWAEnabled)}}
	<div class="card wa-qr-box">
		<div>
			<h2>Scan QR WhatsApp</h2>
			<p class="muted">WhatsApp HP → Perangkat Tertaut → Tautkan Perangkat</p>
		</div>
		<img class="wa-qr-img" src="{{.QRData}}" alt="QR WhatsApp">
	</div>
	{{else}}
	<div class="notice wa-qr-note">
		{{if .Connected}}
		WhatsApp sudah terhubung. QR tidak diperlukan.
		{{else}}
		QR belum tersedia. Klik Reconnect lalu refresh halaman.
		{{end}}
	</div>
	{{end}}

	{{if or (not .DemoMode) .DemoWAEnabled}}
	<div class="card wa-action-card">
		<h2>Kontrol Koneksi</h2>
		<div class="actions wa-actions">
			<form method="post">
				<input type="hidden" name="action" value="reconnect">
				<button class="btn pri" type="submit">Reconnect / Refresh QR</button>
			</form>

			<form method="post" onsubmit="return confirm('Logout WhatsApp dan hapus sesi? Nanti harus scan QR ulang.');">
				<input type="hidden" name="action" value="logout">
				<button class="btn danger" type="submit">Logout WhatsApp</button>
			</form>
		</div>
	</div>

	<div class="card wa-test-card">
		<h2>Test Kirim WhatsApp</h2>

		<form method="post" class="grid wa-test-form">
			<input type="hidden" name="action" value="test">

			<div class="field">
				<label>Nomor WhatsApp</label>
				<input name="number" placeholder="" inputmode="numeric" required>
			</div>

			<div class="field">
				<label>Pesan</label>
				<input name="message" value="Tes WhatsApp VoucherGo dari admin.">
			</div>

			<div class="field">
				<label>&nbsp;</label>
				<button class="btn ok" type="submit">Kirim Test</button>
			</div>
		</form>
	</div>

	{{end}}

	<details class="card wa-raw-card">
		<summary>Raw Status</summary>
		<pre class="wa-raw">{{.StatusCode}} {{.StatusText}}</pre>
	</details>
</div>`, map[string]any{
		"Msg":             msg,
		"WASentCount":     waSentCount,
		"WAPendingCount":  waPendingCount,
		"WAFailedCount":   waFailedCount,
		"WATotalCount":    waTotalCount,
		"WASentDeg":       waSentDeg,
		"WAPendingEndDeg": waPendingEndDeg,
		"StatusCode":      statusCode,
		"StatusText":      statusText,
		"Provider":        provider,
		"Connection":      connection,
		"Connected":       connected,
		"Me":              me,
		"QRData":          qrProxyURL,
		"DemoMode":        demoMode,
		"DemoWAEnabled":   demoWAEnabled,
	})

	page(w, "WhatsApp", body)
}

func formatHealthBytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}

	div, exp := uint64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

// VOUCHERGO_CLEAN_WA_UI_LABEL_V1

// VOUCHERGO_WA_GATEWAY_LABEL_DYNAMIC_V1
