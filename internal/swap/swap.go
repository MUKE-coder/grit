// Package swap replaces a single canonical admin component — a "slot" — with a
// variant from the Grit UI registry.
//
// This is deliberately NOT the same operation as `grit ui add`. Adding gives you
// a new file to import wherever you like. Swapping overwrites the one file every
// call site already imports, so `grit swap button glow-ring` restyles every
// button in the admin without touching a single import.
//
// That power is why most of this package is refusals:
//
//   - a variant whose contract major differs from the installed slot is rejected
//     rather than written, because it type-checks against a shape the call sites
//     no longer use;
//   - a slot file you have edited by hand is not overwritten without --force,
//     because "swap" should never silently mean "discard my work";
//   - the previous file is always backed up first, so revert is a real thing
//     rather than a suggestion to check git;
//   - and by default the project is type-checked afterwards, with an automatic
//     rollback when it fails. A swap that leaves the app not compiling is worse
//     than one that refuses.
package swap

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/uiregistry"
)

// RequiredExports lists the symbols a variant must export to satisfy a contract.
//
// Checked textually before anything is written. It will not catch a variant that
// exports the right names with the wrong behaviour — that is what the type-check
// pass is for — but it does catch the common case of a variant built against an
// older shape, instantly and without a compiler.
var RequiredExports = map[string][]string{
	"button@1": {"Button", "buttonClasses", "ButtonProps", "ButtonVariant", "ButtonSize"},
	"input@1":  {"Input", "inputClasses", "InputProps", "InputSize"},
}

// KnownSlots is the set a project can swap today. Kept explicit so an unknown
// slot fails with a list rather than by writing a file nothing imports.
var KnownSlots = []string{"button", "input"}

const (
	stateDir  = ".grit"
	stateFile = "slots.json"
	backupDir = "swaps"
)

// slotHeader matches the contract declaration every slot file carries:
//
//	/* grit:slot button@1
var slotHeader = regexp.MustCompile(`grit:slot\s+([a-z0-9-]+@\d+)`)

// previewBlock matches the demo the registry site renders. Stripped on swap so
// the file that lands in a project contains only the contract — nobody wants a
// Preview() component sitting in their button.tsx.
var previewBlock = regexp.MustCompile(`(?s)\n?/\*\s*grit:preview-start.*?grit:preview-end\s*\*/\n?`)

// Record is what we wrote for one slot, so a later swap can tell "unchanged
// since I wrote it" from "the user has edited this".
type Record struct {
	Variant   string `json:"variant"`
	Contract  string `json:"contract"`
	SHA256    string `json:"sha256"`
	SwappedAt string `json:"swappedAt"`
}

// State is .grit/slots.json — slot name to record.
type State map[string]Record

// Project is a resolved admin app.
type Project struct {
	Root     string // repository root
	AdminDir string // e.g. <root>/apps/admin
	Label    string // e.g. "apps/admin", for messages
	SrcRoot  string // AdminDir, or AdminDir/src for the TanStack variant
}

// FindProject locates the admin app.
//
// Swapping is admin-only on purpose: the marketing site and the admin have
// different primitives, and a slot that means two different things in two apps
// is not a slot.
func FindProject(root string) (*Project, error) {
	candidates := []struct{ path, label string }{
		{filepath.Join("apps", "admin"), "apps/admin"},
		{"admin", "admin"},
	}
	for _, c := range candidates {
		dir := filepath.Join(root, c.path)
		if st, err := os.Stat(dir); err == nil && st.IsDir() {
			p := &Project{Root: root, AdminDir: dir, Label: c.label, SrcRoot: dir}
			// The TanStack variant nests everything under src/.
			if st, err := os.Stat(filepath.Join(dir, "src")); err == nil && st.IsDir() {
				p.SrcRoot = filepath.Join(dir, "src")
			}
			return p, nil
		}
	}
	return nil, fmt.Errorf(
		"no admin app found (looked for apps/admin and admin/)\n\n" +
			"grit swap only targets the admin panel. Run it from your project root.")
}

// SlotPath is where the slot lives, e.g. <admin>/components/ui/button.tsx.
func (p *Project) SlotPath(slot string) string {
	return filepath.Join(p.SrcRoot, "components", "ui", slot+".tsx")
}

func (p *Project) statePath() string {
	return filepath.Join(p.Root, stateDir, stateFile)
}

// LoadState reads .grit/slots.json. A missing or unreadable file is an empty
// state, not an error — a project that has never swapped is a normal project.
func (p *Project) LoadState() State {
	b, err := os.ReadFile(p.statePath())
	if err != nil {
		return State{}
	}
	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return State{}
	}
	return s
}

func (p *Project) saveState(s State) error {
	if err := os.MkdirAll(filepath.Join(p.Root, stateDir), 0o755); err != nil {
		return fmt.Errorf("create %s: %w", stateDir, err)
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode slot state: %w", err)
	}
	return os.WriteFile(p.statePath(), append(b, '\n'), 0o644)
}

func hashOf(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// ContractOf reads the grit:slot declaration out of a slot file.
func ContractOf(src string) string {
	if m := slotHeader.FindStringSubmatch(src); m != nil {
		return m[1]
	}
	return ""
}

// StripPreview removes the registry-only demo block.
func StripPreview(src string) string {
	return strings.TrimRight(previewBlock.ReplaceAllString(src, "\n"), "\n") + "\n"
}

// Major returns the version from "button@1". Empty when unparseable.
func Major(contract string) string {
	if i := strings.LastIndex(contract, "@"); i >= 0 {
		return contract[i+1:]
	}
	return ""
}

// VerifyExports checks a variant declares everything the contract promises.
func VerifyExports(contract, src string) error {
	required, ok := RequiredExports[contract]
	if !ok {
		// An unknown contract is not a failure — it is a variant published
		// against a newer CLI. The type-check pass is still the real gate.
		return nil
	}
	var missing []string
	for _, name := range required {
		// Matches `export const X`, `export function X`, `export type X`,
		// `export interface X`.
		re := regexp.MustCompile(`export\s+(const|function|type|interface)\s+` + regexp.QuoteMeta(name) + `\b`)
		if !re.MatchString(src) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("variant does not satisfy %s: missing export(s): %s",
			contract, strings.Join(missing, ", "))
	}
	return nil
}

// Options controls one Apply.
type Options struct {
	// Overwrite a slot file that has been edited since it was written.
	Force bool
	// Skip the post-swap type-check.
	SkipCheck bool
}

// Result describes what Apply did.
type Result struct {
	Slot       string
	Variant    string
	Contract   string
	Path       string
	BackupPath string
	Checked    bool
}

// Apply swaps one slot to one variant.
func Apply(ctx context.Context, p *Project, slot, variant string, opts Options) (*Result, error) {
	slotPath := p.SlotPath(slot)
	current, err := os.ReadFile(slotPath)
	if err != nil {
		return nil, fmt.Errorf("no %s slot in %s\n\nExpected %s.\nSlots arrived in v3.115.0: a project scaffolded before that needs its admin regenerating",
			slot, p.Label, filepath.ToSlash(strings.TrimPrefix(slotPath, p.Root+string(filepath.Separator))))
	}

	// The user's edits are theirs. Refuse rather than discard.
	state := p.LoadState()
	if rec, ok := state[slot]; ok && rec.SHA256 != "" && rec.SHA256 != hashOf(current) && !opts.Force {
		return nil, fmt.Errorf(
			"%s has been edited since it was last swapped\n\n"+
				"Swapping would discard those changes. Re-run with --force to overwrite\n"+
				"(the current file is backed up either way).",
			filepath.ToSlash(strings.TrimPrefix(slotPath, p.Root+string(filepath.Separator))))
	}

	installedContract := ContractOf(string(current))

	// Fetch the variant.
	name := fmt.Sprintf("application-ui-%ss-%s", slot, variant)
	comp, err := uiregistry.Fetch(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("fetch variant %q for slot %q: %w", variant, slot, err)
	}
	if len(comp.Files) == 0 || strings.TrimSpace(comp.Files[0].Content) == "" {
		return nil, fmt.Errorf("registry returned no source for %q", name)
	}

	incoming := StripPreview(comp.Files[0].Content)
	incomingContract := ContractOf(incoming)

	if incomingContract == "" {
		return nil, fmt.Errorf("%q is not a swappable variant: it declares no grit:slot contract", variant)
	}
	if want := slot + "@"; !strings.HasPrefix(incomingContract, want) {
		return nil, fmt.Errorf("%q implements %s, which is not the %s slot", variant, incomingContract, slot)
	}
	// Refuse a major mismatch. Writing it would compile against a shape the
	// call sites stopped using two releases ago.
	if installedContract != "" && Major(incomingContract) != Major(installedContract) {
		return nil, fmt.Errorf(
			"contract mismatch: your admin is on %s, %q implements %s\n\n"+
				"Upgrade the CLI, or pick a variant published against %s.",
			installedContract, variant, incomingContract, installedContract)
	}
	if err := VerifyExports(incomingContract, incoming); err != nil {
		return nil, err
	}

	// Back up before writing, always — including under --force, which is
	// exactly when someone most wants the old file back.
	backup, err := p.backup(slot, current)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(slotPath, []byte(incoming), 0o644); err != nil {
		return nil, fmt.Errorf("write %s: %w", slotPath, err)
	}

	res := &Result{
		Slot: slot, Variant: variant, Contract: incomingContract,
		Path: slotPath, BackupPath: backup,
	}

	if !opts.SkipCheck {
		if err := TypeCheck(ctx, p); err != nil {
			// Put it back. A swap that leaves the app not compiling is worse
			// than one that refuses.
			_ = os.WriteFile(slotPath, current, 0o644)
			return nil, fmt.Errorf("%s does not type-check against your admin: reverted\n\n%v", variant, err)
		}
		res.Checked = true
	}

	state[slot] = Record{
		Variant:   variant,
		Contract:  incomingContract,
		SHA256:    hashOf([]byte(incoming)),
		SwappedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := p.saveState(state); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Project) backup(slot string, content []byte) (string, error) {
	dir := filepath.Join(p.Root, stateDir, backupDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create backup dir: %w", err)
	}
	path := filepath.Join(dir, fmt.Sprintf("%s.%s.tsx", slot, time.Now().UTC().Format("20060102-150405")))
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}
	return path, nil
}

// Revert restores the most recent backup for a slot.
func Revert(p *Project, slot string) (string, error) {
	dir := filepath.Join(p.Root, stateDir, backupDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("nothing to revert: no backups in %s/%s", stateDir, backupDir)
	}
	var candidates []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), slot+".") && strings.HasSuffix(e.Name(), ".tsx") {
			candidates = append(candidates, e.Name())
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("nothing to revert: no backup for slot %q", slot)
	}
	// Names embed a sortable UTC timestamp, so lexical order is chronological.
	sort.Strings(candidates)
	newest := candidates[len(candidates)-1]

	content, err := os.ReadFile(filepath.Join(dir, newest))
	if err != nil {
		return "", fmt.Errorf("read backup: %w", err)
	}
	if err := os.WriteFile(p.SlotPath(slot), content, 0o644); err != nil {
		return "", fmt.Errorf("restore %s: %w", slot, err)
	}

	state := p.LoadState()
	delete(state, slot)
	if err := p.saveState(state); err != nil {
		return "", err
	}
	return newest, nil
}

// TypeCheck runs the admin's TypeScript compiler.
//
// This is the real contract test. A variant can export every required symbol and
// still take a size union the call sites do not use; tsc is what notices, and it
// notices at every call site rather than at the one the author happened to think
// of.
func TypeCheck(ctx context.Context, p *Project) error {
	if _, err := os.Stat(filepath.Join(p.AdminDir, "node_modules")); err != nil {
		return fmt.Errorf("dependencies are not installed in %s: run your package manager there first, or pass --skip-check", p.Label)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npx", "tsc", "--noEmit")
	cmd.Dir = p.AdminDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 12 {
		lines = append(lines[:12], fmt.Sprintf("... and %d more", len(lines)-12))
	}
	return fmt.Errorf("%s", strings.Join(lines, "\n"))
}

// Variant is one swappable option in the registry.
type Variant struct {
	Slot        string
	Name        string // "glow-ring"
	Title       string
	Description string
	Contract    string
	Pro         bool
}

// ListVariants returns every swappable variant, optionally filtered to one slot.
func ListVariants(ctx context.Context, slot string) ([]Variant, error) {
	items, err := uiregistry.List(ctx)
	if err != nil {
		return nil, err
	}
	var out []Variant
	for _, it := range items {
		if it.Slot == "" {
			continue
		}
		if slot != "" && it.Slot != slot {
			continue
		}
		// Registry names are flat: application-ui-buttons-glow-ring. The
		// variant is whatever follows "<slot>s-".
		short := it.Name
		if i := strings.Index(short, it.Slot+"s-"); i >= 0 {
			short = short[i+len(it.Slot)+2:]
		}
		out = append(out, Variant{
			Slot: it.Slot, Name: short, Title: it.Title,
			Description: it.Description, Contract: it.Contract, Pro: it.Pro,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Slot != out[j].Slot {
			return out[i].Slot < out[j].Slot
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}
