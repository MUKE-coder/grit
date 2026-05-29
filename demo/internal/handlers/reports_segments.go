package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// SegmentReportHandler hosts cross-segment financial reports for Grit Motors.
// Spares-side reports continue to live on ReportHandler — these are the
// motorcycle-and-loans counterparts plus a combined dashboard roll-up.
type SegmentReportHandler struct {
	DB *gorm.DB
}

// Main returns a single payload that drives the combined home dashboard. Replaces
// the 5-call fan-out the front-end was doing (spares dashboard + 3× loans + 1×
// motorcycles). One round trip, one query plan per metric.
func (h *SegmentReportHandler) Main(c *gin.Context) {
	businessID := c.GetUint("business_id")
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	out := gin.H{}

	// --- Spares (today) ---
	var sparesToday float64
	var sparesTxCount int64
	h.DB.Model(&models.Sale{}).
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("branches.business_id = ? AND sales.created_at >= ?", businessID, todayStart).
		Count(&sparesTxCount).
		Select("COALESCE(SUM(sales.total), 0)").Row().Scan(&sparesToday)

	// --- Motorcycles ---
	var motoCounts []struct {
		Status string
		Count  int64
	}
	h.DB.Model(&models.Motorcycle{}).
		Select("status, COUNT(*) as count").
		Where("business_id = ?", businessID).
		Group("status").
		Scan(&motoCounts)

	motoStatusMap := map[string]int64{}
	for _, r := range motoCounts {
		motoStatusMap[r.Status] = r.Count
	}

	var inventoryValue float64
	h.DB.Model(&models.Motorcycle{}).
		Where("business_id = ? AND status = ?", businessID, models.MotorcycleStatusAvailable).
		Select("COALESCE(SUM(selling_price), 0)").Row().Scan(&inventoryValue)

	// --- Cash sales today ---
	var cashSalesTodayTotal float64
	var cashSalesTodayCount int64
	h.DB.Model(&models.CashSale{}).
		Where("business_id = ? AND created_at >= ?", businessID, todayStart).
		Count(&cashSalesTodayCount).
		Select("COALESCE(SUM(total), 0)").Row().Scan(&cashSalesTodayTotal)

	// --- Loans ---
	var loanCounts []struct {
		Status string
		Count  int64
	}
	h.DB.Model(&models.Loan{}).
		Select("status, COUNT(*) as count").
		Where("business_id = ?", businessID).
		Group("status").
		Scan(&loanCounts)

	loanStatusMap := map[string]int64{}
	for _, r := range loanCounts {
		loanStatusMap[r.Status] = r.Count
	}

	var totalOutstanding float64
	h.DB.Model(&models.Loan{}).
		Where("business_id = ? AND status = ?", businessID, models.LoanStatusActive).
		Select("COALESCE(SUM(balance_remaining), 0)").Row().Scan(&totalOutstanding)

	var disbursedTotal float64
	h.DB.Model(&models.Loan{}).
		Where("business_id = ? AND status IN ?", businessID, []string{models.LoanStatusActive, models.LoanStatusCompleted}).
		Select("COALESCE(SUM(disbursed_amount), 0)").Row().Scan(&disbursedTotal)

	// --- Repayments today (only approved count toward revenue) ---
	var repaymentsTodayTotal float64
	var repaymentsTodayCount int64
	h.DB.Model(&models.Repayment{}).
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("loans.business_id = ? AND repayments.collection_date >= ? AND repayments.status = ?",
			businessID, todayStart, models.RepaymentStatusApproved).
		Count(&repaymentsTodayCount).
		Select("COALESCE(SUM(repayments.amount), 0)").Row().Scan(&repaymentsTodayTotal)

	var pendingRepaymentsCount int64
	h.DB.Model(&models.Repayment{}).
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("loans.business_id = ? AND repayments.status = ?", businessID, models.RepaymentStatusPending).
		Count(&pendingRepaymentsCount)

	// --- Overdue installments ---
	var overdueCount int64
	h.DB.Model(&models.RepaymentSchedule{}).
		Joins("JOIN loans ON loans.id = repayment_schedules.loan_id").
		Where("loans.business_id = ? AND repayment_schedules.is_overdue = ? AND repayment_schedules.is_paid = ?", businessID, true, false).
		Count(&overdueCount)

	// --- Daily Boda today ---
	var dailyBodaTodayTotal float64
	var dailyBodaTodayCount int64
	h.DB.Model(&models.DailyBodaPayment{}).
		Where("business_id = ? AND payment_date >= ?", businessID, todayStart).
		Count(&dailyBodaTodayCount).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dailyBodaTodayTotal)

	// Combined revenue today across all 4 segments.
	totalToday := sparesToday + cashSalesTodayTotal + repaymentsTodayTotal + dailyBodaTodayTotal

	out["spares"] = gin.H{
		"today_total":             sparesToday,
		"today_transaction_count": sparesTxCount,
	}
	out["motorcycles"] = gin.H{
		"available":         motoStatusMap[models.MotorcycleStatusAvailable],
		"reserved":          motoStatusMap[models.MotorcycleStatusReserved],
		"sold":              motoStatusMap[models.MotorcycleStatusSold],
		"on_loan":           motoStatusMap[models.MotorcycleStatusOnLoan],
		"repossessed":       motoStatusMap[models.MotorcycleStatusRepossessed],
		"inventory_value":   inventoryValue,
		"cash_sales_today":  cashSalesTodayTotal,
		"cash_sales_today_count": cashSalesTodayCount,
	}
	out["loans"] = gin.H{
		"pending":          loanStatusMap[models.LoanStatusPending],
		"approved":         loanStatusMap[models.LoanStatusApproved],
		"active":           loanStatusMap[models.LoanStatusActive],
		"completed":        loanStatusMap[models.LoanStatusCompleted],
		"defaulted":        loanStatusMap[models.LoanStatusDefaulted],
		"total_outstanding": totalOutstanding,
		"total_disbursed":  disbursedTotal,
	}
	out["repayments"] = gin.H{
		"collected_today":         repaymentsTodayTotal,
		"collected_today_count":   repaymentsTodayCount,
		"pending_verification":    pendingRepaymentsCount,
		"overdue_installments":    overdueCount,
	}
	out["daily_boda"] = gin.H{
		"today_total": dailyBodaTodayTotal,
		"today_count": dailyBodaTodayCount,
	}
	out["combined"] = gin.H{
		"today_total": totalToday,
		"by_segment": gin.H{
			"spares":          sparesToday,
			"motorcycles_cash": cashSalesTodayTotal,
			"loans":           repaymentsTodayTotal,
			"daily_boda":      dailyBodaTodayTotal,
		},
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// LoansReport returns portfolio-level metrics for the loans book.
func (h *SegmentReportHandler) LoansReport(c *gin.Context) {
	businessID := c.GetUint("business_id")

	// By status
	type StatusBucket struct {
		Status string  `json:"status"`
		Count  int64   `json:"count"`
		Value  float64 `json:"value"` // sum of principal for visibility into book size
	}
	var byStatus []StatusBucket
	h.DB.Model(&models.Loan{}).
		Select("status, COUNT(*) as count, COALESCE(SUM(principal_amount), 0) as value").
		Where("business_id = ?", businessID).
		Group("status").
		Scan(&byStatus)

	// By branch
	type BranchBucket struct {
		BranchID         uint    `json:"branch_id"`
		BranchName       string  `json:"branch_name"`
		ActiveLoans      int64   `json:"active_loans"`
		TotalOutstanding float64 `json:"total_outstanding"`
	}
	var byBranch []BranchBucket
	h.DB.Model(&models.Loan{}).
		Select("loans.branch_id, branches.name as branch_name, COUNT(*) as active_loans, COALESCE(SUM(loans.balance_remaining), 0) as total_outstanding").
		Joins("JOIN branches ON branches.id = loans.branch_id").
		Where("loans.business_id = ? AND loans.status = ?", businessID, models.LoanStatusActive).
		Group("loans.branch_id, branches.name").
		Order("total_outstanding DESC").
		Scan(&byBranch)

	// Top borrowers by outstanding balance
	type TopBorrower struct {
		BorrowerID  uint    `json:"borrower_id"`
		FullName    string  `json:"full_name"`
		Phone       string  `json:"phone"`
		Outstanding float64 `json:"outstanding"`
		LoanCount   int64   `json:"loan_count"`
	}
	var topBorrowers []TopBorrower
	h.DB.Model(&models.Loan{}).
		Select("loans.borrower_id, (borrowers.first_name || ' ' || borrowers.last_name) as full_name, borrowers.phone, COALESCE(SUM(loans.balance_remaining), 0) as outstanding, COUNT(*) as loan_count").
		Joins("JOIN borrowers ON borrowers.id = loans.borrower_id").
		Where("loans.business_id = ? AND loans.status = ?", businessID, models.LoanStatusActive).
		Group("loans.borrower_id, full_name, borrowers.phone").
		Order("outstanding DESC").
		Limit(10).
		Scan(&topBorrowers)

	// Overdue summary
	type OverdueRow struct {
		LoanID       uint    `json:"loan_id"`
		LoanNumber   string  `json:"loan_number"`
		BorrowerName string  `json:"borrower_name"`
		DaysPastDue  int     `json:"days_past_due"`
		AmountDue    float64 `json:"amount_due"`
	}
	var overdue []OverdueRow
	h.DB.Model(&models.RepaymentSchedule{}).
		Select(`
			repayment_schedules.loan_id,
			loans.loan_number,
			(borrowers.first_name || ' ' || borrowers.last_name) as borrower_name,
			MAX(repayment_schedules.days_past_due) as days_past_due,
			COALESCE(SUM(repayment_schedules.total_amount - repayment_schedules.paid_amount), 0) as amount_due
		`).
		Joins("JOIN loans ON loans.id = repayment_schedules.loan_id").
		Joins("JOIN borrowers ON borrowers.id = loans.borrower_id").
		Where("loans.business_id = ? AND repayment_schedules.is_overdue = ? AND repayment_schedules.is_paid = ?", businessID, true, false).
		Group("repayment_schedules.loan_id, loans.loan_number, borrower_name").
		Order("days_past_due DESC").
		Limit(50).
		Scan(&overdue)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"by_status":     byStatus,
			"by_branch":     byBranch,
			"top_borrowers": topBorrowers,
			"overdue":       overdue,
		},
	})
}

// CollectionsReport summarizes repayment collections in a date range.
func (h *SegmentReportHandler) CollectionsReport(c *gin.Context) {
	businessID := c.GetUint("business_id")
	now := time.Now()
	fromDate := c.DefaultQuery("from_date", now.AddDate(0, 0, -30).Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", now.Format("2006-01-02"))

	base := h.DB.Model(&models.Repayment{}).
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("loans.business_id = ? AND repayments.collection_date >= ? AND repayments.collection_date <= ? AND repayments.status = ?",
			businessID, fromDate, toDate+" 23:59:59", models.RepaymentStatusApproved)

	// Total collected
	var totalCollected float64
	var totalCount int64
	base.Count(&totalCount)
	base.Select("COALESCE(SUM(repayments.amount), 0)").Row().Scan(&totalCollected)

	// By payment method
	type MethodBucket struct {
		PaymentMethod string  `json:"payment_method"`
		Total         float64 `json:"total"`
		Count         int64   `json:"count"`
	}
	var byMethod []MethodBucket
	h.DB.Model(&models.Repayment{}).
		Select("repayments.payment_method, COALESCE(SUM(repayments.amount), 0) as total, COUNT(*) as count").
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("loans.business_id = ? AND repayments.collection_date >= ? AND repayments.collection_date <= ? AND repayments.status = ?",
			businessID, fromDate, toDate+" 23:59:59", models.RepaymentStatusApproved).
		Group("repayments.payment_method").
		Scan(&byMethod)

	// By collector
	type CollectorBucket struct {
		UserID uint    `json:"user_id"`
		Name   string  `json:"name"`
		Total  float64 `json:"total"`
		Count  int64   `json:"count"`
	}
	var byCollector []CollectorBucket
	h.DB.Model(&models.Repayment{}).
		Select("repayments.collected_by as user_id, users.name, COALESCE(SUM(repayments.amount), 0) as total, COUNT(*) as count").
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Joins("JOIN users ON users.id = repayments.collected_by").
		Where("loans.business_id = ? AND repayments.collection_date >= ? AND repayments.collection_date <= ? AND repayments.status = ?",
			businessID, fromDate, toDate+" 23:59:59", models.RepaymentStatusApproved).
		Group("repayments.collected_by, users.name").
		Order("total DESC").
		Scan(&byCollector)

	// Daily series
	type DailyRow struct {
		Date  string  `json:"date"`
		Total float64 `json:"total"`
	}
	var daily []DailyRow
	h.DB.Model(&models.Repayment{}).
		Select("DATE(repayments.collection_date) as date, COALESCE(SUM(repayments.amount), 0) as total").
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("loans.business_id = ? AND repayments.collection_date >= ? AND repayments.collection_date <= ? AND repayments.status = ?",
			businessID, fromDate, toDate+" 23:59:59", models.RepaymentStatusApproved).
		Group("DATE(repayments.collection_date)").
		Order("date ASC").
		Scan(&daily)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"from_date":         fromDate,
			"to_date":           toDate,
			"total_collected":   totalCollected,
			"transaction_count": totalCount,
			"by_method":         byMethod,
			"by_collector":      byCollector,
			"daily":             daily,
		},
	})
}

// MotorcyclesReport summarizes motorcycle inventory + cash-sale activity.
func (h *SegmentReportHandler) MotorcyclesReport(c *gin.Context) {
	businessID := c.GetUint("business_id")
	now := time.Now()
	fromDate := c.DefaultQuery("from_date", now.AddDate(0, 0, -30).Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", now.Format("2006-01-02"))

	// Inventory by branch + status
	type InventoryBucket struct {
		BranchID   uint   `json:"branch_id"`
		BranchName string `json:"branch_name"`
		Status     string `json:"status"`
		Count      int64  `json:"count"`
	}
	var inventory []InventoryBucket
	h.DB.Model(&models.Motorcycle{}).
		Select("motorcycles.branch_id, branches.name as branch_name, motorcycles.status, COUNT(*) as count").
		Joins("JOIN branches ON branches.id = motorcycles.branch_id").
		Where("motorcycles.business_id = ?", businessID).
		Group("motorcycles.branch_id, branches.name, motorcycles.status").
		Order("branches.name ASC").
		Scan(&inventory)

	// Cash sales summary in range
	var cashSalesTotal float64
	var cashSalesCount int64
	h.DB.Model(&models.CashSale{}).
		Where("business_id = ? AND created_at >= ? AND created_at <= ?", businessID, fromDate, toDate+" 23:59:59").
		Count(&cashSalesCount).
		Select("COALESCE(SUM(total), 0)").Row().Scan(&cashSalesTotal)

	// Cash sales by branch
	type CashSalesByBranch struct {
		BranchID   uint    `json:"branch_id"`
		BranchName string  `json:"branch_name"`
		Total      float64 `json:"total"`
		Count      int64   `json:"count"`
	}
	var cashByBranch []CashSalesByBranch
	h.DB.Model(&models.CashSale{}).
		Select("cash_sales.branch_id, branches.name as branch_name, COALESCE(SUM(cash_sales.total), 0) as total, COUNT(*) as count").
		Joins("JOIN branches ON branches.id = cash_sales.branch_id").
		Where("cash_sales.business_id = ? AND cash_sales.created_at >= ? AND cash_sales.created_at <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("cash_sales.branch_id, branches.name").
		Order("total DESC").
		Scan(&cashByBranch)

	// Loan-financed motorcycles in range
	var loanSalesTotal float64
	var loanSalesCount int64
	h.DB.Model(&models.Loan{}).
		Where("business_id = ? AND motorcycle_id IS NOT NULL AND status IN ? AND created_at >= ? AND created_at <= ?",
			businessID, []string{models.LoanStatusActive, models.LoanStatusCompleted}, fromDate, toDate+" 23:59:59").
		Count(&loanSalesCount).
		Select("COALESCE(SUM(principal_amount), 0)").Row().Scan(&loanSalesTotal)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"from_date":          fromDate,
			"to_date":            toDate,
			"inventory":          inventory,
			"cash_sales_total":   cashSalesTotal,
			"cash_sales_count":   cashSalesCount,
			"cash_sales_by_branch": cashByBranch,
			"loan_sales_total":   loanSalesTotal,
			"loan_sales_count":   loanSalesCount,
			"motorcycles_moved":  cashSalesCount + loanSalesCount,
		},
	})
}

// DailyBodaReport summarizes the daily-rental fleet's activity.
func (h *SegmentReportHandler) DailyBodaReport(c *gin.Context) {
	businessID := c.GetUint("business_id")
	now := time.Now()
	fromDate := c.DefaultQuery("from_date", now.AddDate(0, 0, -30).Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", now.Format("2006-01-02"))

	var totalCollected float64
	var totalCount int64
	h.DB.Model(&models.DailyBodaPayment{}).
		Where("business_id = ? AND payment_date >= ? AND payment_date <= ?", businessID, fromDate, toDate+" 23:59:59").
		Count(&totalCount).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalCollected)

	var totalBalance float64
	h.DB.Model(&models.DailyBodaPayment{}).
		Where("business_id = ? AND payment_date >= ? AND payment_date <= ? AND balance > 0", businessID, fromDate, toDate+" 23:59:59").
		Select("COALESCE(SUM(balance), 0)").Row().Scan(&totalBalance)

	type DriverStats struct {
		DriverID    uint    `json:"driver_id"`
		FullName    string  `json:"full_name"`
		Phone       string  `json:"phone"`
		TotalPaid   float64 `json:"total_paid"`
		PaymentCount int64  `json:"payment_count"`
	}
	var topDrivers []DriverStats
	h.DB.Model(&models.DailyBodaPayment{}).
		Select("daily_boda_payments.driver_id, daily_boda_drivers.full_name, daily_boda_drivers.phone, COALESCE(SUM(daily_boda_payments.amount), 0) as total_paid, COUNT(*) as payment_count").
		Joins("JOIN daily_boda_drivers ON daily_boda_drivers.id = daily_boda_payments.driver_id").
		Where("daily_boda_payments.business_id = ? AND daily_boda_payments.payment_date >= ? AND daily_boda_payments.payment_date <= ?",
			businessID, fromDate, toDate+" 23:59:59").
		Group("daily_boda_payments.driver_id, daily_boda_drivers.full_name, daily_boda_drivers.phone").
		Order("total_paid DESC").
		Limit(20).
		Scan(&topDrivers)

	type DailyRow struct {
		Date  string  `json:"date"`
		Total float64 `json:"total"`
	}
	var daily []DailyRow
	h.DB.Model(&models.DailyBodaPayment{}).
		Select("DATE(payment_date) as date, COALESCE(SUM(amount), 0) as total").
		Where("business_id = ? AND payment_date >= ? AND payment_date <= ?", businessID, fromDate, toDate+" 23:59:59").
		Group("DATE(payment_date)").
		Order("date ASC").
		Scan(&daily)

	// Fleet status breakdown
	type FleetStat struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var fleet []FleetStat
	h.DB.Model(&models.DailyBodaMotorcycle{}).
		Select("status, COUNT(*) as count").
		Where("business_id = ?", businessID).
		Group("status").
		Scan(&fleet)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"from_date":         fromDate,
			"to_date":           toDate,
			"total_collected":   totalCollected,
			"transaction_count": totalCount,
			"unpaid_balance":    totalBalance,
			"top_drivers":       topDrivers,
			"daily":             daily,
			"fleet":             fleet,
		},
	})
}
