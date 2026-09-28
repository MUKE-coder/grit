package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The scale report answers "what is the one thing to do next". The stage panel
// is the other half of the question, and the one people ask first: what is
// already handled here? Every line of it is measured off the running
// deployment, so a green tick means this app was observed doing the thing.
func TestScaleReportCarriesTheStages(t *testing.T) {
	src := apiScaleHandlerGo()
	if _, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "demo/apps/api"))); err != nil {
		t.Fatalf("scale.go is not valid Go: %v", err)
	}
	for _, want := range []string{
		"type StageStatus struct",
		`Stages  []StageStatus`,
		"r.Stages = stages(r, h.largestTable())",
		"func stages(r ScaleReport, biggest seqScan) []StageStatus {",
		"func (h *ScaleHandler) largestTable() seqScan {",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the scale handler is missing %q", want)
		}
	}
	// All ten, in order, or the panel has a hole in it and nobody notices.
	for n, name := range map[string]string{
		"0": "Measure first",
		"1": "One server, one database",
		"2": "Vertical scaling",
		"3": "Horizontal + load balancer",
		"4": "Stateless servers",
		"5": "Connection pooling",
		"6": "Indexes, then read replicas",
		"7": "Caching",
		"8": "Queues and background jobs",
		"9": "Sharding",
	} {
		if !strings.Contains(src, `add("`+n+`", "`+name+`"`) {
			t.Errorf("stage %s (%s) is missing", n, name)
		}
	}
}

// The panel and the verdict read the same measurements, so a stage the verdict
// is calling out must not be showing a green tick beside it. Shared constants
// are how that is kept true, rather than two copies of 0.8 drifting apart.
func TestScaleStageThresholdsAreShared(t *testing.T) {
	src := apiScaleHandlerGo()
	for _, want := range []string{
		"connectionPressure = 0.8",
		"replicaLagSeconds  = 5.0",
		"poorHitRate        = 0.5",
		"hitRateSample      = 1000",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the shared threshold %q is missing", want)
		}
	}
	// The verdict has to use them too, or they are decoration.
	for _, want := range []string{
		"used > connectionPressure",
		"db.ReplicationLag > replicaLagSeconds",
		"r.Cache.HitRate < poorHitRate",
		"> hitRateSample",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the verdict does not use %q, so the panel can disagree with it", want)
		}
	}
}

// Two things can be genuinely wrong at these stages and both are invisible
// until a second instance exists, so the panel says so rather than ticking.
func TestScaleStagesAreHonestAboutSQLiteAndLocalDisk(t *testing.T) {
	src := apiScaleHandlerGo()
	if !strings.Contains(src, `stage1 = "attention"`) || !strings.Contains(src, "SQLite serialises writes") {
		t.Error("SQLite is not called out at stage 1")
	}
	if !strings.Contains(src, `stage4 = "attention"`) || !strings.Contains(src, "STORAGE_DRIVER=local") {
		t.Error("uploads on a local disk are not called out at stage 4")
	}
	// Stage 9 is the one Grit deliberately does nothing about, and says so.
	if !strings.Contains(src, `add("9", "Sharding", "yours"`) {
		t.Error("stage 9 does not admit that Grit does nothing here")
	}
}
