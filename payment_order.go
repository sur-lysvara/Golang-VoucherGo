package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (a *App) orderPage(w http.ResponseWriter, r *http.Request) {
	pkgID := r.URL.Query().Get("package_id")
	if pkgID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var p Package
	var ro Router

	err := a.db.QueryRow(`SELECT
			p.id,
			p.router_id,
			p.name,
			p.price,
			COALESCE(p.profile,''),
			COALESCE(p.limit_uptime,''),
			r.name,
			r.slug
		FROM packages p
		JOIN routers r ON r.id=p.router_id
		WHERE p.id=? AND p.active=1 AND r.enabled=1`, pkgID).
		Scan(&p.ID, &p.RouterID, &p.Name, &p.Price, &p.Profile, &p.LimitUptime, &ro.Name, &ro.Slug)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	body := renderHTML(`
<div class="hero checkout-hero-v1"> <!-- CHECKOUT_STABLE_CLASSES_V1 -->
	<h1>Checkout Voucher</h1>
	<p>{{.Router.Name}} • {{.Package.Name}}</p>
</div>

<div class="card checkout-form-card-v1" style="max-width:640px;margin:0 auto 18px">
	<div class="notice">
		<div class="muted">Paket dipilih</div>
		<h2 style="margin:6px 0">{{.Package.Name}}</h2>
		<div class="price">{{rp .Package.Price}}</div>
		<p class="muted">Profile: {{.Package.Profile}} {{if .Package.LimitUptime}}• {{.Package.LimitUptime}}{{end}}</p>
	</div>

	<form method="post" action="/checkout" style="margin-top:16px" data-checkout-form-v1>
		<input type="hidden" name="package_id" value="{{.Package.ID}}">


			<input type="hidden" name="client_mac" value="{{.ClientMAC}}">
			<input type="hidden" name="client_ip" value="{{.ClientIP}}">
<div class="field">
			<label>Nama Pembeli</label>
			<input name="customer_name" placeholder="Opsional" autocomplete="name">
		</div>

		<div class="field">
			<label>No WhatsApp</label>
			<input name="customer_phone" placeholder="08xxxx" inputmode="tel" autocomplete="tel" required>
		</div>

		<div class="field">
			<label>Email</label>
			<input name="customer_email" type="email" placeholder="email@example.com" autocomplete="email">
		</div>

		<button class="btn pri" style="width:100%" type="submit" data-checkout-submit-v1>Lanjut Bayar</button>
	</form>

	<br>
	<a class="btn light" href="/{{.Router.Slug}}">← Pilih paket lain</a>
</div>

<script>
(function(){
	var checkoutSubmitting = false;

	function resetCheckoutSubmitState(){
		checkoutSubmitting = false;
		document.querySelectorAll('[data-checkout-form-v1]').forEach(function(form){
			var btn = form.querySelector('[data-checkout-submit-v1]');
			if(!btn) return;
			btn.disabled = false;
			btn.removeAttribute('aria-disabled');
			btn.classList.remove('loading');
			btn.textContent = btn.getAttribute('data-original-text') || 'Lanjut Bayar';
		});
	}

	function bootCheckoutSubmitReset(){
		document.querySelectorAll('[data-checkout-form-v1]').forEach(function(form){
			if(form.dataset.checkoutSubmitBooted === '1') return;
			form.dataset.checkoutSubmitBooted = '1';

			var btn = form.querySelector('[data-checkout-submit-v1]');
			if(btn && !btn.getAttribute('data-original-text')){
				btn.setAttribute('data-original-text', btn.textContent || 'Lanjut Bayar');
			}

			form.addEventListener('submit', function(e){
				if(checkoutSubmitting){
					e.preventDefault();
					return false;
				}
				checkoutSubmitting = true;
				if(btn){
					btn.disabled = true;
					btn.setAttribute('aria-disabled', 'true');
					btn.classList.add('loading');
					btn.textContent = 'Memproses...';
				}
				return true;
			});
		});

		resetCheckoutSubmitState();
	}

	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', bootCheckoutSubmitReset);
	}else{
		bootCheckoutSubmitReset();
	}

	window.addEventListener('pageshow', function(){
		resetCheckoutSubmitState();
	});
})();
</script>`, map[string]any{
		"Package":   p,
		"Router":    ro,
		"ClientMAC": cleanPaymentClientMAC(r.URL.Query().Get("mac")),
		"ClientIP":  cleanPaymentClientIP(r.URL.Query().Get("ip")),
	})

	publicPage(w, "Checkout "+p.Name, body)
}

func (a *App) checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()

	pkgID := r.FormValue("package_id")
	clientMAC := cleanPaymentClientMAC(r.FormValue("client_mac"))
	clientIP := cleanPaymentClientIP(r.FormValue("client_ip"))

	var ro Router
	var amount int64

	if err := a.db.QueryRow(`SELECT
			r.id,
			r.name,
			r.slug,
			COALESCE(r.hotspot_login_url,''),
			COALESCE(r.host,''),
			r.port,
			COALESCE(r.api_mode,'auto'),
			COALESCE(r.username,''),
			COALESCE(r.password,''),
			r.use_ssl,
			r.insecure_ssl,
			r.enabled,
			p.price
		FROM packages p
		JOIN routers r ON r.id=p.router_id
		WHERE p.id=? AND p.active=1 AND r.enabled=1`, pkgID).
		Scan(
			&ro.ID,
			&ro.Name,
			&ro.Slug,
			&ro.HotspotLoginURL,
			&ro.Host,
			&ro.Port,
			&ro.APIMode,
			&ro.Username,
			&ro.Password,
			&ro.UseSSL,
			&ro.InsecureSSL,
			&ro.Enabled,
			&amount,
		); err != nil {
		http.Error(w, "Paket tidak ditemukan", 404)
		return
	}

	returnURL := checkoutOrderURL(pkgID, clientMAC, clientIP)
	customerPhone := r.FormValue("customer_phone")
	blocked, banReason, err := isPhoneBlocked(r.Context(), a.db, customerPhone)
	if err != nil {
		log.Printf("phone ban check failed err=%v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		publicPage(w, "Validasi Nomor WhatsApp", renderHTML(`
<div class="card checkout-form-card-v1" style="max-width:640px;margin:0 auto 18px">
	<div class="notice" style="background:#fee2e2;color:#991b1b;border-color:#fecaca">
		Sistem sedang memvalidasi nomor WhatsApp, silakan coba lagi.
	</div>
	<p style="margin:16px 0 0"><a class="btn" href="{{.ReturnURL}}">Kembali</a></p>
</div>`, map[string]any{"ReturnURL": returnURL}))
		return
	}
	if blocked {
		w.WriteHeader(http.StatusForbidden)
		publicPage(w, "Nomor WhatsApp Diblokir", renderHTML(`
<style>
.blocked-phone-screen-v1{
	min-height:calc(100vh - 160px);
	display:flex;
	align-items:center;
	justify-content:center;
	padding:28px 14px 44px;
}
.blocked-phone-card-v1{
	width:min(100%,560px);
	position:relative;
	overflow:hidden;
	text-align:center;
	padding:30px 24px 24px;
	border-radius:26px;
	color:#fff1f2;
	background:
		radial-gradient(circle at 50% -12%, rgba(248,113,113,.32), transparent 34%),
		linear-gradient(145deg, #3f0712 0%, #190711 52%, #110713 100%);
	border:1px solid rgba(248,113,113,.42);
	box-shadow:0 26px 80px rgba(127,29,29,.34), inset 0 1px 0 rgba(255,255,255,.08);
}
.blocked-phone-card-v1::before{
	content:"";
	position:absolute;
	inset:0;
	pointer-events:none;
	background:linear-gradient(180deg, rgba(255,255,255,.08), transparent 36%);
}
.blocked-phone-inner-v1{
	position:relative;
	z-index:1;
}
.blocked-phone-icon-v1{
	position:relative;
	width:88px;
	height:88px;
	margin:0 auto 16px;
	display:grid;
	place-items:center;
	border-radius:26px;
	background:linear-gradient(145deg, rgba(69,10,10,.98), rgba(23,7,14,.92));
	border:1px solid rgba(254,202,202,.22);
	box-shadow:0 18px 44px rgba(220,38,38,.30);
}
.blocked-phone-icon-v1 svg{
	width:56px;
	height:56px;
	display:block;
	color:#fecaca;
	filter:drop-shadow(0 0 16px rgba(248,113,113,.28));
}
.blocked-phone-cross-v1{
	position:absolute;
	right:-7px;
	top:-7px;
	width:32px;
	height:32px;
	border-radius:999px;
	display:grid;
	place-items:center;
	font-size:20px;
	font-weight:900;
	line-height:1;
	color:#fff;
	background:#dc2626;
	border:2px solid #fecaca;
	box-shadow:0 10px 24px rgba(220,38,38,.42);
}
.blocked-phone-title-v1{
	margin:0;
	font-size:clamp(25px,5vw,34px);
	line-height:1.1;
	font-weight:950;
	color:#fff7f7;
	letter-spacing:0;
}
.blocked-phone-message-v1{
	max-width:430px;
	margin:11px auto 0;
	color:#fecaca;
	font-size:15.5px;
	line-height:1.6;
}
.blocked-phone-reason-v1{
	margin:20px 0 22px;
	padding:16px;
	border-radius:18px;
	text-align:left;
	background:rgba(14,6,12,.74);
	border:1px solid rgba(248,113,113,.26);
	box-shadow:inset 0 1px 0 rgba(255,255,255,.06);
}
.blocked-phone-reason-label-v1{
	margin:0 0 7px;
	color:#fca5a5;
	font-size:12px;
	font-weight:900;
	text-transform:uppercase;
	letter-spacing:.08em;
}
.blocked-phone-reason-text-v1{
	margin:0;
	color:#fff1f2;
	font-size:15px;
	line-height:1.6;
	overflow-wrap:anywhere;
}
.blocked-phone-back-v1{
	width:100%;
	min-height:46px;
	display:inline-flex;
	align-items:center;
	justify-content:center;
	border-radius:14px;
	background:rgba(255,255,255,.08)!important;
	border:1px solid rgba(254,202,202,.28)!important;
	color:#fff7f7!important;
	font-weight:850;
	box-shadow:none!important;
}
.blocked-phone-back-v1:hover{
	background:rgba(255,255,255,.12)!important;
	color:#fff!important;
}
@media(max-width:520px){
	.blocked-phone-screen-v1{
		min-height:calc(100vh - 110px);
		padding:18px 10px 34px;
	}
	.blocked-phone-card-v1{
		border-radius:22px;
		padding:24px 16px 18px;
	}
	.blocked-phone-icon-v1{
		width:76px;
		height:76px;
		border-radius:22px;
		margin-bottom:14px;
	}
	.blocked-phone-icon-v1 svg{
		width:48px;
		height:48px;
	}
	.blocked-phone-cross-v1{
		width:28px;
		height:28px;
		font-size:18px;
	}
	.blocked-phone-reason-v1{
		padding:14px;
		margin:17px 0 18px;
	}
}
</style>
<div class="blocked-phone-screen-v1">
	<section class="blocked-phone-card-v1" role="alert" aria-labelledby="blocked-phone-title-v1">
		<div class="blocked-phone-inner-v1">
			<div class="blocked-phone-icon-v1" aria-hidden="true">
				<svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
					<path d="M32 8c-13.2 0-22 9.3-22 21.2 0 7.3 3.3 12.9 8.3 16.4V53c0 2.2 1.8 4 4 4h19.4c2.2 0 4-1.8 4-4v-7.4c5-3.5 8.3-9.1 8.3-16.4C54 17.3 45.2 8 32 8Z" stroke="currentColor" stroke-width="3.4" stroke-linejoin="round"/>
					<path d="M25.2 33.7c3.2 0 5.8-2.3 5.8-5.1s-2.6-5.1-5.8-5.1-5.8 2.3-5.8 5.1 2.6 5.1 5.8 5.1ZM38.8 33.7c3.2 0 5.8-2.3 5.8-5.1s-2.6-5.1-5.8-5.1-5.8 2.3-5.8 5.1 2.6 5.1 5.8 5.1Z" fill="currentColor"/>
					<path d="M29 43h6M24 49v7M32 49v7M40 49v7" stroke="currentColor" stroke-width="3.4" stroke-linecap="round"/>
				</svg>
				<span class="blocked-phone-cross-v1">×</span>
			</div>
			<h1 id="blocked-phone-title-v1" class="blocked-phone-title-v1">Nomor WhatsApp Diblokir</h1>
			<p class="blocked-phone-message-v1">Maaf, nomor WhatsApp anda dibanned sementara.</p>
			<div class="blocked-phone-reason-v1">
				<p class="blocked-phone-reason-label-v1">Alasan</p>
				<p class="blocked-phone-reason-text-v1">{{.Reason}}</p>
			</div>
			<a class="btn blocked-phone-back-v1" href="{{.ReturnURL}}">Kembali</a>
		</div>
	</section>
</div>
`, map[string]any{"Reason": banReason, "ReturnURL": returnURL}))
		return
	}

	orderID := "VG" + time.Now().Format("060102150405") + strings.ToUpper(randCode(4))
	vuser := a.newVoucherCode()
	vpass := vuser

	res, err := a.db.Exec(
		`INSERT INTO transactions(order_id,router_id,package_id,customer_name,customer_phone,customer_email,hotspot_mac,hotspot_ip,amount,status,voucher_username,voucher_password)
		VALUES(?,?,?,?,?,?,?,?,?,'pending',?,?)`,
		orderID,
		ro.ID,
		pkgID,
		r.FormValue("customer_name"),
		customerPhone,
		r.FormValue("customer_email"),
		clientMAC,
		clientIP,
		amount,
		vuser,
		vpass,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	txID, _ := res.LastInsertId()

	slug := a.setting("pakasir_slug")
	if slug != "" {
		if qris, err := a.createPakasirQRIS(orderID, amount); err != nil {
			log.Printf("pakasir qris create gagal order=%s amount=%d err=%v", orderID, amount, err)
		} else {
			_, _ = a.db.Exec(`UPDATE transactions SET pakasir_txn_id=?, payment_method=?, qris_string=?, qris_total_payment=?, qris_expired_at=? WHERE id=?`, qris.TxnID, qris.PaymentMethod, qris.QRISString, qris.TotalPayment, qris.ExpiredAt, txID)
			go a.sendUnpaidPaymentWA(txID)
			go a.sendUnpaidPaymentEmail(txID)
		}
	}

	setPendingInvoiceCookie(w, r, orderID)

	http.Redirect(w, r, "/invoice/"+orderID, http.StatusSeeOther)
	if txID > 0 {
		go a.sendAdminOrderEmail(txID, orderID)
	} else {
		log.Printf("admin order email failed order_id=%s", orderID)
	}
}

func checkoutOrderURL(pkgID, clientMAC, clientIP string) string {
	q := url.Values{}
	q.Set("package_id", strings.TrimSpace(pkgID))
	if strings.TrimSpace(clientMAC) != "" {
		q.Set("mac", strings.TrimSpace(clientMAC))
	}
	if strings.TrimSpace(clientIP) != "" {
		q.Set("ip", strings.TrimSpace(clientIP))
	}
	return "/order?" + q.Encode()
}
