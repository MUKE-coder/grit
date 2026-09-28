package generate

import (
	"strings"
	"testing"
)

// grit#91 item 4: every string field got gofakeit.Name(), so a catalogue seeded
// forty products called "Emily Gardner", and every whole number got 1 to 100,
// so they cost 37 shillings.
func TestFakerNamesFitTheResource(t *testing.T) {
	product := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name: "Product",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "int"},
				{Name: "stock", Type: "int"},
			},
		},
	}
	lines, _, _, _, _ := product.seederFieldLines("faker")
	for _, want := range []string{
		"gofakeit.ProductName()",
		"Price: gofakeit.Number(1000, 500000)",
		// A count is still a count.
		"Stock: gofakeit.Number(1, 100)",
	} {
		if !strings.Contains(lines, want) {
			t.Errorf("the product seeder is missing %q, got:\n%s", want, lines)
		}
	}

	// A resource whose rows are people keeps a person's name.
	customer := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Customer",
			Fields: []Field{{Name: "name", Type: "string"}},
		},
	}
	lines, _, _, _, _ = customer.seederFieldLines("faker")
	if !strings.Contains(lines, "gofakeit.Name()") {
		t.Errorf("a customer is a person and should get a person's name, got:\n%s", lines)
	}
}

// gofakeit.Word() put "moreover" in every title column.
func TestFakerTitlesReadAsTitles(t *testing.T) {
	post := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Post",
			Fields: []Field{{Name: "title", Type: "string"}},
		},
	}
	lines, _, _, _, _ := post.seederFieldLines("faker")
	if !strings.Contains(lines, "gofakeit.Sentence(6)") {
		t.Errorf("a title should read as one, got:\n%s", lines)
	}
}

func TestIsMoneyColumn(t *testing.T) {
	for _, yes := range []string{"price", "compare_at_price", "total_amount", "unit_cost", "delivery_fee"} {
		if !isMoneyColumn(yes) {
			t.Errorf("%s should be money", yes)
		}
	}
	for _, no := range []string{"stock", "quantity", "position", "depth", "rating"} {
		if isMoneyColumn(no) {
			t.Errorf("%s should not be money", no)
		}
	}
}
