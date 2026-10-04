package policies

import (
	"library/apps/api/internal/authz"
	"library/apps/api/internal/models"
)

// Loans rules.
//
// Three abilities are checked around every loan the API loads:
//
//	loans.read     before returning one
//	loans.update   before an update or a patch
//	loans.delete   before a delete
//
// A rule you do not define is not a rule that denies. The route's permission
// has already said yes by the time one of these runs; a rule decides whether
// this particular record is an exception to that.
//
// This file is the worked example the generated stubs point at: it is what a
// real rule looks like once the comments come off.
func registerLoans() {
	// A returned loan is history.
	//
	// Nothing in the permission catalogue can express this, because it is not
	// about who is asking: a librarian with every permission still should not
	// be able to quietly change the dates on a loan that is already closed, and
	// --owned-by cannot say it either, because the borrower owns the row.
	//
	// Written as a rule, it is four lines and it is in one place. Written as an
	// `if` at the top of a handler, it is four lines in the update handler and
	// missing from the patch handler, which is the bug this exists to prevent.
	authz.Define("loans.update", func(actor authz.Actor, record any) authz.Decision {
		loan, ok := record.(*models.Loan)
		if !ok {
			// Another model reached this rule, which is a wiring mistake rather
			// than a refusal. Allow, and let the ability meant for it decide.
			return authz.Allow()
		}
		if loan.ReturnedAt == nil {
			return authz.Allow()
		}
		return authz.Deny("This loan was returned on " +
			loan.ReturnedAt.Format("2 January 2006") +
			", and a returned loan is a record of what happened rather than something to edit. " +
			"Create a new loan if the book has gone out again.")
	})

	// Deleting is separate from editing, and stricter.
	//
	// A loan is how the library knows where a book went. Deleting one does not
	// correct a mistake, it hides it, and the only honest way to undo a loan
	// that should never have existed is to say so in the record.
	authz.Define("loans.delete", func(actor authz.Actor, record any) authz.Decision {
		loan, ok := record.(*models.Loan)
		if !ok {
			return authz.Allow()
		}
		if loan.ReturnedAt != nil {
			return authz.Deny("A returned loan is part of the book's history and is not deleted. " +
				"If it was recorded in error, note that on the loan instead.")
		}
		return authz.Allow()
	})
}
