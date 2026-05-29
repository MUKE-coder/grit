package handlers

import (
	"strings"

	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// validateBranchInBusiness returns true if the branch belongs to the given business.
// Used by every create/update that takes a branch_id from the request body.
func validateBranchInBusiness(db *gorm.DB, branchID, businessID uint) bool {
	var count int64
	db.Model(&models.Branch{}).Where("id = ? AND business_id = ?", branchID, businessID).Count(&count)
	return count > 0
}

// isUniqueViolation matches the postgres + sqlite error strings for a unique
// constraint violation. We don't have access to driver-specific error codes
// here without importing pgx, so string match is the cheap path.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}

// hasAnyRole returns true if `actual` matches any of `allowed`.
func hasAnyRole(actual interface{}, allowed ...string) bool {
	s, ok := actual.(string)
	if !ok {
		return false
	}
	for _, r := range allowed {
		if s == r {
			return true
		}
	}
	return false
}
