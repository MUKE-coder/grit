package scaffold

import (
	"strings"
	"testing"
)

// mainGoBefore is the google.New call exactly as every release up to
// v3.415.0 wrote it into cmd/server/main.go.
const mainGoBefore = `	var oauthProviders []goth.Provider
	if cfg.GoogleClientID != "" {
		oauthProviders = append(oauthProviders, google.New(
			cfg.GoogleClientID, cfg.GoogleClientSecret,
			cfg.AppURL+"/api/auth/oauth/google/callback",
		))
		log.Println("Google OAuth2 configured")
	}
`

func TestGoogleScopeRepairAddsProfile(t *testing.T) {
	out, changes, warnings := repairGoogleScopeSource(mainGoBefore)

	if !strings.Contains(out, `"email", "profile",`) {
		t.Fatalf("the profile scope was not added:\n%s", out)
	}
	if len(changes) != 1 {
		t.Errorf("want one change reported, got %v", changes)
	}
	if len(warnings) != 0 {
		t.Errorf("want no warnings on a file Grit wrote, got %v", warnings)
	}

	// Running it twice must not add the scopes twice.
	again, changes, _ := repairGoogleScopeSource(out)
	if again != out {
		t.Error("the repair is not idempotent: a second pass changed the file again")
	}
	if len(changes) != 0 {
		t.Errorf("a second pass reported %v", changes)
	}
}

// A project with no Google provider configured is left alone and says
// nothing, rather than warning about a file it has no business in.
func TestGoogleScopeRepairIgnoresAProjectWithoutGoogle(t *testing.T) {
	src := "func main() {\n\t// no social login here\n}\n"
	out, changes, warnings := repairGoogleScopeSource(src)
	if out != src || len(changes) != 0 || len(warnings) != 0 {
		t.Errorf("want no action, got out=%q changes=%v warnings=%v", out, changes, warnings)
	}
}

// A main.go somebody has rewritten is told what to add, not edited.
func TestGoogleScopeRepairWarnsOnAnEditedFile(t *testing.T) {
	edited := `		oauthProviders = append(oauthProviders, google.New(
			cfg.GoogleClientID, cfg.GoogleClientSecret,
			buildRedirectURL(cfg), // we build ours
		))`
	out, changes, warnings := repairGoogleScopeSource(edited)
	if out != edited {
		t.Error("an edited main.go was modified; it must be left alone")
	}
	if len(changes) != 0 {
		t.Errorf("want no changes, got %v", changes)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], `"email", "profile"`) {
		t.Errorf("the warning must name the fix, got %v", warnings)
	}
}

// TestScaffoldAsksForTheProfileScope is the other half: the repair fixes an
// old project, and this proves a NEW one never needs it. Two writers, one
// comparing them.
func TestScaffoldAsksForTheProfileScope(t *testing.T) {
	for name, src := range map[string]string{
		"api_files.go":    apiMainGo(Options{}),
		"single_files.go": singleMainGo(Options{}),
	} {
		if !strings.Contains(src, "google.New(") {
			continue
		}
		if !strings.Contains(src, `"email", "profile",`) {
			t.Errorf("%s writes a google.New without the profile scope, so every "+
				"account created through Google is named \"User\"", name)
		}
	}
}
