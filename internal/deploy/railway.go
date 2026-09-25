package deploy

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// Deploying a Grit API to Railway.
//
// The split between the CLI and the GraphQL API is not a preference. Railway's
// API has no endpoint that accepts local source: a service is built either from
// a connected GitHub repository or from an archive the CLI uploads. So the CLI
// does the deploy, and everything around it, provisioning and variables and the
// domain, is done through the same CLI rather than half through each. One
// dependency, one auth story, and every command below appears verbatim in
// Railway's own documentation.
//
// What this does NOT do is guess. It will not create a project you did not ask
// for, it will not overwrite a variable it did not set, and it prints every
// command before running it, so --dry-run is a real plan rather than a summary.

// RailwayConfig is what `grit deploy --railway` was asked to do.
type RailwayConfig struct {
	// Service is the Railway service to deploy into. Railway requires one when
	// a project has more than a single service, and a Grit project always will
	// once a database is attached.
	Service string
	// Environment is Railway's environment, not the app's. Empty means the
	// project's default.
	Environment string
	// EnvFile is where the variables come from. Defaults to .env.
	EnvFile string
	// APIDir is the directory uploaded as the build context. The API's
	// Dockerfile copies go.mod from the context root, so the context has to be
	// the Go module, not the monorepo.
	APIDir string
	// Provision adds Postgres and Redis if the project has none, and points
	// DATABASE_URL and REDIS_URL at them.
	Provision bool
	// Domain asks Railway for a public URL once the deploy lands.
	Domain bool
	// DryRun prints the plan and changes nothing.
	DryRun bool
}

// railwayOwned are the variables Railway sets itself. Pushing our values over
// them is how a deploy ends up listening on the wrong port, or talking to the
// database it replaced rather than the one it has.
var railwayOwned = map[string]bool{
	"PORT":          true,
	"APP_PORT":      true,
	"RAILWAY_TOKEN": true,
}

// devOnly are the variables that only mean something on a laptop. Carrying
// them into a deploy is how an app ends up pointed at a database on localhost.
var devOnly = map[string]bool{
	"DB_HOST":              true,
	"DB_PORT":              true,
	"DB_USER":              true,
	"DB_PASSWORD":          true,
	"DB_NAME":              true,
	"REDIS_HOST":           true,
	"REDIS_PORT":           true,
	"REDIS_PASSWORD":       true,
	"COMPOSE_PROJECT_NAME": true,
}

// devPrefixes are whole families that exist for docker compose and nothing
// else. MinIO and Mailhog are containers on the developer's machine; on Railway
// the storage and the mail provider are configured with their own keys.
var devPrefixes = []string{
	"MINIO_", "MAILHOG_", "COMPOSE_",
	// The database is reached through DATABASE_URL on Railway, so the discrete
	// host, port and password of a compose container configure nothing and
	// contradict the thing that does.
	"MYSQL_", "POSTGRES_",
}

// localOnly reports whether a key belongs to the compose stack.
func localOnly(key string) bool {
	if devOnly[key] {
		return true
	}
	for _, p := range devPrefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// Railway deploys the API. Every step is announced before it runs.
func Railway(cfg RailwayConfig) error {
	if cfg.EnvFile == "" {
		cfg.EnvFile = ".env"
	}
	if cfg.APIDir == "" {
		cfg.APIDir = filepath.Join("apps", "api")
	}
	if cfg.Service == "" {
		cfg.Service = "api"
	}

	if err := railwayPreflight(cfg); err != nil {
		return err
	}

	vars, err := railwayVariables(cfg)
	if err != nil {
		return err
	}

	step := color.New(color.FgHiCyan, color.Bold)
	muted := color.New(color.FgHiBlack)

	if cfg.DryRun {
		color.New(color.FgHiYellow, color.Bold).Println("\n  Dry run. Nothing below is executed.")
	}

	// 1. The project. `railway status` succeeds only when this directory is
	//    already linked, which is the one thing that cannot be guessed: a new
	//    project or an existing one is the operator's decision, not ours.
	step.Println("\n  1. Project")
	if err := railwayRun(cfg, cfg.APIDir, "status"); err != nil {
		muted.Println("     This directory is not linked to a Railway project yet.")
		muted.Println("     Run one of these, then this command again:")
		color.New(color.FgHiWhite).Printf("\n       cd %s && railway init      # a new project\n", cfg.APIDir)
		color.New(color.FgHiWhite).Printf("       cd %s && railway link      # an existing one\n\n", cfg.APIDir)
		return fmt.Errorf("no Railway project is linked")
	}

	// 2. The databases, only when asked.
	if cfg.Provision {
		step.Println("\n  2. Postgres and Redis")
		for _, db := range []string{"postgres", "redis"} {
			if err := railwayRun(cfg, cfg.APIDir, "add", "--database", db); err != nil {
				// Already present is the common case and is not a failure.
				muted.Printf("     %s: already there, or Railway refused it. Continuing.\n", db)
			}
		}
	}

	// 3. Variables. One call each, because that is the documented interface.
	step.Printf("\n  %d. Variables (%d)\n", provisionStep(cfg, 3), len(vars))
	for _, kv := range vars {
		args := []string{"variable", "set", kv.Key + "=" + kv.Value}
		if cfg.Service != "" {
			args = append(args, "--service", cfg.Service)
		}
		if err := railwayRun(cfg, cfg.APIDir, args...); err != nil {
			return fmt.Errorf("setting %s: %w", kv.Key, err)
		}
	}

	// 4. The deploy itself. --ci streams the build and exits when it is done,
	//    which is what makes this usable from a script.
	step.Printf("\n  %d. Deploy\n", provisionStep(cfg, 4))
	up := []string{"up", "--ci", "--service", cfg.Service}
	if cfg.Environment != "" {
		up = append(up, "--environment", cfg.Environment)
	}
	if err := railwayRun(cfg, cfg.APIDir, up...); err != nil {
		return fmt.Errorf("railway up: %w", err)
	}

	// 5. And a URL to open.
	if cfg.Domain {
		step.Printf("\n  %d. Domain\n", provisionStep(cfg, 5))
		if err := railwayRun(cfg, cfg.APIDir, "domain"); err != nil {
			muted.Println("     Could not generate a domain. `railway domain` in the app directory will.")
		}
	}

	return nil
}

// provisionStep keeps the printed numbering honest when provisioning is off.
func provisionStep(cfg RailwayConfig, n int) int {
	if cfg.Provision {
		return n
	}
	return n - 1
}

// railwayPreflight fails early and says exactly what to do about it.
func railwayPreflight(cfg RailwayConfig) error {
	if _, err := exec.LookPath("railway"); err != nil {
		return fmt.Errorf(
			"the Railway CLI is not installed.\n" +
				"    npm i -g @railway/cli     (or: brew install railway)\n" +
				"    Then: railway login\n" +
				"    Railway's API cannot accept local source, so its CLI is what uploads the build")
	}
	if _, err := os.Stat(cfg.APIDir); err != nil {
		return fmt.Errorf("%s does not exist: run this from the root of a Grit project", cfg.APIDir)
	}
	if _, err := os.Stat(filepath.Join(cfg.APIDir, "Dockerfile")); err != nil {
		return fmt.Errorf("%s has no Dockerfile, so Railway has nothing to build", cfg.APIDir)
	}
	// A token is enough; otherwise there has to be a logged-in session.
	if os.Getenv("RAILWAY_TOKEN") == "" && os.Getenv("RAILWAY_API_TOKEN") == "" {
		if err := exec.Command("railway", "whoami").Run(); err != nil {
			return fmt.Errorf(
				"not signed in to Railway.\n" +
					"    railway login\n" +
					"    Or set RAILWAY_TOKEN, which is what CI uses")
		}
	}
	return nil
}

// EnvVar is one variable on its way to Railway.
type EnvVar struct {
	Key   string
	Value string
}

// railwayVariables reads the env file and decides what belongs in a deploy.
//
// Three rules, and each one exists because the alternative breaks something:
// Railway's own variables are left alone, the values that only mean anything on
// a laptop are dropped, and the database URLs become Railway references so they
// follow the database rather than a copy of its password taken at deploy time.
func railwayVariables(cfg RailwayConfig) ([]EnvVar, error) {
	f, err := os.Open(cfg.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w (run `grit env` to create one)", cfg.EnvFile, err)
	}
	defer f.Close()

	seen := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		// An empty value configures nothing. A generated .env is mostly
		// placeholders for providers the project does not use, and storing
		// them on Railway as empty strings only makes the dashboard unreadable.
		if key == "" || value == "" || railwayOwned[key] || localOnly(key) {
			continue
		}
		seen[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %s: %w", cfg.EnvFile, err)
	}

	// A deploy is production, whatever the file on the laptop says.
	seen["APP_ENV"] = "production"

	// Reference variables, not copies: Railway resolves these against the
	// database services at run time, so rotating a password does not need a
	// redeploy, and the password never sits in this service's own settings.
	if cfg.Provision {
		seen["DATABASE_URL"] = "${{Postgres.DATABASE_URL}}"
		seen["REDIS_URL"] = "${{Redis.REDIS_URL}}"
		seen["DB_PROVIDER"] = "postgres"
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]EnvVar, 0, len(keys))
	for _, k := range keys {
		out = append(out, EnvVar{Key: k, Value: seen[k]})
	}
	return out, nil
}

// railwayRun prints the command and then runs it, unless this is a dry run.
//
// Printed before running, and printed in full, because somebody reading their
// terminal after a failed deploy needs to be able to run the failing step by
// hand. Values are not printed: this is the one place a JWT secret would end up
// in a screenshot.
func railwayRun(cfg RailwayConfig, dir string, args ...string) error {
	shown := make([]string, len(args))
	copy(shown, args)
	for i, a := range shown {
		if k, _, ok := strings.Cut(a, "="); ok && i > 0 && shown[i-1] == "set" {
			shown[i] = k + "=********"
		}
	}
	color.New(color.FgHiBlack).Printf("     $ railway %s\n", strings.Join(shown, " "))

	if cfg.DryRun {
		return nil
	}

	cmd := exec.Command("railway", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
