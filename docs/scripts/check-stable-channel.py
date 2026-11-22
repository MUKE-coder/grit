"""The stable channel, checked against the things that have to agree with it.

Four files are involved in one fact: the config the promoter writes, the route
that serves it, the Go constant the CLI fetches, and the version the site says
is current. Each of them compiles and builds on its own while the set of them
is wrong, which is the shape of bug this project keeps finding.

What this checks:

  - the config parses and names a version
  - stable is not ahead of edge, which would mean the promoter chose an
    unreleased version or the site is behind
  - the promotion date is not before the release date
  - the route that serves /stable.json exists
  - the URL compiled into the CLI is that route

Run: python docs/scripts/check-stable-channel.py
"""
import io
import os
import re
import sys

sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
ROOT = os.path.dirname(DOCS)

CONFIG = os.path.join(DOCS, 'config', 'stable.ts')
SITE = os.path.join(DOCS, 'config', 'site.ts')
ROUTE = os.path.join(DOCS, 'app', 'stable.json', 'route.ts')
CHANNEL_GO = os.path.join(ROOT, 'internal', 'selfupdate', 'channel.go')


def read(path):
    return io.open(path, encoding='utf-8').read()


def version_tuple(v):
    return tuple(int(p) for p in re.findall(r'\d+', v)[:3])


def main():
    problems = []

    src = read(CONFIG)
    fields = {}
    for key in ('version', 'promoted', 'released'):
        m = re.search(r"%s: '([^']+)'," % key, src)
        if not m:
            problems.append('docs/config/stable.ts has no %s' % key)
        else:
            fields[key] = m.group(1)

    edge = re.search(r"GRIT_VERSION = '([^']+)'", read(SITE))
    if not edge:
        problems.append('could not read GRIT_VERSION out of docs/config/site.ts')

    if 'version' in fields and edge:
        stable_v = version_tuple(fields['version'])
        edge_v = version_tuple(edge.group(1))
        if stable_v > edge_v:
            problems.append(
                'stable is v%s and edge is v%s: stable is ahead, which means '
                'either the promoter chose a version that was never released '
                'or the site is behind the channel'
                % (fields['version'], edge.group(1)))
        elif stable_v == edge_v:
            # Possible and fine on a quiet week, but worth printing: it means
            # nothing has shipped in the ageing window.
            print('note: stable and edge are the same release (v%s)' % fields['version'])

    if 'promoted' in fields and 'released' in fields:
        if fields['promoted'] < fields['released']:
            problems.append('promoted %s is before released %s'
                            % (fields['promoted'], fields['released']))

    if not os.path.isfile(ROUTE):
        problems.append('app/stable.json/route.ts is missing, so the channel '
                        'the CLI reads is not served')

    # The URL the CLI fetches has to be the route the site serves. These are
    # in different languages in different directories and nothing else pairs
    # them.
    go = read(CHANNEL_GO)
    m = re.search(r'ChannelURL = "([^"]+)"', go)
    if not m:
        problems.append('internal/selfupdate/channel.go has no ChannelURL')
    else:
        url = m.group(1)
        if not url.endswith('/stable.json'):
            problems.append(
                'the CLI fetches %s, and the site serves /stable.json. '
                'A channel nobody can read is a flag that always fails.' % url)
        else:
            print('the CLI reads %s' % url)

    if problems:
        for problem in problems:
            print('  %s' % problem)
        return 1

    print('stable v%s (released %s, promoted %s), edge v%s: consistent'
          % (fields.get('version'), fields.get('released'),
             fields.get('promoted'), edge.group(1) if edge else '?'))
    return 0


if __name__ == '__main__':
    sys.exit(main())
