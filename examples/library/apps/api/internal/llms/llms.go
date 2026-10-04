// Package llms is what an agent reads before it calls this API.
//
// # Why this exists beside the OpenAPI spec
//
// The spec at /docs/openapi.json is the contract, and it is the better artifact
// for anything that generates a client. It is also 300 KB of JSON that says
// nothing about the conventions every endpoint shares: that a list answers
// {"data": [...], "meta": {...}}, that an error carries a stable code, that an
// unversioned path is rewritten rather than served, that page_size has a
// ceiling. An agent reading the spec learns the shape of each endpoint and none
// of that.
//
// llms.txt is the orientation: a page of prose and a handful of links, at a
// conventional path, so a model pointed at this service has somewhere to start.
// llms-full.txt is every route this process serves, which is a thing the spec
// cannot tell you, because the spec documents the routes somebody wrote an
// override for and this is the router's own answer.
//
// The format is llmstxt.org: an H1, a blockquote summary, prose, then sections
// of links.
//
// # Both are gated with the API reference
//
// They describe the whole surface, admin routes included, so they are mounted
// under the same condition as /docs: off in production unless
// API_DOCS_PUBLIC=true. A route list is not a secret, and it is also not
// something to hand out by default.
package llms

import (
	"fmt"
	"sort"
	"strings"
)

// Route is one entry in the router's own table.
type Route struct {
	Method string
	Path   string
	// Handler is the function the route ends at, as the router names it. It is
	// the honest answer to "what runs", and it is not documentation: a handler
	// named well says a lot and one named Handler1 says nothing.
	Handler string
}

// Config is what the two files say about this particular service.
type Config struct {
	AppName string
	// BaseURL is the origin this API is reached at, without a trailing slash.
	BaseURL string
	// Version is the path segment every endpoint is served under, such as "v1".
	Version string
	// MaxPageSize is the ceiling the list endpoints clamp page_size to.
	MaxPageSize int
}

func (c Config) base() string {
	if c.BaseURL == "" {
		return ""
	}
	return strings.TrimRight(c.BaseURL, "/")
}

// Index renders llms.txt: what this API is, the conventions every endpoint
// shares, and where the machine-readable contract lives.
func Index(cfg Config) string {
	b := cfg.base()
	version := cfg.Version
	if version == "" {
		version = "v1"
	}
	maxPage := cfg.MaxPageSize
	if maxPage <= 0 {
		maxPage = 100
	}

	var out strings.Builder
	fmt.Fprintf(&out, "# %s API\n\n", cfg.AppName)
	fmt.Fprintf(&out, "> A REST API built with Grit (Go, Gin, GORM). This file is the orientation; "+
		"the contract is the OpenAPI spec linked below.\n\n")

	out.WriteString("## Start here\n\n")
	fmt.Fprintf(&out, "- [OpenAPI 3.1 spec](%s/docs/openapi.json): every endpoint, its parameters and its schemas. Generate a client from this.\n", b)
	fmt.Fprintf(&out, "- [Reference with a console](%s/docs): the same spec, rendered, with a request builder.\n", b)
	fmt.Fprintf(&out, "- [Every route this process serves](%s/llms-full.txt): the router's own table, including routes nobody wrote an override for.\n", b)
	fmt.Fprintf(&out, "- [Dependency health](%s/api/health): whether the database, cache, queue, mailer and object store are there, in four states.\n", b)
	out.WriteString("- [The framework](https://gritframework.dev/llms.txt): how a Grit project is laid out and what the CLI does, for changing this service rather than calling it.\n\n")

	fmt.Fprintf(&out, "## Versioning\n\nEvery endpoint is served under `/api/%s`. An unversioned `/api/...` request is\n"+
		"rewritten to the current version, so `/api/users` reaches `/api/%s/users`. Write the\n"+
		"version into your calls: the rewrite is a convenience for browsers, not a promise.\n"+
		"When a breaking change is unavoidable a `v2` group appears beside this one and `%s`\n"+
		"keeps answering the old way.\n\n", version, version, version)

	fmt.Fprintf(&out, "## Authentication\n\n"+
		"A JWT access token, as `Authorization: Bearer <token>` or as the `grit_access`\n"+
		"cookie. Browsers get the cookie, which is httpOnly so a cross-site script cannot\n"+
		"read it; anything else sends the header. `POST /api/%s/auth/login` returns the\n"+
		"token, and `POST /api/%s/auth/refresh` exchanges the refresh cookie for a new one.\n\n"+
		"Server-to-server calls to the public endpoints use an API key instead, as\n"+
		"`X-API-Key: <key>`. A publishable key reaches public endpoints only; a secret key\n"+
		"is never put in a browser.\n\n", version, version)

	out.WriteString("## What a response looks like\n\n" +
		"One item:\n\n" +
		"```json\n{ \"data\": { }, \"message\": \"User created successfully\" }\n```\n\n" +
		"A list, where `meta` is always present so a client never has to infer the page:\n\n" +
		"```json\n{ \"data\": [], \"meta\": { \"total\": 100, \"page\": 1, \"page_size\": 20, \"pages\": 5 } }\n```\n\n" +
		"An error, where `code` is stable and is the thing to switch on. The HTTP status\n" +
		"says what kind of failure it was; the code says which one:\n\n" +
		"```json\n{ \"error\": { \"code\": \"VALIDATION_ERROR\", \"message\": \"Email is required\",\n" +
		"            \"details\": { \"email\": \"This field is required\" } } }\n```\n\n" +
		"Statuses: 200, 201, 400, 401, 403, 404, 409, 422, 429, 500.\n\n")

	fmt.Fprintf(&out, "## Listing, searching and sorting\n\n"+
		"Every list endpoint takes the same parameters: `page`, `page_size`, `search`,\n"+
		"`sort_by`, `sort_order` (`asc` or `desc`), and `filter[field]=value` for the\n"+
		"columns that endpoint allows. `page_size` is clamped to %d rather than refused, so\n"+
		"asking for a million rows gets you %d and not a 400.\n\n"+
		"Which columns may be searched, sorted or filtered is decided per endpoint and is\n"+
		"a whitelist: a column that is not on it is ignored rather than passed to the\n"+
		"database, which is what stops a sort parameter from being an injection.\n\n",
		maxPage, maxPage)

	return out.String()
}

// Full renders llms-full.txt: the index, then every route the router holds,
// grouped by what it is about.
func Full(cfg Config, routes []Route) string {
	var out strings.Builder
	out.WriteString(Index(cfg))
	out.WriteString("---\n\n# Every route\n\n")
	fmt.Fprintf(&out, "%d routes, from the router itself rather than from a list somebody maintains, so a\n"+
		"route added to this service appears here without anybody remembering to add it.\n\n"+
		"They are grouped by what they are about, which is read from the path. The one part\n"+
		"of the authorization boundary that is also in the path is `/admin/`: those need an\n"+
		"administrator. Everything else under `/api/` needs a signed-in user except signing\n"+
		"in itself, registration, password reset, magic links and anything under `/public/`.\n"+
		"Which permission each endpoint wants is per endpoint and is not in a path at all,\n"+
		"so read the reference for that rather than inferring it from here.\n\n", len(routes))

	grouped := map[string][]Route{}
	for _, r := range routes {
		group := GroupOf(r.Path)
		grouped[group] = append(grouped[group], r)
	}
	names := make([]string, 0, len(grouped))
	for name := range grouped {
		names = append(names, name)
	}
	// Alphabetical within the API, and the mounted dashboards last. They are the
	// biggest group by a distance and they are somebody else's: the reference,
	// the profiler, the database browser. A reader scanning for an endpoint of
	// this application should not wade through them first.
	sort.Slice(names, func(i, j int) bool {
		iMounted, jMounted := names[i] == GroupMounted, names[j] == GroupMounted
		if iMounted != jMounted {
			return jMounted
		}
		return names[i] < names[j]
	})

	for _, name := range names {
		rs := grouped[name]
		sort.Slice(rs, func(i, j int) bool {
			if rs[i].Path != rs[j].Path {
				return rs[i].Path < rs[j].Path
			}
			return rs[i].Method < rs[j].Method
		})
		fmt.Fprintf(&out, "## %s\n\n", name)
		for _, r := range rs {
			if r.Handler != "" {
				fmt.Fprintf(&out, "- `%s %s` -> %s\n", r.Method, r.Path, r.Handler)
				continue
			}
			fmt.Fprintf(&out, "- `%s %s`\n", r.Method, r.Path)
		}
		out.WriteString("\n")
	}
	return out.String()
}

// GroupOf names what a path is about: the first segment after /api and the
// version, with /admin/ keeping the segment under it so the administrative
// endpoints for a resource are not mixed in with the ones anyone can reach.
//
// Anything outside /api is something this app mounts rather than part of its
// API: the reference, the database browser, the dashboards, static files.
// GroupMounted is the group for anything outside /api: the reference, the
// profiler, the database browser, static files. Named because two files sort
// on it and a typo in one of them would silently stop matching.
const GroupMounted = "Mounted dashboards and static files"

func GroupOf(path string) string {
	if !strings.HasPrefix(path, "/api/") {
		return GroupMounted
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// parts[0] is "api". Drop a version segment after it, so /api/v1/users and
	// an unversioned /api/users land together.
	parts = parts[1:]
	if len(parts) > 0 && isVersion(parts[0]) {
		parts = parts[1:]
	}
	if len(parts) == 0 {
		return "api"
	}
	if parts[0] == "admin" && len(parts) > 1 {
		return "admin/" + parts[1]
	}
	return parts[0]
}

// isVersion reports whether a segment is a version marker: v1, v2 and so on.
func isVersion(segment string) bool {
	if len(segment) < 2 || segment[0] != 'v' {
		return false
	}
	for _, r := range segment[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ShortHandlerName turns what Gin records, which is a fully qualified function
// name with a closure suffix, into the part worth reading:
// "acme/apps/api/internal/handlers.(*UserHandler).List-fm" becomes
// "UserHandler.List".
//
// Here rather than in routes.go because two things name handlers now, this file
// and the routes explorer, and two spellings of the same handler in two places
// is a difference somebody would try to explain.
func ShortHandlerName(full string) string {
	name := full
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimSuffix(name, "-fm")
	// Drop the package qualifier and the pointer receiver's punctuation.
	if i := strings.Index(name, "."); i >= 0 {
		name = name[i+1:]
	}
	name = strings.ReplaceAll(name, "(*", "")
	name = strings.ReplaceAll(name, ")", "")
	// An anonymous handler ends in func1, func2 and so on, which names nothing:
	// a closure declared inside Setup arrives here as "Setup.func1".
	if last := name[strings.LastIndex(name, ".")+1:]; name == "" || isAnonymous(last) {
		return ""
	}
	return name
}

// isAnonymous reports whether a name segment is Go's funcN for a closure.
func isAnonymous(segment string) bool {
	if !strings.HasPrefix(segment, "func") || len(segment) == len("func") {
		return false
	}
	for _, r := range segment[len("func"):] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
