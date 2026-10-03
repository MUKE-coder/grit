package scaffold

// apiAccessGo emits internal/access/access.go: the types and lookup for the
// generated access table.
//
// The table itself (registry.go beside this) is written by internal/accessgen
// after the scaffold is on disk, because it is read out of routes.go and the
// per-resource route files, which do not exist until then. This half is the
// stable half: it never changes when a route is added.
func apiAccessGo() string { return tmpl("api/access/access.go") }
