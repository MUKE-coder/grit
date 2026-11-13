package doctor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultQueryLog is where internal/dbchain writes when GRIT_QUERY_LOG points
// at the conventional place.
const DefaultQueryLog = ".grit/queries.jsonl"

// repeatThreshold is how many times one query shape has to run inside a single
// request before it is worth saying anything.
//
// Five, not two. A handler legitimately loads the actor, then a role, then a
// parent row, and warning about three is a warning nobody can act on. An N+1
// does not look like three: it looks like one per row, which on any list a
// person would notice is a dozen at least.
const repeatThreshold = 5

// queryEntry is one line of the log. The fields internal/dbchain writes; a
// line with more is read fine, which is what lets the recorder grow without
// this needing to know.
type queryEntry struct {
	At      time.Time `json:"at"`
	Request string    `json:"request"`
	Route   string    `json:"route"`
	Table   string    `json:"table"`
	SQL     string    `json:"sql"`
	Caller  string    `json:"caller"`
}

// repeat is one query shape, from one line of the project, inside one request.
type repeat struct {
	Route  string
	SQL    string
	Caller string
	Count  int
	// Requests is how many distinct requests showed this, which separates
	// "one unlucky request" from "every request on this route".
	Requests int
}

// The same query, many times, inside one request.
//
// This is the only check here that reads something a run produced rather than
// something the repository contains, and it is deliberately still offline: the
// application writes a file, and this reads the file. Nothing here talks to a
// database or a server, which is what keeps `grit doctor` usable in CI.
//
// To produce the file:
//
//	GRIT_QUERY_LOG=.grit/queries.jsonl go test ./...
//	GRIT_QUERY_LOG=.grit/queries.jsonl go run ./cmd/server   # then use the app
//
// Silent when there is no file. A project that has never recorded one is not
// doing anything wrong, and a check that nags about a missing artefact is how
// people learn to stop reading the output.
func checkQueryLog(p *project) []Finding {
	path := p.queryLog
	if path == "" {
		path = filepath.Join(p.root, DefaultQueryLog)
	}
	if _, err := os.Stat(path); err != nil {
		return nil // no log: nothing to say
	}

	entries, newest, parseErr := readQueryLog(path)
	if parseErr != nil {
		return []Finding{{
			Level:    "warning",
			Resource: shortPath(p.root, path),
			Message:  "is not a query log this can read: " + parseErr.Error(),
			Fix:      "delete it and record a new one with GRIT_QUERY_LOG=" + DefaultQueryLog,
		}}
	}
	if len(entries) == 0 {
		return nil
	}

	var out []Finding

	// Does the log describe the code as it is now? A report from before a
	// change is a report about code that no longer exists, and the findings
	// read exactly the same either way, which is the whole danger of an
	// artefact on disk.
	if stale, newestCode := codeIsNewerThan(p, newest); stale {
		out = append(out, Finding{
			Level:    "warning",
			Resource: shortPath(p.root, path),
			Message: fmt.Sprintf("was recorded %s and the project's code is newer (%s), "+
				"so what follows may be about queries that no longer exist",
				newest.Local().Format("2 Jan 15:04"), newestCode.Local().Format("2 Jan 15:04")),
			Fix: "record a new one: GRIT_QUERY_LOG=" + DefaultQueryLog + " go test ./...",
		})
	}

	for _, r := range repeatedQueries(entries) {
		where := r.Caller
		if where == "" {
			where = "an unknown line"
		}
		times := "once"
		if r.Requests > 1 {
			times = fmt.Sprintf("in %d requests", r.Requests)
		}
		out = append(out, Finding{
			Level:    "warning",
			Resource: r.Route,
			Message: fmt.Sprintf("runs the same query %d times in one request, from %s (%s): %s",
				r.Count, where, times, firstLine(r.SQL)),
			Fix: "load them in one query instead: Preload the relation, or fetch the ids " +
				"and use a single IN (...). If it is deliberate, this is the line to put a comment on.",
		})
	}
	return out
}

// readQueryLog streams the file. A test suite writes one line per query, which
// on a real project is hundreds of thousands: this never holds the file.
func readQueryLog(path string) ([]queryEntry, time.Time, error) {
	f, err := os.Open(path) // #nosec G304 -- a path inside the project, or one the user named
	if err != nil {
		return nil, time.Time{}, err
	}
	defer func() { _ = f.Close() }()

	var (
		entries []queryEntry
		newest  time.Time
	)
	scanner := bufio.NewScanner(f)
	// A long IN (...) makes a long line, and the default 64 KiB limit would
	// stop the scan in the middle with no error anybody reads.
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}
		var e queryEntry
		if err := json.Unmarshal([]byte(raw), &e); err != nil {
			return nil, time.Time{}, fmt.Errorf("line %d is not JSON", line)
		}
		// Only what a request issued. A migration, a seeder and a worker all
		// query, and "the same query many times" is what they are for.
		if e.Request == "" {
			continue
		}
		if e.At.After(newest) {
			newest = e.At
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, time.Time{}, err
	}
	return entries, newest, nil
}

// repeatedQueries groups by request, then by shape and the line that issued
// it, and keeps the worst case per (route, shape, caller).
func repeatedQueries(entries []queryEntry) []repeat {
	type key struct{ route, sql, caller string }

	// request -> key -> count
	perRequest := map[string]map[key]int{}
	for _, e := range entries {
		if perRequest[e.Request] == nil {
			perRequest[e.Request] = map[key]int{}
		}
		perRequest[e.Request][key{e.Route, e.SQL, e.Caller}]++
	}

	worst := map[key]*repeat{}
	for _, counts := range perRequest {
		for k, n := range counts {
			if n < repeatThreshold {
				continue
			}
			r, seen := worst[k]
			if !seen {
				r = &repeat{Route: k.route, SQL: k.sql, Caller: k.caller}
				worst[k] = r
			}
			r.Requests++
			if n > r.Count {
				r.Count = n
			}
		}
	}

	out := make([]repeat, 0, len(worst))
	for _, r := range worst {
		if r.Route == "" {
			r.Route = "a request"
		}
		out = append(out, *r)
	}
	// Worst first, then stable by route and line, so two runs of the same log
	// report in the same order.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		if out[i].Route != out[j].Route {
			return out[i].Route < out[j].Route
		}
		return out[i].Caller < out[j].Caller
	})
	return out
}

// codeIsNewerThan reports whether anything under the API's internal/ has been
// written since the log was recorded.
func codeIsNewerThan(p *project, recorded time.Time) (bool, time.Time) {
	if recorded.IsZero() {
		return false, time.Time{}
	}
	var newest time.Time
	root := filepath.Join(p.apiRoot, "internal")
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil //nolint:nilerr // an unreadable file is not this check's problem
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
		return nil
	})
	// A minute of slack: a log written during a test run is a few seconds
	// older than the files that run compiled, and reporting that as drift
	// would make the warning permanent and therefore worthless.
	return newest.After(recorded.Add(time.Minute)), newest
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 120 {
		return s[:117] + "..."
	}
	return s
}

func shortPath(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(path)
}
