package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The templates are held to the rewrite the upgrade repair applies: no
// template may still send err.Error() with a 500, or with one of the 400s
// that passed a server error through.
func TestTemplatesSendNoServerErrorText(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") || name == "server_errors_repair.go" {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, n := hideServerErrors(string(src)); n != 0 {
			t.Errorf("%s still sends err.Error() with a server error in %d places", name, n)
		}
	}
}

func TestRespondLogsServerErrors(t *testing.T) {
	src := apiRespondGo()
	if !strings.Contains(src, "func ServerError(c *gin.Context, code string, err error, message string)") {
		t.Fatal("respond.ServerError is missing")
	}
	if strings.Contains(src, "_ = internalErr") {
		t.Error("respond.Internal still discards its error")
	}
	if !strings.Contains(src, `c.GetString("request_id")`) {
		t.Error("a logged 500 does not carry the request id")
	}
}

const serverErrorsBefore = `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/app/internal/models"
)

func a(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{"code": "DB_ERROR", "message": err.Error()},
	})
}

func b(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "EXPORT_FAILED", "message": err.Error()}})
}

func d(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
		"code": "INTERNAL_ERROR", "message": err.Error(),
	}})
}

func e(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "RETRY_FAILED",
			"message": "Failed to retry job: " + err.Error(),
		},
	})
}

func f(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": gin.H{"code": "STATS_FAILED", "message": err.Error()},
	})
}

func g(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_BODY", "message": err.Error()}})
	_ = models.User{}
}
`

func TestHideServerErrors(t *testing.T) {
	out, changes, warnings := repairServerErrorsSource(serverErrorsBefore, "example.com/app", "fixture.go")
	if len(warnings) != 0 || len(changes) != 1 {
		t.Fatalf("changes %v warnings %v", changes, warnings)
	}
	for _, want := range []string{
		`respond.ServerError(c, "DB_ERROR", err, "Internal server error")`,
		`respond.ServerError(c, "EXPORT_FAILED", err, "Internal server error")`,
		`respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")`,
		`respond.ServerError(c, "RETRY_FAILED", err, "Failed to retry job")`,
		`respond.ServerError(c, "STATS_FAILED", err, "Could not compute the stats for this resource")`,
		`"example.com/app/internal/respond"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s", want)
		}
	}
	// A 400 about the request keeps its message: it tells the caller what to fix.
	if !strings.Contains(out, `"code": "INVALID_BODY", "message": err.Error()`) {
		t.Error("a validation 400 was rewritten")
	}
	if strings.Count(out, "err.Error()") != 1 {
		t.Errorf("want only the validation message left, got:\n%s", out)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Errorf("the rewritten file is not valid Go: %v\n%s", err, out)
	}
	if again, _, _ := repairServerErrorsSource(out, "example.com/app", "fixture.go"); again != out {
		t.Error("the repair is not idempotent")
	}
}

func TestHealthStopsTellingCallersWhy(t *testing.T) {
	before := `package routes

import (
	"log"
)

func health() {
	if err := ping(); err != nil {
		dbStatus.OK = false
		dbStatus.Error = err.Error()
	}
	if err := redisPing(); err != nil {
		redisStatus.OK = false
		redisStatus.Error = err.Error()
	}
	log.Println("ok")
}
`
	out, changes, warnings := repairHealthErrorsSource(before)
	if len(warnings) != 0 || len(changes) != 1 {
		t.Fatalf("changes %v warnings %v", changes, warnings)
	}
	if strings.Contains(out, "err.Error()") {
		t.Errorf("the health endpoint still returns the error:\n%s", out)
	}
	for _, want := range []string{`log.Printf("health: database ping failed: %v", err)`, `log.Printf("health: redis ping failed: %v", err)`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Count(out, `"log"`) != 1 {
		t.Error("the log import was added twice")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Errorf("not valid Go: %v", err)
	}
	if strings.Contains(apiFilesHealthTemplate(), "Status.Error = err.Error()") {
		t.Error("the routes template still returns dependency errors from /api/health")
	}
}

// apiFilesHealthTemplate is the part of api_files.go that holds the health
// route, read from source so the test does not depend on the writer's options.
func apiFilesHealthTemplate() string {
	src, err := os.ReadFile("api_files.go")
	if err != nil {
		return ""
	}
	return string(src)
}
