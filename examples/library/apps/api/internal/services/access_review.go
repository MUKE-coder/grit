package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
)

var (
	// ErrReviewClosed is returned when a decision is attempted on a review that
	// has already been signed off. A completed review is immutable evidence.
	ErrReviewClosed = errors.New("access review is already completed")
	// ErrReviewIncomplete is returned when completion is attempted while items
	// are still pending. You cannot certify a review you have not finished.
	ErrReviewIncomplete = errors.New("access review still has undecided items")
	// ErrItemDecided guards the one-way door: a revoked grant is already gone, so
	// re-approving it here would record a certification the system can't honour.
	ErrItemDecided = errors.New("this item has already been revoked and cannot be changed")
)

// OpenAccessReview snapshots every current role assignment into a new campaign
// of pending items. The snapshot copies user email and role name so the record
// stays legible after later deletions.
func OpenAccessReview(db *gorm.DB, name, note, createdBy, createdByEmail string) (*models.AccessReview, error) {
	review := &models.AccessReview{
		Name:           name,
		Note:           note,
		Status:         "open",
		CreatedBy:      createdBy,
		CreatedByEmail: createdByEmail,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return err
		}

		// Join UserRole to users and roles so the snapshot carries human-readable
		// identifiers, not just opaque ids.
		type grant struct {
			UserID    string
			UserEmail string
			RoleID    string
			RoleName  string
		}
		var grants []grant
		if err := tx.Table("user_roles").
			Select("user_roles.user_id, users.email AS user_email, user_roles.role_id, roles.name AS role_name").
			Joins("LEFT JOIN users ON users.id = user_roles.user_id").
			Joins("LEFT JOIN roles ON roles.id = user_roles.role_id").
			Scan(&grants).Error; err != nil {
			return err
		}

		items := make([]models.AccessReviewItem, 0, len(grants))
		for _, g := range grants {
			items = append(items, models.AccessReviewItem{
				ReviewID:  review.ID,
				UserID:    g.UserID,
				UserEmail: g.UserEmail,
				RoleID:    g.RoleID,
				RoleName:  g.RoleName,
				Decision:  "pending",
			})
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		review.Items = items
		return nil
	})
	if err != nil {
		return nil, err
	}
	return review, nil
}

// DecideAccessReviewItem records an approve/revoke decision. A revoke removes
// the underlying UserRole in the same transaction and logs it to the activity
// trail, so the decision and its effect can never drift apart.
func DecideAccessReviewItem(db *gorm.DB, c *gin.Context, reviewID, itemID, decision, note, reviewerID, reviewerEmail string) (*models.AccessReviewItem, error) {
	if decision != "approved" && decision != "revoked" {
		return nil, fmt.Errorf("decision must be \"approved\" or \"revoked\", got %q", decision)
	}

	var item models.AccessReviewItem
	err := db.Transaction(func(tx *gorm.DB) error {
		var review models.AccessReview
		if err := tx.First(&review, "id = ?", reviewID).Error; err != nil {
			return err
		}
		if review.Status == "completed" {
			return ErrReviewClosed
		}
		if err := tx.First(&item, "id = ? AND review_id = ?", itemID, reviewID).Error; err != nil {
			return err
		}
		// Revoke is terminal: the grant is already gone, so the record must not be
		// walked back to a certification the system can't stand behind.
		if item.Decision == "revoked" {
			return ErrItemDecided
		}

		if decision == "revoked" {
			// Remove the actual grant. Scoped to this user+role so nothing else is
			// touched; a no-op if it was already removed out of band.
			if err := tx.Where("user_id = ? AND role_id = ?", item.UserID, item.RoleID).
				Delete(&models.UserRole{}).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		item.Decision = decision
		item.Note = note
		item.DecidedBy = reviewerID
		item.DecidedByEmail = reviewerEmail
		item.DecidedAt = &now
		return tx.Save(&item).Error
	})
	if err != nil {
		return nil, err
	}

	// Log the revocation to the semantic activity trail (best-effort, outside the
	// transaction — losing the log line must not roll back a real revocation).
	if decision == "revoked" && c != nil {
		LogActivity(db, c, ActivityArgs{
			UserID:       reviewerID,
			Action:       "access_review.revoke",
			Severity:     "warn",
			Summary:      fmt.Sprintf("Revoked role %q from %s during access review", item.RoleName, item.UserEmail),
			ResourceType: "user",
			ResourceID:   item.UserID,
			Metadata: map[string]interface{}{
				"review_id": reviewID,
				"role_id":   item.RoleID,
			},
		})
	}
	return &item, nil
}

// CompleteAccessReview signs off a campaign. It refuses while any item is still
// pending: a review with undecided grants is not a review.
func CompleteAccessReview(db *gorm.DB, reviewID, reviewerID, reviewerEmail string) (*models.AccessReview, error) {
	var review models.AccessReview
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&review, "id = ?", reviewID).Error; err != nil {
			return err
		}
		if review.Status == "completed" {
			return ErrReviewClosed
		}
		var pending int64
		if err := tx.Model(&models.AccessReviewItem{}).
			Where("review_id = ? AND decision = ?", reviewID, "pending").
			Count(&pending).Error; err != nil {
			return err
		}
		if pending > 0 {
			return ErrReviewIncomplete
		}
		now := time.Now()
		review.Status = "completed"
		review.CompletedBy = reviewerID
		review.CompletedByEmail = reviewerEmail
		review.CompletedAt = &now
		return tx.Save(&review).Error
	})
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// AccessReviewSummary is the per-campaign counts the list view shows.
type AccessReviewSummary struct {
	models.AccessReview
	TotalItems    int64 `json:"total_items"`
	PendingItems  int64 `json:"pending_items"`
	ApprovedItems int64 `json:"approved_items"`
	RevokedItems  int64 `json:"revoked_items"`
}

// AccessReviewSummaries returns every campaign newest first, each with its
// decision counts.
//
// Two queries, whatever the number of campaigns. The handler ran four counts per
// campaign, so a year of monthly reviews was forty-nine queries to draw one
// page, and every one of them discarded its error: a count that failed drew a
// campaign as having no items at all, which in an access review reads as nothing
// left to certify.
func AccessReviewSummaries(ctx context.Context, db *gorm.DB) ([]AccessReviewSummary, error) {
	var reviews []models.AccessReview
	// Newest first, and by id within a second, because two campaigns opened in
	// the same second ordered arbitrarily: the ids Grit issues are time-ordered,
	// so they break the tie the way a reader expects. Found by a test that passed
	// alone and failed in a full run, which is the same ambiguity an operator
	// meets as a list that reorders itself between refreshes.
	if err := db.WithContext(ctx).Order("created_at desc, id desc").Find(&reviews).Error; err != nil {
		return nil, err
	}

	type tally struct {
		ReviewID string
		Decision string
		N        int64
	}
	var rows []tally
	if err := db.WithContext(ctx).Model(&models.AccessReviewItem{}).
		Select("review_id, decision, count(*) as n").
		Group("review_id, decision").Scan(&rows).Error; err != nil {
		return nil, err
	}
	counted := make(map[string]map[string]int64, len(reviews))
	for _, r := range rows {
		if counted[r.ReviewID] == nil {
			counted[r.ReviewID] = map[string]int64{}
		}
		counted[r.ReviewID][r.Decision] = r.N
	}

	out := make([]AccessReviewSummary, 0, len(reviews))
	for _, r := range reviews {
		byDecision := counted[r.ID]
		s := AccessReviewSummary{
			AccessReview:  r,
			PendingItems:  byDecision["pending"],
			ApprovedItems: byDecision["approved"],
			RevokedItems:  byDecision["revoked"],
		}
		// The total is every decision there is, including any a later release
		// adds, rather than the three named above plus a separate count.
		for _, n := range byDecision {
			s.TotalItems += n
		}
		out = append(out, s)
	}
	return out, nil
}

// AccessReviewWithItems reads one campaign and its items, in the order the
// reviewer works through them: by person, then by role.
func AccessReviewWithItems(ctx context.Context, db *gorm.DB, id string) (*models.AccessReview, error) {
	var review models.AccessReview
	if err := db.WithContext(ctx).Preload("Items", func(q *gorm.DB) *gorm.DB {
		return q.Order("user_email asc, role_name asc")
	}).First(&review, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}
