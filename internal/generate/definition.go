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
//   "image:file:image,attachments:files:[pdf,doc,image]"
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
//   image:file:image
//   attachment:file:all
//   attachment:file:[pdf,doc,image,video,zip]
//   gallery_images:files:image
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

	// belongs_to: third part is the related model name (optional, inferred from field name)
	// e.g., category:belongs_to → Category, author:belongs_to:User → User
	if typ == "belongs_to" {
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

	// Default: string fields are required
	required := typ == "string"
	unique := false

	// Parse optional modifiers (parts[2], parts[3], etc.)
	for _, mod := range parts[2:] {
		mod = strings.TrimSpace(strings.ToLower(mod))
		switch mod {
		case "unique":
			unique = true
		case "required":
			required = true
		case "optional":
			required = false
		case "":
			// ignore empty modifiers
		default:
			return Field{}, fmt.Errorf("invalid modifier %q for field %q (valid: unique, required, optional)", mod, name)
		}
	}

	return Field{
		Name:     name,
		Type:     typ,
		Required: required,
		Unique:   unique,
	}, nil
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
//   "attachments:files:[pdf,doc,image]" → ["attachments", "files", "[pdf,doc,image]"]
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
