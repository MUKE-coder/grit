package respond

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm/schema"
)

// One validator, for every way into a table.
//
// A handler that binds a request body gets the binding: tags enforced for it:
// gin runs them inside ShouldBindJSON, so a create missing a required field is
// rejected before any service sees it. Code that decodes a value itself does
// not. Bulk create decodes a grid of rows with encoding/json, which knows
// nothing about binding tags, so required, min and oneof simply did not fire:
// the grid inserted a gadget with no category into a table whose model
// requires one, and the edit form then refused to save that row, because the
// form does enforce it. A row you can create and cannot edit is the worst of
// both.
//
// ValidateStruct is that second path: the same rules, run by hand, phrased the
// same way, and returning an error WriteError already knows how to answer.
// Any new way into a table wants it too, and a second set of rules is not an
// option worth having.

// columnNaming is GORM's, so a Go field becomes the name the client sent.
var columnNaming = schema.NamingStrategy{}

// Name the fields the way the wire does.
//
// The validator reports the Go field: "Field validation for 'CategoryID'".
// That is a sentence about a struct, written for whoever wrote the struct, and
// it reaches the client, which sent category_id and has nothing called
// CategoryID. Teaching the validator the json tag once fixes the message for
// every bound request in the app, not only for the paths below.
func init() {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return // somebody replaced gin's validator; jsonName still copes
	}
	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return "" // not on the wire, so it has no name to report
		}
		return name
	})
}

// validationError is a binding failure in the shape the rest of this package
// already handles: it implements FieldErrors, so WriteError answers it with 422
// and a message per field, and a form can put each one under its input. That
// is the whole reason to return it rather than the validator's own error.
type validationError struct {
	message string
	fields  map[string]string
}

func (e *validationError) Error() string { return e.message }

func (e *validationError) FieldErrors() map[string]string { return e.fields }

// ValidateStruct runs the binding: tags gin would have run, on a value that
// did not arrive as a request body.
//
// It returns nil, or an error carrying both a sentence and the field map.
func ValidateStruct(obj interface{}) error {
	err := binding.Validator.ValidateStruct(obj)
	if err == nil {
		return nil
	}
	return &validationError{message: ValidationMessage(err), fields: ValidationFields(err)}
}

// ValidationFields turns a binding failure into {field: sentence}, keyed by the
// name the client sent.
//
// It returns nil for anything that is not a validation failure, so a caller can
// tell "the client sent bad data" from "something else went wrong" by whether
// there is anything in the map.
func ValidationFields(err error) map[string]string {
	// An error already carrying a field map answers for itself, which is what
	// makes this safe to call on the error ValidateStruct returned.
	var carried FieldErrors
	if errors.As(err, &carried) {
		return carried.FieldErrors()
	}
	var invalid validator.ValidationErrors
	if !errors.As(err, &invalid) {
		return nil
	}
	fields := make(map[string]string, len(invalid))
	for _, fe := range invalid {
		fields[jsonName(fe.Field())] = sentenceFor(fe)
	}
	return fields
}

// ValidationMessage is the same failure as one line, for a log, a toast, or a
// row of a grid with no room for a field map.
func ValidationMessage(err error) string {
	fields := ValidationFields(err)
	if len(fields) == 0 {
		return err.Error()
	}
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	// A map's order is deliberately random, and a message that changes between
	// two identical requests cannot be tested, diffed or trusted in a log.
	sort.Strings(names)

	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, name+": "+fields[name])
	}
	return strings.Join(parts, "; ")
}

// jsonName is the field as the client knows it.
//
// With the hook in init() registered, the validator already reports the json
// name and this returns it unchanged. Without it the name is a Go field, and
// GORM's naming gives the same answer the generator wrote into the json tag.
func jsonName(field string) string {
	if field == "" {
		return "_"
	}
	if field == strings.ToLower(field) {
		return field
	}
	return columnNaming.ColumnName("", field)
}

// sentenceFor is what to put beside the input.
//
// Only the tags the generator emits and the ones a developer reaches for are
// phrased here. Anything else falls through to the tag's own name, which is
// terse but true, and better than a sentence that guesses at what the rule
// meant.
func sentenceFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required", "required_with", "required_without":
		return "This field is required"
	case "email":
		return "Enter a valid email address"
	case "url", "uri":
		return "Enter a valid URL"
	case "uuid", "uuid4", "uuid7":
		return "Must be a valid id"
	case "oneof":
		return "Must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "len":
		return countedAs(fe, "Must be exactly %s characters", "Must have exactly %s items", "Must be %s")
	case "min", "gte":
		return countedAs(fe, "Must be at least %s characters", "Must have at least %s items", "Must be %s or more")
	case "max", "lte":
		return countedAs(fe, "Must be at most %s characters", "Must have at most %s items", "Must be %s or less")
	case "gt":
		return countedAs(fe, "Must be longer than %s characters", "Must have more than %s items", "Must be more than %s")
	case "lt":
		return countedAs(fe, "Must be shorter than %s characters", "Must have fewer than %s items", "Must be less than %s")
	case "eqfield":
		return "Must match " + jsonName(fe.Param())
	case "nefield":
		return "Must be different from " + jsonName(fe.Param())
	case "numeric":
		return "Must be a number"
	case "alphanum":
		return "Letters and numbers only"
	case "datetime":
		return "Must be a date"
	}
	return "Failed the " + fe.Tag() + " rule"
}

// countedAs picks the unit, because "at least 3" is characters for a string,
// items for a list, and nothing at all for a number.
func countedAs(fe validator.FieldError, text, list, number string) string {
	switch fe.Kind() {
	case reflect.String:
		return fmt.Sprintf(text, fe.Param())
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf(list, fe.Param())
	}
	return fmt.Sprintf(number, fe.Param())
}
