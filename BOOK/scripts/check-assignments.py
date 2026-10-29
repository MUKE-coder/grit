"""The assignments, checked against the rules that make them scorable.

Every chapter ends with one #assignment, and the Committed app scores a
reader's work by EXACT commit name. That makes the name a key, not a label,
and a key has rules:

  unique      two chapters sharing a name means one of them cannot score.
              Chapters 21 and 54 both used "HTTP Server - Assignment One"
              until this check was written.
  present     a chapter with no assignment is a chapter with nothing on the
              leaderboard, which is how a reader quietly stops.
  one each    the book's own convention, so a reader always knows what is
              being asked of them. Chapter 3 has two and says why.
  consistent  "<Topic> - Assignment One", so the name can be predicted from
              the chapter rather than looked up.
"""
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
CHAPTERS = os.path.join(os.path.dirname(HERE), 'chapters')

ASSIGNMENT = re.compile(r'#assignment\(name:\s*"([^"]*)"')
CHAPTER_NUM = re.compile(r'^(\d\d)-')

# Appendices teach nothing and are assigned nothing.
NO_ASSIGNMENT_EXPECTED = {90, 91}

# Chapters allowed more than one, with the reason. An entry here is a decision
# somebody made, not a gap somebody left.
MULTIPLE_ALLOWED = {
    3: 'variables and types carry two exercises: declaring, then zero values',
}

NAME_SHAPE = re.compile(r'^[A-Z][\w\- ]* - Assignment (One|Two|Three)$')
# The first chapter sets up the repository rather than exercising a topic.
NAME_EXCEPTIONS = {'Setup - Repository'}


def main():
    problems = []
    seen = {}
    total = 0

    for name in sorted(os.listdir(CHAPTERS)):
        if not name.endswith('.typ'):
            continue
        m = CHAPTER_NUM.match(name)
        if not m:
            continue
        number = int(m.group(1))
        text = open(os.path.join(CHAPTERS, name), encoding='utf-8').read()
        names = ASSIGNMENT.findall(text)
        total += len(names)

        if not names:
            if number not in NO_ASSIGNMENT_EXPECTED:
                problems.append(
                    f'{name}: no assignment, so nothing from this chapter can '
                    f'reach the leaderboard')
            continue

        if number in NO_ASSIGNMENT_EXPECTED:
            problems.append(f'{name}: an appendix with an assignment in it')

        if len(names) > 1 and number not in MULTIPLE_ALLOWED:
            problems.append(
                f'{name}: {len(names)} assignments. One per chapter, or add it '
                f'to MULTIPLE_ALLOWED with the reason')

        for a in names:
            if a in seen:
                problems.append(
                    f'{name}: assignment name "{a}" is already used by '
                    f'{seen[a]}. Committed scores by exact commit name, so one '
                    f'of the two cannot score')
            else:
                seen[a] = name
            if a not in NAME_EXCEPTIONS and not NAME_SHAPE.match(a):
                problems.append(
                    f'{name}: assignment name "{a}" does not read '
                    f'"<Topic> - Assignment One"')

    print(f'{total} assignments across {len(seen)} distinct names')
    if not problems:
        print('all good')
        return 0
    print()
    for p in problems:
        print('  ' + p)
    return 1


if __name__ == '__main__':
    sys.exit(main())
