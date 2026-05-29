package services

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gritdemo/internal/models"
)

var refStart = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// almostEqual returns true if a and b are within 1 UGX (rounding tolerance).
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1.0
}

func TestFlat_MonthlySimple(t *testing.T) {
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      1_000_000,
		AnnualRate:     24,
		Duration:       12,
		InterestMethod: models.InterestMethodFlat,
		Cycle:          models.RepaymentCycleMonthly,
		StartDate:      refStart,
	})
	require.NoError(t, err)
	require.Len(t, s.Installments, 12)

	// Flat 24% annual on 1M for 12 months = 240K total interest.
	assert.True(t, almostEqual(s.TotalInterest, 240_000), "total interest = %f", s.TotalInterest)
	assert.True(t, almostEqual(s.TotalAmount, 1_240_000), "total amount = %f", s.TotalAmount)

	// Every installment (except possibly last) has principal = 1M/12, interest = 240K/12.
	first := s.Installments[0]
	assert.True(t, almostEqual(first.PrincipalAmount, 83_333.33))
	assert.True(t, almostEqual(first.InterestAmount, 20_000))
	assert.True(t, almostEqual(first.TotalAmount, 103_333.33))

	// Schedule must end at zero balance, and dates step by 1 month each.
	last := s.Installments[len(s.Installments)-1]
	assert.True(t, almostEqual(last.BalanceAfter, 0), "final balance = %f", last.BalanceAfter)
	assert.Equal(t, refStart, s.Installments[0].DueDate)
	assert.Equal(t, refStart.AddDate(0, 11, 0), last.DueDate)
}

func TestFlat_Weekly(t *testing.T) {
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      520_000,
		AnnualRate:     52,
		Duration:       52,
		InterestMethod: models.InterestMethodFlat,
		Cycle:          models.RepaymentCycleWeekly,
		StartDate:      refStart,
	})
	require.NoError(t, err)
	require.Len(t, s.Installments, 52)

	// 52 weeks = 1 year, 52% rate => 270400 total interest.
	assert.True(t, almostEqual(s.TotalInterest, 520_000*0.52), "total interest = %f", s.TotalInterest)

	// Sum of all principal portions equals the original principal.
	var sumP float64
	for _, inst := range s.Installments {
		sumP += inst.PrincipalAmount
	}
	assert.True(t, almostEqual(sumP, 520_000), "sum of principal = %f", sumP)

	// Weekly steps.
	assert.Equal(t, refStart.AddDate(0, 0, 7), s.Installments[1].DueDate)
}

func TestFlat_Biweekly(t *testing.T) {
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      300_000,
		AnnualRate:     12,
		Duration:       6,
		InterestMethod: models.InterestMethodFlat,
		Cycle:          models.RepaymentCycleBiweekly,
		StartDate:      refStart,
	})
	require.NoError(t, err)

	// 6 biweekly = ~12 weeks ≈ 12/26 of a year, 12% annual.
	expected := 300_000 * 0.12 * (6.0 / 26.0)
	assert.True(t, almostEqual(s.TotalInterest, expected), "got %f, want %f", s.TotalInterest, expected)
	assert.Equal(t, refStart.AddDate(0, 0, 14), s.Installments[1].DueDate)
}

func TestReducingBalance_Monthly(t *testing.T) {
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      1_000_000,
		AnnualRate:     24,
		Duration:       12,
		InterestMethod: models.InterestMethodReducingBalance,
		Cycle:          models.RepaymentCycleMonthly,
		StartDate:      refStart,
	})
	require.NoError(t, err)
	require.Len(t, s.Installments, 12)

	// Constant installment via amortization: 1,000,000 * 0.02 * 1.02^12 / (1.02^12 - 1)
	// 1.02^12 ≈ 1.26824. expected ≈ 94,560.
	expectedInstallment := 94_560.0
	assert.InDelta(t, expectedInstallment, s.InstallmentAmount, 50)

	// Reducing balance should produce LESS total interest than flat at same rate/term.
	flat, _ := GenerateSchedule(ScheduleParams{
		Principal:      1_000_000,
		AnnualRate:     24,
		Duration:       12,
		InterestMethod: models.InterestMethodFlat,
		Cycle:          models.RepaymentCycleMonthly,
		StartDate:      refStart,
	})
	assert.Less(t, s.TotalInterest, flat.TotalInterest, "reducing balance should be cheaper than flat")

	// Final balance is zero.
	last := s.Installments[len(s.Installments)-1]
	assert.True(t, almostEqual(last.BalanceAfter, 0), "final balance = %f", last.BalanceAfter)

	// First installment's interest portion = balance * monthly rate = 1M * 0.02 = 20K.
	assert.InDelta(t, 20_000, s.Installments[0].InterestAmount, 1)

	// Interest portion shrinks each month, principal portion grows.
	assert.Greater(t, s.Installments[0].InterestAmount, s.Installments[11].InterestAmount)
	assert.Less(t, s.Installments[0].PrincipalAmount, s.Installments[11].PrincipalAmount)
}

func TestZeroInterest(t *testing.T) {
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      120_000,
		AnnualRate:     0,
		Duration:       12,
		InterestMethod: models.InterestMethodFlat,
		Cycle:          models.RepaymentCycleMonthly,
		StartDate:      refStart,
	})
	require.NoError(t, err)
	assert.True(t, almostEqual(s.TotalInterest, 0))
	for _, inst := range s.Installments {
		assert.True(t, almostEqual(inst.InterestAmount, 0))
		assert.True(t, almostEqual(inst.PrincipalAmount, 10_000))
	}
}

func TestZeroInterest_ReducingBalance(t *testing.T) {
	// Zero interest with reducing balance must avoid divide-by-zero in the amortization formula.
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      120_000,
		AnnualRate:     0,
		Duration:       12,
		InterestMethod: models.InterestMethodReducingBalance,
		Cycle:          models.RepaymentCycleMonthly,
		StartDate:      refStart,
	})
	require.NoError(t, err)
	assert.True(t, almostEqual(s.TotalInterest, 0))
	assert.True(t, almostEqual(s.Installments[0].PrincipalAmount, 10_000))
}

func TestSingleInstallment(t *testing.T) {
	s, err := GenerateSchedule(ScheduleParams{
		Principal:      500_000,
		AnnualRate:     24,
		Duration:       1,
		InterestMethod: models.InterestMethodFlat,
		Cycle:          models.RepaymentCycleMonthly,
		StartDate:      refStart,
	})
	require.NoError(t, err)
	require.Len(t, s.Installments, 1)
	assert.True(t, almostEqual(s.Installments[0].PrincipalAmount, 500_000))
	assert.True(t, almostEqual(s.Installments[0].BalanceAfter, 0))
}

func TestInvalidInputs(t *testing.T) {
	cases := []ScheduleParams{
		{Principal: 0, AnnualRate: 10, Duration: 12, InterestMethod: models.InterestMethodFlat, Cycle: models.RepaymentCycleMonthly},
		{Principal: 1000, AnnualRate: 10, Duration: 0, InterestMethod: models.InterestMethodFlat, Cycle: models.RepaymentCycleMonthly},
		{Principal: 1000, AnnualRate: -1, Duration: 12, InterestMethod: models.InterestMethodFlat, Cycle: models.RepaymentCycleMonthly},
		{Principal: 1000, AnnualRate: 10, Duration: 12, InterestMethod: "bogus", Cycle: models.RepaymentCycleMonthly},
		{Principal: 1000, AnnualRate: 10, Duration: 12, InterestMethod: models.InterestMethodFlat, Cycle: "bogus"},
	}
	for i, c := range cases {
		_, err := GenerateSchedule(c)
		assert.Error(t, err, "case %d should have failed", i)
	}
}

func TestScheduleSumsExactly(t *testing.T) {
	// Both methods: sum of principal portions must equal original principal,
	// sum of total amounts must equal principal + total interest.
	for _, method := range []string{models.InterestMethodFlat, models.InterestMethodReducingBalance} {
		t.Run(method, func(t *testing.T) {
			s, err := GenerateSchedule(ScheduleParams{
				Principal:      1_337_000,
				AnnualRate:     17.5,
				Duration:       9,
				InterestMethod: method,
				Cycle:          models.RepaymentCycleMonthly,
				StartDate:      refStart,
			})
			require.NoError(t, err)

			var sumP, sumT float64
			for _, inst := range s.Installments {
				sumP += inst.PrincipalAmount
				sumT += inst.TotalAmount
			}
			assert.True(t, almostEqual(sumP, 1_337_000), "sum of principal = %f", sumP)
			assert.True(t, almostEqual(sumT, s.TotalAmount), "sum of totals = %f, want %f", sumT, s.TotalAmount)
		})
	}
}

func TestBalanceMonotonicallyDecreases(t *testing.T) {
	for _, method := range []string{models.InterestMethodFlat, models.InterestMethodReducingBalance} {
		t.Run(method, func(t *testing.T) {
			s, err := GenerateSchedule(ScheduleParams{
				Principal:      800_000,
				AnnualRate:     30,
				Duration:       24,
				InterestMethod: method,
				Cycle:          models.RepaymentCycleWeekly,
				StartDate:      refStart,
			})
			require.NoError(t, err)
			for i := 1; i < len(s.Installments); i++ {
				assert.LessOrEqual(t, s.Installments[i].BalanceAfter, s.Installments[i-1].BalanceAfter)
			}
		})
	}
}
