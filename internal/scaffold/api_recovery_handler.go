package scaffold

func recoveryHandlerGo() string {
	return `package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"{{MODULE}}/internal/sms"
)

// RecoveryHandler manages the recovery contacts on the signed-in account.
//
// Every write here takes the account password, and that is the whole security
// model of this file. A recovery address is a second way in, so somebody who
// gets a live session, from a borrowed laptop or a stolen token, could
// otherwise quietly attach their own address and own the account from then on.
// The password is what they do not have.
type RecoveryHandler struct {
	DB   *gorm.DB
	Mail *mail.Mailer
}

func NewRecoveryHandler(db *gorm.DB, m *mail.Mailer) *RecoveryHandler {
	return &RecoveryHandler{DB: db, Mail: m}
}

type setRecoveryRequest struct {
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
	Email    string ` + "`" + `json:"email"` + "`" + `
	Phone    string ` + "`" + `json:"phone"` + "`" + `
}

type verifyRecoveryRequest struct {
	Code string ` + "`" + `json:"code" binding:"required"` + "`" + `
}

type clearRecoveryRequest struct {
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
}

func fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

// authed loads the signed-in user and checks the supplied password.
func (h *RecoveryHandler) authed(c *gin.Context, password string) (*models.User, bool) {
	var user models.User
	if err := h.DB.First(&user, "id = ?", c.GetString("user_id")).Error; err != nil {
		fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "Not signed in")
		return nil, false
	}
	// An OAuth-only account has no password to check, so it cannot pass this
	// gate. Saying so plainly beats "incorrect password" on an account that
	// has never had one.
	if user.Password == "" {
		fail(c, http.StatusBadRequest, "NO_PASSWORD",
			"This account signs in with a provider. Set a password before adding a recovery contact.")
		return nil, false
	}
	if !user.CheckPassword(password) {
		fail(c, http.StatusForbidden, "INVALID_PASSWORD", "That password is not correct")
		return nil, false
	}
	return &user, true
}

// Overview reports what is switched on and what this deployment can offer.
//
// The capability half matters: phone recovery only exists if somebody wired an
// SMS provider, and the admin uses this to leave the card out entirely rather
// than render a button that cannot work.
func (h *RecoveryHandler) Overview(c *gin.Context) {
	var user models.User
	if err := h.DB.First(&user, "id = ?", c.GetString("user_id")).Error; err != nil {
		fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "Not signed in")
		return
	}

	contacts, err := services.LoadRecoveryContacts(h.DB, user.ID)
	if err != nil {
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not load your security settings")
		return
	}
	email, hasEmail := contacts[models.RecoveryEmail]
	phone, hasPhone := contacts[models.RecoveryPhone]

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"email":                   user.Email,
		"email_verified":          user.EmailVerifiedAt != nil,
		"has_password":            user.Password != "",
		"provider":                user.Provider,
		"recovery_email":          maskEmail(email.Destination),
		"recovery_email_verified": hasEmail,
		"recovery_phone":          maskPhone(phone.Destination),
		"recovery_phone_verified": hasPhone,
		"sms_provider_configured": sms.Configured(),
	}})
}

// maskEmail shows enough to recognise an address without publishing it.
//
// The overview is read by a session, and a session is exactly what an attacker
// on a borrowed laptop has. Returning the full recovery address would tell them
// where to go next.
func maskEmail(address string) string {
	if address == "" {
		return ""
	}
	at := strings.LastIndex(address, "@")
	if at <= 0 {
		return "***"
	}
	local, domain := address[:at], address[at+1:]
	if len(local) <= 2 {
		return "**@" + domain
	}
	return local[:1] + strings.Repeat("*", len(local)-2) + local[len(local)-1:] + "@" + domain
}

func maskPhone(number string) string {
	if number == "" {
		return ""
	}
	if len(number) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(number)-4) + number[len(number)-4:]
}

// SetEmail starts adding or replacing the recovery email.
//
// The contact is not stored yet. It lives on the token until a code proves
// somebody can read mail there, because an unverified recovery address is
// worse than none: it looks like a way back in and is not.
func (h *RecoveryHandler) SetEmail(c *gin.Context) {
	var req setRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	user, ok := h.authed(c, req.Password)
	if !ok {
		return
	}

	address := strings.ToLower(strings.TrimSpace(req.Email))
	if address == "" || !strings.Contains(address, "@") {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "A valid email address is required")
		return
	}
	if err := services.ValidateRecoveryEmail(h.DB, user.ID, user.Email, address); err != nil {
		if errors.Is(err, services.ErrRecoverySameAsPrimary) || errors.Is(err, services.ErrRecoveryInUse) {
			fail(c, http.StatusUnprocessableEntity, "INVALID_RECOVERY_ADDRESS", err.Error())
			return
		}
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not check that address")
		return
	}

	code, err := services.NewRecoveryCode()
	if err != nil {
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not start verification")
		return
	}
	if _, err := services.CreateRecoveryToken(h.DB, user.ID, models.RecoveryEmail, address, code); err != nil {
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not start verification")
		return
	}

	if h.Mail != nil {
		// SendRaw rather than a template: this mail goes to an address that is
		// not yet on the account, so it must carry nothing about the account
		// beyond the code itself.
		_ = h.Mail.SendRaw(c.Request.Context(), address, "Your recovery code",
			"<p>Your recovery code is <strong>"+code+"</strong>. It expires in 15 minutes.</p>"+
				"<p>If you did not ask for this, somebody may have your password. Change it.</p>")
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"sent_to": maskEmail(address)},
		"message": "Enter the code we sent to confirm the address"})
}

// VerifyEmail completes it.
func (h *RecoveryHandler) VerifyEmail(c *gin.Context) {
	h.verify(c, models.RecoveryEmail)
}

// VerifyPhone completes the phone flow.
func (h *RecoveryHandler) VerifyPhone(c *gin.Context) {
	h.verify(c, models.RecoveryPhone)
}

func (h *RecoveryHandler) verify(c *gin.Context, kind models.RecoveryContactKind) {
	var req verifyRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	destination, err := services.ConsumeRecoveryToken(h.DB, c.GetString("user_id"), kind, strings.TrimSpace(req.Code))
	if err != nil {
		// One message for wrong, expired and spent alike. Telling them apart
		// says which codes once existed.
		fail(c, http.StatusUnprocessableEntity, "INVALID_CODE",
			"That code is not valid. Request a new one.")
		return
	}

	masked := maskEmail(destination)
	if kind == models.RecoveryPhone {
		masked = maskPhone(destination)
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"verified": masked}, "message": "Recovery contact confirmed"})
}

// ClearEmail removes the recovery email.
func (h *RecoveryHandler) ClearEmail(c *gin.Context) {
	h.clear(c, models.RecoveryEmail)
}

// ClearPhone removes the recovery phone.
func (h *RecoveryHandler) ClearPhone(c *gin.Context) {
	h.clear(c, models.RecoveryPhone)
}

func (h *RecoveryHandler) clear(c *gin.Context, kind models.RecoveryContactKind) {
	var req clearRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	user, ok := h.authed(c, req.Password)
	if !ok {
		return
	}
	if err := services.ClearRecoveryContact(h.DB, user.ID, kind); err != nil {
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not remove that contact")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recovery contact removed"})
}

// SetPhone starts adding or replacing the recovery phone.
func (h *RecoveryHandler) SetPhone(c *gin.Context) {
	// Checked before the password, so a deployment with no provider says so
	// rather than making somebody type their password to reach a dead end.
	if !sms.Configured() {
		fail(c, http.StatusNotImplemented, "SMS_NOT_CONFIGURED",
			"This deployment has no SMS provider. See internal/sms.")
		return
	}

	var req setRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	user, ok := h.authed(c, req.Password)
	if !ok {
		return
	}

	number := strings.TrimSpace(req.Phone)
	if len(number) < 7 {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "A valid phone number is required")
		return
	}

	code, err := services.NewRecoveryCode()
	if err != nil {
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not start verification")
		return
	}
	if _, err := services.CreateRecoveryToken(h.DB, user.ID, models.RecoveryPhone, number, code); err != nil {
		fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Could not start verification")
		return
	}

	if err := sms.Send(c.Request.Context(), number, "Your recovery code is "+code+". It expires in 15 minutes."); err != nil {
		fail(c, http.StatusBadGateway, "SMS_FAILED", "Could not send the code. Check the number and try again.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"sent_to": maskPhone(number)},
		"message": "Enter the code we texted to confirm the number"})
}
`
}
