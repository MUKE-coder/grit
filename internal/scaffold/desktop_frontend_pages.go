package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDesktopFrontendPageFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "frontend", "src", "routes", "_layout", "index.tsx"):             desktopDashboardRoute(),
		filepath.Join(root, "frontend", "src", "routes", "_layout", "blogs.index.tsx"):       desktopBlogListRoute(),
		filepath.Join(root, "frontend", "src", "routes", "_layout", "blogs.new.tsx"):         desktopBlogNewRoute(),
		filepath.Join(root, "frontend", "src", "routes", "_layout", "blogs.$id.edit.tsx"):    desktopBlogEditRoute(),
		filepath.Join(root, "frontend", "src", "routes", "_layout", "contacts.index.tsx"):    desktopContactListRoute(),
		filepath.Join(root, "frontend", "src", "routes", "_layout", "contacts.new.tsx"):      desktopContactNewRoute(),
		filepath.Join(root, "frontend", "src", "routes", "_layout", "contacts.$id.edit.tsx"): desktopContactEditRoute(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// ── Dashboard Route ─────────────────────────────────────────────────────────

func desktopDashboardRoute() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { FileText, Users, CheckCircle, Activity } from "lucide-react";
import { useAuth } from "../../hooks/use-auth";
// @ts-ignore
import { GetBlogs, GetContacts } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/")({
  component: DashboardPage,
});

function DashboardPage() {
  const { user } = useAuth();
  const { data: blogData } = useQuery({
    queryKey: ["blogs", "stats"],
    queryFn: () => GetBlogs(1, 1000, ""),
  });
  const { data: contactData } = useQuery({
    queryKey: ["contacts", "stats"],
    queryFn: () => GetContacts(1, 1000, ""),
  });

  const stats = [
    { label: "Total Blogs", value: blogData?.total || 0, icon: FileText, color: "text-accent" },
    { label: "Published", value: (blogData?.data || []).filter((b: any) => b.published).length, icon: CheckCircle, color: "text-success" },
    { label: "Total Contacts", value: contactData?.total || 0, icon: Users, color: "text-accent" },
    { label: "Recent Activity", value: (blogData?.total || 0) + (contactData?.total || 0), icon: Activity, color: "text-warning" },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-foreground mb-1">Dashboard</h1>
      <p className="text-text-secondary mb-6">{"Welcome back, " + (user?.name || "User")}</p>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="bg-bg-elevated border border-border rounded-xl p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm text-text-secondary">{stat.label}</span>
                <Icon size={18} className={stat.color} />
              </div>
              <div className="text-3xl font-bold text-foreground">{stat.value}</div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
`
}

// ── Blog List Route ─────────────────────────────────────────────────────────

func desktopBlogListRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Pencil, Trash2, FileDown, FileSpreadsheet, Search } from "lucide-react";
import { toast } from "sonner";
// @ts-ignore
import { GetBlogs, DeleteBlog, ExportBlogsPDF, ExportBlogsExcel } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/blogs/")({
  component: BlogListPage,
});

function BlogListPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const pageSize = 10;

  const { data, isLoading } = useQuery({
    queryKey: ["blogs", page, search],
    queryFn: () => GetBlogs(page, pageSize, search),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => DeleteBlog(id),
    onSuccess: () => {
      toast.success("Blog deleted");
      queryClient.invalidateQueries({ queryKey: ["blogs"] });
    },
    onError: (err: any) => toast.error(err?.message || "Failed to delete"),
  });

  const handleDelete = (id: number, title: string) => {
    if (window.confirm("Delete \"" + title + "\"? This cannot be undone.")) {
      deleteMutation.mutate(id);
    }
  };

  const handleExportPDF = async () => {
    try {
      const path = await ExportBlogsPDF();
      toast.success("PDF exported to " + path);
    } catch (err: any) {
      toast.error(err?.message || "Export failed");
    }
  };

  const handleExportExcel = async () => {
    try {
      const path = await ExportBlogsExcel();
      toast.success("Excel exported to " + path);
    } catch (err: any) {
      toast.error(err?.message || "Export failed");
    }
  };

  const blogs = data?.data || [];
  const total = data?.total || 0;
  const totalPages = data?.pages || 1;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Blogs</h1>
          <p className="text-sm text-text-secondary">{total + " total blog" + (total !== 1 ? "s" : "")}</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleExportPDF}
            className="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            <FileDown size={16} />
            PDF
          </button>
          <button
            onClick={handleExportExcel}
            className="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            <FileSpreadsheet size={16} />
            Excel
          </button>
          <button
            onClick={() => navigate({ to: "/blogs/new" })}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors"
          >
            <Plus size={16} />
            New Blog
          </button>
        </div>
      </div>

      <div className="mb-4">
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            type="text"
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            placeholder="Search blogs..."
            className="w-full pl-9 pr-3 py-2 bg-bg-elevated border border-border rounded-lg text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
          />
        </div>
      </div>

      <div className="bg-bg-elevated border border-border rounded-xl overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border">
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Title</th>
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Slug</th>
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Status</th>
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Created</th>
              <th className="text-right text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-text-muted">Loading...</td>
              </tr>
            ) : blogs.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-text-muted">No blogs found</td>
              </tr>
            ) : (
              blogs.map((blog: any) => (
                <tr key={blog.id} className="border-b border-border last:border-0 hover:bg-bg-hover/50 transition-colors">
                  <td className="px-4 py-3 text-sm text-foreground font-medium">{blog.title}</td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{blog.slug}</td>
                  <td className="px-4 py-3">
                    <span className={"inline-flex items-center px-2 py-0.5 rounded text-xs font-medium " + (blog.published ? "bg-success/10 text-success" : "bg-warning/10 text-warning")}>
                      {blog.published ? "Published" : "Draft"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-text-secondary">
                    {new Date(blog.created_at).toLocaleDateString("en-US", { year: "numeric", month: "short", day: "numeric" })}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => navigate({ to: "/blogs/$id/edit", params: { id: String(blog.id) } })}
                        className="p-1.5 rounded text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors"
                      >
                        <Pencil size={14} />
                      </button>
                      <button
                        onClick={() => handleDelete(blog.id, blog.title)}
                        className="p-1.5 rounded text-text-secondary hover:text-danger hover:bg-danger/10 transition-colors"
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <p className="text-sm text-text-muted">
            {"Page " + page + " of " + totalPages}
          </p>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="px-3 py-1.5 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Previous
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="px-3 py-1.5 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
`
}

// ── Blog New Route ──────────────────────────────────────────────────────────

func desktopBlogNewRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useAuth } from "../../hooks/use-auth";
// @ts-ignore
import { CreateBlog } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/blogs/new")({
  component: BlogNewPage,
});

function BlogNewPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [published, setPublished] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      toast.error("Title is required");
      return;
    }
    setLoading(true);
    try {
      await CreateBlog({
        title: title.trim(),
        content,
        published,
        author_id: user?.id || 1,
      });
      toast.success("Blog created");
      await queryClient.invalidateQueries({ queryKey: ["blogs"] });
      navigate({ to: "/blogs" });
    } catch (err: any) {
      toast.error(err?.message || "Failed to create blog");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-foreground mb-1">New Blog</h1>
      <p className="text-sm text-text-secondary mb-6">Create a new blog post</p>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="bg-bg-elevated border border-border rounded-xl p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Blog post title"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Content</label>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Write your blog content..."
              rows={10}
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors resize-y"
            />
          </div>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setPublished(!published)}
              className={"relative inline-flex h-6 w-11 items-center rounded-full transition-colors " + (published ? "bg-accent" : "bg-border")}
            >
              <span className={"inline-block h-4 w-4 rounded-full bg-white transition-transform " + (published ? "translate-x-6" : "translate-x-1")} />
            </button>
            <label className="text-sm text-foreground">
              {published ? "Published" : "Draft"}
            </label>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-2.5 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Creating..." : "Create Blog"}
          </button>
          <button
            type="button"
            onClick={() => navigate({ to: "/blogs" })}
            className="px-6 py-2.5 text-sm font-medium text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
`
}

// ── Blog Edit Route ─────────────────────────────────────────────────────────

func desktopBlogEditRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useAuth } from "../../hooks/use-auth";
// @ts-ignore
import { GetBlog, UpdateBlog } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/blogs/$id/edit")({
  component: BlogEditPage,
});

function BlogEditPage() {
  const { id } = Route.useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [published, setPublished] = useState(false);
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    GetBlog(Number(id))
      .then((blog: any) => {
        setTitle(blog.title);
        setContent(blog.content || "");
        setPublished(blog.published);
      })
      .catch((err: any) => {
        toast.error(err?.message || "Failed to load blog");
        navigate({ to: "/blogs" });
      })
      .finally(() => setFetching(false));
  }, [id, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      toast.error("Title is required");
      return;
    }
    setLoading(true);
    try {
      await UpdateBlog(Number(id), {
        title: title.trim(),
        content,
        published,
        author_id: user?.id || 1,
      });
      toast.success("Blog updated");
      navigate({ to: "/blogs" });
    } catch (err: any) {
      toast.error(err?.message || "Failed to update blog");
    } finally {
      setLoading(false);
    }
  };

  if (fetching) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="w-6 h-6 border-2 border-accent border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-foreground mb-1">Edit Blog</h1>
      <p className="text-sm text-text-secondary mb-6">Update the blog post details below</p>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="bg-bg-elevated border border-border rounded-xl p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Blog post title"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Content</label>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Write your blog content..."
              rows={10}
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors resize-y"
            />
          </div>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setPublished(!published)}
              className={"relative inline-flex h-6 w-11 items-center rounded-full transition-colors " + (published ? "bg-accent" : "bg-border")}
            >
              <span className={"inline-block h-4 w-4 rounded-full bg-white transition-transform " + (published ? "translate-x-6" : "translate-x-1")} />
            </button>
            <label className="text-sm text-foreground">
              {published ? "Published" : "Draft"}
            </label>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-2.5 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Saving..." : "Update Blog"}
          </button>
          <button
            type="button"
            onClick={() => navigate({ to: "/blogs" })}
            className="px-6 py-2.5 text-sm font-medium text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
`
}

// ── Contact List Route ──────────────────────────────────────────────────────

func desktopContactListRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Pencil, Trash2, FileDown, FileSpreadsheet, Search } from "lucide-react";
import { toast } from "sonner";
// @ts-ignore
import { GetContacts, DeleteContact, ExportContactsPDF, ExportContactsExcel } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/contacts/")({
  component: ContactListPage,
});

function ContactListPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const pageSize = 10;

  const { data, isLoading } = useQuery({
    queryKey: ["contacts", page, search],
    queryFn: () => GetContacts(page, pageSize, search),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => DeleteContact(id),
    onSuccess: () => {
      toast.success("Contact deleted");
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
    onError: (err: any) => toast.error(err?.message || "Failed to delete"),
  });

  const handleDelete = (id: number, name: string) => {
    if (window.confirm("Delete \"" + name + "\"? This cannot be undone.")) {
      deleteMutation.mutate(id);
    }
  };

  const handleExportPDF = async () => {
    try {
      const path = await ExportContactsPDF();
      toast.success("PDF exported to " + path);
    } catch (err: any) {
      toast.error(err?.message || "Export failed");
    }
  };

  const handleExportExcel = async () => {
    try {
      const path = await ExportContactsExcel();
      toast.success("Excel exported to " + path);
    } catch (err: any) {
      toast.error(err?.message || "Export failed");
    }
  };

  const contacts = data?.data || [];
  const total = data?.total || 0;
  const totalPages = data?.pages || 1;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Contacts</h1>
          <p className="text-sm text-text-secondary">{total + " total contact" + (total !== 1 ? "s" : "")}</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleExportPDF}
            className="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            <FileDown size={16} />
            PDF
          </button>
          <button
            onClick={handleExportExcel}
            className="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            <FileSpreadsheet size={16} />
            Excel
          </button>
          <button
            onClick={() => navigate({ to: "/contacts/new" })}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors"
          >
            <Plus size={16} />
            New Contact
          </button>
        </div>
      </div>

      <div className="mb-4">
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            type="text"
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            placeholder="Search contacts..."
            className="w-full pl-9 pr-3 py-2 bg-bg-elevated border border-border rounded-lg text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
          />
        </div>
      </div>

      <div className="bg-bg-elevated border border-border rounded-xl overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border">
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Name</th>
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Email</th>
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Phone</th>
              <th className="text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Company</th>
              <th className="text-right text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-text-muted">Loading...</td>
              </tr>
            ) : contacts.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-text-muted">No contacts found</td>
              </tr>
            ) : (
              contacts.map((contact: any) => (
                <tr key={contact.id} className="border-b border-border last:border-0 hover:bg-bg-hover/50 transition-colors">
                  <td className="px-4 py-3 text-sm text-foreground font-medium">{contact.name}</td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{contact.email || "-"}</td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{contact.phone || "-"}</td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{contact.company || "-"}</td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => navigate({ to: "/contacts/$id/edit", params: { id: String(contact.id) } })}
                        className="p-1.5 rounded text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors"
                      >
                        <Pencil size={14} />
                      </button>
                      <button
                        onClick={() => handleDelete(contact.id, contact.name)}
                        className="p-1.5 rounded text-text-secondary hover:text-danger hover:bg-danger/10 transition-colors"
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <p className="text-sm text-text-muted">
            {"Page " + page + " of " + totalPages}
          </p>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="px-3 py-1.5 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Previous
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="px-3 py-1.5 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
`
}

// ── Contact New Route ───────────────────────────────────────────────────────

func desktopContactNewRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
// @ts-ignore
import { CreateContact } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/contacts/new")({
  component: ContactNewPage,
});

function ContactNewPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [company, setCompany] = useState("");
  const [notes, setNotes] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error("Name is required");
      return;
    }
    setLoading(true);
    try {
      await CreateContact({
        name: name.trim(),
        email: email.trim(),
        phone: phone.trim(),
        company: company.trim(),
        notes,
      });
      toast.success("Contact created");
      await queryClient.invalidateQueries({ queryKey: ["contacts"] });
      navigate({ to: "/contacts" });
    } catch (err: any) {
      toast.error(err?.message || "Failed to create contact");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-foreground mb-1">New Contact</h1>
      <p className="text-sm text-text-secondary mb-6">Add a new contact to your directory</p>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="bg-bg-elevated border border-border rounded-xl p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contact name"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-foreground mb-1.5">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="email@example.com"
                className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-foreground mb-1.5">Phone</label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="+1 (555) 000-0000"
                className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Company</label>
            <input
              type="text"
              value={company}
              onChange={(e) => setCompany(e.target.value)}
              placeholder="Company name"
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Notes</label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Additional notes..."
              rows={4}
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors resize-y"
            />
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-2.5 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Creating..." : "Create Contact"}
          </button>
          <button
            type="button"
            onClick={() => navigate({ to: "/contacts" })}
            className="px-6 py-2.5 text-sm font-medium text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
`
}

// ── Contact Edit Route ──────────────────────────────────────────────────────

func desktopContactEditRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { toast } from "sonner";
// @ts-ignore
import { GetContact, UpdateContact } from "../../../wailsjs/go/main/App";

export const Route = createFileRoute("/_layout/contacts/$id/edit")({
  component: ContactEditPage,
});

function ContactEditPage() {
  const { id } = Route.useParams();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [company, setCompany] = useState("");
  const [notes, setNotes] = useState("");
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    GetContact(Number(id))
      .then((contact: any) => {
        setName(contact.name);
        setEmail(contact.email || "");
        setPhone(contact.phone || "");
        setCompany(contact.company || "");
        setNotes(contact.notes || "");
      })
      .catch((err: any) => {
        toast.error(err?.message || "Failed to load contact");
        navigate({ to: "/contacts" });
      })
      .finally(() => setFetching(false));
  }, [id, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error("Name is required");
      return;
    }
    setLoading(true);
    try {
      await UpdateContact(Number(id), {
        name: name.trim(),
        email: email.trim(),
        phone: phone.trim(),
        company: company.trim(),
        notes,
      });
      toast.success("Contact updated");
      navigate({ to: "/contacts" });
    } catch (err: any) {
      toast.error(err?.message || "Failed to update contact");
    } finally {
      setLoading(false);
    }
  };

  if (fetching) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="w-6 h-6 border-2 border-accent border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-foreground mb-1">Edit Contact</h1>
      <p className="text-sm text-text-secondary mb-6">Update the contact details below</p>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="bg-bg-elevated border border-border rounded-xl p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contact name"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-foreground mb-1.5">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="email@example.com"
                className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-foreground mb-1.5">Phone</label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="+1 (555) 000-0000"
                className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Company</label>
            <input
              type="text"
              value={company}
              onChange={(e) => setCompany(e.target.value)}
              placeholder="Company name"
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Notes</label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Additional notes..."
              rows={4}
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors resize-y"
            />
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-2.5 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Saving..." : "Update Contact"}
          </button>
          <button
            type="button"
            onClick={() => navigate({ to: "/contacts" })}
            className="px-6 py-2.5 text-sm font-medium text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
`
}
