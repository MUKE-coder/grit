package scaffold

import "strings"

// SAML 2.0 service-provider support.
//
// SAML predates OIDC and is clunkier — XML, signed assertions, certificates —
// but it is still what a lot of enterprise procurement asks for, and a customer
// whose IdP only speaks SAML cannot use an OIDC-only app at all.
//
// Everything after the assertion is verified is shared with OIDC: both
// protocols normalize onto externalIdentity and go through the same
// provisioning, identity-linking and role-mapping path, so SAML inherits the
// behaviour that path already has tests for.

// apiSAMLModelGo emits internal/models/saml.go — the service provider's own
// signing keypair.
func apiSAMLModelGo() string {
	src := `package models

import (
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"

	"{{MODULE}}/internal/crypto"
)

// SAMLKeypair is this application's service-provider identity.
//
// SAML requires the SP to have an X.509 keypair: it signs the authentication
// requests it sends, and publishes the certificate in its metadata so the IdP
// can verify them. One keypair serves every connection — it identifies the
// application, not the customer.
//
// It is generated on first use rather than asked for during setup, because
// "generate a self-signed certificate" is the step where SAML integrations
// stall. The private key is encrypted at rest with the same AES-256-GCM field
// encryption used for other secrets, and never leaves the server.
type SAMLKeypair struct {
	ID string ~gorm:"primarykey;size:36" json:"id"~

	// CertPEM is public — it goes in the SP metadata the customer uploads to
	// their IdP.
	CertPEM string ~gorm:"type:text;not null" json:"cert_pem"~
	// KeyPEM is the private half. Encrypted at rest, never serialized.
	KeyPEM crypto.EncryptedString ~gorm:"type:text;not null" json:"-"~

	NotAfter  time.Time ~json:"not_after"~
	CreatedAt time.Time ~json:"created_at"~
	UpdatedAt time.Time ~json:"updated_at"~
}

func (k *SAMLKeypair) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = ids.New()
	}
	return nil
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// apiSAMLServiceGo emits internal/services/saml.go.
func apiSAMLServiceGo() string {
	src := `package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	dsig "github.com/russellhaering/goxmldsig"
	"gorm.io/gorm"

	"{{MODULE}}/internal/crypto"
	"{{MODULE}}/internal/models"
)

// ErrSAMLUnavailable is returned when a slug has no live SAML provider.
var ErrSAMLUnavailable = errors.New("saml connection unavailable")

// SAMLRegistry holds one configured saml.ServiceProvider per SAML connection.
//
// Like the OIDC registry it owns its map behind an RWMutex rather than relying
// on package-level state, so an admin saving a connection while somebody signs
// in can't race a login.
type SAMLRegistry struct {
	mu        sync.RWMutex
	providers map[string]*saml.ServiceProvider
	appURL    string
}

func NewSAMLRegistry(appURL string) *SAMLRegistry {
	return &SAMLRegistry{
		providers: map[string]*saml.ServiceProvider{},
		appURL:    strings.TrimRight(appURL, "/"),
	}
}

// ACSURL is where the IdP POSTs its assertion, and MetadataURL is what the
// customer uploads to their IdP. Both are deliberately UNVERSIONED — they are
// registered in the customer's IdP console, so an API version bump must not
// change them. The unversioned paths are re-dispatched internally.
func (r *SAMLRegistry) ACSURL(slug string) string {
	return r.appURL + "/api/auth/saml/" + slug + "/acs"
}

func (r *SAMLRegistry) MetadataURL(slug string) string {
	return r.appURL + "/api/auth/saml/" + slug + "/metadata"
}

// EntityID identifies this application to the IdP. It is the metadata URL by
// convention, and must stay stable — changing it makes every existing IdP
// configuration stop matching.
func (r *SAMLRegistry) EntityID(slug string) string {
	return r.MetadataURL(slug)
}

// Reload rebuilds every SAML provider from the database.
func (r *SAMLRegistry) Reload(db *gorm.DB) []error {
	kp, err := EnsureSAMLKeypair(db)
	if err != nil {
		return []error{fmt.Errorf("saml keypair: %w", err)}
	}

	var conns []models.SSOConnection
	if err := db.Where("enabled = ? AND protocol = ?", true, "saml").Find(&conns).Error; err != nil {
		return []error{fmt.Errorf("loading saml connections: %w", err)}
	}

	built := map[string]*saml.ServiceProvider{}
	var errs []error
	for _, conn := range conns {
		sp, err := r.build(conn, kp)
		if err != nil {
			errs = append(errs, fmt.Errorf("saml %q: %w", conn.Slug, err))
			continue
		}
		built[conn.Slug] = sp
	}

	r.mu.Lock()
	r.providers = built
	r.mu.Unlock()
	return errs
}

func (r *SAMLRegistry) build(conn models.SSOConnection, kp *loadedKeypair) (*saml.ServiceProvider, error) {
	idp, err := loadIDPMetadata(conn)
	if err != nil {
		return nil, err
	}

	acs, err := url.Parse(r.ACSURL(conn.Slug))
	if err != nil {
		return nil, err
	}
	meta, err := url.Parse(r.MetadataURL(conn.Slug))
	if err != nil {
		return nil, err
	}

	return &saml.ServiceProvider{
		EntityID:    r.EntityID(conn.Slug),
		Key:         kp.Key,
		Certificate: kp.Cert,
		MetadataURL: *meta,
		AcsURL:      *acs,
		IDPMetadata: idp,
		// Sign our authentication requests. Without this the library publishes
		// only an encryption certificate in the SP metadata, so an IdP
		// configured to require signed requests rejects every login and the
		// operator has no way to fix it from this end. RSA-SHA256 is the
		// current interoperable default.
		SignatureMethod:   dsig.RSASHA256SignatureMethod,
		AllowIDPInitiated: conn.AllowIDPInitiated,
	}, nil
}

// loadIDPMetadata resolves the IdP's metadata: pasted XML wins over a URL,
// because a pasted document is an explicit, pinned choice while a URL can start
// serving something different tomorrow.
func loadIDPMetadata(conn models.SSOConnection) (*saml.EntityDescriptor, error) {
	if xmlDoc := strings.TrimSpace(conn.MetadataXML); xmlDoc != "" {
		return samlsp.ParseMetadata([]byte(xmlDoc))
	}
	raw := strings.TrimSpace(conn.MetadataURL)
	if raw == "" {
		return nil, errors.New("either metadata XML or a metadata URL is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("metadata url: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return samlsp.FetchMetadata(ctx, http.DefaultClient, *u)
}

// Provider returns the live service provider for a slug.
func (r *SAMLRegistry) Provider(slug string) (*saml.ServiceProvider, error) {
	r.mu.RLock()
	sp, ok := r.providers[slug]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrSAMLUnavailable
	}
	return sp, nil
}

func (r *SAMLRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.providers)
}

type loadedKeypair struct {
	Key  *rsa.PrivateKey
	Cert *x509.Certificate
}

var keypairOnce sync.Mutex

// EnsureSAMLKeypair returns the app's SP keypair, generating one the first time
// a SAML connection needs it. Guarded so two simultaneous reloads can't both
// decide to generate and leave two keypairs behind — the certificate is
// published to customers' IdPs, so it silently changing would break every
// existing connection.
func EnsureSAMLKeypair(db *gorm.DB) (*loadedKeypair, error) {
	keypairOnce.Lock()
	defer keypairOnce.Unlock()

	var row models.SAMLKeypair
	err := db.Order("created_at desc").First(&row).Error
	if err == nil {
		return parseKeypair(row)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	key, cert, err := generateSelfSigned()
	if err != nil {
		return nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	row = models.SAMLKeypair{
		CertPEM:  string(certPEM),
		KeyPEM:   crypto.EncryptedString(keyPEM),
		NotAfter: cert.NotAfter,
	}
	if err := db.Create(&row).Error; err != nil {
		return nil, err
	}
	return &loadedKeypair{Key: key, Cert: cert}, nil
}

func parseKeypair(row models.SAMLKeypair) (*loadedKeypair, error) {
	certBlock, _ := pem.Decode([]byte(row.CertPEM))
	if certBlock == nil {
		return nil, errors.New("stored certificate is not valid PEM")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode([]byte(row.KeyPEM))
	if keyBlock == nil {
		return nil, errors.New("stored private key is not valid PEM")
	}
	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &loadedKeypair{Key: key, Cert: cert}, nil
}

// generateSelfSigned mints the SP certificate. Self-signed is correct here: the
// IdP trusts it because the customer's admin uploaded this exact certificate,
// not because a public CA vouched for it.
func generateSelfSigned() (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}
	tmpl := x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "grit-saml-sp"},
		NotBefore:             time.Now().Add(-5 * time.Minute),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}

// SAMLAttribute pulls the first value of whichever attribute name matches.
// IdPs disagree wildly on naming — Okta sends "email", Entra ID sends the long
// Microsoft claim URI — so callers pass every name worth trying.
func SAMLAttribute(assertion *saml.Assertion, names ...string) string {
	for _, want := range names {
		w := strings.ToLower(strings.TrimSpace(want))
		if w == "" {
			continue
		}
		for _, st := range assertion.AttributeStatements {
			for _, attr := range st.Attributes {
				if strings.ToLower(attr.Name) == w || strings.ToLower(attr.FriendlyName) == w {
					for _, v := range attr.Values {
						if strings.TrimSpace(v.Value) != "" {
							return v.Value
						}
					}
				}
			}
		}
	}
	return ""
}

// SAMLAttributeValues returns every value of the first matching attribute —
// group membership arrives as repeated values rather than a list.
func SAMLAttributeValues(assertion *saml.Assertion, names ...string) []string {
	for _, want := range names {
		w := strings.ToLower(strings.TrimSpace(want))
		if w == "" {
			continue
		}
		for _, st := range assertion.AttributeStatements {
			for _, attr := range st.Attributes {
				if strings.ToLower(attr.Name) == w || strings.ToLower(attr.FriendlyName) == w {
					out := []string{}
					for _, v := range attr.Values {
						if strings.TrimSpace(v.Value) != "" {
							out = append(out, v.Value)
						}
					}
					if len(out) > 0 {
						return out
					}
				}
			}
		}
	}
	return nil
}

// SAMLSubject is the IdP's stable identifier for the person. NameID is the
// SAML equivalent of OIDC's "sub"; falling back to the email attribute keeps
// logins working with IdPs configured to send a transient NameID, at the cost
// of the rename-safety a persistent identifier gives.
func SAMLSubject(assertion *saml.Assertion, emailNames ...string) string {
	if assertion.Subject != nil && assertion.Subject.NameID != nil {
		if v := strings.TrimSpace(assertion.Subject.NameID.Value); v != "" {
			return v
		}
	}
	return SAMLAttribute(assertion, emailNames...)
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// apiSAMLHandlerGo emits internal/handlers/saml.go.
func apiSAMLHandlerGo() string {
	src := `package handlers

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/crewjam/saml"
	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// SAMLMetadata publishes this application's service-provider metadata.
//
// The customer's IdP admin uploads this document (or its URL) to create the
// application on their side — it carries the entity ID, the ACS endpoint and
// the certificate they must verify our requests against.
//
//	GET /api/auth/saml/:slug/metadata
func (h *SSOHandler) SAMLMetadata(c *gin.Context) {
	slug := strings.ToLower(c.Param("slug"))

	sp, err := h.SAML.Provider(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "no SAML connection named " + slug},
		})
		return
	}

	out, err := xml.MarshalIndent(sp.Metadata(), "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "could not render metadata"},
		})
		return
	}
	c.Data(http.StatusOK, "application/samlmetadata+xml", append([]byte(xml.Header), out...))
}

// SAMLBegin sends the user to their IdP with a signed authentication request.
//
//	GET /api/auth/saml/:slug
func (h *SSOHandler) SAMLBegin(c *gin.Context) {
	slug := strings.ToLower(c.Param("slug"))

	sp, err := h.SAML.Provider(slug)
	if err != nil {
		h.failLogin(c, "That sign-in method is not available.")
		return
	}

	authURL, err := sp.MakeRedirectAuthenticationRequest("")
	if err != nil {
		log.Printf("saml %s: authn request: %v", slug, err)
		h.failLogin(c, "Could not reach the identity provider.")
		return
	}

	c.Redirect(http.StatusFound, authURL.String())
}

// SAMLACS is the Assertion Consumer Service — where the IdP POSTs the signed
// assertion once the user has authenticated.
//
// ParseResponse does the security-critical work: it verifies the signature
// against the certificate in the IdP metadata, checks the audience is us, and
// enforces the assertion's validity window. Anything it rejects is treated as a
// failed login with a generic message, because the detail is only useful to an
// attacker.
//
//	POST /api/auth/saml/:slug/acs
func (h *SSOHandler) SAMLACS(c *gin.Context) {
	slug := strings.ToLower(c.Param("slug"))

	sp, err := h.SAML.Provider(slug)
	if err != nil {
		h.failLogin(c, "That sign-in method is not available.")
		return
	}

	var conn models.SSOConnection
	if err := h.DB.Where("slug = ? AND enabled = ? AND protocol = ?", slug, true, "saml").
		First(&conn).Error; err != nil {
		h.failLogin(c, "That sign-in method is not available.")
		return
	}

	if err := c.Request.ParseForm(); err != nil {
		h.failLogin(c, "Sign-in was not completed.")
		return
	}

	// No tracked request IDs: an IdP-initiated login has no request to match,
	// and crewjam skips the InResponseTo check when the connection allows that.
	// With it disallowed, an unsolicited assertion is refused here.
	assertion, err := sp.ParseResponse(c.Request, []string{})
	if err != nil {
		log.Printf("saml %s: assertion rejected: %v", slug, err)
		h.failLogin(c, "Sign-in could not be verified. Please try again.")
		return
	}

	ident := samlIdentity(&conn, assertion)
	if strings.TrimSpace(ident.Email) == "" {
		h.failLogin(c, "Your identity provider did not release an email address.")
		return
	}
	if strings.TrimSpace(ident.Subject) == "" {
		h.failLogin(c, "Your identity provider did not release an identifier.")
		return
	}

	user, err := h.resolveUser(c, &conn, ident)
	if err != nil {
		h.failLogin(c, err.Error())
		return
	}
	if err := h.applyGroupRoles(&conn, user, ident); err != nil {
		log.Printf("saml %s: role mapping for %s: %v", slug, user.ID, err)
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
		log.Printf("saml: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	services.TouchConnection(h.DB, conn.ID)

	c.Redirect(http.StatusFound, h.Config.OAuthFrontendURL+"/auth/callback")
}

// samlIdentity normalizes an assertion onto the same shape the OIDC callback
// produces, so both protocols share one provisioning path. Each lookup falls
// back to the conventional attribute names when the connection doesn't pin one.
func samlIdentity(conn *models.SSOConnection, assertion *saml.Assertion) externalIdentity {
	emailNames := conn.AttributeOr(conn.EmailAttribute,
		"email", "mail", "emailaddress", "urn:oid:0.9.2342.19200300.100.1.3",
		"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress")
	firstNames := conn.AttributeOr(conn.FirstNameAttribute,
		"firstName", "givenName", "urn:oid:2.5.4.42",
		"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname")
	lastNames := conn.AttributeOr(conn.LastNameAttribute,
		"lastName", "surname", "sn", "urn:oid:2.5.4.4",
		"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname")
	groupNames := conn.AttributeOr(conn.GroupsAttribute,
		"groups", "memberOf", "Role",
		"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups")

	return externalIdentity{
		Subject:   services.SAMLSubject(assertion, emailNames...),
		Email:     strings.ToLower(services.SAMLAttribute(assertion, emailNames...)),
		FirstName: services.SAMLAttribute(assertion, firstNames...),
		LastName:  services.SAMLAttribute(assertion, lastNames...),
		Groups:    services.SAMLAttributeValues(assertion, groupNames...),
	}
}
`
	return strings.ReplaceAll(src, "~", "`")
}
