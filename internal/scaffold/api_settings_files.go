package scaffold

// The settings registry.
//
// Every application ends up needing values that are configuration but not
// environment variables: company name, invoice prefix, default currency,
// whether to email on a new order. Today those have three bad homes. In .env,
// which needs a deploy to change and which no admin can touch. Hardcoded,
// which needs a code change. Or a hand-rolled settings table with a
// hand-rolled admin page, written again in every project.
//
// A setting is declared next to the code that reads it, and everything else
// (storage, the admin page, validation, caching, the resolution order) follows
// from the declaration.

// apiSettingsRegistryGo emits internal/settings/settings.go.
func apiSettingsRegistryGo() string {
	return `// Package settings is the typed application settings registry.
//
// Declare next to the code that uses it:
//
//	settings.Define(settings.Setting{
//	    Key: "invoice.prefix", Type: settings.TypeString, Default: "INV-",
//	    Label: "Invoice number prefix", Group: "Billing",
//	})
//
// and read it anywhere:
//
//	prefix := settings.String(ctx, "invoice.prefix")
package settings

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Kind is a setting's type. It decides the control the admin renders and how
// the stored string is parsed on the way out.
type Kind string

const (
	TypeString Kind = "string"
	TypeText   Kind = "text"
	TypeNumber Kind = "number"
	TypeBool   Kind = "bool"
	TypeSelect Kind = "select"
	TypeColour Kind = "colour"
	// Secret is a string the API never sends back once set. An SMTP password
	// belongs here: an admin can replace it and cannot read it, which is the
	// behaviour people expect and almost never get.
	TypeSecret Kind = "secret"
)

// Option is one choice for a TypeSelect setting.
type Option struct {
	Value string ` + "`" + `json:"value"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
}

// Setting is one declared value.
type Setting struct {
	Key   string ` + "`" + `json:"key"` + "`" + `
	Type  Kind   ` + "`" + `json:"type"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
	Help  string ` + "`" + `json:"help,omitempty"` + "`" + `
	// Group is the admin section this appears under: "Billing", "Email".
	Group string ` + "`" + `json:"group"` + "`" + `
	// Default is used when nothing is stored and no environment variable is
	// set. Always a string; typed accessors parse it.
	Default string   ` + "`" + `json:"default"` + "`" + `
	Options []Option ` + "`" + `json:"options,omitempty"` + "`" + `
	// Env names an environment variable that overrides the stored value.
	// Empty means the key's own SCREAMING_SNAKE form is used.
	Env string ` + "`" + `json:"-"` + "`" + `
	// Scope says how far a value can be overridden below the global one.
	Scope Scope ` + "`" + `json:"scope"` + "`" + `
	// Permission gates editing. Empty means any admin.
	Permission string ` + "`" + `json:"permission,omitempty"` + "`" + `
	// Validate rejects a bad value before it is stored. Nil means anything
	// the type parses.
	Validate func(string) error ` + "`" + `json:"-"` + "`" + `
	// Order sorts a setting within its group. Equal values fall back to Label.
	Order int ` + "`" + `json:"order"` + "`" + `
}

// Scope is how far down a setting may be overridden.
type Scope string

const (
	// Global is one value for the whole installation.
	Global Scope = "global"
	// Tenant allows a per-organisation override, for the multitenant plugin.
	Tenant Scope = "tenant"
	// User allows a per-person override: locale, density, notification
	// preferences.
	User Scope = "user"
)

var (
	mu       sync.RWMutex
	registry = map[string]Setting{}
)

// Define registers a setting.
//
// Panics on a bad declaration, at boot, because a setting nothing can store
// or render is a bug in the code rather than a runtime condition, and finding
// out at startup is much cheaper than finding out from an admin page that
// will not save.
func Define(s Setting) {
	if s.Key == "" {
		panic("settings: a setting needs a key")
	}
	if s.Type == "" {
		s.Type = TypeString
	}
	if s.Label == "" {
		s.Label = humanise(s.Key)
	}
	if s.Group == "" {
		s.Group = "General"
	}
	if s.Scope == "" {
		s.Scope = Global
	}
	if s.Env == "" {
		s.Env = envName(s.Key)
	}
	if s.Type == TypeSelect && len(s.Options) == 0 {
		panic("settings: " + s.Key + " is a select with no options")
	}
	if s.Default != "" && s.Validate != nil {
		if err := s.Validate(s.Default); err != nil {
			panic("settings: the default for " + s.Key + " fails its own validation: " + err.Error())
		}
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[s.Key]; exists {
		panic("settings: " + s.Key + " is defined twice")
	}
	registry[s.Key] = s
}

// Get returns a declaration.
func Get(key string) (Setting, bool) {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := registry[key]
	return s, ok
}

// All returns every declaration, sorted by group then order then label, which
// is the order the admin page renders them in.
func All() []Setting {
	mu.RLock()
	out := make([]Setting, 0, len(registry))
	for _, s := range registry {
		out = append(out, s)
	}
	mu.RUnlock()

	sort.Slice(out, func(i, j int) bool {
		if out[i].Group != out[j].Group {
			return out[i].Group < out[j].Group
		}
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].Label < out[j].Label
	})
	return out
}

// Groups returns the group names in render order.
func Groups() []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range All() {
		if !seen[s.Group] {
			seen[s.Group] = true
			out = append(out, s.Group)
		}
	}
	return out
}

// EnvOverride returns the environment value for a setting, if one is set.
//
// The environment wins over the database on purpose. A container that sets
// SMTP_HOST should not be quietly overridden by a value an admin typed months
// ago, and an operator needs a way to force a value that does not involve the
// application booting far enough to serve an admin page.
func (s Setting) EnvOverride() (string, bool) {
	if s.Env == "" {
		return "", false
	}
	v, ok := os.LookupEnv(s.Env)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

// Parse checks a value against the setting's type and validator.
func (s Setting) Parse(value string) error {
	switch s.Type {
	case TypeNumber:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("%s must be a number", s.Label)
		}
	case TypeBool:
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("%s must be true or false", s.Label)
		}
	case TypeSelect:
		found := false
		for _, o := range s.Options {
			if o.Value == value {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s must be one of the offered choices", s.Label)
		}
	case TypeColour:
		if !strings.HasPrefix(value, "#") || (len(value) != 4 && len(value) != 7) {
			return fmt.Errorf("%s must be a hex colour like #0a0a0f", s.Label)
		}
	}
	if s.Validate != nil {
		return s.Validate(value)
	}
	return nil
}

// ── Validators ───────────────────────────────────────────────────────────

// MaxLen rejects a value longer than n.
func MaxLen(n int) func(string) error {
	return func(v string) error {
		if len(v) > n {
			return fmt.Errorf("must be %d characters or fewer", n)
		}
		return nil
	}
}

// NotEmpty rejects a blank value, for a setting the app cannot run without.
func NotEmpty() func(string) error {
	return func(v string) error {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("cannot be empty")
		}
		return nil
	}
}

// Between rejects a number outside a range.
func Between(lo, hi float64) func(string) error {
	return func(v string) error {
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("must be a number")
		}
		if n < lo || n > hi {
			return fmt.Errorf("must be between %g and %g", lo, hi)
		}
		return nil
	}
}

// envName turns "invoice.prefix" into "INVOICE_PREFIX".
func envName(key string) string {
	r := strings.NewReplacer(".", "_", "-", "_")
	return strings.ToUpper(r.Replace(key))
}

// humanise turns "invoice.prefix" into "Invoice prefix".
func humanise(key string) string {
	parts := strings.FieldsFunc(key, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})
	for i, p := range parts {
		if i == 0 && p != "" {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
`
}
