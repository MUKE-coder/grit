package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// A deploy is a new directory and one symlink.
//
//	/opt/app/releases/20261030-150405/{app,migrate,.env}
//	/opt/app/current -> releases/20261030-150405
//
// Enterprise review 2.2. Before this the binary was scp'd straight over the
// path the running process was executing, which on Linux is ETXTBSY the
// second time anybody deploys; nothing checked that the service came back;
// and there was no previous version to return to when it did not.

// prepareRelease makes the directory this deploy will fill.
//
// The chown is scoped to this release, not to the root. Chowning the root
// recursively takes every older release's .env back off the service user, and
// those files are 0600: a rollback would then put `current` on a release whose
// configuration the service cannot read, and it would start with none of it.
func prepareRelease(cfg Config, l Layout) error {
	return sshCmd(cfg, prepareReleaseCmd(l))
}

// prepareReleaseCmd is the command itself, so a test can read it without a
// server to run it against.
func prepareReleaseCmd(l Layout) string {
	return fmt.Sprintf("sudo mkdir -p %s && sudo chown -R $(whoami) %s",
		quote(l.Release), quote(l.Release))
}

// uploadBinary copies the built artifacts into the new release directory.
//
// Into the NEW directory, which is the point: the running binary is in the
// previous one and is never written over.
func uploadBinary(cfg Config, l Layout) error {
	for _, name := range builtArtifacts(cfg) {
		local := filepath.Join(buildDir(cfg), name)
		if _, err := os.Stat(local); err != nil {
			if name == migrateBinary {
				continue // a project without a migrate command
			}
			return fmt.Errorf("built artifact missing: %s", local)
		}
		args := []string{"-P", cfg.Port}
		if cfg.KeyFile != "" {
			args = append(args, "-i", cfg.KeyFile)
		}
		args = append(args, local, cfg.Host+":"+l.Release+"/"+name)

		cmd := exec.Command("scp", args...)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("uploading %s: %w", name, err)
		}
		if err := sshCmd(cfg, fmt.Sprintf("chmod +x %s",
			quote(l.Release+"/"+name))); err != nil {
			return err
		}
	}
	return nil
}

// envPath is the local file uploaded as the release's .env.
//
// One function, so --dry-run names the file the deploy would actually read
// rather than a guess at it.
func (c Config) envPath() string {
	if c.EnvFile != "" {
		return c.EnvFile
	}
	path := filepath.Join(c.APIDir, ".env")
	if _, err := os.Stat(path); err != nil {
		return ".env"
	}
	return path
}

// uploadEnv puts the project's .env into the release at 0600.
//
// The unit has always named EnvironmentFile=<release>/.env and nothing ever
// uploaded one, so the service started with none of the project's
// configuration: no database URL, no secrets, defaults everywhere. systemd
// does not complain about a missing EnvironmentFile, so it looked fine.
//
// 0600 and owned by the service user, because it holds every secret the app
// has. Over SSH rather than scp so the mode is never briefly wider.
func uploadEnv(cfg Config, l Layout) error {
	path := cfg.envPath()
	body, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s to upload: %w\n\n"+
			"The service reads its configuration from this file. Point --env-file "+
			"at it, or deploy from the directory that has one", path, err)
	}
	cmd, err := writeFileCmd(l.Release+"/.env", string(body), "600")
	if err != nil {
		return err
	}
	if err := sshCmd(cfg, cmd); err != nil {
		return err
	}
	return sshCmd(cfg, fmt.Sprintf("sudo chown www-data:www-data %s",
		quote(l.Release+"/.env")))
}

// runMigrations runs the release's own migrate binary before the service
// switches to it.
//
// Before the switch, deliberately: the old binary is still serving, and a
// migration that fails must not leave the new one running against a schema
// it does not expect. Additive migrations make that ordering safe, which is
// the kind of migration AutoMigrate produces.
func runMigrations(cfg Config, l Layout) error {
	probe := fmt.Sprintf("test -x %s", quote(l.Release+"/"+migrateBinary))
	if err := sshCmd(cfg, probe); err != nil {
		fmt.Println("    no migrate binary in this release, skipping")
		return nil
	}
	return sshCmd(cfg, fmt.Sprintf("cd %s && ./%s",
		quote(l.Release), migrateBinary))
}

// activate points `current` at the new release.
//
// One symlink, replaced atomically by ln -sfn, so a request in flight either
// sees the old release or the new one and never half of a copy.
func activate(cfg Config, l Layout) error { return activateRelease(cfg, l.Release) }

func activateRelease(cfg Config, release string) error {
	root := strings.TrimRight(cfg.RemotePath, "/")
	return sshCmd(cfg, fmt.Sprintf("sudo ln -sfn %s %s",
		quote(release), quote(root+"/current")))
}

// currentRelease reads where `current` points, for the rollback target.
func currentRelease(cfg Config) (string, error) {
	root := strings.TrimRight(cfg.RemotePath, "/")
	out, err := sshOutput(cfg, fmt.Sprintf("readlink -f %s 2>/dev/null || true",
		quote(root+"/current")))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// rollback moves `current` back one release and restarts.
func rollback(cfg Config) error {
	root := strings.TrimRight(cfg.RemotePath, "/")

	now, err := currentRelease(cfg)
	if err != nil {
		return err
	}
	// Newest first, so the second line is the one before the current.
	out, err := sshOutput(cfg, fmt.Sprintf("ls -1 %s 2>/dev/null | sort -r || true",
		quote(root+"/releases")))
	if err != nil {
		return err
	}
	var target string
	for _, name := range strings.Fields(out) {
		candidate := root + "/releases/" + name
		if candidate != now {
			target = candidate
			break
		}
	}
	if target == "" {
		return fmt.Errorf("there is no earlier release on %s to roll back to.\n\n"+
			"Releases live in %s/releases and this server has %s",
			cfg.Host, root, plural(len(strings.Fields(out))))
	}

	fmt.Printf("  → Rolling back to %s...\n", target)
	if err := activateRelease(cfg, target); err != nil {
		return err
	}
	if err := restartService(cfg); err != nil {
		return err
	}
	if err := waitHealthy(cfg); err != nil {
		return fmt.Errorf("rolled back to %s and it is not answering either: %w",
			target, err)
	}
	fmt.Printf("\n  ✓ %s is serving %s\n", cfg.AppName, target)
	return nil
}

func plural(n int) string {
	switch n {
	case 0:
		return "none"
	case 1:
		return "one"
	}
	return fmt.Sprintf("%d", n)
}

// pruneReleases keeps the most recent Keep directories and removes the rest.
//
// After the health check, never before: a tidy-up that runs first is one that
// deletes the release you are about to need.
func pruneReleases(cfg Config, l Layout) error {
	return sshCmd(cfg, fmt.Sprintf(
		"cd %s && ls -1 | sort -r | tail -n +%d | xargs -r -I{} sudo rm -rf -- {}",
		quote(l.Releases), cfg.Keep+1))
}

// waitHealthy polls the app's health endpoint from the server itself.
//
// From the server, through the app's own port, so it tests the service
// rather than the reverse proxy or DNS: a deploy that works and a
// certificate that has not issued yet are different problems and should
// not look the same.
func waitHealthy(cfg Config) error {
	url := fmt.Sprintf("http://127.0.0.1:%s%s", cfg.AppPort, cfg.HealthPath)
	probe := fmt.Sprintf("curl -fsS -m 3 -o /dev/null -w '%%{http_code}' %s", quote(url))

	deadline := time.Now().Add(60 * time.Second)
	var last string
	for time.Now().Before(deadline) {
		out, err := sshOutput(cfg, probe)
		code := strings.TrimSpace(out)
		if err == nil && strings.HasPrefix(code, "2") {
			return nil
		}
		if code != "" {
			last = code
		}
		time.Sleep(2 * time.Second)
	}
	if last != "" {
		return fmt.Errorf("the app answered %s on %s for 60 seconds", last, cfg.HealthPath)
	}
	return fmt.Errorf("the app did not answer on %s within 60 seconds", cfg.HealthPath)
}
