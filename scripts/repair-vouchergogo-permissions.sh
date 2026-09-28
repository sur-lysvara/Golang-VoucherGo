#!/usr/bin/env bash
set -Eeuo pipefail
umask 027

APP_BRAND="VoucherGo"

OWNER_APP_DIR="${VG_OWNER_APP_DIR:-/var/www/vouchergo}"
CLIENTS_DIR="${VG_CLIENTS_DIR:-/var/www/vcr-clients}"
WA_DIR="${VG_WA_DIR:-/var/www/wa-webjs}"

APP_USER="${VG_APP_USER:-www-data}"
APP_GROUP="${VG_APP_GROUP:-www-data}"

log() { echo "== $* =="; }
warn() { echo "WARN: $*" >&2; }
die() { echo "ERROR: $*" >&2; exit 1; }

need_root() {
  [ "$(id -u)" = "0" ] || die "jalankan sebagai root: sudo bash $0"
}

check_user_group() {
  log "Cek user/group"
  id "$APP_USER" >/dev/null 2>&1 || die "user tidak ada: $APP_USER"
  getent group "$APP_GROUP" >/dev/null 2>&1 || die "group tidak ada: $APP_GROUP"
}

fix_tree_owner() {
  local dir="$1"

  if [ ! -d "$dir" ]; then
    warn "skip, folder tidak ada: $dir"
    return 0
  fi

  log "Fix owner app tree: $dir"

  find "$dir" \
    -path "$dir/.git" -prune -o \
    -path "$dir/node_modules" -prune -o \
    -exec chown "$APP_USER:$APP_GROUP" {} +

  find "$dir" \
    -path "$dir/.git" -prune -o \
    -path "$dir/node_modules" -prune -o \
    -type d -exec chmod 0750 {} +

  find "$dir" \
    -path "$dir/.git" -prune -o \
    -path "$dir/node_modules" -prune -o \
    -type f -exec chmod 0640 {} +

  find "$dir" \
    -path "$dir/.git" -prune -o \
    -path "$dir/node_modules" -prune -o \
    -type f \( -name 'tuku' -o -name '*.sh' \) -exec chmod 0750 {} +
}

fix_env_files() {
  log "Fix env files"
  shopt -s nullglob
  local files=(/etc/tuku*.env)

  if [ "${#files[@]}" -eq 0 ]; then
    warn "tidak ada /etc/tuku*.env"
    shopt -u nullglob
    return 0
  fi

  chown root:"$APP_GROUP" "${files[@]}"
  chmod 0640 "${files[@]}"
  ls -lh "${files[@]}"
  shopt -u nullglob
}

fix_systemd_files() {
  log "Fix systemd files"
  shopt -s nullglob
  local files=(/etc/systemd/system/tuku*.service)

  if [ "${#files[@]}" -eq 0 ]; then
    warn "tidak ada /etc/systemd/system/tuku*.service"
    shopt -u nullglob
    return 0
  fi

  chown root:root "${files[@]}"
  chmod 0644 "${files[@]}"
  systemctl daemon-reload
  ls -lh "${files[@]}"
  shopt -u nullglob
}

fix_sbin_helpers() {
  log "Fix helper scripts"
  shopt -s nullglob
  local files=(/usr/local/sbin/vcr-* /usr/local/sbin/tuku-* /usr/local/sbin/vouchergogo-*)

  if [ "${#files[@]}" -eq 0 ]; then
    warn "tidak ada helper vcr/tuku/vouchergogo di /usr/local/sbin"
    shopt -u nullglob
    return 0
  fi

  chown root:root "${files[@]}"
  chmod 0755 "${files[@]}"
  ls -lh "${files[@]}" | head -80
  shopt -u nullglob
}

fix_apache_files() {
  log "Fix Apache site files"
  if [ -d /etc/apache2/sites-available ]; then
    chown -R root:root /etc/apache2/sites-available
    find /etc/apache2/sites-available -type f -exec chmod 0644 {} +
  fi

  if [ -d /etc/apache2/sites-enabled ]; then
    chown -h root:root /etc/apache2/sites-enabled/* 2>/dev/null || true
  fi

  apache2ctl configtest || true
}

final_check() {
  log "Final check"
  echo
  echo "App dirs:"
  ls -ld "$OWNER_APP_DIR" "$CLIENTS_DIR" "$WA_DIR" 2>/dev/null || true

  echo
  echo "VoucherGo services:"
  systemctl --no-pager --full --type=service --all | grep -E 'tuku|wa-webjs|vcr|VoucherGo' || true

  echo
  echo "Repair permission selesai."
}

main() {
  need_root
  check_user_group
  fix_tree_owner "$OWNER_APP_DIR"
  fix_tree_owner "$CLIENTS_DIR"
  fix_tree_owner "$WA_DIR"
  fix_env_files
  fix_systemd_files
  fix_sbin_helpers
  fix_apache_files
  final_check
}

main "$@"
