package generate

import (
	"fmt"
	"os"
	"strings"
)

// pickLabelExpr returns a Go expression that evaluates to a
// human-readable label for a newly-created record, used in the public
// form-share submission response. Best-effort:
//   - if the model has a Name string, use item.Name
//   - else if it has Title, use item.Title
//   - else if it has Subject, use item.Subject
//   - else fall back to item.ID (always safe — every Grit model has it)
//
// The label is intentionally non-sensitive (no email, no password
// fields) so it's safe to echo back to an anonymous visitor.
func pickLabelExpr(fields []Field) string {
	have := map[string]bool{}
	for _, f := range fields {
		have[toSnakeCase(f.Name)] = true
	}
	switch {
	case have["name"]:
		return "item.Name"
	case have["title"]:
		return "item.Title"
	case have["subject"]:
		return "item.Subject"
	default:
		return "item.ID"
	}
}

// ensureDispatchImports patches form_share_dispatch.go so the new
// dispatch case compiles. Each generated case uses encoding/json
// (for the map → struct conversion) and the project's models package
// (referenced as models.Pascal). The base file ships with neither —
// adding them lazily keeps the file lean until at least one resource
// has been generated.
func ensureDispatchImports(path, module string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)

	jsonImport := `"encoding/json"`
	modelsImport := fmt.Sprintf(`"%s/internal/models"`, module)
	// v3.31.43: PublicFields uses reflection + strings.HasSuffix etc.
	// Older projects scaffolded before v3.31.43 don't have these
	// imports yet; add them lazily when we inject the first case.
	reflectImport := `"reflect"`
	stringsImport := `"strings"`

	needJSON := !strings.Contains(content, jsonImport)
	needModels := !strings.Contains(content, modelsImport)
	needReflect := !strings.Contains(content, reflectImport)
	needStrings := !strings.Contains(content, stringsImport)
	if !needJSON && !needModels && !needReflect && !needStrings {
		return nil
	}

	// Find the existing import block. The base file imports "fmt" and
	// "gorm.io/gorm" between import ( … ). We append inside that block.
	const openMarker = "import ("
	openIdx := strings.Index(content, openMarker)
	if openIdx < 0 {
		return fmt.Errorf("import ( not found in %s — file structure unexpected", path)
	}
	closeIdx := strings.Index(content[openIdx:], ")")
	if closeIdx < 0 {
		return fmt.Errorf("import block close ) not found in %s", path)
	}
	closeIdx += openIdx

	additions := ""
	if needJSON {
		additions += "\n\t" + jsonImport
	}
	if needReflect {
		additions += "\n\t" + reflectImport
	}
	if needStrings {
		additions += "\n\t" + stringsImport
	}
	if needModels {
		additions += "\n\n\t" + modelsImport
	}

	updated := content[:closeIdx] + additions + "\n" + content[closeIdx:]
	return os.WriteFile(path, []byte(updated), 0644)
}
