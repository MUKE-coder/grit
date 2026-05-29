package models

import "time"

// CashSale is a direct cash sale of a Motorcycle. Distinct from spares Sale
// (which is a multi-line POS transaction) — one CashSale = one motorcycle.
// Loan-financed motorcycle sales go through Loan, not CashSale.
type CashSale struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	BusinessID     uint    `gorm:"not null;index" json:"business_id"`
	BranchID       uint    `gorm:"not null;index" json:"branch_id"`
	MotorcycleID   uint    `gorm:"not null;uniqueIndex" json:"motorcycle_id"` // a motorcycle can only be sold once
	SaleNumber     string  `gorm:"uniqueIndex" json:"sale_number"`            // KM-MC-000123
	SoldBy         uint    `gorm:"not null;index" json:"sold_by"`

	// Customer details (we don't always have a Borrower record for cash buyers)
	CustomerName    string `gorm:"not null" json:"customer_name"`
	CustomerPhone   string `gorm:"index" json:"customer_phone"`
	CustomerEmail   string `json:"customer_email"`
	CustomerNIN     string `json:"customer_nin"`
	CustomerAddress string `json:"customer_address"`

	// Money
	ListPrice      float64 `json:"list_price"`
	DiscountAmount float64 `gorm:"default:0" json:"discount_amount"`
	Total          float64 `gorm:"not null" json:"total"`
	PaymentMethod  string  `gorm:"not null" json:"payment_method"` // cash | mobile_money | bank_transfer
	TransactionRef string  `json:"transaction_ref"`

	// DGateway fields (if mobile money)
	DGatewayReference string `gorm:"index" json:"dgateway_reference"`

	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Motorcycle Motorcycle `gorm:"foreignKey:MotorcycleID" json:"motorcycle,omitempty"`
	Branch     Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Seller     User       `gorm:"foreignKey:SoldBy" json:"seller,omitempty"`
}
