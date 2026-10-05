package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
)

// NewProduct builds a valid Product without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
//
// CollectionID is left empty: the factory cannot know which parent you mean, and
// inserting one would make every test that touches this resource write rows in
// another table too. Pass it in:
//
//	product := factory.CreateProduct(t, db, func(m *models.Product) {
//		m.CollectionID = parent.ID
//	})
func NewProduct(overrides ...func(*models.Product)) *models.Product {
	n := next()

	record := &models.Product{
		Title:       fmt.Sprintf("Product %d", n),
		Description: fmt.Sprintf("A description for fixture %d", n),
		Price:       money.New(1999+int64(n), "USD"),
		Available:   true,
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateProduct builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateProduct(tb testing.TB, db *gorm.DB, overrides ...func(*models.Product)) *models.Product {
	tb.Helper()
	record := NewProduct(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a product fixture: %v", err)
	}
	return record
}

// CreateProducts writes n of them, for a test about paging or sorting.
func CreateProducts(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Product)) []*models.Product {
	tb.Helper()
	out := make([]*models.Product, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateProduct(tb, db, overrides...))
	}
	return out
}
