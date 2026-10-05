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

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/concurrency"
	"commerce/apps/api/internal/events"
	"commerce/apps/api/internal/export"
	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
	"commerce/apps/api/internal/pdf"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
	"commerce/apps/api/internal/storage"
)

// CollectionHandler serves the collection endpoints. It reads the request, asks
// services.CollectionService, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type CollectionHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	// AppName brands the PDF export. The route file sets it from config.
	AppName string
}

// service is the collection service over this handler's database.
func (h *CollectionHandler) service() *services.CollectionService {
	return &services.CollectionService{DB: h.DB, Storage: h.Storage}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *CollectionHandler) ctx(c *gin.Context) context.Context {
	return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *CollectionHandler) fail(c *gin.Context, err error, fallback string) {
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
		respond.Fail(c, respond.CodeNotFound, "Collection not found")
	default:
		respond.WriteError(c, err, fallback)
	}
}

// List returns a paginated list of collections.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *CollectionHandler) List(c *gin.Context) {
	res, err := h.service().List(h.ctx(c), paginate.Bind(c), c.Query("archived"))
	if err != nil {
		h.fail(c, err, "Failed to fetch collections")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX: the same
// search as List, without pages.
//
//	GET /api/collections/export?format=csv
//	GET /api/collections/export?format=xlsx&search=foo
func (h *CollectionHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	opts := export.Options{
		Sheet: "Collections",
		Columns: []export.Column{
			{Header: "ID", Field: "ID"},
			{Header: "Title", Field: "Title"},
			{Header: "Handle", Field: "Handle"},
			{Header: "Description", Field: "Description"},
			{Header: "Image", Field: "Image"},
			{Header: "Created At", Field: "CreatedAt", Format: "date:2006-01-02"},
		},
	}

	if format == "xlsx" {
		// Rows go into the workbook batch by batch, through a writer that moves
		// them to a temporary file, so memory stays flat however many match.
		sheet, err := export.NewXLSXStream(opts)
		if err != nil {
			h.fail(c, err, "Failed to export collections")
			return
		}
		defer func() {
			if err := sheet.Close(); err != nil {
				log.Printf("export collections as xlsx: removing temporary files: %v", err)
			}
		}()
		if err := h.service().Export(h.ctx(c), search, func(rows []models.Collection) error {
			return sheet.Rows(rows)
		}); err != nil {
			h.fail(c, err, "Failed to export collections")
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", `attachment; filename="collections.xlsx"`)
		if err := sheet.Finish(c.Writer); err != nil {
			log.Printf("export collections as xlsx: %v", err)
		}
		return
	}

	// CSV streams: the header with the first batch, then each batch as it is
	// read.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="collections.csv"`)

	headerWritten := false
	err := h.service().Export(h.ctx(c), search, func(rows []models.Collection) error {
		if !headerWritten {
			headerWritten = true
			return export.CSV(c.Writer, rows, opts)
		}
		return export.CSVRows(c.Writer, rows, opts)
	})
	if err == nil && !headerWritten {
		// Nothing matched: still a CSV, with its header row.
		err = export.CSV(c.Writer, []models.Collection{}, opts)
	}
	if err != nil {
		if !c.Writer.Written() {
			c.Writer.Header().Del("Content-Type")
			c.Writer.Header().Del("Content-Disposition")
			h.fail(c, err, "Failed to export collections")
			return
		}
		// Rows are already on the wire, so the status cannot change. Logged,
		// so a truncated file has an explanation somewhere.
		log.Printf("export collections: %v", err)
	}
}

// GetByID returns a single collection by ID.
func (h *CollectionHandler) GetByID(c *gin.Context) {
	item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to load collection")
		return
	}

	// The version to send back as If-Match when saving.
	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this collection as a print-ready PDF: a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *CollectionHandler) PDF(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service().GetByID(h.ctx(c), id)
	if err != nil {
		h.fail(c, err, "Failed to load collection")
		return
	}

	appName := h.AppName
	if appName == "" {
		appName = "Collection"
	}

	rec := pdf.Record{
		Title:      "COLLECTION",
		Subtitle:   pdf.Value(item.Title),
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
			{Label: "Title", Value: pdf.Value(item.Title)},
			{Label: "Description", Value: pdf.Value(item.Description)},
			{Label: "Created", Value: pdf.Value(item.CreatedAt)},
		},
	}

	out, err := pdf.RenderRecord(rec)
	if err != nil {
		respond.Fail(c, respond.CodePDFError, "could not render the PDF")
		return
	}

	filename := "collection-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// CreateCollectionRequest is the JSON body accepted by POST /collections.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type. routes.go passes
// this type to docs.Route(...).RequestBody().
type CreateCollectionRequest struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	Image       *files.FileRef `json:"image"`
}

// UpdateCollectionRequest is the JSON body accepted by PUT /collections/:id.
// Every field is optional: only what the client sends is applied.
type UpdateCollectionRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Image       **files.FileRef `json:"image"`
}

// Create adds a new collection.
func (h *CollectionHandler) Create(c *gin.Context) {
	var req CreateCollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item := models.Collection{
		Title:       req.Title,
		Description: req.Description,
		Image:       req.Image,
	}

	if err := h.service().Create(h.ctx(c), &item); err != nil {
		h.fail(c, err, "Failed to create collection")
		return
	}

	events.Emitted(c, "collections", "Collection", "created", item.ID, item.Title, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "Collection created successfully",
	})
}

// Update modifies an existing collection. A client that sends the version it
// read as If-Match gets 409 instead of overwriting a newer save.
func (h *CollectionHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}

	item, err := h.service().Update(h.ctx(c), id, updates, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to update collection")
		return
	}

	events.Emitted(c, "collections", "Collection", "updated", item.ID, item.Title, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Collection updated successfully",
	})
}

// Patch applies a partial update to a collection. Used by the admin's grouped
// update view: each form group's Save button sends only the fields it owns, so
// editing "Address" does not rewrite "Pricing".
func (h *CollectionHandler) Patch(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item, updates, err := h.service().Patch(h.ctx(c), c.Param("id"), body, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to patch collection")
		return
	}

	events.Emitted(c, "collections", "Collection", "updated", item.ID, item.Title, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Collection updated successfully",
	})
}

// Delete soft-deletes a collection.
func (h *CollectionHandler) Delete(c *gin.Context) {
	item, err := h.service().Delete(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to delete collection")
		return
	}

	events.Emitted(c, "collections", "Collection", "deleted", item.ID, item.Title, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Collection deleted successfully",
	})
}

// BulkCollectionRequest is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails.
type BulkCollectionRequest struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string `json:"action" binding:"required,oneof=delete archive restore patch"`
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string `json:"ids" binding:"required,min=1,max=500"`
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} `json:"patch"`
}

// Bulk applies one action to many collections in a single transaction.
func (h *CollectionHandler) Bulk(c *gin.Context) {
	var req BulkCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	result, err := h.service().Bulk(h.ctx(c), req.Action, req.IDs, req.Patch)
	if err != nil {
		h.fail(c, err, "Failed to "+req.Action+" collections")
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

	noun := "collections"
	if len(ids) == 1 {
		noun = "collection"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(result.Updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and the
	// count lives in the summary, where it can be read.
	events.Emitted(c, "collections", "Collection", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}
