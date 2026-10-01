package scaffold

// This file scaffolds the FormShare feature (Phase 2 of
// PLAN_FORMS_AND_SHARING.md): public form links with optional bcrypt
// password protection.
//
// What gets generated:
//   - models/form_share.go            — GORM model + helpers
//   - handlers/form_share.go          — admin CRUD + public surface
//   - services/form_share_dispatch.go — marker-driven resource dispatch
//                                       (each grit generate adds a case here)
//
// Public endpoints live under /api/public/forms/:token — no auth, no
// CSRF. The dispatch service is the security boundary: it whitelists
// which resources are reachable via a share token and which fields they
// accept.

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeFormShareFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	// The handler, its service and the models travel on upgrade: see
	// writeFrameworkOwnedFiles. This writes the rest at scaffold time.
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "services", "form_share_dispatch.go"): formShareDispatchGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// formShareModelGo emits the FormShare GORM model.
func formShareModelGo() string {
	return `package models

import (
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FormShare is a public link to a resource's create form. Operators
// generate one of these in the admin to expose a single Grit resource
// (e.g. Contact) without requiring auth — useful for lead forms,
// applications, public submissions.
//
// Optional bcrypt password adds a gate: if PasswordHash is set, the
// public submit endpoint requires the visitor to verify the password
// first.
type FormShare struct {
	ID              string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	ResourceName    string         ` + "`" + `gorm:"size:64;not null;index" json:"resource_name"` + "`" + `
	// Token is the URL-safe identifier (32 chars). Used as
	// /forms/<token> in the public web app and
	// /api/public/forms/:token/* in the API.
	Token           string         ` + "`" + `gorm:"size:64;not null;uniqueIndex" json:"token"` + "`" + `
	// PasswordHash is empty for open-access shares. When set (bcrypt
	// cost 10), visitors must POST the plaintext password to
	// /check-password before /submit succeeds.
	PasswordHash    string         ` + "`" + `gorm:"size:255" json:"-"` + "`" + `
	HasPassword     bool           ` + "`" + `gorm:"-" json:"has_password"` + "`" + ` // computed, not stored
	// not null, but no default: a default made Enabled:false unstorable on
	// create, so a share link meant to start disabled went live instead.
	Enabled         bool           ` + "`" + `gorm:"not null" json:"enabled"` + "`" + `
	SubmissionCount int            ` + "`" + `gorm:"not null;default:0" json:"submission_count"` + "`" + `
	CreatedByUserID string         ` + "`" + `gorm:"size:36;index" json:"created_by_user_id"` + "`" + `
	Label           string         ` + "`" + `gorm:"size:200" json:"label"` + "`" + ` // optional operator-facing label

	// v3.31.50 — operator-customisable surface for the public form.
	// CustomTitle / CustomDescription replace the default heading +
	// subtitle when set; HiddenFields is a list of field keys
	// (matching the json tag on the model) to omit from the
	// rendered form -- useful for optional columns the operator
	// doesn't want anonymous visitors filling in.
	CustomTitle       string                      ` + "`" + `gorm:"size:200" json:"custom_title"` + "`" + `
	CustomDescription string                      ` + "`" + `gorm:"size:500" json:"custom_description"` + "`" + `
	HiddenFields      datatypes.JSONSlice[string] ` + "`" + `gorm:"type:json" json:"hidden_fields"` + "`" + `

	CreatedAt       time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt       time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt       gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (s *FormShare) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}

// AfterFind computes the HasPassword virtual field so the admin UI can
// show a lock icon without exposing the hash itself.
func (s *FormShare) AfterFind(tx *gorm.DB) error {
	s.HasPassword = s.PasswordHash != ""
	return nil
}
`
}

// formSubmissionModelGo emits the FormSubmission audit-log model. One
// row per successful public submission — records which share, which
// resource, which record ID, plus IP + User-Agent for forensics.
//
// The audit row is best-effort: failure to write it does NOT roll
// back the user's submission. They get their record either way; the
// admin just loses one line in the audit trail.
func formSubmissionModelGo() string {
	return `package models

import (
	"time"

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)

// FormSubmission is an audit-log row for one successful public form
// submission. Created by FormShareHandler.PublicSubmit after the
// dispatcher returns a record ID. Never modified; soft-deletable so
// admins can prune old rows without losing the share's history.
type FormSubmission struct {
	ID           string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	ShareID      string         ` + "`" + `gorm:"size:36;not null;index" json:"share_id"` + "`" + `
	ResourceName string         ` + "`" + `gorm:"size:64;not null;index" json:"resource_name"` + "`" + `
	RecordID     string         ` + "`" + `gorm:"size:36;not null;index" json:"record_id"` + "`" + `
	// IP and UserAgent are best-effort — set from gin.Context at
	// submission time. Truncated to fit the column; long UAs are
	// trimmed to 500 chars.
	IP           string         ` + "`" + `gorm:"size:64" json:"ip"` + "`" + `
	UserAgent    string         ` + "`" + `gorm:"size:500" json:"user_agent"` + "`" + `
	CreatedAt    time.Time      ` + "`" + `json:"created_at"` + "`" + `
	DeletedAt    gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (s *FormSubmission) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}
`
}

// formShareDispatchGo emits the marker-driven dispatcher. Each call to
// grit generate resource appends a case here, so the public submit
// handler always knows how to materialise the latest resources.
func formShareDispatchGo() string {
	return `package services

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
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
	Key string ` + "`json:\"key\"`" + `
	// Label is a human-friendly rendering of Key for the form label.
	Label string ` + "`json:\"label\"`" + `
	// Type is the input shape the frontend should render. One of:
	//   "text" | "email" | "tel" | "textarea" | "number" |
	//   "checkbox" | "date" | "datetime" | "file"
	Type string ` + "`json:\"type\"`" + `
	// Required mirrors the binding:"required" tag.
	Required bool ` + "`json:\"required\"`" + `
}

// RegisteredResources returns every resource name that
// PublicFields knows how to render. v3.31.50 adds this so the
// admin New-Share modal can show a dropdown instead of a free-text
// input (typing "Catgeory" instead of "Category" was producing a
// silently-broken share). The generator injects one entry per
// ` + "`grit generate resource`" + ` at the marker below.
func RegisteredResources() []string {
	return []string{
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
`
}

// formShareHandlerGo emits the FormShare admin CRUD + public surface.
func formShareHandlerGo() string { return tmpl("api/handlers/form_share.go") }
