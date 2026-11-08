package scaffold

import (
	"strings"
	"testing"
)

// The "EOF" a reader saw on the invitation page came from here: the request
// carried no session, the client tried to refresh, and Refresh bound an empty
// body and answered 422 with the decoder's message. A caller can act on 401;
// it can do nothing with "EOF".
func TestRefreshWithoutASessionAnswersAnAuthError(t *testing.T) {
	src := withModule(apiAuthHandlerGo())

	if !strings.Contains(src, "errors.Is(err, io.EOF)") {
		t.Error("Refresh still reports an empty body as a validation error, so a " +
			"client's 401-refresh path surfaces \"EOF\" instead of the 401")
	}
	if !strings.Contains(src, `respond.Fail(c, respond.CodeInvalidToken, "No refresh token")`) {
		t.Error("Refresh does not answer INVALID_TOKEN when there is no refresh " +
			"token to spend")
	}
}

// And the client keeps the original failure. A page that receives 401 can ask
// for a sign-in; it cannot do anything with "No refresh token".
func TestTheWebClientKeepsTheOriginalFailureWhenRefreshFails(t *testing.T) {
	src := webAPIClient()
	if !strings.Contains(src, "try {\n    await refreshing;\n  } catch {") {
		t.Fatal("the web client lets the refresh error replace the request's own " +
			"failure, so an anonymous request shows the refresh message")
	}
	if !strings.Contains(src, "throw error;") {
		t.Error("the web client does not rethrow the original error")
	}
}

// The empty-body fix has to reach a handler that was already repaired for the
// session id: `grit upgrade` runs every repair, and a project scaffolded before
// this change has GenerateSessionTokenPair but still binds an empty body.
func TestRefreshEmptyBodyRepairReachesAnAlreadyRepairedHandler(t *testing.T) {
	older := `package handlers

import (
	"errors"
)

func (h *AuthHandler) Refresh(c *gin.Context) {
	var refreshToken string
	if cookieValue, err := c.Cookie("grit_refresh"); err == nil && cookieValue != "" {
		refreshToken = cookieValue
	} else {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respond.Fail(c, respond.CodeValidationError, err.Error())
			return
		}
		refreshToken = req.RefreshToken
	}
	tokens, err := h.AuthService.GenerateSessionTokenPair(user.ID, user.Email, user.Role, sessionID)
	if err != nil {
		return
	}
	_ = tokens
}
`
	out, fixes, warn := repairRefreshHandlerSource(older)
	if len(warn) > 0 {
		t.Fatalf("repair warned: %v", warn)
	}
	if len(fixes) == 0 {
		t.Fatal("the repair did not touch a handler that still binds an empty body")
	}
	if !strings.Contains(out, "errors.Is(err, io.EOF)") {
		t.Error("the repaired handler still answers 422 for an empty body")
	}
	if !strings.Contains(out, `"io"`) {
		t.Error("the repair used io.EOF without importing io, so the project stops compiling")
	}
	if again, fixes2, _ := repairRefreshHandlerSource(out); again != out || len(fixes2) != 0 {
		t.Error("the repair is not idempotent: a second upgrade would rewrite the handler again")
	}
}
