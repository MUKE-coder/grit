package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeStorageFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "storage", "storage.go"):  storageServiceGo(),
		filepath.Join(apiRoot, "internal", "storage", "image.go"):   storageImageGo(),
		filepath.Join(apiRoot, "internal", "handlers", "upload.go"): uploadHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func storageServiceGo() string {
	return `package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"{{MODULE}}/internal/config"
)

// Storage provides S3-compatible file storage operations.
type Storage struct {
	client *s3.Client
	bucket string
	cfg    config.StorageConfig
}

// New creates a new Storage instance using the given config.
// Works with AWS S3, MinIO, Cloudflare R2, and Backblaze B2.
func New(cfg config.StorageConfig) (*Storage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if cfg.Endpoint != "" {
				return aws.Endpoint{
					URL:               cfg.Endpoint,
					HostnameImmutable: true,
					SigningRegion:     cfg.Region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO
	})

	// Verify bucket exists with a quick head request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		// Try to create the bucket
		_, createErr := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket),
		})
		if createErr != nil {
			return nil, fmt.Errorf("bucket %q not accessible and cannot be created: %w", cfg.Bucket, err)
		}
	}

	// Always ensure public-read policy so uploaded files are accessible via URL.
	// This is idempotent — safe to call on every startup.
	policy := fmt.Sprintf(` + "`" + `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": ["*"]},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}` + "`" + `, cfg.Bucket)

	_, _ = client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(cfg.Bucket),
		Policy: aws.String(policy),
	})

	return &Storage{
		client: client,
		bucket: cfg.Bucket,
		cfg:    cfg,
	}, nil
}

// Upload stores a file in the bucket at the given key.
func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("uploading %q: %w", key, err)
	}
	return nil
}

// Download retrieves a file from the bucket.
func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("downloading %q: %w", key, err)
	}
	return result.Body, nil
}

// Delete removes a file from the bucket.
func (s *Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("deleting %q: %w", key, err)
	}
	return nil
}

// GetURL returns the public URL for a stored file.
func (s *Storage) GetURL(key string) string {
	endpoint := strings.TrimRight(s.cfg.Endpoint, "/")
	// Encode each path segment individually to preserve forward slashes
	segments := strings.Split(key, "/")
	for i, seg := range segments {
		segments[i] = url.PathEscape(seg)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, s.bucket, strings.Join(segments, "/"))
}

// GetSignedURL returns a pre-signed URL valid for the given duration.
func (s *Storage) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	result, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))
	if err != nil {
		return "", fmt.Errorf("generating signed URL for %q: %w", key, err)
	}
	return result.URL, nil
}

// PresignPutURL generates a pre-signed PUT URL for direct browser upload.
func (s *Storage) PresignPutURL(ctx context.Context, key, contentType string) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	result, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("generating presigned PUT URL for %q: %w", key, err)
	}
	return result.URL, nil
}
`
}

func storageImageGo() string {
	return `package storage

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/disintegration/imaging"
)

// MaxImageWidth is the maximum width for processed images.
const MaxImageWidth = 1920

// ThumbnailSize is the size of generated thumbnails.
const ThumbnailSize = 300

// ProcessImage resizes an image if it exceeds MaxImageWidth, preserving aspect ratio.
// Returns the processed image bytes and format.
func ProcessImage(reader io.Reader, mimeType string) ([]byte, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() > MaxImageWidth {
		img = imaging.Resize(img, MaxImageWidth, 0, imaging.Lanczos)
	}

	var buf bytes.Buffer
	if err := encodeImage(&buf, img, mimeType); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateThumbnail creates a square thumbnail of the given size.
func GenerateThumbnail(reader io.Reader, mimeType string) ([]byte, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	thumb := imaging.Fill(img, ThumbnailSize, ThumbnailSize, imaging.Center, imaging.Lanczos)

	var buf bytes.Buffer
	if err := encodeImage(&buf, thumb, mimeType); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// IsImageMimeType returns true if the MIME type is a supported image format.
func IsImageMimeType(mimeType string) bool {
	switch strings.ToLower(mimeType) {
	case "image/jpeg", "image/png", "image/gif":
		return true
	}
	return false
}

func encodeImage(buf *bytes.Buffer, img image.Image, mimeType string) error {
	switch strings.ToLower(mimeType) {
	case "image/png":
		return png.Encode(buf, img)
	default:
		return jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	}
}
`
}

func uploadHandlerGo() string {
	return `package handlers

import (
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/jobs"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/storage"
)

// AllowedMimeTypes defines which file types can be uploaded.
var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"video/mp4":       true,
	"video/webm":      true,
	"video/quicktime": true,
	"application/pdf": true,
	"text/plain":      true,
	"text/csv":        true,
	"application/json": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// MaxUploadSize is the maximum file size (50 MB).
const MaxUploadSize = 50 << 20

// UploadHandler handles file upload endpoints.
type UploadHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	Jobs    *jobs.Client
}

// Create handles file upload via multipart form.
func (h *UploadHandler) Create(c *gin.Context) {
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "STORAGE_UNAVAILABLE",
				"message": "File storage is not configured",
			},
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_FILE",
				"message": "No file provided",
			},
		})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > MaxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "FILE_TOO_LARGE",
				"message": fmt.Sprintf("File size exceeds maximum of %d MB", MaxUploadSize/(1<<20)),
			},
		})
		return
	}

	// Validate MIME type
	mimeType := header.Header.Get("Content-Type")
	if !AllowedMimeTypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_FILE_TYPE",
				"message": "File type not allowed",
			},
		})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.TrimSuffix(filepath.Base(header.Filename), ext), ext)
	key := fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01"), filename)

	// Upload to storage
	if err := h.Storage.Upload(c.Request.Context(), key, file, mimeType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPLOAD_FAILED",
				"message": "Failed to upload file",
			},
		})
		return
	}

	userID, _ := c.Get("user_id")

	upload := models.Upload{
		Filename:     filename,
		OriginalName: header.Filename,
		MimeType:     mimeType,
		Size:         header.Size,
		Path:         key,
		URL:          h.Storage.GetURL(key),
		UserID:       userID.(uint),
	}

	if err := h.DB.Create(&upload).Error; err != nil {
		// Clean up uploaded file on DB error
		_ = h.Storage.Delete(c.Request.Context(), key)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to save upload record",
			},
		})
		return
	}

	// Enqueue image processing job if it's an image
	if h.Jobs != nil && storage.IsImageMimeType(mimeType) {
		_ = h.Jobs.EnqueueProcessImage(upload.ID, key, mimeType)
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    upload,
		"message": "File uploaded successfully",
	})
}

// List returns a paginated list of uploads.
func (h *UploadHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := h.DB.Model(&models.Upload{})

	// Filter by MIME type
	if mimeType := c.Query("mime_type"); mimeType != "" {
		query = query.Where("mime_type LIKE ?", mimeType+"%")
	}

	var total int64
	query.Count(&total)

	var uploads []models.Upload
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&uploads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch uploads",
			},
		})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"data": uploads,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
		},
	})
}

// GetByID returns a single upload by ID.
func (h *UploadHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var upload models.Upload
	if err := h.DB.First(&upload, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Upload not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": upload,
	})
}

// Delete removes an upload and its stored file.
func (h *UploadHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var upload models.Upload
	if err := h.DB.First(&upload, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Upload not found",
			},
		})
		return
	}

	// Delete from storage
	if h.Storage != nil {
		_ = h.Storage.Delete(c.Request.Context(), upload.Path)
		// Also delete thumbnail if it exists
		if upload.ThumbnailURL != "" {
			thumbKey := strings.Replace(upload.Path, "uploads/", "thumbnails/", 1)
			_ = h.Storage.Delete(c.Request.Context(), thumbKey)
		}
	}

	if err := h.DB.Delete(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete upload",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload deleted successfully",
	})
}

// Presign generates a presigned PUT URL for direct browser-to-storage upload.
func (h *UploadHandler) Presign(c *gin.Context) {
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "STORAGE_UNAVAILABLE", "message": "File storage is not configured"},
		})
		return
	}

	var req struct {
		Filename    string ` + "`" + `json:"filename" binding:"required"` + "`" + `
		ContentType string ` + "`" + `json:"content_type" binding:"required"` + "`" + `
		FileSize    int64  ` + "`" + `json:"file_size" binding:"required"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if !AllowedMimeTypes[req.ContentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_FILE_TYPE", "message": "File type not allowed"},
		})
		return
	}

	if req.FileSize > MaxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "FILE_TOO_LARGE", "message": fmt.Sprintf("File size exceeds maximum of %d MB", MaxUploadSize/(1<<20))},
		})
		return
	}

	ext := filepath.Ext(req.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.TrimSuffix(filepath.Base(req.Filename), ext), ext)
	key := fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01"), filename)

	presignedURL, err := h.Storage.PresignPutURL(c.Request.Context(), key, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "PRESIGN_FAILED", "message": "Failed to generate upload URL"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"presigned_url": presignedURL,
			"key":           key,
			"public_url":    h.Storage.GetURL(key),
		},
	})
}

// CompleteUpload records a file that was uploaded directly to storage via presigned URL.
func (h *UploadHandler) CompleteUpload(c *gin.Context) {
	var req struct {
		Key         string ` + "`" + `json:"key" binding:"required"` + "`" + `
		Filename    string ` + "`" + `json:"filename" binding:"required"` + "`" + `
		ContentType string ` + "`" + `json:"content_type" binding:"required"` + "`" + `
		Size        int64  ` + "`" + `json:"size" binding:"required"` + "`" + `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	userID, _ := c.Get("user_id")

	upload := models.Upload{
		Filename:     filepath.Base(req.Key),
		OriginalName: req.Filename,
		MimeType:     req.ContentType,
		Size:         req.Size,
		Path:         req.Key,
		URL:          h.Storage.GetURL(req.Key),
		UserID:       userID.(uint),
	}

	if err := h.DB.Create(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to save upload record"},
		})
		return
	}

	// Enqueue image processing job if it's an image
	if h.Jobs != nil && storage.IsImageMimeType(req.ContentType) {
		_ = h.Jobs.EnqueueProcessImage(upload.ID, req.Key, req.ContentType)
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    upload,
		"message": "Upload recorded successfully",
	})
}
`
}
