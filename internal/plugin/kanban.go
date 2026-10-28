package plugin

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

func init() { Register(kanbanPlugin()) }

// kanbanRankSource is the ordering algorithm, emitted into the project
// verbatim.
//
// It is embedded rather than written out as a string literal so that the code
// the project receives is byte-for-byte the code Grit's own test suite runs.
// A plugin that keeps its algorithm in a Go string has no tests for what it
// ships: the string compiles nowhere until somebody installs it. Here
// internal/plugin/kanban/rank_test.go tests these exact bytes on every CI run,
// including the property that a thousand drags leave the board in the order
// the user arranged.
//
//go:embed kanban/rank.go
var kanbanRankSource string

// kanbanPlugin adds boards, columns and cards with drag-and-drop reordering.
//
// It exists as much to be read as to be installed. A Kanban board is the
// smallest feature that needs every part of the plugin system at once: models
// and migrations, a service with real logic, routes, React Query hooks, two
// admin pages, a nav entry, an npm dependency and a set of injections. The
// authoring guide walks through this file.
//
// The one piece worth stealing even if you never install it is the ordering.
// See internal/plugin/kanban/rank.go for why a position column is the wrong
// answer and what to do instead.
func kanbanPlugin() Plugin {
	return Plugin{
		Name:    "kanban",
		Version: "1.0.0",
		Summary: "Boards, columns and cards with drag-and-drop ordering",
		Description: `A working Kanban board in the admin: boards, columns, cards, and dragging
cards between columns with the pointer, touch or the keyboard.

  • Board / BoardColumn / BoardCard models, with migrations
  • A move endpoint that reorders by writing one row, with no renumbering
    and no transaction held across other people's cards
  • Optional work-in-progress limit per column
  • Drag with a pointer, a finger, or the keyboard, with screen reader
    announcements, via @dnd-kit
  • Emits board.card.moved on the domain event bus, so webhooks, the audit
    log and realtime subscribers all hear about it

Boards are visible to every signed-in user and record who created them. If
boards need to be per-tenant, install the multitenant plugin, which scopes
every model including these.

Requires an admin app (triple or full).`,

		Preflight: kanbanPreflight,

		NodeDeps: []Dependency{
			// The admin has no drag-and-drop primitive, and hand-rolled
			// HTML5 drag events do not fire on touch devices and cannot be
			// driven from the keyboard at all. dnd-kit gives all three plus
			// the live-region announcements a screen reader needs, which is
			// the whole reason to take a dependency rather than write 200
			// lines of pointer maths.
			{Name: "@dnd-kit/core", Version: "6.3.1", Workspace: "apps/admin"},
			{Name: "@dnd-kit/sortable", Version: "10.0.0", Workspace: "apps/admin"},
			{Name: "@dnd-kit/utilities", Version: "3.2.2", Workspace: "apps/admin"},
		},

		NextSteps: []string{
			"Install the new admin dependency:  pnpm install",
			"Run migrations:                    cd apps/api && go run cmd/migrate/main.go",
			"Open the admin and click Boards in the sidebar",
		},

		Files:      kanbanFiles,
		Injections: kanbanInjections,
	}
}

// kanbanPreflight refuses the install rather than writing files that cannot
// work.
//
// The admin pages are most of this plugin, and an --api or --mobile project
// has nowhere to put them. Saying so before the first write is the difference
// between "this project has no admin app" and a half-installed project whose
// Go code compiles and whose feature is invisible.
func kanbanPreflight(ctx Context) error {
	if !hasAdmin(ctx) {
		return fmt.Errorf("kanban needs an admin app, and this project is --%s\n\n"+
			"The board UI lives in the admin. Scaffold with --triple, or install the\n"+
			"API half by hand from the files this plugin would have written.", ctx.Architecture)
	}
	// ids.New is what every model in a Grit project uses for its primary key.
	// A project old enough to predate it would produce boards with empty ids,
	// which is a confusing failure at runtime rather than at install.
	idsPkg := filepath.Join(ctx.Root, apiPath(ctx, "internal/ids"))
	if _, err := os.Stat(idsPkg); err != nil {
		return fmt.Errorf("this project has no internal/ids package, which kanban's models need for their primary keys\n\n" +
			"Run `grit upgrade` first.")
	}
	return nil
}

func kanbanFiles(ctx Context) map[string]string {
	files := map[string]string{
		// The ordering algorithm, verbatim from internal/plugin/kanban/rank.go.
		apiPath(ctx, "internal/kanban/rank.go"): kanbanRankSource,

		apiPath(ctx, "internal/models/board.go"):        kanbanModelsGo(ctx),
		apiPath(ctx, "internal/services/board.go"):      kanbanServiceGo(ctx),
		apiPath(ctx, "internal/handlers/board.go"):      kanbanHandlerGo(ctx),
		apiPath(ctx, "internal/handlers/board_test.go"): kanbanHandlerTestGo(ctx),
	}

	// Guarded by Preflight, so this is belt and braces rather than a branch
	// that is expected to be taken.
	if !hasAdmin(ctx) {
		return files
	}

	files[adminFile(ctx, "hooks/use-boards.ts")] = adminSource(ctx, kanbanHook())
	files[adminFile(ctx, "components/kanban/board-view.tsx")] = adminSource(ctx, kanbanBoardView())
	files[adminFile(ctx, "components/kanban/card-composer.tsx")] = adminSource(ctx, kanbanCardComposer())

	list, listRoute := adminPageFile(ctx, "boards", "boards")
	files[list] = adminSource(ctx, kanbanBoardsPage())
	if listRoute != "" {
		files[listRoute] = tanStackRouteShim("boards", "@/pages/boards")
	}

	detail, detailRoute := adminPageFile(ctx, "boards/[id]", "board-detail")
	if ctx.Frontend == "tanstack" {
		// TanStack spells a path parameter $id, not [id].
		detail, detailRoute = adminPageFile(ctx, "boards/$id", "board-detail")
	}
	files[detail] = adminSource(ctx, kanbanBoardDetailPage(ctx))
	if detailRoute != "" {
		files[detailRoute] = tanStackRouteShim("boards/$id", "@/pages/board-detail")
	}

	return files
}

func kanbanInjections(ctx Context) []Injection {
	injections := []Injection{
		{
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&Board{},\n\t\t&BoardColumn{},\n\t\t&BoardCard{},",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\tboardHandler := handlers.NewBoardHandler(db)",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: `		// Kanban boards. Every route is behind auth: a board is team data,
		// not public. See internal/kanban/rank.go for how /move orders cards.
		protected.GET("/boards", boardHandler.ListBoards)
		protected.POST("/boards", boardHandler.CreateBoard)
		protected.GET("/boards/:id", boardHandler.GetBoard)
		protected.PATCH("/boards/:id", boardHandler.UpdateBoard)
		protected.DELETE("/boards/:id", boardHandler.DeleteBoard)

		protected.POST("/boards/:id/columns", boardHandler.CreateColumn)
		protected.POST("/boards/:id/cards", boardHandler.CreateCard)

		protected.PATCH("/board-columns/:id", boardHandler.UpdateColumn)
		protected.DELETE("/board-columns/:id", boardHandler.DeleteColumn)
		protected.POST("/board-columns/:id/move", boardHandler.MoveColumn)

		protected.PATCH("/board-cards/:id", boardHandler.UpdateCard)
		protected.DELETE("/board-cards/:id", boardHandler.DeleteCard)
		protected.POST("/board-cards/:id/move", boardHandler.MoveCard)`,
		},
	}

	if !hasAdmin(ctx) {
		return injections
	}

	admin := adminDir(ctx)
	return append(injections,
		// getIcon falls back to a document icon for a key it does not know, so
		// a nav entry whose icon is missing from both of these renders wrong
		// and nothing reports it.
		Injection{
			File:   admin + "/lib/icons.ts",
			Marker: "// grit:icons:import",
			Code:   "  SquareKanban,",
		},
		Injection{
			File:   admin + "/lib/icons.ts",
			Marker: "// grit:icons:map",
			Code:   "  SquareKanban,",
		},
		Injection{
			File:   admin + "/components/chrome/CollapsibleSidebar.tsx",
			// Content, not system. The board is where the work happens, so it
			// belongs beside the resources in the rail. grit:nav:system would
			// have filed it under System, Extensions, three clicks from the
			// dashboard, which is where the first version of this plugin put
			// it and where nobody found it.
			Marker: "// grit:nav:content",
			// adminOnly is false: a board is for whoever is doing the work,
			// not an administrative screen.
			Code: `  { href: "/boards", label: "Boards", iconKey: "SquareKanban", adminOnly: false },`,
		},
	)
}
