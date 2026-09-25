package deploy

import (
	"os"
	"path/filepath"
	"testing"
)

// What goes to Railway, and what does not.
//
// A generated .env has 107 entries. Most are empty placeholders for providers
// the project does not use, and a good number configure the local compose
// stack. Sending all of them is minutes of CLI calls, a dashboard nobody can
// read, and a deployed app carrying settings that point at a laptop.
func TestRailwayVariablesSendsWhatADeployNeeds(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	body := `# a comment, and a blank line follow

APP_ENV=development
APP_NAME=Shop
JWT_SECRET=s3cret

# Empty placeholders: providers this project does not use.
STRIPE_SECRET=
AWS_SES_REGION=

# Railway sets these itself.
PORT=8080
APP_PORT=8080

# The compose stack, and nothing else.
DB_HOST=db
DB_PORT=5432
MINIO_ROOT_USER=minio
MAILHOG_UI_PORT=8025
POSTGRES_PASSWORD=postgres
MYSQL_HOST=mysql

# Quoted values keep their quotes off.
MAIL_FROM="hello@shop.test"
`
	if err := os.WriteFile(env, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	got := func(provision bool) map[string]string {
		t.Helper()
		vars, err := railwayVariables(RailwayConfig{EnvFile: env, Provision: provision})
		if err != nil {
			t.Fatalf("railwayVariables: %v", err)
		}
		m := map[string]string{}
		for _, v := range vars {
			m[v.Key] = v.Value
		}
		return m
	}

	plain := got(false)

	for _, key := range []string{
		"STRIPE_SECRET", "AWS_SES_REGION", // empty: configure nothing
		"PORT", "APP_PORT", // Railway's own
		"DB_HOST", "DB_PORT", // the compose database
		"MINIO_ROOT_USER", "MAILHOG_UI_PORT", // containers on a laptop
		"POSTGRES_PASSWORD", "MYSQL_HOST", // contradict DATABASE_URL
	} {
		if _, ok := plain[key]; ok {
			t.Errorf("%s was sent to Railway and should not have been", key)
		}
	}

	for key, want := range map[string]string{
		"APP_NAME":   "Shop",
		"JWT_SECRET": "s3cret",
		"MAIL_FROM":  "hello@shop.test", // the quotes are stripped
	} {
		if plain[key] != want {
			t.Errorf("%s = %q, want %q", key, plain[key], want)
		}
	}

	// A deploy is production, whatever the file on the laptop says.
	if plain["APP_ENV"] != "production" {
		t.Errorf("APP_ENV = %q, want production", plain["APP_ENV"])
	}

	// Without --provision nothing is assumed about where the database is.
	if _, ok := plain["DATABASE_URL"]; ok {
		t.Error("DATABASE_URL was invented without --provision")
	}
}

// With --provision the database URLs are Railway references rather than
// copies, so rotating a password does not need a redeploy and the password
// never sits in this service's own settings.
func TestRailwayProvisionUsesReferenceVariables(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	if err := os.WriteFile(env, []byte("DATABASE_URL=postgres://localhost/dev\nREDIS_URL=redis://localhost:6379\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	vars, err := railwayVariables(RailwayConfig{EnvFile: env, Provision: true})
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]string{}
	for _, v := range vars {
		m[v.Key] = v.Value
	}

	if m["DATABASE_URL"] != "${{Postgres.DATABASE_URL}}" {
		t.Errorf("DATABASE_URL = %q, want the Railway reference", m["DATABASE_URL"])
	}
	if m["REDIS_URL"] != "${{Redis.REDIS_URL}}" {
		t.Errorf("REDIS_URL = %q, want the Railway reference", m["REDIS_URL"])
	}
	if m["DB_PROVIDER"] != "postgres" {
		t.Errorf("DB_PROVIDER = %q, want postgres", m["DB_PROVIDER"])
	}
	// The laptop's own connection string is gone rather than carried along.
	if m["DATABASE_URL"] == "postgres://localhost/dev" {
		t.Error("the local database URL was deployed")
	}
}

// The plan is sorted, so two runs print the same thing and a diff of two
// deploys is about what changed rather than about map iteration.
func TestRailwayVariablesAreOrdered(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	if err := os.WriteFile(env, []byte("ZED=1\nALPHA=2\nMID=3\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	vars, err := railwayVariables(RailwayConfig{EnvFile: env})
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(vars); i++ {
		if vars[i-1].Key > vars[i].Key {
			t.Fatalf("not sorted: %s before %s", vars[i-1].Key, vars[i].Key)
		}
	}
}

// A missing env file says what to do about it rather than what went wrong.
func TestRailwayVariablesSaysHowToGetAnEnvFile(t *testing.T) {
	_, err := railwayVariables(RailwayConfig{EnvFile: filepath.Join(t.TempDir(), "nope")})
	if err == nil {
		t.Fatal("a missing .env was accepted")
	}
	if got := err.Error(); !contains(got, "grit env") {
		t.Errorf("the error does not say how to make one: %q", got)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
