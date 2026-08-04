"""Turn results/*.json into the table the docs page publishes.

Reports the MEDIAN of the repetitions, not the best run. Picking the best is how
benchmarks end up unreproducible: whoever re-runs it gets a lower number and
concludes the published one was cherry-picked, which it was.

Also carries the CPU samples through. They are the part that decides whether a
row is a real ceiling or just where the database gave out, and that distinction
matters more than the number itself.
"""

import glob
import io
import json
import os
import re
import statistics
from collections import defaultdict

RESULTS = os.path.join(os.path.dirname(os.path.abspath(__file__)),
                       os.environ.get("RESULTS_DIR", "results/final"))
SCENARIOS = ["list", "show", "mixed", "write"]
APPS = ["grit", "encore", "bun", "express", "nextjs", "laravel", "django"]


def metric(summary, name, field):
    m = summary.get("metrics", {}).get(name, {})
    return m.get(field)


LOGNUM = re.compile(r"http_reqs[.]+:\s+(\d+)\s+([0-9.]+)/s")
DUR = re.compile(r"http_req_duration[.]+:\s+avg=(\S+)\s+min=(\S+)\s+med=(\S+)\s+max=(\S+)\s+p\(90\)=(\S+)\s+p\(95\)=(\S+)")
FAILR = re.compile(r"failed_requests[.]+:\s+([0-9.]+)%")


def ms(v):
    """k6 prints 1.2s / 340ms / 687.46us — normalise to milliseconds."""
    v = v.strip()
    for suffix, mult in (("ms", 1.0), ("us", 0.001), ("µs", 0.001), ("s", 1000.0)):
        if v.endswith(suffix):
            try:
                return float(v[: -len(suffix)]) * mult
            except ValueError:
                return None
    return None


def from_log(path):
    txt = io.open(path, encoding="utf-8", errors="replace").read()
    n = LOGNUM.search(txt)
    d = DUR.search(txt)
    f = FAILR.search(txt)
    if not n:
        return None
    return {
        "count": int(n.group(1)),
        "rps": float(n.group(2)),
        "med": ms(d.group(3)) if d else None,
        "p95": ms(d.group(6)) if d else None,
        "p99": None,
        "max": ms(d.group(4)) if d else None,
        "fail": (float(f.group(1)) / 100.0) if f else 0.0,
    }


def load():
    runs = defaultdict(list)
    for path in sorted(glob.glob(os.path.join(RESULTS, "*.log"))):
        base = os.path.basename(path)[: -len(".log")]
        m = re.match(r"^(grit|encore|bun|express|nextjs|laravel|django)-(list|show|mixed|write)-(\d+)$", base)
        if not m:
            continue
        app, scenario, _rep = m.groups()
        row = from_log(path)
        if not row:
            continue
        row["cpu"] = read_cpu(app, scenario, _rep)
        runs[(app, scenario)].append(row)
    return runs


def read_cpu(app, scenario, rep):
    path = os.path.join(RESULTS, "%s-%s-%s.cpu" % (app, scenario, rep))
    out = {}
    if not os.path.exists(path):
        return out
    for line in io.open(path, encoding="utf-8"):
        parts = line.split()
        if len(parts) == 2 and parts[1].endswith("%"):
            key = re.sub(r"^benchmarks-|-\d+$", "", parts[0])
            out[key] = float(parts[1].rstrip("%"))
    return out


def med(values):
    values = [v for v in values if v is not None]
    return statistics.median(values) if values else None


def main():
    runs = load()
    if not runs:
        print("no results — run ./run.sh first")
        return

    table = {}
    for app in APPS:
        for scenario in SCENARIOS:
            rs = runs.get((app, scenario), [])
            if not rs:
                continue
            table[(app, scenario)] = {
                "n": len(rs),
                "rps": med([r["rps"] for r in rs]),
                "med": med([r["med"] for r in rs]),
                "p95": med([r["p95"] for r in rs]),
                "p99": med([r["p99"] for r in rs]),
                "fail": max(r["fail"] or 0 for r in rs),
                "app_cpu": med([r["cpu"].get(app) for r in rs]),
                "pg_cpu": med([r["cpu"].get("postgres") for r in rs]),
                "spread": (max(r["rps"] for r in rs) - min(r["rps"] for r in rs))
                / med([r["rps"] for r in rs])
                * 100,
            }

    print("%-8s %-7s %3s %10s %9s %9s %9s %7s %9s %9s %7s" % (
        "app", "scen", "n", "rps", "med ms", "p95 ms", "p99 ms",
        "fail%", "app cpu%", "pg cpu%", "spread%"))
    for scenario in SCENARIOS:
        for app in APPS:
            r = table.get((app, scenario))
            if not r:
                continue
            print("%-8s %-7s %3d %10.1f %9.1f %9.1f %9.1f %7.2f %9.0f %9.0f %7.1f" % (
                app, scenario, r["n"], r["rps"], r["med"], r["p95"],
                r["p99"] or 0, (r["fail"] or 0) * 100,
                r["app_cpu"] or 0, r["pg_cpu"] or 0, r["spread"]))
        print()

    print("ratios (grit / laravel), median rps:")
    for scenario in SCENARIOS:
        g = table.get(("grit", scenario))
        l = table.get(("laravel", scenario))
        if not (g and l):
            continue
        bound = "DB-bound — Grit is a FLOOR" if (g["pg_cpu"] or 0) > 700 and (g["app_cpu"] or 0) < 350 \
            else "app-bound"
        print("  %-6s %5.1fx   (%s)" % (scenario, g["rps"] / l["rps"], bound))

    with io.open(os.path.join(RESULTS, "summary.json"), "w", encoding="utf-8") as fh:
        json.dump({"%s-%s" % k: v for k, v in table.items()}, fh, indent=2)
    print("\nwrote results/summary.json")


if __name__ == "__main__":
    main()
