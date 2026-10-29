package scaffold

// apiPDFTestGo emits internal/pdf/pdf_test.go.
//
// Every invoice came out as two pages, the second holding nothing but the
// footer line, and nothing failed: the function returned bytes, the handler
// served them, and the only way to see it was to open the file (grit#136).
//
// The cause is arithmetic. fpdf.New sets an automatic page break 20mm from
// the bottom, 277mm on A4; Footer does SetY(-25) and writes a 5mm cell, which
// ends at exactly 277. One millimetre either way and nobody would have
// noticed for another year.
func apiPDFTestGo() string {
	return `package pdf

import (
	"bytes"
	"testing"
	"time"
)

// pageCount counts /Type /Page objects, which is how many pages the document
// has. Reading the bytes rather than a viewer because this has to fail in CI.
func pageCount(b []byte) int {
	return bytes.Count(b, []byte("/Type /Page\n")) + bytes.Count(b, []byte("/Type /Page "))
}

func oneLineInvoice(number string) Invoice {
	now := time.Now()
	return Invoice{
		Number: number, IssueDate: now, DueDate: now,
		BillTo:   Party{Name: "A Customer"},
		Items:    []LineItem{{Description: "One line", Quantity: 1, UnitPrice: 1, Total: 1}},
		Subtotal: 1, Total: 1, Currency: "UGX",
	}
}

// A short invoice is one page. It used to be two: the footer sat exactly on
// the automatic page break, so fpdf started a new page to write it and the
// invoice ended with an empty space where the footer should have been.
func TestInvoiceIsOnePage(t *testing.T) {
	b, err := RenderInvoice(oneLineInvoice("INV-1"))
	if err != nil {
		t.Fatalf("rendering: %v", err)
	}
	if n := pageCount(b); n != 1 {
		t.Errorf("a one-line invoice came out as %d pages, want 1: the footer "+
			"is being pushed onto a page of its own (grit#136)", n)
	}
}

// And a long one still paginates. The fix turns the automatic break off
// around the footer and puts it back, so a document that genuinely needs more
// pages still gets them; leaving it off would turn this into one very long
// page with the rows running off the bottom.
func TestALongInvoiceStillPaginates(t *testing.T) {
	inv := oneLineInvoice("INV-2")
	inv.Items = make([]LineItem, 120)
	for i := range inv.Items {
		inv.Items[i] = LineItem{Description: "A line of the invoice", Quantity: 1, UnitPrice: 1, Total: 1}
	}
	inv.Subtotal, inv.Total = 120, 120

	b, err := RenderInvoice(inv)
	if err != nil {
		t.Fatalf("rendering: %v", err)
	}
	if n := pageCount(b); n < 2 {
		t.Errorf("120 lines came out as %d page(s): the automatic page break "+
			"has been left off for the whole document (grit#136)", n)
	}
}
`
}
