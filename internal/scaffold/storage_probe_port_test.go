package scaffold

import (
	"strings"
	"testing"
)

// grit#97, filed from a real project: the storage probe asked MinIO on 9000
// while every project publishes it on MINIO_PORT, which defaults to 9002 so
// that several projects can run at once. MinIO was reported as not answering
// while it was answering, and every new project stored uploads on local disk
// in development without anybody choosing that. Nothing failed, which is why
// it survived.
func TestStorageProbeFollowsMinioPort(t *testing.T) {
	if strings.Contains(configResolveStorageDriverFunc, "localhost:9000") {
		t.Error("the storage probe hard-codes 9000, so it ignores MINIO_PORT")
	}
	if !strings.Contains(configResolveStorageDriverFunc, "resolveMinioEndpoint()") {
		t.Error("the storage probe should call resolveMinioEndpoint, which already " +
			"honours MINIO_ENDPOINT then MINIO_PORT")
	}
}

// The project that needs this most is an existing one, and it already has the
// fallback, so the early return in the repair would have skipped it.
func TestRepairFixesTheProbeOnAProjectThatAlreadyHasTheFallback(t *testing.T) {
	src := `package config

func resolveStorageDriver() string {
	driver := getEnv("STORAGE_DRIVER", "minio")
	endpoint := getEnv("MINIO_ENDPOINT", "http://localhost:9000")
	if minioAnswers(endpoint) {
		return driver
	}
	return "local"
}

func minioAnswers(endpoint string) bool { return false }
`
	out, changes, problems := repairStorageFallbackSource(src)
	if len(problems) > 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
	if len(changes) == 0 {
		t.Fatal("the probe was left asking for 9000 and nothing was reported")
	}
	if strings.Contains(out, "localhost:9000") {
		t.Error("the hard-coded probe survived the repair")
	}
	if !strings.Contains(out, "endpoint := resolveMinioEndpoint()") {
		t.Error("the repair did not install the port-aware probe")
	}
}

// And it stays idempotent: a project already repaired reports no change, or
// `grit upgrade` claims work it did not do on every run.
func TestRepairIsQuietOnAnAlreadyCorrectProbe(t *testing.T) {
	src := `package config

func resolveStorageDriver() string {
	endpoint := resolveMinioEndpoint()
	if minioAnswers(endpoint) {
		return "minio"
	}
	return "local"
}

func minioAnswers(endpoint string) bool { return false }
`
	out, changes, problems := repairStorageFallbackSource(src)
	if out != src {
		t.Error("a correct probe was rewritten")
	}
	if len(changes) != 0 || len(problems) != 0 {
		t.Errorf("expected silence, got changes=%v problems=%v", changes, problems)
	}
}
