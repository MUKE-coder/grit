"""The stability table and the feature pages, checked against each other.

Every subsystem on /docs/stability names the page that documents it, and that
page carries a <MaturityBadge> reading the same id. Two files, one fact, which
is the shape of bug this project keeps finding: a test over each side passes
while the pair is broken. So this reads both in one pass.

Three ways it can go wrong, all of them silent:

  - an area names a route that no longer exists, so the badge links nowhere
  - a page loses its badge in an edit, and goes back to saying nothing about
    its own maturity
  - a badge names an id the config does not have, which renders nothing at all
    because MaturityBadge returns null rather than throwing in a reader's face

Run: python docs/scripts/check-maturity-badges.py
"""
import io
import os
import re
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
CONFIG = os.path.join(DOCS, 'config', 'stability.ts')
APP = os.path.join(DOCS, 'app')

# A floor, so a config this stops being able to read fails loudly rather than
# reporting that everything is fine.
MINIMUM_AREAS = 18


def areas():
    """id -> page, as config/stability.ts states it."""
    src = io.open(CONFIG, encoding='utf-8').read()
    ids = re.findall(r"^\s*id: '([a-z-]+)',", src, re.M)
    pairs = dict(re.findall(r"id: '([a-z-]+)',\s*\n\s*page: '([^']+)',", src))
    return ids, pairs


def badges_on(path):
    """The ids a page's <MaturityBadge> names."""
    src = io.open(path, encoding='utf-8').read()
    found = []
    for call in re.findall(r'<MaturityBadge\s+areas=\{\[([^\]]*)\]\}', src):
        found.extend(re.findall(r"'([a-z-]+)'", call))
    return found


def main():
    ids, pages = areas()
    if len(ids) < MINIMUM_AREAS:
        print('read only %d areas out of config/stability.ts, which means this '
              'check is not reading the config rather than that areas were '
              'removed' % len(ids))
        return 2

    problems = []

    # Every area that names a page: the page exists and carries the badge.
    for area_id, route in sorted(pages.items()):
        path = os.path.join(APP, route.lstrip('/').replace('/', os.sep), 'page.tsx')
        if not os.path.isfile(path):
            problems.append('%s names %s, which has no page.tsx' % (area_id, route))
            continue
        if area_id not in badges_on(path):
            problems.append(
                '%s names %s, and that page has no <MaturityBadge> for it. '
                'Add <MaturityBadge areas={[\'%s\']} /> under its heading, or '
                'drop the page field.' % (area_id, route, area_id))

    # And every badge anywhere names an area that exists.
    placed = 0
    for dirpath, dirs, names in os.walk(APP):
        dirs[:] = [d for d in dirs if d not in ('node_modules', '.next')]
        for name in names:
            if not name.endswith('.tsx'):
                continue
            path = os.path.join(dirpath, name)
            for area_id in badges_on(path):
                placed += 1
                if area_id not in ids:
                    rel = os.path.relpath(path, DOCS).replace('\\', '/')
                    problems.append(
                        '%s: a badge names %r, which is not an area in '
                        'config/stability.ts. It renders nothing.'
                        % (rel, area_id))

    print('%d areas, %d of them with a page, %d badges placed'
          % (len(ids), len(pages), placed))

    # An area with no page is allowed: not every subsystem has one page. But
    # the count is printed, because a quiet drift downward is the way this
    # stops covering anything.
    unpaged = [i for i in ids if i not in pages]
    if unpaged:
        print('no page named: %s' % ', '.join(unpaged))

    if problems:
        print()
        for problem in problems:
            print('  %s' % problem)
        return 1

    print('every area with a page has its badge, and every badge names a real area')
    return 0


if __name__ == '__main__':
    sys.exit(main())
