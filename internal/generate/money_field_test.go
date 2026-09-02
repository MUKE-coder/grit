package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A money field has to become two database columns. Everything else about the
// type is downstream of that: an amount the database can SUM, and a currency
// you can filter by.
func TestMoneyFieldTypeMapping(t *testing.T) {
	f := Field{Name: "Price", Type: string(FieldMoney)}

	if got := f.GoType(); got != "money.Money" {
		t.Errorf("GoType() = %q, want money.Money", got)
	}
	// The embedded tag comes from the model template rather than GORMTag, so
	// this stays empty on purpose.
	if got := f.GORMTag(); got != "" {
		t.Errorf("GORMTag() = %q, want empty", got)
	}
}

// A money field has to reach TypeScript as the object, not as a bare number or
// a string. The currency travelling with the amount is the whole point.
func TestMoneyFrontendMapping(t *testing.T) {
	f := Field{Name: "Price", Type: string(FieldMoney)}

	if got := f.TSType(); got != "Money" {
		t.Errorf("TSType() = %q, want Money", got)
	}
	if got := f.ZodType(); !strings.Contains(got, "MoneySchema") {
		t.Errorf("ZodType() = %q, want it to use MoneySchema", got)
	}
	if got := f.FormFieldType(); got != "money" {
		t.Errorf("FormFieldType() = %q, want money", got)
	}
	if got := f.ColumnFormat(); got != "money" {
		t.Errorf("ColumnFormat() = %q, want money", got)
	}
	if !f.IsSortable() {
		t.Error("IsSortable() = false; a price column nobody can sort is a table nobody can use")
	}
}

// End to end through the generator, because the parts that broke in practice
// were the ones no unit test covered: the embedded tag on the struct field,
// and the sort and filter whitelists.
func TestMoneyGeneratesEmbeddedColumnsAndWhitelists(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Product", "name:string,price:money,weight:float")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	// The struct field. Without embeddedPrefix the two halves land as
	// "amount" and "currency", which collide the moment a model has a second
	// money field.
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "product.go"))
	assertContains(t, "model", model,
		"money.Money",
		"embedded;embeddedPrefix:price_",
		module+"/internal/money",
	)
	// A float field alongside it still maps to a plain column. Matched loosely
	// because gofmt pads struct fields into columns.
	if !strings.Contains(model, "Weight") || !strings.Contains(model, "float64") {
		t.Error("model lost the plain float field")
	}

	// The handler interpolates these names straight into ORDER BY and WHERE.
	// "price" is not a column, and asking for it is a 500.
	handler := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "product.go"))
	assertContains(t, "handler", handler, `"price_amount": true`, `"price_currency": true`)
	if strings.Contains(handler, `"price": true`) {
		t.Error(`handler whitelists "price", which is not a column on the table`)
	}

	// The shared type and schema name Money, so they have to import it.
	types := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "product.ts"))
	assertContains(t, "types", types, "price: Money;", `import type { Money } from "./money";`)

	schema := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "product.ts"))
	assertContains(t, "schema", schema, "MoneySchema", `from "./money"`)
}
