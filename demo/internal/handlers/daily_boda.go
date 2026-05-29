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

type DailyBodaHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// ---------- Drivers ----------

func (h *DailyBodaHandler) ListDrivers(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Branch").Where("business_id = ?", businessID)
	if active := c.Query("active"); active == "true" {
		q = q.Where("is_active = ?", true)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(full_name) LIKE ? OR phone LIKE ? OR national_id LIKE ?", like, "%"+search+"%", "%"+search+"%")
	}
	var rows []models.DailyBodaDriver
	q.Order("full_name ASC").Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (h *DailyBodaHandler) GetDriver(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var d models.DailyBodaDriver
	if err := h.DB.Preload("Branch").Where("id = ? AND business_id = ?", id, businessID).First(&d).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (h *DailyBodaHandler) CreateDriver(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageDailyBoda(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	var req struct {
		BranchID   uint    `json:"branch_id" binding:"required"`
		FullName   string  `json:"full_name" binding:"required"`
		Phone      string  `json:"phone" binding:"required"`
		NationalID string  `json:"national_id"`
		Address    string  `json:"address"`
		DailyRate  float64 `json:"daily_rate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch, full_name and phone are required"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}
	if req.DailyRate == 0 {
		req.DailyRate = models.DefaultDailyBodaRate
	}
	d := models.DailyBodaDriver{
		BusinessID: businessID,
		BranchID:   req.BranchID,
		FullName:   strings.TrimSpace(req.FullName),
		Phone:      strings.TrimSpace(req.Phone),
		NationalID: req.NationalID,
		Address:    req.Address,
		DailyRate:  req.DailyRate,
		IsActive:   true,
	}
	if err := h.DB.Create(&d).Error; err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "A driver with that phone already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver"})
		return
	}
	h.Audit.Log(c, models.ActionCreate, "daily_boda_driver", &d.ID, fmt.Sprintf("Added daily-boda driver %s (%s)", d.FullName, d.Phone), gin.H{"branch_id": d.BranchID, "daily_rate": d.DailyRate})
	c.JSON(http.StatusCreated, gin.H{"data": d})
}

func (h *DailyBodaHandler) UpdateDriver(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageDailyBoda(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var d models.DailyBodaDriver
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&d).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	allowed := map[string]bool{
		"branch_id": true, "full_name": true, "phone": true, "national_id": true,
		"address": true, "daily_rate": true, "is_active": true,
	}
	updates := map[string]interface{}{}
	for k, v := range req {
		if allowed[k] {
			updates[k] = v
		}
	}
	if err := h.DB.Model(&d).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver"})
		return
	}
	h.Audit.Log(c, models.ActionUpdate, "daily_boda_driver", &d.ID, fmt.Sprintf("Updated daily-boda driver %s", d.FullName), nil)
	c.JSON(http.StatusOK, gin.H{"data": d})
}

// ---------- Motorcycles ----------

func (h *DailyBodaHandler) ListMotorcycles(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Branch").Where("business_id = ?", businessID)
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	var rows []models.DailyBodaMotorcycle
	q.Order("number_plate ASC").Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (h *DailyBodaHandler) CreateMotorcycle(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageDailyBoda(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	var req struct {
		BranchID    uint   `json:"branch_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		NumberPlate string `json:"number_plate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch, name and number_plate are required"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}
	m := models.DailyBodaMotorcycle{
		BusinessID:  businessID,
		BranchID:    req.BranchID,
		Name:        strings.TrimSpace(req.Name),
		NumberPlate: strings.ToUpper(strings.ReplaceAll(req.NumberPlate, " ", "")),
		Status:      models.DailyBodaStatusAvailable,
	}
	if err := h.DB.Create(&m).Error; err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "A motorcycle with that number plate already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create motorcycle"})
		return
	}
	h.Audit.Log(c, models.ActionCreate, "daily_boda_motorcycle", &m.ID, fmt.Sprintf("Added daily-boda motorcycle %q (%s)", m.Name, m.NumberPlate), gin.H{"branch_id": m.BranchID})
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

// AssignDriver pairs a driver with a motorcycle. Both must be free.
func (h *DailyBodaHandler) AssignDriver(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageDailyBoda(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	motoID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		DriverID uint `json:"driver_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "driver_id is required"})
		return
	}

	var m models.DailyBodaMotorcycle
	var assignedDriver models.DailyBodaDriver
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND business_id = ?", motoID, businessID).First(&m).Error; err != nil {
			return fmt.Errorf("motorcycle not found")
		}
		if m.Status != models.DailyBodaStatusAvailable {
			return fmt.Errorf("motorcycle is not available (status: %s)", m.Status)
		}
		if err := tx.Where("id = ? AND business_id = ?", req.DriverID, businessID).First(&assignedDriver).Error; err != nil {
			return fmt.Errorf("driver not found")
		}
		// Driver must not already be on another bike.
		var occupied int64
		tx.Model(&models.DailyBodaMotorcycle{}).
			Where("assigned_driver_id = ? AND business_id = ?", req.DriverID, businessID).
			Count(&occupied)
		if occupied > 0 {
			return fmt.Errorf("driver is already assigned to another motorcycle")
		}
		return tx.Model(&m).Updates(map[string]interface{}{
			"assigned_driver_id": req.DriverID,
			"status":             models.DailyBodaStatusOccupied,
		}).Error
	})
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	h.Audit.Log(c, models.ActionUpdate, "daily_boda_motorcycle", &m.ID, fmt.Sprintf("Assigned driver %s to motorcycle %q (%s)", assignedDriver.FullName, m.Name, m.NumberPlate), gin.H{"driver_id": req.DriverID})
	c.JSON(http.StatusOK, gin.H{"message": "Driver assigned"})
}

// ReturnMotorcycle marks a motorcycle as returned (free of driver).
func (h *DailyBodaHandler) ReturnMotorcycle(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageDailyBoda(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	motoID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var m models.DailyBodaMotorcycle
	if err := h.DB.Where("id = ? AND business_id = ?", motoID, businessID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motorcycle not found"})
		return
	}
	h.DB.Model(&m).Updates(map[string]interface{}{
		"assigned_driver_id": nil,
		"status":             models.DailyBodaStatusAvailable,
	})
	h.Audit.Log(c, models.ActionUpdate, "daily_boda_motorcycle", &m.ID, fmt.Sprintf("Returned motorcycle %q (%s) from driver", m.Name, m.NumberPlate), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Motorcycle returned"})
}

// ---------- Payments ----------

func (h *DailyBodaHandler) ListPayments(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Driver").Preload("Motorcycle").Preload("Branch").
		Where("business_id = ?", businessID)
	if driverID := c.Query("driver_id"); driverID != "" {
		q = q.Where("driver_id = ?", driverID)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if from := c.Query("from"); from != "" {
		q = q.Where("payment_date >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("payment_date <= ?", to)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 500 {
		perPage = 100
	}

	var total int64
	q.Model(&models.DailyBodaPayment{}).Count(&total)

	var rows []models.DailyBodaPayment
	q.Order("payment_date DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "per_page": perPage})
}

// CreatePayment records a daily collection from a driver. Status auto-set
// based on whether amount covers the daily rate.
func (h *DailyBodaHandler) CreatePayment(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canCollectDailyBoda(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		DriverID      uint       `json:"driver_id" binding:"required"`
		MotorcycleID  uint       `json:"motorcycle_id" binding:"required"`
		BranchID      uint       `json:"branch_id" binding:"required"`
		Amount        float64    `json:"amount" binding:"required"`
		DailyRate     float64    `json:"daily_rate"`
		PaymentDate   *time.Time `json:"payment_date"`
		PaymentMethod string     `json:"payment_method"`
		Notes         string     `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields missing"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}
	if req.DailyRate == 0 {
		// Default to driver's configured rate.
		var d models.DailyBodaDriver
		if err := h.DB.Where("id = ? AND business_id = ?", req.DriverID, businessID).First(&d).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Driver not found"})
			return
		}
		req.DailyRate = d.DailyRate
	}
	if req.PaymentDate == nil {
		now := time.Now()
		req.PaymentDate = &now
	}

	balance := req.DailyRate - req.Amount
	if balance < 0 {
		balance = 0
	}
	status := models.DailyBodaPaymentPaid
	if req.Amount <= 0 {
		status = models.DailyBodaPaymentPending
	} else if req.Amount < req.DailyRate {
		status = models.DailyBodaPaymentPartial
	}

	p := models.DailyBodaPayment{
		BusinessID:    businessID,
		DriverID:      req.DriverID,
		MotorcycleID:  req.MotorcycleID,
		BranchID:      req.BranchID,
		Amount:        req.Amount,
		DailyRate:     req.DailyRate,
		Balance:       balance,
		PaymentDate:   *req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
		Status:        status,
		CollectedBy:   userID,
		Notes:         req.Notes,
	}
	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment"})
		return
	}
	h.Audit.Log(c, models.ActionCreate, "daily_boda_payment", &p.ID,
		fmt.Sprintf("Collected daily-boda payment UGX %s (status: %s)", formatMoney(p.Amount), p.Status),
		gin.H{"driver_id": p.DriverID, "motorcycle_id": p.MotorcycleID, "branch_id": p.BranchID, "amount": p.Amount, "daily_rate": p.DailyRate, "balance": p.Balance})
	c.JSON(http.StatusCreated, gin.H{"data": p, "segment": models.SegmentDailyBoda})
}

// VerifyPayment is the manager's sign-off on a collected payment.
func (h *DailyBodaHandler) VerifyPayment(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if !hasAnyRole(role, models.RoleAdmin, models.RoleManager, models.RoleAccountant) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var p models.DailyBodaPayment
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	if p.VerifiedAt != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Payment already verified"})
		return
	}
	now := time.Now()
	h.DB.Model(&p).Updates(map[string]interface{}{
		"verified_by": userID,
		"verified_at": now,
	})
	h.Audit.Log(c, models.ActionVerify, "daily_boda_payment", &p.ID, fmt.Sprintf("Verified daily-boda payment UGX %s", formatMoney(p.Amount)), gin.H{"driver_id": p.DriverID, "amount": p.Amount})
	c.JSON(http.StatusOK, gin.H{"message": "Payment verified"})
}

func canManageDailyBoda(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager)
}

func canCollectDailyBoda(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager, models.RoleCashier, models.RoleLoanOfficer)
}
