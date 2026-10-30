package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Review 2.23: re-describing an unedited resource with fewer fields dropped
// columns silently.
//
// The report is the fix, not a refusal: regenerating with fewer fields is a
// legitimate thing to do, and the files are in git. What was missing is anybody
// saying it happened, because the output of the second run looks exactly like
// the output of the first.
func TestARegenerateSaysWhichColumnsItStopsReading(t *testing.T) {
	const module = "dropfields/apps/api"
	root := setupMinimalProject(t, module)

	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	// A model as the generator wrote it, with three of the developer's columns
	// and the framework's own.
	if err := os.WriteFile(filepath.Join(models, "contact.go"), []byte(`package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID        string         `+"`"+`gorm:"primarykey;size:36" json:"id"`+"`"+`
	Name      string         `+"`"+`gorm:"size:255" json:"name"`+"`"+`
	Email     string         `+"`"+`gorm:"size:255" json:"email"`+"`"+`
	Company   string         `+"`"+`gorm:"size:255" json:"company"`+"`"+`
	GroupID   *string        `+"`"+`gorm:"size:36;index" json:"group_id"`+"`"+`
	Group     *Group         `+"`"+`gorm:"foreignKey:GroupID" json:"group,omitempty"`+"`"+`
	Version   int            `+"`"+`gorm:"not null;default:1" json:"version"`+"`"+`
	CreatedAt time.Time      `+"`"+`gorm:"index" json:"created_at"`+"`"+`
	UpdatedAt time.Time      `+"`"+`gorm:"index" json:"updated_at"`+"`"+`
	DeletedAt gorm.DeletedAt `+"`"+`gorm:"index" json:"-"`+"`"+`
	ArchivedAt *time.Time    `+"`"+`gorm:"index" json:"archived_at,omitempty"`+"`"+`
}
`), 0o644); err != nil {
		t.Fatal(err)
	}

	// Regenerated with name and the relation only: email and company go.
	def, err := ParseInlineFields("Contact", "name:string,group:belongs_to:Group")
	if err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() {
		newTestGenerator(root, module, def).reportDroppedFields(MakeNames("Contact"))
	})

	for _, want := range []string{"Email", "Company", "stops reading 2 columns"} {
		if !strings.Contains(out, want) {
			t.Errorf("the report does not mention %q:\n%s", want, out)
		}
	}
	// The framework's own columns are not the developer's and a definition
	// never mentions them: reported as dropped, they would be reported on
	// every single regenerate and the message would be worth nothing.
	for _, never := range []string{"Version", "CreatedAt", "UpdatedAt", "DeletedAt", "ArchivedAt", "ID"} {
		if strings.Contains(out, never) {
			t.Errorf("the report names %s, which no definition ever mentions:\n%s", never, out)
		}
	}
	// A belongs_to writes GroupID and Group, and the definition says "group".
	// Both spellings are kept, or every regenerate would claim to drop a
	// relation it is about to write again.
	if strings.Contains(out, "Group") {
		t.Errorf("the report names the relation it is keeping:\n%s", out)
	}
	// And it says what it means, because "stops reading" on its own invites
	// the wrong conclusion about the data.
	if !strings.Contains(out, "still in the database") {
		t.Errorf("the report does not say the data is still there:\n%s", out)
	}
}

// Nothing to say for a resource being generated for the first time, and
// nothing to say when the field list has not shrunk. Those are almost every
// run, and a warning that fires on them is a warning nobody reads.
func TestTheDroppedColumnReportIsSilentWhenItShouldBe(t *testing.T) {
	const module = "dropfields2/apps/api"
	root := setupMinimalProject(t, module)
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}

	def, err := ParseInlineFields("Contact", "name:string,email:email")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)

	// First generation: no model yet.
	if out := captureStdout(t, func() { g.reportDroppedFields(MakeNames("Contact")) }); out != "" {
		t.Errorf("a first generation printed:\n%s", out)
	}

	// Same fields again.
	if err := os.WriteFile(filepath.Join(models, "contact.go"), []byte(`package models

import "time"

type Contact struct {
	ID        string    `+"`"+`json:"id"`+"`"+`
	Name      string    `+"`"+`json:"name"`+"`"+`
	Email     string    `+"`"+`json:"email"`+"`"+`
	CreatedAt time.Time `+"`"+`json:"created_at"`+"`"+`
}
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if out := captureStdout(t, func() { g.reportDroppedFields(MakeNames("Contact")) }); out != "" {
		t.Errorf("a regenerate with the same fields printed:\n%s", out)
	}

	// More fields than before is not a drop either.
	wider, err := ParseInlineFields("Contact", "name:string,email:email,phone:tel")
	if err != nil {
		t.Fatal(err)
	}
	if out := captureStdout(t, func() {
		newTestGenerator(root, module, wider).reportDroppedFields(MakeNames("Contact"))
	}); out != "" {
		t.Errorf("adding a field printed a drop:\n%s", out)
	}
}
