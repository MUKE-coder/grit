package saga

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"library/apps/api/internal/models"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_pragma=busy_timeout(5000)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&models.SagaRun{}, &models.SagaStep{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM saga_steps")
		db.Exec("DELETE FROM saga_runs")
	})
	return db
}

// drain advances a run until it stops moving, with no backoff waiting. The
// runner steps once per tick on purpose, so a test that wants a finished run
// has to tick.
func drain(t *testing.T, db *gorm.DB, id string) *models.SagaRun {
	t.Helper()
	runner := &Runner{DB: db, Batch: 10, ClaimTTL: time.Minute, Owner: "test"}
	for i := 0; i < 200; i++ {
		// Backoff would have the test sleeping, so the wait is removed rather
		// than waited out. What is being tested is the sequence, not the clock.
		db.Model(&models.SagaRun{}).Where("id = ?", id).Update("available_at", time.Now().Add(-time.Hour))
		if _, err := runner.Tick(context.Background()); err != nil {
			t.Fatalf("tick: %v", err)
		}
		var row models.SagaRun
		if err := db.Preload("Steps").First(&row, "id = ?", id).Error; err != nil {
			t.Fatalf("reload: %v", err)
		}
		if row.Finished() {
			return &row
		}
	}
	t.Fatal("the run never finished")
	return nil
}

func stepsByIdx(run *models.SagaRun) map[int]models.SagaStep {
	out := map[int]models.SagaStep{}
	for _, s := range run.Steps {
		out[s.Idx] = s
	}
	return out
}

func TestAHappyRunDoesEveryStepInOrder(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	var order []string
	var mu sync.Mutex
	record := func(name string) func(context.Context, *Run) error {
		return func(_ context.Context, r *Run) error {
			mu.Lock()
			defer mu.Unlock()
			order = append(order, name)
			r.Set(name, "done")
			return nil
		}
	}
	Register(Definition{Name: "checkout", Steps: []Step{
		{Name: "charge", Do: record("charge"), Undo: record("refund")},
		{Name: "reserve", Do: record("reserve"), Undo: record("release")},
		{Name: "ship", Do: record("ship")},
	}})

	started, err := Start(context.Background(), db, "checkout", map[string]any{"order": 7})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	run := drain(t, db, started.ID)

	if run.Status != models.SagaDone {
		t.Fatalf("status = %q, want done: %s", run.Status, run.LastError)
	}
	if len(order) != 3 || order[0] != "charge" || order[1] != "reserve" || order[2] != "ship" {
		t.Fatalf("ran %v", order)
	}
	for idx, step := range stepsByIdx(run) {
		if step.Status != models.SagaStepDone {
			t.Errorf("step %d (%s) = %q, want done", idx, step.Name, step.Status)
		}
	}
	if run.FinishedAt == nil {
		t.Error("a finished run has no finished_at")
	}
}

func TestAFailureUndoesTheCompletedStepsNewestFirst(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	var order []string
	note := func(name string) func(context.Context, *Run) error {
		return func(_ context.Context, r *Run) error {
			order = append(order, name)
			return nil
		}
	}
	Register(Definition{Name: "checkout", MaxAttempts: 1, Steps: []Step{
		{Name: "charge", Do: note("charge"), Undo: note("refund")},
		{Name: "reserve", Do: note("reserve"), Undo: note("release")},
		{Name: "ship", Do: func(context.Context, *Run) error { return errors.New("no courier") }, Undo: note("unship")},
	}})

	started, _ := Start(context.Background(), db, "checkout", nil)
	run := drain(t, db, started.ID)

	if run.Status != models.SagaCompensated {
		t.Fatalf("status = %q, want compensated", run.Status)
	}
	want := []string{"charge", "reserve", "release", "refund"}
	if len(order) != len(want) {
		t.Fatalf("ran %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("ran %v, want %v", order, want)
		}
	}
	// The step that failed is not compensated: it never completed, so there is
	// nothing of it to take back. Undoing it would be undoing something that
	// did not happen.
	steps := stepsByIdx(run)
	if steps[2].Status != models.SagaStepFailed {
		t.Errorf("the failed step is %q, want failed", steps[2].Status)
	}
	if steps[0].Status != models.SagaStepCompensated || steps[1].Status != models.SagaStepCompensated {
		t.Errorf("the completed steps are %q and %q, want compensated", steps[0].Status, steps[1].Status)
	}
}

func TestAStepIsRetriedBeforeItGivesUp(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	attempts := 0
	Register(Definition{Name: "flaky", MaxAttempts: 3, BaseBackoff: time.Millisecond, Steps: []Step{
		{Name: "call", Do: func(context.Context, *Run) error {
			attempts++
			if attempts < 3 {
				return errors.New("timeout")
			}
			return nil
		}},
	}})

	started, _ := Start(context.Background(), db, "flaky", nil)
	run := drain(t, db, started.ID)

	if run.Status != models.SagaDone {
		t.Fatalf("status = %q, want done", run.Status)
	}
	if attempts != 3 {
		t.Errorf("attempted %d times, want 3", attempts)
	}
}

func TestAFatalErrorSkipsTheRetries(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	attempts := 0
	Register(Definition{Name: "declined", MaxAttempts: 10, BaseBackoff: time.Millisecond, Steps: []Step{
		{Name: "charge", Do: func(context.Context, *Run) error {
			attempts++
			return Fatal(errors.New("card declined"))
		}},
	}})

	started, _ := Start(context.Background(), db, "declined", nil)
	run := drain(t, db, started.ID)

	if run.Status != models.SagaCompensated {
		t.Fatalf("status = %q, want compensated", run.Status)
	}
	// A declined card is not going to be accepted on the fourth try, and
	// spending ten attempts on it delays the compensation by an hour.
	if attempts != 1 {
		t.Errorf("attempted %d times, want 1", attempts)
	}
}

func TestACompensationThatCannotCompleteLeavesTheRunStuck(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	Register(Definition{Name: "stuck", MaxAttempts: 2, BaseBackoff: time.Millisecond, Steps: []Step{
		{
			Name: "charge",
			Do:   func(context.Context, *Run) error { return nil },
			Undo: func(context.Context, *Run) error { return errors.New("the refund API is down") },
		},
		{Name: "ship", Do: func(context.Context, *Run) error { return Fatal(errors.New("no courier")) }},
	}})

	started, _ := Start(context.Background(), db, "stuck", nil)
	run := drain(t, db, started.ID)

	// Not "compensated". Something happened that could not be taken back, and
	// the run says so rather than carrying on and leaving the charge in place
	// with a status that reads as resolved.
	if run.Status != models.SagaStuck {
		t.Fatalf("status = %q, want stuck", run.Status)
	}
	if steps := stepsByIdx(run); steps[0].Status != models.SagaStepStuck {
		t.Errorf("the uncompensated step is %q, want stuck", steps[0].Status)
	}
	if run.LastError == "" {
		t.Error("a stuck run should carry the error a person has to act on")
	}
}

func TestAStepWithNoUndoIsSkippedOnTheWayBack(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	refunded := false
	Register(Definition{Name: "email", MaxAttempts: 1, Steps: []Step{
		{
			Name: "charge",
			Do:   func(context.Context, *Run) error { return nil },
			Undo: func(context.Context, *Run) error { refunded = true; return nil },
		},
		// An email cannot be unsent. Declaring no Undo is allowed, and the
		// consequence is the ordering: a step after this one that fails leaves
		// the email sent.
		{Name: "notify", Do: func(context.Context, *Run) error { return nil }},
		{Name: "ship", Do: func(context.Context, *Run) error { return Fatal(errors.New("nope")) }},
	}})

	started, _ := Start(context.Background(), db, "email", nil)
	run := drain(t, db, started.ID)

	if run.Status != models.SagaCompensated {
		t.Fatalf("status = %q, want compensated", run.Status)
	}
	if !refunded {
		t.Error("the compensable step was not undone")
	}
}

func TestStateSurvivesBetweenStepsAndIntoACompensation(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	var refundedCharge string
	Register(Definition{Name: "state", MaxAttempts: 1, Steps: []Step{
		{
			Name: "charge",
			Do:   func(_ context.Context, r *Run) error { r.Set("charge_id", "ch_123"); return nil },
			Undo: func(_ context.Context, r *Run) error { refundedCharge = r.GetString("charge_id"); return nil },
		},
		{
			Name: "reserve",
			Do: func(_ context.Context, r *Run) error {
				if r.GetString("charge_id") != "ch_123" {
					return errors.New("the charge id did not reach the next step")
				}
				return Fatal(errors.New("out of stock"))
			},
		},
	}})

	started, _ := Start(context.Background(), db, "state", nil)
	run := drain(t, db, started.ID)

	if run.Status != models.SagaCompensated {
		t.Fatalf("status = %q, want compensated: %s", run.Status, run.LastError)
	}
	// The whole reason State is a column: the process that charged the card is
	// usually gone by the time the refund runs, so the id has to come from the
	// row and not from a closure.
	if refundedCharge != "ch_123" {
		t.Errorf("the refund got charge id %q, want ch_123", refundedCharge)
	}
}

func TestStateWrittenByAFailingStepIsKept(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	// The dangerous case: the charge went through and then the call timed out.
	// What the step managed to write before failing is exactly what the refund
	// needs, so it must not be rolled back with the step.
	Register(Definition{Name: "partial", MaxAttempts: 1, Steps: []Step{
		{Name: "charge", Do: func(_ context.Context, r *Run) error {
			r.Set("charge_id", "ch_late")
			return errors.New("timed out reading the response")
		}},
	}})

	started, _ := Start(context.Background(), db, "partial", nil)
	run := drain(t, db, started.ID)

	if got := string(run.State); got != `{"charge_id":"ch_late"}` {
		t.Fatalf("state = %s, want the charge id the failing step wrote", got)
	}
}

func TestInputReachesTheSteps(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	type order struct {
		ID    string `json:"id"`
		Total int    `json:"total"`
	}
	var seen order
	Register(Definition{Name: "input", Steps: []Step{
		{Name: "read", Do: func(_ context.Context, r *Run) error { return r.Input(&seen) }},
	}})

	started, _ := Start(context.Background(), db, "input", order{ID: "o_1", Total: 4200})
	if run := drain(t, db, started.ID); run.Status != models.SagaDone {
		t.Fatalf("status = %q", run.Status)
	}
	if seen.ID != "o_1" || seen.Total != 4200 {
		t.Errorf("the step read %+v", seen)
	}
}

func TestTheIdempotencyKeyIsStableAcrossRetries(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	var keys []string
	Register(Definition{Name: "keys", MaxAttempts: 3, BaseBackoff: time.Millisecond, Steps: []Step{
		{Name: "charge", Do: func(_ context.Context, r *Run) error {
			keys = append(keys, r.IdempotencyKey())
			if len(keys) < 3 {
				return errors.New("timeout")
			}
			return nil
		}},
	}})

	started, _ := Start(context.Background(), db, "keys", nil)
	drain(t, db, started.ID)

	if len(keys) != 3 {
		t.Fatalf("ran %d times", len(keys))
	}
	// This is the one thing standing between a retried step and a second
	// charge: the provider deduplicates on this key, so it has to be the same
	// string every time.
	for _, k := range keys {
		if k != keys[0] {
			t.Fatalf("the key changed between attempts: %v", keys)
		}
	}
	if keys[0] == "" {
		t.Error("the key is empty, which deduplicates nothing")
	}
}

func TestStartingTwiceWithAKeyReturnsTheSameRun(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)
	Register(Definition{Name: "once", Steps: []Step{{Name: "a", Do: func(context.Context, *Run) error { return nil }}}})

	first, err := Start(context.Background(), db, "once", nil, Key("order:42"))
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	second, err := Start(context.Background(), db, "once", nil, Key("order:42"))
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("a retried start made a second run: %s and %s", first.ID, second.ID)
	}
	var runs int64
	db.Model(&models.SagaRun{}).Count(&runs)
	if runs != 1 {
		t.Errorf("%d runs in the table, want 1", runs)
	}
}

func TestAPanicInAStepIsThatRunsProblemAndNobodyElses(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)

	Register(Definition{Name: "panics", MaxAttempts: 1, Steps: []Step{
		{Name: "boom", Do: func(context.Context, *Run) error { panic("a nil map write, say") }},
	}})

	started, _ := Start(context.Background(), db, "panics", nil)
	run := drain(t, db, started.ID)

	if run.Status != models.SagaCompensated {
		t.Fatalf("status = %q, want compensated", run.Status)
	}
	if run.LastError == "" {
		t.Error("the panic was not recorded as the failure")
	}
}

func TestARunIsClaimedSoTwoRunnersCannotAdvanceIt(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)
	Register(Definition{Name: "claimed", Steps: []Step{
		{Name: "a", Do: func(context.Context, *Run) error { return nil }},
		{Name: "b", Do: func(context.Context, *Run) error { return nil }},
	}})
	started, _ := Start(context.Background(), db, "claimed", nil)

	one := &Runner{DB: db, ClaimTTL: time.Hour, Owner: "one"}
	two := &Runner{DB: db, ClaimTTL: time.Hour, Owner: "two"}
	one.defaults()
	two.defaults()

	claimedByOne, err := one.claim(context.Background())
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if len(claimedByOne) != 1 {
		t.Fatalf("the first runner claimed %d runs, want 1", len(claimedByOne))
	}
	claimedByTwo, err := two.claim(context.Background())
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	// A step that runs on two replicas at once is a card charged twice.
	if len(claimedByTwo) != 0 {
		t.Fatalf("the second runner also claimed the run %s", started.ID)
	}
}

func TestAStaleClaimIsTakenBack(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)
	Register(Definition{Name: "stale", Steps: []Step{{Name: "a", Do: func(context.Context, *Run) error { return nil }}}})
	started, _ := Start(context.Background(), db, "stale", nil)

	// The replica that held it was killed mid-step. Without this, the run would
	// sit claimed by a process that no longer exists, forever.
	long := time.Now().Add(-time.Hour)
	db.Model(&models.SagaRun{}).Where("id = ?", started.ID).
		Updates(map[string]any{"claimed_by": "a dead replica", "claimed_at": long})

	runner := &Runner{DB: db, ClaimTTL: time.Minute, Owner: "alive"}
	runner.defaults()
	claimed, err := runner.claim(context.Background())
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if len(claimed) != 1 {
		t.Fatal("a run whose claim went stale was not picked up")
	}
}

func TestStartRefusesASagaNobodyRegistered(t *testing.T) {
	Reset()
	defer Reset()
	db := testDB(t)
	if _, err := Start(context.Background(), db, "nothing", nil); !errors.Is(err, ErrUnknownSaga) {
		t.Fatalf("err = %v, want ErrUnknownSaga", err)
	}
}

func TestRegisterRefusesADefinitionThatCannotWork(t *testing.T) {
	Reset()
	defer Reset()
	for name, def := range map[string]Definition{
		"no name":  {Steps: []Step{{Name: "a", Do: func(context.Context, *Run) error { return nil }}}},
		"no steps": {Name: "empty"},
		"a step with no name": {Name: "x", Steps: []Step{
			{Do: func(context.Context, *Run) error { return nil }},
		}},
		"a step with no Do": {Name: "y", Steps: []Step{{Name: "a"}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s was accepted", name)
				}
			}()
			Register(def)
		}()
	}
}

func TestRegisteredListsTheSagas(t *testing.T) {
	Reset()
	defer Reset()
	noop := func(context.Context, *Run) error { return nil }
	Register(Definition{Name: "refund", Steps: []Step{{Name: "a", Do: noop}}})
	Register(Definition{Name: "checkout", Steps: []Step{{Name: "a", Do: noop}}})

	names := Registered()
	if len(names) != 2 || names[0] != "checkout" || names[1] != "refund" {
		t.Fatalf("Registered = %v", names)
	}
}
