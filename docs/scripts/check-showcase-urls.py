"""Every showcase entry links somewhere real.

The landing showcase once carried a card whose URL was '#', so the card
went nowhere. A card without a destination is not a showcase entry, it is
a dead link on the homepage.

This checks the two places showcase entries live: the landing showcase on
the homepage and the showcase page. Every `url:` value in either file must
be absolute (http or https). A relative link, a fragment, or a placeholder
fails the check.

Run: python docs/scripts/check-showcase-urls.py
"""
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)

FILES = (
    os.path.join(DOCS, 'components', 'landing', 'showcase.tsx'),
    os.path.join(DOCS, 'app', 'showcase', 'page.tsx'),
)

URL = re.compile(r"""\burl:\s*['"]([^'"]*)['"]""")


def main():
    checked = 0
    problems = []
    for path in FILES:
        rel = os.path.relpath(path, DOCS)
        try:
            lines = open(path, encoding='utf-8').read().splitlines()
        except OSError as exc:
            print(f'{rel}: cannot read ({exc})')
            return 1
        for i, line in enumerate(lines, 1):
            for m in URL.finditer(line):
                checked += 1
                if not m.group(1).startswith(('http://', 'https://')):
                    problems.append((rel, i, m.group(1), line.strip()[:100]))
    if not problems:
        print(f'every showcase entry has an absolute URL ({checked} checked)')
        return 0
    print(f'{len(problems)} showcase URL(s) are not absolute:\n')
    for rel, line, found, text in problems:
        print(f'  {rel}:{line}  "{found}"')
        print(f'    {text}')
        print('    give the entry its real URL, or remove it; a card that')
        print('    links nowhere is worse than no card\n')
    return 1


if __name__ == '__main__':
    sys.exit(main())
