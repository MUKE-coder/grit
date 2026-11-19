package plugin

// mcpConnectionHandlerGo is internal/handlers/mcp_connection.go.
func mcpConnectionHandlerGo(ctx Context) string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + ctx.Module + `/internal/authz"
	"` + ctx.Module + `/internal/mcp"
	"` + ctx.Module + `/internal/models"
	"` + ctx.Module + `/internal/respond"
	"` + ctx.Module + `/internal/services"
)

// MCPConnectionHandler is the connections a person manages for themselves.
//
// Mounted on the protected group, not the admin one: a connection acts as the
// person who made it and can do no more than they can, so it is theirs to
// create and revoke without an administrator in the loop. The same reason an
// API key is.
type MCPConnectionHandler struct {
	DB *gorm.DB
}

func NewMCPConnectionHandler(db *gorm.DB) *MCPConnectionHandler {
	return &MCPConnectionHandler{DB: db}
}

// List is GET /api/v1/mcp-connections: this person's own connections.
//
// Theirs only. A connection is a credential acting as somebody, and the list
// of whose agents exist is not something one user should read about another.
func (h *MCPConnectionHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	var rows []models.MCPConnection
	if err := h.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		respond.Internal(c, err)
		return
	}

	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		// Named rather than the stored list, so a connection whose resource
		// was removed shows the tools it actually has, and missing says why
		// the count is smaller than it was.
		out = append(out, gin.H{
			"id":             row.ID,
			"name":           row.Name,
			"tools":          []string(row.Tools),
			"missing":        mcp.ToolsMissing(row.Tools),
			"last_used_at":   row.LastUsedAt,
			"call_count":     row.CallCount,
			"revoked_at":     row.RevokedAt,
			"created_at":     row.CreatedAt,
			"read_only":      readOnly(row.Tools),
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// readOnly reports a connection that holds no tool that writes, which is what
// a person actually wants to know at a glance.
func readOnly(names []string) bool {
	for _, tool := range mcp.Named(names) {
		if tool.Write {
			return false
		}
	}
	return true
}

// Tools is GET /api/v1/mcp-connections/tools: what this person may choose.
//
// Filtered by their own grants, so the screen cannot offer a tool their role
// does not allow. The service checks it again: this is the convenience, not
// the boundary.
func (h *MCPConnectionHandler) Tools(c *gin.Context) {
	grants := authz.CallerGrants(c)

	out := make([]gin.H, 0)
	for _, tool := range mcp.All() {
		if tool.Permission != "" && !authz.Granted(grants, tool.Permission) {
			continue
		}
		out = append(out, gin.H{
			"name":        tool.Name,
			"title":       tool.Title,
			"description": tool.Description,
			"group":       tool.Group,
			"write":       tool.Write,
			"destructive": tool.Destructive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

type createConnectionRequest struct {
	Name  string   ` + "`" + `json:"name" binding:"required,min=2,max=120"` + "`" + `
	Tools []string ` + "`" + `json:"tools" binding:"required,min=1"` + "`" + `
}

// Create is POST /api/v1/mcp-connections.
//
// Answers with the key, once. It is hashed on the way in, so there is nothing
// to show later and the response says so: a person who loses it makes a new
// connection rather than recovering this one.
func (h *MCPConnectionHandler) Create(c *gin.Context) {
	var req createConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Invalid(c, err)
		return
	}

	created, err := services.CreateMCPConnection(h.DB, services.NewConnection{
		Name:           req.Name,
		UserID:         c.GetString("user_id"),
		OrganizationID: c.GetString("organization_id"),
		Tools:          req.Tools,
		Grants:         authz.CallerGrants(c),
	})
	if err != nil {
		// A refusal here is the caller's to fix: an unknown tool, or one their
		// role does not grant. Both are 422 with the sentence, not a 500.
		respond.Validation(c, err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"connection": created.Connection,
			"token":      created.Token,
		},
		"message": "Connection created. Copy the key now: it is not stored and " +
			"cannot be shown again.",
	})
}

// Revoke is DELETE /api/v1/mcp-connections/:id.
func (h *MCPConnectionHandler) Revoke(c *gin.Context) {
	userID := c.GetString("user_id")

	var conn models.MCPConnection
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).
		First(&conn).Error; err != nil {
		// Somebody else's connection is not found rather than forbidden, so a
		// wrong guess cannot be told from a right one.
		respond.NotFound(c, "No such connection")
		return
	}

	if err := services.RevokeMCPConnection(h.DB, conn.ID, userID); err != nil {
		respond.Internal(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Connection revoked"})
}

// Calls is GET /api/v1/mcp-connections/:id/calls: what this agent did.
func (h *MCPConnectionHandler) Calls(c *gin.Context) {
	userID := c.GetString("user_id")

	var conn models.MCPConnection
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).
		First(&conn).Error; err != nil {
		respond.NotFound(c, "No such connection")
		return
	}

	var calls []models.MCPToolCall
	if err := h.DB.Where("connection_id = ?", conn.ID).
		Order("created_at DESC").
		Limit(100).
		Find(&calls).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": calls})
}
`
}

// mcpBuiltinToolsGo is internal/mcp/tools_builtin.go: enough that tools/list
// is never empty.
func mcpBuiltinToolsGo(ctx Context) string {
	return `package mcp

import (
	"context"
	"fmt"
)

// Two tools every project has, registered here so a fresh install answers
// tools/list with something and a client can be tested before any resource
// has tools of its own.
//
// whoami is also the honest first thing for an agent to call: it tells the
// model which account it is acting as, which is the fact every later refusal
// will turn on.
func init() {
	Register(
		Tool{
			Name:        "whoami",
			Title:       "Who am I",
			Description: "The account this connection acts as, and its role. Call this " +
				"first if a request is refused: the answer is usually that this " +
				"account does not have the permission.",
			Group:  "Account",
			Routes: []string{"GET /api/v1/auth/me"},
			Schema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
			Exec: func(ctx context.Context, call *Call, args map[string]any) (any, error) {
				result, err := call.Get(ctx, "/api/v1/auth/me", nil)
				if err != nil {
					return nil, err
				}
				if !result.OK() {
					return nil, fmt.Errorf("%s", result.Message())
				}
				var me map[string]any
				if err := result.Decode(&me); err != nil {
					return nil, err
				}
				return me, nil
			},
		},
		Tool{
			Name:        "list_tools",
			Title:       "List the tools I have",
			Description: "Every tool this connection may call, with what each one does. " +
				"Useful when a task does not obviously map to one of them.",
			Group:  "Account",
			Schema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
			Exec: func(ctx context.Context, call *Call, args map[string]any) (any, error) {
				tools := All()
				out := make([]map[string]any, 0, len(tools))
				for _, tool := range tools {
					out = append(out, map[string]any{
						"name":        tool.Name,
						"description": tool.Description,
						"group":       tool.Group,
						"writes":      tool.Write,
					})
				}
				return map[string]any{"tools": out}, nil
			},
		},
	)
}
`
}

// mcpCustomToolsGo is internal/mcp/custom_tools.go: the file Grit never
// rewrites.
func mcpCustomToolsGo(ctx Context) string {
	return `package mcp

// Your own tools.
//
// Grit never rewrites this file, and the generator never touches it, so an
// upgrade cannot take your work. The generated CRUD tools are the base, not
// the ceiling: the interesting tools in any app are the ones that are not
// create-read-update-delete.
//
// A worked example, commented out because it names a resource you may not
// have. Four things in it are worth copying:
//
//  1. Routes lists every endpoint Exec calls. This is a security boundary, not
//     documentation: a connection's key is restricted to the union of its
//     tools' routes, so a route left out here is refused at runtime.
//
//  2. Permission is the key a role must grant. The connection screen will not
//     offer this tool to somebody whose role lacks it, and the service checks
//     again when the connection is made.
//
//  3. Write and Destructive drive the hints a client uses to decide whether to
//     confirm with the person first. Getting these wrong is how an agent
//     deletes something nobody was asked about.
//
//  4. Exec goes through call.Post rather than touching GORM, so tenancy,
//     policies, validation and the audit log all apply. A tool that reaches
//     for the database is a second implementation of every rule in the app.
//
//	func init() {
//		Register(Tool{
//			Name:        "send_invoice",
//			Title:       "Send an invoice",
//			Description: "Emails an invoice to its client and marks it sent. " +
//				"The invoice must be a draft. Takes the invoice number or its id.",
//			Group:       "Invoices",
//			Permission:  "invoices.edit",
//			Write:       true,
//			Routes: []string{
//				"GET /api/v1/invoices",
//				"POST /api/v1/invoices/:id/send",
//			},
//			Schema: map[string]any{
//				"type": "object",
//				"properties": map[string]any{
//					"invoice": map[string]any{
//						"type":        "string",
//						"description": "The invoice number, or its id.",
//					},
//				},
//				"required": []string{"invoice"},
//			},
//			Exec: func(ctx context.Context, call *Call, args map[string]any) (any, error) {
//				term, err := String(args, "invoice")
//				if err != nil {
//					return nil, err
//				}
//				id, err := Resolve(ctx, call, "invoice", "/api/v1/invoices", term)
//				if err != nil {
//					return nil, err
//				}
//				result, err := call.Post(ctx, "/api/v1/invoices/"+id+"/send", nil)
//				if err != nil {
//					return nil, err
//				}
//				if !result.OK() {
//					return nil, fmt.Errorf("%s", result.Message())
//				}
//				return map[string]any{"sent": true, "invoice_id": id}, nil
//			},
//		})
//	}
`
}
