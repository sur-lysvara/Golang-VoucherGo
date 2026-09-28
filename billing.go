package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// CLIENT_BILLING_PHASE1_V1
func (a *App) adminBilling(w http.ResponseWriter, r *http.Request) {
	// CLIENT_BILLING_OWNER_REDIRECT_V1
	if env("OWNER_PANEL", "0") == "1" {
		http.Redirect(w, r, "/admin/client-billing", http.StatusSeeOther)
		return
	}

	slug := strings.TrimSpace(env("CLIENT_SLUG", ""))
	if slug == "" {
		slug = strings.TrimSpace(env("APP_SLUG", ""))
	}
	if slug == "" {
		slug = strings.TrimSpace(env("CLIENT_NAME", ""))
	}
	if slug == "" {
		slug = "client"
	}

	status := strings.ToLower(strings.TrimSpace(env("CLIENT_BILLING_STATUS", "active")))
	if status == "" {
		status = "active"
	}

	dueAt := strings.TrimSpace(env("CLIENT_BILLING_DUE_AT", "-"))
	graceUntil := strings.TrimSpace(env("CLIENT_BILLING_GRACE_UNTIL", "-"))
	price := strings.TrimSpace(env("CLIENT_BILLING_PRICE", "-"))

	ownerPayURL := strings.TrimSpace(env("CLIENT_BILLING_OWNER_URL", ""))
	if ownerPayURL == "" {
		base := strings.TrimRight(env("CLIENT_BILLING_OWNER_BASE_URL", "https://vouchergo.biz.id/billing/pay"), "/")
		ownerPayURL = base + "?client=" + url.QueryEscape(slug)
	}

	statusLabel := "Aktif"
	statusClass := "ok"
	statusNote := "Layanan aplikasi aktif. Semua fitur dapat digunakan normal."

	switch status {
	case "warning":
		statusLabel = "Mendekati Jatuh Tempo"
		statusClass = "warn"
		statusNote = "Masa aktif aplikasi akan segera habis. Silakan lakukan pembayaran sebelum tanggal jatuh tempo."
	case "grace":
		statusLabel = "Masa Tenggang"
		statusClass = "warn"
		statusNote = "Aplikasi masih bisa digunakan sementara. Segera lakukan pembayaran agar layanan tidak ditangguhkan."
	case "expired":
		statusLabel = "Expired"
		statusClass = "bad"
		statusNote = "Masa aktif aplikasi sudah habis. Lakukan pembayaran agar layanan tetap aktif."
	case "suspended":
		statusLabel = "Suspended"
		statusClass = "bad"
		statusNote = "Layanan sedang ditangguhkan. Silakan lakukan pembayaran untuk mengaktifkan kembali aplikasi."
	case "active", "":
		statusLabel = "Aktif"
	default:
		statusLabel = strings.Title(status)
		statusClass = "warn"
		statusNote = "Status langganan membutuhkan pengecekan admin."
	}

	daysLeft := "-"
	if dueAt != "" && dueAt != "-" {
		if t, err := time.Parse("2006-01-02", dueAt); err == nil {
			now := time.Now()
			diff := int(t.Sub(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())).Hours() / 24)
			if diff >= 0 {
				daysLeft = strconv.Itoa(diff) + " hari"
			} else {
				daysLeft = "Lewat " + strconv.Itoa(-diff) + " hari"
			}
		}
	}

	body := renderHTML(`
<div class="client-billing-page-v1">
	<div class="client-billing-hero-v1">
		<div>
			<div class="billing-kicker-v1">Billing Aplikasi</div>
			<h2>Sewa Aplikasi VoucherGo</h2>
			<p>Lihat masa aktif aplikasi dan lakukan pembayaran atau perpanjangan.</p>
		</div>
		<div class="billing-status-pill-v1 {{.StatusClass}}">{{.StatusLabel}}</div>
	</div>

	<div class="grid billing-grid-v1">
		<div class="card billing-stat-v1">
			<b>Aplikasi</b>
			<div class="price">{{.Slug}}</div>
			<p class="muted">Identitas aplikasi</p>
		</div>
		<div class="card billing-stat-v1">
			<b>Jatuh Tempo</b>
			<div class="price">{{.DueAt}}</div>
			<p class="muted">Sisa waktu: {{.DaysLeft}}</p>
		</div>
		<div class="card billing-stat-v1">
			<b>Biaya Sewa</b>
			<div class="price">{{.Price}}</div>
			<p class="muted">Biaya sewa per periode</p>
		</div>
		<div class="card billing-stat-v1">
			<b>Masa Tenggang</b>
			<div class="price">{{.GraceUntil}}</div>
			<p class="muted">Batas akhir sebelum layanan ditangguhkan</p>
		</div>
	</div>

	<div class="card billing-action-card-v1">
		<h3>Status Billing</h3>
		<p>{{.StatusNote}}</p>

		<div class="billing-actions-v1">
			<a class="btn pri" href="{{.OwnerPayURL}}" target="_blank" rel="noopener">Bayar / Perpanjang</a>
			<a class="btn light" href="/admin">Kembali Dashboard</a>
		</div>

		<div class="billing-note-v1">
			<b>Catatan:</b> Pembayaran diproses melalui halaman pembayaran resmi. Setelah berhasil, masa aktif akan diperbarui otomatis.
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
		"DaysLeft":    daysLeft,
		"OwnerPayURL": ownerPayURL,
	})

	page(w, "Billing Aplikasi", body)
}

// CLIENT_BILLING_WARNING_BANNER_V1
func (a *App) adminBillingStatusJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !clientBillingFeatureEnabledV1() {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"enabled": false,
		})
		return
	}

	slug := strings.TrimSpace(env("CLIENT_SLUG", "client"))
	status := strings.ToLower(strings.TrimSpace(env("CLIENT_BILLING_STATUS", "active")))
	dueAt := strings.TrimSpace(env("CLIENT_BILLING_DUE_AT", ""))
	graceUntil := strings.TrimSpace(env("CLIENT_BILLING_GRACE_UNTIL", ""))
	price := strings.TrimSpace(env("CLIENT_BILLING_PRICE", "-"))
	payURL := strings.TrimSpace(env("CLIENT_BILLING_OWNER_URL", ""))
	if payURL == "" {
		payURL = "https://vouchergo.biz.id/billing/pay?client=" + url.QueryEscape(slug)
	}

	daysLeft := 999999
	if dueAt != "" {
		if t, err := time.Parse("2006-01-02", dueAt); err == nil {
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			dueDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
			daysLeft = int(dueDay.Sub(today).Hours() / 24)
		}
	}

	suspended, reason := clientBillingSuspendedNowV1()

	showBanner := false
	level := "ok"
	title := "Billing aplikasi aktif"
	message := "Masa aktif aplikasi masih aman."

	if suspended {
		showBanner = true
		level = "bad"
		title = "Aplikasi ditangguhkan"
		message = "Masa aktif aplikasi sudah melewati masa tenggang. Silakan lakukan pembayaran untuk mengaktifkan kembali aplikasi."
		if reason == "status_suspended" {
			message = "Layanan sedang ditangguhkan. Silakan lakukan pembayaran atau hubungi admin untuk mengaktifkan kembali aplikasi."
		}
	} else if status == "expired" || status == "grace" {
		showBanner = true
		level = "warn"
		title = "Masa aktif masuk masa tenggang"
		message = "Segera lakukan pembayaran agar layanan tidak ditangguhkan."
	} else if daysLeft <= 7 {
		showBanner = true
		level = "warn"
		title = "Masa aktif hampir habis"
		if daysLeft < 0 {
			message = "Masa aktif aplikasi sudah lewat jatuh tempo. Silakan lakukan pembayaran."
		} else if daysLeft == 0 {
			message = "Masa aktif aplikasi habis hari ini. Silakan lakukan pembayaran."
		} else {
			message = "Masa aktif aplikasi akan habis dalam " + strconv.Itoa(daysLeft) + " hari."
		}
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"enabled":     true,
		"show_banner": showBanner,
		"level":       level,
		"title":       title,
		"message":     message,
		"slug":        slug,
		"status":      status,
		"due_at":      dueAt,
		"grace_until": graceUntil,
		"price":       price,
		"days_left":   daysLeft,
		"pay_url":     payURL,
	})
}
