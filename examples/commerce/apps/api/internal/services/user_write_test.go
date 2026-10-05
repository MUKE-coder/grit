package services

import (
	"context"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
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
