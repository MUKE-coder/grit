package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func purchaseRequestGenerator(t *testing.T) (*Generator, string) {
	t.Helper()
	const module = "erp/apps/api"
	root := setupMinimalProject(t, module)
	path := filepath.Join(t.TempDir(), "purchase_request.yaml")
	body := `name: PurchaseRequest
fields:
  - name: title
    type: string
  - name: status
    type: select
    options:
      - value: draft
      - value: pending
      - value: approved
      - value: rejected
    workflow:
      initial: draft
      terminal: [approved, rejected]
      transitions:
        - action: submit
          from: [draft]
          to: pending
        - action: approve
          from: [pending]
          to: approved
          permission: purchase_requests.approve
        - action: reject
          from: [pending]
          to: rejected
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return newTestGenerator(root, module, def), root
}

// Approving a purchase request has to draw down a budget and reserve stock,
// and refuse when it cannot. The transition was one UPDATE followed by an
// in-memory event, so that work could only go in a subscriber that ran after
// the approval had committed and could not refuse it.
func TestTransitionRunsHooksInItsTransaction(t *testing.T) {
	g, root := purchaseRequestGenerator(t)
	if err := g.writeWorkflowService(g.Names(), g.Definition.WorkflowField()); err != nil {
		t.Fatalf("workflow service: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "purchase_request_workflow.go"))
	mustParse(t, "purchase_request_workflow.go", src)

	tx := strings.Index(src, "db.Transaction(func(tx *gorm.DB) error {")
	update := strings.Index(src, "tx.Model(&models.PurchaseRequest{})")
	hooks := strings.Index(src, "workflow.RunHooks(tx, workflow.Move{")
	durable := strings.Index(src, "events.EmitTx(tx, c, &ev)")
	after := strings.Index(src, "events.Emit(c, ev)")
	if tx < 0 || update < tx || hooks < update || durable < hooks || after < durable {
		t.Fatalf("the status update, hooks and durable event are not one transaction in that order:\n%s", src)
	}
	if strings.Contains(src, "db.Model(&models.PurchaseRequest{})") {
		t.Error("the status update still runs outside the transaction")
	}
	if !strings.Contains(src, "workflow.ErrNotPermitted{") {
		t.Error("a missing permission is not a typed error, so the handler cannot tell it from a 500")
	}
}

// Every error the handler did not recognise was answered as a 403, so a
// database failure read as "you are not allowed".
func TestTransitionHandlerClassifiesErrors(t *testing.T) {
	g, _ := purchaseRequestGenerator(t)
	src := g.workflowHandlerMethod(g.Names())
	for _, want := range []string{
		"status, code := workflow.Classify(err)",
		"case http.StatusInternalServerError:",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the transition handler is missing %s", want)
		}
	}
	if strings.Contains(src, `"code": "FORBIDDEN", "message": err.Error()`) {
		t.Error("the handler still answers unknown errors with a 403")
	}
}
