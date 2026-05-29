package services

import (
	"log"

	"gorm.io/gorm"

	"gritdemo/internal/mail"
	"gritdemo/internal/models"
)

// CheckAndAlertLowStock checks if a product's stock in a branch has fallen
// below its threshold, and sends an email alert to the business admin(s).
// This should be called AFTER a stock-deducting operation (sale, transfer_out).
// It runs asynchronously (goroutine) so it doesn't block the request.
func CheckAndAlertLowStock(db *gorm.DB, mailer *mail.Mailer, productID, branchID uint) {
	if mailer == nil {
		return
	}

	go func() {
		// Get current stock
		var stock models.Stock
		if err := db.Where("product_id = ? AND branch_id = ?", productID, branchID).First(&stock).Error; err != nil {
			return
		}

		// Get product for threshold and title
		var product models.Product
		if err := db.First(&product, productID).Error; err != nil {
			return
		}

		// Only alert if at or below threshold
		if stock.Quantity > product.LowStockThreshold {
			return
		}

		// Get branch name
		var branch models.Branch
		if err := db.First(&branch, branchID).Error; err != nil {
			return
		}

		// Find admin(s) of this business
		var adminRoles []models.UserBusinessRole
		db.Where("business_id = ? AND role = ?", product.BusinessID, models.RoleAdmin).Find(&adminRoles)

		for _, role := range adminRoles {
			var user models.User
			if err := db.First(&user, role.UserID).Error; err != nil {
				continue
			}

			if err := mailer.SendLowStockAlert(
				user.Email,
				product.Title,
				branch.Name,
				stock.Quantity,
				product.LowStockThreshold,
			); err != nil {
				log.Printf("Failed to send low stock alert to %s: %v", user.Email, err)
			} else {
				log.Printf("Low stock alert sent to %s for %s (%d units at %s)",
					user.Email, product.Title, stock.Quantity, branch.Name)
			}
		}
	}()
}
