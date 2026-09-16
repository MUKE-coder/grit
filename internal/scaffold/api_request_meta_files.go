package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Contact-app review M29: services took *gin.Context.
//
// services.LogActivity, services.CreateSession and services.RotateSession all
// did, and all three wanted the same three facts from it: the client IP, the
// user agent and the request id. Taking the whole context meant none of them
// could be called from a background job, a seeder or a test without building a
// fake HTTP request, and it meant a service could reach for anything at all in
// the request, which is how business logic ends up in a handler's shape.
//
// RequestMeta is those three facts. Middleware records them once per request;
// a service reads them from the context.Context it already takes.

// writeRequestMetaFiles writes the request-meta plumbing.
//
// Its own writer, called by the scaffold and by grit upgrade, because both
// files are new and the handlers that call into them do not compile without
// them.
func writeRequestMetaFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "services", "request_meta.go"):   servicesRequestMetaGo(),
		filepath.Join(apiRoot, "internal", "middleware", "request_meta.go"): middlewareRequestMetaGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", opts.Module())
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// servicesRequestMetaGo returns internal/services/request_meta.go.
func servicesRequestMetaGo() string {
	return `package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// What a service is allowed to know about the request it is serving.
//
// A service that takes *gin.Context cannot be called from a job, a seeder or a
// test without a fake HTTP request, and it can reach for anything in the
// request, which is how the rules end up shaped like the transport. These three
// facts are all the services here ever wanted.

// RequestMeta is the request, as far as a service is concerned.
type RequestMeta struct {
	// IP is the client address, already resolved through the trusted proxy
	// rules by ResolveClientIP.
	IP string
	// UserAgent is the User-Agent header, unparsed.
	UserAgent string
	// RequestID is the X-Request-ID the client was given, so an audit row and a
	// log line can be lined up afterwards.
	RequestID string
	// UserID is the authenticated actor, empty on an anonymous request.
	UserID string
}

type requestMetaKey struct{}

// WithRequestMeta returns ctx carrying meta. Middleware calls this once per
// request; a job that wants an audit row attributed to something can call it
// too.
func WithRequestMeta(ctx context.Context, meta RequestMeta) context.Context {
	return context.WithValue(ctx, requestMetaKey{}, meta)
}

// MetaFrom returns what was recorded about the request, or the zero value
// outside one. A job has no IP and no user agent, and an audit row written from
// one says so rather than inventing them.
func MetaFrom(ctx context.Context) RequestMeta {
	if ctx == nil {
		return RequestMeta{}
	}
	meta, _ := ctx.Value(requestMetaKey{}).(RequestMeta)
	return meta
}

// MetaOf reads the meta off a gin request.
//
// The one function in this file that knows what gin is, and it exists so that
// nothing else has to: middleware.RequestMeta calls it, and everything
// downstream takes a context.Context.
func MetaOf(c *gin.Context) RequestMeta {
	if c == nil {
		return RequestMeta{}
	}
	meta := RequestMeta{
		IP:        ResolveClientIP(c),
		UserAgent: c.GetHeader("User-Agent"),
		RequestID: c.GetString("request_id"),
	}
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			meta.UserID = s
		}
	}
	return meta
}

// ContextOf is the request's context carrying its meta, for a handler calling a
// service that has not been handed the meta some other way.
func ContextOf(c *gin.Context) context.Context {
	if c == nil || c.Request == nil {
		return WithRequestMeta(context.Background(), MetaOf(c))
	}
	return WithRequestMeta(c.Request.Context(), MetaOf(c))
}

// LogActivityCtx writes a UserActivity row, taking the actor, the IP and the
// user agent from the context rather than from a *gin.Context.
//
// This is the one a job, a CLI command or a test calls. LogActivity is the same
// thing with a gin context in front of it.
//
// Errors are logged, not returned: losing an audit row should never fail a real
// request.
func LogActivityCtx(ctx context.Context, db *gorm.DB, args ActivityArgs) {
	row := activityRowCtx(ctx, args)
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		log.Printf("activity: failed to write %s: %v", args.Action, err)
	}
}

// activityRowCtx builds the row for args from whatever the context knows.
func activityRowCtx(ctx context.Context, args ActivityArgs) models.UserActivity {
	meta := MetaFrom(ctx)

	userID := args.UserID
	if userID == "" {
		userID = meta.UserID
	}

	var metaJSON string
	if args.Metadata != nil {
		if b, err := json.Marshal(args.Metadata); err == nil {
			metaJSON = string(b)
		}
	}

	return models.UserActivity{
		UserID:       userID,
		Action:       args.Action,
		Severity:     args.Severity,
		Summary:      args.Summary,
		ResourceType: args.ResourceType,
		ResourceID:   args.ResourceID,
		IPAddress:    meta.IP,
		UserAgent:    meta.UserAgent,
		Metadata:     metaJSON,
	}
}
`
}

// middlewareRequestMetaGo returns internal/middleware/request_meta.go.
func middlewareRequestMetaGo() string {
	return `package middleware

import (
	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/services"
)

// RequestMeta records the client IP, the user agent and the request id on the
// request's context, so a service can read them without taking a *gin.Context.
//
// It runs after RequestID, which is where the id comes from, and before Auth,
// which is where the actor comes from. The actor is therefore not in here: a
// service that wants it reads it from the context it was handed, which Auth
// refreshes once it knows who is calling.
func RequestMeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(services.ContextOf(c))
		c.Next()
	}
}
`
}
