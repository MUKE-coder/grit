package scaffold

// The one place a user row is inserted.
//
// Registration and the admin's Create User both used to do the same thing:
// SELECT the email, and INSERT when the SELECT found nothing. Two copies of a
// rule is one copy too many, and this particular rule does not work even once.
// Between the SELECT and the INSERT there is a window, and two signups for the
// same address that land inside it both pass the check. One of them then hits
// the unique index on users.email and the handler, which was only prepared for
// its own check to refuse, reports the constraint as a 500.
//
// So the index decides. It is the only thing in the system that can, because it
// is the only thing that sees both writes.

// apiUserWriteServiceGo emits internal/services/user_write.go.
func apiUserWriteServiceGo() string {
	return `package services

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// ErrEmailExists means the address is already registered.
//
// Return it as a 409 with code EMAIL_EXISTS. It is the only error CreateUser
// translates; everything else comes back as the driver wrote it.
var ErrEmailExists = errors.New("a user with this email already exists")

// CreateUser inserts a user and turns a duplicate email into ErrEmailExists.
//
// There is deliberately no SELECT first. The unique index on users.email is
// the check, it is atomic with the insert, and it is the same check whether
// the request arrived alone or alongside five others for the same address.
func CreateUser(ctx context.Context, db *gorm.DB, user *models.User) error {
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		if IsDuplicateKey(err) {
			return ErrEmailExists
		}
		return err
	}
	return nil
}

// IsDuplicateKey reports whether err is a unique-constraint violation.
//
// GORM has gorm.ErrDuplicatedKey, but it only translates the driver's error
// into it when the connection was opened with TranslateError, which a project
// scaffolded before that was set has not got. So the sentinel is checked first
// and the four drivers Grit runs on are matched by message as the fallback.
// Matching on a message is unpleasant; answering a duplicate email with a 500
// is worse.
func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, marker := range []string{
		"duplicate key value",      // Postgres
		"unique constraint failed", // SQLite
		"duplicate entry",          // MySQL
		"violation of unique key",  // SQL Server
	} {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}
`
}

// apiUserWriteServiceTestGo emits the test that holds the rule down: a second
// insert of the same address is a conflict, not a server error.
func apiUserWriteServiceTestGo() string {
	return `package services

import (
	"context"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
)

func userWriteDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}

func TestCreateUserRejectsADuplicateEmail(t *testing.T) {
	db := userWriteDB(t)
	first := models.User{FirstName: "A", LastName: "One", Email: "dup@example.com", Password: "secret123", Role: models.RoleUser, Active: true}
	require.NoError(t, CreateUser(context.Background(), db, &first))

	second := models.User{FirstName: "B", LastName: "Two", Email: "dup@example.com", Password: "secret123", Role: models.RoleUser, Active: true}
	err := CreateUser(context.Background(), db, &second)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmailExists), "a duplicate email must be ErrEmailExists, got %v", err)
}

func TestIsDuplicateKeyIgnoresEverythingElse(t *testing.T) {
	assert.False(t, IsDuplicateKey(nil))
	assert.False(t, IsDuplicateKey(errors.New("connection refused")))
	assert.True(t, IsDuplicateKey(gorm.ErrDuplicatedKey))
	assert.True(t, IsDuplicateKey(errors.New("UNIQUE constraint failed: users.email")))
}
`
}
