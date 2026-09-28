package main

import (
	"net/http"
	"strings"
)

/* OWNER_PANEL_ROUTE_GUARD_V1 */
func ownerPanelRequestAllowed(r *http.Request) bool {
	if env("OWNER_PANEL", "0") != "1" {
		return false
	}

	if r == nil {
		return false
	}

	host := strings.ToLower(r.Host)
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}

	// Safety guard: demo/client domains must never expose owner-only tools,
	// even if an env file is accidentally copied with OWNER_PANEL=1.
	if strings.HasPrefix(host, "demo.") ||
		strings.HasPrefix(host, "client-") ||
		strings.Contains(host, ".client-") {
		return false
	}

	return true
}

func requireOwnerPanelRequest(w http.ResponseWriter, r *http.Request) bool {
	if ownerPanelRequestAllowed(r) {
		return true
	}

	http.NotFound(w, r)
	return false
}

/* OWNER_PANEL_PREAUTH_ROUTE_GUARD_V2 */
func (a *App) requireOwnerPanel(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireOwnerPanelRequest(w, r) {
			return
		}
		next(w, r)
	}
}
