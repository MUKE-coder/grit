#!/usr/bin/env python3
"""Every capability row's stability pointer names a real stability area.

Two files have to agree:

  config/capabilities.ts    a row may carry  stability: '<id>'
  config/stability.ts       declares the areas, each with an id

The matrix shows the area's badge beside the row label, so a row that depends
on a beta subsystem reads beta in both places. A pointer that resolves to
nothing is silent: stabilityOf returns undefined, StabilityBadge renders null,
and the row quietly goes back to looking settled. Which is the bug this pair
was built to close (grit#166), reappearing by a rename.

This is the shape that keeps recurring here: two sources that must agree, with
a test per side that passes while the pair is broken. So the check reads both.

Usage: python3 scripts/check-stability-pointers.py      (from docs/)
"""
import io
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)

CAPS = os.path.join(DOCS, 'config', 'capabilities.ts')
STAB = os.path.join(DOCS, 'config', 'stability.ts')


def read(path):
    return io.open(path, encoding='utf-8').read()


def main():
    caps = read(CAPS)
    stab = read(STAB)

    areas = re.findall(r"^    id: '([^']+)',$", stab, re.M)
    if not areas:
        print('FAIL: config/stability.ts declares no areas; has it been renamed?')
        return 1
    if len(areas) != len(set(areas)):
        dupes = sorted({a for a in areas if areas.count(a) > 1})
        print('FAIL: duplicate stability area ids: %s' % ', '.join(dupes))
        return 1

    # A row is  id: 'x',  optionally followed by  stability: 'y',
    rows = re.findall(
        r"^        id: '([^']+)',\n(?:        stability: '([^']+)',\n)?",
        caps,
        re.M,
    )
    if not rows:
        print('FAIL: config/capabilities.ts declares no rows; has it been renamed?')
        return 1

    pointed = [(row, area) for row, area in rows if area]
    if not pointed:
        print('FAIL: no capability row points at a stability area. There were '
              'eleven when this check was written, so either the pointers were '
              'dropped or this script no longer reads the file it thinks it does.')
        return 1

    bad = [(row, area) for row, area in pointed if area not in areas]
    if bad:
        print('FAIL: %d capability row(s) point at a stability area that does '
              'not exist.\n' % len(bad))
        for row, area in bad:
            print("  capabilities.ts row %-14s -> stability.ts id '%s' (missing)" % (row, area))
        print('\nThe badge on those rows renders nothing, so the matrix reads as '
              'settled\nwhile the stability page may not. Known area ids:\n')
        print('  ' + ', '.join(sorted(areas)))
        return 1

    print('OK: %d capability rows, %d of them pointing at %d stability areas'
          % (len(rows), len(pointed), len(areas)))
    return 0


if __name__ == '__main__':
    sys.exit(main())
