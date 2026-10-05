package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// proxyEcho answers with what the handlers downstream see: the client IP and
// the X-Forwarded-Proto header.
func proxyEcho(t *testing.T, proxies []string, peer string, headers map[string]string) (string, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := TrustProxies(r, proxies); err != nil {
		t.Fatalf("TrustProxies: %v", err)
	}
	var ip, proto string
	r.GET("/", func(c *gin.Context) {
		ip = c.ClientIP()
		proto = c.GetHeader("X-Forwarded-Proto")
		c.Status(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = peer
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	r.ServeHTTP(httptest.NewRecorder(), req)
	return ip, proto
}

var loopbackAndPrivate = []string{"127.0.0.0/8", "::1/128", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7"}

func TestTrustProxiesIgnoresSpoofedHeadersFromTheInternet(t *testing.T) {
	ip, proto := proxyEcho(t, loopbackAndPrivate, "203.0.113.7:51000", map[string]string{
		"X-Forwarded-For":   "127.0.0.1",
		"X-Real-IP":         "10.1.2.3",
		"X-Forwarded-Proto": "https",
	})
	if ip != "203.0.113.7" {
		t.Errorf("ClientIP = %q, want the connecting address 203.0.113.7", ip)
	}
	if proto != "" {
		t.Errorf("X-Forwarded-Proto = %q, want it removed from an untrusted peer", proto)
	}
}

func TestTrustProxiesBelievesAConfiguredProxy(t *testing.T) {
	ip, proto := proxyEcho(t, loopbackAndPrivate, "172.18.0.5:40000", map[string]string{
		"X-Forwarded-For":   "198.51.100.20",
		"X-Forwarded-Proto": "https",
	})
	if ip != "198.51.100.20" {
		t.Errorf("ClientIP = %q, want the address the proxy forwarded", ip)
	}
	if proto != "https" {
		t.Errorf("X-Forwarded-Proto = %q, want https from a trusted proxy", proto)
	}
}

func TestTrustProxiesTakesTheRightmostUntrustedHop(t *testing.T) {
	// The client wrote the left-hand entry; the proxy appended the real peer.
	ip, _ := proxyEcho(t, loopbackAndPrivate, "127.0.0.1:40000", map[string]string{
		"X-Forwarded-For": "1.2.3.4, 198.51.100.20",
	})
	if ip != "198.51.100.20" {
		t.Errorf("ClientIP = %q, want 198.51.100.20", ip)
	}
}

func TestTrustProxiesNoneTrustsNobody(t *testing.T) {
	ip, proto := proxyEcho(t, []string{}, "127.0.0.1:40000", map[string]string{
		"X-Forwarded-For":   "198.51.100.20",
		"X-Forwarded-Proto": "https",
	})
	if ip != "127.0.0.1" || proto != "" {
		t.Errorf("ClientIP = %q, proto = %q, want 127.0.0.1 and no header", ip, proto)
	}
}

func TestTrustProxiesRejectsAMalformedEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := TrustProxies(r, []string{"10.0.0.0/8", "not-an-ip"}); err == nil {
		t.Fatal("TrustProxies accepted a malformed entry")
	}
	var ip string
	r.GET("/", func(c *gin.Context) { ip = c.ClientIP() })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.9:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.20")
	r.ServeHTTP(httptest.NewRecorder(), req)
	if ip != "10.0.0.9" {
		t.Errorf("ClientIP = %q after a bad TRUSTED_PROXIES, want the connecting address", ip)
	}
}
