package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The handler and the service are built from one description of the resource,
// and every shape of resource the generator knows has to come out as Go that
// builds, with the handler free of queries. The shapes are the ones whose
// fragments land in different places: relations, links, money and encrypted
// and date fields, files, a self-reference, an owner, audited reads.
func TestEveryShapeRendersAThinHandlerAndAService(t *testing.T) {
	const module = "shop/apps/api"
	for _, tc := range []struct {
		name   string
		fields string
		setup  func(*ResourceDefinition)
	}{
		{"Plain", "title:string,price:float,live:bool,count:int", nil},
		{"Linked", "name:string,category:belongs_to:Category,tags:many_to_many:Tag", nil},
		{"Priced", "name:string,price:money,secret:text:encrypted,opened_on:date,due_at:datetime", nil},
		{"Filed", "name:string,cover:file,docs:files", nil},
		{"Covered", "name:string,cover:file", nil},
		{"Node", "name:string,parent:belongs_to:Node", nil},
		{"Owned", "number:string", func(d *ResourceDefinition) {
			d.OwnedBy = "user"
			d.Fields = append(d.Fields, Field{Name: "user", Type: "belongs_to", RelatedModel: "User"})
		}},
		{"Audited", "name:string", func(d *ResourceDefinition) { d.AuditReads = true }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := setupMinimalProject(t, module)
			def, err := ParseInlineFields(tc.name, tc.fields)
			if err != nil {
				t.Fatalf("ParseInlineFields: %v", err)
			}
			if tc.setup != nil {
				tc.setup(def)
			}
			g := newTestGenerator(root, module, def)
			names := g.Names()
			if err := g.writeGoService(names); err != nil {
				t.Fatalf("service: %v", err)
			}
			if err := g.writeGoHandler(names); err != nil {
				t.Fatalf("handler: %v", err)
			}
			s := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", names.Snake+".go"))
			h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", names.Snake+".go"))
			importsMatchUse(t, "service", s)
			importsMatchUse(t, "handler", h)

			for _, q := range []string{".Where(", ".First(", ".Find(", ".Transaction(", ".Association(", "FindInBatches("} {
				if strings.Contains(h, q) {
					t.Errorf("the handler still runs a query (%s)", q)
				}
			}
			// The database package imports services for its seeders, so a
			// service that imports database back does not build. Parsing does
			// not catch it; the first live project did.
			if strings.Contains(s, module+`/internal/database"`) {
				t.Error("the service imports the database package, which imports services: an import cycle")
			}
			for _, placeholder := range []string{"{{", "}}"} {
				if strings.Contains(s, placeholder) || strings.Contains(h, placeholder) {
					t.Errorf("an unreplaced placeholder survived (%s)", placeholder)
				}
			}
		})
	}
}

// A resource with a workflow keeps its transition endpoint, and hands it the
// database bound to the request rather than a bare handle.
func TestTheWorkflowHandlerStaysThin(t *testing.T) {
	g, root := purchaseRequestGenerator(t)
	names := g.Names()
	if err := g.writeGoService(names); err != nil {
		t.Fatalf("service: %v", err)
	}
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "purchase_request.go"))
	importsMatchUse(t, "handler", h)
	transition := method(t, h, "func (h *PurchaseRequestHandler) Transition(")
	if !strings.Contains(transition, "services.TransitionPurchaseRequest(h.DB.WithContext(c.Request.Context()), c, id, action, can)") {
		t.Error("the transition does not run on the request's context")
	}
}

var majorVersion = regexp.MustCompile(`^v[0-9]+$`)

// importsMatchUse fails when a file names a package it does not import, or
// imports one it never names: the two ways a generated file parses and still
// does not build. A file field's request struct named files.FileRef in a
// handler that no longer imported files, and only a live project said so.
func importsMatchUse(t *testing.T, name, src string) {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), name, src, 0)
	if err != nil {
		t.Fatalf("%s does not parse: %v", name, err)
	}
	imported := map[string]bool{}
	for _, imp := range f.Imports {
		parts := strings.Split(strings.Trim(imp.Path.Value, `"`), "/")
		pkg := parts[len(parts)-1]
		if majorVersion.MatchString(pkg) && len(parts) > 1 {
			pkg = parts[len(parts)-2]
		}
		if imp.Name != nil {
			pkg = imp.Name.Name
		}
		if pkg != "_" {
			imported[pkg] = true
		}
	}
	// The parser resolves what the file declares; a selector on a name it
	// could not resolve is a package.
	used := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok && id.Obj == nil {
				used[id.Name] = true
			}
		}
		return true
	})
	var missing, unused []string
	for pkg := range used {
		if !imported[pkg] {
			missing = append(missing, pkg)
		}
	}
	for pkg := range imported {
		if !used[pkg] {
			unused = append(unused, pkg)
		}
	}
	sort.Strings(missing)
	sort.Strings(unused)
	if len(missing) > 0 {
		t.Errorf("the %s names %v without importing them", name, missing)
	}
	if len(unused) > 0 {
		t.Errorf("the %s imports %v and never uses them", name, unused)
	}
}
