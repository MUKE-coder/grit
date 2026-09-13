package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// goToolchainFloor is the Go release a generated API builds with. Go 1.26.4 and
// earlier carry eight standard library vulnerabilities the API reaches (H11 in
// the contact-app review), including the encoding/xml behind SAML sign-in.
const goToolchainFloor = "1.26.6"

// goModEdit runs go mod edit in dir. A variable, so a test can see the call
// without a Go toolchain.
var goModEdit = func(dir string, args ...string) error {
	cmd := exec.Command("go", append([]string{"mod", "edit"}, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// repairGoToolchain raises an existing project to goToolchainFloor: the go
// directive in the API's go.mod, which is what setup-go and the release workflow
// install, and the Go pinned in its Dockerfiles and CI workflows.
func repairGoToolchain(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if data, err := os.ReadFile(filepath.Join(apiRoot, "go.mod")); err == nil {
		if current := goDirective(string(data)); current != "" && goVersionLess(current, goToolchainFloor) {
			if err := goModEdit(apiRoot, "-go="+goToolchainFloor); err != nil {
				fmt.Printf("  ⚠ could not raise go.mod to go %s, so run `go mod edit -go=%s` in %s: %v\n",
					goToolchainFloor, goToolchainFloor, apiRoot, err)
			} else {
				fmt.Printf("  ✓ go.mod: go %s → %s, the release without the standard library vulnerabilities govulncheck reports\n",
					current, goToolchainFloor)
			}
		}
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	paths := []string{filepath.Join(apiRoot, "Dockerfile"), filepath.Join(root, "Dockerfile")}
	if matches, err := filepath.Glob(filepath.Join(root, ".github", "workflows", "*.yml")); err == nil {
		paths = append(paths, matches...)
	}
	for _, path := range paths {
		if !fileExists(path) {
			continue
		}
		if err := repairTextFile(root, m, path, repairGoVersionPinsSource); err != nil {
			return err
		}
	}
	return nil
}

var (
	goDirectiveRe = regexp.MustCompile(`(?m)^go (\d+\.\d+(?:\.\d+)?)\s*$`)
	goImageRe     = regexp.MustCompile(`golang:(\d+\.\d+(?:\.\d+)?)-alpine`)
	setupGoRe     = regexp.MustCompile(`go-version: '(\d+\.\d+(?:\.\d+)?)'`)
)

func goDirective(gomod string) string {
	if m := goDirectiveRe.FindStringSubmatch(gomod); m != nil {
		return m[1]
	}
	return ""
}

// goVersionLess compares Go release numbers: "1.25" is less than "1.25.1".
func goVersionLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		var x, y int
		if i < len(pa) {
			x, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			y, _ = strconv.Atoi(pb[i])
		}
		if x != y {
			return x < y
		}
	}
	return false
}

// repairGoVersionPinsSource raises golang:<v>-alpine images and setup-go
// go-version pins below the floor. A newer pin is left as it is.
func repairGoVersionPinsSource(src string) (string, []string, []string) {
	raised := 0
	raise := func(re *regexp.Regexp, format string) func(string) string {
		return func(match string) string {
			v := re.FindStringSubmatch(match)[1]
			if !goVersionLess(v, goToolchainFloor) {
				return match
			}
			raised++
			return fmt.Sprintf(format, goToolchainFloor)
		}
	}
	out := goImageRe.ReplaceAllStringFunc(src, raise(goImageRe, "golang:%s-alpine"))
	out = setupGoRe.ReplaceAllStringFunc(out, raise(setupGoRe, "go-version: '%s'"))
	if raised == 0 {
		return src, nil, nil
	}
	return out, []string{fmt.Sprintf("Go %s (%d %s)", goToolchainFloor, raised, plural(raised, "pin", "pins"))}, nil
}
