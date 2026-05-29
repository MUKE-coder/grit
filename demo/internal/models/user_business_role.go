package models

import "time"

// Role constants for Grit Motors.
// admin/manager/cashier/stock_clerk apply to spares + motorcycles operations.
// loan_officer manages borrowers, loans, and repayments.
// accountant has read-only access to finance reports across both segments.
const (
	RoleAdmin       = "admin"
	RoleManager     = "manager"
	RoleCashier     = "cashier"
	RoleStockClerk  = "stock_clerk"
	RoleLoanOfficer = "loan_officer"
	RoleAccountant  = "accountant"
)

// WorkspaceAccess controls which sidebar layout(s) a user can use.
// - "loans"  = loans + motorcycles + daily-boda only
// - "spares" = spares POS + inventory only
// - "both"   = full access to both workspaces (default)
const (
	WorkspaceLoans  = "loans"
	WorkspaceSpares = "spares"
	WorkspaceBoth   = "both"
)

type UserBusinessRole struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	BusinessID      uint      `gorm:"not null;index" json:"business_id"`
	BranchID        *uint     `gorm:"index" json:"branch_id"` // nil = access to all branches
	Role            string    `gorm:"not null" json:"role"`    // admin|manager|cashier|stock_clerk
	WorkspaceAccess string    `gorm:"type:varchar(16);default:'both'" json:"workspace_access"` // loans|spares|both
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Business        Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	Branch          *Branch   `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
