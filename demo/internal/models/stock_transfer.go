package models

import "time"

type StockTransfer struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BusinessID    uint      `gorm:"not null;index" json:"business_id"`
	FromBranchID  uint      `gorm:"not null" json:"from_branch_id"`
	ToBranchID    uint      `gorm:"not null" json:"to_branch_id"`
	ProductID     uint      `gorm:"not null" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Note          string    `json:"note"`
	TransferredBy uint      `json:"transferred_by"`
	CreatedAt     time.Time `json:"created_at"`
	Product       Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	FromBranch    Branch    `gorm:"foreignKey:FromBranchID" json:"from_branch,omitempty"`
	ToBranch      Branch    `gorm:"foreignKey:ToBranchID" json:"to_branch,omitempty"`
}
