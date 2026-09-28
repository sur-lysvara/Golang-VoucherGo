package main

import (
	"net/http"
)

func (a *App) settings(w http.ResponseWriter, r *http.Request) {
	canEditBranding := env("OWNER_PANEL", "0") == "1" || env("VIP_BRANDING", "0") == "1"
	brandingAllowed := vipBrandingEnabled()
	// ADMIN_SETTINGS_RESPONSIVE_V1
	keys := []string{
		"public_base_url",
		"public_subtitle",
		"maintenance_enabled",
		"maintenance_title",
		"maintenance_message",
		"wa_provider",
		"wa_local_url",
		"wa_local_key",
		"email_enabled",
		"smtp_host",
		"smtp_port",
		"smtp_username",
		"smtp_password",
		"smtp_from_email",
		"smtp_from_name",
		"email_test_to",
	}

	if brandingAllowed {
		if canEditBranding {
			keys = append([]string{"site_name"}, keys...)
		}
	}

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		for _, k := range keys {
			_ = a.setSetting(k, r.FormValue(k))
		}

		http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
		return
	}

	vals := map[string]string{}
	if r.URL.Query().Get("email_test") == "ok" {
		vals["email_test_notice"] = "Test email berhasil terkirim."
	}
	if r.URL.Query().Get("email_test_error") != "" {
		vals["email_test_error"] = r.URL.Query().Get("email_test_error")
	}
	vals["can_edit_branding"] = "0"
	if canEditBranding {
		vals["can_edit_branding"] = "1"
	}
	for _, k := range keys {
		vals[k] = a.setting(k)
	}
	if vals["smtp_port"] == "" {
		vals["smtp_port"] = "587"
	}
	if vals["smtp_from_name"] == "" {
		vals["smtp_from_name"] = "VoucherGo"
	}

	if brandingAllowed {
		vals["branding_allowed"] = "1"
	} else {
		vals["branding_allowed"] = ""
		vals["branding_locked_notice"] = "Branding website hanya tersedia untuk client VIP."
	}

	body := renderHTML(`
<div class="card settings-page-card">
	<h2>Settings</h2>
	<p class="muted">Kelola URL publik, koneksi WhatsApp lokal, dan pengaturan operasional aplikasi.</p>


        <!-- SETTINGS_RELATED_MENU_V1 -->
        <div class="card settings-page-card" style="margin-bottom:16px">
                <h2>Pengaturan Lanjutan</h2>
                <p class="muted">Akses cepat ke pengaturan WhatsApp dan password admin.</p>
                <div class="actions settings-actions" style="display:flex;gap:10px;flex-wrap:wrap">
                        <a class="btn pri" href="/admin/whatsapp">WhatsApp</a>
                        <a class="btn" href="/admin/password">Password</a>
                </div>
        </div>

        <form method="post" class="settings-form-compact-v1">
		<div class="grid">
			{{if index . "branding_allowed"}}
			{{if eq (index . "can_edit_branding") "1"}}
                        <div class="field"><label>Nama Website</label><input name="site_name" value="{{index . "site_name"}}"></div>
                        {{end}}
			{{else}}
			<div class="notice" style="grid-column:1/-1">{{index . "branding_locked_notice"}}</div>
			{{end}}
			<div class="field"><label>Public Base URL</label><input name="public_base_url" value="{{index . "public_base_url"}}"></div>
			<div class="field" style="grid-column:1/-1">
				<label>Subtitle Halaman Publik</label>
				<textarea name="public_subtitle" rows="4" placeholder="Beli voucher hotspot, bayar via QRIS.">{{index . "public_subtitle"}}</textarea>
				<small class="muted">Bisa pakai baris baru agar tampil lebih rapi di halaman publik.</small>
			</div>
			<div style="grid-column:1/-1;margin-top:8px">
				<h3 style="margin:8px 0 4px">Maintenance Publik</h3>
				<p class="muted" style="margin:0 0 8px">Jika aktif, halaman publik pembelian voucher dialihkan ke halaman sedang update. Admin tetap bisa login dan mengubah setting ini.</p>
			</div>
			<div class="field">
				<label>Status Maintenance Publik</label>
				<select name="maintenance_enabled">
					<option value="0">Nonaktif</option>
					<option value="1" {{if eq (index . "maintenance_enabled") "1"}}selected{{end}}>Aktif</option>
				</select>
			</div>
			<div class="field"><label>Judul Maintenance</label><input name="maintenance_title" value="{{index . "maintenance_title"}}" placeholder="Sedang Update"></div>
			<div class="field" style="grid-column:1/-1">
				<label>Pesan Maintenance</label>
				<textarea name="maintenance_message" rows="3" placeholder="Mohon maaf, layanan sedang diperbarui. Silakan coba lagi beberapa saat lagi.">{{index . "maintenance_message"}}</textarea>
			</div>
			<div class="field"><label>WA Provider</label><input name="wa_provider" value="{{index . "wa_provider"}}" placeholder="local"></div>
			<div class="field"><label>WA Local URL</label><input name="wa_local_url" value="{{index . "wa_local_url"}}" placeholder="http://127.0.0.1:8666/send"></div>
			<div class="field secret-field"><label>WA Local Key</label><input name="wa_local_key" type="password" value="{{index . "wa_local_key"}}"><button type="button" class="pass-eye" onclick="toggleSecret(this)">👁</button></div>

			<div style="grid-column:1/-1;margin-top:8px">
				<h3 style="margin:8px 0 4px">Email Otomatis</h3>
				<p class="muted" style="margin:0 0 8px">SMTP digunakan untuk kirim voucher otomatis ke email pelanggan setelah pembayaran berhasil. Rekomendasi gunakan SMTP port 587.</p>
				{{if index . "email_test_notice"}}<div class="notice" style="margin-bottom:10px">{{index . "email_test_notice"}}</div>{{end}}
				{{if index . "email_test_error"}}<div class="notice" style="margin-bottom:10px;background:#fee2e2;color:#991b1b;border-color:#fecaca">Test email gagal: {{index . "email_test_error"}}</div>{{end}}
			</div>

			<div class="field">
				<label>Email Otomatis</label>
				<select name="email_enabled">
					<option value="0">Nonaktif</option>
					<option value="1" {{if eq (index . "email_enabled") "1"}}selected{{end}}>Aktif</option>
				</select>
			</div>
			<div class="field"><label>SMTP Host</label><input name="smtp_host" value="{{index . "smtp_host"}}" placeholder="smtp.domain.com"></div>
			<div class="field"><label>SMTP Port</label><input name="smtp_port" value="{{index . "smtp_port"}}" placeholder="587" inputmode="numeric"></div>
			<div class="field"><label>SMTP Username</label><input name="smtp_username" value="{{index . "smtp_username"}}" autocomplete="username"></div>
			<div class="field secret-field"><label>SMTP Password</label><input name="smtp_password" type="password" value="{{index . "smtp_password"}}" autocomplete="current-password"><button type="button" class="pass-eye" onclick="toggleSecret(this)">👁</button></div>
			<div class="field"><label>From Email</label><input name="smtp_from_email" type="email" value="{{index . "smtp_from_email"}}" placeholder="noreply@domain.com"></div>
			<div class="field"><label>From Name</label><input name="smtp_from_name" value="{{index . "smtp_from_name"}}" placeholder="VoucherGo"></div>
			<div class="field"><label>Email Tujuan Test</label><input name="email_test_to" type="email" value="{{index . "email_test_to"}}" placeholder="admin@email.com"></div>
		</div>

		<div class="actions settings-actions">
			<button class="btn pri" type="submit">Simpan Settings</button>
			<button class="btn" type="submit" formaction="/admin/settings/test-email" formmethod="post">Kirim Test Email</button>
		</div>
	</form>

</div>`, vals)

	page(w, "Settings", body)
}
