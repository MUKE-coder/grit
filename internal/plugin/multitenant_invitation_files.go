package plugin

import "strings"

func mtInvitationModel(ctx Context) string {
	src := `package models

import (
	"time"
	"gorm.io/gorm"
	"{{MODULE}}/internal/ids"
)

// OrganizationInvitation stores only a one-way hash of the accept token.
// PendingKey is non-nil only while pending, allowing one pending invite per
// organization and email while permitting a later re-invitation.
type OrganizationInvitation struct {
	ID             string     ~gorm:"primarykey;size:36" json:"id"~
	OrganizationID string     ~gorm:"size:36;not null;index;uniqueIndex:ux_org_invitation_pending" json:"organization_id"~
	Email          string     ~gorm:"size:255;not null" json:"email"~
	Role           string     ~gorm:"size:24;not null;default:member" json:"role"~
	InvitedBy      string     ~gorm:"size:36;not null" json:"invited_by"~
	TokenHash      string     ~gorm:"size:64;not null;uniqueIndex" json:"-"~
	PendingKey     *string    ~gorm:"size:64;uniqueIndex:ux_org_invitation_pending" json:"-"~
	ExpiresAt      time.Time  ~gorm:"not null;index" json:"expires_at"~
	AcceptedAt     *time.Time ~json:"accepted_at,omitempty"~
	RevokedAt      *time.Time ~json:"revoked_at,omitempty"~
	CreatedAt      time.Time  ~json:"created_at"~
}

func (i *OrganizationInvitation) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" { i.ID = ids.New() }
	return nil
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

func mtInvitationHandler(ctx Context) string {
	src := `package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/tenant"
)

const defaultOrganizationInvitationTTL = 7 * 24 * time.Hour

// The web app's origin, not the API's. In a monorepo APP_URL is the API, so the
// invitation link pointed at a host with no accept page and every invite opened
// a 404. GRIT_WEB_URL wins in every environment; in development a monorepo's web
// app is the scaffolded port; otherwise APP_URL is the only address on record.
const defaultOrganizationWebURL = "{{WEB_DEFAULT}}"

func organizationWebBaseURL() string {
	if value := strings.TrimRight(os.Getenv("GRIT_WEB_URL"), "/"); value != "" {
		return value
	}
	if value := strings.TrimRight(os.Getenv("WEB_URL"), "/"); value != "" {
		return value
	}
	if defaultOrganizationWebURL != "" && os.Getenv("APP_ENV") != "production" {
		return defaultOrganizationWebURL
	}
	if value := strings.TrimRight(os.Getenv("APP_URL"), "/"); value != "" {
		return value
	}
	return "http://localhost:3000"
}

func organizationInvitationTTL() time.Duration {
	configured, err := time.ParseDuration(strings.TrimSpace(os.Getenv("GRIT_ORG_INVITATION_TTL")))
	if err != nil || configured < time.Hour || configured > 30*24*time.Hour {
		return defaultOrganizationInvitationTTL
	}
	return configured
}

type OrganizationInvitationHandler struct {
	DB *gorm.DB
	Mailer *mail.Mailer
}

func NewOrganizationInvitationHandler(db *gorm.DB, mailer *mail.Mailer) *OrganizationInvitationHandler {
	return &OrganizationInvitationHandler{DB: db, Mailer: mailer}
}

type organizationInvitationInput struct {
	Email string ~json:"email" binding:"required,email"~
	Role string ~json:"role"~
}

func (h *OrganizationInvitationHandler) Create(c *gin.Context) {
	orgID := c.Param("id")
	if !h.mayInvite(c, orgID) { return }
	var in organizationInvitationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code":"VALIDATION_ERROR", "message":"A valid email is required"}})
		return
	}
	emailAddress := strings.ToLower(strings.TrimSpace(in.Email))
	role := strings.ToLower(strings.TrimSpace(in.Role))
	if role == "" { role = "member" }
	if role != "member" && role != "admin" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code":"VALIDATION_ERROR", "message":"Invitations may assign the member or admin role"}})
		return
	}
	var org models.Organization
	if err := tenant.Unscoped(h.DB).Where("id = ? AND active = ?", orgID, true).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND", "message":"Organization not found"}})
		return
	}
	userID, _ := c.Get("user_id")
	inviter, _ := userID.(string)
	if role == "admin" && org.OwnerID != inviter && !authz.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code":"FORBIDDEN", "message":"Only the organization owner may invite an admin"}})
		return
	}
	var memberCount int64
	if err := tenant.Unscoped(h.DB).Model(&models.OrganizationMember{}).
		Joins("JOIN users ON users.id = organization_members.user_id").
		Where("organization_members.organization_id = ? AND lower(users.email) = ?", orgID, emailAddress).
		Count(&memberCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL_ERROR", "message":"Could not check organization membership"}})
		return
	}
	if memberCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code":"CONFLICT", "message":"That email already belongs to this organization"}})
		return
	}
	// Expired invitations no longer reserve the unique pending key. Keeping
	// them visible as revoked history avoids making an expired invite block a
	// fresh invitation forever.
	_ = tenant.Unscoped(h.DB).Model(&models.OrganizationInvitation{}).
		Where("organization_id = ? AND accepted_at IS NULL AND revoked_at IS NULL AND expires_at <= ?", orgID, time.Now()).
		Updates(map[string]any{"revoked_at":time.Now(), "pending_key":nil}).Error
	var duplicate models.OrganizationInvitation
	if err := tenant.Unscoped(h.DB).Where("organization_id = ? AND lower(email) = ? AND accepted_at IS NULL AND revoked_at IS NULL AND expires_at > ?", orgID, emailAddress, time.Now()).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code":"CONFLICT", "message":"A pending invitation already exists for that email"}})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL_ERROR", "message":"Could not check pending invitations"}})
		return
	}
	if h.Mailer == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{"code":"MAIL_UNAVAILABLE", "message":"Email delivery is not configured"}})
		return
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL_ERROR", "message":"Could not create invitation"}})
		return
	}
	token := hex.EncodeToString(tokenBytes)
	tokenSum := sha256.Sum256([]byte(token))
	keySum := sha256.Sum256([]byte(orgID + "\x00" + emailAddress))
	pendingKey := hex.EncodeToString(keySum[:])
	invitation := models.OrganizationInvitation{
		OrganizationID: orgID, Email: emailAddress, Role: role, InvitedBy: inviter,
		TokenHash: hex.EncodeToString(tokenSum[:]), PendingKey: &pendingKey,
		ExpiresAt: time.Now().Add(organizationInvitationTTL()),
	}
	if err := tenant.Unscoped(h.DB).Create(&invitation).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code":"CONFLICT", "message":"A pending invitation already exists for that email"}})
		return
	}
	baseURL := organizationWebBaseURL()
	acceptURL := baseURL + "/invitations/accept?token=" + url.QueryEscape(token)
	ttlHours := int(organizationInvitationTTL().Hours())
	body := fmt.Sprintf("<p>You have been invited to join <strong>%s</strong>.</p><p><a href=\"%s\">Review invitation</a></p><p>This link expires in %d hours.</p>", html.EscapeString(org.Name), html.EscapeString(acceptURL), ttlHours)
	if err := h.Mailer.SendRaw(c.Request.Context(), emailAddress, "Invitation to "+org.Name, body); err != nil {
		now := time.Now()
		_ = tenant.Unscoped(h.DB).Model(&models.OrganizationInvitation{}).Where("id = ?", invitation.ID).Updates(map[string]any{"revoked_at": now, "pending_key": nil}).Error
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code":"MAIL_FAILED", "message":"The invitation could not be emailed; no pending invitation was left"}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": invitation, "message":"Invitation sent"})
}

func (h *OrganizationInvitationHandler) List(c *gin.Context) {
	orgID := c.Param("id")
	if !h.mayInvite(c, orgID) { return }
	var rows []models.OrganizationInvitation
	if err := tenant.Unscoped(h.DB).Where("organization_id = ? AND accepted_at IS NULL AND revoked_at IS NULL", orgID).Order("created_at desc").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL_ERROR", "message":"Could not load invitations"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func (h *OrganizationInvitationHandler) Revoke(c *gin.Context) {
	orgID := c.Param("id")
	if !h.mayInvite(c, orgID) { return }
	now := time.Now()
	result := tenant.Unscoped(h.DB).Model(&models.OrganizationInvitation{}).
		Where("id = ? AND organization_id = ? AND accepted_at IS NULL AND revoked_at IS NULL", c.Param("invitationId"), orgID).
		Updates(map[string]any{"revoked_at": now, "pending_key": nil})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL_ERROR", "message":"Could not revoke invitation"}})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND", "message":"Invitation not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message":"Invitation revoked"})
}

func (h *OrganizationInvitationHandler) Resend(c *gin.Context) {
	orgID := c.Param("id")
	if !h.mayInvite(c, orgID) { return }
	if h.Mailer == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{"code":"MAIL_UNAVAILABLE", "message":"Email delivery is not configured"}})
		return
	}
	var invitation models.OrganizationInvitation
	if err := tenant.Unscoped(h.DB).Where("id = ? AND organization_id = ? AND accepted_at IS NULL AND revoked_at IS NULL", c.Param("invitationId"), orgID).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND", "message":"Pending invitation not found"}})
		return
	}
	var org models.Organization
	if err := tenant.Unscoped(h.DB).Where("id = ? AND active = ?", orgID, true).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND", "message":"Organization not found"}})
		return
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL_ERROR", "message":"Could not resend invitation"}})
		return
	}
	token := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	now := time.Now()
	result := tenant.Unscoped(h.DB).Model(&models.OrganizationInvitation{}).
		Where("id = ? AND token_hash = ? AND accepted_at IS NULL AND revoked_at IS NULL", invitation.ID, invitation.TokenHash).
		Updates(map[string]any{"token_hash":hex.EncodeToString(sum[:]), "expires_at":now.Add(organizationInvitationTTL()), "created_at":now})
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code":"CONFLICT", "message":"Invitation is no longer pending"}})
		return
	}
	baseURL := organizationWebBaseURL()
	acceptURL := baseURL + "/invitations/accept?token=" + url.QueryEscape(token)
	ttlHours := int(organizationInvitationTTL().Hours())
	body := fmt.Sprintf("<p>Your invitation to join <strong>%s</strong> is still waiting.</p><p><a href=\"%s\">Review invitation</a></p><p>This link expires in %d hours.</p>", html.EscapeString(org.Name), html.EscapeString(acceptURL), ttlHours)
	if err := h.Mailer.SendRaw(c.Request.Context(), invitation.Email, "Invitation to "+org.Name, body); err != nil {
		_ = tenant.Unscoped(h.DB).Model(&models.OrganizationInvitation{}).
			Where("id = ? AND token_hash = ?", invitation.ID, hex.EncodeToString(sum[:])).
			Updates(map[string]any{"revoked_at":time.Now(), "pending_key":nil}).Error
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code":"MAIL_FAILED", "message":"The invitation email could not be resent"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message":"Invitation resent"})
}

type acceptOrganizationInvitationInput struct { Token string ~json:"token" binding:"required"~ }

var (
	errInvitationEmailMismatch = errors.New("the invitation belongs to a different email")
	errInvitationExpired       = errors.New("the invitation expired")
)

// Preview answers what an invitation is before anybody signs in.
//
// The page needs the organization's name and the invited address to offer the
// right action: sign in, create the account, or say the link is for somebody
// else. A token is 256 bits from crypto/rand, so it is the credential; no
// session is required to look one up, and nothing beyond the invitation's own
// fields is returned.
func (h *OrganizationInvitationHandler) Preview(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "INVITATION_INVALID", "message": "This invitation link is missing its token"}})
		return
	}
	sum := sha256.Sum256([]byte(token))
	var invitation models.OrganizationInvitation
	if err := tenant.Unscoped(h.DB).Where("token_hash = ?", hex.EncodeToString(sum[:])).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "INVITATION_INVALID", "message": "This invitation is invalid"}})
		return
	}
	status := "pending"
	switch {
	case invitation.AcceptedAt != nil:
		status = "accepted"
	case invitation.RevokedAt != nil:
		status = "revoked"
	case time.Now().After(invitation.ExpiresAt):
		status = "expired"
	}
	var org models.Organization
	_ = tenant.Unscoped(h.DB).Where("id = ?", invitation.OrganizationID).First(&org).Error
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"organization_name": org.Name,
		"email":             invitation.Email,
		"role":              invitation.Role,
		"status":            status,
		"expires_at":        invitation.ExpiresAt,
	}})
}

// Accept turns a pending invitation into a membership.
//
// Idempotent on purpose, which is how the reference implementation behaves:
// somebody who clicks the link twice, or returns after the first request
// already added them, gets success and lands in the organization instead of an
// "invalid" page for an invitation they hold. A token from another address is
// refused explicitly, because "belongs to a different email" is something the
// reader can act on.
func (h *OrganizationInvitationHandler) Accept(c *gin.Context) {
	var in acceptOrganizationInvitationInput
	if err := c.ShouldBindJSON(&in); err != nil || strings.TrimSpace(in.Token) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "Invitation token is required"}})
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	var user models.User
	if err := tenant.Unscoped(h.DB).Where("id = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Sign in before accepting this invitation"}})
		return
	}
	if user.EmailVerifiedAt == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "EMAIL_NOT_VERIFIED", "message": "Verify your email before accepting an organization invitation"}})
		return
	}
	sum := sha256.Sum256([]byte(in.Token))
	hash := hex.EncodeToString(sum[:])

	var organizationID, organizationName string
	alreadyMember := false
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.OrganizationInvitation
		if err := tenant.Unscoped(tx).Where("token_hash = ? AND revoked_at IS NULL", hash).First(&invitation).Error; err != nil {
			return err
		}
		if !strings.EqualFold(strings.TrimSpace(user.Email), invitation.Email) {
			return errInvitationEmailMismatch
		}
		var org models.Organization
		if err := tenant.Unscoped(tx).Where("id = ? AND active = ?", invitation.OrganizationID, true).First(&org).Error; err != nil {
			return err
		}
		organizationID, organizationName = org.ID, org.Name
		if invitation.AcceptedAt != nil {
			// The link was already used. It is only a way back in for the
			// person it added; for anybody else it is spent.
			var member models.OrganizationMember
			if err := tenant.Unscoped(tx).Where("user_id = ? AND organization_id = ?", uid, org.ID).First(&member).Error; err != nil {
				return gorm.ErrRecordNotFound
			}
			alreadyMember = true
			return nil
		}
		if time.Now().After(invitation.ExpiresAt) {
			return errInvitationExpired
		}
		var existing models.OrganizationMember
		switch err := tenant.Unscoped(tx).Where("user_id = ? AND organization_id = ?", uid, org.ID).First(&existing).Error; {
		case err == nil:
			alreadyMember = true
		case errors.Is(err, gorm.ErrRecordNotFound):
			member := models.OrganizationMember{UserID: uid, OrganizationID: org.ID, Role: invitation.Role}
			if err := tenant.Unscoped(tx).Create(&member).Error; err != nil {
				return err
			}
		default:
			return err
		}
		now := time.Now()
		res := tenant.Unscoped(tx).Model(&models.OrganizationInvitation{}).
			Where("id = ? AND accepted_at IS NULL AND revoked_at IS NULL AND expires_at > ?", invitation.ID, now).
			Updates(map[string]any{"accepted_at": now, "pending_key": nil})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})

	switch {
	case errors.Is(err, errInvitationEmailMismatch):
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "EMAIL_MISMATCH", "message": "This invitation belongs to a different email address"}})
		return
	case errors.Is(err, errInvitationExpired):
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "INVITATION_INVALID", "message": "This invitation has expired. Ask for a new one."}})
		return
	case err != nil:
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "INVITATION_INVALID", "message": "This invitation is invalid, revoked, or belongs to a different email"}})
		return
	}

	message := "Invitation accepted"
	if alreadyMember {
		message = "You already belong to this organization"
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"organization_id":   organizationID,
		"organization_name": organizationName,
		"already_member":    alreadyMember,
	}, "message": message})
}

func (h *OrganizationInvitationHandler) mayInvite(c *gin.Context, orgID string) bool {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	var org models.Organization
	if err := tenant.Unscoped(h.DB).Where("id = ? AND active = ?", orgID, true).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND", "message":"Organization not found"}})
		return false
	}
	if authz.IsAdmin(c) { return true }
	if org.OwnerID == uid { return true }
	if grants, ok := c.Get("user_grants"); ok {
		if list, ok := grants.([]string); ok {
			for _, grant := range list { if grant == "orgmembers.manage" { return true } }
		}
	}
	c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code":"FORBIDDEN", "message":"Only the organization owner or a member manager may manage invitations"}})
	return false
}
`
	webDefault := ""
	switch ctx.Architecture {
	case "triple", "double", "full":
		webDefault = "http://localhost:3000"
	}
	out := strings.ReplaceAll(src, "{{WEB_DEFAULT}}", webDefault)
	return strings.ReplaceAll(strings.ReplaceAll(out, "~", "`"), "{{MODULE}}", ctx.Module)
}
