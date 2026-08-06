package scaffold

// v3.136 — internationalisation for the API surface.
//
// Scope, stated up front because a half-translated app is worse than an
// honestly monolingual one: this covers what the *API* says. Error codes,
// validation failures, auth responses. The strings a client cannot translate
// for itself because it never sees the reason, only the message.
//
// It does not translate the generated admin's 31,000 lines of chrome. That is a
// separate and much wider change, and claiming it here would be a lie the first
// time somebody switched locale and saw an English table header.
//
// Why not gettext: .po files need a compile step and a toolchain, and the whole
// point of the scaffold is that things work on first run. JSON catalogues embed
// into the binary, need no tooling, and are editable by a translator with a text
// editor.

// i18nGo emits internal/i18n/i18n.go.
func i18nGo() string {
	return `package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed locales/*.json
var catalogues embed.FS

// Default is the locale used when the request asks for nothing, asks for
// something unknown, or asks for a key this build has no translation for.
const Default = "en"

var (
	once     sync.Once
	messages = map[string]map[string]string{}
	loadErr  error
)

// Load reads every embedded catalogue. Called once, lazily, so a handler that
// asks for a string before main() has finished wiring still gets one.
func Load() error {
	once.Do(func() {
		entries, err := catalogues.ReadDir("locales")
		if err != nil {
			loadErr = fmt.Errorf("i18n: reading locales: %w", err)
			return
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			raw, err := catalogues.ReadFile("locales/" + e.Name())
			if err != nil {
				loadErr = fmt.Errorf("i18n: reading %s: %w", e.Name(), err)
				return
			}
			var flat map[string]string
			if err := json.Unmarshal(raw, &flat); err != nil {
				loadErr = fmt.Errorf("i18n: parsing %s: %w", e.Name(), err)
				return
			}
			messages[strings.TrimSuffix(e.Name(), ".json")] = flat
		}
	})
	return loadErr
}

// Available lists the locales this binary was built with, so the frontend can
// render a switcher containing only what actually resolves.
func Available() []string {
	_ = Load()
	out := make([]string, 0, len(messages))
	for k := range messages {
		out = append(out, k)
	}
	return out
}

// Has reports whether a locale was compiled in.
func Has(locale string) bool {
	_ = Load()
	_, ok := messages[locale]
	return ok
}

// T translates key into locale, substituting fmt-style arguments.
//
// Three fallbacks, in order: the requested locale, the default locale, then the
// key itself. Returning the key rather than an empty string matters: a missing
// translation shows up in the UI as "errors.not_found" and gets fixed, where an
// empty string shows up as a blank space and does not.
func T(locale, key string, args ...any) string {
	_ = Load()

	if m, ok := messages[locale]; ok {
		if s, ok := m[key]; ok {
			return format(s, args...)
		}
	}
	if locale != Default {
		if m, ok := messages[Default]; ok {
			if s, ok := m[key]; ok {
				return format(s, args...)
			}
		}
	}
	return key
}

func format(s string, args ...any) string {
	if len(args) == 0 {
		return s
	}
	return fmt.Sprintf(s, args...)
}

// Negotiate picks the best locale from an Accept-Language header.
//
// Deliberately simple: it honours quality weights and falls back from a region
// to its base language, so "fr-CA" resolves to "fr" when only "fr" is compiled
// in. It does not implement the whole of RFC 4647, because the alternative to
// forty lines here is a dependency, and this is the kind of thing that should
// not need one.
func Negotiate(header string) string {
	_ = Load()
	type pick struct {
		tag string
		q   float64
	}
	var picks []pick

	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		tag, q := part, 1.0
		if i := strings.Index(part, ";"); i != -1 {
			tag = strings.TrimSpace(part[:i])
			if _, err := fmt.Sscanf(strings.TrimSpace(part[i+1:]), "q=%f", &q); err != nil {
				q = 1.0
			}
		}
		picks = append(picks, pick{strings.ToLower(tag), q})
	}

	// Highest quality first, stable within equal weights so the header's own
	// order decides ties, which is what a browser means by it.
	for i := 1; i < len(picks); i++ {
		for j := i; j > 0 && picks[j].q > picks[j-1].q; j-- {
			picks[j], picks[j-1] = picks[j-1], picks[j]
		}
	}

	for _, p := range picks {
		if p.tag == "*" {
			return Default
		}
		if Has(p.tag) {
			return p.tag
		}
		if base, _, found := strings.Cut(p.tag, "-"); found && Has(base) {
			return base
		}
	}
	return Default
}
`
}

// i18nMiddlewareGo emits internal/middleware/locale.go.
func i18nMiddlewareGo() string {
	return `package middleware

import (
	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/i18n"
)

// LocaleKey is where the resolved locale lives on the request context.
const LocaleKey = "locale"

// LocaleCookie is the cookie the frontend stores the chosen language in.
//
// It has to be the same name next-intl reads on the Next.js side. The admin
// picks a language, writes this cookie, and every subsequent API call carries
// it automatically, so the page and its error messages agree. Get this wrong
// and you get a French UI reporting English validation failures, which is the
// exact failure mode partial i18n is famous for.
const LocaleCookie = "grit_locale"

// Locale resolves the language for this request, most explicit signal first:
//
//  1. ?lang= on the query string. An override for testing and deep links.
//  2. X-Locale, for native clients with no cookie jar.
//  3. The grit_locale cookie, which is what the web and admin actually use.
//     Cookie rather than a URL prefix is a deliberate choice: routes stay the
//     same in every language, so a link shared between two people works for
//     both of them and analytics does not fragment by locale.
//  4. Accept-Language, which is what a browser sends before anyone has chosen.
//
// A user's stored database preference is deliberately not consulted here. This
// runs before auth, so there is no user yet; handlers that want to prefer a
// saved setting can override the context value after auth has run.
func Locale() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := ""

		if q := c.Query("lang"); q != "" && i18n.Has(q) {
			locale = q
		}
		if locale == "" {
			if h := c.GetHeader("X-Locale"); h != "" && i18n.Has(h) {
				locale = h
			}
		}
		if locale == "" {
			if ck, err := c.Cookie(LocaleCookie); err == nil && i18n.Has(ck) {
				locale = ck
			}
		}
		if locale == "" {
			locale = i18n.Negotiate(c.GetHeader("Accept-Language"))
		}

		c.Set(LocaleKey, locale)
		// Caches and CDNs key on this. Without it, the first visitor's language
		// is served to everyone who follows.
		c.Writer.Header().Add("Vary", "Accept-Language")
		c.Header("Content-Language", locale)
		c.Next()
	}
}

// LocaleOf returns the locale resolved for this request.
func LocaleOf(c *gin.Context) string {
	if v, ok := c.Get(LocaleKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return i18n.Default
}
`
}

// i18nResponseGo emits internal/response/response.go — the helpers handlers use
// so a translated error is the path of least resistance rather than extra work.
func i18nResponseGo() string {
	return `package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/i18n"
	"{{MODULE}}/internal/middleware"
)

// T translates a key in the request's locale.
func T(c *gin.Context, key string, args ...any) string {
	return i18n.T(middleware.LocaleOf(c), key, args...)
}

// Error writes the standard error envelope with a translated message.
//
//	response.Error(c, http.StatusNotFound, "NOT_FOUND", "errors.not_found")
//
// The code stays in English on purpose. Clients switch on it, humans read the
// message, and translating the thing code branches on breaks every client that
// ever shipped.
func Error(c *gin.Context, status int, code, key string, args ...any) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": T(c, key, args...),
		},
	})
}

// ErrorWithDetails is Error plus a per-field map, for validation failures.
// Field names are keys a client matches on, so they stay untranslated; the
// values are messages a person reads, so they do not.
func ErrorWithDetails(c *gin.Context, status int, code, key string, details map[string]string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": T(c, key),
			"details": details,
		},
	})
}

// NotFound, Unauthorized, Forbidden, Validation and Internal cover the five
// shapes the generated handlers actually emit.
func NotFound(c *gin.Context, key string, args ...any) {
	Error(c, http.StatusNotFound, "NOT_FOUND", key, args...)
}

func Unauthorized(c *gin.Context, key string, args ...any) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", key, args...)
}

func Forbidden(c *gin.Context, key string, args ...any) {
	Error(c, http.StatusForbidden, "FORBIDDEN", key, args...)
}

func Validation(c *gin.Context, key string, args ...any) {
	Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", key, args...)
}

func Internal(c *gin.Context, key string, args ...any) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", key, args...)
}

// OK and Created carry a translated message beside the payload.
func OK(c *gin.Context, data any, key string, args ...any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "message": T(c, key, args...)})
}

func Created(c *gin.Context, data any, key string, args ...any) {
	c.JSON(http.StatusCreated, gin.H{"data": data, "message": T(c, key, args...)})
}
`
}
