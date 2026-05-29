package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/database"
	"gritdemo/internal/models"
)

// AdminHandler hosts dangerous one-shot maintenance endpoints. Currently:
//   POST /api/admin/wipe-seed   — admin-only nuclear reset
//
// This file exists because the temporary "Reset Database" button in
// Settings needs a server endpoint to call. Remove the route + this file
// once the demo data has been cleared and you're past first-deploy.
type AdminHandler struct {
	DB *gorm.DB
}

// WipeSeedData truncates every domain table and re-runs the seeder. Used
// to clear leftover Nakawa Fashion demo data from older deploys without
// having to wipe the postgres docker volume manually.
//
// All sessions become invalid because every User row is gone — the caller
// MUST log out client-side after a successful response.
func (h *AdminHandler) WipeSeedData(c *gin.Context) {
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can reset the database"})
		return
	}

	// Order matters: child rows before parents, otherwise FK constraints fire.
	// Soft-delete tables get .Unscoped() so we delete rows for real.
	tables := []interface{}{
		// Loans (most dependent first)
		&models.Repayment{},
		&models.RepaymentSchedule{},
		&models.LoanFee{},
		&models.LoanPenalty{},
		&models.LoanCollateral{},
		&models.Loan{},
		&models.Borrower{},
		&models.LoanProduct{},

		// Motorcycles
		&models.CashSale{},
		&models.Motorcycle{},

		// Daily Boda
		&models.DailyBodaPayment{},
		&models.DailyBodaMotorcycle{},
		&models.DailyBodaDriver{},

		// Spares POS
		&models.SaleItem{},
		&models.Sale{},
		&models.StockTransfer{},
		&models.StockMovement{},
		&models.Stock{},
		&models.Product{},
		&models.Category{},

		// Identity / tenancy (last — everything else FKs into these)
		&models.Invitation{},
		&models.UserBusinessRole{},
		&models.Branch{},
		&models.Business{},
		&models.User{},
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		for _, t := range tables {
			if err := tx.Unscoped().Where("1 = 1").Delete(t).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("admin: wipe-seed failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Re-run the seeder so the user has something to log into.
	if err := database.Seed(h.DB); err != nil {
		log.Printf("admin: re-seed failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wipe succeeded but re-seed failed: " + err.Error()})
		return
	}

	log.Println("admin: database wiped and re-seeded")
	c.JSON(http.StatusOK, gin.H{
		"message":   "Database reset. You'll be signed out — log back in with the fresh admin credentials.",
		"new_admin": "admin@grit.demo",
		"password":  "password123",
	})
}
