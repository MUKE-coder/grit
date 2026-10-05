"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { Check, ArrowLeft } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import type { Ticket as TicketRow, TicketReply } from "@repo/shared/types";

// The rows are the models', from grit sync. Narrowed here only where a Go
// struct cannot say more: the two string enums, and the relations the thread
// preloads, which sync types as unknown.
type Author = { first_name: string; last_name: string; email: string };

type Reply = Omit<TicketReply, "user"> & { user?: Author };

type Ticket = Omit<TicketRow, "status" | "priority" | "user" | "replies"> & {
  status: "open" | "closed";
  priority: "low" | "medium" | "high" | "critical";
  user?: Author;
  replies: Reply[];
};

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
      <div className="animate-pulse space-y-4">
        <div className="h-16 rounded-xl bg-bg-hover" />
        <div className="h-6 w-32 rounded bg-bg-hover" />
        <div className="rounded-xl border border-border bg-bg-elevated p-5 space-y-3">
          <div className="h-4 w-1/3 rounded bg-bg-hover" />
          <div className="h-4 w-full rounded bg-bg-hover" />
          <div className="h-4 w-5/6 rounded bg-bg-hover" />
          <div className="h-4 w-2/3 rounded bg-bg-hover" />
        </div>
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
            className={inputClasses({ multiline: true })}
          />
          <div className="mt-3 flex justify-end">
            <button
              type="submit"
              disabled={replyM.isPending || !reply.trim()}
              className={buttonClasses()}
            >
              {replyM.isPending ? "Sending..." : "Send reply"}
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
