"""Does the book use a thing before it teaches it?

This is the check that matters most for a beginner, and the one no human
reviewer does reliably, because catching it means holding 61 chapters in your
head at once. A reader who meets `go func()` in chapter 5 when goroutines
arrive in chapter 13 does not know they are allowed to not understand it. They
assume they have missed something and lose confidence, which is the single
most expensive thing a teaching book can do.

Each concept is a pattern plus the chapter that owns it. A hit in an earlier
chapter is a forward reference. Prose mentions are fine and expected ("we will
come to goroutines"), so only CODE is scanned.
"""
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
BOOK = os.path.join(os.path.dirname(HERE), 'chapters')
BLOCK = re.compile(r'```go\n(.*?)```', re.S)

# concept -> (regex over code, chapter number that teaches it)
CONCEPTS = {
    'goroutines':      (r'\bgo\s+func\b|\bgo\s+\w+\(', 13),
    'channels':        (r'\bchan\b', 14),
    'sync.WaitGroup':  (r'\bsync\.WaitGroup\b', 15),
    'sync.Mutex':      (r'\bsync\.(Mutex|RWMutex)\b', 15),
    'context':         (r'\bcontext\.(Context|Background|WithCancel|WithTimeout)\b', 16),
    'generics':        (r'\[[A-Z]\w*\s+(any|comparable|cmp\.Ordered)\]', 17),
    'testing':         (r'\bfunc Test\w+\(t \*testing\.T\)', 18),
    'json tags':       (r'`json:"', 19),
    'defer':           (r'^\s*defer\b', 12),
    'errors.Is/As':    (r'\berrors\.(Is|As)\b', 11),
    'interfaces':      (r'^\s*type \w+ interface\b', 10),
    'pointers':        (r'(?<![\w\)\]])\*\w+\s*\{|\&\w+\{', 8),
    'maps':            (r'\bmap\[', 7),
    'slices':          (r'\[\]\w+\{|\bappend\(', 6),
    'structs':         (r'^\s*type \w+ struct\b', 5),
}

CHAPTER = re.compile(r'^(\d\d)-')


def chapter_number(name):
    m = CHAPTER.match(name)
    return int(m.group(1)) if m else None


findings = []
for name in sorted(os.listdir(BOOK)):
    if not name.endswith('.typ'):
        continue
    n = chapter_number(name)
    if n is None or n >= 90:          # appendices teach nothing in order
        continue
    text = open(os.path.join(BOOK, name), encoding='utf-8').read()
    code = '\n'.join(m.group(1) for m in BLOCK.finditer(text))
    if not code.strip():
        continue
    for concept, (pattern, taught) in CONCEPTS.items():
        if n >= taught:
            continue
        m = re.search(pattern, code, re.M)
        if m:
            line = code[:m.start()].count('\n') + 1
            snippet = code.splitlines()[line - 1].strip()[:88]
            findings.append((n, name, concept, taught, snippet))

findings.sort()
print(f'{len(findings)} forward reference(s): code using a concept the book '
      f'has not taught yet\n')
for n, name, concept, taught, snippet in findings:
    print(f'  ch{n:02d} {name}')
    print(f'      uses {concept}, which chapter {taught} teaches')
    print(f'      | {snippet}')
    print()
