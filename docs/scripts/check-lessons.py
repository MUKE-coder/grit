#!/usr/bin/env python3
"""Every lesson advertised as written has a body, and every body is reachable.

Three files have to agree for a lesson to exist:

  config/courses.ts                      says the lesson exists and is written
  lib/lesson-content.ts                  maps its key to an import
  content/lessons/<c>/<ch>/<l>.tsx       is the body

Nothing compared them. loadLessonContent returns null for a key it does not
hold, and the page falls through to <LessonComingSoon />, so a lesson marked
available with no entry renders a "being written" card under a heading that
says it is finished. Nineteen were advertised that way.

The reverse is quieter and also wrong: a lesson with a written body still
marked coming-soon is work nobody can read.

Usage: python3 scripts/check-lessons.py      (from docs/)
"""
import io
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
sys.path.insert(0, HERE)        # -I does not add the script's own directory

from course_catalogue import Shape, lessons, read_catalogue  # noqa: E402

# Both halves of the map: the hand-written one, and the Go Full course's,
# generated from the book beside its lessons.
MAPS = (
    os.path.join(DOCS, 'lib', 'lesson-content.ts'),
    os.path.join(DOCS, 'lib', 'lesson-content.go-full.generated.ts'),
)
BODIES = os.path.join(DOCS, 'content', 'lessons')


def mapped_keys():
    keys = {}
    for path in MAPS:
        src = io.open(path, encoding='utf-8').read()
        for m in re.finditer(
            r"^\s*'([^']+)': \(\) => import\('@/content/lessons/([^']+)'\)", src, re.M
        ):
            keys[m.group(1)] = (m.group(2), os.path.basename(path))
    return keys


def main():
    try:
        catalogue = read_catalogue()
    except Shape as e:
        print('the course config changed shape: %s' % e)
        return 1

    mapped = mapped_keys()
    problems = []
    seen = set()
    total = written = 0

    for key, status, where, _course, _chapter, _lesson in lessons(catalogue):
        total += 1
        if key in seen:
            problems.append('%s: duplicate lesson %s' % (where, key))
        seen.add(key)

        if status == 'MISSING':
            problems.append('%s: %s has no status' % (where, key))
        elif status == 'available':
            written += 1
            if key not in mapped:
                problems.append(
                    '%s: %s is available but has no entry in lesson-content.ts, '
                    'so the page renders the coming-soon card' % (where, key)
                )
                continue
            rel, src = mapped[key]
            body = os.path.join(BODIES, *rel.split('/')) + '.tsx'
            if not os.path.exists(body):
                problems.append(
                    '%s: %s imports %s, which does not exist' % (src, key, rel)
                )
        elif key in mapped:
            problems.append(
                '%s: %s has a written body but is still marked coming-soon, '
                'so nobody can read it' % (where, key)
            )

    for key, (rel, src) in sorted(mapped.items()):
        if key not in seen:
            problems.append(
                '%s: %s is mapped but no course has that lesson, so the body '
                'is unreachable' % (src, key)
            )

    if problems:
        for p in problems:
            print(p)
        print('\n%d problems across %d lessons' % (len(problems), total))
        return 1

    print('%d lessons, %d written, every one with a body' % (total, written))
    return 0


if __name__ == '__main__':
    sys.exit(main())
