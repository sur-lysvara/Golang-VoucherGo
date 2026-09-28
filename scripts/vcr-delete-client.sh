#!/usr/bin/env bash
set -euo pipefail

if [[ "${EUID}" -ne 0 ]]; then
  echo "ERROR: jalankan pakai sudo/root"
  exit 1
fi

if [[ $# -lt 3 ]]; then
  echo "Usage:"
  echo "  sudo vcr-delete-client.sh <slug> <domain> --yes"
  echo
  echo "Example:"
  echo "  sudo vcr-delete-client.sh client-a client-a.vouchergo.biz.id --yes"
  exit 1
fi

SLUG="$1"
DOMAIN="$2"
CONFIRM="$3"

if [[ "$CONFIRM" != "--yes" ]]; then
  echo "ERROR: wajib pakai --yes agar tidak salah hapus"
  exit 1
fi

if ! [[ "$SLUG" =~ ^[a-z0-9-]{3,20}$ ]]; then
  echo "ERROR: slug tidak valid"
  exit 1
fi

NAME_UNDERSCORE="${SLUG//-/_}"

ROOT_DIR="/var/www/vcr-clients/${SLUG}"
DB_NAME="db_vcr_${NAME_UNDERSCORE}"
DB_USER="vcr_${NAME_UNDERSCORE}_u"

APP_SERVICE="tuku-${SLUG}"
WA_SERVICE="wa-webjs-${SLUG}"
LEGACY_WA_SERVICE="tuku-wa-${SLUG}"
ENV_FILE="/etc/tuku-${SLUG}.env"

APACHE_CONF="/etc/apache2/sites-available/${DOMAIN}.conf"
APACHE_SSL_CONF="/etc/apache2/sites-available/${DOMAIN}-le-ssl.conf"
INFO_FILE="/home/xyzsur/VCR-CLIENT-${SLUG}-INFO.txt"

BACKUP_DIR="/home/xyzsur/backup-vcr-clients/deleted-${SLUG}"
TS="$(date +%F_%H%M%S)"

echo "=== VCR DELETE CLIENT ==="
echo "Slug    : $SLUG"
echo "Domain  : $DOMAIN"
echo "Root    : $ROOT_DIR"
echo "DB      : $DB_NAME"
echo "DB user : $DB_USER"
echo

mkdir -p "$BACKUP_DIR"

echo "=== backup sebelum hapus ==="
if mysql -N -e "SHOW DATABASES LIKE '${DB_NAME}';" | grep -q "$DB_NAME"; then
  mysqldump "$DB_NAME" | gzip > "${BACKUP_DIR}/${DB_NAME}-before-delete-${TS}.sql.gz"
fi

if [[ -d "$ROOT_DIR" ]]; then
  tar -czf "${BACKUP_DIR}/${SLUG}-files-before-delete-${TS}.tar.gz" "$ROOT_DIR" 2>/dev/null || true
fi

for f in "$ENV_FILE" "/etc/systemd/system/${APP_SERVICE}.service" "/etc/systemd/system/${WA_SERVICE}.service" "/etc/systemd/system/${LEGACY_WA_SERVICE}.service" "$APACHE_CONF" "$APACHE_SSL_CONF" "$INFO_FILE"; do
  if [[ -e "$f" ]]; then
    cp -a "$f" "$BACKUP_DIR/" || true
  fi
done

echo
echo "=== stop services ==="
systemctl stop "$APP_SERVICE" 2>/dev/null || true
systemctl stop "$WA_SERVICE" 2>/dev/null || true
systemctl stop "$LEGACY_WA_SERVICE" 2>/dev/null || true
systemctl disable "$APP_SERVICE" 2>/dev/null || true
systemctl disable "$WA_SERVICE" 2>/dev/null || true
systemctl disable "$LEGACY_WA_SERVICE" 2>/dev/null || true

echo
echo "=== remove apache site ==="
a2dissite "${DOMAIN}.conf" 2>/dev/null || true
a2dissite "${DOMAIN}-le-ssl.conf" 2>/dev/null || true

rm -f "$APACHE_CONF" "$APACHE_SSL_CONF"

apache2ctl configtest
systemctl reload apache2

echo
echo "=== remove systemd/env/files ==="
rm -f "/etc/systemd/system/${APP_SERVICE}.service"
rm -f "/etc/systemd/system/${WA_SERVICE}.service"
rm -f "/etc/systemd/system/${LEGACY_WA_SERVICE}.service"
rm -f "$ENV_FILE"
rm -f "$INFO_FILE"
rm -rf "$ROOT_DIR"

systemctl daemon-reload
systemctl reset-failed "$APP_SERVICE" 2>/dev/null || true
systemctl reset-failed "$WA_SERVICE" 2>/dev/null || true
systemctl reset-failed "$LEGACY_WA_SERVICE" 2>/dev/null || true

echo
echo "=== drop db/user ==="
mysql <<SQL
DROP DATABASE IF EXISTS ${DB_NAME};
DROP USER IF EXISTS '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
SQL

echo
echo "=== remove certbot certificates ==="
for CERT_NAME in "$DOMAIN" "${SLUG}.vouchergo.biz.id" "${SLUG}.vcr.qithmir.my.id"; do
  if certbot certificates 2>/dev/null | grep -q "Certificate Name: ${CERT_NAME}$"; then
    echo "delete cert: $CERT_NAME"
    certbot delete --cert-name "$CERT_NAME" --non-interactive || true
  else
    echo "cert not found: $CERT_NAME"
  fi
done

echo
echo "=== remove onboarding/info files ==="
rm -f \
  "/home/xyzsur/VoucherGo-CLIENT-${SLUG}-INFO.txt" \
  "/home/xyzsur/VoucherGo-CLIENT-${SLUG}-ONBOARDING.txt" \
  "/home/xyzsur/VoucherGo-CLIENT-${SLUG}-OWNER-NOTES.txt" \
  "/home/xyzsur/VCR-CLIENT-${SLUG}-INFO.txt" \
  "/home/xyzsur/VCR-CLIENT-${SLUG}-ONBOARDING.txt" \
  "/home/xyzsur/VCR-CLIENT-${SLUG}-OWNER-NOTES.txt" || true

echo
echo "=== final check ==="
systemctl status "$APP_SERVICE" --no-pager -l 2>/dev/null | sed -n '1,8p' || true
systemctl status "$WA_SERVICE" --no-pager -l 2>/dev/null | sed -n '1,8p' || true
systemctl status "$LEGACY_WA_SERVICE" --no-pager -l 2>/dev/null | sed -n '1,8p' || true
mysql -N -e "SHOW DATABASES LIKE '${DB_NAME}';" || true
ls -ld "$ROOT_DIR" 2>/dev/null || true

echo
echo "DELETE_CLIENT_OK"
echo "Backup hapus ada di: $BACKUP_DIR"
