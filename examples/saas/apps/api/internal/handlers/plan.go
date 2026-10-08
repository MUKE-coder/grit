package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/authz"
	"saas/apps/api/internal/concurrency"
	"saas/apps/api/internal/events"
	"saas/apps/api/internal/export"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/money"
	"saas/apps/api/internal/paginate"
	"saas/apps/api/internal/pdf"
	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
)

// PlanHandler serves the plan endpoints. It reads the request, asks
// services.PlanService, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type PlanHandler struct {
	DB *gorm.DB
	// AppName brands the PDF export. The route file sets it from config.
	AppName string
}

// service is the plan service over this handler's database.
func (h *PlanHandler) service() *services.PlanService {
	return &services.PlanService{DB: h.DB}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *PlanHandler) ctx(c *gin.Context) context.Context {
	return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *PlanHandler) fail(c *gin.Context, err error, fallback string) {
	var conflict *concurrency.ErrConflict
	var denied *authz.DeniedError
	switch {
	case errors.As(err, &conflict):
		concurrency.WriteConflict(c, conflict.Current)
	case errors.As(err, &denied):
		// A policy refused, and said why. 403 rather than the 404 that
		// ownership uses: the caller is already allowed to see this row,
		// so hiding it says nothing and withholds the reason, which is the
		// thing the rule exists to give.
		respond.Fail(c, respond.CodeForbidden, denied.Reason)
	case errors.Is(err, gorm.ErrRecordNotFound):
		respond.Fail(c, respond.CodeNotFound, "Plan not found")
	default:
		respond.WriteError(c, err, fallback)
	}
}

// List returns a paginated list of plans.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *PlanHandler) List(c *gin.Context) {
	res, err := h.service().List(h.ctx(c), paginate.Bind(c), c.Query("archived"))
	if err != nil {
		h.fail(c, err, "Failed to fetch plans")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX: the same
// search as List, without pages.
//
//	GET /api/plans/export?format=csv
//	GET /api/plans/export?format=xlsx&search=foo
func (h *PlanHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	opts := export.Options{
		Sheet: "Plans",
		Columns: []export.Column{
			{Header: "ID", Field: "ID"},
			{Header: "Name", Field: "Name"},
			{Header: "Slug", Field: "Slug"},
			{Header: "Price", Field: "Price"},
			{Header: "Interval", Field: "Interval"},
			{Header: "Seats", Field: "Seats"},
			{Header: "Blurb", Field: "Blurb"},
			{Header: "Active", Field: "Active", Format: "bool"},
			{Header: "Created At", Field: "CreatedAt", Format: "date:2006-01-02"},
		},
	}

	if format == "xlsx" {
		// Rows go into the workbook batch by batch, through a writer that moves
		// them to a temporary file, so memory stays flat however many match.
		sheet, err := export.NewXLSXStream(opts)
		if err != nil {
			h.fail(c, err, "Failed to export plans")
			return
		}
		defer func() {
			if err := sheet.Close(); err != nil {
				log.Printf("export plans as xlsx: removing temporary files: %v", err)
			}
		}()
		if err := h.service().Export(h.ctx(c), search, func(rows []models.Plan) error {
			return sheet.Rows(rows)
		}); err != nil {
			h.fail(c, err, "Failed to export plans")
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", `attachment; filename="plans.xlsx"`)
		if err := sheet.Finish(c.Writer); err != nil {
			log.Printf("export plans as xlsx: %v", err)
		}
		return
	}

	// CSV streams: the header with the first batch, then each batch as it is
	// read.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="plans.csv"`)

	headerWritten := false
	err := h.service().Export(h.ctx(c), search, func(rows []models.Plan) error {
		if !headerWritten {
			headerWritten = true
			return export.CSV(c.Writer, rows, opts)
		}
		return export.CSVRows(c.Writer, rows, opts)
	})
	if err == nil && !headerWritten {
		// Nothing matched: still a CSV, with its header row.
		err = export.CSV(c.Writer, []models.Plan{}, opts)
	}
	if err != nil {
		if !c.Writer.Written() {
			c.Writer.Header().Del("Content-Type")
			c.Writer.Header().Del("Content-Disposition")
			h.fail(c, err, "Failed to export plans")
			return
		}
		// Rows are already on the wire, so the status cannot change. Logged,
		// so a truncated file has an explanation somewhere.
		log.Printf("export plans: %v", err)
	}
}

// GetByID returns a single plan by ID.
func (h *PlanHandler) GetByID(c *gin.Context) {
	item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to load plan")
		return
	}

	// The version to send back as If-Match when saving.
	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this plan as a print-ready PDF: a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *PlanHandler) PDF(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service().GetByID(h.ctx(c), id)
	if err != nil {
		h.fail(c, err, "Failed to load plan")
		return
	}

	appName := h.AppName
	if appName == "" {
		appName = "Plan"
	}

	rec := pdf.Record{
		Title:      "PLAN",
		Subtitle:   pdf.Value(item.Name),
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
			{Label: "Name", Value: pdf.Value(item.Name)},
			{Label: "Price", Value: pdf.Value(item.Price)},
			{Label: "Interval", Value: pdf.Value(item.Interval)},
			{Label: "Seats", Value: pdf.Value(item.Seats)},
			{Label: "Blurb", Value: pdf.Value(item.Blurb)},
			{Label: "Active", Value: pdf.Value(item.Active)},
			{Label: "Created", Value: pdf.Value(item.CreatedAt)},
		},
	}

	out, err := pdf.RenderRecord(rec)
	if err != nil {
		respond.Fail(c, respond.CodePDFError, "could not render the PDF")
		return
	}

	filename := "plan-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// CreatePlanRequest is the JSON body accepted by POST /plans.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type. routes.go passes
// this type to docs.Route(...).RequestBody().
type CreatePlanRequest struct {
	Name     string      `json:"name" binding:"required"`
	Price    money.Money `json:"price"`
	Interval string      `json:"interval" binding:"required"`
	Seats    int         `json:"seats"`
	Blurb    string      `json:"blurb"`
	Active   bool        `json:"active"`
}

// UpdatePlanRequest is the JSON body accepted by PUT /plans/:id.
// Every field is optional: only what the client sends is applied.
type UpdatePlanRequest struct {
	Name     string       `json:"name"`
	Price    *money.Money `json:"price"`
	Interval string       `json:"interval"`
	Seats    *int         `json:"seats"`
	Blurb    string       `json:"blurb"`
	Active   *bool        `json:"active"`
}

// Create adds a new plan.
func (h *PlanHandler) Create(c *gin.Context) {
	var req CreatePlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item := models.Plan{
		Name:     req.Name,
		Price:    req.Price,
		Interval: req.Interval,
		Seats:    req.Seats,
		Blurb:    req.Blurb,
		Active:   req.Active,
	}

	if err := h.service().Create(h.ctx(c), &item); err != nil {
		h.fail(c, err, "Failed to create plan")
		return
	}

	events.Emitted(c, "plans", "Plan", "created", item.ID, item.Name, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "Plan created successfully",
	})
}

// Update modifies an existing plan. A client that sends the version it
// read as If-Match gets 409 instead of overwriting a newer save.
func (h *PlanHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Interval != "" {
		updates["interval"] = req.Interval
	}
	if req.Seats != nil {
		updates["seats"] = *req.Seats
	}
	if req.Blurb != "" {
		updates["blurb"] = req.Blurb
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}

	item, err := h.service().Update(h.ctx(c), id, updates, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to update plan")
		return
	}

	events.Emitted(c, "plans", "Plan", "updated", item.ID, item.Name, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Plan updated successfully",
	})
}

// Patch applies a partial update to a plan. Used by the admin's grouped
// update view: each form group's Save button sends only the fields it owns, so
// editing "Address" does not rewrite "Pricing".
func (h *PlanHandler) Patch(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item, updates, err := h.service().Patch(h.ctx(c), c.Param("id"), body, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to patch plan")
		return
	}

	events.Emitted(c, "plans", "Plan", "updated", item.ID, item.Name, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Plan updated successfully",
	})
}

// Delete soft-deletes a plan.
func (h *PlanHandler) Delete(c *gin.Context) {
	item, err := h.service().Delete(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to delete plan")
		return
	}

	events.Emitted(c, "plans", "Plan", "deleted", item.ID, item.Name, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan deleted successfully",
	})
}

// BulkPlanRequest is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails.
type BulkPlanRequest struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string `json:"action" binding:"required,oneof=delete archive restore patch"`
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string `json:"ids" binding:"required,min=1,max=500"`
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} `json:"patch"`
}

// Bulk applies one action to many plans in a single transaction.
func (h *PlanHandler) Bulk(c *gin.Context) {
	var req BulkPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	result, err := h.service().Bulk(h.ctx(c), req.Action, req.IDs, req.Patch)
	if err != nil {
		h.fail(c, err, "Failed to "+req.Action+" plans")
		return
	}
	if len(result.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    gin.H{"affected": 0, "requested": len(req.IDs)},
			"message": "Nothing to do",
		})
		return
	}
	ids := result.IDs

	// One audit entry naming the action and the count, not N entries that bury
	// everything else somebody did today. A local map, not a package-level
	// helper: every resource has its own handler file in package handlers.
	past := map[string]string{
		"delete":  "deleted",
		"archive": "archived",
		"restore": "restored",
		"patch":   "updated",
	}[req.Action]

	noun := "plans"
	if len(ids) == 1 {
		noun = "plan"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(result.Updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and the
	// count lives in the summary, where it can be read.
	events.Emitted(c, "plans", "Plan", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}

// The two grid endpoints, for the admin's spreadsheet editors.
//
// Separate from Bulk rather than actions on it, because the request shapes are
// genuinely different: Bulk takes ids and ONE patch for all of them, and these
// take a row each. Folding them in would mean a required ids field that two of
// the five actions ignore.

// BulkEditPlanRequest is the grid's save: one patch per row.
type BulkEditPlanRequest struct {
	// Capped for the same reason Bulk's ids are: a transaction holding five
	// hundred row locks is a transaction other writers wait behind.
	Items []services.PlanGridEdit `json:"items" binding:"required,min=1,max=500"`
}

// BulkEdit saves a grid of edits, one patch per row, in one transaction.
func (h *PlanHandler) BulkEdit(c *gin.Context) {
	var req BulkEditPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	saved, rowErrors, err := h.service().BulkEdit(h.ctx(c), req.Items)
	if err != nil {
		h.fail(c, err, "Failed to save the plans")
		return
	}
	// Row errors are a 422 with the rows in it, so the grid can put each
	// message back on the line it belongs to rather than showing one banner.
	if len(rowErrors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Some rows could not be saved",
			"details": gin.H{"rows": rowErrors},
		}})
		return
	}
	if len(saved) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    gin.H{"affected": 0, "requested": len(req.Items)},
			"message": "Nothing to do",
		})
		return
	}

	noun := "plans"
	if len(saved) == 1 {
		noun = "plan"
	}
	summary := "edited " + strconv.Itoa(len(saved)) + " " + noun + " in the grid"
	events.Emitted(c, "plans", "Plan", "bulk", saved[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(saved), "requested": len(req.Items)},
		"message": strconv.Itoa(len(saved)) + " " + noun + " updated",
	})
}

// BulkCreatePlanRequest is the grid's insert: whole rows, no ids.
type BulkCreatePlanRequest struct {
	Items []map[string]interface{} `json:"items" binding:"required,min=1,max=500"`
}

// BulkCreate inserts a grid of new plans in one transaction.
func (h *PlanHandler) BulkCreate(c *gin.Context) {
	var req BulkCreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	created, rowErrors, err := h.service().BulkCreate(h.ctx(c), req.Items)
	if err != nil {
		h.fail(c, err, "Failed to create the plans")
		return
	}
	if len(rowErrors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Some rows could not be created",
			"details": gin.H{"rows": rowErrors},
		}})
		return
	}
	if len(created) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    gin.H{"affected": 0, "requested": len(req.Items)},
			"message": "Nothing to do",
		})
		return
	}

	noun := "plans"
	if len(created) == 1 {
		noun = "plan"
	}
	summary := "created " + strconv.Itoa(len(created)) + " " + noun + " in the grid"
	events.Emitted(c, "plans", "Plan", "bulk", created[0].ID, summary, summary, nil, nil)

	// The rows as stored, so the grid can show the ids, the generated slugs and
	// anything else a hook filled in rather than guessing at them.
	c.JSON(http.StatusCreated, gin.H{
		"data":    created,
		"message": strconv.Itoa(len(created)) + " " + noun + " created",
	})
}
