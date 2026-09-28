# Local WhatsApp Gateway - whatsapp-web.js

Template local WhatsApp gateway for VoucherGo owner, demo, and client services.

## Runtime

- Owner app path: `/var/www/wa-webjs`
- Owner service: `wa-webjs`
- Owner URL: `http://127.0.0.1:8666`
- Client service pattern: `wa-webjs-<client>`
- Provider: `whatsapp-web.js`
- Port, API key, cache dir, and session dir are set from systemd environment variables.

## VoucherGo Settings

- WA Provider: `local`
- WA Local URL: `http://127.0.0.1:<port>/send`
- WA Local Key: must match the service `API_KEY`.

## Endpoints

- `GET /status`
- `GET /check`
- `GET /qr`
- `POST /send`
- `POST /send-image`
- `POST /`
- `POST /logout`
- `POST /reconnect`

## Notes

`/send-image` and `POST /` are compatibility endpoints for VoucherGo auto invoice QRIS and manual WA resend.

Do not commit:
- `node_modules/`
- `sessions/`
- `.cache/`
- `.config/`
