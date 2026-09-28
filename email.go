package main

import (
	"crypto/tls"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	adminOrderEmailTo   = "suryanasyarif4@gmail.com"
	defaultEmailTimeout = 10 * time.Second
)

func cleanMailHeader(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.TrimSpace(s)
}

func formatRpEmail(n int64) string {
	s := strconv.FormatInt(n, 10)
	out := ""
	for len(s) > 3 {
		out = "." + s[len(s)-3:] + out
		s = s[:len(s)-3]
	}
	return "Rp" + s + out
}

type adminOrderEmailData struct {
	ID                  int64
	OrderID             string
	CustomerName        string
	CustomerPhone       string
	CustomerEmail       string
	HotspotMAC          string
	HotspotIP           string
	CreatedAt           string
	Status              string
	VoucherUsername     string
	PaymentBypassCount  int
	PaymentBypassLastAt string
	PaymentBypassUntil  string
}

func (a *App) saveEmailSettingsFromRequest(r *http.Request) {
	keys := []string{
		"email_enabled",
		"smtp_host",
		"smtp_port",
		"smtp_username",
		"smtp_password",
		"smtp_from_email",
		"smtp_from_name",
		"email_test_to",
	}
	for _, k := range keys {
		_ = a.setSetting(k, r.FormValue(k))
	}
}

func (a *App) emailTestAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = r.ParseForm()
	a.saveEmailSettingsFromRequest(r)

	to := strings.TrimSpace(r.FormValue("email_test_to"))
	if to == "" {
		to = strings.TrimSpace(a.setting("email_test_to"))
	}
	if to == "" {
		to = strings.TrimSpace(a.setting("smtp_username"))
	}
	if to == "" || !strings.Contains(to, "@") {
		http.Redirect(w, r, "/admin/settings?email_test=failed&email_test_error="+url.QueryEscape("Email tujuan test belum valid"), http.StatusSeeOther)
		return
	}

	if err := a.sendEmailDirect(to, "Test Email VoucherGo", `<div style="font-family:Arial,sans-serif;padding:18px">
	<h2 style="margin:0 0 10px;color:#0f172a">Test Email VoucherGo Berhasil</h2>
	<p style="margin:0;color:#334155;line-height:1.6">SMTP sudah bisa dipakai untuk kirim voucher otomatis.</p>
</div>`); err != nil {
		log.Printf("email test failed to=%s err=%v", to, err)
		http.Redirect(w, r, "/admin/settings?email_test=failed&email_test_error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/settings?email_test=ok", http.StatusSeeOther)
}

func (a *App) sendVoucherEmail(txID int64) {
	if strings.TrimSpace(a.setting("email_enabled")) != "1" {
		return
	}

	var orderID, email, name, phone, user, pass, pkg, router string
	var amount int64

	err := a.db.QueryRow(`SELECT
			t.order_id,
			COALESCE(t.customer_email,''),
			COALESCE(t.customer_name,''),
			COALESCE(t.customer_phone,''),
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,''),
			p.name,
			r.name,
			t.amount
		FROM transactions t
		JOIN packages p ON p.id=t.package_id
		JOIN routers r ON r.id=t.router_id
		WHERE t.id=?`, txID).Scan(&orderID, &email, &name, &phone, &user, &pass, &pkg, &router, &amount)

	if err != nil {
		log.Printf("email voucher query failed tx=%d err=%v", txID, err)
		return
	}

	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		log.Printf("email voucher skipped tx=%d order=%s reason=empty email phone=%s", txID, orderID, phone)
		return
	}
	if user == "" {
		log.Printf("email voucher skipped tx=%d order=%s reason=empty voucher", txID, orderID)
		return
	}
	if pass == "" {
		pass = user
	}
	if name == "" {
		name = "Pelanggan"
	}

	subject := "Voucher WiFi " + router + " - " + orderID

	nameHTML := html.EscapeString(name)
	orderHTML := html.EscapeString(orderID)
	routerHTML := html.EscapeString(router)
	pkgHTML := html.EscapeString(pkg)
	totalHTML := html.EscapeString(formatRpEmail(amount))
	voucherHTML := html.EscapeString(user)

	body := fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#f1f5f9;font-family:Arial,sans-serif;color:#0f172a">
	<div style="max-width:560px;margin:0 auto;padding:22px 14px">
		<div style="background:linear-gradient(135deg,#2563eb,#7c3aed 58%%,#06b6d4);border-radius:24px 24px 0 0;padding:24px 22px;color:#ffffff">
			<div style="font-size:13px;font-weight:800;letter-spacing:.08em;text-transform:uppercase;opacity:.86">VoucherGo</div>
			<h1 style="margin:8px 0 0;font-size:26px;line-height:1.15;letter-spacing:-.04em">Voucher WiFi Berhasil Dibeli</h1>
			<p style="margin:8px 0 0;font-size:14px;line-height:1.5;opacity:.92">Pembayaran berhasil, voucher siap digunakan.</p>
		</div>

		<div style="background:#ffffff;border:1px solid #e2e8f0;border-top:0;border-radius:0 0 24px 24px;padding:22px;box-shadow:0 18px 45px rgba(15,23,42,.10)">
			<p style="margin:0 0 14px;font-size:15px;line-height:1.6">Halo <b>%s</b>,</p>

			<div style="border:1px solid #dbeafe;background:#eff6ff;border-radius:18px;padding:18px;margin:0 0 18px;text-align:center">
				<div style="font-size:12px;font-weight:900;letter-spacing:.12em;color:#2563eb;text-transform:uppercase;margin-bottom:8px">Kode Voucher</div>
				<div style="font-family:Consolas,Menlo,monospace;font-size:34px;line-height:1;font-weight:900;letter-spacing:.12em;color:#0f172a">%s</div>
			</div>

			<table role="presentation" style="width:100%%;border-collapse:collapse;margin:0 0 18px">
				<tr>
					<td style="padding:9px 0;color:#64748b;font-size:13px;border-bottom:1px solid #e5e7eb">Order ID</td>
					<td style="padding:9px 0;text-align:right;font-weight:800;font-size:13px;border-bottom:1px solid #e5e7eb">%s</td>
				</tr>
				<tr>
					<td style="padding:9px 0;color:#64748b;font-size:13px;border-bottom:1px solid #e5e7eb">Router</td>
					<td style="padding:9px 0;text-align:right;font-weight:800;font-size:13px;border-bottom:1px solid #e5e7eb">%s</td>
				</tr>
				<tr>
					<td style="padding:9px 0;color:#64748b;font-size:13px;border-bottom:1px solid #e5e7eb">Paket</td>
					<td style="padding:9px 0;text-align:right;font-weight:800;font-size:13px;border-bottom:1px solid #e5e7eb">%s</td>
				</tr>
				<tr>
					<td style="padding:9px 0;color:#64748b;font-size:13px">Total</td>
					<td style="padding:9px 0;text-align:right;font-weight:900;font-size:14px">%s</td>
				</tr>
			</table>

			<div style="background:#f8fafc;border:1px solid #e2e8f0;border-radius:16px;padding:14px;color:#334155;font-size:14px;line-height:1.55">
				Salin <b>Kode Voucher</b> di atas, lalu masukkan di halaman login WiFi.
			</div>
					%s

			<p style="margin:18px 0 0;color:#64748b;font-size:12px;line-height:1.5">
				Email ini dikirim otomatis oleh VoucherGo.
			</p>
		</div>
	</div>
</body>
</html>`, nameHTML, voucherHTML, orderHTML, routerHTML, pkgHTML, totalHTML,
		voucherEmailPhoneHintHTML(phone))

	if err := a.sendEmailDirect(email, subject, body); err != nil {
		log.Printf("email voucher failed tx=%d order=%s to=%s err=%v", txID, orderID, email, err)
		return
	}

	log.Printf("email voucher sent tx=%d order=%s to=%s", txID, orderID, email)
}

func (a *App) sendEmailDirect(to, subject, body string) error {
	return a.sendEmailWithContentType(to, subject, body, "text/html; charset=UTF-8", defaultEmailTimeout)
}

func (a *App) sendEmailPlain(to, subject, body string) error {
	return a.sendEmailWithContentType(to, subject, body, "text/plain; charset=UTF-8", defaultEmailTimeout)
}

func (a *App) sendEmailWithContentType(to, subject, body, contentType string, timeout time.Duration) error {
	host := strings.TrimSpace(a.setting("smtp_host"))
	port := strings.TrimSpace(a.setting("smtp_port"))
	user := strings.TrimSpace(a.setting("smtp_username"))
	pass := strings.TrimSpace(a.setting("smtp_password"))
	fromEmail := strings.TrimSpace(a.setting("smtp_from_email"))
	fromName := strings.TrimSpace(a.setting("smtp_from_name"))

	if host == "" {
		return fmt.Errorf("SMTP host belum diisi")
	}
	if port == "" {
		port = "587"
	}
	if fromEmail == "" {
		fromEmail = user
	}
	if fromEmail == "" {
		return fmt.Errorf("SMTP from email belum diisi")
	}
	if fromName == "" {
		fromName = a.setting("site_name")
	}
	if fromName == "" {
		fromName = "VoucherGo"
	}

	to = cleanMailHeader(to)
	subject = cleanMailHeader(subject)
	fromEmail = cleanMailHeader(fromEmail)
	fromName = cleanMailHeader(fromName)

	if to == "" || !strings.Contains(to, "@") {
		return fmt.Errorf("email tujuan tidak valid")
	}
	if timeout <= 0 {
		timeout = defaultEmailTimeout
	}

	addr := net.JoinHostPort(host, port)
	fromHeader := fmt.Sprintf("%s <%s>", fromName, fromEmail)

	var msg strings.Builder
	msg.WriteString("From: " + fromHeader + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: " + contentType + "\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	dialer := &net.Dialer{Timeout: timeout}
	var conn net.Conn
	var err error
	tlsConfig := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
	if port == "465" {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if port != "465" {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsConfig); err != nil {
				return err
			}
		}
	}

	if user != "" {
		if ok, _ := c.Extension("AUTH"); !ok {
			return fmt.Errorf("SMTP auth tidak didukung server")
		}
		if err := c.Auth(smtp.PlainAuth("", user, pass, host)); err != nil {
			return err
		}
	}
	if err := c.Mail(fromEmail); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := wc.Write([]byte(msg.String())); err != nil {
		_ = wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return c.Quit()
}

func (a *App) sendAdminOrderEmail(txID int64, fallbackOrderID string) {
	orderID := fallbackOrderID
	var d adminOrderEmailData

	err := a.db.QueryRow(`SELECT
			t.id,
			t.order_id,
			COALESCE(t.customer_name,''),
			COALESCE(t.customer_phone,''),
			COALESCE(t.customer_email,''),
			COALESCE(t.hotspot_mac,''),
			COALESCE(t.hotspot_ip,''),
			DATE_FORMAT(t.created_at,'%Y-%m-%d %H:%i:%s'),
			t.status,
			COALESCE(t.voucher_username,''),
			COALESCE(t.payment_bypass_count,0),
			COALESCE(DATE_FORMAT(t.payment_bypass_last_at,'%Y-%m-%d %H:%i:%s'),''),
			COALESCE(DATE_FORMAT(t.payment_bypass_until,'%Y-%m-%d %H:%i:%s'),'')
		FROM transactions t
		WHERE t.id=?`, txID).Scan(
		&d.ID,
		&d.OrderID,
		&d.CustomerName,
		&d.CustomerPhone,
		&d.CustomerEmail,
		&d.HotspotMAC,
		&d.HotspotIP,
		&d.CreatedAt,
		&d.Status,
		&d.VoucherUsername,
		&d.PaymentBypassCount,
		&d.PaymentBypassLastAt,
		&d.PaymentBypassUntil,
	)
	if d.OrderID != "" {
		orderID = d.OrderID
	}
	if err != nil {
		log.Printf("admin order email failed order_id=%s", orderID)
		return
	}

	subject := fmt.Sprintf("[VoucherGo] Order Baru %s", d.OrderID)
	body := fmt.Sprintf(`Nama: %s
No. HP: %s
Email: %s
MAC Address : %s
IP Hotspot  : %s
Waktu order : %s
ID transaksi: %d
Status      : %s
Voucher     : %s
Bypass      : %d
Bypass last : %s
Bypass until: %s
`, d.CustomerName, d.CustomerPhone, d.CustomerEmail, d.HotspotMAC, d.HotspotIP, d.CreatedAt, d.ID, d.Status, d.VoucherUsername, d.PaymentBypassCount, d.PaymentBypassLastAt, d.PaymentBypassUntil)

	if err := a.sendEmailPlain(adminOrderEmailTo, subject, body); err != nil {
		log.Printf("admin order email failed order_id=%s", d.OrderID)
		return
	}

	log.Printf("admin order email sent order_id=%s", d.OrderID)
}

func voucherEmailPhoneHintHTML(phone string) string {
	// VOUCHERGO_EMAIL_PHONE_HINT_V1
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}
	return fmt.Sprintf(`<div style="margin:12px 0 0;text-align:center;font-size:14px;line-height:1.55;color:#cbd5e1">No. WhatsApp: <b style="color:#f8fafc">%s</b></div>`, html.EscapeString(phone))
}
