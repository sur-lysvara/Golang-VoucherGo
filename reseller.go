package main

import (
	"net/http"
	"strings"
)

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		slug := strings.Trim(r.URL.Path, "/")
		if slug != "" && !strings.Contains(slug, "/") {
			a.publicRouter(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	rows, err := a.db.Query(`SELECT name,slug FROM routers WHERE enabled=1 ORDER BY name`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type Row struct {
		Name string
		Slug string
	}

	var list []Row

	for rows.Next() {
		var x Row
		_ = rows.Scan(&x.Name, &x.Slug)
		list = append(list, x)
	}

	body := renderHTML(`
<div class="hero">
	<h1 style="margin:0">{{.Site}}</h1>
	<p>Pilih router / lokasi untuk beli voucher.</p>
</div>

<div class="grid" style="margin-top:14px">
{{range .Routers}}
	<a class="card" href="/{{.Slug}}">
		<b>{{.Name}}</b>
		<p class="muted">Lihat paket voucher</p>
	</a>
{{else}}
	<div class="card">Belum ada router aktif.</div>
{{end}}
</div>`, map[string]any{
		"Site":    a.setting("site_name"),
		"Routers": list,
	})

	publicPage(w, a.setting("site_name"), body)
}
