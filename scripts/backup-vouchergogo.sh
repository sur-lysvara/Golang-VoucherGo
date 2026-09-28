#!/usr/bin/env bash
set -Eeuo pipefail
umask 027

APP_BRAND="VoucherGo"
BACKUP_ROOT="${VG_BACKUP_ROOT:-/home/xyzsur/backups}"
STAMP="$(date +%F_%H%M%S)"
BACKUP_NAME="vouchergogo-backup-${STAMP}"
WORKDIR="${BACKUP_ROOT}/${BACKUP_NAME}"
ARCHIVE="${BACKUP_ROOT}/${BACKUP_NAME}.tar.gz"
LOG_FILE="${BACKUP_ROOT}/${BACKUP_NAME}.log"

OWNER_APP_DIR="${VG_OWNER_APP_DIR:-/var/www/vouchergo}"
CLIENTS_DIR="${VG_CLIENTS_DIR:-/var/www/vcr-clients}"
WA_DIR="${VG_WA_DIR:-/var/www/wa-webjs}"

log() { echo "== $* =="; }
warn() { echo "WARN: $*" >&2; }
die() { echo "ERROR: $*" >&2; exit 1; }

need_root() {
  [ "$(id -u)" = "0" ] || die "jalankan sebagai root: sudo bash $0"
}

start_logging() {
  mkdir -p "$BACKUP_ROOT"
  touch "$LOG_FILE"
  chmod 600 "$LOG_FILE"
  exec > >(tee -a "$LOG_FILE") 2>&1

  log "${APP_BRAND} backup"
  echo "Tanggal: $(date -Is)"
  echo "Backup root: $BACKUP_ROOT"
  echo "Workdir: $WORKDIR"
  echo "Archive: $ARCHIVE"
}

prepare_dirs() {
  log "Prepare folder backup"
  rm -rf "$WORKDIR"
  mkdir -p \
    "$WORKDIR/app" \
    "$WORKDIR/clients" \
    "$WORKDIR/wa" \
    "$WORKDIR/env" \
    "$WORKDIR/systemd" \
    "$WORKDIR/apache" \
    "$WORKDIR/sbin" \
    "$WORKDIR/mysql" \
    "$WORKDIR/meta"
}

copy_if_exists() {
  local src="$1"
  local dst="$2"

  if [ -e "$src" ]; then
    rsync -a \
      --exclude='.git' \
      --exclude='node_modules' \
      --exclude='.cache' \
      --exclude='.config' \
      --exclude='.local' \
      --exclude='*.bak' \
      --exclude='*.bak_*' \
      --exclude='tmp' \
      "$src" "$dst"
  else
    warn "skip, tidak ada: $src"
  fi
}

backup_app_files() {
  log "Backup app/client/WA files"
  copy_if_exists "$OWNER_APP_DIR/" "$WORKDIR/app/"
  copy_if_exists "$CLIENTS_DIR/" "$WORKDIR/clients/"
  copy_if_exists "$WA_DIR/" "$WORKDIR/wa/"
}

backup_env_files() {
  log "Backup env files"
  shopt -s nullglob
  local files=(/etc/tuku*.env)
  if [ "${#files[@]}" -gt 0 ]; then
    cp -a "${files[@]}" "$WORKDIR/env/"
  else
    warn "tidak ada /etc/tuku*.env"
  fi
  shopt -u nullglob
}

backup_systemd() {
  log "Backup systemd service"
  shopt -s nullglob
  local files=(/etc/systemd/system/tuku*.service /etc/systemd/system/wa-webjs*.service)
  if [ "${#files[@]}" -gt 0 ]; then
    cp -a "${files[@]}" "$WORKDIR/systemd/"
  else
    warn "tidak ada /etc/systemd/system/tuku*.service atau wa-webjs*.service"
  fi

  systemctl list-units 'tuku*' --all --no-pager > "$WORKDIR/meta/systemd-tuku-units.txt" || true
  systemctl list-units 'wa-webjs*' --all --no-pager > "$WORKDIR/meta/systemd-wa-webjs-units.txt" || true
  shopt -u nullglob
}

backup_apache() {
  log "Backup Apache config"
  mkdir -p "$WORKDIR/apache/sites-available" "$WORKDIR/apache/sites-enabled"

  if [ -d /etc/apache2/sites-available ]; then
    rsync -a /etc/apache2/sites-available/ "$WORKDIR/apache/sites-available/"
  fi

  if [ -d /etc/apache2/sites-enabled ]; then
    rsync -a /etc/apache2/sites-enabled/ "$WORKDIR/apache/sites-enabled/"
  fi

  apache2ctl -S > "$WORKDIR/meta/apache-vhosts.txt" 2>&1 || true
}

backup_sbin_helpers() {
  log "Backup helper scripts"
  shopt -s nullglob
  local files=(/usr/local/sbin/vcr-* /usr/local/sbin/tuku-* /usr/local/sbin/vouchergogo-*)
  if [ "${#files[@]}" -gt 0 ]; then
    cp -a "${files[@]}" "$WORKDIR/sbin/"
  else
    warn "tidak ada helper vcr/tuku/vouchergogo di /usr/local/sbin"
  fi
  shopt -u nullglob
}

backup_mysql() {
  log "Backup MariaDB all databases"
  systemctl is-active --quiet mariadb || warn "mariadb tidak aktif"

  if command -v mysqldump >/dev/null 2>&1; then
    mysqldump --all-databases --single-transaction --routines --triggers --events \
      > "$WORKDIR/mysql/all-databases.sql"
    gzip -9 "$WORKDIR/mysql/all-databases.sql"
  else
    warn "mysqldump tidak ditemukan, skip database"
  fi
}

write_manifest() {
  log "Write manifest"
  {
    echo "{"
    echo "  \"app\": \"VoucherGo\","
    echo "  \"backup_date\": \"$(date -Is)\","
    echo "  \"hostname\": \"$(hostname)\","
    echo "  \"owner_app_dir\": \"${OWNER_APP_DIR}\","
    echo "  \"clients_dir\": \"${CLIENTS_DIR}\","
    echo "  \"wa_dir\": \"${WA_DIR}\","
    echo "  \"archive\": \"${ARCHIVE}\""
    echo "}"
  } > "$WORKDIR/manifest.json"

  uname -a > "$WORKDIR/meta/uname.txt" || true
  df -h > "$WORKDIR/meta/df-h.txt" || true
  free -h > "$WORKDIR/meta/free-h.txt" || true
  ip addr > "$WORKDIR/meta/ip-addr.txt" || true
}

make_archive() {
  log "Create archive"
  tar -C "$BACKUP_ROOT" -czf "$ARCHIVE" "$BACKUP_NAME"
  chmod 600 "$ARCHIVE"

  log "Backup selesai"
  ls -lh "$ARCHIVE"
  echo "Log: $LOG_FILE"
}

main() {
  need_root
  start_logging
  prepare_dirs
  backup_app_files
  backup_env_files
  backup_systemd
  backup_apache
  backup_sbin_helpers
  backup_mysql
  write_manifest
  make_archive
}

main "$@"
