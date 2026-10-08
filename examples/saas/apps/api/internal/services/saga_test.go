package services

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas/apps/api/internal/models"
	"saas/apps/api/internal/paginate"
)

func sagaTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&models.SagaRun{}, &models.SagaStep{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// stuckRun is a run whose compensation failed past its attempts: the charge was
// not taken back, so there is money standing that should not be.
func stuckRun(t *testing.T, db *gorm.DB) *models.SagaRun {
	t.Helper()
	run := &models.SagaRun{
		Name: "checkout", Status: models.SagaStuck, Cursor: 0, Attempts: 5,
		LastError: "the refund API returned 503",
	}
	if err := db.Create(run).Error; err != nil {
		t.Fatal(err)
	}
	steps := []models.SagaStep{
		{RunID: run.ID, Idx: 0, Name: "charge", Status: models.SagaStepStuck, Attempts: 5, LastError: "the refund API returned 503"},
		{RunID: run.ID, Idx: 1, Name: "ship", Status: models.SagaStepFailed, Attempts: 3},
	}
	if err := db.Create(&steps).Error; err != nil {
		t.Fatal(err)
	}
	return run
}

func TestSagaCountsReportEveryStatusIncludingZero(t *testing.T) {
	db := sagaTestDB(t)
	svc := &SagaService{DB: db}
	for _, status := range []string{models.SagaDone, models.SagaDone, models.SagaStuck} {
		if err := db.Create(&models.SagaRun{Name: "checkout", Status: status}).Error; err != nil {
			t.Fatal(err)
		}
	}

	counts, err := svc.Counts(context.Background())
	if err != nil {
		t.Fatalf("counts: %v", err)
	}
	if counts[models.SagaDone] != 2 || counts[models.SagaStuck] != 1 {
		t.Fatalf("counts = %v", counts)
	}
	// A status with no runs is present at zero rather than absent, so the screen
	// can draw a chip for each without deciding what a missing key means.
	if _, ok := counts[models.SagaRunning]; !ok {
		t.Error("running is missing from the counts, so the chip would not be drawn")
	}
}

func TestRetryPutsAStuckRunBackInFrontOfTheRunner(t *testing.T) {
	db := sagaTestDB(t)
	svc := &SagaService{DB: db}
	run := stuckRun(t, db)

	out, err := svc.Retry(context.Background(), run.ID)
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if out.Status != models.SagaCompensating {
		t.Fatalf("status = %q, want compensating", out.Status)
	}
	if out.Attempts != 0 || out.LastError != "" || out.FinishedAt != nil {
		t.Errorf("the run was not reset: attempts %d, error %q, finished %v", out.Attempts, out.LastError, out.FinishedAt)
	}
	// The step that could not be undone has to go back to compensating, or the
	// engine sees a step it has already given up on and walks past it, leaving
	// the charge standing while the run reports itself compensated.
	var step models.SagaStep
	if err := db.Where("run_id = ? AND idx = 0", run.ID).First(&step).Error; err != nil {
		t.Fatal(err)
	}
	if step.Status != models.SagaStepCompensating {
		t.Errorf("the stuck step is %q, want compensating", step.Status)
	}
	// And the claim is cleared, or the run waits out the whole TTL before any
	// runner will touch it.
	if out.ClaimedBy != "" || out.ClaimedAt != nil {
		t.Error("the claim was not released")
	}
}

func TestRetryRefusesARunThatIsNotStuck(t *testing.T) {
	db := sagaTestDB(t)
	svc := &SagaService{DB: db}
	for _, status := range []string{models.SagaRunning, models.SagaDone, models.SagaCompensated, models.SagaCompensating} {
		run := &models.SagaRun{Name: "checkout", Status: status}
		if err := db.Create(run).Error; err != nil {
			t.Fatal(err)
		}
		// A running one needs no help and a finished one has nothing left to
		// do; retrying either would re-run work that already happened.
		if _, err := svc.Retry(context.Background(), run.ID); err == nil {
			t.Errorf("a %s run was retried", status)
		}
	}
}

func TestByIDReturnsTheStepsInOrder(t *testing.T) {
	db := sagaTestDB(t)
	svc := &SagaService{DB: db}
	run := stuckRun(t, db)

	out, err := svc.ByID(context.Background(), run.ID)
	if err != nil {
		t.Fatalf("by id: %v", err)
	}
	if len(out.Steps) != 2 {
		t.Fatalf("got %d steps", len(out.Steps))
	}
	// The order is the whole point: steps run top to bottom and are undone
	// bottom to top, and a list in insertion order would read as neither.
	if out.Steps[0].Idx != 0 || out.Steps[1].Idx != 1 {
		t.Errorf("the steps came back in the order %d, %d", out.Steps[0].Idx, out.Steps[1].Idx)
	}
}

func TestListFiltersByStatus(t *testing.T) {
	db := sagaTestDB(t)
	svc := &SagaService{DB: db}
	for _, status := range []string{models.SagaDone, models.SagaStuck, models.SagaStuck} {
		if err := db.Create(&models.SagaRun{Name: "checkout", Status: status}).Error; err != nil {
			t.Fatal(err)
		}
	}

	page, err := svc.List(context.Background(), paginate.Params{
		Page: 1, PageSize: 20, Filters: map[string]any{"status": models.SagaStuck},
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if page.Meta.Total != 2 {
		t.Fatalf("total = %d, want 2", page.Meta.Total)
	}
	for _, run := range page.Data {
		if run.Status != models.SagaStuck {
			t.Errorf("the filter let a %s run through", run.Status)
		}
	}
}
