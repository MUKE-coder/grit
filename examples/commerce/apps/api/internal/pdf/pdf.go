// Package pdf is a tiny styled-PDF builder backed by go-pdf/fpdf.
//
// The package exports two layers:
//
//  1. Doc primitives — Header, KV, Table, Totals, Notes, Footer — that
//     apply Grit's default styling (Helvetica, 20mm margins, blue
//     accent, A4 portrait). Compose them to build any document.
//
//  2. Pre-built templates — RenderInvoice (in invoice.go) — for the
//     common business-app cases. Copy + adapt these for receipts,
//     leases, statements, etc.
//
// When the helpers don't fit, the embedded *fpdf.Fpdf gives you the
// full underlying API. Call d.Bytes() at the end to finalize.
package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

// Doc wraps fpdf.Fpdf with section helpers + Grit's default colors.
// Mutate Accent on the returned Doc to retheme.
type Doc struct {
	*fpdf.Fpdf
	Accent [3]int // RGB; default Grit blue (30, 126, 245)
	Muted  [3]int // RGB; default neutral gray (110, 110, 110)
}

// New returns a fresh A4 portrait document with Grit's default styling.
// Adds the first page automatically — call d.AddPage() for additional
// pages.
func New() *Doc {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 18, 20)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 10)
	return &Doc{
		Fpdf:   pdf,
		Accent: [3]int{30, 126, 245},
		Muted:  [3]int{110, 110, 110},
	}
}

// Header writes the standard top-of-document title bar — accent-colored
// title in 22pt + a smaller secondary line below in muted gray.
//
//	d.Header("INVOICE", "INV-202605-0001")
//	d.Header("RECEIPT", "RCT-202605-0042")
func (d *Doc) Header(title, subtitle string) {
	d.SetFont("Helvetica", "B", 22)
	d.SetTextColor(d.Accent[0], d.Accent[1], d.Accent[2])
	d.CellFormat(0, 10, strings.ToUpper(title), "", 1, "L", false, 0, "")
	if subtitle != "" {
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.SetFont("Helvetica", "", 10)
		d.CellFormat(0, 5, subtitle, "", 1, "L", false, 0, "")
	}
	d.SetTextColor(0, 0, 0)
	d.Ln(6)
}

// KV writes a "label: value" pair. Label is bold + small caps style;
// value is regular weight on the next line. Used for "Bill To",
// "Issue Date", "Reference Number", etc.
func (d *Doc) KV(label, value string) {
	d.SetFont("Helvetica", "B", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(0, 5, strings.ToUpper(label), "", 1, "L", false, 0, "")
	d.SetFont("Helvetica", "", 10)
	d.SetTextColor(0, 0, 0)
	d.CellFormat(0, 5, value, "", 1, "L", false, 0, "")
	d.Ln(3)
}

// TwoColumnKV writes two KV pairs side by side — useful for fitting
// "BILL TO" + "ISSUE DATE" or "FROM" + "TO" on one row.
func (d *Doc) TwoColumnKV(leftLabel, leftValue, rightLabel, rightValue string) {
	d.SetFont("Helvetica", "B", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(95, 5, strings.ToUpper(leftLabel), "", 0, "L", false, 0, "")
	d.CellFormat(0, 5, strings.ToUpper(rightLabel), "", 1, "L", false, 0, "")
	d.SetFont("Helvetica", "", 10)
	d.SetTextColor(0, 0, 0)
	d.CellFormat(95, 5, leftValue, "", 0, "L", false, 0, "")
	d.CellFormat(0, 5, rightValue, "", 1, "L", false, 0, "")
	d.Ln(3)
}

// Table writes a styled table. headers + rows are matched by index.
// colWidths are in mm — pass 0 for the last column to fill remaining
// width. Header row gets a light gray background; data rows are plain.
//
//	d.Table(
//	    []string{"DESCRIPTION", "QTY", "UNIT", "TOTAL"},
//	    [][]string{
//	        {"Office rent — June", "1", "1,500,000", "1,500,000"},
//	        {"Service charge",      "1",   "120,000",   "120,000"},
//	    },
//	    []float64{105, 15, 25, 0},
//	    []string{"L", "R", "R", "R"},
//	)
func (d *Doc) Table(headers []string, rows [][]string, colWidths []float64, aligns []string) {
	if len(headers) == 0 || len(colWidths) != len(headers) {
		return
	}
	if len(aligns) != len(headers) {
		// Default all-left if alignment slice is malformed.
		aligns = make([]string, len(headers))
		for i := range aligns {
			aligns[i] = "L"
		}
	}

	// Header row
	d.SetFillColor(244, 244, 245)
	d.SetFont("Helvetica", "B", 9)
	d.SetTextColor(120, 120, 120)
	for i, h := range headers {
		end := 0
		if i == len(headers)-1 {
			end = 1
		}
		d.CellFormat(colWidths[i], 7, h, "", end, aligns[i], true, 0, "")
	}

	// Data rows
	d.SetTextColor(0, 0, 0)
	d.SetFont("Helvetica", "", 10)
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(colWidths) {
				break
			}
			end := 0
			if i == len(row)-1 {
				end = 1
			}
			d.CellFormat(colWidths[i], 6, cell, "", end, aligns[i], false, 0, "")
		}
	}
}

// TotalLine is one entry in a Totals stack.
type TotalLine struct {
	Label string
	Value string // pre-formatted with currency + thousands separators
	Bold  bool   // bold + accent color (for the grand total line)
}

// Totals writes a right-aligned totals stack. The last "Bold" line
// gets accent coloring + a slightly larger size — typically used for
// the grand total or outstanding balance.
func (d *Doc) Totals(lines []TotalLine) {
	for _, line := range lines {
		d.CellFormat(120, 6, "", "", 0, "L", false, 0, "")
		d.SetFont("Helvetica", "", 10)
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.CellFormat(20, 6, line.Label, "", 0, "R", false, 0, "")
		if line.Bold {
			d.SetFont("Helvetica", "B", 11)
			d.SetTextColor(d.Accent[0], d.Accent[1], d.Accent[2])
		} else {
			d.SetFont("Helvetica", "", 10)
			d.SetTextColor(0, 0, 0)
		}
		d.CellFormat(0, 6, line.Value, "", 1, "R", false, 0, "")
	}
	d.SetTextColor(0, 0, 0)
}

// Notes writes a "NOTES" header + the body text wrapped to page width.
// Skipped silently when text is empty.
func (d *Doc) Notes(text string) {
	if text == "" {
		return
	}
	d.Ln(4)
	d.SetFont("Helvetica", "B", 10)
	d.SetTextColor(0, 0, 0)
	d.CellFormat(0, 5, "NOTES", "", 1, "L", false, 0, "")
	d.SetFont("Helvetica", "", 10)
	d.SetTextColor(82, 82, 91)
	d.MultiCell(0, 5, text, "", "L", false)
	d.Ln(4)
}

// Footer writes a centered italic footer 25mm from the bottom of the
// CURRENT page only. For a footer that repeats on every page (and page
// numbers), use RunningFooter instead.
func (d *Doc) Footer(text string) {
	d.SetY(-25)
	d.SetFont("Helvetica", "I", 9)
	d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
	d.CellFormat(0, 5, text, "", 1, "C", false, 0, "")
}

// RunningHeader repeats a header band on EVERY page: the title on the
// left, an optional right-aligned line (document number, date), and a
// thin accent rule beneath. Call it before writing body content — fpdf
// invokes the callback as each page is added, including the first.
//
//	d.RunningHeader("INVOICE", "INV-202605-0001")
func (d *Doc) RunningHeader(title, right string) {
	draw := func() {
		d.SetY(10)
		d.SetFont("Helvetica", "B", 10)
		d.SetTextColor(d.Accent[0], d.Accent[1], d.Accent[2])
		d.CellFormat(0, 5, title, "", 0, "L", false, 0, "")
		if right != "" {
			d.SetFont("Helvetica", "", 9)
			d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
			d.CellFormat(0, 5, right, "", 0, "R", false, 0, "")
		}
		d.Ln(7)
		d.SetDrawColor(d.Accent[0], d.Accent[1], d.Accent[2])
		d.SetLineWidth(0.4)
		leftM, _, rightM, _ := d.GetMargins()
		w, _ := d.GetPageSize()
		y := d.GetY()
		d.Line(leftM, y, w-rightM, y)
		d.Ln(4)
	}
	d.SetHeaderFunc(draw)
	// New() already added page 1 before any header func existed, so fpdf
	// never invoked the callback for it. Draw it once now; the callback
	// covers every page added from here on.
	if d.PageNo() == 1 {
		draw()
	}
}

// RunningFooter repeats a footer on EVERY page: text on the left and
// "Page N of M" on the right. The page count uses fpdf's page-number
// alias, substituted when the document is finalized.
//
//	d.RunningFooter("Generated 2 Jun 2026 · Acme Ltd")
func (d *Doc) RunningFooter(text string) {
	d.AliasNbPages("")
	d.SetFooterFunc(func() {
		d.SetY(-15)
		d.SetFont("Helvetica", "I", 8)
		d.SetTextColor(d.Muted[0], d.Muted[1], d.Muted[2])
		d.CellFormat(0, 5, text, "", 0, "L", false, 0, "")
		d.CellFormat(0, 5, fmt.Sprintf("Page %d of {nb}", d.PageNo()), "", 0, "R", false, 0, "")
	})
}

// Bytes finalizes the document and returns the PDF bytes. Call this
// once at the very end — the underlying fpdf is not reusable after.
func (d *Doc) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := d.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return buf.Bytes(), nil
}
