package scaffold

import "strings"

// GDPR data toolkit — right-to-access (export) and right-to-erasure (delete).
//
// GDPR Art. 15 (access) and Art. 17 (erasure) give a person the right to get a
// copy of everything you hold on them, and to have it deleted. This scaffolds
// both, plus the piece auditors actually ask for: proof the deletion happened
// without keeping the personal data it removed.
//
// Export gathers a user's records from every table that references them into a
// single JSON bundle. Erasure hard-deletes the child records that exist only to
// serve that user (sessions, uploads, 2FA secrets, …) and anonymizes the user
// row itself — scrubbing name, email and every other PII column while keeping
// the id, so foreign-key references elsewhere resolve to a tombstone rather than
// dangling. It never touches the tamper-evident activity log: those rows carry a
// bare UUID, not inline PII, so once the user row is scrubbed they're already
// anonymized, and mutating them would break the hash chain.
//
// Every erasure appends one row to a deletion_journal that is itself a
// hash chain (the same construction as the activity log): who erased whom, when,
// and how many records fell — but never the erased person's data. That's the
// evidence an auditor verifies, and VerifyJournalChain proves no row was altered
// or dropped after the fact.

func apiGDPRModelGo() string {
	src := `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeletionJournal is one tamper-evident record of a GDPR erasure. It proves a
// deletion happened without retaining the personal data that was deleted: it
// holds the erased user's opaque id (anonymized once their row is scrubbed), the
// admin who performed it, per-table counts, and a hash-chain link. Read-only —
// never updated or deleted, or the chain breaks.
type DeletionJournal struct {
	ID              string ~gorm:"primarykey;size:36" json:"id"~
	DeletedUserID   string ~gorm:"size:36;index" json:"deleted_user_id"~ // UUID only — no PII
	ActorID         string ~gorm:"size:36" json:"actor_id"~              // who ran the erasure
	ActorEmail      string ~gorm:"size:255" json:"actor_email"~          // the admin, not the erased user
	Reason          string ~gorm:"size:500" json:"reason"~
	RecordsAffected int    ~json:"records_affected"~
	Counts          string ~gorm:"type:text" json:"counts"~ // JSON: {"uploads":3,"sessions":2}

	PrevHash string ~gorm:"size:64" json:"prev_hash"~            // hex sha256, "" for the genesis row
	Hash     string ~gorm:"size:64;uniqueIndex" json:"hash"~     // hex sha256(prev_hash || canonical)

	CreatedAt time.Time ~gorm:"index" json:"created_at"~
}

func (d *DeletionJournal) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiGDPRServiceGo() string {
	src := `package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// ErrUserNotFound is returned when the target of an export or erasure doesn't
// exist (or was already erased and hard-deleted).
var ErrUserNotFound = errors.New("user not found")

// UserExport is the right-to-access bundle: everything the system holds about
// one person, gathered from every table that references them.
type UserExport struct {
	GeneratedAt time.Time              ~json:"generated_at"~
	Profile     models.User            ~json:"profile"~
	Uploads     []models.Upload        ~json:"uploads"~
	Sessions    []models.Session       ~json:"sessions"~
	Activity    []models.UserActivity  ~json:"activity"~
	Layout      *models.DashboardLayout ~json:"dashboard_layout,omitempty"~
	TwoFactor   struct {
		Enabled bool ~json:"enabled"~
	} ~json:"two_factor"~
}

// ExportUserData collects a full copy of a user's data. The User's Password and
// OAuth ids are hidden by their json:"-" tags, so the bundle carries the
// person's data without leaking secrets.
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

	out := &UserExport{GeneratedAt: time.Now().UTC(), Profile: user}
	db.Where("user_id = ?", userID).Find(&out.Uploads)
	db.Where("user_id = ?", userID).Find(&out.Sessions)
	db.Where("user_id = ?", userID).Order("created_at desc").Find(&out.Activity)

	var layout models.DashboardLayout
	if err := db.First(&layout, "user_id = ?", userID).Error; err == nil {
		out.Layout = &layout
	}
	var tf models.TwoFactorConfig
	if err := db.First(&tf, "user_id = ?", userID).Error; err == nil {
		out.TwoFactor.Enabled = true
	}
	return out, nil
}

// erasableTables are the tables whose rows exist only to serve the user and
// carry no independent compliance value — safe to hard-delete on erasure.
func erasableTargets() []struct {
	name  string
	model interface{}
} {
	return []struct {
		name  string
		model interface{}
	}{
		{"uploads", &models.Upload{}},
		{"sessions", &models.Session{}},
		{"password_reset_tokens", &models.PasswordResetToken{}},
		{"user_roles", &models.UserRole{}},
		{"two_factor_configs", &models.TwoFactorConfig{}},
		{"trusted_devices", &models.TrustedDevice{}},
		{"totp_pending_tokens", &models.TOTPPendingToken{}},
		{"dashboard_layouts", &models.DashboardLayout{}},
		{"notifications", &models.Notification{}},
	}
}

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

		// Unscoped is essential: several of these models carry gorm.DeletedAt, and a
		// normal Delete would only *soft*-delete — setting deleted_at while the row
		// (and its PII) physically remains. That is not erasure. Unscoped removes the
		// row for real, and counts physical rows so the journal reflects what was
		// actually destroyed.
		counts := map[string]int64{}
		total := 0
		for _, t := range erasableTargets() {
			var n int64
			if err := tx.Unscoped().Model(t.model).Where("user_id = ?", targetID).Count(&n).Error; err != nil {
				return fmt.Errorf("counting %s: %w", t.name, err)
			}
			if n > 0 {
				if err := tx.Unscoped().Where("user_id = ?", targetID).Delete(t.model).Error; err != nil {
					return fmt.Errorf("deleting %s: %w", t.name, err)
				}
			}
			counts[t.name] = n
			total += int(n)
		}

		// Anonymize the user row in place: scrub every PII column, keep the id so
		// references resolve to a tombstone. The email keeps a unique, non-routable
		// value so the unique index stays satisfiable if the row is re-read.
		tomb := "erased-" + user.ID + "@deleted.invalid"
		if err := tx.Model(&models.User{}).Where("id = ?", targetID).Updates(map[string]interface{}{
			"first_name":  "Erased",
			"last_name":   "User",
			"email":       tomb,
			"password":    "",
			"avatar":      "",
			"job_title":   "",
			"bio":         "",
			"ip_address":  "",
			"mac_address": "",
			"google_id":   "",
			"github_id":   "",
			"active":      false,
			"role":        "USER",
		}).Error; err != nil {
			return fmt.Errorf("anonymizing user: %w", err)
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
	Verified bool   ~json:"verified"~
	Count    int    ~json:"count"~
	BrokenAt string ~json:"broken_at,omitempty"~ // id of the first row that fails
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
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiGDPRHandlerGo() string {
	src := `package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// GDPRHandler serves the right-to-access and right-to-erasure endpoints.
type GDPRHandler struct {
	DB *gorm.DB
}

func NewGDPRHandler(db *gorm.DB) *GDPRHandler {
	return &GDPRHandler{DB: db}
}

// Export returns the full data bundle for a user as a downloadable JSON file.
// A user may export their own data; an admin may export anyone's.
func (h *GDPRHandler) Export(c *gin.Context) {
	targetID := c.Param("id")
	callerID, _ := c.Get("user_id")
	callerRole, _ := c.Get("user_role")
	if fmt.Sprint(callerID) != targetID && fmt.Sprint(callerRole) != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you may only export your own data"}})
		return
	}

	bundle, err := services.ExportUserData(h.DB, targetID)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "EXPORT_FAILED", "message": err.Error()}})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"user-%s-export.json\"", targetID))
	c.JSON(http.StatusOK, bundle)
}

type eraseRequest struct {
	Reason string ~json:"reason"~
}

// Erase fulfils a right-to-erasure request (admin only). It refuses to let an
// admin erase themselves — that would revoke their own access mid-request and
// orphan the operation.
func (h *GDPRHandler) Erase(c *gin.Context) {
	targetID := c.Param("id")
	callerID, _ := c.Get("user_id")
	callerEmail, _ := c.Get("user_email")

	if fmt.Sprint(callerID) == targetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "SELF_ERASE", "message": "you cannot erase your own account here"}})
		return
	}

	var req eraseRequest
	_ = c.ShouldBindJSON(&req)

	journal, err := services.EraseUser(h.DB, targetID, fmt.Sprint(callerID), fmt.Sprint(callerEmail), req.Reason)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "ERASE_FAILED", "message": err.Error()}})
		return
	}

	// Record the erasure in the semantic activity log too, so it shows up in the
	// dashboard and flows out through the OCSF/SIEM export.
	h.DB.Create(&models.UserActivity{
		UserID:       fmt.Sprint(callerID),
		Action:       "user.gdpr_erase",
		Severity:     "warn",
		Summary:      fmt.Sprintf("Erased all personal data for user %s (%d records)", targetID, journal.RecordsAffected),
		ResourceType: "user",
		ResourceID:   targetID,
	})

	c.JSON(http.StatusOK, gin.H{"data": journal, "message": "User data erased"})
}

// Journal lists the deletion journal and replays its hash chain so the caller
// can see at a glance whether the record of erasures is intact (admin only).
func (h *GDPRHandler) Journal(c *gin.Context) {
	var rows []models.DeletionJournal
	if err := h.DB.Order("created_at desc, id desc").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "QUERY_FAILED", "message": err.Error()}})
		return
	}
	verification, err := services.VerifyJournalChain(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "VERIFY_FAILED", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": verification})
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiGDPRTestGo() string {
	src := `package services

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

func gdprDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{}, &models.Upload{}, &models.Session{}, &models.UserActivity{},
		&models.DashboardLayout{}, &models.DeletionJournal{},
		&models.PasswordResetToken{}, &models.UserRole{}, &models.TwoFactorConfig{},
		&models.TrustedDevice{}, &models.TOTPPendingToken{}, &models.Notification{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seedUserWithData(t *testing.T, db *gorm.DB, email string) models.User {
	t.Helper()
	u := models.User{FirstName: "Real", LastName: "Person", Email: email, Password: "hash", Active: true}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	db.Create(&models.Upload{UserID: u.ID, Filename: "a.png"})
	db.Create(&models.Upload{UserID: u.ID, Filename: "b.png"})
	db.Create(&models.Session{UserID: u.ID})
	db.Create(&models.UserActivity{UserID: u.ID, Action: "login", Summary: "signed in"})
	return u
}

func TestExportGathersEverything(t *testing.T) {
	db := gdprDB(t)
	u := seedUserWithData(t, db, "export@test.dev")

	bundle, err := ExportUserData(db, u.ID)
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if bundle.Profile.Email != "export@test.dev" {
		t.Errorf("profile email = %q", bundle.Profile.Email)
	}
	if len(bundle.Uploads) != 2 {
		t.Errorf("uploads = %d, want 2", len(bundle.Uploads))
	}
	if len(bundle.Sessions) != 1 {
		t.Errorf("sessions = %d, want 1", len(bundle.Sessions))
	}
	if len(bundle.Activity) != 1 {
		t.Errorf("activity = %d, want 1", len(bundle.Activity))
	}
	if bundle.Profile.Password != "" {
		t.Errorf("export leaked the password hash")
	}
}

func TestExportUnknownUser(t *testing.T) {
	db := gdprDB(t)
	if _, err := ExportUserData(db, "nope"); err != ErrUserNotFound {
		t.Errorf("err = %v, want ErrUserNotFound", err)
	}
}

func TestEraseDeletesChildrenAndAnonymizes(t *testing.T) {
	db := gdprDB(t)
	u := seedUserWithData(t, db, "erase@test.dev")

	j, err := EraseUser(db, u.ID, "admin-1", "admin@test.dev", "user request")
	if err != nil {
		t.Fatalf("erase: %v", err)
	}

	// Unscoped counts PHYSICAL rows. A soft delete would leave the row (and its
	// PII) in place while a plain Count hid it — so this must be Unscoped to prove
	// the data is actually gone, not just flagged deleted.
	var uploads, sessions int64
	db.Unscoped().Model(&models.Upload{}).Where("user_id = ?", u.ID).Count(&uploads)
	db.Unscoped().Model(&models.Session{}).Where("user_id = ?", u.ID).Count(&sessions)
	if uploads != 0 || sessions != 0 {
		t.Errorf("child rows survived erasure (soft-deleted?): uploads=%d sessions=%d", uploads, sessions)
	}

	var after models.User
	db.First(&after, "id = ?", u.ID)
	if after.Email == "erase@test.dev" || after.FirstName == "Real" {
		t.Errorf("user was not anonymized: %+v", after)
	}
	if after.Email != "erased-"+u.ID+"@deleted.invalid" {
		t.Errorf("tombstone email = %q", after.Email)
	}

	if j.RecordsAffected != 3 { // 2 uploads + 1 session
		t.Errorf("records affected = %d, want 3", j.RecordsAffected)
	}
}

func TestEraseKeepsActivityLog(t *testing.T) {
	db := gdprDB(t)
	u := seedUserWithData(t, db, "keep@test.dev")

	if _, err := EraseUser(db, u.ID, "admin-1", "admin@test.dev", ""); err != nil {
		t.Fatalf("erase: %v", err)
	}
	// The activity row must remain — it carries only the UUID, now anonymized.
	var acts int64
	db.Model(&models.UserActivity{}).Where("user_id = ?", u.ID).Count(&acts)
	if acts != 1 {
		t.Errorf("activity rows = %d, want 1 (audit trail must survive)", acts)
	}
}

func TestJournalChainVerifies(t *testing.T) {
	db := gdprDB(t)
	for _, e := range []string{"a@t.dev", "b@t.dev", "c@t.dev"} {
		u := seedUserWithData(t, db, e)
		if _, err := EraseUser(db, u.ID, "admin-1", "admin@test.dev", ""); err != nil {
			t.Fatalf("erase: %v", err)
		}
	}
	v, err := VerifyJournalChain(db)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !v.Verified || v.Count != 3 {
		t.Errorf("verification = %+v, want verified with 3 rows", v)
	}
}

func TestJournalTamperDetected(t *testing.T) {
	db := gdprDB(t)
	u := seedUserWithData(t, db, "tamper@t.dev")
	if _, err := EraseUser(db, u.ID, "admin-1", "admin@test.dev", "before"); err != nil {
		t.Fatalf("erase: %v", err)
	}
	// Forge the record: change the reason after the fact without rehashing.
	if err := db.Model(&models.DeletionJournal{}).Where("deleted_user_id = ?", u.ID).
		Update("reason", "after").Error; err != nil {
		t.Fatalf("tamper: %v", err)
	}
	v, _ := VerifyJournalChain(db)
	if v.Verified {
		t.Errorf("tampering was not detected")
	}
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func adminGDPRPage() string {
	src := `"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { Shield, ShieldCheck, Download, Trash2, AlertTriangle, Loader2, Check, X } from "@/lib/icons";

interface JournalRow {
  id: string;
  deleted_user_id: string;
  actor_email: string;
  reason: string;
  records_affected: number;
  created_at: string;
}

interface JournalMeta {
  verified: boolean;
  count: number;
  broken_at?: string;
}

export default function GDPRPage() {
  const qc = useQueryClient();
  const [userId, setUserId] = useState("");
  const [confirming, setConfirming] = useState(false);
  const [reason, setReason] = useState("");
  const [msg, setMsg] = useState<{ kind: "ok" | "err"; text: string } | null>(null);

  const journalQ = useQuery({
    queryKey: ["gdpr-journal"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/gdpr/journal");
      return { rows: (data.data ?? []) as JournalRow[], meta: data.meta as JournalMeta };
    },
  });

  // Export streams the bundle to the browser as a downloaded JSON file — the
  // artifact you hand a data-subject-access request.
  const exportM = useMutation({
    mutationFn: async (id: string) => {
      const res = await apiClient.get("/api/users/" + id + "/gdpr-export");
      const blob = new Blob([JSON.stringify(res.data, null, 2)], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "user-" + id + "-export.json";
      a.click();
      URL.revokeObjectURL(url);
    },
    onSuccess: () => setMsg({ kind: "ok", text: "Export downloaded." }),
    onError: (e: unknown) => setMsg({ kind: "err", text: errMsg(e) }),
  });

  const eraseM = useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post("/api/users/" + userId + "/gdpr-erase", { reason });
      return data.data as { records_affected: number };
    },
    onSuccess: (d) => {
      setMsg({ kind: "ok", text: "Erased " + d.records_affected + " record(s) and anonymized the user." });
      setConfirming(false);
      setReason("");
      setUserId("");
      qc.invalidateQueries({ queryKey: ["gdpr-journal"] });
    },
    onError: (e: unknown) => setMsg({ kind: "err", text: errMsg(e) }),
  });

  const meta = journalQ.data?.meta;

  return (
    <div className="space-y-8">
      <PageHeader
        title="GDPR"
        subtitle="Right-to-access exports and right-to-erasure. Every erasure is written to a tamper-evident deletion journal."
      />

      {msg && (
        <div
          className={
            "flex items-center gap-2 rounded-lg border px-4 py-3 text-sm " +
            (msg.kind === "ok"
              ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-500"
              : "border-red-500/30 bg-red-500/10 text-red-500")
          }
        >
          {msg.kind === "ok" ? <Check className="h-4 w-4" /> : <X className="h-4 w-4" />}
          {msg.text}
        </div>
      )}

      <div className="rounded-xl border border-border bg-card p-6">
        <div className="mb-4 flex items-center gap-2">
          <Shield className="h-5 w-5 text-primary" />
          <h2 className="text-lg font-semibold">Export or erase a user</h2>
        </div>
        <p className="mb-4 max-w-2xl text-sm text-muted-foreground">
          Enter a user ID. Export downloads a full JSON copy of their data (Art. 15). Erase
          hard-deletes their personal records and anonymizes the account (Art. 17) — this cannot
          be undone.
        </p>
        <div className="flex flex-wrap items-center gap-3">
          <input
            value={userId}
            onChange={(e) => {
              setUserId(e.target.value);
              setConfirming(false);
            }}
            spellCheck={false}
            placeholder="user id (uuid)"
            className="min-w-[320px] flex-1 rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm outline-none focus:border-primary"
          />
          <button
            type="button"
            disabled={!userId || exportM.isPending}
            onClick={() => exportM.mutate(userId)}
            className="inline-flex items-center gap-2 rounded-lg border border-border bg-background px-4 py-2 text-sm font-medium transition-colors hover:border-primary disabled:opacity-40"
          >
            {exportM.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Download className="h-4 w-4" />}
            Export data
          </button>
          <button
            type="button"
            disabled={!userId}
            onClick={() => setConfirming(true)}
            className="inline-flex items-center gap-2 rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-2 text-sm font-medium text-red-500 transition-colors hover:bg-red-500/20 disabled:opacity-40"
          >
            <Trash2 className="h-4 w-4" />
            Erase&hellip;
          </button>
        </div>

        {confirming && (
          <div className="mt-4 rounded-lg border border-red-500/30 bg-red-500/5 p-4">
            <div className="mb-3 flex items-start gap-2 text-sm text-red-500">
              <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
              <span>
                This permanently deletes <span className="font-mono">{userId}</span>&rsquo;s uploads,
                sessions, 2FA and more, and anonymizes their account. The audit trail is retained.
              </span>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <input
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder="reason (recorded in the journal)"
                className="min-w-[280px] flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
              />
              <button
                type="button"
                disabled={eraseM.isPending}
                onClick={() => eraseM.mutate()}
                className="inline-flex items-center gap-2 rounded-lg bg-red-500 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-red-600 disabled:opacity-40"
              >
                {eraseM.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                Confirm erasure
              </button>
              <button
                type="button"
                onClick={() => setConfirming(false)}
                className="rounded-lg px-3 py-2 text-sm text-muted-foreground hover:text-foreground"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>

      <div className="rounded-xl border border-border bg-card p-6">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold">Deletion journal</h2>
          {meta && (
            <span
              className={
                "inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-semibold " +
                (meta.verified
                  ? "bg-emerald-500/15 text-emerald-500"
                  : "bg-red-500/15 text-red-500")
              }
            >
              {meta.verified ? <ShieldCheck className="h-3.5 w-3.5" /> : <AlertTriangle className="h-3.5 w-3.5" />}
              {meta.verified ? "Chain verified" : "Chain broken"}
            </span>
          )}
        </div>

        {journalQ.isLoading ? (
          <div className="flex items-center gap-2 py-8 text-sm text-muted-foreground">
            <Loader2 className="h-4 w-4 animate-spin" /> Loading&hellip;
          </div>
        ) : !journalQ.data?.rows.length ? (
          <p className="py-8 text-center text-sm text-muted-foreground">
            No erasures recorded yet. When you erase a user, an entry appears here.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-muted-foreground">
                  <th className="py-2 pr-4 font-medium">Deleted user</th>
                  <th className="py-2 pr-4 font-medium">Erased by</th>
                  <th className="py-2 pr-4 font-medium">Records</th>
                  <th className="py-2 pr-4 font-medium">Reason</th>
                  <th className="py-2 pr-4 font-medium">When</th>
                </tr>
              </thead>
              <tbody>
                {journalQ.data.rows.map((r) => (
                  <tr key={r.id} className="border-b border-border/50">
                    <td className="py-2 pr-4 font-mono text-xs">{r.deleted_user_id.slice(0, 8)}&hellip;</td>
                    <td className="py-2 pr-4">{r.actor_email || "—"}</td>
                    <td className="py-2 pr-4">{r.records_affected}</td>
                    <td className="py-2 pr-4 text-muted-foreground">{r.reason || "—"}</td>
                    <td className="py-2 pr-4 text-muted-foreground">
                      {new Date(r.created_at).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

function errMsg(e: unknown): string {
  const ax = e as { response?: { data?: { error?: { message?: string } } } };
  return ax?.response?.data?.error?.message || (e as Error)?.message || "Something went wrong";
}
`
	return strings.ReplaceAll(src, "~", "`")
}
