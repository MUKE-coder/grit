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
	"library/apps/api/internal/files"
	"library/apps/api/internal/jsontime"
	"library/apps/api/internal/models"
	"library/apps/api/internal/money"
	"library/apps/api/internal/paginate"
	"library/apps/api/internal/pdf"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
	"library/apps/api/internal/storage"
)

// BookHandler serves the book endpoints. It reads the request, asks
// services.BookService, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type BookHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	// AppName brands the PDF export. The route file sets it from config.
	AppName string
}

// service is the book service over this handler's database.
func (h *BookHandler) service() *services.BookService {
	return &services.BookService{DB: h.DB, Storage: h.Storage}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *BookHandler) ctx(c *gin.Context) context.Context {
	return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *BookHandler) fail(c *gin.Context, err error, fallback string) {
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
		respond.Fail(c, respond.CodeNotFound, "Book not found")
	default:
		respond.WriteError(c, err, fallback)
	}
}

// List returns a paginated list of books.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *BookHandler) List(c *gin.Context) {
	res, err := h.service().List(h.ctx(c), paginate.Bind(c).With("author_id", c.Query("author_id")), c.Query("archived"))
	if err != nil {
		h.fail(c, err, "Failed to fetch books")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX: the same
// search as List, without pages.
//
//	GET /api/books/export?format=csv
//	GET /api/books/export?format=xlsx&search=foo
func (h *BookHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	opts := export.Options{
		Sheet: "Books",
		Columns: []export.Column{
			{Header: "ID", Field: "ID"},
			{Header: "Title", Field: "Title"},
			{Header: "Isbn", Field: "Isbn"},
			{Header: "Summary", Field: "Summary"},
			{Header: "Published", Field: "Published"},
			{Header: "Price", Field: "Price"},
			{Header: "Cover", Field: "Cover"},
			{Header: "Genre", Field: "Genre"},
			{Header: "Created At", Field: "CreatedAt", Format: "date:2006-01-02"},
		},
	}

	if format == "xlsx" {
		// Rows go into the workbook batch by batch, through a writer that moves
		// them to a temporary file, so memory stays flat however many match.
		sheet, err := export.NewXLSXStream(opts)
		if err != nil {
			h.fail(c, err, "Failed to export books")
			return
		}
		defer func() {
			if err := sheet.Close(); err != nil {
				log.Printf("export books as xlsx: removing temporary files: %v", err)
			}
		}()
		if err := h.service().Export(h.ctx(c), search, func(rows []models.Book) error {
			return sheet.Rows(rows)
		}); err != nil {
			h.fail(c, err, "Failed to export books")
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", `attachment; filename="books.xlsx"`)
		if err := sheet.Finish(c.Writer); err != nil {
			log.Printf("export books as xlsx: %v", err)
		}
		return
	}

	// CSV streams: the header with the first batch, then each batch as it is
	// read.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="books.csv"`)

	headerWritten := false
	err := h.service().Export(h.ctx(c), search, func(rows []models.Book) error {
		if !headerWritten {
			headerWritten = true
			return export.CSV(c.Writer, rows, opts)
		}
		return export.CSVRows(c.Writer, rows, opts)
	})
	if err == nil && !headerWritten {
		// Nothing matched: still a CSV, with its header row.
		err = export.CSV(c.Writer, []models.Book{}, opts)
	}
	if err != nil {
		if !c.Writer.Written() {
			c.Writer.Header().Del("Content-Type")
			c.Writer.Header().Del("Content-Disposition")
			h.fail(c, err, "Failed to export books")
			return
		}
		// Rows are already on the wire, so the status cannot change. Logged,
		// so a truncated file has an explanation somewhere.
		log.Printf("export books: %v", err)
	}
}

// GetByID returns a single book by ID.
func (h *BookHandler) GetByID(c *gin.Context) {
	item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to load book")
		return
	}

	// The version to send back as If-Match when saving.
	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this book as a print-ready PDF: a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *BookHandler) PDF(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service().GetByID(h.ctx(c), id)
	if err != nil {
		h.fail(c, err, "Failed to load book")
		return
	}

	appName := h.AppName
	if appName == "" {
		appName = "Book"
	}

	rec := pdf.Record{
		Title:      "BOOK",
		Subtitle:   pdf.Value(item.Title),
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
			{Label: "Title", Value: pdf.Value(item.Title)},
			{Label: "Isbn", Value: pdf.Value(item.Isbn)},
			{Label: "Published", Value: pdf.Value(item.Published)},
			{Label: "Price", Value: pdf.Value(item.Price)},
			{Label: "Genre", Value: pdf.Value(item.Genre)},
			{Label: "Author", Value: pdf.Display(item.Author)},
			{Label: "Created", Value: pdf.Value(item.CreatedAt)},
		},
	}

	out, err := pdf.RenderRecord(rec)
	if err != nil {
		respond.Fail(c, respond.CodePDFError, "could not render the PDF")
		return
	}

	filename := "book-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// CreateBookRequest is the JSON body accepted by POST /books.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type. routes.go passes
// this type to docs.Route(...).RequestBody().
type CreateBookRequest struct {
	Title     string         `json:"title" binding:"required"`
	Isbn      string         `json:"isbn" binding:"required"`
	Summary   string         `json:"summary"`
	Published *jsontime.Date `json:"published"`
	Price     money.Money    `json:"price"`
	Cover     *files.FileRef `json:"cover"`
	Genre     string         `json:"genre" binding:"required"`
	AuthorID  string         `json:"author_id" binding:"required"`
}

// UpdateBookRequest is the JSON body accepted by PUT /books/:id.
// Every field is optional: only what the client sends is applied.
type UpdateBookRequest struct {
	Title     string          `json:"title"`
	Isbn      string          `json:"isbn"`
	Summary   string          `json:"summary"`
	Published **jsontime.Date `json:"published"`
	Price     *money.Money    `json:"price"`
	Cover     **files.FileRef `json:"cover"`
	Genre     string          `json:"genre"`
	AuthorID  *string         `json:"author_id"`
}

// Create adds a new book.
func (h *BookHandler) Create(c *gin.Context) {
	var req CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item := models.Book{
		Title:     req.Title,
		Isbn:      req.Isbn,
		Summary:   req.Summary,
		Published: req.Published,
		Price:     req.Price,
		Cover:     req.Cover,
		Genre:     req.Genre,
		AuthorID:  req.AuthorID,
	}

	if err := h.service().Create(h.ctx(c), &item); err != nil {
		h.fail(c, err, "Failed to create book")
		return
	}

	events.Emitted(c, "books", "Book", "created", item.ID, item.Title, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "Book created successfully",
	})
}

// Update modifies an existing book. A client that sends the version it
// read as If-Match gets 409 instead of overwriting a newer save.
func (h *BookHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Isbn != "" {
		updates["isbn"] = req.Isbn
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Published != nil {
		updates["published"] = *req.Published
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Cover != nil {
		updates["cover"] = *req.Cover
	}
	if req.Genre != "" {
		updates["genre"] = req.Genre
	}
	if req.AuthorID != nil {
		updates["author_id"] = *req.AuthorID
	}

	item, err := h.service().Update(h.ctx(c), id, updates, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to update book")
		return
	}

	events.Emitted(c, "books", "Book", "updated", item.ID, item.Title, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Book updated successfully",
	})
}

// Patch applies a partial update to a book. Used by the admin's grouped
// update view: each form group's Save button sends only the fields it owns, so
// editing "Address" does not rewrite "Pricing".
func (h *BookHandler) Patch(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item, updates, err := h.service().Patch(h.ctx(c), c.Param("id"), body, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to patch book")
		return
	}

	events.Emitted(c, "books", "Book", "updated", item.ID, item.Title, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Book updated successfully",
	})
}

// Delete soft-deletes a book.
func (h *BookHandler) Delete(c *gin.Context) {
	item, err := h.service().Delete(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to delete book")
		return
	}

	events.Emitted(c, "books", "Book", "deleted", item.ID, item.Title, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}

// BulkBookRequest is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails.
type BulkBookRequest struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string `json:"action" binding:"required,oneof=delete archive restore patch"`
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string `json:"ids" binding:"required,min=1,max=500"`
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} `json:"patch"`
}

// Bulk applies one action to many books in a single transaction.
func (h *BookHandler) Bulk(c *gin.Context) {
	var req BulkBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	result, err := h.service().Bulk(h.ctx(c), req.Action, req.IDs, req.Patch)
	if err != nil {
		h.fail(c, err, "Failed to "+req.Action+" books")
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

	noun := "books"
	if len(ids) == 1 {
		noun = "book"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(result.Updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and the
	// count lives in the summary, where it can be read.
	events.Emitted(c, "books", "Book", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}
