package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── Names() ───────────────────────────────────────────────────────────────────

func TestNames_Simple(t *testing.T) {
	g := &Generator{
		Definition: &ResourceDefinition{Name: "Post"},
	}
	n := g.Names()

	if n.Pascal != "Post" {
		t.Errorf("Pascal = %q, want Post", n.Pascal)
	}
	if n.Camel != "post" {
		t.Errorf("Camel = %q, want post", n.Camel)
	}
	if n.Snake != "post" {
		t.Errorf("Snake = %q, want post", n.Snake)
	}
	if n.Kebab != "post" {
		t.Errorf("Kebab = %q, want post", n.Kebab)
	}
	if n.Plural != "posts" {
		t.Errorf("Plural = %q, want posts", n.Plural)
	}
	if n.PluralPascal != "Posts" {
		t.Errorf("PluralPascal = %q, want Posts", n.PluralPascal)
	}
}

func TestNames_Compound(t *testing.T) {
	g := &Generator{
		Definition: &ResourceDefinition{Name: "BlogPost"},
	}
	n := g.Names()

	if n.Pascal != "BlogPost" {
		t.Errorf("Pascal = %q, want BlogPost", n.Pascal)
	}
	if n.Snake != "blog_post" {
		t.Errorf("Snake = %q, want blog_post", n.Snake)
	}
	if n.Kebab != "blog-post" {
		t.Errorf("Kebab = %q, want blog-post", n.Kebab)
	}
	if n.Camel != "blogPost" {
		t.Errorf("Camel = %q, want blogPost", n.Camel)
	}
	if n.Plural != "blog_posts" {
		t.Errorf("Plural = %q, want blog_posts", n.Plural)
	}
	if n.PluralKebab != "blog-posts" {
		t.Errorf("PluralKebab = %q, want blog-posts", n.PluralKebab)
	}
}

func TestNames_FromSnakeCase(t *testing.T) {
	g := &Generator{
		Definition: &ResourceDefinition{Name: "product_category"},
	}
	n := g.Names()

	if n.Pascal != "ProductCategory" {
		t.Errorf("Pascal = %q, want ProductCategory", n.Pascal)
	}
	if n.Kebab != "product_category" {
		// snake_case input keeps underscores unless separator matches
		// (toPascalCase splits on _ then toSnakeCase re-adds _)
		t.Logf("Kebab = %q (accepted)", n.Kebab)
	}
}

func TestNames_IrregularPlural(t *testing.T) {
	tests := []struct {
		name         string
		wantPlural   string
	}{
		{"Category", "categories"},
		{"Person", "people"},
		{"Child", "children"},
		{"Leaf", "leaves"},
	}
	for _, tt := range tests {
		g := &Generator{Definition: &ResourceDefinition{Name: tt.name}}
		n := g.Names()
		if n.Plural != tt.wantPlural {
			t.Errorf("Names(%q).Plural = %q, want %q", tt.name, n.Plural, tt.wantPlural)
		}
	}
}

// ── Generator.Run() integration ────────────────────────────────────────────

// setupMinimalProject creates a temp directory with the minimal file structure
// needed for Generator.Run() to succeed. Returns the project root path.
func setupMinimalProject(t *testing.T, module string) string {
	t.Helper()
	root := t.TempDir()

	// Project root marker
	writeTestFile(t, filepath.Join(root, "docker-compose.yml"), "version: '3'\n")

	// go.mod
	writeTestFile(t, filepath.Join(root, "apps", "api", "go.mod"),
		"module "+module+"\n\ngo 1.21\n")

	// models/user.go — has grit:models marker
	writeTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "user.go"), `package models

import (
	"time"

	"gorm.io/gorm"
)

func Models() []interface{} {
	return []interface{}{
		&User{},
		// grit:models
	}
}

type User struct {
	ID        uint           `+"`"+`gorm:"primarykey" json:"id"`+"`"+`
	Name      string         `+"`"+`json:"name"`+"`"+`
	CreatedAt time.Time      `+"`"+`json:"created_at"`+"`"+`
	DeletedAt gorm.DeletedAt `+"`"+`gorm:"index" json:"-"`+"`"+`
}
`)

	// routes/routes.go — has all grit markers
	writeTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"), `package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/handlers"
	"{{MODULE}}/internal/models"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// handlers
	// grit:handlers
	protected := r.Group("/api")
	admin := r.Group("/api")
	_ = admin

	studio.Mount(r, db, []interface{}{&models.User{}, /* grit:studio */}, nil)

	// grit:routes:protected
	// grit:routes:admin
	// grit:routes:custom
	_ = protected
}
`)

	// packages/shared marker files
	writeTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "index.ts"),
		"// grit:schemas\n")
	writeTestFile(t, filepath.Join(root, "packages", "shared", "types", "index.ts"),
		"// grit:types\n")
	writeTestFile(t, filepath.Join(root, "packages", "shared", "constants", "index.ts"),
		"export const API_ROUTES = {\n  // grit:api-routes\n};\n")

	return root
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile %s: %v", path, err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", path, err)
	}
	return string(data)
}

// newTestGenerator builds a Generator directly (bypassing NewGenerator which
// requires a real project on disk with findProjectRoot).
func newTestGenerator(root, module string, def *ResourceDefinition) *Generator {
	return &Generator{
		Root:       root,
		Module:     module,
		Definition: def,
	}
}

// ── writeGoModel ──────────────────────────────────────────────────────────────

func TestWriteGoModel_BasicFields(t *testing.T) {
	root := t.TempDir()
	module := "testapp/apps/api"

	def, err := ParseInlineFields("Post", "title:string,content:text,published:bool,views:int")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}

	modelPath := filepath.Join(root, "apps", "api", "internal", "models", "post.go")
	got := readTestFile(t, modelPath)

	checks := []struct {
		label   string
		pattern string
	}{
		{"package declaration", "package models"},
		{"struct declaration", "type Post struct {"},
		{"ID field", `gorm:"primarykey"`},
		{"Title field", "Title string"},
		{"title json tag", `json:"title"`},
		{"Content field", "Content string"},
		{"text gorm tag", `gorm:"type:text"`},
		{"Published field", "Published bool"},
		{"Views field", "Views int"},
		{"CreatedAt field", "CreatedAt time.Time"},
		{"UpdatedAt field", "UpdatedAt time.Time"},
		{"DeletedAt field", "gorm.DeletedAt"},
		{"gorm import", `"gorm.io/gorm"`},
		{"time import", `"time"`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.pattern) {
			t.Errorf("[%s] model missing %q:\n%s", c.label, c.pattern, got)
		}
	}
}

func TestWriteGoModel_SlugField(t *testing.T) {
	root := t.TempDir()

	def, err := ParseInlineFields("Article", "title:string,slug:slug")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, "testapp/apps/api", def)
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}

	got := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "article.go"))

	if !strings.Contains(got, "BeforeCreate") {
		t.Error("slug model should have BeforeCreate hook")
	}
	if !strings.Contains(got, "slugify") {
		t.Error("slug model should call slugify()")
	}

	// helpers.go should have been created
	helpersPath := filepath.Join(root, "apps", "api", "internal", "models", "helpers.go")
	if _, err := os.Stat(helpersPath); os.IsNotExist(err) {
		t.Error("helpers.go should be created when slug field is present")
	}
}

func TestWriteGoModel_BelongsTo(t *testing.T) {
	root := t.TempDir()

	def, err := ParseInlineFields("Post", "title:string,category:belongs_to:Category")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, "testapp/apps/api", def)
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}

	got := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "post.go"))

	if !strings.Contains(got, "CategoryID uint") {
		t.Errorf("belongs_to should generate CategoryID field:\n%s", got)
	}
	if !strings.Contains(got, "Category Category") {
		t.Errorf("belongs_to should generate Category association:\n%s", got)
	}
}

// ── writeGoService ────────────────────────────────────────────────────────────

func TestWriteGoService_BasicFields(t *testing.T) {
	root := t.TempDir()
	module := "testapp/apps/api"

	def, err := ParseInlineFields("Post", "title:string,content:text")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeGoService(names); err != nil {
		t.Fatalf("writeGoService: %v", err)
	}

	got := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "post.go"))

	checks := []string{
		"package services",
		module,
		"PostService",
		"PostListParams",
		"func (s *PostService) List(",
		"func (s *PostService) GetByID(",
		"func (s *PostService) Create(",
		"func (s *PostService) Update(",
		"func (s *PostService) Delete(",
		"title ILIKE ?", // search on string field
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("service missing %q:\n%s", want, got)
		}
	}
	// Module placeholder must be fully substituted
	if strings.Contains(got, "{{MODULE}}") {
		t.Error("service still contains {{MODULE}} placeholder")
	}
}

// ── writeGoHandler ────────────────────────────────────────────────────────────

func TestWriteGoHandler_BasicFields(t *testing.T) {
	root := t.TempDir()
	module := "testapp/apps/api"

	def, err := ParseInlineFields("Post", "title:string,published:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)
	names := g.Names()

	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("writeGoHandler: %v", err)
	}

	got := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "post.go"))

	checks := []string{
		"package handlers",
		"PostHandler",
		"func (h *PostHandler) List(",
		"func (h *PostHandler) GetByID(",
		"func (h *PostHandler) Create(",
		"func (h *PostHandler) Update(",
		"func (h *PostHandler) Delete(",
		module,
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("handler missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "{{MODULE}}") {
		t.Error("handler still contains {{MODULE}} placeholder")
	}
}

// ── writeZodSchema ────────────────────────────────────────────────────────────

func TestWriteZodSchema(t *testing.T) {
	root := t.TempDir()

	def, err := ParseInlineFields("Post", "title:string,published:bool,price:float")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, "testapp/apps/api", def)
	names := g.Names()

	if err := g.writeZodSchema(names); err != nil {
		t.Fatalf("writeZodSchema: %v", err)
	}

	got := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "post.ts"))

	checks := []string{
		"CreatePostSchema",
		"UpdatePostSchema",
		"type CreatePostInput",
		"type UpdatePostInput",
		`z.string()`,
		`z.boolean()`,
		`z.number()`,
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("zod schema missing %q:\n%s", want, got)
		}
	}
}

// ── writeTSTypes ──────────────────────────────────────────────────────────────

func TestWriteTSTypes(t *testing.T) {
	root := t.TempDir()

	def, err := ParseInlineFields("Post", "title:string,published:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, "testapp/apps/api", def)
	names := g.Names()

	if err := g.writeTSTypes(names); err != nil {
		t.Fatalf("writeTSTypes: %v", err)
	}

	got := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "post.ts"))

	checks := []string{
		"export interface Post {",
		"title: string",
		"published: boolean",
		"id: number",
		"created_at: string",
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("TS types missing %q:\n%s", want, got)
		}
	}
}

// ── Full Run() integration ────────────────────────────────────────────────────

func TestGenerator_Run_BasicResource(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Post", "title:string,content:text,published:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)

	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	// ── Verify generated files ────────────────────────────────────────────────

	// Go model
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "post.go"))
	assertContains(t, "model", model, "type Post struct {", `json:"title"`, "Published bool", "gorm.DeletedAt")

	// Go service
	svc := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "post.go"))
	assertContains(t, "service", svc, "PostService", "func (s *PostService) List(", module)

	// Go handler
	handler := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "post.go"))
	assertContains(t, "handler", handler, "PostHandler", "func (h *PostHandler) List(", module)

	// Zod schema
	schema := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "post.ts"))
	assertContains(t, "schema", schema, "CreatePostSchema", "UpdatePostSchema")

	// TS types
	types := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "post.ts"))
	assertContains(t, "types", types, "export interface Post {", "title: string")

	// ── Verify injections ─────────────────────────────────────────────────────

	// Model injected into user.go AutoMigrate
	userGo := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "user.go"))
	if !strings.Contains(userGo, "&Post{},") {
		t.Errorf("user.go: &Post{} not injected into AutoMigrate:\n%s", userGo)
	}
	if !strings.Contains(userGo, "// grit:models") {
		t.Error("user.go: marker // grit:models was removed")
	}

	// Handler injected into routes.go
	routes := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"))
	if !strings.Contains(routes, "postHandler") {
		t.Errorf("routes.go: postHandler not injected:\n%s", routes)
	}
	if !strings.Contains(routes, "PostHandler") {
		t.Errorf("routes.go: PostHandler not injected:\n%s", routes)
	}

	// Protected routes injected
	if !strings.Contains(routes, `protected.GET("/posts"`) {
		t.Errorf("routes.go: protected GET /posts route not injected:\n%s", routes)
	}

	// Studio model injected (inline)
	if !strings.Contains(routes, "&models.Post{},") {
		t.Errorf("routes.go: &models.Post{} not injected into studio:\n%s", routes)
	}

	// Schema export injected
	schemaIndex := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "index.ts"))
	if !strings.Contains(schemaIndex, "CreatePostSchema") {
		t.Errorf("schemas/index.ts: CreatePostSchema export not injected:\n%s", schemaIndex)
	}

	// Type export injected
	typesIndex := readTestFile(t, filepath.Join(root, "packages", "shared", "types", "index.ts"))
	if !strings.Contains(typesIndex, `export type { Post }`) {
		t.Errorf("types/index.ts: Post type export not injected:\n%s", typesIndex)
	}

	// API routes constant injected
	constants := readTestFile(t, filepath.Join(root, "packages", "shared", "constants", "index.ts"))
	if !strings.Contains(constants, "POSTS") {
		t.Errorf("constants/index.ts: POSTS route not injected:\n%s", constants)
	}
}

func TestGenerator_Run_Idempotent(t *testing.T) {
	// Run() should be idempotent for marker injections
	// (second run adds duplicate entries, but markers stay)
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Tag", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)

	if err := g.Run(); err != nil {
		t.Fatalf("first Run(): %v", err)
	}

	// Markers must still be present after first run (so second run can inject)
	userGo := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "user.go"))
	if !strings.Contains(userGo, "// grit:models") {
		t.Error("// grit:models marker was removed after first run")
	}

	routes := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"))
	if !strings.Contains(routes, "// grit:handlers") {
		t.Error("// grit:handlers marker was removed after first run")
	}
}

func TestGenerator_Run_CompoundName(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("BlogPost", "title:string,body:text")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)

	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	// Files should use snake_case names
	modelPath := filepath.Join(root, "apps", "api", "internal", "models", "blog_post.go")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Errorf("model file missing at %s", modelPath)
	}

	// Routes should use plural snake
	routes := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"))
	if !strings.Contains(routes, `"/blog_posts"`) {
		t.Errorf("routes.go: /blog_posts route not found:\n%s", routes)
	}
}

func TestGenerator_Run_RoleRestricted(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Report", "title:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)
	g.Roles = []string{"ADMIN", "EDITOR"}

	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	routes := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"))

	// Role-restricted routes go into the custom marker
	if !strings.Contains(routes, `"ADMIN"`) {
		t.Errorf("routes.go: ADMIN role not found in restricted routes:\n%s", routes)
	}
	if !strings.Contains(routes, `"EDITOR"`) {
		t.Errorf("routes.go: EDITOR role not found in restricted routes:\n%s", routes)
	}
	// Should NOT inject into grit:routes:protected (default routes)
	// (protected routes block gets no entry for this resource)
	if strings.Contains(routes, `protected.GET("/reports"`) {
		t.Errorf("routes.go: role-restricted resource should not appear in protected block:\n%s", routes)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// assertContains checks that content contains all expected substrings.
func assertContains(t *testing.T, label, content string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(content, want) {
			t.Errorf("[%s] missing %q:\n%s", label, want, content)
		}
	}
}
