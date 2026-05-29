package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// CreateLoanInput is the data needed to draft and persist a new Loan along with
// its full repayment schedule. Status starts at "pending" — call ApproveLoan
// and DisburseLoan to advance it.
type CreateLoanInput struct {
	BusinessID       uint
	BranchID         uint
	BorrowerID       uint
	MotorcycleID     *uint // nil = working capital loan, no asset
	LoanProductID    uint
	PrincipalAmount  float64
	InitialDeposit   float64
	Duration         int       // overrides product default if non-zero
	InterestRate     float64   // overrides product default if non-zero (annual %)
	InterestMethod   string    // overrides product default if non-empty
	RepaymentCycle   string    // overrides product default if non-empty
	FirstPaymentDate time.Time // due date of installment #1
	CreatedBy        uint
}

// CreateLoan persists a Loan and its RepaymentSchedule rows in one transaction.
// Pulls defaults from the LoanProduct, validates against its bounds, generates
// the schedule, and stamps a human-readable LoanNumber.
func CreateLoan(db *gorm.DB, in CreateLoanInput) (*models.Loan, error) {
	if in.PrincipalAmount <= 0 {
		return nil, errors.New("principal amount must be positive")
	}
	if in.InitialDeposit < 0 || in.InitialDeposit >= in.PrincipalAmount {
		return nil, errors.New("initial deposit must be >= 0 and less than principal")
	}
	if in.FirstPaymentDate.IsZero() {
		return nil, errors.New("first payment date is required")
	}

	var product models.LoanProduct
	if err := db.First(&product, in.LoanProductID).Error; err != nil {
		return nil, fmt.Errorf("loan product not found: %w", err)
	}
	if product.BusinessID != in.BusinessID {
		return nil, errors.New("loan product belongs to a different business")
	}
	if !product.IsActive {
		return nil, errors.New("loan product is inactive")
	}
	if in.PrincipalAmount < product.MinAmount || in.PrincipalAmount > product.MaxAmount {
		return nil, fmt.Errorf("principal %.0f outside product range [%.0f, %.0f]",
			in.PrincipalAmount, product.MinAmount, product.MaxAmount)
	}

	// Resolve terms: input wins, otherwise fall back to product defaults.
	duration := in.Duration
	if duration == 0 {
		duration = product.MinDuration
	}
	if duration < product.MinDuration || duration > product.MaxDuration {
		return nil, fmt.Errorf("duration %d outside product range [%d, %d]",
			duration, product.MinDuration, product.MaxDuration)
	}
	rate := in.InterestRate
	if rate == 0 {
		rate = product.InterestRate
	}
	method := in.InterestMethod
	if method == "" {
		method = product.InterestMethod
	}
	cycle := in.RepaymentCycle
	if cycle == "" {
		cycle = product.RepaymentCycle
	}

	disbursed := in.PrincipalAmount - in.InitialDeposit

	summary, err := GenerateSchedule(ScheduleParams{
		Principal:      disbursed,
		AnnualRate:     rate,
		Duration:       duration,
		InterestMethod: method,
		Cycle:          cycle,
		StartDate:      in.FirstPaymentDate,
	})
	if err != nil {
		return nil, fmt.Errorf("schedule: %w", err)
	}

	loan := &models.Loan{
		BusinessID:        in.BusinessID,
		BranchID:          in.BranchID,
		BorrowerID:        in.BorrowerID,
		MotorcycleID:      in.MotorcycleID,
		LoanProductID:     in.LoanProductID,
		PrincipalAmount:   in.PrincipalAmount,
		InitialDeposit:    in.InitialDeposit,
		DisbursedAmount:   disbursed,
		Duration:          duration,
		InterestMethod:    method,
		InterestRate:      rate,
		RepaymentCycle:    cycle,
		TotalInterest:     summary.TotalInterest,
		TotalRepayments:   duration,
		InstallmentAmount: summary.InstallmentAmount,
		TotalAmount:       summary.TotalAmount,
		BalanceRemaining:  summary.TotalAmount, // before any payment
		Status:            models.LoanStatusPending,
		FirstPaymentDate:  &in.FirstPaymentDate,
		NextPaymentDate:   &summary.Installments[0].DueDate,
		MaturityDate:      &summary.MaturityDate,
		CreatedBy:         in.CreatedBy,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(loan).Error; err != nil {
			return err
		}
		// Stamp human-readable loan number now that we have the auto ID.
		loan.LoanNumber = fmt.Sprintf("KM-LN-%06d", loan.ID)
		if err := tx.Model(loan).Update("loan_number", loan.LoanNumber).Error; err != nil {
			return err
		}

		rows := make([]models.RepaymentSchedule, 0, len(summary.Installments))
		for _, inst := range summary.Installments {
			rows = append(rows, models.RepaymentSchedule{
				LoanID:            loan.ID,
				InstallmentNumber: inst.Number,
				DueDate:           inst.DueDate,
				PrincipalAmount:   inst.PrincipalAmount,
				InterestAmount:    inst.InterestAmount,
				TotalAmount:       inst.TotalAmount,
				BalanceBefore:     inst.BalanceBefore,
				BalanceAfter:      inst.BalanceAfter,
			})
		}
		if err := tx.CreateInBatches(rows, 100).Error; err != nil {
			return err
		}

		// If a motorcycle is being financed, reserve it.
		if in.MotorcycleID != nil {
			if err := tx.Model(&models.Motorcycle{}).
				Where("id = ? AND business_id = ?", *in.MotorcycleID, in.BusinessID).
				Update("status", models.MotorcycleStatusReserved).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// ApproveLoan flips a pending loan to approved.
func ApproveLoan(db *gorm.DB, loanID, approverID uint) error {
	now := time.Now()
	return db.Model(&models.Loan{}).
		Where("id = ? AND status = ?", loanID, models.LoanStatusPending).
		Updates(map[string]interface{}{
			"status":      models.LoanStatusApproved,
			"approved_by": approverID,
			"updated_at":  now,
		}).Error
}

// DisburseLoan flips an approved loan to active and stamps the disbursement
// date. Marks the financed motorcycle as on_loan.
func DisburseLoan(db *gorm.DB, loanID, disburserID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var loan models.Loan
		if err := tx.First(&loan, loanID).Error; err != nil {
			return err
		}
		if loan.Status != models.LoanStatusApproved {
			return fmt.Errorf("loan must be approved before disbursement (current: %s)", loan.Status)
		}
		now := time.Now()
		if err := tx.Model(&loan).Updates(map[string]interface{}{
			"status":            models.LoanStatusActive,
			"disbursement_date": now,
			"disbursed_by":      disburserID,
		}).Error; err != nil {
			return err
		}
		if loan.MotorcycleID != nil {
			if err := tx.Model(&models.Motorcycle{}).
				Where("id = ?", *loan.MotorcycleID).
				Update("status", models.MotorcycleStatusOnLoan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ApplyRepayment records an APPROVED Repayment against a Loan: marks the matched
// schedule rows paid, updates the Loan.BalanceRemaining and NextPaymentDate, and
// completes the Loan if balance hits zero. Caller is responsible for creating
// the Repayment record itself (status=approved) and passing its amount/scheduleID.
//
// This is idempotent on the schedule row: if it's already fully paid, the
// remaining amount cascades to the next unpaid row.
func ApplyRepayment(tx *gorm.DB, loanID uint, amount float64, preferredScheduleID *uint) error {
	if amount <= 0 {
		return errors.New("repayment amount must be positive")
	}

	var loan models.Loan
	if err := tx.First(&loan, loanID).Error; err != nil {
		return err
	}

	remaining := amount

	// First, apply to the preferred schedule row if specified.
	if preferredScheduleID != nil {
		applied, err := applyToSchedule(tx, *preferredScheduleID, remaining)
		if err != nil {
			return err
		}
		remaining -= applied
	}

	// Cascade: any leftover applies to oldest unpaid installments.
	if remaining > 0.001 {
		var rows []models.RepaymentSchedule
		if err := tx.Where("loan_id = ? AND is_paid = ?", loanID, false).
			Order("installment_number ASC").Find(&rows).Error; err != nil {
			return err
		}
		for _, row := range rows {
			if remaining <= 0.001 {
				break
			}
			applied, err := applyToSchedule(tx, row.ID, remaining)
			if err != nil {
				return err
			}
			remaining -= applied
		}
	}

	// Update loan balance.
	newBalance := round2(loan.BalanceRemaining - amount)
	if newBalance < 0 {
		newBalance = 0
	}
	updates := map[string]interface{}{
		"balance_remaining": newBalance,
	}
	if newBalance <= 0.001 {
		now := time.Now()
		updates["status"] = models.LoanStatusCompleted
		updates["completed_at"] = now
		updates["next_payment_date"] = nil

		// Mark the financed motorcycle as sold (loan paid off).
		if loan.MotorcycleID != nil {
			if err := tx.Model(&models.Motorcycle{}).
				Where("id = ?", *loan.MotorcycleID).
				Update("status", models.MotorcycleStatusSold).Error; err != nil {
				return err
			}
		}
	} else {
		// Update next payment date to the next unpaid installment.
		var next models.RepaymentSchedule
		err := tx.Where("loan_id = ? AND is_paid = ?", loanID, false).
			Order("installment_number ASC").First(&next).Error
		if err == nil {
			updates["next_payment_date"] = next.DueDate
		}
	}
	return tx.Model(&loan).Updates(updates).Error
}

// applyToSchedule applies up to `amount` to one RepaymentSchedule row.
// Returns the amount actually consumed (capped at the row's outstanding total).
func applyToSchedule(tx *gorm.DB, scheduleID uint, amount float64) (float64, error) {
	var row models.RepaymentSchedule
	if err := tx.First(&row, scheduleID).Error; err != nil {
		return 0, err
	}
	if row.IsPaid {
		return 0, nil
	}
	outstanding := round2(row.TotalAmount - row.PaidAmount)
	if outstanding <= 0 {
		// Already settled but flag wasn't flipped — flip it now.
		return 0, tx.Model(&row).Update("is_paid", true).Error
	}

	apply := amount
	if apply > outstanding {
		apply = outstanding
	}

	updates := map[string]interface{}{
		"paid_amount": round2(row.PaidAmount + apply),
	}
	if round2(row.PaidAmount+apply) >= row.TotalAmount-0.001 {
		now := time.Now()
		updates["is_paid"] = true
		updates["paid_at"] = now
		updates["is_overdue"] = false
		updates["days_past_due"] = 0
	}
	if err := tx.Model(&row).Updates(updates).Error; err != nil {
		return 0, err
	}
	return apply, nil
}

// MarkOverdueSchedules flips any unpaid schedule row whose due date has passed
// to is_overdue=true and updates days_past_due. Designed to run nightly.
// Returns the number of rows touched.
func MarkOverdueSchedules(db *gorm.DB, now time.Time) (int64, error) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	res := db.Exec(`
		UPDATE repayment_schedules
		SET is_overdue = TRUE,
		    days_past_due = GREATEST(0, EXTRACT(DAY FROM (? - due_date))::int),
		    updated_at = ?
		WHERE is_paid = FALSE AND due_date < ?
	`, today, now, today)
	return res.RowsAffected, res.Error
}
