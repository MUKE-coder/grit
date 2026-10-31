package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Review 2.14: two writers produced packages/shared/types/<model>.ts and they
// disagreed, so `grit sync` -- the step the generator itself prints as the next
// thing to run -- silently rewrote what the generator had just written.
//
// What each got wrong was the opposite of the other. The generator knew the
// field types and wrote Money, FileRef and a union of an enum's values; sync
// read the Go struct with a mapping that knows the builtins and wrote unknown
// and string. Sync knew every column the API serialises; the generator built
// the interface from the definition, so version, archived_at and a --tree
// resource's path, depth, position and children were on no generated type at
// all.
//
// Neither is fixed by a second set of rules. There is one writer now, and this
// test is what holds the generator to it: it goes through writeTSTypes, then
// through buildTSType directly, and requires the same bytes.
//
// Which is comparing a function with itself, and that is the point and the
// limit of it. It cannot tell whether buildTSType is right; it fails the moment
// writeTSTypes grows an interface of its own again, which is the thing that was
// wrong. What says the output is right is TestTheTypeHasTheColumnsTheAPISends
// below, which reads the model, and running grit sync over a real project,
// which is how the disagreement was measured in the first place.
func TestGenerateThenSyncChangesNothing(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Product",
		"name:string,body:text,price:money,cover:file,gallery:files,status:select:draft|live,tags:string_array,rank:int,live:bool,due:datetime")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	if err := g.writeTSTypes(names); err != nil {
		t.Fatalf("writeTSTypes: %v", err)
	}

	generated := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "product.ts"))

	// What sync would write, from the same model, through the same code path
	// Sync() uses.
	modelPath := filepath.Join(root, "apps", "api", "internal", "models", "product.go")
	structs, err := parseGoStructs(modelPath)
	if err != nil {
		t.Fatalf("parsing the model: %v", err)
	}
	var synced string
	for _, s := range structs {
		if s.Name == "Product" {
			synced = buildTSType(s)
		}
	}
	if synced == "" {
		t.Fatal("no Product struct in the model that was just written")
	}

	if generated != synced {
		t.Errorf("grit sync would rewrite the file grit generate just wrote.\ngenerate:\n%s\nsync:\n%s", generated, synced)
	}
}

// The columns the generator used to leave out. Every one of them is on the
// model, which means the API sends it and the TypeScript said it did not
// exist.
func TestTheTypeHasTheColumnsTheAPISends(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Chapter", "label:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.Tree = true

	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	if err := g.writeTSTypes(names); err != nil {
		t.Fatalf("writeTSTypes: %v", err)
	}
	got := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "chapter.ts"))

	// Read off the model rather than listed by hand, so a column added to the
	// template is covered without anybody remembering this test.
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "chapter.go"))
	for _, name := range jsonNamesIn(model) {
		if !strings.Contains(got, name+"?:") && !strings.Contains(got, name+":") {
			t.Errorf("the model serialises %q and the TypeScript type has no such property:\n%s", name, got)
		}
	}
}

// zodObject is the body of one z.object in a generated schema file.
//
// A schema file holds three of them now: the entity schema, which describes
// what the API sends, and the create and update pair, which describe what it
// accepts. A test about what a form may be asked for has to say which one it
// means. Before the entity schema was written by the generator as well as by
// sync, "the string is not in the file" stood in for "the field is not in the
// create schema", and the two stopped being the same thing: the owner of an
// --owned-by row and a server_set status are both absent from create and
// present in the response.
func zodObject(t *testing.T, src, name string) string {
	t.Helper()
	head := "export const " + name + " = z.object({\n"
	i := strings.Index(src, head)
	if i < 0 {
		t.Fatalf("no %s in the schema file:\n%s", name, src)
	}
	rest := src[i+len(head):]
	end := strings.Index(rest, "})")
	if end < 0 {
		t.Fatalf("%s is not closed:\n%s", name, src)
	}
	return rest[:end]
}

// jsonNamesIn is every json name a Go source file's struct tags declare, apart
// from the ones marked "-".
func jsonNamesIn(source string) []string {
	var out []string
	for _, line := range strings.Split(source, "\n") {
		raw := extractRawTag(line, "json")
		if raw == "" {
			continue
		}
		name := strings.Split(raw, ",")[0]
		if name == "" || name == "-" {
			continue
		}
		out = append(out, name)
	}
	return out
}

// A belongs_to and a many_to_many naming the same model emitted the import
// twice, which is TS2300 and fails the whole workspace rather than one file.
// Found by generating a resource with both, which no fixture had.
func TestOneImportPerRelatedModel(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	if err := os.MkdirAll(filepath.Join(root, "apps", "api", "internal", "models"), 0o755); err != nil {
		t.Fatal(err)
	}
	// The model, written by hand, because what is being tested is the emitter
	// and this is the shape that broke: one relationship of each kind to the
	// same model, plus a self-reference, which must not be imported at all.
	model := `package models

type Article struct {
	ID       string    ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	AuthorID string    ` + "`" + `gorm:"size:36" json:"author_id"` + "`" + `
	Author   *Author   ` + "`" + `gorm:"foreignKey:AuthorID" json:"author,omitempty"` + "`" + `
	Tags     []Author  ` + "`" + `gorm:"many2many:article_tags" json:"tags,omitempty"` + "`" + `
	ParentID *string   ` + "`" + `gorm:"size:36;default:null" json:"parent_id"` + "`" + `
	Parent   *Article  ` + "`" + `gorm:"foreignKey:ParentID" json:"parent,omitempty"` + "`" + `
}
`
	path := filepath.Join(root, "apps", "api", "internal", "models", "article.go")
	if err := os.WriteFile(path, []byte(model), 0o644); err != nil {
		t.Fatal(err)
	}
	structs, err := parseGoStructs(path)
	if err != nil {
		t.Fatal(err)
	}
	got := buildTSType(structs[0])

	if n := strings.Count(got, `import type { Author }`); n != 1 {
		t.Errorf("Author is imported %d times, want 1:\n%s", n, got)
	}
	if strings.Contains(got, `import type { Article }`) {
		t.Errorf("the interface imports itself, which is TS2440:\n%s", got)
	}
	for _, want := range []string{"author?: Author;", "tags?: Author[];", "parent?: Article;", "parent_id: string | null;"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q:\n%s", want, got)
		}
	}
}
