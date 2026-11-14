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

	// Not omitempty. An absent list and an empty list are the same thing to
	// Go and very different things to a client: one is [] and the other is
	// undefined, and undefined.filter is how this page used to crash on the
	// first review anybody opened.
	Items []AccessReviewItem ~gorm:"foreignKey:ReviewID" json:"items"~
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
import { ShieldCheck, Check, X, Plus, Loader2, UserCheck, Users, KeyRound } from "@/lib/icons";
import { Button } from "@/components/ui/button";
import { Pill, type PillTone } from "@/components/ui/pill";
import { Note } from "@/components/ui/note";
import { EmptyState } from "@/components/ui/empty-state";
import { inputClasses } from "@/components/ui/input";

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

// The decision, as a word somebody would say.
//
// The API's enum is approved / revoked / pending, and "approved" on a screen
// about access reads as "this request was granted" rather than "this access
// was allowed to stay", which is the opposite way round from what happened.
const DECISION: Record<string, { label: string; tone: PillTone }> = {
  approved: { label: "Kept", tone: "good" },
  revoked: { label: "Revoked", tone: "bad" },
  pending: { label: "Undecided", tone: "warn" },
};

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
  // (detail.items ?? []) rather than detail.items: the API sends [] now, and
  // a page that dies when it does not is a page that dies on a deploy skew.
  const items = detail?.items ?? [];
  const pendingCount = items.filter((i) => i.decision === "pending").length;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Access reviews"
        subtitle="Certify who has access to what. Revoking a grant removes it immediately and is logged to the audit trail."
        actions={
          <Button onClick={() => setNewOpen(true)} loading={openM.isPending}>
            {!openM.isPending && <Plus className="h-3.5 w-3.5" />}
            Start a review
          </Button>
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
            <Button variant="outline" onClick={() => setNewOpen(false)}>
              Cancel
            </Button>
            <Button onClick={startReview} disabled={!newName.trim()} loading={openM.isPending}>
              Create review
            </Button>
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
              className={inputClasses()}
            />
            <p className="text-xs text-text-muted">
              Name it for the period you&rsquo;re certifying. It becomes the audit record.
            </p>
          </div>

          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">Note</label>
            <textarea
              value={newNote}
              onChange={(e) => setNewNote(e.target.value)}
              rows={3}
              placeholder="Optional context: scope, who asked for it, a ticket reference"
              className={inputClasses({ multiline: true }) + " resize-y"}
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
            <div className="flex items-center gap-2 rounded-xl border border-border bg-bg-elevated p-4 text-sm text-text-secondary">
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" /> Loading…
            </div>
          ) : (listQ.data ?? []).length === 0 ? (
            <div className="rounded-xl border border-border bg-bg-elevated">
              <EmptyState
                icon={<ShieldCheck className="h-5 w-5" />}
                title="No reviews yet"
                description="A review is a point in time you certified who had what. The panel beside this says what one would cover."
                action={
                  <Button onClick={() => setNewOpen(true)}>
                    <Plus className="h-3.5 w-3.5" />
                    Start a review
                  </Button>
                }
              />
            </div>
          ) : (
            (listQ.data ?? []).map((r) => (
              <button
                key={r.id}
                type="button"
                onClick={() => setSelected(r.id)}
                className={
                  "w-full rounded-xl border p-4 text-left transition-grit " +
                  (selected === r.id
                    ? "border-accent-line bg-accent-soft"
                    : "border-border bg-bg-elevated hover:border-border-strong")
                }
              >
                <div className="flex items-center justify-between gap-2">
                  <span className="font-medium text-foreground">{r.name}</span>
                  <Pill tone={r.status === "completed" ? "good" : "warn"} dot>
                    {r.status === "completed" ? "Signed off" : "In progress"}
                  </Pill>
                </div>
                {/* The three numbers, each with its own word.
                    "4 pending · 0 approved · 0 revoked" was one grey line in
                    which the only number that needs acting on looked exactly
                    like the two that do not. */}
                <p className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs">
                  <span className={r.pending_items > 0 ? "font-medium text-warning" : "text-text-muted"}>
                    {r.pending_items} to decide
                  </span>
                  <span className="text-text-muted">{r.approved_items} kept</span>
                  <span className="text-text-muted">{r.revoked_items} revoked</span>
                </p>
              </button>
            ))
          )}
        </div>

        {/* Detail */}
        <div>
          {!detail ? (
            /* What a review would cover, rather than an instruction to click.
               On a project that has never run one, the only thing here was
               "Select a review", which is a dead end on the one screen where
               somebody is deciding whether this feature is worth their
               afternoon. No numbers are quoted: there is no endpoint that
               counts these without opening a campaign, and a made-up figure
               on a compliance screen is worse than no figure. */
            <div className="flex h-full flex-col justify-center rounded-xl border border-dashed border-border bg-bg-elevated/40 p-8">
              <h2 className="text-[15px] font-semibold text-foreground">
                {(listQ.data ?? []).length === 0
                  ? "What a review would cover"
                  : "Pick a review to certify it"}
              </h2>
              <p className="mt-1 max-w-prose text-[13px] text-text-secondary">
                Opening one takes a snapshot of access as it stands right now and turns every line
                of it into something you have to decide on, one at a time.
              </p>

              <ul className="mt-5 space-y-3">
                {[
                  {
                    icon: <Users className="h-4 w-4" />,
                    title: "Every role somebody holds",
                    body: "One line per person per role. Keep it or revoke it, and a revoke takes the role away immediately rather than filing a note about it.",
                  },
                  {
                    icon: <UserCheck className="h-4 w-4" />,
                    title: "Who decided, and when",
                    body: "Each decision is signed with your name and written to the audit trail, which is the part a reviewer asks to see.",
                  },
                  {
                    icon: <KeyRound className="h-4 w-4" />,
                    title: "What it does not cover",
                    body: "API keys are not roles and are not in here. They live on their own page, with their own expiry and their own revoke.",
                  },
                ].map((row) => (
                  <li key={row.title} className="flex gap-3">
                    <span
                      className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-bg-hover text-text-secondary"
                      aria-hidden="true"
                    >
                      {row.icon}
                    </span>
                    <span className="min-w-0">
                      <span className="block text-sm font-medium text-foreground">{row.title}</span>
                      <span className="mt-0.5 block text-[13px] text-text-secondary">{row.body}</span>
                    </span>
                  </li>
                ))}
              </ul>

              <Note className="mt-5">
                A review is a record of a decision, not a schedule. Nothing opens one for you, which
                is deliberate: an unread campaign that opens itself every quarter is worse evidence
                than no campaign at all.
              </Note>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-border bg-bg-elevated p-4">
                <div>
                  <p className="font-semibold text-foreground">{detail.name}</p>
                  <p className="mt-0.5 text-[13px] text-text-secondary">
                    Opened by {detail.created_by_email}
                    {detail.completed_at
                      ? " · signed off"
                      : pendingCount === 0
                        ? " · every line decided"
                        : " · " + pendingCount + " left to decide"}
                  </p>
                </div>
                {detail.status !== "completed" && (
                  <Button
                    onClick={() => completeM.mutate()}
                    disabled={pendingCount > 0}
                    loading={completeM.isPending}
                    title={
                      pendingCount > 0
                        ? "Decide every line before signing off"
                        : "Sign off this review"
                    }
                  >
                    {!completeM.isPending && <UserCheck className="h-3.5 w-3.5" />}
                    Sign off
                  </Button>
                )}
              </div>

              <div className="overflow-hidden rounded-xl border border-border">
                <table className="w-full text-sm">
                  <thead className="text-left text-xs font-medium uppercase tracking-wider text-text-muted">
                    <tr className="border-b border-border">
                      <th className="px-5 py-3 font-medium">Person</th>
                      <th className="px-4 py-3 font-medium">Role</th>
                      <th className="px-4 py-3 font-medium">Decision</th>
                      <th className="px-5 py-3" />
                    </tr>
                  </thead>
                  <tbody>
                    {items.map((it) => (
                      <tr key={it.id} className="border-b border-border/60 last:border-0">
                        <td className="px-5 py-3 text-foreground">{it.user_email}</td>
                        <td className="px-4 py-3 text-text-secondary">{it.role_name}</td>
                        <td className="px-4 py-3">
                          <Pill tone={DECISION[it.decision]?.tone ?? "neutral"} dot>
                            {DECISION[it.decision]?.label ?? it.decision}
                          </Pill>
                        </td>
                        <td className="px-5 py-3 text-right">
                          {detail.status !== "completed" && it.decision === "pending" && (
                            <div className="inline-flex gap-2">
                              <button
                                type="button"
                                onClick={() => decideM.mutate({ itemId: it.id, decision: "approved" })}
                                disabled={decideM.isPending}
                                className="inline-flex items-center gap-1.5 rounded-md border border-border px-2.5 py-1 text-[13px] text-text-secondary transition-grit hover:border-success-line hover:bg-success-soft hover:text-success disabled:opacity-40"
                              >
                                <Check className="h-3.5 w-3.5" aria-hidden="true" /> Keep
                              </button>
                              <button
                                type="button"
                                onClick={() => decideM.mutate({ itemId: it.id, decision: "revoked" })}
                                disabled={decideM.isPending}
                                className="inline-flex items-center gap-1.5 rounded-md border border-border px-2.5 py-1 text-[13px] text-text-secondary transition-grit hover:border-danger-line hover:bg-danger-soft hover:text-danger disabled:opacity-40"
                              >
                                <X className="h-3.5 w-3.5" aria-hidden="true" /> Revoke
                              </button>
                            </div>
                          )}
                        </td>
                      </tr>
                    ))}
                    {items.length === 0 && (
                      <tr>
                        <td colSpan={4} className="p-0">
                          <EmptyState
                            icon={<Users className="h-5 w-5" />}
                            title="Nothing to certify"
                            description="Nobody held a role when this review was opened, so it is complete by default and still worth keeping as the record that you looked."
                          />
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
