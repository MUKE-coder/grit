"use client";

import { useEffect, useRef, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useSearchParams } from "next/navigation";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { Shield, ShieldCheck, Download, Trash2, AlertTriangle, Loader2, Check, X, Search } from "@/lib/icons";

interface PickedUser {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
}

// UserPicker — type to search the user directory instead of pasting a UUID.
// GDPR requests arrive as "delete john@acme.com", never as an id, so the raw
// id box was the wrong shape for the job. Selecting a user surfaces their id
// for the record.
function UserPicker({
  value,
  onChange,
}: {
  value: PickedUser | null;
  onChange: (u: PickedUser | null) => void;
}) {
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const boxRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onDoc = (e: MouseEvent) => {
      if (boxRef.current && !boxRef.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDoc);
    return () => document.removeEventListener("mousedown", onDoc);
  }, [open]);

  const searchQ = useQuery({
    queryKey: ["gdpr-user-search", query],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/users", {
        params: { search: query, page_size: 8 },
      });
      return (data.data ?? []) as PickedUser[];
    },
    enabled: open,
  });

  const name = (u: PickedUser) =>
    [u.first_name, u.last_name].filter(Boolean).join(" ") || u.email;

  if (value) {
    return (
      <div className="flex min-w-[320px] flex-1 items-center justify-between gap-3 rounded-lg border border-border bg-background px-3 py-2">
        <span className="min-w-0">
          <span className="block truncate text-sm font-medium text-foreground">{name(value)}</span>
          <span className="block truncate font-mono text-xs text-text-muted">
            {value.email} &middot; {value.id}
          </span>
        </span>
        <button
          type="button"
          onClick={() => onChange(null)}
          className="shrink-0 rounded p-1 text-text-muted hover:text-foreground"
          title="Clear"
        >
          <X className="h-4 w-4" />
        </button>
      </div>
    );
  }

  return (
    <div ref={boxRef} className="relative min-w-[320px] flex-1">
      <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
      <input
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        spellCheck={false}
        placeholder="Search a user by name or email…"
        className="w-full rounded-lg border border-border bg-background py-2 pl-9 pr-3 text-sm outline-none focus:border-accent"
      />
      {open && (
        <div className="absolute z-20 mt-1 max-h-72 w-full overflow-y-auto rounded-lg border border-border bg-bg-elevated shadow-xl">
          {searchQ.isLoading ? (
            <p className="flex items-center gap-2 px-3 py-3 text-sm text-text-muted">
              <Loader2 className="h-3.5 w-3.5 animate-spin" /> Searching…
            </p>
          ) : (searchQ.data ?? []).length === 0 ? (
            <p className="px-3 py-3 text-sm text-text-muted">No users match.</p>
          ) : (
            (searchQ.data ?? []).map((u) => (
              <button
                key={u.id}
                type="button"
                onClick={() => {
                  onChange(u);
                  setOpen(false);
                  setQuery("");
                }}
                className="flex w-full flex-col items-start px-3 py-2 text-left hover:bg-bg-tertiary"
              >
                <span className="text-sm text-foreground">{name(u)}</span>
                <span className="font-mono text-xs text-text-muted">{u.email}</span>
              </button>
            ))
          )}
        </div>
      )}
    </div>
  );
}

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
  const [picked, setPicked] = useState<PickedUser | null>(null);
  const [confirming, setConfirming] = useState(false);
  const [reason, setReason] = useState("");
  const [msg, setMsg] = useState<{ kind: "ok" | "err"; text: string } | null>(null);
  const userId = picked?.id ?? "";

  // The Users table's "Erase (GDPR)" action deep-links here with ?user=<id>,
  // so the request arrives with the subject already chosen.
  const searchParams = useSearchParams();
  const deepLinked = searchParams?.get("user");
  // biome-ignore lint/correctness/useExhaustiveDependencies: loads the deep-linked user once; picking another must not re-fetch it.
  useEffect(() => {
    if (!deepLinked || picked) return;
    apiClient
      .get("/api/users/" + deepLinked)
      .then(({ data }) => setPicked(data.data as PickedUser))
      .catch(() => {});
  }, [deepLinked]);

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
      setPicked(null);
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

      <div className="rounded-xl border border-border bg-bg-elevated p-6">
        <div className="mb-4 flex items-center gap-2">
          <Shield className="h-5 w-5 text-accent" />
          <h2 className="text-lg font-semibold">Export or erase a user</h2>
        </div>
        <p className="mb-4 max-w-2xl text-sm text-text-muted">
          Search for the user this request is about. Export downloads a full JSON copy of their
          data (Art. 15). Erase hard-deletes their personal records and anonymizes the account
          (Art. 17) — this cannot be undone.
        </p>
        <div className="flex flex-wrap items-center gap-3">
          <UserPicker
            value={picked}
            onChange={(u) => {
              setPicked(u);
              setConfirming(false);
            }}
          />
          <button
            type="button"
            disabled={!userId || exportM.isPending}
            onClick={() => exportM.mutate(userId)}
            className="inline-flex items-center gap-2 rounded-lg border border-border bg-background px-4 py-2 text-sm font-medium transition-colors hover:border-accent disabled:opacity-40"
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
                This permanently deletes{" "}
                <span className="font-medium">{picked?.email ?? userId}</span>&rsquo;s uploads,
                sessions, 2FA and more, and anonymizes their account. The audit trail is retained.
              </span>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <input
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder="reason (required, recorded in the journal)"
                className="min-w-[280px] flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-accent"
              />
              {/* The API refuses an erasure with no reason (422). The button is
                  disabled to say so before the request rather than after it. */}
              <button
                type="button"
                disabled={eraseM.isPending || reason.trim().length < 3}
                onClick={() => eraseM.mutate()}
                className="inline-flex items-center gap-2 rounded-lg bg-red-500 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-red-600 disabled:opacity-40"
              >
                {eraseM.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                Confirm erasure
              </button>
              <button
                type="button"
                onClick={() => setConfirming(false)}
                className="rounded-lg px-3 py-2 text-sm text-text-muted hover:text-foreground"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>

      <div className="rounded-xl border border-border bg-bg-elevated p-6">
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

        <p className="mb-4 text-xs text-text-muted">
          Erasures only. Deleting a user from the Users page is an ordinary, reversible soft
          delete — it keeps their data and is not recorded here. Use <strong>Erase</strong> above
          for a real Art. 17 erasure.
        </p>

        {journalQ.isLoading ? (
          <div className="flex items-center gap-2 py-8 text-sm text-text-muted">
            <Loader2 className="h-4 w-4 animate-spin" /> Loading&hellip;
          </div>
        ) : !journalQ.data?.rows.length ? (
          <p className="py-8 text-center text-sm text-text-muted">
            No erasures recorded yet. When you erase a user, an entry appears here.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-text-muted">
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
                    <td className="py-2 pr-4 text-text-muted">{r.reason || "—"}</td>
                    <td className="py-2 pr-4 text-text-muted">
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
