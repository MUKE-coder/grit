package generate

// The two files `grit generate migration` writes.
//
// Both open with a comment block, because a .sql file in a framework that
// otherwise has no migration files needs to say what reads it and when. The
// no-transaction marker in particular has to be discoverable from the file
// itself: it is the difference between a concurrent index build and an error
// Postgres gives you at the worst moment.

import (
	"fmt"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// NoTransactionMarkerText is the marker as the runner reads it. One spelling,
// in internal/scaffold, which owns the runner.
const NoTransactionMarkerText = scaffold.NoTransactionMarker

// migrationHeader is the comment block both halves share.
func migrationHeader(version, slug, half string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "-- %s_%s.%s.sql\n", version, slug, half)
	b.WriteString("--\n")
	if half == "up" {
		b.WriteString("-- Applied once, after AutoMigrate, in version order, by `grit migrate`.\n")
		b.WriteString("--\n")
		b.WriteString("-- AutoMigrate adds tables, columns, indexes and constraints and takes nothing\n")
		b.WriteString("-- away, so this file is for what it cannot do: a rename, a type change, a\n")
		b.WriteString("-- backfill, a drop, or an index on a table too big to lock.\n")
	} else {
		b.WriteString("-- The reverse of the up file, run by `grit migrate down --sql`.\n")
		b.WriteString("--\n")
		b.WriteString("-- Leave it empty only when the change genuinely cannot be undone, and say so\n")
		b.WriteString("-- here in a comment: an empty down file rolls back silently and reads as a\n")
		b.WriteString("-- rollback that worked.\n")
	}
	b.WriteString("--\n")
	b.WriteString("-- Statements run in one transaction unless the file's first line is\n")
	b.WriteString("--     " + NoTransactionMarkerText + "\n")
	b.WriteString("-- which is required for CREATE INDEX CONCURRENTLY, and is how you ask for a\n")
	b.WriteString("-- statement Postgres refuses to run inside one.\n")
	return b.String()
}

// migrationUpSQL is the up file: an index build when --index was given, and an
// empty file with instructions otherwise.
func migrationUpSQL(version, slug string, index *IndexSpec) string {
	if index == nil {
		return migrationHeader(version, slug, "up") + `
-- Write the change here. For example:
--
--     ALTER TABLE companies RENAME TO organisations;
--     ALTER TABLE orders ALTER COLUMN total TYPE bigint USING (total * 100)::bigint;
--     UPDATE orders SET total_cents = total * 100 WHERE total_cents IS NULL;
`
	}

	cols := strings.Join(index.Columns, ", ")
	unique := ""
	if index.Unique {
		unique = "UNIQUE "
	}

	var b strings.Builder
	if index.Driver == "postgres" {
		// The marker goes first, because that is where the runner looks for it.
		b.WriteString(NoTransactionMarkerText + "\n")
		b.WriteString(migrationHeader(version, slug, "up"))
		b.WriteString("--\n")
		b.WriteString("-- CONCURRENTLY takes no write lock, so the table stays readable AND writable\n")
		b.WriteString("-- while the index builds. That is the whole point: a plain CREATE INDEX on a\n")
		b.WriteString("-- large busy table is an outage.\n")
		b.WriteString("--\n")
		b.WriteString("-- It is slower, it cannot run inside a transaction (hence the marker on the\n")
		b.WriteString("-- first line), and if it fails it leaves an INVALID index behind. Check with\n")
		b.WriteString("--     SELECT indexrelid::regclass FROM pg_index WHERE NOT indisvalid;\n")
		b.WriteString("-- then drop it and run this again.\n")
		fmt.Fprintf(&b, "CREATE %sINDEX CONCURRENTLY IF NOT EXISTS %s ON %s (%s);\n",
			unique, index.Name, index.Table, cols)
		return b.String()
	}

	b.WriteString(migrationHeader(version, slug, "up"))
	b.WriteString("--\n")
	switch index.Driver {
	case "mysql":
		b.WriteString("-- MySQL has no CONCURRENTLY. InnoDB builds a secondary index online by\n")
		b.WriteString("-- default (ALGORITHM=INPLACE, LOCK=NONE), so writes continue, but confirm it\n")
		b.WriteString("-- for your version before running this on a large table.\n")
	default:
		b.WriteString("-- SQLite has no CONCURRENTLY and locks the database for the build. On a small\n")
		b.WriteString("-- file that is instant; it is still worth knowing before production.\n")
	}
	fmt.Fprintf(&b, "CREATE %sINDEX IF NOT EXISTS %s ON %s (%s);\n",
		unique, index.Name, index.Table, cols)
	return b.String()
}

// migrationDownSQL is the reverse.
func migrationDownSQL(version, slug string, index *IndexSpec) string {
	if index == nil {
		return migrationHeader(version, slug, "down") + `
-- Write the reverse here. For example:
--
--     ALTER TABLE organisations RENAME TO companies;
`
	}

	var b strings.Builder
	if index.Driver == "postgres" {
		b.WriteString(NoTransactionMarkerText + "\n")
		b.WriteString(migrationHeader(version, slug, "down"))
		b.WriteString("--\n")
		b.WriteString("-- DROP INDEX CONCURRENTLY also refuses to run in a transaction.\n")
		fmt.Fprintf(&b, "DROP INDEX CONCURRENTLY IF EXISTS %s;\n", index.Name)
		return b.String()
	}
	b.WriteString(migrationHeader(version, slug, "down"))
	if index.Driver == "mysql" {
		// MySQL's DROP INDEX needs the table and has no IF EXISTS.
		fmt.Fprintf(&b, "DROP INDEX %s ON %s;\n", index.Name, index.Table)
		return b.String()
	}
	fmt.Fprintf(&b, "DROP INDEX IF EXISTS %s;\n", index.Name)
	return b.String()
}
