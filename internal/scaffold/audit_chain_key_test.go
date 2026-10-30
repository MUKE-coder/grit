package scaffold

import (
	"strings"
	"testing"
)

// Review 2.7: the audit chain was unkeyed, so anybody with database write
// access could edit a row, recompute every hash after it, and verification
// would pass. The code's own comment said it defended against "a DBA covering
// tracks" and it did not.
//
// Four pieces have to be in place together, and the project compiles with any
// three: the audit package keys the hash, the model records which algorithm
// hashed each row, config derives the key, and main hands it over. Miss the
// last one and the chain is silently unkeyed, which is why they are checked in
// one test.
func TestTheAuditChainIsKeyed(t *testing.T) {
	audit := strings.ReplaceAll(apiAuditGo(), "{{MODULE}}", "example.com/app")
	for _, want := range []string{
		"func SetChainKey(key []byte) {",
		"m := hmac.New(sha256.New, chainKey)",
		`AlgoHMACSHA256 = "hmac-sha256"`,
		`AlgoSHA256 = "sha256"`,
		// Verification honours each row's own algorithm, because a project
		// that upgraded has both.
		"expected, known := computeUnder(e.HashAlgo, prevHash, canonical)",
		// And refuses a downgrade, which is the way around the key.
		"if sawKeyed && e.HashAlgo != AlgoHMACSHA256 {",
		`"crypto/hmac"`,
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("internal/audit is missing %q", want)
		}
	}

	model := strings.ReplaceAll(apiActivityLogModelGo(), "{{MODULE}}", "example.com/app")
	if !strings.Contains(model, "HashAlgo      string") {
		t.Error("the ActivityLog model has no HashAlgo column, so a project that upgraded cannot verify its older rows")
	}

	cfg := strings.ReplaceAll(apiConfigGo(), "{{MODULE}}", "example.com/app")
	for _, want := range []string{
		"AuditChainKey string",
		`AuditChainKey:      getEnv("AUDIT_CHAIN_KEY", ""),`,
		"func (c *Config) AuditChainKeyBytes() []byte {",
		// Derived from JWT_SECRET when unset, so no project has to set a
		// second secret to get a keyed chain.
		"material = c.JWTSecret",
		// Through a hash, so the chain key and the signing key are never the
		// same bytes.
		`sha256.Sum256([]byte("audit-chain-key:" + material))`,
		`"crypto/sha256"`,
	} {
		if !strings.Contains(cfg, want) {
			t.Errorf("internal/config is missing %q", want)
		}
	}

	main := strings.ReplaceAll(apiMainGo(Options{}), "{{MODULE}}", "example.com/app")
	if !strings.Contains(main, "audit.SetChainKey(cfg.AuditChainKeyBytes())") {
		t.Error("main does not key the chain, so every piece above is in place and the chain is still unkeyed")
	}
	// Before the port is bound, and so before anything can write an entry.
	if at := strings.Index(main, "audit.SetChainKey("); at >= 0 {
		if port := strings.Index(main, "// The port, before anything else."); port >= 0 && at > port {
			t.Error("the chain is keyed after the server starts coming up, so an early entry could be written unkeyed")
		}
	}
}

// An empty key is not a key. Keyed with nothing, the chain would look keyed in
// the column and be computable by anybody, which is worse than being honestly
// unkeyed.
func TestAnEmptySecretDoesNotKeyTheChain(t *testing.T) {
	cfg := strings.ReplaceAll(apiConfigGo(), "{{MODULE}}", "example.com/app")
	if !strings.Contains(cfg, "if material == \"\" {\n\t\t// No secret at all is a development run") {
		t.Error("AuditChainKeyBytes does not refuse an empty secret, so a project with no JWT_SECRET would write rows marked keyed that anybody can recompute")
	}
}
