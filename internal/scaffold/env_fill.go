package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// A project cloned from git has .env.example, with CHANGE_ME where grit new
// wrote each secret, and no .env. The API refuses to start on a placeholder, as
// it should, which left every teammate who cloned a project generating ten
// secrets by hand before it would run. FillEnv does what grit new did.

// placeholderLine is a variable still set to the placeholder.
var placeholderLine = regexp.MustCompile(`(?m)^([A-Z][A-Z0-9_]*)=CHANGE_ME(\r?)$`)

// secretBytes is how many random bytes each secret gets, matching envFile.
// Anything else gets 16 bytes (32 hex characters).
var secretBytes = map[string]int{
	"JWT_SECRET":          32,
	"SENTINEL_SECRET_KEY": 32,
	"SENTINEL_AUDIT_KEY":  32,
	"POSTGRES_PASSWORD":   24,
	"MINIO_SECRET_KEY":    24,
	"REDIS_PASSWORD":      24,
}

// freshSecret is a new value in the shape the variable takes.
func freshSecret(name string) string {
	if name == "FIELD_ENCRYPTION_KEY" {
		return randomBase64Key()
	}
	n, ok := secretBytes[name]
	if !ok {
		n = 16
	}
	return randomHex(n)
}

// fillPlaceholders replaces every CHANGE_ME value in src with a fresh secret
// and returns the names it filled, in order.
func fillPlaceholders(src string) (string, []string) {
	var names []string
	out := placeholderLine.ReplaceAllStringFunc(src, func(line string) string {
		m := placeholderLine.FindStringSubmatch(line)
		names = append(names, m[1])
		return m[1] + "=" + freshSecret(m[1]) + m[2]
	})
	return out, names
}

// FillEnv makes root/.env usable: created from .env.example when there is no
// .env, and with every CHANGE_ME in it replaced by a fresh secret. It returns
// whether it created the file and the names it filled. Values are never
// returned or printed.
func FillEnv(root string) (created bool, filled []string, err error) {
	envPath := filepath.Join(root, ".env")
	raw, err := os.ReadFile(envPath)
	switch {
	case os.IsNotExist(err):
		raw, err = os.ReadFile(filepath.Join(root, ".env.example"))
		if os.IsNotExist(err) {
			return false, nil, fmt.Errorf("neither .env nor .env.example exists in %s", root)
		}
		if err != nil {
			return false, nil, fmt.Errorf("reading .env.example: %w", err)
		}
		created = true
	case err != nil:
		return false, nil, fmt.Errorf("reading .env: %w", err)
	}

	out, names := fillPlaceholders(string(raw))
	if !created && len(names) == 0 {
		return false, nil, nil
	}
	// 0600: this file now holds real secrets.
	if err := os.WriteFile(envPath, []byte(out), 0o600); err != nil {
		return false, nil, fmt.Errorf("writing .env: %w", err)
	}
	return created, names, nil
}

// EnvFillSummary is the line grit env prints for what it filled.
func EnvFillSummary(filled []string) string {
	if len(filled) == 0 {
		return "no placeholders to fill"
	}
	return fmt.Sprintf("%d %s generated: %s", len(filled), plural(len(filled), "secret", "secrets"), strings.Join(filled, ", "))
}
