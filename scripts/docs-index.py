#!/usr/bin/env python3
"""Build the searchable documentation index the CLI embeds.

An agent working in a Grit project needs to find the right API for the version
that project is pinned to. The docs site is the source of truth and it is a
Next.js app: an installed `grit` binary cannot read it, and a developer offline
cannot reach it. So the prose is extracted here and embedded in the binary,
which also means the index is pinned to the CLI's own version and can say so
when the project disagrees.

Run from the repository root, as part of the release ritual:

    python scripts/docs-index.py

Writes internal/docs/index.json. Commit it: the build embeds it, so a stale one
ships silently and the whole point of the index is being right about a version.
"""

import html
import io
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOCS = os.path.join(ROOT, "docs")
OUT = os.path.join(ROOT, "internal", "docs", "index.json")


def read(path):
    with io.open(path, encoding="utf-8", errors="replace") as handle:
        return handle.read()


def cli_version():
    """The version this index belongs to, from the one place that declares it."""
    main_go = read(os.path.join(ROOT, "cmd", "grit", "main.go"))
    match = re.search(r'var version = "([^"]+)"', main_go)
    if not match:
        sys.exit("cannot find `var version` in cmd/grit/main.go")
    return match.group(1)


def metadata():
    """url -> (title, description), from docs/config/docs-metadata.ts.

    Parsed rather than evaluated. It is a TypeScript object literal of string
    pairs and nothing more, so a regular expression is honest here in a way it
    would not be for the pages themselves.
    """
    source = read(os.path.join(DOCS, "config", "docs-metadata.ts"))
    out = {}
    # '/docs/x': { title: '...', description: '...' }
    pattern = re.compile(
        r"'(/docs[^']*)':\s*\{\s*title:\s*(['\"])(.*?)\2\s*,\s*description:\s*$",
        re.S,
    )
    # Simpler and more robust: find each key, then the next title and
    # description string after it.
    for key in re.finditer(r"^\s*'(/docs[^']*)':\s*\{", source, re.M):
        url = key.group(1)
        chunk = source[key.end(): key.end() + 1200]
        title = re.search(r"title:\s*'((?:[^'\\]|\\.)*)'", chunk)
        desc = re.search(r"description:\s*\n?\s*'((?:[^'\\]|\\.)*)'", chunk)
        out[url] = (
            unescape(title.group(1)) if title else "",
            unescape(desc.group(1)) if desc else "",
        )
    _ = pattern
    return out


def unescape(s):
    return s.replace("\\'", "'").replace('\\"', '"').replace("\\n", " ")


# JSX text extraction.
#
# The pages are TSX, so this is lossy by nature. It does not need to be perfect:
# the index exists to tell an agent which page answers a question and to give it
# enough prose to be sure, after which it reads the page or the URL. Being
# roughly right about 181 pages beats being exactly right about none.
TAG = re.compile(r"<[^>]+>")
EXPR = re.compile(r"\{[^{}]*\}")
IMPORT = re.compile(r"^\s*import .*$", re.M)
COMMENT_BLOCK = re.compile(r"/\*.*?\*/", re.S)
COMMENT_LINE = re.compile(r"^\s*//.*$", re.M)
WHITESPACE = re.compile(r"\s+")


def headings(source):
    """Every heading, which is what a reader scans for."""
    found = []
    for match in re.finditer(r"<h([1-6])[^>]*>(.*?)</h\1>", source, re.S):
        text = clean(match.group(2))
        if text:
            found.append(text)
    return found


def prose(source):
    """The paragraph text, flattened.

    Code blocks are left out on purpose. They are the bulk of these pages and
    the least searchable part of them: an agent looking for "how do I place a
    table on another tier" is helped by the sentence that explains it and not by
    forty lines of shell.
    """
    body = IMPORT.sub("", source)
    body = COMMENT_BLOCK.sub("", body)
    body = COMMENT_LINE.sub("", body)
    # Drop CodeBlock elements entirely, including their props.
    body = re.sub(r"<CodeBlock[\s\S]*?(/>|</CodeBlock>)", " ", body)

    parts = []
    for match in re.finditer(r"<(p|li|blockquote)[^>]*>(.*?)</\1>", body, re.S):
        text = clean(match.group(2))
        if len(text) > 2:
            parts.append(text)
    return parts


# A page's tables, when the rows are data rather than markup.
#
# The reference pages render their tables from a const array at the top of the
# file and a .map() in the body, so every <td> holds {row.type} and the text
# that matters is in the array. prose() reads <p>, <li> and <blockquote> and
# clean() strips expressions, so those tables came out of the index as nothing
# at all: the field-types page called itself the single source of truth for
# field mapping and contained no field types (grit#99). An assistant told to
# look a type up in the binary could not, and had to go and read
# internal/generate/field.go instead, which is the exact failure the embedded
# docs exist to prevent.
#
# Parsed rather than evaluated, and lossily: the string values of each object
# literal, in declaration order, joined as a row. That is enough to answer
# "what does money become" and cheap enough to be safe.
DATA_ARRAY = re.compile(r"^const\s+([A-Za-z_][A-Za-z0-9_]*)\s*(?::[^=\n]+)?=\s*\[", re.M)
ROW_VALUE = re.compile(r"[A-Za-z_][A-Za-z0-9_]*\s*:\s*(?:'((?:[^'\\]|\\.)*)'|\"((?:[^\"\\]|\\.)*)\"|`((?:[^`\\]|\\.)*)`)", re.S)


def data_tables(source):
    """Rows from the const arrays this page renders with .map()."""
    out = []
    for match in DATA_ARRAY.finditer(source):
        name = match.group(1)
        # Only arrays the page actually renders. A const that is never mapped
        # is configuration, not content, and indexing it adds noise.
        if name + ".map(" not in source:
            continue

        # Walk to the array's closing bracket, so a nested [] cannot end it
        # early.
        depth, i, end = 0, match.end() - 1, -1
        while i < len(source):
            if source[i] == "[":
                depth += 1
            elif source[i] == "]":
                depth -= 1
                if depth == 0:
                    end = i
                    break
            i += 1
        if end < 0:
            continue

        rows = []
        for obj in objects_in(source[match.end():end]):
            values = [unescape_js(a or b or c) for a, b, c in ROW_VALUE.findall(obj)]
            values = [v for v in values if v]
            if values:
                rows.append(" | ".join(values))
        if rows:
            out.append(humanise(name) + ": " + " // ".join(rows))
    return out


def objects_in(body):
    """Each top-level { ... } in an array body."""
    objs, depth, start = [], 0, -1
    for i, ch in enumerate(body):
        if ch == "{":
            if depth == 0:
                start = i
            depth += 1
        elif ch == "}":
            depth -= 1
            if depth == 0 and start >= 0:
                objs.append(body[start + 1:i])
                start = -1
    return objs


def unescape_js(raw):
    text = raw.replace("\\'", "'").replace('\\"', '"').replace("\\`", "`")
    text = text.replace("\\n", " ").replace("\\\\", "\\")
    return clean(text)


def humanise(const_name):
    """COLUMNS -> Columns, MODIFIER_TAGS -> Modifier tags."""
    words = const_name.replace("_", " ").strip()
    if words.isupper():
        words = words.capitalize()
    return words


def clean(fragment):
    text = EXPR.sub(" ", fragment)
    text = TAG.sub(" ", text)
    # html.unescape rather than a list of replacements. The first version
    # handled five entities by hand and &mdash; went straight through into the
    # search output, which is both ugly and, since this project bans em dashes
    # in prose, a tell that the entity was never really text.
    text = html.unescape(text)
    text = text.replace("{'", "").replace("'}", "")
    # The docs write em dashes as entities; prose here uses punctuation that a
    # terminal can show.
    text = text.replace("—", ", ").replace("–", "-")
    text = text.replace("’", "'").replace("‘", "'")
    text = text.replace("“", '"').replace("”", '"')
    text = text.replace("…", "...").replace(" ", " ")
    return WHITESPACE.sub(" ", text).strip()


def url_for(page_path):
    """app/docs/admin/custom-pages/page.tsx -> /docs/admin/custom-pages"""
    rel = os.path.relpath(page_path, os.path.join(DOCS, "app"))
    rel = rel.replace(os.sep, "/")
    rel = rel[: -len("/page.tsx")] if rel.endswith("/page.tsx") else rel
    return "/" + rel


def doc_pages(meta):
    pages = []
    root = os.path.join(DOCS, "app", "docs")
    for current, dirs, files in os.walk(root):
        dirs[:] = [d for d in dirs if d not in (".next", "node_modules")]
        if "page.tsx" not in files:
            continue
        path = os.path.join(current, "page.tsx")
        source = read(path)
        url = url_for(path)
        title, description = meta.get(url, ("", ""))
        if not title:
            heads = headings(source)
            title = heads[0] if heads else url
        pages.append({
            "url": url,
            "title": title,
            "description": description,
            "headings": headings(source),
            # Budgeted separately, because on a reference page the tables are
            # the content and the prose is the commentary. Sharing one cap let
            # the prose push the flags table off the end, which is the same
            # bug as not indexing it at all for anybody looking for a flag.
            "text": (" ".join(prose(source))[:6000] + " "
                     + " ".join(data_tables(source))[:8000]).strip(),
            "kind": "docs",
        })
    return pages


# The system design pages.
#
# These are data rather than files: one config entry per system, one dynamic
# route, one renderer. The walk above sees a single `[system]` page and indexes
# the template, so without this the whole section is unsearchable, which for
# twenty-six pages of prose is the difference between the index being useful
# and being a list of the pages somebody happened to write by hand.
#
# Parsed rather than evaluated, like docs-metadata.ts above, and lossily on
# purpose: every string literal long enough to be a sentence, which is the
# prose, plus the labels and the fixed section headings a reader scans for.
SYSTEM_SECTIONS = [
    "Problem statement",
    "System requirements",
    "Capacity estimation",
    "High level design",
    "Technology stack",
    "Data model",
    "API design",
    "Low level design",
    "Scalability and performance",
    "Bottlenecks and improvements",
]

TS_STRING = re.compile(r"'((?:[^'\\\n]|\\.)*)'")
UNICODE_ESCAPE = re.compile(r"\\u([0-9a-fA-F]{4})")


def ts_string(raw):
    """One single-quoted TypeScript literal, as the text it stands for."""
    text = UNICODE_ESCAPE.sub(lambda m: chr(int(m.group(1), 16)), raw)
    return unescape(text)


def system_pages():
    config = os.path.join(DOCS, "config")
    if not os.path.isdir(config):
        return []

    sources = []
    for name in sorted(os.listdir(config)):
        if name.startswith("systems") and name.endswith(".ts") and name != "systems.ts":
            sources.append(read(os.path.join(config, name)))

    pages = []
    for source in sources:
        # Each `export const X: SystemDesign = {` opens one system, and the
        # next one closes it.
        starts = [m.start() for m in re.finditer(r"^export const \w+: SystemDesign = \{", source, re.M)]
        for i, start in enumerate(starts):
            block = source[start: starts[i + 1] if i + 1 < len(starts) else len(source)]

            def field(key):
                match = re.search(r"\b%s:\s*'((?:[^'\\\n]|\\.)*)'" % key, block)
                return ts_string(match.group(1)) if match else ""

            slug = field("slug")
            name = field("name")
            if not slug or not name:
                continue

            labels = [ts_string(m) for m in re.findall(r"\blabel:\s*'((?:[^'\\\n]|\\.)*)'", block)]
            classes = [ts_string(m) for m in re.findall(r"\bname:\s*'((?:[^'\\\n]|\\.)*)'", block)]

            # Every literal long enough to be a sentence rather than a key, a
            # slug or a table cell. Short ones are mostly stack rows and node
            # labels, which the headings already carry.
            prose_bits = []
            for raw in TS_STRING.findall(block):
                text = clean(ts_string(raw))
                if len(text) >= 40:
                    prose_bits.append(text)

            pages.append({
                "url": "/docs/systems/%s" % slug,
                "title": "%s System Design" % name,
                "description": field("tagline"),
                "headings": SYSTEM_SECTIONS + sorted(set(labels + classes[1:])),
                "text": " ".join(prose_bits)[:6000],
                "kind": "docs",
            })
    return pages


FRONTMATTER = re.compile(r"^---\s*\n(.*?)\n---\s*\n", re.S)


def markdown_pages():
    """The blog and the lessons, which carry the long explanations."""
    pages = []
    for sub, kind in (("blog", "blog"), ("lessons", "lesson")):
        base = os.path.join(DOCS, "content", sub)
        if not os.path.isdir(base):
            continue
        for name in sorted(os.listdir(base)):
            if not name.endswith(".md"):
                continue
            source = read(os.path.join(base, name))
            title = ""
            front = FRONTMATTER.match(source)
            if front:
                match = re.search(r"^title:\s*['\"]?(.*?)['\"]?\s*$", front.group(1), re.M)
                if match:
                    title = match.group(1)
                source = source[front.end():]
            slug = name[:-3]
            # Strip fenced code, for the same reason as above.
            body = re.sub(r"```[\s\S]*?```", " ", source)
            heads = [h.strip("# ").strip() for h in re.findall(r"^#{1,4} .*$", body, re.M)]
            body = re.sub(r"^#{1,6} .*$", " ", body, flags=re.M)
            body = WHITESPACE.sub(" ", TAG.sub(" ", body)).strip()
            pages.append({
                "url": "/%s/%s" % (sub, slug),
                "title": title or slug,
                "description": "",
                "headings": heads,
                "text": body[:6000],
                "kind": kind,
            })
    return pages


def main():
    if not os.path.isdir(DOCS):
        sys.exit("docs/ not found; run this from the repository root")

    meta = metadata()
    pages = doc_pages(meta) + system_pages() + markdown_pages()
    # The dynamic route itself indexes the template, not a page anybody reads.
    pages = [p for p in pages if "[" not in p["url"]]
    pages.sort(key=lambda p: p["url"])

    index = {
        "version": cli_version(),
        "pages": pages,
    }

    os.makedirs(os.path.dirname(OUT), exist_ok=True)
    # newline="" so Windows does not turn this into CRLF, which would change
    # the embedded bytes on every machine that touched it.
    with io.open(OUT, "w", encoding="utf-8", newline="") as handle:
        json.dump(index, handle, ensure_ascii=False, separators=(",", ":"), sort_keys=True)
        handle.write("\n")

    with_text = sum(1 for p in pages if len(p["text"]) > 200)
    size = os.path.getsize(OUT)
    print("indexed %d pages for v%s (%d with real prose), %.0f KB"
          % (len(pages), index["version"], with_text, size / 1024.0))
    if with_text < len(pages) * 0.5:
        sys.exit("fewer than half the pages produced usable prose; the extractor is wrong")


if __name__ == "__main__":
    main()
