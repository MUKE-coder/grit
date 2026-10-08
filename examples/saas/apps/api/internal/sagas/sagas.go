// Package sagas holds this application's multi-step processes.
//
// A saga is for work that crosses systems and has to either finish or be
// undone: charge a card, reserve stock, book a courier. Each step has a Do and
// an Undo, where the run got to is a row rather than a stack frame, and a
// failure past its retries undoes the completed steps newest first. See
// internal/saga for the engine and the three rules a step has to follow.
//
// `grit generate workflow Checkout --steps "charge,reserve,ship"` writes a file
// here with the steps stubbed out, and adds it to Register below.
package sagas

// Register wires every saga in this package.
//
// Called once from routes.go before the runner starts. Registration is explicit
// rather than an init() per file, because a saga that silently disappears when
// somebody reorders an import is not a thing to debug at four in the morning.
func Register() {
	// grit:sagas
}
