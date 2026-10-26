package generate

import (
	"testing"
)

// grit sync read GORM association structs as missing fields and appended a
// free-text input for each: a resource with five belongs_to got five duplicate
// text boxes beside the relationship-selects the generator had already written
// correctly, plus a duplicate column beside the one that renders the related
// record's label (grit#101). generate and sync disagreed about what a field
// is, so the form was wrong the moment it was generated.
func TestIsAssociation(t *testing.T) {
	fields := []GoField{
		{Name: "ID", GoType: "string"},
		{Name: "Title", GoType: "string"},
		{Name: "ContactID", GoType: "*string"},
		{Name: "Contact", GoType: "*Contact", GORMTag: "foreignKey:ContactID"},
		{Name: "Tags", GoType: "[]Tag", GORMTag: "many2many:deal_tags"},
		{Name: "Items", GoType: "[]DealItem", GORMTag: "foreignKey:DealID"},
		// A hand-written model may carry no tag at all: GORM infers by
		// convention, and so does this.
		{Name: "Owner", GoType: "*User"},
		{Name: "OwnerID", GoType: "string"},
		// Not associations, whatever they look like.
		{Name: "PublishedAt", GoType: "*time.Time"},
		{Name: "DeletedAt", GoType: "gorm.DeletedAt"},
		{Name: "Meta", GoType: "datatypes.JSON"},
		{Name: "Scores", GoType: "[]int"},
	}
	byName := map[string]GoField{}
	for _, f := range fields {
		byName[f.Name] = f
	}

	for _, c := range []struct {
		field string
		want  bool
		why   string
	}{
		{"Contact", true, "a belongs_to struct with a foreignKey tag"},
		{"Tags", true, "a many2many slice"},
		{"Items", true, "a has-many slice"},
		{"Owner", true, "a struct with a matching OwnerID, inferred like GORM does"},
		{"ContactID", false, "the key column is a real field and belongs on the form"},
		{"Title", false, "an ordinary column"},
		{"PublishedAt", false, "a pointer to a stdlib type is a nullable column"},
		{"DeletedAt", false, "gorm bookkeeping, but still not an association"},
		{"Meta", false, "a JSON column"},
		{"Scores", false, "a slice of a builtin is a column, not a relation"},
	} {
		if got := isAssociation(byName[c.field], fields); got != c.want {
			t.Errorf("%s: got %v, want %v (%s)", c.field, got, c.want, c.why)
		}
	}
}

// archived_at is the framework's own bookkeeping, and sync was offering it as
// an editable date on every generated form.
func TestArchivedAtIsNotAFormField(t *testing.T) {
	for _, key := range []string{"id", "created_at", "updated_at", "deleted_at", "archived_at"} {
		if !isAutoField(key) {
			t.Errorf("%s is framework bookkeeping and should never be auto-added to a form", key)
		}
	}
	for _, key := range []string{"title", "amount", "contact_id"} {
		if isAutoField(key) {
			t.Errorf("%s is a real field and should be offered", key)
		}
	}
}
