package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// middleware.Identify, for an existing project.
//
// auth.go is the project's own file: people add claims to it, change where the
// token is read from, and Grit does not overwrite it on upgrade. So Identify
// arrives as an insertion rather than a rewrite, which also means a project
// that has moved things around gets told what to add instead of having its
// edits thrown away.
const identifyAnchor = `// RequireRole guards a route by role name, permission, or both.`

func identifyFunc() string {
	src := apiAuthMiddlewareGo()
	start := strings.Index(src, "// Identify is Auth without the wall")
	end := strings.Index(src, identifyAnchor)
	if start < 0 || end < 0 || end < start {
		return ""
	}
	return src[start:end]
}

func repairIdentifySource(src string) (string, []string, []string) {
	if strings.Contains(src, "func Identify(") {
		return src, nil, nil
	}
	block := identifyFunc()
	if block == "" || strings.Count(src, identifyAnchor) != 1 || !strings.Contains(src, "func Auth(db *gorm.DB") {
		return src, nil, []string{"internal/middleware/auth.go is not the one Grit wrote, so middleware.Identify was not added: copy it from a fresh project if you want a public route to be able to see the reader"}
	}
	return strings.Replace(src, identifyAnchor, block+identifyAnchor, 1),
		[]string{"middleware.Identify reads the session on a public route without refusing anonymous requests"}, nil
}

// repairIdentify applies it.
func repairIdentify(root string) error {
	path := filepath.Join(root, "apps", "api", "internal", "middleware", "auth.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairIdentifySource)
}
