package models

import (
	"time"

	"gorm.io/gorm"
)

// Risk level constants.
const (
	RiskLevelLow      = "low"
	RiskLevelMedium   = "medium"
	RiskLevelHigh     = "high"
	RiskLevelCritical = "critical"
)

// Employment status constants.
const (
	EmploymentEmployed     = "employed"
	EmploymentSelfEmployed = "self_employed"
	EmploymentUnemployed   = "unemployed"
	EmploymentStudent      = "student"
	EmploymentRetired      = "retired"
)

// Borrower is a person who has taken (or applied for) a loan. Phone is unique
// per business — same person joining a different Grit tenant is a new record.
type Borrower struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	BusinessID       uint           `gorm:"not null;index;uniqueIndex:idx_business_phone" json:"business_id"`
	BranchID         uint           `gorm:"not null;index" json:"branch_id"`
	FirstName        string         `gorm:"not null" json:"first_name"`
	LastName         string         `gorm:"not null" json:"last_name"`
	Phone            string         `gorm:"not null;uniqueIndex:idx_business_phone" json:"phone"`
	AltPhone         string         `json:"alt_phone"`
	Email            string         `json:"email"`
	DateOfBirth      *time.Time     `json:"date_of_birth"`
	Gender           string         `json:"gender"`
	NationalID       string         `gorm:"index" json:"national_id"`     // NIN
	Address          string         `gorm:"type:text" json:"address"`
	EmploymentStatus string         `json:"employment_status"`
	Occupation       string         `json:"occupation"`
	Employer         string         `json:"employer"`
	MonthlyIncome    float64        `json:"monthly_income"`
	NextOfKinName    string         `json:"next_of_kin_name"`
	NextOfKinPhone   string         `json:"next_of_kin_phone"`
	NextOfKinRelation string        `json:"next_of_kin_relation"`
	CreditScore      int            `gorm:"default:0" json:"credit_score"` // 0-1000 internal
	RiskLevel        string         `gorm:"default:medium" json:"risk_level"`
	PhotoKey         string         `json:"-"` // R2 key
	IDDocumentKey    string         `json:"-"` // R2 key — NIN scan
	CreatedBy        uint           `json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Branch           Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// FullName is a convenience accessor for "First Last" rendering.
func (b *Borrower) FullName() string {
	return b.FirstName + " " + b.LastName
}
