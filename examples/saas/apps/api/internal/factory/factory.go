package factory

// Package factory builds valid records for tests.
//
// One builder per resource, generated beside the model, so the knowledge of
// what a valid row looks like lives in one place rather than in every test that
// needs one. Adding a required field is then one edit here instead of one per
// test file.

import "sync/atomic"

// counter makes each fixture distinct.
//
// Shared across every resource in the package on purpose: two resources with
// their own counters would both start at 1, and a test creating one of each
// would get two rows claiming to be the first, which is exactly the collision
// a unique column then rejects.
var counter atomic.Int64

// next is the number the builders interpolate into unique fields.
func next() int {
	return int(counter.Add(1))
}

// Reset puts the counter back, for a test that asserts on an exact generated
// value. Rarely what you want: asserting on "Widget 1" rather than on the value
// the factory returned is a test that breaks when somebody adds a fixture above
// it.
func Reset() {
	counter.Store(0)
}
