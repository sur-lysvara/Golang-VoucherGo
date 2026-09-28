#!/usr/bin/env bash
set -Eeuo pipefail
umask 027

APP_BRAND="VoucherGo"
ARCHIVE="${1:-}"

RESTORE_ROOT="${VG_RESTORE_ROOT:-/root/vouchergogo-restore}"
STAMP="$(date +%F_%H%M%S)"
WORKDIR="${RESTORE_ROOT}/work-${STAMP}"
SAFETY_DIR="${RESTORE_ROOT}/safety-before-restore-${STAMP}"
LOG_FILE="${RESTORE_ROOT}/restore-${STAMP}.log"

OWNER_APP_DIR="${VG_OWNER_APP_DIR:-/var/www/vouchergo}"
CLIENTS_DIR="${VG_CLIENTS_DIR:-/var/www/vcr-clients}"
WA_DIR="${VG_WA_DIR:-/var/www/wa-webjs}"

APP_USER="${VG_APP_USER:-www-data}"
APP_GROUP="${VG_APP_GROUP:-www-data}"

RESTORE_DB="${VG_RESTORE_DB:-1}"
RESTORE_APACHE="${VG_RESTORE_APACHE:-1}"
RESTORE_SYSTEMD="${VG_RESTORE_SYSTEMD:-1}"
RESTORE_SBIN="${VG_RESTORE_SBIN:-1}"
REBUILD_APP="${VG_REBUILD_APP:-1}"
INSTALL_PACKAGES="${VG_INSTALL_PACKAGES:-1}"
START_SERVICES="${VG_START_SERVICES:-1}"

log() { echo "== $* =="; }
warn() { echo "WARN: $*" >&2; }
die() { echo "ERROR: $*" >&2; exit 1; }

need_root() {
  [ "$(id -u)" = "0" ] || die "jalankan sebagai root: sudo bash $0 backup.tar.gz"
}

usage() {
  cat <<USAGE
Usage:
  sudo bash scripts/restore-vouchergogo.sh /path/backup.tar.gz

Optional env:
  VG_RESTORE_DB=1|0
  VG_RESTORE_APACHE=1|0
  VG_RESTORE_SYSTEMD=1|0
  VG_RESTORE_SBIN=1|0
  VG_REBUILD_APP=1|0
  VG_INSTALL_PACKAGES=1|0
  VG_START_SERVICES=1|0

Default restore path:
  owner app : ${OWNER_APP_DIR}
  clients   : ${CLIENTS_DIR}
  wa        : ${WA_DIR}
USAGE
}

validate_args() {
  [ -n "$ARCHIVE" ] || { usage; die "file backup belum diisi"; }
  [ -f "$ARCHIVE" ] || die "file backup tidak ditemukan: $ARCHIVE"
  case "$ARCHIVE" in
    *.tar.gz|*.tgz) ;;
    *) die "format backup harus .tar.gz atau .tgz" ;;
  esac
}

start_logging() {
  mkdir -p "$RESTORE_ROOT"
  touch "$LOG_FILE"
  chmod 600 "$LOG_FILE"
  exec > >(tee -a "$LOG_FILE") 2>&1

  log "${APP_BRAND} restore"
  echo "Tanggal: $(date -Is)"
  echo "Archive: $ARCHIVE"
  echo "Workdir: $WORKDIR"
  echo "Safety backup: $SAFETY_DIR"
  echo "Log: $LOG_FILE"
}

install_packages() {
  if [ "$INSTALL_PACKAGES" != "1" ]; then
    log "Skip install package karena VG_INSTALL_PACKAGES=0"
    return 0
  fi

  log "Install dependency dasar"
  apt update
  DEBIAN_FRONTEND=noninteractive apt install -y \
    apache2 mariadb-server certbot python3-certbot-apache \
    git curl wget unzip rsync ca-certificates gnupg \
    build-essential python3 python3-pip golang-go nodejs npm \
    openssl gzip tar
}

extract_archive() {
  log "Extract archive"
  rm -rf "$WORKDIR"
  mkdir -p "$WORKDIR"
  tar -xzf "$ARCHIVE" -C "$WORKDIR"

  BACKUP_DIR="$(find "$WORKDIR" -mindepth 1 -maxdepth 1 -type d -name 'vouchergogo-backup-*' | head -1)"
  [ -n "${BACKUP_DIR:-}" ] || die "folder vouchergogo-backup-* tidak ditemukan di archive"

  echo "Backup dir: $BACKUP_DIR"

  if [ -f "$BACKUP_DIR/manifest.json" ]; then
    echo
    echo "Manifest:"
    cat "$BACKUP_DIR/manifest.json"
    echo
  fi
}

backup_existing_path() {
  local path="$1"

  if [ -e "$path" ]; then
    mkdir -p "$SAFETY_DIR"
    log "Safety backup existing: $path"
    rsync -a "$path" "$SAFETY_DIR/"
  fi
}

restore_dir_contents() {
  local src="$1"
  local dst="$2"

  [ -d "$src" ] || { warn "skip restore, source tidak ada: $src"; return 0; }

  backup_existing_path "$dst"
  mkdir -p "$dst"
  rsync -a --delete "$src"/ "$dst"/
}

restore_files_glob() {
  local src_dir="$1"
  local dst_dir="$2"
  local pattern="$3"

  [ -d "$src_dir" ] || { warn "skip restore files, source tidak ada: $src_dir"; return 0; }

  mkdir -p "$dst_dir"
  shopt -s nullglob
  local files=("$src_dir"/$pattern)
  if [ "${#files[@]}" -eq 0 ]; then
    warn "tidak ada file match: $src_dir/$pattern"
    shopt -u nullglob
    return 0
  fi

  for f in "${files[@]}"; do
    if [ -e "$dst_dir/$(basename "$f")" ]; then
      backup_existing_path "$dst_dir/$(basename "$f")"
    fi
    cp -a "$f" "$dst_dir/"
  done
  shopt -u nullglob
}

restore_app_files() {
  log "Restore app/client/WA files"
  restore_dir_contents "$BACKUP_DIR/app" "$OWNER_APP_DIR"
  restore_dir_contents "$BACKUP_DIR/clients" "$CLIENTS_DIR"
  restore_dir_contents "$BACKUP_DIR/wa" "$WA_DIR"

  chown -R "$APP_USER:$APP_GROUP" "$OWNER_APP_DIR" "$CLIENTS_DIR" "$WA_DIR" 2>/dev/null || true
}

restore_env_files() {
  log "Restore env files"
  restore_files_glob "$BACKUP_DIR/env" "/etc" "tuku*.env"

  shopt -s nullglob
  local files=(/etc/tuku*.env)
  if [ "${#files[@]}" -gt 0 ]; then
    chown root:"$APP_GROUP" "${files[@]}" 2>/dev/null || true
    chmod 640 "${files[@]}" 2>/dev/null || true
  fi
  shopt -u nullglob
}

restore_systemd() {
  if [ "$RESTORE_SYSTEMD" != "1" ]; then
    log "Skip systemd restore karena VG_RESTORE_SYSTEMD=0"
    return 0
  fi

  log "Restore systemd services"
  restore_files_glob "$BACKUP_DIR/systemd" "/etc/systemd/system" "tuku*.service"
  restore_files_glob "$BACKUP_DIR/systemd" "/etc/systemd/system" "wa-webjs*.service"
  systemctl daemon-reload
}

restore_apache() {
  if [ "$RESTORE_APACHE" != "1" ]; then
    log "Skip Apache restore karena VG_RESTORE_APACHE=0"
    return 0
  fi

  log "Restore Apache config"
  a2enmod proxy proxy_http headers ssl rewrite >/dev/null || true

  if [ -d "$BACKUP_DIR/apache/sites-available" ]; then
    backup_existing_path "/etc/apache2/sites-available"
    rsync -a "$BACKUP_DIR/apache/sites-available"/ /etc/apache2/sites-available/
  fi

  if [ -d "$BACKUP_DIR/apache/sites-enabled" ]; then
    backup_existing_path "/etc/apache2/sites-enabled"
    rsync -a "$BACKUP_DIR/apache/sites-enabled"/ /etc/apache2/sites-enabled/
  fi

  apache2ctl configtest
  systemctl reload apache2 || true
}

restore_sbin_helpers() {
  if [ "$RESTORE_SBIN" != "1" ]; then
    log "Skip helper script restore karena VG_RESTORE_SBIN=0"
    return 0
  fi

  log "Restore helper scripts"
  mkdir -p /usr/local/sbin
  restore_files_glob "$BACKUP_DIR/sbin" "/usr/local/sbin" "*"
  chmod +x /usr/local/sbin/vcr-* /usr/local/sbin/tuku-* /usr/local/sbin/vouchergogo-* 2>/dev/null || true
}

restore_mysql() {
  if [ "$RESTORE_DB" != "1" ]; then
    log "Skip database restore karena VG_RESTORE_DB=0"
    return 0
  fi

  log "Restore MariaDB database"
  systemctl enable --now mariadb

  if [ -f "$BACKUP_DIR/mysql/all-databases.sql.gz" ]; then
    gzip -dc "$BACKUP_DIR/mysql/all-databases.sql.gz" | mysql -uroot
  elif [ -f "$BACKUP_DIR/mysql/all-databases.sql" ]; then
    mysql -uroot < "$BACKUP_DIR/mysql/all-databases.sql"
  else
    warn "database dump tidak ditemukan"
  fi
}

install_node_deps() {
  log "Install Node.js deps untuk WA jika ada"
  if [ -f "$WA_DIR/package.json" ]; then
    cd "$WA_DIR"
    npm install --omit=dev
    chown -R "$APP_USER:$APP_GROUP" "$WA_DIR" 2>/dev/null || true
  else
    warn "package.json WA tidak ditemukan di $WA_DIR"
  fi
}

rebuild_app() {
  if [ "$REBUILD_APP" != "1" ]; then
    log "Skip rebuild app karena VG_REBUILD_APP=0"
    return 0
  fi

  log "Rebuild Golang owner app"
  if [ -f "$OWNER_APP_DIR/go.mod" ]; then
    cd "$OWNER_APP_DIR"
    GOCACHE=/tmp/go-build-cache GOMODCACHE=/tmp/go-mod-cache go build -o /tmp/tuku-restore-check .
    install -o "$APP_USER" -g "$APP_GROUP" -m 0755 /tmp/tuku-restore-check "$OWNER_APP_DIR/tuku"
    rm -f /tmp/tuku-restore-check
  else
    warn "go.mod tidak ditemukan di $OWNER_APP_DIR, skip rebuild"
  fi
}

restart_services() {
  if [ "$START_SERVICES" != "1" ]; then
    log "Skip start services karena VG_START_SERVICES=0"
    return 0
  fi

  log "Restart VoucherGo services"
  systemctl daemon-reload

  shopt -s nullglob
  local services=(/etc/systemd/system/tuku*.service /etc/systemd/system/wa-webjs*.service)
  if [ "${#services[@]}" -eq 0 ]; then
    warn "tidak ada tuku*.service atau wa-webjs*.service"
    shopt -u nullglob
    return 0
  fi

  for f in "${services[@]}"; do
    svc="$(basename "$f")"
    echo "Restart: $svc"
    systemctl enable "$svc" >/dev/null || true
    systemctl restart "$svc" || true
  done
  shopt -u nullglob
}

final_check() {
  log "Final check"
  systemctl --no-pager --full --type=service --all | grep -E 'tuku|wa-webjs|vcr|VoucherGo' || true

  echo
  echo "Safety backup existing file/folder:"
  echo "$SAFETY_DIR"

  echo
  echo "Log restore:"
  echo "$LOG_FILE"

  echo
  echo "Restore selesai."
}

main() {
  need_root
  validate_args
  start_logging
  install_packages
  extract_archive
  restore_app_files
  restore_env_files
  restore_systemd
  restore_apache
  restore_sbin_helpers
  restore_mysql
  install_node_deps
  rebuild_app
  restart_services
  final_check
}

main "$@"
