package generate

import (
	"os"
	"testing"
)

// ── Pluralize benchmarks ──────────────────────────────────────────────────────

func BenchmarkPluralize(b *testing.B) {
	words := []string{"Post", "Category", "Person", "child", "leaf", "status", "box"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Pluralize(words[i%len(words)])
	}
}

// ── Field method benchmarks ───────────────────────────────────────────────────

func BenchmarkField_GoType(b *testing.B) {
	fields := []Field{
		{Type: "string"}, {Type: "int"}, {Type: "float"}, {Type: "bool"},
		{Type: "datetime"}, {Type: "belongs_to"}, {Type: "many_to_many"},
		{Type: "string_array"},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fields[i%len(fields)].GoType()
	}
}

func BenchmarkField_ZodType(b *testing.B) {
	f := Field{Type: "string", Required: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.ZodType()
	}
}

func BenchmarkField_GORMTag(b *testing.B) {
	fields := []Field{
		{Type: "string", Unique: true},
		{Type: "text"},
		{Type: "belongs_to"},
		{Type: "slug"},
		{Type: "string_array"},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fields[i%len(fields)].GORMTag()
	}
}

// ── ParseInlineFields benchmarks ──────────────────────────────────────────────

func BenchmarkParseInlineFields(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseInlineFields("Post", "title:string,content:text,published:bool,price:float,author:belongs_to:User")
	}
}

// ── injectBefore benchmark ────────────────────────────────────────────────────

func BenchmarkInjectBefore(b *testing.B) {
	const src = `package models

func Models() []interface{} {
	return []interface{}{
		&User{},
		// grit:models
	}
}
`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		path := writeTempFileBench(b, "bench.go", src)
		_ = injectBefore(path, "// grit:models", "\t\t&Post{},")
	}
}

// ── goTypeToTS benchmark ──────────────────────────────────────────────────────

func BenchmarkGoTypeToTS(b *testing.B) {
	types := []string{
		"string", "int64", "uint", "float64", "bool",
		"*time.Time", "time.Time", "gorm.DeletedAt", "[]uint",
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		goTypeToTS(types[i%len(types)])
	}
}

// writeTempFileBench creates a temp file for benchmarks.
// (writeTempFile in inject_test.go requires *testing.T; benchmarks need *testing.B)
func writeTempFileBench(b *testing.B, name, content string) string {
	b.Helper()
	dir := b.TempDir()
	path := dir + "/" + name
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		b.Fatalf("writeTempFileBench: %v", err)
	}
	return path
}
