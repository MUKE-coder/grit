package scaffold

import (
	"strings"
	"testing"
)

// The threshold is "more than", not "at least", and the difference is one
// address. At the default of 3 an address is banned on its fourth
// high-severity threat, and a flip to `<` would ban it on its third: a silent
// change to who gets locked out of an application, which nothing else here
// would catch.
func TestTheSweepBansAboveTheThresholdAndNotAtIt(t *testing.T) {
	src := securityAutoBanGo()

	if !strings.Contains(src, "if n <= cfg.MinThreats {\n\t\t\tcontinue\n\t\t}") {
		t.Error("the threshold comparison is not `n <= cfg.MinThreats`, so the sweep " +
			"no longer bans strictly above the threshold")
	}
	if !strings.Contains(src, "cfg.MinThreats = 3") {
		t.Error("the default threshold is not 3")
	}
	// The comment is the only place the semantics are written down for the
	// next reader, and a default of 3 that bans at 3 is the misreading.
	if !strings.Contains(src, "is a threshold to exceed") && !strings.Contains(src, "exceed, not") {
		t.Error("nothing in the sweep says the threshold is exceeded rather than met")
	}
}

// Behind a reverse proxy that Sentinel has not been told to trust, every
// request carries the proxy's address. A sweep that banned it would ban every
// user of the application, on a schedule, with nobody watching.
func TestTheSweepNeverBansItsOwnInfrastructure(t *testing.T) {
	src := securityAutoBanGo()

	if !strings.Contains(src, "if IsPrivateOrLoopbackIP(c.ip) {") {
		t.Error("the sweep no longer skips loopback and private addresses")
	}
	for _, check := range []string{"IsLoopback()", "IsPrivate()", "isCarrierNAT"} {
		if !strings.Contains(src, check) {
			t.Errorf("%s is no longer part of the address test", check)
		}
	}
	// An address it cannot parse is not an address it should ban.
	if !strings.Contains(src, "if parsed == nil {") || !strings.Contains(src, "return true") {
		t.Error("an unparseable address is no longer treated as not-bannable")
	}
	// And the hand-ban endpoint refuses the same addresses, through the same
	// function, so the two cannot disagree.
	if !strings.Contains(securityHandlerGo(), "services.IsPrivateOrLoopbackIP(req.IP)") {
		t.Error("the manual ban endpoint no longer refuses private addresses")
	}
}

// An address already banned must not be banned again: Sentinel would accept it
// and the sweep would report work it did not do, which is how a log stops
// being worth reading.
func TestTheSweepSkipsWhatIsAlreadyBanned(t *testing.T) {
	src := securityAutoBanGo()
	if !strings.Contains(src, "if already[c.ip] {") {
		t.Error("the sweep no longer reads the ban list before banning")
	}
	if !strings.Contains(src, "/sentinel/api/ip/blocked") {
		t.Error("the sweep never fetches the current bans")
	}
	if !strings.Contains(src, "MaxBans") {
		t.Error("one sweep is no longer capped, so a log that has gone wrong can " +
			"ban without limit before anybody sees it")
	}
}

// The dashboard tells an automatic ban from a hand-made one by the reason
// string, because Sentinel's BlockedIP carries no flag for it. So the word the
// sweep writes and the word the dashboard reads have to be the same word.
func TestTheBanReasonTheSweepWritesIsTheOneTheDashboardReads(t *testing.T) {
	src := securityAutoBanGo()
	if !strings.Contains(src, `AutoBanReason = "auto-ban"`) {
		t.Fatal("the auto-ban reason prefix changed")
	}
	// Written with the prefix...
	if !strings.Contains(src, `fmt.Sprintf("%s: %d high-severity threats in %s",`) ||
		!strings.Contains(src, "AutoBanReason, c.count") {
		t.Error("the ban reason is no longer built from AutoBanReason, so IsAutoBan " +
			"will not recognise it and every automatic ban reads as manual")
	}
	// ...and read back through the same constant.
	if !strings.Contains(src, "strings.HasPrefix(r, AutoBanReason)") {
		t.Error("IsAutoBan no longer tests the prefix the sweep writes")
	}
	// Sentinel's own burst rule writes its own wording, and that counts as
	// automatic too.
	if !strings.Contains(src, `strings.HasPrefix(r, "auto-blocked")`) {
		t.Error("a ban made by Sentinel's own AutoBlock is no longer counted as automatic")
	}
	if !strings.Contains(securityHandlerGo(), "services.IsAutoBan(b.Reason)") {
		t.Error("the summary no longer splits bans into automatic and manual")
	}
}

// The repairs are no-ops on a project that already has the sweep, which is
// what makes running upgrade twice safe.
func TestTheAutoBanRepairsAreIdempotent(t *testing.T) {
	cases := []struct {
		name string
		fn   func(string) (string, []string, []string)
		src  string
	}{
		{"config", repairAutoBanConfigSource, "\tSentinelTurnstileSecret    string\n\n\t\tSentinelTurnstileSecret:    getEnv(\"SENTINEL_TURNSTILE_SECRET\", \"\"),\n"},
		{"task type", repairAutoBanTaskTypeSource, "\tTypeTrashPurgeExpired     = \"trash:purge_expired\"\n"},
		{"cron start", repairAutoBanCronStartSource, "\ts, err := New(cfg.RedisURL)\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			once, fixed, warn := c.fn(c.src)
			if len(warn) > 0 {
				t.Fatalf("warned on the shape Grit writes: %v", warn)
			}
			if len(fixed) == 0 {
				t.Fatal("the repair did nothing to a project that needs it")
			}
			twice, fixed2, _ := c.fn(once)
			if twice != once || len(fixed2) > 0 {
				t.Error("running upgrade twice changes the file twice")
			}
		})
	}
}

// gofmt aligns struct fields into columns, so an anchor copied from the
// template matches the template and not one real project. lineWith is what
// makes the anchors survive that, and this is the case that caught it.
func TestTheConfigRepairSurvivesGofmtAlignment(t *testing.T) {
	aligned := `type Config struct {
	SentinelTurnstileSecret    string
	SentinelHCaptchaSecret     string
}

func Load() *Config {
	return &Config{
		SentinelTurnstileSecret:    getEnv("SENTINEL_TURNSTILE_SECRET", ""),
	}
}
`
	out, fixed, warn := repairAutoBanConfigSource(aligned)
	if len(warn) > 0 {
		t.Fatalf("an aligned config was refused: %v", warn)
	}
	if len(fixed) != 1 {
		t.Fatalf("fixed %v", fixed)
	}
	for _, want := range []string{
		"SecurityAutoBanEnabled    bool",
		`getEnvInt("SECURITY_AUTO_BAN_MIN_THREATS", 3)`,
		`getEnvDuration("AUTO_BAN_INTERVAL_HOURS", 6, time.Hour)`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired config is missing %q", want)
		}
	}
}

// TestTheConfigRepairWorksWithoutTheTurnstileField is the case the WhatsApp
// blueprint found.
//
// The repair anchored on SentinelTurnstileSecret, a field added to config.go
// after some projects were scaffolded. In one of those, the anchor is missing,
// the repair warns and returns, and the upgrade carries on to write
// internal/services/security_autoban.go and internal/cron/start.go, both of
// which read cfg.SecurityAutoBan*. The upgrade reported success and left a
// project that would not compile:
//
//	cfg.SecurityAutoBanEnabled undefined (type *config.Config has no field
//	or method SecurityAutoBanEnabled)
//
// So the fallback is the end of the struct and the end of the literal, which
// are there whatever the fields have done since.
func TestTheConfigRepairWorksWithoutTheTurnstileField(t *testing.T) {
	// A config.go from before SentinelTurnstileSecret existed.
	old := `type Config struct {
	AppPort string
	AppURL  string

	// Serve the API reference at /docs in production (API_DOCS_PUBLIC).
	APIDocsPublic bool
}

func Load() (*Config, error) {
	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		AppURL:        getEnv("APP_URL", "http://localhost:8080"),
		APIDocsPublic: getEnv("API_DOCS_PUBLIC", "false") == "true",
	}
	return cfg, nil
}
`
	out, fixed, warn := repairAutoBanConfigSource(old)
	if len(warn) > 0 {
		t.Fatalf("a config.go with no Turnstile field was refused: %v", warn)
	}
	if len(fixed) != 1 {
		t.Fatalf("fixed %v", fixed)
	}
	for _, want := range []string{
		"SecurityAutoBanEnabled    bool",
		`SecurityAutoBanEnabled:    getEnv("SECURITY_AUTO_BAN", "true") == "true",`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired config is missing %q\n\n%s", want, out)
		}
	}
	// Inside the struct and inside the literal, not after either of them.
	if strings.Index(out, "SecurityAutoBanEnabled    bool") > strings.Index(out, "func Load()") {
		t.Error("the fields landed outside the Config struct")
	}
	if !strings.Contains(out, "== \"true\",\n\t}\n\treturn cfg, nil") {
		t.Errorf("the loads did not land as the last entries of the literal\n\n%s", out)
	}

	// And it stays idempotent, which is what a second upgrade does.
	twice, fixed2, _ := repairAutoBanConfigSource(out)
	if twice != out || len(fixed2) > 0 {
		t.Error("running upgrade twice changes the file twice")
	}
}

// TestTheWorkerRepairIsAllOrNothing is the second half of what the WhatsApp
// blueprint found.
//
// The repair used to collect warnings and return the source it had already
// modified, so a project missing one anchor got the other edits written: the
// WorkerDeps field without the services import (undefined: services), and the
// handler without the mux line. Because its only guard was the string the mux
// line carries, the next upgrade added the handler again:
//
//	workers.go:305:6: handleSecurityAutoBan redeclared in this block
func TestTheWorkerRepairIsAllOrNothing(t *testing.T) {
	// Grit's WorkerDeps, and a mux this project has made its own.
	partial := `package jobs

import (
	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"app/apps/api/internal/cache"
	"app/apps/api/internal/mail"
	"app/apps/api/internal/storage"
)

type WorkerDeps struct {
	DB      *gorm.DB
	Mailer  *mail.Mailer
	Storage *storage.Storage
	Cache   *cache.Cache
}

func Register(mux *asynq.ServeMux, deps WorkerDeps) {
	mux.HandleFunc(TypeEmailSend, handleEmailSend(deps))
}

// handleBackupWeekly takes the scheduled full-database backup
func handleBackupWeekly(deps WorkerDeps) {}
`
	out, fixed, warn := repairAutoBanWorkerSource(partial)
	if len(warn) == 0 {
		t.Fatal("a file with no TypeTrashPurgeExpired line was repaired silently")
	}
	if len(fixed) != 0 {
		t.Errorf("reported %v while warning", fixed)
	}
	if out != partial {
		t.Error("the file was modified even though an anchor was missing, which is the half-applied " +
			"state that does not compile")
	}

	// Run it again, as the next upgrade would. Still unchanged, still warning,
	// and above all not a second copy of anything.
	twice, _, _ := repairAutoBanWorkerSource(out)
	if twice != partial {
		t.Error("a second upgrade changed the file again")
	}
}

// TestTheWorkerRepairImportsWhatItNames: the field it adds has a type from
// internal/services, and a project whose workers.go never imported that
// package stopped compiling with "undefined: services".
func TestTheWorkerRepairImportsWhatItNames(t *testing.T) {
	src := `package jobs

import (
	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"whatsapp/apps/api/internal/cache"
	"whatsapp/apps/api/internal/mail"
	"whatsapp/apps/api/internal/storage"
)

type WorkerDeps struct {
	DB      *gorm.DB
	Mailer  *mail.Mailer
	Storage *storage.Storage
	Cache   *cache.Cache
}

func Register(mux *asynq.ServeMux, deps WorkerDeps) {
	mux.HandleFunc(TypeTrashPurgeExpired, handleTrashPurgeExpired(deps))
}

// handleBackupWeekly takes the scheduled full-database backup
func handleBackupWeekly(deps WorkerDeps) {}
`
	out, fixed, warn := repairAutoBanWorkerSource(src)
	if len(warn) > 0 {
		t.Fatalf("refused a file with every anchor: %v", warn)
	}
	if len(fixed) != 1 {
		t.Fatalf("fixed %v", fixed)
	}
	if !strings.Contains(out, "AutoBan *services.SecurityAutoBan") {
		t.Error("the WorkerDeps field was not added")
	}
	// The module comes off a sibling import, so it is this project's and not
	// a guess.
	if !strings.Contains(out, `"whatsapp/apps/api/internal/services"`) {
		t.Errorf("the services import was not added, so this file is `undefined: services`\n\n%s", out)
	}
	if strings.Contains(out, `"app/apps/api/internal/services"`) {
		t.Error("the import was written with a module path from somewhere else")
	}

	// Idempotent: the second upgrade does nothing at all.
	twice, fixed2, _ := repairAutoBanWorkerSource(out)
	if twice != out || len(fixed2) > 0 {
		t.Error("running upgrade twice changes the file twice")
	}
}

// TestTheCronArgumentIsFixedEvenWhenTheWorkerIsNot.
//
// main.go has two edits and only one of them depends on workers.go. Gating
// both on the worker skipped the cron.New argument, and cron.go is
// framework-owned: its signature changes on every upgrade whatever the worker
// did, so the call site has to follow.
//
//	cmd/server/main.go:267:23: not enough arguments in call to cron.New
//	        have (string)
//	        want (string, "time".Duration)
func TestTheCronArgumentIsFixedEvenWhenTheWorkerIsNot(t *testing.T) {
	src := `package main

func main() {
	stop, err := jobs.StartWorker(cfg.RedisURL, jobs.WorkerDeps{
		DB:      db,
		Mailer:  mailer,
		Storage: storageService,
		Cache:   cacheService,
	})
	cs, err := cron.New(cfg.RedisURL)
}
`
	out, fixed, warn := repairAutoBanMainWith(src, false)
	if len(warn) > 0 {
		t.Fatalf("warned although the worker was deliberately skipped: %v", warn)
	}
	if !strings.Contains(out, "cron.New(cfg.RedisURL, cfg.SecurityAutoBanEvery)") {
		t.Error("the cron.New argument was not added, so main.go does not compile")
	}
	if strings.Contains(out, "AutoBan: services.NewSecurityAutoBan") {
		t.Error("the WorkerDeps field was set although workers.go has no such field")
	}
	if len(fixed) != 1 {
		t.Errorf("reported %v, wanted just the scheduler", fixed)
	}

	// And with the worker wired, both land.
	both, fixed2, warn2 := repairAutoBanMainWith(src, true)
	if len(warn2) > 0 {
		t.Fatalf("refused the shape Grit writes: %v", warn2)
	}
	if !strings.Contains(both, "AutoBan: services.NewSecurityAutoBan") ||
		!strings.Contains(both, "cfg.SecurityAutoBanEvery") {
		t.Error("a wired project should get both edits")
	}
	if len(fixed2) != 2 {
		t.Errorf("reported %v, wanted both", fixed2)
	}
}
