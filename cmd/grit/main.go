package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v2/internal/generate"
	"github.com/MUKE-coder/grit/v2/internal/project"
	"github.com/MUKE-coder/grit/v2/internal/scaffold"
)

var version = "2.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "grit",
		Short: "Grit — Go + React. Built with Grit.",
		Long:  "Grit is a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript).",
	}

	rootCmd.AddCommand(newCmd())
	rootCmd.AddCommand(newDesktopCmd())
	rootCmd.AddCommand(generateCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(addCmd())
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(compileCmd())
	rootCmd.AddCommand(studioCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(seedCmd())
	rootCmd.AddCommand(upgradeCmd())
	rootCmd.AddCommand(updateCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of Grit CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("grit version %s\n", version)
		},
	}
}

func newCmd() *cobra.Command {
	var apiOnly, includeExpo, mobileOnly, full bool
	var style string

	cmd := &cobra.Command{
		Use:   "new <project-name>",
		Short: "Create a new Grit project",
		Long:  "Scaffold a new Grit monorepo with Go API, Next.js frontend, admin panel, and shared packages.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			if err := scaffold.ValidateProjectName(projectName); err != nil {
				return err
			}

			// Validate mutual exclusivity
			count := 0
			if apiOnly {
				count++
			}
			if includeExpo {
				count++
			}
			if mobileOnly {
				count++
			}
			if full {
				count++
			}
			if count > 1 {
				return fmt.Errorf("only one mode flag can be used at a time (--api, --expo, --mobile, --full)")
			}

			opts := scaffold.Options{
				ProjectName: projectName,
				APIOnly:     apiOnly,
				IncludeExpo: includeExpo,
				MobileOnly:  mobileOnly,
				Full:        full,
				Style:       style,
			}

			if err := opts.ValidateStyle(); err != nil {
				return err
			}

			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Creating new Grit project: %s\n\n", projectName)

			if err := scaffold.Run(opts); err != nil {
				color.Red("\n  Error: %v\n", err)
				return err
			}

			printSuccess(projectName, opts)
			return nil
		},
	}

	cmd.Flags().BoolVar(&apiOnly, "api", false, "Scaffold only the Go API (no frontend apps)")
	cmd.Flags().BoolVar(&includeExpo, "expo", false, "Include Expo mobile app (api + web + admin + shared + expo)")
	cmd.Flags().BoolVar(&mobileOnly, "mobile", false, "Scaffold API + Expo mobile app only")
	cmd.Flags().BoolVar(&full, "full", false, "Scaffold everything including docs site")
	cmd.Flags().StringVar(&style, "style", "default", "Admin panel style variant (default, modern, minimal, glass)")

	return cmd
}

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code for your Grit project",
		Aliases: []string{"g"},
	}

	cmd.AddCommand(generateResourceCmd())

	return cmd
}

func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove components from your Grit project",
		Aliases: []string{"rm"},
	}

	cmd.AddCommand(removeResourceCmd())

	return cmd
}

func removeResourceCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "resource <Name>",
		Short: "Remove a previously generated resource",
		Long:  "Delete generated files and reverse all marker-based injections for a resource.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			if !force {
				fmt.Printf("\n  ⚠ This will remove all files and injections for resource %q.\n", args[0])
				if !generate.ConfirmRemoval() {
					fmt.Println("\n  Cancelled.")
					return nil
				}
			}

			return generate.RemoveResource(args[0])
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add components to your Grit project",
	}

	cmd.AddCommand(addRoleCmd())

	return cmd
}

func addRoleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "role <ROLE_NAME>",
		Short: "Add a new role to the project",
		Long:  "Adds a new role constant to Go models, TypeScript types, Zod schemas, constants, and admin resource definitions.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Adding role: %s\n\n", strings.ToUpper(args[0]))

			return scaffold.AddRole(args[0])
		},
	}
}

func generateResourceCmd() *cobra.Command {
	var fromFile string
	var interactive bool
	var fields string
	var roles string

	cmd := &cobra.Command{
		Use:   "resource <Name>",
		Short: "Generate a full-stack CRUD resource",
		Long:  "Generate Go model, handler, service, Zod schemas, TypeScript types, React Query hooks, and admin page.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			printLogo()

			var def *generate.ResourceDefinition
			var err error

			switch {
			case fromFile != "":
				def, err = generate.LoadFromYAML(fromFile)
				if err != nil {
					return err
				}
				// Override name if provided
				if name != "" {
					def.Name = name
				}
			case fields != "":
				def, err = generate.ParseInlineFields(name, fields)
				if err != nil {
					return err
				}
			case interactive:
				def, err = generate.PromptInteractive(name)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("specify fields with --fields, --from, or use -i for interactive mode\n\nExamples:\n  grit generate resource Post --fields \"title:string,content:text,published:bool\"\n  grit generate resource Post --fields \"title:string,slug:string:unique,views:int\"\n  grit generate resource Post --from post.yaml\n  grit generate resource Post -i")
			}

			// Detect project type and dispatch
			info, _ := project.DetectProject()
			if info != nil && info.Type == project.ProjectDesktop {
				return generate.RunDesktop(def)
			}

			gen, err := generate.NewGenerator(def)
			if err != nil {
				return err
			}

			// Parse --roles flag
			if roles != "" {
				parts := strings.Split(roles, ",")
				for _, r := range parts {
					r = strings.TrimSpace(strings.ToUpper(r))
					if r != "" {
						gen.Roles = append(gen.Roles, r)
					}
				}
			}

			return gen.Run()
		},
	}

	cmd.Flags().StringVar(&fromFile, "from", "", "YAML file defining the resource fields")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively define fields")
	cmd.Flags().StringVar(&fields, "fields", "", "Inline field definitions (e.g., \"title:string,content:text,published:bool\")")
	cmd.Flags().StringVar(&roles, "roles", "", "Restrict routes to specific roles (comma-separated, e.g., \"ADMIN,EDITOR\")")

	return cmd
}

func migrateCmd() *cobra.Command {
	var fresh bool

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  "Connect to the database and run GORM AutoMigrate for all models. Use --fresh to drop all tables first.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}

			goArgs := []string{"run", "cmd/migrate/main.go"}
			if fresh {
				goArgs = append(goArgs, "--fresh")
			}

			c := exec.Command("go", goArgs...)
			c.Dir = apiDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			purple := color.New(color.FgHiMagenta, color.Bold)
			if fresh {
				purple.Println("\n  Running fresh migration (drop + re-migrate)...")
			} else {
				purple.Println("\n  Running database migrations...")
			}

			return c.Run()
		},
	}

	cmd.Flags().BoolVar(&fresh, "fresh", false, "Drop all tables before migrating (migrate:fresh)")

	return cmd
}

func seedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Run database seeders",
		Long:  "Populate the database with initial data (admin user, demo users).",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}

			c := exec.Command("go", "run", "cmd/seed/main.go")
			c.Dir = apiDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Println("\n  Seeding database...")

			return c.Run()
		},
	}
}

// findAPIDir locates the apps/api directory from the project root.
func findAPIDir() (string, error) {
	root, err := scaffold.FindProjectRoot()
	if err != nil {
		return "", err
	}
	apiDir := filepath.Join(root, "apps", "api")
	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		return "", fmt.Errorf("apps/api directory not found in %s", root)
	}
	return apiDir, nil
}

func upgradeCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade an existing Grit project to the latest scaffold templates",
		Long:  "Regenerates framework components (admin panel, web app, configs) while preserving your resource definitions and API code.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Upgrading project to Grit v%s\n\n", version)

			return scaffold.Upgrade(scaffold.UpgradeOptions{
				Force: force,
			})
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite all files without prompting")

	return cmd
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the Grit CLI to the latest version",
		Long:  "Removes the current Grit binary and installs the latest version from GitHub using go install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			green := color.New(color.FgHiGreen, color.Bold)
			spinner := color.New(color.FgHiBlack)

			purple.Printf("\n  Updating Grit CLI (current: v%s)...\n\n", version)

			// Find the current binary path
			binPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("finding current binary: %w", err)
			}
			binPath, err = filepath.EvalSymlinks(binPath)
			if err != nil {
				return fmt.Errorf("resolving binary path: %w", err)
			}

			// On Windows, a running binary can't be deleted but can be renamed
			if runtime.GOOS == "windows" {
				oldPath := binPath + ".old"
				// Clean up any previous .old file
				os.Remove(oldPath)
				spinner.Printf("  → Renaming old binary: %s\n", binPath)
				if err := os.Rename(binPath, oldPath); err != nil {
					return fmt.Errorf("renaming old binary: %w", err)
				}
			} else {
				spinner.Printf("  → Removing old binary: %s\n", binPath)
				if err := os.Remove(binPath); err != nil {
					return fmt.Errorf("removing old binary: %w", err)
				}
			}

			spinner.Println("  → Installing latest version...")
			c := exec.Command("go", "install", "github.com/MUKE-coder/grit/v2/cmd/grit@latest")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("installing latest version: %w", err)
			}

			// Clean up renamed binary on Windows
			if runtime.GOOS == "windows" {
				os.Remove(binPath + ".old")
			}

			fmt.Println()
			green.Println("  ✓ Grit CLI updated successfully!")
			spinner.Println("  Run 'grit version' to verify the new version.")
			fmt.Println()

			return nil
		},
	}
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start development servers",
		Long:  "Start the Go API server or frontend client apps for local development. In a desktop project, runs wails dev.",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := project.DetectProject()
			if err != nil {
				// Not inside a project — show subcommand help
				return cmd.Help()
			}
			if info.Type == project.ProjectDesktop {
				printLogo()
				purple := color.New(color.FgHiMagenta, color.Bold)
				purple.Println("\n  Starting Wails desktop app...")

				c := exec.Command("wails", "dev")
				c.Dir = info.Root
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Stdin = os.Stdin
				return c.Run()
			}
			// Web project — show subcommand help
			return cmd.Help()
		},
	}

	cmd.AddCommand(startClientCmd())
	cmd.AddCommand(startServerCmd())

	return cmd
}

func compileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "compile",
		Short: "Build desktop application executable",
		Long:  "Compile the Wails desktop app into a distributable binary. Only available in desktop projects.",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := project.DetectProject()
			if err != nil {
				return fmt.Errorf("not inside a Grit project: %w", err)
			}
			if info.Type != project.ProjectDesktop {
				return fmt.Errorf("grit compile is only available in desktop projects (wails.json not found)")
			}

			printLogo()
			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Println("\n  Building desktop executable...")

			c := exec.Command("wails", "build")
			c.Dir = info.Root
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			return c.Run()
		},
	}
}

func studioCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "studio",
		Short: "Open GORM Studio database browser",
		Long:  "Launch the database studio. For web projects, opens the browser. For desktop, starts a standalone studio server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := project.DetectProject()
			if err != nil {
				return fmt.Errorf("not inside a Grit project: %w", err)
			}

			printLogo()

			if info.Type == project.ProjectDesktop {
				purple := color.New(color.FgHiMagenta, color.Bold)
				purple.Println("\n  Starting GORM Studio...")

				c := exec.Command("go", "run", "cmd/studio/main.go")
				c.Dir = info.Root
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Stdin = os.Stdin
				return c.Run()
			}

			// Web project — open browser
			gray := color.New(color.FgHiBlack)
			gray.Println("\n  GORM Studio is available at http://localhost:8080/studio")
			gray.Println("  Make sure your API server is running (grit start server)")

			openURL("http://localhost:8080/studio")
			return nil
		},
	}
}

func startClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "client",
		Short: "Start frontend apps (web, admin, expo)",
		Long:  "Runs 'pnpm dev' from the project root to start all frontend apps via Turborepo.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Println("\n  Starting client apps...")

			c := exec.Command("pnpm", "dev")
			c.Dir = root
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}
}

func startServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the Go API server",
		Long:  "Runs 'go run cmd/server/main.go' from the apps/api directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Println("\n  Starting API server...")

			c := exec.Command("go", "run", "cmd/server/main.go")
			c.Dir = apiDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}
}

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync Go types to TypeScript types and Zod schemas",
		Long:  "Parse Go model files and regenerate TypeScript types and Zod schemas in packages/shared.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()
			return generate.Sync()
		},
	}
}

func openURL(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func printLogo() {
	purple := color.New(color.FgHiMagenta, color.Bold)
	purple.Println(`
   ██████╗ ██████╗ ██╗████████╗
  ██╔════╝ ██╔══██╗██║╚══██╔══╝
  ██║  ███╗██████╔╝██║   ██║
  ██║   ██║██╔══██╗██║   ██║
  ╚██████╔╝██║  ██║██║   ██║
   ╚═════╝ ╚═╝  ╚═╝╚═╝   ╚═╝`)
	gray := color.New(color.FgHiBlack)
	gray.Printf("  Go + React. Built with Grit. v%s\n", version)
}

func printSuccess(name string, opts scaffold.Options) {
	green := color.New(color.FgHiGreen, color.Bold)
	white := color.New(color.FgWhite)
	cyan := color.New(color.FgHiCyan)
	gray := color.New(color.FgHiBlack)

	fmt.Println()
	green.Println("  ✓ Project created successfully!")
	fmt.Println()

	white.Println("  Next steps:")
	fmt.Println()
	cyan.Printf("    cd %s\n", name)
	cyan.Println("    docker compose up -d")

	if opts.APIOnly {
		cyan.Println("    cd apps/api && go run cmd/server/main.go")
	} else {
		cyan.Println("    pnpm install")
		cyan.Println("    pnpm dev")
	}

	if opts.ShouldIncludeExpo() {
		cyan.Println("    cd apps/expo && npx expo start")
	}

	fmt.Println()
	gray.Println("  ─────────────────────────────────────")
	gray.Printf("  API:         http://localhost:8080\n")
	gray.Printf("  API Docs:    http://localhost:8080/docs\n")
	gray.Printf("  GORM Studio: http://localhost:8080/studio\n")
	gray.Printf("  Sentinel:    http://localhost:8080/sentinel/ui\n")

	if opts.ShouldIncludeWeb() {
		gray.Printf("  Web App:     http://localhost:3000\n")
	}
	if opts.ShouldIncludeAdmin() {
		gray.Printf("  Admin:       http://localhost:3001\n")
	}
	if opts.ShouldIncludeExpo() {
		gray.Printf("  Expo:        exp://localhost:8081\n")
	}
	if opts.ShouldIncludeDocs() {
		gray.Printf("  Docs:        http://localhost:3002\n")
	}

	gray.Printf("  PostgreSQL:  localhost:5432\n")
	gray.Printf("  Redis:       localhost:6379\n")
	gray.Printf("  MinIO:       http://localhost:9001\n")
	gray.Printf("  Mailhog:     http://localhost:8025\n")
	gray.Println("  ─────────────────────────────────────")
	fmt.Println()
}

func newDesktopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-desktop <project-name>",
		Short: "Create a new Wails desktop application",
		Long:  "Scaffold a standalone desktop app with Go backend (Wails), React frontend, SQLite/PostgreSQL, and local auth.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			if err := scaffold.ValidateProjectName(projectName); err != nil {
				return err
			}

			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Creating new Grit desktop app: %s\n\n", projectName)

			opts := scaffold.DesktopOptions{
				ProjectName: projectName,
			}

			if err := scaffold.RunDesktop(opts); err != nil {
				color.Red("\n  Error: %v\n", err)
				return err
			}

			printDesktopSuccess(projectName)
			return nil
		},
	}

	return cmd
}

func printDesktopSuccess(name string) {
	green := color.New(color.FgHiGreen, color.Bold)
	white := color.New(color.FgWhite)
	cyan := color.New(color.FgHiCyan)
	gray := color.New(color.FgHiBlack)

	fmt.Println()
	green.Println("  ✓ Desktop app created successfully!")
	fmt.Println()

	white.Println("  Next steps:")
	fmt.Println()
	cyan.Printf("    cd %s\n", name)
	cyan.Println("    wails dev")
	fmt.Println()

	gray.Println("  ─────────────────────────────────────")
	gray.Printf("  Desktop App: http://localhost:34115\n")
	gray.Printf("  Database:    SQLite (%s.db)\n", name)
	gray.Println("  ─────────────────────────────────────")
	fmt.Println()

	gray.Println("  To build for production:")
	cyan.Println("    wails build")
	fmt.Println()
}
