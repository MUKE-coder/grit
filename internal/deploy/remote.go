package deploy

import (
	"fmt"
	"regexp"
	"strings"
)

// Everything that becomes part of a remote shell command.
//
// Enterprise review 2.2: AppName, RemotePath, Domain and Port went straight
// into `ssh host "chmod +x %s/%s"` and `sudo systemctl restart %s` by string
// interpolation, and DEPLOY_HOST and DEPLOY_DOMAIN come from .env. Only the
// unit file and the Caddyfile bodies were escaped, and those are the two
// values a person is least likely to be able to influence.
//
// Two defences rather than one, because quoting alone is a single point of
// failure in code that builds shell strings:
//
//  1. Every interpolated value is single-quoted through quote().
//  2. The values that become paths, unit names and domains are checked
//     against what they are allowed to be, before any command is built. A
//     deploy that refuses is recoverable; one that runs `rm -rf /` is not.

// quote wraps a value for POSIX sh.
//
// Single quotes make everything literal except a single quote, which is
// closed, escaped and reopened. The shell on the far side is sh, whatever
// the local one is, so this is correct on Windows too.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

var (
	// A systemd unit name and a directory component. Deliberately narrower
	// than either allows: this is a project name, and one with a space or a
	// slash in it is a mistake worth catching here.
	nameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)

	// An absolute path with no traversal and no shell metacharacters.
	pathRe = regexp.MustCompile(`^/[a-zA-Z0-9._/-]{1,255}$`)

	// A hostname, optionally with a port, optionally a wildcard label. Caddy
	// accepts more than this; a deploy target does not need more than this.
	domainRe = regexp.MustCompile(`^(\*\.)?[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*(:[0-9]{1,5})?$`)

	portRe = regexp.MustCompile(`^[0-9]{1,5}$`)

	// A URL path. It is quoted before it reaches curl, so this is the second
	// defence rather than the only one.
	healthRe = regexp.MustCompile(`^/[a-zA-Z0-9._~/%-]{0,255}$`)
)

// validate refuses a configuration before anything is built or connected.
func validate(cfg Config) error {
	if !nameRe.MatchString(cfg.AppName) {
		return fmt.Errorf("the app name %q cannot be a systemd unit or a directory.\n\n"+
			"Letters, digits, dot, dash and underscore, up to 64 characters", cfg.AppName)
	}
	if !pathRe.MatchString(cfg.RemotePath) {
		return fmt.Errorf("the remote path %q is not an absolute path this command "+
			"will build a shell command around.\n\n"+
			"An absolute path of letters, digits, dot, dash, underscore and slash",
			cfg.RemotePath)
	}
	if strings.Contains(cfg.RemotePath, "..") {
		return fmt.Errorf("the remote path %q contains \"..\"", cfg.RemotePath)
	}
	if !portRe.MatchString(cfg.Port) {
		return fmt.Errorf("the SSH port %q is not a port number", cfg.Port)
	}
	if !portRe.MatchString(cfg.AppPort) {
		return fmt.Errorf("the app port %q is not a port number", cfg.AppPort)
	}
	if !healthRe.MatchString(cfg.HealthPath) {
		return fmt.Errorf("the health path %q is not a URL path.\n\n"+
			"It is fetched on the server after the restart, so it starts with "+
			"a slash: --health-path /api/health", cfg.HealthPath)
	}
	if cfg.Domain != "" && !domainRe.MatchString(cfg.Domain) {
		return fmt.Errorf("the domain %q is not a hostname.\n\n"+
			"DEPLOY_DOMAIN comes from .env, and it ends up in a Caddy site "+
			"block and a shell command", cfg.Domain)
	}
	return nil
}

// Layout is where a release lives on the server.
//
// Releases are directories and `current` is a symlink to one of them, which
// is what makes a deploy atomic and a rollback possible. Before v3.452.0 the
// binary was scp'd over the path the running process was executing, so the
// second deploy to a Linux box failed with ETXTBSY and there was nothing to
// roll back to.
type Layout struct {
	Root     string // /opt/app
	Releases string // /opt/app/releases
	Current  string // /opt/app/current  (symlink)
	Release  string // /opt/app/releases/20261030-150405
}

func layoutFor(cfg Config, stamp string) Layout {
	root := strings.TrimRight(cfg.RemotePath, "/")
	return Layout{
		Root:     root,
		Releases: root + "/releases",
		Current:  root + "/current",
		Release:  root + "/releases/" + stamp,
	}
}

// writeFileCmd builds the command that puts content at a remote path with a
// mode, without the content ever being a shell argument.
//
// Through a heredoc with a quoted delimiter, so the shell expands nothing
// inside it: a unit file or an .env holding $HOME or a backtick arrives as
// written.
//
// The one thing a heredoc cannot survive is its own delimiter alone on a
// line, and since the .env uploaded into a release is the project's rather
// than ours, that is content we do not control. Refused rather than escaped:
// a line reading GRIT_EOF in a .env is far less likely than a clever escape
// being subtly wrong, and the error says exactly what to change.
const heredocDelimiter = "GRIT_EOF"

func writeFileCmd(path, content, mode string) (string, error) {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimRight(line, "\r") == heredocDelimiter {
			return "", fmt.Errorf("the contents destined for %s have a line reading %s, "+
				"which is the marker this command uses to send a file over SSH.\n\n"+
				"Change or remove that line", path, heredocDelimiter)
		}
	}
	return fmt.Sprintf(
		"sudo mkdir -p %s && sudo tee %s > /dev/null <<'%s'\n%s\n%s\nsudo chmod %s %s",
		quote(dirOf(path)), quote(path), heredocDelimiter, content, heredocDelimiter,
		quote(mode), quote(path)), nil
}

func dirOf(path string) string {
	if i := strings.LastIndex(path, "/"); i > 0 {
		return path[:i]
	}
	return "/"
}
