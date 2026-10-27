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
