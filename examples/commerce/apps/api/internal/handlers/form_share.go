package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

type FormShareHandler struct {
	DB *gorm.DB
}

// maxSubmissionUserAgent is the width of form_submissions.user_agent. A longer
// header is cut to fit rather than failing the insert.
const maxSubmissionUserAgent = 500

// ── Admin endpoints ────────────────────────────────────────────────────

// Resources returns the list of resource names the form-share
// dispatcher knows how to render -- the dropdown the admin's New
// Share modal uses (v3.31.50).
func (h *FormShareHandler) Resources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": services.RegisteredResources(),
	})
}

// shares is the service every method below reads and writes through. What a
// share is, and how a submission is counted, is a decision about data: see
// internal/services/form_share.go. Who may see one, and whether a password was
// right, is decided here.
func (h *FormShareHandler) shares() *services.FormShareService {
	return &services.FormShareService{DB: h.DB}
}

// FieldsPreview returns the same field schema PublicGet returns,
// keyed by resource name rather than by share token. Used by the
// admin's New Share modal to render a preview before the share
// even exists. v3.31.50.
func (h *FormShareHandler) FieldsPreview(c *gin.Context) {
	resourceName := c.Param("resource")
	fields := services.PublicFields(resourceName)
	if fields == nil {
		respond.Fail(c, respond.CodeNotFound, "Resource not registered")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"resource_name": resourceName,
			"fields":        fields,
		},
	})
}

// List paginates form shares for the admin dashboard.
func (h *FormShareHandler) List(c *gin.Context) {
	shares, err := h.shares().All(c.Request.Context(), c.Query("resource_name"))
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": shares,
		"meta": gin.H{"total": len(shares), "page": 1, "page_size": len(shares), "pages": 1},
	})
}

// A new public form share.
type CreateFormShareRequest struct {
	ResourceName      string   `json:"resource_name" binding:"required"`
	Label             string   `json:"label"`
	Password          string   `json:"password"`
	CustomTitle       string   `json:"custom_title"`
	CustomDescription string   `json:"custom_description"`
	HiddenFields      []string `json:"hidden_fields"`
}

// Create generates a new share for a resource. Optional password is
// bcrypt-hashed before storage. A 32-char URL-safe token is generated
// automatically.

func (h *FormShareHandler) Create(c *gin.Context) {
	var req CreateFormShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	token, err := randomToken(24)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "token generation failed")
		return
	}

	if req.HiddenFields == nil {
		req.HiddenFields = []string{}
	}
	share := models.FormShare{
		ResourceName:      req.ResourceName,
		Token:             token,
		Label:             req.Label,
		Enabled:           true,
		CustomTitle:       req.CustomTitle,
		CustomDescription: req.CustomDescription,
		HiddenFields:      datatypes.NewJSONSlice(req.HiddenFields),
	}
	if userID, ok := c.Get("user_id"); ok {
		if s, ok := userID.(string); ok {
			share.CreatedByUserID = s
		}
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			respond.Fail(c, respond.CodeInternalError, "password hash failed")
			return
		}
		share.PasswordHash = string(hash)
	}

	if err := h.shares().Create(c.Request.Context(), &share); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	share.HasPassword = share.PasswordHash != ""
	c.JSON(http.StatusCreated, gin.H{"data": share, "message": "Share created"})
}

// Changes to a share.
type UpdateFormShareRequest struct {
	Label             *string   `json:"label"`
	Enabled           *bool     `json:"enabled"`
	Password          *string   `json:"password"`
	CustomTitle       *string   `json:"custom_title"`
	CustomDescription *string   `json:"custom_description"`
	HiddenFields      *[]string `json:"hidden_fields"`
}

// Update toggles enabled/label/password. Pass password="" to leave
// unchanged; pass password="-" to remove an existing password.

func (h *FormShareHandler) Update(c *gin.Context) {
	share, err := h.shares().ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "Share not found")
		return
	}

	var req UpdateFormShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
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
				respond.Fail(c, respond.CodeInternalError, "password hash failed")
				return
			}
			share.PasswordHash = string(hash)
		}
	}
	// v3.31.50 — operator-customisable surface.
	if req.CustomTitle != nil {
		share.CustomTitle = *req.CustomTitle
	}
	if req.CustomDescription != nil {
		share.CustomDescription = *req.CustomDescription
	}
	if req.HiddenFields != nil {
		share.HiddenFields = datatypes.NewJSONSlice(*req.HiddenFields)
	}
	if err := h.shares().Save(c.Request.Context(), share); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	share.HasPassword = share.PasswordHash != ""
	c.JSON(http.StatusOK, gin.H{"data": share, "message": "Share updated"})
}

// Delete soft-deletes a share (token stops working).
func (h *FormShareHandler) Delete(c *gin.Context) {
	if err := h.shares().Delete(c.Request.Context(), c.Param("id")); err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
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
	share, err := h.shares().EnabledByToken(c.Request.Context(), token)
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "Link not found or disabled")
		return
	}
	// v3.31.50 — filter hidden fields out of the schema before
	// returning it. The operator picked these in the New Share
	// modal; they should never show up to anonymous visitors.
	all := services.PublicFields(share.ResourceName)
	hidden := map[string]bool{}
	for _, k := range share.HiddenFields {
		hidden[k] = true
	}
	visible := make([]services.PublicFieldInfo, 0, len(all))
	for _, f := range all {
		if !hidden[f.Key] {
			visible = append(visible, f)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"resource_name":      share.ResourceName,
			"has_password":       share.PasswordHash != "",
			"label":              share.Label,
			"custom_title":       share.CustomTitle,
			"custom_description": share.CustomDescription,
			"fields":             visible,
		},
	})
}

// A submission through a public share.
type PublicFormSubmitRequest struct {
	Password string                 `json:"_password"`
	Fields   map[string]interface{} `json:"fields" binding:"required"`
}

// PublicSubmit accepts the form payload, verifies the password (when
// required), and dispatches to the resource's create service. Returns
// the new record's ID + an opaque label.

func (h *FormShareHandler) PublicSubmit(c *gin.Context) {
	token := c.Param("token")
	share, err := h.shares().EnabledByToken(c.Request.Context(), token)
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "Link not found or disabled")
		return
	}

	var body PublicFormSubmitRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	if share.PasswordHash != "" {
		if body.Password == "" || bcrypt.CompareHashAndPassword([]byte(share.PasswordHash), []byte(body.Password)) != nil {
			respond.Fail(c, respond.CodePasswordRequired, "Password incorrect or missing")
			return
		}
	}

	out, err := services.SubmitSharedForm(h.DB, share.ResourceName, body.Fields)
	if err != nil {
		respond.ServerError(c, "SUBMISSION_FAILED", err, "Your submission could not be saved")
		return
	}

	// Best-effort: a count that fails does not retroactively invalidate somebody's
	// submission. The increment itself is the service's, and is in SQL, because
	// this one used to lose a submission whenever two arrived together.
	if err := h.shares().CountSubmission(c.Request.Context(), share); err != nil {
		log.Printf("form share %s: counting a submission: %v", share.ID, err)
	}

	// v3.31.25 — write the audit row. Best-effort; failure here means
	// the visitor still gets their record, the admin just misses one
	// line in the trail.
	ua := c.GetHeader("User-Agent")
	if len(ua) > maxSubmissionUserAgent {
		ua = ua[:maxSubmissionUserAgent]
	}
	if err := h.shares().RecordSubmission(c.Request.Context(), &models.FormSubmission{
		ShareID:      share.ID,
		ResourceName: share.ResourceName,
		RecordID:     out.ID,
		IP:           c.ClientIP(),
		UserAgent:    ua,
	}); err != nil {
		log.Printf("form share %s: writing the audit row for record %s: %v", share.ID, out.ID, err)
	}

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
	rows, err := h.shares().Submissions(c.Request.Context(), c.Query("share_id"), c.Query("resource_name"), 100)
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rows,
		"meta": gin.H{"total": len(rows), "page": 1, "page_size": len(rows), "pages": 1},
	})
}

// randomToken returns a URL-safe base64 string of about `byteLen*4/3`
// characters. 24 bytes → ~32 chars, plenty of entropy against brute
// force at the public endpoint (paired with Sentinel rate limits).
func randomToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}
