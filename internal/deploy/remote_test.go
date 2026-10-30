package deploy

import (
	"strings"
	"testing"
)

// Enterprise review 2.2. Every value a deploy interpolates into a remote
// command came from the project's .env or its go.mod, and only the unit file
// and the Caddyfile bodies were escaped: the two values a person is least
// likely to be able to influence.
//
// Two defences, tested separately, because quoting alone is a single point of
// failure in code that builds shell strings.

func TestQuoteSurvivesASingleQuote(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"app", `'app'`},
		{"/opt/app", `'/opt/app'`},
		{"a b", `'a b'`},
		{"$HOME", `'$HOME'`},
		{"`id`", "'`id`'"},
		{"a;rm -rf /", `'a;rm -rf /'`},
		// The one character single quotes cannot hold: closed, escaped,
		// reopened. Getting this wrong is how quoting fails.
		{"it's", `'it'\''s'`},
		{"'; id; '", `''\''; id; '\'''`},
	} {
		if got := quote(tc.in); got != tc.want {
			t.Errorf("quote(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

// validate refuses before anything is built or connected, because a deploy
// that refuses is recoverable and one that has already run `rm -rf /` is not.
func TestValidateRefusesWhatBecomesAShellCommand(t *testing.T) {
	// As Run builds it: every default is already applied by the time
	// validate runs, which is why none of these may be empty.
	ok := Config{AppName: "contacts", RemotePath: "/opt/contacts", Port: "22",
		AppPort: "8080", HealthPath: "/api/health"}
	if err := validate(ok); err != nil {
		t.Fatalf("an ordinary configuration was refused: %v", err)
	}

	for name, bad := range map[string]Config{
		"app name with a semicolon":  withName(ok, "app;id"),
		"app name with a slash":      withName(ok, "../../etc/systemd"),
		"app name with a space":      withName(ok, "my app"),
		"app name with a backtick":   withName(ok, "app`id`"),
		"app name with a dollar":     withName(ok, "app$HOME"),
		"empty app name":             withName(ok, ""),
		"relative remote path":       withPath(ok, "opt/app"),
		"remote path with traversal": withPath(ok, "/opt/../etc/app"),
		"remote path with a quote":   withPath(ok, "/opt/a'b"),
		"empty remote path":          withPath(ok, ""),
		"ssh port that is a command": withPort(ok, "22;id"),
		"app port that is a word":    withAppPort(ok, "eighty"),
		"domain with a semicolon":    withDomain(ok, "example.com;id"),
		"domain with a space":        withDomain(ok, "example.com other.com"),
		"domain with a brace":        withDomain(ok, "example.com{reverse_proxy evil}"),
		"health path with a command": withHealth(ok, "/api/health;id"),
		"health path with a space":   withHealth(ok, "/api/health x"),
		"relative health path":       withHealth(ok, "api/health"),
		"empty health path":          withHealth(ok, ""),
	} {
		if err := validate(bad); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}

	// An empty domain means no Caddy step, which is not the same as a bad one.
	if err := validate(withDomain(ok, "")); err != nil {
		t.Errorf("no domain was treated as a bad domain: %v", err)
	}
	// Real domains people deploy to.
	for _, d := range []string{"example.com", "api.example.com", "*.example.com", "example.com:8443"} {
		if err := validate(withDomain(ok, d)); err != nil {
			t.Errorf("the domain %q was refused: %v", d, err)
		}
	}
}

// A file goes over SSH inside a heredoc, so nothing in it is ever a shell
// argument. The one thing that breaks a heredoc is its own delimiter, and
// since v3.452.0 one of the files is the project's .env: content we do not
// control. Refused, with a message that says what to change.
func TestWriteFileRefusesItsOwnDelimiter(t *testing.T) {
	cmd, err := writeFileCmd("/opt/app/.env", "APP_ENV=production\nSECRET=x\n", "600")
	if err != nil {
		t.Fatalf("an ordinary .env was refused: %v", err)
	}
	for _, want := range []string{
		"sudo tee '/opt/app/.env'",
		"<<'GRIT_EOF'",
		"APP_ENV=production",
		"sudo chmod '600' '/opt/app/.env'",
		"sudo mkdir -p '/opt/app'",
	} {
		if !strings.Contains(cmd, want) {
			t.Errorf("the command does not contain %q:\n%s", want, cmd)
		}
	}

	// Nothing inside the heredoc is expanded, so these arrive as written
	// rather than as whatever the far-side shell makes of them.
	cmd, err = writeFileCmd("/opt/app/.env", "SECRET=$(id)\nHOME_LIKE=$HOME\nTICK=`id`\n", "600")
	if err != nil {
		t.Fatalf("a .env with shell characters was refused: %v", err)
	}
	for _, want := range []string{"SECRET=$(id)", "HOME_LIKE=$HOME", "TICK=`id`"} {
		if !strings.Contains(cmd, want) {
			t.Errorf("the command does not pass %q through:\n%s", want, cmd)
		}
	}

	if _, err := writeFileCmd("/opt/app/.env", "A=1\nGRIT_EOF\nB=2\n", "600"); err == nil {
		t.Error("a file holding the heredoc delimiter was accepted, which ends the heredoc early and runs the rest as commands")
	}
	if _, err := writeFileCmd("/opt/app/.env", "A=1\nGRIT_EOF\r\nB=2\n", "600"); err == nil {
		t.Error("the delimiter with a carriage return was accepted: a .env written on Windows has CRLF line endings")
	}
	// Not a bare delimiter line, so not a problem.
	if _, err := writeFileCmd("/opt/app/.env", "NOTE=GRIT_EOF is the marker\n", "600"); err != nil {
		t.Errorf("the delimiter inside a value was refused: %v", err)
	}
}

// A release is a directory of its own and `current` is a symlink to it. That
// is what makes the deploy atomic, the rollback possible, and the second
// deploy to a Linux box work at all: scp onto the file a running process is
// executing is ETXTBSY.
func TestLayoutKeepsReleasesApartFromCurrent(t *testing.T) {
	l := layoutFor(Config{RemotePath: "/opt/contacts"}, "20261030-150405")
	for name, got := range map[string]string{
		"root":     l.Root,
		"releases": l.Releases,
		"current":  l.Current,
		"release":  l.Release,
	} {
		if got == "" {
			t.Errorf("the %s path is empty", name)
		}
	}
	if l.Release == l.Current || l.Release == l.Root {
		t.Error("the release is written where the running binary already is")
	}
	if l.Release != "/opt/contacts/releases/20261030-150405" {
		t.Errorf("release = %q", l.Release)
	}

	// A trailing slash in RemotePath is a person's typo, not two directories.
	if layoutFor(Config{RemotePath: "/opt/contacts/"}, "s").Release != "/opt/contacts/releases/s" {
		t.Error("a trailing slash in the remote path doubled up")
	}
}

// The unit points at `current`, never at a release. A unit naming the release
// would have to be rewritten and reloaded on every deploy, and a rollback
// would need it rewritten again.
func TestSystemdUnitPointsAtCurrent(t *testing.T) {
	cfg := Config{AppName: "contacts", RemotePath: "/opt/contacts", AppPort: "8080"}
	l := layoutFor(cfg, "20261030-150405")
	unit := SystemdUnit(cfg, l)

	for _, want := range []string{
		"WorkingDirectory=/opt/contacts/current",
		"ExecStart=/opt/contacts/current/contacts",
		"EnvironmentFile=/opt/contacts/current/.env",
		"ReadWritePaths=/opt/contacts",
	} {
		if !strings.Contains(unit, want) {
			t.Errorf("the unit does not contain %q:\n%s", want, unit)
		}
	}
	if strings.Contains(unit, l.Release) {
		t.Errorf("the unit names a release, so every deploy would have to rewrite it:\n%s", unit)
	}
}

func withName(c Config, v string) Config    { c.AppName = v; return c }
func withPath(c Config, v string) Config    { c.RemotePath = v; return c }
func withPort(c Config, v string) Config    { c.Port = v; return c }
func withAppPort(c Config, v string) Config { c.AppPort = v; return c }
func withDomain(c Config, v string) Config  { c.Domain = v; return c }
func withHealth(c Config, v string) Config  { c.HealthPath = v; return c }

// A rollback puts `current` on an older release, so that release's files have
// to still be readable by the service. They are 0600 and owned by the service
// user, and prepareRelease used to chown the whole root recursively on every
// deploy, which took them back: a rollback would then start the app with none
// of its configuration, which is the failure mode a rollback exists to avoid.
func TestPreparingAReleaseDoesNotTouchTheOnesBeforeIt(t *testing.T) {
	cfg := Config{AppName: "contacts", RemotePath: "/opt/contacts", Port: "22",
		AppPort: "8080", HealthPath: "/api/health"}
	l := layoutFor(cfg, "20261030-150405")

	cmd := prepareReleaseCmd(l)
	if !strings.Contains(cmd, quote(l.Release)) {
		t.Errorf("the release directory is not created:\n%s", cmd)
	}
	if strings.Contains(cmd, quote(l.Root)) || strings.Contains(cmd, quote(l.Releases)) {
		t.Errorf("a deploy reaches outside its own release, so an older release's .env "+
			"stops being readable by the service and a rollback starts it with no "+
			"configuration:\n%s", cmd)
	}
}
