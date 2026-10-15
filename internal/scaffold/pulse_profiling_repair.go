package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Pulse v1.3.0 added the runtime's own profiles and the switches that turn them
// on. The mount lives in internal/routes/routes.go and the settings in
// internal/config/config.go, and neither file is framework-owned on upgrade:
// both are files a project edits. So a new project would get profiling and
// every existing project would not, which is the half-delivery that this
// repair exists to prevent.
//
// The text below is spliced into the scaffold templates as well, so there is
// one copy of it rather than two that drift.
const (
	configPulseProfilingFields = "\n" +
		"\t// The config half of Pulse's two-gate profiling switch. Pulse also needs\n" +
		"\t// PULSE_PROFILE_ENABLED set in the environment before it will sample, so\n" +
		"\t// one forgotten setting cannot expose the shape of your code.\n" +
		"\tPulseProfiling bool\n" +
		"\t// Mutex and block profiling cost the runtime something while on, so both\n" +
		"\t// stay off until a number is given. Zero means off.\n" +
		"\tPulseMutexProfileFraction int\n" +
		"\tPulseBlockProfileRate     int\n" +
		"\t// Ask the database to explain a slow SELECT. Off by default: it issues\n" +
		"\t// extra read-only queries against your database.\n" +
		"\tPulseExplainSlowQueries bool\n"

	configPulseProfilingLoad = "\n" +
		"\t\tPulseProfiling:            getEnv(\"PULSE_PROFILING\", \"false\") == \"true\",\n" +
		"\t\tPulseMutexProfileFraction: getEnvInt(\"PULSE_MUTEX_PROFILE_FRACTION\", 0),\n" +
		"\t\tPulseBlockProfileRate:     getEnvInt(\"PULSE_BLOCK_PROFILE_RATE\", 0),\n" +
		"\t\tPulseExplainSlowQueries:   getEnv(\"PULSE_EXPLAIN_SLOW_QUERIES\", \"false\") == \"true\",\n"

	routesPulseProfilingBlock = "\n" +
		"\t\t// The runtime's own profiles, from v1.3.0: heap, allocations, mutex,\n" +
		"\t\t// block and CPU, served as standard pprof protobuf so\n" +
		"\t\t// \"go tool pprof http://host/pulse/api/profile/heap\" works and\n" +
		"\t\t// -diff_base compares two captures.\n" +
		"\t\t//\n" +
		"\t\t// Two switches, not one, and deliberately so. This one is the config\n" +
		"\t\t// half; Pulse also refuses to sample unless PULSE_PROFILE_ENABLED is\n" +
		"\t\t// set in the environment. A profile leaks the shape of your code, so\n" +
		"\t\t// one forgotten setting should not be enough to expose it.\n" +
		"\t\tif cfg.PulseProfiling {\n" +
		"\t\t\tpulseOpts = append(pulseOpts, pulse.WithProfiling(), pulse.WithPprofHandlers())\n" +
		"\t\t}\n" +
		"\t\t// Mutex and block profiling cost the runtime something on every\n" +
		"\t\t// contended lock, so the runtime keeps them off and so do we. Set a\n" +
		"\t\t// fraction or a rate only while you are looking: these answer \"what is\n" +
		"\t\t// contended\", which nothing else here can.\n" +
		"\t\tif cfg.PulseMutexProfileFraction > 0 {\n" +
		"\t\t\tpulseOpts = append(pulseOpts, pulse.WithMutexProfiling(cfg.PulseMutexProfileFraction))\n" +
		"\t\t}\n" +
		"\t\tif cfg.PulseBlockProfileRate > 0 {\n" +
		"\t\t\tpulseOpts = append(pulseOpts, pulse.WithBlockProfiling(cfg.PulseBlockProfileRate))\n" +
		"\t\t}\n" +
		"\t\t// Ask the database for the plan behind a slow SELECT. Off by default\n" +
		"\t\t// because it issues extra queries against your database; when on it is\n" +
		"\t\t// EXPLAIN only, never EXPLAIN ANALYZE, so the slow query is not run a\n" +
		"\t\t// second time.\n" +
		"\t\tif cfg.PulseExplainSlowQueries {\n" +
		"\t\t\tpulseOpts = append(pulseOpts, pulse.WithExplainSlowQueries(0))\n" +
		"\t\t}\n"

	envPulseProfiling = "\n" +
		"# Pulse's runtime profiles: heap, allocations, mutex, block and CPU, served as\n" +
		"# standard pprof so \"go tool pprof http://host/pulse/api/profile/heap\" works.\n" +
		"#\n" +
		"# Two switches on purpose. This one, and PULSE_PROFILE_ENABLED in the\n" +
		"# environment where the server runs. A profile leaks the shape of your code, so\n" +
		"# one forgotten setting should not be enough to expose it.\n" +
		"PULSE_PROFILING=false\n" +
		"\n" +
		"# What is contended, which nothing else here can answer. Both cost the runtime\n" +
		"# a little while on, so both stay off until you give a number and are worth\n" +
		"# setting only while you are looking. 100 and 10000 are common starting points.\n" +
		"PULSE_MUTEX_PROFILE_FRACTION=0\n" +
		"PULSE_BLOCK_PROFILE_RATE=0\n" +
		"\n" +
		"# Ask the database for the plan behind a slow SELECT and show it on the\n" +
		"# Database page. Off by default because it issues extra read-only queries\n" +
		"# against your database. EXPLAIN only, never EXPLAIN ANALYZE, so the slow\n" +
		"# query is never run a second time.\n" +
		"PULSE_EXPLAIN_SLOW_QUERIES=false\n"
)

var (
	configPulseDSNFieldRe = regexp.MustCompile(`(?m)^\tPulseStorageDSN string // path for sqlite[^\n]*\n`)
	configPulseDSNLoadRe  = regexp.MustCompile(`(?m)^\t\tPulseStorageDSN: resolveSQLitePath\(getEnv\("PULSE_STORAGE_DSN"[^\n]*\n`)
	routesPulseMountRe    = regexp.MustCompile(`(?m)^\t\tp := pulse\.Mount\(context\.Background\(\), r, db, pulseOpts\.\.\.\)\n`)
	envPulsePasswordRe    = regexp.MustCompile(`(?m)^PULSE_PASSWORD=[^\n]*\n`)
	envPulseProfilingRe   = regexp.MustCompile(`(?m)^(# )?PULSE_PROFILING=`)
)

// repairPulseProfiling brings the Pulse v1.3.0 switches to an existing project:
// the four settings in config.go, the options on the mount in routes.go, and
// the settings documented in .env.example.
func repairPulseProfiling(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	config := filepath.Join(apiRoot, "internal", "config", "config.go")
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(config) || !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if err := repairSourceFile(root, m, config, repairPulseProfilingConfigSource); err != nil {
		return err
	}
	// The mount reads the fields, so it is only changed once config.go has
	// them. A config.go too edited to patch therefore leaves a project that
	// still builds, rather than one referring to fields that do not exist.
	if !fileContains(config, "PulseExplainSlowQueries") {
		return nil
	}
	if err := repairSourceFile(root, m, routes, repairPulseProfilingRoutesSource); err != nil {
		return err
	}
	if env := filepath.Join(root, ".env.example"); fileExists(env) {
		if err := repairTextFile(root, m, env, repairPulseProfilingEnvSource); err != nil {
			return err
		}
	}
	return nil
}

func repairPulseProfilingConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "PulseExplainSlowQueries") {
		return src, nil, nil
	}
	if len(configPulseDSNFieldRe.FindAllStringIndex(src, -1)) != 1 ||
		len(configPulseDSNLoadRe.FindAllStringIndex(src, -1)) != 1 {
		return src, nil, []string{"config.go is not the file Grit wrote: add PulseProfiling, PulseMutexProfileFraction, PulseBlockProfileRate and PulseExplainSlowQueries read from PULSE_PROFILING, PULSE_MUTEX_PROFILE_FRACTION, PULSE_BLOCK_PROFILE_RATE and PULSE_EXPLAIN_SLOW_QUERIES, or the Pulse profile endpoints stay unreachable"}
	}
	out := configPulseDSNFieldRe.ReplaceAllStringFunc(src, func(line string) string {
		return line + configPulseProfilingFields
	})
	out = configPulseDSNLoadRe.ReplaceAllStringFunc(out, func(line string) string {
		return line + configPulseProfilingLoad
	})
	return out, []string{"PULSE_PROFILING, the two contention profiles and PULSE_EXPLAIN_SLOW_QUERIES are read"}, nil
}

func repairPulseProfilingRoutesSource(src string) (string, []string, []string) {
	if strings.Contains(src, "pulse.WithPprofHandlers()") {
		return src, nil, nil
	}
	if len(routesPulseMountRe.FindAllStringIndex(src, -1)) != 1 {
		return src, nil, []string{"routes.go is not the file Grit wrote: pass pulse.WithProfiling() and pulse.WithPprofHandlers() when cfg.PulseProfiling is set, or heap, allocation and CPU profiles cannot be fetched from a running process"}
	}
	out := routesPulseMountRe.ReplaceAllStringFunc(src, func(line string) string {
		return routesPulseProfilingBlock + line
	})
	return out, []string{"heap, allocation, mutex, block and CPU profiles are served when PULSE_PROFILING is on"}, nil
}

func repairPulseProfilingEnvSource(src string) (string, []string, []string) {
	if envPulseProfilingRe.MatchString(src) {
		return src, nil, nil
	}
	fixed := []string{"the Pulse profiling settings are documented"}
	// After PULSE_PASSWORD, so the settings sit with the rest of Pulse. There
	// are two .env templates and the single-file shape has its own, so the
	// anchor is matched rather than assumed, and a file without it gets the
	// block appended instead of nothing.
	loc := envPulsePasswordRe.FindStringIndex(src)
	if loc == nil {
		return strings.TrimRight(src, "\n") + "\n" + envPulseProfiling, fixed, nil
	}
	return src[:loc[1]] + envPulseProfiling + src[loc[1]:], fixed, nil
}
