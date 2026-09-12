// Package doctor audits a Grit project for the mistakes that do not announce
// themselves.
//
// Every check here exists because the mistake it finds was made once, in a real
// project, and cost something to find. An encrypted column with no key is a
// plain column, and nothing says so. A resource whose list is not scoped to its
// owner hands every row to every signed-in user, and the only symptom is that
// it works. A tenant-owned model that forgot tenant.Owned is a table shared
// across organizations. None of that fails a build, a test, or a request log,
// which is why it needs a command.
//
// The rule for a check: it must key on something the generator actually writes,
// and it must stay quiet on a project that is right. A linter that cries wolf is
// a linter people turn off.
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Finding is one thing wrong, or worth knowing, about a project.
type Finding struct {
	Level    string // "error" | "warning"
	Check    string // the check's name, so a finding can be talked about
	Resource string // what it is about: a resource, a file, a variable
	Message  string
	Fix      string
}

// Report is what `grit doctor` found.
type Report struct {
	Findings  []Finding
	Resources []string // the generated resources it inspected
	Checks    int
}

// Errors counts the findings that are a problem now.
func (r Report) Errors() int { return r.count("error") }

// Warnings counts the findings worth reading but not fatal.
func (r Report) Warnings() int { return r.count("warning") }

func (r Report) count(level string) int {
	n := 0
	for _, f := range r.Findings {
		if f.Level == level {
			n++
		}
	}
	return n
}

// checks is every check, in the order they are reported.
var checks = []struct {
	name string
	run  func(*project) []Finding
}{
	{"encryption-key-unset", checkEncryptionKey},
	{"encrypted-column-in-whitelist", checkEncryptedWhitelists},
	{"owned-resource-unscoped", checkOwnedScoping},
	{"owner-settable-from-body", checkOwnerFromBody},
	{"resource-could-be-owned", checkUnownedUserResource},
	{"tenant-shared-resource", checkTenantOwned},
	{"tenant-middleware-missing", checkTenantMiddleware},
	{"database-not-durable", checkDatabaseProvider},
	{"tenant-sso-no-organization", checkTenantSSOProvisioning},
	{"pii-column-not-encrypted", checkPII},
	{"append-only-mutable-routes", checkAppendOnlyRoutes},
	{"studio-unprotected", checkStudio},
	{"default-credentials", checkCredentials},
	{"framework-library-behind", checkFrameworkDeps},
	{"rate-limits-per-process", checkSentinelCounters},
	{"public-allowlist-sensitive", checkPublicAllowlist},
}

// Run audits the project at root.
func Run(root string) (*Report, error) {
	p, err := loadProject(root)
	if err != nil {
		return nil, err
	}
	report := &Report{Resources: p.resources, Checks: len(checks)}
	for _, c := range checks {
		for _, f := range c.run(p) {
			f.Check = c.name
			report.Findings = append(report.Findings, f)
		}
	}
	return report, nil
}

// --------------------------------------------------------------- the project

// field is one column of a model, as the model declares it.
type field struct {
	Go   string // the Go field name
	Type string // the Go type, as written
	JSON string // the json tag, which is the column name
}

type model struct {
	Name        string
	Fields      []field
	byGo        map[string]field
	OwnerColumn string // from GetOwnerID, empty when the resource is shared
	TenantOwned bool
	AppendOnly  bool
}

func (m *model) encrypted() []string {
	var out []string
	for _, f := range m.Fields {
		if strings.Contains(f.Type, "crypto.EncryptedString") {
			out = append(out, f.JSON)
		}
	}
	return out
}

// hasUserRelation reports whether the model points at a user, which is what
// makes "should this be owned?" worth asking.
func (m *model) hasUserRelation() bool {
	for _, f := range m.Fields {
		if f.Type == "User" || f.Type == "*User" || f.JSON == "user_id" {
			return true
		}
	}
	return false
}

type project struct {
	root, apiRoot string
	env           map[string]string
	goMod         string
	routes        string            // every routes file, concatenated
	routeFiles    map[string]string // resource snake name → its own routes file
	models        map[string]*model
	handlers      map[string]string // resource → handlers/<snake>.go
	services      map[string]string // resource → services/<snake>.go
	publics       map[string]string // resource → handlers/<snake>_public.go
	snake         map[string]string // resource → its file name
	resources     []string          // resources with a handler and a service
	hasTenant     bool
}

var (
	structRe   = regexp.MustCompile(`(?m)^type (\w+) struct \{`)
	fieldRe    = regexp.MustCompile("^\\s*(\\w+)\\s+([\\w\\.\\*\\[\\]]+)\\s+`[^`]*json:\"([a-zA-Z0-9_]+)")
	ownerRe    = regexp.MustCompile(`func \(\w+ \*(\w+)\) GetOwnerID\(\) string \{\s*\n?\s*return \w+\.(\w+)`)
	registerRe = regexp.MustCompile(`appendonly\.Register\(&(\w+)\{\}\)`)
)

// loadProject reads everything the checks look at, once.
func loadProject(root string) (*project, error) {
	p := &project{
		root:       root,
		apiRoot:    apiRootOf(root),
		models:     map[string]*model{},
		handlers:   map[string]string{},
		services:   map[string]string{},
		publics:    map[string]string{},
		snake:      map[string]string{},
		routeFiles: map[string]string{},
	}
	internal := filepath.Join(p.apiRoot, "internal")
	if !dirExists(internal) {
		return nil, fmt.Errorf("no Go API at %s: run grit doctor inside a Grit project", p.apiRoot)
	}
	p.env = readEnv(filepath.Join(root, ".env"))
	p.goMod = readFile(filepath.Join(p.apiRoot, "go.mod"))
	p.hasTenant = dirExists(filepath.Join(internal, "tenant"))

	for name, src := range goFiles(filepath.Join(internal, "models")) {
		for _, m := range parseModels(src) {
			p.models[m.Name] = m
			p.snake[m.Name] = strings.TrimSuffix(name, ".go")
		}
	}
	for name, src := range goFiles(filepath.Join(internal, "handlers")) {
		switch {
		case strings.HasSuffix(name, "_public.go"):
			p.publics[strings.TrimSuffix(name, "_public.go")] = src
		case strings.Contains(name, "_"):
			// _import.go, _tree.go, _variant.go: not the resource's own handler.
		default:
			p.handlers[strings.TrimSuffix(name, ".go")] = src
		}
	}
	for name, src := range goFiles(filepath.Join(internal, "services")) {
		if !strings.Contains(name, "_") {
			p.services[strings.TrimSuffix(name, ".go")] = src
		}
	}
	for name, src := range goFiles(filepath.Join(internal, "routes")) {
		p.routes += src
		if strings.HasSuffix(name, "_routes.go") {
			p.routeFiles[strings.TrimSuffix(name, "_routes.go")] = src
		}
	}

	// A generated resource is one the generator wrote, which is not the same as
	// "has a model and a handler": the framework ships handlers and services for
	// its own tables, and a session or a passkey referencing a user is not a
	// user-owned resource. Asked the loose way, every project reported four
	// framework tables as possibly-owned, including a project with no resources
	// of its own at all.
	//
	// Every generated resource declares a bulk request and a list config, and no
	// built-in declares either.
	for name := range p.models {
		snake := p.snake[name]
		code := p.handlers[snake] + p.services[snake]
		if !strings.Contains(code, "Bulk"+name+"Request") && !strings.Contains(code, lowerFirst(name)+"ListConfig") {
			continue
		}
		p.resources = append(p.resources, name)
	}
	sort.Strings(p.resources)
	return p, nil
}

// parseModels reads the structs in one model file, with their fields and the
// markers the flags leave behind.
func parseModels(src string) []*model {
	var out []*model
	lines := strings.Split(src, "\n")
	for _, loc := range structRe.FindAllStringSubmatchIndex(src, -1) {
		name := src[loc[2]:loc[3]]
		m := &model{Name: name, byGo: map[string]field{}}
		start := strings.Count(src[:loc[0]], "\n") + 1
		for i := start; i < len(lines); i++ {
			line := lines[i]
			if strings.HasPrefix(line, "}") {
				break
			}
			if strings.Contains(line, "tenant.Owned") {
				m.TenantOwned = true
			}
			if f := fieldRe.FindStringSubmatch(line); f != nil {
				col := field{Go: f[1], Type: f[2], JSON: f[3]}
				m.Fields = append(m.Fields, col)
				m.byGo[col.Go] = col
			}
		}
		out = append(out, m)
	}
	// GetOwnerID and the append-only registration sit outside the struct.
	for _, o := range ownerRe.FindAllStringSubmatch(src, -1) {
		for _, m := range out {
			if m.Name == o[1] {
				m.OwnerColumn = m.byGo[o[2]].JSON
			}
		}
	}
	for _, r := range registerRe.FindAllStringSubmatch(src, -1) {
		for _, m := range out {
			if m.Name == r[1] {
				m.AppendOnly = true
			}
		}
	}
	return out
}

// code returns the resource's service when it has one, and its handler as well,
// so a check can look wherever the queries happen to live: before v3.224.0 they
// were in the handler.
func (p *project) code(resource string) string {
	snake := p.snake[resource]
	return p.services[snake] + "\n" + p.handlers[snake]
}

// --------------------------------------------------------------- the checks

// An encrypted field with no key is a plain column.
//
// crypto.EncryptedString stores the value as it is when FIELD_ENCRYPTION_KEY is
// unset, so nothing fails and the column is readable in the database. Found
// while verifying v3.224.0: a "secret" column came back in clear text from
// Postgres, and the API had answered 201.
func checkEncryptionKey(p *project) []Finding {
	var with []string
	for _, name := range p.resources {
		if len(p.models[name].encrypted()) > 0 {
			with = append(with, name)
		}
	}
	if len(with) == 0 || strings.TrimSpace(p.env["FIELD_ENCRYPTION_KEY"]) != "" {
		return nil
	}
	return []Finding{{
		Level:    "error",
		Resource: strings.Join(with, ", "),
		Message:  "has encrypted fields, and FIELD_ENCRYPTION_KEY is not set in .env, so they are stored in the clear",
		Fix:      "set FIELD_ENCRYPTION_KEY to 32 random bytes (openssl rand -base64 32), or set it in the deployment's environment",
	}}
}

// Ciphertext differs on every write, so searching, sorting or filtering an
// encrypted column returns nothing or noise. The generator leaves it out of
// the whitelists; a hand-edited one can put it back.
func checkEncryptedWhitelists(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		src := p.code(name)
		for _, col := range p.models[name].encrypted() {
			for _, line := range strings.Split(src, "\n") {
				if !strings.Contains(line, "Searchable:") && !strings.Contains(line, "Sortable:") && !strings.Contains(line, "Filterable:") {
					continue
				}
				if strings.Contains(line, `"`+col+`"`) {
					out = append(out, Finding{
						Level:    "warning",
						Resource: name + "." + col,
						Message:  "is encrypted and is in a search, sort or filter whitelist, where it matches nothing",
						Fix:      "remove " + col + " from that list: ciphertext differs on every write",
					})
					break
				}
			}
		}
	}
	return out
}

// An owned resource has to be scoped on every path that reads or writes by id,
// not only the obvious ones.
//
// The list is the path people forget, and an export is the list without pages.
// PDF and Patch shipped unguarded once, and a second ordinary account printed
// another user's record and rewrote it.
func checkOwnedScoping(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		m := p.models[name]
		if m.OwnerColumn == "" {
			continue
		}
		service, handler := p.services[p.snake[name]], p.handlers[p.snake[name]]
		// Since v3.224.0 the service scopes; before it, the handler did.
		newShape := strings.Contains(service, "authz.ScopeOwned(") || strings.Contains(service, "authz.Owns(")
		oldShape := strings.Contains(handler, "authz.ScopeToOwner(") || strings.Contains(handler, "authz.OwnsOr404(")
		if !newShape && !oldShape {
			out = append(out, Finding{
				Level:    "error",
				Resource: name,
				Message:  "is owned (" + m.OwnerColumn + ") and nothing scopes it: every signed-in caller can read and change every row",
				Fix:      "regenerate it with --owned-by, or scope its queries with authz.ScopeOwned and authz.Owns",
			})
			continue
		}
		if !newShape {
			continue // an older project, guarded the way its generation did
		}
		for _, want := range []struct{ method, call string }{
			{"List", "authz.ScopeOwned("},
			{"Export", "authz.ScopeOwned("},
			{"Bulk", "authz.ScopeOwned("},
			{"GetByID", "authz.Owns("},
			{"load", "authz.Owns("},
		} {
			body, found := goMethod(service, name+"Service) "+want.method+"(")
			if !found || strings.Contains(body, want.call) {
				continue
			}
			out = append(out, Finding{
				Level:    "error",
				Resource: name + "Service." + want.method,
				Message:  "does not scope an owned resource, so it reaches other users' rows",
				Fix:      "add " + want.call + "...) there, as the other methods have it",
			})
		}
	}
	return out
}

// The owner comes from the session, never from the request body. A client that
// can set it can file rows under somebody else's name.
func checkOwnerFromBody(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		m := p.models[name]
		if m.OwnerColumn == "" {
			continue
		}
		if strings.Contains(p.handlers[p.snake[name]], `json:"`+m.OwnerColumn+`"`) {
			out = append(out, Finding{
				Level:    "error",
				Resource: name,
				Message:  "accepts " + m.OwnerColumn + " from the request body, so a caller can give a row to another account",
				Fix:      "drop it from the create and update requests: the service stamps it from the caller",
			})
		}
	}
	return out
}

// A resource that belongs to a user and is not scoped to one is the default
// every signed-in account can read. Sometimes that is right, which is why this
// is a question rather than an error.
func checkUnownedUserResource(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		m := p.models[name]
		if m.OwnerColumn != "" || !m.hasUserRelation() {
			continue
		}
		out = append(out, Finding{
			Level:    "warning",
			Resource: name,
			Message:  "references a user but is not scoped to one: every signed-in caller sees every row",
			Fix:      "if each row belongs to its user, regenerate with --owned-by user; if the resource is shared, nothing to do",
		})
	}
	return out
}

// The multitenant plugin scopes by a column. A model without tenant.Owned has
// no column to scope by, so its table is shared across every organization, and
// the plugin cannot tell that from a table that is meant to be.
func checkTenantOwned(p *project) []Finding {
	if !p.hasTenant {
		return nil
	}
	var out []Finding
	for _, name := range p.resources {
		if p.models[name].TenantOwned {
			continue
		}
		out = append(out, Finding{
			Level:    "warning",
			Resource: name,
			Message:  "is not tenant-owned in a project with the multitenant plugin: its rows are shared across organizations",
			Fix:      "add tenant.Owned to the model and run grit migrate, or leave it if the table is deliberately global",
		})
	}
	return out
}

// checkDatabaseProvider reports a database that will not survive production.
//
// DB_PROVIDER=memory keeps everything in RAM. The app works perfectly until it
// restarts, at which point every row is gone: there is no error to see and no
// file to recover, which is why this is an error rather than a warning when
// APP_ENV says production.
//
// SQLite in production is a decision, not a mistake, so it is a warning that
// names what you are taking on rather than a refusal.
func checkDatabaseProvider(p *project) []Finding {
	provider := strings.ToLower(strings.TrimSpace(p.env["DB_PROVIDER"]))
	if provider == "" {
		return nil
	}
	production := strings.EqualFold(strings.TrimSpace(p.env["APP_ENV"]), "production")
	if !production {
		return nil
	}

	switch provider {
	case "memory", ":memory:":
		return []Finding{{
			Level:    "error",
			Check:    "database-not-durable",
			Resource: "DB_PROVIDER",
			Message: "is memory while APP_ENV is production: the whole database is in RAM and every row " +
				"is gone at the next restart, with nothing to recover from",
			Fix: "set DB_PROVIDER to postgres, mysql or sqlite, with the matching settings in .env",
		}}
	case "sqlite", "sqlite3", "file":
		return []Finding{{
			Level:    "warning",
			Check:    "database-not-durable",
			Resource: "DB_PROVIDER",
			Message: "is sqlite in production: one file on one machine, so no second replica can read it, " +
				"and a container without a mounted volume loses it on redeploy",
			Fix: "keep it only if the file is on a persistent volume and one process writes it; otherwise " +
				"use postgres or mysql",
		}}
	}
	return nil
}

// checkTenantMiddleware looks for the organization resolver on every
// authenticated route group, in the right order.
//
// The multitenant plugin used to mount middleware.Tenant on the protected group
// only. The staff group carries every DELETE and every bulk route, and the admin
// group every admin-only endpoint, so on those the request never resolved an
// organization: the scoping callback failed closed and the answer was a 500. A
// project that installed the plugin before v3.233.0 still looks like that, and
// grit upgrade does not rewrite routes.go, which plugins and people both edit.
//
// Order matters as much as presence. A role held through an organization
// membership has to be resolved before RequireStaff reads the caller's grants, or
// a member who is staff of the organization they are acting in gets a 403.
func checkTenantMiddleware(p *project) []Finding {
	if !p.hasTenant || p.routes == "" {
		return nil
	}

	var out []Finding
	for _, group := range []struct{ name, mount, gate string }{
		{"staff", "staff.Use(middleware.Tenant(db))", "staff.Use(middleware.RequireStaff())"},
		{"admin", "admin.Use(middleware.Tenant(db))", `admin.Use(middleware.RequireRole("ADMIN"))`},
	} {
		at := strings.Index(p.routes, group.mount)
		if at < 0 {
			out = append(out, Finding{
				Level:    "error",
				Check:    "tenant-middleware-missing",
				Resource: group.name + " routes",
				Message: "the " + group.name + " group does not resolve the active organization, so every " +
					"tenant-owned row it touches answers 500: the DELETE and bulk routes live here",
				Fix: "add " + group.mount + " in internal/routes/routes.go, above " + group.gate,
			})
			continue
		}
		if gate := strings.Index(p.routes, group.gate); gate >= 0 && at > gate {
			out = append(out, Finding{
				Level:    "error",
				Check:    "tenant-middleware-missing",
				Resource: group.name + " routes",
				Message: "the organization is resolved after the permission gate, so a role held through an " +
					"organization membership grants nothing on these routes",
				Fix: "move " + group.mount + " above " + group.gate,
			})
		}
	}
	return out
}

// checkTenantSSOProvisioning reports the gap between single sign-on and
// organizations.
//
// A user created by SSO or a social login joins no organization: nothing in the
// provisioning path knows about them, and which organization somebody belongs to
// is a policy decision (their email domain, the connection they came through, a
// group the identity provider released) that the framework cannot make. The
// result, with the rest of v3.233.0 in place, is an honest 400 NO_ORGANIZATION on
// every tenant-owned endpoint rather than a silent empty list, which is the right
// failure and still a surprise if nobody said so.
func checkTenantSSOProvisioning(p *project) []Finding {
	if !p.hasTenant {
		return nil
	}
	if !strings.Contains(p.routes, "samlHandler") && !strings.Contains(p.routes, "/auth/sso") &&
		!strings.Contains(p.routes, "oauthHandler") {
		return nil
	}
	return []Finding{{
		Level:    "warning",
		Check:    "tenant-sso-no-organization",
		Resource: "single sign-on",
		Message: "a user provisioned by SSO or a social login joins no organization, so every tenant-owned " +
			"endpoint answers NO_ORGANIZATION until somebody adds them to one",
		Fix: "add the membership where you know the policy: map the email domain or the SSO connection to an " +
			"organization after sign-in, or have an administrator invite them",
	}}
}

// piiColumns are column names that usually hold something that should not be
// readable in a database dump. Deliberately short: a name that is merely
// personal (email, phone) is not on it, because encrypting what the app filters
// and sorts by breaks the app, and a warning nobody can act on is noise.
var piiColumns = []string{
	"ssn", "social_security", "national_id", "passport", "tax_id",
	"date_of_birth", "dob", "diagnosis", "medical", "prescription",
	"card_number", "cvv", "iban", "bank_account", "account_number", "routing_number",
}

func checkPII(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		for _, f := range p.models[name].Fields {
			if f.Type != "string" || strings.Contains(f.Type, "crypto.") {
				continue
			}
			for _, needle := range piiColumns {
				if !strings.Contains(f.JSON, needle) {
					continue
				}
				out = append(out, Finding{
					Level:    "warning",
					Resource: name + "." + f.JSON,
					Message:  "looks like data that should not be readable in the database, and is a plain column",
					Fix:      "generate the field as " + f.JSON + ":string:encrypted, or leave it if the value is not sensitive here",
				})
				break
			}
		}
	}
	return out
}

// An append-only resource refuses an update in the model and in a database
// trigger. A route that still offers one can only fail, and it tells a reader
// the row can be edited.
func checkAppendOnlyRoutes(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		if !p.models[name].AppendOnly {
			continue
		}
		src, ok := p.routeFiles[p.snake[name]]
		if !ok {
			src = p.routes // an older project keeps every route in routes.go
		}
		for _, verb := range []string{".PUT(", ".PATCH(", ".DELETE("} {
			if !strings.Contains(src, verb) {
				continue
			}
			out = append(out, Finding{
				Level:    "error",
				Resource: name,
				Message:  "is append-only and still mounts " + strings.Trim(verb, ".(") + ": the write is refused, so the endpoint only ever fails",
				Fix:      "remove that route, or stop registering the model with appendonly if corrections are allowed now",
			})
			break
		}
	}
	return out
}

// GORM Studio browses and edits every table. Unauthenticated, it is the
// database on the internet.
func checkStudio(p *project) []Finding {
	if strings.EqualFold(p.env["GORM_STUDIO_ENABLED"], "false") {
		return nil
	}
	production := strings.HasPrefix(strings.ToLower(p.env["APP_ENV"]), "prod")
	user, password := strings.TrimSpace(p.env["GORM_STUDIO_USERNAME"]), strings.TrimSpace(p.env["GORM_STUDIO_PASSWORD"])
	if user == "" || password == "" {
		return []Finding{{
			Level:    "error",
			Resource: "/studio",
			Message:  "is mounted with no login: the scaffold adds basic auth only when GORM_STUDIO_USERNAME and GORM_STUDIO_PASSWORD are both set",
			Fix:      "set both in .env, or set GORM_STUDIO_ENABLED=false",
		}}
	}
	if password == "studio" || strings.Contains(password, "change-me") {
		level := "warning"
		if production {
			level = "error"
		}
		return []Finding{{
			Level:    level,
			Resource: "/studio",
			Message:  "still has its default password, and it can browse and edit every table",
			Fix:      "set GORM_STUDIO_PASSWORD in .env: openssl rand -hex 16",
		}}
	}
	return nil
}

// The credentials that guard the dashboards, and the secret that signs every
// session. Sentinel and Pulse refuse to mount in production on their defaults,
// which is a floor, not a setting anyone should rely on.
func checkCredentials(p *project) []Finding {
	production := strings.HasPrefix(strings.ToLower(p.env["APP_ENV"]), "prod")
	level := "warning"
	if production {
		level = "error"
	}
	var out []Finding
	for _, c := range []struct{ key, bad, what string }{
		{"SENTINEL_PASSWORD", "sentinel", "the Sentinel dashboard"},
		{"PULSE_PASSWORD", "pulse", "the Pulse dashboard"},
		{"SENTINEL_SECRET_KEY", "sentinel-secret-change-me", "Sentinel's dashboard sessions"},
	} {
		v := strings.TrimSpace(p.env[c.key])
		if v == "" || v == c.bad || strings.Contains(v, "change-me") || strings.Contains(v, "generate-a-random") {
			out = append(out, Finding{
				Level:    level,
				Resource: c.key,
				Message:  "is still the default, and it guards " + c.what,
				Fix:      "set " + c.key + " in .env: openssl rand -hex 32",
			})
		}
	}
	// The JWT secret signs access tokens. Short or empty, they can be forged.
	if jwt := strings.TrimSpace(p.env["JWT_SECRET"]); len(jwt) < 32 || strings.Contains(jwt, "change-me") {
		out = append(out, Finding{
			Level:    "error",
			Resource: "JWT_SECRET",
			Message:  "is empty, a placeholder, or under 32 characters: access tokens can be forged",
			Fix:      "set JWT_SECRET in .env: openssl rand -hex 32",
		})
	}
	// v3.226.0 keys the audit chain. Absent, the chain still works, unkeyed.
	if strings.TrimSpace(p.env["SENTINEL_AUDIT_KEY"]) == "" {
		out = append(out, Finding{
			Level:    "warning",
			Resource: "SENTINEL_AUDIT_KEY",
			Message:  "is not set, so the security audit log's chain is unkeyed: someone with database access can rewrite it and relink the chain",
			Fix:      "set SENTINEL_AUDIT_KEY in .env: openssl rand -hex 32",
		})
	}
	return out
}

// The libraries every API mounts, against the versions this CLI ships with.
// grit upgrade raises them; this says so before something goes wrong.
func checkFrameworkDeps(p *project) []Finding {
	var out []Finding
	for _, d := range scaffold.FrameworkDepFloors() {
		m := regexp.MustCompile(`(?m)^\s*(?:require\s+)?` + regexp.QuoteMeta(d.Path) + `\s+(v\S+)`).FindStringSubmatch(p.goMod)
		if m == nil || !scaffold.VersionBelow(m[1], d.Floor) {
			continue
		}
		out = append(out, Finding{
			Level:    "warning",
			Resource: d.Path,
			Message:  "is at " + m[1] + ", below " + d.Floor + ": " + d.Why,
			Fix:      "run grit upgrade, which raises it",
		})
	}
	return out
}

// Rate limits and lockouts counted per process give a client one allowance per
// replica. Only worth saying when the project has the Redis to share them.
func checkSentinelCounters(p *project) []Finding {
	mounted := strings.Contains(p.routes, "sentinel.MountE(") || strings.Contains(p.routes, "sentinel.Mount(")
	hasRedis := p.env["REDIS_URL"] != "" || p.env["REDIS_HOST"] != ""
	if !mounted || !hasRedis || strings.Contains(p.routes, "redisstore.New(") {
		return nil
	}
	return []Finding{{
		Level:    "warning",
		Resource: "sentinel",
		Message:  "counts rate limits and AuthShield lockouts in this process, so each replica allows a client the full limit again",
		Fix:      "pass Counters: redisstore.New(svc.Cache.Client()) to sentinel.Config (see the Sentinel docs, Replicas)",
	}}
}

// A public response is an allowlist. The generator refuses to publish anything
// that smells like cost, stock or internal notes; the file is yours to edit,
// and an edit is how one of those reaches the internet.
var heldBackColumns = []string{
	"cost", "margin", "profit", "internal", "note", "secret", "password",
	"token", "commission", "supplier", "wholesale", "stock", "inventory",
}

func checkPublicAllowlist(p *project) []Finding {
	var out []Finding
	for snake, src := range p.publics {
		block, ok := structBlock(src, "struct {")
		if !ok {
			continue
		}
		for _, line := range strings.Split(block, "\n") {
			f := fieldRe.FindStringSubmatch(line)
			if f == nil {
				continue
			}
			for _, needle := range heldBackColumns {
				if !strings.Contains(f[3], needle) {
					continue
				}
				out = append(out, Finding{
					Level:    "warning",
					Resource: snake + "_public.go: " + f[3],
					Message:  "is published to anonymous callers, and it is the kind of column the generator holds back",
					Fix:      "remove it from the public struct unless the value is meant to be public",
				})
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Resource < out[j].Resource })
	return out
}

// --------------------------------------------------------------- small helpers

// lowerFirst is the generator's camel form of a resource name: Invoice becomes
// invoice, as in invoiceListConfig.
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// goMethod returns one method's body, from its signature to its closing brace.
func goMethod(src, signature string) (string, bool) {
	start := strings.Index(src, signature)
	if start < 0 {
		return "", false
	}
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return src[start:], true
	}
	return src[start : start+end], true
}

// structBlock returns the first struct body in src.
func structBlock(src, opener string) (string, bool) {
	start := strings.Index(src, opener)
	if start < 0 {
		return "", false
	}
	end := strings.Index(src[start:], "\n}")
	if end < 0 {
		return "", false
	}
	return src[start : start+end], true
}

// apiRootOf returns the directory holding the Go API: the project root for a
// single-binary layout, apps/api for a monorepo.
func apiRootOf(root string) string {
	mono := filepath.Join(root, "apps", "api")
	if _, err := os.Stat(filepath.Join(mono, "go.mod")); err == nil {
		return mono
	}
	return root
}

// readEnv reads KEY=VALUE lines. Values keep their '#' unless it is spaced off
// as a trailing comment, which is how the scaffold writes .env.example.
func readEnv(path string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(readFile(path), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if i := strings.Index(value, " #"); i >= 0 {
			value = value[:i]
		}
		out[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	return out
}

// goFiles returns the non-test Go files in dir, by file name.
func goFiles(dir string) map[string]string {
	out := map[string]string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		out[name] = readFile(filepath.Join(dir, name))
	}
	return out
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
