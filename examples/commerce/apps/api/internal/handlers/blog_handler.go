package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/sanitize"
	"commerce/apps/api/internal/services"
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
		respond.Fail(c, respond.CodeInternalError, "Failed to fetch blogs")
		return
	}

	// The admin's stat cards ask for their counts on this request, over the
	// same search as the list.
	countQuery := h.DB.WithContext(c.Request.Context()).Model(&models.Blog{})
	if search != "" {
		countQuery = countQuery.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}
	counts, err := paginate.Counts(c, countQuery)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to count blogs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": blogs,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
			"counts":    counts,
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
		respond.Fail(c, respond.CodeInternalError, "Failed to fetch blogs")
		return
	}

	for i := range blogs {
		blogs[i].Content = sanitize.HTML(blogs[i].Content)
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
		respond.Fail(c, respond.CodeNotFound, "Blog not found")
		return
	}

	blog.Content = sanitize.HTML(blog.Content)

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, gin.H{
		"data": blog,
	})
}

// GetByID powers the admin blog detail page. Skips the
// Cache-Control public hint that GetBySlug sets because admin views
// expect fresh data after each save.
func (h *BlogHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	blog, err := h.Service.GetByID(id)
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "Blog not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": blog})
}

// A new post.
type CreateBlogRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Excerpt   string `json:"excerpt"`
	Published *bool  `json:"published"`
}

// Create adds a new blog (admin).

func (h *BlogHandler) Create(c *gin.Context) {
	var req CreateBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
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
		respond.Fail(c, respond.CodeInternalError, "Failed to create blog")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    blog,
		"message": "Blog created successfully",
	})
}

// Changes to a post.
type UpdateBlogRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Excerpt   string `json:"excerpt"`
	Published *bool  `json:"published"`
}

// Update modifies an existing blog (admin).

func (h *BlogHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Fetch existing blog to check published state
	existing, err := h.Service.GetByID(id)
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "Blog not found")
		return
	}

	var req UpdateBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
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

	blog, err := h.Service.Update(id, updates)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to update blog")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    blog,
		"message": "Blog updated successfully",
	})
}

// Delete soft-deletes a blog (admin).
func (h *BlogHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.Delete(id); err != nil {
		respond.Fail(c, respond.CodeNotFound, "Blog not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Blog deleted successfully",
	})
}
