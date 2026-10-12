package scaffold

// The published key set, at /.well-known/jwks.json.
//
// This is the endpoint that makes asymmetric signing operational rather than
// theoretical. Hardcoding a public key into every service that verifies turns
// rotation into a synchronised redeploy of the whole estate, which in practice
// means the key is never rotated. Publishing it means a verifier fetches the
// set, caches it, and refetches when a token names a key it has not seen.
//
// It is served unauthenticated and that is correct: everything in it is public
// by construction. On a project using a shared secret it serves an empty set
// rather than 404, so a verifier gets "this issuer publishes nothing" instead
// of "this URL is wrong", which are different problems.

// apiJWKSHandlerGo emits internal/handlers/jwks.go.
func apiJWKSHandlerGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/services"
)

// JWKSHandler serves the public keys that verify this API's tokens.
type JWKSHandler struct {
	Keys *services.SigningKeys
}

// Serve answers GET /.well-known/jwks.json.
//
// Cacheable, because the whole point is that a verifier holds these in memory
// rather than fetching one per request. Five minutes is short enough that a
// rotation propagates on its own and long enough that this endpoint is not
// traffic. A verifier that meets an unknown key id should refetch immediately
// rather than waiting for the cache to lapse, which is what makes rotation
// seamless.
func (h *JWKSHandler) Serve(c *gin.Context) {
	set := services.JWKS{Keys: []services.JWK{}}
	if h.Keys != nil {
		set = h.Keys.PublicJWKS()
	}
	c.Header("Cache-Control", "public, max-age=300")
	c.JSON(http.StatusOK, set)
}
`
}

// apiJWKSHandlerTestGo emits internal/handlers/jwks_test.go.
func apiJWKSHandlerTestGo() string {
	return `package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/services"
)

func serveJWKS(t *testing.T, keys *services.SigningKeys) (int, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &JWKSHandler{Keys: keys}
	router.GET("/.well-known/jwks.json", h.Serve)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil))
	return rec.Code, rec.Body.String()
}

// A shared secret must never be published, and the endpoint must still exist,
// so a verifier learns "nothing to fetch" rather than "wrong URL".
func TestJWKSPublishesNothingOnASharedSecret(t *testing.T) {
	keys, err := services.LoadSigningKeys("a-secret-at-least-32-characters-long", "", "", "", "")
	if err != nil {
		t.Fatalf("loading keys: %v", err)
	}

	code, body := serveJWKS(t, keys)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200 even with nothing to publish", code)
	}

	var set services.JWKS
	if err := json.Unmarshal([]byte(body), &set); err != nil {
		t.Fatalf("response is not a key set: %v", err)
	}
	if len(set.Keys) != 0 {
		t.Fatal("an HS256 deployment must publish no keys: the verification key is the signing key")
	}
	// A null rather than an empty array breaks clients that iterate it.
	if !strings.Contains(body, ` + "`" + `"keys":[]` + "`" + `) {
		t.Fatalf("empty set should serialise as [], got %s", body)
	}
}

func TestJWKSIsServedWithoutKeysConfigured(t *testing.T) {
	code, _ := serveJWKS(t, nil)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
}
`
}
