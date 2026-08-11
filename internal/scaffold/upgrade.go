package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// UpgradeOptions holds the configuration for upgrading a project.
type UpgradeOptions struct {
	// Force overwrites files without prompting
	Force bool
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
	}

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
		green.Printf("  ✓ Migration and seed tools updated\n")
		updated += 4

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
	if hasWeb {
		spinner.Printf("  → Updating web app (landing page)...\n")
		if err := writeWebFiles(root, opts); err != nil {
			return fmt.Errorf("updating web files: %w", err)
		}
		green.Printf("  ✓ Web app updated\n")
		updated += 9
	}

	// --- Admin panel (generic components only — preserves resource definitions) ---
	if hasAdmin {
		spinner.Printf("  → Updating admin panel components...\n")
		n, err := upgradeAdminFiles(root, opts, uOpts)
		if err != nil {
			return fmt.Errorf("updating admin files: %w", err)
		}
		updated += n
	}

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

	fmt.Println()
	green.Printf("  ✓ Upgrade complete! Updated %d files.\n", updated)
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

// upgradeAdminFiles regenerates admin panel generic files while preserving
// user-created resource definitions and pages.
func upgradeAdminFiles(root string, opts Options, uOpts UpgradeOptions) (int, error) {
	adminRoot := filepath.Join(root, "apps", "admin")
	green := color.New(color.FgHiGreen)

	// These are all framework-generated files that are safe to overwrite.
	// User-created files (resources/*.ts, app/(dashboard)/resources/*/page.tsx) are NOT touched.
	files := map[string]string{
		// Config files
		filepath.Join(adminRoot, "package.json"):       adminPackageJSON(opts),
		filepath.Join(adminRoot, "tailwind.config.ts"): adminTailwindConfig(),
		filepath.Join(adminRoot, "postcss.config.js"):  adminPostCSSConfig(),
		filepath.Join(adminRoot, "next.config.ts"):     adminNextConfig(),
		filepath.Join(adminRoot, "tsconfig.json"):      adminTSConfig(),

		// App core
		filepath.Join(adminRoot, "app", "globals.css"):                                   adminGlobalCSS(),
		filepath.Join(adminRoot, "app", "layout.tsx"):                                    adminRootLayout(opts),
		filepath.Join(adminRoot, "app", "page.tsx"):                                      adminRedirectPage(),
		filepath.Join(adminRoot, "app", "(auth)", "login", "page.tsx"):                   adminLoginPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(auth)", "sign-up", "page.tsx"):                 adminSignUpPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(auth)", "forgot-password", "page.tsx"):         adminForgotPasswordPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(auth)", "reset-password", "page.tsx"):          adminThemedResetPasswordPage(),
		filepath.Join(adminRoot, "app", "(auth)", "callback", "page.tsx"):                adminAuthCallbackPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "layout.tsx"):                     adminDashboardLayout(),
		filepath.Join(adminRoot, "app", "(dashboard)", "dashboard", "page.tsx"):          adminDashboardPageForStyle(opts.Style),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", "users", "page.tsx"): adminUsersPage(),

		// System pages
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "jobs", "page.tsx"):          adminJobsPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "files", "page.tsx"):         adminFilesPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "cron", "page.tsx"):          adminCronPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "mail", "page.tsx"):          adminMailPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "security", "page.tsx"):      adminSecurityPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "observability", "page.tsx"): adminObservabilityPage(),

		// Lib
		filepath.Join(adminRoot, "lib", "api-client.ts"):   adminAPIClient(),
		filepath.Join(adminRoot, "lib", "query-client.ts"): adminQueryClient(),
		filepath.Join(adminRoot, "lib", "utils.ts"):        adminUtils(),
		filepath.Join(adminRoot, "lib", "resource.ts"):     adminResourceTypes(),
		filepath.Join(adminRoot, "lib", "icons.ts"):        adminIconMap(),

		// Hooks
		filepath.Join(adminRoot, "hooks", "use-auth.ts"):                adminUseAuth(),
		filepath.Join(adminRoot, "hooks", "use-resource.ts"):            adminUseResource(),
		filepath.Join(adminRoot, "hooks", "use-resource-controller.ts"): adminUseResourceController(),
		filepath.Join(adminRoot, "hooks", "use-system.ts"):              adminUseSystem(),

		// Layout components
		filepath.Join(adminRoot, "components", "layout", "admin-layout.tsx"):   adminLayoutComponent(),
		filepath.Join(adminRoot, "components", "layout", "sidebar.tsx"):        adminSidebar(),
		filepath.Join(adminRoot, "components", "layout", "navbar.tsx"):         adminNavbar(),
		filepath.Join(adminRoot, "components", "shared", "theme-provider.tsx"): adminThemeProvider(),
		filepath.Join(adminRoot, "components", "shared", "providers.tsx"):      adminProviders(),

		// Table components
		filepath.Join(adminRoot, "components", "tables", "data-table.tsx"):        adminDataTable(),
		filepath.Join(adminRoot, "components", "tables", "column-header.tsx"):     adminColumnHeader(),
		filepath.Join(adminRoot, "components", "tables", "cell-renderers.tsx"):    adminCellRenderers(),
		filepath.Join(adminRoot, "components", "tables", "table-filters.tsx"):     adminTableFilters(),
		filepath.Join(adminRoot, "components", "tables", "table-toolbar.tsx"):     adminTableToolbar(),
		filepath.Join(adminRoot, "components", "tables", "table-pagination.tsx"):  adminTablePagination(),
		filepath.Join(adminRoot, "components", "tables", "table-skeleton.tsx"):    adminTableSkeleton(),
		filepath.Join(adminRoot, "components", "tables", "table-empty-state.tsx"): adminTableEmptyState(),
		filepath.Join(adminRoot, "lib", "formatters.ts"):                          adminFormatters(),

		// Form components
		filepath.Join(adminRoot, "components", "forms", "form-builder.tsx"):             adminFormBuilder(),
		filepath.Join(adminRoot, "components", "forms", "form-modal.tsx"):               adminFormModal(),
		filepath.Join(adminRoot, "components", "forms", "form-page.tsx"):                adminFormPage(),
		filepath.Join(adminRoot, "components", "forms", "form-sheet.tsx"):               adminFormSheet(),
		filepath.Join(adminRoot, "components", "forms", "form-modal-steps.tsx"):         adminFormModalSteps(),
		filepath.Join(adminRoot, "components", "forms", "form-page-steps.tsx"):          adminFormPageSteps(),
		filepath.Join(adminRoot, "components", "forms", "update-groups.tsx"):            adminUpdateGroups(),
		filepath.Join(adminRoot, "components", "forms", "fields", "text-field.tsx"):     adminTextField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "textarea-field.tsx"): adminTextareaField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "number-field.tsx"):   adminNumberField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "select-field.tsx"):   adminSelectField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "date-field.tsx"):     adminDateField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "toggle-field.tsx"):   adminToggleField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "checkbox-field.tsx"): adminCheckboxField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "radio-field.tsx"):    adminRadioField(),

		// Widget components
		filepath.Join(adminRoot, "components", "widgets", "stats-card.tsx"):      adminStatsCard(),
		filepath.Join(adminRoot, "components", "widgets", "chart-widget.tsx"):    adminChartWidget(),
		filepath.Join(adminRoot, "components", "widgets", "activity-widget.tsx"): adminActivityWidget(),
		filepath.Join(adminRoot, "components", "widgets", "widget-grid.tsx"):     adminWidgetGrid(),

		// Resource components
		filepath.Join(adminRoot, "components", "resource", "resource-page.tsx"):        adminResourcePage(),
		filepath.Join(adminRoot, "components", "resource", "view-modal.tsx"):           adminViewModal(),
		filepath.Join(adminRoot, "components", "resource", "resource-detail-page.tsx"): adminResourceDetailPage(),
		filepath.Join(adminRoot, "components", "ui", "confirm-modal.tsx"):              adminConfirmModal(),

		// Dropzone
		filepath.Join(adminRoot, "components", "ui", "dropzone.tsx"): adminDropzone(),

		// Resource definitions (only the built-in users one)
		filepath.Join(adminRoot, "resources", "users.ts"): adminUsersResource(),
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
		filepath.Join(adminRoot, "resources", "users.custom.tsx"),
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
		fmt.Printf("  ✓ Removed %d unwired i18n file(s) — run `grit add i18n` to set it up properly\n", removed)
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
