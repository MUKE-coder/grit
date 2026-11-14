package scaffold

// adminFormSharesPage emits the /system/form-shares admin page.
//
// What operators do here:
//   - See every public share at a glance (resource, status, password?, hits)
//   - Generate a new share for any resource by typing its Pascal name +
//     an optional bcrypt-protected password
//   - Copy the public URL to the clipboard (toast confirms)
//   - Toggle enabled / disable + soft-delete shares
//
// The page is intentionally a single file with inline state — no
// per-row drawer, no detail page. Form sharing is a low-volume
// operator task; one tidy table beats a five-screen flow.
func adminFormSharesPage() string {
	return `"use client";

import { useEffect, useState } from "react";
import {
  useFormShares,
  useFormSubmissions,
  useFormShareResources,
  useFormShareFields,
  useCreateFormShare,
  useUpdateFormShare,
  useDeleteFormShare,
} from "@/hooks/use-form-shares";
import { getApiErrorMessage } from "@/lib/api-core";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { Plus, Lock, Unlock, Trash2, X, ExternalLink, Activity, Pencil, Package, Shield, Link2 } from "@/lib/icons";
import { Button, buttonClasses } from "@/components/ui/button";
import { Pill } from "@/components/ui/pill";
import { CopyButton } from "@/components/ui/copy-button";
import { EmptyState } from "@/components/ui/empty-state";
import { toast } from "sonner";
import { inputClasses } from "@/components/ui/input";

interface FormShare {
  id: string;
  resource_name: string;
  token: string;
  has_password: boolean;
  enabled: boolean;
  submission_count: number;
  label: string;
  // v3.31.50 -- customisable surface.
  custom_title?: string;
  custom_description?: string;
  hidden_fields?: string[];
  created_at: string;
}

interface FormSubmission {
  id: string;
  share_id: string;
  resource_name: string;
  record_id: string;
  ip: string;
  user_agent: string;
  created_at: string;
}

// Web app origin for the shareable link. NEXT_PUBLIC_WEB_URL is set in
// .env (defaults to http://localhost:3000) — same convention as the
// API_URL constant used by the Sentinel/Pulse external links.
const WEB_URL = process.env.NEXT_PUBLIC_WEB_URL || "http://localhost:3000";

export default function FormSharesPage() {
  const [createOpen, setCreateOpen] = useState(false);

  // Every request this screen makes is in hooks/use-form-shares.ts. It was
  // eight apiClient calls and four useMutation blocks spread over four
  // components here, each with its own query key, so a share created in one
  // modal and edited in another invalidated two different things.
  const { data, isLoading } = useFormShares<FormShare>();

  const shares = data?.data ?? [];

  return (
    <div>
      <PageHeader
        title="Public form sharing"
        subtitle="A link anyone can open to submit one of your forms, with or without a password on it."
        actions={
          <Button onClick={() => setCreateOpen(true)}>
            <Plus className="h-3.5 w-3.5" />
            New share link
          </Button>
        }
      />

      {isLoading ? (
        <SkeletonCards count={3} />
      ) : shares.length === 0 ? (
        <>
          <div className="rounded-2xl border border-border bg-bg-elevated">
            <EmptyState
              icon={<Link2 className="h-5 w-5" />}
              title="No share links yet"
              description="A share link lets somebody fill in one of your forms without an account, and the submission lands as a real record."
              action={
                <Button onClick={() => setCreateOpen(true)}>
                  <Plus className="h-3.5 w-3.5" />
                  New share link
                </Button>
              }
            />
          </div>

          {/* What a share link is, in the order you decide it.
              The empty state said what the feature does and stopped. Three
              decisions follow it and all three are made in the same modal, so
              naming them here is the difference between opening that modal
              knowing what it will ask and discovering it one field at a
              time. */}
          <div className="mt-6 grid gap-4 sm:grid-cols-3">
            {[
              {
                icon: <Package className="h-4 w-4" />,
                title: "Pick a resource",
                body: "The form is built from that resource's own fields, so it validates the same way the admin does. Hide the fields a stranger has no business filling in.",
              },
              {
                icon: <Shield className="h-4 w-4" />,
                title: "Protect it",
                body: "Leave it open for a public contact form, or put a password on it when the link is going to one supplier rather than to the world.",
              },
              {
                icon: <Link2 className="h-4 w-4" />,
                title: "Share the link",
                body: "Every submission is recorded with its time and address, and you can turn the link off without deleting what came through it.",
              },
            ].map((step, i) => (
              <article key={step.title} className="rounded-xl border border-border bg-bg-elevated p-4">
                <div className="mb-3 flex items-center gap-2">
                  <span
                    className="inline-flex h-8 w-8 items-center justify-center rounded-lg bg-bg-hover text-text-secondary"
                    aria-hidden="true"
                  >
                    {step.icon}
                  </span>
                  <span className="font-mono text-[11px] text-text-muted">{"0" + (i + 1)}</span>
                </div>
                <h2 className="text-sm font-semibold text-foreground">{step.title}</h2>
                <p className="mt-1 text-[13px] text-text-secondary">{step.body}</p>
              </article>
            ))}
          </div>
        </>
      ) : (
        <div className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
          <table className="w-full text-sm">
            <thead className="border-b border-border">
              <tr className="text-xs font-medium uppercase tracking-wider text-text-muted">
                <th className="px-5 py-3 text-left font-medium">Resource</th>
                <th className="px-4 py-3 text-left font-medium">Label</th>
                <th className="px-4 py-3 text-left font-medium">Protection</th>
                <th className="px-4 py-3 text-right font-medium">Submissions</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-5 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {shares.map((s) => (
                <ShareRow key={s.id} share={s} />
              ))}
            </tbody>
          </table>
        </div>
      )}

      {createOpen && <CreateShareModal onClose={() => setCreateOpen(false)} />}
    </div>
  );
}

function ShareRow({ share }: { share: FormShare }) {
  const publicURL = WEB_URL + "/forms/" + share.token;
  const [submissionsOpen, setSubmissionsOpen] = useState(false);
  // v3.31.50 — Edit modal lets the operator change the title,
  // description, password mode, and hidden fields after creation.
  const [editOpen, setEditOpen] = useState(false);

  const updateM = useUpdateFormShare();
  const deleteM = useDeleteFormShare();
  const update = (body: Record<string, unknown>) => updateM.mutate({ id: share.id, body });
  const remove = () =>
    deleteM.mutate(share.id, { onSuccess: () => toast.success("Share deleted") });

  return (
    <>
    <tr className="hover:bg-bg-hover">
      <td className="px-5 py-3 font-mono text-xs text-foreground">{share.resource_name}</td>
      <td className="px-4 py-3 text-foreground">{share.label || <span className="text-text-muted">—</span>}</td>
      <td className="px-4 py-3">
        {/* Open is the one worth noticing, so it is the one that is not
            neutral: a link anybody can submit through is a decision, and a
            grey chip reads as the default. */}
        {share.has_password ? (
          <Pill tone="neutral">
            <Lock className="h-3 w-3" aria-hidden="true" /> Password
          </Pill>
        ) : (
          <Pill tone="warn">
            <Unlock className="h-3 w-3" aria-hidden="true" /> Open to anyone
          </Pill>
        )}
      </td>
      <td className="px-4 py-3 text-right tabular-nums text-foreground">{share.submission_count}</td>
      <td className="px-4 py-3">
        <button
          onClick={() => update({ enabled: !share.enabled })}
          title={share.enabled ? "Turn this link off" : "Turn this link back on"}
          className="rounded-md transition-grit hover:opacity-80 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          <Pill tone={share.enabled ? "good" : "neutral"} dot>
            {share.enabled ? "Live" : "Off"}
          </Pill>
        </button>
      </td>
      <td className="px-4 py-3">
        <div className="flex items-center justify-end gap-1">
          <button
            onClick={() => setEditOpen(true)}
            className="inline-flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-2 py-1 text-xs text-text-secondary hover:bg-bg-hover hover:text-foreground"
            title="Edit share"
          >
            <Pencil className="h-3 w-3" /> Edit
          </button>
          <button
            onClick={() => setSubmissionsOpen(true)}
            className="inline-flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-2 py-1 text-xs text-text-secondary hover:bg-bg-hover hover:text-foreground"
            title="View submissions"
          >
            <Activity className="h-3 w-3" /> Audit
          </button>
          {/* The tick lands on the button that was pressed. A toast in the
              corner is for something that happened somewhere else. */}
          <CopyButton
            value={publicURL}
            label={"Copy the public link for " + (share.label || share.resource_name)}
            className="border border-border bg-bg-elevated text-xs"
          >
            Copy
          </CopyButton>
          <a
            href={publicURL}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-2 py-1 text-xs text-text-secondary hover:bg-bg-hover hover:text-foreground"
            title="Open public form"
          >
            <ExternalLink className="h-3 w-3" /> Open
          </a>
          <button
            onClick={() => {
              if (confirm("Delete this share? Existing submissions are kept; the link stops working.")) {
                remove();
              }
            }}
            className="inline-flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-2 py-1 text-xs text-danger hover:bg-danger/10"
            title="Delete share"
          >
            <Trash2 className="h-3 w-3" />
          </button>
        </div>
      </td>
    </tr>
    {submissionsOpen && <SubmissionsModal share={share} onClose={() => setSubmissionsOpen(false)} />}
    {editOpen && <EditShareModal share={share} onClose={() => setEditOpen(false)} />}
    </>
  );
}

// SubmissionsModal — drill-in audit log for one FormShare. Shows the
// 100 most recent submissions with their record ID, IP, and User-Agent.
// v3.31.25.
function SubmissionsModal({ share, onClose }: { share: FormShare; onClose: () => void }) {
  const { data, isLoading } = useFormSubmissions<FormSubmission>(share.id);

  const rows = data?.data ?? [];

  return (
    <tr>
      <td colSpan={6}>
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
          <div className="relative z-10 w-full max-w-3xl rounded-2xl border border-border bg-bg-secondary shadow-2xl">
            <div className="flex items-center justify-between border-b border-border px-6 py-4">
              <div>
                <h2 className="text-lg font-semibold text-foreground">Audit log</h2>
                <p className="text-xs text-text-secondary">
                  {share.resource_name} · {share.label || share.token.slice(0, 12) + "…"}
                </p>
              </div>
              <button onClick={onClose} className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground">
                <X className="h-5 w-5" />
              </button>
            </div>

            <div className="max-h-[60vh] overflow-y-auto p-6">
              {isLoading ? (
                <p className="text-sm text-text-secondary">Loading…</p>
              ) : rows.length === 0 ? (
                <div className="rounded-xl border border-dashed border-border bg-bg-elevated p-8 text-center">
                  <p className="text-foreground font-medium">No submissions yet</p>
                  <p className="mt-1 text-sm text-text-secondary">Audit rows appear here after the first public submission.</p>
                </div>
              ) : (
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border">
                      <th className="px-2 py-2 text-left font-medium text-text-secondary">When</th>
                      <th className="px-2 py-2 text-left font-medium text-text-secondary">Record ID</th>
                      <th className="px-2 py-2 text-left font-medium text-text-secondary">IP</th>
                      <th className="px-2 py-2 text-left font-medium text-text-secondary">User Agent</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {rows.map((row) => (
                      <tr key={row.id}>
                        <td className="px-2 py-2 text-xs text-text-secondary whitespace-nowrap">
                          {new Date(row.created_at).toLocaleString()}
                        </td>
                        <td className="px-2 py-2 font-mono text-xs text-foreground">
                          {row.record_id.slice(0, 8)}…
                        </td>
                        <td className="px-2 py-2 font-mono text-xs text-text-secondary">{row.ip || "—"}</td>
                        <td className="px-2 py-2 text-xs text-text-muted truncate max-w-xs" title={row.user_agent}>
                          {row.user_agent || "—"}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          </div>
        </div>
      </td>
    </tr>
  );
}

function CreateShareModal({ onClose }: { onClose: () => void }) {
  const [resourceName, setResourceName] = useState("");
  const [label, setLabel] = useState("");
  const [password, setPassword] = useState("");
  // v3.31.50 -- customisable public surface.
  const [customTitle, setCustomTitle] = useState("");
  const [customDescription, setCustomDescription] = useState("");
  const [hiddenFields, setHiddenFields] = useState<Set<string>>(new Set());

  const { data: resources, isLoading: resourcesLoading } = useFormShareResources();
  const { data: fields, isLoading: fieldsLoading } = useFormShareFields(resourceName);

  // biome-ignore lint/correctness/useExhaustiveDependencies: clears the hidden fields whenever the resource changes.
  useEffect(() => { setHiddenFields(new Set()); }, [resourceName]);

  const toggleHidden = (key: string, hide: boolean) => {
    setHiddenFields((prev) => {
      const next = new Set(prev);
      if (hide) next.add(key); else next.delete(key);
      return next;
    });
  };

  const createM = useCreateFormShare();
  const isPending = createM.isPending;
  const create = () =>
    createM.mutate(
      {
        resource_name: resourceName,
        label,
        password,
        custom_title: customTitle,
        custom_description: customDescription,
        hidden_fields: Array.from(hiddenFields),
      },
      {
        onSuccess: () => {
          toast.success("Share created");
          onClose();
        },
        onError: (err: unknown) => toast.error(getApiErrorMessage(err, "Failed to create")),
      }
    );

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 my-8 max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-2xl border border-border bg-bg-secondary shadow-2xl">
        <div className="sticky top-0 z-10 flex items-center justify-between border-b border-border bg-bg-secondary px-6 py-4">
          <h2 className="text-lg font-semibold text-foreground">New form share</h2>
          <button onClick={onClose} className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground">
            <X className="h-5 w-5" />
          </button>
        </div>

        <form onSubmit={(e) => { e.preventDefault(); create(); }} className="space-y-5 p-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Resource</label>
            <select value={resourceName} onChange={(e) => setResourceName(e.target.value)} required autoFocus
              className={inputClasses()}>
              <option value="">{resourcesLoading ? "Loading…" : "Select a resource…"}</option>
              {(resources ?? []).map((r) => (<option key={r} value={r}>{r}</option>))}
            </select>
            <p className="text-xs text-text-muted">
              Only resources registered in <code>services/form_share_dispatch.go</code> can be shared publicly.
            </p>
          </div>

          {resourceName && (
            <div className="space-y-2 rounded-xl border border-border bg-bg-elevated/40 p-4">
              <div className="flex items-center justify-between">
                <p className="text-sm font-medium text-foreground">Form preview</p>
                <p className="text-xs text-text-muted">
                  {fieldsLoading ? "Loading fields…" : (fields?.length ?? 0) + " field" + ((fields?.length ?? 0) === 1 ? "" : "s")}
                </p>
              </div>
              {fields && fields.length > 0 && (
                <ul className="space-y-2">
                  {fields.map((f) => {
                    const isHidden = hiddenFields.has(f.key);
                    return (
                      <li key={f.key} className={"flex items-center justify-between gap-3 rounded-lg border px-3 py-2 text-sm transition-colors " + (isHidden ? "border-border bg-bg-secondary opacity-60" : "border-border bg-bg-elevated")}>
                        <div className="min-w-0">
                          <p className="truncate text-foreground">{f.label}{f.required && <span className="ml-1 text-danger">*</span>}</p>
                          <p className="text-[11px] text-text-muted">
                            <code className="font-mono">{f.key}</code><span className="mx-1.5">·</span>{f.type}<span className="mx-1.5">·</span>{f.required ? "required" : "optional"}
                          </p>
                        </div>
                        {f.required ? (
                          <span className="text-[11px] text-text-muted">always shown</span>
                        ) : (
                          <label className="flex shrink-0 items-center gap-2 text-xs text-text-secondary">
                            <input type="checkbox" checked={isHidden} onChange={(e) => toggleHidden(f.key, e.target.checked)} className="h-3.5 w-3.5 rounded border-border" />
                            Hide
                          </label>
                        )}
                      </li>
                    );
                  })}
                </ul>
              )}
              {fields && fields.length === 0 && !fieldsLoading && (
                <p className="rounded-lg border border-warning/30 bg-warning/5 px-3 py-2 text-xs text-warning">
                  This resource has no public-form fields registered. Regenerate it to pick up the field schema.
                </p>
              )}
            </div>
          )}

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Form title <span className="text-text-muted">(optional)</span></label>
            <input type="text" value={customTitle} onChange={(e) => setCustomTitle(e.target.value)} placeholder={resourceName ? "New " + resourceName : "Heading shown on the public form"} className={inputClasses()} />
            <p className="text-xs text-text-muted">Shown as the heading on the public form. Blank = falls back to the label, then to the resource name.</p>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Description <span className="text-text-muted">(optional)</span></label>
            <textarea value={customDescription} onChange={(e) => setCustomDescription(e.target.value)} placeholder="One-line subtitle shown under the heading" rows={2} className={inputClasses({ multiline: true })} />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Internal label <span className="text-text-muted">(optional)</span></label>
            <input type="text" value={label} onChange={(e) => setLabel(e.target.value)} placeholder="Q3 lead form" className={inputClasses()} />
            <p className="text-xs text-text-muted">Operator-facing tag (the visitor never sees this).</p>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Password <span className="text-text-muted">(optional)</span></label>
            <input type="password"
 autoComplete="off" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Leave blank for open access" className={inputClasses()} />
            <p className="text-xs text-text-muted">Stored as bcrypt. Visitors must enter this before the form is shown.</p>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-lg border border-border bg-bg-elevated px-4 py-2 text-sm font-medium text-foreground hover:bg-bg-hover">Cancel</button>
            <button type="submit" disabled={isPending || !resourceName} className={buttonClasses()}>
              {isPending ? "Creating…" : "Create share"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// v3.31.50 -- EditShareModal lets the operator change label,
// password mode, custom title + description, and hidden fields
// after the share is created.
type PasswordMode = "keep" | "set" | "remove";

function EditShareModal({ share, onClose }: { share: FormShare; onClose: () => void }) {
  const [label, setLabel] = useState(share.label ?? "");
  const [customTitle, setCustomTitle] = useState(share.custom_title ?? "");
  const [customDescription, setCustomDescription] = useState(share.custom_description ?? "");
  const [hiddenFields, setHiddenFields] = useState<Set<string>>(() => new Set(share.hidden_fields ?? []));
  const [passwordMode, setPasswordMode] = useState<PasswordMode>("keep");
  const [newPassword, setNewPassword] = useState("");

  const { data: fields, isLoading: fieldsLoading } = useFormShareFields(share.resource_name);

  const toggleHidden = (key: string, hide: boolean) => {
    setHiddenFields((prev) => {
      const next = new Set(prev);
      if (hide) next.add(key); else next.delete(key);
      return next;
    });
  };

  const saveM = useUpdateFormShare();
  const isPending = saveM.isPending;
  const save = () => {
    const body: Record<string, unknown> = {
      label,
      custom_title: customTitle,
      custom_description: customDescription,
      hidden_fields: Array.from(hiddenFields),
    };
    if (passwordMode === "set" && newPassword) body.password = newPassword;
    else if (passwordMode === "remove") body.password = "-";
    saveM.mutate(
      { id: share.id, body },
      {
        onSuccess: () => {
          toast.success("Share updated");
          onClose();
        },
        onError: (err: unknown) => toast.error(getApiErrorMessage(err, "Failed to save")),
      }
    );
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 my-8 max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-2xl border border-border bg-bg-secondary shadow-2xl">
        <div className="sticky top-0 z-10 flex items-center justify-between border-b border-border bg-bg-secondary px-6 py-4">
          <div>
            <h2 className="text-lg font-semibold text-foreground">Edit form share</h2>
            <p className="text-xs text-text-muted">
              <code className="font-mono">{share.resource_name}</code>
              <span className="mx-1.5">·</span>{share.token.slice(0, 12)}…
            </p>
          </div>
          <button onClick={onClose} className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground">
            <X className="h-5 w-5" />
          </button>
        </div>

        <form onSubmit={(e) => { e.preventDefault(); save(); }} className="space-y-5 p-6">
          <div className="space-y-2 rounded-xl border border-border bg-bg-elevated/40 p-4">
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-foreground">Form preview</p>
              <p className="text-xs text-text-muted">
                {fieldsLoading ? "Loading fields…" : (fields?.length ?? 0) + " field" + ((fields?.length ?? 0) === 1 ? "" : "s")}
              </p>
            </div>
            {fields && fields.length > 0 && (
              <ul className="space-y-2">
                {fields.map((f) => {
                  const isHidden = hiddenFields.has(f.key);
                  return (
                    <li key={f.key} className={"flex items-center justify-between gap-3 rounded-lg border px-3 py-2 text-sm transition-colors " + (isHidden ? "border-border bg-bg-secondary opacity-60" : "border-border bg-bg-elevated")}>
                      <div className="min-w-0">
                        <p className="truncate text-foreground">{f.label}{f.required && <span className="ml-1 text-danger">*</span>}</p>
                        <p className="text-[11px] text-text-muted">
                          <code className="font-mono">{f.key}</code><span className="mx-1.5">·</span>{f.type}<span className="mx-1.5">·</span>{f.required ? "required" : "optional"}
                        </p>
                      </div>
                      {f.required ? (
                        <span className="text-[11px] text-text-muted">always shown</span>
                      ) : (
                        <label className="flex shrink-0 items-center gap-2 text-xs text-text-secondary">
                          <input type="checkbox" checked={isHidden} onChange={(e) => toggleHidden(f.key, e.target.checked)} className="h-3.5 w-3.5 rounded border-border" />
                          Hide
                        </label>
                      )}
                    </li>
                  );
                })}
              </ul>
            )}
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Form title <span className="text-text-muted">(optional)</span></label>
            <input type="text" value={customTitle} onChange={(e) => setCustomTitle(e.target.value)} placeholder={"New " + share.resource_name} className={inputClasses()} />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Description <span className="text-text-muted">(optional)</span></label>
            <textarea value={customDescription} onChange={(e) => setCustomDescription(e.target.value)} rows={2} className={inputClasses({ multiline: true })} />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Internal label <span className="text-text-muted">(optional)</span></label>
            <input type="text" value={label} onChange={(e) => setLabel(e.target.value)} placeholder="Q3 lead form" className={inputClasses()} />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Password</label>
            <div className="flex items-center gap-1 rounded-md border border-border bg-bg-tertiary p-0.5">
              {(["keep", "set", "remove"] as PasswordMode[]).map((mode) => {
                const disabled = mode === "remove" && !share.has_password;
                return (
                  <button key={mode} type="button" onClick={() => !disabled && setPasswordMode(mode)} disabled={disabled}
                    className={"flex-1 rounded px-2.5 py-1 text-xs font-medium transition-colors " + (passwordMode === mode ? "bg-accent text-accent-fg" : disabled ? "text-text-muted opacity-40 cursor-not-allowed" : "text-text-secondary hover:text-foreground")}>
                    {mode === "keep" ? "Keep current" : mode === "set" ? "Set password" : "Remove password"}
                  </button>
                );
              })}
            </div>
            {passwordMode === "set" && (
              <input type="password"
 autoComplete="off" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} placeholder="New password" className={inputClasses()} />
            )}
            {passwordMode === "remove" && (
              <p className="rounded-lg border border-warning/30 bg-warning/5 px-3 py-2 text-xs text-warning">
                Removing the password makes the form open to anyone with the link.
              </p>
            )}
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-lg border border-border bg-bg-elevated px-4 py-2 text-sm font-medium text-foreground hover:bg-bg-hover">Cancel</button>
            <button type="submit" disabled={isPending} className={buttonClasses()}>
              {isPending ? "Saving…" : "Save changes"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
`
}
