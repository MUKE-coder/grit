package errorcodes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#208: a project had nowhere to declare an error code of its own, and the
// workaround that suggests itself compiles and answers 500 at runtime.

func writeProjectFile(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(ProjectFile))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestAProjectCodeIsLoaded(t *testing.T) {
	root := writeProjectFile(t, `
codes:
  - code: PLAN_LIMIT_EXCEEDED_HERE
    status: 402
    category: limit
    area: billing
    meaning: The account has used all three of its checks.
    client: Show the upgrade prompt. Retrying will not help.
`)
	codes, err := LoadProject(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 1 {
		t.Fatalf("loaded %d codes, want 1", len(codes))
	}
	got := codes[0]
	if got.Code != "PLAN_LIMIT_EXCEEDED_HERE" || got.Status != 402 ||
		got.Category != CategoryLimit || got.Area != "billing" {
		t.Errorf("got %+v", got)
	}
	if got.Meaning == "" || got.Client == "" {
		t.Error("meaning and client are what the docs table and the TypeScript comment show")
	}
}

// Most projects have no codes of their own, and that is not a failure.
func TestAMissingFileIsNotAnError(t *testing.T) {
	codes, err := LoadProject(t.TempDir())
	if err != nil {
		t.Fatalf("a project with no error-codes.yaml failed to load: %v", err)
	}
	if len(codes) != 0 {
		t.Errorf("got %d codes from nowhere", len(codes))
	}
}

// An empty list is the shape the starter file ships in, so it has to load.
func TestTheStarterFileLoadsAndDeclaresNothing(t *testing.T) {
	root := writeProjectFile(t, StarterProjectFile())
	codes, err := LoadProject(root)
	if err != nil {
		t.Fatalf("the starter file Grit writes does not load: %v", err)
	}
	if len(codes) != 0 {
		t.Errorf("the starter file declares %d codes; it is meant to be an example, commented out",
			len(codes))
	}
	// And the example in it is a comment, not a declaration: a new project
	// must not start out answering a code it has no handler for.
	if !strings.Contains(StarterProjectFile(), "#   - code:") {
		t.Error("the example in the starter file is not commented out")
	}
}

// Every refusal, because each one is a failure that would otherwise be silent
// or confusing later, and that is the bar for rejecting somebody's file.
func TestTheRefusals(t *testing.T) {
	cases := []struct {
		name, body, wants string
	}{
		{
			name: "shadowing one of Grit's codes",
			body: "codes:\n  - code: VALIDATION_ERROR\n    status: 400\n    category: request\n    meaning: m\n    client: c\n",
			// Two entries for one code means its status depends on which wins,
			// which is the exact inconsistency the catalogue ended.
			wants: "already one of Grit's codes",
		},
		{
			name:  "a status that is not an error",
			body:  "codes:\n  - code: MY_CODE\n    status: 200\n    category: request\n    meaning: m\n    client: c\n",
			wants: "not an HTTP error status",
		},
		{
			name:  "a status that is not a status",
			body:  "codes:\n  - code: MY_CODE\n    status: 799\n    category: request\n    meaning: m\n    client: c\n",
			wants: "not an HTTP error status",
		},
		{
			name: "a category that does not exist",
			body: "codes:\n  - code: MY_CODE\n    status: 400\n    category: billing\n    meaning: m\n    client: c\n",
			// Named, because guessing from a list of ten is worse than reading it.
			wants: "The categories are:",
		},
		{
			name:  "a code name no constant could be derived from",
			body:  "codes:\n  - code: plan-limit\n    status: 400\n    category: request\n    meaning: m\n    client: c\n",
			wants: "SHOUTING_SNAKE_CASE",
		},
		{
			name:  "no code at all",
			body:  "codes:\n  - status: 400\n    category: request\n    meaning: m\n    client: c\n",
			wants: "has no code",
		},
		{
			name:  "the same code twice",
			body:  "codes:\n  - code: MY_CODE\n    status: 400\n    category: request\n    meaning: m\n    client: c\n  - code: MY_CODE\n    status: 409\n    category: conflict\n    meaning: m\n    client: c\n",
			wants: "declared twice",
		},
		{
			name:  "no meaning",
			body:  "codes:\n  - code: MY_CODE\n    status: 400\n    category: request\n    client: c\n",
			wants: "has no meaning",
		},
		{
			name:  "no client advice",
			body:  "codes:\n  - code: MY_CODE\n    status: 400\n    category: request\n    meaning: m\n",
			wants: "has no client",
		},
		{
			name:  "not YAML at all",
			body:  "codes: [this is: not valid\n",
			wants: "not valid YAML",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := LoadProject(writeProjectFile(t, c.body))
			if err == nil {
				t.Fatalf("accepted: %s", c.body)
			}
			if !strings.Contains(err.Error(), c.wants) {
				t.Errorf("the message does not say what is wrong.\n got: %v\nwant it to contain: %q",
					err, c.wants)
			}
			// Every message names the file, because the person is looking at a
			// terminal and has to know which file to open.
			if !strings.Contains(err.Error(), ProjectFile) {
				t.Errorf("the message does not name the file: %v", err)
			}
		})
	}
}

// The area is the only optional field, and it has to default to something the
// docs table can group by.
func TestAreaDefaultsToProject(t *testing.T) {
	root := writeProjectFile(t, "codes:\n  - code: MY_CODE\n    status: 400\n    category: request\n    meaning: m\n    client: c\n")
	codes, err := LoadProject(root)
	if err != nil {
		t.Fatal(err)
	}
	if codes[0].Area != "project" {
		t.Errorf("area = %q, want project", codes[0].Area)
	}
}

// Sorted, so the generated files do not churn on every run: two codes in a
// different order would rewrite codes.go and errors.ts with no change in them.
func TestProjectCodesComeBackSorted(t *testing.T) {
	root := writeProjectFile(t, "codes:\n  - code: ZEBRA\n    status: 400\n    category: request\n    meaning: m\n    client: c\n  - code: AARDVARK\n    status: 400\n    category: request\n    meaning: m\n    client: c\n")
	codes, err := LoadProject(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 2 || codes[0].Code != "AARDVARK" {
		t.Errorf("not sorted: %+v", codes)
	}
}

// The new code the issue asked to ship, and the thing about it that matters:
// it is not 429. A client with backoff told to retry against a plan limit
// retries against a wall, and the person is told the wrong thing.
func TestPlanLimitReachedIsPaymentRequiredAndNotRateLimited(t *testing.T) {
	entry, ok := Lookup("PLAN_LIMIT_REACHED")
	if !ok {
		t.Fatal("PLAN_LIMIT_REACHED is not in the catalogue")
	}
	if entry.Status != 402 {
		t.Errorf("status = %d, want 402: 429 means wait, and waiting does not raise a plan limit",
			entry.Status)
	}
	if entry.Status == StatusOf("RATE_LIMITED") {
		t.Error("it shares a status with RATE_LIMITED, so a client cannot tell them apart by status")
	}
	if entry.Category != CategoryLimit {
		t.Errorf("category = %q, want limit", entry.Category)
	}
}
