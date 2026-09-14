package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M6 to M10. The new text lives here once: templates
// splice it in, and repairDeployHardening puts it into an existing project,
// anchored on what Grit generated before.

// M6: SAML metadata is fetched through the SSRF guard.
const (
	samlFetchOld = "\treturn samlsp.FetchMetadata(ctx, http.DefaultClient, *u)\n"
	samlFetchNew = "\t// An administrator types this URL, so it goes through the SSRF guard like\n" +
		"\t// every other server-side fetch: no loopback, private or metadata addresses.\n" +
		"\treturn samlsp.FetchMetadata(ctx, safefetch.Client, *u)\n"
)

// M7: build contexts leave local state out of images.
const (
	dockerIgnoreOld = "node_modules\n.next\n.turbo\ndist\n*.log\n.env\n.env.local\n.git\n"
	dockerIgnoreNew = `# Everything not listed here is sent to the Docker daemon as the build context,
# and a Dockerfile that copies the context bakes it into the image. The **/
# patterns matter: a bare name only matches at the context root, which is how
# apps/web/.env.local, local databases and dev binaries ended up in images.
**/node_modules
**/.next
**/.turbo
dist
**/tmp
**/*.log
**/.env
**/.env.*
!**/.env.example
**/*.db
**/*.db-shm
**/*.db-wal
**/*.sqlite
**/*.exe
**/coverage
.git
.grit
e2e
`
	// The API image builds from apps/api, which the root file does not cover.
	apiDockerIgnore = `# The API image's build context is apps/api. Keep local state out of it: a dev
# database, the air build output, local .env files and binaries.
tmp
**/*.exe
**/*.db
**/*.db-shm
**/*.db-wal
**/*.sqlite
**/.env
**/.env.*
!**/.env.example
**/*.log
coverage.out
`
)

// M8: the production stack uses Postgres, a Redis password and pinned images.
const (
	pgbouncerImageOld = "image: edoburu/pgbouncer:latest"
	pgbouncerImageNew = "image: edoburu/pgbouncer:v1.25.2-p0"
	minioImageOld     = "image: minio/minio\n"
	// Docker Hub no longer serves minio/minio; MinIO publishes to quay.io.
	minioImageNew = "image: quay.io/minio/minio:RELEASE.2025-09-07T16-13-09Z\n"

	composeDBProviderAnchor = "      APP_ENV: production\n"
	composeDBProvider       = "      APP_ENV: production\n" +
		"      # Postgres, whatever DB_PROVIDER a development .env says. SQLite here lives\n" +
		"      # inside the container, and a redeploy starts from an empty database.\n" +
		"      DB_PROVIDER: postgres\n"
	composeRedisURLOld = "      REDIS_URL: redis://redis:6379\n"
	composeRedisURLNew = "      REDIS_URL: redis://:${REDIS_PASSWORD:?set REDIS_PASSWORD in .env}@redis:6379\n"
	composeRedisAuth   = "    # A password, with no default. Anything on the Docker network could otherwise\n" +
		"    # read the cache and the job queue. REDISCLI_AUTH lets the health check in.\n" +
		"    command: [\"redis-server\", \"--requirepass\", \"${REDIS_PASSWORD:?set REDIS_PASSWORD in .env}\"]\n" +
		"    environment:\n" +
		"      REDISCLI_AUTH: ${REDIS_PASSWORD:?set REDIS_PASSWORD in .env}\n"

	sqliteProductionAnchor = "\t\tcfg.GORMStudioDisableSQL = true\n\t}\n"
	sqliteProductionCheck  = "\t\tcfg.GORMStudioDisableSQL = true\n\n" +
		"\t\t// SQLite in production keeps the database inside the container, and a\n" +
		"\t\t// redeploy replaces the container: every deploy began from an empty\n" +
		"\t\t// database while the Postgres the production compose file starts sat\n" +
		"\t\t// unused. A deployment that does want SQLite, on a volume it backs up,\n" +
		"\t\t// says so.\n" +
		"\t\tif strings.HasPrefix(cfg.DatabaseURL, \"sqlite:\") && getEnv(\"ALLOW_SQLITE_IN_PRODUCTION\", \"false\") != \"true\" {\n" +
		"\t\t\treturn nil, fmt.Errorf(\"APP_ENV=production is using SQLite: set DB_PROVIDER=postgres, or ALLOW_SQLITE_IN_PRODUCTION=true if the database file is on a volume you back up\")\n" +
		"\t\t}\n" +
		"\t}\n"
)

// M10: a marker cookie the web app's middleware can read, and the admin gate.
const (
	signedInCookieSet = "\t// grit_signed_in carries no secret. It lives as long as the session so the\n" +
		"\t// web app's middleware can tell a signed-in browser from a stranger on the\n" +
		"\t// same host: grit_access expires with its token and grit_refresh is scoped\n" +
		"\t// to the auth routes, so neither can be seen on an admin page.\n" +
		"\tc.SetCookie(\"grit_signed_in\", \"1\", refreshSeconds, \"/\", \"\", secure, true)\n"
	signedInCookieClear = "\tc.SetCookie(\"grit_signed_in\", \"\", -1, \"/\", \"\", secure, true)\n"

	adminGateTS = `// ── The admin panel ──────────────────────────────────────────────────────
//
// The admin checks the session in the browser, so a visitor who never signed
// in downloaded every admin page before being sent to the login. The API sets
// grit_signed_in, a marker with no secret in it, beside its session cookies. A
// browser only sends it here when the web app and the API share a host
// (localhost in development, or one domain behind a proxy); on separate hosts
// it never arrives, so this gate stands aside and the page's own check still
// applies. The API authorises every request either way.
const ADMIN_PUBLIC_PATHS = [
  "/admin/login",
  "/admin/sign-up",
  "/admin/forgot-password",
  "/admin/reset-password",
  "/admin/verify-email",
  "/admin/callback",
];

function apiHostname(): string {
  try {
    return new URL(process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").hostname;
  } catch {
    return "";
  }
}

function adminGate(request: NextRequest): NextResponse | null {
  const { pathname, hostname } = request.nextUrl;
  if (pathname !== "/admin" && !pathname.startsWith("/admin/")) return null;
  if (ADMIN_PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(p + "/"))) return null;
  if (apiHostname() !== hostname) return null;
  if (request.cookies.has("grit_signed_in") || request.cookies.has("grit_access")) return null;
  const url = request.nextUrl.clone();
  url.pathname = "/admin/login";
  url.search = "";
  return NextResponse.redirect(url);
}
`
)

// webAdminMiddlewareTS is the web app's middleware when only the admin panel
// needs one. grit add web-auth replaces it with webMiddlewareTS, which carries
// the same gate.
func webAdminMiddlewareTS() string {
	return `import { NextResponse, type NextRequest } from "next/server";

` + adminGateTS + `
export function middleware(request: NextRequest) {
  return adminGate(request) ?? NextResponse.next();
}

export const config = {
  matcher: ["/admin", "/admin/:path*"],
};
`
}

// writeAdminEdgeGuard gives a Next.js web app with the embedded admin its
// middleware, when it has none. An existing middleware.ts is the developer's
// and is never replaced; one without the gate gets a note.
func writeAdminEdgeGuard(root string, opts Options) error {
	if !opts.ShouldEmbedAdmin() || opts.UseTanStack() {
		return nil
	}
	web := filepath.Join(root, "apps", "web")
	if !fileExists(filepath.Join(web, "package.json")) {
		return nil
	}
	for _, name := range []string{"middleware.ts", "proxy.ts", "middleware.js", "proxy.js"} {
		if path := filepath.Join(web, name); fileExists(path) {
			if !fileContains(path, "adminGate(") {
				fmt.Printf("  ⚠ apps/web/%s is your own, so the admin panel's sign-in gate was not added to it.\n"+
					"    Copy adminGate from the Grit docs into it to stop signed-out visitors loading admin pages.\n", name)
			}
			return nil
		}
	}
	return writeFile(filepath.Join(web, "middleware.ts"), webAdminMiddlewareTS())
}

var redisServiceHead = regexp.MustCompile(`(?m)^  redis:\n    image: redis:[^\n]*\n    container_name: [^\n]*\n    restart: unless-stopped\n`)

// repairDeployHardening applies M7, M8 and M10 to an existing project. M6 and
// M9 arrive with files upgrade already delivers whole.
func repairDeployHardening(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	if path := filepath.Join(root, ".dockerignore"); fileExists(path) {
		if err := repairTextFile(root, m, path, repairDockerIgnoreSource); err != nil {
			return err
		}
	}
	if dockerfile := filepath.Join(apiRoot, "Dockerfile"); fileExists(dockerfile) {
		if path := filepath.Join(apiRoot, ".dockerignore"); !fileExists(path) {
			if err := os.WriteFile(path, []byte(apiDockerIgnore), 0o644); err != nil {
				return err
			}
			fmt.Println("  ✓ apps/api/.dockerignore keeps local databases, .env files and binaries out of the API image")
		}
	}

	for _, name := range []string{"docker-compose.prod.yml", "docker-compose.yml"} {
		if path := filepath.Join(root, name); fileExists(path) {
			if err := repairTextFile(root, m, path, repairComposeHardeningSource); err != nil {
				return err
			}
		}
	}
	if prod := filepath.Join(root, "docker-compose.prod.yml"); fileContains(prod, "REDIS_PASSWORD") && !fileContains(filepath.Join(root, ".env"), "REDIS_PASSWORD=") {
		fmt.Println("  ⚠ docker-compose.prod.yml now gives Redis a password. Add one to .env before the next production deploy:\n" +
			"      REDIS_PASSWORD=<the output of: openssl rand -hex 24>")
	}

	if path := filepath.Join(apiRoot, "internal", "config", "config.go"); fileExists(path) {
		if err := repairSourceFile(root, m, path, repairSQLiteProductionSource); err != nil {
			return err
		}
	}
	return writeAdminEdgeGuard(root, opts)
}

func repairDockerIgnoreSource(src string) (string, []string, []string) {
	if strings.Contains(src, "**/node_modules") {
		return src, nil, nil
	}
	if strings.ReplaceAll(src, "\r\n", "\n") != dockerIgnoreOld {
		return src, nil, []string{".dockerignore is not the file Grit wrote: use **/ patterns (**/node_modules, **/.env.*, **/*.db, **/*.exe) so nested .env.local files, databases and binaries stay out of images"}
	}
	return dockerIgnoreNew, []string{".dockerignore keeps nested .env files, databases and binaries out of build contexts"}, nil
}

func repairComposeHardeningSource(src string) (string, []string, []string) {
	out := strings.ReplaceAll(src, pgbouncerImageOld, pgbouncerImageNew)
	out = strings.ReplaceAll(out, minioImageOld, minioImageNew)
	var changes, warnings []string
	if out != src {
		changes = append(changes, "pgbouncer and MinIO run pinned image versions, and MinIO comes from quay.io")
	}
	if start, end, ok := serviceBlock(out, "api"); ok && strings.Contains(out[start:end], composeRedisURLOld) {
		block := out[start:end]
		if !strings.Contains(block, "DB_PROVIDER:") && strings.Count(block, composeDBProviderAnchor) == 1 {
			block = strings.Replace(block, composeDBProviderAnchor, composeDBProvider, 1)
		}
		block = strings.Replace(block, composeRedisURLOld, composeRedisURLNew, 1)
		loc := redisServiceHead.FindStringIndex(out[end:])
		if loc == nil {
			warnings = append(warnings, "the redis service is not the one Grit wrote: start it with --requirepass ${REDIS_PASSWORD} before the API's REDIS_URL carries that password")
			return out, changes, warnings
		}
		redisEnd := end + loc[1]
		out = out[:start] + block + out[end:redisEnd] + composeRedisAuth + out[redisEnd:]
		changes = append(changes, "production uses Postgres, and Redis takes REDIS_PASSWORD")
	}
	return out, changes, warnings
}

func repairSQLiteProductionSource(src string) (string, []string, []string) {
	if strings.Contains(src, "ALLOW_SQLITE_IN_PRODUCTION") || !strings.Contains(src, "cfg.GORMStudioDisableSQL = true") {
		return src, nil, nil
	}
	if strings.Count(src, sqliteProductionAnchor) != 1 {
		return src, nil, []string{"config.go is not the file Grit wrote: refuse to start with APP_ENV=production and a sqlite: DATABASE_URL, or a redeploy loses the database"}
	}
	out := strings.Replace(src, sqliteProductionAnchor, sqliteProductionCheck, 1)
	out, ok := withImports(out, "strings", "fmt")
	if !ok {
		return src, nil, []string{"could not add imports to config.go"}
	}
	return out, []string{"a production API refuses to start on SQLite unless ALLOW_SQLITE_IN_PRODUCTION=true"}, nil
}
