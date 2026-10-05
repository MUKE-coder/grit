package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"commerce/apps/api/internal/ids"
)

// Saga run statuses. A run is in exactly one of these at any moment.
const (
	// SagaRunning means steps are still being worked through.
	SagaRunning = "running"
	// SagaCompensating means a step failed for good and the completed steps
	// before it are being undone, newest first.
	SagaCompensating = "compensating"
	// SagaDone means every step completed.
	SagaDone = "done"
	// SagaCompensated means a step failed and every completed step before it
	// was undone. The world is back where it started, which is the good outcome
	// of a bad run.
	SagaCompensated = "compensated"
	// SagaStuck means a compensation itself failed past its attempts. This is
	// the status that needs a person: something happened that could not be
	// undone, and the run will not move again on its own.
	SagaStuck = "stuck"
)

// Saga step statuses.
const (
	SagaStepPending      = "pending"
	SagaStepDone         = "done"
	SagaStepFailed       = "failed"
	SagaStepCompensated  = "compensated"
	SagaStepCompensating = "compensating"
	SagaStepStuck        = "stuck"
)

// SagaRun is one execution of a multi-step process that has to either finish or
// be undone.
//
// A saga is for the work that crosses systems, where a database transaction
// cannot help you: charge a card, reserve stock, book a courier, send the
// receipt. Those happen in four places, none of which rolls back when the next
// one fails, so "undo" has to be a thing you write and a thing that is run for
// you. This row is what makes that possible after a crash: it holds where the
// run got to, so a process that dies between the charge and the reservation
// resumes rather than leaving a customer charged for nothing.
type SagaRun struct {
	ID string `gorm:"type:varchar(36);primaryKey" json:"id"`

	// Name is the registered saga this run belongs to, such as "checkout".
	Name string `gorm:"size:100;not null;index:idx_saga_runs_name" json:"name"`

	// Key is the caller's idempotency key, unique across the table. Starting a
	// saga twice with the same key returns the first run rather than a second
	// one, which is what you want when the start is behind a retried HTTP
	// request or a redelivered webhook.
	//
	// Empty means "no key", and an empty string cannot be unique across many
	// rows, so it is stored as NULL. That is why this is a pointer.
	Key *string `gorm:"size:255;uniqueIndex" json:"key,omitempty"`

	// Input is what the run was started with. It never changes.
	Input datatypes.JSON `gorm:"type:json" json:"input"`

	// State is what the steps have written for the steps after them: the charge
	// id the refund will need, the reservation id the release will need. It is
	// the only thing a compensation can rely on, because the process that ran
	// the step may be long gone by the time the undo runs.
	State datatypes.JSON `gorm:"type:json" json:"state"`

	Status string `gorm:"size:20;not null;default:'running';index:idx_saga_claim,priority:1" json:"status"`

	// Cursor is the index of the next step to run while the run is going
	// forward, and the index of the next step to undo while it is compensating.
	Cursor int `gorm:"not null;default:0" json:"cursor"`

	Attempts  int    `gorm:"not null;default:0" json:"attempts"`
	LastError string `gorm:"type:text" json:"last_error,omitempty"`

	// AvailableAt is when a runner may next pick this up. It moves forward on
	// each failure, which is the backoff.
	AvailableAt time.Time `gorm:"not null;index:idx_saga_claim,priority:2" json:"available_at"`

	// ClaimedBy is the runner holding this run, so two replicas cannot advance
	// the same run at once. A step that runs twice is a card charged twice.
	ClaimedBy  string     `gorm:"size:64" json:"claimed_by,omitempty"`
	ClaimedAt  *time.Time `json:"claimed_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`

	Steps []SagaStep `gorm:"foreignKey:RunID;constraint:OnDelete:CASCADE" json:"steps,omitempty"`

	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SagaRun) TableName() string { return "saga_runs" }

// BeforeCreate fills the id, the first availability and the starting status.
func (r *SagaRun) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ids.New()
	}
	if r.AvailableAt.IsZero() {
		r.AvailableAt = time.Now()
	}
	if r.Status == "" {
		r.Status = SagaRunning
	}
	return nil
}

// Finished reports whether this run will move again on its own. A stuck run is
// finished in that sense and is the one to go and look at.
func (r *SagaRun) Finished() bool {
	return r.Status == SagaDone || r.Status == SagaCompensated || r.Status == SagaStuck
}

// SagaStep is one step of one run, and the record of what it did.
//
// Kept after the run finishes rather than deleted. When a customer says they
// were charged and never shipped, this table is the answer, and when a
// compensation fails these rows are what somebody works from by hand.
type SagaStep struct {
	ID    string `gorm:"type:varchar(36);primaryKey" json:"id"`
	RunID string `gorm:"type:varchar(36);not null;index:idx_saga_steps_run" json:"run_id"`

	// Idx is the step's position in the definition. Compensation walks these
	// downward, which is the whole point of recording it.
	Idx  int    `gorm:"not null" json:"idx"`
	Name string `gorm:"size:100;not null" json:"name"`

	Status    string `gorm:"size:20;not null;default:'pending'" json:"status"`
	Attempts  int    `gorm:"not null;default:0" json:"attempts"`
	LastError string `gorm:"type:text" json:"last_error,omitempty"`

	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SagaStep) TableName() string { return "saga_steps" }

func (s *SagaStep) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	if s.Status == "" {
		s.Status = SagaStepPending
	}
	return nil
}
