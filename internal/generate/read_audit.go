package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// A resource generated with --audit-reads records every read in the
// tamper-evident activity log: which rows each caller was shown, by list, by
// id, as a PDF or as an export.
//
// Found building a patient-records app. The activity log recorded every write
// and no read at all, so "who looked at this patient's chart" had no answer,
// which is the first question an access review asks.

// auditReadSnippets is what the handler template needs for --audit-reads. All
// empty without the flag, so a resource that did not ask for it is unchanged.
type auditReadSnippets struct {
	Import, List, One, ExportDecl, ExportCount, ExportMark, XLSXMark string
}

// Values are spliced in by a single-pass strings.Replacer, so they carry the
// resource name itself rather than a placeholder that would never be expanded.
func (g *Generator) auditReadSnippets(names Names) auditReadSnippets {
	if !g.Definition.AuditReads {
		return auditReadSnippets{}
	}
	plural := names.Plural
	return auditReadSnippets{
		Import: "\"" + g.Module + "/internal/audit\"\n\t",
		List: "\t// --audit-reads: which rows this caller was shown, in the\n" +
			"\t// tamper-evident activity log.\n" +
			"\tids := make([]string, 0, len(res.Data))\n" +
			"\tfor _, row := range res.Data {\n" +
			"\t\tids = append(ids, row.ID)\n" +
			"\t}\n" +
			"\taudit.Read(c, \"" + plural + "\", ids...)\n\n",
		One: "\n\t// --audit-reads: this read goes in the tamper-evident activity log.\n" +
			"\taudit.Read(c, \"" + plural + "\", item.ID)\n",
		ExportDecl:  "\texported := 0\n",
		ExportCount: "\t\texported += len(rows)\n",
		ExportMark: "\t// --audit-reads: an export is recorded as a count. Listing every id\n" +
			"\t// it held would make one entry the size of the table.\n" +
			"\taudit.ReadCount(c, \"" + plural + "\", exported)\n",
		XLSXMark: "\t\taudit.ReadCount(c, \"" + plural + "\", len(all))\n",
	}
}

// prepareReadAudit refuses a project whose activity log cannot record reads.
//
// The generated handlers would mark every read and the project's middleware
// would silently drop the marks, which looks finished and records nothing.
func (g *Generator) prepareReadAudit() error {
	api := g.APIRoot()
	auditGo := filepath.Join(api, "internal", "audit", "audit.go")
	activity := filepath.Join(api, "internal", "middleware", "activity.go")
	if readAuditFileHas(auditGo, "func Read(c *gin.Context") && readAuditFileHas(activity, "audit.ReadMarkOf(c)") {
		return nil
	}
	return fmt.Errorf("--audit-reads needs an activity log that records reads, and this "+
		"project's does not yet.\n\nRun grit upgrade, then generate %s again", g.Definition.Name)
}

func readAuditFileHas(path, needle string) bool {
	data, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(data), needle)
}
