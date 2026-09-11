package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A catalogue from before a key existed never got it, because grit add i18n
// never overwrites a catalogue: the admin asked for nav.systemHub and showed
// English in the middle of a French page.
func TestMergeCatalogueAddsOnlyWhatIsMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "messages", "fr.json")
	mine := `{
  "nav": { "dashboard": "Mon tableau", "custom": "À moi" },
  "resources": { "purchase-requests": { "plural": "Demandes d'achat" } }
}
`
	writeTestFile(t, path, mine)

	added, err := mergeCatalogue(path, i18nMessagesFR())
	if err != nil || !added {
		t.Fatalf("mergeCatalogue: added=%v err=%v", added, err)
	}
	got, _ := os.ReadFile(path)
	var doc map[string]map[string]interface{}
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatalf("the merged catalogue is not valid JSON: %v\n%s", err, got)
	}
	if doc["nav"]["dashboard"] != "Mon tableau" || doc["nav"]["custom"] != "À moi" {
		t.Errorf("an existing translation was changed: %v", doc["nav"])
	}
	if doc["nav"]["systemHub"] != "Centre système" {
		t.Errorf("a missing key was not added: %v", doc["nav"])
	}
	if _, ok := doc["table"]["searchPlaceholder"]; !ok {
		t.Error("a missing section was not added")
	}
	if _, ok := doc["resources"]; !ok {
		t.Error("a section the template does not have was dropped")
	}
	// Order is kept: the user's first key is still first.
	if strings.Index(string(got), `"nav"`) > strings.Index(string(got), `"resources"`) {
		t.Error("the keys were reordered")
	}

	// A second merge has nothing to add.
	if added, _ := mergeCatalogue(path, i18nMessagesFR()); added {
		t.Error("merging twice changed the file again")
	}
}

func TestMergeCatalogueLeavesInvalidJSONAlone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "messages", "en.json")
	writeTestFile(t, path, "{ not json")
	if added, err := mergeCatalogue(path, i18nMessagesEN()); added || err != nil {
		t.Fatalf("an invalid catalogue was touched: added=%v err=%v", added, err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "{ not json" {
		t.Error("an invalid catalogue was rewritten")
	}
}
