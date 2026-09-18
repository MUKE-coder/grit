package scaffold

func passkeyHandlerGo() string {
	return `package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"{{MODULE}}/internal/respond"
)

// PasskeyHandler exposes the two WebAuthn ceremonies.
//
// Registration is behind auth: you add a passkey to an account you are already
// signed in to. Sign-in is public by necessity, and is the one place where the
// challenge stored server-side is doing the real work.
type PasskeyHandler struct {
	DB       *gorm.DB
	Passkeys *services.Passkeys
	Auth     *AuthHandler
}

func NewPasskeyHandler(db *gorm.DB, p *services.Passkeys, auth *AuthHandler) *PasskeyHandler {
	return &PasskeyHandler{DB: db, Passkeys: p, Auth: auth}
}

// available guards every route: a deployment with no usable origin has no
// relying party, and saying so beats a nil dereference.
func (h *PasskeyHandler) available(c *gin.Context) bool {
	if h.Passkeys == nil {
		respond.Fail(c, respond.CodePasskeysNotConfigured, "Passkeys are not configured on this deployment")
		return false
	}
	return true
}

// List returns the signed-in user's passkeys.
func (h *PasskeyHandler) List(c *gin.Context) {
	if !h.available(c) {
		return
	}
	rows, err := h.Passkeys.List(c.GetString("user_id"))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not load your passkeys")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// BeginRegistration hands the browser the creation options.
func (h *PasskeyHandler) BeginRegistration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	opts, sessionID, err := h.Passkeys.BeginRegistration(c.GetString("user_id"))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not start registration")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"session_id": sessionID,
		"options":    opts.Response,
	}})
}

// FinishRegistration verifies the attestation and stores the credential.
//
// The body carries both the envelope and the raw credential the browser
// produced, so it is read once and handed on whole: re-encoding an
// authenticator's response is a good way to invalidate the signature over it.
func (h *PasskeyHandler) FinishRegistration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respond.Fail(c, respond.CodeInvalidBody, "Could not read the request")
		return
	}

	sessionID := c.Query("session")
	name := c.Query("name")
	if sessionID == "" {
		respond.Fail(c, respond.CodeValidationError, "session is required")
		return
	}

	key, err := h.Passkeys.FinishRegistration(sessionID, name, raw)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if errors.Is(err, services.ErrPasskeyChallengeGone) {
			status = http.StatusGone
		}
		c.JSON(status, gin.H{"error": gin.H{
			"code": "PASSKEY_REJECTED", "message": err.Error(),
		}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": key, "message": "Passkey added"})
}

// BeginLogin starts a usernameless sign-in. Public.
func (h *PasskeyHandler) BeginLogin(c *gin.Context) {
	if !h.available(c) {
		return
	}
	opts, sessionID, err := h.Passkeys.BeginLogin()
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not start sign-in")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"session_id": sessionID,
		"options":    opts.Response,
	}})
}

// FinishLogin verifies the assertion and issues the same tokens a password
// sign-in would, so everything downstream is unchanged.
func (h *PasskeyHandler) FinishLogin(c *gin.Context) {
	if !h.available(c) {
		return
	}
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respond.Fail(c, respond.CodeInvalidBody, "Could not read the request")
		return
	}
	sessionID := c.Query("session")
	if sessionID == "" {
		respond.Fail(c, respond.CodeValidationError, "session is required")
		return
	}

	user, err := h.Passkeys.FinishLogin(sessionID, raw)
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, services.ErrPasskeyChallengeGone) {
			status = http.StatusGone
		}
		c.JSON(status, gin.H{"error": gin.H{
			"code": "PASSKEY_REJECTED", "message": err.Error(),
		}})
		return
	}
	if !user.Active {
		respond.Fail(c, respond.CodeAccountDisabled, "That account is disabled")
		return
	}

	// The same tokens a password sign-in issues, recorded as a session so the
	// device shows up in Active Sessions and can be revoked like any other.
	tokens, err := h.Auth.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not complete sign-in")
		return
	}
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("passkey: failed to record session for %s: %v", user.ID, err)
	}
	h.Auth.AuthService.SetAuthCookies(c, tokens)
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"user": user, "tokens": tokens},
		"message": "Signed in with a passkey",
	})
}

type renamePasskeyRequest struct {
	Name string ` + "`" + `json:"name" binding:"required"` + "`" + `
}

// Rename changes the label. Somebody with four passkeys needs to know which is
// the laptop they left at the office.
func (h *PasskeyHandler) Rename(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var req renamePasskeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		respond.Fail(c, respond.CodeValidationError, "A name is required")
		return
	}
	res := h.DB.WithContext(c.Request.Context()).Model(&models.Passkey{}).
		Where("id = ? AND user_id = ?", c.Param("id"), c.GetString("user_id")).
		Update("name", name)
	if res.Error != nil || res.RowsAffected == 0 {
		respond.Fail(c, respond.CodeNotFound, "No such passkey")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Renamed"})
}

// Delete removes one of the signed-in user's passkeys.
func (h *PasskeyHandler) Delete(c *gin.Context) {
	if !h.available(c) {
		return
	}
	if err := h.Passkeys.Delete(c.GetString("user_id"), c.Param("id")); err != nil {
		respond.Fail(c, respond.CodeNotFound, "No such passkey")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Passkey removed"})
}
`
}
