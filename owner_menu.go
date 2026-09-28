package main

import "net/http"

// OWNER_MENU_PAGE_V1
func (a *App) adminOwner(w http.ResponseWriter, r *http.Request) {
	if env("OWNER_PANEL", "0") != "1" {
		http.NotFound(w, r)
		return
	}

	body := renderHTML(`
<div class="card owner-menu-page-v1">
	<div class="owner-menu-head-v1">
		<div>
			<h2>Owner Panel</h2>
			<p class="muted">Pusat pengelolaan demo, client, dan tagihan aplikasi.</p>
		</div>
	</div>

	<div class="grid owner-menu-grid-v1">
		<a class="card owner-menu-card-v1" href="/admin/clients">
			<b>Client Manager</b>
			<span>Kelola client, buat client baru, cek status, action, onboarding, dan hapus client.</span>
		</a>

		<a class="card owner-menu-card-v1" href="/admin/demo">
			<b>Demo Manager</b>
			<span>Kelola demo, reset demo, user demo, PIN WhatsApp demo, dan status layanan demo.</span>
		</a>

		<a class="card owner-menu-card-v1" href="/admin/client-billing">
			<b>Billing Client</b>
			<span>Kelola masa aktif, tagihan, pembayaran, dan penangguhan layanan client.</span>
		</a>
	</div>
</div>
`, nil)

	page(w, "Owner", body)
}
