package services

import (
	"fmt"
	"math"
	"time"

	"gritdemo/internal/models"
)

// Installment is one row of a repayment schedule, computed before the loan is
// persisted. The loan service maps each Installment to a RepaymentSchedule row.
type Installment struct {
	Number          int
	DueDate         time.Time
	PrincipalAmount float64
	InterestAmount  float64
	TotalAmount     float64
	BalanceBefore   float64
	BalanceAfter    float64
}

// ScheduleParams describes a loan's terms enough to generate its schedule.
// Principal is the DISBURSED amount (already net of any initial deposit).
type ScheduleParams struct {
	Principal      float64
	AnnualRate     float64 // % e.g. 24.0 means 24%
	Duration       int     // number of installments
	InterestMethod string  // flat | reducing_balance
	Cycle          string  // weekly | biweekly | monthly
	StartDate      time.Time
}

// ScheduleSummary is the roll-up an installment schedule produces — what we
// store back on the Loan record (TotalInterest, InstallmentAmount, etc).
type ScheduleSummary struct {
	Installments      []Installment
	TotalInterest     float64
	TotalAmount       float64 // Principal + TotalInterest
	InstallmentAmount float64 // representative — first installment's TotalAmount (constant for reducing balance, equal for flat)
	MaturityDate      time.Time
}

// GenerateSchedule produces the installment schedule + summary for a loan.
func GenerateSchedule(p ScheduleParams) (*ScheduleSummary, error) {
	if p.Principal <= 0 {
		return nil, fmt.Errorf("principal must be positive")
	}
	if p.Duration <= 0 {
		return nil, fmt.Errorf("duration must be positive")
	}
	if p.AnnualRate < 0 {
		return nil, fmt.Errorf("annual rate cannot be negative")
	}

	cyclesPerYear, advance, err := cycleParams(p.Cycle)
	if err != nil {
		return nil, err
	}

	var installments []Installment
	switch p.InterestMethod {
	case models.InterestMethodFlat:
		installments, err = flatSchedule(p, cyclesPerYear, advance)
	case models.InterestMethodReducingBalance:
		installments, err = reducingBalanceSchedule(p, cyclesPerYear, advance)
	default:
		return nil, fmt.Errorf("unknown interest method: %q", p.InterestMethod)
	}
	if err != nil {
		return nil, err
	}

	var totalInterest float64
	for _, i := range installments {
		totalInterest += i.InterestAmount
	}
	totalInterest = round2(totalInterest)

	return &ScheduleSummary{
		Installments:      installments,
		TotalInterest:     totalInterest,
		TotalAmount:       round2(p.Principal + totalInterest),
		InstallmentAmount: installments[0].TotalAmount,
		MaturityDate:      installments[len(installments)-1].DueDate,
	}, nil
}

func cycleParams(cycle string) (cyclesPerYear int, advance func(time.Time) time.Time, err error) {
	switch cycle {
	case models.RepaymentCycleWeekly:
		return 52, func(t time.Time) time.Time { return t.AddDate(0, 0, 7) }, nil
	case models.RepaymentCycleBiweekly:
		return 26, func(t time.Time) time.Time { return t.AddDate(0, 0, 14) }, nil
	case models.RepaymentCycleMonthly:
		return 12, func(t time.Time) time.Time { return t.AddDate(0, 1, 0) }, nil
	default:
		return 0, nil, fmt.Errorf("unknown repayment cycle: %q", cycle)
	}
}

// flatSchedule: total interest is computed once across the loan term, then split
// evenly across all installments. Each installment has equal principal (P/n) and
// equal interest (totalInterest/n). The final installment absorbs rounding
// drift so the schedule sums to exactly principal + totalInterest.
func flatSchedule(p ScheduleParams, cyclesPerYear int, advance func(time.Time) time.Time) ([]Installment, error) {
	durationInYears := float64(p.Duration) / float64(cyclesPerYear)
	totalInterest := round2(p.Principal * (p.AnnualRate / 100.0) * durationInYears)

	principalPer := round2(p.Principal / float64(p.Duration))
	interestPer := round2(totalInterest / float64(p.Duration))

	out := make([]Installment, 0, p.Duration)
	balance := p.Principal
	due := p.StartDate

	for i := 1; i <= p.Duration; i++ {
		principal := principalPer
		interest := interestPer
		if i == p.Duration {
			// Last row: pay off the remaining balance exactly + reconcile interest rounding.
			principal = round2(balance)
			interest = round2(totalInterest - interestPer*float64(p.Duration-1))
		}
		balanceAfter := round2(balance - principal)

		out = append(out, Installment{
			Number:          i,
			DueDate:         due,
			PrincipalAmount: principal,
			InterestAmount:  interest,
			TotalAmount:     round2(principal + interest),
			BalanceBefore:   round2(balance),
			BalanceAfter:    balanceAfter,
		})

		balance = balanceAfter
		due = advance(due)
	}
	return out, nil
}

// reducingBalanceSchedule: equal-installment amortization. Periodic interest is
// applied to the outstanding balance; the principal portion of each installment
// grows as the balance shrinks. Closed-form formula keeps installment constant.
func reducingBalanceSchedule(p ScheduleParams, cyclesPerYear int, advance func(time.Time) time.Time) ([]Installment, error) {
	r := (p.AnnualRate / 100.0) / float64(cyclesPerYear)

	var installment float64
	if r == 0 {
		installment = p.Principal / float64(p.Duration)
	} else {
		// P = Pv * r * (1+r)^n / ((1+r)^n - 1)
		f := math.Pow(1+r, float64(p.Duration))
		installment = (p.Principal * r * f) / (f - 1)
	}
	installment = round2(installment)

	out := make([]Installment, 0, p.Duration)
	balance := p.Principal
	due := p.StartDate

	for i := 1; i <= p.Duration; i++ {
		interest := round2(balance * r)
		principal := round2(installment - interest)
		total := installment
		if i == p.Duration {
			principal = round2(balance)
			total = round2(principal + interest)
		}
		balanceAfter := round2(balance - principal)

		out = append(out, Installment{
			Number:          i,
			DueDate:         due,
			PrincipalAmount: principal,
			InterestAmount:  interest,
			TotalAmount:     total,
			BalanceBefore:   round2(balance),
			BalanceAfter:    balanceAfter,
		})

		balance = balanceAfter
		due = advance(due)
	}
	return out, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
