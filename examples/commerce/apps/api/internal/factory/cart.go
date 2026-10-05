package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// NewCart builds a valid Cart without saving it.
//
// Every call produces different values for the unique fields, so two rows made
// in the same test do not collide.
func NewCart(overrides ...func(*models.Cart)) *models.Cart {
	n := next()

	record := &models.Cart{
		Token:    fmt.Sprintf("Token %d", n),
		Currency: fmt.Sprintf("Currency %d", n),
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateCart builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateCart(tb testing.TB, db *gorm.DB, overrides ...func(*models.Cart)) *models.Cart {
	tb.Helper()
	record := NewCart(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a cart fixture: %v", err)
	}
	return record
}

// CreateCarts writes n of them, for a test about paging or sorting.
func CreateCarts(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Cart)) []*models.Cart {
	tb.Helper()
	out := make([]*models.Cart, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateCart(tb, db, overrides...))
	}
	return out
}
