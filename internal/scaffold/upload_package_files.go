package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	uploadpkg "github.com/MUKE-coder/grit/v3/packages/upload"
)

// writeUploadPackageFiles writes packages/upload: client-side image
// optimisation and direct-to-storage uploads.
//
// Uploads go from the browser straight to S3 via a presigned URL and never
// pass through the API, so the optimisation has to happen before the bytes
// leave the device. That is not only where it has to be, it is where it is
// best: the browser has a lossy WebP encoder built in, which is exactly what
// the pure-Go server backend cannot do without cgo. Measured in Chromium, a
// 5 MB photograph comes out at 35 KB, matching libvips with no cgo anywhere.
//
// The TypeScript is embedded from packages/upload rather than held in Go
// string literals here, because that same directory is published to npm as
// @gritframework/upload. Two copies would agree only until somebody edited
// one of them.
func writeUploadPackageFiles(root string, opts Options) error {
	pkgRoot := filepath.Join(root, "packages", "upload")

	entries, err := uploadpkg.Sources.ReadDir("src")
	if err != nil {
		return fmt.Errorf("reading embedded upload sources: %w", err)
	}
	// The optional stylesheet for <Dropzone>.
	css, err := uploadpkg.Sources.ReadFile("styles.css")
	if err != nil {
		return fmt.Errorf("reading embedded upload styles: %w", err)
	}
	if err := writeFile(filepath.Join(pkgRoot, "styles.css"), string(css)); err != nil {
		return fmt.Errorf("writing styles.css: %w", err)
	}
	for _, e := range entries {
		// The component test stays in the library.
		//
		// It needs jsdom, Testing Library and a second copy of react-dom, and
		// under pnpm's strict layout that second copy talks to a different
		// React instance and every render throws. None of that is worth
		// carrying into a project that did not write the component. The logic
		// tests ship, because they need nothing but vitest.
		if strings.HasSuffix(e.Name(), ".test.tsx") {
			continue
		}
		body, err := uploadpkg.Sources.ReadFile("src/" + e.Name())
		if err != nil {
			return fmt.Errorf("reading embedded %s: %w", e.Name(), err)
		}
		content := strings.ReplaceAll(string(body), "{{MODULE}}", opts.Module())
		if err := writeFile(filepath.Join(pkgRoot, "src", e.Name()), content); err != nil {
			return fmt.Errorf("writing %s: %w", e.Name(), err)
		}
	}

	readme, err := uploadpkg.Sources.ReadFile("README.md")
	if err != nil {
		return fmt.Errorf("reading embedded upload README: %w", err)
	}
	// The published name means nothing inside a workspace that resolves it from
	// disk, so the copy in a generated project refers to itself by its
	// workspace name.
	readmeBody := strings.ReplaceAll(string(readme), "@gritframework/upload", "@repo/upload")

	files := map[string]string{
		filepath.Join(pkgRoot, "package.json"):  uploadPackageJSON(),
		filepath.Join(pkgRoot, "tsconfig.json"): uploadTSConfig(),
		filepath.Join(pkgRoot, "README.md"):     readmeBody,
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// uploadPackageJSON is the workspace flavour.
//
// Deliberately not the published package.json. Inside a monorepo the package
// is consumed straight from source through workspace:*, so there is no build
// step, no dist, and no version to keep in step with a registry. The published
// one lives in packages/upload and carries the tsup build and the real name.
func uploadPackageJSON() string {
	return `{
  "name": "@repo/upload",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "exports": {
    ".": "./src/index.ts",
    "./web": "./src/web.ts",
    "./expo": "./src/expo.ts",
    "./react": "./src/react.ts",
    "./ui": "./src/ui.tsx",
    "./transport": "./src/transport.ts",
    "./styles.css": "./styles.css"
  },
  "scripts": {
    "test": "vitest run",
    "typecheck": "tsc --noEmit"
  },
  "peerDependencies": {
    "react": ">=18"
  },
  "peerDependenciesMeta": {
    "react": { "optional": true }
  },
  "devDependencies": {
    "typescript": "^5.6.3",
    "vitest": "^2.1.4"
  }
}
`
}

func uploadTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "noEmit": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "jsx": "react-jsx"
  },
  "include": ["src/**/*.ts"]
}
`
}
