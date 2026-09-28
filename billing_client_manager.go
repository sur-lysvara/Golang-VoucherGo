package main

import (
	"bytes"
	"net/http"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// OWNER_CLIENT_BILLING_MANAGER_V1
type clientBillingRowV1 struct {
	Slug        string
	Status      string
	DueAt       string
	GraceUntil  string
	Price       string
	PayURL      string
	LastPaidAt  string
	StatusClass string
}

func parseBillingStatusOutputV1(out string) map[string]string {
	m := map[string]string{}

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		m[k] = v
	}

	return m
}

func billingStatusClassV1(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active":
		return "ok"
	case "warning", "grace", "expired":
		return "warn"
	case "suspended":
		return "bad"
	default:
		return "warn"
	}
}

func (a *App) clientBillingRowsV1() []clientBillingRowV1 {
	files, _ := filepath.Glob("/etc/tuku-*.env")
	rows := []clientBillingRowV1{}

	for _, f := range files {
		base := filepath.Base(f)
		if base == "tuku-demo.env" || base == "tuku.env" {
			continue
		}

		slug := strings.TrimSuffix(strings.TrimPrefix(base, "tuku-"), ".env")
		if slug == "" || slug == "demo" {
			continue
		}

		cmd := exec.Command("sudo", "/usr/local/sbin/vcr-client-billing.sh", "status", slug)
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		if cmd.Run() != nil {
			continue
		}

		vals := parseBillingStatusOutputV1(buf.String())
		status := vals["CLIENT_BILLING_STATUS"]
		if status == "" {
			status = "active"
		}

		rows = append(rows, clientBillingRowV1{
			Slug:        slug,
			Status:      status,
			DueAt:       defaultDashV1(vals["CLIENT_BILLING_DUE_AT"]),
			GraceUntil:  defaultDashV1(vals["CLIENT_BILLING_GRACE_UNTIL"]),
			Price:       defaultDashV1(vals["CLIENT_BILLING_PRICE"]),
			PayURL:      defaultDashV1(vals["CLIENT_BILLING_OWNER_URL"]),
			LastPaidAt:  defaultDashV1(vals["CLIENT_BILLING_LAST_PAID_AT"]),
			StatusClass: billingStatusClassV1(status),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Slug < rows[j].Slug
	})

	return rows
}

func defaultDashV1(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "-"
	}
	return v
}

func (a *App) adminClientBilling(w http.ResponseWriter, r *http.Request) {
	if env("OWNER_PANEL", "0") != "1" {
		http.NotFound(w, r)
		return
	}

	msg := strings.TrimSpace(r.URL.Query().Get("msg"))

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		slug := strings.TrimSpace(r.FormValue("slug"))
		action := strings.TrimSpace(r.FormValue("action"))
		months := strings.TrimSpace(r.FormValue("months"))
		if months == "" {
			months = "1"
		}

		allowed := map[string]bool{
			"paid":      true,
			"suspend":   true,
			"unsuspend": true,
		}

		if !allowed[action] || slug == "" {
			http.Redirect(w, r, "/admin/client-billing?msg=invalid", http.StatusSeeOther)
			return
		}

		cmd := exec.Command("sudo", "/usr/local/sbin/vcr-client-billing.sh", action, slug, months)
		out, err := cmd.CombinedOutput()
		if err != nil {
			_ = out
			http.Redirect(w, r, "/admin/client-billing?msg=failed", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/admin/client-billing?msg=ok", http.StatusSeeOther)
		return
	}

	body := renderHTML(`
<div class="card owner-client-billing-page-v1">
	<div class="owner-client-billing-head-v1">
		<div>
			<h2>Billing Client</h2>
			<p class="muted">Kelola masa aktif, perpanjangan manual, penangguhan, dan aktivasi ulang layanan client.</p>
		</div>
		<a class="btn light" href="/admin/clients">Client Manager</a>
	</div>

	{{if .Msg}}
	<div class="owner-client-billing-msg-v1">{{.Msg}}</div>
	{{end}}

	<div class="tablewrap owner-client-billing-tablewrap-v1">
		<table class="owner-client-billing-table-v1">
			<thead>
				<tr>
					<th>Client</th>
					<th>Status</th>
					<th>Jatuh Tempo</th>
					<th>Tenggang</th>
					<th>Harga</th>
					<th>Terakhir Bayar</th>
					<th>Aksi</th>
				</tr>
			</thead>
			<tbody>
				{{range .Rows}}
				<tr>
					<td><b>{{.Slug}}</b><br><small>{{.PayURL}}</small></td>
					<td><span class="billing-status-mini-v1 {{.StatusClass}}">{{.Status}}</span></td>
					<td>{{.DueAt}}</td>
					<td>{{.GraceUntil}}</td>
					<td>{{.Price}}</td>
					<td>{{.LastPaidAt}}</td>
					<td>
						<div class="owner-client-billing-actions-v1">
							<form method="post">
								<input type="hidden" name="slug" value="{{.Slug}}">
								<input type="hidden" name="months" value="1">
								<button class="btn ok" name="action" value="paid" type="submit" onclick="return confirm('Tambah 1 Bulan bulan untuk {{.Slug}}?')">Tambah 1 Bulan</button>
							</form>
							<form method="post">
								<input type="hidden" name="slug" value="{{.Slug}}">
								<button class="btn danger" name="action" value="suspend" type="submit" onclick="return confirm('Suspend {{.Slug}}?')">Tangguhkan</button>
							</form>
							<form method="post">
								<input type="hidden" name="slug" value="{{.Slug}}">
								<button class="btn light" name="action" value="unsuspend" type="submit" onclick="return confirm('Unsuspend {{.Slug}}?')">Aktifkan</button>
							</form>
							<a class="btn light" href="/billing/pay?client={{.Slug}}" target="_blank" rel="noopener">Bayar / Invoice</a>
						</div>
					</td>
				</tr>
				{{else}}
				<tr><td colspan="7">Belum ada data billing client.</td></tr>
				{{end}}
			</tbody>
		</table>
	</div>
</div>
`, map[string]any{
		"Rows": a.clientBillingRowsV1(),
		"Msg":  msg,
	})

	page(w, "Billing Client", body)
}
