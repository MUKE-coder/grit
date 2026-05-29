package models

import "time"

type Product struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	BusinessID        uint      `gorm:"not null;index" json:"business_id"`
	CategoryID        uint      `gorm:"not null;index" json:"category_id"`
	Category          Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Title             string    `gorm:"not null" json:"title"`
	Description       string    `gorm:"type:text" json:"description"`
	Barcode           string    `gorm:"index" json:"barcode"`
	SellingPrice      float64   `gorm:"not null" json:"selling_price"`
	CostPrice         float64   `json:"cost_price"`
	ImageKey          string    `json:"-"` // R2 key — never exposed directly
	LowStockThreshold int       `gorm:"default:5" json:"low_stock_threshold"`
	CreatedBy         uint      `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
