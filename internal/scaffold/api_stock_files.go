package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeStockFiles writes internal/stock: the atomic decrement.
func writeStockFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "stock", "stock.go"):      stockGo(),
		filepath.Join(apiRoot, "internal", "stock", "stock_test.go"): stockTestGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func stockGo() string {
	return `// Package stock takes from a counted quantity without overselling.
//
// The obvious version is wrong:
//
//	if product.Stock >= qty {
//	    product.Stock -= qty
//	    db.Save(&product)
//	}
//
// Two orders for the last item both read Stock = 1, both find it sufficient,
// and both write 0. One unit, two sales, and the row does not even look
// strange afterwards. The window is the few milliseconds between the read and
// the write, which is exactly when a checkout is busiest.
//
// The fix is not a mutex, and not a longer transaction. It is to stop reading
// the value at all:
//
//	UPDATE products SET stock = stock - 1 WHERE id = ? AND stock >= 1
//
// The database evaluates the condition and applies the change as one
// statement, so the second order updates zero rows and is told so. There is no
// window because there is no gap between deciding and doing.
package stock

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"gorm.io/gorm"
)

// ErrInsufficient means the row did not have enough to give.
//
// Returned rather than silently taking what is there: a checkout that quietly
// ships three of the five somebody ordered is a worse outcome than one that
// says no.
var ErrInsufficient = errors.New("stock: not enough available")

// ErrNotFound means no row matched the id at all.
//
// Distinct from ErrInsufficient on purpose. Both come back as "zero rows
// updated" from the database, and collapsing them tells a caller their
// catalogue is out of stock when what actually happened is a bad product id.
var ErrNotFound = errors.New("stock: no such row")

// safeColumn guards the one part of the statement that is not a bound
// parameter.
//
// The column name is interpolated into SQL because it cannot be a parameter,
// so it is checked against a pattern rather than trusted. Callers pass a
// literal today; the day one passes something from a request, this is the
// difference between a rejected call and a rewritten WHERE clause.
var safeColumn = regexp.MustCompile(` + "`" + `^[a-z_][a-z0-9_]*$` + "`" + `)

// Take removes qty from a column, or reports that it could not.
//
//	err := db.Transaction(func(tx *gorm.DB) error {
//	    if err := stock.Take(tx, &models.Product{}, id, "stock", 3); err != nil {
//	        return err
//	    }
//	    return tx.Create(&orderLine).Error
//	})
//
// Pass the transaction, not the root handle, whenever the decrement belongs
// with other writes. Taking stock in its own transaction and then failing to
// write the order line loses the units with nothing to show for them.
func Take(db *gorm.DB, model any, id, column string, qty int) error {
	return adjust(db, model, id, column, qty, true)
}

// Release adds qty back, for a cancelled order or an expired reservation.
//
// No ceiling. A maximum would need somewhere to read it from, and inventing
// one here would be a rule nobody asked for enforced in the wrong place.
func Release(db *gorm.DB, model any, id, column string, qty int) error {
	return adjust(db, model, id, column, qty, false)
}

func adjust(db *gorm.DB, model any, id, column string, qty int, take bool) error {
	if qty < 0 {
		return fmt.Errorf("stock: qty must not be negative, got %d", qty)
	}
	if qty == 0 {
		// Not an error, and not a statement either. Zero is what an empty
		// basket line looks like, and a no-op is the honest response.
		return nil
	}
	if !safeColumn.MatchString(column) {
		return fmt.Errorf("stock: %q is not a column name", column)
	}

	// Taking carries a floor in the WHERE clause; releasing does not, because
	// adding back can never go negative.
	op := "+"
	where, args := "id = ?", []any{id}
	if take {
		op = "-"
		where, args = fmt.Sprintf("id = ? AND %s >= ?", column), []any{id, qty}
	}

	// UpdateColumn, not Update: this must be one statement with no hooks and
	// no updated_at, so the value the database compares is the value the
	// database holds, not one this process read a moment ago.
	res := db.Model(model).
		Where(where, args...).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s %s ?", column, op), qty))

	if res.Error != nil {
		return fmt.Errorf("stock: updating %s: %w", column, res.Error)
	}
	if res.RowsAffected == 1 {
		return nil
	}

	// Zero rows means the id was wrong or the guard failed, and the two need
	// different answers. One more read, only on the failure path, so the happy
	// path stays a single statement.
	var exists int64
	if err := db.Model(model).Where("id = ?", id).Count(&exists).Error; err != nil {
		return fmt.Errorf("stock: checking %s: %w", column, err)
	}
	if exists == 0 {
		return ErrNotFound
	}
	return ErrInsufficient
}

// Line is one item in a multi-item take.
type Line struct {
	ID  string
	Qty int
}

// TakeMany removes from several rows, all or nothing.
//
// Runs inside the caller's transaction so a later failure returns every unit.
//
// The lines are sorted by id first, and that is not tidiness. Two concurrent
// orders touching the same two products in opposite orders each hold the row
// the other is waiting for, and the database resolves it by killing one of
// them: a deadlock error on a checkout, under exactly the load that makes it
// most expensive. Locking in a consistent order means the second waits for the
// first instead.
func TakeMany(db *gorm.DB, model any, column string, lines []Line) error {
	ordered := make([]Line, len(lines))
	copy(ordered, lines)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })

	for _, l := range ordered {
		if err := Take(db, model, l.ID, column, l.Qty); err != nil {
			return fmt.Errorf("%s: %w", l.ID, err)
		}
	}
	return nil
}

// Available reads the current value.
//
// For display. Do not read it, decide, and then write: that is the race this
// package exists to remove, and Take already makes the decision atomically.
func Available(db *gorm.DB, model any, id, column string) (int, error) {
	if !safeColumn.MatchString(column) {
		return 0, fmt.Errorf("stock: %q is not a column name", column)
	}
	var n int
	err := db.Model(model).Where("id = ?", id).Pluck(column, &n).Error
	if err != nil {
		return 0, fmt.Errorf("stock: reading %s: %w", column, err)
	}
	return n, nil
}
`
}
