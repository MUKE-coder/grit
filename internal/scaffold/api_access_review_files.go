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

	"{{MODULE}}/internal/ids"
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
		r.ID = ids.New()
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
		i.ID = ids.New()
	}
	if i.Decision == "" {
		i.Decision = "pending"
	}
	return nil
}
`
	return strings.ReplaceAll(src, "~", "`")
}

func apiAccessReviewServiceGo() string { return tmpl("api/services/access_review.go") }

// adminAccessReviewPage emits the admin UI: a campaign list plus a detail panel
// where each grant is approved or revoked. Written in the Next shape; the
// TanStack build gets it via nextToTanStack at wiring time.
func adminAccessReviewPage() string {
	src := `"use client";

import { useState } from "react";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import {
  useAccessReviews,
  useAccessReview,
  useOpenAccessReview,
  useDecideAccessReviewItem,
  useCompleteAccessReview,
} from "@/hooks/use-access-reviews";
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
  const [selected, setSelected] = useState<string | null>(null);
  // "New review" opens a proper form (name + optional note) rather than a
  // browser prompt — the note is stored on the campaign as context for whoever
  // audits it later.
  const [newOpen, setNewOpen] = useState(false);
  const [newName, setNewName] = useState("");
  const [newNote, setNewNote] = useState("");

  // The five requests are in hooks/use-access-reviews.ts, which also knows
  // that a decision moves both the campaign's rows and the list's counts.
  const listQ = useAccessReviews<ReviewSummary>();
  const detailQ = useAccessReview<ReviewDetail>(selected);
  const openM = useOpenAccessReview<ReviewDetail>();
  const decideM = useDecideAccessReviewItem(selected);
  const completeM = useCompleteAccessReview(selected);

  const startReview = () => {
    if (!newName.trim()) return;
    openM.mutate(
      { name: newName.trim(), note: newNote.trim() || undefined },
      {
        onSuccess: (r) => {
          setSelected(r.id);
          setNewOpen(false);
          setNewName("");
          setNewNote("");
        },
      }
    );
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
            onClick={() => setNewOpen(true)}
            disabled={openM.isPending}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
          >
            {openM.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
            New review
          </button>
        }
      />

      {/* New review form. Opening a campaign snapshots every current role
          assignment, so the copy says so up front — it's not a draft. */}
      <ResponsiveSheet
        open={newOpen}
        onClose={() => setNewOpen(false)}
        title="Start an access review"
        description="This snapshots every current role assignment as pending items for you to certify."
        size="md"
        footer={
          <div className="flex items-center justify-end gap-3">
            <button
              type="button"
              onClick={() => setNewOpen(false)}
              className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-text-secondary transition-colors hover:bg-bg-hover"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={startReview}
              disabled={!newName.trim() || openM.isPending}
              className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
            >
              {openM.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
              Create review
            </button>
          </div>
        }
      >
        <div className="space-y-4">
          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">
              Name<span className="ml-1 text-danger">*</span>
            </label>
            <input
              autoFocus
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && newName.trim()) startReview();
              }}
              placeholder="Q3 2026 quarterly"
              className="w-full rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            />
            <p className="text-xs text-text-muted">
              Name it for the period you&rsquo;re certifying — it becomes the audit record.
            </p>
          </div>

          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">Note</label>
            <textarea
              value={newNote}
              onChange={(e) => setNewNote(e.target.value)}
              rows={3}
              placeholder="Optional context — scope, who requested it, ticket reference…"
              className="w-full resize-y rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            />
          </div>

          {openM.isError && (
            <p className="text-sm text-danger">Could not start the review. Please try again.</p>
          )}
        </div>
      </ResponsiveSheet>

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

func apiAccessReviewTestGo() string { return tmpl("api/services/access_review_test.go") }

func apiAccessReviewHandlerGo() string { return tmpl("api/handlers/access_review.go") }
