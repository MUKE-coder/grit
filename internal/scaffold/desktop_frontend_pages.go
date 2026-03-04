package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDesktopFrontendPageFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "frontend", "src", "pages", "login.tsx"):          desktopLoginPage(),
		filepath.Join(root, "frontend", "src", "pages", "register.tsx"):       desktopRegisterPage(),
		filepath.Join(root, "frontend", "src", "pages", "dashboard.tsx"):      desktopDashboardPage(),
		filepath.Join(root, "frontend", "src", "pages", "blogs", "index.tsx"): desktopBlogListPage(),
		filepath.Join(root, "frontend", "src", "pages", "blogs", "form.tsx"):  desktopBlogFormPage(),
		filepath.Join(root, "frontend", "src", "pages", "contacts", "index.tsx"): desktopContactListPage(),
		filepath.Join(root, "frontend", "src", "pages", "contacts", "form.tsx"):  desktopContactFormPage(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// ── Login Page ───────────────────────────────────────────────────────────────

func desktopLoginPage() string {
	return `import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { toast } from "sonner";
import { useAuth } from "../hooks/use-auth";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login(email, password);
      navigate("/");
    } catch (err: any) {
      toast.error(err?.message || "Invalid credentials");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-background">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-foreground">Welcome Back</h1>
          <p className="text-sm text-text-secondary mt-1">Sign in to your account</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 bg-bg-elevated border border-border rounded-xl p-6">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="admin@example.com"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter password"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-accent hover:bg-accent/90 text-white font-medium py-2.5 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>
        <p className="text-center text-sm text-text-secondary mt-4">
          {"Don't have an account? "}
          <Link to="/register" className="text-accent hover:underline">Create one</Link>
        </p>
        <p className="text-center text-xs text-text-muted mt-2">Default: admin@example.com / admin123</p>
      </div>
    </div>
  );
}
`
}

// ── Register Page ────────────────────────────────────────────────────────────

func desktopRegisterPage() string {
	return `import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { toast } from "sonner";
import { useAuth } from "../hooks/use-auth";

export default function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }
    if (password.length < 6) {
      toast.error("Password must be at least 6 characters");
      return;
    }
    setLoading(true);
    try {
      await register(name, email, password);
      navigate("/");
    } catch (err: any) {
      toast.error(err?.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-background">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-foreground">Create Account</h1>
          <p className="text-sm text-text-secondary mt-1">Get started with your new account</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 bg-bg-elevated border border-border rounded-xl p-6">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your name"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="At least 6 characters"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Confirm Password</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Repeat password"
              required
              className="w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-accent hover:bg-accent/90 text-white font-medium py-2.5 rounded-lg transition-colors disabled:opacity-50"
          >
            {loading ? "Creating account..." : "Create Account"}
          </button>
        </form>
        <p className="text-center text-sm text-text-secondary mt-4">
          {"Already have an account? "}
          <Link to="/login" className="text-accent hover:underline">Sign in</Link>
        </p>
      </div>
    </div>
  );
}
`
}

// ── Dashboard Page ───────────────────────────────────────────────────────────

func desktopDashboardPage() string {
	return `import { useQuery } from "@tanstack/react-query";
import { FileText, Users, CheckCircle, Activity } from "lucide-react";
import { useAuth } from "../hooks/use-auth";
// @ts-ignore
import { GetBlogs, GetContacts } from "../../wailsjs/go/main/App";

export default function DashboardPage() {
  const { user } = useAuth();
  const { data: blogData } = useQuery({
    queryKey: ["blogs-stats"],
    queryFn: () => GetBlogs(1, 1000, ""),
  });
  const { data: contactData } = useQuery({
    queryKey: ["contacts-stats"],
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

// ── Blog List Page ───────────────────────────────────────────────────────────

func desktopBlogListPage() string {
	return `import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Pencil, Trash2, FileDown, FileSpreadsheet, Search } from "lucide-react";
import { toast } from "sonner";
// @ts-ignore
import { GetBlogs, DeleteBlog, ExportBlogsPDF, ExportBlogsExcel } from "../../../wailsjs/go/main/App";

export default function BlogListPage() {
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
            onClick={() => navigate("/blogs/new")}
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
                        onClick={() => navigate("/blogs/" + blog.id + "/edit")}
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

// ── Blog Form Page ───────────────────────────────────────────────────────────

func desktopBlogFormPage() string {
	return `import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import { useAuth } from "../../hooks/use-auth";
// @ts-ignore
import { GetBlog, CreateBlog, UpdateBlog } from "../../../wailsjs/go/main/App";

export default function BlogFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const isEditing = Boolean(id);

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [published, setPublished] = useState(false);
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(false);

  useEffect(() => {
    if (isEditing && id) {
      setFetching(true);
      GetBlog(Number(id))
        .then((blog: any) => {
          setTitle(blog.title);
          setContent(blog.content || "");
          setPublished(blog.published);
        })
        .catch((err: any) => {
          toast.error(err?.message || "Failed to load blog");
          navigate("/blogs");
        })
        .finally(() => setFetching(false));
    }
  }, [id, isEditing, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      toast.error("Title is required");
      return;
    }
    setLoading(true);
    try {
      const input = {
        title: title.trim(),
        content,
        published,
        author_id: user?.id || 1,
      };
      if (isEditing && id) {
        await UpdateBlog(Number(id), input);
        toast.success("Blog updated");
      } else {
        await CreateBlog(input);
        toast.success("Blog created");
      }
      navigate("/blogs");
    } catch (err: any) {
      toast.error(err?.message || "Failed to save blog");
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
      <h1 className="text-2xl font-bold text-foreground mb-1">
        {isEditing ? "Edit Blog" : "New Blog"}
      </h1>
      <p className="text-sm text-text-secondary mb-6">
        {isEditing ? "Update the blog post details below" : "Create a new blog post"}
      </p>

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
            {loading ? "Saving..." : (isEditing ? "Update Blog" : "Create Blog")}
          </button>
          <button
            type="button"
            onClick={() => navigate("/blogs")}
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

// ── Contact List Page ────────────────────────────────────────────────────────

func desktopContactListPage() string {
	return `import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Pencil, Trash2, FileDown, FileSpreadsheet, Search } from "lucide-react";
import { toast } from "sonner";
// @ts-ignore
import { GetContacts, DeleteContact, ExportContactsPDF, ExportContactsExcel } from "../../../wailsjs/go/main/App";

export default function ContactListPage() {
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
            onClick={() => navigate("/contacts/new")}
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
                        onClick={() => navigate("/contacts/" + contact.id + "/edit")}
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

// ── Contact Form Page ────────────────────────────────────────────────────────

func desktopContactFormPage() string {
	return `import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
// @ts-ignore
import { GetContact, CreateContact, UpdateContact } from "../../../wailsjs/go/main/App";

export default function ContactFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isEditing = Boolean(id);

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [company, setCompany] = useState("");
  const [notes, setNotes] = useState("");
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(false);

  useEffect(() => {
    if (isEditing && id) {
      setFetching(true);
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
          navigate("/contacts");
        })
        .finally(() => setFetching(false));
    }
  }, [id, isEditing, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error("Name is required");
      return;
    }
    setLoading(true);
    try {
      const input = {
        name: name.trim(),
        email: email.trim(),
        phone: phone.trim(),
        company: company.trim(),
        notes,
      };
      if (isEditing && id) {
        await UpdateContact(Number(id), input);
        toast.success("Contact updated");
      } else {
        await CreateContact(input);
        toast.success("Contact created");
      }
      navigate("/contacts");
    } catch (err: any) {
      toast.error(err?.message || "Failed to save contact");
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
      <h1 className="text-2xl font-bold text-foreground mb-1">
        {isEditing ? "Edit Contact" : "New Contact"}
      </h1>
      <p className="text-sm text-text-secondary mb-6">
        {isEditing ? "Update the contact details below" : "Add a new contact to your directory"}
      </p>

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
            {loading ? "Saving..." : (isEditing ? "Update Contact" : "Create Contact")}
          </button>
          <button
            type="button"
            onClick={() => navigate("/contacts")}
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
