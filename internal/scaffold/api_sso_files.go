package scaffold

import "strings"

// Enterprise SSO (OpenID Connect).
//
// Consumer social login (Google/GitHub buttons) and enterprise SSO look similar
// but are different products: social login is one app-wide provider the app
// owner configures once, while SSO is one connection *per customer*, configured
// at runtime, routed by the user's email domain, with users provisioned on first
// login and their roles derived from the IdP's groups.
//
// This file emits that: an SSOConnection model, a UserIdentity table (the
// scaffold previously had only hardcoded google_id / github_id columns, which
// can't express "this user came from Acme's Okta"), a provider registry, and the
// begin/callback handlers.
//
// Why OIDC and not SAML: every current IdP — Okta, Entra ID, Auth0, Keycloak,
// Google Workspace, Ping, OneLogin — speaks OIDC, and goth (already a
// dependency) ships a generic openidConnect provider that needs only a
// discovery URL. SAML is still common in older enterprise procurement and is
// tracked separately; it needs a real SAML library and its own certificate
// handling.

// apiSSOModelGo emits internal/models/sso.go.
func apiSSOModelGo() string {
	src := `package models

import (
	"strings"
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"

	"{{MODULE}}/internal/crypto"
)

// SSOConnection is one customer's identity provider.
//
// A connection is identified by its slug in URLs (/api/auth/sso/acme), and
// matched to a user by the domain of the email they type on the login page.
// Everything needed to talk to a compliant OIDC provider is here: the discovery
// document URL plus a client ID and secret issued by that provider.
type SSOConnection struct {
	ID   string ~gorm:"primarykey;size:36" json:"id"~
	Slug string ~gorm:"size:64;uniqueIndex;not null" json:"slug" binding:"required"~
	Name string ~gorm:"size:255;not null" json:"name" binding:"required"~

	// Protocol is "oidc" (default) or "saml". OIDC needs an issuer plus client
	// credentials; SAML needs the IdP's metadata and no secret at all, because
	// trust rides on the IdP's signing certificate inside that metadata.
	Protocol string ~gorm:"size:10;default:'oidc'" json:"protocol"~

	// Domains is a comma-separated list of email domains routed to this
	// connection ("acme.com,acme.co.uk"). A user typing bob@acme.com on the
	// login page is sent here. Matching is case-insensitive and exact on the
	// domain part — no wildcards, because "*.com" is a security incident.
	Domains string ~gorm:"size:1000" json:"domains"~

	// IssuerURL is the provider's base issuer (e.g.
	// https://login.microsoftonline.com/<tenant>/v2.0 or https://acme.okta.com).
	// The discovery document is fetched from <issuer>/.well-known/openid-configuration
	// unless DiscoveryURL overrides it.
	IssuerURL    string ~gorm:"size:500;not null" json:"issuer_url" binding:"required"~
	DiscoveryURL string ~gorm:"size:500" json:"discovery_url"~

	ClientID string ~gorm:"size:255;not null" json:"client_id" binding:"required"~
	// ClientSecret is encrypted at rest via the same AES-256-GCM field
	// encryption used for other PII, and never serialized — the API returns
	// HasSecret instead so the admin UI can show "configured" without ever
	// shipping the value back to a browser.
	ClientSecret crypto.EncryptedString ~gorm:"type:text" json:"-"~
	HasSecret    bool                   ~gorm:"-" json:"has_secret"~

	// Scopes requested from the IdP. "openid" is always sent; profile and email
	// are the defaults because the callback needs a name and an address.
	Scopes string ~gorm:"size:500;default:'profile,email'" json:"scopes"~

	// NOTE: no ~default:true~ on these two booleans, deliberately.
	//
	// GORM omits zero-valued fields from an INSERT when the column carries a
	// default, so a ~default:true~ bool can never be stored as false through a
	// struct create — it silently comes back true. On a flag like "is this
	// connection live" or "may this provider create accounts", that turns an
	// operator's explicit "no" into a "yes". Defaults are applied in the handler
	// instead, where they can be expressed without fighting the ORM.
	Enabled bool ~gorm:"index" json:"enabled"~

	// JITProvisioning creates a user on first successful login. With it off, a
	// user who authenticates successfully but has no account is rejected —
	// which is what a customer who pre-provisions their users expects.
	JITProvisioning bool ~json:"jit_provisioning"~

	// DefaultRoleID is granted to users created by JIT provisioning when no
	// group mapping matches. Empty means the app's normal default.
	DefaultRoleID string ~gorm:"size:36" json:"default_role_id"~

	// GroupsClaim is the claim holding the user's IdP groups ("groups" for
	// Okta/Keycloak, "roles" for Entra ID app roles). GroupMappings is JSON of
	// {"<idp group>": "<role name>"}; every matching group's role is granted on
	// each login, so revoking a group in the IdP revokes the role here too.
	GroupsClaim   string ~gorm:"size:120;default:'groups'" json:"groups_claim"~
	GroupMappings string ~gorm:"type:text" json:"group_mappings"~

	// ── SAML 2.0 ────────────────────────────────────────────────────────────
	// The IdP's metadata document, either fetched from a URL or pasted. It
	// carries the sign-in URL and the certificate every assertion is verified
	// against, so this is the trust anchor for the whole connection.
	MetadataURL string ~gorm:"size:500" json:"metadata_url"~
	MetadataXML string ~gorm:"type:text" json:"metadata_xml,omitempty"~

	// SAML carries claims as named attributes, and every IdP names them
	// differently. Blank falls back to the common conventions.
	EmailAttribute     string ~gorm:"size:255" json:"email_attribute"~
	FirstNameAttribute string ~gorm:"size:255" json:"first_name_attribute"~
	LastNameAttribute  string ~gorm:"size:255" json:"last_name_attribute"~
	GroupsAttribute    string ~gorm:"size:255" json:"groups_attribute"~

	// AllowIDPInitiated accepts an assertion the app never asked for — the user
	// clicking the app tile in Okta rather than starting at our login page.
	// That is how most enterprise users actually sign in, and it avoids relying
	// on a cross-site cookie surviving the IdP's POST back. The assertion is
	// still signature-checked, audience-restricted and time-bounded; turn it off
	// if you require every login to begin at this app.
	//
	// The column is allow_id_p_initiated: that is what GORM made of the field
	// name when this shipped, and the update wrote allow_idp_initiated, so
	// turning the flag off failed with a 500 and left it on. Pinned, so nothing
	// depends on that guess again.
	AllowIDPInitiated bool ~gorm:"column:allow_id_p_initiated" json:"allow_idp_initiated"~

	LastUsedAt *time.Time ~json:"last_used_at"~

	CreatedAt time.Time      ~json:"created_at"~
	UpdatedAt time.Time      ~json:"updated_at"~
	DeletedAt gorm.DeletedAt ~gorm:"index" json:"-"~
}

// AfterFind surfaces whether a secret is stored without exposing it.
func (s *SSOConnection) AfterFind(tx *gorm.DB) error {
	s.HasSecret = strings.TrimSpace(string(s.ClientSecret)) != ""
	return nil
}

func (s *SSOConnection) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	s.Slug = strings.ToLower(strings.TrimSpace(s.Slug))
	return nil
}

// DomainList returns the connection's email domains, normalized.
func (s *SSOConnection) DomainList() []string {
	out := []string{}
	for _, d := range strings.Split(s.Domains, ",") {
		d = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(d, "@")))
		if d != "" {
			out = append(out, d)
		}
	}
	return out
}

// OwnsEmail reports whether an address is at one of the connection's domains:
// the only addresses its identity provider is trusted to vouch for.
func (s *SSOConnection) OwnsEmail(email string) bool {
	at := strings.LastIndex(email, "@")
	if at < 0 || at == len(email)-1 {
		return false
	}
	domain := strings.ToLower(strings.TrimSpace(email[at+1:]))
	for _, d := range s.DomainList() {
		if d == domain {
			return true
		}
	}
	return false
}

// ScopeList returns the requested scopes. "openid" is implied and always first.
func (s *SSOConnection) ScopeList() []string {
	out := []string{"openid"}
	for _, sc := range strings.Split(s.Scopes, ",") {
		sc = strings.TrimSpace(sc)
		if sc != "" && sc != "openid" {
			out = append(out, sc)
		}
	}
	return out
}

// IsSAML reports whether this connection speaks SAML rather than OIDC.
func (s *SSOConnection) IsSAML() bool {
	return strings.EqualFold(strings.TrimSpace(s.Protocol), "saml")
}

// AttributeOr returns the configured attribute name or a fallback list of the
// conventional names for that claim.
func (s *SSOConnection) AttributeOr(configured string, fallbacks ...string) []string {
	if c := strings.TrimSpace(configured); c != "" {
		return []string{c}
	}
	return fallbacks
}

// Discovery returns the OIDC discovery document URL.
func (s *SSOConnection) Discovery() string {
	if u := strings.TrimSpace(s.DiscoveryURL); u != "" {
		return u
	}
	return strings.TrimRight(strings.TrimSpace(s.IssuerURL), "/") + "/.well-known/openid-configuration"
}

// UserIdentity links a local user to an external identity.
//
// The subject ("sub") is the only identifier an IdP guarantees is stable —
// email addresses get reassigned when people leave and are renamed when people
// marry. Matching on sub first means a user who changes their email at the IdP
// keeps their account and their data, rather than silently getting a new one.
type UserIdentity struct {
	ID     string ~gorm:"primarykey;size:36" json:"id"~
	UserID string ~gorm:"size:36;index;not null" json:"user_id"~

	// Provider is the SSO connection slug (or "google"/"github" for social).
	Provider string ~gorm:"size:64;not null;uniqueIndex:idx_identity_provider_subject" json:"provider"~
	// Subject is the IdP's immutable identifier for this person.
	Subject string ~gorm:"size:255;not null;uniqueIndex:idx_identity_provider_subject" json:"subject"~

	Email      string     ~gorm:"size:255" json:"email"~
	LastLoginAt *time.Time ~json:"last_login_at"~

	CreatedAt time.Time ~json:"created_at"~
	UpdatedAt time.Time ~json:"updated_at"~
}

func (i *UserIdentity) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = ids.New()
	}
	return nil
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// apiSSOServiceGo emits internal/services/sso.go — the provider registry.
func apiSSOServiceGo() string {
	src := `package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/markbates/goth/providers/openidConnect"
	"gorm.io/gorm"

	"{{MODULE}}/internal/cluster"
	"{{MODULE}}/internal/models"
)

// SSORegistry holds one configured OIDC provider per enabled connection.
//
// It deliberately does NOT use goth's package-level provider map. That map is
// written by goth.UseProviders and read on every login with no lock, so an
// admin saving a connection while somebody signs in is a concurrent map
// read/write — which in Go is a fatal runtime error that takes the process
// down, not a race the detector merely complains about. Owning the map here
// with an RWMutex makes runtime reconfiguration safe.
//
// Building a provider performs OIDC discovery (a network call to the IdP), so
// providers are built once at boot and rebuilt only when a connection changes.
type SSORegistry struct {
	mu        sync.RWMutex
	providers map[string]*openidConnect.Provider
	appURL    string
}

// ErrSSOConnectionUnavailable is returned when a slug has no live provider —
// either it was never configured, it's disabled, or discovery failed at boot.
var ErrSSOConnectionUnavailable = errors.New("sso connection unavailable")

func NewSSORegistry(appURL string) *SSORegistry {
	return &SSORegistry{
		providers: map[string]*openidConnect.Provider{},
		appURL:    strings.TrimRight(appURL, "/"),
	}
}

// CallbackURL is where the IdP sends the user back. It is deliberately
// UNVERSIONED (/api/auth/..., not /api/v1/auth/...) because this exact string
// is registered as a redirect URI in the customer's IdP console — a value they
// control, not us. Versioning it would break every existing connection the day
// the API version bumps. The unversioned path is re-dispatched internally.
func (r *SSORegistry) CallbackURL(slug string) string {
	return r.appURL + "/api/auth/sso/" + slug + "/callback"
}

// Reload rebuilds the whole registry from the database. Called at boot and
// after any connection is created, updated or deleted. A connection whose
// discovery fails is logged and skipped rather than failing the others — one
// customer's misconfigured IdP must not stop everyone else signing in.
func (r *SSORegistry) Reload(db *gorm.DB) []error {
	var conns []models.SSOConnection
	// SAML connections have their own registry; building them here logged a
	// "client id and secret are required" error for each one on every reload.
	if err := db.Where("enabled = ? AND (protocol IS NULL OR protocol <> ?)", true, "saml").
		Find(&conns).Error; err != nil {
		return []error{fmt.Errorf("loading sso connections: %w", err)}
	}

	built := map[string]*openidConnect.Provider{}
	var errs []error
	for _, conn := range conns {
		p, err := r.build(conn)
		if err != nil {
			errs = append(errs, fmt.Errorf("sso %q: %w", conn.Slug, err))
			continue
		}
		built[conn.Slug] = p
	}

	r.mu.Lock()
	r.providers = built
	r.mu.Unlock()
	return errs
}

func (r *SSORegistry) build(conn models.SSOConnection) (*openidConnect.Provider, error) {
	secret := strings.TrimSpace(string(conn.ClientSecret))
	if conn.ClientID == "" || secret == "" {
		return nil, errors.New("client id and secret are required")
	}
	// NewNamed suffixes the name with "-oidc" internally; the slug is what we
	// key on here, so callers never have to know that.
	p, err := openidConnect.NewNamed(
		conn.Slug,
		conn.ClientID,
		secret,
		r.CallbackURL(conn.Slug),
		conn.Discovery(),
		conn.ScopeList()...,
	)
	if err != nil {
		return nil, fmt.Errorf("openid discovery: %w", err)
	}
	return p, nil
}

// Provider returns the live provider for a slug.
func (r *SSORegistry) Provider(slug string) (*openidConnect.Provider, error) {
	r.mu.RLock()
	p, ok := r.providers[slug]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrSSOConnectionUnavailable
	}
	return p, nil
}

// Count reports how many providers are live — used by the admin UI to
// distinguish "no connections" from "connections that all failed discovery".
func (r *SSORegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.providers)
}

// ConnectionForEmail resolves the connection an address should authenticate
// against, matching on the domain after "@". Returns nil when nothing matches,
// which the caller should treat as "fall back to password login" rather than an
// error — most users of most apps are not SSO users.
func ConnectionForEmail(db *gorm.DB, email string) (*models.SSOConnection, error) {
	at := strings.LastIndex(email, "@")
	if at < 0 || at == len(email)-1 {
		return nil, nil
	}
	domain := strings.ToLower(strings.TrimSpace(email[at+1:]))
	if domain == "" {
		return nil, nil
	}

	var conns []models.SSOConnection
	if err := db.Where("enabled = ?", true).Find(&conns).Error; err != nil {
		return nil, err
	}
	for i := range conns {
		for _, d := range conns[i].DomainList() {
			if d == domain {
				return &conns[i], nil
			}
		}
	}
	return nil, nil
}

// WatchSSO rebuilds this replica's SSO registries when another replica changes
// a connection. The registries are built in each process, so a connection
// created or disabled through one API instance was unknown to, or still live
// on, the others until they restarted.
func WatchSSO(db *gorm.DB, oidc *SSORegistry, samlReg *SAMLRegistry) {
	w := cluster.NewWatch(db, "sso", 2*time.Second)
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for range t.C {
			if !w.Changed() {
				continue
			}
			for _, err := range oidc.Reload(db) {
				log.Printf("sso: %v", err)
			}
			if samlReg != nil {
				for _, err := range samlReg.Reload(db) {
					log.Printf("saml: %v", err)
				}
			}
		}
	}()
}

// AnnounceSSOChange tells the other replicas to rebuild their registries.
func AnnounceSSOChange(db *gorm.DB) {
	if err := cluster.Bump(db, "sso"); err != nil {
		log.Printf("sso: could not tell the other replicas about a connection change: %v", err)
	}
}

// GroupRoleNames maps the IdP groups on a login to local role names, using the
// connection's GroupMappings. Unmapped groups are ignored — an IdP typically
// carries far more groups than an app has roles.
func GroupRoleNames(conn *models.SSOConnection, groups []string) []string {
	raw := strings.TrimSpace(conn.GroupMappings)
	if raw == "" || len(groups) == 0 {
		return nil
	}
	mapping := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &mapping); err != nil {
		// Fail closed: a malformed mapping grants nothing rather than
		// everything, matching how Role.GrantsList treats bad JSON.
		return nil
	}

	lower := map[string]string{}
	for k, v := range mapping {
		lower[strings.ToLower(strings.TrimSpace(k))] = v
	}

	seen := map[string]bool{}
	out := []string{}
	for _, g := range groups {
		if role, ok := lower[strings.ToLower(strings.TrimSpace(g))]; ok && !seen[role] {
			seen[role] = true
			out = append(out, role)
		}
	}
	return out
}

// ClaimGroups pulls the group list out of a raw claims map. IdPs disagree on
// shape — a JSON array for Okta/Keycloak, occasionally a single string — so
// both are accepted.
func ClaimGroups(raw map[string]interface{}, claim string) []string {
	if claim == "" {
		claim = "groups"
	}
	v, ok := raw[claim]
	if !ok {
		return nil
	}
	switch t := v.(type) {
	case []interface{}:
		out := []string{}
		for _, item := range t {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return t
	case string:
		if t == "" {
			return nil
		}
		return []string{t}
	}
	return nil
}

// TouchConnection records that a connection was just used to sign somebody in.
func TouchConnection(db *gorm.DB, id string) {
	now := time.Now()
	if err := db.Model(&models.SSOConnection{}).Where("id = ?", id).Update("last_used_at", now).Error; err != nil {
		log.Printf("sso: recording use of connection %s: %v", id, err)
	}
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// apiSSOHandlerGo emits internal/handlers/sso.go — the public login flow plus
// admin CRUD for connections.
func apiSSOHandlerGo() string { return tmpl("api/handlers/sso.go") }
