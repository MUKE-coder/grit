package scaffold

func stockTestGo() string {
	return `package stock

import (
	"errors"
	"sync"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type product struct {
	ID    string ` + "`" + `gorm:"primarykey;size:36"` + "`" + `
	Name  string
	Stock int
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	// One connection: each connection to ":memory:" is its own database, and
	// the concurrency test below would otherwise be several goroutines each
	// happily selling from a private copy of the inventory.
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&product{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seed(t *testing.T, db *gorm.DB, id string, qty int) {
	t.Helper()
	if err := db.Create(&product{ID: id, Name: "widget", Stock: qty}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func stockOf(t *testing.T, db *gorm.DB, id string) int {
	t.Helper()
	n, err := Available(db, &product{}, id, "stock")
	if err != nil {
		t.Fatalf("Available: %v", err)
	}
	return n
}

func TestTakeReducesTheCount(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 10)

	if err := Take(db, &product{}, "p1", "stock", 3); err != nil {
		t.Fatalf("Take: %v", err)
	}
	if got := stockOf(t, db, "p1"); got != 7 {
		t.Errorf("stock = %d, want 7", got)
	}
}

func TestTakeRefusesMoreThanThereIs(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 2)

	err := Take(db, &product{}, "p1", "stock", 3)
	if !errors.Is(err, ErrInsufficient) {
		t.Fatalf("error = %v, want ErrInsufficient", err)
	}
	// And it took nothing. A partial take is worse than a refusal: the units
	// are gone and no order accounts for them.
	if got := stockOf(t, db, "p1"); got != 2 {
		t.Errorf("stock = %d after a refused take, want 2", got)
	}
}

// Exactly the boundary, which is where an off-by-one lives.
func TestTakeAllOfItIsAllowed(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 5)

	if err := Take(db, &product{}, "p1", "stock", 5); err != nil {
		t.Fatalf("Take: %v", err)
	}
	if got := stockOf(t, db, "p1"); got != 0 {
		t.Errorf("stock = %d, want 0", got)
	}
	if err := Take(db, &product{}, "p1", "stock", 1); !errors.Is(err, ErrInsufficient) {
		t.Errorf("taking from an empty row returned %v, want ErrInsufficient", err)
	}
}

// A missing row and an empty row are different problems and need different
// answers. Both look like "zero rows updated" from the database.
func TestUnknownIDIsNotTheSameAsSoldOut(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 0)

	if err := Take(db, &product{}, "nope", "stock", 1); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown id returned %v, want ErrNotFound", err)
	}
	if err := Take(db, &product{}, "p1", "stock", 1); !errors.Is(err, ErrInsufficient) {
		t.Errorf("sold out returned %v, want ErrInsufficient", err)
	}
}

// The reason this package exists.
//
// Twenty goroutines race for ten units. Read-modify-write lets more than ten
// succeed and drives the count negative; a conditional UPDATE cannot.
func TestConcurrentTakesNeverOversell(t *testing.T) {
	db := newDB(t)
	const units, racers = 10, 20
	seed(t, db, "p1", units)

	var wg sync.WaitGroup
	var mu sync.Mutex
	sold, refused := 0, 0

	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := Take(db, &product{}, "p1", "stock", 1)
			mu.Lock()
			defer mu.Unlock()
			switch {
			case err == nil:
				sold++
			case errors.Is(err, ErrInsufficient):
				refused++
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	if sold != units {
		t.Errorf("sold %d of %d units", sold, units)
	}
	if refused != racers-units {
		t.Errorf("refused %d, want %d", refused, racers-units)
	}
	if got := stockOf(t, db, "p1"); got != 0 {
		t.Errorf("stock = %d after the race, want 0; negative means it oversold", got)
	}
}

func TestReleasePutsItBack(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 5)

	if err := Take(db, &product{}, "p1", "stock", 5); err != nil {
		t.Fatalf("Take: %v", err)
	}
	if err := Release(db, &product{}, "p1", "stock", 2); err != nil {
		t.Fatalf("Release: %v", err)
	}
	if got := stockOf(t, db, "p1"); got != 2 {
		t.Errorf("stock = %d, want 2", got)
	}
}

// A rollback has to return the units, which is why Take takes the caller's
// transaction rather than opening its own.
func TestTakeInsideATransactionRollsBack(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 10)
	boom := errors.New("the order failed after the stock was taken")

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := Take(tx, &product{}, "p1", "stock", 4); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("transaction error = %v", err)
	}
	if got := stockOf(t, db, "p1"); got != 10 {
		t.Errorf("stock = %d after a rolled-back order, want 10", got)
	}
}

// All or nothing across several products.
func TestTakeManyIsAllOrNothing(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 5)
	seed(t, db, "p2", 1)

	err := db.Transaction(func(tx *gorm.DB) error {
		return TakeMany(tx, &product{}, "stock", []Line{
			{ID: "p1", Qty: 2},
			{ID: "p2", Qty: 3}, // more than p2 has
		})
	})
	if !errors.Is(err, ErrInsufficient) {
		t.Fatalf("error = %v, want ErrInsufficient", err)
	}
	if got := stockOf(t, db, "p1"); got != 5 {
		t.Errorf("p1 = %d; the first line was taken and not returned when the second failed", got)
	}
	if got := stockOf(t, db, "p2"); got != 1 {
		t.Errorf("p2 = %d, want 1", got)
	}
}

func TestTakeManySucceedsTogether(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 5)
	seed(t, db, "p2", 5)

	err := db.Transaction(func(tx *gorm.DB) error {
		return TakeMany(tx, &product{}, "stock", []Line{
			{ID: "p2", Qty: 1},
			{ID: "p1", Qty: 2},
		})
	})
	if err != nil {
		t.Fatalf("TakeMany: %v", err)
	}
	if got := stockOf(t, db, "p1"); got != 3 {
		t.Errorf("p1 = %d, want 3", got)
	}
	if got := stockOf(t, db, "p2"); got != 4 {
		t.Errorf("p2 = %d, want 4", got)
	}
}

// The column name is the one part of the statement that cannot be a bound
// parameter, so it is checked rather than trusted.
func TestColumnNameIsValidated(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 5)

	for _, bad := range []string{
		"stock; DROP TABLE products",
		"stock = 0 WHERE 1=1 --",
		"",
		"Stock",
	} {
		if err := Take(db, &product{}, "p1", bad, 1); err == nil {
			t.Errorf("Take accepted %q as a column name", bad)
		}
	}
	if got := stockOf(t, db, "p1"); got != 5 {
		t.Errorf("stock = %d; a rejected call still changed the row", got)
	}
}

func TestZeroAndNegativeQuantities(t *testing.T) {
	db := newDB(t)
	seed(t, db, "p1", 5)

	// Zero is an empty basket line, not an error.
	if err := Take(db, &product{}, "p1", "stock", 0); err != nil {
		t.Errorf("Take(0) = %v, want nil", err)
	}
	// Negative would be a take that adds stock, which is how you would smuggle
	// inventory into existence through a checkout.
	if err := Take(db, &product{}, "p1", "stock", -3); err == nil {
		t.Error("Take accepted a negative quantity")
	}
	if got := stockOf(t, db, "p1"); got != 5 {
		t.Errorf("stock = %d, want 5", got)
	}
}
`
}
