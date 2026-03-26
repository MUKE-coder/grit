package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeAPIBlogFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "blog.go"):       blogModelGo(),
		filepath.Join(apiRoot, "internal", "models", "helpers.go"):    modelsHelpersGo(),
		filepath.Join(apiRoot, "internal", "services", "blog_service.go"): blogServiceGo(),
		filepath.Join(apiRoot, "internal", "handlers", "blog_handler.go"): blogHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func blogModelGo() string {
	return `package models

import (
	"time"

	"gorm.io/gorm"
)

// Blog represents a blog post in the system.
type Blog struct {
	gorm.Model
	Title       string     ` + "`" + `gorm:"size:255;not null" json:"title" binding:"required"` + "`" + `
	Slug        string     ` + "`" + `gorm:"size:255;uniqueIndex" json:"slug"` + "`" + `
	Content     string     ` + "`" + `gorm:"type:text" json:"content"` + "`" + `
	Image       string     ` + "`" + `gorm:"size:500" json:"image"` + "`" + `
	Excerpt     string     ` + "`" + `gorm:"size:500" json:"excerpt"` + "`" + `
	Published   bool       ` + "`" + `gorm:"default:false" json:"published"` + "`" + `
	PublishedAt *time.Time ` + "`" + `json:"published_at"` + "`" + `
}

// BeforeCreate auto-generates the slug before inserting.
func (b *Blog) BeforeCreate(tx *gorm.DB) error {
	if b.Slug == "" {
		b.Slug = slugify(b.Title)
	}
	return nil
}
`
}

func blogServiceGo() string {
	return `package services

import (
	"fmt"
	"math"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// BlogService handles business logic for blogs.
type BlogService struct {
	DB *gorm.DB
}

// NewBlogService creates a new BlogService instance.
func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{DB: db}
}

// List returns a paginated list of all blogs (admin).
func (s *BlogService) List(page, pageSize int, search, sortKey, sortDir string) ([]models.Blog, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	if sortKey == "" {
		sortKey = "created_at"
	}

	query := s.DB.Model(&models.Blog{})

	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var blogs []models.Blog
	offset := (page - 1) * pageSize
	if err := query.Order(sortKey + " " + sortDir).Offset(offset).Limit(pageSize).Find(&blogs).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("fetching blogs: %w", err)
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))
	return blogs, total, pages, nil
}

// ListPublished returns a paginated list of published blogs (public).
func (s *BlogService) ListPublished(page, pageSize int) ([]models.Blog, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.DB.Model(&models.Blog{}).Where("published = ?", true)

	var total int64
	query.Count(&total)

	var blogs []models.Blog
	offset := (page - 1) * pageSize
	if err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&blogs).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("fetching published blogs: %w", err)
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))
	return blogs, total, pages, nil
}

// GetByID returns a single blog by ID.
func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.First(&blog, id).Error; err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}
	return &blog, nil
}

// GetBySlug returns a single published blog by slug.
func (s *BlogService) GetBySlug(slug string) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.Where("slug = ? AND published = ?", slug, true).First(&blog).Error; err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}
	return &blog, nil
}

// Create creates a new blog.
func (s *BlogService) Create(blog *models.Blog) error {
	if err := s.DB.Create(blog).Error; err != nil {
		return fmt.Errorf("creating blog: %w", err)
	}
	return nil
}

// Update modifies an existing blog.
func (s *BlogService) Update(id uint, data map[string]interface{}) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.First(&blog, id).Error; err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}

	if err := s.DB.Model(&blog).Updates(data).Error; err != nil {
		return nil, fmt.Errorf("updating blog: %w", err)
	}

	s.DB.First(&blog, id)
	return &blog, nil
}

// Delete soft-deletes a blog.
func (s *BlogService) Delete(id uint) error {
	var blog models.Blog
	if err := s.DB.First(&blog, id).Error; err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}
	if err := s.DB.Delete(&blog).Error; err != nil {
		return fmt.Errorf("deleting blog: %w", err)
	}
	return nil
}
`
}

func blogHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// BlogHandler handles blog endpoints.
type BlogHandler struct {
	DB      *gorm.DB
	Service *services.BlogService
}

// NewBlogHandler creates a new BlogHandler instance.
func NewBlogHandler(db *gorm.DB) *BlogHandler {
	return &BlogHandler{
		DB:      db,
		Service: services.NewBlogService(db),
	}
}

// List returns a paginated list of all blogs (admin).
func (h *BlogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	allowedSorts := map[string]bool{
		"id": true, "title": true, "slug": true, "published": true, "published_at": true, "created_at": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	blogs, total, pages, err := h.Service.List(page, pageSize, search, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch blogs",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": blogs,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
		},
	})
}

// ListPublished returns a paginated list of published blogs (public).
func (h *BlogHandler) ListPublished(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	blogs, total, pages, err := h.Service.ListPublished(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch blogs",
			},
		})
		return
	}

	c.Header("Cache-Control", "public, max-age=300")
	c.JSON(http.StatusOK, gin.H{
		"data": blogs,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
		},
	})
}

// GetBySlug returns a single published blog by slug (public).
func (h *BlogHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	blog, err := h.Service.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Blog not found",
			},
		})
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, gin.H{
		"data": blog,
	})
}

// Create adds a new blog (admin).
func (h *BlogHandler) Create(c *gin.Context) {
	var req struct {
		Title     string ` + "`" + `json:"title" binding:"required"` + "`" + `
		Content   string ` + "`" + `json:"content"` + "`" + `
		Image     string ` + "`" + `json:"image"` + "`" + `
		Excerpt   string ` + "`" + `json:"excerpt"` + "`" + `
		Published *bool  ` + "`" + `json:"published"` + "`" + `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	blog := models.Blog{
		Title:   req.Title,
		Content: req.Content,
		Image:   req.Image,
		Excerpt: req.Excerpt,
	}

	if req.Published != nil && *req.Published {
		blog.Published = true
		now := time.Now()
		blog.PublishedAt = &now
	}

	if err := h.Service.Create(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create blog",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    blog,
		"message": "Blog created successfully",
	})
}

// Update modifies an existing blog (admin).
func (h *BlogHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid blog ID",
			},
		})
		return
	}

	// Fetch existing blog to check published state
	existing, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Blog not found",
			},
		})
		return
	}

	var req struct {
		Title     string ` + "`" + `json:"title"` + "`" + `
		Content   string ` + "`" + `json:"content"` + "`" + `
		Image     string ` + "`" + `json:"image"` + "`" + `
		Excerpt   string ` + "`" + `json:"excerpt"` + "`" + `
		Published *bool  ` + "`" + `json:"published"` + "`" + `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Excerpt != "" {
		updates["excerpt"] = req.Excerpt
	}
	if req.Published != nil {
		updates["published"] = *req.Published
		if *req.Published && !existing.Published {
			// Toggling published to true — set PublishedAt
			now := time.Now()
			updates["published_at"] = &now
		} else if !*req.Published && existing.Published {
			// Toggling published to false — clear PublishedAt
			updates["published_at"] = nil
		}
	}

	blog, err := h.Service.Update(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update blog",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    blog,
		"message": "Blog updated successfully",
	})
}

// Delete soft-deletes a blog (admin).
func (h *BlogHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid blog ID",
			},
		})
		return
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Blog not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Blog deleted successfully",
	})
}
`
}

func modelsHelpersGo() string {
	return `package models

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

// slugify generates a URL-friendly slug with a unique suffix.
func slugify(s string) string {
	slug := strings.ToLower(s)
	re := regexp.MustCompile(` + "`" + `[^a-z0-9]+` + "`" + `)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	b := make([]byte, 4)
	rand.Read(b)
	return slug + "-" + hex.EncodeToString(b)
}
`
}
