package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/paginate"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/saga"
	"library/apps/api/internal/services"
)

// SagaHandler serves the saga runs screen.
//
// A run that goes stuck needs a person: its compensation failed past its
// attempts, so something happened that could not be taken back and the run will
// not move again on its own. Without a screen the only way to find one is SQL
// against saga_runs, which means nobody finds one until a customer complains.
type SagaHandler struct {
	DB      *gorm.DB
	Service *services.SagaService
}

func (h *SagaHandler) service() *services.SagaService {
	if h.Service != nil {
		return h.Service
	}
	return &services.SagaService{DB: h.DB}
}

// List returns one page of runs, with the counts for the status chips.
func (h *SagaHandler) List(c *gin.Context) {
	page, err := h.service().List(c.Request.Context(), paginate.Bind(c))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not read the saga runs")
		return
	}
	counts, err := h.service().Counts(c.Request.Context())
	if err != nil {
		// The list is the answer; the chips are decoration. A failed count
		// should not blank the page, and an empty map draws no chips rather
		// than drawing zeros that cannot be told from a quiet day.
		counts = nil
	}
	c.JSON(http.StatusOK, gin.H{"data": page.Data, "meta": page.Meta, "counts": counts})
}

// GetByID returns one run with its steps in order.
func (h *SagaHandler) GetByID(c *gin.Context) {
	run, err := h.service().ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "Saga run not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": run})
}

// Retry puts a stuck run back in front of the runner, for after the thing that
// was broken has been fixed.
func (h *SagaHandler) Retry(c *gin.Context) {
	run, err := h.service().Retry(c.Request.Context(), c.Param("id"))
	if err != nil {
		// The service refuses a run that is not stuck, and says why. That is a
		// 422 rather than a 500: the request was understood and is not something
		// this run can do.
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": run, "message": "The run is compensating again"})
}

// Definitions lists the registered sagas and their steps.
//
// From the registry rather than from the rows, so the screen can show what a
// saga is meant to do even when nothing has run yet, and so a run of a saga
// whose definition has since changed can be read against the current one.
func (h *SagaHandler) Definitions(c *gin.Context) {
	names := saga.Registered()
	out := make([]gin.H, 0, len(names))
	for _, name := range names {
		def, ok := saga.Lookup(name)
		if !ok {
			continue
		}
		steps := make([]gin.H, 0, len(def.Steps))
		for _, s := range def.Steps {
			steps = append(steps, gin.H{"name": s.Name, "compensable": s.Undo != nil})
		}
		out = append(out, gin.H{"name": def.Name, "steps": steps, "max_attempts": def.MaxAttempts})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
