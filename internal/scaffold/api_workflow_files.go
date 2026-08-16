package scaffold

// Workflows: the difference between a status column and a process.
//
// A select field today accepts every value on every record. draft can jump
// straight to paid, paid can go back to draft, and a support agent can mark an
// invoice collected. Any rule about which of those is allowed lives in the
// author's head or in a hand-written check they have to remember to write in
// every place the field is touched.
//
// A workflow states the rules once, next to the field, and the guard runs in
// the service layer where every caller passes through it.

// APIWorkflowGo is exported so grit generate resource can add the package to
// a project generated before workflows existed.
func APIWorkflowGo() string { return apiWorkflowGo() }

// apiWorkflowGo emits internal/workflow/workflow.go.
func apiWorkflowGo() string {
	return `// Package workflow turns a status field into a state machine.
//
// A definition names the states, the transitions between them and who may
// perform each one. Generated services consult it before writing, so an
// invalid move is a 422 with a usable message rather than a successful write
// that leaves a record in a state the business rules say cannot exist.
package workflow

import (
	"fmt"
	"sort"
	"strings"
)

// Transition is one legal move.
type Transition struct {
	// Action is the verb: "send", "mark_paid", "void". It is what the API
	// takes as POST /invoices/:id/transitions/:action.
	Action string ` + "`" + `json:"action"` + "`" + `
	// Label is what a button says. Defaults to a title-cased Action.
	Label string ` + "`" + `json:"label"` + "`" + `
	// From lists the states this move is legal in. Empty means any state,
	// which is occasionally what you want (a "cancel" from anywhere) and
	// usually a mistake worth noticing.
	From []string ` + "`" + `json:"from"` + "`" + `
	// To is the resulting state.
	To string ` + "`" + `json:"to"` + "`" + `
	// Permission gates the move. Empty means any authenticated caller.
	Permission string ` + "`" + `json:"permission,omitempty"` + "`" + `
	// Confirm asks the user before firing. For the moves that are hard to
	// undo: void, refund, cancel.
	Confirm bool ` + "`" + `json:"confirm,omitempty"` + "`" + `
}

// State is one position in the machine, with how it should render.
type State struct {
	Name  string ` + "`" + `json:"name"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
	// Colour is a badge tone the admin understands: gray, blue, green,
	// amber, red, purple.
	Colour string ` + "`" + `json:"colour,omitempty"` + "`" + `
	// Terminal marks an end state. Nothing leaves it, and the admin can stop
	// showing a "what next" affordance.
	Terminal bool ` + "`" + `json:"terminal,omitempty"` + "`" + `
}

// Definition is one resource's state machine.
type Definition struct {
	// Resource is the plural snake_case name, matching the event bus.
	Resource string ` + "`" + `json:"resource"` + "`" + `
	// Field is the model column holding the state.
	Field string ` + "`" + `json:"field"` + "`" + `
	// Initial is the state a new record starts in.
	Initial     string       ` + "`" + `json:"initial"` + "`" + `
	States      []State      ` + "`" + `json:"states"` + "`" + `
	Transitions []Transition ` + "`" + `json:"transitions"` + "`" + `
}

// ErrInvalidTransition is returned when a move is not legal from the current
// state.
//
// It carries both ends rather than a bare "invalid", because the useful
// message is "an invoice that is already paid cannot be sent", not "error".
type ErrInvalidTransition struct {
	Resource string
	From     string
	Action   string
	Allowed  []string
}

func (e ErrInvalidTransition) Error() string {
	if len(e.Allowed) == 0 {
		return fmt.Sprintf("%s cannot be %s from %s: it is a terminal state",
			e.Resource, e.Action, e.From)
	}
	return fmt.Sprintf("%s cannot be %s from %s (allowed from here: %s)",
		e.Resource, e.Action, e.From, strings.Join(e.Allowed, ", "))
}

// ErrUnknownAction is returned for an action the workflow does not define.
type ErrUnknownAction struct {
	Resource string
	Action   string
	Known    []string
}

func (e ErrUnknownAction) Error() string {
	return fmt.Sprintf("%s has no transition %q (it has: %s)",
		e.Resource, e.Action, strings.Join(e.Known, ", "))
}

// Find returns a transition by action name.
func (d Definition) Find(action string) (Transition, bool) {
	for _, t := range d.Transitions {
		if t.Action == action {
			return t, true
		}
	}
	return Transition{}, false
}

// Can reports whether an action is legal from a state.
func (d Definition) Can(from, action string) bool {
	t, ok := d.Find(action)
	if !ok {
		return false
	}
	return allows(t, from)
}

// Check validates a move and returns the transition, or an error explaining
// why not. This is what a generated service calls.
func (d Definition) Check(from, action string) (Transition, error) {
	t, ok := d.Find(action)
	if !ok {
		known := make([]string, 0, len(d.Transitions))
		for _, x := range d.Transitions {
			known = append(known, x.Action)
		}
		sort.Strings(known)
		return Transition{}, ErrUnknownAction{Resource: d.Resource, Action: action, Known: known}
	}
	if !allows(t, from) {
		allowed := make([]string, 0, len(d.Transitions))
		for _, x := range d.Transitions {
			if allows(x, from) {
				allowed = append(allowed, x.Action)
			}
		}
		sort.Strings(allowed)
		return Transition{}, ErrInvalidTransition{
			Resource: d.Resource, From: from, Action: action, Allowed: allowed,
		}
	}
	return t, nil
}

// Next lists the transitions legal from a state.
//
// This is what the admin renders: buttons for the moves that are possible
// right now, rather than a dropdown of every state that exists. Showing a
// user an option that will be rejected is a worse interface than not showing
// it.
func (d Definition) Next(from string) []Transition {
	var out []Transition
	for _, t := range d.Transitions {
		if allows(t, from) {
			out = append(out, t)
		}
	}
	return out
}

// StateOf returns a state's metadata, or a sensible default for one the
// definition does not name. A record can hold an unknown state after a
// definition changes, and the admin still has to render it.
func (d Definition) StateOf(name string) State {
	for _, s := range d.States {
		if s.Name == name {
			return s
		}
	}
	return State{Name: name, Label: titleise(name), Colour: "gray"}
}

// Validate checks a definition for the mistakes that would otherwise surface
// as a record stuck in a state nothing can leave.
func (d Definition) Validate() error {
	if d.Field == "" {
		return fmt.Errorf("workflow on %s has no field", d.Resource)
	}
	if len(d.States) == 0 {
		return fmt.Errorf("workflow on %s has no states", d.Resource)
	}

	known := map[string]bool{}
	for _, s := range d.States {
		known[s.Name] = true
	}
	if d.Initial != "" && !known[d.Initial] {
		return fmt.Errorf("workflow on %s starts in %q, which is not one of its states", d.Resource, d.Initial)
	}

	for _, t := range d.Transitions {
		if !known[t.To] {
			return fmt.Errorf("transition %q on %s moves to %q, which is not one of its states",
				t.Action, d.Resource, t.To)
		}
		for _, from := range t.From {
			if !known[from] {
				return fmt.Errorf("transition %q on %s comes from %q, which is not one of its states",
					t.Action, d.Resource, from)
			}
		}
	}

	// A non-terminal state nothing can leave is almost always a mistake, and
	// it is invisible until a record lands there and gets stuck.
	for _, s := range d.States {
		if s.Terminal {
			continue
		}
		if len(d.Next(s.Name)) == 0 {
			return fmt.Errorf("state %q on %s has no transitions out and is not marked terminal",
				s.Name, d.Resource)
		}
	}
	return nil
}

func allows(t Transition, from string) bool {
	if len(t.From) == 0 {
		return true
	}
	for _, f := range t.From {
		if f == from {
			return true
		}
	}
	return false
}

func titleise(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' })
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

// ── Registry ─────────────────────────────────────────────────────────────
//
// Generated definitions register themselves in an init, so the admin can ask
// for a resource's workflow without every caller importing every definition.

var registry = map[string]Definition{}

// Register adds a definition. Panics on an invalid one, at boot, because a
// workflow that cannot be satisfied should stop the process rather than
// surface as a record nobody can move.
func Register(d Definition) {
	if err := d.Validate(); err != nil {
		panic(err)
	}
	registry[d.Resource] = d
}

// For returns a resource's workflow.
func For(resource string) (Definition, bool) {
	d, ok := registry[resource]
	return d, ok
}

// All returns every registered workflow, for GET /api/workflows.
func All() map[string]Definition {
	out := make(map[string]Definition, len(registry))
	for k, v := range registry {
		out[k] = v
	}
	return out
}
`
}
