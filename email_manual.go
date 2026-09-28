package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (a *App) sendEmailManual(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seen := map[int64]bool{}
	count := 0

	for _, raw := range r.Form["id"] {
		id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil || id <= 0 || seen[id] {
			continue
		}
		seen[id] = true
		a.sendVoucherEmail(id)
		count++
	}

	log.Printf("manual voucher email requested count=%d", count)
	http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
}
