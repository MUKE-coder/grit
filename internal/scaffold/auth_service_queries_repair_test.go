package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// authServiceBefore is an AuthService as a project before v3.339.0 has it: the
// token and cookie work, and none of the sign-in queries, because the handler
// still ran those itself.
const authServiceBefore = `package services

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example.com/app/internal/ids"
)

// AuthService handles JWT token operations.
type AuthService struct {
	Secret string
	DB     *gorm.DB
}

func (s *AuthService) Mint() (string, error) {
	_ = rand.Reader
	_ = ids.New
	_ = time.Now
	_ = errors.New
	return "", nil
}

func (s *AuthService) ClearAuthCookies(c *gin.Context) {
	c.SetCookie("grit_access", "", -1, "/", "", true, true)
}
`

// The sign-in repair rewrites Login to call h.AuthService.UserByEmail. A project
// whose service has no such method would be handed a file that does not compile
// by an upgrade that reported success, so this runs first and adds the methods.
func TestAuthServiceQueriesRepairAddsTheMethods(t *testing.T) {
	out, fixed, warnings := repairAuthServiceQueriesSource(authServiceBefore)
	if len(warnings) != 0 {
		t.Fatalf("warned on the service Grit wrote: %v", warnings)
	}
	if len(fixed) == 0 {
		t.Fatal("reported no change")
	}
	for _, method := range []string{"UserByEmail", "UserByID", "EnabledTwoFactor", "StartPendingTOTP", "ClearLoginFailures"} {
		if !strings.Contains(out, "func (s *AuthService) "+method+"(") {
			t.Errorf("%s is missing, so the repaired Login would not compile", method)
		}
	}
	// The methods take a context and read models, which a token-only service
	// did not import.
	for _, imp := range []string{`"context"`, `"example.com/app/internal/models"`} {
		if !strings.Contains(out, imp) {
			t.Errorf("import %s is missing", imp)
		}
	}
	// isRequestHTTPS sits beside the methods in the template and is already in
	// every older project; a block that carried it would redeclare it.
	if strings.Contains(out, "func isRequestHTTPS(") {
		t.Error("the block carries isRequestHTTPS, which the file already declares")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("the result is not valid Go: %v", err)
	}
	// repairSourceFile runs every repair over a file others have touched, so a
	// second pass has to be a no-op rather than a second copy of the methods.
	again, _, _ := repairAuthServiceQueriesSource(out)
	if again != out {
		t.Error("a second pass changed the file again")
	}
}

// A service somebody rewrote is left alone with a warning, because adding
// methods to a type that may not be there produces a file that does not build.
func TestAuthServiceQueriesRepairLeavesAForeignServiceAlone(t *testing.T) {
	const theirs = "package services\n\n// Our own auth, nothing of Grit's.\ntype Sessions struct{}\n"
	out, fixed, warnings := repairAuthServiceQueriesSource(theirs)
	if out != theirs || len(fixed) != 0 {
		t.Error("changed a service Grit did not write")
	}
	if len(warnings) != 1 {
		t.Fatalf("said nothing about it: %v", warnings)
	}
	if !strings.Contains(warnings[0], "UserByEmail") {
		t.Errorf("the warning does not name what to add: %q", warnings[0])
	}
}
