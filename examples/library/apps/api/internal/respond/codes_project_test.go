package respond

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every respond.Code("...") written in this project is a declared code.
//
// Code is a string type, so respond.Code("PLAN_LIMIT_REACHED") compiles
// whether or not anything documents that code. It is not harmless: StatusOf
// does not know it, so respond.Fail answers 500, and a client error is
// reported as a server fault. Grit logs it when it happens; this catches it
// before it ships, which is the difference between a test run and a page at
// three in the morning.
//
// To add a code of your own, declare it in .grit/error-codes.yaml and run
// `grit sync`. Both this file's catalogue and the TypeScript the frontend
// switches on are generated from that, so the two cannot drift.
//
// Written as text rather than with go/parser deliberately: the subject is a
// string literal, and a parse would make this test need the package it is
// checking to compile, which it will even when the code is wrong.

// codeLiteral is respond.Code("X"), or Code("X") inside this package.
var codeLiteral = regexp.MustCompile(`\bCode\("([^"]*)"\)`)

func TestEveryCodeLiteralInThisProjectIsDeclared(t *testing.T) {
	root := apiRootForCodeScan(t)

	type site struct{ file, code string }
	var undeclared []site
	seen := map[site]bool{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Vendored and generated trees are not this project's handlers.
			switch d.Name() {
			case "vendor", "node_modules", ".git", "testdata":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// The two files that name codes in prose rather than using them: this
		// one, and the generated codes.go, whose doc comment on Fail shows
		// what an undeclared code looks like. Without this exclusion a fresh
		// project fails this test out of the box on Grit's own example, which
		// is the sort of false positive that gets a test deleted.
		switch filepath.Base(path) {
		case "codes_project_test.go", "codes.go":
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		for _, m := range codeLiteral.FindAllStringSubmatch(string(raw), -1) {
			if _, ok := Lookup(Code(m[1])); ok {
				continue
			}
			// Once per file and code. The same code used five times in one
			// handler is one thing to fix, and five lines saying so is five
			// times the noise for the same information.
			s := site{file: filepath.ToSlash(rel), code: m[1]}
			if seen[s] {
				continue
			}
			seen[s] = true
			undeclared = append(undeclared, s)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", root, err)
	}

	if len(undeclared) == 0 {
		return
	}
	sort.Slice(undeclared, func(i, j int) bool { return undeclared[i].file < undeclared[j].file })
	var lines []string
	for _, u := range undeclared {
		lines = append(lines, "  "+u.file+": "+u.code)
	}
	t.Errorf("%d error code(s) are used and not declared, so respond.Fail answers 500 for each "+
		"of them and the client is told a server fault:\n%s\n\nDeclare them in "+
		".grit/error-codes.yaml and run `grit sync`, or use one of the codes in codes.go.",
		len(undeclared), strings.Join(lines, "\n"))
}

// apiRootForCodeScan walks up from this package to the Go module root, so the
// test works from `go test ./...` and from inside the package alike.
func apiRootForCodeScan(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Skip("no go.mod above this package, so there is no project to scan")
	return ""
}
