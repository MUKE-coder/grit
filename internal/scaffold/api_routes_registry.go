package scaffold

// apiRoutesRegistryGo emits internal/routes/resources.go: the small piece of
// machinery that lets every generated resource own its routes file.
//
// Before this, `grit generate resource` wrote into routes.go in four places:
// the handler construction, the protected block, the admin block, and
// sometimes the public one. routes.go passed a thousand lines in a project
// with a handful of resources, and the four edits for one resource were
// hundreds of lines apart. Finding where to add a route meant reading the
// whole file, and the generator had four chances to inject into the wrong
// block.
//
// Now each resource gets internal/routes/<resource>_routes.go holding
// everything it mounts, and routes.go never changes again.
func apiRoutesRegistryGo() string {
	return `package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/config"
)

// Mount is what a resource's route file receives: the router groups, already
// built and already carrying their middleware, plus the things a handler needs
// to construct itself.
//
// The groups are the same variables routes.go uses, so a route registered here
// is indistinguishable from one written inline. Public is outside the auth
// middleware and behind an API key; Protected takes a JWT or an API key; Admin
// additionally requires the ADMIN role.
type Mount struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Cfg    *config.Config
	Svc    *Services

	// V1 is the /api/v1 group, for a route that fits none of the three below.
	V1        *gin.RouterGroup
	Public    *gin.RouterGroup
	Protected *gin.RouterGroup
	Admin     *gin.RouterGroup
}

// mounts is every resource route file that has registered itself.
var mounts []func(*Mount)

// RegisterRoutes adds a resource's routes to the application.
//
// Generated resource files call this from an init(), which is why creating
// <resource>_routes.go is all it takes to mount a resource and deleting the
// file is all it takes to unmount one. Nothing else in the package refers to
// it by name.
//
// Go runs init() functions in a file-name order fixed by the compiler, so the
// registration order is deterministic. It does not matter in practice: Gin
// routes on a tree, not a list, and prefers a static segment over a parameter
// however the two were added.
func RegisterRoutes(fn func(*Mount)) {
	mounts = append(mounts, fn)
}

// mountResources runs every registered resource file. Called once from Setup,
// after the groups exist and before the legacy alias fallback.
func mountResources(m *Mount) {
	for _, mount := range mounts {
		mount(m)
	}
}
`
}
