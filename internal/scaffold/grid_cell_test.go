package scaffold

import (
	"strings"
	"testing"
)

// Two reports about the bulk-create grid, both reproduced in a running admin
// against a real API before anything was changed.

// A large integer is readable.
//
// `unit_price:int` rendered as 1750000. The cell was <input type="number">,
// and the HTML spec requires that element's value to be a valid
// floating-point number, so a grouped string assigned to it is discarded and
// the field goes blank. No formatting around the cell could have worked: the
// input type had to change.
func TestTheNumberCellGroupsThousands(t *testing.T) {
	src := adminDataGrid()

	if !strings.Contains(src, "function groupDigits(value: number): string {") {
		t.Fatal("there is no grouping helper, so an integer column still reads 1750000")
	}
	if !strings.Contains(src, "toLocaleString(") {
		t.Error("groupDigits does not use toLocaleString, so the separator is not the reader's own")
	}

	// The number cell is its own component, and not a native number input:
	// that element cannot carry a grouped value at all.
	at := strings.Index(src, "function NumberCell(")
	if at == -1 {
		t.Fatal("there is no NumberCell")
	}
	body := src[at:]
	if end := strings.Index(body, "\nfunction "); end != -1 {
		body = body[:end]
	}
	if strings.Contains(body, `type="number"`) {
		t.Error("NumberCell is still a native number input, which cannot display a separator")
	}
	if !strings.Contains(body, `inputMode="decimal"`) {
		t.Error("NumberCell has no inputMode, so a phone shows the full keyboard for a number column")
	}
	// Raw while editing. A separator appearing under the caret moves the
	// caret, and this grid exists to be typed in quickly.
	if !strings.Contains(body, "onFocus=") || !strings.Contains(body, "onBlur=") {
		t.Error("NumberCell does not switch to the raw value while it is being edited")
	}
	// An honest role. role="spinbutton" on a text input makes an empty cell
	// report aria-valuenow 0, which is what the accessibility tree showed.
	if strings.Contains(body, `role="spinbutton"`) {
		t.Error(`NumberCell claims role="spinbutton", so an empty cell announces as 0`)
	}

	// Money is the column most worth grouping and had the same input type.
	money := src[strings.Index(src, "function MoneyCell("):]
	if end := strings.Index(money, "\nfunction "); end != -1 {
		money = money[:end]
	}
	if strings.Contains(money, `type="number"`) {
		t.Error("MoneyCell is still a native number input, so a price column has no separators either")
	}
}

// A pasted number survives, which is the half that was losing data.
//
// parseCell did Number("1,250,000"), which is NaN, which JSON.stringify writes
// as null. Pasting a column out of a spreadsheet, the likeliest way anybody
// fills this grid, silently emptied every cell over 999.
func TestAPastedNumberKeepsItsSeparators(t *testing.T) {
	src := adminDataGrid()

	if !strings.Contains(src, "function toNumber(text: string): number | null {") {
		t.Fatal("there is no separator-tolerant parse, so a pasted 1,250,000 is still NaN")
	}

	at := strings.Index(src, "function parseCell(")
	if at == -1 {
		t.Fatal("there is no parseCell")
	}
	body := src[at:]
	if end := strings.Index(body, "\nfunction "); end != -1 {
		body = body[:end]
	}
	if strings.Contains(body, "return trimmed === \"\" ? null : Number(trimmed);") {
		t.Error("parseCell still calls Number on the raw text, so a formatted paste becomes NaN")
	}
	if !strings.Contains(body, "toNumber(") {
		t.Error("parseCell does not use toNumber")
	}

	// NaN never reaches a row. It serialises to null, so returning it only
	// delays the same answer past the point where somebody could see it.
	if !strings.Contains(src, "return Number.isFinite(n) ? n : null;") {
		t.Error("toNumber can return NaN")
	}
}

// An upload that fails says so where somebody can see it.
//
// The catch set `failed` and the only thing rendering it was
// <span role="alert" className="sr-only">. Verified against an API answering
// 503 to the presign: the grid was pixel-identical before and after the
// upload, which is exactly the report, "it simply fails silently".
func TestAFailedUploadIsVisibleInTheGrid(t *testing.T) {
	src := adminDataGrid()

	if strings.Contains(src, `<span role="alert" className="sr-only">`) {
		t.Error("the upload failure is still screen-reader only, so a sighted operator " +
			"sees the spinner stop and nothing else")
	}
	if !strings.Contains(src, "text-danger") {
		t.Error("nothing in the cell is marked as failed")
	}

	// The message the operator needs is a sentence naming the setting that is
	// wrong, and it is in the response body. An axios rejection's own message
	// is "Request failed with status code 503", which is not actionable.
	if !strings.Contains(src, `getApiErrorMessage(err, "Upload failed")`) {
		t.Error("the cell reports axios's generic message rather than the API's own")
	}
	if !strings.Contains(src, `import { getApiErrorMessage } from "@/lib/api-core";`) {
		t.Error("getApiErrorMessage is used and not imported")
	}

	// A cell is 112 pixels wide, so the sentence goes to the dialog's error
	// bar. Both grids that render a DataGrid have to pass the channel, or the
	// one that does not is silent again.
	if !strings.Contains(src, "onCellError?.(message)") {
		t.Error("the cell does not pass its message up, so the long form is shown nowhere")
	}
	for name, grid := range map[string]string{
		"bulk-create-grid": adminBulkCreateGrid(),
		"bulk-edit-grid":   adminBulkEditGrid(),
	} {
		if !strings.Contains(grid, "onCellError={setError}") {
			t.Errorf("%s does not read onCellError, so an upload failure there is silent", name)
		}
	}
}
