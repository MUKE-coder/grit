package generate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SyncPolicy is the `sync:` block in a resource definition.
//
// It declares how the resource behaves offline, and it lives in the resource
// definition rather than in a separate file on purpose: one definition should
// describe server persistence, local persistence, outbox behaviour and
// conflict strategy. A second file is a second thing to keep in step, and
// keeping things in step is the entire product.
//
//	name: Sale
//	fields:
//	  - name: total
//	    type: float
//	sync:
//	  mode: offline_first
//	  conflict: server_wins
//	  fields: [total, payment_method]
//	  local_only: [draft_note]
//	  max_offline_age: 72h
type SyncPolicy struct {
	Mode      string   `yaml:"mode"`
	Conflict  string   `yaml:"conflict"`
	Fields    []string `yaml:"fields"`
	LocalOnly []string `yaml:"local_only"`
	// MaxOfflineAge is a Go duration string: "72h", "30m". Parsed rather than
	// taken as a number because "72" is ambiguous and somebody will read it as
	// hours while the code reads it as seconds.
	MaxOfflineAge string `yaml:"max_offline_age"`
}

// Valid values, kept as strings so the generator does not have to import the
// scaffolded project's sync package to validate a definition.
var (
	validSyncModes     = []string{"offline_first", "online_only"}
	validSyncConflicts = []string{"manual", "server_wins", "client_wins"}
)

// Validate checks a declared policy and fills in defaults.
//
// Run at parse time, so a typo is a CLI error naming the resource and the
// allowed values, not a wrong conflict decision on somebody's phone weeks
// later. The generator emits code the API panics on at boot for the same
// reason, but catching it here means the bad code is never written.
func (p *SyncPolicy) Validate(resource string) error {
	if p == nil {
		return nil
	}
	if p.Mode == "" {
		p.Mode = "offline_first"
	}
	if p.Conflict == "" {
		p.Conflict = "manual"
	}

	if !contains(validSyncModes, p.Mode) {
		return fmt.Errorf("sync.mode %q on %s: want one of %s",
			p.Mode, resource, strings.Join(validSyncModes, ", "))
	}
	if !contains(validSyncConflicts, p.Conflict) {
		return fmt.Errorf("sync.conflict %q on %s: want one of %s",
			p.Conflict, resource, strings.Join(validSyncConflicts, ", "))
	}

	// The client needs these to version a row, order it, and know it was
	// deleted. Letting somebody mark one local_only produces a mirror that
	// cannot sync and gives no hint why.
	for _, f := range p.LocalOnly {
		switch f {
		case "id", "version", "created_at", "updated_at", "deleted_at":
			return fmt.Errorf("sync.local_only on %s cannot include %q: the client needs it to version and order rows", resource, f)
		}
	}

	for _, f := range p.LocalOnly {
		if contains(p.Fields, f) {
			return fmt.Errorf("%q appears in both sync.fields and sync.local_only on %s", f, resource)
		}
	}

	if p.MaxOfflineAge != "" {
		d, err := time.ParseDuration(p.MaxOfflineAge)
		if err != nil {
			return fmt.Errorf("sync.max_offline_age %q on %s: %w (use a Go duration like 72h or 30m)",
				p.MaxOfflineAge, resource, err)
		}
		if d <= 0 {
			return fmt.Errorf("sync.max_offline_age on %s must be positive, got %q", resource, p.MaxOfflineAge)
		}
	}
	return nil
}

// CheckAgainstFields reports names in the policy that no field on the resource
// has.
//
// A typo here is silent in the worst way: sync.fields lists "totl", the
// allowlist excludes the real "total", and the client mirrors rows with the
// amount missing. Nothing errors anywhere.
func (p *SyncPolicy) CheckAgainstFields(resource string, fields []Field) []string {
	if p == nil {
		return nil
	}
	known := map[string]bool{
		"id": true, "version": true, "created_at": true,
		"updated_at": true, "deleted_at": true,
	}
	for _, f := range fields {
		known[toSnakeCase(f.Name)] = true
	}

	var unknown []string
	for _, name := range append(append([]string{}, p.Fields...), p.LocalOnly...) {
		if !known[name] {
			unknown = append(unknown, name)
		}
	}
	return unknown
}

// GoLiteral renders the policy as the sync.Policy literal the generator
// injects into routes.go. Returns "" for a policy that is entirely defaults,
// so an ordinary resource keeps the plain Register call it always had.
func (p *SyncPolicy) GoLiteral() string {
	if p == nil {
		return ""
	}
	isDefault := (p.Mode == "" || p.Mode == "offline_first") &&
		(p.Conflict == "" || p.Conflict == "manual") &&
		len(p.Fields) == 0 && len(p.LocalOnly) == 0 && p.MaxOfflineAge == ""
	if isDefault {
		return ""
	}

	var parts []string
	parts = append(parts, "Mode: sync."+goConst(p.Mode))
	parts = append(parts, "Conflict: sync."+goConst(p.Conflict))
	if len(p.Fields) > 0 {
		parts = append(parts, "Fields: "+goStringSlice(p.Fields))
	}
	if len(p.LocalOnly) > 0 {
		parts = append(parts, "LocalOnly: "+goStringSlice(p.LocalOnly))
	}
	if p.MaxOfflineAge != "" {
		// Already validated, so the error cannot fire; rendering the parsed
		// nanoseconds rather than the string keeps the emitted code free of a
		// second parse that could fail at runtime.
		if d, err := time.ParseDuration(p.MaxOfflineAge); err == nil {
			parts = append(parts, "MaxOfflineAge: "+strconv.FormatInt(int64(d), 10)+" /* "+p.MaxOfflineAge+" */")
		}
	}
	return "sync.Policy{" + strings.Join(parts, ", ") + "}"
}

// goConst maps a policy value to the exported Go constant that names it.
func goConst(value string) string {
	switch value {
	case "offline_first":
		return "ModeOfflineFirst"
	case "online_only":
		return "ModeOnlineOnly"
	case "manual":
		return "ConflictManual"
	case "server_wins":
		return "ConflictServerWins"
	case "client_wins":
		return "ConflictClientWins"
	}
	return "ConflictManual"
}

func goStringSlice(values []string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = strconv.Quote(v)
	}
	return "[]string{" + strings.Join(quoted, ", ") + "}"
}

func contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
