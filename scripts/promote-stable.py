"""Move the stable channel forward, when the rule says it may.

Grit ships every merged change, several times a day. That is right for the
project and wrong for a team with a change-control process, and it has a real
cost: a handful of releases have fixed a regression from the days before, so
the newest release is occasionally the worst one to be on (grit#245).

The rule, and nothing else decides this:

    stable is the newest release that
      1. is at least MIN_AGE_DAYS old,
      2. is not a prerelease, and
      3. had at least MIN_RUNS scheduled canary runs in the MIN_AGE_DAYS after
         it published, every one of them green.

Rule 3 is bounded to the release's own ageing window on purpose. The first
version of it said "every scheduled canary since it published", and running
that against real data showed the flaw immediately: the canary failed on two
days in November, so every release older than a week had a failure somewhere
after it and nothing at all was promotable, for another week. A failure three
weeks later says nothing about a release that already aged out clean.

What that proves, honestly: the canary scaffolds from main rather than from
the candidate tag, so a green week means the scaffolder has been healthy
through that week, not that this exact release was canaried. What it buys is
what the report asked for: the last week of releases is excluded, and that is
where the regressions it named lived.

No backports. A fix lands on edge and reaches stable when it has aged.

Run:
    python scripts/promote-stable.py             # report, write nothing
    python scripts/promote-stable.py --write     # update docs/config/stable.ts
    python scripts/promote-stable.py --self-test # check the rule, no network

The first two need the gh CLI, authenticated.
"""
import argparse
import datetime as dt
import io
import json
import os
import re
import subprocess
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
CONFIG = os.path.join(ROOT, 'docs', 'config', 'stable.ts')

MIN_AGE_DAYS = 7

# How many scheduled runs have to land inside the window. Seven days should
# give seven, and a run can be missed: GitHub drops scheduled runs under load,
# and the repository can be quiet. Five is enough evidence and does not make
# the channel hostage to the scheduler.
MIN_RUNS = 5

REPO = 'MUKE-coder/grit'

# How far back to look. A hundred releases is about a fortnight here, which is
# twice the window the rule can ever choose from.
RELEASE_LIMIT = 100
CANARY_LIMIT = 60


def gh(args):
    """gh, returning parsed JSON. A failure here is fatal: a promoter that
    guesses when it cannot see the evidence is worse than one that stops."""
    proc = subprocess.run(['gh'] + args, capture_output=True, text=True,
                          shell=(os.name == 'nt'))
    if proc.returncode != 0:
        sys.exit('gh %s failed:\n%s' % (' '.join(args), (proc.stderr or '').strip()[:2000]))
    try:
        return json.loads(proc.stdout or 'null')
    except ValueError:
        sys.exit('gh %s did not return JSON:\n%s' % (' '.join(args), proc.stdout[:2000]))


def parse_time(s):
    """GitHub's ISO-8601, as an aware UTC datetime."""
    return dt.datetime.strptime(s.replace('Z', '+0000'), '%Y-%m-%dT%H:%M:%S%z')


def releases():
    """Published, non-prerelease releases, newest first."""
    rows = gh(['release', 'list', '--repo', REPO, '--limit', str(RELEASE_LIMIT),
               '--json', 'tagName,publishedAt,isPrerelease,isDraft'])
    out = []
    for r in rows or []:
        if r.get('isPrerelease') or r.get('isDraft'):
            continue
        if not r.get('publishedAt'):
            continue
        out.append((r['tagName'].lstrip('v'), parse_time(r['publishedAt'])))
    out.sort(key=lambda p: p[1], reverse=True)
    return out


def canary_runs():
    """Scheduled canary runs, newest first, as (started, conclusion).

    Scheduled only. A workflow_dispatch run is somebody reproducing a failure
    by hand, and counting a green re-run as evidence would let anybody promote
    a release by pressing a button.
    """
    rows = gh(['run', 'list', '--repo', REPO, '--workflow', 'canary.yml',
               '--limit', str(CANARY_LIMIT),
               '--json', 'conclusion,createdAt,event,status'])
    out = []
    for r in rows or []:
        if r.get('event') != 'schedule' or r.get('status') != 'completed':
            continue
        out.append((parse_time(r['createdAt']), r.get('conclusion')))
    out.sort(key=lambda p: p[0], reverse=True)
    return out


def current():
    """What docs/config/stable.ts says today."""
    src = io.open(CONFIG, encoding='utf-8').read()
    m = re.search(r"export const STABLE: StableChannel = \{([^}]*)\}", src)
    if not m:
        sys.exit('could not read STABLE out of docs/config/stable.ts')
    version = re.search(r"version: '([^']+)'", m.group(1))
    return version.group(1) if version else None


def choose(rels, runs, now):
    """The newest release the rule allows, and why each newer one lost.

    The window is the MIN_AGE_DAYS after the release, and only runs inside it
    count. See the note at the top for why it is bounded rather than open.
    """
    cutoff = now - dt.timedelta(days=MIN_AGE_DAYS)
    notes = []
    for version, published in rels:
        if published > cutoff:
            notes.append('%-12s too new (published %s)' % (version, published.date()))
            continue
        window_end = published + dt.timedelta(days=MIN_AGE_DAYS)
        window = [(at, conclusion) for at, conclusion in runs
                  if published < at <= window_end]
        failed = [at for at, conclusion in window if conclusion != 'success']
        if failed:
            notes.append('%-12s %d canary failure(s) in its first %d days, '
                         'the latest %s'
                         % (version, len(failed), MIN_AGE_DAYS, failed[0].date()))
            continue
        if len(window) < MIN_RUNS:
            notes.append('%-12s only %d scheduled canary run(s) in its first '
                         '%d days, want %d'
                         % (version, len(window), MIN_AGE_DAYS, MIN_RUNS))
            continue
        return version, published, len(window), notes
    return None, None, 0, notes


def write(version, published, runs, now):
    src = io.open(CONFIG, 'rb').read().decode('utf-8')
    nl = '\r\n' if '\r\n' in src else '\n'
    block = (
        "export const STABLE: StableChannel = {" + nl +
        "  version: '%s'," % version + nl +
        "  promoted: '%s'," % now.date().isoformat() + nl +
        "  released: '%s'," % published.date().isoformat() + nl +
        "  canaryRuns: %d," % runs + nl +
        "}"
    )
    out = re.sub(r"export const STABLE: StableChannel = \{[^}]*\}", lambda _: block, src, count=1)
    if out == src:
        sys.exit('nothing replaced in docs/config/stable.ts')
    io.open(CONFIG, 'wb').write(out.encode('utf-8'))


def self_test():
    """The rule, against cases chosen because each one went wrong once.

    Run in CI before the real promotion. The rule decides which release every
    cautious team installs, and it reads GitHub, so the only way to exercise
    its edges is to hand it the answers.
    """
    day = dt.datetime(2026, 11, 22, 12, 0, tzinfo=dt.timezone.utc)

    def at(d, h=3):
        return dt.datetime(2026, 11, d, h, 0, tzinfo=dt.timezone.utc)

    green = lambda d: (at(d), 'success')
    red = lambda d: (at(d), 'failure')

    cases = [
        (
            'the newest release old enough with a clean week wins',
            [('3.9.0', at(21)), ('3.8.0', at(14)), ('3.7.0', at(10))],
            [green(d) for d in range(11, 22)],
            '3.8.0',
        ),
        (
            'a release younger than the window is never chosen',
            [('3.9.0', at(21))],
            [green(d) for d in range(11, 22)],
            None,
        ),
        (
            'one failure inside the window disqualifies it',
            [('3.8.0', at(14))],
            [green(14), red(16)] + [green(d) for d in range(17, 22)],
            None,
        ),
        (
            'a failure after the window does not, which is the whole reason '
            'the window is bounded',
            [('3.7.0', at(10))],
            [green(d) for d in range(11, 18)] + [red(20), red(21)],
            '3.7.0',
        ),
        (
            'too few runs in the window is not evidence',
            [('3.8.0', at(14))],
            [green(15), green(16)],
            None,
        ),
        (
            'a quiet week falls through to an older release',
            [('3.9.0', at(14)), ('3.8.0', at(2))],
            [green(3), green(4), green(5), green(6), green(7), green(8),
             green(9), green(15), green(16)],
            '3.8.0',
        ),
    ]

    failures = 0
    for name, rels, runs, want in cases:
        rels = sorted(rels, key=lambda p: p[1], reverse=True)
        runs = sorted(runs, key=lambda p: p[0], reverse=True)
        got, _, _, _ = choose(rels, runs, day)
        ok = got == want
        print('  %s %s' % ('PASS' if ok else 'FAIL', name))
        if not ok:
            print('        got %r, want %r' % (got, want))
            failures += 1
    if failures:
        print('\n%d case(s) failed' % failures)
        return 1
    print('\nthe rule behaves as described')
    return 0


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--write', action='store_true',
                    help='update docs/config/stable.ts when the choice changed')
    ap.add_argument('--self-test', action='store_true',
                    help='check the rule against known cases and exit')
    args = ap.parse_args()

    if args.self_test:
        return self_test()

    now = dt.datetime.now(dt.timezone.utc)
    rels = releases()
    runs = canary_runs()
    if not rels:
        sys.exit('no releases found, which is not a result: stopping rather '
                 'than promoting nothing')
    if not runs:
        sys.exit('no scheduled canary runs found, so rule 3 cannot be '
                 'evaluated. Stopping rather than promoting on two rules.')

    print('%d releases, %d scheduled canary runs, newest %s'
          % (len(rels), len(runs), runs[0][0].date()))

    version, published, count, notes = choose(rels, runs, now)
    # The near misses are the informative part: a wall of "too new" says
    # only that releases are frequent, which is known.
    too_new = [n for n in notes if 'too new' in n]
    other = [n for n in notes if 'too new' not in n]
    if too_new:
        print('    %d release(s) are younger than %d days'
              % (len(too_new), MIN_AGE_DAYS))
    for note in other[-12:]:
        print('    skipped %s' % note)

    if version is None:
        print('\nno release satisfies the rule today. stable stays where it is.')
        return 0

    was = current()
    print('\nstable: %s (published %s, %d green canary runs in its first %d days)'
          % (version, published.date(), count, MIN_AGE_DAYS))
    if was == version:
        print('already what docs/config/stable.ts says. Nothing to do.')
        return 0

    print('moves from %s to %s' % (was, version))
    if not args.write:
        print('(report only; pass --write to update the file)')
        return 0
    write(version, published, count, now)
    print('docs/config/stable.ts updated')
    return 0


if __name__ == '__main__':
    sys.exit(main())
