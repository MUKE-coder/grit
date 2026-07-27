package scaffold

// Field-level encryption — transparent AES-256-GCM on individual model fields.
//
// Some columns hold data you must store and display but should never be readable
// straight out of the database: personal notes, tokens, contact details. This
// gives you a string type, EncryptedString, that GORM stores as ciphertext and
// hands back as plaintext — the column is opaque without the key, and nothing in
// the handler or the model has to think about it.
//
// The scheme is versioned (enc:v1: = AES-256-GCM, fresh 12-byte nonce per write)
// so it can rotate later. Because each write uses a new nonce the ciphertext is
// non-deterministic, so an encrypted column can't be queried by equality — it's
// for data you keep and show, not for keys or lookup columns. With no key
// configured the type passes values through as plaintext, so a project can adopt
// encryption later without a migration.

func apiCryptoFieldGo() string {
	src := `// Package crypto provides transparent field-level encryption for model columns.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

// cipherPrefix tags an encrypted value and versions the scheme, so the algorithm
// or key can change later without ambiguity about how an existing value was
// written. v1 = AES-256-GCM with a random 12-byte nonce.
const cipherPrefix = "enc:v1:"

var (
	keyMu    sync.RWMutex
	fieldKey []byte // nil = encryption disabled; values pass through as plaintext
)

// InitFieldKey configures the process-wide field key from a base64 string that
// decodes to 32 bytes (AES-256). An empty string disables encryption so a
// project can adopt the feature later without migrating existing rows. A
// non-empty key of the wrong length is a hard error: silently running without
// the encryption you asked for is worse than refusing to start.
func InitFieldKey(b64 string) error {
	keyMu.Lock()
	defer keyMu.Unlock()
	if b64 == "" {
		fieldKey = nil
		return nil
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return fmt.Errorf("FIELD_ENCRYPTION_KEY must be base64: %w", err)
	}
	if len(raw) != 32 {
		return fmt.Errorf("FIELD_ENCRYPTION_KEY must decode to 32 bytes for AES-256 (got %d)", len(raw))
	}
	fieldKey = raw
	return nil
}

// EncryptionEnabled reports whether a key is configured.
func EncryptionEnabled() bool {
	keyMu.RLock()
	defer keyMu.RUnlock()
	return fieldKey != nil
}

func aead() (cipher.AEAD, bool, error) {
	keyMu.RLock()
	k := fieldKey
	keyMu.RUnlock()
	if k == nil {
		return nil, false, nil
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, false, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, false, err
	}
	return gcm, true, nil
}

// Encrypt returns the ciphertext form of s. Empty or already-encrypted values,
// and the disabled-key case, return s unchanged so the operation is idempotent
// and safe to call unconditionally.
func Encrypt(s string) (string, error) {
	if s == "" || strings.HasPrefix(s, cipherPrefix) {
		return s, nil
	}
	gcm, on, err := aead()
	if err != nil {
		return "", err
	}
	if !on {
		return s, nil
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(s), nil)
	return cipherPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt reverses Encrypt. A value without the prefix is returned unchanged —
// that is plaintext written before a key was configured, which stays readable.
// A prefixed value with no key configured is an error, not a silent leak.
func Decrypt(s string) (string, error) {
	if !strings.HasPrefix(s, cipherPrefix) {
		return s, nil
	}
	gcm, on, err := aead()
	if err != nil {
		return "", err
	}
	if !on {
		return "", errors.New("value is encrypted but FIELD_ENCRYPTION_KEY is not configured")
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, cipherPrefix))
	if err != nil {
		return "", fmt.Errorf("decoding ciphertext: %w", err)
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("ciphertext shorter than nonce")
	}
	nonce, sealed := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed (wrong key or tampered value): %w", err)
	}
	return string(plain), nil
}

// EncryptedString is a string that is transparently encrypted at rest. Declare a
// model field as crypto.EncryptedString and GORM stores ciphertext while your
// code sees plaintext. The scheme is non-deterministic, so encrypted columns
// cannot be queried by equality — use it for data you store and display but never
// filter on (notes, tokens, personal details), not for keys or lookup columns.
type EncryptedString string

// GormDataType stores the column as text: ciphertext is base64 and longer than
// the plaintext, so a bounded varchar could truncate it.
func (EncryptedString) GormDataType() string { return "text" }

// Value implements driver.Valuer, so GORM writes ciphertext. This also runs for
// map-based Updates — provided the map value is an EncryptedString and not a bare
// string, since a bare string never reaches this method.
func (e EncryptedString) Value() (driver.Value, error) {
	return Encrypt(string(e))
}

// Scan implements sql.Scanner, decrypting the column on read.
func (e *EncryptedString) Scan(src interface{}) error {
	if src == nil {
		*e = ""
		return nil
	}
	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("EncryptedString: cannot scan %T", src)
	}
	dec, err := Decrypt(raw)
	if err != nil {
		return err
	}
	*e = EncryptedString(dec)
	return nil
}

// MarshalJSON / UnmarshalJSON make the type behave as a plain string over the
// wire, so API responses carry plaintext and request binding is unchanged.
func (e EncryptedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

func (e *EncryptedString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*e = EncryptedString(s)
	return nil
}
`
	return src
}

func apiCryptoFieldTestGo() string {
	src := `package crypto

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

const testKey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" // 32 zero bytes, base64

func TestRoundTrip(t *testing.T) {
	if err := InitFieldKey(testKey); err != nil {
		t.Fatalf("init: %v", err)
	}
	defer InitFieldKey("")

	ct, err := Encrypt("hunter2")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if !strings.HasPrefix(ct, "enc:v1:") {
		t.Errorf("ciphertext missing version prefix: %q", ct)
	}
	if strings.Contains(ct, "hunter2") {
		t.Errorf("plaintext leaked into ciphertext: %q", ct)
	}
	pt, err := Decrypt(ct)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if pt != "hunter2" {
		t.Errorf("round trip = %q, want hunter2", pt)
	}
}

func TestNonceIsRandom(t *testing.T) {
	InitFieldKey(testKey)
	defer InitFieldKey("")
	a, _ := Encrypt("same")
	b, _ := Encrypt("same")
	if a == b {
		t.Errorf("identical plaintext produced identical ciphertext — nonce not random")
	}
}

func TestWrongKeyFails(t *testing.T) {
	InitFieldKey(testKey)
	ct, _ := Encrypt("secret")
	// Rotate to a different key and try to read the old value.
	other := base64.StdEncoding.EncodeToString([]byte("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"))
	InitFieldKey(other)
	defer InitFieldKey("")
	if _, err := Decrypt(ct); err == nil {
		t.Errorf("decrypt with the wrong key should fail")
	}
}

func TestDisabledPassthrough(t *testing.T) {
	InitFieldKey("")
	ct, err := Encrypt("plain")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if ct != "plain" {
		t.Errorf("with no key, value should pass through: got %q", ct)
	}
	// A value written earlier without the prefix must still read back.
	pt, _ := Decrypt("plain")
	if pt != "plain" {
		t.Errorf("plaintext passthrough on read = %q", pt)
	}
}

func TestBadKeyRejected(t *testing.T) {
	if err := InitFieldKey("not-base64!!!"); err == nil {
		t.Errorf("non-base64 key should be rejected")
	}
	if err := InitFieldKey(base64.StdEncoding.EncodeToString([]byte("tooshort"))); err == nil {
		t.Errorf("wrong-length key should be rejected")
	}
	InitFieldKey("")
}

func TestJSONIsPlaintext(t *testing.T) {
	InitFieldKey(testKey)
	defer InitFieldKey("")
	e := EncryptedString("visible")
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != "\"visible\"" {
		t.Errorf("JSON = %s, want \"visible\"", b)
	}
	var back EncryptedString
	if err := json.Unmarshal([]byte("\"in\""), &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back != "in" {
		t.Errorf("unmarshal = %q", back)
	}
}

func TestValueProducesCiphertext(t *testing.T) {
	InitFieldKey(testKey)
	defer InitFieldKey("")
	v, err := EncryptedString("db-bound").Value()
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	s, _ := v.(string)
	if !strings.HasPrefix(s, "enc:v1:") {
		t.Errorf("Value() must yield ciphertext for the DB, got %q", s)
	}
}
`
	return src
}
