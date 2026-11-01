package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvKey(t *testing.T) {
	tests := []struct {
		line string
		key  string
		ok   bool
	}{
		{"MAIL_MAILER=log", "MAIL_MAILER", true},
		{"# MAIL_MAILER=", "MAIL_MAILER", true},
		{"# REDIS_URL=redis://localhost:6380", "REDIS_URL", true},
		{"#   DATABASE_URL=postgres://user:pass@host/db", "", false},
		{"# No Redis, because --db sqlite scaffolds a project that needs no servers.", "", false},
		{"  SMTP_HOST=localhost", "", false},
		{"", "", false},
	}
	for _, tc := range tests {
		key, ok := envKey(tc.line)
		if key != tc.key || ok != tc.ok {
			t.Errorf("envKey(%q) = %q, %v; want %q, %v", tc.line, key, ok, tc.key, tc.ok)
		}
	}
}

func TestPatchEnvReplacesValuesAndKeepsDocumentation(t *testing.T) {
	path := writeEnv(t, `# MAIL_MAILER=
RESEND_API_KEY=re_your_api_key_here
#   DATABASE_URL=postgres://user:password@host:5432/db?sslmode=require
# REDIS_URL=redis://localhost:6380
REDIS_URL=
MAIL_FROM=noreply@example.dev
`)

	err := PatchEnv(path, map[string]string{
		"MAIL_MAILER":    "log",
		"RESEND_API_KEY": "",
		"REDIS_URL":      "redis://cache.internal:6380",
	})
	if err != nil {
		t.Fatalf("PatchEnv: %v", err)
	}

	got := readFile(t, path)
	for _, want := range []string{
		"MAIL_MAILER=log",
		"REDIS_URL=redis://cache.internal:6380",
		"RESEND_API_KEY=\n",
		"#   DATABASE_URL=postgres://user:password@host:5432/db?sslmode=require",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "# MAIL_MAILER=") {
		t.Error("commented MAIL_MAILER was not uncommented")
	}
}

func TestPatchEnvRemovesStaleNoRedisComment(t *testing.T) {
	path := writeEnv(t, `# No Redis, because --db sqlite scaffolds a project that needs no servers.
# Delete this line and start one for caching, background jobs and cron.
REDIS_URL=
`)

	if err := PatchEnv(path, map[string]string{"REDIS_URL": "redis://localhost:6380"}); err != nil {
		t.Fatalf("PatchEnv: %v", err)
	}
	got := readFile(t, path)
	if strings.Contains(got, "No Redis, because") || strings.Contains(got, "Delete this line") {
		t.Errorf("stale comment survived:\n%s", got)
	}
	if !strings.Contains(got, "REDIS_URL=redis://localhost:6380") {
		t.Errorf("Redis URL was not written:\n%s", got)
	}
}

func TestPatchEnvAppendsUnknownKey(t *testing.T) {
	path := writeEnv(t, "DB_PROVIDER=sqlite\n")
	if err := PatchEnv(path, map[string]string{"STORAGE_DRIVER": "local"}); err != nil {
		t.Fatalf("PatchEnv: %v", err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "STORAGE_DRIVER=local") {
		t.Errorf("new key was not appended:\n%s", got)
	}
}

func TestEnvUpdatesClearsResendPlaceholderForLocalMail(t *testing.T) {
	opts := Options{Mail: MailConfig{Mode: MailMailhog, From: "noreply@example.dev"}}
	updates := opts.EnvUpdates()
	if updates["MAIL_MAILER"] != "" {
		t.Errorf("MAIL_MAILER = %q, want empty", updates["MAIL_MAILER"])
	}
	if updates["RESEND_API_KEY"] != "" {
		t.Errorf("RESEND_API_KEY = %q, want empty", updates["RESEND_API_KEY"])
	}
	if updates["MAIL_FROM"] != "noreply@example.dev" {
		t.Errorf("MAIL_FROM = %q", updates["MAIL_FROM"])
	}
}

func TestReadEnv(t *testing.T) {
	path := writeEnv(t, "MAIL_MAILER=log\n# REDIS_URL=redis://localhost:6380\n")
	if got := ReadEnv(path, "MAIL_MAILER"); got != "log" {
		t.Errorf("ReadEnv(MAIL_MAILER) = %q, want log", got)
	}
	if got := ReadEnv(path, "REDIS_URL"); got != "redis://localhost:6380" {
		t.Errorf("ReadEnv(REDIS_URL) = %q", got)
	}
}

func writeEnv(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
