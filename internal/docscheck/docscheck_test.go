package docscheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// A page shaped like the real ones: a multi-line block, a single-line prop, prose
// that starts with the word grit, an entity-encoded spec, a continuation, and a
// block that claims to be runnable.
const page = "" +
	"import { CodeBlock } from '@/components/code-block'\n" +
	"export default function Page() {\n" +
	"  return (\n" +
	"    <>\n" +
	"      <p>grit can keep itself current, which is prose and not a command.</p>\n" +
	"      <h3>grit sync: Manual Type Generation</h3>\n" +
	"      <CodeBlock terminal code=\"grit migrate --fresh\" />\n" +
	"      <CodeBlock terminal code={`grit generate resource Post --fields \"title:string,body:text\"\n" +
	"grit migrate\n" +
	"grit generate resource Invoice \\\n" +
	"  --fields \"number:string\"`} />\n" +
	"      <p>Written as entities: grit generate resource Tag --fields &quot;name:string&quot;</p>\n" +
	"      <CodeBlock\n" +
	"        terminal\n" +
	"        verify=\"rollback\"\n" +
	"        code={`grit migrate status\n" +
	"grit migrate down --yes`}\n" +
	"      />\n" +
	"    </>\n" +
	"  )\n" +
	"}\n"

func writePage(t *testing.T, name, body string) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}
	return false
}

func commandsOf(samples []Sample) []string {
	out := make([]string, 0, len(samples))
	for _, sample := range samples {
		out = append(out, sample.Command)
	}
	return out
}

func TestExtractReadsCodeBlocksAndNotProse(t *testing.T) {
	root := writePage(t, "page.tsx", page)
	samples, err := Extract(root)
	if err != nil {
		t.Fatal(err)
	}
	got := commandsOf(samples)

	for _, want := range []string{
		`grit migrate --fresh`,
		`grit generate resource Post --fields "title:string,body:text"`,
		`grit migrate`,
		// A trailing backslash joins the next line, which is how a long spec is
		// shown.
		`grit generate resource Invoice --fields "number:string"`,
		`grit migrate status`,
		`grit migrate down --yes`,
	} {
		if !contains(got, want) {
			t.Errorf("missing %q from %v", want, got)
		}
	}

	// Prose is not a command, however it starts. Both of these were reported by
	// the first version of the extractor.
	for _, unwanted := range []string{
		"grit can keep itself current, which is prose and not a command.",
		"grit sync: Manual Type Generation",
	} {
		if contains(got, unwanted) {
			t.Errorf("prose was read as a command: %q", unwanted)
		}
	}
}

func TestExtractDecodesEntities(t *testing.T) {
	root := writePage(t, "page.tsx",
		"<CodeBlock code={`grit generate resource Tag --fields &quot;name:string&quot;`} />\n")
	samples, err := Extract(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 1 {
		t.Fatalf("expected one command, got %v", commandsOf(samples))
	}
	// Without decoding, the spec's type reads as `string&quot;` and the parser
	// rejects a page that is perfectly correct.
	if want := `grit generate resource Tag --fields "name:string"`; samples[0].Command != want {
		t.Errorf("got %q, want %q", samples[0].Command, want)
	}
}

func TestExtractKeepsTheVerifyMarker(t *testing.T) {
	root := writePage(t, "page.tsx", page)
	samples, err := Extract(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, sample := range samples {
		marked := sample.Command == "grit migrate status" || sample.Command == "grit migrate down --yes"
		if marked && sample.Verify != "rollback" {
			t.Errorf("%q should belong to the rollback flow, got %q", sample.Command, sample.Verify)
		}
		if !marked && sample.Verify != "" {
			t.Errorf("%q is not in a marked block but claims flow %q", sample.Command, sample.Verify)
		}
	}
}

// The command tree stands in for grit's: two commands, a subcommand, an alias, a
// bool flag and a flag that takes a value.
func testTree() *cobra.Command {
	root := &cobra.Command{Use: "grit"}

	migrate := &cobra.Command{Use: "migrate"}
	migrate.Flags().Bool("fresh", false, "")
	down := &cobra.Command{Use: "down"}
	down.Flags().Int("steps", 1, "")
	down.Flags().Bool("yes", false, "")
	migrate.AddCommand(down)

	generate := &cobra.Command{Use: "generate", Aliases: []string{"g"}}
	resource := &cobra.Command{Use: "resource", Args: cobra.ExactArgs(1)}
	resource.Flags().String("fields", "", "")
	generate.AddCommand(resource)

	root.AddCommand(migrate, generate)
	return root
}

func TestCheckFindsWhatTheCLIWouldRefuse(t *testing.T) {
	cases := []struct {
		command string
		problem string // a substring of the expected message; empty means valid
	}{
		{"grit migrate", ""},
		{"grit migrate --fresh", ""},
		{"grit migrate down --steps 3 --yes", ""},
		{"grit g resource Post --fields \"title:string\"", ""},
		{"grit --help", ""}, // cobra adds it at run time, not at build time
		{"grit migrate:fresh", "no such command"},
		{"grit migrate --refresh", "has no --refresh flag"},
		{"grit migrate down --step 2", "has no --step flag"},
		{"grit nonsense", "no such command"},
		{`grit generate resource Post --fields "title:string`, "unclosed"},
	}

	for _, c := range cases {
		samples := []Sample{{File: "docs/x/page.tsx", Line: 1, Command: c.command}}
		problems := Check(samples, testTree())
		switch {
		case c.problem == "" && len(problems) > 0:
			t.Errorf("%q is valid but was reported: %s", c.command, problems[0].Message)
		case c.problem != "" && len(problems) == 0:
			t.Errorf("%q should have been reported as %q", c.command, c.problem)
		case c.problem != "" && !strings.Contains(problems[0].Message, c.problem):
			t.Errorf("%q reported as %q, want something about %q", c.command, problems[0].Message, c.problem)
		}
	}
}

// The changelog describes the past on purpose, flags and all, so checking it
// against today's CLI would mean editing history.
func TestCheckLeavesTheChangelogAlone(t *testing.T) {
	samples := []Sample{{File: "docs/changelog/page.tsx", Line: 9, Command: "grit upgrade --files"}}
	if problems := Check(samples, testTree()); len(problems) != 0 {
		t.Errorf("the changelog should not be checked: %s", problems[0])
	}
}

func TestCheckFieldSpecsUsesTheRealParser(t *testing.T) {
	valid := []Sample{{Command: `grit generate resource Invoice --fields "total:money,status:select:draft=Draft|paid=Paid"`}}
	if problems := CheckFieldSpecs(valid); len(problems) != 0 {
		t.Errorf("a valid spec was reported: %s", problems[0].Message)
	}

	// A select with no options is the bug this found in the docs: every command is
	// spelled correctly and the generator refuses it.
	invalid := []Sample{{Command: `grit generate resource Invoice --fields "status:select"`}}
	if problems := CheckFieldSpecs(invalid); len(problems) != 1 {
		t.Error("a select with no options should be reported")
	}

	unknownType := []Sample{{Command: `grit generate resource Invoice --fields "total:dollars"`}}
	if problems := CheckFieldSpecs(unknownType); len(problems) != 1 {
		t.Error("a field type that does not exist should be reported")
	}
}

func TestFlowsCaptureTheBlockVerbatim(t *testing.T) {
	root := writePage(t, "page.tsx", page)
	flows, err := Flows(root)
	if err != nil {
		t.Fatal(err)
	}
	blocks, ok := flows["rollback"]
	if !ok || len(blocks) != 1 {
		t.Fatalf("expected one block in the rollback flow, got %v", flows)
	}
	want := "grit migrate status\ngrit migrate down --yes"
	if blocks[0].Script != want {
		t.Errorf("script is %q, want %q", blocks[0].Script, want)
	}

	script := Script("rollback", blocks)
	if !strings.Contains(script, "set -euo pipefail") {
		t.Error("the generated script must stop at the first failure")
	}
	if !strings.Contains(script, "page.tsx") {
		t.Error("the script should name the page each block came from")
	}
}

// The testing page shows what a marked block looks like, which puts the text
// verify="migrate-rollback" inside a code sample. Read as a marker, the example
// would join the flow it describes and CI would run that walkthrough twice in one
// project.
func TestFlowsIgnoreAMarkerInsideACodeBlock(t *testing.T) {
	root := writePage(t, "page.tsx", ""+
		"      <p>Mark a block like this:</p>\n"+
		"      <CodeBlock\n"+
		"        language=\"tsx\"\n"+
		"        code={`<CodeBlock\n"+
		"  verify=\"rollback\"\n"+
		"  code={\\`grit migrate status\\`}\n"+
		"/>`}\n"+
		"      />\n")

	flows, err := Flows(root)
	if err != nil {
		t.Fatalf("an example of a marker should not be an error: %v", err)
	}
	if len(flows) != 0 {
		t.Errorf("a marker inside a code sample claimed a flow: %v", flows)
	}

	samples, err := Extract(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, sample := range samples {
		if sample.Verify != "" {
			t.Errorf("%q was put in flow %q by an example marker", sample.Command, sample.Verify)
		}
	}
}

// A marker with nothing runnable after it is an error, not a silent skip: a block
// that believes it is covered and is not is worse than one that never claimed to
// be.
func TestFlowsRefuseAMarkerWithNoBlock(t *testing.T) {
	root := writePage(t, "page.tsx", "<CodeBlock verify=\"nothing\" terminal />\n<p>prose</p>\n")
	if _, err := Flows(root); err == nil {
		t.Error("a verify marker with no code block should be an error")
	}
}
