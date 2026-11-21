#!/usr/bin/env python3
"""Every internal link and image in a guide points at something that exists.

A guide's job is to send a stuck reader somewhere useful, so a link that 404s
costs more here than almost anywhere else on the site. They are also easy to
get wrong: /docs/cli/doctor reads like a route and is not one, because the CLI
reference is a single page with an anchor per command.

Scoped to app/guides. Widening it to the whole site would be useful and would
also surface years of links nobody is fixing today; this checks the pages
being written now, and the walk is three lines from covering more.

A link is good when the App Router would serve it:

  /guides/x        app/guides/x/page.tsx
  /docs/a/b        app/docs/a/b/page.tsx
  /docs/cli#doctor the page part resolves; the fragment is not checked,
                   because anchors are generated from config in several
                   places and a checker that guessed would be wrong

Run: python docs/scripts/check-guide-links.py
"""
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
APP = os.path.join(DOCS, "app")
GUIDES = os.path.join(APP, "guides")

# href="/..." and href={`/...`} and href: '/...'
HREF = re.compile(r"""href[=:]\s*[{]?\s*['"`](/[^'"`#?]*)""")

# A screenshot that 404s is as bad as a link that does, and quieter: the
# page still renders, with a gap where the evidence was.
SRC = re.compile(r"""src[=:]\s*[{]?\s*['"`](/[^'"`#?]*)""")

# Routes the file tree cannot show, because they are produced by a route
# handler or a generateStaticParams rather than a page.tsx.
SPECIAL = {"/", "/llms.txt", "/llms-full.txt", "/robots.txt", "/sitemap.xml"}


def route_exists(path):
    rel = path.strip("/")
    if not rel:
        return True
    candidate = os.path.join(APP, *rel.split("/"))
    if os.path.isfile(os.path.join(candidate, "page.tsx")):
        return True
    # A dynamic segment: app/problems/[id]/page.tsx serves /problems/d1.
    parent = os.path.dirname(candidate)
    if os.path.isdir(parent):
        for name in os.listdir(parent):
            if name.startswith("[") and os.path.isfile(
                os.path.join(parent, name, "page.tsx")
            ):
                return True
    return False


def main():
    if not os.path.isdir(GUIDES):
        print("no app/guides yet, nothing to check")
        return 0

    problems = []
    checked = 0
    for dirpath, _dirs, files in os.walk(GUIDES):
        for name in files:
            if not name.endswith(".tsx"):
                continue
            path = os.path.join(dirpath, name)
            rel = os.path.relpath(path, DOCS).replace(os.sep, "/")
            with open(path, encoding="utf-8") as fh:
                for i, line in enumerate(fh, 1):
                    for m in HREF.finditer(line):
                        href = m.group(1)
                        if href in SPECIAL:
                            continue
                        # An interpolated href is built at render time, so
                        # there is no literal route to resolve. The index
                        # builds its cards that way, and config/guides.ts is
                        # checked against the pages on disk separately.
                        if "$" in href:
                            continue
                        checked += 1
                        if not route_exists(href):
                            problems.append((rel, i, href))
                    for m in SRC.finditer(line):
                        src = m.group(1)
                        if "$" in src:
                            continue
                        checked += 1
                        if not os.path.isfile(
                            os.path.join(DOCS, "public", *src.strip("/").split("/"))
                        ):
                            problems.append((rel, i, src + " (no such file in public/)"))

    if problems:
        print("%d link(s) in the guides point at a route that does not exist:\n"
              % len(problems))
        for rel, line, href in problems:
            print("  %s:%d  %s" % (rel, line, href))
        print("\nEither the page is missing, or the path is not a route: the CLI")
        print("reference, for one, is a single page with an anchor per command.")
        return 1

    print("every internal link in the guides resolves (%d checked)" % checked)
    return 0


if __name__ == "__main__":
    sys.exit(main())
