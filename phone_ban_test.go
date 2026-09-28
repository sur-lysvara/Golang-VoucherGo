package main

import "testing"

func TestNormalizePhoneForBanCheck(t *testing.T) {
	tests := map[string]string{
		"+62 812-3456.7890": "6281234567890",
		"0812 3456 7890":    "6281234567890",
		"8(123)456-7890":    "6281234567890",
		"6281234567890":     "6281234567890",
		"":                  "",
	}

	for in, want := range tests {
		if got := normalizePhoneForBanCheck(in); got != want {
			t.Fatalf("normalizePhoneForBanCheck(%q) = %q, want %q", in, got, want)
		}
	}
}
