package scaffold

// Checked-in SQL migrations, for the changes AutoMigrate cannot make.
//
// grit#162. AutoMigrate is additive and that is the right default: it adds
// tables, columns, indexes and constraints, and the computed history in
// history.go can undo any of it because none of it is destructive. The five
// things it does not do are a rename, a type change, a backfill, a drop, and an
// index on a table too large to lock, and a project that needed one of those
// had to leave the framework's conventions for golang-migrate or goose.
//
// So: a directory of pairs, applied after AutoMigrate, recorded by version.
// Two histories because there are two mechanisms, and pretending otherwise
// would mean a `grit migrate down` that has to interleave a computed rollback
// with a hand-written one.

// NoTransactionMarker is the line that opts a file out of its transaction.
//
// Exported and declared here because two pieces of code have to agree about
// it: the runner this file emits, which looks for it on a file's first line,
// and `grit generate migration --index` in internal/generate, which writes it
// there. internal/generate imports this package, so there is one spelling of
// it rather than two that drift.
const NoTransactionMarker = "-- grit:no-transaction"

// apiMigrationsPackageGo emits migrations/migrations.go.
func apiMigrationsPackageGo() string {
	return `// Package migrations holds the SQL that AutoMigrate cannot write.
//
// Create a pair with ` + "`" + `grit generate migration <name>` + "`" + `. Both halves are named
// <version>_<name>.up.sql and <version>_<name>.down.sql, and ` + "`" + `grit migrate` + "`" + `
// applies the pending ones in version order after AutoMigrate.
//
// See README.md in this directory for the file format.
package migrations

import "embed"

// FS is every file in this directory, compiled into the binary.
//
// The pattern is * rather than *.sql because an embed pattern that matches
// nothing is a build error, and a project with no migrations yet has no .sql
// files at all. internal/migrate filters for the extension.
//
// Embedded rather than read from disk because the migrate binary runs in three
// places with three working directories: ` + "`" + `grit migrate` + "`" + ` runs it with go run
// from apps/api, the production compose file runs ./migrate inside a container
// where only the binary was copied, and somebody may run the built binary from
// anywhere at all. A relative path that works in one of those is a migration
// silently skipped in another, which is the worst way for a schema change to
// fail. Embedding also pins the checksum to the build: the SQL that ran is the
// SQL that was compiled in.
//
//go:embed *
var FS embed.FS
`
}

// apiMigrationsReadme emits migrations/README.md.
//
// A real README rather than a .keep, because this is the one directory in a
// generated project whose contents are hand-written SQL, and the two rules
// that are easy to get wrong (the version order and the transaction marker)
// have to be readable from inside it.
func apiMigrationsReadme() string {
	return "# migrations\n" +
		"\n" +
		"SQL that `AutoMigrate` cannot write.\n" +
		"\n" +
		"`grit migrate` runs `AutoMigrate` first, then every pending pair in this\n" +
		"directory in version order, and records each one in `grit_sql_migrations`.\n" +
		"\n" +
		"```bash\n" +
		"grit generate migration rename_company_to_organisation\n" +
		"grit generate migration add_order_customer_index --index orders:customer_id\n" +
		"\n" +
		"grit migrate                 # AutoMigrate, then the pending pairs\n" +
		"grit migrate status          # which pairs have run\n" +
		"grit migrate down --sql      # undo the most recent pair\n" +
		"```\n" +
		"\n" +
		"## The file format\n" +
		"\n" +
		"```\n" +
		"20261105120000_rename_company_to_organisation.up.sql\n" +
		"20261105120000_rename_company_to_organisation.down.sql\n" +
		"```\n" +
		"\n" +
		"The version is a UTC timestamp and decides the order. The down file is the\n" +
		"reverse of the up file; leave it empty only when the change genuinely cannot\n" +
		"be undone, and say so in a comment, because an empty down file rolls back\n" +
		"silently and reads as a rollback that worked.\n" +
		"\n" +
		"A file's statements run in one transaction. To opt out, make the first line:\n" +
		"\n" +
		"```sql\n" +
		"-- grit:no-transaction\n" +
		"```\n" +
		"\n" +
		"`CREATE INDEX CONCURRENTLY` needs this: Postgres refuses it inside a\n" +
		"transaction. `grit generate migration --index` writes the marker for you.\n" +
		"\n" +
		"## What not to do\n" +
		"\n" +
		"**Do not edit a pair that has already run.** The database has the old version\n" +
		"and the file has the new one, and nothing can reconcile them. `grit migrate`\n" +
		"compares a checksum and says so when it happens; write a new pair instead.\n" +
		"\n" +
		"**Do not do in SQL what AutoMigrate already does.** Adding a column belongs in\n" +
		"the model. This directory is for what the model cannot express.\n" +
		"\n" +
		"**Do not assume a transaction rolls back DDL.** Postgres rolls back a failed\n" +
		"`ALTER TABLE`. MySQL commits implicitly at every DDL statement, so half an\n" +
		"applied file stays applied there.\n"
}

// apiMigrateSQLGo emits internal/migrate/sql.go.
func apiMigrateSQLGo() string {
	return `package migrate

// Checked-in SQL migrations: the changes AutoMigrate cannot make.
//
// AutoMigrate is additive, so history.go can compute the reverse of a run and
// no down script has to be written. A rename, a type change, a backfill, a drop
// or a concurrent index build is none of those things, and each needs an
// explicit pair: a file that does it and a file that undoes it.
//
// They are recorded separately from AutoMigrate's runs, in grit_sql_migrations,
// because they are a different mechanism. One combined history would make
// "undo the last thing" interleave a computed rollback with a hand-written one,
// and the order it chose would be a guess.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// NoTransactionMarker on a file's first line runs that file without a
// transaction.
//
// CREATE INDEX CONCURRENTLY is the reason it exists: Postgres refuses it inside
// one, and the refusal arrives when you run the migration rather than when you
// write it.
const NoTransactionMarker = "` + NoTransactionMarker + `"

// SQLMigration is one pair, as loaded from the embedded directory.
type SQLMigration struct {
	// Version is the file name's leading timestamp, and the sort key.
	Version string
	// Name is the rest of it: rename_company_to_organisation.
	Name string
	Up   string
	Down string
	// UpNoTransaction and DownNoTransaction are set by the marker.
	UpNoTransaction   bool
	DownNoTransaction bool
	// Checksum is of Up. A pair that changed after being applied is a database
	// and a file that disagree, and the only honest thing to do is say so.
	Checksum string
}

// ID is how a migration is named in output.
func (m SQLMigration) ID() string { return m.Version + "_" + m.Name }

// AppliedSQLMigration is one row of grit_sql_migrations.
//
// Keyed by the whole id rather than by the version. Two migrations generated in
// the same second share a version, and with the version as the key the second
// one ran and then could not record itself: "UNIQUE constraint failed". Found
// by generating two in one command. The generator avoids the collision as well,
// by taking the next free second, but a key that cannot collide is the half
// that holds when two people generate on two branches.
type AppliedSQLMigration struct {
	ID        string ` + "`" + `gorm:"primaryKey;size:224"` + "`" + `
	Version   string ` + "`" + `gorm:"size:32;not null;index"` + "`" + `
	Name      string ` + "`" + `gorm:"size:191;not null"` + "`" + `
	Checksum  string ` + "`" + `gorm:"size:64;not null"` + "`" + `
	AppliedAt time.Time
}

// TableName pins the table, so renaming the type cannot orphan the history.
func (AppliedSQLMigration) TableName() string { return "grit_sql_migrations" }

// EnsureSQLHistory creates the table the pairs are recorded in.
func EnsureSQLHistory(db *gorm.DB) error {
	if err := db.AutoMigrate(&AppliedSQLMigration{}); err != nil {
		return fmt.Errorf("preparing the SQL migration history: %w", err)
	}
	return nil
}

// LoadSQLMigrations reads every pair out of fsys, in version order.
//
// A .up.sql with no .down.sql is an error rather than a warning: it is a
// migration that cannot be rolled back, discovered at the moment you need to
// roll it back.
func LoadSQLMigrations(fsys fs.FS) ([]SQLMigration, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("reading the migrations directory: %w", err)
	}

	byID := map[string]*SQLMigration{}
	order := []string{}
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
		version, rest, ok := strings.Cut(id, "_")
		if !ok || version == "" || rest == "" {
			return nil, fmt.Errorf("migrations/%s is not named <version>_<name>.%s.sql", name, half)
		}
		body, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("reading migrations/%s: %w", name, err)
		}
		m, seen := byID[id]
		if !seen {
			m = &SQLMigration{Version: version, Name: rest}
			byID[id] = m
			order = append(order, id)
		}
		if half == "up" {
			m.Up = string(body)
			m.UpNoTransaction = hasNoTransactionMarker(m.Up)
			sum := sha256.Sum256(body)
			m.Checksum = hex.EncodeToString(sum[:])
		} else {
			m.Down = string(body)
			m.DownNoTransaction = hasNoTransactionMarker(m.Down)
		}
	}

	sort.Strings(order)
	out := make([]SQLMigration, 0, len(order))
	for _, id := range order {
		m := byID[id]
		if m.Up == "" {
			return nil, fmt.Errorf("migrations/%s.down.sql has no .up.sql beside it", id)
		}
		if m.Down == "" {
			return nil, fmt.Errorf("migrations/%s.up.sql has no .down.sql beside it; write one, or write a comment in it saying why the change cannot be undone", id)
		}
		out = append(out, *m)
	}
	return out, nil
}

// hasNoTransactionMarker looks only at the first non-empty line. Anywhere else
// it is a comment somebody wrote about transactions, not a request.
func hasNoTransactionMarker(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return strings.EqualFold(line, NoTransactionMarker)
	}
	return false
}

// PendingSQL is the migrations that have not been applied, in version order.
func PendingSQL(db *gorm.DB, all []SQLMigration) ([]SQLMigration, error) {
	applied, err := appliedSQLByID(db)
	if err != nil {
		return nil, err
	}
	var out []SQLMigration
	for _, m := range all {
		if _, done := applied[m.ID()]; !done {
			out = append(out, m)
		}
	}
	return out, nil
}

// ChangedSQL is the applied migrations whose up file no longer matches what
// ran.
//
// Editing an applied migration is the classic mistake, and it is silent: the
// database holds the old version, the file holds the new one, and every
// developer who has not run it yet gets a different schema from everybody who
// has. Reported rather than fatal, because the fix is a new migration and
// refusing to migrate at all would block that.
func ChangedSQL(db *gorm.DB, all []SQLMigration) ([]string, error) {
	applied, err := appliedSQLByID(db)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, m := range all {
		if row, done := applied[m.ID()]; done && row.Checksum != m.Checksum {
			out = append(out, m.ID())
		}
	}
	return out, nil
}

func appliedSQLByID(db *gorm.DB) (map[string]AppliedSQLMigration, error) {
	rows, err := AppliedSQL(db)
	if err != nil {
		return nil, err
	}
	out := make(map[string]AppliedSQLMigration, len(rows))
	for _, r := range rows {
		out[r.ID] = r
	}
	return out, nil
}

// AppliedSQL is the history, oldest first.
//
// Ordered by id, which is version then name, so it matches the order
// LoadSQLMigrations returns and a rollback takes them off the end.
func AppliedSQL(db *gorm.DB) ([]AppliedSQLMigration, error) {
	var rows []AppliedSQLMigration
	if err := db.Order("id asc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("reading the SQL migration history: %w", err)
	}
	return rows, nil
}

// ApplySQL runs every pending migration in order and records each one.
//
// It stops at the first failure and returns what it managed, so the output
// names the migration that broke rather than leaving you to work it out from a
// partial schema.
func ApplySQL(db *gorm.DB, all []SQLMigration) ([]string, error) {
	pending, err := PendingSQL(db, all)
	if err != nil {
		return nil, err
	}
	var done []string
	for _, m := range pending {
		if err := runSQLBody(db, m.Up, m.UpNoTransaction); err != nil {
			return done, fmt.Errorf("migrations/%s.up.sql: %w", m.ID(), err)
		}
		row := AppliedSQLMigration{
			ID:        m.ID(),
			Version:   m.Version,
			Name:      m.Name,
			Checksum:  m.Checksum,
			AppliedAt: time.Now(),
		}
		if err := db.Create(&row).Error; err != nil {
			return done, fmt.Errorf("migrations/%s ran but could not be recorded, so it would run again: %w", m.ID(), err)
		}
		done = append(done, m.ID())
	}
	return done, nil
}

// RollbackSQL runs the down half of the most recently applied migrations.
//
// steps counts migrations, newest first. A migration whose file is no longer in
// the binary cannot be rolled back, and saying so is better than skipping it.
func RollbackSQL(db *gorm.DB, all []SQLMigration, steps int) ([]string, error) {
	if steps <= 0 {
		steps = 1
	}
	rows, err := AppliedSQL(db)
	if err != nil {
		return nil, err
	}
	byID := map[string]SQLMigration{}
	for _, m := range all {
		byID[m.ID()] = m
	}

	var undone []string
	for i := len(rows) - 1; i >= 0 && len(undone) < steps; i-- {
		row := rows[i]
		m, ok := byID[row.ID]
		if !ok {
			return undone, fmt.Errorf("%s was applied but its files are not in this build, so there is nothing to run; check out the commit that added it and try again", row.ID)
		}
		if err := runSQLBody(db, m.Down, m.DownNoTransaction); err != nil {
			return undone, fmt.Errorf("migrations/%s.down.sql: %w", m.ID(), err)
		}
		if err := db.Delete(&AppliedSQLMigration{}, "id = ?", row.ID).Error; err != nil {
			return undone, fmt.Errorf("migrations/%s was rolled back but the history still records it: %w", m.ID(), err)
		}
		undone = append(undone, m.ID())
	}
	return undone, nil
}

// runSQLBody executes a file, in a transaction unless it asked not to be.
//
// Statement by statement, because the drivers require it: pgx refuses more than
// one command in a single Exec, and the MySQL driver needs multiStatements on
// the DSN. Splitting here means the SQL in the file is ordinary SQL.
func runSQLBody(db *gorm.DB, body string, noTransaction bool) error {
	statements := SplitSQL(body)
	if len(statements) == 0 {
		return nil
	}
	exec := func(tx *gorm.DB) error {
		for _, s := range statements {
			if err := tx.Exec(s).Error; err != nil {
				return fmt.Errorf("%s: %w", firstLine(s), err)
			}
		}
		return nil
	}
	if noTransaction {
		return exec(db)
	}
	return db.Transaction(exec)
}

// firstLine is how a failing statement is named, so the error does not carry a
// whole backfill.
func firstLine(s string) string {
	line := s
	if i := strings.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	line = strings.TrimSpace(line)
	if len(line) > 80 {
		line = line[:77] + "..."
	}
	return line
}

// SplitSQL splits a file into statements on semicolons that are not inside
// something.
//
// Quoted strings, quoted identifiers, dollar-quoted bodies, line comments and
// block comments all legitimately contain semicolons, and a split on every
// semicolon turns a working function definition into three syntax errors. The
// dollar-quote case is not hypothetical: it is how every PL/pgSQL function and
// trigger in Postgres is written.
func SplitSQL(body string) []string {
	var out []string
	var current strings.Builder

	flush := func() {
		if s := strings.TrimSpace(current.String()); s != "" && !onlyComments(s) {
			out = append(out, s)
		}
		current.Reset()
	}

	runes := []rune(body)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case c == '-' && i+1 < len(runes) && runes[i+1] == '-':
			// Line comment, kept so the statement still reads like the file.
			for i < len(runes) && runes[i] != '\n' {
				current.WriteRune(runes[i])
				i++
			}
			if i < len(runes) {
				current.WriteRune('\n')
			}
		case c == '/' && i+1 < len(runes) && runes[i+1] == '*':
			for i < len(runes) {
				current.WriteRune(runes[i])
				if runes[i] == '/' && i > 0 && runes[i-1] == '*' && i >= 2 {
					break
				}
				i++
			}
		case c == '\'' || c == '"' || c == '` + "`" + `':
			quote := c
			current.WriteRune(c)
			i++
			for i < len(runes) {
				current.WriteRune(runes[i])
				if runes[i] == quote {
					// Two in a row is an escaped quote, not the end.
					if i+1 < len(runes) && runes[i+1] == quote {
						i++
						current.WriteRune(runes[i])
						i++
						continue
					}
					break
				}
				if runes[i] == '\\' && i+1 < len(runes) {
					i++
					current.WriteRune(runes[i])
				}
				i++
			}
		case c == '$':
			if tag, width := dollarTag(runes, i); width > 0 {
				end := indexRunes(runes, tag, i+width)
				if end < 0 {
					// Unterminated: the database will say so more usefully
					// than this function can.
					current.WriteString(string(runes[i:]))
					i = len(runes)
					break
				}
				current.WriteString(string(runes[i : end+len(tag)]))
				i = end + len(tag) - 1
			} else {
				current.WriteRune(c)
			}
		case c == ';':
			flush()
		default:
			current.WriteRune(c)
		}
	}
	flush()
	return out
}

// dollarTag reads a Postgres dollar quote opener at i: $$ or $tag$.
func dollarTag(runes []rune, i int) (string, int) {
	if runes[i] != '$' {
		return "", 0
	}
	j := i + 1
	for j < len(runes) && (runes[j] == '_' || isTagRune(runes[j])) {
		j++
	}
	if j < len(runes) && runes[j] == '$' {
		return string(runes[i : j+1]), j + 1 - i
	}
	return "", 0
}

func isTagRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func indexRunes(runes []rune, needle string, from int) int {
	n := []rune(needle)
	for i := from; i+len(n) <= len(runes); i++ {
		match := true
		for j := range n {
			if runes[i+j] != n[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// onlyComments reports whether a fragment has nothing to execute, which is what
// the generated placeholder files are until somebody writes in them.
func onlyComments(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		return false
	}
	return true
}
`
}
