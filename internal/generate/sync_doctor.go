package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// SyncFinding is one thing wrong, or worth knowing, about a project's offline
// configuration.
type SyncFinding struct {
	Level   string // "error" | "warning" | "info"
	Model   string
	Message string
	Fix     string
}

// SyncReport is what `grit sync doctor` found.
type SyncReport struct {
	Findings   []SyncFinding
	Registered []string
	Models     []string
}

// Errors reports whether anything found will actually break.
func (r SyncReport) Errors() int {
	n := 0
	for _, f := range r.Findings {
		if f.Level == "error" {
			n++
		}
	}
	return n
}

// Warnings counts the findings worth reading but not fatal.
func (r SyncReport) Warnings() int {
	n := 0
	for _, f := range r.Findings {
		if f.Level == "warning" {
			n++
		}
	}
	return n
}

var (
	syncRegisterRe = regexp.MustCompile(
		`syncRegistry\.(Register|RegisterWithPolicy)\(\s*"([a-z0-9_]+)"\s*,\s*&models\.(\w+)\{\}(.*)`)
	syncPolicyFieldsRe    = regexp.MustCompile(`Fields:\s*\[\]string\{([^}]*)\}`)
	syncPolicyLocalOnlyRe = regexp.MustCompile(`LocalOnly:\s*\[\]string\{([^}]*)\}`)
	syncPolicyConflictRe  = regexp.MustCompile(`Conflict:\s*sync\.(\w+)`)
	syncPolicyModeRe      = regexp.MustCompile(`Mode:\s*sync\.(\w+)`)
	syncQuotedRe          = regexp.MustCompile(`"([^"]*)"`)
	goJSONTagRe           = regexp.MustCompile("json:\"([a-zA-Z0-9_]+)")
)

// SyncDoctor inspects a project's offline configuration and reports what will
// not work.
//
// The mistakes this catches are the ones that are silent. A field allowlist
// naming a column that does not exist does not error anywhere: the allowlist
// simply excludes the real column, and the client mirrors rows with the value
// missing. A model without a Version field cannot detect a conflict, so it
// silently takes whichever write landed last. Neither shows up in a build, a
// test, or a request log.
func SyncDoctor(root string) (*SyncReport, error) {
	report := &SyncReport{}

	routesPath := findSyncRoutesFile(root)
	if routesPath == "" {
		return nil, fmt.Errorf("no internal/routes/routes.go found: run this from a Grit project")
	}
	routes, err := os.ReadFile(routesPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", routesPath, err)
	}

	modelsDir := filepath.Join(filepath.Dir(filepath.Dir(routesPath)), "models")
	columns, structs := readModelColumns(modelsDir)
	for name := range structs {
		report.Models = append(report.Models, name)
	}
	sort.Strings(report.Models)

	for _, match := range syncRegisterRe.FindAllStringSubmatch(string(routes), -1) {
		table, model, tail := match[2], match[3], match[4]
		report.Registered = append(report.Registered, table)

		known, haveModel := columns[model]
		if !haveModel {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "warning",
				Model:   table,
				Message: fmt.Sprintf("registered as models.%s, which was not found in internal/models", model),
				Fix:     "Check the model name, or ignore this if the model lives elsewhere.",
			})
			continue
		}

		// Without a Version column there is nothing to compare, so every push
		// is accepted and the last writer wins with no way to know it happened.
		if !known["version"] {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "error",
				Model:   table,
				Message: "has no Version field, so conflicts cannot be detected at all",
				Fix:     "Add `Version int` to the model and a BeforeUpdate hook that increments it.",
			})
		}
		if !known["id"] {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "error",
				Model:   table,
				Message: "has no id field, so the client cannot key its mirror",
				Fix:     "Give the model a string (UUID) primary key.",
			})
		}
		if !known["updated_at"] {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "error",
				Model:   table,
				Message: "has no UpdatedAt, so cursor pulls cannot tell what changed",
				Fix:     "Add `UpdatedAt time.Time` to the model.",
			})
		}
		if !known["deleted_at"] {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "warning",
				Model:   table,
				Message: "has no soft delete, so deletes never reach offline clients as tombstones",
				Fix:     "Add `DeletedAt gorm.DeletedAt` so a delete becomes a tombstone rather than a row that silently stays in every mirror.",
			})
		}

		report.Findings = append(report.Findings, checkPolicyLiteral(table, tail, known)...)
	}

	// A policy nothing enforces is worse than no policy: routes.go says
	// server_wins, the handler keeps asking users about conflicts, and the
	// declaration reads as though it took effect.
	if strings.Contains(string(routes), "RegisterWithPolicy") {
		handler := filepath.Join(filepath.Dir(filepath.Dir(routesPath)), "handlers", "sync.go")
		if body, err := os.ReadFile(handler); err == nil && !strings.Contains(string(body), "PolicyFor") {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "error",
				Message: "policies are declared but internal/handlers/sync.go does not enforce them",
				Fix:     "That handler predates sync policies. Replace it with the current template so conflict strategy, field allowlists and local_only actually apply.",
			})
		}
		// Without this route a client cannot discover the policy, so it falls
		// back to defaults and renders a manual conflict prompt for a resource
		// declared server_wins.
		if !strings.Contains(string(routes), "syncHandler.Policy") {
			report.Findings = append(report.Findings, SyncFinding{
				Level:   "warning",
				Message: "GET /api/sync/policy is not mounted, so clients cannot read the declared policies",
				Fix:     "Add protected.GET(\"/sync/policy\", syncHandler.Policy) beside the existing sync routes.",
			})
		}
	}

	sort.Strings(report.Registered)
	if len(report.Registered) == 0 {
		report.Findings = append(report.Findings, SyncFinding{
			Level:   "warning",
			Message: "no models are registered for sync",
			Fix:     "Generate a resource, or add syncRegistry.Register(...) to routes.go.",
		})
	}
	return report, nil
}

// checkPolicyLiteral validates the sync.Policy literal on one registration
// against the model's real columns.
func checkPolicyLiteral(table, literal string, known map[string]bool) []SyncFinding {
	var findings []SyncFinding

	for _, m := range syncPolicyFieldsRe.FindAllStringSubmatch(literal, -1) {
		for _, field := range quotedStrings(m[1]) {
			if !known[field] {
				findings = append(findings, SyncFinding{
					Level:   "error",
					Model:   table,
					Message: fmt.Sprintf("sync fields allowlist names %q, which the model does not have", field),
					Fix:     "Fix the name. An allowlist entry that matches nothing silently drops the real column from every mirror.",
				})
			}
		}
	}

	for _, m := range syncPolicyLocalOnlyRe.FindAllStringSubmatch(literal, -1) {
		for _, field := range quotedStrings(m[1]) {
			if !known[field] {
				findings = append(findings, SyncFinding{
					Level:   "warning",
					Model:   table,
					Message: fmt.Sprintf("local_only names %q, which the model does not have", field),
					Fix:     "Harmless, but it suggests a rename the policy did not follow.",
				})
			}
		}
	}

	if m := syncPolicyConflictRe.FindStringSubmatch(literal); m != nil {
		switch m[1] {
		case "ConflictClientWins":
			findings = append(findings, SyncFinding{
				Level:   "info",
				Model:   table,
				Message: "client_wins: an offline edit overwrites whatever the server has, unread",
				Fix:     "Correct when one author owns a record. Wrong for anything two people touch.",
			})
		case "ConflictServerWins":
			findings = append(findings, SyncFinding{
				Level:   "info",
				Model:   table,
				Message: "server_wins: an offline edit is discarded if the server moved on",
				Fix:     "The user is not asked. Make sure the screen says so before they type.",
			})
		}
	}

	if m := syncPolicyModeRe.FindStringSubmatch(literal); m != nil && m[1] == "ModeOnlineOnly" {
		findings = append(findings, SyncFinding{
			Level:   "info",
			Model:   table,
			Message: "online_only: not mirrored, and a pull for it is rejected",
			Fix:     "Screens reading it will be empty offline. Intended, but worth confirming.",
		})
	}

	return findings
}

// readModelColumns maps each model struct to the JSON column names it carries.
//
// Reading the json tags rather than the Go field names, because the json name
// is what crosses the wire and therefore what a policy allowlist has to match.
func readModelColumns(dir string) (map[string]map[string]bool, map[string]bool) {
	columns := map[string]map[string]bool{}
	structs := map[string]bool{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return columns, structs
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		// parseGoStructs takes a path and parses the file itself. Handing it
		// the contents made it parse a filename-shaped string, find nothing,
		// and report every model as missing.
		parsed, parseErr := parseGoStructs(path)
		if parseErr != nil {
			continue
		}
		for _, s := range parsed {
			structs[s.Name] = true
			cols := map[string]bool{}
			for _, f := range s.Fields {
				name := f.JSONName
				if name == "" {
					name = toSnakeCase(f.Name)
				}
				cols[name] = true
			}
			// Embedded gorm.Model and gorm.DeletedAt do not surface as named
			// fields, so pick the columns off the raw text as well.
			for _, m := range goJSONTagRe.FindAllStringSubmatch(string(data), -1) {
				cols[m[1]] = true
			}
			if strings.Contains(string(data), "gorm.DeletedAt") {
				cols["deleted_at"] = true
			}
			columns[s.Name] = cols
		}
	}
	return columns, structs
}

func quotedStrings(list string) []string {
	var out []string
	for _, m := range syncQuotedRe.FindAllStringSubmatch(list, -1) {
		if m[1] != "" {
			out = append(out, m[1])
		}
	}
	return out
}

func findSyncRoutesFile(root string) string {
	for _, candidate := range []string{
		filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"),
		filepath.Join(root, "internal", "routes", "routes.go"),
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}
