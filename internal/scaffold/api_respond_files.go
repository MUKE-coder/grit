package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeRespondFiles writes internal/respond.
//
// Split out of the main scaffold map so grit upgrade can call it too. The
// generator emits respond.WriteError into every resource handler, so a project
// carrying an older copy of this package gets generated code that does not
// compile, and the error names a helper the reader has never heard of.
//
// The same argument as internal/money: upgrade does not regenerate API code in
// general, and makes an exception for the packages the generator writes calls
// into.
func writeRespondFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "respond", "respond.go"): apiRespondGo(),
		// Generated from Grit's error catalogue: the typed codes, the status each
		// one always carries, and respond.Fail, which takes the status from the
		// catalogue so a handler cannot pair a code with the wrong one.
		filepath.Join(apiRoot, "internal", "respond", "codes.go"): apiRespondCodesGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
