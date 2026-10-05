package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/concurrency"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
	"commerce/apps/api/internal/respond"
)

// CartItemService owns every database read and write for cart_items.
//
// The handler reads the request, calls one of these and writes the answer. A
// job, a command or a test calls the same methods with no request at all, so a
// rule enforced here is enforced everywhere, not only on the HTTP route.
//
// Each method takes the context it runs in. It carries the organization the
// multitenant plugin resolved, the actor ownership is scoped by (see
// authz.WithActor, and authz.AsSystem for a job), and the cancellation that
// fires when a client goes away.
type CartItemService struct {
	DB *gorm.DB
}

// cartItemListConfig is what a client may search, sort and filter cart_items
// by. Whitelisted, because each name ends up in SQL.
var cartItemListConfig = paginate.Config{
	Searchable: []string{"variant_id", "title", "variant_label", "image_url"},
	Sortable:   map[string]bool{"id": true, "created_at": true, "variant_id": true, "quantity": true, "unit_price_amount": true, "title": true, "variant_label": true, "image_url": true},
	Filterable: map[string]bool{"id": true, "cart_id": true, "product_id": true, "variant_id": true, "quantity": true, "unit_price_amount": true, "unit_price_currency": true, "title": true, "variant_label": true, "image_url": true},
}

// writableCartItem is every column Patch and Bulk may write. id, the
// timestamps and the version are the framework's, and are dropped.
var writableCartItem = map[string]bool{
	"cart_id":       true,
	"product_id":    true,
	"variant_id":    true,
	"quantity":      true,
	"unit_price":    true,
	"title":         true,
	"variant_label": true,
	"image_url":     true,
}

// db binds the database to ctx, so whatever a middleware put there reaches
// GORM's callbacks: the multitenant plugin scopes by it.
func (s *CartItemService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// write is a session for a single-statement write: no wrapping transaction,
// and RETURNING where the dialect has it. Safe only because every write it is
// used for is exactly one statement, which the generator knew.
func (s *CartItemService) write(db *gorm.DB) *gorm.DB {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if s.returning(db) {
		tx = tx.Clauses(clause.Returning{})
	}
	return tx
}

// returning reports whether the dialect hands back the written row. MySQL
// does not, and does not say so: the clause is dropped and the defaults come
// back empty, so there the row is read again.
func (s *CartItemService) returning(db *gorm.DB) bool {
	switch db.Dialector.Name() {
	case "postgres", "sqlite":
		return true
	}
	return false
}

// relations reads the rows item points at, after a write whose RETURNING
// brought the row itself back. Reading the row again with its relations
// preloaded cost one more query on every create, update and patch.
func (s *CartItemService) relations(db *gorm.DB, item *models.CartItem) error {
	item.Cart = nil
	if item.CartID != "" {
		var related []models.Cart
		if err := db.Where("id = ?", item.CartID).Limit(1).Find(&related).Error; err != nil {
			return err
		}
		if len(related) == 1 {
			item.Cart = &related[0]
		}
	}
	item.Product = nil
	if item.ProductID != "" {
		var related []models.Product
		if err := db.Where("id = ?", item.ProductID).Limit(1).Find(&related).Error; err != nil {
			return err
		}
		if len(related) == 1 {
			item.Product = &related[0]
		}
	}
	return nil
}

// List returns one page of cart_items.
//
//	archived "true" or "1"   only archived rows
//	archived "all"           both
//	anything else            only live rows
func (s *CartItemService) List(ctx context.Context, p paginate.Params, archived string) (paginate.Result[models.CartItem], error) {
	query := s.db(ctx).Model(&models.CartItem{}).Preload("Cart").Preload("Product")

	// Archived rows are excluded by default. Anything else means an operator
	// archives twelve rows, sees the count go down, and finds them again the
	// next time somebody sorts by a different column.
	switch archived {
	case "true", "1":
		query = query.Where("archived_at IS NOT NULL")
	case "all":
		// no filter
	default:
		query = query.Where("archived_at IS NULL")
	}

	return paginate.List[models.CartItem](query, p, cartItemListConfig)
}

// Export hands every matching cartitem to each, a batch at a time, so a large
// table is never in memory at once. search matches the columns List searches.
//
// No ORDER BY of its own: FindInBatches pages by primary key, which is
// creation order for the time-ordered ids Grit issues, and a sort in front of
// that key repeated rows from the second batch on.
func (s *CartItemService) Export(ctx context.Context, search string, each func(rows []models.CartItem) error) error {
	query := s.db(ctx).Model(&models.CartItem{}).Preload("Cart").Preload("Product")
	if search != "" {
		clause := ""
		args := []any{}
		wild := "%" + search + "%"
		for i, col := range cartItemListConfig.Searchable {
			if i > 0 {
				clause += " OR "
			}
			clause += "LOWER(" + col + ") LIKE LOWER(?)"
			args = append(args, wild)
		}
		if clause != "" {
			query = query.Where(clause, args...)
		}
	}

	// FindInBatches fills rows and pages by primary key. The tx it hands the
	// callback is a fresh session with no query on it, so each batch is read
	// from rows: re-reading it through tx.Scan found nothing, and every export
	// was an empty file with a 200.
	var rows []models.CartItem
	return query.FindInBatches(&rows, 1000, func(tx *gorm.DB, batch int) error {
		return each(rows)
	}).Error
}

// GetByID returns one cartitem with the relations shown beside it.
func (s *CartItemService) GetByID(ctx context.Context, id string) (*models.CartItem, error) {
	var item models.CartItem
	if err := s.db(ctx).Preload("Cart").Preload("Product").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if d := authz.Inspect(ctx, "cart_items.read", &item); !d.Allowed() {
		return nil, authz.Denied("cart_items.read", d)
	}
	return &item, nil
}

// load reads the row a write is about to change, without its relations.
//
// ability is what the caller is about to do, "update" or "delete", so a
// policy can refuse one without refusing the other.
func (s *CartItemService) load(ctx context.Context, id, ability string) (*models.CartItem, error) {
	var item models.CartItem
	if err := s.db(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// A policy rule narrows this further than the route's permission did.
	// None defined means no narrowing, which is how the project behaved
	// before it had policies. See internal/policies.
	if d := authz.Inspect(ctx, "cart_items."+ability, &item); !d.Allowed() {
		return nil, authz.Denied("cart_items."+ability, d)
	}
	return &item, nil
}

// conflict is the answer to a write whose precondition failed: the version
// the row is at now.
func (s *CartItemService) conflict(ctx context.Context, id string) error {
	var current models.CartItem
	if err := s.db(ctx).Select("version").First(&current, "id = ?", id).Error; err != nil {
		return err
	}
	return &concurrency.ErrConflict{Current: current.Version}
}

// Create saves a new cartitem and fills item in as it was stored.
func (s *CartItemService) Create(ctx context.Context, item *models.CartItem) error {
	db := s.db(ctx)
	if err := s.write(db).Create(item).Error; err != nil {
		return err
	}
	// RETURNING filled the row in, so only the rows it points at are read.
	if s.returning(db) {
		if err := s.relations(db, item); err != nil {
			return err
		}
	} else if err := db.Preload("Cart").Preload("Product").First(item, "id = ?", item.ID).Error; err != nil {
		return err
	}
	return nil
}

// Update writes updates to one cartitem. With a precondition it lands only
// if the row is still at that version, and otherwise returns an
// *concurrency.ErrConflict naming the version it is at.
func (s *CartItemService) Update(ctx context.Context, id string, updates map[string]interface{}, pre *concurrency.Precondition) (*models.CartItem, error) {
	item, err := s.load(ctx, id, "update")
	if err != nil {
		return nil, err
	}
	db := s.db(ctx)
	written := s.write(db).Model(item).Scopes(pre.Scope).Updates(updates)
	if err := written.Error; err != nil {
		return nil, err
	}
	// pre named a version this record has moved past: someone else saved
	// first. A conflict rather than overwriting their change.
	if pre.Missed(written) {
		return nil, s.conflict(ctx, item.ID)
	}
	// RETURNING filled the row in, so only the rows it points at are read.
	if s.returning(db) {
		if err := s.relations(db, item); err != nil {
			return nil, err
		}
	} else if err := db.Preload("Cart").Preload("Product").First(item, "id = ?", item.ID).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Patch writes only the columns body names, leaving every other one as it
// is. Keys that are not writable columns are dropped. It returns the row and
// the columns it wrote.
func (s *CartItemService) Patch(ctx context.Context, id string, body map[string]interface{}, pre *concurrency.Precondition) (*models.CartItem, map[string]interface{}, error) {
	item, err := s.load(ctx, id, "update")
	if err != nil {
		return nil, nil, err
	}
	db := s.db(ctx)

	updates := map[string]interface{}{}
	for k, v := range body {
		if writableCartItem[k] {
			updates[k] = v
		}
	}
	touchedRelations := false
	if len(updates) == 0 && !touchedRelations {
		return nil, nil, respond.Rule("No writable fields in request body")
	}

	if len(updates) > 0 {
		written := s.write(db).Model(item).Scopes(pre.Scope).Updates(updates)
		if err := written.Error; err != nil {
			return nil, nil, err
		}
		if pre.Missed(written) {
			return nil, nil, s.conflict(ctx, item.ID)
		}
	}
	// RETURNING filled the row in, so only the rows it points at are read.
	if s.returning(db) {
		if err := s.relations(db, item); err != nil {
			return nil, nil, err
		}
	} else if err := db.Preload("Cart").Preload("Product").First(item, "id = ?", item.ID).Error; err != nil {
		return nil, nil, err
	}
	return item, updates, nil
}

// Delete soft-deletes one cartitem and returns it as it was.
func (s *CartItemService) Delete(ctx context.Context, id string) (*models.CartItem, error) {
	item, err := s.load(ctx, id, "delete")
	if err != nil {
		return nil, err
	}
	if err := s.db(ctx).Delete(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// CartItemBulkResult is what a bulk action did.
type CartItemBulkResult struct {
	// IDs are the rows acted on: those requested that exist, that the caller
	// may touch, and that the action applies to.
	IDs []string
	// Updates are the columns a patch wrote, after the whitelist.
	Updates map[string]interface{}
}

// Bulk applies one action to many cart_items in a single transaction: all of
// it lands or none of it does. action is delete, archive, restore or patch.
func (s *CartItemService) Bulk(ctx context.Context, action string, ids []string, patch map[string]interface{}) (CartItemBulkResult, error) {
	db := s.db(ctx)
	var result CartItemBulkResult

	// Unarchived rows for archive, archived for restore: without it a mixed
	// selection reports "12 archived" having changed three.
	scope := db.Model(&models.CartItem{}).Where("id IN ?", ids)
	if action == "restore" {
		scope = scope.Where("archived_at IS NOT NULL")
	} else if action == "archive" {
		scope = scope.Where("archived_at IS NULL")
	}
	var items []models.CartItem
	if err := scope.Find(&items).Error; err != nil {
		return result, fmt.Errorf("loading cart_items: %w", err)
	}
	for _, item := range items {
		result.IDs = append(result.IDs, item.ID)
	}
	if len(result.IDs) == 0 {
		return result, nil
	}

	if action == "patch" {
		// The same whitelist as Patch. Framework-owned columns are dropped
		// rather than refused, so a client sending the whole row is not wrong.
		result.Updates = map[string]interface{}{}
		for k, v := range patch {
			if writableCartItem[k] {
				result.Updates[k] = v
			}
		}
		if len(result.Updates) == 0 {
			return CartItemBulkResult{}, respond.Rule("No writable fields in patch")
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		switch action {
		case "delete":
			return tx.Where("id IN ?", result.IDs).Delete(&models.CartItem{}).Error
		case "archive":
			return tx.Model(&models.CartItem{}).Where("id IN ?", result.IDs).
				Update("archived_at", time.Now()).Error
		case "restore":
			return tx.Model(&models.CartItem{}).Where("id IN ?", result.IDs).
				Update("archived_at", nil).Error
		case "patch":
			return tx.Model(&models.CartItem{}).Where("id IN ?", result.IDs).
				Updates(result.Updates).Error
		}
		return respond.Rule("unknown bulk action %q", action)
	})
	if err != nil {
		return CartItemBulkResult{}, err
	}
	return result, nil
}
