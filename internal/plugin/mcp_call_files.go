package plugin

// mcpCallGo is internal/mcp/call.go: how a tool reaches the app.
func mcpCallGo(ctx Context) string {
	return `package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Call is what a tool is given: the identity it acts as, and a way to reach
// the app.
//
// # Why a tool does not touch the database
//
// Exec builds an HTTP request and runs it in process against the same gin
// engine that serves the browser. Not a shortcut, and not for want of a
// cheaper path: it is the design.
//
// A Grit API decides what a caller may do in middleware and in services:
// authentication, the organization scope, the role, the per-resource policy,
// the module switch, the record-level scope, the validation tags, the audit
// log. A tool that queried GORM directly would bypass every one of them and
// would have to reimplement them, which means two implementations that start
// identical and drift, and the one that drifts is the one nobody is looking
// at. The first hand-built version of this had exactly that shape and it is
// the thing worth not repeating (grit#224).
//
// Running through the router costs a function call and some JSON. What it buys
// is that an agent is structurally incapable of doing something the same
// person could not do in the browser, and that every tool call lands in the
// audit log the same way a click does.
type Call struct {
	// Engine is the app's own router, handed over by Mount.
	Engine *gin.Engine
	// DB is for the bookkeeping around a call, never for answering one.
	DB *gorm.DB

	// Token is the connection's API key. Sent on every replayed request, so
	// the app authenticates the agent exactly as it would any API client, and
	// the key's endpoint allowlist is enforced by the app rather than here.
	Token string
	// OrganizationID scopes a multitenant project. Empty when the plugin is
	// not installed, which is the single-tenant case.
	OrganizationID string

	// ConnectionID and UserID are for the call log.
	ConnectionID string
	UserID       string
}

// Result is one replayed request's answer.
type Result struct {
	Status int
	Body   []byte
}

// OK reports a 2xx.
func (r Result) OK() bool { return r.Status >= 200 && r.Status < 300 }

// Decode unmarshals the body into v. A non-2xx is not decoded: read Status and
// Message instead, because an error envelope has a different shape from the
// payload a caller is expecting and silently decoding one into the other is
// how a tool reports success on a refusal.
func (r Result) Decode(v any) error {
	if !r.OK() {
		return fmt.Errorf("%s", r.Message())
	}
	if len(r.Body) == 0 {
		return nil
	}
	var envelope struct {
		Data json.RawMessage ` + "`" + `json:"data"` + "`" + `
	}
	if err := json.Unmarshal(r.Body, &envelope); err == nil && len(envelope.Data) > 0 {
		return json.Unmarshal(envelope.Data, v)
	}
	// Not the standard envelope: an endpoint that answers a bare object.
	return json.Unmarshal(r.Body, v)
}

// Message is the app's own sentence for a refusal, with its field errors.
//
// The app's words and not ours, which is the point: a model that is told
// "Check the highlighted fields: email must be a valid address" can correct
// itself and try again. A protocol-level error, or a generic "the call
// failed", gives it nothing to act on.
func (r Result) Message() string {
	var envelope struct {
		Error struct {
			Code    string            ` + "`" + `json:"code"` + "`" + `
			Message string            ` + "`" + `json:"message"` + "`" + `
			Details map[string]string ` + "`" + `json:"details"` + "`" + `
		} ` + "`" + `json:"error"` + "`" + `
	}
	if err := json.Unmarshal(r.Body, &envelope); err != nil || envelope.Error.Message == "" {
		return fmt.Sprintf("the request failed with status %d", r.Status)
	}

	msg := envelope.Error.Message
	if len(envelope.Error.Details) > 0 {
		parts := make([]string, 0, len(envelope.Error.Details))
		for field, why := range envelope.Error.Details {
			parts = append(parts, field+": "+why)
		}
		sort.Strings(parts)
		msg += " (" + strings.Join(parts, "; ") + ")"
	}
	if envelope.Error.Code != "" {
		msg += " [" + envelope.Error.Code + "]"
	}
	return msg
}

// Do replays one request through the app.
//
// query may be nil. body may be nil for a GET or DELETE.
func (c *Call) Do(ctx context.Context, method, path string, query url.Values, body any) (Result, error) {
	if c == nil || c.Engine == nil {
		return Result{}, fmt.Errorf("mcp: no engine to call, so this tool cannot reach the app")
	}

	target := path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}

	var reader *bytes.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return Result{}, fmt.Errorf("mcp: encoding the request body: %w", err)
		}
		reader = bytes.NewReader(encoded)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, target, reader).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// The connection's own credential. Not the caller's cookie: there is no
	// browser here, and a key is what can be narrowed to this connection's
	// routes and revoked on its own.
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("X-API-Key", c.Token)
	}
	if c.OrganizationID != "" {
		req.Header.Set("X-Organization-ID", c.OrganizationID)
	}
	// So a handler, a log line or the audit trail can tell an agent's request
	// from a person's.
	req.Header.Set("X-Grit-Client", "mcp")

	rec := httptest.NewRecorder()
	c.Engine.ServeHTTP(rec, req)

	return Result{Status: rec.Code, Body: rec.Body.Bytes()}, nil
}

// Get is Do for a read.
func (c *Call) Get(ctx context.Context, path string, query url.Values) (Result, error) {
	return c.Do(ctx, http.MethodGet, path, query, nil)
}

// Post is Do for a create.
func (c *Call) Post(ctx context.Context, path string, body any) (Result, error) {
	return c.Do(ctx, http.MethodPost, path, nil, body)
}

// Patch is Do for an update.
func (c *Call) Patch(ctx context.Context, path string, body any) (Result, error) {
	return c.Do(ctx, http.MethodPatch, path, nil, body)
}

// Delete is Do for a removal.
func (c *Call) Delete(ctx context.Context, path string) (Result, error) {
	return c.Do(ctx, http.MethodDelete, path, nil, nil)
}
`
}
