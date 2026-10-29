"""Em dashes, counted per area, with a budget that only goes down.

Banned in this project's user-facing writing. There were 3,457 when counting
started, so a check that simply fails is a check somebody disables on the
first run. This one holds a budget per area and fails when an area goes UP,
which makes the rule enforceable today and the backlog shrinkable.

Lower a number when you have cleared some. Never raise one.

Two exclusions, both deliberate:

  the changelog      a record of what was written at the time. Rewriting
                     past entries is rewriting history, and it is 677 of the
                     total.
  a bare placeholder a lone em dash in a table cell, as in {value || "-"},
                     means "no value". That is a typographic dash, not
                     punctuation, and replacing it with a colon would be
                     nonsense.

Why this cannot be a find and replace: an em dash does several jobs. An
appositive wants a colon, a consequence wants a comma and a conjunction, an
independent clause wants a full stop, and a PAIR of them wants commas or
brackets and is ruined by a colon. Each one is a sentence someone has to
read.

Run: python docs/scripts/check-em-dashes.py
"""
import io
import os
import re
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
APP = os.path.join(DOCS, 'app')

EM = '—'
SKIP_DIRS = {'node_modules', '.next', 'out', '.git'}

# A lone em dash standing in for "no value".
PLACEHOLDER = re.compile(r'["\'>]\s*' + EM + r'\s*["\'<]')

# area -> the most em dashes it may contain. Only ever lowered.
BUDGET = {
    'courses': 1281,
    'docs/prerequisites': 179,
    'docs/security': 168,
    'docs/desktop': 123,
    'docs/backend': 102,
    'docs/concepts': 100,
    'compare': 88,
    'docs/getting-started': 88,
    'docs/ai-skill': 77,
    'docs/mobile': 54,
    'docs/testing': 48,
    'docs/admin': 45,
    'course': 40,
    'docs/stack-selector': 37,
    'docs/tech-kits': 37,
    'docs/learnings': 33,
    'docs/tutorials': 32,
    'docs/ai-workflows': 27,
    'docs/plugins': 23,
    'docs/batteries': 20,
    'docs/deployment': 19,
    'docs/page.tsx': 17,
    'docs/demo': 15,
    'docs/frontend': 12,
    'docs/start': 12,
    '(root)': 11,
    'docs/infrastructure': 7,
    'docs/versioning': 5,
    'sponsors': 5,
    'builder': 4,
    'docs/benchmarks': 4,
    'docs/design': 4,
    'docs/governance': 3,
    'sponsor': 3,
    'api': 2,
    'hire': 2,
    'problems': 2,
    'blueprints': 1,
    'docs/blueprints': 1,
    'docs/extending': 1,
    'donate': 1,
    'showcase': 1,
    'zenith': 1,
    'OTHER': 0,
}

EXCLUDED_AREA = 'docs/changelog'


def area_of(rel):
    bits = rel.split('/')
    if bits[0] == 'docs' and len(bits) > 1:
        return 'docs/' + bits[1]
    return bits[0] if len(bits) > 1 else '(root)'


def main():
    counts = {}
    placeholders = 0

    for dirpath, dirs, names in os.walk(APP):
        dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
        for name in names:
            if not name.endswith(('.tsx', '.ts', '.md')):
                continue
            path = os.path.join(dirpath, name)
            rel = os.path.relpath(path, APP).replace('\\', '/')
            area = area_of(rel)
            if area == EXCLUDED_AREA:
                continue
            try:
                text = io.open(path, encoding='utf-8').read()
            except (OSError, UnicodeDecodeError):
                continue

            total = text.count(EM)
            if not total:
                continue
            bare = len(PLACEHOLDER.findall(text))
            placeholders += bare
            real = total - bare
            if real > 0:
                counts[area] = counts.get(area, 0) + real

    over = []
    for area, n in sorted(counts.items(), key=lambda kv: -kv[1]):
        budget = BUDGET.get(area, BUDGET['OTHER'] if area not in BUDGET else 0)
        if area not in BUDGET:
            budget = BUDGET['OTHER']
        flag = ''
        if n > budget:
            flag = '  OVER BUDGET (%d)' % budget
            over.append((area, n, budget))
        print('%5d  %-24s%s' % (n, area, flag))

    print()
    print('%d em dashes outside the changelog, %d bare placeholders ignored'
          % (sum(counts.values()), placeholders))

    if over:
        print()
        for area, n, budget in over:
            print('%s has %d, budget %d. Em dashes are banned in user-facing '
                  'writing here: rewrite the sentence rather than raising the '
                  'number.' % (area, n, budget))
        return 1

    print('no area is over its budget')
    return 0


if __name__ == '__main__':
    sys.exit(main())
