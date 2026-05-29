package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

type CashSaleHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// List returns motorcycle cash sales scoped to the business.
func (h *CashSaleHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Motorcycle").Preload("Branch").Preload("Seller").
		Where("business_id = ?", businessID)

	if branchID := c.Query("branch_id"); branchID != "" {
		q = q.Where("branch_id = ?", branchID)
	}
	if from := c.Query("from"); from != "" {
		q = q.Where("created_at >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("created_at <= ?", to)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(customer_name) LIKE ? OR customer_phone LIKE ? OR sale_number LIKE ?", like, "%"+search+"%", "%"+strings.ToUpper(search)+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 200 {
		perPage = 50
	}

	var total int64
	q.Model(&models.CashSale{}).Count(&total)

	var rows []models.CashSale
	q.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "per_page": perPage})
}

func (h *CashSaleHandler) Get(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var s models.CashSale
	if err := h.DB.Preload("Motorcycle").Preload("Branch").Preload("Seller").
		Where("id = ? AND business_id = ?", id, businessID).First(&s).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cash sale not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

// Create records a cash sale of a motorcycle. Marks the motorcycle as sold in
// the same transaction so we can't double-sell.
func (h *CashSaleHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canSellMotorcycles(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		BranchID        uint    `json:"branch_id" binding:"required"`
		MotorcycleID    uint    `json:"motorcycle_id" binding:"required"`
		CustomerName    string  `json:"customer_name" binding:"required"`
		CustomerPhone   string  `json:"customer_phone"`
		CustomerEmail   string  `json:"customer_email"`
		CustomerNIN     string  `json:"customer_nin"`
		CustomerAddress string  `json:"customer_address"`
		ListPrice       float64 `json:"list_price"`
		DiscountAmount  float64 `json:"discount_amount"`
		Total           float64 `json:"total" binding:"required"`
		PaymentMethod   string  `json:"payment_method" binding:"required"`
		TransactionRef  string  `json:"transaction_ref"`
		Notes           string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields missing"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}

	var sale models.CashSale
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var moto models.Motorcycle
		if err := tx.Where("id = ? AND business_id = ?", req.MotorcycleID, businessID).First(&moto).Error; err != nil {
			return fmt.Errorf("motorcycle not found")
		}
		if moto.Status != models.MotorcycleStatusAvailable && moto.Status != models.MotorcycleStatusReserved {
			return fmt.Errorf("motorcycle is not available for sale (status: %s)", moto.Status)
		}

		sale = models.CashSale{
			BusinessID:      businessID,
			BranchID:        req.BranchID,
			MotorcycleID:    req.MotorcycleID,
			SoldBy:          userID,
			CustomerName:    strings.TrimSpace(req.CustomerName),
			CustomerPhone:   req.CustomerPhone,
			CustomerEmail:   req.CustomerEmail,
			CustomerNIN:     req.CustomerNIN,
			CustomerAddress: req.CustomerAddress,
			ListPrice:       req.ListPrice,
			DiscountAmount:  req.DiscountAmount,
			Total:           req.Total,
			PaymentMethod:   req.PaymentMethod,
			TransactionRef:  req.TransactionRef,
			Notes:           req.Notes,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return err
		}
		sale.SaleNumber = fmt.Sprintf("KM-MC-%06d", sale.ID)
		if err := tx.Model(&sale).Update("sale_number", sale.SaleNumber).Error; err != nil {
			return err
		}
		// Mark motorcycle sold so it can't be sold or financed again.
		return tx.Model(&moto).Update("status", models.MotorcycleStatusSold).Error
	})
	if err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Motorcycle has already been sold"})
			return
		}
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	h.DB.Preload("Motorcycle").Preload("Branch").Preload("Seller").First(&sale, sale.ID)
	h.Audit.Log(c, models.ActionCreate, "cash_sale", &sale.ID,
		fmt.Sprintf("Cash sale %s: %q for UGX %s", sale.SaleNumber, sale.Motorcycle.Name, formatMoney(sale.Total)),
		gin.H{"sale_number": sale.SaleNumber, "motorcycle_id": sale.MotorcycleID, "branch_id": sale.BranchID, "total": sale.Total, "payment_method": sale.PaymentMethod})
	// Stamp segment-aware response so dashboard can roll up correctly.
	c.JSON(http.StatusCreated, gin.H{"data": sale, "segment": models.SegmentMotorcyclesCash})
	_ = time.Now()
}

func canSellMotorcycles(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager, models.RoleCashier)
}
