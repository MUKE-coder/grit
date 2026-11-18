package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The repair moves an existing project's handlers onto respond.Invalid.
//
// It needs to exist at all because handlers are the developer's after
// scaffolding and upgrade preserves them, which is right. Converting the 53
// templates therefore reached new projects and nobody else: every project
// already out there went on answering a failed validation with the validator's
// own sentence and an empty error.details (grit#225).
func TestTheBindFailureRepairConvertsAnOlderProject(t *testing.T) {
	root := t.TempDir()
	apiRoot := filepath.Join(root, "apps", "api")

	// The gate: respond.Invalid has to be there before anything calls it.
	writeBindFile(t, filepath.Join(apiRoot, "go.mod"), "module example.com/app\n")
	writeBindFile(t, filepath.Join(apiRoot, "internal", "respond", "respond.go"),
		"package respond\n\nfunc Invalid(c *gin.Context, err error) {}\n")

	handler := `package handlers

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
}
`
	path := filepath.Join(apiRoot, "internal", "handlers", "auth.go")
	writeBindFile(t, path, handler)

	if err := RepairBindFailures(root); err != nil {
		t.Fatalf("RepairBindFailures: %v", err)
	}

	out := readBindFile(t, path)
	if strings.Contains(out, "respond.Fail(c, respond.CodeValidationError, err.Error())") {
		t.Error("a handler still answers with the validator's own sentence")
	}
	if n := strings.Count(out, "respond.Invalid(c, err)"); n != 2 {
		t.Errorf("converted %d of 2 call sites", n)
	}

	// Twice over the same project changes nothing, because grit upgrade runs
	// every repair on every run.
	if err := RepairBindFailures(root); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if again := readBindFile(t, path); again != out {
		t.Error("a second upgrade changed the file again, so the repair is not idempotent")
	}
}

// With no respond.Invalid in the project, it does nothing.
//
// The gate is on the function being there rather than on a version, because
// the thing that matters is whether the call it writes would compile. A repair
// that runs before its helper arrives leaves a project that does not build,
// which is worse than the defect.
func TestTheBindFailureRepairWaitsForItsHelper(t *testing.T) {
	root := t.TempDir()
	apiRoot := filepath.Join(root, "apps", "api")

	writeBindFile(t, filepath.Join(apiRoot, "go.mod"), "module example.com/app\n")
	writeBindFile(t, filepath.Join(apiRoot, "internal", "respond", "respond.go"),
		"package respond\n\n// no Invalid here yet\n")

	handler := "package handlers\n\nfunc f() {\n\trespond.Fail(c, respond.CodeValidationError, err.Error())\n}\n"
	path := filepath.Join(apiRoot, "internal", "handlers", "auth.go")
	writeBindFile(t, path, handler)

	if err := RepairBindFailures(root); err != nil {
		t.Fatalf("RepairBindFailures: %v", err)
	}
	if readBindFile(t, path) != handler {
		t.Error("the repair wrote a call to respond.Invalid into a project that " +
			"does not have it, so the project no longer compiles")
	}
}

// It leaves internal/respond alone.
//
// Invalid's own doc comment quotes the old line while explaining why not to
// write it, which is the same shape as a guard matching the prose that
// explains it: a blind replace would rewrite the explanation.
func TestTheBindFailureRepairLeavesRespondsOwnProseAlone(t *testing.T) {
	root := t.TempDir()
	apiRoot := filepath.Join(root, "apps", "api")

	writeBindFile(t, filepath.Join(apiRoot, "go.mod"), "module example.com/app\n")
	respondSrc := `package respond

// Invalid answers a binding failure.
//
// Fifty-four handlers answered with
// respond.Fail(c, respond.CodeValidationError, err.Error()), and err.Error()
// on a validator failure is the library talking to itself.
func Invalid(c *gin.Context, err error) {}
`
	path := filepath.Join(apiRoot, "internal", "respond", "respond.go")
	writeBindFile(t, path, respondSrc)

	if err := RepairBindFailures(root); err != nil {
		t.Fatalf("RepairBindFailures: %v", err)
	}
	if readBindFile(t, path) != respondSrc {
		t.Error("the repair rewrote respond's own explanation of the line it replaces")
	}
}

func writeBindFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readBindFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}
