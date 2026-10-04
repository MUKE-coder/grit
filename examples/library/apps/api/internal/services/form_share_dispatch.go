package services

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"encoding/json"

	"library/apps/api/internal/models"
)

// SharedResourceSubmission is the result of a public form submission —
// the created record's ID and human label, both safe to return to
// anonymous visitors.
type SharedResourceSubmission struct {
	ID    string
	Label string
}

// SubmitSharedForm dispatches a public form submission to the right
// resource service based on the FormShare's ResourceName. fields is
// a free-form map (validated by the resource service's own binding
// rules), since public submissions don't carry the operator's typed
// struct context.
//
// Adding a new resource? grit generate resource appends a case to
// the switch below at the auto-dispatch marker. Each case re-marshals
// fields into the typed model via json.Marshal(fields) — that's why
// the parameter is named "fields" rather than "body".
func SubmitSharedForm(db *gorm.DB, resourceName string, fields map[string]interface{}) (*SharedResourceSubmission, error) {
	switch resourceName {
	case "Author":
		item := &models.Author{}
		body, _ := json.Marshal(fields)
		if err := json.Unmarshal(body, item); err != nil {
			return nil, fmt.Errorf("decoding Author body: %w", err)
		}
		if err := db.Create(item).Error; err != nil {
			return nil, fmt.Errorf("creating Author: %w", err)
		}
		return &SharedResourceSubmission{ID: item.ID, Label: item.Name}, nil

	case "Book":
		item := &models.Book{}
		body, _ := json.Marshal(fields)
		if err := json.Unmarshal(body, item); err != nil {
			return nil, fmt.Errorf("decoding Book body: %w", err)
		}
		if err := db.Create(item).Error; err != nil {
			return nil, fmt.Errorf("creating Book: %w", err)
		}
		return &SharedResourceSubmission{ID: item.ID, Label: item.Title}, nil

	case "Loan":
		item := &models.Loan{}
		body, _ := json.Marshal(fields)
		if err := json.Unmarshal(body, item); err != nil {
			return nil, fmt.Errorf("decoding Loan body: %w", err)
		}
		if err := db.Create(item).Error; err != nil {
			return nil, fmt.Errorf("creating Loan: %w", err)
		}
		return &SharedResourceSubmission{ID: item.ID, Label: item.ID}, nil

	// grit:form-share:dispatch
	default:
		return nil, fmt.Errorf("public submission disabled for %q (no dispatch case registered)", resourceName)
	}
}

// PublicFieldInfo describes one form field the public page should
// render. Keep this struct small + JSON-friendly -- the web client
// reads it directly to build inputs.
type PublicFieldInfo struct {
	// Key matches the json tag on the Go model so the field name on
	// the wire matches what SubmitSharedForm's typed unmarshal
	// expects. e.g. "name", "category_id", "image".
	Key string `json:"key"`
	// Label is a human-friendly rendering of Key for the form label.
	Label string `json:"label"`
	// Type is the input shape the frontend should render. One of:
	//   "text" | "email" | "tel" | "textarea" | "number" |
	//   "checkbox" | "date" | "datetime" | "file"
	Type string `json:"type"`
	// Required mirrors the binding:"required" tag.
	Required bool `json:"required"`
}

// RegisteredResources returns every resource name that
// PublicFields knows how to render. v3.31.50 adds this so the
// admin New-Share modal can show a dropdown instead of a free-text
// input (typing "Catgeory" instead of "Category" was producing a
// silently-broken share). The generator injects one entry per
// `grit generate resource` at the marker below.
func RegisteredResources() []string {
	return []string{
		"Author",

		"Book",

		"Loan",

		// grit:form-share:registered
	}
}

// PublicFields returns the field schema for the public form to
// render. v3.31.43: replaces the previous hardcoded shape with the
// actual resource fields. The switch mirrors SubmitSharedForm so the
// generator only has to emit one extra case per new resource at the
// marker comment inside the switch.
func PublicFields(resourceName string) []PublicFieldInfo {
	switch resourceName {
	case "Author":
		return reflectPublicFields(&models.Author{})

	case "Book":
		return reflectPublicFields(&models.Book{})

	case "Loan":
		return reflectPublicFields(&models.Loan{})

	// grit:form-share:fields
	default:
		return nil
	}
}

// reflectPublicFields walks a model's struct fields and returns the
// public-form descriptors. Skips framework columns (id, created_at,
// etc), slug fields (auto-generated), and json:"-" fields. The
// generator-emitted cases in PublicFields call this with the model
// pointer.
func reflectPublicFields(model interface{}) []PublicFieldInfo {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	skip := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
		"version":    true,
		"slug":       true,
	}

	var out []PublicFieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		if skip[jsonTag] {
			continue
		}

		required := false
		for _, part := range strings.Split(f.Tag.Get("binding"), ",") {
			if strings.TrimSpace(part) == "required" {
				required = true
				break
			}
		}

		out = append(out, PublicFieldInfo{
			Key:      jsonTag,
			Label:    humanizePublicLabel(jsonTag),
			Type:     publicTypeFor(jsonTag, f.Type),
			Required: required,
		})
	}
	return out
}

// publicTypeFor maps a Go reflect.Type onto the public form's input
// type. FileRef columns resolve to "file" so the frontend renders
// the "not supported on public shares" state uniformly.
func publicTypeFor(fieldName string, t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typeName := t.String()
	if strings.Contains(typeName, "FileRef") || strings.Contains(typeName, "FileRefs") {
		return "file"
	}
	if strings.Contains(typeName, "time.Time") {
		return "datetime"
	}
	switch t.Kind() {
	case reflect.Bool:
		return "checkbox"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		lower := strings.ToLower(fieldName)
		switch {
		case lower == "email" || strings.HasSuffix(lower, "_email"):
			return "email"
		case lower == "phone" || strings.HasSuffix(lower, "_phone") || lower == "tel":
			return "tel"
		case lower == "description" || lower == "notes" || lower == "message" ||
			lower == "body" || lower == "content" || lower == "bio" ||
			lower == "summary" || strings.HasSuffix(lower, "_description"):
			return "textarea"
		}
		return "text"
	}
	return "text"
}

// humanizePublicLabel turns "category_id" into "Category Id" and
// "first_name" into "First Name".
func humanizePublicLabel(key string) string {
	parts := strings.Split(key, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}
