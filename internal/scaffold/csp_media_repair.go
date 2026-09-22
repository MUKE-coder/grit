package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The CSP had no media-src, so <video> and <audio> fell back to default-src
// 'self' and the browser refused every clip served by the API or by storage.
// Nothing played video before grit plugin add video, whose VideoPlayer was
// blocked in every Next.js app. Found building the Instagram blueprint.
const (
	// cspMediaSrc follows cspImgSrcTight in the Next.js and Vite templates:
	// the same origins a photo loads from.
	cspMediaSrc = "\"media-src 'self' blob: \" + API_ORIGIN + \" \" + STORAGE_ORIGIN,"

	// The nginx CSP a built Vite frontend is served with.
	nginxImgSrc      = "img-src 'self' data: blob: https:; "
	nginxImgMediaSrc = "img-src 'self' data: blob: https:; media-src 'self' blob: https:; "
)

// repairCSPMedia adds media-src to an existing project's frontends.
func repairCSPMedia(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	dirs := []string{
		filepath.Join(root, "apps", "web"),
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "frontend"),
	}
	for _, dir := range dirs {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js", "vite.config.ts"} {
			if path := filepath.Join(dir, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairCSPMediaSource); err != nil {
					return err
				}
			}
		}
		if path := filepath.Join(dir, "nginx.conf"); fileExists(path) {
			if err := repairTextFile(root, m, path, repairCSPMediaNginxSource); err != nil {
				return err
			}
		}
	}
	return nil
}

func repairCSPMediaSource(src string) (string, []string, []string) {
	if strings.Contains(src, "media-src") || !strings.Contains(src, "img-src") {
		return src, nil, nil
	}
	line := "  " + cspImgSrcTight + "\n"
	if strings.Count(src, line) != 1 {
		return src, nil, []string{"the Content-Security-Policy is not the one Grit wrote: add media-src with the API and storage origins, or the browser refuses to play any video or audio they serve"}
	}
	return strings.Replace(src, line, line+"  "+cspMediaSrc+"\n", 1),
		[]string{"media-src admits video and audio from the API and storage, which the browser refused"}, nil
}

func repairCSPMediaNginxSource(src string) (string, []string, []string) {
	if strings.Contains(src, "media-src") || strings.Count(src, nginxImgSrc) != 1 {
		return src, nil, nil
	}
	return strings.Replace(src, nginxImgSrc, nginxImgMediaSrc, 1),
		[]string{"media-src admits video and audio"}, nil
}
