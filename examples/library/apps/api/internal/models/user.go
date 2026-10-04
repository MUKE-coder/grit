package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"library/apps/api/internal/ids"

	"library/apps/api/internal/appendonly"
	"library/apps/api/internal/crypto"
)

// Role constants
const (
	RoleAdmin  = "ADMIN"
	RoleEditor = "EDITOR"
	RoleUser   = "USER"
	// grit:roles
)

// User represents a user in the system.
type User struct {
	ID        string                 `gorm:"primarykey;size:36" json:"id"`
	FirstName string                 `gorm:"size:255;not null" json:"first_name" binding:"required"`
	LastName  string                 `gorm:"size:255;not null" json:"last_name" binding:"required"`
	Email     string                 `gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string                 `gorm:"size:255" json:"-"`
	Role      string                 `gorm:"size:20;default:USER" json:"role"`
	Avatar    string                 `gorm:"size:500" json:"avatar"`
	JobTitle  string                 `gorm:"size:255" json:"job_title"`
	Bio       crypto.EncryptedString `gorm:"type:text" json:"bio"`
	// No gorm default on this bool. GORM omits zero-valued fields from an
	// INSERT when the column carries a default, so default:true made
	// Active:false unstorable on create — an admin creating a deactivated user
	// silently got an active one. Every create path sets this explicitly.
	Active          bool       `gorm:"" json:"active"`
	Provider        string     `gorm:"size:50;default:'local'" json:"provider"`
	GoogleID        string     `gorm:"size:255" json:"-"`
	GithubID        string     `gorm:"size:255" json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`

	// Per-account brute-force protection. Sentinel rate-limits by IP, which
	// does nothing against attempts spread across many addresses at one
	// account — the shape of every credential-stuffing run. FailedLoginCount
	// is reset on any successful sign-in.
	FailedLoginCount int            `gorm:"default:0" json:"-"`
	LockedUntil      *time.Time     `json:"locked_until,omitempty"`
	IPAddress        string         `gorm:"size:45" json:"ip_address"`
	MACAddress       string         `gorm:"size:50" json:"mac_address"`
	Version          int            `gorm:"not null;default:1" json:"version"`
	CreatedAt        time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate generates a UUID and hashes the password before saving.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = ids.New()
	}
	if u.Version == 0 {
		u.Version = 1
	}
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect that
// a record they edited has moved on. Pair with the Idempotency-Key
// middleware + /api/sync/push for safe write replay.
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// BeforeCreate generates a UUID for uploads.
func (u *Upload) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = ids.New()
	}
	return nil
}

// CheckPassword compares the given password with the stored hash.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Models returns the ordered list of all models for migration.
// Models with no foreign key dependencies come first.
func Models() []interface{} {
	return []interface{}{
		&User{},
		// Server-side refresh sessions — must exist before anything logs in.
		&Session{},
		&PasswordResetToken{},
		&EmailVerificationToken{},
		&RecoveryContactToken{},
		&RecoveryContact{},
		&Passkey{},
		&WebAuthnSession{},
		&APIKey{},
		// Role/UserRole must migrate before anything authorises a request.
		&Role{},
		&UserRole{},
		&Upload{},
		&Blog{},
		&TwoFactorConfig{},
		&TrustedDevice{},
		&TOTPPendingToken{},
		&MagicLinkToken{},
		&ActivityLog{},
		&WebhookEvent{},
		&FeatureFlag{},
		&FlagExposure{},
		&Notification{},
		// v3.30
		&UserActivity{},
		&AccessReview{},
		&AccessReviewItem{},
		&DeletionJournal{},
		&Ticket{},
		&TicketReply{},
		// v3.31.20 — public form sharing (Phase 2)
		&FormShare{},
		// v3.31.25 — audit log for public submissions
		&FormSubmission{},
		// v3.31.40 — per-user dashboard customisation
		&DashboardLayout{},
		// v3.31.68 — background CSV import job tracking
		&ImportJob{},
		// v3.31.77 — full-database backup index
		&Backup{},
		// backup schedule (period + time-of-day for automatic backups)
		&BackupSchedule{},
		// enterprise SSO: one OIDC connection per customer, plus the external
		// identities linking their users to local accounts
		&SSOConnection{},
		&UserIdentity{},
		// the service provider's own signing keypair, generated on first use
		&SAMLKeypair{},
		&Setting{},
		// The transactional outbox. Registered here like every other table so
		// AutoMigrate creates it and the backup writer includes it; an outbox
		// missing from a backup loses events nobody knows were pending.
		&OutboxMessage{},
		// Sagas: the multi-step processes that have to finish or be undone, and
		// the step-by-step record of each run. Here like every other table, so
		// AutoMigrate creates them and a backup includes them.
		&SagaRun{},
		&SagaStep{},
		&Author{},
		&Book{},
		&Loan{},
		// grit:models
	}
}

// Migrate runs AutoMigrate for every registered model. For tables that
// already exist, GORM ALTERs them to add missing columns — we snapshot
// the column set before and after so the deploy log surfaces exactly
// what changed. Silent migrations are gone: if a column you expected
// didn't land, the diff makes it obvious.
//
//	================================================================
//	DATABASE MIGRATION — 8 model(s) registered
//	================================================================
//	  + created models.Building
//	  ~ models.User — added 2 column(s): is_vip, vip_notes
//	----------------------------------------------------------------
//	Migration done — 1 table(s) created, 1 altered (+2 column(s)), 6 unchanged.
//	================================================================
func Migrate(db *gorm.DB) error {
	models := Models()
	separator := strings.Repeat("=", 64)
	thinSep := strings.Repeat("-", 64)

	log.Println(separator)
	log.Printf("DATABASE MIGRATION: %d model(s) registered", len(models))
	log.Println(separator)

	// Silent logger keeps the schema-inspection SQL noise out of the diff log.
	silentDB := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
	mig := silentDB.Migrator()

	created := 0
	altered := 0
	columnsAdded := 0
	unchanged := 0

	for _, model := range models {
		existed := mig.HasTable(model)

		var before map[string]bool
		if existed {
			before = make(map[string]bool)
			cols, err := mig.ColumnTypes(model)
			if err == nil {
				for _, c := range cols {
					before[c.Name()] = true
				}
			}
		}

		if err := silentDB.AutoMigrate(model); err != nil {
			return fmt.Errorf("migrating %T: %w", model, err)
		}

		if !existed {
			log.Printf("  + created %T", model)
			created++
			continue
		}

		// Diff columns to surface anything AutoMigrate added.
		after, err := mig.ColumnTypes(model)
		if err != nil {
			unchanged++
			continue
		}
		var added []string
		for _, c := range after {
			if !before[c.Name()] {
				added = append(added, c.Name())
			}
		}
		if len(added) == 0 {
			unchanged++
			continue
		}
		log.Printf("  ~ %T: added %d column(s): %s", model, len(added), strings.Join(added, ", "))
		altered++
		columnsAdded += len(added)
	}

	log.Println(thinSep)
	log.Printf("Migration done: %d table(s) created, %d altered (+%d column(s)), %d unchanged.",
		created, altered, columnsAdded, unchanged)

	// Seed the default roles here rather than in database.Seed(): authorization
	// must work on a freshly migrated database, without anyone remembering to
	// run "grit seed". SeedRoles is idempotent and never overwrites an existing
	// role's grants.
	// And at the database, for whatever does not go through GORM: the SQL
	// editor in Studio, psql, a script somebody writes next year.
	if err := appendonly.InstallTriggers(db); err != nil {
		return fmt.Errorf("installing append-only triggers: %w", err)
	}

	if err := SeedRoles(db); err != nil {
		return fmt.Errorf("seeding default roles: %w", err)
	}

	log.Println(separator)
	return nil
}
