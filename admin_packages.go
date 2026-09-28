package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"log"
)

func (a *App) packages(w http.ResponseWriter, r *http.Request) {
	// ADMIN_PACKAGE_SEPARATE_PAGES_V1
	role, scopedRouterID := a.currentAdminRole(r)

	var routersRows *sql.Rows
	var err error

	if role == "router" && scopedRouterID > 0 {
		routersRows, err = a.db.Query(`SELECT id,name FROM routers WHERE id=? ORDER BY name`, scopedRouterID)
	} else {
		routersRows, err = a.db.Query(`SELECT id,name FROM routers ORDER BY name`)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer routersRows.Close()

	type Opt struct {
		ID   int64
		Name string
	}

	var opts []Opt

	for routersRows.Next() {
		var x Opt
		if err := routersRows.Scan(&x.ID, &x.Name); err != nil {
			log.Printf("routers option scan error: %v", err)
			continue
		}
		opts = append(opts, x)
	}

	selectedRouterID, _ := strconv.ParseInt(r.URL.Query().Get("router_id"), 10, 64)
	routerLocked := role == "router" && scopedRouterID > 0
	if routerLocked {
		selectedRouterID = scopedRouterID
	}

	if selectedRouterID > 0 && !routerLocked {
		found := false
		for _, opt := range opts {
			if opt.ID == selectedRouterID {
				found = true
				break
			}
		}
		if !found {
			selectedRouterID = 0
		}
	}

	hasRouters := len(opts) > 0
	var list []Package

	if selectedRouterID > 0 {
		rows, err := a.db.Query(`SELECT p.id,p.router_id,r.name,p.name,p.price,p.profile,p.limit_uptime,p.shared_users,p.active
		FROM packages p
		JOIN routers r ON r.id=p.router_id
		WHERE p.router_id=?
		ORDER BY p.id DESC`, selectedRouterID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var x Package
			if err := rows.Scan(&x.ID, &x.RouterID, &x.RouterName, &x.Name, &x.Price, &x.Profile, &x.LimitUptime, &x.SharedUsers, &x.Active); err != nil {
				log.Printf("package scan error: %v", err)
				continue
			}
			list = append(list, x)
		}
	}

	createMode := r.URL.Query().Get("new") == "1"
	editMode := false

	edit := Package{
		Active:      true,
		SharedUsers: 1,
	}

	if selectedRouterID > 0 {
		edit.RouterID = selectedRouterID
	} else if len(opts) > 0 {
		edit.RouterID = opts[0].ID
	}

	if id := r.URL.Query().Get("edit"); id != "" {
		editID, _ := strconv.ParseInt(id, 10, 64)

		if role == "router" && scopedRouterID > 0 {
			err := a.db.QueryRow(`SELECT id,router_id,name,price,profile,limit_uptime,shared_users,active FROM packages WHERE id=? AND router_id=?`, editID, scopedRouterID).
				Scan(&edit.ID, &edit.RouterID, &edit.Name, &edit.Price, &edit.Profile, &edit.LimitUptime, &edit.SharedUsers, &edit.Active)
			if err != nil {
				http.Error(w, "Forbidden: paket ini bukan milik router user ini", http.StatusForbidden)
				return
			}
		} else {
			err := a.db.QueryRow(`SELECT id,router_id,name,price,profile,limit_uptime,shared_users,active FROM packages WHERE id=?`, editID).
				Scan(&edit.ID, &edit.RouterID, &edit.Name, &edit.Price, &edit.Profile, &edit.LimitUptime, &edit.SharedUsers, &edit.Active)
			if err != nil {
				http.Error(w, "Paket tidak ditemukan", http.StatusNotFound)
				return
			}
		}

		editMode = true
	}

	body := renderHTML(`
<div class="card">
	<h2>Paket</h2>

	{{if .CreateMode}}
		<p class="muted">Tambah paket baru.</p>
	{{else if .EditMode}}
		<p class="muted">Edit paket #{{.Edit.ID}}.</p>
	{{else}}
		<div class="actions">
			<a class="btn pri" href="/admin/packages?new=1">Tambah Paket</a>
		</div>
	{{end}}
</div>

{{if .CreateMode}}
<div class="card">
	<h2>Tambah Paket</h2>

	<form method="post" action="/admin/packages/save" class="package-form-compact-v1">
		<input type="hidden" name="id" value="0">

		<div class="grid">
			<div class="field">
				<label>Router</label>
				<select name="router_id">
					{{range .Routers}}
					<option value="{{.ID}}" {{if eq .ID $.Edit.RouterID}}selected{{end}}>{{.Name}}</option>
					{{end}}
				</select>
			</div>

			<div class="field"><label>Nama Paket</label><input name="name" placeholder="1 Hari" required></div>
			<div class="field"><label>Harga</label><input name="price" required></div>
			<div class="field"><label>Profile MikroTik</label><input name="profile" placeholder="PAKET-1H"></div>
			<div class="field"><label>Limit Uptime</label><input name="limit_uptime" placeholder="1d / 12h"></div>
			<div class="field"><label>Shared Users</label><input name="shared_users" value="{{.Edit.SharedUsers}}"></div>
		</div>

		<div class="router-checks">
			<label>Aktif <input type="checkbox" name="active" {{checked .Edit.Active}}></label>
		</div>

		<br>
		<div class="actions">
			<button class="btn pri" type="submit">Simpan Paket</button>
			<a class="btn light" href="/admin/packages">Batal</a>
		</div>
	</form>
</div>
{{else if .EditMode}}
<div class="card">
	<h2>Edit Paket</h2>

	<form method="post" action="/admin/packages/save" class="package-form-compact-v1">
		<input type="hidden" name="id" value="{{.Edit.ID}}">

		<div class="grid">
			<div class="field">
				<label>Router</label>
				<select name="router_id">
					{{range .Routers}}
					<option value="{{.ID}}" {{if eq .ID $.Edit.RouterID}}selected{{end}}>{{.Name}}</option>
					{{end}}
				</select>
			</div>

			<div class="field"><label>Nama Paket</label><input name="name" value="{{.Edit.Name}}" required></div>
			<div class="field"><label>Harga</label><input name="price" value="{{.Edit.Price}}" required></div>
			<div class="field"><label>Profile MikroTik</label><input name="profile" value="{{.Edit.Profile}}"></div>
			<div class="field"><label>Limit Uptime</label><input name="limit_uptime" value="{{.Edit.LimitUptime}}"></div>
			<div class="field"><label>Shared Users</label><input name="shared_users" value="{{.Edit.SharedUsers}}"></div>
		</div>

		<div class="router-checks">
			<label>Aktif <input type="checkbox" name="active" {{checked .Edit.Active}}></label>
		</div>

		<br>
		<div class="actions">
			<button class="btn ok" type="submit">Simpan Perubahan</button>
			<a class="btn light" href="/admin/packages">Batal</a>
		</div>
	</form>
</div>
{{else}}
<div class="card">
	<h2>Lihat List Paket</h2>
	<p class="muted">Pilih router terlebih dahulu untuk melihat paket voucher.</p>

	<form method="get" action="/admin/packages" class="package-router-filter-v3">
		<label>Router</label>
		<div class="actions">
			<select name="router_id" {{if .RouterLocked}}disabled{{end}}>
				<option value="">Pilih router</option>
				{{range .Routers}}
				<option value="{{.ID}}" {{if eq .ID $.SelectedRouterID}}selected{{end}}>{{.Name}}</option>
				{{end}}
			</select>
			{{if .RouterLocked}}
			<input type="hidden" name="router_id" value="{{.SelectedRouterID}}">
			{{end}}
			<button class="btn pri" type="submit">Lihat List Paket</button>
		</div>
	</form>
	<!-- PACKAGE_ROUTER_SELECT_FORM_V3 -->

	{{if .SelectedRouterID}}
	<div class="tablewrap">
	<table class="packages-readonly-table">
		<tr>
			<th>Router</th>
			<th>Paket</th>
			<th>Harga</th>
			<th>Profile</th>
			<th>Status</th>
			<th>Aksi</th>
		</tr>

		{{range .List}}
		<tr>
			<td>{{.RouterName}}</td>
			<td><b>{{.Name}}</b></td>
			<td>{{rp .Price}}</td>
			<td>{{.Profile}}</td>
			<td>{{if .Active}}<span class="badge success">aktif</span>{{else}}<span class="badge failed">off</span>{{end}}</td>
			<td>
				<div class="actions">
					<a class="btn light" href="/admin/packages?edit={{.ID}}">Edit</a>
					<form method="post" action="/admin/packages/delete" onsubmit="return confirm('Hapus paket?')">
						<input type="hidden" name="id" value="{{.ID}}">
						<button class="btn danger package-delete-btn user-delete-btn" type="submit">Hapus</button>
					</form>
				</div>
			</td>
		</tr>
		{{end}}
	</table>
	</div>
	{{else}}
	<div class="package-empty-v3">
		{{if .HasRouters}}
		Pilih router dulu, lalu klik <b>Lihat List Paket</b>.
		{{else}}
		Belum ada router. Tambahkan router dulu sebelum membuat paket.
		{{end}}
	</div>
	{{end}}
</div>
{{end}}`, map[string]any{
		"Routers":          opts,
		"List":             list,
		"Edit":             edit,
		"SelectedRouterID": selectedRouterID,
		"RouterLocked":     routerLocked,
		"HasRouters":       hasRouters,
		"CreateMode":       createMode,
		"EditMode":         editMode,
	})

	page(w, "Paket", body)
}

func (a *App) packageSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
		return
	}

	// ADMIN_PACKAGE_SAVE_SCOPE_V1
	role, scopedRouterID := a.currentAdminRole(r)

	_ = r.ParseForm()

	if role == "router" && scopedRouterID > 0 {
		formRouterID, _ := strconv.ParseInt(r.FormValue("router_id"), 10, 64)
		if formRouterID != scopedRouterID {
			http.Error(w, "Forbidden: tidak boleh menyimpan paket router lain", http.StatusForbidden)
			return
		}

		if id := r.FormValue("id"); id != "" && id != "0" {
			pkgID, _ := strconv.ParseInt(id, 10, 64)
			var ownerID int64
			err := a.db.QueryRow(`SELECT router_id FROM packages WHERE id=?`, pkgID).Scan(&ownerID)
			if err != nil || ownerID != scopedRouterID {
				http.Error(w, "Forbidden: paket ini bukan milik router user ini", http.StatusForbidden)
				return
			}
		}
	}

	id := r.FormValue("id")

	price, _ := strconv.ParseInt(r.FormValue("price"), 10, 64)
	su, _ := strconv.Atoi(r.FormValue("shared_users"))

	if su == 0 {
		su = 1
	}

	active := r.FormValue("active") == "on"

	if id == "0" || id == "" {
		_, _ = a.db.Exec(
			`INSERT INTO packages(router_id,name,price,profile,limit_uptime,shared_users,active) VALUES(?,?,?,?,?,?,?)`,
			r.FormValue("router_id"),
			r.FormValue("name"),
			price,
			r.FormValue("profile"),
			r.FormValue("limit_uptime"),
			su,
			active,
		)
	} else {
		_, _ = a.db.Exec(
			`UPDATE packages SET router_id=?,name=?,price=?,profile=?,limit_uptime=?,shared_users=?,active=? WHERE id=?`,
			r.FormValue("router_id"),
			r.FormValue("name"),
			price,
			r.FormValue("profile"),
			r.FormValue("limit_uptime"),
			su,
			active,
			id,
		)
	}

	http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
}

func (a *App) packageDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()

	// ADMIN_PACKAGE_DELETE_SCOPE_V1
	id, _ := strconv.ParseInt(strings.TrimSpace(r.FormValue("id")), 10, 64)
	if !a.adminCanAccessPackage(r, id) {
		http.Error(w, "Forbidden: paket ini bukan milik router user ini", http.StatusForbidden)
		return
	}

	_, _ = a.db.Exec(`DELETE FROM packages WHERE id=?`, id)
	http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
}
