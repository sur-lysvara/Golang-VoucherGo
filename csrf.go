package main

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const csrfCookieName = "tuku_csrf"

func (a *App) makeCSRFToken(uid int64) string {
	exp := time.Now().Add(12 * time.Hour).Unix()
	payload := fmt.Sprintf("%d|%d", uid, exp)
	sig := a.sign("csrf|" + payload)
	return base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + sig))
}

func (a *App) validCSRFToken(r *http.Request, token string) bool {
	if token == "" {
		return false
	}

	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return false
	}

	uid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || uid == 0 || uid != a.currentUser(r) {
		return false
	}

	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}

	payload := parts[0] + "|" + parts[1]
	want := a.sign("csrf|" + payload)
	return hmac.Equal([]byte(parts[2]), []byte(want))
}

func (a *App) ensureCSRFCookie(w http.ResponseWriter, r *http.Request) {
	uid := a.currentUser(r)
	if uid == 0 {
		return
	}

	if c, err := r.Cookie(csrfCookieName); err == nil && a.validCSRFToken(r, c.Value) {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    a.makeCSRFToken(uid),
		Path:     "/",
		MaxAge:   43200,
		SameSite: http.SameSiteLaxMode,
		Secure:   secureCookieFlag(r),
	})
}

func sameOriginAdminPost(r *http.Request) bool {
	host := r.Host
	if host == "" {
		return false
	}

	checkURL := func(raw string) bool {
		if strings.TrimSpace(raw) == "" {
			return false
		}
		u, err := url.Parse(raw)
		if err != nil {
			return false
		}
		return strings.EqualFold(u.Host, host)
	}

	if checkURL(r.Header.Get("Origin")) {
		return true
	}
	if checkURL(r.Header.Get("Referer")) {
		return true
	}

	return false
}

func (a *App) requireCSRF(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			fn(w, r)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "CSRF invalid form", http.StatusForbidden)
			return
		}

		formToken := r.FormValue("_csrf")
		if formToken == "" {
			formToken = r.Header.Get("X-CSRF-Token")
		}

		cookieToken := ""
		if c, err := r.Cookie(csrfCookieName); err == nil {
			cookieToken = c.Value
		}

		if formToken == "" || cookieToken == "" || !hmac.Equal([]byte(formToken), []byte(cookieToken)) || !a.validCSRFToken(r, formToken) {
			// CSRF_SAME_ORIGIN_FALLBACK_V1
			// Tetap blok cross-site request, tapi jangan bikin form admin lama/dinamis rusak
			// saat hidden _csrf belum tersisip oleh JS.
			if !sameOriginAdminPost(r) {
				http.Error(w, "CSRF token invalid", http.StatusForbidden)
				return
			}
		}

		fn(w, r)
	}
}
