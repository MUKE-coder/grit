package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/markbates/goth/providers/openidConnect"
	"gorm.io/gorm"

	"library/apps/api/internal/cluster"
	"library/apps/api/internal/models"
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

// SSOService owns the connection and identity tables.
//
// The registry above holds the live providers; this holds the rows they are
// built from, and the external identities a sign-in matches against. The handler
// decides what a refusal says and which fields a request may set. What a
// connection is, and which account an assertion belongs to, is decided here, so
// a console command that provisions a connection and the admin screen that
// creates one apply the same rules.
type SSOService struct {
	DB *gorm.DB
}

func (s *SSOService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// All lists every connection, newest first, which is the order the admin screen
// shows them in.
func (s *SSOService) All(ctx context.Context) ([]models.SSOConnection, error) {
	var conns []models.SSOConnection
	if err := s.db(ctx).Order("created_at desc").Find(&conns).Error; err != nil {
		return nil, err
	}
	return conns, nil
}

// ByID reads one connection.
func (s *SSOService) ByID(ctx context.Context, id string) (*models.SSOConnection, error) {
	var conn models.SSOConnection
	if err := s.db(ctx).Where("id = ?", id).First(&conn).Error; err != nil {
		return nil, err
	}
	return &conn, nil
}

// EnabledBySlug reads the connection a callback belongs to, and only while it is
// enabled: a connection turned off mid-flow stops being a way in.
func (s *SSOService) EnabledBySlug(ctx context.Context, slug string) (*models.SSOConnection, error) {
	var conn models.SSOConnection
	if err := s.db(ctx).Where("slug = ? AND enabled = ?", slug, true).First(&conn).Error; err != nil {
		return nil, err
	}
	return &conn, nil
}

// EnabledBySlugAndProtocol is EnabledBySlug for a flow that only one protocol
// can answer: a SAML assertion arriving for a connection that has been switched
// to OIDC is not a sign-in, it is a stale bookmark.
func (s *SSOService) EnabledBySlugAndProtocol(ctx context.Context, slug, protocol string) (*models.SSOConnection, error) {
	var conn models.SSOConnection
	if err := s.db(ctx).Where("slug = ? AND enabled = ? AND protocol = ?", slug, true, protocol).
		First(&conn).Error; err != nil {
		return nil, err
	}
	return &conn, nil
}

// Create saves a new connection.
func (s *SSOService) Create(ctx context.Context, conn *models.SSOConnection) error {
	return s.db(ctx).Create(conn).Error
}

// Update writes the columns a request sent and returns the row as it now
// stands, because the admin screen renders what comes back and a connection
// carries fields the request never named.
func (s *SSOService) Update(ctx context.Context, conn *models.SSOConnection, updates map[string]interface{}) (*models.SSOConnection, error) {
	if err := s.db(ctx).Model(conn).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.ByID(ctx, conn.ID)
}

// Delete removes a connection.
func (s *SSOService) Delete(ctx context.Context, id string) error {
	return s.db(ctx).Where("id = ?", id).Delete(&models.SSOConnection{}).Error
}

// DomainClaimedElsewhere returns a domain in domains that another connection
// already claims, and that connection's name.
//
// Discovery sends an address to the first connection claiming its domain, so two
// claiming one made which identity provider a customer reaches depend on row
// order. The comparison is on the parsed domain list rather than the stored
// string, because "acme.com, acme.co.uk" and "acme.co.uk,acme.com" are the same
// claim.
func (s *SSOService) DomainClaimedElsewhere(ctx context.Context, domains, exceptID string) (string, string, error) {
	wanted := map[string]bool{}
	for _, d := range (&models.SSOConnection{Domains: domains}).DomainList() {
		wanted[d] = true
	}
	if len(wanted) == 0 {
		return "", "", nil
	}
	q := s.db(ctx)
	if exceptID != "" {
		q = q.Where("id <> ?", exceptID)
	}
	var others []models.SSOConnection
	if err := q.Find(&others).Error; err != nil {
		return "", "", err
	}
	for _, o := range others {
		for _, d := range o.DomainList() {
			if wanted[d] {
				return d, o.Name, nil
			}
		}
	}
	return "", "", nil
}

// IdentityBySubject reads the link between a connection and the subject its
// identity provider asserts.
//
// The subject is matched before the email because it is the only identifier an
// identity provider guarantees stable: somebody who changes their address there
// keeps their account rather than quietly getting a second one.
func (s *SSOService) IdentityBySubject(ctx context.Context, provider, subject string) (*models.UserIdentity, error) {
	var identity models.UserIdentity
	if err := s.db(ctx).Where("provider = ? AND subject = ?", provider, subject).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

// RecordIdentityLogin notes when a link was last used, and the address the
// identity provider asserted this time.
func (s *SSOService) RecordIdentityLogin(ctx context.Context, identity *models.UserIdentity, at time.Time, email string) error {
	return s.db(ctx).Model(identity).
		Updates(map[string]interface{}{"last_login_at": at, "email": email}).Error
}

// LinkIdentity records the link, so the next sign-in matches on the subject.
func (s *SSOService) LinkIdentity(ctx context.Context, link *models.UserIdentity) error {
	return s.db(ctx).Create(link).Error
}
