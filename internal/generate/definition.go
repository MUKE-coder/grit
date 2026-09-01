package generate

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ResourceDefinition describes a resource to generate.
type ResourceDefinition struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`

	// Items, when set (via `--items Child:"fields"`), is an inline has-many
	// child resource. It is generated as a normal resource (so it gets its own
	// model/handler/routes and a belongs_to back to this parent), AND the
	// parent gains a has-many slice + an inline line-items form field, so its
	// rows are created atomically with the parent. Nil for ordinary resources.
	Items *ResourceDefinition `yaml:"items,omitempty"`

	// Hidden marks a resource that should be generated fully but kept OUT of the
	// admin sidebar (set on inline `--items` children — you manage them through
	// the parent's form + detail page, not a top-level nav entry).
	Hidden bool `yaml:"-"`

	// Public adds a read-only surface outside the auth middleware, guarded by
	// an API key, for clients with no logged-in user. Set by --public.
	//
	// It does not change the protected routes: the admin panel calls those,
	// and gin panics at boot on two handlers for one method and path.
	Public bool `yaml:"public,omitempty"`

	// Sync, when set, declares how the resource behaves offline: mirrored or
	// not, what happens on a version conflict, which fields cross the wire,
	// and how stale the mirror may get. Nil means the defaults every project
	// had before policies existed.
	Sync *SyncPolicy `yaml:"sync,omitempty"`

	// Tree makes the resource hierarchical: Electronics above Cameras above
	// Lenses. Set by --tree.
	//
	// It adds parent_id, a materialized path, a depth and a sibling position,
	// and generates the queries a hierarchy actually needs: the roots, the
	// whole tree in one query, a node's descendants, its breadcrumbs, and a
	// move that carries its subtree with it.
	//
	// A materialized path rather than a recursive CTE because Grit supports
	// Postgres, MySQL and SQLite: WHERE path LIKE '/1/2/%' is one indexable
	// comparison that behaves identically on all three, where CTE support and
	// syntax do not.
	Tree bool `yaml:"tree,omitempty"`
}

// TreeParentField returns the field the tree hangs from, adding it if --tree was
// passed without one.
//
// A tree needs a self-referential belongs_to, and asking somebody to write
// --fields "parent:belongs_to:Category" *and* --tree is asking them to say the
// same thing twice, with one spelling that silently produces a flat list.
func (d *ResourceDefinition) TreeParentField() *Field {
	for i := range d.Fields {
		f := &d.Fields[i]
		if f.IsBelongsTo() && f.RelatedModelName() == toPascalCase(d.Name) {
			return f
		}
	}
	return nil
}

// LoadFromYAML reads a resource definition from a YAML file.
func LoadFromYAML(path string) (*ResourceDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var def ResourceDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	if def.Name == "" {
		return nil, fmt.Errorf("resource name is required in YAML")
	}
	if len(def.Fields) == 0 {
		return nil, fmt.Errorf("at least one field is required")
	}

	for i, f := range def.Fields {
		if f.Name == "" {
			return nil, fmt.Errorf("field %d: name is required", i+1)
		}
		if f.Type == "" {
			return nil, fmt.Errorf("field %q: type is required", f.Name)
		}
		if !isValidType(f.Type) {
			return nil, fmt.Errorf("field %q: invalid type %q (valid: %s)", f.Name, f.Type, strings.Join(ValidFieldTypes(), ", "))
		}
	}

	// Validated here rather than at generation time, so a bad policy is a
	// message naming the resource and the allowed values, and no code with
	// the mistake baked in is ever written.
	if err := def.Sync.Validate(def.Name); err != nil {
		return nil, err
	}

	// Same for workflows. A state nothing can leave is invisible until a
	// record lands there in production, so it is caught here.
	for _, f := range def.Fields {
		if err := f.Workflow.Validate(def.Name, f.Name, f.OptionValues()); err != nil {
			return nil, err
		}
	}
	if unknown := def.Sync.CheckAgainstFields(def.Name, def.Fields); len(unknown) > 0 {
		return nil, fmt.Errorf(
			"sync policy on %s names field(s) the resource does not have: %s",
			def.Name, strings.Join(unknown, ", "))
	}

	return &def, nil
}

// PromptInteractive guides the user through defining fields interactively.
func PromptInteractive(name string) (*ResourceDefinition, error) {
	reader := bufio.NewReader(os.Stdin)
	def := &ResourceDefinition{Name: name}

	fmt.Println()
	fmt.Printf("  Defining fields for %s\n", name)
	fmt.Println("  Enter fields as name:type[:modifiers] (e.g., title:string, slug:slug:name)")
	fmt.Printf("  Valid types: %s\n", strings.Join(ValidFieldTypes(), ", "))
	fmt.Println("  Valid modifiers: unique, required, optional")
	fmt.Println("  Slug fields: slug:slug (auto-detect source) or slug:slug:name (explicit source)")
	fmt.Println("  Relationships: category:belongs_to or author:belongs_to:User, tags:many_to_many:Tag")
	fmt.Println("  Press Enter with no input when done.")
	fmt.Println()

	for {
		fmt.Print("  > ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("reading input: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		field, err := parseFieldInput(line)
		if err != nil {
			fmt.Printf("  ⚠ %s\n", err)
			continue
		}

		def.Fields = append(def.Fields, field)
		fmt.Printf("  ✓ Added %s (%s)\n", field.Name, field.Type)
	}

	if len(def.Fields) == 0 {
		return nil, fmt.Errorf("at least one field is required")
	}

	return def, nil
}

// ParseInlineFields parses a comma-separated list of field definitions.
// Format: "title:string,content:text,published:bool"
//
// v3.31.30: bracket-aware. File-field type lists use brackets to scope
// commas so they don't collide with the top-level field separator:
//
//	"image:file:image,attachments:files:[pdf,doc,image]"
//
// Without bracket awareness, the `pdf,doc,image` would split across
// three "fields" and produce nonsense.
func ParseInlineFields(name string, fieldStr string) (*ResourceDefinition, error) {
	def := &ResourceDefinition{Name: name}

	parts := splitTopLevelCommas(fieldStr)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field, err := parseFieldInput(part)
		if err != nil {
			return nil, err
		}
		def.Fields = append(def.Fields, field)
	}

	if len(def.Fields) == 0 {
		return nil, fmt.Errorf("at least one field is required")
	}

	return def, nil
}

// ParseItems parses the `--items Child:"field:type,..."` spec into an inline
// child ResourceDefinition owned by the parent. The child gets a synthesized
// belongs_to back to the parent (so its FK column + list filter exist) and is
// flagged Hidden (generated fully, but kept out of the sidebar). Example:
//
//	ParseItems("Invoice", `InvoiceItem:description:string,qty:int,unit_rate:float`)
func ParseItems(parentName, spec string) (*ResourceDefinition, error) {
	spec = strings.TrimSpace(spec)
	// Split off the child name (everything up to the first colon); the rest is
	// the field spec passed to ParseInlineFields.
	idx := strings.Index(spec, ":")
	if idx < 0 {
		return nil, fmt.Errorf("--items must be \"ChildName:field:type,...\" (got %q)", spec)
	}
	childName := strings.TrimSpace(spec[:idx])
	fieldStr := strings.TrimSpace(spec[idx+1:])
	if childName == "" || fieldStr == "" {
		return nil, fmt.Errorf("--items must be \"ChildName:field:type,...\" (got %q)", spec)
	}

	child, err := ParseInlineFields(childName, fieldStr)
	if err != nil {
		return nil, fmt.Errorf("parsing --items fields: %w", err)
	}

	// Synthesize the belongs_to back-link to the parent so the child owns a
	// <parent>_id FK column and a matching list filter.
	parentPascal := toPascalCase(parentName)
	child.Fields = append(child.Fields, Field{
		Name:         toSnakeCase(parentPascal),
		Type:         string(FieldBelongsTo),
		Required:     true,
		RelatedModel: parentPascal,
	})
	child.Hidden = true

	return child, nil
}

// splitTopLevelCommas splits a string on commas that are NOT inside
// square brackets. Used so file-field type lists like [pdf,doc,image]
// stay glued together.
func splitTopLevelCommas(s string) []string {
	var out []string
	var buf strings.Builder
	depth := 0
	for _, r := range s {
		switch r {
		case '[':
			depth++
			buf.WriteRune(r)
		case ']':
			if depth > 0 {
				depth--
			}
			buf.WriteRune(r)
		case ',':
			if depth == 0 {
				out = append(out, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(r)
			}
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		out = append(out, buf.String())
	}
	return out
}

// parseFieldInput parses a field definition string.
// Format: "name:type" or "name:type:modifier1:modifier2"
// Valid modifiers: unique, required, optional
//
// v3.31.30: file/files type. Third part is the accept-list (single
// alias or bracketed alias list):
//
//	image:file:image
//	attachment:file:all
//	attachment:file:[pdf,doc,image,video,zip]
//	gallery_images:files:image
func parseFieldInput(input string) (Field, error) {
	parts := splitFieldParts(input)
	if len(parts) < 2 {
		return Field{}, fmt.Errorf("expected format name:type[:modifiers], got %q", input)
	}

	name := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(parts[1])

	if name == "" {
		return Field{}, fmt.Errorf("field name cannot be empty")
	}
	if !isValidType(typ) {
		return Field{}, fmt.Errorf("invalid type %q for field %q (valid: %s)", typ, name, strings.Join(ValidFieldTypes(), ", "))
	}

	// file / files: third part is the accept-list. Bare alias for single
	// type (image:file:image) or bracketed list for multi
	// (attachments:files:[pdf,doc,image]). Defaults to "all" if omitted.
	if typ == "file" || typ == "files" {
		acceptStr := "all"
		if len(parts) >= 3 && strings.TrimSpace(parts[2]) != "" {
			acceptStr = strings.TrimSpace(parts[2])
		}
		accepts, err := parseFileAccepts(acceptStr)
		if err != nil {
			return Field{}, fmt.Errorf("field %q: %w", name, err)
		}
		return Field{
			Name:        name,
			Type:        typ,
			Required:    false,
			FileAccepts: accepts,
		}, nil
	}

	// Slug fields: third part is the source field name, not a modifier
	if typ == "slug" {
		slugSource := ""
		if len(parts) >= 3 && strings.TrimSpace(parts[2]) != "" {
			slugSource = strings.TrimSpace(parts[2])
		}
		return Field{
			Name:       name,
			Type:       typ,
			Required:   false,
			Unique:     true,
			SlugSource: slugSource,
		}, nil
	}

	// belongs_to and one_to_one: third part is the related model name, optional
	// and inferred from the field name.
	//   category:belongs_to        → Category
	//   author:belongs_to:User     → User
	//   user:one_to_one:User       → User, with a unique foreign key
	//
	// Parsed together because they are the same shape. They differ only in the
	// GORM tag, where one_to_one gets uniqueIndex.
	if typ == "belongs_to" || typ == "one_to_one" {
		relatedModel := ""
		if len(parts) >= 3 && strings.TrimSpace(parts[2]) != "" {
			relatedModel = strings.TrimSpace(parts[2])
		}
		return Field{
			Name:         name,
			Type:         typ,
			Required:     true,
			RelatedModel: relatedModel,
		}, nil
	}

	// many_to_many: third part is the related model name (required)
	// e.g., tags:many_to_many:Tag
	if typ == "many_to_many" {
		if len(parts) < 3 || strings.TrimSpace(parts[2]) == "" {
			return Field{}, fmt.Errorf("many_to_many requires a related model name (e.g., tags:many_to_many:Tag)")
		}
		relatedModel := strings.TrimSpace(parts[2])
		return Field{
			Name:         name,
			Type:         typ,
			Required:     false,
			RelatedModel: relatedModel,
		}, nil
	}

	// select / check: third part is the value=Label option list, pipe-separated
	// (e.g. status:select:draft=Draft|sent=Sent|paid=Paid). A bare value gets a
	// humanized label.
	if typ == "select" || typ == "radio" || typ == "check" {
		if len(parts) < 3 || strings.TrimSpace(parts[2]) == "" {
			return Field{}, fmt.Errorf("%s field %q needs options, e.g. %s:%s:draft=Draft|sent=Sent", typ, name, name, typ)
		}
		options, err := parseFieldOptions(strings.TrimSpace(parts[2]))
		if err != nil {
			return Field{}, fmt.Errorf("field %q: %w", name, err)
		}
		return Field{
			Name:     name,
			Type:     typ,
			Required: typ == "select" || typ == "radio", // single-choice fields default to required
			Options:  options,
		}, nil
	}

	// Default: string fields are required
	required := typ == "string"
	unique := false
	auto := false
	autoPrefix := ""

	// Parse optional modifiers (parts[2], parts[3], etc.). Index-based so `auto`
	// can consume an optional prefix that follows it (name:string:auto:INV).
	mods := parts[2:]
	for i := 0; i < len(mods); i++ {
		mod := strings.TrimSpace(mods[i])
		switch strings.ToLower(mod) {
		case "unique":
			unique = true
		case "required":
			required = true
		case "optional":
			required = false
		case "auto":
			// Auto-generated from a sequence in BeforeCreate. Never user-required.
			auto = true
			required = false
			// An immediately following segment is the sequence prefix (case-preserved).
			if i+1 < len(mods) && strings.TrimSpace(mods[i+1]) != "" {
				autoPrefix = strings.TrimSpace(mods[i+1])
				i++
			}
		case "":
			// ignore empty modifiers
		default:
			return Field{}, fmt.Errorf("invalid modifier %q for field %q (valid: unique, required, optional, auto)", mod, name)
		}
	}

	if auto && typ != "string" {
		return Field{}, fmt.Errorf("field %q: auto is only valid on string fields (got %q)", name, typ)
	}

	return Field{
		Name:       name,
		Type:       typ,
		Required:   required,
		Unique:     unique,
		Auto:       auto,
		AutoPrefix: autoPrefix,
	}, nil
}

// parseFieldOptions turns "draft=Draft|sent=Sent|paid=Paid" into options. A bare
// value ("draft") is kept as the value with a humanized label ("Draft").
func parseFieldOptions(s string) ([]FieldOption, error) {
	var out []FieldOption
	seen := map[string]bool{}
	for _, pair := range strings.Split(s, "|") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		value, label := pair, ""
		if i := strings.Index(pair, "="); i >= 0 {
			value = strings.TrimSpace(pair[:i])
			label = strings.TrimSpace(pair[i+1:])
		}
		if value == "" {
			return nil, fmt.Errorf("empty option value in %q", s)
		}
		if seen[value] {
			return nil, fmt.Errorf("duplicate option value %q", value)
		}
		seen[value] = true
		if label == "" {
			label = humanizeLabel(value)
		}
		out = append(out, FieldOption{Value: value, Label: label})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no options parsed from %q", s)
	}
	return out, nil
}

func isValidType(t string) bool {
	for _, valid := range ValidFieldTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// splitFieldParts splits a field definition on colons, BUT keeps anything
// inside square brackets as a single token. Needed because file/files
// types use brackets to scope accept-lists:
//
//	"attachments:files:[pdf,doc,image]" → ["attachments", "files", "[pdf,doc,image]"]
//
// Without bracket awareness the third part would just be "[pdf" and the
// rest of the accept-list would split into invalid parts.
func splitFieldParts(input string) []string {
	var out []string
	var buf strings.Builder
	depth := 0
	for _, r := range input {
		switch r {
		case '[':
			depth++
			buf.WriteRune(r)
		case ']':
			if depth > 0 {
				depth--
			}
			buf.WriteRune(r)
		case ':':
			if depth == 0 {
				out = append(out, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(r)
			}
		default:
			buf.WriteRune(r)
		}
	}
	out = append(out, buf.String())
	return out
}

// validFileAccepts is the set of accept-aliases recognised by the file
// CLI syntax. These are HIGH-LEVEL aliases — the scaffolded API
// translates them to concrete MIME types at request time.
var validFileAccepts = map[string]bool{
	"image":   true,
	"video":   true,
	"audio":   true,
	"pdf":     true,
	"doc":     true,
	"excel":   true,
	"csv":     true,
	"zip":     true,
	"archive": true,
	"all":     true,
}

// parseFileAccepts parses the accept-list of a file/files field type.
// Accepts a bare alias ("image", "all", ...) or a bracketed list
// ("[pdf,doc,image]"). Returns the normalised lowercase aliases.
func parseFileAccepts(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("file accept-list cannot be empty")
	}
	// Strip outer brackets if present.
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		s = s[1 : len(s)-1]
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(p))
		if p == "" {
			continue
		}
		if !validFileAccepts[p] {
			valids := []string{"image", "video", "audio", "pdf", "doc", "excel", "csv", "zip", "archive", "all"}
			return nil, fmt.Errorf("invalid file accept alias %q (valid: %s)", p, strings.Join(valids, ", "))
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("file accept-list cannot be empty")
	}
	return out, nil
}
