#!/usr/bin/env python3
"""Copy the community plugin index into the docs site.

`internal/plugin/community.json` is the single source of truth: the CLI
embeds it for `grit plugin list`. The docs site cannot import a file from
outside its own directory, so it gets a copy, and this script makes the copy.

Run it after editing the index. `check-community-plugins.py` fails CI when
the two have drifted, which is the only thing that stops the marketplace page
from quietly describing a list nobody has anymore.
"""

import json
import pathlib
import sys

HERE = pathlib.Path(__file__).resolve().parent
REPO = HERE.parent.parent
SOURCE = REPO / "internal" / "plugin" / "community.json"
TARGET = REPO / "docs" / "config" / "community-plugins.json"


def read_source() -> str:
    if not SOURCE.exists():
        sys.exit(f"missing {SOURCE}")
    raw = SOURCE.read_text(encoding="utf-8")
    try:
        json.loads(raw)
    except json.JSONDecodeError as exc:
        sys.exit(f"{SOURCE} does not parse: {exc}")
    return raw


def main() -> None:
    raw = read_source()
    TARGET.parent.mkdir(parents=True, exist_ok=True)
    # newline="" so the file keeps LF on Windows, where open(..., "w")
    # otherwise writes CRLF and the check script then always disagrees.
    with open(TARGET, "w", encoding="utf-8", newline="") as handle:
        handle.write(raw)
    count = len(json.loads(raw).get("plugins", []))
    print(f"synced {count} community plugin(s) to {TARGET.relative_to(REPO)}")


if __name__ == "__main__":
    main()
