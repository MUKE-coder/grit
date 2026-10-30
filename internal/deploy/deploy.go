package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Config holds deployment configuration.
type Config struct {
	Host       string // SSH host (e.g. user@server.com)
	Port       string // SSH port (default: 22)
	KeyFile    string // Path to SSH private key
	AppName    string // Application name
	RemotePath string // Remote deployment directory (default: /opt/{app-name})
	Domain     string // Domain name for Caddy reverse proxy
	AppPort    string // Port the app runs on (default: 8080)

	// APIDir is the directory holding the Go module, which is not the directory
	// deploy is run from in any layout but the oldest. apps/api in a monorepo,
	// api/ in a single, the project root in a single scaffolded before that.
	//
	// Without it, `go build` ran wherever the user happened to be standing: from
	// a monorepo root that is a directory with no go.mod, and the deploy failed
	// on "go.mod file not found" after printing that it was deploying.
	APIDir string

	// WebDir is the frontend to build before the binary, for a project whose
	// binary embeds one. Empty when the frontend deploys separately, which is
	// every monorepo: those are Next.js apps with their own hosting.
	WebDir string

	// EnvFile is the local .env uploaded into the release, at 0600.
	//
	// The unit has always said EnvironmentFile=<path>/.env and nothing ever
	// put one there, so the service started with none of the project's
	// configuration and systemd did not complain, because a missing
	// EnvironmentFile is only fatal with a leading "-" removed. Empty means
	// the project's own .env beside APIDir.
	EnvFile string

	// Keep is how many old releases stay on the server. Enough to roll back
	// more than once, few enough not to fill a five dollar VPS.
	Keep int

	// Rollback flips `current` back to the previous release and restarts,
	// without building or uploading anything.
	Rollback bool

	// HealthPath is probed after the restart, through the app's own port. A
	// deploy whose service comes back up and answers 500 is a failed deploy,
	// and before v3.452.0 it was a successful one.
	HealthPath string

	// DryRun prints what each step would do and does none of it.
	//
	// The flag existed and reached the Railway path only. On this one it was
	// accepted and ignored: a command whose help says it "runs none of them"
	// built the binary and opened an SSH connection to the host, which is the
	// opposite of what somebody asking what would happen is asking for.
	DryRun bool
}

// Run executes the full deployment pipeline.
func Run(cfg Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("deploy host is required (set DEPLOY_HOST in .env or use --host)")
	}

	if cfg.Port == "" {
		cfg.Port = "22"
	}
	if cfg.RemotePath == "" {
		cfg.RemotePath = "/opt/" + cfg.AppName
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	if cfg.Keep <= 0 {
		cfg.Keep = 5
	}
	if cfg.HealthPath == "" {
		cfg.HealthPath = "/api/health"
	}

	// Before anything is built or connected. Every value below becomes part
	// of a remote shell command (review 2.2).
	if err := validate(cfg); err != nil {
		return err
	}

	if cfg.Rollback {
		return rollback(cfg)
	}

	if cfg.DryRun {
		return printPlan(cfg)
	}

	stamp := time.Now().UTC().Format("20060102-150405")
	l := layoutFor(cfg, stamp)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Building application", func() error { return buildBinary(cfg) }},
		{"Preparing the release directory", func() error { return prepareRelease(cfg, l) }},
		{"Uploading binary", func() error { return uploadBinary(cfg, l) }},
		{"Uploading .env", func() error { return uploadEnv(cfg, l) }},
		{"Running migrations", func() error { return runMigrations(cfg, l) }},
		{"Setting up systemd service", func() error { return setupSystemd(cfg, l) }},
		{"Switching to the new release", func() error { return activate(cfg, l) }},
		{"Restarting service", func() error { return restartService(cfg) }},
	}

	if cfg.Domain != "" {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Configuring Caddy reverse proxy", func() error { return setupCaddy(cfg) }})
	}

	// What `current` pointed at before this deploy, so a failed health check
	// can put it back. Empty on the first deploy, where there is nothing to
	// go back to and the right answer is to leave the evidence.
	previous, _ := currentRelease(cfg)

	for _, step := range steps {
		fmt.Printf("  → %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	fmt.Printf("  → Checking the app answers on %s...\n", cfg.HealthPath)
	if err := waitHealthy(cfg); err != nil {
		if previous == "" {
			return fmt.Errorf("%w\n\nThis is the first release, so there is nothing "+
				"to roll back to. The release is at %s and the service is running "+
				"from it; read the logs with:\n  ssh %s 'journalctl -u %s -n 50'",
				err, l.Release, cfg.Host, cfg.AppName)
		}
		fmt.Printf("  ⚠ %v\n", err)
		fmt.Printf("  → Rolling back to %s...\n", previous)
		if rerr := activateRelease(cfg, previous); rerr != nil {
			return fmt.Errorf("the new release is unhealthy AND the rollback failed "+
				"(%v).\n\nThe server is serving %s. Put it back by hand:\n"+
				"  ssh %s 'sudo ln -sfn %s %s/current && sudo systemctl restart %s'",
				rerr, l.Release, cfg.Host, previous, l.Root, cfg.AppName)
		}
		if rerr := restartService(cfg); rerr != nil {
			return fmt.Errorf("rolled back to %s but the restart failed: %w", previous, rerr)
		}
		return fmt.Errorf("the new release did not come up, so the server was "+
			"rolled back to the one before it.\n\nThe failed release is still at "+
			"%s. Read why with:\n  ssh %s 'journalctl -u %s -n 50'",
			l.Release, cfg.Host, cfg.AppName)
	}

	// Only once it is known good. A tidy-up that runs before the health check
	// is a tidy-up that deletes the release you are about to need.
	if err := pruneReleases(cfg, l); err != nil {
		fmt.Printf("  ⚠ could not remove old releases: %v\n", err)
	}

	fmt.Printf("\n  ✓ %s is live at %s\n", cfg.AppName, l.Release)
	return nil
}

// printPlan says what a deploy would do, and does none of it.
func printPlan(cfg Config) error {
	pkg := "./cmd/server"
	if _, err := os.Stat(filepath.Join(cfg.APIDir, "cmd", "server")); err != nil {
		pkg = "."
	}
	l := layoutFor(cfg, "<timestamp>")

	fmt.Println("  Nothing below has been run.")
	fmt.Println()
	if cfg.WebDir != "" {
		fmt.Printf("  cd %s && pnpm build\n", cfg.WebDir)
	}
	fmt.Printf("  cd %s && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/%s %s\n",
		cfg.APIDir, cfg.AppName, pkg)
	if statOK(filepath.Join(cfg.APIDir, "cmd", "migrate")) {
		fmt.Printf("  cd %s && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/%s ./cmd/migrate\n",
			cfg.APIDir, migrateBinary)
	}
	fmt.Println()
	fmt.Printf("  ssh %s 'mkdir -p %s'\n", cfg.Host, l.Release)
	for _, name := range builtArtifacts(cfg) {
		fmt.Printf("  scp bin/%s %s:%s/\n", name, cfg.Host, l.Release)
	}
	fmt.Printf("  scp %s %s:%s/.env   (0600)\n", cfg.envPath(), cfg.Host, l.Release)
	if statOK(filepath.Join(cfg.APIDir, "cmd", "migrate")) {
		fmt.Printf("  ssh %s 'cd %s && ./%s'\n", cfg.Host, l.Release, migrateBinary)
	}
	fmt.Printf("  ssh %s 'write /etc/systemd/system/%s.service pointing at %s/current'\n",
		cfg.Host, cfg.AppName, l.Root)
	fmt.Printf("  ssh %s 'ln -sfn %s %s/current'\n", cfg.Host, l.Release, l.Root)
	fmt.Printf("  ssh %s 'systemctl restart %s'\n", cfg.Host, cfg.AppName)
	if cfg.Domain != "" {
		fmt.Printf("  ssh %s 'write /etc/caddy/sites/%s.caddy for %s -> :%s, then caddy validate && systemctl reload caddy'\n",
			cfg.Host, cfg.AppName, cfg.Domain, cfg.AppPort)
	}
	fmt.Println()
	fmt.Printf("  then GET http://127.0.0.1:%s%s on the server, and if it does not\n",
		cfg.AppPort, cfg.HealthPath)
	fmt.Printf("  answer, point current back at the release before this one and restart.\n")
	fmt.Printf("  Old releases beyond the newest %d are removed, but only once the\n", cfg.Keep)
	fmt.Printf("  new one has answered.\n")
	fmt.Println()
	fmt.Println("  Run it without --dry-run to do all of that.")
	return nil
}

// statOK is "this path exists", for deciding whether a step is in the plan.
func statOK(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// buildBinary cross-compiles the Go binary for Linux.
func buildBinary(cfg Config) error {
	binaryName := cfg.AppName
	outputPath, err := filepath.Abs(filepath.Join("bin", binaryName))
	if err != nil {
		return err
	}

	// The frontend first, for a binary that carries one: the build has to be in
	// place before //go:embed reads it, or the binary ships the placeholder.
	if cfg.WebDir != "" {
		fmt.Println("    Building frontend...")
		feCmd := exec.Command("pnpm", "build")
		feCmd.Dir = cfg.WebDir
		feCmd.Stdout = os.Stdout
		feCmd.Stderr = os.Stderr
		if err := feCmd.Run(); err != nil {
			return fmt.Errorf("building frontend: %w", err)
		}
	}

	// The main package, which is cmd/server in a monorepo and the top of the
	// module in a single.
	pkg := "./cmd/server"
	if _, err := os.Stat(filepath.Join(cfg.APIDir, "cmd", "server")); err != nil {
		pkg = "."
	}

	cmd := exec.Command("go", "build", "-o", outputPath, pkg)
	cmd.Dir = cfg.APIDir
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// And the migrate command, so the deploy can bring the schema with it.
	// A deploy that ships code expecting a column nobody added is the most
	// common way a working release breaks on arrival, and there was no
	// migration step at all before v3.452.0 (review 2.2).
	//
	// Absent in a project that has no cmd/migrate, which is fine: the deploy
	// skips the step rather than failing on it.
	migrateSrc := filepath.Join(cfg.APIDir, "cmd", "migrate")
	if _, err := os.Stat(migrateSrc); err != nil {
		return nil
	}
	migrateOut, err := filepath.Abs(filepath.Join(buildDir(cfg), migrateBinary))
	if err != nil {
		return err
	}
	mcmd := exec.Command("go", "build", "-o", migrateOut, "./cmd/migrate")
	mcmd.Dir = cfg.APIDir
	mcmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0")
	mcmd.Stdout = os.Stdout
	mcmd.Stderr = os.Stderr
	return mcmd.Run()
}

// migrateBinary is the schema tool's name inside a release.
const migrateBinary = "grit-migrate"

// buildDir is where the local build puts its artifacts.
func buildDir(cfg Config) string { return "bin" }

// builtArtifacts are the files a release is made of, in upload order.
func builtArtifacts(cfg Config) []string {
	return []string{cfg.AppName, migrateBinary}
}

// setupSystemd creates a systemd service unit file on the remote server.
func setupSystemd(cfg Config, l Layout) error {
	unit := SystemdUnit(cfg, l)
	serviceName := cfg.AppName + ".service"
	path := "/etc/systemd/system/" + serviceName

	// Through a heredoc rather than echo with hand-escaped quotes: the unit
	// holds paths and Environment lines, and a dollar sign in one of those
	// used to be a shell expansion on the far side.
	cmd, err := writeFileCmd(path, unit, "644")
	if err != nil {
		return err
	}
	if err := sshCmd(cfg, cmd); err != nil {
		return err
	}
	return sshCmd(cfg, fmt.Sprintf("sudo systemctl daemon-reload && sudo systemctl enable %s",
		quote(serviceName)))
}

// restartService restarts the systemd service.
func restartService(cfg Config) error {
	return sshCmd(cfg, fmt.Sprintf("sudo systemctl restart %s", quote(cfg.AppName)))
}

// setupCaddy writes this app's site block and nothing else.
//
// It used to write the whole /etc/caddy/Caddyfile, so deploying a second app
// to a box, or deploying one to a box that already served anything, threw the
// other site away without a word (review 2.2). For a command whose pitch is
// "a five dollar VPS you own", that is the opposite of owning it.
//
// The site lives in /etc/caddy/sites/<app>.caddy and the main Caddyfile gets
// one import line if it does not have one. Appended, never rewritten: that
// file is the operator's, and every other site in it stays.
func setupCaddy(cfg Config) error {
	site := "/etc/caddy/sites/" + cfg.AppName + ".caddy"
	cmd, err := writeFileCmd(site, Caddyfile(cfg), "644")
	if err != nil {
		return err
	}
	if err := sshCmd(cfg, cmd); err != nil {
		return err
	}

	// grep before appending, so a second deploy does not add the line twice.
	const importLine = "import sites/*.caddy"
	ensure := fmt.Sprintf(
		"sudo touch /etc/caddy/Caddyfile && "+
			"(grep -qF %s /etc/caddy/Caddyfile || "+
			"echo %s | sudo tee -a /etc/caddy/Caddyfile > /dev/null)",
		quote(importLine), quote(importLine))
	if err := sshCmd(cfg, ensure); err != nil {
		return err
	}

	// Validate before reloading. A Caddyfile that does not parse takes every
	// site on the box down, including the ones this deploy has nothing to do
	// with, and finding out from "reload failed" is too late.
	if err := sshCmd(cfg, "sudo caddy validate --config /etc/caddy/Caddyfile"); err != nil {
		return fmt.Errorf("the Caddy configuration does not validate with this site "+
			"added, so it was not reloaded and the server still serves what it "+
			"did. The site file is at %s: %w", site, err)
	}
	return sshCmd(cfg, "sudo systemctl reload caddy")
}

// sshOutput runs a command on the server and returns what it printed.
//
// For the three things a deploy has to read rather than do: where `current`
// points, what releases exist, and what the health endpoint answered.
func sshOutput(cfg Config, command string) (string, error) {
	args := []string{"-p", cfg.Port}
	if cfg.KeyFile != "" {
		args = append(args, "-i", cfg.KeyFile)
	}
	args = append(args, cfg.Host, command)

	cmd := exec.Command("ssh", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(out), err
}

// sshCmd executes a command on the remote server.
func sshCmd(cfg Config, command string) error {
	args := []string{"-p", cfg.Port}
	if cfg.KeyFile != "" {
		args = append(args, "-i", cfg.KeyFile)
	}
	args = append(args, cfg.Host, command)

	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SystemdUnit generates a systemd service unit file.
func SystemdUnit(cfg Config, l Layout) string {
	// Every path is <root>/current, not the release, so the unit is written
	// once and a deploy moves only the symlink. Writing the release path in
	// here would mean rewriting and reloading the unit on every deploy, and
	// a rollback would need the unit rewritten too.
	current := l.Root + "/current"

	return fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=%s
ExecStart=%s/%s
Restart=on-failure
RestartSec=5
Environment=APP_ENV=production
Environment=APP_PORT=%s
EnvironmentFile=%s/.env

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=%s

[Install]
WantedBy=multi-user.target
`, cfg.AppName, current, current, cfg.AppName, cfg.AppPort, current, l.Root)
}

// Caddyfile generates a Caddy reverse proxy configuration with auto-TLS.
func Caddyfile(cfg Config) string {
	return fmt.Sprintf(`%s {
    reverse_proxy localhost:%s

    encode gzip

    header {
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }

    log {
        output file /var/log/caddy/%s.log {
            roll_size 10mb
            roll_keep 5
        }
    }
}
`, cfg.Domain, cfg.AppPort, cfg.AppName)
}

// Timestamp returns a deployment timestamp for logging.
func Timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
