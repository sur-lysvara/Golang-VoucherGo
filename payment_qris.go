package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

type pakasirQRISData struct {
	TxnID         string
	QRISString    string
	TotalPayment  int64
	ExpiredAt     string
	PaymentMethod string
}

func appPublicBaseURL(a *App) string {
	base := strings.TrimSpace(a.setting("public_base_url"))
	if base == "" {
		base = "https://vouchergo.biz.id"
	}
	return strings.TrimRight(base, "/")
}

func (a *App) invoiceURL(orderID string) string {
	return appPublicBaseURL(a) + "/invoice/" + url.PathEscape(orderID)
}

func (a *App) qrisImageURL(orderID string) string {
	return appPublicBaseURL(a) + "/qris/" + url.PathEscape(orderID) + ".png"
}

func (a *App) setPendingInvoiceWAStatus(txID int64, status string, errMsg string) bool {
	if len(errMsg) > 1200 {
		errMsg = errMsg[:1200]
	}

	var (
		res sql.Result
		err error
	)

	switch status {
	case "sent":
		res, err = a.db.Exec(`UPDATE transactions SET wa_status='sent', wa_error=NULL, wa_sent_at=NOW() WHERE id=? AND status='pending'`, txID)
	default:
		res, err = a.db.Exec(`UPDATE transactions SET wa_status='failed', wa_error=?, wa_sent_at=NULL WHERE id=? AND status='pending' AND COALESCE(wa_status,'') <> 'sent'`, errMsg, txID)
	}

	if err != nil {
		log.Printf("unpaid wa status update failed tx=%d status=%s err=%v", txID, status, err)
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("unpaid wa status ignored tx=%d status=%s reason=transaction no longer pending or already sent", txID, status)
		return false
	}

	return true
}

func (a *App) createPakasirQRIS(orderID string, amount int64) (*pakasirQRISData, error) {
	project := strings.TrimSpace(a.setting("pakasir_slug"))
	apiKey := strings.TrimSpace(a.setting("pakasir_api_key"))
	if project == "" || apiKey == "" {
		return nil, fmt.Errorf("pakasir_slug/api_key kosong")
	}

	payload := map[string]any{
		"method": "qris",
		"amount": amount,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal pakasir v2 request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	endpoint := "https://app.pakasir.com/api/v2/create-transaction/" +
		url.PathEscape(project) + "/" + url.PathEscape(orderID)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := (&http.Client{Timeout: 25 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rb, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("pakasir v2 HTTP %d: %s", resp.StatusCode, string(rb))
	}

	var raw any
	if err := json.Unmarshal(rb, &raw); err != nil {
		return nil, fmt.Errorf("decode pakasir v2 response: %w body=%s", err, string(rb))
	}

	txnID := jsonFindString(raw, "txn_id")
	if txnID == "" {
		return nil, fmt.Errorf("txn_id tidak ditemukan di response Pakasir V2: %s", string(rb))
	}

	qris := jsonFindString(raw, "qr_string")
	if qris == "" {
		return nil, fmt.Errorf("qr_string tidak ditemukan di response Pakasir V2: %s", string(rb))
	}

	total := jsonFindInt(raw, "total_payment")
	if total <= 0 {
		total = amount
	}

	expired := jsonFindString(raw, "expired_at")
	method := jsonFindString(raw, "payment_method")
	if method == "" {
		method = "qris"
	}

	return &pakasirQRISData{
		TxnID:         txnID,
		QRISString:    qris,
		TotalPayment:  total,
		ExpiredAt:     expired,
		PaymentMethod: method,
	}, nil
}
func jsonFindString(v any, names ...string) string {
	wanted := map[string]bool{}
	for _, n := range names {
		wanted[strings.ToLower(n)] = true
	}

	var walk func(any) string
	walk = func(x any) string {
		switch t := x.(type) {
		case map[string]any:
			for k, v := range t {
				if wanted[strings.ToLower(k)] {
					switch vv := v.(type) {
					case string:
						return strings.TrimSpace(vv)
					case float64:
						return strconv.FormatInt(int64(vv), 10)
					}
				}
			}
			for _, v := range t {
				if got := walk(v); got != "" {
					return got
				}
			}
		case []any:
			for _, v := range t {
				if got := walk(v); got != "" {
					return got
				}
			}
		}
		return ""
	}

	return walk(v)
}

func jsonFindInt(v any, names ...string) int64 {
	s := jsonFindString(v, names...)
	if s != "" {
		s = strings.ReplaceAll(s, ".", "")
		s = strings.ReplaceAll(s, ",", "")
		n, _ := strconv.ParseInt(s, 10, 64)
		if n > 0 {
			return n
		}
	}

	wanted := map[string]bool{}
	for _, n := range names {
		wanted[strings.ToLower(n)] = true
	}

	var walk func(any) int64
	walk = func(x any) int64 {
		switch t := x.(type) {
		case map[string]any:
			for k, v := range t {
				if wanted[strings.ToLower(k)] {
					switch vv := v.(type) {
					case float64:
						return int64(vv)
					case int64:
						return vv
					case int:
						return int64(vv)
					}
				}
			}
			for _, v := range t {
				if got := walk(v); got > 0 {
					return got
				}
			}
		case []any:
			for _, v := range t {
				if got := walk(v); got > 0 {
					return got
				}
			}
		}
		return 0
	}

	return walk(v)
}

func (a *App) qrisPNG(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/qris/")
	orderID = strings.TrimSuffix(orderID, ".png")
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		http.NotFound(w, r)
		return
	}

	var qris string
	err := a.db.QueryRow(`SELECT COALESCE(qris_string,'') FROM transactions WHERE order_id=?`, orderID).Scan(&qris)
	if err != nil || strings.TrimSpace(qris) == "" {
		http.NotFound(w, r)
		return
	}

	png, err := qrcode.Encode(qris, qrcode.Medium, 640)
	if err != nil {
		http.Error(w, "failed generate qris", 500)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(png)
}

func (a *App) paymentBypassStartJSON(w http.ResponseWriter, r *http.Request) {
	writeJSON := func(code int, payload map[string]any) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(payload)
	}

	orderID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/payment-bypass-start/"))
	if orderID == "" {
		writeJSON(http.StatusBadRequest, map[string]any{"ok": false, "error": "order_id kosong"})
		return
	}

	var ro Router
	var status, mac, ip string

	err := a.db.QueryRow(`SELECT
			r.id,
			r.name,
			r.slug,
			COALESCE(r.hotspot_login_url,''),
			COALESCE(r.host,''),
			r.port,
			COALESCE(r.api_mode,'auto'),
			COALESCE(r.username,''),
			COALESCE(r.password,''),
			r.use_ssl,
			r.insecure_ssl,
			r.enabled,
			t.status,
			COALESCE(t.hotspot_mac,''),
			COALESCE(t.hotspot_ip,'')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		WHERE t.order_id=?`, orderID).
		Scan(
			&ro.ID,
			&ro.Name,
			&ro.Slug,
			&ro.HotspotLoginURL,
			&ro.Host,
			&ro.Port,
			&ro.APIMode,
			&ro.Username,
			&ro.Password,
			&ro.UseSSL,
			&ro.InsecureSSL,
			&ro.Enabled,
			&status,
			&mac,
			&ip,
		)

	if err != nil {
		writeJSON(http.StatusNotFound, map[string]any{"ok": false, "error": "invoice tidak ditemukan"})
		return
	}

	if status == "paid" {
		writeJSON(http.StatusOK, map[string]any{
			"ok":      true,
			"skipped": "already_paid",
			"limit":   paymentBypassLimit,
		})
		return
	}

	mac = cleanPaymentClientMAC(mac)
	ip = cleanPaymentClientIP(ip)

	if mac == "" {
		writeJSON(http.StatusBadRequest, map[string]any{
			"ok":    false,
			"error": "hotspot_mac kosong",
			"limit": paymentBypassLimit,
		})
		return
	}

	used, remaining, alreadyActive, allowed, err := a.reservePaymentBypassAttempt(orderID)
	if err != nil {
		writeJSON(http.StatusInternalServerError, map[string]any{
			"ok":    false,
			"error": "gagal cek status bypass",
			"limit": paymentBypassLimit,
		})
		return
	}

	if alreadyActive {
		writeJSON(http.StatusOK, map[string]any{
			"ok":                true,
			"already_active":    true,
			"remaining_seconds": remaining,
			"used":              used,
			"limit":             paymentBypassLimit,
		})
		return
	}

	if !allowed {
		writeJSON(http.StatusTooManyRequests, map[string]any{
			"ok":    false,
			"error": "Jatah bypass sudah habis. Bypass maksimal 2 kali per invoice.",
			"used":  used,
			"limit": paymentBypassLimit,
		})
		return
	}

	if err := a.startPaymentBypass90(ro, orderID, mac, ip); err != nil {
		a.rollbackPaymentBypassAttempt(orderID, used)

		log.Printf("payment qris bypass gagal order=%s router=%s mac=%s ip=%s: %v", orderID, ro.Slug, mac, ip, err)
		writeJSON(http.StatusInternalServerError, map[string]any{
			"ok":    false,
			"error": "bypass failed",
			"used":  used - 1,
			"limit": paymentBypassLimit,
		})
		return
	}

	log.Printf("payment qris bypass OK order=%s router=%s mac=%s ip=%s seconds=%d used=%d/%d", orderID, ro.Slug, mac, ip, paymentBypassSeconds, used, paymentBypassLimit)

	writeJSON(http.StatusOK, map[string]any{
		"ok":                true,
		"remaining_seconds": paymentBypassSeconds,
		"used":              used,
		"limit":             paymentBypassLimit,
	})
}

func (a *App) sendUnpaidPaymentWA(txID int64) {
	var orderID, phone, name, pkg, router, expired, txStatus string
	var total int64

	err := a.db.QueryRow(`SELECT
			t.order_id,
			COALESCE(t.customer_phone,''),
			COALESCE(t.customer_name,''),
			p.name,
			r.name,
			COALESCE(t.qris_total_payment, t.amount),
			COALESCE(t.qris_expired_at,''),
			COALESCE(t.status,'')
		FROM transactions t
		JOIN packages p ON p.id=t.package_id
		JOIN routers r ON r.id=t.router_id
		WHERE t.id=?`, txID).
		Scan(&orderID, &phone, &name, &pkg, &router, &total, &expired, &txStatus)

	if err != nil {
		log.Printf("unpaid wa query failed tx=%d err=%v", txID, err)
		return
	}

	if strings.ToLower(strings.TrimSpace(txStatus)) != "pending" {
		log.Printf("unpaid wa skipped tx=%d order=%s reason=status %s", txID, orderID, txStatus)
		return
	}

	phone = normalizePhoneID(phone)
	if phone == "" {
		log.Printf("unpaid wa skipped tx=%d order=%s reason=empty phone", txID, orderID)
		return
	}

	if name == "" {
		name = "Pelanggan"
	}

	caption := fmt.Sprintf(
		"Invoice VoucherGo - Belum Dibayar\n\nHalo %s,\nPaket: %s\nRouter: %s\nTotal: %s\nOrder ID: %s\n\nSilakan scan QRIS pada gambar ini untuk menyelesaikan pembayaran.\n\nKoneksi perangkat Anda dibypass selama %d detik untuk menyelesaikan pembayaran.\n\nLink invoice:\n%s",
		name,
		pkg,
		router,
		formatRpEmail(total),
		orderID,
		paymentBypassSeconds,
		a.invoiceURL(orderID),
	)
	if expired != "" {
		caption += "\n\nBerlaku sampai:\n" + expired
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	start := time.Now()
	code, body, err := a.localWACall(ctx, http.MethodPost, "/send-image", map[string]string{
		"number":    phone,
		"caption":   caption,
		"image_url": a.qrisImageURL(orderID),
	})
	dur := time.Since(start)
	if err != nil {
		errMsg := err.Error()
		a.setPendingInvoiceWAStatus(txID, "failed", errMsg)
		log.Printf("unpaid wa image failed tx=%d order=%s phone=%s duration=%s err=%v", txID, orderID, phone, dur, err)
		return
	}
	if code < 200 || code >= 300 {
		errMsg := fmt.Sprintf("HTTP %d: %s", code, truncateLocalWAString(body, 500))
		a.setPendingInvoiceWAStatus(txID, "failed", errMsg)
		log.Printf("unpaid wa image failed tx=%d order=%s phone=%s http=%d duration=%s body=%s", txID, orderID, phone, code, dur, truncateLocalWAString(body, 500))
		return
	}

	// VOUCHERGO_UNPAID_WA_STATUS_UPDATE_V1
	a.setPendingInvoiceWAStatus(txID, "sent", "")
	log.Printf("unpaid wa image sent tx=%d order=%s phone=%s duration=%s", txID, orderID, phone, dur)
}

func (a *App) sendUnpaidPaymentEmail(txID int64) {
	if strings.TrimSpace(a.setting("email_enabled")) != "1" {
		return
	}

	var orderID, email, name, pkg, router, expired string
	var total int64

	err := a.db.QueryRow(`SELECT
			t.order_id,
			COALESCE(t.customer_email,''),
			COALESCE(t.customer_name,''),
			p.name,
			r.name,
			COALESCE(t.qris_total_payment, t.amount),
			COALESCE(t.qris_expired_at,'')
		FROM transactions t
		JOIN packages p ON p.id=t.package_id
		JOIN routers r ON r.id=t.router_id
		WHERE t.id=?`, txID).
		Scan(&orderID, &email, &name, &pkg, &router, &total, &expired)

	if err != nil {
		log.Printf("unpaid email query failed tx=%d err=%v", txID, err)
		return
	}

	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		return
	}
	if name == "" {
		name = "Pelanggan"
	}

	subject := "Invoice QRIS VoucherGo - " + orderID
	body := fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#f1f5f9;font-family:Arial,sans-serif;color:#0f172a">
	<div style="max-width:560px;margin:0 auto;padding:22px 14px">
		<div style="background:#ffffff;border:1px solid #e2e8f0;border-radius:24px;padding:22px;text-align:center">
			<h2 style="margin:0 0 8px">Invoice Belum Dibayar</h2>
			<p style="margin:0 0 14px;color:#475569">Halo <b>%s</b>, silakan scan QRIS berikut.</p>
			<img src="%s" alt="QRIS Pembayaran" style="width:280px;max-width:100%%;border:1px solid #e5e7eb;border-radius:18px;padding:12px;background:#fff">
			<table role="presentation" style="width:100%%;border-collapse:collapse;margin:18px 0;text-align:left">
				<tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;color:#64748b">Order ID</td><td style="padding:8px;border-bottom:1px solid #e5e7eb;text-align:right;font-weight:700">%s</td></tr>
				<tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;color:#64748b">Router</td><td style="padding:8px;border-bottom:1px solid #e5e7eb;text-align:right;font-weight:700">%s</td></tr>
				<tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;color:#64748b">Paket</td><td style="padding:8px;border-bottom:1px solid #e5e7eb;text-align:right;font-weight:700">%s</td></tr>
				<tr><td style="padding:8px;color:#64748b">Total</td><td style="padding:8px;text-align:right;font-weight:900">%s</td></tr>
			</table>
			<div style="background:#fef9c3;border:1px solid #fde68a;border-radius:16px;padding:13px;color:#713f12;font-size:14px;line-height:1.5">
				Koneksi perangkat Anda dibypass selama <b>%d detik</b> untuk menyelesaikan pembayaran.
			</div>
			<p style="margin:16px 0 0"><a href="%s" style="display:inline-block;background:#2563eb;color:#fff;text-decoration:none;border-radius:999px;padding:12px 18px;font-weight:700">Buka Invoice</a></p>
			%s
		</div>
	</div>
</body>
</html>`,
		html.EscapeString(name),
		html.EscapeString(a.qrisImageURL(orderID)),
		html.EscapeString(orderID),
		html.EscapeString(router),
		html.EscapeString(pkg),
		html.EscapeString(formatRpEmail(total)),
		paymentBypassSeconds,
		html.EscapeString(a.invoiceURL(orderID)),
		func() string {
			if expired == "" {
				return ""
			}
			return `<p style="color:#64748b;font-size:12px">Berlaku sampai: ` + html.EscapeString(expired) + `</p>`
		}(),
	)

	if err := a.sendEmailDirect(email, subject, body); err != nil {
		log.Printf("unpaid email failed tx=%d order=%s to=%s err=%v", txID, orderID, email, err)
		return
	}

	log.Printf("unpaid email sent tx=%d order=%s to=%s", txID, orderID, email)
}
