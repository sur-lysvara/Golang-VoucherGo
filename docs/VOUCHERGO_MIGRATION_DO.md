# VoucherGo Migration Kit - DigitalOcean / VPS Baru

Dokumen ini dipakai untuk install, backup, restore, dan pindah server VoucherGo ke VPS/server baru.

## Script Utama

- scripts/install-vouchergogo.sh
  Install dependency dan setup awal server baru.

- scripts/backup-vouchergogo.sh
  Backup full VoucherGo dari server lama.

- scripts/restore-vouchergogo.sh
  Restore backup ke server baru.

- scripts/repair-vouchergogo-permissions.sh
  Repair permission setelah restore.

- scripts/doctor-vouchergogo.sh
  Cek kesehatan server setelah install/restore.

- scripts/verify-vouchergogo-backup.sh
  Verifikasi isi file backup sebelum restore.

## WhatsApp Gateway

- Owner WA: `/var/www/wa-webjs`, service `wa-webjs.service`, port default `8666`.
- Client/demo WA baru memakai template `scripts/wa-webjs` dan service pattern `wa-webjs-<slug>`.
- Template `scripts/tuku-wa` hanya arsip legacy Baileys dan tidak dipakai untuk client/demo baru.

## 1. Backup dari Server Lama

Jalankan di server lama:

    cd /var/www/vouchergo || exit 1
    sudo bash scripts/backup-vouchergogo.sh

Hasil backup ada di:

    /home/xyzsur/backups/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz

Cek hasil backup:

    ls -lh /home/xyzsur/backups/ | tail -20

## 2. Upload Backup ke Server Baru

Dari server lama:

    scp /home/xyzsur/backups/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz root@IP_SERVER_BARU:/root/

## 3. Siapkan Server Baru

Rekomendasi server:

- Ubuntu 22.04 LTS atau Ubuntu 24.04 LTS
- RAM minimal 2GB
- Disk minimal 25GB

Login ke server baru:

    ssh root@IP_SERVER_BARU

Install Git:

    apt update && apt upgrade -y
    apt install -y git

Clone repo VoucherGo:

    cd /root
    git clone https://github.com/sur-tailedig/gozan-vouchergo.git vouchergogo
    cd /root/vouchergogo

## 4. Install Dependency Server Baru

Kalau DNS sudah mengarah ke server baru:

    sudo VG_DOMAIN="domain-baru.com" VG_EMAIL="email@domain.com" bash scripts/install-vouchergogo.sh

Kalau DNS belum mengarah atau SSL mau belakangan:

    sudo VG_DOMAIN="domain-baru.com" VG_SKIP_SSL=1 bash scripts/install-vouchergogo.sh

## 5. Restore Backup

Jalankan di server baru:

    cd /root/vouchergogo
    sudo bash scripts/restore-vouchergogo.sh /root/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz

Restore akan mengembalikan:

- /var/www/vouchergo
- /var/www/vcr-clients
- /var/www/wa-webjs
- /etc/tuku*.env
- /etc/systemd/system/tuku*.service
- /etc/systemd/system/wa-webjs*.service
- /etc/apache2/sites-available
- /etc/apache2/sites-enabled
- /usr/local/sbin/vcr-*
- database MariaDB

## 6. Repair Permission

Setelah restore:

    cd /root/vouchergogo
    sudo bash scripts/repair-vouchergogo-permissions.sh

## 7. Restart dan Cek Service

    systemctl daemon-reload
    systemctl restart apache2
    systemctl restart mariadb
    systemctl list-units 'tuku*' 'wa-webjs*' --all --no-pager
    systemctl status tuku --no-pager
    systemctl status wa-webjs --no-pager

Cek port aplikasi dan WA:

    ss -ltnp | grep -E '8097|8666'

Cek lokal:

    curl -I http://127.0.0.1:8097

## 8. Arahkan DNS ke Server Baru

Di panel DNS, arahkan A record ke IP server baru.

Contoh:

    vcr.domain.com        A    IP_SERVER_BARU
    demo.vcr.domain.com   A    IP_SERVER_BARU
    client.domain.com     A    IP_SERVER_BARU

Setelah DNS aktif, install SSL:

    certbot --apache -d domain-baru.com

Multi-domain:

    certbot --apache -d vcr.domain.com -d demo.vcr.domain.com

## 9. Checklist Setelah Migrasi

- Web admin bisa dibuka
- Login admin sukses
- Database transaksi terbaca
- Paket voucher muncul
- Router MikroTik muncul
- Kirim WA berhasil
- Client Manager jalan
- Demo Manager jalan
- Public order page jalan
- SSL aktif
- Semua service tuku* running
- Service wa-webjs* running
- Backup baru berhasil dibuat di server baru

Tes backup di server baru:

    cd /root/vouchergogo
    sudo bash scripts/backup-vouchergogo.sh

## 10. Mode Opsional

Restore tanpa database:

    sudo VG_RESTORE_DB=0 bash scripts/restore-vouchergogo.sh /root/backup.tar.gz

Restore tanpa Apache:

    sudo VG_RESTORE_APACHE=0 bash scripts/restore-vouchergogo.sh /root/backup.tar.gz

Restore tanpa start service:

    sudo VG_START_SERVICES=0 bash scripts/restore-vouchergogo.sh /root/backup.tar.gz

Install tanpa SSL:

    sudo VG_SKIP_SSL=1 VG_DOMAIN="domain.com" bash scripts/install-vouchergogo.sh

## 11. Urutan Cepat Migrasi

Server lama:

    cd /var/www/vouchergo
    sudo bash scripts/backup-vouchergogo.sh
    scp /home/xyzsur/backups/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz root@IP_SERVER_BARU:/root/

Server baru:

    apt update && apt install -y git
    cd /root
    git clone https://github.com/sur-tailedig/gozan-vouchergo.git vouchergogo
    cd /root/vouchergogo

    sudo VG_DOMAIN="domain.com" VG_EMAIL="email@domain.com" bash scripts/install-vouchergogo.sh
    sudo bash scripts/restore-vouchergogo.sh /root/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz
    sudo bash scripts/repair-vouchergogo-permissions.sh

    systemctl list-units 'tuku*' 'wa-webjs*' --all --no-pager
    curl -I http://127.0.0.1:8097


## 12. Verifikasi Backup

Sebelum backup dipindahkan atau direstore, cek dulu isinya:

    bash scripts/verify-vouchergogo-backup.sh /home/xyzsur/backups/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz

Hasil yang diharapkan:

    OK - archive bisa dibaca
    OK - database dump
    OK - env files
    OK - owner app files
    OK - WA directory
    OK - systemd service
    OK - apache config
    OK - helper scripts
    OK - backup terlihat lengkap

## 13. Doctor Check Setelah Restore

Setelah install/restore/repair permission, jalankan:

    sudo bash scripts/doctor-vouchergogo.sh

Doctor akan cek:

- OS dan resource server
- command penting
- folder app/client/WA
- file /etc/tuku*.env
- service tuku* dan wa-webjs*
- Apache
- MariaDB
- port 8097 dan 8666
- curl lokal app
- status WA local
- Git status

## 14. Urutan Final Pindah Server

Server lama:

    cd /var/www/vouchergo
    sudo bash scripts/backup-vouchergogo.sh
    bash scripts/verify-vouchergogo-backup.sh /home/xyzsur/backups/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz
    scp /home/xyzsur/backups/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz root@IP_SERVER_BARU:/root/

Server baru:

    apt update && apt install -y git
    cd /root
    git clone https://github.com/sur-tailedig/gozan-vouchergo.git vouchergogo
    cd /root/vouchergogo

    sudo VG_DOMAIN="domain.com" VG_EMAIL="email@domain.com" bash scripts/install-vouchergogo.sh
    sudo bash scripts/restore-vouchergogo.sh /root/vouchergogo-backup-YYYY-MM-DD_HHMMSS.tar.gz
    sudo bash scripts/repair-vouchergogo-permissions.sh
    sudo bash scripts/doctor-vouchergogo.sh

    systemctl list-units 'tuku*' 'wa-webjs*' --all --no-pager
    curl -I http://127.0.0.1:8097
