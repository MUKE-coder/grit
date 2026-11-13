package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeCodegenRuntimeFiles writes the framework packages that generated code
// compiles against.
//
// These are the packages a `grit generate resource` handler, service or model
// imports: pagination, the event bus, CSV/XLSX export, ID generation, PDF
// rendering. Nothing here is meant to be edited, and everything here is a
// direct dependency of code the generator emits tomorrow.
//
// They live in their own function because they have to run on upgrade, not
// only on scaffold. The generator always emits against the CLI's current
// templates, so a project whose internal/paginate is six versions old gets a
// handler referencing a paginate.Config field that its own copy of the struct
// does not declare, and the API stops compiling. The user did nothing wrong
// and there is nothing in the error to point at the real cause.
//
// This is the same gap that shipped media, recovery contacts and passkeys half
// finished: `grit upgrade` does not regenerate API code in general, which is
// right for anything a person might have edited and wrong for a package that
// exists purely to be imported. Splitting those two sets is what this
// function is.
//
// writeFile is manifest-guarded, so a file someone has genuinely edited is
// reported as a conflict rather than overwritten.
func writeCodegenRuntimeFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "paginate", "paginate.go"):      apiPaginateGo(),
		filepath.Join(apiRoot, "internal", "paginate", "paginate_test.go"): apiPaginateTestGo(),
		filepath.Join(apiRoot, "internal", "paginate", "counts_test.go"):   apiPaginateCountsTestGo(),
		filepath.Join(apiRoot, "internal", "paginate", "count.go"):         paginateCountGo(),
		// The insights panel's chart data: a count per period and a count per
		// value, over the list's own query so the chart describes the rows the
		// table is showing.
		filepath.Join(apiRoot, "internal", "paginate", "insights.go"):      tmpl("api/paginate/insights.go"),
		filepath.Join(apiRoot, "internal", "paginate", "insights_test.go"): tmpl("api/paginate/insights_test.go"),
		filepath.Join(apiRoot, "internal", "paginate", "count_test.go"):    paginateCountTestGo(),
		filepath.Join(apiRoot, "internal", "paginate", "search_index.go"):  paginateSearchIndexGo(),
		// ?filter=, which is the expression half of a list request. The
		// engine is internal/rules below; this is the twenty lines that put
		// a compiled expression on the query.
		filepath.Join(apiRoot, "internal", "paginate", "filter.go"):      tmpl("api/paginate/filter.go"),
		filepath.Join(apiRoot, "internal", "paginate", "filter_test.go"): tmpl("api/paginate/filter_test.go"),
		// ?fields= projects the marshalled rows; ?expand= adds a Preload per
		// relation, resolved through the same schema the filter uses.
		filepath.Join(apiRoot, "internal", "paginate", "project.go"):      tmpl("api/paginate/project.go"),
		filepath.Join(apiRoot, "internal", "paginate", "expand.go"):       tmpl("api/paginate/expand.go"),
		filepath.Join(apiRoot, "internal", "paginate", "project_test.go"): tmpl("api/paginate/project_test.go"),

		// The rate-limit table: labelled rules matched most specific first,
		// so a route nobody thought about is still bounded by a prefix. One
		// table for the whole API rather than a limit written at each route,
		// which is a limit somebody has to remember.
		filepath.Join(apiRoot, "internal", "ratelimit", "rules.go"):      tmpl("api/ratelimit/rules.go"),
		filepath.Join(apiRoot, "internal", "ratelimit", "labels.go"):     tmpl("api/ratelimit/labels.go"),
		filepath.Join(apiRoot, "internal", "ratelimit", "limiter.go"):    tmpl("api/ratelimit/limiter.go"),
		filepath.Join(apiRoot, "internal", "ratelimit", "rules_test.go"): tmpl("api/ratelimit/rules_test.go"),
		filepath.Join(apiRoot, "internal", "middleware", "ratelimit.go"): tmpl("api/middleware/ratelimit.go"),

		// The expression engine: one parser and one compiler, used by
		// ?filter= on a list endpoint and by a declared access rule. Both
		// resolve every identifier against the model's own GORM schema and
		// bind every value, which is the whole of the safety argument.
		filepath.Join(apiRoot, "internal", "rules", "lex.go"):        tmpl("api/rules/lex.go"),
		filepath.Join(apiRoot, "internal", "rules", "parse.go"):      tmpl("api/rules/parse.go"),
		filepath.Join(apiRoot, "internal", "rules", "schema.go"):     tmpl("api/rules/schema.go"),
		filepath.Join(apiRoot, "internal", "rules", "compile.go"):    tmpl("api/rules/compile.go"),
		filepath.Join(apiRoot, "internal", "rules", "context.go"):    tmpl("api/rules/context.go"),
		filepath.Join(apiRoot, "internal", "rules", "macros.go"):     tmpl("api/rules/macros.go"),
		filepath.Join(apiRoot, "internal", "rules", "values.go"):     tmpl("api/rules/values.go"),
		filepath.Join(apiRoot, "internal", "rules", "rules_test.go"): tmpl("api/rules/rules_test.go"),
		filepath.Join(apiRoot, "internal", "imports", "limit.go"):          importsLimitGo(),
		filepath.Join(apiRoot, "internal", "imports", "limit_test.go"):     importsLimitTestGo(),
		filepath.Join(apiRoot, "internal", "events", "events.go"):          apiEventsGo(),
		filepath.Join(apiRoot, "internal", "events", "durable.go"):         apiEventsDurableGo(),
		filepath.Join(apiRoot, "internal", "workflow", "workflow.go"):      apiWorkflowGo(),
		filepath.Join(apiRoot, "internal", "export", "export.go"):          apiExportGo(),
		filepath.Join(apiRoot, "internal", "export", "export_test.go"):     apiExportTestGo(),
		filepath.Join(apiRoot, "internal", "ids", "ids.go"):                apiIDsGo(),
		filepath.Join(apiRoot, "internal", "ids", "ids_test.go"):           apiIDsTestGo(),
		filepath.Join(apiRoot, "internal", "pdf", "pdf.go"):                apiPDFGo(),
		filepath.Join(apiRoot, "internal", "pdf", "invoice.go"):            apiPDFInvoiceGo(),
		filepath.Join(apiRoot, "internal", "pdf", "pdf_test.go"):           apiPDFTestGo(),
		filepath.Join(apiRoot, "internal", "pdf", "record.go"):             apiPDFRecordGo(),
		// Word documents. internal/pdf renders a styled page; this renders a
		// document with a page size, a named font, a line spacing and a page
		// break before each chapter, which is what a submission or a marking
		// scheme specifies and what no list export can produce (grit#212).
		// perm() and openToAnySignedIn(): what a hand-written route says it
		// requires, in a form `grit routes --check` can read (grit#179).
		filepath.Join(apiRoot, "internal", "routes", "perm.go"): tmpl("api/routes/perm.go"),

		filepath.Join(apiRoot, "internal", "docx", "docx.go"):      tmpl("api/docx/docx.go"),
		filepath.Join(apiRoot, "internal", "docx", "report.go"):    tmpl("api/docx/report.go"),
		filepath.Join(apiRoot, "internal", "docx", "docx_test.go"): tmpl("api/docx/docx_test.go"),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
