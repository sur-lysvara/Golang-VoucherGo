package main

import (
	"database/sql"
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"

	"log"
)

func (a *App) transactions(w http.ResponseWriter, r *http.Request) {
	// ADMIN_TRANSACTIONS_SCOPE_V1
	// ADMIN_TX_RESPONSIVE_V1
	role, scopedRouterID := a.currentAdminRole(r)

	var rows *sql.Rows
	var err error

	if role == "router" && scopedRouterID > 0 {
		rows, err = a.db.Query(`SELECT
			t.id,
			t.order_id,
			r.name,
			p.name,
			t.customer_name,
			t.customer_phone,
			t.amount,
			t.status,
			COALESCE(t.payment_method,''),
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,''),
			t.provision_status,
			COALESCE(t.provision_error,''),
			COALESCE(t.wa_status,'pending'),
			COALESCE(t.wa_error,''),
			COALESCE(DATE_FORMAT(t.wa_sent_at,'%Y-%m-%d %H:%i'),''),
			DATE_FORMAT(t.created_at,'%Y-%m-%d %H:%i'),
			DATE_FORMAT(t.paid_at,'%Y-%m-%d %H:%i')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		JOIN packages p ON p.id=t.package_id
		WHERE t.router_id=?
		ORDER BY t.id DESC
		LIMIT 200`, scopedRouterID)
	} else {
		rows, err = a.db.Query(`SELECT
			t.id,
			t.order_id,
			r.name,
			p.name,
			t.customer_name,
			t.customer_phone,
			t.amount,
			t.status,
			COALESCE(t.payment_method,''),
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,''),
			t.provision_status,
			COALESCE(t.provision_error,''),
			COALESCE(t.wa_status,'pending'),
			COALESCE(t.wa_error,''),
			COALESCE(DATE_FORMAT(t.wa_sent_at,'%Y-%m-%d %H:%i'),''),
			DATE_FORMAT(t.created_at,'%Y-%m-%d %H:%i'),
			DATE_FORMAT(t.paid_at,'%Y-%m-%d %H:%i')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		JOIN packages p ON p.id=t.package_id
		ORDER BY t.id DESC
		LIMIT 200`)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var list []Tx

	for rows.Next() {
		var x Tx

		if err := rows.Scan(
			&x.ID,
			&x.OrderID,
			&x.RouterName,
			&x.PackageName,
			&x.CustomerName,
			&x.CustomerPhone,
			&x.Amount,
			&x.Status,
			&x.PaymentMethod,
			&x.VoucherUsername,
			&x.VoucherPassword,
			&x.ProvisionStatus,
			&x.ProvisionError,
			&x.WAStatus,
			&x.WAError,
			&x.WASentAt,
			&x.CreatedAt,
			&x.PaidAt,
		); err != nil {
			log.Printf("tx scan error: %v", err)
			continue
		}

		list = append(list, x)
	}

	body := renderHTML(`
<div class="card transactions-page">
	<div class="tx-toolbar">
		<h2 style="margin:0">Transaksi</h2>
	</div>

	<form id="txBulkMarkPaidForm" method="post" action="/admin/transactions/mark-paid"></form>
<form id="txBulkProvisionForm" method="post" action="/admin/transactions/provision"></form>
	<form id="txBulkWAForm" method="post" action="/admin/transactions/send-wa"></form>
	<form id="txBulkDeleteForm" method="post" action="/admin/transactions/delete"></form>

	<div class="tx-static-filter tx-static-filter-flat" id="txStaticFilter">
		<div class="field tx-search-field">
			<label>Search</label>
			<input type="search" id="txSearch" placeholder="Cari data...">
		</div>

		<div class="field">
			<label>Dari Tanggal</label>
			<input type="date" id="txDateFrom">
		</div>

		<div class="field">
			<label>Sampai Tanggal</label>
			<input type="date" id="txDateTo">
		</div>

		<div class="field tx-action-field">
			<label>&nbsp;</label>
			<button class="btn" type="button" id="txResetFilter">Reset</button>
		</div>

		<div class="field tx-action-field">
			<label>&nbsp;</label>
			<button class="btn pri" type="button" id="txExportFilter">Export Filter</button>
		</div>

		<div class="field tx-action-field">
			<label>&nbsp;</label>
			<button class="btn" type="button" id="txPrintFilter">Print Filter</button>
		</div>

		<div class="field tx-action-field">
			<label>&nbsp;</label>
			<button class="btn light" type="submit" form="txBulkProvisionForm" onclick="return confirm('Push ulang transaksi terpilih?')">Push MikroTik</button>
		</div>

		<div class="field tx-action-field">
			<label>&nbsp;</label>
			<button class="btn ok" type="submit" form="txBulkWAForm" onclick="return confirm('Kirim WhatsApp untuk transaksi terpilih?')">Kirim WA</button>
		</div>

				<div class="field tx-action-field">
					<label>&nbsp;</label>
					<button class="btn" type="submit" form="txBulkDeleteForm" formaction="/admin/transactions/send-email" formmethod="post" onclick="return confirm('Kirim Email voucher untuk transaksi terpilih?')">Kirim Email</button>
				</div>

		<div class="field tx-action-field">
			<label>&nbsp;</label>
			<button class="btn danger user-delete-btn" type="submit" form="txBulkDeleteForm" onclick="return confirm('Hapus transaksi terpilih? Data yang dihapus tidak bisa dikembalikan.')">Hapus Terpilih</button>
		</div>
	</div>
	<div class="vg-filter-meta" id="txFilterMeta"></div>

	<div class="tablewrap">
	<table class="transactions-compact-table">
		<tr>
			<th>Order</th>
			<th>Router/Paket</th>
			<th>Pembeli</th>
			<th>Total</th>
			<th>Status</th>
			<th>Voucher</th>
			<th>Provision</th>
			<th>WA</th>
			<th>Hapus</th>
		</tr>

		{{range .}}
		<tr>
			<td>
					<div class="tx-order-wrap">
						<input type="checkbox" name="id" value="{{.ID}}" form="txBulkDeleteForm" data-tx-id="{{.ID}}">
						<div class="tx-order-text"><b>{{.OrderID}}</b><br><small>{{.CreatedAt}}</small></div>
					</div>
				</td>
			<td>{{.RouterName}}<br><span class="muted">{{.PackageName}}</span></td>
			<td>{{.CustomerName}}<br><span class="muted">{{.CustomerPhone}}</span></td>
			<td>{{rp .Amount}}</td>
			<td><span class="badge {{.Status}}">{{.Status}}</span></td>
			<td>
				<b style="font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:18px">{{.VoucherUsername}}</b><br>
				<span class="muted">{{.VoucherPassword}}</span>
			</td>
			<td><span class="badge {{.ProvisionStatus}}">{{.ProvisionStatus}}</span><br><small>{{.ProvisionError}}</small></td>
			<td>
				<span class="badge {{.WAStatus}}">{{.WAStatus}}</span>
				{{if .WASentAt}}<br><small>{{.WASentAt}}</small>{{end}}
				{{if .WAError}}<br><small>{{.WAError}}</small>{{end}}
			</td>
			<td>
				<div class="actions tx-row-trash-only">
					<form method="post" action="/admin/transactions/delete" onsubmit="return confirm('Hapus transaksi ini?')">
										<input type="hidden" name="id" value="{{.ID}}">
										<button class="tx-trash-icon" title="Hapus transaksi" aria-label="Hapus transaksi">Hapus</button>
									</form>
				</div>
			</td>
		</tr>
		{{end}}
	</table>
	</div>
</div>`, list)

	page(w, "Transaksi", body)
}

func (a *App) provisionManual(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()

	ids := r.Form["id"]
	for _, raw := range ids {
		id, _ := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if id <= 0 {
			continue
		}

		// ADMIN_PROVISION_SCOPE_V1
		if !a.adminCanAccessTransaction(r, id) {
			continue
		}

		_ = a.provisionVoucher(id)
	}

	http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
}

func (a *App) transactionMarkPaid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()
	ids := r.Form["tx_id"]
	if len(ids) == 0 {
		http.Redirect(w, r, "/admin/transactions?msg=pilih_transaksi_dulu", http.StatusSeeOther)
		return
	}

	okCount := 0
	for _, raw := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil || id <= 0 {
			continue
		}

		if !a.adminCanAccessTransaction(r, id) {
			continue
		}

		// ADMIN_TX_MANUAL_PAID_SCOPE_V1
		res, err := a.db.Exec(`UPDATE transactions
			SET status='paid'
			WHERE id=?`, id)
		if err != nil {
			continue
		}

		if n, _ := res.RowsAffected(); n > 0 {
			okCount++
		}
	}

	_ = okCount
	http.Redirect(w, r, "/admin/transactions?msg=tandai_lunas_berhasil", http.StatusSeeOther)
}

func (a *App) transactionDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()

	ids := r.Form["id"]
	if len(ids) == 0 {
		if single := strings.TrimSpace(r.FormValue("id")); single != "" {
			ids = []string{single}
		}
	}

	for _, raw := range ids {
		id, _ := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if id <= 0 {
			continue
		}

		if !a.adminCanAccessTransaction(r, id) {
			continue
		}

		_, _ = a.db.Exec(`DELETE FROM transactions WHERE id=?`, id)
	}

	http.Redirect(w, r, "/admin/transactions", http.StatusSeeOther)
}

func (a *App) transactionsExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=tuku-transaksi.csv")

	cw := csv.NewWriter(w)
	defer cw.Flush()

	_ = cw.Write([]string{
		"order_id",
		"router",
		"paket",
		"customer_name",
		"customer_phone",
		"amount",
		"status",
		"voucher_username",
		"voucher_password",
		"provision_status",
		"provision_error",
		"created_at",
		"paid_at",
	})

	rows, err := a.db.Query(`SELECT
			t.order_id,
			r.name,
			p.name,
			COALESCE(t.customer_name,''),
			COALESCE(t.customer_phone,''),
			t.amount,
			t.status,
			COALESCE(t.voucher_username,''),
			COALESCE(t.voucher_password,''),
			t.provision_status,
			COALESCE(t.provision_error,''),
			DATE_FORMAT(t.created_at,'%Y-%m-%d %H:%i:%s'),
			COALESCE(DATE_FORMAT(t.paid_at,'%Y-%m-%d %H:%i:%s'),'')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		JOIN packages p ON p.id=t.package_id
		ORDER BY t.id DESC
		LIMIT 5000`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var orderID, router, paket, name, phone, status, vu, vp, ps, pe, created, paid string
		var amount int64

		if err := rows.Scan(&orderID, &router, &paket, &name, &phone, &amount, &status, &vu, &vp, &ps, &pe, &created, &paid); err != nil {
			log.Printf("tx export scan error: %v", err)
			continue
		}

		_ = cw.Write([]string{
			orderID,
			router,
			paket,
			name,
			phone,
			strconv.FormatInt(amount, 10),
			status,
			vu,
			vp,
			ps,
			pe,
			created,
			paid,
		})
	}
}
