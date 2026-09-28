package main

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

func moveInvoiceVoucherBox(h string) string {
	boxClass := strings.Index(h, `class="voucher-box"`)
	if boxClass < 0 {
		return h
	}

	boxStart := strings.LastIndex(h[:boxClass], "<div")
	if boxStart < 0 {
		return h
	}

	boxEnd := findMatchingDivEndHTML(h, boxStart)
	if boxEnd < 0 {
		return h
	}

	gridIndex := strings.Index(h, `<div class="grid">`)
	if gridIndex < 0 || boxStart < gridIndex {
		return h
	}

	voucherBox := strings.TrimSpace(h[boxStart:boxEnd])
	h = h[:boxStart] + h[boxEnd:]

	gridIndex = strings.Index(h, `<div class="grid">`)
	if gridIndex < 0 {
		return h
	}

	return h[:gridIndex] + "\n" + voucherBox + "\n" + h[gridIndex:]
}

func (a *App) invoice(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/invoice/")

	var x Tx
	var routerSlug string

	err := a.db.QueryRow(`SELECT
			t.id,
			t.order_id,
			r.name,
			r.slug,
			p.name,
			t.customer_name,
			t.customer_phone,
			t.amount,
			t.status,
			COALESCE(t.payment_method,''),
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,''),
			t.provision_status,
			COALESCE(t.provision_error,''),
			DATE_FORMAT(t.created_at,'%Y-%m-%d %H:%i'),
			DATE_FORMAT(t.paid_at,'%Y-%m-%d %H:%i'),
			COALESCE(t.qris_string,''),
			COALESCE(t.qris_total_payment,0),
			COALESCE(t.qris_expired_at,'')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		JOIN packages p ON p.id=t.package_id
		WHERE t.order_id=?`, orderID).
		Scan(
			&x.ID,
			&x.OrderID,
			&x.RouterName,
			&routerSlug,
			&x.PackageName,
			&x.CustomerName,
			&x.CustomerPhone,
			&x.Amount,
			&x.Status,
			&x.PaymentMethod,
			&x.VoucherUsername,
			&x.VoucherPassword,
			&x.ProvisionStatus,
			&x.ProvisionError,
			&x.CreatedAt,
			&x.PaidAt,
			&x.QRISString,
			&x.QRISTotalPayment,
			&x.QRISExpiredAt,
		)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	if !strings.EqualFold(strings.TrimSpace(x.Status), "pending") {
		clearPendingInvoiceCookieIfMatches(w, r, x.OrderID)
	}

	if strings.EqualFold(strings.TrimSpace(x.Status), "pending") && r.URL.Query().Get("cancel") == "1" {
		_, _ = a.db.Exec(`UPDATE transactions SET status='cancelled' WHERE order_id=? AND status='pending'`, x.OrderID)
		clearPendingInvoiceCookieIfMatches(w, r, x.OrderID)

		target := "/"
		if strings.TrimSpace(routerSlug) != "" {
			target = "/" + url.PathEscape(routerSlug) + "?new_order=1"
		}

		http.Redirect(w, r, target, http.StatusSeeOther)
		return
	}

	payURL := ""
	loginURL := "/go-login"
	if routerSlug != "" {
		loginURL = "/go-login?router=" + url.QueryEscape(routerSlug)
	}

	slug := a.setting("pakasir_slug")

	if slug != "" && strings.TrimSpace(x.QRISString) == "" {
		q := url.Values{}
		q.Set("order_id", x.OrderID)

		if a.setting("pakasir_qris_only") == "1" {
			q.Set("qris_only", "1")
		}

		payURL = "/pay-start/" + url.PathEscape(x.OrderID)
	}

	choosePackageURL := "/"
	if strings.TrimSpace(routerSlug) != "" {
		choosePackageURL = "/" + url.PathEscape(routerSlug) + "?new_order=1"
	}
	cancelInvoiceURL := "/invoice/" + url.PathEscape(x.OrderID) + "?cancel=1"

	body := renderHTML(`
<link href="https://fonts.googleapis.com/css2?family=Ubuntu:wght@400;500;600;700&display=swap" rel="stylesheet">
<div class="invoice-wrap" style="font-family: 'Ubuntu', sans-serif;">
	<style>
		.invoice-main-card {
			padding: 24px; 
			border-radius: 24px; 
			box-shadow: 0 10px 30px rgba(15,23,42,0.06);
			background: var(--bg-card, #fff);
		}
		.invoice-info-grid {
			display: grid !important;
			grid-template-columns: repeat(3, minmax(0, 1fr)) !important;
		}
		.qris-download-note {
			color: #1e40af;
			background: #eff6ff;
			border: 1px solid #bfdbfe;
		}
		html[data-theme="dark"] .qris-download-note {
			color: #93c5fd !important;
			background: rgba(59, 130, 246, 0.15) !important;
			border-color: rgba(59, 130, 246, 0.3) !important;
		}
		@media(max-width: 640px) {
			.invoice-wrap .invoice-main-card {
				padding: 0 !important;
				border: none !important;
				box-shadow: none !important;
				background: transparent !important;
			}
			.invoice-wrap .invoice-info-grid {
				grid-template-columns: repeat(2, minmax(0, 1fr)) !important;
				gap: 8px !important;
			}
			.invoice-wrap .invoice-info-grid .provision-card {
				display: none !important;
			}
			.invoice-wrap .invoice-info-grid .notice {
				padding: 14px 12px !important;
			}
			.invoice-wrap .invoice-info-grid .price {
				font-size: 18px !important;
			}
			.invoice-wrap .invoice-info-grid .muted {
				font-size: 10.5px !important;
			}
			.invoice-wrap .invoice-info-grid b {
				font-size: 13px !important;
			}
			.invoice-pending-actions-v1 {
				display: grid !important;
				grid-template-columns: 1fr 1fr !important;
				gap: 8px !important;
			}
			.invoice-pending-actions-v1 .btn {
				width: 100% !important;
				padding: 10px 4px !important;
				font-size: 12px !important;
				margin: 0 !important;
			}
		}
	/* VG_VOUCHER_CARD_TEXT_CLEAN_V2 */
</style>
	<div class="card invoice-main-card">
		<div class="actions" style="justify-content:space-between; align-items:flex-start;">
			<div style="display:flex; align-items:flex-start; gap:12px;">
				<div style="width:40px; height:40px; background:rgba(37,99,235,0.1); border-radius:12px; display:flex; align-items:center; justify-content:center; flex-shrink:0;">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
				</div>
				<div>
					<h2 style="margin:0; font-size:14.5px; font-weight:600; color:var(--ink); line-height:1.2;">{{.Tx.OrderID}}</h2>
					<p class="muted" style="margin:4px 0 0; font-size:12px; font-weight:400; line-height:1.2;">{{.Tx.RouterName}} • {{.Tx.PackageName}}</p>
					{{/* VG_INVOICE_DATE_PENDING_ONLY_V1 */}}
					{{if and .Tx.QRISExpiredAt (eq .Tx.Status "pending")}}
					<div class="muted" style="font-size:11px; font-weight:600; display:flex; align-items:center; gap:4px; margin-top:6px;">
						<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
						Batas bayar: <span class="qris-expire-date" data-date="{{.Tx.QRISExpiredAt}}">{{.Tx.QRISExpiredAt}}</span>
					</div>
					{{end}}
				</div>
			</div>
{{/* VG_INVOICE_BADGE_FORCE_GREEN_V1 */}}
		{{/* VG_INVOICE_LABEL_PEMBELI_USER_V1 */}}
			{{/* VG_INVOICE_BADGE_PENDING_YELLOW_V1 */}}
			<span class="invoice-status-badge {{.Tx.Status}}" style="padding:6px 14px;border-radius:999px;font-weight:700;font-size:12.5px;white-space:nowrap;{{if or (eq .Tx.Status "paid") (eq .Tx.Status "success") (eq .Tx.Status "lunas")}}background:rgba(22,163,74,.12)!important;color:#16a34a!important;border:1px solid rgba(22,163,74,.24)!important;box-shadow:none!important;{{else if eq .Tx.Status "pending"}}background:rgba(234,179,8,.14)!important;color:#facc15!important;border:1px solid rgba(234,179,8,.30)!important;box-shadow:none!important;{{else}}background:rgba(220,38,38,.10)!important;color:#dc2626!important;border:1px solid rgba(220,38,38,.18)!important;{{end}}">{{if eq .Tx.Status "pending"}}Belum dibayar{{else if or (eq .Tx.Status "paid") (eq .Tx.Status "success") (eq .Tx.Status "lunas")}}Lunas{{else}}{{.Tx.Status}}{{end}}</span>
		</div>

		<hr style="border:0;border-top:2px dashed #e2e8f0;margin:24px 0">

		<div class="grid invoice-info-grid" style="gap:10px;">
			<div class="notice" style="background: linear-gradient(135deg, rgba(37,99,235,0.06), rgba(37,99,235,0.12)); border: 1px solid rgba(37,99,235,0.1); box-shadow: 0 4px 12px rgba(37,99,235,0.05); border-radius: 18px; padding: 18px;">
				<div class="muted" style="display:flex; align-items:center; gap:6px; font-size:12px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px; color:#3b82f6;">
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="6" width="20" height="12" rx="2"></rect><circle cx="12" cy="12" r="2"></circle><path d="M6 12h.01M18 12h.01"></path></svg> Total Bayar
				</div>
				<div class="price" style="color: #1e40af; font-size: 26px; font-weight: 700; margin-top: 6px; text-shadow: 0 2px 4px rgba(0,0,0,0.02);">{{if gt .Tx.QRISTotalPayment 0}}{{rp .Tx.QRISTotalPayment}}{{else}}{{rp .Tx.Amount}}{{end}}</div>
			</div>
			<div class="notice" style="background: rgba(248,250,252,0.6); border: 1px solid #e2e8f0; border-radius: 18px; padding: 16px;">
				<div class="muted" style="display:flex; align-items:center; gap:6px; font-size:12px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px;">
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg> User
				</div>
				<b style="font-size:16px; display:block; margin-top:6px; margin-bottom:0; line-height:1.2; color:var(--ink);">{{.Tx.CustomerName}}</b>
				<span class="muted" style="font-size:13px; display:block; margin-top:2px;">{{.Tx.CustomerPhone}}</span>
			</div>
			<div class="notice provision-card" style="background: rgba(248,250,252,0.6); border: 1px solid #e2e8f0; border-radius: 18px; padding: 16px;">
				<div class="muted" style="display:flex; align-items:center; gap:6px; font-size:12px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px;">
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg> Provision
				</div>
				<span class="badge {{.Tx.ProvisionStatus}}" style="margin-top:6px; display:inline-block; border-radius:8px; padding:4px 10px; font-weight:700;">{{.Tx.ProvisionStatus}}</span><br>
				<small class="muted" style="font-size:12px; display:block; margin-top:4px;">{{.Tx.ProvisionError}}</small>
			</div>
		</div>

		{{/* VG_VOUCHER_BUTTONS_SALIN_BELI_LAGI_V1 */}}
		{{/* VG_SALIN_KODE_ORANGE_TRANSPARENT_V1 */}}
		{{/* VG_COPY_NOTIF_KODE_DISALIN_NORMAL_V1 */}}
		{{if eq .Tx.Status "paid"}}
			<div class="voucher-box" style="margin-top: 28px; background: linear-gradient(135deg, #ffffff, #f8fafc); border: 2px dashed #cbd5e1; border-radius: 20px; padding: 32px 24px; text-align: center; position: relative; box-shadow: 0 12px 30px rgba(15,23,42,0.05); overflow: hidden;">
				<!-- Fake cutouts for physical ticket feel -->
				<div style="position:absolute; left:-12px; top:50%; width:24px; height:24px; background:var(--bg); border-radius:50%; transform:translateY(-50%); box-shadow: inset -2px 0 4px rgba(0,0,0,0.06);"></div>
				<div style="position:absolute; right:-12px; top:50%; width:24px; height:24px; background:var(--bg); border-radius:50%; transform:translateY(-50%); box-shadow: inset 2px 0 4px rgba(0,0,0,0.06);"></div>

				<div class="muted" style="font-weight: 700; letter-spacing: 1.5px; text-transform: uppercase; font-size: 11px; margin-bottom: 8px; color:#64748b;">Kode Voucher</div>
				<div class="voucher-code" style="font-size: 36px; font-weight: 900; color: #0f172a; letter-spacing: 3px; padding: 14px 28px; background: #fff; border-radius: 16px; border: 1px solid #e2e8f0; display: inline-block; box-shadow: 0 4px 12px rgba(15,23,42,0.04); margin-bottom: 12px; font-family: monospace;">{{.Tx.VoucherUsername}}</div>
				<p style="margin:0; font-size:13px; color:#64748b; font-weight:500;">Password sama dengan kode voucher.</p>

				<div id="voucher-copy-text" class="copy-text" style="display:none;">Username: {{.Tx.VoucherUsername}}
Password: {{.Tx.VoucherPassword}}</div>

				<div class="actions" style="justify-content:center; margin-top: 24px; gap: 12px;">
					<button class="btn pri" type="button" onclick="copyVoucher()" style="flex:1; border-radius:14px; box-shadow:0 6px 16px rgba(249,115,22,0.10); transition:all 0.2s; font-weight:800; font-size:14px; border:1px solid rgba(249,115,22,.34)!important; background:rgba(249,115,22,.13)!important; color:#fb923c!important;">
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="margin-right:6px;"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg> Salin Kode
					</button>
					<a class="btn light" href="{{.ChoosePackageURL}}" style="flex:1; border-radius:14px; box-shadow:0 6px 16px rgba(15,23,42,0.06); transition: all 0.2s; font-weight:800; font-size:14px; border:1px solid #e2e8f0; background:#fff; color:#0f172a; text-decoration:none; display:inline-flex; align-items:center; justify-content:center;">
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="margin-right:6px;"><path d="M3 12a9 9 0 0 1 15.5-6.3"></path><path d="M18 3v5h-5"></path><path d="M21 12a9 9 0 0 1-15.5 6.3"></path><path d="M6 21v-5h5"></path></svg> Beli Lagi
					</a>
				</div>
				<style>
					.actions .btn.pri:hover { transform: translateY(-2px); box-shadow: 0 8px 20px rgba(37,99,235,0.4) !important; }
					.actions .btn.light:hover { transform: translateY(-2px); box-shadow: 0 8px 20px rgba(15,23,42,0.1) !important; }
				</style>

				<p id="copy-status" style="margin:12px 0 0; font-size:13px; color:#16a34a; font-weight:800;"></p>
			</div>

			<div class="notice" style="margin-top:20px; background: rgba(37,99,235,0.06); color:#1e40af; border:1px solid rgba(37,99,235,0.12); border-radius:16px; display:flex; gap:14px; align-items:flex-start; padding:16px;">
				<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink:0; margin-top:1px;"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
				<!-- VG_VOUCHER_NOTE_SEMESTA_V1 -->
				<!-- VG_VOUCHER_NOTE_NOT_BOLD_V1 -->
				<div style="font-size:13.5px; line-height:1.55; font-weight:500;">Terima kasih sudah order 🙏<br>Voucher sudah berhasil dibuat.<br>Silakan online sebelum semesta berubah pikiran. 🗿</div>
			</div>
		{{else}}
			<div class="notice" style="margin-top:12px; border-radius:20px; padding:24px; background:linear-gradient(180deg, #ffffff, #f8fafc); border:1px solid #e2e8f0; box-shadow:0 10px 30px rgba(15,23,42,0.05);">
				{{if .Tx.QRISString}}
					<div class="qris-pay-box" style="text-align:center; padding:12px 0 24px;">
						<div style="margin-bottom:16px;">
							<h3 style="margin:0 0 4px; font-size:20px; font-weight:700; color:var(--ink);">Selesaikan Pembayaran</h3>
						</div>
						
						<a href="{{.QRImageURL}}" download="qris-{{.Tx.OrderID}}.png" target="_blank" style="display:inline-block; border:1px solid #e2e8f0; border-radius:12px; padding:12px; background:#fff; box-shadow:0 8px 24px rgba(15,23,42,0.04); transition:transform 0.2s;">
							<img id="qris-pay-img" src="{{.QRImageURL}}" alt="QRIS Pembayaran" style="width:260px; max-width:100%; border-radius:12px; display:block;">
						</a>
						
						<div style="margin-top:8px; margin-bottom:8px;">
							<div style="font-size:13.5px; color:#475569; line-height:1.5;">
								Scan kode QRIS di atas melalui aplikasi m-Banking atau e-Wallet kesayangan Anda.
							</div>
						</div>
						
						<div class="qris-download-note" style="margin-top:8px; font-size:12.5px; color:#1e40af; background:#eff6ff; border:1px solid #bfdbfe; padding:8px 16px; border-radius:999px; display:inline-block; font-weight:600;">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="vertical-align:middle; margin-right:6px; margin-top:-2px;"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
							Ketuk gambar QRIS untuk menyimpan
						</div>
						
						{{if eq .Tx.Status "pending"}}
						<div class="actions invoice-pending-actions-v1" style="margin-top:12px; gap:12px;">
							<a class="btn light" href="{{.ChoosePackageURL}}" style="flex:1; border-radius:14px; font-weight:800; background:#fff; border:1px solid #e2e8f0; color:#0f172a; padding:12px;">Pilih Paket Lagi</a>
							<a class="btn" href="{{.CancelInvoiceURL}}" style="flex:1; border-radius:14px; background:#fee2e2; color:#dc2626; border:none; font-weight:800; padding:12px;">Batalkan</a>
						</div>
						{{end}}

						<div class="qris-bypass-box-v11" style="margin-top:24px; text-align:left; background:#fff; border:1px solid #e2e8f0; border-radius:18px; padding:18px; box-shadow:0 4px 16px rgba(15,23,42,0.03);">
							<button type="button" id="qrisBypassBtnV11" class="qris-bypass-btn-v11" style="width:100%; border-radius:14px; background:#f8fafc; color:#0f172a; border:2px solid #e2e8f0; font-weight:800; padding:12px; cursor:pointer; transition:all 0.2s; font-size:14px;">
								Aktifkan Bypass 60 Detik
							</button>
							<div id="qrisBypassStatusV11" class="qris-bypass-status-v11" style="font-size:12px; color:#94a3b8; margin-top:10px; text-align:center; font-weight:500;">
								Gunakan jika browser belum konek otomatis.
							</div>
						</div>
						
						<div style="margin-top:24px; padding-top:20px; border-top:2px dashed #e2e8f0;">
							<p class="muted" style="margin:0; font-size:13px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px;">Total Tagihan</p>
							<p style="margin:6px 0 0; font-size:26px; font-weight:700; color:#2563eb;">{{rp .Tx.QRISTotalPayment}}</p>
						</div>
					</div>
				{{else}}
					<div style="text-align:center; padding:24px 0;">
						<div style="width:72px; height:72px; background:#eff6ff; border-radius:50%; display:flex; align-items:center; justify-content:center; margin:0 auto 20px;">
							<svg xmlns="http://www.w3.org/2000/svg" width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
						</div>
						<h3 style="margin:0 0 10px; font-size:18px; font-weight:700; color:var(--ink);">Menunggu Pembayaran</h3>
						<p style="margin:0; font-size:14px; color:#64748b; line-height:1.6; font-weight:500;">Status tagihan Anda saat ini <b>{{.Tx.Status}}</b>.<br>Sistem akan menampilkan kode voucher otomatis setelah pelunasan berhasil.</p>
					</div>
				{{end}}
			</div>

			{{if .PayURL}}
				<div style="margin-top:24px;">
					<a class="btn pri" href="{{.PayURL}}" style="display:flex; justify-content:center; width:100%; text-align:center; padding:16px; font-size:16px; border-radius:18px; font-weight:800; box-shadow:0 10px 28px rgba(37,99,235,0.35); transition:transform 0.2s; box-sizing:border-box;">
						<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="margin-right:10px;"><rect x="2" y="5" width="20" height="14" rx="2"></rect><line x1="2" y1="10" x2="22" y2="10"></line></svg>
						Bayar Sekarang
					</a>
				</div>
			{{else}}
				{{if not .Tx.QRISString}}
				<div style="margin-top:24px; padding:18px; background:#fff1f2; border:1px dashed #f43f5e; border-radius:18px; color:#be123c; text-align:center; font-size:13.5px; font-weight:600; box-shadow:0 4px 12px rgba(225,29,72,0.05);">
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="vertical-align:middle; margin-right:6px; margin-top:-2px;"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
					Gateway Pembayaran (Pakasir) belum di-setup oleh Admin.
				</div>
				{{end}}
			{{end}}
		{{end}}
	</div>
</div>
<script>
document.querySelectorAll('.qris-expire-date').forEach(function(el){
	var d = new Date(el.getAttribute('data-date'));
	if(!isNaN(d)){
		el.innerText = d.toLocaleString('id-ID', {day:'numeric', month:'short', year:'numeric', hour:'2-digit', minute:'2-digit'}) + ' WIB';
	}
});
</script>


<style>
.invoice-info-grid-v2{
	display:grid;
	grid-template-columns:repeat(2,minmax(0,1fr));
	gap:12px;
	margin:14px 0 16px;
}
.invoice-info-grid-v2 > .notice,
.invoice-info-grid-v2 > .card,
.invoice-info-grid-v2 > div{
	margin:0 !important;
	min-width:0;
}
.invoice-info-grid-v2 .invoice-user-v2{
	grid-column:1 / -1;
}
@media(max-width:640px){
	.invoice-info-grid-v2{
		grid-template-columns:repeat(2,minmax(0,1fr));
		gap:10px;
		margin:10px 0 12px;
	}
	.invoice-info-grid-v2 > .notice,
	.invoice-info-grid-v2 > .card,
	.invoice-info-grid-v2 > div{
		padding:14px !important;
		border-radius:18px !important;
	}
	.invoice-info-grid-v2 .price{
		font-size:30px !important;
		line-height:1.05 !important;
		white-space:nowrap;
	}
	.invoice-info-grid-v2 h3{
		font-size:24px !important;
		margin:6px 0 0 !important;
	}
}
</style>

<script>
(function(){
	function norm(s){
		return (s || '').replace(/\s+/g, ' ').trim();
	}
	function findExact(label){
		var nodes = document.querySelectorAll('.invoice-wrap *');
		for(var i=0;i<nodes.length;i++){
			if(norm(nodes[i].textContent) === label){
				return nodes[i];
			}
		}
		return null;
	}
	function findBlock(label){
		var el = findExact(label);
		if(!el) return null;

		var b = el.closest('.notice');
		if(b) return b;

		// fallback: ambil div parent yang paling dekat, jangan sampai card utama invoice
		b = el.parentElement;
		while(b && b.parentElement && !b.classList.contains('invoice-wrap')){
			var txt = norm(b.textContent);
			if(txt.indexOf(label) === 0 && txt.length < 180){
				return b;
			}
			b = b.parentElement;
		}
		return el.parentElement;
	}
	function applyInvoiceGrid(){
		if(document.querySelector('.invoice-info-grid-v2')) return;

		var amount = findBlock('Nominal Pembelian');
		var user = findBlock('Data Pengguna');
		var status = findBlock('Status');

		if(!amount || !user || !status) return;
		if(amount === user || amount === status || user === status) return;

		var parent = amount.parentNode;
		if(!parent) return;

		var grid = document.createElement('div');
		grid.className = 'invoice-info-grid-v2';

		parent.insertBefore(grid, amount);
		grid.appendChild(amount);
		grid.appendChild(status);

		user.classList.add('invoice-user-v2');
		grid.appendChild(user);
	}

	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', applyInvoiceGrid);
	}else{
		applyInvoiceGrid();
	}
	setTimeout(applyInvoiceGrid, 300);
})();
</script>



<!-- INVOICE_HEADER_AMOUNT_V5 -->
<style>
.invoice-head-amount-v5{
	margin-left:auto;
	width:min(240px,46vw);
	min-width:150px;
}
.invoice-head-amount-v5 > div,
.invoice-head-amount-v5 > .notice,
.invoice-head-amount-v5 > .card{
	margin:0 !important;
	padding:14px 16px !important;
	border-radius:18px !important;
}
.invoice-head-amount-v5 .muted{
	font-size:13px !important;
	margin:0 0 4px !important;
}
.invoice-head-amount-v5 .price{
	font-size:28px !important;
	line-height:1.05 !important;
	white-space:nowrap;
}
.invoice-bottom-hide-v5{
	display:none !important;
}
.invoice-user-only-v5{
	margin-top:12px !important;
}
@media(max-width:640px){
	.invoice-head-amount-v5{
		width:46%;
		min-width:135px;
	}
	.invoice-head-amount-v5 > div,
	.invoice-head-amount-v5 > .notice,
	.invoice-head-amount-v5 > .card{
		padding:12px !important;
		border-radius:16px !important;
	}
	.invoice-head-amount-v5 .price{
		font-size:24px !important;
	}
	.invoice-head-amount-v5 .muted{
		font-size:12px !important;
	}
}
</style>

<script>
(function(){
	function norm(s){
		return (s || '').replace(/\s+/g, ' ').trim();
	}

	function findExact(label){
		var nodes = document.querySelectorAll('.invoice-wrap *');
		for(var i=0;i<nodes.length;i++){
			if(norm(nodes[i].textContent) === label){
				return nodes[i];
			}
		}
		return null;
	}

	function findBlock(label){
		var el = findExact(label);
		if(!el) return null;

		var b = el.closest('.notice');
		if(b) return b;

		b = el.parentElement;
		while(b && b.parentElement && !b.classList.contains('invoice-wrap')){
			var txt = norm(b.textContent);
			if(txt.indexOf(label) === 0 && txt.length < 240){
				return b;
			}
			b = b.parentElement;
		}
		return el.parentElement;
	}

	function findHeaderActions(){
		var wrap = document.querySelector('.invoice-wrap');
		if(!wrap) return null;

		var actions = wrap.querySelectorAll('.actions');
		for(var i=0;i<actions.length;i++){
			if(actions[i].querySelector('h2')){
				return actions[i];
			}
		}
		return null;
	}

	function applyHeaderAmount(){
		if(document.querySelector('.invoice-head-amount-v5')) return;

		var amount = findBlock('Nominal Pembelian');
		var status = findBlock('Status');
		var user = findBlock('Data Pengguna');
		var header = findHeaderActions();

		if(!amount || !header) return;

		var cloneWrap = document.createElement('div');
		cloneWrap.className = 'invoice-head-amount-v5';

		var clone = amount.cloneNode(true);
		cloneWrap.appendChild(clone);

		header.appendChild(cloneWrap);

		amount.classList.add('invoice-bottom-hide-v5');
		if(status) status.classList.add('invoice-bottom-hide-v5');
		if(user) user.classList.add('invoice-user-only-v5');
	}

	function boot(){
		applyHeaderAmount();
		setTimeout(applyHeaderAmount, 350);
		setTimeout(applyHeaderAmount, 800);
	}

	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', boot);
	}else{
		boot();
	}
})();
</script>


<!-- INVOICE_LIGHT_POLISH_V6 -->
<style>
/* Invoice light-mode final polish */
html[data-theme="light"] .invoice-wrap,
html[data-theme="light"] .invoice-wrap *,
body[data-theme="light"] .invoice-wrap,
body[data-theme="light"] .invoice-wrap *{
	opacity:1 !important;
}

html[data-theme="light"] .invoice-wrap .card,
body[data-theme="light"] .invoice-wrap .card{
	background:rgba(255,255,255,.94) !important;
	color:#0f172a !important;
	border-color:rgba(15,23,42,.10) !important;
	box-shadow:0 18px 45px rgba(15,23,42,.08) !important;
}

html[data-theme="light"] .invoice-wrap h1,
html[data-theme="light"] .invoice-wrap h2,
html[data-theme="light"] .invoice-wrap h3,
html[data-theme="light"] .invoice-wrap .price,
html[data-theme="light"] .invoice-wrap .voucher-code,
body[data-theme="light"] .invoice-wrap h1,
body[data-theme="light"] .invoice-wrap h2,
body[data-theme="light"] .invoice-wrap h3,
body[data-theme="light"] .invoice-wrap .price,
body[data-theme="light"] .invoice-wrap .voucher-code{
	color:#0f172a !important;
	text-shadow:none !important;
}

html[data-theme="light"] .invoice-wrap .muted,
body[data-theme="light"] .invoice-wrap .muted{
	color:#64748b !important;
}

/* Badge pending jangan jadi oval besar */
.invoice-wrap .badge,
.invoice-wrap .badge.pending,
.invoice-wrap span.pending{
	display:inline-flex !important;
	align-items:center !important;
	justify-content:center !important;
	width:auto !important;
	height:auto !important;
	min-width:0 !important;
	min-height:0 !important;
	max-width:max-content !important;
	padding:8px 15px !important;
	border-radius:999px !important;
	font-size:15px !important;
	font-weight:900 !important;
	line-height:1 !important;
	letter-spacing:0 !important;
}

.invoice-wrap .badge.pending,
.invoice-wrap span.pending{
	background:#fff7ed !important;
	color:#b45309 !important;
	border:1px solid #fdba74 !important;
	box-shadow:0 6px 16px rgba(245,158,11,.10) !important;
}

/* Nominal kanan atas */
.invoice-head-amount-v5{
	align-self:center !important;
	width:clamp(145px,42vw,220px) !important;
}
.invoice-head-amount-v5 > div,
.invoice-head-amount-v5 > .notice,
.invoice-head-amount-v5 > .card{
	background:rgba(255,255,255,.84) !important;
	border:1px solid rgba(148,163,184,.30) !important;
	box-shadow:none !important;
	padding:12px 14px !important;
	border-radius:16px !important;
}
html[data-theme="light"] .invoice-head-amount-v5 > div,
body[data-theme="light"] .invoice-head-amount-v5 > div{
	background:#ffffff !important;
}
.invoice-head-amount-v5 .muted{
	font-size:12px !important;
	line-height:1.2 !important;
	margin:0 0 5px !important;
}
.invoice-head-amount-v5 .price{
	font-size:25px !important;
	line-height:1.05 !important;
	white-space:nowrap !important;
}

/* Data pengguna lebih compact */
.invoice-user-only-v5,
.invoice-info-grid-v2 .invoice-user-v2{
	margin-top:12px !important;
}
.invoice-user-only-v5,
.invoice-user-only-v5.notice,
.invoice-user-only-v5.card{
	padding:16px 18px !important;
	border-radius:18px !important;
}

/* QRIS card light/dark lebih jelas */
.invoice-wrap .qris-pay-box{
	color:inherit !important;
}
.invoice-wrap .qris-pay-box h3{
	color:inherit !important;
	margin-top:0 !important;
}
html[data-theme="light"] .invoice-wrap .qris-pay-box h3,
body[data-theme="light"] .invoice-wrap .qris-pay-box h3{
	color:#0f172a !important;
}
html[data-theme="light"] .invoice-wrap .qris-pay-box .muted,
body[data-theme="light"] .invoice-wrap .qris-pay-box .muted{
	color:#475569 !important;
}
.invoice-wrap #qris-pay-img{
	background:#ffffff !important;
	border:1px solid #e2e8f0 !important;
	box-shadow:0 10px 28px rgba(15,23,42,.08) !important;
}
.qris-download-note{
	color:#64748b !important;
}

/* Mobile compact */
@media(max-width:640px){
	.invoice-wrap .card{
		padding:18px !important;
	}
	.invoice-wrap .actions{
		gap:10px !important;
		align-items:flex-start !important;
	}
	.invoice-wrap .actions h2{
		font-size:26px !important;
		line-height:1.1 !important;
	}
	.invoice-wrap .actions .muted{
		font-size:14px !important;
		margin-top:8px !important;
	}
	.invoice-wrap .badge,
	.invoice-wrap .badge.pending,
	.invoice-wrap span.pending{
		padding:7px 13px !important;
		font-size:14px !important;
	}
	.invoice-head-amount-v5{
		width:42% !important;
		min-width:132px !important;
	}
	.invoice-head-amount-v5 > div,
	.invoice-head-amount-v5 > .notice,
	.invoice-head-amount-v5 > .card{
		padding:11px 12px !important;
		border-radius:15px !important;
	}
	.invoice-head-amount-v5 .muted{
		font-size:11px !important;
	}
	.invoice-head-amount-v5 .price{
		font-size:22px !important;
	}
	.invoice-user-only-v5,
	.invoice-user-only-v5.notice,
	.invoice-user-only-v5.card{
		padding:15px !important;
	}
	.invoice-wrap .qris-pay-box h3{
		font-size:27px !important;
		line-height:1.15 !important;
	}
}
</style>


<!-- INVOICE_MOBILE_COMPACT_V7 -->
<style>
.invoice-user-phone-v7{
	font-size:18px !important;
	font-weight:800 !important;
	line-height:1.2 !important;
	color:inherit !important;
	margin-top:6px !important;
}
.invoice-user-only-v5,
.invoice-user-only-v5.notice,
.invoice-user-only-v5.card{
	margin-bottom:10px !important;
}
.invoice-wrap .qris-pay-box{
	margin-top:6px !important;
}
.invoice-wrap .qris-pay-box h3{
	margin-bottom:10px !important;
}

@media(max-width:640px){
	.invoice-wrap .actions{
		align-items:flex-start !important;
		gap:10px !important;
	}
	.invoice-wrap .actions h2{
		font-size:22px !important;
		line-height:1.08 !important;
		letter-spacing:-.02em !important;
		margin:0 6px 0 0 !important;
		overflow-wrap:anywhere !important;
		word-break:break-word !important;
	}
	.invoice-wrap .actions .muted{
		font-size:14px !important;
		margin-top:6px !important;
	}
	.invoice-head-amount-v5{
		width:40% !important;
		min-width:126px !important;
	}
	.invoice-head-amount-v5 > div,
	.invoice-head-amount-v5 > .notice,
	.invoice-head-amount-v5 > .card{
		padding:10px 12px !important;
	}
	.invoice-head-amount-v5 .muted{
		font-size:11px !important;
		line-height:1.15 !important;
		white-space:normal !important;
	}
	.invoice-head-amount-v5 .price{
		font-size:20px !important;
		line-height:1.05 !important;
	}

	.invoice-user-only-v5,
	.invoice-user-only-v5.notice,
	.invoice-user-only-v5.card{
		padding:12px 14px !important;
		margin-top:10px !important;
		margin-bottom:8px !important;
		border-radius:16px !important;
	}
	.invoice-user-only-v5 .muted{
		font-size:13px !important;
		margin-bottom:2px !important;
	}
	.invoice-user-phone-v7{
		font-size:16px !important;
		margin-top:4px !important;
	}
	.invoice-wrap .qris-pay-box{
		margin-top:2px !important;
	}
	.invoice-wrap .qris-pay-box h3{
		font-size:21px !important;
		line-height:1.15 !important;
		margin:0 0 8px !important;
	}
	.invoice-wrap #qris-pay-img{
		width:260px !important;
		max-width:100% !important;
	}
	.qris-download-note{
		font-size:12px !important;
		margin-top:6px !important;
	}
}
</style>

<script>
(function(){
	function norm(s){
		return (s || '').replace(/\s+/g, ' ').trim();
	}
	function findExact(label){
		var nodes = document.querySelectorAll('.invoice-wrap *');
		for(var i=0;i<nodes.length;i++){
			if(norm(nodes[i].textContent) === label){
				return nodes[i];
			}
		}
		return null;
	}
	function findBlock(label){
		var el = findExact(label);
		if(!el) return null;

		var b = el.closest('.notice');
		if(b) return b;

		b = el.parentElement;
		while(b && b.parentElement && !b.classList.contains('invoice-wrap')){
			var txt = norm(b.textContent);
			if(txt.indexOf(label) === 0 && txt.length < 260){
				return b;
			}
			b = b.parentElement;
		}
		return el.parentElement;
	}
	function extractPhone(txt){
		if(!txt) return '';
		var m = txt.match(/(?:\+62|62|0)\d[\d\s-]{7,}/);
		return m ? m[0].replace(/\s+/g,'').trim() : '';
	}
	function compactUserBlock(){
		var user = findBlock('Data Pengguna');
		if(!user) return;
		if(user.getAttribute('data-v7-compact') === '1') return;

		var phone = extractPhone(user.textContent || '');
		if(!phone) return;

		user.setAttribute('data-v7-compact','1');
		user.innerHTML =
			'<div class="muted">No. WhatsApp Anda</div>' +
			'<div class="invoice-user-phone-v7">' + phone + '</div>';
	}
	function boot(){
		compactUserBlock();
		setTimeout(compactUserBlock, 250);
		setTimeout(compactUserBlock, 700);
	}
	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', boot);
	}else{
		boot();
	}
})();
</script>


<!-- INVOICE_MOBILE_TIGHTEN_V8 -->
<style>
@media(max-width:640px){
	/* area header lebih rapat */
	.invoice-wrap .actions{
		gap:8px !important;
		margin-bottom:4px !important;
	}
	.invoice-wrap .actions h2{
		font-size:21px !important;
		line-height:1.05 !important;
		margin:0 4px 0 0 !important;
	}
	.invoice-wrap .actions .muted{
		margin-top:4px !important;
		font-size:13px !important;
	}

	/* badge + nominal lebih mepet */
	.invoice-head-amount-v5{
		width:39% !important;
		min-width:122px !important;
	}
	.invoice-head-amount-v5 > div,
	.invoice-head-amount-v5 > .notice,
	.invoice-head-amount-v5 > .card{
		padding:9px 11px !important;
		border-radius:14px !important;
	}
	.invoice-head-amount-v5 .muted{
		font-size:10px !important;
		line-height:1.1 !important;
		margin:0 0 4px !important;
	}
	.invoice-head-amount-v5 .price{
		font-size:19px !important;
		line-height:1.02 !important;
	}

	.invoice-wrap .badge,
	.invoice-wrap .badge.pending,
	.invoice-wrap span.pending{
		padding:6px 12px !important;
		font-size:13px !important;
	}

	/* garis divider dan blok bawah dipadatkan */
	.invoice-wrap hr{
		margin:10px 0 !important;
	}
	.invoice-user-only-v5,
	.invoice-user-only-v5.notice,
	.invoice-user-only-v5.card{
		margin-top:6px !important;
		margin-bottom:4px !important;
		padding:10px 12px !important;
		border-radius:14px !important;
	}
	.invoice-user-only-v5 .muted{
		font-size:12px !important;
		margin:0 0 2px !important;
	}
	.invoice-user-phone-v7{
		font-size:15px !important;
		line-height:1.15 !important;
		margin-top:2px !important;
	}

	/* QRIS naik dan lebih compact */
	.invoice-wrap .qris-pay-box{
		margin-top:0 !important;
		padding-top:0 !important;
	}
	.invoice-wrap .qris-pay-box h3{
		font-size:18px !important;
		line-height:1.12 !important;
		margin:0 0 6px !important;
	}
	.invoice-wrap #qris-pay-img{
		width:240px !important;
		max-width:100% !important;
	}
	.qris-download-note{
		font-size:11px !important;
		margin-top:4px !important;
		margin-bottom:2px !important;
	}
}
</style>


<!-- INVOICE_NOTE_WA_TIGHT_V9 -->
<style>
.invoice-wrap .qris-download-note{
	color:#cbd5e1 !important;
	font-size:13px !important;
	font-weight:700 !important;
	line-height:1.35 !important;
	margin-top:8px !important;
	margin-bottom:4px !important;
	text-align:center !important;
	opacity:1 !important;
}
html[data-theme="light"] .invoice-wrap .qris-download-note,
body[data-theme="light"] .invoice-wrap .qris-download-note{
	color:#334155 !important;
}

@media(max-width:640px){
	/* rapetin lagi area tengah -> qris */
	.invoice-wrap hr{
		margin:8px 0 !important;
	}
	.invoice-user-only-v5,
	.invoice-user-only-v5.notice,
	.invoice-user-only-v5.card{
		margin-top:4px !important;
		margin-bottom:2px !important;
		padding:9px 12px !important;
	}
	.invoice-user-only-v5 .muted{
		font-size:11px !important;
		margin:0 0 1px !important;
	}
	.invoice-user-phone-v7{
		font-size:14px !important;
		margin-top:1px !important;
	}

	.invoice-wrap .qris-pay-box{
		margin-top:-2px !important;
		padding-top:0 !important;
	}
	.invoice-wrap .qris-pay-box h3{
		font-size:17px !important;
		line-height:1.10 !important;
		margin:0 0 5px !important;
	}
	.invoice-wrap #qris-pay-img{
		width:236px !important;
		max-width:100% !important;
	}

	.qris-download-note{
		font-size:12px !important;
		line-height:1.32 !important;
		margin-top:6px !important;
		margin-bottom:3px !important;
	}

	.invoice-wrap .qris-pay-box .muted{
		font-size:11px !important;
		line-height:1.28 !important;
	}
}
</style>


<style>
/* QRIS_BYPASS_BUTTON_V11 */
.qris-bypass-box-v11{
	margin:12px auto 8px;
	padding:12px 13px;
	border:1px solid rgba(59,130,246,.30);
	background:rgba(37,99,235,.10);
	border-radius:16px;
	text-align:center;
	max-width:360px;
}
.qris-bypass-text-v11{
	font-size:13px;
	line-height:1.45;
	font-weight:800;
	color:#dbeafe;
	margin-bottom:10px;
}
.qris-bypass-btn-v11{
	width:100%;
	border:0;
	border-radius:14px;
	padding:11px 14px;
	background:#2563eb;
	color:#fff;
	font-weight:900;
	font-size:13px;
	cursor:pointer;
}
.qris-bypass-btn-v11:disabled{
	opacity:.75;
	cursor:not-allowed;
}
.qris-bypass-status-v11{
	margin-top:8px;
	font-size:12px;
	line-height:1.35;
	color:#cbd5e1;
}
html[data-theme="light"] .qris-bypass-box-v11,
body[data-theme="light"] .qris-bypass-box-v11{
	background:#eff6ff;
	border-color:#bfdbfe;
}
html[data-theme="light"] .qris-bypass-text-v11,
body[data-theme="light"] .qris-bypass-text-v11{
	color:#1e3a8a;
}
html[data-theme="light"] .qris-bypass-status-v11,
body[data-theme="light"] .qris-bypass-status-v11{
	color:#475569;
}
@media(max-width:640px){
	.qris-bypass-box-v11{
		margin:9px auto 6px;
		padding:10px 11px;
		border-radius:14px;
	}
	.qris-bypass-text-v11{
		font-size:12px;
		margin-bottom:8px;
	}
	.qris-bypass-btn-v11{
		padding:10px 12px;
		font-size:12px;
	}
	.qris-bypass-status-v11{
		font-size:11px;
		margin-top:7px;
	}
}
</style>

<script>
(function(){
	if(window.__qrisBypassButtonV11) return;
	window.__qrisBypassButtonV11 = true;

	function bootBypassButton(){
		var btn = document.getElementById('qrisBypassBtnV11');
		var st = document.getElementById('qrisBypassStatusV11');
		if(!btn) return;

		var running = false;

		function setStatus(text){
			if(st) st.textContent = text;
		}

		function callBypass(manual){
			if(running) return;
			running = true;

			btn.disabled = true;
			btn.textContent = 'Mengaktifkan...';
			setStatus('Sedang mengaktifkan bypass internet 60 detik. Maksimal 2 kali per invoice...');

			fetch('/payment-bypass-start/{{.Tx.OrderID}}', {
				cache: 'no-store',
				keepalive: true
			}).then(function(resp){
				if(!resp.ok) throw new Error('HTTP ' + resp.status);
				btn.textContent = 'Bypass Aktif 60 Detik';
				setStatus('Bypass aktif. Silakan selesaikan pembayaran QRIS.');
			}).catch(function(){
				running = false;
				btn.disabled = false;
				btn.textContent = 'Coba Aktifkan Lagi';
				setStatus('Bypass belum aktif. Tekan tombol ini lagi atau refresh halaman.');
			});
		}

		btn.addEventListener('click', function(e){
			e.preventDefault();
			callBypass(true);
		});

		// Auto coba aktifkan setelah QRIS tampil, tombol tetap jadi fallback manual.
		var img = document.getElementById('qris-pay-img');
		function autoStart(){
			setTimeout(function(){
				callBypass(false);
			}, 500);
		}

		if(img && img.complete){
			autoStart();
		}else if(img){
			img.addEventListener('load', autoStart, {once:true});
		}else{
			setTimeout(autoStart, 700);
		}
	}

	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', bootBypassButton);
	}else{
		bootBypassButton();
	}
})();
</script>


<style>
/* QRIS_BYPASS_REFRESH_BUTTON_V12 */
.qris-bypass-box-v11{
	margin:7px auto 4px !important;
	padding:0 !important;
	border:0 !important;
	background:transparent !important;
	box-shadow:none !important;
	border-radius:0 !important;
	max-width:360px !important;
}
.qris-bypass-text-v11{
	margin:0 0 7px !important;
	font-size:12px !important;
	line-height:1.35 !important;
	font-weight:800 !important;
	color:#cbd5e1 !important;
}
.qris-bypass-btn-v11{
	width:auto !important;
	min-width:220px !important;
	max-width:100% !important;
	border-radius:999px !important;
	padding:10px 16px !important;
	font-size:12px !important;
	line-height:1 !important;
	cursor:pointer !important;
	opacity:1 !important;
}
.qris-bypass-status-v11{
	margin-top:6px !important;
	font-size:11px !important;
	line-height:1.3 !important;
	color:#94a3b8 !important;
}
@media(max-width:640px){
	.qris-bypass-box-v11{
		margin:5px auto 3px !important;
	}
	.qris-bypass-text-v11{
		font-size:11px !important;
		margin-bottom:6px !important;
	}
	.qris-bypass-btn-v11{
		min-width:205px !important;
		padding:9px 14px !important;
		font-size:11px !important;
	}
	.qris-bypass-status-v11{
		font-size:10px !important;
		margin-top:5px !important;
	}
}
html[data-theme="light"] .qris-bypass-text-v11,
body[data-theme="light"] .qris-bypass-text-v11{
	color:#334155 !important;
}
html[data-theme="light"] .qris-bypass-status-v11,
body[data-theme="light"] .qris-bypass-status-v11{
	color:#64748b !important;
}
</style>

<script>
/* QRIS_BYPASS_REFRESH_BUTTON_V12 */
document.addEventListener('DOMContentLoaded', function(){
	function setupRefreshButton(){
		var oldBtn = document.getElementById('qrisBypassBtnV11');
		var st = document.getElementById('qrisBypassStatusV11');
		if(!oldBtn) return false;

		// Clone untuk buang listener lama yang tembak fetch/disable.
		var btn = oldBtn.cloneNode(true);
		oldBtn.parentNode.replaceChild(btn, oldBtn);

		btn.disabled = false;
		btn.textContent = 'Refresh untuk Bypass 60 Detik (maks 2x)';
		btn.dataset.refreshButtonV12 = '1';

		if(st){
			st.textContent = 'Jika internet belum terbuka, tekan tombol ini untuk refresh. Bypass maksimal 2 kali per invoice.';
		}

		btn.addEventListener('click', function(e){
			e.preventDefault();
			btn.textContent = 'Refresh...';
			window.location.reload();
		});

		return true;
	}

	setupRefreshButton();
	setTimeout(setupRefreshButton, 300);
	setTimeout(setupRefreshButton, 900);
});
</script>


<style>
/* INVOICE_SINGLE_BOX_V13 */

/* outer wrapper tetap jadi box utama */
.invoice-wrap{
	padding:14px !important;
}

/* inner box / box kedua dibikin rata, jadi gak numpuk lagi */
.invoice-wrap > div:first-child,
.invoice-wrap .invoice-card,
.invoice-wrap .invoice-inner,
.invoice-wrap .invoice-body{
	background:transparent !important;
	border:0 !important;
	box-shadow:none !important;
	border-radius:0 !important;
	padding:0 !important;
}

/* rapihin area atas */
.invoice-wrap h1,
.invoice-wrap h2,
.invoice-wrap .invoice-title{
	margin-top:0 !important;
	margin-bottom:4px !important;
	line-height:1.1 !important;
}

/* garis pemisah jangan terlalu boros */
.invoice-wrap hr{
	margin:12px 0 !important;
}

/* card kecil nominal / nomor WA lebih compact */
.invoice-wrap .mini-card,
.invoice-wrap .info-card,
.invoice-wrap .invoice-info-card{
	padding:12px 14px !important;
	border-radius:18px !important;
}

/* area qris lebih hemat tempat */
.invoice-wrap .qris-pay-box{
	margin-top:12px !important;
	padding:12px 12px 10px !important;
	border-radius:22px !important;
}

/* qr image naik dikit */
.invoice-wrap #qris-pay-img{
	margin:6px auto 8px !important;
}

/* note-note bawah lebih rapat */
.invoice-wrap .qris-download-note{
	margin-top:4px !important;
	margin-bottom:4px !important;
	line-height:1.3 !important;
}

.invoice-wrap .qris-bypass-box-v11{
	margin-top:4px !important;
	margin-bottom:4px !important;
}

.invoice-wrap .muted{
	line-height:1.25 !important;
}

/* mobile */
@media(max-width:640px){
	.invoice-wrap{
		padding:12px !important;
		border-radius:20px !important;
	}

	.invoice-wrap > div:first-child,
	.invoice-wrap .invoice-card,
	.invoice-wrap .invoice-inner,
	.invoice-wrap .invoice-body{
		padding:0 !important;
	}

	.invoice-wrap h1,
	.invoice-wrap h2,
	.invoice-wrap .invoice-title{
		font-size:15px !important;
	}

	.invoice-wrap .mini-card,
	.invoice-wrap .info-card,
	.invoice-wrap .invoice-info-card{
		padding:10px 12px !important;
		border-radius:16px !important;
	}

	.invoice-wrap .qris-pay-box{
		margin-top:10px !important;
		padding:10px 10px 8px !important;
		border-radius:18px !important;
	}

	.invoice-wrap #qris-pay-img{
		width:220px !important;
		margin:4px auto 6px !important;
	}

	.invoice-wrap .qris-download-note{
		font-size:11px !important;
		margin-top:3px !important;
		margin-bottom:3px !important;
	}

	.invoice-wrap .muted{
		font-size:10px !important;
		line-height:1.2 !important;
	}
}
</style>


<style>
/* INVOICE_TUNE_V14 */

/* box utama lebih compact */
.invoice-wrap{
	padding:12px !important;
	border-radius:22px !important;
}

/* inner content rapihin */
.invoice-wrap > div:first-child{
	padding:10px 10px 12px !important;
	border-radius:18px !important;
}

/* cukup 1 garis rgb: matikan rgb line / top line di dalam */
.invoice-wrap > div:first-child::before,
.invoice-wrap .invoice-card::before,
.invoice-wrap .invoice-inner::before,
.invoice-wrap .invoice-body::before{
	content:none !important;
	display:none !important;
	background:none !important;
}

/* kalau ada border-top / shadow tambahan di inner, matikan */
.invoice-wrap > div:first-child,
.invoice-wrap .invoice-card,
.invoice-wrap .invoice-inner,
.invoice-wrap .invoice-body{
	box-shadow:none !important;
}

/* judul invoice jangan terlalu makan tempat */
.invoice-wrap h1,
.invoice-wrap h2,
.invoice-wrap .invoice-title{
	margin:0 0 4px !important;
	line-height:1.06 !important;
	font-size:15px !important;
}

/* subjudul kecil */
.invoice-wrap .muted,
.invoice-wrap .invoice-subtitle{
	line-height:1.22 !important;
}

/* area status + nominal lebih rapat */
.invoice-wrap .invoice-info-grid-v2,
.invoice-wrap .invoice-top-grid,
.invoice-wrap .invoice-grid{
	gap:10px !important;
	margin:10px 0 !important;
}

/* garis pemisah rapat */
.invoice-wrap hr{
	margin:10px 0 !important;
}

/* card kecil */
.invoice-wrap .mini-card,
.invoice-wrap .info-card,
.invoice-wrap .invoice-info-card{
	padding:10px 12px !important;
	border-radius:16px !important;
}

/* khusus nomor WA */
.invoice-wrap .invoice-info-card b,
.invoice-wrap .info-card b{
	line-height:1.05 !important;
}

/* kotak QRIS lebih hemat */
.invoice-wrap .qris-pay-box{
	margin-top:10px !important;
	padding:10px 10px 8px !important;
	border-radius:18px !important;
}

.invoice-wrap .qris-pay-box h3,
.invoice-wrap .qris-pay-box h4{
	margin:0 0 8px !important;
	line-height:1.08 !important;
	font-size:13px !important;
}

.invoice-wrap #qris-pay-img{
	margin:4px auto 6px !important;
	width:220px !important;
	max-width:100% !important;
}

.invoice-wrap .qris-download-note{
	margin:2px 0 4px !important;
	font-size:11px !important;
	line-height:1.28 !important;
}

/* bypass box compact dan teks gak pecah aneh */
.qris-bypass-box-v11{
	margin:4px auto 3px !important;
	max-width:320px !important;
}

.qris-bypass-text-v11{
	margin:0 auto 6px !important;
	max-width:300px !important;
	font-size:11px !important;
	line-height:1.32 !important;
	text-align:center !important;
	word-break:normal !important;
	overflow-wrap:break-word !important;
}

.qris-bypass-btn-v11{
	min-width:210px !important;
	padding:9px 14px !important;
	font-size:11px !important;
}

.qris-bypass-status-v11{
	margin-top:5px !important;
	font-size:10px !important;
	line-height:1.28 !important;
	max-width:300px !important;
	margin-left:auto !important;
	margin-right:auto !important;
}

/* total dan expired lebih rapat */
.invoice-wrap .qris-pay-box .muted{
	font-size:10px !important;
	line-height:1.18 !important;
	margin-top:2px !important;
}

/* status bawah juga rapat */
.invoice-wrap .payment-status-note,
.invoice-wrap .payment-note,
.invoice-wrap .status-note{
	margin-top:8px !important;
	line-height:1.28 !important;
}

/* light mode tetap nyaman */
html[data-theme="light"] .invoice-wrap .qris-pay-box,
body[data-theme="light"] .invoice-wrap .qris-pay-box{
	background:rgba(255,255,255,.62) !important;
}

/* mobile */
@media(max-width:640px){
	.invoice-wrap{
		padding:10px !important;
		border-radius:20px !important;
	}
	.invoice-wrap > div:first-child{
		padding:8px 8px 10px !important;
		border-radius:16px !important;
	}
	.invoice-wrap h1,
	.invoice-wrap h2,
	.invoice-wrap .invoice-title{
		font-size:14px !important;
		margin-bottom:3px !important;
	}
	.invoice-wrap .mini-card,
	.invoice-wrap .info-card,
	.invoice-wrap .invoice-info-card{
		padding:9px 11px !important;
		border-radius:15px !important;
	}
	.invoice-wrap .qris-pay-box{
		margin-top:8px !important;
		padding:8px 8px 7px !important;
		border-radius:16px !important;
	}
	.invoice-wrap #qris-pay-img{
		width:210px !important;
		margin:3px auto 5px !important;
	}
	.invoice-wrap .qris-download-note{
		font-size:10px !important;
		margin:2px 0 3px !important;
	}
	.qris-bypass-box-v11{
		margin:3px auto 2px !important;
	}
	.qris-bypass-text-v11{
		font-size:10px !important;
		margin-bottom:5px !important;
	}
	.qris-bypass-btn-v11{
		min-width:195px !important;
		padding:8px 12px !important;
		font-size:10px !important;
	}
	.qris-bypass-status-v11{
		font-size:9.5px !important;
		margin-top:4px !important;
	}
}
</style>


<style>
/* INVOICE_WA_BOTTOM_V15 */
.invoice-wa-bottom-v15{
	margin:8px auto 4px !important;
	padding:8px 12px !important;
	border-radius:14px !important;
	max-width:360px !important;
	text-align:center !important;
	background:rgba(15,23,42,.12) !important;
}
.invoice-wa-bottom-v15 .muted{
	font-size:10px !important;
	line-height:1.2 !important;
	margin:0 0 2px !important;
}
.invoice-wa-bottom-v15 .invoice-user-phone-v7{
	font-size:13px !important;
	line-height:1.15 !important;
	margin:0 !important;
}
html[data-theme="light"] .invoice-wa-bottom-v15,
body[data-theme="light"] .invoice-wa-bottom-v15{
	background:rgba(255,255,255,.62) !important;
}
.invoice-wrap .qris-pay-box{
	margin-top:6px !important;
}
@media(max-width:640px){
	.invoice-wa-bottom-v15{
		margin:6px auto 3px !important;
		padding:7px 10px !important;
		max-width:310px !important;
	}
	.invoice-wa-bottom-v15 .muted{
		font-size:9.5px !important;
	}
	.invoice-wa-bottom-v15 .invoice-user-phone-v7{
		font-size:12px !important;
	}
	.invoice-wrap .qris-pay-box{
		margin-top:4px !important;
	}
}

/* QRIS_PAY_TITLE_V1 */
.qris-pay-title-v1{
  margin:0 0 6px;
  font-size:clamp(24px,5vw,34px);
  line-height:1.12;
  font-weight:900;
  letter-spacing:-.04em;
  text-align:center;
}
.qris-pay-subtitle-v1{
  margin:0 0 14px;
  text-align:center;
  font-size:15px;
  opacity:.82;
}


/* INVOICE_PENDING_ACTIONS_V1 */
.invoice-pending-actions-v1{
  justify-content:center;
  margin:14px 0 4px;
  gap:10px;
}
.invoice-pending-actions-v1 .btn{
  min-width:140px;
}
@media(max-width:640px){
  .invoice-pending-actions-v1{
    display:grid;
    grid-template-columns:1fr;
  }
  .invoice-pending-actions-v1 .btn{
    width:100%;
  }
}

</style>

<script>
/* INVOICE_WA_BOTTOM_V15 */
(function(){
	function norm(s){
		return (s || '').replace(/\s+/g, ' ').trim();
	}

	function findUserBlock(){
		var nodes = document.querySelectorAll('.invoice-wrap *');
		for(var i=0;i<nodes.length;i++){
			var t = norm(nodes[i].textContent);
			if(t === 'No. WhatsApp Anda' || t === 'Nomor HP' || t.indexOf('No. WhatsApp Anda') === 0){
				var b = nodes[i].closest('.notice') || nodes[i].closest('.card') || nodes[i].parentElement;
				return b;
			}
		}
		return null;
	}

	function moveWA(){
		var user = findUserBlock();
		var qris = document.querySelector('.invoice-wrap .qris-pay-box');
		if(!user || !qris) return false;
		if(user.classList.contains('invoice-wa-bottom-v15')) return true;

		user.classList.add('invoice-wa-bottom-v15');

		// Pindah ke bawah box QRIS supaya QRIS naik ke atas.
		if(qris.parentNode){
			qris.parentNode.insertBefore(user, qris.nextSibling);
		}
		return true;
	}

	function boot(){
		moveWA();
		setTimeout(moveWA, 200);
		setTimeout(moveWA, 700);
	}

	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', boot);
	}else{
		boot();
	}
})();
</script>

<script>
function copyVoucher(){
  var el = document.getElementById('voucher-copy-text');
  var text = el ? el.innerText : '';
  function done(){
    var s = document.getElementById('copy-status');
    if(s){ s.innerText = 'Kode Voucher berhasil disalin.'; }
  }
  if(navigator.clipboard && window.isSecureContext){
    navigator.clipboard.writeText(text).then(done);
  }else{
    var ta = document.createElement('textarea');
    ta.value = text;
    document.body.appendChild(ta);
    ta.select();
    document.execCommand('copy');
    document.body.removeChild(ta);
    done();
  }
}
</script>`, map[string]any{
		"Tx":               x,
		"PayURL":           payURL,
		"QRImageURL":       template.URL("/qris/" + url.PathEscape(x.OrderID) + ".png"),
		"ChoosePackageURL": choosePackageURL,
		"CancelInvoiceURL": cancelInvoiceURL,
	})

	body = polishInvoiceHTML(body)
	body = template.HTML(strings.ReplaceAll(string(body), `href="/go-login"`, `href="`+loginURL+`"`))
	publicPage(w, "Invoice", body)
}
