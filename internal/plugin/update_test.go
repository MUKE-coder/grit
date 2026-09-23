package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// updatable builds a project with one plugin installed, and returns the
// context plus the plugin as it was at install time.
func updatable(t *testing.T) (Context, Plugin) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "apps", "api", "internal"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "grit.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	// The file an injection goes into.
	routes := "package routes\n\nfunc Setup() {\n\t// grit:handlers\n}\n"
	if err := os.WriteFile(filepath.Join(root, "apps", "api", "internal", "routes.go"), []byte(routes), 0644); err != nil {
		t.Fatal(err)
	}

	p := Plugin{
		Name:    "demo",
		Version: "1.0.0",
		Files: func(ctx Context) map[string]string {
			return map[string]string{
				"apps/api/internal/one.go": "package internal\n\nconst One = 1\n",
				"apps/api/internal/two.go": "package internal\n\nconst Two = 2\n",
			}
		},
		Injections: func(ctx Context) []Injection {
			return []Injection{{
				File:   "apps/api/internal/routes.go",
				Marker: "// grit:handlers",
				Code:   "\tdemoHandler := handlers.NewDemo()",
			}}
		},
	}

	ctx := Context{Root: root, APIRoot: filepath.Join(root, "apps", "api")}
	if _, err := Install(ctx, p); err != nil {
		t.Fatal(err)
	}
	return ctx, p
}

// The whole reason this command exists: a plugin gains a file, and projects
// that installed it before must be able to have it.
func TestUpdateBringsANewFileAndRewritesUntouchedOnes(t *testing.T) {
	ctx, p := updatable(t)

	p.Version = "1.1.0"
	p.Files = func(ctx Context) map[string]string {
		return map[string]string{
			"apps/api/internal/one.go":   "package internal\n\nconst One = 11\n",
			"apps/api/internal/two.go":   "package internal\n\nconst Two = 2\n",
			"apps/api/internal/three.go": "package internal\n\nconst Three = 3\n",
		}
	}

	res, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Added) != 1 || res.Added[0] != "apps/api/internal/three.go" {
		t.Errorf("added = %v, want the new file", res.Added)
	}
	if len(res.Replaced) != 1 || res.Replaced[0] != "apps/api/internal/one.go" {
		t.Errorf("replaced = %v, want the changed file", res.Replaced)
	}
	if len(res.Edited) != 0 || len(res.Unverified) != 0 {
		t.Errorf("nothing was edited, but got %v %v", res.Edited, res.Unverified)
	}

	got, err := os.ReadFile(filepath.Join(ctx.Root, "apps/api/internal/one.go"))
	if err != nil || !strings.Contains(string(got), "One = 11") {
		t.Errorf("one.go was not updated: %q %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(ctx.Root, "apps/api/internal/three.go")); err != nil {
		t.Errorf("the new file did not arrive: %v", err)
	}

	// The lockfile knows about the new file, so removal takes it away too.
	lock, err := LoadLock(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	entry, _ := lock.Find("demo")
	if entry.Version != "1.1.0" {
		t.Errorf("lock version = %s", entry.Version)
	}
	var found bool
	for _, f := range entry.Files {
		if f == "apps/api/internal/three.go" {
			found = true
		}
	}
	if !found {
		t.Error("the new file is not in the lockfile, so removing the plugin would leave it behind")
	}

	// Running it again changes nothing.
	again, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if again.Changed() {
		t.Errorf("a second update did something: %+v", again)
	}
}

// The rule that makes this safe to run: your edits are never overwritten.
func TestUpdateKeepsAFileYouEdited(t *testing.T) {
	ctx, p := updatable(t)

	mine := "package internal\n\nconst One = 1 // mine\n"
	if err := os.WriteFile(filepath.Join(ctx.Root, "apps/api/internal/one.go"), []byte(mine), 0644); err != nil {
		t.Fatal(err)
	}

	p.Version = "1.1.0"
	p.Files = func(ctx Context) map[string]string {
		return map[string]string{
			"apps/api/internal/one.go": "package internal\n\nconst One = 11\n",
			"apps/api/internal/two.go": "package internal\n\nconst Two = 2\n",
		}
	}

	res, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Edited) != 1 || res.Edited[0] != "apps/api/internal/one.go" {
		t.Errorf("edited = %v, want the file the user changed", res.Edited)
	}
	got, _ := os.ReadFile(filepath.Join(ctx.Root, "apps/api/internal/one.go"))
	if string(got) != mine {
		t.Error("an edited file was overwritten, which is the one thing this must never do")
	}
}

// A project set up before fingerprints existed cannot be vouched for, so its
// differing files are reported rather than rewritten.
func TestUpdateWillNotGuessWithoutFingerprints(t *testing.T) {
	ctx, p := updatable(t)

	lock, err := LoadLock(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	entry, _ := lock.Find("demo")
	entry.FileHashes = nil
	if err := lock.Save(ctx.Root); err != nil {
		t.Fatal(err)
	}

	p.Version = "1.1.0"
	p.Files = func(ctx Context) map[string]string {
		return map[string]string{
			"apps/api/internal/one.go": "package internal\n\nconst One = 11\n",
			"apps/api/internal/two.go": "package internal\n\nconst Two = 2\n",
		}
	}

	res, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unverified) != 1 || res.Unverified[0] != "apps/api/internal/one.go" {
		t.Errorf("unverified = %v, want the file it cannot vouch for", res.Unverified)
	}
	if len(res.Replaced) != 0 {
		t.Errorf("replaced %v without knowing what it wrote", res.Replaced)
	}

	// two.go is identical to the template, so it can be fingerprinted now and
	// updated properly next time.
	lock, _ = LoadLock(ctx.Root)
	entry, _ = lock.Find("demo")
	if entry.FileHashes["apps/api/internal/two.go"] == "" {
		t.Error("an unchanged file was not fingerprinted, so it stays unverifiable forever")
	}
}

// A new patch is applied; the one already there is not applied twice and not
// moved, because moving an injection is how a project stops compiling.
func TestUpdateAppliesOnlyNewInjections(t *testing.T) {
	ctx, p := updatable(t)

	routes := filepath.Join(ctx.Root, "apps/api/internal/routes.go")
	before, err := os.ReadFile(routes)
	if err != nil {
		t.Fatal(err)
	}

	p.Version = "1.1.0"
	p.Injections = func(ctx Context) []Injection {
		return []Injection{
			{File: "apps/api/internal/routes.go", Marker: "// grit:handlers", Code: "\tdemoHandler := handlers.NewDemo()"},
			{File: "apps/api/internal/routes.go", Marker: "// grit:handlers", Code: "\tinvoiceHandler := handlers.NewInvoice()"},
		}
	}

	res, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Injected) != 1 {
		t.Errorf("injected = %v, want only the new patch", res.Injected)
	}

	after, err := os.ReadFile(routes)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(after), "demoHandler := handlers.NewDemo()") != 1 {
		t.Error("the patch already in place was applied a second time")
	}
	if !strings.Contains(string(after), "invoiceHandler := handlers.NewInvoice()") {
		t.Error("the new patch was not applied")
	}
	if len(before) >= len(after) {
		t.Error("the file did not grow, so nothing was inserted")
	}
}

// Updating something that was never installed is a mistake worth naming.
func TestUpdateRefusesAPluginThatIsNotInstalled(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "grit.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Update(Context{Root: root}, Plugin{Name: "ghost", Version: "1.0.0"})
	if err == nil || !strings.Contains(err.Error(), "not installed") {
		t.Errorf("err = %v, want it to say the plugin is not installed", err)
	}
}

// Being careful is not always enough: a plugin that adds a file usually
// changes another to use it, and a project left half-updated does not compile.
// --overwrite is the way out, and it has to take both kinds of left-behind file.
func TestUpdateOverwriteTakesEverything(t *testing.T) {
	ctx, p := updatable(t)

	mine := "package internal\n\nconst One = 1 // mine\n"
	if err := os.WriteFile(filepath.Join(ctx.Root, "apps/api/internal/one.go"), []byte(mine), 0644); err != nil {
		t.Fatal(err)
	}
	// And a second file this project cannot vouch for, as an old lockfile.
	lock, err := LoadLock(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	entry, _ := lock.Find("demo")
	delete(entry.FileHashes, "apps/api/internal/two.go")
	if err := lock.Save(ctx.Root); err != nil {
		t.Fatal(err)
	}

	p.Version = "1.1.0"
	p.Files = func(ctx Context) map[string]string {
		return map[string]string{
			"apps/api/internal/one.go": "package internal\n\nconst One = 11\n",
			"apps/api/internal/two.go": "package internal\n\nconst Two = 22\n",
		}
	}

	res, err := UpdateWith(ctx, p, UpdateOptions{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Edited) != 0 || len(res.Unverified) != 0 {
		t.Errorf("--overwrite left something behind: %v %v", res.Edited, res.Unverified)
	}
	if len(res.Replaced) != 2 {
		t.Errorf("replaced = %v, want both files", res.Replaced)
	}
	for file, want := range map[string]string{
		"apps/api/internal/one.go": "One = 11",
		"apps/api/internal/two.go": "Two = 22",
	} {
		got, err := os.ReadFile(filepath.Join(ctx.Root, file))
		if err != nil || !strings.Contains(string(got), want) {
			t.Errorf("%s = %q, want %s", file, got, want)
		}
	}
}

// A file that differs only in its line endings is not an edit. Git rewrites
// them on checkout on Windows, and counting that as an edit would make an
// update refuse to touch a project nobody has touched.
func TestUpdateIgnoresLineEndings(t *testing.T) {
	ctx, p := updatable(t)

	crlf := strings.ReplaceAll("package internal\n\nconst One = 1\n", "\n", "\r\n")
	if err := os.WriteFile(filepath.Join(ctx.Root, "apps/api/internal/one.go"), []byte(crlf), 0644); err != nil {
		t.Fatal(err)
	}

	res, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Edited) != 0 || len(res.Unverified) != 0 {
		t.Errorf("a line ending was read as an edit: %v %v", res.Edited, res.Unverified)
	}
}

// Half an update is worse than none. A project too old to vouch for its files
// is left entirely alone, rather than gaining a new file whose companion
// change was skipped: that combination does not compile, which is exactly what
// happened to a shop the first time this ran.
func TestUpdateDoesNothingRatherThanHalfSomething(t *testing.T) {
	ctx, p := updatable(t)

	lock, err := LoadLock(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	entry, _ := lock.Find("demo")
	entry.FileHashes = nil
	if err := lock.Save(ctx.Root); err != nil {
		t.Fatal(err)
	}

	p.Version = "1.1.0"
	p.Files = func(ctx Context) map[string]string {
		return map[string]string{
			"apps/api/internal/one.go":   "package internal\n\nconst One = 11\n",
			"apps/api/internal/two.go":   "package internal\n\nconst Two = 2\n",
			"apps/api/internal/three.go": "package internal\n\nconst Three = 3\n",
		}
	}

	res, err := Update(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatal("an unverifiable project was updated anyway")
	}
	if len(res.Added) != 0 || len(res.Replaced) != 0 {
		t.Errorf("something was written: added %v, replaced %v", res.Added, res.Replaced)
	}
	if _, err := os.Stat(filepath.Join(ctx.Root, "apps/api/internal/three.go")); err == nil {
		t.Error("the new file was written even though its companion change could not be")
	}

	// And --overwrite gets the whole update, in one step.
	res, err = UpdateWith(ctx, p, UpdateOptions{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Blocked || len(res.Added) != 1 || len(res.Replaced) != 1 {
		t.Errorf("--overwrite did not complete the update: %+v", res)
	}
}
