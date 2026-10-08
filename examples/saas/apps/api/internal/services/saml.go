package services

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
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	dsig "github.com/russellhaering/goxmldsig"
	"gorm.io/gorm"

	"saas/apps/api/internal/crypto"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/safefetch"
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
	// An administrator types this URL, so it goes through the SSRF guard like
	// every other server-side fetch: no loopback, private or metadata addresses.
	return samlsp.FetchMetadata(ctx, safefetch.Client, *u)
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
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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
		id := assertion.Subject.NameID
		// A transient NameID is a new random value on every login, which is what
		// the format means. Used as the subject it made every sign-in a stranger
		// and linked another identity each time, so it is passed over for the
		// email, exactly as an absent NameID is.
		if v := strings.TrimSpace(id.Value); v != "" && id.Format != TransientNameID {
			return v
		}
	}
	return SAMLAttribute(assertion, emailNames...)
}

// TransientNameID is the NameID format whose value changes on every login.
const TransientNameID = "urn:oasis:names:tc:SAML:2.0:nameid-format:transient"
