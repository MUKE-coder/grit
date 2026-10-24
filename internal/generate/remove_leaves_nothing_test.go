package generate

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestGenerateThenRemoveLeavesNothingBehind generates a resource with the
// features that write extra files, removes it with the real remover, and then
// reads the whole project back looking for the resource.
//
// Every removal bug this repository has had was found by a person running the
// commands, one or two at a time: the policy file and its register call
// (v3.411.1), then the seeder, the tree files and the workflow service
// (v3.411.2). The tests that existed could not find them because each one
// checked a removal somebody had already thought of, and the round trip
// reversed the injections by hand instead of calling the remover.
//
// This asks the opposite question. Not "is this file deleted" but "is the
// resource gone", which is the question `grit remove` claims to answer when it
// prints a tick. Anything that still names the resource is listed, so a new
// injection with no removal fails here rather than in somebody's project.
//
// The name is two words on purpose. BlogPost spells its kebab and snake
// variants differently, and a one-word fixture cannot tell them apart: that is
// exactly how the seeder survived.
func TestGenerateThenRemoveLeavesNothingBehind(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("BlogPost", "title:string,status:select[draft,published]")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	// --tree, which writes a service, a handler, a test and five routes.
	def.Tree = true

	// A status field with a workflow, which writes internal/workflow/<name>.go,
	// a transition service and two more routes.
	for i := range def.Fields {
		if def.Fields[i].Name == "status" {
			def.Fields[i].Workflow = &WorkflowSpec{
				Initial: "draft",
				Transitions: []TransitionSpec{
					{Action: "publish", Label: "Publish", From: []string{"draft"}, To: "published"},
				},
			}
		}
	}

	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Run(): %v", err)
	}
	// --seed, whose file is named for the plural kebab and registered in seed.go.
	if err := g.WriteSeeder(SeederOptions{}); err != nil {
		t.Fatalf("WriteSeeder: %v", err)
	}

	names := MakeNames("BlogPost")
	if !mentionsResource(t, root, names) {
		t.Fatal("generate wrote nothing that names the resource, so this test " +
			"would pass whatever remove did")
	}

	t.Chdir(root)
	if err := RemoveResource("BlogPost"); err != nil {
		t.Fatalf("RemoveResource: %v", err)
	}

	for _, line := range survivors(t, root, names) {
		t.Errorf("remove left this behind: %s", line)
	}
}

// resourceVariants is every spelling of the resource a generated file might
// use. A removal that handles one spelling and not another is the bug this
// list exists to catch.
func resourceVariants(names Names) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range []string{
		names.Pascal, names.Camel, names.Snake, names.Kebab, names.Lower,
		names.Plural, names.PluralPascal, names.PluralSnake, names.PluralKebab,
		strings.ToUpper(names.PluralSnake),
	} {
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// mentionsResource reports whether anything in the project names the resource.
func mentionsResource(t *testing.T, root string, names Names) bool {
	t.Helper()
	return len(survivors(t, root, names)) > 0
}

// survivors lists every file whose name mentions the resource and every line
// that does, as "path:line: text".
func survivors(t *testing.T, root string, names Names) []string {
	t.Helper()
	variants := resourceVariants(names)

	var found []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", ".next", "dist":
				return filepath.SkipDir
			}
			return nil
		}

		if len(matches(d.Name(), variants)) > 0 {
			found = append(found, rel+" (file still exists)")
			return nil
		}

		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil // a file we cannot read cannot be checked, and is not a finding
		}
		for i, line := range strings.Split(string(body), "\n") {
			if len(matches(line, variants)) == 0 {
				continue
			}
			found = append(found, rel+":"+itoa(i+1)+": "+strings.TrimSpace(line))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the project: %v", err)
	}
	sort.Strings(found)
	return found
}

// matches returns the variants that appear in s, as whole words, so that
// "post" does not match "posted" and a resource called Blog does not match
// every mention of a blog in a comment.
func matches(s string, variants []string) []string {
	var hit []string
	for _, v := range variants {
		for i := 0; i+len(v) <= len(s); i++ {
			if s[i:i+len(v)] != v {
				continue
			}
			if (i == 0 || !isWordByte(s[i-1])) &&
				(i+len(v) == len(s) || !isWordByte(s[i+len(v)])) {
				hit = append(hit, v)
				break
			}
		}
	}
	return hit
}

func isWordByte(b byte) bool {
	return b == '_' || b >= '0' && b <= '9' ||
		b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
