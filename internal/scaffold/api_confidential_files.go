package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The confidential registry: which tables hold rows an administrator may not
// read, and the permission each one answers to.
//
// Its own package, and a leaf one, for the same reason erasure is: the model
// files register themselves from their own init(), so the package they call
// must import nothing from the app. authz cannot hold it, because authz
// imports models and a model importing authz is a cycle. That was found by
// compiling a generated project rather than by reasoning about it.

// WriteConfidentialPackage writes internal/confidential under apiRoot.
// overwrite=false leaves an existing copy alone, for the generator running on
// a project that predates the package.
func WriteConfidentialPackage(apiRoot, module string, overwrite bool) error {
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "confidential", "confidential.go"): apiConfidentialGo(),
	}
	for path, content := range files {
		if !overwrite {
			if _, err := os.Stat(path); err == nil {
				continue
			}
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func apiConfidentialGo() string {
	return `// Package confidential is the set of tables an administrator may not read.
//
// A resource generated with --confidential registers itself here from its own
// model file, so removing the resource removes the registration with it and
// there is no central list to forget.
//
// It imports nothing from the app, because the models import it. That is not
// a style preference: authz, where the scoping helpers live, imports models,
// so a model file importing authz is a compile error.
//
// What needs this registry is everything that reads BELOW the service layer
// and cannot be expressed as a scoped query. A service scopes its own
// queries and asks nothing here. A backup archive is a file containing every
// table, and the only question it can answer is whether this caller may take
// it off the server at all.
package confidential

import "sync"

var (
	mu     sync.RWMutex
	tables = map[string]string{}
)

// Register marks table as holding rows readable only by their owner or by a
// holder of permission. Called from a model file's init().
func Register(table, permission string) {
	mu.Lock()
	defer mu.Unlock()
	tables[table] = permission
}

// Tables is every registered table and the permission it answers to.
func Tables() map[string]string {
	mu.RLock()
	defer mu.RUnlock()
	out := make(map[string]string, len(tables))
	for t, p := range tables {
		out[t] = p
	}
	return out
}

// Any reports whether this project has any confidential table at all, so a
// caller can skip the work when it has none.
func Any() bool {
	mu.RLock()
	defer mu.RUnlock()
	return len(tables) > 0
}
`
}
