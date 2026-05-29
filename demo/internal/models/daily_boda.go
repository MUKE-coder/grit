package models

import (
	"time"

	"gorm.io/gorm"
)

// Daily Boda — separate motorcycle/driver pool for daily-rental income (NOT
// loan inventory). Drivers pay a fixed daily rate to use the bike.

// DailyBodaMotorcycle status.
const (
	DailyBodaStatusAvailable = "available"
	DailyBodaStatusOccupied  = "occupied"
	DailyBodaStatusReturned  = "returned"   // out of rotation, parked
	DailyBodaStatusInService = "in_service" // mechanic
)

// DailyBodaPayment status.
const (
	DailyBodaPaymentPaid    = "paid"
	DailyBodaPaymentPartial = "partial"
	DailyBodaPaymentPending = "pending"
)

// Default daily rental rate (UGX). Override per-driver via DailyBodaDriver.DailyRate.
const DefaultDailyBodaRate = 15000.0

// DailyBodaMotorcycle is a motorcycle in the daily-rental fleet (separate from
// inventory Motorcycle which is loan/cash-sale inventory).
type DailyBodaMotorcycle struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	BusinessID        uint           `gorm:"not null;index;uniqueIndex:idx_dboda_plate" json:"business_id"`
	BranchID          uint           `gorm:"not null;index" json:"branch_id"`
	Name              string         `gorm:"not null" json:"name"`
	NumberPlate       string         `gorm:"not null;uniqueIndex:idx_dboda_plate" json:"number_plate"`
	Status            string         `gorm:"default:available;index" json:"status"`
	AssignedDriverID  *uint          `gorm:"uniqueIndex" json:"assigned_driver_id"` // one driver at a time
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	Branch            Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// DailyBodaDriver is a registered rider in the daily-rental program.
type DailyBodaDriver struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	BusinessID uint           `gorm:"not null;index;uniqueIndex:idx_dboda_driver_phone" json:"business_id"`
	BranchID   uint           `gorm:"not null;index" json:"branch_id"`
	FullName   string         `gorm:"not null" json:"full_name"`
	Phone      string         `gorm:"not null;uniqueIndex:idx_dboda_driver_phone" json:"phone"`
	NationalID string         `json:"national_id"`
	Address    string         `gorm:"type:text" json:"address"`
	DailyRate  float64        `gorm:"default:15000" json:"daily_rate"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	PhotoKey   string         `json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Branch     Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// DailyBodaPayment records one driver's payment for one day. Balance lets us
// track partial payments — driver may pay 10K today, owes 5K of the 15K daily.
type DailyBodaPayment struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	BusinessID    uint       `gorm:"not null;index" json:"business_id"`
	DriverID      uint       `gorm:"not null;index" json:"driver_id"`
	MotorcycleID  uint       `gorm:"not null;index" json:"motorcycle_id"`
	BranchID      uint       `gorm:"not null;index" json:"branch_id"`
	Amount        float64    `gorm:"not null" json:"amount"`        // paid this transaction
	DailyRate     float64    `gorm:"not null" json:"daily_rate"`    // expected for this day
	Balance       float64    `json:"balance"`                       // remaining for the day after this payment
	PaymentDate   time.Time  `gorm:"not null;index" json:"payment_date"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `gorm:"default:pending;index" json:"status"`
	CollectedBy   uint       `json:"collected_by"`
	VerifiedBy    *uint      `json:"verified_by"`
	VerifiedAt    *time.Time `json:"verified_at"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Driver     DailyBodaDriver     `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
	Motorcycle DailyBodaMotorcycle `gorm:"foreignKey:MotorcycleID" json:"motorcycle,omitempty"`
	Branch     Branch              `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
