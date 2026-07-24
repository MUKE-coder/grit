package plugin

import "strings"

// webhooksWireGo emits internal/handlers/webhooks_wire.go — the thin bridge
// between the app and the grit-webhooks module. It lives in the handlers
// package so routes.go can wire it without a new import (routes.go already
// imports handlers), and it holds the service as a package singleton so any
// handler can fire an event with handlers.DispatchWebhook(...).
func webhooksWireGo(ctx Context) string {
	src := `package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	webhooks "github.com/MUKE-coder/grit-plugins/grit-webhooks"
	"gorm.io/gorm"
)

// webhookService is process-wide so DispatchWebhook works from any handler
// without threading the service through every constructor.
var webhookService *webhooks.Service

// InitWebhooks migrates the webhook tables and builds the service. Call once at
// startup (the plugin injects this into routes.go).
func InitWebhooks(db *gorm.DB) *webhooks.Service {
	// subscriptions + delivery_attempts. Safe to run every boot.
	_ = webhooks.AutoMigrate(db)

	webhookService = webhooks.NewService(db, webhooks.Config{
		MaxRetries:       8,
		RetryBackoffBase: 5 * time.Second, // 5s, 10s, 20s, ... with jitter
		TimeoutSeconds:   10,
		WorkerCount:      4,
	})
	return webhookService
}

// RegisterWebhookRoutes mounts the subscription CRUD + delivery-log endpoints
// under the given group. Wrapped here so routes.go needs no webhooks import.
func RegisterWebhookRoutes(rg *gin.RouterGroup, svc *webhooks.Service) {
	webhooks.RegisterRoutes(rg, svc)
}

// DispatchWebhook fans a domain event out to every subscription listening for
// it. A no-op if InitWebhooks hasn't run, so calling it can never panic a
// request — a webhook is a side effect, not the point of the handler.
//
//	handlers.DispatchWebhook("invoice.paid", gin.H{"id": inv.ID, "amount": 4900})
func DispatchWebhook(eventType string, payload interface{}) error {
	if webhookService == nil {
		return nil
	}
	return webhookService.Dispatch(eventType, payload)
}
`
	return strings.ReplaceAll(src, "{{MODULE}}", ctx.Module)
}

// webhooksAdminPage returns the Webhooks admin page: manage subscriptions, watch
// the per-subscription delivery log, send a test, and resend a failed delivery.
// Template literals are written as string concatenation so the whole thing fits
// in a Go raw string without needing backtick escaping.
func webhooksAdminPage() string {
	return `"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { useMe } from "@/hooks/use-auth";
import { Plus, Trash2, RefreshCw, Play, Webhook } from "@/lib/icons";

interface Subscription {
  id: number;
  url: string;
  events: string[];
  secret: string;
  active: boolean;
  description?: string;
}

interface Delivery {
  id: number;
  event_type: string;
  status: string;
  response_code: number;
  error: string;
  attempt_number: number;
  created_at: string;
}

export default function WebhooksPage() {
  const { data: me } = useMe();
  // useMe may return the user directly or wrapped in { data }.
  const meAny = me as unknown as { id?: string; data?: { id?: string } } | undefined;
  const userId = meAny?.data?.id ?? meAny?.id;
  const qc = useQueryClient();
  const [url, setUrl] = useState("");
  const [events, setEvents] = useState("");
  const [openLog, setOpenLog] = useState<number | null>(null);
  // The signing secret is returned exactly once, on creation. Hold it so we can
  // show it in a one-time banner — it can't be retrieved again.
  const [newSecret, setNewSecret] = useState<string | null>(null);

  const subsQ = useQuery({
    queryKey: ["webhooks", userId],
    enabled: !!userId,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/webhooks?user_id=" + userId);
      return (data.subscriptions ?? []) as Subscription[];
    },
  });

  const createM = useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post("/api/webhooks", {
        user_id: userId,
        url,
        events: events.split(",").map((e) => e.trim()).filter(Boolean),
      });
      return data as { secret?: string };
    },
    onSuccess: (data) => {
      setUrl("");
      setEvents("");
      if (data?.secret) setNewSecret(data.secret);
      qc.invalidateQueries({ queryKey: ["webhooks"] });
    },
  });

  const deleteM = useMutation({
    mutationFn: async (id: number) => {
      await apiClient.delete("/api/webhooks/" + id);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["webhooks"] }),
  });

  const testM = useMutation({
    mutationFn: async (id: number) => {
      await apiClient.post("/api/webhooks/" + id + "/test");
    },
    onSuccess: (_d, id) => qc.invalidateQueries({ queryKey: ["deliveries", id] }),
  });

  const subs = subsQ.data ?? [];

  return (
    <div>
      <PageHeader title="Webhooks" subtitle="Outbound event subscriptions and delivery log" />

      <div className="rounded-xl border border-border bg-bg-elevated p-5 mb-6">
        <h2 className="mb-3 flex items-center gap-2 text-sm font-semibold text-foreground">
          <Webhook className="h-4 w-4 text-accent" /> New subscription
        </h2>
        <div className="grid gap-3 sm:grid-cols-[2fr_2fr_auto]">
          <input
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://acme.com/hooks"
            className="rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-sm text-foreground outline-none focus:border-accent"
          />
          <input
            value={events}
            onChange={(e) => setEvents(e.target.value)}
            placeholder="invoice.paid, user.created"
            className="rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-sm text-foreground outline-none focus:border-accent"
          />
          <button
            onClick={() => createM.mutate()}
            disabled={!url || !events || createM.isPending}
            className="inline-flex items-center gap-1.5 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50"
          >
            <Plus className="h-4 w-4" /> Add
          </button>
        </div>
        <p className="mt-2 text-xs text-text-muted">
          Comma-separated event names. A signing secret (whsec_…) is generated automatically.
        </p>
      </div>

      {newSecret && (
        <div className="mb-6 rounded-xl border border-warning/40 bg-warning/10 p-4">
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0">
              <p className="text-sm font-semibold text-foreground">Signing secret — copy it now</p>
              <p className="mt-0.5 text-xs text-text-muted">
                Shown once. Configure it on your endpoint to verify deliveries. It can&apos;t be retrieved later.
              </p>
              <code className="mt-2 block break-all rounded-lg border border-border bg-bg-tertiary px-3 py-2 font-mono text-xs text-foreground">
                {newSecret}
              </code>
            </div>
            <div className="flex shrink-0 gap-1.5">
              <button
                onClick={() => navigator.clipboard?.writeText(newSecret)}
                className="rounded-lg border border-border px-2.5 py-1.5 text-xs text-text-secondary hover:bg-bg-hover"
              >
                Copy
              </button>
              <button
                onClick={() => setNewSecret(null)}
                className="rounded-lg p-1.5 text-text-muted hover:bg-bg-hover hover:text-foreground"
                aria-label="Dismiss"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      )}

      {subsQ.isLoading ? (
        <p className="p-4 text-sm text-text-muted">Loading…</p>
      ) : subs.length === 0 ? (
        <p className="rounded-xl border border-border bg-bg-elevated p-8 text-center text-sm text-text-muted">
          No webhook subscriptions yet.
        </p>
      ) : (
        <div className="space-y-3">
          {subs.map((s) => (
            <div key={s.id} className="rounded-xl border border-border bg-bg-elevated">
              <div className="flex flex-wrap items-center justify-between gap-3 p-4">
                <div className="min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="truncate font-mono text-sm text-foreground">{s.url}</span>
                    <span className={"rounded-full px-2 py-0.5 text-[10px] font-medium " + (s.active ? "bg-success/15 text-success" : "bg-bg-hover text-text-muted")}>
                      {s.active ? "active" : "paused"}
                    </span>
                  </div>
                  <div className="mt-1 flex flex-wrap gap-1">
                    {s.events.map((e) => (
                      <span key={e} className="rounded bg-bg-hover px-1.5 py-0.5 text-[11px] text-text-secondary">{e}</span>
                    ))}
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <button onClick={() => testM.mutate(s.id)} title="Send test" className="inline-flex items-center gap-1 rounded-lg border border-border px-2.5 py-1.5 text-xs text-text-secondary hover:bg-bg-hover">
                    <Play className="h-3.5 w-3.5" /> Test
                  </button>
                  <button onClick={() => setOpenLog(openLog === s.id ? null : s.id)} className="inline-flex items-center gap-1 rounded-lg border border-border px-2.5 py-1.5 text-xs text-text-secondary hover:bg-bg-hover">
                    Deliveries
                  </button>
                  <button onClick={() => { if (confirm("Delete this subscription?")) deleteM.mutate(s.id); }} title="Delete" className="rounded-lg p-1.5 text-text-muted hover:text-danger hover:bg-bg-hover">
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
              {openLog === s.id && <DeliveryLog subscriptionId={s.id} />}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function DeliveryLog({ subscriptionId }: { subscriptionId: number }) {
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["deliveries", subscriptionId],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/webhooks/" + subscriptionId + "/deliveries?page=1&page_size=20");
      return (data.deliveries ?? []) as Delivery[];
    },
  });

  const retryM = useMutation({
    mutationFn: async (deliveryId: number) => {
      await apiClient.post("/api/webhooks/" + subscriptionId + "/retry/" + deliveryId);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["deliveries", subscriptionId] }),
  });

  const rows = q.data ?? [];

  return (
    <div className="border-t border-border p-4">
      {q.isLoading ? (
        <p className="text-xs text-text-muted">Loading deliveries…</p>
      ) : rows.length === 0 ? (
        <p className="text-xs text-text-muted">No deliveries yet — send a test or fire an event.</p>
      ) : (
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left text-xs text-text-muted">
              <th className="py-1.5 font-medium">Event</th>
              <th className="py-1.5 font-medium">Status</th>
              <th className="py-1.5 font-medium">Code</th>
              <th className="py-1.5 font-medium">Attempt</th>
              <th className="py-1.5" />
            </tr>
          </thead>
          <tbody>
            {rows.map((d) => (
              <tr key={d.id} className="border-b border-border/50 last:border-0">
                <td className="py-1.5 font-mono text-xs text-foreground">{d.event_type}</td>
                <td className="py-1.5">
                  <span className={"rounded-full px-2 py-0.5 text-[10px] font-medium " + (d.status === "success" ? "bg-success/15 text-success" : d.status === "failed" ? "bg-danger/15 text-danger" : "bg-warning/15 text-warning")}>
                    {d.status}
                  </span>
                </td>
                <td className="py-1.5 text-xs text-text-secondary">{d.response_code || "—"}</td>
                <td className="py-1.5 text-xs text-text-secondary">#{d.attempt_number}</td>
                <td className="py-1.5 text-right">
                  {d.status === "failed" && (
                    <button onClick={() => retryM.mutate(d.id)} title="Resend" className="inline-flex items-center gap-1 rounded border border-border px-2 py-1 text-[11px] text-accent hover:bg-bg-hover">
                      <RefreshCw className="h-3 w-3" /> Resend
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
`
}

// webhooksTanStackRoute returns the TanStack route wrapper that renders the
// Webhooks page component.
func webhooksTanStackRoute() string {
	return `import { createFileRoute } from '@tanstack/react-router'
import WebhooksPage from '@/pages/system/webhooks'

export const Route = createFileRoute('/_dashboard/system/webhooks')({
  component: WebhooksPage,
})
`
}
