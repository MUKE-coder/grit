package generate

import (
	"strings"
	"testing"
)

// A unique column seeds from the row number. Four random letters and four
// random digits collided about a hundred times in a million rows, and with
// batched inserts each collision fails its whole batch.
func TestUniqueColumnSeedsFromTheRowNumber(t *testing.T) {
	def, err := ParseInlineFields("Product", "name:string,sku:string:unique")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := &Generator{Module: "acme/apps/api", Definition: def}
	out := g.seederContent(MakeNames("Product"), SeederOptions{Faker: true, Count: 50})

	if !strings.Contains(out, `Sku: fmt.Sprintf("SKU-%07d", i+1)`) {
		t.Errorf("the unique sku does not seed from the row number:\n%s", out)
	}
	if strings.Contains(out, "gofakeit.LetterN(4)") {
		t.Error("the unique sku still seeds from random characters")
	}
	if !strings.Contains(out, "\t\"fmt\"\n") || !strings.Contains(out, "i := int(n)") {
		t.Errorf("the seeder does not import fmt and declare i for the sku:\n%s", out)
	}
	if strings.Contains(out, "\t\"log\"\n") {
		t.Error("the faker seeder imports log, which only the static seeder uses")
	}
}

// A resource whose fields never read i must not declare it: an unused variable
// does not compile.
func TestSeederDeclaresTheRowIndexOnlyWhenUsed(t *testing.T) {
	def, err := ParseInlineFields("Note", "title:string,body:text,done:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := &Generator{Module: "acme/apps/api", Definition: def}
	out := g.seederContent(MakeNames("Note"), SeederOptions{Faker: true, Count: 10})
	if strings.Contains(out, "i := int(n)") {
		t.Errorf("a seeder with no row-dependent field declares i:\n%s", out)
	}
	for _, want := range []string{"return SeedNotesTo(db, 10)", `RegisterSeeder("Note", SeedNotesTo`, "SeedTopUp(db, \"notes\""} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}
