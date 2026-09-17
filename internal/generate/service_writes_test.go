package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// renderService is the service the generator writes today, formatted as it
// lands on disk.
func renderService(t *testing.T, module, name, fields string, setup func(*ResourceDefinition)) string {
	t.Helper()
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields(name, fields)
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	if setup != nil {
		setup(def)
	}
	g := newTestGenerator(root, module, def)
	return codefmt.Go(g.serviceSource(g.Names()))
}

func fixture(t *testing.T, name string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "l17", name))
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(raw), "\r\n", "\n")
}

// L17: a write with belongs_to relations is one statement with RETURNING, and
// reads only the rows it points at. Patch writes the same way for every resource
// without links or lines.
func TestGeneratedWritesDoNotReadTheRowBack(t *testing.T) {
	const module = "l19x-old/apps/api"
	contact := renderService(t, module, "Contact", "name:string,phone:string,group:belongs_to:Group", nil)
	for _, want := range []string{
		"if err := s.write(db).Create(item).Error; err != nil {",
		"written := s.write(db).Model(item).Scopes(pre.Scope).Updates(updates)",
		"func (s *ContactService) relations(db *gorm.DB, item *models.Contact) error {",
		`if err := db.Where("id = ?", item.GroupID).Limit(1).Find(&related).Error; err != nil {`,
		`} else if err := db.Preload("Group").First(item, "id = ?", item.ID).Error; err != nil {`,
	} {
		if !strings.Contains(contact, want) {
			t.Errorf("the Contact service is missing %q", want)
		}
	}
	if strings.Count(contact, "if err := s.relations(db, item); err != nil {") != 3 {
		t.Error("create, update and patch do not all read the relations after RETURNING")
	}
	if strings.Contains(contact, "\tif err := db.Preload(\"Group\").First(item,") {
		t.Error("a write still reads its row back on every dialect")
	}

	group := renderService(t, module, "Group", "name:string", nil)
	if strings.Count(group, "written := s.write(db).Model(item)") != 2 || strings.Contains(group, "relations(") {
		t.Error("the Group service's update and patch do not both write with RETURNING")
	}

	// A self-reference is a nullable key; links and trees keep the transaction
	// and the reload.
	node := renderService(t, module, "Node", "name:string,parent:belongs_to:Node", nil)
	if !strings.Contains(node, `if item.ParentID != nil && *item.ParentID != "" {`) {
		t.Error("a self-reference's relation does not check its nullable key")
	}
	linked := renderService(t, module, "Linked", "name:string,category:belongs_to:Category,tags:many_to_many:Tag", nil)
	if strings.Contains(linked, "relations(") {
		t.Error("a resource with links reads relations after RETURNING, but its write is a transaction")
	}
	menu := renderService(t, module, "Menu", "name:string,parent:belongs_to:Menu", func(d *ResourceDefinition) { d.Tree = true })
	if strings.Contains(menu, "relations(") || strings.Contains(menu, "s.write(db)") {
		t.Error("a tree writes with RETURNING, but its hooks move a subtree")
	}

	// grit upgrade: a service the generator wrote before becomes the one it
	// writes now, and a second pass changes nothing.
	for _, tc := range []struct{ service, model, want string }{
		{"contact_service.go.txt", "contact_model.go.txt", contact},
		{"group_service.go.txt", "group_model.go.txt", group},
	} {
		old := fixture(t, tc.service)
		model := fixture(t, tc.model)
		out, changed := scaffold.RepairGeneratedServiceWrites(old, model)
		if !changed {
			t.Fatalf("%s was not repaired", tc.service)
		}
		if got := codefmt.Go(out); got != tc.want {
			t.Errorf("repairing %s did not produce the generated service.\n--- repaired\n%s\n--- generated\n%s", tc.service, got, tc.want)
		}
		if again, changed := scaffold.RepairGeneratedServiceWrites(codefmt.Go(out), model); changed || again != codefmt.Go(out) {
			t.Errorf("repairing %s twice changed it again", tc.service)
		}
	}
	if _, changed := scaffold.RepairGeneratedServiceWrites(linked, ""); changed {
		t.Error("a service with links was rewritten")
	}
}
