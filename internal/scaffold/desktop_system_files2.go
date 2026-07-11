package scaffold

// Second half of the desktop system pages: Blogs (list + editor), Support
// (list + thread), and Dashboard settings. Same offline-graceful pattern as
// desktop_system_files.go.

// desktopClientSystemBlogsPage lists blog posts and creates new ones.
func desktopClientSystemBlogsPage() string {
	return `import { useState } from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { FileText, Plus } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { EmptyState, relTime } from "@/components/system-ui";
import { ResourceDrawer } from "@/components/resource-drawer";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/blogs")({
  component: SystemBlogsPage,
});

interface Blog {
  id: string; title: string; slug?: string; excerpt?: string;
  image?: string; published?: boolean; created_at: string;
}

function SystemBlogsPage() {
  const navigate = useNavigate();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [excerpt, setExcerpt] = useState("");

  const { data = [], isLoading } = useQuery<Blog[]>({
    queryKey: ["system", "blogs"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Blog[] }>("/admin/blogs?page_size=100");
        return data.data ?? [];
      } catch { return []; }
    },
    refetchInterval: 60_000,
  });

  const create = useMutation({
    mutationFn: (input: { title: string; excerpt: string }) => apiClient.post<{ data: Blog }>("/admin/blogs", input),
    onSuccess: (res) => {
      qc.invalidateQueries({ queryKey: ["system", "blogs"] });
      setOpen(false); setTitle(""); setExcerpt("");
      const id = res.data?.data?.id;
      if (id) navigate({ to: "/app/system/blogs/$id", params: { id: String(id) } });
    },
  });

  return (
    <div>
      <PageHeader
        title="Blogs"
        description="Write & publish posts"
        actions={
          <button
            onClick={() => setOpen(true)}
            className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover"
          >
            <Plus className="h-4 w-4" /> New Blog
          </button>
        }
      />

      <div className="mt-6 space-y-3">
        {isLoading ? (
          <div className="rounded-xl border border-border bg-surface px-5 py-12 text-center text-[13px] text-foreground-muted">Loading…</div>
        ) : data.length === 0 ? (
          <div className="rounded-xl border border-border bg-surface">
            <EmptyState icon={FileText} title="No posts yet." hint="Create your first blog post." />
          </div>
        ) : (
          data.map((b) => (
            <button
              key={b.id}
              onClick={() => navigate({ to: "/app/system/blogs/$id", params: { id: String(b.id) } })}
              className="flex w-full items-center gap-4 rounded-xl border border-border bg-surface p-4 text-left transition-colors hover:bg-surface-hover"
            >
              <span className="flex h-14 w-14 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-surface-2">
                {b.image ? <img src={b.image} alt="" className="h-full w-full object-cover" /> : <FileText className="h-5 w-5 text-foreground-muted" />}
              </span>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2">
                  <p className="truncate text-[14px] font-semibold text-foreground">{b.title}</p>
                  <span className={"rounded-full px-2 py-0.5 text-[11px] font-medium " + (b.published ? "bg-success/10 text-success" : "bg-surface-2 text-foreground-muted")}>
                    {b.published ? "Published" : "Draft"}
                  </span>
                </div>
                {b.excerpt && <p className="line-clamp-1 text-[12px] text-foreground-secondary">{b.excerpt}</p>}
                <p className="mt-0.5 text-[11px] text-foreground-muted">{relTime(b.created_at)}</p>
              </div>
            </button>
          ))
        )}
      </div>

      <ResourceDrawer open={open} title="New Blog" description="Create a new post" onClose={() => setOpen(false)}>
        <form
          onSubmit={(e) => { e.preventDefault(); create.mutate({ title, excerpt }); }}
          className="flex h-full flex-col"
        >
          <div className="flex-1 space-y-4">
            <div>
              <label className="mb-1.5 block text-[13px] font-medium text-foreground">Title</label>
              <input value={title} onChange={(e) => setTitle(e.target.value)} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
            </div>
            <div>
              <label className="mb-1.5 block text-[13px] font-medium text-foreground">Excerpt</label>
              <textarea value={excerpt} onChange={(e) => setExcerpt(e.target.value)} rows={4} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
            </div>
          </div>
          <div className="mt-6 flex justify-end gap-3 border-t border-border pt-4">
            <button type="button" onClick={() => setOpen(false)} className="rounded-lg border border-border px-4 py-2 text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover">Cancel</button>
            <button type="submit" disabled={create.isPending || !title} className="rounded-lg bg-accent px-4 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover disabled:opacity-50">
              {create.isPending ? "Creating…" : "Create Blog"}
            </button>
          </div>
        </form>
      </ResourceDrawer>
    </div>
  );
}
`
}

// desktopClientSystemBlogDetailPage is the blog editor.
func desktopClientSystemBlogDetailPage() string {
	return `import { useEffect, useState } from "react";
import { createFileRoute, useNavigate, useParams } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Save, Trash2 } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/blogs/$id")({
  component: SystemBlogDetailPage,
});

interface Blog {
  id: string; title: string; excerpt?: string; content?: string; published?: boolean; created_at: string;
}

function SystemBlogDetailPage() {
  const { id } = useParams({ from: "/app/system/blogs/$id" });
  const navigate = useNavigate();
  const qc = useQueryClient();
  const [title, setTitle] = useState("");
  const [excerpt, setExcerpt] = useState("");
  const [content, setContent] = useState("");
  const [published, setPublished] = useState(false);

  const { data, isLoading } = useQuery<Blog | null>({
    queryKey: ["system", "blog", id],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Blog }>("/admin/blogs/" + id);
        return data.data;
      } catch { return null; }
    },
  });

  useEffect(() => {
    if (!data) return;
    setTitle(data.title ?? "");
    setExcerpt(data.excerpt ?? "");
    setContent(data.content ?? "");
    setPublished(!!data.published);
  }, [data]);

  const save = useMutation({
    mutationFn: () => apiClient.put("/admin/blogs/" + id, { title, excerpt, content, published }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["system", "blog", id] }),
  });
  const remove = useMutation({
    mutationFn: () => apiClient.delete("/admin/blogs/" + id),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["system", "blogs"] }); navigate({ to: "/app/system/blogs" }); },
  });

  if (isLoading) return <div className="text-[13px] text-foreground-muted">Loading…</div>;

  return (
    <div>
      <button onClick={() => navigate({ to: "/app/system/blogs" })} className="mb-3 flex items-center gap-1.5 text-[13px] text-foreground-secondary hover:text-foreground">
        <ArrowLeft className="h-4 w-4" /> Back to blogs
      </button>
      <PageHeader
        title="Edit post"
        description="Update this blog post"
        actions={
          <div className="flex items-center gap-2">
            <button onClick={() => { if (confirm("Delete this post?")) remove.mutate(); }} className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-[13px] text-danger hover:bg-danger/10">
              <Trash2 className="h-4 w-4" /> Delete
            </button>
            <button onClick={() => save.mutate()} disabled={save.isPending} className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover disabled:opacity-50">
              <Save className="h-4 w-4" /> {save.isPending ? "Saving…" : "Save"}
            </button>
          </div>
        }
      />

      <div className="mt-6 max-w-3xl space-y-4">
        <div>
          <label className="mb-1.5 block text-[13px] font-medium text-foreground">Title</label>
          <input value={title} onChange={(e) => setTitle(e.target.value)} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
        </div>
        <div>
          <label className="mb-1.5 block text-[13px] font-medium text-foreground">Excerpt</label>
          <textarea value={excerpt} onChange={(e) => setExcerpt(e.target.value)} rows={2} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
        </div>
        <div>
          <label className="mb-1.5 block text-[13px] font-medium text-foreground">Content</label>
          <textarea value={content} onChange={(e) => setContent(e.target.value)} rows={12} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 font-mono text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
        </div>
        <label className="flex items-center gap-2 text-[13px] text-foreground">
          <input type="checkbox" checked={published} onChange={(e) => setPublished(e.target.checked)} className="h-4 w-4 accent-accent" />
          Published
        </label>
      </div>
    </div>
  );
}
`
}

// desktopClientSystemSupportPage lists support tickets.
func desktopClientSystemSupportPage() string {
	return `import { useState } from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { MessageSquare, Plus } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { EmptyState, relTime } from "@/components/system-ui";
import { ResourceDrawer } from "@/components/resource-drawer";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/support")({
  component: SystemSupportPage,
});

interface Ticket {
  id: string; subject: string; description?: string; status: string;
  priority: "low" | "medium" | "high" | "critical"; created_at: string;
  user?: { first_name?: string; last_name?: string; email?: string };
}

const prio: Record<Ticket["priority"], string> = {
  low: "bg-surface-2 text-foreground-muted",
  medium: "bg-info/10 text-info",
  high: "bg-warning/10 text-warning",
  critical: "bg-danger/10 text-danger",
};

function SystemSupportPage() {
  const navigate = useNavigate();
  const qc = useQueryClient();
  const [status, setStatus] = useState<"open" | "closed">("open");
  const [open, setOpen] = useState(false);
  const [subject, setSubject] = useState("");
  const [priority, setPriority] = useState<Ticket["priority"]>("medium");
  const [description, setDescription] = useState("");

  const { data = [], isLoading } = useQuery<Ticket[]>({
    queryKey: ["system", "tickets", status],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Ticket[] }>("/tickets?status=" + status);
        return data.data ?? [];
      } catch { return []; }
    },
    refetchInterval: 30_000,
  });

  const create = useMutation({
    mutationFn: () => apiClient.post<{ data: Ticket }>("/tickets", { subject, priority, description }),
    onSuccess: (res) => {
      qc.invalidateQueries({ queryKey: ["system", "tickets"] });
      setOpen(false); setSubject(""); setDescription(""); setPriority("medium");
      const id = res.data?.data?.id;
      if (id) navigate({ to: "/app/system/support/$id", params: { id: String(id) } });
    },
  });

  return (
    <div>
      <PageHeader
        title="Support"
        description="Tickets & conversations"
        actions={
          <button onClick={() => setOpen(true)} className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover">
            <Plus className="h-4 w-4" /> New ticket
          </button>
        }
      />

      <div className="mt-6 flex gap-1">
        {(["open", "closed"] as const).map((s) => (
          <button
            key={s}
            onClick={() => setStatus(s)}
            className={"rounded-lg px-3 py-1.5 text-[13px] font-medium capitalize transition-colors " + (status === s ? "bg-accent/10 text-accent" : "text-foreground-secondary hover:bg-surface-hover")}
          >
            {s}
          </button>
        ))}
      </div>

      <div className="mt-4 space-y-3">
        {isLoading ? (
          <div className="rounded-xl border border-border bg-surface px-5 py-12 text-center text-[13px] text-foreground-muted">Loading…</div>
        ) : data.length === 0 ? (
          <div className="rounded-xl border border-border bg-surface">
            <EmptyState icon={MessageSquare} title={"No " + status + " tickets."} />
          </div>
        ) : (
          data.map((t) => (
            <button
              key={t.id}
              onClick={() => navigate({ to: "/app/system/support/$id", params: { id: String(t.id) } })}
              className="flex w-full items-center justify-between gap-4 rounded-xl border border-border bg-surface p-4 text-left transition-colors hover:bg-surface-hover"
            >
              <div className="min-w-0">
                <p className="truncate text-[14px] font-semibold text-foreground">{t.subject}</p>
                <p className="text-[12px] text-foreground-muted">
                  {t.user ? ([t.user.first_name, t.user.last_name].filter(Boolean).join(" ") || t.user.email) : "—"} · {relTime(t.created_at)}
                </p>
              </div>
              <span className={"shrink-0 rounded-full px-2 py-0.5 text-[11px] font-medium capitalize " + prio[t.priority]}>{t.priority}</span>
            </button>
          ))
        )}
      </div>

      <ResourceDrawer open={open} title="New ticket" description="Open a support ticket" onClose={() => setOpen(false)}>
        <form onSubmit={(e) => { e.preventDefault(); create.mutate(); }} className="flex h-full flex-col">
          <div className="flex-1 space-y-4">
            <div>
              <label className="mb-1.5 block text-[13px] font-medium text-foreground">Subject</label>
              <input value={subject} onChange={(e) => setSubject(e.target.value)} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
            </div>
            <div>
              <label className="mb-1.5 block text-[13px] font-medium text-foreground">Priority</label>
              <select value={priority} onChange={(e) => setPriority(e.target.value as Ticket["priority"])} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent">
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="critical">Critical</option>
              </select>
            </div>
            <div>
              <label className="mb-1.5 block text-[13px] font-medium text-foreground">Description</label>
              <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={5} className="w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent" />
            </div>
          </div>
          <div className="mt-6 flex justify-end gap-3 border-t border-border pt-4">
            <button type="button" onClick={() => setOpen(false)} className="rounded-lg border border-border px-4 py-2 text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover">Cancel</button>
            <button type="submit" disabled={create.isPending || !subject} className="rounded-lg bg-accent px-4 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover disabled:opacity-50">
              {create.isPending ? "Creating…" : "Create ticket"}
            </button>
          </div>
        </form>
      </ResourceDrawer>
    </div>
  );
}
`
}

// desktopClientSystemSupportThreadPage is the ticket conversation view.
func desktopClientSystemSupportThreadPage() string {
	return `import { useState } from "react";
import { createFileRoute, useNavigate, useParams } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Send } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { relTime } from "@/components/system-ui";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/support/$id")({
  component: SystemSupportThreadPage,
});

interface Reply { id: string; body: string; created_at: string; author?: string }
interface Ticket {
  id: string; subject: string; description?: string; status: string;
  priority: string; created_at: string; replies?: Reply[];
}

function SystemSupportThreadPage() {
  const { id } = useParams({ from: "/app/system/support/$id" });
  const navigate = useNavigate();
  const qc = useQueryClient();
  const [body, setBody] = useState("");

  const { data, isLoading } = useQuery<Ticket | null>({
    queryKey: ["system", "ticket", id],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Ticket }>("/tickets/" + id);
        return data.data;
      } catch { return null; }
    },
    refetchInterval: 15_000,
  });

  const reply = useMutation({
    mutationFn: () => apiClient.post("/tickets/" + id + "/reply", { body }),
    onSuccess: () => { setBody(""); qc.invalidateQueries({ queryKey: ["system", "ticket", id] }); },
  });
  const setStatus = useMutation({
    mutationFn: (action: "close" | "reopen") => apiClient.patch("/tickets/" + id + "/" + action),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["system", "ticket", id] }),
  });

  if (isLoading) return <div className="text-[13px] text-foreground-muted">Loading…</div>;
  if (!data) return <div className="text-[13px] text-foreground-muted">Ticket not found (or offline).</div>;

  const closed = data.status === "closed";
  const replies = data.replies ?? [];

  return (
    <div>
      <button onClick={() => navigate({ to: "/app/system/support" })} className="mb-3 flex items-center gap-1.5 text-[13px] text-foreground-secondary hover:text-foreground">
        <ArrowLeft className="h-4 w-4" /> Back to support
      </button>
      <PageHeader
        title={data.subject}
        description={"Priority: " + data.priority + " · " + data.status}
        actions={
          <button
            onClick={() => setStatus.mutate(closed ? "reopen" : "close")}
            className="rounded-lg border border-border px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover"
          >
            {closed ? "Reopen" : "Close ticket"}
          </button>
        }
      />

      <div className="mt-6 max-w-3xl space-y-3">
        <div className="rounded-xl border border-border bg-surface p-4">
          <p className="text-[13px] text-foreground">{data.description}</p>
          <p className="mt-2 text-[11px] text-foreground-muted">Opened {relTime(data.created_at)}</p>
        </div>
        {replies.map((r) => (
          <div key={r.id} className="rounded-xl border border-border bg-surface-2 p-4">
            <p className="text-[13px] text-foreground">{r.body}</p>
            <p className="mt-2 text-[11px] text-foreground-muted">{r.author ?? "Reply"} · {relTime(r.created_at)}</p>
          </div>
        ))}

        {!closed && (
          <form onSubmit={(e) => { e.preventDefault(); if (body.trim()) reply.mutate(); }} className="flex items-end gap-2">
            <textarea
              value={body}
              onChange={(e) => setBody(e.target.value)}
              rows={2}
              placeholder="Write a reply…"
              className="flex-1 rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground outline-none focus:border-accent focus:ring-1 focus:ring-accent"
            />
            <button type="submit" disabled={reply.isPending || !body.trim()} className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2.5 text-[13px] font-semibold text-white hover:bg-accent-hover disabled:opacity-50">
              <Send className="h-4 w-4" /> Send
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
`
}

// desktopClientSystemDashboardSettingsPage lets the user toggle which
// dashboard widgets show. Derived locally from the nav resources; the choice
// persists to localStorage (offline) and syncs to the server when online.
func desktopClientSystemDashboardSettingsPage() string {
	return `import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { Save } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { SectionCard } from "@/components/system-ui";
import { NAV_SECTIONS } from "@/lib/nav-config";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/dashboard-settings")({
  component: SystemDashboardSettingsPage,
});

const STORAGE_KEY = "grit-dashboard-widgets";

const SECTIONS = [
  { key: "cards", label: "Stat cards", items: ["Users", "Events (24h)", "Notifications", "Sync status"] },
  { key: "charts", label: "Charts", items: ["Activity (7 days)", "Severity mix"] },
  { key: "feeds", label: "Feeds", items: ["Recent activity", "Quick access"] },
];

function resourceWidgets(): string[] {
  const manage = NAV_SECTIONS.find((s) => s.title === "Manage");
  return (manage?.items ?? []).map((i) => i.label);
}

function SystemDashboardSettingsPage() {
  const [enabled, setEnabled] = useState<Record<string, boolean>>({});
  const [saved, setSaved] = useState(false);

  const allWidgets = [...SECTIONS.flatMap((s) => s.items), ...resourceWidgets()];

  useEffect(() => {
    let stored: Record<string, boolean> = {};
    try { stored = JSON.parse(localStorage.getItem(STORAGE_KEY) || "{}"); } catch { /* ignore */ }
    const init: Record<string, boolean> = {};
    for (const w of allWidgets) init[w] = stored[w] !== false;
    setEnabled(init);
    // Best-effort server load (ignored offline).
    apiClient.get("/dashboard-layout").catch(() => undefined);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const toggle = (w: string) => setEnabled((e) => ({ ...e, [w]: !e[w] }));

  const save = () => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(enabled));
    apiClient.put("/dashboard-layout", { widgets: enabled }).catch(() => undefined);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  const renderGroup = (label: string, items: string[]) => (
    <SectionCard key={label} title={label}>
      <div className="divide-y divide-border">
        {items.map((w) => (
          <label key={w} className="flex cursor-pointer items-center justify-between px-5 py-3 text-[13px] text-foreground hover:bg-surface-hover">
            {w}
            <input type="checkbox" checked={enabled[w] ?? true} onChange={() => toggle(w)} className="h-4 w-4 accent-accent" />
          </label>
        ))}
      </div>
    </SectionCard>
  );

  return (
    <div>
      <PageHeader
        title="Dashboard settings"
        description="Choose which widgets appear on your dashboard"
        actions={
          <button onClick={save} className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover">
            <Save className="h-4 w-4" /> {saved ? "Saved" : "Save"}
          </button>
        }
      />
      <div className="mt-6 space-y-4">
        {SECTIONS.map((s) => renderGroup(s.label, s.items))}
        {resourceWidgets().length > 0 && renderGroup("By resource", resourceWidgets())}
      </div>
    </div>
  );
}
`
}
