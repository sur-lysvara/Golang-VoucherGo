package main

import "database/sql"

type Router struct {
	ID              int64
	Name            string
	Slug            string
	HotspotLoginURL string
	Host            string
	Port            int
	APIMode         string
	Username        string
	Password        string
	UseSSL          bool
	InsecureSSL     bool
	Enabled         bool
	LastTestStatus  string
	LastTestMessage string
	LastTestAt      string
	LastTestOK      bool
	LastTestFail    bool
}

type Package struct {
	ID          int64
	RouterID    int64
	RouterName  string
	Name        string
	Price       int64
	Profile     string
	LimitUptime string
	SharedUsers int
	Active      bool
}

type Tx struct {
	ID               int64
	OrderID          string
	RouterName       string
	PackageName      string
	CustomerName     string
	CustomerPhone    string
	QRISString       string
	QRISTotalPayment int64
	QRISExpiredAt    string
	Amount           int64
	Status           string
	PaymentMethod    string
	VoucherUsername  string
	VoucherPassword  string
	ProvisionStatus  string
	ProvisionError   string
	WAStatus         string
	WAError          string
	WASentAt         string
	CreatedAt        string
	PaidAt           sql.NullString
}

type AdminUserRow struct {
	ID         int64
	Username   string
	Role       string
	RouterID   int64
	RouterName string
	RouterSlug string
	Enabled    int
	CreatedAt  string
}

type AdminRouterOption struct {
	ID   int64
	Name string
	Slug string
}

type PakasirWebhookV2 struct {
	TxnID       string `json:"txn_id"`
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	IsSandbox   bool   `json:"is_sandbox"`
	Status      string `json:"status"`
	CompletedAt string `json:"completed_at"`
}

type PakasirV2Status struct {
	TxnID       string `json:"txn_id"`
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	IsSandbox   bool   `json:"is_sandbox"`
	Status      string `json:"status"`
	CompletedAt string `json:"completed_at"`
}
