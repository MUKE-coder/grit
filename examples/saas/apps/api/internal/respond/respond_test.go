package respond

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type answer struct {
	Error Error `json:"error"`
}

// writeErrorOn answers err through WriteError on a route shaped like a
// generated resource's, and returns the status and the envelope.
func writeErrorOn(t *testing.T, route, path string, err error) (int, Error) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST(route, func(c *gin.Context) { WriteError(c, err, "Failed to save contact") })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, nil))
	var body answer
	if decodeErr := json.Unmarshal(w.Body.Bytes(), &body); decodeErr != nil {
		t.Fatalf("decoding %q: %v", w.Body.String(), decodeErr)
	}
	return w.Code, body.Error
}

type dupContact struct {
	ID   uint
	Code string `gorm:"uniqueIndex"`
}

// The error a real driver returns, not a string written to look like one.
func TestWriteErrorAnswersADuplicateWith409(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&dupContact{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&dupContact{Code: "A-1"}).Error; err != nil {
		t.Fatal(err)
	}
	dup := db.Create(&dupContact{Code: "A-1"}).Error
	if dup == nil {
		t.Fatal("the second insert should have been refused")
	}

	status, body := writeErrorOn(t, "/api/dup-contacts", "/api/dup-contacts", dup)
	if status != http.StatusConflict || body.Code != string(CodeConflict) {
		t.Fatalf("got %d %s, want 409 CONFLICT", status, body.Code)
	}
	if body.Details["code"] == "" {
		t.Errorf("the field was not named: %+v", body)
	}
}

// Postgres, MySQL and SQL Server name the index, not the column. GORM names
// it idx_<table>_<column> or uni_<table>_<column>, and the route names the
// table.
func TestWriteErrorNamesTheFieldFromTheIndex(t *testing.T) {
	for name, msg := range map[string]string{
		"postgres":        "ERROR: duplicate key value violates unique constraint \"idx_contacts_code\" (SQLSTATE 23505)",
		"postgres uni":    "ERROR: duplicate key value violates unique constraint \"uni_contacts_code\" (SQLSTATE 23505)",
		"mysql":           "Error 1062 (23000): Duplicate entry 'A-1' for key 'contacts.idx_contacts_code'",
		"mysql, no table": "Error 1062 (23000): Duplicate entry 'A-1' for key 'idx_contacts_code'",
	} {
		status, body := writeErrorOn(t, "/api/contacts/:id", "/api/contacts/7", errors.New(msg))
		if status != http.StatusConflict {
			t.Errorf("%s: got %d, want 409", name, status)
			continue
		}
		if body.Details["code"] == "" {
			t.Errorf("%s: the field was not named: %+v", name, body)
		}
	}
}

// An index the route does not explain is still a 409, without a guess at
// which field.
func TestWriteErrorDoesNotGuessAField(t *testing.T) {
	msg := "ERROR: duplicate key value violates unique constraint \"contacts_pkey\" (SQLSTATE 23505)"
	status, body := writeErrorOn(t, "/api/contacts", "/api/contacts", errors.New(msg))
	if status != http.StatusConflict {
		t.Fatalf("got %d, want 409", status)
	}
	if len(body.Details) != 0 {
		t.Errorf("guessed a field: %+v", body.Details)
	}
}

func TestWriteErrorKnowsGormsSentinel(t *testing.T) {
	status, _ := writeErrorOn(t, "/api/contacts", "/api/contacts", gorm.ErrDuplicatedKey)
	if status != http.StatusConflict {
		t.Errorf("gorm.ErrDuplicatedKey: got %d, want 409", status)
	}
}

// Everything else stays an opaque 500: the driver's text is for the log.
func TestWriteErrorKeepsOtherErrorsOpaque(t *testing.T) {
	status, body := writeErrorOn(t, "/api/contacts", "/api/contacts", errors.New("connection refused"))
	if status != http.StatusInternalServerError || body.Message != "Failed to save contact" {
		t.Errorf("got %d %q, want 500 with the fallback message", status, body.Message)
	}
}
