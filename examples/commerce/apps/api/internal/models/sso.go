package models

import (
	"strings"
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"

	"commerce/apps/api/internal/crypto"
)

// SSOConnection is one customer's identity provider.
//
// A connection is identified by its slug in URLs (/api/auth/sso/acme), and
// matched to a user by the domain of the email they type on the login page.
// Everything needed to talk to a compliant OIDC provider is here: the discovery
// document URL plus a client ID and secret issued by that provider.
type SSOConnection struct {
	ID   string `gorm:"primarykey;size:36" json:"id"`
	Slug string `gorm:"size:64;uniqueIndex;not null" json:"slug" binding:"required"`
	Name string `gorm:"size:255;not null" json:"name" binding:"required"`

	// Protocol is "oidc" (default) or "saml". OIDC needs an issuer plus client
	// credentials; SAML needs the IdP's metadata and no secret at all, because
	// trust rides on the IdP's signing certificate inside that metadata.
	Protocol string `gorm:"size:10;default:'oidc'" json:"protocol"`

	// Domains is a comma-separated list of email domains routed to this
	// connection ("acme.com,acme.co.uk"). A user typing bob@acme.com on the
	// login page is sent here. Matching is case-insensitive and exact on the
	// domain part — no wildcards, because "*.com" is a security incident.
	Domains string `gorm:"size:1000" json:"domains"`

	// IssuerURL is the provider's base issuer (e.g.
	// https://login.microsoftonline.com/<tenant>/v2.0 or https://acme.okta.com).
	// The discovery document is fetched from <issuer>/.well-known/openid-configuration
	// unless DiscoveryURL overrides it.
	IssuerURL    string `gorm:"size:500;not null" json:"issuer_url" binding:"required"`
	DiscoveryURL string `gorm:"size:500" json:"discovery_url"`

	ClientID string `gorm:"size:255;not null" json:"client_id" binding:"required"`
	// ClientSecret is encrypted at rest via the same AES-256-GCM field
	// encryption used for other PII, and never serialized — the API returns
	// HasSecret instead so the admin UI can show "configured" without ever
	// shipping the value back to a browser.
	ClientSecret crypto.EncryptedString `gorm:"type:text" json:"-"`
	HasSecret    bool                   `gorm:"-" json:"has_secret"`

	// Scopes requested from the IdP. "openid" is always sent; profile and email
	// are the defaults because the callback needs a name and an address.
	Scopes string `gorm:"size:500;default:'profile,email'" json:"scopes"`

	// NOTE: no `default:true` on these two booleans, deliberately.
	//
	// GORM omits zero-valued fields from an INSERT when the column carries a
	// default, so a `default:true` bool can never be stored as false through a
	// struct create — it silently comes back true. On a flag like "is this
	// connection live" or "may this provider create accounts", that turns an
	// operator's explicit "no" into a "yes". Defaults are applied in the handler
	// instead, where they can be expressed without fighting the ORM.
	Enabled bool `gorm:"index" json:"enabled"`

	// JITProvisioning creates a user on first successful login. With it off, a
	// user who authenticates successfully but has no account is rejected —
	// which is what a customer who pre-provisions their users expects.
	JITProvisioning bool `json:"jit_provisioning"`

	// DefaultRoleID is granted to users created by JIT provisioning when no
	// group mapping matches. Empty means the app's normal default.
	DefaultRoleID string `gorm:"size:36" json:"default_role_id"`

	// GroupsClaim is the claim holding the user's IdP groups ("groups" for
	// Okta/Keycloak, "roles" for Entra ID app roles). GroupMappings is JSON of
	// {"<idp group>": "<role name>"}; every matching group's role is granted on
	// each login, so revoking a group in the IdP revokes the role here too.
	GroupsClaim   string `gorm:"size:120;default:'groups'" json:"groups_claim"`
	GroupMappings string `gorm:"type:text" json:"group_mappings"`

	// ── SAML 2.0 ────────────────────────────────────────────────────────────
	// The IdP's metadata document, either fetched from a URL or pasted. It
	// carries the sign-in URL and the certificate every assertion is verified
	// against, so this is the trust anchor for the whole connection.
	MetadataURL string `gorm:"size:500" json:"metadata_url"`
	MetadataXML string `gorm:"type:text" json:"metadata_xml,omitempty"`

	// SAML carries claims as named attributes, and every IdP names them
	// differently. Blank falls back to the common conventions.
	EmailAttribute     string `gorm:"size:255" json:"email_attribute"`
	FirstNameAttribute string `gorm:"size:255" json:"first_name_attribute"`
	LastNameAttribute  string `gorm:"size:255" json:"last_name_attribute"`
	GroupsAttribute    string `gorm:"size:255" json:"groups_attribute"`

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
	AllowIDPInitiated bool `gorm:"column:allow_id_p_initiated" json:"allow_idp_initiated"`

	LastUsedAt *time.Time `json:"last_used_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
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
	ID     string `gorm:"primarykey;size:36" json:"id"`
	UserID string `gorm:"size:36;index;not null" json:"user_id"`

	// Provider is the SSO connection slug (or "google"/"github" for social).
	Provider string `gorm:"size:64;not null;uniqueIndex:idx_identity_provider_subject" json:"provider"`
	// Subject is the IdP's immutable identifier for this person.
	Subject string `gorm:"size:255;not null;uniqueIndex:idx_identity_provider_subject" json:"subject"`

	Email       string     `gorm:"size:255" json:"email"`
	LastLoginAt *time.Time `json:"last_login_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *UserIdentity) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = ids.New()
	}
	return nil
}
