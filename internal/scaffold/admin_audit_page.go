package scaffold

// The tamper-evident activity log, as a screen.
//
// This is not the same thing as /system/activity. Two logs ship, deliberately:
//
//   UserActivity  — semantic events ("ticket.create"), soft-deletable, the
//                   operator's timeline. That is what /system/activity shows.
//   ActivityLog   — every authenticated mutation, hash-chained, append-only.
//                   Evidence. That is this page.
//
// The chain was implemented, the verify endpoint was implemented, and nothing
// ever called either — so the one feature a compliance reviewer actually wants
// to see was invisible. A hash chain nobody can inspect is indistinguishable
// from a table with a hash column.
//
// The verify button is deliberately explicit rather than run-on-load: the walk
// is O(n) over the whole table, and firing it every time someone opens a page
// is how you turn an audit feature into a performance incident.

func adminAuditPage() string {
	return `"use client";

import { useState } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { ShieldCheck, AlertTriangle, Loader2, Check, Search } from "@/lib/icons";
import { Button } from "@/components/ui/button";

interface AuditEntry {
  id: string;
  user_id: string;
  method: string;
  path: string;
  status: number;
  payload_digest: string;
  ip_address: string;
  user_agent: string;
  duration_ms: number;
  prev_hash: string;
  hash: string;
  created_at: string;
}

interface ChainStatus {
  valid: boolean;
  total_entries: number;
  broken_at?: number;
  broken_at_id?: string;
  expected?: string;
  got?: string;
  message?: string;
}

const METHOD_TONE: Record<string, string> = {
  POST: "bg-success/15 text-success",
  PUT: "bg-info/15 text-info",
  PATCH: "bg-info/15 text-info",
  DELETE: "bg-danger/15 text-danger",
};

function statusTone(status: number): string {
  if (status >= 500) return "text-danger";
  if (status >= 400) return "text-warning";
  return "text-success";
}

export default function AuditLogPage() {
  const [method, setMethod] = useState("");
  const [path, setPath] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["audit-log", method, path, page],
    queryFn: async () => {
      const params = new URLSearchParams({ page: String(page), page_size: "25" });
      if (method) params.set("method", method);
      if (path) params.set("path", path);
      const { data: res } = await apiClient.get("/api/admin/activity?" + params.toString());
      return res as { data: AuditEntry[]; meta?: { total: number; pages: number } };
    },
  });

  // Explicit, not on mount — see the file header.
  const verify = useMutation({
    mutationFn: async () => {
      const { data: res } = await apiClient.get("/api/admin/activity/integrity");
      return (res.data ?? res) as ChainStatus;
    },
  });

  const entries = data?.data ?? [];
  const pages = data?.meta?.pages ?? 1;
  const result = verify.data;

  return (
    <div>
      <PageHeader
        title="Audit log"
        subtitle="Every authenticated write, hash-chained so a changed row can be detected."
      />

      <div className="px-6 py-6 space-y-5">
        {/* ── Integrity ───────────────────────────────────────────── */}
        <div className="rounded-xl border border-border bg-bg-secondary p-5">
          <div className="flex flex-wrap items-start gap-4">
            <ShieldCheck className="mt-0.5 h-4 w-4 shrink-0 text-accent" />
            <div className="min-w-0 flex-1">
              <h2 className="font-semibold text-foreground">Chain integrity</h2>
              <p className="text-sm text-foreground-secondary">
                Recomputes every hash from the first entry. A mismatch means a row was
                modified, deleted, or inserted out of order — the result names the first one.
              </p>
            </div>
            <Button onClick={() => verify.mutate()} loading={verify.isPending}>
              {verify.isPending ? "Verifying…" : "Verify chain"}
            </Button>
          </div>

          {result && (
            <div
              className={
                "mt-4 rounded-lg border p-4 " +
                (result.valid
                  ? "border-success/40 bg-success/[0.06]"
                  : "border-danger/40 bg-danger/[0.06]")
              }
            >
              <div className="flex items-center gap-2">
                {result.valid ? (
                  <Check className="h-4 w-4 text-success" />
                ) : (
                  <AlertTriangle className="h-4 w-4 text-danger" />
                )}
                <span className="text-sm font-medium text-foreground">
                  {result.valid
                    ? "Chain verified — " + result.total_entries.toLocaleString() + " entries"
                    : "Chain broken"}
                </span>
              </div>
              {!result.valid && (
                <div className="mt-2 space-y-1 text-xs text-foreground-secondary">
                  <p>{result.message}</p>
                  <p>
                    First bad entry: position {result.broken_at}, id{" "}
                    <code className="text-foreground">{result.broken_at_id}</code>
                  </p>
                  <p className="font-mono break-all">expected {result.expected}</p>
                  <p className="font-mono break-all">got {result.got}</p>
                  <p className="pt-1">
                    Everything before that position is still trustworthy.
                  </p>
                </div>
              )}
            </div>
          )}

          {verify.isError && (
            <p className="mt-3 text-sm text-danger">
              Could not run the check. Very large logs can exceed the 60-second limit — run it
              from a cron job instead.
            </p>
          )}
        </div>

        {/* ── Filters ─────────────────────────────────────────────── */}
        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-muted" />
            <input
              value={path}
              onChange={(e) => { setPath(e.target.value); setPage(1); }}
              placeholder="Path starts with /api/v1/…"
              className="w-72 rounded-lg border border-border bg-bg-primary py-2 pl-9 pr-3 text-sm text-foreground"
            />
          </div>
          {["", "POST", "PUT", "PATCH", "DELETE"].map((m) => (
            <button
              key={m || "all"}
              type="button"
              onClick={() => { setMethod(m); setPage(1); }}
              className={
                "rounded-lg border px-3 py-1.5 text-xs font-medium transition-colors " +
                (method === m
                  ? "border-accent/40 bg-accent/10 text-foreground"
                  : "border-border text-text-muted hover:text-foreground")
              }
            >
              {m || "All"}
            </button>
          ))}
        </div>

        {/* ── Entries ─────────────────────────────────────────────── */}
        <div className="overflow-hidden rounded-xl border border-border bg-bg-secondary">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-xs uppercase tracking-wider text-text-muted">
                  <th className="px-4 py-3 font-medium">When</th>
                  <th className="px-4 py-3 font-medium">Request</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-4 py-3 font-medium">Who</th>
                  <th className="px-4 py-3 font-medium">Body digest</th>
                </tr>
              </thead>
              <tbody>
                {isLoading && (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center text-text-muted">
                      <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                    </td>
                  </tr>
                )}
                {!isLoading && entries.length === 0 && (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center text-text-muted">
                      No entries yet. Every authenticated create, update or delete lands here.
                    </td>
                  </tr>
                )}
                {entries.map((e) => (
                  <tr key={e.id} className="border-b border-border/50 last:border-0">
                    <td className="whitespace-nowrap px-4 py-3 text-xs text-text-muted">
                      {new Date(e.created_at).toLocaleString()}
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={
                          "mr-2 rounded px-1.5 py-0.5 font-mono text-[10px] " +
                          (METHOD_TONE[e.method] ?? "bg-bg-tertiary text-text-muted")
                        }
                      >
                        {e.method}
                      </span>
                      <span className="font-mono text-xs text-foreground">{e.path}</span>
                    </td>
                    <td className={"px-4 py-3 font-mono text-xs " + statusTone(e.status)}>
                      {e.status}
                      <span className="ml-2 text-text-muted">{e.duration_ms}ms</span>
                    </td>
                    <td className="px-4 py-3 text-xs text-text-muted">
                      <span className="font-mono">{e.user_id ? e.user_id.slice(0, 8) : "—"}</span>
                      <span className="ml-2">{e.ip_address}</span>
                    </td>
                    <td className="px-4 py-3">
                      {/* The body itself is never stored — only a digest of it. */}
                      <code className="font-mono text-[11px] text-text-muted">
                        {e.payload_digest ? e.payload_digest.slice(0, 12) : "—"}
                      </code>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {pages > 1 && (
            <div className="flex items-center justify-between border-t border-border px-4 py-3 text-sm">
              <span className="text-text-muted">
                Page {page} of {pages}
              </span>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="rounded-lg border border-border px-3 py-1.5 text-xs disabled:opacity-40"
                >
                  Previous
                </button>
                <button
                  type="button"
                  onClick={() => setPage((p) => Math.min(pages, p + 1))}
                  disabled={page >= pages}
                  className="rounded-lg border border-border px-3 py-1.5 text-xs disabled:opacity-40"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>

        <p className="text-xs text-text-muted">
          Request bodies are stored as a SHA-256 digest, not verbatim — evidence of what was
          sent without keeping the personal data in it. This log is append-only; the weekly{" "}
          <code>audit:prune</code> job trims old entries and re-anchors the chain.
        </p>
      </div>
    </div>
  );
}
`
}
