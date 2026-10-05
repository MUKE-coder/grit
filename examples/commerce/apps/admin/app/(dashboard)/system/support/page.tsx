"use client";

import { useState } from "react";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import { IconButton } from "@/components/ui/IconButton";
import { Plus, MessageSquare, AlertCircle } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import type { Ticket as TicketRow } from "@repo/shared/types";

// The row is the model's, from grit sync. Narrowed here only where a Go struct
// cannot say more: the two string enums, and the user the list preloads.
type Ticket = Omit<TicketRow, "status" | "priority" | "user"> & {
  status: "open" | "closed";
  priority: "low" | "medium" | "high" | "critical";
  user?: { first_name: string; last_name: string; email: string };
};

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
              (status === s ? "bg-accent text-accent-fg" : "text-text-secondary hover:text-foreground")
            }
          >
            {s === "open" ? <MessageSquare className="h-3.5 w-3.5" /> : <AlertCircle className="h-3.5 w-3.5" />}
            {s}
          </button>
        ))}
      </div>

      {isLoading ? (
        <ul className="space-y-2 animate-pulse">
          {Array.from({ length: 4 }).map((_, i) => (
            <li key={i} className="flex gap-4 rounded-xl border border-border bg-bg-elevated p-4">
              <div className="h-16 w-32 shrink-0 rounded-lg bg-bg-hover" />
              <div className="flex-1 space-y-2">
                <div className="h-4 w-1/3 rounded bg-bg-hover" />
                <div className="h-3 w-3/4 rounded bg-bg-hover" />
                <div className="h-3 w-1/4 rounded bg-bg-hover" />
              </div>
            </li>
          ))}
        </ul>
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
              className={buttonClasses({ className: "mt-4" })}
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
            className={buttonClasses()}
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
            className={inputClasses()}
          />
        </Field>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <Field label="Priority">
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value as Ticket["priority"])}
              className={inputClasses()}
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
              className={inputClasses()}
            />
          </Field>
        </div>

        <Field label="Describe what's happening" required>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={5}
            placeholder="Steps to reproduce, expected vs. actual behaviour, anything that helps."
            className={inputClasses({ multiline: true })}
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
