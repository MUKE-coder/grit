package scaffold

import "strings"

// apiVariantServiceGo emits internal/services/<resource>_variants.go.
//
// Everything a variant system needs that is not CRUD: resolving a price,
// building the option payload a storefront renders, and generating the
// combinations so nobody types twelve rows by hand.
func apiVariantServiceGo(module, pascal, snake, plural string) string {
	lower := strings.ToLower(pascal)
	return `package services

import (
	"fmt"
	"sort"

	"gorm.io/gorm"

	"` + module + `/internal/models"
)

// ` + pascal + `VariantService owns the questions a variant system is for.
//
// The important one is what a combination costs, and the answer is computed
// rather than stored. A stored price is a copy that goes stale the moment the
// ` + lower + `'s own price changes, and the symptom is a receipt that disagrees
// with the page the customer bought from.
type ` + pascal + `VariantService struct {
	DB *gorm.DB
}

func New` + pascal + `VariantService(db *gorm.DB) *` + pascal + `VariantService {
	return &` + pascal + `VariantService{DB: db}
}

// ResolvePrice returns what one variant costs.
//
// An override wins outright: it is the escape hatch for a combination priced by
// hand, and second-guessing it would make the field useless.
//
// Otherwise the price is the ` + lower + `'s own, plus the delta of every value
// whose OPTION declares AffectsPrice. Reading the flag from the option rather
// than the value is what stops a half-configured shop charging extra for a
// colour because somebody typed a delta on one swatch.
func (s *` + pascal + `VariantService) ResolvePrice(base float64, variant models.` + pascal + `Variant, optionsByID map[string]models.Option) float64 {
	if variant.PriceOverride != nil {
		return *variant.PriceOverride
	}
	price := base
	for _, value := range variant.OptionValues {
		option, ok := optionsByID[value.OptionID]
		if !ok || !option.AffectsPrice {
			continue
		}
		price += value.PriceDelta
	}
	return price
}

// OptionsFor returns the options a ` + lower + ` offers, each with its values,
// in the order the storefront should draw them.
//
// Two queries, not one per option. A picker with four axes is four round trips
// under the obvious implementation, on the page with the highest traffic in the
// shop.
func (s *` + pascal + `VariantService) OptionsFor(` + snake + `ID string) ([]models.Option, error) {
	var links []models.` + pascal + `Option
	err := s.DB.Where("` + snake + `_id = ?", ` + snake + `ID).
		Order("position asc").
		Find(&links).Error
	if err != nil {
		return nil, fmt.Errorf("loading ` + lower + ` options: %w", err)
	}
	if len(links) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(links))
	order := make(map[string]int, len(links))
	for i, link := range links {
		ids = append(ids, link.OptionID)
		order[link.OptionID] = i
	}

	var options []models.Option
	err = s.DB.Preload("Values", func(db *gorm.DB) *gorm.DB {
		return db.Order("position asc, label asc")
	}).Where("id IN ?", ids).Find(&options).Error
	if err != nil {
		return nil, fmt.Errorf("loading options: %w", err)
	}

	// IN does not promise an order, and the order here was somebody's
	// decision about which axis to show first.
	sort.SliceStable(options, func(i, j int) bool {
		return order[options[i].ID] < order[options[j].ID]
	})
	return options, nil
}

// VariantsFor returns every variant of a ` + lower + ` with its values attached.
func (s *` + pascal + `VariantService) VariantsFor(` + snake + `ID string) ([]models.` + pascal + `Variant, error) {
	var variants []models.` + pascal + `Variant
	err := s.DB.Preload("OptionValues").
		Where("` + snake + `_id = ?", ` + snake + `ID).
		Order("position asc").
		Find(&variants).Error
	if err != nil {
		return nil, fmt.Errorf("loading variants: %w", err)
	}
	return variants, nil
}

// FindByValues returns the variant defined by exactly this set of option value
// ids, which is the lookup a storefront makes when somebody picks Red and XXL.
//
// Exactly, in both directions: a variant with more values than were asked for
// is a different combination, not a match. Comparing only the asked-for ids
// would return the first variant that happens to include them.
func (s *` + pascal + `VariantService) FindByValues(` + snake + `ID string, valueIDs []string) (*models.` + pascal + `Variant, error) {
	variants, err := s.VariantsFor(` + snake + `ID)
	if err != nil {
		return nil, err
	}

	wanted := make(map[string]bool, len(valueIDs))
	for _, id := range valueIDs {
		wanted[id] = true
	}

	for i := range variants {
		if len(variants[i].OptionValues) != len(wanted) {
			continue
		}
		match := true
		for _, value := range variants[i].OptionValues {
			if !wanted[value.ID] {
				match = false
				break
			}
		}
		if match {
			return &variants[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// PriceRange returns the cheapest and dearest in-stock variant price, for the
// "from X" a listing shows.
//
// Out-of-stock variants are excluded, because a "from 49" that can only be had
// by buying something unavailable is a lie the customer discovers at the last
// step.
func (s *` + pascal + `VariantService) PriceRange(base float64, variants []models.` + pascal + `Variant, optionsByID map[string]models.Option) (low, high float64) {
	first := true
	for _, variant := range variants {
		if !variant.InStock() {
			continue
		}
		price := s.ResolvePrice(base, variant, optionsByID)
		if first {
			low, high, first = price, price, false
			continue
		}
		if price < low {
			low = price
		}
		if price > high {
			high = price
		}
	}
	if first {
		return base, base
	}
	return low, high
}

// Generate creates the missing combinations of the ` + lower + `'s options.
//
// Existing variants are left exactly as they are: their sku, stock, price and
// photographs are somebody's work, and the point of this is to add the rows
// nobody has filled in yet, not to reset the ones they have.
//
// It refuses past a cap, because the cartesian product of four options with
// five values each is 625 rows, and a UI that silently writes those is a UI
// that has destroyed the page it was meant to help with.
func (s *` + pascal + `VariantService) Generate(` + snake + `ID string, limit int) (created int, err error) {
	if limit <= 0 {
		limit = 200
	}

	options, err := s.OptionsFor(` + snake + `ID)
	if err != nil {
		return 0, err
	}
	if len(options) == 0 {
		return 0, fmt.Errorf("` + lower + ` %s offers no options, so it has no combinations", ` + snake + `ID)
	}

	// The cartesian product, built iteratively rather than recursively: the
	// number of axes is data, and a recursive version needs a depth nobody
	// declared.
	combinations := [][]models.OptionValue{{}}
	for _, option := range options {
		if len(option.Values) == 0 {
			return 0, fmt.Errorf("option %q has no values, so no combination can include it", option.Name)
		}
		next := make([][]models.OptionValue, 0, len(combinations)*len(option.Values))
		for _, existing := range combinations {
			for _, value := range option.Values {
				row := make([]models.OptionValue, len(existing), len(existing)+1)
				copy(row, existing)
				next = append(next, append(row, value))
			}
		}
		combinations = next
		if len(combinations) > limit {
			return 0, fmt.Errorf(
				"that would be %d or more combinations, past the limit of %d. Narrow the options first",
				len(combinations), limit)
		}
	}

	existing, err := s.VariantsFor(` + snake + `ID)
	if err != nil {
		return 0, err
	}
	seen := make(map[string]bool, len(existing))
	for _, variant := range existing {
		seen[fingerprint(variant.OptionValues)] = true
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		for i, combination := range combinations {
			if seen[fingerprint(combination)] {
				continue
			}
			variant := models.` + pascal + `Variant{
				` + pascal + `ID: ` + snake + `ID,
				Active:   true,
				Position: len(existing) + i,
			}
			if err := tx.Create(&variant).Error; err != nil {
				return fmt.Errorf("creating variant: %w", err)
			}
			if err := tx.Model(&variant).Association("OptionValues").Replace(combination); err != nil {
				return fmt.Errorf("attaching values: %w", err)
			}
			created++
		}
		return nil
	})
	return created, err
}

// fingerprint identifies a combination by its value ids, order-independently.
func fingerprint(values []models.OptionValue) string {
	ids := make([]string, 0, len(values))
	for _, value := range values {
		ids = append(ids, value.ID)
	}
	sort.Strings(ids)
	out := ""
	for _, id := range ids {
		out += id + "|"
	}
	return out
}
`
}

// APIVariantServiceGo is exported for grit add variants.
func APIVariantServiceGo(module, pascal, snake, plural string) string {
	return apiVariantServiceGo(module, pascal, snake, plural)
}
