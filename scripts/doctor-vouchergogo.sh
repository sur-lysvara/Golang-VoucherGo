#!/usr/bin/env bash
set -Eeuo pipefail

OWNER_APP_DIR="${VG_OWNER_APP_DIR:-/var/www/vouchergo}"
CLIENTS_DIR="${VG_CLIENTS_DIR:-/var/www/vcr-clients}"
WA_DIR="${VG_WA_DIR:-/var/www/wa-webjs}"
APP_PORT="${VG_APP_PORT:-8097}"
WA_PORT="${VG_WA_PORT:-8666}"
WA_SERVICE="${VG_WA_SERVICE:-wa-webjs.service}"
WA_API_KEY="${VG_WA_API_KEY:-}"

ok() { echo "OK   - $*"; }
warn() { echo "WARN - $*"; }
fail() { echo "FAIL - $*"; }

section() {
  echo
  echo "== $* =="
}

check_cmd() {
  if command -v "$1" >/dev/null 2>&1; then
    ok "command ada: $1"
  else
    warn "command tidak ada: $1"
  fi
}

check_service() {
  local svc="$1"
  if systemctl list-unit-files "$svc" >/dev/null 2>&1; then
    if systemctl is-active --quiet "$svc"; then
      ok "service running: $svc"
    else
      warn "service tidak running: $svc"
      systemctl status "$svc" --no-pager -n 10 || true
    fi
  else
    warn "service tidak ditemukan: $svc"
  fi
}

section "VoucherGo Doctor"
echo "Tanggal : $(date -Is)"
echo "Host    : $(hostname)"
echo "User    : $(whoami)"

section "OS"
if [ -r /etc/os-release ]; then
  . /etc/os-release
  echo "${PRETTY_NAME:-unknown}"
else
  warn "/etc/os-release tidak ditemukan"
fi

section "Resource"
free -h || true
df -h / /var /home 2>/dev/null || df -h || true

section "Command Check"
for c in git go mysql mysqldump apache2ctl certbot node npm curl ss tar gzip rsync; do
  check_cmd "$c"
done

section "Folder Check"
for d in "$OWNER_APP_DIR" "$CLIENTS_DIR" "$WA_DIR"; do
  if [ -d "$d" ]; then
    ok "folder ada: $d"
    ls -ld "$d"
  else
    warn "folder tidak ada: $d"
  fi
done

section "Env Check"
shopt -s nullglob
envs=(/etc/tuku*.env)
if [ "${#envs[@]}" -gt 0 ]; then
  ok "env ditemukan: ${#envs[@]} file"
  ls -lh "${envs[@]}"
else
  warn "tidak ada /etc/tuku*.env"
fi
shopt -u nullglob

section "Systemd Services"
systemctl list-units 'tuku*' --all --no-pager || true
systemctl list-units 'wa-webjs*' --all --no-pager || true

shopt -s nullglob
svcs=(/etc/systemd/system/tuku*.service /etc/systemd/system/wa-webjs*.service)
if [ "${#svcs[@]}" -gt 0 ]; then
  for f in "${svcs[@]}"; do
    check_service "$(basename "$f")"
  done
else
  warn "tidak ada /etc/systemd/system/tuku*.service atau wa-webjs*.service"
fi
shopt -u nullglob

section "Apache"
if command -v apache2ctl >/dev/null 2>&1; then
  apache2ctl configtest || true
  apache2ctl -S 2>&1 | head -80 || true
else
  warn "apache2ctl tidak ditemukan"
fi

section "MariaDB"
check_service mariadb.service
if command -v mysql >/dev/null 2>&1; then
  mysql -uroot -e "SHOW DATABASES;" 2>/dev/null | head -40 || warn "mysql root tanpa password tidak bisa akses atau MariaDB belum siap"
fi

section "Port Check"
if command -v ss >/dev/null 2>&1; then
  ss -ltnp | grep -E ":(${APP_PORT}|${WA_PORT})\b" || warn "port ${APP_PORT}/${WA_PORT} belum listen"
else
  warn "ss tidak ditemukan"
fi

section "Local HTTP Check"
curl -I --max-time 8 "http://127.0.0.1:${APP_PORT}" || warn "app lokal port ${APP_PORT} belum respon"

section "WA Local Check"
if [ -z "$WA_API_KEY" ] && systemctl cat "$WA_SERVICE" >/dev/null 2>&1; then
  WA_API_KEY="$(systemctl cat "$WA_SERVICE" | awk -F= '/Environment=API_KEY=/ {print $3; exit}')"
fi
if [ -n "$WA_API_KEY" ]; then
  code="$(curl -sS --max-time 8 -o /tmp/vg-wa-status.json -w '%{http_code}' -H "X-API-Key: ${WA_API_KEY}" "http://127.0.0.1:${WA_PORT}/status" || true)"
  if [ "$code" = "200" ]; then
    ok "WA service port ${WA_PORT} respon dengan auth"
  else
    warn "WA service port ${WA_PORT} respon tidak normal: HTTP ${code:-curl-error}"
  fi
  rm -f /tmp/vg-wa-status.json
else
  code="$(curl -sS --max-time 8 -o /tmp/vg-wa-status.json -w '%{http_code}' "http://127.0.0.1:${WA_PORT}/status" || true)"
  if [ "$code" = "401" ]; then
    ok "WA service port ${WA_PORT} respon dan butuh auth"
  else
    warn "WA service port ${WA_PORT} belum respon normal: HTTP ${code:-curl-error}"
  fi
  rm -f /tmp/vg-wa-status.json
fi
echo

section "Git Status"
if [ -d "$OWNER_APP_DIR/.git" ]; then
  git -C "$OWNER_APP_DIR" status --short || true
  git -C "$OWNER_APP_DIR" log --oneline -5 || true
else
  warn "folder owner app bukan git repo atau .git tidak ikut restore"
fi

section "Done"
echo "Doctor selesai."
