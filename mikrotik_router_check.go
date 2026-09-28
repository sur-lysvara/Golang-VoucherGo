package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (a *App) routerTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method tidak valid", http.StatusMethodNotAllowed)
		return
	}

	_ = r.ParseForm()

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	ro, err := a.routerFromTestRequest(r, id)

	ok := false
	msg := ""

	if err != nil {
		msg = err.Error()
	} else {
		ok, msg = routerConnectionStatus(ro)
	}

	status := "fail"
	if ok {
		status = "ok"
	}

	if id > 0 {
		_, _ = a.db.Exec(`UPDATE routers SET last_test_status=?, last_test_message=?, last_test_at=NOW() WHERE id=?`, status, msg, id)
	}

	target := strings.TrimSpace(r.FormValue("return_to"))
	if target == "" {
		target = "/admin/routers"
		if id > 0 {
			target = fmt.Sprintf("/admin/routers?edit=%d", id)
		}
	}

	// Safety: hanya boleh redirect internal ke halaman Router.
	if target != "/admin/routers" && !strings.HasPrefix(target, "/admin/routers?") {
		target = "/admin/routers"
	}

	sep := "?"
	if strings.Contains(target, "?") {
		sep = "&"
	}

	http.Redirect(w, r, target+sep+"router_test="+status+"&router_msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

func (a *App) routerFromTestRequest(r *http.Request, id int64) (Router, error) {
	if _, exists := r.Form["host"]; exists {
		apiMode := normalizeRouterAPIMode(r.FormValue("api_mode"))

		port, _ := strconv.Atoi(r.FormValue("port"))
		if port <= 0 {
			switch apiMode {
			case "classic":
				port = 8728
			case "classic_ssl":
				port = 8729
			default:
				port = 443
			}
		}

		password := r.FormValue("password")
		if strings.TrimSpace(password) == "" && id > 0 {
			_ = a.db.QueryRow(`SELECT COALESCE(password,'') FROM routers WHERE id=?`, id).Scan(&password)
		}

		return Router{
			ID:          id,
			Host:        strings.TrimSpace(r.FormValue("host")),
			Port:        port,
			APIMode:     apiMode,
			Username:    strings.TrimSpace(r.FormValue("username")),
			Password:    password,
			UseSSL:      r.FormValue("use_ssl") == "on",
			InsecureSSL: r.FormValue("insecure_ssl") == "on",
			Enabled:     true,
		}, nil
	}

	if id <= 0 {
		return Router{}, fmt.Errorf("router belum dipilih")
	}

	var ro Router
	err := a.db.QueryRow(`SELECT
			id,
			name,
			slug,
			COALESCE(hotspot_login_url,''),
			COALESCE(host,''),
			port,
			COALESCE(api_mode,'auto'),
			COALESCE(username,''),
			COALESCE(password,''),
			use_ssl,
			insecure_ssl,
			enabled
		FROM routers
		WHERE id=?`, id).
		Scan(
			&ro.ID,
			&ro.Name,
			&ro.Slug,
			&ro.HotspotLoginURL,
			&ro.Host,
			&ro.Port,
			&ro.APIMode,
			&ro.Username,
			&ro.Password,
			&ro.UseSSL,
			&ro.InsecureSSL,
			&ro.Enabled,
		)

	if err != nil {
		return Router{}, err
	}

	return ro, nil
}

func routerConnectionStatus(ro Router) (bool, string) {
	if strings.TrimSpace(ro.Host) == "" {
		return false, "Host/IP MikroTik kosong"
	}

	if strings.TrimSpace(ro.Username) == "" {
		return false, "Username MikroTik kosong"
	}

	mode := selectedMTAPIMode(ro.APIMode)

	switch mode {
	case "rest":
		if err := routerTestREST(ro); err != nil {
			return false, "REST API gagal: " + shortRouterErr(err)
		}
		return true, fmt.Sprintf("OK via REST API ke %s:%d", cleanRouterHost(ro.Host), ro.Port)

	case "classic":
		if err := routerTestClassic(ro, false, true); err != nil {
			return false, "Classic API gagal: " + shortRouterErr(err)
		}
		return true, fmt.Sprintf("OK via Classic API ke %s:%d", cleanRouterHost(ro.Host), mtClassicRouterPort(ro, false, true))

	case "classic_ssl":
		if err := routerTestClassic(ro, true, true); err != nil {
			return false, "Classic API SSL gagal: " + shortRouterErr(err)
		}
		return true, fmt.Sprintf("OK via Classic API SSL ke %s:%d", cleanRouterHost(ro.Host), mtClassicRouterPort(ro, true, true))

	default:
		restErr := routerTestREST(ro)
		if restErr == nil {
			return true, fmt.Sprintf("OK via REST API ke %s:%d", cleanRouterHost(ro.Host), ro.Port)
		}

		classicErr := routerTestClassic(ro, false, false)
		if classicErr == nil {
			return true, fmt.Sprintf("OK via Classic API ke %s:%d", cleanRouterHost(ro.Host), mtClassicRouterPort(ro, false, false))
		}

		return false, "Auto gagal. REST: " + shortRouterErr(restErr) + " | Classic: " + shortRouterErr(classicErr)
	}
}

func routerTestREST(ro Router) error {
	if ro.Port <= 0 {
		ro.Port = 443
	}

	scheme := "http"
	if ro.UseSSL {
		scheme = "https"
	}

	u := fmt.Sprintf("%s://%s:%d/rest/system/resource", scheme, cleanRouterHost(ro.Host), ro.Port)

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.SetBasicAuth(ro.Username, ro.Password)

	resp, err := mtRestClientForRouter(ro).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rb, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		{
			// VOUCHERGO_MIKROTIK_DUPLICATE_USER_IDEMPOTENT_V1_CALL
			mtBodyText := string(rb)
			if isMikroTikDuplicateUserText(mtBodyText) {
				return nil
			}
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, mtBodyText)
		}
	}

	return nil
}

func routerTestClassic(ro Router, useTLS bool, useRouterPort bool) error {
	addr := net.JoinHostPort(cleanRouterHost(ro.Host), strconv.Itoa(mtClassicRouterPort(ro, useTLS, useRouterPort)))

	c, err := mtClassicConnect(addr, ro.Username, ro.Password, useTLS, ro.InsecureSSL)
	if err != nil {
		return err
	}
	defer c.Close()

	return c.Run("/system/resource/print")
}

func shortRouterErr(err error) string {
	if err == nil {
		return ""
	}

	s := strings.TrimSpace(err.Error())
	if len(s) > 240 {
		s = s[:240] + "..."
	}

	return s
}
