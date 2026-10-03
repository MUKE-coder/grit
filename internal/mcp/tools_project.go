package mcp

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/doctor"
	"github.com/MUKE-coder/grit/v3/internal/generate"
	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// projectInfoTool adapts projectInfo to the registry's signature. The argument
// is ignored: the tool takes none.
func (s *Server) projectInfoTool(json.RawMessage) (string, error) { return s.projectInfo() }

// listResources reports what `grit generate resource` has made in this project.
//
// Found by the marker the generator leaves rather than by "has a model and a
// handler". The framework ships handlers and services for its own tables, so the
// loose reading reported SAML, sessions and passkeys as things somebody
// generated, and missed the sample blog, whose files are named the older way.
//
// Every generated resource declares a Bulk<Name>Request and no built-in does,
// which is the same marker `grit doctor` uses to decide what it is auditing.
func (s *Server) listResources(json.RawMessage) (string, error) {
	api := apiRoot(s.Root)
	modelsDir := filepath.Join(api, "internal", "models")
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return "", fmt.Errorf("reading models directory: %w", err)
	}

	// Where the marker could be: the handler and the service, whatever they are
	// called. Read once rather than per model.
	code := map[string]string{}
	for _, dir := range []string{"handlers", "services"} {
		path := filepath.Join(api, "internal", dir)
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
				continue
			}
			full := filepath.Join(path, f.Name())
			data, err := os.ReadFile(full)
			if err != nil {
				continue
			}
			code[full] = string(data)
		}
	}

	out := []map[string]interface{}{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		modelFile := filepath.Join(modelsDir, entry.Name())
		structs, err := generate.ParseGoStructs(modelFile)
		if err != nil {
			// One unparseable file should not blank the whole answer.
			continue
		}
		for _, st := range structs {
			marker := "Bulk" + st.Name + "Request"
			found := false
			for path, src := range code {
				if !strings.Contains(src, marker) {
					continue
				}
				found = true
				// Only to prove the resource is generated. Which file is which
				// is settled by filesFor, which also finds the ones that never
				// mention the bulk request.
				_ = path
			}
			if !found {
				continue
			}
			names := generate.MakeNames(st.Name)
			files := s.filesFor(names)
			// The model file is the one already in hand, whatever it is called.
			files["model"] = relOrEmpty(s.Root, modelFile)
			out = append(out, map[string]interface{}{
				"name":  st.Name,
				"table": names.PluralSnake,
				"files": files,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i]["name"].(string) < out[j]["name"].(string)
	})
	return asJSON(map[string]interface{}{
		"count":     len(out),
		"resources": out,
		"note": "These are the resources `grit generate resource` wrote, found by the bulk-request " +
			"type it declares for each one. The framework's own endpoints, such as auth and users, " +
			"are not listed: they are not resources anybody generated and regenerating them is not " +
			"a thing you do.",
	})
}

// filesFor finds the files a generated resource owns.
//
// A list of candidates rather than a rule, because there are several layouts in
// the wild: a shared package with and without src/, an admin app of its own and
// an admin inside the web app, and handlers named widget.go or blog_handler.go.
// The first path that exists wins, and a label with nothing behind it is left
// out rather than reported as a file that is not there.
func (s *Server) filesFor(names generate.Names) map[string]string {
	api, err := filepath.Rel(s.Root, apiRoot(s.Root))
	if err != nil {
		api = "."
	}
	snake, plural := names.Snake, names.PluralKebab

	candidates := map[string][]string{
		"model": {
			filepath.Join(api, "internal", "models", snake+".go"),
		},
		"service": {
			filepath.Join(api, "internal", "services", snake+".go"),
			filepath.Join(api, "internal", "services", snake+"_service.go"),
		},
		"handler": {
			filepath.Join(api, "internal", "handlers", snake+".go"),
			filepath.Join(api, "internal", "handlers", snake+"_handler.go"),
		},
		"csv_import": {
			filepath.Join(api, "internal", "handlers", snake+"_import.go"),
		},
		"routes": {
			filepath.Join(api, "internal", "routes", snake+"_routes.go"),
		},
		"admin_resource": {
			filepath.Join("apps", "admin", "resources", plural, plural+".ts"),
			filepath.Join("apps", "admin", "src", "resources", plural, plural+".ts"),
			filepath.Join("apps", "web", "admin-panel", "resources", plural, plural+".ts"),
		},
		"admin_overlay": {
			filepath.Join("apps", "admin", "resources", plural, plural+".custom.tsx"),
			filepath.Join("apps", "admin", "src", "resources", plural, plural+".custom.tsx"),
			filepath.Join("apps", "web", "admin-panel", "resources", plural, plural+".custom.tsx"),
		},
		"hooks": {
			filepath.Join("apps", "admin", "hooks", "use-"+plural+".ts"),
			filepath.Join("apps", "admin", "src", "hooks", "use-"+plural+".ts"),
			filepath.Join("apps", "web", "hooks", "use-"+plural+".ts"),
		},
		"schema": {
			filepath.Join("packages", "shared", "schemas", names.Kebab+".ts"),
			filepath.Join("packages", "shared", "src", "schemas", names.Kebab+".ts"),
			filepath.Join("packages", "shared", "schemas", plural+".ts"),
		},
		"types": {
			filepath.Join("packages", "shared", "types", names.Kebab+".ts"),
			filepath.Join("packages", "shared", "src", "types", names.Kebab+".ts"),
			filepath.Join("packages", "shared", "types", plural+".ts"),
		},
	}

	out := map[string]string{}
	for label, paths := range candidates {
		for _, rel := range paths {
			if fileExists(filepath.Join(s.Root, rel)) {
				out[label] = filepath.ToSlash(rel)
				break
			}
		}
	}
	// The admin page is a directory, not a file.
	for _, rel := range []string{
		filepath.Join("apps", "admin", "app", "(dashboard)", "resources", plural),
		filepath.Join("apps", "admin", "src", "pages", "resources", plural),
		filepath.Join("apps", "web", "app", "admin", "(dashboard)", "resources", plural),
	} {
		if info, err := os.Stat(filepath.Join(s.Root, rel)); err == nil && info.IsDir() {
			out["admin_page"] = filepath.ToSlash(rel)
			break
		}
	}
	return out
}

// fileOwnership answers the question only the manifest can: which generated
// files are still Grit's and which the developer has taken over.
func (s *Server) fileOwnership(args json.RawMessage) (string, error) {
	var a struct {
		Status   string `json:"status"`
		Contains string `json:"contains"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}
	m, err := manifest.Load(s.Root)
	if err != nil {
		return "", err
	}

	want := strings.ToLower(strings.TrimSpace(a.Status))
	contains := strings.ToLower(strings.TrimSpace(a.Contains))
	counts := map[string]int{}
	files := []map[string]string{}

	for key := range m.Files {
		status := m.StatusOf(s.Root, key).String()
		counts[status]++
		if want != "" && status != want {
			continue
		}
		if contains != "" && !strings.Contains(strings.ToLower(key), contains) {
			continue
		}
		files = append(files, map[string]string{"path": key, "status": status})
	}
	sort.Slice(files, func(i, j int) bool { return files[i]["path"] < files[j]["path"] })

	note := "unchanged means the bytes on disk are the bytes Grit wrote, so grit upgrade will " +
		"keep updating the file. modified means somebody edited it, so an upgrade reports a " +
		"conflict there instead of overwriting the edit. missing means it was deleted."
	if len(m.Files) == 0 {
		note = "This project has no manifest, which means it was generated before v3.147.0 or " +
			"the .grit directory was removed. Nothing can be said about which files are still Grit's."
	}

	return asJSON(map[string]interface{}{
		"tracked": len(m.Files),
		"counts":  counts,
		"files":   files,
		"note":    note,
	})
}

// doctorReport runs the same checks as `grit doctor` and returns them as data.
func (s *Server) doctorReport(args json.RawMessage) (string, error) {
	var a struct {
		Level string `json:"level"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}
	report, err := doctor.Run(s.Root)
	if err != nil {
		return "", err
	}

	// A nil slice encodes as null, and a model reading "resources": null has to
	// decide whether that means none or means the field was not filled in.
	resources := report.Resources
	if resources == nil {
		resources = []string{}
	}

	want := strings.ToLower(strings.TrimSpace(a.Level))
	findings := []map[string]string{}
	for _, f := range report.Findings {
		if want != "" && f.Level != want {
			continue
		}
		findings = append(findings, map[string]string{
			"level":    f.Level,
			"check":    f.Check,
			"resource": f.Resource,
			"message":  f.Message,
			"fix":      f.Fix,
		})
	}

	return asJSON(map[string]interface{}{
		"checks":    report.Checks,
		"errors":    report.Errors(),
		"warnings":  report.Warnings(),
		"resources": resources,
		"findings":  findings,
	})
}

// listPermissions reads the permission catalogue out of the project's own
// authz package.
//
// Parsed with go/ast rather than matched with a regex, because the file
// documents itself with examples and a pattern over the text reports those as
// permissions that nothing checks.
//
// Every literal in that file elides its type: coreModules returns []Module{{Key:
// "access", Groups: []Group{{...}}}}, so the elements carry no Ident to match
// on. The walk keys off the one type that is written down, []Module, and then
// reads field names downward.
func (s *Server) listPermissions(json.RawMessage) (string, error) {
	path := filepath.Join(apiRoot(s.Root), "internal", "authz", "permissions.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return "", fmt.Errorf("parsing %s: %w", relOrEmpty(s.Root, path), err)
	}

	modules := []permModule{}
	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok || !isSliceOf(lit.Type, "Module") {
			return true
		}
		for _, el := range lit.Elts {
			if m, ok := readModule(el); ok {
				modules = append(modules, m)
			}
		}
		return true
	})

	sort.Slice(modules, func(i, j int) bool { return modules[i].Key < modules[j].Key })
	features, total := 0, 0
	for _, m := range modules {
		for _, g := range m.Groups {
			features += len(g.Features)
			for _, f := range g.Features {
				total += len(f.Keys)
			}
		}
	}

	return asJSON(map[string]interface{}{
		"modules":     len(modules),
		"features":    features,
		"permissions": total,
		"source":      relOrEmpty(s.Root, path),
		"catalogue":   modules,
		"note": "A role holds permission keys like products.view. A wildcard, products.*, covers " +
			"every action on that feature. A key that is not in this catalogue matches nothing and " +
			"fails silently, which is why it is worth reading this before writing a guard. Modules " +
			"and groups only shape the admin's permission tree and never appear in a key.",
	})
}

type permModule struct {
	Key    string      `json:"key"`
	Name   string      `json:"name,omitempty"`
	Groups []permGroup `json:"groups"`
}

type permGroup struct {
	Key      string        `json:"key"`
	Name     string        `json:"name,omitempty"`
	Features []permFeature `json:"features"`
}

type permFeature struct {
	Key     string   `json:"key"`
	Name    string   `json:"name,omitempty"`
	Actions []string `json:"actions"`
	Keys    []string `json:"permission_keys"`
}

// isSliceOf reports whether an expression is []Name.
func isSliceOf(expr ast.Expr, name string) bool {
	arr, ok := expr.(*ast.ArrayType)
	if !ok {
		return false
	}
	ident, ok := arr.Elt.(*ast.Ident)
	return ok && ident.Name == name
}

// fields reads a composite literal's keyed fields, whether or not its type was
// written out.
func fields(expr ast.Expr) (map[string]ast.Expr, bool) {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	out := map[string]ast.Expr{}
	for _, el := range lit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if key, ok := kv.Key.(*ast.Ident); ok {
			out[key.Name] = kv.Value
		}
	}
	return out, true
}

func readModule(expr ast.Expr) (permModule, bool) {
	f, ok := fields(expr)
	if !ok {
		return permModule{}, false
	}
	m := permModule{Key: stringLit(f["Key"]), Name: stringLit(f["Name"]), Groups: []permGroup{}}
	if groups, ok := f["Groups"].(*ast.CompositeLit); ok {
		for _, el := range groups.Elts {
			if g, ok := readGroup(el); ok {
				m.Groups = append(m.Groups, g)
			}
		}
	}
	return m, m.Key != "" || len(m.Groups) > 0
}

func readGroup(expr ast.Expr) (permGroup, bool) {
	f, ok := fields(expr)
	if !ok {
		return permGroup{}, false
	}
	g := permGroup{Key: stringLit(f["Key"]), Name: stringLit(f["Name"]), Features: []permFeature{}}
	if feats, ok := f["Features"].(*ast.CompositeLit); ok {
		for _, el := range feats.Elts {
			if feat, ok := readFeature(el); ok {
				g.Features = append(g.Features, feat)
			}
		}
	}
	return g, true
}

func readFeature(expr ast.Expr) (permFeature, bool) {
	f, ok := fields(expr)
	if !ok {
		return permFeature{}, false
	}
	feat := permFeature{Key: stringLit(f["Key"]), Name: stringLit(f["Name"])}
	if feat.Key == "" {
		return permFeature{}, false
	}
	feat.Actions = actionNames(f["Actions"])
	for _, action := range feat.Actions {
		feat.Keys = append(feat.Keys, feat.Key+"."+action)
	}
	return feat, true
}

// actionNames reads an Actions value: either a named set such as AllActions or
// a literal slice of the Action constants.
func actionNames(expr ast.Expr) []string {
	known := map[string][]string{
		"AllActions": {"create", "view", "edit", "delete"},
		"ViewOnly":   {"view"},
		"ViewEdit":   {"view", "edit"},
	}
	switch v := expr.(type) {
	case *ast.Ident:
		return known[v.Name]
	case *ast.CompositeLit:
		out := []string{}
		for _, el := range v.Elts {
			switch e := el.(type) {
			case *ast.Ident:
				out = append(out, strings.ToLower(strings.TrimPrefix(e.Name, "Action")))
			case *ast.BasicLit:
				out = append(out, stringLit(e))
			}
		}
		return out
	}
	return nil
}

func stringLit(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return ""
	}
	return s
}

// envKeys lists the variables this project reads, and whether .env sets them.
//
// No value is ever returned, from either file. .env holds the database
// password, the JWT signing key and the storage credentials, and a tool that
// returns one puts it in a transcript, a log and a model's context in a single
// call. Whether a key is set is the useful half and carries nothing.
func (s *Server) envKeys(args json.RawMessage) (string, error) {
	var a struct {
		Contains string `json:"contains"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}

	documented, err := parseEnvFile(filepath.Join(s.Root, ".env.example"))
	if err != nil {
		return "", fmt.Errorf("reading .env.example: %w", err)
	}
	// Errors ignored on purpose: a project with no .env is one nobody has
	// configured yet, which is a thing worth reporting rather than failing on.
	local, _ := parseEnvFile(filepath.Join(s.Root, ".env"))

	contains := strings.ToLower(strings.TrimSpace(a.Contains))
	keys := []map[string]interface{}{}
	for _, e := range documented {
		if contains != "" && !strings.Contains(strings.ToLower(e.Key), contains) {
			continue
		}
		set := false
		for _, l := range local {
			if l.Key == e.Key && l.HasValue {
				set = true
				break
			}
		}
		entry := map[string]interface{}{"key": e.Key, "set_in_env": set}
		if e.Comment != "" {
			entry["doc"] = e.Comment
		}
		keys = append(keys, entry)
	}

	return asJSON(map[string]interface{}{
		"count":    len(keys),
		"has_env":  fileExists(filepath.Join(s.Root, ".env")),
		"keys":     keys,
		"note":     "Values are never returned, only whether .env sets one. Read the file yourself if you need a value.",
		"commands": "grit env writes a .env from .env.example with fresh secrets.",
	})
}

type envEntry struct {
	Key      string
	Comment  string
	HasValue bool
}

// parseEnvFile reads keys and the comment block above each one. A commented-out
// key, "# MAIL_MAILER=", is a key the project reads with no default, so it is
// reported rather than treated as prose.
func parseEnvFile(path string) ([]envEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []envEntry
	var comment []string
	seen := map[string]bool{}

	for _, raw := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			comment = nil
			continue
		}
		if strings.HasPrefix(line, "#") {
			body := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			// A commented assignment is a documented key, not a sentence.
			if key, value, ok := splitEnvLine(body); ok {
				if !seen[key] {
					seen[key] = true
					out = append(out, envEntry{Key: key, Comment: strings.Join(comment, " "), HasValue: value != ""})
				}
				continue
			}
			if body != "" {
				comment = append(comment, body)
			}
			continue
		}
		key, value, ok := splitEnvLine(line)
		if !ok {
			comment = nil
			continue
		}
		if !seen[key] {
			seen[key] = true
			out = append(out, envEntry{Key: key, Comment: strings.Join(comment, " "), HasValue: value != ""})
		}
		comment = nil
	}
	return out, nil
}

// splitEnvLine reports KEY=VALUE, where the key looks like an environment
// variable rather than like the start of an English sentence.
func splitEnvLine(line string) (key, value string, ok bool) {
	at := strings.Index(line, "=")
	if at <= 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:at])
	if key == "" {
		return "", "", false
	}
	for _, r := range key {
		if r != '_' && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return "", "", false
		}
	}
	return key, strings.TrimSpace(line[at+1:]), true
}

// cliReference reports the commands of the binary running this server.
//
// From the running binary rather than from a document: an agent proposing
// `grit generate resource --embed` because some version had it costs the user a
// confusing error, and the binary is the only thing that knows what this version
// accepts.
func (s *Server) cliReference(args json.RawMessage) (string, error) {
	var a struct {
		Command string `json:"command"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}
	if len(s.Commands) == 0 {
		return "", fmt.Errorf("this server was built without a command list")
	}

	want := strings.TrimSpace(a.Command)
	out := []CommandInfo{}
	for _, c := range s.Commands {
		if want != "" && c.Path != want && !strings.HasPrefix(c.Path, want+" ") {
			continue
		}
		out = append(out, c)
	}
	if want != "" && len(out) == 0 {
		return "", fmt.Errorf("no command %q in grit %s", want, s.Version)
	}

	return asJSON(map[string]interface{}{
		"grit_version": s.Version,
		"count":        len(out),
		"commands":     out,
	})
}

// resourceNameRe is the shape of a name `grit generate resource` accepts: a
// letter, then letters and digits. Checked here rather than left to the CLI
// because the name becomes an argument to a process, and an argument built from
// a model's output is checked before it is passed, not after.
var resourceNameRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

// fieldSpecRe is the shape of an inline field list: name:type pairs, with the
// type possibly carrying its own colon-separated detail, separated by commas.
var fieldSpecRe = regexp.MustCompile(`^[A-Za-z0-9_:,.\- ]*$`)

// generateResource runs the generator. Registered only in write mode.
func (s *Server) generateResource(args json.RawMessage) (string, error) {
	var a struct {
		Name    string `json:"name"`
		Fields  string `json:"fields"`
		OwnedBy string `json:"owned_by"`
		Seed    bool   `json:"seed"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}

	name := strings.TrimSpace(a.Name)
	if !resourceNameRe.MatchString(name) {
		return "", fmt.Errorf("resource name %q is not a name: letters and digits, starting with a letter, e.g. Product", a.Name)
	}
	fields := strings.TrimSpace(a.Fields)
	if fields != "" {
		if !fieldSpecRe.MatchString(fields) {
			return "", fmt.Errorf("field list %q has characters a field list does not: use name:type pairs separated by commas", a.Fields)
		}
		// Parsed before anything runs, so a bad type is an error the agent can
		// read rather than a half-generated resource it has to clean up.
		if _, err := generate.ParseInlineFields(name, fields); err != nil {
			return "", fmt.Errorf("field list: %w", err)
		}
	}
	ownedBy := strings.TrimSpace(a.OwnedBy)
	if ownedBy != "" && !resourceNameRe.MatchString(ownedBy) {
		return "", fmt.Errorf("owned_by %q is not a field name", a.OwnedBy)
	}

	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("finding the grit binary: %w", err)
	}
	argv := []string{"generate", "resource", name}
	if fields != "" {
		argv = append(argv, "--fields", fields)
	}
	if ownedBy != "" {
		argv = append(argv, "--owned-by", ownedBy)
	}
	if a.Seed {
		argv = append(argv, "--seed")
	}

	cmd := exec.Command(exe, argv...)
	cmd.Dir = s.Root
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return "", fmt.Errorf("grit %s failed: %v\n%s", strings.Join(argv, " "), runErr, output)
	}

	return asJSON(map[string]interface{}{
		"command": "grit " + strings.Join(argv, " "),
		"output":  string(output),
		"next": "Run grit migrate to create the table, and grit_doctor to check the resource " +
			"against the project's own rules. The generated files are yours to edit.",
	})
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
