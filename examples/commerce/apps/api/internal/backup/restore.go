package backup

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"commerce/apps/api/internal/appendonly"
	"commerce/apps/api/internal/erasure"
	"commerce/apps/api/internal/models"
)

// insertTable returns the table an INSERT statement writes to, or "" for any
// other statement.
func insertTable(stmt string) string {
	s := strings.TrimSpace(stmt)
	const prefix = "INSERT INTO \""
	if !strings.HasPrefix(s, prefix) {
		return ""
	}
	rest := s[len(prefix):]
	if i := strings.Index(rest, "\""); i > 0 {
		return rest[:i]
	}
	return ""
}

// replayOrder puts a dump's INSERTs into an order the foreign keys accept.
//
// The dump writes tables in models.Models() order, and the resource generator
// registers an --items child before its parent, so a straight replay inserted
// journal_lines before the journal_entries they point at and failed on the
// foreign key. Every project with line items had backups it could not restore.
// Reordering here rather than only in the backup rescues archives that were
// already written.
func replayOrder(db *gorm.DB, stmts []string) []string {
	groups := map[string][]string{}
	var tables, others []string
	for _, s := range stmts {
		t := insertTable(s)
		if t == "" {
			others = append(others, s)
			continue
		}
		if _, seen := groups[t]; !seen {
			tables = append(tables, t)
		}
		groups[t] = append(groups[t], s)
	}
	out := append([]string{}, others...)
	for _, t := range dependencyOrder(db, tables) {
		out = append(out, groups[t]...)
	}
	return out
}

// dependencyOrder sorts tables so each comes after the tables its foreign keys
// point at, read from the models' own relationships. Unrelated tables keep the
// order they had, and a cycle, which a foreign key cannot form without a
// nullable column, falls back to that order for the tables caught in it.
func dependencyOrder(db *gorm.DB, tables []string) []string {
	needs := map[string]map[string]bool{}
	link := func(child, parent string) {
		if child == "" || parent == "" || child == parent {
			return
		}
		if needs[child] == nil {
			needs[child] = map[string]bool{}
		}
		needs[child][parent] = true
	}
	for _, m := range models.Models() {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(m); err != nil {
			continue
		}
		s := stmt.Schema
		for _, rel := range s.Relationships.Relations {
			if rel.FieldSchema == nil {
				continue
			}
			switch rel.Type {
			case schema.BelongsTo:
				link(s.Table, rel.FieldSchema.Table)
			case schema.HasOne, schema.HasMany:
				link(rel.FieldSchema.Table, s.Table)
			case schema.Many2Many:
				if rel.JoinTable != nil {
					link(rel.JoinTable.Table, s.Table)
					link(rel.JoinTable.Table, rel.FieldSchema.Table)
				}
			}
		}
	}

	present := map[string]bool{}
	for _, t := range tables {
		present[t] = true
	}
	placed := map[string]bool{}
	var out []string
	for len(out) < len(tables) {
		progressed := false
		for _, t := range tables {
			if placed[t] {
				continue
			}
			ready := true
			for p := range needs[t] {
				if present[p] && !placed[p] {
					ready = false
					break
				}
			}
			if ready {
				out = append(out, t)
				placed[t] = true
				progressed = true
			}
		}
		if !progressed {
			for _, t := range tables {
				if !placed[t] {
					out = append(out, t)
					placed[t] = true
				}
			}
		}
	}
	return out
}

// reapplyErasures puts back the deletion-journal entries made after the backup
// was taken, and erases those people again.
//
// The journal is only ever appended to, so the archive's copy is a prefix of
// the one that was just replaced: the missing entries go back verbatim, hashes
// and all, and the chain verifies exactly as it did. Each person is then
// scrubbed again, because the archive holds them as they were before they asked
// to be forgotten.
func reapplyErasures(tx *gorm.DB, before []models.DeletionJournal) (int, error) {
	if len(before) == 0 {
		return 0, nil
	}
	var restored []models.DeletionJournal
	if err := tx.Select("hash").Find(&restored).Error; err != nil {
		return 0, fmt.Errorf("reading the restored journal: %w", err)
	}
	have := make(map[string]bool, len(restored))
	for _, j := range restored {
		have[j.Hash] = true
	}
	n := 0
	for _, j := range before {
		if have[j.Hash] {
			continue
		}
		entry := j
		if err := tx.Create(&entry).Error; err != nil {
			return n, fmt.Errorf("restoring journal entry %s: %w", j.ID, err)
		}
		if _, _, err := erasure.Scrub(tx, j.DeletedUserID); err != nil {
			return n, fmt.Errorf("erasing %s again: %w", j.DeletedUserID, err)
		}
		n++
	}
	return n, nil
}

// SplitStatements splits our generated dump.sql into executable statements.
//
// This is NOT a general SQL parser — it doesn't need to be. We generate the file
// ourselves and only ever emit numbers, NULL/TRUE/FALSE, and single-quoted
// literals with ” escaping. Tracking quote state is therefore exact. Splitting
// naively on ";" would corrupt any value containing a semicolon.
func SplitStatements(script string) []string {
	var out []string
	var cur strings.Builder
	inString := false

	rs := []rune(script)
	for i := 0; i < len(rs); i++ {
		c := rs[i]
		if inString {
			cur.WriteRune(c)
			if c == '\'' {
				// '' is an escaped quote, not the end of the literal.
				if i+1 < len(rs) && rs[i+1] == '\'' {
					cur.WriteRune('\'')
					i++
					continue
				}
				inString = false
			}
			continue
		}
		switch c {
		case '\'':
			inString = true
			cur.WriteRune(c)
		case ';':
			if s := strings.TrimSpace(cur.String()); s != "" {
				out = append(out, s)
			}
			cur.Reset()
		default:
			cur.WriteRune(c)
		}
	}
	if s := strings.TrimSpace(cur.String()); s != "" {
		out = append(out, s)
	}
	return out
}

// stripComments removes SQL "--" comments, but ONLY when they're real comments
// and not part of a quoted string value. The naive line-based version dropped
// any line beginning with "--", which corrupted multi-line string values whose
// continuation line happened to start with "--" (e.g. a note field). This walks
// the script tracking string state (with ” escape handling), so text inside a
// literal is never touched.
func stripComments(script string) string {
	var b strings.Builder
	inString := false
	for i := 0; i < len(script); i++ {
		c := script[i]
		if inString {
			b.WriteByte(c)
			if c == '\'' {
				// A doubled '' is an escaped quote — stays inside the string.
				if i+1 < len(script) && script[i+1] == '\'' {
					b.WriteByte(script[i+1])
					i++
					continue
				}
				inString = false
			}
			continue
		}
		if c == '\'' {
			inString = true
			b.WriteByte(c)
			continue
		}
		// Outside a string, "--" begins a comment that runs to end of line.
		if c == '-' && i+1 < len(script) && script[i+1] == '-' {
			for i < len(script) && script[i] != '\n' {
				i++
			}
			if i < len(script) {
				b.WriteByte('\n')
			}
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

// Restore replays a backup archive into the connected database inside a single
// transaction: either every row lands or nothing does.
//
// The archive carries DATA, not schema — run migrations on the target database
// first (cmd/restore does this for you).
func Restore(db *gorm.DB, zipPath string) (Manifest, error) {
	var man Manifest

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return man, fmt.Errorf("opening %s: %w", zipPath, err)
	}
	defer zr.Close()

	var dump string
	for _, f := range zr.File {
		switch f.Name {
		case "dump.sql", "metadata.json":
			rc, err := f.Open()
			if err != nil {
				return man, err
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return man, err
			}
			if f.Name == "dump.sql" {
				dump = string(data)
			} else {
				_ = json.Unmarshal(data, &man)
			}
		}
	}
	if dump == "" {
		return man, errors.New("dump.sql not found in archive")
	}

	stmts := SplitStatements(stripComments(dump))
	return man, db.Transaction(func(tx *gorm.DB) error {
		// Clear every table the dump repopulates BEFORE replaying it. Migrations
		// run first (cmd/restore) and seed baseline rows — the default roles —
		// and the dump carries its own authoritative copy of those same rows.
		// Without this, replaying the dump's ADMIN/EDITOR/USER inserts collides
		// with the freshly seeded ones on the unique role-name index and the
		// whole restore aborts: the backup becomes unrestorable. Truncating first
		// makes the restored database match the backup exactly, and makes restore
		// idempotent onto a non-empty schema. Table names come from the manifest,
		// derived from models.Models() — never user input — so this is not an
		// injection surface. RESTART IDENTITY resets sequences; CASCADE handles
		// foreign keys regardless of order.
		// The deletion journal as it stands, before the archive replaces it. An
		// erasure made after the backup was taken is here and not in the archive,
		// and replaying the archive alone brought the person back and lost the
		// proof they had ever asked to be forgotten.
		var journal []models.DeletionJournal
		if tx.Migrator().HasTable(&models.DeletionJournal{}) {
			if err := tx.Order("created_at asc, id asc").Find(&journal).Error; err != nil {
				return fmt.Errorf("reading the deletion journal: %w", err)
			}
		}

		// Append-only tables refuse TRUNCATE by trigger, which is right for every
		// writer except this one: a restore rewrites them wholesale. Suspended for
		// this transaction only, and switched back on before it commits.
		resume, err := appendonly.Suspend(tx)
		if err != nil {
			return err
		}
		if len(man.Tables) > 0 {
			quoted := make([]string, len(man.Tables))
			for i, t := range man.Tables {
				quoted[i] = "\"" + t + "\""
			}
			if err := tx.Exec("TRUNCATE " + strings.Join(quoted, ", ") + " RESTART IDENTITY CASCADE").Error; err != nil {
				return fmt.Errorf("clearing tables before restore: %w", err)
			}
		}

		for _, s := range replayOrder(tx, stmts) {
			switch strings.ToUpper(strings.TrimSpace(s)) {
			case "BEGIN", "COMMIT":
				continue // we own the transaction
			}
			if err := tx.Exec(s).Error; err != nil {
				head := s
				if len(head) > 120 {
					head = head[:120] + "..."
				}
				return fmt.Errorf("executing %q: %w", head, err)
			}
		}
		reapplied, err := reapplyErasures(tx, journal)
		if err != nil {
			return err
		}
		if reapplied > 0 {
			log.Printf("Re-applied %d erasure(s) made after this backup was taken", reapplied)
		}
		return resume()
	})
}
