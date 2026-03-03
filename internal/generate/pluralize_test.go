package generate

import "testing"

func TestPluralize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Irregulars
		{"person", "people"},
		{"Person", "People"},
		{"child", "children"},
		{"mouse", "mice"},
		{"man", "men"},
		{"woman", "women"},
		{"tooth", "teeth"},
		{"foot", "feet"},
		{"datum", "data"},
		{"index", "indices"},
		{"matrix", "matrices"},

		// -y → -ies (consonant before y)
		{"Category", "Categories"},
		{"category", "categories"},
		{"story", "stories"},
		{"country", "countries"},
		{"city", "cities"},
		{"party", "parties"},

		// -y stays (vowel before y)
		{"day", "days"},
		{"key", "keys"},
		{"boy", "boys"},

		// -s, -sh, -ch, -x, -z → -es
		{"class", "classes"},
		{"bush", "bushes"},
		{"church", "churches"},
		{"box", "boxes"},
		{"buzz", "buzzes"},

		// -f / -fe → -ves
		{"leaf", "leaves"},
		{"knife", "knives"},
		{"life", "lives"},

		// Default: add -s
		{"Post", "Posts"},
		{"post", "posts"},
		{"Product", "Products"},
		{"User", "Users"},
		{"Comment", "Comments"},
		{"Tag", "Tags"},
		{"Blog", "Blogs"},

		// Edge cases
		{"", ""},
	}

	for _, tt := range tests {
		got := Pluralize(tt.input)
		if got != tt.want {
			t.Errorf("Pluralize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
