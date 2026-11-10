package scaffold

import (
	"strings"
	"testing"
)

// --storage reaches .env, which is the only place it has any effect.
//
// The picker and the flag agreeing is one thing (internal/prompt tests that);
// the chosen value being written is another, and it is the half that decides
// whether the feature does anything at all.
func TestTheChosenStorageDriverIsWrittenToEnv(t *testing.T) {
	for _, driver := range StorageDriverOrder {
		opts := Options{ProjectName: "shop", DBProvider: "postgres", StorageDriver: driver}
		env := envFile(opts)

		want := "\nSTORAGE_DRIVER=" + driver + "\n"
		if !strings.Contains(env, want) {
			t.Errorf("--storage %s did not write %q into .env", driver, strings.TrimSpace(want))
		}
		// Exactly one active line. The file also carries commented examples
		// for the other providers, and two uncommented ones would mean the
		// last wins, silently.
		if n := strings.Count(env, "\nSTORAGE_DRIVER="); n != 1 {
			t.Errorf("--storage %s left %d active STORAGE_DRIVER lines in .env", driver, n)
		}
	}
}

// An unset flag keeps the default every existing project has.
//
// This is the compatibility check: scaffolding with no --storage must produce
// the same .env it produced before the flag existed, or the flag is a silent
// change to what `grit new` does for everybody who does not pass it.
func TestNoStorageChoiceKeepsMinio(t *testing.T) {
	opts := Options{ProjectName: "shop", DBProvider: "postgres"}
	if env := envFile(opts); !strings.Contains(env, "\nSTORAGE_DRIVER=minio\n") {
		t.Error("a project scaffolded without --storage no longer defaults to minio")
	}
}

// The spellings people reach for, and the one that must be refused.
func TestValidateStorageDriverNormalises(t *testing.T) {
	for given, want := range map[string]string{
		"":           "minio",
		"local":      "local",
		"Local":      "local",
		"  disk  ":   "local",
		"filesystem": "local",
		"aws":        "s3",
		"cloudflare": "r2",
		"backblaze":  "b2",
	} {
		opts := Options{StorageDriver: given}
		if err := opts.ValidateStorageDriver(); err != nil {
			t.Errorf("--storage %q was refused: %v", given, err)
			continue
		}
		if opts.StorageDriver != want {
			t.Errorf("--storage %q became %q, want %q", given, opts.StorageDriver, want)
		}
	}

	opts := Options{StorageDriver: "dropbox"}
	if err := opts.ValidateStorageDriver(); err == nil {
		t.Error("--storage dropbox was accepted, and would be written to .env as a driver " +
			"storage.Open refuses, which disables uploads with no error at scaffold time")
	}
}
