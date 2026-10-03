package scaffold

import (
	"strings"
	"testing"
)

// v3350Routes is routes.go as v3.350.0 wrote it: the four states, and nothing
// serving /llms.txt yet.
func v3350Routes(t *testing.T) string {
	t.Helper()
	fresh := apiRoutesGo()
	old := strings.Replace(fresh, routesDocsBlock+llmsRoutesBlock, routesDocsBlock, 1)
	old = strings.Replace(old, llmsRegisterFunc, "", 1)
	old = strings.Replace(old, "\t\"{{MODULE}}/internal/llms\"\n", "", 1)
	old = strings.Replace(old, "\t\"{{MODULE}}/internal/paginate\"\n", "", 1)
	if old == fresh || strings.Contains(old, "registerLLMSFiles") || strings.Contains(old, "/internal/llms\"") {
		t.Fatal("could not rebuild routes.go from before /llms.txt was mounted")
	}
	return old
}

func TestRepairLLMSFiles(t *testing.T) {
	out, fixed, warn := repairLLMSFilesSource("demo")(v3350Routes(t))
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{
		`"demo/internal/llms"`,
		`"demo/internal/paginate"`,
		"registerLLMSFiles(r, cfg)",
		"func registerLLMSFiles(r *gin.Engine, cfg *config.Config) {",
		`r.GET("/llms.txt"`,
		`r.GET("/llms-full.txt"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired routes.go is missing %q", want)
		}
	}
	// Mounted behind the reference's gate, not beside it: both files describe
	// the whole surface, admin routes included.
	gate := strings.Index(out, "registerAPIDocs(r, db, cfg)")
	mount := strings.Index(out, "registerLLMSFiles(r, cfg)")
	if gate < 0 || mount < gate {
		t.Error("the mount is not after the reference's gate")
	}
	// And below the imports it is the file a fresh scaffold writes, down to the
	// blank lines. An upgraded project that differs from a new one by whitespace
	// is a diff somebody has to read on every later upgrade.
	//
	// The imports themselves are not compared: addImportGroup appends a group
	// rather than sorting a new line into the one it belongs in, which gofmt
	// leaves alone and the compiler does not care about.
	fresh := strings.ReplaceAll(apiRoutesGo(), "{{MODULE}}", "demo")
	if belowImports(t, out) != belowImports(t, fresh) {
		t.Error("the repaired routes.go is not what a fresh scaffold writes")
	}
}

// belowImports is a Go file after its import block, which is the part two
// routes.go files should agree on byte for byte.
func belowImports(t *testing.T, src string) string {
	t.Helper()
	at := strings.Index(src, "\n)\n")
	if at < 0 {
		t.Fatal("no import block in routes.go")
	}
	return src[at:]
}

func TestRepairLLMSFilesIsIdempotent(t *testing.T) {
	first, _, _ := repairLLMSFilesSource("demo")(v3350Routes(t))
	again, fixed, warn := repairLLMSFilesSource("demo")(first)
	if again != first || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("a second upgrade changed routes.go again (fixed %v, warn %v)", fixed, warn)
	}
}

func TestRepairLLMSFilesWaitsForTheFourStates(t *testing.T) {
	// The handler is written beside eventBusStatus, so a routes.go that has not
	// taken the four-state repair yet is left alone and catches up next time
	// rather than having a function dropped at a line number.
	src := strings.Replace(v3350Routes(t), "func eventBusStatus() health.Report {", "func eventBusStatus() interface{} {", 1)
	out, fixed, warn := repairLLMSFilesSource("demo")(src)
	if out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestRepairLLMSFilesSkipsAProjectWithNoReference(t *testing.T) {
	// A project that turned its API reference off has said it does not want its
	// route list served, and these two files are the same claim in another form.
	src := strings.Replace(v3350Routes(t), routesDocsBlock, "", 1)
	out, fixed, warn := repairLLMSFilesSource("demo")(src)
	if out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshRoutesServeTheLLMSFiles(t *testing.T) {
	src := apiRoutesGo()
	for _, want := range []string{
		`"{{MODULE}}/internal/llms"`,
		"registerLLMSFiles(r, cfg)",
		"llms.Index(conf)",
		"llms.Full(conf, table())",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("a fresh routes.go is missing %q", want)
		}
	}
	if out, fixed, warn := repairLLMSFilesSource("demo")(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("a fresh routes.go still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
}

func TestLLMSPackageIsFrameworkOwned(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "demo", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writeFrameworkOwnedFiles: %v", err)
	}
	for _, rel := range []string{"internal/llms/llms.go", "internal/llms/llms_test.go"} {
		if !fileExists(opts.APIRoot(root) + "/" + rel) {
			t.Errorf("%s was not written", rel)
		}
	}
}
