package scaffold

import (
	"strings"
	"testing"
)

func TestTemplatesServePulseProfiles(t *testing.T) {
	config := apiConfigGo()
	for _, want := range []string{
		"PulseProfiling bool",
		"PulseMutexProfileFraction int",
		"PulseBlockProfileRate     int",
		"PulseExplainSlowQueries bool",
		`getEnv("PULSE_PROFILING", "false") == "true"`,
		`getEnvInt("PULSE_MUTEX_PROFILE_FRACTION", 0)`,
		`getEnvInt("PULSE_BLOCK_PROFILE_RATE", 0)`,
		`getEnv("PULSE_EXPLAIN_SLOW_QUERIES", "false") == "true"`,
	} {
		if !strings.Contains(config, want) {
			t.Errorf("config.go lacks %q", want)
		}
	}

	routes := apiRoutesGo()
	for _, want := range []string{
		"pulse.WithProfiling(), pulse.WithPprofHandlers()",
		"pulse.WithMutexProfiling(cfg.PulseMutexProfileFraction)",
		"pulse.WithBlockProfiling(cfg.PulseBlockProfileRate)",
		"pulse.WithExplainSlowQueries(0)",
	} {
		if !strings.Contains(routes, want) {
			t.Errorf("routes.go lacks %q", want)
		}
	}
	// The options have to reach pulse.Mount, which they do not if they are
	// appended after it.
	if strings.Index(routes, "pulse.WithPprofHandlers()") > strings.Index(routes, "p := pulse.Mount(") {
		t.Error("the profiling options are appended after pulse.Mount")
	}
	// Each of the three costly ones is behind its own switch rather than on by
	// default: a profile anyone can ask for describes the program, and a 30s
	// CPU profile fetched repeatedly is a denial of service.
	for _, gate := range []string{
		"if cfg.PulseProfiling {",
		"if cfg.PulseMutexProfileFraction > 0 {",
		"if cfg.PulseBlockProfileRate > 0 {",
		"if cfg.PulseExplainSlowQueries {",
	} {
		if !strings.Contains(routes, gate) {
			t.Errorf("routes.go does not gate profiling on %q", gate)
		}
	}

	for _, env := range []string{envFile(Options{ProjectName: "app"}), envExampleFile(Options{ProjectName: "app"})} {
		for _, want := range []string{
			"PULSE_PROFILING=false",
			"PULSE_MUTEX_PROFILE_FRACTION=0",
			"PULSE_BLOCK_PROFILE_RATE=0",
			"PULSE_EXPLAIN_SLOW_QUERIES=false",
		} {
			if !strings.Contains(env, want) {
				t.Errorf("an env template lacks %q", want)
			}
		}
	}
}

func TestPulseProfilingRepairs(t *testing.T) {
	config := apiConfigGo()
	oldConfig := strings.Replace(config, configPulseProfilingFields, "", 1)
	oldConfig = strings.Replace(oldConfig, configPulseProfilingLoad, "", 1)
	checkRepair(t, "config.go", repairPulseProfilingConfigSource, oldConfig, config)

	routes := apiRoutesGo()
	checkRepair(t, "routes.go", repairPulseProfilingRoutesSource,
		strings.Replace(routes, routesPulseProfilingBlock, "", 1), routes)

	env := envFile(Options{ProjectName: "app"})
	checkRepair(t, ".env", repairPulseProfilingEnvSource, strings.Replace(env, envPulseProfiling, "", 1), env)

	// A config.go whose Pulse storage settings somebody rewrote is named, not
	// patched into something that does not compile.
	odd := strings.Replace(oldConfig,
		`	PulseStorageDSN string // path for sqlite, e.g. "pulse.db" or ":memory:"`,
		"\tPulseStorageDSN string // where the pulse database lives", 1)
	if out, _, warnings := repairPulseProfilingConfigSource(odd); out != odd || len(warnings) != 1 {
		t.Errorf("an unrecognised config.go was changed or not reported: %v", warnings)
	}

	// And a routes.go that mounts Pulse its own way keeps its mount.
	oddRoutes := strings.Replace(strings.Replace(routes, routesPulseProfilingBlock, "", 1),
		"\t\tp := pulse.Mount(context.Background(), r, db, pulseOpts...)",
		"\t\tp := mountPulseMyWay(r, db, pulseOpts)", 1)
	if out, _, warnings := repairPulseProfilingRoutesSource(oddRoutes); out != oddRoutes || len(warnings) != 1 {
		t.Errorf("an unrecognised routes.go was changed or not reported: %v", warnings)
	}

	// A .env with no PULSE_PASSWORD still gets the settings, appended, rather
	// than silently getting nothing.
	noAnchor := "APP_NAME=app\n"
	out, fixed, _ := repairPulseProfilingEnvSource(noAnchor)
	if len(fixed) != 1 || !strings.Contains(out, "PULSE_PROFILING=false") {
		t.Error("an .env without the Pulse block did not get the profiling settings")
	}
}
