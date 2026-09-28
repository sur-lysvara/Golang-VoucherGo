package main

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"net"
	"sync"
)

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := crand.Read(salt); err != nil {
		return "", err
	}

	iter := 120000
	dk := pbkdf2SHA256([]byte(password), salt, iter, 32)

	return fmt.Sprintf(
		"pbkdf2_sha256$%d$%s$%s",
		iter,
		base64.RawURLEncoding.EncodeToString(salt),
		base64.RawURLEncoding.EncodeToString(dk),
	), nil
}

func checkPassword(password, stored string) bool {
	// BCRYPT_LOGIN_COMPAT_V1: support new users created with bcrypt.
	if strings.HasPrefix(stored, "$2a$") || strings.HasPrefix(stored, "$2b$") || strings.HasPrefix(stored, "$2y$") {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) == nil
	}

	parts := strings.Split(stored, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2_sha256" {
		return false
	}

	iter, err := strconv.Atoi(parts[1])
	if err != nil || iter < 10000 {
		return false
	}

	salt, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	want, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}

	got := pbkdf2SHA256([]byte(password), salt, iter, len(want))
	return hmac.Equal(got, want)
}

func pbkdf2SHA256(password, salt []byte, iter, keyLen int) []byte {
	const hLen = 32
	numBlocks := (keyLen + hLen - 1) / hLen
	dk := make([]byte, 0, numBlocks*hLen)

	for block := 1; block <= numBlocks; block++ {
		mac := hmac.New(sha256.New, password)
		mac.Write(salt)

		var ibuf [4]byte
		binary.BigEndian.PutUint32(ibuf[:], uint32(block))
		mac.Write(ibuf[:])

		u := mac.Sum(nil)
		t := append([]byte(nil), u...)

		for i := 1; i < iter; i++ {
			mac = hmac.New(sha256.New, password)
			mac.Write(u)
			u = mac.Sum(nil)

			for j := range t {
				t[j] ^= u[j]
			}
		}

		dk = append(dk, t...)
	}

	return dk[:keyLen]
}

func makePasswordHash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (a *App) seedAdmin() error {
	var n int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	pass := env("ADMIN_DEFAULT_PASSWORD", "admin12345")

	hash, err := hashPassword(pass)
	if err != nil {
		return err
	}

	_, err = a.db.Exec(`INSERT INTO users(username,password_hash) VALUES('admin',?)`, hash)
	return err
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "tuku_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookieFlag(r)})

	// LOGOUT_ROLE_COOKIE_CLEAN_V1
	http.SetCookie(w, &http.Cookie{
		Name:     "tuku_admin_role",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Secure:   secureCookieFlag(r)})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *App) setSession(w http.ResponseWriter, uid int64) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	payload := fmt.Sprintf("%d|%d", uid, exp)
	sig := a.sign(payload)

	val := base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + sig))

	http.SetCookie(w, &http.Cookie{
		Name:     "tuku_session",
		Value:    val,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode})
}

func (a *App) sign(s string) string {
	h := hmac.New(sha256.New, a.sessionKey)
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (a *App) currentUser(r *http.Request) int64 {
	c, err := r.Cookie("tuku_session")
	if err != nil {
		return 0
	}

	raw, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return 0
	}

	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return 0
	}

	payload := parts[0] + "|" + parts[1]

	if !hmac.Equal([]byte(parts[2]), []byte(a.sign(payload))) {
		return 0
	}

	exp, _ := strconv.ParseInt(parts[1], 10, 64)
	if time.Now().Unix() > exp {
		return 0
	}

	id, _ := strconv.ParseInt(parts[0], 10, 64)
	return id
}

func (a *App) requireSuperAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.currentUser(r) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !a.isSuperAdmin(r) {
			http.Error(w, "Forbidden: hanya superadmin yang boleh akses halaman ini", http.StatusForbidden)
			return
		}

		// ADMIN_SUPER_COOKIE_CLEAN_V1
		http.SetCookie(w, &http.Cookie{
			Name:     "tuku_admin_role",
			Value:    "superadmin",
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookieFlag(r)})

		a.ensureCSRFCookie(w, r)

		fn(w, r)
	}
}

func (a *App) password(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Message string
		OK      bool
	}

	render := func(msg string, ok bool) {
		body := renderHTML(`

<div class="card" style="max-width:560px;margin:auto">
	<h2>Ganti Password Admin</h2>

	{{if .Message}}
		<div class="notice" style="margin-bottom:14px">
			<b>{{if .OK}}Berhasil{{else}}Gagal{{end}}:</b> {{.Message}}
		</div>
	{{end}}

	<form method="post">
		<div class="field">
			<label>Password Lama</label>
			<input name="old_password" type="password" required>
		</div>
		<div class="field">
			<label>Password Baru</label>
			<input name="new_password" type="password" minlength="6" required>
		</div>
		<div class="field">
			<label>Ulangi Password Baru</label>
			<input name="confirm_password" type="password" minlength="6" required>
		</div>
		<button class="btn pri">Simpan Password</button>
	</form>
</div>`, Data{Message: msg, OK: ok})

		page(w, "Ganti Password", body)
	}

	if r.Method != http.MethodPost {
		render("", false)
		return
	}

	_ = r.ParseForm()

	oldPass := r.FormValue("old_password")
	newPass := r.FormValue("new_password")
	confirm := r.FormValue("confirm_password")

	if newPass != confirm {
		render("Konfirmasi password baru tidak sama.", false)
		return
	}

	if len(newPass) < 6 {
		render("Password baru minimal 6 karakter.", false)
		return
	}

	uid := a.currentUser(r)

	var currentHash string
	err := a.db.QueryRow(`SELECT password_hash FROM users WHERE id=?`, uid).Scan(&currentHash)
	if err != nil || !checkPassword(oldPass, currentHash) {
		render("Password lama salah.", false)
		return
	}

	hash, err := hashPassword(newPass)
	if err != nil {
		render("Gagal membuat hash password.", false)
		return
	}

	_, err = a.db.Exec(`UPDATE users SET password_hash=? WHERE id=?`, hash, uid)
	if err != nil {
		render("Gagal menyimpan password.", false)
		return
	}

	render("Password admin berhasil diganti.", true)
}

func secureCookieFlag(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "https" || proto == "HTTPS" {
		return true
	}
	return r.Header.Get("X-Forwarded-Port") == "443"
}

// LOGIN_RATE_LIMIT_V1
type loginRateEntry struct {
	FailedUntil time.Time
	FailCount   int
	LastFail    time.Time
}

var loginRateMu sync.Mutex
var loginRateMap = map[string]loginRateEntry{}

func loginRateClientIP(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP"} {
		v := strings.TrimSpace(r.Header.Get(h))
		if v == "" {
			continue
		}

		if h == "X-Forwarded-For" {
			v = strings.TrimSpace(strings.Split(v, ",")[0])
		}

		if ip := net.ParseIP(v); ip != nil {
			return ip.String()
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}

func loginRateKey(r *http.Request, username string) string {
	ip := loginRateClientIP(r)
	u := strings.ToLower(strings.TrimSpace(username))
	if u == "" {
		u = "-"
	}
	return ip + "|" + u
}

func loginRateBlocked(r *http.Request, username string) (bool, time.Duration) {
	loginRateMu.Lock()
	defer loginRateMu.Unlock()

	now := time.Now()
	key := loginRateKey(r, username)
	x := loginRateMap[key]

	if !x.FailedUntil.IsZero() && now.Before(x.FailedUntil) {
		return true, time.Until(x.FailedUntil)
	}

	if !x.LastFail.IsZero() && now.Sub(x.LastFail) > 30*time.Minute {
		delete(loginRateMap, key)
	}

	return false, 0
}

func loginRateFail(r *http.Request, username string) {
	loginRateMu.Lock()
	defer loginRateMu.Unlock()

	now := time.Now()
	key := loginRateKey(r, username)
	x := loginRateMap[key]

	if !x.LastFail.IsZero() && now.Sub(x.LastFail) > 30*time.Minute {
		x = loginRateEntry{}
	}

	x.FailCount++
	x.LastFail = now

	if x.FailCount >= 5 {
		x.FailedUntil = now.Add(10 * time.Minute)
	}

	loginRateMap[key] = x
}

func loginRateSuccess(r *http.Request, username string) {
	loginRateMu.Lock()
	defer loginRateMu.Unlock()

	delete(loginRateMap, loginRateKey(r, username))
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	renderLogin := func(message string, statusCode int) {
		body := renderHTML(`
<div class="card" style="max-width:420px;margin:50px auto">
	<h2>Login Admin</h2>

	{{if .Message}}
		<div class="notice" style="margin:0 0 14px;padding:11px 13px;border-radius:14px;font-size:13px;line-height:1.4;border-left:5px solid #dc2626;background:rgba(254,242,242,.95);color:#7f1d1d">
			<b>Login gagal:</b><br>{{.Message}}
		</div>
	{{end}}

	<form method="post" action="/login">
		<div class="field">
			<label>Username</label>
			<input name="username" value="{{.Username}}" autofocus autocomplete="username">
		</div>
		<div class="field">
			<label>Password</label>
			<div class="login-pass-wrap" style="position:relative;width:100%;margin:0;">
				<input name="password" type="password" autocomplete="current-password" style="width:100%;box-sizing:border-box;padding-right:45px;">
				<span class="login-pass-eye" aria-label="Lihat password" onclick="var i=this.parentNode.querySelector('input'); if(!i)return; var show=i.type==='password'; i.type=show?'text':'password'; this.querySelector('.eye-on').style.display=show?'none':'block'; this.querySelector('.eye-off').style.display=show?'block':'none';" style="position:absolute;right:12px;top:50%;transform:translateY(-50%);cursor:pointer;color:#64748b;display:flex;align-items:center;justify-content:center;outline:none;-webkit-tap-highlight-color:transparent;">
					<svg class="eye-on" style="display:block;" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
					<svg class="eye-off" style="display:none;" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
				</span>
			</div>
		</div>

		<button class="btn pri" type="submit">Login</button>

		<div class="notice login-forgot-hint" style="margin:10px 0 12px;padding:10px 12px;border-radius:14px;font-size:13px;line-height:1.35">
			<b>Lupa password?</b><br>
			Hubungi owner / penyedia aplikasi untuk reset akun.
			<!-- LOGIN_FORGOT_PASSWORD_HINT_V1 -->
		</div>
	</form>
</div>`, map[string]any{
			"Message":  message,
			"Username": strings.TrimSpace(r.FormValue("username")),
		})

		if statusCode > 0 && statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}

		publicPage(w, "VoucherGo Admin", body)
	}

	if r.Method == http.MethodGet {
		renderLogin("", http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method tidak valid", http.StatusMethodNotAllowed)
		return
	}

	_ = r.ParseForm()

	u := strings.TrimSpace(r.FormValue("username"))
	p := r.FormValue("password")

	if blocked, wait := loginRateBlocked(r, u); blocked {
		mins := int(wait.Minutes()) + 1
		renderLogin(fmt.Sprintf("Terlalu banyak percobaan login. Coba lagi dalam %d menit.", mins), http.StatusTooManyRequests)
		return
	}

	var id int64
	var hash string

	err := a.db.QueryRow(`SELECT id,password_hash FROM users WHERE username=? AND COALESCE(enabled,1)=1`, u).Scan(&id, &hash)
	if err != nil || !checkPassword(p, hash) {
		loginRateFail(r, u)
		renderLogin("Username atau password salah.", http.StatusUnauthorized)
		return
	}

	loginRateSuccess(r, u)

	// LOGIN_ROLE_CLEAN_V1
	var loginRole string
	_ = a.db.QueryRow(`SELECT COALESCE(role,'superadmin') FROM users WHERE id=?`, id).Scan(&loginRole)
	if loginRole == "" {
		loginRole = "router"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "tuku_admin_role",
		Value:    loginRole,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   secureCookieFlag(r)})

	a.setSession(w, id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (a *App) requireAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.currentUser(r) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// ADMIN_ROLE_COOKIE_CLEAN_V1
		scopeRole, _ := a.currentAdminRole(r)
		if scopeRole == "" {
			scopeRole = "router"
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "tuku_admin_role",
			Value:    scopeRole,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookieFlag(r)})

		a.ensureCSRFCookie(w, r)

		fn(w, r)
	}
}

func (a *App) currentAdminRole(r *http.Request) (string, int64) {
	uid := a.currentUser(r)
	if uid == 0 {
		return "", 0
	}

	var role string
	var routerID sql.NullInt64

	err := a.db.QueryRow(`SELECT COALESCE(role,'superadmin'), router_id FROM users WHERE id=? AND COALESCE(enabled,1)=1`, uid).Scan(&role, &routerID)
	if err != nil {
		return "", 0
	}

	if role == "" {
		role = "superadmin"
	}

	if routerID.Valid {
		return role, routerID.Int64
	}

	return role, 0
}

func (a *App) isSuperAdmin(r *http.Request) bool {
	role, _ := a.currentAdminRole(r)
	return role == "superadmin"
}
