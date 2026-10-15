package scaffold

// Asymmetric token signing, and the key set that lets other services verify
// without holding anything that can mint.
//
// A shared secret is the right default for one service: it is one environment
// variable and there is nothing to distribute. It stops being right the moment
// a second service has to verify the tokens, because HMAC verification needs
// the same key that signs. Every service that can check a token can also forge
// one, so the blast radius of any one of them being compromised is the whole
// estate, and an attacker who reads the secret out of the least important
// service can mint an administrator token for all of them.
//
// Signing with a private key and verifying with the public half removes that:
// a verifier holds something that cannot sign. The cost is key management,
// which is why this is opt-in rather than the default. Set JWT_PRIVATE_KEY and
// the API signs RS256 and publishes its public half at
// /.well-known/jwks.json; leave it unset and nothing changes.

// apiJWTKeysGo emits internal/services/jwtkeys.go.
func apiJWTKeysGo() string {
	return `package services

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ErrNoVerificationKey is returned when a token names a key this process does
// not have. With rotation configured that means the key is older than the one
// before the current one, or the token came from somewhere else entirely.
var ErrNoVerificationKey = errors.New("no verification key for this token")

// SigningKeys is what signs tokens and what verifies them, which are the same
// thing on HS256 and deliberately not on RS256.
//
// The zero value is not usable: build one with LoadSigningKeys.
type SigningKeys struct {
	// Algorithm is "HS256" or "RS256".
	Algorithm string

	// secret signs and verifies on HS256.
	secret []byte

	// private signs on RS256. nil on HS256, and nil in a service that only
	// verifies, which is the entire point of the arrangement.
	private *rsa.PrivateKey

	// public verifies on RS256, keyed by the kid in the token header. Holds
	// the current key and, while a rotation is in flight, the previous one, so
	// tokens signed before the switch keep working until they expire.
	public map[string]*rsa.PublicKey

	// currentKID names the key new tokens are signed with.
	currentKID string

	// Issuer and Audience are written into every token this process mints.
	// Empty means the claim is not written, which is the old behaviour.
	Issuer   string
	Audience string
}

// LoadSigningKeys decides how this process signs, from configuration.
//
// With no private key it is HS256 on the shared secret, which is what every
// project gets until it asks for something else. With one it is RS256, and the
// public half is derived from the private key rather than configured
// separately, because two values that must correspond are two values that can
// disagree.
func LoadSigningKeys(secret, privatePEM, previousPublicPEM, issuer, audience string) (*SigningKeys, error) {
	k := &SigningKeys{Issuer: issuer, Audience: audience}

	privatePEM, err := expandKeyMaterial("JWT_PRIVATE_KEY", privatePEM)
	if err != nil {
		return nil, err
	}
	privatePEM = strings.TrimSpace(privatePEM)
	if privatePEM == "" {
		if secret == "" {
			return nil, errors.New("JWT_SECRET is required when JWT_PRIVATE_KEY is not set")
		}
		k.Algorithm = "HS256"
		k.secret = []byte(secret)
		return k, nil
	}

	private, parseErr := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePEM))
	if parseErr != nil {
		return nil, fmt.Errorf("parsing JWT_PRIVATE_KEY: %w", parseErr)
	}
	if size := private.N.BitLen(); size < 2048 {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY is %d bits; 2048 is the minimum", size)
	}

	k.Algorithm = "RS256"
	k.private = private
	k.public = map[string]*rsa.PublicKey{}
	k.currentKID = thumbprint(&private.PublicKey)
	k.public[k.currentKID] = &private.PublicKey

	// The key being rotated away from. Tokens it signed are still valid until
	// they expire, so without this a rotation signs out everybody holding one.
	previous, err := expandKeyMaterial("JWT_PREVIOUS_PUBLIC_KEY", previousPublicPEM)
	if err != nil {
		return nil, err
	}
	if previous = strings.TrimSpace(previous); previous != "" {
		pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(previous))
		if err != nil {
			return nil, fmt.Errorf("parsing JWT_PREVIOUS_PUBLIC_KEY: %w", err)
		}
		k.public[thumbprint(pub)] = pub
	}
	return k, nil
}

// expandKeyMaterial lets a key be given inline or as a path to a file.
//
// A PEM block is multi-line, which environment variables and process managers
// handle inconsistently: \n ends up literal, or the value is truncated at the
// first newline and the error is a parse failure nobody can explain. So a value
// holding a PEM block is used as it is, with any literal \n turned back into a
// newline, and anything else is treated as a path and read.
//
// An unreadable path is reported rather than passed on. Handing a file path to
// a PEM parser produces "invalid key: must be PEM encoded", which says nothing
// about the actual mistake, and the actual mistake is nearly always a relative
// path resolved against a working directory the operator did not expect: the
// server runs from the API directory, so a path that works for a CLI command
// run from the project root does not work here.
func expandKeyMaterial(name, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	if strings.Contains(value, "-----BEGIN") {
		return strings.ReplaceAll(value, "\\n", "\n"), nil
	}
	data, err := os.ReadFile(value)
	if err != nil {
		cwd, _ := os.Getwd()
		return "", fmt.Errorf("%s=%q is neither a PEM block nor a readable file (looked from %s): %w", name, value, cwd, err)
	}
	return string(data), nil
}

// Method is the signing method new tokens use.
func (k *SigningKeys) Method() jwt.SigningMethod {
	if k.Algorithm == "RS256" {
		return jwt.SigningMethodRS256
	}
	return jwt.SigningMethodHS256
}

// ValidMethods is what a parser will accept, and it is exactly one entry.
//
// Pinning this is what closes the algorithm confusion attack: given a parser
// that accepts both, an attacker takes a public key that is published on
// purpose, signs a token with it as an HMAC secret, and a verifier that trusts
// the header's alg accepts it.
func (k *SigningKeys) ValidMethods() []string {
	return []string{k.Method().Alg()}
}

// SignKey is what to sign with: the secret on HS256, the private key on RS256.
func (k *SigningKeys) SignKey() interface{} {
	if k.Algorithm == "RS256" {
		return k.private
	}
	return k.secret
}

// CurrentKID names the signing key, for the token header. Empty on HS256,
// where there is only ever one key and nothing to choose between.
func (k *SigningKeys) CurrentKID() string { return k.currentKID }

// VerifyKey returns the key a token signed under kid should be checked with.
func (k *SigningKeys) VerifyKey(kid string) (interface{}, error) {
	if k.Algorithm != "RS256" {
		return k.secret, nil
	}
	// No kid: a token from before kids were written, or from another issuer.
	// Try the current key, which is right for the first and wrong for the
	// second in a way the signature check will catch.
	if kid == "" {
		return k.public[k.currentKID], nil
	}
	pub, ok := k.public[kid]
	if !ok {
		return nil, fmt.Errorf("%w: kid %q", ErrNoVerificationKey, kid)
	}
	return pub, nil
}

// CanMint reports whether this process holds something that can sign. A
// verify-only deployment answers false, which is the property asymmetric
// signing is bought for.
func (k *SigningKeys) CanMint() bool {
	return k.Algorithm != "RS256" || k.private != nil
}

// JWK is one key in the published set, in the shape RFC 7517 defines.
type JWK struct {
	Kty string ` + "`" + `json:"kty"` + "`" + `
	Use string ` + "`" + `json:"use"` + "`" + `
	Alg string ` + "`" + `json:"alg"` + "`" + `
	Kid string ` + "`" + `json:"kid"` + "`" + `
	N   string ` + "`" + `json:"n"` + "`" + `
	E   string ` + "`" + `json:"e"` + "`" + `
}

// JWKS is the document served at /.well-known/jwks.json.
type JWKS struct {
	Keys []JWK ` + "`" + `json:"keys"` + "`" + `
}

// PublicJWKS is every key a verifier should accept right now: the current one
// and, during a rotation, the previous one.
//
// On HS256 it is empty, and deliberately so. The verification key there is the
// signing secret, and publishing it would hand out the ability to mint.
func (k *SigningKeys) PublicJWKS() JWKS {
	set := JWKS{Keys: []JWK{}}
	if k.Algorithm != "RS256" {
		return set
	}
	// Current first, so a client that takes the first match uses the newest.
	if pub, ok := k.public[k.currentKID]; ok {
		set.Keys = append(set.Keys, jwkFor(k.currentKID, pub))
	}
	for kid, pub := range k.public {
		if kid != k.currentKID {
			set.Keys = append(set.Keys, jwkFor(kid, pub))
		}
	}
	return set
}

func jwkFor(kid string, pub *rsa.PublicKey) JWK {
	return JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: kid,
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(exponentBytes(pub.E)),
	}
}

// exponentBytes is an RSA public exponent as the minimal big-endian bytes a
// JWK's "e" member carries.
//
// The guard is why this is a function rather than two inline conversions. E is
// an int, and uint64(E) on a negative int is a silent sign change into an
// enormous number: the "e" published would then not be the exponent the key
// signs with, and nothing would say so. No key x509 will parse has a
// non-positive exponent, so this is unreachable rather than defensive, and an
// unchecked sign conversion is still worth not writing. A zero byte is what
// comes back, which every client rejects as an invalid key, which is the right
// outcome for a key that should not have got this far.
func exponentBytes(e int) []byte {
	if e <= 0 {
		return []byte{0}
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(e))
	return trimLeadingZeros(b)
}

func trimLeadingZeros(b []byte) []byte {
	for len(b) > 1 && b[0] == 0 {
		b = b[1:]
	}
	return b
}

// thumbprint derives a stable key id from the key itself, per RFC 7638.
//
// Derived rather than configured, so the kid in a token and the kid in the
// published set cannot disagree: both are computed from the same key. It also
// means rotating a key changes its id without anybody having to remember to.
func thumbprint(pub *rsa.PublicKey) string {
	canonical := fmt.Sprintf(` + "`" + `{"e":"%s","kty":"RSA","n":"%s"}` + "`" + `,
		base64.RawURLEncoding.EncodeToString(exponentBytes(pub.E)),
		base64.RawURLEncoding.EncodeToString(pub.N.Bytes()))
	sum := sha256.Sum256([]byte(canonical))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
`
}

// apiJWTKeysTestGo emits internal/services/jwtkeys_test.go.
func apiJWTKeysTestGo() string {
	return `package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
)

func rsaPEM(t *testing.T, bits int) (string, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		t.Fatalf("generating a key: %v", err)
	}
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	return string(pem.EncodeToMemory(block)), key
}

func publicPEM(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("marshalling the public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

func TestDefaultsToSharedSecret(t *testing.T) {
	keys, err := LoadSigningKeys("a-secret-at-least-32-characters-long", "", "", "", "")
	if err != nil {
		t.Fatalf("loading: %v", err)
	}
	if keys.Algorithm != "HS256" {
		t.Fatalf("Algorithm = %q, want HS256 when no private key is set", keys.Algorithm)
	}
	// The whole point of the default: nothing is published, because the
	// verification key is the signing key.
	if got := len(keys.PublicJWKS().Keys); got != 0 {
		t.Fatalf("published %d keys on HS256; a shared secret must never be published", got)
	}
}

func TestSecretIsRequiredWithoutAPrivateKey(t *testing.T) {
	if _, err := LoadSigningKeys("", "", "", "", ""); err == nil {
		t.Fatal("loading with neither a secret nor a private key should fail rather than sign with nothing")
	}
}

func TestPrivateKeySwitchesToRS256AndPublishesOnlyThePublicHalf(t *testing.T) {
	private, key := rsaPEM(t, 2048)
	keys, err := LoadSigningKeys("", private, "", "https://auth.example.com", "example-api")
	if err != nil {
		t.Fatalf("loading: %v", err)
	}
	if keys.Algorithm != "RS256" {
		t.Fatalf("Algorithm = %q, want RS256", keys.Algorithm)
	}
	if keys.ValidMethods()[0] != "RS256" || len(keys.ValidMethods()) != 1 {
		t.Fatalf("ValidMethods = %v, want exactly [RS256]", keys.ValidMethods())
	}

	set := keys.PublicJWKS()
	if len(set.Keys) != 1 {
		t.Fatalf("published %d keys, want 1", len(set.Keys))
	}
	if set.Keys[0].Kid != keys.CurrentKID() {
		t.Fatal("the published kid must be the one tokens are signed with, or no verifier can match them")
	}
	if set.Keys[0].Kty != "RSA" || set.Keys[0].Use != "sig" || set.Keys[0].Alg != "RS256" {
		t.Fatalf("published key is not a signing RSA key: %+v", set.Keys[0])
	}
	// Nothing in the document may be derivable into a signing key.
	if set.Keys[0].N == "" || set.Keys[0].E == "" {
		t.Fatal("a published key needs its modulus and exponent to be usable")
	}
	_ = key
}

func TestShortKeysAreRefused(t *testing.T) {
	private, _ := rsaPEM(t, 1024)
	if _, err := LoadSigningKeys("", private, "", "", ""); err == nil {
		t.Fatal("a 1024 bit key should be refused rather than quietly accepted")
	}
}

// The property rotation exists for: tokens signed by the key being retired keep
// working until they expire, instead of everybody being signed out at once.
func TestRotationKeepsThePreviousKeyVerifiable(t *testing.T) {
	oldPrivate, oldKey := rsaPEM(t, 2048)
	newPrivate, _ := rsaPEM(t, 2048)

	before, err := LoadSigningKeys("", oldPrivate, "", "", "")
	if err != nil {
		t.Fatalf("loading the old key: %v", err)
	}
	oldKID := before.CurrentKID()

	after, err := LoadSigningKeys("", newPrivate, publicPEM(t, oldKey), "", "")
	if err != nil {
		t.Fatalf("loading after rotation: %v", err)
	}
	if after.CurrentKID() == oldKID {
		t.Fatal("a new key must get a new id, or verifiers cannot tell them apart")
	}
	if _, err := after.VerifyKey(oldKID); err != nil {
		t.Fatalf("the previous key must still verify during a rotation: %v", err)
	}
	if _, err := after.VerifyKey(after.CurrentKID()); err != nil {
		t.Fatalf("the current key must verify: %v", err)
	}
	if len(after.PublicJWKS().Keys) != 2 {
		t.Fatal("both keys must be published while a rotation is in flight")
	}
}

func TestUnknownKeyIsRefused(t *testing.T) {
	private, _ := rsaPEM(t, 2048)
	keys, err := LoadSigningKeys("", private, "", "", "")
	if err != nil {
		t.Fatalf("loading: %v", err)
	}
	if _, err := keys.VerifyKey("a-kid-from-somewhere-else"); !errors.Is(err, ErrNoVerificationKey) {
		t.Fatalf("err = %v, want ErrNoVerificationKey", err)
	}
}

// A path that does not resolve is the commonest mistake here, because the
// server runs from the API directory and a path that works from the project
// root does not. The error has to name the value and where it looked, or it
// surfaces as "invalid key: must be PEM encoded" and says nothing useful.
func TestUnreadableKeyPathSaysSo(t *testing.T) {
	_, err := LoadSigningKeys("", "config/keys/jwt-private.pem", "", "", "")
	if err == nil {
		t.Fatal("a path that cannot be read should fail")
	}
	for _, want := range []string{"JWT_PRIVATE_KEY", "config/keys/jwt-private.pem", "readable file"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "PEM encoded") {
		t.Errorf("a path mistake should not be reported as a PEM parse failure: %v", err)
	}
}

func TestInlineKeyWithEscapedNewlines(t *testing.T) {
	private, _ := rsaPEM(t, 2048)
	// How a PEM block survives most process managers.
	inline := ""
	for _, r := range private {
		if r == '\n' {
			inline += "\\n"
			continue
		}
		inline += string(r)
	}
	if _, err := LoadSigningKeys("", inline, "", "", ""); err != nil {
		t.Fatalf("a key with escaped newlines should load: %v", err)
	}
}
`
}
