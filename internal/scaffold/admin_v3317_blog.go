package scaffold

// The Word-style editor, components/forms/word-editor.tsx.
//
// v3.31.7 wrote it for a blog detail page of its own. Since contact-app review
// M42 and M43 it is the panel's only rich text editor: the rich text form field
// renders it, so the Blog resource runs on <ResourcePage> like every other
// resource, and its content schema lives in lib/tiptap-extensions.ts.

// adminWordEditor generates the Word-style rich text editor. It's a
// thin wrapper around Tiptap so apps can drop <WordEditor value onChange />
// anywhere a full-page content surface is wanted.
func adminWordEditor() string {
	return `"use client";

import { useEditor, useEditorState, EditorContent, type Editor } from "@tiptap/react";
import { useEffect, useMemo, useRef } from "react";
import { uploadFile } from "@/lib/api-client";
import { richTextExtensions } from "@/lib/tiptap-extensions";
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
  /** Called when the editor loses focus, handy for autosave. */
  onBlur?: () => void;
  /** Force editor height when used outside a flex layout. */
  minHeight?: number;
  /** Id of the element that names the editor, for screen readers. */
  labelledBy?: string;
  /** Marks the editor invalid, for a form field with an error. */
  invalid?: boolean;
}

// What the toolbar shows. Tiptap 3 no longer re-renders on every transaction,
// so the toolbar subscribes to exactly these and re-renders when one changes.
function toolbarState(editor: Editor | null) {
  if (!editor) return null;
  return {
    canUndo: editor.can().undo(),
    canRedo: editor.can().redo(),
    h1: editor.isActive("heading", { level: 1 }),
    h2: editor.isActive("heading", { level: 2 }),
    h3: editor.isActive("heading", { level: 3 }),
    bold: editor.isActive("bold"),
    italic: editor.isActive("italic"),
    underline: editor.isActive("underline"),
    strike: editor.isActive("strike"),
    alignLeft: editor.isActive({ textAlign: "left" }),
    alignCenter: editor.isActive({ textAlign: "center" }),
    alignRight: editor.isActive({ textAlign: "right" }),
    alignJustify: editor.isActive({ textAlign: "justify" }),
    bulletList: editor.isActive("bulletList"),
    orderedList: editor.isActive("orderedList"),
    blockquote: editor.isActive("blockquote"),
    codeBlock: editor.isActive("codeBlock"),
    link: editor.isActive("link"),
  };
}

/**
 * MS-Word-style Tiptap editor.
 *
 * Toolbar: undo/redo, headings 1-3, bold/italic/underline/strikethrough,
 * text color, highlight, alignment 4-way, bullet/ordered lists,
 * blockquote/code, link/image/table, horizontal rule.
 *
 * The container is a single rounded card so it reads like a Word page;
 * the toolbar sticks to the top of the card so it stays accessible while
 * scrolling through long content.
 */
export function WordEditor({ value, onChange, placeholder, onBlur, minHeight, labelledBy, invalid }: WordEditorProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Memoised: a new extension array on every render reads to Tiptap as new
  // options, and it would reconfigure the editor on each keystroke.
  const extensions = useMemo(
    () => richTextExtensions({ placeholder: placeholder || "Start writing..." }),
    [placeholder],
  );

  const attributes: Record<string, string> = {
    role: "textbox",
    "aria-multiline": "true",
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
  };
  if (labelledBy) attributes["aria-labelledby"] = labelledBy;
  if (invalid) attributes["aria-invalid"] = "true";

  const editor = useEditor({
    extensions,
    content: value || "",
    onUpdate: ({ editor }) => onChange(editor.getHTML()),
    onBlur: () => onBlur?.(),
    editorProps: { attributes },
    // Tiptap renders on the server otherwise, and its DOM never matches what
    // React produces on the client. The editor mounts in an effect instead.
    immediatelyRender: false,
  });

  const state = useEditorState({ editor, selector: ({ editor }) => toolbarState(editor) });

  useEffect(() => {
    if (editor && value !== editor.getHTML()) {
      // emitUpdate: false suppresses the onUpdate echo.
      editor.commands.setContent(value || "", { emitUpdate: false });
    }
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
      // when wired by the caller; the editor itself stays silent.
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

  if (!editor || !state) {
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-8 text-sm text-text-muted">
        Loading editor...
      </div>
    );
  }

  return (
    <div className={"rounded-xl border bg-bg-elevated shadow-sm overflow-hidden " + (invalid ? "border-danger" : "border-border")}>
      {/* Sticky toolbar, anchored to the top of the editor card
          so it's always accessible while writing long articles. */}
      <div className="sticky top-0 z-10 flex flex-wrap items-center gap-0.5 border-b border-border bg-bg-elevated/95 backdrop-blur px-3 py-2">
        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().undo().run()} aria-label="Undo" disabled={!state.canUndo}>
            <Undo className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().redo().run()} aria-label="Redo" disabled={!state.canRedo}>
            <Redo className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()} active={state.h1} aria-label="Heading 1">
            <Heading1 className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()} active={state.h2} aria-label="Heading 2">
            <Heading2 className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()} active={state.h3} aria-label="Heading 3">
            <Heading3 className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleBold().run()} active={state.bold} aria-label="Bold">
            <Bold className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleItalic().run()} active={state.italic} aria-label="Italic">
            <Italic className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleUnderline().run()} active={state.underline} aria-label="Underline">
            <UnderlineIcon className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleStrike().run()} active={state.strike} aria-label="Strikethrough">
            <Strikethrough className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        {/* Text color + highlight: pick from a small native palette so
            we don't pull in a separate picker component. */}
        <ToolbarGroup>
          <label className="relative inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-md text-text-secondary hover:bg-bg-hover" title="Text color">
            <Palette className="h-4 w-4" />
            <input
              type="color"
              aria-label="Text color"
              onInput={(e) => editor.chain().focus().setColor((e.target as HTMLInputElement).value).run()}
              className="absolute inset-0 cursor-pointer opacity-0"
            />
          </label>
          <label className="relative inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-md text-text-secondary hover:bg-bg-hover" title="Highlight">
            <Highlighter className="h-4 w-4" />
            <input
              type="color"
              aria-label="Highlight"
              defaultValue="#fef08a"
              onInput={(e) => editor.chain().focus().toggleHighlight({ color: (e.target as HTMLInputElement).value }).run()}
              className="absolute inset-0 cursor-pointer opacity-0"
            />
          </label>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("left").run()} active={state.alignLeft} aria-label="Align left">
            <AlignLeft className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("center").run()} active={state.alignCenter} aria-label="Align center">
            <AlignCenter className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("right").run()} active={state.alignRight} aria-label="Align right">
            <AlignRight className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().setTextAlign("justify").run()} active={state.alignJustify} aria-label="Justify">
            <AlignJustify className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleBulletList().run()} active={state.bulletList} aria-label="Bullet list">
            <List className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleOrderedList().run()} active={state.orderedList} aria-label="Ordered list">
            <ListOrdered className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleBlockquote().run()} active={state.blockquote} aria-label="Blockquote">
            <Quote className="h-4 w-4" />
          </ToolbarBtn>
          <ToolbarBtn onClick={() => editor.chain().focus().toggleCodeBlock().run()} active={state.codeBlock} aria-label="Code block">
            <Code className="h-4 w-4" />
          </ToolbarBtn>
        </ToolbarGroup>

        <ToolbarSep />

        <ToolbarGroup>
          <ToolbarBtn onClick={insertLink} active={state.link} aria-label="Insert link">
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
