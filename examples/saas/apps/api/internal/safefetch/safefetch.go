// Package safefetch performs HTTP requests against user-supplied URLs
// with SSRF defences. Use this any time you fetch a URL the caller
// chose — webhook delivery, "fetch image from URL", PDF render from a
// URL, OEmbed expansion, etc.
//
// Coverage: OWASP Top 10:2025 A01 (SSRF was absorbed into Broken Access
// Control in 2025). The classic SSRF impact — proxying requests to the
// cloud metadata service (169.254.169.254) to steal IAM credentials —
// is blocked at TCP-connect time even if DNS resolution races.
package safefetch

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"
)

// Client is the default safefetch HTTP client. Timeouts are short on
// purpose — long-running fetches are the easiest amplification vector.
var Client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
			Control:   controlSafeDial,
		}).DialContext,
		MaxIdleConns:        20,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	},
	// Validate every redirect target — attackers love to 302 you into
	// 169.254.169.254 after a clean validation pass.
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("too many redirects")
		}
		return validateURL(req.URL)
	},
}

// ErrBlocked is returned when a target is rejected by the SSRF guard.
var ErrBlocked = errors.New("safefetch: target blocked")

// Get fetches the URL with SSRF defences. Returns ErrBlocked (wrapped)
// for any disallowed target.
func Get(ctx context.Context, rawURL string) (*http.Response, error) {
	req, err := newRequest(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	// #nosec G704 -- the destination was resolved and checked against the
	// private and metadata ranges above. Refusing an internal target is
	// what this package is for.
	return Client.Do(req)
}

// Do runs a pre-built request through the safe client. Use this when
// you need POST/PUT or custom headers.
func Do(req *http.Request) (*http.Response, error) {
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}
	// #nosec G704 -- the destination was resolved and checked against the
	// private and metadata ranges above. Refusing an internal target is
	// what this package is for.
	return Client.Do(req)
}

func newRequest(ctx context.Context, method, rawURL string, body interface {
	Read(p []byte) (n int, err error)
}) (*http.Request, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("safefetch: parsing url: %w", err)
	}
	if err := validateURL(u); err != nil {
		return nil, err
	}
	if body == nil {
		return http.NewRequestWithContext(ctx, method, u.String(), nil)
	}
	return http.NewRequestWithContext(ctx, method, u.String(), body)
}

// validateURL enforces the scheme allowlist and rejects hostnames that
// are obviously private. The IP-level check happens again at dial time
// via controlSafeDial.
func validateURL(u *url.URL) error {
	if u == nil {
		return fmt.Errorf("%w: nil url", ErrBlocked)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		// ok
	default:
		return fmt.Errorf("%w: scheme %q not allowed", ErrBlocked, u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("%w: empty host", ErrBlocked)
	}
	// Cloud metadata hostnames — rejected by name in case DNS hasn't
	// been consulted yet (e.g. proxy mode).
	lower := strings.ToLower(host)
	if lower == "metadata.google.internal" || lower == "metadata" || lower == "instance-data" {
		return fmt.Errorf("%w: cloud metadata host", ErrBlocked)
	}
	// If the host parses as a literal IP, check it now. (Hostnames are
	// re-checked at dial time after resolution.)
	if ip := net.ParseIP(host); ip != nil {
		if isPrivateOrLoopback(ip) {
			return fmt.Errorf("%w: literal IP %s in private range", ErrBlocked, ip)
		}
	}
	return nil
}

// controlSafeDial runs after DNS resolution, immediately before connect.
// This is where we close the DNS-rebinding hole — even if a hostname
// validated cleanly, the resolved IP at connect time may not have.
func controlSafeDial(network, address string, _ syscall.RawConn) error {
	if !strings.HasPrefix(network, "tcp") {
		return fmt.Errorf("%w: non-tcp network %q", ErrBlocked, network)
	}
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("%w: splitting host:port: %v", ErrBlocked, err)
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("%w: dial address %q is not an IP", ErrBlocked, host)
	}
	if isPrivateOrLoopback(ip) {
		return fmt.Errorf("%w: resolved IP %s in private range", ErrBlocked, ip)
	}
	return nil
}

// isPrivateOrLoopback returns true for any IP an SSRF guard must reject:
// loopback, link-local, multicast, unspecified, RFC1918 private,
// carrier-grade NAT, and the IMDS endpoints (169.254.169.254 + the IPv6
// AWS metadata fd00:ec2::254).
func isPrivateOrLoopback(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified() || ip.IsInterfaceLocalMulticast() {
		return true
	}
	// Go 1.17+ exposes IsPrivate for RFC1918 / RFC4193.
	if ip.IsPrivate() {
		return true
	}
	// Carrier-grade NAT (RFC 6598).
	_, cgnat, _ := net.ParseCIDR("100.64.0.0/10")
	if cgnat.Contains(ip) {
		return true
	}
	// AWS IPv6 metadata endpoint.
	_, awsV6, _ := net.ParseCIDR("fd00:ec2::/32")
	if awsV6.Contains(ip) {
		return true
	}
	return false
}
