package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func envFromFile(path string) map[string]string {
	out := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		v = strings.Trim(v, `"`)
		out[k] = v
	}
	return out
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  vcr-demo-user list")
	fmt.Println("  vcr-demo-user create USER PASS ROLE ENABLED")
	fmt.Println("  vcr-demo-user update ID USER ROLE ENABLED")
	fmt.Println("  vcr-demo-user passwd ID PASS")
	fmt.Println("  vcr-demo-user enable ID")
	fmt.Println("  vcr-demo-user disable ID")
	fmt.Println("  vcr-demo-user delete ID")
	os.Exit(2)
}

func validRole(role string) string {
	role = strings.TrimSpace(role)
	switch role {
	case "superadmin", "admin", "operator", "router":
		return role
	default:
		return "superadmin"
	}
}

func validEnabled(v string) int {
	v = strings.TrimSpace(v)
	if v == "0" || strings.EqualFold(v, "false") || strings.EqualFold(v, "off") {
		return 0
	}
	return 1
}

func mustID(s string) int64 {
	id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil || id <= 0 {
		fmt.Println("ERROR: ID tidak valid")
		os.Exit(1)
	}
	return id
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	env := envFromFile("/etc/tuku-demo.env")
	dsn := strings.TrimSpace(env["DB_DSN"])
	if dsn == "" {
		fmt.Println("ERROR: DB_DSN tidak ketemu di /etc/tuku-demo.env")
		os.Exit(1)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("ERROR open DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("ERROR ping DB:", err)
		os.Exit(1)
	}

	action := os.Args[1]

	switch action {
	case "list":
		rows, err := db.Query(`
			SELECT id, username, COALESCE(role,'superadmin'), COALESCE(router_id,0), COALESCE(enabled,1), COALESCE(created_at,'')
			FROM users
			ORDER BY id ASC
		`)
		if err != nil {
			fmt.Println("ERROR list:", err)
			os.Exit(1)
		}
		defer rows.Close()

		fmt.Println("ID\tUSERNAME\tROLE\tROUTER_ID\tENABLED\tCREATED_AT")
		for rows.Next() {
			var id int64
			var username, role, created string
			var routerID, enabled int
			_ = rows.Scan(&id, &username, &role, &routerID, &enabled, &created)
			fmt.Printf("%d\t%s\t%s\t%d\t%d\t%s\n", id, username, role, routerID, enabled, created)
		}

	case "create":
		if len(os.Args) < 6 {
			usage()
		}
		username := strings.TrimSpace(os.Args[2])
		pass := os.Args[3]
		role := validRole(os.Args[4])
		enabled := validEnabled(os.Args[5])

		if username == "" || pass == "" {
			fmt.Println("ERROR: username/password wajib diisi")
			os.Exit(1)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("ERROR bcrypt:", err)
			os.Exit(1)
		}

		_, err = db.Exec(`INSERT INTO users(username,password_hash,role,router_id,enabled) VALUES(?,?,?,?,?)`, username, string(hash), role, nil, enabled)
		if err != nil {
			fmt.Println("ERROR create:", err)
			os.Exit(1)
		}
		fmt.Println("OK: user demo dibuat:", username)

	case "update":
		if len(os.Args) < 6 {
			usage()
		}
		id := mustID(os.Args[2])
		username := strings.TrimSpace(os.Args[3])
		role := validRole(os.Args[4])
		enabled := validEnabled(os.Args[5])

		if username == "" {
			fmt.Println("ERROR: username wajib diisi")
			os.Exit(1)
		}

		_, err := db.Exec(`UPDATE users SET username=?, role=?, enabled=? WHERE id=?`, username, role, enabled, id)
		if err != nil {
			fmt.Println("ERROR update:", err)
			os.Exit(1)
		}
		fmt.Println("OK: user demo diupdate:", id)

	case "passwd":
		if len(os.Args) < 4 {
			usage()
		}
		id := mustID(os.Args[2])
		pass := os.Args[3]

		if pass == "" {
			fmt.Println("ERROR: password kosong")
			os.Exit(1)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("ERROR bcrypt:", err)
			os.Exit(1)
		}

		_, err = db.Exec(`UPDATE users SET password_hash=? WHERE id=?`, string(hash), id)
		if err != nil {
			fmt.Println("ERROR password:", err)
			os.Exit(1)
		}
		fmt.Println("OK: password user demo diganti:", id)

	case "enable":
		if len(os.Args) < 3 {
			usage()
		}
		id := mustID(os.Args[2])
		_, err := db.Exec(`UPDATE users SET enabled=1 WHERE id=?`, id)
		if err != nil {
			fmt.Println("ERROR enable:", err)
			os.Exit(1)
		}
		fmt.Println("OK: user demo diaktifkan:", id)

	case "disable":
		if len(os.Args) < 3 {
			usage()
		}
		id := mustID(os.Args[2])
		_, err := db.Exec(`UPDATE users SET enabled=0 WHERE id=?`, id)
		if err != nil {
			fmt.Println("ERROR disable:", err)
			os.Exit(1)
		}
		fmt.Println("OK: user demo dinonaktifkan:", id)

	case "delete":
		if len(os.Args) < 3 {
			usage()
		}
		id := mustID(os.Args[2])
		_, err := db.Exec(`DELETE FROM users WHERE id=?`, id)
		if err != nil {
			fmt.Println("ERROR delete:", err)
			os.Exit(1)
		}
		fmt.Println("OK: user demo dihapus:", id)

	default:
		usage()
	}
}
