package plugin

import "path/filepath"

// mcpMountGo is internal/mcp/mount.go: the server, its auth, and its log.
func mcpMountGo(ctx Context) string {
	// The fallback the server reports when APP_NAME is unset. The project's own
	// directory, so an agent is told the name of the app it is working in
	// rather than Grit's.
	name := filepath.Base(ctx.Root)
	if name == "." || name == string(filepath.Separator) {
		name = "this app"
	}

	return `package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + ctx.Module + `/internal/ids"
	"` + ctx.Module + `/internal/models"
)

// Deps is what the server needs and must not import.
//
// Verifying an API key lives in internal/services, and the service that
// creates a connection needs this package to read its tools: importing each
// other is a cycle, and Go says so at build time. So the server is told how to
// check a key rather than reaching for the function itself. routes/mcp_wire.go
// supplies these, because routes is above both.
//
// It also makes the server testable without a database: a test supplies its
// own VerifyKey and never mints anything.
type Deps struct {
	// VerifyKey turns a token into the id of the key it belongs to, or an
	// error for one that is unknown, expired or revoked.
	VerifyKey func(token string) (keyID string, err error)
	// TouchKey records that a key was used. Best effort by nature: a failure
	// to update "last used" must not fail the call.
	TouchKey func(keyID string)
}

// Server answers /mcp.
type Server struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Deps   Deps

	AppName    string
	AppVersion string
}

// Mount puts the server on the router.
//
// At the root and not under /api, because it carries its own credential and
// speaks JSON-RPC rather than REST: the /api groups' middleware would
// authenticate it twice, with cookies it does not have.
//
// It is handed the engine so a tool can replay a request through it. Taken
// mid-Setup and used only at request time, by which point the engine is fully
// built: see call.go for why a tool goes through the router at all.
func Mount(r *gin.Engine, db *gorm.DB, deps Deps) *Server {
	s := &Server{
		Engine:     r,
		DB:         db,
		Deps:       deps,
		AppName:    appName(),
		AppVersion: appVersion(),
	}

	r.POST("/mcp", s.Handle)
	// Both answered rather than left to gin's 404, so a client that opens the
	// stream half of streamable HTTP is told to use POST instead of being told
	// the endpoint does not exist.
	r.GET("/mcp", Reject)
	r.DELETE("/mcp", Reject)

	log.Printf("MCP server on POST /mcp with %d tool(s)", len(All()))
	return s
}

func appName() string {
	if name := strings.TrimSpace(os.Getenv("APP_NAME")); name != "" {
		return name
	}
	return "` + name + `"
}

func appVersion() string {
	if v := strings.TrimSpace(os.Getenv("APP_VERSION")); v != "" {
		return v
	}
	return "0.0.0"
}

// Instructions is what the model reads before it does anything.
//
// Worth spending words on: it is the only place to say what this app is for,
// what today is, and the conventions a tool schema cannot express. A model
// without the date guesses at "this month", and a model that does not know
// drafts are for a person to send will send them.
func (s *Server) Instructions() string {
	var b strings.Builder
	fmt.Fprintf(&b, "You are working inside %s, through its own API.\n\n", s.AppName)
	fmt.Fprintf(&b, "Today is %s.\n\n", time.Now().Format("Monday, 2 January 2006"))
	b.WriteString("Every tool call runs as the person who created this connection, " +
		"with their permissions. A refusal is theirs to resolve, not something to " +
		"work around: do not try another tool to reach the same data.\n\n")
	b.WriteString("Ids are UUIDs. Where a tool takes a name instead, it will say so.\n")
	if extra := strings.TrimSpace(os.Getenv("MCP_INSTRUCTIONS")); extra != "" {
		b.WriteString("\n" + extra + "\n")
	}
	return b.String()
}

// connection is an authenticated caller: which connection, acting as whom,
// with which tools.
type connection struct {
	ID             string
	UserID         string
	OrganizationID string
	Token          string
	Tools          []string
}

func (c connection) ToolNames() []string { return c.Tools }

// Allows reports whether this connection was given a tool.
//
// Its own list and not the registry's: a tool can exist in the process and
// still not be this connection's to call, which is what makes a read-only
// connection read-only.
func (c connection) Allows(name string) bool {
	for _, allowed := range c.Tools {
		if allowed == name {
			return true
		}
	}
	return false
}

// authenticate turns the request's key into a connection.
//
// The key is verified by the app's own service, so a revoked or expired key is
// refused here for the same reason it would be refused on any other endpoint.
// The connection row is then looked up by the key it was minted with, which
// is what ties a call to a person and a tool list.
func (s *Server) authenticate(c *gin.Context) (connection, error) {
	token := bearer(c)
	if token == "" {
		s.challenge(c)
		return connection{}, fmt.Errorf("this endpoint needs a connection key. Make one " +
			"in the app under Settings, then send it as Authorization: Bearer <key>.")
	}

	if s.Deps.VerifyKey == nil {
		return connection{}, fmt.Errorf("mcp: this server was mounted without a way " +
			"to verify a key, so no connection can be authenticated")
	}
	keyID, err := s.Deps.VerifyKey(token)
	if err != nil || keyID == "" {
		s.challenge(c)
		return connection{}, fmt.Errorf("that key is not valid, or has been revoked. " +
			"Make a new connection in the app.")
	}

	var row models.MCPConnection
	if err := s.DB.Where("api_key_id = ? AND revoked_at IS NULL", keyID).
		First(&row).Error; err != nil {
		s.challenge(c)
		return connection{}, fmt.Errorf("that key is valid but is not an MCP connection. " +
			"Use the key a connection gave you, not an ordinary API key.")
	}

	if s.Deps.TouchKey != nil {
		s.Deps.TouchKey(keyID)
	}
	// last_used_at on every request, including tools/list, because "when did
	// this agent last connect" is the question it answers.
	//
	// call_count is NOT bumped here. It is bumped in record, on an actual tool
	// call, because that is what the name says and what somebody reading
	// "three calls" expects. Counting the handshake and every tools/list made
	// an idle client look busy.
	s.DB.Model(&row).UpdateColumn("last_used_at", time.Now())

	return connection{
		ID:             row.ID,
		UserID:         row.UserID,
		OrganizationID: row.OrganizationID,
		Token:          token,
		Tools:          []string(row.Tools),
	}, nil
}

// challenge sets WWW-Authenticate, so a client knows this was an auth failure
// rather than a bad request and can say so to the person using it.
func (s *Server) challenge(c *gin.Context) {
	c.Header("WWW-Authenticate", ` + "`" + `Bearer realm="mcp"` + "`" + `)
	c.Status(http.StatusOK) // the JSON-RPC error carries the detail
}

// record writes one call to the log.
//
// Every call, whether it worked or not, with the arguments. An owner looking
// at a connection needs to see what an agent did, and "it called create_invoice
// nine times" is not an answer without the arguments.
//
// Best effort: a failure to log must not fail the call, but it is logged
// locally so a silent gap in the trail is visible to whoever runs the server.
func (s *Server) record(c *gin.Context, conn connection, tool Tool, args map[string]any, execErr error, took time.Duration) {
	if s.DB == nil {
		return
	}

	encoded, err := json.Marshal(args)
	if err != nil {
		encoded = []byte("{}")
	}

	entry := models.MCPToolCall{
		ID:             ids.New(),
		ConnectionID:   conn.ID,
		UserID:         conn.UserID,
		OrganizationID: conn.OrganizationID,
		Tool:           tool.Name,
		Arguments:      string(encoded),
		OK:             execErr == nil,
		DurationMS:     took.Milliseconds(),
	}
	if execErr != nil {
		entry.Error = execErr.Error()
	}

	if err := s.DB.Create(&entry).Error; err != nil {
		log.Printf("mcp: could not record the call to %s: %v", tool.Name, err)
	}

	// The count a person reads, so it counts tool calls and nothing else.
	s.DB.Model(&models.MCPConnection{}).
		Where("id = ?", conn.ID).
		UpdateColumn("call_count", gorm.Expr("call_count + 1"))
}
`
}
