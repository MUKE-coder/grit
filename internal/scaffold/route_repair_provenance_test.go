package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// handWrittenInterchangeRoutes is the file from grit#213, as near as the issue
// reports it: written by hand, mounting two student-facing endpoints that are
// scoped by owner inside the handler rather than by authz.ScopeToOwner.
//
// Every anchor the authorization repair looks for matches. It holds
// `h := &handlers.<Pascal>Handler{`, its handler is not in the owned set, and
// its two methods are Import and Export, which are names in routePermission.
// That is why the repair rewrote it, and why no amount of care in the anchors
// would have stopped it: the generator's shape is the convention the docs
// teach, so hand-written code that follows the convention is indistinguishable
// from generated code by inspection.
const handWrittenInterchangeRoutes = `package routes

import (
	"viva/apps/api/internal/handlers"
)

// Hand-written. Students import and export their own reference library; the
// handler narrows every query to the caller.
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.ReferenceInterchangeHandler{
			DB: m.DB,
		}

		m.Protected.POST("/projects/:id/references/import", h.Import)
		m.Protected.GET("/projects/:id/references/export", h.Export)
	})
}
`

// writeProvenanceProject lays out a project with one generated route file and
// one hand-written one, and records only the generated one in the manifest.
func writeProvenanceProject(t *testing.T, recordGenerated bool) (root string, generated, handWritten string) {
	t.Helper()
	root = t.TempDir()
	routesDir := filepath.Join(root, "apps", "api", "internal", "routes")
	if err := os.MkdirAll(routesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	generated = filepath.Join(routesDir, "contact_routes.go")
	handWritten = filepath.Join(routesDir, "reference_interchange_routes.go")
	writeFileT(t, generated, contactRoutesV239)
	writeFileT(t, handWritten, handWrittenInterchangeRoutes)

	m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
	if recordGenerated {
		m.Record("apps/api/internal/routes/contact_routes.go", "resource:Contact", "3.239.0", contactRoutesV239)
	}
	if err := m.Save(root, "3.239.0"); err != nil {
		t.Fatal(err)
	}
	return root, generated, handWritten
}

// The fault in grit#213: an upgrade rewrote a hand-written route file and
// changed the authorization on two endpoints, and the summary did not say so.
func TestAHandWrittenRouteFileIsNotRewritten(t *testing.T) {
	root, generated, handWritten := writeProvenanceProject(t, true)
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}

	p := routeFileProvenance(root, m, []string{generated, handWritten})
	if !p.known {
		t.Fatal("the manifest tracks a route file, so provenance should be authoritative here")
	}
	if allowed, _ := p.mayRewrite(generated); !allowed {
		t.Error("the generated route file must still be repaired: it is the exposure the repair exists for")
	}
	allowed, why := p.mayRewrite(handWritten)
	if allowed {
		t.Error("a route file Grit never wrote was cleared for an authorization rewrite")
	}
	if !strings.Contains(why, "not one Grit wrote") {
		t.Errorf("the refusal reads %q, which does not say why", why)
	}
}

// And the refusal has to say what it declined to do. "Nothing happened" on a
// file holding two endpoints the repair would have moved is not useful; the
// developer has to be able to decide.
func TestTheRefusalNamesTheRoutesItLeftAlone(t *testing.T) {
	would := untouchedRouteRewrites(handWrittenInterchangeRoutes)
	if len(would) != 2 {
		t.Fatalf("named %d routes, want 2: %v", len(would), would)
	}
	joined := strings.Join(would, "\n")
	for _, want := range []string{
		"POST /projects/:id/references/import",
		"GET /projects/:id/references/export",
		"perm:projects.create",
		"perm:projects.view",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("the refusal does not mention %q:\n%s", want, joined)
		}
	}
}

// The inverse failure, which is the one a gate like this usually causes: a
// project generated before v3.147.0 has no manifest, so every file in it reads
// as Untracked. Gating on the manifest alone would stop repairing exactly the
// projects whose resource routes are still open to any signed-in account.
func TestAProjectWithNoManifestIsStillRepaired(t *testing.T) {
	root, generated, handWritten := writeProvenanceProject(t, false)
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}

	p := routeFileProvenance(root, m, []string{generated, handWritten})
	if p.known {
		t.Fatal("no route file is recorded, so the manifest cannot be authoritative")
	}
	for _, path := range []string{generated, handWritten} {
		if allowed, why := p.mayRewrite(path); !allowed {
			t.Errorf("%s was skipped in a project with no manifest coverage (%s): an old project is the one that needs the fix",
				filepath.Base(path), why)
		}
	}
}

// A file Grit wrote and the developer then edited is still Grit's file. The
// repair anchors on exact generator text, so it leaves the developer's
// additions alone, and holding it back would leave the open routes open.
func TestAnEditedGeneratedRouteFileIsStillRepaired(t *testing.T) {
	root, generated, _ := writeProvenanceProject(t, true)
	edited := contactRoutesV239 + "\n// A comment the developer added.\n"
	writeFileT(t, generated, edited)

	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := m.StatusOf(root, "apps/api/internal/routes/contact_routes.go"); got != manifest.Modified {
		t.Fatalf("status is %v, want modified", got)
	}
	p := routeFileProvenance(root, m, []string{generated})
	if allowed, why := p.mayRewrite(generated); !allowed {
		t.Errorf("an edited generated file was skipped (%s); the gate is on provenance, not on cleanliness", why)
	}
}

// handWrittenRouteFiles is what reporting reads. It must be silent rather than
// wrong on a project whose manifest knows nothing about routes.
func TestHandWrittenRouteFilesListsOnlyTheUntracked(t *testing.T) {
	root, _, _ := writeProvenanceProject(t, true)
	apiRoot := filepath.Join(root, "apps", "api")

	got, err := handWrittenRouteFiles(root, apiRoot)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"apps/api/internal/routes/reference_interchange_routes.go"}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("listed %v, want %v", got, want)
	}

	bare, _, _ := writeProvenanceProject(t, false)
	got, err = handWrittenRouteFiles(bare, filepath.Join(bare, "apps", "api"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("listed %v for a project with no manifest coverage, want nothing: every file reads as untracked there", got)
	}
}
