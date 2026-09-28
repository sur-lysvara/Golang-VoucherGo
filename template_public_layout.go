package main

const publicPageLayout = `<!doctype html>
<html lang="id" data-theme="dark">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}}</title>
<style>` + publicPageCSS + `</style>

<script>
// ADMIN_HIDE_SUPER_DEFAULT_V1
(function(){
  function getRole(){
    var m = document.cookie.match(/(?:^|;\s*)tuku_admin_role=([^;]+)/);
    return m ? decodeURIComponent(m[1]) : "";
  }

  function applyRole(){
    var role = getRole();

    document.documentElement.classList.remove("admin-role-router");
    document.documentElement.classList.remove("admin-role-superadmin");

    if(role === "superadmin"){
      document.documentElement.classList.add("admin-role-superadmin");
      return;
    }

    document.documentElement.classList.add("admin-role-router");

    document.querySelectorAll('.super-only,a[href="/admin/whatsapp"],a[href="/admin/settings"],a[href="/admin/users"]').forEach(function(el){
      el.style.display = "none";
      el.remove();
    });
  }

  applyRole();
  document.addEventListener("DOMContentLoaded", applyRole);
  setTimeout(applyRole, 50);
  setTimeout(applyRole, 250);
  setTimeout(applyRole, 800);
})();
</script>


<script>
// ADMIN_ROUTER_SETTINGS_COMPACT_V1
function toggleSecret(btn){
  var wrap = btn.closest(".secret-field");
  if(!wrap) return;
  var input = wrap.querySelector("input");
  if(!input) return;
  input.type = input.type === "password" ? "text" : "password";
}
</script>






<style>
/* VOUCHERGO_SECRET_EYE_POLISH_V1 */
.secret-field{
	position:relative!important;
}
.secret-field input{
	padding-right:58px!important;
}
.secret-field .pass-eye{
	position:absolute!important;
	right:10px!important;
	bottom:9px!important;
	top:auto!important;
	transform:none!important;
	width:42px!important;
	height:42px!important;
	min-width:42px!important;
	min-height:42px!important;
	max-width:42px!important;
	padding:0!important;
	margin:0!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	border-radius:14px!important;
	font-size:15px!important;
	line-height:1!important;
	z-index:5!important;
	cursor:pointer!important;
	background:rgba(15,23,42,.42)!important;
	border:1px solid rgba(148,163,184,.22)!important;
	color:#e5edf8!important;
	box-shadow:none!important;
}
.secret-field .pass-eye:hover{
	background:rgba(37,99,235,.20)!important;
	border-color:rgba(96,165,250,.45)!important;
}
html[data-theme="light"] .secret-field .pass-eye{
	background:#ffffff!important;
	border-color:#dbe4f0!important;
	color:#0f172a!important;
	box-shadow:0 4px 14px rgba(15,23,42,.06)!important;
}
html[data-theme="light"] .secret-field .pass-eye:hover{
	background:#eff6ff!important;
	border-color:#93c5fd!important;
}
@media(max-width:760px){
	.secret-field .pass-eye{
		width:40px!important;
		height:40px!important;
		min-width:40px!important;
		min-height:40px!important;
		right:8px!important;
		bottom:8px!important;
	}
	.secret-field input{
		padding-right:54px!important;
	}
}
</style>

<script>
// VOUCHERGO_SECRET_TOGGLE_MONKEY_V1
(function(){
	window.toggleSecret = function(btn){
		if(!btn) return;

		var wrap = btn.closest ? btn.closest(".secret-field") : null;
		var input = wrap ? wrap.querySelector("input") : null;

		if(!input && btn.parentElement){
			input = btn.parentElement.querySelector("input");
		}
		if(!input) return;

		var willShow = input.type === "password";
		input.type = willShow ? "text" : "password";

		btn.textContent = willShow ? "🙈" : "👁";
		btn.setAttribute("aria-label", willShow ? "Sembunyikan" : "Tampilkan");
		btn.setAttribute("title", willShow ? "Sembunyikan" : "Tampilkan");
	};

	function initSecretButtons(){
		document.querySelectorAll(".secret-field .pass-eye").forEach(function(btn){
			if(!btn.textContent.trim()){
				btn.textContent = "👁";
			}
			btn.setAttribute("type", "button");
			btn.setAttribute("aria-label", "Tampilkan");
			btn.setAttribute("title", "Tampilkan");
		});
	}

	document.addEventListener("DOMContentLoaded", initSecretButtons);
	setTimeout(initSecretButtons, 300);
})();
</script>






<script>
// VOUCHERGO_THEME_FINAL_V6
(function(){
	var KEY = "vouchergo_theme";

	function valid(t){
		return t === "light" || t === "dark";
	}

	function read(){
		try {
			var t = localStorage.getItem(KEY) || localStorage.getItem("theme") || localStorage.getItem("tuku_theme");
			return valid(t) ? t : "";
		} catch(e) {
			return "";
		}
	}

	function save(t){
		if (!valid(t)) return;
		try {
			localStorage.setItem(KEY, t);
			localStorage.setItem("theme", t);
			localStorage.setItem("tuku_theme", t);
		} catch(e) {}
	}

	function apply(t){
		if (!valid(t)) return;

		var html = document.documentElement;
		var body = document.body;

		html.setAttribute("data-theme", t);
		html.setAttribute("data-admin-theme", t);
		html.classList.toggle("theme-light", t === "light");
		html.classList.toggle("theme-dark", t === "dark");

		if (body) {
			body.setAttribute("data-theme", t);
			body.setAttribute("data-admin-theme", t);
			body.classList.toggle("theme-light", t === "light");
			body.classList.toggle("theme-dark", t === "dark");
		}

		document.querySelectorAll(".theme-toggle, [data-theme-toggle], #themeToggle, .theme-btn, .theme-switch").forEach(function(btn){
			btn.setAttribute("data-current-theme", t);
			btn.setAttribute("aria-label", "Theme " + t);
			btn.setAttribute("title", "Theme " + t);
		});
	}

	function current(){
		var html = document.documentElement;
		var body = document.body;

		var t = html.getAttribute("data-theme") || html.getAttribute("data-admin-theme");
		if (valid(t)) return t;

		if (body) {
			t = body.getAttribute("data-theme") || body.getAttribute("data-admin-theme");
			if (valid(t)) return t;
		}

		if (html.classList.contains("theme-light")) return "light";
		if (html.classList.contains("theme-dark")) return "dark";
		if (body && body.classList.contains("theme-light")) return "light";
		if (body && body.classList.contains("theme-dark")) return "dark";

		return read() || "dark";
	}

	function setTheme(t){
		if (!valid(t)) return;
		save(t);
		apply(t);
	}

	function toggleTheme(){
		setTheme(current() === "light" ? "dark" : "light");
	}

	window.setVoucherGoTheme = setTheme;
	window.toggleVoucherGoTheme = toggleTheme;

	function themeButton(target){
		if (!target || !target.closest) return null;

		var explicit = target.closest(".theme-toggle, [data-theme-toggle], #themeToggle, .theme-btn, .theme-switch");
		if (explicit) return explicit;

		var btn = target.closest("button, a");
		if (!btn) return null;

		var txt = (btn.textContent || "").trim();

		// Hanya icon-only. Tidak membaca class "light", jadi tombol Edit/Shortcut aman.
		if (txt === "🌙" || txt === "☀" || txt === "☀️" || txt === "🌗" || txt === "🌞" || txt === "🌚") {
			return btn;
		}

		return null;
	}

	// Apply theme tersimpan sedini mungkin.
	var saved = read();
	if (saved) apply(saved);

	document.addEventListener("DOMContentLoaded", function(){
		var saved = read();
		if (saved) apply(saved);

		document.addEventListener("click", function(ev){
			var btn = themeButton(ev.target);
			if (!btn) return;

			ev.preventDefault();
			ev.stopPropagation();
			if (ev.stopImmediatePropagation) ev.stopImmediatePropagation();

			toggleTheme();
		}, true);
	});
})();
</script>


<style>
/* PUBLIC_BEACH_MOBILE_POLISH_V1 */
html[data-theme="light"] body:not(.admin-body){
	background:
		linear-gradient(180deg,
			#7dd8ff 0%,
			#b7f3ff 40%,
			#f7f5d1 63%,
			#f4c96f 64%,
			#f2c163 100%) !important;
	min-height:100vh;
	overflow-x:hidden;
}

html[data-theme="light"] body:not(.admin-body)::before{
	content:"";
	position:fixed;
	left:0;
	right:0;
	bottom:30%;
	height:150px;
	pointer-events:none;
	z-index:0;
	background:
		radial-gradient(ellipse at 14% 88%, rgba(255,255,255,.55) 0 10%, transparent 11%),
		linear-gradient(180deg, rgba(35,181,211,.78), rgba(32,169,199,.72));
	border-radius:50% 50% 0 0 / 34% 34% 0 0;
	transform:translateY(52%);
	box-shadow:0 -18px 38px rgba(18,126,160,.14);
}

html[data-theme="light"] body:not(.admin-body)::after{
	content:"";
	position:fixed;
	right:7vw;
	top:11vh;
	width:86px;
	height:86px;
	pointer-events:none;
	z-index:0;
	border-radius:50%;
	background:#ffd34d;
	box-shadow:
		0 0 0 16px rgba(255,211,77,.22),
		0 18px 42px rgba(255,174,0,.22);
}

html[data-theme="light"] .public-wrap,
html[data-theme="light"] .public-page,
html[data-theme="light"] .public-shell,
html[data-theme="light"] main{
	position:relative;
	z-index:2;
}

html[data-theme="light"] .public-shell::before,
html[data-theme="light"] .public-shell::after{
	content:"";
	position:fixed;
	pointer-events:none;
	z-index:1;
	border-radius:999px;
	background:rgba(255,255,255,.72);
	filter:blur(.2px);
}

html[data-theme="light"] .public-shell::before{
	left:4vw;
	top:22vh;
	width:135px;
	height:36px;
	box-shadow:
		38px -18px 0 2px rgba(255,255,255,.62),
		82px 5px 0 -3px rgba(255,255,255,.66);
}

html[data-theme="light"] .public-shell::after{
	left:10vw;
	bottom:26%;
	width:120px;
	height:18px;
	opacity:.28;
	background:rgba(255,255,255,.9);
	box-shadow:0 0 28px 12px rgba(255,255,255,.38);
}

html[data-theme="light"] .public-card,
html[data-theme="light"] .router-card,
html[data-theme="light"] .package-card{
	background:rgba(255,255,255,.72) !important;
	backdrop-filter:blur(14px);
	-webkit-backdrop-filter:blur(14px);
	border:1px solid rgba(255,255,255,.58) !important;
	box-shadow:0 18px 45px rgba(15,55,80,.16) !important;
}

html[data-theme="light"] .hero,
html[data-theme="light"] .public-hero{
	box-shadow:0 18px 45px rgba(15,55,80,.18) !important;
}

@media(max-width:760px){
	html[data-theme="light"] body:not(.admin-body){
		background:
			linear-gradient(180deg,
				#78d9ff 0%,
				#aef1ff 45%,
				#f8f4cf 67%,
				#f4c86d 68%,
				#efbd59 100%) !important;
	}

	html[data-theme="light"] body:not(.admin-body)::before{
		bottom:24%;
		height:110px;
		transform:translateY(45%);
		border-radius:55% 55% 0 0 / 28% 28% 0 0;
		background:linear-gradient(180deg, rgba(34,181,213,.74), rgba(27,158,191,.70));
	}

	html[data-theme="light"] body:not(.admin-body)::after{
		right:9vw;
		top:15vh;
		width:54px;
		height:54px;
		box-shadow:
			0 0 0 10px rgba(255,211,77,.20),
			0 12px 28px rgba(255,174,0,.20);
	}

	html[data-theme="light"] .public-shell::before{
		left:8vw;
		top:26vh;
		width:86px;
		height:25px;
		opacity:.72;
		box-shadow:
			25px -11px 0 1px rgba(255,255,255,.62),
			55px 4px 0 -4px rgba(255,255,255,.66);
	}

	html[data-theme="light"] .public-shell::after{
		display:none;
	}

	html[data-theme="light"] .public-wrap,
	html[data-theme="light"] .public-page,
	html[data-theme="light"] .public-shell{
		padding-left:14px !important;
		padding-right:14px !important;
	}

	html[data-theme="light"] .public-card,
	html[data-theme="light"] .router-card,
	html[data-theme="light"] .package-card{
		background:rgba(255,255,255,.78) !important;
		box-shadow:0 12px 28px rgba(15,55,80,.14) !important;
	}

	html[data-theme="light"] .hero,
	html[data-theme="light"] .public-hero{
		border-radius:18px !important;
	}
}


/* PUBLIC_CARD_TEXT_CONTRAST_V1_START */
html[data-theme="light"] .router-card,
html[data-theme="light"] .package-card,
html[data-theme="light"] .public-card{
        color:#0f172a !important;
}

html[data-theme="light"] .router-card :is(h1,h2,h3,h4,strong,.title,.name,.router-name),
html[data-theme="light"] .public-card.router-card :is(h1,h2,h3,h4,strong,.title,.name,.router-name){
        color:#0f172a !important;
        font-weight:800 !important;
        letter-spacing:-0.02em;
}

html[data-theme="light"] .router-card :is(p,.muted,.desc,small,span:not(.btn):not(.badge)),
html[data-theme="light"] .public-card.router-card :is(p,.muted,.desc,small,span:not(.btn):not(.badge)){
        color:#64748b !important;
        opacity:1 !important;
        font-weight:600 !important;
}

html[data-theme="light"] .package-card :is(h1,h2,h3,h4,strong,.title,.name,.pkg-name),
html[data-theme="light"] .public-card.package-card :is(h1,h2,h3,h4,strong,.title,.name,.pkg-name){
        color:#0f172a !important;
        font-weight:800 !important;
        opacity:1 !important;
}

html[data-theme="light"] .package-card :is(.price,.amount,.pkg-price,.voucher-price),
html[data-theme="light"] .public-card.package-card :is(.price,.amount,.pkg-price,.voucher-price){
        color:#0f172a !important;
        font-weight:900 !important;
        opacity:1 !important;
        text-shadow:none !important;
}

html[data-theme="light"] .package-card :is(p,.muted,.desc,small,label,span:not(.btn):not(.badge)),
html[data-theme="light"] .public-card.package-card :is(p,.muted,.desc,small,label,span:not(.btn):not(.badge)){
        color:#64748b !important;
        opacity:1 !important;
        font-weight:600 !important;
}

html[data-theme="light"] .router-card .btn,
html[data-theme="light"] .router-card button,
html[data-theme="light"] .package-card .btn,
html[data-theme="light"] .package-card button,
html[data-theme="light"] .public-card .btn,
html[data-theme="light"] .public-card button{
        color:#fff !important;
}

@media (min-width: 761px){
        html[data-theme="light"] .router-card :is(h1,h2,h3,h4,strong,.title,.name,.router-name){
                font-size:36px !important;
                line-height:1.08 !important;
        }

        html[data-theme="light"] .router-card :is(p,.muted,.desc,small,span:not(.btn):not(.badge)){
                font-size:19px !important;
                line-height:1.45 !important;
        }

        html[data-theme="light"] .package-card :is(h1,h2,h3,h4,strong,.title,.name,.pkg-name){
                font-size:28px !important;
                line-height:1.12 !important;
        }

        html[data-theme="light"] .package-card :is(.price,.amount,.pkg-price,.voucher-price){
                font-size:54px !important;
                line-height:1.02 !important;
        }

        html[data-theme="light"] .package-card :is(p,.muted,.desc,small,label,span:not(.btn):not(.badge)){
                font-size:18px !important;
                line-height:1.4 !important;
        }
}

@media (max-width: 760px){
        html[data-theme="light"] .router-card :is(h1,h2,h3,h4,strong,.title,.name,.router-name){
                font-size:24px !important;
                line-height:1.12 !important;
        }

        html[data-theme="light"] .router-card :is(p,.muted,.desc,small,span:not(.btn):not(.badge)){
                font-size:14px !important;
                line-height:1.4 !important;
        }

        html[data-theme="light"] .package-card :is(h1,h2,h3,h4,strong,.title,.name,.pkg-name){
                font-size:20px !important;
                line-height:1.1 !important;
        }

        html[data-theme="light"] .package-card :is(.price,.amount,.pkg-price,.voucher-price){
                font-size:36px !important;
                line-height:1.05 !important;
        }

        html[data-theme="light"] .package-card :is(p,.muted,.desc,small,label,span:not(.btn):not(.badge)){
                font-size:14px !important;
                line-height:1.35 !important;
        }
}
/* PUBLIC_CARD_TEXT_CONTRAST_V1_END */

</style>

</head>
<body>
<div class="public-day-scene" aria-hidden="true">
  <div class="day-sun"></div>
  <div class="day-cloud"></div>
  <div class="day-cat">
    <div class="umbrella"></div>
    <div class="mat"></div>
    <div class="cat-face">
      <span class="eye left"></span>
      <span class="eye right"></span>
      <span class="nose"></span>
    </div>
  </div>
</div>

<figure class="catpen-scene" aria-hidden="true">
  <figcaption>Animated cat astronaut sitting on the moon</figcaption>
  <div class="catpen-moon"></div>
  <div class="catpen-bubble"></div>
  <div class="catpen-backpack"></div>

  <article class="catpen-cat">
    <span id="alt">cat astronaut</span>

    <div class="catpen-tail"></div>

    <div class="catpen-body">
      <div class="catpen-leg"></div>
      <div class="catpen-paw"></div>
      <div class="catpen-paw"></div>
    </div>

    <div class="catpen-ear"></div>
    <div class="catpen-ear"></div>

    <div class="catpen-head">
      <div class="catpen-whisker"></div>
      <div class="catpen-whisker"></div>
      <div class="catpen-whisker"></div>
      <div class="catpen-whisker"></div>

      <div class="catpen-eye"></div>
      <div class="catpen-eye"></div>

      <div class="catpen-nose"></div>
    </div>
  </article>
</figure>

<div class="tuku-cat">
    <div class="tuku-cat-tail"></div>
    <div class="tuku-cat-body"></div>
    <div class="tuku-cat-head">
      <span class="ear left"></span>
      <span class="ear right"></span>
      <span class="eye left"></span>
      <span class="eye right"></span>
      <span class="nose"></span>
    </div>
  </div>
</div>





<!-- PUBLIC_CARD_TEXT_LIVE_V3_START -->
<style>
html[data-theme="light"] .wrap > .grid > .card{
	color:#0f172a !important;
	background:rgba(255,255,255,.88) !important;
	border:1px solid rgba(255,255,255,.78) !important;
	box-shadow:0 18px 42px rgba(15,55,80,.18) !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .grid > .card b{
	display:block !important;
	color:#0f172a !important;
	font-size:24px !important;
	line-height:1.15 !important;
	font-weight:950 !important;
	letter-spacing:-.02em !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .grid > .card h3{
	color:#0f172a !important;
	font-size:24px !important;
	line-height:1.12 !important;
	font-weight:950 !important;
	margin:0 0 8px !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .grid > .card .price{
	color:#0f172a !important;
	font-size:34px !important;
	line-height:1.05 !important;
	font-weight:950 !important;
	margin:4px 0 8px !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .grid > .card p,
html[data-theme="light"] .wrap > .grid > .card .muted{
	color:#475569 !important;
	opacity:1 !important;
	font-size:17px !important;
	line-height:1.45 !important;
	font-weight:750 !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .grid > .card .btn,
html[data-theme="light"] .wrap > .grid > .card .buy-btn{
	color:#fff !important;
	font-size:16px !important;
	font-weight:950 !important;
	min-height:48px !important;
	border-radius:14px !important;
	text-shadow:none !important;
}

@media(max-width:760px){
	html[data-theme="light"] .wrap > .grid{
		gap:12px !important;
	}

	html[data-theme="light"] .wrap > .grid > .card{
		padding:16px !important;
		border-radius:18px !important;
		min-height:96px !important;
	}

	html[data-theme="light"] .wrap > .grid > .card b{
		font-size:18px !important;
		line-height:1.15 !important;
	}

	html[data-theme="light"] .wrap > .grid > .card h3{
		font-size:18px !important;
		margin:0 0 6px !important;
	}

	html[data-theme="light"] .wrap > .grid > .card .price{
		font-size:28px !important;
		margin:2px 0 6px !important;
	}

	html[data-theme="light"] .wrap > .grid > .card p,
	html[data-theme="light"] .wrap > .grid > .card .muted{
		font-size:13px !important;
		line-height:1.35 !important;
	}

	html[data-theme="light"] .wrap > .grid > .card .btn,
	html[data-theme="light"] .wrap > .grid > .card .buy-btn{
		font-size:14px !important;
		min-height:42px !important;
		border-radius:13px !important;
	}
}
</style>
<!-- PUBLIC_CARD_TEXT_LIVE_V3_END -->


<!-- PUBLIC_CHECKOUT_CONTRAST_V1_START -->
<style>
html[data-theme="light"] .wrap > .card{
	color:#0f172a !important;
	background:rgba(255,255,255,.94) !important;
	border:1px solid rgba(255,255,255,.80) !important;
	box-shadow:0 22px 55px rgba(15,55,80,.20) !important;
	backdrop-filter:blur(14px) !important;
	-webkit-backdrop-filter:blur(14px) !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .card .card{
	color:#0f172a !important;
	background:rgba(255,255,255,.96) !important;
	border:1px solid rgba(203,213,225,.80) !important;
	box-shadow:0 10px 28px rgba(15,55,80,.10) !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .card h1,
html[data-theme="light"] .wrap > .card h2,
html[data-theme="light"] .wrap > .card h3,
html[data-theme="light"] .wrap > .card b,
html[data-theme="light"] .wrap > .card strong,
html[data-theme="light"] .wrap > .card label,
html[data-theme="light"] .wrap > .card .price{
	color:#0f172a !important;
	opacity:1 !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .card h3{
	font-size:26px !important;
	line-height:1.12 !important;
	font-weight:900 !important;
	margin:0 0 8px !important;
}

html[data-theme="light"] .wrap > .card .price{
	font-size:42px !important;
	line-height:1.05 !important;
	font-weight:950 !important;
	margin:4px 0 10px !important;
}

html[data-theme="light"] .wrap > .card p,
html[data-theme="light"] .wrap > .card .muted,
html[data-theme="light"] .wrap > .card small{
	color:#475569 !important;
	opacity:1 !important;
	font-weight:700 !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .card label{
	color:#334155 !important;
	font-weight:900 !important;
	font-size:14px !important;
}

html[data-theme="light"] .wrap > .card input,
html[data-theme="light"] .wrap > .card select,
html[data-theme="light"] .wrap > .card textarea{
	color:#0f172a !important;
	background:rgba(255,255,255,.98) !important;
	border:1px solid rgba(148,163,184,.45) !important;
	box-shadow:0 8px 22px rgba(15,55,80,.06) !important;
}

html[data-theme="light"] .wrap > .card input::placeholder,
html[data-theme="light"] .wrap > .card textarea::placeholder{
	color:#64748b !important;
	opacity:.95 !important;
}

html[data-theme="light"] .wrap > .card .btn,
html[data-theme="light"] .wrap > .card button{
	color:#fff !important;
	font-weight:950 !important;
	text-shadow:none !important;
}

@media(max-width:760px){
	html[data-theme="light"] .wrap > .card{
		padding:16px !important;
		border-radius:20px !important;
		background:rgba(255,255,255,.95) !important;
	}

	html[data-theme="light"] .wrap > .card .card{
		padding:14px !important;
		border-radius:16px !important;
	}

	html[data-theme="light"] .wrap > .card h3{
		font-size:21px !important;
	}

	html[data-theme="light"] .wrap > .card .price{
		font-size:32px !important;
	}

	html[data-theme="light"] .wrap > .card p,
	html[data-theme="light"] .wrap > .card .muted,
	html[data-theme="light"] .wrap > .card small{
		font-size:13px !important;
		line-height:1.35 !important;
	}

	html[data-theme="light"] .wrap > .card label{
		font-size:12px !important;
	}

	html[data-theme="light"] .wrap > .card input,
	html[data-theme="light"] .wrap > .card select,
	html[data-theme="light"] .wrap > .card textarea{
		min-height:44px !important;
		border-radius:13px !important;
		font-size:14px !important;
	}
}
</style>
<!-- PUBLIC_CHECKOUT_CONTRAST_V1_END -->





<!-- PUBLIC_COMPACT_VISIBLE_V2_START -->
<style>
/* PUBLIC_COMPACT_VISIBLE_V2 */

/* ===== desktop/laptop: lebih kecil & rapat ===== */
html[data-theme="light"] .wrap{
        max-width: 980px !important;
        padding-top: 10px !important;
        padding-bottom: 18px !important;
}

html[data-theme="light"] .hero{
        max-width: 760px !important;
        margin: 0 auto 10px !important;
        padding: 18px 22px !important;
        border-radius: 20px !important;
}

html[data-theme="light"] .hero h1{
        font-size: 28px !important;
        line-height: 1.08 !important;
        margin: 0 0 4px !important;
}

html[data-theme="light"] .hero p{
        font-size: 14px !important;
        line-height: 1.3 !important;
        margin: 0 !important;
}

/* body checkout diperkecil agar lebih muat di layar laptop */
html[data-theme="light"] .wrap > .card{
        max-width: 620px !important;
        margin: 8px auto 0 !important;
        padding: 14px !important;
        border-radius: 18px !important;
        background: rgba(255,255,255,.96) !important;
        border: 1px solid rgba(255,255,255,.82) !important;
        box-shadow: 0 14px 34px rgba(15,55,80,.14) !important;
}

html[data-theme="light"] .wrap > .card .card{
        padding: 14px !important;
        border-radius: 16px !important;
        background: rgba(255,255,255,.98) !important;
        border: 1px solid rgba(203,213,225,.78) !important;
        box-shadow: 0 8px 18px rgba(15,55,80,.06) !important;
}

html[data-theme="light"] .wrap > .card h3{
        font-size: 18px !important;
        line-height: 1.08 !important;
        margin: 0 0 2px !important;
        color: #0f172a !important;
        font-weight: 900 !important;
}

html[data-theme="light"] .wrap > .card .price{
        font-size: 30px !important;
        line-height: 1 !important;
        margin: 4px 0 8px !important;
        color: #0f172a !important;
        font-weight: 950 !important;
}

html[data-theme="light"] .wrap > .card p,
html[data-theme="light"] .wrap > .card .muted,
html[data-theme="light"] .wrap > .card small{
        font-size: 13px !important;
        line-height: 1.28 !important;
        color: #475569 !important;
        opacity: 1 !important;
        font-weight: 700 !important;
}

html[data-theme="light"] .wrap > .card label{
        font-size: 12px !important;
        margin-bottom: 4px !important;
        color: #334155 !important;
        font-weight: 900 !important;
}

html[data-theme="light"] .wrap > .card input,
html[data-theme="light"] .wrap > .card select,
html[data-theme="light"] .wrap > .card textarea{
        min-height: 38px !important;
        padding: 8px 13px !important;
        font-size: 14px !important;
        border-radius: 12px !important;
        color: #0f172a !important;
        background: #ffffff !important;
        border: 1px solid rgba(148,163,184,.50) !important;
        box-shadow: 0 4px 12px rgba(15,55,80,.04) !important;
}

html[data-theme="light"] .wrap > .card .btn,
html[data-theme="light"] .wrap > .card button{
        min-height: 40px !important;
        border-radius: 13px !important;
        font-size: 14px !important;
        font-weight: 900 !important;
}

/* tombol/link "pilih paket lain" dipaksa kelihatan di mode light */
html[data-theme="light"] .wrap > .card a:not(.btn):not(.buy-btn){
        display: inline-flex !important;
        align-items: center !important;
        justify-content: center !important;
        gap: 6px !important;
        min-height: 38px !important;
        padding: 0 14px !important;
        border-radius: 12px !important;
        background: #eef2ff !important;
        border: 1px solid rgba(99,102,241,.18) !important;
        color: #1e293b !important;
        font-size: 13px !important;
        font-weight: 900 !important;
        text-decoration: none !important;
        opacity: 1 !important;
        visibility: visible !important;
        box-shadow: 0 6px 14px rgba(15,23,42,.06) !important;
}

html[data-theme="light"] .wrap > .card a:not(.btn):not(.buy-btn):hover{
        background: #e0e7ff !important;
        color: #0f172a !important;
}

/* dark juga dirapikan sedikit */
html[data-theme="dark"] .wrap > .card{
        max-width: 620px !important;
        margin: 8px auto 0 !important;
        padding: 14px !important;
        border-radius: 18px !important;
}

html[data-theme="dark"] .wrap > .card .card{
        padding: 14px !important;
        border-radius: 16px !important;
}

html[data-theme="dark"] .wrap > .card h3{
        font-size: 18px !important;
        margin: 0 0 2px !important;
}

html[data-theme="dark"] .wrap > .card .price{
        font-size: 30px !important;
        margin: 4px 0 8px !important;
}

html[data-theme="dark"] .wrap > .card label{
        font-size: 12px !important;
        margin-bottom: 4px !important;
}

html[data-theme="dark"] .wrap > .card input,
html[data-theme="dark"] .wrap > .card select,
html[data-theme="dark"] .wrap > .card textarea{
        min-height: 38px !important;
        padding: 8px 13px !important;
        font-size: 14px !important;
        border-radius: 12px !important;
}

html[data-theme="dark"] .wrap > .card .btn,
html[data-theme="dark"] .wrap > .card button{
        min-height: 40px !important;
        font-size: 14px !important;
        border-radius: 13px !important;
}

/* mobile tetap aman */
@media (max-width: 760px){
        html[data-theme="light"] .wrap,
        html[data-theme="dark"] .wrap{
                max-width: none !important;
                padding: 12px 12px 18px !important;
        }

        html[data-theme="light"] .hero,
        html[data-theme="dark"] .hero{
                max-width: none !important;
                margin: 0 0 10px !important;
                padding: 16px !important;
                border-radius: 18px !important;
        }

        html[data-theme="light"] .hero h1,
        html[data-theme="dark"] .hero h1{
                font-size: 24px !important;
        }

        html[data-theme="light"] .hero p,
        html[data-theme="dark"] .hero p{
                font-size: 13px !important;
        }

        html[data-theme="light"] .wrap > .card,
        html[data-theme="dark"] .wrap > .card{
                max-width: none !important;
                margin: 8px 0 0 !important;
                padding: 12px !important;
                border-radius: 16px !important;
        }

        html[data-theme="light"] .wrap > .card .card,
        html[data-theme="dark"] .wrap > .card .card{
                padding: 12px !important;
                border-radius: 14px !important;
        }

        html[data-theme="light"] .wrap > .card h3,
        html[data-theme="dark"] .wrap > .card h3{
                font-size: 17px !important;
        }

        html[data-theme="light"] .wrap > .card .price,
        html[data-theme="dark"] .wrap > .card .price{
                font-size: 28px !important;
        }
}
</style>
<!-- PUBLIC_COMPACT_VISIBLE_V2_END -->


<!-- PUBLIC_CHECKOUT_BACK_BUTTON_FIX_V3_START -->
<style>
/* PUBLIC_CHECKOUT_BACK_BUTTON_FIX_V3 */

/* checkout makin compact */
html[data-theme="light"] .wrap > .card,
html[data-theme="dark"] .wrap > .card{
	max-width:560px !important;
	margin-top:6px !important;
	padding:12px !important;
	border-radius:17px !important;
}

html[data-theme="light"] .wrap > .card .card,
html[data-theme="dark"] .wrap > .card .card{
	padding:12px !important;
	border-radius:15px !important;
}

html[data-theme="light"] .wrap > .card .price,
html[data-theme="dark"] .wrap > .card .price{
	font-size:28px !important;
	margin:2px 0 6px !important;
}

html[data-theme="light"] .wrap > .card h3,
html[data-theme="dark"] .wrap > .card h3{
	font-size:17px !important;
	margin:0 0 2px !important;
}

html[data-theme="light"] .wrap > .card input,
html[data-theme="dark"] .wrap > .card input{
	min-height:36px !important;
	height:36px !important;
	padding:7px 12px !important;
	font-size:13px !important;
}

html[data-theme="light"] .wrap > .card .btn,
html[data-theme="dark"] .wrap > .card .btn{
	min-height:38px !important;
}

/* tombol/link pilih paket lain: paksa kelihatan */
html[data-theme="light"] .wrap > .card a.btn:not(.pri),
html[data-theme="light"] .wrap > .card a:not(.buy-btn):not(.admin-mini),
html[data-theme="light"] .wrap > .card a[href^="/"]:not(.buy-btn):not(.admin-mini){
	display:inline-flex !important;
	align-items:center !important;
	justify-content:center !important;
	width:auto !important;
	min-width:150px !important;
	min-height:38px !important;
	padding:0 14px !important;
	margin-top:10px !important;
	border-radius:13px !important;
	background:#e0e7ff !important;
	border:1px solid rgba(99,102,241,.35) !important;
	color:#0f172a !important;
	font-size:14px !important;
	font-weight:950 !important;
	text-decoration:none !important;
	opacity:1 !important;
	visibility:visible !important;
	box-shadow:0 8px 18px rgba(15,23,42,.08) !important;
	text-shadow:none !important;
}

html[data-theme="light"] .wrap > .card a.btn:not(.pri)::before,
html[data-theme="light"] .wrap > .card a:not(.buy-btn):not(.admin-mini)::before{
	color:#0f172a !important;
}

html[data-theme="dark"] .wrap > .card a.btn:not(.pri),
html[data-theme="dark"] .wrap > .card a:not(.buy-btn):not(.admin-mini),
html[data-theme="dark"] .wrap > .card a[href^="/"]:not(.buy-btn):not(.admin-mini){
	display:inline-flex !important;
	align-items:center !important;
	justify-content:center !important;
	width:auto !important;
	min-width:150px !important;
	min-height:38px !important;
	padding:0 14px !important;
	margin-top:10px !important;
	border-radius:13px !important;
	background:rgba(15,23,42,.80) !important;
	border:1px solid rgba(148,163,184,.26) !important;
	color:#ffffff !important;
	font-size:14px !important;
	font-weight:950 !important;
	text-decoration:none !important;
	opacity:1 !important;
	visibility:visible !important;
	text-shadow:none !important;
}

@media(max-height:820px) and (min-width:761px){
	html[data-theme="light"] .wrap,
	html[data-theme="dark"] .wrap{
		padding-top:8px !important;
		padding-bottom:12px !important;
	}

	html[data-theme="light"] .hero,
	html[data-theme="dark"] .hero{
		max-width:720px !important;
		padding:15px 20px !important;
		margin-bottom:8px !important;
	}

	html[data-theme="light"] .hero h1,
	html[data-theme="dark"] .hero h1{
		font-size:26px !important;
	}

	html[data-theme="light"] .wrap > .card,
	html[data-theme="dark"] .wrap > .card{
		max-width:540px !important;
		padding:11px !important;
	}
}

@media(max-width:760px){
	html[data-theme="light"] .wrap > .card,
	html[data-theme="dark"] .wrap > .card{
		max-width:none !important;
		padding:12px !important;
	}

	html[data-theme="light"] .wrap > .card a.btn:not(.pri),
	html[data-theme="light"] .wrap > .card a:not(.buy-btn):not(.admin-mini),
	html[data-theme="dark"] .wrap > .card a.btn:not(.pri),
	html[data-theme="dark"] .wrap > .card a:not(.buy-btn):not(.admin-mini){
		width:100% !important;
		min-width:0 !important;
	}
}
</style>
<!-- PUBLIC_CHECKOUT_BACK_BUTTON_FIX_V3_END -->



<div class="PUBLIC_LAYOUT_ACTIVE">PUBLIC_LAYOUT_ACTIVE</div>
<div class="wrap">
  <div class="public-top">
    <div class="public-brand">{{.Title}}</div>
    <div class="public-actions">
      <div style="display:inline-flex; align-items:center;">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right:8px; opacity:0.8; cursor:pointer;" onclick="this.nextElementSibling.click()"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2"></path><path d="M12 20v2"></path><path d="m4.93 4.93 1.41 1.41"></path><path d="m17.66 17.66 1.41 1.41"></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path></svg>
        <button class="theme-toggle" type="button" aria-label="Ganti tema" title="Ganti tema"><span class="sun">☀️</span><span class="moon">🌙</span></button>
      </div>
      <a class="admin-mini" href="/login" aria-label="Admin" title="Admin"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg></a>
    </div>
  </div>
  {{.Body}}
</div>


<style>
/* PUBLIC_DISABLE_LEGACY_ANIMATION_V1 */
/* Matikan animasi bawaan public seperti tukuPageIn/tukuCardIn agar tidak numpuk. */
.wrap,
.wrap > h1,
.hero,
.public-topbar,
.grid,
.grid > .card,
.grid > a,
.card,
.router-card,
.package-card,
.public-card,
.checkout-card,
.invoice-wrap .card,
.voucher-box,
form.card,
.pay-card{
	animation:none!important;
	opacity:1!important;
	visibility:visible!important;
	filter:none!important;
}

.grid > .card,
.grid > a,
.router-card,
.package-card,
.public-card{
	transform:none;
}
</style>


<style>
/* PUBLIC_CHECKOUT_INVOICE_POLISH_V1 */
.wrap > .card,
.invoice-wrap,
.checkout-card,
.pay-card{
	position:relative;
	overflow:hidden;
}

html[data-theme="light"] .wrap > .card,
html[data-theme="light"] .checkout-card,
html[data-theme="light"] .invoice-wrap{
	background:rgba(255,255,255,.94)!important;
	border:1px solid rgba(255,255,255,.84)!important;
	box-shadow:
		0 22px 60px rgba(15,55,80,.16),
		inset 0 1px 0 rgba(255,255,255,.86)!important;
}

html[data-theme="dark"] .wrap > .card,
html[data-theme="dark"] .checkout-card,
html[data-theme="dark"] .invoice-wrap{
	background:rgba(15,23,42,.78)!important;
	border:1px solid rgba(148,163,184,.18)!important;
	box-shadow:
		0 24px 70px rgba(0,0,0,.34),
		inset 0 1px 0 rgba(255,255,255,.06)!important;
}

.wrap > .card::before,
.invoice-wrap::before,
.checkout-card::before{
	content:"";
	position:absolute;
	inset:0 0 auto 0;
	height:4px;
	background:linear-gradient(90deg,#38bdf8,#6366f1,#a78bfa,#38bdf8);
	opacity:.9;
	pointer-events:none;
}

.wrap > .card h1,
.wrap > .card h2,
.wrap > .card h3,
.invoice-wrap h1,
.invoice-wrap h2,
.invoice-wrap h3{
	letter-spacing:-.035em!important;
	font-weight:950!important;
}

.wrap > .card input,
.wrap > .card select,
.wrap > .card textarea,
.checkout-card input,
.checkout-card select,
.checkout-card textarea{
	transition:border-color .18s ease, box-shadow .18s ease, transform .18s ease;
}

.wrap > .card input:focus,
.wrap > .card select:focus,
.wrap > .card textarea:focus,
.checkout-card input:focus,
.checkout-card select:focus,
.checkout-card textarea:focus{
	outline:none!important;
	border-color:rgba(59,130,246,.72)!important;
	box-shadow:
		0 0 0 4px rgba(59,130,246,.13),
		0 10px 24px rgba(15,55,80,.08)!important;
	transform:translateY(-1px);
}

.wrap > .card .btn,
.wrap > .card button,
.checkout-card .btn,
.checkout-card button,
.invoice-wrap .btn,
.invoice-wrap button{
	font-weight:950!important;
	letter-spacing:-.01em!important;
	box-shadow:0 12px 28px rgba(37,99,235,.18)!important;
}

.invoice-wrap .voucher-box,
.voucher-box{
	background:
		linear-gradient(180deg, rgba(239,246,255,.96), rgba(255,255,255,.94))!important;
	border:2px dashed rgba(59,130,246,.44)!important;
	box-shadow:
		0 16px 38px rgba(37,99,235,.10),
		inset 0 1px 0 rgba(255,255,255,.86)!important;
}

html[data-theme="dark"] .invoice-wrap .voucher-box,
html[data-theme="dark"] .voucher-box{
	background:
		linear-gradient(180deg, rgba(30,41,59,.88), rgba(15,23,42,.86))!important;
	border-color:rgba(96,165,250,.42)!important;
	box-shadow:
		0 18px 46px rgba(0,0,0,.28),
		inset 0 1px 0 rgba(255,255,255,.06)!important;
}

.voucher-code{
	display:inline-block!important;
	padding:8px 14px!important;
	border-radius:18px!important;
	background:rgba(15,23,42,.06)!important;
	color:#0f172a!important;
	text-shadow:none!important;
}

html[data-theme="dark"] .voucher-code{
	background:rgba(255,255,255,.08)!important;
	color:#f8fafc!important;
}

.invoice-wrap .muted,
.wrap > .card .muted{
	opacity:.86!important;
}

@media(max-width:760px){
	.wrap > .card,
	.checkout-card,
	.invoice-wrap{
		border-radius:18px!important;
		padding:14px!important;
	}

	.wrap > .card::before,
	.invoice-wrap::before,
	.checkout-card::before{
		height:3px;
	}

	.wrap > .card h1,
	.wrap > .card h2,
	.wrap > .card h3,
	.invoice-wrap h1,
	.invoice-wrap h2,
	.invoice-wrap h3{
		font-size:20px!important;
		line-height:1.12!important;
	}

	.voucher-code{
		font-size:30px!important;
		letter-spacing:.08em!important;
		padding:6px 10px!important;
		border-radius:14px!important;
		max-width:100%!important;
		overflow-wrap:anywhere!important;
	}

	.invoice-wrap .voucher-box,
	.voucher-box{
		padding:16px!important;
		border-radius:18px!important;
	}
}
</style>


<style>
/* PUBLIC_INVOICE_STATUS_POLISH_V2 */
.invoice-wrap{
	max-width:860px!important;
	margin-left:auto!important;
	margin-right:auto!important;
}

.invoice-wrap .card,
.invoice-wrap .voucher-box,
.invoice-wrap .pay-card{
	position:relative!important;
	overflow:hidden!important;
}

.invoice-wrap .card::before,
.invoice-wrap .pay-card::before{
	content:"";
	position:absolute;
	inset:0 0 auto 0;
	height:3px;
	background:linear-gradient(90deg,#38bdf8,#6366f1,#a78bfa);
	opacity:.9;
	pointer-events:none;
}

/* VG_VOUCHER_CARD_TEXT_CLEAN_V2 */
.invoice-wrap .voucher-box::before{
	content:none;
	display:none;
	align-items:center;
	justify-content:center;
	margin:0;
	padding:0;
	border-radius:999px;
	background:rgba(37,99,235,.10);
	color:#2563eb;
	font-size:12px;
	font-weight:950;
	letter-spacing:.02em;
}

html[data-theme="dark"] .invoice-wrap .voucher-box::before{
	background:rgba(96,165,250,.14);
	color:#bfdbfe;
}

.invoice-wrap .voucher-code{
	margin:8px auto 10px!important;
	box-shadow:inset 0 1px 0 rgba(255,255,255,.72), 0 10px 24px rgba(15,23,42,.08)!important;
}

.invoice-wrap .copy-text{
	position:absolute!important;
	left:-9999px!important;
	top:auto!important;
	width:1px!important;
	height:1px!important;
	overflow:hidden!important;
}

.invoice-wrap .btn,
.invoice-wrap button{
	min-height:44px!important;
	border-radius:14px!important;
}

.invoice-wrap .muted b,
.invoice-wrap p b{
	font-weight:950!important;
}

.vg-public-status-badge-v2{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-height:28px!important;
	padding:5px 11px!important;
	border-radius:999px!important;
	font-size:12px!important;
	font-weight:950!important;
	line-height:1!important;
	letter-spacing:.01em!important;
	vertical-align:middle!important;
}

.vg-public-status-paid-v2,
.vg-public-status-success-v2{
	background:rgba(34,197,94,.13)!important;
	border:1px solid rgba(34,197,94,.28)!important;
	color:#15803d!important;
}

.vg-public-status-pending-v2{
	background:rgba(245,158,11,.14)!important;
	border:1px solid rgba(245,158,11,.32)!important;
	color:#b45309!important;
}

.vg-public-status-failed-v2{
	background:rgba(239,68,68,.12)!important;
	border:1px solid rgba(239,68,68,.28)!important;
	color:#b91c1c!important;
}

html[data-theme="dark"] .vg-public-status-paid-v2,
html[data-theme="dark"] .vg-public-status-success-v2{
	color:#86efac!important;
}

html[data-theme="dark"] .vg-public-status-pending-v2{
	color:#fcd34d!important;
}

html[data-theme="dark"] .vg-public-status-failed-v2{
	color:#fca5a5!important;
}

@media(max-width:760px){
	.invoice-wrap{
		max-width:none!important;
	}

	.invoice-wrap .card,
	.invoice-wrap .pay-card{
		padding:14px!important;
		border-radius:18px!important;
	}

	.invoice-wrap .voucher-box::before{
		font-size:11px;
		padding:5px 10px;
		margin-bottom:8px;
	}

	.invoice-wrap .btn,
	.invoice-wrap button{
		width:100%!important;
		min-height:42px!important;
		font-size:14px!important;
	}

	.vg-public-status-badge-v2{
		min-height:26px!important;
		font-size:11px!important;
		padding:5px 9px!important;
	}
}
</style>
<script>
/* PUBLIC_INVOICE_STATUS_POLISH_V2 */
(function(){
	function ready(fn){
		if(document.readyState === "loading"){
			document.addEventListener("DOMContentLoaded", fn);
		}else{
			fn();
		}
	}

	function classify(text){
		var t = String(text || "").trim().toLowerCase();
		if(!t) return "";

		if(t === "paid" || t === "lunas" || t === "success" || t === "sukses"){
			return "vg-public-status-paid-v2";
		}
		if(t === "pending" || t === "unpaid" || t === "menunggu"){
			return "vg-public-status-pending-v2";
		}
		if(t === "failed" || t === "gagal" || t === "expired"){
			return "vg-public-status-failed-v2";
		}
		return "";
	}

	ready(function(){
		document.querySelectorAll(".invoice-wrap b, .invoice-wrap strong, .invoice-wrap span").forEach(function(el){
			var cls = classify(el.textContent);
			if(!cls) return;

			el.classList.add("vg-public-status-badge-v2", cls);
		});
	});
})();
</script>







<style>
/* OWNER_BILLING_PAY_PHASE2_UI_V1 */
.owner-billing-pay-v1{
	display:grid!important;
	gap:16px!important;
	max-width:920px!important;
	margin:0 auto!important;
}

.owner-billing-pay-hero-v1{
	display:flex!important;
	justify-content:space-between!important;
	align-items:flex-start!important;
	gap:18px!important;
	padding:20px!important;
	border-radius:24px!important;
	border:1px solid rgba(148,163,184,.22)!important;
	background:
		radial-gradient(circle at 0% 0%, rgba(56,189,248,.22), transparent 34%),
		radial-gradient(circle at 100% 0%, rgba(124,58,237,.18), transparent 34%),
		rgba(255,255,255,.72)!important;
	box-shadow:0 22px 56px rgba(15,55,80,.14)!important;
}

html[data-theme="dark"] .owner-billing-pay-hero-v1{
	background:
		radial-gradient(circle at 0% 0%, rgba(56,189,248,.16), transparent 34%),
		radial-gradient(circle at 100% 0%, rgba(124,58,237,.18), transparent 34%),
		rgba(15,23,42,.76)!important;
	box-shadow:0 22px 60px rgba(0,0,0,.32)!important;
}

.owner-billing-kicker-v1{
	display:inline-flex!important;
	padding:6px 11px!important;
	border-radius:999px!important;
	background:rgba(37,99,235,.12)!important;
	color:#2563eb!important;
	font-size:12px!important;
	font-weight:950!important;
}

html[data-theme="dark"] .owner-billing-kicker-v1{
	color:#bfdbfe!important;
	background:rgba(96,165,250,.14)!important;
}

.owner-billing-pay-hero-v1 h1{
	margin:8px 0 6px!important;
	font-size:34px!important;
	font-weight:950!important;
	letter-spacing:-.05em!important;
	line-height:1.05!important;
}

.owner-billing-pay-hero-v1 p{
	margin:0!important;
	font-weight:750!important;
	opacity:.78!important;
	line-height:1.45!important;
}

.owner-billing-status-v1{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:110px!important;
	padding:9px 13px!important;
	border-radius:999px!important;
	background:rgba(34,197,94,.14)!important;
	border:1px solid rgba(34,197,94,.28)!important;
	color:#15803d!important;
	font-size:13px!important;
	font-weight:950!important;
	text-transform:capitalize!important;
}

html[data-theme="dark"] .owner-billing-status-v1{
	color:#86efac!important;
}

.owner-billing-pay-grid-v1{
	display:grid!important;
	grid-template-columns:repeat(3,minmax(0,1fr))!important;
	gap:14px!important;
}

.owner-billing-pay-grid-v1 .card{
	min-height:130px!important;
}

.owner-billing-pay-card-v1{
	display:grid!important;
	gap:12px!important;
}

.owner-billing-pay-card-v1 h2{
	margin:0!important;
	font-size:24px!important;
	font-weight:950!important;
	letter-spacing:-.04em!important;
}

.owner-billing-alert-v1{
	padding:12px 13px!important;
	border-radius:16px!important;
	border:1px solid rgba(245,158,11,.32)!important;
	background:rgba(245,158,11,.12)!important;
	color:#92400e!important;
	font-size:13px!important;
	font-weight:800!important;
	line-height:1.45!important;
}

html[data-theme="dark"] .owner-billing-alert-v1{
	color:#fcd34d!important;
	background:rgba(245,158,11,.10)!important;
}

@media(max-width:760px){
	.owner-billing-pay-v1{
		gap:12px!important;
	}

	.owner-billing-pay-hero-v1{
		display:grid!important;
		padding:16px!important;
		border-radius:20px!important;
	}

	.owner-billing-pay-hero-v1 h1{
		font-size:25px!important;
	}

	.owner-billing-pay-grid-v1{
		grid-template-columns:1fr!important;
		gap:10px!important;
	}

	.owner-billing-pay-grid-v1 .card{
		min-height:auto!important;
	}

	.owner-billing-pay-card-v1 .actions{
		display:grid!important;
		grid-template-columns:1fr!important;
	}

	.owner-billing-pay-card-v1 .btn{
		width:100%!important;
	}
}
</style>


<style>
/* OWNER_BILLING_PAY_CONTRAST_FIX_V2 */
html[data-theme="light"] .owner-billing-pay-v1,
html[data-theme="light"] .owner-billing-pay-v1 .card,
html[data-theme="light"] .owner-billing-pay-v1 .owner-billing-pay-card-v1{
	color:#0f172a!important;
}

html[data-theme="light"] .owner-billing-pay-v1 .card{
	background:rgba(255,255,255,.90)!important;
	border:1px solid rgba(255,255,255,.78)!important;
	box-shadow:
		0 18px 42px rgba(15,55,80,.14),
		inset 0 1px 0 rgba(255,255,255,.90)!important;
}

html[data-theme="light"] .owner-billing-pay-v1 .card b,
html[data-theme="light"] .owner-billing-pay-v1 .card h1,
html[data-theme="light"] .owner-billing-pay-v1 .card h2,
html[data-theme="light"] .owner-billing-pay-v1 .card h3{
	color:#0f172a!important;
	font-weight:950!important;
}

html[data-theme="light"] .owner-billing-pay-v1 .card .price{
	color:#0f172a!important;
	text-shadow:none!important;
	font-weight:950!important;
}

html[data-theme="light"] .owner-billing-pay-v1 .card p,
html[data-theme="light"] .owner-billing-pay-v1 .card .muted{
	color:#475569!important;
	opacity:1!important;
	font-weight:750!important;
}

html[data-theme="light"] .owner-billing-pay-card-v1{
	background:rgba(255,255,255,.88)!important;
	color:#0f172a!important;
}

html[data-theme="light"] .owner-billing-pay-card-v1 h2{
	color:#0f172a!important;
}

html[data-theme="light"] .owner-billing-pay-card-v1 p{
	color:#334155!important;
	opacity:1!important;
}

html[data-theme="dark"] .owner-billing-pay-v1 .card .price,
html[data-theme="dark"] .owner-billing-pay-v1 .card b,
html[data-theme="dark"] .owner-billing-pay-card-v1 h2{
	color:#f8fafc!important;
}

html[data-theme="dark"] .owner-billing-pay-v1 .card p,
html[data-theme="dark"] .owner-billing-pay-v1 .card .muted,
html[data-theme="dark"] .owner-billing-pay-card-v1 p{
	color:#cbd5e1!important;
	opacity:1!important;
}
</style>


<style>
/* OWNER_BILLING_HISTORY_PHASE3C_V1 */
.owner-billing-history-v1{
	display:grid!important;
	gap:12px!important;
	background:rgba(255,255,255,.88)!important;
	border:1px solid rgba(255,255,255,.76)!important;
}

html[data-theme="dark"] .owner-billing-history-v1{
	background:rgba(15,23,42,.76)!important;
	border-color:rgba(148,163,184,.18)!important;
}

.owner-billing-history-head-v1{
	display:flex!important;
	align-items:flex-start!important;
	justify-content:space-between!important;
	gap:12px!important;
}

.owner-billing-history-head-v1 h2{
	margin:0!important;
	font-size:23px!important;
	font-weight:950!important;
	letter-spacing:-.04em!important;
	color:#0f172a!important;
}

html[data-theme="dark"] .owner-billing-history-head-v1 h2{
	color:#f8fafc!important;
}

.owner-billing-history-list-v1{
	display:grid!important;
	gap:9px!important;
}

.owner-billing-history-item-v1{
	display:grid!important;
	grid-template-columns:minmax(0,1.5fr) minmax(110px,.45fr) auto!important;
	align-items:center!important;
	gap:12px!important;
	padding:12px!important;
	border-radius:16px!important;
	border:1px solid rgba(148,163,184,.20)!important;
	background:rgba(248,250,252,.68)!important;
}

html[data-theme="dark"] .owner-billing-history-item-v1{
	background:rgba(2,6,23,.28)!important;
	border-color:rgba(148,163,184,.16)!important;
}

.owner-billing-history-main-v1{
	display:grid!important;
	gap:3px!important;
	min-width:0!important;
}

.owner-billing-history-main-v1 b{
	color:#0f172a!important;
	font-weight:950!important;
	overflow-wrap:anywhere!important;
}

html[data-theme="dark"] .owner-billing-history-main-v1 b{
	color:#f8fafc!important;
}

.owner-billing-history-main-v1 span,
.owner-billing-history-meta-v1 small{
	color:#64748b!important;
	font-size:12px!important;
	font-weight:750!important;
}

html[data-theme="dark"] .owner-billing-history-main-v1 span,
html[data-theme="dark"] .owner-billing-history-meta-v1 small{
	color:#cbd5e1!important;
}

.owner-billing-history-meta-v1{
	display:grid!important;
	gap:2px!important;
	text-align:right!important;
}

.owner-billing-history-meta-v1 strong{
	color:#0f172a!important;
	font-size:18px!important;
	font-weight:950!important;
	letter-spacing:-.03em!important;
}

html[data-theme="dark"] .owner-billing-history-meta-v1 strong{
	color:#f8fafc!important;
}

.owner-billing-history-action-v1{
	display:flex!important;
	align-items:center!important;
	justify-content:flex-end!important;
	gap:8px!important;
}

.owner-billing-mini-status-v1{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:72px!important;
	padding:7px 10px!important;
	border-radius:999px!important;
	font-size:12px!important;
	font-weight:950!important;
	text-transform:capitalize!important;
}

.owner-billing-mini-status-v1.ok{
	background:rgba(34,197,94,.14)!important;
	border:1px solid rgba(34,197,94,.28)!important;
	color:#15803d!important;
}

.owner-billing-mini-status-v1.warn{
	background:rgba(245,158,11,.14)!important;
	border:1px solid rgba(245,158,11,.32)!important;
	color:#b45309!important;
}

.owner-billing-mini-status-v1.bad{
	background:rgba(239,68,68,.12)!important;
	border:1px solid rgba(239,68,68,.30)!important;
	color:#b91c1c!important;
}

.owner-billing-empty-v1{
	padding:14px!important;
	border-radius:16px!important;
	border:1px dashed rgba(148,163,184,.38)!important;
	color:#64748b!important;
	font-weight:800!important;
	background:rgba(248,250,252,.52)!important;
}

html[data-theme="dark"] .owner-billing-empty-v1{
	color:#cbd5e1!important;
	background:rgba(2,6,23,.22)!important;
}

@media(max-width:760px){
	.owner-billing-history-item-v1{
		grid-template-columns:1fr!important;
		align-items:start!important;
		gap:9px!important;
		padding:11px!important;
	}

	.owner-billing-history-meta-v1{
		text-align:left!important;
	}

	.owner-billing-history-action-v1{
		justify-content:space-between!important;
	}

	.owner-billing-history-action-v1 .btn{
		min-height:36px!important;
		padding:8px 12px!important;
	}
}
</style>


<style>
/* CLIENT_BILLING_SUSPENDED_UI_V1 */
.billing-suspended-page-v1{
	min-height:68vh!important;
	display:grid!important;
	place-items:center!important;
	padding:16px!important;
}

.billing-suspended-card-v1{
	width:min(720px,100%)!important;
	position:relative!important;
	overflow:hidden!important;
	border-radius:28px!important;
	padding:26px!important;
	text-align:center!important;
	border:1px solid rgba(148,163,184,.22)!important;
	background:
		radial-gradient(circle at 0% 0%, rgba(239,68,68,.18), transparent 34%),
		radial-gradient(circle at 100% 0%, rgba(245,158,11,.16), transparent 34%),
		rgba(255,255,255,.90)!important;
	box-shadow:0 26px 70px rgba(15,23,42,.18)!important;
}

html[data-theme="dark"] .billing-suspended-card-v1{
	background:
		radial-gradient(circle at 0% 0%, rgba(239,68,68,.16), transparent 34%),
		radial-gradient(circle at 100% 0%, rgba(245,158,11,.14), transparent 34%),
		rgba(15,23,42,.82)!important;
	box-shadow:0 26px 80px rgba(0,0,0,.38)!important;
}

.billing-suspended-icon-v1{
	width:62px!important;
	height:62px!important;
	margin:0 auto 12px!important;
	display:grid!important;
	place-items:center!important;
	border-radius:22px!important;
	background:rgba(239,68,68,.12)!important;
	border:1px solid rgba(239,68,68,.26)!important;
	font-size:30px!important;
}

.billing-suspended-card-v1 h1{
	margin:0 0 8px!important;
	font-size:34px!important;
	font-weight:950!important;
	letter-spacing:-.05em!important;
	color:#0f172a!important;
}

html[data-theme="dark"] .billing-suspended-card-v1 h1{
	color:#f8fafc!important;
}

.billing-suspended-lead-v1{
	max-width:560px!important;
	margin:0 auto 18px!important;
	color:#475569!important;
	font-weight:800!important;
	line-height:1.5!important;
}

html[data-theme="dark"] .billing-suspended-lead-v1{
	color:#cbd5e1!important;
}

.billing-suspended-grid-v1{
	display:grid!important;
	grid-template-columns:repeat(5,minmax(0,1fr))!important;
	gap:10px!important;
	margin:18px 0!important;
	text-align:left!important;
}

.billing-suspended-grid-v1 div{
	padding:12px!important;
	border-radius:16px!important;
	border:1px solid rgba(148,163,184,.22)!important;
	background:rgba(248,250,252,.72)!important;
	min-width:0!important;
}

html[data-theme="dark"] .billing-suspended-grid-v1 div{
	background:rgba(2,6,23,.24)!important;
	border-color:rgba(148,163,184,.16)!important;
}

.billing-suspended-grid-v1 b{
	display:block!important;
	color:#64748b!important;
	font-size:11px!important;
	font-weight:950!important;
	margin-bottom:4px!important;
}

.billing-suspended-grid-v1 span{
	display:block!important;
	color:#0f172a!important;
	font-weight:950!important;
	overflow-wrap:anywhere!important;
}

html[data-theme="dark"] .billing-suspended-grid-v1 span{
	color:#f8fafc!important;
}

.billing-suspended-actions-v1{
	display:flex!important;
	flex-wrap:wrap!important;
	justify-content:center!important;
	gap:10px!important;
	margin-top:14px!important;
}

.billing-suspended-note-v1{
	margin-top:16px!important;
	padding:12px 14px!important;
	border-radius:16px!important;
	border:1px solid rgba(245,158,11,.32)!important;
	background:rgba(245,158,11,.12)!important;
	color:#92400e!important;
	font-size:13px!important;
	font-weight:800!important;
	line-height:1.45!important;
	text-align:left!important;
}

html[data-theme="dark"] .billing-suspended-note-v1{
	color:#fcd34d!important;
	background:rgba(245,158,11,.10)!important;
}

@media(max-width:760px){
	.billing-suspended-page-v1{
		min-height:70vh!important;
		padding:10px!important;
	}

	.billing-suspended-card-v1{
		border-radius:22px!important;
		padding:18px!important;
	}

	.billing-suspended-card-v1 h1{
		font-size:26px!important;
	}

	.billing-suspended-grid-v1{
		grid-template-columns:repeat(2,minmax(0,1fr))!important;
	}

	.billing-suspended-actions-v1{
		display:grid!important;
		grid-template-columns:1fr!important;
	}

	.billing-suspended-actions-v1 .btn{
		width:100%!important;
	}
}
</style>


<style>


/* PUBLIC_PACKAGE_LAYOUT_CLEAN_V1 */
/* Consolidated package page cards: mobile 2 kolom, tablet 2 kolom, desktop 3 kolom. */

.grid:has(.buy-btn){
  display:grid!important;
  grid-template-columns:repeat(2,minmax(0,1fr))!important;
  gap:10px!important;
  align-items:stretch!important;
  width:100%!important;
}

.grid:has(.buy-btn) > .card{
  position:relative!important;
  overflow:hidden!important;
  isolation:isolate!important;
  min-width:0!important;
  padding:12px!important;
  border-radius:18px!important;
}

.grid:has(.buy-btn) > .card::before{
  content:"";
  position:absolute;
  inset:0;
  z-index:-1;
  background:
    radial-gradient(circle at 18% 12%, rgba(56,189,248,.20), transparent 32%),
    radial-gradient(circle at 92% 0%, rgba(168,85,247,.18), transparent 34%);
  opacity:.88;
  pointer-events:none;
}

.grid:has(.buy-btn) > .card::after{
  content:"";
  position:absolute;
  inset:auto 16px 0 16px;
  height:1px;
  background:linear-gradient(90deg, transparent, rgba(255,255,255,.70), transparent);
  opacity:.8;
  pointer-events:none;
}

html[data-theme="light"] .grid:has(.buy-btn) > .card{
  background:rgba(255,255,255,.88)!important;
  border:1px solid rgba(255,255,255,.78)!important;
  box-shadow:
    0 18px 42px rgba(15,55,80,.14),
    inset 0 1px 0 rgba(255,255,255,.90)!important;
}

html[data-theme="dark"] .grid:has(.buy-btn) > .card{
  background:rgba(15,23,42,.72)!important;
  border:1px solid rgba(148,163,184,.18)!important;
  box-shadow:
    0 22px 56px rgba(0,0,0,.28),
    inset 0 1px 0 rgba(255,255,255,.06)!important;
}

.grid:has(.buy-btn) > .card h3{
  font-size:18px!important;
  line-height:1.25!important;
  margin:0 0 5px!important;
  letter-spacing:-.035em!important;
}

.grid:has(.buy-btn) > .card .price{
  display:inline-flex!important;
  align-items:center!important;
  justify-content:flex-start!important;
  gap:6px!important;
  font-size:27px!important;
  line-height:1.15!important;
  margin:2px 0 6px!important;
  letter-spacing:-.055em!important;
}

.grid:has(.buy-btn) > .card .muted,
.grid:has(.buy-btn) > .card p{
  font-size:12px!important;
  line-height:1.25!important;
  font-weight:750!important;
  margin:0 0 10px!important;
}

.grid:has(.buy-btn) > .card .buy-btn{
  position:relative!important;
  overflow:hidden!important;
  width:100%!important;
  min-height:38px!important;
  padding:9px 10px!important;
  border:0!important;
  border-radius:12px!important;
  font-size:13px!important;
  background:
    linear-gradient(135deg,#2563eb 0%,#7c3aed 58%,#38bdf8 120%)!important;
  box-shadow:
    0 14px 28px rgba(37,99,235,.24),
    inset 0 1px 0 rgba(255,255,255,.22)!important;
}

.grid:has(.buy-btn) > .card .buy-btn::before{
  content:"";
  position:absolute;
  inset:-60% auto auto -45%;
  width:42%;
  height:220%;
  background:linear-gradient(90deg, transparent, rgba(255,255,255,.42), transparent);
  transform:rotate(22deg) translateX(-140%);
  transition:transform .62s cubic-bezier(.22,1,.36,1);
  pointer-events:none;
}

.grid:has(.buy-btn) > .card .buy-btn:hover::before{
  transform:rotate(22deg) translateX(450%);
}

.grid:has(.buy-btn) > .card:hover{
  border-color:rgba(96,165,250,.38)!important;
}

@media(min-width:761px){
  .grid:has(.buy-btn){
    gap:14px!important;
  }

  .grid:has(.buy-btn) > .card{
    padding:15px!important;
  }

  .grid:has(.buy-btn) > .card h3{
    font-size:17px!important;
  }

  .grid:has(.buy-btn) > .card .price{
    font-size:24px!important;
  }
}

@media(min-width:1024px){
  html[data-theme="light"] .wrap > .grid:has(> .card .buy-btn),
  html[data-theme="dark"] .wrap > .grid:has(> .card .buy-btn){
    grid-template-columns:repeat(3,minmax(0,1fr))!important;
    max-width:1120px!important;
    margin-left:auto!important;
    margin-right:auto!important;
  }

  html[data-theme="light"] .wrap > .grid:has(> .card .buy-btn) > .card h3,
  html[data-theme="dark"] .wrap > .grid:has(> .card .buy-btn) > .card h3{
    font-size:20px!important;
  }

  html[data-theme="light"] .wrap > .grid:has(> .card .buy-btn) > .card .price,
  html[data-theme="dark"] .wrap > .grid:has(> .card .buy-btn) > .card .price{
    font-size:31px!important;
  }

  html[data-theme="light"] .wrap > .grid:has(> .card .buy-btn) > .card .buy-btn,
  html[data-theme="dark"] .wrap > .grid:has(> .card .buy-btn) > .card .buy-btn{
    min-height:42px!important;
  }
}

@media(max-width:360px){
  .grid:has(.buy-btn){
    gap:8px!important;
  }

  .grid:has(.buy-btn) > .card{
    padding:10px!important;
    border-radius:15px!important;
  }

  .grid:has(.buy-btn) > .card h3{
    font-size:14px!important;
  }

  .grid:has(.buy-btn) > .card .price{
    font-size:18px!important;
  }

  .grid:has(.buy-btn) > .card .buy-btn{
    font-size:12px!important;
    padding:8px!important;
  }
}
/* /PUBLIC_PACKAGE_LAYOUT_CLEAN_V1 */

/* PUBLIC_HOME_DARK_SYNC_LIGHT_LAYOUT_V1 */

/*
  Khusus halaman public root:
  - posisi dark disamakan dengan light
  - hero tetap center 900px
  - grid router tetap 1190px
  - card dark dibuat glass/dark, bukan putih
  Selector dibuat lebih spesifik supaya menang dari override lama.
*/

html[data-theme="light"] body .wrap,
html[data-theme="dark"] body .wrap{
  width:100%!important;
  max-width:1190px!important;
  margin-left:auto!important;
  margin-right:auto!important;
}

html[data-theme="light"] body .wrap > .public-top,
html[data-theme="dark"] body .wrap > .public-top{
  width:100%!important;
  max-width:1190px!important;
  margin-left:auto!important;
  margin-right:auto!important;
}

html[data-theme="light"] body .wrap > .hero,
html[data-theme="dark"] body .wrap > .hero{
  width:100%!important;
  max-width:900px!important;
  margin:0 auto 18px!important;
  padding:26px!important;
  border-radius:26px!important;
}

html[data-theme="light"] body .wrap > .hero h1,
html[data-theme="dark"] body .wrap > .hero h1{
  font-size:32px!important;
  line-height:1.05!important;
  font-weight:950!important;
  letter-spacing:-.05em!important;
  margin:0!important;
}

html[data-theme="light"] body .wrap > .hero p,
html[data-theme="dark"] body .wrap > .hero p{
  font-size:18px!important;
  line-height:1.35!important;
  font-weight:500!important;
  margin:6px 0 0!important;
}

html[data-theme="light"] body .wrap > .grid,
html[data-theme="dark"] body .wrap > .grid{
  width:100%!important;
  max-width:1190px!important;
  margin:14px auto 0!important;
  display:grid!important;
  grid-template-columns:repeat(2,minmax(0,1fr))!important;
  gap:14px!important;
}

html[data-theme="light"] body .wrap > .grid > a.card,
html[data-theme="dark"] body .wrap > .grid > a.card{
  min-height:156px!important;
  padding:24px!important;
  border-radius:22px!important;
  margin:0!important;
  display:block!important;
  transform:none!important;
}

html[data-theme="light"] body .wrap > .grid > a.card h3,
html[data-theme="dark"] body .wrap > .grid > a.card h3{
  font-size:30px!important;
  line-height:1.14!important;
  font-weight:950!important;
  letter-spacing:-.035em!important;
  margin:0 0 22px!important;
}

html[data-theme="light"] body .wrap > .grid > a.card p,
html[data-theme="light"] body .wrap > .grid > a.card .muted,
html[data-theme="dark"] body .wrap > .grid > a.card p,
html[data-theme="dark"] body .wrap > .grid > a.card .muted{
  font-size:20px!important;
  line-height:1.35!important;
  font-weight:800!important;
  margin:0!important;
  opacity:1!important;
}

/* Warna light tetap seperti sekarang */
html[data-theme="light"] body .wrap > .grid > a.card{
  color:#172033!important;
  background:rgba(255,255,255,.78)!important;
  border:1px solid rgba(255,255,255,.76)!important;
  box-shadow:0 24px 70px rgba(14,116,144,.18)!important;
  backdrop-filter:blur(14px) saturate(130%)!important;
  -webkit-backdrop-filter:blur(14px) saturate(130%)!important;
}

html[data-theme="light"] body .wrap > .grid > a.card h3{
  color:#172033!important;
}

html[data-theme="light"] body .wrap > .grid > a.card p,
html[data-theme="light"] body .wrap > .grid > a.card .muted{
  color:#475569!important;
}

/* Warna dark: dark glass/transparan, bukan putih */
html[data-theme="dark"] body .wrap > .grid > a.card{
  color:#f8fafc!important;
  background:
    linear-gradient(135deg,rgba(15,23,42,.72),rgba(30,41,59,.42))!important;
  border:1px solid rgba(148,163,184,.26)!important;
  box-shadow:
    0 24px 70px rgba(0,0,0,.34),
    inset 0 1px 0 rgba(255,255,255,.08)!important;
  backdrop-filter:blur(18px) saturate(140%)!important;
  -webkit-backdrop-filter:blur(18px) saturate(140%)!important;
}

html[data-theme="dark"] body .wrap > .grid > a.card h3{
  color:#f8fafc!important;
  text-shadow:none!important;
}

html[data-theme="dark"] body .wrap > .grid > a.card p,
html[data-theme="dark"] body .wrap > .grid > a.card .muted{
  color:#cbd5e1!important;
  text-shadow:none!important;
}

@media(max-width:760px){
  html[data-theme="light"] body .wrap,
  html[data-theme="dark"] body .wrap{
    max-width:480px!important;
    padding:12px!important;
  }

  html[data-theme="light"] body .wrap > .hero,
  html[data-theme="dark"] body .wrap > .hero{
    max-width:100%!important;
    padding:18px!important;
    border-radius:22px!important;
    margin-bottom:12px!important;
  }

  html[data-theme="light"] body .wrap > .hero h1,
  html[data-theme="dark"] body .wrap > .hero h1{
    font-size:25px!important;
  }

  html[data-theme="light"] body .wrap > .hero p,
  html[data-theme="dark"] body .wrap > .hero p{
    font-size:13px!important;
  }

  html[data-theme="light"] body .wrap > .grid,
  html[data-theme="dark"] body .wrap > .grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:10px!important;
  }

  html[data-theme="light"] body .wrap > .grid > a.card,
  html[data-theme="dark"] body .wrap > .grid > a.card{
    min-height:118px!important;
    padding:12px!important;
    border-radius:18px!important;
  }

  html[data-theme="light"] body .wrap > .grid > a.card h3,
  html[data-theme="dark"] body .wrap > .grid > a.card h3{
    font-size:19px!important;
    line-height:1.16!important;
    margin-bottom:10px!important;
  }

  html[data-theme="light"] body .wrap > .grid > a.card p,
  html[data-theme="light"] body .wrap > .grid > a.card .muted,
  html[data-theme="dark"] body .wrap > .grid > a.card p,
  html[data-theme="dark"] body .wrap > .grid > a.card .muted{
    font-size:14px!important;
  }
}

/* /PUBLIC_HOME_DARK_SYNC_LIGHT_LAYOUT_V1 */




/* PUBLIC_ROUTER_B_TITLE_SIZE_V1 */

/* Judul router public pakai tag <b>, bukan h3 */
html[data-theme="light"] body .wrap > .grid > a.card > b,
html[data-theme="dark"] body .wrap > .grid > a.card > b,
html[data-theme="light"] body .wrap > .grid > .card > b,
html[data-theme="dark"] body .wrap > .grid > .card > b{
  display:block!important;
  font-size:30px!important;
  line-height:1.14!important;
  font-weight:950!important;
  letter-spacing:-.035em!important;
  margin:0 0 22px!important;
}

html[data-theme="light"] body .wrap > .grid > a.card > b,
html[data-theme="light"] body .wrap > .grid > .card > b{
  color:#172033!important;
}

html[data-theme="dark"] body .wrap > .grid > a.card > b,
html[data-theme="dark"] body .wrap > .grid > .card > b{
  color:#f8fafc!important;
  text-shadow:none!important;
}

@media(max-width:760px){
  html[data-theme="light"] body .wrap > .grid > a.card > b,
  html[data-theme="dark"] body .wrap > .grid > a.card > b,
  html[data-theme="light"] body .wrap > .grid > .card > b,
  html[data-theme="dark"] body .wrap > .grid > .card > b{
    font-size:19px!important;
    line-height:1.16!important;
    margin-bottom:10px!important;
  }
}

/* /PUBLIC_ROUTER_B_TITLE_SIZE_V1 */


/* VOUCHERGO_DARK_PALETTE_V1 */
html[data-theme="dark"]{
  --vg-bg: #04142B;
  --vg-surface: #231832;
  --vg-surface-2: #512438;
  --vg-border: rgba(88,72,104,.62);
  --vg-primary: #F6762F;
  --vg-primary-2: #BD4C3E;
  --vg-warning: #EEC207;
  --vg-muted: #584868;
  --vg-text: #F8FAFC;
  --vg-text-muted: rgba(248,250,252,.72);
}

html[data-theme="dark"] body:not(.admin-body){
  background:
    radial-gradient(circle at 18% 12%, rgba(246,118,47,.22), transparent 28%),
    radial-gradient(circle at 88% 20%, rgba(189,76,62,.30), transparent 30%),
    radial-gradient(circle at 72% 82%, rgba(238,194,7,.13), transparent 24%),
    linear-gradient(145deg, #04142B 0%, #231832 48%, #04142B 100%) !important;
  color: var(--vg-text) !important;
}

html[data-theme="dark"] .wrap,
html[data-theme="dark"] .public-wrap,
html[data-theme="dark"] .public-page,
html[data-theme="dark"] .public-shell,
html[data-theme="dark"] main{
  color: var(--vg-text) !important;
}

html[data-theme="dark"] .hero,
html[data-theme="dark"] .public-hero{
  background:
    linear-gradient(135deg, rgba(35,24,50,.86), rgba(81,36,56,.72)),
    radial-gradient(circle at top right, rgba(246,118,47,.26), transparent 38%) !important;
  border: 1px solid rgba(88,72,104,.55) !important;
  box-shadow: 0 18px 45px rgba(4,20,43,.42) !important;
}

html[data-theme="dark"] .wrap > .card,
html[data-theme="dark"] .public-card,
html[data-theme="dark"] .router-card,
html[data-theme="dark"] .package-card,
html[data-theme="dark"] .wrap > .grid > .card{
  background:
    linear-gradient(145deg, rgba(35,24,50,.94), rgba(81,36,56,.72)) !important;
  border: 1px solid rgba(88,72,104,.70) !important;
  box-shadow: 0 18px 42px rgba(4,20,43,.46) !important;
  color: var(--vg-text) !important;
}

html[data-theme="dark"] .wrap > .card .card{
  background: rgba(4,20,43,.58) !important;
  border: 1px solid rgba(88,72,104,.56) !important;
}

html[data-theme="dark"] h1,
html[data-theme="dark"] h2,
html[data-theme="dark"] h3,
html[data-theme="dark"] h4,
html[data-theme="dark"] strong,
html[data-theme="dark"] b,
html[data-theme="dark"] label,
html[data-theme="dark"] .title,
html[data-theme="dark"] .name,
html[data-theme="dark"] .price,
html[data-theme="dark"] .amount{
  color: var(--vg-text) !important;
}

html[data-theme="dark"] p,
html[data-theme="dark"] small,
html[data-theme="dark"] .muted,
html[data-theme="dark"] .desc{
  color: var(--vg-text-muted) !important;
}

html[data-theme="dark"] input,
html[data-theme="dark"] select,
html[data-theme="dark"] textarea{
  background: rgba(4,20,43,.70) !important;
  border: 1px solid rgba(88,72,104,.78) !important;
  color: var(--vg-text) !important;
}

html[data-theme="dark"] input:focus,
html[data-theme="dark"] select:focus,
html[data-theme="dark"] textarea:focus{
  border-color: rgba(246,118,47,.86) !important;
  box-shadow: 0 0 0 3px rgba(246,118,47,.18) !important;
}

html[data-theme="dark"] .btn,
html[data-theme="dark"] button,
html[data-theme="dark"] .buy-btn,
html[data-theme="dark"] .btn.pri{
  background: linear-gradient(135deg, #F6762F, #BD4C3E) !important;
  color: #fff !important;
  border: 1px solid rgba(246,118,47,.76) !important;
  box-shadow: 0 12px 28px rgba(246,118,47,.22) !important;
}

html[data-theme="dark"] .btn:not(.pri):not(.buy-btn),
html[data-theme="dark"] a.btn:not(.pri):not(.buy-btn){
  background: rgba(81,36,56,.72) !important;
  color: var(--vg-text) !important;
  border: 1px solid rgba(189,76,62,.58) !important;
  box-shadow: none !important;
}

html[data-theme="dark"] .badge,
html[data-theme="dark"] .status,
html[data-theme="dark"] .tag{
  background: rgba(238,194,7,.15) !important;
  color: #FFE36A !important;
  border-color: rgba(238,194,7,.35) !important;
}
/* /VOUCHERGO_DARK_PALETTE_V1 */

/* VOUCHERGO_PUBLIC_BUTTON_SIMPLE_HOVER_V1 */
html[data-theme="dark"] body:not(.admin-body) .wrap .grid:has(.buy-btn) > .card .buy-btn,
html[data-theme="dark"] body:not(.admin-body) .grid:has(.buy-btn) > .card .buy-btn,
html[data-theme="dark"] body:not(.admin-body) a.btn.pri.buy-btn,
html[data-theme="dark"] body:not(.admin-body) .buy-btn{
  position:relative!important;
  overflow:hidden!important;
  isolation:isolate!important;
  background:#F6762F!important;
  background-image:none!important;
  color:#fff!important;
  border:1px solid rgba(246,118,47,.72)!important;
  box-shadow:0 10px 22px rgba(246,118,47,.22)!important;
  transform:none!important;
  filter:none!important;
  transition:
    background-color .16s ease,
    border-color .16s ease,
    box-shadow .16s ease,
    opacity .16s ease!important;
}

html[data-theme="dark"] body:not(.admin-body) .wrap .grid:has(.buy-btn) > .card .buy-btn::before,
html[data-theme="dark"] body:not(.admin-body) .wrap .grid:has(.buy-btn) > .card .buy-btn::after,
html[data-theme="dark"] body:not(.admin-body) .grid:has(.buy-btn) > .card .buy-btn::before,
html[data-theme="dark"] body:not(.admin-body) .grid:has(.buy-btn) > .card .buy-btn::after,
html[data-theme="dark"] body:not(.admin-body) a.btn.pri.buy-btn::before,
html[data-theme="dark"] body:not(.admin-body) a.btn.pri.buy-btn::after,
html[data-theme="dark"] body:not(.admin-body) .buy-btn::before,
html[data-theme="dark"] body:not(.admin-body) .buy-btn::after{
  content:none!important;
  display:none!important;
  opacity:0!important;
  animation:none!important;
  transition:none!important;
  transform:none!important;
}

html[data-theme="dark"] body:not(.admin-body) .wrap .grid:has(.buy-btn) > .card .buy-btn:hover,
html[data-theme="dark"] body:not(.admin-body) .grid:has(.buy-btn) > .card .buy-btn:hover,
html[data-theme="dark"] body:not(.admin-body) a.btn.pri.buy-btn:hover,
html[data-theme="dark"] body:not(.admin-body) .buy-btn:hover{
  background:#BD4C3E!important;
  background-image:none!important;
  border-color:rgba(189,76,62,.88)!important;
  box-shadow:0 10px 22px rgba(189,76,62,.24)!important;
  transform:none!important;
  filter:none!important;
}

html[data-theme="dark"] body:not(.admin-body) .wrap .grid:has(.buy-btn) > .card .buy-btn:active,
html[data-theme="dark"] body:not(.admin-body) .grid:has(.buy-btn) > .card .buy-btn:active,
html[data-theme="dark"] body:not(.admin-body) a.btn.pri.buy-btn:active,
html[data-theme="dark"] body:not(.admin-body) .buy-btn:active{
  opacity:.88!important;
  transform:none!important;
  filter:none!important;
}
/* /VOUCHERGO_PUBLIC_BUTTON_SIMPLE_HOVER_V1 */
</style>



<style>
/* PUBLIC_CHECKOUT_CHROME_GLITCH_FIX_V1 */
/* Stabilkan first paint Chrome di halaman order/checkout. */
.checkout-hero-v1,
.checkout-form-card-v1,
.checkout-form-card-v1 .notice{
	animation:none!important;
	transition:none!important;
	filter:none!important;
	-webkit-backdrop-filter:none!important;
	backdrop-filter:none!important;
	transform:none!important;
	will-change:auto!important;
	contain:paint;
}

html:not([data-theme="light"]) .checkout-hero-v1,
html[data-theme="dark"] .checkout-hero-v1,
html[data-admin-theme="dark"] .checkout-hero-v1{
	background:
		radial-gradient(circle at 18% 8%, rgba(56,189,248,.14), transparent 32%),
		radial-gradient(circle at 88% 14%, rgba(168,85,247,.13), transparent 34%),
		linear-gradient(135deg, rgba(18,18,42,.96), rgba(75,28,55,.86))!important;
	border:1px solid rgba(148,163,184,.16)!important;
	box-shadow:0 18px 48px rgba(2,6,23,.34)!important;
	color:#f9fafb!important;
}

html:not([data-theme="light"]) .checkout-hero-v1 h1,
html[data-theme="dark"] .checkout-hero-v1 h1,
html[data-admin-theme="dark"] .checkout-hero-v1 h1{
	color:#f9fafb!important;
	text-shadow:none!important;
}

html:not([data-theme="light"]) .checkout-hero-v1 p,
html[data-theme="dark"] .checkout-hero-v1 p,
html[data-admin-theme="dark"] .checkout-hero-v1 p{
	color:#e5e7eb!important;
	text-shadow:none!important;
}

html:not([data-theme="light"]) .checkout-form-card-v1,
html[data-theme="dark"] .checkout-form-card-v1,
html[data-admin-theme="dark"] .checkout-form-card-v1{
	background:
		radial-gradient(circle at 12% 0%, rgba(56,189,248,.11), transparent 30%),
		linear-gradient(145deg, rgba(46,22,52,.94), rgba(83,31,61,.84))!important;
	border:1px solid rgba(148,163,184,.18)!important;
	box-shadow:0 22px 58px rgba(2,6,23,.36)!important;
	color:#f9fafb!important;
}

html:not([data-theme="light"]) .checkout-form-card-v1 .notice,
html[data-theme="dark"] .checkout-form-card-v1 .notice,
html[data-admin-theme="dark"] .checkout-form-card-v1 .notice{
	background:#071126!important;
	border:1px solid rgba(148,163,184,.14)!important;
	box-shadow:none!important;
	color:#f9fafb!important;
}

html:not([data-theme="light"]) .checkout-form-card-v1 label,
html[data-theme="dark"] .checkout-form-card-v1 label,
html[data-admin-theme="dark"] .checkout-form-card-v1 label,
html:not([data-theme="light"]) .checkout-form-card-v1 h2,
html[data-theme="dark"] .checkout-form-card-v1 h2,
html[data-admin-theme="dark"] .checkout-form-card-v1 h2,
html:not([data-theme="light"]) .checkout-form-card-v1 .price,
html[data-theme="dark"] .checkout-form-card-v1 .price,
html[data-admin-theme="dark"] .checkout-form-card-v1 .price{
	color:#f9fafb!important;
}

html:not([data-theme="light"]) .checkout-form-card-v1 .muted,
html[data-theme="dark"] .checkout-form-card-v1 .muted,
html[data-admin-theme="dark"] .checkout-form-card-v1 .muted{
	color:#aeb8c8!important;
}

html:not([data-theme="light"]) .checkout-form-card-v1 input,
html[data-theme="dark"] .checkout-form-card-v1 input,
html[data-admin-theme="dark"] .checkout-form-card-v1 input{
	background:#0b1024!important;
	border-color:rgba(96,165,250,.24)!important;
	color:#f9fafb!important;
	box-shadow:none!important;
}

html:not([data-theme="light"]) .checkout-form-card-v1 input::placeholder,
html[data-theme="dark"] .checkout-form-card-v1 input::placeholder,
html[data-admin-theme="dark"] .checkout-form-card-v1 input::placeholder{
	color:#6b7280!important;
}

@media(max-width:720px){
	.checkout-hero-v1{
		border-radius:24px!important;
		padding:28px 24px!important;
	}

	.checkout-form-card-v1{
		border-radius:24px!important;
	}

	.checkout-form-card-v1 .notice{
		border-radius:20px!important;
	}
}
</style>

</body>
</html>`
