package scaffold

import (
	"strings"
	"testing"
)

const loginV244 = `package handlers

import (
	"log"
	"net/http"
)

func (h *AuthHandler) Login(c *gin.Context) {
	if !user.CheckPassword(req.Password) {
		h.registerFailedLogin(&user)
		return
	}

	// Any successful password check clears the counter, including one that
	// still has 2FA ahead of it — the password was correct, which is what this
	// counter measures.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		h.DB.Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil})
	}

	// Check if user has TOTP enabled
	var totpConfig models.TwoFactorConfig
	if err := h.DB.Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err == nil {
		if !IsTrustedDevice(c, h.DB, user.ID) {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"totp_required": true,
				},
				"message": "Two-factor authentication required",
			})
			return
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	_ = tokens
	_ = err
}
`

func TestRepairLoginClearsTheCountAfterTwoFactor(t *testing.T) {
	out, fixed, warn := repairLoginFailureCountSource(loginV244)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "handlers/auth.go", out, fixed, "The failure count is cleared when the sign-in completes")
	if strings.Contains(out, "Any successful password check clears the counter") {
		t.Error("the early reset is still there")
	}
	reset := strings.Index(out, "The failure count is cleared")
	branch := strings.Index(out, `"Two-factor authentication required"`)
	mint := strings.Index(out, "GenerateTokenPair(")
	if !(branch < reset && reset < mint) {
		t.Errorf("the reset is not between the 2FA branch and the token pair:\n%s", out)
	}
	if again, fixed, _ := repairLoginFailureCountSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed Login again")
	}
}

func TestFreshTemplatesNeedNoTwoFactorRepair(t *testing.T) {
	src := apiAuthHandlerGo()
	if out, fixed, warn := repairLoginFailureCountSource(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("the scaffold's Login still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
	h := totpHandlerGo()
	for _, want := range []string{"func (h *TOTPHandler) failSecondFactor(", "last_used_step < ?", "TOTP_ALREADY_ENABLED", "user.LockedUntil != nil"} {
		if !strings.Contains(h, want) {
			t.Errorf("the TOTP handler template is missing %s", want)
		}
	}
	if !strings.Contains(totpServiceGo(), "func ValidateCodeStep(") || !strings.Contains(twoFactorModelsGo(), "LastUsedStep") {
		t.Error("the TOTP package or models do not carry the step")
	}
}
