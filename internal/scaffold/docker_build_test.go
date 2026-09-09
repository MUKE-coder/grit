package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The frontend images must build on a project that has never run pnpm install.
//
// That is the deploy path the VPS guide describes: scaffold, push, build. Two
// things stopped it, and both were invisible until an image was actually built.
//
// The Dockerfile did `COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./`
// and `grit new` does not write pnpm-lock.yaml, so the build failed at the
// first COPY with "/pnpm-lock.yaml: not found".
//
// Then it copied two workspace manifests by name, and both apps depend on
// @repo/upload as well as @repo/shared, so pnpm stopped with
// ERR_PNPM_WORKSPACE_PKG_NOT_FOUND. A list of workspace members maintained by
// hand in a Dockerfile goes stale the first time one is added, and had.
func TestFrontendDockerfilesBuildWithoutALockfile(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeDockerFiles(root, opts); err != nil {
		t.Fatalf("writeDockerFiles: %v", err)
	}

	for _, app := range []string{"web", "admin"} {
		path := filepath.Join(root, "apps", app, "Dockerfile")
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s: %v", app, err)
		}
		src := string(b)

		// A bare pnpm-lock.yaml in a COPY fails the build outright when the
		// file is absent; the glob is what makes it optional.
		if strings.Contains(src, "COPY package.json pnpm-lock.yaml ") ||
			strings.Contains(src, "pnpm-lock.yaml pnpm-workspace.yaml ./") {
			t.Errorf("%s COPYs pnpm-lock.yaml without a glob, so a project that "+
				"has not run pnpm install cannot build an image", app)
		}
		if !strings.Contains(src, "pnpm-lock.yaml*") {
			t.Errorf("%s does not tolerate a missing lockfile", app)
		}

		// --frozen-lockfile is right when a lockfile exists and fatal when it
		// does not, so it cannot be the only mode.
		if strings.Contains(src, "RUN pnpm install --frozen-lockfile\n") {
			t.Errorf("%s runs --frozen-lockfile unconditionally; there is no "+
				"lockfile to freeze to on a fresh project", app)
		}

		// Naming workspace members is what went stale.
		if strings.Contains(src, "COPY packages/shared/package.json") {
			t.Errorf("%s copies workspace members by name; @repo/upload was "+
				"added and this list did not follow", app)
		}
		if !strings.Contains(src, "COPY packages ./packages") {
			t.Errorf("%s does not copy the whole workspace, so a package added "+
				"later will not be in the image", app)
		}
	}
}
