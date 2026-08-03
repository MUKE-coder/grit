package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The button slot is only worth having if `grit swap button <variant>` actually
// restyles the admin. That holds when call sites route their classes through
// buttonClasses() instead of hand-writing bg-accent — and it breaks silently,
// because the templates are Go raw strings that compile no matter what TSX they
// contain.
//
// Two failure modes have already happened once each and are what these tests
// exist to catch:
//
//  1. A page uses buttonClasses() with no import. The Go build is happy; the
//     admin build fails on "buttonClasses is not defined".
//  2. The import lands in the wrong emitted file. Multi-function template files
//     hold several TSX sources, so an anchor-based insert can put the import in
//     a sibling page — which then carries an unused import while the real call
//     site has none.

// adminTSX scaffolds one admin variant and returns every emitted .tsx/.ts file.
func adminTSX(t *testing.T, frontend Frontend) map[string]string {
	t.Helper()

	root := t.TempDir()
	opts := Options{ProjectName: "slot-test", Architecture: ArchTriple, Frontend: frontend}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}

	write := writeAdminFiles
	if frontend == FrontendTanStack {
		write = writeAdminTanStackFiles
	}
	if err := write(root, opts); err != nil {
		t.Fatalf("write admin files: %v", err)
	}

	out := map[string]string{}
	adminDir := filepath.Join(root, "apps", "admin")
	err := filepath.Walk(adminDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || (!strings.HasSuffix(path, ".tsx") && !strings.HasSuffix(path, ".ts")) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(adminDir, path)
		out[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", adminDir, err)
	}
	if len(out) == 0 {
		t.Fatalf("no .tsx files emitted for frontend %q", frontend)
	}
	return out
}

// isSlot reports whether a file is the button slot itself, which defines
// buttonClasses and must not import it.
func isSlot(path string) bool {
	return strings.HasSuffix(path, "components/ui/button.tsx")
}

func TestButtonClassesImportsMatchUsage(t *testing.T) {
	for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
		t.Run(string(frontend), func(t *testing.T) {
			for path, src := range adminTSX(t, frontend) {
				if isSlot(path) {
					continue
				}

				imported := false
				for _, line := range strings.Split(src, "\n") {
					if strings.Contains(line, `from "@/components/ui/button"`) &&
						strings.Contains(line, "buttonClasses") {
						imported = true
						break
					}
				}

				// Count calls, not mentions — the import line mentions the name
				// too, so a substring check would call every import a use.
				used := strings.Contains(src, "buttonClasses(")

				switch {
				case used && !imported:
					t.Errorf("%s calls buttonClasses() but does not import it — the admin build will fail", path)
				case imported && !used:
					t.Errorf("%s imports buttonClasses but never calls it — the import probably belongs in a sibling file", path)
				}
			}
		})
	}
}

// A primary action that hand-writes bg-accent is invisible to `grit swap`. This
// pins the count so the number can only go down: adding a new page with an
// inline accent button now fails here rather than quietly eroding the slot.
//
// The remaining sites are ones a class-string swap cannot express — buttons
// whose className is a runtime expression rather than a literal, and the ones
// mixing accent tints (bg-accent/10) that are not primary actions at all.
func TestInlineAccentButtonsDoNotGrow(t *testing.T) {
	// 13 at the time of writing. The headroom is for a genuinely bespoke
	// control, not for a page that skipped the slot out of habit.
	const budget = 18

	for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
		t.Run(string(frontend), func(t *testing.T) {
			count := 0
			worst := ""
			worstN := 0

			for path, src := range adminTSX(t, frontend) {
				if isSlot(path) {
					continue
				}
				n := 0
				for _, tag := range inlineButtonTags(src) {
					if strings.Contains(tag, `bg-accent"`) || strings.Contains(tag, "bg-accent ") {
						n++
					}
				}
				count += n
				if n > worstN {
					worstN, worst = n, path
				}
			}

			if count > budget {
				t.Errorf("inline bg-accent <button> sites: %d, budget %d (worst: %s with %d).\n"+
					"New primary actions should call buttonClasses() so `grit swap button` reaches them.",
					count, budget, worst, worstN)
			}
			t.Logf("inline bg-accent <button> sites: %d (budget %d)", count, budget)
		})
	}
}

// inlineButtonTags returns each <button ...> opening tag in src. It tracks
// braces and quotes because `onClick={() => x}` contains a > that a naive scan
// mistakes for the end of the tag.
func inlineButtonTags(src string) []string {
	var tags []string
	for i := 0; ; {
		start := strings.Index(src[i:], "<button")
		if start < 0 {
			return tags
		}
		start += i

		depth := 0
		var quote byte
		j := start + len("<button")
		for ; j < len(src); j++ {
			ch := src[j]
			switch {
			case quote != 0:
				if ch == quote && src[j-1] != '\\' {
					quote = 0
				}
			case ch == '"' || ch == '\'' || ch == '`':
				quote = ch
			case ch == '{':
				depth++
			case ch == '}':
				depth--
			case ch == '>' && depth == 0:
				tags = append(tags, src[start:j+1])
			}
			if ch == '>' && depth == 0 && quote == 0 {
				break
			}
		}
		i = j + 1
	}
}
