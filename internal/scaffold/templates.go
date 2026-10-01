package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// Templates as files, rather than as string literals inside Go functions.
//
// internal/scaffold is 174,000 lines and most of it is Go, TypeScript, CSS and
// YAML inside Go raw strings, with one file at 10,437 lines. Nothing checks any
// of it until a project is scaffolded in CI, which is why a colour utility that
// compiled to nothing, a hook after an early return and a respond call without
// its import all shipped. A product review in October 2026 scored the
// maintainability of the framework itself 4 out of 10 for this, and it is the
// reason a second contributor cannot find the template for a given output
// without grepping a file the size of a small book.
//
// This is the first directory moved. The handler templates live under
// templates/api/handlers as ordinary files, which buys three things at once:
// an editor highlights and navigates them, gofmt and go/parser can be pointed
// at them (see templates_test.go), and `ls` answers "where does auth.go come
// from" without a grep.
//
// What it does not buy is a full type check. These files reference packages
// that exist only in a scaffolded project, so the compiler cannot see them
// here; that remains the job of the live suite, which scaffolds an app and
// builds it on every push. Syntax and formatting are caught here, at the point
// of editing, which is where the bugs above were introduced.
//
// Moving the rest is a directory at a time. The live suite protects each move,
// because a template that changes by one byte changes the app it generates.

//go:embed templates
var templateFS embed.FS

// tmpl returns a template by its path under templates/, without the .tmpl
// suffix: tmpl("api/handlers/auth.go").
//
// It panics on a missing template rather than returning an error. The argument
// is a constant in this package, so a miss is a typo a developer introduced
// thirty seconds ago and not a condition any caller could handle. Every one of
// them is covered by TestEveryTemplateResolves, so the panic is a developer's
// and never a user's.
func tmpl(name string) string {
	data, err := templateFS.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		panic(fmt.Sprintf("scaffold: no template %q: %v (see internal/scaffold/templates/)", name, err))
	}
	// Files on disk end with a newline and the string literals they replaced
	// did not, and a scaffolded project has to be byte-identical across the
	// move. One trailing newline is added back by the writer.
	return strings.TrimSuffix(string(data), "\n")
}

// templateNames lists every embedded template, for the tests that check them
// all rather than the ones somebody remembered.
func templateNames() ([]string, error) {
	var out []string
	err := fs.WalkDir(templateFS, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".tmpl") {
			return nil
		}
		out = append(out, strings.TrimSuffix(strings.TrimPrefix(p, "templates/"), ".tmpl"))
		return nil
	})
	return out, err
}

// isGoTemplate reports whether a template produces Go source, so the tests know
// which ones to parse.
func isGoTemplate(name string) bool {
	return path.Ext(name) == ".go"
}
