package crypto

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
