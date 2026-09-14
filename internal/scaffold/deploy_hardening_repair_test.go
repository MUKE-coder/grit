package scaffold

import (
	"go/format"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestSAMLMetadataGoesThroughSafefetch(t *testing.T) {
	src := apiSAMLServiceGo()
	if strings.Contains(src, samlFetchOld) || !strings.Contains(src, samlFetchNew) {
		t.Error("SAML metadata is still fetched with http.DefaultClient")
	}
	// Parsing does not catch an unused import; the generated API then fails to
	// build. http.DefaultClient was the file's only use of net/http.
	if strings.Contains(src, "\"net/http\"") && !strings.Contains(src, "http.") {
		t.Error("the SAML service imports net/http without using it")
	}
	if _, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("the SAML service template is not valid Go: %v", err)
	}
}

func TestDockerIgnoreIsRecursive(t *testing.T) {
	if dockerIgnore() != dockerIgnoreNew {
		t.Error("the root .dockerignore template is not the recursive one")
	}
	out, changes, warnings := repairDockerIgnoreSource(dockerIgnoreOld)
	if out != dockerIgnoreNew || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("the old .dockerignore was not replaced: %v %v", changes, warnings)
	}
	if _, _, warnings := repairDockerIgnoreSource("node_modules\nmy-own-thing\n"); len(warnings) != 1 {
		t.Error("an edited .dockerignore should be left alone with a note")
	}
	for _, want := range []string{"**/.env.*", "**/*.db", "**/*.exe", "**/node_modules"} {
		if !strings.Contains(dockerIgnoreNew, want) {
			t.Errorf("the root .dockerignore misses %s", want)
		}
	}
}

func TestProductionComposeHardening(t *testing.T) {
	fresh := dockerComposeProd(Options{ProjectName: "app", Architecture: ArchDouble})
	for _, want := range []string{"DB_PROVIDER: postgres", composeRedisURLNew, "--requirepass", pgbouncerImageNew, minioImageNew} {
		if !strings.Contains(fresh, want) {
			t.Errorf("the production compose template is missing %q", want)
		}
	}
	if out, changes, _ := repairComposeHardeningSource(fresh); out != fresh || len(changes) != 0 {
		t.Errorf("a fresh production compose needs no repair: %v", changes)
	}

	old := strings.Replace(fresh, composeDBProvider, composeDBProviderAnchor, 1)
	old = strings.Replace(old, composeRedisURLNew, composeRedisURLOld, 1)
	old = strings.Replace(old, composeRedisAuth, "", 1)
	old = strings.ReplaceAll(old, pgbouncerImageNew, pgbouncerImageOld)
	old = strings.ReplaceAll(old, minioImageNew, minioImageOld)
	out, changes, warnings := repairComposeHardeningSource(old)
	if len(warnings) != 0 || out != fresh {
		t.Errorf("repairing the previous compose file did not give the fresh one (changes %v, warnings %v)", changes, warnings)
	}
}

func TestProductionRefusesSQLite(t *testing.T) {
	files, _ := os.ReadDir(".")
	found := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_test.go") || f.Name() == "deploy_hardening_repair.go" || !strings.HasSuffix(f.Name(), ".go") {
			continue
		}
		raw, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(raw), "sqliteProductionCheck") {
			found = true
		}
		// The anchor must not survive anywhere a template could still emit it.
		if strings.Contains(string(raw), "\t\tcfg.GORMStudioDisableSQL = true\n\t}\n") {
			t.Errorf("%s still emits the production block without the SQLite check", f.Name())
		}
	}
	if !found {
		t.Fatal("no template splices in the SQLite production check")
	}

	src := "package config\n\nimport (\n\t\"fmt\"\n\t\"strings\"\n)\n\nfunc load(cfg *Config) (*Config, error) {\n\tif cfg.AppEnv == \"production\" {\n" + sqliteProductionAnchor + "\treturn cfg, nil\n}\n"
	out, changes, warnings := repairSQLiteProductionSource(src)
	if len(warnings) != 0 || len(changes) != 1 || !strings.Contains(out, "ALLOW_SQLITE_IN_PRODUCTION") {
		t.Fatalf("changes %v warnings %v", changes, warnings)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Errorf("not valid Go: %v", err)
	}
}

func TestWorkflowsArePinnedAndScoped(t *testing.T) {
	unpinned := regexp.MustCompile(`uses: [\w./-]+@v\d+\b(?:\s|$)`)
	for _, name := range []string{"api_ci_files.go", "api_lint_files.go", "api_release_workflow.go"} {
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		src := string(raw)
		for _, line := range strings.Split(src, "\n") {
			if strings.Contains(line, "uses:") && !regexp.MustCompile(`@[0-9a-f]{40}`).MatchString(line) && unpinned.MatchString(line) {
				t.Errorf("%s: action not pinned to a commit: %s", name, strings.TrimSpace(line))
			}
		}
		for _, banned := range []string{"@latest", "version: latest", "get.anchore.io"} {
			if strings.Contains(src, banned) {
				t.Errorf("%s still has %q", name, banned)
			}
		}
	}
	raw, _ := os.ReadFile("api_release_workflow.go")
	if strings.Contains(string(raw), "permissions:\n  contents: write\n") {
		t.Error("release.yml still grants contents: write to the whole workflow")
	}
	raw, _ = os.ReadFile("api_lint_files.go")
	if !strings.Contains(string(raw), "permissions:\n  contents: read\n") {
		t.Error("lint.yml has no read-only permissions block")
	}
}

func TestAdminSignInGate(t *testing.T) {
	src := apiAuthServiceGo()
	if !strings.Contains(src, signedInCookieSet) || !strings.Contains(src, signedInCookieClear) {
		t.Error("the API does not set and clear grit_signed_in")
	}
	for name, mw := range map[string]string{"web-auth middleware": webMiddlewareTS(), "admin middleware": webAdminMiddlewareTS()} {
		if !strings.Contains(mw, "adminGate(request)") || !strings.Contains(mw, `"/admin/:path*"`) {
			t.Errorf("the %s does not gate /admin", name)
		}
	}
	if !strings.Contains(adminLayoutComponent(), "permissions.length > 0") {
		t.Error("a USER with no grants can still open admin pages other than the profile")
	}
}
