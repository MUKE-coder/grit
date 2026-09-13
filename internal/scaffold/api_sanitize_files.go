package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The HTML sanitiser internal/sanitize wraps, at the version the scaffold pins.
const (
	bluemondayModule  = "github.com/microcosm-cc/bluemonday"
	bluemondayVersion = "v1.0.27"
)

// sanitizeFiles are internal/sanitize and its test, which ship with every API.
func sanitizeFiles(apiRoot string) map[string]string {
	return map[string]string{
		filepath.Join(apiRoot, "internal", "sanitize", "html.go"):      apiSanitizeHTMLGo(),
		filepath.Join(apiRoot, "internal", "sanitize", "html_test.go"): apiSanitizeHTMLTestGo(),
	}
}

const sanitizeConnectHook = `	// Rich text is sanitised on its way into the database: every field tagged
	// sanitize:"html", through every write that has a model (create, update,
	// PATCH, bulk edit, CSV import, sync push, GORM Studio's row editor).
	if err := sanitize.Install(db); err != nil {
		return nil, fmt.Errorf("installing the HTML sanitiser: %w", err)
	}

`

// EnsureSanitizeWiring makes the project call sanitize.Install when it
// connects. Idempotent, and it refuses rather than guesses when the file has
// been reshaped past recognition.
func EnsureSanitizeWiring(apiRoot, module string) error {
	path := filepath.Join(apiRoot, "internal", "database", "database.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("the HTML sanitiser needs %s: %w", path, err)
	}
	crlf := strings.Contains(string(raw), "\r\n")
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	if strings.Contains(content, "sanitize.Install(db)") {
		return nil
	}
	const anchor = "\tsqlDB, err := db.DB()"
	if !strings.Contains(content, anchor) {
		return fmt.Errorf("could not find where to install the HTML sanitiser in %s.\n\n"+
			"Add this call after connecting, or rich text is stored as it was sent:\n\n  sanitize.Install(db)", path)
	}
	content = strings.Replace(content, anchor, sanitizeConnectHook+anchor, 1)
	var ok bool
	if content, ok = addImportGroup(content, module+"/internal/sanitize"); !ok {
		return fmt.Errorf("could not add the sanitize import to %s: add it by hand", path)
	}
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Wired the HTML sanitiser into database.go")
	return nil
}

func apiSanitizeHTMLGo() string {
	return `package sanitize

import (
	"reflect"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// HTML returns s with everything but safe markup removed.
//
// It is the policy for rich text that somebody else will read. The formatting
// the admin editor produces survives: headings, lists, links, images, tables,
// code blocks with their language, text alignment, colours and highlights.
// Scripts, event handlers, javascript: URLs, iframes, and any style beyond
// those few properties do not.
func HTML(s string) string {
	return policy.Sanitize(s)
}

var policy = newPolicy()

func newPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	color := regexp.MustCompile("^(#[0-9a-fA-F]{3,8}|rgba?\\(\\s*\\d{1,3}%?\\s*,\\s*\\d{1,3}%?\\s*,\\s*\\d{1,3}%?\\s*(,\\s*(0|1|0?\\.\\d+)\\s*)?\\)|var\\(--[\\w-]+\\))$")
	p.AllowStyles("text-align").MatchingEnum("left", "center", "right", "justify").
		OnElements("p", "h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowStyles("color").Matching(color).OnElements("span", "mark")
	p.AllowStyles("background-color").Matching(color).OnElements("mark")
	p.AllowAttrs("data-color").Matching(color).OnElements("mark")
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[\\w+#-]+$")).OnElements("code")
	p.AllowAttrs("target").Matching(regexp.MustCompile("^_blank$")).OnElements("a")
	p.RequireNoReferrerOnLinks(true)
	return p
}

// Install sanitises every field tagged sanitize:"html" on its way into the
// database. Call it once, straight after connecting.
//
// It sits in the write path rather than in handlers because there are too many
// ways in to remember: create, update, PATCH, bulk edit, the CSV importer, sync
// push and GORM Studio's row editor all write through this handle. A blog post
// stored as sent was rendered with dangerouslySetInnerHTML, so anybody allowed
// to edit posts could run script in the browser of an admin who read one.
//
// It needs the model to know the column: db.Model(&row).Updates is covered,
// db.Table("x").Updates is not, and neither is raw SQL.
func Install(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:create").Register("sanitize:html_create", sanitizeWrite); err != nil {
		return err
	}
	return db.Callback().Update().Before("gorm:update").Register("sanitize:html_update", sanitizeWrite)
}

func htmlFields(s *schema.Schema) []*schema.Field {
	var out []*schema.Field
	for _, f := range s.Fields {
		if f.Tag.Get("sanitize") == "html" {
			out = append(out, f)
		}
	}
	return out
}

func sanitizeWrite(tx *gorm.DB) {
	stmt := tx.Statement
	if stmt == nil || stmt.Schema == nil {
		return
	}
	fields := htmlFields(stmt.Schema)
	if len(fields) == 0 {
		return
	}
	switch dest := stmt.Dest.(type) {
	case map[string]interface{}:
		sanitizeMap(stmt.Schema, dest)
		return
	case *map[string]interface{}:
		if dest != nil {
			sanitizeMap(stmt.Schema, *dest)
		}
		return
	}
	// The row being written, and the value handed to Updates when that is a
	// different struct. Sanitising a value twice leaves it as it was.
	sanitizeValue(tx, fields, stmt.ReflectValue)
	if stmt.Dest != nil {
		sanitizeValue(tx, fields, reflect.ValueOf(stmt.Dest))
	}
}

// sanitizeMap covers Updates(map) and Update(column, value), which is how the
// generated update, PATCH and bulk handlers write.
func sanitizeMap(s *schema.Schema, values map[string]interface{}) {
	for key, value := range values {
		f := s.LookUpField(key)
		if f == nil || f.Tag.Get("sanitize") != "html" {
			continue
		}
		switch v := value.(type) {
		case string:
			values[key] = HTML(v)
		case *string:
			if v != nil {
				clean := HTML(*v)
				values[key] = &clean
			}
		}
	}
}

func sanitizeValue(tx *gorm.DB, fields []*schema.Field, v reflect.Value) {
	if !v.IsValid() {
		return
	}
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			sanitizeValue(tx, fields, v.Index(i))
		}
	case reflect.Struct:
		if !v.CanAddr() || v.Type() != fields[0].Schema.ModelType {
			return
		}
		ctx := tx.Statement.Context
		for _, f := range fields {
			value, zero := f.ValueOf(ctx, v)
			if zero {
				continue
			}
			var err error
			switch s := value.(type) {
			case string:
				if clean := HTML(s); clean != s {
					err = f.Set(ctx, v, clean)
				}
			case *string:
				if s != nil {
					if clean := HTML(*s); clean != *s {
						err = f.Set(ctx, v, &clean)
					}
				}
			}
			if err != nil {
				_ = tx.AddError(err)
			}
		}
	}
}
`
}

func apiSanitizeHTMLTestGo() string {
	return `package sanitize

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type post struct {
	ID    uint
	Body  string  ` + "`" + `gorm:"type:text" sanitize:"html"` + "`" + `
	Intro *string ` + "`" + `sanitize:"html"` + "`" + `
	Title string
}

const attack = "<img src=x onerror=alert(1)><p>kept</p>"

func TestPolicy(t *testing.T) {
	for in, want := range map[string]string{
		attack:                                    "<img src=\"x\"><p>kept</p>",
		"<script>alert(1)</script>":               "",
		"<a href=\"javascript:alert(1)\">x</a>":   "x",
		"<p style=\"position: fixed\">x</p>":      "<p>x</p>",
		"<p style=\"text-align: center\">c</p>":   "<p style=\"text-align: center\">c</p>",
		"<span style=\"color: #958DF1\">c</span>": "<span style=\"color: #958DF1\">c</span>",
		"<pre><code class=\"language-go\">x</code></pre>": "<pre><code class=\"language-go\">x</code></pre>",
	} {
		if got := HTML(in); got != want {
			t.Errorf("HTML(%q) = %q, want %q", in, got, want)
		}
	}
}

// Every way a generated API writes a row. Before Install a blog post was stored
// exactly as sent, and rendered raw.
func TestEveryWriteIsSanitised(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := Install(db); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := db.AutoMigrate(&post{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	column := func(name string, id uint) string {
		var s string
		db.Raw("SELECT "+name+" FROM posts WHERE id = ?", id).Scan(&s)
		return s
	}
	clean := func(what, got string) {
		t.Helper()
		if strings.Contains(got, "onerror") || !strings.Contains(got, "<p>kept</p>") {
			t.Errorf("%s stored %q", what, got)
		}
	}

	intro := attack
	row := post{Body: attack, Intro: &intro, Title: attack}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	clean("Create", column("body", row.ID))
	clean("Create (pointer field)", column("intro", row.ID))
	if got := column("title", row.ID); got != attack {
		t.Errorf("an untagged field was changed: %q", got)
	}

	writes := []struct {
		name  string
		write func() error
	}{
		{"Updates(map)", func() error { return db.Model(&row).Updates(map[string]interface{}{"body": attack}).Error }},
		{"Update(column)", func() error { return db.Model(&row).Update("body", attack).Error }},
		{"Updates(struct)", func() error { return db.Model(&row).Updates(&post{Body: attack}).Error }},
		{"Save", func() error { row.Body = attack; return db.Save(&row).Error }},
	}
	for _, w := range writes {
		if err := w.write(); err != nil {
			t.Fatalf("%s: %v", w.name, err)
		}
		clean(w.name, column("body", row.ID))
	}

	rows := []post{{Body: attack}, {Body: attack}}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("batch create: %v", err)
	}
	for _, r := range rows {
		clean("Create(slice)", column("body", r.ID))
	}
}
`
}
