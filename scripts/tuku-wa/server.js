import express from 'express'
import QRCode from 'qrcode'
import pino from 'pino'
import fs from 'fs'
import path from 'path'
import makeWASocket, {
  Browsers,
  DisconnectReason,
  fetchLatestBaileysVersion,
  useMultiFileAuthState
} from '@whiskeysockets/baileys'

const app = express()
app.use(express.json({ limit: '2mb' }))
app.use(express.urlencoded({ extended: true }))

const PORT = Number(process.env.PORT || process.env.WA_PORT || 8098)
const HOST = process.env.WA_HOST || process.env.HOST || '127.0.0.1'
const API_KEY = process.env.WA_API_KEY || process.env.API_KEY || process.env.TUKU_WA_API_KEY || 'tuku-local-wa-key'

function findAuthDir() {
  if (process.env.WA_AUTH_DIR) return process.env.WA_AUTH_DIR
  if (process.env.AUTH_DIR) return process.env.AUTH_DIR

  const candidates = [
    './auth_info_baileys',
    './auth',
    './session',
    './sessions',
    './baileys_auth'
  ]

  for (const d of candidates) {
    if (fs.existsSync(path.join(d, 'creds.json'))) return d
  }

  return './auth_info_baileys'
}

const AUTH_DIR = findAuthDir()
let sock = null
let starting = false
let lastQR = ''
let lastPairCode = ''
let lastPairPhone = ''
let connected = false
let lastState = 'starting'
let lastError = ''

function cleanPhone(v) {
  return String(v || '').replace(/[^0-9]/g, '')
}

function normalizeJid(v) {
  let s = String(v || '').trim()
  if (!s) return ''
  if (s.includes('@')) return s
  s = cleanPhone(s)
  return s ? `${s}@s.whatsapp.net` : ''
}

function requireKey(req, res, next) {
  const key = req.headers['x-api-key'] || req.query.key || req.body?.key
  if (String(key || '') !== API_KEY) {
    return res.status(401).json({ ok: false, error: 'unauthorized' })
  }
  next()
}


// WA_AUTO_RESET_401_V1
let autoResettingSession = false

async function resetSessionAndStart(reason = 'manual') {
  if (autoResettingSession) return
  autoResettingSession = true

  try {
    lastState = 'resetting-session'
    lastError = reason
    connected = false
    lastQR = ''
    lastPairCode = ''
    lastPairPhone = ''

    try {
      if (sock?.ws?.close) sock.ws.close()
    } catch (_) {}

    try {
      if (sock?.end) sock.end(new Error(reason))
    } catch (_) {}

    sock = null

    fs.rmSync(AUTH_DIR, { recursive: true, force: true })
    fs.mkdirSync(AUTH_DIR, { recursive: true, mode: 0o700 })

    setTimeout(() => {
      autoResettingSession = false
      startSock(true).catch(e => {
        lastError = e?.message || String(e)
      })
    }, 1500)
  } catch (e) {
    autoResettingSession = false
    lastError = e?.message || String(e)
  }
}


async function startSock(force = false) {
  if (starting) return sock
  if (sock && !force) return sock

  starting = true
  lastState = 'connecting'
  lastError = ''

  try {
    const { state, saveCreds } = await useMultiFileAuthState(AUTH_DIR)
    const { version } = await fetchLatestBaileysVersion()

    sock = makeWASocket({
      version,
      auth: state,
      logger: pino({ level: process.env.LOG_LEVEL || 'silent' }),
      browser: Browsers.ubuntu('VoucherGo WA'),
      printQRInTerminal: false,
      markOnlineOnConnect: false,
      syncFullHistory: false,
      generateHighQualityLinkPreview: false
    })

    sock.ev.on('creds.update', saveCreds)

    sock.ev.on('connection.update', async (u) => {
      if (u.qr) {
        lastQR = u.qr
        lastState = 'qr'
        connected = false
      }

      if (u.connection === 'open') {
        connected = true
        lastQR = ''
        lastPairCode = ''
        lastState = 'connected'
      }

      if (u.connection === 'close') {
        connected = false
        const code = u.lastDisconnect?.error?.output?.statusCode
        lastState = `closed${code ? ':' + code : ''}`
        lastError = u.lastDisconnect?.error?.message || ''

        const loggedOut = code === DisconnectReason.loggedOut || code === 401
        sock = null

        if (loggedOut) {
          resetSessionAndStart('logged out / session invalid / 401')
        } else {
          setTimeout(() => startSock(true).catch(() => {}), 3000)
        }
      }
    })

    return sock
  } finally {
    starting = false
  }
}

app.get('/status', async (req, res) => {
  res.json({
    ok: true,
    provider: 'baileys',
    connection: connected ? 'open' : lastState,
    connected,
    qr_available: !!lastQR,
    has_qr: !!lastQR,
    has_pair_code: !!lastPairCode,
    pair_phone: lastPairPhone,
    me: sock?.user || null,
    auth_dir: AUTH_DIR,
    last_error: lastError,
    error: lastError,
    state: lastState
  })
})

app.get('/check', async (req, res) => {
  res.json({ ok: true, connected, state: lastState })
})

app.get('/qr', async (req, res) => {
  if (!lastQR) {
    return res.status(404).json({ ok: false, error: 'QR belum tersedia / WhatsApp sudah connected' })
  }

  const png = await QRCode.toBuffer(lastQR, { type: 'png', margin: 2, width: 320 })
  res.setHeader('Content-Type', 'image/png')
  res.send(png)
})

app.get('/qr-page', async (req, res) => {
  const qrImg = lastQR ? await QRCode.toDataURL(lastQR, { margin: 2, width: 320 }) : ''
  res.setHeader('Content-Type', 'text/html; charset=utf-8')
  res.end(`<!doctype html>
<html lang="id">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>VoucherGo WhatsApp Login</title>
<style>
body{margin:0;font-family:system-ui,-apple-system,Segoe UI,sans-serif;background:#0f172a;color:#e5e7eb}
.wrap{max-width:760px;margin:auto;padding:18px}
.card{background:rgba(15,23,42,.9);border:1px solid rgba(148,163,184,.25);border-radius:18px;padding:16px;margin:12px 0;box-shadow:0 20px 50px rgba(0,0,0,.25)}
h2{margin:0 0 6px;font-size:20px}
p{color:#94a3b8;line-height:1.45}
input{width:100%;box-sizing:border-box;border:1px solid rgba(148,163,184,.35);border-radius:14px;background:#020617;color:#fff;padding:13px;font-size:16px;outline:none}
button{border:0;border-radius:14px;padding:12px 14px;font-weight:700;cursor:pointer;background:#22c55e;color:#052e16;margin-top:10px;width:100%}
button.secondary{background:#334155;color:#e5e7eb}
.code{font-size:30px;letter-spacing:4px;font-weight:900;text-align:center;background:#020617;border:1px dashed #64748b;border-radius:16px;padding:16px;margin-top:12px;color:#facc15}
.badge{display:inline-flex;padding:6px 10px;border-radius:999px;background:#1e293b;color:#cbd5e1;font-size:13px}
.qr{text-align:center}.qr img{background:white;border-radius:14px;padding:10px;max-width:100%;box-sizing:border-box}
.small{font-size:13px;color:#94a3b8}
</style>
</head>
<body>
<div class="wrap">
  <div class="card">
    <h2>Login WhatsApp VoucherGo</h2>
    <div class="badge">Status: ${connected ? 'Connected' : lastState}</div>
    <p>Metode utama: scan QR. Login pakai nomor HP tersedia sebagai opsi kedua jika QR sulit digunakan.</p>
  </div>

  <div class="card qr">
    <h2>Login pakai QR</h2>
    <p class="small">Buka WhatsApp di HP utama → Perangkat tertaut → Tautkan perangkat → scan QR ini.</p>
    ${qrImg ? `<img src="${qrImg}" alt="QR WhatsApp">` : `<p>QR belum tersedia atau WhatsApp sudah connected.</p>`}
    <button class="secondary" onclick="location.reload()">Refresh QR</button>
  </div>

  <div class="card">
    <h2>Opsi kedua: Login pakai Nomor HP</h2>
    <p>Gunakan ini kalau QR sulit discan. Isi nomor WhatsApp dengan kode negara. Contoh: <b>6281234567890</b></p>
    <input id="phone" placeholder="6281234567890" inputmode="numeric">
    <button onclick="pair()">Minta Kode Pairing</button>
    <div id="result"></div>
    <p class="small">Di HP utama: WhatsApp → Perangkat tertaut → Tautkan dengan nomor telepon → masukkan kode.</p>
  </div>
</div>
<script>
async function pair(){
  const phone = document.getElementById('phone').value.replace(/[^0-9]/g,'')
  const result = document.getElementById('result')
  result.innerHTML = '<p>Meminta kode...</p>'
  try{
    const r = await fetch('/pair-code', {
      method:'POST',
      headers:{'Content-Type':'application/json','X-API-Key':'${API_KEY}'},
      body:JSON.stringify({phone})
    })
    const j = await r.json()
    if(!j.ok) throw new Error(j.error || 'gagal')
    result.innerHTML = '<div class="code">'+j.code+'</div><p class="small">Masukkan kode ini di WhatsApp HP utama.</p>'
  }catch(e){
    result.innerHTML = '<p style="color:#fca5a5">Gagal: '+e.message+'</p>'
  }
}
</script>
</body>
</html>`)
})

app.post('/pair-code', requireKey, async (req, res) => {
  try {
    const phone = cleanPhone(req.body?.phone || req.body?.number || req.query.phone)

    if (!/^[1-9][0-9]{8,17}$/.test(phone)) {
      return res.status(400).json({
        ok: false,
        error: 'Nomor tidak valid. Gunakan format 62812xxxx tanpa +/spasi/strip.'
      })
    }

    const s = await startSock(false)

    if (connected || s?.authState?.creds?.registered) {
      return res.json({ ok: true, connected: true, message: 'WhatsApp sudah connected/registered' })
    }

    await new Promise(r => setTimeout(r, 1200))
    const code = await s.requestPairingCode(phone)

    lastPairCode = code
    lastPairPhone = phone
    lastState = 'pairing-code'

    res.json({ ok: true, phone, code })
  } catch (e) {
    lastError = e?.message || String(e)
    res.status(500).json({ ok: false, error: lastError })
  }
})


// WA_SEND_WAIT_READY_V1
async function waitWAReady(timeoutMs = 10000) {
  if (!sock) {
    try {
      await startSock(false)
    } catch (e) {
      lastError = e?.message || String(e)
    }
  }

  const end = Date.now() + timeoutMs
  while (Date.now() < end) {
    if (sock && connected) return sock
    await new Promise(r => setTimeout(r, 500))
  }

  if (sock && connected) return sock
  return null
}

app.post('/send', requireKey, async (req, res) => {
  try {
    const readySock = await waitWAReady(10000)
    if (!readySock) {
      return res.status(503).json({
        ok: false,
        error: 'WhatsApp belum connected',
        state: lastState,
        last_error: lastError
      })
    }

    const to = normalizeJid(req.body?.to || req.body?.number || req.body?.phone)
    const text = String(req.body?.message || req.body?.text || '').trim()

    if (!to || !text) return res.status(400).json({ ok: false, error: 'to/message wajib diisi' })

    const r = await readySock.sendMessage(to, { text })
    res.json({ ok: true, id: r?.key?.id || '' })
  } catch (e) {
    lastError = e?.message || String(e)
    res.status(500).json({ ok: false, error: lastError })
  }
})


// WA_SEND_IMAGE_COMPAT_V1
app.post('/send-image', requireKey, async (req, res) => {
  try {
    const readySock = await waitWAReady(10000)
    if (!readySock) {
      return res.status(503).json({
        ok: false,
        error: 'WhatsApp belum connected',
        state: lastState,
        last_error: lastError
      })
    }

    const to = normalizeJid(req.body?.to || req.body?.number || req.body?.phone)
    const imageUrl = String(req.body?.image_url || req.body?.imageUrl || req.body?.url || req.body?.image || '').trim()
    const caption = String(req.body?.caption || req.body?.message || req.body?.text || '').trim()

    if (!to) return res.status(400).json({ ok: false, error: 'to wajib diisi' })

    if (!imageUrl) {
      if (!caption) return res.status(400).json({ ok: false, error: 'image_url/message wajib diisi' })
      const r = await readySock.sendMessage(to, { text: caption })
      return res.json({ ok: true, jid: to, id: r?.key?.id || '', message_id: r?.key?.id || '', fallback: 'text' })
    }

    const imgResp = await fetch(imageUrl)
    if (!imgResp.ok) {
      return res.status(502).json({ ok: false, error: 'gagal ambil image', status: imgResp.status, url: imageUrl })
    }

    const mime = imgResp.headers.get('content-type') || 'image/png'
    const arr = await imgResp.arrayBuffer()
    const buf = Buffer.from(arr)

    const r = await readySock.sendMessage(to, {
      image: buf,
      mimetype: mime,
      caption
    })

    res.json({
      ok: true,
      jid: to,
      id: r?.key?.id || '',
      message_id: r?.key?.id || ''
    })
  } catch (e) {
    lastError = e?.message || String(e)
    res.status(500).json({ ok: false, error: lastError })
  }
})


app.post('/logout', requireKey, async (req, res) => {
  try {
    if (sock) await sock.logout()
  } catch (_) {}

  connected = false
  sock = null
  lastQR = ''
  lastPairCode = ''
  lastState = 'logged-out'

  res.json({ ok: true })
})

app.post('/reconnect', requireKey, async (req, res) => {
  sock = null
  connected = false
  await startSock(true)
  res.json({ ok: true, state: lastState })
})

app.get('/', (req, res) => res.redirect('/qr-page'))

app.listen(PORT, HOST, async () => {
  console.log(`VoucherGo WA service listening on http://${HOST}:${PORT}`)
  console.log(`Auth dir: ${AUTH_DIR}`)
  startSock().catch(e => console.error('startSock error:', e))
})
