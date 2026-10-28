package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// injectBefore is idempotent against the lines directly above its own marker.
// That is the right question for a snippet which may legitimately repeat, and
// the wrong one for a short variable declaration: `x := ...` can exist once in
// a scope, so a second one does not compile however the lines around the
// marker have moved.
//
// The Kaya Goods storefront had moved the Stripe block up and extended it, so
// the marker check did not recognise it and `grit upgrade` injected a second
// copy fourteen lines below the first:
//
//	internal/routes/routes.go:904:17: no new variables on left side of :=
func TestInjectRefusesADeclarationTheFileAlreadyHas(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "routes.go")

	// The project's own copy, moved and extended, so the lines above the
	// marker are not the injection.
	existing := `package routes

func Register(db *gorm.DB) {
	paymentService := payments.NewService(db)
	paymentService.OnSucceeded = shopService.MarkPaid

	// grit:handlers
}
`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	code := "\tpaymentService := payments.NewService(db)\n\tpaymentService.Register()"
	err := injectBefore(path, "// grit:handlers", code)
	if err == nil {
		t.Fatal("injected a second paymentService declaration, which cannot compile")
	}

	var declared *AlreadyDeclaredError
	if !asAlreadyDeclared(err, &declared) {
		t.Fatalf("wanted an AlreadyDeclaredError, got %T: %v", err, err)
	}
	if declared.Name != "paymentService" {
		t.Errorf("named %q, wanted paymentService", declared.Name)
	}

	after, _ := os.ReadFile(path)
	if string(after) != existing {
		t.Error("the file was modified despite the refusal")
	}
	if strings.Count(string(after), "paymentService :=") != 1 {
		t.Error("more than one declaration survived")
	}
}

// And it still injects when the file does not declare that name, which is the
// ordinary case this guard must not break.
func TestInjectStillAddsADeclarationTheFileLacks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "routes.go")
	if err := os.WriteFile(path, []byte("package routes\n\nfunc Register() {\n\t// grit:handlers\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := injectBefore(path, "// grit:handlers", "\tpaymentService := payments.NewService(db)"); err != nil {
		t.Fatalf("refused an injection into a file that does not have it: %v", err)
	}
	after, _ := os.ReadFile(path)
	if !strings.Contains(string(after), "paymentService := payments.NewService(db)") {
		t.Error("the injection did not land")
	}
}

// A block that declares nothing, such as an import or a struct entry, is not
// this guard's business and goes in as before.
func TestInjectIgnoresBlocksThatDeclareNothing(t *testing.T) {
	if name, ok := declaresVariable("\t\"app/internal/payments\""); ok {
		t.Errorf("read %q as a declaration", name)
	}
	if name, ok := declaresVariable("\t\t&Payment{},\n\t\t&Subscription{},"); ok {
		t.Errorf("read %q as a declaration", name)
	}
	// Two different names is a shape this does not judge.
	if _, ok := declaresVariable("\ta := f()\n\tb := g()"); ok {
		t.Error("judged a block declaring two names")
	}
}

// asAlreadyDeclared keeps errors.As out of the test's import list, which the
// repository's older files do not carry.
func asAlreadyDeclared(err error, target **AlreadyDeclaredError) bool {
	if e, ok := err.(*AlreadyDeclaredError); ok {
		*target = e
		return true
	}
	return false
}
