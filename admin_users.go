package main

import (
	"net/http"
	"strconv"
	"strings"
)

func (a *App) adminCanAccessTransaction(r *http.Request, txID int64) bool {
	role, routerID := a.currentAdminRole(r)
	if role == "superadmin" {
		return true
	}
	if role != "router" || routerID <= 0 || txID <= 0 {
		return false
	}

	var ownerRouterID int64
	err := a.db.QueryRow(`SELECT router_id FROM transactions WHERE id=?`, txID).Scan(&ownerRouterID)
	return err == nil && ownerRouterID == routerID
}

func (a *App) adminCanAccessPackage(r *http.Request, pkgID int64) bool {
	role, routerID := a.currentAdminRole(r)
	if role == "superadmin" {
		return true
	}
	if role != "router" || routerID <= 0 || pkgID <= 0 {
		return false
	}

	var ownerRouterID int64
	err := a.db.QueryRow(`SELECT router_id FROM packages WHERE id=?`, pkgID).Scan(&ownerRouterID)
	return err == nil && ownerRouterID == routerID
}

func (a *App) adminUsers(w http.ResponseWriter, r *http.Request) {
	// ADMIN_USERS_SEPARATE_PAGES_CLEAN_V1
	msg := ""

	createMode := r.URL.Query().Get("new") == "1"
	editID, _ := strconv.ParseInt(r.URL.Query().Get("edit_user"), 10, 64)
	editMode := false
	var editUser AdminUserRow

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		action := r.FormValue("action")
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirm")
		role := strings.TrimSpace(r.FormValue("role"))
		routerID, _ := strconv.ParseInt(r.FormValue("router_id"), 10, 64)

		enabled := 0
		if r.FormValue("enabled") == "1" {
			enabled = 1
		}

		if role != "router" {
			role = "superadmin"
			routerID = 0
		}

		var routerValue any = nil
		if routerID > 0 {
			routerValue = routerID
		}

		switch action {
		case "create":
			if username == "" || password == "" {
				msg = "Username dan password wajib diisi."
				break
			}

			if password != passwordConfirm {
				msg = "Konfirmasi password tidak sama."
				break
			}

			hash, err := makePasswordHash(password)
			if err != nil {
				msg = "Gagal hash password: " + err.Error()
				break
			}

			_, err = a.db.Exec(`INSERT INTO users(username,password_hash,role,router_id,enabled) VALUES(?,?,?,?,?)`, username, hash, role, routerValue, enabled)
			if err != nil {
				msg = "Gagal tambah user: " + err.Error()
				break
			}

			msg = "User berhasil ditambahkan."
			createMode = false

		case "update":
			if id <= 0 || username == "" {
				msg = "Data user tidak valid."
				break
			}

			if password != "" {
				if password != passwordConfirm {
					msg = "Konfirmasi password baru tidak sama."
					break
				}

				hash, err := makePasswordHash(password)
				if err != nil {
					msg = "Gagal hash password: " + err.Error()
					break
				}

				_, err = a.db.Exec(`UPDATE users SET username=?, password_hash=?, role=?, router_id=?, enabled=? WHERE id=?`, username, hash, role, routerValue, enabled, id)
				if err != nil {
					msg = "Gagal update user: " + err.Error()
					break
				}
			} else {
				_, err := a.db.Exec(`UPDATE users SET username=?, role=?, router_id=?, enabled=? WHERE id=?`, username, role, routerValue, enabled, id)
				if err != nil {
					msg = "Gagal update user: " + err.Error()
					break
				}
			}

			msg = "User berhasil diupdate."
			editMode = false
			editID = 0

		case "disable":
			currentID := a.currentUser(r)
			if id == currentID {
				msg = "User yang sedang login tidak bisa dinonaktifkan."
				break
			}

			_, err := a.db.Exec(`UPDATE users SET enabled=0 WHERE id=?`, id)
			if err != nil {
				msg = "Gagal disable user: " + err.Error()
				break
			}

			msg = "User berhasil dinonaktifkan."

		case "delete":
			currentID := a.currentUser(r)

			if id <= 0 {
				msg = "ID user tidak valid."
				break
			}

			if id == currentID {
				msg = "User yang sedang login tidak bisa dihapus."
				break
			}

			var targetRole string
			var targetEnabled int

			err := a.db.QueryRow(`SELECT COALESCE(role,'superadmin'), COALESCE(enabled,1) FROM users WHERE id=?`, id).Scan(&targetRole, &targetEnabled)
			if err != nil {
				msg = "User tidak ditemukan: " + err.Error()
				break
			}

			if targetRole == "superadmin" && targetEnabled == 1 {
				var activeSuper int
				_ = a.db.QueryRow(`SELECT COUNT(*) FROM users WHERE COALESCE(role,'superadmin')='superadmin' AND COALESCE(enabled,1)=1`).Scan(&activeSuper)

				if activeSuper <= 1 {
					msg = "Tidak bisa hapus superadmin aktif terakhir."
					break
				}
			}

			_, err = a.db.Exec(`DELETE FROM users WHERE id=?`, id)
			if err != nil {
				_, err2 := a.db.Exec(`UPDATE users SET enabled=0 WHERE id=?`, id)
				if err2 != nil {
					msg = "Gagal hapus user: " + err.Error()
					break
				}

				msg = "User tidak bisa dihapus permanen, jadi dinonaktifkan."
				break
			}

			msg = "User berhasil dihapus."
		}
	}

	routerRows, err := a.db.Query(`SELECT id,name,slug FROM routers WHERE enabled=1 ORDER BY name`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer routerRows.Close()

	var routers []AdminRouterOption

	for routerRows.Next() {
		var ro AdminRouterOption
		_ = routerRows.Scan(&ro.ID, &ro.Name, &ro.Slug)
		routers = append(routers, ro)
	}

	rows, err := a.db.Query(`SELECT
			u.id,
			u.username,
			COALESCE(u.role,'superadmin'),
			COALESCE(u.router_id,0),
			COALESCE(r.name,'-'),
			COALESCE(r.slug,''),
			COALESCE(u.enabled,1),
			DATE_FORMAT(u.created_at,'%Y-%m-%d %H:%i')
		FROM users u
		LEFT JOIN routers r ON r.id=u.router_id
		ORDER BY u.id`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []AdminUserRow

	for rows.Next() {
		var u AdminUserRow
		_ = rows.Scan(&u.ID, &u.Username, &u.Role, &u.RouterID, &u.RouterName, &u.RouterSlug, &u.Enabled, &u.CreatedAt)

		if editID > 0 && u.ID == editID {
			editUser = u
			editMode = true
		}

		users = append(users, u)
	}

	body := renderHTML(`
<div class="card">
	<h2>Kelola User Admin</h2>
	{{if .Msg}}<div class="notice">{{.Msg}}</div>{{end}}

	{{if .CreateMode}}
		<p class="muted">Tambah user admin baru.</p>
	{{else if .EditMode}}
		<p class="muted">Edit user #{{.Edit.ID}}.</p>
	{{else}}
		<div class="actions">
			<a class="btn pri" href="/admin/users?new=1">Tambah User</a>
		</div>
	{{end}}
</div>

{{if .CreateMode}}
<div class="card">
	<h2>Tambah User</h2>

	<form method="post" action="/admin/users" class="grid">
		<input type="hidden" name="action" value="create">

		<div class="field">
			<label>Username</label>
			<input name="username" placeholder="contoh: admin-x86" required>
		</div>

		<div class="field">
			<label>Password</label>
			<div class="pass-wrap">
				<input name="password" type="password" placeholder="password awal" required>
				<button class="pass-eye" type="button" onclick="togglePass(this)">👁</button>
			</div>
		</div>

		<div class="field">
			<label>Ulangi Password</label>
			<div class="pass-wrap">
				<input name="password_confirm" type="password" placeholder="ulangi password" required>
				<button class="pass-eye" type="button" onclick="togglePass(this)">👁</button>
			</div>
		</div>

		<div class="field">
			<label>Role</label>
			<select name="role">
				<option value="superadmin">superadmin - akses semua</option>
				<option value="router">router - khusus 1 router</option>
			</select>
		</div>

		<div class="field">
			<label>Router untuk role router</label>
			<select name="router_id">
				<option value="0">- pilih router -</option>
				{{range .Routers}}
				<option value="{{.ID}}">{{.Name}} / {{.Slug}}</option>
				{{end}}
			</select>
		</div>

		<div class="field">
			<label>Status</label>
			<select name="enabled">
				<option value="1">Aktif</option>
				<option value="0">Nonaktif</option>
			</select>
		</div>

		<div class="field">
			<label>&nbsp;</label>
			<div class="actions">
				<button class="btn pri" type="submit">Tambah User</button>
				<a class="btn light" href="/admin/users">Batal</a>
			</div>
		</div>
	</form>
</div>
{{else if .EditMode}}
<div class="card">
	<h2>Edit User</h2>

	<form method="post" action="/admin/users" class="grid">
		<input type="hidden" name="action" value="update">
		<input type="hidden" name="id" value="{{.Edit.ID}}">

		<div class="field">
			<label>Username</label>
			<input name="username" value="{{.Edit.Username}}" required>
		</div>

		<div class="field">
			<label>Role</label>
			<select name="role">
				<option value="superadmin" {{if eq .Edit.Role "superadmin"}}selected{{end}}>superadmin - akses semua</option>
				<option value="router" {{if eq .Edit.Role "router"}}selected{{end}}>router - khusus 1 router</option>
			</select>
		</div>

		<div class="field">
			<label>Router untuk role router</label>
			<select name="router_id">
				<option value="0">- semua / kosong -</option>
				{{range .Routers}}
				<option value="{{.ID}}" {{if eq .ID $.Edit.RouterID}}selected{{end}}>{{.Name}} / {{.Slug}}</option>
				{{end}}
			</select>
		</div>

		<div class="field">
			<label>Status</label>
			<select name="enabled">
				<option value="1" {{if eq .Edit.Enabled 1}}selected{{end}}>Aktif</option>
				<option value="0" {{if eq .Edit.Enabled 0}}selected{{end}}>Nonaktif</option>
			</select>
		</div>

		<div class="field">
			<label>Password Baru</label>
			<div class="pass-wrap">
				<input name="password" type="password" placeholder="kosongkan jika tidak diganti">
				<button class="pass-eye" type="button" onclick="togglePass(this)">👁</button>
			</div>
		</div>

		<div class="field">
			<label>Ulangi Password Baru</label>
			<div class="pass-wrap">
				<input name="password_confirm" type="password" placeholder="ulangi password baru">
				<button class="pass-eye" type="button" onclick="togglePass(this)">👁</button>
			</div>
		</div>

		<div class="field">
			<label>&nbsp;</label>
			<div class="actions">
				<button class="btn ok" type="submit">Simpan Perubahan</button>
				<a class="btn light" href="/admin/users">Batal</a>
			</div>
		</div>
	</form>
</div>
{{else}}
<div class="card">
	<h2>Daftar User</h2>


<div class="tx-manual-note-v1">
	<b>Catatan:</b> Tandai Lunas dipakai jika pembeli bayar di luar Pakasir. Setelah itu gunakan Push MikroTik untuk memasukkan voucher ke router.
</div>
<!-- ADMIN_TX_MANUAL_OVERRIDE_NOTE_V1 -->
<div class="tablewrap">
	<table class="users-readonly-table">
		<tr>
			<th>ID</th>
			<th>Username</th>
			<th>Role</th>
			<th>Router</th>
			<th>Status</th>
			<th>Dibuat</th>
			<th>Aksi</th>
		</tr>

		{{range .Users}}
		<tr>
			<td>{{.ID}}</td>
			<td><b>{{.Username}}</b></td>
			<td>{{.Role}}</td>
			<td>{{.RouterName}} {{if .RouterSlug}}/{{.RouterSlug}}{{end}}</td>
			<td>{{if eq .Enabled 1}}Aktif{{else}}Nonaktif{{end}}</td>
			<td>{{.CreatedAt}}</td>
			<td>
				<div class="actions">
					<a class="btn light" href="/admin/users?edit_user={{.ID}}">Edit</a>

					{{if ne .ID $.CurrentID}}
					<form method="post" action="/admin/users" onsubmit="return confirm('Hapus user ini?')">
						<input type="hidden" name="action" value="delete">
						<input type="hidden" name="id" value="{{.ID}}">
						<button class="btn danger user-delete-btn" type="submit" title="Hapus user" aria-label="Hapus user">Hapus</button>
					</form>
					{{end}}
				</div>
			</td>
		</tr>
		{{end}}
	</table>
	</div>
</div>
{{end}}

<script>
function togglePass(btn){
  const wrap = btn.closest('.pass-wrap');
  if(!wrap) return;
  const input = wrap.querySelector('input');
  if(!input) return;
  input.type = input.type === 'password' ? 'text' : 'password';
  btn.textContent = input.type === 'password' ? '👁' : '🙈';
}
</script>
`, map[string]any{
		"Msg":        msg,
		"Routers":    routers,
		"Users":      users,
		"CurrentID":  a.currentUser(r),
		"CreateMode": createMode,
		"EditMode":   editMode,
		"Edit":       editUser,
	})

	page(w, "Users", body)
}
