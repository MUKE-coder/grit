package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var nextPublicRe = regexp.MustCompile(`NEXT_PUBLIC_[A-Z0-9_]+`)

// The variables a generated app actually reads, found by scaffolding one and
// grepping it.
//
// This is the half that makes the test worth having. A test over
// nextPublicBuildArgs alone would assert that the list equals itself. The
// question is whether the list still matches the apps, and the only way to
// answer it is to read the apps (grit#216).
func varsReadBy(t *testing.T, root, app string) []string {
	t.Helper()
	dir := filepath.Join(root, "apps", app)
	found := map[string]bool{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "node_modules" || info.Name() == ".next" {
				return filepath.SkipDir
			}
			return nil
		}
		switch filepath.Ext(path) {
		case ".ts", ".tsx", ".js", ".mjs", ".cjs":
		default:
			return nil
		}
		// The Dockerfile and .env.example name them too; this is about what
		// the code reads.
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		for _, name := range nextPublicRe.FindAllString(string(raw), -1) {
			found[name] = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking apps/%s: %v", app, err)
	}

	out := make([]string, 0, len(found))
	for name := range found {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// scaffoldForBuildArgs writes a triple project with both Next apps.
//
// A real scaffold rather than rendering the two templates in isolation,
// because the question this file asks is whether the Dockerfile matches the
// app, and rendering only the Dockerfile would answer a different one.
func scaffoldForBuildArgs(t *testing.T) string {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	if err := Run(Options{
		ProjectName:  "buildargs",
		Architecture: ArchTriple,
		Frontend:     "next",
		Version:      DefaultVersion,
	}); err != nil {
		t.Fatalf("scaffolding: %v", err)
	}
	return filepath.Join(dir, "buildargs")
}

// Every NEXT_PUBLIC_* the generated apps read has to be a build arg, because
// Next inlines them at build time and there is no second chance at runtime.
//
// The one that cost this: NEXT_PUBLIC_STORAGE_URL. apps/web/next.config.ts
// builds the Content-Security-Policy from it and falls back to a MinIO on
// localhost, so a production image permitted only the developer's laptop as a
// storage origin. Uploads are presigned PUTs straight to object storage, so
// every upload and every stored image was blocked, with no HTTP status and no
// server log line to grep for.
func TestEveryNextPublicVarIsABuildArg(t *testing.T) {
	root := scaffoldForBuildArgs(t)

	for _, app := range []string{"web", "admin"} {
		read := varsReadBy(t, root, app)
		if len(read) == 0 {
			t.Fatalf("apps/%s reads no NEXT_PUBLIC_ variables, so this test is checking nothing", app)
		}

		declared := map[string]bool{}
		for _, name := range NextPublicBuildArgNames(app) {
			declared[name] = true
		}
		for _, name := range read {
			if !declared[name] {
				t.Errorf("apps/%s reads %s and nextPublicBuildArgs does not list it for that app, "+
					"so a production build takes its development default and nothing says so", app, name)
			}
		}

		// And the Dockerfile actually carries them, which is the other side of
		// the pair: a list nothing renders is a list that is still wrong.
		raw, err := os.ReadFile(filepath.Join(root, "apps", app, "Dockerfile"))
		if err != nil {
			t.Fatal(err)
		}
		dockerfile := string(raw)
		for _, name := range read {
			if !strings.Contains(dockerfile, "ARG "+name) {
				t.Errorf("apps/%s/Dockerfile has no ARG %s", app, name)
			}
			if !strings.Contains(dockerfile, "ENV "+name+"=$"+name) {
				t.Errorf("apps/%s/Dockerfile declares ARG %s without the matching ENV, so the build cannot see it", app, name)
			}
		}
	}
}

// The compose file has to pass them, or the ARGs are accepted and never
// supplied.
func TestComposePassesEveryBuildArg(t *testing.T) {
	root := scaffoldForBuildArgs(t)
	raw, err := os.ReadFile(filepath.Join(root, "docker-compose.prod.yml"))
	if err != nil {
		t.Fatal(err)
	}
	compose := string(raw)

	for _, app := range []string{"web", "admin"} {
		for _, name := range NextPublicBuildArgNames(app) {
			if !strings.Contains(compose, name+":") {
				t.Errorf("docker-compose.prod.yml never passes %s, which apps/%s reads", name, app)
			}
		}
	}

	// The one with no safe default is required, so a deploy fails loudly
	// rather than shipping a bundle whose CSP blocks every upload. An empty
	// string is falsy in JavaScript, so a default of "" would fire the
	// localhost fallback exactly as an unset value does and reintroduce the
	// whole bug.
	if !strings.Contains(compose, "${NEXT_PUBLIC_STORAGE_URL:?") {
		t.Error("NEXT_PUBLIC_STORAGE_URL has a default or is optional; it must be required, " +
			"because an empty value silently restores the localhost CSP")
	}
}

// The workspace packages' own node_modules reach the builder stage.
//
// hoisted puts one copy of each dependency at the root, except where two
// packages need different versions, and then it nests one. A packages/* member
// whose dependency pins a peer more tightly than the root does gets its own
// node_modules, and without it the build fails with forty lines of structural
// type mismatch naming two paths that are the same library at two versions:
// it reads as a bug in code you did not touch (grit#216).
func TestTheBuilderStageCopiesTheWorkspacePackages(t *testing.T) {
	root := scaffoldForBuildArgs(t)

	for _, app := range []string{"web", "admin"} {
		raw, err := os.ReadFile(filepath.Join(root, "apps", app, "Dockerfile"))
		if err != nil {
			t.Fatal(err)
		}
		dockerfile := string(raw)

		const copyPackages = "COPY --from=deps /app/packages ./packages"
		if !strings.Contains(dockerfile, copyPackages) {
			t.Errorf("apps/%s/Dockerfile does not copy the workspace packages from deps, "+
				"so a package with a nested dependency fails to build", app)
			continue
		}

		// After `COPY . .`, or the context copy puts a source tree with no
		// node_modules back over it. The order is the whole fix.
		context := strings.Index(dockerfile, "\nCOPY . .\n")
		packages := strings.Index(dockerfile, copyPackages)
		if context < 0 {
			t.Errorf("apps/%s/Dockerfile has no COPY . . to order against", app)
			continue
		}
		if packages < context {
			t.Errorf("apps/%s/Dockerfile copies the packages before COPY . . , which overwrites them", app)
		}

		// And never a node_modules inside packages/, which pnpm does not
		// create for a package with no nested dependency: COPY --from has no
		// "if it exists" and the build stops (grit#177).
		if strings.Contains(dockerfile, "/app/packages/shared/node_modules") {
			t.Errorf("apps/%s/Dockerfile names packages/shared/node_modules, which pnpm often does not create", app)
		}
	}
}

// dockerInstructions are the keywords a line may start with. Anything else at
// the start of a line is a spliced fragment.
var dockerInstructions = []string{
	"FROM", "RUN", "CMD", "LABEL", "EXPOSE", "ENV", "ADD", "COPY", "ENTRYPOINT",
	"VOLUME", "USER", "WORKDIR", "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL",
}

// The generated Dockerfile has to be a Dockerfile.
//
// This exists because the first version of the fix above passed every test in
// this file while producing a file that could not build. The ARG block went in
// at the wrong %s, landing in the middle of
//
//	COPY --from=deps /app/apps/web/node_modules ./apps/
//
// and leaving the rest of that path, and a bare `web`, on lines of their own.
// Every Contains check still matched, because the strings were all present:
// they were simply in pieces.
//
// A format string with twelve %s and one of them a different value is exactly
// the kind of thing to get wrong once, so the check is structural rather than
// another substring.
func TestTheGeneratedDockerfileIsStructurallyADockerfile(t *testing.T) {
	root := scaffoldForBuildArgs(t)

	for _, app := range []string{"web", "admin"} {
		raw, err := os.ReadFile(filepath.Join(root, "apps", app, "Dockerfile"))
		if err != nil {
			t.Fatal(err)
		}
		body := strings.ReplaceAll(string(raw), "\r\n", "\n")

		// Go's own marker for a format string that did not line up.
		if strings.Contains(body, "%!") {
			t.Errorf("apps/%s/Dockerfile holds a format error marker", app)
		}

		continued := false
		for i, line := range strings.Split(body, "\n") {
			trimmed := strings.TrimSpace(line)
			wasContinued := continued
			continued = strings.HasSuffix(trimmed, "\\")

			if trimmed == "" || strings.HasPrefix(trimmed, "#") || wasContinued {
				continue
			}
			known := false
			for _, kw := range dockerInstructions {
				if strings.HasPrefix(trimmed, kw+" ") || strings.HasPrefix(trimmed, kw+"\t") {
					known = true
					break
				}
			}
			if !known {
				t.Errorf("apps/%s/Dockerfile line %d is not an instruction, a comment or a continuation: %q",
					app, i+1, line)
			}
		}

		// And the two lines the splice broke, whole.
		for _, want := range []string{
			"COPY --from=deps /app/apps/" + app + "/node_modules ./apps/" + app + "/node_modules",
			"RUN pnpm --filter " + app + " build",
		} {
			if !strings.Contains(body, want) {
				t.Errorf("apps/%s/Dockerfile is missing the whole line %q", app, want)
			}
		}
	}
}
