"""No colour token may be wrapped in hsl().

The site stored colour as HSL triplets once, consumed as hsl(var(--x)). The
redesign moved every token to hex and the implementation guide was explicit
that the two must not be mixed. One component did not get the message.

LaneFlow kept `primary: 'var(--primary)'` inside an hsl() wrapper. The moment
--primary became var(--fg), every hsl(var(--primary) / 0.1) became invalid
CSS, and a browser resolves an invalid fill to black and an invalid stroke to
none. 275 boxes across 121 pages, the ones marked as Grit's own code, rendered
as black holes with no border. Nothing failed. No test noticed. The site built
and deployed for as long as it took somebody to say the diagrams looked wrong.

That is the shape this file exists to catch: a colour that is silently
invalid. Put alpha on a hex token with rgb(from var(--x) r g b / a), which is
what tailwind.config.ts already does.
"""
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)

# hsl() holding either a var() or a template interpolation.
#
# The first version of this looked for a literal `hsl(...var(--`, which NEITHER
# real bug looked like: both built the colour by interpolation,
# `hsl(${TONES[tone]})`, with the var() hiding in a table above. It passed on
# the exact code it existed to catch. Found by putting the bug back and
# watching the check report success.
#
# The literal triplets in globals.css, hsl(240 10% 20%), contain neither and
# stay legal.
BAD = re.compile(r'hsl\([^)]*(?:var\(|\$\{)')

# A line carrying this, or the line above it, is deliberate. Exceptions stay
# visible and next to their reason rather than living in a list in here.
ALLOW = 'colour-tokens: deliberate'

SKIP_DIRS = {'node_modules', '.next', 'out', '.git'}

# The changelog is a historical record and quotes the code it is reporting on,
# including the broken code. Same reason the CLI claim checker skips it: a file
# whose job is to describe a past mistake will contain that mistake.
SKIP_FILES = {os.path.join('app', 'docs', 'changelog', 'page.tsx')}
EXTS = ('.tsx', '.ts', '.css', '.jsx', '.js')


def main():
    problems = []
    for dirpath, dirs, names in os.walk(DOCS):
        dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
        for name in names:
            if not name.endswith(EXTS):
                continue
            path = os.path.join(dirpath, name)
            if any(os.path.relpath(path, DOCS).endswith(skip) for skip in SKIP_FILES):
                continue
            try:
                text = open(path, encoding='utf-8').read()
            except (OSError, UnicodeDecodeError):
                continue
            lines = text.splitlines()
            for i, line in enumerate(lines, 1):
                # A line that only talks about the problem is not the problem.
                stripped = line.strip()
                if stripped.startswith(('//', '*', '/*', '#')):
                    continue
                # The marker may sit a line or two above, because the
                # reason usually needs a sentence and a sentence wraps.
                above = lines[max(0, i - 4):i - 1]
                if ALLOW in line or any(ALLOW in a for a in above):
                    continue
                if BAD.search(line):
                    rel = os.path.relpath(path, os.path.dirname(DOCS))
                    problems.append((rel.replace('\\', '/'), i, stripped[:110]))

    if not problems:
        print('no colour token is wrapped in hsl()')
        return 0

    print(f'{len(problems)} colour token(s) wrapped in hsl(). Every token in '
          f'this site is hex, so each of these resolves to black:\n')
    for rel, line, text in problems:
        print(f'  {rel}:{line}')
        print(f'    {text}')
        print('    use rgb(from var(--token) r g b / <alpha>) instead\n')
    return 1


if __name__ == '__main__':
    sys.exit(main())
