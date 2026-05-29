package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/mail"
	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

type SaleHandler struct {
	DB     *gorm.DB
	Mailer *mail.Mailer
	Audit  *services.AuditService
}

type CreateSaleRequest struct {
	BranchID       uint   `json:"branch_id" binding:"required"`
	PaymentMethod  string `json:"payment_method" binding:"required"`
	CustomerPhone  string `json:"customer_phone"`
	DiscountAmount float64 `json:"discount_amount"`
	Items          []struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

// Create processes a sale with stock deduction in a transaction.
func (h *SaleHandler) Create(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")

	if role != models.RoleAdmin && role != models.RoleManager && role != models.RoleCashier {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins, managers, and cashiers can process sales"})
		return
	}

	var req CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if req.PaymentMethod != "cash" && req.PaymentMethod != "mobile_money" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method must be 'cash' or 'mobile_money'"})
		return
	}

	// Validate branch belongs to business
	var branch models.Branch
	if err := h.DB.Where("id = ? AND business_id = ?", req.BranchID, businessID).First(&branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch not found"})
		return
	}

	var sale models.Sale

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validate all products and check stock
		var subtotal float64
		saleItems := make([]models.SaleItem, 0, len(req.Items))

		for _, item := range req.Items {
			var product models.Product
			if err := tx.Where("id = ? AND business_id = ?", item.ProductID, businessID).First(&product).Error; err != nil {
				return fmt.Errorf("product %d not found", item.ProductID)
			}

			// Check stock
			ok, err := services.CheckSufficientStock(tx, item.ProductID, req.BranchID, item.Quantity)
			if err != nil {
				return fmt.Errorf("checking stock for product %d: %w", item.ProductID, err)
			}
			if !ok {
				return fmt.Errorf("insufficient stock for '%s': not enough units in %s", product.Title, branch.Name)
			}

			lineTotal := product.SellingPrice * float64(item.Quantity)
			subtotal += lineTotal

			saleItems = append(saleItems, models.SaleItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: product.SellingPrice,
				UnitCost:  product.CostPrice,
				LineTotal: lineTotal,
			})
		}

		// 2. Create Sale
		total := subtotal - req.DiscountAmount
		if total < 0 {
			total = 0
		}

		sale = models.Sale{
			BranchID:       req.BranchID,
			CashierID:      userID,
			PaymentMethod:  req.PaymentMethod,
			CustomerPhone:  req.CustomerPhone,
			Subtotal:       subtotal,
			DiscountAmount: req.DiscountAmount,
			Total:          total,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return fmt.Errorf("creating sale: %w", err)
		}

		// 3. Create SaleItems and deduct stock
		for i := range saleItems {
			saleItems[i].SaleID = sale.ID
			if err := tx.Create(&saleItems[i]).Error; err != nil {
				return fmt.Errorf("creating sale item: %w", err)
			}

			// Deduct stock
			saleID := sale.ID
			if err := services.DeductStock(tx, saleItems[i].ProductID, req.BranchID,
				saleItems[i].Quantity, models.MovementSale, &saleID, userID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Check low stock alerts for each sold item (async, non-blocking)
	for _, item := range req.Items {
		services.CheckAndAlertLowStock(h.DB, h.Mailer, item.ProductID, req.BranchID)
	}

	// Reload sale with items and relations
	h.DB.Preload("Items").Preload("Items.Product").
		Preload("Branch").Preload("Cashier").
		First(&sale, sale.ID)

	h.Audit.Log(c, models.ActionCreate, "sale", &sale.ID,
		fmt.Sprintf("Sold %d item(s) for UGX %s (%s)", len(sale.Items), formatMoney(sale.Total), sale.PaymentMethod),
		gin.H{"items": len(sale.Items), "total": sale.Total, "method": sale.PaymentMethod, "branch_id": sale.BranchID},
	)

	c.JSON(http.StatusCreated, gin.H{"data": h.saleResponse(sale)})
}

// List returns paginated sales history.
func (h *SaleHandler) List(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	query := h.DB.
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ?", businessID)

	if bid := c.Query("branch_id"); bid != "" {
		query = query.Where("sales.branch_id = ?", bid)
	}
	if cid := c.Query("cashier_id"); cid != "" {
		query = query.Where("sales.cashier_id = ?", cid)
	}
	if pm := c.Query("payment_method"); pm != "" {
		query = query.Where("sales.payment_method = ?", pm)
	}
	if from := c.Query("from_date"); from != "" {
		query = query.Where("sales.created_at >= ?", from)
	}
	if to := c.Query("to_date"); to != "" {
		query = query.Where("sales.created_at <= ?", to+" 23:59:59")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 50
	if page < 1 {
		page = 1
	}

	var total int64
	query.Model(&models.Sale{}).Count(&total)

	var sales []models.Sale
	query.Preload("Items").Preload("Items.Product").
		Preload("Branch").Preload("Cashier").
		Order("sales.created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&sales)

	results := make([]gin.H, 0, len(sales))
	for _, s := range sales {
		results = append(results, h.saleResponse(s))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     results,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// CreateReturn records a customer return against a sale: validates the line
// items, restocks via services.AddStock with MovementReturn, and persists
// Return + ReturnItem rows in one transaction.
func (h *SaleHandler) CreateReturn(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager && role != models.RoleCashier {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins, managers, and cashiers can process returns"})
		return
	}

	saleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	var req struct {
		Items []struct {
			SaleItemID uint `json:"sale_item_id" binding:"required"`
			Quantity   int  `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
		Reason         string `json:"reason"`
		PaymentMethod  string `json:"payment_method"`
		TransactionRef string `json:"transaction_ref"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "items[] is required"})
		return
	}

	lines := make([]services.ReturnLineInput, 0, len(req.Items))
	for _, it := range req.Items {
		lines = append(lines, services.ReturnLineInput{
			SaleItemID: it.SaleItemID,
			Quantity:   it.Quantity,
		})
	}

	ret, err := services.ProcessReturn(h.DB, services.ProcessReturnInput{
		BusinessID:     businessID,
		SaleID:         uint(saleID),
		Items:          lines,
		Reason:         req.Reason,
		PaymentMethod:  req.PaymentMethod,
		TransactionRef: req.TransactionRef,
		ProcessedBy:    userID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sid := uint(saleID)
	h.Audit.Log(c, models.ActionReturn, "sale", &sid,
		fmt.Sprintf("Refunded UGX %s on Sale #%d", formatMoney(ret.RefundedTotal), sid),
		gin.H{"items": len(ret.Items), "refunded_total": ret.RefundedTotal, "reason": req.Reason},
	)

	c.JSON(http.StatusCreated, gin.H{"data": ret})
}

// ListReturns returns all return events recorded against a sale (for the
// "Returns" history shown on the sale detail page).
func (h *SaleHandler) ListReturns(c *gin.Context) {
	businessID := c.GetUint("business_id")
	saleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// Verify sale ownership before returning anything (cheap defence-in-depth).
	var ok int64
	h.DB.Model(&models.Sale{}).
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("sales.id = ? AND branches.business_id = ?", saleID, businessID).
		Count(&ok)
	if ok == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	var returns []models.Return
	h.DB.Preload("Items.Product").Preload("Processor").
		Where("sale_id = ?", saleID).
		Order("created_at DESC").
		Find(&returns)
	c.JSON(http.StatusOK, gin.H{"data": returns})
}

// Get returns a single sale with items (receipt view).
func (h *SaleHandler) Get(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	var sale models.Sale
	if err := h.DB.
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("sales.id = ? AND branches.business_id = ?", id, businessID).
		Preload("Items").Preload("Items.Product").
		Preload("Branch").Preload("Cashier").
		First(&sale).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.saleResponse(sale)})
}

func (h *SaleHandler) saleResponse(s models.Sale) gin.H {
	items := make([]gin.H, 0, len(s.Items))
	for _, item := range s.Items {
		items = append(items, gin.H{
			"id":           item.ID,
			"product_id":   item.ProductID,
			"product_name": item.Product.Title,
			"quantity":     item.Quantity,
			"unit_price":   item.UnitPrice,
			"unit_cost":    item.UnitCost,
			"line_total":   item.LineTotal,
		})
	}

	resp := gin.H{
		"id":              s.ID,
		"branch_id":       s.BranchID,
		"branch_name":     s.Branch.Name,
		"cashier_id":      s.CashierID,
		"cashier_name":    s.Cashier.Name,
		"payment_method":  s.PaymentMethod,
		"customer_phone":  s.CustomerPhone,
		"subtotal":        s.Subtotal,
		"discount_amount": s.DiscountAmount,
		"total":           s.Total,
		"items":           items,
		"created_at":      s.CreatedAt,
	}
	return resp
}
