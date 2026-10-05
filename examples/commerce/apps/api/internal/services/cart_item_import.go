package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"commerce/apps/api/internal/money"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/imports"
	"commerce/apps/api/internal/models"
)

// StartImport records a CSV import of cart_items about to run: the job a client
// polls for progress.
func (s *CartItemService) StartImport(ctx context.Context, total int) (*models.ImportJob, error) {
	// CreatedBy is who may follow the job: GET /imports/:id answers only them,
	// or an admin.
	job := models.ImportJob{Resource: "cart_items", Status: "processing", Total: total, CreatedBy: authz.UserIDFrom(ctx)}
	if err := s.db(ctx).Create(&job).Error; err != nil {
		return nil, fmt.Errorf("starting the cart_items import: %w", err)
	}
	return &job, nil
}

// ImportCSV streams the CSV at path, creating cart_items in batches and
// updating the ImportJob jobID as it goes. belongs_to columns are resolved by
// their natural key (or id); unique-conflict rows are skipped; per-row failures
// are recorded. The file is removed when done.
//
// It outlives the request that started it, so give it a context that is not
// cancelled with the response. context.WithoutCancel of the request's keeps the
// caller, whom an owned resource's rows belong to, and the organization the
// multitenant plugin stamps them with.
func (s *CartItemService) ImportCSV(ctx context.Context, jobID, path string) {
	db := s.db(ctx)
	// record writes the job's state for the client polling it. A failure is
	// logged rather than dropped: the import carries on, and a job stuck on
	// its last state is at least explained.
	record := func(fields map[string]interface{}) {
		if err := db.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(fields).Error; err != nil {
			log.Printf("import job %s: recording progress: %v", jobID, err)
		}
	}
	defer os.Remove(path)
	// This runs in a bare goroutine, so gin.Recovery() does NOT cover it: an
	// unrecovered panic here would crash the whole server. Recover, and mark
	// the job failed so the client's poll terminates instead of hanging.
	defer func() {
		if r := recover(); r != nil {
			record(map[string]interface{}{
				"status":  "failed",
				"message": fmt.Sprintf("import crashed: %v", r),
			})
		}
	}()

	// Imports take turns: each holds a database connection and writes in
	// batches for its whole run, and a burst of them drained the pool every
	// request shares. See internal/imports.
	release := imports.Wait(func() {
		record(map[string]interface{}{"message": "Waiting for another import to finish"})
	})
	defer release()

	f, err := os.Open(path)
	if err != nil {
		record(map[string]interface{}{
			"status": "failed", "message": "could not reopen upload",
		})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		record(map[string]interface{}{
			"status": "failed", "message": "empty or invalid CSV",
		})
		return
	}
	idx := map[string]int{}
	for i, name := range headers {
		idx[strings.TrimSpace(strings.ToLower(name))] = i
	}
	get := func(rec []string, key string) (string, bool) {
		if i, ok := idx[key]; ok && i < len(rec) {
			return strings.TrimSpace(rec[i]), true
		}
		return "", false
	}

	// cartIDs remembers the Cart each token in the CSV resolved to, so a
	// distinct token costs one lookup rather than one per row. A missing Cart
	// is created; any other error fails the row instead of passing for missing.
	cartIDs := map[string]string{}
	resolveCart := func(v string) (string, error) {
		if id, ok := cartIDs[v]; ok {
			return id, nil
		}
		var rel models.Cart
		err := db.Where("token = ?", v).First(&rel).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rel = models.Cart{Token: v}
			if err = db.Create(&rel).Error; err != nil {
				// Another import may have created it in the meantime.
				if again := db.Where("token = ?", v).First(&rel).Error; again == nil {
					err = nil
				}
			}
		}
		if err != nil {
			return "", fmt.Errorf("cart %q: %w", v, err)
		}
		cartIDs[v] = rel.ID
		return rel.ID, nil
	}

	// productIDs remembers the Product each title in the CSV resolved to, so a
	// distinct title costs one lookup rather than one per row. A missing Product
	// is created; any other error fails the row instead of passing for missing.
	productIDs := map[string]string{}
	resolveProduct := func(v string) (string, error) {
		if id, ok := productIDs[v]; ok {
			return id, nil
		}
		var rel models.Product
		err := db.Where("title = ?", v).First(&rel).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rel = models.Product{Title: v}
			if err = db.Create(&rel).Error; err != nil {
				// Another import may have created it in the meantime.
				if again := db.Where("title = ?", v).First(&rel).Error; again == nil {
					err = nil
				}
			}
		}
		if err != nil {
			return "", fmt.Errorf("product %q: %w", v, err)
		}
		productIDs[v] = rel.ID
		return rel.ID, nil
	}

	created, skipped, failed := 0, 0, 0
	rowErrors := []map[string]interface{}{}

	// checkpoint writes current progress so the client's poll sees movement.
	checkpoint := func(status, message string) {
		errsJSON, _ := json.Marshal(rowErrors)
		record(map[string]interface{}{
			"status":    status,
			"processed": created + skipped + failed,
			"created":   created,
			"skipped":   skipped,
			"failed":    failed,
			"errors":    string(errsJSON),
			"message":   message,
		})
	}

	const batchSize = 200
	type pendingRow struct {
		item   models.CartItem
		rowNum int
	}
	batch := make([]pendingRow, 0, batchSize)

	// flush inserts the accumulated batch. CreateInBatches (with OnConflict
	// DoNothing) amortises the per-row fsync that makes large SQLite imports
	// crawl. If the whole batch errors (a bad row can poison it), we fall back
	// to per-row inserts so created/skipped/failed stay accurate and only the
	// offending row is dropped.
	flush := func() {
		if len(batch) == 0 {
			return
		}
		items := make([]models.CartItem, len(batch))
		for i := range batch {
			items[i] = batch[i].item
		}
		res := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(items, len(items))
		if res.Error == nil {
			created += int(res.RowsAffected)
			skipped += len(items) - int(res.RowsAffected)
		} else {
			for i := range batch {
				one := batch[i].item
				r := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&one)
				switch {
				case r.Error != nil:
					failed++
					if len(rowErrors) < 50 {
						rowErrors = append(rowErrors, map[string]interface{}{"row": batch[i].rowNum, "message": r.Error.Error()})
					}
				case r.RowsAffected == 0:
					skipped++
				default:
					created++
				}
			}
		}
		batch = batch[:0]
	}

	rowNum := 1 // header was row 1
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			failed++
			if len(rowErrors) < 50 {
				rowErrors = append(rowErrors, map[string]interface{}{"row": rowNum, "message": err.Error()})
			}
			continue
		}

		item := models.CartItem{}
		if v, ok := get(rec, "cart"); ok && v != "" {
			id, err := resolveCart(v)
			if err != nil {
				failed++
				if len(rowErrors) < 50 {
					rowErrors = append(rowErrors, map[string]interface{}{"row": rowNum, "message": err.Error()})
				}
				continue
			}
			item.CartID = id
		}
		if v, ok := get(rec, "product"); ok && v != "" {
			id, err := resolveProduct(v)
			if err != nil {
				failed++
				if len(rowErrors) < 50 {
					rowErrors = append(rowErrors, map[string]interface{}{"row": rowNum, "message": err.Error()})
				}
				continue
			}
			item.ProductID = id
		}
		if v, ok := get(rec, "variant_id"); ok {
			item.VariantID = v
		}
		if v, ok := get(rec, "quantity"); ok {
			n, _ := strconv.Atoi(v)
			item.Quantity = n
		}
		if v, ok := get(rec, "unit_price"); ok && v != "" {
			cur := "USD"
			if c, ok := get(rec, "unit_price_currency"); ok && c != "" {
				cur = c
			}
			major, _ := strconv.ParseFloat(v, 64)
			item.UnitPrice = money.FromMajor(major, cur)
		}
		if v, ok := get(rec, "title"); ok {
			item.Title = v
		}
		if v, ok := get(rec, "variant_label"); ok {
			item.VariantLabel = v
		}
		if v, ok := get(rec, "image_url"); ok {
			item.ImageURL = v
		}
		batch = append(batch, pendingRow{item: item, rowNum: rowNum})

		if len(batch) >= batchSize {
			flush()
			checkpoint("processing", "")
		}
	}
	flush()

	checkpoint("completed", fmt.Sprintf("Imported %d, skipped %d, failed %d", created, skipped, failed))
}
