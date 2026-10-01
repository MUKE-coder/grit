package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairAuthServiceQueries gives an older project the AuthService methods its
// repaired sign-in is about to call.
//
// The auth handler used to run its own queries, which is the thing the
// generator forbids for every resource it writes and which grit doctor checks
// for. Moving them onto the service fixed the handler and broke the upgrade
// path in the same change: repairLoginDisclosureSource rewrites an older
// project's Login to what the template says today, and what it says today is
// h.AuthService.UserByEmail. A project whose services/auth.go predates those
// methods would be handed calls to methods that are not there, and would stop
// compiling on an upgrade that reported success.
//
// So this runs first, and that is why the ordering in upgrade.go is not
// arbitrary.
//
// The methods come from the template rather than a copy kept here, for the
// reason the whole repair layer is moving that way: a copy drifts, and the
// drift is invisible because the two halves are joined by a Go "+" that no
// search for the text can follow.
func repairAuthServiceQueries(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "services", "auth.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairAuthServiceQueriesSource)
}

// serviceInternalImportRe matches an import of one of the project's own
// packages and captures the module path in front of it. A regexp over import
// lines rather than a search for "/internal/", so a path mentioned in a comment
// is not read as the module.
var serviceInternalImportRe = regexp.MustCompile(`(?m)^\t"([^"]+)/internal/[^"]+"$`)

func repairAuthServiceQueriesSource(src string) (string, []string, []string) {
	// UserByEmail is the one the repaired Login calls first, so its absence is
	// what makes the upgrade unsafe.
	if strings.Contains(src, "func (s *AuthService) UserByEmail(") {
		return src, nil, nil
	}
	if !strings.Contains(src, "type AuthService struct") {
		return src, nil, []string{"is not the AuthService Grit wrote, so the sign-in repair would call methods that are not there: add UserByEmail, UserByID, EnabledTwoFactor, StartPendingTOTP and ClearLoginFailures by hand, or take Grit's copy with grit upgrade --force"}
	}

	out := src
	// context and models, which the methods need and a service that only minted
	// tokens did not. Added to the existing group rather than one of their own,
	// so the file reads like gofmt left it.
	if !strings.Contains(out, "\n\t\"context\"\n") {
		const anchor = "import (\n"
		i := strings.Index(out, anchor)
		if i < 0 {
			return src, nil, []string{"has no import block, so the sign-in queries cannot be added to it: take Grit's copy of internal/services/auth.go"}
		}
		out = out[:i+len(anchor)] + "\t\"context\"\n" + out[i+len(anchor):]
	}
	if !strings.Contains(out, "/internal/models\"") {
		// Beside the project's other internal imports, whose prefix is whatever
		// this project's module is called. Read off an import rather than go.mod:
		// this function is given the file and nothing else, and every copy of
		// services/auth.go imports at least one of the project's own packages.
		m := serviceInternalImportRe.FindStringSubmatchIndex(out)
		if m == nil {
			return src, nil, []string{"imports none of the project's own packages, so its module path cannot be read to import models: take Grit's copy of internal/services/auth.go"}
		}
		module := out[m[2]:m[3]]
		out = out[:m[1]] + "\n\t\"" + module + "/internal/models\"" + out[m[1]:]
	}

	out = strings.TrimRight(out, "\n") + "\n\n" + repairBlock("api/services/auth.go", "auth-service-queries")
	return out, []string{"AuthService has the sign-in queries, so the handler can stop running them itself"}, nil
}
