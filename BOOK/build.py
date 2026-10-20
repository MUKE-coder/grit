#!/usr/bin/env python3
"""Build the book and report progress against the 500 page target.

    python build.py            compile book.pdf and print the page count
    python build.py --pages 6  also render pages 1..6 as PNGs to look at

Page count comes from Typst's own PNG pagination rather than from counting
/Type /Page in the PDF: that regex double counts on this writer, which made an
early draft look twice as long as it was.
"""

import argparse
import os
import re
import shutil
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
TARGET_PAGES = 500


def run(cmd, **kw):
    return subprocess.run(cmd, cwd=HERE, capture_output=True, text=True, **kw)


def page_count():
    """Pages, counted by asking Typst to paginate into a throwaway directory."""
    tmp = tempfile.mkdtemp(prefix="bookpages")
    try:
        r = run(["typst", "compile", "book.typ",
                 os.path.join(tmp, "p-{n}.png"), "--format", "png", "--ppi", "40"])
        if r.returncode != 0:
            return None
        return len([f for f in os.listdir(tmp) if f.endswith(".png")])
    finally:
        shutil.rmtree(tmp, ignore_errors=True)


def chapter_report():
    """Lines of source per chapter, so a stub is obvious at a glance."""
    chdir = os.path.join(HERE, "chapters")
    if not os.path.isdir(chdir):
        return []
    rows = []
    for name in sorted(os.listdir(chdir)):
        if not name.endswith(".typ"):
            continue
        path = os.path.join(chdir, name)
        with open(path, encoding="utf-8") as fh:
            body = fh.read()
        lines = len(body.splitlines())
        # A chapter that is still a placeholder says so in one line.
        stub = "Placeholder." in body or lines < 25
        rows.append((name, lines, stub))
    return rows


def lint_prose():
    """Find bare # in prose, which Typst reads as the start of an expression.

    Written after 'C#' broke the build twice. Typst's own error points at the
    character after the hash and says 'unexpected comma', which does not lead
    anybody to the real problem quickly. Code blocks are exempt: Typst leaves
    their contents alone, so a shell comment is fine.
    """
    chdir = os.path.join(HERE, "chapters")
    if not os.path.isdir(chdir):
        return []
    known = ("import", "bridge", "gotcha", "project", "tryit", "assignment",
             "backmatter", "compare",
             "table", "link", "raw", "text", "set", "show", "let", "block",
             "grid", "align", "upper", "line", "counter", "outline", "include",
             "v(", "h(", "pagebreak")
    problems = []
    for name in sorted(os.listdir(chdir)):
        if not name.endswith(".typ"):
            continue
        in_code = False
        with open(os.path.join(chdir, name), encoding="utf-8") as fh:
            for n, line in enumerate(fh, 1):
                if line.lstrip().startswith("```"):
                    in_code = not in_code
                    continue
                if in_code:
                    continue
                for m in re.finditer(r"(?<!\\)#(\w*)", line):
                    word = m.group(1)
                    if any(word.startswith(k.rstrip("(")) for k in known):
                        continue
                    if line.count('"') >= 2 and m.start() > line.find('"'):
                        continue  # inside a string, e.g. rgb("#e3e6ea")
                    problems.append("%s:%d  %s" % (name, n, line.strip()[:60]))
    return problems


def lint_code_width():
    """Find code lines too wide for the measure they are set in.

    Body code is Consolas 8.8pt in a 5.35in column, which fits 80 characters, so 78 leaves a safety margin.
    Inside a callout the inset takes 27pt off, and the type drops to 8.2pt to
    compensate, so the budget is the same. Anything longer wraps mid-token and
    reads as nonsense, which is worse in a listing than small type.
    """
    chdir = os.path.join(HERE, "chapters")
    if not os.path.isdir(chdir):
        return []
    LIMIT = 75
    problems = []
    for name in sorted(os.listdir(chdir)):
        if not name.endswith(".typ"):
            continue
        in_code = False
        with open(os.path.join(chdir, name), encoding="utf-8") as fh:
            for n, line in enumerate(fh, 1):
                stripped = line.rstrip("\n")
                if stripped.lstrip().startswith("```"):
                    in_code = not in_code
                    continue
                if in_code and len(stripped) > LIMIT:
                    problems.append("%s:%d  %d cols  %s"
                                    % (name, n, len(stripped),
                                       stripped.strip()[:46]))
    return problems


def lint_markdown_bold():
    """Find Markdown's ** in prose. Typst's bold is a single star, so ** is a
    habit from writing Markdown and renders as an empty emphasis plus a stray
    star. Typst warns rather than failing, so it is easy to ship."""
    chdir = os.path.join(HERE, "chapters")
    if not os.path.isdir(chdir):
        return []
    out = []
    for name in sorted(os.listdir(chdir)):
        if not name.endswith(".typ"):
            continue
        in_code = False
        with open(os.path.join(chdir, name), encoding="utf-8") as fh:
            for n, line in enumerate(fh, 1):
                if line.lstrip().startswith("```"):
                    in_code = not in_code
                    continue
                if not in_code and "**" in line:
                    out.append("%s:%d  %s" % (name, n, line.strip()[:50]))
    return out


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--pages", type=int, default=0,
                    help="also render this many leading pages as PNGs")
    args = ap.parse_args()

    for p in lint_prose():
        print("  bare # in prose (escape it as \\#): " + p)
    for p in lint_markdown_bold():
        print("  Markdown ** in prose (Typst bold is *single*): " + p)
    wide = lint_code_width()
    for p in wide[:12]:
        print("  code too wide: " + p)
    if len(wide) > 12:
        print("  ... and %d more over 75 columns" % (len(wide) - 12))

    r = run(["typst", "compile", "book.typ", "book.pdf"])
    if r.returncode != 0:
        print(r.stdout)
        print(r.stderr, file=sys.stderr)
        return 1
    if r.stderr.strip():
        print(r.stderr.strip())

    size = os.path.getsize(os.path.join(HERE, "book.pdf"))
    pages = page_count()

    # Publish straight to the docs site. Two copies of a PDF in one repository
    # is two copies that drift, so the one under docs/public is the only one
    # committed and this is the only thing that writes it.
    published = os.path.join(
        HERE, "..", "docs", "public", "book", "go-for-javascript-developers.pdf")
    os.makedirs(os.path.dirname(published), exist_ok=True)
    shutil.copyfile(os.path.join(HERE, "book.pdf"), published)
    print("published  docs/public/book/go-for-javascript-developers.pdf")

    rows = chapter_report()
    written = [r for r in rows if not r[2]]
    stubs = [r for r in rows if r[2]]

    print("book.pdf   %d KB" % (size // 1024))
    if pages:
        pct = 100.0 * pages / TARGET_PAGES
        per = (pages / len(written)) if written else 0
        print("pages      %d of %d target (%.0f%%)" % (pages, TARGET_PAGES, pct))
        print("chapters   %d written, %d stubs" % (len(written), len(stubs)))
        if per:
            print("density    %.1f pages per written chapter" % per)
            need = TARGET_PAGES - pages
            print("remaining  about %d more chapters at this density"
                  % round(need / per))
    if stubs:
        print("stubs      " + ", ".join(s[0] for s in stubs))

    if args.pages:
        for f in os.listdir(HERE):
            if re.match(r"page-\d+\.png$", f):
                os.remove(os.path.join(HERE, f))
        run(["typst", "compile", "book.typ", "page-{n}.png",
             "--format", "png", "--ppi", "110", "--pages", "1-%d" % args.pages])
        print("rendered   page-01.png .. page-%02d.png" % args.pages)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
