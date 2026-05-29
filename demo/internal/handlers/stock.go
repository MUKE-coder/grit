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

type StockHandler struct {
	DB     *gorm.DB
	Mailer *mail.Mailer
	Audit  *services.AuditService
}

// Levels returns stock levels for the active business.
func (h *StockHandler) Levels(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	query := h.DB.
		Joins("JOIN products ON products.id = stocks.product_id").
		Joins("JOIN branches ON branches.id = stocks.branch_id").
		Where("products.business_id = ?", businessID)

	if branchID := c.Query("branch_id"); branchID != "" {
		query = query.Where("stocks.branch_id = ?", branchID)
	}
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("products.category_id = ?", catID)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("LOWER(products.title) LIKE LOWER(?)", like)
	}

	var stocks []models.Stock
	query.Preload("Product").Preload("Product.Category").Preload("Branch").
		Find(&stocks)

	type StockLevel struct {
		ProductID     uint    `json:"product_id"`
		ProductTitle  string  `json:"product_title"`
		CategoryName  string  `json:"category_name"`
		BranchID      uint    `json:"branch_id"`
		BranchName    string  `json:"branch_name"`
		Quantity      int     `json:"quantity"`
		Threshold     int     `json:"threshold"`
		Status        string  `json:"status"`
		SellingPrice  float64 `json:"selling_price"`
	}

	results := make([]StockLevel, 0, len(stocks))
	for _, s := range stocks {
		status := "ok"
		if s.Quantity == 0 {
			status = "out"
		} else if s.Quantity <= s.Product.LowStockThreshold {
			status = "low"
		}

		// Filter by status if requested
		if statusFilter := c.Query("status"); statusFilter != "" {
			if statusFilter != status {
				continue
			}
		}

		results = append(results, StockLevel{
			ProductID:    s.ProductID,
			ProductTitle: s.Product.Title,
			CategoryName: s.Product.Category.Name,
			BranchID:     s.BranchID,
			BranchName:   s.Branch.Name,
			Quantity:     s.Quantity,
			Threshold:    s.Product.LowStockThreshold,
			Status:       status,
			SellingPrice: s.Product.SellingPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// StockIn records stock arrival.
func (h *StockHandler) StockIn(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager && role != models.RoleStockClerk {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins, managers, and stock clerks can record stock-in"})
		return
	}

	var req struct {
		ProductID         uint    `json:"product_id" binding:"required"`
		BranchID          uint    `json:"branch_id" binding:"required"`
		Quantity          int     `json:"quantity" binding:"required,min=1"`
		CostPriceOverride *float64 `json:"cost_price_override"`
		Note              string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id, branch_id, and quantity (>0) are required"})
		return
	}

	// Validate product belongs to business
	var product models.Product
	if err := h.DB.Where("id = ? AND business_id = ?", req.ProductID, businessID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found in this business"})
		return
	}

	// Validate branch belongs to business
	var branch models.Branch
	if err := h.DB.Where("id = ? AND business_id = ?", req.BranchID, businessID).First(&branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch not found in this business"})
		return
	}

	// Update cost price if override provided
	if req.CostPriceOverride != nil {
		h.DB.Model(&product).Update("cost_price", *req.CostPriceOverride)
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		return services.AddStock(tx, req.ProductID, req.BranchID, req.Quantity,
			models.MovementStockIn, req.Note, nil, userID)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock-in: " + err.Error()})
		return
	}

	// Get updated stock level
	qty, _ := services.GetStockLevel(h.DB, req.ProductID, req.BranchID)

	h.Audit.Log(c, models.ActionUpdate, "stock", &product.ID,
		fmt.Sprintf("Stocked in %d × %s at %s (now %d)", req.Quantity, product.Title, branch.Name, qty),
		gin.H{"product_id": req.ProductID, "branch_id": req.BranchID, "quantity": req.Quantity},
	)

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"product_id":    req.ProductID,
			"branch_id":     req.BranchID,
			"quantity_added": req.Quantity,
			"current_stock": qty,
		},
	})
}

// Transfer moves stock between branches (admin only).
func (h *StockHandler) Transfer(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can transfer stock"})
		return
	}

	var req struct {
		ProductID    uint   `json:"product_id" binding:"required"`
		FromBranchID uint   `json:"from_branch_id" binding:"required"`
		ToBranchID   uint   `json:"to_branch_id" binding:"required"`
		Quantity     int    `json:"quantity" binding:"required,min=1"`
		Note         string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id, from_branch_id, to_branch_id, and quantity are required"})
		return
	}

	if req.FromBranchID == req.ToBranchID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to the same branch"})
		return
	}

	// Validate product + branches belong to business
	var product models.Product
	if err := h.DB.Where("id = ? AND business_id = ?", req.ProductID, businessID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	for _, bid := range []uint{req.FromBranchID, req.ToBranchID} {
		var b models.Branch
		if err := h.DB.Where("id = ? AND business_id = ?", bid, businessID).First(&b).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Branch not found"})
			return
		}
	}

	// Check sufficient stock
	ok, _ := services.CheckSufficientStock(h.DB, req.ProductID, req.FromBranchID, req.Quantity)
	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Insufficient stock in source branch"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Create transfer record
		transfer := models.StockTransfer{
			BusinessID:    businessID.(uint),
			FromBranchID:  req.FromBranchID,
			ToBranchID:    req.ToBranchID,
			ProductID:     req.ProductID,
			Quantity:      req.Quantity,
			Note:          req.Note,
			TransferredBy: userID,
		}
		if err := tx.Create(&transfer).Error; err != nil {
			return err
		}

		refID := transfer.ID

		// Deduct from source
		if err := services.DeductStock(tx, req.ProductID, req.FromBranchID, req.Quantity,
			models.MovementTransferOut, &refID, userID); err != nil {
			return err
		}

		// Add to destination
		return services.AddStock(tx, req.ProductID, req.ToBranchID, req.Quantity,
			models.MovementTransferIn, req.Note, &refID, userID)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed: " + err.Error()})
		return
	}

	// Check low stock on source branch (transfer_out reduces stock)
	services.CheckAndAlertLowStock(h.DB, h.Mailer, req.ProductID, req.FromBranchID)

	fromQty, _ := services.GetStockLevel(h.DB, req.ProductID, req.FromBranchID)
	toQty, _ := services.GetStockLevel(h.DB, req.ProductID, req.ToBranchID)

	h.Audit.Log(c, models.ActionTransfer, "stock", &product.ID,
		fmt.Sprintf("Transferred %d × %s between branches", req.Quantity, product.Title),
		gin.H{"product_id": req.ProductID, "quantity": req.Quantity, "from_branch_id": req.FromBranchID, "to_branch_id": req.ToBranchID},
	)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"product_id":       req.ProductID,
			"quantity":         req.Quantity,
			"from_branch_stock": fromQty,
			"to_branch_stock":   toQty,
		},
	})
}

// Movements returns paginated stock movement history.
func (h *StockHandler) Movements(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	query := h.DB.
		Joins("JOIN products ON products.id = stock_movements.product_id").
		Where("products.business_id = ?", businessID)

	if pid := c.Query("product_id"); pid != "" {
		query = query.Where("stock_movements.product_id = ?", pid)
	}
	if bid := c.Query("branch_id"); bid != "" {
		query = query.Where("stock_movements.branch_id = ?", bid)
	}
	if mt := c.Query("type"); mt != "" {
		query = query.Where("stock_movements.movement_type = ?", mt)
	}
	if from := c.Query("from_date"); from != "" {
		query = query.Where("stock_movements.created_at >= ?", from)
	}
	if to := c.Query("to_date"); to != "" {
		query = query.Where("stock_movements.created_at <= ?", to+" 23:59:59")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 50
	if page < 1 {
		page = 1
	}

	var total int64
	query.Model(&models.StockMovement{}).Count(&total)

	var movements []models.StockMovement
	query.Preload("Product").Preload("Branch").
		Order("stock_movements.created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&movements)

	c.JSON(http.StatusOK, gin.H{
		"data":     movements,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// Transfers returns stock transfer history.
func (h *StockHandler) Transfers(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	query := h.DB.Where("business_id = ?", businessID)

	if bid := c.Query("branch_id"); bid != "" {
		query = query.Where("from_branch_id = ? OR to_branch_id = ?", bid, bid)
	}
	if from := c.Query("from_date"); from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to := c.Query("to_date"); to != "" {
		query = query.Where("created_at <= ?", to+" 23:59:59")
	}

	var transfers []models.StockTransfer
	query.Preload("Product").Preload("FromBranch").Preload("ToBranch").
		Order("created_at DESC").
		Find(&transfers)

	c.JSON(http.StatusOK, gin.H{"data": transfers})
}
