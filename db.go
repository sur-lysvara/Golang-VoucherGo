package main

import "database/sql"

func (a *App) setting(k string) string {
	var v sql.NullString
	_ = a.db.QueryRow(`SELECT v FROM settings WHERE k=?`, k).Scan(&v)
	return v.String
}

func (a *App) setSetting(k, v string) error {
	_, err := a.db.Exec(
		`INSERT INTO settings(k,v) VALUES(?,?) ON DUPLICATE KEY UPDATE v=VALUES(v)`,
		k,
		v,
	)
	return err
}

func (a *App) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(30) NOT NULL DEFAULT 'superadmin',
			router_id BIGINT NULL,
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS settings (
			k VARCHAR(100) PRIMARY KEY,
			v TEXT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS routers (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			slug VARCHAR(120) NOT NULL UNIQUE,
			host VARCHAR(180) NULL,
			port INT DEFAULT 443,
			api_mode VARCHAR(20) NOT NULL DEFAULT 'auto',
			username VARCHAR(120) NULL,
			password TEXT NULL,
			use_ssl TINYINT(1) DEFAULT 1,
			insecure_ssl TINYINT(1) DEFAULT 1,
			enabled TINYINT(1) DEFAULT 1,
			last_test_status VARCHAR(20) NULL,
			last_test_message TEXT NULL,
			last_test_at DATETIME NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS packages (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			router_id BIGINT NOT NULL,
			name VARCHAR(120) NOT NULL,
			price BIGINT NOT NULL,
			profile VARCHAR(120) NULL,
			limit_uptime VARCHAR(50) NULL,
			shared_users INT DEFAULT 1,
			active TINYINT(1) DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX(router_id),
			CONSTRAINT fk_packages_router FOREIGN KEY (router_id) REFERENCES routers(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS transactions (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			order_id VARCHAR(80) NOT NULL UNIQUE,
			router_id BIGINT NOT NULL,
			package_id BIGINT NOT NULL,
			customer_name VARCHAR(120) NULL,
			customer_phone VARCHAR(50) NULL,
			amount BIGINT NOT NULL,
			status ENUM('pending','paid','failed','expired','cancelled') DEFAULT 'pending',
			payment_method VARCHAR(50) NULL,
			voucher_username VARCHAR(80) NULL,
			voucher_password VARCHAR(80) NULL,
			provision_status ENUM('pending','success','failed','skipped') DEFAULT 'pending',
			provision_error TEXT NULL,
			raw_webhook JSON NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			paid_at DATETIME NULL,
			INDEX(router_id),
			INDEX(package_id),
			INDEX(status),
			CONSTRAINT fk_tx_router FOREIGN KEY (router_id) REFERENCES routers(id),
			CONSTRAINT fk_tx_package FOREIGN KEY (package_id) REFERENCES packages(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,

		`CREATE TABLE IF NOT EXISTS blocked_customer_phones (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			phone_normalized VARCHAR(32) NOT NULL UNIQUE,
			phone_display VARCHAR(32) NULL,
			reason TEXT NOT NULL,
			is_active TINYINT(1) NOT NULL DEFAULT 1,
			banned_until DATETIME NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}

	for _, s := range stmts {
		if _, err := a.db.Exec(s); err != nil {
			return err
		}
	}

	// USERS_ROLE_SCOPE_V1
	// Ignore duplicate-column errors so old installs can migrate safely.
	_, _ = a.db.Exec(`ALTER TABLE users ADD COLUMN role VARCHAR(30) NOT NULL DEFAULT 'superadmin' AFTER password_hash`)
	_, _ = a.db.Exec(`ALTER TABLE users ADD COLUMN router_id BIGINT NULL AFTER role`)
	_, _ = a.db.Exec(`ALTER TABLE users ADD COLUMN enabled TINYINT(1) NOT NULL DEFAULT 1 AFTER router_id`)

	// ROUTER_LOGIN_URL_V1
	// Ignore error so repeated migration remains safe when column already exists.
	_, _ = a.db.Exec(`ALTER TABLE routers ADD COLUMN hotspot_login_url VARCHAR(255) NULL AFTER slug`)
	_, _ = a.db.Exec(`ALTER TABLE routers ADD COLUMN api_mode VARCHAR(20) NOT NULL DEFAULT 'auto' AFTER port`)
	_, _ = a.db.Exec(`ALTER TABLE routers ADD COLUMN last_test_status VARCHAR(20) NULL AFTER enabled`)
	_, _ = a.db.Exec(`ALTER TABLE routers ADD COLUMN last_test_message TEXT NULL AFTER last_test_status`)
	_, _ = a.db.Exec(`ALTER TABLE routers ADD COLUMN last_test_at DATETIME NULL AFTER last_test_message`)

	// PAYMENT_TEMP_BYPASS_V1
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN customer_email VARCHAR(180) NULL AFTER customer_phone`)
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN hotspot_mac VARCHAR(32) NULL AFTER customer_email`)
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN hotspot_ip VARCHAR(64) NULL AFTER hotspot_mac`)

	// PAKASIR_INTERNAL_QRIS_V1
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN qris_string TEXT NULL AFTER hotspot_ip`)
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN qris_total_payment BIGINT NULL AFTER qris_string`)
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN qris_expired_at VARCHAR(80) NULL AFTER qris_total_payment`)

	// PAKASIR_V2_TXN_ID
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN pakasir_txn_id VARCHAR(100) NULL AFTER qris_expired_at`)

	// PAKASIR_V2_BILLING
	_, _ = a.db.Exec(`ALTER TABLE client_billing_invoices ADD COLUMN pakasir_txn_id VARCHAR(100) NULL AFTER amount`)
	_, _ = a.db.Exec(`ALTER TABLE client_billing_invoices ADD COLUMN qris_string TEXT NULL AFTER pakasir_txn_id`)
	_, _ = a.db.Exec(`ALTER TABLE client_billing_invoices ADD COLUMN qris_total_payment BIGINT NULL AFTER qris_string`)
	_, _ = a.db.Exec(`ALTER TABLE client_billing_invoices ADD COLUMN qris_expired_at VARCHAR(80) NULL AFTER qris_total_payment`)

	// PAYMENT_BYPASS_LIMIT_V1
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN payment_bypass_count INT NOT NULL DEFAULT 0 AFTER qris_expired_at`)
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN payment_bypass_until DATETIME NULL AFTER payment_bypass_count`)
	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN payment_bypass_last_at DATETIME NULL AFTER payment_bypass_until`)

	// UNUSED_VOUCHER_REMINDER_DB_V1

	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN unused_reminder_sent_at DATETIME NULL AFTER payment_bypass_last_at`)

	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN unused_reminder_channel VARCHAR(80) NULL AFTER unused_reminder_sent_at`)

	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN first_used_at DATETIME NULL AFTER unused_reminder_channel`)

	_, _ = a.db.Exec(`ALTER TABLE transactions ADD COLUMN last_usage_check_at DATETIME NULL AFTER first_used_at`)

	defaults := map[string]string{

		"unused_voucher_reminder_enabled": "0",

		"unused_voucher_reminder_minutes": "60",

		"unused_voucher_reminder_channel": "wa,email",

		"unused_voucher_reminder_limit": "25",
		"site_name":                     "VoucherGo",
		"public_base_url":               "https://vouchergo.biz.id",
		"public_subtitle":               "Beli voucher hotspot, bayar via QRIS.",
		"maintenance_enabled":           "0",
		"maintenance_title":             "Sedang Update",
		"maintenance_message":           "Mohon maaf, layanan sedang diperbarui. Silakan coba lagi beberapa saat lagi.",
		"wa_provider":                   "local",
		"wa_local_url":                  "http://127.0.0.1:8666/send",
		"wa_local_key":                  "local-wa-key",
		"payment_provider":              "pakasir",
		"midtrans_mode":                 "sandbox",
		"midtrans_server_key":           "",
		"midtrans_client_key":           "",
		"midtrans_merchant_id":          "",
		"pakasir_slug":                  "",
		"pakasir_api_key":               "",
		"pakasir_webhook_secret":        "",
		"pakasir_qris_only":             "1",
		"wanesia_enabled":               "0",
		"wanesia_base_url":              "",
		"wanesia_api_key":               "",
		"wanesia_sender":                "",
	}

	for k, v := range defaults {
		if _, err := a.db.Exec(`INSERT IGNORE INTO settings(k,v) VALUES(?,?)`, k, v); err != nil {
			return err
		}
	}

	return nil
}
