package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type unusedVoucherCandidate struct {
	ID              int64
	OrderID         string
	CustomerName    string
	CustomerPhone   string
	CustomerEmail   string
	VoucherUsername string
	PackageName     string
	RouterName      string
	Router          Router
}

func (a *App) startUnusedVoucherReminderWorker() {
	log.Printf("unused voucher reminder worker started")

	timer := time.NewTimer(90 * time.Second)
	defer timer.Stop()

	for {
		<-timer.C

		if err := a.processUnusedVoucherReminders(); err != nil {
			log.Printf("unused voucher reminder worker error: %v", err)
		}

		timer.Reset(5 * time.Minute)
	}
}

func (a *App) processUnusedVoucherReminders() error {
	if strings.TrimSpace(a.setting("unused_voucher_reminder_enabled")) != "1" {
		return nil
	}

	afterMinutes := unusedReminderInt(a.setting("unused_voucher_reminder_minutes"), 60)
	if afterMinutes < 5 {
		afterMinutes = 5
	}
	if afterMinutes > 1440 {
		afterMinutes = 1440
	}

	limit := unusedReminderInt(a.setting("unused_voucher_reminder_limit"), 25)
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	rows, err := a.db.Query(`
		SELECT
			t.id,
			t.order_id,
			COALESCE(t.customer_name,''),
			COALESCE(t.customer_phone,''),
			COALESCE(t.customer_email,''),
			COALESCE(t.voucher_username,''),
			COALESCE(p.name,''),
			COALESCE(r.name,''),
			r.id,
			COALESCE(r.name,''),
			COALESCE(r.host,''),
			COALESCE(r.port,0),
			COALESCE(r.api_mode,'auto'),
			COALESCE(r.username,''),
			COALESCE(r.password,''),
			COALESCE(r.use_ssl,0),
			COALESCE(r.insecure_ssl,0)
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		LEFT JOIN packages p ON p.id=t.package_id
		WHERE t.status='paid'
		  AND COALESCE(t.provision_status,'')='success'
		  AND COALESCE(t.voucher_username,'') <> ''
		  AND t.unused_reminder_sent_at IS NULL
		  AND t.first_used_at IS NULL
		  AND TIMESTAMPDIFF(MINUTE, t.created_at, NOW()) >= ?
		ORDER BY t.id ASC
		LIMIT ?`, afterMinutes, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var c unusedVoucherCandidate
		var useSSL, insecureSSL bool

		err := rows.Scan(
			&c.ID,
			&c.OrderID,
			&c.CustomerName,
			&c.CustomerPhone,
			&c.CustomerEmail,
			&c.VoucherUsername,
			&c.PackageName,
			&c.RouterName,
			&c.Router.ID,
			&c.Router.Name,
			&c.Router.Host,
			&c.Router.Port,
			&c.Router.APIMode,
			&c.Router.Username,
			&c.Router.Password,
			&useSSL,
			&insecureSSL,
		)
		if err != nil {
			log.Printf("unused voucher scan failed: %v", err)
			continue
		}

		c.Router.UseSSL = useSSL
		c.Router.InsecureSSL = insecureSSL

		_, _ = a.db.Exec("UPDATE transactions SET last_usage_check_at=NOW() WHERE id=?", c.ID)

		unused, err := unusedVoucherCheckRouter(c.Router, c.VoucherUsername)
		if err != nil {
			log.Printf("unused voucher check failed tx=%d order=%s router=%s voucher=%s err=%v", c.ID, c.OrderID, c.RouterName, c.VoucherUsername, err)
			continue
		}

		if !unused {
			_, _ = a.db.Exec("UPDATE transactions SET first_used_at=COALESCE(first_used_at,NOW()) WHERE id=?", c.ID)
			continue
		}

		sent, failed := a.sendUnusedVoucherReminder(c)
		if len(sent) > 0 {
			channel := strings.Join(sent, ",")
			_, _ = a.db.Exec("UPDATE transactions SET unused_reminder_sent_at=NOW(), unused_reminder_channel=? WHERE id=?", channel, c.ID)
			log.Printf("unused voucher reminder sent tx=%d order=%s channel=%s", c.ID, c.OrderID, channel)
			continue
		}

		if len(failed) == 0 {
			_, _ = a.db.Exec("UPDATE transactions SET unused_reminder_sent_at=NOW(), unused_reminder_channel='skipped:no-contact' WHERE id=?", c.ID)
			log.Printf("unused voucher reminder skipped no contact tx=%d order=%s", c.ID, c.OrderID)
		} else {
			log.Printf("unused voucher reminder failed tx=%d order=%s failed=%s", c.ID, c.OrderID, strings.Join(failed, ","))
		}
	}

	return rows.Err()
}

func unusedReminderInt(v string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return def
	}
	return n
}

func (a *App) sendUnusedVoucherReminder(c unusedVoucherCandidate) ([]string, []string) {
	channels := unusedReminderChannels(a.setting("unused_voucher_reminder_channel"))

	var sent []string
	var failed []string

	for _, ch := range channels {
		switch ch {
		case "wa":
			phone := unusedReminderNormalizePhone(c.CustomerPhone)
			if phone == "" {
				continue
			}
			if err := a.sendUnusedVoucherReminderWA(phone, a.unusedVoucherReminderText(c)); err != nil {
				failed = append(failed, "wa:"+err.Error())
			} else {
				sent = append(sent, "wa")
			}

		case "email":
			email := strings.TrimSpace(c.CustomerEmail)
			if email == "" || !strings.Contains(email, "@") {
				continue
			}
			if strings.TrimSpace(a.setting("email_enabled")) != "1" {
				continue
			}
			if err := a.sendEmailDirect(email, "Voucher WiFi belum digunakan", a.unusedVoucherReminderHTML(c)); err != nil {
				failed = append(failed, "email:"+err.Error())
			} else {
				sent = append(sent, "email")
			}
		}
	}

	return sent, failed
}

func unusedReminderChannels(v string) []string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		v = "wa"
	}

	v = strings.ReplaceAll(v, "both", "wa,email")
	v = strings.ReplaceAll(v, "+", ",")
	v = strings.ReplaceAll(v, ";", ",")

	seen := map[string]bool{}
	var out []string

	for _, p := range strings.Split(v, ",") {
		p = strings.TrimSpace(p)
		if p == "whatsapp" {
			p = "wa"
		}
		if (p == "wa" || p == "email") && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}

	if len(out) == 0 {
		out = append(out, "wa")
	}

	return out
}

func unusedReminderNormalizePhone(v string) string {
	var b strings.Builder

	for _, r := range v {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	s := b.String()
	if s == "" {
		return ""
	}

	if strings.HasPrefix(s, "0") {
		s = "62" + strings.TrimPrefix(s, "0")
	} else if strings.HasPrefix(s, "8") {
		s = "62" + s
	}

	if len(s) < 10 {
		return ""
	}

	return s
}

func (a *App) unusedVoucherReminderText(c unusedVoucherCandidate) string {
	name := strings.TrimSpace(c.CustomerName)
	if name == "" {
		name = "Pelanggan"
	}

	pkg := strings.TrimSpace(c.PackageName)
	if pkg == "" {
		pkg = "-"
	}

	return fmt.Sprintf("Halo %s,\n\nVoucher WiFi kamu sudah aktif, tapi belum terlihat digunakan.\n\nPaket: %s\nKode Voucher: %s\n\nSilakan sambungkan ke WiFi lalu login memakai kode voucher tersebut.\n\nJika voucher sudah digunakan, abaikan pesan ini.\n\nVoucherGo", name, pkg, c.VoucherUsername)
}

func (a *App) unusedVoucherReminderHTML(c unusedVoucherCandidate) string {
	name := strings.TrimSpace(c.CustomerName)
	if name == "" {
		name = "Pelanggan"
	}

	pkg := strings.TrimSpace(c.PackageName)
	if pkg == "" {
		pkg = "-"
	}

	name = html.EscapeString(name)
	pkg = html.EscapeString(pkg)
	code := html.EscapeString(c.VoucherUsername)

	return fmt.Sprintf("<div style=\"font-family:Arial,sans-serif;padding:18px;color:#0f172a\"><h2 style=\"margin:0 0 12px\">Voucher WiFi belum digunakan</h2><p>Halo <b>%s</b>,</p><p>Voucher WiFi kamu sudah aktif, tapi belum terlihat digunakan.</p><div style=\"background:#f8fafc;border:1px solid #e2e8f0;border-radius:12px;padding:14px;margin:14px 0\"><div>Paket: <b>%s</b></div><div>Kode Voucher: <b style=\"font-size:20px\">%s</b></div></div><p>Silakan sambungkan ke WiFi lalu login memakai kode voucher tersebut.</p><p style=\"color:#64748b;font-size:13px\">Jika voucher sudah digunakan, abaikan email ini.</p><p style=\"color:#64748b;font-size:12px;margin-top:20px\">Email ini dikirim otomatis oleh VoucherGo.</p></div>", name, pkg, code)
}

func (a *App) sendUnusedVoucherReminderWA(phone, message string) error {
	endpoint := strings.TrimSpace(a.setting("wa_local_url"))
	if endpoint == "" {
		return errors.New("wa_local_url kosong")
	}

	endpoint = strings.TrimRight(endpoint, "/")
	if !strings.HasSuffix(endpoint, "/send") {
		endpoint += "/send"
	}

	// UNUSED_REMINDER_WA_LOCALKEY_V1
	apiKey := strings.TrimSpace(a.localWAKey())

	payload := map[string]string{
		"to":      phone,
		"number":  phone,
		"phone":   phone,
		"message": message,
		"text":    message,
	}

	raw, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 800))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncateLocalWAString(string(body), 300))
	}

	return nil
}

func unusedVoucherCheckRouter(ro Router, voucher string) (bool, error) {
	mode := normalizeRouterAPIMode(ro.APIMode)

	switch mode {
	case "rest":
		return unusedVoucherCheckREST(ro, voucher)
	case "classic":
		return unusedVoucherCheckClassic(ro, voucher, false)
	case "classic_ssl":
		return unusedVoucherCheckClassic(ro, voucher, true)
	default:
		if unused, err := unusedVoucherCheckREST(ro, voucher); err == nil {
			return unused, nil
		}
		return unusedVoucherCheckClassic(ro, voucher, false)
	}
}

func unusedVoucherCheckREST(ro Router, voucher string) (bool, error) {
	host := cleanRouterHost(ro.Host)
	port := ro.Port
	scheme := "http"

	if ro.UseSSL {
		scheme = "https"
		if port <= 0 {
			port = 443
		}
	} else if port <= 0 {
		port = 80
	}

	endpoint := fmt.Sprintf("%s://%s:%d/rest/ip/hotspot/user?name=%s", scheme, host, port, url.QueryEscape(voucher))

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}

	req.SetBasicAuth(ro.Username, ro.Password)

	client := &http.Client{
		Timeout: 12 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: ro.InsecureSSL},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 128*1024))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("REST HTTP %d: %s", resp.StatusCode, truncateLocalWAString(string(raw), 300))
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(raw, &items); err != nil {
		return false, err
	}

	for _, item := range items {
		if strings.TrimSpace(fmt.Sprint(item["name"])) != voucher {
			continue
		}

		uptime := strings.TrimSpace(fmt.Sprint(item["uptime"]))
		return unusedVoucherUptimeIsZero(uptime), nil
	}

	return false, fmt.Errorf("voucher %s tidak ditemukan via REST", voucher)
}

func unusedVoucherUptimeIsZero(v string) bool {
	v = strings.TrimSpace(strings.ToLower(v))

	switch v {
	case "", "0", "0s", "00:00:00", "0:00:00":
		return true
	}

	return false
}

func unusedVoucherCheckClassic(ro Router, voucher string, useTLS bool) (bool, error) {
	conn, err := unusedMTDial(ro, useTLS)
	if err != nil {
		return false, err
	}

	if err := unusedMTLoginNew(conn, ro.Username, ro.Password); err != nil {
		_ = conn.Close()

		conn, err = unusedMTDial(ro, useTLS)
		if err != nil {
			return false, err
		}

		if err := unusedMTLoginLegacy(conn, ro.Username, ro.Password); err != nil {
			_ = conn.Close()
			return false, err
		}
	}

	defer conn.Close()

	if err := unusedMTWriteSentence(conn, []string{
		"/ip/hotspot/user/print",
		"?name=" + voucher,
		"=.proplist=name,uptime",
	}); err != nil {
		return false, err
	}

	records, _, err := unusedMTReadReply(conn)
	if err != nil {
		return false, err
	}

	for _, rec := range records {
		if strings.TrimSpace(rec["name"]) != voucher {
			continue
		}

		return unusedVoucherUptimeIsZero(rec["uptime"]), nil
	}

	return false, fmt.Errorf("voucher %s tidak ditemukan via Classic API", voucher)
}

func unusedMTDial(ro Router, useTLS bool) (net.Conn, error) {
	host := cleanRouterHost(ro.Host)
	port := ro.Port

	if port <= 0 {
		if useTLS {
			port = 8729
		} else {
			port = 8728
		}
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	dialer := &net.Dialer{Timeout: 10 * time.Second}

	if useTLS {
		return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: ro.InsecureSSL,
		})
	}

	return dialer.Dial("tcp", addr)
}

func unusedMTLoginNew(conn net.Conn, username, password string) error {
	if err := unusedMTWriteSentence(conn, []string{"/login", "=name=" + username, "=password=" + password}); err != nil {
		return err
	}

	_, _, err := unusedMTReadReply(conn)
	return err
}

func unusedMTLoginLegacy(conn net.Conn, username, password string) error {
	if err := unusedMTWriteSentence(conn, []string{"/login"}); err != nil {
		return err
	}

	_, ret, err := unusedMTReadReply(conn)
	if err != nil {
		return err
	}

	if ret == "" {
		return errors.New("legacy login challenge kosong")
	}

	challenge, err := hex.DecodeString(ret)
	if err != nil {
		return err
	}

	h := md5.New()
	h.Write([]byte{0})
	h.Write([]byte(password))
	h.Write(challenge)

	response := "00" + hex.EncodeToString(h.Sum(nil))

	if err := unusedMTWriteSentence(conn, []string{"/login", "=name=" + username, "=response=" + response}); err != nil {
		return err
	}

	_, _, err = unusedMTReadReply(conn)
	return err
}

func unusedMTWriteSentence(w io.Writer, words []string) error {
	for _, word := range words {
		if err := unusedMTWriteWord(w, word); err != nil {
			return err
		}
	}

	return unusedMTWriteWord(w, "")
}

func unusedMTWriteWord(w io.Writer, word string) error {
	if err := unusedMTWriteLen(w, len(word)); err != nil {
		return err
	}

	if word == "" {
		return nil
	}

	_, err := io.WriteString(w, word)
	return err
}

func unusedMTWriteLen(w io.Writer, l int) error {
	var b []byte

	switch {
	case l < 0x80:
		b = []byte{byte(l)}
	case l < 0x4000:
		b = []byte{byte((l >> 8) | 0x80), byte(l)}
	case l < 0x200000:
		b = []byte{byte((l >> 16) | 0xC0), byte(l >> 8), byte(l)}
	case l < 0x10000000:
		b = []byte{byte((l >> 24) | 0xE0), byte(l >> 16), byte(l >> 8), byte(l)}
	default:
		b = []byte{0xF0, byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}
	}

	_, err := w.Write(b)
	return err
}

func unusedMTReadReply(r io.Reader) ([]map[string]string, string, error) {
	var records []map[string]string
	var ret string

	for {
		sentence, err := unusedMTReadSentence(r)
		if err != nil {
			return records, ret, err
		}

		if len(sentence) == 0 {
			continue
		}

		switch sentence[0] {
		case "!re":
			records = append(records, unusedMTSentenceMap(sentence))

		case "!done":
			m := unusedMTSentenceMap(sentence)
			if m["ret"] != "" {
				ret = m["ret"]
			}
			return records, ret, nil

		case "!trap", "!fatal":
			m := unusedMTSentenceMap(sentence)
			msg := m["message"]
			if msg == "" {
				msg = strings.Join(sentence, " ")
			}
			return records, ret, errors.New(msg)
		}
	}
}

func unusedMTReadSentence(r io.Reader) ([]string, error) {
	var words []string

	for {
		l, err := unusedMTReadLen(r)
		if err != nil {
			return words, err
		}

		if l == 0 {
			return words, nil
		}

		buf := make([]byte, l)
		if _, err := io.ReadFull(r, buf); err != nil {
			return words, err
		}

		words = append(words, string(buf))
	}
}

func unusedMTReadLen(r io.Reader) (int, error) {
	var b [1]byte

	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}

	c := int(b[0])

	switch {
	case c&0x80 == 0:
		return c, nil

	case c&0xC0 == 0x80:
		var x [1]byte
		if _, err := io.ReadFull(r, x[:]); err != nil {
			return 0, err
		}
		return ((c & ^0xC0) << 8) + int(x[0]), nil

	case c&0xE0 == 0xC0:
		var x [2]byte
		if _, err := io.ReadFull(r, x[:]); err != nil {
			return 0, err
		}
		return ((c & ^0xE0) << 16) + (int(x[0]) << 8) + int(x[1]), nil

	case c&0xF0 == 0xE0:
		var x [3]byte
		if _, err := io.ReadFull(r, x[:]); err != nil {
			return 0, err
		}
		return ((c & ^0xF0) << 24) + (int(x[0]) << 16) + (int(x[1]) << 8) + int(x[2]), nil

	default:
		var x [4]byte
		if _, err := io.ReadFull(r, x[:]); err != nil {
			return 0, err
		}
		return (int(x[0]) << 24) + (int(x[1]) << 16) + (int(x[2]) << 8) + int(x[3]), nil
	}
}

func unusedMTSentenceMap(sentence []string) map[string]string {
	m := map[string]string{}

	for _, word := range sentence[1:] {
		if !strings.HasPrefix(word, "=") {
			continue
		}

		kv := strings.SplitN(strings.TrimPrefix(word, "="), "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}

	return m
}
