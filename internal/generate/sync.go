package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Sync parses Go model files and regenerates TypeScript types and Zod schemas.
func Sync() error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}

	modelsDir := filepath.Join(root, "apps", "api", "internal", "models")
	if !dirExists(modelsDir) {
		return fmt.Errorf("models directory not found at %s", modelsDir)
	}

	sharedRoot := filepath.Join(root, "packages", "shared")
	if !dirExists(sharedRoot) {
		return fmt.Errorf("shared package not found at %s — sync requires a monorepo project", sharedRoot)
	}

	fmt.Print("\n  Syncing Go types → TypeScript...\n\n")

	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return fmt.Errorf("reading models directory: %w", err)
	}

	var synced int
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(modelsDir, entry.Name())
		structs, err := parseGoStructs(filePath)
		if err != nil {
			fmt.Printf("  ⚠ Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		for _, s := range structs {
			// Skip User — it has custom schemas already
			if s.Name == "User" {
				continue
			}

			kebab := strings.ReplaceAll(toSnakeCase(s.Name), "_", "-")

			// Write TypeScript type
			tsContent := buildTSType(s)
			tsPath := filepath.Join(sharedRoot, "types", kebab+".ts")
			if err := writeFileWithDirs(tsPath, tsContent); err != nil {
				return fmt.Errorf("writing %s: %w", tsPath, err)
			}
			fmt.Printf("  ✓ packages/shared/types/%s.ts\n", kebab)

			// Write Zod schema
			zodContent := buildZodSchema(s)
			zodPath := filepath.Join(sharedRoot, "schemas", kebab+".ts")
			if err := writeFileWithDirs(zodPath, zodContent); err != nil {
				return fmt.Errorf("writing %s: %w", zodPath, err)
			}
			fmt.Printf("  ✓ packages/shared/schemas/%s.ts\n", kebab)

			synced++
		}
	}

	if synced == 0 {
		fmt.Println("  No models found to sync (User is skipped).")
	} else {
		fmt.Printf("\n  ✅ Synced %d model(s) to TypeScript + Zod\n\n", synced)
	}

	return nil
}

// GoStruct represents a parsed Go struct.
type GoStruct struct {
	Name   string
	Fields []GoField
}

// GoField represents a parsed Go struct field.
type GoField struct {
	Name     string
	GoType   string
	JSONName string
	GORMTag  string
}

// parseGoStructs parses all struct declarations from a Go file.
func parseGoStructs(filePath string) ([]GoStruct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", filePath, err)
	}

	var structs []GoStruct

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			s := GoStruct{Name: typeSpec.Name.Name}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue // embedded field
				}

				goField := GoField{
					Name:   field.Names[0].Name,
					GoType: typeToString(field.Type),
				}

				// Parse struct tags
				if field.Tag != nil {
					tag := field.Tag.Value
					tag = strings.Trim(tag, "`")

					goField.JSONName = extractTag(tag, "json")
					goField.GORMTag = extractTag(tag, "gorm")
				}

				// Skip fields with json:"-"
				if goField.JSONName == "-" {
					continue
				}

				// Default JSON name to snake_case of Go name
				if goField.JSONName == "" {
					goField.JSONName = toSnakeCase(goField.Name)
				}

				s.Fields = append(s.Fields, goField)
			}

			structs = append(structs, s)
		}
	}

	return structs, nil
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.IndexExpr:
		return typeToString(t.X) + "[" + typeToString(t.Index) + "]"
	default:
		return "interface{}"
	}
}

func extractTag(tag, key string) string {
	search := key + `:"`
	idx := strings.Index(tag, search)
	if idx == -1 {
		return ""
	}
	rest := tag[idx+len(search):]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return ""
	}
	value := rest[:end]
	// json:"name,omitempty" → just "name"
	if comma := strings.Index(value, ","); comma != -1 {
		value = value[:comma]
	}
	return value
}

func buildTSType(s GoStruct) string {
	var fields strings.Builder
	for _, f := range s.Fields {
		tsType := goTypeToTS(f.GoType)
		fields.WriteString(fmt.Sprintf("  %s: %s;\n", f.JSONName, tsType))
	}

	return fmt.Sprintf("export interface %s {\n%s}\n", s.Name, fields.String())
}

func buildZodSchema(s GoStruct) string {
	var createFields strings.Builder
	var updateFields strings.Builder

	for _, f := range s.Fields {
		// Skip auto-generated fields
		if isAutoField(f.JSONName) {
			continue
		}

		zodType := goTypeToZod(f.GoType, f.GORMTag)
		camelName := toCamelCase(f.JSONName)

		createFields.WriteString(fmt.Sprintf("  %s: %s,\n", camelName, zodType))

		updateZod := zodType
		if !strings.Contains(updateZod, ".optional()") && !strings.Contains(updateZod, ".nullable()") {
			updateZod += ".optional()"
		}
		updateFields.WriteString(fmt.Sprintf("  %s: %s,\n", camelName, updateZod))
	}

	return fmt.Sprintf(`import { z } from "zod";

export const Create%sSchema = z.object({
%s});

export const Update%sSchema = z.object({
%s});

export type Create%sInput = z.infer<typeof Create%sSchema>;
export type Update%sInput = z.infer<typeof Update%sSchema>;
`, s.Name, createFields.String(), s.Name, updateFields.String(),
		s.Name, s.Name, s.Name, s.Name)
}

func goTypeToTS(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "*time.Time":
		return "string | null"
	case "time.Time":
		return "string"
	case "gorm.DeletedAt":
		return "string | null"
	default:
		if strings.Contains(goType, "JSONSlice[string]") {
			return "string[]"
		}
		if strings.HasPrefix(goType, "[]") {
			inner := goTypeToTS(goType[2:])
			return inner + "[]"
		}
		if strings.HasPrefix(goType, "*") {
			inner := goTypeToTS(goType[1:])
			return inner + " | null"
		}
		return "unknown"
	}
}

func goTypeToZod(goType, gormTag string) string {
	switch goType {
	case "string":
		if strings.Contains(gormTag, "type:text") {
			return "z.string()"
		}
		return "z.string()"
	case "int", "int8", "int16", "int32", "int64":
		return "z.number().int()"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "z.number().int().nonnegative()"
	case "float32", "float64":
		return "z.number()"
	case "bool":
		return "z.boolean()"
	case "*time.Time":
		return "z.string().nullable()"
	case "time.Time":
		return "z.string()"
	default:
		if strings.Contains(goType, "JSONSlice[string]") {
			return "z.array(z.string()).optional()"
		}
		return "z.unknown()"
	}
}

func isAutoField(jsonName string) bool {
	auto := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
	}
	return auto[jsonName]
}
