package crypto

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type secretRow struct {
	ID    uint
	Notes EncryptedString
}

// Generated update and PATCH handlers write through maps of plain strings, and
// GORM only calls a column type's Value() when the value already has that
// type. Before Install the first edit to an encrypted column stored plaintext,
// and nothing showed it: reads still came back right, because Scan passes a
// value without the enc:v1: prefix straight through.
func TestMapUpdatesAreEncrypted(t *testing.T) {
	if err := InitFieldKey(base64.StdEncoding.EncodeToString(make([]byte, 32))); err != nil {
		t.Fatalf("key: %v", err)
	}
	t.Cleanup(func() { _ = InitFieldKey("") })

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := Install(db); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := db.AutoMigrate(&secretRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	row := secretRow{Notes: "created"}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	raw := func() string {
		var s string
		db.Raw("SELECT notes FROM secret_rows WHERE id = ?", row.ID).Scan(&s)
		return s
	}

	writes := []struct {
		name  string
		write func() error
	}{
		{"Updates(map)", func() error {
			return db.Model(&row).Updates(map[string]interface{}{"notes": "updated"}).Error
		}},
		{"Update(column)", func() error {
			return db.Model(&row).Update("notes", "patched").Error
		}},
	}
	for _, w := range writes {
		if err := w.write(); err != nil {
			t.Fatalf("%s: %v", w.name, err)
		}
		if got := raw(); !strings.HasPrefix(got, "enc:v1:") {
			t.Errorf("%s stored %q: an encrypted column written in plaintext", w.name, got)
		}
	}

	var back secretRow
	if err := db.First(&back, row.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(back.Notes) != "patched" {
		t.Errorf("read back %q, want the last value written", back.Notes)
	}
}

// Rows written before a key was set stay plaintext until something saves them
// again. EncryptExisting, which grit migrate runs, is that something, and a
// value it encrypts still reads back as it was.
func TestEncryptExistingEncryptsPlaintextRows(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&secretRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	for _, notes := range []string{"first", "second", ""} {
		if err := db.Create(&secretRow{Notes: EncryptedString(notes)}).Error; err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	if n, err := EncryptExisting(db, &secretRow{}); err != nil || n != 0 {
		t.Fatalf("without a key: wrote %d, %v; want nothing", n, err)
	}

	if err := InitFieldKey(base64.StdEncoding.EncodeToString(make([]byte, 32))); err != nil {
		t.Fatalf("key: %v", err)
	}
	t.Cleanup(func() { _ = InitFieldKey("") })

	n, err := EncryptExisting(db, &secretRow{})
	if err != nil {
		t.Fatalf("encrypt existing: %v", err)
	}
	if n != 2 {
		t.Errorf("encrypted %d values, want the 2 non-empty ones", n)
	}
	var raws []string
	db.Raw("SELECT notes FROM secret_rows WHERE notes <> '' ORDER BY id").Scan(&raws)
	for _, raw := range raws {
		if !strings.HasPrefix(raw, "enc:v1:") {
			t.Errorf("still plaintext after EncryptExisting: %q", raw)
		}
	}
	if again, err := EncryptExisting(db, &secretRow{}); err != nil || again != 0 {
		t.Errorf("a second run wrote %d, %v; want nothing", again, err)
	}

	var rows []secretRow
	if err := db.Order("id").Find(&rows).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(rows) != 3 || rows[0].Notes != "first" || rows[1].Notes != "second" || rows[2].Notes != "" {
		t.Errorf("read back %+v, want the values as they were written", rows)
	}
}
