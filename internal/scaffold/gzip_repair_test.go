package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

const securityHeadersAnchor = "// SecurityHeaders adds production security headers to every response."

func oldLoggerMiddleware(t *testing.T) string {
	t.Helper()
	fresh := apiLoggerMiddlewareGo()
	old := strings.Replace(fresh, securityHeadersAnchor, oldGzipMiddleware+securityHeadersAnchor, 1)
	old = strings.Replace(old, "import (\n\t\"fmt\"\n", "import (\n\t\"compress/gzip\"\n\t\"fmt\"\n", 1)
	if old == fresh || !strings.Contains(old, "\"compress/gzip\"") || !strings.Contains(old, "func Gzip()") {
		t.Fatal("could not rebuild logger.go from before the repair")
	}
	return old
}

func TestRepairGzipLogger(t *testing.T) {
	out, fixed, warn := repairGzipLoggerSource(oldLoggerMiddleware(t))
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != apiLoggerMiddlewareGo() {
		t.Error("the repaired logger.go differs from a fresh one")
	}
	if again, fixed, _ := repairGzipLoggerSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed logger.go again")
	}
}

func TestRepairGzipLoggerWarnsOnAnEditedGzip(t *testing.T) {
	src := strings.Replace(oldLoggerMiddleware(t), "gzip.BestSpeed", "gzip.BestCompression", 1)
	if out, fixed, warn := repairGzipLoggerSource(src); out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshTemplatesNeedNoGzipRepair(t *testing.T) {
	logger := apiLoggerMiddlewareGo()
	if strings.Contains(logger, "func Gzip(") || strings.Contains(logger, "compress/gzip") {
		t.Error("logger.go still carries a gzip middleware")
	}
	if out, fixed, warn := repairGzipLoggerSource(logger); out != logger || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("logger.go still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
	for name, src := range map[string]string{"gzip.go": middlewareGzipGo(), "gzip_test.go": middlewareGzipTestGo(), "logger.go": logger} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
