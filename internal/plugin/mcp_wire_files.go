package plugin

// mcpWireGo is internal/routes/mcp_wire.go: the adapter that breaks the cycle.
func mcpWireGo(ctx Context) string {
	return `package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + ctx.Module + `/internal/mcp"
	"` + ctx.Module + `/internal/services"
)

// mountMCP puts the MCP server on the router, with the one thing it cannot
// reach for itself.
//
// internal/mcp must not import internal/services, because the service that
// creates a connection imports internal/mcp to read a tool's permission and
// routes. Go rejects that pair outright, which is the right outcome: the
// server's job is to speak the protocol, not to know how this app stores
// credentials.
//
// So the server is told how to check a key. This file is where the two meet,
// in package routes, which is above both. It also means a test can mount a
// server with its own verifier and never touch a database.
func mountMCP(r *gin.Engine, db *gorm.DB) *mcp.Server {
	return mcp.Mount(r, db, mcp.Deps{
		VerifyKey: func(token string) (string, error) {
			key, err := services.VerifyAPIKey(db, token)
			if err != nil {
				return "", err
			}
			return key.ID, nil
		},
		TouchKey: func(keyID string) {
			services.TouchAPIKey(db, keyID)
		},
	})
}
`
}
