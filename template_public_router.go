package main

const publicRouterPageTemplate = `
<!-- PUBLIC_ROUTER_MOBILE_POLISH_V1 -->
<style>
.public-router-v1{
	display:grid;
	gap:14px;
}

.public-router-hero-v1{
	position:relative;
	overflow:hidden;
	border-radius:28px;
	padding:22px 20px;
	background:
		radial-gradient(circle at 20% 10%, rgba(56,189,248,.28), transparent 34%),
		radial-gradient(circle at 88% 22%, rgba(167,139,250,.24), transparent 34%),
		linear-gradient(135deg, rgba(15,23,42,.88), rgba(30,41,59,.72));
	border:1px solid rgba(148,163,184,.18);
	box-shadow:0 18px 45px rgba(0,0,0,.24);
}

.public-router-kicker-v1{
	display:inline-flex;
	align-items:center;
	gap:7px;
	padding:7px 10px;
	border-radius:999px;
	background:rgba(34,197,94,.12);
	border:1px solid rgba(34,197,94,.20);
	color:#86efac;
	font-size:12px;
	font-weight:900;
	margin-bottom:12px;
}

.public-router-hero-v1 h1{
	margin:0;
	font-size:31px;
	line-height:1.05;
	letter-spacing:-1px;
	color:#f8fafc;
}

.public-router-hero-v1 p{
	margin:9px 0 0;
	color:#cbd5e1;
	font-size:14px;
	line-height:1.5;
	max-width:560px;
}

.public-package-grid-v1{
	display:grid;
	grid-template-columns:repeat(2,minmax(0,1fr));
	gap:12px;
}

.public-package-card-v1{
	position:relative;
	overflow:hidden;
	border-radius:24px;
	padding:16px;
	background:rgba(15,23,42,.74);
	border:1px solid rgba(148,163,184,.18);
	box-shadow:0 14px 36px rgba(0,0,0,.20);
}

.public-package-card-v1::before{
	content:"";
	position:absolute;
	width:90px;
	height:90px;
	right:-34px;
	top:-34px;
	border-radius:999px;
	background:rgba(59,130,246,.18);
	pointer-events:none;
}

.public-package-name-v1{
	position:relative;
	z-index:1;
	margin:0 0 8px;
	font-size:17px;
	line-height:1.15;
	font-weight:950;
	color:#f8fafc;
}

.public-package-price-v1{
	position:relative;
	z-index:1;
	font-size:26px;
	line-height:1;
	font-weight:950;
	letter-spacing:-.8px;
	color:#fff;
	margin-bottom:10px;
}

.public-package-meta-v1{
	position:relative;
	z-index:1;
	display:flex;
	flex-wrap:wrap;
	gap:6px;
	margin-bottom:13px;
}

.public-package-chip-v1{
	display:inline-flex;
	align-items:center;
	min-height:24px;
	padding:0 8px;
	border-radius:999px;
	background:rgba(148,163,184,.13);
	border:1px solid rgba(148,163,184,.16);
	color:#cbd5e1;
	font-size:11px;
	font-weight:800;
}

.public-package-card-v1 .buy-btn{
	position:relative;
	z-index:1;
	width:100%!important;
	min-height:43px;
	border-radius:15px!important;
	font-weight:950!important;
}

.public-empty-v1{
	border-radius:22px;
	padding:18px;
	background:rgba(15,23,42,.72);
	border:1px solid rgba(148,163,184,.18);
	color:#cbd5e1;
	font-weight:800;
}

html[data-theme="light"] .public-router-hero-v1,
html[data-admin-theme="light"] .public-router-hero-v1{
	background:
		radial-gradient(circle at 20% 10%, rgba(14,165,233,.18), transparent 34%),
		radial-gradient(circle at 88% 22%, rgba(99,102,241,.14), transparent 34%),
		linear-gradient(135deg, rgba(255,255,255,.96), rgba(241,245,249,.92))!important;
	border-color:#dbe4f0!important;
	box-shadow:0 16px 36px rgba(15,23,42,.10)!important;
}

html[data-theme="light"] .public-router-hero-v1 h1,
html[data-admin-theme="light"] .public-router-hero-v1 h1,
html[data-theme="light"] .public-package-name-v1,
html[data-admin-theme="light"] .public-package-name-v1,
html[data-theme="light"] .public-package-price-v1,
html[data-admin-theme="light"] .public-package-price-v1{
	color:#0f172a!important;
}

html[data-theme="light"] .public-router-hero-v1 p,
html[data-admin-theme="light"] .public-router-hero-v1 p{
	color:#475569!important;
}

html[data-theme="light"] .public-package-card-v1,
html[data-admin-theme="light"] .public-package-card-v1,
html[data-theme="light"] .public-empty-v1,
html[data-admin-theme="light"] .public-empty-v1{
	background:rgba(255,255,255,.94)!important;
	border-color:#dbe4f0!important;
	box-shadow:0 14px 34px rgba(15,23,42,.09)!important;
	color:#475569!important;
}

html[data-theme="light"] .public-package-chip-v1,
html[data-admin-theme="light"] .public-package-chip-v1{
	background:#f1f5f9!important;
	border-color:#e2e8f0!important;
	color:#475569!important;
}

@media(max-width:720px){
	.public-router-v1{
		gap:10px;
	}

	.public-router-hero-v1{
		border-radius:22px;
		padding:17px 15px;
	}

	.public-router-kicker-v1{
		font-size:11px;
		padding:6px 9px;
		margin-bottom:10px;
	}

	.public-router-hero-v1 h1{
		font-size:25px;
	}

	.public-router-hero-v1 p{
		font-size:13px;
		line-height:1.4;
		margin-top:7px;
	}

	.public-package-grid-v1{
		grid-template-columns:repeat(2,minmax(0,1fr));
		gap:9px;
	}

	.public-package-card-v1{
		border-radius:19px;
		padding:12px;
	}

	.public-package-name-v1{
		font-size:14px;
		margin-bottom:7px;
	}

	.public-package-price-v1{
		font-size:21px;
		margin-bottom:9px;
	}

	.public-package-meta-v1{
		gap:5px;
		margin-bottom:10px;
	}

	.public-package-chip-v1{
		min-height:21px;
		padding:0 7px;
		font-size:10px;
	}

	.public-package-card-v1 .buy-btn{
		min-height:38px;
		border-radius:13px!important;
		font-size:13px!important;
		padding:9px 10px!important;
	}
}

@media(max-width:360px){
	.public-package-grid-v1{
		gap:7px;
	}

	.public-package-card-v1{
		padding:10px;
		border-radius:17px;
	}

	.public-package-price-v1{
		font-size:19px;
	}
}
</style>

<div class="public-router-v1">
	<div class="public-router-hero-v1">
		<!-- PUBLIC_ROUTER_HIDE_DUPLICATE_H1_V1 -->
		<div class="public-notice-v3">
			<div class="public-notice-title-v3">Info pembayaran</div>
			<p class="public-subtitle">{{.PublicSubtitle}}</p>
		</div>
	</div>

	<div class="public-package-grid-v1">
	{{range .Packages}}
		<div class="public-package-card-v1">
			<h3 class="public-package-name-v1"><span class="public-package-prefix-v2">Durasi</span> {{.Name}}</h3>
			<div class="public-package-price-v1">{{rp .Price}}</div>
			<div class="public-package-quota-v2">Kuota <b>Unlimited</b></div>
			<a class="btn pri buy-btn" href="/order?package_id={{.ID}}{{$.BuyExtra}}">Pilih</a>
		</div>
	{{else}}
		<div class="public-empty-v1">Paket belum tersedia.</div>
	{{end}}
	</div>


<!-- PUBLIC_ROUTER_NOTICE_FULL_V3 -->
<style>
.public-notice-v3{
	margin-top:10px;
	padding:11px 12px;
	border-radius:17px;
	background:rgba(15,23,42,.34);
	border:1px solid rgba(148,163,184,.15);
}

.public-notice-title-v3{
	display:inline-flex;
	align-items:center;
	gap:6px;
	margin-bottom:6px;
	font-size:12px;
	line-height:1;
	font-weight:950;
	color:#93c5fd;
}

.public-notice-title-v3::before{
	content:"i";
	display:inline-grid;
	place-items:center;
	width:17px;
	height:17px;
	border-radius:999px;
	background:rgba(59,130,246,.18);
	color:#bfdbfe;
	font-size:11px;
	font-weight:950;
}

.public-notice-v3 .public-subtitle{
	margin:0!important;
	padding:0!important;
	background:transparent!important;
	border:0!important;
	color:#cbd5e1!important;
	font-size:13px!important;
	line-height:1.42!important;
	max-width:none!important;
	display:block!important;
	overflow:visible!important;
	-webkit-line-clamp:unset!important;
	-webkit-box-orient:unset!important;
}

html[data-theme="light"] .public-notice-v3,
html[data-admin-theme="light"] .public-notice-v3{
	background:rgba(241,245,249,.88)!important;
	border-color:#e2e8f0!important;
}

html[data-theme="light"] .public-notice-title-v3,
html[data-admin-theme="light"] .public-notice-title-v3{
	color:#2563eb!important;
}

html[data-theme="light"] .public-notice-v3 .public-subtitle,
html[data-admin-theme="light"] .public-notice-v3 .public-subtitle{
	color:#475569!important;
}

@media(max-width:720px){
	.public-router-hero-v1{
		padding:15px 14px!important;
		border-radius:21px!important;
	}

	.public-router-kicker-v1{
		margin-bottom:8px!important;
		min-height:28px!important;
		padding:5px 9px!important;
		font-size:11px!important;
	}

	.public-router-hero-v1 h1{
		font-size:27px!important;
		line-height:1!important;
	}

	.public-notice-v3{
		margin-top:9px!important;
		padding:10px 11px!important;
		border-radius:15px!important;
	}

	.public-notice-title-v3{
		font-size:11px!important;
		margin-bottom:5px!important;
	}

	.public-notice-v3 .public-subtitle{
		font-size:12px!important;
		line-height:1.36!important;
	}
}

@media(max-width:380px){
	.public-router-hero-v1{
		padding:13px 12px!important;
	}

	.public-router-hero-v1 h1{
		font-size:25px!important;
	}

	.public-notice-v3 .public-subtitle{
		font-size:11.5px!important;
		line-height:1.34!important;
	}
}
</style>


<!-- PUBLIC_ROUTER_NOTICE_LIST_V4 -->
<style>
@media(max-width:720px){
	.public-router-kicker-v1{
		display:none!important;
	}

	.public-router-hero-v1{
		padding:14px 14px!important;
		border-radius:21px!important;
	}

	.public-router-hero-v1 h1{
		font-size:28px!important;
		line-height:1!important;
		margin:0 0 10px!important;
	}

	.public-notice-v3{
		margin-top:0!important;
		padding:10px 11px!important;
		border-radius:15px!important;
	}

	.public-notice-title-v3{
		margin-bottom:7px!important;
		font-size:12px!important;
	}

	.public-notice-v3 .public-subtitle{
		display:none!important;
	}

	.public-notice-list-v4{
		margin:0!important;
		padding:0!important;
		list-style:none!important;
		display:grid!important;
		gap:6px!important;
	}

	.public-notice-list-v4 li{
		position:relative;
		margin:0!important;
		padding-left:15px!important;
		color:#cbd5e1!important;
		font-size:12px!important;
		line-height:1.34!important;
	}

	.public-notice-list-v4 li::before{
		content:"";
		position:absolute;
		left:0;
		top:.55em;
		width:5px;
		height:5px;
		border-radius:999px;
		background:#93c5fd;
	}

	html[data-theme="light"] .public-notice-list-v4 li,
	html[data-admin-theme="light"] .public-notice-list-v4 li{
		color:#475569!important;
	}

	html[data-theme="light"] .public-notice-list-v4 li::before,
	html[data-admin-theme="light"] .public-notice-list-v4 li::before{
		background:#2563eb!important;
	}
}

@media(max-width:380px){
	.public-router-hero-v1{
		padding:12px!important;
	}

	.public-router-hero-v1 h1{
		font-size:25px!important;
		margin-bottom:9px!important;
	}

	.public-notice-list-v4 li{
		font-size:11.5px!important;
		line-height:1.32!important;
	}
}
</style>

<script>
(function(){
	function ready(fn){
		if(document.readyState === "loading") document.addEventListener("DOMContentLoaded", fn);
		else fn();
	}

	ready(function(){
		document.querySelectorAll(".public-notice-v3 .public-subtitle").forEach(function(p){
			if(p.dataset.noticeListReady === "1") return;
			p.dataset.noticeListReady = "1";

			var raw = (p.textContent || "").replace(/\s+/g, " ").trim();
			raw = raw.replace(/^Info pembayaran\s*:?\s*/i, "").trim();

			var parts = raw.split(/\s*•\s*/).map(function(x){
				return x.trim();
			}).filter(function(x){
				return x.length > 0;
			});

			if(!parts.length) return;

			var ul = document.createElement("ul");
			ul.className = "public-notice-list-v4";

			parts.forEach(function(text){
				var li = document.createElement("li");
				li.textContent = text;
				ul.appendChild(li);
			});

			p.insertAdjacentElement("afterend", ul);
		});
	});
})();
</script>


<!-- PUBLIC_PACKAGE_DURATION_QUOTA_V2 -->
<style>
.public-package-prefix-v2{
	color:#93c5fd;
	font-size:.82em;
	font-weight:900;
}

.public-package-quota-v2{
	position:relative;
	z-index:1;
	margin:0 0 13px;
	color:#cbd5e1;
	font-size:13px;
	line-height:1.2;
	font-weight:800;
}

.public-package-quota-v2 b{
	color:#f8fafc;
	font-weight:950;
}

html[data-theme="light"] .public-package-prefix-v2,
html[data-admin-theme="light"] .public-package-prefix-v2{
	color:#2563eb!important;
}

html[data-theme="light"] .public-package-quota-v2,
html[data-admin-theme="light"] .public-package-quota-v2{
	color:#64748b!important;
}

html[data-theme="light"] .public-package-quota-v2 b,
html[data-admin-theme="light"] .public-package-quota-v2 b{
	color:#0f172a!important;
}

@media(max-width:720px){
	.public-package-card-v1{
		padding:12px!important;
	}

	.public-package-name-v1{
		font-size:14px!important;
		margin-bottom:7px!important;
		line-height:1.15!important;
	}

	.public-package-price-v1{
		font-size:23px!important;
		margin-bottom:7px!important;
	}

	.public-package-quota-v2{
		font-size:11.5px!important;
		margin-bottom:10px!important;
	}

	.public-package-card-v1 .buy-btn{
		min-height:38px!important;
	}
}

@media(max-width:360px){
	.public-package-name-v1{
		font-size:13px!important;
	}

	.public-package-price-v1{
		font-size:21px!important;
	}

	.public-package-quota-v2{
		font-size:11px!important;
	}
}
</style>

<!-- PUBLIC_ROUTER_DARK_PALETTE_LOCK_V1 -->
<style>
/* kunci dark palette final: harus berada paling bawah agar mengalahkan style lama */
html:not([data-theme="light"]) .public-router-hero-v1,
html[data-theme="dark"] .public-router-hero-v1,
html[data-admin-theme="dark"] .public-router-hero-v1{
	background:
		radial-gradient(circle at 18% 8%, rgba(34,211,238,.18), transparent 34%),
		radial-gradient(circle at 86% 18%, rgba(129,140,248,.16), transparent 34%),
		linear-gradient(135deg, rgba(6,12,26,.96), rgba(12,18,34,.94))!important;
	border-color:rgba(148,163,184,.14)!important;
	box-shadow:0 18px 45px rgba(2,6,23,.40)!important;
}

html:not([data-theme="light"]) .public-package-card-v1,
html[data-theme="dark"] .public-package-card-v1,
html[data-admin-theme="dark"] .public-package-card-v1,
html:not([data-theme="light"]) .public-empty-v1,
html[data-theme="dark"] .public-empty-v1,
html[data-admin-theme="dark"] .public-empty-v1{
	background:rgba(8,14,28,.82)!important;
	border-color:rgba(148,163,184,.14)!important;
	box-shadow:0 14px 36px rgba(2,6,23,.34)!important;
}

html:not([data-theme="light"]) .public-notice-v3,
html[data-theme="dark"] .public-notice-v3,
html[data-admin-theme="dark"] .public-notice-v3{
	background:rgba(8,14,28,.52)!important;
	border-color:rgba(148,163,184,.14)!important;
}

html:not([data-theme="light"]) .public-router-hero-v1 h1,
html[data-theme="dark"] .public-router-hero-v1 h1,
html[data-admin-theme="dark"] .public-router-hero-v1 h1,
html:not([data-theme="light"]) .public-package-name-v1,
html[data-theme="dark"] .public-package-name-v1,
html[data-admin-theme="dark"] .public-package-name-v1,
html:not([data-theme="light"]) .public-package-price-v1,
html[data-theme="dark"] .public-package-price-v1,
html[data-admin-theme="dark"] .public-package-price-v1,
html:not([data-theme="light"]) .public-package-quota-v2 b,
html[data-theme="dark"] .public-package-quota-v2 b,
html[data-admin-theme="dark"] .public-package-quota-v2 b{
	color:#f9fafb!important;
}

html:not([data-theme="light"]) .public-router-hero-v1 p,
html[data-theme="dark"] .public-router-hero-v1 p,
html[data-admin-theme="dark"] .public-router-hero-v1 p,
html:not([data-theme="light"]) .public-package-quota-v2,
html[data-theme="dark"] .public-package-quota-v2,
html[data-admin-theme="dark"] .public-package-quota-v2,
html:not([data-theme="light"]) .public-notice-v3 .public-subtitle,
html[data-theme="dark"] .public-notice-v3 .public-subtitle,
html[data-admin-theme="dark"] .public-notice-v3 .public-subtitle,
html:not([data-theme="light"]) .public-notice-list-v4 li,
html[data-theme="dark"] .public-notice-list-v4 li,
html[data-admin-theme="dark"] .public-notice-list-v4 li{
	color:#d7dde8!important;
}

html:not([data-theme="light"]) .public-router-kicker-v1,
html[data-theme="dark"] .public-router-kicker-v1,
html[data-admin-theme="dark"] .public-router-kicker-v1{
	background:rgba(34,197,94,.10)!important;
	border-color:rgba(34,197,94,.18)!important;
	color:#bbf7d0!important;
}

html:not([data-theme="light"]) .public-package-prefix-v2,
html[data-theme="dark"] .public-package-prefix-v2,
html[data-admin-theme="dark"] .public-package-prefix-v2,
html:not([data-theme="light"]) .public-notice-title-v3,
html[data-theme="dark"] .public-notice-title-v3,
html[data-admin-theme="dark"] .public-notice-title-v3{
	color:#7dd3fc!important;
}

html:not([data-theme="light"]) .public-notice-title-v3::before,
html[data-theme="dark"] .public-notice-title-v3::before,
html[data-admin-theme="dark"] .public-notice-title-v3::before{
	background:rgba(14,165,233,.16)!important;
	color:#bae6fd!important;
}

html:not([data-theme="light"]) .public-notice-list-v4 li::before,
html[data-theme="dark"] .public-notice-list-v4 li::before,
html[data-admin-theme="dark"] .public-notice-list-v4 li::before{
	background:#7dd3fc!important;
}
</style>

</div>`
