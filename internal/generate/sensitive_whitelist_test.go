package generate

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A :sensitive column must not be reachable from the query string.
//
// grit#306: the model refused to marshal it and the rules engine refused to
// resolve it, and then the generator wrote it into the service's Sortable and
// Filterable maps anyway, so `?sort=key_hash` ordered rows by a secret and
// `?key_hash=<guess>` answered one row or none. :encrypted was excluded from
// those maps already, which is what made the inconsistency visible: two
// markers that both mean "this is a secret", enforced in different places.
//
// The five whitelists are asserted together on purpose. Each one is a separate
// write path to the same SQL, and a test per path passes while the set is
// broken.
func TestSensitiveColumnIsInNoQueryWhitelist(t *testing.T) {
	const module = "vault/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "ProjectKey",
		"name:string,key_hash:string:unique:sensitive,key_prefix:string,uses:int"))
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("model: %v", err)
	}
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeGoService(names); err != nil {
		t.Fatalf("service: %v", err)
	}

	api := filepath.Join(root, "apps", "api", "internal")
	service := readTestFile(t, filepath.Join(api, "services", "project_key.go"))

	// The whole list config, which is the thing a request reaches.
	cfg := service
	if i := strings.Index(service, "ListConfig = paginate.Config{"); i >= 0 {
		if j := strings.Index(service[i:], "\n}"); j >= 0 {
			cfg = service[i : i+j]
		}
	}
	if strings.Contains(cfg, "key_hash") {
		t.Errorf("the sensitive column is in the service list config:\n%s", cfg)
	}
	// The columns that are not secret are still there, or the fix was a
	// whitelist that excludes everything.
	for _, want := range []string{`"name"`, `"key_prefix"`} {
		if !strings.Contains(cfg, want) {
			t.Errorf("%s is missing from the list config, so the exclusion took too much:\n%s", want, cfg)
		}
	}

	// Each builder separately, so a failure names the path.
	if strings.Contains(g.buildHandlerSearchCols(), "key_hash") {
		t.Error("a sensitive column is searchable")
	}
	if strings.Contains(g.buildSortableSet(), "key_hash") {
		t.Error("a sensitive column is in the sortable set")
	}
	p := g.crud(names)
	if strings.Contains(p.sortCols, "key_hash") {
		t.Error("a sensitive column is sortable")
	}
	if strings.Contains(p.filterCols, "key_hash") {
		t.Error("a sensitive column is filterable")
	}

	// `grit generate field` is the fourth write path to the same maps, and it
	// had the same hole: a column added later was sortable and filterable
	// even though the same column declared at generate time was not.
	secret := mustFields(t, "ProjectKey", "recovery_code:string:sensitive").Fields[0]
	if sortableInList(secret) {
		t.Error("grit generate field would add a sensitive column to Sortable")
	}

	// And the public endpoint, which holds the stronger invariant: a column
	// held back from the response must not be reachable through a query
	// string either. publishableByDefault decided on the column's NAME, so a
	// sensitive column whose name matched none of its needles was published.
	for _, f := range g.Definition.Fields {
		if f.Sensitive && publishableByDefault(f) {
			t.Error("a sensitive column is published on the public endpoint by default")
		}
	}
}

// The same assertions for :encrypted, so the pair cannot drift apart again in
// the other direction.
func TestEncryptedColumnIsInNoQueryWhitelist(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Patient", "name:string,notes:text:encrypted"))
	names := g.Names()

	if err := g.writeGoService(names); err != nil {
		t.Fatalf("service: %v", err)
	}
	service := readTestFile(t, filepath.Join(filepath.Join(root, "apps", "api", "internal"), "services", "patient.go"))
	if regexp.MustCompile(`(?s)ListConfig = paginate\.Config\{.*?\n\}`).FindString(service) == "" {
		t.Fatal("no list config in the generated service")
	}
	cfg := regexp.MustCompile(`(?s)ListConfig = paginate\.Config\{.*?\n\}`).FindString(service)
	if strings.Contains(cfg, "notes") {
		t.Errorf("the encrypted column is in the service list config:\n%s", cfg)
	}
	for _, f := range g.Definition.Fields {
		if f.Encrypted && publishableByDefault(f) {
			t.Error("an encrypted column is published on the public endpoint by default")
		}
	}
}
