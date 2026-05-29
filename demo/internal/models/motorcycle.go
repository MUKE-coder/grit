package models

import (
	"time"

	"gorm.io/gorm"
)

// Motorcycle status constants.
const (
	MotorcycleStatusAvailable  = "available"
	MotorcycleStatusReserved   = "reserved"   // Pending loan approval / pending cash sale checkout
	MotorcycleStatusSold       = "sold"       // Sold via cash sale or fully paid loan
	MotorcycleStatusOnLoan     = "on_loan"    // Loan disbursed, balance still outstanding
	MotorcycleStatusRepossessed = "repossessed"
)

// Motorcycle represents a single physical motorcycle in inventory, identified by
// its number plate. Unlike spare-parts Products (quantity-tracked stock), each
// Motorcycle row IS one unit — there is no quantity field.
type Motorcycle struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	BusinessID   uint           `gorm:"not null;index;uniqueIndex:idx_business_plate" json:"business_id"`
	BranchID     uint           `gorm:"not null;index" json:"branch_id"`
	Name         string         `gorm:"not null" json:"name"`         // brand/model, e.g. "KEVLA"
	NumberPlate  string         `gorm:"not null;uniqueIndex:idx_business_plate" json:"number_plate"`
	ChassisNo    string         `json:"chassis_no"`
	EngineNo     string         `json:"engine_no"`
	Color        string         `json:"color"`
	YearOfMake   int            `json:"year_of_make"`
	CostPrice    float64        `json:"cost_price"`
	SellingPrice float64        `gorm:"not null" json:"selling_price"`
	ImageKey     string         `json:"-"` // R2 key
	Status       string         `gorm:"default:available;index" json:"status"`
	Notes        string         `gorm:"type:text" json:"notes"`
	CreatedBy    uint           `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Branch       Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
