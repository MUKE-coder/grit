package scaffold

// The API-keys screen.
//
// One rule drives the whole design: the raw key exists for exactly one render.
// The server stores a hash, so there is no "show key" affordance to add later
// and no support path that recovers it. The reveal panel therefore has to be
// unmissable and refuses to close until the key has been copied or downloaded
// — the same contract as the 2FA backup codes, for the same reason.
//
// Revoked keys stay in the list rather than disappearing. Revocation is a
// timestamp, not a delete, because "which key did this?" is a question people
// ask about keys that were turned off months ago.

func adminAPIKeysPage() string {
	return `"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { KeyRound, Plus, Copy, Check, Loader2, Trash2, AlertTriangle } from "@/lib/icons";
import { Button } from "@/components/ui/button";

interface APIKey {
  id: string;
  name: string;
  prefix: string;
  scopes: string[] | null;
  last_used_at?: string | null;
  expires_at?: string | null;
  revoked_at?: string | null;
  created_at: string;
}

function relative(iso?: string | null): string {
  if (!iso) return "never";
  const diff = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return mins + "m ago";
  const hours = Math.floor(mins / 60);
  if (hours < 24) return hours + "h ago";
  return Math.floor(hours / 24) + "d ago";
}

/** The one and only render of the raw key. See the file header. */
function RevealPanel({ token, onDone }: { token: string; onDone: () => void }) {
  const [saved, setSaved] = useState(false);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(token);
    } catch {
      // Insecure origin or an unfocused document. Marking it saved anyway —
      // refusing would trap the user in a panel with no way out, and Download
      // is still there.
    }
    setSaved(true);
  };

  const download = () => {
    const blob = new Blob([token + "\n"], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "api-key.txt";
    a.click();
    URL.revokeObjectURL(url);
    setSaved(true);
  };

  return (
    <div className="mb-6 rounded-xl border border-warning/40 bg-warning/[0.06] p-5">
      <div className="mb-2 flex items-center gap-2">
        <AlertTriangle className="h-4 w-4 text-warning" />
        <h2 className="font-semibold text-foreground">Copy your key now</h2>
      </div>
      <p className="mb-3 text-sm text-foreground-secondary">
        This is the only time it will be shown. The server keeps a hash, so it cannot be
        recovered — if you lose it, revoke this key and make another.
      </p>

      <code className="mb-3 block break-all rounded-lg bg-bg-tertiary px-4 py-3 font-mono text-sm text-foreground">
        {token}
      </code>

      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          onClick={copy}
          className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm hover:bg-bg-hover"
        >
          {saved ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
          Copy
        </button>
        <button
          type="button"
          onClick={download}
          className="rounded-lg border border-border px-3 py-1.5 text-sm hover:bg-bg-hover"
        >
          Download
        </button>
        <Button
          size="sm"
          onClick={onDone}
          disabled={!saved}
          title={saved ? undefined : "Copy or download it first"}
          className="ml-auto"
        >
          I have saved it
        </Button>
      </div>
    </div>
  );
}

export default function APIKeysPage() {
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [expiresIn, setExpiresIn] = useState("");
  const [token, setToken] = useState<string | null>(null);
  const [pendingRevoke, setPendingRevoke] = useState<APIKey | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ["api-keys"],
    queryFn: async () => {
      const { data: res } = await apiClient.get("/api/api-keys");
      return (res.data ?? []) as APIKey[];
    },
  });

  const create = useMutation({
    mutationFn: async () => {
      const body: Record<string, unknown> = { name: name.trim() };
      const days = parseInt(expiresIn, 10);
      if (!Number.isNaN(days) && days > 0) body.expires_in_days = days;
      const { data: res } = await apiClient.post("/api/api-keys", body);
      return res.data as { token: string };
    },
    onSuccess: (d) => {
      setToken(d.token);
      setName("");
      setExpiresIn("");
      queryClient.invalidateQueries({ queryKey: ["api-keys"] });
    },
    onError: () => toast.error("Could not create the key."),
  });

  const revoke = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/api-keys/" + id);
    },
    onSuccess: () => {
      toast.success("Key revoked. Requests using it will now be refused.");
      queryClient.invalidateQueries({ queryKey: ["api-keys"] });
    },
    onError: () => toast.error("Could not revoke the key."),
  });

  const keys = data ?? [];

  return (
    <div>
      <PageHeader
        title="API keys"
        subtitle="Long-lived credentials for scripts and integrations, instead of a login."
      />

      <div className="px-6 py-6">
        {token && <RevealPanel token={token} onDone={() => setToken(null)} />}

        {/* ── Create ── */}
        <div className="mb-6 rounded-xl border border-border bg-bg-secondary p-5">
          <h2 className="mb-1 font-semibold text-foreground">New key</h2>
          <p className="mb-4 text-sm text-foreground-secondary">
            A key acts as you — same roles, same permissions. Name it after the thing that
            will use it, so it can be revoked without guesswork.
          </p>
          <div className="flex flex-wrap items-end gap-3">
            <div className="min-w-56 flex-1">
              <label htmlFor="key-name" className="mb-1 block text-xs uppercase tracking-wider text-text-muted">
                Name
              </label>
              <input
                id="key-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="nightly export"
                className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-sm text-foreground"
              />
            </div>
            <div className="w-40">
              <label htmlFor="key-expiry" className="mb-1 block text-xs uppercase tracking-wider text-text-muted">
                Expires in (days)
              </label>
              <input
                id="key-expiry"
                value={expiresIn}
                onChange={(e) => setExpiresIn(e.target.value)}
                placeholder="never"
                inputMode="numeric"
                className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-sm text-foreground"
              />
            </div>
            <Button
              onClick={() => create.mutate()}
              disabled={!name.trim()}
              loading={create.isPending}
            >
              {!create.isPending && <Plus className="h-3.5 w-3.5" />}
              Create key
            </Button>
          </div>
        </div>

        {/* ── List ── */}
        <div className="overflow-hidden rounded-xl border border-border bg-bg-secondary">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-xs uppercase tracking-wider text-text-muted">
                  <th className="px-4 py-3 font-medium">Name</th>
                  <th className="px-4 py-3 font-medium">Key</th>
                  <th className="px-4 py-3 font-medium">Last used</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody>
                {isLoading && (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center">
                      <Loader2 className="mx-auto h-4 w-4 animate-spin text-text-muted" />
                    </td>
                  </tr>
                )}
                {!isLoading && keys.length === 0 && (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center text-text-muted">
                      No API keys yet. Create one for a script, a cron job, or an integration.
                    </td>
                  </tr>
                )}
                {keys.map((k) => {
                  const expired = !!k.expires_at && new Date(k.expires_at) < new Date();
                  const dead = !!k.revoked_at || expired;
                  return (
                    <tr key={k.id} className="border-b border-border/50 last:border-0">
                      <td className="px-4 py-3">
                        <span className={dead ? "text-text-muted line-through" : "text-foreground"}>
                          {k.name}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        {/* Only the prefix survives — the rest was never stored. */}
                        <code className="font-mono text-xs text-text-muted">
                          grit_{k.prefix}_&hellip;
                        </code>
                      </td>
                      <td className="px-4 py-3 text-xs text-text-muted">{relative(k.last_used_at)}</td>
                      <td className="px-4 py-3">
                        {k.revoked_at ? (
                          <span className="rounded-full bg-danger/15 px-2 py-0.5 text-xs text-danger">Revoked</span>
                        ) : expired ? (
                          <span className="rounded-full bg-warning/15 px-2 py-0.5 text-xs text-warning">Expired</span>
                        ) : (
                          <span className="rounded-full bg-success/15 px-2 py-0.5 text-xs text-success">Active</span>
                        )}
                      </td>
                      <td className="px-4 py-3 text-right">
                        {!dead && (
                          <button
                            type="button"
                            onClick={() => setPendingRevoke(k)}
                            className="text-xs text-text-muted hover:text-danger"
                          >
                            <Trash2 className="mr-1 inline h-3 w-3" />
                            Revoke
                          </button>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>

        <p className="mt-4 flex items-start gap-2 text-xs text-text-muted">
          <KeyRound className="mt-0.5 h-3.5 w-3.5 shrink-0" />
          <span>
            Send a key as <code>X-API-Key: grit_&hellip;</code> or{" "}
            <code>Authorization: Bearer grit_&hellip;</code>. Revoked keys stay listed, because
            &ldquo;which key did this?&rdquo; is a question people ask about keys that were
            turned off months ago.
          </span>
        </p>
      </div>

      <ConfirmModal
        open={!!pendingRevoke}
        title="Revoke this key?"
        description={
          "Anything using \"" + (pendingRevoke?.name ?? "") +
          "\" will start getting 401s immediately. This cannot be undone."
        }
        confirmLabel="Revoke"
        variant="danger"
        onCancel={() => setPendingRevoke(null)}
        onConfirm={() => {
          if (pendingRevoke) revoke.mutate(pendingRevoke.id);
          setPendingRevoke(null);
        }}
      />
    </div>
  );
}
`
}
