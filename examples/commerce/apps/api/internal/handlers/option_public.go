package handlers

import (
	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
)

// The published shape of the option library.
//
// An allowlist rather than the model, the same default the rest of the public
// surface takes: a column added to options next month is private until somebody
// adds it here.

// publicOptionValue is one choice, as a storefront may see it.
type publicOptionValue struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Slug  string `json:"slug"`

	// Swatch is the CSS colour a colour picker paints. Empty for every other
	// kind of option.
	Swatch string `json:"swatch,omitempty"`

	// Image is the swatch picture, not the photograph of the product in this
	// colour. That one belongs to the variant, because it is a picture of a
	// combination rather than of a value.
	Image *files.FileRef `json:"image,omitempty"`

	// PriceDelta lets a picker label a choice "+ 20" before anything is
	// selected. Zero unless the option declares AffectsPrice, so the label can
	// never disagree with the price the server resolves.
	PriceDelta float64 `json:"price_delta"`
}

// publicOption is one axis of choice, with its values in display order.
//
// Kind travels because it is how a storefront knows what to draw: "swatch" for
// colours, "size" for a row of boxes, "select" for a dropdown. Deciding that on
// the client from the option's name is how one shop ends up rendering Colour as
// a dropdown because somebody called it "Color".
type publicOption struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Slug         string              `json:"slug"`
	Kind         string              `json:"kind"`
	AffectsPrice bool                `json:"affects_price"`
	Values       []publicOptionValue `json:"values"`
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
