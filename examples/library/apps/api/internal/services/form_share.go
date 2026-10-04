package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
)

// FormShareService owns the shares and the submission rows behind them.
//
// The handler decides who may see a share and whether a password was right; this
// is the table. The two are worth separating here because a share is reachable
// without an account: the public form posts to it, and the one write that
// happens on every post used to lose submissions under any concurrency at all.
type FormShareService struct {
	DB *gorm.DB
}

func (s *FormShareService) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// All lists shares newest first, optionally for one resource.
func (s *FormShareService) All(ctx context.Context, resourceName string) ([]models.FormShare, error) {
	q := s.db(ctx).Order("created_at DESC, id DESC")
	if resourceName != "" {
		q = q.Where("resource_name = ?", resourceName)
	}
	var shares []models.FormShare
	if err := q.Find(&shares).Error; err != nil {
		return nil, err
	}
	return shares, nil
}

// ByID reads one share for the admin screen.
func (s *FormShareService) ByID(ctx context.Context, id string) (*models.FormShare, error) {
	var share models.FormShare
	if err := s.db(ctx).First(&share, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &share, nil
}

// EnabledByToken reads the share a public link names, and only while it is
// enabled: turning a share off is how an operator stops a form being filled in.
func (s *FormShareService) EnabledByToken(ctx context.Context, token string) (*models.FormShare, error) {
	var share models.FormShare
	if err := s.db(ctx).First(&share, "token = ? AND enabled = ?", token, true).Error; err != nil {
		return nil, err
	}
	return &share, nil
}

// Create saves a new share.
func (s *FormShareService) Create(ctx context.Context, share *models.FormShare) error {
	return s.db(ctx).Create(share).Error
}

// Save writes a share the caller has changed.
func (s *FormShareService) Save(ctx context.Context, share *models.FormShare) error {
	return s.db(ctx).Save(share).Error
}

// Delete soft-deletes a share, which is what stops its token working.
func (s *FormShareService) Delete(ctx context.Context, id string) error {
	return s.db(ctx).Delete(&models.FormShare{}, "id = ?", id).Error
}

// CountSubmission adds one to a share's submission count.
//
// The increment is in SQL. It used to write share.SubmissionCount + 1 from the
// row the request had read, so two visitors submitting at the same moment both
// wrote the same number and the count lost one of them. A public form is exactly
// where that happens: nobody is taking turns.
func (s *FormShareService) CountSubmission(ctx context.Context, share *models.FormShare) error {
	return s.db(ctx).Model(share).UpdateColumns(map[string]interface{}{
		"submission_count": gorm.Expr("submission_count + 1"),
		"updated_at":       time.Now(),
	}).Error
}

// RecordSubmission writes the audit row for one submission.
func (s *FormShareService) RecordSubmission(ctx context.Context, row *models.FormSubmission) error {
	return s.db(ctx).Create(row).Error
}

// Submissions returns the most recent audit rows, optionally for one share or
// one resource. limit caps them, because this is a trail rather than a report.
func (s *FormShareService) Submissions(ctx context.Context, shareID, resourceName string, limit int) ([]models.FormSubmission, error) {
	if limit <= 0 {
		limit = 100
	}
	q := s.db(ctx).Order("created_at DESC, id DESC").Limit(limit)
	if shareID != "" {
		q = q.Where("share_id = ?", shareID)
	}
	if resourceName != "" {
		q = q.Where("resource_name = ?", resourceName)
	}
	var rows []models.FormSubmission
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
