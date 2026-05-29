package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

type ReportHandler struct {
	DB *gorm.DB
}

// Dashboard returns the dashboard summary for the active business.
func (h *ReportHandler) Dashboard(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Today's sales
	var todaySalesTotal float64
	var todayTxCount int64
	h.DB.Model(&models.Sale{}).
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ?", businessID, todayStart).
		Count(&todayTxCount).
		Select("COALESCE(SUM(sales.total), 0)").Row().Scan(&todaySalesTotal)

	// Total stock value (selling_price × quantity)
	var totalStockValue float64
	h.DB.Model(&models.Stock{}).
		Joins("JOIN products ON products.id = stocks.product_id").
		Where("products.business_id = ?", businessID).
		Select("COALESCE(SUM(products.selling_price * stocks.quantity), 0)").
		Row().Scan(&totalStockValue)

	// Total capital invested (cost_price × quantity from stock_in movements)
	var totalCapital float64
	h.DB.Model(&models.StockMovement{}).
		Joins("JOIN products ON products.id = stock_movements.product_id").
		Where("products.business_id = ? AND stock_movements.movement_type = ?",
			businessID, models.MovementStockIn).
		Select("COALESCE(SUM(products.cost_price * stock_movements.quantity), 0)").
		Row().Scan(&totalCapital)

	// Estimated profit today (revenue - COGS)
	var todayCOGS float64
	h.DB.Model(&models.SaleItem{}).
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ?", businessID, todayStart).
		Select("COALESCE(SUM(sale_items.unit_cost * sale_items.quantity), 0)").
		Row().Scan(&todayCOGS)

	// Low stock items (first 10)
	type LowStockItem struct {
		ProductID    uint   `json:"product_id"`
		ProductTitle string `json:"product_title"`
		BranchID     uint   `json:"branch_id"`
		BranchName   string `json:"branch_name"`
		Quantity     int    `json:"quantity"`
		Threshold    int    `json:"threshold"`
	}

	var lowStockItems []LowStockItem
	h.DB.Model(&models.Stock{}).
		Select("stocks.product_id, products.title as product_title, stocks.branch_id, branches.name as branch_name, stocks.quantity, products.low_stock_threshold as threshold").
		Joins("JOIN products ON products.id = stocks.product_id").
		Joins("JOIN branches ON branches.id = stocks.branch_id").
		Where("products.business_id = ? AND stocks.quantity <= products.low_stock_threshold", businessID).
		Order("stocks.quantity ASC").
		Limit(10).
		Scan(&lowStockItems)

	// Recent sales (last 10) — preload Items + Product so the dashboard can
	// show what products were actually sold per sale, not just the cashier
	// name + total. The Preload chain is two queries total (sales + items),
	// not N+1.
	var recentSales []models.Sale
	h.DB.Preload("Branch").Preload("Cashier").
		Preload("Items.Product").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ?", businessID).
		Order("sales.created_at DESC").
		Limit(10).
		Find(&recentSales)

	type RecentSaleItem struct {
		ProductID    uint   `json:"product_id"`
		ProductTitle string `json:"product_title"`
		Quantity     int    `json:"quantity"`
	}
	type RecentSale struct {
		ID            uint             `json:"id"`
		CashierName   string           `json:"cashier_name"`
		BranchName    string           `json:"branch_name"`
		Total         float64          `json:"total"`
		PaymentMethod string           `json:"payment_method"`
		CreatedAt     time.Time        `json:"created_at"`
		ItemCount     int              `json:"item_count"`
		Items         []RecentSaleItem `json:"items"`
	}
	recent := make([]RecentSale, 0, len(recentSales))
	for _, s := range recentSales {
		items := make([]RecentSaleItem, 0, len(s.Items))
		for _, it := range s.Items {
			items = append(items, RecentSaleItem{
				ProductID:    it.ProductID,
				ProductTitle: it.Product.Title,
				Quantity:     it.Quantity,
			})
		}
		recent = append(recent, RecentSale{
			ID:            s.ID,
			CashierName:   s.Cashier.Name,
			BranchName:    s.Branch.Name,
			Total:         s.Total,
			PaymentMethod: s.PaymentMethod,
			CreatedAt:     s.CreatedAt,
			ItemCount:     len(items),
			Items:         items,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"today_sales_total":      todaySalesTotal,
			"today_transaction_count": todayTxCount,
			"total_stock_value":      totalStockValue,
			"total_capital_invested": totalCapital,
			"estimated_profit_today": todaySalesTotal - todayCOGS,
			"low_stock_items":        lowStockItems,
			"recent_sales":           recent,
		},
	})
}

// Daily returns the daily sales report.
func (h *ReportHandler) Daily(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	now := time.Now()
	fromDate := c.DefaultQuery("from_date", now.Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", now.Format("2006-01-02"))

	baseQuery := h.DB.Model(&models.Sale{}).
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59")

	if bid := c.Query("branch_id"); bid != "" {
		baseQuery = baseQuery.Where("sales.branch_id = ?", bid)
	}

	// Totals
	var totalSales float64
	var txCount int64
	baseQuery.Count(&txCount)
	baseQuery.Select("COALESCE(SUM(sales.total), 0)").Row().Scan(&totalSales)

	// By payment method
	type PaymentBreakdown struct {
		PaymentMethod string  `json:"payment_method"`
		Total         float64 `json:"total"`
		Count         int64   `json:"count"`
	}
	var paymentBreakdown []PaymentBreakdown
	h.DB.Model(&models.Sale{}).
		Select("sales.payment_method, COALESCE(SUM(sales.total), 0) as total, COUNT(*) as count").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("sales.payment_method").
		Scan(&paymentBreakdown)

	// Top products
	type TopProduct struct {
		ProductID    uint    `json:"product_id"`
		ProductTitle string  `json:"product_title"`
		QuantitySold int     `json:"quantity_sold"`
		Revenue      float64 `json:"revenue"`
	}
	var topProducts []TopProduct
	h.DB.Model(&models.SaleItem{}).
		Select("sale_items.product_id, products.title as product_title, SUM(sale_items.quantity) as quantity_sold, SUM(sale_items.line_total) as revenue").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("sale_items.product_id, products.title").
		Order("quantity_sold DESC").
		Limit(10).
		Scan(&topProducts)

	// Sales per day (chart data)
	type DailySale struct {
		Date  string  `json:"date"`
		Total float64 `json:"total"`
	}
	var salesPerDay []DailySale
	h.DB.Model(&models.Sale{}).
		Select("DATE(sales.created_at) as date, COALESCE(SUM(sales.total), 0) as total").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("DATE(sales.created_at)").
		Order("date ASC").
		Scan(&salesPerDay)

	// Per-product breakdown with cash vs mobile money split. The user's flow:
	// "I want to see ALL products sold and how much each made by cash vs MoMo".
	// Single SUM(CASE WHEN ...) aggregation per product, no extra round-trips.
	type ProductBreakdown struct {
		ProductID    uint    `json:"product_id"`
		ProductTitle string  `json:"product_title"`
		CashQty      int     `json:"cash_qty"`
		CashRevenue  float64 `json:"cash_revenue"`
		MomoQty      int     `json:"momo_qty"`
		MomoRevenue  float64 `json:"momo_revenue"`
		TotalQty     int     `json:"total_qty"`
		TotalRevenue float64 `json:"total_revenue"`
	}
	var productsBreakdown []ProductBreakdown
	h.DB.Model(&models.SaleItem{}).
		Select(`
			sale_items.product_id,
			products.title as product_title,
			COALESCE(SUM(CASE WHEN sales.payment_method = 'cash' THEN sale_items.quantity ELSE 0 END), 0) as cash_qty,
			COALESCE(SUM(CASE WHEN sales.payment_method = 'cash' THEN sale_items.line_total ELSE 0 END), 0) as cash_revenue,
			COALESCE(SUM(CASE WHEN sales.payment_method = 'mobile_money' THEN sale_items.quantity ELSE 0 END), 0) as momo_qty,
			COALESCE(SUM(CASE WHEN sales.payment_method = 'mobile_money' THEN sale_items.line_total ELSE 0 END), 0) as momo_revenue,
			COALESCE(SUM(sale_items.quantity), 0) as total_qty,
			COALESCE(SUM(sale_items.line_total), 0) as total_revenue
		`).
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("sale_items.product_id, products.title").
		Order("total_revenue DESC").
		Scan(&productsBreakdown)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_sales":       totalSales,
			"transaction_count": txCount,
			"by_payment_method": paymentBreakdown,
			"top_products":      topProducts,
			"products_breakdown": productsBreakdown,
			"sales_per_day":     salesPerDay,
		},
	})
}

// StockReport returns stock levels report.
func (h *ReportHandler) StockReport(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	query := h.DB.Model(&models.Stock{}).
		Select("stocks.product_id, products.title as product_title, categories.name as category_name, stocks.branch_id, branches.name as branch_name, stocks.quantity, products.low_stock_threshold as threshold, products.selling_price").
		Joins("JOIN products ON products.id = stocks.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Joins("JOIN branches ON branches.id = stocks.branch_id").
		Where("products.business_id = ?", businessID)

	if bid := c.Query("branch_id"); bid != "" {
		query = query.Where("stocks.branch_id = ?", bid)
	}
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("products.category_id = ?", catID)
	}

	type StockReportItem struct {
		ProductID    uint    `json:"product_id"`
		ProductTitle string  `json:"product_title"`
		CategoryName string  `json:"category_name"`
		BranchID     uint    `json:"branch_id"`
		BranchName   string  `json:"branch_name"`
		Quantity     int     `json:"quantity"`
		Threshold    int     `json:"threshold"`
		SellingPrice float64 `json:"selling_price"`
		StockValue   float64 `json:"stock_value"`
		Status       string  `json:"status"`
	}

	var items []StockReportItem
	query.Order("products.title ASC, branches.name ASC").Scan(&items)

	// Calculate stock value and status, filter by status if needed
	statusFilter := c.Query("status")
	results := make([]StockReportItem, 0, len(items))
	for i := range items {
		items[i].StockValue = items[i].SellingPrice * float64(items[i].Quantity)
		if items[i].Quantity == 0 {
			items[i].Status = "out"
		} else if items[i].Quantity <= items[i].Threshold {
			items[i].Status = "low"
		} else {
			items[i].Status = "ok"
		}

		if statusFilter == "" || statusFilter == items[i].Status {
			results = append(results, items[i])
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// PnL returns profit & loss / capital report (admin only).
func (h *ReportHandler) PnL(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view P&L reports"})
		return
	}

	now := time.Now()
	fromDate := c.DefaultQuery("from_date", now.AddDate(0, -1, 0).Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", now.Format("2006-01-02"))

	// Total capital invested (all-time stock-in cost)
	var totalCapital float64
	h.DB.Model(&models.StockMovement{}).
		Joins("JOIN products ON products.id = stock_movements.product_id").
		Where("products.business_id = ? AND stock_movements.movement_type = ?",
			businessID, models.MovementStockIn).
		Select("COALESCE(SUM(products.cost_price * stock_movements.quantity), 0)").
		Row().Scan(&totalCapital)

	// Revenue in range
	var totalRevenue float64
	h.DB.Model(&models.Sale{}).
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Select("COALESCE(SUM(sales.total), 0)").
		Row().Scan(&totalRevenue)

	// COGS in range
	var totalCOGS float64
	h.DB.Model(&models.SaleItem{}).
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Select("COALESCE(SUM(sale_items.unit_cost * sale_items.quantity), 0)").
		Row().Scan(&totalCOGS)

	grossProfit := totalRevenue - totalCOGS
	var grossMargin float64
	if totalRevenue > 0 {
		grossMargin = (grossProfit / totalRevenue) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"from_date":              fromDate,
			"to_date":                toDate,
			"total_capital_invested": totalCapital,
			"total_revenue":          totalRevenue,
			"total_cogs":             totalCOGS,
			"gross_profit":           grossProfit,
			"gross_margin_percent":   grossMargin,
		},
	})
}

// Cashiers returns sales by cashier/employee report (admin only).
func (h *ReportHandler) Cashiers(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view cashier reports"})
		return
	}

	now := time.Now()
	fromDate := c.DefaultQuery("from_date", now.Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", now.Format("2006-01-02"))

	query := h.DB.Model(&models.Sale{}).
		Select("sales.cashier_id as user_id, users.name, branches.name as branch_name, COALESCE(SUM(sales.total), 0) as total_sales, COUNT(*) as transaction_count").
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Joins("JOIN users ON users.id = sales.cashier_id").
		Where("branches.business_id = ? AND sales.created_at >= ? AND sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("sales.cashier_id, users.name, branches.name")

	if bid := c.Query("branch_id"); bid != "" {
		query = query.Where("sales.branch_id = ?", bid)
	}

	type CashierStat struct {
		UserID           uint    `json:"user_id"`
		Name             string  `json:"name"`
		BranchName       string  `json:"branch_name"`
		TotalSales       float64 `json:"total_sales"`
		TransactionCount int64   `json:"transaction_count"`
	}

	var stats []CashierStat
	query.Order("total_sales DESC").Scan(&stats)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
