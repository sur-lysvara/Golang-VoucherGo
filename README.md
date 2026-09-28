# VoucherGo

> Manajemen voucher WiFi MikroTik yang dapat di-hosting sendiri.

VoucherGo adalah aplikasi web ringan untuk mengelola dan menjual voucher WiFi / hotspot MikroTik, dengan integrasi pembayaran QRIS dan notifikasi WhatsApp.

## Fitur

- Penjualan voucher WiFi
- Dukungan multi-router MikroTik
- RouterOS REST API & Classic API
- Pembayaran QRIS melalui Pakasir
- Aktivasi voucher otomatis
- Gateway notifikasi WhatsApp
- Manajemen admin & reseller
- Manajemen paket voucher
- Manajemen transaksi & billing
- Manajemen router dan client
- Tampilan responsif untuk perangkat mobile
- Mode maintenance
- Perlindungan CSRF dan konfigurasi keamanan

## Tech Stack

- Go
- MariaDB / MySQL
- MikroTik RouterOS
- Pakasir Payment Gateway
- whatsapp-web.js
- Apache
- systemd

## Quick Start

### Persyaratan

- Go 1.20+
- MariaDB / MySQL
- MikroTik RouterOS
- Linux Server
- Node.js

### Instalasi

Clone repository:

    git clone https://github.com/sur-lysvara/Golang-VoucherGo.git
    cd Golang-VoucherGo

Build:

    go build -o vouchergo .

Run:

    ./vouchergo

> Deployment production mungkin memerlukan konfigurasi tambahan untuk database, MikroTik, WhatsApp gateway, HTTPS, dan systemd.

## Donasi

Jika VoucherGo bermanfaat dan Anda ingin mendukung pengembangan project ini, belikan saya segelas kopi:

| Aset | Network | Alamat |
|---|---|---|
| Bitcoin (BTC) | Bitcoin | `bc1qd9eumepe2nhx5tazcvqm5r43jgz7ny6c0mcz28` |
| Tether (USDT) | TRC20 | `TKgD73uLYDtEDYE6ZX6sDMhKGd8WtVpCm7` |
| Monero (XMR) | Monero | `89EMb1LGtZDY9PfuBkWJ7gH9ig45K1a3P4VbF6BdChx1izW55w3qRX8M9t9Q6x1JhLifoMgzcx6ryUTZHgkkqibyMspR8a8` |