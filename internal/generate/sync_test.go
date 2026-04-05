package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── goTypeToTS ────────────────────────────────────────────────────────────────

func TestGoTypeToTS(t *testing.T) {
	tests := []struct {
		goType string
		want   string
	}{
		{"string", "string"},
		{"int", "number"},
		{"int64", "number"},
		{"uint", "number"},
		{"uint64", "number"},
		{"float32", "number"},
		{"float64", "number"},
		{"bool", "boolean"},
		{"time.Time", "string"},
		{"*time.Time", "string | null"},
		{"gorm.DeletedAt", "string | null"},
		{"datatypes.JSONSlice[string]", "string[]"},
		{"[]string", "string[]"},
		{"[]uint", "number[]"},
		{"*string", "string | null"},
		{"*int", "number | null"},
	}

	for _, tt := range tests {
		got := goTypeToTS(tt.goType)
		if got != tt.want {
			t.Errorf("goTypeToTS(%q) = %q, want %q", tt.goType, got, tt.want)
		}
	}
}

// ── goTypeToZod ───────────────────────────────────────────────────────────────

func TestGoTypeToZod(t *testing.T) {
	tests := []struct {
		goType      string
		gormTag     string
		wantContain string
	}{
		{"string", "", "z.string()"},
		{"string", "type:text", "z.string()"},
		{"int", "", "z.number().int()"},
		{"int64", "", "z.number().int()"},
		{"uint", "", "z.number().int().nonnegative()"},
		{"float64", "", "z.number()"},
		{"bool", "", "z.boolean()"},
		{"*time.Time", "", "z.string().nullable()"},
		{"time.Time", "", "z.string()"},
		{"datatypes.JSONSlice[string]", "", "z.array(z.string())"},
	}

	for _, tt := range tests {
		got := goTypeToZod(tt.goType, tt.gormTag)
		if !strings.Contains(got, tt.wantContain) {
			t.Errorf("goTypeToZod(%q, %q) = %q, does not contain %q",
				tt.goType, tt.gormTag, got, tt.wantContain)
		}
	}
}

// ── extractTag ────────────────────────────────────────────────────────────────

func TestExtractTag(t *testing.T) {
	tests := []struct {
		tag  string
		key  string
		want string
	}{
		{`gorm:"primarykey" json:"id"`, "json", "id"},
		{`gorm:"primarykey" json:"id"`, "gorm", "primarykey"},
		{`json:"first_name" binding:"required,min=2"`, "json", "first_name"},
		{`json:"email,omitempty"`, "json", "email"},   // omitempty stripped
		{`json:"-"`, "json", "-"},                     // skip sentinel preserved
		{`gorm:"size:255"`, "binding", ""},            // missing key
		{``, "json", ""},                              // empty tag
	}

	for _, tt := range tests {
		got := extractTag(tt.tag, tt.key)
		if got != tt.want {
			t.Errorf("extractTag(%q, %q) = %q, want %q", tt.tag, tt.key, got, tt.want)
		}
	}
}

// ── isAutoField ───────────────────────────────────────────────────────────────

func TestIsAutoField(t *testing.T) {
	auto := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, f := range auto {
		if !isAutoField(f) {
			t.Errorf("isAutoField(%q) = false, want true", f)
		}
	}

	notAuto := []string{"title", "email", "published", "price", "user_id"}
	for _, f := range notAuto {
		if isAutoField(f) {
			t.Errorf("isAutoField(%q) = true, want false", f)
		}
	}
}

// ── buildTSType ───────────────────────────────────────────────────────────────

func TestBuildTSType(t *testing.T) {
	s := GoStruct{
		Name: "Post",
		Fields: []GoField{
			{Name: "ID", GoType: "uint", JSONName: "id"},
			{Name: "Title", GoType: "string", JSONName: "title"},
			{Name: "Published", GoType: "bool", JSONName: "published"},
			{Name: "CreatedAt", GoType: "time.Time", JSONName: "created_at"},
		},
	}

	got := buildTSType(s)

	if !strings.Contains(got, "export interface Post {") {
		t.Errorf("buildTSType missing interface declaration:\n%s", got)
	}
	if !strings.Contains(got, "id: number;") {
		t.Errorf("buildTSType missing id field:\n%s", got)
	}
	if !strings.Contains(got, "title: string;") {
		t.Errorf("buildTSType missing title field:\n%s", got)
	}
	if !strings.Contains(got, "published: boolean;") {
		t.Errorf("buildTSType missing published field:\n%s", got)
	}
	if !strings.Contains(got, "created_at: string;") {
		t.Errorf("buildTSType missing created_at field:\n%s", got)
	}
}

// ── buildZodSchema ────────────────────────────────────────────────────────────

func TestBuildZodSchema(t *testing.T) {
	s := GoStruct{
		Name: "Product",
		Fields: []GoField{
			{Name: "ID", GoType: "uint", JSONName: "id"},
			{Name: "Name", GoType: "string", JSONName: "name"},
			{Name: "Price", GoType: "float64", JSONName: "price"},
			{Name: "Active", GoType: "bool", JSONName: "active"},
			{Name: "CreatedAt", GoType: "time.Time", JSONName: "created_at"},
		},
	}

	got := buildZodSchema(s)

	// Both schemas must be present
	if !strings.Contains(got, "export const CreateProductSchema") {
		t.Errorf("buildZodSchema missing CreateProductSchema:\n%s", got)
	}
	if !strings.Contains(got, "export const UpdateProductSchema") {
		t.Errorf("buildZodSchema missing UpdateProductSchema:\n%s", got)
	}

	// Auto fields (id, created_at) must be excluded from create/update schemas
	if strings.Contains(got, "id:") {
		t.Errorf("buildZodSchema should not include auto field 'id':\n%s", got)
	}
	if strings.Contains(got, "createdAt:") {
		t.Errorf("buildZodSchema should not include auto field 'created_at':\n%s", got)
	}

	// Non-auto fields must be present
	if !strings.Contains(got, "name:") {
		t.Errorf("buildZodSchema missing 'name' field:\n%s", got)
	}
	if !strings.Contains(got, "price:") {
		t.Errorf("buildZodSchema missing 'price' field:\n%s", got)
	}

	// Update schema fields should have .optional()
	// (find the UpdateProductSchema block and check)
	updateIdx := strings.Index(got, "UpdateProductSchema")
	if updateIdx == -1 {
		t.Fatal("UpdateProductSchema not found")
	}
	updateBlock := got[updateIdx:]
	if !strings.Contains(updateBlock, ".optional()") {
		t.Errorf("UpdateProductSchema fields should use .optional():\n%s", updateBlock)
	}

	// Type exports must be present
	if !strings.Contains(got, "type CreateProductInput") {
		t.Errorf("buildZodSchema missing CreateProductInput type:\n%s", got)
	}
	if !strings.Contains(got, "type UpdateProductInput") {
		t.Errorf("buildZodSchema missing UpdateProductInput type:\n%s", got)
	}
}

// ── parseGoStructs ────────────────────────────────────────────────────────────

func TestParseGoStructs(t *testing.T) {
	t.Run("simple struct", func(t *testing.T) {
		content := `package models

import "time"

type Product struct {
	ID        string    ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Name      string    ` + "`" + `gorm:"size:255" json:"name"` + "`" + `
	Price     float64   ` + "`" + `json:"price"` + "`" + `
	Active    bool      ` + "`" + `json:"active"` + "`" + `
	Password  string    ` + "`" + `json:"-"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
}
`
		f := writeTempFile(t, "product.go", content)
		structs, err := parseGoStructs(f)
		if err != nil {
			t.Fatalf("parseGoStructs error: %v", err)
		}
		if len(structs) != 1 {
			t.Fatalf("len(structs) = %d, want 1", len(structs))
		}

		s := structs[0]
		if s.Name != "Product" {
			t.Errorf("Name = %q, want Product", s.Name)
		}

		// Password (json:"-") must be excluded
		for _, f := range s.Fields {
			if f.JSONName == "-" {
				t.Error("field with json:\"-\" should have been skipped")
			}
			if f.Name == "Password" {
				t.Error("Password field (json:\"-\") should have been excluded")
			}
		}

		// id, name, price, active, created_at should be present
		names := make(map[string]bool)
		for _, f := range s.Fields {
			names[f.JSONName] = true
		}
		for _, want := range []string{"id", "name", "price", "active", "created_at"} {
			if !names[want] {
				t.Errorf("field %q missing from parsed struct", want)
			}
		}
	})

	t.Run("multiple structs", func(t *testing.T) {
		content := `package models

type Foo struct {
	ID uint ` + "`" + `json:"id"` + "`" + `
}

type Bar struct {
	Name string ` + "`" + `json:"name"` + "`" + `
}
`
		f := writeTempFile(t, "multi.go", content)
		structs, err := parseGoStructs(f)
		if err != nil {
			t.Fatalf("parseGoStructs error: %v", err)
		}
		if len(structs) != 2 {
			t.Fatalf("len(structs) = %d, want 2", len(structs))
		}
	})

	t.Run("invalid Go file returns error", func(t *testing.T) {
		f := writeTempFile(t, "bad.go", "this is not valid Go {{{")
		_, err := parseGoStructs(f)
		if err == nil {
			t.Error("expected error for invalid Go file, got nil")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		_, err := parseGoStructs("/tmp/does-not-exist-grit-parse.go")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})
}

// ── buildTSType + buildZodSchema round-trip ──────────────────────────────────

func TestSyncRoundTrip(t *testing.T) {
	goSrc := `package models

import (
	"time"
	"gorm.io/datatypes"
)

type Article struct {
	ID        uint                      ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Title     string                    ` + "`" + `gorm:"size:255;not null" json:"title"` + "`" + `
	Body      string                    ` + "`" + `gorm:"type:text" json:"body"` + "`" + `
	Tags      datatypes.JSONSlice[string] ` + "`" + `gorm:"type:json" json:"tags"` + "`" + `
	Published bool                      ` + "`" + `json:"published"` + "`" + `
	CreatedAt time.Time                 ` + "`" + `json:"created_at"` + "`" + `
}
`
	dir := t.TempDir()
	goFile := filepath.Join(dir, "article.go")
	if err := os.WriteFile(goFile, []byte(goSrc), 0644); err != nil {
		t.Fatalf("writing go file: %v", err)
	}

	structs, err := parseGoStructs(goFile)
	if err != nil {
		t.Fatalf("parseGoStructs: %v", err)
	}
	if len(structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(structs))
	}

	s := structs[0]
	tsOutput := buildTSType(s)
	zodOutput := buildZodSchema(s)

	// TS output checks
	if !strings.Contains(tsOutput, "export interface Article {") {
		t.Errorf("TS output missing interface:\n%s", tsOutput)
	}
	if !strings.Contains(tsOutput, "tags: string[];") {
		t.Errorf("TS output missing tags as string[]:\n%s", tsOutput)
	}

	// Zod output checks
	if !strings.Contains(zodOutput, "CreateArticleSchema") {
		t.Errorf("Zod output missing CreateArticleSchema:\n%s", zodOutput)
	}
	if !strings.Contains(zodOutput, "tags: z.array(z.string())") {
		t.Errorf("Zod output missing tags as z.array(z.string()):\n%s", zodOutput)
	}
}
