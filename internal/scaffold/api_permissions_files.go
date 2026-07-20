package scaffold

// authzPermissionsGo emits internal/authz/permissions.go — the permission
// catalog and matcher.
//
// Design notes (why this shape):
//
//   - A permission key is "<resource>.<action>", e.g. "products.create". Two
//     segments, always. Grit resources are flat, so a deeper
//     module.submodule.feature.action scheme would be noise for generated code.
//     The module/group in the catalog is DISPLAY metadata for the admin UI, not
//     part of the key.
//
//   - Grants may contain wildcards ("*", "products.*", "*.view"). Roles store
//     what the operator authored, wildcards included, so a role granted
//     "products.*" automatically picks up an action added later. (Shoppleet
//     expands to concrete leaves on save, which is why its roles silently stop
//     inheriting new permissions.)
//
//   - Matching is strictly segment-wise. Shoppleet's matcher returns true at the
//     FIRST "*" in a pattern and ignores the rest, so "a.*.c.view" silently
//     behaves like "a.*" and over-grants. Granted() below does not do that, and
//     there is a regression test for exactly that case.
func authzPermissionsGo() string {
	return `package authz

import (
	"sort"
	"strings"
)

// Action is one verb in a permission key.
type Action string

const (
	ActionCreate Action = "create"
	ActionView   Action = "view"
	ActionEdit   Action = "edit"
	ActionDelete Action = "delete"
)

// Action sets, so a catalog entry can declare which verbs are meaningful.
// A report can be viewed but not created; the admin UI greys out the rest.
var (
	AllActions = []Action{ActionCreate, ActionView, ActionEdit, ActionDelete}
	ViewOnly   = []Action{ActionView}
	ViewEdit   = []Action{ActionView, ActionEdit}
)

// Feature is one permissionable thing — usually a resource.
// Key is the FIRST segment of the permission key: Key="products" yields
// "products.create", "products.view", ...
type Feature struct {
	Key     string   ` + "`json:\"key\"`" + `
	Name    string   ` + "`json:\"name\"`" + `
	Actions []Action ` + "`json:\"actions\"`" + `
}

// Group and Module exist only to shape the admin UI's permission tree.
// Neither appears in a permission key.
type Group struct {
	Key      string    ` + "`json:\"key\"`" + `
	Name     string    ` + "`json:\"name\"`" + `
	Features []Feature ` + "`json:\"features\"`" + `
}

type Module struct {
	Key    string  ` + "`json:\"key\"`" + `
	Name   string  ` + "`json:\"name\"`" + `
	Groups []Group ` + "`json:\"groups\"`" + `
}

// Catalog is every permission this application understands: the built-in core
// plus whatever ` + "`grit generate resource`" + ` has registered.
func Catalog() []Module {
	return append(coreModules(), generatedModules()...)
}

// coreModules are the permissions for what the scaffold itself ships.
// Edit freely — this is your application's catalog, not framework internals.
func coreModules() []Module {
	return []Module{
		{
			Key:  "access",
			Name: "Access",
			Groups: []Group{
				{
					Key:  "identity",
					Name: "Identity",
					Features: []Feature{
						{Key: "users", Name: "Users", Actions: AllActions},
						{Key: "roles", Name: "Roles & permissions", Actions: AllActions},
					},
				},
			},
		},
		{
			Key:  "content",
			Name: "Content",
			Groups: []Group{
				{
					Key:  "media",
					Name: "Media",
					Features: []Feature{
						{Key: "uploads", Name: "Uploads", Actions: []Action{ActionCreate, ActionView, ActionDelete}},
					},
				},
			},
		},
		{
			Key:  "system",
			Name: "System",
			Groups: []Group{
				{
					Key:  "operations",
					Name: "Operations",
					Features: []Feature{
						{Key: "system", Name: "System health", Actions: ViewOnly},
						{Key: "jobs", Name: "Background jobs", Actions: ViewEdit},
						{Key: "backups", Name: "Backups", Actions: AllActions},
						{Key: "audit", Name: "Activity log", Actions: ViewOnly},
					},
				},
			},
		},
	}
}

// generatedModules is maintained by ` + "`grit generate resource`" + `.
// Everything between the markers is machine-written — ` + "`grit remove resource`" + `
// takes it back out. Add hand-written permissions to coreModules() instead.
func generatedModules() []Module {
	return []Module{
		// grit:perms:auto-start
		// grit:perms:auto-end
	}
}

// Keys returns every concrete permission key in the catalog, sorted.
// This is the denominator in the admin UI's "N / total granted".
func Keys() []string {
	var keys []string
	for _, m := range Catalog() {
		for _, g := range m.Groups {
			for _, f := range g.Features {
				for _, a := range f.Actions {
					keys = append(keys, f.Key+"."+string(a))
				}
			}
		}
	}
	sort.Strings(keys)
	return keys
}

// Granted reports whether any of the grants authorises key.
//
// Matching is segment-wise:
//   - "*" matches everything.
//   - a "*" segment matches exactly one segment ("*.view" matches
//     "products.view" but not "products.variants.view").
//   - a pattern ending in "*" also matches longer keys, so "products.*"
//     covers "products.create".
//
// The middle-segment rule is the important one: a pattern must be checked to
// its END. Bailing out at the first "*" (as Shoppleet does) turns
// "products.*.view" into "products.*" and grants delete along with view.
func Granted(grants []string, key string) bool {
	for _, g := range grants {
		if matches(g, key) {
			return true
		}
	}
	return false
}

func matches(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == key {
		return true
	}

	p := strings.Split(pattern, ".")
	k := strings.Split(key, ".")

	for i, seg := range p {
		// A trailing "*" absorbs whatever is left of the key.
		if seg == "*" && i == len(p)-1 {
			return i < len(k)
		}
		if i >= len(k) {
			return false
		}
		if seg != "*" && seg != k[i] {
			return false
		}
	}
	// Every pattern segment matched; the key must not have extra segments.
	return len(p) == len(k)
}

// Expand resolves grants (wildcards included) against the catalog, returning the
// concrete keys they authorise. The API hands this to clients so the frontend
// never has to reimplement matching — a duplicated matcher is how Shoppleet's Go
// and TypeScript rules drifted apart.
func Expand(grants []string) []string {
	// Non-nil so a role with no grants marshals to [] rather than null. The
	// roles UI maps over this and reads .length; null crashed the whole page
	// for the default USER role, which grants nothing.
	out := []string{}
	for _, key := range Keys() {
		if Granted(grants, key) {
			out = append(out, key)
		}
	}
	return out
}

// HasAll reports whether grants cover the entire catalog — used to show an
// "everything" badge rather than "175 / 175".
func HasAll(grants []string) bool {
	for _, g := range grants {
		if g == "*" {
			return true
		}
	}
	keys := Keys()
	if len(keys) == 0 {
		return false
	}
	return len(Expand(grants)) == len(keys)
}
`
}

// authzPermissionsTestGo emits tests for the matcher. These are shipped into the
// generated project deliberately: the catalog is application data that teams
// edit, and this pins the semantics they rely on.
func authzPermissionsTestGo() string {
	return `package authz

import "testing"

func TestGranted(t *testing.T) {
	cases := []struct {
		name  string
		grant string
		key   string
		want  bool
	}{
		{"exact", "products.create", "products.create", true},
		{"exact mismatch", "products.create", "products.delete", false},
		{"superuser", "*", "anything.at.all", true},
		{"resource wildcard", "products.*", "products.create", true},
		{"resource wildcard is scoped", "products.*", "orders.create", false},
		{"action wildcard", "*.view", "products.view", true},
		{"action wildcard is scoped", "*.view", "products.delete", false},
		{"no partial segment match", "prod.*", "products.view", false},

		// The Shoppleet bug: a matcher that returns true at the first "*"
		// treats this as "products.*" and wrongly grants delete.
		{"middle wildcard checks the tail", "products.*.view", "products.variants.delete", false},
		{"middle wildcard matches its tail", "products.*.view", "products.variants.view", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Granted([]string{tc.grant}, tc.key); got != tc.want {
				t.Errorf("Granted([%q], %q) = %v, want %v", tc.grant, tc.key, got, tc.want)
			}
		})
	}
}

func TestGrantedAnyOf(t *testing.T) {
	grants := []string{"uploads.view", "products.*"}
	if !Granted(grants, "products.delete") {
		t.Error("expected products.* to authorise products.delete")
	}
	if !Granted(grants, "uploads.view") {
		t.Error("expected the exact grant to match")
	}
	if Granted(grants, "users.delete") {
		t.Error("ungranted key must not pass")
	}
	if Granted(nil, "products.view") {
		t.Error("no grants must authorise nothing")
	}
}

func TestKeysAndExpand(t *testing.T) {
	keys := Keys()
	if len(keys) == 0 {
		t.Fatal("catalog produced no keys")
	}

	// Every catalog key must be reachable via the superuser grant.
	if got := len(Expand([]string{"*"})); got != len(keys) {
		t.Errorf("Expand(*) = %d keys, want %d", got, len(keys))
	}
	if !HasAll([]string{"*"}) {
		t.Error("HasAll(*) must be true")
	}

	// A resource wildcard expands only within that resource.
	for _, k := range Expand([]string{"users.*"}) {
		if len(k) < 6 || k[:6] != "users." {
			t.Errorf("users.* leaked %q", k)
		}
	}
	if HasAll([]string{"users.*"}) {
		t.Error("a single resource must not satisfy HasAll")
	}
}
`
}
