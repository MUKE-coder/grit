"""Do the book's Go snippets parse?

Not compile: most are fragments, and a fragment has no imports to resolve.
Parsing is the check that matters for "was this rushed" anyway, because it
catches what a reader hits when they type the snippet in: an unclosed brace, a
stray comma, a string literal broken across two lines.

The hard part is not parsing, it is not crying wolf. A teaching snippet is
rarely a whole file: it mixes a declaration with the lines that use it, it
elides the boring middle, and some of it is deliberately illegal because the
chapter is explaining what does not compile. A checker that calls those errors
gets switched off within a week.

So: each snippet is tried in several shapes and passes if ANY parses; elided
snippets are skipped; and a snippet the author marked as intentionally broken
is skipped and SAID so, which is how the exception stays visible.
"""
import os
import re
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
BOOK = os.path.join(os.path.dirname(HERE), 'chapters')
BLOCK = re.compile(r'```go\n(.*?)```', re.S)

# An elision: the author cut code out on purpose, so what is left cannot be
# expected to parse. Both the line form and the inline forms.
ELIDED = re.compile(
    r'^\s*(\.\.\.|//\s*\.\.\.)\s*$'     # a line that is just ...
    r'|\{\s*\.\.\.\s*\}'                 # func(...) { ... }
    r'|\(\s*\.\.\.\s*\)'                 # Create(...) error
    r'|,\s*\.\.\.\s*\)',                # Open(x, ...)
    re.M)

# Code that is wrong on purpose, because the prose is explaining why. The
# marker is a comment on any line of the block.
DELIBERATE = re.compile(
    r'//.*\b(illegal|does not compile|will not compile|wrong|broken|bad|'
    r'do not do this|never do this|ignored|fails to compile)\b', re.I)

TOPLEVEL_START = re.compile(r'^(func|type|var|const|import|package)\b')


def strip_lead_comments(code):
    lines = code.splitlines()
    i = 0
    while i < len(lines) and (not lines[i].strip() or lines[i].lstrip().startswith('//')):
        i += 1
    return '\n'.join(lines[i:])


def split_decls(code):
    decls, rest = [], []
    depth = 0
    current = None
    for line in code.splitlines():
        if depth == 0 and TOPLEVEL_START.match(line.strip()):
            current = decls
        if current is None:
            rest.append(line)
        else:
            decls.append(line)
        depth += line.count('{') - line.count('}')
        depth += line.count('(') - line.count(')')
        if depth <= 0:
            depth = 0
            # A declaration ends when its block closes, OR when it was a
            # one-liner that opened and closed on the same line, such as
            # `import ("slices"; "maps")`. Only the first case was handled,
            # so a one-line import swallowed the whole rest of the snippet.
            if current is decls:
                current = None
    return '\n'.join(decls), '\n'.join(rest)


def candidates(code):
    s = strip_lead_comments(code).strip()
    if not s:
        return []
    out = []
    if s.startswith('package '):
        out.append(s)
    out.append('package main\n\n' + s)
    out.append('package main\n\nfunc main() {\n' + s + '\n}\n')
    decls, rest = split_decls(s)
    if decls.strip() and rest.strip():
        out.append('package main\n\n' + decls + '\n\nfunc main() {\n' + rest + '\n}\n')

    # Struct fields on their own. A book about Go shows these constantly:
    #     Tags []Tag `gorm:"many2many:tasks_tags" json:"tags"`
    out.append('package main\n\ntype T struct {\n' + s + '\n}\n')

    # A switch arm on its own, for the same reason.
    out.append('package main\n\nfunc main() {\n\tswitch {\n' + s + '\n\t}\n}\n')

    # Interface methods on their own.
    out.append('package main\n\ntype I interface {\n' + s + '\n}\n')
    return out


def is_continuation(code):
    """True when the snippet is one part of a listing split across blocks.

    A chapter often opens a function in one block, explains the middle in
    prose, and closes it two blocks later. The first has more { than }, the
    last has more } than {. Both are correct and neither parses alone.
    """
    depth = 0
    lowest = 0
    for ch in code:
        if ch == '{':
            depth += 1
        elif ch == '}':
            depth -= 1
            lowest = min(lowest, depth)
    return depth != 0 or lowest < 0

def parses(code):
    with tempfile.NamedTemporaryFile('w', suffix='.go', delete=False,
                                     encoding='utf-8') as fh:
        fh.write(code)
        path = fh.name
    try:
        r = subprocess.run(['gofmt', '-e', path], capture_output=True,
                           text=True, timeout=30)
        return r.returncode == 0, (r.stderr or '').strip()
    finally:
        os.unlink(path)


total = elided = deliberate = continued = commentary = 0
bad = []
for name in sorted(os.listdir(BOOK)):
    if not name.endswith('.typ'):
        continue
    text = open(os.path.join(BOOK, name), encoding='utf-8').read()
    for m in BLOCK.finditer(text):
        code = m.group(1)
        total += 1
        if ELIDED.search(code):
            elided += 1
            continue
        if is_continuation(code):
            continued += 1
            continue
        if not strip_lead_comments(code).strip():
            commentary += 1
            continue
        errs = []
        ok = False
        for cand in candidates(code):
            good, err = parses(cand)
            if good:
                ok = True
                break
            errs.append(err)
        if ok:
            continue
        if DELIBERATE.search(code):
            deliberate += 1
            continue
        line = text[:m.start()].count('\n') + 1
        bad.append((name, line, errs, code))

print(f'{total} Go blocks')
print(f'  {elided} elided, not checked')
print(f'  {deliberate} wrong on purpose (the prose says so), not checked')
print(f'  {continued} part of a listing split across blocks, not checked')
print(f'  {commentary} comment only, nothing to parse')
print(f'  {len(bad)} DID NOT PARSE\n')

for name, line, errs, code in bad:
    print(f'--- {name}:{line}')
    msg = min(errs, key=len) if errs else ''
    for e in msg.splitlines()[:2]:
        print('    ' + re.sub(r'^.*?\.go:', 'line ', e))
    for l in [l for l in code.strip().splitlines() if l.strip()][:4]:
        print(f'    | {l}')
    print()
