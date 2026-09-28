package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type PhoneBan struct {
	ID                 int64
	PhoneNormalized    string
	PhoneDisplay       string
	Reason             string
	IsActive           bool
	IsCurrentlyBlocked bool
	BannedUntil        sql.NullString
	CreatedAt          string
	UpdatedAt          string
}

func normalizePhoneForBanCheck(phone string) string {
	phone = strings.TrimSpace(phone)
	replacer := strings.NewReplacer(
		" ", "",
		"-", "",
		"(", "",
		")", "",
		".", "",
	)
	phone = replacer.Replace(phone)
	if strings.HasPrefix(phone, "+62") {
		phone = "62" + strings.TrimPrefix(phone, "+62")
	} else if strings.HasPrefix(phone, "08") {
		phone = "628" + strings.TrimPrefix(phone, "08")
	} else if strings.HasPrefix(phone, "8") {
		phone = "628" + strings.TrimPrefix(phone, "8")
	}
	return phone
}

func getActivePhoneBan(ctx context.Context, db *sql.DB, phone string) (*PhoneBan, error) {
	normalized := normalizePhoneForBanCheck(phone)
	if normalized == "" {
		return nil, nil
	}

	var ban PhoneBan
	var active int
	err := db.QueryRowContext(ctx, `SELECT
			id,
			phone_normalized,
			COALESCE(phone_display,''),
			reason,
			is_active,
			COALESCE(DATE_FORMAT(banned_until,'%Y-%m-%d %H:%i:%s'),''),
			DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(updated_at,'%Y-%m-%d %H:%i:%s')
		FROM blocked_customer_phones
		WHERE phone_normalized=?
		  AND is_active=1
		  AND (banned_until IS NULL OR banned_until > NOW())
		LIMIT 1`, normalized).Scan(
		&ban.ID,
		&ban.PhoneNormalized,
		&ban.PhoneDisplay,
		&ban.Reason,
		&active,
		&ban.BannedUntil.String,
		&ban.CreatedAt,
		&ban.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ban.IsActive = active == 1
	ban.IsCurrentlyBlocked = ban.IsActive
	ban.BannedUntil.Valid = ban.BannedUntil.String != ""
	return &ban, nil
}

func isPhoneBlocked(ctx context.Context, db *sql.DB, phone string) (bool, string, error) {
	ban, err := getActivePhoneBan(ctx, db, phone)
	if err != nil {
		return false, "", err
	}
	if ban == nil {
		return false, "", nil
	}
	return true, ban.Reason, nil
}
