package services_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
)

func newReviewDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(tb, err)
	require.NoError(tb, db.AutoMigrate(
		&models.User{}, &models.Role{}, &models.UserRole{},
		&models.AccessReview{}, &models.AccessReviewItem{}, &models.UserActivity{},
	))
	return db
}

// seedGrant creates a user, a role, and the assignment between them.
func seedGrant(tb testing.TB, db *gorm.DB, email, roleName string) (userID, roleID string) {
	tb.Helper()
	u := models.User{Email: email, FirstName: "T", LastName: "U", Password: "x"}
	require.NoError(tb, db.Create(&u).Error)
	r := models.Role{Name: roleName}
	require.NoError(tb, db.Create(&r).Error)
	require.NoError(tb, db.Create(&models.UserRole{UserID: u.ID, RoleID: r.ID}).Error)
	return u.ID, r.ID
}

func TestOpenAccessReviewSnapshotsGrants(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	seedGrant(t, db, "bob@x.dev", "VIEWER")

	review, err := services.OpenAccessReview(db, "Q3 2026", "quarterly", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	require.Len(t, review.Items, 2, "one item per current grant")

	for _, it := range review.Items {
		assert.Equal(t, "pending", it.Decision)
		assert.NotEmpty(t, it.UserEmail, "snapshot must copy the email")
		assert.NotEmpty(t, it.RoleName, "snapshot must copy the role name")
	}
}

// The record must survive later deletion of the user and role it reviewed —
// that is the entire reason for snapshotting.
func TestReviewSnapshotSurvivesDeletion(t *testing.T) {
	db := newReviewDB(t)
	userID, roleID := seedGrant(t, db, "gone@x.dev", "TEMP")

	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	require.NoError(t, db.Delete(&models.User{}, "id = ?", userID).Error)
	require.NoError(t, db.Delete(&models.Role{}, "id = ?", roleID).Error)

	var item models.AccessReviewItem
	require.NoError(t, db.First(&item, "review_id = ?", review.ID).Error)
	assert.Equal(t, "gone@x.dev", item.UserEmail)
	assert.Equal(t, "TEMP", item.RoleName)
}

func TestRevokeRemovesTheGrant(t *testing.T) {
	db := newReviewDB(t)
	userID, roleID := seedGrant(t, db, "ada@x.dev", "EDITOR")

	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	item := review.Items[0]

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, item.ID, "revoked", "no longer needed", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	var count int64
	db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	assert.Equal(t, int64(0), count, "revoke must delete the UserRole")
}

func TestApproveKeepsTheGrant(t *testing.T) {
	db := newReviewDB(t)
	userID, roleID := seedGrant(t, db, "ada@x.dev", "EDITOR")

	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "approved", "still valid", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	var count int64
	db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	assert.Equal(t, int64(1), count, "approve must keep the grant")
}

// A revoked grant is gone; the decision cannot be walked back to a certification
// the system can't honour.
func TestRevokeIsTerminal(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	item := review.Items[0]

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, item.ID, "revoked", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, item.ID, "approved", "oops", "admin-1", "admin@x.dev")
	assert.ErrorIs(t, err, services.ErrItemDecided)
}

func TestCannotCompleteWithPendingItems(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	seedGrant(t, db, "bob@x.dev", "VIEWER")
	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	// Decide only one of the two.
	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "approved", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.CompleteAccessReview(db, review.ID, "admin-1", "admin@x.dev")
	assert.ErrorIs(t, err, services.ErrReviewIncomplete)

	// Decide the second, then completion succeeds.
	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[1].ID, "revoked", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	done, err := services.CompleteAccessReview(db, review.ID, "admin-1", "admin@x.dev")
	require.NoError(t, err)
	assert.Equal(t, "completed", done.Status)
	require.NotNil(t, done.CompletedAt)
}

// A completed review is immutable evidence.
func TestCannotDecideOnCompletedReview(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "approved", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	_, err = services.CompleteAccessReview(db, review.ID, "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "revoked", "", "admin-1", "admin@x.dev")
	assert.ErrorIs(t, err, services.ErrReviewClosed)
}

// The list view's counts come from one grouped query rather than four per
// campaign, and they have to add up to the same answers.
//
// The handler ran 4n counts and discarded every error, so a count that failed
// drew a campaign as having nothing left to certify, which in an access review is
// the one wrong answer that matters.
func TestAccessReviewSummariesCountsEveryDecision(t *testing.T) {
	db := newReviewDB(t)
	ctx := context.Background()
	seedGrant(t, db, "a@example.com", "EDITOR")
	seedGrant(t, db, "b@example.com", "VIEWER")

	first, err := services.OpenAccessReview(db, "Q1", "", "admin", "admin@example.com")
	require.NoError(t, err)
	second, err := services.OpenAccessReview(db, "Q2", "", "admin", "admin@example.com")
	require.NoError(t, err)

	// Decide one item in the first campaign, leaving the other pending.
	var items []models.AccessReviewItem
	require.NoError(t, db.Where("review_id = ?", first.ID).Order("user_email asc").Find(&items).Error)
	require.Len(t, items, 2)
	require.NoError(t, db.Model(&items[0]).Update("decision", "approved").Error)

	out, err := services.AccessReviewSummaries(ctx, db)
	require.NoError(t, err)
	require.Len(t, out, 2)

	byID := map[string]services.AccessReviewSummary{}
	for _, s := range out {
		byID[s.ID] = s
	}
	one := byID[first.ID]
	assert.Equal(t, int64(2), one.TotalItems)
	assert.Equal(t, int64(1), one.ApprovedItems)
	assert.Equal(t, int64(1), one.PendingItems)
	assert.Equal(t, int64(0), one.RevokedItems)

	two := byID[second.ID]
	assert.Equal(t, int64(2), two.TotalItems, "a second campaign's items were counted against the first")
	assert.Equal(t, int64(2), two.PendingItems)

	// Newest first, which is the order the screen shows them in.
	assert.Equal(t, second.ID, out[0].ID)
}

// The detail view reads the items in the order a reviewer works through them: by
// person, then by role.
func TestAccessReviewWithItemsIsOrderedForTheReviewer(t *testing.T) {
	db := newReviewDB(t)
	ctx := context.Background()
	seedGrant(t, db, "zoe@example.com", "VIEWER")
	seedGrant(t, db, "alice@example.com", "EDITOR")
	review, err := services.OpenAccessReview(db, "Q3", "", "admin", "admin@example.com")
	require.NoError(t, err)

	loaded, err := services.AccessReviewWithItems(ctx, db, review.ID)
	require.NoError(t, err)
	require.Len(t, loaded.Items, 2)
	assert.Equal(t, "alice@example.com", loaded.Items[0].UserEmail)

	_, err = services.AccessReviewWithItems(ctx, db, "no-such-review")
	assert.Error(t, err)
}
