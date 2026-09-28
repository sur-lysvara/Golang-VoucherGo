package main

import (
	"bytes"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

func renderHTML(tpl string, data any) template.HTML {
	t := template.Must(template.New("body").Funcs(template.FuncMap{
		"rp": func(v int64) string {
			return rupiah(v)
		},
		"checked": func(b bool) template.HTML {
			if b {
				return `checked`
			}
			return ``
		},
		"datetimeLocal": datetimeLocal,
	}).Parse(tpl))

	var b bytes.Buffer
	_ = t.Execute(&b, data)
	return template.HTML(b.String())
}

func findMatchingDivEndHTML(h string, start int) int {
	if start < 0 || start >= len(h) {
		return -1
	}

	depth := 0
	i := start

	for i < len(h) {
		nextOpen := strings.Index(h[i:], "<div")
		nextClose := strings.Index(h[i:], "</div>")

		if nextClose < 0 {
			return -1
		}

		if nextOpen >= 0 && nextOpen < nextClose {
			depth++
			i += nextOpen + len("<div")
			continue
		}

		depth--
		i += nextClose + len("</div>")

		if depth == 0 {
			return i
		}
	}

	return -1
}

func polishInvoiceHTML(body template.HTML) template.HTML {
	h := string(body)

	h = strings.ReplaceAll(h, "Gunakan kode voucher sebagai username dan password login hotspot.", "Salin kode voucher dan masukan di halaman login WiFi")
	h = strings.ReplaceAll(h, "Gunakan kode voucher sebagai username dan password login hotspot", "Salin kode voucher dan masukan di halaman login WiFi")
	h = strings.ReplaceAll(h, "Gunakan username dan password tersebut untuk login hotspot.", "Salin kode voucher dan masukan di halaman login WiFi")

	h = strings.ReplaceAll(h, "Password sama dengan kode voucher.", "")
	h = strings.ReplaceAll(h, "Password sama dengan kode voucher", "")
	h = strings.ReplaceAll(h, "Password sama dengan kode voucher.", "")

	h = strings.ReplaceAll(h, "Copy Voucher", "Salin Voucher")
	h = strings.ReplaceAll(h, ">Copy<", ">Salin<")

	h = strings.ReplaceAll(h, "TUKU", "VCR")
	h = strings.ReplaceAll(h, ">Total<", ">Nominal Pembelian<")
	h = strings.ReplaceAll(h, ">Pembeli<", ">Data Pengguna<")
	h = strings.ReplaceAll(h, ">Nama Pengguna<", ">Data Pengguna<")
	h = strings.ReplaceAll(h, ">Provision<", ">Status<")

	h = strings.ReplaceAll(h, `<span class="badge paid">paid</span>`, `<span class="invoice-status-lunas">Lunas</span>`)
	h = strings.ReplaceAll(h, `<span class="badge paid">dibayar</span>`, `<span class="invoice-status-lunas">Lunas</span>`)
	h = strings.ReplaceAll(h, `class="badge paid">paid`, `class="invoice-status-lunas">Lunas`)
	h = strings.ReplaceAll(h, `class="badge paid">dibayar`, `class="invoice-status-lunas">Lunas`)
	h = strings.ReplaceAll(h, ">paid<", ">Lunas<")
	h = strings.ReplaceAll(h, ">dibayar<", ">Lunas<")

	// Clipboard sebelumnya berisi "Username: xxxx Password: xxxx".
	// Paksa jadi kode voucher saja.
	reCopy := regexp.MustCompile(`Username:\s*([A-Za-z0-9]+)\s*Password:\s*([A-Za-z0-9]+)`)
	h = reCopy.ReplaceAllString(h, `$1`)

	// Replace instruction text even if username/password are wrapped in HTML tags.
	reLoginInstruction := regexp.MustCompile(`(?is)Gunakan\s+kode\s+voucher\s+sebagai\s*(?:<[^>]+>)*\s*username\s*(?:</[^>]+>)*\s*dan\s*(?:<[^>]+>)*\s*password\s*(?:</[^>]+>)*\s*login\s+hotspot\.?`)
	h = reLoginInstruction.ReplaceAllString(h, "Salin kode voucher dan masukan di halaman login WiFi")

	if !strings.Contains(h, "invoice-login-btn") {
		loginBtn := `<a class="btn light invoice-login-btn" href="/go-login">Masuk Login WiFi</a>`
		rePrintBtn := regexp.MustCompile(`(?is)(<[^>]*(?:button|a)[^>]*>\s*Print\s*</(?:button|a)>)`)
		h = rePrintBtn.ReplaceAllString(h, `$1`+loginBtn)
	}

	h = moveInvoiceVoucherBox(h)

	return template.HTML(h)
}

func rupiah(v int64) string {
	s := strconv.FormatInt(v, 10)
	var out []byte

	for i, c := range reverse(s) {
		if i > 0 && i%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, byte(c))
	}

	return "Rp" + reverse(string(out))
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
