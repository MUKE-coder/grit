package scaffold

// apiExportTestGo emits internal/export/export_test.go.
func apiExportTestGo() string {
	return `package export

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"testing"

	"github.com/xuri/excelize/v2"
)

type exportRow struct {
	Name  string
	Email string
}

var exportColumns = []Column{{Header: "Name", Field: "Name"}, {Header: "Email", Field: "Email"}}

func TestXLSXStreamWritesEveryRow(t *testing.T) {
	sheet, err := NewXLSXStream(Options{Sheet: "People", Columns: exportColumns})
	if err != nil {
		t.Fatal(err)
	}
	defer sheet.Close()
	batch := make([]exportRow, 1000)
	for i := 0; i < 5; i++ {
		for j := range batch {
			batch[j] = exportRow{Name: fmt.Sprintf("person %d", i*1000+j), Email: "someone@example.com"}
		}
		if err := sheet.Rows(batch); err != nil {
			t.Fatal(err)
		}
	}
	if sheet.Written() != 5000 {
		t.Errorf("Written() = %d, want 5000", sheet.Written())
	}
	var buf bytes.Buffer
	if err := sheet.Finish(&buf); err != nil {
		t.Fatal(err)
	}
	book, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatal(err)
	}
	defer book.Close()
	rows, err := book.GetRows("People")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 5001 || rows[0][0] != "Name" || rows[5000][0] != "person 4999" {
		t.Errorf("the workbook has %d rows, header %q, last %q", len(rows), rows[0][0], rows[len(rows)-1][0])
	}
}

// Rows are not kept in the heap: 200,000 of them must fit in the memory of a
// few thousand.
func TestXLSXStreamMemoryStaysFlat(t *testing.T) {
	sheet, err := NewXLSXStream(Options{Columns: exportColumns})
	if err != nil {
		t.Fatal(err)
	}
	defer sheet.Close()
	batch := make([]exportRow, 1000)
	for j := range batch {
		batch[j] = exportRow{Name: "a name that takes up some room", Email: "someone.somewhere@example.com"}
	}
	var stats runtime.MemStats
	var peak uint64
	for i := 0; i < 200; i++ {
		if err := sheet.Rows(batch); err != nil {
			t.Fatal(err)
		}
		if i%40 == 39 {
			runtime.GC()
			runtime.ReadMemStats(&stats)
			if stats.HeapInuse > peak {
				peak = stats.HeapInuse
			}
		}
	}
	if peak > 64<<20 {
		t.Errorf("the heap reached %d MB while adding 200,000 rows; rows are being held in memory", peak>>20)
	}
	if err := sheet.Finish(io.Discard); err != nil {
		t.Fatal(err)
	}
}

func TestXLSXStreamRefusesMoreRowsThanASheetHolds(t *testing.T) {
	sheet, err := NewXLSXStream(Options{Columns: exportColumns})
	if err != nil {
		t.Fatal(err)
	}
	defer sheet.Close()
	sheet.next = excelize.TotalRows + 1
	if err := sheet.Rows([]exportRow{{Name: "one too many"}}); !errors.Is(err, ErrTooManyRows) {
		t.Errorf("got %v, want ErrTooManyRows", err)
	}
}

func TestXLSXStillWritesASlice(t *testing.T) {
	var buf bytes.Buffer
	if err := XLSX(&buf, []exportRow{{Name: "only", Email: "one@example.com"}}, Options{Columns: exportColumns}); err != nil {
		t.Fatal(err)
	}
	book, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatal(err)
	}
	defer book.Close()
	if rows, err := book.GetRows("Sheet1"); err != nil || len(rows) != 2 || rows[1][0] != "only" {
		t.Errorf("rows %v, err %v", rows, err)
	}
}
`
}
