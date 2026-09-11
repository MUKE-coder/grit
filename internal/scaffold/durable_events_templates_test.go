package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func durableMustContain(t *testing.T, what, src string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(src, want) {
			t.Errorf("%s is missing %q", what, want)
		}
	}
}

func durableMustFormat(t *testing.T, what, src string) {
	t.Helper()
	if _, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Fatalf("%s is not valid Go: %v", what, err)
	}
}

// Approving a purchase request has to draw down a budget and reserve stock, and
// refuse when it cannot. A transition was one UPDATE and an in-memory event, so
// that work had nowhere to go but an async subscriber that could neither refuse
// the move nor be relied on to run.
func TestWorkflowPackageHasTransitionHooks(t *testing.T) {
	src := apiWorkflowGo()
	durableMustFormat(t, "workflow.go", src)
	durableMustContain(t, "workflow.go", src,
		"func OnTransition(resource, action string, h Hook)",
		"func RunHooks(tx *gorm.DB, m Move) error",
		"func Refuse(format string, args ...interface{}) error",
		"func Classify(err error) (int, string)",
		`"TRANSITION_REFUSED"`,
		"return http.StatusInternalServerError, \"INTERNAL_ERROR\"",
	)
}

// Generated services depend on workflow.go and events/durable.go, so both have
// to travel on upgrade with the other packages generated code imports.
func TestDurableFilesTravelOnUpgrade(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchAPI}
	if err := writeCodegenRuntimeFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"internal/events/durable.go", "internal/workflow/workflow.go"} {
		if _, err := os.Stat(filepath.Join(opts.APIRoot(root), filepath.FromSlash(rel))); err != nil {
			t.Errorf("%s is not written by the upgrade path: %v", rel, err)
		}
	}
}

func TestEventsHaveDurableDelivery(t *testing.T) {
	events := apiEventsGo()
	durableMustFormat(t, "events.go", events)
	durableMustContain(t, "events.go", events,
		"\tDurable\n)",
		"b.enqueueDurable(e)",
		// A subscriber registered from init() ran before Init and was dropped.
		"pending = append(pending, subscription{",
		"for _, s := range pending {",
	)

	durable := apiEventsDurableGo()
	durableMustFormat(t, "durable.go", durable)
	durableMustContain(t, "durable.go", durable,
		"func StartRelay(db *gorm.DB)",
		"func EmitTx(tx *gorm.DB, c *gin.Context, e *Event) error",
		"outbox.Enqueue(tx, durableTopicPrefix+name, e)",
		"TopicPrefix: durableTopicPrefix",
		"func (e Event) DecodeAfter(v interface{}) error",
		// A project has no package-level database handle, so a subscriber
		// registered from init() had nothing to write with.
		"func OnDurable(pattern, name string, h func(tx *gorm.DB, e Event) error)",
	)

	relay := outboxRelayGo()
	durableMustFormat(t, "relay.go", relay)
	durableMustContain(t, "relay.go", relay, "TopicPrefix string", `q.Where("topic LIKE ?", r.TopicPrefix+"%")`)
	// The old claim chained .Or onto the status filter, which a topic filter
	// added after it would only have narrowed half of.
	if strings.Contains(relay, "Or(\"status = ? AND claimed_at") {
		t.Error("the claim still ORs the stale-claim half outside the group")
	}

	routes := apiRoutesGo()
	register := strings.Index(routes, "services.RegisterEventSubscribers(db, realtimeHub, nil)")
	start := strings.Index(routes, "events.StartRelay(db)")
	if register < 0 || start < register {
		t.Error("routes.go does not start the relay after registering the subscribers")
	}
}

func TestEnsureEventRelay(t *testing.T) {
	api := t.TempDir()
	routes := filepath.Join(api, "internal", "routes", "routes.go")
	durable := filepath.Join(api, "internal", "events", "durable.go")
	for _, dir := range []string{filepath.Dir(routes), filepath.Dir(durable)} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	src := "package routes\n\nfunc Setup() {\n\tevents.Init(4)\n\tservices.RegisterEventSubscribers(db, realtimeHub, nil)\n}\n"
	if err := os.WriteFile(routes, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	// Without durable.go the call would not compile, so routes.go is left alone.
	if changed, err := EnsureEventRelay(api); err != nil || changed {
		t.Fatalf("changed routes.go before durable.go existed: %v %v", changed, err)
	}

	if err := os.WriteFile(durable, []byte("package events\n\nfunc StartRelay(db any) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := EnsureEventRelay(api); err != nil {
			t.Fatal(err)
		}
	}
	got, _ := os.ReadFile(routes)
	if n := strings.Count(string(got), "events.StartRelay(db)"); n != 1 {
		t.Fatalf("StartRelay appears %d times:\n%s", n, got)
	}
	if !strings.Contains(string(got), "nil)\n\tevents.StartRelay(db)\n}") {
		t.Errorf("the relay is not started on the line after the subscribers:\n%s", got)
	}
}
