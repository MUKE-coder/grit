// Package explain answers "why is this file here, and may I edit it?"
//
// Two halves, and they come from different places on purpose.
//
// What a file IS comes from a catalogue of path shapes, below. Those are
// conventions, stable across projects, and a sentence about each is knowledge
// a reader does not have.
//
// Whether you may EDIT it comes from the manifest, which is a record of what
// Grit actually wrote in THIS project. That matters: a hardcoded list of
// "framework-owned" paths would be a guess about a project it has not read,
// and would be wrong for any file the owner has already changed. The manifest
// knows the difference between a file Grit wrote and you have not touched, one
// you have edited, and one that was never Grit's.
//
// The rule that decides the advice is the one `grit upgrade` actually follows:
// a file Grit wrote and nobody edited may be rewritten; a file you edited is
// left alone. So "may I edit this" really means "do I mind losing the
// improvements, or maintaining it myself from now on", and that is what the
// output says.
package explain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Kind is the editing advice for a file.
type Kind int

const (
	// Generated means Grit wrote it and has not been edited since. An
	// upgrade may rewrite it.
	Generated Kind = iota
	// Edited means Grit wrote it and somebody changed it. An upgrade leaves
	// it alone, and it stops receiving improvements.
	Edited
	// Yours means Grit has no record of writing it.
	Yours
	// Gone means Grit wrote it and it is not there any more.
	Gone
)

// Explanation is the answer for one file.
type Explanation struct {
	// Rel is the path relative to the project root, with forward slashes.
	Rel string
	// What the file is for, one or two sentences. Empty when the path
	// matches no known shape, which is not a failure: a project has files
	// Grit never heard of, and saying nothing is better than guessing.
	Purpose string
	// Edit is when you would reasonably change it.
	Edit string
	// Read is where to read more, as a URL or a doc path. May be empty.
	Read string

	Kind Kind
	// Generator is who wrote it, from the manifest: "scaffold",
	// "resource:Product", "plugin:webhooks". Empty when untracked.
	Generator string
	// Version is the Grit version that last wrote it.
	Version string
	// Source is the file this one was generated from, when there is one.
	Source string
}

// pattern is one path shape and what to say about it.
type pattern struct {
	// match is tested against the forward-slashed relative path.
	match func(rel string) bool
	purpose,
	edit,
	read string
}

func dir(prefixes ...string) func(string) bool {
	return func(rel string) bool {
		for _, p := range prefixes {
			if strings.Contains(rel, p) {
				return true
			}
		}
		return false
	}
}

func base(names ...string) func(string) bool {
	return func(rel string) bool {
		b := filepath.Base(rel)
		for _, n := range names {
			if b == n {
				return true
			}
		}
		return false
	}
}

// patterns are tested in order, so the specific ones come first.
//
// Every entry here is a convention a generated project actually follows. A
// pattern that matched nothing would be harmless and useless; one that
// matches the wrong thing is worse than none, which is why the file-name
// matches come before the directory ones.
var patterns = []pattern{
	{
		match:   base("grit.json"),
		purpose: "The project's own description: its architecture, its frontend, and which features are on. Every Grit command reads it to find the project root and to know the shape of what it is working on.",
		edit:    "Rarely by hand. Commands that change the project's shape update it themselves.",
		read:    "/docs/getting-started/configuration",
	},
	{
		match:   dir(".grit/manifest.json"),
		purpose: "The record of every file Grit has written here and what it looked like. It is how an upgrade can tell a file you edited from one you have not touched, and leave yours alone.",
		edit:    "Never by hand. Deleting it makes an upgrade treat every file as yours, so you stop receiving framework fixes.",
		read:    "/docs/upgrades",
	},
	{
		match:   dir(".grit/error-codes.yaml"),
		purpose: "Your project's own API error codes. `grit sync` folds them into the Go constants and the TypeScript union, so both sides know the same set.",
		edit:    "Yes, this one is yours. Add a code here, then run grit sync.",
		read:    "/docs/backend/errors",
	},
	{
		match:   dir("internal/models/"),
		purpose: "A GORM model: the table, its columns, its validation rules and its hooks. This is the source the generator derives the handler, the shared type and the admin screen from.",
		edit:    "Yes. This is where a field's type, tags and rules live. Run grit sync afterwards so the TypeScript and Zod follow.",
		read:    "/docs/backend/models",
	},
	{
		match:   dir("internal/handlers/"),
		purpose: "The HTTP layer for one resource: parse the request, call a service, write the response. It holds no business logic by design.",
		edit:    "When the HTTP shape changes: a new route, a different status, an extra query parameter. Business rules belong in the service.",
		read:    "/docs/backend/handlers",
	},
	{
		match:   dir("internal/services/"),
		purpose: "The business rules for one resource, and where a transaction begins and ends. Called by handlers, by jobs and by the CLI alike.",
		edit:    "Yes, constantly. This is where most of what your product actually does should live.",
		read:    "/docs/backend/handlers",
	},
	{
		match:   dir("internal/policies/"),
		purpose: "Per-resource authorisation: a function that sees both the user and the record and decides whether the action is allowed.",
		edit:    "Yes. Rules a role cannot express go here, such as an author editing their own post for thirty minutes.",
		read:    "/docs/security/policies",
	},
	{
		match:   dir("internal/access/"),
		purpose: "The static table of every route and what guards it. It exists because a web framework cannot report its own middleware chain at runtime, so without this nothing can answer whether an endpoint is protected.",
		edit:    "No. It is regenerated from the routes, and grit doctor reads it to find an unguarded route.",
		read:    "/docs/security/authorization",
	},
	{
		match:   dir("internal/middleware/"),
		purpose: "Request middleware: authentication, CORS, rate limiting, the request id, logging.",
		edit:    "When you add a cross-cutting concern. Register it before the routes it should guard, or it does not run for them.",
		read:    "/docs/backend/middleware",
	},
	{
		match:   dir("internal/routes/"),
		purpose: "The route table: which path maps to which handler, and which middleware guards it. Generated resources are injected here at a marker.",
		edit:    "Yes, for routes you add by hand. Leave the markers in place or future resources stop appearing.",
		read:    "/docs/backend/handlers",
	},
	{
		match:   dir("internal/jobs/"),
		purpose: "A background job: its payload, and the handler that runs it off the queue.",
		edit:    "Yes. The handler must be safe to run twice, because a retry will.",
		read:    "/docs/batteries/jobs",
	},
	{
		match:   dir("migrations/"),
		purpose: "A versioned SQL migration pair. These are the record of how the database got to its current shape, and the only safe way to change it in production.",
		edit:    "Never after it has run anywhere. Write a new migration instead; the runner checksums these and will refuse a changed one.",
		read:    "/docs/backend/migrations",
	},
	{
		match:   dir("packages/shared/types/", "packages/shared/schemas/"),
		purpose: "Generated from the Go models: the TypeScript type and the Zod schema for one resource, so the frontend cannot disagree with the backend about a shape.",
		edit:    "No. grit sync rewrites these from the Go, so an edit is lost on the next run. Change the model instead.",
		read:    "/docs/frontend/shared-package",
	},
	{
		match:   dir("/resources/"),
		purpose: "An admin resource definition: which columns the table shows, which filters it has, and which fields the form carries. The table, the form and the detail view are all derived from it.",
		edit:    "Yes, this is the main thing you change about an admin screen. A .custom file beside it is yours and is never overwritten.",
		read:    "/docs/admin/resources",
	},
	{
		match:   dir("/hooks/use-"),
		purpose: "A React Query hook for one resource: the fetching, the caching and the invalidation after a write.",
		edit:    "Usually not. It is regenerated with the resource. Put bespoke fetching in a hook of your own beside it.",
		read:    "/docs/frontend/hooks",
	},
	{
		match:   base("docker-compose.yml", "docker-compose.prod.yml"),
		purpose: "The local infrastructure: Postgres, Redis, MinIO and Mailhog, on ports chosen so a native install of any of them does not collide.",
		edit:    "Yes, for a service you add. Note that the production file publishes no database port at all, on purpose.",
		read:    "/docs/infrastructure/docker",
	},
	{
		match:   base(".env", ".env.example"),
		purpose: "Configuration. .env is yours and is not committed; .env.example is the list of what a new machine needs.",
		edit:    "Yes. Add a variable to both, or somebody cloning the project will not know it exists.",
		read:    "/docs/getting-started/configuration",
	},
	{
		match:   base("main.go"),
		purpose: "A binary's entry point: construct the dependencies, wire them together and start. There is one per cmd/ directory, usually the server and the worker.",
		edit:    "When you add a dependency that everything else needs.",
		read:    "/docs/getting-started/project-structure",
	},
	{
		match:   dir("internal/config/"),
		purpose: "The environment read into a struct and validated at startup, so a missing variable stops the process at boot rather than at the first request that needs it.",
		edit:    "Yes, when you add a setting.",
		read:    "/docs/getting-started/configuration",
	},
	{
		match:   dir("internal/seeders/", "/seeders/"),
		purpose: "Starting data: an admin account, reference rows, or realistic data for development. Run with grit seed.",
		edit:    "Yes. Keep it idempotent, because it will be run twice.",
		read:    "/docs/backend/models",
	},
	{
		match:   dir("/components/ui/"),
		purpose: "The hand-written primitives the admin is built from, against the design tokens rather than a component library.",
		edit:    "Yes, they are yours. They are also shared by every screen, so a change here is a change everywhere.",
		read:    "/docs/admin/overview",
	},
	{
		match:   dir("_test.go", ".test.ts", ".spec.ts"),
		purpose: "A test. The generated ones cover the shape of what was generated; the valuable ones are the cases you add.",
		edit:    "Yes. Make one fail on purpose once, so you know it can.",
		read:    "/docs/testing",
	},
}

// Explain answers for one path. root is the project root, path is anything
// the user typed: absolute, or relative to where they are standing.
func Explain(root, path string) (*Explanation, error) {
	abs := path
	if !filepath.IsAbs(abs) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("working directory: %w", err)
		}
		abs = filepath.Join(wd, path)
	}
	abs = filepath.Clean(abs)

	rel, ok := manifest.Rel(root, abs)
	if !ok {
		return nil, fmt.Errorf("%s is outside this project", path)
	}
	rel = filepath.ToSlash(rel)

	e := &Explanation{Rel: rel}

	for _, p := range patterns {
		if p.match(rel) {
			e.Purpose, e.Edit, e.Read = p.purpose, p.edit, p.read
			break
		}
	}

	// The manifest decides the editing advice, because it is the only thing
	// that knows what happened in THIS project.
	if m, err := manifest.Load(root); err == nil {
		if entry, found := m.Files[manifest.Key(rel)]; found {
			e.Generator, e.Version = entry.Generator, entry.Grit
			if entry.From != nil {
				e.Source = entry.From.Path
			}
		}
		switch m.StatusOf(root, rel) {
		case manifest.Unchanged:
			e.Kind = Generated
		case manifest.Modified:
			e.Kind = Edited
		case manifest.Missing:
			e.Kind = Gone
		default:
			e.Kind = Yours
		}
	} else {
		e.Kind = Yours
	}

	return e, nil
}

// Advice is the sentence about editing, derived from what the manifest knows
// rather than from the path.
func (e *Explanation) Advice() string {
	switch e.Kind {
	case Generated:
		return "Grit wrote this and it has not been edited, so an upgrade may rewrite it. " +
			"Edit it and your change is kept, and it stops receiving framework improvements."
	case Edited:
		return "Grit wrote this and it has been edited since, so an upgrade will leave it alone. " +
			"It no longer receives framework improvements."
	case Gone:
		return "Grit wrote this and it is not there any more. An upgrade will put it back."
	default:
		return "Grit has no record of writing this, so it is yours and nothing will touch it."
	}
}
