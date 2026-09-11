package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// The tree endpoints query through the tree service, and the service through
// the request's context.
//
// The service queried s.DB directly, so no query carried the organization the
// multitenant plugin resolves onto the request, and the plugin refuses a
// tenant-owned query without one. And the handler read the node's parent itself
// before a reorder, outside the move's transaction.
func TestTheTreeRunsOnTheRequestContext(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Folder", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.Tree = true
	g := newTestGenerator(root, module, def)
	if err := g.writeTreeService(g.Names()); err != nil {
		t.Fatalf("tree: %v", err)
	}
	dir := filepath.Join(root, "apps", "api", "internal")
	svc := readTestFile(t, filepath.Join(dir, "services", "folder_tree.go"))
	h := readTestFile(t, filepath.Join(dir, "handlers", "folder_tree.go"))
	tests := readTestFile(t, filepath.Join(dir, "services", "folder_tree_test.go"))
	importsMatchUse(t, "tree service", svc)
	importsMatchUse(t, "tree handler", h)
	importsMatchUse(t, "tree tests", tests)

	for _, line := range strings.Split(svc, "\n") {
		if strings.Contains(line, "s.DB") && !strings.Contains(line, "return s.DB.WithContext(ctx)") {
			t.Errorf("a tree query bypasses the request context:\n  %s", strings.TrimSpace(line))
		}
	}
	if strings.Contains(h, "h.DB.") {
		t.Error("the tree handler still runs a query of its own")
	}
	for _, call := range []string{
		"h.Tree.Tree(c.Request.Context())",
		"h.Tree.Breadcrumbs(c.Request.Context(), ",
		"h.Tree.Move(c.Request.Context(), c.Param(\"id\"), req.ParentID, req.Position)",
		"h.Tree.Reorder(c.Request.Context(), ",
		"h.Tree.RebuildPaths(c.Request.Context())",
	} {
		if !strings.Contains(h, call) {
			t.Errorf("the tree handler does not call %s", call)
		}
	}
	// The parent a reorder keeps is read inside the move's transaction.
	move := method(t, svc, "func (s *FolderTreeService) Move(")
	if !strings.Contains(move, "} else if node.ParentID != nil {") {
		t.Error("Move does not keep the node's parent when none is given")
	}
	if !strings.Contains(tests, "func TestFolderTreeMoveWithoutAParentKeepsIt(") {
		t.Error("the project's tree tests do not cover a reorder that keeps the parent")
	}
}
