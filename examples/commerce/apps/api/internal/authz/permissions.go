package authz

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
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Actions []Action `json:"actions"`
}

// Group and Module exist only to shape the admin UI's permission tree.
// Neither appears in a permission key.
type Group struct {
	Key      string    `json:"key"`
	Name     string    `json:"name"`
	Features []Feature `json:"features"`
}

type Module struct {
	Key    string  `json:"key"`
	Name   string  `json:"name"`
	Groups []Group `json:"groups"`
}

// Catalog is every permission this application understands: the built-in core
// plus whatever `grit generate resource` has registered.
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
					Key:  "publishing",
					Name: "Publishing",
					Features: []Feature{
						// Blog ships as a built-in resource (/api/admin/blogs).
						// It belongs in the catalog like any other so a role can
						// be granted blog access and the UI can gate on blogs.*.
						{Key: "blogs", Name: "Blogs", Actions: AllActions},
					},
				},
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

// generatedModules is maintained by `grit generate resource`.
// Everything between the markers is machine-written — `grit remove resource`
// takes it back out. Add hand-written permissions to coreModules() instead.
func generatedModules() []Module {
	return []Module{
		// grit:perms:auto-start
		{
			Key:  "collections",
			Name: "Collections",
			Groups: []Group{
				{
					Key:  "collections",
					Name: "Collections",
					Features: []Feature{
						{Key: "collections", Name: "Collections", Actions: AllActions},
					},
				},
			},
		},

		{
			Key:  "products",
			Name: "Products",
			Groups: []Group{
				{
					Key:  "products",
					Name: "Products",
					Features: []Feature{
						{Key: "products", Name: "Products", Actions: AllActions},
					},
				},
			},
		},

		{
			Key:  "pages",
			Name: "Pages",
			Groups: []Group{
				{
					Key:  "pages",
					Name: "Pages",
					Features: []Feature{
						{Key: "pages", Name: "Pages", Actions: AllActions},
					},
				},
			},
		},

		{
			Key:  "carts",
			Name: "Carts",
			Groups: []Group{
				{
					Key:  "carts",
					Name: "Carts",
					Features: []Feature{
						{Key: "carts", Name: "Carts", Actions: AllActions},
					},
				},
			},
		},

		{
			Key:  "cart_items",
			Name: "CartItems",
			Groups: []Group{
				{
					Key:  "cart_items",
					Name: "CartItems",
					Features: []Feature{
						{Key: "cart_items", Name: "CartItems", Actions: AllActions},
					},
				},
			},
		},

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
