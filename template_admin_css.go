package main

const adminPageCSS = `:root{--bg:#f5f7fb;--card:#fff;--ink:#0f172a;--muted:#64748b;--pri:#2563eb;--pri2:#1d4ed8;--danger:#dc2626;--ok:#16a34a;--line:#e5e7eb;--shadow:0 18px 50px rgba(15,23,42,.08)}
*{box-sizing:border-box}
body{margin:0;font-family:system-ui,-apple-system,Segoe UI,Arial,sans-serif;background:var(--bg);color:var(--ink);font-size:15px}
a{color:var(--pri);text-decoration:none}
.wrap{max-width:1440px;margin:auto;padding:18px}
.nav{position:sticky;top:0;z-index:20;display:flex;gap:14px;align-items:center;justify-content:space-between;padding:10px 0 16px;background:linear-gradient(180deg,var(--bg) 75%,rgba(245,247,251,0))}
.brand{font-weight:900;font-size:24px;color:#0f172a;letter-spacing:-.03em}
.menu{display:flex;gap:10px;flex-wrap:wrap;align-items:center}
.menu a,.btn{display:inline-flex;align-items:center;justify-content:center;gap:8px;border:0;border-radius:14px;padding:11px 16px;background:#111827;color:#fff;font-weight:800;cursor:pointer;line-height:1.1;min-height:42px}
.menu a{background:#fff;color:#111827;border:1px solid var(--line);box-shadow:0 6px 18px rgba(15,23,42,.04)}
.btn.pri{background:linear-gradient(135deg,var(--pri),var(--pri2))}
.btn.danger{background:var(--danger)}
.btn.light{background:#fff;color:#111827;border:1px solid var(--line)}
.btn.ok{background:var(--ok)}
.card{background:var(--card);border:1px solid var(--line);border-radius:24px;padding:22px;box-shadow:var(--shadow);margin-bottom:16px}
.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(230px,1fr));gap:16px}
.field{display:flex;flex-direction:column;gap:7px;margin-bottom:12px}
.field label{font-size:13px;color:var(--muted);font-weight:800}
input,select{width:100%;padding:12px 14px;border:1px solid var(--line);border-radius:14px;background:#fff;font-size:15px;outline:none}
.tablewrap{width:100%;overflow:auto;border-radius:18px;border:1px solid var(--line)}
table{width:100%;border-collapse:collapse;background:#fff}
th,td{padding:14px 12px;border-bottom:1px solid var(--line);text-align:left;font-size:14px;vertical-align:middle}
th{color:var(--muted);font-size:12px;text-transform:uppercase;letter-spacing:.04em;background:#f8fafc}
.badge{display:inline-flex;border-radius:999px;padding:5px 10px;font-size:12px;font-weight:900;background:#eef2ff;color:#3730a3}
.paid,.success{background:#dcfce7;color:#166534}
.pending{background:#fef9c3;color:#854d0e}
.failed{background:#fee2e2;color:#991b1b}
.skipped{background:#e5e7eb;color:#374151}
.muted{color:var(--muted)}
.actions{display:flex;gap:8px;flex-wrap:wrap;align-items:center}
.hero{padding:30px;border-radius:28px;background:linear-gradient(135deg,#0f172a,#1d4ed8);color:#fff;box-shadow:var(--shadow)}
.price{font-size:26px;font-weight:950;letter-spacing:-.04em}
.notice{border-radius:16px;padding:12px 14px;background:#f8fafc;border:1px solid var(--line)}
@media(max-width:760px){
  .wrap{padding:12px}
  .nav{align-items:flex-start;flex-direction:column;padding-top:8px}
  .brand{font-size:22px}
  .menu{width:100%;display:grid;grid-template-columns:repeat(3,1fr);gap:8px}
  .menu a{padding:10px 8px;font-size:13px;border-radius:13px;min-height:40px}
  .card{border-radius:18px;padding:16px}
  .grid{grid-template-columns:1fr;gap:12px}
  th,td{padding:12px 10px;white-space:nowrap}
  .btn{padding:10px 12px;border-radius:13px;min-height:40px}
}

/* ADMIN_UI_POLISH_V1 */
body{
  background:
    radial-gradient(circle at top left, rgba(37,99,235,.08), transparent 28%),
    linear-gradient(180deg,#f8fafc,#eef2ff)!important;
}

.wrap{
  max-width:1500px!important;
  padding:16px!important;
}

.nav{
  position:sticky!important;
  top:0!important;
  z-index:100!important;
  margin:-16px -16px 16px!important;
  padding:12px 16px!important;
  background:rgba(248,250,252,.88)!important;
  backdrop-filter:blur(14px)!important;
  border-bottom:1px solid rgba(226,232,240,.9)!important;
}

.brand{
  font-size:22px!important;
  font-weight:950!important;
}

.menu{
  gap:8px!important;
}

.menu a{
  min-height:38px!important;
  padding:9px 13px!important;
  border-radius:13px!important;
  font-size:14px!important;
  font-weight:850!important;
}

.card{
  border-radius:20px!important;
  padding:18px!important;
  box-shadow:0 14px 34px rgba(15,23,42,.07)!important;
}

.card h1,
.card h2{
  margin-top:0!important;
  letter-spacing:-.04em!important;
}

.grid{
  gap:14px!important;
}

.tablewrap{
  border-radius:16px!important;
  overflow:auto!important;
  background:#fff!important;
}

table{
  min-width:900px!important;
}

th{
  position:sticky!important;
  top:0!important;
  z-index:2!important;
  padding:11px 12px!important;
  font-size:11px!important;
  background:#f8fafc!important;
}

td{
  padding:11px 12px!important;
  font-size:13px!important;
}

tr:hover td{
  background:#f8fafc!important;
}

.btn{
  min-height:38px!important;
  padding:9px 13px!important;
  border-radius:12px!important;
  font-size:13px!important;
}

input,
select{
  min-height:42px!important;
  border-radius:12px!important;
}

.badge{
  padding:4px 9px!important;
}

.actions{
  gap:7px!important;
}

@media(max-width:760px){
  .wrap{
    padding:10px!important;
  }

  .nav{
    margin:-10px -10px 12px!important;
    padding:10px!important;
    gap:8px!important;
  }

  .brand{
    font-size:18px!important;
    line-height:1.1!important;
  }

  .menu{
    width:100%!important;
    display:flex!important;
    flex-wrap:nowrap!important;
    overflow-x:auto!important;
    gap:7px!important;
    padding-bottom:2px!important;
    scrollbar-width:none!important;
  }

  .menu::-webkit-scrollbar{
    display:none!important;
  }

  .menu a{
    flex:0 0 auto!important;
    min-height:34px!important;
    padding:8px 11px!important;
    font-size:12px!important;
    border-radius:999px!important;
    white-space:nowrap!important;
  }

  .card{
    padding:13px!important;
    border-radius:16px!important;
    margin-bottom:10px!important;
  }

  .card h1,
  .card h2{
    font-size:20px!important;
    margin-bottom:12px!important;
  }

  .grid{
    grid-template-columns:1fr!important;
    gap:10px!important;
  }

  .tablewrap{
    margin-left:-2px!important;
    margin-right:-2px!important;
    border-radius:14px!important;
  }

  table{
    min-width:760px!important;
  }

  th,
  td{
    padding:9px 10px!important;
    font-size:12px!important;
  }

  .btn{
    min-height:34px!important;
    padding:8px 10px!important;
    font-size:12px!important;
    border-radius:10px!important;
  }

  input,
  select{
    min-height:40px!important;
    font-size:14px!important;
  }
}


/* ADMIN_DARKMODE_V1 */
.admin-theme-toggle{
  display:inline-flex;
  align-items:center;
  justify-content:space-between;
  gap:4px;
  width:54px;
  min-width:54px;
  height:38px;
  padding:0 8px;
  border-radius:999px;
  border:1px solid var(--line);
  background:#fff;
  color:#111827;
  box-shadow:0 6px 18px rgba(15,23,42,.04);
  cursor:pointer;
  font-size:13px;
  font-weight:850;
}

html[data-admin-theme="light"] .moon{
  opacity:.35;
}

html[data-admin-theme="dark"] .sun{
  opacity:.35;
}

html[data-admin-theme="dark"]{
  --bg:#020617;
  --card:#0f172a;
  --ink:#e5e7eb;
  --muted:#94a3b8;
  --line:#1f2937;
  --shadow:0 18px 50px rgba(0,0,0,.35);
}

html[data-admin-theme="dark"] body{
  background:
    radial-gradient(circle at top left, rgba(37,99,235,.16), transparent 28%),
    radial-gradient(circle at top right, rgba(6,182,212,.10), transparent 30%),
    linear-gradient(180deg,#020617,#0b1120)!important;
  color:#e5e7eb!important;
}

html[data-admin-theme="dark"] .nav{
  background:rgba(2,6,23,.88)!important;
  border-bottom-color:#1f2937!important;
}

html[data-admin-theme="dark"] .brand,
html[data-admin-theme="dark"] h1,
html[data-admin-theme="dark"] h2,
html[data-admin-theme="dark"] h3{
  color:#f8fafc!important;
}

html[data-admin-theme="dark"] .card{
  background:linear-gradient(180deg,#0f172a,#0b1120)!important;
  border-color:#1f2937!important;
  box-shadow:0 18px 45px rgba(0,0,0,.32)!important;
}

html[data-admin-theme="dark"] .menu a,
html[data-admin-theme="dark"] .btn.light,
html[data-admin-theme="dark"] .admin-theme-toggle,
html[data-admin-theme="dark"] input,
html[data-admin-theme="dark"] select,
html[data-admin-theme="dark"] .notice{
  background:#111827!important;
  color:#e5e7eb!important;
  border-color:#1f2937!important;
}

html[data-admin-theme="dark"] input::placeholder{
  color:#64748b!important;
}

html[data-admin-theme="dark"] .tablewrap,
html[data-admin-theme="dark"] table{
  background:#0f172a!important;
  border-color:#1f2937!important;
}

html[data-admin-theme="dark"] th{
  background:#111827!important;
  color:#94a3b8!important;
  border-bottom-color:#1f2937!important;
}

html[data-admin-theme="dark"] td{
  color:#e5e7eb!important;
  border-bottom-color:#1f2937!important;
}

html[data-admin-theme="dark"] tr:hover td{
  background:#111827!important;
}

html[data-admin-theme="dark"] .muted{
  color:#94a3b8!important;
}

html[data-admin-theme="dark"] .badge{
  border:1px solid rgba(255,255,255,.06)!important;
}

html[data-admin-theme="dark"] .pending{
  background:rgba(234,179,8,.18)!important;
  color:#fde68a!important;
}

html[data-admin-theme="dark"] .paid,
html[data-admin-theme="dark"] .success{
  background:rgba(34,197,94,.18)!important;
  color:#86efac!important;
}

html[data-admin-theme="dark"] .failed{
  background:rgba(239,68,68,.18)!important;
  color:#fca5a5!important;
}

html[data-admin-theme="dark"] .skipped{
  background:rgba(148,163,184,.16)!important;
  color:#cbd5e1!important;
}

@media(max-width:760px){
  .admin-theme-toggle{
    flex:0 0 auto!important;
    width:48px!important;
    min-width:48px!important;
    height:34px!important;
    padding:0 7px!important;
    font-size:11px!important;
  }
}


/* ADMIN_TRANSACTIONS_UI_V1 */
.tablewrap{
  box-shadow:0 10px 26px rgba(15,23,42,.05)!important;
}

table{
  border-spacing:0!important;
}

th:first-child,
td:first-child{
  padding-left:16px!important;
}

th:last-child,
td:last-child{
  padding-right:16px!important;
}

td b{
  letter-spacing:-.02em!important;
}

td small{
  color:#64748b!important;
  display:inline-block!important;
  max-width:520px!important;
  line-height:1.35!important;
}

.badge{
  text-transform:lowercase!important;
  white-space:nowrap!important;
}

.actions form{
  margin:0!important;
}

.actions .btn,
td .btn{
  min-width:88px!important;
  justify-content:center!important;
}

td .actions{
  align-items:center!important;
  flex-wrap:nowrap!important;
}

@media(min-width:900px){
  .tablewrap table{
    font-size:13px!important;
  }

  .tablewrap td:nth-child(1){
    min-width:210px!important;
  }

  .tablewrap td:nth-child(2){
    min-width:130px!important;
  }

  .tablewrap td:nth-child(6){
    min-width:120px!important;
  }

  .tablewrap td:nth-child(7){
    min-width:260px!important;
  }

  .tablewrap td:nth-child(8){
    min-width:210px!important;
  }
}

@media(max-width:760px){
  .tablewrap{
    border-radius:15px!important;
    overflow-x:auto!important;
    -webkit-overflow-scrolling:touch!important;
  }

  .tablewrap::after{
    content:"Geser tabel ke kanan →";
    display:block;
    position:sticky;
    left:0;
    bottom:0;
    padding:8px 10px;
    font-size:11px;
    color:#64748b;
    background:rgba(248,250,252,.94);
    border-top:1px solid #e5e7eb;
  }

  html[data-admin-theme="dark"] .tablewrap::after{
    background:rgba(15,23,42,.94)!important;
    color:#94a3b8!important;
    border-top-color:#1f2937!important;
  }

  th{
    font-size:10.5px!important;
  }

  td{
    font-size:12px!important;
    line-height:1.25!important;
  }

  td small{
    max-width:260px!important;
    font-size:10.5px!important;
  }

  .badge{
    padding:4px 8px!important;
    font-size:10.5px!important;
  }

  td .actions{
    gap:6px!important;
  }

  td .btn{
    min-width:74px!important;
    min-height:32px!important;
    padding:7px 8px!important;
    font-size:11px!important;
  }
}


/* ADMIN_WA_STATUS_UI_V1 */
.sent{
  background:#dcfce7!important;
  color:#166534!important;
}
.sending{
  background:#dbeafe!important;
  color:#1d4ed8!important;
}
html[data-admin-theme="dark"] .sent{
  background:rgba(34,197,94,.18)!important;
  color:#86efac!important;
}
html[data-admin-theme="dark"] .sending{
  background:rgba(59,130,246,.18)!important;
  color:#93c5fd!important;
}


/* ADMIN_MOBILE_TABLET_FIX_V2 */

/* desktop/tablet navbar: jangan wrap berantakan */
.nav{
  display:grid!important;
  grid-template-columns:auto minmax(0,1fr)!important;
  align-items:center!important;
  gap:12px!important;
}

.menu{
  min-width:0!important;
  display:flex!important;
  align-items:center!important;
  justify-content:flex-end!important;
  flex-wrap:nowrap!important;
  overflow-x:auto!important;
  gap:8px!important;
  padding-bottom:2px!important;
  scrollbar-width:none!important;
}

.menu::-webkit-scrollbar{
  display:none!important;
}

.menu a,
.admin-theme-toggle{
  flex:0 0 auto!important;
  white-space:nowrap!important;
}

.admin-theme-toggle{
  order:99!important;
}

/* tablet */
@media(max-width:1180px){
  .wrap{
    max-width:100%!important;
    padding:12px!important;
  }

  .nav{
    margin:-12px -12px 12px!important;
    padding:10px 12px!important;
    grid-template-columns:1fr!important;
    gap:8px!important;
  }

  .brand{
    font-size:20px!important;
    line-height:1.1!important;
  }

  .menu{
    justify-content:flex-start!important;
    width:100%!important;
    gap:7px!important;
  }

  .menu a{
    min-height:36px!important;
    padding:8px 12px!important;
    border-radius:999px!important;
    font-size:13px!important;
  }

  .admin-theme-toggle{
    width:50px!important;
    min-width:50px!important;
    height:36px!important;
  }

  .card{
    border-radius:18px!important;
    padding:15px!important;
  }

  .card h1,
  .card h2{
    font-size:24px!important;
    line-height:1.1!important;
  }

  .tablewrap{
    width:100%!important;
    max-width:100%!important;
    overflow-x:auto!important;
    -webkit-overflow-scrolling:touch!important;
    border-radius:15px!important;
  }

  table{
    min-width:980px!important;
  }

  th{
    padding:10px 11px!important;
    font-size:10.5px!important;
  }

  td{
    padding:10px 11px!important;
    font-size:12.5px!important;
  }

  td small{
    font-size:11px!important;
    line-height:1.25!important;
  }

  .btn{
    min-height:34px!important;
    padding:8px 11px!important;
    font-size:12px!important;
    border-radius:11px!important;
  }
}

/* mobile */
@media(max-width:760px){
  body{
    font-size:14px!important;
  }

  .wrap{
    padding:9px!important;
  }

  .nav{
    margin:-9px -9px 10px!important;
    padding:9px!important;
    position:sticky!important;
    top:0!important;
  }

  .brand{
    font-size:19px!important;
    max-width:100%!important;
  }

  .menu{
    gap:6px!important;
    padding-bottom:3px!important;
  }

  .menu a{
    min-height:32px!important;
    padding:7px 10px!important;
    border-radius:999px!important;
    font-size:12px!important;
  }

  .admin-theme-toggle{
    width:46px!important;
    min-width:46px!important;
    height:32px!important;
    padding:0 7px!important;
    font-size:10px!important;
  }

  .card{
    padding:12px!important;
    border-radius:15px!important;
    margin-bottom:10px!important;
  }

  .card h1,
  .card h2{
    font-size:21px!important;
    margin-bottom:10px!important;
  }

  .grid{
    grid-template-columns:1fr!important;
    gap:9px!important;
  }

  .tablewrap{
    margin-left:0!important;
    margin-right:0!important;
    border-radius:13px!important;
  }

  table{
    min-width:880px!important;
  }

  th:first-child,
  td:first-child{
    padding-left:12px!important;
  }

  th:last-child,
  td:last-child{
    padding-right:12px!important;
  }

  th{
    padding:8px 9px!important;
    font-size:10px!important;
  }

  td{
    padding:9px!important;
    font-size:12px!important;
    vertical-align:top!important;
  }

  td b{
    font-size:12px!important;
  }

  td small{
    font-size:10.5px!important;
    max-width:220px!important;
  }

  .badge{
    padding:4px 8px!important;
    font-size:10px!important;
  }

  .actions{
    gap:5px!important;
    flex-wrap:nowrap!important;
  }

  .actions .btn,
  td .btn{
    min-width:68px!important;
    min-height:28px!important;
    padding:7px 8px!important;
    font-size:10.5px!important;
    border-radius:9px!important;
  }

  input,
  select{
    min-height:38px!important;
    font-size:14px!important;
    border-radius:10px!important;
  }
}

/* layar sangat kecil */
@media(max-width:420px){
  .brand{
    font-size:18px!important;
  }

  .menu a{
    font-size:11px!important;
    padding:7px 9px!important;
  }

  table{
    min-width:820px!important;
  }

  .card h1,
  .card h2{
    font-size:20px!important;
  }
}

/* dark mode adjustment */
html[data-admin-theme="dark"] .nav{
  background:rgba(2,6,23,.90)!important;
}

html[data-admin-theme="dark"] .tablewrap::after{
  background:rgba(15,23,42,.96)!important;
  border-color:#1f2937!important;
}

html[data-admin-theme="dark"] .menu a,
html[data-admin-theme="dark"] .admin-theme-toggle{
  background:#111827!important;
  border-color:#1f2937!important;
  color:#e5e7eb!important;
}


/* ADMIN_USERS_ROLE_V1 */
.notice{
  margin:10px 0 14px!important;
}

.tablewrap select,
.tablewrap input{
  min-width:150px!important;
}

@media(max-width:760px){
  .tablewrap select,
  .tablewrap input{
    min-width:130px!important;
  }
}


/* ADMIN_USER_PASS_EYE_CONFIRM_V1 */
.pass-wrap{
  display:flex!important;
  align-items:center!important;
  gap:6px!important;
  width:100%!important;
}

.pass-wrap input{
  min-width:0!important;
  flex:1 1 auto!important;
}

.pass-wrap.mini{
  margin-bottom:6px!important;
}

.pass-eye{
  flex:0 0 auto!important;
  width:40px!important;
  height:40px!important;
  border:1px solid var(--line)!important;
  border-radius:12px!important;
  background:#fff!important;
  cursor:pointer!important;
  display:inline-flex!important;
  align-items:center!important;
  justify-content:center!important;
  font-size:15px!important;
}

html[data-admin-theme="dark"] .pass-eye{
  background:#111827!important;
  color:#e5e7eb!important;
  border-color:#1f2937!important;
}

@media(max-width:760px){
  .pass-eye{
    width:36px!important;
    height:36px!important;
    border-radius:10px!important;
    font-size:13px!important;
  }

  .pass-wrap.mini input{
    min-width:170px!important;
  }
}


/* ADMIN_WHATSAPP_MENU_V1 */
.wa-status-grid{
  display:grid!important;
  grid-template-columns:repeat(3,minmax(0,1fr))!important;
  gap:12px!important;
  margin:14px 0!important;
}

.wa-stat{
  border:1px solid var(--line)!important;
  border-radius:16px!important;
  padding:14px!important;
  background:rgba(255,255,255,.65)!important;
}

.wa-stat b{
  display:block!important;
  color:var(--muted)!important;
  font-size:12px!important;
  margin-bottom:7px!important;
}

.wa-stat span{
  font-weight:900!important;
}

.wa-qr-box{
  text-align:center!important;
  border:1px dashed var(--line)!important;
  border-radius:20px!important;
  padding:18px!important;
  margin:14px 0!important;
}

.wa-qr-img{
  width:360px!important;
  max-width:100%!important;
  background:#fff!important;
  border-radius:18px!important;
  padding:14px!important;
}

.wa-actions form{
  margin:0!important;
}

.wa-raw{
  overflow:auto!important;
  white-space:pre-wrap!important;
  word-break:break-word!important;
  border:1px solid var(--line)!important;
  border-radius:16px!important;
  padding:14px!important;
  background:#0f172a!important;
  color:#e5e7eb!important;
  font-size:12px!important;
  line-height:1.45!important;
}

html[data-admin-theme="dark"] .wa-stat{
  background:#111827!important;
  border-color:#1f2937!important;
}

@media(max-width:760px){
  .wa-status-grid{
    grid-template-columns:1fr!important;
  }

  .wa-qr-img{
    width:100%!important;
  }
}


/* ADMIN_HIDE_SUPER_DEFAULT_V1 */
.super-only,
a[href="/admin/whatsapp"].super-only,
a[href="/admin/settings"].super-only,
a[href="/admin/users"].super-only{
  display:none!important;
}

html.admin-role-superadmin .super-only{
  display:inline-flex!important;
}

html.admin-role-router .super-only,
html.admin-role-router a[href="/admin/whatsapp"],
html.admin-role-router a[href="/admin/settings"],
html.admin-role-router a[href="/admin/users"]{
  display:none!important;
}


/* ADMIN_COMPACT_STAGE1D_V1 */
.card{
  padding:18px!important;
  border-radius:18px!important;
}

.grid{
  gap:12px!important;
}

.price{
  font-size:26px!important;
  line-height:1.05!important;
}

.field label,
.card b,
.card h2,
.card h3{
  margin-bottom:8px!important;
}

.field input,
.field select,
.field textarea{
  min-height:42px!important;
  padding:10px 12px!important;
}

.actions{
  gap:8px!important;
}

.btn{
  padding:10px 14px!important;
  min-height:40px!important;
}

.wa-raw{
  max-height:260px!important;
  overflow:auto!important;
  font-size:12px!important;
  line-height:1.45!important;
}

@media(max-width:760px){
  .card{
    padding:14px!important;
    border-radius:16px!important;
  }

  .price{
    font-size:24px!important;
  }

  .grid{
    gap:10px!important;
  }
}


/* ADMIN_ROUTER_SETTINGS_COMPACT_V1 */
.secret-field{
  position:relative!important;
}
.secret-field input{
  padding-right:52px!important;
}
.secret-field .pass-eye{
  position:absolute!important;
  right:8px!important;
  bottom:7px!important;
  width:36px!important;
  height:36px!important;
  border:1px solid var(--line)!important;
  border-radius:12px!important;
  background:var(--card)!important;
  color:var(--text)!important;
  cursor:pointer!important;
}

.router-checks{
  display:grid!important;
  grid-template-columns:repeat(3,minmax(0,1fr))!important;
  gap:10px!important;
  align-items:center!important;
}

.router-checks label{
  display:flex!important;
  align-items:center!important;
  justify-content:space-between!important;
  gap:10px!important;
  min-height:46px!important;
  padding:10px 12px!important;
  border:1px solid var(--line)!important;
  border-radius:14px!important;
  background:rgba(255,255,255,.55)!important;
  font-weight:800!important;
}

.router-checks input[type="checkbox"]{
  width:22px!important;
  height:22px!important;
}

.router-form-compact{
  display:grid!important;
  grid-template-columns:repeat(5,minmax(0,1fr))!important;
  gap:12px!important;
}

.router-form-compact .router-password{
  grid-column:span 2!important;
}

.router-form-compact .router-checks{
  grid-column:span 2!important;
}

.router-form-compact .router-submit{
  align-self:end!important;
}

@media(max-width:980px){
  .router-form-compact{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }
  .router-form-compact .router-password,
  .router-form-compact .router-checks{
    grid-column:span 2!important;
  }
}

@media(max-width:640px){
  .router-form-compact,
  .router-checks{
    grid-template-columns:1fr!important;
  }
  .router-form-compact .router-password,
  .router-form-compact .router-checks{
    grid-column:span 1!important;
  }
}


/* ADMIN_ROUTER_COMPACT_V2 */
.router-form-compact-v2{
  margin:0!important;
}

.router-form-compact-v2 .grid{
  display:grid!important;
  grid-template-columns:repeat(4,minmax(0,1fr))!important;
  gap:12px!important;
  align-items:end!important;
}

.router-form-compact-v2 .field{
  margin:0!important;
}

.router-form-compact-v2 input:not([type="checkbox"]){
  min-height:40px!important;
  height:40px!important;
  padding:9px 12px!important;
}

.router-form-compact-v2 input[type="checkbox"]{
  appearance:auto!important;
  -webkit-appearance:checkbox!important;
  width:18px!important;
  height:18px!important;
  min-width:18px!important;
  min-height:18px!important;
  max-width:18px!important;
  max-height:18px!important;
  padding:0!important;
  margin:0!important;
  display:inline-block!important;
  vertical-align:middle!important;
  accent-color:#2563eb!important;
  transform:none!important;
}

.router-form-compact-v2 .field:has(input[type="checkbox"]){
  min-height:40px!important;
  height:40px!important;
  padding:8px 12px!important;
  border:1px solid var(--line)!important;
  border-radius:14px!important;
  display:flex!important;
  align-items:center!important;
  justify-content:space-between!important;
  gap:10px!important;
  background:rgba(255,255,255,.55)!important;
}

.router-form-compact-v2 .field:has(input[type="checkbox"]) label{
  margin:0!important;
  font-size:13px!important;
  line-height:1.2!important;
}

.router-form-compact-v2 button[type="submit"],
.router-form-compact-v2 .btn{
  min-height:40px!important;
  height:40px!important;
  padding:9px 14px!important;
}

@media(max-width:980px){
  .router-form-compact-v2 .grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }
}

@media(max-width:640px){
  .router-form-compact-v2 .grid{
    grid-template-columns:1fr!important;
  }
}


/* ADMIN_MOBILE_DRAWER_ACTIVE_V1 */
.mobile-adminbar,
.mobile-drawer-overlay,
.mobile-drawer{
  display:none;
}

.mobile-adminbar{
  position:sticky;
  top:0;
  z-index:80;
  align-items:center;
  justify-content:space-between;
  gap:10px;
  padding:10px 0 12px;
  background:linear-gradient(180deg,var(--bg) 78%,rgba(245,247,251,0));
}

.mobile-menu-btn,
.mobile-drawer-close{
  width:44px;
  height:44px;
  border:1px solid var(--line);
  border-radius:16px;
  background:var(--card);
  color:var(--text);
  font-size:22px;
  font-weight:900;
  cursor:pointer;
  box-shadow:0 10px 24px rgba(15,23,42,.08);
}

.mobile-admin-title{
  flex:1;
  text-decoration:none;
  color:var(--text);
  font-size:20px;
  font-weight:950;
  letter-spacing:-.04em;
  white-space:nowrap;
  overflow:hidden;
  text-overflow:ellipsis;
}

.mobile-drawer-overlay{
  position:fixed;
  inset:0;
  z-index:90;
  background:rgba(2,6,23,.54);
  backdrop-filter:blur(4px);
}

.mobile-drawer{
  position:fixed;
  top:0;
  left:-320px;
  bottom:0;
  z-index:100;
  width:min(310px,86vw);
  padding:16px;
  background:var(--card);
  border-right:1px solid var(--line);
  box-shadow:24px 0 60px rgba(15,23,42,.22);
  transition:left .22s ease;
  flex-direction:column;
  gap:14px;
}

.mobile-drawer-head{
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:10px;
  padding-bottom:12px;
  border-bottom:1px solid var(--line);
}

.mobile-drawer-head strong{
  font-size:20px;
  font-weight:950;
  letter-spacing:-.04em;
}

.mobile-drawer-links{
  display:flex;
  flex-direction:column;
  gap:10px;
  overflow:auto;
  padding-bottom:20px;
}

.mobile-drawer-links a{
  display:flex;
  align-items:center;
  min-height:46px;
  padding:12px 14px;
  border:1px solid var(--line);
  border-radius:16px;
  text-decoration:none;
  color:var(--text);
  font-weight:900;
  background:rgba(248,250,252,.75);
}

body.admin-drawer-open{
  overflow:hidden;
}

body.admin-drawer-open .mobile-drawer-overlay{
  display:block;
}

body.admin-drawer-open .mobile-drawer{
  left:0;
  display:flex;
}

html.admin-role-router .mobile-drawer-links .super-only{
  display:none!important;
}

html[data-admin-theme="dark"] .mobile-adminbar{
  background:linear-gradient(180deg,var(--bg) 78%,rgba(2,6,23,0));
}

html[data-admin-theme="dark"] .mobile-drawer-links a{
  background:rgba(15,23,42,.78);
}

@media(max-width:760px), (min-width:768px) and (max-width:1199px){
  .mobile-adminbar{
    display:flex!important;
  }

  .nav{
    display:none!important;
  }

  body{
    padding-top:0!important;
  }
}


/* ADMIN_MOBILE_SUPER_COMPACT_V1 */
@media(max-width:760px){
  body{
    padding-left:10px!important;
    padding-right:10px!important;
  }

  .card{
    padding:12px!important;
    border-radius:16px!important;
    margin:8px 0!important;
  }

  .card h2{
    font-size:18px!important;
    margin-bottom:12px!important;
  }

  .card h3{
    font-size:15px!important;
    margin-bottom:8px!important;
  }

  .grid{
    display:grid!important;
    grid-template-columns:1fr!important;
    gap:9px!important;
  }

  .field{
    margin:0!important;
  }

  .field label{
    font-size:11.5px!important;
    margin-bottom:5px!important;
  }

  input,
  select,
  textarea{
    min-height:38px!important;
    height:38px!important;
    padding:8px 10px!important;
    font-size:13px!important;
    border-radius:12px!important;
  }

  textarea{
    height:auto!important;
    min-height:72px!important;
  }

  .btn{
    min-height:38px!important;
    height:38px!important;
    padding:8px 11px!important;
    font-size:12.5px!important;
    border-radius:12px!important;
  }

  .pass-wrap{
    gap:6px!important;
  }

  .pass-wrap input{
    min-width:0!important;
  }

  .pass-eye{
    width:38px!important;
    height:38px!important;
    min-width:38px!important;
    padding:0!important;
    border-radius:12px!important;
    font-size:13px!important;
  }

  .tablewrap{
    overflow:visible!important;
    border:0!important;
    border-radius:0!important;
    background:transparent!important;
  }

  .tablewrap::after{
    display:none!important;
    content:none!important;
  }

  .tablewrap table{
    display:block!important;
    width:100%!important;
    min-width:0!important;
    border:0!important;
    background:transparent!important;
  }

  .tablewrap tbody,
  .tablewrap tr{
    display:block!important;
    width:100%!important;
  }

  .tablewrap tr:first-child{
    display:none!important;
  }

  .tablewrap tr{
    margin:0 0 10px!important;
    padding:10px!important;
    border:1px solid var(--line)!important;
    border-radius:16px!important;
    background:var(--card)!important;
    box-shadow:0 8px 22px rgba(15,23,42,.05)!important;
  }

  .tablewrap td{
    display:grid!important;
    grid-template-columns:88px minmax(0,1fr)!important;
    gap:8px!important;
    align-items:center!important;
    width:100%!important;
    border:0!important;
    padding:6px 0!important;
    font-size:12.5px!important;
    min-width:0!important;
  }

  .tablewrap td::before{
    color:var(--muted)!important;
    font-size:10px!important;
    font-weight:900!important;
    text-transform:uppercase!important;
    letter-spacing:.04em!important;
  }

  .tablewrap td:nth-child(1)::before{content:"ID";}
  .tablewrap td:nth-child(2)::before{content:"Username";}
  .tablewrap td:nth-child(3)::before{content:"Role";}
  .tablewrap td:nth-child(4)::before{content:"Router";}
  .tablewrap td:nth-child(5)::before{content:"Status";}
  .tablewrap td:nth-child(6)::before{content:"Password";}
  .tablewrap td:nth-child(7)::before{content:"Aksi";}

  .tablewrap td input,
  .tablewrap td select{
    width:100%!important;
    min-width:0!important;
    height:36px!important;
    min-height:36px!important;
    padding:7px 9px!important;
    font-size:12.5px!important;
  }

  .tablewrap td small{
    display:block!important;
    margin-top:4px!important;
    max-width:100%!important;
    font-size:10.5px!important;
    line-height:1.25!important;
  }

  .tablewrap td .pass-wrap{
    margin-bottom:6px!important;
  }

  .tablewrap td .actions{
    display:grid!important;
    grid-template-columns:1fr 1fr!important;
    gap:8px!important;
    width:100%!important;
  }

  .tablewrap td .actions .btn{
    width:100%!important;
    min-width:0!important;
  }

  .tablewrap td:last-child{
    grid-template-columns:88px minmax(0,1fr)!important;
    padding-top:10px!important;
    border-top:1px dashed var(--line)!important;
  }

  .router-form-compact-v2 .grid{
    grid-template-columns:1fr!important;
    gap:9px!important;
  }

  .router-form-compact-v2 .field:has(input[type="checkbox"]){
    height:38px!important;
    min-height:38px!important;
    padding:7px 10px!important;
  }

  .router-form-compact-v2 input[type="checkbox"]{
    width:16px!important;
    height:16px!important;
    min-width:16px!important;
    min-height:16px!important;
  }

  .wa-raw{
    max-height:180px!important;
    font-size:10.5px!important;
  }
}

@media(max-width:420px){
  .mobile-admin-title{
    font-size:17px!important;
  }

  .tablewrap td{
    grid-template-columns:76px minmax(0,1fr)!important;
  }

  .tablewrap td:last-child{
    grid-template-columns:76px minmax(0,1fr)!important;
  }

  .tablewrap td .actions{
    grid-template-columns:1fr!important;
  }
}


/* ADMIN_HEALTH_MOBILE_THEME_V1 */
.mobile-theme-btn{
  width:44px;
  height:44px;
  border:1px solid var(--line);
  border-radius:16px;
  background:var(--card);
  color:var(--text);
  font-size:18px;
  font-weight:900;
  cursor:pointer;
  box-shadow:0 10px 24px rgba(15,23,42,.08);
}

.mobile-theme-btn .moon{
  display:none;
}

html[data-admin-theme="dark"] .mobile-theme-btn .sun{
  display:none;
}

html[data-admin-theme="dark"] .mobile-theme-btn .moon{
  display:inline;
}

.health-page{
  display:flex;
  flex-direction:column;
  gap:12px;
}

.health-hero{
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:14px;
}

.health-hero h2{
  margin:0 0 4px!important;
}

.health-hero p{
  margin:0!important;
}

.health-pill{
  display:inline-flex;
  align-items:center;
  justify-content:center;
  min-width:94px;
  min-height:38px;
  padding:8px 12px;
  border-radius:999px;
  font-weight:950;
  font-size:13px;
}

.health-pill.ok{
  background:#dcfce7;
  color:#166534;
}

.health-pill.bad{
  background:#fee2e2;
  color:#991b1b;
}

.health-grid{
  gap:12px!important;
}

.health-top-grid{
  grid-template-columns:repeat(3,minmax(0,1fr))!important;
}

.health-stat-grid{
  grid-template-columns:repeat(5,minmax(0,1fr))!important;
}

.health-card .price,
.health-mini .price{
  font-size:24px!important;
}

.health-mini{
  padding:14px!important;
}

.health-mini b,
.health-card b{
  display:block;
  font-size:12px;
  color:var(--muted);
  text-transform:uppercase;
  letter-spacing:.04em;
}

.health-setting-grid{
  grid-template-columns:repeat(3,minmax(0,1fr))!important;
}

.health-raw-card{
  padding:0!important;
  overflow:hidden;
}

.health-raw-card summary{
  cursor:pointer;
  padding:16px 18px;
  font-weight:950;
  list-style:none;
}

.health-raw-card summary::-webkit-details-marker{
  display:none;
}

.health-raw-card summary::after{
  content:"Buka";
  float:right;
  color:var(--muted);
  font-size:12px;
}

.health-raw-card[open] summary::after{
  content:"Tutup";
}

.health-raw-card .wa-raw{
  margin:0!important;
  border-top:1px solid var(--line);
  border-radius:0!important;
}

html[data-admin-theme="dark"] .health-pill.ok{
  background:rgba(34,197,94,.18);
  color:#86efac;
}

html[data-admin-theme="dark"] .health-pill.bad{
  background:rgba(239,68,68,.18);
  color:#fca5a5;
}

@media(max-width:980px){
  .health-top-grid,
  .health-stat-grid,
  .health-setting-grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }
}

@media(max-width:760px){
  .mobile-theme-btn{
    display:inline-flex!important;
    align-items:center!important;
    justify-content:center!important;
    flex:0 0 auto!important;
  }

  .health-page{
    gap:9px!important;
  }

  .health-hero{
    align-items:flex-start!important;
    padding:12px!important;
  }

  .health-hero h2{
    font-size:18px!important;
  }

  .health-hero p{
    font-size:12px!important;
    line-height:1.35!important;
  }

  .health-pill{
    min-width:72px!important;
    min-height:32px!important;
    padding:6px 9px!important;
    font-size:11px!important;
  }

  .health-top-grid,
  .health-stat-grid,
  .health-setting-grid{
    grid-template-columns:1fr!important;
    gap:8px!important;
  }

  .health-card,
  .health-mini,
  .health-settings{
    padding:12px!important;
  }

  .health-card .price,
  .health-mini .price{
    font-size:22px!important;
  }

  .health-mini{
    display:flex!important;
    align-items:center!important;
    justify-content:space-between!important;
  }

  .health-mini b{
    margin:0!important;
  }

  .health-raw-card summary{
    padding:13px 14px!important;
    font-size:14px!important;
  }

  .health-raw-card .wa-raw{
    max-height:190px!important;
    font-size:10.5px!important;
  }
}


/* ADMIN_HEALTH_MOBILE_2COL_V1 */
@media(max-width:760px){
  .health-page{
    gap:7px!important;
  }

  .health-hero{
    padding:9px 10px!important;
    border-radius:14px!important;
    align-items:center!important;
  }

  .health-hero h2{
    font-size:14px!important;
    line-height:1.1!important;
    margin:0 0 3px!important;
  }

  .health-hero p{
    font-size:10.5px!important;
    line-height:1.2!important;
    max-width:190px!important;
  }

  .health-pill{
    min-width:58px!important;
    min-height:28px!important;
    padding:5px 8px!important;
    font-size:10px!important;
    white-space:nowrap!important;
  }

  .health-top-grid,
  .health-stat-grid,
  .health-setting-grid{
    display:grid!important;
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:7px!important;
  }

  .health-card,
  .health-mini,
  .health-settings{
    min-height:0!important;
    padding:9px 10px!important;
    border-radius:14px!important;
  }

  .health-card b,
  .health-mini b{
    font-size:9.5px!important;
    letter-spacing:.05em!important;
    margin:0 0 4px!important;
  }

  .health-card .price,
  .health-mini .price{
    font-size:18px!important;
    line-height:1!important;
    margin:0!important;
  }

  .health-card p,
  .health-mini p,
  .health-settings p{
    font-size:10.5px!important;
    line-height:1.25!important;
    margin-top:6px!important;
  }

  .health-mini{
    display:block!important;
  }

  .health-settings h2{
    font-size:14px!important;
    margin-bottom:8px!important;
  }

  .health-settings .field label{
    font-size:9.5px!important;
    margin-bottom:4px!important;
  }

  .health-settings input{
    height:32px!important;
    min-height:32px!important;
    font-size:11px!important;
    padding:6px 8px!important;
  }

  .health-raw-card{
    border-radius:14px!important;
  }

  .health-raw-card summary{
    padding:10px 12px!important;
    font-size:12px!important;
  }

  .health-raw-card .wa-raw{
    max-height:120px!important;
    font-size:9.5px!important;
    line-height:1.35!important;
  }
}

@media(max-width:420px){
  .health-top-grid,
  .health-stat-grid,
  .health-setting-grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }

  .health-card,
  .health-mini{
    padding:8px 9px!important;
  }

  .health-card .price,
  .health-mini .price{
    font-size:17px!important;
  }
}


/* ADMIN_TABLET_RESPONSIVE_V1 */
@media(min-width:761px) and (max-width:1024px){
  body{
    padding-left:14px!important;
    padding-right:14px!important;
  }

  .mobile-adminbar{
    display:flex!important;
  }

  .nav{
    display:none!important;
  }

  .mobile-admin-title{
    font-size:22px!important;
  }

  .wrap,
  main,
  .container{
    max-width:100%!important;
  }

  .card{
    padding:14px!important;
    border-radius:18px!important;
    margin:10px 0!important;
  }

  .card h2{
    font-size:20px!important;
    margin-bottom:12px!important;
  }

  .grid{
    display:grid!important;
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:10px!important;
  }

  .field label{
    font-size:12px!important;
    margin-bottom:6px!important;
  }

  input,
  select,
  textarea{
    min-height:40px!important;
    height:40px!important;
    padding:9px 11px!important;
    font-size:13px!important;
    border-radius:13px!important;
  }

  .btn{
    min-height:39px!important;
    height:39px!important;
    padding:9px 12px!important;
    font-size:12.5px!important;
    border-radius:13px!important;
  }

  .actions{
    gap:8px!important;
    flex-wrap:wrap!important;
  }

  .price{
    font-size:24px!important;
  }

  .tablewrap{
    overflow-x:auto!important;
    border-radius:16px!important;
  }

  table{
    min-width:780px!important;
  }

  th{
    font-size:10.5px!important;
    padding:9px 10px!important;
  }

  td{
    font-size:12px!important;
    padding:9px 10px!important;
  }

  td .actions{
    display:flex!important;
    flex-wrap:nowrap!important;
  }

  td .btn{
    min-width:72px!important;
  }

  .router-form-compact-v2 .grid,
  .package-form-compact-v1 .grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }

  .router-checks{
    display:grid!important;
    grid-template-columns:repeat(3,minmax(0,1fr))!important;
    gap:8px!important;
  }

  .router-checks label{
    min-height:40px!important;
    padding:8px 10px!important;
    border-radius:14px!important;
    font-size:12px!important;
  }

  .router-checks input[type="checkbox"]{
    width:16px!important;
    height:16px!important;
  }

  .health-top-grid,
  .health-stat-grid,
  .health-setting-grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }

  .health-mini{
    display:flex!important;
    align-items:center!important;
    justify-content:space-between!important;
  }

  .health-card,
  .health-mini{
    padding:13px!important;
  }

  .health-card .price,
  .health-mini .price{
    font-size:22px!important;
  }
}

/* laptop kecil */
@media(min-width:1025px) and (max-width:1280px){
  body{
    padding-left:18px!important;
    padding-right:18px!important;
  }

  .nav{
    gap:8px!important;
  }

  .nav a{
    padding:10px 13px!important;
    font-size:13px!important;
  }

  .card{
    padding:16px!important;
  }

  .grid{
    gap:10px!important;
  }

  table{
    font-size:12px!important;
  }

  th,
  td{
    padding:9px 10px!important;
  }
}


/* ADMIN_TX_SETTINGS_RESPONSIVE_V1 */
.tx-toolbar{
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:12px;
  margin-bottom:14px;
}

.tx-toolbar h2{
  margin:0!important;
}

.transactions-compact-table td b{
  font-weight:950!important;
}

.transactions-compact-table td .actions form{
  margin:0!important;
}

.transactions-compact-table td .actions{
  display:flex;
  gap:8px;
  flex-wrap:wrap;
}

.transactions-compact-table .badge{
  white-space:nowrap;
}

.settings-page-card > .muted{
  margin-top:-4px!important;
  margin-bottom:14px!important;
}

.settings-form-compact-v1 .grid{
  display:grid!important;
  grid-template-columns:repeat(2,minmax(0,1fr))!important;
  gap:12px!important;
}

.settings-form-compact-v1 .settings-actions{
  margin-top:14px!important;
}

@media(min-width:761px) and (max-width:1024px){
  .tx-toolbar{
    align-items:flex-start!important;
  }

  .tx-toolbar > .actions{
    justify-content:flex-end!important;
  }

  .transactions-compact-table{
    min-width:860px!important;
  }

  .transactions-compact-table th{
    font-size:10px!important;
  }

  .transactions-compact-table td{
    font-size:11.5px!important;
  }

  .transactions-compact-table td .btn{
    min-width:78px!important;
    padding-left:9px!important;
    padding-right:9px!important;
  }

  .settings-form-compact-v1 .grid{
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
  }
}

@media(max-width:760px){
  .transactions-page{
    padding:11px!important;
  }

  .tx-toolbar{
    display:grid!important;
    grid-template-columns:1fr!important;
    gap:8px!important;
    margin-bottom:10px!important;
  }

  .tx-toolbar h2{
    font-size:18px!important;
  }

  .tx-toolbar > .actions{
    display:grid!important;
    grid-template-columns:1fr 1fr!important;
    gap:8px!important;
    width:100%!important;
  }

  .tx-toolbar > .actions .btn{
    width:100%!important;
  }

  .transactions-compact-table td:nth-child(1)::before{content:"Order"!important;}
  .transactions-compact-table td:nth-child(2)::before{content:"Router/Paket"!important;}
  .transactions-compact-table td:nth-child(3)::before{content:"Pembeli"!important;}
  .transactions-compact-table td:nth-child(4)::before{content:"Total"!important;}
  .transactions-compact-table td:nth-child(5)::before{content:"Status"!important;}
  .transactions-compact-table td:nth-child(6)::before{content:"Voucher"!important;}
  .transactions-compact-table td:nth-child(7)::before{content:"Provision"!important;}
  .transactions-compact-table td:nth-child(8)::before{content:"WA"!important;}
  .transactions-compact-table td:nth-child(9)::before{content:"Aksi"!important;}

  .transactions-compact-table td{
    grid-template-columns:82px minmax(0,1fr)!important;
    padding:5px 0!important;
  }

  .transactions-compact-table td b{
    font-size:13px!important;
    word-break:break-word!important;
  }

  .transactions-compact-table td:nth-child(6) b{
    font-size:16px!important;
  }

  .transactions-compact-table td small,
  .transactions-compact-table td .muted{
    font-size:10.5px!important;
    line-height:1.25!important;
    word-break:break-word!important;
  }

  .transactions-compact-table td .badge{
    padding:4px 7px!important;
    font-size:10px!important;
  }

  .transactions-compact-table td .actions{
    display:grid!important;
    grid-template-columns:1fr 1fr!important;
    gap:8px!important;
    width:100%!important;
  }

  .transactions-compact-table td .actions form,
  .transactions-compact-table td .actions button{
    width:100%!important;
  }

  .settings-page-card{
    padding:12px!important;
  }

  .settings-page-card h2{
    font-size:18px!important;
    margin-bottom:6px!important;
  }

  .settings-page-card > .muted{
    font-size:12px!important;
    line-height:1.35!important;
    margin-bottom:10px!important;
  }

  .settings-form-compact-v1 .grid{
    grid-template-columns:1fr!important;
    gap:9px!important;
  }

  .settings-form-compact-v1 .field label{
    font-size:11px!important;
  }

  .settings-form-compact-v1 input{
    height:38px!important;
    min-height:38px!important;
    font-size:12.5px!important;
  }

  .settings-actions{
    display:grid!important;
    grid-template-columns:1fr!important;
  }

  .settings-actions .btn{
    width:100%!important;
  }
}

@media(max-width:420px){
  .transactions-compact-table td{
    grid-template-columns:74px minmax(0,1fr)!important;
  }

  .transactions-compact-table td .actions{
    grid-template-columns:1fr!important;
  }
}


/* ADMIN_WHATSAPP_RESPONSIVE_V1 */
.wa-page{
  display:flex;
  flex-direction:column;
  gap:12px;
}

.wa-hero{
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:14px;
}

.wa-hero h2{
  margin:0 0 4px!important;
}

.wa-hero p{
  margin:0!important;
}

.wa-hero-pill{
  display:inline-flex;
  align-items:center;
  justify-content:center;
  min-width:92px;
  min-height:38px;
  padding:8px 12px;
  border-radius:999px;
  font-size:12px;
  font-weight:950;
}

.wa-hero-pill.ok{
  background:#dcfce7;
  color:#166534;
}

.wa-hero-pill.bad{
  background:#fee2e2;
  color:#991b1b;
}

.wa-status-grid{
  grid-template-columns:repeat(3,minmax(0,1fr))!important;
  gap:12px!important;
}

.wa-stat{
  padding:14px!important;
}

.wa-stat b,
.wa-device-card b{
  display:block;
  font-size:12px;
  color:var(--muted);
  text-transform:uppercase;
  letter-spacing:.04em;
  margin-bottom:8px;
}

.wa-stat span{
  font-size:20px;
  font-weight:950;
}

.wa-device-card small{
  display:block;
  font-size:12px;
  line-height:1.45;
  word-break:break-word;
}

.wa-qr-box{
  display:grid;
  grid-template-columns:minmax(0,1fr) auto;
  align-items:center;
  gap:18px;
}

.wa-qr-box h2{
  margin-bottom:6px!important;
}

.wa-qr-img{
  width:220px!important;
  max-width:100%!important;
  border-radius:18px;
  background:#fff;
  padding:10px;
  border:1px solid var(--line);
}

.wa-action-card h2,
.wa-test-card h2{
  margin-bottom:12px!important;
}

.wa-actions{
  display:grid!important;
  grid-template-columns:1fr 1fr!important;
  gap:10px!important;
}

.wa-actions form{
  margin:0!important;
}

.wa-actions button{
  width:100%!important;
}

.wa-test-form{
  grid-template-columns:1fr 2fr auto!important;
  align-items:end!important;
}

.wa-test-form .btn{
  width:100%!important;
}

.wa-raw-card{
  padding:0!important;
  overflow:hidden;
}

.wa-raw-card summary{
  cursor:pointer;
  padding:16px 18px;
  font-weight:950;
  list-style:none;
}

.wa-raw-card summary::-webkit-details-marker{
  display:none;
}

.wa-raw-card summary::after{
  content:"Buka";
  float:right;
  color:var(--muted);
  font-size:12px;
}

.wa-raw-card[open] summary::after{
  content:"Tutup";
}

.wa-raw-card .wa-raw{
  margin:0!important;
  border-top:1px solid var(--line);
  border-radius:0!important;
  max-height:240px!important;
}

html[data-admin-theme="dark"] .wa-hero-pill.ok{
  background:rgba(34,197,94,.18);
  color:#86efac;
}

html[data-admin-theme="dark"] .wa-hero-pill.bad{
  background:rgba(239,68,68,.18);
  color:#fca5a5;
}

@media(min-width:761px) and (max-width:1024px){
  .wa-status-grid{
    grid-template-columns:repeat(3,minmax(0,1fr))!important;
  }

  .wa-stat span{
    font-size:18px!important;
  }

  .wa-qr-box{
    grid-template-columns:1fr auto!important;
  }

  .wa-qr-img{
    width:190px!important;
  }

  .wa-test-form{
    grid-template-columns:1fr 1fr!important;
  }

  .wa-test-form .field:last-child{
    grid-column:span 2!important;
  }
}

@media(max-width:760px){
  .wa-page{
    gap:8px!important;
  }

  .wa-hero{
    padding:11px 12px!important;
    border-radius:15px!important;
  }

  .wa-hero h2{
    font-size:18px!important;
  }

  .wa-hero p{
    font-size:12px!important;
    line-height:1.3!important;
  }

  .wa-hero-pill{
    min-width:72px!important;
    min-height:32px!important;
    padding:6px 9px!important;
    font-size:10.5px!important;
  }

  .wa-status-grid{
    grid-template-columns:repeat(3,minmax(0,1fr))!important;
    gap:7px!important;
  }

  .wa-stat{
    padding:9px!important;
    border-radius:14px!important;
  }

  .wa-stat b,
  .wa-device-card b{
    font-size:9.5px!important;
    margin-bottom:5px!important;
  }

  .wa-stat span{
    font-size:13px!important;
    line-height:1.1!important;
    word-break:break-word!important;
  }

  .wa-stat .badge{
    font-size:9.5px!important;
    padding:4px 6px!important;
  }

  .wa-device-card{
    padding:10px 12px!important;
  }

  .wa-device-card small{
    font-size:10.5px!important;
    line-height:1.3!important;
  }

  .wa-qr-box{
    grid-template-columns:1fr!important;
    text-align:center!important;
    padding:12px!important;
  }

  .wa-qr-box h2{
    font-size:16px!important;
  }

  .wa-qr-box p{
    font-size:11px!important;
  }

  .wa-qr-img{
    width:210px!important;
    margin:auto!important;
    border-radius:16px!important;
  }

  .wa-action-card,
  .wa-test-card{
    padding:12px!important;
  }

  .wa-action-card h2,
  .wa-test-card h2{
    font-size:16px!important;
    margin-bottom:10px!important;
  }

  .wa-actions{
    grid-template-columns:1fr!important;
    gap:8px!important;
  }

  .wa-test-form{
    grid-template-columns:1fr!important;
    gap:8px!important;
  }

  .wa-test-form .field label{
    font-size:11px!important;
  }

  .wa-test-form input{
    height:38px!important;
    min-height:38px!important;
    font-size:12.5px!important;
  }

  .wa-raw-card summary{
    padding:12px 14px!important;
    font-size:13px!important;
  }

  .wa-raw-card .wa-raw{
    max-height:150px!important;
    font-size:10px!important;
  }
}

@media(max-width:420px){
  .wa-status-grid{
    grid-template-columns:1fr 1fr!important;
  }

  .wa-stat:nth-child(3){
    grid-column:span 2!important;
  }
}

/* ROUTER_PAGE_CSS_EXTRACTED_V1 */
/* ROUTER_TEST_RESULT_V1 */
.router-test-result-v1{
        border-left:6px solid #94a3b8;
}
.router-test-result-v1.ok{
        border-left-color:#16a34a;
        background:rgba(240,253,244,.88);
}
.router-test-result-v1.bad{
        border-left-color:#dc2626;
        background:rgba(254,242,242,.88);
}
.router-test-result-v1 b{
        display:block;
        margin-bottom:6px;
        color:#0f172a;
        font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace;
        font-size:15px;
        font-weight:700;
        letter-spacing:.2px;
}
.router-test-result-v1 div{
        color:#334155;
        line-height:1.6;
        font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace;
        font-size:14px;
        white-space:pre-wrap;
}

/* ROUTER_API_FORM_FINAL_V4 */
.router-api-form-v2 .grid{
        display:grid !important;
        grid-template-columns:repeat(4,minmax(0,1fr)) !important;
        gap:16px !important;
        align-items:start !important;
}
.router-api-form-v2 .field{
        display:flex !important;
        flex-direction:column !important;
        gap:8px !important;
        min-width:0 !important;
}
.router-api-form-v2 .field label{
        display:flex !important;
        align-items:center !important;
        height:22px !important;
        min-height:22px !important;
        max-height:22px !important;
        margin:0 !important;
        padding:0 !important;
        line-height:22px !important;
}
.router-api-form-v2 input,
.router-api-form-v2 select.router-api-select-v3{
        display:block !important;
        width:100% !important;
        height:50px !important;
        min-height:50px !important;
        max-height:50px !important;
        margin:0 !important;
        padding-top:0 !important;
        padding-bottom:0 !important;
        padding-left:16px !important;
        box-sizing:border-box !important;
        border-radius:16px !important;
        border:1px solid rgba(148,163,184,.35) !important;
        background-color:#fff !important;
        color:#0f172a !important;
        font-size:16px !important;
        line-height:50px !important;
        vertical-align:top !important;
}
.router-api-form-v2 select.router-api-select-v3{
        padding-right:42px !important;
        -webkit-appearance:none !important;
        -moz-appearance:none !important;
        appearance:none !important;
        background-image:linear-gradient(45deg, transparent 50%, #334155 50%),linear-gradient(135deg, #334155 50%, transparent 50%) !important;
        background-position:calc(100% - 20px) 22px,calc(100% - 14px) 22px !important;
        background-size:6px 6px,6px 6px !important;
        background-repeat:no-repeat !important;
}
.router-api-form-v2 .router-checks{
        display:grid !important;
        grid-template-columns:repeat(3,minmax(0,1fr)) !important;
        gap:12px !important;
        margin-top:16px !important;
}
.router-api-guide-v4{
        margin-top:14px;
        padding:14px 16px;
        border:1px solid rgba(148,163,184,.35);
        background:rgba(248,250,252,.9);
        border-radius:16px;
        color:#334155;
        font-size:14px;
        line-height:1.5;
}
.router-api-guide-v4 b{color:#0f172a}
.router-api-guide-v4 code{
        padding:2px 6px;
        border-radius:8px;
        background:#eef2ff;
        color:#1e293b;
}
@media (max-width:1100px){
        .router-api-form-v2 .grid{
                grid-template-columns:repeat(2,minmax(0,1fr)) !important;
        }
}
@media (max-width:640px){
        .router-api-form-v2 .grid,
        .router-api-form-v2 .router-checks{
                grid-template-columns:1fr !important;
        }
}

/* ROUTER_API_EQUAL_HEIGHT_FIX_V5 */
.router-api-form-v2 .field > input,
.router-api-form-v2 .field > select.router-api-select-v3{
        height:50px !important;
        min-height:50px !important;
        max-height:50px !important;
        padding-top:0 !important;
        padding-bottom:0 !important;
        margin-top:0 !important;
        margin-bottom:0 !important;
        box-sizing:border-box !important;
}
.router-api-form-v2 .field > select.router-api-select-v3{
        line-height:normal !important;
        padding-left:16px !important;
        padding-right:42px !important;
}



/* ROUTER_API_DARK_POLISH_V1 */
html.dark .router-api-form-v2 input,
html.dark .router-api-form-v2 select.router-api-select-v3,
body.dark .router-api-form-v2 input,
body.dark .router-api-form-v2 select.router-api-select-v3,
.dark .router-api-form-v2 input,
.dark .router-api-form-v2 select.router-api-select-v3,
[data-theme="dark"] .router-api-form-v2 input,
[data-theme="dark"] .router-api-form-v2 select.router-api-select-v3{
        background-color:#111827 !important;
        color:#e5e7eb !important;
        border-color:rgba(148,163,184,.22) !important;
        box-shadow:inset 0 1px 0 rgba(255,255,255,.025) !important;
}

html.dark .router-api-form-v2 input::placeholder,
body.dark .router-api-form-v2 input::placeholder,
.dark .router-api-form-v2 input::placeholder,
[data-theme="dark"] .router-api-form-v2 input::placeholder{
        color:#64748b !important;
}

html.dark .router-api-form-v2 select.router-api-select-v3,
body.dark .router-api-form-v2 select.router-api-select-v3,
.dark .router-api-form-v2 select.router-api-select-v3,
[data-theme="dark"] .router-api-form-v2 select.router-api-select-v3{
        background-image:linear-gradient(45deg, transparent 50%, #cbd5e1 50%),linear-gradient(135deg, #cbd5e1 50%, transparent 50%) !important;
}

html.dark .router-api-form-v2 .router-checks label,
body.dark .router-api-form-v2 .router-checks label,
.dark .router-api-form-v2 .router-checks label,
[data-theme="dark"] .router-api-form-v2 .router-checks label{
        background:#111827 !important;
        color:#e5e7eb !important;
        border:1px solid rgba(148,163,184,.22) !important;
        box-shadow:inset 0 1px 0 rgba(255,255,255,.025) !important;
}

html.dark .router-api-guide-v4,
body.dark .router-api-guide-v4,
.dark .router-api-guide-v4,
[data-theme="dark"] .router-api-guide-v4{
        background:#0f172a !important;
        border-color:rgba(148,163,184,.22) !important;
        color:#cbd5e1 !important;
}

html.dark .router-api-guide-v4 b,
body.dark .router-api-guide-v4 b,
.dark .router-api-guide-v4 b,
[data-theme="dark"] .router-api-guide-v4 b{
        color:#f8fafc !important;
}

html.dark .router-api-guide-v4 code,
body.dark .router-api-guide-v4 code,
.dark .router-api-guide-v4 code,
[data-theme="dark"] .router-api-guide-v4 code{
        background:#1e293b !important;
        color:#e2e8f0 !important;
}

html.dark .router-test-result-v1 b,
body.dark .router-test-result-v1 b,
.dark .router-test-result-v1 b,
[data-theme="dark"] .router-test-result-v1 b{
        color:#f8fafc !important;
}

html.dark .router-test-result-v1 div,
body.dark .router-test-result-v1 div,
.dark .router-test-result-v1 div,
[data-theme="dark"] .router-test-result-v1 div{
        color:#cbd5e1 !important;
}

/* ROUTER_TEST_RESULT_COLOR_V1 */
.router-test-result-v1.ok b{
        color:#166534 !important;
}
.router-test-result-v1.ok div{
        color:#15803d !important;
}
.router-test-result-v1.bad b{
        color:#991b1b !important;
}
.router-test-result-v1.bad div{
        color:#b91c1c !important;
}
html.dark .router-test-result-v1.ok b,
body.dark .router-test-result-v1.ok b,
.dark .router-test-result-v1.ok b,
[data-theme="dark"] .router-test-result-v1.ok b{
        color:#86efac !important;
}
html.dark .router-test-result-v1.ok div,
body.dark .router-test-result-v1.ok div,
.dark .router-test-result-v1.ok div,
[data-theme="dark"] .router-test-result-v1.ok div{
        color:#bbf7d0 !important;
}
html.dark .router-test-result-v1.bad b,
body.dark .router-test-result-v1.bad b,
.dark .router-test-result-v1.bad b,
[data-theme="dark"] .router-test-result-v1.bad b{
        color:#fca5a5 !important;
}
html.dark .router-test-result-v1.bad div,
body.dark .router-test-result-v1.bad div,
.dark .router-test-result-v1.bad div,
[data-theme="dark"] .router-test-result-v1.bad div{
        color:#fecaca !important;
}
/* /ROUTER_PAGE_CSS_EXTRACTED_V1 */


/* NAVBAR_PALETTE_MATCH_V1 */
.nav{
  background:
    linear-gradient(135deg, rgba(8,15,34,.94), rgba(36,18,47,.92)) !important;
  border:1px solid rgba(138,92,246,.16) !important;
  border-radius:24px !important;
  padding:14px 18px !important;
  box-shadow:
    0 16px 42px rgba(0,0,0,.28),
    inset 0 1px 0 rgba(255,255,255,.04) !important;
  backdrop-filter: blur(14px) !important;
}

.nav .brand{
  color:#ffffff !important;
  font-weight:950 !important;
  letter-spacing:-.02em !important;
}

.nav .menu{
  gap:10px !important;
}

.nav .menu a,
.admin-theme-toggle{
  background:linear-gradient(180deg, rgba(255,255,255,.07), rgba(255,255,255,.03)) !important;
  border:1px solid rgba(255,255,255,.09) !important;
  color:#eef2ff !important;
  border-radius:16px !important;
  box-shadow:
    inset 0 1px 0 rgba(255,255,255,.04),
    0 8px 20px rgba(0,0,0,.18) !important;
  transition:all .18s ease !important;
}

.nav .menu a:hover,
.admin-theme-toggle:hover{
  transform:translateY(-1px) !important;
  border-color:rgba(167,139,250,.30) !important;
  background:linear-gradient(180deg, rgba(139,92,246,.18), rgba(56,189,248,.10)) !important;
  color:#ffffff !important;
  box-shadow:
    0 12px 26px rgba(0,0,0,.24),
    0 0 0 1px rgba(167,139,250,.08) inset !important;
}

.nav .menu a[href="/admin"],
.nav .menu a.active{
  background:linear-gradient(180deg, rgba(59,130,246,.22), rgba(139,92,246,.16)) !important;
  border-color:rgba(96,165,250,.24) !important;
  color:#ffffff !important;
}

html[data-admin-theme="light"] .nav,
html[data-theme="light"] .nav{
  background:
    linear-gradient(135deg, rgba(255,255,255,.96), rgba(241,245,249,.96)) !important;
  border-color:#dbe4f0 !important;
  box-shadow:
    0 14px 34px rgba(15,23,42,.08),
    inset 0 1px 0 rgba(255,255,255,.80) !important;
}

html[data-admin-theme="light"] .nav .brand,
html[data-theme="light"] .nav .brand{
  color:#0f172a !important;
}

html[data-admin-theme="light"] .nav .menu a,
html[data-theme="light"] .nav .menu a,
html[data-admin-theme="light"] .admin-theme-toggle,
html[data-theme="light"] .admin-theme-toggle{
  background:linear-gradient(180deg,#ffffff,#f8fafc) !important;
  border-color:#dbe4f0 !important;
  color:#0f172a !important;
  box-shadow:0 6px 18px rgba(15,23,42,.06) !important;
}

html[data-admin-theme="light"] .nav .menu a:hover,
html[data-theme="light"] .nav .menu a:hover,
html[data-admin-theme="light"] .admin-theme-toggle:hover,
html[data-theme="light"] .admin-theme-toggle:hover{
  background:linear-gradient(180deg,#eff6ff,#eef2ff) !important;
  border-color:#c7d2fe !important;
  color:#1e293b !important;
}

@media(max-width:760px){
  .nav{
    border-radius:20px !important;
    padding:10px 12px !important;
  }

  .nav .menu{
    gap:8px !important;
  }

  .nav .menu a,
  .admin-theme-toggle{
    border-radius:14px !important;
  }
}
/* /NAVBAR_PALETTE_MATCH_V1 */


/* NAVBAR_PALETTE_MATCH_FIX_V2 */
/* Samakan navbar dengan palette card dashboard, dan hapus efek active permanen Dashboard */
.nav{
  background:
    radial-gradient(circle at 12% 0%, rgba(139,92,246,.16), transparent 34%),
    radial-gradient(circle at 88% 100%, rgba(244,63,94,.14), transparent 34%),
    linear-gradient(135deg, rgba(36,22,49,.96), rgba(62,25,54,.92) 58%, rgba(15,25,48,.94)) !important;
  border:1px solid rgba(168,85,247,.18) !important;
  box-shadow:
    0 16px 42px rgba(0,0,0,.24),
    inset 0 1px 0 rgba(255,255,255,.055) !important;
}

/* reset tombol Dashboard yang sebelumnya dibuat active permanen */
.nav .menu a[href="/admin"]{
  background:linear-gradient(180deg, rgba(255,255,255,.07), rgba(255,255,255,.035)) !important;
  border:1px solid rgba(255,255,255,.10) !important;
  color:#eef2ff !important;
  box-shadow:
    inset 0 1px 0 rgba(255,255,255,.04),
    0 8px 20px rgba(0,0,0,.16) !important;
}

/* hover hanya saat mouse diarahkan */
.nav .menu a:hover,
.nav .menu a[href="/admin"]:hover,
.admin-theme-toggle:hover{
  transform:translateY(-1px) !important;
  background:
    radial-gradient(circle at 18% 0%, rgba(56,189,248,.20), transparent 42%),
    linear-gradient(180deg, rgba(139,92,246,.22), rgba(244,63,94,.11)) !important;
  border-color:rgba(196,181,253,.34) !important;
  color:#ffffff !important;
  box-shadow:
    0 12px 26px rgba(0,0,0,.25),
    inset 0 1px 0 rgba(255,255,255,.08) !important;
}

.nav .brand{
  color:#ffffff !important;
  text-shadow:0 8px 22px rgba(168,85,247,.18) !important;
}

.admin-theme-toggle{
  background:linear-gradient(180deg, rgba(255,255,255,.07), rgba(255,255,255,.035)) !important;
  border-color:rgba(255,255,255,.10) !important;
}

/* light mode tetap bersih */
html[data-admin-theme="light"] .nav,
html[data-theme="light"] .nav{
  background:
    radial-gradient(circle at 12% 0%, rgba(59,130,246,.10), transparent 34%),
    linear-gradient(135deg, rgba(255,255,255,.98), rgba(241,245,249,.96)) !important;
  border-color:#dbe4f0 !important;
}

html[data-admin-theme="light"] .nav .menu a[href="/admin"],
html[data-theme="light"] .nav .menu a[href="/admin"]{
  background:linear-gradient(180deg,#ffffff,#f8fafc) !important;
  border-color:#dbe4f0 !important;
  color:#0f172a !important;
}

html[data-admin-theme="light"] .nav .menu a[href="/admin"]:hover,
html[data-theme="light"] .nav .menu a[href="/admin"]:hover{
  background:linear-gradient(180deg,#eff6ff,#eef2ff) !important;
  border-color:#c7d2fe !important;
}
/* /NAVBAR_PALETTE_MATCH_FIX_V2 */

/* ADMIN_DESIGN_SYSTEM_SMOOTH_V1 */
:root{
  --ui-radius:14px;
  --ui-radius-card:18px;
  --ui-focus:0 0 0 4px rgba(37,99,235,.14);
  --ui-soft-shadow:0 12px 30px rgba(15,23,42,.08);
  --ui-hover-shadow:0 16px 36px rgba(15,23,42,.12);
}

.btn,
button.btn,
a.btn,
input[type="submit"].btn{
  position:relative!important;
  border:1px solid rgba(15,23,42,.08)!important;
  border-radius:var(--ui-radius)!important;
  font-weight:850!important;
  letter-spacing:0!important;
  box-shadow:0 7px 18px rgba(15,23,42,.08)!important;
  transition:transform .16s ease, box-shadow .16s ease, border-color .16s ease, background .16s ease, color .16s ease!important;
  -webkit-tap-highlight-color:transparent!important;
}

.btn:hover,
button.btn:hover,
a.btn:hover{
  transform:translateY(-1px)!important;
  box-shadow:var(--ui-hover-shadow)!important;
}

.btn:active,
button.btn:active,
a.btn:active{
  transform:translateY(0)!important;
  box-shadow:0 6px 14px rgba(15,23,42,.10)!important;
}

.btn:focus-visible,
button:focus-visible,
a.btn:focus-visible,
input:focus-visible,
select:focus-visible,
textarea:focus-visible{
  outline:none!important;
  border-color:rgba(37,99,235,.70)!important;
  box-shadow:var(--ui-focus)!important;
}

.btn.pri{
  background:linear-gradient(135deg,#2563eb,#1d4ed8)!important;
  color:#fff!important;
}

.btn.ok{
  background:linear-gradient(135deg,#16a34a,#15803d)!important;
  color:#fff!important;
}

.btn.danger{
  background:linear-gradient(135deg,#dc2626,#b91c1c)!important;
  color:#fff!important;
}

.btn.warn{
  background:linear-gradient(135deg,#f59e0b,#d97706)!important;
  color:#111827!important;
}

.btn.light{
  background:linear-gradient(180deg,#ffffff,#f8fafc)!important;
  color:#0f172a!important;
  border-color:#dbe4f0!important;
}

button:disabled,
.btn:disabled,
.btn[aria-disabled="true"],
.btn.loading,
button.loading{
  cursor:not-allowed!important;
  opacity:.68!important;
  transform:none!important;
  box-shadow:none!important;
}

.btn.loading::after,
button.loading::after,
.loading-state::before{
  content:""!important;
  width:14px!important;
  height:14px!important;
  border-radius:999px!important;
  border:2px solid currentColor!important;
  border-right-color:transparent!important;
  display:inline-block!important;
  animation:vgSpin .75s linear infinite!important;
}

@keyframes vgSpin{
  to{transform:rotate(360deg)}
}

input,
select,
textarea{
  border:1px solid #dbe4f0!important;
  border-radius:var(--ui-radius)!important;
  background:#fff!important;
  color:#0f172a!important;
  box-shadow:0 1px 0 rgba(15,23,42,.02)!important;
  transition:border-color .16s ease, box-shadow .16s ease, background .16s ease!important;
}

textarea{
  min-height:92px!important;
  resize:vertical!important;
  line-height:1.5!important;
}

input::placeholder,
textarea::placeholder{
  color:#94a3b8!important;
}

.field{
  gap:7px!important;
}

.field label{
  color:#475569!important;
  font-size:12px!important;
  font-weight:850!important;
}

.card{
  border-radius:var(--ui-radius-card)!important;
  border:1px solid rgba(219,228,240,.95)!important;
  background:rgba(255,255,255,.92)!important;
  box-shadow:var(--ui-soft-shadow)!important;
}

.card .card{
  box-shadow:0 8px 22px rgba(15,23,42,.06)!important;
}

.tablewrap{
  border-radius:var(--ui-radius-card)!important;
  border-color:#dbe4f0!important;
  background:#fff!important;
  box-shadow:0 10px 26px rgba(15,23,42,.06)!important;
}

table{
  background:#fff!important;
}

th{
  background:linear-gradient(180deg,#f8fafc,#f1f5f9)!important;
  color:#64748b!important;
  font-weight:900!important;
  letter-spacing:.02em!important;
}

td{
  color:#1e293b!important;
}

tbody tr{
  transition:background .14s ease!important;
}

tbody tr:hover td{
  background:#f8fbff!important;
}

.notice,
.alert,
.toast{
  border-radius:var(--ui-radius)!important;
  border:1px solid #dbe4f0!important;
  background:linear-gradient(180deg,#ffffff,#f8fafc)!important;
  color:#334155!important;
  box-shadow:0 8px 22px rgba(15,23,42,.06)!important;
}

.notice.success,
.alert.success,
.toast.success{
  background:#ecfdf5!important;
  border-color:#bbf7d0!important;
  color:#166534!important;
}

.notice.error,
.alert.error,
.toast.error{
  background:#fef2f2!important;
  border-color:#fecaca!important;
  color:#991b1b!important;
}

.notice.warn,
.alert.warn,
.toast.warn{
  background:#fffbeb!important;
  border-color:#fde68a!important;
  color:#92400e!important;
}

html[data-admin-theme="dark"]{
  --ui-focus:0 0 0 4px rgba(246,118,47,.18);
  --ui-soft-shadow:0 16px 38px rgba(4,20,43,.34);
  --ui-hover-shadow:0 18px 44px rgba(4,20,43,.44);
}

html[data-admin-theme="dark"] .btn.light{
  background:linear-gradient(180deg,rgba(255,255,255,.08),rgba(255,255,255,.04))!important;
  color:#f8fafc!important;
  border-color:rgba(255,255,255,.12)!important;
}

html[data-admin-theme="dark"] input,
html[data-admin-theme="dark"] select,
html[data-admin-theme="dark"] textarea{
  background:rgba(4,20,43,.72)!important;
  border-color:rgba(148,163,184,.20)!important;
  color:#f8fafc!important;
  box-shadow:inset 0 1px 0 rgba(255,255,255,.04)!important;
}

html[data-admin-theme="dark"] input::placeholder,
html[data-admin-theme="dark"] textarea::placeholder{
  color:#94a3b8!important;
}

html[data-admin-theme="dark"] .field label{
  color:#cbd5e1!important;
}

html[data-admin-theme="dark"] .card{
  background:linear-gradient(145deg, rgba(35,24,50,.96), rgba(81,36,56,.68))!important;
  border-color:rgba(148,163,184,.16)!important;
}

html[data-admin-theme="dark"] .tablewrap,
html[data-admin-theme="dark"] table{
  background:rgba(35,24,50,.86)!important;
  border-color:rgba(148,163,184,.16)!important;
}

html[data-admin-theme="dark"] th{
  background:rgba(81,36,56,.72)!important;
  color:#cbd5e1!important;
}

html[data-admin-theme="dark"] td{
  color:#f8fafc!important;
}

html[data-admin-theme="dark"] tbody tr:hover td{
  background:rgba(246,118,47,.08)!important;
}

html[data-admin-theme="dark"] .notice,
html[data-admin-theme="dark"] .alert,
html[data-admin-theme="dark"] .toast{
  background:rgba(4,20,43,.66)!important;
  border-color:rgba(148,163,184,.18)!important;
  color:#e5e7eb!important;
}

@media(max-width:760px){
  :root{
    --ui-radius:12px;
    --ui-radius-card:16px;
  }

  .btn,
  button.btn,
  a.btn{
    box-shadow:0 5px 14px rgba(15,23,42,.08)!important;
  }

  textarea{
    min-height:84px!important;
  }
}
/* /ADMIN_DESIGN_SYSTEM_SMOOTH_V1 */

/* ADMIN_DASHBOARD_ENTRY_MOTION_V1 */
@keyframes vgDashFadeUp{
  from{
    opacity:0;
    transform:translateY(14px) scale(.985);
    filter:blur(2px);
  }
  to{
    opacity:1;
    transform:translateY(0) scale(1);
    filter:blur(0);
  }
}

@keyframes vgDashFadeDown{
  from{
    opacity:0;
    transform:translateY(-10px);
    filter:blur(2px);
  }
  to{
    opacity:1;
    transform:translateY(0);
    filter:blur(0);
  }
}

@keyframes vgDashGlowIn{
  from{
    opacity:0;
    transform:scale(.78);
  }
  to{
    opacity:1;
    transform:scale(1);
  }
}

html.is-admin-dashboard .nav,
html.is-admin-dashboard .mobile-adminbar{
  animation:vgDashFadeDown .34s cubic-bezier(.2,.8,.2,1) both!important;
}

html.is-admin-dashboard .dashboard-top-card-v1,
html.is-admin-dashboard .dashboard-chart-card-v1,
html.is-admin-dashboard .dashboard-revenue-card-v1,
html.is-admin-dashboard .dashboard-activity-card-v1{
  opacity:0;
  animation:vgDashFadeUp .42s cubic-bezier(.2,.8,.2,1) both!important;
  will-change:transform, opacity;
}

html.is-admin-dashboard .dashboard-top-card-v1:nth-child(1){animation-delay:.04s!important}
html.is-admin-dashboard .dashboard-top-card-v1:nth-child(2){animation-delay:.08s!important}
html.is-admin-dashboard .dashboard-top-card-v1:nth-child(3){animation-delay:.12s!important}
html.is-admin-dashboard .dashboard-top-card-v1:nth-child(4){animation-delay:.16s!important}
html.is-admin-dashboard .dashboard-chart-card-v1{animation-delay:.20s!important}
html.is-admin-dashboard .dashboard-revenue-card-v1{animation-delay:.24s!important}
html.is-admin-dashboard .dashboard-activity-card-v1{animation-delay:.28s!important}

html.is-admin-dashboard .dashboard-top-card-v1:before{
  animation:vgDashGlowIn .50s cubic-bezier(.2,.8,.2,1) both!important;
  animation-delay:.20s!important;
}

html.is-admin-dashboard .dashboard-donut-v1{
  animation:vgDashGlowIn .46s cubic-bezier(.2,.8,.2,1) both!important;
  animation-delay:.30s!important;
}

html.is-admin-dashboard .dash-income-item-v1,
html.is-admin-dashboard .dashboard-legend-v1 span,
html.is-admin-dashboard .dashboard-activity-item-v1{
  opacity:0;
  animation:vgDashFadeUp .34s cubic-bezier(.2,.8,.2,1) both!important;
}

html.is-admin-dashboard .dashboard-legend-v1 span:nth-child(1),
html.is-admin-dashboard .dash-income-item-v1:nth-child(1),
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(1){animation-delay:.34s!important}
html.is-admin-dashboard .dashboard-legend-v1 span:nth-child(2),
html.is-admin-dashboard .dash-income-item-v1:nth-child(2),
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(2){animation-delay:.38s!important}
html.is-admin-dashboard .dashboard-legend-v1 span:nth-child(3),
html.is-admin-dashboard .dash-income-item-v1:nth-child(3),
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(3){animation-delay:.42s!important}
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(4){animation-delay:.46s!important}
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(5){animation-delay:.50s!important}
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(6){animation-delay:.54s!important}
html.is-admin-dashboard .dashboard-activity-item-v1:nth-child(7){animation-delay:.58s!important}

html.is-admin-dashboard .dashboard-top-card-v1,
html.is-admin-dashboard .dashboard-activity-item-v1,
html.is-admin-dashboard .dash-income-item-v1{
  transition:transform .18s ease, box-shadow .18s ease, border-color .18s ease, background .18s ease!important;
}

html.is-admin-dashboard .dashboard-top-card-v1:hover{
  transform:translateY(-3px)!important;
}

html.is-admin-dashboard .dashboard-activity-item-v1:hover,
html.is-admin-dashboard .dash-income-item-v1:hover{
  transform:translateY(-2px)!important;
}

@media (prefers-reduced-motion: reduce){
  html.is-admin-dashboard .nav,
  html.is-admin-dashboard .mobile-adminbar,
  html.is-admin-dashboard .dashboard-top-card-v1,
  html.is-admin-dashboard .dashboard-chart-card-v1,
  html.is-admin-dashboard .dashboard-revenue-card-v1,
  html.is-admin-dashboard .dashboard-activity-card-v1,
  html.is-admin-dashboard .dashboard-top-card-v1:before,
  html.is-admin-dashboard .dashboard-donut-v1,
  html.is-admin-dashboard .dash-income-item-v1,
  html.is-admin-dashboard .dashboard-legend-v1 span,
  html.is-admin-dashboard .dashboard-activity-item-v1{
    opacity:1!important;
    animation:none!important;
    transform:none!important;
    filter:none!important;
  }
}
/* /ADMIN_DASHBOARD_ENTRY_MOTION_V1 */

/* ADMIN_MOBILE_DRAWER_MOTION_V1 */
@media(max-width:760px), (min-width:768px) and (max-width:1199px){
  .mobile-drawer-overlay{
    display:block!important;
    opacity:0;
    pointer-events:none;
    backdrop-filter:blur(0);
    -webkit-backdrop-filter:blur(0);
    transition:
      opacity .30s cubic-bezier(.2,.8,.2,1),
      backdrop-filter .30s cubic-bezier(.2,.8,.2,1)!important;
    will-change:opacity, backdrop-filter;
  }

  .mobile-drawer{
    left:0!important;
    display:flex!important;
    transform:translate3d(-104%,0,0);
    opacity:.96;
    pointer-events:none;
    transition:
      transform .46s cubic-bezier(.16,1,.3,1),
      opacity .34s cubic-bezier(.2,.8,.2,1),
      box-shadow .46s cubic-bezier(.16,1,.3,1)!important;
    will-change:transform, opacity;
    contain:layout paint;
    backface-visibility:hidden;
  }

  body.admin-drawer-open .mobile-drawer-overlay{
    opacity:1;
    pointer-events:auto;
    backdrop-filter:blur(10px);
    -webkit-backdrop-filter:blur(10px);
  }

  body.admin-drawer-open .mobile-drawer{
    transform:translate3d(0,0,0);
    opacity:1;
    pointer-events:auto;
    box-shadow:30px 0 70px rgba(2,8,23,.34);
  }

  .mobile-drawer-head,
  .mobile-drawer-links a{
    opacity:0;
    transform:translate3d(-10px,0,0);
    transition:
      transform .36s cubic-bezier(.16,1,.3,1),
      opacity .32s cubic-bezier(.2,.8,.2,1),
      background .16s ease,
      border-color .16s ease!important;
    will-change:transform, opacity;
  }

  body.admin-drawer-open .mobile-drawer-head,
  body.admin-drawer-open .mobile-drawer-links a{
    opacity:1;
    transform:translate3d(0,0,0);
  }

  body.admin-drawer-open .mobile-drawer-head{transition-delay:.07s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(1){transition-delay:.10s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(2){transition-delay:.13s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(3){transition-delay:.16s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(4){transition-delay:.19s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(5){transition-delay:.22s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(6){transition-delay:.25s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(7){transition-delay:.28s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(8){transition-delay:.31s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(9){transition-delay:.34s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(10){transition-delay:.37s!important}
  body.admin-drawer-open .mobile-drawer-links a:nth-child(n+11){transition-delay:.40s!important}

  .mobile-drawer-links a:active{
    transform:translate3d(1px,0,0) scale(.99)!important;
  }

  .mobile-menu-btn,
  .mobile-drawer-close{
    transition:transform .16s ease, box-shadow .16s ease, background .16s ease!important;
    will-change:transform;
  }

  .mobile-menu-btn:active,
  .mobile-drawer-close:active{
    transform:scale(.96)!important;
  }
}

@media(max-width:760px) and (prefers-reduced-motion: reduce), (min-width:768px) and (max-width:1199px) and (prefers-reduced-motion: reduce){
  .mobile-drawer-overlay,
  .mobile-drawer,
  .mobile-drawer-head,
  .mobile-drawer-links a,
  .mobile-menu-btn,
  .mobile-drawer-close{
    transition:none!important;
    animation:none!important;
  }
}
/* /ADMIN_MOBILE_DRAWER_MOTION_V1 */

/* ADMIN_HAMBURGER_X_MOTION_V1 */
@media(max-width:760px), (min-width:768px) and (max-width:1199px){
  .mobile-menu-btn{
    position:relative!important;
    overflow:hidden!important;
    display:inline-flex!important;
    align-items:center!important;
    justify-content:center!important;
  }

  .mobile-menu-btn-area{
    position:relative!important;
    display:block!important;
    width:26px!important;
    height:22px!important;
    transition:transform .46s cubic-bezier(.16,1,.3,1)!important;
    transform-style:preserve-3d!important;
    will-change:transform!important;
  }

  .mobile-menu-btn-area span{
    position:absolute!important;
    left:4px!important;
    display:block!important;
    width:18px!important;
    height:3px!important;
    border-radius:999px!important;
    background:currentColor!important;
    transition:
      top .38s cubic-bezier(.16,1,.3,1),
      left .38s cubic-bezier(.16,1,.3,1),
      width .38s cubic-bezier(.16,1,.3,1),
      transform .42s cubic-bezier(.16,1,.3,1),
      opacity .20s ease!important;
    will-change:transform, opacity, top, left, width!important;
  }

  .mobile-menu-btn-area span:nth-child(1){top:3px!important}
  .mobile-menu-btn-area span:nth-child(2){top:10px!important}
  .mobile-menu-btn-area span:nth-child(3){top:17px!important}

  body.admin-drawer-open .mobile-menu-btn-area{
    transform:rotateX(360deg)!important;
  }

  body.admin-drawer-open .mobile-menu-btn-area span:nth-child(1){
    width:18px!important;
    top:4px!important;
    left:4px!important;
    transform:translateY(6px) rotate(-135deg)!important;
  }

  body.admin-drawer-open .mobile-menu-btn-area span:nth-child(2){
    opacity:0!important;
    transform:scaleX(.25)!important;
  }

  body.admin-drawer-open .mobile-menu-btn-area span:nth-child(3){
    width:18px!important;
    top:16px!important;
    left:4px!important;
    transform:translateY(-6px) rotate(135deg)!important;
  }
}

@media(max-width:760px) and (prefers-reduced-motion: reduce), (min-width:768px) and (max-width:1199px) and (prefers-reduced-motion: reduce){
  .mobile-menu-btn-area,
  .mobile-menu-btn-area span{
    transition:none!important;
  }
}
/* /ADMIN_HAMBURGER_X_MOTION_V1 */

/* ADMIN_DRAWER_BELOW_TOPBAR_V1 */
@media(max-width:760px), (min-width:768px) and (max-width:1199px){
  :root{
    --mobile-adminbar-space:68px;
  }

  .mobile-adminbar{
    position:sticky!important;
    top:0!important;
    z-index:130!important;
  }

  .mobile-drawer-overlay{
    top:var(--mobile-adminbar-space)!important;
    z-index:100!important;
  }

  .mobile-drawer{
    top:var(--mobile-adminbar-space)!important;
    bottom:0!important;
    height:calc(100dvh - var(--mobile-adminbar-space))!important;
    max-height:calc(100dvh - var(--mobile-adminbar-space))!important;
    z-index:120!important;
    border-top:1px solid var(--line)!important;
    border-radius:0 22px 0 0!important;
  }

  body.admin-drawer-open .mobile-adminbar{
    box-shadow:
      0 14px 32px rgba(2,8,23,.18),
      inset 0 1px 0 rgba(255,255,255,.08)!important;
  }

  body.admin-drawer-open .mobile-menu-btn{
    z-index:131!important;
  }
}

@media(max-width:420px){
  :root{
    --mobile-adminbar-space:64px;
  }
}
/* /ADMIN_DRAWER_BELOW_TOPBAR_V1 */

/* ADMIN_DRAWER_HAMBURGER_TOGGLE_V1 */
@media(max-width:760px), (min-width:768px) and (max-width:1199px){
  .mobile-drawer-head{
    justify-content:flex-start!important;
  }

  .mobile-drawer-close{
    display:none!important;
  }
}
/* TOGGLE_SLIDER_OVERRIDE */
button.admin-theme-toggle, button.mobile-theme-btn {
    position: relative !important;
    width: 40px !important;
    min-width: 40px !important;
    max-width: 40px !important;
    height: 22px !important;
    min-height: 22px !important;
    max-height: 22px !important;
    border-radius: 22px !important;
    background: #cbd5e1 !important;
    border: none !important;
    padding: 0 !important;
    box-shadow: none !important;
    transition: background 0.3s ease !important;
    display: inline-flex !important;
    align-items: center !important;
    overflow: hidden !important;
    outline: none !important;
}

html[data-admin-theme="dark"] button.admin-theme-toggle,
html[data-admin-theme="dark"] button.mobile-theme-btn,
html[data-theme="dark"] button.admin-theme-toggle,
html[data-theme="dark"] button.mobile-theme-btn {
    background: #2563eb !important;
}

button.admin-theme-toggle::before, button.mobile-theme-btn::before {
    content: "" !important;
    position: absolute !important;
    width: 16px !important;
    height: 16px !important;
    border-radius: 50% !important;
    background: #fff !important;
    left: 3px !important;
    top: 3px !important;
    box-shadow: 0 2px 5px rgba(0,0,0,0.2) !important;
    transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1) !important;
    z-index: 2 !important;
}

html[data-admin-theme="dark"] button.admin-theme-toggle::before,
html[data-admin-theme="dark"] button.mobile-theme-btn::before,
html[data-theme="dark"] button.admin-theme-toggle::before,
html[data-theme="dark"] button.mobile-theme-btn::before {
    transform: translateX(18px) !important;
}

button.admin-theme-toggle .sun, button.mobile-theme-btn .sun {
    position: absolute !important;
    left: 6px !important;
    font-size: 11px !important;
    z-index: 1 !important;
    opacity: 1 !important;
    display: inline-block !important;
    transition: opacity 0.3s !important;
}

button.admin-theme-toggle .moon, button.mobile-theme-btn .moon {
    position: absolute !important;
    right: 5px !important;
    font-size: 11px !important;
    z-index: 1 !important;
    opacity: 0 !important;
    display: inline-block !important;
    transition: opacity 0.3s !important;
}

html[data-admin-theme="dark"] button.admin-theme-toggle .sun,
html[data-admin-theme="dark"] button.mobile-theme-btn .sun,
html[data-theme="dark"] button.admin-theme-toggle .sun,
html[data-theme="dark"] button.mobile-theme-btn .sun {
    opacity: 0 !important;
}

html[data-admin-theme="dark"] button.admin-theme-toggle .moon,
html[data-admin-theme="dark"] button.mobile-theme-btn .moon,
html[data-theme="dark"] button.admin-theme-toggle .moon,
html[data-theme="dark"] button.mobile-theme-btn .moon {
    opacity: 1 !important;
}
/* /TOGGLE_SLIDER_OVERRIDE */
`
