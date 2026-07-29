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
		filepath.Join(root, "internal", "models", "slug.go"):    desktopSlugGo(),
		filepath.Join(root, "internal", "ids", "ids.go"):        desktopIDsGo(),
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
	// APIPort is the loopback port the embedded REST API listens on, so curl
	// and other clients can reach the same endpoints the webview uses.
	APIPort string
}

func Load() *Config {
	return &Config{
		DBDriver: getEnv("DB_DRIVER", "sqlite"),
		DBDSN:    getEnv("DB_DSN", "app.db"),
		AppName:  getEnv("APP_NAME", "Grit Desktop"),
		// 34999 (not 34115): the Wails dev server binds 34115 during
		// "wails dev", so the embedded API must not collide with it.
		APIPort:  getEnv("API_PORT", "34999"),
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
	"github.com/glebarez/sqlite"
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

	"<MODULE>/internal/ids"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Name      string         ` + "`" + `gorm:"size:200;not null" json:"name"` + "`" + `
	Email     string         ` + "`" + `gorm:"size:200;uniqueIndex;not null" json:"email"` + "`" + `
	Password  string         ` + "`" + `gorm:"size:200;not null" json:"-"` + "`" + `
	Role      string         ` + "`" + `gorm:"size:50;default:USER" json:"role"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = ids.New()
	}
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

	"<MODULE>/internal/ids"
	"gorm.io/gorm"
)

type Blog struct {
	ID        string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Title     string         ` + "`" + `gorm:"size:300;not null" json:"title"` + "`" + `
	Slug      string         ` + "`" + `gorm:"size:300;uniqueIndex;not null" json:"slug"` + "`" + `
	Content   string         ` + "`" + `gorm:"type:text" json:"content"` + "`" + `
	Published bool           ` + "`" + `gorm:"default:false" json:"published"` + "`" + `
	AuthorID  string         ` + "`" + `gorm:"size:36;index" json:"author_id"` + "`" + `
	Author    User           ` + "`" + `gorm:"foreignKey:AuthorID" json:"author,omitempty"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (b *Blog) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = ids.New()
	}
	return nil
}
`
}

func desktopContactModel() string {
	return `package models

import (
	"time"

	"<MODULE>/internal/ids"
	"gorm.io/gorm"
)

type Contact struct {
	ID        string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Name      string         ` + "`" + `gorm:"size:200;not null" json:"name"` + "`" + `
	Email     string         ` + "`" + `gorm:"size:200" json:"email"` + "`" + `
	Phone     string         ` + "`" + `gorm:"size:50" json:"phone"` + "`" + `
	Company   string         ` + "`" + `gorm:"size:200" json:"company"` + "`" + `
	Notes     string         ` + "`" + `gorm:"type:text" json:"notes"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (c *Contact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = ids.New()
	}
	return nil
}
`
}

// desktopSlugGo puts slugify in the models package so a generated model's
// BeforeCreate hook (for a slug field) can call it. The service package has its
// own copy for input handling; models needs one too since Go scopes by package.
func desktopSlugGo() string {
	return `package models

import (
	"regexp"
	"strings"
)

var slugRe = regexp.MustCompile(` + "`" + `[^a-z0-9]+` + "`" + `)

// slugify turns "Hello World!" into "hello-world".
func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
`
}

func desktopTypesGo() string {
	return `package models

import (
	// grit:input-imports
)

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
	AuthorID  string ` + "`" + `json:"author_id"` + "`" + `
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

// desktopIDsGo emits internal/ids/ids.go for the desktop module.
//
// The desktop app is its own Go module, so it needs its own copy rather than
// importing the API's. Keeping both on UUIDv7 matters more here than anywhere
// else: an offline client mints IDs with no server contact, and those rows sync
// into the same tables as server-created ones. If one side generated v4 and the
// other v7, ordering would silently depend on which device created the record.
func desktopIDsGo() string {
	return `// Package ids generates primary keys.
//
// ids.New() returns a UUIDv7: a standard 128-bit UUID any client can generate
// offline with no coordination, with a millisecond timestamp in the high bits
// so IDs sort chronologically and insert with good index locality.
//
// Trade-off: a v7 identifier encodes when it was created. If you expose raw IDs
// publicly you are also publishing creation times.
package ids

import "github.com/google/uuid"

// New returns a time-ordered UUIDv7 string.
//
// uuid.NewV7 only fails if the OS entropy source does, which is a situation
// where nothing else works either. It falls back to v4 rather than returning an
// error into every BeforeCreate hook: an unordered id is a performance
// regression, an empty primary key is data corruption.
func New() string {
	if u, err := uuid.NewV7(); err == nil {
		return u.String()
	}
	return uuid.New().String()
}
`
}
