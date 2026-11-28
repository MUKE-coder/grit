package doctor

// An owned resource whose erasure takes the row and leaves the file.
//
// A resource generated with --owned-by registers itself for right-to-erasure,
// so erasing a user deletes the rows they own here. That is a DELETE, and a
// file column is not in the database: the row names an object in storage.
//
// Erasing the user already took their own uploads and every object those owned
// (grit#202, grit#221). What that cannot see is a file somebody else uploaded
// onto one of these rows, an administrator above all, which is exactly what
// the owner check allows an administrator to do. The uploads row belongs to
// the administrator, so the erasure leaves it alone; the record that claimed
// it is hard-deleted, so the orphan sweep can never consider it either. The
// file, its renditions and the kept original stay in storage for good with
// nothing pointing at them, and the deletion journal records the erasure as
// complete (grit#311).
//
// An error rather than a warning, and not because the data is certainly
// personal: the erasure is reported as done. A record that says a file is gone
// while the file is in the bucket is the one outcome the feature exists to
// prevent.
//
// Resources generated from v3.583.0 carry the cleanup. This finds the ones
// generated before it, which `grit upgrade` cannot fix for them: the
// registration lives in the resource's own model file, which is the project's
// to edit.
func checkOwnedFileErasure(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		m := p.models[name]
		// Not registered at all is a different state, and an older one: the
		// rows outlive their owner entirely. Saying "and its files too" about
		// a resource whose rows are not erased would be a second complaint
		// about a bigger thing.
		if !m.Erased || m.ErasesFiles {
			continue
		}
		if m.OwnerColumn == "" || m.AppendOnly || !m.hasFileColumn() {
			continue
		}
		out = append(out, Finding{
			Level:    "error",
			Resource: name,
			Message: "is erased with its owner and its files are not: a file somebody else uploaded " +
				"onto one of its rows, an administrator above all, stays in storage with nothing pointing at it",
			Fix: "regenerate it, which writes the cleanup, or change its model file to " +
				"erasure.RegisterWithCleanup(&" + name + "{}, \"" + m.OwnerColumn + "\", erase" + name + "Files) " +
				"with a cleanup that calls files.EraseOwnedFiles",
		})
	}
	return out
}
