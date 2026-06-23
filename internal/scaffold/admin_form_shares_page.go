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

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { Plus, Copy, Lock, Unlock, Trash2, X, ExternalLink } from "@/lib/icons";
import { toast } from "sonner";

interface FormShare {
  id: string;
  resource_name: string;
  token: string;
  has_password: boolean;
  enabled: boolean;
  submission_count: number;
  label: string;
  created_at: string;
}

// Web app origin for the shareable link. NEXT_PUBLIC_WEB_URL is set in
// .env (defaults to http://localhost:3000) — same convention as the
// API_URL constant used by the Sentinel/Pulse external links.
const WEB_URL = process.env.NEXT_PUBLIC_WEB_URL || "http://localhost:3000";

export default function FormSharesPage() {
  const [createOpen, setCreateOpen] = useState(false);

  const { data, isLoading } = useQuery<{ data: FormShare[] }>({
    queryKey: ["form-shares"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/form-shares");
      return data;
    },
  });

  const shares = data?.data ?? [];

  return (
    <div>
      <PageHeader
        title="Public form sharing"
        subtitle="Generate share links so anyone with the URL can submit a resource's form — with or without a password gate."
        actions={
          <button
            onClick={() => setCreateOpen(true)}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-sm font-semibold text-white hover:bg-accent-hover"
          >
            <Plus className="h-4 w-4" />
            New share
          </button>
        }
      />

      {isLoading ? (
        <SkeletonCards count={3} />
      ) : shares.length === 0 ? (
        <div className="mt-6 rounded-xl border border-dashed border-border bg-bg-elevated p-12 text-center">
          <p className="text-foreground font-medium">No shares yet</p>
          <p className="mt-1 text-sm text-text-secondary">
            Create a share to let visitors submit forms for one of your resources without an admin login.
          </p>
        </div>
      ) : (
        <div className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
          <table className="w-full text-sm">
            <thead className="border-b border-border bg-bg-secondary">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-text-secondary">Resource</th>
                <th className="px-4 py-3 text-left font-medium text-text-secondary">Label</th>
                <th className="px-4 py-3 text-left font-medium text-text-secondary">Protection</th>
                <th className="px-4 py-3 text-right font-medium text-text-secondary">Submissions</th>
                <th className="px-4 py-3 text-left font-medium text-text-secondary">Status</th>
                <th className="px-4 py-3 text-right font-medium text-text-secondary">Actions</th>
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
  const qc = useQueryClient();
  const publicURL = WEB_URL + "/forms/" + share.token;

  const { mutate: update } = useMutation({
    mutationFn: async (body: Record<string, unknown>) => {
      await apiClient.patch("/api/admin/form-shares/" + share.id, body);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["form-shares"] }),
  });

  const { mutate: remove } = useMutation({
    mutationFn: async () => {
      await apiClient.delete("/api/admin/form-shares/" + share.id);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["form-shares"] });
      toast.success("Share deleted");
    },
  });

  const copyLink = async () => {
    try {
      await navigator.clipboard.writeText(publicURL);
      toast.success("Link copied");
    } catch {
      toast.error("Copy failed — your browser blocked clipboard access");
    }
  };

  return (
    <tr className="hover:bg-bg-hover">
      <td className="px-4 py-3 font-mono text-xs text-foreground">{share.resource_name}</td>
      <td className="px-4 py-3 text-foreground">{share.label || <span className="text-text-muted">—</span>}</td>
      <td className="px-4 py-3">
        {share.has_password ? (
          <span className="inline-flex items-center gap-1 rounded-full bg-warning/15 px-2 py-0.5 text-xs font-medium text-warning">
            <Lock className="h-3 w-3" /> Password
          </span>
        ) : (
          <span className="inline-flex items-center gap-1 rounded-full bg-bg-secondary px-2 py-0.5 text-xs font-medium text-text-secondary">
            <Unlock className="h-3 w-3" /> Open
          </span>
        )}
      </td>
      <td className="px-4 py-3 text-right tabular-nums text-foreground">{share.submission_count}</td>
      <td className="px-4 py-3">
        <button
          onClick={() => update({ enabled: !share.enabled })}
          className={
            share.enabled
              ? "inline-flex items-center gap-1 rounded-full bg-success/15 px-2 py-0.5 text-xs font-medium text-success hover:bg-success/25"
              : "inline-flex items-center gap-1 rounded-full bg-bg-secondary px-2 py-0.5 text-xs font-medium text-text-muted hover:bg-bg-hover"
          }
        >
          {share.enabled ? "Enabled" : "Disabled"}
        </button>
      </td>
      <td className="px-4 py-3">
        <div className="flex items-center justify-end gap-1">
          <button
            onClick={copyLink}
            className="inline-flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-2 py-1 text-xs text-text-secondary hover:bg-bg-hover hover:text-foreground"
            title="Copy public link"
          >
            <Copy className="h-3 w-3" /> Copy
          </button>
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
  );
}

function CreateShareModal({ onClose }: { onClose: () => void }) {
  const qc = useQueryClient();
  const [resourceName, setResourceName] = useState("");
  const [label, setLabel] = useState("");
  const [password, setPassword] = useState("");

  const { mutate: create, isPending } = useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post("/api/admin/form-shares", {
        resource_name: resourceName,
        label,
        password,
      });
      return data;
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["form-shares"] });
      toast.success("Share created");
      onClose();
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || "Failed to create");
    },
  });

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-bg-secondary shadow-2xl">
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-lg font-semibold text-foreground">New form share</h2>
          <button onClick={onClose} className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground">
            <X className="h-5 w-5" />
          </button>
        </div>

        <form
          onSubmit={(e) => { e.preventDefault(); create(); }}
          className="space-y-4 p-6"
        >
          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Resource name</label>
            <input
              type="text"
              value={resourceName}
              onChange={(e) => setResourceName(e.target.value)}
              placeholder="Contact"
              autoFocus
              required
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2 text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/30"
            />
            <p className="text-xs text-text-muted">
              PascalCase, must match a resource your generator has registered (a {"&"}quot;case{"&"}quot; in
              {" "}<code>services/form_share_dispatch.go</code>).
            </p>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Label <span className="text-text-muted">(optional)</span></label>
            <input
              type="text"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="Q3 lead form"
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2 text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/30"
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-foreground">Password <span className="text-text-muted">(optional)</span></label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Leave blank for open access"
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2 text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/30"
            />
            <p className="text-xs text-text-muted">
              Stored as bcrypt. Visitors must enter this before the form is shown.
            </p>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-border bg-bg-elevated px-4 py-2 text-sm font-medium text-foreground hover:bg-bg-hover"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isPending || !resourceName}
              className="rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
            >
              {isPending ? "Creating…" : "Create share"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
`
}
