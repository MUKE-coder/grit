package scaffold

// v3.31.7 — Blog two-step flow and Word-style editor.
//
// New flow:
//   1. List page "New Blog" -> sheet with Title + Cover Image + Excerpt
//   2. Submit -> POST /api/blogs -> redirect to /resources/blogs/[id]
//   3. Detail page renders a richer Tiptap editor (Word-style toolbar)
//      with autosave on blur and a manual Save button. Publish/Unpublish
//      + Delete are on the page header.
//
// Files generated:
//   components/forms/word-editor.tsx       (Word-style Tiptap editor)
//   app/(dashboard)/resources/blogs/[id]/page.tsx (detail + editor)
//   app/(dashboard)/resources/blogs/page.tsx       (redirects create -> detail)

// adminWordEditor generates the Word-style rich text editor. It's a
// thin wrapper around Tiptap so apps can drop <WordEditor value onChange />
// anywhere a full-page content surface is wanted (not just blogs).
func adminWordEditor() string {
	return `"use client";

import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import Link from "@tiptap/extension-link";
import TextAlign from "@tiptap/extension-text-align";
import TextStyle from "@tiptap/extension-text-style";
import Color from "@tiptap/extension-color";
import Highlight from "@tiptap/extension-highlight";
import Image from "@tiptap/extension-image";
import Table from "@tiptap/extension-table";
import TableRow from "@tiptap/extension-table-row";
import TableCell from "@tiptap/extension-table-cell";
import TableHeader from "@tiptap/extension-table-header";
import Placeholder from "@tiptap/extension-placeholder";
import { useEffect, useRef } from "react";
import { uploadFile } from "@/lib/api-client";
import {
  Bold, Italic, Underline as UnderlineIcon, Strikethrough,
  Heading1, Heading2, Heading3,
  List, ListOrdered, Quote, Code,
  Link as LinkIcon, Image as ImageIcon, Undo, Redo,
  AlignLeft, AlignCenter, AlignRight, AlignJustify,
  Highlighter, Palette, Table as TableIcon, Minus,
} from "@/lib/icons";

interface WordEditorProps {
  value: string;
  onChange: (html: string) => void;
  placeholder?: string;
  /** Called when the editor loses focus — handy for autosave. */
  onBlur?: () => void;
  /** Force editor height when used outside a flex layout. */
  minHeight?: number;
}

/**
 * MS-Word-style Tiptap editor.
 *
 * Toolbar: undo/redo · headings 1-3 · bold/italic/underline/strikethrough ·
 * text color · highlight · alignment 4-way · bullet/ordered lists ·
 * blockquote/code · link/image/table · horizontal rule.
 *
 * The container is a single rounded card so it reads like a Word page;
 * the toolbar sticks to the top of the card so it stays accessible while
 * scrolling through long content.
 */
export function WordEditor({ value, onChange, placeholder, onBlur, minHeight }: WordEditorProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const editor = useEditor({
    extensions: [
      StarterKit.configure({ heading: { levels: [1, 2, 3] } }),
      Underline,
      Link.configure({
        openOnClick: false,
        HTMLAttributes: { class: "text-accent underline" },
      }),
      TextAlign.configure({ types: ["heading", "paragraph"] }),
      TextStyle,
      Color,
      Highlight.configure({ multicolor: true }),
      Image.configure({
        HTMLAttributes: { class: "rounded-lg my-3 max-w-full" },
        allowBase64: false,
      }),
      Table.configure({ resizable: true, HTMLAttributes: { class: "border-collapse border border-border my-3" } }),
      TableRow,
      TableHeader.configure({ HTMLAttributes: { class: "border border-border bg-bg-hover px-3 py-2 font-semibold" } }),
      TableCell.configure({ HTMLAttributes: { class: "border border-border px-3 py-2 align-top" } }),
      Placeholder.configure({ placeholder: placeholder || "Start writing..." }),
    ],
    content: value || "",
    onUpdate: ({ editor }) => onChange(editor.getHTML()),
    onBlur: () => onBlur?.(),
    editorProps: {
      attributes: {
        // Word-page feel: white-ish surface, generous padding, prose
        // typography that survives the dark/light flip via theme tokens.
        class:
          "prose max-w-none p-8 focus:outline-none text-foreground " +
          "prose-headings:text-foreground prose-headings:font-bold " +
          "prose-p:text-foreground prose-strong:text-foreground prose-em:text-foreground " +
          "prose-li:text-foreground prose-a:text-accent " +
          "prose-blockquote:text-text-secondary prose-blockquote:border-l-4 prose-blockquote:border-accent prose-blockquote:bg-bg-hover prose-blockquote:px-4 prose-blockquote:py-1 " +
          "prose-code:text-accent prose-code:bg-bg-hover prose-code:rounded prose-code:px-1.5 prose-code:py-0.5 prose-code:before:content-none prose-code:after:content-none " +
          "prose-pre:bg-bg-secondary prose-pre:border prose-pre:border-border prose-pre:rounded-lg",
      },
    },
    immediatelyRender: false,
  });

  useEffect(() => {
    if (editor && value !== editor.getHTML()) {
      // emitUpdate=false suppresses the onUpdate echo.
      editor.commands.setContent(value || "", false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value, editor]);

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !editor) return;
    try {
      const res = await uploadFile(file);
      const url = (res.data as Record<string, unknown>)?.url as string;
      if (url) editor.chain().focus().setImage({ src: url, alt: file.name }).run();
    } catch {
      // Image upload errors surface via the toaster from useToastedMutation
      // when wired by the caller — the editor itself stays silent.
    } finally {
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  const insertLink = () => {
    if (!editor) return;
    const prev = editor.getAttributes("link").href as string | undefined;
    const url = window.prompt("Link URL", prev || "https://");
    if (url === null) return;
    if (url === "") {
      editor.chain().focus().extendMarkRange("link").unsetLink().run();
      return;
    }
    editor.chain().focus().extendMarkRange("link").setLink({ href: url }).run();
  };

  const insertTable = () => {
    editor?.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run();
  };

  if (!editor) {
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-8 text-sm text-text-muted">
        Loading editor...
      </div>
    );
  }

  return (
    <div className="rounded-xl border border-border bg-bg-elevated shadow-sm overflow-hidden">
      {/* Sticky toolbar — stays anchored to the top of the editor card
          so it's always accessible while writing long articles. */}
      <div className="sticky top-0 z-10 flex flex-wrap items-center gap-0.5 border-b border-border bg-bg-elevated/95 backdrop-blur px-3 py-2">
        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().undo().run()} aria-label="Undo" disabled={!editor.can().undo()}>
            <Undo className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().redo().run()} aria-label="Redo" disabled={!editor.can().redo()}>
            <Redo className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()} active={editor.isActive("heading", { level: 1 })} aria-label="Heading 1">
            <Heading1 className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()} active={editor.isActive("heading", { level: 2 })} aria-label="Heading 2">
            <Heading2 className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()} active={editor.isActive("heading", { level: 3 })} aria-label="Heading 3">
            <Heading3 className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleBold().run()} active={editor.isActive("bold")} aria-label="Bold">
            <Bold className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleItalic().run()} active={editor.isActive("italic")} aria-label="Italic">
            <Italic className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleUnderline().run()} active={editor.isActive("underline")} aria-label="Underline">
            <UnderlineIcon className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleStrike().run()} active={editor.isActive("strike")} aria-label="Strikethrough">
            <Strikethrough className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        {/* Text color + highlight — pick from a small native palette so
            we don't pull in a separate picker component for v3.31. */}
        <ToolbarGroup>
          <label className="relative inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-md text-text-secondary hover:bg-bg-hover" title="Text color">
            <Palette className="h-4 w-4" />
            <input
              type="color"
              onInput={(e) => editor.chain().focus().setColor((e.target as HTMLInputElement).value).run()}
              className="absolute inset-0 cursor-pointer opacity-0"
            />
          </label>
          <label className="relative inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-md text-text-secondary hover:bg-bg-hover" title="Highlight">
            <Highlighter className="h-4 w-4" />
            <input
              type="color"
              defaultValue="#fef08a"
              onInput={(e) => editor.chain().focus().toggleHighlight({ color: (e.target as HTMLInputElement).value }).run()}
              className="absolute inset-0 cursor-pointer opacity-0"
            />
          </label>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("left").run()} active={editor.isActive({ textAlign: "left" })} aria-label="Align left">
            <AlignLeft className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("center").run()} active={editor.isActive({ textAlign: "center" })} aria-label="Align center">
            <AlignCenter className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("right").run()} active={editor.isActive({ textAlign: "right" })} aria-label="Align right">
            <AlignRight className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("justify").run()} active={editor.isActive({ textAlign: "justify" })} aria-label="Justify">
            <AlignJustify className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleBulletList().run()} active={editor.isActive("bulletList")} aria-label="Bullet list">
            <List className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleOrderedList().run()} active={editor.isActive("orderedList")} aria-label="Ordered list">
            <ListOrdered className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleBlockquote().run()} active={editor.isActive("blockquote")} aria-label="Blockquote">
            <Quote className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleCodeBlock().run()} active={editor.isActive("codeBlock")} aria-label="Code block">
            <Code className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={insertLink} active={editor.isActive("link")} aria-label="Insert link">
            <LinkIcon className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => fileInputRef.current?.click()} aria-label="Insert image">
            <ImageIcon className="h-4 w-4" />
          </ToolbarBtn>
          <input ref={fileInputRef} type="file" accept="image/*" className="hidden" onChange={handleImageUpload} />
          <ToolbarBtn onClick={insertTable} aria-label="Insert table">
            <TableIcon className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setHorizontalRule().run()} aria-label="Horizontal rule">
            <Minus className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>
      </div>

      <div style={minHeight ? { minHeight: minHeight + "px" } : undefined}>
        <EditorContent editor={editor} />
      </div>
    </div>
  );
}

function ToolbarGroup({ children }: { children: React.ReactNode }) {
  return <div className="flex items-center gap-0.5">{children}</div>;
}

function ToolbarSep() {
  return <span className="mx-1 h-5 w-px shrink-0 bg-border" aria-hidden />;
}

interface ToolbarBtnProps {
  children: React.ReactNode;
  onClick: () => void;
  active?: boolean;
  disabled?: boolean;
  "aria-label": string;
}

function ToolbarBtn({ children, onClick, active, disabled, ...rest }: ToolbarBtnProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      aria-pressed={active}
      className={
        "inline-flex h-8 w-8 items-center justify-center rounded-md transition-colors disabled:opacity-30 disabled:hover:bg-transparent " +
        (active ? "bg-accent/10 text-accent" : "text-text-secondary hover:bg-bg-hover hover:text-foreground")
      }
      {...rest}
    >
      {children}
    </button>
  );
}
`
}

// adminBlogDetailPage generates the blog detail/edit page. Loads the
// blog by id, surfaces the cover + title + excerpt at the top, then
// renders the WordEditor over the content field. Autosaves on blur +
// has explicit Save / Publish / Delete actions.
func adminBlogDetailPage() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { WordEditor } from "@/components/forms/word-editor";
import { IconButton } from "@/components/ui/IconButton";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import { apiClient, uploadFile } from "@/lib/api-client";
import { ArrowLeft, Save, Trash2, Upload, Check } from "@/lib/icons";

interface Blog {
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
  const coverInputRef = useRef<HTMLInputElement>(null);

  const { data: blog, isLoading } = useQuery<Blog>({
    queryKey: ["blog", params.id],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<Blog>>("/api/blogs/" + params.id);
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
      const { data } = await apiClient.put<ApiResponse<Blog>>("/api/blogs/" + params.id, patch);
      return data.data;
    },
    successMessage: "Saved",
    silentSuccess: true,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["blog", params.id] }),
  });

  const publish = useToastedMutation({
    mutationFn: async (next: boolean) => {
      const { data } = await apiClient.put<ApiResponse<Blog>>("/api/blogs/" + params.id, { published: next });
      return data.data;
    },
    successMessage: (b) => b.published ? "Published" : "Moved back to draft",
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["blog", params.id] }),
  });

  const del = useToastedMutation({
    mutationFn: async () => apiClient.delete("/api/blogs/" + params.id),
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
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
        <div className="inline-block h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
      </div>
    );
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
              onClick={() => { if (window.confirm("Delete this blog? This cannot be undone.")) del.mutate(); }}
            />
          </>
        }
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
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-base font-semibold text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            />
          </Field>
          <Field label="Excerpt">
            <textarea
              value={excerpt}
              onChange={(e) => setExcerpt(e.target.value)}
              onBlur={() => { if (excerpt !== blog.excerpt) save.mutate({ excerpt }); }}
              rows={3}
              placeholder="A short summary readers see in lists and social previews."
              className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            />
          </Field>
        </div>
      </section>

      {/* Word-style editor */}
      <section>
        <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-text-muted">Content</p>
        <WordEditor
          value={content}
          onChange={setContent}
          placeholder="Start your article here. Use the toolbar to format headings, lists, tables, images, and more."
          minHeight={500}
          onBlur={() => { if (content !== blog.content) save.mutate({ content }); }}
        />
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
`
}

// adminBlogsListPage generates the blogs list page (replaces the stock
// resource page for blogs). Renders the standard list + stats, but
// intercepts the New Blog button to open a sheet with just title +
// cover + excerpt; on submit it POSTs and redirects to the detail page.
func adminBlogsListPage() string {
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
            className="mt-4 inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover"
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
            className="rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
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
            className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
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
            className="w-full rounded-lg border border-border bg-bg-elevated px-3 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
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
