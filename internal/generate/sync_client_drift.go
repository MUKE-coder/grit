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
// grit sync keeps the shared TypeScript type, the Zod schema and the admin
// resource definition in step with the model, and since grit#87 the desktop
// columns through their marker fences and the Expo form through a rebuild. The
// Expo list is the one left, and this reports it.
//
// What that report is for: the failure it came from was a half-synced project
// that said nothing about it. Adding a column to a Go model and running sync
// gave a correct type, a correct schema, a correct admin screen, and a mobile
// create form still showing the original three fields. Verified on an
// emulator: the field was in the API response and absent from the form.
//
// So a screen this cannot update says so and names the command that resolves
// it, which is the difference between a known limitation and a bug somebody
// finds in a demo. It also still covers the form in the one case a rebuild is
// refused, which is a form somebody has edited.
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
	fmt.Println("    Types, schemas and the admin are current, and so is a desktop")
	fmt.Println("    screen carrying the grit:cols markers: sync adds new columns to")
	fmt.Println("    those in place.")
	fmt.Println("    An Expo form is rebuilt from the model, unless you have edited it,")
	fmt.Println("    in which case it is named above and left alone.")
	fmt.Println("    An Expo list is not: it renders a header block and a row block that")
	fmt.Println("    have to line up under a summed table width, so adding a column in")
	fmt.Println("    place would be three things to keep in step and a misaligned table")
	fmt.Println("    when one slips.")
	fmt.Println()

	for _, resource := range order {
		for _, d := range byResource[resource] {
			fmt.Printf("    %s\n      missing: %s\n", d.File, strings.Join(d.Missing, ", "))
		}
		fmt.Printf("      to rebuild them from the model, which no longer means retyping the\n")
		fmt.Printf("      field list:\n")
		fmt.Printf("        grit generate resource %s --force\n", resource)
		fmt.Println("      it prints the field list it read before it uses it. That rewrites")
		fmt.Println("      the screens and discards any edits you made to them, so add the")
		fmt.Println("      field by hand instead if you have customised those files.")
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

// syncMobileForm rebuilds a resource's Expo create/edit form from its model.
//
// Shared by `grit sync` and `grit generate field`, because a field added by
// either has to reach the same screen and two writers for one file is how the
// two come to disagree about what it looks like.
//
// Returns the path to report as updated, or the path to report as kept, never
// both. Both empty means there was nothing to do: no Expo app, or no form for
// this resource.
//
// Whether a rebuild is allowed is the manifest's answer, not a new rule: it
// holds the hash of what Grit last wrote, so a form nobody has touched is
// Grit's to rewrite and one somebody has edited is theirs to keep. The same
// rule grit upgrade applies to every framework-owned file.
func syncMobileForm(root string, def *ResourceDefinition, names Names) (updated, kept string) {
	// Only Definition is read by mobileFormSource, and Root by the two path
	// helpers. A full Generator is not available on the sync path, which has a
	// parsed model and no flags.
	g := &Generator{Root: root, Definition: def}
	if !g.hasMobileApp() {
		return "", ""
	}
	path := filepath.Join(g.mobileRoot(), "components", "resource-forms",
		names.PluralKebab+"-form.tsx")
	if _, err := os.Stat(path); err != nil {
		return "", ""
	}

	// Three different reasons not to write, and they are worth telling apart:
	// "you have edited it" is advice, and saying it about a file Grit has no
	// record of is a guess presented as fact. A project older than the
	// manifest, or one whose form predates it, lands on Untracked.
	label := relFromRoot(root, path)
	switch manifest.StatusOfPath(path) {
	case manifest.Unchanged:
	case manifest.Untracked:
		return "", label + "  (Grit has no record of writing it, so it is left alone)"
	default:
		return "", label + "  (you have edited it)"
	}

	body := g.mobileFormSource(names)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		// Not fatal for either caller. On the generate-field path the field is
		// already in the model, the schema and the type by now, so failing
		// here would leave the command half done with nothing to resume it.
		return "", label + "  (could not be written: " + err.Error() + ")"
	}
	manifest.Note(path, body)
	return label, ""
}
