// Package mcp implements a Model Context Protocol server over stdio, exposing
// what Grit knows about a project to an AI coding agent.
//
// Every tool here is READ-ONLY and STATIC: answers come from parsing the
// project's source, never from a running server or a database. That is a
// deliberate constraint rather than a first-draft limitation. It means the
// server works on a checkout that has never been started, needs no credentials,
// cannot mutate the repo, and cannot be talked into running a migration by a
// prompt-injected README. An agent that wants to change the project still has
// to call the CLI, where the change is visible in the diff.
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// jsonrpcVersion is the only JSON-RPC version MCP uses.
const jsonrpcVersion = "2.0"

// defaultProtocolVersion is the MCP revision this server is written against.
// A client that asks for a different revision gets its own echoed back when we
// recognise it, because these tools use only the stable initialize/tools
// surface that has not changed across these revisions.
const defaultProtocolVersion = "2025-06-18"

var knownProtocolVersions = map[string]bool{
	"2024-11-05": true,
	"2025-03-26": true,
	"2025-06-18": true,
}

// JSON-RPC error codes (see the JSON-RPC 2.0 spec).
const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
	codeInternalError  = -32603
)

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server answers MCP requests about one Grit project.
type Server struct {
	// Root is the project root — the directory holding grit.json.
	Root string
	// Version is the Grit CLI version, reported in serverInfo.
	Version string
}

// Serve runs the stdio loop until in is exhausted or closed.
//
// Nothing may be written to out except protocol messages: an MCP client parses
// stdout as a stream of JSON, so a stray banner or progress line desynchronises
// the session. Diagnostics belong on stderr, which the client ignores.
func (s *Server) Serve(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	// Route lists and model dumps exceed bufio's default 64KB line limit.
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)

	enc := json.NewEncoder(out)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			if encErr := enc.Encode(response{
				JSONRPC: jsonrpcVersion,
				Error:   &rpcError{Code: codeParseError, Message: "invalid JSON"},
			}); encErr != nil {
				return encErr
			}
			continue
		}

		resp, ok := s.handle(req)
		// A notification (no id) gets no reply at all — sending one is a
		// protocol violation that some clients treat as fatal.
		if !ok {
			continue
		}
		if err := enc.Encode(resp); err != nil {
			return fmt.Errorf("writing response: %w", err)
		}
	}

	return scanner.Err()
}

// handle dispatches one request. The second return is false when the message
// is a notification and must not be answered.
func (s *Server) handle(req request) (response, bool) {
	isNotification := len(req.ID) == 0

	reply := func(result interface{}) (response, bool) {
		return response{JSONRPC: jsonrpcVersion, ID: req.ID, Result: result}, !isNotification
	}
	fail := func(code int, msg string) (response, bool) {
		return response{
			JSONRPC: jsonrpcVersion,
			ID:      req.ID,
			Error:   &rpcError{Code: code, Message: msg},
		}, !isNotification
	}

	if req.JSONRPC != "" && req.JSONRPC != jsonrpcVersion {
		return fail(codeInvalidRequest, "unsupported jsonrpc version "+req.JSONRPC)
	}

	switch req.Method {
	case "initialize":
		return reply(s.initialize(req.Params))

	case "notifications/initialized", "notifications/cancelled":
		// Acknowledged by doing nothing; these carry no id.
		return response{}, false

	case "ping":
		return reply(map[string]interface{}{})

	case "tools/list":
		return reply(map[string]interface{}{"tools": toolDefinitions()})

	case "tools/call":
		var p struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &p); err != nil {
				return fail(codeInvalidParams, "invalid params: "+err.Error())
			}
		}
		if p.Name == "" {
			return fail(codeInvalidParams, "missing tool name")
		}

		text, err := s.callTool(p.Name, p.Arguments)
		if err != nil {
			// A tool that fails reports through the result with isError set,
			// not as a protocol error: the agent should see the message and be
			// able to try something else, rather than have the call look like a
			// transport fault.
			return reply(toolResult(err.Error(), true))
		}
		return reply(toolResult(text, false))

	default:
		return fail(codeMethodNotFound, "unknown method "+req.Method)
	}
}

func (s *Server) initialize(params json.RawMessage) map[string]interface{} {
	protocol := defaultProtocolVersion
	if len(params) > 0 {
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		if err := json.Unmarshal(params, &p); err == nil && knownProtocolVersions[p.ProtocolVersion] {
			protocol = p.ProtocolVersion
		}
	}

	return map[string]interface{}{
		"protocolVersion": protocol,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "grit",
			"version": s.Version,
		},
	}
}

// toolResult builds the content envelope MCP expects from tools/call.
func toolResult(text string, isError bool) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
		"isError": isError,
	}
}
