package scaffold

import "testing"

// The progress report is only as good as the line it reads, and go's wording
// is not something this repo controls. So the parser is pinned to the shapes
// the toolchain actually prints, including the ones that must NOT be shown.
func TestDownloadingModuleReadsGoProgressLines(t *testing.T) {
	cases := []struct {
		line string
		want string
		ok   bool
	}{
		{"go: downloading github.com/gin-gonic/gin v1.10.0", "github.com/gin-gonic/gin", true},
		{"go: downloading gorm.io/gorm v1.25.12", "gorm.io/gorm", true},
		// No version, which happens for some module queries.
		{"go: downloading golang.org/x/text", "golang.org/x/text", true},

		// Everything else is not progress. A warning or an error rendered as
		// "downloading ..." on the running line would be worse than silence,
		// because it reads as the thing going fine.
		{"go: finding module for package github.com/x/y", "", false},
		{"go: added github.com/gin-gonic/gin v1.10.0", "", false},
		{"go: downloading", "", false},
		{"go: downloading ", "", false},
		{"", "", false},
		{"  go: downloading github.com/a/b v1", "", false}, // leading space: not go's format
		{"go: upgraded github.com/a/b v1.0.0 => v1.1.0", "", false},
	}

	for _, c := range cases {
		got, ok := downloadingModule(c.line)
		if ok != c.ok {
			t.Errorf("downloadingModule(%q) ok = %v, want %v", c.line, ok, c.ok)
			continue
		}
		if got != c.want {
			t.Errorf("downloadingModule(%q) = %q, want %q", c.line, got, c.want)
		}
	}
}
