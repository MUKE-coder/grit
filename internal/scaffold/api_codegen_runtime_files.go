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
		filepath.Join(apiRoot, "internal", "paginate", "paginate.go"): apiPaginateGo(),
		filepath.Join(apiRoot, "internal", "events", "events.go"):     apiEventsGo(),
		filepath.Join(apiRoot, "internal", "export", "export.go"):     apiExportGo(),
		filepath.Join(apiRoot, "internal", "ids", "ids.go"):           apiIDsGo(),
		filepath.Join(apiRoot, "internal", "ids", "ids_test.go"):      apiIDsTestGo(),
		filepath.Join(apiRoot, "internal", "pdf", "pdf.go"):           apiPDFGo(),
		filepath.Join(apiRoot, "internal", "pdf", "invoice.go"):       apiPDFInvoiceGo(),
		filepath.Join(apiRoot, "internal", "pdf", "record.go"):        apiPDFRecordGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
