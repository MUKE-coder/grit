package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// writeInto is newFixtureProject's writer, for a test that adds to the fixture.
func writeInto(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

// call runs one tool directly and decodes its JSON answer.
func call(t *testing.T, s *Server, name string, args string) map[string]interface{} {
	t.Helper()
	var raw json.RawMessage
	if args != "" {
		raw = json.RawMessage(args)
	}
	text, err := s.callTool(name, raw)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("%s returned text that is not JSON: %v\n%s", name, err, text)
	}
	return out
}

func TestListResourcesFindsGeneratedOnesByTheirMarker(t *testing.T) {
	root := newFixtureProject(t)
	api := filepath.Join("apps", "api")

	// A generated resource: the model, and the bulk request only the generator
	// writes. The fixture's Note model is already there.
	writeInto(t, root, filepath.Join(api, "internal", "handlers", "note.go"),
		"package handlers\n\ntype BulkNoteRequest struct{ IDs []string }\n")
	writeInto(t, root, filepath.Join(api, "internal", "services", "note.go"),
		"package services\n\nfunc BulkDelete(req handlers.BulkNoteRequest) {}\n")
	writeInto(t, root, filepath.Join(api, "internal", "handlers", "note_import.go"),
		"package handlers\n\n// CSV import for BulkNoteRequest's resource.\n")
	writeInto(t, root, filepath.Join("packages", "shared", "schemas", "note.ts"), "export {}\n")

	// A framework table with a model, a handler and a service, and no marker.
	writeInto(t, root, filepath.Join(api, "internal", "models", "saml.go"),
		"package models\n\ntype SAMLConnection struct{ ID string }\n")
	writeInto(t, root, filepath.Join(api, "internal", "handlers", "saml.go"), "package handlers\n")
	writeInto(t, root, filepath.Join(api, "internal", "services", "saml.go"), "package services\n")

	got := call(t, &Server{Root: root}, "grit_list_resources", "")
	if got["count"] != float64(1) {
		t.Fatalf("count = %v, want 1: %v", got["count"], got["resources"])
	}
	first := got["resources"].([]interface{})[0].(map[string]interface{})
	if first["name"] != "Note" {
		t.Fatalf("found %v, want Note", first["name"])
	}
	if first["table"] != "notes" {
		t.Errorf("table = %v, want notes", first["table"])
	}
	files := first["files"].(map[string]interface{})
	for _, want := range []string{"model", "service", "handler", "schema"} {
		if files[want] == nil {
			t.Errorf("the resource does not report its %s: %v", want, files)
		}
	}
	// The CSV import belongs to the resource rather than being one of its own.
	if files["csv_import"] == nil {
		t.Errorf("the import handler was not attributed to the resource: %v", files)
	}
}

func TestFileOwnershipSaysThereIsNoManifest(t *testing.T) {
	// A project from before manifests, or one whose .grit was deleted. The
	// honest answer is that nothing can be said, not an empty list that reads
	// as "you have edited nothing".
	got := call(t, &Server{Root: newFixtureProject(t)}, "grit_file_ownership", "")
	if got["tracked"] != float64(0) {
		t.Fatalf("tracked = %v, want 0", got["tracked"])
	}
	if !strings.Contains(got["note"].(string), "no manifest") {
		t.Errorf("the note does not explain the empty answer: %q", got["note"])
	}
}

func TestFileOwnershipSeparatesPristineFromEdited(t *testing.T) {
	root := newFixtureProject(t)
	writeInto(t, root, "kept.txt", "one\n")
	writeInto(t, root, "edited.txt", "one\n")
	// The manifest records the hash of what Grit wrote. Written by hand here
	// rather than by running a scaffold, which this package does not do.
	writeInto(t, root, filepath.Join(".grit", "manifest.json"), `{
		"format": 1,
		"files": {
			"kept.txt":    {"hash": "`+sha256Of("one\n")+`"},
			"edited.txt":  {"hash": "`+sha256Of("something else\n")+`"},
			"deleted.txt": {"hash": "`+sha256Of("gone\n")+`"}
		}
	}`)

	got := call(t, &Server{Root: root}, "grit_file_ownership", "")
	counts := got["counts"].(map[string]interface{})
	for status, want := range map[string]float64{"unchanged": 1, "modified": 1, "missing": 1} {
		if counts[status] != want {
			t.Errorf("%s = %v, want %v", status, counts[status], want)
		}
	}

	only := call(t, &Server{Root: root}, "grit_file_ownership", `{"status":"modified"}`)
	files := only["files"].([]interface{})
	if len(files) != 1 || files[0].(map[string]interface{})["path"] != "edited.txt" {
		t.Errorf("the modified filter returned %v", files)
	}
}

func TestListPermissionsReadsTheCatalogueNotTheComments(t *testing.T) {
	root := newFixtureProject(t)
	writeInto(t, root, filepath.Join("apps", "api", "internal", "authz", "permissions.go"), `package authz

type Action string

const (
	ActionCreate Action = "create"
	ActionView   Action = "view"
)

var (
	AllActions = []Action{ActionCreate, ActionView, ActionEdit, ActionDelete}
	ViewOnly   = []Action{ActionView}
)

// An example in a comment: Feature{Key: "invented", Actions: AllActions} must
// not be reported as a permission, because nothing checks it.
func coreModules() []Module {
	return []Module{
		{
			Key: "access",
			Groups: []Group{
				{
					Key: "people",
					Features: []Feature{
						{Key: "users", Name: "Users", Actions: AllActions},
						{Key: "reports", Name: "Reports", Actions: ViewOnly},
						{Key: "flags", Name: "Flags", Actions: []Action{ActionView, ActionEdit}},
					},
				},
			},
		},
	}
}
`)

	got := call(t, &Server{Root: root}, "grit_list_permissions", "")
	if got["features"] != float64(3) {
		t.Fatalf("features = %v, want 3", got["features"])
	}
	if got["permissions"] != float64(4+1+2) {
		t.Fatalf("permissions = %v, want 7", got["permissions"])
	}
	body, _ := json.Marshal(got)
	for _, want := range []string{`"users.create"`, `"users.delete"`, `"reports.view"`, `"flags.edit"`} {
		if !strings.Contains(string(body), want) {
			t.Errorf("the catalogue is missing %s", want)
		}
	}
	if strings.Contains(string(body), "invented") {
		t.Error("an example in a comment was reported as a real permission")
	}
	if strings.Contains(string(body), `"reports.create"`) {
		t.Error("a view-only feature was given the actions it does not have")
	}
}

func TestEnvKeysNeverReturnAValue(t *testing.T) {
	root := newFixtureProject(t)
	writeInto(t, root, ".env.example", `# The database.
# Postgres in production, SQLite for a quick start.
DB_PROVIDER=postgres

# Signing key for access tokens.
JWT_SECRET=change-me

# MAIL_MAILER picks how mail is sent.
# MAIL_MAILER=
`)
	writeInto(t, root, ".env", "DB_PROVIDER=postgres\nJWT_SECRET=hunter2-the-real-one\n")

	got := call(t, &Server{Root: root}, "grit_env_keys", "")
	body, _ := json.Marshal(got)
	if strings.Contains(string(body), "hunter2") {
		t.Fatal("a value from .env reached the answer")
	}
	if strings.Contains(string(body), "change-me") {
		t.Fatal("a value from .env.example reached the answer")
	}
	if got["count"] != float64(3) {
		t.Errorf("count = %v, want 3 (a commented-out key is still a key)", got["count"])
	}

	keys := map[string]map[string]interface{}{}
	for _, k := range got["keys"].([]interface{}) {
		entry := k.(map[string]interface{})
		keys[entry["key"].(string)] = entry
	}
	if keys["JWT_SECRET"]["set_in_env"] != true {
		t.Error("JWT_SECRET is set in .env and was not reported as set")
	}
	if keys["MAIL_MAILER"]["set_in_env"] != false {
		t.Error("MAIL_MAILER is not in .env and was reported as set")
	}
	if doc, _ := keys["DB_PROVIDER"]["doc"].(string); !strings.Contains(doc, "Postgres in production") {
		t.Errorf("the comment above a key was not carried through: %q", doc)
	}
}

func TestEnvKeysFilters(t *testing.T) {
	root := newFixtureProject(t)
	writeInto(t, root, ".env.example", "DB_PROVIDER=postgres\nMAIL_FROM=noreply@test\nMAIL_MAILER=log\n")
	got := call(t, &Server{Root: root}, "grit_env_keys", `{"contains":"mail"}`)
	if got["count"] != float64(2) {
		t.Errorf("count = %v, want 2", got["count"])
	}
}

func TestCLIReferenceComesFromTheRunningBinary(t *testing.T) {
	s := &Server{Root: newFixtureProject(t), Version: "3.351.0", Commands: []CommandInfo{
		{Path: "generate", Short: "Generate code"},
		{Path: "generate resource", Short: "Generate a full resource", Flags: []FlagInfo{
			{Name: "--fields", Type: "string", Usage: "Inline field definitions"},
		}},
		{Path: "migrate", Short: "Run database migrations"},
	}}

	all := call(t, s, "grit_cli_reference", "")
	if all["count"] != float64(3) || all["grit_version"] != "3.351.0" {
		t.Fatalf("got %v", all)
	}

	one := call(t, s, "grit_cli_reference", `{"command":"generate resource"}`)
	if one["count"] != float64(1) {
		t.Fatalf("count = %v, want 1", one["count"])
	}
	body, _ := json.Marshal(one)
	if !strings.Contains(string(body), "--fields") {
		t.Error("the flags did not come through")
	}

	// A command this version does not have is an error, not an empty list: an
	// agent should be told the flag it remembered is gone.
	if _, err := s.callTool("grit_cli_reference", json.RawMessage(`{"command":"deploy to mars"}`)); err == nil {
		t.Error("an unknown command was answered with an empty list")
	}
}

func TestCLIReferenceSaysWhenItWasNotWired(t *testing.T) {
	// The mcp package cannot read cobra's tree itself; the command layer hands
	// it in. A server built without it should say so rather than answer "no
	// commands", which reads as a CLI with nothing in it.
	if _, err := (&Server{Root: t.TempDir()}).callTool("grit_cli_reference", nil); err == nil {
		t.Error("a server with no command list answered as though there were none")
	}
}

func TestGenerateResourceChecksItsArgumentsBeforeRunningAnything(t *testing.T) {
	s := &Server{Root: newFixtureProject(t), Mode: ModeWrite}
	for name, args := range map[string]string{
		"a name with a shell metacharacter": `{"name":"Product; rm -rf /"}`,
		"a name with a path separator":      `{"name":"../../etc/passwd"}`,
		"an empty name":                     `{"name":""}`,
		"a field list with a quote":         `{"name":"Product","fields":"title:string\"}`,
		"a field type that does not exist":  `{"name":"Product","fields":"title:unicorn"}`,
		"an owned_by that is not a field":   `{"name":"Product","owned_by":"user; whoami"}`,
	} {
		if _, err := s.callTool("grit_generate_resource", json.RawMessage(args)); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

func sha256Of(content string) string {
	return manifestHash(content)
}

// manifestHash is manifest.Hash, so the fixture records what Grit would.
func manifestHash(content string) string { return manifest.Hash(content) }
