package prompt

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// enterThrough answers every question with Enter and returns the Options a
// real `grit new` would go on to scaffold from.
//
// The form is driven as the bubbletea model it is, and applyAnswers is the
// same function RunNewProjectPrompt calls, so this is the whole path minus the
// terminal. A test that read the option lists instead would be asserting that
// a list contains a string, which is not the question.
func enterThrough(t *testing.T) scaffold.Options {
	t.Helper()

	opts := &scaffold.Options{}
	arch, frontend, theme, db, storage := "", "", "", "", ""
	preselect(&arch, &frontend, &theme, &db, &storage)

	form := newProjectForm(opts, &arch, &frontend, &theme, &db, &storage).
		WithTheme(ui.PromptTheme())
	form.Update(tea.WindowSizeMsg{Width: testCols, Height: testRows})
	form.Init()

	enter := tea.KeyMsg{Type: tea.KeyEnter}
	for i := 0; i < 20 && form.State == huh.StateNormal; i++ {
		model, cmd := form.Update(enter)
		form = model.(*huh.Form)
		// huh submits a group through a command, so a form driven without a
		// running program stalls after the last question unless the queued
		// message is delivered.
		for cmd != nil {
			msg := cmd()
			if msg == nil {
				break
			}
			model, cmd = form.Update(msg)
			form = model.(*huh.Form)
		}
	}
	if form.State == huh.StateNormal {
		t.Fatalf("twenty Enters did not finish the form; state is %v", form.State)
	}

	applyAnswers(opts, arch, frontend, theme, db, storage)
	return *opts
}

// The issue this exists for: `grit new` opened on Full while /docs/start told
// a newcomer to run `grit new my-app --triple --next`, so the recommendation
// and the thing Enter gave you were different projects (grit#240).
//
// Both sides are read in one test on purpose. A test over the picker alone
// passes while the page says something else, and a test over the page alone
// passes while the picker opens somewhere else, which is how the pair came
// apart in the first place.
func TestEnterThroughThePromptGivesTheDocumentedProject(t *testing.T) {
	got := enterThrough(t)

	documented := documentedRecommendation(t)
	for flag := range documented {
		switch flag {
		case "triple", "double", "single", "api", "mobile", "full":
			if got.Architecture != scaffold.Architecture(flag) && !(flag == "full" && got.Full) {
				t.Errorf("/docs/start recommends --%s and pressing Enter gives architecture %q",
					flag, got.Architecture)
			}
		case "next", "vite", "tanstack":
			want := scaffold.FrontendNext
			if flag != "next" {
				want = scaffold.FrontendTanStack
			}
			if got.Frontend != want {
				t.Errorf("/docs/start recommends --%s and pressing Enter gives frontend %q",
					flag, got.Frontend)
			}
		}
	}

	// And the rest of the answers, which the page describes in prose rather
	// than as flags: "Go API, Next.js web app, admin panel, Postgres, Redis".
	if got.DBProvider != scaffold.DefaultDBProvider {
		t.Errorf("Enter gives DB_PROVIDER=%q, want the default %q",
			got.DBProvider, scaffold.DefaultDBProvider)
	}
	if got.StorageDriver != scaffold.DefaultStorageDriver {
		t.Errorf("Enter gives STORAGE_DRIVER=%q, want the default %q",
			got.StorageDriver, scaffold.DefaultStorageDriver)
	}
	if got.Theme != scaffold.DefaultTheme {
		t.Errorf("Enter gives THEME=%q, want the default %q", got.Theme, scaffold.DefaultTheme)
	}
}

// Whatever Recommended names has to be an answer the validators accept.
// Opening a picker on a value the next step refuses would turn pressing Enter
// into an error message.
func TestTheRecommendationIsValid(t *testing.T) {
	opts := scaffold.Options{
		Architecture:  scaffold.Architecture(Recommended.Architecture),
		Frontend:      scaffold.Frontend(Recommended.Frontend),
		Theme:         Recommended.Theme,
		DBProvider:    Recommended.DBProvider,
		StorageDriver: Recommended.StorageDriver,
	}
	if err := opts.ValidateDBProvider(); err != nil {
		t.Errorf("recommended database: %v", err)
	}
	if err := opts.ValidateStorageDriver(); err != nil {
		t.Errorf("recommended storage: %v", err)
	}
	if err := opts.ValidateTheme(); err != nil {
		t.Errorf("recommended theme: %v", err)
	}
	if !opts.ShouldIncludeAdmin() {
		t.Error("the recommended architecture has no admin panel, which is the " +
			"thing /docs/start recommends it for")
	}
}

// Exactly one option per picker carries the marker, and it is the recommended
// one. Two marked options is worse than none: it is a recommendation that
// recommends nothing.
func TestOneOptionPerPickerIsMarkedRecommended(t *testing.T) {
	lists := map[string]struct {
		options []huh.Option[string]
		want    string
	}{
		"database": {databaseOptions(), Recommended.DBProvider},
		"storage":  {storageOptions(), Recommended.StorageDriver},
	}
	for name, list := range lists {
		marked := []string{}
		for _, o := range list.options {
			if strings.Contains(o.Key, recommendedSuffix) {
				marked = append(marked, o.Value)
			}
		}
		if len(marked) != 1 || marked[0] != list.want {
			t.Errorf("the %s picker marks %v as recommended, want exactly [%s]",
				name, marked, list.want)
		}
	}

	// The architecture, frontend and theme lists are written out by hand in
	// prompt.go, so they are counted in the source.
	source := readSource(t, "prompt.go")
	if n := strings.Count(source, "Recommended.Architecture)"); n != 1 {
		t.Errorf("%d architecture options are marked against Recommended.Architecture, want 1", n)
	}
	if n := strings.Count(source, `, Recommended.Theme)`); n != 1 {
		t.Errorf("%d theme options are marked against Recommended.Theme, want 1", n)
	}
}

// The hint has one job: tell somebody with no preference that Enter is the way
// out. It is checked here rather than through stdout because what matters is
// the words, not the escape codes around them.
func TestTheHintNamesEnterAndTheRecommendation(t *testing.T) {
	hint := Hint()
	for _, want := range []string{"Enter", Recommended.Architecture, Recommended.Frontend} {
		if !strings.Contains(hint, want) {
			t.Errorf("the hint does not mention %q:\n%s", want, hint)
		}
	}
}

// documentedRecommendation returns the flags of the `grit new` command
// /docs/start puts under "If you do not have a preference, use this".
func documentedRecommendation(t *testing.T) map[string]bool {
	t.Helper()

	page := readRepoFile(t, filepath.Join("docs", "app", "docs", "start", "page.tsx"))
	// The recommended command is the one in the panel, and it is the only
	// `grit new` on the page with flags on it.
	re := regexp.MustCompile(`grit new [a-z0-9-.]+((?: --[a-z-]+)+)`)
	m := re.FindStringSubmatch(page)
	if m == nil {
		t.Fatal("no flagged `grit new` command on /docs/start. If the page " +
			"stopped recommending one, this test has nothing to compare the " +
			"picker against and should be deleted rather than loosened.")
	}
	flags := map[string]bool{}
	for _, f := range strings.Fields(m[1]) {
		flags[strings.TrimPrefix(f, "--")] = true
	}
	if len(flags) == 0 {
		t.Fatalf("parsed no flags out of %q", m[0])
	}
	return flags
}

// readSource reads a file from this package's own directory.
func readSource(t *testing.T, name string) string {
	t.Helper()
	raw, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return strings.ReplaceAll(string(raw), "\r\n", "\n")
}

// readRepoFile reads a file by its path from the repository root.
//
// Line endings are normalised, because git checks the docs out with CRLF on
// Windows and a test that searched for a newline-containing string would pass
// in CI and fail for every Windows contributor.
func readRepoFile(t *testing.T, rel string) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("working directory: %v", err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		dir = filepath.Dir(dir)
	}
	f, err := os.Open(filepath.Join(dir, rel))
	if err != nil {
		t.Fatalf("opening %s: %v", rel, err)
	}
	defer f.Close()
	raw, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("reading %s: %v", rel, err)
	}
	return strings.ReplaceAll(string(raw), "\r\n", "\n")
}
