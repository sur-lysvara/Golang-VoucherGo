package main

import (
	"fmt"
	"net/http"
)

func (a *App) admin(w http.ResponseWriter, r *http.Request) {
	// ADMIN_DASHBOARD_SCOPE_V1
	// DASHBOARD_PREMIUM_STATS_V1
	role, scopedRouterID := a.currentAdminRole(r)

	var routers, packages, txPending, txPaid, txFailed int
	var revenueToday, revenueMonth, revenueTotal int64

	publicLink := "/"
	isSuper := role == "superadmin"

	if role == "router" && scopedRouterID > 0 {
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM routers WHERE id=?`, scopedRouterID).Scan(&routers)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM packages WHERE router_id=?`, scopedRouterID).Scan(&packages)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status='pending' AND router_id=?`, scopedRouterID).Scan(&txPending)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status='paid' AND router_id=?`, scopedRouterID).Scan(&txPaid)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status IN ('failed','cancelled','expired') AND router_id=?`, scopedRouterID).Scan(&txFailed)

		_ = a.db.QueryRow(`
			SELECT
				COALESCE(SUM(CASE WHEN status='paid' AND DATE(COALESCE(paid_at, created_at))=CURDATE() THEN amount ELSE 0 END),0),
				COALESCE(SUM(CASE WHEN status='paid' AND DATE_FORMAT(COALESCE(paid_at, created_at),'%Y-%m')=DATE_FORMAT(CURDATE(),'%Y-%m') THEN amount ELSE 0 END),0),
				COALESCE(SUM(CASE WHEN status='paid' THEN amount ELSE 0 END),0)
			FROM transactions
			WHERE router_id=?
		`, scopedRouterID).Scan(&revenueToday, &revenueMonth, &revenueTotal)

		var slug string
		_ = a.db.QueryRow(`SELECT slug FROM routers WHERE id=?`, scopedRouterID).Scan(&slug)
		if slug != "" {
			publicLink = "/" + slug
		}
	} else {
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM routers`).Scan(&routers)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM packages`).Scan(&packages)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status='pending'`).Scan(&txPending)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status='paid'`).Scan(&txPaid)
		_ = a.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE status IN ('failed','cancelled','expired')`).Scan(&txFailed)

		_ = a.db.QueryRow(`
			SELECT
				COALESCE(SUM(CASE WHEN status='paid' AND DATE(COALESCE(paid_at, created_at))=CURDATE() THEN amount ELSE 0 END),0),
				COALESCE(SUM(CASE WHEN status='paid' AND DATE_FORMAT(COALESCE(paid_at, created_at),'%Y-%m')=DATE_FORMAT(CURDATE(),'%Y-%m') THEN amount ELSE 0 END),0),
				COALESCE(SUM(CASE WHEN status='paid' THEN amount ELSE 0 END),0)
			FROM transactions
		`).Scan(&revenueToday, &revenueMonth, &revenueTotal)
	}

	formatRp := func(v int64) string {
		raw := fmt.Sprintf("%d", v)
		n := len(raw)
		out := ""
		for i, ch := range raw {
			if i > 0 && (n-i)%3 == 0 {
				out += "."
			}
			out += string(ch)
		}
		return "Rp" + out
	}

	txTotal := txPaid + txPending + txFailed
	paidDeg := "0"
	pendingEndDeg := "0"
	if txTotal > 0 {
		pd := float64(txPaid) / float64(txTotal) * 360
		pn := float64(txPending) / float64(txTotal) * 360
		paidDeg = fmt.Sprintf("%.2f", pd)
		pendingEndDeg = fmt.Sprintf("%.2f", pd+pn)
	}

	// DASHBOARD_RECENT_ACTIVITY_V2
	type dashboardActivity struct {
		Time   string
		Status string
		Class  string
		Title  string
		Meta   string
	}

	activities := []dashboardActivity{}

	addActivity := func(at, status, className, title, meta string) {
		if len(activities) >= 7 {
			return
		}
		if at == "" {
			at = "--:--"
		}
		if title == "" {
			title = "-"
		}
		if meta == "" {
			meta = "-"
		}
		activities = append(activities, dashboardActivity{
			Time:   at,
			Status: status,
			Class:  className,
			Title:  title,
			Meta:   meta,
		})
	}

	joinMeta := func(parts ...string) string {
		out := ""
		for _, part := range parts {
			if part == "" || part == "-" {
				continue
			}
			if out != "" {
				out += " • "
			}
			out += part
		}
		return out
	}

	activitySQL := `
		SELECT
			DATE_FORMAT(t.created_at, '%H:%i'),
			DATE_FORMAT(COALESCE(t.paid_at, t.created_at), '%H:%i'),
			COALESCE(t.status, 'pending'),
			COALESCE(p.name, ''),
			COALESCE(t.voucher_username, ''),
			LOWER(COALESCE(t.provision_status, '')),
			LOWER(COALESCE(t.wa_status, '')),
			COALESCE(rt.slug, ''),
			COALESCE(t.order_id, '')
		FROM transactions t
		LEFT JOIN packages p ON p.id=t.package_id
		LEFT JOIN routers rt ON rt.id=t.router_id
	`
	activityArgs := []any{}
	if role == "router" && scopedRouterID > 0 {
		activitySQL += ` WHERE t.router_id=? `
		activityArgs = append(activityArgs, scopedRouterID)
	}
	activitySQL += ` ORDER BY t.id DESC LIMIT 10 `

	if rows, err := a.db.Query(activitySQL, activityArgs...); err == nil {
		defer rows.Close()
		for rows.Next() {
			var createdAt, paidAt, status, pkgName, voucher, provision, wa, routerSlug, orderID string
			if err := rows.Scan(&createdAt, &paidAt, &status, &pkgName, &voucher, &provision, &wa, &routerSlug, &orderID); err != nil {
				continue
			}

			pkgLabel := pkgName
			if pkgLabel == "" {
				pkgLabel = "Paket"
			}

			voucherLabel := ""
			if voucher != "" {
				voucherLabel = "kode " + voucher
			}

			orderLabel := orderID
			if len(orderLabel) > 14 {
				orderLabel = orderLabel[:14] + "…"
			}

			switch status {
			case "paid":
				addActivity(paidAt, "Paid", "paid", "Pembayaran masuk", joinMeta(pkgLabel, voucherLabel, routerSlug))
			case "pending":
				addActivity(createdAt, "Pending", "pending", "Menunggu pembayaran", joinMeta(pkgLabel, orderLabel, routerSlug))
			case "failed":
				addActivity(createdAt, "Gagal", "failed", "Pembayaran gagal", joinMeta(pkgLabel, orderLabel, routerSlug))
			case "expired":
				addActivity(createdAt, "Expired", "failed", "Pembayaran expired", joinMeta(pkgLabel, orderLabel, routerSlug))
			case "cancelled":
				addActivity(createdAt, "Batal", "failed", "Pembayaran dibatalkan", joinMeta(pkgLabel, orderLabel, routerSlug))
			}

			if provision == "success" && voucher != "" {
				addActivity(paidAt, "MT", "mikrotik", "Voucher dibuat di MikroTik", joinMeta(pkgLabel, voucherLabel, routerSlug))
			} else if provision == "failed" {
				addActivity(paidAt, "MT", "failed", "Provision MikroTik gagal", joinMeta(pkgLabel, orderLabel, routerSlug))
			}

			if wa == "sent" || wa == "success" || wa == "ok" || wa == "delivered" {
				addActivity(paidAt, "WA", "wa", "WA voucher terkirim", joinMeta(voucherLabel, routerSlug))
			} else if wa == "failed" || wa == "error" {
				addActivity(paidAt, "WA", "failed", "WA gagal terkirim", joinMeta(voucherLabel, routerSlug))
			}

			if len(activities) >= 7 {
				break
			}
		}
	}

	body := renderHTML(`
<style>
/* DASHBOARD_TOP_CARDS_POLISH_V1 */
.dashboard-top-cards-v1{
	display:grid;
	grid-template-columns:repeat(4,minmax(0,1fr));
	gap:12px;
	margin-bottom:16px;
}
.dashboard-top-card-v1{
	position:relative;
	display:block;
	text-decoration:none;
	padding:18px 18px;
	border-radius:18px;
	background:linear-gradient(135deg,rgba(255,255,255,.07),rgba(255,255,255,.035));
	border:1px solid rgba(255,255,255,.09);
	box-shadow:0 14px 35px rgba(0,0,0,.18);
	overflow:hidden;
	min-height:92px;
}
.dashboard-top-card-v1:before{
	content:"";
	position:absolute;
	inset:auto -35px -45px auto;
	width:110px;
	height:110px;
	border-radius:999px;
	background:var(--glow);
	opacity:.16;
}
.dashboard-top-card-v1:hover{
	transform:translateY(-1px);
	border-color:rgba(255,255,255,.18);
}
.dashboard-top-card-v1 .label{
	position:relative;
	z-index:1;
	font-weight:800;
	color:rgba(255,255,255,.88);
	font-size:14px;
}
.dashboard-top-card-v1 .num{
	position:relative;
	z-index:1;
	font-size:31px;
	font-weight:950;
	line-height:1;
	margin-top:8px;
	color:#fff;
}
.dashboard-top-card-v1 .hint{
	position:absolute;
	right:14px;
	top:14px;
	width:27px;
	height:27px;
	border-radius:999px;
	display:grid;
	place-items:center;
	background:rgba(255,255,255,.08);
	color:rgba(255,255,255,.75);
	font-weight:900;
}
@media(max-width:900px){
	.dashboard-top-cards-v1{grid-template-columns:repeat(2,minmax(0,1fr));}
}
@media(max-width:420px){
	.dashboard-top-cards-v1{gap:8px;margin-bottom:10px;}
	.dashboard-top-card-v1{padding:13px 13px;min-height:78px;border-radius:15px;}
	.dashboard-top-card-v1 .num{font-size:25px;}
	.dashboard-top-card-v1 .label{font-size:12px;}
	.dashboard-top-card-v1 .hint{width:23px;height:23px;right:10px;top:10px;font-size:12px;}
}

/* LIGHT_DASHBOARD_TOP_TEXT_FIX_V1 */
html[data-admin-theme="light"] .dashboard-top-card-v1,
html[data-theme="light"] .dashboard-top-card-v1{
	background:linear-gradient(135deg,#ffffff,#f8fafc)!important;
	border-color:#dbe4f0!important;
	box-shadow:0 12px 28px rgba(15,23,42,.08)!important;
}

html[data-admin-theme="light"] .dashboard-top-card-v1 .label,
html[data-theme="light"] .dashboard-top-card-v1 .label{
	color:#334155!important;
}

html[data-admin-theme="light"] .dashboard-top-card-v1 .num,
html[data-theme="light"] .dashboard-top-card-v1 .num{
	color:#0f172a!important;
	text-shadow:none!important;
}

html[data-admin-theme="light"] .dashboard-top-card-v1 .hint,
html[data-theme="light"] .dashboard-top-card-v1 .hint{
	background:#e2e8f0!important;
	color:#334155!important;
}

html[data-admin-theme="light"] .dashboard-top-card-v1:before,
html[data-theme="light"] .dashboard-top-card-v1:before{
	opacity:.24!important;
}
/* /LIGHT_DASHBOARD_TOP_TEXT_FIX_V1 */

</style>

<div class="dashboard-top-cards-v1">
	<a class="dashboard-top-card-v1" style="--glow:#38bdf8" href="/admin/routers">
		<div class="label">Router</div>
		<div class="num">{{.Routers}}</div>
		<div class="hint">›</div>
	</a>
	<!-- DASHBOARD_FAILED_CARD_V1 -->
    <a class="dashboard-top-card-v1" style="--glow:#ef4444" href="/admin/transactions?status=failed">
            <div class="label">Gagal/Batal</div>
            <div class="num">{{.Failed}}</div>
            <div class="hint">›</div>
    </a>
	<a class="dashboard-top-card-v1" style="--glow:#ff8a0a" href="/admin/transactions?status=pending">
		<div class="label">Pending</div>
		<div class="num">{{.Pending}}</div>
		<div class="hint">›</div>
	</a>
	<a class="dashboard-top-card-v1" style="--glow:#22c55e" href="/admin/transactions?status=paid">
		<div class="label">Paid</div>
		<div class="num">{{.Paid}}</div>
		<div class="hint">›</div>
	</a>
</div>

<style>
/* DASHBOARD_PREMIUM_STATS_V1 */
body:has(.dashboard-premium-stats-v1) .dashboard-premium-stats-v1{
	display:grid;
	grid-template-columns:minmax(280px,1fr) minmax(300px,.9fr);
	gap:18px;
	align-items:stretch;
	margin:16px 0;
}
body:has(.dashboard-premium-stats-v1) .dashboard-donut-wrap-v1{
	display:flex;
	align-items:center;
	justify-content:center;
	min-height:210px;
}
body:has(.dashboard-premium-stats-v1) .dashboard-donut-v1{
	width:210px;
	height:210px;
	border-radius:999px;
	background:conic-gradient(#22c55e 0deg {{.PaidDeg}}deg, #ff8a0a {{.PaidDeg}}deg {{.PendingEndDeg}}deg, #ef4444 {{.PendingEndDeg}}deg 360deg);
	position:relative;
	box-shadow:0 18px 45px rgba(0,0,0,.25);
}
body:has(.dashboard-premium-stats-v1) .dashboard-donut-v1:after{
	content:"";
	position:absolute;
	inset:58px;
	border-radius:999px;
	background:var(--card, #20142d);
	box-shadow:inset 0 0 0 1px rgba(255,255,255,.07);
}
body:has(.dashboard-premium-stats-v1) .dashboard-legend-v1{
	display:flex;
	justify-content:center;
	gap:14px;
	flex-wrap:wrap;
	margin:12px 0 0;
	font-size:13px;
	color:var(--muted, #94a3b8);
}
body:has(.dashboard-premium-stats-v1) .dashboard-legend-v1 span,
body:has(.dashboard-premium-stats-v1) .dash-income-item-v1{
	display:flex;
	align-items:center;
	gap:8px;
}
body:has(.dashboard-premium-stats-v1) .dash-dot-v1{
	width:13px;
	height:13px;
	border-radius:999px;
	display:inline-block;
	flex:0 0 auto;
}
body:has(.dashboard-premium-stats-v1) .dashboard-revenue-list-v1{
	display:grid;
	gap:12px;
	margin-top:14px;
}
body:has(.dashboard-premium-stats-v1) .dash-income-item-v1{
	padding:15px 16px;
	border-radius:18px;
	background:rgba(255,255,255,.055);
	border:1px solid rgba(255,255,255,.075);
}
body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .num{
	font-size:25px;
	font-weight:900;
	line-height:1.1;
}
body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .label{
	color:var(--muted, #94a3b8);
	margin-top:4px;
}
@media(max-width:720px){
	body:has(.dashboard-premium-stats-v1) .dashboard-premium-stats-v1{
		grid-template-columns:1fr;
		gap:10px;
		margin:10px 0;
	}
	body:has(.dashboard-premium-stats-v1) .dashboard-donut-wrap-v1{
		min-height:175px;
	}
	body:has(.dashboard-premium-stats-v1) .dashboard-donut-v1{
		width:165px;
		height:165px;
	}
	body:has(.dashboard-premium-stats-v1) .dashboard-donut-v1:after{
		inset:46px;
	}
	body:has(.dashboard-premium-stats-v1) .dash-income-item-v1{
		padding:12px 14px;
	}
	body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .num{
		font-size:21px;
	}
}
</style>

<div class="dashboard-premium-stats-v1">
	<div class="card dashboard-chart-card-v1">
		<h3>Statistik Transaksi</h3>
		<div class="dashboard-donut-wrap-v1">
			<div class="dashboard-donut-v1" aria-label="Statistik transaksi"></div>
		</div>
		<div class="dashboard-legend-v1">
			<span><i class="dash-dot-v1" style="background:#22c55e"></i>Paid {{.Paid}}</span>
			<span><i class="dash-dot-v1" style="background:#ff8a0a"></i>Pending {{.Pending}}</span>
			<span><i class="dash-dot-v1" style="background:#ef4444"></i>Gagal/Batal {{.Failed}}</span>
		</div>
	</div>

	<div class="card dashboard-revenue-card-v1">
		<h3>Pendapatan</h3>
		<div class="dashboard-revenue-list-v1">
			<div class="dash-income-item-v1">
				<i class="dash-dot-v1" style="background:#22c55e"></i>
				<div><div class="num">{{.RevenueToday}}</div><div class="label">Hari Ini</div></div>
			</div>
			<div class="dash-income-item-v1">
				<i class="dash-dot-v1" style="background:#38bdf8"></i>
				<div><div class="num">{{.RevenueMonth}}</div><div class="label">Bulan Ini</div></div>
			</div>
			<div class="dash-income-item-v1">
				<i class="dash-dot-v1" style="background:#a78bfa"></i>
				<div><div class="num">{{.RevenueTotal}}</div><div class="label">Total Lunas</div></div>
			</div>
		</div>
	</div>
</div>

<!-- DASHBOARD_RECENT_ACTIVITY_CARD_V1 -->
<style>
.dashboard-activity-card-v1{
	margin-top:14px;
}
.dashboard-activity-head-v1{
	display:flex;
	align-items:center;
	justify-content:space-between;
	gap:10px;
	margin-bottom:10px;
}
.dashboard-activity-head-v1 h3{
	margin:0!important;
}
.dashboard-activity-head-v1 a{
	font-size:13px;
	font-weight:800;
	text-decoration:none;
	color:#93c5fd;
}
.dashboard-activity-list-v1{
	display:grid;
	gap:8px;
}
.dashboard-activity-item-v1{
	display:grid;
	grid-template-columns:48px 74px minmax(0,1fr);
	gap:9px;
	align-items:center;
	padding:10px 11px;
	border-radius:15px;
	background:rgba(255,255,255,.045);
	border:1px solid rgba(255,255,255,.075);
}
.dashboard-activity-time-v1{
	color:var(--muted,#94a3b8);
	font-weight:800;
	font-size:12px;
}
.dashboard-activity-badge-v1{
	display:inline-flex;
	align-items:center;
	justify-content:center;
	min-height:24px;
	padding:0 8px;
	border-radius:999px;
	font-size:11px;
	font-weight:900;
	background:rgba(148,163,184,.14);
	color:#e5e7eb;
}
.dashboard-activity-badge-v1.paid{
	background:rgba(34,197,94,.16);
	color:#86efac;
}
.dashboard-activity-badge-v1.pending{
	background:rgba(255,138,10,.15);
	color:#fdba74;
}
.dashboard-activity-badge-v1.failed{
	background:rgba(239,68,68,.15);
	color:#fca5a5;
}
/* DASHBOARD_ACTIVITY_BADGE_V2 */
.dashboard-activity-badge-v1.wa{
	background:rgba(14,165,233,.16);
	color:#7dd3fc;
}
.dashboard-activity-badge-v1.mikrotik{
	background:rgba(167,139,250,.16);
	color:#c4b5fd;
}
.dashboard-activity-main-v1{
	min-width:0;
}
.dashboard-activity-title-v1{
	font-weight:900;
	color:#f8fafc;
	white-space:nowrap;
	overflow:hidden;
	text-overflow:ellipsis;
}
.dashboard-activity-meta-v1{
	margin-top:2px;
	font-size:12px;
	color:var(--muted,#94a3b8);
	white-space:nowrap;
	overflow:hidden;
	text-overflow:ellipsis;
}
.dashboard-activity-empty-v1{
	padding:12px;
	border-radius:15px;
	background:rgba(255,255,255,.045);
	color:var(--muted,#94a3b8);
	font-weight:700;
}
html[data-admin-theme="light"] .dashboard-activity-item-v1,
html[data-theme="light"] .dashboard-activity-item-v1,
html[data-admin-theme="light"] .dashboard-activity-empty-v1,
html[data-theme="light"] .dashboard-activity-empty-v1{
	background:#f8fafc!important;
	border-color:#e2e8f0!important;
}
html[data-admin-theme="light"] .dashboard-activity-title-v1,
html[data-theme="light"] .dashboard-activity-title-v1{
	color:#0f172a!important;
}

@media(max-width:768px){
	.dashboard-activity-card-v1{
		margin-top:10px;
		padding:14px 13px!important;
		border-radius:20px!important;
	}
	.dashboard-activity-head-v1{
		margin-bottom:9px;
	}
	.dashboard-activity-head-v1 h3{
		font-size:17px!important;
	}
	.dashboard-activity-list-v1{
		gap:7px;
	}
	.dashboard-activity-item-v1{
		grid-template-columns:42px 64px minmax(0,1fr);
		gap:7px;
		padding:9px 9px;
		border-radius:14px;
	}
	.dashboard-activity-time-v1{
		font-size:11px;
	}
	.dashboard-activity-badge-v1{
		min-height:22px;
		padding:0 7px;
		font-size:10px;
	}
	.dashboard-activity-title-v1{
		font-size:13px;
	}
	.dashboard-activity-meta-v1{
		font-size:11px;
	}
}
</style>

<div class="card dashboard-activity-card-v1">
	<div class="dashboard-activity-head-v1">
		<h3>Aktivitas Terakhir</h3>
		<a href="/admin/transactions">Lihat semua</a>
	</div>
	{{if .Activities}}
	<div class="dashboard-activity-list-v1">
		{{range .Activities}}
		<a class="dashboard-activity-item-v1" href="/admin/transactions">
			<div class="dashboard-activity-time-v1">{{.Time}}</div>
			<div><span class="dashboard-activity-badge-v1 {{.Class}}">{{.Status}}</span></div>
			<div class="dashboard-activity-main-v1">
				<div class="dashboard-activity-title-v1">{{.Title}}</div>
				<div class="dashboard-activity-meta-v1">{{.Meta}}</div>
			</div>
		</a>
		{{end}}
	</div>
	{{else}}
	<div class="dashboard-activity-empty-v1">Belum ada aktivitas transaksi.</div>
	{{end}}
</div>`, map[string]any{
		"Routers":       routers,
		"Packages":      packages,
		"Pending":       txPending,
		"Paid":          txPaid,
		"Failed":        txFailed,
		"RevenueToday":  formatRp(revenueToday),
		"RevenueMonth":  formatRp(revenueMonth),
		"RevenueTotal":  formatRp(revenueTotal),
		"PaidDeg":       paidDeg,
		"PendingEndDeg": pendingEndDeg,
		"PublicLink":    publicLink,
		"IsSuper":       isSuper,
		"Activities":    activities,
	})

	page(w, "VoucherGo Admin", body)
}
