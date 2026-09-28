#!/usr/bin/env bash
set -euo pipefail

ACTION="${1:-}"
SLUG="${2:-}"
MONTHS="${3:-1}"

fail(){ echo "ERROR: $*" >&2; exit 1; }

[[ "$ACTION" =~ ^(status|paid|suspend|unsuspend)$ ]] || fail "action harus status/paid/suspend/unsuspend"
[[ "$SLUG" =~ ^[a-zA-Z0-9][a-zA-Z0-9_-]{0,60}$ ]] || fail "slug tidak valid"
[[ "$MONTHS" =~ ^[0-9]+$ ]] || fail "months tidak valid"
[ "$MONTHS" -ge 1 ] && [ "$MONTHS" -le 24 ] || fail "months harus 1-24"

ENVF="/etc/tuku-${SLUG}.env"
SVC="tuku-${SLUG}"

[ -f "$ENVF" ] || fail "env tidak ditemukan: $ENVF"

if [ "$ACTION" = "status" ]; then
  grep -E '^(CLIENT_SLUG|CLIENT_BILLING_STATUS|CLIENT_BILLING_DUE_AT|CLIENT_BILLING_GRACE_UNTIL|CLIENT_BILLING_PRICE|CLIENT_BILLING_OWNER_URL|CLIENT_BILLING_LAST_PAID_AT)=' "$ENVF" || true
  exit 0
fi

if [ "$ACTION" = "suspend" ] || [ "$ACTION" = "unsuspend" ]; then
  NEW_STATUS="suspended"
  [ "$ACTION" = "unsuspend" ] && NEW_STATUS="active"

  python3 - "$ENVF" "$NEW_STATUS" <<'PY'
from pathlib import Path
import sys

envf = Path(sys.argv[1])
new_status = sys.argv[2]

s = envf.read_text()
lines = s.splitlines()
out = []
found = False

for line in lines:
    raw = line.strip()
    if raw.startswith("CLIENT_BILLING_STATUS="):
        out.append(f'CLIENT_BILLING_STATUS="{new_status}"')
        found = True
    else:
        out.append(line)

if not found:
    out.append(f'CLIENT_BILLING_STATUS="{new_status}"')

envf.write_text("\n".join(out).rstrip() + "\n")
print(f"CLIENT_BILLING_STATUS={new_status}")
PY

  chown root:root "$ENVF"
  chmod 600 "$ENVF"

  if command -v setfacl >/dev/null 2>&1; then
    setfacl -m u:www-data:r "$ENVF" || true
  fi

  systemctl restart "$SVC"
  systemctl is-active "$SVC" >/dev/null

  echo "OK: billing $ACTION applied for $SLUG"
  exit 0
fi

python3 - "$ENVF" "$MONTHS" <<'PY'
from pathlib import Path
from datetime import date, timedelta, datetime
from calendar import monthrange
import sys

envf = Path(sys.argv[1])
months = int(sys.argv[2])
s = envf.read_text()

def get_value(text, key):
    for line in text.splitlines():
        raw = line.strip()
        if raw.startswith(key + "="):
            return raw.split("=", 1)[1].strip().strip('"').strip("'")
    return ""

def parse_date(v):
    try:
        return datetime.strptime(v, "%Y-%m-%d").date()
    except Exception:
        return None

def add_months(d, m):
    month = d.month - 1 + m
    year = d.year + month // 12
    month = month % 12 + 1
    day = min(d.day, monthrange(year, month)[1])
    return date(year, month, day)

today = date.today()
old_due = parse_date(get_value(s, "CLIENT_BILLING_DUE_AT"))
base = old_due if old_due and old_due > today else today
new_due = add_months(base, months)
new_grace = new_due + timedelta(days=3)

updates = {
    "CLIENT_BILLING_STATUS": "active",
    "CLIENT_BILLING_DUE_AT": new_due.isoformat(),
    "CLIENT_BILLING_GRACE_UNTIL": new_grace.isoformat(),
    "CLIENT_BILLING_LAST_PAID_AT": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
}

lines = s.splitlines()
seen = set()
out = []

for line in lines:
    stripped = line.strip()
    if "=" in stripped and not stripped.startswith("#"):
        key = stripped.split("=", 1)[0].strip()
        if key in updates:
            out.append(f'{key}="{updates[key]}"')
            seen.add(key)
            continue
    out.append(line)

for key, val in updates.items():
    if key not in seen:
        out.append(f'{key}="{val}"')

envf.write_text("\n".join(out).rstrip() + "\n")

print(f"CLIENT_BILLING_DUE_AT={new_due.isoformat()}")
print(f"CLIENT_BILLING_GRACE_UNTIL={new_grace.isoformat()}")
PY

chown root:root "$ENVF"
chmod 600 "$ENVF"

if command -v setfacl >/dev/null 2>&1; then
  setfacl -m u:www-data:r "$ENVF" || true
fi

systemctl restart "$SVC"
systemctl is-active "$SVC" >/dev/null

echo "OK: billing paid applied for $SLUG"
