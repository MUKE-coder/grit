package models

import "time"

type Sale struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	BranchID       uint       `gorm:"not null;index" json:"branch_id"`
	CashierID      uint       `gorm:"not null;index" json:"cashier_id"`
	Segment        string     `gorm:"default:spares;index" json:"segment"` // spares (POS today)
	PaymentMethod  string     `json:"payment_method"`                      // cash | mobile_money
	CustomerPhone  string     `json:"customer_phone"`
	Subtotal       float64    `json:"subtotal"`
	DiscountAmount float64    `gorm:"default:0" json:"discount_amount"`
	Total          float64    `json:"total"`
	Items          []SaleItem `gorm:"foreignKey:SaleID" json:"items,omitempty"`
	Branch         Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Cashier        User       `gorm:"foreignKey:CashierID" json:"cashier,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
