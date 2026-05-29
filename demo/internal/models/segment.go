package models

// Revenue segment tags for separating money flows between Grit's two business models.
// Every money-moving record (Sale, CashSale, Repayment, DailyBodaPayment) belongs to one segment.
const (
	SegmentSpares          = "spares"           // POS sales of spare parts (existing Sale)
	SegmentMotorcyclesCash = "motorcycles_cash" // Direct cash sale of a motorcycle (CashSale)
	SegmentLoans           = "loans"            // Loan repayment collections (Repayment)
	SegmentDailyBoda       = "daily_boda"       // Daily boda rental income (DailyBodaPayment)
)
