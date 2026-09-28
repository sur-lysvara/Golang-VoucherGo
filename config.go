package main

import (
	"log"
	"os"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func env(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func vipBrandingEnabled() bool {
	return env("OWNER_PANEL", "0") == "1" || env("VIP_BRANDING", "0") == "1"
}
