package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// next/image refuses to fetch an image from a host that resolves to a private
// IP, and in development every host is one: the API on localhost serves stored
// files, and MinIO is localhost too. The optimizer answers 400, the page shows
// alt text, and the reason is a line in the terminal nobody is reading. So
// every scaffolded app showed broken images the moment it had an upload. Found
// building the storefront blueprint, whose every product is a photograph.
const (
	nextImagesOld = `  images: {
    remotePatterns: nextImageHosts,
  },`
	nextImagesNew = `  images: {
    remotePatterns: nextImageHosts,
    // Next refuses an upstream image whose host resolves to a private IP,
    // which in development is all of them: the API and MinIO both live on
    // localhost. Without this the optimizer answers 400 and the page shows
    // alt text, with the reason only in the terminal. Never in production,
    // where storage is a real origin and this guard is worth having.
    dangerouslyAllowLocalIP: isDev,
  },`
)

func repairNextImageLocalSource(src string) (string, []string, []string) {
	if strings.Contains(src, "dangerouslyAllowLocalIP") || !strings.Contains(src, "remotePatterns") {
		return src, nil, nil
	}
	if strings.Count(src, nextImagesOld) != 1 {
		return src, nil, []string{"the next/image config is not the one Grit wrote: add dangerouslyAllowLocalIP: isDev to it, or stored images do not load in development"}
	}
	return strings.Replace(src, nextImagesOld, nextImagesNew, 1),
		[]string{"next/image loads stored images from localhost in development, which it had started refusing"}, nil
}

// repairNextImageLocal applies it to an existing project's Next.js apps.
func repairNextImageLocal(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, app := range []string{"web", "admin", "docs"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js"} {
			if path := filepath.Join(root, "apps", app, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairNextImageLocalSource); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
