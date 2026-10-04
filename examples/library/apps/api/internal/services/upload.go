package services

import (
	"context"

	"gorm.io/gorm"

	"library/apps/api/internal/models"
	"library/apps/api/internal/paginate"
)

// UploadService owns the uploads table.
//
// The handler decides whether a caller may see everybody's files or only their
// own, because that reads the request; it then says which by passing an ownerID,
// and every method here is scoped by it. The storage bucket is the handler's too:
// a row and an object are two systems, and only the caller knows which order to
// fail in.
type UploadService struct {
	DB *gorm.DB
}

// scope is the uploads ownerID may reach. An empty ownerID means every upload,
// which is what an ADMIN or a holder of the matching uploads permission gets.
//
// Before this scoping existed, GET and DELETE /uploads/:id handed the path
// segment to GORM as a raw condition and List covered every user's files.
func (s *UploadService) scope(ctx context.Context, ownerID string) *gorm.DB {
	q := s.DB.WithContext(ctx).Model(&models.Upload{})
	if ownerID == "" {
		return q
	}
	return q.Where("user_id = ?", ownerID)
}

// uploadListConfig is what the files screen may search, sort and filter by.
// Whitelisted, because each name ends up in SQL.
var uploadListConfig = paginate.Config{
	Searchable:   []string{"filename", "mime_type"},
	Sortable:     map[string]bool{"created_at": true, "filename": true, "size": true, "mime_type": true},
	Filterable:   map[string]bool{"mime_type": true},
	DefaultSort:  "created_at",
	DefaultOrder: "desc",
}

// List returns one page of uploads.
//
// Through paginate, like every other list in the framework: this one clamped the
// page and did the offset arithmetic itself, which meant it answered none of the
// search, sort or counts parameters the rest of the API does.
func (s *UploadService) List(ctx context.Context, p paginate.Params, ownerID string) (paginate.Result[models.Upload], error) {
	return paginate.List[models.Upload](s.scope(ctx, ownerID), p, uploadListConfig)
}

// ByID reads one upload the caller may reach.
func (s *UploadService) ByID(ctx context.Context, id, ownerID string) (*models.Upload, error) {
	var upload models.Upload
	if err := s.scope(ctx, ownerID).Where("id = ?", id).First(&upload).Error; err != nil {
		return nil, err
	}
	return &upload, nil
}

// Record saves the row for a stored object.
func (s *UploadService) Record(ctx context.Context, upload *models.Upload) error {
	return s.DB.WithContext(ctx).Create(upload).Error
}

// Delete removes the row. The object in the bucket is the caller's to delete.
func (s *UploadService) Delete(ctx context.Context, upload *models.Upload) error {
	return s.DB.WithContext(ctx).Delete(upload).Error
}

// PathRecorded reports whether an object already has a row.
//
// Once per key: a second row for one object would let one delete remove a file
// another row still points at.
func (s *UploadService) PathRecorded(ctx context.Context, path string) (bool, error) {
	var n int64
	if err := s.DB.WithContext(ctx).Model(&models.Upload{}).Where("path = ?", path).Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

// UploadKind is one bucket of the files screen's breakdown.
type UploadKind struct {
	Kind  string `gorm:"column:kind" json:"kind"`
	Count int64  `gorm:"column:count" json:"count"`
	Size  int64  `gorm:"column:size" json:"size"`
}

// UploadStats is what the files screen shows above the list.
type UploadStats struct {
	TotalCount int64        `json:"total_count"`
	TotalSize  int64        `json:"total_size"`
	ByKind     []UploadKind `json:"by_kind"`
}

// uploadKindExpr buckets a MIME type. CASE in SQL keeps the breakdown a single
// scan on both Postgres and SQLite.
const uploadKindExpr = `CASE
		WHEN mime_type LIKE 'image/%' THEN 'image'
		WHEN mime_type LIKE 'video/%' THEN 'video'
		WHEN mime_type LIKE 'audio/%' THEN 'audio'
		WHEN mime_type = 'application/pdf' THEN 'pdf'
		WHEN mime_type LIKE '%spreadsheet%' OR mime_type LIKE '%excel%' OR mime_type = 'text/csv' THEN 'spreadsheet'
		WHEN mime_type LIKE '%wordprocessing%' OR mime_type = 'application/msword' THEN 'document'
		ELSE 'other'
	END`

// Stats counts the caller's files, their bytes, and how they break down by kind.
//
// Every one of the three reports its error. The total bytes and the breakdown
// used to discard theirs, so a failure drew a storage page of zeros, which reads
// as an empty bucket rather than as a question the database did not answer.
func (s *UploadService) Stats(ctx context.Context, ownerID string) (UploadStats, error) {
	var stats UploadStats
	if err := s.scope(ctx, ownerID).Count(&stats.TotalCount).Error; err != nil {
		return stats, err
	}
	if err := s.scope(ctx, ownerID).Select("COALESCE(SUM(size), 0)").Scan(&stats.TotalSize).Error; err != nil {
		return stats, err
	}
	rows := []UploadKind{}
	if err := s.scope(ctx, ownerID).
		Select(uploadKindExpr + " AS kind, COUNT(*) AS count, COALESCE(SUM(size), 0) AS size").
		Group("kind").Scan(&rows).Error; err != nil {
		return stats, err
	}
	stats.ByKind = rows
	return stats, nil
}
