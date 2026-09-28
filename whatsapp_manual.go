package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func normalizePhoneID(s string) string {
	// VOUCHERGO_WA_PHONE_NORMALIZE_ID_V2
	// Terima format umum Indonesia:
	// 08xxxxxxxxxx  -> 628xxxxxxxxxx
	// +628xxxxxxxxx -> 628xxxxxxxxx
	// 628xxxxxxxxxx -> 628xxxxxxxxxx
	// 8xxxxxxxxxx   -> 628xxxxxxxxxx
	digits := ""
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digits += string(r)
		}
	}

	if digits == "" {
		return ""
	}

	// 6208xxx biasanya akibat normalisasi dobel/salah input.
	if strings.HasPrefix(digits, "6208") {
		digits = "62" + strings.TrimPrefix(digits[2:], "0")
	} else if strings.HasPrefix(digits, "08") {
		digits = "62" + strings.TrimPrefix(digits, "0")
	} else if strings.HasPrefix(digits, "8") {
		digits = "62" + digits
	}

	// Nomor Indonesia WA umumnya 62 + 9-13 digit setelahnya.
	// Reject nomor kependekan seperti "5326" atau "628".
	if !strings.HasPrefix(digits, "62") {
		return ""
	}
	if len(digits) < 10 || len(digits) > 15 {
		return ""
	}

	return digits
}

func isSafePhoneText(s string) bool {
	if len(s) < 8 || len(s) > 24 {
		return false
	}

	digits := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digits++
			continue
		}
		if r == '+' || r == '-' || r == ' ' || r == '(' || r == ')' {
			continue
		}
		return false
	}

	return digits >= 10 && digits <= 16
}

func (a *App) sendWAManual(w http.ResponseWriter, r *http.Request) {
	log.Printf("local wa manual handler called method=%s path=%s", r.Method, r.URL.Path)

	f, _ := os.OpenFile("/tmp/tuku-wa-click.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		fmt.Fprintf(f, "%s handler called method=%s path=%s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		defer f.Close()
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = r.ParseForm()

	id, _ := strconv.ParseInt(strings.TrimSpace(r.FormValue("id")), 10, 64)

	// ADMIN_SEND_WA_SCOPE_V1
	if !a.adminCanAccessTransaction(r, id) {
		http.Error(w, "Forbidden: transaksi ini bukan milik router user ini", http.StatusForbidden)
		return
	}

	log.Printf("local wa manual request tx=%d raw_id=%q", id, r.FormValue("id"))

	if f != nil {
		fmt.Fprintf(f, "%s raw_id=%q parsed_id=%d\n", time.Now().Format("2006-01-02 15:04:05"), r.FormValue("id"), id)
	}

	if id > 0 {
		// VOUCHERGO_MANUAL_WA_PENDING_SEND_INVOICE_V1
		// Pending = kirim ulang invoice/QRIS belum dibayar.
		// Paid/lunas = kirim voucher pembayaran berhasil.
		var txStatus string
		_ = a.db.QueryRow(`SELECT COALESCE(status,'') FROM transactions WHERE id=?`, id).Scan(&txStatus)
		txStatus = strings.ToLower(strings.TrimSpace(txStatus))

		switch txStatus {
		case "paid", "lunas", "success", "sukses", "settlement", "capture":
			a.sendWanesia(id)
		default:
			a.sendUnpaidPaymentWA(id)
		}
	} else {
		log.Printf("local wa manual skip reason=invalid tx id")
	}

	http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
}
