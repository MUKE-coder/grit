package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// ReturnLineInput is one line item the customer is bringing back.
type ReturnLineInput struct {
	SaleItemID uint
	Quantity   int
}

// ProcessReturnInput is the data needed to record a customer return.
type ProcessReturnInput struct {
	BusinessID     uint
	SaleID         uint
	Items          []ReturnLineInput
	Reason         string
	PaymentMethod  string
	TransactionRef string
	ProcessedBy    uint
}

// ProcessReturn validates the requested return against the original sale,
// restocks the products, and persists Return + ReturnItem rows in a single
// transaction. Idempotent guards:
//   - sale must belong to the active business
//   - each sale_item_id must belong to that sale
//   - quantity returned can't exceed (sold qty − already-returned qty) per line
//   - quantity must be positive
//
// Returns the persisted Return with items preloaded.
func ProcessReturn(db *gorm.DB, in ProcessReturnInput) (*models.Return, error) {
	if len(in.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Load sale and verify ownership.
	var sale models.Sale
	err := db.Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("sales.id = ? AND branches.business_id = ?", in.SaleID, in.BusinessID).
		First(&sale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sale not found")
		}
		return nil, err
	}

	// Pull all sale items up-front so we can validate without N round-trips.
	var saleItems []models.SaleItem
	if err := db.Where("sale_id = ?", sale.ID).Find(&saleItems).Error; err != nil {
		return nil, err
	}
	saleItemByID := make(map[uint]models.SaleItem, len(saleItems))
	for _, si := range saleItems {
		saleItemByID[si.ID] = si
	}

	// How much of each sale_item has already been returned in prior returns?
	type alreadyRow struct {
		SaleItemID uint
		Total      int
	}
	var already []alreadyRow
	db.Model(&models.ReturnItem{}).
		Select("sale_item_id, COALESCE(SUM(quantity), 0) as total").
		Joins("JOIN returns ON returns.id = return_items.return_id").
		Where("returns.sale_id = ?", sale.ID).
		Group("sale_item_id").
		Scan(&already)
	prior := make(map[uint]int, len(already))
	for _, a := range already {
		prior[a.SaleItemID] = a.Total
	}

	// Build the planned ReturnItem rows + total refund amount.
	plannedItems := make([]models.ReturnItem, 0, len(in.Items))
	var refundTotal float64
	for _, line := range in.Items {
		if line.Quantity <= 0 {
			return nil, fmt.Errorf("return quantity must be positive (sale_item_id=%d)", line.SaleItemID)
		}
		si, ok := saleItemByID[line.SaleItemID]
		if !ok {
			return nil, fmt.Errorf("sale_item %d does not belong to this sale", line.SaleItemID)
		}
		remaining := si.Quantity - prior[line.SaleItemID]
		if line.Quantity > remaining {
			return nil, fmt.Errorf("can't return %d of '%s' — only %d remaining returnable on this line",
				line.Quantity, productNameOrID(si.ProductID, db), remaining)
		}
		lineTotal := round2(si.UnitPrice * float64(line.Quantity))
		refundTotal += lineTotal
		plannedItems = append(plannedItems, models.ReturnItem{
			SaleItemID: si.ID,
			ProductID:  si.ProductID,
			Quantity:   line.Quantity,
			UnitPrice:  si.UnitPrice,
			LineTotal:  lineTotal,
		})
	}

	// Persist the return + restock in one transaction.
	ret := models.Return{
		BusinessID:     in.BusinessID,
		SaleID:         sale.ID,
		BranchID:       sale.BranchID,
		RefundedTotal:  round2(refundTotal),
		PaymentMethod:  in.PaymentMethod,
		TransactionRef: in.TransactionRef,
		Reason:         in.Reason,
		ProcessedBy:    in.ProcessedBy,
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ret).Error; err != nil {
			return err
		}
		// Stamp ReturnID on each item, persist, and restock.
		for i := range plannedItems {
			plannedItems[i].ReturnID = ret.ID
		}
		if err := tx.CreateInBatches(plannedItems, 50).Error; err != nil {
			return err
		}
		for _, line := range plannedItems {
			refID := ret.ID
			if err := AddStock(tx, line.ProductID, sale.BranchID, line.Quantity,
				models.MovementReturn, "Customer return", &refID, in.ProcessedBy); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return the fully-loaded record.
	var out models.Return
	if err := db.Preload("Items.Product").Preload("Branch").Preload("Processor").
		First(&out, ret.ID).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func productNameOrID(productID uint, db *gorm.DB) string {
	var p models.Product
	if err := db.Select("title").First(&p, productID).Error; err == nil && p.Title != "" {
		return p.Title
	}
	return fmt.Sprintf("product #%d", productID)
}
