// Package expose scaffolds public-facing Next.js pages in apps/web/
// that consume an already-generated Grit resource. Phase 3 of the
// PLAN_FORMS_AND_SHARING roadmap.
//
// Two entry points:
//
//	expose.Form(opts)  — emit a form page that creates records via the
//	                     auto-generated useCreate<Resource>() hook
//	expose.Table(opts) — emit a list page that reads via the
//	                     auto-generated use<Resources>() hook
//
// Unlike `grit generate resource`, these commands write a SINGLE file
// at an operator-chosen path. They don't touch the model, the API, the
// schemas, or any indexes — the resource is assumed to already exist.
package expose

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/generate"
)

// Opts is the common configuration for both Form and Table.
type Opts struct {
	// Resource is the PascalCase resource name (e.g. "Contact"). Must
	// match a model file at apps/api/internal/models/<snake>.go.
	Resource string

	// To is the destination path relative to the project root or an
	// absolute path. Must end in .tsx. The directory is created if it
	// doesn't exist; an existing file is refused unless Force is set.
	To string

	// Root is the project root (where grit.config.ts lives). Discovered
	// automatically by the caller via scaffold.FindProjectRoot().
	Root string

	// Force overwrites an existing file at To. Required for
	// idempotent CI runs; humans usually want the safety.
	Force bool
}

// Form scaffolds a Next.js client page that renders a form for the
// given resource and submits via the auto-generated React Query hook.
//
// The page validates with the shared Zod schema (Create<Resource>Schema
// from @repo/shared/schemas) and shows per-field errors. On success it
// shows a "Thank you" state and resets the form.
func Form(opts Opts) error {
	struct_, err := loadStruct(opts)
	if err != nil {
		return err
	}

	target, err := resolveTarget(opts)
	if err != nil {
		return err
	}

	content := buildFormPage(opts.Resource, struct_)
	return writeOnce(target, content, opts.Force)
}

// Table scaffolds a Next.js client page that renders a paginated,
// searchable list using the auto-generated React Query hook.
//
// The page is web-styled (plain Tailwind), not the heavy admin
// DataTable — it's meant to live on a public marketing site or a
// customer-facing dashboard, where the admin's CRM aesthetic is
// inappropriate.
func Table(opts Opts) error {
	struct_, err := loadStruct(opts)
	if err != nil {
		return err
	}

	target, err := resolveTarget(opts)
	if err != nil {
		return err
	}

	content := buildTablePage(opts.Resource, struct_)
	return writeOnce(target, content, opts.Force)
}

// loadStruct reads the resource's Go model file + extracts the struct
// matching the PascalCase resource name. Used by both Form and Table.
func loadStruct(opts Opts) (*generate.GoStruct, error) {
	if opts.Resource == "" {
		return nil, fmt.Errorf("resource name is required")
	}
	if !isPascal(opts.Resource) {
		return nil, fmt.Errorf("resource name must be PascalCase (got %q)", opts.Resource)
	}

	snake := pascalToSnake(opts.Resource)
	modelFile := filepath.Join(opts.Root, "apps", "api", "internal", "models", snake+".go")
	if _, err := os.Stat(modelFile); err != nil {
		return nil, fmt.Errorf("resource %q not found at %s — run `grit generate resource %s` first",
			opts.Resource, modelFile, opts.Resource)
	}

	structs, err := generate.ParseGoStructs(modelFile)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", modelFile, err)
	}

	for _, s := range structs {
		if s.Name == opts.Resource {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("struct %s not found in %s — generator output may be stale; run `grit sync`",
		opts.Resource, modelFile)
}

// resolveTarget normalises opts.To: enforces .tsx, normalises
// separators, returns the absolute final path.
func resolveTarget(opts Opts) (string, error) {
	if opts.To == "" {
		return "", fmt.Errorf("--to is required (e.g. --to apps/web/app/contact-us/page.tsx)")
	}
	if !strings.HasSuffix(opts.To, ".tsx") {
		return "", fmt.Errorf("--to must end in .tsx (got %q)", opts.To)
	}
	target := opts.To
	if !filepath.IsAbs(target) {
		target = filepath.Join(opts.Root, target)
	}
	return filepath.Clean(target), nil
}

// writeOnce creates the parent directory + writes the file. Refuses
// to overwrite without --force so an operator can't accidentally
// blow away a hand-customised page.
func writeOnce(path, content string, force bool) error {
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("refusing to overwrite %s — pass --force to replace", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ── small string helpers (duplicated to keep this package self-contained) ──

func isPascal(s string) bool {
	if s == "" {
		return false
	}
	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func pascalToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r - 'A' + 'a')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// pluralKebab converts "Contact" → "contacts", "Person" → "people"-ish
// (very small heuristic — matches the generator's Pluralize for the
// common cases). When the heuristic fails the operator can rename the
// generated file by hand.
func pluralKebab(s string) string {
	snake := pascalToSnake(s)
	switch {
	case strings.HasSuffix(snake, "y"):
		return strings.TrimSuffix(snake, "y") + "ies"
	case strings.HasSuffix(snake, "s"):
		return snake + "es"
	default:
		return snake + "s"
	}
}

// pluralPascal returns the PascalCase plural — what the React Query
// hooks are exported as (`useContacts`, `useGetContact`, etc).
func pluralPascal(s string) string {
	plural := pluralKebab(s)
	parts := strings.Split(plural, "_")
	for i := range parts {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
