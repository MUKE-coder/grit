// Package cron schedules background tasks for Grit Motors.
//
// The scheduler is intentionally tiny: one robfig/cron/v3 instance owned by
// the server, jobs registered at boot in main.go via Register(). Each job is
// idempotent — re-running mid-day after a restart should be safe.
package cron

import (
	"log"
	"time"

	rcron "github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"gritdemo/internal/database"
	"gritdemo/internal/services"
)

// Scheduler wraps a robfig cron with project-specific job registration.
type Scheduler struct {
	c        *rcron.Cron
	db       *gorm.DB
	demoMode bool
}

// New constructs a scheduler that runs in the local timezone.
func New(db *gorm.DB) *Scheduler {
	return &Scheduler{
		c:  rcron.New(rcron.WithLocation(time.Local)),
		db: db,
	}
}

// WithDemoMode enables the nightly demo-reset job. Call before Register.
func (s *Scheduler) WithDemoMode(on bool) *Scheduler {
	s.demoMode = on
	return s
}

// Register wires up the standard set of jobs and returns the scheduler so
// the caller can Start/Stop it.
func (s *Scheduler) Register() error {
	// Nightly at 00:30 — flip overdue installment rows.
	if _, err := s.c.AddFunc("30 0 * * *", s.runMarkOverdue); err != nil {
		return err
	}
	// Nightly at 00:00 — reset the demo database when DEMO_MODE is on.
	// Wipes mutable rows and re-runs the seeder, leaving the admin user +
	// canonical demo cohort in place so demo.gritframework.dev stays honest.
	if s.demoMode {
		if _, err := s.c.AddFunc("0 0 * * *", s.runDemoReset); err != nil {
			return err
		}
		log.Println("cron: nightly demo-reset job registered (DEMO_MODE=true)")
	}
	return nil
}

func (s *Scheduler) Start() {
	s.c.Start()
}

// Stop gracefully halts the scheduler. Returns a context that's done once any
// in-flight jobs have finished.
func (s *Scheduler) Stop() {
	ctx := s.c.Stop()
	<-ctx.Done()
}

func (s *Scheduler) runMarkOverdue() {
	n, err := services.MarkOverdueSchedules(s.db, time.Now())
	if err != nil {
		log.Printf("cron: mark overdue failed: %v", err)
		return
	}
	if n > 0 {
		log.Printf("cron: marked %d schedule rows overdue", n)
	}
}

// runDemoReset wipes the mutable demo data (sales, payments, loans,
// stock movements, notifications, audit log) and re-runs the seeder.
// Only registered in DEMO_MODE. The admin user is preserved because the
// seeder treats it as canonical and tops it up on every run.
func (s *Scheduler) runDemoReset() {
	log.Println("cron: demo reset starting…")
	if err := database.WipeMutableData(s.db); err != nil {
		log.Printf("cron: demo wipe failed: %v", err)
		return
	}
	if err := database.Seed(s.db); err != nil {
		log.Printf("cron: demo reseed failed: %v", err)
		return
	}
	log.Println("cron: demo reset complete")
}
