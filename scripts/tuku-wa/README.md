# VoucherGo Local WhatsApp Service

Snapshot service WhatsApp lokal Baileys untuk VoucherGo.

Live path:
- /var/www/tuku-wa/server.js
- systemd: tuku-wa.service
- localhost: http://127.0.0.1:8098
- auth/session live ada di /var/www/tuku-wa/auth dan tidak ikut Git.

Catatan:
- QR scan adalah metode utama.
- Pairing code nomor HP hanya opsi kedua.
- Endpoint kompatibel: /status, /qr, /qr-page, /send, /send-image, /logout, /reconnect, /reset-session.
