package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairCSPUploadOrigin admits a bucket whose presigned writes go somewhere
// other than where its files are read from.
//
// A CSP source is an origin, so listing only the read origin blocks the PUT: the
// API signs a perfectly good URL and the browser refuses it, reported as a
// console violation and never as an HTTP status, so the upload looks like it
// simply did nothing. A managed bucket arranges this by default, serving reads
// from its own host and signing writes for the underlying S3 endpoint. Reported
// against Laravel Cloud in grit#92, where reads come from fls-*.laravel.cloud
// and the signature is for *.r2.cloudflarestorage.com.
//
// Unset, NEXT_PUBLIC_STORAGE_UPLOAD_URL (VITE_STORAGE_UPLOAD_URL on a TanStack
// frontend) adds nothing to the policy, so this costs a project that does not
// need it one dead branch.
//
// Runs after repairCSPWebSocket, whose shape it matches on.
func repairCSPUploadOrigin(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, app := range []string{"admin", "web"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js"} {
			if path := filepath.Join(root, "apps", app, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairCSPUploadOriginNextSource); err != nil {
					return err
				}
			}
		}
		if path := filepath.Join(root, "apps", app, "vite.config.ts"); fileExists(path) {
			if err := repairTextFile(root, m, path, repairCSPUploadOriginViteSource); err != nil {
				return err
			}
		}
	}
	return nil
}

const nextUploadOriginBlock = `// Where presigned uploads are actually sent, when that is not where the files
// are read back from. A CSP source is an origin, so listing only the read origin
// blocks the PUT: the API signs a perfectly good URL and the browser refuses it,
// which shows up as a console violation and never as an HTTP status. A managed
// bucket does this by default, serving reads from its own CDN host and signing
// writes for the underlying S3 endpoint. Unset, it costs nothing.
const STORAGE_UPLOAD_ORIGIN = process.env.NEXT_PUBLIC_STORAGE_UPLOAD_URL
  ? " " + toOrigin(process.env.NEXT_PUBLIC_STORAGE_UPLOAD_URL)
  : "";
`

func repairCSPUploadOriginNextSource(src string) (string, []string, []string) {
	if strings.Contains(src, nextUploadOriginDecl) || !strings.Contains(src, "connect-src") {
		return src, nil, nil
	}
	if strings.Count(src, nextConnectSrcNew) != 1 || !strings.Contains(src, "const STORAGE_ORIGIN = toOrigin(") {
		return src, nil, []string{"the Content-Security-Policy is not the one Grit wrote: if presigned uploads go to a different host from the one files are served from, add that origin to connect-src or the browser blocks every upload"}
	}
	// After the STORAGE_ORIGIN declaration it reads from, wherever that line ends.
	at := strings.Index(src, "const STORAGE_ORIGIN = toOrigin(")
	line := strings.Index(src[at:], "\n")
	if line < 0 {
		return src, nil, nil
	}
	at += line + 1
	out := src[:at] + nextUploadOriginBlock + src[at:]
	out = strings.Replace(out, nextConnectSrcNew, nextConnectSrcUpload, 1)
	return out, []string{"connect-src admits NEXT_PUBLIC_STORAGE_UPLOAD_URL, for a bucket whose presigned writes and public reads are on different hosts"}, nil
}

const viteUploadOriginDecl = "const STORAGE_UPLOAD_ORIGIN = viteEnv.VITE_STORAGE_UPLOAD_URL"

const viteUploadOriginBlock = `
// Where presigned uploads are actually sent, when that is not where the files
// are read back from. A CSP source is an origin, so listing only the read origin
// blocks the PUT: the API signs a perfectly good URL and the browser refuses it,
// which shows up as a console violation and never as an HTTP status. A managed
// bucket does this by default, serving reads from its own CDN host and signing
// writes for the underlying S3 endpoint. Unset, it costs nothing.
const STORAGE_UPLOAD_ORIGIN = viteEnv.VITE_STORAGE_UPLOAD_URL
  ? ' ' + toOrigin(viteEnv.VITE_STORAGE_UPLOAD_URL)
  : ''`

const (
	viteStorageOriginDecl = "const STORAGE_ORIGIN = toOrigin(viteEnv.VITE_STORAGE_URL || 'http://localhost:9002')"
	viteConnectSrcOld     = `"connect-src 'self' ws: wss: " + API_ORIGIN + " " + STORAGE_ORIGIN + (isDev`
	viteConnectSrcNew     = `"connect-src 'self' ws: wss: " + API_ORIGIN + " " + STORAGE_ORIGIN + STORAGE_UPLOAD_ORIGIN + (isDev`
)

func repairCSPUploadOriginViteSource(src string) (string, []string, []string) {
	if strings.Contains(src, viteUploadOriginDecl) || !strings.Contains(src, "connect-src") {
		return src, nil, nil
	}
	if strings.Count(src, viteStorageOriginDecl) != 1 || strings.Count(src, viteConnectSrcOld) != 1 {
		return src, nil, []string{"the Content-Security-Policy is not the one Grit wrote: if presigned uploads go to a different host from the one files are served from, add that origin to connect-src or the browser blocks every upload"}
	}
	out := strings.Replace(src, viteStorageOriginDecl, viteStorageOriginDecl+viteUploadOriginBlock, 1)
	out = strings.Replace(out, viteConnectSrcOld, viteConnectSrcNew, 1)
	return out, []string{"connect-src admits VITE_STORAGE_UPLOAD_URL, for a bucket whose presigned writes and public reads are on different hosts"}, nil
}
