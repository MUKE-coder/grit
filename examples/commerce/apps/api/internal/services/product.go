package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/concurrency"
	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/storage"
)

// ProductService owns every database read and write for products.
//
// The handler reads the request, calls one of these and writes the answer. A
// job, a command or a test calls the same methods with no request at all, so a
// rule enforced here is enforced everywhere, not only on the HTTP route.
//
// Each method takes the context it runs in. It carries the organization the
// multitenant plugin resolved, the actor ownership is scoped by (see
// authz.WithActor, and authz.AsSystem for a job), and the cancellation that
// fires when a client goes away.
type ProductService struct {
	DB *gorm.DB
	// Storage removes files a write replaced and claims the ones it keeps.
	Storage *storage.Storage
}

// productListConfig is what a client may search, sort and filter products
// by. Whitelisted, because each name ends up in SQL.
var productListConfig = paginate.Config{
	Searchable: []string{"title", "handle", "description"},
	Sortable:   map[string]bool{"id": true, "created_at": true, "title": true, "handle": true, "description": true, "price_amount": true},
	Filterable: map[string]bool{"id": true, "title": true, "handle": true, "description": true, "price_amount": true, "price_currency": true, "featured_image": true, "images": true, "available": true, "collection_id": true},
}

// writableProduct is every column Patch and Bulk may write. id, the
// timestamps and the version are the framework's, and are dropped.
var writableProduct = map[string]bool{
	"title":          true,
	"description":    true,
	"price":          true,
	"featured_image": true,
	"images":         true,
	"available":      true,
	"collection_id":  true,
}

// db binds the database to ctx, so whatever a middleware put there reaches
// GORM's callbacks: the multitenant plugin scopes by it.
func (s *ProductService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// write is a session for a single-statement write: no wrapping transaction,
// and RETURNING where the dialect has it. Safe only because every write it is
// used for is exactly one statement, which the generator knew.
func (s *ProductService) write(db *gorm.DB) *gorm.DB {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if s.returning(db) {
		tx = tx.Clauses(clause.Returning{})
	}
	return tx
}

// returning reports whether the dialect hands back the written row. MySQL
// does not, and does not say so: the clause is dropped and the defaults come
// back empty, so there the row is read again.
func (s *ProductService) returning(db *gorm.DB) bool {
	switch db.Dialector.Name() {
	case "postgres", "sqlite":
		return true
	}
	return false
}

// relations reads the rows item points at, after a write whose RETURNING
// brought the row itself back. Reading the row again with its relations
// preloaded cost one more query on every create, update and patch.
func (s *ProductService) relations(db *gorm.DB, item *models.Product) error {
	item.Collection = nil
	if item.CollectionID != "" {
		var related []models.Collection
		if err := db.Where("id = ?", item.CollectionID).Limit(1).Find(&related).Error; err != nil {
			return err
		}
		if len(related) == 1 {
			item.Collection = &related[0]
		}
	}
	return nil
}

// List returns one page of products.
//
//	archived "true" or "1"   only archived rows
//	archived "all"           both
//	anything else            only live rows
func (s *ProductService) List(ctx context.Context, p paginate.Params, archived string) (paginate.Result[models.Product], error) {
	query := s.db(ctx).Model(&models.Product{}).Preload("Collection")

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

	return paginate.List[models.Product](query, p, productListConfig)
}

// Export hands every matching product to each, a batch at a time, so a large
// table is never in memory at once. search matches the columns List searches.
//
// No ORDER BY of its own: FindInBatches pages by primary key, which is
// creation order for the time-ordered ids Grit issues, and a sort in front of
// that key repeated rows from the second batch on.
func (s *ProductService) Export(ctx context.Context, search string, each func(rows []models.Product) error) error {
	query := s.db(ctx).Model(&models.Product{}).Preload("Collection")
	if search != "" {
		clause := ""
		args := []any{}
		wild := "%" + search + "%"
		for i, col := range productListConfig.Searchable {
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
	var rows []models.Product
	return query.FindInBatches(&rows, 1000, func(tx *gorm.DB, batch int) error {
		return each(rows)
	}).Error
}

// GetByID returns one product with the relations shown beside it.
func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	var item models.Product
	if err := s.db(ctx).Preload("Collection").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if d := authz.Inspect(ctx, "products.read", &item); !d.Allowed() {
		return nil, authz.Denied("products.read", d)
	}
	return &item, nil
}

// load reads the row a write is about to change, without its relations.
//
// ability is what the caller is about to do, "update" or "delete", so a
// policy can refuse one without refusing the other.
func (s *ProductService) load(ctx context.Context, id, ability string) (*models.Product, error) {
	var item models.Product
	if err := s.db(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// A policy rule narrows this further than the route's permission did.
	// None defined means no narrowing, which is how the project behaved
	// before it had policies. See internal/policies.
	if d := authz.Inspect(ctx, "products."+ability, &item); !d.Allowed() {
		return nil, authz.Denied("products."+ability, d)
	}
	return &item, nil
}

// conflict is the answer to a write whose precondition failed: the version
// the row is at now.
func (s *ProductService) conflict(ctx context.Context, id string) error {
	var current models.Product
	if err := s.db(ctx).Select("version").First(&current, "id = ?", id).Error; err != nil {
		return err
	}
	return &concurrency.ErrConflict{Current: current.Version}
}

// Create saves a new product and fills item in as it was stored.
func (s *ProductService) Create(ctx context.Context, item *models.Product) error {
	db := s.db(ctx)
	if err := s.write(db).Create(item).Error; err != nil {
		return err
	}
	// RETURNING filled the row in, so only the rows it points at are read.
	if s.returning(db) {
		if err := s.relations(db, item); err != nil {
			return err
		}
	} else if err := db.Preload("Collection").First(item, "id = ?", item.ID).Error; err != nil {
		return err
	}
	if s.Storage != nil {
		files.ClaimRefs(ctx, db, item)
	}
	return nil
}

// Update writes updates to one product. With a precondition it lands only
// if the row is still at that version, and otherwise returns an
// *concurrency.ErrConflict naming the version it is at.
func (s *ProductService) Update(ctx context.Context, id string, updates map[string]interface{}, pre *concurrency.Precondition) (*models.Product, error) {
	item, err := s.load(ctx, id, "update")
	if err != nil {
		return nil, err
	}
	db := s.db(ctx)
	oldItem := *item // the files before the write, to find the ones it replaced
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
	} else if err := db.Preload("Collection").First(item, "id = ?", item.ID).Error; err != nil {
		return nil, err
	}
	if s.Storage != nil {
		files.CleanupRemoved(ctx, s.Storage, &oldItem, item)
		files.ClaimRefs(ctx, db, item)
	}
	return item, nil
}

// Patch writes only the columns body names, leaving every other one as it
// is. Keys that are not writable columns are dropped. It returns the row and
// the columns it wrote.
func (s *ProductService) Patch(ctx context.Context, id string, body map[string]interface{}, pre *concurrency.Precondition) (*models.Product, map[string]interface{}, error) {
	item, err := s.load(ctx, id, "update")
	if err != nil {
		return nil, nil, err
	}
	db := s.db(ctx)

	updates := map[string]interface{}{}
	for k, v := range body {
		if writableProduct[k] {
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
	} else if err := db.Preload("Collection").First(item, "id = ?", item.ID).Error; err != nil {
		return nil, nil, err
	}
	return item, updates, nil
}

// Delete soft-deletes one product and returns it as it was.
func (s *ProductService) Delete(ctx context.Context, id string) (*models.Product, error) {
	item, err := s.load(ctx, id, "delete")
	if err != nil {
		return nil, err
	}
	if err := s.db(ctx).Delete(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// ProductBulkResult is what a bulk action did.
type ProductBulkResult struct {
	// IDs are the rows acted on: those requested that exist, that the caller
	// may touch, and that the action applies to.
	IDs []string
	// Updates are the columns a patch wrote, after the whitelist.
	Updates map[string]interface{}
}

// Bulk applies one action to many products in a single transaction: all of
// it lands or none of it does. action is delete, archive, restore or patch.
func (s *ProductService) Bulk(ctx context.Context, action string, ids []string, patch map[string]interface{}) (ProductBulkResult, error) {
	db := s.db(ctx)
	var result ProductBulkResult

	// Unarchived rows for archive, archived for restore: without it a mixed
	// selection reports "12 archived" having changed three.
	scope := db.Model(&models.Product{}).Where("id IN ?", ids)
	if action == "restore" {
		scope = scope.Where("archived_at IS NOT NULL")
	} else if action == "archive" {
		scope = scope.Where("archived_at IS NULL")
	}
	var items []models.Product
	if err := scope.Find(&items).Error; err != nil {
		return result, fmt.Errorf("loading products: %w", err)
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
			if writableProduct[k] {
				result.Updates[k] = v
			}
		}
		if len(result.Updates) == 0 {
			return ProductBulkResult{}, respond.Rule("No writable fields in patch")
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		switch action {
		case "delete":
			return tx.Where("id IN ?", result.IDs).Delete(&models.Product{}).Error
		case "archive":
			return tx.Model(&models.Product{}).Where("id IN ?", result.IDs).
				Update("archived_at", time.Now()).Error
		case "restore":
			return tx.Model(&models.Product{}).Where("id IN ?", result.IDs).
				Update("archived_at", nil).Error
		case "patch":
			return tx.Model(&models.Product{}).Where("id IN ?", result.IDs).
				Updates(result.Updates).Error
		}
		return respond.Rule("unknown bulk action %q", action)
	})
	if err != nil {
		return ProductBulkResult{}, err
	}
	return result, nil
}

// publicScope narrows a query to what an anonymous caller may see.
//
// Every public read below goes through it. Add a condition here and it applies
// to the list, the detail page and everything derived from them at once.
func (s *ProductService) publicScope(q *gorm.DB) *gorm.DB {
	q = q.Where("archived_at IS NULL")
	return q
}

// ListPublic returns one page of what the public surface may list: rows that are
// not archived. cfg is the allowlist the public handler keeps.
func (s *ProductService) ListPublic(ctx context.Context, params paginate.Params, cfg paginate.Config) (paginate.Result[models.Product], error) {
	query := s.publicScope(s.db(ctx).Model(&models.Product{}))
	return paginate.List[models.Product](query, params, cfg)
}

// GetPublic returns the product whose handle is key, if it is
// not archived.
//
// A row outside that is a 404 rather than a 403: whether it exists is itself
// not public.
func (s *ProductService) GetPublic(ctx context.Context, key string) (*models.Product, error) {
	var item models.Product
	if err := s.publicScope(s.db(ctx).Model(&models.Product{})).
		Where("handle = ?", key).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// RelatedPublic returns up to limit live products sharing item's
// collection, newest first, item excluded. A row with no collection has no
// siblings, so it gets the newest rows instead: a worse recommendation than a
// real match, and a better one than an empty strip.
func (s *ProductService) RelatedPublic(ctx context.Context, item *models.Product, limit int) ([]models.Product, error) {
	query := s.publicScope(s.db(ctx).Model(&models.Product{})).
		Where("id <> ?", item.ID)
	if item.CollectionID != "" {
		query = query.Where("collection_id = ?", item.CollectionID)
	}
	var rows []models.Product
	err := query.Order("created_at desc").Limit(limit).Find(&rows).Error
	return rows, err
}
