package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	paymentBypassSeconds           = 60
	paymentBypassStartDelaySeconds = 15
	paymentBypassLimit             = 2
)

// PAYMENT_BYPASS_LIMIT_V1
func (a *App) reservePaymentBypassAttempt(orderID string) (used int, remainingSeconds int, alreadyActive bool, allowed bool, err error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return 0, 0, false, false, fmt.Errorf("order_id kosong")
	}

	for i := 0; i < 2; i++ {
		var count int
		var remaining int

		err := a.db.QueryRow(`SELECT
				COALESCE(payment_bypass_count,0),
				GREATEST(COALESCE(TIMESTAMPDIFF(SECOND, NOW(), payment_bypass_until),0),0)
			FROM transactions
			WHERE order_id=?`, orderID).Scan(&count, &remaining)
		if err != nil {
			return 0, 0, false, false, err
		}

		if remaining > 0 {
			return count, remaining, true, true, nil
		}

		if count >= paymentBypassLimit {
			return count, 0, false, false, nil
		}

		newCount := count + 1
		res, err := a.db.Exec(`UPDATE transactions
			SET payment_bypass_count=?,
			    payment_bypass_until=DATE_ADD(NOW(), INTERVAL 60 SECOND),
			    payment_bypass_last_at=NOW()
			WHERE order_id=?
			  AND COALESCE(payment_bypass_count,0)=?
			  AND (payment_bypass_until IS NULL OR payment_bypass_until <= NOW())`,
			newCount,
			orderID,
			count,
		)
		if err != nil {
			return count, 0, false, false, err
		}

		affected, _ := res.RowsAffected()
		if affected > 0 {
			return newCount, paymentBypassSeconds, false, true, nil
		}
	}

	return 0, 0, false, false, fmt.Errorf("status bypass berubah, coba refresh")
}

func (a *App) rollbackPaymentBypassAttempt(orderID string, used int) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" || used <= 0 {
		return
	}

	_, _ = a.db.Exec(`UPDATE transactions
		SET payment_bypass_count=GREATEST(payment_bypass_count-1,0),
		    payment_bypass_until=NULL
		WHERE order_id=? AND payment_bypass_count=?`, orderID, used)
}

func cleanPaymentClientMAC(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	v = strings.ReplaceAll(v, "-", ":")
	hw, err := net.ParseMAC(v)
	if err != nil {
		return ""
	}
	return strings.ToUpper(hw.String())
}

func cleanPaymentClientIP(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	ip := net.ParseIP(v)
	if ip == nil {
		return ""
	}
	return ip.String()
}

func publicBuyExtraQuery(r *http.Request, routerSlug string) template.URL {
	q := url.Values{}

	routerSlug = strings.TrimSpace(routerSlug)
	if routerSlug != "" {
		q.Set("router", routerSlug)
	}

	if mac := cleanPaymentClientMAC(r.URL.Query().Get("mac")); mac != "" {
		q.Set("mac", mac)
	}
	if ip := cleanPaymentClientIP(r.URL.Query().Get("ip")); ip != "" {
		q.Set("ip", ip)
	}

	if len(q) == 0 {
		return template.URL("")
	}
	return template.URL("&" + q.Encode())
}

func (a *App) startPaymentBypass90(ro Router, orderID, mac, ip string) error {
	mac = cleanPaymentClientMAC(mac)
	ip = cleanPaymentClientIP(ip)

	if mac == "" {
		return nil
	}
	if ro.Host == "" || ro.Username == "" {
		return fmt.Errorf("router host/username kosong")
	}

	name := safePaymentBypassName(orderID)
	comment := name

	mode := selectedMTAPIMode(ro.APIMode)
	switch mode {
	case "rest":
		return paymentBypassREST(ro, name, comment, mac)
	case "classic":
		return paymentBypassClassic(ro, name, comment, mac, false, true)
	case "classic_ssl":
		return paymentBypassClassic(ro, name, comment, mac, true, true)
	default:
		if err := paymentBypassREST(ro, name, comment, mac); err == nil {
			return nil
		}
		return paymentBypassClassic(ro, name, comment, mac, false, false)
	}
}

func safePaymentBypassName(orderID string) string {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		orderID = randCode(8)
	}

	var b strings.Builder
	for _, r := range orderID {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}

	v := b.String()
	if v == "" {
		v = randCode(8)
	}
	if len(v) > 32 {
		v = v[:32]
	}
	return "VG-PAY-TEMP-" + v
}

func paymentBypassEvent(name, comment string) string {
	return fmt.Sprintf(`/ip hotspot ip-binding remove [find comment="%s"]; /system scheduler remove [find name="%s"]`, comment, name)
}

func paymentBypassREST(ro Router, name, comment, mac string) error {
	scheme := "http"
	if ro.UseSSL {
		scheme = "https"
	}

	if ro.Port <= 0 {
		ro.Port = 443
		if !ro.UseSSL {
			ro.Port = 80
		}
	}

	base := fmt.Sprintf("%s://%s:%d/rest", scheme, ro.Host, ro.Port)

	if err := mtRESTPut(base+"/system/scheduler", ro, map[string]string{
		"name":     name,
		"interval": strconv.Itoa(paymentBypassSeconds) + "s",
		"on-event": paymentBypassEvent(name, comment),
		"comment":  "VoucherGo payment temp bypass",
	}); err != nil {
		return err
	}

	return mtRESTPut(base+"/ip/hotspot/ip-binding", ro, map[string]string{
		"mac-address": mac,
		"type":        "bypassed",
		"comment":     comment,
	})
}

func mtRESTPut(endpoint string, ro Router, body map[string]string) error {
	b, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

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

func paymentBypassClassic(ro Router, name, comment, mac string, useTLS bool, useRouterPort bool) error {
	addr := fmt.Sprintf("%s:%d", cleanRouterHost(ro.Host), mtClassicRouterPort(ro, useTLS, useRouterPort))

	c, err := mtClassicConnect(addr, ro.Username, ro.Password, useTLS, ro.InsecureSSL)
	if err != nil {
		return err
	}
	defer c.Close()

	if err := c.Run(
		"/system/scheduler/add",
		"=name="+name,
		"=interval="+strconv.Itoa(paymentBypassSeconds)+"s",
		"=on-event="+paymentBypassEvent(name, comment),
		"=comment=VoucherGo payment temp bypass",
	); err != nil {
		return err
	}

	return c.Run(
		"/ip/hotspot/ip-binding/add",
		"=mac-address="+mac,
		"=type=bypassed",
		"=comment="+comment,
	)
}
