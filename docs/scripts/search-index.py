#!/usr/bin/env python3
"""Build the docs site's search index from the pages that actually exist.

The index was a hand-written list in components/search-dialog.tsx, and it
drifted the way every hand-written list of a growing thing drifts: 181 doc
routes on disk, 94 of them in the index. Searching "fields" could not return
the Field Types page because the Field Types page was not in the index at all,
so course lessons with the word in their title were the best it could do.

So the list is generated. Every route under app/docs is in it, titled from
config/docs-metadata.ts, with that page's description as keywords. The curated
keyword blobs from the old list are kept and merged by href, because "admin
table list sort filter search pagination columns rows bulk actions" is real
actions" is real knowledge about what people type, and a description cannot
replace it.

Run it after adding a doc page:

    python scripts/search-index.py

It rewrites the index and prints what changed. scripts/check-search-index.py
fails when it is stale, and CI runs that, so a page added without running this
is caught rather than left silently unsearchable.
"""
import glob
import json
import os
import re
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
APP_DOCS = os.path.join(DOCS, "app", "docs")
APP_COURSES = os.path.join(DOCS, "app", "courses")
METADATA = os.path.join(DOCS, "config", "docs-metadata.ts")
CURATED = os.path.join(DOCS, "config", "search-keywords.json")
OUT = os.path.join(DOCS, "config", "search-index.generated.ts")

# A route's section, by the directory it lives in. The names are the ones the
# old hand-written list used, because they are what the sidebar says.
SECTIONS = {
    "": "Getting Started",
    "getting-started": "Getting Started",
    "concepts": "Core Concepts",
    "backend": "Backend (Go API)",
    "frontend": "Frontend (Next.js)",
    "admin": "Admin Panel",
    "batteries": "Batteries",
    "infrastructure": "Infrastructure",
    "deployment": "Deployment",
    "security": "Security & Testing",
    "testing": "Security & Testing",
    "stability": "Security & Testing",
    "governance": "Security & Testing",
    "plugins": "Plugins",
    "tutorials": "Tutorials",
    "tech-kits": "Tech Kits",
    "design": "Design System",
    "desktop": "Desktop (Wails)",
    "mobile": "Mobile (Expo)",
    "ai-integration": "AI Workflows",
    "ai-workflows": "AI Workflows",
    "ai-skill": "For AI Assistants",
    "systems": "Systems (Design)",
    "learnings": "Learnings",
    "demo": "Demo",
    "scaling": "Scaling",
    "upgrades": "Upgrades",
    "versioning": "Upgrades",
    "changelog": "Changelog",
    "benchmarks": "Benchmarks",
    "cli": "Core Concepts",
    "start": "Getting Started",
    "prerequisites": "Getting Started",
    "stack-selector": "Getting Started",
    "api": "Backend (Go API)",
    "courses": "Courses",
}


def routes():
    """Every static /docs route, from the pages on disk."""
    out = []
    for dirpath, _dirs, files in os.walk(APP_DOCS):
        if "page.tsx" not in files:
            continue
        rel = os.path.relpath(dirpath, APP_DOCS).replace(os.sep, "/")
        if "[" in rel:  # one directory, many pages: see dynamic() below
            continue
        href = "/docs" if rel == "." else "/docs/" + rel
        out.append(href)
    return sorted(out)


# The three sections rendered from a config file through one [param] template,
# and the number of entries each holds. Walking app/docs finds one directory
# for each and no pages at all, which is how 43 pages, the whole of Systems
# included, were missing from search and from the sitemap: "studio" could not
# find GORM Studio because the page had no file of its own.
#
# The count is asserted. If one of these files changes shape, this fails and
# says so, rather than quietly indexing nothing: being invisible is the failure
# mode that nobody notices, because a search that finds too little looks
# exactly like a search for something that isn't there.
DYNAMIC = (
    ("/docs/systems/", "Systems (Design)", "systems-*.ts", 29),
    ("/docs/deployment/", "Deployment", "deployment-providers.ts", 8),
    ("/docs/benchmarks/", "Benchmarks", "benchmarks.ts", 6),
)

# slug, then the name and one-liner that follow it in the same object. All
# three configs happen to agree on these field names, which is why one reader
# serves them.
ENTRY = re.compile(
    r"slug:\s*'([a-z0-9-]+)',"
    r"(?:(?!\bslug:).)*?\bname:\s*'((?:[^'\\]|\\.)*)',"
    r"(?:(?!\bslug:).)*?\btagline:\s*\n?\s*'((?:[^'\\]|\\.)*)'",
    re.S,
)


def dynamic():
    """Every page of the three config-driven sections."""
    out = []
    for prefix, section, pattern, expected in DYNAMIC:
        found = []
        for path in sorted(glob.glob(os.path.join(DOCS, "config", pattern))):
            if path.endswith(("systems-types.ts", "systems.ts")):
                continue
            src = open(path, encoding="utf-8").read()
            for slug, name, tagline in ENTRY.findall(src):
                found.append((prefix + slug, name.replace("\\'", "'"),
                              re.sub(r"\s+", " ", tagline.replace("\\'", "'"))))
        if len(found) != expected:
            raise SystemExit(
                "%s: read %d entries from %s, expected %d. The config changed "
                "shape; fix the reader in scripts/search-index.py, or these "
                "pages go missing from search without anything saying so."
                % (prefix, len(found), pattern, expected)
            )
        out.extend((href, name, tagline, section) for href, name, tagline in found)
    return out


# The 41 standalone course pages under app/courses: Stripe payments, GORM
# mastery, the Grit Web / Desktop / Mobile series. Real pages with their own
# files, and none of them was in search, because routes() walks app/docs and
# stops there. A search for "stripe" found the plugin page and not the course
# about integrating it.
#
# The courses in config/courses.ts are deliberately NOT here. search-dialog.tsx
# builds those live from the config, with the module title, the minutes and the
# difficulty in their keywords, and its map keeps the later entry: adding them
# here would ship 300 entries to the browser that are then thrown away.
PAGE_METADATA = re.compile(
    r"export const metadata(?::\s*Metadata)?\s*=\s*\{"
    r".*?\btitle:\s*(['\"])(.*?)\1"
    r".*?\bdescription:\s*(['\"])(.*?)\3",
    re.S,
)
STANDALONE_COURSES_MIN = 40


def standalone_courses():
    """Every course page that has a file of its own, titled from its metadata."""
    out = []
    for dirpath, _dirs, files in os.walk(APP_COURSES):
        if "page.tsx" not in files:
            continue
        rel = os.path.relpath(dirpath, APP_COURSES).replace(os.sep, "/")
        if rel == "." or "[" in rel:   # /courses itself, and the [param] tree
            continue
        href = "/courses/" + rel
        src = open(os.path.join(dirpath, "page.tsx"), encoding="utf-8").read()
        m = PAGE_METADATA.search(src)
        if not m:
            raise SystemExit(
                "%s: no metadata title and description, so there is nothing to "
                "index it by. Add them to the page, or teach the reader in "
                "scripts/search-index.py where its title lives." % href
            )
        title = re.sub(r"\s+", " ", m.group(2)).strip()
        description = re.sub(r"\s+", " ", m.group(4)).strip()
        out.append((href, title, description, "Courses"))
    if len(out) < STANDALONE_COURSES_MIN:
        raise SystemExit(
            "standalone courses: found %d pages under app/courses, expected at "
            "least %d. Either pages were deleted, or the walk broke and they go "
            "missing from search without anything saying so."
            % (len(out), STANDALONE_COURSES_MIN)
        )
    return sorted(out)


def metadata():
    """title and description per route, from config/docs-metadata.ts."""
    src = open(METADATA, encoding="utf-8").read()
    entries = {}
    # '/docs/x': { title: '...', description: '...' } across several lines.
    pattern = re.compile(
        r"'(/docs[^']*)':\s*\{\s*title:\s*(['\"])(.*?)\2,\s*description:\s*(['\"])(.*?)\4",
        re.S,
    )
    for m in pattern.finditer(src):
        href, title, desc = m.group(1), m.group(3), m.group(5)
        entries[href] = (title.strip(), re.sub(r"\s+", " ", desc).strip())
    return entries


def curated():
    """Hand-maintained keywords, by href.

    The old list held two things at once: which pages exist, which it got
    wrong, and what people type when hunting for one, which it got right. The
    first is generated now. The second cannot be, because "bulk create grid
    spreadsheet excel paste fill down" is knowledge about the people searching
    and no page description contains it. So it is kept by hand, here, and the
    href is the only thing the two halves need to agree on.
    """
    return json.load(open(CURATED, encoding="utf-8"))


def section_for(href):
    rest = href[len("/docs"):].strip("/")
    top = rest.split("/")[0] if rest else ""
    return SECTIONS.get(top, "Documentation")


def title_from_href(href):
    """A readable title for a page with no metadata entry."""
    rest = href[len("/docs"):].strip("/")
    leaf = rest.split("/")[-1] if rest else "Introduction"
    words = [w for w in re.split(r"[-_]", leaf) if w]
    return " ".join(w[:1].upper() + w[1:] for w in words) or "Introduction"


# Words a description contributes nothing by repeating.
STOP = set("""a an the and or of for to in on with without by from as at is are
be it its this that these those you your we our can will into over under
use used using get set new how what when where which why""".split())


def keywords_from(description, title, href):
    """Searchable words for a page: its description, path and title."""
    words = []
    for source in (description, href.replace("/", " ").replace("-", " "), title):
        for w in re.split(r"[^a-z0-9+]+", source.lower()):
            if len(w) > 1 and w not in STOP and w not in words:
                words.append(w)
    return " ".join(words[:60])


def main():
    meta = metadata()
    hand = curated()
    items = []
    missing_meta = []

    def add(href, title, description, section):
        entry = hand.get(href, {})
        words = keywords_from(description, title, href)
        if entry.get("keywords"):
            # The curated blob first: it is what people actually type.
            words = entry["keywords"] + " " + words
        # Dedupe, keeping order.
        seen, ordered = set(), []
        for w in words.split():
            if w not in seen:
                seen.add(w)
                ordered.append(w)
        items.append(
            {
                "title": title,
                "href": href,
                "section": section,
                "keywords": " ".join(ordered),
            }
        )

    for href in routes():
        entry = hand.get(href, {})
        title, description = meta.get(href, (None, ""))
        if title is None:
            missing_meta.append(href)
            # Only now is the old list's title worth having: the page has no
            # metadata entry, so there is nothing better. Where metadata does
            # have one it wins, because the old list had File Storage filed
            # under "CSV / Excel Export" and a search for "upload" could not
            # reach the page about uploads.
            title = entry.get("title") or title_from_href(href)
        add(href, title, description, entry.get("section") or section_for(href))

    shadowed = []
    seen_href = {i["href"] for i in items}
    for href, name, tagline, section in dynamic() + standalone_courses():
        if href in seen_href:
            # A static page of the same path, which is what Next serves: the
            # [param] template never renders this slug. One entry, not two.
            shadowed.append(href)
            continue
        add(href, name, tagline, section)

    # Anchors and non-/docs entries the hand-written list carried, which are
    # not pages on disk and would otherwise be dropped.
    extra = [v | {"href": h} for h, v in hand.items() if h not in {i["href"] for i in items}]
    for e in sorted(extra, key=lambda e: e["href"]):
        items.append({"title": e["title"], "href": e["href"], "section": e["section"], "keywords": e["keywords"]})

    # The entries are JSON, and the TypeScript beside them is a typed import of
    # that file. One list, readable by the test and by anything else that wants
    # it, with no reader for TypeScript object literals: the afternoon spent
    # writing one of those ended on a page title containing a quotation mark.
    data = json.dumps(items, indent=2, ensure_ascii=False) + "\n"
    open(OUT.replace(".ts", ".json"), "w", encoding="utf-8", newline="").write(data)

    out = '''// Generated by scripts/search-index.py. Do not edit by hand.
//
// Every static route under app/docs, titled from config/docs-metadata.ts, with
// that page's description, path and curated keywords as the searchable text.
// Regenerate after adding a page:
//
//     python scripts/search-index.py
//
// The entries live in search-index.generated.json beside this file, so a test
// or a script can read them without parsing TypeScript.
// scripts/check-search-index.py fails when this is stale, and CI runs it: the
// hand-written list this replaced had drifted to 94 of 224 pages, and nothing
// anywhere said so.

import entries from './search-index.generated.json'

export interface SearchItem {
  title: string
  href: string
  section: string
  keywords: string
}

export const SEARCH_INDEX: SearchItem[] = entries
'''
    before = ""
    if os.path.exists(OUT):
        before = open(OUT, encoding="utf-8").read()
    open(OUT, "w", encoding="utf-8", newline="").write(out)

    print("indexed %d pages (%d static routes, %d config-driven, %d course "
          "pages, %d curated kept)"
          % (len(items), len(routes()), len(dynamic()),
             len(standalone_courses()), len(hand)))
    if shadowed:
        print("a static page of the same path wins, so the config entry is skipped:")
        for href in shadowed:
            print("   ", href)
    if missing_meta:
        print("no entry in config/docs-metadata.ts, titled from the path:")
        for href in missing_meta:
            print("   ", href)
    if before and before != out:
        print("search-index.generated.ts changed")
    return 0


if __name__ == "__main__":
    sys.exit(main())
