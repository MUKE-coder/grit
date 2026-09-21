package scaffold

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A clone has only .env.example. grit env gives it a .env with every
// placeholder replaced, in the shape each variable takes, and nothing else
// changed.
func TestFillEnvCreatesEnvFromTheExample(t *testing.T) {
	root := t.TempDir()
	example := envExampleFile(Options{ProjectName: "demo"})
	if !strings.Contains(example, "=CHANGE_ME") {
		t.Fatal("the example has no placeholders to fill")
	}
	if err := os.WriteFile(filepath.Join(root, ".env.example"), []byte(example), 0o644); err != nil {
		t.Fatal(err)
	}

	created, filled, err := FillEnv(root)
	if err != nil {
		t.Fatal(err)
	}
	if !created || len(filled) == 0 {
		t.Fatalf("created %v, filled %v", created, filled)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	env := string(raw)
	if strings.Contains(env, "CHANGE_ME") {
		t.Error(".env still has a placeholder")
	}
	key := regexp.MustCompile(`(?m)^FIELD_ENCRYPTION_KEY=(\S+)$`).FindStringSubmatch(env)
	if key == nil {
		t.Fatal("no FIELD_ENCRYPTION_KEY in .env")
	}
	if b, err := base64.StdEncoding.DecodeString(key[1]); err != nil || len(b) != 32 {
		t.Errorf("FIELD_ENCRYPTION_KEY is not 32 bytes of base64 (%d bytes, %v)", len(b), err)
	}
	if jwt := regexp.MustCompile(`(?m)^JWT_SECRET=([0-9a-f]+)$`).FindStringSubmatch(env); jwt == nil || len(jwt[1]) != 64 {
		t.Error("JWT_SECRET is not 64 hex characters")
	}
	// Every other line is the example's, untouched.
	exampleLines, envLines := strings.Split(example, "\n"), strings.Split(env, "\n")
	if len(exampleLines) != len(envLines) {
		t.Fatalf("%d lines became %d", len(exampleLines), len(envLines))
	}
	for i := range exampleLines {
		if !strings.HasSuffix(exampleLines[i], "=CHANGE_ME") && exampleLines[i] != envLines[i] {
			t.Errorf("line %d changed: %q became %q", i+1, exampleLines[i], envLines[i])
		}
	}
}

// Run again, grit env fills what is still a placeholder and leaves every value
// someone has set alone.
func TestFillEnvKeepsValuesYouSet(t *testing.T) {
	root := t.TempDir()
	env := "JWT_SECRET=mine-do-not-touch\nSENTINEL_PASSWORD=CHANGE_ME\r\nAPP_ENV=development\n"
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte(env), 0o600); err != nil {
		t.Fatal(err)
	}
	created, filled, err := FillEnv(root)
	if err != nil {
		t.Fatal(err)
	}
	if created || len(filled) != 1 || filled[0] != "SENTINEL_PASSWORD" {
		t.Fatalf("created %v, filled %v; want only SENTINEL_PASSWORD", created, filled)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	if !strings.Contains(got, "JWT_SECRET=mine-do-not-touch\n") || !strings.Contains(got, "APP_ENV=development\n") {
		t.Errorf("a value that was set changed:\n%s", got)
	}
	if !regexp.MustCompile(`SENTINEL_PASSWORD=[0-9a-f]{32}\r\n`).MatchString(got) {
		t.Errorf("SENTINEL_PASSWORD was not filled, or lost its CRLF:\n%q", got)
	}
	if _, filled, _ := FillEnv(root); len(filled) != 0 {
		t.Errorf("a second run filled %v", filled)
	}
}

func TestFillEnvWithNothingToCopyFromIsAnError(t *testing.T) {
	if _, _, err := FillEnv(t.TempDir()); err == nil {
		t.Error("an empty directory was not an error")
	}
}
