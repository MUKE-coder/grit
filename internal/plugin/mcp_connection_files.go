package plugin

// mcpModelsGo is internal/models/mcp.go.
func mcpModelsGo(ctx Context) string {
	return `package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MCPConnection is one agent's access to this app.
//
// A connection, not a key, is the unit a person thinks in: "Claude Code on my
// laptop", with a list of what it may do. The key is an implementation detail
// of it, minted when the connection is made and revoked with it.
//
// One per client, rather than one shared secret for the app. That is the whole
// difference: a shared secret cannot say who did what, cannot be narrowed to
// the twelve tools one person needs, and cannot be revoked without breaking
// everybody.
type MCPConnection struct {
	ID string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `

	// Name is for humans, shown in the list so a connection can be revoked
	// without guessing what it is.
	Name string ` + "`" + `gorm:"size:120;not null" json:"name"` + "`" + `

	// UserID is who it acts as. Every tool call runs with this person's
	// permissions, which is why the screen cannot offer a tool their role does
	// not grant: an agent must not be able to do more than the person who set
	// it up.
	UserID string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	User   User   ` + "`" + `gorm:"foreignKey:UserID" json:"-"` + "`" + `

	// OrganizationID scopes it in a multitenant project. Empty single-tenant.
	OrganizationID string ` + "`" + `gorm:"size:36;index" json:"organization_id,omitempty"` + "`" + `

	// APIKeyID is the key this connection authenticates with. The key holds
	// the endpoint allowlist, so the narrowing is enforced by the same
	// middleware that narrows any other API key.
	APIKeyID string ` + "`" + `gorm:"size:36;index;not null" json:"api_key_id"` + "`" + `

	// Tools is what it may call, by name. The authority at call time: a tool
	// can exist in the process and still not be this connection's, which is
	// what makes a read-only connection read-only.
	Tools datatypes.JSONSlice[string] ` + "`" + `json:"tools"` + "`" + `

	// LastUsedAt and CallCount are for the screen: a connection nobody has
	// used in three months is one to revoke.
	LastUsedAt *time.Time ` + "`" + `json:"last_used_at,omitempty"` + "`" + `
	CallCount  int64      ` + "`" + `gorm:"not null;default:0" json:"call_count"` + "`" + `

	// RevokedAt ends it without deleting the record, so the call log keeps
	// pointing at something.
	RevokedAt *time.Time ` + "`" + `json:"revoked_at,omitempty"` + "`" + `

	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// MCPToolCall is one thing an agent did.
//
// Every call, whether it worked or not, with the arguments it was given. An
// owner asking "what has this agent been doing" cannot be answered by a count:
// "it called create_invoice nine times" is not an answer without the
// arguments, and a refusal is as interesting as a success.
type MCPToolCall struct {
	ID string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `

	ConnectionID string ` + "`" + `gorm:"size:36;index;not null" json:"connection_id"` + "`" + `
	UserID       string ` + "`" + `gorm:"size:36;index" json:"user_id"` + "`" + `
	// OrganizationID is denormalised onto the row so an operator can read one
	// organization's agent activity without joining through the connection.
	OrganizationID string ` + "`" + `gorm:"size:36;index" json:"organization_id,omitempty"` + "`" + `

	Tool string ` + "`" + `gorm:"size:120;index;not null" json:"tool"` + "`" + `
	// Arguments as JSON text rather than a typed column: the shape differs per
	// tool, and this is read by a person rather than queried.
	Arguments string ` + "`" + `gorm:"type:text" json:"arguments"` + "`" + `

	OK         bool   ` + "`" + `gorm:"not null" json:"ok"` + "`" + `
	Error      string ` + "`" + `gorm:"type:text" json:"error,omitempty"` + "`" + `
	DurationMS int64  ` + "`" + `gorm:"not null;default:0" json:"duration_ms"` + "`" + `

	CreatedAt time.Time ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
}
`
}

// mcpConnectionServiceGo is internal/services/mcp_connection.go: making one.
func mcpConnectionServiceGo(ctx Context) string {
	return `package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"` + ctx.Module + `/internal/authz"
	"` + ctx.Module + `/internal/ids"
	"` + ctx.Module + `/internal/mcp"
	"` + ctx.Module + `/internal/models"
)

// NewConnection is what a connection costs: a name, a person, and the tools.
type NewConnection struct {
	Name           string
	UserID         string
	OrganizationID string
	Tools          []string
	// Grants is the permissions the person actually holds. A tool needing one
	// they lack is refused rather than silently dropped, because a connection
	// quietly missing the tool somebody chose is worse than being told why.
	Grants []string
	// ExpiresAt ends the key on its own. Optional, and worth setting for an
	// agent somebody is trying out.
	ExpiresAt *time.Time
}

// CreatedConnection carries the plaintext key, which exists for exactly as
// long as this value does.
type CreatedConnection struct {
	Connection *models.MCPConnection
	// Token is shown once and never stored. Same rule as any secret key: it is
	// hashed on the way in, so there is nothing to show later.
	Token string
}

// CreateMCPConnection mints a connection.
//
// The key it creates is narrowed twice over, and both narrowings matter:
//
//   - Scopes, to the permissions its own tools declare, so a key cannot be
//     used to reach a permission no chosen tool needs.
//   - Endpoints, to the routes those tools call. This is the one that makes a
//     read-only connection structurally read-only: it holds no POST route, so
//     a write is refused by the API key middleware before any handler runs,
//     whatever the agent asks for.
//
// Both are derived from the tools rather than chosen, so nobody can widen a
// connection by accident, and adding a tool to a resource later does not
// silently widen the connections that already exist.
func CreateMCPConnection(db *gorm.DB, req NewConnection) (*CreatedConnection, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("a connection needs a name, so it can be told apart from the others")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("a connection needs a user to act as")
	}
	if len(req.Tools) == 0 {
		return nil, fmt.Errorf("a connection with no tools can do nothing: choose at least one")
	}

	if missing := mcp.ToolsMissing(req.Tools); len(missing) > 0 {
		return nil, fmt.Errorf("no such tool: %v", missing)
	}

	tools := mcp.Named(req.Tools)

	// An agent must not be able to do more than the person who set it up.
	// Checked here and not only in the UI, because this is an API endpoint and
	// the UI is a convenience over it.
	for _, tool := range tools {
		if tool.Permission == "" {
			continue
		}
		if !authz.Granted(req.Grants, tool.Permission) {
			return nil, fmt.Errorf("%s needs the %s permission, which you do not have",
				tool.Name, tool.Permission)
		}
	}

	key, err := GenerateAPIKey(db, KeyOptions{
		UserID:    req.UserID,
		Name:      "MCP: " + req.Name,
		Kind:      models.KindSecret,
		Scopes:    mcp.PermissionsFor(tools),
		Endpoints: mcp.RoutesFor(tools),
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("minting the connection's key: %w", err)
	}

	conn := models.MCPConnection{
		ID:             ids.New(),
		Name:           req.Name,
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		APIKeyID:       key.Record.ID,
		Tools:          req.Tools,
	}
	if err := db.Create(&conn).Error; err != nil {
		// The key would otherwise be a credential with no connection behind
		// it, which authenticate() would refuse but which would still sit in
		// the keys list looking legitimate.
		_ = db.Delete(&models.APIKey{}, "id = ?", key.Record.ID).Error
		return nil, fmt.Errorf("saving the connection: %w", err)
	}

	return &CreatedConnection{Connection: &conn, Token: key.Token}, nil
}

// RevokeMCPConnection ends a connection and the key with it.
//
// Both, in that order, and the key is the one that matters: a connection
// marked revoked whose key still verifies is a connection that still works.
func RevokeMCPConnection(db *gorm.DB, id, userID string) error {
	var conn models.MCPConnection
	if err := db.Where("id = ?", id).First(&conn).Error; err != nil {
		return fmt.Errorf("no connection with that id")
	}

	if err := RevokeAPIKey(db, conn.APIKeyID, conn.UserID); err != nil {
		return fmt.Errorf("revoking the connection's key: %w", err)
	}

	now := time.Now()
	return db.Model(&conn).UpdateColumn("revoked_at", now).Error
}
`
}
