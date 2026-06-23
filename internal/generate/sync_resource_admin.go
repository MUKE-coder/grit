package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SyncAdminResource walks one Go struct → its admin resource file and
// appends any model fields that aren't represented as columns or form
// fields yet.
//
// The strategy is intentionally conservative:
//   - Only INSERTS, never deletes. If you remove a field from the Go
//     model, the admin file keeps showing it (and your customisations).
//   - Only fires when a field's "key:" string isn't found anywhere in
//     the file — so an operator who hand-edited a row to change its
//     label or format never has their changes overwritten.
//   - Insertion happens inside the marker fences:
//        // grit:cols:auto-start ... // grit:cols:auto-end
//        // grit:fields:auto-start ... // grit:fields:auto-end
//     If those markers are absent (file was generated before v3.31.16),
//     this function gives up gracefully and the operator is told what
//     to add by hand. No silent corruption.
//
// Returns the number of fields added across columns + form, plus any
// human-readable warnings to surface in the CLI output.
func SyncAdminResource(root string, s GoStruct) (added int, warnings []string, err error) {
	plural := Pluralize(toSnakeCase(s.Name))
	kebab := strings.ReplaceAll(plural, "_", "-")
	path := filepath.Join(root, "apps", "admin", "resources", kebab+".ts")

	data, statErr := os.ReadFile(path)
	if os.IsNotExist(statErr) {
		// No admin resource file for this model — operator either never
		// generated one or removed it. Either way: not our problem.
		return 0, nil, nil
	}
	if statErr != nil {
		return 0, nil, fmt.Errorf("reading %s: %w", path, statErr)
	}

	content := string(data)

	// Markers came in at v3.31.16. Older files don't have them; fall
	// back to a warning so the operator knows the auto-add can't help
	// them on this resource without a small one-time edit.
	if !strings.Contains(content, "// grit:cols:auto-end") ||
		!strings.Contains(content, "// grit:fields:auto-end") {
		return 0, []string{fmt.Sprintf(
			"  ⚠ %s: missing grit:cols:auto-end / grit:fields:auto-end markers — auto-add skipped.\n     Re-run `grit generate` to regenerate the resource file, or add the markers by hand to enable sync auto-add.",
			path,
		)}, nil
	}

	updated := content
	for _, f := range s.Fields {
		fieldKey := f.JSONName
		if fieldKey == "" {
			fieldKey = toSnakeCase(f.Name)
		}
		if isAutoField(fieldKey) || f.Name == "Version" {
			continue
		}

		// Skip if already referenced anywhere in the file — covers both
		// inside-the-markers and a customised version the operator moved
		// outside the fence.
		needle := fmt.Sprintf(`key: "%s"`, fieldKey)
		if strings.Contains(updated, needle) {
			continue
		}

		colLine := buildAutoColumnLine(f, fieldKey)
		formLine := buildAutoFormFieldLine(f, fieldKey)

		updated = insertBeforeMarker(updated, "// grit:cols:auto-end", colLine)
		updated = insertBeforeMarker(updated, "// grit:fields:auto-end", formLine)
		added++
	}

	if added > 0 {
		if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
			return added, warnings, fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return added, warnings, nil
}

// insertBeforeMarker injects a new line above the given marker comment,
// preserving the marker's indentation so the result stays neatly
// aligned.
func insertBeforeMarker(content, marker, line string) string {
	idx := strings.Index(content, marker)
	if idx < 0 {
		return content // marker missing — caller already warned
	}
	// Walk back to the start of the marker's line to pick up its indent.
	lineStart := strings.LastIndex(content[:idx], "\n") + 1
	indent := content[lineStart:idx]
	insert := indent + line + "\n"
	return content[:lineStart] + insert + content[lineStart:]
}

// buildAutoColumnLine emits a sensible default { key, label, ... } column
// entry for a Go field. Best-effort — operators are expected to refine
// later (badges, formats, custom cells), but the field at least shows up
// without further action.
func buildAutoColumnLine(f GoField, key string) string {
	label := humanLabel(f.Name)
	parts := []string{
		fmt.Sprintf(`key: "%s"`, key),
		fmt.Sprintf(`label: "%s"`, label),
	}

	switch guessFormat(f) {
	case "boolean":
		parts = append(parts, `format: "boolean"`)
	case "currency":
		parts = append(parts, `format: "currency", currencyPrefix: "$"`)
	case "date":
		parts = append(parts, `format: "date"`)
	case "relative":
		parts = append(parts, `format: "relative"`)
	case "image":
		parts = append(parts, `format: "image"`)
	case "email":
		parts = append(parts, `format: "email"`)
	}

	if isSortable(f) {
		parts = append(parts, `sortable: true`)
	}
	if isSearchable(f) {
		parts = append(parts, `searchable: true`)
	}

	return "{ " + strings.Join(parts, ", ") + " },"
}

// buildAutoFormFieldLine emits a sensible default { key, label, type } form
// entry. Operators tighten the validation/help text later.
func buildAutoFormFieldLine(f GoField, key string) string {
	label := humanLabel(f.Name)
	formType := guessFormFieldType(f)
	parts := []string{
		fmt.Sprintf(`key: "%s"`, key),
		fmt.Sprintf(`label: "%s"`, label),
		fmt.Sprintf(`type: "%s"`, formType),
	}
	return "{ " + strings.Join(parts, ", ") + " },"
}

// humanLabel turns "PhoneNumber" or "phone_number" into "Phone number".
func humanLabel(name string) string {
	// Convert to space-separated words, capitalize first letter only.
	words := splitPascal(toPascalCase(name))
	if len(words) == 0 {
		return name
	}
	words[0] = strings.ToUpper(words[0][:1]) + strings.ToLower(words[0][1:])
	for i := 1; i < len(words); i++ {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, " ")
}

// guessFormat reads the field's Go type + name to suggest a column
// format. Conservative — defaults to "text" when nothing matches.
func guessFormat(f GoField) string {
	name := strings.ToLower(f.Name)

	if f.GoType == "bool" {
		return "boolean"
	}
	if strings.Contains(f.GoType, "time.Time") {
		// CreatedAt-shaped names get "relative"; everything else "date".
		if strings.HasSuffix(name, "_at") || name == "createdat" || name == "updatedat" {
			return "relative"
		}
		return "date"
	}
	if f.GoType == "float64" || f.GoType == "float32" {
		if isMoneyField(name) {
			return "currency"
		}
	}
	if f.GoType == "string" && isURLField(name) {
		// URL-named strings are often image URLs in practice.
		switch name {
		case "avatar", "photo", "image", "thumbnail", "logo", "cover", "icon", "banner":
			return "image"
		}
	}
	if name == "email" {
		return "email"
	}
	return ""
}

// guessFormFieldType picks an admin form input type from the Go type.
func guessFormFieldType(f GoField) string {
	name := strings.ToLower(f.Name)
	switch {
	case f.GoType == "bool":
		return "toggle"
	case f.GoType == "int" || f.GoType == "uint" || f.GoType == "float64" || f.GoType == "float32":
		return "number"
	case strings.Contains(f.GoType, "time.Time"):
		if strings.HasSuffix(name, "_at") {
			return "datetime"
		}
		return "date"
	case isLongTextField(name):
		return "textarea"
	default:
		return "text"
	}
}

// isSortable mirrors Field.IsSortable() but for the parser's GoField type.
func isSortable(f GoField) bool {
	switch f.GoType {
	case "string", "int", "uint", "float64", "float32":
		return true
	}
	if strings.Contains(f.GoType, "time.Time") {
		return true
	}
	return false
}

// isSearchable mirrors Field.IsSearchable() for GoField.
func isSearchable(f GoField) bool {
	return f.GoType == "string"
}
