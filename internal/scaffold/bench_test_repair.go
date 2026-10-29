package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// repairBenchTestGo fixes a benchmark that could never have run.
//
// Two faults in internal/handlers/bench_test.go, one hiding the other:
//
//   - BenchmarkAuthLogin registered with the password "benchpassword99" and
//     the surname "Bench". Registration refuses a password containing any
//     four character piece of the name or address, so it answered 422 and the
//     benchmark failed on its own require.
//   - Behind that, both benchmarks built an AuthService with no DB. Login
//     reads the user through AuthService.UserByEmail, so the first request
//     dereferenced nil. The 422 meant nobody reached it.
//
// The project's own ci.yml runs the benchmarks, so this was a red step in
// every scaffolded project from the first push.
//
// Only rewritten where Grit still owns the file; the upgrade guard keeps an
// edited one, and a project that deleted it keeps it deleted.
func repairBenchTestGo(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "handlers", "bench_test.go")
	if !fileExists(path) {
		return nil
	}
	// Both placeholders. {{MODULE}} was missed the first time, and the file
	// it wrote carried a literal "{{MODULE}}/internal/config" import that
	// stopped the project compiling.
	content := strings.ReplaceAll(apiBenchTestGo(), "{{MODULE}}", opts.Module())
	content = strings.ReplaceAll(content, "{{PROJECT}}", opts.ProjectName)
	if err := writeFile(path, content); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
