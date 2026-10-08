package respond

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

// The struct a grid row is decoded into, with the tags the generator writes.
type gadget struct {
	Name       string   `json:"name" binding:"required"`
	CategoryID string   `json:"category_id" binding:"required"`
	Email      string   `json:"email" binding:"omitempty,email"`
	Stock      int      `json:"stock" binding:"gte=0"`
	Code       string   `json:"code" binding:"omitempty,min=3"`
	Status     string   `json:"status" binding:"omitempty,oneof=draft live"`
	Tags       []string `json:"tags" binding:"omitempty,min=2"`
	Internal   string   `json:"-" binding:"omitempty,email"`
}

func TestValidateStructCatchesWhatBindingWouldHave(t *testing.T) {
	// The bug this exists for: a bulk-create row with no category, which
	// encoding/json accepted happily and the edit form then refused to save.
	err := ValidateStruct(&gadget{Name: "Whisk"})
	if err == nil {
		t.Fatal("a row with no category was accepted")
	}
	if !strings.Contains(err.Error(), "category_id") {
		t.Fatalf("the message does not name the field the client sent: %q", err.Error())
	}
	if strings.Contains(err.Error(), "CategoryID") {
		t.Fatalf("the message names a Go field: %q", err.Error())
	}
}

func TestValidateStructPassesAGoodRow(t *testing.T) {
	if err := ValidateStruct(&gadget{Name: "Whisk", CategoryID: "abc"}); err != nil {
		t.Fatalf("a complete row was rejected: %v", err)
	}
}

// The point of returning an error that implements FieldErrors: a handler that
// already calls WriteError answers 422 with a message per input, with nothing
// added to the handler.
func TestWriteErrorAnswersAValidationFailurePerField(t *testing.T) {
	status, body := writeErrorOn(t, "/gadgets", "/gadgets", ValidateStruct(&gadget{Name: "Whisk"}))
	if status != http.StatusUnprocessableEntity {
		t.Fatalf("status %d, want 422", status)
	}
	if body.Code != string(CodeValidationError) {
		t.Errorf("code %q", body.Code)
	}
	if body.Details["category_id"] != "This field is required" {
		t.Errorf("details %v", body.Details)
	}
}

func TestValidationFieldsKeysByJSONName(t *testing.T) {
	err := validationErrorOf(t, &gadget{Stock: -1, Code: "ab", Status: "nope", Email: "x"})
	fields := ValidationFields(err)

	for _, want := range []string{"name", "category_id", "stock", "code", "status", "email"} {
		if _, ok := fields[want]; !ok {
			t.Errorf("no entry for %q, got %v", want, keysOf(fields))
		}
	}
	if fields["name"] != "This field is required" {
		t.Errorf("required reads %q", fields["name"])
	}
	if fields["email"] != "Enter a valid email address" {
		t.Errorf("email reads %q", fields["email"])
	}
	if fields["stock"] != "Must be 0 or more" {
		t.Errorf("a number's gte reads %q", fields["stock"])
	}
	if fields["code"] != "Must be at least 3 characters" {
		t.Errorf("a string's min reads %q", fields["code"])
	}
	if fields["status"] != "Must be one of: draft, live" {
		t.Errorf("oneof reads %q", fields["status"])
	}
}

func TestValidationFieldsCountsAListInItems(t *testing.T) {
	err := validationErrorOf(t, &gadget{Name: "a", CategoryID: "b", Tags: []string{"one"}})
	if got := ValidationFields(err)["tags"]; got != "Must have at least 2 items" {
		t.Errorf("a slice's min reads %q", got)
	}
}

// Calling it on the error ValidateStruct returned has to give the same map,
// since that error carries one and is no longer the validator's.
func TestValidationFieldsReadsItsOwnError(t *testing.T) {
	fields := ValidationFields(ValidateStruct(&gadget{Name: "Whisk"}))
	if fields["category_id"] != "This field is required" {
		t.Errorf("got %v", fields)
	}
}

// A field kept off the wire has no name to report, and reporting it as the
// empty string would produce an envelope with an unkeyed entry in it.
func TestValidationFieldsNamesAFieldThatIsNotOnTheWire(t *testing.T) {
	err := validationErrorOf(t, &gadget{Name: "a", CategoryID: "b", Internal: "nope"})
	fields := ValidationFields(err)
	if len(fields) != 1 {
		t.Fatalf("expected the one failure, got %v", fields)
	}
	for name, sentence := range fields {
		if name == "" {
			t.Error("the entry has no key")
		}
		if sentence != "Enter a valid email address" {
			t.Errorf("it reads %q", sentence)
		}
	}
}

func TestValidationFieldsIgnoresEverythingElse(t *testing.T) {
	if fields := ValidationFields(errors.New("the database is on fire")); fields != nil {
		t.Errorf("a plain error produced field errors: %v", fields)
	}
	if got := ValidationMessage(errors.New("the database is on fire")); got != "the database is on fire" {
		t.Errorf("a plain error was rewritten to %q", got)
	}
}

func TestValidationMessageIsStableAcrossRuns(t *testing.T) {
	// A map's iteration order is random, so this read differently every time
	// it was logged until the names were sorted.
	first := ValidationMessage(validationErrorOf(t, &gadget{}))
	for i := 0; i < 20; i++ {
		if again := ValidationMessage(validationErrorOf(t, &gadget{})); again != first {
			t.Fatalf("the message changed between runs:\n%q\n%q", first, again)
		}
	}
	if !strings.HasPrefix(first, "category_id") {
		t.Errorf("the fields are not in order: %q", first)
	}
}

// validationErrorOf returns the error gin's validator produces, which is what
// a handler receiving one from ShouldBindJSON has in hand.
func validationErrorOf(t *testing.T, obj interface{}) error {
	t.Helper()
	err := binding.Validator.ValidateStruct(obj)
	if err == nil {
		t.Fatalf("%#v was expected to fail validation", obj)
	}
	return err
}

func keysOf(fields map[string]string) []string {
	out := make([]string, 0, len(fields))
	for name := range fields {
		out = append(out, name)
	}
	return out
}
