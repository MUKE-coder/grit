package doctor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeLog writes a query log the way internal/dbchain writes one: JSON
// Lines, one statement per line.
func writeLog(t *testing.T, root string, entries []queryEntry) string {
	t.Helper()
	dir := filepath.Join(root, ".grit")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "queries.jsonl")
	var b strings.Builder
	for _, e := range entries {
		raw, err := json.Marshal(e)
		if err != nil {
			t.Fatal(err)
		}
		b.Write(raw)
		b.WriteByte('\n')
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func logProject(t *testing.T) *project {
	t.Helper()
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api", "internal")
	if err := os.MkdirAll(api, 0o755); err != nil {
		t.Fatal(err)
	}
	// One source file, dated well before anything the tests record, so the
	// staleness warning stays out of the way unless a test asks for it.
	src := filepath.Join(api, "thing.go")
	if err := os.WriteFile(src, []byte("package internal\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(src, old, old); err != nil {
		t.Fatal(err)
	}
	return &project{root: root, apiRoot: filepath.Join(root, "apps", "api")}
}

func repeated(route, sql, caller string, n int) []queryEntry {
	var out []queryEntry
	for i := 0; i < n; i++ {
		out = append(out, queryEntry{
			At: time.Now(), Request: "r1", Route: route, SQL: sql, Caller: caller,
		})
	}
	return out
}

// The same query, many times, inside one request. The whole point.
func TestTheSameQueryManyTimesInOneRequestIsReported(t *testing.T) {
	p := logProject(t)
	writeLog(t, p.root, repeated(
		"GET /api/v1/posts",
		"SELECT * FROM `users` WHERE id = ?",
		"internal/services/post.go:88", 23))

	findings := checkQueryLog(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %+v", len(findings), findings)
	}
	f := findings[0]
	if f.Resource != "GET /api/v1/posts" {
		t.Errorf("the finding does not name the route: %q", f.Resource)
	}
	// The count and the line are the two things that make this actionable.
	if !strings.Contains(f.Message, "23 times") {
		t.Errorf("the finding does not say how many: %q", f.Message)
	}
	if !strings.Contains(f.Message, "internal/services/post.go:88") {
		t.Errorf("the finding does not name the line that issued it: %q", f.Message)
	}
}

// A handful is not an N+1. A handler legitimately loads the actor, a role and
// a parent row, and warning about three is a warning nobody can act on.
func TestAFewRepeatsAreNotReported(t *testing.T) {
	p := logProject(t)
	writeLog(t, p.root, repeated(
		"GET /api/v1/posts", "SELECT * FROM `users` WHERE id = ?", "post.go:88", 3))
	if findings := checkQueryLog(p); len(findings) != 0 {
		t.Errorf("three repeats were reported: %q", findings[0].Message)
	}
}

// Queries outside a request are not an N+1: a migration, a seeder and a
// worker all run the same statement many times, and that is what they do.
func TestQueriesWithNoRequestAreIgnored(t *testing.T) {
	p := logProject(t)
	var entries []queryEntry
	for i := 0; i < 50; i++ {
		entries = append(entries, queryEntry{
			At: time.Now(), SQL: "INSERT INTO `rows` VALUES (?)", Caller: "cmd/seed/main.go:20",
		})
	}
	writeLog(t, p.root, entries)
	if findings := checkQueryLog(p); len(findings) != 0 {
		t.Errorf("a seeder was reported as an N+1: %q", findings[0].Message)
	}
}

// Two requests each running it six times is one finding, not two, and it says
// so: "in 2 requests" separates a pattern from one unlucky call.
func TestTheWorstCaseIsReportedOncePerPlace(t *testing.T) {
	p := logProject(t)
	var entries []queryEntry
	for _, req := range []string{"r1", "r2"} {
		for i := 0; i < 6; i++ {
			entries = append(entries, queryEntry{
				At: time.Now(), Request: req, Route: "GET /api/v1/posts",
				SQL: "SELECT * FROM `users` WHERE id = ?", Caller: "post.go:88",
			})
		}
	}
	writeLog(t, p.root, entries)
	findings := checkQueryLog(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %+v", len(findings), findings)
	}
	if !strings.Contains(findings[0].Message, "in 2 requests") {
		t.Errorf("the finding does not say how widespread it is: %q", findings[0].Message)
	}
}

// No log is not a fault. A check that nags about a missing artefact teaches
// people to stop reading the output.
func TestNoLogIsSilent(t *testing.T) {
	p := logProject(t)
	if findings := checkQueryLog(p); len(findings) != 0 {
		t.Errorf("a project with no query log was reported: %q", findings[0].Message)
	}
}

// A log recorded before the code changed describes code that may no longer
// exist, and the findings read exactly the same either way. That is the whole
// danger of an artefact on disk, so it is said out loud.
func TestALogOlderThanTheCodeIsCalledOut(t *testing.T) {
	p := logProject(t)
	entries := repeated("GET /api/v1/posts", "SELECT 1", "post.go:88", 8)
	for i := range entries {
		entries[i].At = time.Now().Add(-2 * time.Hour)
	}
	writeLog(t, p.root, entries)
	// A source file written since.
	src := filepath.Join(p.apiRoot, "internal", "thing.go")
	now := time.Now()
	if err := os.Chtimes(src, now, now); err != nil {
		t.Fatal(err)
	}

	findings := checkQueryLog(p)
	if len(findings) != 2 {
		t.Fatalf("got %d findings, want the staleness note and the repeat: %+v", len(findings), findings)
	}
	if !strings.Contains(findings[0].Message, "newer") {
		t.Errorf("the first finding is not the staleness note: %q", findings[0].Message)
	}
	if !strings.Contains(findings[0].Fix, "GRIT_QUERY_LOG") {
		t.Errorf("the staleness note does not say how to record a new one: %q", findings[0].Fix)
	}
}

// A file that is not a query log is reported as that, rather than producing
// nothing and leaving somebody to wonder why the check is quiet.
func TestAFileThatIsNotALogIsReported(t *testing.T) {
	p := logProject(t)
	dir := filepath.Join(p.root, ".grit")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "queries.jsonl"), []byte("not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings := checkQueryLog(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1", len(findings))
	}
	if !strings.Contains(findings[0].Message, "not a query log") {
		t.Errorf("the finding does not say what is wrong: %q", findings[0].Message)
	}
}

// --queries points it somewhere else, which is how a CI job reads the log its
// test step produced.
func TestTheCallerCanNameTheLog(t *testing.T) {
	p := logProject(t)
	elsewhere := filepath.Join(t.TempDir(), "run.jsonl")
	raw, err := json.Marshal(queryEntry{
		At: time.Now(), Request: "r1", Route: "GET /x", SQL: "SELECT 1", Caller: "a.go:1",
	})
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for i := 0; i < 9; i++ {
		b.Write(raw)
		b.WriteByte('\n')
	}
	if err := os.WriteFile(elsewhere, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	p.queryLog = elsewhere
	findings := checkQueryLog(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings from the named log, want 1: %+v", len(findings), findings)
	}
	if !strings.Contains(findings[0].Message, "9 times") {
		t.Errorf("the finding does not describe the named log: %q", findings[0].Message)
	}
}
