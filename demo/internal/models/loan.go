package models

import "time"

// Loan status constants.
const (
	LoanStatusPending   = "pending"   // Created, not yet approved
	LoanStatusApproved  = "approved"  // Approved, not yet disbursed
	LoanStatusActive    = "active"    // Disbursed, repayments ongoing
	LoanStatusCompleted = "completed" // Fully paid off
	LoanStatusDefaulted = "defaulted" // Written off
	LoanStatusRejected  = "rejected"
	LoanStatusCancelled = "cancelled"
)

// Loan represents a credit extended to a Borrower, optionally tied to a
// Motorcycle (asset financing). RepaymentSchedule rows are generated when the
// loan is created; Repayment rows record actual collections.
type Loan struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	BusinessID    uint   `gorm:"not null;index" json:"business_id"`
	BranchID      uint   `gorm:"not null;index" json:"branch_id"`
	BorrowerID    uint   `gorm:"not null;index" json:"borrower_id"`
	MotorcycleID  *uint  `gorm:"index" json:"motorcycle_id"` // nil = working-capital loan, no asset
	LoanProductID uint   `gorm:"not null;index" json:"loan_product_id"`
	LoanNumber    string `gorm:"uniqueIndex" json:"loan_number"` // human-readable, e.g. KM-LN-000123

	// Principal & terms snapshot (copied from LoanProduct at creation)
	PrincipalAmount float64 `gorm:"not null" json:"principal_amount"`
	InitialDeposit  float64 `gorm:"default:0" json:"initial_deposit"` // paid up-front by borrower, reduces disbursed
	DisbursedAmount float64 `gorm:"not null" json:"disbursed_amount"` // = principal - initial_deposit
	Duration        int     `gorm:"not null" json:"duration"`         // # of repayment cycles
	InterestMethod  string  `gorm:"not null" json:"interest_method"`
	InterestRate    float64 `gorm:"not null" json:"interest_rate"` // annual %
	RepaymentCycle  string  `gorm:"not null" json:"repayment_cycle"`

	// Computed amounts
	TotalInterest    float64 `json:"total_interest"`
	TotalRepayments  int     `json:"total_repayments"` // = Duration; kept for clarity
	InstallmentAmount float64 `json:"installment_amount"`
	TotalAmount      float64 `json:"total_amount"`      // disbursed + total_interest
	BalanceRemaining float64 `json:"balance_remaining"` // updated as Repayments come in

	// Status & lifecycle
	Status            string     `gorm:"default:pending;index" json:"status"`
	DisbursementDate  *time.Time `json:"disbursement_date"`
	FirstPaymentDate  *time.Time `json:"first_payment_date"`
	NextPaymentDate   *time.Time `json:"next_payment_date"`
	MaturityDate      *time.Time `json:"maturity_date"`
	CompletedAt       *time.Time `json:"completed_at"`

	// Audit
	CreatedBy   uint      `json:"created_by"`
	ApprovedBy  *uint     `json:"approved_by"`
	DisbursedBy *uint     `json:"disbursed_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Borrower    Borrower    `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
	Motorcycle  *Motorcycle `gorm:"foreignKey:MotorcycleID" json:"motorcycle,omitempty"`
	LoanProduct LoanProduct `gorm:"foreignKey:LoanProductID" json:"loan_product,omitempty"`
	Branch      Branch      `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
