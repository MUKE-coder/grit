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
	} else if current, err := os.ReadFile(pkgPath); err == nil && !strings.Contains(string(current), "func RunHooks(") {
		// The transition service below runs hooks and classifies errors with
		// functions an older copy does not have.
		if g.refreshIfUnchanged(pkgPath, scaffold.APIWorkflowGo()) {
			fmt.Println("  ✓ Updated internal/workflow/workflow.go")
		} else {
			fmt.Println("  ⚠ internal/workflow/workflow.go is from an earlier release and has no transition")
			fmt.Println("    hooks. The generated service will not compile until that file is replaced")
			fmt.Println("    with the current template (grit upgrade does it when it is unedited).")
		}
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

	// An owned resource's transition route sits on the authenticated group like
	// every other, and it loaded the row by id with no check, so any account
	// could move anyone's record through its workflow.
	ownerCheck, authzImport := "", ""
	if g.Definition.IsOwned() {
		authzImport = "\t\"" + g.Module + "/internal/authz\"\n"
		ownerCheck = `
	// --owned-by: only the owner may move a row, ADMIN excepted, and somebody
	// else's row is not found rather than forbidden. c is nil when a job or a
	// command makes the move, and those are trusted.
	if c != nil && !authz.IsAdmin(c) && item.GetOwnerID() != authz.CurrentUserID(c) {
		return nil, gorm.ErrRecordNotFound
	}
`
	}

	body := `package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

` + authzImport + `	"` + g.Module + `/internal/events"
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
` + ownerCheck + `
	from := item.` + goField + `
	transition, err := workflow.` + names.Pascal + `Workflow.Check(from, action)
	if err != nil {
		return nil, err
	}
	if transition.Permission != "" && can != nil && !can(transition.Permission) {
		return nil, workflow.ErrNotPermitted{
			Resource: "` + names.Plural + `", Action: action, Label: transition.Label, Noun: "` + names.Lower + `",
		}
	}

	// The status change, the hooks and the durable copy of the event commit
	// together or not at all. A hook that reserves stock and then finds the
	// budget short takes the reservation back out with the approval.
	var ev events.Event
	err = db.Transaction(func(tx *gorm.DB) error {
		// Guarded by the current state as well as the id. Two people pressing
		// Send at the same moment would otherwise both pass the check above and
		// both write; this makes the second one affect no rows, before any hook
		// has run for it.
		result := tx.Model(&models.` + names.Pascal + `{}).
			Where("id = ? AND ` + col + ` = ?", id, from).
			Update("` + col + `", transition.To)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return workflow.ErrInvalidTransition{
				Resource: "` + names.Plural + `",
				From:     from,
				Action:   action,
				Allowed:  actionsFrom(workflow.` + names.Pascal + `Workflow, from),
			}
		}
		if err := tx.First(&item, "id = ?", id).Error; err != nil {
			return err
		}

		// Hooks registered with workflow.OnTransition do the rest of the move's
		// work here, in this transaction, and can refuse it.
		if err := workflow.RunHooks(tx, workflow.Move{
			Resource: "` + names.Plural + `", Action: action, From: from, To: transition.To,
			ID: item.ID, Record: &item, C: c,
		}); err != nil {
			return err
		}

		// The transition is its own event, not an "updated". A subscriber that
		// cares about invoices being paid should not have to diff two versions
		// of a record to find out that is what happened. Durable subscribers
		// get it in this transaction, so it reaches them if and only if the
		// move commits.
		ev = events.Event{
			Name:     "` + names.Plural + `." + action,
			Resource: "` + names.Plural + `",
			Entity:   "` + names.Pascal + `",
			ID:       item.ID,
			Label:    ` + identExpr(names, g.Definition) + `,
			Detail:   fmt.Sprintf("%s: %s to %s", transition.Label, from, transition.To),
			Before:   map[string]interface{}{"` + col + `": from},
			After:    item,
		}
		return events.EmitTx(tx, c, &ev)
	})
	if err != nil {
		return nil, err
	}

	// The activity feed, realtime and webhooks hear about it once it has
	// committed.
	events.Emit(c, ev)
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

	// A project with a route file per resource mounts both routes in that file.
	// This injection is the fallback for a project from before route files, and
	// it used to run regardless: it wrote referralHandler.Transition into
	// routes.go, where no referralHandler exists, and the API stopped
	// compiling. A re-run takes those lines back out.
	if fileExists(filepath.Join(g.APIRoot(), "internal", "routes", names.Snake+"_routes.go")) {
		if !strings.Contains(content, names.Camel+"Handler.Transition") {
			return nil
		}
		for _, stale := range []string{names.Camel + "Handler.Workflow)", names.Camel + "Handler.Transition)"} {
			if err := removeLinesContaining(path, stale); err != nil {
				fmt.Printf("  Could not remove %s from routes.go: %v\n", stale, err)
			}
		}
		manifest.Refresh(path)
		fmt.Println("  ✓ Removed the transition routes an earlier run put in routes.go")
		return nil
	}

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
