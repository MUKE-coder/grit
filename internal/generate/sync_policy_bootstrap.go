package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Reading a project's declared sync policies, so the desktop app can be told
// what it may write before it has ever reached the server (grit#214).
//
// The server is the authority and GET /sync/policy is how a client learns the
// table. That endpoint is useless in the one situation the refusal matters,
// which is a device with no connection and nothing stored yet, so the policies
// are also compiled into the desktop app and refreshed by `grit sync`.
//
// # Why this reads the whole routes directory
//
// SyncDoctor reads internal/routes/routes.go, which is where the generator
// injects a registration. A project is free to put RegisterWithPolicy in a
// file of its own applied after the // grit:sync marker, and that is a good
// pattern rather than a misuse: a regeneration cannot quietly put a table back
// to DefaultPolicy. The project in grit#214 does exactly that, so a reader
// that only opened routes.go would have found none of its policies and written
// an empty bootstrap, which is the failure mode where everything looks fine.

var (
	// syncRegisterWithPolicyRe finds a policy registration and keeps the rest
	// of the call for the field readers below. Both spellings of the receiver,
	// because a project that assembles its own registry does not have to call
	// the variable syncRegistry.
	syncRegisterWithPolicyRe = regexp.MustCompile(
		`(?s)\w+\.RegisterWithPolicy\(\s*"([a-z0-9_]+)"\s*,\s*&models\.(\w+)\{\}\s*,\s*(sync\.Policy\{.*?\n\s*\}|sync\.Policy\{[^{}]*\})`)
	// And a plain registration, which is offline_first by definition. Worth
	// recording: it is the difference between "this table is writable" and
	// "nobody has said", and only the second one should ever be guessed.
	syncRegisterPlainRe = regexp.MustCompile(
		`\w+\.Register\(\s*"([a-z0-9_]+)"\s*,\s*&models\.(\w+)\{\}\s*\)`)

	policyModeRe     = regexp.MustCompile(`Mode:\s*sync\.(\w+)`)
	policyConflictRe = regexp.MustCompile(`Conflict:\s*sync\.(\w+)`)
	// MaxOfflineAge has three spellings in real code, and reading only one of
	// them silently returns zero, which means "no limit". The generator emits
	// raw nanoseconds with the duration in a comment, because rendering the
	// parsed value keeps a second parse out of the running project:
	//
	//	MaxOfflineAge: 259200000000000 /* 72h */
	//
	// A person writes the thing a person writes:
	//
	//	MaxOfflineAge: 72 * time.Hour
	//	MaxOfflineAge: time.Hour
	policyMaxAgeRe = regexp.MustCompile(
		`MaxOfflineAge:\s*(?:(\d+)\s*\*\s*)?time\.(Hour|Minute|Second)|MaxOfflineAge:\s*(\d+)`)
)

// goConstToWire turns the Go constant name a project writes into the string
// the wire format and the client use. sync.ModeReadOnly is "read_only".
var goConstToWire = map[string]string{
	"ModeOfflineFirst":   "offline_first",
	"ModeOnlineOnly":     "online_only",
	"ModeReadOnly":       "read_only",
	"ConflictManual":     "manual",
	"ConflictServerWins": "server_wins",
	"ConflictClientWins": "client_wins",
}

// ReadDeclaredSyncPolicies reads every sync registration in a project's routes
// directory and returns what each table declares.
//
// A registration with no policy literal is reported with the default mode, so
// the bootstrap can distinguish a table the project registered as writable
// from one nobody mentioned.
func ReadDeclaredSyncPolicies(root string) ([]scaffold.SyncPolicyDeclaration, error) {
	routesDir := findRoutesDir(root)
	if routesDir == "" {
		return nil, nil
	}
	entries, err := os.ReadDir(routesDir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", routesDir, err)
	}

	// By table, so a project that registers the same table twice gets the
	// stricter answer rather than whichever file sorted last. Taking the
	// stricter one is the safe direction: a guard that is wrongly strict shows
	// the user a refusal, and one that is wrongly lax queues a doomed write.
	byTable := map[string]scaffold.SyncPolicyDeclaration{}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(routesDir, name))
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", name, err)
		}
		src := string(raw)

		for _, m := range syncRegisterWithPolicyRe.FindAllStringSubmatch(src, -1) {
			table, literal := m[1], m[3]
			declared := scaffold.SyncPolicyDeclaration{
				Table:    table,
				Mode:     wireValue(policyModeRe, literal, "offline_first"),
				Conflict: wireValue(policyConflictRe, literal, "manual"),
				Explicit: true,
			}
			declared.MaxOfflineAgeSeconds = readMaxOfflineAge(literal)
			byTable[table] = stricter(byTable[table], declared)
		}

		for _, m := range syncRegisterPlainRe.FindAllStringSubmatch(src, -1) {
			table := m[1]
			declared := scaffold.SyncPolicyDeclaration{
				Table: table, Mode: "offline_first", Conflict: "manual",
			}
			byTable[table] = stricter(byTable[table], declared)
		}
	}

	out := make([]scaffold.SyncPolicyDeclaration, 0, len(byTable))
	for _, d := range byTable {
		out = append(out, d)
	}
	return out, nil
}

// readMaxOfflineAge returns the declared limit in seconds, or 0 for none.
func readMaxOfflineAge(literal string) int {
	m := policyMaxAgeRe.FindStringSubmatch(literal)
	if m == nil {
		return 0
	}
	// Group 3: a bare number, which is nanoseconds, the form the generator
	// emits.
	if m[3] != "" {
		n, err := strconv.ParseInt(m[3], 10, 64)
		if err != nil {
			return 0
		}
		return int(time.Duration(n).Seconds())
	}
	// Groups 1 and 2: a multiple of a time unit. A bare time.Hour has no
	// count, which is one.
	count := int64(1)
	if m[1] != "" {
		n, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return 0
		}
		count = n
	}
	var unit time.Duration
	switch m[2] {
	case "Hour":
		unit = time.Hour
	case "Minute":
		unit = time.Minute
	case "Second":
		unit = time.Second
	default:
		return 0
	}
	return int((time.Duration(count) * unit).Seconds())
}

// wireValue reads one `Field: sync.Const` from a policy literal. A numeric or
// string value a project wrote directly is taken as it is.
func wireValue(re *regexp.Regexp, literal, fallback string) string {
	m := re.FindStringSubmatch(literal)
	if m == nil {
		return fallback
	}
	if wire, ok := goConstToWire[m[1]]; ok {
		return wire
	}
	// A constant this reader does not recognise: pass it through rather than
	// guessing. The client treats anything that is not offline_first as not
	// writable, so an unknown mode errs toward refusing.
	return m[1]
}

// stricter picks the less permissive of two declarations for the same table,
// and on a tie prefers the one that was actually written down.
//
// Both halves matter, and the second is the one that bit. A project registers
// a table bare in routes.go, because the generator put it there, and then
// declares a policy for it in a file of its own. Both are offline_first, so
// they tie on strictness, and before this the tie went to whichever file the
// directory listing reached first: the bare one, whose Conflict and
// MaxOfflineAge are the defaults rather than what the project asked for. The
// result looked right and quietly said manual.
func stricter(a, b scaffold.SyncPolicyDeclaration) scaffold.SyncPolicyDeclaration {
	if a.Table == "" {
		return b
	}
	rank := func(mode string) int {
		switch mode {
		case "offline_first":
			return 0
		case "online_only":
			return 1
		case "read_only":
			return 2
		default:
			// Unrecognised: stricter than anything named, because the client
			// refuses a mode it does not know.
			return 3
		}
	}
	switch {
	case rank(b.Mode) > rank(a.Mode):
		return b
	case rank(a.Mode) > rank(b.Mode):
		return a
	case b.Explicit && !a.Explicit:
		return b
	default:
		return a
	}
}

// findRoutesDir locates internal/routes in any of the project shapes.
func findRoutesDir(root string) string {
	for _, candidate := range []string{
		filepath.Join(root, "apps", "api", "internal", "routes"),
		filepath.Join(root, "api", "internal", "routes"),
		filepath.Join(root, "internal", "routes"),
	} {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return ""
}

// findDesktopRoot locates the desktop app, or "" when a project has none.
func findDesktopRoot(root string) string {
	for _, candidate := range []string{
		filepath.Join(root, "apps", "desktop"),
		filepath.Join(root, "desktop"),
	} {
		if info, err := os.Stat(filepath.Join(candidate, "sync")); err == nil && info.IsDir() {
			return candidate
		}
	}
	return ""
}

// SyncDesktopPolicyBootstrap rewrites the desktop app's compiled-in policy
// table from the project's registrations.
//
// Returns a line for `grit sync`'s output, or "" when there is nothing to say:
// no desktop app, or a bootstrap that already matches.
func SyncDesktopPolicyBootstrap(root string) (string, error) {
	desktopRoot := findDesktopRoot(root)
	if desktopRoot == "" {
		return "", nil
	}
	declared, err := ReadDeclaredSyncPolicies(root)
	if err != nil {
		return "", err
	}
	n, changed, err := scaffold.WriteSyncPolicyBootstrap(desktopRoot, declared)
	if err != nil {
		return "", err
	}
	if !changed {
		return "", nil
	}

	guarded := 0
	for _, d := range declared {
		if d.Mode != "offline_first" {
			guarded++
		}
	}
	line := fmt.Sprintf("desktop sync policies: %d table%s compiled in", n, plural(n))
	if guarded > 0 {
		verb := "refuses"
		if guarded > 1 {
			verb = "refuse"
		}
		line += fmt.Sprintf(", %d of them %s a local write before the first sync", guarded, verb)
	}
	return line, nil
}
