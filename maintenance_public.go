package main

import (
	"html/template"
	"net/http"
	"strings"
)

func (a *App) publicMaintenance(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.publicMaintenanceEnabled() {
			a.publicMaintenancePage(w, r)
			return
		}

		next(w, r)
	}
}

func (a *App) publicMaintenanceEnabled() bool {
	return strings.TrimSpace(a.setting("maintenance_enabled")) == "1"
}

func (a *App) publicMaintenancePage(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(a.setting("maintenance_title"))
	if title == "" {
		title = "Sedang Update"
	}

	message := strings.TrimSpace(a.setting("maintenance_message"))
	if message == "" {
		message = "Mohon maaf, layanan sedang diperbarui. Silakan coba lagi beberapa saat lagi."
	}

	w.Header().Set("Retry-After", "300")
	w.WriteHeader(http.StatusServiceUnavailable)

	body := renderHTML(`
<div class="hero">
	<h1 style="margin:0">{{.Title}}</h1>
	<p>{{.Message}}</p>
</div>

<div class="card" style="max-width:620px;margin:14px auto 0;text-align:center">
	<b>{{.SiteName}}</b>
	<p class="muted" style="margin-bottom:0">Pembelian voucher sementara tidak tersedia.</p>
</div>`, map[string]any{
		"Title":    title,
		"Message":  template.HTML(strings.ReplaceAll(template.HTMLEscapeString(message), "\n", "<br>")),
		"SiteName": a.setting("site_name"),
	})

	publicPage(w, title, body)
}
