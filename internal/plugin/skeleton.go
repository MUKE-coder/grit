package plugin

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// grit#143: a community plugin author started from an empty directory.
//
// The authoring docs describe the manifest and point at internal/plugin for an
// example, which is inside the CLI and written in Go. So somebody writing a
// JSON plugin outside this repository had the shape described to them and
// nothing to copy, and the failure modes the issue lists are all late ones: a
// marker typo injects nothing, a missing when clause writes an admin file into
// a project with no admin, a path that climbs out is refused at install.
//
// The issue asks for a template repository, a `grit plugin new` command, or
// both, and makes the case for the command itself: a repository drifts from
// the manifest schema the first time the schema gains a field, and a skeleton
// generated from templates inside the CLI cannot. The test below is the other
// half of that argument. It writes the skeleton and loads it with LoadDir, the
// same function `grit plugin add ./dir` calls, so a change to the schema that
// the skeleton does not follow fails here rather than in somebody's terminal.
//
// What the skeleton exercises, which is the list the issue asked for: one file
// written unconditionally, one gated by a when clause, one Go test, two
// injections (an import and a route, because one without the other does not
// compile), and the install and remove loop in CI.

// SkeletonMarkers are the markers the generated plugin injects at.
//
// Exported because the test that matters compares them against the scaffold's
// own routes.go, rather than against a copy of the list.
var SkeletonMarkers = []string{"// grit:imports", "// grit:routes:protected"}

var pluginNameRe = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

// ValidatePluginName accepts the names `grit plugin add` can later refer to.
//
// Lowercase and hyphenated, because the name is a directory, a key in the
// lockfile, an argument on a command line and part of a URL path, and the one
// spelling that is safe in all four is this one. Refused early: a name with a
// separator or a capital in it installs and then cannot be removed by the name
// the person typed.
func ValidatePluginName(name string) error {
	switch {
	case strings.TrimSpace(name) == "":
		return fmt.Errorf("a plugin needs a name")
	case strings.ContainsAny(name, `/\`):
		return fmt.Errorf("%q: a plugin name is a single name, not a path", name)
	case name != strings.ToLower(name):
		return fmt.Errorf("%q: lowercase only, so %q", name, strings.ToLower(name))
	case !pluginNameRe.MatchString(name):
		return fmt.Errorf("%q: letters, digits and single hyphens, starting with a letter (for example \"audit-trail\")", name)
	}
	// A built-in shadows a directory of the same name: grit plugin add refuses
	// to read a path as a name, but grit plugin remove takes the name, and two
	// things answering to one name is a bad place to start.
	for _, built := range All() {
		if name == built.Name {
			return fmt.Errorf("%q is a built-in plugin; pick another name", name)
		}
	}
	return nil
}

// goPackage is the Go package name for a plugin name: audit-trail becomes
// audittrail, because a hyphen is not legal in an identifier and the
// alternative, audit_trail, is not idiomatic Go.
func goPackage(name string) string {
	return strings.ReplaceAll(name, "-", "")
}

// SkeletonFiles is the plugin directory, keyed by path relative to it.
//
// Returned as a map rather than written here so the command can refuse to
// overwrite before anything lands on disk, and so the test can load it from a
// temporary directory.
func SkeletonFiles(name string) map[string]string {
	pkg := goPackage(name)
	title := strings.ReplaceAll(name, "-", " ")

	return map[string]string{
		"plugin.json":                                 skeletonManifest(name, pkg),
		path.Join("files", "api", pkg+".go"):          skeletonService(name, pkg),
		path.Join("files", "api", pkg+"_test.go"):     skeletonServiceTest(pkg),
		path.Join("files", "admin", name+"-card.tsx"): skeletonAdminCard(name, title),
		"README.md":    skeletonReadme(name, pkg),
		"CHANGELOG.md": skeletonChangelog(),
		".gitignore":   "node_modules/\n.DS_Store\n",
		path.Join(".github", "workflows", "plugin.yml"): skeletonWorkflow(name),
	}
}

// skeletonManifest is plugin.json.
//
// Every field the validator requires, and nothing it does not, with the
// optional ones named in comments the JSON cannot carry. The README explains
// them instead, because a JSON file with no comments is the one place an
// author cannot be told anything.
func skeletonManifest(name, pkg string) string {
	return `{
  "name": "` + name + `",
  "version": "0.1.0",
  "summary": "One line. This is what grit plugin list shows.",
  "description": "A few sentences on what this installs and why somebody would want it. grit plugin info shows this.",
  "nextSteps": [
    "Restart the API server so the new route is served.",
    "GET /api/` + name + ` returns the message."
  ],
  "files": [
    {
      "from": "files/api/` + pkg + `.go",
      "to": "{{API_ROOT}}/internal/` + pkg + `/` + pkg + `.go"
    },
    {
      "from": "files/api/` + pkg + `_test.go",
      "to": "{{API_ROOT}}/internal/` + pkg + `/` + pkg + `_test.go"
    },
    {
      "from": "files/admin/` + name + `-card.tsx",
      "to": "apps/admin/components/` + name + `/` + name + `-card.tsx",
      "when": { "architecture": ["triple"], "frontend": ["next"] }
    }
  ],
  "injections": [
    {
      "file": "{{API_ROOT}}/internal/routes/routes.go",
      "marker": "// grit:imports",
      "code": "\t\"{{MODULE}}/internal/` + pkg + `\""
    },
    {
      "file": "{{API_ROOT}}/internal/routes/routes.go",
      "marker": "// grit:routes:protected",
      "code": "\t\tprotected.GET(\"/` + name + `\", ` + pkg + `.Handler())"
    }
  ]
}
`
}

// skeletonService is the one Go file the plugin installs.
//
// It answers through the project's respond package rather than c.JSON, because
// that is the API contract every other endpoint follows and a plugin that
// writes its own envelope is the first thing a reviewer would send back.
func skeletonService(name, pkg string) string {
	return `// Package ` + pkg + ` was installed by the ` + name + ` plugin.
//
// A plugin writes code into the project rather than shipping a runtime
// dependency, so this file is now yours: edit it, and grit plugin update will
// leave it alone and tell you it did.
package ` + pkg + `

import (
	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/respond"
)

// Message is what the endpoint returns. A real plugin would read something.
type Message struct {
	Plugin  string ` + "`" + `json:"plugin"` + "`" + `
	Project string ` + "`" + `json:"project"` + "`" + `
}

// Handler serves GET /api/` + name + `.
//
// respond.OK, not c.JSON: every endpoint in a Grit project answers with the
// same envelope ({"data": ...}), and the frontend's error helper walks one
// shape. See /docs/reference/api-design.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		respond.OK(c, Message{
			Plugin:  "` + name + `",
			Project: "{{PROJECT}}",
		})
	}
}
`
}

// skeletonServiceTest is the Go test that ships with it.
//
// It runs in the installing project, so it is the check that the plugin's own
// CI loop ends with: `go test ./...` after an install either passes or the
// plugin wrote something that does not work there.
func skeletonServiceTest(pkg string) string {
	return `package ` + pkg + `

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// The endpoint answers in the project's envelope. Worth asserting rather than
// assuming: respond.OK wraps the value in {"data": ...}, and a frontend that
// reads res.data.plugin breaks if this ever answers with the bare struct.
func TestHandlerAnswersInTheEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/` + pkg + `", Handler())

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/` + pkg + `", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Data Message ` + "`" + `json:"data"` + "`" + `
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("the response is not the envelope: %v\n%s", err, rec.Body.String())
	}
	if body.Data.Plugin == "" {
		t.Error("the response has no plugin name, so the placeholder was not substituted")
	}
}
`
}

// skeletonAdminCard is the file behind the when clause.
//
// apps/admin, so the clause is architecture: triple. An embedded panel keeps
// its components under a different root and behind a different import alias
// (@admin/ rather than @/), so one file cannot serve both: the README says so
// and shows the second entry to add.
func skeletonAdminCard(name, title string) string {
	return `// Installed by the ` + name + ` plugin, into a triple project's admin panel.
//
// Drop it on a page: <` + componentName(name) + ` />
"use client";

import { useQuery } from "@tanstack/react-query";

import { apiClient } from "@/lib/api-client";

interface Message {
  plugin: string;
  project: string;
}

export function ` + componentName(name) + `() {
  const { data, isLoading } = useQuery({
    queryKey: ["` + name + `"],
    queryFn: async () => {
      const res = await apiClient.get<{ data: Message }>("/api/` + name + `");
      return res.data.data;
    },
  });

  return (
    <div className="rounded-xl border border-border bg-card p-5">
      <h3 className="text-sm font-medium text-muted-foreground">` + title + `</h3>
      <p className="mt-2 text-2xl font-semibold text-foreground">
        {isLoading ? "..." : data?.project}
      </p>
    </div>
  );
}
`
}

// componentName turns audit-trail into AuditTrailCard.
func componentName(name string) string {
	var b strings.Builder
	for _, part := range strings.Split(name, "-") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	b.WriteString("Card")
	return b.String()
}
