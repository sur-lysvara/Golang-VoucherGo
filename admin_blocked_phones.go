package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type blockedPhonePageData struct {
	Msg        string
	CreateMode bool
	EditMode   bool
	Edit       PhoneBan
	List       []PhoneBan
}

func (a *App) adminBlockedPhones(w http.ResponseWriter, r *http.Request) {
	msg := strings.TrimSpace(r.URL.Query().Get("msg"))
	createMode := r.URL.Query().Get("new") == "1"
	editID, _ := strconv.ParseInt(r.URL.Query().Get("edit"), 10, 64)
	editMode := false
	var edit PhoneBan

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		action := strings.TrimSpace(r.FormValue("action"))
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		phoneDisplay := strings.TrimSpace(r.FormValue("phone_display"))
		phoneNormalized := normalizePhoneForBanCheck(phoneDisplay)
		reason := strings.TrimSpace(r.FormValue("reason"))
		bannedUntil := strings.TrimSpace(r.FormValue("banned_until"))
		bannedUntilValue := blockedPhoneBannedUntilValue(bannedUntil)
		isActive := 0
		if r.FormValue("is_active") == "1" {
			isActive = 1
		}

		switch action {
		case "create":
			if phoneNormalized == "" || reason == "" {
				msg = "Nomor WhatsApp dan alasan wajib diisi."
				break
			}

			_, err := a.db.Exec(`INSERT INTO blocked_customer_phones(phone_normalized,phone_display,reason,is_active,banned_until)
				VALUES(?,?,?,?,?)
				ON DUPLICATE KEY UPDATE phone_display=VALUES(phone_display), reason=VALUES(reason), is_active=VALUES(is_active), banned_until=VALUES(banned_until)`,
				phoneNormalized, phoneDisplay, reason, isActive, bannedUntilValue)
			if err != nil {
				log.Printf("blocked phone create failed normalized=%s err=%v", phoneNormalized, err)
				msg = "Gagal menyimpan nomor diblokir."
				break
			}

			http.Redirect(w, r, "/admin/blocked-phones?msg=Nomor+berhasil+disimpan", http.StatusSeeOther)
			return

		case "update":
			if id <= 0 || phoneNormalized == "" || reason == "" {
				msg = "Data nomor diblokir tidak valid."
				break
			}

			_, err := a.db.Exec(`UPDATE blocked_customer_phones
				SET phone_normalized=?, phone_display=?, reason=?, is_active=?, banned_until=?
				WHERE id=?`, phoneNormalized, phoneDisplay, reason, isActive, bannedUntilValue, id)
			if err != nil {
				log.Printf("blocked phone update failed id=%d normalized=%s err=%v", id, phoneNormalized, err)
				msg = "Gagal mengupdate nomor diblokir."
				break
			}

			http.Redirect(w, r, "/admin/blocked-phones?msg=Nomor+berhasil+diupdate", http.StatusSeeOther)
			return

		case "disable":
			if id <= 0 {
				msg = "ID nomor diblokir tidak valid."
				break
			}

			if _, err := a.db.Exec(`UPDATE blocked_customer_phones SET is_active=0 WHERE id=?`, id); err != nil {
				log.Printf("blocked phone disable failed id=%d err=%v", id, err)
				msg = "Gagal menonaktifkan blokir."
				break
			}

			http.Redirect(w, r, "/admin/blocked-phones?msg=Blokir+dinonaktifkan", http.StatusSeeOther)
			return

		case "enable":
			if id <= 0 {
				msg = "ID nomor diblokir tidak valid."
				break
			}

			if _, err := a.db.Exec(`UPDATE blocked_customer_phones SET is_active=1 WHERE id=?`, id); err != nil {
				log.Printf("blocked phone enable failed id=%d err=%v", id, err)
				msg = "Gagal mengaktifkan blokir."
				break
			}

			http.Redirect(w, r, "/admin/blocked-phones?msg=Blokir+diaktifkan", http.StatusSeeOther)
			return
		}
	}

	rows, err := a.db.Query(`SELECT
			id,
			phone_normalized,
			COALESCE(phone_display,''),
			reason,
			is_active,
			CASE WHEN is_active=1 AND (banned_until IS NULL OR banned_until > NOW()) THEN 1 ELSE 0 END,
			COALESCE(DATE_FORMAT(banned_until,'%Y-%m-%d %H:%i:%s'),''),
			DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(updated_at,'%Y-%m-%d %H:%i:%s')
		FROM blocked_customer_phones
		ORDER BY is_active DESC, id DESC
		LIMIT 500`)
	if err != nil {
		log.Printf("blocked phone list failed err=%v", err)
		http.Error(w, "Gagal memuat nomor diblokir", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []PhoneBan
	for rows.Next() {
		var b PhoneBan
		var active int
		var activeNow int
		if err := rows.Scan(&b.ID, &b.PhoneNormalized, &b.PhoneDisplay, &b.Reason, &active, &activeNow, &b.BannedUntil.String, &b.CreatedAt, &b.UpdatedAt); err != nil {
			log.Printf("blocked phone scan failed err=%v", err)
			continue
		}
		b.IsActive = active == 1
		b.IsCurrentlyBlocked = activeNow == 1
		b.BannedUntil.Valid = b.BannedUntil.String != ""
		if editID > 0 && b.ID == editID {
			edit = b
			editMode = true
		}
		list = append(list, b)
	}

	if editID > 0 && !editMode {
		msg = "Nomor diblokir tidak ditemukan."
	}

	body := renderHTML(`
<div class="card blocked-phone-head-v1">
	<div>
		<h2>Nomor Diblokir</h2>
		<p class="muted">Kelola nomor WhatsApp yang tidak boleh membuat order voucher.</p>
	</div>
	<div class="actions">
		{{if or .CreateMode .EditMode}}
			<a class="btn light" href="/admin/blocked-phones">Lihat Daftar</a>
		{{else}}
			<a class="btn pri" href="/admin/blocked-phones?new=1">Tambah Nomor</a>
		{{end}}
	</div>
</div>

{{if .Msg}}<div class="notice">{{.Msg}}</div>{{end}}

{{if .CreateMode}}
<div class="card">
	<h2>Tambah Nomor Diblokir</h2>
	<form method="post" action="/admin/blocked-phones" class="grid">
		<input type="hidden" name="action" value="create">
		<div class="field"><label>No. WhatsApp</label><input name="phone_display" placeholder="0812 3456 7890" required></div>
		<div class="field"><label>Status</label><select name="is_active"><option value="1">Aktif diblokir</option><option value="0">Nonaktif</option></select></div>
		<div class="field"><label>Banned sampai</label><input name="banned_until" type="datetime-local"><small class="muted">Kosongkan untuk permanen sampai admin buka.</small></div>
		<div class="field" style="grid-column:1/-1"><label>Alasan</label><textarea name="reason" rows="4" required placeholder="Contoh: Penyalahgunaan bypass pembayaran"></textarea></div>
		<div class="field"><label>&nbsp;</label><div class="actions"><button class="btn pri" type="submit">Simpan</button><a class="btn light" href="/admin/blocked-phones">Batal</a></div></div>
	</form>
</div>
{{else if .EditMode}}
<div class="card">
	<h2>Edit Nomor Diblokir</h2>
	<form method="post" action="/admin/blocked-phones" class="grid">
		<input type="hidden" name="action" value="update">
		<input type="hidden" name="id" value="{{.Edit.ID}}">
		<div class="field"><label>No. WhatsApp</label><input name="phone_display" value="{{if .Edit.PhoneDisplay}}{{.Edit.PhoneDisplay}}{{else}}{{.Edit.PhoneNormalized}}{{end}}" required></div>
		<div class="field"><label>Nomor Normalized</label><input value="{{.Edit.PhoneNormalized}}" disabled></div>
		<div class="field"><label>Status</label><select name="is_active"><option value="1" {{if .Edit.IsActive}}selected{{end}}>Aktif diblokir</option><option value="0" {{if not .Edit.IsActive}}selected{{end}}>Nonaktif</option></select></div>
		<div class="field"><label>Banned sampai</label><input name="banned_until" type="datetime-local" value="{{datetimeLocal .Edit.BannedUntil.String}}"><small class="muted">Kosongkan untuk permanen sampai admin buka.</small></div>
		<div class="field" style="grid-column:1/-1"><label>Alasan</label><textarea name="reason" rows="4" required>{{.Edit.Reason}}</textarea></div>
		<div class="field"><label>&nbsp;</label><div class="actions"><button class="btn ok" type="submit">Simpan Perubahan</button><a class="btn light" href="/admin/blocked-phones">Batal</a></div></div>
	</form>
</div>
{{else}}
<div class="card">
	<div class="tablewrap">
	<table class="blocked-phone-table-v1">
		<tr><th>Nomor</th><th>Status</th><th>Banned sampai</th><th>Alasan</th><th>Updated</th><th>Aksi</th></tr>
		{{range .List}}
		<tr>
			<td><b>{{if .PhoneDisplay}}{{.PhoneDisplay}}{{else}}{{.PhoneNormalized}}{{end}}</b><br><span class="muted">{{.PhoneNormalized}}</span></td>
			<td>{{if .IsCurrentlyBlocked}}<span class="badge failed">Aktif</span>{{else if .IsActive}}<span class="badge pending">Kedaluwarsa</span>{{else}}<span class="badge skipped">Nonaktif</span>{{end}}</td>
			<td>{{if .BannedUntil.Valid}}{{.BannedUntil.String}}{{else}}Permanen{{end}}</td>
			<td>{{.Reason}}</td>
			<td><span class="muted">{{.UpdatedAt}}</span></td>
			<td>
				<div class="actions">
					<a class="btn light" href="/admin/blocked-phones?edit={{.ID}}">Edit</a>
					{{if .IsActive}}
					<form method="post" action="/admin/blocked-phones" onsubmit="return confirm('Nonaktifkan blokir nomor ini?')">
						<input type="hidden" name="action" value="disable"><input type="hidden" name="id" value="{{.ID}}">
						<button class="btn" type="submit">Nonaktif</button>
					</form>
					{{else}}
					<form method="post" action="/admin/blocked-phones" onsubmit="return confirm('Aktifkan blokir nomor ini?')">
						<input type="hidden" name="action" value="enable"><input type="hidden" name="id" value="{{.ID}}">
						<button class="btn danger" type="submit">Aktifkan</button>
					</form>
					{{end}}
				</div>
			</td>
		</tr>
		{{else}}
		<tr><td colspan="6"><div class="notice">Belum ada nomor diblokir.</div></td></tr>
		{{end}}
	</table>
	</div>
</div>
{{end}}
`, blockedPhonePageData{
		Msg:        msg,
		CreateMode: createMode,
		EditMode:   editMode,
		Edit:       edit,
		List:       list,
	})

	page(w, "Nomor Diblokir", body)
}

func blockedPhoneBannedUntilValue(v string) any {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	v = strings.Replace(v, "T", " ", 1)
	if len(v) == len("2006-01-02 15:04") {
		v += ":00"
	}
	return v
}

func datetimeLocal(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	v = strings.Replace(v, " ", "T", 1)
	if len(v) >= len("2006-01-02T15:04") {
		return v[:len("2006-01-02T15:04")]
	}
	return v
}
