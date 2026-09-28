package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

/* VOUCHERGO_PUBLIC_CEK_VOUCHER_V1 */

var voucherCheckQueryRE = regexp.MustCompile(`^[A-Za-z0-9@._+-]{3,128}$`)

type voucherCheckView struct {
	Query           string
	Searched        bool
	Found           bool
	Error           string
	OrderID         string
	StatusRaw       string
	StatusLabel     string
	Tone            string
	PackageName     string
	Amount          int64
	VoucherCode     string
	ProvisionStatus string
	CreatedAt       string
	InvoiceURL      string
	Note            string
}

var voucherCheckTpl = template.Must(template.New("voucher-check").Funcs(template.FuncMap{
	"rupiah": voucherCheckRupiah,
}).Parse(`<!doctype html>
<html lang="id">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Cek Kode Voucher - VoucherGo</title>
	<style>
		:root {
			--bg:#070707; --card:#121212; --card2:#181818; --text:#f8fafc; --muted:#a3a3a3;
			--yellow:#facc15; --yellow2:#eab308; --border:rgba(250,204,21,.22);
			--ok:#22c55e; --warn:#f59e0b; --bad:#ef4444; --info:#38bdf8;
		}
		*{box-sizing:border-box}
		body{
			margin:0; min-height:100vh; font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;
			color:var(--text);
			background:
				radial-gradient(circle at 20% 10%, rgba(250,204,21,.16), transparent 28%),
				radial-gradient(circle at 85% 18%, rgba(234,179,8,.12), transparent 26%),
				linear-gradient(135deg,#030303,#111 52%,#050505);
			display:flex; align-items:center; justify-content:center; padding:22px;
		}
		.wrap{width:100%; max-width:720px}
		.hero{text-align:center; margin-bottom:18px}
		.badge{
			display:inline-flex; gap:8px; align-items:center; padding:7px 12px; border:1px solid var(--border);
			border-radius:999px; color:var(--yellow); background:rgba(250,204,21,.08); font-weight:800; font-size:13px;
		}
		h1{margin:14px 0 6px; font-size:clamp(28px,6vw,44px); line-height:1.05}
		p{color:var(--muted); margin:0; line-height:1.55}
		.card{
			background:linear-gradient(180deg,rgba(24,24,24,.96),rgba(12,12,12,.96));
			border:1px solid var(--border); border-radius:24px; padding:18px;
			box-shadow:0 24px 80px rgba(0,0,0,.45);
		}
		form{display:flex; gap:10px; margin:4px 0 16px}
		input{
			flex:1; min-width:0; border-radius:16px; border:1px solid rgba(255,255,255,.12);
			background:#070707; color:var(--text); padding:15px 14px; font-size:16px; outline:none;
		}
		input:focus{border-color:var(--yellow); box-shadow:0 0 0 4px rgba(250,204,21,.12)}
		button,a.btn{
			border:0; border-radius:16px; padding:15px 18px; cursor:pointer;
			background:linear-gradient(135deg,var(--yellow),var(--yellow2)); color:#111; font-weight:900;
			text-decoration:none; display:inline-flex; justify-content:center; align-items:center;
		}
		.result{border-radius:20px; border:1px solid rgba(255,255,255,.1); background:var(--card2); padding:16px; margin-top:12px}
		.status{display:inline-flex; padding:7px 11px; border-radius:999px; font-size:13px; font-weight:900; margin-bottom:12px}
		.status.ok{background:rgba(34,197,94,.14); color:#86efac; border:1px solid rgba(34,197,94,.28)}
		.status.warn{background:rgba(245,158,11,.14); color:#fcd34d; border:1px solid rgba(245,158,11,.28)}
		.status.bad{background:rgba(239,68,68,.14); color:#fca5a5; border:1px solid rgba(239,68,68,.28)}
		.status.info{background:rgba(56,189,248,.14); color:#7dd3fc; border:1px solid rgba(56,189,248,.28)}
		.grid{display:grid; grid-template-columns:1fr 1fr; gap:10px}
		.item{border:1px solid rgba(255,255,255,.08); border-radius:16px; padding:12px; background:rgba(0,0,0,.18)}
		.label{font-size:12px; color:var(--muted); margin-bottom:5px}
		.value{font-size:15px; font-weight:850; word-break:break-word}
		.voucher{font-size:30px; letter-spacing:4px; color:var(--yellow)}
		.note{margin-top:12px; color:#d4d4d4}
		.actions{display:flex; gap:10px; flex-wrap:wrap; margin-top:14px}
		.err{color:#fecaca; background:rgba(239,68,68,.12); border:1px solid rgba(239,68,68,.22); padding:12px; border-radius:16px}
		.empty{color:#fde68a; background:rgba(250,204,21,.1); border:1px solid rgba(250,204,21,.2); padding:12px; border-radius:16px}
		.foot{text-align:center; margin-top:14px; font-size:12px; color:#777}
		@media(max-width:560px){
			body{padding:14px; align-items:flex-start}
			form{flex-direction:column}
			button,a.btn{width:100%}
			.grid{grid-template-columns:1fr}
			.card{border-radius:20px; padding:15px}
			.voucher{font-size:26px}
		}
	
		/* VOUCHERGO_CEK_VOUCHER_MOBILE_POLISH_V1 */
		@media(max-width:720px){
			body{
				display:block;
				min-height:100svh;
				padding:14px;
				overflow-x:hidden;
				background:
					radial-gradient(circle at 20% 0%, rgba(250,204,21,.18), transparent 28%),
					radial-gradient(circle at 90% 8%, rgba(234,179,8,.12), transparent 25%),
					linear-gradient(135deg,#030303,#111 58%,#050505);
			}
			.wrap{
				max-width:100%;
				margin:0 auto;
				padding-top:18px;
			}
			.hero{
				margin-bottom:14px;
			}
			.badge{
				padding:6px 11px;
				font-size:12px;
			}
			h1{
				margin:12px 0 6px;
				font-size:34px;
				letter-spacing:-.8px;
			}
			p{
				font-size:14px;
			}
			.card{
				padding:14px;
				border-radius:22px;
				box-shadow:0 18px 55px rgba(0,0,0,.42);
			}
			form{
				display:grid;
				grid-template-columns:1fr 72px;
				gap:8px;
				margin:0 0 12px;
			}
			input{
				width:100%;
				height:52px;
				padding:0 14px;
				border-radius:15px;
				font-size:16px;
			}
			button,a.btn{
				min-height:52px;
				padding:0 14px;
				border-radius:15px;
				font-size:14px;
			}
			.result{
				padding:13px;
				border-radius:18px;
				margin-top:10px;
			}
			.status{
				margin-bottom:10px;
				font-size:12px;
			}
			.grid{
				grid-template-columns:1fr;
				gap:8px;
			}
			.item{
				padding:11px;
				border-radius:15px;
			}
			.label{
				font-size:11px;
			}
			.value{
				font-size:14px;
			}
			.voucher{
				font-size:28px;
				letter-spacing:3px;
				line-height:1.15;
			}
			.note{
				font-size:13px;
				margin-top:10px;
			}
			.actions{
				display:grid;
				grid-template-columns:1fr;
				gap:8px;
				margin-top:12px;
			}
			.empty,.err{
				padding:12px;
				border-radius:15px;
				font-size:14px;
				line-height:1.45;
			}
			.foot{
				margin-top:12px;
				padding-bottom:8px;
			}
		}

		@media(max-width:380px){
			body{padding:10px}
			.wrap{padding-top:10px}
			h1{font-size:30px}
			form{grid-template-columns:1fr 66px; gap:7px}
			button,a.btn{font-size:13px}
			.voucher{font-size:24px; letter-spacing:2px}
		}

		@media(max-height:640px) and (max-width:720px){
			.wrap{padding-top:8px}
			.hero{margin-bottom:10px}
			h1{font-size:30px; margin:8px 0 4px}
			.badge{display:none}
		}

	
		/* VOUCHERGO_CEK_VOUCHER_SIMPLE_TITLE_V1 */
		.hero-simple-title-v1{
			margin-bottom:14px;
		}
		.hero-simple-title-v1 h1{
			margin:0;
			font-size:clamp(25px,5vw,36px);
			font-weight:650;
			letter-spacing:-.45px;
			line-height:1.12;
		}
		@media(max-width:720px){
			.hero-simple-title-v1{
				margin-bottom:12px;
				padding-top:4px;
			}
			.hero-simple-title-v1 h1{
				font-size:28px;
				font-weight:600;
				letter-spacing:-.25px;
			}
		}
		@media(max-width:380px){
			.hero-simple-title-v1 h1{
				font-size:25px;
			}
		}

	</style>
</head>
<body>
	<main class="wrap">
		<section class="hero hero-simple-title-v1">
			<h1>Cek Kode Voucher</h1>
		</section>

		<section class="card">
			<form method="get" action="/cek-voucher">
				<input name="q" value="{{.Query}}" placeholder="Contoh: abcd atau VG..." autocomplete="off" autofocus>
				<button type="submit">Cek</button>
			</form>

			{{if .Error}}
				<div class="err">{{.Error}}</div>
			{{else if .Searched}}
				{{if .Found}}
					<div class="result">
						<div class="status {{.Tone}}">{{.StatusLabel}}</div>
						<div class="grid">
							<div class="item">
								<div class="label">Kode Voucher</div>
								<div class="value voucher">{{if .VoucherCode}}{{.VoucherCode}}{{else}}-{{end}}</div>
							</div>
							<div class="item">
								<div class="label">Order ID</div>
								<div class="value">{{.OrderID}}</div>
							</div>
							<div class="item">
								<div class="label">Paket</div>
								<div class="value">{{if .PackageName}}{{.PackageName}}{{else}}-{{end}}</div>
							</div>
							<div class="item">
								<div class="label">Nominal</div>
								<div class="value">{{rupiah .Amount}}</div>
							</div>
							<div class="item">
								<div class="label">Tanggal</div>
								<div class="value">{{if .CreatedAt}}{{.CreatedAt}}{{else}}-{{end}}</div>
							</div>
							<div class="item">
								<div class="label">Provision</div>
								<div class="value">{{if .ProvisionStatus}}{{.ProvisionStatus}}{{else}}-{{end}}</div>
							</div>
						</div>
						<p class="note">{{.Note}}</p>
						<div class="actions">
							<a class="btn" href="/cek-voucher" onclick="if(history.length>1){history.back();return false;}">Kembali</a> <!-- VOUCHERGO_CEK_VOUCHER_BACK_FOOTER_V1 -->
							<a class="btn" href="/x86">Beli Voucher Lagi</a> <!-- VOUCHERGO_CEK_VOUCHER_BUY_AGAIN_X86_V1 -->
						</div>
					</div>
				{{else}}
					<div class="empty">Kode voucher tidak ditemukan. Pastikan data yang dimasukkan sudah benar. <!-- VOUCHERGO_CEK_VOUCHER_NOTFOUND_TEXT_V1 --></div>
				{{end}}
			{{else}}
				<div class="empty">Silakan masukkan kode voucher yang sudah dibeli.</div>
			{{end}}
		</section>

		<div class="foot">VoucherGo 2026</div>
	</main>
</body>
</html>`))

func (a *App) publicVoucherCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	view := voucherCheckView{
		Query:    q,
		Searched: q != "",
	}

	if q != "" {
		if !voucherCheckQueryRE.MatchString(q) {
			view.Error = "Format kode tidak valid. Gunakan kode voucher atau Order ID tanpa spasi."
		} else {
			// VOUCHERGO_PUBLIC_CEK_BY_PHONE_V1
			var orderID, status, voucherUser, voucherPass, provision, pkgName, createdAt sql.NullString
			var amount sql.NullInt64

			phoneDigits := voucherCheckDigits(q)
			phone08 := phoneDigits
			phone62 := phoneDigits
			if strings.HasPrefix(phoneDigits, "0") {
				phone62 = "62" + strings.TrimPrefix(phoneDigits, "0")
			} else if strings.HasPrefix(phoneDigits, "62") {
				phone08 = "0" + strings.TrimPrefix(phoneDigits, "62")
			}

			err := a.db.QueryRow(`
				SELECT
					COALESCE(t.order_id, ''),
					COALESCE(t.status, ''),
					COALESCE(t.voucher_username, ''),
					COALESCE(t.voucher_password, ''),
					COALESCE(t.provision_status, ''),
					COALESCE(p.name, ''),
					COALESCE(t.amount, 0),
					COALESCE(DATE_FORMAT(t.created_at, '%d/%m/%Y %H:%i'), '')
				FROM transactions t
				LEFT JOIN packages p ON p.id = t.package_id
				WHERE LOWER(t.voucher_username) = LOWER(?)
				   OR LOWER(t.voucher_password) = LOWER(?)
				   OR LOWER(t.order_id) = LOWER(?)
				   OR LOWER(COALESCE(t.customer_email, '')) = LOWER(?) -- VOUCHERGO_PUBLIC_CEK_BY_EMAIL_V1
				   OR REPLACE(REPLACE(REPLACE(REPLACE(COALESCE(t.customer_phone,''), '+', ''), '-', ''), ' ', ''), '.', '') IN (?, ?, ?)
				ORDER BY t.id DESC
				LIMIT 1
			`, q, q, q, q, phoneDigits, phone08, phone62).Scan(&orderID, &status, &voucherUser, &voucherPass, &provision, &pkgName, &amount, &createdAt)

			if err == sql.ErrNoRows {
				view.Found = false
			} else if err != nil {
				log.Printf("cek voucher error q=%q: %v", q, err)
				view.Error = "Terjadi kesalahan saat mengecek voucher. Silakan coba lagi."
			} else {
				// VOUCHERGO_CEK_PENDING_CONTACT_REDIRECT_V1
				code := strings.TrimSpace(voucherUser.String)
				if code == "" {
					code = strings.TrimSpace(voucherPass.String)
				}

				statusLower := strings.ToLower(strings.TrimSpace(status.String))
				isPaid := voucherCheckIsPaidStatus(statusLower)
				isPending := voucherCheckIsPendingStatus(statusLower)
				isContactLookup := strings.Contains(q, "@") || len(phoneDigits) >= 8

				// Kalau user cek via WA/email dan order terbaru belum dibayar,
				// arahkan langsung ke invoice agar lanjut bayar, bukan tampilkan voucher.
				if isContactLookup && isPending && orderID.String != "" {
					http.Redirect(w, r, "/invoice/"+template.URLQueryEscaper(orderID.String), http.StatusSeeOther) // VOUCHERGO_CEK_VOUCHER_INVOICE_PATH_V1
					return
				}

				// Status belum lunas tidak boleh menampilkan kode voucher.
				if !isPaid {
					code = ""
				}

				view.Found = true
				view.OrderID = orderID.String
				view.StatusRaw = status.String
				view.PackageName = pkgName.String
				view.Amount = amount.Int64
				view.VoucherCode = code
				view.ProvisionStatus = provision.String
				view.CreatedAt = createdAt.String
				if view.OrderID != "" {
					view.InvoiceURL = "/invoice/" + template.URLQueryEscaper(view.OrderID)
				}
				view.StatusLabel, view.Tone, view.Note = voucherCheckStatus(status.String, provision.String, code)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := voucherCheckTpl.Execute(w, view); err != nil {
		log.Printf("render cek voucher error: %v", err)
	}
}

func voucherCheckDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func voucherCheckIsPaidStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "paid", "lunas", "success", "settlement", "capture":
		return true
	default:
		return false
	}
}

func voucherCheckIsPendingStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending", "unpaid", "menunggu", "waiting", "wait_payment":
		return true
	default:
		return false
	}
}

func voucherCheckStatus(status, provision, code string) (string, string, string) {
	s := strings.ToLower(strings.TrimSpace(status))
	p := strings.ToLower(strings.TrimSpace(provision))

	switch s {
	case "paid", "lunas", "success", "settlement", "capture":
		if code != "" && (p == "" || p == "ok" || p == "success" || p == "done") {
			return "Voucher Aktif", "ok", "Pembayaran sudah lunas dan kode voucher sudah tersedia."
		}
		if code != "" {
			return "Voucher Dibuat", "ok", "Pembayaran sudah lunas. Jika voucher belum bisa dipakai, hubungi admin hotspot."
		}
		return "Lunas, Voucher Diproses", "warn", "Pembayaran sudah lunas, tetapi kode voucher belum tercatat di sistem."
	case "pending", "unpaid", "menunggu":
		return "Menunggu Pembayaran", "warn", "Order masih menunggu pembayaran. Selesaikan pembayaran dari halaman invoice."
	case "expired", "cancelled", "canceled", "failed", "gagal":
		return "Tidak Aktif", "bad", "Order ini gagal, dibatalkan, atau sudah kedaluwarsa."
	default:
		if s == "" {
			return "Status Tidak Diketahui", "info", "Data voucher ditemukan, tetapi status pembayaran belum terbaca."
		}
		return "Status: " + status, "info", "Data voucher ditemukan."
	}
}

func voucherCheckRupiah(n int64) string {
	if n <= 0 {
		return "-"
	}
	s := fmt.Sprintf("%d", n)
	out := ""
	for i, c := range reverseString(s) {
		if i > 0 && i%3 == 0 {
			out += "."
		}
		out += string(c)
	}
	return "Rp" + reverseString(out)
}

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
