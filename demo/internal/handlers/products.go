package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
	"gritdemo/internal/services"
	"gritdemo/internal/storage"
)

type ProductHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	Audit   *services.AuditService
}

// productResponse builds a product response with image URL instead of image key.
func (h *ProductHandler) productResponse(p models.Product) gin.H {
	resp := gin.H{
		"id":                  p.ID,
		"business_id":         p.BusinessID,
		"category_id":         p.CategoryID,
		"title":               p.Title,
		"description":         p.Description,
		"barcode":             p.Barcode,
		"selling_price":       p.SellingPrice,
		"cost_price":          p.CostPrice,
		"low_stock_threshold": p.LowStockThreshold,
		"created_by":          p.CreatedBy,
		"created_at":          p.CreatedAt,
		"updated_at":          p.UpdatedAt,
	}

	if p.Category.ID != 0 {
		resp["category"] = p.Category
	}

	// Image URL
	if p.ImageKey != "" && h.Storage != nil {
		url, err := h.Storage.GetSignedURL(context.Background(), p.ImageKey, 1*time.Hour)
		if err == nil {
			resp["image_url"] = url
		}
	}

	return resp
}

// List returns products for the active business with stock per branch.
func (h *ProductHandler) List(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	query := h.DB.Preload("Category").Where("business_id = ?", businessID)

	// Filters
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR barcode = ?", like, search)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	var total int64
	query.Model(&models.Product{}).Count(&total)

	var products []models.Product
	query.Order("title ASC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&products)

	// Get stock per branch for each product
	type StockInfo struct {
		BranchID   uint   `json:"branch_id"`
		BranchName string `json:"branch_name"`
		Quantity   int    `json:"quantity"`
	}

	// Was N+1: one stock query per product. Now: one query for the whole page.
	productIDs := make([]uint, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}
	stocksByProduct := make(map[uint][]StockInfo)
	if len(productIDs) > 0 {
		var stocks []models.Stock
		h.DB.Preload("Branch").Where("product_id IN ?", productIDs).Find(&stocks)
		for _, s := range stocks {
			stocksByProduct[s.ProductID] = append(stocksByProduct[s.ProductID], StockInfo{
				BranchID:   s.BranchID,
				BranchName: s.Branch.Name,
				Quantity:   s.Quantity,
			})
		}
	}

	results := make([]gin.H, 0, len(products))
	for _, p := range products {
		resp := h.productResponse(p)
		stocks := stocksByProduct[p.ID]
		if stocks == nil {
			stocks = make([]StockInfo, 0)
		}
		resp["stock"] = stocks
		results = append(results, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     results,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// Get returns a single product with full details and stock per branch.
func (h *ProductHandler) Get(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.DB.Preload("Category").Where("id = ? AND business_id = ?", id, businessID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	resp := h.productResponse(product)

	// Stock per branch
	var stocks []models.Stock
	h.DB.Preload("Branch").Where("product_id = ?", product.ID).Find(&stocks)

	type StockInfo struct {
		BranchID   uint   `json:"branch_id"`
		BranchName string `json:"branch_name"`
		Quantity   int    `json:"quantity"`
	}
	stockInfos := make([]StockInfo, 0, len(stocks))
	for _, s := range stocks {
		stockInfos = append(stockInfos, StockInfo{
			BranchID:   s.BranchID,
			BranchName: s.Branch.Name,
			Quantity:   s.Quantity,
		})
	}
	resp["stock"] = stockInfos

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// POS returns a lightweight product list optimized for the POS screen.
func (h *ProductHandler) POS(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	branchIDStr := c.Query("branch_id")
	if branchIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "branch_id is required"})
		return
	}
	branchID, _ := strconv.ParseUint(branchIDStr, 10, 64)

	query := h.DB.Preload("Category").Where("products.business_id = ?", businessID)

	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR barcode = ?", like, search)
	}
	if barcode := c.Query("barcode"); barcode != "" {
		query = query.Where("barcode = ?", barcode)
	}

	var products []models.Product
	query.Order("title ASC").Find(&products)

	type POSProduct struct {
		ID            uint    `json:"id"`
		Title         string  `json:"title"`
		SellingPrice  float64 `json:"selling_price"`
		CostPrice     float64 `json:"cost_price"`
		CategoryID    uint    `json:"category_id"`
		CategoryName  string  `json:"category_name"`
		Barcode       string  `json:"barcode"`
		ImageURL      string  `json:"image_url,omitempty"`
		StockInBranch int     `json:"stock_in_branch"`
	}

	// Was N+1: one stock query per product. Single query with IN clause now.
	productIDs := make([]uint, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}
	stockByProduct := make(map[uint]int, len(productIDs))
	if len(productIDs) > 0 {
		var stocks []models.Stock
		h.DB.Where("product_id IN ? AND branch_id = ?", productIDs, branchID).Find(&stocks)
		for _, s := range stocks {
			stockByProduct[s.ProductID] = s.Quantity
		}
	}

	results := make([]POSProduct, 0, len(products))
	for _, p := range products {
		pp := POSProduct{
			ID:           p.ID,
			Title:        p.Title,
			SellingPrice: p.SellingPrice,
			CostPrice:    p.CostPrice,
			CategoryID:   p.CategoryID,
			CategoryName: p.Category.Name,
			Barcode:      p.Barcode,
			StockInBranch: stockByProduct[p.ID], // 0 if no stock row exists
		}

		// Signed-URL generation is local crypto in the AWS SDK so doing it
		// per row is cheap, but we could batch later if it becomes a hotspot.
		if p.ImageKey != "" && h.Storage != nil {
			if url, err := h.Storage.GetSignedURL(context.Background(), p.ImageKey, 1*time.Hour); err == nil {
				pp.ImageURL = url
			}
		}

		results = append(results, pp)
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// Import bulk-creates products from a parsed Excel sheet.
//
// The frontend parses the .xlsx client-side (SheetJS) and posts the rows
// as JSON. Each row carries: code (e.g. '21420H33000H000'), title (e.g.
// 'PLATE GROUP CLUTCH DRIVE'), cost_price, selling_price, quantity. The
// caller picks a target branch_id; quantity is added to that branch
// through the existing AddStock service so movements are audit-tracked.
//
// Behaviour:
//   - The 'Others' category is auto-created on first use.
//   - Duplicate barcodes (within the same business) are skipped.
//   - Rows with empty title or non-positive selling price are skipped.
//   - Everything runs in one transaction; partial failure rolls back.
//
// Response includes counts so the UI can show "Imported N, skipped M".
func (h *ProductHandler) Import(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can import products"})
		return
	}

	var req struct {
		BranchID uint `json:"branch_id" binding:"required"`
		Items    []struct {
			Code         string  `json:"code"`
			Title        string  `json:"title"`
			CostPrice    float64 `json:"cost_price"`
			SellingPrice float64 `json:"selling_price"`
			Quantity     int     `json:"quantity"`
		} `json:"items" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "branch_id and items[] are required"})
		return
	}

	// Validate branch.
	var branch models.Branch
	if err := h.DB.Where("id = ? AND business_id = ?", req.BranchID, businessID).First(&branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch not found in this business"})
		return
	}

	type SkippedRow struct {
		Index  int    `json:"index"`
		Title  string `json:"title"`
		Reason string `json:"reason"`
	}
	var (
		created  []models.Product
		skipped  []SkippedRow
		othersID uint
	)

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Find or create the "Others" category for this business.
		var others models.Category
		err := tx.Where("business_id = ? AND name = ?", businessID, "Others").First(&others).Error
		if err == gorm.ErrRecordNotFound {
			others = models.Category{BusinessID: businessID, Name: "Others"}
			if err := tx.Create(&others).Error; err != nil {
				return fmt.Errorf("creating Others category: %w", err)
			}
		} else if err != nil {
			return err
		}
		othersID = others.ID

		// Pre-load existing barcodes so we can skip duplicates without one
		// query per row.
		var existing []string
		tx.Model(&models.Product{}).
			Where("business_id = ? AND barcode <> ''", businessID).
			Pluck("barcode", &existing)
		seen := make(map[string]bool, len(existing))
		for _, b := range existing {
			seen[strings.ToLower(b)] = true
		}

		for i, row := range req.Items {
			title := strings.TrimSpace(row.Title)
			code := strings.TrimSpace(row.Code)
			if title == "" {
				skipped = append(skipped, SkippedRow{Index: i, Title: title, Reason: "empty title"})
				continue
			}
			if row.SellingPrice <= 0 {
				skipped = append(skipped, SkippedRow{Index: i, Title: title, Reason: "selling price must be > 0"})
				continue
			}
			if code != "" && seen[strings.ToLower(code)] {
				skipped = append(skipped, SkippedRow{Index: i, Title: title, Reason: "duplicate code in this business"})
				continue
			}

			p := models.Product{
				BusinessID:        businessID,
				CategoryID:        othersID,
				Title:             title,
				Barcode:           code,
				SellingPrice:      row.SellingPrice,
				CostPrice:         row.CostPrice,
				LowStockThreshold: 5,
				CreatedBy:         userID,
			}
			if err := tx.Create(&p).Error; err != nil {
				return fmt.Errorf("creating product %q: %w", title, err)
			}
			if code != "" {
				seen[strings.ToLower(code)] = true
			}
			if row.Quantity > 0 {
				if err := services.AddStock(tx, p.ID, req.BranchID, row.Quantity,
					models.MovementStockIn, "Excel import", nil, userID); err != nil {
					return fmt.Errorf("stocking %q: %w", title, err)
				}
			}
			created = append(created, p)
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Audit.Log(c, models.ActionImport, "product", nil,
		fmt.Sprintf("Imported %d product(s) from spreadsheet (%d skipped)", len(created), len(skipped)),
		gin.H{"created_count": len(created), "skipped_count": len(skipped), "branch_id": req.BranchID},
	)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"created_count":      len(created),
			"skipped_count":      len(skipped),
			"others_category_id": othersID,
			"branch_id":          req.BranchID,
			"skipped":            skipped,
		},
	})
}

// Create creates a new product (admin/manager).
func (h *ProductHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can create products"})
		return
	}

	var req struct {
		Title             string  `json:"title" binding:"required"`
		SellingPrice      float64 `json:"selling_price" binding:"required"`
		CategoryID        uint    `json:"category_id" binding:"required"`
		Description       string  `json:"description"`
		Barcode           string  `json:"barcode"`
		CostPrice         float64 `json:"cost_price"`
		LowStockThreshold int     `json:"low_stock_threshold"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, selling_price, and category_id are required"})
		return
	}

	// Validate category belongs to business
	var cat models.Category
	if err := h.DB.Where("id = ? AND business_id = ?", req.CategoryID, businessID).First(&cat).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	product := models.Product{
		BusinessID:        businessID,
		CategoryID:        req.CategoryID,
		Title:             strings.TrimSpace(req.Title),
		Description:       strings.TrimSpace(req.Description),
		Barcode:           strings.TrimSpace(req.Barcode),
		SellingPrice:      req.SellingPrice,
		CostPrice:         req.CostPrice,
		LowStockThreshold: req.LowStockThreshold,
		CreatedBy:         userID,
	}

	if product.LowStockThreshold == 0 {
		product.LowStockThreshold = 5
	}

	if err := h.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Reload with category
	h.DB.Preload("Category").First(&product, product.ID)

	h.Audit.Log(c, models.ActionCreate, "product", &product.ID,
		fmt.Sprintf("Created product %q", product.Title), nil)

	c.JSON(http.StatusCreated, gin.H{"data": h.productResponse(product)})
}

// Update updates a product (admin/manager).
func (h *ProductHandler) Update(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can update products"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req struct {
		Title             *string  `json:"title"`
		SellingPrice      *float64 `json:"selling_price"`
		CategoryID        *uint    `json:"category_id"`
		Description       *string  `json:"description"`
		Barcode           *string  `json:"barcode"`
		CostPrice         *float64 `json:"cost_price"`
		LowStockThreshold *int     `json:"low_stock_threshold"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Title != nil {
		product.Title = strings.TrimSpace(*req.Title)
	}
	if req.SellingPrice != nil {
		product.SellingPrice = *req.SellingPrice
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Description != nil {
		product.Description = strings.TrimSpace(*req.Description)
	}
	if req.Barcode != nil {
		product.Barcode = strings.TrimSpace(*req.Barcode)
	}
	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}
	if req.LowStockThreshold != nil {
		product.LowStockThreshold = *req.LowStockThreshold
	}

	h.DB.Save(&product)
	h.DB.Preload("Category").First(&product, product.ID)

	h.Audit.Log(c, models.ActionUpdate, "product", &product.ID,
		fmt.Sprintf("Updated product %q", product.Title), nil)

	c.JSON(http.StatusOK, gin.H{"data": h.productResponse(product)})
}

// Delete soft-deletes a product (admin only).
func (h *ProductHandler) Delete(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete products"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result := h.DB.Where("id = ? AND business_id = ?", id, businessID).Delete(&models.Product{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	pid := uint(id)
	h.Audit.Log(c, models.ActionDelete, "product", &pid, fmt.Sprintf("Deleted product #%d", pid), nil)

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// UploadImage handles product image upload to R2.
func (h *ProductHandler) UploadImage(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can upload images"})
		return
	}

	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage not configured"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// Delete old image if exists
	if product.ImageKey != "" {
		_ = h.Storage.Delete(context.Background(), product.ImageKey)
	}

	// Upload new image
	key := fmt.Sprintf("products/%d/%d-%s", businessID, product.ID, header.Filename)
	if err := h.Storage.Upload(context.Background(), key, file, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	product.ImageKey = key
	h.DB.Save(&product)

	// Return signed URL
	imageURL, _ := h.Storage.GetSignedURL(context.Background(), key, 1*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"image_url": imageURL,
		},
	})
}

// DeleteImage removes a product's image from R2.
func (h *ProductHandler) DeleteImage(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can delete images"})
		return
	}

	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage not configured"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.ImageKey != "" {
		_ = h.Storage.Delete(context.Background(), product.ImageKey)
		product.ImageKey = ""
		h.DB.Save(&product)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image removed"})
}
