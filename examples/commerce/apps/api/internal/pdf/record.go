package pdf

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Value formats a single model field for printing: times as "2 Jan 2006",
// booleans as Yes/No, nil as an em dash, everything else via %v. Keeps the
// generated handlers free of per-type formatting noise.
func Value(v any) string {
	switch t := v.(type) {
	case nil:
		return "—"
	case time.Time:
		if t.IsZero() {
			return "—"
		}
		return t.Format("2 Jan 2006")
	case *time.Time:
		if t == nil || t.IsZero() {
			return "—"
		}
		return t.Format("2 Jan 2006")
	case bool:
		if t {
			return "Yes"
		}
		return "No"
	case string:
		if strings.TrimSpace(t) == "" {
			return "—"
		}
		return t
	}
	s := fmt.Sprintf("%v", v)
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

// Display renders a related record (a belongs_to association) as a human
// label. It reflects for the first of Name / Title / Subject / Label / Email /
// Number / Code that exists and is non-empty, falling back to the ID. Taking
// "any" means a handler can pass any association without the renderer knowing
// that model's shape.
func Display(v any) string {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "—"
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return Value(v)
	}
	for _, name := range []string{"Name", "Title", "Subject", "Label", "Email", "Number", "Code"} {
		f := rv.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.String && strings.TrimSpace(f.String()) != "" {
			return f.String()
		}
	}
	if f := rv.FieldByName("ID"); f.IsValid() && f.Kind() == reflect.String {
		if s := f.String(); s != "" {
			return s
		}
	}
	return "—"
}

// Field is one label/value pair in a record's detail grid.
type Field struct {
	Label string
	Value string
}

// Section is a titled table inside a record — an invoice's line items, an
// order's shipments, a booking's guests.
type Section struct {
	Title   string
	Headers []string
	Rows    [][]string
	// Aligns is per-column: "L", "C" or "R". Empty means all left.
	Aligns []string
	// Widths is per-column in mm. Empty means evenly distributed.
	Widths []float64
}

// Record is a whole printable document. Build one in a handler from your
// model and hand it to RenderRecord — nothing here is model-specific, so the
// same shape prints an invoice, a receipt, a work order or a patient chart.
type Record struct {
	// Title is the big word at the top ("INVOICE", "ORDER").
	Title string
	// Subtitle sits under the title — usually the record's identifier.
	Subtitle string
	// Brand is the app/company name shown in the repeating header.
	Brand string
	// Fields render as a two-column detail grid under the header.
	Fields []Field
	// Sections render in order as titled tables.
	Sections []Section
	// Totals renders a right-aligned stack after the sections.
	Totals []TotalLine
	// Notes is free text at the end.
	Notes string
	// FooterNote sits bottom-left on every page, opposite the page number.
	FooterNote string
}

// RenderRecord lays a Record out as PDF bytes: a repeating header band
// (brand + identifier) and footer (note + "Page N of M") on every page, the
// title block, a two-up field grid, each section as a table, then totals and
// notes. Long tables page-break naturally and the header/footer follow.
func RenderRecord(r Record) ([]byte, error) {
	d := New()

	header := r.Brand
	if header == "" {
		header = r.Title
	}
	d.RunningHeader(header, r.Subtitle)
	footer := r.FooterNote
	if footer == "" {
		footer = r.Title
	}
	d.RunningFooter(footer)

	d.Header(r.Title, r.Subtitle)

	// Detail grid, two pairs per row so a record with many columns stays
	// compact instead of running one-per-line down the page.
	for i := 0; i < len(r.Fields); i += 2 {
		if i+1 < len(r.Fields) {
			d.TwoColumnKV(
				r.Fields[i].Label, r.Fields[i].Value,
				r.Fields[i+1].Label, r.Fields[i+1].Value,
			)
			continue
		}
		d.KV(r.Fields[i].Label, r.Fields[i].Value)
	}

	for _, s := range r.Sections {
		if len(s.Rows) == 0 {
			continue
		}
		if s.Title != "" {
			d.Ln(3)
			d.SetFont("Helvetica", "B", 11)
			d.SetTextColor(0, 0, 0)
			d.CellFormat(0, 6, s.Title, "", 1, "L", false, 0, "")
			d.Ln(1)
		}
		aligns := s.Aligns
		if len(aligns) == 0 {
			aligns = make([]string, len(s.Headers))
			for i := range aligns {
				aligns[i] = "L"
			}
		}
		widths := s.Widths
		if len(widths) == 0 {
			widths = evenWidths(d, len(s.Headers))
		}
		d.Table(s.Headers, s.Rows, widths, aligns)
	}

	if len(r.Totals) > 0 {
		d.Totals(r.Totals)
	}
	if strings.TrimSpace(r.Notes) != "" {
		d.Notes(r.Notes)
	}

	return d.Bytes()
}

// evenWidths splits the printable width evenly across n columns.
func evenWidths(d *Doc, n int) []float64 {
	if n <= 0 {
		return nil
	}
	left, _, right, _ := d.GetMargins()
	pageW, _ := d.GetPageSize()
	each := (pageW - left - right) / float64(n)
	out := make([]float64, n)
	for i := range out {
		out[i] = each
	}
	return out
}
