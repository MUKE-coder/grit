package models

import "time"

// Movement type constants
const (
	MovementStockIn     = "stock_in"
	MovementSale        = "sale"
	MovementReturn      = "return" // customer return — adds stock back, reference_id = Return.ID
	MovementTransferOut = "transfer_out"
	MovementTransferIn  = "transfer_in"
	MovementAdjustment  = "adjustment"
)

type StockMovement struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProductID    uint      `gorm:"not null;index" json:"product_id"`
	BranchID     uint      `gorm:"not null;index" json:"branch_id"`
	MovementType string    `json:"movement_type"` // stock_in | sale | transfer_out | transfer_in | adjustment
	Quantity     int       `gorm:"not null" json:"quantity"` // positive = in, negative = out
	ReferenceID  *uint     `json:"reference_id"` // sale_id or transfer_id
	Note         string    `json:"note"`
	CreatedBy    uint      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	Product      Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Branch       Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
