package health

import (
	"encoding/json"
	"testing"
)

func TestOverallIgnoresOffAndUnknown(t *testing.T) {
	// The whole point of four states: a deployment with no Redis and an
	// unsampled queue is healthy, and said so without anybody having to
	// remember to exclude those two cases at the call site.
	comps := map[string]Report{
		"database": OK("", nil),
		"redis":    Off("no REDIS_URL"),
		"jobs":     Unknown("no snapshot yet"),
	}
	if got := Overall(comps); got != StateOK {
		t.Fatalf("Overall = %q, want ok: off and unknown must not lower it", got)
	}
}

func TestOverallIsDegradedWhenAnyComponentIs(t *testing.T) {
	comps := map[string]Report{
		"database": OK("", nil),
		"redis":    Degraded("ping failed", nil),
		"email":    Off("no mailer"),
	}
	if got := Overall(comps); got != StateDegraded {
		t.Fatalf("Overall = %q, want degraded", got)
	}
}

func TestOverallWithNoComponentsIsOK(t *testing.T) {
	if got := Overall(nil); got != StateOK {
		t.Fatalf("Overall(nil) = %q, want ok", got)
	}
}

func TestReportFlattensItsFields(t *testing.T) {
	// A reader sees {"state":"ok","latency_ms":3}, not a nested object, because
	// the clients that read latency_ms predate this package.
	raw, err := json.Marshal(OK("", map[string]any{"latency_ms": 3, "tables": 41}))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["state"] != "ok" || got["ok"] != true || got["latency_ms"] != float64(3) || got["tables"] != float64(41) {
		t.Fatalf("flattened report = %v", got)
	}
	if _, ok := got["detail"]; ok {
		t.Fatalf("an empty detail should be left out, got %v", got)
	}
}

func TestLegacyOKFlagIsTrueOnlyForOK(t *testing.T) {
	// A load balancer or an old dashboard that reads only "ok" must see what it
	// saw before this package existed, which is false for anything but ok.
	for _, tc := range []struct {
		name string
		r    Report
		want bool
	}{
		{"ok", OK("", nil), true},
		{"degraded", Degraded("down", nil), false},
		{"off", Off("not configured"), false},
		{"unknown", Unknown("not sampled"), false},
		{"zero value", Report{}, false},
	} {
		raw, err := json.Marshal(tc.r)
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.name, err)
		}
		var got map[string]any
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("%s: unmarshal: %v", tc.name, err)
		}
		if got["ok"] != tc.want {
			t.Errorf("%s: ok = %v, want %v", tc.name, got["ok"], tc.want)
		}
	}
}

func TestZeroReportIsUnknownNotOK(t *testing.T) {
	// Report{} is what a struct literal somebody forgot to fill looks like, and
	// the safe reading of it is "nobody said", not "everything is fine".
	raw, err := json.Marshal(Report{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["state"] != "unknown" {
		t.Fatalf("zero Report state = %v, want unknown", got["state"])
	}
}

func TestRegisterAndSnapshot(t *testing.T) {
	Reset()
	defer Reset()

	Register("billing", func() Report { return OK("", map[string]any{"settled": 2}) })
	Register("search", func() Report { return Off("no SEARCH_URL") })

	if names := Registered(); len(names) != 2 || names[0] != "billing" || names[1] != "search" {
		t.Fatalf("Registered = %v", names)
	}
	snap := Snapshot()
	if snap["billing"].State != StateOK || snap["search"].State != StateOff {
		t.Fatalf("Snapshot = %v", snap)
	}
	if snap["billing"].Fields["settled"] != 2 {
		t.Fatalf("a probe's fields should survive the snapshot, got %v", snap["billing"].Fields)
	}
}

func TestRegisterTwiceReplaces(t *testing.T) {
	Reset()
	defer Reset()

	Register("cache", func() Report { return Off("first") })
	Register("cache", func() Report { return OK("second", nil) })

	if names := Registered(); len(names) != 1 {
		t.Fatalf("one name expected, got %v", names)
	}
	if got := Snapshot()["cache"]; got.State != StateOK || got.Detail != "second" {
		t.Fatalf("second registration should win, got %+v", got)
	}
}

func TestAPanickingProbeDoesNotTakeTheEndpointDown(t *testing.T) {
	Reset()
	defer Reset()

	Register("bad", func() Report { panic("a nil map write, say") })
	Register("database", func() Report { return OK("", nil) })

	snap := Snapshot()
	if snap["bad"].State != StateUnknown {
		t.Fatalf("a panicking probe should be unknown, got %+v", snap["bad"])
	}
	if snap["database"].State != StateOK {
		t.Fatalf("the other components must still be reported, got %+v", snap["database"])
	}
	// And it must not make the deployment look broken: the probe failed, not
	// the component.
	if got := Overall(snap); got != StateOK {
		t.Fatalf("Overall = %q, want ok", got)
	}
}

func TestRegisterIgnoresNonsense(t *testing.T) {
	Reset()
	defer Reset()

	Register("", func() Report { return OK("", nil) })
	Register("nil", nil)

	if names := Registered(); len(names) != 0 {
		t.Fatalf("Registered = %v, want none", names)
	}
}
