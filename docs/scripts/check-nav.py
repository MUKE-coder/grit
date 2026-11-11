#!/usr/bin/env python3
"""Check the site header: the docs nav stays in the docs, and lights up once.

The header carries two sets of links, and which one it shows depends on
whether the reader is under /docs. That is deliberate. What is not is a link
in the docs set that points outside /docs: clicking it swaps the whole header
to the product nav, which reads as a second navbar appearing from nowhere.
"Blueprints" did exactly that and was reported as two different navbars.

Two things are checked, and the second is the one that bites quietly:

1. Every link in the in-docs set starts with /docs.

2. Every entry in that set, apart from Guides, is in DOCS_SECTIONS. Guides is
   "every /docs page that is not one of the others", so a section missing from
   that list underlines Guides as well as its own entry, on every page it owns.

Run: python docs/scripts/check-nav.py
"""

import pathlib
import re
import sys

HEADER = pathlib.Path(__file__).resolve().parent.parent / "components" / "site-header.tsx"


def main() -> None:
    src = HEADER.read_text(encoding="utf-8")
    problems: list[str] = []

    sections_block = re.search(r"const DOCS_SECTIONS = \[(.*?)\]", src, re.S)
    if not sections_block:
        sys.exit("DOCS_SECTIONS is not where this script expects it; has the header been rewritten?")
    sections = set(re.findall(r"'([^']+)'", sections_block.group(1)))

    # The in-docs set is the first arm of the `inDocs ? [...] : [...]`.
    arms = re.search(r"\{\(inDocs\s*\?\s*\[(.*?)\]\s*:\s*\[(.*?)\]\s*\)", src, re.S)
    if not arms:
        sys.exit("could not find the two link sets; has the header been rewritten?")

    in_docs = re.findall(r"\{ label: '([^']+)', href: '([^']+)'", arms.group(1))
    if not in_docs:
        sys.exit("the in-docs nav parsed as empty, which cannot be right")

    for label, href in in_docs:
        if not href.startswith("/docs"):
            problems.append(
                f"the docs nav's {label!r} points at {href!r}, which leaves /docs. "
                f"Clicking it swaps the header to the product nav, which is what "
                f"somebody reports as two different navbars"
            )
        # Guides is the catch-all and is deliberately absent from DOCS_SECTIONS.
        if label == "Guides":
            continue
        if not any(href.startswith(p) for p in sections):
            problems.append(
                f"{label!r} ({href}) is in the docs nav and not in DOCS_SECTIONS, "
                f"so its pages underline Guides as well as {label!r}"
            )

    # The other direction: a section listed for a nav entry that no longer exists.
    for prefix in sorted(sections):
        if not any(href.startswith(prefix) for _, href in in_docs):
            problems.append(
                f"DOCS_SECTIONS lists {prefix!r}, which no nav entry points at. "
                f"Guides will not light up on those pages and nothing else will either"
            )

    if problems:
        print("site header problems:\n")
        for p in problems:
            print(f"  - {p}")
        sys.exit(1)

    print(f"site header is consistent: {len(in_docs)} docs links, all inside /docs, all accounted for")


if __name__ == "__main__":
    main()
