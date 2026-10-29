#!/usr/bin/env python3
"""Read the course catalogue out of the two config files that describe it.

One reader, used by everything that needs to know what courses exist:

    scripts/check-lessons.py    every written lesson has a body
    scripts/search-index.py     every course page is findable

Two files describe the catalogue. config/courses.ts is hand-written;
config/go-full-course.generated.ts is built from the book by
scripts/go-course-build.py. They nest identically at different indentation,
so this works off the array markers rather than columns.

Nothing here validates. It parses and returns, and the callers decide what
counts as wrong. What it will not do is guess: a shape it does not recognise
raises Shape, because a parser that silently reads half a file produces a
check that silently passes.
"""
import io
import os
import re

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)

CONFIGS = (
    os.path.join(DOCS, 'config', 'courses.ts'),
    os.path.join(DOCS, 'config', 'go-full-course.generated.ts'),
)

# Any key whose array opens and does not close on the same line. Tracking all
# of them, not just the three that matter, keeps the stack honest: a
# successCriteria array that pushed nothing would pop a chapters array on its
# way out.
OPENS = re.compile(r'^(\s*)(\w+): \[\s*$')
CLOSES = re.compile(r'^(\s*)\],')
SLUG = re.compile(r"^(\s*)(?:\{ )?slug: '([^']+)'")
ASSIGNMENT = re.compile(r'^(\s*)assignment: \{\s*$')
FIELD = re.compile(r"\b(name|shortName|title|tagline|status|brief): "
                   r"""(?:'((?:[^'\\]|\\.)*)'|"((?:[^"\\]|\\.)*)")""")

NESTS = ('chapters', 'modules', 'lessons')


class Shape(Exception):
    """The config changed shape, so the parse cannot be trusted."""


def _unescape(s):
    return s.replace("\\'", "'").replace('\\"', '"')


def _fields(lines, n, indent, wanted):
    """The fields of the object that opens on line n, at that indentation."""
    found = {}

    def take(line):
        for m in FIELD.finditer(line):
            found.setdefault(m.group(1), _unescape(m.group(2) or m.group(3) or ''))

    take(lines[n])
    if found and all(k in found for k in wanted):
        return found                             # written on one line
    end = re.compile(r'^\s{%d}\}' % max(indent - 2, 0))
    for line in lines[n + 1:n + 40]:
        if end.match(line):
            break
        take(line)
        if all(k in found for k in wanted):
            break
    return found


def _read(path):
    """[{slug, name, shortName, tagline, chapters: [...]}] from one config."""
    lines = io.open(path, encoding='utf-8').read().splitlines()
    # The interface declarations at the top of courses.ts have slug lines too.
    first = next(
        (n for n, l in enumerate(lines)
         if re.match(r'^(?:const|export const) \w+', l) and ': Course' in l),
        0,
    )

    stack = []       # (key, indent) per open multi-line array
    courses = []
    for n, line in enumerate(lines[first:], start=first):
        where = '%s:%d' % (os.path.basename(path), n + 1)

        m = OPENS.match(line)
        if m:
            stack.append((m.group(2), len(m.group(1))))
            continue
        m = CLOSES.match(line)
        if m and stack and len(m.group(1)) == stack[-1][1]:
            stack.pop()
            continue

        m = ASSIGNMENT.match(line)
        if m:
            if not courses or not courses[-1]['chapters']:
                raise Shape('%s: an assignment outside a chapter' % where)
            f = _fields(lines, n, len(m.group(1)) + 2, ('title', 'brief'))
            courses[-1]['chapters'][-1]['assignment'] = {
                'title': f.get('title', 'Assignment'),
                'brief': re.sub(r'\s+', ' ', f.get('brief', '')),
            }
            continue

        m = SLUG.match(line)
        if not m:
            continue
        indent, slug = len(m.group(1)), m.group(2)
        kind = stack[-1][0] if stack else 'course'

        if kind == 'course':
            f = _fields(lines, n, indent, ('name', 'shortName', 'tagline'))
            courses.append({
                'slug': slug,
                'name': f.get('name', slug),
                'shortName': f.get('shortName') or f.get('name', slug),
                'tagline': f.get('tagline', ''),
                'where': where,
                'chapters': [],
            })
        elif kind == 'chapters':
            if not courses:
                raise Shape('%s: chapter %s before any course' % (where, slug))
            f = _fields(lines, n, indent, ('title', 'tagline'))
            courses[-1]['chapters'].append({
                'slug': slug,
                'title': f.get('title', slug),
                'tagline': f.get('tagline', ''),
                'where': where,
                'assignment': None,
                'lessons': [],
            })
        elif kind == 'lessons':
            if not courses or not courses[-1]['chapters']:
                raise Shape('%s: lesson %s outside a chapter' % (where, slug))
            f = _fields(lines, n, indent, ('title', 'tagline', 'status'))
            courses[-1]['chapters'][-1]['lessons'].append({
                'slug': slug,
                'title': f.get('title', slug),
                'tagline': f.get('tagline', ''),
                'status': f.get('status', 'MISSING'),
                'where': where,
            })
        else:
            # A slug only ever belongs to a course, a chapter or a lesson.
            # Module objects carry a title and nothing else identifying.
            raise Shape('%s: slug inside a %s array: %s' % (where, kind, slug))

    if not courses:
        raise Shape('%s: no courses found' % os.path.basename(path))
    return courses


def read_catalogue():
    """Every course, in the order the two configs declare them."""
    out = []
    for path in CONFIGS:
        out.extend(_read(path))
    return out


def lessons(catalogue=None):
    """Flat: (key, status, where, course, chapter, lesson) per lesson."""
    for course in catalogue if catalogue is not None else read_catalogue():
        for chapter in course['chapters']:
            for lesson in chapter['lessons']:
                key = '%s/%s/%s' % (course['slug'], chapter['slug'], lesson['slug'])
                yield key, lesson['status'], lesson['where'], course, chapter, lesson


if __name__ == '__main__':
    cat = read_catalogue()
    chapters = [ch for c in cat for ch in c['chapters']]
    print('%d courses, %d chapters, %d assignments, %d lessons (%d written)' % (
        len(cat), len(chapters),
        sum(1 for ch in chapters if ch['assignment']),
        sum(len(ch['lessons']) for ch in chapters),
        sum(1 for ch in chapters for l in ch['lessons']
            if l['status'] == 'available'),
    ))
    for c in cat:
        print('  %-14s %-40s %-18s %d chapters'
              % (c['slug'], c['name'][:40], c['shortName'][:18], len(c['chapters'])))
