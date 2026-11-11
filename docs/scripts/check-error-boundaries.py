#!/usr/bin/env python3
"""Check the site has the three fallbacks every Grit project is given.

`grit new` writes app/error.tsx, app/global-error.tsx and app/not-found.tsx
into every frontend it scaffolds. This site, which is where people come to read
about that, shipped for a long time with only not-found.tsx: an exception
anywhere under the root layout fell through to Next's own fallback, one line of
unstyled black-on-white reading "Application error: a client-side exception has
occurred (see the browser console for more information)".

That is the worst possible answer on a framework's own site, because the
visitor cannot tell whether the framework is broken or the page is.

Run: python docs/scripts/check-error-boundaries.py
"""

import pathlib
import sys

APP = pathlib.Path(__file__).resolve().parent.parent / "app"

REQUIRED = {
    "error.tsx": "an exception in a page falls through to Next's unstyled fallback",
    "global-error.tsx": "an exception in the root layout has nothing to render at all",
    "not-found.tsx": "a 404 is Next's default rather than this site's",
}


def main() -> None:
    problems: list[str] = []

    for name, consequence in REQUIRED.items():
        path = APP / name
        if not path.exists():
            problems.append(f"app/{name} is missing, so {consequence}")
            continue
        if not path.read_text(encoding="utf-8").strip():
            problems.append(f"app/{name} is empty, so {consequence}")

    # global-error replaces the root layout, so the stylesheet that layout
    # imports is not loaded. A Tailwind class there renders unstyled, which is
    # the failure this file exists to prevent, in the one place it is hardest
    # to notice: it only renders when something has already gone very wrong.
    global_error = APP / "global-error.tsx"
    if global_error.exists():
        src = global_error.read_text(encoding="utf-8")
        if "className=" in src:
            problems.append(
                "app/global-error.tsx uses className. It replaces the root layout, so "
                "globals.css is never loaded and those classes do nothing: style it inline"
            )
        if "<html" not in src or "<body" not in src:
            problems.append(
                "app/global-error.tsx does not render <html> and <body>. It replaces the "
                "root layout, so it has to provide them itself or the page is blank"
            )

    if problems:
        print("error boundary problems:\n")
        for p in problems:
            print(f"  - {p}")
        sys.exit(1)

    print("the site has all three fallbacks, and global-error stands on its own")


if __name__ == "__main__":
    main()
