#!/usr/bin/env python3
"""Write CHANGELOG.md from the changelog page, so the history is greppable.

The canonical changelog is docs/app/docs/changelog/page.tsx, about 37,000 lines
of TSX, and CHANGELOG.md was a pointer to the website. So the history of this
project could not be grepped from a checkout, diffed between two tags, read on a
plane, or indexed by `grit docs` (grit#167).

The direction matters and is the whole design: this is GENERATED from the page,
never maintained beside it. A hand-written mirror is two changelogs within a
week, and the second one is always the stale one. The page keeps its code blocks
and its styling and stays canonical; this is the text of the same entries.

    python scripts/changelog-md.py            # write it
    python scripts/changelog-md.py --check     # fail if it is out of date (CI)

--check is what stops the mirror drifting: it regenerates in memory and compares,
so a release that forgets to run this, or an edit made to CHANGELOG.md by hand,
fails rather than quietly diverging.
"""
import io
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
sys.path.insert(0, HERE)  # -I does not add the script's own directory

# The page parser already exists, and reusing it is the point: one reader means
# the Markdown cannot find a different set of entries from the rollup.
import importlib.util

spec = importlib.util.spec_from_file_location(
    'changelog_rollup', os.path.join(HERE, 'changelog-rollup.py'))
rollup = importlib.util.module_from_spec(spec)
spec.loader.exec_module(rollup)

CHANGELOG_PAGE = os.path.join(ROOT, 'docs', 'app', 'docs', 'changelog', 'page.tsx')
OUT = os.path.join(ROOT, 'CHANGELOG.md')

SITE = 'https://gritframework.com'

HEADER = """# Changelog

Every release, newest first: what was wrong, why it was wrong, and what changed.
Not a list of commit subjects.

**This file is generated** from the changelog page by
`scripts/changelog-md.py`, which runs as part of the release. Edit the page at
`docs/app/docs/changelog/page.tsx`, not this file: an edit here is overwritten
by the next release, and CI fails on one that has not been.

The rendered version, with syntax highlighting and the weekly rollup, is at
**[gritframework.com/docs/changelog]({site}/docs/changelog)**.

Also worth knowing:

- **[Releases](https://github.com/MUKE-coder/grit/releases)** on GitHub, with the
  binaries, their checksums and the signing material described in
  [SECURITY.md](SECURITY.md).
- **[What is stable and what is not]({site}/docs/stability)**, per subsystem.
- `grit update` brings the CLI to the latest release; `grit upgrade` brings an
  existing project's framework files up to the CLI it is run with.

---
""".replace('{site}', SITE)


# ── JSX to Markdown, for the subset the changelog actually uses ─────────────

def code_blocks(body):
    """Pull every <CodeBlock .../> out, leaving a placeholder.

    Done first and separately because the code is the one part that must not be
    touched by any of the tag rewriting below: a Go snippet containing `<-` or
    a TSX one containing `</div>` would be eaten by a tag regex.
    """
    blocks = []

    def take(match):
        attrs, code = match.group('attrs'), match.group('code')
        lang = ''
        lang_match = re.search(r'language="([^"]+)"', attrs)
        if lang_match:
            lang = lang_match.group(1)
        elif 'terminal' in attrs:
            lang = 'bash'
        filename = ''
        name_match = re.search(r'filename="([^"]+)"', attrs)
        if name_match:
            filename = name_match.group(1)
        blocks.append((lang, filename, code))
        return '\x00CODE%d\x00' % (len(blocks) - 1)

    # code={`...`} with the backticks, which is how every block on the page is
    # written.
    #
    # The body is a JS template literal, so it can contain an ESCAPED backtick,
    # and the changelog's blocks routinely do: a TSX snippet showing
    # aria-label={`Remove ${file.name}`} carries two of them. A non-greedy
    # `.*?` stops at the first \`} and silently truncates the snippet mid-line,
    # which is exactly what the first version of this script did. So the body
    # is "any escaped pair, or any character that is neither a backtick nor a
    # backslash", which cannot end early.
    pattern = re.compile(
        r'<CodeBlock(?P<attrs>[^>]*?)code=\{`(?P<code>(?:\\.|[^`\\])*)`\}(?:[^>]*?)/?>',
        re.S)
    body = pattern.sub(take, body)

    # One block on the page is a single-quoted string rather than a template
    # literal, which the pattern above does not match. Finding that out was
    # luck: the block was dropped and nothing said so, which is the quiet
    # half-result this repository keeps paying for. Hence the second form, and
    # the check in render() that fails on a third.
    quoted = re.compile(
        r"<CodeBlock(?P<attrs>[^>]*?)code=\{'(?P<code>(?:\\.|[^'\\])*)'\}(?:[^>]*?)/?>",
        re.S)
    body = quoted.sub(take, body)

    return body, blocks


def restore_code(text, blocks):
    for i, (lang, filename, code) in enumerate(blocks):
        code = code.replace('\\`', '`').replace('\\$', '$').replace('\\\\', '\\')
        fence = '```' + lang
        parts = []
        if filename:
            parts.append('`%s`\n' % filename)
        parts.append('%s\n%s\n```' % (fence, code.strip('\n')))
        text = text.replace('\x00CODE%d\x00' % i, '\n'.join(parts))
    return text


ENTITIES = [
    ('&apos;', "'"), ('&quot;', '"'), ('&mdash;', ': '), ('&ndash;', '-'),
    ('&nbsp;', ' '), ('&lt;', '<'), ('&gt;', '>'), ('&amp;', '&'),
    ('&hellip;', '...'), ('&rarr;', '->'), ('&times;', 'x'),
]


def inline(text):
    """Inline JSX to Markdown."""
    # JSX string escapes: {' '} is a space, {'...'} is its contents.
    text = re.sub(r"\{'\s*'\}", ' ', text)
    text = re.sub(r"\{'([^']*)'\}", r'\1', text)
    text = re.sub(r'\{&apos;([^}]*)\}', r"'\1", text)

    # Links. An internal href becomes absolute, so the file is useful offline
    # and on GitHub, where a bare /docs/... path resolves to the wrong place.
    def link(match):
        href, label = match.group('href'), match.group('label')
        if href.startswith('/'):
            href = SITE + href
        return '[%s](%s)' % (label.strip(), href)

    text = re.sub(
        r'<(?:Link|a)[^>]*href="(?P<href>[^"]+)"[^>]*>(?P<label>.*?)</(?:Link|a)>',
        link, text, flags=re.S)

    text = re.sub(r'<code[^>]*>(.*?)</code>', lambda m: '`%s`' % m.group(1), text, flags=re.S)
    text = re.sub(r'<strong[^>]*>(.*?)</strong>', r'**\1**', text, flags=re.S)
    text = re.sub(r'<em[^>]*>(.*?)</em>', r'*\1*', text, flags=re.S)

    text = re.sub(r'<[^>]+>', '', text)
    for entity, char in ENTITIES:
        text = text.replace(entity, char)
    # Numeric entities. The changelog uses &#123; and &#125; for the braces in
    # a JSX snippet written as prose, and leaving them raw puts "&#123;" in the
    # Markdown, which is worse than the brace it was escaping.
    text = re.sub(r'&#(\d+);', lambda m: chr(int(m.group(1))), text)
    text = re.sub(r'&#x([0-9a-fA-F]+);', lambda m: chr(int(m.group(1), 16)), text)
    return ' '.join(text.split())


def body_markdown(body):
    """One entry's body as Markdown."""
    body, blocks = code_blocks(body)

    out = []
    # Walk the block-level elements in order. Anything not matched here is not
    # prose: a <div> wrapper, a <DocsShot>, the version badge.
    pattern = re.compile(
        r'<h3[^>]*>(?P<h3>.*?)</h3>'
        r'|<h4[^>]*>(?P<h4>.*?)</h4>'
        r'|<p[^>]*>(?P<p>.*?)</p>'
        r'|<li[^>]*>(?P<li>.*?)</li>'
        r'|(?P<code>\x00CODE\d+\x00)',
        re.S)

    for match in pattern.finditer(body):
        if match.group('h3') is not None:
            text = inline(match.group('h3'))
            if text:
                out.append('### ' + text)
        elif match.group('h4') is not None:
            text = inline(match.group('h4'))
            if text:
                out.append('#### ' + text)
        elif match.group('p') is not None:
            text = inline(match.group('p'))
            if text:
                out.append(text)
        elif match.group('li') is not None:
            text = inline(match.group('li'))
            if text:
                out.append('- ' + text)
        else:
            out.append(match.group('code'))

    # Collapse runs of list items into one block so Markdown renders a list
    # rather than separated paragraphs.
    merged = []
    for chunk in out:
        if chunk.startswith('- ') and merged and merged[-1].startswith('- '):
            merged[-1] = merged[-1] + '\n' + chunk
        else:
            merged.append(chunk)

    return restore_code('\n\n'.join(merged), blocks)


def render(src):
    entries, broken, matches = rollup.read_entries(src)
    if broken:
        print('FAIL: %d entries have no readable title: %s'
              % (len(broken), ', '.join(broken[:10])), file=sys.stderr)
        return None, entries

    # read_entries gives the title and the date; the body has to be re-cut,
    # because the Markdown needs the prose and the rollup only needed a heading.
    bodies = {}
    for i, match in enumerate(matches):
        end = matches[i + 1].start() if i + 1 < len(matches) else len(src)
        bodies[match.group('version')] = src[match.end():end]

    # A CodeBlock form code_blocks does not recognise is dropped, not reported.
    # That happened once already, with a single-quoted code= instead of a
    # template literal, and the only reason anyone noticed was counting the
    # fences afterwards. Count them here instead, every run.
    wanted = len(re.findall(r'<CodeBlock', src))
    stripped, found = code_blocks(src)
    if len(found) != wanted or '<CodeBlock' in stripped:
        print('FAIL: %d <CodeBlock> on the page, %d understood.\n'
              % (wanted, len(found)), file=sys.stderr)
        for match in re.finditer(r'<CodeBlock', stripped):
            print('  not understood, at offset %d:' % match.start(), file=sys.stderr)
            print('    %s' % stripped[match.start():match.start() + 160].replace('\n', ' '),
                  file=sys.stderr)
        print('\n  Teach code_blocks() the form, rather than letting the block '
              'vanish from\n  CHANGELOG.md with nothing said.', file=sys.stderr)
        return None, entries

    parts = [HEADER]
    for entry in entries:
        version = entry['version']
        date = ''
        date_match = rollup.DATE.search(bodies[version])
        if date_match:
            date = ' - ' + ' '.join(date_match.group(1).split())
        parts.append('\n## v%s%s\n' % (version, date))
        text = body_markdown(bodies[version])
        parts.append(text if text else '_No prose on this entry._')
        parts.append('')

    return '\n'.join(parts).rstrip() + '\n', entries


def main():
    check = '--check' in sys.argv

    src = io.open(CHANGELOG_PAGE, encoding='utf-8').read()
    text, entries = render(src)
    if text is None:
        return 1

    if check:
        if not os.path.exists(OUT):
            print('FAIL: CHANGELOG.md is missing. Run python scripts/changelog-md.py')
            return 1
        have = io.open(OUT, encoding='utf-8').read()
        if have != text:
            # Say which versions differ rather than "they differ", because the
            # useful answer is almost always "you added an entry and did not
            # regenerate".
            on_disk = set(re.findall(r'^## v([0-9][0-9.]*)', have, re.M))
            in_page = {e['version'] for e in entries}
            missing = sorted(in_page - on_disk)
            extra = sorted(on_disk - in_page)
            print('FAIL: CHANGELOG.md is out of date.')
            if missing:
                print('  on the page and not in the file: %s' % ', '.join(missing))
            if extra:
                print('  in the file and not on the page: %s' % ', '.join(extra))
            if not missing and not extra:
                print('  the same %d versions, different text: somebody edited '
                      'CHANGELOG.md by hand, or the page prose changed.' % len(in_page))
            print('\n  Run: python scripts/changelog-md.py')
            return 1
        print('OK: CHANGELOG.md matches the page (%d releases)' % len(entries))
        return 0

    io.open(OUT, 'w', encoding='utf-8', newline='').write(text)
    print('wrote CHANGELOG.md: %d releases, %d KB'
          % (len(entries), len(text.encode('utf-8')) // 1024))
    return 0


if __name__ == '__main__':
    sys.exit(main())
