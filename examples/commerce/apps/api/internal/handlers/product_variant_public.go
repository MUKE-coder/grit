package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
	"commerce/apps/api/internal/respond"
)

// The public product variant surface.
//
// One endpoint, because a picker needs one payload. Fetching the options, then
// the variants, then a price on every swatch click is three round trips per
// interaction on the busiest page in the shop, and the figure it lands on would
// still be a client's own arithmetic rather than the server's.

// publicProductVariant is one buyable combination.
//
// Stock is published as a boolean and never as a count, the rule the rest of
// the public surface follows: "in stock" is what the page renders, and the
// number is a business fact competitors enjoy.
//
// OptionValueIDs rather than nested values, because the values are already in
// the options list and a picker matches a selection to a variant by comparing
// ids. Nesting them would send the same objects twice.
type publicProductVariant struct {
	ID             string         `json:"id"`
	SKU            string         `json:"sku,omitempty"`
	Price          money.Money    `json:"price"`
	InStock        bool           `json:"in_stock"`
	Images         files.FileRefs `json:"images,omitempty"`
	OptionValueIDs []string       `json:"option_value_ids"`
}

// ListPublic handles GET /api/v1/public/products/:key/variants.
//
// Looked up by handle or id, so one URL works from a catalogue
// link and from a relation the client had already resolved.
//
// Returns the options to draw, the combinations to match a selection against,
// and the price range a listing card needs. A product with no variants
// gets empty lists and a range of its own price, which is what lets a storefront
// render one component whether or not variants were ever set up.
func (h *ProductVariantHandler) ListPublic(c *gin.Context) {
	var product models.Product
	if err := h.DB.WithContext(c.Request.Context()).Where("handle = ? OR id = ?", c.Param("key"), c.Param("key")).Where("archived_at IS NULL").First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": "Product not found",
		}})
		return
	}

	options, err := h.Variants.OptionsFor(product.ID)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to load the options")
		return
	}
	byID := make(map[string]models.Option, len(options))
	for _, option := range options {
		byID[option.ID] = option
	}

	variants, err := h.Variants.VariantsFor(product.ID)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to load the variants")
		return
	}

	// Inactive combinations are dropped rather than marked. A variant somebody
	// switched off is not something the shop sells, and publishing it greyed out
	// invites a client to render a choice that can never be completed.
	published := make([]models.ProductVariant, 0, len(variants))
	out := make([]publicProductVariant, 0, len(variants))
	for _, variant := range variants {
		if !variant.Active {
			continue
		}
		published = append(published, variant)

		valueIDs := make([]string, 0, len(variant.OptionValues))
		for _, value := range variant.OptionValues {
			valueIDs = append(valueIDs, value.ID)
		}
		out = append(out, publicProductVariant{
			ID:             variant.ID,
			SKU:            variant.SKU,
			Price:          h.Variants.ResolvePrice(product.Price, variant, byID),
			InStock:        variant.InStock(),
			Images:         variant.Images,
			OptionValueIDs: valueIDs,
		})
	}

	low, high := h.Variants.PriceRange(product.Price, published, byID)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"options":  toPublicOptions(options),
		"variants": out,
		"price_range": gin.H{
			"low":  low,
			"high": high,
			// True when every buyable combination costs the same, which is the
			// question a listing card asks before choosing between "49" and
			// "from 49".
			"single": low == high,
		},
	}})
}
