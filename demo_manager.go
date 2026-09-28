package main

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func runDemoManagerCmd(args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", args...)
	out, err := cmd.CombinedOutput()

	var b strings.Builder
	if len(out) > 0 {
		b.Write(out)
	}
	if ctx.Err() == context.DeadlineExceeded {
		b.WriteString("\nERROR: command timeout")
	} else if err != nil {
		b.WriteString("\nERROR: " + err.Error())
	}

	return strings.TrimSpace(b.String())
}

func demoServiceState(service string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "systemctl", "is-active", service)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(out))
	}
	return strings.TrimSpace(string(out))
}

func (a *App) adminDemo(w http.ResponseWriter, r *http.Request) {
	if !requireOwnerPanelRequest(w, r) {
		return
	}

	msg := ""

	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		action := strings.TrimSpace(r.FormValue("action"))

		switch action {
		case "status":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-status.sh")
		case "restart_app":
			msg = runDemoManagerCmd("systemctl", "restart", "tuku-demo")
			if msg == "" {
				msg = "OK: tuku-demo direstart."
			}
		case "restart_wa":
			msg = runDemoManagerCmd("systemctl", "restart", "wa-webjs-demo")
			if msg == "" {
				msg = "OK: wa-webjs-demo direstart."
			}
		case "wa_private":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-wa-private.sh")
		case "wa_public":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-wa-public.sh")
		case "reset_demo":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-reset.sh")
		case "pin_status":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-pin.sh", "status")
		case "pin_show":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-pin.sh", "show")
		case "pin_set":
			newPin := strings.TrimSpace(r.FormValue("demo_pin"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-pin.sh", "set", newPin)
		case "pin_delete":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-pin.sh", "delete")
		case "user_list":
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "list")
		case "user_create":
			username := strings.TrimSpace(r.FormValue("demo_user_username"))
			password := strings.TrimSpace(r.FormValue("demo_user_password"))
			role := strings.TrimSpace(r.FormValue("demo_user_role"))
			enabled := strings.TrimSpace(r.FormValue("demo_user_enabled"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "create", username, password, role, enabled)
		case "user_update":
			id := strings.TrimSpace(r.FormValue("demo_user_id"))
			username := strings.TrimSpace(r.FormValue("demo_user_username"))
			role := strings.TrimSpace(r.FormValue("demo_user_role"))
			enabled := strings.TrimSpace(r.FormValue("demo_user_enabled"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "update", id, username, role, enabled)
		case "user_passwd":
			id := strings.TrimSpace(r.FormValue("demo_user_id"))
			password := strings.TrimSpace(r.FormValue("demo_user_password"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "passwd", id, password)
		case "user_enable":
			id := strings.TrimSpace(r.FormValue("demo_user_id"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "enable", id)
		case "user_disable":
			id := strings.TrimSpace(r.FormValue("demo_user_id"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "disable", id)
		case "user_delete":
			id := strings.TrimSpace(r.FormValue("demo_user_id"))
			msg = runDemoManagerCmd("/usr/local/sbin/vcr-demo-user", "delete", id)
		default:
			msg = "Action tidak dikenal."
		}
	}

	demoBaseURL := strings.TrimRight(env("DEMO_PUBLIC_URL", "https://demo.vouchergo.biz.id"), "/")

	data := map[string]string{
		"Msg":       msg,
		"AppState":  demoServiceState("tuku-demo"),
		"WAState":   demoServiceState("wa-webjs-demo"),
		"DemoURL":   demoBaseURL + "/",
		"DemoWAURL": demoBaseURL + "/admin/whatsapp",
		"DemoAdmin": demoBaseURL + "/admin",
	}

	const tpl = `
<style>
/* DEMO_MANAGER_UI_POLISH_V1 */
.demo-manager-wrap{
	display:grid;
	gap:16px;
}
.demo-manager-hero{
	border:1px solid rgba(148,163,184,.22);
	background:linear-gradient(135deg,rgba(59,130,246,.10),rgba(16,185,129,.08));
}
.demo-manager-hero h2{
	margin-bottom:6px;
}
.demo-status-grid{
	display:grid;
	grid-template-columns:repeat(auto-fit,minmax(180px,1fr));
	gap:12px;
	margin-top:14px;
}
.demo-status-box{
	border:1px solid rgba(148,163,184,.22);
	border-radius:16px;
	padding:12px;
	background:rgba(255,255,255,.04);
}
.demo-status-box label{
	display:block;
	font-size:12px;
	opacity:.75;
	margin-bottom:6px;
}
.demo-status-box b{
	font-size:15px;
}
.demo-action-grid{
	display:grid!important;
	grid-template-columns:repeat(auto-fit,minmax(170px,1fr))!important;
	gap:10px!important;
}
.demo-mini-title{
	font-size:14px;
	font-weight:800;
	margin:14px 0 8px;
	opacity:.9;
}
.demo-link-actions{
	display:flex;
	gap:10px;
	flex-wrap:wrap;
	margin-top:14px;
}
.demo-manager-wrap pre{
	max-height:320px;
	border-radius:14px;
	padding:12px;
	background:rgba(15,23,42,.08);
}
@media(max-width:640px){
	.demo-action-grid{
		grid-template-columns:1fr!important;
	}
	.demo-link-actions .btn{
		width:100%;
		text-align:center;
	}
}
</style>
<div class="demo-manager-wrap">
<div class="card settings-page-card demo-manager-hero">
	<h2>Demo Manager</h2>
	<p class="muted">Kelola aplikasi demo VoucherGo dari owner panel.</p>

	<div class="demo-status-grid">
		<div class="demo-status-box">
			<label>Demo App</label>
			<b>{{.AppState}}</b>
		</div>
		<div class="demo-status-box">
			<label>Demo WhatsApp</label>
			<b>{{.WAState}}</b>
		</div>
	</div>

	<div class="demo-mini-title">Kontrol Demo</div>
	<form method="post" class="settings-form-compact-v1" style="margin-top:0">
		<div class="actions settings-actions demo-action-grid">
			<button class="btn ok" name="action" value="status" type="submit">Cek Status Demo</button>
			<button class="btn pri" name="action" value="restart_app" type="submit" onclick="return confirm('Restart aplikasi demo?')">Restart Demo App</button>
			<button class="btn pri" name="action" value="restart_wa" type="submit" onclick="return confirm('Restart WhatsApp demo?')">Restart Demo WA</button>
			<button class="btn ok" name="action" value="wa_private" type="submit" onclick="return confirm('Aktifkan private trial WA demo?')">WA Private Trial</button>
			<button class="btn warn" name="action" value="wa_public" type="submit" onclick="return confirm('Kunci WA demo untuk public?')">WA Public Aman</button>
			<button class="btn danger" name="action" value="reset_demo" type="submit" onclick="return confirm('Reset data demo? Lanjutkan?')">Reset Demo</button>
		</div>
	</form>

	<div class="demo-link-actions">
		<a class="btn" href="{{.DemoWAURL}}" target="_blank" rel="noopener">Buka WA / QR Demo</a>
		<a class="btn" href="{{.DemoURL}}" target="_blank" rel="noopener">Buka Halaman Public Demo</a>
		<a class="btn" href="{{.DemoAdmin}}" target="_blank" rel="noopener">Buka Admin Demo</a>
	</div>

	<div class="card settings-page-card" style="margin-top:16px">
		<h2>User Demo</h2>
		<p class="muted">Cek, tambah, edit, aktif/nonaktifkan, ganti password, atau hapus user login demo.</p>

		<form method="post" class="settings-form-compact-v1" style="margin-top:14px">
			<div class="grid">
				<div class="field">
					<label>ID User</label>
					<input name="demo_user_id" placeholder="isi untuk edit/delete/password">
				</div>
				<div class="field">
					<label>Username</label>
					<input name="demo_user_username" placeholder="contoh: admin-demo">
				</div>
				<div class="field">
					<label>Password Baru</label>
					<input name="demo_user_password" type="password" autocomplete="new-password" placeholder="isi untuk tambah/ganti password">
				</div>
				<div class="field">
					<label>Role</label>
					<select name="demo_user_role">
						<option value="superadmin">superadmin</option>
						<option value="admin">admin</option>
						<option value="operator">operator</option>
						<option value="router">router</option>
					</select>
				</div>
				<div class="field">
					<label>Status</label>
					<select name="demo_user_enabled">
						<option value="1">Aktif</option>
						<option value="0">Nonaktif</option>
					</select>
				</div>
			</div>

			<div class="actions settings-actions demo-action-grid" style="margin-top:12px">
				<button class="btn ok" name="action" value="user_list" type="submit">Cek User Demo</button>
				<button class="btn pri" name="action" value="user_create" type="submit">Tambah User</button>
				<button class="btn" name="action" value="user_update" type="submit">Edit User</button>
				<button class="btn" name="action" value="user_passwd" type="submit">Ganti Password</button>
				<button class="btn ok" name="action" value="user_enable" type="submit">Aktifkan</button>
				<button class="btn warn" name="action" value="user_disable" type="submit">Nonaktifkan</button>
				<button class="btn danger" name="action" value="user_delete" type="submit" onclick="return confirm('Hapus user demo ini?')">Delete User</button>
			</div>
		</form>
	</div>

	<div class="card settings-page-card" style="margin-top:16px">
		<h2>PIN WhatsApp Demo</h2>
		<p class="muted">Atur PIN untuk membuka mode Private Trial WA di halaman demo.</p>

		<form method="post" class="settings-form-compact-v1" style="margin-top:14px">
			<div class="grid">
				<div class="field">
					<label>PIN Baru</label>
					<input name="demo_pin" autocomplete="off" placeholder="contoh: 123456">
				</div>
			</div>

			<div class="actions settings-actions demo-action-grid" style="margin-top:12px">
				<button class="btn ok" name="action" value="pin_status" type="submit">Cek PIN</button>
				<button class="btn" name="action" value="pin_show" type="submit" onclick="return confirm('Tampilkan PIN demo sekarang?')">Tampilkan PIN</button>
				<button class="btn pri" name="action" value="pin_set" type="submit">Simpan PIN Baru</button>
				<button class="btn danger" name="action" value="pin_delete" type="submit" onclick="return confirm('Hapus PIN demo? Private Trial WA tidak bisa dibuka sampai PIN diset lagi.')">Hapus PIN</button>
			</div>
		</form>
	</div>

	{{if .Msg}}
		<div class="notice" style="margin-top:16px">
			<b>Output:</b>
			<pre style="white-space:pre-wrap;overflow:auto;margin:10px 0 0">{{.Msg}}</pre>
		</div>
	{{end}}
</div>
</div>
`

	t, err := template.New("demo-manager").Parse(tpl)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	page(w, "Demo Manager", template.HTML(b.String()))
}
