package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func orderStates() []string {
	return []string{"draft", "submitted", "approved", "shipped", "cancelled"}
}

func orderWorkflow() *WorkflowSpec {
	return &WorkflowSpec{
		Initial:  "draft",
		Terminal: []string{"shipped", "cancelled"},
		Transitions: []TransitionSpec{
			{Action: "submit", From: []string{"draft"}, To: "submitted"},
			{Action: "approve", From: []string{"submitted"}, To: "approved", Permission: "orders.approve"},
			{Action: "ship", From: []string{"approved"}, To: "shipped"},
			{Action: "cancel", From: []string{"draft", "submitted", "approved"}, To: "cancelled", Confirm: true},
		},
	}
}

func TestWorkflowValidateAcceptsAGoodMachine(t *testing.T) {
	w := orderWorkflow()
	if err := w.Validate("Order", "status", orderStates()); err != nil {
		t.Fatalf("a valid workflow was rejected: %v", err)
	}
	// Labels are derived so a definition does not have to repeat itself.
	if w.Transitions[1].Label != "Approve" {
		t.Errorf("label not derived: %q", w.Transitions[1].Label)
	}
}

// The mistake this whole check exists for. A state with no way out is
// invisible until a record lands there in production and cannot be moved.
func TestWorkflowRejectsAStateNothingCanLeave(t *testing.T) {
	w := &WorkflowSpec{
		Initial: "draft",
		Transitions: []TransitionSpec{
			{Action: "submit", From: []string{"draft"}, To: "submitted"},
		},
		// submitted has no exit and is not declared terminal.
	}
	err := w.Validate("Order", "status", []string{"draft", "submitted"})
	if err == nil {
		t.Fatal("a state with no exit should be rejected")
	}
	if !strings.Contains(err.Error(), "submitted") {
		t.Errorf("the error should name the stuck state: %v", err)
	}
	// And it should say how to fix it, since both fixes are legitimate.
	if !strings.Contains(err.Error(), "terminal") {
		t.Errorf("the error should mention the terminal escape hatch: %v", err)
	}

	// Declaring it terminal is the other fix, and it must work.
	w.Terminal = []string{"submitted"}
	if err := w.Validate("Order", "status", []string{"draft", "submitted"}); err != nil {
		t.Fatalf("declaring it terminal should satisfy the check: %v", err)
	}
}

func TestWorkflowRejectsStatesTheFieldDoesNotHave(t *testing.T) {
	for _, tc := range []struct {
		name string
		spec *WorkflowSpec
		want string
	}{
		{
			"transition to an unknown state",
			&WorkflowSpec{Transitions: []TransitionSpec{{Action: "send", From: []string{"draft"}, To: "psoted"}}},
			"psoted",
		},
		{
			"transition from an unknown state",
			&WorkflowSpec{Transitions: []TransitionSpec{{Action: "send", From: []string{"drfat"}, To: "sent"}}},
			"drfat",
		},
		{
			"initial state not in the options",
			&WorkflowSpec{Initial: "nope", Transitions: []TransitionSpec{{Action: "send", From: []string{"draft"}, To: "sent"}}},
			"nope",
		},
		{
			"terminal state not in the options",
			&WorkflowSpec{
				Transitions: []TransitionSpec{{Action: "send", From: []string{"draft"}, To: "sent"}},
				Terminal:    []string{"snet"},
			},
			"snet",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.Validate("Invoice", "status", []string{"draft", "sent"})
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error should name %q: %v", tc.want, err)
			}
		})
	}
}

func TestWorkflowRejectsDuplicateActions(t *testing.T) {
	w := &WorkflowSpec{
		Transitions: []TransitionSpec{
			{Action: "send", From: []string{"draft"}, To: "sent"},
			{Action: "send", From: []string{"sent"}, To: "draft"},
		},
	}
	err := w.Validate("Invoice", "status", []string{"draft", "sent"})
	if err == nil || !strings.Contains(err.Error(), "twice") {
		t.Fatalf("a duplicate action should be rejected, got %v", err)
	}
}

func TestWorkflowRejectsAnUnknownBadgeColour(t *testing.T) {
	w := orderWorkflow()
	w.Colours = map[string]string{"draft": "chartreuse"}
	err := w.Validate("Order", "status", orderStates())
	if err == nil || !strings.Contains(err.Error(), "chartreuse") {
		t.Fatalf("an unknown colour should be rejected, got %v", err)
	}
}

// The generated per-resource file lives in package workflow, where
// workflow.Definition does not resolve. Getting this wrong produced a file
// that would not compile.
func TestWorkflowGoLiteralQualifier(t *testing.T) {
	w := orderWorkflow()
	if err := w.Validate("Order", "status", orderStates()); err != nil {
		t.Fatal(err)
	}

	inside := w.GoLiteral("orders", "status", orderStates(), "")
	if strings.Contains(inside, "workflow.Definition") {
		t.Errorf("a file inside package workflow must not qualify its own types:\n%s", inside)
	}
	if !strings.Contains(inside, "Definition{") || !strings.Contains(inside, "[]State{") {
		t.Errorf("unqualified literal is malformed:\n%s", inside)
	}

	outside := w.GoLiteral("orders", "status", orderStates(), "workflow.")
	if !strings.Contains(outside, "workflow.Definition{") || !strings.Contains(outside, "[]workflow.State{") {
		t.Errorf("qualified literal is malformed:\n%s", outside)
	}
}

func TestWorkflowGoLiteralCarriesEverything(t *testing.T) {
	w := orderWorkflow()
	if err := w.Validate("Order", "status", orderStates()); err != nil {
		t.Fatal(err)
	}
	got := w.GoLiteral("orders", "status", orderStates(), "")

	for _, want := range []string{
		`Resource: "orders"`,
		`Field: "status"`,
		`Initial: "draft"`,
		`Permission: "orders.approve"`,
		`Confirm: true`,
		`Terminal: true`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("literal missing %s:\n%s", want, got)
		}
	}
	// Colours are guessed when not declared, so a badge is never blank.
	if !strings.Contains(got, `Colour: "green"`) {
		t.Errorf("shipped should default to a green badge:\n%s", got)
	}
	if !strings.Contains(got, `Colour: "red"`) {
		t.Errorf("cancelled should default to a red badge:\n%s", got)
	}
}

func TestWorkflowLoadsFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "order.yaml")
	body := `name: Order
fields:
  - name: reference
    type: string
  - name: status
    type: select
    options:
      - value: draft
        label: Draft
      - value: sent
        label: Sent
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
	field := def.WorkflowField()
	if field == nil {
		t.Fatal("the workflow field was not found")
	}
	if field.Name != "status" {
		t.Errorf("wrong field: %q", field.Name)
	}
	if len(field.Workflow.Transitions) != 1 {
		t.Errorf("transitions: %+v", field.Workflow.Transitions)
	}
}

func TestLoadFromYAMLRejectsABrokenWorkflow(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "order.yaml")
	body := `name: Order
fields:
  - name: status
    type: select
    options:
      - value: draft
      - value: sent
    workflow:
      transitions:
        - action: send
          from: [draft]
          to: snet
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFromYAML(path)
	if err == nil {
		t.Fatal("a transition to a state the field does not have should fail to load")
	}
	if !strings.Contains(err.Error(), "snet") {
		t.Errorf("the error should name the typo: %v", err)
	}
}

func TestResourceWithoutAWorkflowIsUnaffected(t *testing.T) {
	def := &ResourceDefinition{
		Name: "Post",
		Fields: []Field{
			{Name: "title", Type: "string"},
			{Name: "status", Type: "select", Options: []FieldOption{{Value: "draft"}, {Value: "live"}}},
		},
	}
	if def.WorkflowField() != nil {
		t.Fatal("a plain select must not be treated as a state machine")
	}
	var nilSpec *WorkflowSpec
	if got := nilSpec.GoLiteral("posts", "status", []string{"draft"}, ""); got != "" {
		t.Errorf("a nil workflow should emit nothing, got %q", got)
	}
	if got := nilSpec.TSLiteral("status", []string{"draft"}); got != "" {
		t.Errorf("a nil workflow should emit no TS, got %q", got)
	}
}
