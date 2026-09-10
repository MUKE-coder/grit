package scaffold

import (
	"strings"
	"testing"
)

// Restoring a backup failed on every project with --items (children replayed
// before the parents their foreign keys name) and on every project with
// --append-only (restore begins with TRUNCATE, which the table refuses).
func TestRestoreReplaysInDependencyOrderWithTriggersSuspended(t *testing.T) {
	src := backupRestoreGo()
	for _, want := range []string{
		"for _, s := range replayOrder(tx, stmts) {",
		"case schema.BelongsTo:",
		"resume, err := appendonly.Suspend(tx)",
		"return resume()",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("restore is missing %s", want)
		}
	}
	if !strings.Contains(appendOnlyGo(), "func Suspend(tx *gorm.DB) (func() error, error) {") {
		t.Error("appendonly has no way to let a restore through")
	}
}

// Generated update and PATCH handlers wrote plain strings into encrypted
// columns. Every new project installs the fix when it connects.
func TestConnectInstallsFieldEncryption(t *testing.T) {
	if !strings.Contains(apiDatabaseGo(), "crypto.Install(db)") {
		t.Error("Connect never installs the map-update encryption")
	}
	crypto := apiCryptoFieldGo()
	for _, want := range []string{
		"func Install(db *gorm.DB) error {",
		"stmt.Schema.LookUpField(key)",
		"values[key] = EncryptedString(v)",
	} {
		if !strings.Contains(crypto, want) {
			t.Errorf("the crypto package is missing %s", want)
		}
	}
}
