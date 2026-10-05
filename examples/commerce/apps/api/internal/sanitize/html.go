package sanitize

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
