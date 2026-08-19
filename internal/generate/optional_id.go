package generate

// OPTIONAL_ID_HELPER_SRC is emitted into a handler that has a nullable foreign
// key, which today means a resource with a self-reference (a tree).
//
// It exists because of one mismatch between HTTP and SQL. An HTML select with
// nothing chosen posts an empty string; a JSON client that means "no parent"
// sends "" for the same reason. SQL has no such value for an absent reference:
// a foreign key column holds a real key or it holds NULL.
//
// Writing "" instead fails outright on Postgres:
//
//	ERROR: insert or update on table "categories" violates foreign key
//	constraint "fk_categories_children" (SQLSTATE 23503)
//
// and on a database with constraints switched off it succeeds, storing a
// reference to a row that does not exist, which is worse: nothing fails until
// something tries to follow it.
const OPTIONAL_ID_HELPER_SRC = `
// optionalID turns an empty id from a form or a JSON body into NULL.
//
// A nullable foreign key column holds a real key or NULL, and "" is neither.
// Postgres refuses it with SQLSTATE 23503; SQLite accepts it and stores a
// dangling reference. Both are answered here, once, at the edge.
func optionalID(id string) *string {
	if id == "" {
		return nil
	}
	return &id
}
`
