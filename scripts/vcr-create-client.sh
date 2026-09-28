#!/usr/bin/env bash
set -euo pipefail

if [[ "${EUID}" -ne 0 ]]; then
  echo "ERROR: jalankan pakai sudo/root"
  exit 1
fi

if [[ $# -lt 4 ]]; then
  echo "Usage:"
  echo "  sudo vcr-create-client.sh <slug> <domain> <app_port> <wa_port> [tag]"
  echo
  echo "Example:"
  echo "  sudo vcr-create-client.sh client-a client-a.vouchergo.biz.id 8111 8112 main"
  exit 1
fi

SLUG="$1"
DOMAIN="$2"

EXPECTED_NEW="${SLUG}.vouchergo.biz.id"
EXPECTED_OLD="${SLUG}.vcr.qithmir.my.id"

if [[ "$DOMAIN" != "$EXPECTED_NEW" && "$DOMAIN" != "$EXPECTED_OLD" ]]; then
  echo "ERROR: domain client harus format: <slug>.vouchergo.biz.id"
  echo "DEBUG: slug=[$SLUG] domain=[$DOMAIN] expected=[$EXPECTED_NEW]"
  exit 1
fi

APP_PORT="$3"
WA_PORT="$4"
TAG="${5:-main}"

if ! [[ "$SLUG" =~ ^[a-z0-9-]{3,20}$ ]]; then
  echo "ERROR: slug harus 3-20 karakter, huruf kecil/angka/dash. Contoh: client-a"
  exit 1
fi

if ! [[ "$APP_PORT" =~ ^[0-9]+$ ]] || ! [[ "$WA_PORT" =~ ^[0-9]+$ ]]; then
  echo "ERROR: port harus angka"
  exit 1
fi

if [[ "$APP_PORT" == "$WA_PORT" ]]; then
  echo "ERROR: app_port dan wa_port tidak boleh sama"
  exit 1
fi

NAME_UNDERSCORE="${SLUG//-/_}"

BASE_REPO="/var/www/vouchergo"
BASE_WA_TEMPLATE="${BASE_REPO}/scripts/wa-webjs"

ROOT_DIR="/var/www/vcr-clients/${SLUG}"
APP_DIR="${ROOT_DIR}/app"
WA_DIR="${ROOT_DIR}/wa"

DB_NAME="db_vcr_${NAME_UNDERSCORE}"
DB_USER="vcr_${NAME_UNDERSCORE}_u"
DB_PASS="$(openssl rand -hex 18)"

APP_SERVICE="tuku-${SLUG}"
WA_SERVICE="wa-webjs-${SLUG}"
ENV_FILE="/etc/tuku-${SLUG}.env"

WA_KEY="$(openssl rand -hex 18)"
SESSION_KEY="$(openssl rand -hex 32)"
ADMIN_USER="admin"
ADMIN_PASS="Vcr$(openssl rand -hex 4)!A1"

APACHE_CONF="/etc/apache2/sites-available/${DOMAIN}.conf"
INFO_FILE="/home/xyzsur/VoucherGo-CLIENT-${SLUG}-INFO.txt"

echo "=== VoucherGo CREATE CLIENT ==="
echo "Slug       : $SLUG"
echo "Domain     : $DOMAIN"
echo "App port   : $APP_PORT"
echo "WA port    : $WA_PORT"
echo "Tag        : $TAG"
echo "App dir    : $APP_DIR"
echo "WA dir     : $WA_DIR"
echo "DB         : $DB_NAME"
echo "App service: $APP_SERVICE"
echo "WA service : $WA_SERVICE"
echo

echo "=== preflight ==="
if [[ ! -d "$BASE_REPO/.git" ]]; then
  echo "ERROR: base repo tidak ditemukan: $BASE_REPO"
  exit 1
fi

if [[ ! -f "$BASE_WA_TEMPLATE/server.js" ]]; then
  echo "ERROR: template WA webjs tidak ditemukan: $BASE_WA_TEMPLATE"
  exit 1
fi

if git -C "$BASE_REPO" rev-parse "$TAG" >/dev/null 2>&1; then
  echo "OK: tag $TAG ditemukan"
else
  echo "ERROR: tag $TAG tidak ditemukan di $BASE_REPO"
  exit 1
fi

if ss -ltn | awk '{print $4}' | grep -qE "(:${APP_PORT}|:${WA_PORT})$"; then
  echo "ERROR: port $APP_PORT atau $WA_PORT sudah dipakai"
  ss -ltnp | grep -E ":${APP_PORT}|:${WA_PORT}" || true
  exit 1
fi

if [[ -e "$ROOT_DIR" || -e "$ENV_FILE" || -e "$APACHE_CONF" ]]; then
  echo "ERROR: folder/env/apache client sudah ada. Stop agar tidak menimpa."
  echo "$ROOT_DIR"
  echo "$ENV_FILE"
  echo "$APACHE_CONF"
  exit 1
fi

if systemctl list-unit-files | grep -q "^${APP_SERVICE}.service"; then
  echo "ERROR: service $APP_SERVICE sudah ada"
  exit 1
fi

if systemctl list-unit-files | grep -q "^${WA_SERVICE}.service"; then
  echo "ERROR: service $WA_SERVICE sudah ada"
  exit 1
fi

if mysql -N -e "SHOW DATABASES LIKE '${DB_NAME}';" | grep -q "$DB_NAME"; then
  echo "ERROR: database $DB_NAME sudah ada"
  exit 1
fi

echo "PREFLIGHT_OK"
echo

echo "=== clone app source ==="
mkdir -p "$ROOT_DIR"
git clone --branch "$TAG" --depth 1 "file://${BASE_REPO}" "$APP_DIR"

chown -R root:www-data "$APP_DIR"
chmod -R g+rX "$APP_DIR"

echo
echo "=== create database + user ==="
mysql <<SQL
CREATE DATABASE ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASS}';
GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
SQL

echo
echo "=== import schema from production ==="
mysqldump --no-data db_tuku | mysql "$DB_NAME"

echo
echo "=== seed client settings ==="
mysql "$DB_NAME" <<SQL
SET FOREIGN_KEY_CHECKS=0;
DELETE FROM transactions;
DELETE FROM packages;
DELETE FROM routers;
DELETE FROM users;
SET FOREIGN_KEY_CHECKS=1;

INSERT INTO settings(k,v) VALUES
('site_name','VoucherGo'),
('public_base_url','https://${DOMAIN}'),
('pakasir_slug','${SLUG}'),
('pakasir_api_key',''),
('pakasir_qris_only','1'),
('wa_provider','local'),
('wa_local_url','http://127.0.0.1:${WA_PORT}/send'),
('wa_local_key','${WA_KEY}'),
('wanesia_enabled','0'),
('wanesia_api_key',''),
('wanesia_sender',''),
('wanesia_base_url','disabled-client')
ON DUPLICATE KEY UPDATE v=VALUES(v);
SQL

echo
echo "=== create env ==="
cat > "$ENV_FILE" <<ENV
APP_ADDR=127.0.0.1:${APP_PORT}
DB_DSN="${DB_USER}:${DB_PASS}@unix(/var/run/mysqld/mysqld.sock)/${DB_NAME}?parseTime=true&loc=Local"
SESSION_KEY=${SESSION_KEY}
ADMIN_DEFAULT_USERNAME=${ADMIN_USER}
ADMIN_DEFAULT_PASSWORD=${ADMIN_PASS}
DEMO_MODE=0
ENV

chmod 600 "$ENV_FILE"
chown root:root "$ENV_FILE"

echo
echo "=== build app ==="
cd "$APP_DIR"
env "PATH=$PATH" go build -buildvcs=false -o tuku.new . || {
  echo "BUILD_FAILED"
  exit 1
}
install -m 755 tuku.new tuku
rm -f tuku.new

echo
echo "=== create app service ==="
cat > "/etc/systemd/system/${APP_SERVICE}.service" <<EOFAPP
[Unit]
Description=VoucherGo Client App ${SLUG}
After=network.target mariadb.service mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=${APP_DIR}
EnvironmentFile=${ENV_FILE}
ExecStart=${APP_DIR}/tuku
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOFAPP

echo
echo "=== clone WA service ==="
mkdir -p "$WA_DIR"
rsync -a \
  --exclude 'node_modules' \
  --exclude 'sessions' \
  --exclude '.cache' \
  --exclude '.config' \
  --exclude '.local' \
  "$BASE_WA_TEMPLATE/" "$WA_DIR/"

mkdir -p "$WA_DIR/sessions" "$WA_DIR/.cache/puppeteer" "$WA_DIR/.config" "$WA_DIR/.local"
cd "$WA_DIR"
npm install --omit=dev
chown -R www-data:www-data "$WA_DIR"

cat > "/etc/systemd/system/${WA_SERVICE}.service" <<EOFWA
[Unit]
Description=VoucherGo Client WhatsApp ${SLUG}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=${WA_DIR}
Environment=NODE_ENV=production
Environment=HOST=127.0.0.1
Environment=HOME=${WA_DIR}
Environment=PORT=${WA_PORT}
Environment=API_KEY=${WA_KEY}
Environment=AUTH_DIR=${WA_DIR}/sessions
Environment=PUPPETEER_CACHE_DIR=${WA_DIR}/.cache/puppeteer
ExecStart=/usr/bin/node ${WA_DIR}/server.js
Restart=always
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ReadWritePaths=${WA_DIR}

[Install]
WantedBy=multi-user.target
EOFWA

echo
echo "=== start services ==="
systemctl daemon-reload
systemctl enable "$APP_SERVICE" "$WA_SERVICE"
systemctl restart "$WA_SERVICE"
systemctl restart "$APP_SERVICE"
sleep 2

systemctl status "$APP_SERVICE" --no-pager -l | sed -n '1,14p'
echo
systemctl status "$WA_SERVICE" --no-pager -l | sed -n '1,14p'

echo
echo "=== create apache vhost HTTP ==="
a2enmod proxy proxy_http rewrite headers ssl >/dev/null

cat > "$APACHE_CONF" <<EOFAPACHE
<VirtualHost *:80>
    ServerName ${DOMAIN}

    ProxyPreserveHost On

    ProxyPass "/wa-qr" "http://127.0.0.1:${WA_PORT}/qr-page"
    ProxyPassReverse "/wa-qr" "http://127.0.0.1:${WA_PORT}/qr-page"

    ProxyPass / http://127.0.0.1:${APP_PORT}/
    ProxyPassReverse / http://127.0.0.1:${APP_PORT}/

    ErrorLog \${APACHE_LOG_DIR}/${DOMAIN}_error.log
    CustomLog \${APACHE_LOG_DIR}/${DOMAIN}_access.log combined
</VirtualHost>
EOFAPACHE

a2ensite "${DOMAIN}.conf"
apache2ctl configtest
systemctl reload apache2

echo
echo "=== local checks ==="
curl -fsSI "http://127.0.0.1:${APP_PORT}/admin" | sed -n '1,8p'
curl -fsSI -H "Host: ${DOMAIN}" "http://127.0.0.1/admin" | sed -n '1,8p'
curl -s -H "X-API-Key: ${WA_KEY}" "http://127.0.0.1:${WA_PORT}/status" || true
echo

echo
echo "=== save client info ==="
cat > "$INFO_FILE" <<EOFINFO
VoucherGo Client Info
===============

Client slug : ${SLUG}
Domain      : ${DOMAIN}
App URL     : http://${DOMAIN}/admin
App port    : ${APP_PORT}
WA port     : ${WA_PORT}

Admin login:
Username: ${ADMIN_USER}
Password: ${ADMIN_PASS}

Services:
App: ${APP_SERVICE}
WA : ${WA_SERVICE}

Folders:
App: ${APP_DIR}
WA : ${WA_DIR}

Database:
DB user : ${DB_USER}
DB name : ${DB_NAME}

WA:
Local status : http://127.0.0.1:${WA_PORT}/status
QR path      : http://${DOMAIN}/wa-qr
WA API key   : ${WA_KEY}

Next SSL command after DNS points to this server:
sudo certbot --apache -d ${DOMAIN}

Useful commands:
sudo systemctl status ${APP_SERVICE} --no-pager -l
sudo systemctl status ${WA_SERVICE} --no-pager -l
sudo systemctl restart ${APP_SERVICE}
sudo systemctl restart ${WA_SERVICE}
EOFINFO

chown xyzsur:xyzsur "$INFO_FILE" 2>/dev/null || true
chmod 600 "$INFO_FILE"

echo
echo "=== generate onboarding ==="
if [[ -x /usr/local/sbin/vcr-client-onboarding.sh ]]; then
  /usr/local/sbin/vcr-client-onboarding.sh "$SLUG" || true
else
  echo "WARN: /usr/local/sbin/vcr-client-onboarding.sh tidak ditemukan"
fi

echo
echo "CREATE_CLIENT_OK"
echo "Info file: $INFO_FILE"
echo
echo "Admin URL : http://${DOMAIN}/admin"
echo "Username  : ${ADMIN_USER}"
echo "Password  : ${ADMIN_PASS}"
echo
echo "Set DNS A record first, then run:"
echo "sudo certbot --apache -d ${DOMAIN}"


# VCR_CLIENT_BILLING_DEFAULTS_V1_START
# Tambahkan env billing default untuk client baru.
# Aman/idempotent: tidak menimpa nilai billing yang sudah ada.
billing_slug="${SLUG:-${CLIENT_SLUG:-${1:-}}}"

if [ -n "$billing_slug" ]; then
  billing_env="/etc/tuku-${billing_slug}.env"

  if [ -f "$billing_env" ]; then
    billing_due="${CLIENT_BILLING_DUE_AT:-$(date -d '+30 days' +%F)}"
    billing_grace="${CLIENT_BILLING_GRACE_UNTIL:-$(date -d "${billing_due} +3 days" +%F)}"
    billing_price="${CLIENT_BILLING_PRICE:-Rp150.000}"
    billing_owner_url="${CLIENT_BILLING_OWNER_URL:-https://vouchergo.biz.id/billing/pay?client=${billing_slug}}"

    python3 - "$billing_env" "$billing_slug" "$billing_due" "$billing_grace" "$billing_price" "$billing_owner_url" <<'PYBILL'
from pathlib import Path
import sys

envf = Path(sys.argv[1])
slug, due, grace, price, owner_url = sys.argv[2:7]

s = envf.read_text()
lines = s.splitlines()

defaults = {
    "CLIENT_SLUG": slug,
    "CLIENT_BILLING_STATUS": "active",
    "CLIENT_BILLING_DUE_AT": due,
    "CLIENT_BILLING_GRACE_UNTIL": grace,
    "CLIENT_BILLING_PRICE": price,
    "CLIENT_BILLING_OWNER_URL": owner_url,
}

existing = {}
for line in lines:
    raw = line.strip()
    if "=" in raw and not raw.startswith("#"):
        k, v = raw.split("=", 1)
        existing[k.strip()] = v.strip().strip('"').strip("'")

out = []
seen = set()

for line in lines:
    raw = line.strip()
    if "=" in raw and not raw.startswith("#"):
        k = raw.split("=", 1)[0].strip()
        if k in defaults:
            val = defaults[k]
            # CLIENT_SLUG dipaksa benar, field lain jangan timpa kalau sudah ada.
            if k != "CLIENT_SLUG" and existing.get(k):
                val = existing[k]
            out.append(f'{k}="{val}"')
            seen.add(k)
            continue
    out.append(line)

for k, val in defaults.items():
    if k not in seen:
        out.append(f'{k}="{val}"')

envf.write_text("\n".join(out).rstrip() + "\n")
PYBILL

    chown root:root "$billing_env"
    chmod 600 "$billing_env"

    if command -v setfacl >/dev/null 2>&1; then
      setfacl -m u:www-data:r "$billing_env" || true
    else
      chgrp www-data "$billing_env" || true
      chmod 640 "$billing_env" || true
    fi

    echo "OK: billing defaults applied to $billing_env"
  else
    echo "WARN: billing env target belum ada: $billing_env"
  fi
fi
# VCR_CLIENT_BILLING_DEFAULTS_V1_END
