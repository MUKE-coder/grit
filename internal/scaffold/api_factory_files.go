package scaffold

// Test factories for the framework's own models.
//
// `grit generate resource X --fields "...:belongs_to:User"` emits a fixture
// line for every parent:
//
//	author := factory.CreateUser(t, db)
//
// and factories are generated per resource, so the parent needs one. User and
// Upload are the framework's, not a generated resource's, so nothing ever wrote
// theirs and the emitted line referred to a function that did not exist:
//
//	vet: internal/handlers/note_test.go:62:18: undefined: factory.CreateUser
//
// `go build` passes, because the file is a _test.go. `go test` runs a subset of
// vet as an error before it compiles, so the package's tests refuse to build
// instead. That is why this survived: every local `go build` was clean, and the
// only thing that sees it is a job that runs the generated project's tests.
//
// Framework-owned, and in writeFrameworkOwnedFiles rather than written once at
// scaffold time, so a project generated before this gets them from
// `grit upgrade` too.

// apiFactoryCounterGo is the shared counter the factories interpolate.
//
// internal/generate writes this same file when it generates the first
// resource's factory, and reads this constant to do it, so the two cannot
// differ. The generator's write is conditional on the file's absence, so the
// scaffold shipping it first is not a conflict.
func apiFactoryCounterGo() string {
	return `package factory

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
`
}

// apiUserFactoryGo is the User fixture every `belongs_to:User` needs.
func apiUserFactoryGo() string {
	return `package factory

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// NewUser builds a valid User without saving it.
//
// Email is unique in the schema, so every call gets its own.
//
// Password is left empty, which means User.BeforeCreate does not hash it. That
// is deliberate: almost every fixture exists to satisfy a foreign key, and
// bcrypt at the default cost is tens of milliseconds a row, which a suite that
// creates users in a loop pays for nothing. A test that signs in asks for one:
//
//	user := factory.CreateUser(t, db, func(u *models.User) {
//		u.Password = "correct-horse-battery-staple"
//	})
//
// EmailVerifiedAt is set, because a fixture cannot receive a link and anything
// gated on a confirmed address would otherwise refuse it.
func NewUser(overrides ...func(*models.User)) *models.User {
	n := next()
	verifiedAt := time.Now()

	record := &models.User{
		FirstName:       fmt.Sprintf("Test%d", n),
		LastName:        "User",
		Email:           fmt.Sprintf("user%d@example.test", n),
		Role:            models.RoleUser,
		Active:          true,
		EmailVerifiedAt: &verifiedAt,
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateUser builds one and writes it, failing the test if it will not save.
//
// Failing here rather than returning an error is deliberate: a fixture that
// cannot be created is not a test failure to be asserted on, it is a broken
// test, and the stack should point at the line that asked for it.
func CreateUser(tb testing.TB, db *gorm.DB, overrides ...func(*models.User)) *models.User {
	tb.Helper()
	record := NewUser(overrides...)
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating a user fixture: %v", err)
	}
	return record
}

// CreateUsers writes n of them, for a test about paging or sorting.
func CreateUsers(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.User)) []*models.User {
	tb.Helper()
	out := make([]*models.User, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, CreateUser(tb, db, overrides...))
	}
	return out
}

// CreateAdmin is the one override asked for often enough to name.
func CreateAdmin(tb testing.TB, db *gorm.DB, overrides ...func(*models.User)) *models.User {
	tb.Helper()
	all := append([]func(*models.User){func(u *models.User) {
		u.Role = models.RoleAdmin
	}}, overrides...)
	return CreateUser(tb, db, all...)
}
`
}

// apiUploadFactoryGo is the Upload fixture a `belongs_to:Upload` needs.
func apiUploadFactoryGo() string {
	return `package factory

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// NewUpload builds a valid Upload without saving it.
//
// UserID is left for CreateUpload to fill, because it cannot be filled without
// writing a row: see the note there.
func NewUpload(overrides ...func(*models.Upload)) *models.Upload {
	n := next()

	record := &models.Upload{
		Filename:     fmt.Sprintf("fixture-%d.png", n),
		OriginalName: fmt.Sprintf("Fixture %d.png", n),
		MimeType:     "image/png",
		Size:         int64(1024 * n),
		Path:         fmt.Sprintf("uploads/fixture-%d.png", n),
		URL:          fmt.Sprintf("http://localhost:8080/files/uploads/fixture-%d.png", n),
	}

	for _, apply := range overrides {
		apply(record)
	}
	return record
}

// CreateUpload builds one and writes it, failing the test if it will not save.
//
// Unlike a generated resource's factory, this one creates its parent. Upload's
// user_id is NOT NULL, so a fixture with an empty one is rejected by Postgres
// and the test fails on the insert rather than on anything it meant to assert.
// Pass a user to attach it to an existing one:
//
//	upload := factory.CreateUpload(t, db, func(u *models.Upload) {
//		u.UserID = owner.ID
//	})
func CreateUpload(tb testing.TB, db *gorm.DB, overrides ...func(*models.Upload)) *models.Upload {
	tb.Helper()
	record := NewUpload(overrides...)
	if record.UserID == "" {
		record.UserID = CreateUser(tb, db).ID
	}
	if err := db.Create(record).Error; err != nil {
		tb.Fatalf("creating an upload fixture: %v", err)
	}
	return record
}

// CreateUploads writes n of them, for a test about paging or sorting.
//
// They share one owner: n users for n uploads is rows nothing asked for, and a
// test that wants them apart passes a UserID.
func CreateUploads(tb testing.TB, db *gorm.DB, n int, overrides ...func(*models.Upload)) []*models.Upload {
	tb.Helper()
	owner := CreateUser(tb, db)
	out := make([]*models.Upload, 0, n)
	for i := 0; i < n; i++ {
		all := append([]func(*models.Upload){func(u *models.Upload) {
			u.UserID = owner.ID
		}}, overrides...)
		out = append(out, CreateUpload(tb, db, all...))
	}
	return out
}
`
}

// FrameworkFactoryModels is the set of models the scaffold ships a factory for.
//
// internal/generate reads this to decide whether a `belongs_to` parent already
// has a fixture builder, and warns when it does not. Two things that have to
// agree: the files written above, and this list. TestFrameworkFactoriesMatchTheList
// reads both.
var FrameworkFactoryModels = []string{"User", "Upload"}

// FactoryCounterGo is the counter file, for internal/generate.
//
// The generator writes it when it generates the first resource's factory in a
// project old enough not to have it already. Exported so there is one copy of
// the text rather than two that can drift.
func FactoryCounterGo() string {
	return apiFactoryCounterGo()
}
