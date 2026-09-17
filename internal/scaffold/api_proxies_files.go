package scaffold

// Contact-app review L6. gin.New() believes X-Forwarded-For and X-Real-IP from
// any peer until SetTrustedProxies is called, and nothing called it, so any
// client chose the address its session rows, audit entries and rate limits were
// keyed on. X-Forwarded-Proto was read the same way. middleware/proxies.go
// narrows all three headers to the proxies TRUSTED_PROXIES names.

// middlewareProxiesGo emits internal/middleware/proxies.go.
func middlewareProxiesGo() string {
	return `package middleware

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/gin-gonic/gin"
)

// forwardedHeaders are the request headers a reverse proxy writes. Only a proxy
// in TRUSTED_PROXIES gets to set them: from anyone else they are the client's
// own claim about itself.
var forwardedHeaders = []string{"X-Forwarded-For", "X-Real-IP", "X-Forwarded-Proto", "X-Forwarded-Host"}

// TrustProxies makes the engine believe forwarded headers only from the given
// proxies (IPs or CIDRs, config.TrustedProxies).
//
// c.ClientIP() then returns the connecting address unless that address is a
// trusted proxy, and X-Forwarded-Proto is removed from requests that did not
// come through one, so code that reads it (the Secure cookie flag, HSTS) cannot
// be told a plain HTTP request is HTTPS. Register it before any other
// middleware.
//
// An entry that is neither an IP nor a CIDR is an error, and the engine is left
// trusting no proxy at all: the safe direction, where the worst case is that
// requests carry the proxy's address.
func TrustProxies(r *gin.Engine, proxies []string) error {
	prefixes, err := parseProxies(proxies)
	if err != nil {
		if resetErr := r.SetTrustedProxies(nil); resetErr != nil {
			return fmt.Errorf("TRUSTED_PROXIES: %w (and clearing it: %v)", err, resetErr)
		}
		r.Use(stripForwardedHeaders(nil))
		return fmt.Errorf("TRUSTED_PROXIES: %w, so no proxy is trusted", err)
	}
	if err := r.SetTrustedProxies(proxies); err != nil {
		return fmt.Errorf("TRUSTED_PROXIES: %w", err)
	}
	r.Use(stripForwardedHeaders(prefixes))
	return nil
}

func parseProxies(proxies []string) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, 0, len(proxies))
	for _, raw := range proxies {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		if p, err := netip.ParsePrefix(entry); err == nil {
			prefixes = append(prefixes, p.Masked())
			continue
		}
		addr, err := netip.ParseAddr(entry)
		if err != nil {
			return nil, fmt.Errorf("%q is neither an IP address nor a CIDR", entry)
		}
		prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
	}
	return prefixes, nil
}

// stripForwardedHeaders removes the forwarded headers from a request whose
// connecting address is not a trusted proxy.
func stripForwardedHeaders(trusted []netip.Prefix) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !peerTrusted(c.RemoteIP(), trusted) {
			for _, h := range forwardedHeaders {
				c.Request.Header.Del(h)
			}
		}
		c.Next()
	}
}

func peerTrusted(ip string, trusted []netip.Prefix) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	// A dual-stack listener reports an IPv4 peer as ::ffff:a.b.c.d, which an
	// IPv4 prefix never contains.
	addr = addr.WithZone("").Unmap()
	for _, p := range trusted {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}
`
}

// middlewareProxiesTestGo emits internal/middleware/proxies_test.go.
func middlewareProxiesTestGo() string {
	return `package middleware

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
`
}
