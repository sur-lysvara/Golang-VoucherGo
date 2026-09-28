# VoucherGo

VoucherGo adalah aplikasi web untuk penjualan voucher WiFi/hotspot MikroTik dengan pembayaran QRIS, provisioning otomatis ke router, dan notifikasi WhatsApp lokal.

Tagline: Jual Voucher WiFi Lebih Mudah

## Fitur Utama

- Penjualan voucher WiFi berbasis web
- Paket voucher per router MikroTik
- Integrasi MikroTik RouterOS REST dan Classic API
- Pembayaran QRIS via Pakasir
- Invoice publik dengan QRIS dan status pembayaran
- Provisioning voucher otomatis setelah pembayaran sukses
- WhatsApp lokal berbasis whatsapp-web.js
- Login WhatsApp via QR sebagai metode utama
- Reset session / QR baru dari panel admin
- Mode maintenance publik untuk menutup sementara halaman pembelian
- Admin transaksi, paket, router, user, settings, client manager, dan demo manager
- UI admin responsif untuk desktop dan mobile

## Stack

- Go
- MariaDB / MySQL
- Apache reverse proxy
- systemd
- MikroTik RouterOS REST / Classic API
- Pakasir QRIS
- whatsapp-web.js local WhatsApp gateway

## Struktur Penting

- /var/www/vouchergo : aplikasi utama VoucherGo
- /var/www/wa-webjs : service WhatsApp lokal whatsapp-web.js
- /etc/tuku.env : environment owner app
- /etc/tuku-<client>.env : environment client app
- /usr/local/sbin/vcr-*.sh : helper script client/demo/backup

## Service Utama

- tuku.service : owner app
- wa-webjs.service : WhatsApp owner
- tuku-demo.service : demo app
- wa-webjs-demo.service : WhatsApp demo
- tuku-<client>.service : client app
- wa-webjs-<client>.service : WhatsApp client

## WhatsApp Lokal

Service WhatsApp berjalan di localhost:

http://127.0.0.1:8666

Endpoint kompatibel:

- /status
- /qr
- /qr-page
- /send
- /send-image
- /logout
- /reconnect
- /reset-session

Catatan:

- QR scan adalah metode login utama.
- Pairing code tidak dipakai di gateway whatsapp-web.js saat ini.
- Folder session WhatsApp tidak ikut Git.
- Jika session invalid/logout, QR baru bisa dibuat dari panel admin.

## Build

Jalankan:

    go build -o /tmp/tuku-check .

## Deploy Aman

Contoh deploy owner app:

    go build -o /tmp/tuku-new .
    BIN="$(systemctl cat tuku | awk -F= '/^ExecStart=/{print $2; exit}' | awk '{print $1}')"
    cp -a "$BIN" "$BIN.bak_$(date +%F_%H%M%S)"
    install -o www-data -g www-data -m 755 /tmp/tuku-new "$BIN"
    systemctl restart tuku

## Release Terbaru

v1.0-rc15

## Catatan Keamanan

- Jangan commit file .env
- Jangan commit folder auth WhatsApp
- Jangan commit binary hasil build
- Jangan commit backup .bak
- Gunakan SESSION_KEY unik di environment production
- Webhook Pakasir wajib memakai verifikasi API key
