package scaffold

import (
	"strings"
	"testing"
)

func TestProjectInitial(t *testing.T) {
	for name, want := range map[string]string{
		"contacts":   "C",
		"ui":         "U",
		"2fa-portal": "F",
		"":           "G",
		"123":        "G",
	} {
		if got := projectInitial(Options{ProjectName: name}); got != want {
			t.Errorf("%q: got %q, want %q", name, got, want)
		}
	}
}

// The logo mark beside the project's name is the project's letter.
func TestLogoMarkIsNotTheFrameworks(t *testing.T) {
	opts := Options{ProjectName: "contacts", Architecture: ArchTriple, Frontend: FrontendNext}
	for name, src := range map[string]string{
		"the web navbar":        webNavbar(opts),
		"the auth navbar":       webNavbarWithAuth(opts),
		"the single's navbar":   singleViteNavbar(Options{ProjectName: "contacts", Architecture: ArchSingle}),
		"the desktop title bar": desktopClientTitleBar(opts),
	} {
		if strings.Contains(src, `bold text-sm">G</span>`) || strings.Contains(src, `bold text-white">G</span>`) {
			t.Errorf("%s still shows the framework's initial", name)
		}
		if !strings.Contains(src, ">C</span>") {
			t.Errorf("%s does not show the project's initial", name)
		}
	}
}
