package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorkflowOptions is what `grit generate workflow` was asked for.
type WorkflowOptions struct {
	// Name is the saga, in PascalCase: Checkout, Onboarding, RefundOrder.
	Name string
	// Steps are the step names in order, comma separated: "charge,reserve,ship".
	Steps string
}

// GenerateWorkflow writes a saga: an ordered list of steps, each with a Do and
// an Undo, and registers it.
//
// The bodies are stubs that return an error rather than nil. A stub that
// returned nil would be a saga that reports success without doing anything, and
// the first time you would find out is when a customer says the parcel never
// came. Failing loudly means the step is unmissable until it is written.
func GenerateWorkflow(opts WorkflowOptions) error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	arch, _ := readGritJSON(root)
	apiRoot := root
	if arch != "single" {
		apiRoot = filepath.Join(root, "apps", "api")
	}
	module, err := readModulePath(root, arch)
	if err != nil {
		return err
	}
	return generateWorkflowAt(apiRoot, module, opts)
}

func generateWorkflowAt(apiRoot, module string, opts WorkflowOptions) error {
	if strings.TrimSpace(opts.Name) == "" {
		return fmt.Errorf("a workflow name is required (e.g. grit generate workflow Checkout --steps \"charge,reserve,ship\")")
	}
	names := MakeNames(opts.Name)
	if names.Pascal == "" || !isGoIdentifier(names.Pascal) {
		return fmt.Errorf("%q is not a usable workflow name: use letters and digits, starting with a letter (e.g. Checkout)", opts.Name)
	}

	steps, err := parseWorkflowSteps(opts.Steps)
	if err != nil {
		return err
	}

	sagasDir := filepath.Join(apiRoot, "internal", "sagas")
	registry := filepath.Join(sagasDir, "sagas.go")
	if !fileHas(registry, "// grit:sagas") {
		return fmt.Errorf("internal/sagas has no registry to add this to, so the saga would be written and never run: run grit upgrade first")
	}
	if !fileHas(filepath.Join(apiRoot, "internal", "saga", "saga.go"), "func Register(") {
		return fmt.Errorf("internal/saga is missing, so there is no engine to run this: run grit upgrade first")
	}

	path := filepath.Join(sagasDir, names.Snake+".go")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("internal/sagas/%s.go already exists. It holds your saga, so it is not overwritten", names.Snake)
	}

	fmt.Printf("\n  Generating workflow: %s (%d step(s))\n\n", names.Pascal, len(steps))

	if err := writeFileWithDirs(path, workflowSource(module, names, steps)); err != nil {
		return fmt.Errorf("writing the saga: %w", err)
	}
	fmt.Printf("  ✓ internal/sagas/%s.go\n", names.Snake)

	if err := injectWorkflowRegistration(registry, module, names); err != nil {
		return err
	}
	fmt.Printf("  ✓ internal/sagas/sagas.go: %s registered\n", names.Pascal)

	fmt.Printf("\n  Each step returns an error until you write it, so a half-built saga\n")
	fmt.Printf("  fails loudly instead of reporting success and doing nothing.\n\n")
	fmt.Printf("  Start one with:\n\n")
	fmt.Printf("      saga.Start(ctx, db, %q, input, saga.Key(\"order:\"+order.ID))\n\n", names.Snake)
	fmt.Printf("  Read internal/saga's package comment first: a step has to be idempotent,\n")
	fmt.Printf("  and the steps you cannot undo go last.\n\n")
	return nil
}

// workflowStep is one parsed step.
type workflowStep struct {
	Names Names
	// Raw is what the user typed, for the message when it is not usable.
	Raw string
}

// parseWorkflowSteps reads "charge,reserve stock,ship" into step names.
func parseWorkflowSteps(spec string) ([]workflowStep, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, fmt.Errorf("a workflow needs steps: --steps \"charge,reserve,ship\"")
	}
	var out []workflowStep
	seen := map[string]bool{}
	for _, raw := range strings.Split(spec, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		names := MakeNames(raw)
		if names.Pascal == "" || !isGoIdentifier(names.Pascal) {
			return nil, fmt.Errorf("%q is not a usable step name: use letters, digits and underscores (e.g. charge, reserve_stock)", raw)
		}
		// Two steps with one name would share a row in saga_steps and a
		// compensation would undo whichever the query happened to return.
		if seen[names.Snake] {
			return nil, fmt.Errorf("two steps are both called %q, and the run table keys on the name", names.Snake)
		}
		seen[names.Snake] = true
		out = append(out, workflowStep{Names: names, Raw: raw})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("a workflow needs steps: --steps \"charge,reserve,ship\"")
	}
	return out, nil
}

// injectWorkflowRegistration adds the saga to sagas.Register.
func injectWorkflowRegistration(path, module string, names Names) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading the saga registry: %w", err)
	}
	src := string(data)
	crlf := strings.Contains(src, "\r\n")
	src = strings.ReplaceAll(src, "\r\n", "\n")

	const marker = "\t// grit:sagas\n"
	if !strings.Contains(src, marker) {
		return fmt.Errorf("internal/sagas/sagas.go has no // grit:sagas marker, so there is nowhere to register %s", names.Pascal)
	}
	line := fmt.Sprintf("\tsaga.Register(%s())\n", names.Pascal)
	if strings.Contains(src, line) {
		return nil
	}
	src = strings.Replace(src, marker, line+marker, 1)

	// The first saga in the project is the one that adds the import: the
	// registry is scaffolded with no import block, because an unused import
	// does not compile and an empty registry uses nothing.
	if !strings.Contains(src, "/internal/saga\"") {
		src = strings.Replace(src, "package sagas\n", "package sagas\n\nimport (\n\t\""+module+"/internal/saga\"\n)\n", 1)
	}

	if crlf {
		src = strings.ReplaceAll(src, "\n", "\r\n")
	}
	return os.WriteFile(path, []byte(src), 0o644)
}

func workflowSource(module string, names Names, steps []workflowStep) string {
	var b strings.Builder

	fmt.Fprintf(&b, "package sagas\n\n")
	fmt.Fprintf(&b, "import (\n\t\"context\"\n\t\"fmt\"\n\n\t\"%s/internal/saga\"\n)\n\n", module)

	fmt.Fprintf(&b, "// %s is a saga: %d step(s) that have to either all finish or be undone.\n", names.Pascal, len(steps))
	fmt.Fprintf(&b, "//\n")
	fmt.Fprintf(&b, "// Start one with:\n")
	fmt.Fprintf(&b, "//\n")
	fmt.Fprintf(&b, "//\tsaga.Start(ctx, db, %q, input, saga.Key(\"...\"))\n", names.Snake)
	fmt.Fprintf(&b, "//\n")
	fmt.Fprintf(&b, "// Three things to get right, and they are not optional:\n")
	fmt.Fprintf(&b, "//\n")
	fmt.Fprintf(&b, "//  1. Every Do must be idempotent. A crash between the side effect and the\n")
	fmt.Fprintf(&b, "//     record of it will happen, so a resumed run sometimes repeats a step.\n")
	fmt.Fprintf(&b, "//     Pass r.IdempotencyKey() to whatever you are calling.\n")
	fmt.Fprintf(&b, "//  2. Every Undo must be idempotent too, and must tolerate a Do that never\n")
	fmt.Fprintf(&b, "//     finished. Check before you reverse.\n")
	fmt.Fprintf(&b, "//  3. Steps you cannot undo go last. An email cannot be unsent, so a step\n")
	fmt.Fprintf(&b, "//     after it that fails leaves it sent.\n")
	fmt.Fprintf(&b, "func %s() saga.Definition {\n", names.Pascal)
	fmt.Fprintf(&b, "\treturn saga.Definition{\n")
	fmt.Fprintf(&b, "\t\tName: %q,\n", names.Snake)
	fmt.Fprintf(&b, "\t\tSteps: []saga.Step{\n")

	for i, step := range steps {
		fmt.Fprintf(&b, "\t\t\t{\n")
		fmt.Fprintf(&b, "\t\t\t\tName: %q,\n", step.Names.Snake)
		fmt.Fprintf(&b, "\t\t\t\tDo: func(ctx context.Context, r *saga.Run) error {\n")
		fmt.Fprintf(&b, "\t\t\t\t\t// TODO: %s.\n", step.Names.Snake)
		fmt.Fprintf(&b, "\t\t\t\t\t//\n")
		if i == 0 {
			fmt.Fprintf(&b, "\t\t\t\t\t// Read what the saga was started with:\n")
			fmt.Fprintf(&b, "\t\t\t\t\t//\n")
			fmt.Fprintf(&b, "\t\t\t\t\t//\tvar in %sInput\n", names.Pascal)
			fmt.Fprintf(&b, "\t\t\t\t\t//\tif err := r.Input(&in); err != nil { return err }\n")
			fmt.Fprintf(&b, "\t\t\t\t\t//\n")
		} else {
			fmt.Fprintf(&b, "\t\t\t\t\t// Read what an earlier step left, and leave what the next\n")
			fmt.Fprintf(&b, "\t\t\t\t\t// one and this step's Undo will need:\n")
			fmt.Fprintf(&b, "\t\t\t\t\t//\n")
			fmt.Fprintf(&b, "\t\t\t\t\t//\tid := r.GetString(\"charge_id\")\n")
			fmt.Fprintf(&b, "\t\t\t\t\t//\tr.Set(%q, result.ID)\n", step.Names.Snake+"_id")
			fmt.Fprintf(&b, "\t\t\t\t\t//\n")
		}
		fmt.Fprintf(&b, "\t\t\t\t\t// Pass r.IdempotencyKey() to whatever you call, so a retry of\n")
		fmt.Fprintf(&b, "\t\t\t\t\t// this step does not do it twice. Return saga.Fatal(err) for a\n")
		fmt.Fprintf(&b, "\t\t\t\t\t// failure that retrying will not fix.\n")
		fmt.Fprintf(&b, "\t\t\t\t\treturn fmt.Errorf(%q)\n", "saga "+names.Snake+": step "+step.Names.Snake+" is not implemented")
		fmt.Fprintf(&b, "\t\t\t\t},\n")
		fmt.Fprintf(&b, "\t\t\t\t// Undo takes this step back when a later one fails for good. Set\n")
		fmt.Fprintf(&b, "\t\t\t\t// it to nil if this step cannot be undone, and then move it to\n")
		fmt.Fprintf(&b, "\t\t\t\t// the end of the list, because everything before it will be\n")
		fmt.Fprintf(&b, "\t\t\t\t// compensated and this will not.\n")
		fmt.Fprintf(&b, "\t\t\t\tUndo: func(ctx context.Context, r *saga.Run) error {\n")
		fmt.Fprintf(&b, "\t\t\t\t\t// TODO: take back %s, if it happened. It may not have:\n", step.Names.Snake)
		fmt.Fprintf(&b, "\t\t\t\t\t// compensation runs because something went wrong, and that\n")
		fmt.Fprintf(&b, "\t\t\t\t\t// includes not knowing whether the Do landed. Check first.\n")
		fmt.Fprintf(&b, "\t\t\t\t\treturn nil\n")
		fmt.Fprintf(&b, "\t\t\t\t},\n")
		fmt.Fprintf(&b, "\t\t\t},\n")
	}

	fmt.Fprintf(&b, "\t\t},\n")
	fmt.Fprintf(&b, "\t}\n")
	fmt.Fprintf(&b, "}\n\n")

	fmt.Fprintf(&b, "// %sInput is what this saga is started with. It is stored as JSON on the\n", names.Pascal)
	fmt.Fprintf(&b, "// run and never changes, so put the identifiers here and read the records\n")
	fmt.Fprintf(&b, "// themselves inside the steps: a copy of a row taken at the start is stale by\n")
	fmt.Fprintf(&b, "// the time a compensation reads it an hour later.\n")
	fmt.Fprintf(&b, "type %sInput struct {\n", names.Pascal)
	fmt.Fprintf(&b, "\t// TODO: the identifiers this saga works from.\n")
	fmt.Fprintf(&b, "\tID string `json:\"id\"`\n")
	fmt.Fprintf(&b, "}\n")

	return b.String()
}
