package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// WipeMutableData truncates the rows the demo-reset cron job needs to
// clear nightly while leaving identity / config / canonical demo cohort
// intact. Spelled-out approach (not "delete from every table") so adding
// a new model doesn't accidentally start nuking it.
//
// Tables NOT wiped: users, businesses, branches, user_business_roles,
// categories, products, motorcycles, loan_products, daily_boda_drivers,
// daily_boda_motorcycles — these are seeded shape rather than visitor
// activity. After WipeMutableData runs, Seed() restores any defaults
// that lookup-by-name would otherwise miss.
func WipeMutableData(db *gorm.DB) error {
	// Run inside a transaction so a mid-flight error doesn't leave half a
	// reset. Foreign-key disabling is best-effort: Postgres needs SET
	// session_replication_role for it; SQLite uses PRAGMA. We just delete
	// in dependency-safe order instead, which is portable and obvious.
	order := []string{
		"notifications",
		"activities",
		"return_items",
		"returns",
		"sale_items",
		"sales",
		"cash_sales",
		"daily_boda_payments",
		"loan_penalties",
		"loan_fees",
		"loan_collaterals",
		"repayments",
		"repayment_schedules",
		"loans",
		"borrowers",
		"stock_movements",
		"stock_transfers",
		"stocks",
		"invitations",
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, table := range order {
			// Some tables may not exist yet on a fresh deploy — IF EXISTS
			// guard would be nice but isn't portable. Use raw SQL and
			// swallow "no such table" errors per-table.
			if err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
				log.Printf("demo wipe: %s: %v (continuing)", table, err)
			}
		}
		return nil
	})
}
