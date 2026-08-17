package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Trees, and why they are shaped like this.
//
// A category tree is Electronics above Cameras above Lenses, and the query a
// storefront actually asks is "every product in Electronics or anything under
// it". There are three ways to answer that:
//
//   - Recursive CTE. Reads well, and is the wrong choice here: Grit supports
//     Postgres, MySQL and SQLite, and CTE support, syntax and performance
//     differ across all three. A generator that emits one query per dialect is
//     a generator with three bugs.
//   - Nested sets. One indexed range query, and every insert rewrites half the
//     table. Fine for a read-only taxonomy, miserable for a category tree
//     somebody reorders in the admin.
//   - Materialized path. One indexed LIKE, identical on all three dialects, and
//     a move rewrites only the subtree that moved. This is what Grit emits.
//
// The path is "/id/id/id/", ids delimited on both sides so a prefix match
// cannot half-match an id, and so "is X above me" is one strings.Contains.
// depth is stored rather than derived, because counting separators in SQL is
// three different expressions across three dialects for something a move can
// keep correct with a single delta.

// treeFields returns the struct fields --tree adds to the model.
//
// position is deliberately part of this rather than something a project adds
// later: a category tree with no sibling order renders in insertion order,
// which means the admin's drag handles have nothing to write to.
func treeFields(names Names) string {
	return fmt.Sprintf(`	// --- tree (grit --tree) ---------------------------------------------------
	// Children is a convenience for one level. The whole tree comes from the
	// service in a single query rather than a preload per level.
	Children []%s `+"`"+`gorm:"foreignKey:%s" json:"children,omitempty"`+"`"+`
	// Path is "/id/id/id/", this row's id last, ids delimited on both sides so
	// a prefix match cannot half-match an id. Descendants are one indexed LIKE.
	Path string `+"`"+`gorm:"size:1024;index" json:"path"`+"`"+`
	// Depth is 0 for a root. Stored rather than counted from Path, because a
	// move can keep it correct with one delta and counting separators in SQL is
	// three expressions across three dialects.
	Depth int `+"`"+`gorm:"index" json:"depth"`+"`"+`
	// Position orders siblings. Without it a tree renders in insertion order
	// and the admin's drag handles have nowhere to write.
	Position int `+"`"+`gorm:"index;default:0" json:"position"`+"`"+`
`, names.Pascal, treeParentFK(names))
}

// treeParentFK is the FK column holding the parent id. Fixed rather than
// derived from the field name, because the tree service, the move endpoint and
// the admin all have to agree on it.
func treeParentFK(names Names) string { return "ParentID" }

// treeCreateHook is appended to BeforeCreate, after the id is assigned.
//
// Path cannot be computed before the id exists, and the id is assigned in the
// same hook a few lines above this. That ordering is the only reason this is
// not its own hook.
func treeCreateHook() string {
	return `	// Materialized path, computed here because it needs the id assigned above.
	if err := m.resolveTreePath(tx); err != nil {
		return err
	}
`
}

// treeModelMethods returns the path resolution shared by create and update.
func treeModelMethods(names Names) string {
	return fmt.Sprintf(`
// resolveTreePath sets Path and Depth from the parent, and refuses a cycle.
//
// The cycle check is one string comparison: the parent's path contains every id
// above it, so if it contains this row's id then this row is being moved under
// its own descendant. Without the check, that move detaches a whole subtree
// from the tree and no query ever finds it again.
func (m *%[1]s) resolveTreePath(tx *gorm.DB) error {
	if m.ParentID == "" {
		m.Path = "/" + m.ID + "/"
		m.Depth = 0
		return nil
	}
	if m.ParentID == m.ID {
		return fmt.Errorf("a %[2]s cannot be its own parent")
	}

	var parent %[1]s
	// NewDB so this lookup does not inherit the conditions of the statement
	// being built, which would silently scope it to the row being written.
	err := tx.Session(&gorm.Session{NewDB: true}).
		Select("id", "path", "depth").
		Where("id = ?", m.ParentID).
		First(&parent).Error
	if err != nil {
		return fmt.Errorf("parent %%s does not exist: %%w", m.ParentID, err)
	}
	if m.ID != "" && strings.Contains(parent.Path, "/"+m.ID+"/") {
		return fmt.Errorf("cannot move a %[2]s under its own descendant")
	}

	m.Path = parent.Path + m.ID + "/"
	m.Depth = parent.Depth + 1
	return nil
}

// IsRoot reports whether this node sits at the top of the tree.
func (m *%[1]s) IsRoot() bool { return m.ParentID == "" }

// AncestorIDs returns the ids above this node, outermost first, read straight
// from the path. No query: that is the point of storing it.
func (m *%[1]s) AncestorIDs() []string {
	parts := strings.Split(strings.Trim(m.Path, "/"), "/")
	if len(parts) <= 1 {
		return nil
	}
	return parts[:len(parts)-1] // everything except this node
}
`, names.Pascal, names.Lower)
}

// treeServiceSource generates internal/services/<resource>_tree.go.
//
// These live in a service rather than the handler because the admin, the public
// surface and the importer all need the same answers, and because Move is the
// one operation with enough invariants to deserve tests of its own.
func treeServiceSource(module string, names Names) string {
	return fmt.Sprintf(`package services

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"%[2]s/internal/models"
)

// %[1]sTreeService answers the questions a hierarchy is for.
//
// Every query here is one round trip. A tree rendered by walking parents is the
// N+1 that makes people give up on hierarchies.
type %[1]sTreeService struct {
	DB *gorm.DB
}

func New%[1]sTreeService(db *gorm.DB) *%[1]sTreeService {
	return &%[1]sTreeService{DB: db}
}

// %[1]sNode is a %[3]s with its children attached, for rendering a tree.
type %[1]sNode struct {
	models.%[1]s
	Children []*%[1]sNode `+"`"+`json:"children"`+"`"+`
}

// Tree returns the whole hierarchy in ONE query, assembled in Go.
//
// Ordered by depth so a parent is always seen before its children, then by
// position and name so siblings come back in the order the admin arranged them.
func (s *%[1]sTreeService) Tree() ([]*%[1]sNode, error) {
	var rows []models.%[1]s
	err := s.DB.Where("archived_at IS NULL").
		Order("depth asc, position asc, name asc").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("loading %[4]s: %%w", err)
	}

	byID := make(map[string]*%[1]sNode, len(rows))
	roots := make([]*%[1]sNode, 0)
	for i := range rows {
		node := &%[1]sNode{%[1]s: rows[i]}
		byID[node.ID] = node
	}
	// Depth order means a parent is already in the map when its child arrives.
	// A child whose parent is missing (archived, or deleted mid-read) is
	// promoted to a root rather than dropped: losing a subtree from a tree view
	// is worse than showing it slightly out of place.
	for i := range rows {
		node := byID[rows[i].ID]
		if parent, ok := byID[rows[i].ParentID]; ok && rows[i].ParentID != "" {
			parent.Children = append(parent.Children, node)
			continue
		}
		roots = append(roots, node)
	}
	return roots, nil
}

// Roots returns the top level only.
//
// NULL counts as no parent, not just the empty string. Adding --tree to a
// resource that already has rows is the normal case, and AutoMigrate fills the
// new column with NULL, so a query testing only for '' finds nothing at all.
func (s *%[1]sTreeService) Roots() ([]models.%[1]s, error) {
	var rows []models.%[1]s
	err := s.DB.Where("(parent_id = '' OR parent_id IS NULL) AND archived_at IS NULL").
		Order("position asc, name asc").
		Find(&rows).Error
	return rows, err
}

// Children returns one level below a node.
func (s *%[1]sTreeService) Children(id string) ([]models.%[1]s, error) {
	var rows []models.%[1]s
	err := s.DB.Where("parent_id = ? AND archived_at IS NULL", id).
		Order("position asc, name asc").
		Find(&rows).Error
	return rows, err
}

// Descendants returns everything below a node, at any depth, in one query.
//
// The LIKE is anchored at the start with a delimited prefix, so /1/ matches
// /1/2/ and /1/2/3/ and never /11/.
func (s *%[1]sTreeService) Descendants(id string) ([]models.%[1]s, error) {
	var node models.%[1]s
	if err := s.DB.Select("id", "path").Where("id = ?", id).First(&node).Error; err != nil {
		return nil, fmt.Errorf("%[3]s %%s not found: %%w", id, err)
	}
	var rows []models.%[1]s
	err := s.DB.Where("path LIKE ? AND id <> ? AND archived_at IS NULL", node.Path+"%%", id).
		Order("depth asc, position asc, name asc").
		Find(&rows).Error
	return rows, err
}

// DescendantIDs returns the node's id plus every id below it, which is what a
// filter wants: "products in Electronics" means Electronics and everything
// under it.
func (s *%[1]sTreeService) DescendantIDs(id string) ([]string, error) {
	var node models.%[1]s
	if err := s.DB.Select("id", "path").Where("id = ?", id).First(&node).Error; err != nil {
		return nil, fmt.Errorf("%[3]s %%s not found: %%w", id, err)
	}
	var ids []string
	err := s.DB.Model(&models.%[1]s{}).
		Where("path LIKE ? AND archived_at IS NULL", node.Path+"%%").
		Pluck("id", &ids).Error
	return ids, err
}

// Breadcrumbs returns the path from the root down to and including this node.
//
// The ids come from the stored path, so this is one IN query however deep the
// tree is, and the result is re-sorted into path order because IN does not
// promise one.
func (s *%[1]sTreeService) Breadcrumbs(id string) ([]models.%[1]s, error) {
	var node models.%[1]s
	if err := s.DB.Where("id = ?", id).First(&node).Error; err != nil {
		return nil, fmt.Errorf("%[3]s %%s not found: %%w", id, err)
	}
	ids := strings.Split(strings.Trim(node.Path, "/"), "/")
	if len(ids) == 0 {
		return []models.%[1]s{node}, nil
	}

	var rows []models.%[1]s
	if err := s.DB.Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	byID := make(map[string]models.%[1]s, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
	}
	ordered := make([]models.%[1]s, 0, len(ids))
	for _, wanted := range ids {
		if r, ok := byID[wanted]; ok {
			ordered = append(ordered, r)
		}
	}
	return ordered, nil
}

// Move reparents a node and carries its subtree with it.
//
// Three things have to happen together, which is why this is a transaction and
// not three calls:
//
//  1. The node gets its new parent, path and depth.
//  2. Every descendant's path is rewritten, because their paths embed the old
//     prefix. REPLACE is used rather than string surgery in Go: it is one
//     UPDATE, it exists on Postgres, MySQL and SQLite alike, and it cannot race
//     with a concurrent read the way read-modify-write can.
//  3. Every descendant's depth shifts by the same delta, because a subtree
//     keeps its shape when it moves.
//
// Refuses a move that would make the node its own ancestor. Without that check
// the subtree is detached from the tree and no query finds it again.
func (s *%[1]sTreeService) Move(id, newParentID string, position int) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if id == newParentID {
		return fmt.Errorf("a %[3]s cannot be its own parent")
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		var node models.%[1]s
		if err := tx.Where("id = ?", id).First(&node).Error; err != nil {
			return fmt.Errorf("%[3]s %%s not found: %%w", id, err)
		}
		oldPath := node.Path
		oldDepth := node.Depth

		newPath := "/" + id + "/"
		newDepth := 0
		if newParentID != "" {
			var parent models.%[1]s
			if err := tx.Select("id", "path", "depth").Where("id = ?", newParentID).First(&parent).Error; err != nil {
				return fmt.Errorf("parent %%s does not exist: %%w", newParentID, err)
			}
			if strings.Contains(parent.Path, "/"+id+"/") {
				return fmt.Errorf("cannot move a %[3]s under its own descendant")
			}
			newPath = parent.Path + id + "/"
			newDepth = parent.Depth + 1
		}

		err := tx.Model(&models.%[1]s{}).Where("id = ?", id).Updates(map[string]any{
			"parent_id": newParentID,
			"path":      newPath,
			"depth":     newDepth,
			"position":  position,
		}).Error
		if err != nil {
			return fmt.Errorf("moving %[3]s: %%w", err)
		}

		if oldPath == newPath {
			return nil // reordered among its siblings, nothing below it moved
		}
		// A node with no path has no subtree to rewrite, and asking anyway is
		// destructive rather than merely pointless: "" as a LIKE prefix matches
		// every row in the table, so the depth increment below lands on all of
		// them. That is not hypothetical. Dragging a category that predated
		// --tree added one to the depth of every category in the database, and
		// nothing looked wrong until somebody read the rows.
		//
		// A row without a path is a row that predates --tree. RebuildPaths is
		// what fixes those, and this row now has a correct path either way.
		if oldPath == "" {
			return nil
		}

		// The subtree. Excludes the node itself, which is already written.
		delta := newDepth - oldDepth
		return tx.Model(&models.%[1]s{}).
			Where("path LIKE ? AND id <> ?", oldPath+"%%", id).
			Updates(map[string]any{
				"path":  gorm.Expr("REPLACE(path, ?, ?)", oldPath, newPath),
				"depth": gorm.Expr("depth + ?", delta),
			}).Error
	})
}

// Reorder writes a new sibling order in one transaction, which is what a
// drag-and-drop tree sends: the ids of one parent's children, in order.
//
// The parent is part of the WHERE on purpose, so a stale client cannot reorder
// a node into a parent it no longer belongs to. It has to treat NULL as no
// parent, though: rows that predate --tree have a NULL parent_id, and
// "parent_id = ''" matched none of them. The reorder then returned 200 having
// updated nothing, which is the worst kind of bug, the one that looks like it
// worked.
func (s *%[1]sTreeService) Reorder(parentID string, orderedIDs []string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for position, id := range orderedIDs {
			q := tx.Model(&models.%[1]s{}).Where("id = ?", id)
			if parentID == "" {
				q = q.Where("(parent_id = '' OR parent_id IS NULL)")
			} else {
				q = q.Where("parent_id = ?", parentID)
			}
			// parent_id is normalised on the way past, so the next reorder of
			// this row has one thing to compare rather than two.
			updates := map[string]any{"position": position}
			if parentID == "" {
				updates["parent_id"] = ""
			}
			if err := q.Updates(updates).Error; err != nil {
				return fmt.Errorf("reordering %%s: %%w", id, err)
			}
		}
		return nil
	})
}

// RebuildPaths recomputes every path and depth from parent_id alone.
//
// The escape hatch for a tree whose paths were written by a bulk import that
// went around the hooks, or one that predates --tree.
//
// Deliberately not SQL. The obvious version is an UPDATE per level joining each
// row to its parent, and it is a trap: string concatenation is || on Postgres
// and SQLite but CONCAT on MySQL, UPDATE aliasing differs, and the WHERE that
// decides "not yet computed" has to be exactly right or it rewrites every row
// on every pass and leaves the table worse than it found it. That is not a
// hypothetical, it is what the first version of this function did.
//
// So: one read, the arithmetic in Go, one UPDATE per row, all in a transaction.
// This runs rarely, on a table with hundreds of rows, and being obviously
// correct on all three dialects is worth more here than being fast.
func (s *%[1]sTreeService) RebuildPaths() error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		return s.rebuildPathsInGo(tx)
	})
}

func (s *%[1]sTreeService) rebuildPathsInGo(tx *gorm.DB) error {
	var rows []models.%[1]s
	if err := tx.Order("parent_id asc").Find(&rows).Error; err != nil {
		return err
	}
	paths := make(map[string]string, len(rows))
	depths := make(map[string]int, len(rows))

	// Repeat until nothing changes: a child may be seen before its parent, and
	// the number of passes needed is at most the depth of the tree.
	for pass := 0; pass < 64; pass++ {
		progress := false
		for _, r := range rows {
			if _, done := paths[r.ID]; done {
				continue
			}
			if r.ParentID == "" {
				paths[r.ID] = "/" + r.ID + "/"
				depths[r.ID] = 0
				progress = true
				continue
			}
			if parentPath, ok := paths[r.ParentID]; ok {
				paths[r.ID] = parentPath + r.ID + "/"
				depths[r.ID] = depths[r.ParentID] + 1
				progress = true
			}
		}
		if !progress {
			break
		}
	}

	for id, path := range paths {
		fields := map[string]any{"path": path, "depth": depths[id]}
		// Normalise NULL to '' while we are here, so nothing downstream has to
		// keep asking about both. A row migrated into a tree arrives with NULL.
		if depths[id] == 0 {
			fields["parent_id"] = ""
		}
		if err := tx.Model(&models.%[1]s{}).Where("id = ?", id).Updates(fields).Error; err != nil {
			return err
		}
	}
	return nil
}
`, names.Pascal, module, names.Lower, names.PluralSnake)
}

// writeTreeService writes internal/services/<resource>_tree.go, and mounts the
// endpoints the admin's tree needs.
//
// Its own file rather than an addition to the resource service, so regenerating
// the resource cannot clobber it and a project can hand-edit one without
// touching the other.
func (g *Generator) writeTreeService(names Names) error {
	if !g.Definition.Tree {
		return nil
	}

	path := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+"_tree.go")
	if err := writeFileWithDirs(path, treeServiceSource(g.Module, names)); err != nil {
		return fmt.Errorf("writing tree service: %w", err)
	}
	fmt.Printf("  ✓ internal/services/%s_tree.go (tree, descendants, breadcrumbs, move)\n", names.Snake)

	// The behaviour tests ship with the service. They need a database to be
	// worth anything, and running them in the project means running them
	// against the dialect the project actually uses.
	testPath := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+"_tree_test.go")
	if err := writeFileWithDirs(testPath, treeTestSource(g.Module, names)); err != nil {
		return fmt.Errorf("writing tree tests: %w", err)
	}
	fmt.Printf("  ✓ internal/services/%s_tree_test.go (8 tests: paths, move, cycles, rebuild)\n", names.Snake)

	// Routes are NOT mounted here. ensureTreeRoutes runs after injectAll,
	// because injectAll decides a resource is already wired by finding its
	// handler in routes.go, and a tree route mounted first makes it skip every
	// CRUD route for the resource.
	return g.writeTreeHandler(names)
}

// writeTreeHandler writes the endpoints over the tree service.
func (g *Generator) writeTreeHandler(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_tree.go")
	if err := writeFileWithDirs(path, treeHandlerSource(g.Module, names)); err != nil {
		return fmt.Errorf("writing tree handler: %w", err)
	}
	fmt.Printf("  ✓ internal/handlers/%s_tree.go\n", names.Snake)
	return nil
}

// treeHandlerSource generates the tree endpoints.
func treeHandlerSource(module string, names Names) string {
	return fmt.Sprintf(`package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"%[2]s/internal/services"
)

// %[1]sTreeHandler serves the hierarchy: the whole tree for a picker, the
// breadcrumbs for a detail page, and the two writes a drag-and-drop tree makes.
type %[1]sTreeHandler struct {
	DB   *gorm.DB
	Tree *services.%[1]sTreeService
}

func New%[1]sTreeHandler(db *gorm.DB) *%[1]sTreeHandler {
	return &%[1]sTreeHandler{DB: db, Tree: services.New%[1]sTreeService(db)}
}

// GetTree handles GET /api/v1/%[3]s/tree.
func (h *%[1]sTreeHandler) GetTree(c *gin.Context) {
	nodes, err := h.Tree.Tree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to load the tree",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nodes})
}

// GetBreadcrumbs handles GET /api/v1/%[3]s/:id/breadcrumbs.
func (h *%[1]sTreeHandler) GetBreadcrumbs(c *gin.Context) {
	rows, err := h.Tree.Breadcrumbs(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": err.Error(),
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// Move handles PATCH /api/v1/%[3]s/:id/move.
//
// parent_id is a pointer so the JSON can distinguish "move to the root"
// (parent_id: "") from "leave the parent alone" (field absent). Without that
// distinction a reorder within one parent would silently promote the node.
func (h *%[1]sTreeHandler) Move(c *gin.Context) {
	var req struct {
		ParentID *string `+"`"+`json:"parent_id"`+"`"+`
		Position int     `+"`"+`json:"position"`+"`"+`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}

	id := c.Param("id")
	parentID := ""
	if req.ParentID != nil {
		parentID = *req.ParentID
	} else {
		// Field absent: keep the parent it has, and treat this as a reorder.
		var current struct{ ParentID string }
		if err := h.DB.Table("%[4]s").Select("parent_id").Where("id = ?", id).Scan(&current).Error; err == nil {
			parentID = current.ParentID
		}
	}

	if err := h.Tree.Move(id, parentID, req.Position); err != nil {
		// A refused move is the caller's mistake, not a server fault: it is
		// almost always an attempt to drop a node inside its own subtree.
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
			"code": "INVALID_MOVE", "message": err.Error(),
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Moved"})
}

// RebuildPaths handles POST /api/v1/%[3]s/rebuild-tree.
//
// Exists for one specific moment: --tree was added to a resource that already
// had rows, so every one of them has a NULL path and the tree renders flat.
// Also the repair for a bulk import that went around the hooks.
func (h *%[1]sTreeHandler) RebuildPaths(c *gin.Context) {
	if err := h.Tree.RebuildPaths(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": err.Error(),
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Paths rebuilt"})
}

// Reorder handles POST /api/v1/%[3]s/reorder, the ids of one parent's children
// in their new order.
func (h *%[1]sTreeHandler) Reorder(c *gin.Context) {
	var req struct {
		ParentID string   `+"`"+`json:"parent_id"`+"`"+`
		IDs      []string `+"`"+`json:"ids" binding:"required"`+"`"+`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}
	if err := h.Tree.Reorder(req.ParentID, req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to reorder",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reordered"})
}
`, names.Pascal, module, names.Plural, names.PluralSnake)
}

// ensureTreeRoutes mounts the tree endpoints.
//
// Runs after injectAll for the same reason the public and workflow routes do:
// injectAll decides a resource is already wired by finding its handler in
// routes.go, and a tree route added first makes it skip every CRUD route.
func (g *Generator) ensureTreeRoutes(names Names) {
	if !g.Definition.Tree {
		return
	}
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(data)

	// Checked route by route rather than "does the handler exist".
	//
	// The all-or-nothing version is wrong in a way that keeps recurring: a
	// project generated before rebuild-tree existed has the handler and four of
	// the five routes, and one glance at the handler name declared it done. The
	// endpoint then exists in Go and is reachable from nowhere.
	handlerDecl := fmt.Sprintf("\t\t%sTreeHandler := handlers.New%sTreeHandler(db)", names.Camel, names.Pascal)
	routes := []struct {
		method string
		path   string
		call   string
	}{
		{"GET", "/" + names.Plural + "/tree", "GetTree"},
		{"GET", "/" + names.Plural + "/:id/breadcrumbs", "GetBreadcrumbs"},
		{"PATCH", "/" + names.Plural + "/:id/move", "Move"},
		{"POST", "/" + names.Plural + "/reorder", "Reorder"},
		{"POST", "/" + names.Plural + "/rebuild-tree", "RebuildPaths"},
	}

	lines := []string{}
	// The handler is constructed inline in the route block rather than beside the
	// other handlers, so a first run is one injection at one marker instead of
	// two that can half-apply.
	if !strings.Contains(content, handlerDecl) {
		lines = append(lines, handlerDecl)
	}
	added := []string{}
	for _, r := range routes {
		call := fmt.Sprintf("%sTreeHandler.%s", names.Camel, r.call)
		if strings.Contains(content, call) {
			continue
		}
		lines = append(lines, fmt.Sprintf("\t\tprotected.%s(%q, %s)", r.method, r.path, call))
		added = append(added, r.method+" "+r.path)
	}
	if len(added) == 0 {
		return
	}

	if err := injectBefore(path, "// grit:routes:custom", strings.Join(lines, "\n")); err != nil {
		fmt.Printf("  Could not mount the tree routes: %v\n", err)
		return
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ %s\n", strings.Join(added, ", "))
}

// treeTestSource generates internal/services/<resource>_tree_test.go.
//
// These ship into the project rather than living in the generator for two
// reasons. They need a database to be worth anything: the questions are "does
// REPLACE rewrite the subtree" and "does this LIKE prefix match what I think",
// and no amount of string comparison on generated source answers either. And
// running in the project means running against the dialect the project actually
// uses, which is the only place a dialect bug shows up.
//
// This suite earned its keep before it shipped: it caught RebuildPaths
// rewriting every row on every pass and leaving the table worse than it found
// it.
func treeTestSource(module string, names Names) string {
	return fmt.Sprintf(`package services_test

import (
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"%[2]s/internal/models"
	"%[2]s/internal/services"
)

// Generated by grit --tree. The hierarchy has invariants worth holding onto:
// a move carries its subtree, and a node can never end up inside itself.

func %[1]sTreeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %%v", err)
	}
	if err := db.AutoMigrate(&models.%[1]s{}); err != nil {
		t.Fatalf("migrate: %%v", err)
	}
	return db
}

func make%[1]s(t *testing.T, db *gorm.DB, name, parentID string) models.%[1]s {
	t.Helper()
	row := models.%[1]s{Name: name, ParentID: parentID}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create %%s: %%v", name, err)
	}
	return row
}

func Test%[1]sTreePathsAndDepths(t *testing.T) {
	db := %[1]sTreeTestDB(t)

	top := make%[1]s(t, db, "Top", "")
	middle := make%[1]s(t, db, "Middle", top.ID)
	bottom := make%[1]s(t, db, "Bottom", middle.ID)

	if top.Path != "/"+top.ID+"/" || top.Depth != 0 {
		t.Errorf("root: path=%%q depth=%%d", top.Path, top.Depth)
	}
	if want := "/" + top.ID + "/" + middle.ID + "/"; middle.Path != want || middle.Depth != 1 {
		t.Errorf("child: path=%%q depth=%%d, want %%q depth 1", middle.Path, middle.Depth, want)
	}
	if want := "/" + top.ID + "/" + middle.ID + "/" + bottom.ID + "/"; bottom.Path != want || bottom.Depth != 2 {
		t.Errorf("grandchild: path=%%q depth=%%d, want %%q depth 2", bottom.Path, bottom.Depth, want)
	}

	// Read from the stored path, no query.
	if got := bottom.AncestorIDs(); len(got) != 2 || got[0] != top.ID || got[1] != middle.ID {
		t.Errorf("AncestorIDs = %%v", got)
	}
	if !top.IsRoot() || bottom.IsRoot() {
		t.Error("IsRoot is wrong")
	}
}

func Test%[1]sTreeRefusesAMissingParent(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	row := models.%[1]s{Name: "Orphan", ParentID: "does-not-exist"}
	if err := db.Create(&row).Error; err == nil {
		t.Error("a child of a parent that does not exist must be refused")
	}
}

func Test%[1]sTreeQueries(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	top := make%[1]s(t, db, "Top", "")
	middle := make%[1]s(t, db, "Middle", top.ID)
	bottom := make%[1]s(t, db, "Bottom", middle.ID)
	make%[1]s(t, db, "Other", "")

	roots, err := svc.Roots()
	if err != nil {
		t.Fatalf("Roots: %%v", err)
	}
	if len(roots) != 2 {
		t.Errorf("Roots = %%d, want 2", len(roots))
	}

	desc, err := svc.Descendants(top.ID)
	if err != nil {
		t.Fatalf("Descendants: %%v", err)
	}
	if len(desc) != 2 {
		t.Errorf("Descendants = %%d, want 2", len(desc))
	}

	// What a filter wants: this node AND everything under it.
	ids, err := svc.DescendantIDs(top.ID)
	if err != nil {
		t.Fatalf("DescendantIDs: %%v", err)
	}
	if len(ids) != 3 {
		t.Errorf("DescendantIDs = %%d, want 3 (self + 2)", len(ids))
	}

	crumbs, err := svc.Breadcrumbs(bottom.ID)
	if err != nil {
		t.Fatalf("Breadcrumbs: %%v", err)
	}
	if len(crumbs) != 3 || crumbs[0].Name != "Top" || crumbs[2].Name != "Bottom" {
		t.Errorf("Breadcrumbs came back wrong: %%+v", crumbs)
	}

	tree, err := svc.Tree()
	if err != nil {
		t.Fatalf("Tree: %%v", err)
	}
	if len(tree) != 2 {
		t.Fatalf("Tree roots = %%d, want 2", len(tree))
	}
	for _, node := range tree {
		if node.Name != "Top" {
			continue
		}
		if len(node.Children) != 1 || len(node.Children[0].Children) != 1 {
			t.Error("the tree is not nested three deep")
		}
	}
}

// The point of a materialized path: a move rewrites the subtree with it.
func Test%[1]sTreeMoveCarriesTheSubtree(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	top := make%[1]s(t, db, "Top", "")
	middle := make%[1]s(t, db, "Middle", top.ID)
	bottom := make%[1]s(t, db, "Bottom", middle.ID)
	other := make%[1]s(t, db, "Other", "")

	if err := svc.Move(middle.ID, other.ID, 0); err != nil {
		t.Fatalf("Move: %%v", err)
	}

	var movedMiddle, movedBottom models.%[1]s
	db.Where("id = ?", middle.ID).First(&movedMiddle)
	db.Where("id = ?", bottom.ID).First(&movedBottom)

	wantMiddle := "/" + other.ID + "/" + middle.ID + "/"
	if movedMiddle.Path != wantMiddle || movedMiddle.Depth != 1 || movedMiddle.ParentID != other.ID {
		t.Errorf("moved node: path=%%q depth=%%d parent=%%q", movedMiddle.Path, movedMiddle.Depth, movedMiddle.ParentID)
	}
	if want := wantMiddle + bottom.ID + "/"; movedBottom.Path != want || movedBottom.Depth != 2 {
		t.Errorf("descendant: path=%%q depth=%%d, want %%q depth 2", movedBottom.Path, movedBottom.Depth, want)
	}

	if ids, _ := svc.DescendantIDs(other.ID); len(ids) != 3 {
		t.Errorf("the new parent should hold 3 ids, holds %%d", len(ids))
	}
	if ids, _ := svc.DescendantIDs(top.ID); len(ids) != 1 {
		t.Errorf("the old parent should be alone, holds %%d", len(ids))
	}
}

func Test%[1]sTreeMoveToRoot(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	top := make%[1]s(t, db, "Top", "")
	child := make%[1]s(t, db, "Child", top.ID)

	if err := svc.Move(child.ID, "", 0); err != nil {
		t.Fatalf("Move to root: %%v", err)
	}
	var moved models.%[1]s
	db.Where("id = ?", child.ID).First(&moved)
	if moved.Path != "/"+child.ID+"/" || moved.Depth != 0 || moved.ParentID != "" {
		t.Errorf("promoted node: path=%%q depth=%%d parent=%%q", moved.Path, moved.Depth, moved.ParentID)
	}
}

// The invariant that matters most: without it the subtree is detached from the
// tree and no query ever finds it again.
func Test%[1]sTreeRefusesACycle(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	top := make%[1]s(t, db, "Top", "")
	middle := make%[1]s(t, db, "Middle", top.ID)
	bottom := make%[1]s(t, db, "Bottom", middle.ID)

	if err := svc.Move(top.ID, bottom.ID, 0); err == nil {
		t.Error("moving a node under its own descendant must be refused")
	}
	if err := svc.Move(top.ID, top.ID, 0); err == nil {
		t.Error("a node cannot be its own parent")
	}

	var check models.%[1]s
	db.Where("id = ?", top.ID).First(&check)
	if check.ParentID != "" || check.Depth != 0 {
		t.Errorf("a refused move must leave the row untouched: parent=%%q depth=%%d", check.ParentID, check.Depth)
	}
}

// Rows that predate --tree arrive with a NULL parent_id, because AutoMigrate
// fills a new column with NULL. Every query that means "is a root" has to say so
// in a way that matches those rows: "parent_id = ''" alone silently matches none
// of them, and a reorder that updates nothing still answers 200.
func Test%[1]sTreeToleratesNullParents(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	a := make%[1]s(t, db, "A", "")
	b := make%[1]s(t, db, "B", "")
	// Exactly what adding --tree to an existing table leaves behind.
	db.Exec("UPDATE %[3]s SET parent_id = NULL")

	roots, err := svc.Roots()
	if err != nil {
		t.Fatalf("Roots: %%v", err)
	}
	if len(roots) != 2 {
		t.Errorf("Roots with NULL parents = %%d, want 2", len(roots))
	}

	if err := svc.Reorder("", []string{b.ID, a.ID}); err != nil {
		t.Fatalf("Reorder: %%v", err)
	}
	ordered, _ := svc.Roots()
	if len(ordered) != 2 || ordered[0].Name != "B" {
		names := []string{}
		for _, r := range ordered {
			names = append(names, r.Name)
		}
		t.Errorf("order = %%v, want [B A]", names)
	}
}

// Moving a node that has no path must not touch any other row.
//
// The bug this pins down: "" as a LIKE prefix matches every row in the table, so
// the subtree rewrite that follows a move landed on all of them and added one to
// every depth in the database. A node with no path is one that predates --tree,
// which makes this the normal case the first time somebody drags a row after
// adding the flag to an existing table.
func Test%[1]sTreeMoveOfAPathlessNodeLeavesOthersAlone(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	top := make%[1]s(t, db, "Top", "")
	middle := make%[1]s(t, db, "Middle", top.ID)
	bottom := make%[1]s(t, db, "Bottom", middle.ID)

	// A row as AutoMigrate leaves it when --tree is added to an existing table.
	stray := make%[1]s(t, db, "Stray", "")
	db.Model(&models.%[1]s{}).Where("id = ?", stray.ID).
		Updates(map[string]any{"path": "", "depth": 0})

	if err := svc.Move(stray.ID, top.ID, 0); err != nil {
		t.Fatalf("Move: %%v", err)
	}

	var movedStray, checkTop, checkMiddle, checkBottom models.%[1]s
	db.Where("id = ?", stray.ID).First(&movedStray)
	db.Where("id = ?", top.ID).First(&checkTop)
	db.Where("id = ?", middle.ID).First(&checkMiddle)
	db.Where("id = ?", bottom.ID).First(&checkBottom)

	// The moved row is repaired.
	if movedStray.Depth != 1 || movedStray.Path != top.Path+stray.ID+"/" {
		t.Errorf("moved node: path=%%q depth=%%d", movedStray.Path, movedStray.Depth)
	}
	// And nothing else moved an inch.
	if checkTop.Depth != 0 || checkMiddle.Depth != 1 || checkBottom.Depth != 2 {
		t.Errorf("depths of untouched rows changed: %%d %%d %%d, want 0 1 2",
			checkTop.Depth, checkMiddle.Depth, checkBottom.Depth)
	}
	if checkBottom.Path != checkMiddle.Path+bottom.ID+"/" {
		t.Errorf("an untouched path was rewritten: %%q", checkBottom.Path)
	}
}

func Test%[1]sTreeReorderAndRebuild(t *testing.T) {
	db := %[1]sTreeTestDB(t)
	svc := services.New%[1]sTreeService(db)

	root := make%[1]s(t, db, "Root", "")
	a := make%[1]s(t, db, "A", root.ID)
	b := make%[1]s(t, db, "B", root.ID)
	c := make%[1]s(t, db, "C", root.ID)

	if err := svc.Reorder(root.ID, []string{c.ID, a.ID, b.ID}); err != nil {
		t.Fatalf("Reorder: %%v", err)
	}
	children, _ := svc.Children(root.ID)
	if got := fmt.Sprintf("%%s,%%s,%%s", children[0].Name, children[1].Name, children[2].Name); got != "C,A,B" {
		t.Errorf("order = %%s, want C,A,B", got)
	}

	// A bulk import that went around the hooks leaves paths empty, and Rebuild
	// has to reconstruct them from parent_id alone.
	db.Exec("UPDATE %[3]s SET path = '', depth = 0")
	if err := svc.RebuildPaths(); err != nil {
		t.Fatalf("RebuildPaths: %%v", err)
	}
	var rebuilt models.%[1]s
	db.Where("id = ?", a.ID).First(&rebuilt)
	if want := "/" + root.ID + "/" + a.ID + "/"; rebuilt.Path != want || rebuilt.Depth != 1 {
		t.Errorf("rebuilt: path=%%q depth=%%d, want %%q depth 1", rebuilt.Path, rebuilt.Depth, want)
	}
	_ = b
}
`, names.Pascal, module, names.PluralSnake)
}
