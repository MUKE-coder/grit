package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
)

// NewCartItem builds a valid CartItem without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// CartID and ProductID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	cartItem := factory.CreateCartItem(t, db, func(m *models.CartItem) {
//		m.CartID = parent.ID
//	})
func NewCartItem(overrides ...func(*models.CartItem)) *models.CartItem {
	n := next()

	record := &models.CartItem{
		VariantID:    fmt.Sprintf("VariantID %d", n),
		Quantity:     n,
		UnitPrice:    money.New(1999+int64(n), "USD"),
		Title:        fmt.Sprintf("CartItem %d", n),
		VariantLabel: fmt.Sprintf("VariantLabel %d", n),
		ImageURL:     fmt.Sprintf("ImageURL %d", n),
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateCartItem builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateCartItem(tb testing.TB, db *gorm.DB, overrides ...func(*models.CartItem)) *models.CartItem {
	tb.Helper()
	record := NewCartItem(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a cartitem fixture: %v", err)
	}
	return record
}

// CreateCartItems writes n of them, for a test about paging or sorting.
func CreateCartItems(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.CartItem)) []*models.CartItem {
	tb.Helper()
	out := make([]*models.CartItem, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateCartItem(tb, db, overrides...))
	}
	return out
}
