package scaffold

// The blog pages Grit wrote from v3.31.7 until contact-app review M42, when the
// Blog resource moved onto <ResourcePage>. Kept for the tests of the repairs
// that recognise them.

// legacyAdminBlogDetailPage generates the blog detail/edit page. Loads the
// blog by id, surfaces the cover + title + excerpt at the top, then
// renders the WordEditor over the content field. Autosaves on blur +
// has explicit Save / Publish / Delete actions.
func legacyAdminBlogDetailPage() string {
	return `"use client";

` + blogEditorNewReactImport + `
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import { apiClient, uploadFile } from "@/lib/api-client";
import { ArrowLeft, Save, Trash2, Upload, Check } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";

` + blogEditorDynamic + `interface Blog {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  image: string;
  published: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
}

interface ApiResponse<T> { data: T }

export default function BlogDetailPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();

  const [title, setTitle] = useState("");
  const [excerpt, setExcerpt] = useState("");
  const [content, setContent] = useState("");
  const [image, setImage] = useState("");
  const [coverUploading, setCoverUploading] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const coverInputRef = useRef<HTMLInputElement>(null);

  const { data: blog, isLoading } = useQuery<Blog>({
    queryKey: ["blog", params.id],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<Blog>>("/api/admin/blogs/" + params.id);
      return data.data;
    },
    enabled: !!params.id,
  });

  // Sync local form state when the blog loads. We track state locally
  // rather than threading useForm because the WordEditor is heavy + we
  // want autosave on blur, not on every keystroke.
  useEffect(() => {
    if (!blog) return;
    setTitle(blog.title || "");
    setExcerpt(blog.excerpt || "");
    setContent(blog.content || "");
    setImage(blog.image || "");
  }, [blog]);

  const save = useToastedMutation({
    mutationFn: async (patch: Partial<Blog>) => {
      const { data } = await apiClient.put<ApiResponse<Blog>>("/api/admin/blogs/" + params.id, patch);
      return data.data;
    },
    successMessage: "Saved",
    silentSuccess: true,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["blog", params.id] }),
  });

  const publish = useToastedMutation({
    mutationFn: async (next: boolean) => {
      const { data } = await apiClient.put<ApiResponse<Blog>>("/api/admin/blogs/" + params.id, { published: next });
      return data.data;
    },
    successMessage: (b) => b.published ? "Published" : "Moved back to draft",
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["blog", params.id] }),
  });

  const del = useToastedMutation({
    mutationFn: async () => apiClient.delete("/api/admin/blogs/" + params.id),
    successMessage: "Deleted",
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["blogs"] });
      router.push("/resources/blogs");
    },
  });

  const handleCoverUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setCoverUploading(true);
    try {
      const res = await uploadFile(file);
      const url = (res.data as Record<string, unknown>)?.url as string;
      if (url) {
        setImage(url);
        save.mutate({ image: url });
      }
    } finally {
      setCoverUploading(false);
      if (coverInputRef.current) coverInputRef.current.value = "";
    }
  };

  if (isLoading || !blog) {
    return <BlogDetailSkeleton />;
  }

  return (
    <div>
      <PageHeader
        title={title || "Untitled blog"}
        subtitle={blog.published ? "Published" : "Draft"}
        actions={
          <>
            <IconButton
              variant="secondary"
              icon={<ArrowLeft className="h-4 w-4" />}
              label="Back"
              onClick={() => router.push("/resources/blogs")}
            />
            <IconButton
              variant="secondary"
              icon={<Save className="h-4 w-4" />}
              label="Save"
              onClick={() => save.mutate({ title, excerpt, content, image })}
              disabled={save.isPending}
            />
            {blog.published ? (
              <IconButton
                variant="secondary"
                icon={<Check className="h-4 w-4" />}
                label="Unpublish"
                onClick={() => publish.mutate(false)}
                disabled={publish.isPending}
              />
            ) : (
              <IconButton
                icon={<Check className="h-4 w-4" />}
                label="Publish"
                onClick={() => publish.mutate(true)}
                disabled={publish.isPending}
              />
            )}
            <IconButton
              variant="danger"
              icon={<Trash2 className="h-4 w-4" />}
              label="Delete"
              onClick={() => setConfirmDelete(true)}
            />
          </>
        }
      />

      <ConfirmModal
        open={confirmDelete}
        title="Delete this blog?"
        description="This cannot be undone."
        confirmLabel="Delete"
        variant="danger"
        loading={del.isPending}
        onCancel={() => setConfirmDelete(false)}
        onConfirm={() => { setConfirmDelete(false); del.mutate(); }}
      />

      {/* Meta panel — cover, title, excerpt */}
      <section className="mb-6 grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div className="lg:col-span-1">
          <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-text-muted">Cover image</p>
          <div className="relative aspect-video overflow-hidden rounded-xl border border-dashed border-border bg-bg-elevated">
            {image ? (
              <img src={image} alt={title} className="h-full w-full object-cover" />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-text-muted">
                <span className="text-xs">No cover image</span>
              </div>
            )}
            <button
              type="button"
              onClick={() => coverInputRef.current?.click()}
              disabled={coverUploading}
              className="absolute bottom-2 right-2 inline-flex items-center gap-1.5 rounded-lg bg-bg-elevated/90 px-3 py-1.5 text-xs font-semibold text-foreground shadow-sm backdrop-blur hover:bg-bg-elevated disabled:opacity-50"
            >
              <Upload className="h-3.5 w-3.5" />
              {coverUploading ? "Uploading..." : (image ? "Replace" : "Upload")}
            </button>
            <input ref={coverInputRef} type="file" accept="image/*" className="hidden" onChange={handleCoverUpload} />
          </div>
        </div>

        <div className="lg:col-span-2 space-y-4">
          <Field label="Title">
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              onBlur={() => { if (title !== blog.title) save.mutate({ title }); }}
              placeholder="Article title..."
              className={inputClasses({ className: "text-base font-semibold" })}
            />
          </Field>
          <Field label="Excerpt">
            <textarea
              value={excerpt}
              onChange={(e) => setExcerpt(e.target.value)}
              onBlur={() => { if (excerpt !== blog.excerpt) save.mutate({ excerpt }); }}
              rows={3}
              placeholder="A short summary readers see in lists and social previews."
              className={inputClasses({ multiline: true })}
            />
          </Field>
        </div>
      </section>

      {/* Word-style editor */}
      <section>
        <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-text-muted">Content</p>
` + blogEditorNewUsage + `
      </section>

      <div className="mt-3 flex items-center justify-between text-xs text-text-muted">
        <p>Autosaves when you leave a field. Last updated {new Date(blog.updated_at).toLocaleString()}.</p>
        <Link href={"/blog/" + blog.slug} target="_blank" rel="noopener noreferrer" className="text-accent hover:text-accent-hover">
          View public page →
        </Link>
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1 block text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</span>
      {children}
    </label>
  );
}

// BlogDetailSkeleton paints the same shape the loaded page will use:
// header band on top, a 3-col grid for cover + title/excerpt, then a
// tall editor placeholder. Keeps layout from jumping when data arrives.
function BlogDetailSkeleton() {
  return (
    <div className="animate-pulse">
      <div className="mb-6 h-16 rounded-xl bg-bg-hover" />
      <section className="mb-6 grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div className="aspect-video rounded-xl bg-bg-hover lg:col-span-1" />
        <div className="space-y-3 lg:col-span-2">
          <div className="h-10 rounded-lg bg-bg-hover" />
          <div className="h-20 rounded-lg bg-bg-hover" />
        </div>
      </section>
      <div className="h-[500px] rounded-xl bg-bg-hover" />
    </div>
  );
}
`
}

// legacyAdminBlogsListPage generates the blogs list page (replaces the stock
// resource page for blogs). Renders the standard list + stats, but
// intercepts the New Blog button to open a sheet with just title +
// cover + excerpt; on submit it POSTs and redirects to the detail page.
func legacyAdminBlogsListPage() string {
	return `"use client";

import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import { IconButton } from "@/components/ui/IconButton";
import { SkeletonTable } from "@/components/ui/Skeleton";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import { apiClient, uploadFile } from "@/lib/api-client";
import { Plus, Upload, FileText, Loader2 } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";

interface Blog {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  image: string;
  published: boolean;
  created_at: string;
}

interface ListResponse { data: Blog[] }
interface ApiResponse<T> { data: T }

export default function BlogsPage() {
  const [open, setOpen] = useState(false);

  const { data, isLoading } = useQuery<Blog[]>({
    queryKey: ["blogs"],
    queryFn: async () => {
      const { data } = await apiClient.get<ListResponse>("/api/admin/blogs?page_size=100");
      return data.data;
    },
  });

  return (
    <div>
      <PageHeader
        title="Blogs"
        subtitle="Articles, drafts, and published posts."
        actions={
          <IconButton
            icon={<Plus className="h-4 w-4" />}
            label="New Blog"
            onClick={() => setOpen(true)}
          />
        }
      />

      {isLoading ? (
        <SkeletonTable rows={6} columns={3} />
      ) : (data?.length ?? 0) === 0 ? (
        <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
          <FileText className="mx-auto h-10 w-10 text-text-muted" />
          <p className="mt-3 text-base font-medium text-foreground">No blogs yet</p>
          <p className="mt-1 text-sm text-text-muted">Click New Blog to draft your first article.</p>
          <button
            type="button"
            onClick={() => setOpen(true)}
            className={buttonClasses({ className: "mt-4" })}
          >
            <Plus className="h-4 w-4" />
            New Blog
          </button>
        </div>
      ) : (
        <ul className="space-y-2">
          {(data ?? []).map((b) => (
            <li key={b.id}>
              <Link
                href={"/resources/blogs/" + b.id}
                className="flex gap-4 rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover hover:border-accent/30"
              >
                <div className="h-20 w-32 shrink-0 overflow-hidden rounded-lg bg-bg-hover">
                  {b.image ? (
                    <img src={b.image} alt={b.title} className="h-full w-full object-cover" />
                  ) : (
                    <div className="flex h-full w-full items-center justify-center text-text-muted">
                      <FileText className="h-6 w-6" />
                    </div>
                  )}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <p className="truncate text-sm font-semibold text-foreground">{b.title}</p>
                    <span className={"shrink-0 rounded px-1.5 py-0.5 text-[10px] font-semibold uppercase " + (b.published ? "bg-success/10 text-success" : "bg-bg-hover text-text-muted")}>
                      {b.published ? "Published" : "Draft"}
                    </span>
                  </div>
                  {b.excerpt && <p className="mt-1 line-clamp-2 text-sm text-text-secondary">{b.excerpt}</p>}
                  <p className="mt-2 text-xs text-text-muted">Created {new Date(b.created_at).toLocaleString()}</p>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}

      <NewBlogSheet open={open} onClose={() => setOpen(false)} />
    </div>
  );
}

function NewBlogSheet({ open, onClose }: { open: boolean; onClose: () => void }) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [excerpt, setExcerpt] = useState("");
  const [image, setImage] = useState("");
  const [uploading, setUploading] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);

  const create = useToastedMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post<ApiResponse<Blog>>("/api/admin/blogs", {
        title,
        excerpt,
        image,
        content: "",
        published: false,
      });
      return data.data;
    },
    successMessage: "Draft created — opening editor",
    onSuccess: (blog) => {
      queryClient.invalidateQueries({ queryKey: ["blogs"] });
      // Reset for next open
      setTitle(""); setExcerpt(""); setImage("");
      onClose();
      router.push("/resources/blogs/" + blog.id);
    },
  });

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploading(true);
    try {
      const res = await uploadFile(file);
      const url = (res.data as Record<string, unknown>)?.url as string;
      if (url) setImage(url);
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = "";
    }
  };

  const submit = () => {
    if (!title.trim()) return;
    create.mutate();
  };

  return (
    <ResponsiveSheet
      open={open}
      onClose={onClose}
      title="New blog"
      description="Add the title and cover image. You'll write the article on the next screen."
      footer={
        <>
          <button type="button" onClick={onClose} className="rounded-lg px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover">
            Cancel
          </button>
          <button
            type="button"
            onClick={submit}
            disabled={!title.trim() || create.isPending}
            className={buttonClasses()}
          >
            {create.isPending ? "Creating..." : "Continue to editor"}
          </button>
        </>
      }
    >
      <form onSubmit={(e) => { e.preventDefault(); submit(); }} className="space-y-4">
        <Field label="Title" required>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Your article's headline"
            autoFocus
            className={inputClasses()}
          />
        </Field>

        <Field label="Cover image">
          <div className="space-y-2">
            <div className="aspect-video overflow-hidden rounded-lg border border-dashed border-border bg-bg-elevated">
              {image ? (
                <img src={image} alt="Cover preview" className="h-full w-full object-cover" />
              ) : (
                <div className="flex h-full w-full items-center justify-center text-xs text-text-muted">
                  16:9 cover image
                </div>
              )}
            </div>
            <input ref={fileRef} type="file" accept="image/*" className="hidden" onChange={handleUpload} />
            <button
              type="button"
              onClick={() => fileRef.current?.click()}
              disabled={uploading}
              className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm font-medium text-foreground hover:bg-bg-hover disabled:opacity-50"
            >
              {uploading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Upload className="h-4 w-4" />}
              {image ? "Replace cover" : "Upload cover"}
            </button>
          </div>
        </Field>

        <Field label="Excerpt">
          <textarea
            value={excerpt}
            onChange={(e) => setExcerpt(e.target.value)}
            rows={3}
            placeholder="A short summary readers see in lists and social previews."
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
`
}
