package main

import "strings"

// VOUCHERGO_MIKROTIK_DUPLICATE_USER_IDEMPOTENT_V1
func isMikroTikDuplicateUserText(s string) bool {
	x := strings.ToLower(s)
	return strings.Contains(x, "already have user with this name") ||
		(strings.Contains(x, "already exists") && strings.Contains(x, "user"))
}
