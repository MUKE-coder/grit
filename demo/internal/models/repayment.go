package models

import "time"

// Repayment status constants. Cash collections start as PENDING (cashier records,
// manager verifies). DGateway-collected mobile-money payments auto-APPROVE on webhook.
const (
	RepaymentStatusPending  = "pending"
	RepaymentStatusApproved = "approved"
	RepaymentStatusRejected = "rejected"
	RepaymentStatusFailed   = "failed" // mobile-money attempt that didn't go through
)

// Payment method constants for repayments.
const (
	RepaymentMethodCash         = "cash"
	RepaymentMethodMobileMoney  = "mobile_money" // via DGateway
	RepaymentMethodBank         = "bank_transfer"
	RepaymentMethodCheque       = "cheque"
)

// RepaymentSchedule is one row per planned installment. Generated when a Loan is
// created; never deleted. Updated as Repayments are applied.
type RepaymentSchedule struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	LoanID            uint       `gorm:"not null;index;uniqueIndex:idx_loan_installment" json:"loan_id"`
	InstallmentNumber int        `gorm:"not null;uniqueIndex:idx_loan_installment" json:"installment_number"`
	DueDate           time.Time  `gorm:"not null;index" json:"due_date"`
	PrincipalAmount   float64    `json:"principal_amount"`
	InterestAmount    float64    `json:"interest_amount"`
	TotalAmount       float64    `json:"total_amount"` // principal + interest for this installment
	BalanceBefore     float64    `json:"balance_before"`
	BalanceAfter      float64    `json:"balance_after"`
	IsPaid            bool       `gorm:"default:false;index" json:"is_paid"`
	PaidAmount        float64    `gorm:"default:0" json:"paid_amount"`
	PaidAt            *time.Time `json:"paid_at"`
	IsOverdue         bool       `gorm:"default:false;index" json:"is_overdue"`
	DaysPastDue       int        `gorm:"default:0" json:"days_past_due"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Repayment is one actual money-in event applied against a Loan. May span
// multiple schedule rows (rare) or be a partial payment of one row.
type Repayment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	LoanID         uint      `gorm:"not null;index" json:"loan_id"`
	ScheduleID     *uint     `gorm:"index" json:"schedule_id"` // nil = unallocated/across-multiple
	Amount         float64   `gorm:"not null" json:"amount"`
	CollectionDate time.Time `gorm:"not null;index" json:"collection_date"`
	PaymentMethod  string    `gorm:"not null" json:"payment_method"` // cash | mobile_money | bank_transfer | cheque
	TransactionRef string    `gorm:"index" json:"transaction_ref"`   // bank ref, MoMo txn id, cheque #
	Receipt        string    `json:"receipt"`                        // receipt # we issue
	Status         string    `gorm:"default:pending;index" json:"status"`
	Notes          string    `gorm:"type:text" json:"notes"`

	// DGateway integration fields (mobile money via Desispay)
	DGatewayReference  string     `gorm:"index" json:"dgateway_reference"` // dgw_xxxxx
	DGatewayProvider   string     `json:"dgateway_provider"`               // iotec | relworx
	DGatewayPhone      string     `json:"dgateway_phone"`
	DGatewayFee        float64    `json:"dgateway_fee"`     // platform commission deducted
	DGatewayNetAmount  float64    `json:"dgateway_net_amount"`
	DGatewayFailReason string     `json:"dgateway_fail_reason"`
	DGatewayConfirmedAt *time.Time `json:"dgateway_confirmed_at"`

	// Cashier/verifier audit
	CollectedBy uint       `gorm:"not null;index" json:"collected_by"`
	VerifiedBy  *uint      `json:"verified_by"`
	VerifiedAt  *time.Time `json:"verified_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Loan      Loan               `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
	Schedule  *RepaymentSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Collector User               `gorm:"foreignKey:CollectedBy" json:"collector,omitempty"`
}
