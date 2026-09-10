package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A project that moved REDIS_PORT but kept REDIS_URL on the old port was using
// another project's Redis without a single error. Found on a ledger whose jobs
// and cache were landing in a storefront's Redis.

func dotEnvProject(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestEnvPortDriftNamesTheMismatch(t *testing.T) {
	root := dotEnvProject(t, "REDIS_PORT=6780\nREDIS_URL=redis://localhost:6380\n"+
		"MINIO_PORT=9402              # moved\nMINIO_ENDPOINT=\"http://localhost:9002\"\n")
	got := envPortDrift(root)
	if len(got) != 2 {
		t.Fatalf("want both services reported, got %d:\n%s", len(got), strings.Join(got, "\n"))
	}
	if !strings.Contains(got[0], "REDIS_URL=redis://localhost:6780") {
		t.Errorf("the Redis warning does not give the line to set: %s", got[0])
	}
	if !strings.Contains(got[1], "MINIO_ENDPOINT=http://localhost:9402") {
		t.Errorf("the MinIO warning does not give the line to set: %s", got[1])
	}
}

func TestEnvPortDriftStaysQuietWhenThingsAgree(t *testing.T) {
	for name, body := range map[string]string{
		"matching ports":   "REDIS_PORT=6780\nREDIS_URL=redis://localhost:6780\n",
		"remote host":      "REDIS_PORT=6780\nREDIS_URL=redis://cache.internal:6380\n",
		"url commented":    "REDIS_PORT=6780\n# REDIS_URL=redis://localhost:6380\n",
		"redis turned off": "REDIS_PORT=6780\nREDIS_URL=\n",
		"no port set":      "REDIS_URL=redis://localhost:6380\n",
		"no .env at all":   "",
	} {
		root := dotEnvProject(t, body)
		if body == "" {
			_ = os.Remove(filepath.Join(root, ".env"))
		}
		if got := envPortDrift(root); len(got) != 0 {
			t.Errorf("%s: false alarm: %v", name, got)
		}
	}
}

// The config a new project ships with builds both addresses from the ports,
// and no longer pins either.
func TestConfigFollowsThePorts(t *testing.T) {
	src := apiConfigGo()
	for _, want := range []string{
		`getEnv("REDIS_PORT", "6380")`,
		`getEnv("MINIO_PORT", "9002")`,
		"Endpoint:  resolveMinioEndpoint(),",
		`warnPortMismatch("REDIS_URL", v, "REDIS_PORT")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("config template is missing %s", want)
		}
	}
	if strings.Contains(src, `return "redis://localhost:6380"`) {
		t.Error("the Redis address is still a fixed default")
	}
	if strings.Contains(src, `getEnv("MINIO_ENDPOINT", "http://localhost:9002")`) {
		t.Error("the MinIO endpoint is still a fixed default")
	}
}
