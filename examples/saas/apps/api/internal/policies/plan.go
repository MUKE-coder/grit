package policies

import (
	"saas/apps/api/internal/authz"
	"saas/apps/api/internal/models"
)

// Plans rules.
//
// Three abilities are checked around every plan the API loads:
//
//	plans.read     before returning one
//	plans.update   before an update or a patch
//	plans.delete   before a delete
//
// A rule you do not define is not a rule that denies. The route's permission
// has already said yes by the time one of these runs; a rule decides whether
// this particular record is an exception to that. So these are restrictions you
// add, never permissions you forget to grant, and an undefined ability behaves
// exactly as this project did before it had policies.
//
// Uncomment what you need. The actor is who is asking, and record is the row as
// the service loaded it, so a rule can read any column on it.
func registerPlans() {
	// Only the owner may change a plan once it is no longer a draft. The kind of
	// rule that lives at the top of two handlers and is missing from the third.
	//
	// authz.Define("plans.update", func(actor authz.Actor, record any) authz.Decision {
	// 	item, ok := record.(*models.Plan)
	// 	if !ok {
	// 		// Another model reached this rule, which is a wiring mistake
	// 		// rather than a refusal. Allow, and let the ability that is
	// 		// meant for it decide.
	// 		return authz.Allow()
	// 	}
	// 	if item.Status == "draft" {
	// 		return authz.Allow()
	// 	}
	// 	return authz.Deny("This plan has been published, so it is edited through a revision rather than in place")
	// })

	// A rule that refuses a delete without refusing an edit.
	//
	// authz.Define("plans.delete", func(actor authz.Actor, record any) authz.Decision {
	// 	item, ok := record.(*models.Plan)
	// 	if !ok {
	// 		return authz.Allow()
	// 	}
	// 	if item.Archived() {
	// 		return authz.Allow()
	// 	}
	// 	return authz.Deny("Archive this plan before deleting it, so anything referring to it keeps working")
	// })

	// Reads, for a resource where some rows are not everybody's business.
	//
	// authz.Define("plans.read", func(actor authz.Actor, record any) authz.Decision {
	// 	return authz.Allow()
	// })
}

// Referenced so the import compiles while every rule above is commented out.
// Delete this line once you have uncommented one.
var _ = authz.Allow

// Referenced for the same reason.
var _ = models.Plan{}
