package scaffold

import (
	"fmt"
	"sort"
	"strings"
)

// The NEXT_PUBLIC_* variables a generated Next app reads, in one list.
//
// grit#216: the Dockerfiles accepted `NEXT_PUBLIC_API_URL` and nothing else,
// while the two apps read eleven between them. The other ten could not be
// supplied to a production build at all, so they took their development
// defaults in the deployed image and nothing said so.
//
// One of them breaks every upload, silently.
// `apps/web/next.config.ts` builds the Content-Security-Policy from
// NEXT_PUBLIC_STORAGE_URL, falling back to a MinIO on localhost. Uploads are
// presigned PUTs the browser makes straight to object storage, so a production
// image ships a CSP that permits only the developer's laptop. The API signs a
// good URL, the browser refuses it, and the trace is a console violation with
// no HTTP status and no server log line: nothing to grep for.
//
// # Why a list rather than more ARG lines
//
// Because the drift is the bug. .env.example documents these, the apps read
// them, the Dockerfiles accept them and docker-compose.prod.yml passes them:
// four places that have to agree, and they stopped agreeing without anything
// failing. nextPublicBuildArgs is what the Dockerfile and the compose file are
// both generated from, and TestEveryNextPublicVarIsABuildArg reads the
// generated apps and fails if one is missing from it, so the next variable
// somebody adds cannot repeat this.
//
// # Why build args at all
//
// Next inlines NEXT_PUBLIC_* at build time. Setting one as an `ENV` on the
// running container, or in a platform's environment variables panel, is too
// late: the value is already baked into the bundle. That is the most confusing
// property of deploying this stack, and it is why .env.example telling you to
// set the variable was not enough on its own.

// nextPublicVar is one build-time variable, with how the compose file supplies
// it.
type nextPublicVar struct {
	// Name is the NEXT_PUBLIC_* variable.
	Name string
	// Apps is which generated apps read it: "web", "admin", or both.
	Apps []string
	// Compose is the right-hand side in docker-compose.prod.yml's build.args.
	//
	// Required for the ones with no safe default, so a deploy fails loudly
	// rather than shipping a broken bundle. An empty string is falsy in
	// JavaScript, so `|| "http://localhost:9002"` fires on an empty value
	// exactly as it does on an unset one: a default of "" here would
	// reintroduce the whole bug.
	Compose string
}

var nextPublicBuildArgs = []nextPublicVar{
	{
		Name:    "NEXT_PUBLIC_API_URL",
		Apps:    []string{"web", "admin", "docs"},
		Compose: "${API_URL:-http://localhost:8080}",
	},
	{
		Name: "NEXT_PUBLIC_STORAGE_URL",
		Apps: []string{"web", "admin", "docs"},
		// Required. Without it the CSP permits only localhost and every
		// upload and every stored image is blocked, with no server-side
		// trace (grit#216).
		Compose: "${NEXT_PUBLIC_STORAGE_URL:?set this to the browser-facing origin files are served from, or the CSP blocks every upload}",
	},
	{
		Name:    "NEXT_PUBLIC_STORAGE_UPLOAD_URL",
		Apps:    []string{"web", "admin", "docs"},
		Compose: "${NEXT_PUBLIC_STORAGE_UPLOAD_URL:-}",
	},
	{
		Name:    "NEXT_PUBLIC_IMAGE_ORIGINS",
		Apps:    []string{"web", "admin", "docs"},
		Compose: "${NEXT_PUBLIC_IMAGE_ORIGINS:-}",
	},
	{
		Name:    "NEXT_PUBLIC_THEME",
		Apps:    []string{"web", "admin", "docs"},
		Compose: "${THEME:-atlas}",
	},
	{
		Name:    "NEXT_PUBLIC_SOCIAL_AUTH_ENABLED",
		Apps:    []string{"web", "admin"},
		Compose: "${NEXT_PUBLIC_SOCIAL_AUTH_ENABLED:-}",
	},
	{
		// The marketing navbar, the nav menu and the dev-links component all
		// link to it, so without it a deployed public site has links pointing
		// at the visitor's own machine.
		Name:    "NEXT_PUBLIC_ADMIN_URL",
		Apps:    []string{"web"},
		Compose: "${NEXT_PUBLIC_ADMIN_URL:-}",
	},
	{
		Name:    "NEXT_PUBLIC_WEB_URL",
		Apps:    []string{"admin"},
		Compose: "${NEXT_PUBLIC_WEB_URL:-}",
	},
	{
		Name:    "NEXT_PUBLIC_IP_HINT",
		Apps:    []string{"admin"},
		Compose: "${NEXT_PUBLIC_IP_HINT:-}",
	},
	{
		Name:    "NEXT_PUBLIC_SESSION_IDLE_MS",
		Apps:    []string{"admin"},
		Compose: "${NEXT_PUBLIC_SESSION_IDLE_MS:-}",
	},
	{
		Name:    "NEXT_PUBLIC_SESSION_COUNTDOWN_MS",
		Apps:    []string{"admin"},
		Compose: "${NEXT_PUBLIC_SESSION_COUNTDOWN_MS:-}",
	},
}

// readsVar reports whether an app reads this variable.
func (v nextPublicVar) readsVar(app string) bool {
	for _, a := range v.Apps {
		if a == app {
			return true
		}
	}
	return false
}

// nextPublicVarsFor returns the variables one app reads, sorted by name so the
// generated file is stable.
func nextPublicVarsFor(app string) []nextPublicVar {
	var out []nextPublicVar
	for _, v := range nextPublicBuildArgs {
		if v.readsVar(app) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// nextPublicDockerfileArgs renders the ARG/ENV pairs for one app's builder
// stage.
func nextPublicDockerfileArgs(app string) string {
	var b strings.Builder
	b.WriteString("# Next inlines NEXT_PUBLIC_* at BUILD time, so each one has to arrive as a\n")
	b.WriteString("# build arg. Setting one in the platform's environment panel, or as an ENV\n")
	b.WriteString("# on the running container, is too late: the value is already in the bundle.\n")
	b.WriteString("# Changing one needs a rebuild, not a restart.\n")
	b.WriteString("#\n")
	b.WriteString("# Generated from one list. These used to be a single ARG while the app read\n")
	b.WriteString("# eleven, and the missing NEXT_PUBLIC_STORAGE_URL left a CSP that blocked\n")
	b.WriteString("# every upload with nothing in any log (grit#216).\n")
	for _, v := range nextPublicVarsFor(app) {
		b.WriteString(fmt.Sprintf("ARG %s\nENV %s=$%s\n", v.Name, v.Name, v.Name))
	}
	return b.String()
}

// nextPublicComposeArgs renders one app's build.args block for
// docker-compose.prod.yml, at the indentation that block sits at.
func nextPublicComposeArgs(app string) string {
	var lines []string
	for _, v := range nextPublicVarsFor(app) {
		lines = append(lines, fmt.Sprintf("        %s: %s", v.Name, v.Compose))
	}
	return strings.Join(lines, "\n")
}

// NextPublicBuildArgNames is the list of variables, for grit doctor and for
// anything that reports on them.
func NextPublicBuildArgNames(app string) []string {
	vars := nextPublicVarsFor(app)
	out := make([]string, 0, len(vars))
	for _, v := range vars {
		out = append(out, v.Name)
	}
	return out
}
