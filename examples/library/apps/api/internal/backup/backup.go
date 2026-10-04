package backup

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
	"library/apps/api/internal/storage"
)

// KeepBackups is how many READY archives rolling retention keeps.
const KeepBackups = 4

// ErrStorageUnconfigured is returned when object storage isn't set up. The weekly
// cron treats it as a silent skip so local dev doesn't spam failures.
var ErrStorageUnconfigured = errors.New("object storage is not configured")

// Service produces full-database backups and uploads them to object storage.
type Service struct {
	DB      *gorm.DB
	Storage *storage.Storage
}

// DiskName is the storage disk archives go to when STORAGE_DISKS defines it.
// An archive holds every row of the database, so it belongs in a private
// bucket of its own rather than beside the files the app serves.
const DiskName = "backups"

// Store is where new archives go: the disk named DiskName when one is
// configured, and Storage otherwise.
func (s *Service) Store() *storage.Storage {
	if named := storage.Disks.Get(DiskName); named != nil {
		return named
	}
	return s.Storage
}

// stores are the places an archive may be: the backups disk, and the default
// store for archives taken before that disk was configured.
func (s *Service) stores() []*storage.Storage {
	var out []*storage.Storage
	store := s.Store()
	if store != nil {
		out = append(out, store)
	}
	if s.Storage != nil && s.Storage != store {
		out = append(out, s.Storage)
	}
	return out
}

// Locate is the store holding key, looking on the backups disk first.
func (s *Service) Locate(ctx context.Context, key string) *storage.Storage {
	for _, store := range s.stores() {
		if ok, err := store.Disk().Exists(ctx, key); err == nil && ok {
			return store
		}
	}
	return s.Store()
}

// Manifest is metadata.json — enough to verify a restore landed everything.
type Manifest struct {
	GeneratedAt time.Time      `json:"generated_at"`
	Tables      []string       `json:"tables"`
	RowCounts   map[string]int `json:"row_counts"`
	TotalRows   int            `json:"total_rows"`
}

// Tables returns every registered model's table name in registration order —
// parents before children, which is the order dump.sql must INSERT in for the
// foreign keys to hold.
//
// The list is derived from models.Models(), NEVER from user input, so a table
// name can't be injected into the raw SQL below. It also means every
// `grit generate resource` is automatically included in the next backup.
func Tables(db *gorm.DB) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	// joinTables are collected in a second bucket so they sort AFTER their
	// owning models — the m2m rows reference both sides, so both parent tables
	// must be inserted first on restore.
	var joinTables []string
	for _, m := range models.Models() {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(m); err != nil {
			return nil, fmt.Errorf("parsing model %T: %w", m, err)
		}
		if t := stmt.Schema.Table; t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
		// many_to_many fields create an implicit join table that is NOT in
		// models.Models(); without this its rows are silently dropped from
		// every backup (the relationship data vanishes on restore).
		for _, rel := range stmt.Schema.Relationships.Relations {
			if rel.JoinTable != nil {
				if t := rel.JoinTable.Table; t != "" && !seen[t] {
					seen[t] = true
					joinTables = append(joinTables, t)
				}
			}
		}
	}
	return append(out, joinTables...), nil
}

// streamTable reads one table row-at-a-time with raw database/sql, writing each
// row to the CSV entry AND an INSERT line to the dump buffer as it goes. Columns
// are scanned dynamically, so it works for any registered model (and join
// table). Crucially it holds ONE row in memory at a time — the whole table (let
// alone the whole database) is never materialised, which is what keeps the
// weekly backup from OOMing on a large database. Returns the row count.
//
// table comes from Tables() — the model registry, not user input — so it can't
// be injected into the raw SQL below.
func streamTable(ctx context.Context, sqlDB *sql.DB, table string, csvw *csv.Writer, dbuf *bufio.Writer) (int, error) {
	// #nosec G202 -- table comes from Tables(), the model registry, and never
	// from a request. A placeholder cannot carry a table name.
	rows, err := sqlDB.QueryContext(ctx, "SELECT * FROM \""+table+"\"")
	if err != nil {
		return 0, fmt.Errorf("select %s: %w", table, err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return 0, err
	}
	types := make([]string, len(colTypes))
	for i, ct := range colTypes {
		types[i] = strings.ToUpper(ct.DatabaseTypeName())
	}

	if err := csvw.Write(cols); err != nil {
		return 0, err
	}

	quoted := make([]string, len(cols))
	for i, c := range cols {
		quoted[i] = "\"" + c + "\""
	}
	prefix := "INSERT INTO \"" + table + "\" (" + strings.Join(quoted, ", ") + ") VALUES ("

	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	rec := make([]string, len(cols))
	sqlVals := make([]string, len(cols))

	count := 0
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return count, fmt.Errorf("scanning %s: %w", table, err)
		}
		// Both formatters consume each value immediately (into a string), so
		// the driver reusing its []byte buffer between rows is harmless — no
		// per-row copy needed.
		for i, v := range vals {
			rec[i] = csvFormat(v)
			sqlVals[i] = sqlFormat(v, types[i])
		}
		if err := csvw.Write(rec); err != nil {
			return count, err
		}
		if _, err := dbuf.WriteString(prefix + strings.Join(sqlVals, ", ") + ");\n"); err != nil {
			return count, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return count, err
	}
	if count > 0 {
		dbuf.WriteString("\n")
	}
	return count, nil
}

// csvFormat renders a scanned value for a spreadsheet. NULL becomes an empty cell.
func csvFormat(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case time.Time:
		return t.Format(time.RFC3339)
	case []byte:
		return string(t)
	case bool:
		return strconv.FormatBool(t)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// sqlFormat renders a scanned value as a SQL literal.
//
// The database type name matters: Postgres drivers hand back JSONB, TEXT and
// BYTEA all as []byte. Mislabelling one breaks the restore with errors like
// "invalid input syntax for type json", so we branch on the column type.
func sqlFormat(v any, dbType string) string {
	if v == nil {
		return "NULL"
	}
	switch t := v.(type) {
	case time.Time:
		return quote(t.Format(time.RFC3339Nano))
	case bool:
		if t {
			return "TRUE"
		}
		return "FALSE"
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []byte:
		switch {
		case strings.Contains(dbType, "BYTEA") || strings.Contains(dbType, "BLOB"):
			return "'\\x" + hex.EncodeToString(t) + "'"
		default: // JSON, JSONB, TEXT, VARCHAR — all arrive as bytes
			return quote(string(t))
		}
	case string:
		return quote(t)
	default:
		return quote(fmt.Sprintf("%v", v))
	}
}

// quote wraps a SQL string literal, doubling embedded single quotes. The restore
// splitter understands exactly this escaping and nothing fancier.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// ArchiveTo streams the backup ZIP to w, holding one row in memory at a time:
//
//	tables/<table>.csv — one per table, opens in any spreadsheet
//	dump.sql           — INSERTs parent->child, wrapped in BEGIN/COMMIT
//	metadata.json      — row counts, for verifying a restore
//
// A zip.Writer only allows one entry open at a time, but dump.sql spans every
// table while the CSV entries are per-table. So the SQL is streamed to a temp
// file alongside the CSVs, then copied into its own entry at the end. Nothing
// buffers the whole database — this is the memory-safe path the weekly cron and
// object-storage upload use.
func (s *Service) ArchiveTo(ctx context.Context, w io.Writer) (Manifest, error) {
	man := Manifest{GeneratedAt: time.Now().UTC(), RowCounts: map[string]int{}}

	tables, err := Tables(s.DB)
	if err != nil {
		return man, err
	}
	sqlDB, err := s.DB.DB()
	if err != nil {
		return man, err
	}

	dumpFile, err := os.CreateTemp("", "grit-dump-*.sql")
	if err != nil {
		return man, err
	}
	dumpPath := dumpFile.Name()
	defer os.Remove(dumpPath)
	defer dumpFile.Close()
	dbuf := bufio.NewWriter(dumpFile)

	dbuf.WriteString("-- Grit full-database backup\n")
	dbuf.WriteString("-- Restore: run migrations on an empty database, then replay this file.\n")
	dbuf.WriteString("--   grit restore backup.zip     (or)     psql \"$DATABASE_URL\" < dump.sql\n\n")
	dbuf.WriteString("BEGIN;\n\n")

	zw := zip.NewWriter(w)

	for _, table := range tables {
		cw, err := zw.Create("tables/" + table + ".csv")
		if err != nil {
			return man, err
		}
		csvw := csv.NewWriter(cw)

		count, err := streamTable(ctx, sqlDB, table, csvw, dbuf)
		if err != nil {
			return man, err
		}
		csvw.Flush()
		if err := csvw.Error(); err != nil {
			return man, err
		}

		man.RowCounts[table] = count
		man.TotalRows += count
	}

	dbuf.WriteString("COMMIT;\n")
	if err := dbuf.Flush(); err != nil {
		return man, err
	}
	if _, err := dumpFile.Seek(0, io.SeekStart); err != nil {
		return man, err
	}

	dw, err := zw.Create("dump.sql")
	if err != nil {
		return man, err
	}
	if _, err := io.Copy(dw, dumpFile); err != nil {
		return man, err
	}

	man.Tables = tables
	mw, err := zw.Create("metadata.json")
	if err != nil {
		return man, err
	}
	enc := json.NewEncoder(mw)
	enc.SetIndent("", "  ")
	if err := enc.Encode(man); err != nil {
		return man, err
	}

	return man, zw.Close()
}

// Archive builds the whole archive in memory. Prefer ArchiveTo for large
// databases; this convenience wrapper is fine for small/local use and keeps the
// original signature for any caller that already depends on it.
func (s *Service) Archive(ctx context.Context) ([]byte, Manifest, error) {
	var buf bytes.Buffer
	man, err := s.ArchiveTo(ctx, &buf)
	if err != nil {
		return nil, man, err
	}
	return buf.Bytes(), man, nil
}

// Start inserts the RUNNING row so callers can return it immediately and let the
// client poll while Run does the slow part.
func (s *Service) Start(kind string) (*models.Backup, error) {
	rec := &models.Backup{Kind: kind, Status: "RUNNING"}
	if err := s.DB.Create(rec).Error; err != nil {
		return nil, err
	}
	return rec, nil
}

// Run builds the archive, uploads it, flips the row to READY, then prunes old
// archives. A failed prune never fails the backup — the archive is already safe.
func (s *Service) Run(ctx context.Context, rec *models.Backup) error {
	store := s.Store()
	if store == nil {
		s.fail(rec, ErrStorageUnconfigured)
		return ErrStorageUnconfigured
	}

	// Stream the archive to a temp file (constant memory) and upload from it,
	// rather than building the whole ZIP in RAM.
	tmp, err := os.CreateTemp("", "grit-backup-*.zip")
	if err != nil {
		s.fail(rec, err)
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	defer tmp.Close()

	man, err := s.ArchiveTo(ctx, tmp)
	if err != nil {
		s.fail(rec, err)
		return err
	}
	size, err := tmp.Seek(0, io.SeekEnd)
	if err != nil {
		s.fail(rec, err)
		return err
	}
	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		s.fail(rec, err)
		return err
	}

	key := fmt.Sprintf("backups/%s-%s.zip", time.Now().UTC().Format("2006-01-02"), rec.ID)
	// Private, and refused if STORAGE_PUBLIC_PREFIXES would make backups/ public.
	if err := store.Disk().Put(ctx, key, tmp, storage.PutOptions{ContentType: "application/zip", Visibility: storage.VisibilityPrivate}); err != nil {
		s.fail(rec, err)
		return err
	}

	counts, _ := json.Marshal(man.RowCounts)
	now := time.Now()
	rec.Status = "READY"
	rec.StorageKey = key
	rec.SizeBytes = size
	rec.TableCount = len(man.Tables)
	rec.RowCount = man.TotalRows
	rec.RowCounts = string(counts)
	rec.CompletedAt = &now
	if err := s.DB.Save(rec).Error; err != nil {
		return err
	}

	if err := s.RollingCleanup(ctx, KeepBackups); err != nil {
		log.Printf("[backup] retention cleanup failed (archive is safe): %v", err)
	}
	return nil
}

// Generate is Start + Run, for callers that don't need the row up front (the
// weekly cron and the CLI).
func (s *Service) Generate(ctx context.Context, kind string) (*models.Backup, error) {
	rec, err := s.Start(kind)
	if err != nil {
		return nil, err
	}
	if err := s.Run(ctx, rec); err != nil {
		return rec, err
	}
	return rec, nil
}

func (s *Service) fail(rec *models.Backup, cause error) {
	msg := cause.Error()
	if len(msg) > 1000 {
		msg = msg[:1000]
	}
	now := time.Now()
	rec.Status = "FAILED"
	rec.Error = msg
	rec.CompletedAt = &now
	_ = s.DB.Save(rec).Error
}

// RollingCleanup keeps the newest `keep` READY archives and deletes the rest from
// object storage. Rows are marked PURGED rather than removed, so the audit trail
// still shows a backup ran that week.
func (s *Service) RollingCleanup(ctx context.Context, keep int) error {
	stores := s.stores()
	if len(stores) == 0 {
		return nil
	}
	var ready []models.Backup
	if err := s.DB.Where("status = ?", "READY").Order("created_at desc").Find(&ready).Error; err != nil {
		return err
	}
	if len(ready) <= keep {
		return nil
	}
	for _, b := range ready[keep:] {
		if b.StorageKey != "" {
			// From every store it may be on. A key already gone is not an error.
			for _, store := range stores {
				if err := store.Delete(ctx, b.StorageKey); err != nil {
					return fmt.Errorf("deleting %s: %w", b.StorageKey, err)
				}
			}
		}
		if err := s.DB.Model(&models.Backup{}).Where("id = ?", b.ID).
			Updates(map[string]any{"status": "PURGED", "storage_key": ""}).Error; err != nil {
			return err
		}
	}
	return nil
}

// ManualRateLimited reports whether a MANUAL backup was already taken inside the
// window. Weekly (cron) backups bypass this.
func (s *Service) ManualRateLimited(window time.Duration) (bool, error) {
	var count int64
	err := s.DB.Model(&models.Backup{}).
		Where("kind = ? AND created_at > ?", "MANUAL", time.Now().Add(-window)).
		Count(&count).Error
	return count > 0, err
}

// validFrequencies is the closed set of backup periods the UI offers.
var validFrequencies = map[string]bool{"daily": true, "weekly": true, "monthly": true, "yearly": true}

// GetSchedule returns the singleton backup-schedule row, seeding the default
// (weekly at 02:00, enabled) on first read.
func (s *Service) GetSchedule() (models.BackupSchedule, error) {
	var sc models.BackupSchedule
	err := s.DB.First(&sc, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		sc = models.BackupSchedule{ID: 1, Frequency: "weekly", Time: "02:00", Enabled: true}
		if cerr := s.DB.Create(&sc).Error; cerr != nil {
			return sc, cerr
		}
		return sc, nil
	}
	return sc, err
}

// SaveSchedule validates and persists the backup schedule (upsert of the
// singleton row).
func (s *Service) SaveSchedule(frequency, tod string, enabled bool) (models.BackupSchedule, error) {
	if !validFrequencies[frequency] {
		return models.BackupSchedule{}, fmt.Errorf("invalid frequency %q", frequency)
	}
	if _, _, err := parseHHMM(tod); err != nil {
		return models.BackupSchedule{}, fmt.Errorf("invalid time %q (want HH:MM)", tod)
	}
	sc := models.BackupSchedule{ID: 1, Frequency: frequency, Time: tod, Enabled: enabled}
	if err := s.DB.Save(&sc).Error; err != nil {
		return sc, err
	}
	return sc, nil
}

// DueNow reports whether an automatic backup should run at time `now`. It's
// called on every scheduler tick: a backup is due when the schedule is enabled,
// the current period's scheduled time has passed, and no SCHEDULED backup has
// been taken yet in this period (which also makes a missed run catch up on the
// next tick).
func (s *Service) DueNow(now time.Time) (bool, error) {
	sc, err := s.GetSchedule()
	if err != nil {
		return false, err
	}
	if !sc.Enabled {
		return false, nil
	}
	anchor := scheduleAnchor(now, sc.Frequency, sc.Time)
	if now.Before(anchor) {
		return false, nil
	}
	var last models.Backup
	err = s.DB.Where("kind = ? AND status IN ? AND created_at >= ?",
		"SCHEDULED", []string{"READY", "RUNNING"}, anchor).First(&last).Error
	if err == nil {
		return false, nil // already ran this period
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

// parseHHMM parses a "HH:MM" 24-hour time-of-day.
func parseHHMM(s string) (int, int, error) {
	var h, m int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
		return 0, 0, err
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, fmt.Errorf("out of range")
	}
	return h, m, nil
}

// scheduleAnchor returns the most recent scheduled datetime for the current
// period (server-local): today for daily, this week's Sunday for weekly, the
// 1st for monthly, Jan 1 for yearly — each at the configured time-of-day.
func scheduleAnchor(now time.Time, frequency, tod string) time.Time {
	h, m, err := parseHHMM(tod)
	if err != nil {
		h, m = 2, 0
	}
	y, mon, d := now.Date()
	loc := now.Location()
	switch frequency {
	case "daily":
		return time.Date(y, mon, d, h, m, 0, 0, loc)
	case "monthly":
		return time.Date(y, mon, 1, h, m, 0, 0, loc)
	case "yearly":
		return time.Date(y, time.January, 1, h, m, 0, 0, loc)
	default: // weekly — rewind to this week's Sunday
		start := time.Date(y, mon, d, h, m, 0, 0, loc)
		return start.AddDate(0, 0, -int(now.Weekday()))
	}
}
