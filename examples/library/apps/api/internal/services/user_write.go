package services

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
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
