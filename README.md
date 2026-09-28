# VoucherGo

> Self-hosted WiFi voucher management for MikroTik.

VoucherGo is a lightweight web application for managing and selling MikroTik WiFi / hotspot vouchers with integrated QRIS payments and WhatsApp notifications.

## Features

- WiFi voucher sales
- Multi-router MikroTik support
- RouterOS REST API & Classic API
- QRIS payment via Pakasir
- Automatic voucher provisioning
- WhatsApp notification gateway
- Admin & reseller management
- Voucher package management
- Transaction & billing management
- Router and client management
- Responsive mobile-friendly UI
- Maintenance mode
- CSRF protection and configurable security settings

## Tech Stack

- Go
- MariaDB / MySQL
- MikroTik RouterOS
- Pakasir QRIS
- whatsapp-web.js
- Apache
- systemd

## Quick Start

### Requirements

- Go 1.20+
- MariaDB / MySQL
- MikroTik RouterOS
- Linux server
- Node.js

### Installation

Clone the repository:

    git clone https://github.com/sur-lysvara/Golang-VoucherGo.git
    cd Golang-VoucherGo

Build:

    go build -o vouchergo .

Run:

    ./vouchergo

> Production deployments may require additional configuration for the database, MikroTik, WhatsApp gateway, HTTPS, and systemd.


## Donations

If you find VoucherGo useful and would like to support the project, you can donate using:

| Asset | Network | Address |
|---|---|---|
| Bitcoin (BTC) | Bitcoin | `bc1qd9eumepe2nhx5tazcvqm5r43jgz7ny6c0mcz28` |
| Tether (USDT) | TRC20 | `TKgD73uLYDtEDYE6ZX6sDMhKGd8WtVpCm7` |