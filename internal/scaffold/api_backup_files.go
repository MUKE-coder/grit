package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeBackupFiles scaffolds the full-database backup subsystem:
//
//	models/backup.go       — the index row for each archive
//	internal/backup/*.go   — archive builder + restore replayer
//	handlers/backup.go     — list / generate / signed-URL download
//	cmd/backup/main.go     — `grit backup`  (upload, or --output a local file)
//	cmd/restore/main.go    — `grit restore` (migrate + replay a dump)
//
// Templates are written with "~" where a backtick belongs (Go raw strings can't
// contain backticks) and swapped back below.
func writeBackupFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "backup.go"):   backupModelGo(),
		filepath.Join(apiRoot, "internal", "backup", "backup.go"):   backupServiceGo(),
		filepath.Join(apiRoot, "internal", "backup", "restore.go"):  backupRestoreGo(),
		filepath.Join(apiRoot, "internal", "handlers", "backup.go"): backupHandlerGo(),
		filepath.Join(apiRoot, "cmd", "backup", "main.go"):          backupCmdMainGo(),
		filepath.Join(apiRoot, "cmd", "restore", "main.go"):         restoreCmdMainGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "~", "`")
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// backupModelGo is the index row for one archive. The archive itself lives in
// object storage; this table is what the API, admin panel and CLI list.
func backupModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Backup indexes one full-database snapshot. The archive lives in object storage
// (R2 / S3 / MinIO) under the "backups/" prefix; this row is what the admin UI,
// the REST API and the CLI read. Lifecycle:
//
//	RUNNING -> READY | FAILED -> PURGED
//
// PURGED means rolling retention deleted the object but we kept the row, so the
// audit trail still shows a backup existed on that date.
type Backup struct {
	ID          string     ~gorm:"primarykey;size:36" json:"id"~
	Kind        string     ~gorm:"size:20;index" json:"kind"~   // WEEKLY | MANUAL | CLI
	Status      string     ~gorm:"size:20;index" json:"status"~ // RUNNING | READY | FAILED | PURGED
	StorageKey  string     ~gorm:"size:512" json:"-"~
	SizeBytes   int64      ~json:"size_bytes"~
	TableCount  int        ~json:"table_count"~
	RowCount    int        ~json:"row_count"~
	RowCounts   string     ~gorm:"type:text" json:"-"~ // JSON map of table -> rows
	Error       string     ~gorm:"size:1000" json:"error,omitempty"~
	CreatedAt   time.Time  ~json:"created_at"~
	CompletedAt *time.Time ~json:"completed_at,omitempty"~
}

// BeforeCreate assigns a UUID so backup ids are opaque in download URLs.
func (m *Backup) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Status == "" {
		m.Status = "RUNNING"
	}
	return nil
}
`
}

// backupServiceGo builds the archive. Table names come from the model registry —
// never from user input — which is the single most important security property
// here: a dynamic table name would leak the whole database.
func backupServiceGo() string {
	return `package backup

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/storage"
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

// Manifest is metadata.json — enough to verify a restore landed everything.
type Manifest struct {
	GeneratedAt time.Time      ~json:"generated_at"~
	Tables      []string       ~json:"tables"~
	RowCounts   map[string]int ~json:"row_counts"~
	TotalRows   int            ~json:"total_rows"~
}

// Tables returns every registered model's table name in registration order —
// parents before children, which is the order dump.sql must INSERT in for the
// foreign keys to hold.
//
// The list is derived from models.Models(), NEVER from user input, so a table
// name can't be injected into the raw SQL below. It also means every
// ~grit generate resource~ is automatically included in the next backup.
func Tables(db *gorm.DB) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, m := range models.Models() {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(m); err != nil {
			return nil, fmt.Errorf("parsing model %T: %w", m, err)
		}
		if t := stmt.Schema.Table; t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out, nil
}

type tableData struct {
	Name    string
	Columns []string
	Types   []string
	Rows    [][]any
}

// dumpTable reads a whole table with raw database/sql, scanning columns
// dynamically. No struct mapping, so it works for any registered model.
func dumpTable(ctx context.Context, sqlDB *sql.DB, table string) (*tableData, error) {
	// table comes from Tables() — the model registry, not user input.
	rows, err := sqlDB.QueryContext(ctx, "SELECT * FROM \"" + table + "\"")
	if err != nil {
		return nil, fmt.Errorf("select %s: %w", table, err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	types := make([]string, len(colTypes))
	for i, ct := range colTypes {
		types[i] = strings.ToUpper(ct.DatabaseTypeName())
	}

	td := &tableData{Name: table, Columns: cols, Types: types}
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scanning %s: %w", table, err)
		}
		// Copy []byte: drivers may reuse the underlying buffer between rows.
		row := make([]any, len(cols))
		for i, v := range vals {
			if b, ok := v.([]byte); ok {
				cp := make([]byte, len(b))
				copy(cp, b)
				row[i] = cp
			} else {
				row[i] = v
			}
		}
		td.Rows = append(td.Rows, row)
	}
	return td, rows.Err()
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

// Archive builds the backup ZIP in memory:
//
//	tables/<table>.csv — one per table, opens in any spreadsheet
//	dump.sql           — INSERTs parent->child, wrapped in BEGIN/COMMIT
//	metadata.json      — row counts, for verifying a restore
func (s *Service) Archive(ctx context.Context) ([]byte, Manifest, error) {
	man := Manifest{GeneratedAt: time.Now().UTC(), RowCounts: map[string]int{}}

	tables, err := Tables(s.DB)
	if err != nil {
		return nil, man, err
	}
	sqlDB, err := s.DB.DB()
	if err != nil {
		return nil, man, err
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	var dump strings.Builder
	dump.WriteString("-- Grit full-database backup\n")
	dump.WriteString("-- Restore: run migrations on an empty database, then replay this file.\n")
	dump.WriteString("--   grit restore backup.zip     (or)     psql \"$DATABASE_URL\" < dump.sql\n\n")
	dump.WriteString("BEGIN;\n\n")

	for _, table := range tables {
		td, err := dumpTable(ctx, sqlDB, table)
		if err != nil {
			return nil, man, err
		}

		// CSV — one file per table.
		w, err := zw.Create("tables/" + table + ".csv")
		if err != nil {
			return nil, man, err
		}
		cw := csv.NewWriter(w)
		if err := cw.Write(td.Columns); err != nil {
			return nil, man, err
		}
		for _, r := range td.Rows {
			rec := make([]string, len(r))
			for i, v := range r {
				rec[i] = csvFormat(v)
			}
			if err := cw.Write(rec); err != nil {
				return nil, man, err
			}
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			return nil, man, err
		}

		// SQL — INSERTs in registration (parent -> child) order.
		if len(td.Rows) > 0 {
			quoted := make([]string, len(td.Columns))
			for i, c := range td.Columns {
				quoted[i] = "\"" + c + "\""
			}
			prefix := "INSERT INTO \"" + table + "\" (" + strings.Join(quoted, ", ") + ") VALUES ("
			for _, r := range td.Rows {
				vals := make([]string, len(r))
				for i, v := range r {
					vals[i] = sqlFormat(v, td.Types[i])
				}
				dump.WriteString(prefix + strings.Join(vals, ", ") + ");\n")
			}
			dump.WriteString("\n")
		}

		man.RowCounts[table] = len(td.Rows)
		man.TotalRows += len(td.Rows)
	}

	man.Tables = tables
	dump.WriteString("COMMIT;\n")

	dw, err := zw.Create("dump.sql")
	if err != nil {
		return nil, man, err
	}
	if _, err := dw.Write([]byte(dump.String())); err != nil {
		return nil, man, err
	}

	mw, err := zw.Create("metadata.json")
	if err != nil {
		return nil, man, err
	}
	enc := json.NewEncoder(mw)
	enc.SetIndent("", "  ")
	if err := enc.Encode(man); err != nil {
		return nil, man, err
	}

	if err := zw.Close(); err != nil {
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
	if s.Storage == nil {
		s.fail(rec, ErrStorageUnconfigured)
		return ErrStorageUnconfigured
	}

	data, man, err := s.Archive(ctx)
	if err != nil {
		s.fail(rec, err)
		return err
	}

	key := fmt.Sprintf("backups/%s-%s.zip", time.Now().UTC().Format("2006-01-02"), rec.ID)
	if err := s.Storage.Upload(ctx, key, bytes.NewReader(data), "application/zip"); err != nil {
		s.fail(rec, err)
		return err
	}

	counts, _ := json.Marshal(man.RowCounts)
	now := time.Now()
	rec.Status = "READY"
	rec.StorageKey = key
	rec.SizeBytes = int64(len(data))
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

// RollingCleanup keeps the newest ~keep~ READY archives and deletes the rest from
// object storage. Rows are marked PURGED rather than removed, so the audit trail
// still shows a backup ran that week.
func (s *Service) RollingCleanup(ctx context.Context, keep int) error {
	if s.Storage == nil {
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
			if err := s.Storage.Delete(ctx, b.StorageKey); err != nil {
				return fmt.Errorf("deleting %s: %w", b.StorageKey, err)
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
`
}

// backupRestoreGo replays an archive. Restore is the real deliverable — a backup
// you've never restored is a rumour — so it ships as a first-class command.
func backupRestoreGo() string {
	return `package backup

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"gorm.io/gorm"
)

// SplitStatements splits our generated dump.sql into executable statements.
//
// This is NOT a general SQL parser — it doesn't need to be. We generate the file
// ourselves and only ever emit numbers, NULL/TRUE/FALSE, and single-quoted
// literals with '' escaping. Tracking quote state is therefore exact. Splitting
// naively on ";" would corrupt any value containing a semicolon.
func SplitStatements(script string) []string {
	var out []string
	var cur strings.Builder
	inString := false

	rs := []rune(script)
	for i := 0; i < len(rs); i++ {
		c := rs[i]
		if inString {
			cur.WriteRune(c)
			if c == '\'' {
				// '' is an escaped quote, not the end of the literal.
				if i+1 < len(rs) && rs[i+1] == '\'' {
					cur.WriteRune('\'')
					i++
					continue
				}
				inString = false
			}
			continue
		}
		switch c {
		case '\'':
			inString = true
			cur.WriteRune(c)
		case ';':
			if s := strings.TrimSpace(cur.String()); s != "" {
				out = append(out, s)
			}
			cur.Reset()
		default:
			cur.WriteRune(c)
		}
	}
	if s := strings.TrimSpace(cur.String()); s != "" {
		out = append(out, s)
	}
	return out
}

// stripComments drops our "--" header lines before splitting.
func stripComments(script string) string {
	var b strings.Builder
	for _, line := range strings.Split(script, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

// Restore replays a backup archive into the connected database inside a single
// transaction: either every row lands or nothing does.
//
// The archive carries DATA, not schema — run migrations on the target database
// first (cmd/restore does this for you).
func Restore(db *gorm.DB, zipPath string) (Manifest, error) {
	var man Manifest

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return man, fmt.Errorf("opening %s: %w", zipPath, err)
	}
	defer zr.Close()

	var dump string
	for _, f := range zr.File {
		switch f.Name {
		case "dump.sql", "metadata.json":
			rc, err := f.Open()
			if err != nil {
				return man, err
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return man, err
			}
			if f.Name == "dump.sql" {
				dump = string(data)
			} else {
				_ = json.Unmarshal(data, &man)
			}
		}
	}
	if dump == "" {
		return man, errors.New("dump.sql not found in archive")
	}

	stmts := SplitStatements(stripComments(dump))
	return man, db.Transaction(func(tx *gorm.DB) error {
		for _, s := range stmts {
			switch strings.ToUpper(strings.TrimSpace(s)) {
			case "BEGIN", "COMMIT":
				continue // we own the transaction
			}
			if err := tx.Exec(s).Error; err != nil {
				head := s
				if len(head) > 120 {
					head = head[:120] + "..."
				}
				return fmt.Errorf("executing %q: %w", head, err)
			}
		}
		return nil
	})
}
`
}

// backupHandlerGo exposes list / generate / download over REST. Mounted on the
// admin group — backups are an operator feature.
func backupHandlerGo() string {
	return `package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/backup"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/storage"
)

const (
	// manualBackupWindow rate-limits on-demand backups. The weekly cron bypasses it.
	manualBackupWindow = 24 * time.Hour
	// downloadURLTTL: long enough for a slow phone, short enough that a leaked
	// link stops working before anyone can use it.
	downloadURLTTL = 15 * time.Minute
	// backupTimeout bounds a single run so a hung upload can't wedge the worker.
	backupTimeout = 30 * time.Minute
)

// BackupHandler serves the full-database backup index.
type BackupHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
}

func (h *BackupHandler) svc() *backup.Service {
	return &backup.Service{DB: h.DB, Storage: h.Storage}
}

// List returns backups newest-first. Poll it while one is RUNNING.
func (h *BackupHandler) List(c *gin.Context) {
	var items []models.Backup
	if err := h.DB.Order("created_at desc").Limit(50).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to list backups"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// Generate starts a manual backup in the background and returns the RUNNING row
// immediately — a full dump can take a while. Poll List until it flips to READY.
func (h *BackupHandler) Generate(c *gin.Context) {
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "STORAGE_UNAVAILABLE", "message": "Object storage is not configured"},
		})
		return
	}

	limited, err := h.svc().ManualRateLimited(manualBackupWindow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to check rate limit"},
		})
		return
	}
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{"code": "RATE_LIMITED", "message": "A manual backup was already taken in the last 24 hours"},
		})
		return
	}

	rec, err := h.svc().Start("MANUAL")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to start backup"},
		})
		return
	}

	go func(r models.Backup) {
		ctx, cancel := context.WithTimeout(context.Background(), backupTimeout)
		defer cancel()
		_ = h.svc().Run(ctx, &r)
	}(*rec)

	c.JSON(http.StatusAccepted, gin.H{"data": rec, "message": "Backup started"})
}

// Download mints a short-lived pre-signed URL so the client pulls the archive
// straight from object storage — no proxying a multi-hundred-MB file through the API.
func (h *BackupHandler) Download(c *gin.Context) {
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "STORAGE_UNAVAILABLE", "message": "Object storage is not configured"},
		})
		return
	}

	var b models.Backup
	if err := h.DB.First(&b, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Backup not found"},
		})
		return
	}
	if b.Status != "READY" || b.StorageKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "NOT_AVAILABLE", "message": "This backup is not available for download"},
		})
		return
	}

	url, err := h.Storage.GetSignedURL(c.Request.Context(), b.StorageKey, downloadURLTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to sign download URL"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"url": url, "expires_in": int(downloadURLTTL.Seconds())},
	})
}
`
}

// backupCmdMainGo powers ~grit backup~.
func backupCmdMainGo() string {
	return `package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"{{MODULE}}/internal/backup"
	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/storage"
)

// Backs up every registered model to a ZIP (CSV per table + dump.sql +
// metadata.json). By default it uploads to object storage and records the row;
// --output writes a local file instead and touches nothing else.
func main() {
	out := flag.String("output", "", "Write the archive to this local path instead of uploading")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	svc := &backup.Service{DB: db}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if *out != "" {
		data, man, err := svc.Archive(ctx)
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		if err := os.WriteFile(*out, data, 0o644); err != nil {
			log.Fatalf("Failed to write %s: %v", *out, err)
		}
		fmt.Printf("Backup written to %s — %d tables, %d rows, %.1f KB\n",
			*out, len(man.Tables), man.TotalRows, float64(len(data))/1024)
		return
	}

	st, err := storage.New(cfg.Storage)
	if err != nil {
		log.Fatalf("Object storage is not configured: %v\n(use --output <file> to write a local archive)", err)
	}
	svc.Storage = st

	rec, err := svc.Generate(ctx, "CLI")
	if err != nil {
		log.Fatalf("Backup failed: %v", err)
	}
	fmt.Printf("Backup %s uploaded — %d tables, %d rows, %.1f KB\n",
		rec.ID, rec.TableCount, rec.RowCount, float64(rec.SizeBytes)/1024)
}
`
}

// restoreCmdMainGo powers ~grit restore~ — the path you must test before you
// trust any of this.
func restoreCmdMainGo() string {
	return `package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"{{MODULE}}/internal/backup"
	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/models"
)

// Replays a backup archive into the configured database. Runs migrations first
// (the archive carries data, not schema), then executes dump.sql in one
// transaction: every row lands, or none does.
func main() {
	migrate := flag.Bool("migrate", true, "Run migrations before restoring")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage: restore [--migrate=false] <backup.zip>")
	}
	path := flag.Arg(0)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *migrate {
		fmt.Println("Running migrations...")
		if err := models.Migrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}

	fmt.Printf("Restoring %s ...\n", path)
	man, err := backup.Restore(db, path)
	if err != nil {
		log.Fatalf("Restore failed: %v", err)
	}

	fmt.Printf("Restored %d tables, %d rows (archive generated %s)\n",
		len(man.Tables), man.TotalRows, man.GeneratedAt.Format(time.RFC3339))
}
`
}
