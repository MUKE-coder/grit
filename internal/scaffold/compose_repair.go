package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairStorageSecrets brings a project scaffolded before v3.248.0 up to the fix
// for H10 in the contact-app review: MinIO ran on minioadmin/minioadmin, in
// development on a port open to the LAN and in production behind the public
// proxy, and the production stack handed Postgres and MinIO the whole .env.
//
// The compose files and config.go are the developer's, so each edit anchors on
// what Grit wrote. .env holds the developer's secrets and is never rewritten:
// a default MinIO password there is reported, and a production server now
// refuses to start with it.
func repairStorageSecrets(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, s := range []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(root, "docker-compose.yml"), repairDevComposeSource},
		{filepath.Join(root, "docker-compose.prod.yml"), repairProdComposeSource},
	} {
		if fileExists(s.path) {
			if err := repairTextFile(root, m, s.path, s.fn); err != nil {
				return err
			}
		}
	}
	if config := filepath.Join(opts.APIRoot(root), "internal", "config", "config.go"); fileExists(config) {
		if err := repairSourceFile(root, m, config, repairConfigMinioSecretSource); err != nil {
			return err
		}
	}
	if env, err := os.ReadFile(filepath.Join(root, ".env")); err == nil &&
		regexp.MustCompile(`(?m)^MINIO_SECRET_KEY=minioadmin\s*$`).Match(env) {
		fmt.Println("  • MINIO_SECRET_KEY in .env is still minioadmin. A production server now refuses to start with it:\n" +
			"    set MINIO_ACCESS_KEY and MINIO_SECRET_KEY (openssl rand -hex 24); MinIO takes them on its next start.")
	}
	return nil
}

const (
	devMinioCredentials = "      MINIO_ROOT_USER: minioadmin\n      MINIO_ROOT_PASSWORD: minioadmin\n"
	envMinioCredentials = "      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY:?set MINIO_ACCESS_KEY in .env}\n      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY:?set MINIO_SECRET_KEY in .env}\n"
)

func repairDevComposeSource(src string) (string, []string, []string) {
	if !strings.Contains(src, devMinioCredentials) {
		return src, nil, nil
	}
	return strings.Replace(src, devMinioCredentials, envMinioCredentials, 1),
		[]string{"MinIO takes its root credentials from .env instead of minioadmin"}, nil
}

// serviceBlock returns the bounds of a top-level compose service's block.
func serviceBlock(src, name string) (int, int, bool) {
	start := strings.Index(src, "\n  "+name+":\n")
	if start < 0 {
		return 0, 0, false
	}
	next := regexp.MustCompile(`\n  [A-Za-z0-9_-]+:\n|\n[A-Za-z]`).FindStringIndex(src[start+1:])
	end := len(src)
	if next != nil {
		end = start + 1 + next[0]
	}
	return start, end, true
}

func repairProdComposeSource(src string) (string, []string, []string) {
	out := src
	var fixed []string
	for _, name := range []string{"postgres", "minio"} {
		start, end, ok := serviceBlock(out, name)
		if !ok {
			continue
		}
		block := out[start:end]
		cleaned := strings.Replace(block, "    env_file:\n      - .env\n", "", 1)
		if cleaned != block {
			out = out[:start] + cleaned + out[end:]
			fixed = append(fixed, name+" no longer receives all of .env")
		}
	}
	replacer := strings.NewReplacer(
		"${POSTGRES_PASSWORD:-grit}", "${POSTGRES_PASSWORD:?set POSTGRES_PASSWORD in .env}",
		"${MINIO_ACCESS_KEY:-minioadmin}", "${MINIO_ACCESS_KEY:?set MINIO_ACCESS_KEY in .env}",
		"${MINIO_SECRET_KEY:-minioadmin}", "${MINIO_SECRET_KEY:?set MINIO_SECRET_KEY in .env}",
	)
	if next := replacer.Replace(out); next != out {
		out = next
		fixed = append(fixed, "the Postgres and MinIO passwords have no default, so a stack without them refuses to start")
	}
	return out, fixed, nil
}

const (
	checkSecretsPulseLine = "\t\t{\"PULSE_PASSWORD\", cfg.PulsePassword, cfg.PulseEnabled},\n"
	checkSecretsMinioLine = "\t\t// MinIO's root password, when MinIO is the store. S3, R2 and B2 keys are\n\t\t// issued by the provider and are not guessable defaults.\n\t\t{\"MINIO_SECRET_KEY\", cfg.Storage.SecretKey, cfg.StorageDriver == \"minio\"},\n"
	weakSecretsLastLine   = "\t\"change-me\": true, \"change_me\": true, \"change-me-in-prod\": true, \"sentinel-secret-change-me\": true,\n"
)

// repairConfigMinioSecretSource extends checkSecrets, added in v3.243.0, to
// MinIO's root password.
func repairConfigMinioSecretSource(src string) (string, []string, []string) {
	// The marker is the checkSecrets entry itself: "MINIO_SECRET_KEY" alone
	// also names the setting getEnv reads, so it is in every config.go.
	if !strings.Contains(src, "func checkSecrets(") || strings.Contains(src, `{"MINIO_SECRET_KEY", cfg.Storage.SecretKey`) {
		return src, nil, nil
	}
	out, ok := applyInsertions(src, []insertion{
		{anchor: checkSecretsPulseLine, text: checkSecretsMinioLine},
		{anchor: weakSecretsLastLine, text: "\t\"minioadmin\": true,\n"},
	})
	if !ok {
		return src, nil, []string{"checkSecrets is not the one Grit wrote: add MINIO_SECRET_KEY to it, or production accepts minioadmin"}
	}
	return out, []string{"outside development, a default or short MinIO secret stops the server from starting"}, nil
}

// repairTextFile is repairSourceFile for a file that is not Go.
func repairTextFile(root string, m *manifest.Manifest, path string, fn func(string) (string, []string, []string)) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	crlf := strings.Contains(string(raw), "\r\n")
	src := strings.ReplaceAll(string(raw), "\r\n", "\n")
	out, fixed, warnings := fn(src)
	shown := path
	if rel, err := filepath.Rel(root, path); err == nil {
		shown = filepath.ToSlash(rel)
	}
	for _, w := range warnings {
		fmt.Printf("  ⚠ %s: %s\n", shown, w)
	}
	if out == src {
		return nil
	}
	pristine := false
	if key, inside := manifest.Rel(root, path); inside {
		pristine = m.StatusOf(root, key) == manifest.Unchanged
	}
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if pristine {
		manifest.Refresh(path)
	}
	fmt.Printf("  ✓ %s: %s\n", shown, strings.Join(fixed, "; "))
	return nil
}
