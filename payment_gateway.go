package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (a *App) paymentStart(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/pay-start/"))
	if orderID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	slug := strings.TrimSpace(a.setting("pakasir_slug"))
	apiKey := strings.TrimSpace(a.setting("pakasir_api_key"))
	if slug == "" || apiKey == "" {
		http.Redirect(w, r, "/invoice/"+url.PathEscape(orderID), http.StatusSeeOther)
		return
	}

	var ro Router
	var amount int64
	var status, mac, ip string
	var pakasirTxnID string

	err := a.db.QueryRow(`SELECT
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
			t.amount,
			t.status,
			COALESCE(t.hotspot_mac,''),
			COALESCE(t.hotspot_ip,''),
			COALESCE(t.pakasir_txn_id,'')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		WHERE t.order_id=?`, orderID).
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
			&status,
			&mac,
			&ip,
			&pakasirTxnID,
		)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	if status == "paid" {
		http.Redirect(w, r, "/invoice/"+url.PathEscape(orderID), http.StatusSeeOther)
		return
	}

	if strings.TrimSpace(pakasirTxnID) == "" {
		qris, err := a.createPakasirQRIS(orderID, amount)
		if err != nil {
			log.Printf("pakasir v2 create gagal order=%s amount=%d: %v", orderID, amount, err)
			http.Error(w, "gagal membuat transaksi pembayaran", http.StatusBadGateway)
			return
		}

		_, err = a.db.Exec(`
UPDATE transactions
SET pakasir_txn_id=?,
    payment_method=?,
    qris_string=?,
    qris_total_payment=?,
    qris_expired_at=?
WHERE order_id=? AND status<>'paid'`,
			qris.TxnID,
			qris.PaymentMethod,
			qris.QRISString,
			qris.TotalPayment,
			qris.ExpiredAt,
			orderID,
		)
		if err != nil {
			log.Printf("pakasir v2 simpan transaksi gagal order=%s txn=%s: %v", orderID, qris.TxnID, err)
			http.Error(w, "gagal menyimpan transaksi pembayaran", http.StatusInternalServerError)
			return
		}

		log.Printf("pakasir v2 create OK order=%s txn=%s amount=%d total=%d", orderID, qris.TxnID, amount, qris.TotalPayment)
	}

	mac = cleanPaymentClientMAC(mac)
	ip = cleanPaymentClientIP(ip)

	if mac != "" {
		go func(ro Router, orderID, mac, ip string) {
			time.Sleep(time.Duration(paymentBypassStartDelaySeconds) * time.Second)

			used, remaining, alreadyActive, allowed, err := a.reservePaymentBypassAttempt(orderID)
			if err != nil {
				log.Printf("payment delayed bypass cek gagal order=%s router=%s mac=%s ip=%s: %v", orderID, ro.Slug, mac, ip, err)
				return
			}
			if alreadyActive {
				log.Printf("payment delayed bypass skip masih aktif order=%s router=%s remaining=%ds used=%d/%d", orderID, ro.Slug, remaining, used, paymentBypassLimit)
				return
			}
			if !allowed {
				log.Printf("payment delayed bypass skip limit habis order=%s router=%s used=%d/%d", orderID, ro.Slug, used, paymentBypassLimit)
				return
			}

			if err := a.startPaymentBypass90(ro, orderID, mac, ip); err != nil {
				a.rollbackPaymentBypassAttempt(orderID, used)
				log.Printf("payment delayed bypass gagal order=%s router=%s mac=%s ip=%s: %v", orderID, ro.Slug, mac, ip, err)
				return
			}

			log.Printf("payment delayed bypass OK order=%s router=%s mac=%s ip=%s delay=%ds seconds=%d used=%d/%d", orderID, ro.Slug, mac, ip, paymentBypassStartDelaySeconds, paymentBypassSeconds, used, paymentBypassLimit)
		}(ro, orderID, mac, ip)
	} else {
		log.Printf("payment delayed bypass skip order=%s: hotspot_mac kosong", orderID)
	}

	http.Redirect(w, r, "/invoice/"+url.PathEscape(orderID), http.StatusSeeOther)
}

func (a *App) pakasirWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}

	expectedSecret := strings.TrimSpace(a.setting("pakasir_webhook_secret"))
	if expectedSecret == "" {
		http.Error(w, "pakasir webhook secret belum dikonfigurasi", http.StatusServiceUnavailable)
		return
	}

	gotSecret := strings.TrimSpace(r.Header.Get("X-Secret"))
	if gotSecret == "" || gotSecret != expectedSecret {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	raw, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))

	var h PakasirWebhookV2
	if err := json.Unmarshal(raw, &h); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	h.TxnID = strings.TrimSpace(h.TxnID)
	h.OrderID = strings.TrimSpace(h.OrderID)
	h.Status = strings.TrimSpace(h.Status)

	if h.TxnID == "" || h.OrderID == "" || h.Amount <= 0 {
		http.Error(w, "payload tidak lengkap", http.StatusBadRequest)
		return
	}

	if h.Status != "completed" {
		http.Error(w, "transaction belum completed", http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(h.OrderID, "BILL-") {
		if !a.handleClientBillingPakasirWebhookV2(h, raw) {
			http.Error(w, "billing invoice not valid", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"type":"client_billing"}`))
		return
	}

	var txID int64
	var amount int64
	var pakasirTxnID string
	var paymentMethod string
	var oldStatus string

	err := a.db.QueryRow(`
SELECT id,
       amount,
       COALESCE(pakasir_txn_id,''),
       COALESCE(payment_method,''),
       status
FROM transactions
WHERE order_id=?`,
		h.OrderID,
	).Scan(&txID, &amount, &pakasirTxnID, &paymentMethod, &oldStatus)

	if err != nil {
		http.Error(w, "tx not found", http.StatusNotFound)
		return
	}

	if pakasirTxnID == "" || pakasirTxnID != h.TxnID || amount != h.Amount {
		_, _ = a.db.Exec(`UPDATE transactions SET raw_webhook=? WHERE id=?`, string(raw), txID)
		http.Error(w, "transaction mismatch", http.StatusBadRequest)
		return
	}

	if oldStatus == "paid" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"duplicate":true}`))
		return
	}

	if !a.verifyPakasirV2(h.TxnID, h.OrderID, amount) {
		_, _ = a.db.Exec(`UPDATE transactions SET raw_webhook=? WHERE id=?`, string(raw), txID)
		http.Error(w, "transaction not valid", http.StatusBadRequest)
		return
	}

	res, err := a.db.Exec(`
UPDATE transactions
SET status='paid',
    payment_method=?,
    raw_webhook=?,
    paid_at=IFNULL(paid_at,NOW())
WHERE id=? AND status<>'paid'`,
		paymentMethod,
		string(raw),
		txID,
	)
	if err != nil {
		http.Error(w, "gagal update transaction", http.StatusInternalServerError)
		return
	}

	if n, _ := res.RowsAffected(); n == 0 {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"duplicate":true}`))
		return
	}

	if err := a.provisionVoucher(txID); err == nil {
		go a.sendVoucherEmail(txID)
	}

	go a.sendWanesia(txID)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true}`))
}

func (a *App) verifyPakasirV2(txnID, orderID string, amount int64) bool {
	project := strings.TrimSpace(a.setting("pakasir_slug"))
	apiKey := strings.TrimSpace(a.setting("pakasir_api_key"))
	txnID = strings.TrimSpace(txnID)
	if project == "" || apiKey == "" || txnID == "" || orderID == "" || amount <= 0 {
		return false
	}

	endpoint := "https://app.pakasir.com/api/v2/transaction-status/" + url.PathEscape(project) + "/" + url.PathEscape(txnID)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}

	var d PakasirV2Status
	if json.NewDecoder(resp.Body).Decode(&d) != nil {
		return false
	}

	return d.TxnID == txnID &&
		d.OrderID == orderID &&
		d.Amount == amount &&
		d.Status == "completed"
}

func (a *App) paymentGatewayAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		provider := strings.TrimSpace(r.FormValue("payment_provider"))
		if provider != "midtrans" {
			provider = "pakasir"
		}

		_ = a.setSetting("payment_provider", provider)
		http.Redirect(w, r, "/admin/payment-gateway", http.StatusSeeOther)
		return
	}

	provider := strings.TrimSpace(a.setting("payment_provider"))
	if provider != "midtrans" {
		provider = "pakasir"
	}

	vals := map[string]any{
		"Provider": provider,
	}

	if provider == "midtrans" {
		vals["PakasirChecked"] = ""
		vals["MidtransChecked"] = "checked"
		vals["ActiveName"] = "Midtrans"
		vals["ActiveDesc"] = "Midtrans dipilih sebagai provider aktif."
	} else {
		vals["PakasirChecked"] = "checked"
		vals["MidtransChecked"] = ""
		vals["ActiveName"] = "Pakasir"
		vals["ActiveDesc"] = "Pakasir dipilih sebagai provider aktif."
	}

	body := renderHTML(`
<style>
.pg-home{display:grid;gap:18px}

/* PAYMENT_GATEWAY_LIGHT_POLISH_V1 */
html[data-theme="light"] .pg-status{
	background:#eff6ff!important;
	border-color:#bfdbfe!important;
	color:#1e3a8a!important;
	box-shadow:none!important;
}
html[data-theme="light"] .pg-status strong{
	color:#0f172a!important;
}
html[data-theme="light"] .pg-select-card,
html[data-theme="light"] .pg-config-card{
	background:#ffffff!important;
	border-color:#dbe4f0!important;
	color:#0f172a!important;
	box-shadow:0 10px 28px rgba(15,23,42,.06)!important;
}
html[data-theme="light"] .pg-select-card h3,
html[data-theme="light"] .pg-config-card h3{
	color:#0f172a!important;
}
html[data-theme="light"] .pg-select-card p,
html[data-theme="light"] .pg-config-card p{
	color:#475569!important;
}
html[data-theme="light"] .pg-select-card:has(input:checked){
	background:linear-gradient(135deg,#eff6ff,#ffffff)!important;
	border-color:#3b82f6!important;
	box-shadow:0 18px 40px rgba(37,99,235,.14)!important;
}
html[data-theme="light"] .pg-select-card:not(:has(input:checked)){
	opacity:1!important;
}
html[data-theme="light"] .pg-config-card .btn{
	background:#0f172a!important;
	color:#ffffff!important;
	border-color:#0f172a!important;
}
html[data-theme="light"] .pg-config-card .btn:hover{
	background:#2563eb!important;
	border-color:#2563eb!important;
}
html[data-theme="light"] .settings-page-card .muted{
	color:#475569!important;
}

.pg-top{display:flex;justify-content:space-between;gap:16px;align-items:flex-start;flex-wrap:wrap}
.pg-status{border:1px solid rgba(96,165,250,.32);background:rgba(37,99,235,.12);padding:12px 14px;border-radius:18px;color:#bfdbfe;min-width:220px}
.pg-status strong{display:block;color:white;margin-bottom:4px}
.pg-select-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px}
.pg-select-card{position:relative;display:block;border:1px solid rgba(148,163,184,.20);background:rgba(15,23,42,.50);border-radius:22px;padding:18px;cursor:pointer;transition:.15s ease}
.pg-select-card:hover{border-color:rgba(96,165,250,.55);transform:translateY(-1px)}
.pg-select-card input{position:absolute;opacity:0;pointer-events:none}
.pg-select-card h3{margin:0 0 8px}
.pg-select-card p{margin:0;color:#9fb6d8;line-height:1.45}
.pg-select-card:has(input:checked){border-color:rgba(59,130,246,.90);background:linear-gradient(135deg,rgba(37,99,235,.22),rgba(15,23,42,.64));box-shadow:0 18px 45px rgba(37,99,235,.16)}
.pg-select-card:has(input:checked)::after{content:"Aktif";position:absolute;top:14px;right:14px;background:#2563eb;color:white;border-radius:999px;padding:5px 10px;font-size:12px;font-weight:900}
.pg-config-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px}
.pg-config-card{border:1px solid rgba(148,163,184,.18);background:rgba(15,23,42,.35);border-radius:22px;padding:18px;display:flex;flex-direction:column;gap:12px}
.pg-config-card h3{margin:0}
.pg-config-card p{margin:0;color:#9fb6d8;line-height:1.45}
.pg-config-card .btn{align-self:flex-start;margin-top:auto;text-decoration:none}
@media(max-width:760px){.pg-select-grid,.pg-config-grid{grid-template-columns:1fr}.pg-top{display:block}.pg-status{margin-top:12px}}

.inv-top-grid{
	display:grid;
	grid-template-columns:repeat(2,minmax(0,1fr));
	gap:12px;
	margin:12px 0 14px;
}
.inv-top-grid > div{
	margin:0 !important;
	min-width:0;
}
.inv-top-grid .inv-span-2{
	grid-column:1 / -1;
}
@media (max-width:640px){
	.inv-top-grid{
		grid-template-columns:repeat(2,minmax(0,1fr));
		gap:10px;
	}
	.inv-top-grid > div{
		padding:16px !important;
	}
	.qris-download-note{
		font-size:12px;
	}
}

</style>

<div class="card settings-page-card pg-home">
	<div class="pg-top">
		<div>
			<h2>Payment Gateway</h2>
			<p class="muted">Pilih provider aktif di sini. Konfigurasi tiap provider dibuat di halaman terpisah agar lebih mudah dipahami client.</p>
		</div>
		<div class="pg-status">
			<strong>Aktif: {{.ActiveName}}</strong>
			<span>{{.ActiveDesc}}</span>
		</div>
	</div>

	<form method="post" class="settings-form-compact-v1">
		<div class="pg-select-grid">
			<label class="pg-select-card">
				<input type="radio" name="payment_provider" value="pakasir" {{.PakasirChecked}}>
				<h3>Pakasir</h3>
				<p>Flow lama tetap berjalan. Cocok untuk gateway sederhana berbasis redirect.</p>
			</label>

			<label class="pg-select-card">
				<input type="radio" name="payment_provider" value="midtrans" {{.MidtransChecked}}>
				<h3>Midtrans</h3>
				<p>Untuk Snap/redirect Midtrans. Aktivasi checkout dibuat setelah konfigurasi siap.</p>
			</label>
		</div>

		<div class="actions settings-actions">
			<button class="btn pri" type="submit">Simpan Provider Aktif</button>
		</div>
	</form>

	<div class="pg-config-grid">
		<div class="pg-config-card">
			<h3>Setting Pakasir</h3>
			<p>Atur project/slug, API key, dan opsi QRIS Pakasir.</p>
			<a class="btn" href="/admin/payment-gateway/pakasir">Buka Setting Pakasir</a>
		</div>

		<div class="pg-config-card">
			<h3>Setting Midtrans</h3>
			<p>Atur mode sandbox/production, server key, client key, dan merchant ID.</p>
			<a class="btn" href="/admin/payment-gateway/midtrans">Buka Setting Midtrans</a>
		</div>
	</div>
</div>`, vals)

	page(w, "Payment Gateway", body)
}

func (a *App) paymentGatewayPakasirAdmin(w http.ResponseWriter, r *http.Request) {
	keys := []string{
		"pakasir_slug",
		"pakasir_api_key",
		"pakasir_qris_only",
	}

	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		for _, k := range keys {
			_ = a.setSetting(k, r.FormValue(k))
		}
		if r.FormValue("make_active") == "1" {
			_ = a.setSetting("payment_provider", "pakasir")
		}
		http.Redirect(w, r, "/admin/payment-gateway/pakasir", http.StatusSeeOther)
		return
	}

	vals := map[string]any{}
	for _, k := range keys {
		vals[k] = a.setting(k)
	}
	vals["IsActive"] = a.setting("payment_provider") != "midtrans"

	body := renderHTML(`



<style>
/* PAYMENT_PROVIDER_FORM_CONTROLS_V2 */

/* PAYMENT_PROVIDER_FORM_FINAL_POLISH_V1 */
html[data-theme="light"] .pg-provider-form .pg-radio-row label{
	color:#0f172a!important;
}
html[data-theme="light"] .pg-provider-form .pg-radio-row input[type="radio"]{
	border:1px solid #94a3b8!important;
}
html[data-theme="light"] .pg-provider-form .pg-check-row{
	background:#f8fafc!important;
	border-color:#cbd5e1!important;
	color:#0f172a!important;
	box-shadow:0 8px 22px rgba(15,23,42,.06)!important;
}
html[data-theme="light"] .pg-provider-form .pg-check-row:hover{
	border-color:#60a5fa!important;
	background:#eff6ff!important;
}
html[data-theme="light"] .pg-provider-form .pg-check-row input[type="checkbox"]{
	border:1px solid #94a3b8!important;
}
html[data-theme="light"] .pg-provider-form .field label{
	color:#334155!important;
}
html[data-theme="light"] .pg-provider-form input{
	color:#0f172a!important;
	background:#ffffff!important;
	border-color:#dbe4f0!important;
}
html[data-theme="light"] .pg-provider-form .muted{
	color:#475569!important;
}
html[data-theme="light"] .pg-provider-form b{
	color:#0f172a!important;
}
html[data-theme="dark"] .pg-provider-form .pg-check-row{
	background:rgba(37,99,235,.10)!important;
	border-color:rgba(96,165,250,.28)!important;
	color:#dbeafe!important;
}
html[data-theme="dark"] .pg-provider-form .pg-check-row:hover{
	border-color:rgba(96,165,250,.55)!important;
	background:rgba(37,99,235,.18)!important;
}
.pg-provider-form .pg-check-row{
	cursor:pointer!important;
	transition:.15s ease!important;
}
.pg-provider-form .pg-check-row:hover{
	transform:translateY(-1px)!important;
}
.pg-provider-form .actions.settings-actions{
	margin-top:18px!important;
	display:flex!important;
	gap:10px!important;
	flex-wrap:wrap!important;
	align-items:center!important;
}
@media(max-width:760px){
	.pg-provider-form .pg-check-row{
		width:100%!important;
	}
	.pg-provider-form .actions.settings-actions .btn{
		width:100%!important;
		text-align:center!important;
	}
}

.pg-provider-form .pg-radio-row{
	display:flex!important;
	flex-direction:row!important;
	gap:18px!important;
	align-items:center!important;
	flex-wrap:wrap!important;
	padding:10px 0 0!important;
}
.pg-provider-form .pg-radio-row label,
.pg-provider-form .pg-check-row{
	display:inline-flex!important;
	flex-direction:row!important;
	align-items:center!important;
	justify-content:flex-start!important;
	gap:10px!important;
	width:auto!important;
	max-width:100%!important;
	margin:0!important;
	color:#dbeafe!important;
	font-weight:800!important;
	line-height:1.35!important;
}
.pg-provider-form .pg-check-row{
	margin-top:16px!important;
	padding:12px 14px!important;
	border:1px solid rgba(148,163,184,.18)!important;
	background:rgba(15,23,42,.35)!important;
	border-radius:16px!important;
}
.pg-provider-form input[type="radio"],
.pg-provider-form input[type="checkbox"]{
	-webkit-appearance:auto!important;
	appearance:auto!important;
	display:inline-block!important;
	position:static!important;
	opacity:1!important;
	pointer-events:auto!important;
	width:17px!important;
	height:17px!important;
	min-width:17px!important;
	min-height:17px!important;
	max-width:17px!important;
	max-height:17px!important;
	padding:0!important;
	margin:0!important;
	border-radius:initial!important;
	box-shadow:none!important;
	accent-color:#2563eb!important;
	flex:0 0 17px!important;
}
.pg-provider-form input[type="checkbox"]{
	width:18px!important;
	height:18px!important;
	min-width:18px!important;
	min-height:18px!important;
	max-width:18px!important;
	max-height:18px!important;
	flex-basis:18px!important;
}
.pg-provider-form .pg-check-row span{
	display:inline!important;
}
</style>


<style>
/* PAYMENT_SECRET_EYE_FORCE_V2 */
.pg-provider-form .secret-field{
	position:relative!important;
}
.pg-provider-form .secret-field input{
	padding-right:46px!important;
}
.pg-provider-form .secret-field > .pass-eye{
	position:absolute!important;
	right:10px!important;
	top:auto!important;
	bottom:10px!important;
	transform:none!important;
	width:28px!important;
	height:28px!important;
	min-width:28px!important;
	min-height:28px!important;
	max-width:28px!important;
	max-height:28px!important;
	padding:0!important;
	margin:0!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	border-radius:9px!important;
	font-size:12px!important;
	line-height:1!important;
	z-index:20!important;
	cursor:pointer!important;
	background:rgba(15,23,42,.12)!important;
	border:1px solid rgba(148,163,184,.24)!important;
	color:#0f172a!important;
	box-shadow:none!important;
}
.pg-provider-form .secret-field > .pass-eye:hover{
	background:#eff6ff!important;
	border-color:#93c5fd!important;
}
html[data-theme="dark"] .pg-provider-form .secret-field > .pass-eye{
	background:rgba(15,23,42,.34)!important;
	border-color:rgba(148,163,184,.24)!important;
	color:#e5edf8!important;
}
html[data-theme="dark"] .pg-provider-form .secret-field > .pass-eye:hover{
	background:rgba(37,99,235,.20)!important;
	border-color:rgba(96,165,250,.45)!important;
}
@media(max-width:760px){
	.pg-provider-form .secret-field input{
		padding-right:44px!important;
	}
	.pg-provider-form .secret-field > .pass-eye{
		width:26px!important;
		height:26px!important;
		min-width:26px!important;
		min-height:26px!important;
		max-width:26px!important;
		max-height:26px!important;
		right:9px!important;
		bottom:11px!important;
		font-size:11px!important;
	}
}
</style>

<script>
// PAYMENT_SECRET_TOGGLE_FORCE_V2
(function(){
	window.toggleSecret = function(btn){
		if(!btn) return;

		var wrap = btn.closest ? btn.closest(".secret-field") : null;
		var input = wrap ? wrap.querySelector("input") : null;

		if(!input && btn.parentElement){
			input = btn.parentElement.querySelector("input");
		}
		if(!input) return;

		var show = input.type === "password";
		input.type = show ? "text" : "password";
		btn.textContent = show ? "🙈" : "👁";
		btn.setAttribute("aria-label", show ? "Sembunyikan" : "Tampilkan");
		btn.setAttribute("title", show ? "Sembunyikan" : "Tampilkan");
	};
})();
</script>

<div class="card settings-page-card pg-provider-form">
	<h2>Setting Pakasir</h2>
	<p class="muted">Konfigurasi khusus Pakasir. Tidak ada field Midtrans di halaman ini.</p>

	<form method="post" class="settings-form-compact-v1">
		<div class="grid">
			<div class="field"><label>Pakasir Slug / Project</label><input name="pakasir_slug" value="{{index . "pakasir_slug"}}"></div>
			<div class="field secret-field"><label>Pakasir API Key</label><input name="pakasir_api_key" type="password" value="{{index . "pakasir_api_key"}}"><button type="button" class="pass-eye" onclick="toggleSecret(this)">👁</button></div>
			<div class="field secret-field"><label>Pakasir Webhook Secret</label><input name="pakasir_webhook_secret" type="password" value="{{index . "pakasir_webhook_secret"}}"><button type="button" class="pass-eye" onclick="toggleSecret(this)">👁</button></div>
			<div class="field"><label>Pakasir QRIS Only 1/0</label><input name="pakasir_qris_only" value="{{index . "pakasir_qris_only"}}"></div>
		</div>

		<label class="pg-check-row">
			<input type="checkbox" name="make_active" value="1">
			<span>Jadikan Pakasir sebagai provider aktif</span>
		</label>

		<div class="actions settings-actions">
			<button class="btn pri" type="submit">Simpan Setting Pakasir</button>
			<a class="btn" href="/admin/payment-gateway">Kembali</a>
		</div>
	</form>

	<p class="muted">Webhook Pakasir: <b>/webhook/pakasir</b></p>
</div>`, vals)

	page(w, "Setting Pakasir", body)
}

func (a *App) paymentGatewayMidtransAdmin(w http.ResponseWriter, r *http.Request) {
	keys := []string{
		"midtrans_mode",
		"midtrans_server_key",
		"midtrans_client_key",
		"midtrans_merchant_id",
	}

	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		mode := strings.TrimSpace(r.FormValue("midtrans_mode"))
		if mode != "production" {
			mode = "sandbox"
		}

		_ = a.setSetting("midtrans_mode", mode)
		_ = a.setSetting("midtrans_server_key", r.FormValue("midtrans_server_key"))
		_ = a.setSetting("midtrans_client_key", r.FormValue("midtrans_client_key"))
		_ = a.setSetting("midtrans_merchant_id", r.FormValue("midtrans_merchant_id"))

		if r.FormValue("make_active") == "1" {
			_ = a.setSetting("payment_provider", "midtrans")
		}

		http.Redirect(w, r, "/admin/payment-gateway/midtrans", http.StatusSeeOther)
		return
	}

	vals := map[string]any{}
	for _, k := range keys {
		vals[k] = a.setting(k)
	}

	mode := strings.TrimSpace(a.setting("midtrans_mode"))
	if mode == "production" {
		vals["SandboxChecked"] = ""
		vals["ProductionChecked"] = "checked"
	} else {
		vals["SandboxChecked"] = "checked"
		vals["ProductionChecked"] = ""
	}

	body := renderHTML(`

<style>
/* PAYMENT_PROVIDER_FORM_CONTROLS_V2 */
.pg-provider-form .pg-radio-row{
	display:flex!important;
	flex-direction:row!important;
	gap:18px!important;
	align-items:center!important;
	flex-wrap:wrap!important;
	padding:10px 0 0!important;
}
.pg-provider-form .pg-radio-row label,
.pg-provider-form .pg-check-row{
	display:inline-flex!important;
	flex-direction:row!important;
	align-items:center!important;
	justify-content:flex-start!important;
	gap:10px!important;
	width:auto!important;
	max-width:100%!important;
	margin:0!important;
	color:#dbeafe!important;
	font-weight:800!important;
	line-height:1.35!important;
}
.pg-provider-form .pg-check-row{
	margin-top:16px!important;
	padding:12px 14px!important;
	border:1px solid rgba(148,163,184,.18)!important;
	background:rgba(15,23,42,.35)!important;
	border-radius:16px!important;
}
.pg-provider-form input[type="radio"],
.pg-provider-form input[type="checkbox"]{
	-webkit-appearance:auto!important;
	appearance:auto!important;
	display:inline-block!important;
	position:static!important;
	opacity:1!important;
	pointer-events:auto!important;
	width:17px!important;
	height:17px!important;
	min-width:17px!important;
	min-height:17px!important;
	max-width:17px!important;
	max-height:17px!important;
	padding:0!important;
	margin:0!important;
	border-radius:initial!important;
	box-shadow:none!important;
	accent-color:#2563eb!important;
	flex:0 0 17px!important;
}
.pg-provider-form input[type="checkbox"]{
	width:18px!important;
	height:18px!important;
	min-width:18px!important;
	min-height:18px!important;
	max-width:18px!important;
	max-height:18px!important;
	flex-basis:18px!important;
}
.pg-provider-form .pg-check-row span{
	display:inline!important;
}
</style>


<style>
/* PAYMENT_MIDTRANS_FORM_POLISH_V1 */
.pg-midtrans-form .pg-radio-row{
	display:flex!important;
	flex-direction:row!important;
	gap:18px!important;
	align-items:center!important;
	flex-wrap:wrap!important;
	padding:10px 0 0!important;
}
.pg-midtrans-form .pg-radio-row label,
.pg-midtrans-form .pg-check-row{
	display:inline-flex!important;
	flex-direction:row!important;
	align-items:center!important;
	justify-content:flex-start!important;
	gap:10px!important;
	width:auto!important;
	max-width:100%!important;
	margin:0!important;
	font-weight:800!important;
	line-height:1.35!important;
}
.pg-midtrans-form .pg-check-row{
	margin-top:16px!important;
	padding:12px 14px!important;
	border-radius:16px!important;
	cursor:pointer!important;
	transition:.15s ease!important;
}
.pg-midtrans-form .pg-check-row:hover{
	transform:translateY(-1px)!important;
}
.pg-midtrans-form input[type="radio"],
.pg-midtrans-form input[type="checkbox"]{
	-webkit-appearance:auto!important;
	appearance:auto!important;
	display:inline-block!important;
	position:static!important;
	opacity:1!important;
	pointer-events:auto!important;
	width:17px!important;
	height:17px!important;
	min-width:17px!important;
	min-height:17px!important;
	max-width:17px!important;
	max-height:17px!important;
	padding:0!important;
	margin:0!important;
	box-shadow:none!important;
	accent-color:#2563eb!important;
	flex:0 0 17px!important;
}
.pg-midtrans-form input[type="checkbox"]{
	width:18px!important;
	height:18px!important;
	min-width:18px!important;
	min-height:18px!important;
	max-width:18px!important;
	max-height:18px!important;
	flex-basis:18px!important;
}
html[data-theme="light"] .pg-midtrans-form .pg-radio-row label,
html[data-theme="light"] .pg-midtrans-form .pg-check-row{
	color:#0f172a!important;
}
html[data-theme="light"] .pg-midtrans-form .pg-check-row{
	background:#f8fafc!important;
	border:1px solid #cbd5e1!important;
	box-shadow:0 8px 22px rgba(15,23,42,.06)!important;
}
html[data-theme="light"] .pg-midtrans-form .pg-check-row:hover{
	background:#eff6ff!important;
	border-color:#60a5fa!important;
}
html[data-theme="light"] .pg-midtrans-form .field label{
	color:#334155!important;
}
html[data-theme="light"] .pg-midtrans-form input{
	color:#0f172a!important;
	background:#ffffff!important;
	border-color:#dbe4f0!important;
}
html[data-theme="light"] .pg-midtrans-form .muted,
html[data-theme="light"] .pg-midtrans-form p{
	color:#475569!important;
}
html[data-theme="light"] .pg-midtrans-form b{
	color:#0f172a!important;
}
html[data-theme="dark"] .pg-midtrans-form .pg-check-row{
	background:rgba(37,99,235,.10)!important;
	border:1px solid rgba(96,165,250,.28)!important;
	color:#dbeafe!important;
}
html[data-theme="dark"] .pg-midtrans-form .pg-check-row:hover{
	border-color:rgba(96,165,250,.55)!important;
	background:rgba(37,99,235,.18)!important;
}
@media(max-width:760px){
	.pg-midtrans-form .pg-check-row{
		width:100%!important;
	}
	.pg-midtrans-form .actions.settings-actions .btn{
		width:100%!important;
		text-align:center!important;
	}
}
</style>


<style>
/* PAYMENT_SECRET_EYE_FORCE_V2 */
.pg-provider-form .secret-field{
	position:relative!important;
}
.pg-provider-form .secret-field input{
	padding-right:46px!important;
}
.pg-provider-form .secret-field > .pass-eye{
	position:absolute!important;
	right:10px!important;
	top:auto!important;
	bottom:10px!important;
	transform:none!important;
	width:28px!important;
	height:28px!important;
	min-width:28px!important;
	min-height:28px!important;
	max-width:28px!important;
	max-height:28px!important;
	padding:0!important;
	margin:0!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	border-radius:9px!important;
	font-size:12px!important;
	line-height:1!important;
	z-index:20!important;
	cursor:pointer!important;
	background:rgba(15,23,42,.12)!important;
	border:1px solid rgba(148,163,184,.24)!important;
	color:#0f172a!important;
	box-shadow:none!important;
}
.pg-provider-form .secret-field > .pass-eye:hover{
	background:#eff6ff!important;
	border-color:#93c5fd!important;
}
html[data-theme="dark"] .pg-provider-form .secret-field > .pass-eye{
	background:rgba(15,23,42,.34)!important;
	border-color:rgba(148,163,184,.24)!important;
	color:#e5edf8!important;
}
html[data-theme="dark"] .pg-provider-form .secret-field > .pass-eye:hover{
	background:rgba(37,99,235,.20)!important;
	border-color:rgba(96,165,250,.45)!important;
}
@media(max-width:760px){
	.pg-provider-form .secret-field input{
		padding-right:44px!important;
	}
	.pg-provider-form .secret-field > .pass-eye{
		width:26px!important;
		height:26px!important;
		min-width:26px!important;
		min-height:26px!important;
		max-width:26px!important;
		max-height:26px!important;
		right:9px!important;
		bottom:11px!important;
		font-size:11px!important;
	}
}
</style>

<script>
// PAYMENT_SECRET_TOGGLE_FORCE_V2
(function(){
	window.toggleSecret = function(btn){
		if(!btn) return;

		var wrap = btn.closest ? btn.closest(".secret-field") : null;
		var input = wrap ? wrap.querySelector("input") : null;

		if(!input && btn.parentElement){
			input = btn.parentElement.querySelector("input");
		}
		if(!input) return;

		var show = input.type === "password";
		input.type = show ? "text" : "password";
		btn.textContent = show ? "🙈" : "👁";
		btn.setAttribute("aria-label", show ? "Sembunyikan" : "Tampilkan");
		btn.setAttribute("title", show ? "Sembunyikan" : "Tampilkan");
	};
})();
</script>

<div class="card settings-page-card pg-provider-form pg-midtrans-form">
	<h2>Setting Midtrans</h2>
	<p class="muted">Konfigurasi khusus Midtrans. Tidak ada field Pakasir di halaman ini.</p>

	<form method="post" class="settings-form-compact-v1">
		<div class="grid">
			<div class="field">
				<label>Mode</label>
				<div class="pg-radio-row">
					<label><input type="radio" name="midtrans_mode" value="sandbox" {{.SandboxChecked}}> Sandbox</label>
					<label><input type="radio" name="midtrans_mode" value="production" {{.ProductionChecked}}> Production</label>
				</div>
			</div>
			<div class="field"><label>Midtrans Merchant ID</label><input name="midtrans_merchant_id" value="{{index . "midtrans_merchant_id"}}"></div>
			<div class="field secret-field"><label>Midtrans Server Key</label><input name="midtrans_server_key" type="password" value="{{index . "midtrans_server_key"}}"><button type="button" class="pass-eye" onclick="toggleSecret(this)">👁</button></div>
			<div class="field secret-field"><label>Midtrans Client Key</label><input name="midtrans_client_key" type="password" value="{{index . "midtrans_client_key"}}"><button type="button" class="pass-eye" onclick="toggleSecret(this)">👁</button></div>
		</div>

		<label class="pg-check-row">
			<input type="checkbox" name="make_active" value="1">
			<span>Jadikan Midtrans sebagai provider aktif</span>
		</label>

		<div class="actions settings-actions">
			<button class="btn pri" type="submit">Simpan Setting Midtrans</button>
			<a class="btn" href="/admin/payment-gateway">Kembali</a>
		</div>
	</form>

	<p class="muted">Webhook Midtrans: <b>/webhook/midtrans</b></p>
	<p class="muted">Catatan: checkout Midtrans akan diaktifkan pada tahap berikutnya.</p>
</div>`, vals)

	page(w, "Setting Midtrans", body)
}
