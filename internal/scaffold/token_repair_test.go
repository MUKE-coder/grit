package scaffold

import (
	"strings"
	"testing"
)

func TestRepairRefreshHandlerKeepsTheSession(t *testing.T) {
	src := "package handlers\n\nfunc (h *AuthHandler) Login(c *gin.Context) {\n" +
		"\ttokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)\n}\n\n" +
		"func (h *AuthHandler) Refresh(c *gin.Context) {\n" +
		"\tclaims, err := h.AuthService.ValidateToken(refreshToken)\n" +
		"\ttokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)\n}\n\n" +
		"func (h *AuthHandler) Logout(c *gin.Context) {}\n"
	out, fixed, warn := repairRefreshHandlerSource(src)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "handlers/auth.go", out, fixed,
		"h.AuthService.ValidateRefreshToken(refreshToken)",
		"GenerateSessionTokenPair(user.ID, user.Email, user.Role, sessionID)",
	)
	if strings.Count(out, "GenerateTokenPair(user.ID") != 1 {
		t.Error("Login's token pair was changed too; only Refresh keeps a session id")
	}
}

func TestRepairTOTPRecordsSessions(t *testing.T) {
	block := "\ttokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)\n" +
		"\tif err != nil {\n\t\tc.JSON(500, nil)\n\t\treturn\n\t}\n"
	src := "package handlers\n\nimport (\n\t\"net/http\"\n)\n\ntype TOTPHandler struct{}\n\n" +
		"func (h *TOTPHandler) Verify(c *gin.Context) {\n" + block + "\th.AuthService.SetAuthCookies(c, tokens)\n}\n\n" +
		"func (h *TOTPHandler) VerifyBackupCode(c *gin.Context) {\n" + block + "}\n"
	out, fixed, warn := repairTOTPSessionSource(src)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "handlers/totp.go", out, fixed, "\"log\"")
	if n := strings.Count(out, "services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken)"); n != 2 {
		t.Errorf("%d sessions recorded, want 2", n)
	}
	if again, fixed, _ := repairTOTPSessionSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed totp.go again")
	}
}

func TestRepairImpersonateRecordsSessions(t *testing.T) {
	mint := func(owner, id string) string {
		return "\tpair, err := h.Auth.GenerateTokenPair(" + owner + "." + id + ", " + owner + ".Email, " + owner + ".Role)\n" +
			"\tif err != nil {\n\t\trespond.Internal(c, err)\n\t\treturn\n\t}\n"
	}
	src := "package handlers\n\ntype ImpersonateHandler struct{}\n\n" +
		"func (h *ImpersonateHandler) Start(c *gin.Context) {\n" + mint("target", "ID") + "}\n\n" +
		"func (h *ImpersonateHandler) Stop(c *gin.Context) {\n\tclaims, err := h.Auth.ValidateToken(adminToken)\n" + mint("claims", "UserID") + "}\n"
	out, fixed, _ := repairImpersonateSessionSource(src)
	assertRepaired(t, "handlers/impersonate.go", out, fixed,
		"services.CreateSession(h.DB, c, target.ID, pair.RefreshToken)",
		"services.CreateSession(h.DB, c, claims.UserID, pair.RefreshToken)",
		"h.Auth.ValidateAccessToken(adminToken)",
	)
}

func TestRepairAuthServiceWiring(t *testing.T) {
	src := "package routes\n\nfunc Setup() {\n\tauthService := &services.AuthService{\n" +
		"\t\tSecret:        cfg.JWTSecret,\n\t\tAccessExpiry:  cfg.JWTAccessExpiry,\n\t\tRefreshExpiry: cfg.JWTRefreshExpiry,\n\t}\n\t_ = authService\n}\n"
	out, fixed, warn := repairAuthServiceWiringSource(src)
	if len(warn) > 0 {
		t.Fatalf("warnings: %v", warn)
	}
	assertRepaired(t, "routes.go", out, fixed, "DB: db,", `services.RefreshCookiePath = "/api/" + APIVersion + "/auth"`)
}

// A fresh project needs none of it.
func TestFreshTemplatesNeedNoTokenRepair(t *testing.T) {
	for name, c := range map[string]struct {
		src string
		fn  func(string) (string, []string, []string)
	}{
		"middleware/auth.go":   {apiAuthMiddlewareGo(), repairAuthMiddlewareTokenSource},
		"handlers/realtime.go": {apiRealtimeHandlerGo(), repairRealtimeTokenSource},
		"handlers/auth.go":     {apiAuthHandlerGo(), repairRefreshHandlerSource},
		"handlers/totp.go":     {totpHandlerGo(), repairTOTPSessionSource},
		"routes/routes.go":     {apiRoutesGo(), repairAuthServiceWiringSource},
	} {
		out, fixed, warn := c.fn(c.src)
		if out != c.src || len(fixed) > 0 || len(warn) > 0 {
			t.Errorf("%s: the scaffold template still needs the repair (fixed %v, warn %v)", name, fixed, warn)
		}
	}
}
