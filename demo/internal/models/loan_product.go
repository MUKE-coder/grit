package models

import "time"

// Interest method constants.
const (
	InterestMethodFlat             = "flat"              // Equal interest split across all installments
	InterestMethodReducingBalance  = "reducing_balance"  // Monthly rate applied to outstanding principal
)

// Repayment cycle constants.
const (
	RepaymentCycleWeekly   = "weekly"
	RepaymentCycleBiweekly = "biweekly"
	RepaymentCycleMonthly  = "monthly"
)

// LoanProduct is a loan template — the rules a Loan is created from.
// e.g. "Boda Boda 12-month loan: 200K–5M, 24% flat, weekly".
type LoanProduct struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	BusinessID        uint      `gorm:"not null;index" json:"business_id"`
	Name              string    `gorm:"not null" json:"name"`
	Description       string    `gorm:"type:text" json:"description"`
	MinAmount         float64   `gorm:"not null" json:"min_amount"`
	MaxAmount         float64   `gorm:"not null" json:"max_amount"`
	MinDuration       int       `gorm:"not null" json:"min_duration"` // in repayment cycles
	MaxDuration       int       `gorm:"not null" json:"max_duration"`
	InterestMethod    string    `gorm:"not null" json:"interest_method"` // flat | reducing_balance
	InterestRate      float64   `gorm:"not null" json:"interest_rate"`   // annual %
	RepaymentCycle    string    `gorm:"not null" json:"repayment_cycle"` // weekly | biweekly | monthly
	RequiresCollateral bool     `gorm:"default:false" json:"requires_collateral"`
	GracePeriodDays   int       `gorm:"default:0" json:"grace_period_days"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
