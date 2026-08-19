package scaffold

import "strings"

// The public half of the variant system.
//
// The admin endpoints in api_variant_handler.go sit behind auth, which is right
// for the matrix editor and useless to a storefront: the page that has to draw
// the colour swatches has no logged-in user. So there is a second, narrower
// surface here, mounted in the public group beside the resource's own public
// routes and guarded by the same API key.
//
// It is narrower in the way the rest of the public surface is narrow. Stock
// counts do not go out, only whether a combination can be bought; inactive
// variants do not go out at all; and the response is a hand-written view rather
// than the model, so a column added to the variant table next month is private
// until somebody says otherwise.
//
// The one thing it publishes generously is price, resolved per variant and
// summarised as a range. That is not a leak, it is the point: a picker that
// cannot say what a combination costs until it is in the cart is the pattern
// this system exists to replace.

// APIOptionPublicGo emits internal/handlers/option_public.go: the published
// shape of an option and its values.
//
// Its own file, written once however many resources offer variants, for the
// same reason the option MODEL is: options are shop-wide. A per-resource copy
// of these types would also be a second declaration of publicOption in package
// handlers, which is a compile error the moment a project runs
// `grit add variants` twice.
func APIOptionPublicGo(module string) string {
	return `package handlers

import (
	"` + module + `/internal/files"
	"` + module + `/internal/models"
)

// The published shape of the option library.
//
// An allowlist rather than the model, the same default the rest of the public
// surface takes: a column added to options next month is private until somebody
// adds it here.

// publicOptionValue is one choice, as a storefront may see it.
type publicOptionValue struct {
	ID    string ` + "`" + `json:"id"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
	Slug  string ` + "`" + `json:"slug"` + "`" + `

	// Swatch is the CSS colour a colour picker paints. Empty for every other
	// kind of option.
	Swatch string ` + "`" + `json:"swatch,omitempty"` + "`" + `

	// Image is the swatch picture, not the photograph of the product in this
	// colour. That one belongs to the variant, because it is a picture of a
	// combination rather than of a value.
	Image *files.FileRef ` + "`" + `json:"image,omitempty"` + "`" + `

	// PriceDelta lets a picker label a choice "+ 20" before anything is
	// selected. Zero unless the option declares AffectsPrice, so the label can
	// never disagree with the price the server resolves.
	PriceDelta float64 ` + "`" + `json:"price_delta"` + "`" + `
}

// publicOption is one axis of choice, with its values in display order.
//
// Kind travels because it is how a storefront knows what to draw: "swatch" for
// colours, "size" for a row of boxes, "select" for a dropdown. Deciding that on
// the client from the option's name is how one shop ends up rendering Colour as
// a dropdown because somebody called it "Color".
type publicOption struct {
	ID           string              ` + "`" + `json:"id"` + "`" + `
	Name         string              ` + "`" + `json:"name"` + "`" + `
	Slug         string              ` + "`" + `json:"slug"` + "`" + `
	Kind         string              ` + "`" + `json:"kind"` + "`" + `
	AffectsPrice bool                ` + "`" + `json:"affects_price"` + "`" + `
	Values       []publicOptionValue ` + "`" + `json:"values"` + "`" + `
}

// toPublicOptions maps the library for the wire.
func toPublicOptions(options []models.Option) []publicOption {
	out := make([]publicOption, 0, len(options))
	for _, option := range options {
		values := make([]publicOptionValue, 0, len(option.Values))
		for _, value := range option.Values {
			values = append(values, publicOptionValue{
				ID:         value.ID,
				Label:      value.Label,
				Slug:       value.Slug,
				Swatch:     value.Swatch,
				Image:      value.Image,
				PriceDelta: deltaWhenPriced(option, value),
			})
		}
		out = append(out, publicOption{
			ID: option.ID, Name: option.Name, Slug: option.Slug,
			Kind: option.Kind, AffectsPrice: option.AffectsPrice,
			Values: values,
		})
	}
	return out
}

// deltaWhenPriced returns a value's delta only where its option declares that
// the axis affects price.
//
// The same guard ResolvePrice applies, repeated rather than trusted, because
// the two numbers are read side by side. Without it a picker labels a swatch
// "+ 20" from a delta typed by mistake and then resolves to the base price,
// having told the customer something untrue in the space of one click.
func deltaWhenPriced(option models.Option, value models.OptionValue) float64 {
	if !option.AffectsPrice {
		return 0
	}
	return value.PriceDelta
}
`
}

// APIVariantPublicGo emits internal/handlers/<snake>_variant_public.go.
//
// hasSlug and hasArchivedAt are read off the resource's model rather than
// assumed. A resource generated without a slug has no such column, and a lookup
// on it would be a SQL error rather than a 404.
func APIVariantPublicGo(module, pascal, snake, plural string, hasSlug, hasArchivedAt bool) string {
	return apiVariantPublicGo(module, pascal, snake, plural, hasSlug, hasArchivedAt)
}

func apiVariantPublicGo(module, pascal, snake, plural string, hasSlug, hasArchivedAt bool) string {
	lower := strings.ToLower(pascal)

	// A slug is what a public URL should carry, and an id is what a client has
	// when it followed a relation, so both work where both exist.
	lookup := `h.DB.Where("id = ?", c.Param("key"))`
	lookupNote := "Looked up by id."
	if hasSlug {
		lookup = `h.DB.Where("slug = ? OR id = ?", c.Param("key"), c.Param("key"))`
		lookupNote = "Looked up by slug or id, so one URL works from a catalogue\n" +
			"// link and from a relation the client had already resolved."
	}
	if hasArchivedAt {
		lookup += `.Where("archived_at IS NULL")`
	}

	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"` + module + `/internal/files"
	"` + module + `/internal/models"
)

// The public ` + lower + ` variant surface.
//
// One endpoint, because a picker needs one payload. Fetching the options, then
// the variants, then a price on every swatch click is three round trips per
// interaction on the busiest page in the shop, and the figure it lands on would
// still be a client's own arithmetic rather than the server's.

// public` + pascal + `Variant is one buyable combination.
//
// Stock is published as a boolean and never as a count, the rule the rest of
// the public surface follows: "in stock" is what the page renders, and the
// number is a business fact competitors enjoy.
//
// OptionValueIDs rather than nested values, because the values are already in
// the options list and a picker matches a selection to a variant by comparing
// ids. Nesting them would send the same objects twice.
type public` + pascal + `Variant struct {
	ID             string         ` + "`" + `json:"id"` + "`" + `
	SKU            string         ` + "`" + `json:"sku,omitempty"` + "`" + `
	Price          float64        ` + "`" + `json:"price"` + "`" + `
	InStock        bool           ` + "`" + `json:"in_stock"` + "`" + `
	Images         files.FileRefs ` + "`" + `json:"images,omitempty"` + "`" + `
	OptionValueIDs []string       ` + "`" + `json:"option_value_ids"` + "`" + `
}

// ListPublic handles GET /api/v1/public/` + plural + `/:key/variants.
//
// ` + lookupNote + `
//
// Returns the options to draw, the combinations to match a selection against,
// and the price range a listing card needs. A ` + lower + ` with no variants
// gets empty lists and a range of its own price, which is what lets a storefront
// render one component whether or not variants were ever set up.
func (h *` + pascal + `VariantHandler) ListPublic(c *gin.Context) {
	var ` + snake + ` models.` + pascal + `
	if err := ` + lookup + `.First(&` + snake + `).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": "` + pascal + ` not found",
		}})
		return
	}

	options, err := h.Variants.OptionsFor(` + snake + `.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to load the options",
		}})
		return
	}
	byID := make(map[string]models.Option, len(options))
	for _, option := range options {
		byID[option.ID] = option
	}

	variants, err := h.Variants.VariantsFor(` + snake + `.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to load the variants",
		}})
		return
	}

	// Inactive combinations are dropped rather than marked. A variant somebody
	// switched off is not something the shop sells, and publishing it greyed out
	// invites a client to render a choice that can never be completed.
	published := make([]models.` + pascal + `Variant, 0, len(variants))
	out := make([]public` + pascal + `Variant, 0, len(variants))
	for _, variant := range variants {
		if !variant.Active {
			continue
		}
		published = append(published, variant)

		valueIDs := make([]string, 0, len(variant.OptionValues))
		for _, value := range variant.OptionValues {
			valueIDs = append(valueIDs, value.ID)
		}
		out = append(out, public` + pascal + `Variant{
			ID:             variant.ID,
			SKU:            variant.SKU,
			Price:          h.Variants.ResolvePrice(` + snake + `.Price, variant, byID),
			InStock:        variant.InStock(),
			Images:         variant.Images,
			OptionValueIDs: valueIDs,
		})
	}

	low, high := h.Variants.PriceRange(` + snake + `.Price, published, byID)

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
`
}
