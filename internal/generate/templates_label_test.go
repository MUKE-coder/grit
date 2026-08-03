package generate

import (
	"strings"
	"testing"
)

// splitPascal drives every generated label, so an acronym it mangles shows up
// on screen in the admin ("Portfolio U R L") and in the API reference.
func TestSplitPascalKeepsAcronymsTogether(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"PortfolioURL", "Portfolio URL"},
		{"URL", "URL"},
		{"APIKey", "API Key"},
		{"PDFExport", "PDF Export"},
		{"UserID", "User ID"},
		{"ID", "ID"},
		// Ordinary PascalCase must be unaffected.
		{"FullName", "Full Name"},
		{"Name", "Name"},
		{"ExpectedSalary", "Expected Salary"},
		{"HTTPSProxyHost", "HTTPS Proxy Host"},
	}

	for _, c := range cases {
		got := strings.Join(splitPascal(c.in), " ")
		if got != c.want {
			t.Errorf("splitPascal(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
