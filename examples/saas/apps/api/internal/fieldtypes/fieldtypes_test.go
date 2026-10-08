package fieldtypes

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRules(t *testing.T) {
	good := map[string][2]string{
		"email":   {"  Ada@Example.COM ", "ada@example.com"},
		"url":     {"https://example.com/pricing?x=1", "https://example.com/pricing?x=1"},
		"domain":  {"https://www.Example.co.ug/about", "www.example.co.ug"},
		"country": {"ug", "UG"},
		"color":   {"#ABC", "#aabbcc"},
		"time":    {"9:05:30", "09:05"},
	}
	for format, c := range good {
		got, err := Normalize(format, c[0])
		if err != nil || got != c[1] {
			t.Errorf("%s(%q) = %v, %v; want %q", format, c[0], got, err, c[1])
		}
	}
	if got, err := Normalize("domain", "münchen.de"); err != nil || got != "xn--mnchen-3ya.de" {
		t.Errorf("an IDN is stored as punycode: got %v, %v", got, err)
	}

	bad := map[string]string{
		"email":   "Ada <ada@example.com>",
		"url":     "javascript:alert(1)",
		"domain":  "not a domain",
		"country": "XX",
		"color":   "blue",
		"time":    "25:00",
		"percent": "101",
		"rating":  "4.5",
		"json":    "",
	}
	for format, v := range bad {
		var value any = v
		if format == "json" {
			value = datatypes.JSON("{nope")
		}
		if _, err := Normalize(format, value); err == nil {
			t.Errorf("%s accepted %q", format, v)
		}
	}
	if _, err := Normalize("rating:10", 7.0); err != nil {
		t.Errorf("rating:10 refused 7: %v", err)
	}
	if _, err := Normalize("rating:5", 7.0); err == nil {
		t.Error("rating:5 accepted 7")
	}
	if got, _ := Normalize("percent", 12.345); got != 12.35 {
		t.Errorf("percent rounds to two decimals: got %v", got)
	}
}

func TestSamplesPassTheirRules(t *testing.T) {
	seen := map[string]bool{}
	for n := 0; n < 5000; n++ {
		values := map[string]any{
			"email":     SampleEmail("Ada", "Lovelace", "example.com", n),
			"url":       SampleURL("example.com", n),
			"domain":    SampleDomain("acme", n),
			"country":   SampleCountry(n),
			"color":     SampleColor(n),
			"time":      SampleTime(n),
			"percent":   SamplePercent(n),
			"rating:5":  SampleRating(n, 5),
			"rating:10": SampleRating(n, 10),
			"json":      SampleJSON(n),
		}
		for format, v := range values {
			if _, err := Normalize(format, v); err != nil {
				t.Fatalf("sample %d for %s (%v) fails its own rule: %v", n, format, v, err)
			}
		}
		for _, unique := range []string{values["email"].(string), values["url"].(string), values["domain"].(string), values["color"].(string)} {
			if seen[unique] {
				t.Fatalf("sample %d repeats %q", n, unique)
			}
			seen[unique] = true
		}
	}
}

type contact struct {
	ID      uint
	Email   string         `json:"email" format:"email"`
	Color   string         `json:"color" format:"color"`
	Score   int            `json:"score" format:"rating:5"`
	Meta    datatypes.JSON `json:"meta" format:"json"`
	Comment string
}

// Every way a generated API writes a row.
func TestEveryWriteIsChecked(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := Install(db); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := db.AutoMigrate(&contact{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	row := contact{Email: "Ada@Example.com", Color: "#ABC", Score: 4}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	var stored contact
	db.First(&stored, row.ID)
	if stored.Email != "ada@example.com" || stored.Color != "#aabbcc" {
		t.Errorf("create stored %q and %q", stored.Email, stored.Color)
	}

	var fe *Error
	err = db.Create(&contact{Email: "nope"}).Error
	if !errors.As(err, &fe) || fe.Field != "email" {
		t.Fatalf("an invalid create was not refused with the field: %v", err)
	}

	if err := db.Model(&row).Updates(map[string]interface{}{"color": "FFF"}).Error; err != nil {
		t.Fatalf("Updates(map): %v", err)
	}
	db.First(&stored, row.ID)
	if stored.Color != "#ffffff" {
		t.Errorf("Updates(map) stored %q", stored.Color)
	}
	if err := db.Model(&row).Updates(map[string]interface{}{"score": 9}).Error; !errors.As(err, &fe) {
		t.Errorf("Updates(map) accepted 9 stars out of 5: %v", err)
	}
	if err := db.Model(&row).Update("email", "bad").Error; !errors.As(err, &fe) {
		t.Errorf("Update(column) accepted a bad email: %v", err)
	}
	// A PATCH body holds decoded JSON, which is stored as JSON text.
	if err := db.Model(&row).Updates(map[string]interface{}{"meta": map[string]interface{}{"plan": "pro"}}).Error; err != nil {
		t.Fatalf("Updates(map) with a decoded object: %v", err)
	}
	db.First(&stored, row.ID)
	if string(stored.Meta) != "{\"plan\":\"pro\"}" {
		t.Errorf("json stored %s", stored.Meta)
	}
	// An untagged column is not touched.
	if err := db.Model(&row).Update("comment", "anything").Error; err != nil {
		t.Errorf("an untagged column was checked: %v", err)
	}
}
