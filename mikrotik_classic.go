package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type mtClassicClient struct {
	conn net.Conn
}

func provisionVoucherClassic(ro Router, p Package, orderID, user, pass string, useTLS bool, useRouterPort bool) error {
	addr := net.JoinHostPort(cleanRouterHost(ro.Host), strconv.Itoa(mtClassicRouterPort(ro, useTLS, useRouterPort)))

	c, err := mtClassicConnect(addr, ro.Username, ro.Password, useTLS, ro.InsecureSSL)
	if err != nil {
		return err
	}
	defer c.Close()

	args := []string{
		"/ip/hotspot/user/add",
		"=name=" + user,
		"=password=" + pass,
		"=profile=" + p.Profile,
		"=comment=" + voucherGoComment(orderID),
	}

	if p.LimitUptime != "" {
		args = append(args, "=limit-uptime="+p.LimitUptime)
	}

	return c.Run(args...)
}

func mtClassicConnect(addr, user, pass string, useTLS bool, insecure bool) (*mtClassicClient, error) {
	c, err := mtClassicDial(addr, useTLS, insecure)
	if err != nil {
		return nil, err
	}

	if err := c.loginNew(user, pass); err == nil {
		return c, nil
	}

	_ = c.Close()

	c, err = mtClassicDial(addr, useTLS, insecure)
	if err != nil {
		return nil, err
	}

	if err := c.loginLegacy(user, pass); err != nil {
		_ = c.Close()
		return nil, err
	}

	return c, nil
}

func mtClassicDial(addr string, useTLS bool, insecure bool) (*mtClassicClient, error) {
	d := net.Dialer{Timeout: 12 * time.Second}

	var conn net.Conn
	var err error

	if useTLS {
		cfg := &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: insecure,
		}
		conn, err = tls.DialWithDialer(&d, "tcp", addr, cfg)
	} else {
		conn, err = d.Dial("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	_ = conn.SetDeadline(time.Now().Add(20 * time.Second))

	return &mtClassicClient{conn: conn}, nil
}

func (c *mtClassicClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *mtClassicClient) loginNew(user, pass string) error {
	return c.Run("/login", "=name="+user, "=password="+pass)
}

func (c *mtClassicClient) loginLegacy(user, pass string) error {
	replies, err := c.runRaw("/login")
	if err != nil {
		return err
	}

	ret := ""
	for _, r := range replies {
		if v := r["ret"]; v != "" {
			ret = v
			break
		}
	}

	if ret == "" {
		return errors.New("legacy login challenge kosong")
	}

	challenge, err := hex.DecodeString(ret)
	if err != nil {
		return fmt.Errorf("legacy login challenge invalid: %w", err)
	}

	h := md5.New() // RouterOS legacy API challenge-response.
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(pass))
	_, _ = h.Write(challenge)

	response := "00" + hex.EncodeToString(h.Sum(nil))

	return c.Run("/login", "=name="+user, "=response="+response)
}

func (c *mtClassicClient) Run(words ...string) error {
	_, err := c.runRaw(words...)
	return err
}

func (c *mtClassicClient) runRaw(words ...string) ([]map[string]string, error) {
	if err := c.writeSentence(words); err != nil {
		return nil, err
	}

	var replies []map[string]string

	for {
		sentence, err := c.readSentence()
		if err != nil {
			return replies, err
		}

		if len(sentence) == 0 {
			continue
		}

		kind := sentence[0]
		attrs := parseMTAttrs(sentence[1:])
		replies = append(replies, attrs)

		switch kind {
		case "!done":
			return replies, nil
		case "!trap", "!fatal":
			msg := attrs["message"]
			if msg == "" {
				msg = strings.Join(sentence, " ")
			}
			return replies, errors.New(msg)
		}
	}
}

func (c *mtClassicClient) writeSentence(words []string) error {
	for _, w := range words {
		if err := writeMTWord(c.conn, w); err != nil {
			return err
		}
	}
	return writeMTWord(c.conn, "")
}

func (c *mtClassicClient) readSentence() ([]string, error) {
	var words []string

	for {
		w, err := readMTWord(c.conn)
		if err != nil {
			return words, err
		}

		if w == "" {
			return words, nil
		}

		words = append(words, w)
	}
}

func parseMTAttrs(words []string) map[string]string {
	m := make(map[string]string)

	for _, w := range words {
		if !strings.HasPrefix(w, "=") {
			continue
		}

		x := strings.TrimPrefix(w, "=")
		parts := strings.SplitN(x, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}

	return m
}

func writeMTWord(w io.Writer, s string) error {
	b := []byte(s)

	if err := writeMTLen(w, len(b)); err != nil {
		return err
	}

	if len(b) == 0 {
		return nil
	}

	_, err := w.Write(b)
	return err
}

func readMTWord(r io.Reader) (string, error) {
	n, err := readMTLen(r)
	if err != nil {
		return "", err
	}

	if n == 0 {
		return "", nil
	}

	b := make([]byte, n)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func writeMTLen(w io.Writer, n int) error {
	var b []byte

	switch {
	case n < 0x80:
		b = []byte{byte(n)}
	case n < 0x4000:
		b = []byte{byte((n >> 8) | 0x80), byte(n)}
	case n < 0x200000:
		b = []byte{byte((n >> 16) | 0xC0), byte(n >> 8), byte(n)}
	case n < 0x10000000:
		b = []byte{byte((n >> 24) | 0xE0), byte(n >> 16), byte(n >> 8), byte(n)}
	default:
		b = []byte{0xF0, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	}

	_, err := w.Write(b)
	return err
}

func readMTLen(r io.Reader) (int, error) {
	var first [1]byte
	if _, err := io.ReadFull(r, first[:]); err != nil {
		return 0, err
	}

	c := first[0]

	switch {
	case c&0x80 == 0x00:
		return int(c), nil
	case c&0xC0 == 0x80:
		var b [1]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return 0, err
		}
		return int(c&^0xC0)<<8 + int(b[0]), nil
	case c&0xE0 == 0xC0:
		var b [2]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return 0, err
		}
		return int(c&^0xE0)<<16 + int(b[0])<<8 + int(b[1]), nil
	case c&0xF0 == 0xE0:
		var b [3]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return 0, err
		}
		return int(c&^0xF0)<<24 + int(b[0])<<16 + int(b[1])<<8 + int(b[2]), nil
	default:
		var b [4]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return 0, err
		}
		return int(b[0])<<24 + int(b[1])<<16 + int(b[2])<<8 + int(b[3]), nil
	}
}

type routerAPIModeOption struct {
	Value    string
	Label    string
	Selected bool
}

func routerAPIModeOptions(current string) []routerAPIModeOption {
	current = normalizeRouterAPIMode(current)

	items := []routerAPIModeOption{
		{Value: "auto", Label: "Auto"},
		{Value: "rest", Label: "REST API (RouterOS7)"},
		{Value: "classic", Label: "Classic API (RouterOS6/RouterOS7)"},
		{Value: "classic_ssl", Label: "Classic API SSL (RouterOS6/RouterOS7)"},
	}

	for i := range items {
		items[i].Selected = items[i].Value == current
	}

	return items
}

func selectedMTAPIMode(routerMode string) string {
	mode := normalizeRouterAPIMode(routerMode)
	if mode != "auto" {
		return mode
	}

	if strings.TrimSpace(os.Getenv("MIKROTIK_API_MODE")) != "" {
		return mtAPIMode()
	}

	return "auto"
}

func normalizeRouterAPIMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))

	switch mode {
	case "rest":
		return "rest"
	case "classic", "api", "api_classic", "classic_api":
		return "classic"
	case "classic_ssl", "api_ssl", "ssl":
		return "classic_ssl"
	default:
		return "auto"
	}
}

func mtClassicRouterPort(ro Router, useTLS bool, useRouterPort bool) int {
	if useRouterPort && ro.Port > 0 && ro.Port != 443 {
		return ro.Port
	}

	return mtClassicAPIPort(useTLS)
}

func mtAPIMode() string {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("MIKROTIK_API_MODE")))

	switch mode {
	case "rest":
		return "rest"
	case "classic", "api", "api_classic", "classic_api":
		return "classic"
	case "classic_ssl", "api_ssl", "ssl":
		return "classic_ssl"
	default:
		return "auto"
	}
}

func mtClassicAPIPort(useTLS bool) int {
	keys := []string{"MIKROTIK_API_PORT", "MIKROTIK_CLASSIC_API_PORT"}
	def := 8728

	if useTLS {
		keys = []string{"MIKROTIK_API_SSL_PORT", "MIKROTIK_CLASSIC_API_SSL_PORT", "MIKROTIK_API_PORT", "MIKROTIK_CLASSIC_API_PORT"}
		def = 8729
	}

	for _, k := range keys {
		v := strings.TrimSpace(os.Getenv(k))
		if v == "" {
			continue
		}

		n, err := strconv.Atoi(v)
		if err == nil && n > 0 {
			return n
		}
	}

	return def
}

func cleanRouterHost(host string) string {
	host = strings.TrimSpace(host)
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")

	if i := strings.Index(host, "/"); i >= 0 {
		host = host[:i]
	}

	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}

	return strings.Trim(host, "[]")
}
