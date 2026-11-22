"""The coverage table and the workflows that are supposed to back it.

docs/config/coverage.ts tells a reader which architectures CI scaffolds and
how deeply. That is the kind of claim that starts true and drifts into
flattery: a workflow is edited, a shape stops being built, and the page goes
on saying it is. Nothing would notice, because the page compiles either way.

So this reads both sides in one pass. For every shape the table says a
workflow scaffolds, that workflow has to contain a `grit new` line producing
that shape; and every shape a workflow does scaffold has to be claimed by some
row, so coverage that quietly appears is recorded rather than lost.

Run: python docs/scripts/check-ci-coverage.py
"""
import io
import os
import re
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
ROOT = os.path.dirname(DOCS)
CONFIG = os.path.join(DOCS, 'config', 'coverage.ts')
WORKFLOWS = os.path.join(ROOT, '.github', 'workflows')

# The architecture flags `grit new` takes. --triple is split by frontend
# because the two resolve different dependency graphs, which is why the
# canary runs both.
ARCH_FLAGS = ['--full', '--triple', '--double', '--single', '--api', '--mobile']

# A workflow that scaffolds nothing is not evidence of anything, so it is not
# read. release.yml and docs.yml build the CLI and the site.
IGNORE = {'release.yml', 'scorecard.yml'}


def shapes_in_workflow(src):
    """The shapes a workflow's `grit new` lines produce."""
    found = set()
    for line in src.splitlines():
        if 'grit new' not in line:
            continue
        # A comment explaining the command is not the command.
        stripped = line.strip()
        if stripped.startswith('#'):
            continue
        arch = next((f for f in ARCH_FLAGS if f in line), None)
        if arch is None:
            # A matrix row supplies the flags: `grit new app ${{ matrix.flags }}`
            continue
        if arch == '--triple':
            if '--vite' in line:
                found.add('--triple --vite')
            elif '--next' in line:
                found.add('--triple --next')
            else:
                found.add('--triple')
        else:
            found.add(arch)
    # Matrix flags, which is how the canary names its two shapes.
    for flags in re.findall(r'flags:\s*"([^"]*)"', src):
        arch = next((f for f in ARCH_FLAGS if f in flags), None)
        if arch == '--triple' and '--vite' in flags:
            found.add('--triple --vite')
        elif arch == '--triple' and '--next' in flags:
            found.add('--triple --next')
        elif arch:
            found.add(arch)
    return found


def claimed():
    """flag -> the workflows coverage.ts says scaffold it."""
    src = io.open(CONFIG, encoding='utf-8').read()
    out = {}
    for block in re.findall(r"\{\s*\n\s*flag: '([^']+)',[\s\S]*?workflows: \[([^\]]*)\]", src):
        flag, files = block
        out[flag] = re.findall(r"'([^']+)'", files)
    return out


def main():
    table = claimed()
    if len(table) < 5:
        print('read only %d rows out of coverage.ts, so this is not reading '
              'the config rather than finding it empty' % len(table))
        return 2

    by_workflow = {}
    for name in sorted(os.listdir(WORKFLOWS)):
        if not name.endswith('.yml') or name in IGNORE:
            continue
        src = io.open(os.path.join(WORKFLOWS, name), encoding='utf-8').read()
        shapes = shapes_in_workflow(src)
        if shapes:
            by_workflow[name] = shapes

    print('%d rows in coverage.ts; %d workflows scaffold a project'
          % (len(table), len(by_workflow)))
    for name, shapes in sorted(by_workflow.items()):
        print('    %-14s %s' % (name, ', '.join(sorted(shapes))))
    print()

    problems = []

    # Every claim is backed.
    for flag, files in sorted(table.items()):
        for name in files:
            shapes = by_workflow.get(name)
            if shapes is None:
                problems.append(
                    '%s claims %s scaffolds it, and that workflow scaffolds '
                    'nothing (or does not exist)' % (flag, name))
                continue
            # --triple --next is covered by --full, which is a triple project
            # with more in it. The table says so in prose; accept it here.
            if flag in shapes:
                continue
            if flag == '--triple --next' and '--full' in shapes:
                continue
            problems.append(
                '%s claims %s scaffolds it, and no grit new line there does. '
                'That workflow scaffolds: %s'
                % (flag, name, ', '.join(sorted(shapes)) or 'nothing'))

    # And every shape a workflow scaffolds is claimed by some row, so new
    # coverage is recorded rather than going unmentioned.
    all_claimed = set(table)
    for name, shapes in by_workflow.items():
        for shape in shapes:
            if shape in all_claimed:
                continue
            if shape == '--triple' and ('--triple --next' in all_claimed or
                                        '--triple --vite' in all_claimed):
                continue
            problems.append(
                '%s scaffolds %s and no row in coverage.ts mentions that '
                'shape, so the page understates what is covered' % (name, shape))

    if problems:
        for problem in problems:
            print('  %s' % problem)
        return 1

    print('every coverage claim is backed by a grit new line, and every '
          'shape a workflow scaffolds is on the page')
    return 0


if __name__ == '__main__':
    sys.exit(main())
