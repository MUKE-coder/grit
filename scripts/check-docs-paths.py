"""Every project path the docs name, checked against a real project.

A page that tells you to open apps/api/internal/foo/bar.go, when the
scaffold stopped writing that file three releases ago, sends the reader
looking for something that is not there. They do not conclude the page is
wrong; they conclude they scaffolded it incorrectly.

Checked against a GENERATED PROJECT, never the templates. The templates are
what Grit intends to write and a project is what it wrote, and the gap
between the two is exactly the kind of thing this is looking for.

Set GRIT_REFERENCE_PROJECT to a scaffolded project. The fullest one you have
is the best choice: a path missing only because this project has no Expo app
is noise, so anything under an app the reference lacks is skipped.

Run: GRIT_REFERENCE_PROJECT=/path/to/project python scripts/check-docs-paths.py

The nightly canary runs this against the --full project it scaffolds, which
is the only place a real generated project exists in CI.
"""
import io
import os
import re
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
DOCS = os.path.join(ROOT, 'docs', 'app')

REFERENCE = os.environ.get('GRIT_REFERENCE_PROJECT', '')

# A path inside a scaffolded project: apps/api/internal/x/y.go, packages/
# shared/types/index.ts, and so on. Only the roots a project actually has,
# so prose like "internal state" never matches.
PATH = re.compile(
    r'\b((?:apps/(?:api|web|admin|expo|desktop)|packages/shared|cmd|e2e)'
    r'/[A-Za-z0-9_./\[\]-]+\.(?:go|tsx|ts|json|css|sql|yaml|yml|md))')

# The changelog records what was true then.
SKIP_PARTS = ('app/docs/changelog',)
SKIP_DIRS = {'node_modules', '.next', 'out', '.git'}

# A path with one of these in it is a shape, not a file.
PLACEHOLDER = re.compile(r'\{|\}|<|>|\.\.\.|\bName\b|\bResource\b|xxx', re.I)

# Directories whose contents come from `grit generate resource`, not from the
# scaffold. A page about generating a Post naming models/post.go is correct;
# the file is simply absent from a project where nobody generated one.
# Checking those buries the signal worth having: a framework file that moved.
GENERATED_DIRS = (
    'apps/api/internal/models/',
    'apps/api/internal/handlers/',
    'apps/api/internal/services/',
    'apps/api/internal/routes/',
    'apps/admin/resources/',
    'apps/admin/app/',
    'apps/web/app/',
    'apps/expo/app/',
    'packages/shared/schemas/',
    'packages/shared/types/',
    'e2e/',
    'apps/web/hooks/',
    'apps/admin/hooks/',
    'apps/expo/hooks/',
)


# Paths a scaffolded project correctly does not have, each with the reason.
# Without this the checker prints the same twenty-three lines every run, and a
# checker nobody reads is a checker that misses the twenty-fourth.
#
# Add an entry only after opening the page and confirming the path is absent
# BY DESIGN. "I could not find it" is the thing this script is for.
EXPECTED_ABSENT = {
    # The reader writes these, in the lesson that names them.
    'apps/api/internal/jobs/report.go': 'written by the reader in the jobs lesson',
    'apps/api/internal/jobs/email.go': 'written by the reader in the jobs lesson',
    'apps/api/internal/database/categories_seeder.go': 'written by the reader in the custom-endpoints tutorial',
    'apps/api/internal/jobs/order_confirmation.go': 'written by the reader in the ecommerce tutorial',
    'apps/api/internal/jobs/task_assigned.go': 'written by the reader in the saas tutorial',
    'apps/admin/components/dashboard/revenue-stats.tsx': 'written by the reader in the ecommerce tutorial',
    'apps/admin/components/dashboard/project-stats.tsx': 'written by the reader in the saas tutorial',
    'apps/web/lib/catalog.ts': 'written by the reader in the custom-endpoints tutorial',
    'apps/web/components/ChatRoom.tsx': 'written by the reader in the realtime-chat course',
    'apps/web/components/RoomMembers.tsx': 'written by the reader in the realtime-chat course',
    'apps/web/components/Typing.tsx': 'written by the reader in the realtime-chat course',

    # Created by something you opt into.
    'apps/api/internal/middleware/tenant.go': 'created by grit plugin add multitenant',
    'apps/admin/components/calendar.tsx': 'created by the example plugin on the authoring page',
    'apps/admin/messages/fr.json': 'created by grit add i18n',

    'apps/api/internal/middleware/rate_limit.go': 'the example middleware the reader writes on the middleware page',
    'apps/admin/lib/uploader.ts': 'written by the reader, over the @repo/upload package',
    'apps/web/lib/products-public.ts': 'written by the reader in the public-api page',
    'apps/web/components/conflict-list.tsx': 'written by the reader in the offline-sync page',

    # The destination of a command on the same page.
    'apps/web/lib/icons.ts': 'the target of a cp in the standalone-usage page',
}


def app_of(path):
    """apps/web/... -> apps/web, so a missing app can be skipped wholesale."""
    bits = path.split('/')
    if bits[0] == 'apps' and len(bits) > 1:
        return 'apps/' + bits[1]
    return bits[0]


def main():
    if not REFERENCE or not os.path.isdir(REFERENCE):
        print('Set GRIT_REFERENCE_PROJECT to a scaffolded project.')
        print('Nothing checked, which is not the same as nothing wrong.')
        return 0

    present = {d for d in ('apps/api', 'apps/web', 'apps/admin', 'apps/expo',
                           'apps/desktop', 'packages/shared', 'cmd', 'e2e')
               if os.path.isdir(os.path.join(REFERENCE, d.replace('/', os.sep)))}
    print('reference: %s' % REFERENCE)
    print('apps present: %s\n' % ', '.join(sorted(present)))

    missing = {}
    excused = set()
    checked = 0

    for dirpath, dirs, names in os.walk(DOCS):
        dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
        for name in names:
            if not name.endswith(('.tsx', '.ts', '.md')):
                continue
            full = os.path.join(dirpath, name)
            rel = os.path.relpath(full, ROOT).replace('\\', '/')
            if any(p in rel for p in SKIP_PARTS):
                continue
            try:
                text = io.open(full, encoding='utf-8').read()
            except (OSError, UnicodeDecodeError):
                continue

            for lineno, line in enumerate(text.splitlines(), 1):
                for m in PATH.finditer(line):
                    claimed = m.group(1)
                    if PLACEHOLDER.search(claimed):
                        continue
                    if app_of(claimed) not in present:
                        continue
                    if claimed.startswith(GENERATED_DIRS):
                        continue
                    checked += 1
                    on_disk = os.path.join(REFERENCE, claimed.replace('/', os.sep))
                    if os.path.exists(on_disk):
                        continue
                    if claimed in EXPECTED_ABSENT:
                        excused.add(claimed)
                        continue
                    missing.setdefault(claimed, []).append((rel, lineno))

    print('%d project paths named in the docs, %d of them not in the '
          'reference project, %d expected absent\n'
          % (checked, len(missing), len(excused)))

    for claimed in sorted(missing, key=lambda c: -len(missing[c])):
        where = missing[claimed]
        print('  %s' % claimed)
        for rel, lineno in where[:4]:
            print('      %s:%d' % (rel, lineno))
        if len(where) > 4:
            print('      ... and %d more' % (len(where) - 4))
        print()

    # An exemption for a path no page names any more is an exemption that
    # will one day excuse a real mistake.
    stale = sorted(set(EXPECTED_ABSENT) - excused)
    if stale:
        print('EXPECTED_ABSENT lists paths no page names any more. Remove '
              'them, or they will excuse a real one later:')
        for claimed in stale:
            print('    %s  (%s)' % (claimed, EXPECTED_ABSENT[claimed]))
        print()

    if not missing and not stale:
        print('every path the docs name is in the reference project, or is '
              'listed with a reason')

    return 1 if (missing or stale) else 0


if __name__ == '__main__':
    sys.exit(main())
