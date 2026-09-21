package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestCountryCodesAreISO3166(t *testing.T) {
	if len(CountryCodes) != 249 {
		t.Fatalf("CountryCodes has %d entries, ISO 3166-1 assigns 249", len(CountryCodes))
	}
	if !sort.StringsAreSorted(CountryCodes) {
		t.Error("CountryCodes is not sorted")
	}
	seen := map[string]bool{}
	for _, c := range CountryCodes {
		if !regexp.MustCompile("^[A-Z]{2}$").MatchString(c) || seen[c] {
			t.Errorf("bad or repeated code %q", c)
		}
		seen[c] = true
	}
	if !IsCountryCode("ug") || IsCountryCode("XX") {
		t.Error("IsCountryCode is wrong about UG or XX")
	}
}

// The API, the shared schema and the admin picker carry the same list, or the
// form offers a country the API refuses.
func TestEveryCopyOfTheCountryListAgrees(t *testing.T) {
	want := countryCodesLiteral("")
	codes := regexp.MustCompile(`"([A-Z]{2})",`)
	for name, src := range map[string]string{
		"internal/fieldtypes": apiFieldTypesCountriesGo(),
		"shared schema":       sharedFieldFormatsSchema(),
		"admin lib":           adminCountriesLib(),
	} {
		got := codes.FindAllStringSubmatch(src, -1)
		if len(got) != len(CountryCodes) {
			t.Errorf("%s has %d codes, want %d", name, len(got), len(CountryCodes))
			continue
		}
		for i, m := range got {
			if m[1] != CountryCodes[i] {
				t.Errorf("%s: code %d is %s, want %s", name, i, m[1], CountryCodes[i])
				break
			}
		}
	}
	if !strings.Contains(want, `"ZW",`) {
		t.Error("the literal lost its last code")
	}
}

// A database.go from before the field types gets the call and the import, and
// ends up exactly as a new project's is. Twice is the same as once.
func TestEnsureFieldTypesWiringMatchesTheTemplate(t *testing.T) {
	const module = "acme/apps/api"
	fresh := strings.ReplaceAll(apiDatabaseGo(), "{{MODULE}}", module)
	old := strings.Replace(fresh, fieldTypesConnectHook, "", 1)
	old = strings.Replace(old, "\t\""+module+"/internal/fieldtypes\"\n", "", 1)
	if old == fresh || strings.Contains(old, "fieldtypes") {
		t.Fatal("could not reconstruct the old database.go")
	}

	apiRoot := t.TempDir()
	path := filepath.Join(apiRoot, "internal", "database", "database.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := EnsureFieldTypesWiring(apiRoot, module); err != nil {
			t.Fatalf("EnsureFieldTypesWiring: %v", err)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != fresh {
		t.Errorf("the repaired database.go differs from the template")
	}
	if n := strings.Count(string(got), "fieldtypes.Install(db)"); n != 1 {
		t.Errorf("fieldtypes.Install appears %d times", n)
	}
}

func TestFieldTypesShipWithEveryAPI(t *testing.T) {
	db := apiDatabaseGo()
	for _, want := range []string{"fieldtypes.Install(db)", "{{MODULE}}/internal/fieldtypes"} {
		if !strings.Contains(db, want) {
			t.Errorf("database.go template is missing %q", want)
		}
	}
	files := fieldTypesFiles("api")
	if len(files) != 4 {
		t.Errorf("internal/fieldtypes has %d files", len(files))
	}
	for path, body := range files {
		if strings.Contains(body, "—") {
			t.Errorf("%s has an em dash", path)
		}
	}
	if !strings.Contains(apiRespondGo(), "Validation(c, fields.Error(), fields.FieldErrors())") {
		t.Error("respond.WriteError does not answer a FieldErrors error with its details")
	}
	// internal/phone is not in every project: its module is only fetched for a
	// tel field.
	if strings.Contains(apiGoMod(tripleOptions()), PhoneNumbersModule) {
		t.Error("every project's go.mod requires the phonenumbers module")
	}
}

func TestAdminShipsTheFieldInputs(t *testing.T) {
	root := t.TempDir()
	for name, files := range map[string]map[string]string{
		"Next": adminFileMap(root, tripleOptions()),
		"Vite": adminTanStackFileMap(root, viteTripleOptions()),
	} {
		var builder, cells, pkg string
		found := map[string]bool{}
		for path, body := range files {
			slash := filepath.ToSlash(path)
			switch {
			case strings.HasSuffix(slash, "components/forms/form-builder.tsx"):
				builder = body
			case strings.HasSuffix(slash, "components/tables/cell-renderers.tsx"):
				cells = body
			case strings.HasSuffix(slash, "apps/admin/package.json"):
				pkg = body
			}
			for _, f := range []string{"phone-field.tsx", "country-select.tsx", "country-field.tsx", "rating-field.tsx", "json-field.tsx", "color-field.tsx", "percent-field.tsx", "format-text-field.tsx", "phone-cell.tsx", "lib/countries.ts", "lib/field-formats.ts"} {
				if strings.HasSuffix(slash, f) {
					found[f] = true
				}
			}
			if strings.Contains(body, "—") && (strings.Contains(slash, "field-formats") || strings.Contains(slash, "countries.ts") || strings.Contains(slash, "phone-")) {
				t.Errorf("%s: %s has an em dash", name, slash)
			}
		}
		if len(found) != 11 {
			t.Errorf("%s admin is missing field inputs: has %v", name, found)
		}
		for _, typ := range []string{"email", "url", "domain", "time", "tel", "country", "color", "percent", "rating", "json"} {
			if !strings.Contains(builder, `case "`+typ+`":`) {
				t.Errorf("%s form builder has no case for %s", name, typ)
			}
		}
		for _, format := range []string{"domain", "tel", "country", "percent", "rating", "time", "json"} {
			if !strings.Contains(cells, `case "`+format+`":`) {
				t.Errorf("%s cell renderers have no case for %s", name, format)
			}
		}
		for _, dep := range []string{`"@base-ui/react"`, `"libphonenumber-js"`} {
			if !strings.Contains(pkg, dep) {
				t.Errorf("%s admin package.json lacks %s", name, dep)
			}
		}
	}
	if !strings.Contains(sharedPackageJSON(tripleOptions()), `"libphonenumber-js"`) {
		t.Error("packages/shared does not depend on libphonenumber-js, which phone.ts imports")
	}
	if !strings.Contains(webAdminDependencies(doubleOptions()), `"@base-ui/react"`) {
		t.Error("the panel embedded in a web app lacks @base-ui/react")
	}
}

// A project whose shared package predates the formatted types gets the two
// schema modules, and libphonenumber-js once a tel field needs it.
func TestWriteSharedFieldFormats(t *testing.T) {
	shared := t.TempDir()
	pkg := filepath.Join(shared, "package.json")
	if err := os.WriteFile(pkg, []byte("{\n  \"dependencies\": {\n    \"zod\": \"^3.22.0\"\n  }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	added, err := WriteSharedFieldFormats(shared, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(added) != 3 {
		t.Errorf("added %v", added)
	}
	body, _ := os.ReadFile(pkg)
	if !strings.Contains(string(body), `"libphonenumber-js": "`+libphonenumberVersion+`",`) {
		t.Errorf("package.json = %s", body)
	}
	again, err := WriteSharedFieldFormats(shared, true)
	if err != nil || len(again) != 0 {
		t.Errorf("a second run added %v, %v", again, err)
	}
}
