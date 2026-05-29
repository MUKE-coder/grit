package services

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gritdemo/internal/models"
)

// GetStockLevel returns the current stock quantity for a product in a branch.
func GetStockLevel(db *gorm.DB, productID, branchID uint) (int, error) {
	var stock models.Stock
	err := db.Where("product_id = ? AND branch_id = ?", productID, branchID).First(&stock).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return stock.Quantity, nil
}

// CheckSufficientStock returns true if the branch has enough stock.
func CheckSufficientStock(db *gorm.DB, productID, branchID uint, needed int) (bool, error) {
	qty, err := GetStockLevel(db, productID, branchID)
	if err != nil {
		return false, err
	}
	return qty >= needed, nil
}

// AddStock increases stock for a product in a branch (upsert) and logs a movement.
func AddStock(tx *gorm.DB, productID, branchID uint, qty int, movementType, note string, referenceID *uint, createdBy uint) error {
	if qty <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Upsert stock record
	stock := models.Stock{
		ProductID: productID,
		BranchID:  branchID,
	}
	result := tx.Where("product_id = ? AND branch_id = ?", productID, branchID).First(&stock)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		stock.Quantity = qty
		if err := tx.Create(&stock).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(&stock).Update("quantity", gorm.Expr("quantity + ?", qty)).Error; err != nil {
			return err
		}
	}

	// Log movement
	movement := models.StockMovement{
		ProductID:    productID,
		BranchID:     branchID,
		MovementType: movementType,
		Quantity:     qty,
		ReferenceID:  referenceID,
		Note:         note,
		CreatedBy:    createdBy,
	}
	return tx.Create(&movement).Error
}

// DeductStock decreases stock for a product in a branch and logs a movement.
func DeductStock(tx *gorm.DB, productID, branchID uint, qty int, movementType string, referenceID *uint, createdBy uint) error {
	if qty <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	var stock models.Stock
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ? AND branch_id = ?", productID, branchID).
		First(&stock).Error; err != nil {
		return fmt.Errorf("stock record not found")
	}

	if stock.Quantity < qty {
		return fmt.Errorf("insufficient stock: have %d, need %d", stock.Quantity, qty)
	}

	if err := tx.Model(&stock).Update("quantity", gorm.Expr("quantity - ?", qty)).Error; err != nil {
		return err
	}

	movement := models.StockMovement{
		ProductID:    productID,
		BranchID:     branchID,
		MovementType: movementType,
		Quantity:     -qty,
		ReferenceID:  referenceID,
		CreatedBy:    createdBy,
	}
	return tx.Create(&movement).Error
}
