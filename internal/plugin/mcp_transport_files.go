package plugin

// mcpTransportGo is internal/mcp/transport.go: the protocol.
func mcpTransportGo(ctx Context) string {
	return `package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// The protocol versions this server speaks, newest first.
//
// A client sends the one it wants in initialize; we answer with the newest we
// both know. Answering with a version the client did not offer is how a
// session starts and then fails on the first call it cannot parse.
var protocolVersions = []string{
	"2025-11-25",
	"2025-06-18",
	"2025-03-26",
	"2024-11-05",
}

// rpcRequest is one JSON-RPC 2.0 call.
//
// ID is json.RawMessage because the spec allows a string or a number and says
// to echo back what you were given. Decoding it into either loses the
// distinction and a strict client rejects the answer.
type rpcRequest struct {
	JSONRPC string          ` + "`" + `json:"jsonrpc"` + "`" + `
	ID      json.RawMessage ` + "`" + `json:"id,omitempty"` + "`" + `
	Method  string          ` + "`" + `json:"method"` + "`" + `
	Params  json.RawMessage ` + "`" + `json:"params,omitempty"` + "`" + `
}

// isNotification reports a request with no id, which by the spec gets no
// answer at all.
func (r rpcRequest) isNotification() bool {
	return len(r.ID) == 0 || string(r.ID) == "null"
}

type rpcError struct {
	Code    int    ` + "`" + `json:"code"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
	Data    any    ` + "`" + `json:"data,omitempty"` + "`" + `
}

type rpcResponse struct {
	JSONRPC string          ` + "`" + `json:"jsonrpc"` + "`" + `
	ID      json.RawMessage ` + "`" + `json:"id,omitempty"` + "`" + `
	Result  any             ` + "`" + `json:"result,omitempty"` + "`" + `
	Error   *rpcError       ` + "`" + `json:"error,omitempty"` + "`" + `
}

// The JSON-RPC error codes, which are the spec's and not ours.
const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
	codeInternalError  = -32603
)

func ok(id json.RawMessage, result any) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", ID: id, Result: result}
}

func fail(id json.RawMessage, code int, message string) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: message}}
}

// Handle is POST /mcp.
//
// Streamable HTTP, and deliberately the plain half of it: every request gets a
// JSON answer rather than an event stream. The transport allows that, and no
// tool here reports progress, so a stream would be a second code path with
// nothing to carry. If a tool ever needs progress, this is where it changes.
//
// Stateless. Every request carries its own key, so there is no session to
// create, lose, or pin to one process: the server can be behind a load
// balancer and scaled to zero without a client noticing.
func (s *Server) Handle(c *gin.Context) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodyBytes))
	if err != nil {
		writeOne(c, http.StatusBadRequest, fail(nil, codeParseError, "the request body could not be read"))
		return
	}

	// A batch is an array, a single call is an object. The spec allows both and
	// a client that sends a batch of two expects an array of two back, so the
	// shape of the answer follows the shape of the request.
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		writeOne(c, http.StatusBadRequest, fail(nil, codeInvalidRequest, "empty request"))
		return
	}

	if trimmed[0] == '[' {
		var batch []rpcRequest
		if err := json.Unmarshal(trimmed, &batch); err != nil {
			writeOne(c, http.StatusBadRequest, fail(nil, codeParseError, "invalid JSON"))
			return
		}
		if len(batch) == 0 {
			writeOne(c, http.StatusBadRequest, fail(nil, codeInvalidRequest, "empty batch"))
			return
		}
		var answers []rpcResponse
		for _, req := range batch {
			if answer, answered := s.dispatch(c, req); answered {
				answers = append(answers, answer)
			}
		}
		if len(answers) == 0 {
			// Every member was a notification. Accepted, nothing to say.
			accepted(c)
			return
		}
		c.JSON(http.StatusOK, answers)
		return
	}

	var req rpcRequest
	if err := json.Unmarshal(trimmed, &req); err != nil {
		writeOne(c, http.StatusBadRequest, fail(nil, codeParseError, "invalid JSON"))
		return
	}
	answer, answered := s.dispatch(c, req)
	if !answered {
		accepted(c)
		return
	}
	writeOne(c, http.StatusOK, answer)
}

// maxBodyBytes caps a request. A tools/call argument object is small; anything
// approaching this is a mistake or an attack, and reading it to find out is
// the cost either way.
const maxBodyBytes = 1 << 20

func writeOne(c *gin.Context, status int, answer rpcResponse) {
	c.JSON(status, answer)
}

// accepted is 202 with no body, for a notification.
//
// WriteHeaderNow as well as Status, because Status only records the code on
// gin's writer and the write happens when a body is written or the chain
// flushes. With no body there is nothing to trigger it, so a handler called
// directly answers 200 and a client that checks the status is told its
// notification was a normal response.
func accepted(c *gin.Context) {
	c.Status(http.StatusAccepted)
	c.Writer.WriteHeaderNow()
}

// dispatch answers one request, and reports whether there is anything to send.
func (s *Server) dispatch(c *gin.Context, req rpcRequest) (rpcResponse, bool) {
	if req.JSONRPC != "" && req.JSONRPC != "2.0" {
		return fail(req.ID, codeInvalidRequest, "this server speaks JSON-RPC 2.0"), !req.isNotification()
	}

	// A notification is answered by doing the work and saying nothing. Only
	// the handshake ones arrive in practice, and they need no work at all.
	if req.isNotification() {
		return rpcResponse{}, false
	}

	switch req.Method {
	case "initialize":
		return s.initialize(req), true
	case "ping":
		return ok(req.ID, map[string]any{}), true
	case "tools/list":
		return s.listTools(c, req), true
	case "tools/call":
		return s.callTool(c, req), true
	default:
		return fail(req.ID, codeMethodNotFound,
			fmt.Sprintf("this server does not implement %q. It has initialize, ping, "+
				"tools/list and tools/call.", req.Method)), true
	}
}

func (s *Server) initialize(req rpcRequest) rpcResponse {
	var params struct {
		ProtocolVersion string ` + "`" + `json:"protocolVersion"` + "`" + `
	}
	_ = json.Unmarshal(req.Params, &params)

	version := negotiate(params.ProtocolVersion)

	return ok(req.ID, map[string]any{
		"protocolVersion": version,
		"capabilities": map[string]any{
			// listChanged false: the catalogue is fixed for the life of the
			// process, so promising notifications we never send would make a
			// client wait for one.
			"tools": map[string]any{"listChanged": false},
		},
		"serverInfo": map[string]any{
			"name":    s.AppName,
			"version": s.AppVersion,
		},
		// Read by the model before it does anything, so it is worth spending
		// words on: what this app is, what today is, and the house rules a
		// model cannot infer from a tool schema.
		"instructions": s.Instructions(),
	})
}

// negotiate picks the newest version both sides know.
//
// An unknown version gets the oldest we speak rather than a refusal: a client
// from the future that sends a date we have never heard of can still talk to
// us over 2024-11-05, and a refusal here would look like the server being down.
func negotiate(asked string) string {
	for _, known := range protocolVersions {
		if asked == known {
			return known
		}
	}
	if asked == "" {
		return protocolVersions[0]
	}
	return protocolVersions[len(protocolVersions)-1]
}

func (s *Server) listTools(c *gin.Context, req rpcRequest) rpcResponse {
	conn, err := s.authenticate(c)
	if err != nil {
		return fail(req.ID, codeInvalidRequest, err.Error())
	}

	tools := Named(conn.ToolNames())
	out := make([]map[string]any, 0, len(tools))
	for _, tool := range tools {
		schema := tool.Schema
		if schema == nil {
			schema = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out = append(out, map[string]any{
			"name":        tool.Name,
			"title":       tool.Title,
			"description": tool.Description,
			"inputSchema": schema,
			"annotations": map[string]any{
				"readOnlyHint":    tool.ReadOnly(),
				"destructiveHint": tool.Destructive,
			},
		})
	}
	return ok(req.ID, map[string]any{"tools": out})
}

func (s *Server) callTool(c *gin.Context, req rpcRequest) rpcResponse {
	conn, err := s.authenticate(c)
	if err != nil {
		return fail(req.ID, codeInvalidRequest, err.Error())
	}

	var params struct {
		Name      string         ` + "`" + `json:"name"` + "`" + `
		Arguments map[string]any ` + "`" + `json:"arguments"` + "`" + `
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return fail(req.ID, codeInvalidParams, "params must be an object with name and arguments")
	}
	if params.Name == "" {
		return fail(req.ID, codeInvalidParams, "no tool name")
	}

	tool, found := Lookup(params.Name)
	if !found {
		return fail(req.ID, codeInvalidParams,
			fmt.Sprintf("there is no tool called %q", params.Name))
	}
	// The connection's own list is the authority, not the registry: a tool
	// exists in the process and is still not this connection's to call.
	if !conn.Allows(params.Name) {
		return fail(req.ID, codeInvalidParams,
			fmt.Sprintf("this connection was not given %q", params.Name))
	}

	call := &Call{
		Engine:         s.Engine,
		DB:             s.DB,
		Token:          conn.Token,
		OrganizationID: conn.OrganizationID,
		ConnectionID:   conn.ID,
		UserID:         conn.UserID,
	}

	started := time.Now()
	result, execErr := tool.Exec(c.Request.Context(), call, params.Arguments)
	elapsed := time.Since(started)

	s.record(c, conn, tool, params.Arguments, execErr, elapsed)

	if execErr != nil {
		// isError and not a protocol error, which is the distinction worth
		// getting right: a refusal is a result the model should read and act
		// on, and a protocol error is the session breaking. Answering 403 as
		// a protocol error makes a client give up instead of trying something
		// it is allowed to do.
		return ok(req.ID, map[string]any{
			"isError": true,
			"content": []map[string]any{{"type": "text", "text": execErr.Error()}},
		})
	}

	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fail(req.ID, codeInternalError, "the tool's answer could not be encoded")
	}

	answer := map[string]any{
		"content": []map[string]any{{"type": "text", "text": string(encoded)}},
	}
	// structuredContent as well as text, for a client that can use it. Only
	// for an object: the field is specified as one, and a list or a scalar
	// there is rejected by a strict client.
	if object, isObject := result.(map[string]any); isObject {
		answer["structuredContent"] = object
	}
	return ok(req.ID, answer)
}

// Reject answers the methods this transport does not take, with the header the
// spec asks for.
//
// A GET on the endpoint is how a client opens the server-sent-events half of
// streamable HTTP. This server does not implement it, and 405 with Allow: POST
// is the answer that tells a client to use the other half rather than leaving
// it waiting on a stream that will never open.
func Reject(c *gin.Context) {
	c.Header("Allow", "POST")
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": gin.H{
			"code": "METHOD_NOT_ALLOWED",
			"message": "this MCP endpoint takes POST. It answers every JSON-RPC " +
				"request with JSON rather than opening an event stream.",
		},
	})
}

// bearer reads the key from either header a client might use.
func bearer(c *gin.Context) string {
	if header := c.GetHeader("Authorization"); header != "" {
		if token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer")); token != header {
			return token
		}
		return strings.TrimSpace(header)
	}
	return strings.TrimSpace(c.GetHeader("X-API-Key"))
}
`
}
