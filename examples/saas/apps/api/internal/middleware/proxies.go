package middleware

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
