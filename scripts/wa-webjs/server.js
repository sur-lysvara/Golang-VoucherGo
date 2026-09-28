const express = require("express");
const QRCode = require("qrcode");
const path = require("path");
const { Client, LocalAuth, MessageMedia } = require("whatsapp-web.js");

const HOST = process.env.HOST || "127.0.0.1";
const PORT = parseInt(process.env.PORT || "8666", 10);
const API_KEY = process.env.API_KEY || "local-wa-key";
const AUTH_DIR = process.env.AUTH_DIR || path.join(__dirname, "sessions");

let client = null;
let connection = "starting";
let connected = false;
let qrText = "";
let qrDataUrl = "";
let lastError = "";
let me = null;
let restarting = false;

function auth(req, res, next) {
  const key = req.get("X-API-Key") || req.query.key || req.query.api_key || "";
  if (key !== API_KEY) return res.status(401).json({ ok: false, error: "unauthorized" });
  next();
}

function normalizeNumber(input) {
  let n = String(input || "").trim();
  n = n.replace(/@c\.us$/i, "");
  n = n.replace(/[^\d]/g, "");
  if (!n) return "";
  if (n.startsWith("0")) n = "62" + n.slice(1);
  if (!n.startsWith("62") && n.length >= 8) n = "62" + n;
  return n + "@c.us";
}

function maskChatId(chatId) {
  const s = String(chatId || "");
  return s.replace(/^(\d{4})\d+(\d{3}@c\.us)$/i, "$1***$2");
}

function safeImageHint(imageUrl) {
  try {
    const u = new URL(String(imageUrl || ""));
    return `${u.protocol}//${u.host}${u.pathname}`.slice(0, 180);
  } catch (err) {
    return "[invalid-url]";
  }
}

function statusPayload() {
  return {
    ok: true,
    provider: "whatsapp-web.js",
    connection,
    state: connection,
    connected,
    qr_available: !!qrText,
    has_qr: !!qrText,
    has_pair_code: false,
    pair_phone: "",
    me,
    auth_dir: AUTH_DIR,
    last_error: lastError,
    error: lastError
  };
}

async function startClient() {
  lastError = "";
  connection = "starting";
  connected = false;
  me = null;
  qrText = "";
  qrDataUrl = "";

  client = new Client({
    authStrategy: new LocalAuth({
      clientId: "local-wa-webjs",
      dataPath: AUTH_DIR
    }),
    webVersionCache: { type: "none" },
    puppeteer: {
      headless: true,
      args: [
        "--no-sandbox",
        "--disable-setuid-sandbox",
        "--disable-dev-shm-usage",
        "--disable-gpu",
        "--disable-extensions",
        "--disable-background-networking",
        "--disable-sync",
        "--disable-default-apps",
        "--no-first-run",
        "--no-zygote"
      ]
    }
  });

  client.on("qr", async (qr) => {
    connection = "qr";
    connected = false;
    qrText = qr;
    qrDataUrl = await QRCode.toDataURL(qr);
    console.log("[WEBJS] QR updated");
  });

  client.on("authenticated", () => {
    connection = "authenticated";
    lastError = "";
    console.log("[WEBJS] authenticated");
  });

  client.on("ready", () => {
    connection = "connected";
    connected = true;
    qrText = "";
    qrDataUrl = "";
    lastError = "";
    me = client.info || null;
    console.log("[WEBJS] ready");
  });

  client.on("auth_failure", (msg) => {
    connection = "auth_failure";
    connected = false;
    lastError = String(msg || "auth_failure");
    console.error("[WEBJS] auth_failure", lastError);
  });

  client.on("disconnected", (reason) => {
    connection = "disconnected";
    connected = false;
    lastError = String(reason || "disconnected");
    console.error("[WEBJS] disconnected", lastError);
  });

  try {
    await client.initialize();
  } catch (err) {
    connection = "error";
    connected = false;
    lastError = err && err.message ? err.message : String(err);
    console.error("[WEBJS] initialize error", lastError);
  }
}

async function restartClient() {
  if (restarting) return;
  restarting = true;
  try {
    connection = "restarting";
    connected = false;
    qrText = "";
    qrDataUrl = "";
    if (client) {
      try { await client.destroy(); } catch (e) {}
    }
    client = null;
    await startClient();
  } finally {
    restarting = false;
  }
}

const app = express();
app.use(express.json({ limit: "1mb" }));
app.use(express.urlencoded({ extended: true }));

app.get("/", (req, res) => {
  res.json({ ok: true, service: "Local WhatsApp Web JS Gateway", port: PORT });
});

app.get("/status", auth, (req, res) => res.json(statusPayload()));

app.get("/check", auth, (req, res) => {
  if (!connected) return res.status(503).json({ ok: false, error: "not_connected", state: connection });
  res.json({ ok: true, connected: true, provider: "whatsapp-web.js", me });
});

app.get("/qr", auth, (req, res) => {
  // WEBJS_QR_DATAURL_COMPAT_V1
  // VoucherGo lama biasanya render field "qr" langsung ke <img src="">.
  // Jadi "qr" harus berupa data:image/png;base64, bukan teks QR mentah.
  res.json({
    ok: true,
    state: connection,
    connection,
    qr: qrDataUrl,
    qr_text: qrText,
    qr_raw: qrText,
    qr_data_url: qrDataUrl,
    qrDataUrl,
    data_url: qrDataUrl,
    image: qrDataUrl
  });
});

app.get("/qr-page", (req, res) => {
  const img = qrDataUrl ? `<img src="${qrDataUrl}" style="width:320px;max-width:90vw;border-radius:18px;background:#fff;padding:14px">` : "";
  res.send(`<!doctype html>
<html><head><meta name="viewport" content="width=device-width,initial-scale=1">
<meta http-equiv="refresh" content="8">
<title>WA Web JS QR</title>
<style>
body{margin:0;min-height:100vh;display:grid;place-items:center;background:#0f172a;color:#fff;font-family:Arial,sans-serif}
.card{max-width:520px;text-align:center;padding:26px}
.badge{display:inline-block;background:#1e293b;border:1px solid #334155;border-radius:999px;padding:8px 14px;margin:12px 0;color:#cbd5e1}
pre{white-space:pre-wrap;color:#fecaca;text-align:left}
</style></head>
<body><div class="card">
<h2>WA Web JS Gateway</h2>
<div class="badge">Status: ${connection}</div>
<div>${img || "<p>QR belum tersedia. Refresh otomatis...</p>"}</div>
${lastError ? `<pre>${lastError}</pre>` : ""}
</div></body></html>`);
});

app.post("/send", auth, async (req, res) => {
  try {
    if (!connected || !client) {
      return res.status(503).json({ ok: false, error: "not_connected", state: connection });
    }

    const to = req.body.to || req.body.phone || req.body.number || req.query.to || req.query.phone || "";
    const message = req.body.message || req.body.text || req.query.message || req.query.text || "";
    const chatId = normalizeNumber(to);

    if (!chatId) return res.status(400).json({ ok: false, error: "missing_to" });
    if (!String(message).trim()) return res.status(400).json({ ok: false, error: "missing_message" });

    const result = await client.sendMessage(chatId, String(message));
    res.json({
      ok: true,
      provider: "whatsapp-web.js",
      to: chatId,
      id: result && result.id ? result.id._serialized : ""
    });
  } catch (err) {
    lastError = err && err.message ? err.message : String(err);
    res.status(500).json({ ok: false, error: lastError, state: connection });
  }
});

// WEBJS_SEND_IMAGE_COMPAT_V1
// Kompatibel dengan VoucherGo auto notif QRIS yang POST ke /send-image.
app.post("/send-image", auth, async (req, res) => {
  let startedAt = 0;

  try {
    if (!connected || !client) {
      return res.status(503).json({ ok: false, error: "not_connected", state: connection });
    }

    const to = req.body.to || req.body.phone || req.body.number || req.query.to || req.query.phone || "";
    const caption = req.body.caption || req.body.message || req.body.text || req.query.caption || req.query.message || req.query.text || "";
    const imageUrl =
      req.body.image_url || req.body.imageUrl || req.body.image ||
      req.body.url || req.body.media_url || req.body.media ||
      req.query.image_url || req.query.imageUrl || req.query.image ||
      req.query.url || req.query.media_url || req.query.media || "";

    const chatId = normalizeNumber(to);
    if (!chatId) return res.status(400).json({ ok: false, error: "missing_to" });
    if (!String(imageUrl).trim()) return res.status(400).json({ ok: false, error: "missing_image_url" });

    startedAt = Date.now();
    console.log("[WEBJS] send-image start", {
      to: maskChatId(chatId),
      image: safeImageHint(imageUrl)
    });

    const media = await MessageMedia.fromUrl(String(imageUrl), { unsafeMime: true });
    const result = await client.sendMessage(chatId, media, { caption: String(caption || "") });

    console.log("[WEBJS] send-image ok", {
      to: maskChatId(chatId),
      duration_ms: Date.now() - startedAt
    });

    res.json({
      ok: true,
      provider: "whatsapp-web.js",
      to: chatId,
      image_url: imageUrl,
      id: result && result.id ? result.id._serialized : ""
    });
  } catch (err) {
    lastError = err && err.message ? err.message : String(err);
    console.error("[WEBJS] send-image error", {
      duration_ms: startedAt ? Date.now() - startedAt : null,
      error: lastError
    });
    res.status(500).json({ ok: false, error: lastError, state: connection });
  }
});


// WEBJS_ROOT_SEND_COMPAT_V1
// Kompatibel untuk handler lama VoucherGo yang POST ke base URL "/".
app.post("/", auth, async (req, res) => {
  try {
    if (!connected || !client) {
      return res.status(503).json({ ok: false, error: "not_connected", state: connection });
    }

    const to =
      req.body.to || req.body.phone || req.body.number ||
      req.body.receiver || req.body.recipient || req.body.destination ||
      req.query.to || req.query.phone || req.query.number ||
      req.query.receiver || req.query.recipient || req.query.destination || "";

    const message =
      req.body.message || req.body.text || req.body.msg || req.body.body || req.body.caption ||
      req.query.message || req.query.text || req.query.msg || req.query.body || req.query.caption || "";

    const imageUrl =
      req.body.image_url || req.body.imageUrl || req.body.image ||
      req.body.url || req.body.media_url || req.body.media ||
      req.query.image_url || req.query.imageUrl || req.query.image ||
      req.query.url || req.query.media_url || req.query.media || "";

    const chatId = normalizeNumber(to);
    if (!chatId) return res.status(400).json({ ok: false, error: "missing_to" });

    let result;
    if (String(imageUrl).trim()) {
      const media = await MessageMedia.fromUrl(String(imageUrl), { unsafeMime: true });
      result = await client.sendMessage(chatId, media, { caption: String(message || "") });
    } else {
      if (!String(message).trim()) return res.status(400).json({ ok: false, error: "missing_message" });
      result = await client.sendMessage(chatId, String(message));
    }

    res.json({
      ok: true,
      provider: "whatsapp-web.js",
      to: chatId,
      id: result && result.id ? result.id._serialized : ""
    });
  } catch (err) {
    lastError = err && err.message ? err.message : String(err);
    res.status(500).json({ ok: false, error: lastError, state: connection });
  }
});

app.post("/logout", auth, async (req, res) => {
  try {
    if (client) await client.logout();
    connection = "logged_out";
    connected = false;
    qrText = "";
    qrDataUrl = "";
    res.json({ ok: true });
  } catch (err) {
    lastError = err && err.message ? err.message : String(err);
    res.status(500).json({ ok: false, error: lastError });
  }
});

app.post("/reconnect", auth, (req, res) => {
  restartClient().catch(e => console.error("[WEBJS] reconnect error", e));
  res.json({ ok: true, state: "restarting" });
});

app.get("/reconnect", auth, (req, res) => {
  restartClient().catch(e => console.error("[WEBJS] reconnect error", e));
  res.json({ ok: true, state: "restarting" });
});

app.listen(PORT, HOST, () => {
  console.log(`WA Web JS listening on http://${HOST}:${PORT}`);
  console.log(`Auth dir: ${AUTH_DIR}`);
  startClient();
});
