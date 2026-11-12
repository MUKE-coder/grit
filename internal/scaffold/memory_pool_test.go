package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// A generated test that opens its own ":memory:" database and then runs a
// goroutine has to cap the connection pool at one.
//
// Every connection to ":memory:" gets its OWN empty database. The pool opens a
// second as soon as two queries overlap, so a concurrent test lands on a fresh
// empty database and reports "no such table". It passes or fails on timing:
// TestTwoExchangesRaceAndOnlyOneWins went green on this machine and red in CI
// on the same commit, and the error it printed was about the feature rather
// than about the harness.
//
// internal/testkit caps its pool and says why. A helper that opens its own
// database is outside that, and internal/services cannot use testkit anyway
// because testkit imports services.
//
// Only the concurrent ones. Fifty-odd helpers open a bare ":memory:" and all
// but one are sequential, where a second connection is never asked for. A
// sweep over the rest would be fifty files touched to prevent nothing, and
// this test is what catches the first of them to grow a goroutine.
var (
	bareMemoryOpen  = regexp.MustCompile(`sqlite\.Open\("` + ":memory:" + `"\)`)
	poolCapped      = regexp.MustCompile(`SetMaxOpenConns\(1\)`)
	spawnsGoroutine = regexp.MustCompile(`\bgo func\(|sync\.WaitGroup`)
)

func TestAConcurrentGeneratedTestCapsItsConnectionPool(t *testing.T) {
	var offenders []string

	for _, path := range goAndTemplateFiles(t, ".") {
		// testkit is the one that documents the cap, and it has it.
		if strings.Contains(filepath.ToSlash(path), "templates/api/testkit") {
			continue
		}
		// This file names the pattern in its own prose.
		if strings.HasSuffix(path, "memory_pool_test.go") {
			continue
		}
		body := readFileOrEmpty(path)
		if !bareMemoryOpen.MatchString(body) {
			continue
		}
		if !spawnsGoroutine.MatchString(body) {
			continue // sequential: a second connection is never asked for
		}
		if poolCapped.MatchString(body) {
			continue
		}
		rel, err := filepath.Rel(".", path)
		if err != nil {
			rel = path
		}
		offenders = append(offenders, filepath.ToSlash(rel))
	}
	sort.Strings(offenders)

	if len(offenders) > 0 {
		t.Errorf("%d generated test helper(s) open a bare \":memory:\" database and run a "+
			"goroutine without capping the pool at one connection:\n  %s\n\n"+
			"Each overlapping query gets its own empty database, so the test passes or fails "+
			"on timing and the error it prints is about the feature. Add, after gorm.Open:\n\n"+
			"    sqlDB, err := db.DB()\n    if err != nil { t.Fatalf(...) }\n"+
			"    sqlDB.SetMaxOpenConns(1)",
			len(offenders), strings.Join(offenders, "\n  "))
	}
}

// goAndTemplateFiles is every .go and .tmpl file under root, which is where a
// generated test helper can live: some are Go string literals in this package
// and some are template files.
func goAndTemplateFiles(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	for _, dir := range []string{root, filepath.Join("..", "generate"), filepath.Join("..", "plugin")} {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // a directory that is not there is not a failure
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".tmpl") {
				out = append(out, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walking %s: %v", dir, err)
		}
	}
	return out
}
