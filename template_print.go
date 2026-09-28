package main

import (
	"fmt"
	"net/http"

	"log"
)

func (a *App) transactionsPrint(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query(`SELECT
			t.order_id,
			r.name,
			p.name,
			COALESCE(t.customer_name,''),
			COALESCE(t.customer_phone,''),
			t.amount,
			t.status,
			COALESCE(t.voucher_username,''),
			t.provision_status,
			DATE_FORMAT(t.created_at,'%Y-%m-%d %H:%i')
		FROM transactions t
		JOIN routers r ON r.id=t.router_id
		JOIN packages p ON p.id=t.package_id
		ORDER BY t.id DESC
		LIMIT 1000`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type Row struct {
		OrderID string
		Router  string
		Paket   string
		Name    string
		Phone   string
		Amount  int64
		Status  string
		Voucher string
		Prov    string
		Created string
	}

	var list []Row

	for rows.Next() {
		var x Row
		if err := rows.Scan(&x.OrderID, &x.Router, &x.Paket, &x.Name, &x.Phone, &x.Amount, &x.Status, &x.Voucher, &x.Prov, &x.Created); err != nil {
			log.Printf("print tx scan error: %v", err)
			continue
		}
		list = append(list, x)
	}

	body := renderHTML(`
<h2>Data Transaksi VoucherGo</h2>
<table>
<tr>
	<th>Order</th>
	<th>Tanggal</th>
	<th>Router</th>
	<th>Paket</th>
	<th>Pembeli</th>
	<th>Total</th>
	<th>Status</th>
	<th>Voucher</th>
	<th>Provision</th>
</tr>
{{range .}}
<tr>
	<td>{{.OrderID}}</td>
	<td>{{.Created}}</td>
	<td>{{.Router}}</td>
	<td>{{.Paket}}</td>
	<td>{{.Name}}<br>{{.Phone}}</td>
	<td>{{rp .Amount}}</td>
	<td>{{.Status}}</td>
	<td><b>{{.Voucher}}</b></td>
	<td>{{.Prov}}</td>
</tr>
{{end}}`, list)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="id">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Print Transaksi</title>
<style>
*{box-sizing:border-box}
body{
  margin:0;
  padding:20px;
  font-family:Arial,sans-serif;
  color:#111827;
  background:#ffffff;
}
.print-actions{
  display:flex;
  gap:8px;
  margin:0 0 14px;
}
button{
  border:0;
  border-radius:10px;
  padding:9px 14px;
  background:#2563eb;
  color:#ffffff;
  font-weight:800;
  cursor:pointer;
}
h2{
  margin:0 0 14px;
  font-size:20px;
}
table{
  width:100%%;
  border-collapse:collapse;
}
th,td{
  border:1px solid #d1d5db;
  padding:7px 8px;
  font-size:12px;
  text-align:left;
  vertical-align:top;
}
th{
  background:#f3f4f6;
  font-weight:900;
}
td b{
  font-family:ui-monospace,SFMono-Regular,Menlo,monospace;
  font-size:12px;
}
@media print{
  body{padding:0}
  .print-actions{display:none!important}
  h2{font-size:18px}
  th,td{font-size:10.5px;padding:5px 6px}
}
</style>
</head>
<body>
<div class="print-actions">
  <button type="button" onclick="window.print()">Print</button>
  <button type="button" onclick="history.back()">Kembali</button>
</div>
%s
<script>
window.addEventListener("load", function(){
  setTimeout(function(){ window.print(); }, 250);
});
</script>
</body>
</html>`, body)
}
