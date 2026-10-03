package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// RepairGeneratedServicePolicies adds the policy check to a service generated
// before policies existed.
//
// # Why a repair and not just new code
//
// A service is the developer's file. The upgrade will not rewrite it, which is
// why a resource generated two years ago still has the handler its author
// edited. But the policy check has to be in it, or policies are a feature only
// new resources get, and a project with forty resources would have to regenerate
// all forty to use one rule.
//
// So this does what every repair in this project does: finds the exact text the
// generator wrote, replaces it with the exact text the generator writes today,
// and when the text is not what the generator wrote, changes nothing and says
// so. A service somebody has restructured is left alone; a rule added to it by
// hand is one line in the same place.
//
// # What it changes
//
//   - load gains the ability it is loading for, so a rule can refuse a delete
//     without refusing an update.
//   - Update and Patch pass "update", Delete passes "delete".
//   - GetByID checks <plural>.read, load checks <plural>.<ability>.
//   - authz is imported if it was not already.
//
// It does not touch the handler. A service that returns a DeniedError to a
// handler with no case for it falls through to WriteError, which answers 500
// instead of 403: wrong, but visible and safe. RepairGeneratedHandlerDenials
// fixes that half.
func RepairGeneratedServicePolicies(src, module string) (string, bool) {
	m := policyLoadRe.FindStringSubmatch(src)
	if m == nil {
		return src, false
	}
	pascal := m[1]
	if m[2] != pascal {
		return src, false
	}

	// Already repaired, or generated after policies shipped.
	if strings.Contains(src, "authz.Inspect(ctx,") {
		return src, false
	}

	plural := serviceOwnsPlural(src, pascal)
	if plural == "" {
		return src, false
	}

	out := src

	// 1. load's signature, and its check.
	oldDoc := "// load reads the row a write is about to change, without its relations.\n"
	newDoc := oldDoc +
		"//\n" +
		"// ability is what the caller is about to do, \"update\" or \"delete\", so a\n" +
		"// policy can refuse one without refusing the other.\n"
	if strings.Count(out, oldDoc) != 1 {
		return src, false
	}
	out = strings.Replace(out, oldDoc, newDoc, 1)

	oldSig := "func (s *" + pascal + "Service) load(ctx context.Context, id string) (*models." + pascal + ", error) {"
	newSig := "func (s *" + pascal + "Service) load(ctx context.Context, id, ability string) (*models." + pascal + ", error) {"
	if strings.Count(out, oldSig) != 1 {
		return src, false
	}
	out = strings.Replace(out, oldSig, newSig, 1)

	loadBody, ok := methodBody(out, newSig)
	if !ok {
		return src, false
	}
	loadCheck := "\t// A policy rule narrows this further than the route's permission did.\n" +
		"\t// None defined means no narrowing, which is how the project behaved\n" +
		"\t// before it had policies. See internal/policies.\n" +
		"\tif d := authz.Inspect(ctx, \"" + plural + ".\"+ability, &item); !d.Allowed() {\n" +
		"\t\treturn nil, authz.Denied(\"" + plural + ".\"+ability, d)\n\t}\n"
	patchedLoad, ok := insertBeforeFinalReturn(loadBody, loadCheck)
	if !ok {
		return src, false
	}
	out = strings.Replace(out, loadBody, patchedLoad, 1)

	// 2. GetByID's read check.
	getSig := "func (s *" + pascal + "Service) GetByID(ctx context.Context, id string) (*models." + pascal + ", error) {"
	if strings.Count(out, getSig) != 1 {
		return src, false
	}
	getBody, ok := methodBody(out, getSig)
	if !ok {
		return src, false
	}
	readCheck := "\tif d := authz.Inspect(ctx, \"" + plural + ".read\", &item); !d.Allowed() {\n" +
		"\t\treturn nil, authz.Denied(\"" + plural + ".read\", d)\n\t}\n"
	patchedGet, ok := insertBeforeFinalReturn(getBody, readCheck)
	if !ok {
		return src, false
	}
	out = strings.Replace(out, getBody, patchedGet, 1)

	// 3. The three call sites, each inside the method that makes it. Replacing
	//    them positionally rather than by count, because all three read the
	//    same and a blind replace would give Delete the update ability.
	for _, call := range []struct{ method, ability string }{
		{"Update", "update"}, {"Patch", "update"}, {"Delete", "delete"},
	} {
		sig := "func (s *" + pascal + "Service) " + call.method + "(ctx context.Context,"
		body, found := methodBodyByPrefix(out, sig)
		if !found || strings.Count(body, "s.load(ctx, id)") != 1 {
			return src, false
		}
		out = strings.Replace(out, body,
			strings.Replace(body, "s.load(ctx, id)", "s.load(ctx, id, \""+call.ability+"\")", 1), 1)
	}

	// 4. The import, if ownership scoping did not already bring it in.
	authzImport := "\t\"" + module + "/internal/authz\"\n"
	if !strings.Contains(out, authzImport) {
		anchor := "\t\"" + module + "/internal/concurrency\"\n"
		if strings.Count(out, anchor) != 1 {
			return src, false
		}
		out = strings.Replace(out, anchor, authzImport+anchor, 1)
	}

	return out, true
}

var (
	// policyLoadRe matches a generated service's load, which every one has and
	// nothing else in the file looks like.
	policyLoadRe = regexp.MustCompile(`func \(s \*(\w+)Service\) load\(ctx context\.Context, id string\) \(\*models\.(\w+), error\) \{`)
	// servicePluralRe reads the resource's plural out of the doc comment the
	// generator writes, which is the only place in the file that has it.
	servicePluralRe = regexp.MustCompile(`// (\w+)Service owns every database read and write for ([\w-]+)\.`)
)

// serviceOwnsPlural is the plural the permission keys use, read from the
// service's own doc comment so the ability matches the permission beside it.
func serviceOwnsPlural(src, pascal string) string {
	for _, m := range servicePluralRe.FindAllStringSubmatch(src, -1) {
		if m[1] == pascal {
			return m[2]
		}
	}
	return ""
}

// methodBody returns the text of the method starting at sig, through its
// closing brace at column zero.
func methodBody(src, sig string) (string, bool) {
	start := strings.Index(src, sig)
	if start < 0 {
		return "", false
	}
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return "", false
	}
	return src[start : start+end+3], true
}

// methodBodyByPrefix is methodBody for a signature whose tail varies, which the
// write methods' do: Update takes link and item parameters only for some
// resources.
func methodBodyByPrefix(src, prefix string) (string, bool) {
	start := strings.Index(src, prefix)
	if start < 0 {
		return "", false
	}
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return "", false
	}
	return src[start : start+end+3], true
}

// insertBeforeFinalReturn puts block ahead of a method's `return &item, nil`.
//
// The last occurrence, because an owned resource's load returns early when the
// caller does not own the row, and the policy check belongs after that: there is
// no point explaining a refusal for a row the caller is not allowed to know
// exists.
func insertBeforeFinalReturn(body, block string) (string, bool) {
	const ret = "\treturn &item, nil\n"
	idx := strings.LastIndex(body, ret)
	if idx < 0 {
		return "", false
	}
	return body[:idx] + block + body[idx:], true
}

// RepairGeneratedHandlerDenials teaches a generated handler's fail helper to
// answer a policy denial with 403 and the rule's own reason.
//
// Without it the denial falls through to WriteError and comes back as a 500,
// which is wrong but visible: the request is still refused, and the server log
// carries the ability and the reason. Visible and wrong beats silent and right
// here, because a 500 gets investigated and a quietly-swallowed rule does not.
func RepairGeneratedHandlerDenials(src, module string) (string, bool) {
	m := handlerFailRe.FindStringSubmatch(src)
	if m == nil {
		return src, false
	}
	pascal := m[1]

	if strings.Contains(src, "var denied *authz.DeniedError") {
		return src, false
	}

	old := "\tvar conflict *concurrency.ErrConflict\n\tswitch {\n\tcase errors.As(err, &conflict):\n"
	if strings.Count(src, old) != 1 {
		return src, false
	}
	// 403 and not the 404 ownership uses: the caller is already allowed to see
	// this row, so hiding it says nothing and withholds the reason, which is the
	// thing the rule exists to give.
	out := strings.Replace(src, old,
		"\tvar conflict *concurrency.ErrConflict\n"+
			"\tvar denied *authz.DeniedError\n"+
			"\tswitch {\n"+
			"\tcase errors.As(err, &conflict):\n", 1)

	caseAnchor := "\t\tconcurrency.WriteConflict(c, conflict.Current)\n"
	if strings.Count(out, caseAnchor) != 1 {
		return src, false
	}
	out = strings.Replace(out, caseAnchor,
		caseAnchor+
			"\tcase errors.As(err, &denied):\n"+
			"\t\t// A policy refused, and said why. 403 rather than the 404 that\n"+
			"\t\t// ownership uses: the caller is already allowed to see this row,\n"+
			"\t\t// so hiding it says nothing and withholds the reason, which is the\n"+
			"\t\t// thing the rule exists to give.\n"+
			"\t\trespond.Fail(c, respond.CodeForbidden, denied.Reason)\n", 1)

	authzImport := "\t\"" + module + "/internal/authz\"\n"
	if !strings.Contains(out, authzImport) {
		anchor := "\t\"" + module + "/internal/concurrency\"\n"
		if strings.Count(out, anchor) != 1 {
			return src, false
		}
		out = strings.Replace(out, anchor, authzImport+anchor, 1)
	}

	_ = pascal
	return out, true
}

var handlerFailRe = regexp.MustCompile(`func \(h \*(\w+)Handler\) fail\(c \*gin\.Context, err error, fallback string\) \{`)

// repairGeneratedPolicies applies the policy check to every generated service
// and handler in an existing project.
//
// Both halves in one pass, because a service that returns a DeniedError to a
// handler with no case for it answers 500 instead of 403. Repairing the
// service without the handler would be a regression for one upgrade, and an
// upgrade that is briefly worse is one people learn not to run.
func repairGeneratedPolicies(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module, err := detectModule(apiRoot)
	if err != nil || module == "" {
		// No module path means the import cannot be written, and a repair that
		// adds a call without its import is a project that does not compile.
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	services, err := goSources(filepath.Join(apiRoot, "internal", "services"))
	if err != nil {
		return err
	}
	for _, path := range services {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out, changed := RepairGeneratedServicePolicies(src, module)
			if !changed {
				return src, nil, nil
			}
			return out, []string{"a rule in internal/policies can now refuse a read, an update or a delete of one row, and say why"}, nil
		}); err != nil {
			return err
		}
	}

	handlers, err := goSources(filepath.Join(apiRoot, "internal", "handlers"))
	if err != nil {
		return err
	}
	for _, path := range handlers {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out, changed := RepairGeneratedHandlerDenials(src, module)
			if !changed {
				return src, nil, nil
			}
			return out, []string{"a policy denial is answered with 403 and the rule's own reason, instead of a 500"}, nil
		}); err != nil {
			return err
		}
	}

	return nil
}
