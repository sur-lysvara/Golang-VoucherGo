package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func truncateLocalWAString(v string, max int) string {
	if len(v) <= max {
		return v
	}
	return v[:max] + "..."
}

func (a *App) localWABaseURL() string {
	u := strings.TrimSpace(a.setting("wa_local_url"))
	if u == "" {
		u = "http://127.0.0.1:8666/send"
	}

	u = strings.TrimRight(u, "/")
	u = strings.TrimSuffix(u, "/send")

	return u
}

func (a *App) localWAKey() string {
	key := strings.TrimSpace(a.setting("wa_local_key"))
	if key == "" {
		key = "local-wa-key"
	}
	return key
}

func (a *App) localWACall(ctx context.Context, method string, path string, payload any) (int, string, error) {
	base := a.localWABaseURL()

	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, "", err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, base+path, body)
	if err != nil {
		return 0, "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", a.localWAKey())

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := (&http.Client{Timeout: 45 * time.Second}).Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	return resp.StatusCode, strings.TrimSpace(string(b)), nil
}

func maskHealthWARaw(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	var st map[string]any
	if err := json.Unmarshal([]byte(raw), &st); err != nil {
		if len(raw) > 500 {
			return raw[:500] + "..."
		}
		return raw
	}

	if _, ok := st["me"]; ok {
		st["me"] = "[masked]"
	}
	if _, ok := st["lid"]; ok {
		st["lid"] = "[masked]"
	}

	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return raw
	}

	return string(b)
}

func (a *App) setWAStatus(txID int64, status string, errMsg string) {
	if len(errMsg) > 1200 {
		errMsg = errMsg[:1200]
	}

	switch status {
	case "sent":
		_, _ = a.db.Exec(`UPDATE transactions SET wa_status='sent', wa_error=NULL, wa_sent_at=NOW() WHERE id=?`, txID)
	case "sending":
		_, _ = a.db.Exec(`UPDATE transactions SET wa_status='sending', wa_error=NULL, wa_sent_at=NULL WHERE id=?`, txID)
	case "skipped":
		_, _ = a.db.Exec(`UPDATE transactions SET wa_status='skipped', wa_error=?, wa_sent_at=NULL WHERE id=?`, errMsg, txID)
	default:
		_, _ = a.db.Exec(`UPDATE transactions SET wa_status='failed', wa_error=?, wa_sent_at=NULL WHERE id=?`, errMsg, txID)
	}
}

func (a *App) sendWanesia(txID int64) {
	// LOCAL_WA_BAILEYS_V1
	a.setWAStatus(txID, "sending", "")

	var orderID, phone, user, pass, pkg, router string

	err := a.db.QueryRow(`SELECT
			t.order_id,
			COALESCE(t.customer_phone,''),
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,''),
			p.name,
			r.name
		FROM transactions t
		JOIN packages p ON p.id=t.package_id
		JOIN routers r ON r.id=t.router_id
		WHERE t.id=?`, txID).Scan(&orderID, &phone, &user, &pass, &pkg, &router)

	if err != nil {
		msg := "query transaction: " + err.Error()
		log.Printf("local wa failed tx=%d reason=query err=%v", txID, err)
		a.setWAStatus(txID, "failed", msg)
		return
	}

	phone = normalizePhoneID(phone)

	if phone == "" || user == "" {
		msg := "phone or voucher empty"
		log.Printf("local wa skip tx=%d order=%s reason=phone/user empty phone=%q user=%q", txID, orderID, phone, user)
		a.setWAStatus(txID, "failed", msg)
		return
	}

	endpoint := strings.TrimSpace(a.setting("wa_local_url"))
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8666/send"
	}

	apiKey := strings.TrimSpace(a.setting("wa_local_key"))
	if apiKey == "" {
		apiKey = "local-wa-key"
	}

	message := fmt.Sprintf(
		"✅ Pembayaran Berhasil\n\nVoucher WiFi %s\n\nOrder ID: %s\nPaket: %s\nKode Voucher: %s\n\nSalin kode voucher dan masukkan di halaman login WiFi.\n\nLihat kode voucher di web:\n%s",
		router,
		orderID,
		pkg,
		user,
		a.invoiceURL(orderID),
	)

	payload, err := json.Marshal(map[string]string{
		"number":  phone,
		"message": message,
	})
	if err != nil {
		msg := "json payload: " + err.Error()
		log.Printf("local wa failed tx=%d order=%s err=%v", txID, orderID, err)
		a.setWAStatus(txID, "failed", msg)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		msg := "create request: " + err.Error()
		log.Printf("local wa failed tx=%d order=%s err=%v", txID, orderID, err)
		a.setWAStatus(txID, "failed", msg)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	start := time.Now()
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	dur := time.Since(start)

	if err != nil {
		msg := "request error: " + err.Error()
		log.Printf("local wa failed tx=%d order=%s phone=%s duration=%s err=%v", txID, orderID, phone, dur, err)
		a.setWAStatus(txID, "failed", msg)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := strings.TrimSpace(string(bodyBytes))

	if len(body) > 1500 {
		body = body[:1500]
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !strings.Contains(body, `"ok":true`) {
		msg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, body)
		waStatus := "failed"
		bodyLower := strings.ToLower(body)
		if strings.Contains(bodyLower, "account_restricted") || strings.Contains(bodyLower, "restricted") {
			waStatus = "restricted"
			msg = "Akun WhatsApp gateway sedang dibatasi WhatsApp: " + msg
		} else if strings.Contains(bodyLower, "no lid for user") {
			waStatus = "invalid_number"
			msg = "Nomor tidak terdaftar di WhatsApp: " + msg
		}
		// VOUCHERGO_WA_RESTRICTED_STATUS_V1
		log.Printf("local wa failed tx=%d order=%s phone=%s status=%d duration=%s body=%s", txID, orderID, phone, resp.StatusCode, dur, body)
		a.setWAStatus(txID, waStatus, msg)
		return
	}

	log.Printf("local wa success tx=%d order=%s phone=%s status=%d duration=%s body=%s", txID, orderID, phone, resp.StatusCode, dur, body)
	a.setWAStatus(txID, "sent", "")
}
