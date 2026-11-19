package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ─── /healthz and /readyz on a project that predates them ───────────────────
//
// internal/routes/probes_routes.go travels with writeFrameworkOwnedFiles, so an
// upgraded project gets the file. Nothing calls it: routes.go is the
// developer's, held back by the manifest guard the moment anybody edits it, so
// the one line that mounts the probes has to be put there by a repair.
//
// A file with a function nobody calls still compiles, which is exactly why this
// is worth a repair rather than a release note. The symptom would be a project
// that upgraded, got the code, and answered 404 on /readyz, with nothing
// anywhere saying why.
//
// The anchor is the pair of health registrations, which every project has had
// since M38 and which repairHealthVersionedRoute writes for anything older. It
// runs after that repair for exactly that reason.

// probesMountOld is the anchor: the two /api/health registrations, together.
//
// Both lines rather than one, because the insert has to land after the second
// one and a single-line anchor would put it between them.
const probesMountOld = "\tr.GET(\"/api/health\", healthCheck)\n\tr.GET(\"/api/\"+APIVersion+\"/health\", healthCheck)\n"

// probesMountNew is the same two lines with the mount after them. Kept here as
// a literal rather than built from the template's text: a repair that reads the
// template is a repair whose behaviour changes when somebody reformats a
// comment.
const probesMountNew = probesMountOld + `
	// /healthz and /readyz, for the orchestrator rather than for the System
	// Health page. See internal/routes/probes_routes.go for why liveness probes
	// nothing and readiness gates on the database alone.
	mountProbes(r, db)
`

// repairProbesMount makes an existing project serve /healthz and /readyz.
func repairProbesMount(root string, opts Options) error {
	routes := filepath.Join(opts.APIRoot(root), "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairProbesMountSource)
}

func repairProbesMountSource(src string) (string, []string, []string) {
	// Guarded on the call, not on a comment near it: a Contains check that
	// matches the prose explaining the mount is a repair that silently
	// no-ops for ever.
	if strings.Contains(src, "mountProbes(r, db)") {
		return src, nil, nil
	}
	if strings.Count(src, probesMountOld) != 1 {
		return src, nil, []string{
			"routes.go: could not find the /api/health registrations, so /healthz " +
				"and /readyz are not mounted. Add mountProbes(r, db) beside them; " +
				"internal/routes/probes_routes.go is already there.",
		}
	}
	out := strings.Replace(src, probesMountOld, probesMountNew, 1)
	fixed := []string{
		"/healthz and /readyz are mounted: liveness that touches no dependency, " +
			"and readiness that gates on the database",
	}
	if excluded, ok := excludeProbesFromPulse(out); ok {
		out = excluded
		fixed = append(fixed,
			"the probes are excluded from Pulse, so a kubelet calling them every "+
				"few seconds does not fill the request log or flatten the latency charts")
	}
	return out, fixed, nil
}

// excludeProbesFromPulse keeps the probes out of Pulse's request log.
//
// A kubelet calls /healthz and /readyz every few seconds, so one replica
// contributes thousands of records a day. They would dominate the log, and a
// flood of 2ms requests drags the RED charts down until the latency of the
// thing people actually use is invisible inside it.
//
// Best effort, and silent when it does not apply. A project that has edited
// its Pulse options, or turned Pulse off, has made a decision this repair has
// no business overriding, and it is not a reason to warn: the probes work
// either way.
func excludeProbesFromPulse(src string) (string, bool) {
	const old = `pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),`
	if strings.Contains(src, `"/healthz", "/readyz"),`) || strings.Count(src, old) != 1 {
		return src, false
	}
	const updated = `pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*",
				"/healthz", "/readyz"),`
	return strings.Replace(src, old, updated, 1), true
}
