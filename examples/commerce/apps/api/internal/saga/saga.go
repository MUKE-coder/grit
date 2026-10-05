// Package saga runs a multi-step process that has to either finish or be undone.
//
// # What this is for
//
// A database transaction is the right tool when every write is in one database.
// It is no help at all when the steps are in four places: charge a card, reserve
// stock, book a courier, send the receipt. The card does not roll back when the
// courier refuses, and the process that was holding all of this in its head is
// exactly the thing that crashes.
//
// A saga is the answer to that. Each step has a Do and an Undo. The steps run in
// order, and if one fails for good, the completed ones are undone newest first.
// Where the run got to is a row, not a stack frame, so a process that dies
// between the charge and the reservation resumes rather than leaving a customer
// charged for nothing.
//
// # The three rules
//
//  1. A step must be idempotent. A crash between "the side effect happened" and
//     "we wrote that it happened" is not preventable, so a resumed run will
//     sometimes run a step a second time. Use Run.IdempotencyKey, which is
//     stable for a given run and step, as the key your payment provider or your
//     courier API deduplicates on.
//
//  2. An Undo must be idempotent too, for the same reason, and it must tolerate
//     a Do that never completed. The whole reason compensation exists is that
//     something went wrong, and "went wrong" includes "I do not know whether the
//     charge landed". Write the undo to check first.
//
//  3. Put the steps you cannot undo last. A step with no Undo is declared with
//     Undo: nil and that is a legitimate thing: an email cannot be unsent. But a
//     step after it that fails will compensate everything before it and leave
//     that email sent, so the ordering is the design. Sending the receipt before
//     the parcel is booked is a bug you will only find in production.
//
// # What a step may assume
//
// Exactly one runner advances a run at a time, held by a claim, so a step is not
// racing another copy of itself across replicas. It may assume the steps before
// it completed, and must not assume anything about the process that ran them:
// read what they left in Run.State, not from memory. Run.State is the only thing
// a compensation has, because by then the process that charged the card is
// usually gone.
//
// # What happens when a step fails
//
// It is retried with exponential backoff up to MaxAttempts. Past that the run
// turns to compensating and the completed steps are undone in reverse. If a
// compensation itself fails past its attempts the run goes stuck, which is the
// status that needs a person: something happened that could not be undone. A
// compensation is never skipped, because skipping one is how money goes missing.
package saga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// Step is one unit of work and the way to take it back.
type Step struct {
	// Name identifies the step in the saga_steps table and in the logs. It is
	// written down, so renaming one leaves the old name on existing rows.
	Name string

	// Do performs the step. Returning an error retries it; returning a Fatal
	// error gives up at once and starts compensating.
	Do func(ctx context.Context, r *Run) error

	// Undo takes the step back. Nil means this step cannot be undone, which is
	// allowed and is a design decision: see the package comment on ordering.
	Undo func(ctx context.Context, r *Run) error

	// MaxAttempts overrides the definition's default for this step. For a step
	// that calls something known to be flaky, or one that must not be retried
	// at all (set it to 1).
	MaxAttempts int
}

// Definition is a saga: an ordered list of steps under a name.
type Definition struct {
	// Name is what Start takes and what the saga_runs rows carry.
	Name  string
	Steps []Step

	// MaxAttempts is the default for a step that does not set its own.
	MaxAttempts int
	// BaseBackoff is the first retry delay; it doubles each attempt up to
	// MaxBackoff, with jitter so a thousand runs waiting on the same outage do
	// not all retry in the same millisecond.
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
}

func (d *Definition) defaults() {
	if d.MaxAttempts <= 0 {
		d.MaxAttempts = 5
	}
	if d.BaseBackoff <= 0 {
		d.BaseBackoff = 5 * time.Second
	}
	if d.MaxBackoff <= 0 {
		d.MaxBackoff = 15 * time.Minute
	}
}

func (d *Definition) attemptsFor(i int) int {
	if i >= 0 && i < len(d.Steps) && d.Steps[i].MaxAttempts > 0 {
		return d.Steps[i].MaxAttempts
	}
	return d.MaxAttempts
}

// Fatal wraps an error that must not be retried: a card declined, a postcode
// that does not exist, a payload the other side will never accept. The run
// starts compensating at once instead of spending its attempts on something
// that cannot succeed.
func Fatal(err error) error { return fatal{err} }

type fatal struct{ err error }

func (f fatal) Error() string { return f.err.Error() }
func (f fatal) Unwrap() error { return f.err }

// IsFatal reports whether an error says not to retry.
func IsFatal(err error) bool {
	var f fatal
	return errors.As(err, &f)
}

var (
	mu          sync.RWMutex
	definitions = map[string]*Definition{}
)

// Register adds a saga. Call it once while the application is being built, from
// the file that defines the saga. Registering the same name twice replaces the
// first, so a reload cannot leave two definitions fighting over one name.
func Register(d Definition) {
	if d.Name == "" {
		panic("saga: a definition needs a name")
	}
	if len(d.Steps) == 0 {
		panic("saga: " + d.Name + " has no steps")
	}
	for i, s := range d.Steps {
		if s.Name == "" {
			panic(fmt.Sprintf("saga: %s step %d has no name", d.Name, i))
		}
		if s.Do == nil {
			panic(fmt.Sprintf("saga: %s step %q has no Do", d.Name, s.Name))
		}
	}
	d.defaults()
	mu.Lock()
	defer mu.Unlock()
	definitions[d.Name] = &d
}

// Registered lists the saga names, sorted. For a health probe, an admin screen,
// or a test that wants to know the wiring happened.
func Registered() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(definitions))
	for name := range definitions {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Lookup returns a registered definition.
func Lookup(name string) (*Definition, bool) {
	mu.RLock()
	defer mu.RUnlock()
	d, ok := definitions[name]
	return d, ok
}

// Reset clears the registry. For tests only.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	definitions = map[string]*Definition{}
}

// Run is the handle a step is given.
//
// It is the step's whole world: the input the saga was started with, the state
// the steps before it left, and a key to deduplicate on. There is deliberately
// no way from here to reach the step list or to skip ahead; a step that decides
// what runs next is a state machine, and internal/workflow is the thing for
// that.
type Run struct {
	// Row is the run as the database has it. Read it; the engine writes it.
	Row *models.SagaRun
	// DB is the handle, for a step that has its own tables to write.
	DB *gorm.DB
	// StepIndex and StepName are the step currently executing.
	StepIndex int
	StepName  string

	state map[string]any
}

// Input decodes what the saga was started with into v.
func (r *Run) Input(v any) error {
	if len(r.Row.Input) == 0 {
		return nil
	}
	return json.Unmarshal(r.Row.Input, v)
}

// Get reads a value an earlier step left. The second return is false when
// nothing is under that key, which a compensation should check rather than
// assume: it may be undoing a step whose Do never got far enough to write.
func (r *Run) Get(key string) (any, bool) {
	v, ok := r.state[key]
	return v, ok
}

// GetString is Get for the common case, returning "" when the key is absent or
// is not a string.
func (r *Run) GetString(key string) string {
	if v, ok := r.state[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Set records a value for the steps after this one, and for this step's own
// Undo. It is written to the database when the step completes, so a value set
// by a step that then fails is kept: that is deliberate, because a charge id
// obtained just before a crash is exactly what the refund needs.
func (r *Run) Set(key string, value any) {
	if r.state == nil {
		r.state = map[string]any{}
	}
	r.state[key] = value
}

// IdempotencyKey is a stable key for this run and step, for the API you are
// about to call. The same run retrying the same step produces the same key, so
// a provider that deduplicates on it will not charge twice.
func (r *Run) IdempotencyKey() string {
	return fmt.Sprintf("%s:%d:%s", r.Row.ID, r.StepIndex, r.StepName)
}

func (r *Run) loadState() error {
	r.state = map[string]any{}
	if len(r.Row.State) == 0 {
		return nil
	}
	return json.Unmarshal(r.Row.State, &r.state)
}

func (r *Run) encodeState() ([]byte, error) {
	if len(r.state) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(r.state)
}

// StartOption configures a Start call.
type StartOption func(*models.SagaRun)

// Key makes a start idempotent. Starting the same saga with the same key
// returns the run that already exists instead of a second one.
func Key(k string) StartOption {
	return func(r *models.SagaRun) {
		if k != "" {
			r.Key = &k
		}
	}
}

// After delays the first step.
func After(d time.Duration) StartOption {
	return func(r *models.SagaRun) { r.AvailableAt = time.Now().Add(d) }
}

// ErrUnknownSaga is returned by Start for a name nothing registered.
var ErrUnknownSaga = errors.New("saga: no definition with that name")

// Start records a run and returns at once. A runner advances it.
//
// It does not execute anything inline, and that is the point: a saga exists
// because the process doing the work can die, so the work does not belong to the
// request that asked for it. Start in the same transaction as the write that
// justifies it and the two commit together, the way outbox.Enqueue does.
func Start(ctx context.Context, db *gorm.DB, name string, input any, opts ...StartOption) (*models.SagaRun, error) {
	def, ok := Lookup(name)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownSaga, name)
	}

	payload := []byte("null")
	if input != nil {
		encoded, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("encoding saga input: %w", err)
		}
		payload = encoded
	}

	run := &models.SagaRun{
		Name:        def.Name,
		Input:       payload,
		State:       []byte("{}"),
		Status:      models.SagaRunning,
		AvailableAt: time.Now(),
	}
	for _, opt := range opts {
		opt(run)
	}

	if err := db.WithContext(ctx).Create(run).Error; err != nil {
		// A duplicate key means somebody already started this run, which is
		// what the key is for. Hand back the one that exists.
		if run.Key != nil && isDuplicateKey(err) {
			var existing models.SagaRun
			if lookupErr := db.WithContext(ctx).Where("`key` = ?", *run.Key).First(&existing).Error; lookupErr == nil {
				return &existing, nil
			}
		}
		return nil, fmt.Errorf("starting saga %s: %w", name, err)
	}

	// The step rows are created up front so a run's shape is visible before it
	// has done anything, which is what an operator wants from a stuck run.
	steps := make([]models.SagaStep, 0, len(def.Steps))
	for i, s := range def.Steps {
		steps = append(steps, models.SagaStep{RunID: run.ID, Idx: i, Name: s.Name, Status: models.SagaStepPending})
	}
	if err := db.WithContext(ctx).Create(&steps).Error; err != nil {
		return nil, fmt.Errorf("recording saga steps: %w", err)
	}
	return run, nil
}

// Runner advances runs. One per process; several processes are fine, because a
// run is claimed before it is touched.
type Runner struct {
	DB *gorm.DB
	// Interval is how often to look for work.
	Interval time.Duration
	// Batch is how many runs to claim at once.
	Batch int
	// ClaimTTL is how long a claim is honoured before another runner may take
	// the run back. It has to be longer than the slowest step, or two runners
	// will run that step at once.
	ClaimTTL time.Duration
	// Owner identifies this runner in claimed_by. Defaults to the hostname and
	// the process id.
	Owner string
}

func (r *Runner) defaults() {
	if r.Interval <= 0 {
		r.Interval = 5 * time.Second
	}
	if r.Batch <= 0 {
		r.Batch = 10
	}
	if r.ClaimTTL <= 0 {
		r.ClaimTTL = 5 * time.Minute
	}
	if r.Owner == "" {
		host, _ := os.Hostname()
		r.Owner = fmt.Sprintf("%s/%d", host, os.Getpid())
	}
}

// Start runs the loop until ctx is cancelled.
func (r *Runner) Start(ctx context.Context) {
	r.defaults()
	go func() {
		ticker := time.NewTicker(r.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := r.Tick(ctx); err != nil {
					log.Printf("saga: %v", err)
				}
			}
		}
	}()
}

// Tick claims the runs that are due and advances each by one step. It returns
// how many it moved.
//
// One step per tick on purpose. A run that advanced every step in a loop would
// hold its claim for the length of the whole saga, and a crash in the middle
// would leave it unclaimable for the TTL; stepping once and releasing means the
// next tick, on this process or another, picks up exactly where it is.
func (r *Runner) Tick(ctx context.Context) (int, error) {
	r.defaults()
	runs, err := r.claim(ctx)
	if err != nil {
		return 0, err
	}
	for i := range runs {
		r.advance(ctx, &runs[i])
	}
	return len(runs), nil
}

// claim takes ownership of the due runs, so two replicas cannot advance the
// same one. The UPDATE is the lock: whichever transaction writes claimed_by
// first wins, and the loser's update matches no rows.
func (r *Runner) claim(ctx context.Context) ([]models.SagaRun, error) {
	now := time.Now()
	stale := now.Add(-r.ClaimTTL)

	var candidates []models.SagaRun
	err := r.DB.WithContext(ctx).
		Where("status IN ?", []string{models.SagaRunning, models.SagaCompensating}).
		Where("available_at <= ?", now).
		Where("claimed_at IS NULL OR claimed_at < ?", stale).
		Order("available_at asc").
		Limit(r.Batch).
		Find(&candidates).Error
	if err != nil {
		return nil, fmt.Errorf("finding saga runs: %w", err)
	}

	claimed := make([]models.SagaRun, 0, len(candidates))
	for _, c := range candidates {
		res := r.DB.WithContext(ctx).Model(&models.SagaRun{}).
			Where("id = ?", c.ID).
			Where("claimed_at IS NULL OR claimed_at < ?", stale).
			Updates(map[string]any{"claimed_by": r.Owner, "claimed_at": now})
		if res.Error != nil {
			return nil, fmt.Errorf("claiming saga run %s: %w", c.ID, res.Error)
		}
		if res.RowsAffected == 1 {
			c.ClaimedBy = r.Owner
			c.ClaimedAt = &now
			claimed = append(claimed, c)
		}
	}
	return claimed, nil
}

// advance moves one run by one step, forward or backward.
func (r *Runner) advance(ctx context.Context, row *models.SagaRun) {
	def, ok := Lookup(row.Name)
	if !ok {
		// A run whose definition is gone: a saga was deleted, or this replica is
		// an older build. Release it rather than failing it, because the other
		// replica may know what it is.
		r.release(ctx, row)
		log.Printf("saga: run %s names %q, which nothing registered here", row.ID, row.Name)
		return
	}

	if row.Status == models.SagaCompensating {
		r.undoOne(ctx, def, row)
		return
	}
	r.doOne(ctx, def, row)
}

func (r *Runner) doOne(ctx context.Context, def *Definition, row *models.SagaRun) {
	if row.Cursor >= len(def.Steps) {
		r.finish(ctx, row, models.SagaDone)
		return
	}
	idx := row.Cursor
	step := def.Steps[idx]

	run := &Run{Row: row, DB: r.DB, StepIndex: idx, StepName: step.Name}
	if err := run.loadState(); err != nil {
		r.fail(ctx, def, row, idx, fmt.Errorf("reading saga state: %w", err))
		return
	}

	r.markStep(ctx, row.ID, idx, map[string]any{"status": models.SagaStepPending, "started_at": time.Now()})
	err := call(ctx, step.Do, run)

	// The state is saved whether the step succeeded or not. A charge id written
	// just before a timeout is exactly what the refund will need.
	state, encodeErr := run.encodeState()
	if encodeErr != nil {
		r.fail(ctx, def, row, idx, fmt.Errorf("encoding saga state: %w", encodeErr))
		return
	}
	row.State = state

	if err != nil {
		r.fail(ctx, def, row, idx, err)
		return
	}

	now := time.Now()
	r.markStep(ctx, row.ID, idx, map[string]any{
		"status": models.SagaStepDone, "finished_at": now, "last_error": "",
	})
	row.Cursor = idx + 1
	row.Attempts = 0
	row.LastError = ""
	if row.Cursor >= len(def.Steps) {
		r.save(ctx, row, map[string]any{
			"state": row.State, "cursor": row.Cursor, "attempts": 0, "last_error": "",
			"status": models.SagaDone, "finished_at": now, "claimed_by": "", "claimed_at": nil,
		})
		return
	}
	r.save(ctx, row, map[string]any{
		"state": row.State, "cursor": row.Cursor, "attempts": 0, "last_error": "",
		"available_at": time.Now(), "claimed_by": "", "claimed_at": nil,
	})
}

// fail records a step's failure and decides whether to retry it or to start
// undoing everything before it.
func (r *Runner) fail(ctx context.Context, def *Definition, row *models.SagaRun, idx int, cause error) {
	attempts := row.Attempts + 1
	limit := def.attemptsFor(idx)
	r.markStep(ctx, row.ID, idx, map[string]any{"attempts": attempts, "last_error": cause.Error()})

	if !IsFatal(cause) && attempts < limit {
		r.save(ctx, row, map[string]any{
			"state": row.State, "attempts": attempts, "last_error": cause.Error(),
			"available_at": time.Now().Add(backoff(def, attempts)),
			"claimed_by":   "", "claimed_at": nil,
		})
		return
	}

	// Out of attempts, or told not to retry. The step did not complete, so it is
	// not compensated; everything before it is.
	r.markStep(ctx, row.ID, idx, map[string]any{"status": models.SagaStepFailed, "finished_at": time.Now()})
	log.Printf("saga: run %s step %q failed for good after %d attempt(s): %v, compensating", row.ID, def.Steps[idx].Name, attempts, cause)

	row.Status = models.SagaCompensating
	row.Cursor = idx - 1
	r.save(ctx, row, map[string]any{
		"state": row.State, "status": models.SagaCompensating, "cursor": row.Cursor,
		"attempts": 0, "last_error": cause.Error(),
		"available_at": time.Now(), "claimed_by": "", "claimed_at": nil,
	})
}

func (r *Runner) undoOne(ctx context.Context, def *Definition, row *models.SagaRun) {
	if row.Cursor < 0 {
		r.finish(ctx, row, models.SagaCompensated)
		return
	}
	idx := row.Cursor
	if idx >= len(def.Steps) {
		// The definition shrank under a running saga. Walk down rather than
		// index past the end.
		row.Cursor = len(def.Steps) - 1
		r.save(ctx, row, map[string]any{"cursor": row.Cursor, "claimed_by": "", "claimed_at": nil})
		return
	}
	step := def.Steps[idx]

	// Nothing to undo: either the step declares no Undo, or it never completed.
	//
	// "Completed" includes a step already marked compensating, which is one
	// whose Undo failed and is being retried. Reading only "done" here made the
	// second attempt decide there was nothing to take back, so a refund that
	// failed once was silently abandoned and the run reported itself
	// compensated with the charge still standing.
	var recorded models.SagaStep
	err := r.DB.WithContext(ctx).Where("run_id = ? AND idx = ?", row.ID, idx).First(&recorded).Error
	needsUndo := recorded.Status == models.SagaStepDone || recorded.Status == models.SagaStepCompensating
	if step.Undo == nil || err != nil || !needsUndo {
		r.save(ctx, row, map[string]any{
			"cursor": idx - 1, "attempts": 0,
			"available_at": time.Now(), "claimed_by": "", "claimed_at": nil,
		})
		return
	}

	run := &Run{Row: row, DB: r.DB, StepIndex: idx, StepName: step.Name}
	if err := run.loadState(); err != nil {
		log.Printf("saga: run %s cannot read its state to compensate: %v", row.ID, err)
	}
	r.markStep(ctx, row.ID, idx, map[string]any{"status": models.SagaStepCompensating})

	if err := call(ctx, step.Undo, run); err != nil {
		attempts := row.Attempts + 1
		limit := def.attemptsFor(idx)
		r.markStep(ctx, row.ID, idx, map[string]any{"attempts": attempts, "last_error": err.Error()})
		if attempts < limit {
			r.save(ctx, row, map[string]any{
				"attempts": attempts, "last_error": err.Error(),
				"available_at": time.Now().Add(backoff(def, attempts)),
				"claimed_by":   "", "claimed_at": nil,
			})
			return
		}
		// A compensation that cannot be completed is not skipped. The run stops
		// here and waits for a person, because the alternative is to carry on
		// and quietly leave the thing it was undoing in place.
		r.markStep(ctx, row.ID, idx, map[string]any{"status": models.SagaStepStuck, "finished_at": time.Now()})
		log.Printf("saga: run %s is STUCK: compensating step %q failed %d time(s): %v", row.ID, step.Name, attempts, err)
		r.save(ctx, row, map[string]any{
			"status": models.SagaStuck, "last_error": err.Error(), "attempts": attempts,
			"finished_at": time.Now(), "claimed_by": "", "claimed_at": nil,
		})
		return
	}

	state, _ := run.encodeState()
	row.State = state
	r.markStep(ctx, row.ID, idx, map[string]any{
		"status": models.SagaStepCompensated, "finished_at": time.Now(), "last_error": "",
	})
	r.save(ctx, row, map[string]any{
		"state": state, "cursor": idx - 1, "attempts": 0,
		"available_at": time.Now(), "claimed_by": "", "claimed_at": nil,
	})
}

func (r *Runner) finish(ctx context.Context, row *models.SagaRun, status string) {
	now := time.Now()
	r.save(ctx, row, map[string]any{
		"status": status, "finished_at": now, "claimed_by": "", "claimed_at": nil,
	})
}

func (r *Runner) release(ctx context.Context, row *models.SagaRun) {
	r.save(ctx, row, map[string]any{"claimed_by": "", "claimed_at": nil})
}

func (r *Runner) save(ctx context.Context, row *models.SagaRun, fields map[string]any) {
	if err := r.DB.WithContext(ctx).Model(&models.SagaRun{}).Where("id = ?", row.ID).Updates(fields).Error; err != nil {
		log.Printf("saga: writing run %s: %v", row.ID, err)
	}
}

func (r *Runner) markStep(ctx context.Context, runID string, idx int, fields map[string]any) {
	err := r.DB.WithContext(ctx).Model(&models.SagaStep{}).
		Where("run_id = ? AND idx = ?", runID, idx).Updates(fields).Error
	if err != nil {
		log.Printf("saga: writing step %d of run %s: %v", idx, runID, err)
	}
}

// call runs a step body and turns a panic into an error.
//
// A step is application code calling somebody else's API, which is where nil
// maps and unchecked type assertions live. A panic that took the runner down
// would stop every other saga in the process, so it becomes this run's failure
// and nothing else's.
func call(ctx context.Context, fn func(context.Context, *Run) error, run *Run) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("step panicked: %v", rec)
		}
	}()
	return fn(ctx, run)
}

// backoff is exponential with jitter, so a thousand runs waiting on the same
// outage do not all retry in the same millisecond and knock it over again.
func backoff(def *Definition, attempt int) time.Duration {
	d := float64(def.BaseBackoff) * math.Pow(2, float64(attempt-1))
	if d > float64(def.MaxBackoff) || math.IsInf(d, 0) {
		d = float64(def.MaxBackoff)
	}
	// #nosec G404 -- the jitter spreads retries across a window so a thousand
	// runs waiting on one outage do not all come back in the same
	// millisecond. Nothing is guessed from it and nothing is protected by it,
	// and crypto/rand here would be a syscall per retry to pick a number
	// nobody can act on.
	jitter := 1 + (rand.Float64()-0.5)*0.4 // +/- 20%
	return time.Duration(d * jitter)
}

// isDuplicateKey reports whether an insert lost a race on a unique index. The
// drivers word it differently and none of them gives a typed error worth
// matching, so the message is what there is.
func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, s := range []string{"duplicate key", "UNIQUE constraint", "Duplicate entry", "1062"} {
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}
