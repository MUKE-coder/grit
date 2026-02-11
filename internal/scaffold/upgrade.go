package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	opts := Options{ProjectName: projectName}

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
		filepath.Join(root, "pnpm-workspace.yaml"): pnpmWorkspace(),
	}
	n, err := writeUpgradeFiles(rootFiles, uOpts.Force)
	if err != nil {
		return fmt.Errorf("updating root files: %w", err)
	}
	updated += n

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

		// API documentation (Scalar)
		spinner.Printf("  → Updating API documentation (Scalar)...\n")
		if err := writeScalarFiles(root, opts); err != nil {
			return fmt.Errorf("updating Scalar docs files: %w", err)
		}
		green.Printf("  ✓ API documentation updated\n")
		updated += 2
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
		filepath.Join(adminRoot, "package.json"):      adminPackageJSON(opts),
		filepath.Join(adminRoot, "tailwind.config.ts"): adminTailwindConfig(),
		filepath.Join(adminRoot, "postcss.config.mjs"): adminPostCSSConfig(),
		filepath.Join(adminRoot, "next.config.ts"):     adminNextConfig(),
		filepath.Join(adminRoot, "tsconfig.json"):      adminTSConfig(),

		// App core
		filepath.Join(adminRoot, "app", "globals.css"):                              adminGlobalCSS(),
		filepath.Join(adminRoot, "app", "layout.tsx"):                               adminRootLayout(opts),
		filepath.Join(adminRoot, "app", "page.tsx"):                                 adminRedirectPage(),
		filepath.Join(adminRoot, "app", "(auth)", "login", "page.tsx"):              adminLoginPage(),
		filepath.Join(adminRoot, "app", "(auth)", "sign-up", "page.tsx"):            adminSignUpPage(),
		filepath.Join(adminRoot, "app", "(auth)", "forgot-password", "page.tsx"):    adminForgotPasswordPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "layout.tsx"):                adminDashboardLayout(),
		filepath.Join(adminRoot, "app", "(dashboard)", "dashboard", "page.tsx"):     adminDashboardPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", "users", "page.tsx"): adminUsersPage(),

		// System pages
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "jobs", "page.tsx"):  adminJobsPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "files", "page.tsx"): adminFilesPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "cron", "page.tsx"):  adminCronPage(),
		filepath.Join(adminRoot, "app", "(dashboard)", "system", "mail", "page.tsx"):  adminMailPage(),

		// Lib
		filepath.Join(adminRoot, "lib", "api-client.ts"):  adminAPIClient(),
		filepath.Join(adminRoot, "lib", "query-client.ts"): adminQueryClient(),
		filepath.Join(adminRoot, "lib", "utils.ts"):        adminUtils(),
		filepath.Join(adminRoot, "lib", "resource.ts"):     adminResourceTypes(),
		filepath.Join(adminRoot, "lib", "icons.ts"):        adminIconMap(),

		// Hooks
		filepath.Join(adminRoot, "hooks", "use-auth.ts"):    adminUseAuth(),
		filepath.Join(adminRoot, "hooks", "use-resource.ts"): adminUseResource(),
		filepath.Join(adminRoot, "hooks", "use-system.ts"):  adminUseSystem(),

		// Layout components
		filepath.Join(adminRoot, "components", "layout", "AdminLayout.tsx"):   adminLayoutComponent(),
		filepath.Join(adminRoot, "components", "layout", "Sidebar.tsx"):       adminSidebar(),
		filepath.Join(adminRoot, "components", "layout", "Navbar.tsx"):        adminNavbar(),
		filepath.Join(adminRoot, "components", "layout", "ThemeProvider.tsx"): adminThemeProvider(),
		filepath.Join(adminRoot, "components", "shared", "Providers.tsx"):     adminProviders(),

		// Table components
		filepath.Join(adminRoot, "components", "tables", "DataTable.tsx"):       adminDataTable(),
		filepath.Join(adminRoot, "components", "tables", "ColumnHeader.tsx"):    adminColumnHeader(),
		filepath.Join(adminRoot, "components", "tables", "CellRenderers.tsx"):   adminCellRenderers(),
		filepath.Join(adminRoot, "components", "tables", "TableFilters.tsx"):    adminTableFilters(),
		filepath.Join(adminRoot, "components", "tables", "TableToolbar.tsx"):    adminTableToolbar(),
		filepath.Join(adminRoot, "components", "tables", "TablePagination.tsx"): adminTablePagination(),
		filepath.Join(adminRoot, "components", "tables", "Skeleton.tsx"):        adminTableSkeleton(),
		filepath.Join(adminRoot, "components", "tables", "EmptyState.tsx"):      adminTableEmptyState(),
		filepath.Join(adminRoot, "components", "tables", "Formatters.ts"):       adminFormatters(),

		// Form components
		filepath.Join(adminRoot, "components", "forms", "FormBuilder.tsx"):         adminFormBuilder(),
		filepath.Join(adminRoot, "components", "forms", "FormModal.tsx"):           adminFormModal(),
		filepath.Join(adminRoot, "components", "forms", "fields", "TextField.tsx"):     adminTextField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "TextAreaField.tsx"): adminTextareaField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "NumberField.tsx"):   adminNumberField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "SelectField.tsx"):   adminSelectField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "DateField.tsx"):     adminDateField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "ToggleField.tsx"):   adminToggleField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "CheckboxField.tsx"): adminCheckboxField(),
		filepath.Join(adminRoot, "components", "forms", "fields", "RadioField.tsx"):    adminRadioField(),

		// Widget components
		filepath.Join(adminRoot, "components", "widgets", "StatsCard.tsx"):    adminStatsCard(),
		filepath.Join(adminRoot, "components", "widgets", "ChartWidget.tsx"):  adminChartWidget(),
		filepath.Join(adminRoot, "components", "widgets", "ActivityWidget.tsx"): adminActivityWidget(),
		filepath.Join(adminRoot, "components", "widgets", "WidgetGrid.tsx"):   adminWidgetGrid(),

		// Resource components
		filepath.Join(adminRoot, "components", "resource", "ResourcePage.tsx"): adminResourcePage(),
		filepath.Join(adminRoot, "components", "resource", "view-modal.tsx"):   adminViewModal(),
		filepath.Join(adminRoot, "components", "ui", "confirm-modal.tsx"):      adminConfirmModal(),

		// Dropzone
		filepath.Join(adminRoot, "components", "ui", "dropzone.tsx"): adminDropzone(),

		// Resource definitions (only the built-in users one)
		filepath.Join(adminRoot, "resources", "users.ts"):    adminUsersResource(),
		filepath.Join(adminRoot, "resources", "registry.ts"): adminResourceRegistry(),
	}

	n, err := writeUpgradeFiles(files, uOpts.Force)
	if err != nil {
		return 0, err
	}

	green.Printf("  ✓ Admin panel updated (%d files)\n", n)
	return n, nil
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

// FindProjectRoot walks up from the current directory looking for a Grit project.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}

	for {
		// A Grit project has turbo.json + apps/api directory
		if fileExists(filepath.Join(dir, "turbo.json")) && dirExists(filepath.Join(dir, "apps", "api")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("not inside a Grit project (no turbo.json + apps/api found)\n\nRun this command from your Grit project root or any subdirectory")
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
