package scaffold

import "strings"

// Access review (access recertification).
//
// SOC 2 CC6.2/CC6.3 and ISO 27001 A.9.2.5 all require the same thing: someone
// with authority periodically looks at who has access to what and certifies
// each grant as still-appropriate or revokes it — and can produce evidence that
// the review happened, who did it, and what changed.
//
// A campaign snapshots every current role assignment into a list of items. A
// reviewer decides each one: approve (keep) or revoke (remove the grant now).
// The snapshot copies the user's email and the role's name at open time, so the
// record stays readable — and auditable — even after the user or role is later
// deleted. Revocations are written to the tamper-evident activity log, so the
// review is cross-checkable against the same trail everything else lands in.

func apiAccessReviewModelGo() string {
	src := `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccessReview is one recertification campaign: a point-in-time snapshot of all
// role assignments, reviewed grant by grant.
type AccessReview struct {
	ID   string ~gorm:"primarykey;size:36" json:"id"~
	Name string ~gorm:"size:200;not null" json:"name" binding:"required"~
	Note string ~gorm:"size:1000" json:"note"~

	// Status is "open" while items are being decided, "completed" once the
	// reviewer signs it off. A completed review is the artifact an auditor asks
	// for; it is never reopened — a new period is a new campaign.
	Status string ~gorm:"size:16;index;default:open" json:"status"~

	CreatedBy      string ~gorm:"size:36;index" json:"created_by"~
	CreatedByEmail string ~gorm:"size:255" json:"created_by_email"~

	CompletedBy      string     ~gorm:"size:36" json:"completed_by,omitempty"~
	CompletedByEmail string     ~gorm:"size:255" json:"completed_by_email,omitempty"~
	CompletedAt      *time.Time ~json:"completed_at,omitempty"~

	CreatedAt time.Time ~json:"created_at"~
	UpdatedAt time.Time ~json:"updated_at"~

	Items []AccessReviewItem ~gorm:"foreignKey:ReviewID" json:"items,omitempty"~
}

func (r *AccessReview) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.Status == "" {
		r.Status = "open"
	}
	return nil
}

// AccessReviewItem is one role assignment under review. UserEmail and RoleName
// are snapshots taken when the campaign opens — the whole point of a review is
// that its record survives later changes to the things it reviewed.
type AccessReviewItem struct {
	ID       string ~gorm:"primarykey;size:36" json:"id"~
	ReviewID string ~gorm:"size:36;index;not null" json:"review_id"~

	UserID    string ~gorm:"size:36;index" json:"user_id"~
	UserEmail string ~gorm:"size:255" json:"user_email"~
	RoleID    string ~gorm:"size:36;index" json:"role_id"~
	RoleName  string ~gorm:"size:80" json:"role_name"~

	// Decision is "pending", "approved", or "revoked". A revoked item has already
	// had its grant removed — the decision and the effect are one transaction.
	Decision        string     ~gorm:"size:16;index;default:pending" json:"decision"~
	DecidedBy       string     ~gorm:"size:36" json:"decided_by,omitempty"~
	DecidedByEmail  string     ~gorm:"size:255" json:"decided_by_email,omitempty"~
	DecidedAt       *time.Time ~json:"decided_at,omitempty"~
	Note            string     ~gorm:"size:1000" json:"note"~

	CreatedAt time.Time ~json:"created_at"~
}

func (i *AccessReviewItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	if i.Decision == "" {
		i.Decision = "pending"
	}
	return nil
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiAccessReviewServiceGo() string {
	src := `package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

var (
	// ErrReviewClosed is returned when a decision is attempted on a review that
	// has already been signed off. A completed review is immutable evidence.
	ErrReviewClosed = errors.New("access review is already completed")
	// ErrReviewIncomplete is returned when completion is attempted while items
	// are still pending. You cannot certify a review you have not finished.
	ErrReviewIncomplete = errors.New("access review still has undecided items")
	// ErrItemDecided guards the one-way door: a revoked grant is already gone, so
	// re-approving it here would record a certification the system can't honour.
	ErrItemDecided = errors.New("this item has already been revoked and cannot be changed")
)

// OpenAccessReview snapshots every current role assignment into a new campaign
// of pending items. The snapshot copies user email and role name so the record
// stays legible after later deletions.
func OpenAccessReview(db *gorm.DB, name, note, createdBy, createdByEmail string) (*models.AccessReview, error) {
	review := &models.AccessReview{
		Name:           name,
		Note:           note,
		Status:         "open",
		CreatedBy:      createdBy,
		CreatedByEmail: createdByEmail,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return err
		}

		// Join UserRole to users and roles so the snapshot carries human-readable
		// identifiers, not just opaque ids.
		type grant struct {
			UserID    string
			UserEmail string
			RoleID    string
			RoleName  string
		}
		var grants []grant
		if err := tx.Table("user_roles").
			Select("user_roles.user_id, users.email AS user_email, user_roles.role_id, roles.name AS role_name").
			Joins("LEFT JOIN users ON users.id = user_roles.user_id").
			Joins("LEFT JOIN roles ON roles.id = user_roles.role_id").
			Scan(&grants).Error; err != nil {
			return err
		}

		items := make([]models.AccessReviewItem, 0, len(grants))
		for _, g := range grants {
			items = append(items, models.AccessReviewItem{
				ReviewID:  review.ID,
				UserID:    g.UserID,
				UserEmail: g.UserEmail,
				RoleID:    g.RoleID,
				RoleName:  g.RoleName,
				Decision:  "pending",
			})
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		review.Items = items
		return nil
	})
	if err != nil {
		return nil, err
	}
	return review, nil
}

// DecideAccessReviewItem records an approve/revoke decision. A revoke removes
// the underlying UserRole in the same transaction and logs it to the activity
// trail, so the decision and its effect can never drift apart.
func DecideAccessReviewItem(db *gorm.DB, c *gin.Context, reviewID, itemID, decision, note, reviewerID, reviewerEmail string) (*models.AccessReviewItem, error) {
	if decision != "approved" && decision != "revoked" {
		return nil, fmt.Errorf("decision must be \"approved\" or \"revoked\", got %q", decision)
	}

	var item models.AccessReviewItem
	err := db.Transaction(func(tx *gorm.DB) error {
		var review models.AccessReview
		if err := tx.First(&review, "id = ?", reviewID).Error; err != nil {
			return err
		}
		if review.Status == "completed" {
			return ErrReviewClosed
		}
		if err := tx.First(&item, "id = ? AND review_id = ?", itemID, reviewID).Error; err != nil {
			return err
		}
		// Revoke is terminal: the grant is already gone, so the record must not be
		// walked back to a certification the system can't stand behind.
		if item.Decision == "revoked" {
			return ErrItemDecided
		}

		if decision == "revoked" {
			// Remove the actual grant. Scoped to this user+role so nothing else is
			// touched; a no-op if it was already removed out of band.
			if err := tx.Where("user_id = ? AND role_id = ?", item.UserID, item.RoleID).
				Delete(&models.UserRole{}).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		item.Decision = decision
		item.Note = note
		item.DecidedBy = reviewerID
		item.DecidedByEmail = reviewerEmail
		item.DecidedAt = &now
		return tx.Save(&item).Error
	})
	if err != nil {
		return nil, err
	}

	// Log the revocation to the semantic activity trail (best-effort, outside the
	// transaction — losing the log line must not roll back a real revocation).
	if decision == "revoked" && c != nil {
		LogActivity(db, c, ActivityArgs{
			UserID:       reviewerID,
			Action:       "access_review.revoke",
			Severity:     "warn",
			Summary:      fmt.Sprintf("Revoked role %q from %s during access review", item.RoleName, item.UserEmail),
			ResourceType: "user",
			ResourceID:   item.UserID,
			Metadata: map[string]interface{}{
				"review_id": reviewID,
				"role_id":   item.RoleID,
			},
		})
	}
	return &item, nil
}

// CompleteAccessReview signs off a campaign. It refuses while any item is still
// pending: a review with undecided grants is not a review.
func CompleteAccessReview(db *gorm.DB, reviewID, reviewerID, reviewerEmail string) (*models.AccessReview, error) {
	var review models.AccessReview
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&review, "id = ?", reviewID).Error; err != nil {
			return err
		}
		if review.Status == "completed" {
			return ErrReviewClosed
		}
		var pending int64
		if err := tx.Model(&models.AccessReviewItem{}).
			Where("review_id = ? AND decision = ?", reviewID, "pending").
			Count(&pending).Error; err != nil {
			return err
		}
		if pending > 0 {
			return ErrReviewIncomplete
		}
		now := time.Now()
		review.Status = "completed"
		review.CompletedBy = reviewerID
		review.CompletedByEmail = reviewerEmail
		review.CompletedAt = &now
		return tx.Save(&review).Error
	})
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// AccessReviewSummary is the per-campaign counts the list view shows.
type AccessReviewSummary struct {
	models.AccessReview
	TotalItems    int64 ~json:"total_items"~
	PendingItems  int64 ~json:"pending_items"~
	ApprovedItems int64 ~json:"approved_items"~
	RevokedItems  int64 ~json:"revoked_items"~
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// adminAccessReviewPage emits the admin UI: a campaign list plus a detail panel
// where each grant is approved or revoked. Written in the Next shape; the
// TanStack build gets it via nextToTanStack at wiring time.
func adminAccessReviewPage() string {
	src := `"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { ShieldCheck, Check, X, Plus, Loader2, UserCheck } from "@/lib/icons";

interface ReviewSummary {
  id: string;
  name: string;
  status: string;
  created_by_email: string;
  created_at: string;
  completed_at?: string;
  total_items: number;
  pending_items: number;
  approved_items: number;
  revoked_items: number;
}

interface ReviewItem {
  id: string;
  user_email: string;
  role_name: string;
  decision: string;
  decided_by_email?: string;
  note?: string;
}

interface ReviewDetail extends ReviewSummary {
  items: ReviewItem[];
}

export default function AccessReviewPage() {
  const qc = useQueryClient();
  const [selected, setSelected] = useState<string | null>(null);

  const listQ = useQuery({
    queryKey: ["access-reviews"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/access-reviews");
      return (data.data ?? []) as ReviewSummary[];
    },
  });

  const detailQ = useQuery({
    queryKey: ["access-review", selected],
    enabled: !!selected,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/access-reviews/" + selected);
      return data.data as ReviewDetail;
    },
  });

  const openM = useMutation({
    mutationFn: async (name: string) => {
      const { data } = await apiClient.post("/api/access-reviews", { name });
      return data.data as ReviewDetail;
    },
    onSuccess: (r) => {
      qc.invalidateQueries({ queryKey: ["access-reviews"] });
      setSelected(r.id);
    },
  });

  const decideM = useMutation({
    mutationFn: async (args: { itemId: string; decision: "approved" | "revoked" }) => {
      await apiClient.post(
        "/api/access-reviews/" + selected + "/items/" + args.itemId + "/decision",
        { decision: args.decision }
      );
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["access-review", selected] });
      qc.invalidateQueries({ queryKey: ["access-reviews"] });
    },
  });

  const completeM = useMutation({
    mutationFn: async () => {
      await apiClient.post("/api/access-reviews/" + selected + "/complete");
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["access-review", selected] });
      qc.invalidateQueries({ queryKey: ["access-reviews"] });
    },
  });

  const startReview = () => {
    const name = window.prompt("Name this access review (e.g. \"Q3 2026 quarterly\")");
    if (name && name.trim()) openM.mutate(name.trim());
  };

  const detail = detailQ.data;
  // The detail endpoint returns items but not the aggregate counts (those live on
  // the list summary), so derive the pending count from the items we already
  // have. This keeps the header and the Complete button honest even mid-review.
  const pendingCount = detail ? detail.items.filter((i) => i.decision === "pending").length : 0;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Access reviews"
        subtitle="Certify who has access to what. Revoking a grant removes it immediately and is logged to the audit trail."
        actions={
          <button
            type="button"
            onClick={startReview}
            disabled={openM.isPending}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
          >
            {openM.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
            New review
          </button>
        }
      />

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-[320px_1fr]">
        {/* Campaign list */}
        <div className="space-y-2">
          {listQ.isLoading ? (
            <div className="flex items-center gap-2 p-4 text-sm text-text-secondary">
              <Loader2 className="h-4 w-4 animate-spin" /> Loading…
            </div>
          ) : (listQ.data ?? []).length === 0 ? (
            <p className="rounded-xl border border-border bg-bg-secondary/40 p-6 text-sm text-text-secondary">
              No access reviews yet. Start one to snapshot every current role assignment for certification.
            </p>
          ) : (
            (listQ.data ?? []).map((r) => (
              <button
                key={r.id}
                type="button"
                onClick={() => setSelected(r.id)}
                className={
                  "w-full rounded-xl border p-4 text-left transition-colors " +
                  (selected === r.id ? "border-accent bg-accent/5" : "border-border bg-bg-secondary/40 hover:border-accent/40")
                }
              >
                <div className="flex items-center justify-between gap-2">
                  <span className="font-medium text-foreground">{r.name}</span>
                  <span
                    className={
                      "rounded-full px-2 py-0.5 text-[11px] font-semibold " +
                      (r.status === "completed" ? "bg-success/10 text-success" : "bg-warning/10 text-warning")
                    }
                  >
                    {r.status}
                  </span>
                </div>
                <p className="mt-1 text-xs text-text-secondary">
                  {r.pending_items} pending · {r.approved_items} approved · {r.revoked_items} revoked
                </p>
              </button>
            ))
          )}
        </div>

        {/* Detail */}
        <div>
          {!detail ? (
            <div className="flex h-full items-center justify-center rounded-xl border border-dashed border-border p-12 text-sm text-text-secondary">
              <ShieldCheck className="mr-2 h-5 w-5" /> Select a review to certify its access.
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-border bg-bg-secondary/40 p-4">
                <div>
                  <p className="font-semibold text-foreground">{detail.name}</p>
                  <p className="text-xs text-text-secondary">
                    Opened by {detail.created_by_email}
                    {detail.completed_at ? " · completed" : " · " + pendingCount + " pending"}
                  </p>
                </div>
                {detail.status !== "completed" && (
                  <button
                    type="button"
                    onClick={() => completeM.mutate()}
                    disabled={pendingCount > 0 || completeM.isPending}
                    title={pendingCount > 0 ? "Decide every item before completing" : "Sign off this review"}
                    className="inline-flex items-center gap-2 rounded-lg border border-success/40 bg-success/10 px-4 py-2 text-sm font-semibold text-success hover:bg-success/20 disabled:opacity-40"
                  >
                    <UserCheck className="h-4 w-4" />
                    Complete review
                  </button>
                )}
              </div>

              <div className="overflow-hidden rounded-xl border border-border">
                <table className="w-full text-sm">
                  <thead className="bg-bg-tertiary text-left text-xs text-text-secondary">
                    <tr>
                      <th className="px-4 py-2 font-medium">User</th>
                      <th className="px-4 py-2 font-medium">Role</th>
                      <th className="px-4 py-2 font-medium">Decision</th>
                      <th className="px-4 py-2" />
                    </tr>
                  </thead>
                  <tbody>
                    {detail.items.map((it) => (
                      <tr key={it.id} className="border-t border-border">
                        <td className="px-4 py-2.5 text-foreground">{it.user_email}</td>
                        <td className="px-4 py-2.5 text-text-secondary">{it.role_name}</td>
                        <td className="px-4 py-2.5">
                          <span
                            className={
                              "rounded-full px-2 py-0.5 text-[11px] font-semibold " +
                              (it.decision === "approved"
                                ? "bg-success/10 text-success"
                                : it.decision === "revoked"
                                ? "bg-danger/10 text-danger"
                                : "bg-bg-hover text-text-secondary")
                            }
                          >
                            {it.decision}
                          </span>
                        </td>
                        <td className="px-4 py-2.5 text-right">
                          {detail.status !== "completed" && it.decision === "pending" && (
                            <div className="inline-flex gap-1">
                              <button
                                type="button"
                                onClick={() => decideM.mutate({ itemId: it.id, decision: "approved" })}
                                disabled={decideM.isPending}
                                className="inline-flex items-center gap-1 rounded-md border border-border px-2 py-1 text-xs text-text-secondary hover:border-success/40 hover:text-success disabled:opacity-40"
                              >
                                <Check className="h-3 w-3" /> Keep
                              </button>
                              <button
                                type="button"
                                onClick={() => decideM.mutate({ itemId: it.id, decision: "revoked" })}
                                disabled={decideM.isPending}
                                className="inline-flex items-center gap-1 rounded-md border border-border px-2 py-1 text-xs text-text-secondary hover:border-danger/40 hover:text-danger disabled:opacity-40"
                              >
                                <X className="h-3 w-3" /> Revoke
                              </button>
                            </div>
                          )}
                        </td>
                      </tr>
                    ))}
                    {detail.items.length === 0 && (
                      <tr>
                        <td colSpan={4} className="px-4 py-6 text-center text-text-secondary">
                          No role assignments existed when this review was opened.
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiAccessReviewTestGo() string {
	src := `package services_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func newReviewDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(tb, err)
	require.NoError(tb, db.AutoMigrate(
		&models.User{}, &models.Role{}, &models.UserRole{},
		&models.AccessReview{}, &models.AccessReviewItem{}, &models.UserActivity{},
	))
	return db
}

// seedGrant creates a user, a role, and the assignment between them.
func seedGrant(tb testing.TB, db *gorm.DB, email, roleName string) (userID, roleID string) {
	tb.Helper()
	u := models.User{Email: email, FirstName: "T", LastName: "U", Password: "x"}
	require.NoError(tb, db.Create(&u).Error)
	r := models.Role{Name: roleName}
	require.NoError(tb, db.Create(&r).Error)
	require.NoError(tb, db.Create(&models.UserRole{UserID: u.ID, RoleID: r.ID}).Error)
	return u.ID, r.ID
}

func TestOpenAccessReviewSnapshotsGrants(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	seedGrant(t, db, "bob@x.dev", "VIEWER")

	review, err := services.OpenAccessReview(db, "Q3 2026", "quarterly", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	require.Len(t, review.Items, 2, "one item per current grant")

	for _, it := range review.Items {
		assert.Equal(t, "pending", it.Decision)
		assert.NotEmpty(t, it.UserEmail, "snapshot must copy the email")
		assert.NotEmpty(t, it.RoleName, "snapshot must copy the role name")
	}
}

// The record must survive later deletion of the user and role it reviewed —
// that is the entire reason for snapshotting.
func TestReviewSnapshotSurvivesDeletion(t *testing.T) {
	db := newReviewDB(t)
	userID, roleID := seedGrant(t, db, "gone@x.dev", "TEMP")

	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	require.NoError(t, db.Delete(&models.User{}, "id = ?", userID).Error)
	require.NoError(t, db.Delete(&models.Role{}, "id = ?", roleID).Error)

	var item models.AccessReviewItem
	require.NoError(t, db.First(&item, "review_id = ?", review.ID).Error)
	assert.Equal(t, "gone@x.dev", item.UserEmail)
	assert.Equal(t, "TEMP", item.RoleName)
}

func TestRevokeRemovesTheGrant(t *testing.T) {
	db := newReviewDB(t)
	userID, roleID := seedGrant(t, db, "ada@x.dev", "EDITOR")

	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	item := review.Items[0]

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, item.ID, "revoked", "no longer needed", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	var count int64
	db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	assert.Equal(t, int64(0), count, "revoke must delete the UserRole")
}

func TestApproveKeepsTheGrant(t *testing.T) {
	db := newReviewDB(t)
	userID, roleID := seedGrant(t, db, "ada@x.dev", "EDITOR")

	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "approved", "still valid", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	var count int64
	db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	assert.Equal(t, int64(1), count, "approve must keep the grant")
}

// A revoked grant is gone; the decision cannot be walked back to a certification
// the system can't honour.
func TestRevokeIsTerminal(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	item := review.Items[0]

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, item.ID, "revoked", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, item.ID, "approved", "oops", "admin-1", "admin@x.dev")
	assert.ErrorIs(t, err, services.ErrItemDecided)
}

func TestCannotCompleteWithPendingItems(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	seedGrant(t, db, "bob@x.dev", "VIEWER")
	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	// Decide only one of the two.
	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "approved", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.CompleteAccessReview(db, review.ID, "admin-1", "admin@x.dev")
	assert.ErrorIs(t, err, services.ErrReviewIncomplete)

	// Decide the second, then completion succeeds.
	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[1].ID, "revoked", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	done, err := services.CompleteAccessReview(db, review.ID, "admin-1", "admin@x.dev")
	require.NoError(t, err)
	assert.Equal(t, "completed", done.Status)
	require.NotNil(t, done.CompletedAt)
}

// A completed review is immutable evidence.
func TestCannotDecideOnCompletedReview(t *testing.T) {
	db := newReviewDB(t)
	seedGrant(t, db, "ada@x.dev", "EDITOR")
	review, err := services.OpenAccessReview(db, "R", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "approved", "", "admin-1", "admin@x.dev")
	require.NoError(t, err)
	_, err = services.CompleteAccessReview(db, review.ID, "admin-1", "admin@x.dev")
	require.NoError(t, err)

	_, err = services.DecideAccessReviewItem(db, nil, review.ID, review.Items[0].ID, "revoked", "", "admin-1", "admin@x.dev")
	assert.ErrorIs(t, err, services.ErrReviewClosed)
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiAccessReviewHandlerGo() string {
	src := `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// AccessReviewHandler exposes the recertification workflow. Admin-only: the
// point of an access review is that a privileged reviewer certifies access, and
// its evidence must not be editable by the people whose access it covers.
type AccessReviewHandler struct {
	DB *gorm.DB
}

func NewAccessReviewHandler(db *gorm.DB) *AccessReviewHandler {
	return &AccessReviewHandler{DB: db}
}

// List returns campaigns newest first, each with its decision counts.
func (h *AccessReviewHandler) List(c *gin.Context) {
	var reviews []models.AccessReview
	if err := h.DB.Order("created_at desc").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "failed to load reviews"},
		})
		return
	}

	out := make([]services.AccessReviewSummary, 0, len(reviews))
	for _, r := range reviews {
		s := services.AccessReviewSummary{AccessReview: r}
		h.DB.Model(&models.AccessReviewItem{}).Where("review_id = ?", r.ID).Count(&s.TotalItems)
		h.DB.Model(&models.AccessReviewItem{}).Where("review_id = ? AND decision = ?", r.ID, "pending").Count(&s.PendingItems)
		h.DB.Model(&models.AccessReviewItem{}).Where("review_id = ? AND decision = ?", r.ID, "approved").Count(&s.ApprovedItems)
		h.DB.Model(&models.AccessReviewItem{}).Where("review_id = ? AND decision = ?", r.ID, "revoked").Count(&s.RevokedItems)
		out = append(out, s)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Get returns one campaign with all its items.
func (h *AccessReviewHandler) Get(c *gin.Context) {
	var review models.AccessReview
	if err := h.DB.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("user_email asc, role_name asc")
	}).First(&review, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "access review not found"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": review})
}

type openReviewRequest struct {
	Name string ~json:"name" binding:"required"~
	Note string ~json:"note"~
}

// Open starts a campaign, snapshotting all current role assignments.
func (h *AccessReviewHandler) Open(c *gin.Context) {
	var req openReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}
	review, err := services.OpenAccessReview(h.DB, req.Name, req.Note, c.GetString("user_id"), c.GetString("user_email"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "failed to open review"},
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data":    review,
		"message": "Access review opened",
	})
}

type decideRequest struct {
	Decision string ~json:"decision" binding:"required"~
	Note     string ~json:"note"~
}

// Decide records an approve/revoke on one item. A revoke removes the grant.
func (h *AccessReviewHandler) Decide(c *gin.Context) {
	var req decideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}
	item, err := services.DecideAccessReviewItem(
		h.DB, c, c.Param("id"), c.Param("itemId"),
		req.Decision, req.Note,
		c.GetString("user_id"), c.GetString("user_email"),
	)
	if err != nil {
		status, code := http.StatusBadRequest, "INVALID_DECISION"
		switch err {
		case services.ErrReviewClosed:
			code = "REVIEW_CLOSED"
		case services.ErrItemDecided:
			code = "ITEM_LOCKED"
		case gorm.ErrRecordNotFound:
			status, code = http.StatusNotFound, "NOT_FOUND"
		}
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item, "message": "Decision recorded"})
}

// Complete signs off a campaign once every item is decided.
func (h *AccessReviewHandler) Complete(c *gin.Context) {
	review, err := services.CompleteAccessReview(h.DB, c.Param("id"), c.GetString("user_id"), c.GetString("user_email"))
	if err != nil {
		status, code := http.StatusBadRequest, "CANNOT_COMPLETE"
		switch err {
		case services.ErrReviewIncomplete:
			code = "REVIEW_INCOMPLETE"
		case services.ErrReviewClosed:
			code = "REVIEW_CLOSED"
		case gorm.ErrRecordNotFound:
			status, code = http.StatusNotFound, "NOT_FOUND"
		}
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": review, "message": "Access review completed"})
}
`
	return strings.ReplaceAll(src, "~", "`")
}
