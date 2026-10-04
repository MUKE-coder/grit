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

	"library/apps/api/internal/authz"
	"library/apps/api/internal/concurrency"
	"library/apps/api/internal/events"
	"library/apps/api/internal/export"
	"library/apps/api/internal/models"
	"library/apps/api/internal/paginate"
	"library/apps/api/internal/pdf"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

// AuthorHandler serves the author endpoints. It reads the request, asks
// services.AuthorService, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type AuthorHandler struct {
	DB *gorm.DB
	// AppName brands the PDF export. The route file sets it from config.
	AppName string
}

// service is the author service over this handler's database.
func (h *AuthorHandler) service() *services.AuthorService {
	return &services.AuthorService{DB: h.DB}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *AuthorHandler) ctx(c *gin.Context) context.Context {
	return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *AuthorHandler) fail(c *gin.Context, err error, fallback string) {
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
		respond.Fail(c, respond.CodeNotFound, "Author not found")
	default:
		respond.WriteError(c, err, fallback)
	}
}

// List returns a paginated list of authors.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *AuthorHandler) List(c *gin.Context) {
	res, err := h.service().List(h.ctx(c), paginate.Bind(c), c.Query("archived"))
	if err != nil {
		h.fail(c, err, "Failed to fetch authors")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX: the same
// search as List, without pages.
//
//	GET /api/authors/export?format=csv
//	GET /api/authors/export?format=xlsx&search=foo
func (h *AuthorHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	opts := export.Options{
		Sheet: "Authors",
		Columns: []export.Column{
			{Header: "ID", Field: "ID"},
			{Header: "Name", Field: "Name"},
			{Header: "Bio", Field: "Bio"},
			{Header: "Website", Field: "Website"},
			{Header: "Created At", Field: "CreatedAt", Format: "date:2006-01-02"},
		},
	}

	if format == "xlsx" {
		// Rows go into the workbook batch by batch, through a writer that moves
		// them to a temporary file, so memory stays flat however many match.
		sheet, err := export.NewXLSXStream(opts)
		if err != nil {
			h.fail(c, err, "Failed to export authors")
			return
		}
		defer func() {
			if err := sheet.Close(); err != nil {
				log.Printf("export authors as xlsx: removing temporary files: %v", err)
			}
		}()
		if err := h.service().Export(h.ctx(c), search, func(rows []models.Author) error {
			return sheet.Rows(rows)
		}); err != nil {
			h.fail(c, err, "Failed to export authors")
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", `attachment; filename="authors.xlsx"`)
		if err := sheet.Finish(c.Writer); err != nil {
			log.Printf("export authors as xlsx: %v", err)
		}
		return
	}

	// CSV streams: the header with the first batch, then each batch as it is
	// read.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="authors.csv"`)

	headerWritten := false
	err := h.service().Export(h.ctx(c), search, func(rows []models.Author) error {
		if !headerWritten {
			headerWritten = true
			return export.CSV(c.Writer, rows, opts)
		}
		return export.CSVRows(c.Writer, rows, opts)
	})
	if err == nil && !headerWritten {
		// Nothing matched: still a CSV, with its header row.
		err = export.CSV(c.Writer, []models.Author{}, opts)
	}
	if err != nil {
		if !c.Writer.Written() {
			c.Writer.Header().Del("Content-Type")
			c.Writer.Header().Del("Content-Disposition")
			h.fail(c, err, "Failed to export authors")
			return
		}
		// Rows are already on the wire, so the status cannot change. Logged,
		// so a truncated file has an explanation somewhere.
		log.Printf("export authors: %v", err)
	}
}

// GetByID returns a single author by ID.
func (h *AuthorHandler) GetByID(c *gin.Context) {
	item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to load author")
		return
	}

	// The version to send back as If-Match when saving.
	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this author as a print-ready PDF: a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *AuthorHandler) PDF(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service().GetByID(h.ctx(c), id)
	if err != nil {
		h.fail(c, err, "Failed to load author")
		return
	}

	appName := h.AppName
	if appName == "" {
		appName = "Author"
	}

	rec := pdf.Record{
		Title:      "AUTHOR",
		Subtitle:   pdf.Value(item.Name),
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
			{Label: "Name", Value: pdf.Value(item.Name)},
			{Label: "Bio", Value: pdf.Value(item.Bio)},
			{Label: "Website", Value: pdf.Value(item.Website)},
			{Label: "Created", Value: pdf.Value(item.CreatedAt)},
		},
	}

	out, err := pdf.RenderRecord(rec)
	if err != nil {
		respond.Fail(c, respond.CodePDFError, "could not render the PDF")
		return
	}

	filename := "author-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// CreateAuthorRequest is the JSON body accepted by POST /authors.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type. routes.go passes
// this type to docs.Route(...).RequestBody().
type CreateAuthorRequest struct {
	Name    string `json:"name" binding:"required"`
	Bio     string `json:"bio"`
	Website string `json:"website"`
}

// UpdateAuthorRequest is the JSON body accepted by PUT /authors/:id.
// Every field is optional: only what the client sends is applied.
type UpdateAuthorRequest struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Website string `json:"website"`
}

// Create adds a new author.
func (h *AuthorHandler) Create(c *gin.Context) {
	var req CreateAuthorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item := models.Author{
		Name:    req.Name,
		Bio:     req.Bio,
		Website: req.Website,
	}

	if err := h.service().Create(h.ctx(c), &item); err != nil {
		h.fail(c, err, "Failed to create author")
		return
	}

	events.Emitted(c, "authors", "Author", "created", item.ID, item.Name, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "Author created successfully",
	})
}

// Update modifies an existing author. A client that sends the version it
// read as If-Match gets 409 instead of overwriting a newer save.
func (h *AuthorHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAuthorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}

	item, err := h.service().Update(h.ctx(c), id, updates, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to update author")
		return
	}

	events.Emitted(c, "authors", "Author", "updated", item.ID, item.Name, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Author updated successfully",
	})
}

// Patch applies a partial update to a author. Used by the admin's grouped
// update view: each form group's Save button sends only the fields it owns, so
// editing "Address" does not rewrite "Pricing".
func (h *AuthorHandler) Patch(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item, updates, err := h.service().Patch(h.ctx(c), c.Param("id"), body, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to patch author")
		return
	}

	events.Emitted(c, "authors", "Author", "updated", item.ID, item.Name, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Author updated successfully",
	})
}

// Delete soft-deletes a author.
func (h *AuthorHandler) Delete(c *gin.Context) {
	item, err := h.service().Delete(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to delete author")
		return
	}

	events.Emitted(c, "authors", "Author", "deleted", item.ID, item.Name, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Author deleted successfully",
	})
}

// BulkAuthorRequest is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails.
type BulkAuthorRequest struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string `json:"action" binding:"required,oneof=delete archive restore patch"`
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string `json:"ids" binding:"required,min=1,max=500"`
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} `json:"patch"`
}

// Bulk applies one action to many authors in a single transaction.
func (h *AuthorHandler) Bulk(c *gin.Context) {
	var req BulkAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	result, err := h.service().Bulk(h.ctx(c), req.Action, req.IDs, req.Patch)
	if err != nil {
		h.fail(c, err, "Failed to "+req.Action+" authors")
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

	noun := "authors"
	if len(ids) == 1 {
		noun = "author"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(result.Updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and the
	// count lives in the summary, where it can be read.
	events.Emitted(c, "authors", "Author", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}
