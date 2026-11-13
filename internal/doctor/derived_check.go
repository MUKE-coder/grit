package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// A generated file that is older than the thing it was generated from.
//
// The manifest has always answered "have you edited this file", which is what
// lets `grit upgrade` leave your work alone. It could not answer the other
// question: is this file still what the model says. Nothing could.
//
// The gap is specific and it is not theoretical. `grit sync` regenerates the
// shared types and the Zod schemas from the models, so those stay honest. It
// does not regenerate the handler's create request, the service's writable
// column list, or the admin form. Add a column to a model by hand, as a
// person reasonably does, and those three describe the model as it was: the
// admin shows a form missing a field, the API drops it if the form sends it
// anyway, and the only symptom is that the value does not save.
//
// So the generator records what each file came out of, and this compares. The
// finding is per source rather than per file, because one edited model leaves
// a dozen files behind and a dozen findings about one cause is a wall of text
// that gets scrolled past.
func checkDerivedFiles(p *project) []Finding {
	m, err := manifest.Load(p.root)
	if err != nil || len(m.Files) == 0 {
		return nil // no manifest: a project from before one, or not a Grit project
	}

	// source path -> the files generated from it that are now behind
	behind := map[string][]string{}
	// source path -> the files generated from it that have been hand-edited,
	// which cannot simply be regenerated
	edited := map[string][]string{}
	// Hashed once per source, not once per file: a resource has a dozen.
	current := map[string]string{}

	for _, rel := range m.Entries() {
		entry := m.Files[rel]
		if entry.From == nil || entry.From.Path == "" {
			continue
		}
		if !followsTheFields(rel) {
			continue
		}
		src := entry.From.Path

		now, seen := current[src]
		if !seen {
			data, readErr := os.ReadFile(filepath.Join(p.root, filepath.FromSlash(src))) // #nosec G304
			if readErr != nil {
				// The model was deleted or renamed. That is `grit remove`'s
				// business, not this check's, and saying so here would be a
				// second complaint about one event.
				current[src] = ""
				continue
			}
			now = manifest.Hash(string(data))
			current[src] = now
		}
		if now == "" || now == entry.From.Hash {
			continue
		}

		// A file nobody has touched can be regenerated. One that has been
		// edited cannot, not without taking somebody's work with it, and the
		// finding has to say which.
		if m.StatusOf(p.root, rel) == manifest.Modified {
			edited[src] = append(edited[src], rel)
			continue
		}
		behind[src] = append(behind[src], rel)
	}

	var out []Finding
	for _, src := range sortedKeys(behind, edited) {
		stale, mine := behind[src], edited[src]
		sort.Strings(stale)
		sort.Strings(mine)

		var parts []string
		if len(stale) > 0 {
			parts = append(parts, fmt.Sprintf("%s behind it", countFiles(len(stale))))
		}
		if len(mine) > 0 {
			parts = append(parts, fmt.Sprintf("%s you have edited", countFiles(len(mine))))
		}

		out = append(out, Finding{
			Level:    "warning",
			Resource: src,
			Message: fmt.Sprintf("has changed since the files generated from it: %s. %s",
				strings.Join(parts, ", "), firstFew(append(append([]string{}, stale...), mine...))),
			Fix: derivedFix(src, len(mine) > 0),
		})
	}
	return out
}

// derivedFix names the command for this source, and says when regenerating
// would take an edit with it.
func derivedFix(src string, hasEdits bool) string {
	name := modelNameFrom(src)
	fix := "run grit sync for the shared types and schemas"
	if name != "" {
		fix += ", and grit generate resource " + name +
			" to bring the handler, service and admin screen back in line"
	}
	if hasEdits {
		fix += ". Some of those files you have edited: regenerating overwrites them, " +
			"so diff first or add the column with grit generate field instead"
	}
	return fix
}

// modelNameFrom turns internal/models/blog_post.go into BlogPost. Empty when
// the source is not a model, which is every other kind of source this may
// grow later.
func modelNameFrom(src string) string {
	src = filepath.ToSlash(src)
	i := strings.Index(src, "internal/models/")
	if i < 0 || !strings.HasSuffix(src, ".go") {
		return ""
	}
	base := strings.TrimSuffix(src[i+len("internal/models/"):], ".go")
	if base == "" || strings.Contains(base, "/") {
		return ""
	}
	var b strings.Builder
	for _, word := range strings.Split(base, "_") {
		if word == "" {
			continue
		}
		b.WriteString(strings.ToUpper(word[:1]) + word[1:])
	}
	return b.String()
}

func countFiles(n int) string {
	if n == 1 {
		return "1 file"
	}
	return fmt.Sprintf("%d files", n)
}

// firstFew names up to three, because a reader needs to recognise the place
// and does not need the list.
//
// Ordered by how much the file depends on the model rather than
// alphabetically. The first version named three admin page wrappers, which
// render a resource by name and say nothing about its fields, while the
// handler and the service were in the "and 12 more".
func firstFew(files []string) string {
	sort.SliceStable(files, func(i, j int) bool {
		ri, rj := tellingRank(files[i]), tellingRank(files[j])
		if ri != rj {
			return ri < rj
		}
		return files[i] < files[j]
	})
	if len(files) <= 3 {
		return strings.Join(files, ", ")
	}
	return fmt.Sprintf("%s and %d more", strings.Join(files[:3], ", "), len(files)-3)
}

// followsTheFields reports whether a change to the model can leave this file
// saying something untrue.
//
// Generated from the model is not the same as invalidated by it. The route
// file lists endpoints, which do not move when a column arrives. The React
// Query hooks are generic CRUD over one endpoint. The admin page wrappers
// render a resource by name. Those three are generated out of the model and
// there is nothing in any of them a new column would change, so reporting
// them is noise beside the handler and the form, which really are wrong.
func followsTheFields(rel string) bool {
	rel = filepath.ToSlash(rel)
	switch {
	case strings.HasSuffix(rel, "_routes.go"):
		return false
	case strings.Contains(rel, "/hooks/use-"), strings.HasPrefix(rel, "hooks/use-"):
		return false
	case strings.HasSuffix(rel, "/page.tsx"), strings.HasSuffix(rel, "/page.jsx"):
		return false
	default:
		return true
	}
}

// tellingRank is how much a generated file's content follows the model's
// fields. Lower is more telling.
//
// The handler and the service carry the create request and the writable
// column list, which is where a missing column actually costs something: the
// admin sends the value and the API drops it. The resource definition carries
// the form. A page wrapper carries a name.
func tellingRank(path string) int {
	path = filepath.ToSlash(path)
	switch {
	case strings.Contains(path, "/internal/handlers/"),
		strings.Contains(path, "/internal/services/"):
		return 0
	case strings.HasSuffix(path, ".ts") && strings.Contains(path, "/resources/"):
		return 1
	case strings.Contains(path, "/packages/shared/"), strings.HasPrefix(path, "packages/shared/"):
		return 2
	case strings.Contains(path, "/internal/routes/"), strings.Contains(path, "/internal/factory/"):
		return 3
	default:
		return 4
	}
}

func sortedKeys(maps ...map[string][]string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range maps {
		for k := range m {
			if !seen[k] {
				seen[k] = true
				out = append(out, k)
			}
		}
	}
	sort.Strings(out)
	return out
}
