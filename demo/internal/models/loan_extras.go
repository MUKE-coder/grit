package models

import "time"

// LoanFee constants.
const (
	FeeTypeProcessing  = "processing"
	FeeTypeInsurance   = "insurance"
	FeeTypeRegistration = "registration"
	FeeTypeOther       = "other"
)

// LoanFee is a one-off charge on a Loan (processing, insurance, etc).
type LoanFee struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	LoanID    uint       `gorm:"not null;index" json:"loan_id"`
	Type      string     `gorm:"not null" json:"type"`
	Amount    float64    `gorm:"not null" json:"amount"`
	IsPaid    bool       `gorm:"default:false" json:"is_paid"`
	PaidAt    *time.Time `json:"paid_at"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
}

// LoanPenalty is a late-payment fine attached to a missed schedule installment.
type LoanPenalty struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	LoanID     uint       `gorm:"not null;index" json:"loan_id"`
	ScheduleID *uint      `gorm:"index" json:"schedule_id"`
	Amount     float64    `gorm:"not null" json:"amount"`
	Reason     string     `json:"reason"`
	IsPaid     bool       `gorm:"default:false" json:"is_paid"`
	PaidAt     *time.Time `json:"paid_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// LoanCollateral types.
const (
	CollateralLandTitle      = "land_title"
	CollateralVehicleLogbook = "vehicle_logbook"
	CollateralGuarantor      = "guarantor"
	CollateralBusinessAsset  = "business_asset"
	CollateralOther          = "other"
)

// LoanCollateral is security pledged against a loan. The motorcycle being
// financed is implicit collateral via Loan.MotorcycleID — this table is for
// additional collateral (land titles, guarantors, etc).
type LoanCollateral struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	LoanID        uint      `gorm:"not null;index" json:"loan_id"`
	Type          string    `gorm:"not null" json:"type"`
	Description   string    `gorm:"type:text" json:"description"`
	EstimatedValue float64  `json:"estimated_value"`
	DocumentKey   string    `json:"-"` // R2 key for collateral document scan
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
