package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The scaffolded cors.go carries a list, not a literal, so anything that
// introduces a request header can add it from an init() (grit#107).
func TestScaffoldedCORSCarriesAHeaderList(t *testing.T) {
	src := apiCorsMiddlewareGo()

	for _, want := range []string{
		"var allowedRequestHeaders = []string{",
		"var allowedRequestHeadersValue = strings.Join(allowedRequestHeaders",
		"func AllowRequestHeader(name string) {",
		`c.Header("Access-Control-Allow-Headers", allowedRequestHeadersValue)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the CORS template is missing %q", want)
		}
	}
	// The literal is what made every new header somebody's job to remember.
	if strings.Contains(src, `"Access-Control-Allow-Headers", "Origin`) {
		t.Error("the allow list is still written as a string literal")
	}
	// And nothing dropped out of it while it changed shape.
	for _, h := range []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"X-API-Key", "X-CSRF-Token", "Idempotency-Key", "X-Public-IP-Hint",
	} {
		if !strings.Contains(src, "\t\""+h+"\",") {
			t.Errorf("%s is no longer in the allow list", h)
		}
	}
}

// The repair brings an existing project's cors.go to the same shape, and keeps
// whatever that project had in its literal: it reshapes the list, it does not
// replace the contents.
func TestCORSRepairKeepsTheHeadersItFinds(t *testing.T) {
	before := `package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Internal-Key")
		c.Next()
	}
}
`
	out, fixed, warnings := repairCORSHeaderListSource(before)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(fixed) == 0 {
		t.Fatal("the repair reported nothing")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("the repaired file is not valid Go: %v\n%s", err, out)
	}
	for _, want := range []string{
		`"Origin",`,
		`"Content-Type",`,
		// The hand-added one. Dropping it would break that project's own
		// client while fixing a bug it did not have.
		`"X-Internal-Key",`,
		"func AllowRequestHeader(name string) {",
		`c.Header("Access-Control-Allow-Headers", allowedRequestHeadersValue)`,
		`"strings"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired file is missing %q", want)
		}
	}
}

// Twice is a no-op. An upgrade runs on a project that has already been
// upgraded more often than on one that has not.
func TestCORSRepairIsIdempotent(t *testing.T) {
	once, _, _ := repairCORSHeaderListSource(`package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSDynamic(resolve func() []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Origin, Authorization")
		c.Next()
	}
}
`)
	twice, fixed, warnings := repairCORSHeaderListSource(once)
	if twice != once {
		t.Error("the second run changed the file again")
	}
	if len(fixed) != 0 || len(warnings) != 0 {
		t.Errorf("the second run reported something: %v %v", fixed, warnings)
	}
}

// A file that has been rewritten by hand is left alone and says so. Guessing
// at a rewrite of a security-relevant middleware is not worth being wrong,
// and the consequence of skipping it is invisible otherwise.
func TestCORSRepairWarnsRatherThanGuessing(t *testing.T) {
	before := `package middleware

func CORS() {}
`
	out, fixed, warnings := repairCORSHeaderListSource(before)
	if out != before {
		t.Error("the file was changed")
	}
	if len(fixed) != 0 {
		t.Errorf("it claimed a fix: %v", fixed)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "AllowRequestHeader") {
		t.Errorf("expected one warning naming AllowRequestHeader, got %v", warnings)
	}
}
