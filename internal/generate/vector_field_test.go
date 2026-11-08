package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// grit#191: an embedding had nowhere first-class to live, so the retrieval half
// of an AI feature was four hand-built pieces in every project that tried.

// TestVectorNeedsItsDimension. It is part of the column type, so a default
// would be a guess that costs a migration and a full re-embed to correct.
func TestVectorNeedsItsDimension(t *testing.T) {
	if _, err := ParseInlineFields("Doc", "embedding:vector"); err == nil {
		t.Error("a vector with no dimension was accepted")
	} else if !strings.Contains(err.Error(), "1536") {
		t.Errorf("the error does not say what a dimension looks like: %v", err)
	}

	for _, bad := range []string{
		"embedding:vector[0]",
		"embedding:vector[-1]",
		"embedding:vector[many]",
		"embedding:vector[20000]",
	} {
		if _, err := ParseInlineFields("Doc", bad); err == nil {
			t.Errorf("%q was accepted", bad)
		}
	}

	def, err := ParseInlineFields("Doc", "embedding:vector[1536]")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	f := def.Fields[0]
	if !f.IsVector() {
		t.Fatal("the field is not a vector")
	}
	if f.VectorDim() != 1536 {
		t.Errorf("dimension %d, want 1536", f.VectorDim())
	}
	if f.VectorTag() != "1536" {
		t.Errorf("tag %q, want 1536", f.VectorTag())
	}
}

// TestVectorTypesAcrossTheStack: one declaration, four languages.
func TestVectorTypesAcrossTheStack(t *testing.T) {
	def, err := ParseInlineFields("Doc", "embedding:vector[768]")
	if err != nil {
		t.Fatal(err)
	}
	f := def.Fields[0]

	if got := f.GoType(); got != "vectors.Vector" {
		t.Errorf("GoType %q", got)
	}
	if got := f.TSType(); got != "number[]" {
		t.Errorf("TSType %q", got)
	}
	if got := f.ZodType(); !strings.HasPrefix(got, "z.array(z.number())") {
		t.Errorf("ZodType %q", got)
	}
	// And the reverse, which `grit sync` uses when it reads the model back.
	// Without these the type came back as unknown and sync rewrote what
	// generate had just written.
	if got := goTypeToTS("vectors.Vector"); got != "number[]" {
		t.Errorf("goTypeToTS %q", got)
	}
	if got := goTypeToZod("vectors.Vector", ""); got != "z.array(z.number())" {
		t.Errorf("goTypeToZod %q", got)
	}
}

// TestBothModelWritersEmitTheVectorTag reads both sides.
//
// The dimension can only live in a struct tag, because Go has no way to put a
// length in a slice type, and internal/vectors reads it to decide between
// vector(n) and text. Two different functions write a model field line:
// templates.go for a new resource and add_field.go for `grit generate field`.
// One of them emitting the tag and the other not would give a column with no
// dimension, which pgvector allows and no index does.
func TestBothModelWritersEmitTheVectorTag(t *testing.T) {
	const module = "myapp/apps/api"

	// Writer one: a new resource.
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Doc", "title:string,embedding:vector[1536]")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "doc.go"))
	if !strings.Contains(model, `vector:"1536"`) {
		t.Errorf("the generated model has no vector tag, so the column has no dimension:\n%s", model)
	}
	if !strings.Contains(model, "internal/vectors") {
		t.Error("the model names vectors.Vector and does not import it")
	}

	// Writer two: `grit generate field`, which injects the line into a model
	// that already exists.
	root2 := setupMinimalProject(t, module)
	def2, err := ParseInlineFields("Note", "title:string")
	if err != nil {
		t.Fatal(err)
	}
	g2 := newTestGenerator(root2, module, def2)
	if err := g2.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}
	added, err := ParseInlineFields("Note", "embedding:vector[1536]")
	if err != nil {
		t.Fatal(err)
	}
	names2 := MakeNames("Note")
	if err := g2.injectModelField(names2, added.Fields[0]); err != nil {
		t.Fatalf("injectModelField: %v", err)
	}
	model2 := readTestFile(t, filepath.Join(root2, "apps", "api", "internal", "models", "note.go"))
	if !strings.Contains(model2, `vector:"1536"`) {
		t.Errorf("grit generate field wrote the column without its dimension, so it would "+
			"be a bare vector column no index can be built on:\n%s", model2)
	}
}

// TestVectorIsNotAnAdminField: 1536 numbers in a table cell, and a text input
// nobody can fill in.
func TestVectorIsNotAnAdminField(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Doc", "title:string,embedding:vector[1536]")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	// The content the admin resource file is written from. Read from the
	// generator rather than from disk: the minimal test project has no
	// apps/admin, and a missing-path error would say nothing about the fields.
	resource := g.resourceDefinitionFileContent(MakeNames("Doc"))
	if strings.Contains(resource, "embedding") {
		t.Errorf("the admin resource names the embedding, so the list page renders 1536 "+
			"numbers in a cell and the form asks somebody to type them:\n%s", resource)
	}
	// And the title is still there, so this is not passing by generating nothing.
	if !strings.Contains(resource, `key: "title"`) {
		t.Error("the admin resource has no title column either, so this test is vacuous")
	}
}

// TestVectorIsNotImported: a CSV import of an embedding is not a thing. The
// importer assigned the cell's string straight to the field, which did not
// compile.
func TestVectorIsNotImported(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Doc", "title:string,embedding:vector[1536]")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "doc_import.go"))
	if strings.Contains(src, "item.Embedding") {
		t.Errorf("the CSV importer assigns to the embedding:\n%s", src)
	}
	if !strings.Contains(src, "item.Title") {
		t.Error("the importer does not assign the title either, so this test is vacuous")
	}
}

// TestVectorIndexSQLNamesTheOperatorClass.
//
// The operator class has to match the distance operator the query uses or the
// index is silently not used, which is the thing the issue asked not to have to
// get right by hand.
func TestVectorIndexSQLNamesTheOperatorClass(t *testing.T) {
	up := vectorIndexUpSQL("20260101000000", "index_docs_embedding", "docs", "embedding", "idx_docs_embedding", 1536)

	for _, want := range []string{
		NoTransactionMarkerText,
		"USING hnsw (embedding vector_cosine_ops)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_docs_embedding",
		"vector_l2_ops",
		"ivfflat",
	} {
		if !strings.Contains(up, want) {
			t.Errorf("the index migration does not mention %q:\n%s", want, up)
		}
	}

	down := vectorIndexDownSQL("20260101000000", "index_docs_embedding", "idx_docs_embedding")
	if !strings.Contains(down, "DROP INDEX CONCURRENTLY IF EXISTS idx_docs_embedding") {
		t.Errorf("the rollback does not drop the index:\n%s", down)
	}

	// Past the index limit it explains itself instead of writing a statement
	// that fails.
	wide := vectorIndexUpSQL("20260101000000", "index_docs_embedding", "docs", "embedding", "idx_docs_embedding", 4096)
	// USING hnsw rather than CREATE INDEX: the migration header's own
	// boilerplate explains that the no-transaction marker is required for
	// CREATE INDEX CONCURRENTLY, so that phrase is in every one of these files.
	if strings.Contains(wide, "USING hnsw") || strings.Contains(wide, "USING ivfflat") {
		t.Errorf("an index was written for 4096 dimensions, which pgvector refuses:\n%s", wide)
	}
	if !strings.Contains(wide, "2000") {
		t.Errorf("the file does not say what the limit is:\n%s", wide)
	}
}
