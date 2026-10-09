package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The realtime socket is the API origin over ws: or wss:, and a CSP source
// matches its scheme exactly: http://api.example.com does not admit
// ws://api.example.com. Development hid it, because the dev CSP allows every
// ws: and wss: for the hot-reload socket. Under next start the browser blocked
// the socket, so notifications, presence and every other live update in a
// Next.js web app or admin panel were dead in production. Found building the
// WhatsApp blueprint.
const (
	// One origin, read the same way the client reads it.
	//
	// The CSP authorises what the app is about to call, so the two must agree
	// about where the API is. This one read only the public name while the
	// client fell back to API_URL, which is the name the generated .env
	// actually sets: move APP_PORT and the browser blocks every request,
	// naming a port the developer never configured.
	nextAPIOriginLine = "const API_ORIGIN = toOrigin(\n  process.env.NEXT_PUBLIC_API_URL || process.env.API_URL || \"http://localhost:8080\",\n);\n"

	// What projects written before v3.395.0 have. Kept so the repairs below
	// still recognise a config they wrote, and so an upgrade can replace it.
	nextAPIOriginLineOld = "const API_ORIGIN = toOrigin(process.env.NEXT_PUBLIC_API_URL || \"http://localhost:8080\");\n"
	nextAPIWSOrigin   = `// The realtime socket is the API origin over ws: or wss:. A CSP source matches
// its scheme exactly, so API_ORIGIN does not admit it, and without this the
// browser blocks every live update once the dev-only ws: wss: below is gone.
const API_WS_ORIGIN = API_ORIGIN.replace(/^http/, "ws");
`
	nextConnectSrcOld = `"connect-src 'self' " + API_ORIGIN + " " + STORAGE_ORIGIN + (isDev`
	nextConnectSrcNew = `"connect-src 'self' " + API_ORIGIN + " " + API_WS_ORIGIN + " " + STORAGE_ORIGIN + (isDev`

	// The upload origin, for a bucket whose presigned PUTs go somewhere other
	// than where its files are served from.
	nextUploadOriginDecl = `const STORAGE_UPLOAD_ORIGIN = process.env.NEXT_PUBLIC_STORAGE_UPLOAD_URL`
	nextConnectSrcUpload = `"connect-src 'self' " + API_ORIGIN + " " + API_WS_ORIGIN + " " + STORAGE_ORIGIN + STORAGE_UPLOAD_ORIGIN + (isDev`

	// The nginx CSP a Vite frontend is served with allowed http: and https:
	// connections but no socket at all.
	nginxConnectSrcOld = "connect-src 'self' https: http:;"
	nginxConnectSrcNew = "connect-src 'self' https: http: wss: ws:;"
)

// repairCSPWebSocket gives an existing project's frontends the socket origin.
func repairCSPWebSocket(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, app := range []string{"admin", "web"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js"} {
			if path := filepath.Join(root, "apps", app, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairCSPWebSocketNextSource); err != nil {
					return err
				}
			}
		}
		if path := filepath.Join(root, "apps", app, "nginx.conf"); fileExists(path) {
			if err := repairTextFile(root, m, path, repairCSPWebSocketNginxSource); err != nil {
				return err
			}
		}
	}
	return nil
}

func repairCSPWebSocketNextSource(src string) (string, []string, []string) {
	if strings.Contains(src, "API_WS_ORIGIN") || !strings.Contains(src, "connect-src") {
		return src, nil, nil
	}
	// Either spelling of the origin line: a project may predate the API_URL
	// fallback or already have it.
	anchor := nextAPIOriginLine
	if strings.Count(src, anchor) != 1 {
		anchor = nextAPIOriginLineOld
	}
	if strings.Count(src, anchor) != 1 || strings.Count(src, nextConnectSrcOld) != 1 {
		return src, nil, []string{"the Content-Security-Policy is not the one Grit wrote: add the API origin with ws: or wss: in place of http: or https: to connect-src, or the browser blocks the realtime socket in production"}
	}
	out := strings.Replace(src, anchor, anchor+nextAPIWSOrigin, 1)
	out = strings.Replace(out, nextConnectSrcOld, nextConnectSrcNew, 1)
	return out, []string{"connect-src admits the API's realtime socket, which the browser blocked in production"}, nil
}

// repairCSPAPIURLFallback gives an existing project's frontends the API_URL
// fallback in their CSP origin.
//
// Without it the policy authorises one host and the app calls another the
// moment APP_PORT moves, and the browser reports that as a console violation
// rather than a status: a sign-in that does nothing, with nothing in the
// network tab that looks like a failure.
func repairCSPAPIURLFallback(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, app := range []string{"admin", "web"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js"} {
			if path := filepath.Join(root, "apps", app, name); fileExists(path) {
				if err := repairTextFile(root, m, path, repairCSPAPIURLFallbackSource); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func repairCSPAPIURLFallbackSource(src string) (string, []string, []string) {
	if strings.Contains(src, "process.env.API_URL") {
		return src, nil, nil
	}
	if strings.Count(src, nextAPIOriginLineOld) != 1 {
		return src, nil, nil
	}
	return strings.Replace(src, nextAPIOriginLineOld, nextAPIOriginLine, 1),
		[]string{"the Content-Security-Policy reads API_URL, so moving APP_PORT no longer blocks every request"}, nil
}

func repairCSPWebSocketNginxSource(src string) (string, []string, []string) {
	if strings.Count(src, nginxConnectSrcOld) != 1 {
		return src, nil, nil
	}
	return strings.Replace(src, nginxConnectSrcOld, nginxConnectSrcNew, 1),
		[]string{"connect-src admits the API's realtime socket"}, nil
}

// buildArtefactIgnores are written by tsc, next and expo on every build or dev
// run. Unignored, they turned up in the first commit of every project that
// built before committing.
var buildArtefactIgnores = []string{"*.tsbuildinfo", "next-env.d.ts", ".expo/"}

const buildArtefactComment = "# TypeScript, Next.js and Expo write these on every build or dev run."

// ignoreBuildArtefacts adds whichever buildArtefactIgnores .gitignore lacks,
// line by line, so a project that already has the first two (v3.297.0) still
// gets .expo/.
func ignoreBuildArtefacts(root string) error {
	path := filepath.Join(root, ".gitignore")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading .gitignore: %w", err)
	}
	src := string(raw)
	have := map[string]bool{}
	for _, line := range strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n") {
		have[strings.TrimSpace(line)] = true
	}
	var missing []string
	for _, line := range buildArtefactIgnores {
		if !have[line] {
			missing = append(missing, line)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
	}
	block := strings.Join(missing, nl) + nl
	if !have[buildArtefactComment] {
		block = buildArtefactComment + nl + block
	}
	if src != "" && !strings.HasSuffix(src, "\n") {
		src += nl
	}
	if err := os.WriteFile(path, []byte(src+nl+block), 0o644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}
	fmt.Printf("  ✓ .gitignore: %s ignored\n", strings.Join(missing, ", "))
	fmt.Println("    One already committed stays tracked until you git rm --cached it.")
	return nil
}
