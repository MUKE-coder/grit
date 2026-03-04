package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopGoServiceFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "internal", "service", "auth.go"):    desktopAuthService(),
		filepath.Join(root, "internal", "service", "blog.go"):    desktopBlogService(),
		filepath.Join(root, "internal", "service", "contact.go"): desktopContactService(),
		filepath.Join(root, "internal", "service", "export.go"):  desktopExportService(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "<MODULE>", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func desktopAuthService() string {
	return `package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"<MODULE>/internal/models"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(name, email, password string) (*models.AuthResponse, error) {
	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("email already registered")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "USER",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token := generateToken()
	return &models.AuthResponse{User: user, Token: token}, nil
}

func (s *AuthService) Login(email, password string) (*models.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	token := generateToken()
	return &models.AuthResponse{User: user, Token: token}, nil
}

func (s *AuthService) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (s *AuthService) UpdateUser(id uint, name, email string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	user.Name = name
	user.Email = email
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &user, nil
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
`
}

func desktopBlogService() string {
	return `package service

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"<MODULE>/internal/models"
	"gorm.io/gorm"
)

type BlogService struct {
	db *gorm.DB
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{db: db}
}

func (s *BlogService) List(page, pageSize int, search string) (*models.PaginatedResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64
	query := s.db.Model(&models.Blog{})
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)

	var blogs []models.Blog
	offset := (page - 1) * pageSize
	if err := query.Preload("Author").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&blogs).Error; err != nil {
		return nil, fmt.Errorf("failed to list blogs: %w", err)
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))
	return &models.PaginatedResult{
		Data:     blogs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

func (s *BlogService) ListAll() ([]models.Blog, error) {
	var blogs []models.Blog
	if err := s.db.Preload("Author").Order("created_at DESC").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := s.db.Preload("Author").First(&blog, id).Error; err != nil {
		return nil, fmt.Errorf("blog not found")
	}
	return &blog, nil
}

func (s *BlogService) Create(input models.BlogInput) (*models.Blog, error) {
	blog := models.Blog{
		Title:     input.Title,
		Slug:      slugify(input.Title),
		Content:   input.Content,
		Published: input.Published,
		AuthorID:  input.AuthorID,
	}
	if err := s.db.Create(&blog).Error; err != nil {
		return nil, fmt.Errorf("failed to create blog: %w", err)
	}
	s.db.Preload("Author").First(&blog, blog.ID)
	return &blog, nil
}

func (s *BlogService) Update(id uint, input models.BlogInput) (*models.Blog, error) {
	var blog models.Blog
	if err := s.db.First(&blog, id).Error; err != nil {
		return nil, fmt.Errorf("blog not found")
	}
	blog.Title = input.Title
	blog.Slug = slugify(input.Title)
	blog.Content = input.Content
	blog.Published = input.Published
	if err := s.db.Save(&blog).Error; err != nil {
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}
	s.db.Preload("Author").First(&blog, blog.ID)
	return &blog, nil
}

func (s *BlogService) Delete(id uint) error {
	if err := s.db.Delete(&models.Blog{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	return nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(` + "`" + `[^a-z0-9]+` + "`" + `)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
`
}

func desktopContactService() string {
	return `package service

import (
	"fmt"
	"math"

	"<MODULE>/internal/models"
	"gorm.io/gorm"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{db: db}
}

func (s *ContactService) List(page, pageSize int, search string) (*models.PaginatedResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64
	query := s.db.Model(&models.Contact{})
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR company LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)

	var contacts []models.Contact
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&contacts).Error; err != nil {
		return nil, fmt.Errorf("failed to list contacts: %w", err)
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))
	return &models.PaginatedResult{
		Data:     contacts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

func (s *ContactService) ListAll() ([]models.Contact, error) {
	var contacts []models.Contact
	if err := s.db.Order("created_at DESC").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *ContactService) GetByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, fmt.Errorf("contact not found")
	}
	return &contact, nil
}

func (s *ContactService) Create(input models.ContactInput) (*models.Contact, error) {
	contact := models.Contact{
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Company: input.Company,
		Notes:   input.Notes,
	}
	if err := s.db.Create(&contact).Error; err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}
	return &contact, nil
}

func (s *ContactService) Update(id uint, input models.ContactInput) (*models.Contact, error) {
	var contact models.Contact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, fmt.Errorf("contact not found")
	}
	contact.Name = input.Name
	contact.Email = input.Email
	contact.Phone = input.Phone
	contact.Company = input.Company
	contact.Notes = input.Notes
	if err := s.db.Save(&contact).Error; err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}
	return &contact, nil
}

func (s *ContactService) Delete(id uint) error {
	if err := s.db.Delete(&models.Contact{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	return nil
}
`
}

func desktopExportService() string {
	return `package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

func (s *ExportService) ExportPDF(title string, headers []string, rows [][]string) (string, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 16)
	pdf.AddPage()
	pdf.Cell(0, 10, title)
	pdf.Ln(15)

	// Table headers
	pdf.SetFont("Arial", "B", 10)
	colWidth := 270.0 / float64(len(headers))
	for _, h := range headers {
		pdf.CellFormat(colWidth, 8, h, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	for _, row := range rows {
		for _, cell := range row {
			pdf.CellFormat(colWidth, 7, cell, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Write to temp file
	tmpDir := os.TempDir()
	filename := fmt.Sprintf("%s_%s.pdf", title, time.Now().Format("20060102_150405"))
	path := filepath.Join(tmpDir, filename)

	if err := pdf.OutputFileAndClose(path); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return path, nil
}

func (s *ExportService) ExportExcel(title string, headers []string, rows [][]string) (string, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	// Write headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Style headers bold
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	for i := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellStyle(sheet, cell, cell, style)
	}

	// Write rows
	for r, row := range rows {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	// Write to temp file
	tmpDir := os.TempDir()
	filename := fmt.Sprintf("%s_%s.xlsx", title, time.Now().Format("20060102_150405"))
	path := filepath.Join(tmpDir, filename)

	if err := f.SaveAs(path); err != nil {
		return "", fmt.Errorf("failed to generate Excel: %w", err)
	}

	return path, nil
}
`
}
