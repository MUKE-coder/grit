package scaffold

import (
	"strings"
	"testing"
)

// A new project must take an upload without Docker running.
//
// resolveStorageDriver falls back to the local disk outside production so that
// "a new project stores uploads before Docker is running instead of answering
// each one with STORAGE_UNAVAILABLE". It tested for a missing MINIO_ACCESS_KEY,
// and the scaffold writes that key into .env, so the condition could never be
// true and the fallback never fired: every new project on a machine without
// Docker answered every upload with "File storage is not configured" until the
// developer discovered STORAGE_DRIVER=local for themselves.
func TestStorageFallsBackWhenMinIOIsNotAnswering(t *testing.T) {
	src := configResolveStorageDriverFunc

	// The question is whether MinIO answers.
	if !strings.Contains(src, "func minioAnswers(endpoint string) bool {") {
		t.Fatal("there is no reachability check, so the fallback still cannot fire " +
			"for a project whose MINIO_ACCESS_KEY the scaffold filled in")
	}
	if !strings.Contains(src, "if minioAnswers(endpoint) {") {
		t.Error("resolveStorageDriver does not consult minioAnswers")
	}
	if !strings.Contains(src, "net.DialTimeout(") {
		t.Error("minioAnswers does not actually try to connect")
	}

	// A missing key still falls back: that path was right, it was just never
	// the only one that mattered.
	if !strings.Contains(src, `getEnv("MINIO_ACCESS_KEY", "") == ""`) {
		t.Error("the missing-credentials fallback was dropped")
	}

	// Production never falls back, and nothing but minio does. An s3, r2 or b2
	// endpoint that is briefly unreachable is not a reason to start writing to
	// a local disk.
	if !strings.Contains(src, `if driver != "minio" || getEnv("APP_ENV", "production") == "production" {`) {
		t.Error("the fallback is not confined to minio outside production")
	}

	// And it says what happened, with both ways out.
	for _, phrase := range []string{
		"is not answering",
		"docker compose up -d minio",
		"STORAGE_DRIVER=local",
	} {
		if !strings.Contains(src, phrase) {
			t.Errorf("the log line does not mention %q", phrase)
		}
	}
}

// APP_URL naming a different port from the one the server listens on is always
// a mistake on localhost, and a silent one: every URL built from it is wrong,
// and with STORAGE_DRIVER=local that means every uploaded file's URL points at
// a port where nothing is listening.
func TestAppURLPortMismatchIsReported(t *testing.T) {
	src := apiConfigGo()

	if !strings.Contains(src, "func resolveAppURL(port string) string {") {
		t.Fatal("APP_URL is read without checking it against the listening port")
	}
	if !strings.Contains(src, "AppURL:             resolveAppURL(") &&
		!strings.Contains(src, "AppURL:      resolveAppURL(") {
		t.Error("the config does not resolve APP_URL through resolveAppURL")
	}

	// Only on localhost. Behind a proxy, APP_URL is the public address and a
	// different port is correct, which is the normal production shape.
	for _, host := range []string{`"localhost"`, `"127.0.0.1"`} {
		if !strings.Contains(src, host) {
			t.Errorf("resolveAppURL does not confine the warning to %s", host)
		}
	}
	if !strings.Contains(src, "uploaded files included") {
		t.Error("the warning does not say what breaks")
	}
}

// The error a developer actually sees has to name the fix.
func TestStorageUnavailableNamesTheFix(t *testing.T) {
	src := uploadHandlerGo()
	const want = "Set STORAGE_DRIVER=local in .env"
	if !strings.Contains(src, want) {
		t.Errorf("the storage-unavailable message does not say %q, so a developer "+
			"with no object storage is told what is wrong and not what to do", want)
	}
}

// Both storage repairs must turn an old config.go into exactly the template.
//
// config.go is repaired on upgrade rather than rewritten, so a function added
// to the template reaches new projects only unless a repair puts it there, and
// the project that needs these two is by definition an existing one: a new one
// has APP_URL and APP_PORT agreeing, and a new one's fallback already works.
func TestStorageConfigRepairsProduceTheTemplate(t *testing.T) {
	current := strings.ReplaceAll(apiConfigGo(), "{{MODULE}}", "example.com/app")

	t.Run("the MinIO fallback", func(t *testing.T) {
		// A project on the version whose fallback could never fire.
		old := strings.Replace(current, configResolveStorageDriverFunc, oldResolveStorageDriver, 1)
		if old == current {
			t.Fatal("the reconstruction did not change the template")
		}
		checkPerfLowRepair(t, "config.go", current, old, repairStorageFallbackSource)
	})

	t.Run("the APP_URL port check", func(t *testing.T) {
		// Formatted first: a real project's config.go is gofmt'd, so the struct
		// field the repair anchors on is aligned there and not in the template.
		current := formatPerfLowGo(t, "config.go", current)

		// A project from before resolveAppURL existed.
		start := strings.Index(current, "// resolveAppURL reads APP_URL")
		end := strings.Index(current, "// warnProviderMismatch says so when")
		if start < 0 || end < 0 || end < start {
			t.Fatal("resolveAppURL is not where the template puts it")
		}
		old := current[:start] + current[end:]
		old = strings.Replace(old,
			`AppURL:             resolveAppURL(firstNonEmpty(os.Getenv("PORT"), os.Getenv("APP_PORT"), "8080")),`,
			`AppURL:             getEnv("APP_URL", "http://localhost:8080"),`, 1)
		checkPerfLowRepair(t, "config.go", current, old, repairAppURLPortSource)
	})
}

// oldResolveStorageDriver is the version whose fallback could never fire,
// kept here so the repair above has something real to be tested against.
const oldResolveStorageDriver = `// resolveStorageDriver picks the storage driver from STORAGE_DRIVER.
//
// Outside production, minio with no MINIO_ACCESS_KEY becomes local, so a new
// project stores uploads before Docker is running instead of answering each
// one with STORAGE_UNAVAILABLE. Production never falls back: a server that
// lost its credentials should say so, not start writing to its own disk.
func resolveStorageDriver() string {
	driver := getEnv("STORAGE_DRIVER", "minio")
	if driver == "minio" && getEnv("MINIO_ACCESS_KEY", "") == "" && getEnv("APP_ENV", "production") != "production" {
		log.Println("MinIO has no credentials (MINIO_ACCESS_KEY), so files are kept on the local disk. Set STORAGE_DRIVER=local to make that the choice, or set the MinIO credentials to use MinIO")
		return "local"
	}
	return driver
}

`
