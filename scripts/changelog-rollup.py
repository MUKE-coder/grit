#!/usr/bin/env python3
"""Build the changelog's weekly rollup from the changelog itself.

The changelog is the most useful thing Grit publishes and the hardest to scan:
hundreds of entries, several on the same day, each with a paragraph explaining
what went wrong. Somebody deciding whether to adopt Grit today wants the week,
not the postmortem.

So this reads docs/app/docs/changelog/page.tsx, groups every entry by the week it
shipped in, and writes docs/components/changelog-rollup.tsx, which the changelog
page renders above the full entries. It also gives every entry an id
(id="v3.228.0") so a line in the rollup can link to the entry it summarises.

Run it after adding a changelog entry:

    python scripts/changelog-rollup.py

It rewrites both files in place, prints what it found, and fails if an entry
cannot be read. Failing is the point: the first version of this script required
an <h3> title and silently skipped the 160 older entries that use
<p><strong>...</strong></p> instead, which is exactly the kind of quiet
half-result this repository keeps finding the hard way.
"""
import datetime
import json
import os
import re
import sys

HERE = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CHANGELOG = os.path.join(HERE, 'docs', 'app', 'docs', 'changelog', 'page.tsx')
COMPONENT = os.path.join(HERE, 'docs', 'components', 'changelog-rollup.tsx')

# An entry is a version comment followed by the entry container.
HEADER = re.compile(
    r'\{/\* v(?P<version>[0-9][0-9.]*) \*/\}\s*\n'
    r'[ \t]*<div className="mb-12"(?P<attrs>[^>]*)>')

DATE = re.compile(r'<span className="text-sm text-muted-foreground">([^<]+)</span>')
# Newer entries use an <h3>. Older ones open with a bold paragraph, and one
# (v3.79.0, a batch of admin fixes) leads with a plain paragraph and puts the
# first fix in bold inside a list. Three shapes, tried in that order.
TITLE_H3 = re.compile(r'<h3[^>]*>(.*?)</h3>', re.S)
TITLE_STRONG = re.compile(r'<p>\s*<strong>(.*?)</strong>', re.S)
TITLE_LIST = re.compile(r'<li>\s*<strong>(.*?)</strong>', re.S)


def clean(title: str) -> str:
    """A heading as prose: no tags, no JSX braces, one line."""
    text = re.sub(r"\{'([^']*)'\}", r'\1', title)
    text = re.sub(r'\{&apos;[^}]*\}', "'", text)
    text = re.sub(r'<[^>]+>', '', text)
    for entity, char in (('&apos;', '’'), ('&quot;', '"'), ('&amp;', '&'),
                         ('&lt;', '<'), ('&gt;', '>'), ('&mdash;', ', '),
                         ('&ndash;', '-'), ('&nbsp;', ' ')):
        text = text.replace(entity, char)
    return ' '.join(text.split())


def parse_date(raw: str):
    """Entries carry 'September 12, 2026'; the oldest only 'February 2026'."""
    raw = ' '.join(raw.split())
    for fmt, exact in (('%B %d, %Y', True), ('%d %B %Y', True), ('%B %Y', False)):
        try:
            return datetime.datetime.strptime(raw, fmt).date(), exact
        except ValueError:
            continue
    return None, False


def read_entries(src):
    """Every entry, with the body running to the next entry."""
    matches = list(HEADER.finditer(src))
    entries, broken = [], []
    for i, match in enumerate(matches):
        end = matches[i + 1].start() if i + 1 < len(matches) else len(src)
        body = src[match.end():end]

        title_match = (TITLE_H3.search(body) or TITLE_STRONG.search(body)
                       or TITLE_LIST.search(body))
        date_match = DATE.search(body)
        version = match.group('version')
        if not title_match:
            broken.append(version)
            continue
        day, exact = parse_date(date_match.group(1)) if date_match else (None, False)
        entries.append({'version': version, 'date': day, 'exact': exact,
                        'title': clean(title_match.group(1)), 'dated': date_match is not None})
    return entries, broken, matches


def anchor(src, matches):
    """Give every entry container an id, leaving any it already has."""
    added = 0
    out = []
    last = 0
    for match in matches:
        out.append(src[last:match.start()])
        text = match.group(0)
        if 'id=' not in match.group('attrs'):
            text = text.replace(
                '<div className="mb-12"%s>' % match.group('attrs'),
                '<div className="mb-12"%s id="v%s">' % (match.group('attrs'), match.group('version')),
                1)
            added += 1
        out.append(text)
        last = match.end()
    out.append(src[last:])
    return ''.join(out), added


def bucket(entry):
    if entry['date'] is None:
        return ('unknown', datetime.date.min)
    if entry['exact']:
        monday = entry['date'] - datetime.timedelta(days=entry['date'].weekday())
        return ('week', monday)
    return ('month', entry['date'].replace(day=1))


def label(key):
    kind, value = key
    if kind == 'unknown':
        return 'Undated releases'
    if kind == 'month':
        return value.strftime('%B %Y')
    end = value + datetime.timedelta(days=6)
    if value.month == end.month:
        return '%s %d to %d, %d' % (value.strftime('%B'), value.day, end.day, value.year)
    return '%s %d to %s %d, %d' % (value.strftime('%B'), value.day,
                                   end.strftime('%B'), end.day, end.year)


def main():
    src = open(CHANGELOG, encoding='utf-8').read()
    entries, broken, matches = read_entries(src)

    if not entries:
        sys.exit('no changelog entries found: has the entry markup changed?')
    if broken:
        sys.exit('could not read a title for %d entry/entries: %s\n'
                 'Add an <h3> or a bold first paragraph, or teach this script the new markup.'
                 % (len(broken), ', '.join('v' + v for v in broken[:12])))

    src, added = anchor(src, matches)
    open(CHANGELOG, 'w', encoding='utf-8', newline='').write(src)

    buckets = {}
    for entry in entries:
        buckets.setdefault(bucket(entry), []).append(entry)
    ordered = sorted(buckets.items(), key=lambda kv: kv[0][1], reverse=True)

    lines = [
        '// Generated by scripts/changelog-rollup.py from the entries on the changelog',
        '// page. Do not edit by hand: run the script after adding an entry.',
        '//',
        '// The changelog is worth reading and hard to scan, which is what this is for:',
        '// the week, not the postmortem, with a link into each entry.',
        '',
        'export type RollupEntry = { version: string; title: string }',
        'export type RollupWeek = { label: string; count: number; entries: RollupEntry[] }',
        '',
        'export const changelogRollup: RollupWeek[] = [',
    ]
    for key, group in ordered:
        lines.append('  {')
        lines.append('    label: %s,' % json.dumps(label(key), ensure_ascii=False))
        lines.append('    count: %d,' % len(group))
        lines.append('    entries: [')
        for entry in group:
            lines.append('      { version: %s, title: %s },' % (
                json.dumps(entry['version'], ensure_ascii=False),
                json.dumps(entry['title'], ensure_ascii=False)))
        lines.append('    ],')
        lines.append('  },')
    lines.append(']')
    lines.append('')
    open(COMPONENT, 'w', encoding='utf-8', newline='').write('\n'.join(lines))

    undated = sum(1 for e in entries if not e['dated'])
    print('entries: %d' % len(entries))
    print('anchors added: %d (every entry now has one)' % added)
    print('buckets: %d' % len(ordered))
    for key, group in ordered[:6]:
        print('  %-28s %d' % (label(key), len(group)))
    if undated:
        print('entries with no date: %d' % undated)


if __name__ == '__main__':
    main()
