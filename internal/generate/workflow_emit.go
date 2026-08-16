package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// WorkflowField returns the field carrying a state machine, if any.
//
// One per resource on purpose. Two state machines on one record is a thing
// people ask for and almost always means the second one is a different
// resource: an invoice with a payment status and a fulfilment status is an
// invoice and a shipment.
func (d *ResourceDefinition) WorkflowField() *Field {
	for i := range d.Fields {
		if d.Fields[i].HasWorkflow() {
			return &d.Fields[i]
		}
	}
	return nil
}

// writeWorkflow emits internal/workflow/<resource>.go and the transition
// service, and makes sure the workflow package exists.
func (g *Generator) writeWorkflow(names Names) error {
	field := g.Definition.WorkflowField()
	if field == nil {
		return nil
	}

	apiRoot := g.APIRoot()

	// The package itself, for a project generated before workflows existed.
	pkgPath := filepath.Join(apiRoot, "internal", "workflow", "workflow.go")
	if !fileExists(pkgPath) {
		if err := writeFileWithDirs(pkgPath, scaffold.APIWorkflowGo()); err != nil {
			return fmt.Errorf("writing the workflow package: %w", err)
		}
		fmt.Println("  ✓ Added internal/workflow/workflow.go")
	}

	options := field.OptionValues()
	// Empty qualifier: this file is package workflow, where workflow.Definition
	// does not resolve.
	def := field.Workflow.GoLiteral(names.Plural, toSnakeCase(field.Name), options, "")

	body := `package workflow

// ` + names.Pascal + `Workflow is the state machine on ` + names.Pascal + `.` + toPascalCase(field.Name) + `.
//
// Generated from the workflow: block in the resource definition. Edit that and
// regenerate rather than editing here, so the admin, the API and this file
// cannot disagree about which moves are legal.
var ` + names.Pascal + `Workflow = ` + def + `

func init() {
	Register(` + names.Pascal + `Workflow)
}
`
	path := filepath.Join(apiRoot, "internal", "workflow", names.Snake+".go")
	if err := writeFileWithDirs(path, body); err != nil {
		return fmt.Errorf("writing the %s workflow: %w", names.Lower, err)
	}
	fmt.Printf("  ✓ internal/workflow/%s.go (%d states, %d transitions)\n",
		names.Snake, len(options), len(field.Workflow.Transitions))

	// The route is NOT mounted here. injectAll decides whether a resource is
	// already wired by looking for its handler in routes.go, so a transition
	// route added first makes it skip every CRUD route for that resource.
	// ensureWorkflowRoutes runs after injectAll instead.
	return g.writeWorkflowService(names, field)
}

// writeWorkflowService emits the transition method.
//
// The guard lives here rather than in the handler because a handler is one
// caller. A job, a CLI command, an importer and a sync push all reach the
// service, and a rule enforced at only one entrance is not enforced.
func (g *Generator) writeWorkflowService(names Names, field *Field) error {
	col := toSnakeCase(field.Name)
	goField := toPascalCase(field.Name)

	body := `package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"` + g.Module + `/internal/events"
	"` + g.Module + `/internal/models"
	"` + g.Module + `/internal/workflow"
)

// Transition` + names.Pascal + ` moves one ` + names.Lower + ` through its workflow.
//
// Returns the updated record, or an error describing why the move is not
// legal. The caller turns that into a 422: an invalid transition is the user
// asking for something the process does not allow, not a server fault.
//
// Permission is checked here as well as on the route, because the route is
// one of several ways in.
func Transition` + names.Pascal + `(db *gorm.DB, c *gin.Context, id, action string, can func(string) bool) (*models.` + names.Pascal + `, error) {
	var item models.` + names.Pascal + `
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	from := item.` + goField + `
	transition, err := workflow.` + names.Pascal + `Workflow.Check(from, action)
	if err != nil {
		return nil, err
	}
	if transition.Permission != "" && can != nil && !can(transition.Permission) {
		return nil, fmt.Errorf("you do not have permission to %s this ` + names.Lower + `", transition.Label)
	}

	// Guarded by the current state as well as the id. Two people pressing
	// Send at the same moment would otherwise both pass the check above and
	// both write; this makes the second one affect no rows.
	result := db.Model(&models.` + names.Pascal + `{}).
		Where("id = ? AND ` + col + ` = ?", id, from).
		Update("` + col + `", transition.To)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, workflow.ErrInvalidTransition{
			Resource: "` + names.Plural + `",
			From:     from,
			Action:   action,
			Allowed:  actionsFrom(workflow.` + names.Pascal + `Workflow, from),
		}
	}

	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// The transition is its own event, not an "updated". A subscriber that
	// cares about invoices being paid should not have to diff two versions of
	// a record to find out that is what happened.
	events.Emitted(c, "` + names.Plural + `", "` + names.Pascal + `", action, item.ID,
		` + identExpr(names, g.Definition) + `,
		fmt.Sprintf("%s: %s to %s", transition.Label, from, transition.To),
		map[string]interface{}{"` + col + `": from},
		item)

	return &item, nil
}

// actionsFrom lists the actions legal from a state, for an error message that
// tells the caller what they could have done instead.
func actionsFrom(d workflow.Definition, from string) []string {
	var out []string
	for _, t := range d.Next(from) {
		out = append(out, t.Action)
	}
	return out
}
`
	path := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+"_workflow.go")
	if err := writeFileWithDirs(path, body); err != nil {
		return fmt.Errorf("writing the %s transition service: %w", names.Lower, err)
	}
	fmt.Printf("  ✓ internal/services/%s_workflow.go\n", names.Snake)
	return nil
}

// identExpr picks the expression for a record's human label, matching what the
// CRUD handlers use.
func identExpr(names Names, def *ResourceDefinition) string {
	for _, f := range def.Fields {
		switch strings.ToLower(f.Name) {
		case "name", "title", "reference", "number", "code", "slug", "email":
			return "item." + toPascalCase(f.Name)
		}
	}
	return "item.ID"
}

// ensureWorkflowRoutes mounts POST /<resource>/:id/transitions/:action.
//
// Injected at the protected-routes marker the generator already owns, so it
// lands beside the resource's other routes.
func (g *Generator) ensureWorkflowRoutes(names Names) error {
	if g.Definition.WorkflowField() == nil {
		return nil
	}
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	// Both routes in one injection. gin's tree prefers a static segment over a
	// param at the same position, so /orders/workflow and /orders/:id coexist.
	route := fmt.Sprintf(
		"\t\tprotected.GET(\"/%s/workflow\", %sHandler.Workflow)\n"+
			"\t\tprotected.POST(\"/%s/:id/transitions/:action\", %sHandler.Transition)",
		names.Plural, names.Camel, names.Plural, names.Camel)
	if strings.Contains(content, names.Camel+"Handler.Transition") {
		return nil
	}

	if err := injectBefore(path, "// grit:routes:protected", route); err != nil {
		// Not fatal: the resource is still generated and usable, it just has
		// no transition endpoint until the route is added by hand.
		fmt.Printf("  Could not mount the transition route: %v\n", err)
		return nil
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ GET /api/%s/workflow and POST /api/%s/:id/transitions/:action\n",
		names.Plural, names.Plural)
	return nil
}
