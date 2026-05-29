package models

import "time"

// Return is a customer-initiated reversal of part or all of a Sale. Each
// ReturnItem references an original SaleItem and records how many units came
// back. The original Sale is preserved untouched — returns are tracked as
// distinct events for audit and reporting.
//
// On successful processing, ReturnItem.Quantity is added back to Stock via a
// StockMovement of type MovementReturn, with reference_id pointing at the
// Return row.
type Return struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	BusinessID     uint    `gorm:"not null;index" json:"business_id"`
	SaleID         uint    `gorm:"not null;index" json:"sale_id"`
	BranchID       uint    `gorm:"not null;index" json:"branch_id"`
	RefundedTotal  float64 `gorm:"not null" json:"refunded_total"`
	PaymentMethod  string  `json:"payment_method"`  // cash | mobile_money | store_credit
	TransactionRef string  `json:"transaction_ref"` // for mobile-money refunds
	Reason         string  `gorm:"type:text" json:"reason"`
	ProcessedBy    uint    `gorm:"not null;index" json:"processed_by"`
	CreatedAt      time.Time `json:"created_at"`

	Sale      Sale         `gorm:"foreignKey:SaleID" json:"sale,omitempty"`
	Branch    Branch       `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Processor User         `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
	Items     []ReturnItem `gorm:"foreignKey:ReturnID" json:"items,omitempty"`
}

// ReturnItem is one line of a Return — the link back to the original
// SaleItem so we can compute "how much of this line has already been
// returned" when validating subsequent partial returns of the same sale.
type ReturnItem struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	ReturnID   uint    `gorm:"not null;index" json:"return_id"`
	SaleItemID uint    `gorm:"not null;index" json:"sale_item_id"`
	ProductID  uint    `gorm:"not null;index" json:"product_id"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	LineTotal  float64 `json:"line_total"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
