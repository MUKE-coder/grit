package middleware

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"saas/apps/api/internal/respond"
)

// RequestID injects a unique X-Request-ID header into every request and
// stores it in the context for downstream logging and tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// #nosec G404 -- a request id has to be unique for tracing, not
			// unpredictable. crypto/rand here would be a syscall on every
			// request for no security gain.
			requestID = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// SecurityHeaders adds production security headers to every response.
//
// Coverage against OWASP Top 10:2025 — A02 Security Misconfiguration,
// A05 Injection (XSS hardening via CSP), A04 Cryptographic Failures
// (HSTS forces TLS), plus mitigations for clickjacking, MIME sniffing,
// referrer leakage, and Spectre-class cross-origin attacks.
//
// CSP is deliberately strict-by-default. The scaffold's SPA serves /api
// from the same origin, so 'self' covers the normal case. Customise
// CSPDirectives via config when adding a CDN / inline scripts.
func SecurityHeaders() gin.HandlerFunc {
	// What script-src allows, which depends on what this binary serves.
	//
	// An API serving JSON needs no inline script at all, and the strict value
	// is right for it. A single project serves its own frontend from this same
	// binary, and Next.js delivers its bootstrap as an inline <script>: under
	// script-src 'self' the browser blocks it, React never starts, and what you
	// get is a page that renders perfectly and does nothing. A form submits
	// natively, which puts the password in the URL.
	//
	// The frontend's own config has carried this caveat in a comment for a long
	// time; it was the Go side that did not know it had become a frontend host.
	scriptSrc := "script-src 'self'; "

	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		// Spectre-class defence: isolate this origin from cross-origin reads
		// and require explicit opt-in for cross-origin embedders.
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")
		// Content-Security-Policy — strict default, blocks inline script
		// (XSS A05 hardening). Skip on /docs and /studio which serve
		// vendored UIs that rely on inline styles.
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/docs") && !strings.HasPrefix(path, "/studio") && !strings.HasPrefix(path, "/sentinel") && !strings.HasPrefix(path, "/pulse") {
			c.Header("Content-Security-Policy",
				"default-src 'self'; "+
					scriptSrc+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data: blob: https:; "+
					"font-src 'self' data:; "+
					"connect-src 'self'; "+
					"frame-ancestors 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'; "+
					"object-src 'none'")
		}
		// HSTS only when actually on HTTPS (don't break dev on http://).
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}
		c.Next()
	}
}

// MaxBodySize limits the request body to prevent abuse.
func MaxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > limit {
			respond.Fail(c, respond.CodePayloadTooLarge, fmt.Sprintf("Request body exceeds %dMB limit", limit/(1024*1024)))
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

// Logger creates a structured logging middleware with request ID correlation.
// Silently skips internal dashboard paths to keep the terminal readable.
func Logger() gin.HandlerFunc {
	// Paths that generate noise and aren't useful to see in dev logs
	skipPrefixes := []string{
		"/studio/",
		"/pulse/",
		"/pulse",
		"/sentinel/",
		"/docs/",
		"/docs",
		"/api/health",
		"/favicon.ico",
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Skip noisy internal paths
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(path, prefix) || path == prefix {
				c.Next()
				return
			}
		}

		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		requestID, _ := c.Get("request_id")

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%d] %s %s | %s | %v | id=%v",
			status,
			method,
			path,
			clientIP,
			latency,
			requestID,
		)
	}
}
