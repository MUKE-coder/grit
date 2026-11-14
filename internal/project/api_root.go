package project

import (
	"os"
	"path/filepath"
)

// APIRoot returns the directory holding the project's Go module.
//
// This exists because three places resolved it independently and two of them
// were wrong. doctor's copy learned about api/ in v3.517.1 and the other two
// did not, so on a --single project (and on the --api --admin shape that
// followed it) grit new --i18n failed with "no Go API found" one second after
// writing the module, and grit versions reported every framework library as
// "absent, not in go.mod" while reading a path that holds no go.mod at all.
// Neither failure looked like a path bug from the outside, which is what makes
// a duplicated resolver expensive: each copy fails in the vocabulary of
// whatever command owns it.
//
// The directory on disk decides, never grit.json, so a project scaffolded
// before a layout changed still resolves. Candidates are ordered most specific
// first: a monorepo that also happens to hold an api/ directory is still a
// monorepo.
//
// The project root is the last candidate rather than an error, because an
// api-only project keeps its module there, and so does anything scaffolded
// before the module moved. A caller that needs to know whether it found a real
// module should check for go.mod itself: this returns a best guess, and
// HasAPIModule answers the stricter question.
func APIRoot(root string) string {
	for _, candidate := range []string{
		filepath.Join(root, "apps", "api"), // triple, double, mobile, api
		filepath.Join(root, "api"),         // single, and --api --admin
	} {
		if hasGoMod(candidate) {
			return candidate
		}
	}
	return root
}

// HasAPIModule reports whether APIRoot found an actual Go module rather than
// falling through to the project root.
func HasAPIModule(root string) bool {
	return hasGoMod(APIRoot(root))
}

func hasGoMod(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "go.mod"))
	return err == nil && !info.IsDir()
}
