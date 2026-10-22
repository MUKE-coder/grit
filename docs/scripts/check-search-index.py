#!/usr/bin/env python3
"""Fail when the search index does not match the pages on disk.

The index used to be a hand-written array in components/search-dialog.tsx. It
ended up covering 94 of 224 pages, and nobody knew, because a search that
cannot find a page looks exactly like a search for something that is not
there. Someone searched "fields", got three course lessons, and concluded the
search was bad. It was: the Field Types reference was not in the index at all.

So the index is generated, and this asserts it is current. A page added without
re-running the generator fails here instead of being quietly unsearchable.

    python scripts/check-search-index.py

Nothing is written. On a mismatch it prints what is missing and exits 1; run
scripts/search-index.py to fix it.
"""
import io
import json
import os
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
INDEX = os.path.join(DOCS, "config", "search-index.generated.json")
GENERATOR = os.path.join(HERE, "search-index.py")


def main():
    before = io.open(INDEX, encoding="utf-8").read()

    # The generator is the only description of what the index should hold, so
    # the check runs it rather than reimplementing the walk. Restored either
    # way, so a failing check never leaves the tree modified.
    result = subprocess.run(
        [sys.executable, "-I", GENERATOR],
        capture_output=True,
        text=True,
        cwd=DOCS,
    )
    after = io.open(INDEX, encoding="utf-8").read()
    io.open(INDEX, "w", encoding="utf-8", newline="").write(before)

    if result.returncode != 0:
        print(result.stdout, end="")
        print(result.stderr, end="", file=sys.stderr)
        print("\nscripts/search-index.py failed, so the index cannot be checked.")
        return 1

    if before == after:
        print("search index is current (%d pages)" % len(json.loads(before)))
        return 0

    was = {i["href"]: i for i in json.loads(before)}
    now = {i["href"]: i for i in json.loads(after)}

    print("The search index is stale. Run:\n")
    print("    python scripts/search-index.py\n")

    missing = sorted(set(now) - set(was))
    if missing:
        print("These pages exist and are not in the index, so nothing can find them:")
        for href in missing:
            print("    %s  (%s)" % (href, now[href]["title"]))

    gone = sorted(set(was) - set(now))
    if gone:
        print("\nIn the index with no page behind them:")
        for href in gone:
            print("    %s" % href)

    changed = sorted(h for h in set(was) & set(now) if was[h] != now[h])
    if changed:
        print("\nTitle or keywords moved on:")
        for href in changed:
            print("    %s" % href)
            if was[href]["title"] != now[href]["title"]:
                print("        %r -> %r" % (was[href]["title"], now[href]["title"]))

    return 1


if __name__ == "__main__":
    sys.exit(main())
