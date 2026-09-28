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