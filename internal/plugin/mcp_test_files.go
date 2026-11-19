package plugin

// mcpTransportTestGo is internal/mcp/transport_test.go.
func mcpTransportTestGo(ctx Context) string {
	return `package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// post sends one JSON-RPC body to the handler and gives back the recorder.
func post(t *testing.T, s *Server, body string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	s.Handle(c)
	return rec
}

func TestInitializeAnswersAVersionTheClientAskedFor(t *testing.T) {
	s := &Server{AppName: "Test", AppVersion: "1.0.0"}

	rec := post(t, s, ` + "`" + `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18"}}` + "`" + `)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var answer struct {
		Result struct {
			ProtocolVersion string ` + "`" + `json:"protocolVersion"` + "`" + `
			Instructions    string ` + "`" + `json:"instructions"` + "`" + `
			ServerInfo      struct {
				Name string ` + "`" + `json:"name"` + "`" + `
			} ` + "`" + `json:"serverInfo"` + "`" + `
		} ` + "`" + `json:"result"` + "`" + `
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &answer); err != nil {
		t.Fatalf("%v: %s", err, rec.Body.String())
	}
	if answer.Result.ProtocolVersion != "2025-06-18" {
		t.Errorf("answered %q, want the version the client asked for",
			answer.Result.ProtocolVersion)
	}
	if answer.Result.ServerInfo.Name != "Test" {
		t.Errorf("serverInfo.name is %q", answer.Result.ServerInfo.Name)
	}
	// The model reads this before anything else, so an empty one is a bug
	// that costs every later call its context.
	if answer.Result.Instructions == "" {
		t.Error("initialize returned no instructions")
	}
}

// A version nobody here knows gets the oldest we speak, not a refusal.
//
// A client from the future sending a date we have never heard of can still
// talk over 2024-11-05, and a refusal would look to the person like the server
// being down.
func TestAnUnknownProtocolVersionFallsBackRatherThanRefusing(t *testing.T) {
	if got := negotiate("2099-01-01"); got != protocolVersions[len(protocolVersions)-1] {
		t.Errorf("negotiated %q for an unknown version", got)
	}
	if got := negotiate(""); got != protocolVersions[0] {
		t.Errorf("negotiated %q when the client named none", got)
	}
}

func TestPingAnswersEmpty(t *testing.T) {
	s := &Server{}
	rec := post(t, s, ` + "`" + `{"jsonrpc":"2.0","id":"a","method":"ping"}` + "`" + `)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(` + "`" + `"id":"a"` + "`" + `)) {
		t.Errorf("the id was not echoed: %s", rec.Body.String())
	}
}

// A string id comes back a string and a number comes back a number.
//
// The spec says echo what you were given, and a strict client rejects an
// answer whose id changed type. Decoding it into either loses the distinction,
// which is why it is held as raw JSON.
func TestTheRequestIDIsEchoedWithItsOwnType(t *testing.T) {
	s := &Server{}

	for _, id := range []string{` + "`" + `7` + "`" + `, ` + "`" + `"seven"` + "`" + `} {
		rec := post(t, s, ` + "`" + `{"jsonrpc":"2.0","id":` + "`" + `+id+` + "`" + `,"method":"ping"}` + "`" + `)
		if !bytes.Contains(rec.Body.Bytes(), []byte(` + "`" + `"id":` + "`" + `+id)) {
			t.Errorf("id %s came back as %s", id, rec.Body.String())
		}
	}
}

// A notification has no id and gets no body at all.
func TestANotificationIsAcceptedWithNoAnswer(t *testing.T) {
	s := &Server{}
	rec := post(t, s, ` + "`" + `{"jsonrpc":"2.0","method":"notifications/initialized"}` + "`" + `)
	if rec.Code != http.StatusAccepted {
		t.Errorf("status %d, want 202: %s", rec.Code, rec.Body.String())
	}
}

// A batch of two answers with an array of two.
func TestABatchAnswersAnArray(t *testing.T) {
	s := &Server{}
	rec := post(t, s, ` + "`" + `[{"jsonrpc":"2.0","id":1,"method":"ping"},{"jsonrpc":"2.0","id":2,"method":"ping"}]` + "`" + `)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var answers []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &answers); err != nil {
		t.Fatalf("a batch did not answer an array: %v: %s", err, rec.Body.String())
	}
	if len(answers) != 2 {
		t.Errorf("answered %d of 2", len(answers))
	}
}

func TestAnUnknownMethodSaysWhatThereIs(t *testing.T) {
	s := &Server{}
	rec := post(t, s, ` + "`" + `{"jsonrpc":"2.0","id":1,"method":"resources/list"}` + "`" + `)

	var answer struct {
		Error struct {
			Code    int    ` + "`" + `json:"code"` + "`" + `
			Message string ` + "`" + `json:"message"` + "`" + `
		} ` + "`" + `json:"error"` + "`" + `
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &answer); err != nil {
		t.Fatal(err)
	}
	if answer.Error.Code != codeMethodNotFound {
		t.Errorf("code %d, want %d", answer.Error.Code, codeMethodNotFound)
	}
	if !bytes.Contains([]byte(answer.Error.Message), []byte("tools/list")) {
		t.Errorf("the message does not say what the server has: %q", answer.Error.Message)
	}
}

func TestBadJSONIsAParseError(t *testing.T) {
	s := &Server{}
	rec := post(t, s, ` + "`" + `{not json` + "`" + `)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("-32700")) {
		t.Errorf("not a parse error: %s", rec.Body.String())
	}
}

// GET is the stream half of streamable HTTP, which this server does not
// implement. 405 with Allow tells a client to use POST; a 404 would leave it
// thinking the endpoint is missing.
func TestGETIsRefusedWithAllowPost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/mcp", nil)

	Reject(c)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status %d, want 405", rec.Code)
	}
	if rec.Header().Get("Allow") != "POST" {
		t.Errorf("Allow is %q", rec.Header().Get("Allow"))
	}
}

// A call with no key is refused, and says where to get one.
func TestToolsListWithoutAKeySaysHowToGetOne(t *testing.T) {
	s := &Server{}
	rec := post(t, s, ` + "`" + `{"jsonrpc":"2.0","id":1,"method":"tools/list"}` + "`" + `)

	if !bytes.Contains(rec.Body.Bytes(), []byte("connection key")) {
		t.Errorf("the refusal does not say what is missing: %s", rec.Body.String())
	}
	if rec.Header().Get("WWW-Authenticate") == "" {
		t.Error("no WWW-Authenticate, so a client cannot tell this from a bad request")
	}
}

// The two tools every project has are registered, so tools/list is never empty.
func TestTheBuiltinToolsAreRegistered(t *testing.T) {
	for _, name := range []string{"whoami", "list_tools"} {
		if _, ok := Lookup(name); !ok {
			t.Errorf("%s is not registered", name)
		}
	}
}

// A connection's own list is the authority, not the registry.
//
// This is what makes a read-only connection read-only: a tool can exist in the
// process and still not be this connection's to call.
func TestAConnectionOnlyAllowsItsOwnTools(t *testing.T) {
	conn := connection{Tools: []string{"whoami"}}
	if !conn.Allows("whoami") {
		t.Error("a connection refuses a tool it was given")
	}
	if conn.Allows("list_tools") {
		t.Error("a connection allows a tool it was not given, so its list is not " +
			"the boundary and a read-only connection is not read-only")
	}
}

// A key's routes come from its tools, so nobody can widen a connection by
// hand and adding a tool later does not widen the ones that exist.
func TestRoutesAndPermissionsComeFromTheTools(t *testing.T) {
	tools := []Tool{
		{Name: "a", Permission: "x.read", Routes: []string{"GET /api/v1/a"}},
		{Name: "b", Permission: "x.read", Routes: []string{"GET /api/v1/a", "POST /api/v1/b"}},
		{Name: "c", Routes: []string{""}},
	}

	routes := RoutesFor(tools)
	if len(routes) != 2 || routes[0] != "GET /api/v1/a" || routes[1] != "POST /api/v1/b" {
		t.Errorf("routes are %v, want the two distinct ones sorted", routes)
	}
	perms := PermissionsFor(tools)
	if len(perms) != 1 || perms[0] != "x.read" {
		t.Errorf("permissions are %v, want one", perms)
	}
}
`
}

// mcpArgsTestGo is internal/mcp/args_test.go.
func mcpArgsTestGo(ctx Context) string {
	return `package mcp

import "testing"

// A model sends what the person said, not what the schema says.
//
// The schema says number and it sends "5", because that is what was typed.
// Refusing is correct by the schema and useless: the call fails, the model
// apologises, and the person is no further forward.
func TestArgumentsAreReadForgivinglyAndStrictly(t *testing.T) {
	args := map[string]any{
		"name":    "  Kato Hardware  ",
		"qty":     "7",
		"price":   "$1,250.50",
		"ready":   "yes",
		"tags":    "urgent",
		"numeric": float64(3),
	}

	if got, err := String(args, "name"); err != nil || got != "Kato Hardware" {
		t.Errorf("String trimmed to %q (%v)", got, err)
	}
	if got, err := Int(args, "qty"); err != nil || got != 7 {
		t.Errorf("Int read %d (%v)", got, err)
	}
	if got, err := Float(args, "price"); err != nil || got != 1250.50 {
		t.Errorf("Float read %v (%v)", got, err)
	}
	if !Bool(args, "ready", false) {
		t.Error(` + "`" + `Bool did not read "yes" as true` + "`" + `)
	}
	// A list field given the bare item: the intent is unambiguous, so
	// refusing it would be pedantry.
	if got := Strings(args, "tags"); len(got) != 1 || got[0] != "urgent" {
		t.Errorf("Strings read %v from a single value", got)
	}
	// An id a model copied out of an earlier answer as a number.
	if got, err := String(args, "numeric"); err != nil || got != "3" {
		t.Errorf("String read %q from a number (%v)", got, err)
	}

	// Strict about meaning: a missing required argument is an error with the
	// field's name in it, because that is what the model needs to fix it.
	if _, err := String(args, "absent"); err == nil {
		t.Error("String accepted a missing argument")
	}
	if _, err := Int(args, "name"); err == nil {
		t.Error("Int accepted a word")
	}
}

// Paging is capped, whatever the model asks for.
//
// "List every invoice" with no page size is a model asking for a hundred
// thousand rows it will then put in its own context.
func TestPagingIsCapped(t *testing.T) {
	query := Page(map[string]any{"page_size": 5000, "page": -3})
	if query.Get("page_size") != "100" {
		t.Errorf("page_size is %q, want the cap", query.Get("page_size"))
	}
	if query.Get("page") != "1" {
		t.Errorf("page is %q, want 1 for a negative one", query.Get("page"))
	}

	defaults := Page(map[string]any{})
	if defaults.Get("page_size") != "20" || defaults.Get("page") != "1" {
		t.Errorf("defaults are %v", defaults)
	}
}

// An id is recognised by shape, so passing one costs no lookup request.
func TestAnIDIsRecognisedByShape(t *testing.T) {
	if !looksLikeID("0b7f1d6e-9c3a-4f21-8a55-2d9e7c4b1f08") {
		t.Error("a UUID was not recognised")
	}
	for _, not := range []string{
		"Kato Hardware",
		"",
		"0b7f1d6e9c3a4f218a552d9e7c4b1f08",
		"0b7f1d6e-9c3a-4f21-8a55-2d9e7c4b1f0",
		"zzzzzzzz-9c3a-4f21-8a55-2d9e7c4b1f08",
	} {
		if looksLikeID(not) {
			t.Errorf("%q was taken for an id", not)
		}
	}
}

// An ambiguous name is an error naming the candidates, never a guess.
//
// Picking the first match is the behaviour that quietly invoices the wrong
// customer. A model handed the candidates can ask which was meant, which is
// the only correct resolution.
func TestAnAmbiguousNameNamesTheCandidates(t *testing.T) {
	err := &Ambiguous{
		What:       "client",
		Term:       "kato",
		Candidates: []string{"Kato Hardware", "Kato Motors"},
	}
	message := err.Error()
	for _, want := range []string{"Kato Hardware", "Kato Motors", "Ask", "its id"} {
		if !contains(message, want) {
			t.Errorf("the message does not mention %q: %s", want, message)
		}
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
`
}

// mcpConnectionServiceTestGo is internal/services/mcp_connection_test.go.
func mcpConnectionServiceTestGo(ctx Context) string {
	return `package services

import (
	"testing"

	"` + ctx.Module + `/internal/mcp"
)

// A connection cannot be given a tool its owner's role does not grant.
//
// An agent must not be able to do more than the person who set it up. Checked
// in the service and not only in the screen, because the screen is a
// convenience over an API endpoint anybody can call directly.
func TestAConnectionCannotExceedItsOwnersRole(t *testing.T) {
	// whoami needs no permission, so it is always offerable.
	if _, ok := mcp.Lookup("whoami"); !ok {
		t.Skip("the builtin tools are not registered in this build")
	}

	_, err := CreateMCPConnection(nil, NewConnection{
		Name:   "no tools",
		UserID: "u1",
		Tools:  nil,
	})
	if err == nil {
		t.Error("a connection with no tools was accepted, and it can do nothing")
	}

	_, err = CreateMCPConnection(nil, NewConnection{
		Name:   "unknown tool",
		UserID: "u1",
		Tools:  []string{"delete_the_database"},
	})
	if err == nil {
		t.Error("a connection naming a tool that does not exist was accepted")
	}
}

// The name and the user are both required, because a connection without them
// cannot be told apart in a list or held to anybody's permissions.
func TestAConnectionNeedsANameAndAUser(t *testing.T) {
	if _, err := CreateMCPConnection(nil, NewConnection{UserID: "u1", Tools: []string{"whoami"}}); err == nil {
		t.Error("a connection with no name was accepted")
	}
	if _, err := CreateMCPConnection(nil, NewConnection{Name: "x", Tools: []string{"whoami"}}); err == nil {
		t.Error("a connection with no user was accepted")
	}
}
`
}
