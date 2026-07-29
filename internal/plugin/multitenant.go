package plugin

import "strings"

func init() { Register(multitenantPlugin()) }

// multitenantPlugin adds organizations, per-organization roles, and automatic
// query scoping.
//
// Deliberately a plugin rather than core: most apps are single-tenant, and
// putting org_id on every table is a schema decision no framework should make
// for you.
//
// Design notes:
//
//   - **No subdomains.** The active organization comes from a header
//     (X-Organization-ID) or the session. Subdomain routing forces DNS and TLS
//     decisions on every deployment and breaks local development.
//
//   - **Membership carries the role**, so a user is Editor in one organization
//     and Viewer in another. This reuses the existing roles table rather than
//     inventing a parallel permission system, and it's why user↔role was made
//     many-to-many in v3.66.0 — no breaking migration was needed here.
//
//   - **Scoping is automatic**, via a GORM callback. The reference
//     implementation this was modelled on hand-writes `Where("business_id = ?")`
//     447 times across 33 files; one missed call is a cross-tenant data leak.
//     Making it structural means a developer cannot forget it — they have to
//     opt OUT explicitly, which is visible in review.
func multitenantPlugin() Plugin {
	return Plugin{
		Name:    "multitenant",
		Version: "1.0.0",
		Summary: "Organizations, per-org roles, and automatic query scoping",
		Description: `Adds multi-tenancy where one user belongs to many organizations, holding a
different role in each.

  • Organization + OrganizationMember models
  • Active organization resolved from a header or the session — no subdomains
  • Automatic org scoping on every query via a GORM callback, so a forgotten
    WHERE cannot leak data across tenants
  • Per-organization permissions, reusing the roles you already have

Requires the roles system (Grit v3.66.0+).`,

		NextSteps: []string{
			"Run migrations:  cd apps/api && go run cmd/migrate/main.go",
			"Every tenant-scoped model needs an OrgID field — add `tenant.Owned` to it",
			"Clients send the active org as the X-Organization-ID header",
			"Read internal/tenant/tenant.go for how to opt a query OUT of scoping",
		},

		Files:      multitenantFiles,
		Injections: multitenantInjections,
	}
}

func multitenantFiles(ctx Context) map[string]string {
	api := "apps/api"
	if ctx.Architecture == "single" {
		api = "."
	}
	p := func(rel string) string {
		if api == "." {
			return rel
		}
		return api + "/" + rel
	}

	return map[string]string{
		p("internal/models/organization.go"):   mtOrganizationModel(ctx),
		p("internal/tenant/tenant.go"):         mtTenantPackage(ctx),
		p("internal/tenant/tenant_test.go"):    mtTenantTests(ctx),
		p("internal/middleware/tenant.go"):     mtMiddleware(ctx),
		p("internal/handlers/organization.go"): mtHandler(ctx),
	}
}

func multitenantInjections(ctx Context) []Injection {
	api := "apps/api"
	if ctx.Architecture == "single" {
		api = "."
	}
	p := func(rel string) string {
		if api == "." {
			return rel
		}
		return api + "/" + rel
	}

	return []Injection{
		{
			File:   p("internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&Organization{},\n\t\t&OrganizationMember{},",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\torgHandler := handlers.NewOrganizationHandler(db)",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: `		// Organizations. Listing and switching are available to any signed-in
		// user — they need to see the orgs they belong to before one is active.
		protected.GET("/organizations", orgHandler.MyOrganizations)
		protected.POST("/organizations", orgHandler.Create)
		protected.GET("/organizations/:id/members", orgHandler.Members)
		protected.POST("/organizations/:id/members", orgHandler.AddMember)
		protected.DELETE("/organizations/:id/members/:userId", orgHandler.RemoveMember)`,
		},
	}
}

// mtOrganizationModel emits the Organization + membership models.
func mtOrganizationModel(ctx Context) string {
	src := `package models

import (
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// Organization is a tenant.
//
// Slug exists for human-readable URLs and API calls; it is NOT used for
// subdomain routing, which this plugin deliberately avoids.
type Organization struct {
	ID          string ~gorm:"primarykey;size:36" json:"id"~
	Name        string ~gorm:"size:160;not null" json:"name" binding:"required"~
	Slug        string ~gorm:"size:160;uniqueIndex;not null" json:"slug"~
	Description string ~gorm:"size:500" json:"description"~

	// OwnerID is the user who created the organization. Kept separate from
	// membership roles so an owner can't be locked out by a permission change.
	OwnerID string ~gorm:"size:36;index" json:"owner_id"~

	// not null, but no default: a column default makes Active:false unstorable
	// on create, because GORM omits zero-valued fields from the INSERT. An
	// organization created suspended would come back live.
	Active    bool           ~gorm:"not null" json:"active"~
	Version   int            ~gorm:"not null;default:1" json:"version"~
	CreatedAt time.Time      ~json:"created_at"~
	UpdatedAt time.Time      ~json:"updated_at"~
	DeletedAt gorm.DeletedAt ~gorm:"index" json:"-"~
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = ids.New()
	}
	return nil
}

func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// OrganizationMember links a user to an organization WITH a role.
//
// The role lives on the membership, not the user, which is what lets someone be
// an Editor in one organization and a Viewer in another. RoleID points at the
// existing roles table — this plugin adds no parallel permission system.
type OrganizationMember struct {
	UserID         string ~gorm:"primaryKey;size:36;index" json:"user_id"~
	OrganizationID string ~gorm:"primaryKey;size:36;index" json:"organization_id"~
	RoleID         string ~gorm:"size:36;index" json:"role_id"~

	// No DeletedAt: soft-delete tombstones collide with the composite primary
	// key when someone is re-invited, which surfaces as a unique-constraint
	// error on an operation that should be idempotent.
	CreatedAt time.Time ~json:"created_at"~
	UpdatedAt time.Time ~json:"updated_at"~
}

func (OrganizationMember) TableName() string { return "organization_members" }
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// mtTenantPackage emits the scoping engine — the security-critical part.
func mtTenantPackage(ctx Context) string {
	src := `// Package tenant provides organization scoping.
//
// The important piece here is RegisterScoping: it installs GORM callbacks that
// add "WHERE org_id = ?" to every query against a tenant-owned model, using the
// active organization on the request context.
//
// Why automatic rather than per-query: hand-written scoping fails open. A
// developer who forgets one WHERE clause creates a cross-tenant data leak that
// no test catches, because the query returns MORE rows rather than erroring.
// Making it structural inverts that — you must opt out explicitly, and the
// opt-out is visible in code review.
package tenant

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ctxKey string

const orgKey ctxKey = "grit.org_id"

// ErrNoOrganization is returned when an operation needs an active organization
// and none is set.
var ErrNoOrganization = errors.New("tenant: no active organization")

// Owned marks a model as belonging to an organization.
//
// Embed it in any model that holds tenant data:
//
//	type Invoice struct {
//	    ID string
//	    tenant.Owned
//	    ...
//	}
//
// The scoping callbacks look for the OrgID column, so embedding this is what
// opts a table into isolation.
type Owned struct {
	OrgID string ~gorm:"size:36;index;not null" json:"org_id"~
}

// WithOrg returns a context carrying the active organization.
func WithOrg(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, orgKey, orgID)
}

// FromContext reads the active organization.
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(orgKey).(string)
	return id, ok && id != ""
}

// skipKey marks a *gorm.DB session as exempt from scoping.
const skipKey = "tenant:skip"

// Unscoped returns a session that bypasses organization scoping.
//
// For genuinely cross-tenant work — platform admin screens, background jobs
// that sweep every organization, migrations. Every call is a deliberate,
// greppable decision:
//
//	tenant.Unscoped(db).Find(&all)
func Unscoped(db *gorm.DB) *gorm.DB {
	return db.Set(skipKey, true)
}

// RegisterScoping installs the callbacks. Call once at startup, after opening
// the database:
//
//	tenant.RegisterScoping(db)
func RegisterScoping(db *gorm.DB) error {
	if err := db.Callback().Query().Before("gorm:query").Register("tenant:query", scope); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("tenant:update", scope); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register("tenant:delete", scope); err != nil {
		return err
	}
	// Create stamps the org id rather than filtering, so a new row can't be
	// written into the wrong tenant (or none).
	return db.Callback().Create().Before("gorm:create").Register("tenant:create", stamp)
}

// hasOrgColumn reports whether the statement's model is tenant-owned.
func hasOrgColumn(db *gorm.DB) bool {
	if db.Statement == nil || db.Statement.Schema == nil {
		return false
	}
	return db.Statement.Schema.LookUpField("OrgID") != nil
}

func skipped(db *gorm.DB) bool {
	if v, ok := db.Get(skipKey); ok {
		if b, ok := v.(bool); ok && b {
			return true
		}
	}
	return false
}

// scope adds the organization filter to reads and writes.
func scope(db *gorm.DB) {
	if db.Error != nil || skipped(db) || !hasOrgColumn(db) {
		return
	}
	orgID, ok := FromContext(db.Statement.Context)
	if !ok {
		// No active organization on a tenant-owned model. Fail CLOSED: erroring
		// is noisy, but returning every tenant's rows is a data breach.
		// Use tenant.Unscoped(db) where that's genuinely intended.
		db.AddError(ErrNoOrganization)
		return
	}
	db.Statement.AddClause(orgClause(db, orgID))
}

// stamp sets OrgID on insert.
func stamp(db *gorm.DB) {
	if db.Error != nil || skipped(db) || !hasOrgColumn(db) {
		return
	}
	orgID, ok := FromContext(db.Statement.Context)
	if !ok {
		db.AddError(ErrNoOrganization)
		return
	}
	field := db.Statement.Schema.LookUpField("OrgID")
	if field == nil || db.Statement.ReflectValue.Kind() == 0 {
		return
	}
	// Only fill an empty value, so an explicit assignment still wins.
	if v, zero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); zero || v == "" {
		_ = field.Set(db.Statement.Context, db.Statement.ReflectValue, orgID)
	}
}

// orgClause builds the WHERE fragment, qualified with the table name so it
// stays correct when the query joins.
func orgClause(db *gorm.DB, orgID string) clause.Where {
	table := db.Statement.Table
	if table == "" && db.Statement.Schema != nil {
		table = db.Statement.Schema.Table
	}
	return clause.Where{Exprs: []clause.Expression{
		clause.Eq{
			Column: clause.Column{Table: table, Name: "org_id"},
			Value:  orgID,
		},
	}}
}
`
	return strings.ReplaceAll(src, "~", "`")
}
