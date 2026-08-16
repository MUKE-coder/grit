package generate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// WorkflowSpec is the `workflow:` block on a select or enum field.
//
//   - name: status
//     type: select
//     options: [draft, sent, paid, void]
//     workflow:
//     initial: draft
//     transitions:
//   - action: send
//     from: [draft]
//     to: sent
//     permission: invoices.send
//   - action: mark_paid
//     from: [sent]
//     to: paid
//   - action: void
//     from: [draft, sent]
//     to: void
//     confirm: true
//     terminal: [paid, void]
//
// States come from the field's own options rather than being repeated here.
// Two lists of the same thing drift, and the drift is silent: a transition to
// a state the dropdown does not offer.
type WorkflowSpec struct {
	Initial     string            `yaml:"initial"`
	Transitions []TransitionSpec  `yaml:"transitions"`
	Terminal    []string          `yaml:"terminal"`
	Colours     map[string]string `yaml:"colours"`
	Labels      map[string]string `yaml:"labels"`
}

// TransitionSpec is one legal move.
type TransitionSpec struct {
	Action     string   `yaml:"action"`
	Label      string   `yaml:"label"`
	From       []string `yaml:"from"`
	To         string   `yaml:"to"`
	Permission string   `yaml:"permission"`
	Confirm    bool     `yaml:"confirm"`
}

var validBadgeColours = []string{"gray", "blue", "green", "amber", "red", "purple"}

// Validate checks a workflow against the field that carries it.
//
// Run at parse time so a state machine that cannot be satisfied is a CLI error
// naming the resource, rather than a panic at boot or, worse, a record that
// lands somewhere nothing can move it out of.
func (w *WorkflowSpec) Validate(resource, field string, options []string) error {
	if w == nil {
		return nil
	}
	if len(options) == 0 {
		return fmt.Errorf("workflow on %s.%s needs the field to declare its options", resource, field)
	}

	known := map[string]bool{}
	for _, o := range options {
		known[o] = true
	}

	if w.Initial == "" {
		w.Initial = options[0]
	}
	if !known[w.Initial] {
		return fmt.Errorf("workflow on %s.%s starts in %q, which is not one of its options (%s)",
			resource, field, w.Initial, strings.Join(options, ", "))
	}

	if len(w.Transitions) == 0 {
		return fmt.Errorf("workflow on %s.%s declares no transitions, which makes it a plain select", resource, field)
	}

	seen := map[string]bool{}
	for i := range w.Transitions {
		t := &w.Transitions[i]
		if t.Action == "" {
			return fmt.Errorf("workflow on %s.%s has a transition with no action", resource, field)
		}
		if seen[t.Action] {
			return fmt.Errorf("workflow on %s.%s declares %q twice", resource, field, t.Action)
		}
		seen[t.Action] = true

		if t.To == "" {
			return fmt.Errorf("transition %q on %s.%s has no target state", t.Action, resource, field)
		}
		if !known[t.To] {
			return fmt.Errorf("transition %q on %s.%s moves to %q, which is not one of its options (%s)",
				t.Action, resource, field, t.To, strings.Join(options, ", "))
		}
		for _, from := range t.From {
			if !known[from] {
				return fmt.Errorf("transition %q on %s.%s comes from %q, which is not one of its options (%s)",
					t.Action, resource, field, from, strings.Join(options, ", "))
			}
		}
		if t.Label == "" {
			t.Label = humanLabel(t.Action)
		}
	}

	for _, s := range w.Terminal {
		if !known[s] {
			return fmt.Errorf("workflow on %s.%s marks %q terminal, which is not one of its options", resource, field, s)
		}
	}
	for name, colour := range w.Colours {
		if !known[name] {
			return fmt.Errorf("workflow on %s.%s colours %q, which is not one of its options", resource, field, name)
		}
		if !contains(validBadgeColours, colour) {
			return fmt.Errorf("workflow on %s.%s gives %q the colour %q (want one of %s)",
				resource, field, name, colour, strings.Join(validBadgeColours, ", "))
		}
	}

	// A state with no way out that is not declared terminal is the mistake
	// this whole check exists for: it is invisible until a record lands there
	// and cannot be moved, and by then it is in production.
	terminal := map[string]bool{}
	for _, s := range w.Terminal {
		terminal[s] = true
	}
	var stuck []string
	for _, o := range options {
		if terminal[o] {
			continue
		}
		hasExit := false
		for _, t := range w.Transitions {
			if len(t.From) == 0 || contains(t.From, o) {
				hasExit = true
				break
			}
		}
		if !hasExit {
			stuck = append(stuck, o)
		}
	}
	if len(stuck) > 0 {
		sort.Strings(stuck)
		return fmt.Errorf(
			"workflow on %s.%s: nothing can leave %s. Add a transition out, or list %s under terminal:",
			resource, field, strings.Join(stuck, " or "), strings.Join(stuck, " and "))
	}

	return nil
}

// GoLiteral renders the workflow.Definition the generator writes into
// internal/workflow/<resource>.go.
// qualifier is "workflow." for a file outside the package and "" for one
// inside it. The generated per-resource definitions live in package workflow
// itself, where workflow.Definition is undefined.
func (w *WorkflowSpec) GoLiteral(resource, field string, options []string, qualifier string) string {
	if w == nil {
		return ""
	}

	terminal := map[string]bool{}
	for _, s := range w.Terminal {
		terminal[s] = true
	}

	var states []string
	for _, o := range options {
		parts := []string{
			"Name: " + strconv.Quote(o),
			"Label: " + strconv.Quote(labelFor(w.Labels, o)),
			"Colour: " + strconv.Quote(colourFor(w.Colours, o, terminal[o])),
		}
		if terminal[o] {
			parts = append(parts, "Terminal: true")
		}
		states = append(states, "\t\t{"+strings.Join(parts, ", ")+"},")
	}

	var transitions []string
	for _, t := range w.Transitions {
		parts := []string{
			"Action: " + strconv.Quote(t.Action),
			"Label: " + strconv.Quote(t.Label),
			"From: " + goStringSlice(t.From),
			"To: " + strconv.Quote(t.To),
		}
		if t.Permission != "" {
			parts = append(parts, "Permission: "+strconv.Quote(t.Permission))
		}
		if t.Confirm {
			parts = append(parts, "Confirm: true")
		}
		transitions = append(transitions, "\t\t{"+strings.Join(parts, ", ")+"},")
	}

	return qualifier + "Definition{\n" +
		"\tResource: " + strconv.Quote(resource) + ",\n" +
		"\tField: " + strconv.Quote(field) + ",\n" +
		"\tInitial: " + strconv.Quote(w.Initial) + ",\n" +
		"\tStates: []" + qualifier + "State{\n" + strings.Join(states, "\n") + "\n\t},\n" +
		"\tTransitions: []" + qualifier + "Transition{\n" + strings.Join(transitions, "\n") + "\n\t},\n" +
		"}"
}

// TSLiteral renders the same definition for the admin resource file, so the
// detail page can render badges and the transitions legal right now without a
// round trip to ask.
func (w *WorkflowSpec) TSLiteral(field string, options []string) string {
	if w == nil {
		return ""
	}
	terminal := map[string]bool{}
	for _, s := range w.Terminal {
		terminal[s] = true
	}

	var states []string
	for _, o := range options {
		entry := "{ name: " + strconv.Quote(o) +
			", label: " + strconv.Quote(labelFor(w.Labels, o)) +
			", colour: " + strconv.Quote(colourFor(w.Colours, o, terminal[o]))
		if terminal[o] {
			entry += ", terminal: true"
		}
		states = append(states, "      "+entry+" },")
	}

	var transitions []string
	for _, t := range w.Transitions {
		entry := "{ action: " + strconv.Quote(t.Action) +
			", label: " + strconv.Quote(t.Label) +
			", from: " + tsStringArray(t.From) +
			", to: " + strconv.Quote(t.To)
		if t.Permission != "" {
			entry += ", permission: " + strconv.Quote(t.Permission)
		}
		if t.Confirm {
			entry += ", confirm: true"
		}
		transitions = append(transitions, "      "+entry+" },")
	}

	return "{\n" +
		"    field: " + strconv.Quote(field) + ",\n" +
		"    initial: " + strconv.Quote(w.Initial) + ",\n" +
		"    states: [\n" + strings.Join(states, "\n") + "\n    ],\n" +
		"    transitions: [\n" + strings.Join(transitions, "\n") + "\n    ],\n" +
		"  }"
}

func labelFor(labels map[string]string, state string) string {
	if l, ok := labels[state]; ok && l != "" {
		return l
	}
	return humanLabel(state)
}

// colourFor picks a badge tone when the author has not.
//
// The guesses are conventional rather than clever: things that mean "done" go
// green, things that mean "stopped" go red, everything else is neutral. An
// author who disagrees sets colours: explicitly.
func colourFor(colours map[string]string, state string, terminal bool) string {
	if c, ok := colours[state]; ok && c != "" {
		return c
	}
	switch strings.ToLower(state) {
	case "draft", "new", "pending", "todo":
		return "gray"
	case "sent", "submitted", "in_progress", "processing", "review", "assigned":
		return "blue"
	case "paid", "approved", "completed", "resolved", "done", "delivered", "active":
		return "green"
	case "overdue", "waiting", "on_hold", "partially_paid":
		return "amber"
	case "void", "cancelled", "canceled", "rejected", "failed", "refunded":
		return "red"
	}
	if terminal {
		return "green"
	}
	return "gray"
}

func tsStringArray(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = strconv.Quote(v)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}
