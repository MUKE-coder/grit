package scaffold

import "strings"

// APIVariantHandlerGo emits the admin endpoints for options and variants.
func APIVariantHandlerGo(module, pascal, snake, plural string) string {
	return apiVariantHandlerGo(module, pascal, snake, plural)
}

func apiVariantHandlerGo(module, pascal, snake, plural string) string {
	lower := strings.ToLower(pascal)
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + module + `/internal/models"
	"` + module + `/internal/services"
)

// ` + pascal + `VariantHandler serves the option library and one ` + lower + `'s
// variant matrix.
//
// Options are shop-wide, so they are managed on their own endpoints rather than
// nested under a ` + lower + `. Nesting them would suggest each ` + lower + ` owns
// its colours, which is the per-product-options mistake the schema exists to
// avoid.
type ` + pascal + `VariantHandler struct {
	DB       *gorm.DB
	Variants *services.` + pascal + `VariantService
}

func New` + pascal + `VariantHandler(db *gorm.DB) *` + pascal + `VariantHandler {
	return &` + pascal + `VariantHandler{DB: db, Variants: services.New` + pascal + `VariantService(db)}
}

// ListOptions handles GET /api/v1/options: the whole library with its values.
func (h *` + pascal + `VariantHandler) ListOptions(c *gin.Context) {
	var options []models.Option
	err := h.DB.Preload("Values", func(db *gorm.DB) *gorm.DB {
		return db.Order("position asc, label asc")
	}).Order("position asc, name asc").Find(&options).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to load options",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": options})
}

// CreateOption handles POST /api/v1/options.
func (h *` + pascal + `VariantHandler) CreateOption(c *gin.Context) {
	var req struct {
		Name         string ` + "`" + `json:"name" binding:"required"` + "`" + `
		Kind         string ` + "`" + `json:"kind"` + "`" + `
		AffectsPrice bool   ` + "`" + `json:"affects_price"` + "`" + `
		Position     int    ` + "`" + `json:"position"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}
	if req.Kind == "" {
		req.Kind = "select"
	}

	option := models.Option{
		Name: req.Name, Kind: req.Kind,
		AffectsPrice: req.AffectsPrice, Position: req.Position,
	}
	if err := h.DB.Create(&option).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to create the option",
		}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": option, "message": "Option created"})
}

// CreateOptionValue handles POST /api/v1/options/:id/values.
func (h *` + pascal + `VariantHandler) CreateOptionValue(c *gin.Context) {
	var req struct {
		Label      string  ` + "`" + `json:"label" binding:"required"` + "`" + `
		Swatch     string  ` + "`" + `json:"swatch"` + "`" + `
		PriceDelta float64 ` + "`" + `json:"price_delta"` + "`" + `
		Position   int     ` + "`" + `json:"position"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}

	value := models.OptionValue{
		OptionID: c.Param("id"), Label: req.Label, Swatch: req.Swatch,
		PriceDelta: req.PriceDelta, Position: req.Position,
	}
	if err := h.DB.Create(&value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to add the value",
		}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": value, "message": "Value added"})
}

// DeleteOptionValue handles DELETE /api/v1/option-values/:id.
//
// Refuses while any variant still uses it. Deleting the value out from under a
// combination leaves a variant that means "Red, and something", which is not a
// thing anybody can buy or ship.
func (h *` + pascal + `VariantHandler) DeleteOptionValue(c *gin.Context) {
	var used int64
	h.DB.Table("variant_option_values").
		Where("option_value_id = ?", c.Param("id")).
		Count(&used)
	if used > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "VALUE_IN_USE",
			"message": "That value is part of existing variants. Delete those first.",
		}})
		return
	}

	if err := h.DB.Delete(&models.OptionValue{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to delete the value",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Value deleted"})
}

// SetOptions handles PUT /api/v1/` + plural + `/:id/options: which options this
// ` + lower + ` offers, in order.
func (h *` + pascal + `VariantHandler) SetOptions(c *gin.Context) {
	var req struct {
		OptionIDs []string ` + "`" + `json:"option_ids"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}

	` + snake + `ID := c.Param("id")
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Removing an option a variant is built from would leave that variant
		// describing a choice the ` + lower + ` no longer offers, so the
		// existing matrix goes with it.
		if err := tx.Where("` + snake + `_id = ?", ` + snake + `ID).
			Delete(&models.` + pascal + `Option{}).Error; err != nil {
			return err
		}
		for i, optionID := range req.OptionIDs {
			link := models.` + pascal + `Option{` + pascal + `ID: ` + snake + `ID, OptionID: optionID, Position: i}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to set the options",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Options set"})
}

// List handles GET /api/v1/` + plural + `/:id/variants, with prices resolved.
//
// The resolved price is included rather than left to the caller, because every
// caller computing it themselves is every caller getting a chance to compute it
// differently.
func (h *` + pascal + `VariantHandler) List(c *gin.Context) {
	` + snake + `ID := c.Param("id")

	var ` + snake + ` models.` + pascal + `
	if err := h.DB.Where("id = ?", ` + snake + `ID).First(&` + snake + `).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": "` + pascal + ` not found",
		}})
		return
	}

	options, err := h.Variants.OptionsFor(` + snake + `ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": err.Error(),
		}})
		return
	}
	byID := make(map[string]models.Option, len(options))
	for _, option := range options {
		byID[option.ID] = option
	}

	variants, err := h.Variants.VariantsFor(` + snake + `ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": err.Error(),
		}})
		return
	}

	type row struct {
		models.` + pascal + `Variant
		Price float64 ` + "`" + `json:"price"` + "`" + `
	}
	out := make([]row, 0, len(variants))
	for _, variant := range variants {
		out = append(out, row{
			` + pascal + `Variant: variant,
			Price:    h.Variants.ResolvePrice(` + snake + `.Price, variant, byID),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"options": options, "variants": out}})
}

// GenerateMatrix handles POST /api/v1/` + plural + `/:id/variants/generate.
func (h *` + pascal + `VariantHandler) GenerateMatrix(c *gin.Context) {
	var req struct {
		Limit int ` + "`" + `json:"limit"` + "`" + `
	}
	_ = c.ShouldBindJSON(&req)

	created, err := h.Variants.Generate(c.Param("id"), req.Limit)
	if err != nil {
		// The caller asked for something that cannot be done rather than the
		// server failing, and the message says which.
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
			"code": "CANNOT_GENERATE", "message": err.Error(),
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"created": created},
		"message": "Matrix generated",
	})
}

// Update handles PATCH /api/v1/variants/:id: sku, stock, price and images.
func (h *` + pascal + `VariantHandler) Update(c *gin.Context) {
	var req struct {
		SKU           *string  ` + "`" + `json:"sku"` + "`" + `
		Stock         *int     ` + "`" + `json:"stock"` + "`" + `
		PriceOverride *float64 ` + "`" + `json:"price_override"` + "`" + `
		Active        *bool    ` + "`" + `json:"active"` + "`" + `
		// ClearPrice removes an override. A nil price_override cannot mean
		// "clear it", because nil is also what a partial update sends for a
		// field it is not touching.
		ClearPrice bool ` + "`" + `json:"clear_price"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}

	updates := map[string]any{}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.ClearPrice {
		updates["price_override"] = nil
	} else if req.PriceOverride != nil {
		updates["price_override"] = *req.PriceOverride
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nothing to update"})
		return
	}

	err := h.DB.Model(&models.` + pascal + `Variant{}).
		Where("id = ?", c.Param("id")).
		Updates(updates).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to update the variant",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Variant updated"})
}
`
}
