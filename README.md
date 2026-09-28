cd ~/vouchergo-open && python3 -c 'from pathlib import Path; Path("README.md").write_text("""# VoucherGo

> Self-hosted WiFi voucher management for MikroTik.

VoucherGo is a lightweight web application for managing and selling MikroTik WiFi / hotspot vouchers with integrated QRIS payments and WhatsApp notifications.

## ✨ Features

- 🎟️ WiFi voucher sales
- 🌐 Multi-router MikroTik support
- 🔌 RouterOS REST API & Classic API
- 💳 QRIS payment via Pakasir
- ⚡ Automatic voucher provisioning
- 📱 WhatsApp notification gateway
- 👤 Admin & reseller management
- 📦 Voucher package management
- 💰 Transaction & billing management
- 📊 Router and client management
- 📱 Responsive mobile-friendly UI
- 🛠️ Maintenance mode
- 🔐 CSRF protection and configurable security settings

## 🧱 Tech Stack

- **Go**
- **MariaDB / MySQL**
- **MikroTik RouterOS**
- **Pakasir QRIS**
- **whatsapp-web.js**
- **Apache**
- **systemd**

## 🚀 Getting Started

### Requirements

- Go 1.20+
- MariaDB / MySQL
- MikroTik RouterOS
- Linux server
- Node.js

### Clone

```bash
git clone https://github.com/sur-lysvara/Golang-VoucherGo.git
cd Golang-VoucherGo
