package scaffold

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// The exclusion is the part that is not optional: v2.8.0 reads any path ending
// in .env as an env-file probe, so without this an upgrade that changed nothing
// visible starts answering 403 for a stored object.
func TestTheLocalFileRouteIsSteppedAsideForTheWAF(t *testing.T) {
	old := `func wafExcludedRoutes() []string {
	prefix := "/api/" + APIVersion
	paths := []string{
		"/blogs", "/blogs/*",
		"/uploads", "/uploads/*",
	}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		out = append(out, prefix+p)
	}
	return out
}
`
	out, fixed, warn := repairWAFFilesExclusionSource(old)
	if len(warn) > 0 {
		t.Fatalf("warned: %v", warn)
	}
	if !strings.Contains(out, `return append(out, "/files", "/files/*")`) {
		t.Errorf("the local file route was not excluded:\n%s", out)
	}
	// And it must NOT go in the prefixed literal, which is the mistake that
	// made the first attempt at this dead config: an entry there comes back
	// as /api/v1/files, and the route is mounted at /files.
	paths := out[strings.Index(out, "paths := []string{"):strings.Index(out, "out := make")]
	if strings.Contains(paths, "/files") {
		t.Error("the entry is inside the prefixed literal, so it becomes " +
			"/api/v1/files and never matches a request")
	}
	mustParse(t, out)
	if len(fixed) != 1 {
		t.Errorf("fixed %v, want one line about it", fixed)
	}
	// Twice over the same file changes nothing the second time.
	again, fixed2, _ := repairWAFFilesExclusionSource(out)
	if again != out || len(fixed2) > 0 {
		t.Error("the repair is not idempotent")
	}
}

// A project whose wafExcludedRoutes somebody rewrote is left alone and told,
// because guessing where to put the entry is how a hand-edited file gets
// broken.
func TestAnEditedExclusionListIsLeftAloneAndNamed(t *testing.T) {
	theirs := `func wafExcludedRoutes() []string {
	return []string{"/api/v1/blogs"}
}
`
	out, fixed, warn := repairWAFFilesExclusionSource(theirs)
	if out != theirs || len(fixed) > 0 {
		t.Error("an edited list was rewritten")
	}
	if len(warn) != 1 || !strings.Contains(warn[0], "/files") {
		t.Errorf("warnings %v, want one naming the entry to add", warn)
	}
}

// And auto-block, which needs two edits that have to agree: the field, and the
// local the field reads.
func TestAutoBlockIsTurnedOnOutsideDevelopment(t *testing.T) {
	old := `// mountSentinel mounts the security suite: WAF, rate limiting, auth shield and
// anomaly detection, with its dashboard at /sentinel.
func mountSentinel(r *gin.Engine, db *gorm.DB, cfg *config.Config, svc *Services) {
	if cfg.SentinelEnabled {
		isDev := cfg.AppEnv == "development"
		// Sentinel v2 — use MountE so we can recover gracefully on
		// misconfiguration in dev instead of log.Fatalf-ing the host.
		if err := sentinel.MountE(r, db, sentinel.Config{
			Storage: sentinelStorage,
			Anomaly: sentinel.AnomalyConfig{Enabled: !isDev},
			Geo:     sentinel.GeoConfig{Enabled: !isDev},
		}); err != nil {
			log.Printf("Warning: Sentinel mount failed: %v", err)
		}
	}
}
`
	out, fixed, warn := repairSentinelAutoBlockSource(old)
	if len(warn) > 0 {
		t.Fatalf("warned: %v", warn)
	}
	if !strings.Contains(out, "AutoBlock: sentinel.AutoBlockConfig{") {
		t.Fatalf("auto-block was not added:\n%s", out)
	}
	if !strings.Contains(out, "func autoBlockExempt(cfg *config.Config) []string {") {
		t.Error("the exempt list the field reads was not declared")
	}
	// The whole point of the default list: a proxy address must not be barred.
	for _, cidr := range []string{"127.0.0.0/8", "10.0.0.0/8", "192.168.0.0/16"} {
		if !strings.Contains(out, cidr) {
			t.Errorf("%s is auto-blockable: behind an untrusted proxy that bars every user", cidr)
		}
	}
	// The two edits have to agree about where a statement may go, and
	// comparing their offsets does not test that: a declaration inserted
	// inside the config literal still comes before the field that reads it,
	// and is not Go. Parsing is the property, so parse it.
	//
	// repairSourceFile would catch this in production by refusing to write a
	// file that does not format, which means the repair would quietly never
	// apply. A test that passes while the repair can never run is worse than
	// no test.
	mustParse(t, out)
	if len(fixed) != 1 {
		t.Errorf("fixed %v, want one line about it", fixed)
	}
	again, fixed2, _ := repairSentinelAutoBlockSource(out)
	if again != out || len(fixed2) > 0 {
		t.Error("the repair is not idempotent")
	}
}

// config.go gains the field and its env read, in that order.
func TestTheAutoBlockExemptListIsReadFromTheEnvironment(t *testing.T) {
	old := `type Config struct {
	SentinelLockoutExemptIPs []string
	SentinelTurnstileSecret  string
}

func Load() *Config {
	return &Config{
		SentinelLockoutExemptIPs: splitCSV(getEnv("SENTINEL_LOCKOUT_EXEMPT_IPS", "")),
	}
}
`
	out, fixed, warn := repairSentinelAutoBlockConfigSource(old)
	if len(warn) > 0 {
		t.Fatalf("warned: %v", warn)
	}
	if !strings.Contains(out, "SentinelAutoBlockExemptIPs []string") {
		t.Error("the field is missing")
	}
	if !strings.Contains(out, `getEnv("SENTINEL_AUTOBLOCK_EXEMPT_IPS", "")`) {
		t.Error("the field is never loaded, so the variable does nothing")
	}
	if len(fixed) != 1 {
		t.Errorf("fixed %v", fixed)
	}
	again, fixed2, _ := repairSentinelAutoBlockConfigSource(out)
	if again != out || len(fixed2) > 0 {
		t.Error("the repair is not idempotent")
	}
}

// A freshly generated project is already right, so the repair must do nothing
// to it. Without this the repair could pass every test above by appending
// unconditionally, and double every entry on a new project.
func TestTheRepairsAreNoOpsOnAFreshProject(t *testing.T) {
	fresh := apiRoutesGo()
	out, fixed, warn := repairSentinelV280RoutesSource(fresh)
	if out != fresh {
		t.Error("the repair changed a freshly generated routes.go")
	}
	if len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("fixed %v, warned %v on a fresh file", fixed, warn)
	}

	mustParse(t, out)

	freshCfg := apiConfigGo()
	outCfg, fixedCfg, warnCfg := repairSentinelAutoBlockConfigSource(freshCfg)
	if outCfg != freshCfg {
		t.Error("the repair changed a freshly generated config.go")
	}
	if len(fixedCfg) > 0 || len(warnCfg) > 0 {
		t.Errorf("fixed %v, warned %v on a fresh config.go", fixedCfg, warnCfg)
	}
}

// mustParse fails unless src is Go. The fixtures are function bodies rather
// than whole files, so they are wrapped in the smallest package that holds
// them.
func mustParse(t *testing.T, src string) {
	t.Helper()
	if !strings.HasPrefix(src, "package ") {
		src = "package routes\n\n" + src
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "routes.go", src, parser.SkipObjectResolution); err != nil {
		t.Fatalf("the repaired file is not Go: %v\n\n%s", err, src)
	}
}

// The exclusion names a path, and the route it has to match is declared
// somewhere else entirely: the storage template's LocalFilesRoute. Two
// writers, nothing comparing them, which is how /api/v1/files happened. This
// compares them.
func TestTheWAFExclusionMatchesWhereTheFileRouteIsMounted(t *testing.T) {
	const decl = "const LocalFilesRoute = "
	src := storageLocalGo()
	at := strings.Index(src, decl)
	if at < 0 {
		t.Fatal("LocalFilesRoute is no longer declared in the storage template")
	}
	route := src[at+len(decl):]
	route = route[strings.Index(route, `"`)+1:]
	route = route[:strings.Index(route, `"`)]
	// "/files/*key" -> "/files"
	base := route[:strings.Index(route[1:], "/")+1]

	routes := apiRoutesGo()
	want := `return append(out, "` + base + `", "` + base + `/*")`
	if !strings.Contains(routes, want) {
		t.Errorf("the file route is mounted at %q, and wafExcludedRoutes does "+
			"not exclude %q: a stored file whose name ends in .env is refused "+
			"403 in production (Sentinel v2.8.0 reads it as an env probe)",
			route, base)
	}
	// And it is excluded unprefixed, because the route is not under the API
	// prefix.
	start := strings.Index(routes, "func wafExcludedRoutes() []string {")
	end := strings.Index(routes[start:], "\n}\n")
	paths := routes[start : start+end]
	if strings.Contains(paths[:strings.Index(paths, "out := make")], base) {
		t.Errorf("%q is in the prefixed literal, so it comes back as "+
			"/api/<version>%s, which no request ever has", base, base)
	}
}
