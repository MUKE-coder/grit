#!/usr/bin/env python3
"""Fail when the docs site's community plugin list has drifted from the CLI's.

Two copies of a list with nothing comparing them is how the marketplace page
ends up advertising a plugin `grit plugin list` has never heard of. This is
the thing comparing them.
"""

import json
import pathlib
import sys

HERE = pathlib.Path(__file__).resolve().parent
REPO = HERE.parent.parent
SOURCE = REPO / "internal" / "plugin" / "community.json"
TARGET = REPO / "docs" / "config" / "community-plugins.json"


def main() -> None:
    if not SOURCE.exists():
        sys.exit(f"missing {SOURCE}")
    if not TARGET.exists():
        sys.exit(
            f"missing {TARGET}\n\nRun: python docs/scripts/sync-community-plugins.py"
        )

    source = SOURCE.read_text(encoding="utf-8")
    target = TARGET.read_text(encoding="utf-8")

    if source == target:
        count = len(json.loads(source).get("plugins", []))
        print(f"community plugin index is in sync ({count} plugin(s))")
        return

    # Say which way they differ, because "they differ" sends people to diff
    # two files by hand to find out that one has an extra entry.
    try:
        src_names = {p["name"] for p in json.loads(source).get("plugins", [])}
        tgt_names = {p["name"] for p in json.loads(target).get("plugins", [])}
    except (json.JSONDecodeError, KeyError, TypeError) as exc:
        sys.exit(f"one of the two indexes does not parse: {exc}")

    only_source = sorted(src_names - tgt_names)
    only_target = sorted(tgt_names - src_names)

    lines = ["the community plugin index and its docs copy have drifted apart"]
    if only_source:
        lines.append(f"  in the CLI but not the docs: {', '.join(only_source)}")
    if only_target:
        lines.append(f"  in the docs but not the CLI: {', '.join(only_target)}")
    if not only_source and not only_target:
        lines.append("  the same plugins, but the files differ (formatting or a field)")
    lines.append("")
    lines.append("Run: python docs/scripts/sync-community-plugins.py")
    sys.exit("\n".join(lines))


if __name__ == "__main__":
    main()
