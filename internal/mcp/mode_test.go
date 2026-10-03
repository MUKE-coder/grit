package mcp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseMode(t *testing.T) {
	for in, want := range map[string]Mode{"": ModeRead, "read": ModeRead, "write": ModeWrite} {
		got, err := ParseMode(in)
		if err != nil || got != want {
			t.Errorf("ParseMode(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := ParseMode("readonly"); err == nil {
		t.Error("an unknown mode should be refused rather than quietly treated as read")
	}
}

func TestTheZeroServerIsReadOnly(t *testing.T) {
	// A Server built without thinking about the mode is the one that cannot
	// write. The default has to be the safe one, because the unsafe one is
	// reached by asking for it.
	s := &Server{Root: t.TempDir()}
	if s.mode() != ModeRead {
		t.Fatalf("the zero Mode is %q, want read", s.mode())
	}
}

// mutatingToolNames is every tool that writes, taken from the one list, so a
// tool added later is covered by these tests without anybody editing them.
func mutatingToolNames() []string {
	var out []string
	for _, t := range allTools() {
		if t.mutates {
			out = append(out, t.name)
		}
	}
	return out
}

func TestThereIsAtLeastOneMutatingTool(t *testing.T) {
	// Otherwise the two tests below pass by having nothing to check, which is
	// the way this kind of guarantee quietly stops being one.
	if len(mutatingToolNames()) == 0 {
		t.Fatal("no tool is marked mutates, so the mode split is not being tested")
	}
}

func TestAReadOnlyServerDoesNotContainTheWritingTools(t *testing.T) {
	s := &Server{Root: t.TempDir(), Mode: ModeRead}

	listed := map[string]bool{}
	for _, def := range s.toolDefinitions() {
		listed[def["name"].(string)] = true
	}
	for _, name := range mutatingToolNames() {
		if listed[name] {
			t.Errorf("%s is listed by a read-only server", name)
		}
		if _, ok := s.tools()[name]; ok {
			t.Errorf("%s is registered on a read-only server, so something can still call it", name)
		}
	}
	if len(listed) == 0 {
		t.Fatal("a read-only server listed no tools at all")
	}
}

func TestCallingAWritingToolOnAReadOnlyServerSaysWhy(t *testing.T) {
	s := &Server{Root: t.TempDir(), Mode: ModeRead}
	for _, name := range mutatingToolNames() {
		_, err := s.callTool(name, nil)
		if err == nil {
			t.Fatalf("%s ran on a read-only server", name)
		}
		// Not "unknown tool": a model told that would conclude the capability
		// does not exist and hand-write the files the generator owns.
		if !strings.Contains(err.Error(), "--mode write") {
			t.Errorf("%s: the refusal does not say how to allow it: %v", name, err)
		}
	}
}

func TestAWriteServerHasEverything(t *testing.T) {
	read := &Server{Root: t.TempDir(), Mode: ModeRead}
	write := &Server{Root: t.TempDir(), Mode: ModeWrite}

	if len(write.tools()) != len(allTools()) {
		t.Errorf("a write server has %d of %d tools", len(write.tools()), len(allTools()))
	}
	if len(write.tools()) <= len(read.tools()) {
		t.Error("write mode added nothing, so the split is not doing anything")
	}
	// Everything a read server has, a write server also has: write is read plus
	// the generators, not a different server.
	for name := range read.tools() {
		if _, ok := write.tools()[name]; !ok {
			t.Errorf("%s is missing from a write server", name)
		}
	}
}

func TestListAndCallAgreeAboutWhatExists(t *testing.T) {
	// The bug this prevents: a tool listed but not dispatchable, or dispatchable
	// but not listed. Both come from two places deciding separately.
	for _, mode := range []Mode{ModeRead, ModeWrite} {
		s := &Server{Root: t.TempDir(), Mode: mode}
		set := s.tools()
		for _, def := range s.toolDefinitions() {
			if _, ok := set[def["name"].(string)]; !ok {
				t.Errorf("%s mode lists %s and cannot call it", mode, def["name"])
			}
		}
		if len(set) != len(s.toolDefinitions()) {
			t.Errorf("%s mode: %d callable, %d listed", mode, len(set), len(s.toolDefinitions()))
		}
	}
}

func TestEveryToolIsDescribedAndSchemad(t *testing.T) {
	for _, tool := range allTools() {
		if !strings.HasPrefix(tool.name, "grit_") {
			t.Errorf("%s does not carry the grit_ prefix, so it collides with other servers' tools", tool.name)
		}
		if len(tool.description) < 80 {
			t.Errorf("%s has a description too short to help a model choose: %q", tool.name, tool.description)
		}
		if tool.schema == nil || tool.schema["type"] != "object" {
			t.Errorf("%s has no object input schema", tool.name)
		}
		if tool.run == nil {
			t.Errorf("%s has no handler", tool.name)
		}
	}
}

func TestToolsListOverTheProtocolRespectsTheMode(t *testing.T) {
	// Through handle, not through the helper, because that is the path a client
	// takes.
	for _, tc := range []struct {
		mode     Mode
		expected bool
	}{{ModeRead, false}, {ModeWrite, true}} {
		s := &Server{Root: t.TempDir(), Mode: tc.mode}
		resp, ok := s.handle(request{JSONRPC: jsonrpcVersion, ID: json.RawMessage(`1`), Method: "tools/list"})
		if !ok {
			t.Fatal("tools/list got no reply")
		}
		body, err := json.Marshal(resp.Result)
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range mutatingToolNames() {
			if strings.Contains(string(body), name) != tc.expected {
				t.Errorf("%s mode: %s present = %v, want %v", tc.mode, name, !tc.expected, tc.expected)
			}
		}
	}
}
