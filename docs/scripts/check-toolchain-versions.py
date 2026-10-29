"""Does the documentation ask for a toolchain older than Grit accepts?

A requirement is a promise: "Go 1.21 or newer" tells a reader that 1.21 will
work. Grit's toolchain floor is 1.26.10 and `grit upgrade` rewrites a go.mod
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

Run: python docs/scripts/check-toolchain-versions.py
"""
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
ROOT = os.path.dirname(DOCS)

MARKER = 'toolchain-version: historical'

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


GO = re.compile(r'\bGo (\d+)\.(\d+)')
NEXT = re.compile(r'\bNext\.js (\d+)')


def main():
    go_major, go_minor = pinned_go()
    next_major = pinned_next()

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
                near = lines[max(0, i - 3):i]
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

    if not problems:
        print(f'every stated version is current '
              f'(Go {go_major}.{go_minor}+, Next.js {next_major})')
        return 0

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
