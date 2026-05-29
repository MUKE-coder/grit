package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Models returns all models for migration in dependency order.
func Models() []interface{} {
	return []interface{}{
		// Identity & tenancy
		&User{},
		&Business{},
		&Branch{},
		&UserBusinessRole{},
		&Invitation{},

		// Spares (existing POS)
		&Category{},
		&Product{},
		&Stock{},
		&StockMovement{},
		&StockTransfer{},
		&Sale{},
		&SaleItem{},
		&Return{},
		&ReturnItem{},

		// Motorcycles inventory + cash sales
		&Motorcycle{},
		&CashSale{},

		// Loans
		&LoanProduct{},
		&Borrower{},
		&Loan{},
		&RepaymentSchedule{},
		&Repayment{},
		&LoanFee{},
		&LoanPenalty{},
		&LoanCollateral{},

		// Daily Boda (rental fleet)
		&DailyBodaDriver{},
		&DailyBodaMotorcycle{},
		&DailyBodaPayment{},

		// Audit log
		&Activity{},

		// In-app notifications surfaced from Sentinel + Pulse pollers
		&Notification{},
	}
}

// Migrate runs AutoMigrate for all models.
func Migrate(db *gorm.DB) error {
	models := Models()
	log.Printf("Migrating %d models...", len(models))

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("migrating %T: %w", model, err)
		}
		log.Printf("  ✓ %T", model)
	}

	log.Println("All migrations complete.")
	return nil
}
