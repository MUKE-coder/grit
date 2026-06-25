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

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "form_share.go"):              formShareModelGo(),
		filepath.Join(apiRoot, "internal", "models", "form_submission.go"):         formSubmissionModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "form_share.go"):            formShareHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "form_share_dispatch.go"):   formShareDispatchGo(),
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

	"github.com/google/uuid"
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
	Enabled         bool           ` + "`" + `gorm:"not null;default:true" json:"enabled"` + "`" + `
	SubmissionCount int            ` + "`" + `gorm:"not null;default:0" json:"submission_count"` + "`" + `
	CreatedByUserID string         ` + "`" + `gorm:"size:36;index" json:"created_by_user_id"` + "`" + `
	Label           string         ` + "`" + `gorm:"size:200" json:"label"` + "`" + ` // optional operator-facing label
	CreatedAt       time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt       time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt       gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (s *FormShare) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
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

	"github.com/google/uuid"
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
		s.ID = uuid.New().String()
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
func formShareHandlerGo() string {
	return `package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

type FormShareHandler struct {
	DB *gorm.DB
}

// ── Admin endpoints ────────────────────────────────────────────────────

// List paginates form shares for the admin dashboard.
func (h *FormShareHandler) List(c *gin.Context) {
	var shares []models.FormShare
	q := h.DB.Order("created_at DESC")
	if rn := c.Query("resource_name"); rn != "" {
		q = q.Where("resource_name = ?", rn)
	}
	if err := q.Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": shares,
		"meta": gin.H{"total": len(shares), "page": 1, "page_size": len(shares), "pages": 1},
	})
}

// Create generates a new share for a resource. Optional password is
// bcrypt-hashed before storage. A 32-char URL-safe token is generated
// automatically.
func (h *FormShareHandler) Create(c *gin.Context) {
	var req struct {
		ResourceName string ` + "`" + `json:"resource_name" binding:"required"` + "`" + `
		Label        string ` + "`" + `json:"label"` + "`" + `
		Password     string ` + "`" + `json:"password"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	token, err := randomToken(24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "token generation failed"},
		})
		return
	}

	share := models.FormShare{
		ResourceName: req.ResourceName,
		Token:        token,
		Label:        req.Label,
		Enabled:      true,
	}
	if userID, ok := c.Get("user_id"); ok {
		if s, ok := userID.(string); ok {
			share.CreatedByUserID = s
		}
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "INTERNAL_ERROR", "message": "password hash failed"},
			})
			return
		}
		share.PasswordHash = string(hash)
	}

	if err := h.DB.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	share.HasPassword = share.PasswordHash != ""
	c.JSON(http.StatusCreated, gin.H{"data": share, "message": "Share created"})
}

// Update toggles enabled/label/password. Pass password="" to leave
// unchanged; pass password="-" to remove an existing password.
func (h *FormShareHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var share models.FormShare
	if err := h.DB.First(&share, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Share not found"},
		})
		return
	}

	var req struct {
		Label    *string ` + "`" + `json:"label"` + "`" + `
		Enabled  *bool   ` + "`" + `json:"enabled"` + "`" + `
		Password *string ` + "`" + `json:"password"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}
	if req.Label != nil {
		share.Label = *req.Label
	}
	if req.Enabled != nil {
		share.Enabled = *req.Enabled
	}
	if req.Password != nil {
		if *req.Password == "-" {
			share.PasswordHash = ""
		} else if *req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "INTERNAL_ERROR", "message": "password hash failed"},
				})
				return
			}
			share.PasswordHash = string(hash)
		}
	}
	if err := h.DB.Save(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	share.HasPassword = share.PasswordHash != ""
	c.JSON(http.StatusOK, gin.H{"data": share, "message": "Share updated"})
}

// Delete soft-deletes a share (token stops working).
func (h *FormShareHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.FormShare{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Share deleted"})
}

// ── Public endpoints ───────────────────────────────────────────────────

// PublicGet returns the resource name + has_password + the field
// schema so the public web page can render the correct inputs.
// Does NOT expose the hash or submission stats.
//
// v3.31.43: now also returns the resource's field shape via
// services.PublicFields. Before this release the public page
// rendered a hardcoded name/email/phone/message contact-form --
// which had nothing to do with the resource being submitted. Now
// a Category share shows Name + Image, a Product share shows
// Name + Price + Description + Category, etc.
func (h *FormShareHandler) PublicGet(c *gin.Context) {
	token := c.Param("token")
	var share models.FormShare
	if err := h.DB.First(&share, "token = ? AND enabled = ?", token, true).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Link not found or disabled"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"resource_name": share.ResourceName,
			"has_password":  share.PasswordHash != "",
			"label":         share.Label,
			"fields":        services.PublicFields(share.ResourceName),
		},
	})
}

// PublicSubmit accepts the form payload, verifies the password (when
// required), and dispatches to the resource's create service. Returns
// the new record's ID + an opaque label.
func (h *FormShareHandler) PublicSubmit(c *gin.Context) {
	token := c.Param("token")
	var share models.FormShare
	if err := h.DB.First(&share, "token = ? AND enabled = ?", token, true).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Link not found or disabled"},
		})
		return
	}

	var body struct {
		Password string                 ` + "`" + `json:"_password"` + "`" + `
		Fields   map[string]interface{} ` + "`" + `json:"fields" binding:"required"` + "`" + `
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if share.PasswordHash != "" {
		if body.Password == "" || bcrypt.CompareHashAndPassword([]byte(share.PasswordHash), []byte(body.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "PASSWORD_REQUIRED", "message": "Password incorrect or missing"},
			})
			return
		}
	}

	out, err := services.SubmitSharedForm(h.DB, share.ResourceName, body.Fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "SUBMISSION_FAILED", "message": err.Error()},
		})
		return
	}

	// Bump submission count (best-effort — failure here doesn't
	// retroactively invalidate the user's submission).
	h.DB.Model(&share).UpdateColumns(map[string]interface{}{
		"submission_count": share.SubmissionCount + 1,
		"updated_at":       time.Now(),
	})

	// v3.31.25 — write the audit row. Best-effort; failure here means
	// the visitor still gets their record, the admin just misses one
	// line in the trail. We truncate UA at 500 chars (column width).
	ua := c.GetHeader("User-Agent")
	if len(ua) > 500 {
		ua = ua[:500]
	}
	_ = h.DB.Create(&models.FormSubmission{
		ShareID:      share.ID,
		ResourceName: share.ResourceName,
		RecordID:     out.ID,
		IP:           c.ClientIP(),
		UserAgent:    ua,
	}).Error

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":    out.ID,
			"label": out.Label,
		},
		"message": "Submitted",
	})
}

// ListSubmissions returns audit-log rows for one or more shares.
// Filterable by share_id and resource_name; defaults to the 100 most
// recent rows. v3.31.25.
func (h *FormShareHandler) ListSubmissions(c *gin.Context) {
	var rows []models.FormSubmission
	q := h.DB.Order("created_at DESC").Limit(100)

	if shareID := c.Query("share_id"); shareID != "" {
		q = q.Where("share_id = ?", shareID)
	}
	if resourceName := c.Query("resource_name"); resourceName != "" {
		q = q.Where("resource_name = ?", resourceName)
	}

	if err := q.Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rows,
		"meta": gin.H{"total": len(rows), "page": 1, "page_size": len(rows), "pages": 1},
	})
}

// randomToken returns a URL-safe base64 string of about ` + "`" + `byteLen*4/3` + "`" + `
// characters. 24 bytes → ~32 chars, plenty of entropy against brute
// force at the public endpoint (paired with Sentinel rate limits).
func randomToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}
`
}
