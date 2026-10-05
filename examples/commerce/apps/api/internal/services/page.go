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

// PageService owns every database read and write for pages.
//
// The handler reads the request, calls one of these and writes the answer. A
// job, a command or a test calls the same methods with no request at all, so a
// rule enforced here is enforced everywhere, not only on the HTTP route.
//
// Each method takes the context it runs in. It carries the organization the
// multitenant plugin resolved, the actor ownership is scoped by (see
// authz.WithActor, and authz.AsSystem for a job), and the cancellation that
// fires when a client goes away.
type PageService struct {
	DB *gorm.DB
}

// pageListConfig is what a client may search, sort and filter pages
// by. Whitelisted, because each name ends up in SQL.
var pageListConfig = paginate.Config{
	Searchable: []string{"title", "handle", "body"},
	Sortable:   map[string]bool{"id": true, "created_at": true, "title": true, "handle": true, "body": true},
	Filterable: map[string]bool{"id": true, "title": true, "handle": true, "body": true},
}

// writablePage is every column Patch and Bulk may write. id, the
// timestamps and the version are the framework's, and are dropped.
var writablePage = map[string]bool{
	"title": true,
	"body":  true,
}

// db binds the database to ctx, so whatever a middleware put there reaches
// GORM's callbacks: the multitenant plugin scopes by it.
func (s *PageService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// write is a session for a single-statement write: no wrapping transaction,
// and RETURNING where the dialect has it. Safe only because every write it is
// used for is exactly one statement, which the generator knew.
func (s *PageService) write(db *gorm.DB) *gorm.DB {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if s.returning(db) {
		tx = tx.Clauses(clause.Returning{})
	}
	return tx
}

// returning reports whether the dialect hands back the written row. MySQL
// does not, and does not say so: the clause is dropped and the defaults come
// back empty, so there the row is read again.
func (s *PageService) returning(db *gorm.DB) bool {
	switch db.Dialector.Name() {
	case "postgres", "sqlite":
		return true
	}
	return false
}

// List returns one page of pages.
//
//	archived "true" or "1"   only archived rows
//	archived "all"           both
//	anything else            only live rows
func (s *PageService) List(ctx context.Context, p paginate.Params, archived string) (paginate.Result[models.Page], error) {
	query := s.db(ctx).Model(&models.Page{})

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

	return paginate.List[models.Page](query, p, pageListConfig)
}

// Export hands every matching page to each, a batch at a time, so a large
// table is never in memory at once. search matches the columns List searches.
//
// No ORDER BY of its own: FindInBatches pages by primary key, which is
// creation order for the time-ordered ids Grit issues, and a sort in front of
// that key repeated rows from the second batch on.
func (s *PageService) Export(ctx context.Context, search string, each func(rows []models.Page) error) error {
	query := s.db(ctx).Model(&models.Page{})
	if search != "" {
		clause := ""
		args := []any{}
		wild := "%" + search + "%"
		for i, col := range pageListConfig.Searchable {
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
	var rows []models.Page
	return query.FindInBatches(&rows, 1000, func(tx *gorm.DB, batch int) error {
		return each(rows)
	}).Error
}

// GetByID returns one page with the relations shown beside it.
func (s *PageService) GetByID(ctx context.Context, id string) (*models.Page, error) {
	var item models.Page
	if err := s.db(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if d := authz.Inspect(ctx, "pages.read", &item); !d.Allowed() {
		return nil, authz.Denied("pages.read", d)
	}
	return &item, nil
}

// load reads the row a write is about to change, without its relations.
//
// ability is what the caller is about to do, "update" or "delete", so a
// policy can refuse one without refusing the other.
func (s *PageService) load(ctx context.Context, id, ability string) (*models.Page, error) {
	var item models.Page
	if err := s.db(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// A policy rule narrows this further than the route's permission did.
	// None defined means no narrowing, which is how the project behaved
	// before it had policies. See internal/policies.
	if d := authz.Inspect(ctx, "pages."+ability, &item); !d.Allowed() {
		return nil, authz.Denied("pages."+ability, d)
	}
	return &item, nil
}

// conflict is the answer to a write whose precondition failed: the version
// the row is at now.
func (s *PageService) conflict(ctx context.Context, id string) error {
	var current models.Page
	if err := s.db(ctx).Select("version").First(&current, "id = ?", id).Error; err != nil {
		return err
	}
	return &concurrency.ErrConflict{Current: current.Version}
}

// Create saves a new page and fills item in as it was stored.
func (s *PageService) Create(ctx context.Context, item *models.Page) error {
	db := s.db(ctx)
	if err := s.write(db).Create(item).Error; err != nil {
		return err
	}
	if !s.returning(db) {
		if err := db.First(item, "id = ?", item.ID).Error; err != nil {
			return err
		}
	}
	return nil
}

// Update writes updates to one page. With a precondition it lands only
// if the row is still at that version, and otherwise returns an
// *concurrency.ErrConflict naming the version it is at.
func (s *PageService) Update(ctx context.Context, id string, updates map[string]interface{}, pre *concurrency.Precondition) (*models.Page, error) {
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
	if !s.returning(db) {
		if err := db.First(item, "id = ?", item.ID).Error; err != nil {
			return nil, err
		}
	}
	return item, nil
}

// Patch writes only the columns body names, leaving every other one as it
// is. Keys that are not writable columns are dropped. It returns the row and
// the columns it wrote.
func (s *PageService) Patch(ctx context.Context, id string, body map[string]interface{}, pre *concurrency.Precondition) (*models.Page, map[string]interface{}, error) {
	item, err := s.load(ctx, id, "update")
	if err != nil {
		return nil, nil, err
	}
	db := s.db(ctx)

	updates := map[string]interface{}{}
	for k, v := range body {
		if writablePage[k] {
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
	if !s.returning(db) {
		if err := db.First(item, "id = ?", item.ID).Error; err != nil {
			return nil, nil, err
		}
	}
	return item, updates, nil
}

// Delete soft-deletes one page and returns it as it was.
func (s *PageService) Delete(ctx context.Context, id string) (*models.Page, error) {
	item, err := s.load(ctx, id, "delete")
	if err != nil {
		return nil, err
	}
	if err := s.db(ctx).Delete(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// PageBulkResult is what a bulk action did.
type PageBulkResult struct {
	// IDs are the rows acted on: those requested that exist, that the caller
	// may touch, and that the action applies to.
	IDs []string
	// Updates are the columns a patch wrote, after the whitelist.
	Updates map[string]interface{}
}

// Bulk applies one action to many pages in a single transaction: all of
// it lands or none of it does. action is delete, archive, restore or patch.
func (s *PageService) Bulk(ctx context.Context, action string, ids []string, patch map[string]interface{}) (PageBulkResult, error) {
	db := s.db(ctx)
	var result PageBulkResult

	// Unarchived rows for archive, archived for restore: without it a mixed
	// selection reports "12 archived" having changed three.
	scope := db.Model(&models.Page{}).Where("id IN ?", ids)
	if action == "restore" {
		scope = scope.Where("archived_at IS NOT NULL")
	} else if action == "archive" {
		scope = scope.Where("archived_at IS NULL")
	}
	var items []models.Page
	if err := scope.Find(&items).Error; err != nil {
		return result, fmt.Errorf("loading pages: %w", err)
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
			if writablePage[k] {
				result.Updates[k] = v
			}
		}
		if len(result.Updates) == 0 {
			return PageBulkResult{}, respond.Rule("No writable fields in patch")
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		switch action {
		case "delete":
			return tx.Where("id IN ?", result.IDs).Delete(&models.Page{}).Error
		case "archive":
			return tx.Model(&models.Page{}).Where("id IN ?", result.IDs).
				Update("archived_at", time.Now()).Error
		case "restore":
			return tx.Model(&models.Page{}).Where("id IN ?", result.IDs).
				Update("archived_at", nil).Error
		case "patch":
			return tx.Model(&models.Page{}).Where("id IN ?", result.IDs).
				Updates(result.Updates).Error
		}
		return respond.Rule("unknown bulk action %q", action)
	})
	if err != nil {
		return PageBulkResult{}, err
	}
	return result, nil
}

// publicScope narrows a query to what an anonymous caller may see.
//
// Every public read below goes through it. Add a condition here and it applies
// to the list, the detail page and everything derived from them at once.
func (s *PageService) publicScope(q *gorm.DB) *gorm.DB {
	q = q.Where("archived_at IS NULL")
	return q
}

// ListPublic returns one page of what the public surface may list: rows that are
// not archived. cfg is the allowlist the public handler keeps.
func (s *PageService) ListPublic(ctx context.Context, params paginate.Params, cfg paginate.Config) (paginate.Result[models.Page], error) {
	query := s.publicScope(s.db(ctx).Model(&models.Page{}))
	return paginate.List[models.Page](query, params, cfg)
}

// GetPublic returns the page whose handle is key, if it is
// not archived.
//
// A row outside that is a 404 rather than a 403: whether it exists is itself
// not public.
func (s *PageService) GetPublic(ctx context.Context, key string) (*models.Page, error) {
	var item models.Page
	if err := s.publicScope(s.db(ctx).Model(&models.Page{})).
		Where("handle = ?", key).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
