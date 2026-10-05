import { afterEach, describe, expect, it } from "vitest";
import { Editor, type JSONContent } from "@tiptap/core";
import { richTextExtensions } from "@/lib/tiptap-extensions";

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
