package main

import (
	"html/template"
	"net/http"
)

func page(w http.ResponseWriter, title string, body template.HTML) {
	const layout = `<!doctype html>
<html lang="id" data-admin-theme="light">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<title>{{.Title}}</title>
<script>
// VOUCHERGO_ADMIN_THEME_BOOTSTRAP_HEAD_V1
(function(){
        var keys = ["tuku-admin-theme", "vouchergo_theme", "theme", "tuku_theme", "admin-theme", "adminTheme", "tuku_admin_theme"];
        var root = document.documentElement;
        for(var i=0;i<keys.length;i++){
                try{
                        var v = localStorage.getItem(keys[i]);
                        if(v === "light" || v === "dark"){
                                root.setAttribute("data-admin-theme", v);
                                root.setAttribute("data-theme", v);
                                return;
                        }
                }catch(e){}
        }
})();
</script>
<style>{{.AdminCSS}}</style>

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

    document.querySelectorAll('.super-only,a[href="/admin/whatsapp"],a[href="/admin/settings"],a[href="/admin/users"],a[href="/admin/blocked-phones"]').forEach(function(el){
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
/* NAV_BG_FINAL_OVERRIDE_V1 */
.nav{
  background-color:#2a1530!important;
  background-image:
    radial-gradient(circle at 10% 0%, rgba(139,92,246,.18), transparent 34%),
    radial-gradient(circle at 90% 100%, rgba(244,63,94,.16), transparent 36%),
    linear-gradient(135deg, rgba(40,22,50,.98), rgba(66,27,56,.95) 58%, rgba(20,26,52,.96))!important;
  border:1px solid rgba(168,85,247,.22)!important;
  box-shadow:
    0 16px 42px rgba(0,0,0,.24),
    inset 0 1px 0 rgba(255,255,255,.06)!important;
}

.nav .menu a,
.admin-theme-toggle{
  background:linear-gradient(180deg, rgba(255,255,255,.085), rgba(255,255,255,.04))!important;
  border:1px solid rgba(255,255,255,.11)!important;
  color:#eef2ff!important;
}

.nav .menu a[href="/admin"]{
  background:linear-gradient(180deg, rgba(255,255,255,.085), rgba(255,255,255,.04))!important;
  border-color:rgba(255,255,255,.11)!important;
}

.nav .menu a:hover,
.nav .menu a[href="/admin"]:hover,
.admin-theme-toggle:hover{
  background:
    radial-gradient(circle at 18% 0%, rgba(56,189,248,.20), transparent 42%),
    linear-gradient(180deg, rgba(139,92,246,.24), rgba(244,63,94,.13))!important;
  border-color:rgba(196,181,253,.35)!important;
  color:#fff!important;
}

html[data-admin-theme="light"] .nav,
html[data-theme="light"] .nav{
  background-color:#ffffff!important;
  background-image:linear-gradient(135deg,#ffffff,#f1f5f9)!important;
  border-color:#dbe4f0!important;
}
/* /NAV_BG_FINAL_OVERRIDE_V1 */
</style>

</head>
<body>
<div class="wrap">

<div class="mobile-adminbar">
  <button class="mobile-menu-btn" type="button" onclick="toggleAdminDrawer()" aria-label="Buka menu" aria-expanded="false"><span class="mobile-menu-btn-area" aria-hidden="true"><span></span><span></span><span></span></span></button>
  <a class="mobile-admin-title" href="/admin">{{.Title}}</a>
  <div style="display:inline-flex; align-items:center;">
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right:8px; opacity:0.8; cursor:pointer;" onclick="this.nextElementSibling.click()"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2"></path><path d="M12 20v2"></path><path d="m4.93 4.93 1.41 1.41"></path><path d="m17.66 17.66 1.41 1.41"></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path></svg>
    <button class="mobile-theme-btn" type="button" onclick="toggleAdminMobileTheme()" aria-label="Ganti tema"><span class="sun">☀</span><span class="moon">🌙</span></button>
  </div>
</div>
<div class="mobile-drawer-overlay" id="adminDrawerOverlay" onclick="closeAdminDrawer()"></div>
<aside class="mobile-drawer" id="adminDrawer" aria-hidden="true">
  <div class="mobile-drawer-head">
    <strong>{{.Title}}</strong>
  </div>
  <div class="mobile-drawer-links" id="adminDrawerLinks"></div>
</aside>

<div class="nav">
<a class="brand" href="/admin">{{.Title}}</a>
<div class="menu">
<a href="/admin">Dashboard</a>
<a href="/admin/routers">Router</a>
<a href="/admin/packages">Paket</a>
<a href="/admin/transactions">Transaksi</a>
<a hidden class="super-only owner-menu-navbar-v1" href="/admin/owner">Owner</a>
<a href="/admin/billing">Billing</a>
<a class="super-only" href="/admin/health">Health</a>
<a class="super-only" href="/admin/settings">Settings</a>
<a class="super-only" href="/admin/payment-gateway">Payment Gateway</a>
<a class="super-only" href="/admin/blocked-phones">Nomor Diblokir</a>
<a class="super-only" href="/admin/users">Users</a>
<a href="/logout">Logout</a>
<div style="display:inline-flex; align-items:center;">
<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right:8px; opacity:0.8; cursor:pointer;" onclick="this.nextElementSibling.click()"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2"></path><path d="M12 20v2"></path><path d="m4.93 4.93 1.41 1.41"></path><path d="m17.66 17.66 1.41 1.41"></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path></svg>
<button class="admin-theme-toggle" type="button" id="adminThemeToggle" title="Ganti mode tampilan" aria-label="Ganti mode tampilan"><span class="sun">☀️</span><span class="moon">🌙</span></button>
</div>
</div>
</div>
{{.Body}}
</div>

<script>
// ADMIN_THEME_SYNC_V1
(function(){
	var keys = ["tuku-admin-theme", "vouchergo_theme", "theme", "tuku_theme"];
	var root = document.documentElement;

	function valid(v){
		return v === "light" || v === "dark";
	}

	function readSaved(){
		for(var i=0;i<keys.length;i++){
			try{
				var v = localStorage.getItem(keys[i]);
				if(valid(v)) return v;
			}catch(e){}
		}
		return "light";
	}

	function save(v){
		if(!valid(v)) return;
		for(var i=0;i<keys.length;i++){
			try{ localStorage.setItem(keys[i], v); }catch(e){}
		}
	}

	function apply(v){
		if(!valid(v)) return;

		root.setAttribute("data-admin-theme", v);
		root.setAttribute("data-theme", v);
		root.classList.toggle("theme-light", v === "light");
		root.classList.toggle("theme-dark", v === "dark");

		if(document.body){
			document.body.setAttribute("data-admin-theme", v);
			document.body.setAttribute("data-theme", v);
			document.body.classList.toggle("theme-light", v === "light");
			document.body.classList.toggle("theme-dark", v === "dark");
		}
	}

	function setTheme(v){
		if(!valid(v)) return;
		save(v);
		apply(v);
	}

	function current(){
		var v = root.getAttribute("data-admin-theme") || root.getAttribute("data-theme");
		return valid(v) ? v : readSaved();
	}

	window.setVoucherGoTheme = setTheme;
	window.toggleVoucherGoTheme = function(){
		setTheme(current() === "dark" ? "light" : "dark");
	};

	setTheme(readSaved());

	document.addEventListener("DOMContentLoaded", function(){
		setTheme(readSaved());

		var btn = document.getElementById("adminThemeToggle");
		if(btn){
			btn.addEventListener("click", function(ev){
				ev.preventDefault();
				ev.stopPropagation();
				if(ev.stopImmediatePropagation) ev.stopImmediatePropagation();
				window.toggleVoucherGoTheme();
			}, true);
		}

		document.querySelectorAll(".mobile-theme-btn").forEach(function(mbtn){
			mbtn.addEventListener("click", function(ev){
				ev.preventDefault();
				ev.stopPropagation();
				if(ev.stopImmediatePropagation) ev.stopImmediatePropagation();
				window.toggleVoucherGoTheme();
			}, true);
		});
	});
})();
</script>


<script>
// ADMIN_MOBILE_DRAWER_ACTIVE_V1
(function(){
  function getRole(){
    var m = document.cookie.match(/(?:^|;\s*)tuku_admin_role=([^;]+)/);
    return m ? decodeURIComponent(m[1]) : "";
  }

  function buildAdminDrawer(){
    var target = document.getElementById("adminDrawerLinks");
    var nav = document.querySelector(".nav");
    if(!target || !nav) return;

    target.innerHTML = "";

    nav.querySelectorAll("a").forEach(function(a){
      if(a.classList.contains("brand")) return;
      if(getRole() !== "superadmin" && a.classList.contains("super-only")) return;

      var clone = a.cloneNode(true);
      clone.addEventListener("click", closeAdminDrawer);
      target.appendChild(clone);
    });
  }

  window.openAdminDrawer = function(){
    buildAdminDrawer();
    document.body.classList.add("admin-drawer-open");
    var drawer = document.getElementById("adminDrawer");
    if(drawer) drawer.setAttribute("aria-hidden", "false");
    var btn = document.querySelector(".mobile-menu-btn");
    if(btn) {
      btn.setAttribute("aria-expanded", "true");
      btn.setAttribute("aria-label", "Tutup menu");
    }
  };

  window.closeAdminDrawer = function(){
    document.body.classList.remove("admin-drawer-open");
    var drawer = document.getElementById("adminDrawer");
    if(drawer) drawer.setAttribute("aria-hidden", "true");
    var btn = document.querySelector(".mobile-menu-btn");
    if(btn) {
      btn.setAttribute("aria-expanded", "false");
      btn.setAttribute("aria-label", "Buka menu");
    }
  };

  window.toggleAdminDrawer = function(){
    if(document.body.classList.contains("admin-drawer-open")){
      closeAdminDrawer();
      return;
    }
    openAdminDrawer();
  };

  document.addEventListener("DOMContentLoaded", buildAdminDrawer);
  window.addEventListener("resize", function(){
    if(window.innerWidth > 1199) closeAdminDrawer();
  });
  document.addEventListener("keydown", function(e){
    if(e.key === "Escape") closeAdminDrawer();
  });

  setTimeout(buildAdminDrawer, 150);
  setTimeout(buildAdminDrawer, 600);
})();
</script>


<script>
// ADMIN_HEALTH_MOBILE_THEME_V1
(function(){
  function currentTheme(){
    return document.documentElement.getAttribute("data-admin-theme") || "light";
  }

  function setTheme(v){
    if(v !== "dark") v = "light";
    document.documentElement.setAttribute("data-admin-theme", v);
    try{
      localStorage.setItem("admin-theme", v);
      localStorage.setItem("adminTheme", v);
      localStorage.setItem("tuku_admin_theme", v);
    }catch(e){}
  }

  window.toggleAdminMobileTheme = function(){
    setTheme(currentTheme() === "dark" ? "light" : "dark");
  };

  try{
    var saved = localStorage.getItem("admin-theme") || localStorage.getItem("adminTheme") || localStorage.getItem("tuku_admin_theme");
    if(saved === "dark" || saved === "light") setTheme(saved);
  }catch(e){}
})();
</script>


<style>
/* ADMIN_TABLE_FILTER_SEARCH_V2 */
.vg-table-filter{
	display:grid;
	grid-template-columns:minmax(220px,1.4fr) repeat(5,minmax(130px,.7fr)) auto auto;
	gap:10px;
	align-items:end;
	margin:0 0 14px;
	padding:12px;
	border:1px solid rgba(148,163,184,.22);
	border-radius:18px;
	background:rgba(255,255,255,.035);
}
.vg-table-filter .field{display:grid;gap:6px;margin:0}
.vg-table-filter label{font-size:12px;font-weight:800;opacity:.78}
.vg-table-filter input,
.vg-table-filter select{
	width:100%;
	min-height:40px;
	border-radius:13px;
	border:1px solid rgba(148,163,184,.24);
	background:rgba(15,23,42,.16);
	color:inherit;
	padding:0 12px;
	outline:none;
}
.vg-table-filter input:focus,
.vg-table-filter select:focus{
	border-color:rgba(37,99,235,.72);
	box-shadow:0 0 0 3px rgba(37,99,235,.16);
}
.vg-filter-meta{
	font-size:12px;
	font-weight:800;
	opacity:.72;
	padding:10px 2px 0;
}
.vg-filter-actions{
	display:flex;
	gap:8px;
	flex-wrap:wrap;
}
.vg-filter-actions .btn{
	min-height:40px;
	white-space:nowrap;
}
@media(max-width:1180px){
	.vg-table-filter{grid-template-columns:1fr 1fr 1fr}
}
@media(max-width:720px){
	.vg-table-filter{grid-template-columns:1fr}
	.vg-filter-actions .btn{width:100%}
}
</style>
<script>
/* ADMIN_TABLE_FILTER_SEARCH_V2 */
document.addEventListener('DOMContentLoaded', function(){
	var path = window.location.pathname;
	if(['/admin/routers','/admin/packages','/admin/transactions'].indexOf(path) === -1) return;

	var wrap = document.querySelector('.tablewrap');
	if(!wrap) return;

	var table = wrap.querySelector('table');
	if(!table || table.dataset.vgFilterReady === '1') return;
	table.dataset.vgFilterReady = '1';

	var rows = Array.prototype.slice.call(table.querySelectorAll('tbody tr'));
	if(!rows.length) return;

	var headers = Array.prototype.slice.call(table.querySelectorAll('thead th')).map(function(th){
		return (th.textContent || '').trim().toLowerCase();
	});

	function findCol(names){
		for(var i=0;i<names.length;i++){
			var name = names[i];
			for(var j=0;j<headers.length;j++){
				if(headers[j] === name || headers[j].indexOf(name) >= 0) return j;
			}
		}
		return -1;
	}

	var statusIdx = findCol(['status']);
	var provisionIdx = findCol(['provision']);
	var waIdx = findCol(['wa']);
	var orderIdx = findCol(['order']);
	var isTx = path === '/admin/transactions';

	function cellText(row, idx){
		if(idx < 0) return '';
		var cell = row.children[idx];
		return cell ? (cell.textContent || '').trim() : '';
	}

	function rowDate(row){
		var source = isTx ? cellText(row, orderIdx) : '';
		var m = source.match(/[0-9]{4}-[0-9]{2}-[0-9]{2}/);
		return m ? m[0] : '';
	}

	function uniqueValues(idx){
		var set = {};
		rows.forEach(function(row){
			var v = cellText(row, idx).replace(/\s+/g, ' ').trim();
			if(v) set[v] = true;
		});
		return Object.keys(set).sort(function(a,b){ return a.localeCompare(b); });
	}

	var toolbar = document.createElement('div');
	toolbar.className = 'vg-table-filter';
	toolbar.innerHTML =
		'<div class="field">' +
			'<label>Search</label>' +
			'<input type="search" data-vg-search placeholder="Cari data...">' +
		'</div>' +
		'<div class="field" data-vg-status-wrap>' +
			'<label>Status</label>' +
			'<select data-vg-status><option value="">Semua status</option></select>' +
		'</div>' +
		'<div class="field" data-vg-provision-wrap>' +
			'<label>Provision</label>' +
			'<select data-vg-provision><option value="">Semua provision</option></select>' +
		'</div>' +
		'<div class="field" data-vg-wa-wrap>' +
			'<label>WA</label>' +
			'<select data-vg-wa><option value="">Semua WA</option></select>' +
		'</div>' +
		'<div class="field" data-vg-date-from-wrap>' +
			'<label>Dari Tanggal</label>' +
			'<input type="date" data-vg-date-from>' +
		'</div>' +
		'<div class="field" data-vg-date-to-wrap>' +
			'<label>Sampai Tanggal</label>' +
			'<input type="date" data-vg-date-to>' +
		'</div>' +
		'<div class="field">' +
			'<label>&nbsp;</label>' +
			'<button class="btn" type="button" data-vg-reset>Reset</button>' +
		'</div>' +
		'<div class="field" data-vg-export-wrap>' +
			'<label>&nbsp;</label>' +
			'<div class="vg-filter-actions">' +
				'<button class="btn pri" type="button" data-vg-export>Export Filter</button>' +
				'<button class="btn" type="button" data-vg-print>Print Filter</button>' +
			'</div>' +
		'</div>' +
		'<div class="field" data-vg-bulk-wrap>' +
			'<label>Aksi Terpilih</label>' +
			'<div class="tx-filter-bulk-actions">' +
				'<button class="btn light" type="submit" form="txBulkProvisionForm" onclick="return confirm(\\'Push ulang transaksi terpilih?\\')">Push</button>' +
				'<button class="btn ok" type="submit" form="txBulkWAForm" onclick="return confirm(\\'Kirim WhatsApp untuk transaksi terpilih?\\')">WA</button>' +
				'<button class="btn danger user-delete-btn" type="submit" form="txBulkDeleteForm" onclick="return confirm(\\'Hapus transaksi terpilih? Data yang dihapus tidak bisa dikembalikan.\\')">🗑 Hapus</button>' +
			'</div>' +
		'</div>';

	wrap.parentNode.insertBefore(toolbar, wrap);

	var search = toolbar.querySelector('[data-vg-search]');
	var status = toolbar.querySelector('[data-vg-status]');
	var provision = toolbar.querySelector('[data-vg-provision]');
	var wa = toolbar.querySelector('[data-vg-wa]');
	var dateFrom = toolbar.querySelector('[data-vg-date-from]');
	var dateTo = toolbar.querySelector('[data-vg-date-to]');

	var meta = document.createElement('div');
	meta.className = 'vg-filter-meta';
	toolbar.parentNode.insertBefore(meta, wrap);

	function fillSelect(select, values){
		values.forEach(function(v){
			var opt = document.createElement('option');
			opt.value = v.toLowerCase();
			opt.textContent = v;
			select.appendChild(opt);
		});
	}

	if(statusIdx >= 0) fillSelect(status, uniqueValues(statusIdx));
	else toolbar.querySelector('[data-vg-status-wrap]').style.display = 'none';

	if(isTx && provisionIdx >= 0) fillSelect(provision, uniqueValues(provisionIdx));
	else toolbar.querySelector('[data-vg-provision-wrap]').style.display = 'none';

	if(isTx && waIdx >= 0) fillSelect(wa, uniqueValues(waIdx));
	else toolbar.querySelector('[data-vg-wa-wrap]').style.display = 'none';

	if(!isTx){
		toolbar.querySelector('[data-vg-date-from-wrap]').style.display = 'none';
		toolbar.querySelector('[data-vg-date-to-wrap]').style.display = 'none';
		toolbar.querySelector('[data-vg-export-wrap]').style.display = 'none';
		toolbar.querySelector('[data-vg-bulk-wrap]').style.display = 'none';
	}

	function visibleRows(){
		return rows.filter(function(row){ return row.style.display !== 'none'; });
	}


	function syncTxBulkForms(){
		if(!isTx) return;

		var ids = [];
		table.querySelectorAll('tbody input[type="checkbox"][name="id"]:checked').forEach(function(cb){
			ids.push(cb.value);
		});

		['txBulkProvisionForm','txBulkWAForm'].forEach(function(formID){
			var form = document.getElementById(formID);
			if(!form) return;

			form.querySelectorAll('input[name="id"]').forEach(function(el){
				el.remove();
			});

			ids.forEach(function(id){
				var input = document.createElement('input');
				input.type = 'hidden';
				input.name = 'id';
				input.value = id;
				form.appendChild(input);
			});
		});
	}

	table.addEventListener('change', function(e){
		if(e.target && e.target.matches('input[type="checkbox"][name="id"]')){
			syncTxBulkForms();
		}
	});

	document.addEventListener('change', function(e){
		if(e.target && e.target.id === 'txCheckAll'){
			setTimeout(syncTxBulkForms, 0);
		}
	});


	function applyFilter(){
		var q = (search.value || '').toLowerCase().trim();
		var s = status.value || '';
		var p = provision.value || '';
		var w = wa.value || '';
		var df = dateFrom.value || '';
		var dt = dateTo.value || '';
		var shown = 0;

		rows.forEach(function(row){
			var all = (row.textContent || '').toLowerCase();
			var st = cellText(row, statusIdx).toLowerCase();
			var pr = cellText(row, provisionIdx).toLowerCase();
			var wt = cellText(row, waIdx).toLowerCase();
			var rd = rowDate(row);

			var ok = (!q || all.indexOf(q) >= 0) &&
				(!s || st.indexOf(s) >= 0) &&
				(!p || pr.indexOf(p) >= 0) &&
				(!w || wt.indexOf(w) >= 0) &&
				(!df || (rd && rd >= df)) &&
				(!dt || (rd && rd <= dt));

			row.style.display = ok ? '' : 'none';
			if(ok) shown++;
		});

		meta.textContent = 'Tampil ' + shown + ' dari ' + rows.length + ' data';
	}

	function csvEscape(v){
		v = (v || '').replace(/\s+/g, ' ').trim();
		if(/[",\n]/.test(v)) return '"' + v.replace(/"/g, '""') + '"';
		return v;
	}

	function exportFilteredCsv(){
		var cols = headers.map(function(h){ return h.toUpperCase(); });
		var lines = [cols.map(csvEscape).join(',')];

		visibleRows().forEach(function(row){
			var vals = Array.prototype.slice.call(row.children).map(function(td){
				return csvEscape(td.textContent || '');
			});
			lines.push(vals.join(','));
		});

		var blob = new Blob([lines.join('\n')], {type:'text/csv;charset=utf-8;'});
		var a = document.createElement('a');
		var date = new Date().toISOString().slice(0,10);
		a.href = URL.createObjectURL(blob);
		a.download = 'vouchergo-transaksi-filter-' + date + '.csv';
		document.body.appendChild(a);
		a.click();
		setTimeout(function(){
			URL.revokeObjectURL(a.href);
			a.remove();
		}, 200);
	}

	function printFiltered(){
		var clone = table.cloneNode(true);
		Array.prototype.slice.call(clone.querySelectorAll('tbody tr')).forEach(function(row, idx){
			if(rows[idx] && rows[idx].style.display === 'none') row.remove();
		});

		var win = window.open('', '_blank');
		if(!win){
			alert('Popup print diblokir browser.');
			return;
		}

		win.document.write('<!doctype html><html><head><title>Print Filter Transaksi</title>');
		win.document.write('<style>body{font-family:Arial,sans-serif;padding:20px}table{width:100%;border-collapse:collapse}th,td{border:1px solid #ddd;padding:8px;text-align:left}th{background:#f3f4f6}.meta{margin:0 0 14px;color:#555}</style>');
		win.document.write('</head><body>');
		win.document.write('<h2>VoucherGo - Hasil Filter Transaksi</h2>');
		win.document.write('<div class="meta">' + meta.textContent + '</div>');
		win.document.write(clone.outerHTML);
		win.document.write('
<script>
/* ADMIN_CSRF_AUTO_INJECT_V1 */
(function(){
	function readCookie(name){
		var parts = document.cookie ? document.cookie.split(";") : [];
		for(var i=0;i<parts.length;i++){
			var p = parts[i].trim();
			if(p.indexOf(name + "=") === 0){
				return decodeURIComponent(p.substring(name.length + 1));
			}
		}
		return "";
	}

	function injectCSRF(){
		var token = readCookie("tuku_csrf");
		if(!token) return;

		document.querySelectorAll("form").forEach(function(form){
			var method = String(form.getAttribute("method") || form.method || "get").toLowerCase();
			if(method !== "post") return;

			var input = form.querySelector('input[name="_csrf"]');
			if(!input){
				input = document.createElement("input");
				input.type = "hidden";
				input.name = "_csrf";
				form.appendChild(input);
			}
			input.value = token;
		});
	}

	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", injectCSRF);
	}else{
		injectCSRF();
	}

	setTimeout(injectCSRF, 250);
	setTimeout(injectCSRF, 800);
})();
</script>


<script>
(function(){
  // WA_GUI_LIVE_POLL_V1
  if (location.pathname !== "/admin/whatsapp") return;

  function textOf(v){ return (v == null ? "" : String(v)); }

  function ensureLiveBox(){
    var box = document.getElementById("wa-live-status-box-v1");
    if (box) return box;

    var grid = document.querySelector(".wa-status-grid") || document.querySelector(".grid");
    if (!grid) return null;

    box = document.createElement("div");
    box.id = "wa-live-status-box-v1";
    box.className = "card";
    box.style.border = "1px solid rgba(56,189,248,.35)";
    box.innerHTML =
      '<h2 style="margin-bottom:8px">Live Status</h2>' +
      '<div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">' +
      '<b id="wa-live-badge-v1">CHECKING</b>' +
      '<span class="muted" id="wa-live-conn-v1">Mengecek...</span>' +
      '</div>' +
      '<p class="muted" id="wa-live-device-v1" style="margin-top:8px"></p>' +
      '<p class="muted" id="wa-live-error-v1" style="margin-top:6px"></p>';

    grid.insertBefore(box, grid.firstChild);
    return box;
  }

  function updateMatchingTexts(online, connText){
    document.querySelectorAll("span,b,strong,small,div,p").forEach(function(el){
      var t = (el.textContent || "").trim();
      if (t === "ONLINE" || t === "OFFLINE") {
        el.textContent = online ? "ONLINE" : "OFFLINE";
      }
      if (t === "Connected" || t === "QR Required" || t === "connecting" || t === "connected") {
        el.textContent = connText;
      }
    });
  }

  function setLive(j){
    ensureLiveBox();

    var online = !!j.connected || j.connection === "open" || j.state === "connected";
    var conn = online ? "Connected" : textOf(j.connection || j.state || "QR Required");
    var badge = document.getElementById("wa-live-badge-v1");
    var connEl = document.getElementById("wa-live-conn-v1");
    var devEl = document.getElementById("wa-live-device-v1");
    var errEl = document.getElementById("wa-live-error-v1");

    if (badge) {
      badge.textContent = online ? "ONLINE" : "OFFLINE";
      badge.style.padding = "6px 10px";
      badge.style.borderRadius = "999px";
      badge.style.background = online ? "rgba(34,197,94,.16)" : "rgba(239,68,68,.16)";
      badge.style.color = online ? "#22c55e" : "#ef4444";
    }

    if (connEl) connEl.textContent = conn;
    if (devEl) {
      var id = j.me && j.me.id ? j.me.id : "";
      var name = j.me && j.me.name ? " / " + j.me.name : "";
      devEl.textContent = id ? ("Device: " + id + name) : "";
    }
    if (errEl) errEl.textContent = j.error || j.last_error ? ("Error: " + (j.error || j.last_error)) : "";

    updateMatchingTexts(online, conn);

    var img = document.querySelector(".wa-qr-img");
    if (img && !online && (j.has_qr || j.qr_available)) {
      img.src = "/admin/whatsapp/qr.png?t=" + Date.now();
      img.style.display = "";
    }
    if (img && online) {
      img.style.display = "none";
    }
  }

  async function poll(){
    try {
      var r = await fetch("/admin/whatsapp/status.json?t=" + Date.now(), {
        credentials: "same-origin",
        cache: "no-store"
      });
      var j = await r.json();
      setLive(j);
    } catch(e) {
      setLive({connected:false, connection:"error", error:"Gagal polling status WA"});
    }
  }

  poll();
  setInterval(poll, 2000);
})();
</script>




</body></html>');
		win.document.close();
		win.focus();
		setTimeout(function(){ win.print(); }, 300);
	}

	[search,status,provision,wa,dateFrom,dateTo].forEach(function(el){
		if(!el) return;
		el.addEventListener('input', applyFilter);
		el.addEventListener('change', applyFilter);
	});

	toolbar.querySelector('[data-vg-reset]').addEventListener('click', function(){
		search.value = '';
		status.value = '';
		provision.value = '';
		wa.value = '';
		dateFrom.value = '';
		dateTo.value = '';
		applyFilter();
	});

	var exportBtn = toolbar.querySelector('[data-vg-export]');
	var printBtn = toolbar.querySelector('[data-vg-print]');
	if(exportBtn) exportBtn.addEventListener('click', exportFilteredCsv);
	if(printBtn) printBtn.addEventListener('click', printFiltered);

	applyFilter();
});
</script>


<style>
/* ADMIN_TX_STATIC_TOOLBAR_FINAL_V1 */
.tx-static-filter{
	display:grid;
	grid-template-columns:minmax(220px,1.3fr) minmax(145px,.75fr) minmax(145px,.75fr) minmax(110px,.55fr) minmax(210px,.85fr) minmax(390px,1.5fr);
	gap:10px;
	align-items:end;
	margin:12px 0 14px;
	padding:12px;
	border:1px solid rgba(148,163,184,.22);
	border-radius:18px;
	background:rgba(255,255,255,.035);
}
.tx-static-filter .field{
	display:grid;
	gap:6px;
	margin:0;
}
.tx-static-filter label{
	font-size:12px;
	font-weight:800;
	opacity:.78;
}
.tx-static-filter input{
	width:100%;
	min-height:40px;
	border-radius:13px;
	border:1px solid rgba(148,163,184,.24);
	background:rgba(15,23,42,.16);
	color:inherit;
	padding:0 12px;
	outline:none;
}
.tx-filter-actions{
	display:grid;
	grid-template-columns:1fr 1fr;
	gap:8px;
}
.tx-bulk-actions-static{
	grid-template-columns:1fr 1fr 1fr;
}
.tx-filter-actions .btn{
	width:100%!important;
	min-height:40px!important;
	border-radius:13px!important;
	font-size:12px!important;
	white-space:nowrap!important;
}
.tx-head-check{
	display:flex!important;
	gap:7px!important;
	align-items:center!important;
	font-size:10px!important;
	font-weight:900!important;
	letter-spacing:.045em!important;
}
.tx-head-check input[type="checkbox"],
.tx-order-wrap input[type="checkbox"]{
	appearance:auto!important;
	-webkit-appearance:checkbox!important;
	display:inline-block!important;
	width:13px!important;
	height:13px!important;
	min-width:13px!important;
	max-width:13px!important;
	min-height:13px!important;
	max-height:13px!important;
	margin:0!important;
	padding:0!important;
	accent-color:#2563eb;
	cursor:pointer;
	position:static!important;
	box-shadow:none!important;
}
.tx-order-wrap{
	display:grid!important;
	grid-template-columns:15px minmax(0,1fr)!important;
	gap:7px!important;
	align-items:start!important;
}
.tx-order-text b{
	display:block;
	font-size:11px!important;
	line-height:1.15!important;
	white-space:nowrap!important;
	overflow:hidden!important;
	text-overflow:ellipsis!important;
}
.tx-order-text small{
	display:block;
	font-size:10px!important;
	line-height:1.15!important;
	opacity:.72;
}
.transactions-compact-table{
	table-layout:fixed!important;
	width:100%!important;
	font-size:12px!important;
}
.transactions-compact-table th,
.transactions-compact-table td{
	padding:8px 7px!important;
	vertical-align:middle!important;
	line-height:1.18!important;
}
.transactions-compact-table th{
	font-size:10px!important;
	letter-spacing:.045em!important;
}
.transactions-compact-table td:nth-child(1), .transactions-compact-table th:nth-child(1){width:17%!important}
.transactions-compact-table td:nth-child(2), .transactions-compact-table th:nth-child(2){width:13%!important}
.transactions-compact-table td:nth-child(3), .transactions-compact-table th:nth-child(3){width:12%!important}
.transactions-compact-table td:nth-child(4), .transactions-compact-table th:nth-child(4){width:8%!important}
.transactions-compact-table td:nth-child(5), .transactions-compact-table th:nth-child(5){width:8%!important}
.transactions-compact-table td:nth-child(6), .transactions-compact-table th:nth-child(6){width:11%!important}
.transactions-compact-table td:nth-child(7), .transactions-compact-table th:nth-child(7){width:11%!important}
.transactions-compact-table td:nth-child(8), .transactions-compact-table th:nth-child(8){width:12%!important}
.transactions-compact-table td:nth-child(9), .transactions-compact-table th:nth-child(9){width:8%!important}
.transactions-compact-table td:last-child .actions,
.transactions-compact-table td:last-child .tx-row-trash-only{
	display:flex!important;
	align-items:center!important;
	justify-content:center!important;
}
.tx-trash-icon{
	appearance:none!important;
	-webkit-appearance:none!important;
	border:0!important;
	background:transparent!important;
	box-shadow:none!important;
	outline:none!important;
	width:28px!important;
	height:28px!important;
	min-width:28px!important;
	min-height:28px!important;
	padding:0!important;
	margin:0!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	cursor:pointer!important;
	font-size:15px!important;
	line-height:1!important;
	color:#ef4444!important;
	border-radius:8px!important;
}
.tx-trash-icon:hover{
	background:rgba(239,68,68,.12)!important;
}
@media(max-width:1200px){
	.tx-static-filter{
		grid-template-columns:1fr 1fr!important;
	}
	.tx-static-filter .field:last-child{
		grid-column:1 / -1!important;
	}
}
@media(max-width:760px){
	.tx-static-filter{
		grid-template-columns:1fr!important;
		padding:10px!important;
		border-radius:16px!important;
	}
	.tx-filter-actions,
	.tx-bulk-actions-static{
		grid-template-columns:1fr!important;
	}
	.transactions-compact-table{
		min-width:900px!important;
	}
}
</style>
<script>
/* ADMIN_TX_STATIC_TOOLBAR_FINAL_V1 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	var table = document.querySelector('table.transactions-compact-table');
	var toolbar = document.getElementById('txStaticFilter');
	if(!table || !toolbar) return;

	// Hapus filter global lama kalau masih kebentuk dari script global.
	document.querySelectorAll('.vg-table-filter').forEach(function(el){
		el.remove();
	});

	var rows = Array.prototype.slice.call(table.querySelectorAll('tbody tr'));
	if(!rows.length){
		rows = Array.prototype.slice.call(table.querySelectorAll('tr')).slice(1);
	}

	var firstTh = table.querySelector('tr:first-child th:first-child');
	if(firstTh){
		var all = document.createElement('input');
		all.type = 'checkbox';
		all.id = 'txCheckAll';

		var label = document.createElement('label');
		label.className = 'tx-head-check';
		label.appendChild(all);

		var span = document.createElement('span');
		span.textContent = 'ORDER';
		label.appendChild(span);

		firstTh.innerHTML = '';
		firstTh.appendChild(label);

		all.addEventListener('change', function(){
			rows.forEach(function(row){
				if(row.style.display === 'none') return;
				var cb = row.querySelector('input[type="checkbox"][name="id"]');
				if(cb) cb.checked = all.checked;
			});
			syncBulkForms();
		});
	}

	var search = document.getElementById('txSearch');
	var from = document.getElementById('txDateFrom');
	var to = document.getElementById('txDateTo');
	var meta = document.getElementById('txFilterMeta');

	function rowDate(row){
		var txt = row.children[0] ? row.children[0].textContent : '';
		var m = txt.match(/[0-9]{4}-[0-9]{2}-[0-9]{2}/);
		return m ? m[0] : '';
	}

	function applyFilter(){
		var q = (search.value || '').toLowerCase().trim();
		var df = from.value || '';
		var dt = to.value || '';
		var shown = 0;

		rows.forEach(function(row){
			var all = (row.textContent || '').toLowerCase();
			var rd = rowDate(row);
			var ok = (!q || all.indexOf(q) >= 0) &&
				(!df || (rd && rd >= df)) &&
				(!dt || (rd && rd <= dt));

			row.style.display = ok ? '' : 'none';
			if(ok) shown++;
		});

		if(meta) meta.textContent = 'Tampil ' + shown + ' dari ' + rows.length + ' data';
		syncBulkForms();
	}

	function syncBulkForms(){
		var ids = [];
		rows.forEach(function(row){
			if(row.style.display === 'none') return;
			var cb = row.querySelector('input[type="checkbox"][name="id"]:checked');
			if(cb) ids.push(cb.value);
		});

		['txBulkProvisionForm','txBulkWAForm'].forEach(function(formID){
			var form = document.getElementById(formID);
			if(!form) return;

			form.querySelectorAll('input[name="id"]').forEach(function(el){ el.remove(); });

			ids.forEach(function(id){
				var input = document.createElement('input');
				input.type = 'hidden';
				input.name = 'id';
				input.value = id;
				form.appendChild(input);
			});
		});
	}

	table.addEventListener('change', function(e){
		if(e.target && e.target.matches('input[type="checkbox"][name="id"]')) syncBulkForms();
	});

	[search, from, to].forEach(function(el){
		if(!el) return;
		el.addEventListener('input', applyFilter);
		el.addEventListener('change', applyFilter);
	});

	var reset = document.getElementById('txResetFilter');
	if(reset){
		reset.addEventListener('click', function(){
			search.value = '';
			from.value = '';
			to.value = '';
			applyFilter();
		});
	}

	function visibleRows(){
		return rows.filter(function(row){ return row.style.display !== 'none'; });
	}

	function csvEscape(v){
		v = (v || '').replace(/\s+/g, ' ').trim();
		if(/[",\n]/.test(v)) return '"' + v.replace(/"/g, '""') + '"';
		return v;
	}

	var exportBtn = document.getElementById('txExportFilter');
	if(exportBtn){
		exportBtn.addEventListener('click', function(){
			var headers = Array.prototype.slice.call(table.querySelectorAll('tr:first-child th')).map(function(th){
				return csvEscape(th.textContent || '');
			});
			var lines = [headers.join(',')];

			visibleRows().forEach(function(row){
				var vals = Array.prototype.slice.call(row.children).map(function(td){
					return csvEscape(td.textContent || '');
				});
				lines.push(vals.join(','));
			});

			var blob = new Blob([lines.join('\n')], {type:'text/csv;charset=utf-8;'});
			var a = document.createElement('a');
			a.href = URL.createObjectURL(blob);
			a.download = 'vouchergo-transaksi-filter-' + new Date().toISOString().slice(0,10) + '.csv';
			document.body.appendChild(a);
			a.click();
			setTimeout(function(){ URL.revokeObjectURL(a.href); a.remove(); }, 200);
		});
	}

	var printBtn = document.getElementById('txPrintFilter');
	if(printBtn){
		printBtn.addEventListener('click', function(){
			var clone = table.cloneNode(true);
			Array.prototype.slice.call(clone.querySelectorAll('tr')).slice(1).forEach(function(row, idx){
				if(rows[idx] && rows[idx].style.display === 'none') row.remove();
			});

			var win = window.open('', '_blank');
			if(!win){
				alert('Popup print diblokir browser.');
				return;
			}
			win.document.write('<!doctype html><html><head><title>Print Filter Transaksi</title>');
			win.document.write('<style>body{font-family:Arial,sans-serif;padding:20px}table{width:100%;border-collapse:collapse}th,td{border:1px solid #ddd;padding:8px;text-align:left}th{background:#f3f4f6}.meta{margin:0 0 14px;color:#555}</style>');
			win.document.write('</head><body><h2>VoucherGo - Hasil Filter Transaksi</h2>');
			win.document.write('<div class="meta">' + (meta ? meta.textContent : '') + '</div>');
			win.document.write(clone.outerHTML);
			win.document.write('</body></html>');
			win.document.close();
			win.focus();
			setTimeout(function(){ win.print(); }, 300);
		});
	}

	applyFilter();
});
</script>


<style>
/* ADMIN_TX_FILTER_FLAT_FINAL_V1 */
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
	display:grid!important;
	grid-template-columns:minmax(260px,1.35fr) 145px 145px 86px 122px 112px 108px 108px 132px!important;
	gap:8px!important;
	align-items:end!important;
	padding:10px!important;
	border-radius:18px!important;
	overflow:hidden!important;
}
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .field{
	min-width:0!important;
	gap:5px!important;
}
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat label{
	font-size:10px!important;
	line-height:1!important;
	opacity:.72!important;
}
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat input{
	width:100%!important;
	height:34px!important;
	min-height:34px!important;
	border-radius:11px!important;
	font-size:12px!important;
	padding:0 10px!important;
}
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field .btn{
	width:100%!important;
	min-width:0!important;
	height:34px!important;
	min-height:34px!important;
	padding:0 8px!important;
	border-radius:11px!important;
	font-size:11px!important;
	line-height:1!important;
	white-space:nowrap!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
}
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field .btn.danger{
	padding:0 9px!important;
}
@media(max-width:1500px){
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
		grid-template-columns:minmax(220px,1fr) 135px 135px repeat(3,105px)!important;
	}
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field:nth-last-child(-n+3){
		grid-column:auto!important;
	}
}
@media(max-width:1050px){
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
		grid-template-columns:1fr 1fr!important;
	}
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-search-field{
		grid-column:1 / -1!important;
	}
}
@media(max-width:760px){
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
		grid-template-columns:1fr!important;
	}
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field .btn{
		width:100%!important;
	}
}
</style>



<style>
/* ADMIN_TX_STATUS_CHIPS_V1 */
body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
	grid-template-columns:minmax(260px,1.4fr) repeat(2,112px) repeat(4,120px)!important;
}

body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-date-field-v1{
	display:none!important;
}

body:has(.transactions-compact-table) #txStatusChips{
	grid-column:1 / -1;
	display:flex;
	align-items:center;
	gap:8px;
	flex-wrap:wrap;
	margin-top:2px;
}

body:has(.transactions-compact-table) #txStatusChips .tx-status-chip{
	border:1px solid rgba(148,163,184,.22);
	background:rgba(15,23,42,.34);
	color:inherit;
	border-radius:999px;
	padding:8px 13px;
	font-size:12px;
	font-weight:900;
	line-height:1;
	cursor:pointer;
}

body:has(.transactions-compact-table) #txStatusChips .tx-status-chip.active{
	background:#2563eb;
	border-color:#2563eb;
	color:#fff;
}

body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field .btn{
	font-weight:900!important;
}

@media(max-width:760px){
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
		gap:10px!important;
		padding:10px!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-search-field{
		grid-column:1 / -1!important;
	}

	body:has(.transactions-compact-table) #txStatusChips{
		grid-column:1 / -1!important;
		display:grid!important;
		grid-template-columns:repeat(4,minmax(0,1fr))!important;
		gap:7px!important;
		margin:0!important;
	}

	body:has(.transactions-compact-table) #txStatusChips .tx-status-chip{
		width:100%;
		padding:9px 6px!important;
		font-size:11px!important;
		text-align:center!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field{
		min-width:0!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field .btn{
		height:38px!important;
		min-height:38px!important;
		border-radius:13px!important;
		font-size:12px!important;
		padding:0 8px!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat label{
		font-size:0!important;
		height:0!important;
		line-height:0!important;
		overflow:hidden!important;
		margin:0!important;
	}
}
</style>

<script>
/* ADMIN_TX_STATUS_CHIPS_V1 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	function boot(){
		var table = document.querySelector('table.transactions-compact-table');
		var toolbar = document.getElementById('txStaticFilter');
		if(!table || !toolbar) return false;

		if(toolbar.dataset.statusChipsV1 === '1') return true;
		toolbar.dataset.statusChipsV1 = '1';

		var rows = Array.prototype.slice.call(table.querySelectorAll('tbody tr'));
		if(!rows.length){
			rows = Array.prototype.slice.call(table.querySelectorAll('tr')).slice(1);
		}

		function fieldByInput(id){
			var el = document.getElementById(id);
			return el ? el.closest('.field') : null;
		}

		var search = document.getElementById('txSearch');
		var from = document.getElementById('txDateFrom');
		var to = document.getElementById('txDateTo');
		var meta = document.getElementById('txFilterMeta');

		var fromField = fieldByInput('txDateFrom');
		var toField = fieldByInput('txDateTo');
		if(fromField) fromField.classList.add('tx-date-field-v1');
		if(toField) toField.classList.add('tx-date-field-v1');

		var exportBtn = document.getElementById('txExportFilter');
		var printBtn = document.getElementById('txPrintFilter');
		if(exportBtn) exportBtn.textContent = 'Export';
		if(printBtn) printBtn.textContent = 'Print';

		document.querySelectorAll('button').forEach(function(btn){
			var t = (btn.textContent || '').replace(/\s+/g,' ').trim().toLowerCase();
			if(t === 'lunas') btn.textContent = 'Tandai Lunas';
			if(t === 'wa') btn.textContent = 'WA Ulang';
			if(t === 'push') btn.textContent = 'Push MikroTik';
			if(t === 'push ulang' || t === 'push mikrotik') btn.textContent = 'Push MikroTik';
		});

		var chips = document.createElement('div');
		chips.id = 'txStatusChips';
		chips.innerHTML =
			'<button type="button" class="tx-status-chip active" data-status="">Semua</button>' +
			'<button type="button" class="tx-status-chip" data-status="pending">Pending</button>' +
			'<button type="button" class="tx-status-chip" data-status="paid">Lunas</button>' +
			'<button type="button" class="tx-status-chip" data-status="failed">Gagal</button>';

		var searchField = search ? search.closest('.field') : null;
		if(searchField && searchField.parentNode){
			searchField.parentNode.insertBefore(chips, searchField.nextSibling);
		}else{
			toolbar.insertBefore(chips, toolbar.firstChild);
		}

		var selectedStatus = '';

		function headerIndex(names){
			var ths = Array.prototype.slice.call(table.querySelectorAll('tr:first-child th'));
			var labels = ths.map(function(th){ return (th.textContent || '').replace(/\s+/g,' ').trim().toLowerCase(); });
			for(var n=0;n<names.length;n++){
				var key = names[n].toLowerCase();
				for(var i=0;i<labels.length;i++){
					if(labels[i] === key || labels[i].indexOf(key) >= 0) return i;
				}
			}
			return -1;
		}

		var statusIdx = headerIndex(['status']);

		function rowStatus(row){
			var txt = '';
			if(statusIdx >= 0 && row.children[statusIdx]){
				txt = row.children[statusIdx].textContent || '';
			}else{
				txt = row.textContent || '';
			}
			txt = txt.toLowerCase();

			if(txt.indexOf('paid') >= 0 || txt.indexOf('lunas') >= 0 || txt.indexOf('success') >= 0) return 'paid';
			if(txt.indexOf('pending') >= 0 || txt.indexOf('menunggu') >= 0 || txt.indexOf('unpaid') >= 0) return 'pending';
			if(txt.indexOf('failed') >= 0 || txt.indexOf('gagal') >= 0 || txt.indexOf('expired') >= 0 || txt.indexOf('cancelled') >= 0) return 'failed';
			return '';
		}

		function rowDate(row){
			var txt = row.children[0] ? row.children[0].textContent : '';
			var m = txt.match(/[0-9]{4}-[0-9]{2}-[0-9]{2}/);
			return m ? m[0] : '';
		}

		function syncBulkForms(){
			var ids = [];
			rows.forEach(function(row){
				if(row.style.display === 'none') return;
				var cb = row.querySelector('input[type="checkbox"][name="id"]:checked');
				if(cb) ids.push(cb.value);
			});

			['txBulkProvisionForm','txBulkWAForm','txBulkDeleteForm','txBulkMarkPaidForm'].forEach(function(formID){
				var form = document.getElementById(formID);
				if(!form) return;
				form.querySelectorAll('input[name="id"]').forEach(function(el){ el.remove(); });
				ids.forEach(function(id){
					var input = document.createElement('input');
					input.type = 'hidden';
					input.name = 'id';
					input.value = id;
					form.appendChild(input);
				});
			});
		}

		function applyFilterV1(){
			var q = search ? (search.value || '').toLowerCase().trim() : '';
			var df = from ? (from.value || '') : '';
			var dt = to ? (to.value || '') : '';
			var shown = 0;

			rows.forEach(function(row){
				var all = (row.textContent || '').toLowerCase();
				var rd = rowDate(row);
				var st = rowStatus(row);

				var ok = (!q || all.indexOf(q) >= 0) &&
					(!selectedStatus || st === selectedStatus) &&
					(!df || (rd && rd >= df)) &&
					(!dt || (rd && rd <= dt));

				row.style.display = ok ? '' : 'none';
				if(ok) shown++;
			});

			if(meta) meta.textContent = 'Tampil ' + shown + ' dari ' + rows.length + ' data';
			syncBulkForms();
		}

		chips.addEventListener('click', function(e){
			var btn = e.target.closest('.tx-status-chip');
			if(!btn) return;

			selectedStatus = btn.dataset.status || '';
			chips.querySelectorAll('.tx-status-chip').forEach(function(b){ b.classList.remove('active'); });
			btn.classList.add('active');

			applyFilterV1();
		});

		if(search){
			search.addEventListener('input', applyFilterV1);
			search.addEventListener('change', applyFilterV1);
		}
		if(from){
			from.addEventListener('input', applyFilterV1);
			from.addEventListener('change', applyFilterV1);
		}
		if(to){
			to.addEventListener('input', applyFilterV1);
			to.addEventListener('change', applyFilterV1);
		}

		var reset = document.getElementById('txResetFilter');
		if(reset){
			reset.addEventListener('click', function(){
				selectedStatus = '';
				chips.querySelectorAll('.tx-status-chip').forEach(function(b){ b.classList.remove('active'); });
				var allChip = chips.querySelector('[data-status=""]');
				if(allChip) allChip.classList.add('active');
				applyFilterV1();
			});
		}

		applyFilterV1();
		return true;
	}

	if(!boot()){
		setTimeout(boot, 300);
		setTimeout(boot, 900);
	}
});
</script>



<style>
/* ADMIN_TX_STATUS_FIX_V2 */
.tx-row-status-chip-v2{
	display:inline-flex;
	align-items:center;
	justify-content:center;
	margin-top:5px;
	padding:4px 8px;
	border-radius:999px;
	font-size:10px;
	font-weight:900;
	line-height:1;
	border:1px solid rgba(148,163,184,.22);
	background:rgba(148,163,184,.12);
	color:#cbd5e1;
	white-space:nowrap;
}
.tx-row-status-chip-v2.pending{
	background:rgba(245,158,11,.14);
	border-color:rgba(245,158,11,.35);
	color:#facc15;
}
.tx-row-status-chip-v2.paid{
	background:rgba(34,197,94,.14);
	border-color:rgba(34,197,94,.35);
	color:#4ade80;
}
.tx-row-status-chip-v2.failed{
	background:rgba(239,68,68,.14);
	border-color:rgba(239,68,68,.35);
	color:#f87171;
}
@media(max-width:760px){
	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(6) b{
		display:block;
	}
	.tx-row-status-chip-v2{
		font-size:9px!important;
		padding:3px 7px!important;
		margin-top:4px!important;
	}
}
</style>

<script>
/* ADMIN_TX_STATUS_FIX_V2 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	function norm(s){
		return (s || '').replace(/\s+/g,' ').trim();
	}

	function bootStatusFix(){
		var table = document.querySelector('table.transactions-compact-table');
		var toolbar = document.getElementById('txStaticFilter');
		var chips = document.getElementById('txStatusChips');
		if(!table || !toolbar || !chips) return false;

		var rows = Array.prototype.slice.call(table.querySelectorAll('tbody tr'));
		if(!rows.length){
			rows = Array.prototype.slice.call(table.querySelectorAll('tr')).slice(1);
		}

		function headers(){
			return Array.prototype.slice.call(table.querySelectorAll('tr:first-child th')).map(function(th){
				return norm(th.textContent).toLowerCase();
			});
		}

		function headerIndex(names){
			var h = headers();
			for(var n=0;n<names.length;n++){
				var key = names[n].toLowerCase();
				for(var i=0;i<h.length;i++){
					if(h[i] === key || h[i].indexOf(key) >= 0) return i;
				}
			}
			return -1;
		}

		var statusIdx = headerIndex(['status']);
		var voucherIdx = headerIndex(['voucher']);
		if(voucherIdx < 0) voucherIdx = 5;

		function classify(txt){
			txt = norm(txt).toLowerCase();
			if(!txt) return '';

			if(/\b(pending|menunggu|unpaid)\b/.test(txt)) return 'pending';
			if(/\b(failed|gagal|expired|cancelled|canceled|batal)\b/.test(txt)) return 'failed';
			if(/\b(paid|lunas|success|sukses)\b/.test(txt) && !/\bunpaid\b/.test(txt)) return 'paid';

			return '';
		}

		function rowStatus(row){
			// 1. Prioritas dari kolom status asli.
			if(statusIdx >= 0 && row.children[statusIdx]){
				var direct = classify(row.children[statusIdx].textContent || '');
				if(direct) return direct;
			}

			// 2. Cari cell yang isinya status pendek.
			for(var i=0;i<row.children.length;i++){
				var t = norm(row.children[i].textContent || '').toLowerCase();
				if(t.length <= 40){
					var st = classify(t);
					if(st) return st;
				}
			}

			// 3. Fallback badge class/text.
			var badges = row.querySelectorAll('.badge, span[class*="pending"], span[class*="paid"], span[class*="failed"], span[class*="success"]');
			for(var b=0;b<badges.length;b++){
				var bt = norm((badges[b].className || '') + ' ' + (badges[b].textContent || '')).toLowerCase();
				var bs = classify(bt);
				if(bs) return bs;
			}

			return '';
		}

		function statusLabel(st){
			if(st === 'paid') return 'Lunas';
			if(st === 'pending') return 'Pending';
			if(st === 'failed') return 'Gagal';
			return '-';
		}

		function renderRowStatusBadges(){
			rows.forEach(function(row){
				if(row.dataset.statusBadgeV2 === '1') return;
				row.dataset.statusBadgeV2 = '1';

				var st = rowStatus(row);
				var target = row.children[voucherIdx];
				if(!target || !st) return;

				var chip = document.createElement('div');
				chip.className = 'tx-row-status-chip-v2 ' + st;
				chip.textContent = statusLabel(st);
				target.appendChild(chip);
			});
		}

		function activeStatus(){
			var active = chips.querySelector('.tx-status-chip.active');
			return active ? (active.dataset.status || '') : '';
		}

		function showRow(row, ok){
			if(ok){
				row.style.setProperty('display', 'table-row', 'important');
			}else{
				row.style.setProperty('display', 'none', 'important');
			}
		}

		function syncBulkForms(){
			var ids = [];
			rows.forEach(function(row){
				if(row.style.getPropertyValue('display') === 'none') return;
				var cb = row.querySelector('input[type="checkbox"][name="id"]:checked');
				if(cb) ids.push(cb.value);
			});

			['txBulkProvisionForm','txBulkWAForm','txBulkDeleteForm','txBulkMarkPaidForm'].forEach(function(formID){
				var form = document.getElementById(formID);
				if(!form) return;

				form.querySelectorAll('input[name="id"]').forEach(function(el){ el.remove(); });

				ids.forEach(function(id){
					var input = document.createElement('input');
					input.type = 'hidden';
					input.name = 'id';
					input.value = id;
					form.appendChild(input);
				});
			});
		}

		function applyStatusFilterV2(){
			var search = document.getElementById('txSearch');
			var meta = document.getElementById('txFilterMeta');
			var q = search ? (search.value || '').toLowerCase().trim() : '';
			var selected = activeStatus();
			var shown = 0;

			rows.forEach(function(row){
				var all = (row.textContent || '').toLowerCase();
				var st = rowStatus(row);
				var ok = (!q || all.indexOf(q) >= 0) && (!selected || st === selected);

				showRow(row, ok);
				if(ok) shown++;
			});

			if(meta) meta.textContent = 'Tampil ' + shown + ' dari ' + rows.length + ' data';
			syncBulkForms();
		}

		function fixLabels(){
			document.querySelectorAll('button').forEach(function(btn){
				var t = norm(btn.textContent).toLowerCase();
				if(t === 'lunas') btn.textContent = 'Tandai Lunas';
				if(t === 'wa') btn.textContent = 'WA Ulang';
				if(t === 'push') btn.textContent = 'Push MikroTik';
			});
		}

		renderRowStatusBadges();
		fixLabels();

		var freshChips = chips.cloneNode(true);
		chips.parentNode.replaceChild(freshChips, chips);
		chips = freshChips;

		chips.addEventListener('click', function(e){
			var btn = e.target.closest('.tx-status-chip');
			if(!btn) return;

			chips.querySelectorAll('.tx-status-chip').forEach(function(b){ b.classList.remove('active'); });
			btn.classList.add('active');

			applyStatusFilterV2();
		});

		var search = document.getElementById('txSearch');
		if(search){
			search.addEventListener('input', applyStatusFilterV2);
			search.addEventListener('change', applyStatusFilterV2);
		}

		var reset = document.getElementById('txResetFilter');
		if(reset){
			reset.addEventListener('click', function(){
				if(search) search.value = '';
				chips.querySelectorAll('.tx-status-chip').forEach(function(b){ b.classList.remove('active'); });
				var all = chips.querySelector('[data-status=""]');
				if(all) all.classList.add('active');
				applyStatusFilterV2();
			});
		}

		setTimeout(function(){
			fixLabels();
			renderRowStatusBadges();
			applyStatusFilterV2();
		}, 300);

		applyStatusFilterV2();
		return true;
	}

	if(!bootStatusFix()){
		setTimeout(bootStatusFix, 400);
		setTimeout(bootStatusFix, 1000);
	}
});
</script>



<style>
/* ADMIN_TX_LUNAS_CHIP_FIX_V3 */
@media(max-width:760px){
	body:has(.transactions-compact-table) #txStatusChips{
		grid-template-columns:repeat(4,minmax(0,1fr))!important;
	}
	body:has(.transactions-compact-table) #txStatusChips .tx-status-chip{
		display:inline-flex!important;
		align-items:center!important;
		justify-content:center!important;
	}
}
</style>

<script>
/* ADMIN_TX_LUNAS_CHIP_FIX_V3 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	function norm(s){
		return (s || '').replace(/\s+/g,' ').trim();
	}

	function fixTxToolbarV3(){
		var chips = document.getElementById('txStatusChips');
		if(!chips) return false;

		var current = norm(chips.textContent).toLowerCase();
		if(
			current.indexOf('semua') < 0 ||
			current.indexOf('pending') < 0 ||
			current.indexOf('lunas') < 0 ||
			current.indexOf('gagal') < 0 ||
			current.indexOf('tandai lunas') >= 0
		){
			chips.innerHTML =
				'<button type="button" class="tx-status-chip active" data-status="">Semua</button>' +
				'<button type="button" class="tx-status-chip" data-status="pending">Pending</button>' +
				'<button type="button" class="tx-status-chip" data-status="paid">Lunas</button>' +
				'<button type="button" class="tx-status-chip" data-status="failed">Gagal</button>';
		}

		// Rename tombol aksi, tapi jangan sentuh status chips.
		document.querySelectorAll('button').forEach(function(btn){
			if(btn.closest('#txStatusChips')) return;

			var t = norm(btn.textContent).toLowerCase();
			var form = (btn.getAttribute('form') || '').toLowerCase();

			if(t === 'lunas' || form.indexOf('markpaid') >= 0 || form.indexOf('mark-paid') >= 0){
				btn.textContent = 'Tandai Lunas';
			}
			if(t === 'wa'){
				btn.textContent = 'WA Ulang';
			}
			if(t === 'push'){
				btn.textContent = 'Push MikroTik';
			}
		});

		return true;
	}

	fixTxToolbarV3();
	setTimeout(fixTxToolbarV3, 300);
	setTimeout(fixTxToolbarV3, 900);
	setTimeout(fixTxToolbarV3, 1500);
});
</script>



<style>
/* ADMIN_TX_HIDE_HEAD_CHECKBOX_V5 */
/* Pakai tombol "Pilih semua" di atas tabel saja, jangan dobel dengan checkbox header TGL */
body:has(.transactions-compact-table) .transactions-compact-table tr:first-child th:first-child input[type="checkbox"],
body:has(.transactions-compact-table) .transactions-compact-table tr:first-child th:first-child .tx-head-check input,
body:has(.transactions-compact-table) .transactions-compact-table tr:first-child th:first-child #txCheckAll{
	display:none!important;
}

body:has(.transactions-compact-table) .transactions-compact-table tr:first-child th:first-child .tx-head-check{
	gap:0!important;
}
</style>

<script>
/* ADMIN_TX_HIDE_HEAD_CHECKBOX_V5 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	function hideHeaderCheckbox(){
		var table = document.querySelector('table.transactions-compact-table');
		if(!table) return false;

		var headCb = table.querySelector('tr:first-child th:first-child input[type="checkbox"]');
		if(headCb){
			headCb.checked = false;
			headCb.indeterminate = false;
			headCb.style.display = 'none';
		}

		return true;
	}

	hideHeaderCheckbox();
	setTimeout(hideHeaderCheckbox, 300);
	setTimeout(hideHeaderCheckbox, 900);
});
</script>





<style>
/* ADMIN_TX_SELECT_BUTTON_V10 */
body:has(.transactions-compact-table) .transactions-compact-table tr:first-child th:first-child input[type="checkbox"],
body:has(.transactions-compact-table) .transactions-compact-table tr:first-child th:first-child #txCheckAll{
	display:none!important;
}

.tx-meta-select-row-v10{
	display:flex;
	align-items:center;
	justify-content:space-between;
	gap:10px;
	margin:12px 0 8px;
}

.tx-meta-select-row-v10 #txFilterMeta{
	margin:0!important;
}

#txSelectVisibleBtn{
	border:1px solid rgba(148,163,184,.24);
	background:rgba(15,23,42,.35);
	color:inherit;
	border-radius:999px;
	padding:9px 13px;
	font-size:12px;
	font-weight:900;
	line-height:1;
	cursor:pointer;
	white-space:nowrap;
}

#txSelectVisibleBtn.active{
	background:#2563eb;
	border-color:#2563eb;
	color:#fff;
}

@media(max-width:760px){
	.tx-meta-select-row-v10{
		margin:10px 0 7px!important;
	}
	#txSelectVisibleBtn{
		padding:9px 11px!important;
		font-size:11px!important;
	}
}
</style>

<script>
/* ADMIN_TX_SELECT_BUTTON_V10 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	function bootSelectButtonV10(){
		var table = document.querySelector('table.transactions-compact-table');
		var meta = document.getElementById('txFilterMeta');
		if(!table || !meta) return false;

		// hapus sisa wrapper lama
		document.querySelectorAll('.tx-meta-select-row-v9, .tx-meta-select-row-v4, .tx-select-visible-v4').forEach(function(el){
			if(el.id === 'txFilterMeta') return;
			el.remove();
		});

		var rows = Array.prototype.slice.call(table.querySelectorAll('tbody tr'));
		if(!rows.length){
			rows = Array.prototype.slice.call(table.querySelectorAll('tr')).slice(1);
		}

		var btn = document.getElementById('txSelectVisibleBtn');
		if(!btn){
			var row = document.createElement('div');
			row.className = 'tx-meta-select-row-v10';

			meta.parentNode.insertBefore(row, meta);
			row.appendChild(meta);

			btn = document.createElement('button');
			btn.type = 'button';
			btn.id = 'txSelectVisibleBtn';
			btn.textContent = 'Pilih semua';
			btn.dataset.checkedAll = '0';
			row.appendChild(btn);
		}

		function visibleRows(){
			return rows.filter(function(tr){
				var d = tr.style.getPropertyValue('display');
				return d !== 'none';
			});
		}

		function syncBulkForms(){
			var ids = [];

			rows.forEach(function(tr){
				if(tr.style.getPropertyValue('display') === 'none') return;

				var cb = tr.querySelector('input[type="checkbox"][name="id"]:checked');
				if(cb) ids.push(cb.value);
			});

			['txBulkProvisionForm','txBulkWAForm','txBulkDeleteForm','txBulkMarkPaidForm'].forEach(function(formID){
				var form = document.getElementById(formID);
				if(!form) return;

				form.querySelectorAll('input[name="id"]').forEach(function(el){ el.remove(); });

				ids.forEach(function(id){
					var input = document.createElement('input');
					input.type = 'hidden';
					input.name = 'id';
					input.value = id;
					form.appendChild(input);
				});
			});
		}

		function setBtnState(active){
			btn.dataset.checkedAll = active ? '1' : '0';
			btn.classList.toggle('active', !!active);
			btn.textContent = active ? 'Batal pilih' : 'Pilih semua';
		}

		function resetBtnState(){
			setBtnState(false);
		}

		function toggleVisible(){
			var shouldCheck = btn.dataset.checkedAll !== '1';

			visibleRows().forEach(function(tr){
				var cb = tr.querySelector('input[type="checkbox"][name="id"]');
				if(cb) cb.checked = shouldCheck;
			});

			syncBulkForms();
			setBtnState(shouldCheck);
		}

		if(btn.dataset.boundV10 !== '1'){
			btn.dataset.boundV10 = '1';
			btn.addEventListener('click', function(e){
				e.preventDefault();
				e.stopPropagation();
				toggleVisible();
			}, true);
		}

		table.addEventListener('change', function(e){
			if(e.target && e.target.matches('input[type="checkbox"][name="id"]')){
				syncBulkForms();

				// Kalau user manual uncheck salah satu row, tombol balik ke Pilih semua.
				var allVisibleChecked = visibleRows().length > 0 && visibleRows().every(function(tr){
					var cb = tr.querySelector('input[type="checkbox"][name="id"]');
					return cb && cb.checked;
				});
				setBtnState(allVisibleChecked);
			}
		});

		// Saat filter berubah, reset tombol supaya klik berikutnya selalu check semua yang tampil.
		document.addEventListener('click', function(e){
			if(e.target && (e.target.closest('#txStatusChips') || e.target.id === 'txResetFilter')){
				setTimeout(resetBtnState, 80);
				setTimeout(syncBulkForms, 120);
			}
		});

		var search = document.getElementById('txSearch');
		if(search){
			search.addEventListener('input', function(){
				setTimeout(resetBtnState, 80);
				setTimeout(syncBulkForms, 120);
			});
			search.addEventListener('change', function(){
				setTimeout(resetBtnState, 80);
				setTimeout(syncBulkForms, 120);
			});
		}

		resetBtnState();
		syncBulkForms();

		return true;
	}

	if(!bootSelectButtonV10()){
		setTimeout(bootSelectButtonV10, 300);
		setTimeout(bootSelectButtonV10, 900);
	}
});
</script>


<style>
/* ADMIN_PACKAGE_DELETE_TEXT_V1 */
.package-delete-btn{
	min-width:86px!important;
	min-height:38px!important;
	padding:9px 14px!important;
	border-radius:12px!important;
	font-size:13px!important;
	font-weight:900!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	background:#dc2626!important;
	color:#fff!important;
	border:0!important;
}
.package-delete-btn:hover{
	filter:brightness(1.05);
}
</style>


<style>
/* ADMIN_MOBILE_TABLE_CARDS_V1 */
@media(max-width:760px){
	body{
		overflow-x:hidden!important;
	}

	.card{
		border-radius:18px!important;
	}

	.tablewrap{
		overflow:visible!important;
		border:0!important;
		background:transparent!important;
	}

	.vg-mobile-card-table{
		display:block!important;
		width:100%!important;
		min-width:0!important;
		border:0!important;
		background:transparent!important;
	}

	.vg-mobile-card-table tbody{
		display:grid!important;
		gap:10px!important;
		width:100%!important;
	}

	.vg-mobile-card-table tr{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
		gap:9px 10px!important;
		width:100%!important;
		border:1px solid rgba(148,163,184,.22)!important;
		border-radius:18px!important;
		padding:12px!important;
		background:rgba(15,23,42,.20)!important;
		box-shadow:0 10px 24px rgba(0,0,0,.10)!important;
	}

	.vg-mobile-card-table tr.vg-mobile-header-row{
		display:none!important;
	}

	.vg-mobile-card-table th{
		display:none!important;
	}

	.vg-mobile-card-table td{
		display:grid!important;
		grid-template-columns:1fr!important;
		gap:4px!important;
		width:auto!important;
		padding:0!important;
		border:0!important;
		min-width:0!important;
		font-size:13px!important;
		line-height:1.25!important;
		word-break:break-word!important;
	}

	.vg-mobile-card-table td::before{
		content:attr(data-label);
		display:block!important;
		font-size:10px!important;
		font-weight:900!important;
		letter-spacing:.06em!important;
		text-transform:uppercase!important;
		opacity:.55!important;
		margin-bottom:1px!important;
	}

	.vg-mobile-card-table td:first-child{
		grid-column:1 / -1!important;
		padding:8px 9px!important;
		border-radius:14px!important;
		background:rgba(255,255,255,.035)!important;
		border:1px solid rgba(148,163,184,.14)!important;
	}

	.vg-mobile-card-table td:last-child{
		grid-column:1 / -1!important;
	}

	.vg-mobile-card-table td:last-child .actions,
	.vg-mobile-card-table td:last-child div.actions{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
		gap:8px!important;
		width:100%!important;
	}

	.vg-mobile-card-table td:last-child form{
		margin:0!important;
	}

	.vg-mobile-card-table td:last-child .btn,
	.vg-mobile-card-table td:last-child button{
		width:100%!important;
		min-height:38px!important;
		border-radius:13px!important;
		font-size:12px!important;
		padding:8px 10px!important;
	}

	.vg-mobile-card-table .badge{
		width:max-content!important;
		font-size:10px!important;
		padding:5px 8px!important;
		border-radius:999px!important;
	}

	.vg-mobile-card-table small,
	.vg-mobile-card-table .muted{
		font-size:11px!important;
		opacity:.72!important;
	}

	/* transaksi tetap lebih khusus */
	.transactions-compact-table.vg-mobile-card-table td:first-child{
		grid-column:1 / -1!important;
	}

	.transactions-compact-table.vg-mobile-card-table td:nth-child(9){
		grid-column:1 / -1!important;
	}

	.transactions-compact-table.vg-mobile-card-table td:nth-child(9) .actions,
	.transactions-compact-table.vg-mobile-card-table td:nth-child(9) .tx-row-trash-only{
		display:flex!important;
		justify-content:flex-start!important;
		width:auto!important;
	}

	.transactions-compact-table.vg-mobile-card-table .tx-trash-icon{
		width:38px!important;
		height:38px!important;
		min-width:38px!important;
		min-height:38px!important;
		border:1px solid rgba(239,68,68,.28)!important;
		background:rgba(239,68,68,.08)!important;
	}

	.transactions-compact-table.vg-mobile-card-table .tx-order-wrap{
		grid-template-columns:18px minmax(0,1fr)!important;
		gap:8px!important;
	}

	.transactions-compact-table.vg-mobile-card-table .tx-order-text b{
		font-size:12px!important;
		white-space:normal!important;
		overflow:visible!important;
		text-overflow:clip!important;
		word-break:break-all!important;
	}

	/* paket/router tombol hapus teks tetap enak */
	.package-delete-btn{
		width:100%!important;
		min-height:38px!important;
	}

	/* filter transaksi mobile */
	#txStaticFilter.tx-static-filter-flat{
		grid-template-columns:1fr!important;
	}

	#txStaticFilter.tx-static-filter-flat .tx-action-field .btn{
		width:100%!important;
	}
}
</style>
<script>
/* ADMIN_MOBILE_TABLE_CARDS_V1 */
document.addEventListener('DOMContentLoaded', function(){
	var path = window.location.pathname;
	var allowed = ['/admin/routers','/admin/packages'];
	if(allowed.indexOf(path) === -1) return;

	document.querySelectorAll('.tablewrap table').forEach(function(table){
		if(table.dataset.mobileCardsReady === '1') return;
		table.dataset.mobileCardsReady = '1';
		table.classList.add('vg-mobile-card-table');

		var rows = Array.prototype.slice.call(table.querySelectorAll('tr'));
		if(rows.length < 2) return;

		var headerRow = rows[0];
		headerRow.classList.add('vg-mobile-header-row');

		var headers = Array.prototype.slice.call(headerRow.querySelectorAll('th')).map(function(th){
			return (th.textContent || '').trim();
		});

		rows.slice(1).forEach(function(row){
			Array.prototype.slice.call(row.children).forEach(function(td, idx){
				var label = headers[idx] || '';
				if(label) td.setAttribute('data-label', label);
			});
		});
	});
});
</script>


<style>
/* ADMIN_TX_MOBILE_TABLE_COMPACT_V1 */
@media(max-width:760px){
	body:has(.transactions-compact-table) .transactions-page{
		padding:12px!important;
		border-radius:18px!important;
	}

	body:has(.transactions-compact-table) .tablewrap{
		overflow-x:auto!important;
		overflow-y:visible!important;
		-webkit-overflow-scrolling:touch!important;
		border-radius:16px!important;
		border:1px solid rgba(148,163,184,.18)!important;
		background:rgba(15,23,42,.16)!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table{
		display:table!important;
		min-width:880px!important;
		width:880px!important;
		table-layout:fixed!important;
		font-size:11px!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table tr{
		display:table-row!important;
		border:0!important;
		padding:0!important;
		background:transparent!important;
		box-shadow:none!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table th,
	body:has(.transactions-compact-table) .transactions-compact-table td{
		display:table-cell!important;
		padding:8px 6px!important;
		border-bottom:1px solid rgba(148,163,184,.13)!important;
		vertical-align:middle!important;
		font-size:11px!important;
		line-height:1.16!important;
		width:auto!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table th{
		font-size:9px!important;
		letter-spacing:.045em!important;
		opacity:.8!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table td::before{
		content:none!important;
		display:none!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(1),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(1){width:155px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(2),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(2){width:115px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(3),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(3){width:125px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(4),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(4){width:75px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(5),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(5){width:70px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(6),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(6){width:95px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(7),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(7){width:95px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(8),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(8){width:105px!important}

	body:has(.transactions-compact-table) .transactions-compact-table td:nth-child(9),
	body:has(.transactions-compact-table) .transactions-compact-table th:nth-child(9){width:45px!important}

	body:has(.transactions-compact-table) .tx-order-wrap{
		display:grid!important;
		grid-template-columns:14px minmax(0,1fr)!important;
		gap:6px!important;
		align-items:start!important;
		padding:0!important;
		background:transparent!important;
		border:0!important;
	}

	body:has(.transactions-compact-table) .tx-order-text b{
		font-size:10px!important;
		white-space:nowrap!important;
		overflow:hidden!important;
		text-overflow:ellipsis!important;
	}

	body:has(.transactions-compact-table) .tx-order-text small,
	body:has(.transactions-compact-table) .transactions-compact-table small,
	body:has(.transactions-compact-table) .transactions-compact-table .muted{
		font-size:9px!important;
		opacity:.68!important;
	}

	body:has(.transactions-compact-table) .transactions-compact-table .badge{
		font-size:9px!important;
		padding:4px 7px!important;
	}

	body:has(.transactions-compact-table) .tx-trash-icon{
		width:28px!important;
		height:28px!important;
		min-width:28px!important;
		min-height:28px!important;
		font-size:14px!important;
		border:0!important;
		background:transparent!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
		grid-template-columns:1fr!important;
		gap:8px!important;
		padding:10px!important;
		border-radius:16px!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat .tx-action-field .btn{
		width:100%!important;
	}
}
</style>


<style>
/* ADMIN_TX_MOBILE_FINALFIX_V5 */
@media(max-width:760px){
	body:has(.transactions-compact-table) #txStaticFilter.tx-static-filter-flat{
		display:grid!important;
		grid-template-columns:repeat(3,minmax(0,1fr))!important;
		gap:7px!important;
		padding:9px!important;
		border-radius:16px!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter .field{
		display:grid!important;
		min-width:0!important;
		gap:4px!important;
		margin:0!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter .tx-search-field{
		grid-column:1 / -1!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter label{
		display:block!important;
		font-size:9px!important;
		line-height:1!important;
		opacity:.68!important;
		white-space:nowrap!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter input{
		width:100%!important;
		height:30px!important;
		min-height:30px!important;
		border-radius:9px!important;
		font-size:10px!important;
		padding:0 8px!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter button{
		width:100%!important;
		height:30px!important;
		min-height:30px!important;
		padding:0 5px!important;
		border-radius:9px!important;
		font-size:10px!important;
		font-weight:900!important;
		line-height:1!important;
		white-space:nowrap!important;
		display:inline-flex!important;
		align-items:center!important;
		justify-content:center!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter .tx-mobile-short-btn{
		font-size:10px!important;
	}

	body:has(.transactions-compact-table) #txStaticFilter .tx-mobile-short-btn::before{
		content:none!important;
		display:none!important;
	}

	body:has(.transactions-compact-table) .tablewrap{
		overflow-x:hidden!important;
		overflow-y:visible!important;
		border-radius:16px!important;
		border:1px solid rgba(148,163,184,.16)!important;
		background:rgba(15,23,42,.16)!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table{
		display:table!important;
		width:100%!important;
		min-width:0!important;
		max-width:100%!important;
		table-layout:fixed!important;
		border-collapse:collapse!important;
		background:transparent!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table tbody{
		display:table-row-group!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table tr{
		display:table-row!important;
		border:0!important;
		padding:0!important;
		background:transparent!important;
		box-shadow:none!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table th,
	body:has(.transactions-compact-table) table.transactions-compact-table td{
		display:table-cell!important;
		padding:7px 5px!important;
		border-bottom:1px solid rgba(148,163,184,.13)!important;
		vertical-align:middle!important;
		font-size:10px!important;
		line-height:1.12!important;
		word-break:normal!important;
		width:auto!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table td::before{
		content:none!important;
		display:none!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table th{
		font-size:8px!important;
		letter-spacing:.04em!important;
		opacity:.78!important;
		white-space:nowrap!important;
	}

	/* Mobile: tampil cuma TGL | NAMA | VOUCHER | DETAIL */
	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(2),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(2),
	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(4),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(4),
	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(5),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(5),
	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(7),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(7),
	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(8),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(8){
		display:none!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(1),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(1){width:30%!important}

	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(3),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(3){width:27%!important}

	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(6),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(6){width:24%!important}

	body:has(.transactions-compact-table) table.transactions-compact-table th:nth-child(9),
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(9){
		width:19%!important;
		text-align:center!important;
	}

	body:has(.transactions-compact-table) .tx-order-wrap{
		display:grid!important;
		grid-template-columns:13px minmax(0,1fr)!important;
		gap:5px!important;
		align-items:center!important;
		padding:0!important;
		border:0!important;
		background:transparent!important;
	}

	body:has(.transactions-compact-table) .tx-order-wrap input[type="checkbox"]{
		width:12px!important;
		height:12px!important;
		min-width:12px!important;
		min-height:12px!important;
		margin:0!important;
	}

	body:has(.transactions-compact-table) .tx-order-text b{
		display:none!important;
	}

	body:has(.transactions-compact-table) .tx-order-text small{
		display:block!important;
		font-size:8px!important;
		line-height:1.1!important;
		opacity:.76!important;
		white-space:normal!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(3){
		font-weight:900!important;
		overflow:hidden!important;
		text-overflow:ellipsis!important;
		white-space:nowrap!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(3) br,
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(3) small,
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(3) .muted{
		display:none!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(6){
		font-weight:950!important;
		font-size:11px!important;
		white-space:nowrap!important;
		overflow:hidden!important;
		text-overflow:ellipsis!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(6) br,
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(6) small,
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(6) .muted{
		display:none!important;
	}

	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(9) .actions,
	body:has(.transactions-compact-table) table.transactions-compact-table td:nth-child(9) .tx-row-trash-only{
		display:flex!important;
		flex-direction:column!important;
		gap:4px!important;
		align-items:center!important;
		justify-content:center!important;
		width:100%!important;
	}

	body:has(.transactions-compact-table) .tx-mobile-detail-btn{
		display:inline-flex!important;
		align-items:center!important;
		justify-content:center!important;
		width:100%!important;
		height:26px!important;
		min-height:26px!important;
		padding:0 5px!important;
		border-radius:8px!important;
		border:1px solid rgba(59,130,246,.28)!important;
		background:rgba(37,99,235,.18)!important;
		color:#dbeafe!important;
		font-size:9px!important;
		font-weight:900!important;
		line-height:1!important;
		cursor:pointer!important;
	}

	body:has(.transactions-compact-table) .tx-trash-icon{
		width:24px!important;
		height:24px!important;
		min-width:24px!important;
		min-height:24px!important;
		margin:0 auto!important;
		font-size:12px!important;
	}

	body:has(.transactions-compact-table) tr.tx-mobile-detail-row{
		display:none!important;
	}

	body:has(.transactions-compact-table) tr.tx-mobile-detail-row.open{
		display:table-row!important;
	}

	body:has(.transactions-compact-table) tr.tx-mobile-detail-row td{
		display:table-cell!important;
		padding:0!important;
		border-bottom:1px solid rgba(148,163,184,.13)!important;
		width:100%!important;
	}

	body:has(.transactions-compact-table) .tx-mobile-detail-box{
		padding:9px!important;
		background:rgba(255,255,255,.035)!important;
	}

	body:has(.transactions-compact-table) .tx-mobile-detail-grid{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
		gap:7px 9px!important;
	}

	body:has(.transactions-compact-table) .tx-mobile-detail-item{
		display:grid!important;
		gap:2px!important;
	}

	body:has(.transactions-compact-table) .tx-mobile-detail-item b{
		font-size:8px!important;
		line-height:1!important;
		opacity:.72!important;
		text-transform:uppercase!important;
		letter-spacing:.05em!important;
	}

	body:has(.transactions-compact-table) .tx-mobile-detail-item span{
		font-size:10px!important;
		line-height:1.2!important;
		word-break:break-word!important;
	}
}
</style>
<script>
/* ADMIN_TX_MOBILE_FINALFIX_V5 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;
	if(window.innerWidth > 760) return;

	var filter = document.getElementById('txStaticFilter');
	if(filter){
		var labels = [
			['txResetFilter', 'Reset'],
			['txExportFilter', 'Export'],
			['txPrintFilter', 'Print']
		];

		labels.forEach(function(item){
			var btn = document.getElementById(item[0]);
			if(btn){
				btn.classList.add('tx-mobile-short-btn');
				btn.textContent = item[1];
			}
		});

		var pushBtn = filter.querySelector('button[form="txBulkProvisionForm"]');
		var waBtn = filter.querySelector('button[form="txBulkWAForm"]');
		var delBtn = filter.querySelector('button[form="txBulkDeleteForm"]');

		if(pushBtn) pushBtn.textContent = 'Push';
		if(waBtn) waBtn.textContent = 'WA';
		if(delBtn) delBtn.textContent = 'Hapus';
	}

	var table = document.querySelector('table.transactions-compact-table');
	if(!table) return;

	table.classList.remove('vg-mobile-card-table');
	table.querySelectorAll('[data-label]').forEach(function(el){
		el.removeAttribute('data-label');
	});

	// Hapus detail row lama dari script sebelumnya.
	table.querySelectorAll('tr.tx-mobile-detail-row').forEach(function(row){
		row.remove();
	});

	var ths = table.querySelectorAll('tr:first-child th');
	if(ths.length >= 9){
		ths[0].innerHTML = '<label class="tx-head-check"><input type="checkbox" id="txCheckAll"><span>TGL</span></label>';
		ths[2].textContent = 'NAMA';
		ths[5].textContent = 'VOUCHER';
		ths[8].textContent = 'DETAIL';
	}

	function cleanText(node){
		return ((node && node.textContent) || '').replace(/\s+/g, ' ').trim();
	}

	function esc(s){
		return String(s || '').replace(/[&<>"']/g, function(c){
			return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c];
		});
	}

	function mainText(cell){
		if(!cell) return '';
		var clone = cell.cloneNode(true);
		clone.querySelectorAll('small,.muted,br,input,button,form').forEach(function(el){ el.remove(); });
		return cleanText(clone);
	}

	var rows = Array.prototype.slice.call(table.querySelectorAll('tr')).slice(1);

	rows.forEach(function(row){
		if(row.classList.contains('tx-mobile-detail-row')) return;

		var cells = row.children;
		if(cells.length < 9) return;

		// Buang tombol detail lama yang mungkin dobel.
		cells[8].querySelectorAll('.tx-mobile-detail-btn').forEach(function(btn){ btn.remove(); });

		var dateText = cleanText(cells[0].querySelector('small')) || cleanText(cells[0]);
		var orderText = cleanText(cells[0].querySelector('b')) || cleanText(cells[0]);
		var routerText = cleanText(cells[1]);
		var buyerText = mainText(cells[2]) || cleanText(cells[2]);
		var totalText = cleanText(cells[3]);
		var statusText = cleanText(cells[4]);
		var voucherText = mainText(cells[5]) || cleanText(cells[5]);
		var provisionText = cleanText(cells[6]);
		var waText = cleanText(cells[7]);

		var btn = document.createElement('button');
		btn.type = 'button';
		btn.className = 'tx-mobile-detail-btn';
		btn.textContent = 'Detail';
		cells[8].insertBefore(btn, cells[8].firstChild);

		var detailRow = document.createElement('tr');
		detailRow.className = 'tx-mobile-detail-row';
		detailRow.innerHTML =
			'<td colspan="9">' +
				'<div class="tx-mobile-detail-box">' +
					'<div class="tx-mobile-detail-grid">' +
						'<div class="tx-mobile-detail-item"><b>Tanggal</b><span>' + esc(dateText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Order</b><span>' + esc(orderText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Nama</b><span>' + esc(buyerText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Voucher</b><span>' + esc(voucherText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Router/Paket</b><span>' + esc(routerText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Total</b><span>' + esc(totalText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Status</b><span>' + esc(statusText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>Provision</b><span>' + esc(provisionText) + '</span></div>' +
						'<div class="tx-mobile-detail-item"><b>WA</b><span>' + esc(waText) + '</span></div>' +
					'</div>' +
				'</div>' +
			'</td>';

		row.parentNode.insertBefore(detailRow, row.nextSibling);

		btn.addEventListener('click', function(e){
			e.preventDefault();
			e.stopPropagation();

			var open = detailRow.classList.contains('open');

			table.querySelectorAll('tr.tx-mobile-detail-row.open').forEach(function(openRow){
				openRow.classList.remove('open');
				var prev = openRow.previousElementSibling;
				if(prev){
					var b = prev.querySelector('.tx-mobile-detail-btn');
					if(b) b.textContent = 'Detail';
				}
			});

			if(!open){
				detailRow.classList.add('open');
				btn.textContent = 'Tutup';
			}else{
				detailRow.classList.remove('open');
				btn.textContent = 'Detail';
			}
		});
	});

	var all = document.getElementById('txCheckAll');
	if(all){
		all.addEventListener('change', function(){
			table.querySelectorAll('tr:not(.tx-mobile-detail-row) input[type="checkbox"][name="id"]').forEach(function(cb){
				var row = cb.closest('tr');
				if(row && row.style.display === 'none') return;
				cb.checked = all.checked;
			});
		});
	}
});
</script>


<style>
/* ADMIN_TX_MOBILE_DETAIL_BULK_GUARD_V6 */
body:has(.transactions-compact-table) tr.tx-mobile-detail-row:not(.open){
	display:none!important;
}
body:has(.transactions-compact-table) tr.tx-mobile-detail-row.open{
	display:table-row!important;
}
body:has(.transactions-compact-table) #txStaticFilter .tx-bulk-disabled{
	opacity:.48!important;
	filter:grayscale(.35)!important;
	cursor:not-allowed!important;
}
body:has(.transactions-compact-table) #txStaticFilter .tx-bulk-disabled:hover{
	filter:grayscale(.35)!important;
}
</style>
<script>
/* ADMIN_TX_MOBILE_DETAIL_BULK_GUARD_V6 */
document.addEventListener('DOMContentLoaded', function(){
	if(window.location.pathname !== '/admin/transactions') return;

	var tries = 0;

	function initGuard(){
		tries++;

		var table = document.querySelector('table.transactions-compact-table');
		var filter = document.getElementById('txStaticFilter');

		if((!table || !filter) && tries < 30){
			setTimeout(initGuard, 100);
			return;
		}
		if(!table || !filter) return;

		function mainRows(){
			return Array.prototype.slice.call(
				table.querySelectorAll('tr:not(.tx-mobile-detail-row)')
			).filter(function(row){
				return row.querySelector('input[type="checkbox"][name="id"]');
			});
		}

		function selectedIDs(){
			var ids = [];
			mainRows().forEach(function(row){
				if(row.style.display === 'none') return;
				var cb = row.querySelector('input[type="checkbox"][name="id"]');
				if(cb && cb.checked && cb.value) ids.push(cb.value);
			});
			return ids;
		}

		function syncBulkForms(){
			var ids = selectedIDs();

			['txBulkProvisionForm','txBulkWAForm'].forEach(function(formID){
				var form = document.getElementById(formID);
				if(!form) return;

				form.querySelectorAll('input[name="id"]').forEach(function(el){
					el.remove();
				});

				ids.forEach(function(id){
					var input = document.createElement('input');
					input.type = 'hidden';
					input.name = 'id';
					input.value = id;
					form.appendChild(input);
				});
			});

			var canUse = ids.length > 0;
			[
				filter.querySelector('button[form="txBulkProvisionForm"]'),
				filter.querySelector('button[form="txBulkWAForm"]'),
				filter.querySelector('button[form="txBulkDeleteForm"]')
			].forEach(function(btn){
				if(!btn) return;
				btn.classList.toggle('tx-bulk-disabled', !canUse);
				btn.setAttribute('aria-disabled', canUse ? 'false' : 'true');
			});
		}

		function closeAllDetails(){
			table.querySelectorAll('tr.tx-mobile-detail-row').forEach(function(detail){
				detail.classList.remove('open');
				detail.style.display = 'none';

				var prev = detail.previousElementSibling;
				if(prev){
					var btn = prev.querySelector('.tx-mobile-detail-btn');
					if(btn) btn.textContent = 'Detail';
				}
			});
		}

		function fixDetailButtons(){
			table.querySelectorAll('tr.tx-mobile-detail-row').forEach(function(detail){
				detail.classList.remove('open');
				detail.style.display = 'none';
			});

			table.querySelectorAll('.tx-mobile-detail-btn').forEach(function(oldBtn){
				if(oldBtn.dataset.v6Ready === '1') return;

				var btn = oldBtn.cloneNode(true);
				btn.dataset.v6Ready = '1';
				btn.textContent = 'Detail';
				oldBtn.replaceWith(btn);

				btn.addEventListener('click', function(e){
					e.preventDefault();
					e.stopPropagation();
					e.stopImmediatePropagation();

					var row = btn.closest('tr');
					if(!row) return;

					var detail = row.nextElementSibling;
					if(!(detail && detail.classList.contains('tx-mobile-detail-row'))) return;

					var isOpen = detail.classList.contains('open') || detail.style.display === 'table-row';

					closeAllDetails();

					if(!isOpen){
						detail.classList.add('open');
						detail.style.display = 'table-row';
						btn.textContent = 'Tutup';
					}else{
						detail.classList.remove('open');
						detail.style.display = 'none';
						btn.textContent = 'Detail';
					}
				}, true);
			});
		}

		function guardBulkButton(btn){
			if(!btn || btn.dataset.bulkGuardReady === '1') return;
			btn.dataset.bulkGuardReady = '1';

			btn.addEventListener('click', function(e){
				syncBulkForms();

				if(selectedIDs().length === 0){
					e.preventDefault();
					e.stopPropagation();
					e.stopImmediatePropagation();
					alert('Harap pilih transaksi dulu.');
					return false;
				}
			}, true);
		}

		guardBulkButton(filter.querySelector('button[form="txBulkProvisionForm"]'));
		guardBulkButton(filter.querySelector('button[form="txBulkWAForm"]'));
		guardBulkButton(filter.querySelector('button[form="txBulkDeleteForm"]'));

		table.addEventListener('change', function(e){
			if(e.target && e.target.matches('input[type="checkbox"][name="id"]')){
				syncBulkForms();
			}
			if(e.target && e.target.id === 'txCheckAll'){
				setTimeout(syncBulkForms, 0);
			}
		});

		var observer = new MutationObserver(function(){
			mainRows().forEach(function(row){
				var detail = row.nextElementSibling;
				if(!(detail && detail.classList.contains('tx-mobile-detail-row'))) return;

				if(row.style.display === 'none'){
					detail.classList.remove('open');
					detail.style.display = 'none';

					var btn = row.querySelector('.tx-mobile-detail-btn');
					if(btn) btn.textContent = 'Detail';
				}
			});
			syncBulkForms();
		});

		mainRows().forEach(function(row){
			observer.observe(row, {attributes:true, attributeFilter:['style','class']});
		});

		fixDetailButtons();
		syncBulkForms();
	}

	initGuard();
});
</script>


<style>
/* ADMIN_TX_BULK_BUTTON_NORMAL_COLOR_V1 */
body:has(.transactions-compact-table) #txStaticFilter .tx-bulk-disabled{
	opacity:1!important;
	filter:none!important;
	cursor:pointer!important;
}
body:has(.transactions-compact-table) #txStaticFilter .tx-bulk-disabled:hover{
	filter:brightness(1.05)!important;
}
</style>


<style>
/* OWNER_ONLY_MENU_HIDE_NON_OWNER_HOST_V1 */
html.vg-non-owner-host a[href="/admin/demo"],
html.vg-non-owner-host a[href="/admin/clients"]{
	display:none!important;
}
</style>
<script>
/* OWNER_ONLY_MENU_HIDE_NON_OWNER_HOST_V1 */
(function(){
	function applyOwnerMenuGuard(){
		var h = (location.hostname || '').toLowerCase();

		// Owner panel resmi saat ini.
		var isOwnerHost = (h === 'vouchergo.biz.id' || h === 'www.vouchergo.biz.id');

		// Demo/client tidak boleh lihat owner-only menu.
		var nonOwner = !isOwnerHost || h.indexOf('demo.') === 0 || h.indexOf('client-') === 0;

		if(!nonOwner) return;

		document.documentElement.classList.add('vg-non-owner-host');

		document.querySelectorAll('a[href="/admin/demo"],a[href="/admin/clients"]').forEach(function(a){
			a.remove();
		});
	}

	if(document.readyState === 'loading'){
		document.addEventListener('DOMContentLoaded', applyOwnerMenuGuard);
	}else{
		applyOwnerMenuGuard();
	}
	setTimeout(applyOwnerMenuGuard, 300);
})();
</script>


<style>
/* USER_DELETE_TEXT_BUTTON_V1 */
.user-delete-btn{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:82px!important;
	min-height:38px!important;
	padding:0 14px!important;
	border-radius:13px!important;
	background:#ef4444!important;
	border:1px solid rgba(239,68,68,.35)!important;
	color:#fff!important;
	font-weight:900!important;
	font-size:13px!important;
	line-height:1!important;
	text-decoration:none!important;
	box-shadow:0 8px 18px rgba(239,68,68,.16)!important;
}
.user-delete-btn:hover{
	filter:brightness(1.04)!important;
}
@media(max-width:760px){
	.user-delete-btn{
		width:100%!important;
		min-width:0!important;
	}
}
</style>


<style>
/* USER_DELETE_FORCE_BUTTON_V2 */
form[action="/admin/users/delete"]{
	display:inline-flex!important;
	margin:0!important;
	padding:0!important;
	vertical-align:middle!important;
}

form[action="/admin/users/delete"] button,
form[action="/admin/users/delete"] input[type="submit"],
a[href="/admin/users/delete"],
a[href^="/admin/users/delete?"],
.user-delete-btn{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:86px!important;
	min-height:38px!important;
	padding:0 14px!important;
	border-radius:13px!important;
	background:#ef4444!important;
	border:1px solid rgba(239,68,68,.42)!important;
	color:#fff!important;
	font-weight:900!important;
	font-size:13px!important;
	line-height:1!important;
	text-decoration:none!important;
	box-shadow:0 8px 18px rgba(239,68,68,.16)!important;
	cursor:pointer!important;
	appearance:none!important;
	-webkit-appearance:none!important;
}

form[action="/admin/users/delete"] button:hover,
form[action="/admin/users/delete"] input[type="submit"]:hover,
a[href="/admin/users/delete"]:hover,
a[href^="/admin/users/delete?"]:hover,
.user-delete-btn:hover{
	filter:brightness(1.04)!important;
	color:#fff!important;
}
</style>


<style>
/* USER_DELETE_REAL_BUTTON_V4 */
.user-delete-btn,
button.user-delete-btn{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:86px!important;
	min-height:38px!important;
	padding:0 14px!important;
	border-radius:13px!important;
	background:#ef4444!important;
	border:1px solid rgba(239,68,68,.42)!important;
	color:#fff!important;
	font-weight:900!important;
	font-size:13px!important;
	line-height:1!important;
	text-decoration:none!important;
	box-shadow:0 8px 18px rgba(239,68,68,.16)!important;
	cursor:pointer!important;
	appearance:none!important;
	-webkit-appearance:none!important;
}
.user-delete-btn:hover{
	filter:brightness(1.04)!important;
	color:#fff!important;
}
</style>


<style>
/* HEALTH_HIDE_WANESIA_LEGACY_V1 */
.health-hide-wanesia-legacy{
	display:none!important;
}
@media(min-width:761px){
	body:has(.HEALTH_HIDE_WANESIA_LEGACY_V1_ACTIVE) .grid{
		align-items:stretch;
	}
}
</style>
<script>
/* HEALTH_HIDE_WANESIA_LEGACY_V1 */
(function(){
	function run(){
		if(location.pathname !== "/admin/health") return;

		var marker = document.createElement("div");
		marker.className = "HEALTH_HIDE_WANESIA_LEGACY_V1_ACTIVE";
		marker.style.display = "none";
		document.body.appendChild(marker);

		document.querySelectorAll("label, th, div, span, b").forEach(function(el){
			var t = (el.textContent || "").replace(/\s+/g, " ").trim().toLowerCase();
			if(t !== "wanesia enabled") return;

			var box =
				el.closest(".field") ||
				el.closest(".form-field") ||
				el.closest(".input-group") ||
				el.parentElement;

			if(box){
				box.classList.add("health-hide-wanesia-legacy");
			}
		});
	}

	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", run);
	}else{
		run();
	}
	setTimeout(run, 300);
})();
</script>


<style>
/* PACKAGE_ROUTER_SELECT_V3 */
.package-router-filter-v3{
	display:grid!important;
	gap:8px!important;
	margin:12px 0 16px!important;
	padding:12px!important;
	border:1px solid rgba(148,163,184,.22)!important;
	border-radius:16px!important;
	background:rgba(255,255,255,.035)!important;
}
.package-router-filter-v3 label{
	font-weight:900!important;
	font-size:13px!important;
	opacity:.78!important;
}
.package-router-filter-v3 .actions{
	display:grid!important;
	grid-template-columns:minmax(220px,1fr) auto!important;
	gap:10px!important;
	align-items:center!important;
}
.package-router-filter-v3 select{
	width:100%!important;
	min-height:42px!important;
	border-radius:13px!important;
}
.package-empty-v3{
	padding:18px!important;
	border-radius:16px!important;
	border:1px dashed rgba(148,163,184,.35)!important;
	background:rgba(148,163,184,.08)!important;
	color:inherit!important;
	font-weight:800!important;
}
@media(max-width:760px){
	.package-router-filter-v3 .actions{
		grid-template-columns:1fr!important;
	}
	.package-router-filter-v3 .btn{
		width:100%!important;
	}
}
</style>


<style>
/* ADMIN_PREMIUM_NAV_DASHBOARD_MOBILE_V1 */
.nav{
	position:sticky!important;
	top:0!important;
	z-index:80!important;
	margin:0 0 20px!important;
	padding:14px 18px!important;
	border:1px solid rgba(148,163,184,.18)!important;
	border-radius:0 0 24px 24px!important;
	background:
		linear-gradient(135deg, rgba(15,23,42,.94), rgba(15,23,42,.82) 48%, rgba(2,132,199,.22))!important;
	backdrop-filter:blur(18px)!important;
	-webkit-backdrop-filter:blur(18px)!important;
	box-shadow:
		0 18px 45px rgba(2,8,23,.18),
		inset 0 1px 0 rgba(255,255,255,.08)!important;
}

html[data-admin-theme="light"] .nav,
html[data-theme="light"] .nav{
	background:
		linear-gradient(135deg, rgba(255,255,255,.86), rgba(248,250,252,.78) 48%, rgba(219,234,254,.78))!important;
	border-color:rgba(148,163,184,.26)!important;
	box-shadow:
		0 18px 45px rgba(15,23,42,.10),
		inset 0 1px 0 rgba(255,255,255,.88)!important;
}

.nav .brand{
	font-weight:950!important;
	letter-spacing:-.035em!important;
	font-size:24px!important;
	text-shadow:0 8px 22px rgba(15,23,42,.16)!important;
}

.nav .menu{
	gap:8px!important;
}

.nav .menu a,
.nav .menu button,
.admin-theme-toggle{
	border-radius:16px!important;
	border:1px solid rgba(148,163,184,.22)!important;
	background:rgba(15,23,42,.56)!important;
	box-shadow:
		0 10px 26px rgba(2,8,23,.12),
		inset 0 1px 0 rgba(255,255,255,.06)!important;
	transition:transform .15s ease, box-shadow .15s ease, background .15s ease!important;
}

html[data-admin-theme="light"] .nav .menu a,
html[data-theme="light"] .nav .menu a,
html[data-admin-theme="light"] .nav .menu button,
html[data-theme="light"] .nav .menu button,
html[data-admin-theme="light"] .admin-theme-toggle,
html[data-theme="light"] .admin-theme-toggle{
	background:rgba(255,255,255,.72)!important;
	border-color:rgba(148,163,184,.28)!important;
	box-shadow:
		0 10px 24px rgba(15,23,42,.08),
		inset 0 1px 0 rgba(255,255,255,.85)!important;
}

.nav .menu a:hover,
.nav .menu button:hover,
.admin-theme-toggle:hover{
	transform:translateY(-1px)!important;
	box-shadow:
		0 14px 30px rgba(37,99,235,.16),
		inset 0 1px 0 rgba(255,255,255,.10)!important;
}

.mobile-adminbar{
	position:sticky!important;
	top:0!important;
	z-index:90!important;
	margin:0 0 12px!important;
	padding:10px 12px!important;
	border-radius:0 0 22px 22px!important;
	border:1px solid rgba(148,163,184,.18)!important;
	background:
		linear-gradient(135deg, rgba(15,23,42,.94), rgba(15,23,42,.84) 50%, rgba(2,132,199,.24))!important;
	backdrop-filter:blur(18px)!important;
	-webkit-backdrop-filter:blur(18px)!important;
	box-shadow:
		0 16px 36px rgba(2,8,23,.20),
		inset 0 1px 0 rgba(255,255,255,.08)!important;
}

html[data-admin-theme="light"] .mobile-adminbar,
html[data-theme="light"] .mobile-adminbar{
	background:
		linear-gradient(135deg, rgba(255,255,255,.90), rgba(248,250,252,.82) 50%, rgba(219,234,254,.78))!important;
	border-color:rgba(148,163,184,.24)!important;
	box-shadow:
		0 16px 36px rgba(15,23,42,.10),
		inset 0 1px 0 rgba(255,255,255,.88)!important;
}

.mobile-menu-btn,
.mobile-theme-btn{
	border-radius:16px!important;
	box-shadow:0 10px 24px rgba(2,8,23,.14)!important;
}

.mobile-admin-title{
	font-weight:950!important;
	letter-spacing:-.03em!important;
}

/* khusus dashboard mobile: kartu statistik jadi 2 kolom */
@media(max-width:760px){
	html.is-admin-dashboard .grid{
		display:grid!important;
		grid-template-columns:repeat(2, minmax(0, 1fr))!important;
		gap:10px!important;
	}

	html.is-admin-dashboard .grid .card{
		min-height:74px!important;
		padding:13px 12px!important;
		border-radius:16px!important;
	}

	html.is-admin-dashboard .grid .card b{
		font-size:11px!important;
		line-height:1.1!important;
		letter-spacing:.01em!important;
	}

	html.is-admin-dashboard .grid .card .price{
		font-size:24px!important;
		line-height:1.05!important;
		margin-top:4px!important;
	}

	html.is-admin-dashboard .card:has(h3){
		padding:14px!important;
		border-radius:18px!important;
	}

	html.is-admin-dashboard .card:has(h3) .actions{
		display:grid!important;
		grid-template-columns:1fr!important;
		gap:8px!important;
	}

	html.is-admin-dashboard .card:has(h3) .actions .btn{
		width:100%!important;
		min-height:38px!important;
	}
}
</style>
<script>
/* ADMIN_PREMIUM_NAV_DASHBOARD_MOBILE_V1 */
(function(){
	function mark(){
		var p = window.location.pathname;
		document.documentElement.classList.add("admin-premium-ui");
		if(p === "/admin" || p === "/admin/"){
			document.documentElement.classList.add("is-admin-dashboard");
		}else{
			document.documentElement.classList.remove("is-admin-dashboard");
		}
	}
	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", mark);
	}else{
		mark();
	}
})();
</script>


<style>
/* CLIENT_CREATE_AUTODEFAULTS_WARNING_V1 */
.client-autodefault-warning-v1{
	margin:10px 0 14px!important;
	padding:12px 13px!important;
	border-radius:15px!important;
	border:1px solid rgba(245,158,11,.35)!important;
	background:rgba(245,158,11,.12)!important;
	color:#fbbf24!important;
	font-size:13px!important;
	line-height:1.45!important;
	font-weight:750!important;
}
.client-autodefault-warning-v1 b{
	color:#fde68a!important;
	font-weight:950!important;
}
html[data-admin-theme="light"] .client-autodefault-warning-v1,
html[data-theme="light"] .client-autodefault-warning-v1{
	background:#fffbeb!important;
	border-color:rgba(245,158,11,.38)!important;
	color:#92400e!important;
}
html[data-admin-theme="light"] .client-autodefault-warning-v1 b,
html[data-theme="light"] .client-autodefault-warning-v1 b{
	color:#78350f!important;
}
</style>


<style>
/* ADMIN_TX_TOOLBAR_DOM_FINAL_V4 */
@media(min-width:761px){
	.tx-toolbar-dom-final-v4{
		display:grid!important;
		grid-template-columns:minmax(280px,1.6fr) minmax(150px,.65fr) minmax(150px,.65fr) auto auto auto!important;
		gap:10px!important;
		align-items:end!important;
	}
	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4{
		grid-column:1/-1!important;
		display:flex!important;
		justify-content:flex-end!important;
		align-items:center!important;
		gap:10px!important;
		width:100%!important;
		margin-top:4px!important;
	}
	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4 button{
		min-height:42px!important;
		padding:0 18px!important;
		border-radius:14px!important;
		white-space:nowrap!important;
	}
}
@media(max-width:760px){
	.tx-bulk-row-dom-final-v4{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
		gap:8px!important;
		width:100%!important;
	}
	.tx-bulk-row-dom-final-v4 button{
		width:100%!important;
		min-height:38px!important;
	}
}
</style>
<script>
/* ADMIN_TX_TOOLBAR_DOM_FINAL_V4 */
(function(){
	function txText(el){
		return (el && el.textContent || "").replace(/\s+/g, " ").trim();
	}
	function txNorm(el){
		return txText(el).toLowerCase();
	}
	function selectedTxValues(){
		var vals = [];
		Array.prototype.slice.call(document.querySelectorAll('input[type="checkbox"]:checked')).forEach(function(cb){
			var v = String(cb.value || "").trim();
			if(!/^[0-9]+$/.test(v)) return;
			if(vals.indexOf(v) === -1) vals.push(v);
		});
		return vals;
	}
	function ensureSelected(){
		if(selectedTxValues().length > 0) return true;
		alert("Harap pilih transaksi terlebih dahulu.");
		return false;
	}
	function fillBulkForm(formId){
		var form = document.getElementById(formId);
		if(!form) return;
		form.querySelectorAll('input[data-tx-auto-v4="1"]').forEach(function(el){ el.remove(); });
		selectedTxValues().forEach(function(v){
			var input = document.createElement("input");
			input.type = "hidden";
			input.name = "tx_id";
			input.value = v;
			input.setAttribute("data-tx-auto-v4", "1");
			form.appendChild(input);
		});
	}
	function submitMarkPaid(){
		if(!ensureSelected()) return;
		if(!confirm("Tandai transaksi terpilih sebagai lunas manual?")) return;

		var form = document.createElement("form");
		form.method = "post";
		form.action = "/admin/transactions/mark-paid";

		selectedTxValues().forEach(function(v){
			var input = document.createElement("input");
			input.type = "hidden";
			input.name = "tx_id";
			input.value = v;
			form.appendChild(input);
		});

		document.body.appendChild(form);
		form.submit();
	}
	function firstButtonByText(list){
		var wanted = list.map(function(x){ return x.toLowerCase(); });
		return Array.prototype.slice.call(document.querySelectorAll("button")).find(function(btn){
			return wanted.indexOf(txNorm(btn)) >= 0;
		}) || null;
	}
	function cleanDuplicates(text){
		var low = text.toLowerCase();
		var btns = Array.prototype.slice.call(document.querySelectorAll("button")).filter(function(btn){
			return txNorm(btn) === low;
		});
		if(btns.length <= 1) return btns[0] || null;

		var keep = btns.find(function(btn){ return btn.classList.contains("tx-mark-paid-btn"); }) || btns[0];
		btns.forEach(function(btn){
			if(btn !== keep) btn.remove();
		});
		return keep;
	}
	function findToolbarFrom(btn){
		var n = btn;
		while(n && n !== document.body){
			if(
				n.querySelector &&
				n.querySelector('input[placeholder*="Cari"],input[data-vg-search],input[type="date"]') &&
				n.querySelector('button[form="txBulkProvisionForm"]') &&
				!n.querySelector('table')
			){
				return n;
			}
			n = n.parentElement;
		}
		return null;
	}
	function ensureMarkButton(anchor){
		var btn = document.querySelector(".tx-mark-paid-btn") || firstButtonByText(["Tandai Lunas"]);
		if(btn) return btn;

		btn = document.createElement("button");
		btn.type = "button";
		btn.className = "btn tx-mark-paid-btn";
		btn.textContent = "Tandai Lunas";
		btn.addEventListener("click", submitMarkPaid);

		if(anchor && anchor.parentElement){
			anchor.parentElement.insertBefore(btn, anchor);
		}else{
			document.body.appendChild(btn);
		}
		return btn;
	}
	function move(btn, row){
		if(!btn || !row) return;
		row.appendChild(btn);
	}
	function hookGuards(){
		document.querySelectorAll('button[form="txBulkProvisionForm"],button[form="txBulkSendWAForm"],button[form="txBulkDeleteForm"]').forEach(function(btn){
			if(btn.dataset.txGuardV4) return;
			btn.dataset.txGuardV4 = "1";
			btn.addEventListener("click", function(e){
				if(!ensureSelected()){
					e.preventDefault();
					e.stopPropagation();
					return false;
				}
				fillBulkForm(btn.getAttribute("form"));
			}, true);
		});
	}
	function arrange(){
		if(location.pathname !== "/admin/transactions") return;

		document.querySelectorAll('button[form="txBulkProvisionForm"]').forEach(function(btn){
			btn.textContent = "Push MikroTik";
		});
		document.querySelectorAll("button").forEach(function(btn){
			if(txNorm(btn) === "voucher mikrotik") btn.textContent = "Push MikroTik";
		});

		var pushBtn = document.querySelector('button[form="txBulkProvisionForm"]') || firstButtonByText(["Push MikroTik", "Voucher MikroTik"]);
		if(!pushBtn) return;

		var markBtn = ensureMarkButton(pushBtn);
		cleanDuplicates("Tandai Lunas");

		var toolbar = findToolbarFrom(pushBtn);
		if(!toolbar) return;

		toolbar.classList.add("tx-toolbar-dom-final-v4");

		var row = toolbar.querySelector(".tx-bulk-row-dom-final-v4");
		if(!row){
			row = document.createElement("div");
			row.className = "tx-bulk-row-dom-final-v4";
			toolbar.appendChild(row);
		}

		var waBtn = document.querySelector('button[form="txBulkSendWAForm"]') || firstButtonByText(["Kirim WA", "WA"]);
		var delBtn = document.querySelector('button[form="txBulkDeleteForm"]') || firstButtonByText(["Hapus Terpilih"]);

		move(markBtn, row);
		move(pushBtn, row);
		move(waBtn, row);
		move(delBtn, row);

		cleanDuplicates("Tandai Lunas");
		hookGuards();
	}
	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", arrange);
	}else{
		arrange();
	}
	setTimeout(arrange, 300);
	setTimeout(arrange, 900);
	setTimeout(arrange, 1500);
	setTimeout(arrange, 2500);
})();
</script>


<style>
/* ADMIN_TX_TOOLBAR_LEFT_ALIGN_V6 */
@media(min-width:761px){
	.tx-toolbar-dom-final-v4{
		padding:12px!important;
		gap:8px 10px!important;
		align-items:end!important;
		grid-template-columns:minmax(300px,1.7fr) minmax(145px,.6fr) minmax(145px,.6fr) 110px 150px 140px!important;
	}

	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4{
		grid-column:1/-1!important;
		display:flex!important;
		justify-content:flex-start!important;
		align-items:center!important;
		gap:8px!important;
		width:100%!important;
		margin-top:0!important;
		padding-top:0!important;
	}

	.tx-toolbar-dom-final-v4 label{
		margin-bottom:5px!important;
		font-size:12px!important;
		line-height:1.1!important;
	}

	.tx-toolbar-dom-final-v4 input,
	.tx-toolbar-dom-final-v4 select{
		min-height:38px!important;
		height:38px!important;
		border-radius:12px!important;
	}

	.tx-toolbar-dom-final-v4 button,
	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4 button{
		min-height:38px!important;
		height:38px!important;
		border-radius:12px!important;
		padding:0 14px!important;
		font-size:13px!important;
		line-height:1!important;
		white-space:nowrap!important;
	}

	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4 .tx-mark-paid-btn{
		min-width:132px!important;
	}

	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4 button[form="txBulkProvisionForm"]{
		min-width:136px!important;
	}

	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4 button[form="txBulkSendWAForm"]{
		min-width:120px!important;
	}

	.tx-toolbar-dom-final-v4 .tx-bulk-row-dom-final-v4 button[form="txBulkDeleteForm"]{
		min-width:138px!important;
	}
}
</style>


<style>
/* ADMIN_TX_PUSH_MIKROTIK_BLUEGRAY_V1 */
button[form="txBulkProvisionForm"],
.tx-bulk-row-dom-final-v4 button[form="txBulkProvisionForm"],
.tx-bulk-row-final-v3 button[form="txBulkProvisionForm"],
.tx-bulk-row-safe-v2 button[form="txBulkProvisionForm"]{
	background:linear-gradient(135deg,#64748b,#334155)!important;
	border-color:rgba(100,116,139,.55)!important;
	color:#fff!important;
	font-weight:950!important;
	box-shadow:0 8px 18px rgba(51,65,85,.18)!important;
}

button[form="txBulkProvisionForm"]:hover,
.tx-bulk-row-dom-final-v4 button[form="txBulkProvisionForm"]:hover{
	background:linear-gradient(135deg,#718096,#475569)!important;
	color:#fff!important;
	transform:translateY(-1px);
}
</style>


<style>
/* ADMIN_TX_MOBILE_TOOLBAR_POLISH_V7 */
@media(max-width:760px){
	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4{
		display:grid!important;
		grid-template-columns:repeat(3,minmax(0,1fr))!important;
		gap:8px!important;
		padding:10px!important;
		overflow:hidden!important;
		border-radius:18px!important;
	}

	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 .field,
	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 > div{
		min-width:0!important;
		width:100%!important;
	}

	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 .field:first-child{
		grid-column:1/-1!important;
	}

	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 label{
		font-size:10px!important;
		margin-bottom:4px!important;
		line-height:1!important;
	}

	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 input,
	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 select,
	html.tx-mobile-toolbar-polish-v7 .tx-toolbar-dom-final-v4 button{
		width:100%!important;
		min-width:0!important;
		height:34px!important;
		min-height:34px!important;
		padding:0 8px!important;
		border-radius:11px!important;
		font-size:11px!important;
		line-height:1!important;
		white-space:nowrap!important;
	}

	html.tx-mobile-toolbar-polish-v7 .tx-bulk-row-dom-final-v4{
		grid-column:1/-1!important;
		display:grid!important;
		grid-template-columns:repeat(2,minmax(0,1fr))!important;
		gap:8px!important;
		width:100%!important;
		min-width:0!important;
		margin:0!important;
		padding:0!important;
		justify-content:stretch!important;
		justify-items:stretch!important;
		overflow:hidden!important;
	}

	html.tx-mobile-toolbar-polish-v7 .tx-bulk-row-dom-final-v4 button{
		width:100%!important;
		min-width:0!important;
		height:36px!important;
		min-height:36px!important;
		padding:0 8px!important;
		font-size:11px!important;
		border-radius:12px!important;
		overflow:hidden!important;
		text-overflow:ellipsis!important;
	}

	html.tx-mobile-toolbar-polish-v7 table td:last-child a,
	html.tx-mobile-toolbar-polish-v7 table td:last-child button:not(.tx-mobile-detail-btn):not([data-tx-detail]),
	html.tx-mobile-toolbar-polish-v7 table td:last-child form{
		display:none!important;
	}
}
</style>
<script>
/* ADMIN_TX_MOBILE_TOOLBAR_POLISH_V7 */
(function(){
	function t(el){
		return (el && el.textContent || "").replace(/\s+/g," ").trim();
	}
	function mobile(){
		return window.matchMedia && window.matchMedia("(max-width: 760px)").matches;
	}
	function apply(){
		if(location.pathname !== "/admin/transactions") return;

		document.documentElement.classList.add("tx-mobile-toolbar-polish-v7");

		if(!mobile()) return;

		document.querySelectorAll("button").forEach(function(btn){
			var text = t(btn).toLowerCase();

			if(text === "tandai lunas"){
				btn.textContent = "Lunas";
			}
			if(text === "push mikrotik" || text === "push mikrotik" || text === "voucher mikrotik"){
				btn.textContent = "Push MT";
			}
			if(text === "kirim wa"){
				btn.textContent = "WA";
			}
			if(text === "hapus terpilih"){
				btn.textContent = "Hapus";
			}
			if(text === "export filter"){
				btn.textContent = "Export";
			}
			if(text === "print filter"){
				btn.textContent = "Print";
			}
		});

		// Sembunyikan Hapus per-row di tabel mobile, supaya kolom DETAIL bersih.
		document.querySelectorAll("table td:last-child a, table td:last-child button").forEach(function(el){
			if(t(el).toLowerCase() === "hapus"){
				el.style.display = "none";
			}
		});
	}

	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", apply);
	}else{
		apply();
	}
	setTimeout(apply,300);
	setTimeout(apply,900);
	setTimeout(apply,1600);
})();
</script>


<style>
/* CLIENT_BILLING_PHASE1_UI_V1 */
.client-billing-page-v1{
	display:grid!important;
	gap:16px!important;
}
.client-billing-hero-v1{
	display:flex!important;
	justify-content:space-between!important;
	align-items:flex-start!important;
	gap:16px!important;
	padding:18px!important;
	border-radius:22px!important;
	border:1px solid rgba(148,163,184,.22)!important;
	background:
		radial-gradient(circle at 0% 0%, rgba(59,130,246,.18), transparent 34%),
		radial-gradient(circle at 100% 0%, rgba(168,85,247,.14), transparent 30%),
		rgba(255,255,255,.04)!important;
	box-shadow:0 18px 44px rgba(15,23,42,.10)!important;
}
.client-billing-hero-v1 h2{
	margin:4px 0 6px!important;
	font-size:28px!important;
	font-weight:950!important;
	letter-spacing:-.04em!important;
}
.client-billing-hero-v1 p{
	margin:0!important;
	opacity:.76!important;
	font-weight:700!important;
	line-height:1.45!important;
}
.billing-kicker-v1{
	display:inline-flex!important;
	padding:6px 10px!important;
	border-radius:999px!important;
	background:rgba(59,130,246,.13)!important;
	color:#60a5fa!important;
	font-size:12px!important;
	font-weight:950!important;
}
.billing-status-pill-v1{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:120px!important;
	padding:9px 13px!important;
	border-radius:999px!important;
	font-weight:950!important;
	font-size:13px!important;
	white-space:nowrap!important;
}
.billing-status-pill-v1.ok{
	background:rgba(34,197,94,.14)!important;
	color:#22c55e!important;
	border:1px solid rgba(34,197,94,.28)!important;
}
.billing-status-pill-v1.warn{
	background:rgba(245,158,11,.14)!important;
	color:#f59e0b!important;
	border:1px solid rgba(245,158,11,.32)!important;
}
.billing-status-pill-v1.bad{
	background:rgba(239,68,68,.13)!important;
	color:#ef4444!important;
	border:1px solid rgba(239,68,68,.30)!important;
}
.billing-grid-v1{
	grid-template-columns:repeat(4,minmax(0,1fr))!important;
}
.billing-stat-v1 .price{
	font-size:24px!important;
	line-height:1.08!important;
	letter-spacing:-.04em!important;
	overflow-wrap:anywhere!important;
}
.billing-action-card-v1{
	display:grid!important;
	gap:12px!important;
}
.billing-action-card-v1 h3{
	margin:0!important;
	font-size:22px!important;
	font-weight:950!important;
	letter-spacing:-.035em!important;
}
.billing-action-card-v1 p{
	margin:0!important;
	opacity:.82!important;
	font-weight:700!important;
	line-height:1.45!important;
}
.billing-actions-v1{
	display:flex!important;
	flex-wrap:wrap!important;
	gap:10px!important;
}
.billing-note-v1{
	padding:12px 13px!important;
	border-radius:15px!important;
	border:1px solid rgba(59,130,246,.22)!important;
	background:rgba(59,130,246,.09)!important;
	font-size:13px!important;
	font-weight:750!important;
	line-height:1.45!important;
}
.billing-note-v1 b{
	font-weight:950!important;
}
@media(max-width:900px){
	.billing-grid-v1{
		grid-template-columns:repeat(2,minmax(0,1fr))!important;
	}
}
@media(max-width:760px){
	.client-billing-hero-v1{
		display:grid!important;
		padding:14px!important;
		border-radius:18px!important;
	}
	.client-billing-hero-v1 h2{
		font-size:22px!important;
	}
	.billing-grid-v1{
		grid-template-columns:repeat(2,minmax(0,1fr))!important;
		gap:10px!important;
	}
	.billing-stat-v1{
		padding:13px!important;
		border-radius:16px!important;
	}
	.billing-stat-v1 .price{
		font-size:19px!important;
	}
	.billing-actions-v1{
		display:grid!important;
		grid-template-columns:1fr!important;
	}
	.billing-actions-v1 .btn{
		width:100%!important;
	}
}
</style>


<style>
/* OWNER_CLIENT_BILLING_MANAGER_UI_V1 */
.owner-client-billing-page-v1{
	display:grid!important;
	gap:14px!important;
}
.owner-client-billing-head-v1{
	display:flex!important;
	align-items:flex-start!important;
	justify-content:space-between!important;
	gap:12px!important;
}
.owner-client-billing-head-v1 h2{
	margin:0!important;
	font-weight:950!important;
	letter-spacing:-.04em!important;
}
.owner-client-billing-msg-v1{
	padding:11px 13px!important;
	border-radius:15px!important;
	background:rgba(59,130,246,.10)!important;
	border:1px solid rgba(59,130,246,.22)!important;
	font-weight:850!important;
}
.owner-client-billing-table-v1 td small{
	display:block!important;
	max-width:300px!important;
	overflow:hidden!important;
	text-overflow:ellipsis!important;
	white-space:nowrap!important;
	color:#64748b!important;
	font-weight:700!important;
}
.billing-status-mini-v1{
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	min-width:78px!important;
	padding:7px 10px!important;
	border-radius:999px!important;
	font-size:12px!important;
	font-weight:950!important;
	text-transform:capitalize!important;
}
.billing-status-mini-v1.ok{
	background:rgba(34,197,94,.14)!important;
	border:1px solid rgba(34,197,94,.28)!important;
	color:#15803d!important;
}
.billing-status-mini-v1.warn{
	background:rgba(245,158,11,.14)!important;
	border:1px solid rgba(245,158,11,.32)!important;
	color:#b45309!important;
}
.billing-status-mini-v1.bad{
	background:rgba(239,68,68,.12)!important;
	border:1px solid rgba(239,68,68,.30)!important;
	color:#b91c1c!important;
}
.owner-client-billing-actions-v1{
	display:flex!important;
	flex-wrap:wrap!important;
	gap:8px!important;
}
.owner-client-billing-actions-v1 form{
	margin:0!important;
}
.owner-client-billing-actions-v1 .btn{
	min-height:36px!important;
	padding:8px 11px!important;
	border-radius:12px!important;
	font-size:12px!important;
}
@media(max-width:760px){
	.owner-client-billing-head-v1{
		display:grid!important;
	}
	.owner-client-billing-table-v1,
	.owner-client-billing-table-v1 thead,
	.owner-client-billing-table-v1 tbody,
	.owner-client-billing-table-v1 tr,
	.owner-client-billing-table-v1 th,
	.owner-client-billing-table-v1 td{
		display:block!important;
		width:100%!important;
	}
	.owner-client-billing-table-v1 thead{
		display:none!important;
	}
	.owner-client-billing-table-v1 tr{
		margin-bottom:12px!important;
		padding:12px!important;
		border-radius:18px!important;
		background:rgba(248,250,252,.72)!important;
		border:1px solid rgba(148,163,184,.18)!important;
	}
	.owner-client-billing-table-v1 td{
		padding:7px 0!important;
		border:0!important;
	}
	.owner-client-billing-actions-v1{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
	}
	.owner-client-billing-actions-v1 .btn,
	.owner-client-billing-actions-v1 form{
		width:100%!important;
	}
}
</style>
<script>
/* OWNER_CLIENT_BILLING_MANAGER_UI_V1 */
(function(){
	function apply(){
		var h = location.hostname;
		if(h !== "vouchergo.biz.id" && h !== "www.vouchergo.biz.id"){
			document.querySelectorAll('a[href="/admin/client-billing"], .owner-billing-menu-v1').forEach(function(el){
				el.style.display = "none";
			});
		}
	}
	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", apply);
	}else{
		apply();
	}
	setTimeout(apply, 300);
})();
</script>


<style>
/* CLIENT_BILLING_WARNING_BANNER_UI_V1 */
.client-billing-warning-banner-v1{
	display:grid!important;
	grid-template-columns:1fr auto!important;
	align-items:center!important;
	gap:12px!important;
	margin:0 0 14px!important;
	padding:13px 14px!important;
	border-radius:18px!important;
	border:1px solid rgba(245,158,11,.34)!important;
	background:
		radial-gradient(circle at 0% 0%, rgba(245,158,11,.20), transparent 34%),
		rgba(255,251,235,.92)!important;
	color:#78350f!important;
	box-shadow:0 14px 36px rgba(120,53,15,.10)!important;
}
.client-billing-warning-banner-v1.bad{
	border-color:rgba(239,68,68,.34)!important;
	background:
		radial-gradient(circle at 0% 0%, rgba(239,68,68,.18), transparent 34%),
		rgba(254,242,242,.94)!important;
	color:#7f1d1d!important;
}
.client-billing-warning-banner-v1 b{
	display:block!important;
	font-size:14px!important;
	font-weight:950!important;
	margin-bottom:2px!important;
}
.client-billing-warning-banner-v1 span{
	display:block!important;
	font-size:13px!important;
	font-weight:800!important;
	line-height:1.35!important;
}
.client-billing-warning-actions-v1{
	display:flex!important;
	align-items:center!important;
	gap:8px!important;
}
.client-billing-warning-actions-v1 a{
	min-height:36px!important;
	padding:8px 12px!important;
	border-radius:12px!important;
	font-size:12px!important;
	white-space:nowrap!important;
}
html[data-admin-theme="dark"] .client-billing-warning-banner-v1{
	background:
		radial-gradient(circle at 0% 0%, rgba(245,158,11,.18), transparent 34%),
		rgba(69,26,3,.72)!important;
	color:#fde68a!important;
	border-color:rgba(245,158,11,.28)!important;
}
html[data-admin-theme="dark"] .client-billing-warning-banner-v1.bad{
	background:
		radial-gradient(circle at 0% 0%, rgba(239,68,68,.18), transparent 34%),
		rgba(69,10,10,.72)!important;
	color:#fecaca!important;
	border-color:rgba(239,68,68,.28)!important;
}
@media(max-width:760px){
	.client-billing-warning-banner-v1{
		grid-template-columns:1fr!important;
		padding:12px!important;
		border-radius:16px!important;
	}
	.client-billing-warning-actions-v1{
		display:grid!important;
		grid-template-columns:1fr 1fr!important;
	}
	.client-billing-warning-actions-v1 a{
		width:100%!important;
	}
}
</style>
<script>
/* CLIENT_BILLING_WARNING_BANNER_UI_V1 */
(function(){
	function ready(fn){
		if(document.readyState === "loading"){
			document.addEventListener("DOMContentLoaded", fn);
		}else{
			fn();
		}
	}

	function esc(s){
		return String(s || "").replace(/[&<>"']/g, function(c){
			return {"&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;","'":"&#39;"}[c];
		});
	}

	ready(function(){
		if(!location.pathname.startsWith("/admin")) return;
		if(location.pathname === "/admin/billing/status.json") return;

		fetch("/admin/billing/status.json", {credentials:"same-origin"})
			.then(function(r){ return r.ok ? r.json() : null; })
			.then(function(d){
				if(!d || !d.enabled || !d.show_banner) return;
				if(document.querySelector(".client-billing-warning-banner-v1")) return;

				var div = document.createElement("div");
				div.className = "client-billing-warning-banner-v1 " + (d.level === "bad" ? "bad" : "warn");
				div.innerHTML =
					'<div>' +
						'<b>' + esc(d.title) + '</b>' +
						'<span>' + esc(d.message) + ' Jatuh tempo: ' + esc(d.due_at || '-') + ' · Tenggang: ' + esc(d.grace_until || '-') + '</span>' +
					'</div>' +
					'<div class="client-billing-warning-actions-v1">' +
						'<a class="btn pri" href="' + esc(d.pay_url || "/admin/billing") + '" target="_blank" rel="noopener">Bayar</a>' +
						'<a class="btn light" href="/admin/billing">Billing</a>' +
					'</div>';

				var nav = document.querySelector(".nav");
				if(nav && nav.parentNode){
					nav.parentNode.insertBefore(div, nav.nextSibling);
				}else{
					document.body.insertBefore(div, document.body.firstChild);
				}
			})
			.catch(function(){});
	});
})();
</script>


<style>
/* CLIENT_MANAGER_BILLING_ENTRY_V1 */
.client-manager-billing-entry-v1{
	background:linear-gradient(135deg,#2563eb,#38bdf8)!important;
	color:#fff!important;
	border-color:rgba(37,99,235,.32)!important;
	box-shadow:0 12px 26px rgba(37,99,235,.18)!important;
}
.client-manager-billing-entry-v1:hover{
	transform:translateY(-1px);
}
</style>


<style>
/* OWNER_NAVBAR_LINK_PAGE_V1 */
.owner-menu-navbar-v1[hidden]{
	display:none!important;
}
.owner-menu-page-v1{
	display:grid!important;
	gap:16px!important;
}
.owner-menu-head-v1 h2{
	margin:0!important;
	font-size:28px!important;
	font-weight:950!important;
	letter-spacing:-.04em!important;
}
.owner-menu-grid-v1{
	grid-template-columns:repeat(3,minmax(0,1fr))!important;
	gap:14px!important;
}
.owner-menu-card-v1{
	display:grid!important;
	gap:8px!important;
	text-decoration:none!important;
	color:inherit!important;
	min-height:150px!important;
	position:relative!important;
	overflow:hidden!important;
}
.owner-menu-card-v1::before{
	content:"";
	position:absolute;
	inset:0;
	background:
		radial-gradient(circle at 0% 0%, rgba(37,99,235,.16), transparent 34%),
		radial-gradient(circle at 100% 0%, rgba(56,189,248,.13), transparent 34%);
	pointer-events:none;
}
.owner-menu-card-v1 b,
.owner-menu-card-v1 span{
	position:relative;
	z-index:1;
}
.owner-menu-card-v1 b{
	font-size:20px!important;
	font-weight:950!important;
	letter-spacing:-.035em!important;
}
.owner-menu-card-v1 span{
	color:#64748b!important;
	font-weight:750!important;
	line-height:1.45!important;
}
html[data-admin-theme="dark"] .owner-menu-card-v1 span,
html.dark .owner-menu-card-v1 span{
	color:#cbd5e1!important;
}
@media(max-width:900px){
	.owner-menu-grid-v1{
		grid-template-columns:1fr!important;
	}
}
</style>
<script>
/* OWNER_NAVBAR_LINK_PAGE_V1 */
(function(){
	function apply(){
		var isOwner = ["vouchergo.biz.id","www.vouchergo.biz.id"].includes(location.hostname);

		document.querySelectorAll('.owner-menu-navbar-v1, a[href="/admin/owner"]').forEach(function(el){
			if(isOwner){
				el.removeAttribute("hidden");
				el.style.removeProperty("display");
			}else{
				el.setAttribute("hidden", "hidden");
				el.style.setProperty("display", "none", "important");
			}
		});

		// Pastikan menu lama owner-specific tidak bocor di non-owner.
		if(!isOwner){
			document.querySelectorAll('a[href="/admin/demo"], a[href="/admin/clients"], a[href="/admin/client-billing"]').forEach(function(el){
				el.remove();
			});
		}
	}

	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", apply);
	}else{
		apply();
	}

	setTimeout(apply, 100);
	setTimeout(apply, 400);
})();
</script>


<style>
/* DASHBOARD_CARDS_CLICKABLE_V1 */
.dashboard-card-clickable-v1{
	cursor:pointer!important;
	position:relative!important;
	transition:
		transform .16s ease,
		box-shadow .16s ease,
		border-color .16s ease!important;
}

.dashboard-card-clickable-v1:hover{
	transform:translateY(-2px)!important;
	border-color:rgba(37,99,235,.34)!important;
	box-shadow:
		0 18px 44px rgba(37,99,235,.12),
		0 8px 20px rgba(15,23,42,.08)!important;
}

.dashboard-card-clickable-v1:active{
	transform:translateY(0)!important;
}

.dashboard-card-clickable-v1::after{
	content:"↗";
	position:absolute;
	right:16px;
	top:14px;
	width:26px;
	height:26px;
	display:grid;
	place-items:center;
	border-radius:999px;
	font-size:12px;
	font-weight:950;
	opacity:.72;
	background:rgba(37,99,235,.10);
	color:#2563eb;
	pointer-events:none;
}

html[data-admin-theme="dark"] .dashboard-card-clickable-v1::after,
html.dark .dashboard-card-clickable-v1::after{
	background:rgba(96,165,250,.14);
	color:#bfdbfe;
}

@media(max-width:760px){
	.dashboard-card-clickable-v1::after{
		right:12px;
		top:12px;
		width:23px;
		height:23px;
		font-size:11px;
	}
}

/* VOUCHERGO_ADMIN_DARK_PALETTE_V1 */
html[data-admin-theme="dark"]{
  --bg: #04142B;
  --card: #231832;
  --card2: #512438;
  --line: rgba(88,72,104,.70);
  --primary: #F6762F;
  --primary2: #BD4C3E;
  --warning: #EEC207;
  --muted: #AFA4BD;
  --text: #F8FAFC;
}

html[data-admin-theme="dark"] body,
html[data-admin-theme="dark"] .admin-body{
  background:
    radial-gradient(circle at 15% 10%, rgba(246,118,47,.16), transparent 25%),
    radial-gradient(circle at 85% 18%, rgba(189,76,62,.18), transparent 30%),
    linear-gradient(145deg, #04142B 0%, #231832 48%, #04142B 100%) !important;
  color: var(--text) !important;
}

html[data-admin-theme="dark"] .card,
html[data-admin-theme="dark"] .panel,
html[data-admin-theme="dark"] .box,
html[data-admin-theme="dark"] .toolbar,
html[data-admin-theme="dark"] .table-wrap{
  background: linear-gradient(145deg, rgba(35,24,50,.96), rgba(81,36,56,.68)) !important;
  border-color: rgba(88,72,104,.70) !important;
  box-shadow: 0 18px 42px rgba(4,20,43,.42) !important;
}

html[data-admin-theme="dark"] .sidebar,
html[data-admin-theme="dark"] aside{
  background: linear-gradient(180deg, #04142B, #231832 56%, #512438) !important;
  border-color: rgba(88,72,104,.70) !important;
}

html[data-admin-theme="dark"] .btn,
html[data-admin-theme="dark"] button,
html[data-admin-theme="dark"] input[type="submit"]{
  background: linear-gradient(135deg, #F6762F, #BD4C3E) !important;
  border-color: rgba(246,118,47,.76) !important;
  color: #fff !important;
}

html[data-admin-theme="dark"] input,
html[data-admin-theme="dark"] select,
html[data-admin-theme="dark"] textarea{
  background: rgba(4,20,43,.70) !important;
  border-color: rgba(88,72,104,.78) !important;
  color: var(--text) !important;
}

html[data-admin-theme="dark"] table,
html[data-admin-theme="dark"] th,
html[data-admin-theme="dark"] td{
  border-color: rgba(88,72,104,.55) !important;
}

html[data-admin-theme="dark"] th{
  background: rgba(81,36,56,.62) !important;
}

html[data-admin-theme="dark"] .muted,
html[data-admin-theme="dark"] small{
  color: rgba(248,250,252,.70) !important;
}
/* /VOUCHERGO_ADMIN_DARK_PALETTE_V1 */
</style>
<script>
/* DASHBOARD_CARDS_CLICKABLE_V1 */
(function(){
	function norm(s){
		return String(s || "").replace(/\s+/g, " ").trim().toLowerCase();
	}

	function pickLink(txt){
		txt = norm(txt);

		if(/^router\b/.test(txt)) return "/admin/routers";
		if(/^paket\b/.test(txt)) return "/admin/packages";
		if(/^pending\b/.test(txt)) return "/admin/transactions?status=pending";
		if(/^paid\b/.test(txt)) return "/admin/transactions?status=paid";
		if(/^transaksi\b/.test(txt)) return "/admin/transactions";
		if(/^whatsapp\b/.test(txt)) return "/admin/whatsapp";
		if(/^database\b/.test(txt)) return "/admin/health";

		return "";
	}

	function apply(){
		if(location.pathname !== "/admin" && location.pathname !== "/admin/") return;

		document.querySelectorAll(".card").forEach(function(card){
			if(card.closest("form")) return;
			if(card.querySelector("a,button,input,select,textarea")) return;
			if(card.dataset.dashboardClickReady === "1") return;

			var text = card.innerText || card.textContent || "";
			var href = pickLink(text);
			if(!href) return;

			card.dataset.dashboardClickReady = "1";
			card.dataset.href = href;
			card.classList.add("dashboard-card-clickable-v1");
			card.setAttribute("role", "link");
			card.setAttribute("tabindex", "0");
			card.setAttribute("title", "Buka " + text.split(/\n/)[0].trim());

			card.addEventListener("click", function(){
				location.href = href;
			});

			card.addEventListener("keydown", function(e){
				if(e.key === "Enter" || e.key === " "){
					e.preventDefault();
					location.href = href;
				}
			});
		});
	}

	if(document.readyState === "loading"){
		document.addEventListener("DOMContentLoaded", apply);
	}else{
		apply();
	}

	setTimeout(apply, 150);
	setTimeout(apply, 500);
})();
</script>


<!-- VOUCHERGO_MOBILE_BOTTOM_NAV_V2 -->
<style>
@media (min-width: 768px){
  .vg-mobile-bottom-nav,
  .vg-mobile-more-sheet,
  .vg-mobile-more-backdrop{display:none!important}
}

@media (max-width: 767.98px){
  body{
    padding-bottom: calc(82px + env(safe-area-inset-bottom, 0px)) !important;
  }

  .vg-mobile-bottom-nav{
    position:fixed;
    left:0;
    right:0;
    bottom:0;
    z-index:99999;
    height:calc(66px + env(safe-area-inset-bottom, 0px));
    padding:7px 10px calc(7px + env(safe-area-inset-bottom, 0px));
    background:rgba(10,15,28,.96);
    border-top:1px solid rgba(148,163,184,.22);
    box-shadow:0 -14px 34px rgba(0,0,0,.38);
    backdrop-filter:blur(16px);
    -webkit-backdrop-filter:blur(16px);
    display:grid;
    grid-template-columns:repeat(5,1fr);
    gap:5px;
  }

  .vg-mobile-bottom-nav a,
  .vg-mobile-bottom-nav button{
    appearance:none;
    border:0;
    background:transparent;
    color:rgba(226,232,240,.70);
    text-decoration:none;
    display:flex;
    flex-direction:column;
    align-items:center;
    justify-content:center;
    gap:4px;
    min-width:0;
    border-radius:16px;
    font:inherit;
    font-weight:800;
    font-size:10.5px;
    line-height:1;
    padding:6px 2px;
    -webkit-tap-highlight-color:transparent;
  }

  .vg-mobile-bottom-nav .vg-ico{
    width:24px;
    height:24px;
    display:grid;
    place-items:center;
    font-size:20px;
    line-height:1;
  }

  .vg-mobile-bottom-nav a.is-active,
  .vg-mobile-bottom-nav button.is-active{
    color:#6ea8ff;
    background:rgba(60,112,255,.14);
  }

  .vg-mobile-bottom-nav a:active,
  .vg-mobile-bottom-nav button:active{
    transform:translateY(1px);
    background:rgba(96,165,250,.16);
  }

  .vg-mobile-more-backdrop{
    position:fixed;
    inset:0;
    z-index:99997;
    background:rgba(2,6,23,.50);
    opacity:0;
    pointer-events:none;
    transition:.16s ease;
  }

  .vg-mobile-more-backdrop.is-open{
    opacity:1;
    pointer-events:auto;
  }

  .vg-mobile-more-sheet{
    position:fixed;
    left:12px;
    right:12px;
    bottom:calc(78px + env(safe-area-inset-bottom, 0px));
    z-index:99998;
    background:rgba(15,23,42,.98);
    border:1px solid rgba(148,163,184,.22);
    box-shadow:0 20px 55px rgba(0,0,0,.45);
    border-radius:24px;
    padding:12px;
    transform:translateY(18px);
    opacity:0;
    pointer-events:none;
    transition:.18s ease;
    backdrop-filter:blur(18px);
    -webkit-backdrop-filter:blur(18px);
  }

  .vg-mobile-more-sheet.is-open{
    transform:translateY(0);
    opacity:1;
    pointer-events:auto;
  }

  .vg-mobile-more-title{
    color:#e5e7eb;
    font-weight:900;
    font-size:14px;
    padding:4px 6px 10px;
  }

  .vg-mobile-more-grid{
    display:grid;
    grid-template-columns:repeat(2,minmax(0,1fr));
    gap:8px;
  }

  .vg-mobile-more-grid a{
    display:flex;
    align-items:center;
    gap:10px;
    min-height:44px;
    padding:10px 12px;
    border-radius:16px;
    text-decoration:none;
    color:#e5e7eb;
    background:rgba(30,41,59,.76);
    border:1px solid rgba(148,163,184,.14);
    font-weight:800;
    font-size:13px;
  }

  .vg-mobile-more-grid a:active{
    transform:translateY(1px);
    background:rgba(59,130,246,.18);
  }
}
</style>

<div class="vg-mobile-more-backdrop" id="vgMobileMoreBackdrop" aria-hidden="true"></div>

<nav class="vg-mobile-more-sheet" id="vgMobileMoreSheet" aria-label="Menu lainnya">
  <div class="vg-mobile-more-title">Menu lainnya</div>
  <div class="vg-mobile-more-grid">
    <a href="/admin/whatsapp">💬 WhatsApp</a>
    <a href="/admin/settings">⚙️ Settings</a>
    <a href="/admin/blocked-phones">🚫 Nomor Diblokir</a>
    <a href="/admin/users">👥 Users</a>
    <a href="/admin/health">🩺 Health</a>
    <a href="/admin/clients">🏢 Client</a>
    <a href="/logout">🚪 Logout</a>
  </div>
</nav>

<nav class="vg-mobile-bottom-nav" aria-label="Navigasi mobile">
  <a href="/admin" data-vg-nav="dashboard"><span class="vg-ico">🏠</span><span>Home</span></a>
  <a href="/admin/transactions" data-vg-nav="transactions"><span class="vg-ico">📄</span><span>Transaksi</span></a>
  <a href="/admin/packages" data-vg-nav="packages"><span class="vg-ico">🎫</span><span>Paket</span></a>
  <a href="/admin/routers" data-vg-nav="routers"><span class="vg-ico">📡</span><span>Router</span></a>
  <button type="button" id="vgMobileMoreBtn" data-vg-nav="more" aria-expanded="false"><span class="vg-ico">☰</span><span>More</span></button>
</nav>

<script>
(function(){
  var path = window.location.pathname || "";
  var active = "dashboard";

  if(path.indexOf("/admin/transactions") === 0) active = "transactions";
  else if(path.indexOf("/admin/packages") === 0) active = "packages";
  else if(path.indexOf("/admin/routers") === 0) active = "routers";
  else if(["/admin/whatsapp","/admin/settings","/admin/users","/admin/health","/admin/clients"].some(function(x){
    return path.indexOf(x) === 0;
  })) active = "more";

  var activeEl = document.querySelector('[data-vg-nav="'+active+'"]');
  if(activeEl) activeEl.classList.add("is-active");

  var btn = document.getElementById("vgMobileMoreBtn");
  var sheet = document.getElementById("vgMobileMoreSheet");
  var backdrop = document.getElementById("vgMobileMoreBackdrop");

  function closeMore(){
    if(!sheet || !backdrop || !btn) return;
    sheet.classList.remove("is-open");
    backdrop.classList.remove("is-open");
    btn.classList.toggle("is-active", active === "more");
    btn.setAttribute("aria-expanded","false");
  }

  function toggleMore(){
    if(!sheet || !backdrop || !btn) return;
    var open = !sheet.classList.contains("is-open");
    sheet.classList.toggle("is-open", open);
    backdrop.classList.toggle("is-open", open);
    btn.classList.toggle("is-active", open || active === "more");
    btn.setAttribute("aria-expanded", open ? "true" : "false");
  }

  if(btn) btn.addEventListener("click", toggleMore);
  if(backdrop) backdrop.addEventListener("click", closeMore);
  document.addEventListener("keydown", function(e){
    if(e.key === "Escape") closeMore();
  });
})();
</script>


<!-- VOUCHERGO_MOBILE_BOTTOM_NAV_MORE_BUTTON_FIX_V1 -->
<style>
@media (max-width: 768px){
  .vg-mobile-bottom-nav #vgMobileMoreBtn{
    appearance:none!important;
    -webkit-appearance:none!important;
    border:0!important;
    outline:0!important;
    box-shadow:none!important;
    background:transparent!important;
    color:rgba(226,232,240,.70)!important;
    text-decoration:none!important;
    display:flex!important;
    flex-direction:column!important;
    align-items:center!important;
    justify-content:center!important;
    gap:4px!important;
    min-width:0!important;
    width:100%!important;
    height:auto!important;
    border-radius:16px!important;
    font-family:inherit!important;
    font-weight:800!important;
    font-size:10.5px!important;
    line-height:1!important;
    padding:6px 2px!important;
    margin:0!important;
    -webkit-tap-highlight-color:transparent!important;
  }

  .vg-mobile-bottom-nav #vgMobileMoreBtn.is-active{
    color:#6ea8ff!important;
    background:rgba(60,112,255,.14)!important;
  }

  .vg-mobile-bottom-nav #vgMobileMoreBtn:active{
    transform:translateY(1px)!important;
    background:rgba(96,165,250,.16)!important;
  }
}
</style>





<!-- VOUCHERGO_MOBILE_ADMINBAR_SIMPLE_V2 -->
<style>
@media (max-width: 768px){
  .mobile-adminbar{
    margin:0 10px 10px!important;
    padding:8px 10px!important;
    min-height:54px!important;
    border-radius:0 0 18px 18px!important;
    background:linear-gradient(135deg, rgba(15,23,42,.84), rgba(15,23,42,.66))!important;
    border:1px solid rgba(148,163,184,.16)!important;
    border-top:0!important;
    box-shadow:0 8px 22px rgba(0,0,0,.18)!important;
    backdrop-filter:blur(14px)!important;
    -webkit-backdrop-filter:blur(14px)!important;
  }

  .mobile-admin-title{
    font-size:18px!important;
    line-height:1.1!important;
    letter-spacing:-.45px!important;
    font-weight:900!important;
    color:#f8fafc!important;
    text-shadow:none!important;
  }

  .mobile-menu-btn,
  .mobile-theme-btn{
    width:38px!important;
    height:38px!important;
    min-width:38px!important;
    min-height:38px!important;
    padding:0!important;
    margin:0!important;
    border-radius:13px!important;
    border:1px solid rgba(226,232,240,.12)!important;
    background:rgba(15,23,42,.42)!important;
    color:#e5e7eb!important;
    box-shadow:none!important;
    outline:0!important;
    transform:none!important;
    display:grid!important;
    place-items:center!important;
    -webkit-tap-highlight-color:transparent!important;
  }

  .mobile-menu-btn{
    font-size:21px!important;
    line-height:1!important;
  }

  .mobile-theme-btn{
    font-size:18px!important;
    line-height:1!important;
  }

  .mobile-theme-btn .sun,
  .mobile-theme-btn .moon{
    filter:none!important;
    text-shadow:none!important;
  }

  .mobile-menu-btn:hover,
  .mobile-menu-btn:focus,
  .mobile-theme-btn:hover,
  .mobile-theme-btn:focus{
    background:rgba(30,41,59,.58)!important;
    border-color:rgba(226,232,240,.18)!important;
    box-shadow:none!important;
    transform:none!important;
  }

  .mobile-menu-btn:active,
  .mobile-theme-btn:active{
    transform:translateY(1px)!important;
    background:rgba(51,65,85,.66)!important;
  }

  html[data-admin-theme="light"] .mobile-adminbar,
  html[data-theme="light"] .mobile-adminbar{
    background:rgba(255,255,255,.82)!important;
    border-color:rgba(15,23,42,.10)!important;
    box-shadow:0 8px 20px rgba(15,23,42,.09)!important;
  }

  html[data-admin-theme="light"] .mobile-admin-title,
  html[data-theme="light"] .mobile-admin-title{
    color:#0f172a!important;
  }

  html[data-admin-theme="light"] .mobile-menu-btn,
  html[data-admin-theme="light"] .mobile-theme-btn,
  html[data-theme="light"] .mobile-menu-btn,
  html[data-theme="light"] .mobile-theme-btn{
    background:rgba(15,23,42,.06)!important;
    color:#0f172a!important;
    border-color:rgba(15,23,42,.10)!important;
  }
}
</style>


<!-- VOUCHERGO_MOBILE_ADMINBAR_KILL_ORANGE_V3 -->
<style>
@media (max-width: 768px){
  .mobile-adminbar .mobile-menu-btn,
  .mobile-adminbar .mobile-theme-btn{
    background-image:none!important;
    background-color:rgba(15,23,42,.46)!important;
    background:rgba(15,23,42,.46)!important;
    color:#e5e7eb!important;
    border-color:rgba(226,232,240,.14)!important;
    box-shadow:none!important;
  }

  .mobile-adminbar .mobile-menu-btn::before,
  .mobile-adminbar .mobile-menu-btn::after,
  .mobile-adminbar .mobile-theme-btn::before,
  .mobile-adminbar .mobile-theme-btn::after{
    content:none!important;
    display:none!important;
    background:none!important;
    background-image:none!important;
    box-shadow:none!important;
  }

  .mobile-adminbar .mobile-menu-btn:hover,
  .mobile-adminbar .mobile-menu-btn:focus,
  .mobile-adminbar .mobile-theme-btn:hover,
  .mobile-adminbar .mobile-theme-btn:focus{
    background-image:none!important;
    background-color:rgba(30,41,59,.62)!important;
    background:rgba(30,41,59,.62)!important;
    box-shadow:none!important;
  }

  html[data-admin-theme="light"] .mobile-adminbar .mobile-menu-btn,
  html[data-admin-theme="light"] .mobile-adminbar .mobile-theme-btn,
  html[data-theme="light"] .mobile-adminbar .mobile-menu-btn,
  html[data-theme="light"] .mobile-adminbar .mobile-theme-btn{
    background-image:none!important;
    background-color:rgba(15,23,42,.07)!important;
    background:rgba(15,23,42,.07)!important;
    color:#0f172a!important;
    border-color:rgba(15,23,42,.12)!important;
    box-shadow:none!important;
  }
}
</style>


<!-- VOUCHERGO_MOBILE_REVENUE_COMPACT_V1 -->
<style>
@media (max-width: 768px){
  body:has(.dashboard-premium-stats-v1) .dashboard-revenue-card-v1{
    padding:14px 13px!important;
    border-radius:20px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dashboard-revenue-card-v1 h3{
    margin:0 0 10px!important;
    font-size:17px!important;
    line-height:1.15!important;
  }

  body:has(.dashboard-premium-stats-v1) .dashboard-revenue-list-v1{
    display:grid!important;
    grid-template-columns:repeat(3,minmax(0,1fr))!important;
    gap:8px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1{
    min-height:76px!important;
    padding:10px 8px!important;
    border-radius:16px!important;
    display:flex!important;
    flex-direction:column!important;
    align-items:flex-start!important;
    justify-content:center!important;
    gap:4px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .dash-dot-v1{
    width:10px!important;
    height:10px!important;
    min-width:10px!important;
    margin:0 0 2px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .num{
    font-size:17px!important;
    line-height:1.05!important;
    letter-spacing:-.7px!important;
    white-space:nowrap!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .label{
    font-size:11px!important;
    line-height:1.1!important;
    white-space:nowrap!important;
  }
}

@media (max-width: 380px){
  body:has(.dashboard-premium-stats-v1) .dashboard-revenue-list-v1{
    gap:6px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1{
    padding:9px 6px!important;
    min-height:70px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .num{
    font-size:15px!important;
  }

  body:has(.dashboard-premium-stats-v1) .dash-income-item-v1 .label{
    font-size:10px!important;
  }
}
</style>


<!-- VOUCHERGO_MOBILE_TOPCARDS_COMPACT_V1 -->
<style>
@media (max-width: 768px){
  .dashboard-top-cards-v1{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:8px!important;
    margin:0 0 12px!important;
  }

  .dashboard-top-card-v1{
    min-height:76px!important;
    padding:11px 12px!important;
    border-radius:17px!important;
    overflow:hidden!important;
  }

  .dashboard-top-card-v1 .label{
    font-size:13px!important;
    line-height:1.1!important;
    margin-bottom:7px!important;
    letter-spacing:-.15px!important;
  }

  .dashboard-top-card-v1 .num{
    font-size:31px!important;
    line-height:1!important;
    letter-spacing:-1px!important;
  }

  .dashboard-top-card-v1 .hint{
    width:24px!important;
    height:24px!important;
    right:10px!important;
    top:10px!important;
    font-size:12px!important;
    opacity:.72!important;
  }

  .dashboard-top-card-v1:before{
    width:70px!important;
    height:70px!important;
    right:-18px!important;
    bottom:-20px!important;
    opacity:.62!important;
  }
}

@media (max-width: 380px){
  .dashboard-top-cards-v1{
    gap:7px!important;
  }

  .dashboard-top-card-v1{
    min-height:70px!important;
    padding:10px!important;
    border-radius:15px!important;
  }

  .dashboard-top-card-v1 .label{
    font-size:12px!important;
    margin-bottom:6px!important;
  }

  .dashboard-top-card-v1 .num{
    font-size:28px!important;
  }

  .dashboard-top-card-v1 .hint{
    width:22px!important;
    height:22px!important;
    right:8px!important;
    top:8px!important;
  }

  .dashboard-top-card-v1:before{
    width:62px!important;
    height:62px!important;
  }
}
</style>

</body>
</html>`

	t := template.Must(template.New("admin-layout").Parse(layout))
	if err := t.Execute(w, map[string]any{
		"Title":    title,
		"Body":     body,
		"AdminCSS": template.CSS(adminPageCSS),
	}); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
