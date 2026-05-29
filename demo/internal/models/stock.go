package models

import "time"

type Stock struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ProductID uint    `gorm:"not null;uniqueIndex:idx_product_branch" json:"product_id"`
	BranchID  uint    `gorm:"not null;uniqueIndex:idx_product_branch" json:"branch_id"`
	Quantity  int     `gorm:"default:0" json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Branch    Branch  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
