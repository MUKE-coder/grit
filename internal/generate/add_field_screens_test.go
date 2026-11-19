package generate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// expoProject stands up a project with an Expo app and one generated form,
// recorded in the manifest as Grit wrote it.
//
// The recording is the point: every decision syncMobileForm makes comes from
// what the manifest says about the file, so a fixture that skips it tests the
// untracked path and nothing else.
func expoProject(t *testing.T, body string) (root, form string) {
	t.Helper()
	root = t.TempDir()
	form = filepath.Join(root, "apps", "expo", "components", "resource-forms", "tasks-form.tsx")

	release, err := manifest.Start(root, "3.551.0", "generate")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = release() })

	writeTestFile(t, form, body)
	manifest.Note(form, body)
	return root, form
}

// taskDefinition is a resource as it stood before a field was added.
func taskDefinition() *ResourceDefinition {
	return &ResourceDefinition{
		Name: "Task",
		Fields: []Field{
			{Name: "title", Type: "string", Required: true},
		},
	}
}

// A field reaches the Expo form.
//
// This is the whole of grit#87 as it was reported: the field was in the API
// response and absent from the form on an Android emulator. Nothing updated
// that file, ever.
func TestTheExpoFormIsRebuiltWithTheNewField(t *testing.T) {
	root, form := expoProject(t, "the form Grit wrote\n")

	def := taskDefinition()
	def.Fields = append(def.Fields, Field{Name: "priority", Type: "int"})

	updated, kept := syncMobileForm(root, def, MakeNames("Task"))
	if updated == "" {
		t.Fatalf("the form was not rebuilt (kept: %q)", kept)
	}

	body := read(t, form)
	if !strings.Contains(body, "priority") {
		t.Errorf("the rebuilt form does not mention the new field:\n%s", body)
	}
	// And the field it already had is still there, which is what distinguishes
	// a rebuild from a reset.
	if !strings.Contains(body, "title") {
		t.Error("the rebuilt form lost the field it already had")
	}
}

// The rebuild is recorded, so the next one is still allowed.
//
// Writing the file without noting it leaves the manifest holding the hash of
// the previous version, which reads as "you have edited it": the command works
// once and refuses from then on, blaming the user for its own write.
func TestARebuiltFormIsRecorded(t *testing.T) {
	root, form := expoProject(t, "the form Grit wrote\n")
	def := taskDefinition()

	if updated, kept := syncMobileForm(root, def, MakeNames("Task")); updated == "" {
		t.Fatalf("first rebuild refused (kept: %q)", kept)
	}
	if !manifest.IsUnchanged(form) {
		t.Fatal("the rebuild did not record what it wrote, so the next one is refused")
	}

	def.Fields = append(def.Fields, Field{Name: "notes", Type: "text"})
	if updated, kept := syncMobileForm(root, def, MakeNames("Task")); updated == "" {
		t.Errorf("second rebuild refused (kept: %q)", kept)
	}
	if !strings.Contains(read(t, form), "notes") {
		t.Error("the second rebuild did not land")
	}
}

// A form somebody has edited is theirs.
//
// The manifest decides this, which is the same rule grit upgrade applies to
// every framework-owned file. A rebuild that ignores it is a command that
// silently discards work, and this one runs as a side effect of adding a
// field, which is not a moment anybody expects to lose a screen.
func TestAnEditedExpoFormIsKept(t *testing.T) {
	root, form := expoProject(t, "the form Grit wrote\n")
	mine := "the form Grit wrote\n// and my own control\n"
	writeTestFile(t, form, mine)

	updated, kept := syncMobileForm(root, taskDefinition(), MakeNames("Task"))
	if updated != "" {
		t.Error("an edited form was rebuilt")
	}
	if !strings.Contains(kept, "you have edited it") {
		t.Errorf("the report does not say why it was kept: %q", kept)
	}
	if got := read(t, form); got != mine {
		t.Errorf("the edited form was written over:\n%s", got)
	}
}

// A form Grit has no record of is left alone, and not called an edit.
//
// A project older than the manifest lands here. "You have edited it" would be
// a guess presented as fact, and the advice that follows it (take Grit's
// version with --force) is wrong for a file Grit never wrote.
func TestAnUntrackedExpoFormIsNotCalledAnEdit(t *testing.T) {
	root := t.TempDir()
	form := filepath.Join(root, "apps", "expo", "components", "resource-forms", "tasks-form.tsx")
	writeTestFile(t, form, "hand written\n")

	release, err := manifest.Start(root, "3.551.0", "generate")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = release() }()

	updated, kept := syncMobileForm(root, taskDefinition(), MakeNames("Task"))
	if updated != "" {
		t.Fatal("a form Grit never wrote was rebuilt")
	}
	if strings.Contains(kept, "you have edited it") {
		t.Errorf("an untracked form is reported as an edit: %q", kept)
	}
	if !strings.Contains(kept, "no record") {
		t.Errorf("the report does not say what happened: %q", kept)
	}
}

// Nothing is written into a project with no Expo app.
func TestNoExpoAppMeansNothingToDo(t *testing.T) {
	root := t.TempDir()
	updated, kept := syncMobileForm(root, taskDefinition(), MakeNames("Task"))
	if updated != "" || kept != "" {
		t.Errorf("reported %q / %q for a project with no Expo app", updated, kept)
	}
}

// grit generate field renders the definition WITH the field it just added.
//
// The generator's own definition is the resource as it was before the field,
// because DefinitionFromModel read it at the top of AddField. Rendering the
// form from that rebuilds the form exactly as it was, and reports "rebuilt
// from the model" while the field is not on it: worse than not trying, because
// it claims the job is done. Caught by running it, and this is the guard.
func TestAddingAFieldRendersTheFieldAndNotTheOldDefinition(t *testing.T) {
	root, form := expoProject(t, "the form Grit wrote\n")

	g := &Generator{
		Root:         root,
		Architecture: "triple",
		Module:       "example.com/app",
		Definition:   taskDefinition(),
	}
	added := Field{Name: "priority", Type: "int"}

	result := g.updateFieldScreens(MakeNames("Task"), added)

	if len(result.Updated) == 0 {
		t.Fatalf("nothing was updated (kept: %v)", result.Kept)
	}
	if !strings.Contains(read(t, form), "priority") {
		t.Error("the form was rebuilt without the field that was just added")
	}
	// And the generator's own definition did not grow underneath the
	// injections that already ran against it.
	if len(g.Definition.Fields) != 1 {
		t.Errorf("the generator's definition was mutated: %d fields", len(g.Definition.Fields))
	}
}

// One writer for the Expo form, called by both commands.
//
// grit sync and grit generate field both have to put a field on that screen.
// Two implementations of one rebuild is how the two come to disagree about
// what the form looks like, and the disagreement shows up as a form that
// changes shape depending on which command you ran last.
func TestBothCommandsUseTheSameFormWriter(t *testing.T) {
	for _, caller := range []string{"sync.go", "add_field_screens.go"} {
		body := read(t, caller)
		if !strings.Contains(body, "syncMobileForm(") {
			t.Errorf("%s does not call syncMobileForm", caller)
		}
		if strings.Contains(body, "mobileFormSource(") {
			t.Errorf("%s renders the form itself rather than going through "+
				"syncMobileForm, so there are two writers for one file", caller)
		}
	}
}

// The desktop half goes through the function grit sync already uses.
//
// SyncDesktopScreen inserts a column between the grit:cols fences. A second
// insert written for grit generate field would be two writers for one array,
// which is how the two would come to disagree about what a column looks like.
func TestTheDesktopHalfReusesSyncDesktopScreen(t *testing.T) {
	body := read(t, "add_field_screens.go")
	if !strings.Contains(body, "SyncDesktopScreen(") {
		t.Error("grit generate field does not call SyncDesktopScreen")
	}
	if strings.Contains(body, desktopColsStart) {
		t.Error("grit generate field writes into the desktop fences itself")
	}
}
