package main

import (
	"net/http"
	"net/url"
	"strings"
)

const pendingInvoiceCookieName = "tuku_pending_invoice"

func setPendingInvoiceCookie(w http.ResponseWriter, r *http.Request, orderID string) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     pendingInvoiceCookieName,
		Value:    orderID,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secureCookieFlag(r),
	})
}

func clearPendingInvoiceCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     pendingInvoiceCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secureCookieFlag(r),
	})
}

func clearPendingInvoiceCookieIfMatches(w http.ResponseWriter, r *http.Request, orderID string) {
	c, err := r.Cookie(pendingInvoiceCookieName)
	if err != nil {
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Value), strings.TrimSpace(orderID)) {
		clearPendingInvoiceCookie(w, r)
	}
}

func (a *App) pendingInvoiceRedirectForRouter(r *http.Request, routerSlug string) (target string, clear bool) {
	c, err := r.Cookie(pendingInvoiceCookieName)
	if err != nil {
		return "", false
	}

	orderID := strings.TrimSpace(c.Value)
	if orderID == "" || len(orderID) > 80 || strings.ContainsAny(orderID, "/?#& \t\r\n") {
		return "", true
	}

	var status string
	var txRouterSlug string

	err = a.db.QueryRow(`
		SELECT LOWER(COALESCE(t.status,'')), r.slug
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		WHERE t.order_id=?
		  AND t.created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY)
		LIMIT 1`, orderID).Scan(&status, &txRouterSlug)

	if err != nil {
		return "", true
	}

	if !strings.EqualFold(txRouterSlug, strings.TrimSpace(routerSlug)) {
		return "", false
	}

	if status == "pending" {
		return "/invoice/" + url.PathEscape(orderID), false
	}

	return "", true
}
