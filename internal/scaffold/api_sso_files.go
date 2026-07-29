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
	AllowIDPInitiated bool ~json:"allow_idp_initiated"~

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
	"strings"
	"sync"
	"time"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"gorm.io/gorm"

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
	if err := db.Where("enabled = ?", true).Find(&conns).Error; err != nil {
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
	db.Model(&models.SSOConnection{}).Where("id = ?", id).Update("last_used_at", now)
}

// providerSession is the goth session marshalled between the redirect to the
// IdP and the callback. Kept as its own type so the handler can be explicit
// about what it stores in the cookie.
type providerSession = goth.Session
`
	return strings.ReplaceAll(src, "~", "`")
}

// apiSSOHandlerGo emits internal/handlers/sso.go — the public login flow plus
// admin CRUD for connections.
func apiSSOHandlerGo() string {
	src := `package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/crypto"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/respond"
	"{{MODULE}}/internal/services"
)

// ssoStateCookie carries the marshalled provider session plus the CSRF state
// between the redirect out to the IdP and the callback. It is HttpOnly and
// short-lived: it exists only for the seconds the user spends at the IdP.
// externalIdentity is one person as an identity provider described them,
// normalized so OIDC and SAML converge before any account is touched. Keeping
// provisioning, identity linking and role mapping on this one shape means SAML
// inherits the behaviour OIDC already has tests for, instead of growing a
// second, subtly different copy.
type externalIdentity struct {
	Subject   string // the IdP's immutable ID: OIDC "sub", SAML NameID
	Email     string
	FirstName string
	LastName  string
	Avatar    string
	Groups    []string
}

const ssoStateCookie = "grit_sso"

// ssoCookiePath is deliberately "/api" rather than "/api/auth".
//
// The redirect out to the IdP and the return trip can legitimately arrive on
// either the versioned path (/api/v1/auth/sso/...) or the unversioned alias
// (/api/auth/sso/..., which is what gets registered in the IdP console). A
// cookie scoped to /api/auth is not sent to /api/v1/auth/..., so scoping it
// that tightly makes the flow work or break depending on which URL the caller
// happened to use. The cookie is HttpOnly and lives ~10 minutes, so the wider
// path costs nothing.
const ssoCookiePath = "/api"

type SSOHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Config      *config.Config
	Registry    *services.SSORegistry
	SAML        *services.SAMLRegistry
}

func NewSSOHandler(db *gorm.DB, auth *services.AuthService, cfg *config.Config,
	reg *services.SSORegistry, samlReg *services.SAMLRegistry) *SSOHandler {
	return &SSOHandler{DB: db, AuthService: auth, Config: cfg, Registry: reg, SAML: samlReg}
}

// ── Public: discovery ────────────────────────────────────────────────────────

// Discover answers "should this email address use SSO, and if so where?".
//
// The login form calls it as the user submits their address. A miss is a normal
// answer, not an error — most people at most apps sign in with a password.
//
//	POST /api/auth/sso/discover  {"email": "bob@acme.com"}
//	200 {"data": {"sso": true, "slug": "acme", "name": "Acme Okta", "redirect_url": "..."}}
func (h *SSOHandler) Discover(c *gin.Context) {
	var req struct {
		Email string ~json:"email" binding:"required,email"~
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, "a valid email address is required")
		return
	}

	conn, err := services.ConnectionForEmail(h.DB, req.Email)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	if conn == nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"sso": false}})
		return
	}

	// The two protocols start at different URLs, so the caller is told which
	// one to send the browser to rather than having to infer it.
	start := "/api/auth/sso/" + conn.Slug
	if conn.IsSAML() {
		start = "/api/auth/saml/" + conn.Slug
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"sso":          true,
		"slug":         conn.Slug,
		"name":         conn.Name,
		"protocol":     conn.Protocol,
		"redirect_url": start,
	}})
}

// ── Public: the login flow ───────────────────────────────────────────────────

// Begin redirects the user to their identity provider.
//
//	GET /api/auth/sso/:slug
func (h *SSOHandler) Begin(c *gin.Context) {
	slug := strings.ToLower(c.Param("slug"))

	provider, err := h.Registry.Provider(slug)
	if err != nil {
		h.failLogin(c, "That sign-in method is not available.")
		return
	}

	state, err := randomState()
	if err != nil {
		h.failLogin(c, "Could not start sign-in. Please try again.")
		return
	}

	sess, err := provider.BeginAuth(state)
	if err != nil {
		log.Printf("sso %s: begin: %v", slug, err)
		h.failLogin(c, "Could not reach the identity provider.")
		return
	}
	authURL, err := sess.GetAuthURL()
	if err != nil {
		log.Printf("sso %s: auth url: %v", slug, err)
		h.failLogin(c, "Could not reach the identity provider.")
		return
	}

	// The session and the state travel together in one HttpOnly cookie. The
	// state is echoed back by the IdP and compared on return, which is what
	// stops an attacker replaying someone else's callback.
	payload, _ := json.Marshal(map[string]string{
		"slug":    slug,
		"state":   state,
		"session": sess.Marshal(),
	})
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(ssoStateCookie, base64.RawURLEncoding.EncodeToString(payload),
		600, ssoCookiePath, "", isSecureRequest(c), true)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback completes the login: verifies the response, resolves or provisions
// the user, applies role mapping, and issues the same cookies a password login
// would.
//
//	GET /api/auth/sso/:slug/callback
func (h *SSOHandler) Callback(c *gin.Context) {
	slug := strings.ToLower(c.Param("slug"))

	raw, err := c.Cookie(ssoStateCookie)
	if err != nil {
		h.failLogin(c, "Your sign-in session expired. Please try again.")
		return
	}
	// One-shot: clear it whatever happens next.
	c.SetCookie(ssoStateCookie, "", -1, ssoCookiePath, "", isSecureRequest(c), true)

	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		h.failLogin(c, "Your sign-in session was invalid. Please try again.")
		return
	}
	var stash map[string]string
	if err := json.Unmarshal(decoded, &stash); err != nil {
		h.failLogin(c, "Your sign-in session was invalid. Please try again.")
		return
	}

	if stash["slug"] != slug {
		h.failLogin(c, "Your sign-in session did not match. Please try again.")
		return
	}
	if s := c.Query("state"); s == "" || s != stash["state"] {
		h.failLogin(c, "Your sign-in session did not match. Please try again.")
		return
	}

	provider, err := h.Registry.Provider(slug)
	if err != nil {
		h.failLogin(c, "That sign-in method is not available.")
		return
	}

	sess, err := provider.UnmarshalSession(stash["session"])
	if err != nil {
		h.failLogin(c, "Your sign-in session was invalid. Please try again.")
		return
	}
	if _, err := sess.Authorize(provider, c.Request.URL.Query()); err != nil {
		log.Printf("sso %s: authorize: %v", slug, err)
		h.failLogin(c, "Sign-in was not completed.")
		return
	}

	external, err := provider.FetchUser(sess)
	if err != nil {
		log.Printf("sso %s: fetch user: %v", slug, err)
		h.failLogin(c, "Could not read your profile from the identity provider.")
		return
	}
	if strings.TrimSpace(external.Email) == "" {
		h.failLogin(c, "Your identity provider did not release an email address.")
		return
	}

	var conn models.SSOConnection
	if err := h.DB.Where("slug = ? AND enabled = ?", slug, true).First(&conn).Error; err != nil {
		h.failLogin(c, "That sign-in method is not available.")
		return
	}

	ident := externalIdentity{
		Subject:   external.UserID,
		Email:     external.Email,
		FirstName: firstNonBlank(external.FirstName, external.NickName),
		LastName:  external.LastName,
		Avatar:    external.AvatarURL,
		Groups:    services.ClaimGroups(external.RawData, conn.GroupsClaim),
	}

	user, err := h.resolveUser(c, &conn, ident)
	if err != nil {
		h.failLogin(c, err.Error())
		return
	}

	if err := h.applyGroupRoles(&conn, user, ident); err != nil {
		// Role mapping failing shouldn't strand a user who authenticated
		// correctly — log it and let them in with whatever they already have.
		log.Printf("sso %s: role mapping for %s: %v", slug, user.ID, err)
	}

	// Reload so the token carries any role the mapping just assigned.
	if err := h.DB.Where("id = ?", user.ID).First(user).Error; err != nil {
		h.failLogin(c, "Could not complete sign-in.")
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		h.failLogin(c, "Could not complete sign-in.")
		return
	}
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("sso: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	services.TouchConnection(h.DB, conn.ID)

	c.Redirect(http.StatusTemporaryRedirect, h.Config.OAuthFrontendURL+"/auth/callback")
}

// resolveUser finds the account this login belongs to, provisioning one when
// the connection allows it.
//
// Order matters. The IdP subject is matched first because it is the only
// identifier guaranteed stable — someone who changes their email at the IdP
// must keep their account rather than silently getting a second one. Email is
// the fallback, which is also how an existing password user gets linked to SSO
// the first time their company turns it on.
func (h *SSOHandler) resolveUser(c *gin.Context, conn *models.SSOConnection, external externalIdentity) (*models.User, error) {
	now := time.Now()
	email := strings.ToLower(strings.TrimSpace(external.Email))

	var identity models.UserIdentity
	err := h.DB.Where("provider = ? AND subject = ?", conn.Slug, external.Subject).First(&identity).Error
	if err == nil {
		var user models.User
		if err := h.DB.Where("id = ?", identity.UserID).First(&user).Error; err != nil {
			return nil, fmt.Errorf("Your account could not be found.")
		}
		if !user.Active {
			return nil, fmt.Errorf("Your account has been disabled.")
		}
		h.DB.Model(&identity).Updates(map[string]interface{}{"last_login_at": now, "email": email})
		return &user, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("Could not complete sign-in.")
	}

	var user models.User
	err = h.DB.Where("email = ?", email).First(&user).Error
	switch {
	case err == nil:
		if !user.Active {
			return nil, fmt.Errorf("Your account has been disabled.")
		}
	case err == gorm.ErrRecordNotFound:
		if !conn.JITProvisioning {
			return nil, fmt.Errorf("No account exists for %s. Ask your administrator to create one.", email)
		}
		user = models.User{
			FirstName:       firstNonBlank(external.FirstName, "User"),
			LastName:        external.LastName,
			Email:           email,
			Avatar:          external.Avatar,
			Provider:        conn.Slug,
			Role:            models.RoleUser,
			Active:          true,
			EmailVerifiedAt: &now, // the IdP asserted it
			IPAddress:       c.ClientIP(),
		}
		if err := h.DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("Could not create your account.")
		}
	default:
		return nil, fmt.Errorf("Could not complete sign-in.")
	}

	// Link the external identity so subsequent logins match on subject.
	link := models.UserIdentity{
		UserID:      user.ID,
		Provider:    conn.Slug,
		Subject:     external.Subject,
		Email:       email,
		LastLoginAt: &now,
	}
	if err := h.DB.Create(&link).Error; err != nil {
		log.Printf("sso %s: linking identity for %s: %v", conn.Slug, user.ID, err)
	}
	return &user, nil
}

// applyGroupRoles grants the roles the IdP's groups map to.
//
// Roles are replaced, not merged, so removing someone from a group at the IdP
// removes the role here on their next login — which is the only reason to map
// groups at all. Connections with no mapping configured are left alone so an
// admin's manual grants survive.
func (h *SSOHandler) applyGroupRoles(conn *models.SSOConnection, user *models.User, external externalIdentity) error {
	names := services.GroupRoleNames(conn, external.Groups)

	if len(names) == 0 {
		if strings.TrimSpace(conn.GroupMappings) != "" {
			return nil // mapping configured but nothing matched — leave as-is
		}
		if conn.DefaultRoleID == "" {
			return nil
		}
		var role models.Role
		if err := h.DB.Where("id = ?", conn.DefaultRoleID).First(&role).Error; err != nil {
			return err
		}
		names = []string{role.Name}
	}

	var roles []models.Role
	if err := h.DB.Where("name IN ?", names).Find(&roles).Error; err != nil {
		return err
	}
	if len(roles) == 0 {
		return nil
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		for _, r := range roles {
			if err := tx.Create(&models.UserRole{UserID: user.ID, RoleID: r.ID}).Error; err != nil {
				return err
			}
		}
		// The legacy users.role string is what GenerateTokenPair bakes into the
		// JWT, so it has to move in step with the join table or the token and
		// the grants disagree.
		return tx.Model(&models.User{}).Where("id = ?", user.ID).
			Update("role", strings.ToUpper(roles[0].Name)).Error
	})
	if err != nil {
		return err
	}
	authz.Invalidate()
	return nil
}

func (h *SSOHandler) failLogin(c *gin.Context, message string) {
	c.Redirect(http.StatusTemporaryRedirect,
		fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape(message)))
}

// ── Admin CRUD ───────────────────────────────────────────────────────────────

type ssoConnectionInput struct {
	Slug            string ~json:"slug"~
	Name            string ~json:"name"~
	Domains         string ~json:"domains"~
	IssuerURL       string ~json:"issuer_url"~
	DiscoveryURL    string ~json:"discovery_url"~
	ClientID        string ~json:"client_id"~
	ClientSecret    string ~json:"client_secret"~
	Scopes          string ~json:"scopes"~
	Enabled         *bool  ~json:"enabled"~
	JITProvisioning *bool  ~json:"jit_provisioning"~
	DefaultRoleID   string ~json:"default_role_id"~
	GroupsClaim     string ~json:"groups_claim"~
	GroupMappings   string ~json:"group_mappings"~

	Protocol           string ~json:"protocol"~
	MetadataURL        string ~json:"metadata_url"~
	MetadataXML        string ~json:"metadata_xml"~
	EmailAttribute     string ~json:"email_attribute"~
	FirstNameAttribute string ~json:"first_name_attribute"~
	LastNameAttribute  string ~json:"last_name_attribute"~
	GroupsAttribute    string ~json:"groups_attribute"~
	AllowIDPInitiated  *bool  ~json:"allow_idp_initiated"~
}

func (h *SSOHandler) List(c *gin.Context) {
	var conns []models.SSOConnection
	if err := h.DB.Order("created_at desc").Find(&conns).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": conns,
		"meta": gin.H{"live": h.Registry.Count()},
	})
}

func (h *SSOHandler) Create(c *gin.Context) {
	var in ssoConnectionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, err.Error())
		return
	}
	protocol := strings.ToLower(strings.TrimSpace(in.Protocol))
	if protocol == "" {
		protocol = "oidc"
	}
	if in.Slug == "" || in.Name == "" {
		respond.BadRequest(c, "slug and name are required")
		return
	}
	// The two protocols need different things: OIDC needs client credentials,
	// SAML needs the IdP metadata and no secret at all.
	if protocol == "saml" {
		if in.MetadataURL == "" && in.MetadataXML == "" {
			respond.BadRequest(c, "a metadata URL or metadata XML is required for SAML")
			return
		}
	} else if in.IssuerURL == "" || in.ClientID == "" || in.ClientSecret == "" {
		respond.BadRequest(c, "issuer_url, client_id and client_secret are required for OIDC")
		return
	}

	conn := models.SSOConnection{
		Protocol:        protocol,
		MetadataURL:     in.MetadataURL,
		MetadataXML:     in.MetadataXML,
		EmailAttribute:  in.EmailAttribute,
		FirstNameAttribute: in.FirstNameAttribute,
		LastNameAttribute:  in.LastNameAttribute,
		GroupsAttribute:    in.GroupsAttribute,
		AllowIDPInitiated:  in.AllowIDPInitiated == nil || *in.AllowIDPInitiated,
		Slug:            strings.ToLower(strings.TrimSpace(in.Slug)),
		Name:            in.Name,
		Domains:         in.Domains,
		IssuerURL:       in.IssuerURL,
		DiscoveryURL:    in.DiscoveryURL,
		ClientID:        in.ClientID,
		ClientSecret:    crypto.EncryptedString(in.ClientSecret),
		Scopes:          firstNonBlank(in.Scopes, "profile,email"),
		Enabled:         in.Enabled == nil || *in.Enabled,
		JITProvisioning: in.JITProvisioning == nil || *in.JITProvisioning,
		DefaultRoleID:   in.DefaultRoleID,
		GroupsClaim:     firstNonBlank(in.GroupsClaim, "groups"),
		GroupMappings:   in.GroupMappings,
	}
	if err := h.DB.Create(&conn).Error; err != nil {
		respond.BadRequest(c, "could not save the connection — is the slug already in use?")
		return
	}
	// AfterFind hasn't run on a freshly-created struct, so set the flag the UI
	// reads rather than reporting "no secret" on the record we just stored one on.
	conn.HasSecret = true

	h.reload()
	c.JSON(http.StatusCreated, gin.H{"data": conn, "message": "SSO connection created"})
}

func (h *SSOHandler) Update(c *gin.Context) {
	var conn models.SSOConnection
	if err := h.DB.Where("id = ?", c.Param("id")).First(&conn).Error; err != nil {
		respond.NotFound(c, "SSO connection not found")
		return
	}

	var in ssoConnectionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Domains != "" {
		updates["domains"] = in.Domains
	}
	if in.IssuerURL != "" {
		updates["issuer_url"] = in.IssuerURL
	}
	updates["discovery_url"] = in.DiscoveryURL
	if in.ClientID != "" {
		updates["client_id"] = in.ClientID
	}
	// An empty secret means "leave the stored one alone" — the UI never has the
	// current value to send back, so treating blank as a clear would wipe it on
	// every unrelated edit.
	if in.ClientSecret != "" {
		updates["client_secret"] = crypto.EncryptedString(in.ClientSecret)
	}
	if in.Scopes != "" {
		updates["scopes"] = in.Scopes
	}
	if in.Enabled != nil {
		updates["enabled"] = *in.Enabled
	}
	if in.JITProvisioning != nil {
		updates["jit_provisioning"] = *in.JITProvisioning
	}
	updates["default_role_id"] = in.DefaultRoleID
	if in.GroupsClaim != "" {
		updates["groups_claim"] = in.GroupsClaim
	}
	updates["group_mappings"] = in.GroupMappings
	updates["metadata_url"] = in.MetadataURL
	if in.MetadataXML != "" {
		updates["metadata_xml"] = in.MetadataXML
	}
	updates["email_attribute"] = in.EmailAttribute
	updates["first_name_attribute"] = in.FirstNameAttribute
	updates["last_name_attribute"] = in.LastNameAttribute
	updates["groups_attribute"] = in.GroupsAttribute
	if in.AllowIDPInitiated != nil {
		updates["allow_idp_initiated"] = *in.AllowIDPInitiated
	}

	if err := h.DB.Model(&conn).Updates(updates).Error; err != nil {
		respond.Internal(c, err)
		return
	}

	h.reload()
	h.DB.Where("id = ?", conn.ID).First(&conn)
	c.JSON(http.StatusOK, gin.H{"data": conn, "message": "SSO connection updated"})
}

func (h *SSOHandler) Delete(c *gin.Context) {
	if err := h.DB.Where("id = ?", c.Param("id")).Delete(&models.SSOConnection{}).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	h.reload()
	c.JSON(http.StatusOK, gin.H{"message": "SSO connection deleted"})
}

// Test re-runs discovery for one connection so an admin can tell a typo in the
// issuer URL from a wrong client secret without waiting for a user to fail.
func (h *SSOHandler) Test(c *gin.Context) {
	var conn models.SSOConnection
	if err := h.DB.Where("id = ?", c.Param("id")).First(&conn).Error; err != nil {
		respond.NotFound(c, "SSO connection not found")
		return
	}
	if _, err := h.Registry.Provider(conn.Slug); err != nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"ok":      false,
			"message": "Not live. Check the issuer URL and credentials, then save again.",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"ok":           true,
		"message":      "Connection is live.",
		"callback_url": h.Registry.CallbackURL(conn.Slug),
	}})
}

func (h *SSOHandler) reload() {
	for _, err := range h.Registry.Reload(h.DB) {
		log.Printf("sso: %v", err)
	}
	if h.SAML != nil {
		for _, err := range h.SAML.Reload(h.DB) {
			log.Printf("sso: %v", err)
		}
	}
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func firstNonBlank(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func isSecureRequest(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
`
	return strings.ReplaceAll(src, "~", "`")
}
