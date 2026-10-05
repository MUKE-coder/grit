package pdf

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Invoice is the data shape RenderInvoice consumes. Build one from
// your domain model in the handler and pass it through.
type Invoice struct {
	Number    string // "INV-202605-0001"
	IssueDate time.Time
	DueDate   time.Time
	BillTo    Party
	From      Party // your company — optional, shown in the header area
	Items     []LineItem
	Subtotal  float64
	Tax       float64 // tax amount (not rate)
	Total     float64
	Paid      float64 // amount already paid; if > 0, an "Outstanding" line is added
	Currency  string  // "UGX", "USD", etc. — prefixed to every amount
	Notes     string  // free-text footer notes
	Status    string  // shown in the document footer ("paid", "overdue", "draft")
}

// Party is a name + free-text contact block (phone, email, address).
type Party struct {
	Name    string
	Contact string
}

// LineItem is one row in the invoice's items table.
type LineItem struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	Total       float64
}

// RenderInvoice returns the invoice as PDF bytes ready to stream to
// the response writer. Composition over inheritance: it's just a Doc
// with the section helpers called in order — copy this file as a
// starting point for receipts / leases / statements / quotes.
//
//	GET /api/invoices/:id/pdf
//	    inv := h.Service.GetByID(c.Param("id"))
//	    bytes, _ := pdf.RenderInvoice(toInvoice(inv))
//	    c.Data(200, "application/pdf", bytes)
func RenderInvoice(inv Invoice) ([]byte, error) {
	d := New()

	d.Header("INVOICE", inv.Number)
	d.TwoColumnKV("BILL TO", inv.BillTo.Name, "ISSUE DATE", inv.IssueDate.Format("2 Jan 2006"))
	if inv.BillTo.Contact != "" {
		d.SetFont("Helvetica", "", 10)
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.CellFormat(95, 5, inv.BillTo.Contact, "", 0, "L", false, 0, "")
		d.SetTextColor(0, 0, 0)
		d.SetFont("Helvetica", "B", 9)
		d.CellFormat(0, 5, "DUE DATE", "", 1, "L", false, 0, "")
		d.CellFormat(95, 5, "", "", 0, "L", false, 0, "")
		d.SetFont("Helvetica", "", 10)
		d.CellFormat(0, 5, inv.DueDate.Format("2 Jan 2006"), "", 1, "L", false, 0, "")
	}
	d.Ln(8)

	// Items table
	rows := make([][]string, len(inv.Items))
	for i, it := range inv.Items {
		rows[i] = []string{
			it.Description,
			strconv.FormatFloat(it.Quantity, 'f', -1, 64),
			formatAmount(it.UnitPrice),
			formatAmount(it.Total),
		}
	}
	d.Table(
		[]string{"DESCRIPTION", "QTY", "UNIT", "TOTAL"},
		rows,
		[]float64{105, 15, 25, 0},
		[]string{"L", "R", "R", "R"},
	)
	d.Ln(4)

	// Totals
	totals := []TotalLine{
		{Label: "Subtotal", Value: inv.Currency + " " + formatAmount(inv.Subtotal)},
	}
	if inv.Tax > 0 {
		totals = append(totals, TotalLine{Label: "Tax", Value: inv.Currency + " " + formatAmount(inv.Tax)})
	}
	totals = append(totals, TotalLine{Label: "Total", Value: inv.Currency + " " + formatAmount(inv.Total), Bold: true})
	if inv.Paid > 0 {
		totals = append(totals,
			TotalLine{Label: "Paid", Value: inv.Currency + " " + formatAmount(inv.Paid)},
			TotalLine{Label: "Outstanding", Value: inv.Currency + " " + formatAmount(inv.Total-inv.Paid), Bold: true},
		)
	}
	d.Totals(totals)

	d.Notes(inv.Notes)

	footer := fmt.Sprintf("Generated %s", time.Now().Format("2 Jan 2006 15:04"))
	if inv.Status != "" {
		footer += " · Status: " + inv.Status
	}
	d.Footer(footer)

	return d.Bytes()
}

// formatAmount renders 1234567.89 as "1,234,567.89" — matches the
// thousands-separator convention used by the export package.
func formatAmount(n float64) string {
	s := strconv.FormatFloat(n, 'f', 2, 64)
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	neg := strings.HasPrefix(intPart, "-")
	if neg {
		intPart = intPart[1:]
	}
	var out []byte
	// Indexed by byte rather than ranged by rune: intPart is digits from
	// FormatFloat, so the two are the same here, and mixing a rune loop with
	// byte arithmetic is how the first multi-byte character someone passes in
	// gets truncated.
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, intPart[i])
	}
	result := string(out) + "." + parts[1]
	if neg {
		return "-" + result
	}
	return result
}
