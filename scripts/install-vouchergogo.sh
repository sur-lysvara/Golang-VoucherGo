#!/usr/bin/env bash
set -Eeuo pipefail
umask 027

APP_BRAND="VoucherGo"
APP_SERVICE="${VG_SERVICE:-tuku}"
WA_SERVICE="${VG_WA_SERVICE:-wa-webjs}"
APP_USER="${VG_APP_USER:-www-data}"
APP_GROUP="${VG_APP_GROUP:-www-data}"
APP_DIR="${VG_APP_DIR:-/var/www/vouchergo}"
WA_DIR="${VG_WA_DIR:-/var/www/wa-webjs}"
APP_PORT="${VG_APP_PORT:-8097}"
WA_PORT="${VG_WA_PORT:-8666}"
ENV_FILE="${VG_ENV_FILE:-/etc/tuku.env}"
DB_NAME="${VG_DB_NAME:-vouchergogo}"
DB_USER="${VG_DB_USER:-vouchergogo}"
DB_PASS="${VG_DB_PASS:-$(openssl rand -base64 32 | tr -d '=+/ ' | cut -c1-24)}"
SESSION_KEY="${VG_SESSION_KEY:-$(openssl rand -hex 32)}"
WA_API_KEY="${VG_WA_API_KEY:-local-wa-key}"
DOMAIN="${VG_DOMAIN:-}"
EMAIL="${VG_EMAIL:-}"
SKIP_SSL="${VG_SKIP_SSL:-0}"
STRICT_DNS="${VG_STRICT_DNS:-0}"
LOG_FILE="${VG_LOG_FILE:-/root/vouchergogo-install.log}"
INSTALL_VERSION="v0.2"
# VOUCHERGO_INSTALLER_SAFETY_V02

log() { echo "== $* =="; }
warn() { echo "WARN: $*" >&2; }
die() { echo "ERROR: $*" >&2; exit 1; }

start_logging() {
  mkdir -p "$(dirname "$LOG_FILE")"
  touch "$LOG_FILE"
  chmod 600 "$LOG_FILE"
  exec > >(tee -a "$LOG_FILE") 2>&1

  echo
  log "${APP_BRAND} installer ${INSTALL_VERSION}"
  echo "Tanggal: $(date -Is)"
  echo "Log: ${LOG_FILE}"
}

preflight_os() {
  log "Preflight OS"
  command -v apt >/dev/null 2>&1 || die "installer ini butuh OS berbasis apt/Ubuntu"

  if [ -r /etc/os-release ]; then
    . /etc/os-release
    echo "OS: ${PRETTY_NAME:-unknown}"
    if [ "${ID:-}" != "ubuntu" ]; then
      warn "OS bukan Ubuntu. Installer tetap lanjut, tapi rekomendasi: Ubuntu 22.04/24.04"
    fi
  fi
}

check_port_free() {
  log "Cek port aplikasi dan WhatsApp"
  if command -v ss >/dev/null 2>&1; then
    if ss -ltnH | awk '{print $4}' | grep -Eq "(^|:|\])(${APP_PORT}|${WA_PORT})$"; then
      ss -ltnp || true
      die "port ${APP_PORT}/${WA_PORT} sudah dipakai. Set port lain dengan VG_APP_PORT/VG_WA_PORT"
    fi
  else
    warn "command ss tidak ditemukan, skip cek port"
  fi
}

check_domain_dns() {
  log "Cek DNS domain"
  RESOLVED_IPS="$(getent ahostsv4 "$DOMAIN" 2>/dev/null | awk '{print $1}' | sort -u | tr '\n' ' ' || true)"
  SERVER_IP="$(curl -4 -fsS --max-time 8 https://api.ipify.org 2>/dev/null || true)"

  echo "Domain: ${DOMAIN}"
  echo "DNS IP: ${RESOLVED_IPS:-tidak terdeteksi}"
  echo "Server public IP: ${SERVER_IP:-tidak terdeteksi}"

  if [ -z "$RESOLVED_IPS" ]; then
    if [ "$STRICT_DNS" = "1" ]; then
      die "DNS domain belum resolve. Matikan strict dengan VG_STRICT_DNS=0"
    fi
    warn "DNS domain belum resolve. SSL bisa gagal kalau domain belum mengarah ke server."
    return 0
  fi

  if [ -n "$SERVER_IP" ] && ! printf '%s\n' "$RESOLVED_IPS" | grep -qw "$SERVER_IP"; then
    if [ "$STRICT_DNS" = "1" ]; then
      die "DNS domain tidak mengarah ke IP server"
    fi
    warn "DNS domain belum mengarah ke public IP server. SSL bisa gagal."
  fi
}

need_root() {
  [ "$(id -u)" = "0" ] || die "jalankan sebagai root: sudo bash $0"
}

ask_input() {
  if [ -z "$DOMAIN" ]; then
    read -r -p "Domain VoucherGo, contoh vcr.domain.com: " DOMAIN
  fi
  [ -n "$DOMAIN" ] || die "domain kosong"

  if [ -z "$EMAIL" ]; then
    read -r -p "Email SSL Let's Encrypt, boleh kosong: " EMAIL || true
  fi
}

install_packages() {
  log "Install dependency Ubuntu"
  apt update
  DEBIAN_FRONTEND=noninteractive apt install -y \
    apache2 mariadb-server certbot python3-certbot-apache \
    git curl wget unzip rsync ca-certificates gnupg \
    build-essential python3 python3-pip golang-go nodejs npm \
    openssl ufw
}

setup_database() {
  log "Setup MariaDB database"
  systemctl enable --now mariadb

  mysql -uroot <<SQL
CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASS}';
ALTER USER '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASS}';
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
SQL
}

copy_source() {
  log "Copy source VoucherGo"
  mkdir -p "$APP_DIR"

  SRC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
  [ -f "$SRC_DIR/go.mod" ] || die "go.mod tidak ditemukan di source: $SRC_DIR"
  find "$SRC_DIR" -maxdepth 1 -name '*.go' | grep -q . || die "file .go tidak ditemukan di source: $SRC_DIR"

  rsync -a --delete \
    --exclude='.git' \
    --exclude='tuku' \
    --exclude='*.bak' \
    --exclude='*.bak_*' \
    --exclude='node_modules' \
    --exclude='tmp' \
    "$SRC_DIR"/ "$APP_DIR"/

  chown -R "$APP_USER:$APP_GROUP" "$APP_DIR"
}

setup_env() {
  log "Setup env file"
  if [ -f "$ENV_FILE" ]; then
    cp -a "$ENV_FILE" "$ENV_FILE.bak_$(date +%F_%H%M%S)"
  fi

  cat > "$ENV_FILE" <<ENV
APP_NAME=VoucherGo
APP_ADDR=127.0.0.1:${APP_PORT}
DB_DSN=${DB_USER}:${DB_PASS}@tcp(127.0.0.1:3306)/${DB_NAME}?parseTime=true
SESSION_KEY=${SESSION_KEY}
PUBLIC_BASE_URL=https://${DOMAIN}
WA_BASE_URL=http://127.0.0.1:${WA_PORT}/send
WA_API_KEY=${WA_API_KEY}
ENV

  chown root:"$APP_GROUP" "$ENV_FILE"
  chmod 640 "$ENV_FILE"

  cat > /root/vouchergogo-credentials.txt <<CREDS
VoucherGo Install Credentials
Date: $(date -Is)
Domain: ${DOMAIN}
App dir: ${APP_DIR}
Service: ${APP_SERVICE}
Port: ${APP_PORT}
Env: ${ENV_FILE}
WA dir: ${WA_DIR}
WA service: ${WA_SERVICE}
WA port: ${WA_PORT}
WA API key: ${WA_API_KEY}

Database:
DB_NAME=${DB_NAME}
DB_USER=${DB_USER}
DB_PASS=${DB_PASS}

Application secret:
SESSION_KEY=${SESSION_KEY}
CREDS
  chmod 600 /root/vouchergogo-credentials.txt
}

build_app() {
  log "Build Golang app"
  cd "$APP_DIR"
  sudo -u "$APP_USER" go build -o /tmp/tuku-install-check .
  install -o "$APP_USER" -g "$APP_GROUP" -m 0755 /tmp/tuku-install-check "$APP_DIR/tuku"
  rm -f /tmp/tuku-install-check
}

setup_wa_gateway() {
  log "Setup WhatsApp gateway wa-webjs"
  local template_dir="${APP_DIR}/scripts/wa-webjs"
  [ -f "$template_dir/server.js" ] || die "template wa-webjs tidak ditemukan: $template_dir"

  mkdir -p "$WA_DIR"
  rsync -a --delete \
    --exclude='node_modules' \
    --exclude='sessions' \
    --exclude='.cache' \
    --exclude='.config' \
    --exclude='.local' \
    "$template_dir"/ "$WA_DIR"/

  mkdir -p "$WA_DIR/sessions" "$WA_DIR/.cache/puppeteer" "$WA_DIR/.config" "$WA_DIR/.local"
  cd "$WA_DIR"
  npm install --omit=dev
  chown -R "$APP_USER:$APP_GROUP" "$WA_DIR"

  cat > "/etc/systemd/system/${WA_SERVICE}.service" <<SERVICE
[Unit]
Description=Local WhatsApp Web JS Gateway
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${APP_USER}
Group=${APP_GROUP}
WorkingDirectory=${WA_DIR}
Environment=NODE_ENV=production
Environment=HOST=127.0.0.1
Environment=PORT=${WA_PORT}
Environment=API_KEY=${WA_API_KEY}
Environment=AUTH_DIR=${WA_DIR}/sessions
Environment=PUPPETEER_CACHE_DIR=${WA_DIR}/.cache/puppeteer
Environment=HOME=${WA_DIR}
ExecStart=/usr/bin/node ${WA_DIR}/server.js
Restart=always
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ReadWritePaths=${WA_DIR}

[Install]
WantedBy=multi-user.target
SERVICE

  systemctl daemon-reload
  systemctl enable --now "$WA_SERVICE"
}

setup_systemd() {
  log "Setup systemd service"
  cat > "/etc/systemd/system/${APP_SERVICE}.service" <<SERVICE
[Unit]
Description=VoucherGo App
After=network.target mariadb.service

[Service]
Type=simple
WorkingDirectory=${APP_DIR}
EnvironmentFile=${ENV_FILE}
ExecStart=${APP_DIR}/tuku
Restart=always
RestartSec=3
User=${APP_USER}
Group=${APP_GROUP}

[Install]
WantedBy=multi-user.target
SERVICE

  systemctl daemon-reload
  systemctl enable --now "$APP_SERVICE"
}

seed_wa_settings() {
  log "Seed setting WhatsApp lokal"
  sleep 2
  mysql -uroot "$DB_NAME" <<SQL
INSERT INTO settings(k,v) VALUES
('wa_provider','local'),
('wa_local_url','http://127.0.0.1:${WA_PORT}/send'),
('wa_local_key','${WA_API_KEY}')
ON DUPLICATE KEY UPDATE v=VALUES(v);
SQL
}

setup_apache() {
  log "Setup Apache reverse proxy"
  a2enmod proxy proxy_http headers ssl rewrite >/dev/null

  cat > "/etc/apache2/sites-available/${DOMAIN}.conf" <<APACHE
<VirtualHost *:80>
    ServerName ${DOMAIN}

    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:${APP_PORT}/
    ProxyPassReverse / http://127.0.0.1:${APP_PORT}/

    ErrorLog \${APACHE_LOG_DIR}/${DOMAIN}-error.log
    CustomLog \${APACHE_LOG_DIR}/${DOMAIN}-access.log combined
</VirtualHost>
APACHE

  a2ensite "${DOMAIN}.conf" >/dev/null
  apache2ctl configtest
  systemctl reload apache2
}

setup_ssl() {
  if [ "$SKIP_SSL" = "1" ]; then
    log "Skip SSL karena VG_SKIP_SSL=1"
    return 0
  fi

  log "Setup SSL Let's Encrypt"
  if [ -n "$EMAIL" ]; then
    certbot --apache -d "$DOMAIN" --non-interactive --agree-tos -m "$EMAIL" --redirect
  else
    certbot --apache -d "$DOMAIN" --register-unsafely-without-email --non-interactive --agree-tos --redirect
  fi
}

final_check() {
  log "Final check"
  systemctl restart "$APP_SERVICE"
  systemctl restart "$WA_SERVICE" || true
  sleep 2

  systemctl --no-pager --full status "$APP_SERVICE" || true
  systemctl --no-pager --full status "$WA_SERVICE" || true

  echo
  echo "Cek lokal:"
  curl -I "http://127.0.0.1:${APP_PORT}" || true
  curl -sS -o /dev/null -w "WA HTTP %{http_code}\n" -H "X-API-Key: ${WA_API_KEY}" "http://127.0.0.1:${WA_PORT}/status" || true

  echo
  echo "Selesai."
  echo "Credential tersimpan di: /root/vouchergogo-credentials.txt"
}

main() {
  need_root
  start_logging
  preflight_os
  ask_input
  install_packages
  check_port_free
  check_domain_dns
  setup_database
  copy_source
  setup_env
  build_app
  setup_wa_gateway
  setup_systemd
  seed_wa_settings
  setup_apache
  setup_ssl
  final_check
}

main "$@"
