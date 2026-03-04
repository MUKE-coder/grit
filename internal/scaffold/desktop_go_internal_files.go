package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopGoInternalFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "internal", "config", "config.go"):  desktopConfigGo(),
		filepath.Join(root, "internal", "db", "db.go"):          desktopDBGo(),
		filepath.Join(root, "internal", "models", "user.go"):    desktopUserModel(),
		filepath.Join(root, "internal", "models", "blog.go"):    desktopBlogModel(),
		filepath.Join(root, "internal", "models", "contact.go"): desktopContactModel(),
		filepath.Join(root, "internal", "models", "types.go"):   desktopTypesGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "<MODULE>", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func desktopConfigGo() string {
	return `package config

import "os"

type Config struct {
	DBDriver string
	DBDSN    string
	AppName  string
}

func Load() *Config {
	return &Config{
		DBDriver: getEnv("DB_DRIVER", "sqlite"),
		DBDSN:    getEnv("DB_DSN", "app.db"),
		AppName:  getEnv("APP_NAME", "Grit Desktop"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
`
}

func desktopDBGo() string {
	return `package db

import (
	"log"

	"<MODULE>/internal/config"
	"<MODULE>/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	if cfg.DBDriver == "postgres" {
		db, err = gorm.Open(postgres.Open(cfg.DBDSN), gormCfg)
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DBDSN), gormCfg)
	}

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Blog{},
		&models.Contact{},
		// grit:models
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Seed default admin user if no users exist
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		admin := models.User{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: "admin123",
			Role:     "ADMIN",
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("warning: failed to seed admin user: %v", err)
		}
	}

	return db
}
`
}

func desktopUserModel() string {
	return `package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Name      string         ` + "`" + `gorm:"size:200;not null" json:"name"` + "`" + `
	Email     string         ` + "`" + `gorm:"size:200;uniqueIndex;not null" json:"email"` + "`" + `
	Password  string         ` + "`" + `gorm:"size:200;not null" json:"-"` + "`" + `
	Role      string         ` + "`" + `gorm:"size:50;default:USER" json:"role"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
`
}

func desktopBlogModel() string {
	return `package models

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Title     string         ` + "`" + `gorm:"size:300;not null" json:"title"` + "`" + `
	Slug      string         ` + "`" + `gorm:"size:300;uniqueIndex;not null" json:"slug"` + "`" + `
	Content   string         ` + "`" + `gorm:"type:text" json:"content"` + "`" + `
	Published bool           ` + "`" + `gorm:"default:false" json:"published"` + "`" + `
	AuthorID  uint           ` + "`" + `gorm:"index" json:"author_id"` + "`" + `
	Author    User           ` + "`" + `gorm:"foreignKey:AuthorID" json:"author,omitempty"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}
`
}

func desktopContactModel() string {
	return `package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Name      string         ` + "`" + `gorm:"size:200;not null" json:"name"` + "`" + `
	Email     string         ` + "`" + `gorm:"size:200" json:"email"` + "`" + `
	Phone     string         ` + "`" + `gorm:"size:50" json:"phone"` + "`" + `
	Company   string         ` + "`" + `gorm:"size:200" json:"company"` + "`" + `
	Notes     string         ` + "`" + `gorm:"type:text" json:"notes"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}
`
}

func desktopTypesGo() string {
	return `package models

type AuthResponse struct {
	User  User   ` + "`" + `json:"user"` + "`" + `
	Token string ` + "`" + `json:"token"` + "`" + `
}

type PaginatedResult struct {
	Data     interface{} ` + "`" + `json:"data"` + "`" + `
	Total    int64       ` + "`" + `json:"total"` + "`" + `
	Page     int         ` + "`" + `json:"page"` + "`" + `
	PageSize int         ` + "`" + `json:"page_size"` + "`" + `
	Pages    int         ` + "`" + `json:"pages"` + "`" + `
}

type BlogInput struct {
	Title     string ` + "`" + `json:"title"` + "`" + `
	Content   string ` + "`" + `json:"content"` + "`" + `
	Published bool   ` + "`" + `json:"published"` + "`" + `
	AuthorID  uint   ` + "`" + `json:"author_id"` + "`" + `
}

type ContactInput struct {
	Name    string ` + "`" + `json:"name"` + "`" + `
	Email   string ` + "`" + `json:"email"` + "`" + `
	Phone   string ` + "`" + `json:"phone"` + "`" + `
	Company string ` + "`" + `json:"company"` + "`" + `
	Notes   string ` + "`" + `json:"notes"` + "`" + `
}

// grit:input-types
`
}
