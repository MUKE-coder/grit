package services

import (
	"fmt"
	"sort"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
)

// ProductVariantService owns the questions a variant system is for.
//
// The important one is what a combination costs, and the answer is computed
// rather than stored. A stored price is a copy that goes stale the moment the
// product's own price changes, and the symptom is a receipt that disagrees
// with the page the customer bought from.
type ProductVariantService struct {
	DB *gorm.DB
}

func NewProductVariantService(db *gorm.DB) *ProductVariantService {
	return &ProductVariantService{DB: db}
}

// ResolvePrice returns what one variant costs.
//
// An override wins outright: it is the escape hatch for a combination priced by
// hand, and second-guessing it would make the field useless.
//
// Otherwise the price is the product's own, plus the delta of every value
// whose OPTION declares AffectsPrice. Reading the flag from the option rather
// than the value is what stops a half-configured shop charging extra for a
// colour because somebody typed a delta on one swatch.
func (s *ProductVariantService) ResolvePrice(base money.Money, variant models.ProductVariant, optionsByID map[string]models.Option) money.Money {
	if variant.PriceOverride != nil {
		return money.FromMajor(*variant.PriceOverride, base.Currency)
	}
	delta := 0.0
	for _, value := range variant.OptionValues {
		option, ok := optionsByID[value.OptionID]
		if !ok || !option.AffectsPrice {
			continue
		}
		delta += value.PriceDelta
	}
	if delta == 0 {
		return base
	}
	return money.FromMajor(base.Major()+delta, base.Currency)
}

// OptionsFor returns the options a product offers, each with its values,
// in the order the storefront should draw them.
//
// Two queries, not one per option. A picker with four axes is four round trips
// under the obvious implementation, on the page with the highest traffic in the
// shop.
func (s *ProductVariantService) OptionsFor(productID string) ([]models.Option, error) {
	var links []models.ProductOption
	err := s.DB.Where("product_id = ?", productID).
		Order("position asc").
		Find(&links).Error
	if err != nil {
		return nil, fmt.Errorf("loading product options: %w", err)
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

	// And now only the values this product offers. Every caller reads its
	// options through here, so narrowing once covers the matrix generator, the
	// public payload and the admin picker alike.
	if err := s.narrowValues(productID, options); err != nil {
		return nil, err
	}
	return options, nil
}

// narrowValues drops the values this product does not offer.
//
// An option with no rows for this product is offered whole, which is what
// makes the feature additive: a catalogue that has never narrowed anything
// behaves exactly as it did, and nothing had to be backfilled when the table
// arrived.
func (s *ProductVariantService) narrowValues(productID string, options []models.Option) error {
	var chosen []models.ProductOptionValue
	if err := s.DB.Where("product_id = ?", productID).Find(&chosen).Error; err != nil {
		return fmt.Errorf("loading the values this product offers: %w", err)
	}
	if len(chosen) == 0 {
		return nil
	}

	picked := make(map[string]bool, len(chosen))
	for _, row := range chosen {
		picked[row.OptionValueID] = true
	}

	for i := range options {
		// Which of this option's values were picked, if any. An option nobody
		// narrowed keeps all of them.
		kept := make([]models.OptionValue, 0, len(options[i].Values))
		for _, value := range options[i].Values {
			if picked[value.ID] {
				kept = append(kept, value)
			}
		}
		if len(kept) > 0 {
			options[i].Values = kept
		}
	}
	return nil
}

// SetOfferedValues records which values of its options a product offers.
//
// Passing none for an option means all of them, so clearing the list is how a
// product goes back to the whole axis.
func (s *ProductVariantService) SetOfferedValues(productID string, valueIDs []string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		keep, err := OfferedProductValues(tx, productID, nil, valueIDs)
		if err != nil {
			return err
		}
		return WriteProductOfferedValues(tx, productID, keep)
	})
}

// OfferedProductValues works out what to store for a set of ticked values.
//
// Two rules, and both matter. Ids belonging to an option the product does
// not offer are dropped, because the picker sends what it has and an axis
// removed in the same edit is not an error. And an option whose every value is
// ticked stores nothing at all: "all of them" has to keep meaning all of them,
// or a colour added to the shop next month would be quietly missing from every
// product that had said yes to the whole axis.
//
// optionIDs is the set the product will offer after this edit. Nil reads
// it from what it offers now, which is what a caller outside the options
// endpoint wants.
func OfferedProductValues(tx *gorm.DB, productID string, optionIDs []string, valueIDs []string) ([]string, error) {
	if len(valueIDs) == 0 {
		return nil, nil
	}
	if optionIDs == nil {
		var links []models.ProductOption
		if err := tx.Where("product_id = ?", productID).Find(&links).Error; err != nil {
			return nil, fmt.Errorf("loading product options: %w", err)
		}
		for _, link := range links {
			optionIDs = append(optionIDs, link.OptionID)
		}
	}
	if len(optionIDs) == 0 {
		return nil, nil
	}

	var values []models.OptionValue
	if err := tx.Where("id IN ? AND option_id IN ?", valueIDs, optionIDs).Find(&values).Error; err != nil {
		return nil, fmt.Errorf("loading the chosen values: %w", err)
	}

	byOption := map[string][]string{}
	for _, value := range values {
		byOption[value.OptionID] = append(byOption[value.OptionID], value.ID)
	}

	keep := make([]string, 0, len(values))
	for optionID, chosen := range byOption {
		var total int64
		if err := tx.Model(&models.OptionValue{}).Where("option_id = ?", optionID).Count(&total).Error; err != nil {
			return nil, fmt.Errorf("counting the values of an option: %w", err)
		}
		if int64(len(chosen)) >= total {
			continue
		}
		keep = append(keep, chosen...)
	}
	sort.Strings(keep)
	return keep, nil
}

// WriteProductOfferedValues replaces the stored set with these ids.
func WriteProductOfferedValues(tx *gorm.DB, productID string, valueIDs []string) error {
	if err := tx.Unscoped().Where("product_id = ?", productID).
		Delete(&models.ProductOptionValue{}).Error; err != nil {
		return fmt.Errorf("clearing the offered values: %w", err)
	}
	for _, valueID := range valueIDs {
		row := models.ProductOptionValue{ProductID: productID, OptionValueID: valueID}
		if err := tx.Create(&row).Error; err != nil {
			return fmt.Errorf("recording an offered value: %w", err)
		}
	}
	return nil
}

// OfferedValueIDs is what the admin picker ticks: the values this product
// has been narrowed to, which is empty for every option it offers whole.
func (s *ProductVariantService) OfferedValueIDs(productID string) ([]string, error) {
	var rows []models.ProductOptionValue
	if err := s.DB.Where("product_id = ?", productID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("loading the values this product offers: %w", err)
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.OptionValueID)
	}
	sort.Strings(out)
	return out, nil
}

// VariantsFor returns every variant of a product with its values attached.
func (s *ProductVariantService) VariantsFor(productID string) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := s.DB.Preload("OptionValues").
		Where("product_id = ?", productID).
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
func (s *ProductVariantService) FindByValues(productID string, valueIDs []string) (*models.ProductVariant, error) {
	variants, err := s.VariantsFor(productID)
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
func (s *ProductVariantService) PriceRange(base money.Money, variants []models.ProductVariant, optionsByID map[string]models.Option) (low, high money.Money) {
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
		if price.Amount < low.Amount {
			low = price
		}
		if price.Amount > high.Amount {
			high = price
		}
	}
	if first {
		return base, base
	}
	return low, high
}

// Generate creates the missing combinations of the product's options.
//
// Existing variants are left exactly as they are: their sku, stock, price and
// photographs are somebody's work, and the point of this is to add the rows
// nobody has filled in yet, not to reset the ones they have.
//
// It refuses past a cap, because the cartesian product of four options with
// five values each is 625 rows, and a UI that silently writes those is a UI
// that has destroyed the page it was meant to help with.
func (s *ProductVariantService) Generate(productID string, limit int) (created int, err error) {
	if limit <= 0 {
		limit = 200
	}

	options, err := s.OptionsFor(productID)
	if err != nil {
		return 0, err
	}
	if len(options) == 0 {
		return 0, fmt.Errorf("product %s offers no options, so it has no combinations", productID)
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

	existing, err := s.VariantsFor(productID)
	if err != nil {
		return 0, err
	}
	seen := make(map[string]bool, len(existing))
	for _, variant := range existing {
		seen[productFingerprint(variant.OptionValues)] = true
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		for i, combination := range combinations {
			if seen[productFingerprint(combination)] {
				continue
			}
			variant := models.ProductVariant{
				ProductID: productID,
				Active:    true,
				Position:  len(existing) + i,
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

// productFingerprint identifies a combination by its value ids, order-independently.
// Named for the resource, because a second resource with variants puts a
// second one in the same package.
func productFingerprint(values []models.OptionValue) string {
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
