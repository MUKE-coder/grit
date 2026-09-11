package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// UpgradeOptions holds the configuration for upgrading a project.
type UpgradeOptions struct {
	// Force overwrites every file, including ones you have edited. This is
	// what upgrade did unconditionally before v3.147.0.
	Force bool
	// ShowDiff prints what the upgrade would have changed in the files it
	// left alone.
	ShowDiff bool
	// Version is the Grit version performing the upgrade, recorded against
	// every file it rewrites.
	Version string
}

// Upgrade regenerates generic scaffold files in an existing Grit project.
// It preserves user-generated code (resource definitions, API handlers, .env)
// while updating framework components to the latest version.
func Upgrade(uOpts UpgradeOptions) error {
	root, err := FindProjectRoot()
	if err != nil {
		return err
	}

	projectName, err := readProjectName(root)
	if err != nil {
		return err
	}

	spinner := color.New(color.FgHiBlack)
	green := color.New(color.FgHiGreen)
	cyan := color.New(color.FgHiCyan)

	opts := Options{
		ProjectName: projectName,
		Style:       readProjectStyle(root),
		Version:     uOpts.Version,
		Frontend:    readProjectFrontend(root),
	}
	opts.Normalize()

	// Two pieces of machinery, both process-wide for the length of the
	// upgrade. The guard decides what may be overwritten; the recorder notes
	// what was, so the next upgrade knows.
	if err := startGuard(root, uOpts.Force); err != nil {
		return err
	}
	release, err := manifest.Start(root, opts.Version, "scaffold")
	if err != nil {
		return err
	}
	// Every step below can return early. Without this, a failure halfway
	// through would leave the files it had already written unrecorded, and the
	// next upgrade would treat all of them as hand-edited.
	finished := false
	defer func() {
		if !finished {
			stopGuard()
			_ = release()
		}
	}()

	// Detect which apps exist
	hasWeb := dirExists(filepath.Join(root, "apps", "web"))
	hasAdmin := dirExists(filepath.Join(root, "apps", "admin"))
	hasExpo := dirExists(filepath.Join(root, "apps", "expo"))
	hasDocs := dirExists(filepath.Join(root, "apps", "docs"))
	hasShared := dirExists(filepath.Join(root, "packages", "shared"))

	var updated int

	// --- Root config files ---
	spinner.Printf("  → Updating root configuration...\n")
	rootFiles := map[string]string{
		filepath.Join(root, "turbo.json"):          turboJSON(),
		filepath.Join(root, "pnpm-workspace.yaml"): pnpmWorkspace(dirExists(filepath.Join(root, "apps", "desktop", "frontend"))),
		filepath.Join(root, ".npmrc"):              rootNpmrc(),
	}
	n, err := writeUpgradeFiles(rootFiles, uOpts.Force)
	if err != nil {
		return fmt.Errorf("updating root files: %w", err)
	}
	updated += n

	// --- React version pin ---
	// react and react-dom must be the exact same version (React 19 hard-errors
	// otherwise). Older scaffolds left them at "^19.0.0" and relied on a pnpm
	// override that pnpm 10 doesn't apply to react-dom — so they drift onto
	// different 19.x lines and the app white-screens. Surgically pin both in
	// every frontend package.json without clobbering the user's other deps.
	pinReactVersions(root, spinner)

	// --- Docker files ---
	spinner.Printf("  → Updating Docker configuration...\n")
	dockerFiles := map[string]string{
		filepath.Join(root, "docker-compose.yml"):      dockerCompose(opts),
		filepath.Join(root, "docker-compose.prod.yml"): dockerComposeProd(opts),
	}
	n, err = writeUpgradeFiles(dockerFiles, uOpts.Force)
	if err != nil {
		return fmt.Errorf("updating Docker files: %w", err)
	}
	updated += n

	// --- API migrate/seed tools ---
	hasAPI := dirExists(filepath.Join(root, "apps", "api"))
	if hasAPI {
		spinner.Printf("  → Updating migration and seed tools...\n")
		if err := writeMigrateSeedFiles(root, opts); err != nil {
			return fmt.Errorf("updating migrate/seed files: %w", err)
		}
		if err := writeCodegenRuntimeFiles(root, opts); err != nil {
			return fmt.Errorf("updating codegen runtime packages: %w", err)
		}
		if err := writeFrameworkOwnedFiles(root, opts); err != nil {
			return fmt.Errorf("updating framework models: %w", err)
		}
		if err := ensureRouteRegistry(root, opts); err != nil {
			return fmt.Errorf("adding the route registry: %w", err)
		}
		green.Printf("  ✓ Migration and seed tools updated\n")
		updated += 4

		// The media pipeline, and the storage/upload code it hooks into.
		//
		// Narrow on purpose. Upgrade does not regenerate API code in general,
		// and widening it here would be that project rather than this feature.
		// These are listed because the pipeline is useless without them:
		// internal/media is new, and the upload handler is where the transform
		// actually happens.
		//
		// writeFile is manifest-guarded, so a file the reader has edited is
		// reported as a conflict rather than overwritten.
		spinner.Printf("  → Updating the media pipeline...\n")
		if err := writeUploadPackageFiles(root, opts); err != nil {
			return fmt.Errorf("updating upload package: %w", err)
		}
		if err := writeAdminSecurityFiles(root, opts); err != nil {
			return fmt.Errorf("updating security page: %w", err)
		}
		if err := ensurePasskeyWiring(root, opts); err != nil {
			return fmt.Errorf("wiring passkeys: %w", err)
		}
		if err := ensureRecoveryWiring(root, opts); err != nil {
			return fmt.Errorf("wiring recovery contacts: %w", err)
		}
		if err := writeAdminPasskeyFiles(root, opts); err != nil {
			return fmt.Errorf("updating passkey UI: %w", err)
		}
		if err := writePasskeyFiles(root, opts); err != nil {
			return fmt.Errorf("updating passkey files: %w", err)
		}
		if err := writeRecoveryFiles(root, opts); err != nil {
			return fmt.Errorf("updating recovery files: %w", err)
		}
		if err := writeMoneyFiles(root, opts); err != nil {
			return fmt.Errorf("updating money files: %w", err)
		}
		// Same reason as money: the generator writes calls into this package,
		// so an older copy of it means generated code that does not compile.
		if err := writeRespondFiles(root, opts); err != nil {
			return fmt.Errorf("updating respond files: %w", err)
		}
		// The generator writes calls into appendonly for --append-only, and a
		// model registering with a guard the project never installs protects
		// nothing. So the package, and the two calls that install it.
		if err := writeAppendOnlyFiles(root, opts); err != nil {
			return fmt.Errorf("updating append-only files: %w", err)
		}
		if err := EnsureAppendOnlyWiring(opts.APIRoot(root), opts.Module()); err != nil {
			fmt.Printf("  ⚠ %v\n", err)
		}
		// The leak fix for encrypted columns lives in the crypto package and a
		// call at connect time; an older project has neither.
		if err := writeCryptoFiles(root, opts); err != nil {
			return fmt.Errorf("updating crypto files: %w", err)
		}
		if err := EnsureFieldEncryptionWiring(opts.APIRoot(root), opts.Module()); err != nil {
			fmt.Printf("  ⚠ %v\n", err)
		}
		// Restore replays in foreign-key order and lets itself past the
		// append-only triggers. Without it a project with --items or
		// --append-only has backups it cannot restore, which nobody finds out
		// until the day they need one. Only refreshed where the package is
		// already there: a project older than backups may lack what it imports.
		if fileExists(filepath.Join(opts.APIRoot(root), "internal", "backup", "restore.go")) {
			if err := writeBackupFiles(root, opts); err != nil {
				return fmt.Errorf("updating backup files: %w", err)
			}
		}
		// Erasure became a registry so owned resources are erased with their
		// owner, and restore re-applies later erasures through it. Only where the
		// GDPR service is already there.
		if fileExists(filepath.Join(opts.APIRoot(root), "internal", "services", "gdpr.go")) {
			if err := writeErasureFiles(root, opts); err != nil {
				return fmt.Errorf("updating erasure files: %w", err)
			}
		}
		// Resource handlers are the developer's and are not regenerated, but
		// generated text that is still recognisable gets two fixes: owned
		// resources check the owner on every path that reaches a row, and CSV
		// export returns its rows. Anything too changed to patch is named.
		if err := repairResourceHandlers(root, opts); err != nil {
			fmt.Printf("  ⚠ repairing resource handlers: %v\n", err)
		}
		// The activity-log chain was hashed at a precision Postgres and MySQL
		// do not store, so it failed verification on its first row. Only where
		// the audit package is already there.
		if fileExists(filepath.Join(opts.APIRoot(root), "internal", "audit", "audit.go")) {
			if err := writeAuditChainFiles(root, opts); err != nil {
				return fmt.Errorf("updating the audit chain: %w", err)
			}
		}
		// Replicas behind a load balancer kept separate permission caches and
		// SSO registries, and each ran every scheduled job. The package and the
		// files that use it go first: the SSO files below import it.
		if err := writeClusterFiles(root, opts); err != nil {
			return fmt.Errorf("updating replica coordination: %w", err)
		}
		// Generated update routes call this package for If-Match; a resource
		// generated after this release does not compile without it.
		if err := WriteConcurrencyPackage(opts.APIRoot(root), opts.Module(), true); err != nil {
			return fmt.Errorf("updating the concurrency package: %w", err)
		}
		// Enterprise SSO trusted a customer's identity provider for any email
		// address, so it could sign in as your own administrator. Only where
		// SSO is already there.
		if fileExists(filepath.Join(opts.APIRoot(root), "internal", "handlers", "sso.go")) {
			if err := writeSSOFiles(root, opts); err != nil {
				return fmt.Errorf("updating SSO: %w", err)
			}
		}
		if err := ensureClusterWiring(opts.APIRoot(root), opts.Module()); err != nil {
			fmt.Printf("  ⚠ %v\n", err)
		}
		// Custom roles reach the admin API through the staff group a new
		// routes.go has. An existing routes.go is not restructured, so say so.
		if f := filepath.Join(opts.APIRoot(root), "internal", "routes", "routes.go"); fileExists(f) && !fileContains(f, "RequireStaff()") {
			fmt.Println("  • Custom roles cannot reach the admin API in this project yet: its routes.go\n" +
				"    predates the staff group. The Roles & permissions docs show the change (Staff routes).")
		}
		// Read-only. .env holds secrets and is the developer's, so a URL that
		// disagrees with the port compose binds is reported, not rewritten.
		for _, w := range envPortDrift(root) {
			fmt.Printf("  ⚠ %s\n", w)
		}
		if err := writeOutboxFiles(root, opts); err != nil {
			return fmt.Errorf("updating outbox files: %w", err)
		}
		if err := writeStockFiles(root, opts); err != nil {
			return fmt.Errorf("updating stock files: %w", err)
		}
		if err := ensureMoneyFrontend(root, opts); err != nil {
			return fmt.Errorf("wiring money frontend: %w", err)
		}
		if err := writeJSONTimeFiles(root, opts); err != nil {
			return fmt.Errorf("updating jsontime files: %w", err)
		}
		if err := writeMediaFiles(root, opts); err != nil {
			return fmt.Errorf("updating media files: %w", err)
		}
		// The browser half of realtime: one shared socket, reconnection, and
		// the React Query bindings. The hub has shipped for a long time with
		// nothing on the client able to consume it.
		if err := writeRealtimeClientFiles(root, opts); err != nil {
			return err
		}
		if err := writeStorageFiles(root, opts); err != nil {
			return fmt.Errorf("updating storage files: %w", err)
		}
		green.Printf("  ✓ Media pipeline updated\n")
		green.Printf("    Then run: cd apps/api && go mod tidy\n")
		updated += 3

		// API documentation (gin-docs — now configured in routes.go)
		green.Printf("  ✓ API documentation (gin-docs) configured in routes.go\n")
	}

	// --- Shared package ---
	if hasShared {
		spinner.Printf("  → Updating shared package config...\n")
		sharedFiles := map[string]string{
			filepath.Join(root, "packages", "shared", "package.json"):  sharedPackageJSON(opts),
			filepath.Join(root, "packages", "shared", "tsconfig.json"): sharedTSConfig(),
		}
		n, err = writeUpgradeFiles(sharedFiles, uOpts.Force)
		if err != nil {
			return fmt.Errorf("updating shared files: %w", err)
		}
		updated += n
	}

	// --- Web app (landing page — always safe to fully regenerate) ---
	if hasWeb && opts.Frontend == FrontendTanStack {
		skipViteApp("web")
	} else if hasWeb {
		spinner.Printf("  → Updating web app (landing page)...\n")
		if err := writeWebFiles(root, opts); err != nil {
			return fmt.Errorf("updating web files: %w", err)
		}
		green.Printf("  ✓ Web app updated\n")
		updated += 9
	}

	// --- Admin panel (generic components only — preserves resource definitions) ---
	if hasAdmin && opts.Frontend == FrontendTanStack {
		skipViteApp("admin")
	} else if hasAdmin {
		// Before anything is written: resources moved from resources/<name>.ts
		// to resources/<name>/<name>.ts in v3.143.0. Writing users/users.ts into
		// a project still holding a flat users.ts would leave two definitions and
		// a registry importing the stale one.
		migrateResourceFolders(root, spinner, green)
		spinner.Printf("  → Updating admin panel components...\n")
		n, err := upgradeAdminFiles(root, opts, uOpts)
		if err != nil {
			return fmt.Errorf("updating admin files: %w", err)
		}
		updated += n
	}

	// --- shadcn config for every frontend ---
	//
	// Created when missing, never overwritten. Without it,
	// "npx shadcn@latest add https://ui.gritframework.dev/r/<block>.json"
	// abandons the install and starts an interactive setup instead, asking a
	// project that already knows all the answers to answer them again.
	//
	// Its own step rather than part of the app file maps because there is no
	// upgrade path for apps/web at all: this is a config file, so creating a
	// missing one is safe in a way that rewriting a page never is.
	updated += ensureShadcnConfigs(root, spinner, green)

	// --- Docs ---
	if hasDocs {
		spinner.Printf("  → Updating documentation...\n")
		if err := writeDocsFiles(root, opts); err != nil {
			return fmt.Errorf("updating docs files: %w", err)
		}
		green.Printf("  ✓ Documentation updated\n")
		updated += 15
	}

	// --- Expo (if exists) ---
	if hasExpo {
		spinner.Printf("  → Updating Expo app...\n")
		opts.IncludeExpo = true
		if err := writeExpoFiles(root, opts); err != nil {
			return fmt.Errorf("updating Expo files: %w", err)
		}
		green.Printf("  ✓ Expo app updated\n")
		updated += 10
	}

	written, skipped := stopGuard()
	if saveErr := release(); saveErr != nil {
		spinner.Printf("  Could not write .grit/manifest.json: %v\n", saveErr)
	}
	finished = true
	_ = updated // superseded by the guard's count, which counts actual writes

	// Record the version in grit.json, which is where people look for it.
	// Never fatal: the upgrade itself has already happened, and failing here
	// over a label would be the wrong trade.
	if err := stampProjectVersion(root, uOpts.Version); err != nil {
		spinner.Printf("  Could not update grit.json: %v\n", err)
	}

	fmt.Println()
	green.Printf("  ✓ Upgrade complete. Updated %d files.\n", written)

	if len(skipped) > 0 {
		fmt.Println()
		yellow := color.New(color.FgHiYellow)
		yellow.Printf("  ⚠ Left alone, because you have edited %s since Grit wrote %s:\n",
			plural(len(skipped), "this file", "these files"),
			plural(len(skipped), "it", "them"))
		for _, s := range skipped {
			fmt.Printf("      %s\n", s.Rel)
		}
		fmt.Println()
		if uOpts.ShowDiff {
			for _, s := range skipped {
				cyan.Printf("  ── %s ──\n", s.Rel)
				fmt.Println(diffIndent(s.Diff()))
			}
		} else {
			cyan.Println("    grit upgrade --diff     # see what the new version would change")
		}
		cyan.Println("    grit upgrade --force    # take the new version and lose your edits")
	}

	fmt.Println()
	cyan.Println("  Next steps:")
	if hasWeb || hasAdmin {
		cyan.Println("    pnpm install    # Install any new dependencies")
		cyan.Println("    pnpm dev        # Restart development servers")
	}
	fmt.Println()

	spinner.Println("  Note: Resource definitions and API code were preserved.")
	spinner.Println("  Run 'grit sync' to regenerate TypeScript types from Go models.")
	fmt.Println()

	return nil
}

// plural picks a word based on a count, so the summary reads as English rather
// than as "1 file(s)".
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// diffIndent pads a unified diff so it lines up with the rest of the CLI
// output without disturbing the leading +/- that makes it readable.
func diffIndent(diff string) string {
	lines := strings.Split(strings.TrimRight(diff, "\n"), "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
}

// migrateResourceFolders gives every resource its own folder, for both the
// Next.js and the Vite admin. A failure is a warning rather than a stopped
// upgrade, because the rest of it is still worth having.
func migrateResourceFolders(root string, spinner, green *color.Color) {
	for _, rel := range [][]string{
		{"apps", "admin", "resources"},
		{"apps", "admin", "src", "resources"},
	} {
		dir := filepath.Join(append([]string{root}, rel...)...)
		if !dirExists(dir) {
			continue
		}
		moved, err := MigrateResourceLayout(dir)
		if err != nil {
			spinner.Printf("  Could not reorganise %s: %v\n", filepath.Join(rel...), err)
			continue
		}
		if moved > 0 {
			green.Printf("  Moved %d resource(s) into their own folders\n", moved)
		}
	}
}

// upgradeAdminFiles regenerates admin panel generic files while preserving
// user-created resource definitions and pages.
func upgradeAdminFiles(root string, opts Options, uOpts UpgradeOptions) (int, error) {
	adminRoot := filepath.Join(root, "apps", "admin")
	green := color.New(color.FgHiGreen)

	// One list, shared with the scaffold. Keeping a second copy here is what
	// left 59 files unreachable by upgrade, including a page that shipped
	// broken: the fix went out, and upgrade had no idea the file existed.
	//
	// Everything framework-owned is refreshed. Two things protect the reader's
	// work. The manifest guard at writeFile refuses to overwrite a file whose
	// contents are not exactly what Grit last wrote there, so local edits
	// survive whatever is in this map. And userOwnedAdminFiles below excludes
	// the handful of files a person is *expected* to edit, so they are not even
	// offered as conflicts on every upgrade.
	files := map[string]string{}
	for path, body := range adminFileMap(root, opts) {
		if isUserOwnedAdminFile(adminRoot, path) {
			continue
		}
		files[path] = body
	}

	n, err := writeUpgradeFiles(files, uOpts.Force)
	if err != nil {
		return 0, err
	}

	// resources/users.ts now imports ./users.custom, so the overlay has to be
	// there or the admin will not compile after an upgrade. createIfMissing,
	// never write: the whole promise of that file is that it is never
	// overwritten, and an upgrade is exactly when someone would lose work.
	created, err := createIfMissing(
		filepath.Join(adminRoot, "resources", "users", "users.custom.tsx"),
		AdminResourceCustomStub("User", "User"),
	)
	if err != nil {
		return n, err
	}
	if created {
		n++
	}

	n += pruneAdminStrays(adminRoot)
	n += pruneUnwiredI18n(adminRoot)

	green.Printf("  ✓ Admin panel updated (%d files)\n", n)
	return n, nil
}

// Files earlier upgrades wrote to paths the scaffold does not use.
//
// The admin components were renamed to kebab-case a long time ago, but this
// command's path list was not, so every `grit upgrade` since has written
// components/tables/DataTable.tsx next to the real data-table.tsx and left it
// there. Nothing imported them. The visible symptom was the opposite of an
// error: upgrade reported dozens of files updated while the components the app
// actually renders were never touched, so no fix shipped in one ever arrived.
//
// They are safe to delete precisely because nothing ever imported them — but
// only when the real file is present, so a half-finished rename cannot take the
// last copy with it.
var adminStrayFiles = map[string]string{
	"components/tables/DataTable.tsx":           "components/tables/data-table.tsx",
	"components/tables/ColumnHeader.tsx":        "components/tables/column-header.tsx",
	"components/tables/CellRenderers.tsx":       "components/tables/cell-renderers.tsx",
	"components/tables/TableFilters.tsx":        "components/tables/table-filters.tsx",
	"components/tables/TableToolbar.tsx":        "components/tables/table-toolbar.tsx",
	"components/tables/TablePagination.tsx":     "components/tables/table-pagination.tsx",
	"components/tables/Skeleton.tsx":            "components/tables/table-skeleton.tsx",
	"components/tables/EmptyState.tsx":          "components/tables/table-empty-state.tsx",
	"components/tables/Formatters.ts":           "lib/formatters.ts",
	"components/forms/FormBuilder.tsx":          "components/forms/form-builder.tsx",
	"components/forms/FormModal.tsx":            "components/forms/form-modal.tsx",
	"components/forms/fields/TextField.tsx":     "components/forms/fields/text-field.tsx",
	"components/forms/fields/TextAreaField.tsx": "components/forms/fields/textarea-field.tsx",
	"components/forms/fields/NumberField.tsx":   "components/forms/fields/number-field.tsx",
	"components/forms/fields/SelectField.tsx":   "components/forms/fields/select-field.tsx",
	"components/forms/fields/DateField.tsx":     "components/forms/fields/date-field.tsx",
	"components/forms/fields/ToggleField.tsx":   "components/forms/fields/toggle-field.tsx",
	"components/forms/fields/CheckboxField.tsx": "components/forms/fields/checkbox-field.tsx",
	"components/forms/fields/RadioField.tsx":    "components/forms/fields/radio-field.tsx",
	"components/layout/AdminLayout.tsx":         "components/layout/admin-layout.tsx",
	"components/layout/Sidebar.tsx":             "components/layout/sidebar.tsx",
	"components/layout/Navbar.tsx":              "components/layout/navbar.tsx",
	"components/layout/ThemeProvider.tsx":       "components/shared/theme-provider.tsx",
	"components/shared/Providers.tsx":           "components/shared/providers.tsx",
	"components/widgets/StatsCard.tsx":          "components/widgets/stats-card.tsx",
	"components/widgets/ChartWidget.tsx":        "components/widgets/chart-widget.tsx",
	"components/widgets/ActivityWidget.tsx":     "components/widgets/activity-widget.tsx",
	"components/widgets/WidgetGrid.tsx":         "components/widgets/widget-grid.tsx",
	"components/resource/ResourcePage.tsx":      "components/resource/resource-page.tsx",
	"resources/registry.ts":                     "resources/index.ts",
	"postcss.config.mjs":                        "postcss.config.js",
}

// i18n files the scaffold used to write without the dependency that makes them
// compile. `grit add i18n` writes the same six files AND adds next-intl, wraps
// the layout and registers the plugin — and it skips files that already exist,
// so these copies did not merely fail to type-check, they blocked the command
// that would have fixed them. Removed only when next-intl is absent, which is
// exactly the case where they cannot be doing anything useful.
func pruneUnwiredI18n(adminRoot string) int {
	pkg, err := os.ReadFile(filepath.Join(adminRoot, "package.json"))
	if err != nil || strings.Contains(string(pkg), `"next-intl"`) {
		return 0
	}
	removed := 0
	for _, dead := range []string{
		"i18n/request.ts",
		"lib/locale.ts",
		"components/language-switcher.tsx",
		"messages/en.json",
		"messages/fr.json",
		"messages/sw.json",
	} {
		if err := os.Remove(filepath.Join(adminRoot, filepath.FromSlash(dead))); err == nil {
			removed++
		}
	}
	if removed > 0 {
		fmt.Printf("  ✓ Removed %d unwired i18n file(s): run `grit add i18n` to set it up properly\n", removed)
	}
	return removed
}

// existsExact reports whether path exists with exactly this spelling.
//
// os.Stat is not good enough here. Windows and macOS match file names without
// regard to case, so Stat("providers.tsx") happily returns the entry for
// "Providers.tsx" — and a prune that trusts it deletes the only copy while
// believing it kept one. Reading the directory and comparing names byte for
// byte is the only answer that means the same thing on every platform.
func existsExact(path string) bool {
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		return false
	}
	want := filepath.Base(path)
	for _, e := range entries {
		if e.Name() == want {
			return true
		}
	}
	return false
}

func pruneAdminStrays(adminRoot string) int {
	removed := 0
	for stray, real := range adminStrayFiles {
		// Pairs that differ only in case (Providers.tsx / providers.tsx) need
		// no special handling once existsExact is doing the looking: on a
		// case-insensitive filesystem only one of the two spellings is really
		// there, so the "is the real one present?" check below fails and the
		// single copy survives. On a case-sensitive one both exist and the
		// stray is genuinely a stray.
		strayPath := filepath.Join(adminRoot, filepath.FromSlash(stray))
		if !existsExact(strayPath) {
			continue
		}
		if !existsExact(filepath.Join(adminRoot, filepath.FromSlash(real))) {
			continue // the real one is missing — leave the copy alone
		}
		if err := os.Remove(strayPath); err == nil {
			removed++
		}
	}
	if removed > 0 {
		fmt.Printf("  ✓ Removed %d stale duplicate component file(s) from earlier upgrades\n", removed)
	}
	return removed
}

// ensureShadcnConfigs gives every frontend a components.json if it has none.
//
// The shadcn CLI refuses to install into a project without one, and Grit UI is
// distributed as a shadcn registry, so a project scaffolded before this shipped
// cannot install a block without answering four setup questions first, one of
// which offers a component library that is wrong for the project.
//
// Layout decides the contents: an App Router app declares rsc true and points
// at app/globals.css, a Vite app declares false and points at src/globals.css.
// Detected from the file that is actually there rather than from a flag, so a
// project that switched frontends still gets the right one.
func ensureShadcnConfigs(root string, spinner, green *color.Color) int {
	created := 0
	for _, app := range []string{"web", "admin"} {
		appRoot := filepath.Join(root, "apps", app)
		if !dirExists(appRoot) {
			continue
		}

		body := viteComponentsJSON()
		if fileExists(filepath.Join(appRoot, "app", "globals.css")) {
			body = nextComponentsJSON()
		}

		wrote, err := createIfMissing(filepath.Join(appRoot, "components.json"), body)
		if err == nil && wrote {
			created++
			green.Printf("  ✓ apps/%s/components.json (shadcn and Grit UI installs)\n", app)
		}
	}
	return created
}

// createIfMissing writes a file only when it is absent, and reports whether it
// did. Used for the files an upgrade must guarantee exist without ever
// clobbering their contents.
func createIfMissing(path, content string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("checking %s: %w", path, err)
	}
	if err := writeFile(path, content); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}

// writeUpgradeFiles writes files, creating directories as needed.
// Returns the number of files written.
func writeUpgradeFiles(files map[string]string, force bool) (int, error) {
	count := 0
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return count, fmt.Errorf("writing %s: %w", path, err)
		}
		count++
	}
	return count, nil
}

// pinReactVersions surgically pins react + react-dom to one exact version in
// every frontend package.json, leaving the rest of the file untouched. Older
// scaffolds used "^19.0.0" (which drifts) or a "19.1.0" pin (which pnpm 10
// applies to react only), producing a react/react-dom mismatch that
// white-screens the app.
func pinReactVersions(root string, spinner *color.Color) {
	const target = "19.2.7"
	replacer := strings.NewReplacer(
		`"react": "^19.0.0"`, `"react": "`+target+`"`,
		`"react-dom": "^19.0.0"`, `"react-dom": "`+target+`"`,
		`"react": "19.1.0"`, `"react": "`+target+`"`,
		`"react-dom": "19.1.0"`, `"react-dom": "`+target+`"`,
	)
	paths := []string{
		filepath.Join(root, "package.json"), // single-app scaffold
		filepath.Join(root, "apps", "admin", "package.json"),
		filepath.Join(root, "apps", "web", "package.json"),
		filepath.Join(root, "apps", "desktop", "frontend", "package.json"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		out := replacer.Replace(string(data))
		if out != string(data) {
			if err := os.WriteFile(p, []byte(out), 0o644); err == nil {
				spinner.Printf("  ✓ pinned react/react-dom to %s in %s\n", target, filepath.Base(filepath.Dir(p)))
			}
		}
	}
}

// FindProjectRoot walks up from the current directory looking for a Grit project.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}

	for {
		// grit.json is the universal project marker — present in every mode
		// (single / api-only / monorepo). Monorepo modes also have turbo.json +
		// apps/api, which we still accept as a fallback for older projects.
		if fileExists(filepath.Join(dir, "grit.json")) {
			return dir, nil
		}
		if fileExists(filepath.Join(dir, "turbo.json")) && dirExists(filepath.Join(dir, "apps", "api")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("not inside a Grit project (no grit.json found)\n\nRun this command from your Grit project root or any subdirectory")
}

// readProjectName reads the project name from the root package.json.
func readProjectName(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		// Fall back to directory name
		return filepath.Base(root), nil
	}

	var pkg struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return filepath.Base(root), nil
	}

	name := pkg.Name
	// Remove @scope/ prefix if present
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}

	if name == "" {
		return filepath.Base(root), nil
	}

	return name, nil
}

// readProjectStyle reads the style variant from grit.config.ts.
// Returns "default" if the file doesn't exist or the style isn't found.
func readProjectStyle(root string) string {
	data, err := os.ReadFile(filepath.Join(root, "grit.config.ts"))
	if err != nil {
		return "default"
	}

	re := regexp.MustCompile(`style:\s*"([^"]+)"`)
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return "default"
	}

	style := string(matches[1])
	for _, valid := range ValidStyles {
		if style == valid {
			return style
		}
	}

	return "default"
}

// fileExists returns true if a file exists.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists returns true if a directory exists.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isUserOwnedAdminFile reports whether a path is one the reader owns.
//
// These are excluded from upgrade even though the scaffold writes them, because
// they exist to be edited: a resource definition is where columns, filters and
// form fields are configured, and a .custom overlay is a promise that Grit will
// never touch it.
//
// The users resource definition is the one exception and is refreshed
// deliberately: it is the built-in one, it gains fields as the framework does,
// and the manifest guard still refuses to touch it if it has been edited.
func isUserOwnedAdminFile(adminRoot, path string) bool {
	rel, err := filepath.Rel(adminRoot, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)

	if rel == "resources/users/users.ts" {
		return false
	}
	// Every .custom overlay, whatever the resource.
	if strings.HasSuffix(rel, ".custom.tsx") {
		return true
	}
	// Resource definitions and their generated pages.
	if strings.HasPrefix(rel, "resources/") {
		return true
	}
	if strings.HasPrefix(rel, "app/(dashboard)/resources/") {
		return true
	}
	// The demo Blog resource, which people delete or rewrite.
	if strings.Contains(rel, "/blogs/") {
		return true
	}
	return false
}
