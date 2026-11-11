package generate

import (
	"os"
	"path/filepath"
	"strings"
)

// The grid was the one way into an owned table that ignored ownership.
//
// BulkEdit takes a list of ids and loaded exactly those rows, with no scope at
// all, so a caller holding the resource's edit permission could rewrite any
// other user's row by putting its id in the request. List refused, GetByID
// refused, Patch refused, and this one did it and answered "2 records
// updated". BulkCreate had the matching hole in the other direction: it
// decoded the row it was given, owner column included, so the grid could file
// rows under somebody else's account. Both shipped with the grid in v3.372.0.
//
// The generator emits the scope and the stamp now, which fixes a resource
// generated from here on and nothing that exists. This repair is the rest of
// it: a project with twenty owned resources is not going to regenerate twenty
// resources, and "regenerate everything" is not a fix anybody applies.
//
// Its own pass rather than a step inside RepairMoneyUpdates, which is gated on
// the money expander being absent: a project already repaired there would
// never have been visited again.

// RepairGridOwnerScoping scopes BulkEdit and stamps BulkCreate in every
// generated service under apiRoot that is owned and lacks them.
//
// It returns one line per service changed, and one per service that is owned,
// needs the fix and could not be given it, which is the case worth printing:
// that service is the vulnerable one.
func RepairGridOwnerScoping(apiRoot string) (changed []string, blocked []string, err error) {
	dir := filepath.Join(apiRoot, "internal", "services")
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		// No services directory. An upgrade runs against projects of every
		// shape, and a project with no generated resources has nothing here.
		return nil, nil, nil
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		raw, e := os.ReadFile(path)
		if e != nil {
			continue
		}
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")

		match := writableDecl.FindStringSubmatch(src)
		if match == nil {
			continue // not a generated resource service
		}
		pascal := match[1]

		next, note := repairGridOwnerSource(src, pascal, strings.ToLower(pascal))
		if note != "" {
			blocked = append(blocked, name+": "+note)
			continue
		}
		if next == src {
			continue
		}
		// Through writeRepairedGo, which formats and tells the manifest. A
		// repair that writes behind the manifest makes the file it fixed read
		// as hand-edited, so every later framework fix to it is skipped: the
		// repair poisons the file it repaired (review 2.5).
		if e := writeRepairedGo(path, next); e != nil {
			return changed, blocked, e
		}
		changed = append(changed, "internal/services/"+name+": the grid now scopes to the owner")
	}
	return changed, blocked, nil
}

// repairGridOwnerSource is RepairGridOwnerScoping for one service, as text.
//
// Returns the source unchanged with no note when there is nothing to do: a
// shared resource, or one already repaired. Returns a note only when the
// service is owned, still has the hole, and the anchors are not the ones Grit
// wrote, because that is a service somebody must look at by hand.
func repairGridOwnerSource(src, pascal, lower string) (string, string) {
	col, perm := ownerScopeOf(src)
	if col == "" {
		return src, "" // shared, or scoped in a way this cannot read: nothing to scope by
	}
	out := src

	// BulkEdit: the load of the rows the grid named.
	//
	// Both texts come from gridRowsLoad, the function the generator writes
	// them with, so a repaired service is byte-identical to a generated one.
	// Spelling either of them out here would be a second copy of the text and
	// the thing that drifts from it.
	plural := Pluralize(toSnakeCase(pascal))
	if strings.Contains(out, "func (s *"+pascal+"Service) BulkEdit(") &&
		!strings.Contains(out, "rows = authz.Scope") {
		anchor := gridRowsLoad(pascal, plural, "")
		if strings.Count(out, anchor) != 1 {
			return src, "BulkEdit loads its rows in a way Grit did not write, so it was left alone: " +
				"check by hand that a caller cannot edit another user's row by its id"
		}
		call := "authz.ScopeOwned(ctx, rows, \"" + col + "\")"
		if perm != "" {
			call = "authz.ScopeConfidential(ctx, rows, \"" + col + "\", \"" + perm + "\")"
		}
		scope := ownerGridScopeSource(ownerNameOf(col), lower, call)
		out = strings.Replace(out, anchor, gridRowsLoad(pascal, plural, scope), 1)
	}

	// BulkCreate: the owner of each row.
	field := pascalFromSnake(col)
	if strings.Contains(out, "func (s *"+pascal+"Service) BulkCreate(") &&
		!strings.Contains(out, "\t\titem."+field+" = authz.UserIDFrom(ctx)") {
		// Inside BulkCreate only: the three lines the generator writes to stop
		// the request setting what the database owns. Update has its own copy
		// of nothing like it, so this text appears once.
		anchor := "\t\titem.ID = \"\"\n\t\titem.CreatedAt = time.Time{}\n\t\titem.UpdatedAt = time.Time{}\n"
		if strings.Count(out, anchor) != 1 {
			return src, "BulkCreate does not reset its columns the way Grit wrote it, so it was left " +
				"alone: check by hand that the grid cannot file a row under another account"
		}
		out = strings.Replace(out, anchor, anchor+ownerGridStampSource(ownerNameOf(col), col, field), 1)
	}

	if out == src {
		return src, ""
	}
	// Both insertions call authz. Every owned service imports it already,
	// because its List is scoped, but a service that somehow does not would
	// not compile and that is worse than not repairing it.
	if !strings.Contains(out, "/internal/authz\"") {
		return src, "the service does not import internal/authz, so the scope could not be added"
	}
	return out, ""
}

// ownerNameOf is user_id -> user: the --owned-by field name, which is what the
// generator writes in the comment. The repair has only the column to work
// from, so the two agree exactly when the column is the field plus "_id",
// which is every shape --owned-by produces.
func ownerNameOf(column string) string {
	return strings.TrimSuffix(column, "_id")
}
