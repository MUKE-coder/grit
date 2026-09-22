package plugin

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripePluginIsRegistered(t *testing.T) {
	p, err := Get("stripe")
	if err != nil {
		t.Fatalf("the stripe plugin is not registered: %v", err)
	}
	for _, d := range p.NodeDeps {
		if d.Workspace != "apps/web" {
			t.Errorf("%s goes to %s, want apps/web", d.Name, d.Workspace)
		}
	}
	if len(p.GoDeps) != 0 {
		t.Errorf("GoDeps = %+v: payments talks to Stripe's REST API directly", p.GoDeps)
	}
}

// The API half always installs; the payment form only into a Next.js web app,
// the one whose CSP the plugin knows how to extend.
func TestStripeFilesFollowTheProject(t *testing.T) {
	root := t.TempDir()
	ctx := Context{Root: root, Module: "shop/apps/api", Architecture: "triple", Frontend: "next"}
	files := stripeFiles(ctx)
	for _, want := range []string{"apps/api/internal/models/payment.go", "apps/api/internal/payments/stripe.go", "apps/api/internal/payments/service.go", "apps/api/internal/payments/payments_test.go", "apps/api/internal/handlers/payment.go"} {
		if _, ok := files[want]; !ok {
			t.Errorf("missing %s", want)
		}
	}
	for path := range files {
		if strings.HasPrefix(path, "apps/web/") {
			t.Errorf("%s was written to a project without a web app", path)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "apps", "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, ok := stripeFiles(ctx)["apps/web/components/stripe-checkout.tsx"]; !ok {
		t.Error("no payment form once the web app exists")
	}
	ctx.Frontend = "tanstack"
	if _, ok := stripeFiles(ctx)["apps/web/components/stripe-checkout.tsx"]; ok {
		t.Error("the payment form was written to a Vite app, whose CSP it cannot extend")
	}
}

func TestStripeTemplatesAreFormattedGo(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "apps", "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := Context{Root: root, Module: "shop/apps/api", Architecture: "triple", Frontend: "next"}
	for path, src := range stripeFiles(ctx) {
		if strings.Contains(src, "{{MODULE}}") {
			t.Errorf("%s still has {{MODULE}}", path)
		}
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		formatted, err := format.Source([]byte(src))
		if err != nil {
			t.Errorf("%s does not parse: %v", path, err)
			continue
		}
		if string(formatted) != src {
			t.Errorf("%s is not gofmt-clean", path)
		}
	}
}

// Stripe.js needs its origins in the CSP, written at the marker the scaffold
// puts in next.config.ts; an older config is reported, not a failed install.
func TestStripeAddsItsOriginsToTheCSP(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "apps", "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := Context{Root: root, Module: "shop/apps/api", Architecture: "triple", Frontend: "next"}
	for _, inj := range stripeInjections(ctx) {
		if inj.File != "apps/web/next.config.ts" {
			continue
		}
		for _, want := range []string{`["script-src", "https://js.stripe.com"]`, `["frame-src", "https://js.stripe.com"]`, `["connect-src", "https://api.stripe.com"]`} {
			if !strings.Contains(inj.Code, want) {
				t.Errorf("the CSP injection lacks %s", want)
			}
		}
		if inj.Marker != "// grit:csp-origins" || !inj.Optional {
			t.Errorf("CSP injection = %+v", inj)
		}
		return
	}
	t.Error("no CSP injection")
}

// Two plugins inject their imports at the same marker, and the second one
// lands after the first however it sorts. The file is formatted on the way
// out, so a project with both installed is still gofmt-clean.
func TestInjectingLeavesGoFileFormatted(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "routes.go")
	source := "package routes\n\nimport (\n\t\"shop/apps/api/internal/models\"\n\t\"shop/apps/api/internal/webhooks\"\n\t// Imports added by plugins.\n\t// grit:imports\n)\n"
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"\t\"shop/apps/api/internal/video\"", "\t\"shop/apps/api/internal/payments\""} {
		if err := injectBefore(path, "// grit:imports", code); err != nil {
			t.Fatal(err)
		}
	}
	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	formatted, err := format.Source(written)
	if err != nil {
		t.Fatalf("the injected file does not parse: %v", err)
	}
	if string(formatted) != string(written) {
		t.Errorf("not gofmt-clean after two injections:\n%s", written)
	}
	for _, want := range []string{"internal/video", "internal/payments"} {
		if !strings.Contains(string(written), want) {
			t.Errorf("%s was lost", want)
		}
	}
}
