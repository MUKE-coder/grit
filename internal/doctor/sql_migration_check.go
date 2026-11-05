package doctor

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// A checked-in SQL migration that cannot be rolled back, or cannot be loaded.
//
// internal/migrate refuses to load a half pair, which is the right behaviour
// and the wrong moment: it happens when somebody runs `grit migrate`, which on
// a deploy is after AutoMigrate has already changed the schema. Both of these
// are readable from the files alone, so they belong here, where they are found
// before the deploy rather than during it.
//
// The empty-down case is the one that matters in practice. An up file with
// statements and a down file with none rolls back silently and reads as a
// rollback that worked, which is the worst possible answer to "undo that".
// A pair that is empty on both sides is a freshly generated placeholder nobody
// has written yet, and is not reported.
func checkSQLMigrations(p *project) []Finding {
	dir := filepath.Join(p.apiRoot, "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		// No directory at all is a project from before SQL migrations existed.
		return nil
	}

	type pair struct{ up, down string }
	pairs := map[string]*pair{}
	ids := []string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		var half string
		switch {
		case strings.HasSuffix(name, ".up.sql"):
			half = "up"
		case strings.HasSuffix(name, ".down.sql"):
			half = "down"
		default:
			continue
		}
		id := strings.TrimSuffix(strings.TrimSuffix(name, ".sql"), "."+half)
		if pairs[id] == nil {
			pairs[id] = &pair{}
			ids = append(ids, id)
		}
		if half == "up" {
			pairs[id].up = name
		} else {
			pairs[id].down = name
		}
	}
	sort.Strings(ids)

	var out []Finding
	for _, id := range ids {
		pr := pairs[id]
		switch {
		case pr.up == "":
			out = append(out, Finding{
				Level:    "error",
				Resource: "migrations/" + id,
				Message: "has a .down.sql with no .up.sql beside it, so grit migrate " +
					"cannot load the directory at all and every pending migration stops",
				Fix: "add " + id + ".up.sql, or delete the down file",
			})
		case pr.down == "":
			out = append(out, Finding{
				Level:    "error",
				Resource: "migrations/" + id,
				Message: "has an .up.sql with no .down.sql beside it, so grit migrate " +
					"cannot load the directory at all and every pending migration stops",
				Fix: "add " + id + ".down.sql, even if it is only a comment explaining " +
					"why the change cannot be undone",
			})
		default:
			if hasStatements(filepath.Join(dir, pr.up)) && !hasStatements(filepath.Join(dir, pr.down)) {
				out = append(out, Finding{
					Level:    "warning",
					Resource: "migrations/" + id,
					Message: "changes the schema and its down file has nothing in it, so " +
						"grit migrate down --sql reports a rollback that did not happen",
					Fix: "write the reverse, or say in a comment why there is not one; " +
						"an empty down file cannot be told apart from one somebody forgot",
				})
			}
		}
	}
	return out
}

// hasStatements reports whether a .sql file has anything but comments and
// blank lines in it.
func hasStatements(path string) bool {
	body, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		return true
	}
	return false
}
