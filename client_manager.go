package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func isSafeClientSlug(s string) bool {
	if len(s) < 3 || len(s) > 20 {
		return false
	}

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '-' {
			continue
		}
		return false
	}

	return true
}

func isSafeClientDomain(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))

	if len(s) < len("a.vouchergo.biz.id") || len(s) > 120 {
		return false
	}

	if strings.Contains(s, "://") || strings.Contains(s, "/") || strings.Contains(s, " ") || strings.Contains(s, "_") {
		return false
	}

	const suffix = ".vouchergo.biz.id"
	if !strings.HasSuffix(s, suffix) {
		return false
	}

	prefix := strings.TrimSuffix(s, suffix)
	if prefix == "" || len(prefix) > 40 {
		return false
	}

	if strings.HasPrefix(prefix, "-") || strings.HasSuffix(prefix, "-") || strings.Contains(prefix, "--") {
		return false
	}

	for _, r := range prefix {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '-' {
			continue
		}
		return false
	}

	return true
}

func defaultClientDomain(slug string) string {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if !isSafeClientSlug(slug) {
		slug = "client-a"
	}
	return slug + ".vouchergo.biz.id"
}

func isSafeClientPortText(s string) bool {
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}

	if n < 8111 || n > 65535 {
		return false
	}

	switch n {
	case 8097, 8098, 8099, 8100, 8666:
		return false
	}

	return true
}

func isSafeReleaseTag(s string) bool {
	if len(s) < 2 || len(s) > 40 {
		return false
	}

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '.' || r == '-' || r == '_' {
			continue
		}
		return false
	}

	return true
}

func runVCRClientScript(script string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	all := append([]string{"-n", script}, args...)
	cmd := exec.CommandContext(ctx, "/usr/bin/sudo", all...)
	out, err := cmd.CombinedOutput()

	text := strings.TrimSpace(string(out))

	if ctx.Err() == context.DeadlineExceeded {
		if text != "" {
			text += "\n"
		}
		text += "ERROR: command timeout"
		return text
	}

	if err != nil {
		if text != "" {
			text += "\n"
		}
		text += "ERROR: " + err.Error()
	}

	if strings.TrimSpace(text) == "" {
		text = "(tidak ada output)"
	}

	return text
}

func clientOnboardingPath(slug string) string {
	return "/home/xyzsur/VCR-CLIENT-" + slug + "-ONBOARDING.txt"
}

func (a *App) adminClientOnboardingDownload(w http.ResponseWriter, r *http.Request) {
	if !requireOwnerPanelRequest(w, r) {
		return
	}

	// ADMIN_CLIENT_MANAGER_ONBOARDING_DOWNLOAD_V1
	if env("OWNER_PANEL", "0") != "1" {
		http.NotFound(w, r)
		return
	}

	slug := strings.TrimSpace(r.URL.Query().Get("slug"))
	if !isSafeClientSlug(slug) {
		http.Error(w, "slug tidak valid", http.StatusBadRequest)
		return
	}

	_ = runVCRClientScript("/usr/local/sbin/vcr-client-onboarding.sh", slug)

	path := clientOnboardingPath(slug)
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "onboarding file tidak ditemukan", http.StatusNotFound)
		return
	}

	filename := "VoucherGo-CLIENT-" + slug + "-ONBOARDING.txt"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	_, _ = w.Write(data)
}

// CLIENT_CREATE_AUTODEFAULTS_V3_HELPERS
func vcrLatestReleaseTagAutoV3() string {
	// CLIENT_LATEST_TAG_SAFE_GIT_V1
	bestTag := ""
	bestN := -1

	consider := func(tag string) {
		tag = strings.TrimSpace(tag)
		tag = strings.TrimSuffix(tag, "^{}")
		tag = strings.TrimPrefix(tag, "refs/tags/")
		if !strings.HasPrefix(tag, "v1.0-rc") {
			return
		}

		n, err := strconv.Atoi(strings.TrimPrefix(tag, "v1.0-rc"))
		if err != nil {
			return
		}

		if n > bestN {
			bestN = n
			bestTag = tag
		}
	}

	candidates := []string{}
	if exe, err := os.Executable(); err == nil && exe != "" {
		candidates = append(candidates, filepath.Dir(exe))
	}
	candidates = append(candidates, "/var/www/vouchergo", ".")

	seen := map[string]bool{}
	for _, dir := range candidates {
		dir = strings.TrimSpace(dir)
		if dir == "" || seen[dir] {
			continue
		}
		seen[dir] = true

		// Pakai safe.directory supaya tetap bisa baca repo walau service user beda ownership.
		cmd := exec.Command("git", "-c", "safe.directory="+dir, "-C", dir, "tag", "--list", "v1.0-rc*", "--sort=-version:refname")
		out, err := cmd.Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				consider(line)
			}
		}

		gitDir := filepath.Join(dir, ".git")

		// Fallback 1: loose refs.
		files, _ := filepath.Glob(filepath.Join(gitDir, "refs", "tags", "v1.0-rc*"))
		for _, file := range files {
			consider(filepath.Base(file))
		}

		// Fallback 2: packed refs.
		if b, err := os.ReadFile(filepath.Join(gitDir, "packed-refs")); err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "^") {
					continue
				}
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					consider(parts[1])
				}
			}
		}
	}

	if bestTag != "" {
		return bestTag
	}

	return "main"
}

func vcrUsedPortsFromEnvAutoV3() map[int]bool {
	used := map[int]bool{}

	files, _ := filepath.Glob("/etc/tuku-*.env")
	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

			if key != "PORT" && key != "APP_PORT" && key != "WA_PORT" && key != "WA_LOCAL_PORT" {
				continue
			}

			n, err := strconv.Atoi(val)
			if err == nil && n > 0 {
				used[n] = true
			}
		}
	}

	return used
}

func vcrPortBusyAutoV3(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return true
	}
	_ = ln.Close()
	return false
}

func vcrNextFreeClientPortsAutoV3() (int, int) {
	used := vcrUsedPortsFromEnvAutoV3()

	for appPort := 8111; appPort <= 8997; appPort += 2 {
		waPort := appPort + 1

		if used[appPort] || used[waPort] {
			continue
		}
		if vcrPortBusyAutoV3(appPort) || vcrPortBusyAutoV3(waPort) {
			continue
		}

		return appPort, waPort
	}

	return 0, 0
}

func (a *App) adminClients(w http.ResponseWriter, r *http.Request) {
	// CLIENT_CREATE_AUTODEFAULTS_V3_ACTIVE
	defaultReleaseTag := vcrLatestReleaseTagAutoV3()
	defaultAppPort, defaultWAPort := vcrNextFreeClientPortsAutoV3()

	if !requireOwnerPanelRequest(w, r) {
		return
	}

	// ADMIN_CLIENT_MANAGER_STAGE1_V1
	// ADMIN_CLIENT_MANAGER_CREATE_V1
	// ADMIN_CLIENT_MANAGER_STAGE3_ACTIONS_V1
	// ADMIN_CLIENT_MANAGER_STAGE4_DELETE_V1
	// ADMIN_CLIENT_MANAGER_STAGE5_SSL_V1
	// ADMIN_CLIENT_MANAGER_STAGE6_ONBOARDING_V1
	if env("OWNER_PANEL", "0") != "1" {
		http.NotFound(w, r)
		return
	}

	msg := ""
	slug := strings.TrimSpace(r.URL.Query().Get("slug"))
	domain := ""
	appPort := "8111"
	waPort := "8112"
	tag := "v1.0-rc7"
	output := ""

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		action := strings.TrimSpace(r.FormValue("action"))
		slug = strings.TrimSpace(r.FormValue("slug"))
		domain = strings.ToLower(strings.TrimSpace(r.FormValue("domain")))
		if domain == "" && isSafeClientSlug(slug) {
			domain = defaultClientDomain(slug)
		}
		confirmText := strings.TrimSpace(r.FormValue("confirm_text"))
		appPort = strings.TrimSpace(r.FormValue("app_port"))
		waPort = strings.TrimSpace(r.FormValue("wa_port"))
		tag = strings.TrimSpace(r.FormValue("tag"))

		if tag == "" {
			tag = "v1.0-rc7"
		}

		switch action {
		case "list":
			output = runVCRClientScript("/usr/local/sbin/vcr-client-list.sh")

		case "status":
			if !isSafeClientSlug(slug) {
				msg = "Slug client tidak valid. Contoh: client-a"
			} else {
				output = runVCRClientScript("/usr/local/sbin/vcr-client-status.sh", slug)
			}

		case "restart_app", "restart_wa", "wa_status", "info", "vip_branding_on", "vip_branding_off", "vip_branding_status":
			if !isSafeClientSlug(slug) {
				msg = "Slug client tidak valid. Contoh: client-a"
			} else {
				scriptAction := map[string]string{
					"restart_app":         "restart-app",
					"restart_wa":          "restart-wa",
					"wa_status":           "wa-status",
					"info":                "info",
					"vip_branding_on":     "vip-branding-on",
					"vip_branding_off":    "vip-branding-off",
					"vip_branding_status": "vip-branding-status",
					"vip-branding-on":     "vip-branding-on",
					"vip-branding-off":    "vip-branding-off",
				}[action]
				output = runVCRClientScript("/usr/local/sbin/vcr-client-action.sh", slug, scriptAction)
			}

		case "onboarding":
			if !isSafeClientSlug(slug) {
				msg = "Slug client tidak valid. Contoh: client-a"
			} else {
				gen := runVCRClientScript("/usr/local/sbin/vcr-client-onboarding.sh", slug)
				data, err := os.ReadFile(clientOnboardingPath(slug))
				if err != nil {
					output = gen + "\n\nERROR: gagal baca onboarding file: " + err.Error()
				} else {
					output = string(data)
				}
			}

		case "send_onboarding":
			phone := strings.TrimSpace(r.FormValue("phone"))
			if !isSafeClientSlug(slug) {
				msg = "Slug client tidak valid. Contoh: client-a"
			} else if !isSafePhoneText(phone) {
				msg = "Nomor WhatsApp tidak valid."
			} else {
				output = runVCRClientScript("/usr/local/sbin/vcr-client-onboarding-send.sh", slug, phone)
			}

		case "dns_check", "install_ssl":
			if !isSafeClientSlug(slug) {
				msg = "Slug client tidak valid. Contoh: client-a"
			} else {
				scriptAction := map[string]string{
					"dns_check":   "check-dns",
					"install_ssl": "install-ssl",
				}[action]
				output = runVCRClientScript("/usr/local/sbin/vcr-client-ssl.sh", slug, scriptAction)
			}

		case "create":
			switch {
			case !isSafeClientSlug(slug):
				msg = "Slug client tidak valid. Gunakan huruf kecil, angka, dash. Contoh: client-a"
			case !isSafeClientDomain(domain):
				msg = "Domain tidak valid. Contoh: client-a.vouchergo.biz.id"
			case !isSafeClientPortText(appPort):
				msg = "App port tidak valid. Gunakan port >= 8111."
			case !isSafeClientPortText(waPort):
				msg = "WA port tidak valid. Gunakan port >= 8111."
			case appPort == waPort:
				msg = "App port dan WA port tidak boleh sama."
			case !isSafeReleaseTag(tag):
				msg = "Tag release tidak valid."
			default:
				output = runVCRClientScript("/usr/local/sbin/vcr-create-client.sh", slug, domain, appPort, waPort, tag)
			}

		case "delete":
			expectedConfirm := "HAPUS " + slug
			switch {
			case !isSafeClientSlug(slug):
				msg = "Slug client tidak valid. Contoh: client-a"
			case !isSafeClientDomain(domain):
				msg = "Domain tidak valid. Contoh: client-a.vouchergo.biz.id"
			case confirmText != expectedConfirm:
				msg = "Konfirmasi tidak cocok. Ketik persis: " + expectedConfirm
			default:
				output = runVCRClientScript("/usr/local/sbin/vcr-delete-client.sh", slug, domain, "--yes")
			}

		default:
			msg = "Action tidak dikenal."
		}
	} else {
		output = runVCRClientScript("/usr/local/sbin/vcr-client-list.sh")
	}

	if domain == "" {
		domain = defaultClientDomain(slug)
	}

	body := renderHTML(`
<div class="card clients-manager">
	<style>
/* CLIENT_MANAGER_UI_POLISH_V2 */
.client-manager-note{
	margin:6px 0 14px;
}
.client-manager-actions,
.client-action-grid{
	display:grid!important;
	grid-template-columns:repeat(auto-fit,minmax(145px,1fr))!important;
	gap:8px!important;
	align-items:stretch!important;
}
.mini-stack{
	display:grid!important;
	gap:8px!important;
	margin:0!important;
}
.mini-stack .btn,
.client-manager-actions .btn,
.client-action-grid .btn{
	width:100%;
	min-height:38px;
	display:inline-flex;
	align-items:center;
	justify-content:center;
	text-align:center;
	white-space:normal;
	line-height:1.2;
}
.btn.vip-brand-on{
	background:linear-gradient(135deg,#16a34a,#22c55e)!important;
	color:#fff!important;
	border:0!important;
}
.btn.vip-brand-off{
	background:linear-gradient(135deg,#dc2626,#ef4444)!important;
	color:#fff!important;
	border:0!important;
}
.client-manager-badge{
	display:inline-flex;
	align-items:center;
	gap:6px;
	padding:5px 9px;
	border-radius:999px;
	font-size:12px;
	font-weight:800;
	background:rgba(59,130,246,.10);
	border:1px solid rgba(59,130,246,.20);
}
.tablewrap table td{
	vertical-align:top;
}
.tablewrap td form + form{
	margin-top:8px;
}
@media(max-width:760px){
	.client-manager-actions,
	.client-action-grid{
		grid-template-columns:1fr!important;
	}
	.tablewrap{
		border-radius:18px;
	}
	.tablewrap table{
		font-size:13px;
	}
	.mini-stack .btn{
		min-height:42px;
	}
}
</style>

<h2>Client Manager</h2>

<style>
/* CLIENT_MANAGER_PAGED_UI_V1 */
.client-page-tabs{
	display:flex;
	gap:8px;
	flex-wrap:wrap;
	margin:12px 0 16px;
}
.client-page-tabs a{
	display:inline-flex;
	align-items:center;
	justify-content:center;
	min-height:36px;
	padding:8px 12px;
	border-radius:999px;
	text-decoration:none;
	font-weight:800;
	border:1px solid rgba(148,163,184,.28);
	background:rgba(255,255,255,.04);
	color:inherit;
}
.client-page-tabs a.active{
	background:linear-gradient(135deg,#2563eb,#1d4ed8);
	color:#fff;
	border-color:transparent;
}
.client-page-hidden{
	display:none!important;
}
.client-page-section{
	animation:clientFade .16s ease-out;
}
@keyframes clientFade{
	from{opacity:.45;transform:translateY(4px)}
	to{opacity:1;transform:none}
}
.client-page-tip{
	margin:0 0 14px;
	padding:10px 12px;
	border-radius:14px;
	border:1px solid rgba(148,163,184,.20);
	background:rgba(59,130,246,.08);
	font-size:13px;
}
@media(max-width:720px){
	.client-page-tabs{
		display:grid;
		grid-template-columns:1fr 1fr;
	}
	.client-page-tabs a{
		width:100%;
	}
}
</style>
<div class="client-page-tabs" id="clientPageTabs">
	<a href="/admin/clients?tab=list" data-tab="list">List</a>
	<a href="/admin/clients?tab=status" data-tab="status">Status</a>
	<a href="/admin/clients?tab=action" data-tab="action">Action</a>
	<a href="/admin/clients?tab=create" data-tab="create">Buat Client</a>
	<a href="/admin/clients?tab=onboarding" data-tab="onboarding">Onboarding</a>
  <a class="client-manager-billing-entry-v1" href="/admin/client-billing">Billing Client</a>
	<a href="/admin/clients?tab=delete" data-tab="delete">Hapus</a>
</div>
<div class="client-page-tip" id="clientPageTip"></div>
<script>
document.addEventListener('DOMContentLoaded', function(){
	const params = new URLSearchParams(location.search);
	const tab = params.get('tab') || 'list';

	const tips = {
		list: 'Menampilkan daftar client aktif.',
		status: 'Cek status service/app client berdasarkan slug.',
		action: 'Restart app/WA, cek WA, QR, VIP Branding, DNS, dan SSL.',
		create: 'Buat client baru beserta domain, port, dan release tag.',
		onboarding: 'Lihat, download, atau kirim onboarding client.',
		delete: 'Hapus client. Gunakan hati-hati.'
	};

	document.querySelectorAll('#clientPageTabs a').forEach(function(a){
		if(a.dataset.tab === tab) a.classList.add('active');
	});

	const tip = document.getElementById('clientPageTip');
	if(tip) tip.textContent = tips[tab] || tips.list;

	const map = {
		'List Client': 'list',
		'Status Client': 'status',
		'Quick Action': 'action',
		'Buat Client': 'create',
		'Onboarding Client': 'onboarding',
		'Hapus Client': 'delete'
	};

	document.querySelectorAll('.card').forEach(function(card){
		const h = card.querySelector('h2,h3');
		if(!h) return;
		const title = h.textContent.trim();
		const section = map[title];
		if(!section) return;

		card.classList.add('client-page-section');
		if(section !== tab) card.classList.add('client-page-hidden');
	});

	// Output tetap tampil kalau ada isi hasil command
	document.querySelectorAll('.card').forEach(function(card){
		const h = card.querySelector('h2,h3');
		if(h && h.textContent.trim() === 'Output'){
			const txt = card.textContent.replace('Output','').trim();
			if(!txt) card.classList.add('client-page-hidden');
		}
	});
});
</script>

	<p class="muted client-manager-note">Stage 2: list, status, create, action, onboarding, billing, dan hapus client dengan konfirmasi.</p>

	{{if .Msg}}
	<div class="notice">{{.Msg}}</div>
	{{end}}

	<div class="grid clients-grid">
		<form method="post" class="card client-tool-card">
			<h3>List Client</h3>
			<input type="hidden" name="action" value="list">
			<p class="muted">Tampilkan semua client aktif.</p>
			<button class="btn pri" type="submit">Refresh List</button>
		</form>

		<form method="post" class="card client-tool-card">
			<h3>Status Client</h3>
			<input type="hidden" name="action" value="status">
			<div class="field">
				<label>Slug Client</label>
				<input name="slug" value="{{.Slug}}" placeholder="client-a">
			</div>
			<button class="btn ok" type="submit">Cek Status</button>
		</form>

		<form method="post" class="card client-tool-card client-action-card">
			<h3>Quick Action</h3>
			<div class="field">
				<label>Slug Client</label>
				<input name="slug" value="{{.Slug}}" placeholder="client-a" required>
			</div>

			<div class="client-action-grid">
				<button class="btn pri" name="action" value="restart_app" type="submit" onclick="return confirm('Restart app client ini?')">Restart App</button>
				<button class="btn pri" name="action" value="restart_wa" type="submit" onclick="return confirm('Restart WA client ini?')">Restart WA</button>
				<button class="btn ok" name="action" value="wa_status" type="submit">WA Status</button>
				<button class="btn" name="action" value="info" type="submit">Info / QR</button>
				<button class="btn ok" name="action" value="vip_branding_status" type="submit">Cek VIP</button>
				<button class="btn vip-brand-on" name="action" value="vip_branding_on" type="submit" onclick="return confirm('Aktifkan VIP Branding untuk client ini?')">VIP Branding ON</button>
				<button class="btn vip-brand-off" name="action" value="vip_branding_off" type="submit" onclick="return confirm('Matikan VIP Branding untuk client ini?')">VIP Branding OFF</button>
				<button class="btn" name="action" value="dns_check" type="submit">DNS Check</button>
				<button class="btn pri" name="action" value="install_ssl" type="submit" onclick="return confirm('Install SSL untuk client ini? Pastikan DNS sudah benar.');">Install SSL</button>
			</div>
		</form>

		<form method="post" class="card client-tool-card client-create-card" onsubmit="return confirm('Buat client baru? Pastikan domain dan port benar.');">
			<h3>Buat Client</h3>
			<input type="hidden" name="action" value="create">

			<div class="field">
				<label>Slug</label>
				<input name="slug" value="{{.Slug}}" placeholder="client-a" required>
			</div>

			<div class="field">
				<label>Domain</label>
				<input name="domain" value="{{.Domain}}" placeholder="client-a.vouchergo.biz.id" required>
							<small>Wajib format: <b>slug.vouchergo.biz.id</b></small>
			</div>

			<div class="grid mini-grid">
				<div class="field">
					<label>App Port</label>
					<input name="app_port" value="{{.DefaultAppPort}}" required>
				</div>

				<div class="field">
					<label>WA Port</label>
					<input name="wa_port" value="{{.DefaultWAPort}}" required>
				</div>
			</div>

			<div class="field">
				<label>Release Tag</label>
				<input name="tag" value="{{.Tag}}" required>
<div class="client-autodefault-warning-v1">
	<b>⚠️ Otomatis:</b> Release Tag, App Port, dan WA Port sudah diisi otomatis.
	App Port / WA Port mencari port kosong berikutnya, sedangkan Release Tag mengambil tag terbaru.
	Ubah manual hanya kalau benar-benar perlu.
</div>
<!-- CLIENT_CREATE_AUTODEFAULTS_WARNING_V1 -->

			</div>

			<button class="btn pri" type="submit">Buat Client</button>
		</form>
		<div class="card client-tool-card client-onboarding-card">
			<h3>Onboarding Client</h3>
			<p class="muted">Generate panduan, download file, atau kirim ke WhatsApp client.</p>

			<form method="post" class="mini-stack">
				<input type="hidden" name="action" value="onboarding">
				<div class="field">
					<label>Slug Client</label>
					<input name="slug" value="{{.Slug}}" placeholder="client-a" required>
				</div>
				<button class="btn ok" type="submit">Lihat Onboarding</button>
			</form>

			<form method="get" action="/admin/clients/onboarding-download" class="mini-stack">
				<div class="field">
					<label>Slug Client</label>
					<input name="slug" value="{{.Slug}}" placeholder="client-a" required>
				</div>
				<button class="btn pri" type="submit">Download Onboarding</button>
			</form>

			<form method="post" class="mini-stack" onsubmit="return confirm('Kirim onboarding ke WhatsApp client? File berisi username dan password awal.');">
				<input type="hidden" name="action" value="send_onboarding">
				<div class="field">
					<label>Slug Client</label>
					<input name="slug" value="{{.Slug}}" placeholder="client-a" required>
				</div>
				<div class="field">
					<label>No WhatsApp Client</label>
					<input name="phone" placeholder="08xxxxxxxxxx" required>
				</div>
				<button class="btn ok" type="submit">Kirim WhatsApp ke Client</button>
			</form>
		</div>

		<form method="post" class="card client-tool-card client-delete-card" onsubmit="return confirm('YAKIN hapus client ini? Data, service, database, dan folder client akan dibackup lalu dihapus.');">
			<h3>Hapus Client</h3>
			<input type="hidden" name="action" value="delete">

			<div class="field">
				<label>Slug</label>
				<input name="slug" value="{{.Slug}}" placeholder="client-a" required>
			</div>

			<div class="field">
				<label>Domain</label>
				<input name="domain" value="{{.Domain}}" placeholder="client-a.vouchergo.biz.id" required>
			</div>

			<div class="field">
				<label>Konfirmasi</label>
				<input name="confirm_text" placeholder="HAPUS client-a" required>
			</div>

			<button class="btn danger" type="submit">Hapus Client</button>
		</form>
	</div>
</div>

<div class="card clients-output-card">
	<h2>Output</h2>
	<pre class="client-output">{{.Output}}</pre>
</div>

<style>
.clients-grid{
  grid-template-columns:repeat(3,minmax(0,1fr))!important;
  gap:12px!important;
}
.client-tool-card{
  margin:0!important;
  padding:14px!important;
}
.client-tool-card h3{
  margin-bottom:8px!important;
}
.client-create-card .mini-grid{
  grid-template-columns:1fr 1fr!important;
  gap:8px!important;
}

/* CLIENT_MANAGER_VIP_BRANDING_UI_V1 */
.vip-brand-on{
	background:#16a34a!important;
	color:#fff!important;
}
.vip-brand-off{
	background:#f97316!important;
	color:#fff!important;
}
html[data-admin-theme="dark"] .vip-brand-on{
	background:#15803d!important;
}
html[data-admin-theme="dark"] .vip-brand-off{
	background:#ea580c!important;
}

.client-action-grid{
  display:grid;
  grid-template-columns:1fr 1fr;
  gap:8px;
}


.client-onboarding-card{
  border-color:rgba(34,197,94,.24)!important;
}
.client-onboarding-card .mini-stack{
  display:grid;
  gap:8px;
  margin-top:10px;
  padding-top:10px;
  border-top:1px solid var(--line);
}
.client-onboarding-card .mini-stack:first-of-type{
  border-top:0;
  padding-top:0;
}

.client-delete-card{
  border-color:rgba(220,38,38,.28)!important;
}
.client-delete-card h3{
  color:#dc2626!important;
}
html[data-admin-theme="dark"] .client-delete-card h3{
  color:#fca5a5!important;
}

@media(max-width:760px){
  .client-action-grid{
    grid-template-columns:1fr!important;
  }
}

.client-output{
  max-height:560px;
  overflow:auto;
  white-space:pre-wrap;
  word-break:break-word;
  font-size:12px;
  line-height:1.45;
}
@media(max-width:1024px){
  .clients-grid{
    grid-template-columns:1fr!important;
  }
}
@media(max-width:760px){
  .client-output{
    max-height:460px;
    font-size:10.5px;
  }
}
</style>
`, map[string]any{
		"DefaultAppPort": defaultAppPort,
		"DefaultWAPort":  defaultWAPort,
		"Msg":            msg,
		"Slug":           slug,
		"Domain":         domain,
		"AppPort":        appPort,
		"WAPort":         waPort,
		"Tag":            defaultReleaseTag,
		"Output":         output,
	})

	page(w, "Clients", body)
}
