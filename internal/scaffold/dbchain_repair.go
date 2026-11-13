package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// internal/dbchain reaches an existing project.
//
// The package itself arrives through writeCodegenRuntimeFiles, which both
// `grit new` and `grit upgrade` call, so a project gets the files whatever its
// age. Two edits do not arrive that way, because they are in files a project
// is allowed to have changed:
//
//   - database.go has to call dbchain.Apply, or a registered link is on no
//     connection at all and the package is decoration.
//   - the RequestID middleware has to put the id and the matched route on the
//     request's own context. It sets them on the gin context today, and a GORM
//     callback cannot see that: gin.Context values and Request.Context() values
//     are two different stores, and only the second travels into a service and
//     down to the database. Without it every recorded query is unattributed,
//     and "the same query many times in one request" has no request to be
//     many times inside of.
//
// Its own pass with its own gate rather than a check inside an existing
// repair, because a check added inside one that has already run on a project
// reaches new projects only.
func repairDBChain(root string, opts Options) error {
	apiRoot, found := dbChainAPIRoot(root, opts)
	if !found {
		return nil // no Go API here: a frontend-only or mobile-only project
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	// database.go: Apply the chain after the framework's own installs.
	// dbChainAPIRoot already found it, which is what makes this unconditional.
	database := filepath.Join(apiRoot, "internal", "database", "database.go")
	module := opts.Module()
	if err := repairSourceFile(root, m, database, func(src string) (string, []string, []string) {
		return repairDBChainConnect(src, module)
	}); err != nil {
		return err
	}

	// The RequestID middleware. In the same file as the logger since the
	// generator wrote them together.
	logger := filepath.Join(apiRoot, "internal", "middleware", "logger.go")
	if fileExists(logger) {
		if err := repairSourceFile(root, m, logger, func(src string) (string, []string, []string) {
			return repairDBChainRequestID(src, module)
		}); err != nil {
			return err
		}
	}

	return nil
}

// dbChainAPIRoot finds the API, in whichever layout this project uses.
//
// opts.APIRoot answers for the architecture grit.json names, and that was
// enough until it was not: a legacy flat single keeps its Go code at the
// project root while opts says api/, and the repair then looked at a path
// that does not exist and returned quietly. Every candidate is tried and the
// first one holding a database.go wins, which is the same thing `grit doctor`
// had to learn (v3.517.1).
func dbChainAPIRoot(root string, opts Options) (string, bool) {
	seen := map[string]bool{}
	for _, candidate := range []string{
		opts.APIRoot(root),
		filepath.Join(root, "apps", "api"),
		filepath.Join(root, "api"),
		root,
	} {
		if seen[candidate] {
			continue
		}
		seen[candidate] = true
		if fileExists(filepath.Join(candidate, "internal", "database", "database.go")) {
			return candidate, true
		}
	}
	return "", false
}

// dbChainConnectHook is the block Connect ends with, used by the template
// and by the repair below.
//
// One copy, not two. The repair has to produce the bytes the template holds,
// or a project that upgrades into it differs from a project generated with
// it, and the only thing that would notice is a test comparing the two.
const dbChainConnectHook = `	// Last, so a project's own callbacks see what the framework's have
	// already done: an encrypted column is ciphertext by here, rich text is
	// sanitised, an append-only table has already refused, and every
	// framework callback is registered, which is what makes an Anchor naming
	// one of them mean anything.
	//
	// This is also where GRIT_QUERY_LOG turns query recording on, with no
	// code change: see internal/dbchain.
	if err := dbchain.Apply(db); err != nil {
		return nil, fmt.Errorf("wiring the database callback chain: %w", err)
	}

`

func repairDBChainConnect(src, module string) (string, []string, []string) {
	// Guarded on the call, not on the word: a comment mentioning dbchain
	// matches a Contains("dbchain") and makes the repair a silent no-op.
	if strings.Contains(src, "dbchain.Apply(db)") {
		return src, nil, nil
	}

	// At the end of Connect, where the template has it: after every framework
	// callback is registered, so an Anchor naming one of them can be honoured.
	// Anchored on the line that closes the connection work, which is the one
	// part of the function every version of it has had.
	const anchor = "\tsqlDB, err := db.DB()"
	at := strings.Index(src, anchor)
	if at < 0 {
		return src, nil, []string{"could not find the end of Connect, so internal/dbchain is " +
			"not wired: add dbchain.Apply(db) before the pool settings by hand"}
	}
	out := src[:at] + dbChainConnectHook + src[at:]

	out, added := addModuleImport(out, module, "dbchain")
	if !added {
		return src, nil, []string{"could not add the dbchain import, so the file was left alone"}
	}
	return out, []string{"wires internal/dbchain onto every connection"}, nil
}

// requestContextLine is what RequestID gains.
const requestContextLine = `		// And onto the request's own context, which is the one that travels
		// into a service and down to the database. c.Set stores on the gin
		// context, and a GORM callback never sees that, so a recorded query
		// could not say which request issued it.
		//
		// FullPath is the matched route pattern (/api/posts/:id), not the
		// URL, so a report groups by endpoint rather than by row id.
		c.Request = c.Request.WithContext(
			dbchain.WithRequest(c.Request.Context(), requestID, c.Request.Method+" "+c.FullPath()))
`

func repairDBChainRequestID(src, module string) (string, []string, []string) {
	if strings.Contains(src, "dbchain.WithRequest(") {
		return src, nil, nil
	}
	const anchor = `		c.Header("X-Request-ID", requestID)
`
	at := strings.Index(src, anchor)
	if at < 0 {
		return src, nil, []string{"could not find the RequestID middleware, so a recorded query " +
			"cannot say which request issued it: see internal/dbchain.WithRequest"}
	}
	at += len(anchor)
	out := src[:at] + requestContextLine + src[at:]

	out, added := addModuleImport(out, module, "dbchain")
	if !added {
		return src, nil, []string{"could not add the dbchain import, so the file was left alone"}
	}
	return out, []string{"puts the request id and route where a database callback can read them"}, nil
}

// addModuleImport puts "<module>/internal/<pkg>" in the import block.
//
// Built by inserting a line, never by rewriting an anchor line: a
// strings.Replace that matches nothing succeeds and changes nothing, and the
// result is a file that names a package it does not import.
func addModuleImport(src, module, pkg string) (string, bool) {
	imp := fmt.Sprintf("\t%q\n", module+"/internal/"+pkg)
	if strings.Contains(src, imp) {
		return src, true
	}
	// Among the project's own imports, where the generated file has it. The
	// first version put it at the top of the block, which is inside the
	// standard-library group: it compiles, and gofmt leaves it there, because
	// gofmt sorts within a group and never moves a line between groups. The
	// upgraded file and the generated one then differed, and nothing looks.
	if out, ok := insertSortedImport(src, module, pkg); ok {
		return out, true
	}
	// No project group to join: append one after the last import line.
	open := strings.Index(src, "import (\n")
	if open < 0 {
		return src, false
	}
	// Named end, not close: shadowing a builtin is a lint finding and a
	// reader's double take for no gain.
	end := strings.Index(src[open:], "\n)\n")
	if end < 0 {
		return src, false
	}
	at := open + end + 1
	return src[:at] + "\n" + imp + src[at:], true
}

// insertSortedImport puts "<module>/internal/<pkg>" among the project's own
// imports, in the position gofmt would leave it.
//
// The project's own group is found by its module prefix, which is the one
// thing that distinguishes it from the standard library and the third-party
// group without parsing the file.
func insertSortedImport(src, module, pkg string) (string, bool) {
	want := module + "/internal/" + pkg
	line := "\t" + strconv.Quote(want) + "\n"
	if strings.Contains(src, line) {
		return src, true
	}

	prefix := "\t\"" + module + "/internal/"
	lines := strings.SplitAfter(src, "\n")
	for i, l := range lines {
		if !strings.HasPrefix(l, prefix) && !strings.HasPrefix(l, "\t_ \""+module+"/internal/") {
			continue
		}
		// The first of the project's own imports that sorts after this one.
		if strings.TrimSpace(l) > strconv.Quote(want) {
			return strings.Join(lines[:i], "") + line + strings.Join(lines[i:], ""), true
		}
	}
	return src, false
}

// ignoreQueryLog keeps a recorded query log out of the repository. It is tens
// of megabytes after a test run and it names the lines of a developer's own
// checkout.
//
// The template carries the line too, which covers every project whose
// .gitignore is untouched. This is for the ones people have edited, which the
// upgrade holds back, and it runs at the very end: the first version ran with
// the API repairs and a later step wrote .gitignore again, so the line was in
// the function and never in the project.
func ignoreQueryLog(root string) error {
	path := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(path) // #nosec G304 -- the project's own .gitignore
	if err != nil {
		return nil // no .gitignore is not this repair's problem
	}
	raw := string(data)
	if strings.Contains(raw, ".grit/queries.jsonl") {
		return nil
	}
	out := strings.TrimRight(raw, "\n") +
		"\n\n# A recorded query log: a development artefact, read by `grit doctor --queries`.\n" +
		".grit/queries.jsonl\n"
	return writeFile(path, out)
}
