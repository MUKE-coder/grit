package scaffold

import (
	"strings"
	"testing"
)

// grit#211: ResourceCustomisation could patch columns, fields, components and
// bulkActions, and not table.actions.
//
// So the one override a project genuinely needs, which buttons a resource
// offers, was the one that had to be edited in the generated file and was
// therefore lost on every `grit generate resource --force`. The reporter lost
// the same edit three times.
//
// It matters because the framework cannot know a project's own rules. A
// resource whose rows the server authors, refused to clients by a create
// callback of the project's own, still generates a "New" button that can only
// ever 403, and a button that always fails teaches an operator the admin is
// broken.

func TestTheCustomFileCanReplaceTheTableActions(t *testing.T) {
	src := adminResourceTypes()

	if !strings.Contains(src, "actions?: TableAction[];") {
		t.Error("ResourceCustomisation has no actions, so taking a button away still means " +
			"editing the generated file and losing it to the next regeneration")
	}
	if !strings.Contains(src, `actions: custom?.actions ?? config.table.actions ?? ["create", "view", "edit", "delete"]`) {
		t.Error("defineResource does not read the custom actions, or reads them in the wrong " +
			"order: the custom file has to win over the generated config")
	}
}

// Replaced, not merged, unlike columns and fields. The point of setting it is
// almost always to take something away, and a patch that can only add would
// not do that.
func TestTheActionsOverrideReplacesRatherThanMerges(t *testing.T) {
	src := adminResourceTypes()
	// The merging shapes the other two use, which this must not be.
	for _, wrong := range []string{
		"...config.table.actions, ...custom",
		"custom?.actions ? [...config.table.actions",
		"concat(custom?.actions",
	} {
		if strings.Contains(src, wrong) {
			t.Errorf("the actions override merges (%q), so a project cannot remove a button", wrong)
		}
	}
}

// The three tiers, in the order defineResource resolves them. Asserted on the
// emitted text because the resolution is one line and the thing that breaks is
// its order.
func TestTheActionsPrecedenceIsCustomThenConfigThenEverything(t *testing.T) {
	src := adminResourceTypes()
	line := ""
	for _, l := range strings.Split(src, "\n") {
		if strings.Contains(l, "actions:") && strings.Contains(l, "custom?.actions") {
			line = l
		}
	}
	if line == "" {
		t.Fatal("defineResource does not resolve actions from the custom file at all")
	}
	custom := strings.Index(line, "custom?.actions")
	config := strings.Index(line, "config.table.actions")
	fallback := strings.Index(line, `"create", "view", "edit", "delete"`)
	if custom > config || config > fallback {
		t.Errorf("the precedence is wrong: the custom file must win over the generated config, "+
			"which must win over the default.\n  %s", strings.TrimSpace(line))
	}
}

// everything gated on an action follows from this list, so the override
// reaches the detail page's Edit and the bulk toolbar's delete without any of
// them being told separately.
func TestCanIsStillDerivedFromTheResolvedActions(t *testing.T) {
	// The controller hook, not the types file: this is where the override has
	// to arrive for the buttons to follow it.
	src := adminUseResourceController()
	if !strings.Contains(src, `const actions = resource.table.actions ?? ["create", "view", "edit", "delete"]`) {
		t.Error("the controller no longer reads actions off the resolved resource, so a custom " +
			"override would reach the list's buttons and not can()")
	}
	if !strings.Contains(src, "const can = useCallback((action: TableAction) => actions.includes(action)") {
		t.Error("can() is no longer derived from the actions list, so each surface would have " +
			"to be told about the override separately and each one would forget differently")
	}
}

// The stub is where people look. An extension point documented only in the
// docs site is one they will not find while reading the file it belongs in.
func TestTheCustomStubDocumentsTheActionsOverride(t *testing.T) {
	stub := AdminResourceCustomStub("AuditEntry", "AuditEntry")

	if !strings.Contains(stub, `actions: ["view"]`) {
		t.Error("the .custom stub does not show how to take a button away")
	}
	// And says why, because the reason is the part that is not obvious: the
	// generated actions are right for the framework's rules and wrong for the
	// project's.
	if !strings.Contains(stub, "always fails") {
		t.Error("the stub shows the override without saying what it is for")
	}
	// Still a valid stub: the example is prose, not a declaration, so a new
	// resource does not start out with its buttons removed.
	if strings.Contains(stub, "= {\n  actions:") {
		t.Error("the stub ships with actions actually set, which would remove buttons from " +
			"every newly generated resource")
	}
}
