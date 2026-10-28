package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Sentinel v2.8.0 arrives through upgrade_deps, so a project that never
// touched its security config still gets the new behaviour. Two parts of it
// need the project's routes.go to change with it.
//
// The first is not optional. v2.8.0 reads any path ending in .env as an
// env-file probe, which is what a scanner asks for and is also what the local
// storage driver serves when an "accepts all" upload field keeps a file's
// original extension. In production the WAF blocks, so a stored object would
// start answering 403 on an upgrade that changed nothing anyone could see.
//
// The second is AutoBlock, off until it is turned on.

// wafFilesReturn replaces wafExcludedRoutes' return, adding the paths that are
// not under the API prefix.
//
// It goes after the loop rather than into the paths literal because that
// literal is prefixed: "/files" in it becomes "/api/v1/files", which no
// request ever has. The local storage route is mounted on the engine itself.
const wafFilesOldReturn = "\treturn out\n}\n"

const wafFilesNewReturn = `	// The local storage driver's file route, which is mounted on the engine
	// rather than under the API prefix, so it must not be run through it.
	// It serves a stored object under the key it was stored with, and an
	// "accepts all" upload field keeps a file's original extension, so a
	// stored config.env is a real path this app has to answer. Sentinel
	// v2.8.0 reads any path ending in .env as an env-file probe, and in
	// production the WAF blocks: that object answered 403.
	//
	// Traversal is not what is given up here. The storage layer refuses a key
	// with a . or .. segment, a backslash, a Windows drive letter, or any key
	// resolving outside the disk's root, whether or not the WAF is looking.
	return append(out, "/files", "/files/*")
}
`

// sentinelAutoBlockBlock is the config field, and the comment that says why it
// is on.
const sentinelAutoBlockBlock = `			// An address that keeps attacking is blocked, on what it has
			// done to this application rather than on a third party's
			// reputation feed, which is the only thing that stops a scanner
			// nobody has reported yet. Sentinel's defaults: ten
			// high-severity threats inside ten minutes bars the address for
			// a day. A high severity is what the scanner patterns score (LFI,
			// SQLi, path traversal), so ordinary odd traffic does not
			// accumulate towards a block.
			//
			// Off in development, where the thing it would block is you.
			// Counted through Counters above, so threats spread over replicas
			// add up to one total and a block made anywhere holds everywhere.
			AutoBlock: sentinel.AutoBlockConfig{
				Enabled:   !isDev,
				ExemptIPs: autoBlockExempt(cfg),
			},
`

// sentinelAutoBlockExemptFunc is the exempt list the field reads, and the
// reason the private ranges are in it.
const sentinelAutoBlockExemptFunc = `// autoBlockNeverBlocked are the addresses auto-block is never allowed to bar,
// whatever they do.
//
// Auto-block bars an address, and an address is only the attacker's when
// Sentinel is reading the real client. Behind a reverse proxy it has not been
// told to trust, every request carries the proxy's own address instead, and
// TrustedProxies is empty by default because trusting a forwarded header you
// have not configured is worse. So on a deployment where Caddy, Traefik,
// nginx or an ingress terminates TLS and SENTINEL_TRUSTED_PROXIES was never
// set, ten probes through that proxy would bar the proxy: every user, for a
// day, from one scanner.
//
// Loopback and the private ranges are where such a proxy appears, so they are
// exempt. Nothing is given up by it that you would want: a client genuinely on
// one of these ranges is on your own network, and barring one address there is
// not the control you want anyway.
//
// Set SENTINEL_TRUSTED_PROXIES and the real client address is read, after
// which this list stops applying to anybody who matters.
var autoBlockNeverBlocked = []string{
	"127.0.0.0/8", "::1/128",
	"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
	"100.64.0.0/10", // carrier NAT, and what some ingress controllers use
	"fc00::/7",      // IPv6 unique local
}

// autoBlockExempt is that list plus whatever this project named.
func autoBlockExempt(cfg *config.Config) []string {
	named := cfg.SentinelAutoBlockExemptIPs
	if len(named) == 0 {
		// Whose egress address is thousands of people is one question, and a
		// project that answered it for the sign-in lockout has answered it for
		// auto-block too. SENTINEL_AUTOBLOCK_EXEMPT_IPS parts them.
		named = cfg.SentinelLockoutExemptIPs
	}
	out := make([]string, 0, len(autoBlockNeverBlocked)+len(named))
	out = append(out, autoBlockNeverBlocked...)
	return append(out, named...)
}

`

// repairSentinelV280 brings a project's routes.go and config.go up to what
// Sentinel v2.8.0 expects.
func repairSentinelV280(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairSentinelV280RoutesSource); err != nil {
			return err
		}
	}
	cfg := filepath.Join(apiRoot, "internal", "config", "config.go")
	if fileExists(cfg) {
		if err := repairSourceFile(root, m, cfg, repairSentinelAutoBlockConfigSource); err != nil {
			return err
		}
	}
	return nil
}

// repairWAFFilesExclusionSource adds the local file route to the WAF's
// exclusions. It is anchored on the uploads entry, which every generated
// wafExcludedRoutes has.
func repairWAFFilesExclusionSource(src string) (string, []string, []string) {
	const decl = "func wafExcludedRoutes() []string {"
	start := strings.Index(src, decl)
	if start < 0 || strings.Contains(src, `"/files", "/files/*"`) {
		return src, nil, nil
	}
	pos := strings.Index(src[start:], wafFilesOldReturn)
	if pos < 0 {
		return src, nil, []string{`Sentinel v2.8.0 reads a path ending in .env as a probe, and wafExcludedRoutes does not end where Grit wrote it: make it return append(out, "/files", "/files/*"), or a stored file whose name ends in .env is refused 403 in production`}
	}
	at := start + pos
	out := src[:at] + wafFilesNewReturn + src[at+len(wafFilesOldReturn):]
	// The capacity hint too, so an upgraded wafExcludedRoutes is the same
	// bytes as a freshly generated one. It changes nothing a request sees;
	// a file that differs from the template is a file the next repair has
	// to special-case.
	out = strings.Replace(out, "out := make([]string, 0, len(paths))",
		"out := make([]string, 0, len(paths)+2)", 1)
	return out,
		[]string{"the local storage route is no longer read as an env-file probe (Sentinel v2.8.0)"}, nil
}

// repairSentinelAutoBlockSource turns auto-block on outside development.
func repairSentinelAutoBlockSource(src string) (string, []string, []string) {
	if strings.Contains(src, "sentinel.AutoBlockConfig{") {
		return src, nil, nil
	}
	if !strings.Contains(src, "sentinel.MountE(r, db, sentinel.Config{") {
		return src, nil, nil
	}
	const anchor = "\t\t\tAnomaly: sentinel.AnomalyConfig{Enabled: !isDev},\n"
	at := strings.Index(src, anchor)
	if at < 0 {
		return src, nil, []string{"its Sentinel config is not the one Grit wrote, so auto-block was left off: add AutoBlock: sentinel.AutoBlockConfig{Enabled: !isDev} to turn it on"}
	}
	// The field reads autoBlockExempt, so the function has to exist. It goes
	// above mountSentinel, which is the only thing that calls it.
	const funcAnchor = "// mountSentinel mounts the security suite"
	funcAt := strings.Index(src, funcAnchor)
	if funcAt < 0 {
		return src, nil, []string{"its mountSentinel is not the one Grit wrote, so auto-block was left off"}
	}
	out := src[:at] + sentinelAutoBlockBlock + src[at:]
	funcAt = strings.Index(out, funcAnchor)
	out = out[:funcAt] + sentinelAutoBlockExemptFunc + out[funcAt:]
	return out,
		[]string{"an address that keeps attacking is blocked outside development (Sentinel v2.8.0 AutoBlock)"}, nil
}

// repairSentinelV280RoutesSource is both routes.go edits, so one read and one
// write carry them.
func repairSentinelV280RoutesSource(src string) (string, []string, []string) {
	out, fixed, warn := repairWAFFilesExclusionSource(src)
	out2, fixed2, warn2 := repairSentinelAutoBlockSource(out)
	return out2, append(fixed, fixed2...), append(warn, warn2...)
}

// addBothExemptLists gives a project that has neither exempt list both of
// them, anchored on the end of the Config struct and the end of the literal
// the loader returns rather than on any field, since by definition the fields
// this would anchor on are the ones that are missing.
func addBothExemptLists(src string) (string, []string, []string) {
	at, ok := endOfStruct(src, "type Config struct {")
	if !ok {
		return src, nil, []string{"its config.go has no `type Config struct` that Grit can find, so the " +
			"two exempt lists cannot be added: add SentinelLockoutExemptIPs and " +
			"SentinelAutoBlockExemptIPs ([]string) by hand, or internal/routes will not compile"}
	}
	src = src[:at] + "\t// Addresses or CIDRs the sign-in lockout and auto-block never bar,\n" +
		"\t// whatever they do: your own monitoring, a partner's integration, an\n" +
		"\t// office range. Auto-block falls back to the lockout list when its own\n" +
		"\t// is empty, because both answer the same question.\n" +
		"\tSentinelLockoutExemptIPs   []string\n" +
		"\tSentinelAutoBlockExemptIPs []string\n" + src[at:]

	lat, ok := endOfStruct(src, "cfg := &Config{")
	if !ok {
		lat, ok = endOfStruct(src, "return &Config{")
	}
	if !ok {
		return src, nil, []string{"its config.go never builds a Config literal that Grit can find, so the " +
			"exempt lists are never read from the environment: add the two splitCSV(getEnv(...)) " +
			"lines by hand"}
	}
	src = src[:lat] +
		"\t\tSentinelLockoutExemptIPs:   splitCSV(getEnv(\"SENTINEL_LOCKOUT_EXEMPT_IPS\", \"\")),\n" +
		"\t\tSentinelAutoBlockExemptIPs: splitCSV(getEnv(\"SENTINEL_AUTOBLOCK_EXEMPT_IPS\", \"\")),\n" +
		src[lat:]
	return src, []string{"both Sentinel exempt lists are read from the environment"}, nil
}

// repairSentinelAutoBlockConfigSource adds the field auto-block's exempt list
// is read from.
func repairSentinelAutoBlockConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "SentinelAutoBlockExemptIPs") {
		return src, nil, nil
	}
	if !strings.Contains(src, "SentinelLockoutExemptIPs") {
		// A project from before either exempt list. This used to return here,
		// reasoning that auto-block still works and simply has nothing to
		// exempt. That was true while nothing else named the fields. routes.go
		// names both now and is rewritten from the template on every upgrade,
		// so returning here leaves routes.go referring to two fields that do
		// not exist and an API that does not compile.
		return addBothExemptLists(src)
	}
	const field = "\tSentinelLockoutExemptIPs []string\n"
	at := strings.Index(src, field)
	if at < 0 {
		return src, nil, []string{"its SentinelLockoutExemptIPs field is not where Grit wrote it, so SENTINEL_AUTOBLOCK_EXEMPT_IPS is not read"}
	}
	src = src[:at+len(field)] + "\t// Addresses or CIDRs that auto-block never bars, whatever they do: your\n" +
		"\t// own monitoring, a partner's integration, an office range. Empty, the\n" +
		"\t// sign-in lockout's exempt list above is used instead, because both\n" +
		"\t// answer the same question and a project that answered it once should\n" +
		"\t// not have to answer it twice.\n" +
		"\tSentinelAutoBlockExemptIPs []string\n" + src[at+len(field):]

	const load = "\t\tSentinelLockoutExemptIPs: splitCSV(getEnv(\"SENTINEL_LOCKOUT_EXEMPT_IPS\", \"\")),\n"
	lat := strings.Index(src, load)
	if lat < 0 {
		return src, nil, []string{"its SentinelLockoutExemptIPs is never loaded from the environment, so SENTINEL_AUTOBLOCK_EXEMPT_IPS is not read either"}
	}
	src = src[:lat+len(load)] +
		"\t\tSentinelAutoBlockExemptIPs: splitCSV(getEnv(\"SENTINEL_AUTOBLOCK_EXEMPT_IPS\", \"\")),\n" +
		src[lat+len(load):]
	return src, []string{"SENTINEL_AUTOBLOCK_EXEMPT_IPS is read"}, nil
}
