package scaffold

import (
	"strings"
	"testing"
)

func TestOAuthHookRepairAddsTheCallSite(t *testing.T) {
	fresh := apiAuthOAuthGo()
	mustFormatGo(t, "auth_oauth.go", fresh)

	if !strings.Contains(fresh, "services.RunOAuthLoginHooks(") {
		t.Fatal("the template does not run the social-login hooks")
	}

	old := strings.Replace(fresh, oauthHookCall, "", 1)
	if old == fresh {
		t.Fatal("the fixture is not older than the template")
	}

	out, changes, warnings := repairOAuthHookSource(old)
	if len(warnings) != 0 {
		t.Fatalf("repairing an untouched auth_oauth.go warned: %v", warnings)
	}
	if out != fresh || len(changes) != 1 {
		t.Errorf("the repaired file is not the template (changes %v)", changes)
	}
	mustFormatGo(t, "auth_oauth.go", out)

	if again, changes, _ := repairOAuthHookSource(out); again != out || len(changes) != 0 {
		t.Error("the social-login hook repair is not idempotent")
	}
}

func TestOAuthHookRepairLeavesAnEditedFileAlone(t *testing.T) {
	edited := "package handlers\n\nfunc (h *AuthHandler) OAuthCallback() {}\n"
	out, changes, warnings := repairOAuthHookSource(edited)
	if out != edited || len(changes) != 0 || len(warnings) != 1 {
		t.Errorf("an edited auth_oauth.go should be left alone with a note (changes %v, warnings %v)", changes, warnings)
	}
}
