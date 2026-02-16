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
func ParseInlineFields(name string, fieldStr string) (*ResourceDefinition, error) {
	def := &ResourceDefinition{Name: name}

	parts := strings.Split(fieldStr, ",")
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

// parseFieldInput parses a field definition string.
// Format: "name:type" or "name:type:modifier1:modifier2"
// Valid modifiers: unique, required, optional
func parseFieldInput(input string) (Field, error) {
	parts := strings.Split(input, ":")
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
