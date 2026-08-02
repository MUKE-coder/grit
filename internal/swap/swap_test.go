package swap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContractOf(t *testing.T) {
	cases := []struct{ src, want string }{
		{"/* grit:slot button@1\n * blah */", "button@1"},
		{"// grit:slot   input@2", "input@2"},
		{"/* grit:slot multi-word@11 */", "multi-word@11"},
		{"no marker here", ""},
		// A version is required — "grit:slot button" is not a contract.
		{"/* grit:slot button */", ""},
	}
	for _, c := range cases {
		if got := ContractOf(c.src); got != c.want {
			t.Errorf("ContractOf(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}

func TestMajor(t *testing.T) {
	for in, want := range map[string]string{
		"button@1": "1", "input@12": "12", "nope": "", "": "",
	} {
		if got := Major(in); got != want {
			t.Errorf("Major(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestStripPreview(t *testing.T) {
	src := `export const Button = 1

/* grit:preview-start */
export default function Preview() {
  return <div>demo</div>
}
/* grit:preview-end */
`
	got := StripPreview(src)
	if strings.Contains(got, "Preview") {
		t.Fatalf("preview survived the strip:\n%s", got)
	}
	if !strings.Contains(got, "export const Button") {
		t.Fatalf("strip ate the contract:\n%s", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Error("stripped file should end with exactly one newline")
	}
	if strings.HasSuffix(got, "\n\n") {
		t.Error("stripped file should not end with a blank line")
	}
}

func TestStripPreviewWithLeadingComment(t *testing.T) {
	// The real files put prose after the start marker, on the same comment.
	src := "export const X = 1\n\n/* grit:preview-start\n   Everything below is the demo. */\nexport default function P() { return null }\n/* grit:preview-end */\n"
	got := StripPreview(src)
	if strings.Contains(got, "demo") || strings.Contains(got, "function P") {
		t.Fatalf("preview survived:\n%s", got)
	}
}

func TestStripPreviewNoMarkers(t *testing.T) {
	// A file with no preview is already exactly what should land on disk.
	src := "export const Button = 1\n"
	if got := StripPreview(src); got != src {
		t.Errorf("StripPreview mangled a marker-less file: %q", got)
	}
}

func TestVerifyExports(t *testing.T) {
	full := `
export type ButtonVariant = "primary"
export type ButtonSize = "sm"
export function buttonClasses() { return "" }
export interface ButtonProps {}
export const Button = 1
`
	if err := VerifyExports("button@1", full); err != nil {
		t.Fatalf("complete variant rejected: %v", err)
	}

	// The failure this exists to catch: a variant that dropped a size union.
	partial := strings.Replace(full, `export type ButtonSize = "sm"`, "", 1)
	err := VerifyExports("button@1", partial)
	if err == nil {
		t.Fatal("variant missing ButtonSize was accepted")
	}
	if !strings.Contains(err.Error(), "ButtonSize") {
		t.Errorf("error should name the missing symbol, got: %v", err)
	}

	// An unknown contract is a variant from a newer CLI, not a failure.
	if err := VerifyExports("carousel@9", "export const Whatever = 1"); err != nil {
		t.Errorf("unknown contract should pass through, got: %v", err)
	}
}

func TestVerifyExportsDoesNotMatchSubstrings(t *testing.T) {
	// "ButtonSizes" must not satisfy a requirement for "ButtonSize" — the \b
	// in the pattern is what stops a near-miss variant from being accepted.
	src := `
export type ButtonVariant = "primary"
export type ButtonSizes = "sm"
export function buttonClasses() { return "" }
export interface ButtonProps {}
export const Button = 1
`
	if err := VerifyExports("button@1", src); err == nil {
		t.Fatal("ButtonSizes was accepted for ButtonSize")
	}
}

func TestFindProject(t *testing.T) {
	root := t.TempDir()

	if _, err := FindProject(root); err == nil {
		t.Error("a directory with no admin should be an error")
	}

	admin := filepath.Join(root, "apps", "admin")
	if err := os.MkdirAll(admin, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := FindProject(root)
	if err != nil {
		t.Fatalf("apps/admin not found: %v", err)
	}
	if p.SrcRoot != admin {
		t.Errorf("SrcRoot = %q, want %q", p.SrcRoot, admin)
	}
	want := filepath.Join(admin, "components", "ui", "button.tsx")
	if got := p.SlotPath("button"); got != want {
		t.Errorf("SlotPath = %q, want %q", got, want)
	}

	// The TanStack variant nests under src/, and the slot has to follow it —
	// otherwise the swap writes a file Vite never compiles.
	if err := os.MkdirAll(filepath.Join(admin, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	p2, err := FindProject(root)
	if err != nil {
		t.Fatal(err)
	}
	wantSrc := filepath.Join(admin, "src", "components", "ui", "button.tsx")
	if got := p2.SlotPath("button"); got != wantSrc {
		t.Errorf("TanStack SlotPath = %q, want %q", got, wantSrc)
	}
}

func TestStateRoundTrip(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "apps", "admin"), 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := FindProject(root)
	if err != nil {
		t.Fatal(err)
	}

	// A project that has never swapped is a normal project, not an error.
	if got := p.LoadState(); len(got) != 0 {
		t.Errorf("fresh project should have empty state, got %v", got)
	}

	want := State{"button": {Variant: "glow-ring", Contract: "button@1", SHA256: "abc"}}
	if err := p.saveState(want); err != nil {
		t.Fatal(err)
	}
	got := p.LoadState()
	if got["button"].Variant != "glow-ring" || got["button"].SHA256 != "abc" {
		t.Errorf("round trip lost data: %+v", got)
	}

	// Corrupt state must not take the CLI down with it.
	if err := os.WriteFile(p.statePath(), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := p.LoadState(); len(got) != 0 {
		t.Errorf("corrupt state should read as empty, got %v", got)
	}
}

func TestBackupAndRevert(t *testing.T) {
	root := t.TempDir()
	admin := filepath.Join(root, "apps", "admin", "components", "ui")
	if err := os.MkdirAll(admin, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := FindProject(root)
	if err != nil {
		t.Fatal(err)
	}
	slotPath := p.SlotPath("button")
	original := []byte("/* grit:slot button@1 */\nexport const Button = \"ORIGINAL\"\n")
	if err := os.WriteFile(slotPath, original, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Revert(p, "button"); err == nil {
		t.Error("revert with no backups should be an error, not a silent no-op")
	}

	if _, err := p.backup("button", original); err != nil {
		t.Fatal(err)
	}
	// Simulate the swap having overwritten it.
	if err := os.WriteFile(slotPath, []byte("SWAPPED"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Revert(p, "button"); err != nil {
		t.Fatalf("revert failed: %v", err)
	}
	back, err := os.ReadFile(slotPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(back) != string(original) {
		t.Errorf("revert restored %q, want the original", string(back))
	}
}
