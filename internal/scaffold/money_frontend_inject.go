package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ensureMoneyFrontend teaches an existing project's frontend about the money
// field type.
//
// The Go half of a money field arrives on upgrade for free, because
// internal/money is a new package and new packages just get written. The
// frontend half does not: the shared Money type, the table cell that formats
// it, and the form field that edits it all live in files an upgrade
// deliberately leaves alone, because people edit them.
//
// Without this, `grit generate resource Product --fields price:money` in an
// upgraded project produces a model that imports a Money type nothing exports,
// a schema that references a MoneySchema nothing defines, a table column that
// renders [object Object], and a form field that falls through to a text box.
// That is the same half-delivered shape media, recovery contacts and passkeys
// each arrived in, so it gets the same treatment: an injector that runs on
// upgrade and does nothing on a project that already has the pieces.
//
// Every step is a no-op when its marker is already present, so running an
// upgrade twice is safe.
func ensureMoneyFrontend(root string, opts Options) error {
	if !opts.ShouldIncludeAdmin() {
		return nil
	}

	if err := ensureSharedMoneyModule(filepath.Join(root, "packages", "shared")); err != nil {
		return err
	}

	adminRoot := filepath.Join(root, "apps", "admin")
	// TanStack keeps the same tree one level down, under src/.
	base := adminRoot
	if opts.UseTanStack() {
		base = filepath.Join(adminRoot, "src")
	}

	if err := ensureMoneyFieldComponent(base, opts); err != nil {
		return err
	}
	if err := ensureMoneyFormBuilderCase(filepath.Join(base, "components", "forms", "form-builder.tsx")); err != nil {
		return err
	}
	if err := ensureMoneyCellRenderer(filepath.Join(base, "components", "tables", "cell-renderers.tsx")); err != nil {
		return err
	}
	return ensureMoneyResourceTypes(filepath.Join(base, "lib", "resource.ts"))
}

// ensureSharedMoneyModule writes the shared schema and type, and adds them to
// the two barrels.
func ensureSharedMoneyModule(sharedRoot string) error {
	if _, err := os.Stat(sharedRoot); err != nil {
		return nil
	}

	files := map[string]string{
		filepath.Join(sharedRoot, "schemas", "money.ts"): sharedMoneySchema(),
		filepath.Join(sharedRoot, "types", "money.ts"):   sharedMoneyTypes(),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	if err := appendExport(
		filepath.Join(sharedRoot, "schemas", "index.ts"),
		`export { MoneySchema, type Money } from "./money";`,
		"MoneySchema",
	); err != nil {
		return err
	}
	return appendExport(
		filepath.Join(sharedRoot, "types", "index.ts"),
		"export {\n  type Money,\n  currencyExponent,\n  toMajor,\n  fromMajor,\n  formatMoney,\n  zeroMoney,\n} from \"./money\";",
		"formatMoney",
	)
}

// appendExport adds a line to a barrel file unless the marker is already there.
func appendExport(path, line, marker string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		// No barrel means no project to wire; the scaffold path writes it.
		return nil
	}
	content := string(body)
	if strings.Contains(content, marker) {
		return nil
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content+line+"\n"), 0o644)
}

func ensureMoneyFieldComponent(base string, opts Options) error {
	path := filepath.Join(base, "components", "forms", "fields", "money-field.tsx")
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		return nil
	}
	src := adminMoneyField()
	if opts.UseTanStack() {
		src = nextToTanStack(src)
	}
	return writeFile(path, src)
}

func ensureMoneyFormBuilderCase(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(body)
	if strings.Contains(content, "MoneyField") {
		return nil
	}

	const importAnchor = `import { NumberField } from "./fields/number-field";`
	if !strings.Contains(content, importAnchor) {
		fmt.Println("  ⚠ form-builder.tsx looks hand-edited; add the money case by hand")
		return nil
	}
	content = strings.Replace(content, importAnchor,
		importAnchor+"\n"+`import { MoneyField } from "./fields/money-field";`, 1)

	// Inserted before the select case rather than appended, because the switch
	// ends in a default that would swallow anything after it.
	const caseAnchor = "    case \"select\":"
	if !strings.Contains(content, caseAnchor) {
		fmt.Println("  ⚠ form-builder.tsx switch not recognised; add the money case by hand")
		return nil
	}
	moneyCase := "    case \"money\":\n" +
		"      return (\n" +
		"        <Controller\n" +
		"          name={field.key}\n" +
		"          control={control}\n" +
		"          rules={field.required ? { required: `${field.label} is required` } : undefined}\n" +
		"          render={({ field: formField }) => (\n" +
		"            <MoneyField field={field} value={formField.value ?? null} onChange={formField.onChange} error={error} />\n" +
		"          )}\n" +
		"        />\n" +
		"      );\n"
	content = strings.Replace(content, caseAnchor, moneyCase+caseAnchor, 1)

	return os.WriteFile(path, []byte(content), 0o644)
}

func ensureMoneyCellRenderer(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(body)
	if strings.Contains(content, "MoneyCell") {
		return nil
	}

	const importAnchor = `import { formatDate, formatRelative, formatCurrency } from "@/lib/formatters";`
	if !strings.Contains(content, importAnchor) {
		fmt.Println("  ⚠ cell-renderers.tsx looks hand-edited; add the money cell by hand")
		return nil
	}
	content = strings.Replace(content, importAnchor,
		importAnchor+"\n"+`import { formatMoney, type Money } from "@repo/shared/types";`, 1)

	const caseAnchor = "    case \"date\":"
	if !strings.Contains(content, caseAnchor) {
		fmt.Println("  ⚠ cell-renderers.tsx switch not recognised; add the money cell by hand")
		return nil
	}
	moneyCase := "    case \"money\":\n" +
		"      content = <MoneyCell value={value as Money | null} />;\n" +
		"      break;\n"
	content = strings.Replace(content, caseAnchor, moneyCase+caseAnchor, 1)

	const cellAnchor = "function DateCell({ value }: { value: string }) {"
	if !strings.Contains(content, cellAnchor) {
		fmt.Println("  ⚠ cell-renderers.tsx has no DateCell; add MoneyCell by hand")
		return nil
	}
	cell := "function MoneyCell({ value }: { value: Money | null }) {\n" +
		"  return (\n" +
		"    <span className=\"block text-right font-mono text-sm tabular-nums\">\n" +
		"      {formatMoney(value)}\n" +
		"    </span>\n" +
		"  );\n" +
		"}\n\n"
	content = strings.Replace(content, cellAnchor, cell+cellAnchor, 1)

	return os.WriteFile(path, []byte(content), 0o644)
}

// ensureMoneyResourceTypes widens the two unions a generated resource writes
// into. Without them the resource file is valid at runtime and rejected by
// tsc, which is the worse of the two failure modes: the admin runs in dev and
// the build fails in CI.
func ensureMoneyResourceTypes(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(body)
	changed := false

	if !strings.Contains(content, `"currency" | "money"`) {
		if strings.Contains(content, `"badge" | "currency" |`) {
			content = strings.Replace(content, `"badge" | "currency" |`, `"badge" | "currency" | "money" |`, 1)
			changed = true
		}
	}
	if !strings.Contains(content, `"number" | "money"`) {
		if strings.Contains(content, `"number" | "select"`) {
			content = strings.Replace(content, `"number" | "select"`, `"number" | "money" | "select"`, 1)
			changed = true
		}
	}
	if !strings.Contains(content, "currencies?: string[];") {
		const anchor = `  numberKind?: "int" | "uint" | "float";`
		if strings.Contains(content, anchor) {
			content = strings.Replace(content, anchor,
				anchor+"\n\n"+
					"  /** Currency codes a money field offers. Unset shows a short default list. */\n"+
					"  currencies?: string[];\n"+
					"  /** Which of them a new record starts on. Unset uses the first. */\n"+
					"  defaultCurrency?: string;", 1)
			changed = true
		}
	}

	if !changed {
		return nil
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
