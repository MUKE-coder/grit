"""Mirror of routePoints in components/system-flow.tsx, run over every diagram.

Reports an edge whose drawn route would pass through a box it does not join.
"""
import re, glob, os

COL_W, COL_GAP, ROW_H, ROW_GAP, PAD, CLEAR = 220, 56, 72, 60, 16, 6


def box(n):
    x = PAD + n['col'] * (COL_W + COL_GAP)
    y = PAD + n['row'] * (ROW_H + ROW_GAP)
    return dict(x=x, y=y, w=COL_W, h=ROW_H, cx=x + COL_W / 2, cy=y + ROW_H / 2)


def segment_hits(p, q, n):
    B = box(n)
    left, right = B['x'] - CLEAR, B['x'] + B['w'] + CLEAR
    top, bottom = B['y'] - CLEAR, B['y'] + B['h'] + CLEAR
    if p[1] == q[1]:
        if p[1] <= top or p[1] >= bottom:
            return False
        return max(min(p[0], q[0]), left) < min(max(p[0], q[0]), right)
    if p[0] == q[0]:
        if p[0] <= left or p[0] >= right:
            return False
        return max(min(p[1], q[1]), top) < min(max(p[1], q[1]), bottom)
    return False


def route_is_clear(pts, a, b, nodes):
    for n in nodes.values():
        if n['id'] in (a['id'], b['id']):
            continue
        for i in range(1, len(pts)):
            if segment_hits(pts[i - 1], pts[i], n):
                return False
    return True


def route_points(a, b, bend, nodes):
    A, B = box(a), box(b)

    # Straight along a row or down a column still has to be checked: two boxes
    # three columns apart have a box between them.
    if a['row'] == b['row']:
        left = A['x'] < B['x']
        straight = [(A['x'] + A['w'] if left else A['x'], A['cy']),
                    (B['x'] if left else B['x'] + B['w'], B['cy'])]
        if route_is_clear(straight, a, b, nodes):
            return straight, True
        above = A['y'] - ROW_GAP / 2
        below = A['y'] + A['h'] + ROW_GAP / 2
        for lane in ([above, below] if a['row'] > 0 else [below, above]):
            up = lane < A['y']
            detour = [(A['cx'], A['y'] if up else A['y'] + A['h']),
                      (A['cx'], lane),
                      (B['cx'], lane),
                      (B['cx'], B['y'] if up else B['y'] + B['h'])]
            if route_is_clear(detour, a, b, nodes):
                return detour, True
        return straight, False
    if a['col'] == b['col']:
        down = A['y'] < B['y']
        straight = [(A['cx'], A['y'] + A['h'] if down else A['y']),
                    (B['cx'], B['y'] if down else B['y'] + B['h'])]
        if route_is_clear(straight, a, b, nodes):
            return straight, True
        left_lane = A['x'] - COL_GAP / 2
        right_lane = A['x'] + A['w'] + COL_GAP / 2
        for lane in ([left_lane, right_lane] if a['col'] > 0 else [right_lane, left_lane]):
            leftward = lane < A['x']
            detour = [(A['x'] if leftward else A['x'] + A['w'], A['cy']),
                      (lane, A['cy']),
                      (lane, B['cy']),
                      (B['x'] if leftward else B['x'] + B['w'], B['cy'])]
            if route_is_clear(detour, a, b, nodes):
                return detour, True
        return straight, False

    right = A['x'] < B['x']
    down = A['y'] < B['y']
    exit_x = A['x'] + A['w'] if right else A['x']
    exit_y = A['y'] + A['h'] if down else A['y']
    enter_x = B['x'] if right else B['x'] + B['w']
    enter_y = B['y'] if down else B['y'] + B['h']
    lane_x = A['x'] + A['w'] + COL_GAP / 2 if right else A['x'] - COL_GAP / 2
    lane_y = A['y'] + A['h'] + ROW_GAP / 2 if down else A['y'] - ROW_GAP / 2
    near_x = B['x'] - COL_GAP / 2 if right else B['x'] + B['w'] + COL_GAP / 2
    near_y = B['y'] - ROW_GAP / 2 if down else B['y'] + B['h'] + ROW_GAP / 2

    h = [(exit_x, A['cy']), (B['cx'], A['cy']), (B['cx'], enter_y)]
    v = [(A['cx'], exit_y), (A['cx'], B['cy']), (enter_x, B['cy'])]
    h_lane = [(exit_x, A['cy']), (lane_x, A['cy']), (lane_x, B['cy']), (enter_x, B['cy'])]
    v_lane = [(A['cx'], exit_y), (A['cx'], lane_y), (B['cx'], lane_y), (B['cx'], enter_y)]
    h_near = [(exit_x, A['cy']), (near_x, A['cy']), (near_x, B['cy']), (enter_x, B['cy'])]
    v_near = [(A['cx'], exit_y), (A['cx'], near_y), (B['cx'], near_y), (B['cx'], enter_y)]

    candidates = ([v, v_lane, v_near, h, h_lane, h_near] if bend == 'v'
                  else [h, h_lane, h_near, v, v_lane, v_near])
    for pts in candidates:
        if route_is_clear(pts, a, b, nodes):
            return pts, True
    return candidates[0], False


OBJ = re.compile(r'\{[^{}]*\}')
KV = re.compile(r"(\w+)\s*:\s*('(?:[^'\\]|\\.)*'|true|false|\d+)")


def parse_objs(block):
    out = []
    for m in OBJ.finditer(block):
        d = {}
        for k, v in KV.findall(m.group(0)):
            if v in ('true', 'false'):
                d[k] = v == 'true'
            elif v[0] == "'":
                d[k] = v[1:-1]
            else:
                d[k] = int(v)
        if d:
            out.append(d)
    return out


def grab(src, start_idx):
    i = src.index('[', start_idx)
    depth, j = 0, i
    while j < len(src):
        if src[j] == '[':
            depth += 1
        elif src[j] == ']':
            depth -= 1
            if depth == 0:
                return src[i:j + 1], j
        j += 1
    raise ValueError('unbalanced')


# Which SystemDesign field backs each section id the renderer emits, so a
# question that promises an answer "in the data model" can be checked against
# the page actually having a data model.
SECTION_FIELD = {
    'problem': 'problem',
    'requirements': 'functional',
    'capacity': 'capacity',
    'high-level-design': 'highLevel',
    'stack': 'stack',
    'data-model': 'dataModel',
    'api': 'api',
    'low-level-design': 'lowLevel',
    'scaling': 'scaling',
    'bottlenecks': 'bottlenecks',
}


def check_interview(src, path, report):
    """Every interview question must point at a section this page really has.

    A question listed with no section behind it is a promise the page does not
    keep, which is worse than not listing the question at all.
    """
    bad = 0
    for m in re.finditer(r'^export const (\w+): SystemDesign = \{', src, re.M):
        start = m.start()
        nxt = re.search(r'^export const \w+: SystemDesign = \{', src[start + 10:], re.M)
        block = src[start: start + 10 + nxt.start()] if nxt else src[start:]
        name = m.group(1)
        if 'interview:' not in block:
            continue
        for see in re.findall(r"\bsee:\s*'([^']+)'", block):
            field = SECTION_FIELD.get(see)
            if field is None:
                report('INTERVIEW   %-24s [%s] unknown section id "%s"' % (
                    os.path.basename(path), name, see))
                bad += 1
            elif not re.search(r'^  %s[?]?:' % re.escape(field), block, re.M):
                report('INTERVIEW   %-24s [%s] points at "%s" but the page has no %s' % (
                    os.path.basename(path), name, see, field))
                bad += 1
    return bad


problems = 0
diagrams = 0
rerouted = 0
for path in sorted(glob.glob('config/systems*.ts')):
    src = open(path, encoding='utf-8').read()
    problems += check_interview(src, path, print)
    for m in re.finditer(r'\bnodes:\s*\[', src):
        diagrams += 1
        ntext, nend = grab(src, m.start())
        etext, _ = grab(src, src.index('edges:', nend))
        nodes = {n['id']: n for n in parse_objs(ntext) if 'id' in n}
        edges = [e for e in parse_objs(etext) if 'from' in e]
        head = src.rfind('title:', 0, m.start())
        t = re.search(r"title:\s*'([^']*)'", src[head:head + 200])
        title = t.group(1) if t else '?'

        seen_steps = []
        for e in edges:
            if 'step' in e:
                seen_steps.append(e['step'])
            a, b = nodes.get(e['from']), nodes.get(e['to'])
            if not a or not b:
                print('%s [%s] edge %s references a missing node' % (path, title, e))
                problems += 1
                continue
            pts, clear = route_points(a, b, e.get('bend'), nodes)
            direct, _ = route_points(a, b, e.get('bend'), {})
            if pts != direct:
                rerouted += 1
            if not clear:
                print('UNROUTABLE  %-24s [%s] step %s: %s -> %s' % (
                    os.path.basename(path), title, e.get('step', '-'), e['from'], e['to']))
                problems += 1

        want = list(range(1, len(seen_steps) + 1))
        if sorted(seen_steps) != want:
            print('STEPS       %-24s [%s] numbered %s' % (
                os.path.basename(path), title, sorted(seen_steps)))
            problems += 1

print('')
print('%d diagrams, %d edges rerouted around a box, %d problems' % (diagrams, rerouted, problems))
