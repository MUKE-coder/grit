package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/erasure"
	"library/apps/api/internal/models"
)

// ErrUserNotFound is returned when the target of an export or erasure doesn't
// exist (or was already erased and hard-deleted).
var ErrUserNotFound = errors.New("user not found")

// UserExport is the right-to-access bundle: everything the system holds about
// one person, gathered from every table that references them.
type UserExport struct {
	GeneratedAt time.Time               `json:"generated_at"`
	Profile     models.User             `json:"profile"`
	Uploads     []models.Upload         `json:"uploads"`
	Sessions    []models.Session        `json:"sessions"`
	Layout      *models.DashboardLayout `json:"dashboard_layout,omitempty"`
	TwoFactor   struct {
		Enabled bool `json:"enabled"`
	} `json:"two_factor"`

	// The activity log is not held here. A long-lived account has tens of
	// thousands of rows, and the export used to read every one of them into
	// memory before writing the first byte. WriteJSON reads them a page at a
	// time as it writes, with the database handle the export was made with.
	db     *gorm.DB
	userID string
}

// ExportUserData collects a full copy of a user's data. The User's Password and
// OAuth ids are hidden by their json:"-" tags, so the bundle carries the
// person's data without leaking secrets. Bind db to the request's context, so
// an export the client abandoned stops reading.
//
// Every read's error is returned: an export that quietly left out a table would
// be handed to the person as complete, and a right-of-access answer that is
// missing data is the failure GDPR Art. 15 exists to prevent.
func ExportUserData(db *gorm.DB, userID string) (*UserExport, error) {
	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Scrub secrets from the in-memory copy too. The json:"-" tags keep these out
	// of the HTTP response, but clearing them here means no code path can leak the
	// password hash or OAuth ids out of the bundle.
	user.Password = ""
	user.GoogleID = ""
	user.GithubID = ""

	out := &UserExport{GeneratedAt: time.Now().UTC(), Profile: user, db: db, userID: userID}
	if err := db.Where("user_id = ?", userID).Find(&out.Uploads).Error; err != nil {
		return nil, fmt.Errorf("export uploads: %w", err)
	}
	if err := db.Where("user_id = ?", userID).Find(&out.Sessions).Error; err != nil {
		return nil, fmt.Errorf("export sessions: %w", err)
	}

	var layouts []models.DashboardLayout
	if err := db.Where("user_id = ?", userID).Limit(1).Find(&layouts).Error; err != nil {
		return nil, fmt.Errorf("export dashboard layout: %w", err)
	}
	if len(layouts) == 1 {
		out.Layout = &layouts[0]
	}
	// Only whether two-factor is on: the secret and the backup codes never leave
	// the table, not even into this process's memory.
	var enabled []bool
	if err := db.Model(&models.TwoFactorConfig{}).Where("user_id = ?", userID).Limit(1).Pluck("enabled", &enabled).Error; err != nil {
		return nil, fmt.Errorf("export two-factor: %w", err)
	}
	out.TwoFactor.Enabled = len(enabled) == 1 && enabled[0]
	return out, nil
}

// exportActivityPage is how many activity rows an export holds at once.
const exportActivityPage = 500

// exportHead is UserExport without its methods, for the part of the bundle
// that is marshalled in one piece.
type exportHead UserExport

// WriteJSON writes the bundle to w, the activity log last and a page at a time,
// newest first. Its output is the bundle as one JSON object.
//
// Once the first byte is written a failure cannot change the response status,
// so a failed page leaves the JSON unterminated: the file does not parse, rather
// than parsing as a complete export with rows missing.
func (e *UserExport) WriteJSON(w io.Writer) error {
	head, err := json.Marshal((*exportHead)(e))
	if err != nil {
		return err
	}
	// head is an object; the activity array goes in before its closing brace.
	if _, err := w.Write(head[:len(head)-1]); err != nil {
		return err
	}
	if _, err := io.WriteString(w, `,"activity":[`); err != nil {
		return err
	}

	// Pages by id, newest first. Ids are time-ordered, so this is creation
	// order. A keyset cursor rather than OFFSET, which would have the database
	// count past every row already written before returning the next page.
	first, after := true, ""
	for {
		var page []models.UserActivity
		q := e.db.Where("user_id = ?", e.userID)
		if after != "" {
			q = q.Where("id < ?", after)
		}
		if err := q.Order("id desc").Limit(exportActivityPage).Find(&page).Error; err != nil {
			return fmt.Errorf("export activity: %w", err)
		}
		for i := range page {
			row, err := json.Marshal(&page[i])
			if err != nil {
				return err
			}
			if !first {
				if _, err := io.WriteString(w, ","); err != nil {
					return err
				}
			}
			first = false
			if _, err := w.Write(row); err != nil {
				return err
			}
		}
		if len(page) < exportActivityPage {
			break
		}
		after = page[len(page)-1].ID
	}
	_, err = io.WriteString(w, "]}")
	return err
}

// MarshalJSON is WriteJSON into memory, for a caller that hands the bundle to
// json.Marshal or gin's c.JSON. The rows are still read a page at a time.
func (e *UserExport) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := e.WriteJSON(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// The tables erasure deletes from are registered with internal/erasure:
// the framework's own in models/erasable.go, and every resource generated
// with --owned-by from its model file.

// EraseUser fulfils a right-to-erasure request in one transaction: hard-delete
// the user's child PII, anonymize the user row, and append a tamper-evident
// journal entry. The activity log is deliberately left untouched — its rows hold
// only a UUID, so scrubbing the user row anonymizes them too, and editing them
// would break the audit hash chain.
func EraseUser(db *gorm.DB, targetID, actorID, actorEmail, reason string) (*models.DeletionJournal, error) {
	var journal *models.DeletionJournal

	err := db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, "id = ?", targetID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		// What goes is whatever registered with the erasure package: the
		// framework's own tables and every resource generated with --owned-by.
		// This was a list of nine framework tables, and erasing a user left every
		// app record they owned in place while the journal called it done.
		counts, total, err := erasure.Scrub(tx, targetID)
		if err != nil {
			return err
		}

		countsJSON, _ := json.Marshal(counts)

		var prev models.DeletionJournal
		prevHash := ""
		if err := tx.Order("created_at desc, id desc").First(&prev).Error; err == nil {
			prevHash = prev.Hash
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		j := &models.DeletionJournal{
			DeletedUserID:   targetID,
			ActorID:         actorID,
			ActorEmail:      actorEmail,
			Reason:          reason,
			RecordsAffected: total,
			Counts:          string(countsJSON),
			PrevHash:        prevHash,
			CreatedAt:       time.Now().UTC(),
		}
		j.Hash = JournalHash(prevHash, j)
		if err := tx.Create(j).Error; err != nil {
			return err
		}
		journal = j
		return nil
	})
	if err != nil {
		return nil, err
	}
	return journal, nil
}

// canonicalJournal is the stable serialization hashed into the chain. It
// excludes ID (random, uncorrelated) and PrevHash/Hash (derived) so the digest
// depends only on the entry's content.
func canonicalJournal(j *models.DeletionJournal) string {
	return fmt.Sprintf("%s|%s|%s|%s|%d|%s|%d",
		j.DeletedUserID, j.ActorID, j.ActorEmail, j.Reason,
		j.RecordsAffected, j.Counts, j.CreatedAt.UTC().Unix())
}

// JournalHash returns hex(sha256(prevHash || canonical(entry))) — the same
// construction the activity log uses, so the two chains verify identically.
func JournalHash(prevHash string, j *models.DeletionJournal) string {
	sum := sha256.Sum256([]byte(prevHash + canonicalJournal(j)))
	return hex.EncodeToString(sum[:])
}

// JournalVerification is the result of replaying the deletion journal.
type JournalVerification struct {
	Verified bool   `json:"verified"`
	Count    int    `json:"count"`
	BrokenAt string `json:"broken_at,omitempty"` // id of the first row that fails
}

// VerifyJournalChain replays the whole journal in order and confirms each row's
// stored hash equals a fresh recomputation, and that each PrevHash links to the
// row before it. Any post-hoc edit, reorder, or deletion surfaces here.
func VerifyJournalChain(db *gorm.DB) (*JournalVerification, error) {
	var rows []models.DeletionJournal
	if err := db.Order("created_at asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	// Guard against a tie on created_at producing a nondeterministic order by
	// sorting on the same keys the writer used.
	sort.SliceStable(rows, func(a, b int) bool {
		if rows[a].CreatedAt.Equal(rows[b].CreatedAt) {
			return rows[a].ID < rows[b].ID
		}
		return rows[a].CreatedAt.Before(rows[b].CreatedAt)
	})

	prevHash := ""
	for i := range rows {
		r := rows[i]
		if r.PrevHash != prevHash {
			return &JournalVerification{Verified: false, Count: len(rows), BrokenAt: r.ID}, nil
		}
		if JournalHash(prevHash, &r) != r.Hash {
			return &JournalVerification{Verified: false, Count: len(rows), BrokenAt: r.ID}, nil
		}
		prevHash = r.Hash
	}
	return &JournalVerification{Verified: true, Count: len(rows)}, nil
}
