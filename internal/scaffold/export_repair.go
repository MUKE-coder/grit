package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// oldXLSXExportBlock is the XLSX branch grit generate wrote into a resource's
// Export handler before v3.255.0. {{AUDIT}} is empty, or the --audit-reads line.
const oldXLSXExportBlock = `	if format == "xlsx" {
		// excelize has no streaming writer, so the sheet is built in memory.
		var all []models.{{Pascal}}
		if err := h.service().Export(h.ctx(c), search, func(rows []models.{{Pascal}}) error {
			all = append(all, rows...)
			return nil
		}); err != nil {
			h.fail(c, err, "Failed to export {{plural}}")
			return
		}
{{AUDIT}}		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{plural}}.xlsx"` + "`" + `)
		if err := export.XLSX(c.Writer, all, opts); err != nil {
			log.Printf("export {{plural}} as xlsx: %v", err)
		}
		return
	}
`

// newXLSXExportBlock is the same branch as grit generate writes it now.
const newXLSXExportBlock = `	if format == "xlsx" {
		// Rows go into the workbook batch by batch, through a writer that moves
		// them to a temporary file, so memory stays flat however many match.
		sheet, err := export.NewXLSXStream(opts)
		if err != nil {
			h.fail(c, err, "Failed to export {{plural}}")
			return
		}
		defer func() {
			if err := sheet.Close(); err != nil {
				log.Printf("export {{plural}} as xlsx: removing temporary files: %v", err)
			}
		}()
		if err := h.service().Export(h.ctx(c), search, func(rows []models.{{Pascal}}) error {
			return sheet.Rows(rows)
		}); err != nil {
			h.fail(c, err, "Failed to export {{plural}}")
			return
		}
{{AUDIT}}		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{plural}}.xlsx"` + "`" + `)
		if err := sheet.Finish(c.Writer); err != nil {
			log.Printf("export {{plural}} as xlsx: %v", err)
		}
		return
	}
`

var xlsxExportHeadRe = regexp.MustCompile(`\tif format == "xlsx" \{\n\t\t// excelize has no streaming writer, so the sheet is built in memory\.\n\t\tvar all \[\]models\.(\w+)\n[\s\S]*?"Failed to export ([^"]+)"\)`)

// repairXLSXExports brings generated resources up to the fix for H17 in the
// contact-app review: an XLSX export collected every matching row, then built
// the sheet cell by cell in memory. 300,000 contacts took an API from 80 MB to
// 1.3 GB, and a few concurrent exports would take it down.
//
// internal/export is framework code and arrives whole. The handlers are
// generated but the developer's, so only a branch that is exactly what grit
// generate wrote is replaced.
func repairXLSXExports(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileContains(filepath.Join(apiRoot, "internal", "export", "export.go"), "func NewXLSXStream(") {
		return nil
	}
	handlers, err := filepath.Glob(filepath.Join(apiRoot, "internal", "handlers", "*.go"))
	if err != nil {
		return err
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, path := range handlers {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.Contains(string(data), "export.XLSX(c.Writer, all, opts)") {
			continue
		}
		if err := repairSourceFile(root, m, path, repairXLSXExportSource); err != nil {
			return err
		}
	}
	return nil
}

func xlsxExportBlock(block, pascal, plural, audit string) string {
	return strings.NewReplacer("{{Pascal}}", pascal, "{{plural}}", plural, "{{AUDIT}}", audit).Replace(block)
}

func repairXLSXExportSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "export.XLSX(c.Writer, all, opts)") {
		return src, nil, nil
	}
	out, replaced := src, 0
	for _, match := range xlsxExportHeadRe.FindAllStringSubmatch(src, -1) {
		pascal, plural := match[1], match[2]
		for _, audit := range []string{"", "\t\taudit.ReadCount(c, \"" + plural + "\", len(all))\n"} {
			old := xlsxExportBlock(oldXLSXExportBlock, pascal, plural, audit)
			if strings.Count(out, old) != 1 {
				continue
			}
			newAudit := strings.Replace(audit, "len(all)", "sheet.Written()", 1)
			out = strings.Replace(out, old, xlsxExportBlock(newXLSXExportBlock, pascal, plural, newAudit), 1)
			replaced++
			break
		}
	}
	var warn []string
	if strings.Contains(out, "export.XLSX(c.Writer, all, opts)") {
		warn = []string{"an XLSX export here is not the code grit generate wrote, so it still builds the whole sheet in memory: use export.NewXLSXStream and add each batch with Rows"}
	}
	if replaced == 0 {
		return src, nil, warn
	}
	return out, []string{"the XLSX export streams rows into the workbook instead of holding them all in memory"}, warn
}
