package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairReviewCriticals brings a project scaffolded before v3.241.0 up to the
// fixes for the critical findings of the contact-app security review.
//
//   - A generated resource asked for a permission only to delete. List, get,
//     export, import, create and update were open to any signed-in account, and
//     with open registration that is anybody: the review exported a whole
//     address book from a fresh account.
//   - The default sync registry held users and uploads, and a push was a
//     generic write, so one request made any account ADMIN.
//   - .env.example was a copy of .env, real secrets included, and *.db files
//     were not ignored.
//
// Route files and routes.go belong to the developer, so each edit anchors on
// the exact text the generator wrote and leaves anything else alone. The sync
// handler is framework code and is delivered whole, behind the manifest guard.
func repairReviewCriticals(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	routesDir := filepath.Join(apiRoot, "internal", "routes")
	staff := false
	if data, err := os.ReadFile(filepath.Join(routesDir, "resources.go")); err == nil {
		staff = reviewStaffFieldRe.Match(data)
	}
	owned, err := ownedResourceHandlers(filepath.Join(apiRoot, "internal", "handlers"))
	if err != nil {
		return err
	}

	files, err := goSources(routesDir)
	if err != nil {
		return err
	}
	for _, path := range files {
		if !strings.HasSuffix(path, "_routes.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairRouteSource(src, owned, staff, opts.Module())
		}); err != nil {
			return err
		}
	}

	if f := filepath.Join(routesDir, "routes.go"); fileExists(f) {
		if err := repairSourceFile(root, m, f, func(src string) (string, []string, []string) {
			out, fixed, warn := repairSyncRegistrySource(src)
			out, treeFixed, treeWarn := repairTreeRoutesSource(out, owned)
			return out, append(fixed, treeFixed...), append(warn, treeWarn...)
		}); err != nil {
			return err
		}
	}
	if f := filepath.Join(apiRoot, "internal", "handlers", "sync.go"); fileExists(f) {
		if err := writeFile(f, strings.ReplaceAll(apiSyncHandlerGo(), "{{MODULE}}", opts.Module())); err != nil {
			return err
		}
	}

	// A desktop app pulled users and uploads in the background, which the API
	// now refuses, so every sync round would fail on them.
	for _, f := range []string{
		filepath.Join(root, "apps", "desktop", "app.go"),
		filepath.Join(root, "app.go"),
	} {
		if fileExists(f) {
			if err := repairSourceFile(root, m, f, repairDesktopSyncTablesSource); err != nil {
				return err
			}
		}
	}

	if err := scrubEnvExample(root); err != nil {
		return err
	}
	return ignoreLocalDatabases(root)
}

var (
	desktopSyncTablesRe    = regexp.MustCompile(`(?s)var syncTables = \[\]string\{\n.*?\n\}`)
	desktopIdentityTableRe = regexp.MustCompile(`(?m)^[ \t]*"(?:users|uploads)",[ \t]*\n`)
)

// repairDesktopSyncTablesSource takes users and uploads out of a desktop app's
// background sync list.
func repairDesktopSyncTablesSource(src string) (string, []string, []string) {
	block := desktopSyncTablesRe.FindString(src)
	if block == "" {
		return src, nil, nil
	}
	cleaned := desktopIdentityTableRe.ReplaceAllString(block, "")
	if cleaned == block {
		return src, nil, nil
	}
	return strings.Replace(src, block, cleaned, 1), []string{"the desktop app no longer syncs users or uploads, which the API refuses"}, nil
}

var (
	reviewStaffFieldRe = regexp.MustCompile(`Staff\s+\*gin\.RouterGroup`)
	routeHandlerRe     = regexp.MustCompile(`h := &handlers\.(\w+)Handler\{`)
	handlerPascalRe    = regexp.MustCompile(`func \(h \*(\w+)Handler\) List\(`)
	// A route the generator wrote on the protected group: m.Protected.GET("/contacts/:id", h.GetByID)
	protectedRouteRe = regexp.MustCompile(`(?m)^([ \t]*)m\.Protected\.(GET|POST|PUT|PATCH)\("/(\w+)((?:/[^"]*)?)", h\.(\w+)\)[ \t]*$`)
	syncIdentityRe   = regexp.MustCompile(`(?m)^[ \t]*syncRegistry\.Register\("(?:users|uploads)", &models\.(?:User|Upload)\{\}\)[ \t]*\n`)
)

// routePermission is the permission each generated handler asks for, the same
// table the generator uses. A handler not listed was written by the developer
// and its route is left as it is.
var routePermission = map[string]string{
	"List":       "view",
	"Export":     "view",
	"Import":     "create",
	"Template":   "view",
	"GetByID":    "view",
	"PDF":        "view",
	"Create":     "create",
	"Update":     "edit",
	"Patch":      "edit",
	"Workflow":   "view",
	"Transition": "edit",
}

// ownedResourceHandlers names the resources whose handlers scope rows to their
// owner. Their routes are meant to be reachable by any signed-in account.
func ownedResourceHandlers(dir string) (map[string]bool, error) {
	files, err := goSources(dir)
	if err != nil {
		return nil, err
	}
	owned := map[string]bool{}
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		src := string(raw)
		if !strings.Contains(src, "authz.ScopeToOwner(") && !strings.Contains(src, "authz.OwnsOr404(") {
			continue
		}
		for _, match := range handlerPascalRe.FindAllStringSubmatch(src, -1) {
			owned[match[1]] = true
		}
	}
	// --tenant-owned: the tenant middleware narrows every query to the caller's
	// organization, so those resources are scoped too.
	models, err := goSources(filepath.Join(filepath.Dir(dir), "models"))
	if err != nil {
		return nil, err
	}
	for _, path := range models {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		for _, match := range tenantModelRe.FindAllStringSubmatch(string(raw), -1) {
			owned[match[1]] = true
		}
	}
	return owned, nil
}

// tenantModelRe finds a model that embeds tenant.Owned.
var tenantModelRe = regexp.MustCompile(`type (\w+) struct \{[^}]*\n\s*tenant\.Owned\s*\n`)

// repairRouteSource moves a shared resource's generated routes off the
// protected group: onto the staff group behind the permission for each verb, or
// onto the admin group in a project that has no staff group.
func repairRouteSource(src string, owned map[string]bool, staff bool, module string) (string, []string, []string) {
	handler := routeHandlerRe.FindStringSubmatch(src)
	if handler == nil || owned[handler[1]] {
		return src, nil, nil
	}

	moved := 0
	out := protectedRouteRe.ReplaceAllStringFunc(src, func(line string) string {
		p := protectedRouteRe.FindStringSubmatch(line)
		indent, method, plural, rest, name := p[1], p[2], p[3], p[4], p[5]
		perm, generated := routePermission[name]
		if !generated {
			return line
		}
		moved++
		if staff {
			return fmt.Sprintf(`%sm.Staff.%s("/%s%s", middleware.RequireRole("ADMIN", "perm:%s.%s"), h.%s)`,
				indent, method, plural, rest, plural, perm, name)
		}
		return fmt.Sprintf(`%sm.Admin.%s("/%s%s", h.%s)`, indent, method, plural, rest, name)
	})
	if moved == 0 {
		return src, nil, nil
	}

	if staff && !strings.Contains(out, `/internal/middleware"`) {
		handlers := fmt.Sprintf("%q", module+"/internal/handlers")
		if !strings.Contains(out, handlers) {
			return src, nil, []string{"its shared routes are open to any signed-in account, and its imports are not the generated ones, so it was left alone: move them to m.Staff with a permission by hand"}
		}
		out = strings.Replace(out, handlers, handlers+"\n\t"+fmt.Sprintf("%q", module+"/internal/middleware"), 1)
	}

	where := "the admin group"
	if staff {
		where = "a permission"
	}
	return out, []string{fmt.Sprintf("%d shared %s now ask for %s instead of any signed-in account",
		moved, plural(moved, "route", "routes"), where)}, nil
}

// repairSyncRegistrySource takes users and uploads out of the sync registry.
func repairSyncRegistrySource(src string) (string, []string, []string) {
	out := syncIdentityRe.ReplaceAllString(src, "")
	if out == src {
		return src, nil, nil
	}
	return out, []string{"users and uploads are no longer synced (a push could make any account ADMIN)"}, nil
}

// scrubEnvExample replaces generated secrets in .env.example with placeholders.
func scrubEnvExample(root string) error {
	path := filepath.Join(root, ".env.example")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading .env.example: %w", err)
	}
	src := string(raw)
	found := len(generatedSecret.FindAllString(src, -1)) + len(encryptionKeyLine.FindAllString(src, -1))
	if found == 0 {
		return nil
	}
	out := generatedSecret.ReplaceAllString(src, "${1}=CHANGE_ME${2}")
	out = encryptionKeyLine.ReplaceAllString(out, "${1}=CHANGE_ME${2}")
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing .env.example: %w", err)
	}
	fmt.Printf("  ✓ .env.example: %d generated %s replaced with CHANGE_ME\n", found, plural(found, "secret", "secrets"))
	fmt.Println("    They are the values in .env. If .env.example was ever committed or shared, rotate them.")
	return nil
}

const localDatabaseIgnores = `# Local databases. SQLite files hold real data: users with their password
# hashes, sessions and API keys, as well as Sentinel's WAF log.
*.db
*.db-shm
*.db-wal
*.sqlite
*.sqlite3
`

// ignoreLocalDatabases adds the database patterns to .gitignore once.
func ignoreLocalDatabases(root string) error {
	path := filepath.Join(root, ".gitignore")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading .gitignore: %w", err)
	}
	src := string(raw)
	for _, line := range strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "*.db" {
			return nil
		}
	}
	block := localDatabaseIgnores
	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
		block = strings.ReplaceAll(block, "\n", nl)
	}
	if src != "" && !strings.HasSuffix(src, "\n") {
		src += nl
	}
	if err := os.WriteFile(path, []byte(src+nl+block), 0o644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}
	fmt.Println("  ✓ .gitignore: local databases (*.db, *.sqlite) are ignored")
	fmt.Println("    A database already committed stays in history until you git rm --cached it.")
	return nil
}
