package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newFixtureProject builds a minimal but realistic monorepo-shaped project:
// grit.json at the root, the Go module under apps/api, a routes.go that mounts
// everything under a constant-built prefix, and one model.
func newFixtureProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api")

	write := func(rel, content string) {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	write("grit.json", `{"architecture":"triple","frontend":"next","version":"3.109.0"}`)
	write(filepath.Join("apps", "api", "go.mod"), "module fixture/apps/api\n\ngo 1.24\n")
	write(filepath.Join("apps", "api", "internal", "routes", "routes.go"), `package routes

const APIVersion = "v1"

func Setup(r *gin.Engine) {
	v1 := r.Group("/api/" + APIVersion)
	v1.GET("/notes", noteHandler.List)
	admin := v1.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.DELETE("/notes/:id", noteHandler.Delete)
}
`)
	write(filepath.Join("apps", "api", "internal", "models", "note.go"), `package models

import "time"

type Note struct {
	ID        string    `+"`"+`gorm:"primarykey;size:36" json:"id"`+"`"+`
	Title     string    `+"`"+`gorm:"size:255" json:"title"`+"`"+`
	CreatedAt time.Time `+"`"+`json:"created_at"`+"`"+`
}
`)

	_ = api
	return root
}

// run feeds lines to the server and returns the decoded responses.
func run(t *testing.T, root string, lines ...string) []map[string]interface{} {
	t.Helper()
	srv := &Server{Root: root, Version: "test"}

	var out strings.Builder
	if err := srv.Serve(strings.NewReader(strings.Join(lines, "\n")+"\n"), &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}

	var got []map[string]interface{}
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("server emitted a non-JSON line %q: %v", line, err)
		}
		got = append(got, m)
	}
	return got
}

// toolText pulls the text payload out of a tools/call result.
func toolText(t *testing.T, resp map[string]interface{}) (string, bool) {
	t.Helper()
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("response has no result: %v", resp)
	}
	isError, _ := result["isError"].(bool)
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("result has no content: %v", result)
	}
	first, _ := content[0].(map[string]interface{})
	text, _ := first["text"].(string)
	return text, isError
}

func TestInitializeReportsCapabilitiesAndServerInfo(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18"}}`)

	if len(got) != 1 {
		t.Fatalf("expected 1 response, got %d", len(got))
	}
	result := got[0]["result"].(map[string]interface{})
	if result["protocolVersion"] != "2025-06-18" {
		t.Errorf("protocolVersion = %v, want the version the client asked for", result["protocolVersion"])
	}
	if _, ok := result["capabilities"].(map[string]interface{})["tools"]; !ok {
		t.Error("server did not advertise the tools capability")
	}
	if info := result["serverInfo"].(map[string]interface{}); info["name"] != "grit" {
		t.Errorf("serverInfo.name = %v, want grit", info["name"])
	}
}

// An unrecognised protocol revision must fall back to one we actually
// implement, rather than echoing something we cannot honour.
func TestInitializeFallsBackForAnUnknownProtocolVersion(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"1999-01-01"}}`)

	result := got[0]["result"].(map[string]interface{})
	if result["protocolVersion"] != defaultProtocolVersion {
		t.Errorf("protocolVersion = %v, want %v", result["protocolVersion"], defaultProtocolVersion)
	}
}

// Notifications carry no id and must draw no reply. Answering one is a protocol
// violation that some clients treat as fatal.
func TestNotificationsGetNoResponse(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":7,"method":"ping"}`)

	if len(got) != 1 {
		t.Fatalf("expected exactly 1 response (the ping), got %d: %v", len(got), got)
	}
	if got[0]["id"].(float64) != 7 {
		t.Errorf("the reply was not for the ping: %v", got[0])
	}
}

func TestToolsListAdvertisesEveryToolWithASchema(t *testing.T) {
	got := run(t, newFixtureProject(t), `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)

	tools := got[0]["result"].(map[string]interface{})["tools"].([]interface{})
	seen := map[string]bool{}
	for _, raw := range tools {
		tool := raw.(map[string]interface{})
		name, _ := tool["name"].(string)
		seen[name] = true

		if desc, _ := tool["description"].(string); desc == "" {
			t.Errorf("%s has no description — a model cannot choose it", name)
		}
		schema, ok := tool["inputSchema"].(map[string]interface{})
		if !ok || schema["type"] != "object" {
			t.Errorf("%s has no object inputSchema", name)
		}
	}

	for _, want := range []string{"grit_project_info", "grit_list_routes", "grit_describe_models"} {
		if !seen[want] {
			t.Errorf("tool %s not advertised", want)
		}
	}
}

func TestListRoutesReportsTheVersionedPathAndAccessLevel(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grit_list_routes","arguments":{}}}`)

	text, isErr := toolText(t, got[0])
	if isErr {
		t.Fatalf("tool reported an error: %s", text)
	}
	// The prefix comes from `"/api/" + APIVersion`; without constant resolution
	// this silently reads /notes and sends an agent to a 404.
	if !strings.Contains(text, `"/api/v1/notes"`) {
		t.Errorf("route missing its versioned prefix:\n%s", text)
	}
	if !strings.Contains(text, `"admin"`) {
		t.Errorf("admin access level not reported:\n%s", text)
	}
}

func TestListRoutesFilters(t *testing.T) {
	root := newFixtureProject(t)

	got := run(t, root,
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grit_list_routes","arguments":{"method":"delete"}}}`)
	text, _ := toolText(t, got[0])
	if strings.Contains(text, `"GET"`) {
		t.Errorf("method filter did not exclude GET:\n%s", text)
	}
	if !strings.Contains(text, `"DELETE"`) {
		t.Errorf("method filter is case-sensitive; lowercase input dropped everything:\n%s", text)
	}
}

func TestDescribeModelsReportsFieldsAndTags(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grit_describe_models","arguments":{"model":"Note"}}}`)

	text, isErr := toolText(t, got[0])
	if isErr {
		t.Fatalf("tool reported an error: %s", text)
	}
	for _, want := range []string{`"Note"`, `"Title"`, `"json_name": "title"`, `primarykey;size:36`} {
		if !strings.Contains(text, want) {
			t.Errorf("model output missing %s:\n%s", want, text)
		}
	}
}

func TestProjectInfoReportsArchitectureAndModule(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grit_project_info","arguments":{}}}`)

	text, isErr := toolText(t, got[0])
	if isErr {
		t.Fatalf("tool reported an error: %s", text)
	}
	for _, want := range []string{`"architecture": "triple"`, `"go_module": "fixture/apps/api"`, `"api_dir": "apps/api"`} {
		if !strings.Contains(text, want) {
			t.Errorf("project info missing %s:\n%s", want, text)
		}
	}
}

// A failing tool reports through the result with isError set, so the agent sees
// the message and can try something else. Returning a JSON-RPC error would make
// it look like a transport fault.
func TestToolFailureIsReportedInTheResultNotAsAProtocolError(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"grit_describe_models","arguments":{"model":"Nope"}}}`)

	if _, hasErr := got[0]["error"]; hasErr {
		t.Errorf("tool failure surfaced as a protocol error: %v", got[0])
	}
	text, isErr := toolText(t, got[0])
	if !isErr {
		t.Error("isError was not set on a failed call")
	}
	if !strings.Contains(text, "Nope") {
		t.Errorf("error text does not name the missing model: %s", text)
	}
}

func TestUnknownToolAndUnknownMethod(t *testing.T) {
	root := newFixtureProject(t)

	got := run(t, root, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"nope","arguments":{}}}`)
	if text, isErr := toolText(t, got[0]); !isErr || !strings.Contains(text, "nope") {
		t.Errorf("unknown tool not reported: %v", got[0])
	}

	got = run(t, root, `{"jsonrpc":"2.0","id":2,"method":"totally/unknown"}`)
	errObj, ok := got[0]["error"].(map[string]interface{})
	if !ok || errObj["code"].(float64) != codeMethodNotFound {
		t.Errorf("unknown method should be %d: %v", codeMethodNotFound, got[0])
	}
}

// A malformed message must be reported and the session must continue — one bad
// line should not end the conversation.
func TestMalformedMessageDoesNotKillTheSession(t *testing.T) {
	got := run(t, newFixtureProject(t),
		`{"jsonrpc":"2.0","id":1,"method":"ping"}`,
		`{ this is not json`,
		`{"jsonrpc":"2.0","id":3,"method":"ping"}`)

	if len(got) != 3 {
		t.Fatalf("expected 3 responses (ping, parse error, ping), got %d: %v", len(got), got)
	}
	errObj, ok := got[1]["error"].(map[string]interface{})
	if !ok || errObj["code"].(float64) != codeParseError {
		t.Errorf("second response should be a parse error: %v", got[1])
	}
	if got[2]["id"].(float64) != 3 {
		t.Errorf("session did not continue after the bad line: %v", got[2])
	}
}

// stdout is the protocol channel. Anything else on it desynchronises the client,
// so every emitted line must be a complete JSON object.
func TestStdoutCarriesOnlyProtocolMessages(t *testing.T) {
	srv := &Server{Root: newFixtureProject(t), Version: "test"}
	var out strings.Builder
	in := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"grit_list_routes","arguments":{}}}`,
	}, "\n") + "\n"

	if err := srv.Serve(strings.NewReader(in), &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}

	for i, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Errorf("line %d is not a JSON object: %q", i, line)
			continue
		}
		if m["jsonrpc"] != "2.0" {
			t.Errorf("line %d is not a JSON-RPC message: %q", i, line)
		}
	}
}
