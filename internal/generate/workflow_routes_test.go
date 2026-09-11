package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func workflowDef(t *testing.T) *ResourceDefinition {
	t.Helper()
	path := filepath.Join(t.TempDir(), "referral.yaml")
	body := `name: Referral
fields:
  - name: reason
    type: text
  - name: status
    type: select
    options:
      - value: draft
      - value: sent
    workflow:
      initial: draft
      terminal: [sent]
      transitions:
        - action: send
          from: [draft]
          to: sent
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return def
}

// A resource with a workflow did not compile on a current project: its
// transition routes were injected into routes.go against a referralHandler
// variable that only projects from before route files declare. They belong in
// the resource's own route file, which has the handler as h.
func TestWorkflowRoutesLiveInTheResourceRouteFile(t *testing.T) {
	const module = "clinic/apps/api"
	g := newTestGenerator(setupMinimalProject(t, module), module, workflowDef(t))
	names := g.Names()

	src, err := g.resourceRoutesSource(names)
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	for _, want := range []string{
		`m.Protected.GET("/referrals/workflow", h.Workflow)`,
		`m.Protected.POST("/referrals/:id/transitions/:action", h.Transition)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the route file is missing %s", want)
		}
	}

	g.Roles = []string{"DOCTOR"}
	src, err = g.resourceRoutesSource(names)
	if err != nil {
		t.Fatalf("routes with roles: %v", err)
	}
	if !strings.Contains(src, `g.POST("/:id/transitions/:action", h.Transition)`) {
		t.Error("a role-restricted resource's transitions are not behind its role guard")
	}
}

// A project already broken by the old injection is repaired by generating the
// resource again.
func TestWorkflowRoutesAreTakenBackOutOfRoutesGo(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, workflowDef(t))
	names := g.Names()

	routesDir := filepath.Join(g.APIRoot(), "internal", "routes")
	if err := os.MkdirAll(routesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	routesGo := filepath.Join(routesDir, "routes.go")
	broken := "package routes\n\nfunc mount() {\n" +
		"\t\tprotected.GET(\"/referrals/workflow\", referralHandler.Workflow)\n" +
		"\t\tprotected.POST(\"/referrals/:id/transitions/:action\", referralHandler.Transition)\n" +
		"\t\t// grit:routes:protected\n}\n"
	if err := os.WriteFile(routesGo, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(routesDir, "referral_routes.go"), []byte("package routes\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := g.ensureWorkflowRoutes(names); err != nil {
		t.Fatalf("ensureWorkflowRoutes: %v", err)
	}
	got := readTestFile(t, routesGo)
	if strings.Contains(got, "referralHandler") {
		t.Errorf("routes.go still refers to a handler it never declares:\n%s", got)
	}
	if !strings.Contains(got, "// grit:routes:protected") {
		t.Error("the marker went with the stale lines")
	}
}
