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
import { ShieldCheck, AlertTriangle, Loader2, Check, Search, FileText } from "@/lib/icons";
import { Button } from "@/components/ui/button";
import { Pill, type PillTone } from "@/components/ui/pill";
import { Note } from "@/components/ui/note";
import { EmptyState } from "@/components/ui/empty-state";
import { Segmented } from "@/components/ui/segmented";
import { inputClasses } from "@/components/ui/input";

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
  resource?: string;
  resource_ids?: string;
  record_count?: number;
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

// The method, as the same pill the rest of the panel uses for a state.
//
// These were six hand-written class pairs at four different opacities, which
// is how a POST here ended up a different green from a Published there.
const METHOD_TONE: Record<string, PillTone> = {
  SECURITY: "warn",
  GET: "accent",
  POST: "good",
  PUT: "info",
  PATCH: "info",
  DELETE: "bad",
};

type AuditMethod = "" | "GET" | "POST" | "PUT" | "PATCH" | "DELETE" | "SECURITY";

/** The method filter. "" is every method, which is where it starts. */
const METHODS: ReadonlyArray<{ value: AuditMethod; label: string }> = [
  { value: "", label: "All" },
  { value: "GET", label: "GET" },
  { value: "POST", label: "POST" },
  { value: "PUT", label: "PUT" },
  { value: "PATCH", label: "PATCH" },
  { value: "DELETE", label: "DELETE" },
  { value: "SECURITY", label: "Security" },
];

function statusTone(status: number): string {
  if (status >= 500) return "text-danger";
  if (status >= 400) return "text-warning";
  return "text-success";
}

// A reseal trusts every entry from the break onward as it stands now, so it
// is deliberate: the id has to be typed, and the server records who did it.
function ResealPanel({ fromId, onDone }: { fromId?: string; onDone: () => void }) {
  const [confirm, setConfirm] = useState("");
  const reseal = useMutation({
    mutationFn: async () => {
      const { data: res } = await apiClient.post("/api/admin/activity/reseal", { from_id: fromId });
      return res as { data: { resealed: number } };
    },
    onSuccess: () => onDone(),
  });
  if (!fromId) return null;
  const short = fromId.slice(0, 8);
  const failure = (reseal.error as { response?: { data?: { error?: { message?: string } } } } | null)
    ?.response?.data?.error?.message;

  return (
    <div className="mt-3 rounded-lg border border-border bg-bg-elevated p-3">
      <p className="text-xs text-text-secondary">
        If nobody changed these entries, you can reseal from this one. A chain written before
        v3.215.0 on Postgres or MySQL fails on its first entry, because the database rounded the
        timestamps it was hashed with. A reseal trusts every entry from here on as it stands now,
        and is itself recorded in the chain under your name. Type{" "}
        <code className="text-foreground">{short}</code> to confirm.
      </p>
      <div className="mt-2 flex flex-wrap items-center gap-2">
        <input
          value={confirm}
          onChange={(e) => setConfirm(e.target.value)}
          placeholder={short}
          aria-label="Type the first eight characters of the entry id to confirm the reseal"
          className={inputClasses({ inputSize: "sm", fullWidth: false }) + " w-40 font-mono text-xs"}
        />
        <Button onClick={() => reseal.mutate()} loading={reseal.isPending} disabled={confirm !== short}>
          Reseal from this entry
        </Button>
      </div>
      {reseal.isError && (
        <p className="mt-2 text-xs text-danger">{failure ?? "Could not reseal the chain."}</p>
      )}
      {reseal.isSuccess && (
        <p className="mt-2 text-xs text-success">
          Resealed {reseal.data.data.resealed.toLocaleString()} entries.
        </p>
      )}
    </div>
  );
}

export default function AuditLogPage() {
  const [method, setMethod] = useState<AuditMethod>("");
  const [path, setPath] = useState("");
  const [record, setRecord] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["audit-log", method, path, record, page],
    queryFn: async () => {
      const params = new URLSearchParams({ page: String(page), page_size: "25" });
      if (method) params.set("method", method);
      if (path) params.set("path", path);
      if (record) params.set("record", record);
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

  // Walking an empty chain finds no mismatch, so the endpoint answers valid.
  // That is true and it is not evidence, and a green Verified badge over "0
  // entries" claims it is. Said plainly instead, in the colour of a thing that
  // has not happened.
  const empty = !!result && result.valid && result.total_entries === 0;

  return (
    <div>
      <PageHeader
        title="Audit log"
        subtitle="Every authenticated write, and every read of a resource generated with --audit-reads, hash-chained so a changed row can be detected."
      />

      <div className="space-y-6">
        {/* Chain integrity.
            The one card on the page that answers a yes-or-no question, so the
            one card that changes colour: green when the walk passed, red when
            it did not, and the neutral card it starts as until somebody asks.

            The explanation sits above the button rather than in a tooltip,
            because the check is O(n) over the whole table and that is worth
            knowing before pressing it. */}
        <section
          className={
            "rounded-2xl border p-5 transition-grit " +
            (!result || empty
              ? "border-border bg-bg-elevated"
              : result.valid
                ? "border-success-line bg-success-soft"
                : "border-danger-line bg-danger-soft")
          }
        >
          <div className="flex flex-wrap items-start gap-4">
            <span
              className={
                "mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg " +
                (!result || empty
                  ? "bg-accent-soft text-accent"
                  : result.valid
                    ? "bg-success-soft text-success"
                    : "bg-danger-soft text-danger")
              }
              aria-hidden="true"
            >
              {!result || empty ? (
                <ShieldCheck className="h-4 w-4" />
              ) : result.valid ? (
                <Check className="h-4 w-4" />
              ) : (
                <AlertTriangle className="h-4 w-4" />
              )}
            </span>
            <div className="min-w-0 flex-1">
              <h2 className="text-[15px] font-semibold text-foreground">Chain integrity</h2>
              <p className="mt-0.5 text-[13px] text-text-secondary">
                Recomputes every hash from the first entry. A mismatch means a row was modified,
                deleted, or inserted out of order, and the result names the first one.
              </p>
            </div>
            <Button onClick={() => verify.mutate()} loading={verify.isPending}>
              {verify.isPending ? "Verifying…" : result ? "Verify again" : "Verify chain"}
            </Button>
          </div>

          {result && (
            <div className="mt-4 border-t border-border/60 pt-4">
              <div className="flex flex-wrap items-center gap-2">
                <Pill tone={empty ? "neutral" : result.valid ? "good" : "bad"} dot>
                  {empty ? "Nothing to check" : result.valid ? "Verified" : "Broken"}
                </Pill>
                <span className="text-sm text-foreground">
                  {empty
                    ? "The log is empty, so there is no chain to verify yet."
                    : result.valid
                      ? result.total_entries.toLocaleString() +
                        (result.total_entries === 1 ? " entry" : " entries") +
                        " checked just now"
                      : "A row does not match the hash before it"}
                </span>
              </div>
              {!result.valid && (
                <div className="mt-3 space-y-1 text-xs text-text-secondary">
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
                  <ResealPanel fromId={result.broken_at_id} onDone={() => verify.mutate()} />
                </div>
              )}
            </div>
          )}

          {verify.isError && (
            <p className="mt-3 text-sm text-danger">
              Could not run the check. Very large logs can exceed the 60-second limit, so run it
              from a cron job instead.
            </p>
          )}
        </section>

        {/* The entries, with their filters in the card's own header.
            They used to float above the table as a loose row of controls,
            which reads as page chrome rather than as "these three narrow the
            thing directly below them". */}
        <section className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
          <header className="flex flex-wrap items-center gap-2 border-b border-border px-5 py-4">
            <div className="relative min-w-56 flex-1">
              <Search
                className="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-text-muted"
                aria-hidden="true"
              />
              <input
                value={path}
                onChange={(e) => { setPath(e.target.value); setPage(1); }}
                placeholder="Path starts with /api/v1/…"
                aria-label="Filter by the start of the request path"
                className={inputClasses({ inputSize: "sm", className: "pl-9" })}
              />
            </div>
            <input
              value={record}
              onChange={(e) => { setRecord(e.target.value); setPage(1); }}
              placeholder="Record id"
              aria-label="Show everyone who read or changed this record"
              className={inputClasses({ inputSize: "sm", fullWidth: false }) + " w-52 font-mono text-xs"}
            />
            <Segmented
              label="Method"
              value={method}
              onChange={(m) => { setMethod(m); setPage(1); }}
              options={METHODS}
            />
          </header>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-xs font-medium uppercase tracking-wider text-text-muted">
                  <th className="px-5 py-3 font-medium">When</th>
                  <th className="px-4 py-3 font-medium">Request</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-4 py-3 font-medium">Who</th>
                  <th className="px-5 py-3 font-medium">Digest</th>
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
                    <td colSpan={5} className="p-0">
                      {/* Two different empty tables, and they are not the same
                          news: nothing has happened yet, or the filters have
                          excluded everything that has. One sentence told both
                          stories and was wrong about one of them. */}
                      {method || path || record ? (
                        <EmptyState
                          icon={<Search className="h-5 w-5" />}
                          title="Nothing matches these filters"
                          description="No entry matches that method, path and record together. Widen one of them."
                        />
                      ) : (
                        <EmptyState
                          icon={<FileText className="h-5 w-5" />}
                          title="Nothing recorded yet"
                          description="Every authenticated create, update and delete lands here, as does every read of a resource generated with --audit-reads."
                          hint="grit generate resource Invoice --audit-reads"
                        />
                      )}
                    </td>
                  </tr>
                )}
                {entries.map((e) => (
                  <tr key={e.id} className="border-b border-border/50 last:border-0">
                    <td className="whitespace-nowrap px-5 py-3 font-mono text-xs text-text-muted">
                      {new Date(e.created_at).toLocaleString()}
                    </td>
                    <td className="px-4 py-3">
                      <Pill tone={METHOD_TONE[e.method] ?? "neutral"} className="mr-2 font-mono">
                        {e.method}
                      </Pill>
                      <span className="font-mono text-xs text-foreground">{e.path}</span>
                      {e.resource && (
                        <span className="ml-2 text-xs text-text-muted">
                          {e.resource}
                          {typeof e.record_count === "number" && e.record_count > 0
                            ? " · " + e.record_count.toLocaleString() + (e.record_count === 1 ? " row" : " rows")
                            : ""}
                        </span>
                      )}
                    </td>
                    <td className={"px-4 py-3 font-mono text-xs " + statusTone(e.status)}>
                      {e.status}
                      <span className="ml-2 text-text-muted">{e.duration_ms}ms</span>
                    </td>
                    <td className="px-4 py-3 text-xs text-text-muted">
                      <span className="font-mono">{e.user_id ? e.user_id.slice(0, 8) : "—"}</span>
                      <span className="ml-2">{e.ip_address}</span>
                    </td>
                    <td className="px-5 py-3">
                      {/* The body itself is never stored, only a digest of it. */}
                      <code
                        className="font-mono text-[11px] text-text-muted"
                        title={e.payload_digest || undefined}
                      >
                        {e.payload_digest ? e.payload_digest.slice(0, 12) : "—"}
                      </code>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {pages > 1 && (
            <div className="flex items-center justify-between border-t border-border px-5 py-3 text-sm">
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
        </section>

        <Note>
          Request bodies are stored as a SHA-256 digest rather than verbatim: evidence of what was
          sent, without keeping the personal data in it. Reads keep a digest of the query string,
          so a name typed into a search box is not stored either. This log is append-only, and the
          weekly <code className="font-mono">audit:prune</code> job trims old entries and
          re-anchors the chain.
        </Note>
      </div>
    </div>
  );
}
`
}
