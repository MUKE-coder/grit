package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// writeModel drops a model file shaped like the one the generator writes.
func writeModel(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "product.go")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("writing the model: %v", err)
	}
	return path
}

// modelShape has to find the slug field whatever it was named.
//
// It used to look for a field literally called Slug. A shop calls it handle, so
// a Product with handle:slug read as having no slug at all: the public variants
// endpoint was generated to look up by id while its sibling
// /public/products/:key looks up by handle, which is a 404 on the detail page
// of every shop anyone would build.
func TestModelShapeFindsTheSlugFieldByItsRealName(t *testing.T) {
	cases := []struct {
		name  string
		model string
		want  string
	}{
		{
			name: "handle, which is what a shop calls it",
			model: `package models

type Product struct {
	ID     string ` + "`json:\"id\"`" + `
	Title  string ` + "`json:\"title\"`" + `
	Handle string ` + "`gorm:\"size:255;uniqueIndex\" json:\"handle\"`" + `
}

func (m *Product) BeforeCreate(tx *gorm.DB) error {
	if m.Handle == "" {
		m.Handle = slugify(fmt.Sprintf("%v", m.Title))
	}
	return nil
}
`,
			want: "Handle",
		},
		{
			name: "slug, the name the old check assumed",
			model: `package models

type Product struct {
	Slug string ` + "`json:\"slug\"`" + `
}

func (m *Product) BeforeCreate(tx *gorm.DB) error {
	if m.Slug == "" {
		m.Slug = slugify(fmt.Sprintf("%v", m.Name))
	}
	return nil
}
`,
			want: "Slug",
		},
		{
			name: "a two-word name",
			model: `package models

type Product struct {
	PublicPath string ` + "`json:\"public_path\"`" + `
}

func (m *Product) BeforeCreate(tx *gorm.DB) error {
	m.PublicPath = slugify(fmt.Sprintf("%v", m.Title))
	return nil
}
`,
			want: "PublicPath",
		},
		{
			name: "no slug at all",
			model: `package models

type Product struct {
	ID    string ` + "`json:\"id\"`" + `
	Title string ` + "`json:\"title\"`" + `
}
`,
			want: "",
		},
		{
			name: "a unique string column is not a slug",
			model: `package models

type Product struct {
	SKU string ` + "`gorm:\"size:64;uniqueIndex\" json:\"sku\"`" + `
}
`,
			want: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := modelShape(writeModel(t, c.model))
			if got != c.want {
				t.Errorf("slug field %q, want %q", got, c.want)
			}
		})
	}
}

func TestModelShapeReadsArchivedAt(t *testing.T) {
	with := `package models

type Product struct {
	ArchivedAt *time.Time ` + "`json:\"archived_at\"`" + `
}
`
	if _, archived := modelShape(writeModel(t, with)); !archived {
		t.Error("a model with ArchivedAt read as having none")
	}
	if _, archived := modelShape(writeModel(t, "package models\n\ntype Product struct{}\n")); archived {
		t.Error("a model without ArchivedAt read as having one")
	}
}

// The slug field's name has to reach the generated SQL and the generated SKUs,
// which is the whole reason for finding it.
func TestVariantGeneratorsUseTheSlugColumnTheyWereGiven(t *testing.T) {
	const module = "shop/apps/api"

	public := scaffold.APIVariantPublicGo(module, "Product", "product", "products", "Handle", true, true)
	if !strings.Contains(public, `Where("handle = ? OR id = ?"`) {
		t.Error("the public variants endpoint does not look up by handle, so a detail " +
			"page routed on the handle gets a 404")
	}
	if strings.Contains(public, `Where("slug = ?`) {
		t.Error("the public variants endpoint still looks up a column called slug")
	}

	seeder := scaffold.APIVariantSeederGo(module, "Product", "product", "products", "Handle")
	// productSkuPrefix, not skuPrefix: perResource renames the helpers so two
	// resources with variants do not declare the same function twice.
	if !strings.Contains(seeder, "SkuPrefix(row.Handle, row.ID)") {
		t.Error("the seeder does not read row.Handle, so it names a field the model does not have")
	}
	if strings.Contains(seeder, "row.Slug") {
		t.Error("the seeder still reads row.Slug")
	}

	// And a resource with no slug keeps the id lookup rather than inventing a
	// column, which would be a SQL error instead of the 404 the caller deserves.
	none := scaffold.APIVariantPublicGo(module, "Gadget", "gadget", "gadgets", "", false, false)
	if !strings.Contains(none, `Where("id = ?", c.Param("key"))`) {
		t.Error("a resource with no slug should be looked up by id alone")
	}
	if strings.Contains(none, "OR id = ?") {
		t.Error("a resource with no slug got a two-column lookup")
	}
}
