import type { Extensions } from "@tiptap/core";
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
