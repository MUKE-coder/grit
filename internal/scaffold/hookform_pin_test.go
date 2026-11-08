package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// Every app that ships FormBuilder has to install a resolvers version that has
// the subpath FormBuilder imports.
//
// v3.475.0 gave FormBuilder a Standard Schema resolver and raised the pin in
// two of the three package.json templates that carry that component. The
// TanStack admin kept ^3.9.1, which has no
// @hookform/resolvers/standard-schema, so from v3.475.0 to v3.484.0 every
// --vite project failed to build its admin:
//
//	[vite]: Rollup failed to resolve import "@hookform/resolvers/standard-schema"
//	from src/components/forms/form-builder.tsx
//
// Nothing caught it. CI type-checks the Next admin and the two web apps and not
// this one, the Go side compiled, and `grit doctor` reads the API. It was found
// by type-checking a Vite project while reviewing somebody else's pull request.
//
// So: one constant, and this test reads BOTH the component's import and every
// template that installs for it. A test on either side alone passes while the
// pair is broken, which is exactly what happened.
func TestEveryAppWithFormBuilderPinsAResolverThatHasTheSubpath(t *testing.T) {
	const subpath = "@hookform/resolvers/standard-schema"

	// The side that needs it.
	if !strings.Contains(adminFormBuilder(), subpath) {
		t.Fatalf("FormBuilder no longer imports %s, so this test checks nothing: "+
			"either the resolver went away, in which case drop the pin with it, or "+
			"the import moved and this test has to follow it", subpath)
	}

	// The major the subpath arrived in. Below this, the import does not resolve
	// and the app does not build.
	const firstMajorWithSubpath = 4
	major, err := majorOfPin(HookformResolversVersion)
	if err != nil {
		t.Fatalf("HookformResolversVersion = %q: %v", HookformResolversVersion, err)
	}
	if major < firstMajorWithSubpath {
		t.Errorf("HookformResolversVersion is %q, and %s arrived in %d.x. Every "+
			"app shipping FormBuilder would fail to build.",
			HookformResolversVersion, subpath, firstMajorWithSubpath)
	}

	// And the side that installs it: the three package.json templates that
	// carry FormBuilder. Each is checked against the constant rather than
	// against a literal, so raising the constant is the whole change.
	opts := Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext}
	for name, pkg := range map[string]string{
		"the Next admin":     adminPackageJSON(opts),
		"the TanStack admin": adminTanStackPackageJSON(opts),
	} {
		pin := resolversPin(pkg)
		if pin == "" {
			t.Errorf("%s does not install @hookform/resolvers at all, and it ships "+
				"FormBuilder", name)
			continue
		}
		if pin != HookformResolversVersion {
			t.Errorf("%s pins @hookform/resolvers at %q and the constant is %q. "+
				"FormBuilder is the same component in both; one of them cannot build it.",
				name, pin, HookformResolversVersion)
		}
	}
}

var resolversPinRe = regexp.MustCompile(`"@hookform/resolvers"\s*:\s*"([^"]+)"`)

// resolversPin reads the pin out of a package.json, or "" when it is absent.
func resolversPin(pkg string) string {
	m := resolversPinRe.FindStringSubmatch(pkg)
	if m == nil {
		return ""
	}
	return m[1]
}
