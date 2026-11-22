package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// `grit sync --check`: what sync would write, compared with what is there.
//
// grit sync regenerates packages/shared from the Go models, and the whole
// end-to-end type story rests on somebody having run it. Nothing checked that
// anybody had. A developer who added a field to a model and forgot the command
// got a green build, a green type-check and a client whose type did not have
// the field, which surfaces as undefined at runtime in whichever screen reads
// it first (grit#255).
//
// It writes nothing, deliberately. A check that repairs what it finds cannot
// be run on a pull request, which is the only place it is worth running: the
// point is to fail, in CI, before the mismatch is merged.
//
// Scope is the two files sync derives from a model, the TypeScript type and
// the Zod schema, and only where one is already there.
//
// Not "every model implies two files", which is the literal reading and was
// measured first: on a freshly scaffolded triple project it reports 78
// missing, because packages/shared does not mirror every framework model and
// was never meant to. Nothing on a client validates an OAuthExchangeCode. A
// checker that fails on a project's first day is one nobody runs, and the
// useful question is the one the issue describes: the files that are there,
// are they what their model now says. Every generated resource has both, so
// a resource is covered either way.
//
// Not the admin resource, the desktop screen or the Expo form either, because
// sync only adds to those and leaves edits alone, so "what sync would write"
// is not a single answer for them. And not orphans: a .ts matching no model
// is usually one somebody wrote.

// DriftKind says how a file disagrees with its model.
type DriftKind string

const (
	// DriftStale: the file is there and is not what the model implies.
	DriftStale DriftKind = "stale"
)

// SharedDrift is one file in packages/shared that does not match its model.
type SharedDrift struct {
	// Rel is the path relative to the project root, which is what a person
	// reading CI output can act on.
	Rel   string
	Model string
	Kind  DriftKind
	// Want and Got are the whole files. The first differing line is what
	// gets printed; the rest is here so a caller can diff properly.
	Want string
	Got  string
}

// FirstDifference returns the 1-indexed line where Want and Got part company,
// with both versions of it.
func (d SharedDrift) FirstDifference() (line int, want, got string) {
	wantLines := strings.Split(strings.ReplaceAll(d.Want, "\r\n", "\n"), "\n")
	gotLines := strings.Split(strings.ReplaceAll(d.Got, "\r\n", "\n"), "\n")
	for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
		w, g := "", ""
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			return i + 1, w, g
		}
	}
	return 0, "", ""
}

// SyncCheck reports every shared file that does not match its Go model.
//
// An empty slice and a nil error is the pass.
func SyncCheck() ([]SharedDrift, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}
	return syncCheckIn(root)
}

// syncCheckIn is SyncCheck against a known root, so a test can point it at a
// fixture without changing the working directory.
func syncCheckIn(root string) ([]SharedDrift, error) {
	modelsDir := filepath.Join(apiRootFrom(root), "internal", "models")
	if !dirExists(modelsDir) {
		return nil, fmt.Errorf("models directory not found at %s", modelsDir)
	}
	sharedRoot := sharedRootFrom(root)
	if !dirExists(sharedRoot) {
		return nil, fmt.Errorf("shared package not found at %s: sync requires a monorepo project", sharedRoot)
	}

	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("reading models directory: %w", err)
	}

	var drift []SharedDrift
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		structs, parseErr := parseGoStructs(filepath.Join(modelsDir, entry.Name()))
		if parseErr != nil {
			// Sync skips a model it cannot parse with a warning rather than
			// failing, and a check that failed where sync shrugs would be
			// reporting a different project from the one sync writes.
			continue
		}
		for _, s := range structs {
			if skipSharedType(s.Name) {
				continue
			}
			for _, want := range []struct {
				path    string
				content string
			}{
				{sharedTypePath(sharedRoot, s.Name), buildTSType(s)},
				{sharedSchemaPath(sharedRoot, s.Name), buildZodSchema(s)},
			} {
				raw, readErr := os.ReadFile(want.path)
				if readErr != nil {
					// Not drift. See the scope note at the top: a model
					// packages/shared does not mirror is a decision, and
					// sync creating the file is not the same as the project
					// being wrong without it.
					continue
				}
				rel := relFromRoot(root, want.path)
				// Normalised, because git checks these out with CRLF on
				// Windows and a byte comparison would then report every file
				// as stale on every Windows machine.
				got := strings.ReplaceAll(string(raw), "\r\n", "\n")
				expected := strings.ReplaceAll(want.content, "\r\n", "\n")
				if got != expected {
					drift = append(drift, SharedDrift{
						Rel: rel, Model: s.Name, Kind: DriftStale,
						Want: expected, Got: got,
					})
				}
			}
		}
	}
	return drift, nil
}

// PrintSyncCheck writes the result and reports whether it passed.
func PrintSyncCheck(drift []SharedDrift) bool {
	if len(drift) == 0 {
		fmt.Print("\n  ✓ packages/shared matches the Go models.\n\n")
		return true
	}

	fmt.Printf("\n  %d file%s in packages/shared %s not what the models imply:\n\n",
		len(drift), plural(len(drift)), isAre(len(drift)))
	for _, d := range drift {
		line, want, got := d.FirstDifference()
		fmt.Printf("  ✗ %s\n", d.Rel)
		fmt.Printf("      line %d, from model %s\n", line, d.Model)
		fmt.Printf("        on disk:   %s\n", strings.TrimSpace(got))
		fmt.Printf("        the model: %s\n\n", strings.TrimSpace(want))
	}
	fmt.Print("  Run grit sync and commit what it writes.\n\n")
	return false
}

// isAre, so the sentence reads.
func isAre(n int) string {
	if n == 1 {
		return "is"
	}
	return "are"
}
