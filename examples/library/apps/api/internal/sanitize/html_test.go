package sanitize

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type post struct {
	ID    uint
	Body  string  `gorm:"type:text" sanitize:"html"`
	Intro *string `sanitize:"html"`
	Title string
}

const attack = "<img src=x onerror=alert(1)><p>kept</p>"

func TestPolicy(t *testing.T) {
	for in, want := range map[string]string{
		attack:                                            "<img src=\"x\"><p>kept</p>",
		"<script>alert(1)</script>":                       "",
		"<a href=\"javascript:alert(1)\">x</a>":           "x",
		"<p style=\"position: fixed\">x</p>":              "<p>x</p>",
		"<p style=\"text-align: center\">c</p>":           "<p style=\"text-align: center\">c</p>",
		"<span style=\"color: #958DF1\">c</span>":         "<span style=\"color: #958DF1\">c</span>",
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
