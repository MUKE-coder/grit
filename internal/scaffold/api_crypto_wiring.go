package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// writeCryptoFiles writes internal/crypto into an upgraded project.
//
// Upgrade does not regenerate API code in general. This is the exception for
// the same reason as internal/money: the package gained crypto.Install, the
// fix for encrypted columns written in plaintext, and a project carrying the
// old copy keeps the leak. Manifest-guarded, so an edited copy is a conflict.
func writeCryptoFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "crypto", "field.go"):           apiCryptoFieldGo(),
		filepath.Join(apiRoot, "internal", "crypto", "field_test.go"):      apiCryptoFieldTestGo(),
		filepath.Join(apiRoot, "internal", "crypto", "map_update_test.go"): apiCryptoMapUpdateTestGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", opts.Module())
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

const fieldEncryptionConnectHook = `	// Map-based updates to an EncryptedString column are encrypted too. GORM only
	// runs a column type's Value() when the value already has that type, and the
	// generated update and PATCH handlers write maps of plain strings, so without
	// this the first edit to an encrypted column stored plaintext.
	if err := crypto.Install(db); err != nil {
		return nil, fmt.Errorf("installing field encryption for map updates: %w", err)
	}

`

// EnsureFieldEncryptionWiring makes the project call crypto.Install when it
// connects. Idempotent, and it refuses rather than guesses when the file has
// been reshaped past recognition.
func EnsureFieldEncryptionWiring(apiRoot, module string) error {
	path := filepath.Join(apiRoot, "internal", "database", "database.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("field encryption needs %s: %w", path, err)
	}
	crlf := strings.Contains(string(raw), "\r\n")
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	if strings.Contains(content, "crypto.Install(db)") {
		return nil
	}
	const anchor = "\tsqlDB, err := db.DB()"
	if !strings.Contains(content, anchor) {
		return fmt.Errorf("could not find where to install field encryption in %s.\n\n"+
			"Add this call after connecting, or map updates to encrypted columns store plaintext:\n\n  crypto.Install(db)", path)
	}
	content = strings.Replace(content, anchor, fieldEncryptionConnectHook+anchor, 1)
	var ok bool
	if content, ok = addImportGroup(content, module+"/internal/crypto"); !ok {
		return fmt.Errorf("could not add the crypto import to %s: add it by hand", path)
	}
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Wired field encryption into database.go")
	return nil
}

// apiCryptoMapUpdateTestGo is the regression test for the plaintext leak, and
// it ships with every project so the behaviour is checked where it runs.
func apiCryptoMapUpdateTestGo() string {
	return `package crypto

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
`
}
