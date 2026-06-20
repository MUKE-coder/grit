package scaffold

// v3.30 admin pages.
//
//   app/(dashboard)/system/activity/page.tsx — semantic activity log
//   app/(dashboard)/system/support/page.tsx  — ticket list + create
//   app/(dashboard)/system/support/[id]/page.tsx — ticket thread + reply
//
// All three use the v3.29 chrome (PageHeader, ResponsiveSheet, IconButton)
// + responsive table pattern so they look right on every theme + viewport.

// adminActivityDashboardPage emits the semantic activity log dashboard.
func adminActivityDashboardPage() string {
	return `"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ResponsiveTable, type TableColumn } from "@/components/ui/ResponsiveTable";
import { IconButton } from "@/components/ui/IconButton";
import { Download, AlertCircle, AlertTriangle, Info } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";
import { exportToExcel } from "@/lib/export";

interface ActivityRow {
  id: string;
  user_id: string;
  action: string;
  severity: "info" | "warn" | "critical";
  summary: string;
  resource_type: string;
  resource_id: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
}

interface StatsResponse {
  data: { info: number; warn: number; critical: number; total: number };
}

interface ListResponse {
  data: ActivityRow[];
  meta?: { total: number; page: number; page_size: number; pages: number };
}

const SEVERITIES = ["all", "info", "warn", "critical"] as const;
type Severity = (typeof SEVERITIES)[number];

const severityClass: Record<ActivityRow["severity"], string> = {
  info: "bg-info/10 text-info",
  warn: "bg-warning/10 text-warning",
  critical: "bg-danger/10 text-danger",
};

const severityIcon: Record<ActivityRow["severity"], typeof Info> = {
  info: Info,
  warn: AlertTriangle,
  critical: AlertCircle,
};

export default function ActivityPage() {
  const [search, setSearch] = useState("");
  const [severity, setSeverity] = useState<Severity>("all");

  const { data: stats } = useQuery<StatsResponse["data"]>({
    queryKey: ["user-activity", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get<StatsResponse>("/api/user-activity/stats");
      return data.data;
    },
    refetchInterval: 60_000,
  });

  const { data, isLoading } = useQuery<ActivityRow[]>({
    queryKey: ["user-activity", "list", search, severity],
    queryFn: async () => {
      const params = new URLSearchParams({ page_size: "100" });
      if (search) params.set("q", search);
      if (severity !== "all") params.set("severity", severity);
      const { data } = await apiClient.get<ListResponse>("/api/user-activity?" + params.toString());
      return data.data;
    },
  });

  const rows = data || [];

  const onExport = async () => {
    const payload = rows.map((r) => ({
      Time: new Date(r.created_at).toLocaleString(),
      Severity: r.severity,
      Action: r.action,
      Summary: r.summary,
      User: r.user_id,
      Resource: r.resource_type + ":" + r.resource_id,
      IP: r.ip_address,
      UserAgent: r.user_agent,
    }));
    await exportToExcel(payload, "user-activity-" + new Date().toISOString().slice(0, 10));
  };

  const columns: TableColumn<ActivityRow>[] = [
    {
      key: "severity",
      header: "Severity",
      cell: (r) => {
        const Icon = severityIcon[r.severity];
        return (
          <span className={"inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-xs font-medium uppercase tracking-wide " + severityClass[r.severity]}>
            <Icon className="h-3 w-3" />
            {r.severity}
          </span>
        );
      },
    },
    { key: "action", header: "Action", cell: (r) => <code className="font-mono text-xs">{r.action}</code> },
    { key: "summary", header: "Summary", cell: (r) => <span className="text-sm">{r.summary}</span> },
    { key: "ip", header: "IP", cell: (r) => <span className="font-mono text-xs text-text-muted">{r.ip_address || "—"}</span>, hideOnMobile: true },
    {
      key: "time",
      header: "When",
      cell: (r) => <span className="text-xs text-text-muted">{new Date(r.created_at).toLocaleString()}</span>,
      align: "right",
    },
  ];

  return (
    <div>
      <PageHeader
        title="User Activity"
        subtitle="Every recorded action across the platform — auth events, writes, and operator actions"
        searchPlaceholder="Search by summary..."
        searchValue={search}
        onSearchChange={setSearch}
        actions={
          <IconButton
            variant="secondary"
            icon={<Download className="h-4 w-4" />}
            label="Export"
            onClick={onExport}
          />
        }
      />

      {/* Stats chips */}
      <div className="mb-6 grid grid-cols-2 gap-3 md:grid-cols-4">
        <StatChip label="Last 24h" value={stats?.total ?? 0} tone="muted" />
        <StatChip label="Info" value={stats?.info ?? 0} tone="info" />
        <StatChip label="Warn" value={stats?.warn ?? 0} tone="warn" />
        <StatChip label="Critical" value={stats?.critical ?? 0} tone="critical" />
      </div>

      {/* Severity filter tabs */}
      <div className="mb-4 flex flex-wrap gap-2">
        {SEVERITIES.map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setSeverity(s)}
            className={
              "rounded-full border px-3 py-1 text-xs font-medium capitalize transition-colors " +
              (severity === s
                ? "border-accent bg-accent text-white"
                : "border-border bg-bg-elevated text-text-secondary hover:bg-bg-hover")
            }
          >
            {s}
          </button>
        ))}
      </div>

      <ResponsiveTable
        columns={columns}
        rows={rows}
        rowKey={(r) => r.id}
        loading={isLoading}
        emptyMessage="No activity matches the current filter."
      />
    </div>
  );
}

function StatChip({ label, value, tone }: { label: string; value: number; tone: "muted" | "info" | "warn" | "critical" }) {
  const toneClass = {
    muted: "border-border bg-bg-elevated",
    info: "border-info/30 bg-info/5",
    warn: "border-warning/30 bg-warning/5",
    critical: "border-danger/30 bg-danger/5",
  }[tone];

  return (
    <div className={"rounded-xl border px-4 py-3 " + toneClass}>
      <p className="text-xs font-medium uppercase tracking-wide text-text-muted">{label}</p>
      <p className="mt-1 text-2xl font-bold text-foreground">{value}</p>
    </div>
  );
}
`
}

// adminSupportListPage emits the ticket list page with status tabs and
// a new-ticket sheet form.
func adminSupportListPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import { IconButton } from "@/components/ui/IconButton";
import { Plus, MessageSquare, AlertCircle } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface Ticket {
  id: string;
  user_id: string;
  subject: string;
  description: string;
  status: "open" | "closed";
  priority: "low" | "medium" | "high" | "critical";
  labels: string;
  assignee_id: string;
  last_reply_at: string | null;
  created_at: string;
  user?: { first_name: string; last_name: string; email: string };
}

interface ListResponse {
  data: Ticket[];
  meta?: { total: number };
}

const priorityClass: Record<Ticket["priority"], string> = {
  low: "bg-bg-hover text-text-secondary",
  medium: "bg-info/10 text-info",
  high: "bg-warning/10 text-warning",
  critical: "bg-danger/10 text-danger",
};

export default function SupportPage() {
  const [status, setStatus] = useState<"open" | "closed">("open");
  const [openSheet, setOpenSheet] = useState(false);

  const { data, isLoading } = useQuery<Ticket[]>({
    queryKey: ["tickets", "list", status],
    queryFn: async () => {
      const { data } = await apiClient.get<ListResponse>("/api/tickets?status=" + status);
      return data.data;
    },
  });

  const tickets = data || [];

  return (
    <div>
      <PageHeader
        title="Support"
        subtitle="Open a ticket to reach our team. We'll reply on the ticket and you'll also get an email + an in-app notification."
        actions={
          <IconButton
            icon={<Plus className="h-4 w-4" />}
            label="New ticket"
            onClick={() => setOpenSheet(true)}
          />
        }
      />

      {/* Status tabs */}
      <div className="mb-6 flex w-fit rounded-lg border border-border bg-bg-elevated p-1">
        {(["open", "closed"] as const).map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setStatus(s)}
            className={
              "inline-flex items-center gap-2 rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-colors " +
              (status === s ? "bg-accent text-white" : "text-text-secondary hover:text-foreground")
            }
          >
            {s === "open" ? <MessageSquare className="h-3.5 w-3.5" /> : <AlertCircle className="h-3.5 w-3.5" />}
            {s}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
          <div className="inline-block h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
        </div>
      ) : tickets.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
          <MessageSquare className="mx-auto h-10 w-10 text-text-muted" />
          <p className="mt-3 text-base font-medium text-foreground">No {status} tickets</p>
          <p className="mt-1 text-sm text-text-muted">
            {status === "open"
              ? "When you open a support ticket it'll show up here and our team gets notified."
              : "Closed tickets will appear here once they're resolved."}
          </p>
          {status === "open" && (
            <button
              type="button"
              onClick={() => setOpenSheet(true)}
              className="mt-4 inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover"
            >
              <Plus className="h-4 w-4" />
              New ticket
            </button>
          )}
        </div>
      ) : (
        <ul className="space-y-2">
          {tickets.map((t) => (
            <li key={t.id}>
              <Link
                href={"/system/support/" + t.id}
                className="block rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className={"inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium capitalize " + priorityClass[t.priority]}>
                        {t.priority}
                      </span>
                      <p className="truncate text-sm font-semibold text-foreground">{t.subject}</p>
                    </div>
                    <p className="mt-1 line-clamp-2 text-sm text-text-secondary">{t.description}</p>
                    <p className="mt-2 text-xs text-text-muted">
                      Opened {new Date(t.created_at).toLocaleString()}
                      {t.user && " by " + t.user.first_name + " " + t.user.last_name}
                    </p>
                  </div>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}

      <NewTicketSheet open={openSheet} onClose={() => setOpenSheet(false)} />
    </div>
  );
}

function NewTicketSheet({ open, onClose }: { open: boolean; onClose: () => void }) {
  const queryClient = useQueryClient();
  const [subject, setSubject] = useState("");
  const [priority, setPriority] = useState<Ticket["priority"]>("medium");
  const [labels, setLabels] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  const create = useMutation({
    mutationFn: async () => {
      return apiClient.post("/api/tickets", { subject, priority, labels, description });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tickets"] });
      setSubject(""); setPriority("medium"); setLabels(""); setDescription(""); setError("");
      onClose();
    },
    onError: (err: unknown) => {
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Couldn't open ticket. Please try again.");
    },
  });

  const submit = () => {
    if (!subject.trim() || !description.trim()) {
      setError("Subject and description are required.");
      return;
    }
    create.mutate();
  };

  return (
    <ResponsiveSheet
      open={open}
      onClose={onClose}
      title="New ticket"
      description="Tell us what's going wrong (or what you'd like to see)."
      footer={
        <>
          <button type="button" onClick={onClose} className="rounded-lg px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover">
            Cancel
          </button>
          <button
            type="button"
            onClick={submit}
            disabled={create.isPending}
            className="rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
          >
            {create.isPending ? "Opening..." : "Open ticket"}
          </button>
        </>
      }
    >
      <form
        onSubmit={(e) => { e.preventDefault(); submit(); }}
        className="space-y-4"
      >
        {error && (
          <div className="rounded-lg border border-danger/20 bg-danger/10 px-3 py-2 text-sm text-danger">{error}</div>
        )}

        <Field label="Subject" required>
          <input
            type="text"
            value={subject}
            onChange={(e) => setSubject(e.target.value)}
            placeholder="One-line summary of the problem"
            className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </Field>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <Field label="Priority">
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value as Ticket["priority"])}
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="critical">Critical</option>
            </select>
          </Field>
          <Field label="Labels (comma-separated, up to 8)">
            <input
              type="text"
              value={labels}
              onChange={(e) => setLabels(e.target.value)}
              placeholder="bug, billing, mobile"
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            />
          </Field>
        </div>

        <Field label="Describe what's happening" required>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={5}
            placeholder="Steps to reproduce, expected vs. actual behaviour, anything that helps."
            className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </Field>
      </form>
    </ResponsiveSheet>
  );
}

function Field({ label, required, children }: { label: string; required?: boolean; children: React.ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-xs font-semibold uppercase tracking-wide text-text-muted">
        {label}
        {required && <span className="ml-1 text-danger">*</span>}
      </span>
      {children}
    </label>
  );
}
`
}

// adminTicketThreadPage emits the per-ticket thread view with reply form.
func adminTicketThreadPage() string {
	return `"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { Check, ArrowLeft } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface Reply {
  id: string;
  ticket_id: string;
  user_id: string;
  body: string;
  is_admin_reply: boolean;
  created_at: string;
  user?: { first_name: string; last_name: string; email: string };
}

interface Ticket {
  id: string;
  subject: string;
  description: string;
  status: "open" | "closed";
  priority: "low" | "medium" | "high" | "critical";
  labels: string;
  created_at: string;
  user?: { first_name: string; last_name: string; email: string };
  replies: Reply[];
}

const priorityClass: Record<Ticket["priority"], string> = {
  low: "bg-bg-hover text-text-secondary",
  medium: "bg-info/10 text-info",
  high: "bg-warning/10 text-warning",
  critical: "bg-danger/10 text-danger",
};

export default function TicketThreadPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [reply, setReply] = useState("");

  const { data: ticket, isLoading } = useQuery<Ticket>({
    queryKey: ["ticket", params.id],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: Ticket }>("/api/tickets/" + params.id);
      return data.data;
    },
    enabled: !!params.id,
  });

  const replyM = useMutation({
    mutationFn: async () => apiClient.post("/api/tickets/" + params.id + "/reply", { body: reply }),
    onSuccess: () => {
      setReply("");
      queryClient.invalidateQueries({ queryKey: ["ticket", params.id] });
    },
  });

  const close = useMutation({
    mutationFn: async () => apiClient.patch("/api/tickets/" + params.id + "/close"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["ticket", params.id] }),
  });

  const reopen = useMutation({
    mutationFn: async () => apiClient.patch("/api/tickets/" + params.id + "/reopen"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["ticket", params.id] }),
  });

  if (isLoading) {
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
        <div className="inline-block h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
      </div>
    );
  }

  if (!ticket) {
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
        <p className="text-base font-medium text-foreground">Ticket not found</p>
        <p className="mt-1 text-sm text-text-muted">It may have been deleted or you don't have access.</p>
      </div>
    );
  }

  const labels = ticket.labels ? ticket.labels.split(",").map((l) => l.trim()).filter(Boolean) : [];

  return (
    <div>
      <PageHeader
        title={ticket.subject}
        subtitle={"Opened " + new Date(ticket.created_at).toLocaleString()}
        actions={
          <>
            <IconButton
              variant="secondary"
              icon={<ArrowLeft className="h-4 w-4" />}
              label="Back"
              onClick={() => router.push("/system/support")}
            />
            {ticket.status === "open" ? (
              <IconButton
                variant="secondary"
                icon={<Check className="h-4 w-4" />}
                label="Close"
                onClick={() => close.mutate()}
              />
            ) : (
              <IconButton
                variant="secondary"
                icon={<Check className="h-4 w-4" />}
                label="Reopen"
                onClick={() => reopen.mutate()}
              />
            )}
          </>
        }
      />

      {/* Meta chips */}
      <div className="mb-6 flex flex-wrap items-center gap-2">
        <span className={"inline-flex items-center rounded-md px-2.5 py-1 text-xs font-medium uppercase " + (ticket.status === "open" ? "bg-success/10 text-success" : "bg-text-muted/10 text-text-muted")}>
          {ticket.status}
        </span>
        <span className={"inline-flex items-center rounded-md px-2.5 py-1 text-xs font-medium capitalize " + priorityClass[ticket.priority]}>
          {ticket.priority} priority
        </span>
        {labels.map((l) => (
          <span key={l} className="inline-flex items-center rounded-md border border-border bg-bg-elevated px-2.5 py-1 text-xs font-medium text-text-secondary">
            {l}
          </span>
        ))}
      </div>

      {/* Original message */}
      <article className="rounded-xl border border-border bg-bg-elevated p-5">
        <header className="mb-3 flex items-center justify-between">
          <p className="text-sm font-semibold text-foreground">
            {ticket.user ? ticket.user.first_name + " " + ticket.user.last_name : "Unknown"}
          </p>
          <p className="text-xs text-text-muted">{new Date(ticket.created_at).toLocaleString()}</p>
        </header>
        <p className="whitespace-pre-wrap text-sm text-foreground">{ticket.description}</p>
      </article>

      {/* Replies */}
      {ticket.replies && ticket.replies.length > 0 && (
        <ul className="mt-4 space-y-3">
          {ticket.replies.map((r) => (
            <li
              key={r.id}
              className={
                "rounded-xl border p-5 " +
                (r.is_admin_reply
                  ? "border-accent/30 bg-accent/5"
                  : "border-border bg-bg-elevated")
              }
            >
              <header className="mb-3 flex items-center justify-between">
                <p className="text-sm font-semibold text-foreground">
                  {r.user ? r.user.first_name + " " + r.user.last_name : "Unknown"}
                  {r.is_admin_reply && (
                    <span className="ml-2 rounded bg-accent/20 px-1.5 py-0.5 text-[10px] font-bold uppercase text-accent">
                      Staff
                    </span>
                  )}
                </p>
                <p className="text-xs text-text-muted">{new Date(r.created_at).toLocaleString()}</p>
              </header>
              <p className="whitespace-pre-wrap text-sm text-foreground">{r.body}</p>
            </li>
          ))}
        </ul>
      )}

      {/* Reply form */}
      {ticket.status === "open" && (
        <form
          onSubmit={(e) => { e.preventDefault(); if (reply.trim()) replyM.mutate(); }}
          className="mt-6 rounded-xl border border-border bg-bg-elevated p-5"
        >
          <label className="mb-2 block text-xs font-semibold uppercase tracking-wide text-text-muted">
            Add a reply
          </label>
          <textarea
            value={reply}
            onChange={(e) => setReply(e.target.value)}
            rows={4}
            placeholder="Write a reply..."
            className="w-full rounded-lg border border-border bg-bg-secondary px-3 py-2.5 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
          <div className="mt-3 flex justify-end">
            <button
              type="submit"
              disabled={replyM.isPending || !reply.trim()}
              className="rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
            >
              {replyM.isPending ? "Sending..." : "Send reply"}
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
`
}
