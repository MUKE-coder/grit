package services

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
)

// UserService owns every database read and write behind the user endpoints.
//
// The handler reads the request, calls one of these and writes the answer, the
// same way every generated resource works. That matters more here than
// anywhere: a job that deactivates accounts, a console command that fixes a
// role, and PUT /api/users/:id have to apply the same rules, and a rule written
// in a handler applies to the HTTP route alone.
//
// Each method takes the context it runs in, which carries the organization a
// multitenant project resolved and the cancellation that fires when a client
// goes away.
type UserService struct {
	DB *gorm.DB
}

// db binds the database to ctx, so whatever a middleware put there reaches
// GORM's callbacks.
func (s *UserService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// userListConfig is the one description of what the users list may be
// searched, sorted and filtered by.
//
// The filters are the point. The admin's users page ships a Role dropdown, a
// Status toggle and a Provider dropdown, and its "Active Users" stat card asks
// for /api/users?active=true&page_size=1. The hand-rolled list read none of
// them: every filter returned the whole table and the Active Users card
// reported the total user count.
var userListConfig = paginate.Config{
	Searchable: []string{"first_name", "last_name", "email"},
	Sortable: map[string]bool{
		"id": true, "first_name": true, "last_name": true,
		"email": true, "role": true, "created_at": true,
	},
	Filterable:   map[string]bool{"role": true, "active": true, "provider": true},
	DefaultSort:  "created_at",
	DefaultOrder: "desc",
}

// List returns one page of users.
func (s *UserService) List(ctx context.Context, p paginate.Params) (paginate.Result[models.User], error) {
	return paginate.List[models.User](s.db(ctx).Model(&models.User{}), p, userListConfig)
}

// GetByID reads one account.
//
// The error comes back as GORM wrote it, so the caller can tell a missing row
// from a database that is down: the handler answers the first with a 404 and
// the second with a 500, and for a long time answered both with a 404.
func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := s.db(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail is the lookup by address, for a caller holding one: an SSO assertion,
// an invitation, a support tool. It returns the error GORM gave it, so a missing
// row and a database that is down stay distinguishable.
func (s *UserService) ByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// EmailTakenByAnother reports whether some other account already uses this
// address, case-insensitively.
//
// This is a check before a write rather than the write's own constraint,
// because the answer is shown to the person typing. A create relies on the
// unique index instead: see CreateUser.
func (s *UserService) EmailTakenByAnother(ctx context.Context, email, exceptID string) (bool, error) {
	var taken int64
	if err := s.db(ctx).Model(&models.User{}).
		Where("LOWER(email) = LOWER(?) AND id <> ?", email, exceptID).Count(&taken).Error; err != nil {
		return false, err
	}
	return taken > 0, nil
}

// Update writes updates to one account, and, when roleName is not empty, makes
// the user_roles assignment match it.
//
// Both writes, or neither.
//
// Grant resolution prefers the user_roles table and only falls back to
// users.role, so changing the Role dropdown has to change the assignment as
// well or it is a silent no-op. They used to be two separate writes, the
// assignment first: when the Updates that followed failed, the user kept the
// permissions of the new role while users.role still read as the old one, and
// nothing anywhere said so. One transaction now, so a failure leaves the
// account exactly as it was.
func (s *UserService) Update(ctx context.Context, user *models.User, updates map[string]interface{}, roleName string) error {
	if err := s.db(ctx).Transaction(func(tx *gorm.DB) error {
		if roleName != "" {
			if err := writeUserRoleAssignment(tx, user.ID, roleName); err != nil {
				return err
			}
		}
		return tx.Model(user).Updates(updates).Error
	}); err != nil {
		return err
	}
	// Permissions may just have changed for this user; drop the cached grants.
	// After the commit, never inside it: a request that resolved grants while
	// the transaction was still open would cache rows that had yet to exist.
	if roleName != "" {
		authz.Invalidate()
	}
	return nil
}

// SetPassword stores an already-hashed password.
//
// Hashing is the caller's: a service that took a plaintext password would be a
// second place where the cost parameter is chosen, and the two would drift.
func (s *UserService) SetPassword(ctx context.Context, userID, hashed string) error {
	return s.db(ctx).Model(&models.User{}).Where("id = ?", userID).
		Update("password", hashed).Error
}

// ClearLockout unlocks an account early and reports whether there was an account
// to unlock.
//
// Waiting out the window is the normal path. This is the support call that
// follows somebody locking themselves out five minutes before a demo, which is
// why it is an administrative action on a user rather than part of signing in.
func (s *UserService) ClearLockout(ctx context.Context, userID string) (bool, error) {
	res := s.db(ctx).Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"locked_until": nil, "failed_login_count": 0})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// Delete soft-deletes an account.
func (s *UserService) Delete(ctx context.Context, user *models.User) error {
	return s.db(ctx).Delete(user).Error
}

// SyncRoleAssignment makes the user_roles table reflect a single role name.
//
// The admin UI edits a user's role as one string, while authorization resolves
// through the many-to-many user_roles table. This bridges the two so the simple
// dropdown keeps working and actually takes effect.
//
// Assign several roles to one user with PUT /api/users/:id/roles instead: that
// endpoint is the full many-to-many path and this is its one-role case.
//
// A role name with no matching row (a custom legacy string) clears the
// assignments and lets grant resolution fall back to users.role, rather than
// failing the update.
func (s *UserService) SyncRoleAssignment(ctx context.Context, userID, roleName string) error {
	if err := s.db(ctx).Transaction(func(tx *gorm.DB) error {
		return writeUserRoleAssignment(tx, userID, roleName)
	}); err != nil {
		return err
	}
	authz.Invalidate()
	return nil
}

// writeUserRoleAssignment is the same work with no transaction of its own, for
// a caller that already has one. Update needs it: the assignment and the row
// are one change and have to succeed or fail together.
//
// It does not invalidate the grant cache. Only the caller knows when its
// transaction commits, and that is the moment the cache is stale.
func writeUserRoleAssignment(tx *gorm.DB, userID, roleName string) error {
	var role models.Role
	err := tx.Where("name = ?", roleName).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}
	if role.ID == "" {
		return nil // unknown name: fall back to the legacy string
	}
	return tx.Create(&models.UserRole{UserID: userID, RoleID: role.ID}).Error
}
