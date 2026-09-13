package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func handlerWithExport(block string) string {
	return "package handlers\n\nfunc (h *ContactHandler) Export(c *gin.Context) {\n\tformat := c.DefaultQuery(\"format\", \"csv\")\n\n" +
		block + "\n\t// CSV streams\n}\n"
}

func TestRepairXLSXExport(t *testing.T) {
	for name, audit := range map[string]string{
		"plain":       "",
		"audit-reads": "\t\taudit.ReadCount(c, \"contacts\", len(all))\n",
	} {
		old := handlerWithExport(xlsxExportBlock(oldXLSXExportBlock, "Contact", "contacts", audit))
		want := handlerWithExport(xlsxExportBlock(newXLSXExportBlock, "Contact", "contacts",
			strings.Replace(audit, "len(all)", "sheet.Written()", 1)))
		out, fixed, warn := repairXLSXExportSource(old)
		if len(warn) > 0 || len(fixed) != 1 {
			t.Fatalf("%s: fixed %v, warned %v", name, fixed, warn)
		}
		if out != want {
			t.Errorf("%s: the repaired handler is not the streaming one:\n%s", name, out)
		}
		if strings.Contains(out, "append(all, rows...)") || strings.Contains(out, "len(all)") {
			t.Errorf("%s: rows are still collected", name)
		}
		if _, err := format.Source([]byte(out)); err != nil {
			t.Errorf("%s: not valid Go: %v", name, err)
		}
		if again, fixed, _ := repairXLSXExportSource(out); again != out || len(fixed) > 0 {
			t.Errorf("%s: a second upgrade changed the handler again", name)
		}
	}
}

func TestRepairXLSXExportsEveryResourceInAFile(t *testing.T) {
	src := handlerWithExport(xlsxExportBlock(oldXLSXExportBlock, "Contact", "contacts", "")) +
		handlerWithExport(xlsxExportBlock(oldXLSXExportBlock, "Group", "groups", ""))[len("package handlers\n\n"):]
	out, fixed, warn := repairXLSXExportSource(src)
	if len(warn) > 0 || len(fixed) != 1 || strings.Count(out, "export.NewXLSXStream(opts)") != 2 {
		t.Errorf("fixed %v, warned %v, streams %d", fixed, warn, strings.Count(out, "export.NewXLSXStream(opts)"))
	}
}

func TestRepairXLSXExportLeavesAnEditedBranch(t *testing.T) {
	src := handlerWithExport(strings.Replace(xlsxExportBlock(oldXLSXExportBlock, "Contact", "contacts", ""),
		"return\n\t\t}\n\t\tc.Header", "return\n\t\t}\n\t\tlog.Println(len(all))\n\t\tc.Header", 1))
	if out, fixed, warn := repairXLSXExportSource(src); out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshExportPackageStreams(t *testing.T) {
	src := apiExportGo()
	for _, want := range []string{"func NewXLSXStream(", "NewStreamWriter(", "ErrTooManyRows", "func (x *XLSXStream) Finish("} {
		if !strings.Contains(src, want) {
			t.Errorf("export.go is missing %s", want)
		}
	}
	if strings.Contains(src, "f.SetCellValue(sheet, cell, val)") {
		t.Error("export.go still builds the sheet cell by cell")
	}
	for name, s := range map[string]string{"export.go": src, "export_test.go": apiExportTestGo()} {
		if _, err := format.Source([]byte(s)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
