package models

type SaleItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	SaleID    uint    `gorm:"not null;index" json:"sale_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"unit_price"`
	UnitCost  float64 `json:"unit_cost"`
	LineTotal float64 `json:"line_total"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
