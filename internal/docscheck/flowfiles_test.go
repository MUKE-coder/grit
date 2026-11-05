package docscheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFlow(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReadFlowFiles(t *testing.T) {
	dir := t.TempDir()
	writeFlow(t, dir, "two-steps.md", `---
page: /docs/cli
title: Two blocks, in order
---

Prose about why.

`+"```bash"+`
grit generate resource Widget --fields "name:string"
grit migrate
`+"```"+`

More prose.

`+"```go"+`
// Not a shell block, so not run.
func main() {}
`+"```"+`

`+"```bash"+`
grit remove resource Widget --force
`+"```"+`
`)
	// Skipped: the directory's own README.
	writeFlow(t, dir, "_README.md", "---\npage: /docs/testing\ntitle: x\n---\n\n```bash\nnot-run\n```\n")

	flows, err := ReadFlowFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(flows) != 1 {
		t.Fatalf("read %d flow(s), want 1 (an _-prefixed file is documentation)", len(flows))
	}
	flow := flows[0]
	if flow.ID != "two-steps" {
		t.Errorf("id is %q, want the file name", flow.ID)
	}
	if flow.Page != "/docs/cli" || flow.Title != "Two blocks, in order" {
		t.Errorf("frontmatter read as page=%q title=%q", flow.Page, flow.Title)
	}
	if len(flow.Blocks) != 2 {
		t.Fatalf("read %d block(s), want 2: a go block is documentation, not a step", len(flow.Blocks))
	}
	if !strings.Contains(flow.Blocks[0].Script, "grit migrate") {
		t.Errorf("first block is %q", flow.Blocks[0].Script)
	}
	if !strings.Contains(flow.Blocks[1].Script, "grit remove resource") {
		t.Errorf("second block is %q", flow.Blocks[1].Script)
	}
	for _, b := range flow.Blocks {
		if b.Flow != "two-steps" {
			t.Errorf("block carries flow %q", b.Flow)
		}
	}
}

// A missing directory is the state this started from, not a failure.
func TestReadFlowFilesToleratesNoDirectory(t *testing.T) {
	flows, err := ReadFlowFiles(filepath.Join(t.TempDir(), "nothing-here"))
	if err != nil {
		t.Fatalf("a missing directory is an error: %v", err)
	}
	if len(flows) != 0 {
		t.Errorf("read %d flows from nowhere", len(flows))
	}
}

func TestFlowFileRefusals(t *testing.T) {
	for name, tc := range map[string]struct{ body, want string }{
		"no frontmatter": {
			body: "```bash\ngrit migrate\n```\n",
			want: "opens with ---",
		},
		"frontmatter never closed": {
			body: "---\npage: /docs/cli\ntitle: x\n",
			want: "never closed",
		},
		"no page": {
			body: "---\ntitle: x\n---\n\n```bash\ngrit migrate\n```\n",
			want: "no page:",
		},
		"no title": {
			body: "---\npage: /docs/cli\n---\n\n```bash\ngrit migrate\n```\n",
			want: "no title:",
		},
		// The one outcome worse than no flow: a flow that passes having run
		// nothing, and is counted in the coverage.
		"no commands": {
			body: "---\npage: /docs/cli\ntitle: x\n---\n\nAll prose, no block.\n",
			want: "no ```bash block",
		},
		"block never closed": {
			body: "---\npage: /docs/cli\ntitle: x\n---\n\n```bash\ngrit migrate\n",
			want: "never closed",
		},
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			writeFlow(t, dir, "broken.md", tc.body)
			_, err := ReadFlowFiles(dir)
			if err == nil {
				t.Fatal("no error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error is %q, want it to mention %q", err, tc.want)
			}
		})
	}
}

// A page's verify= id and a flow file's name must not collide.
func TestAllFlowsRefusesADuplicateID(t *testing.T) {
	app := t.TempDir()
	page := filepath.Join(app, "docs", "cli")
	if err := os.MkdirAll(page, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFlow(t, page, "page.tsx", "<CodeBlock verify=\"clash\" code={`grit migrate`} />\n")

	flowDir := t.TempDir()
	writeFlow(t, flowDir, "clash.md", "---\npage: /docs/cli\ntitle: x\n---\n\n```bash\ngrit seed\n```\n")

	if _, err := AllFlows(app, flowDir); err == nil {
		t.Error("two sources declaring the same flow id were merged silently")
	} else if !strings.Contains(err.Error(), "declared twice") {
		t.Errorf("error is %q", err)
	}
}

// A flow whose page has moved still runs and still passes, and has stopped
// being evidence about anything a reader can find.
func TestFlowFilePagesFindsADanglingPointer(t *testing.T) {
	app := t.TempDir()
	if err := os.MkdirAll(filepath.Join(app, "docs", "cli"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFlow(t, filepath.Join(app, "docs", "cli"), "page.tsx", "export default function P() {}\n")

	flowDir := t.TempDir()
	writeFlow(t, flowDir, "here.md", "---\npage: /docs/cli\ntitle: x\n---\n\n```bash\ngrit seed\n```\n")
	writeFlow(t, flowDir, "gone.md", "---\npage: /docs/moved-away\ntitle: x\n---\n\n```bash\ngrit seed\n```\n")

	dangling, err := FlowFilePages(app, flowDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(dangling) != 1 {
		t.Fatalf("reported %d dangling pointer(s), want 1: %v", len(dangling), dangling)
	}
	if !strings.Contains(dangling[0], "/docs/moved-away") {
		t.Errorf("reported %q", dangling[0])
	}
}

// This repository's own flow files, read with the real reader. A broken one
// fails here rather than in CI twenty minutes into a Postgres job.
func TestThisRepositorysFlowFilesParse(t *testing.T) {
	flows, err := ReadFlowFiles(filepath.Join("..", "..", FlowFileDir))
	if err != nil {
		t.Fatal(err)
	}
	if len(flows) == 0 {
		t.Fatal("no flow files, and this repository has some")
	}
	seen := map[string]bool{}
	for _, flow := range flows {
		if seen[flow.ID] {
			t.Errorf("two flows called %q", flow.ID)
		}
		seen[flow.ID] = true
		for _, block := range flow.Blocks {
			if strings.TrimSpace(block.Script) == "" {
				t.Errorf("%s has an empty block", flow.Path)
			}
		}
	}

	dangling, err := FlowFilePages(filepath.Join("..", "..", "docs", "app"), filepath.Join("..", "..", FlowFileDir))
	if err != nil {
		t.Fatal(err)
	}
	for _, problem := range dangling {
		t.Error(problem)
	}
}
