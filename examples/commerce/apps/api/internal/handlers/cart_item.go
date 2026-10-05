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
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
	"commerce/apps/api/internal/paginate"
	"commerce/apps/api/internal/pdf"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

// CartItemHandler serves the cartitem endpoints. It reads the request, asks
// services.CartItemService, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type CartItemHandler struct {
	DB *gorm.DB
	// AppName brands the PDF export. The route file sets it from config.
	AppName string
}

// service is the cartitem service over this handler's database.
func (h *CartItemHandler) service() *services.CartItemService {
	return &services.CartItemService{DB: h.DB}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *CartItemHandler) ctx(c *gin.Context) context.Context {
	return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *CartItemHandler) fail(c *gin.Context, err error, fallback string) {
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
		respond.Fail(c, respond.CodeNotFound, "CartItem not found")
	default:
		respond.WriteError(c, err, fallback)
	}
}

// List returns a paginated list of cart_items.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *CartItemHandler) List(c *gin.Context) {
	res, err := h.service().List(h.ctx(c), paginate.Bind(c).With("cart_id", c.Query("cart_id")).With("product_id", c.Query("product_id")), c.Query("archived"))
	if err != nil {
		h.fail(c, err, "Failed to fetch cart_items")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX: the same
// search as List, without pages.
//
//	GET /api/cart_items/export?format=csv
//	GET /api/cart_items/export?format=xlsx&search=foo
func (h *CartItemHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	opts := export.Options{
		Sheet: "CartItems",
		Columns: []export.Column{
			{Header: "ID", Field: "ID"},
			{Header: "Variant ID", Field: "VariantID"},
			{Header: "Quantity", Field: "Quantity"},
			{Header: "UnitPrice", Field: "UnitPrice"},
			{Header: "Title", Field: "Title"},
			{Header: "VariantLabel", Field: "VariantLabel"},
			{Header: "ImageURL", Field: "ImageURL"},
			{Header: "Created At", Field: "CreatedAt", Format: "date:2006-01-02"},
		},
	}

	if format == "xlsx" {
		// Rows go into the workbook batch by batch, through a writer that moves
		// them to a temporary file, so memory stays flat however many match.
		sheet, err := export.NewXLSXStream(opts)
		if err != nil {
			h.fail(c, err, "Failed to export cart_items")
			return
		}
		defer func() {
			if err := sheet.Close(); err != nil {
				log.Printf("export cart_items as xlsx: removing temporary files: %v", err)
			}
		}()
		if err := h.service().Export(h.ctx(c), search, func(rows []models.CartItem) error {
			return sheet.Rows(rows)
		}); err != nil {
			h.fail(c, err, "Failed to export cart_items")
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", `attachment; filename="cart_items.xlsx"`)
		if err := sheet.Finish(c.Writer); err != nil {
			log.Printf("export cart_items as xlsx: %v", err)
		}
		return
	}

	// CSV streams: the header with the first batch, then each batch as it is
	// read.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="cart_items.csv"`)

	headerWritten := false
	err := h.service().Export(h.ctx(c), search, func(rows []models.CartItem) error {
		if !headerWritten {
			headerWritten = true
			return export.CSV(c.Writer, rows, opts)
		}
		return export.CSVRows(c.Writer, rows, opts)
	})
	if err == nil && !headerWritten {
		// Nothing matched: still a CSV, with its header row.
		err = export.CSV(c.Writer, []models.CartItem{}, opts)
	}
	if err != nil {
		if !c.Writer.Written() {
			c.Writer.Header().Del("Content-Type")
			c.Writer.Header().Del("Content-Disposition")
			h.fail(c, err, "Failed to export cart_items")
			return
		}
		// Rows are already on the wire, so the status cannot change. Logged,
		// so a truncated file has an explanation somewhere.
		log.Printf("export cart_items: %v", err)
	}
}

// GetByID returns a single cartitem by ID.
func (h *CartItemHandler) GetByID(c *gin.Context) {
	item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to load cartitem")
		return
	}

	// The version to send back as If-Match when saving.
	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this cartitem as a print-ready PDF: a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *CartItemHandler) PDF(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service().GetByID(h.ctx(c), id)
	if err != nil {
		h.fail(c, err, "Failed to load cartitem")
		return
	}

	appName := h.AppName
	if appName == "" {
		appName = "CartItem"
	}

	rec := pdf.Record{
		Title:      "CART ITEM",
		Subtitle:   pdf.Value(item.Title),
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
			{Label: "Cart", Value: pdf.Display(item.Cart)},
			{Label: "Product", Value: pdf.Display(item.Product)},
			{Label: "Variant ID", Value: pdf.Value(item.VariantID)},
			{Label: "Quantity", Value: pdf.Value(item.Quantity)},
			{Label: "Unit Price", Value: pdf.Value(item.UnitPrice)},
			{Label: "Title", Value: pdf.Value(item.Title)},
			{Label: "Variant Label", Value: pdf.Value(item.VariantLabel)},
			{Label: "Image URL", Value: pdf.Value(item.ImageURL)},
			{Label: "Created", Value: pdf.Value(item.CreatedAt)},
		},
	}

	out, err := pdf.RenderRecord(rec)
	if err != nil {
		respond.Fail(c, respond.CodePDFError, "could not render the PDF")
		return
	}

	filename := "cart-item-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// CreateCartItemRequest is the JSON body accepted by POST /cart_items.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type. routes.go passes
// this type to docs.Route(...).RequestBody().
type CreateCartItemRequest struct {
	CartID       string      `json:"cart_id" binding:"required"`
	ProductID    string      `json:"product_id" binding:"required"`
	VariantID    string      `json:"variant_id" binding:"required"`
	Quantity     int         `json:"quantity"`
	UnitPrice    money.Money `json:"unit_price"`
	Title        string      `json:"title" binding:"required"`
	VariantLabel string      `json:"variant_label" binding:"required"`
	ImageURL     string      `json:"image_url" binding:"required"`
}

// UpdateCartItemRequest is the JSON body accepted by PUT /cart_items/:id.
// Every field is optional: only what the client sends is applied.
type UpdateCartItemRequest struct {
	CartID       *string      `json:"cart_id"`
	ProductID    *string      `json:"product_id"`
	VariantID    string       `json:"variant_id"`
	Quantity     *int         `json:"quantity"`
	UnitPrice    *money.Money `json:"unit_price"`
	Title        string       `json:"title"`
	VariantLabel string       `json:"variant_label"`
	ImageURL     string       `json:"image_url"`
}

// Create adds a new cartitem.
func (h *CartItemHandler) Create(c *gin.Context) {
	var req CreateCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item := models.CartItem{
		CartID:       req.CartID,
		ProductID:    req.ProductID,
		VariantID:    req.VariantID,
		Quantity:     req.Quantity,
		UnitPrice:    req.UnitPrice,
		Title:        req.Title,
		VariantLabel: req.VariantLabel,
		ImageURL:     req.ImageURL,
	}

	if err := h.service().Create(h.ctx(c), &item); err != nil {
		h.fail(c, err, "Failed to create cartitem")
		return
	}

	events.Emitted(c, "cart_items", "CartItem", "created", item.ID, item.Title, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "CartItem created successfully",
	})
}

// Update modifies an existing cartitem. A client that sends the version it
// read as If-Match gets 409 instead of overwriting a newer save.
func (h *CartItemHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.CartID != nil {
		updates["cart_id"] = *req.CartID
	}
	if req.ProductID != nil {
		updates["product_id"] = *req.ProductID
	}
	if req.VariantID != "" {
		updates["variant_id"] = req.VariantID
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.UnitPrice != nil {
		updates["unit_price"] = *req.UnitPrice
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.VariantLabel != "" {
		updates["variant_label"] = req.VariantLabel
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}

	item, err := h.service().Update(h.ctx(c), id, updates, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to update cartitem")
		return
	}

	events.Emitted(c, "cart_items", "CartItem", "updated", item.ID, item.Title, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "CartItem updated successfully",
	})
}

// Patch applies a partial update to a cartitem. Used by the admin's grouped
// update view: each form group's Save button sends only the fields it owns, so
// editing "Address" does not rewrite "Pricing".
func (h *CartItemHandler) Patch(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	item, updates, err := h.service().Patch(h.ctx(c), c.Param("id"), body, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to patch cartitem")
		return
	}

	events.Emitted(c, "cart_items", "CartItem", "updated", item.ID, item.Title, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "CartItem updated successfully",
	})
}

// Delete soft-deletes a cartitem.
func (h *CartItemHandler) Delete(c *gin.Context) {
	item, err := h.service().Delete(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to delete cartitem")
		return
	}

	events.Emitted(c, "cart_items", "CartItem", "deleted", item.ID, item.Title, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "CartItem deleted successfully",
	})
}

// BulkCartItemRequest is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails.
type BulkCartItemRequest struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string `json:"action" binding:"required,oneof=delete archive restore patch"`
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string `json:"ids" binding:"required,min=1,max=500"`
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} `json:"patch"`
}

// Bulk applies one action to many cart_items in a single transaction.
func (h *CartItemHandler) Bulk(c *gin.Context) {
	var req BulkCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	result, err := h.service().Bulk(h.ctx(c), req.Action, req.IDs, req.Patch)
	if err != nil {
		h.fail(c, err, "Failed to "+req.Action+" cart_items")
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

	noun := "cart_items"
	if len(ids) == 1 {
		noun = "cartitem"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(result.Updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and the
	// count lives in the summary, where it can be read.
	events.Emitted(c, "cart_items", "CartItem", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}
