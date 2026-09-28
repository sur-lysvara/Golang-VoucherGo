package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CLIENT_BILLING_SOFT_SUSPEND_V1
func clientBillingFeatureEnabledV1() bool {
	if env("OWNER_PANEL", "0") == "1" {
		return false
	}
	if env("DEMO_MODE", "0") == "1" {
		return false
	}
	if strings.TrimSpace(env("CLIENT_SLUG", "")) == "" {
		return false
	}
	return true
}

func clientBillingSuspendedNowV1() (bool, string) {
	if !clientBillingFeatureEnabledV1() {
		return false, ""
	}

	status := strings.ToLower(strings.TrimSpace(env("CLIENT_BILLING_STATUS", "active")))
	if status == "suspended" {
		return true, "status_suspended"
	}

	grace := strings.TrimSpace(env("CLIENT_BILLING_GRACE_UNTIL", ""))
	if grace != "" {
		if t, err := time.Parse("2006-01-02", grace); err == nil {
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			graceDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, now.Location())
			if today.After(graceDay) {
				return true, "grace_expired"
			}
		}
	}

	return false, ""
}

func billingAllowedDuringSuspendV1(path string) bool {
	if path == "/billing/suspended" {
		return true
	}
	if path == "/admin/billing" {
		return true
	}
	if path == "/admin/billing/status.json" {
		return true
	}
	if path == "/login" || path == "/logout" {
		return true
	}
	if strings.HasPrefix(path, "/assets/") {
		return true
	}
	if strings.HasPrefix(path, "/static/") {
		return true
	}
	return false
}

func (a *App) requireBillingActive(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		suspended, _ := clientBillingSuspendedNowV1()
		if !suspended || billingAllowedDuringSuspendV1(r.URL.Path) {
			next(w, r)
			return
		}

		http.Redirect(w, r, "/billing/suspended?from="+url.QueryEscape(r.URL.Path), http.StatusSeeOther)
	}
}

func (a *App) billingSuspended(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(env("CLIENT_SLUG", "client"))
	status := strings.TrimSpace(env("CLIENT_BILLING_STATUS", "suspended"))
	dueAt := strings.TrimSpace(env("CLIENT_BILLING_DUE_AT", "-"))
	graceUntil := strings.TrimSpace(env("CLIENT_BILLING_GRACE_UNTIL", "-"))
	price := strings.TrimSpace(env("CLIENT_BILLING_PRICE", "-"))
	payURL := strings.TrimSpace(env("CLIENT_BILLING_OWNER_URL", ""))
	if payURL == "" {
		payURL = "https://vouchergo.biz.id/billing/pay?client=" + url.QueryEscape(slug)
	}

	body := renderHTML(`
<div class="billing-suspended-page-v1">
	<div class="billing-suspended-card-v1">
		<div class="billing-suspended-icon-v1">⏸</div>
		<h1>Aplikasi Ditangguhkan</h1>
		<p class="billing-suspended-lead-v1">Masa aktif VoucherGo untuk client ini sudah berakhir atau sedang ditangguhkan.</p>

		<div class="billing-suspended-grid-v1">
			<div>
				<b>Aplikasi</b>
				<span>{{.Slug}}</span>
			</div>
			<div>
				<b>Status</b>
				<span>{{.Status}}</span>
			</div>
			<div>
				<b>Jatuh Tempo</b>
				<span>{{.DueAt}}</span>
			</div>
			<div>
				<b>Masa Tenggang</b>
				<span>{{.GraceUntil}}</span>
			</div>
			<div>
				<b>Biaya Sewa</b>
				<span>{{.Price}}</span>
			</div>
		</div>

		<div class="billing-suspended-actions-v1">
			<a class="btn pri" href="{{.PayURL}}">Bayar Sekarang</a>
			<a class="btn light" href="/admin/billing">Lihat Billing</a>
			<a class="btn light" href="/logout">Logout</a>
		</div>

		<div class="billing-suspended-note-v1">
			<b>Catatan:</b> Setelah pembayaran sukses, sistem owner akan mengaktifkan kembali aplikasi secara otomatis.
		</div>
	</div>
</div>
`, map[string]any{
		"Slug":       slug,
		"Status":     status,
		"DueAt":      dueAt,
		"GraceUntil": graceUntil,
		"Price":      price,
		"PayURL":     payURL,
	})

	publicPage(w, "Aplikasi Ditangguhkan", body)
}
