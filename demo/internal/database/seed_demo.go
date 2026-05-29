package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// SeedDemo creates the rich motorcycle-dealer demo cohort on top of the
// canonical Seed() output (admin + Grit Motors + Main Branch). Called
// from main.go when DEMO_MODE=true and from the nightly demo-reset cron
// after WipeMutableData has truncated the activity tables.
//
// Each step is idempotent: it checks for existence by stable identifier
// (email / phone / name) and skips the create when the row is already
// there. The same source of truth is used for the daily reset, so the
// cohort is restored deterministically every midnight.
//
// What this lays down:
//   - 4 staff users (manager / cashier / stock-clerk / loan-officer)
//     with predictable credentials.
//   - 6 spare-parts categories + 12 products + opening stock.
//   - 3 loan products (12/18/24-month) with realistic rates.
//   - 8 motorcycles for inventory + 6 borrowers + 5 active loans.
//   - Repayment schedules generated; ~50% of due installments paid.
//   - 3 cash motorcycle sales + 8 historical spare-part POS sales.
//   - 3 daily-boda riders + 3 motorcycles + a week of payment history.
func SeedDemo(db *gorm.DB) error {
	// Resolve the canonical business + main branch + admin user. These
	// always exist (Seed() runs before us).
	var business models.Business
	if err := db.Where("name = ?", gritMotorsBusinessName).First(&business).Error; err != nil {
		return fmt.Errorf("seed_demo: no business: %w", err)
	}
	var mainBranch models.Branch
	if err := db.Where("business_id = ? AND is_default = ?", business.ID, true).First(&mainBranch).Error; err != nil {
		return fmt.Errorf("seed_demo: no default branch: %w", err)
	}
	var admin models.User
	if err := db.Where("email = ?", gritDemoAdminEmail).First(&admin).Error; err != nil {
		return fmt.Errorf("seed_demo: no admin: %w", err)
	}

	// Deterministic RNG so the demo cohort looks the same across resets.
	rng := rand.New(rand.NewSource(20260529))

	staff, err := seedStaff(db, business, mainBranch)
	if err != nil {
		return fmt.Errorf("staff: %w", err)
	}
	categories, err := seedCategories(db, business)
	if err != nil {
		return fmt.Errorf("categories: %w", err)
	}
	products, err := seedProducts(db, business, categories, admin.ID)
	if err != nil {
		return fmt.Errorf("products: %w", err)
	}
	if err := seedStocks(db, products, mainBranch); err != nil {
		return fmt.Errorf("stocks: %w", err)
	}
	motorcycles, err := seedMotorcycles(db, business, mainBranch, admin.ID)
	if err != nil {
		return fmt.Errorf("motorcycles: %w", err)
	}
	loanProducts, err := seedLoanProducts(db, business)
	if err != nil {
		return fmt.Errorf("loan products: %w", err)
	}
	borrowers, err := seedBorrowers(db, business, mainBranch)
	if err != nil {
		return fmt.Errorf("borrowers: %w", err)
	}
	if err := seedLoans(db, business, mainBranch, borrowers, motorcycles, loanProducts, rng); err != nil {
		return fmt.Errorf("loans: %w", err)
	}
	if err := seedCashSales(db, business, mainBranch, motorcycles, staff); err != nil {
		return fmt.Errorf("cash sales: %w", err)
	}
	if err := seedSpareSales(db, business, mainBranch, products, staff, rng); err != nil {
		return fmt.Errorf("spare sales: %w", err)
	}
	if err := seedDailyBoda(db, business, mainBranch, staff, rng); err != nil {
		return fmt.Errorf("daily boda: %w", err)
	}

	log.Println("seed_demo: cohort ready")
	return nil
}

// ─── staff ────────────────────────────────────────────────────────────
type staffSet struct {
	manager     models.User
	cashier     models.User
	stockClerk  models.User
	loanOfficer models.User
}

func seedStaff(db *gorm.DB, biz models.Business, branch models.Branch) (staffSet, error) {
	rows := []struct {
		name, email, password, role string
	}{
		{"Maya Manager", "maya@grit.demo", "password123", models.RoleManager},
		{"Carlos Cashier", "carlos@grit.demo", "password123", models.RoleCashier},
		{"Stella Stock", "stella@grit.demo", "password123", models.RoleStockClerk},
		{"Lawrence Loans", "lawrence@grit.demo", "password123", models.RoleLoanOfficer},
	}
	out := make([]models.User, 0, len(rows))
	for _, r := range rows {
		u, err := upsertUser(db, r.name, r.email, r.password)
		if err != nil {
			return staffSet{}, err
		}
		if err := upsertRole(db, u.ID, biz.ID, branch.ID, r.role); err != nil {
			return staffSet{}, err
		}
		out = append(out, u)
	}
	return staffSet{manager: out[0], cashier: out[1], stockClerk: out[2], loanOfficer: out[3]}, nil
}

func upsertUser(db *gorm.DB, name, email, password string) (models.User, error) {
	var u models.User
	if err := db.Where("email = ?", email).First(&u).Error; err == nil {
		return u, nil
	}
	u = models.User{Name: name, Email: email}
	if err := u.SetPassword(password); err != nil {
		return u, err
	}
	if err := db.Create(&u).Error; err != nil {
		return u, err
	}
	log.Printf("seed_demo: user %s / %s", email, password)
	return u, nil
}

func upsertRole(db *gorm.DB, userID, bizID, branchID uint, role string) error {
	var r models.UserBusinessRole
	err := db.Where("user_id = ? AND business_id = ?", userID, bizID).First(&r).Error
	if err == nil {
		return nil
	}
	bid := branchID
	return db.Create(&models.UserBusinessRole{
		UserID: userID, BusinessID: bizID, BranchID: &bid, Role: role, WorkspaceAccess: "both",
	}).Error
}

// ─── categories + products + stocks ───────────────────────────────────
func seedCategories(db *gorm.DB, biz models.Business) ([]models.Category, error) {
	names := []string{"Engine", "Brakes", "Tyres & Tubes", "Lights & Electrical", "Lubricants", "Body & Frame"}
	out := make([]models.Category, 0, len(names))
	for _, name := range names {
		var c models.Category
		if err := db.Where("business_id = ? AND name = ?", biz.ID, name).First(&c).Error; err == nil {
			out = append(out, c)
			continue
		}
		c = models.Category{BusinessID: biz.ID, Name: name}
		if err := db.Create(&c).Error; err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

type productSpec struct {
	categoryIdx int
	title       string
	sellingUGX  float64
	costUGX     float64
	threshold   int
}

func seedProducts(db *gorm.DB, biz models.Business, cats []models.Category, createdBy uint) ([]models.Product, error) {
	specs := []productSpec{
		{0, "Air filter (Boxer)", 12000, 7000, 8},
		{0, "Spark plug NGK", 8000, 4500, 12},
		{0, "Engine oil 1L (10W-40)", 22000, 14000, 10},
		{1, "Front brake pads", 25000, 15000, 6},
		{1, "Rear brake shoes", 18000, 11000, 6},
		{2, "Tube (17 inch)", 14000, 8500, 10},
		{2, "Tyre 90/90-17", 95000, 65000, 4},
		{3, "Headlight bulb 35W", 9000, 5500, 10},
		{3, "Battery 5Ah", 110000, 78000, 3},
		{4, "Chain lube spray", 16000, 9500, 8},
		{4, "Engine flush 250ml", 12000, 7000, 6},
		{5, "Rear-view mirror (pair)", 28000, 17000, 5},
	}
	out := make([]models.Product, 0, len(specs))
	for _, s := range specs {
		var p models.Product
		if err := db.Where("business_id = ? AND title = ?", biz.ID, s.title).First(&p).Error; err == nil {
			out = append(out, p)
			continue
		}
		p = models.Product{
			BusinessID:        biz.ID,
			CategoryID:        cats[s.categoryIdx].ID,
			Title:             s.title,
			Description:       s.title + " — replacement part, sold per unit.",
			SellingPrice:      s.sellingUGX,
			CostPrice:         s.costUGX,
			LowStockThreshold: s.threshold,
			CreatedBy:         createdBy,
		}
		if err := db.Create(&p).Error; err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func seedStocks(db *gorm.DB, products []models.Product, branch models.Branch) error {
	quantities := []int{40, 60, 30, 20, 25, 35, 15, 28, 12, 22, 18, 14}
	for i, p := range products {
		var st models.Stock
		if err := db.Where("product_id = ? AND branch_id = ?", p.ID, branch.ID).First(&st).Error; err == nil {
			continue
		}
		if err := db.Create(&models.Stock{
			ProductID: p.ID, BranchID: branch.ID, Quantity: quantities[i%len(quantities)],
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── motorcycles ──────────────────────────────────────────────────────
func seedMotorcycles(db *gorm.DB, biz models.Business, branch models.Branch, createdBy uint) ([]models.Motorcycle, error) {
	rows := []struct {
		name, plate, chassis, engine, color string
		year                                int
		cost, price                         float64
	}{
		{"Bajaj Boxer 100", "UBE 401K", "MD2A18AY8KCH00111", "JKZWBN02411", "Red", 2024, 4_200_000, 5_400_000},
		{"Bajaj Boxer 100", "UBE 412K", "MD2A18AY8KCH00112", "JKZWBN02412", "Black", 2024, 4_200_000, 5_400_000},
		{"TVS HLX 125", "UBE 718L", "MD625GF52L1A00211", "OE7RA20211", "Blue", 2024, 4_950_000, 6_300_000},
		{"Honda CG 125", "UBF 220A", "JH2JC32A4LK000321", "JC32E000321", "Red", 2023, 5_400_000, 7_200_000},
		{"Bajaj Discover 100", "UBE 909M", "MD2A21BY1MCH00401", "JKZWHN00401", "Silver", 2025, 4_800_000, 6_100_000},
		{"TVS HLX 150", "UBE 880F", "MD625GF52L1A00511", "OE7RA20511", "Yellow", 2024, 5_400_000, 6_950_000},
		{"Bajaj Boxer 150X", "UBF 333J", "MD2A18BY8KCH00601", "JKZWBN02601", "Green", 2025, 5_100_000, 6_400_000},
		{"Honda Ace CB 125", "UBF 144P", "JH2JC32A4LK000711", "JC32E000711", "Black", 2024, 5_800_000, 7_500_000},
	}
	out := make([]models.Motorcycle, 0, len(rows))
	for _, r := range rows {
		var m models.Motorcycle
		if err := db.Where("business_id = ? AND number_plate = ?", biz.ID, r.plate).First(&m).Error; err == nil {
			out = append(out, m)
			continue
		}
		m = models.Motorcycle{
			BusinessID:   biz.ID,
			BranchID:     branch.ID,
			Name:         r.name,
			NumberPlate:  r.plate,
			ChassisNo:    r.chassis,
			EngineNo:     r.engine,
			Color:        r.color,
			YearOfMake:   r.year,
			CostPrice:    r.cost,
			SellingPrice: r.price,
			Status:       "available",
			CreatedBy:    createdBy,
		}
		if err := db.Create(&m).Error; err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

// ─── loan products ───────────────────────────────────────────────────
func seedLoanProducts(db *gorm.DB, biz models.Business) ([]models.LoanProduct, error) {
	rows := []models.LoanProduct{
		{
			BusinessID:     biz.ID,
			Name:           "Boda Starter — 12 months",
			Description:    "Entry-level boda financing. 12 monthly installments, reducing balance.",
			MinAmount:      3_500_000, MaxAmount: 6_500_000, MinDuration: 12, MaxDuration: 12,
			InterestMethod: "reducing_balance", InterestRate: 24.0, RepaymentCycle: "monthly",
			GracePeriodDays: 5, IsActive: true,
		},
		{
			BusinessID:     biz.ID,
			Name:           "Standard Boda — 18 months",
			Description:    "Standard 18-month boda financing with reducing balance.",
			MinAmount:      4_000_000, MaxAmount: 7_500_000, MinDuration: 18, MaxDuration: 18,
			InterestMethod: "reducing_balance", InterestRate: 22.0, RepaymentCycle: "monthly",
			GracePeriodDays: 7, IsActive: true,
		},
		{
			BusinessID:     biz.ID,
			Name:           "Long Boda — 24 months",
			Description:    "24-month boda financing. Lower instalments, longer total cost.",
			MinAmount:      4_500_000, MaxAmount: 8_500_000, MinDuration: 24, MaxDuration: 24,
			InterestMethod: "reducing_balance", InterestRate: 21.0, RepaymentCycle: "monthly",
			GracePeriodDays: 10, IsActive: true,
		},
	}
	out := make([]models.LoanProduct, 0, len(rows))
	for _, lp := range rows {
		var ex models.LoanProduct
		if err := db.Where("business_id = ? AND name = ?", biz.ID, lp.Name).First(&ex).Error; err == nil {
			out = append(out, ex)
			continue
		}
		if err := db.Create(&lp).Error; err != nil {
			return nil, err
		}
		out = append(out, lp)
	}
	return out, nil
}

// ─── borrowers ────────────────────────────────────────────────────────
func seedBorrowers(db *gorm.DB, biz models.Business, branch models.Branch) ([]models.Borrower, error) {
	rows := []models.Borrower{
		{FirstName: "Joseph", LastName: "Kato", Phone: "+256700111201", Email: "joseph.kato@grit.demo", Gender: "M", NationalID: "CM89000001JK", Address: "Nakawa, Kampala", EmploymentStatus: "self_employed", Occupation: "Boda rider", MonthlyIncome: 1_400_000, NextOfKinName: "Sarah Kato", NextOfKinPhone: "+256700111202", NextOfKinRelation: "Spouse", CreditScore: 720},
		{FirstName: "Esther", LastName: "Nakato", Phone: "+256700111203", Email: "esther.nakato@grit.demo", Gender: "F", NationalID: "CF91000002EN", Address: "Mukono, Central", EmploymentStatus: "employed", Occupation: "Shop attendant", Employer: "Mukono Supermart", MonthlyIncome: 950_000, NextOfKinName: "Daniel Nakato", NextOfKinPhone: "+256700111204", NextOfKinRelation: "Brother", CreditScore: 680},
		{FirstName: "Patrick", LastName: "Mugisha", Phone: "+256700111205", Email: "patrick.mugisha@grit.demo", Gender: "M", NationalID: "CM87000003PM", Address: "Ntinda, Kampala", EmploymentStatus: "self_employed", Occupation: "Boda rider", MonthlyIncome: 1_600_000, NextOfKinName: "Joy Mugisha", NextOfKinPhone: "+256700111206", NextOfKinRelation: "Spouse", CreditScore: 760},
		{FirstName: "Sandra", LastName: "Atim", Phone: "+256700111207", Email: "sandra.atim@grit.demo", Gender: "F", NationalID: "CF93000004SA", Address: "Lira, Northern", EmploymentStatus: "employed", Occupation: "Teacher", Employer: "St Mary College", MonthlyIncome: 1_100_000, NextOfKinName: "Peter Atim", NextOfKinPhone: "+256700111208", NextOfKinRelation: "Father", CreditScore: 700},
		{FirstName: "Brian", LastName: "Okello", Phone: "+256700111209", Email: "brian.okello@grit.demo", Gender: "M", NationalID: "CM85000005BO", Address: "Jinja Town", EmploymentStatus: "self_employed", Occupation: "Boda rider", MonthlyIncome: 1_500_000, NextOfKinName: "Grace Okello", NextOfKinPhone: "+256700111210", NextOfKinRelation: "Sister", CreditScore: 690},
		{FirstName: "Faith", LastName: "Nansubuga", Phone: "+256700111211", Email: "faith.nansubuga@grit.demo", Gender: "F", NationalID: "CF92000006FN", Address: "Entebbe Road", EmploymentStatus: "employed", Occupation: "Nurse", Employer: "Entebbe Hospital", MonthlyIncome: 1_250_000, NextOfKinName: "Eric Nansubuga", NextOfKinPhone: "+256700111212", NextOfKinRelation: "Husband", CreditScore: 740},
	}
	out := make([]models.Borrower, 0, len(rows))
	for _, b := range rows {
		var ex models.Borrower
		if err := db.Where("business_id = ? AND phone = ?", biz.ID, b.Phone).First(&ex).Error; err == nil {
			out = append(out, ex)
			continue
		}
		b.BusinessID = biz.ID
		b.BranchID = branch.ID
		if err := db.Create(&b).Error; err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

// ─── loans + schedules + repayments ───────────────────────────────────
func seedLoans(db *gorm.DB, biz models.Business, branch models.Branch, borrowers []models.Borrower, motorcycles []models.Motorcycle, products []models.LoanProduct, rng *rand.Rand) error {
	if len(borrowers) < 5 || len(motorcycles) < 5 || len(products) < 3 {
		return nil
	}
	// 5 active loans — span new-to-aging so the dashboard shows life.
	startOffsets := []int{30, 75, 150, 210, 280} // days ago
	productIdx := []int{0, 1, 2, 1, 0}           // mix product mix
	for i := 0; i < 5; i++ {
		borrower := borrowers[i]
		moto := motorcycles[i]
		prod := products[productIdx[i]]
		loanNumber := fmt.Sprintf("GM-LN-%06d", 100000+i)

		// Idempotency by loan number.
		var ex models.Loan
		if err := db.Where("loan_number = ?", loanNumber).First(&ex).Error; err == nil {
			continue
		}

		principal := moto.SellingPrice
		deposit := principal * 0.20 // 20% down
		disbursed := principal - deposit
		months := prod.MinDuration

		// Reducing-balance: total interest ≈ disbursed × rate × (months+1) / 24
		totalInterest := disbursed * prod.InterestRate / 100 * float64(months+1) / 24
		total := disbursed + totalInterest
		installment := total / float64(months)

		startDate := time.Now().AddDate(0, 0, -startOffsets[i])
		l := models.Loan{
			BusinessID: biz.ID, BranchID: branch.ID,
			BorrowerID: borrower.ID, MotorcycleID: &moto.ID, LoanProductID: prod.ID,
			LoanNumber:        loanNumber,
			PrincipalAmount:   principal,
			InitialDeposit:    deposit,
			DisbursedAmount:   disbursed,
			Duration:          months,
			InterestMethod:    prod.InterestMethod,
			InterestRate:      prod.InterestRate,
			RepaymentCycle:    prod.RepaymentCycle,
			TotalInterest:     totalInterest,
			TotalRepayments:   months,
			InstallmentAmount: installment,
			TotalAmount:       total,
			BalanceRemaining:  total,
			Status:            "active",
		}
		if err := db.Create(&l).Error; err != nil {
			return err
		}

		// Mark the motorcycle as sold so it disappears from the
		// "available" filter — matches real-world behaviour.
		db.Model(&models.Motorcycle{}).Where("id = ?", moto.ID).Update("status", "on_loan")

		// Generate schedules; mark roughly the months that should be
		// past as paid (with some realistic partials so it looks live).
		monthsElapsed := int(time.Since(startDate).Hours() / 24 / 30)
		if monthsElapsed > months {
			monthsElapsed = months
		}
		balance := total
		for j := 1; j <= months; j++ {
			due := startDate.AddDate(0, j, 0)
			interestPortion := balance * prod.InterestRate / 100 / 12
			principalPortion := installment - interestPortion
			balanceBefore := balance
			balance -= principalPortion
			if balance < 0 {
				balance = 0
			}
			sched := models.RepaymentSchedule{
				LoanID: l.ID, InstallmentNumber: j, DueDate: due,
				PrincipalAmount: principalPortion, InterestAmount: interestPortion,
				TotalAmount:   installment,
				BalanceBefore: balanceBefore, BalanceAfter: balance,
			}
			if j <= monthsElapsed {
				// Most installments paid in full; ~15% partial; ~5% overdue.
				switch r := rng.Intn(100); {
				case r < 80:
					paidAt := due.AddDate(0, 0, -rng.Intn(3))
					sched.IsPaid = true
					sched.PaidAmount = installment
					sched.PaidAt = &paidAt
				case r < 95:
					sched.PaidAmount = installment * 0.5
					sched.IsOverdue = due.Before(time.Now())
				default:
					sched.IsOverdue = due.Before(time.Now())
					sched.DaysPastDue = int(time.Since(due).Hours() / 24)
				}
			}
			if err := db.Create(&sched).Error; err != nil {
				return err
			}
			if sched.IsPaid {
				if err := db.Create(&models.Repayment{
					LoanID: l.ID, ScheduleID: &sched.ID,
					Amount: installment, CollectionDate: *sched.PaidAt,
				}).Error; err != nil {
					return err
				}
			} else if sched.PaidAmount > 0 {
				if err := db.Create(&models.Repayment{
					LoanID: l.ID, ScheduleID: &sched.ID,
					Amount: sched.PaidAmount, CollectionDate: due,
				}).Error; err != nil {
					return err
				}
			}
		}

		// Update loan balance based on payments so far.
		var paid float64
		db.Model(&models.Repayment{}).Where("loan_id = ?", l.ID).
			Select("COALESCE(SUM(amount), 0)").Row().Scan(&paid)
		db.Model(&models.Loan{}).Where("id = ?", l.ID).Update("balance_remaining", total-paid)
	}
	return nil
}

// ─── cash sales (motorcycles sold for cash, not financed) ──────────────
func seedCashSales(db *gorm.DB, biz models.Business, branch models.Branch, motorcycles []models.Motorcycle, staff staffSet) error {
	if len(motorcycles) < 8 {
		return nil
	}
	// Last 3 motorcycles in the inventory sold for cash.
	cashIdx := []int{5, 6, 7}
	customers := []struct{ name, phone string }{
		{"Robert Mukasa", "+256700222301"},
		{"Aisha Namusoke", "+256700222302"},
		{"Charles Kaggwa", "+256700222303"},
	}
	for k, idx := range cashIdx {
		moto := motorcycles[idx]
		var ex models.CashSale
		if err := db.Where("motorcycle_id = ?", moto.ID).First(&ex).Error; err == nil {
			continue
		}
		c := customers[k]
		soldAt := time.Now().AddDate(0, 0, -7*(k+1))
		err := db.Create(&models.CashSale{
			BusinessID:    biz.ID,
			BranchID:      branch.ID,
			MotorcycleID:  moto.ID,
			SaleNumber:    fmt.Sprintf("GM-MC-%06d", 200000+k),
			SoldBy:        staff.manager.ID,
			CustomerName:  c.name,
			CustomerPhone: c.phone,
			ListPrice:     moto.SellingPrice,
			Total:         moto.SellingPrice,
			PaymentMethod: "cash",
			Notes:         "Cash purchase, full payment on day of pickup.",
			CreatedAt:     soldAt,
		}).Error
		if err != nil {
			return err
		}
		// Take the motorcycle off the available pool.
		db.Model(&models.Motorcycle{}).Where("id = ?", moto.ID).Update("status", "sold")
	}
	return nil
}

// ─── spare-parts POS history ──────────────────────────────────────────
func seedSpareSales(db *gorm.DB, biz models.Business, branch models.Branch, products []models.Product, staff staffSet, rng *rand.Rand) error {
	// Skip if there's already POS history (idempotency via count).
	var existing int64
	db.Model(&models.Sale{}).Where("branch_id = ?", branch.ID).Count(&existing)
	if existing > 0 {
		_ = biz // business is implied via branch; placate vet about unused param shape.
		return nil
	}

	for i := 0; i < 8; i++ {
		// 2–4 items per sale, varied product mix.
		nItems := 2 + rng.Intn(3)
		items := make([]models.SaleItem, 0, nItems)
		subtotal := 0.0
		for j := 0; j < nItems; j++ {
			p := products[rng.Intn(len(products))]
			qty := 1 + rng.Intn(3)
			lt := p.SellingPrice * float64(qty)
			items = append(items, models.SaleItem{
				ProductID: p.ID, Quantity: qty,
				UnitPrice: p.SellingPrice, UnitCost: p.CostPrice, LineTotal: lt,
			})
			subtotal += lt
		}
		pm := "cash"
		if rng.Intn(2) == 0 {
			pm = "mobile_money"
		}
		s := models.Sale{
			BranchID: branch.ID, CashierID: staff.cashier.ID,
			Segment: "spares", PaymentMethod: pm,
			Subtotal: subtotal, Total: subtotal,
			Items:     items,
			CreatedAt: time.Now().AddDate(0, 0, -rng.Intn(14)),
		}
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── daily boda (rental fleet + driver payments) ──────────────────────
func seedDailyBoda(db *gorm.DB, biz models.Business, branch models.Branch, staff staffSet, rng *rand.Rand) error {
	// Drivers — 3 stable identities.
	driverRows := []struct {
		name, phone, nin string
	}{
		{"Tom Lubega", "+256700333401", "CM86010001TL"},
		{"Ivan Kiggundu", "+256700333402", "CM88010002IK"},
		{"Henry Mwanga", "+256700333403", "CM84010003HM"},
	}
	drivers := make([]models.DailyBodaDriver, 0, len(driverRows))
	for _, d := range driverRows {
		var ex models.DailyBodaDriver
		if err := db.Where("business_id = ? AND phone = ?", biz.ID, d.phone).First(&ex).Error; err == nil {
			drivers = append(drivers, ex)
			continue
		}
		dr := models.DailyBodaDriver{
			BusinessID: biz.ID, BranchID: branch.ID,
			FullName: d.name, Phone: d.phone, NationalID: d.nin,
			DailyRate: models.DefaultDailyBodaRate, IsActive: true,
		}
		if err := db.Create(&dr).Error; err != nil {
			return err
		}
		drivers = append(drivers, dr)
	}

	motoRows := []struct {
		name, plate string
	}{
		{"Bajaj Boxer 100 (rental)", "UDA 101R"},
		{"TVS HLX 125 (rental)", "UDA 102R"},
		{"Bajaj Discover (rental)", "UDA 103R"},
	}
	bikes := make([]models.DailyBodaMotorcycle, 0, len(motoRows))
	for i, m := range motoRows {
		var ex models.DailyBodaMotorcycle
		if err := db.Where("business_id = ? AND number_plate = ?", biz.ID, m.plate).First(&ex).Error; err == nil {
			bikes = append(bikes, ex)
			continue
		}
		bike := models.DailyBodaMotorcycle{
			BusinessID: biz.ID, BranchID: branch.ID,
			Name: m.name, NumberPlate: m.plate,
			Status: models.DailyBodaStatusOccupied,
		}
		if err := db.Create(&bike).Error; err != nil {
			return err
		}
		// Assign each bike to its same-index driver.
		did := drivers[i].ID
		bike.AssignedDriverID = &did
		db.Save(&bike)
		bikes = append(bikes, bike)
	}

	// 7 days of payment history per driver. Most full, some partial/pending.
	for i, dr := range drivers {
		if i >= len(bikes) {
			break
		}
		bike := bikes[i]
		for day := 7; day >= 1; day-- {
			date := time.Now().AddDate(0, 0, -day)
			// Idempotency check by (driver, day, motorcycle).
			var ex models.DailyBodaPayment
			if err := db.Where("driver_id = ? AND payment_date BETWEEN ? AND ?",
				dr.ID,
				date.Truncate(24*time.Hour),
				date.Truncate(24*time.Hour).Add(24*time.Hour),
			).First(&ex).Error; err == nil {
				continue
			}
			r := rng.Intn(100)
			var amount, balance float64
			status := models.DailyBodaPaymentPaid
			switch {
			case r < 75:
				amount = dr.DailyRate
			case r < 92:
				amount = dr.DailyRate * 0.6
				balance = dr.DailyRate - amount
				status = models.DailyBodaPaymentPartial
			default:
				amount = 0
				balance = dr.DailyRate
				status = models.DailyBodaPaymentPending
			}
			if err := db.Create(&models.DailyBodaPayment{
				BusinessID: biz.ID, DriverID: dr.ID, MotorcycleID: bike.ID, BranchID: branch.ID,
				Amount: amount, DailyRate: dr.DailyRate, Balance: balance,
				PaymentDate:   date,
				PaymentMethod: "mobile_money",
				Status:        status,
				CollectedBy:   staff.loanOfficer.ID,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
