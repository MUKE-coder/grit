package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Review 2.7: keying the audit chain in a project that already exists.
//
// writeAuditChainFiles already rewrites the model and the audit package, so an
// upgraded project gets the hash_algo column and the keyed code path from
// those. config.go and main.go are the developer's files, so they are repaired
// rather than replaced, and without them nothing calls SetChainKey: the
// project builds, its tests pass, and the chain stays unkeyed.
//
// That is a safe way to fail and a useless way to finish. An unkeyed chain is
// one anybody with database write access can recompute, which is the whole
// finding.

const (
	auditKeyConfigField = `
	// AuditChainKey keys the activity log's hash chain.
	//
	// Empty derives it from JWTSecret, so an existing project keeps working
	// without setting anything. Set it explicitly when the audit chain should
	// outlive a JWT_SECRET rotation, or when the two should be held by
	// different people.
	//
	// Without a key the chain is a plain SHA-256 over public inputs, which
	// anybody with database write access can recompute: the tampering then
	// verifies, which is the opposite of what the chain is for (review 2.7).
	AuditChainKey string
`

	auditKeyConfigRead = "\t\tAuditChainKey:      getEnv(\"AUDIT_CHAIN_KEY\", \"\"),\n"

	auditKeyConfigMethod = `
// AuditChainKeyBytes is the key the activity log's hash chain is computed
// under. main hands it to audit.SetChainKey at startup.
//
// AUDIT_CHAIN_KEY when set. Otherwise derived from JWT_SECRET, which every
// project already has and which checkSecrets already refuses to leave at a
// default outside development: the chain key is exactly as strong as that
// secret, and no project has to set a second one to get a keyed chain.
//
// Through SHA-256 rather than used directly, so the chain key and the token
// signing key are never the same bytes.
func (c *Config) AuditChainKeyBytes() []byte {
	material := c.AuditChainKey
	if material == "" {
		material = c.JWTSecret
	}
	if material == "" {
		// No secret at all is a development run that has not been configured.
		// An unkeyed chain is better than a chain keyed with nothing, because
		// an empty key would look keyed and protect nothing.
		return nil
	}
	sum := sha256.Sum256([]byte("audit-chain-key:" + material))
	return sum[:]
}
`

	auditKeyMainCall = `	// Key the activity log's hash chain, before anything can write to it.
	//
	// Unkeyed, the chain was a plain SHA-256 over inputs anybody can read off
	// this source, so a row edited in the database could be re-hashed and the
	// chain would still verify: the thing it exists to detect (review 2.7).
	//
	// The key is AUDIT_CHAIN_KEY, or derived from JWT_SECRET when that is
	// unset, so an existing project gets a keyed chain without setting
	// anything.
	audit.SetChainKey(cfg.AuditChainKeyBytes())

`
)

// repairAuditChainKey gives an existing project the key its audit package is
// now written to use.
func repairAuditChainKey(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()

	cfg := filepath.Join(apiRoot, "internal", "config", "config.go")
	if fileExists(cfg) {
		if err := repairSourceFile(root, m, cfg, repairAuditKeyConfigSource); err != nil {
			return err
		}
	}

	// After config, and only if config now has the method: a main that calls
	// a method which is not there is a project that does not compile, and an
	// upgrade must not be the reason for that.
	if !fileContains(cfg, "func (c *Config) AuditChainKeyBytes()") {
		return nil
	}
	main := filepath.Join(apiRoot, "cmd", "server", "main.go")
	if !fileExists(main) {
		return nil
	}
	return repairSourceFile(root, m, main, func(src string) (string, []string, []string) {
		return repairAuditKeyMainSource(src, module)
	})
}

func repairAuditKeyConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "AuditChainKeyBytes()") {
		return src, nil, nil
	}
	const warn = "config.go is not the file Grit wrote, so the audit chain stays unkeyed: add an AuditChainKey field read from AUDIT_CHAIN_KEY and an AuditChainKeyBytes method deriving from JWTSecret, or anybody with database access can recompute the chain"

	// The field, beside the other key this project holds.
	const fieldAnchor = "\tFieldEncryptionKey string\n"
	if strings.Count(src, fieldAnchor) != 1 {
		return src, nil, []string{warn}
	}
	out := strings.Replace(src, fieldAnchor, fieldAnchor+auditKeyConfigField, 1)

	// The read, beside the same key's.
	const readAnchor = "\t\tFieldEncryptionKey: getEnv(\"FIELD_ENCRYPTION_KEY\", \"\"),\n"
	if strings.Count(out, readAnchor) != 1 {
		return src, nil, []string{warn}
	}
	out = strings.Replace(out, readAnchor, readAnchor+auditKeyConfigRead, 1)

	// The method, at the end, where a reader looking for it will find it
	// beside the rest of Config's own helpers.
	out = strings.TrimRight(out, "\n") + "\n" + auditKeyConfigMethod

	withImport, ok := addImportGroup(out, "crypto/sha256")
	if !ok {
		return src, nil, []string{warn}
	}
	return withImport, []string{"the audit chain is keyed from AUDIT_CHAIN_KEY, or from JWT_SECRET when that is unset, so a row edited in the database can no longer be re-hashed to match"}, nil
}

func repairAuditKeyMainSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "audit.SetChainKey(") {
		return src, nil, nil
	}
	// Only a main that loads a config has one to take the key from.
	const anchor = "\tcfg, err := config.Load()\n\tif err != nil {\n\t\tlog.Fatalf(\"Failed to load config: %v\", err)\n\t}\n\n"
	if strings.Count(src, anchor) != 1 {
		return src, nil, []string{"main.go does not load the config the way Grit wrote it, so the audit chain stays unkeyed: call audit.SetChainKey(cfg.AuditChainKeyBytes()) after the config loads"}
	}
	out := strings.Replace(src, anchor, anchor+auditKeyMainCall, 1)

	if !strings.Contains(out, module+"/internal/audit\"") {
		var ok bool
		out, ok = addImportGroup(out, module+"/internal/audit")
		if !ok {
			return src, nil, []string{"could not import internal/audit into main.go, so the audit chain stays unkeyed"}
		}
	}
	return out, []string{"the audit chain is keyed at startup, before anything can write an entry"}, nil
}
