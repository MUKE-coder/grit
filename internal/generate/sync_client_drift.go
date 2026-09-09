package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// clientUIDrift reports model fields that never reached a client's generated UI.
//
// grit sync keeps three things in step: the shared TypeScript type, the Zod
// schema, and the admin resource definition, which SyncAdminResource updates
// through markers so customised entries survive. The mobile forms and screens
// and the desktop columns have no equivalent, so they are frozen at the moment
// the resource was generated.
//
// The result is a half-synced project that says nothing about it. Adding a
// column to a Go model and running sync gave a correct type, a correct schema,
// a correct admin screen, and a mobile create form still showing the original
// three fields. Verified on an emulator: the field was in the API response and
// absent from the form.
//
// This does not fix the drift, which needs the same marker-based regeneration
// the admin already has. It says the drift exists and names the command that
// resolves it, which is the difference between a known limitation and a bug
// somebody finds in a demo.
type clientUIDrift struct {
	Client   string   // "mobile" or "desktop"
	Resource string   // Pascal name
	File     string   // path relative to the project root, for the message
	Missing  []string // json field names the file never mentions
}

// clientUIFiles are the generated files that render a resource's fields, and
// that sync cannot currently update.
//
// The admin is deliberately absent: SyncAdminResource keeps it current.
func clientUIFiles(root string, s GoStruct) []struct{ client, path string } {
	pluralKebab := strings.ReplaceAll(Pluralize(toSnakeCase(s.Name)), "_", "-")
	return []struct{ client, path string }{
		{"mobile", filepath.Join("apps", "expo", "components", "resource-forms", pluralKebab+"-form.tsx")},
		{"mobile", filepath.Join("apps", "expo", "app", pluralKebab, "index.tsx")},
		{"desktop", filepath.Join("apps", "desktop", "frontend", "src", "routes", "app", pluralKebab+".index.tsx")},
	}
}

// findClientUIDrift looks for model fields no generated client file mentions.
//
// Matching on the JSON name is deliberately loose: a file that mentions the
// field at all, in a column, a form input or a filter, is treated as knowing
// about it. The question worth answering is "was this field never generated
// into the UI", and a loose match answers it without guessing at each client's
// shape. False negatives are the right way to be wrong here, since a spurious
// warning about a field somebody deliberately left out is worse than silence.
func findClientUIDrift(root string, s GoStruct) []clientUIDrift {
	// Only a resource the generator wrote can be rebuilt by the command this
	// message suggests. Built-ins like Blog have hand-written screens that
	// were never derived from the model, so --force would be wrong advice and
	// the "missing" fields are not missing, they were never meant to be there.
	if !generatedResource(root, s.Name) {
		return nil
	}

	var out []clientUIDrift

	for _, target := range clientUIFiles(root, s) {
		abs := filepath.Join(root, target.path)
		data, err := os.ReadFile(abs)
		if err != nil {
			continue // this project does not have that client
		}
		body := string(data)

		var missing []string
		for _, f := range s.Fields {
			if isAutoField(f.JSONName) || frameworkColumn(f.JSONName) ||
				f.JSONName == "" || f.JSONName == "-" {
				continue
			}
			if !strings.Contains(body, f.JSONName) {
				missing = append(missing, f.JSONName)
			}
		}
		if len(missing) > 0 {
			out = append(out, clientUIDrift{
				Client:   target.client,
				Resource: s.Name,
				File:     filepath.ToSlash(target.path),
				Missing:  missing,
			})
		}
	}
	return out
}

// printClientUIDrift explains the drift and what to do about it.
func printClientUIDrift(drifts []clientUIDrift) {
	if len(drifts) == 0 {
		return
	}

	// One line per file, grouped by resource so the fix command appears once
	// per resource rather than once per file.
	byResource := map[string][]clientUIDrift{}
	var order []string
	for _, d := range drifts {
		if _, seen := byResource[d.Resource]; !seen {
			order = append(order, d.Resource)
		}
		byResource[d.Resource] = append(byResource[d.Resource], d)
	}

	fmt.Println("\n  ⚠ Some generated screens do not know about every field.")
	fmt.Println("    Types, schemas and the admin are current. Mobile and desktop")
	fmt.Println("    screens are only written when the resource is generated, and")
	fmt.Println("    neither `grit sync` nor `grit generate field` updates them.")
	fmt.Println()

	for _, resource := range order {
		for _, d := range byResource[resource] {
			fmt.Printf("    %s\n      missing: %s\n", d.File, strings.Join(d.Missing, ", "))
		}
		fmt.Printf("      to rebuild them, regenerate the resource with its full field list:\n")
		fmt.Printf("        grit generate resource %s --fields \"...\" --force\n", resource)
		fmt.Println("      that rewrites the screens and discards any edits you made to them,")
		fmt.Println("      so add the field to those files by hand if you have customised them.")
		fmt.Println()
	}
}

// frameworkColumn reports a column the framework owns, which no generated
// screen is expected to render.
//
// Reporting these was most of the first version's output, and a warning that
// is mostly wrong is one people learn to scroll past.
func frameworkColumn(jsonName string) bool {
	switch jsonName {
	case "version", // optimistic concurrency
		"archived_at",               // archive state, handled by its own control
		"org_id",                    // multitenant scoping, never user-editable
		"path", "depth", "position": // --tree bookkeeping
		return true
	}
	return false
}

// generatedResource reports whether the manifest attributes this resource's
// model to `grit generate resource`.
func generatedResource(root, pascal string) bool {
	m, err := manifest.Load(root)
	if err != nil || len(m.Files) == 0 {
		// No manifest: a project older than v3.147.0. Say nothing rather than
		// guess, since a wrong --force suggestion destroys hand-written files.
		return false
	}
	want := "resource:" + pascal
	for _, entry := range m.Files {
		if entry.Generator == want {
			return true
		}
	}
	return false
}
