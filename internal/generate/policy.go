package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writePolicy emits internal/policies/<snake>.go: where this resource's rules go.
//
// Written as a stub with the rules commented out, not as an empty file and not
// as a working rule.
//
// An empty file says nothing about what belongs in it, and the developer who
// opens it looking for where to put "only the author may edit a draft" closes it
// again. A working rule would be a guess about the application's domain that is
// wrong more often than right, and a wrong rule that compiles is one somebody
// ships. A commented example is the only one of the three that teaches without
// deciding.
func (g *Generator) writePolicy(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "policies", names.Snake+".go")

	// Created once, never rewritten. The second `grit generate resource` for a
	// resource is how anybody adds a field, and taking somebody's rules away
	// because they re-ran the command is the trade nobody wants.
	if fileExists(path) {
		return nil
	}
	return writeFileWithDirs(path, g.policySource(names))
}

func (g *Generator) policySource(names Names) string {
	var b strings.Builder

	fmt.Fprintf(&b, `package policies

import (
	%q
	%q
)

// %s rules.
//
// Three abilities are checked around every %s the API loads:
//
//	%s.read     before returning one
//	%s.update   before an update or a patch
//	%s.delete   before a delete
//
// A rule you do not define is not a rule that denies. The route's permission
// has already said yes by the time one of these runs; a rule decides whether
// this particular record is an exception to that. So these are restrictions you
// add, never permissions you forget to grant, and an undefined ability behaves
// exactly as this project did before it had policies.
//
// Uncomment what you need. The actor is who is asking, and record is the row as
// the service loaded it, so a rule can read any column on it.
func register%s() {
	// Only the owner may change a %s once it is no longer a draft. The kind of
	// rule that lives at the top of two handlers and is missing from the third.
	//
	// authz.Define("%s.update", func(actor authz.Actor, record any) authz.Decision {
	// 	item, ok := record.(*models.%s)
	// 	if !ok {
	// 		// Another model reached this rule, which is a wiring mistake
	// 		// rather than a refusal. Allow, and let the ability that is
	// 		// meant for it decide.
	// 		return authz.Allow()
	// 	}
	// 	if item.Status == "draft" {
	// 		return authz.Allow()
	// 	}
	// 	return authz.Deny("This %s has been published, so it is edited through a revision rather than in place")
	// })

	// A rule that refuses a delete without refusing an edit.
	//
	// authz.Define("%s.delete", func(actor authz.Actor, record any) authz.Decision {
	// 	item, ok := record.(*models.%s)
	// 	if !ok {
	// 		return authz.Allow()
	// 	}
	// 	if item.Archived() {
	// 		return authz.Allow()
	// 	}
	// 	return authz.Deny("Archive this %s before deleting it, so anything referring to it keeps working")
	// })

	// Reads, for a resource where some rows are not everybody's business.
	//
	// authz.Define("%s.read", func(actor authz.Actor, record any) authz.Decision {
	// 	return authz.Allow()
	// })
}

// Referenced so the import compiles while every rule above is commented out.
// Delete this line once you have uncommented one.
var _ = authz.Allow

// Referenced for the same reason.
var _ = models.%s{}
`,
		g.Module+"/internal/authz",
		g.Module+"/internal/models",
		names.PluralPascal, names.Lower,
		names.Plural, names.Plural, names.Plural,
		names.PluralPascal, names.Lower,
		names.Plural, names.Pascal, names.Lower,
		names.Plural, names.Pascal, names.Lower,
		names.Plural,
		names.Pascal,
	)

	return b.String()
}

// policiesRegistrySource is internal/policies/policies.go: the one function
// main calls.
//
// A registry rather than init() functions, because the order rules are defined
// in decides which Before hook sees what, and a package full of inits has no
// order anybody can read.
func policiesRegistrySource(module string) string {
	return `package policies

// Package policies holds the rules about who may do what to which record.
//
// Three kinds of authorization live in this project, and they answer different
// questions:
//
//   - The permission catalogue, in internal/permissions, answers "may this user
//     touch this resource at all". It is checked by the route's middleware,
//     before a handler runs.
//   - --owned-by scoping, in internal/authz, answers "is this their row". It is
//     checked by the service, against the owner column.
//   - These rules answer everything conditional on the record itself: editable
//     while it is a draft, approvable by anyone but the raiser, readable by its
//     reviewers. They are checked by the service, around the row it loaded.
//
// A rule you do not write is not a rule that denies. See internal/authz/gate.go.

import "` + module + `/internal/authz"

// Register defines every rule. Called once from main, before the server starts.
func Register() {
	// An ADMIN is not exempt by default.
	//
	// It would be easy to add, and it is the wrong default: the rules here are
	// about the state of a record rather than the rank of a user, and "this
	// invoice is already paid" is as true for an administrator as for anybody
	// else. Add the exemption deliberately, for the abilities that want it:
	//
	// authz.Before(func(actor authz.Actor, ability string, record any) (authz.Decision, bool) {
	// 	if actor.Admin && strings.HasPrefix(ability, "drafts.") {
	// 		return authz.Allow(), true
	// 	}
	// 	return authz.Decision{}, false
	// })

	// grit:policies
}

// Referenced so the import compiles while Register is empty.
var _ = authz.Allow
`
}

// ensurePoliciesRegistry writes internal/policies/policies.go when it is not
// there yet, and registers the resource in it.
func (g *Generator) ensurePoliciesRegistry() error {
	path := filepath.Join(g.APIRoot(), "internal", "policies", "policies.go")
	if fileExists(path) {
		return nil
	}
	if err := writeFileWithDirs(path, policiesRegistrySource(g.Module)); err != nil {
		return fmt.Errorf("writing the policy registry: %w", err)
	}
	return nil
}

// injectPolicyRegistration adds this resource's register call to Register().
//
// injectBefore skips a call that is already there, so re-running the generator
// for a resource does not register it twice.
func (g *Generator) injectPolicyRegistration(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "policies", "policies.go")
	// The marker without its indentation: injectBefore matches it as a
	// standalone line, so leading whitespace in the needle never matches.
	return injectBefore(path, "// grit:policies", "\tregister"+names.PluralPascal+"()\n")
}
