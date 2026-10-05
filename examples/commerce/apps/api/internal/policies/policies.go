package policies

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

import "commerce/apps/api/internal/authz"

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

	registerCollections()

	registerProducts()

	registerPages()

	registerCarts()

	registerCartItems()

	// grit:policies
}

// Referenced so the import compiles while Register is empty.
var _ = authz.Allow
