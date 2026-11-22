"""Advisories against the packages Grit pins exactly, at any severity.

The canary already runs `pnpm audit --audit-level high` over a scaffolded
project, advisory-only, and that is the right setting for the graph: a
generated app resolves a few thousand packages, most of them build-time
tooling nobody ships, and failing the drift signal on a transitive HIGH would
make the canary useless within a week.

It is the wrong setting for the handful of versions Grit chose. Those are the
project's own decision, they reach every app anybody generates, and moving one
is a one-line change here. On 2026-11-22 the pinned Next carried seven
MODERATE advisories, every one patched three weeks earlier, and nothing said
so: `--audit-level high` filtered all seven out, and the pin is only ever read
by a human who goes looking.

The rule, which needs no list to maintain: a dependency written as an exact
version in a template's package.json is a version Grit picked, so an advisory
against it is Grit's to fix. A dependency written as a range resolved to
whatever npm had this morning, and that is the graph.

Run:

    python scripts/audit-pinned.py /path/to/scaffolded/project

Exit 1 when a pinned package has an advisory. Prints the rest for the log.
"""
import io
import json
import os
import subprocess
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

# A version is a pin when it is just a version: no ^, ~, >, <, ||, x or *.
# "16.4.6" is a decision; "^16.4.6" is a floor.
RANGE_CHARS = set('^~><|*x ')

# Workspace protocol, file paths and git URLs are not npm versions.
NOT_A_VERSION = ('workspace:', 'file:', 'link:', 'git+', 'github:', 'npm:', 'catalog:')


def is_pin(spec):
    if not spec or spec.startswith(NOT_A_VERSION):
        return False
    if any(c in RANGE_CHARS for c in spec):
        return False
    # A bare "latest" or "next" is not a pin either.
    return spec[0].isdigit()


def pinned_packages(root):
    """Every dependency written as an exact version, and where it was found."""
    found = {}
    for dirpath, dirs, names in os.walk(root):
        dirs[:] = [d for d in dirs if d not in ('node_modules', '.next', '.git', 'dist', 'out')]
        if 'package.json' not in names:
            continue
        path = os.path.join(dirpath, 'package.json')
        try:
            data = json.load(io.open(path, encoding='utf-8'))
        except (OSError, ValueError):
            continue
        where = os.path.relpath(path, root).replace('\\', '/')
        for section in ('dependencies', 'devDependencies'):
            for name, spec in (data.get(section) or {}).items():
                if is_pin(spec):
                    found.setdefault(name, []).append('%s (%s)' % (where, spec))
    return found


def audit(root):
    """pnpm audit as JSON. An audit that cannot run is a failure, not a pass."""
    proc = subprocess.run(
        ['pnpm', 'audit', '--json'],
        cwd=root, capture_output=True, text=True, shell=(os.name == 'nt'))
    out = proc.stdout.strip()
    if not out:
        print('pnpm audit produced no output:\n%s' % (proc.stderr or '').strip()[:2000])
        sys.exit(2)
    try:
        return json.loads(out)
    except ValueError:
        # pnpm prints a human summary when the registry refuses the request.
        print('pnpm audit did not return JSON:\n%s' % out[:2000])
        sys.exit(2)


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        return 2
    root = os.path.abspath(sys.argv[1])
    if not os.path.isdir(root):
        print('not a directory: %s' % root)
        return 2

    pins = pinned_packages(root)
    if not pins:
        print('no exactly-pinned dependencies found under %s, which is not a '
              'result this check can report on: it means the package.json '
              'files were not read.' % root)
        return 2
    print('%d exactly-pinned dependencies: %s\n'
          % (len(pins), ', '.join(sorted(pins))))

    report = audit(root)
    advisories = report.get('advisories') or {}

    ours, graph = [], []
    for advisory in advisories.values():
        module = advisory.get('module_name')
        (ours if module in pins else graph).append(advisory)

    if graph:
        print('%d advisories against the dependency graph. Advisory only: a '
              'generated app resolves thousands of packages and most of them '
              'are build-time tooling.' % len(graph))
        for a in sorted(graph, key=lambda a: a.get('module_name') or ''):
            print('    %-10s %-24s %s' % (a.get('severity'), a.get('module_name'),
                                          (a.get('title') or '')[:90]))
        print()

    if not ours:
        print('no advisories against a pinned package. The versions Grit '
              'chose are clean.')
        return 0

    print('%d advisories against packages Grit pins. These are the project\'s '
          'own to move:\n' % len(ours))
    for a in sorted(ours, key=lambda a: a.get('module_name') or ''):
        module = a.get('module_name')
        print('  %s %s' % (a.get('severity', '?').upper(), module))
        print('    %s' % (a.get('title') or '').strip())
        print('    affected %s, patched in %s'
              % (a.get('vulnerable_versions'), a.get('patched_versions')))
        for where in pins.get(module, []):
            print('    pinned in %s' % where)
        if a.get('url'):
            print('    %s' % a['url'])
        print()
    print('Move the pin in the template, not in the generated project: the '
          'project is output.')
    return 1


if __name__ == '__main__':
    sys.exit(main())
