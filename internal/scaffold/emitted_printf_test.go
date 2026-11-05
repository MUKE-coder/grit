package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// printfLikeEmittedHelpers are the functions Grit writes that take a format
// string and its arguments. Passing one a pre-formatted string is what vet's
// printf check calls a non-constant format string.
var printfLikeEmittedHelpers = []string{"respond.Rule", "respond.Refuse"}

// No emitted code may hand a printf-like helper a fmt.Sprintf.
//
// `go test` runs a subset of vet, printf included, and a non-constant format
// string is an error rather than a warning there. So one generated line
//
//	respond.Rule(fmt.Sprintf("%s must be %s", typeErr.Field, ...))
//
// made `go test ./internal/services/` fail to build in every project with a
// generated resource. The service compiled; only its tests did not, which is
// why it survived: `go build ./...` does not run vet, and the generator's own
// tests parse the emitted Go rather than vetting it.
//
// Found by running the test suite of a project upgraded from v3.474.0, not by
// reading the template. A sweep, so the next printf-like helper that acquires a
// caller is caught where this one was not.
func TestEmittedCodeNeverSprintfsIntoAPrintfHelper(t *testing.T) {
	patterns := make([]*regexp.Regexp, 0, len(printfLikeEmittedHelpers))
	for _, helper := range printfLikeEmittedHelpers {
		patterns = append(patterns, regexp.MustCompile(regexp.QuoteMeta(helper)+`\(\s*fmt\.Sprintf\(`))
	}

	checked := 0
	check := func(path string, data []byte) {
		checked++
		for i, p := range patterns {
			for _, loc := range p.FindAllIndex(data, -1) {
				before := string(data[:loc[0]])
				line := 1 + strings.Count(before, "\n")
				// A repair holds the broken text so it can recognise it. That
				// is a record of the fault, not an emission of it.
				start := strings.LastIndex(before, "\n") + 1
				if strings.Contains(before[start:], "Old = ") {
					continue
				}
				t.Errorf("%s:%d: emitted code calls %s with a fmt.Sprintf.\n\n"+
					"    %s takes a format string and its arguments, so vet reads the\n"+
					"    Sprintf's result as a non-constant format string and `go test`\n"+
					"    refuses to build the package. Pass the format and the arguments\n"+
					"    straight through instead of formatting them first.",
					filepath.ToSlash(path), line,
					printfLikeEmittedHelpers[i], printfLikeEmittedHelpers[i])
			}
		}
	}

	// The template source in both writers, and the .tmpl files beside it.
	for _, root := range []string{".", filepath.Join("..", "generate")} {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatalf("reading %s: %v", root, err)
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			path := filepath.Join(root, name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			check(path, data)
		}
	}
	err := filepath.WalkDir("templates", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		check(path, data)
		return nil
	})
	if err != nil {
		t.Fatalf("walking templates: %v", err)
	}

	if checked < 400 {
		t.Fatalf("only %d files were read; the sweep is not seeing the templates", checked)
	}
	t.Logf("%d template files checked", checked)
}

// The repair and the generator have to agree on the fixed line.
//
// Two separate assertions, one per side, would both pass while the pair was
// broken: the generator could emit a third form and the repair would go on
// producing the second. So this reads both.
func TestThePrintfRepairProducesWhatTheGeneratorWrites(t *testing.T) {
	if strings.Contains(printfRuleNew, "fmt.Sprintf") {
		t.Fatal("the repair's replacement still formats the message first")
	}
	old := "func (s *ProductService) Patch() {\n\t" + printfRuleOld + "\n}\n"
	out, changed := repairServicePrintfRule(old)
	if !changed {
		t.Fatal("the repair did not recognise the line it exists to fix")
	}
	if !strings.Contains(out, printfRuleNew) {
		t.Errorf("the repair produced:\n%s\nwant it to contain:\n%s", out, printfRuleNew)
	}
	if again, ok := repairServicePrintfRule(out); ok || again != out {
		t.Error("the repair is not idempotent, so every upgrade reports it again")
	}
}
