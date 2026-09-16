package scaffold

import (
	"fmt"
	"strings"
)

// Contact-app review M43, and the move to Tiptap 3.
//
// The panel had two Tiptap editors with two content schemas: the form field
// knew StarterKit and Link, the blog's Word-style editor added underline,
// alignment, colour, highlight, images and tables. A post written in one and
// saved from the other lost everything the second did not know, silently,
// because ProseMirror drops what its schema cannot represent.
//
// Now there is one schema, lib/tiptap-extensions.ts, and one editor,
// components/forms/word-editor.tsx. The rich text field is that editor with a
// label.
//
// Tiptap 2 is also gone: @tiptap/core below 3.30.4 carries GHSA-cp6q-959q-f8rh
// (mergeAttributes() turns an own __proto__ key into inherited DOM attributes),
// so every generated project failed pnpm audit.

// tiptapVersion is the one version every @tiptap package is pinned to.
//
// Exact, not a range: each Tiptap 3 package names its siblings as exact peer
// dependencies (@tiptap/react 3.31.3 wants @tiptap/core 3.31.3), so two that
// drift apart install two copies of core and the editor throws on a keyed
// plugin registered twice.
const tiptapVersion = "3.31.3"

// tiptapPackages are the packages the editor imports, sorted.
//
// Tiptap 3 folded several of the old ones in: Link and Underline are in
// StarterKit, TableRow, TableCell and TableHeader are in @tiptap/extension-table,
// Color is in @tiptap/extension-text-style, and Placeholder is in
// @tiptap/extensions. @tiptap/core is listed because the extension list imports
// its types, and because every other package wants it as a peer.
var tiptapPackages = []string{
	"@tiptap/core",
	"@tiptap/extension-highlight",
	"@tiptap/extension-image",
	"@tiptap/extension-table",
	"@tiptap/extension-text-align",
	"@tiptap/extension-text-style",
	"@tiptap/extensions",
	"@tiptap/pm",
	"@tiptap/react",
	"@tiptap/starter-kit",
}

// tiptapDependencyLines is the package.json lines for tiptapPackages, each
// indented by indent and followed by a comma.
func tiptapDependencyLines(indent string) string {
	var b strings.Builder
	for _, pkg := range tiptapPackages {
		fmt.Fprintf(&b, "%s%q: %q,\n", indent, pkg, tiptapVersion)
	}
	return b.String()
}

// adminTiptapExtensions is lib/tiptap-extensions.ts.
func adminTiptapExtensions() string {
	return `import type { Extensions } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import { Color, TextStyle } from "@tiptap/extension-text-style";
import TextAlign from "@tiptap/extension-text-align";
import Highlight from "@tiptap/extension-highlight";
import Image from "@tiptap/extension-image";
import { TableKit } from "@tiptap/extension-table";
import { Placeholder } from "@tiptap/extensions";

export interface RichTextExtensionOptions {
  /** Shown while the document is empty. Unset shows nothing. */
  placeholder?: string;
}

/**
 * The content schema of every rich text surface in the panel.
 *
 * One list, because ProseMirror silently drops whatever its schema cannot
 * represent. When the rich text field and the Word-style editor each had their
 * own, a post written with a table in one lost the table the moment it was
 * saved from the other. Add an extension here and every editor can hold it.
 *
 * Link and Underline come with StarterKit in Tiptap 3, and TableKit brings the
 * row, header and cell nodes with the table.
 */
export function richTextExtensions({ placeholder }: RichTextExtensionOptions = {}): Extensions {
  const extensions: Extensions = [
    StarterKit.configure({
      heading: { levels: [1, 2, 3] },
      link: {
        openOnClick: false,
        HTMLAttributes: { class: "text-accent underline" },
      },
    }),
    TextAlign.configure({ types: ["heading", "paragraph"] }),
    TextStyle,
    Color,
    Highlight.configure({ multicolor: true }),
    Image.configure({
      HTMLAttributes: { class: "rounded-lg my-3 max-w-full" },
      allowBase64: false,
    }),
    TableKit.configure({
      table: { resizable: true, HTMLAttributes: { class: "border-collapse border border-border my-3" } },
      tableHeader: { HTMLAttributes: { class: "border border-border bg-bg-hover px-3 py-2 font-semibold" } },
      tableCell: { HTMLAttributes: { class: "border border-border px-3 py-2 align-top" } },
    }),
  ];
  if (placeholder) {
    extensions.push(Placeholder.configure({ placeholder }));
  }
  return extensions;
}
`
}

// adminRichTextField is components/forms/fields/rich-text-field.tsx: the one
// editor, with the label and error every other field has.
func adminRichTextField() string {
	return `"use client";

import { useId } from "react";
import { WordEditor } from "@/components/forms/word-editor";

interface RichTextFieldProps {
  field: { key: string; label: string; required?: boolean; placeholder?: string; description?: string };
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/**
 * A rich text form field. It is the Word-style editor rather than a smaller
 * one of its own, so a record edited here keeps every table, colour and
 * alignment it was written with: both read and write lib/tiptap-extensions.ts.
 */
export function RichTextField({ field, value, onChange, error }: RichTextFieldProps) {
  const labelId = useId();
  return (
    <div className="space-y-1.5">
      <p id={labelId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="ml-1 text-danger">*</span>}
      </p>
      <WordEditor
        value={value}
        onChange={onChange}
        placeholder={field.placeholder}
        minHeight={200}
        labelledBy={labelId}
        invalid={!!error}
      />
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

// adminTiptapExtensionsTest is __tests__/tiptap-extensions.test.ts: a post with
// everything the toolbar makes, loaded and saved twice through the one schema.
// libImport is where lib/ resolves from in this admin.
func adminTiptapExtensionsTest(libImport string) string {
	return strings.ReplaceAll(`import { afterEach, describe, expect, it } from "vitest";
import { Editor, type JSONContent } from "@tiptap/core";
import { richTextExtensions } from "{{LIB}}/tiptap-extensions";

// A post with everything the editor's toolbar can make. When the rich text
// field and the Word-style editor had two extension lists, saving this from
// the smaller one dropped the table, the colours, the highlight and the
// alignment without a word.
const post = [
  "<h1>Release notes</h1>",
  "<h2>What changed</h2>",
  '<p style="text-align: center">Centred <strong>bold</strong>, <em>italic</em>, <u>underlined</u> and <s>struck</s>.</p>',
  '<p><a href="https://example.com">a link</a>, <span style="color: #ff0000">red text</span> and <mark data-color="#fef08a" style="background-color: #fef08a; color: inherit">a highlight</mark></p>',
  "<ul><li><p>one</p></li><li><p>two</p></li></ul>",
  "<ol><li><p>first</p></li></ol>",
  "<blockquote><p>quoted</p></blockquote>",
  "<pre><code>go run .</code></pre>",
  '<img src="https://example.com/cover.png" alt="cover">',
  "<table><tbody><tr><th><p>Plan</p></th><th><p>Price</p></th></tr><tr><td><p>Pro</p></td><td><p>$9</p></td></tr></tbody></table>",
].join("");

const editors: Editor[] = [];

function load(content: string): Editor {
  const editor = new Editor({ extensions: richTextExtensions(), content });
  editors.push(editor);
  return editor;
}

function kinds(node: JSONContent, found = new Set<string>()): Set<string> {
  if (node.type) found.add(node.type);
  for (const mark of node.marks ?? []) found.add("mark:" + mark.type);
  if (node.attrs?.textAlign && node.attrs.textAlign !== "left") found.add("align:" + node.attrs.textAlign);
  for (const child of node.content ?? []) kinds(child, found);
  return found;
}

afterEach(() => {
  while (editors.length) editors.pop()?.destroy();
});

describe("the rich text schema", () => {
  it("keeps every kind of formatting the toolbar makes", () => {
    const found = kinds(load(post).getJSON());
    for (const kind of [
      "heading", "paragraph", "bulletList", "orderedList", "listItem", "blockquote", "codeBlock",
      "image", "table", "tableRow", "tableHeader", "tableCell",
      "mark:bold", "mark:italic", "mark:underline", "mark:strike", "mark:link", "mark:textStyle", "mark:highlight",
      "align:center",
    ]) {
      expect(found, kind).toContain(kind);
    }
  });

  it("round-trips a saved post without losing anything", () => {
    const first = load(post).getHTML();
    const second = load(first).getHTML();
    expect(second).toBe(first);
    expect(first).toContain("<table");
    // jsdom writes the colour back as rgb(), so match the span, not the value.
    expect(first).toMatch(/<span style="color: [^"]+">red text<\/span>/);
    expect(first).toContain("<mark");
    expect(first).toContain("text-align: center");
    expect(first).toContain('href="https://example.com"');
  });
});
`, "{{LIB}}", libImport)
}
