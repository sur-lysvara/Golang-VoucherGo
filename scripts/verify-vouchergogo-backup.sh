#!/usr/bin/env bash
set -Eeuo pipefail

ARCHIVE="${1:-}"

ok() { echo "OK   - $*"; }
warn() { echo "WARN - $*"; }
fail() { echo "FAIL - $*"; FAILED=1; }

usage() {
  echo "Usage:"
  echo "  bash scripts/verify-vouchergogo-backup.sh /path/vouchergogo-backup.tar.gz"
}

check_contains() {
  local label="$1"
  local pattern="$2"

  if grep -Eq "$pattern" "$LIST_FILE"; then
    ok "$label"
  else
    fail "$label tidak ditemukan"
  fi
}

FAILED=0

[ -n "$ARCHIVE" ] || { usage; exit 1; }
[ -f "$ARCHIVE" ] || { echo "ERROR: file tidak ditemukan: $ARCHIVE" >&2; exit 1; }

case "$ARCHIVE" in
  *.tar.gz|*.tgz) ;;
  *) echo "ERROR: file harus .tar.gz atau .tgz" >&2; exit 1 ;;
esac

echo "== VoucherGo Backup Verify =="
echo "Archive: $ARCHIVE"
ls -lh "$ARCHIVE"

TMPDIR="$(mktemp -d)"
LIST_FILE="$TMPDIR/list.txt"
trap 'rm -rf "$TMPDIR"' EXIT

echo
echo "== Test tar integrity =="
tar -tzf "$ARCHIVE" > "$LIST_FILE"
ok "archive bisa dibaca"

ROOT_DIR="$(head -1 "$LIST_FILE" | cut -d/ -f1)"
echo "Root dir: ${ROOT_DIR:-unknown}"

echo
echo "== Cek isi penting =="
check_contains "manifest.json" '^vouchergogo-backup-.*/manifest\.json$'
check_contains "database dump" '^vouchergogo-backup-.*/mysql/all-databases\.sql(\.gz)?$'
check_contains "env files" '^vouchergogo-backup-.*/env/tuku.*\.env$'
check_contains "owner app files" '^vouchergogo-backup-.*/app/'
check_contains "Go module/source" '^vouchergogo-backup-.*/app/go\.mod$|^vouchergogo-backup-.*/app/.*\.go$'
check_contains "client directory" '^vouchergogo-backup-.*/clients/'
check_contains "WA directory" '^vouchergogo-backup-.*/wa/'
check_contains "systemd service" '^vouchergogo-backup-.*/systemd/(tuku|wa-webjs).*\.service$'
check_contains "apache config" '^vouchergogo-backup-.*/apache/sites-available/'
check_contains "helper scripts" '^vouchergogo-backup-.*/sbin/'

echo
echo "== Ringkasan jumlah file =="
echo "Total item archive : $(wc -l < "$LIST_FILE")"
echo "Env files          : $(grep -Ec '^vouchergogo-backup-.*/env/tuku.*\.env$' "$LIST_FILE" || true)"
echo "Systemd services   : $(grep -Ec '^vouchergogo-backup-.*/systemd/(tuku|wa-webjs).*\.service$' "$LIST_FILE" || true)"
echo "Helper scripts     : $(grep -Ec '^vouchergogo-backup-.*/sbin/' "$LIST_FILE" || true)"

echo
echo "== Preview manifest =="
if tar -xOf "$ARCHIVE" "${ROOT_DIR}/manifest.json" 2>/dev/null; then
  echo
else
  warn "manifest tidak bisa dibaca"
fi

echo
if [ "$FAILED" = "0" ]; then
  ok "backup terlihat lengkap"
  exit 0
else
  fail "backup belum lengkap, cek warning di atas"
  exit 1
fi
