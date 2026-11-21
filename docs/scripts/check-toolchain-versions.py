"""Does the documentation ask for a toolchain older than Grit accepts?

A requirement is a promise: "Go 1.21 or newer" tells a reader that 1.21 will
work. Grit's toolchain floor is 1.26.11 and `grit upgrade` rewrites a go.mod
below it, so that promise has been false for a long time, on a page somebody
reads before they install anything. The same for Next: the scaffold pins
16.3.8 and three pages still say 14 or 15.

Nothing catches this class, because a version in prose is just prose. So the
real pins are read from the Go source that sets them, and every version the
docs state is compared against them.

The hard part is that not every version mention is a requirement. "Generics
arrived in Go 1.18" is a true sentence that will stay true, and a checker
that rewrites it is worse than no checker. Two defences:

  by path     the changelog and the dated learnings posts are history
  by marker   `toolchain-version: historical` on or just above the line

It also checks the constants the guides print from config/site.ts, which is
the same concern from the other end: REQUIRED_GO and PINNED_PNPM are copies of
Go constants, and a copy of a version drifts. Node is not pinned anywhere in
the CLI, so REQUIRED_NODE is compared against Next's own engines field in an
installed copy, and skipped with a note when there is none, rather than
asserting a number this has not read.

Run: python docs/scripts/check-toolchain-versions.py
"""
import glob
import json
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
ROOT = os.path.dirname(DOCS)

MARKER = 'toolchain-version: historical'

# How far above a line the marker may sit.
#
# It was two lines, which made the marker unusable on the pages that needed it
# most. A sentence in JSX prose spans several lines, and a comment inserted
# between two of them removes the space between the words either side of it:
# JSX drops the whitespace next to an expression at a line boundary, so
# "matters when" and "you are stuck" render as "matters whenyou are stuck".
# There was nowhere safe to put it, and eight true sentences about the Go 1.22
# router failed the check with no remedy that did not break the page.
#
# Eight lines reaches the top of an enclosing JSX element or object property,
# which is somewhere a comment is always safe. The cost is that a marker could
# cover a second version mention just below it; that is a smaller problem than
# a marker that cannot be applied at all.
WINDOW = 8

# Paths whose job is to record the past.
SKIP_PARTS = (
    os.path.join('app', 'docs', 'changelog'),
    os.path.join('app', 'docs', 'learnings'),
    os.path.join('app', 'blog'),
)
SKIP_DIRS = {'node_modules', '.next', 'out', '.git'}


def read(path):
    with open(path, encoding='utf-8') as fh:
        return fh.read()


def pinned_go():
    """The floor from internal/scaffold/toolchain_repair.go."""
    src = read(os.path.join(ROOT, 'internal', 'scaffold', 'toolchain_repair.go'))
    m = re.search(r'goToolchainFloor\s*=\s*"(\d+)\.(\d+)', src)
    if not m:
        sys.exit('could not read goToolchainFloor; this checker is now lying')
    return int(m.group(1)), int(m.group(2))


def pinned_next():
    """The major from internal/scaffold/next_version.go."""
    src = read(os.path.join(ROOT, 'internal', 'scaffold', 'next_version.go'))
    m = re.search(r'nextVersion\s*=\s*"(\d+)\.', src)
    if not m:
        sys.exit('could not read nextVersion; this checker is now lying')
    return int(m.group(1))



def pinned_go_patch():
    """The exact go directive, from internal/scaffold/go_toolchain.go."""
    src = read(os.path.join(ROOT, 'internal', 'scaffold', 'go_toolchain.go'))
    m = re.search(r'GoToolchainPin\s*=\s*"([^"]+)"', src)
    if not m:
        sys.exit('could not read GoToolchainPin; this checker is now lying')
    return m.group(1)


def pinned_pnpm():
    """The pnpm every generated project installs with."""
    src = read(os.path.join(ROOT, 'internal', 'scaffold', 'infra_hygiene_repair.go'))
    m = re.search(r'pnpmVersion\s*=\s*"([^"]+)"', src)
    if not m:
        sys.exit('could not read pnpmVersion; this checker is now lying')
    return m.group(1)


def stated_constants():
    """What config/site.ts tells a reader to install."""
    src = read(os.path.join(DOCS, 'config', 'site.ts'))
    out = {}
    for name in ('REQUIRED_GO', 'REQUIRED_NODE', 'PINNED_PNPM'):
        m = re.search(name + r" = '([^']+)'", src)
        if not m:
            sys.exit('config/site.ts no longer exports %s, so the guides are '
                     'printing something this cannot check' % name)
        out[name] = m.group(1)
    return out


def constant_problems():
    """config/site.ts against the Go constants it copies."""
    stated = stated_constants()
    found = []
    for name, actual, source in (
        ('REQUIRED_GO', pinned_go_patch(), 'internal/scaffold.GoToolchainPin'),
        ('PINNED_PNPM', pinned_pnpm(), "internal/scaffold's pnpmVersion"),
    ):
        if stated[name] != actual:
            found.append('config/site.ts %s is %s, %s is %s'
                         % (name, stated[name], source, actual))

    # Next's engines field, from any installed copy. Read, not assumed.
    installed = sorted(glob.glob(os.path.join(
        ROOT, 'examples', '*', 'node_modules', 'next', 'package.json')))
    if installed:
        with open(installed[0], encoding='utf-8') as fh:
            engines = json.load(fh).get('engines', {}).get('node', '')
        floor = re.search(r'(\d+\.\d+)', engines)
        if floor and floor.group(1) != stated['REQUIRED_NODE']:
            found.append("config/site.ts REQUIRED_NODE is %s, but Next's "
                         'engines field says %s'
                         % (stated['REQUIRED_NODE'], engines))
    else:
        print('note: no installed Next under examples/*/node_modules, so '
              'REQUIRED_NODE (%s) was not checked against its engines field'
              % stated['REQUIRED_NODE'])
    return found


GO = re.compile(r'\bGo (\d+)\.(\d+)')
NEXT = re.compile(r'\bNext\.js (\d+)')


def main():
    go_major, go_minor = pinned_go()
    next_major = pinned_next()

    constants = constant_problems()

    problems = []
    for dirpath, dirs, names in os.walk(os.path.join(DOCS, 'app')):
        dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
        for name in names:
            if not name.endswith(('.tsx', '.ts', '.md')):
                continue
            path = os.path.join(dirpath, name)
            rel = os.path.relpath(path, ROOT).replace('\\', '/')
            if any(part.replace('\\', '/') in rel for part in
                   (p.replace('\\', '/') for p in SKIP_PARTS)):
                continue
            try:
                lines = read(path).splitlines()
            except (OSError, UnicodeDecodeError):
                continue
            for i, line in enumerate(lines, 1):
                near = lines[max(0, i - WINDOW - 1):i]
                if MARKER in line or any(MARKER in n for n in near):
                    continue

                for m in GO.finditer(line):
                    major, minor = int(m.group(1)), int(m.group(2))
                    if (major, minor) < (go_major, go_minor):
                        problems.append((
                            rel, i, m.group(0),
                            f'Grit requires Go {go_major}.{go_minor}+',
                            line.strip()[:100]))
                for m in NEXT.finditer(line):
                    if int(m.group(1)) < next_major:
                        problems.append((
                            rel, i, m.group(0),
                            f'the scaffold pins Next.js {next_major}',
                            line.strip()[:100]))

    if constants:
        print('the versions the guides print disagree with the CLI:')
        print('')
        for line in constants:
            print('  ' + line)
        print('')
        print('update docs/config/site.ts, or move the Go constant on '
              'purpose and then update it')
        print('')

    if not problems:
        if not constants:
            print(f'every stated version is current '
                  f'(Go {go_major}.{go_minor}+, Next.js {next_major})')
            return 0
        return 1

    print(f'{len(problems)} stale version(s). Grit needs Go '
          f'{go_major}.{go_minor}+ and pins Next.js {next_major}:\n')
    for rel, line, found, why, text in problems:
        print(f'  {rel}:{line}  "{found}" ({why})')
        print(f'    {text}')
        print(f'    if this is a historical fact, put "{MARKER}" on or above '
              f'the line\n')
    return 1


if __name__ == '__main__':
    sys.exit(main())
