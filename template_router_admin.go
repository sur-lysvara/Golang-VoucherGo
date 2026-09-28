package main

const routerAdminPageTemplate = `


<div class="card">
	<h2>Router</h2>

	{{if .CreateMode}}
		<p class="muted">Tambah router baru.</p>
	{{else if .EditMode}}
		<p class="muted">Edit router #{{.Edit.ID}}.</p>
	{{else}}
		{{if .IsSuper}}
		<div class="actions">
			<a class="btn pri" href="/admin/routers?new=1">Tambah Router</a>
		</div>
		{{end}}
	{{end}}
</div>


{{if .TestMessage}}
<div class="card router-test-result-v1 {{if .TestOK}}ok{{else}}bad{{end}}">
        <b>{{if .TestNotice}}Tidak ada perubahan{{else if .TestOK}}Test koneksi berhasil{{else}}Test koneksi gagal{{end}}</b>
        <div>{{.TestMessage}}</div>
</div>
{{end}}

{{if .CreateMode}}
<div class="card">
	<h2>Tambah Router</h2>

	<form method="post" action="/admin/routers/save" class="router-form-compact-v2 router-api-form-v2">
		<input type="hidden" name="id" value="0">

		<div class="grid">
			<div class="field"><label>Nama</label><input name="name" required></div>
			<div class="field"><label>Slug URL</label><input name="slug" placeholder="x86 / krw / kobak" required></div>
				<div class="field"><label>URL Login WiFi</label><input name="hotspot_login_url" placeholder="http://192.168.88.1/login"></div>
			<div class="field"><label>Host/IP MikroTik</label><input name="host" placeholder="103.x.x.x"></div>
			<div class="field"><label>Port</label><input name="port" value="{{.Edit.Port}}" placeholder="Contoh: 443 / 6665"></div>
						<div class="field"><label>Mode API</label><select class="router-api-select-v3" name="api_mode">
							{{range .APIModes}}<option value="{{.Value}}" {{if .Selected}}selected{{end}}>{{.Label}}</option>{{end}}
						</select></div>
			<div class="field"><label>Username</label><input name="username"></div>
			<div class="field"><label>Password</label><input name="password"></div>
		</div>

		<div class="router-checks">
			<label>HTTPS REST/TLS <input type="checkbox" name="use_ssl" {{checked .Edit.UseSSL}}></label>
			<label>Allow self-signed TLS <input type="checkbox" name="insecure_ssl" {{checked .Edit.InsecureSSL}}></label>
			<label>Aktif <input type="checkbox" name="enabled" {{checked .Edit.Enabled}}></label>
		</div>

		<br>
		<div class="actions">
			<button class="btn pri" type="submit">Simpan Router</button>
			<a class="btn light" href="/admin/routers">Batal</a>
		</div>
	</form>
</div>
{{else if .EditMode}}
<div class="card">
	<h2>Edit Router</h2>

	<form method="post" action="/admin/routers/save" class="router-form-compact-v2 router-api-form-v2">
		<input type="hidden" name="id" value="{{.Edit.ID}}">

		<div class="grid">
			<div class="field"><label>Nama</label><input name="name" value="{{.Edit.Name}}" required></div>
			<div class="field"><label>Slug URL</label><input name="slug" value="{{.Edit.Slug}}" required></div>
						<div class="field"><label>URL Login WiFi</label><input name="hotspot_login_url" value="{{.Edit.HotspotLoginURL}}" placeholder="http://192.168.88.1/login"></div>
			<div class="field"><label>Host/IP MikroTik</label><input name="host" value="{{.Edit.Host}}"></div>
			<div class="field"><label>Port</label><input name="port" value="{{.Edit.Port}}" placeholder="Contoh: 443 / 6665"></div>
						<div class="field"><label>Mode API</label><select class="router-api-select-v3" name="api_mode">
							{{range .APIModes}}<option value="{{.Value}}" {{if .Selected}}selected{{end}}>{{.Label}}</option>{{end}}
						</select></div>
			<div class="field"><label>Username</label><input name="username" value="{{.Edit.Username}}"></div>
			<div class="field"><label>Password</label><input name="password" placeholder="Kosongkan jika tidak diubah"></div>
		</div>

		<div class="router-checks">
			<label>HTTPS REST/TLS <input type="checkbox" name="use_ssl" {{checked .Edit.UseSSL}}></label>
			<label>Allow self-signed TLS <input type="checkbox" name="insecure_ssl" {{checked .Edit.InsecureSSL}}></label>
			<label>Aktif <input type="checkbox" name="enabled" {{checked .Edit.Enabled}}></label>
		</div>

		<br>
		<div class="actions">
			<button class="btn ok" type="submit">Simpan Perubahan</button>
                        <button class="btn light" type="submit" formaction="/admin/routers/test" formmethod="post">Test Koneksi</button>
			<a class="btn light" href="/admin/routers">Batal</a>
		</div>
	</form>
</div>
{{else}}
<div class="card">
	<h2>List Router</h2>

	<div class="tablewrap">
	<table class="routers-readonly-table">
		<tr>
			<th>Nama</th>
			<th>URL</th>
			<th>Host</th>
			<th>Status</th>
			{{if .IsSuper}}<th>Aksi</th>{{end}}
		</tr>


	{{range .List}}
		<tr>
			<td><b>{{.Name}}</b></td>
			<td><a href="/{{.Slug}}" target="_blank">/{{.Slug}}</a>{{if .HotspotLoginURL}}<br><small>Login: {{.HotspotLoginURL}}</small>{{end}}</td>
			<td>{{.Host}}:{{.Port}}<br><small>API: {{.APIMode}}</small></td>
			<td>
                                {{if .Enabled}}<span class="badge success">config aktif</span>{{else}}<span class="badge failed">config off</span>{{end}}
                                {{if .LastTestOK}}<br><span class="badge success">koneksi OK</span>{{else if .LastTestFail}}<br><span class="badge failed">koneksi gagal</span>{{else}}<br><small>koneksi belum dites</small>{{end}}
                                {{if .LastTestAt}}<br><small>{{.LastTestAt}}</small>{{end}}
                        </td>
			{{if $.IsSuper}}
			<td>
				<div class="actions">
					<a class="btn light" href="/admin/routers?edit={{.ID}}">Edit</a>
                                        <form class="router-test-list-form" method="post" action="/admin/routers/test">
                                                <input type="hidden" name="id" value="{{.ID}}">
                                                <input type="hidden" name="return_to" value="/admin/routers">
                                                <button class="btn light" type="submit">Test</button>
                                        </form>
					<form method="post" action="/admin/routers/delete" onsubmit="return confirm('Hapus router?')">
						<input type="hidden" name="id" value="{{.ID}}">
						<button class="btn danger" type="submit">Hapus</button>
					</form>
				</div>
			</td>
			{{end}}
		</tr>
		{{end}}
	</table>
	</div>
</div>
{{end}}`
