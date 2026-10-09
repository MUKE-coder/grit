package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The README is the first file anybody opens, and for an --api project every
// command in its Quick Start failed.
//
// It was one fixed triple-tier page whatever shape was generated: `pnpm
// install` and `pnpm dev` on a project with no package.json, `docker compose
// up -d` for PostgreSQL on a SQLite project, a structure naming apps/web,
// apps/admin and turbo.json, and a Services table listing a web app and an
// admin panel that were not there.
//
// `grit new` already prints the right steps per shape, so the knowledge
// existed and only the README did not ask for it.
//
// Found by generating every remaining architecture and reading what each one
// tells its developer to do first.

func readmeFor(opts Options) string {
	opts.ProjectName = "app"
	opts.Version = "3.0.0"
	return readmeFile(opts)
}

// An API-only project is told nothing about a frontend it does not have.
func TestAnAPIOnlyReadmeDoesNotMentionAFrontend(t *testing.T) {
	src := readmeFor(Options{Architecture: ArchAPI, DBProvider: "sqlite"})

	for _, forbidden := range []string{
		"pnpm install",
		"pnpm dev",
		"apps/web",
		"apps/admin",
		"turbo.json",
		"Turborepo",
	} {
		if strings.Contains(src, forbidden) {
			t.Errorf("the --api README mentions %q, and an --api project has no such thing",
				forbidden)
		}
	}
	if !strings.Contains(src, "grit start server") {
		t.Error("the --api README does not say how to start the API, which is the only " +
			"thing it can start")
	}
}

// A SQLite project is not told to start PostgreSQL.
func TestASQLiteReadmeDoesNotStartADatabaseServer(t *testing.T) {
	src := readmeFor(Options{Architecture: ArchTriple, Frontend: FrontendNext, DBProvider: "sqlite"})

	if strings.Contains(src, "docker compose up -d") {
		t.Error("the README tells a SQLite project to start services in Docker; there are none")
	}
	if strings.Contains(src, "PostgreSQL") {
		t.Error("the README names PostgreSQL on a SQLite project")
	}
	if !strings.Contains(src, "SQLite") {
		t.Error("the README does not name the database the project actually uses")
	}
	// And the cloud-services section is about replacing Docker, which this
	// project is not using.
	if strings.Contains(src, "Neon") {
		t.Error("the README offers a cloud Postgres to a project that has no Postgres")
	}
}

// A Postgres project keeps the Docker instructions, because it needs them.
func TestAPostgresReadmeKeepsItsServices(t *testing.T) {
	src := readmeFor(Options{Architecture: ArchTriple, Frontend: FrontendNext, DBProvider: "postgres"})

	for _, want := range []string{"docker compose up -d", "PostgreSQL", "Redis", "Neon"} {
		if !strings.Contains(src, want) {
			t.Errorf("the Postgres README does not mention %q", want)
		}
	}
}

// A TanStack project is not described as a Next.js one.
func TestAViteReadmeNamesItsOwnFrontend(t *testing.T) {
	src := readmeFor(Options{Architecture: ArchTriple, Frontend: FrontendTanStack, DBProvider: "sqlite"})

	if strings.Contains(src, "Next.js") {
		t.Error("the README calls a TanStack project's frontend Next.js")
	}
	if !strings.Contains(src, "TanStack Router") {
		t.Error("the README does not name the router the project is built on")
	}
}

// And the extra apps appear only when they exist.
func TestTheReadmeListsOnlyTheAppsThatWereGenerated(t *testing.T) {
	plain := readmeFor(Options{Architecture: ArchTriple, Frontend: FrontendNext, DBProvider: "sqlite"})
	for _, forbidden := range []string{"apps/expo", "apps/desktop", "wails dev"} {
		if strings.Contains(plain, forbidden) {
			t.Errorf("a project without them mentions %q", forbidden)
		}
	}

	everything := readmeFor(Options{
		Architecture: ArchTriple, Frontend: FrontendNext, DBProvider: "postgres",
		IncludeExpo: true, IncludeDesktop: true,
	})
	for _, want := range []string{"apps/expo", "apps/desktop", "wails dev", "grit start expo"} {
		if !strings.Contains(everything, want) {
			t.Errorf("a project with them does not mention %q", want)
		}
	}
}

// A single project's Go module is at api/, not apps/api/.
func TestASingleProjectReadmePointsAtItsOwnLayout(t *testing.T) {
	src := readmeFor(Options{Architecture: ArchSingle, Frontend: FrontendTanStack, DBProvider: "sqlite"})

	if !strings.Contains(src, "cd api && air") {
		t.Error("the single-project README does not point at api/, where its Go module is")
	}
	if strings.Contains(src, "cd apps/api") {
		t.Error("the single-project README points at apps/api, which it does not have")
	}
}

// nextToTanStack rewrites next/link and friends to "@/lib/next-compat", so
// every project that runs a component through it has to have that module.
//
// A single Vite project put its only copy under the admin panel, at
// @admin/lib/next-compat, and converted one site component, which then
// imported a module that was not there. The project did not type-check, and
// nothing noticed because that shape had no type-check script either.
func TestASingleViteProjectShipsTheCompatShimWhereItIsImported(t *testing.T) {
	root := t.TempDir()
	opts := Options{
		ProjectName:  "app",
		Architecture: ArchSingle,
		Frontend:     FrontendTanStack,
		DBProvider:   "sqlite",
	}
	if err := writeSingleFrontendFiles(root, opts); err != nil {
		t.Fatalf("writeSingleFrontendFiles: %v", err)
	}

	mega := readTestFileAt(t, filepath.Join(root, "src", "components", "mega-menu.tsx"))
	if !strings.Contains(mega, `"@/lib/next-compat"`) {
		t.Skip("the converter no longer points at @/lib/next-compat; this test is about " +
			"the module existing where it points")
	}
	if _, err := os.Stat(filepath.Join(root, "src", "lib", "next-compat.tsx")); err != nil {
		t.Errorf("mega-menu.tsx imports @/lib/next-compat and the module is not there: %v. "+
			"The project does not type-check.", err)
	}
}
