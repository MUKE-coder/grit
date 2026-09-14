package scaffold

import (
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type ciShape struct {
	name string
	opts Options
}

func ciShapes() []ciShape {
	return []ciShape{
		{"single", Options{ProjectName: "demo", Architecture: ArchSingle}},
		{"double", Options{ProjectName: "demo", Architecture: ArchDouble}},
		{"triple", Options{ProjectName: "demo", Architecture: ArchTriple}},
		{"api", Options{ProjectName: "demo", Architecture: ArchAPI}},
		{"mobile", Options{ProjectName: "demo", Architecture: ArchMobile}},
		{"double+desktop", Options{ProjectName: "demo", Architecture: ArchDouble, IncludeDesktop: true}},
	}
}

func parseYAML(t *testing.T, name, src string) map[string]interface{} {
	t.Helper()
	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(src), &doc); err != nil {
		t.Fatalf("%s is not valid YAML: %v\n%s", name, err, src)
	}
	// Grit's placeholders, not GitHub's ${{ expressions }}.
	if loc := gritPlaceholderRe.FindString(src); loc != "" {
		t.Errorf("%s has an unreplaced placeholder %s", name, loc)
	}
	return doc
}

var gritPlaceholderRe = regexp.MustCompile(`\{\{[A-Z_]+\}\}`)

func jobsOf(doc map[string]interface{}) map[string]interface{} {
	jobs, _ := doc["jobs"].(map[string]interface{})
	return jobs
}

func TestCIFilesFollowTheLayout(t *testing.T) {
	for _, shape := range ciShapes() {
		l := ciLayoutFor(shape.opts)
		dependabot := dependabotYAML(shape.opts)
		parseYAML(t, shape.name+" dependabot.yml", dependabot)
		wantGo := `directory: "/apps/api"`
		if l.single {
			wantGo = `directory: "/"`
		}
		if !strings.Contains(dependabot, "package-ecosystem: gomod\n    "+wantGo) {
			t.Errorf("%s: Dependabot does not watch the Go module at %s:\n%s", shape.name, wantGo, dependabot)
		}
		if strings.Contains(dependabot, "package-ecosystem: npm") != l.hasFrontend {
			t.Errorf("%s: npm updates present = %v, want %v", shape.name, !l.hasFrontend, l.hasFrontend)
		}

		for file, src := range map[string]string{
			"security.yml": securityCIYAML(shape.opts),
			"ci.yml":       ciYAML(shape.opts),
			"release.yml":  releaseCIYAML(shape.opts),
			"lint.yml":     lintCIYAML(shape.opts),
		} {
			doc := parseYAML(t, shape.name+" "+file, src)
			if !l.single && strings.Contains(src, "frontend") {
				t.Errorf("%s %s mentions frontend/, which a monorepo does not have", shape.name, file)
			}
			if l.single && strings.Contains(src, "apps/api") {
				t.Errorf("%s %s mentions apps/api, which a single app does not have", shape.name, file)
			}
			if strings.Contains(src, "go-version-file") && !strings.Contains(src, "go-version-file: "+l.apiDir+"/go.mod") {
				t.Errorf("%s %s reads Go's version from somewhere other than %s/go.mod", shape.name, file, l.apiDir)
			}
			for name, job := range jobsOf(doc) {
				if spec, ok := job.(map[string]interface{}); ok {
					if cond, ok := spec["if"].(string); ok && strings.Contains(cond, "hashFiles") {
						t.Errorf("%s %s job %s uses hashFiles in a job-level if, which GitHub rejects", shape.name, file, name)
					}
				}
			}
		}

		ciJobs := jobsOf(parseYAML(t, shape.name+" ci.yml", ciYAML(shape.opts)))
		if _, ok := ciJobs["web"]; ok != l.hasFrontend {
			t.Errorf("%s: ci.yml web job present = %v, want %v", shape.name, ok, l.hasFrontend)
		}
		securityJobs := jobsOf(parseYAML(t, shape.name+" security.yml", securityCIYAML(shape.opts)))
		if _, ok := securityJobs["pnpm-audit"]; ok != l.hasFrontend {
			t.Errorf("%s: pnpm audit job present = %v, want %v", shape.name, ok, l.hasFrontend)
		}
		releaseJobs := jobsOf(parseYAML(t, shape.name+" release.yml", releaseCIYAML(shape.opts)))
		if _, ok := releaseJobs["desktop"]; ok != shape.opts.ShouldIncludeDesktop() {
			t.Errorf("%s: release desktop job present = %v, want %v", shape.name, ok, shape.opts.ShouldIncludeDesktop())
		}
		for _, src := range []string{ciYAML(shape.opts), securityCIYAML(shape.opts), lintCIYAML(shape.opts)} {
			if strings.Contains(src, "Placeholder for the embedded frontend") != l.single {
				t.Errorf("%s: the embedded-frontend placeholder should be there only for a single app", shape.name)
			}
		}
		if l.single && !strings.Contains(releaseCIYAML(shape.opts), "Build the frontend") {
			t.Errorf("%s: the release does not build the frontend it embeds", shape.name)
		}
	}
}

func TestTurboHasATestTask(t *testing.T) {
	if !strings.Contains(turboJSON(), `"test": {`) {
		t.Error("turbo.json has no test task, so pnpm test fails")
	}
}
