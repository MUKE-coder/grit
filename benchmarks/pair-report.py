"""Turn one pair's results into the numbers the docs site publishes.

    python pair-report.py bun
    python pair-report.py            # every pair measured so far

Reports the MEDIAN of the repetitions, never the best run. Publishing a best run
is how a benchmark becomes something nobody else can reproduce: whoever re-runs
it gets a lower number and concludes, correctly, that the published one was
cherry-picked.

Carries the CPU samples through, because they are what separates a real ceiling
from the point where the database gave out. A row where the app sat at 118%
while Postgres was pinned at 823% is not "the framework does 769 req/s" — it is
"at least 769, with the database as the limit".
"""

import glob
import io
import json
import os
import re
import statistics
import sys
from collections import defaultdict

ROOT = os.path.dirname(os.path.abspath(__file__))
PAIRS_DIR = os.path.join(ROOT, "results", "pairs")
SCENARIOS = ["show", "write", "list", "mixed"]

RPS = re.compile(r"http_reqs[.]+:\s+(\d+)\s+([0-9.]+)/s")
DUR = re.compile(
    r"http_req_duration[.]+:\s+avg=(\S+)\s+min=(\S+)\s+med=(\S+)\s+max=(\S+)"
    r"\s+p\(90\)=(\S+)\s+p\(95\)=(\S+)"
)
FAIL = re.compile(r"failed_requests[.]+:\s+([0-9.]+)%")


def ms(value):
    """k6 prints 1.2s / 340ms / 687.46µs — normalise to milliseconds."""
    value = value.strip()
    for suffix, mult in (("ms", 1.0), ("µs", 0.001), ("us", 0.001), ("s", 1000.0)):
        if value.endswith(suffix):
            try:
                return float(value[: -len(suffix)]) * mult
            except ValueError:
                return None
    return None


def read_log(path):
    txt = io.open(path, encoding="utf-8", errors="replace").read()
    rps, dur, fail = RPS.search(txt), DUR.search(txt), FAIL.search(txt)
    if not rps:
        return None
    return {
        "rps": float(rps.group(2)),
        "med": ms(dur.group(3)) if dur else None,
        "p95": ms(dur.group(6)) if dur else None,
        "fail": (float(fail.group(1)) / 100.0) if fail else 0.0,
    }


def read_cpu(path, app):
    out = {}
    if not os.path.exists(path):
        return out
    for line in io.open(path, encoding="utf-8", errors="replace"):
        parts = line.split()
        if len(parts) == 2 and parts[1].endswith("%"):
            name = re.sub(r"^benchmarks-|-\d+$", "", parts[0])
            key = "app" if name == app else name
            out[key] = float(parts[1].rstrip("%"))
    return out


def collect(pair_dir, app):
    rows = defaultdict(list)
    for path in sorted(glob.glob(os.path.join(pair_dir, "%s-*.log" % app))):
        base = os.path.basename(path)[: -len(".log")]
        m = re.match(r"^%s-(show|write|list|mixed)-(\d+)$" % re.escape(app), base)
        if not m:
            continue
        row = read_log(path)
        if not row:
            continue
        row["cpu"] = read_cpu(path[: -len(".log")] + ".cpu", app)
        rows[m.group(1)].append(row)
    return rows


def med(values):
    values = [v for v in values if v is not None]
    return statistics.median(values) if values else None


def summarise(rows):
    out = {}
    for scenario, reps in rows.items():
        if not reps:
            continue
        out[scenario] = {
            "n": len(reps),
            "rps": med([r["rps"] for r in reps]),
            "med_ms": med([r["med"] for r in reps]),
            "p95_ms": med([r["p95"] for r in reps]),
            "worst_fail": max(r["fail"] or 0 for r in reps),
            "app_cpu": med([r["cpu"].get("app") for r in reps]),
            "db_cpu": med([r["cpu"].get("postgres") for r in reps]),
        }
    return out


def bound(entry):
    """A number is only a ceiling when the app itself is what ran out."""
    app, db = entry.get("app_cpu") or 0, entry.get("db_cpu") or 0
    if app > 350:
        return "app-bound (a real ceiling)"
    if db > 700:
        return "DB-bound (a floor — could go higher)"
    return "neither saturated"


def report(opponent):
    pair_dir = os.path.join(PAIRS_DIR, opponent)
    grit = summarise(collect(pair_dir, "grit"))
    other = summarise(collect(pair_dir, opponent))
    if not grit or not other:
        print("  no complete data for %s yet" % opponent)
        return None

    print("== grit vs %s ==" % opponent)
    print("  %-7s %10s %10s %8s %10s %10s  %s"
          % ("scen", "grit r/s", opponent[:10] + " r/s", "ratio",
             "grit med", "%s med" % opponent[:6], "bound by (grit)"))

    for scenario in SCENARIOS:
        g, o = grit.get(scenario), other.get(scenario)
        if not (g and o):
            continue
        warn = ""
        if max(g["worst_fail"], o["worst_fail"]) > 0.001:
            warn = "  !! failures above zero — do not publish"
        print("  %-7s %10.0f %10.0f %7.2fx %9.1fms %9.1fms  %s%s"
              % (scenario, g["rps"], o["rps"], g["rps"] / o["rps"],
                 g["med_ms"] or 0, o["med_ms"] or 0, bound(g), warn))

    print("  reps: %s" % ", ".join("%s=%d" % (s, grit[s]["n"]) for s in SCENARIOS if s in grit))
    print()
    return {"grit": grit, opponent: other}


def main():
    wanted = sys.argv[1:] or sorted(
        d for d in os.listdir(PAIRS_DIR) if os.path.isdir(os.path.join(PAIRS_DIR, d))
    ) if os.path.isdir(PAIRS_DIR) else []

    combined = {}
    for opponent in wanted:
        result = report(opponent)
        if result:
            combined[opponent] = result

    if combined:
        out = os.path.join(PAIRS_DIR, "summary.json")
        with io.open(out, "w", encoding="utf-8") as fh:
            json.dump(combined, fh, indent=2)
        print("wrote %s" % out)


if __name__ == "__main__":
    main()
