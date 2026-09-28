package main

const publicPageCSS = `:root{
  --bg:#f5f7fb;
  --card:#ffffff;
  --ink:#0f172a;
  --muted:#64748b;
  --pri:#2563eb;
  --pri2:#1d4ed8;
  --line:#e5e7eb;
  --shadow:0 14px 40px rgba(15,23,42,.08);
}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--ink);font-family:system-ui,-apple-system,Segoe UI,Arial,sans-serif}
a{text-decoration:none;color:var(--pri)}
.wrap{max-width:980px;margin:auto;padding:16px}
.PUBLIC_LAYOUT_ACTIVE{display:none}
.public-top{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:14px}
.public-brand{font-size:22px;font-weight:950;letter-spacing:-.04em;max-width:78%;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.admin-mini{display:inline-flex;align-items:center;justify-content:center;height:30px;padding:0 9px;border-radius:999px;background:#fff;color:#111827;border:1px solid var(--line);font-size:11px;font-weight:850;box-shadow:0 8px 24px rgba(15,23,42,.07)}
.hero{padding:26px;border-radius:26px;background:linear-gradient(135deg,#0f172a,#1d4ed8);color:#fff;box-shadow:var(--shadow);margin-bottom:16px}
.hero h1{margin:0;font-size:32px;letter-spacing:-.05em}
.hero p{opacity:.9}
.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px}
.card{background:var(--card);border:1px solid var(--line);border-radius:22px;padding:18px;box-shadow:var(--shadow);margin-bottom:14px}
.price{font-size:30px;font-weight:950;letter-spacing:-.05em;margin:8px 0}
.muted{color:var(--muted)}
.field{display:flex;flex-direction:column;gap:7px;margin-bottom:12px}
.field label{font-size:13px;color:var(--muted);font-weight:850}
input,select{width:100%;padding:12px 13px;border:1px solid var(--line);border-radius:14px;background:#fff;font-size:16px}
.btn{display:inline-flex;align-items:center;justify-content:center;border:0;border-radius:14px;padding:12px 15px;min-height:44px;background:#111827;color:#fff;font-weight:900;cursor:pointer}
.btn.pri{background:linear-gradient(135deg,#2563eb,#1d4ed8)}
.btn.light{background:#fff;color:#111827;border:1px solid var(--line)}
.actions{display:flex;gap:8px;flex-wrap:wrap}
.badge{display:inline-flex;border-radius:999px;padding:5px 10px;font-size:12px;font-weight:900;background:#eef2ff;color:#3730a3}
.paid,.success{background:#dcfce7;color:#166534}.pending{background:#fef9c3;color:#854d0e}.failed{background:#fee2e2;color:#991b1b}.skipped{background:#e5e7eb;color:#374151}
.invoice-wrap{max-width:760px;margin:auto}
.voucher-box{border:2px dashed #93c5fd;background:#eff6ff;border-radius:22px;padding:22px;text-align:center;margin:16px 0}
.voucher-code{font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:50px;font-weight:950;letter-spacing:.16em}
.copy-text{position:absolute;left:-9999px}
.notice{border-radius:15px;padding:12px 14px;background:#f8fafc;border:1px solid var(--line)}
@media(max-width:760px){
  .wrap{padding:14px}
  .public-brand{font-size:20px}
  .admin-mini{position:fixed;top:10px;right:10px;z-index:50;background:rgba(255,255,255,.92);backdrop-filter:blur(8px)}
  .hero{padding:22px;border-radius:22px}
  .hero h1{font-size:28px}
  .grid{grid-template-columns:1fr}
  .card{border-radius:18px}
  .voucher-code{font-size:42px}
}

/* PUBLIC_UI_POLISH_V1 */
.public-top{margin-bottom:10px}
.public-brand{line-height:1.1}
.admin-mini{
  width:32px!important;
  height:32px!important;
  min-width:32px!important;
  padding:0!important;
  font-size:14px!important;
  line-height:1!important;
}
.hero{
  padding:22px 22px!important;
  border-radius:24px!important;
  margin-bottom:14px!important;
}
.hero h1{
  font-size:30px!important;
  line-height:1.05!important;
}
.hero p{
  margin:8px 0 0!important;
  font-size:15px!important;
}
.grid{
  grid-template-columns:repeat(auto-fit,minmax(260px,1fr))!important;
  gap:14px!important;
}
.card{
  padding:18px!important;
  border-radius:22px!important;
}
.card h2,.card h3{
  margin:0 0 6px!important;
  font-size:22px!important;
  line-height:1.1!important;
}
.price{
  font-size:32px!important;
  margin:6px 0 8px!important;
}
.field{
  margin-bottom:10px!important;
}
.field label{
  font-size:12px!important;
}
input,select{
  min-height:46px!important;
  font-size:16px!important;
}
.card form{
  margin-top:12px!important;
}
.card form .btn,.buy-btn{
  width:100%!important;
  min-height:46px!important;
  margin-top:4px!important;
}
@media(max-width:760px){
  .wrap{
    padding:12px 12px calc(18px + env(safe-area-inset-bottom))!important;
  }
  .public-top{
    margin-bottom:8px!important;
  }
  .public-brand{
    font-size:18px!important;
    max-width:72%!important;
  }
  .admin-mini{
    position:fixed!important;
    top:10px!important;
    right:10px!important;
    z-index:50!important;
    width:30px!important;
    height:30px!important;
    min-width:30px!important;
    padding:0!important;
    border-radius:999px!important;
    background:rgba(255,255,255,.92)!important;
    backdrop-filter:blur(8px)!important;
  }
  .hero{
    padding:16px 15px!important;
    border-radius:19px!important;
    margin-bottom:12px!important;
  }
  .hero h1{
    font-size:23px!important;
    line-height:1.05!important;
  }
  .hero p{
    font-size:13px!important;
    margin-top:6px!important;
  }
  .grid{
    grid-template-columns:1fr!important;
    gap:12px!important;
  }
  .card{
    padding:15px!important;
    border-radius:18px!important;
  }
  .card h2,.card h3{
    font-size:20px!important;
  }
  .price{
    font-size:28px!important;
  }
  input,select{
    min-height:44px!important;
  }
  .card form .btn,.buy-btn{
    min-height:46px!important;
  }
}


/* PUBLIC_CARD_COMPACT_V1 */
@media(max-width:760px){
  .wrap{
    padding:10px 10px calc(14px + env(safe-area-inset-bottom))!important;
  }

  .public-top{
    margin-bottom:7px!important;
  }

  .public-brand{
    font-size:18px!important;
    line-height:1.1!important;
  }

  .admin-mini{
    width:28px!important;
    height:28px!important;
    min-width:28px!important;
    padding:0!important;
    font-size:12px!important;
  }

  .hero{
    padding:14px 14px!important;
    border-radius:17px!important;
    margin-bottom:9px!important;
  }

  .hero h1{
    font-size:22px!important;
    line-height:1.05!important;
  }

  .hero p{
    font-size:12px!important;
    margin:5px 0 0!important;
  }

  .grid{
    gap:8px!important;
  }

  .card{
    padding:12px!important;
    border-radius:15px!important;
    margin-bottom:8px!important;
  }

  .card h2,
  .card h3{
    font-size:18px!important;
    line-height:1.1!important;
    margin:0 0 3px!important;
  }

  .price{
    font-size:25px!important;
    line-height:1!important;
    margin:4px 0 6px!important;
  }

  .muted{
    font-size:12.5px!important;
    line-height:1.25!important;
    margin:4px 0 8px!important;
  }

  .btn,
  .buy-btn{
    min-height:39px!important;
    padding:9px 12px!important;
    border-radius:12px!important;
    font-size:14px!important;
  }
}


/* PUBLIC_COMPACT_DARK_V2 */
@media(max-width:760px){
  body{
    font-size:13px!important;
  }

  .wrap{
    padding:7px 8px calc(10px + env(safe-area-inset-bottom))!important;
  }

  .public-top{
    margin-bottom:5px!important;
    min-height:24px!important;
  }

  .public-brand{
    font-size:16px!important;
    line-height:1!important;
    max-width:74%!important;
  }

  .admin-mini{
    width:24px!important;
    height:24px!important;
    min-width:24px!important;
    padding:0!important;
    font-size:11px!important;
    top:7px!important;
    right:7px!important;
    box-shadow:0 4px 14px rgba(15,23,42,.10)!important;
  }

  .hero{
    padding:12px 12px!important;
    border-radius:14px!important;
    margin-bottom:7px!important;
  }

  .hero h1{
    font-size:20px!important;
    line-height:1!important;
    letter-spacing:-.04em!important;
  }

  .hero p{
    font-size:11.5px!important;
    margin:4px 0 0!important;
    line-height:1.2!important;
  }

  .grid{
    gap:7px!important;
  }

  .card{
    padding:10px!important;
    border-radius:12px!important;
    margin-bottom:7px!important;
  }

  .card h2,
  .card h3{
    font-size:16px!important;
    line-height:1!important;
    margin:0 0 2px!important;
  }

  .price{
    font-size:22px!important;
    line-height:1!important;
    margin:3px 0 4px!important;
  }

  .muted{
    font-size:10.5px!important;
    line-height:1.2!important;
    margin:2px 0 6px!important;
  }

  .btn,
  .buy-btn{
    min-height:32px!important;
    padding:7px 10px!important;
    border-radius:9px!important;
    font-size:12px!important;
    line-height:1!important;
  }
}

@media(prefers-color-scheme: dark){
  :root{
    --bg:#0b1220;
    --card:#111827;
    --ink:#e5e7eb;
    --muted:#94a3b8;
    --line:#1f2937;
    --shadow:0 14px 40px rgba(0,0,0,.25);
  }

  body{
    background:var(--bg)!important;
    color:var(--ink)!important;
  }

  .public-brand,
  .brand{
    color:var(--ink)!important;
  }

  .card{
    background:linear-gradient(180deg,#111827,#0f172a)!important;
    border-color:#1f2937!important;
  }

  .admin-mini,
  input,
  select,
  .btn.light,
  .notice{
    background:#0f172a!important;
    color:#e5e7eb!important;
    border-color:#1f2937!important;
  }

  input::placeholder{
    color:#64748b!important;
  }

  .hero{
    background:linear-gradient(135deg,#020617,#1d4ed8)!important;
  }

  .muted{
    color:#94a3b8!important;
  }
}


/* PUBLIC_THEME_TOGGLE_V1 */
.public-actions{
  display:flex;
  align-items:center;
  justify-content:flex-end;
  gap:7px;
}
.theme-toggle{
  display:inline-flex;
  align-items:center;
  justify-content:center;
  gap:3px;
  width:44px;
  height:30px;
  padding:0;
  border-radius:999px;
  border:1px solid var(--line);
  background:#fff;
  color:#111827;
  box-shadow:0 8px 24px rgba(15,23,42,.07);
  cursor:pointer;
  font-size:12px;
  font-weight:850;
}
html[data-theme="light"] body{
  background:#f5f7fb!important;
  color:#0f172a!important;
}
html[data-theme="light"] .card{
  background:#fff!important;
  color:#0f172a!important;
  border-color:#e5e7eb!important;
}
html[data-theme="light"] input,
html[data-theme="light"] select,
html[data-theme="light"] .admin-mini,
html[data-theme="light"] .theme-toggle,
html[data-theme="light"] .btn.light,
html[data-theme="light"] .notice{
  background:#fff!important;
  color:#111827!important;
  border-color:#e5e7eb!important;
}
html[data-theme="light"] .public-brand,
html[data-theme="light"] .brand{
  color:#0f172a!important;
}
html[data-theme="light"] .hero{
  background:linear-gradient(135deg,#0f172a,#1d4ed8)!important;
}
html[data-theme="light"] .muted{
  color:#64748b!important;
}
html[data-theme="light"] .moon{
  opacity:.35;
}
html[data-theme="dark"]{
  --bg:#0b1220;
  --card:#111827;
  --ink:#e5e7eb;
  --muted:#94a3b8;
  --line:#1f2937;
  --shadow:0 14px 40px rgba(0,0,0,.25);
}
html[data-theme="dark"] body{
  background:#0b1220!important;
  color:#e5e7eb!important;
}
html[data-theme="dark"] .card{
  background:linear-gradient(180deg,#111827,#0f172a)!important;
  color:#e5e7eb!important;
  border-color:#1f2937!important;
}
html[data-theme="dark"] input,
html[data-theme="dark"] select,
html[data-theme="dark"] .admin-mini,
html[data-theme="dark"] .theme-toggle,
html[data-theme="dark"] .btn.light,
html[data-theme="dark"] .notice{
  background:#0f172a!important;
  color:#e5e7eb!important;
  border-color:#1f2937!important;
}
html[data-theme="dark"] input::placeholder{
  color:#64748b!important;
}
html[data-theme="dark"] .public-brand,
html[data-theme="dark"] .brand{
  color:#e5e7eb!important;
}
html[data-theme="dark"] .hero{
  background:linear-gradient(135deg,#020617,#1d4ed8)!important;
}
html[data-theme="dark"] .muted{
  color:#94a3b8!important;
}
html[data-theme="dark"] .sun{
  opacity:.35;
}
@media(max-width:760px){
  .public-actions{
    position:fixed!important;
    top:7px!important;
    right:7px!important;
    z-index:60!important;
    gap:5px!important;
  }
  .theme-toggle{
    width:39px!important;
    height:24px!important;
    font-size:10px!important;
    background:rgba(255,255,255,.92)!important;
    backdrop-filter:blur(8px)!important;
  }
  html[data-theme="dark"] .theme-toggle{
    background:rgba(15,23,42,.92)!important;
  }
}


/* PUBLIC_2COL_ANIM_V1 */
@keyframes tukuPageIn{
  from{opacity:0;transform:translateY(8px)}
  to{opacity:1;transform:translateY(0)}
}
@keyframes tukuCardIn{
  from{opacity:0;transform:translateY(10px) scale(.985)}
  to{opacity:1;transform:translateY(0) scale(1)}
}
.wrap{
  animation:tukuPageIn .35s ease-out both;
}
.hero,
.card{
  animation:tukuCardIn .38s ease-out both;
}
.card:nth-child(2){animation-delay:.04s}
.card:nth-child(3){animation-delay:.08s}
.card:nth-child(4){animation-delay:.12s}
.card:nth-child(5){animation-delay:.16s}
@media(prefers-reduced-motion: reduce){
  .wrap,.hero,.card{animation:none!important}
}

@media(max-width:760px){
  .wrap{
    padding:7px 7px calc(10px + env(safe-area-inset-bottom))!important;
  }

  .hero{
    padding:11px 12px!important;
    border-radius:13px!important;
    margin-bottom:7px!important;
  }

  .hero h1{
    font-size:19px!important;
  }

  .hero p{
    font-size:11px!important;
    margin-top:4px!important;
  }

  .grid{
    display:grid!important;
    grid-template-columns:repeat(2,minmax(0,1fr))!important;
    gap:7px!important;
  }

  .card{
    padding:9px!important;
    border-radius:13px!important;
    margin-bottom:0!important;
    min-width:0!important;
  }

  .card h2,
  .card h3{
    font-size:13.5px!important;
    line-height:1.05!important;
    margin:0 0 4px!important;
    white-space:nowrap!important;
    overflow:hidden!important;
    text-overflow:ellipsis!important;
  }

  .price{
    font-size:20px!important;
    line-height:1!important;
    margin:2px 0 5px!important;
    letter-spacing:-.04em!important;
  }

  .muted{
    font-size:9.5px!important;
    line-height:1.15!important;
    margin:2px 0 7px!important;
    white-space:nowrap!important;
    overflow:hidden!important;
    text-overflow:ellipsis!important;
  }

  .btn,
  .buy-btn{
    min-height:30px!important;
    padding:7px 8px!important;
    border-radius:9px!important;
    font-size:11px!important;
    line-height:1!important;
    width:100%!important;
  }

  .public-brand{
    font-size:16px!important;
  }

  .admin-mini{
    width:24px!important;
    height:24px!important;
    min-width:24px!important;
  }

  .theme-toggle{
    width:38px!important;
    height:24px!important;
  }
}


/* PUBLIC_MODERN_COLORS_V1 */
:root{
  --bg:#f8fafc;
  --card:#ffffff;
  --ink:#0f172a;
  --muted:#64748b;
  --pri:#2563eb;
  --pri2:#06b6d4;
  --accent:#7c3aed;
  --line:#e2e8f0;
  --shadow:0 18px 45px rgba(15,23,42,.10);
}

html[data-theme="light"] body{
  background:
    radial-gradient(circle at top left, rgba(37,99,235,.12), transparent 28%),
    radial-gradient(circle at top right, rgba(6,182,212,.16), transparent 30%),
    linear-gradient(180deg,#f8fafc,#eef2ff)!important;
  color:#0f172a!important;
}

html[data-theme="light"] .hero{
  background:
    radial-gradient(circle at top right, rgba(255,255,255,.28), transparent 25%),
    linear-gradient(135deg,#0f172a 0%,#1d4ed8 46%,#06b6d4 100%)!important;
  box-shadow:0 18px 45px rgba(37,99,235,.22)!important;
}

html[data-theme="light"] .card{
  background:rgba(255,255,255,.92)!important;
  border-color:rgba(226,232,240,.95)!important;
  box-shadow:0 14px 34px rgba(15,23,42,.08)!important;
}

html[data-theme="light"] .price{
  color:#0f172a!important;
}

html[data-theme="light"] .btn.pri,
html[data-theme="light"] .buy-btn{
  background:linear-gradient(135deg,#2563eb 0%,#7c3aed 55%,#06b6d4 120%)!important;
  box-shadow:0 10px 22px rgba(37,99,235,.24)!important;
}

html[data-theme="light"] input,
html[data-theme="light"] select{
  background:#ffffff!important;
  border-color:#dbe3ef!important;
}

html[data-theme="light"] .admin-mini,
html[data-theme="light"] .theme-toggle{
  background:rgba(255,255,255,.86)!important;
  border-color:rgba(226,232,240,.95)!important;
  box-shadow:0 10px 28px rgba(15,23,42,.10)!important;
}

html[data-theme="dark"] body{
  background:
    radial-gradient(circle at top left, rgba(37,99,235,.20), transparent 26%),
    radial-gradient(circle at top right, rgba(124,58,237,.20), transparent 30%),
    linear-gradient(180deg,#020617,#0b1120)!important;
  color:#e5e7eb!important;
}

html[data-theme="dark"] .hero{
  background:
    radial-gradient(circle at top right, rgba(6,182,212,.25), transparent 26%),
    linear-gradient(135deg,#020617 0%,#312e81 48%,#0891b2 115%)!important;
  box-shadow:0 18px 45px rgba(6,182,212,.14)!important;
}

html[data-theme="dark"] .card{
  background:linear-gradient(180deg,rgba(15,23,42,.96),rgba(2,6,23,.96))!important;
  border-color:rgba(51,65,85,.9)!important;
  box-shadow:0 18px 45px rgba(0,0,0,.28)!important;
}

html[data-theme="dark"] .price{
  color:#f8fafc!important;
}

html[data-theme="dark"] .btn.pri,
html[data-theme="dark"] .buy-btn{
  background:linear-gradient(135deg,#2563eb 0%,#7c3aed 58%,#06b6d4 125%)!important;
  box-shadow:0 10px 24px rgba(37,99,235,.22)!important;
}

html[data-theme="dark"] input,
html[data-theme="dark"] select{
  background:#0b1120!important;
  border-color:#334155!important;
  color:#e5e7eb!important;
}

html[data-theme="dark"] .admin-mini,
html[data-theme="dark"] .theme-toggle{
  background:rgba(15,23,42,.88)!important;
  border-color:rgba(51,65,85,.95)!important;
  box-shadow:0 10px 28px rgba(0,0,0,.24)!important;
}

@media(max-width:760px){
  html[data-theme="light"] .hero{
    box-shadow:0 10px 24px rgba(37,99,235,.18)!important;
  }

  html[data-theme="light"] .card{
    box-shadow:0 8px 20px rgba(15,23,42,.07)!important;
  }

  html[data-theme="dark"] .card{
    box-shadow:0 10px 24px rgba(0,0,0,.24)!important;
  }
}


/* PUBLIC_CODE_BACKGROUND_V1 */
body{
  position:relative!important;
  overflow-x:hidden!important;
}

body::before{
  content:""!important;
  position:fixed!important;
  inset:0!important;
  z-index:-2!important;
  pointer-events:none!important;
  background:
    radial-gradient(circle at 12% 8%, rgba(37,99,235,.16), transparent 28%),
    radial-gradient(circle at 88% 10%, rgba(6,182,212,.16), transparent 30%),
    radial-gradient(circle at 50% 96%, rgba(124,58,237,.12), transparent 34%),
    linear-gradient(180deg,#f8fafc 0%,#eef2ff 100%)!important;
}

body::after{
  content:""!important;
  position:fixed!important;
  inset:0!important;
  z-index:-1!important;
  pointer-events:none!important;
  opacity:.55!important;
  background-image:
    linear-gradient(rgba(37,99,235,.055) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37,99,235,.055) 1px, transparent 1px),
    radial-gradient(circle at center, rgba(15,23,42,.08) 1px, transparent 1.5px)!important;
  background-size:
    34px 34px,
    34px 34px,
    18px 18px!important;
  mask-image:linear-gradient(180deg, rgba(0,0,0,.85), rgba(0,0,0,.25))!important;
}

html[data-theme="light"] body::before{
  background:
    radial-gradient(circle at 12% 8%, rgba(37,99,235,.18), transparent 28%),
    radial-gradient(circle at 88% 10%, rgba(6,182,212,.18), transparent 30%),
    radial-gradient(circle at 50% 96%, rgba(124,58,237,.12), transparent 34%),
    linear-gradient(180deg,#f8fafc 0%,#eef2ff 100%)!important;
}

html[data-theme="light"] body::after{
  opacity:.50!important;
}

html[data-theme="dark"] body::before{
  background:
    radial-gradient(circle at 10% 8%, rgba(37,99,235,.26), transparent 30%),
    radial-gradient(circle at 90% 12%, rgba(6,182,212,.18), transparent 30%),
    radial-gradient(circle at 50% 98%, rgba(124,58,237,.22), transparent 36%),
    linear-gradient(180deg,#020617 0%,#0b1120 100%)!important;
}

html[data-theme="dark"] body::after{
  opacity:.32!important;
  background-image:
    linear-gradient(rgba(148,163,184,.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148,163,184,.08) 1px, transparent 1px),
    radial-gradient(circle at center, rgba(148,163,184,.10) 1px, transparent 1.5px)!important;
}

.hero,
.card,
.admin-mini,
.theme-toggle{
  backdrop-filter:blur(10px)!important;
}

html[data-theme="light"] .card{
  background:rgba(255,255,255,.86)!important;
}

html[data-theme="dark"] .card{
  background:linear-gradient(180deg,rgba(15,23,42,.88),rgba(2,6,23,.90))!important;
}

@media(max-width:760px){
  body::after{
    background-size:
      26px 26px,
      26px 26px,
      16px 16px!important;
    opacity:.38!important;
  }

  body::before{
    background:
      radial-gradient(circle at 8% 4%, rgba(37,99,235,.18), transparent 26%),
      radial-gradient(circle at 96% 8%, rgba(6,182,212,.16), transparent 28%),
      linear-gradient(180deg,#f8fafc 0%,#eef2ff 100%)!important;
  }

  html[data-theme="dark"] body::before{
    background:
      radial-gradient(circle at 8% 4%, rgba(37,99,235,.22), transparent 27%),
      radial-gradient(circle at 96% 8%, rgba(124,58,237,.18), transparent 30%),
      linear-gradient(180deg,#020617 0%,#0b1120 100%)!important;
  }
}


/* PUBLIC_BG_STRONG_V2 */
html[data-theme="light"] body{
  background:
    radial-gradient(circle at 18% 8%, rgba(37,99,235,.30) 0, transparent 26%),
    radial-gradient(circle at 88% 14%, rgba(6,182,212,.30) 0, transparent 28%),
    radial-gradient(circle at 45% 92%, rgba(124,58,237,.18) 0, transparent 35%),
    linear-gradient(180deg,#eef6ff 0%,#f8fafc 45%,#eef2ff 100%)!important;
}

html[data-theme="dark"] body{
  background:
    radial-gradient(circle at 16% 8%, rgba(37,99,235,.35) 0, transparent 28%),
    radial-gradient(circle at 88% 14%, rgba(6,182,212,.24) 0, transparent 30%),
    radial-gradient(circle at 45% 95%, rgba(124,58,237,.28) 0, transparent 36%),
    linear-gradient(180deg,#020617 0%,#0f172a 55%,#111827 100%)!important;
}

body::before{
  content:""!important;
  position:fixed!important;
  inset:0!important;
  z-index:-2!important;
  pointer-events:none!important;
  background-image:
    linear-gradient(rgba(37,99,235,.10) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37,99,235,.10) 1px, transparent 1px)!important;
  background-size:32px 32px!important;
  opacity:.55!important;
}

body::after{
  content:""!important;
  position:fixed!important;
  inset:auto -80px -120px auto!important;
  width:260px!important;
  height:260px!important;
  z-index:-1!important;
  pointer-events:none!important;
  border-radius:999px!important;
  background:linear-gradient(135deg,rgba(37,99,235,.28),rgba(6,182,212,.22),rgba(124,58,237,.18))!important;
  filter:blur(28px)!important;
}

html[data-theme="dark"] body::before{
  background-image:
    linear-gradient(rgba(148,163,184,.075) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148,163,184,.075) 1px, transparent 1px)!important;
  opacity:.35!important;
}

html[data-theme="light"] .card{
  background:rgba(255,255,255,.82)!important;
  backdrop-filter:blur(12px)!important;
}

html[data-theme="dark"] .card{
  background:rgba(15,23,42,.82)!important;
  backdrop-filter:blur(12px)!important;
}

@media(max-width:760px){
  body::before{
    background-size:24px 24px!important;
    opacity:.48!important;
  }

  body::after{
    width:190px!important;
    height:190px!important;
    inset:auto -70px -80px auto!important;
    filter:blur(24px)!important;
  }

  html[data-theme="light"] body{
    background:
      radial-gradient(circle at 10% 5%, rgba(37,99,235,.34) 0, transparent 28%),
      radial-gradient(circle at 95% 10%, rgba(6,182,212,.30) 0, transparent 30%),
      linear-gradient(180deg,#eef6ff 0%,#f8fafc 45%,#eef2ff 100%)!important;
  }
}


/* PUBLIC_BG_VISIBLE_V3 */
html,
body{
  min-height:100%!important;
}

body{
  background-color:#dbeafe!important;
  background-image:
    radial-gradient(circle at 8% 8%, rgba(37,99,235,.42) 0 0, transparent 24%),
    radial-gradient(circle at 92% 10%, rgba(6,182,212,.38) 0 0, transparent 26%),
    radial-gradient(circle at 50% 100%, rgba(124,58,237,.32) 0 0, transparent 35%),
    linear-gradient(135deg, rgba(37,99,235,.10) 25%, transparent 25%),
    linear-gradient(225deg, rgba(6,182,212,.10) 25%, transparent 25%),
    linear-gradient(180deg,#dbeafe 0%,#f8fafc 48%,#ede9fe 100%)!important;
  background-size:
    auto,
    auto,
    auto,
    34px 34px,
    34px 34px,
    auto!important;
  background-attachment:fixed!important;
}

html[data-theme="light"] body{
  background-color:#dbeafe!important;
  background-image:
    radial-gradient(circle at 8% 8%, rgba(37,99,235,.42) 0 0, transparent 24%),
    radial-gradient(circle at 92% 10%, rgba(6,182,212,.38) 0 0, transparent 26%),
    radial-gradient(circle at 50% 100%, rgba(124,58,237,.32) 0 0, transparent 35%),
    linear-gradient(135deg, rgba(37,99,235,.10) 25%, transparent 25%),
    linear-gradient(225deg, rgba(6,182,212,.10) 25%, transparent 25%),
    linear-gradient(180deg,#dbeafe 0%,#f8fafc 48%,#ede9fe 100%)!important;
}

html[data-theme="dark"] body{
  background-color:#020617!important;
  background-image:
    radial-gradient(circle at 8% 8%, rgba(37,99,235,.45) 0 0, transparent 25%),
    radial-gradient(circle at 92% 10%, rgba(6,182,212,.30) 0 0, transparent 28%),
    radial-gradient(circle at 50% 100%, rgba(124,58,237,.38) 0 0, transparent 38%),
    linear-gradient(135deg, rgba(96,165,250,.08) 25%, transparent 25%),
    linear-gradient(225deg, rgba(34,211,238,.07) 25%, transparent 25%),
    linear-gradient(180deg,#020617 0%,#0f172a 55%,#111827 100%)!important;
  background-size:
    auto,
    auto,
    auto,
    34px 34px,
    34px 34px,
    auto!important;
}

body::before{
  content:""!important;
  position:fixed!important;
  inset:0!important;
  pointer-events:none!important;
  z-index:0!important;
  opacity:.35!important;
  background-image:
    radial-gradient(circle, rgba(15,23,42,.18) 1.2px, transparent 1.4px)!important;
  background-size:18px 18px!important;
}

html[data-theme="dark"] body::before{
  opacity:.22!important;
  background-image:
    radial-gradient(circle, rgba(226,232,240,.22) 1.1px, transparent 1.4px)!important;
}

.wrap{
  position:relative!important;
  z-index:1!important;
}

html[data-theme="light"] .card{
  background:rgba(255,255,255,.78)!important;
  border-color:rgba(255,255,255,.72)!important;
  backdrop-filter:blur(14px)!important;
  box-shadow:0 12px 28px rgba(30,64,175,.12)!important;
}

html[data-theme="dark"] .card{
  background:rgba(15,23,42,.78)!important;
  border-color:rgba(51,65,85,.75)!important;
  backdrop-filter:blur(14px)!important;
}

html[data-theme="light"] .hero{
  background:
    radial-gradient(circle at top right, rgba(255,255,255,.30), transparent 28%),
    linear-gradient(135deg,#172554 0%,#2563eb 48%,#06b6d4 100%)!important;
}

html[data-theme="dark"] .hero{
  background:
    radial-gradient(circle at top right, rgba(34,211,238,.30), transparent 30%),
    linear-gradient(135deg,#020617 0%,#312e81 48%,#0891b2 100%)!important;
}

@media(max-width:760px){
  body{
    background-size:
      auto,
      auto,
      auto,
      24px 24px,
      24px 24px,
      auto!important;
  }

  body::before{
    opacity:.30!important;
    background-size:14px 14px!important;
  }

  html[data-theme="light"] .card{
    background:rgba(255,255,255,.74)!important;
  }

  html[data-theme="dark"] .card{
    background:rgba(15,23,42,.74)!important;
  }
}


/* PUBLIC_ANIMATED_BG_V1 */
@keyframes tukuAuroraMove{
  0%{
    background-position:0% 40%,100% 20%,50% 100%,0% 0%;
  }
  50%{
    background-position:65% 15%,30% 60%,80% 70%,50% 50%;
  }
  100%{
    background-position:100% 55%,0% 30%,35% 85%,100% 100%;
  }
}

@keyframes tukuGridMove{
  0%{transform:translate3d(0,0,0)}
  100%{transform:translate3d(-32px,-32px,0)}
}

@keyframes tukuBlobFloat{
  0%{transform:translate3d(0,0,0) scale(1) rotate(0deg)}
  33%{transform:translate3d(-42px,24px,0) scale(1.08) rotate(18deg)}
  66%{transform:translate3d(36px,-28px,0) scale(.96) rotate(-14deg)}
  100%{transform:translate3d(0,0,0) scale(1) rotate(0deg)}
}

html,body{
  min-height:100%!important;
}

body{
  position:relative!important;
  overflow-x:hidden!important;
  background-color:#dbeafe!important;
  background-image:
    radial-gradient(circle at 15% 12%, rgba(37,99,235,.48) 0 0, transparent 26%),
    radial-gradient(circle at 85% 18%, rgba(6,182,212,.42) 0 0, transparent 28%),
    radial-gradient(circle at 50% 96%, rgba(124,58,237,.34) 0 0, transparent 36%),
    linear-gradient(135deg,#dbeafe 0%,#f8fafc 45%,#ede9fe 100%)!important;
  background-size:
    170% 170%,
    180% 180%,
    190% 190%,
    100% 100%!important;
  animation:tukuAuroraMove 14s ease-in-out infinite alternate!important;
}

body::before{
  content:""!important;
  position:fixed!important;
  inset:-40px!important;
  pointer-events:none!important;
  z-index:0!important;
  opacity:.42!important;
  background-image:
    linear-gradient(rgba(37,99,235,.14) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37,99,235,.14) 1px, transparent 1px),
    radial-gradient(circle, rgba(15,23,42,.20) 1.1px, transparent 1.5px)!important;
  background-size:32px 32px,32px 32px,18px 18px!important;
  animation:tukuGridMove 18s linear infinite!important;
}

body::after{
  content:""!important;
  position:fixed!important;
  right:-90px!important;
  bottom:-110px!important;
  width:310px!important;
  height:310px!important;
  pointer-events:none!important;
  z-index:0!important;
  border-radius:999px!important;
  background:
    radial-gradient(circle at 35% 35%, rgba(6,182,212,.45), transparent 42%),
    radial-gradient(circle at 60% 60%, rgba(124,58,237,.35), transparent 48%),
    radial-gradient(circle at 50% 50%, rgba(37,99,235,.30), transparent 62%)!important;
  filter:blur(26px)!important;
  animation:tukuBlobFloat 11s ease-in-out infinite!important;
}

.wrap{
  position:relative!important;
  z-index:1!important;
}

html[data-theme="dark"] body{
  background-color:#020617!important;
  background-image:
    radial-gradient(circle at 15% 12%, rgba(37,99,235,.52) 0 0, transparent 27%),
    radial-gradient(circle at 85% 18%, rgba(6,182,212,.34) 0 0, transparent 30%),
    radial-gradient(circle at 50% 96%, rgba(124,58,237,.42) 0 0, transparent 38%),
    linear-gradient(135deg,#020617 0%,#0f172a 50%,#111827 100%)!important;
}

html[data-theme="dark"] body::before{
  opacity:.26!important;
  background-image:
    linear-gradient(rgba(148,163,184,.10) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148,163,184,.10) 1px, transparent 1px),
    radial-gradient(circle, rgba(226,232,240,.20) 1px, transparent 1.5px)!important;
}

html[data-theme="dark"] body::after{
  background:
    radial-gradient(circle at 35% 35%, rgba(6,182,212,.34), transparent 42%),
    radial-gradient(circle at 60% 60%, rgba(124,58,237,.42), transparent 48%),
    radial-gradient(circle at 50% 50%, rgba(37,99,235,.34), transparent 62%)!important;
}

html[data-theme="light"] .card{
  background:rgba(255,255,255,.72)!important;
  backdrop-filter:blur(16px)!important;
}

html[data-theme="dark"] .card{
  background:rgba(15,23,42,.74)!important;
  backdrop-filter:blur(16px)!important;
}

@media(max-width:760px){
  body{
    background-size:
      145% 145%,
      155% 155%,
      160% 160%,
      100% 100%!important;
    animation-duration:12s!important;
  }

  body::before{
    inset:-24px!important;
    opacity:.34!important;
    background-size:24px 24px,24px 24px,14px 14px!important;
    animation-duration:16s!important;
  }

  body::after{
    width:210px!important;
    height:210px!important;
    right:-80px!important;
    bottom:-70px!important;
    filter:blur(22px)!important;
  }
}

@media(prefers-reduced-motion: reduce){
  body,
  body::before,
  body::after{
    animation:none!important;
  }
}


/* PUBLIC_CAT_ANIM_V1 */
/* Cleaned: hanya CSS .tuku-cat aktif. Wrapper lama .tuku-cat-scene, moon, star sudah tidak dipakai. */

@keyframes tailWiggle{
  0%,100%{transform:rotate(-18deg)}
  50%{transform:rotate(13deg)}
}

@keyframes catBlink{
  0%,88%,100%{transform:scaleY(1)}
  91%,94%{transform:scaleY(.08)}
}

.tuku-cat{
  position:absolute!important;
  left:54px!important;
  bottom:23px!important;
  width:58px!important;
  height:66px!important;
}

.tuku-cat-body{
  position:absolute!important;
  left:10px!important;
  bottom:0!important;
  width:39px!important;
  height:45px!important;
  border-radius:45% 45% 36% 36%!important;
  background:linear-gradient(180deg,#111827,#020617)!important;
}

.tuku-cat-head{
  position:absolute!important;
  left:8px!important;
  bottom:35px!important;
  width:43px!important;
  height:36px!important;
  border-radius:45% 45% 42% 42%!important;
  background:linear-gradient(180deg,#111827,#020617)!important;
}

.tuku-cat-head .ear{
  position:absolute!important;
  top:-10px!important;
  width:17px!important;
  height:17px!important;
  background:#111827!important;
  transform:rotate(45deg)!important;
  border-radius:3px!important;
}

.tuku-cat-head .ear.left{
  left:4px!important;
}

.tuku-cat-head .ear.right{
  right:4px!important;
}

.tuku-cat-head .eye{
  position:absolute!important;
  top:15px!important;
  width:5px!important;
  height:7px!important;
  border-radius:999px!important;
  background:#fef08a!important;
  transform-origin:center!important;
  animation:catBlink 4.5s infinite!important;
}

.tuku-cat-head .eye.left{
  left:12px!important;
}

.tuku-cat-head .eye.right{
  right:12px!important;
}

.tuku-cat-head .nose{
  position:absolute!important;
  left:50%!important;
  top:23px!important;
  width:5px!important;
  height:4px!important;
  transform:translateX(-50%)!important;
  border-radius:999px!important;
  background:#fb7185!important;
}

.tuku-cat-tail{
  position:absolute!important;
  left:-7px!important;
  bottom:10px!important;
  width:16px!important;
  height:42px!important;
  border-radius:999px!important;
  border-left:8px solid #020617!important;
  transform-origin:bottom center!important;
  animation:tailWiggle 2.7s ease-in-out infinite!important;
}

@media(max-width:760px){
  .tuku-cat{
    left:41px!important;
    bottom:17px!important;
    transform:scale(.75)!important;
    transform-origin:bottom center!important;
  }
}

@media(prefers-reduced-motion: reduce){
  .tuku-cat-tail,
  .tuku-cat-head .eye{
    animation:none!important;
  }
}



/* CAT_1TO1_FROM_CODEPEN_V1 */
@keyframes catpenBlink{
  0%,25%,28%,100%{height:0}
  26.5%{height:100%}
}

@keyframes catpenFloat{
  0%,100%{transform:translate3d(0,0,0)}
  50%{transform:translate3d(0,-8px,0)}
}

@keyframes catpenTwinkle{
  0%,100%{filter:brightness(.85)}
  50%{filter:brightness(1.25)}
}

.catpen-scene,
.catpen-scene *,
.catpen-scene *::before,
.catpen-scene *::after{
  position:absolute!important;
  box-sizing:border-box!important;
}

.catpen-scene{
  position:fixed!important;
  right:18px!important;
  bottom:-18px!important;
  width:220px!important;
  aspect-ratio:1!important;
  margin:0!important;
  padding:0!important;
  z-index:0!important;
  pointer-events:none!important;
  overflow:visible!important;
  font-size:2.25px!important;
  animation:catpenFloat 5.5s ease-in-out infinite!important;
  opacity:.92!important;
}

.catpen-scene::before{
  content:""!important;
  position:absolute!important;
  inset:0!important;
  z-index:-2!important;
  border-radius:50%!important;
  background:
    radial-gradient(circle at 20% 20%, #ffc8 .20em, #0000 0),
    radial-gradient(circle at 70% 50%, #ffc8 .25em, #0000 0),
    radial-gradient(circle at 10% 80%, #ffc8 .20em, #0000 0),
    radial-gradient(circle at 80% 30%, #ffc8 .30em, #0000 0),
    radial-gradient(circle at 80% 80%, #ffc8 .23em, #0000 0),
    radial-gradient(circle at 50% 10%, #ff8b .3em, #0000 0),
    radial-gradient(circle at 40% 60%, #ffc8 .20em, #0000 0),
    radial-gradient(circle at 30% 30%, #ffc8 .20em, #0000 0) 0 0 / 70% 70%,
    radial-gradient(circle at 30% 30%, #ffc8 .20em, #0000 0) 40% 30% / 50% 70%;
  animation:catpenTwinkle 4s ease-in-out infinite!important;
}

.catpen-scene figcaption,
.catpen-cat #alt{
  width:1px!important;
  height:1px!important;
  overflow:hidden!important;
}

.catpen-moon{
  width:120%!important;
  height:40%!important;
  left:50%!important;
  bottom:0!important;
  translate:-50% 60%!important;
  border-radius:50% / 100% 100% 0 0!important;
  box-shadow:0 0 3em #ffc,0 0 9em #ffc8!important;
  background:
    radial-gradient(at 38% 20%, #0001 2%, #0000 0),
    radial-gradient(13% 8% at 50% 0, #0001 99%, #0000),
    #ffc!important;
}

.catpen-moon::before{
  content:""!important;
  width:18%!important;
  height:15%!important;
  background:#0001!important;
  border-radius:50%!important;
  box-shadow:inset 0 1em #0002!important;
  rotate:-15deg!important;
  top:15%!important;
  left:20%!important;
  scale:.8!important;
}

.catpen-moon::after{
  content:""!important;
  width:4%!important;
  height:5%!important;
  background:#0001!important;
  border-radius:50%!important;
  box-shadow:-7em 3em 0 2em #0001!important;
  rotate:30deg!important;
  top:25%!important;
  left:80%!important;
  scale:.8!important;
}

.catpen-bubble{
  width:119%!important;
  height:103%!important;
  border-radius:50%!important;
  bottom:25%!important;
  left:49%!important;
  translate:-50%!important;
  background:
    radial-gradient(circle at 25% 25%, #fff3, #fff0 25%),
    radial-gradient(circle at 15% 45%, #fff1, #fff0 15%),
    radial-gradient(circle at 75% 35%, #fff1, #fff0 50%),
    #fff1!important;
  box-shadow:inset 0 0 1em 1em #acf,inset 0 0 5em #fff!important;
}

.catpen-backpack{
  width:40%!important;
  left:50%!important;
  height:30%!important;
  background:#bcd!important;
  bottom:0!important;
  translate:-50%!important;
  border-radius:3em!important;
}

.catpen-cat{
  --fur:#111;
  --fur-dark:#000;
  --skin:pink;
  --suit:#fff;
  --suit-dark:#ddd;
  font-size:1.12px!important;
  width:80em!important;
  aspect-ratio:1!important;
  bottom:15%!important;
  left:50%!important;
  translate:-50%!important;
}

.catpen-tail{
  width:50%!important;
  height:50%!important;
  border-radius:50%!important;
  border:7em solid #0000!important;
  border-top-color:var(--suit-dark)!important;
  border-left-color:var(--suit-dark)!important;
  clip-path:polygon(100% 0,100% 100%,0 30%,0 0)!important;
  top:75%!important;
  left:52%!important;
}

.catpen-tail::before{
  content:""!important;
  width:7em!important;
  aspect-ratio:1!important;
  background:var(--suit-dark)!important;
  border-radius:50%!important;
  left:81%!important;
  top:-9%!important;
}

.catpen-body{
  left:50%!important;
  translate:-50%!important;
  bottom:0!important;
  width:35%!important;
  height:40%!important;
  background:
    radial-gradient(circle at 17% 55%, #36c 2em, #0000 0),
    radial-gradient(100% 70% at 50% 0, var(--fur-dark) 50%, #0000 0),
    radial-gradient(150% 70% at 49% 0, var(--fur) 50%, #d99 0 59%, #0000 calc(59% + 1px)),
    var(--suit)!important;
  border-radius:100% / 200% 200% 20% 20%!important;
}

.catpen-leg{
  width:165%!important;
  height:38%!important;
  bottom:0!important;
  left:50%!important;
  translate:-50%!important;
  border-radius:8em 8em 100% 100% / 10em 10em 16em 16em!important;
  scale:1 -1!important;
  background:
    radial-gradient(20% 200% at 2% 10%, var(--suit-dark) 20%, #0000 0),
    radial-gradient(20% 200% at 98% 10%, var(--suit-dark) 20%, #0000 0),
    var(--suit)!important;
}

.catpen-paw{
  width:35%!important;
  height:49%!important;
  border:.75em solid #777!important;
  border-top:0!important;
  border-radius:0 0 5em 5em!important;
  top:51%!important;
  rotate:-10deg!important;
  left:10%!important;
  clip-path:polygon(0 20%,100% 30%,100% 100%,0 100%)!important;
}

.catpen-paw + .catpen-paw{
  right:7%!important;
  height:50%!important;
  left:auto!important;
  rotate:13deg!important;
  scale:-1 1!important;
  clip-path:polygon(0 25%,100% 15%,100% 100%,0 100%)!important;
}

.catpen-ear{
  width:40%!important;
  aspect-ratio:1!important;
  border:4em solid var(--fur)!important;
  border-radius:5% 90% 10% 80%!important;
  background:var(--skin)!important;
  top:5%!important;
  left:10%!important;
}

.catpen-ear + .catpen-ear{
  scale:-1 1!important;
  left:auto!important;
  right:10%!important;
}

.catpen-head{
  width:80%!important;
  aspect-ratio:1.1!important;
  background:linear-gradient(#0003,#0000 50%),var(--fur)!important;
  left:50%!important;
  translate:-50%!important;
  border-radius:100% / 125% 125% 80% 75%!important;
}

.catpen-whisker{
  width:30%!important;
  height:30%!important;
  border-radius:50%!important;
  border:2em solid #0000!important;
  border-top-color:var(--fur)!important;
  border-left-color:var(--fur)!important;
  clip-path:polygon(100% 0,100% 100%,0 30%,0 0)!important;
}

.catpen-whisker:nth-child(1){
  top:45%!important;
  translate:-65%!important;
}

.catpen-whisker:nth-child(2){
  top:55%!important;
  translate:-40%!important;
  rotate:-20deg!important;
}

.catpen-whisker:nth-child(3){
  right:0%!important;
  top:45%!important;
  translate:65%!important;
  rotate:10deg!important;
}

.catpen-whisker:nth-child(4){
  right:-2%!important;
  top:55%!important;
  translate:40%!important;
  rotate:24deg!important;
}

.catpen-nose{
  width:10%!important;
  height:7%!important;
  background:var(--skin)!important;
  border-radius:50%!important;
  left:50%!important;
  translate:-50% -50%!important;
  top:55%!important;
}

.catpen-eye{
  --pos:25%;
  --x1:50%;
  --x2:42%;
  width:35%!important;
  aspect-ratio:1!important;
  border-radius:50%!important;
  background:
    radial-gradient(50% 50% at var(--x1) 32%, #fff 25%, #0000 calc(25% + 1px)),
    radial-gradient(50% 50% at var(--x2) 51%, #fff 12%, #0000 calc(12% + 1px)),
    radial-gradient(circle at 60% 40%, #000 35%, #0000 calc(35% + 1px)),
    white!important;
  left:var(--pos)!important;
  translate:-50% -50%!important;
  top:43%!important;
  overflow:hidden!important;
}

.catpen-eye + .catpen-eye{
  --x1:70%;
  --x2:78%;
  left:calc(100% - var(--pos))!important;
  scale:-1 1!important;
}

.catpen-eye::before{
  top:-30%!important;
  left:50%!important;
  translate:-50%!important;
  width:150%!important;
  height:0%!important;
  content:""!important;
  background:var(--fur)!important;
  rotate:-10deg!important;
  animation:catpenBlink 10s linear infinite!important;
}

.catpen-eye::after{
  bottom:-10%!important;
  width:150%!important;
  height:0%!important;
  content:""!important;
  background:var(--fur)!important;
  rotate:-10deg!important;
  left:50%!important;
  translate:-50%!important;
  animation:catpenBlink 10s linear infinite!important;
}

html[data-theme="light"] .catpen-scene{
  opacity:.78!important;
}

html[data-theme="dark"] .catpen-scene{
  opacity:.94!important;
}

@media(max-width:760px){
  .catpen-scene{
    width:145px!important;
    right:4px!important;
    bottom:-15px!important;
    opacity:.72!important;
  }

  .catpen-cat{
    font-size:.82px!important;
  }
}

@media(prefers-reduced-motion:reduce){
  .catpen-scene,
  .catpen-eye::before,
  .catpen-eye::after,
  .catpen-scene::before{
    animation:none!important;
  }
}


/* CATPEN_SIZE_HELMET_FIX_V2 */
.catpen-scene{
  width:250px!important;
  right:16px!important;
  bottom:-8px!important;
}

.catpen-bubble{
  width:92%!important;
  height:80%!important;
  left:50%!important;
  bottom:18%!important;
  translate:-50% 0!important;
  border-radius:50%!important;
  box-shadow:
    inset 0 0 .8em .8em rgba(147,197,253,.82),
    inset 0 0 3.6em rgba(255,255,255,.88)!important;
}

.catpen-cat{
  font-size:1.48px!important;
  width:92em!important;
  bottom:14%!important;
  left:50%!important;
  translate:-50% 0!important;
}

.catpen-head{
  width:84%!important;
  top:2%!important;
}

.catpen-ear{
  width:38%!important;
  top:4%!important;
}

.catpen-body{
  width:34%!important;
  height:38%!important;
}

.catpen-eye{
  width:33%!important;
  top:42%!important;
}

.catpen-nose{
  top:56%!important;
}

.catpen-tail{
  top:73%!important;
  left:54%!important;
}

@media(max-width:760px){
  .catpen-scene{
    width:185px!important;
    right:2px!important;
    bottom:-8px!important;
    opacity:.82!important;
  }

  .catpen-bubble{
    width:88%!important;
    height:74%!important;
    bottom:20%!important;
  }

  .catpen-cat{
    font-size:1.12px!important;
    width:94em!important;
    bottom:15%!important;
  }

  .catpen-head{
    width:86%!important;
  }

  .catpen-ear{
    width:40%!important;
  }

  .catpen-eye{
    width:34%!important;
  }
}


/* PUBLIC_BG_NIGHT_V4 */
html,
body{
  min-height:100%!important;
}

html[data-theme="light"] body,
html[data-theme="dark"] body{
  position:relative!important;
  overflow-x:hidden!important;
  background:
    radial-gradient(circle at 50% -15%, rgba(30,64,175,.28) 0%, rgba(30,64,175,.12) 18%, transparent 38%),
    radial-gradient(circle at 85% 14%, rgba(34,211,238,.16) 0%, transparent 26%),
    linear-gradient(180deg,#05263c 0%,#062840 38%,#0b2439 72%,#071a2b 100%)!important;
}

html[data-theme="light"] body::before,
html[data-theme="dark"] body::before{
  content:""!important;
  position:fixed!important;
  inset:0!important;
  z-index:0!important;
  pointer-events:none!important;
  background:
    radial-gradient(circle at 10% 18%, rgba(255,244,180,.95) 0 1px, transparent 1.6px),
    radial-gradient(circle at 22% 12%, rgba(255,255,255,.85) 0 1px, transparent 1.6px),
    radial-gradient(circle at 34% 28%, rgba(255,244,180,.85) 0 1.2px, transparent 1.8px),
    radial-gradient(circle at 45% 10%, rgba(255,255,255,.9) 0 1px, transparent 1.6px),
    radial-gradient(circle at 58% 20%, rgba(255,244,180,.88) 0 1px, transparent 1.7px),
    radial-gradient(circle at 69% 8%, rgba(255,255,255,.84) 0 1px, transparent 1.6px),
    radial-gradient(circle at 80% 24%, rgba(255,244,180,.92) 0 1.2px, transparent 1.8px),
    radial-gradient(circle at 92% 15%, rgba(255,255,255,.82) 0 1px, transparent 1.6px),
    radial-gradient(circle at 14% 44%, rgba(255,255,255,.78) 0 1px, transparent 1.6px),
    radial-gradient(circle at 28% 52%, rgba(255,244,180,.88) 0 1px, transparent 1.7px),
    radial-gradient(circle at 50% 40%, rgba(255,255,255,.75) 0 1px, transparent 1.6px),
    radial-gradient(circle at 72% 48%, rgba(255,244,180,.84) 0 1px, transparent 1.7px),
    radial-gradient(circle at 88% 38%, rgba(255,255,255,.78) 0 1px, transparent 1.6px)!important;
  opacity:.95!important;
}

html[data-theme="light"] body::after,
html[data-theme="dark"] body::after{
  content:""!important;
  position:fixed!important;
  inset:0!important;
  z-index:0!important;
  pointer-events:none!important;
  background:
    radial-gradient(circle at 50% 100%, rgba(255,245,170,.06) 0%, transparent 20%),
    linear-gradient(180deg, rgba(255,255,255,.03), rgba(255,255,255,0) 20%, rgba(0,0,0,.14) 100%)!important;
}

.wrap{
  position:relative!important;
  z-index:1!important;
}

.public-brand{
  color:#f8fafc!important;
  text-shadow:0 2px 10px rgba(0,0,0,.28)!important;
}

html[data-theme="light"] .card,
html[data-theme="dark"] .card{
  background:rgba(8,20,39,.60)!important;
  border-color:rgba(148,163,184,.18)!important;
  box-shadow:0 14px 34px rgba(0,0,0,.24)!important;
  backdrop-filter:blur(10px)!important;
}

html[data-theme="light"] .hero,
html[data-theme="dark"] .hero{
  background:linear-gradient(135deg,rgba(31,41,55,.90),rgba(67,56,202,.82) 50%,rgba(6,182,212,.78) 100%)!important;
  box-shadow:0 18px 40px rgba(0,0,0,.28)!important;
}

html[data-theme="light"] .muted,
html[data-theme="dark"] .muted{
  color:#cbd5e1!important;
}

html[data-theme="light"] .price,
html[data-theme="dark"] .price,
html[data-theme="light"] .card h2,
html[data-theme="dark"] .card h2,
html[data-theme="light"] .card h3,
html[data-theme="dark"] .card h3,
html[data-theme="light"] .hero h1,
html[data-theme="dark"] .hero h1,
html[data-theme="light"] .hero p,
html[data-theme="dark"] .hero p{
  color:#f8fafc!important;
}

html[data-theme="light"] .btn.pri,
html[data-theme="dark"] .btn.pri,
html[data-theme="light"] .buy-btn,
html[data-theme="dark"] .buy-btn{
  background:linear-gradient(135deg,#2563eb 0%,#8b5cf6 58%,#38bdf8 100%)!important;
  box-shadow:0 10px 24px rgba(59,130,246,.22)!important;
}

html[data-theme="light"] .theme-toggle,
html[data-theme="light"] .admin-mini,
html[data-theme="dark"] .theme-toggle,
html[data-theme="dark"] .admin-mini{
  background:rgba(8,20,39,.78)!important;
  color:#f8fafc!important;
  border-color:rgba(148,163,184,.22)!important;
  box-shadow:0 10px 24px rgba(0,0,0,.22)!important;
}

@media(max-width:760px){
  html[data-theme="light"] body,
  html[data-theme="dark"] body{
    background:
      radial-gradient(circle at 50% -8%, rgba(30,64,175,.30) 0%, rgba(30,64,175,.13) 18%, transparent 40%),
      radial-gradient(circle at 88% 10%, rgba(34,211,238,.14) 0%, transparent 22%),
      linear-gradient(180deg,#05263c 0%,#062840 38%,#0b2439 72%,#071a2b 100%)!important;
  }

  html[data-theme="light"] .card,
  html[data-theme="dark"] .card{
    background:rgba(8,20,39,.56)!important;
  }

  html[data-theme="light"] .hero,
  html[data-theme="dark"] .hero{
    background:linear-gradient(135deg,rgba(31,41,55,.94),rgba(67,56,202,.84) 52%,rgba(6,182,212,.78) 100%)!important;
  }
}


/* CATPEN_BLINK_FIX_V3 */
@keyframes catpenBlinkClear{
  0%, 12%, 18%, 45%, 51%, 100% {
    height:0%;
  }
  14.5%, 16% {
    height:58%;
  }
  47.5%, 49% {
    height:58%;
  }
}

.catpen-eye{
  overflow:hidden!important;
}

.catpen-eye::before,
.catpen-eye::after{
  content:""!important;
  position:absolute!important;
  left:50%!important;
  width:155%!important;
  background:var(--fur)!important;
  z-index:5!important;
  translate:-50%!important;
  animation:catpenBlinkClear 4.8s ease-in-out infinite!important;
}

.catpen-eye::before{
  top:-8%!important;
  height:0%;
  border-radius:0 0 50% 50%!important;
  rotate:-7deg!important;
}

.catpen-eye::after{
  bottom:-8%!important;
  height:0%;
  border-radius:50% 50% 0 0!important;
  rotate:-7deg!important;
}

.catpen-eye:nth-of-type(2)::before,
.catpen-eye:nth-of-type(2)::after{
  animation-delay:.13s!important;
}

@media(prefers-reduced-motion:reduce){
  .catpen-eye::before,
  .catpen-eye::after{
    animation:none!important;
    height:0%!important;
  }
}


/* PUBLIC_ICON_SPACING_V1 */
.public-actions{
  gap:10px!important;
}

.theme-toggle{
  width:52px!important;
  min-width:52px!important;
  padding:0 7px!important;
  justify-content:space-between!important;
}

.admin-mini{
  margin-left:2px!important;
}

@media(max-width:760px){
  .public-actions{
    top:8px!important;
    right:8px!important;
    gap:9px!important;
    padding:3px!important;
    border-radius:999px!important;
    background:rgba(255,255,255,.10)!important;
    backdrop-filter:blur(10px)!important;
  }

  .theme-toggle{
    width:48px!important;
    min-width:48px!important;
    height:28px!important;
    padding:0 7px!important;
    font-size:11px!important;
  }

  .admin-mini{
    width:28px!important;
    min-width:28px!important;
    height:28px!important;
    margin-left:1px!important;
    font-size:12px!important;
  }

  html[data-theme="dark"] .public-actions{
    background:rgba(2,6,23,.18)!important;
  }
}


/* FIX_THEME_GEAR_SPACING_V2 */
.public-actions{
  display:flex!important;
  align-items:center!important;
  justify-content:flex-end!important;
  gap:12px!important;
  background:transparent!important;
  padding:0!important;
}

.public-actions .theme-toggle{
  flex:0 0 auto!important;
  width:48px!important;
  min-width:48px!important;
  height:30px!important;
  margin:0!important;
}

.public-actions .admin-mini{
  flex:0 0 auto!important;
  width:30px!important;
  min-width:30px!important;
  height:30px!important;
  margin:0!important;
}

@media(max-width:760px){
  .public-actions{
    position:fixed!important;
    top:8px!important;
    right:8px!important;
    z-index:80!important;
    width:88px!important;
    height:30px!important;
    display:block!important;
    background:transparent!important;
    padding:0!important;
    border-radius:0!important;
    backdrop-filter:none!important;
  }

  .public-actions .theme-toggle{
    position:absolute!important;
    left:0!important;
    right:auto!important;
    top:0!important;
    width:48px!important;
    min-width:48px!important;
    height:30px!important;
    padding:0 7px!important;
    margin:0!important;
    justify-content:space-between!important;
  }

  .public-actions .admin-mini{
    position:absolute!important;
    right:0!important;
    left:auto!important;
    top:0!important;
    width:30px!important;
    min-width:30px!important;
    height:30px!important;
    padding:0!important;
    margin:0!important;
  }
}


/* PUBLIC_SINGLE_NIGHT_THEME_V1 */
.theme-toggle{
  display:none!important;
}

.public-actions{
  display:contents!important;
}

.admin-mini{
  width:30px!important;
  min-width:30px!important;
  height:30px!important;
  padding:0!important;
  border-radius:999px!important;
}

@media(max-width:760px){
  .admin-mini{
    position:fixed!important;
    top:8px!important;
    right:8px!important;
    z-index:80!important;
    width:30px!important;
    min-width:30px!important;
    height:30px!important;
    padding:0!important;
    font-size:12px!important;
  }
}


/* INVOICE_PUBLIC_FIX_V1 */
body:has(.invoice-wrap) .catpen-scene{
  display:none!important;
}

body:has(.invoice-wrap) .wrap{
  max-width:920px!important;
}

.invoice-wrap{
  max-width:860px!important;
  margin:0 auto!important;
}

.invoice-wrap > .card{
  background:rgba(8,20,39,.72)!important;
  border:1px solid rgba(148,163,184,.20)!important;
  box-shadow:0 18px 48px rgba(0,0,0,.28)!important;
  backdrop-filter:blur(12px)!important;
}

.invoice-wrap h1,
.invoice-wrap h2,
.invoice-wrap h3{
  color:#f8fafc!important;
}

.invoice-wrap .muted,
.invoice-wrap small{
  color:#cbd5e1!important;
}

.invoice-wrap .grid{
  gap:14px!important;
}

.invoice-wrap .grid .card{
  background:rgba(15,23,42,.68)!important;
  border-color:rgba(148,163,184,.18)!important;
  box-shadow:none!important;
}

.invoice-wrap .voucher-box{
  background:
    radial-gradient(circle at top right, rgba(56,189,248,.20), transparent 32%),
    linear-gradient(180deg,rgba(15,23,42,.96),rgba(2,6,23,.96))!important;
  border:2px dashed rgba(125,211,252,.72)!important;
  border-radius:24px!important;
  padding:28px!important;
  box-shadow:0 18px 42px rgba(0,0,0,.28)!important;
}

.invoice-wrap .voucher-box .muted{
  color:#bae6fd!important;
  font-weight:800!important;
  letter-spacing:.03em!important;
}

.invoice-wrap .voucher-code{
  color:#ffffff!important;
  text-shadow:0 0 18px rgba(56,189,248,.36)!important;
  font-size:64px!important;
  letter-spacing:.18em!important;
  margin:12px 0 8px!important;
}

.invoice-wrap .voucher-box p{
  color:#e2e8f0!important;
}

.invoice-wrap .voucher-box .btn{
  width:100%!important;
  max-width:620px!important;
  margin-left:auto!important;
  margin-right:auto!important;
}

.invoice-wrap .notice{
  background:rgba(15,23,42,.78)!important;
  border-color:rgba(148,163,184,.20)!important;
  color:#e5e7eb!important;
  font-size:22px!important;
  line-height:1.45!important;
}

.invoice-wrap .notice b{
  color:#ffffff!important;
}

@media(max-width:760px){
  body:has(.invoice-wrap) .wrap{
    padding:12px!important;
  }

  .invoice-wrap > .card{
    border-radius:20px!important;
    padding:16px!important;
  }

  .invoice-wrap h1{
    font-size:22px!important;
    line-height:1.15!important;
  }

  .invoice-wrap .grid{
    grid-template-columns:1fr!important;
    gap:10px!important;
  }

  .invoice-wrap .grid .card{
    padding:14px!important;
    border-radius:16px!important;
  }

  .invoice-wrap .voucher-box{
    padding:20px 14px!important;
    border-radius:20px!important;
  }

  .invoice-wrap .voucher-code{
    font-size:48px!important;
    letter-spacing:.14em!important;
  }

  .invoice-wrap .notice{
    font-size:18px!important;
    padding:14px!important;
    border-radius:16px!important;
  }
}

@media(max-width:420px){
  .invoice-wrap .voucher-code{
    font-size:42px!important;
    letter-spacing:.12em!important;
  }
}


/* INVOICE_LABEL_COPY_FIX_V1 */
/* VG_LUNAS_BADGE_GREEN_TRANSPARENT_V1 */
.invoice-status-lunas{
	background:rgba(22,163,74,.12)!important;
	color:#16a34a!important;
	border:1px solid rgba(22,163,74,.24)!important;
	box-shadow:none!important;
}

.invoice-wrap .voucher-box p:empty{
  display:none!important;
}

.invoice-wrap .voucher-box .btn.pri{
  font-size:18px!important;
}

@media(max-width:760px){
  /* VG_LUNAS_BADGE_GREEN_TRANSPARENT_V1 */
.invoice-status-lunas{
	background:rgba(22,163,74,.12)!important;
	color:#16a34a!important;
	border:1px solid rgba(22,163,74,.24)!important;
	box-shadow:none!important;
}

  .invoice-wrap .voucher-box .btn.pri{
    font-size:16px!important;
  }
}


/* HOTSPOT_LOGIN_REDIRECT_V1 */
.invoice-login-btn{
  width:100%!important;
  max-width:620px!important;
  margin:12px auto 0!important;
  background:rgba(15,23,42,.92)!important;
  color:#f8fafc!important;
  border:1px solid rgba(148,163,184,.24)!important;
  font-size:18px!important;
}

.invoice-login-btn:hover{
  filter:brightness(1.08)!important;
}

@media(max-width:760px){
  .invoice-login-btn{
    font-size:16px!important;
    margin-top:10px!important;
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


/* PUBLIC_DAY_NIGHT_SCENE_V1 */
.public-actions{
	display:flex!important;
	align-items:center!important;
	gap:8px!important;
}
.public-actions .theme-toggle,
.public-actions .admin-mini{
	position:static!important;
	width:42px!important;
	height:42px!important;
	min-width:42px!important;
	padding:0!important;
	border-radius:16px!important;
	display:inline-flex!important;
	align-items:center!important;
	justify-content:center!important;
	font-size:16px!important;
	z-index:60!important;
}
.public-actions .theme-toggle{
	border:1px solid rgba(148,163,184,.24)!important;
	background:rgba(15,23,42,.35)!important;
	color:#fff!important;
	box-shadow:0 14px 34px rgba(15,23,42,.18)!important;
	backdrop-filter:blur(10px)!important;
	cursor:pointer!important;
}
.public-actions .admin-mini{
	background:rgba(15,23,42,.35)!important;
	color:#fff!important;
	border:1px solid rgba(148,163,184,.24)!important;
	box-shadow:0 14px 34px rgba(15,23,42,.18)!important;
	backdrop-filter:blur(10px)!important;
}
html[data-theme="light"] .public-actions .theme-toggle,
html[data-theme="light"] .public-actions .admin-mini{
	background:rgba(255,255,255,.86)!important;
	color:#0f172a!important;
	border-color:rgba(15,23,42,.10)!important;
	box-shadow:0 12px 28px rgba(15,23,42,.10)!important;
}
.theme-toggle .sun,
.theme-toggle .moon{
	line-height:1!important;
}
html[data-theme="light"] .theme-toggle .moon{display:none!important}
html[data-theme="dark"] .theme-toggle .sun{display:none!important}

.public-day-scene{
	display:none;
	position:fixed;
	inset:0;
	z-index:-2;
	overflow:hidden;
	pointer-events:none;
	background:
		linear-gradient(180deg,#79d7ff 0%,#a7ecff 38%,#fff7cf 62%,#f6d38a 100%);
}
.public-day-scene::before{
	content:"";
	position:absolute;
	left:-10%;
	right:-10%;
	bottom:22%;
	height:110px;
	background:
		radial-gradient(120px 40px at 20% 60%,rgba(255,255,255,.9),transparent 70%),
		linear-gradient(180deg,#23b7d8,#0f8fb8);
	border-radius:50% 50% 0 0;
	opacity:.75;
}
.public-day-scene::after{
	content:"";
	position:absolute;
	left:-15%;
	right:-15%;
	bottom:0;
	height:28%;
	background:
		radial-gradient(circle at 18% 35%,rgba(255,255,255,.45) 0 2px,transparent 3px),
		radial-gradient(circle at 55% 55%,rgba(255,255,255,.40) 0 2px,transparent 3px),
		linear-gradient(180deg,#f8d68e,#eebf69);
}
.day-sun{
	position:absolute;
	top:72px;
	right:42px;
	width:78px;
	height:78px;
	border-radius:50%;
	background:#ffd34d;
	box-shadow:0 0 0 14px rgba(255,211,77,.18),0 0 60px rgba(255,180,0,.42);
}
.day-cloud{
	position:absolute;
	top:110px;
	left:34px;
	width:92px;
	height:34px;
	border-radius:999px;
	background:rgba(255,255,255,.78);
	box-shadow:36px -14px 0 8px rgba(255,255,255,.78),78px 3px 0 2px rgba(255,255,255,.78);
}
.day-cat{
	position:absolute;
	right:22px;
	bottom:52px;
	width:176px;
	height:142px;
}
.day-cat .umbrella{
	position:absolute;
	left:18px;
	top:-8px;
	width:150px;
	height:76px;
	border-radius:90px 90px 12px 12px;
	background:conic-gradient(from 180deg,#2563eb 0 25%,#fff 0 50%,#f97316 0 75%,#fff 0 100%);
	transform:rotate(-8deg);
	box-shadow:0 12px 25px rgba(15,23,42,.16);
}
.day-cat .cat-face{
	position:absolute;
	right:28px;
	bottom:26px;
	width:92px;
	height:78px;
	border-radius:45% 45% 42% 42%;
	background:#111827;
	box-shadow:0 18px 32px rgba(15,23,42,.20);
}
.day-cat .cat-face::before,
.day-cat .cat-face::after{
	content:"";
	position:absolute;
	top:-16px;
	width:28px;
	height:28px;
	background:#111827;
	clip-path:polygon(50% 0,0 100%,100% 100%);
}
.day-cat .cat-face::before{left:8px;transform:rotate(-18deg)}
.day-cat .cat-face::after{right:8px;transform:rotate(18deg)}
.day-cat .eye{
	position:absolute;
	top:26px;
	width:18px;
	height:18px;
	border-radius:50%;
	background:#f8fafc;
	box-shadow:inset -4px -4px 0 rgba(15,23,42,.16);
}
.day-cat .eye.left{left:24px}
.day-cat .eye.right{right:24px}
.day-cat .nose{
	position:absolute;
	top:48px;
	left:50%;
	width:12px;
	height:9px;
	border-radius:50%;
	background:#fb7185;
	transform:translateX(-50%);
}
.day-cat .mat{
	position:absolute;
	right:10px;
	bottom:0;
	width:150px;
	height:28px;
	border-radius:50%;
	background:rgba(37,99,235,.20);
}

html[data-theme="light"] .public-day-scene{display:block}
html[data-theme="dark"] .public-day-scene{display:none}
html[data-theme="light"] .catpen-scene,
html[data-theme="light"] .tuku-cat{
	display:none!important;
}
html[data-theme="dark"] .catpen-scene,
html[data-theme="dark"] .tuku-cat{
	display:block!important;
}
html[data-theme="light"] body{
	background:#eaf8ff!important;
}
html[data-theme="light"] .hero{
	background:linear-gradient(135deg,#0ea5e9,#22c55e)!important;
	color:#fff!important;
}
html[data-theme="light"] .card{
	background:rgba(255,255,255,.78)!important;
	backdrop-filter:blur(12px)!important;
	border-color:rgba(15,23,42,.10)!important;
}
html[data-theme="light"] .public-brand{
	color:#0f172a!important;
}
@media(max-width:760px){
	.public-actions{
		gap:7px!important;
	}
	.public-actions .theme-toggle,
	.public-actions .admin-mini{
		width:40px!important;
		height:40px!important;
		min-width:40px!important;
		border-radius:15px!important;
	}
	.day-sun{
		width:60px;
		height:60px;
		top:76px;
		right:30px;
	}
	.day-cloud{
		top:145px;
		left:28px;
		transform:scale(.82);
		transform-origin:left top;
	}
	.day-cat{
		right:12px;
		bottom:38px;
		transform:scale(.86);
		transform-origin:right bottom;
	}
}

/* PUBLIC_SUBTITLE_READABLE_V1 */
.hero .public-subtitle{
  margin:12px 0 0;
  max-width:820px;
  white-space:pre-line;
  line-height:1.48;
  font-size:clamp(15px,2.9vw,18px);
  font-weight:650;
  opacity:.92;
}
@media(max-width:640px){
  .hero .public-subtitle{
    line-height:1.42;
    font-size:15px;
    font-weight:650;
  }
}
/* TOGGLE_SLIDER_OVERRIDE_PUBLIC */
button.theme-toggle, button.theme-switch {
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

html[data-theme="dark"] button.theme-toggle,
html[data-theme="dark"] button.theme-switch {
    background: #2563eb !important;
}

button.theme-toggle::before, button.theme-switch::before {
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

html[data-theme="dark"] button.theme-toggle::before,
html[data-theme="dark"] button.theme-switch::before {
    transform: translateX(18px) !important;
}

button.theme-toggle .sun, button.theme-switch .sun {
    position: absolute !important;
    left: 6px !important;
    font-size: 11px !important;
    z-index: 1 !important;
    opacity: 1 !important;
    display: inline-block !important;
    transition: opacity 0.3s !important;
}

button.theme-toggle .moon, button.theme-switch .moon {
    position: absolute !important;
    right: 5px !important;
    font-size: 11px !important;
    z-index: 1 !important;
    opacity: 0 !important;
    display: inline-block !important;
    transition: opacity 0.3s !important;
}

html[data-theme="dark"] button.theme-toggle .sun,
html[data-theme="dark"] button.theme-switch .sun {
    opacity: 0 !important;
}

html[data-theme="dark"] button.theme-toggle .moon,
html[data-theme="dark"] button.theme-switch .moon {
    opacity: 1 !important;
}
/* /TOGGLE_SLIDER_OVERRIDE_PUBLIC */
`
