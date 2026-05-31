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

	"github.com/MUKE-coder/grit/v3/internal/deploy"
	"github.com/MUKE-coder/grit/v3/internal/generate"
	"github.com/MUKE-coder/grit/v3/internal/maintenance"
	"github.com/MUKE-coder/grit/v3/internal/project"
	"github.com/MUKE-coder/grit/v3/internal/prompt"
	"github.com/MUKE-coder/grit/v3/internal/routeparser"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/selfupdate"
)

var version = "3.25.2"

func main() {
	rootCmd := &cobra.Command{
		Use:   "grit",
		Short: "Grit — Go + React. Built with Grit.",
		Long:  "Grit is a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript).",
	}

	rootCmd.AddCommand(newCmd())
	rootCmd.AddCommand(newDesktopCmd())
	rootCmd.AddCommand(initCmd())
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
	rootCmd.AddCommand(routesCmd())
	rootCmd.AddCommand(downCmd())
	rootCmd.AddCommand(upCmd())
	rootCmd.AddCommand(deployCmd())

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

func initCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Write CLAUDE.md / AGENTS.md framework convention docs to the project root",
		Long: `Write the framework's hard-rules documentation as CLAUDE.md and AGENTS.md
in the current directory. Both files have the same content — different LLM
tooling looks for different filenames.

Skips files that already exist unless --force is passed. Re-run with --force
after a major framework upgrade to refresh the rules.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			written, err := scaffold.WriteAgentsDoc(cwd, force)
			if err != nil {
				return err
			}
			if len(written) == 0 {
				fmt.Println("\n  CLAUDE.md and AGENTS.md already exist — re-run with --force to overwrite.")
				return nil
			}
			fmt.Println()
			for _, name := range written {
				fmt.Printf("  ✓ %s\n", name)
			}
			fmt.Printf("\n  ✅ Framework conventions written. Commit these so contributors\n")
			fmt.Printf("     (and AI assistants) get the rules right on first PR.\n\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	return cmd
}

func newCmd() *cobra.Command {
	// New architecture/frontend flags
	var archFlag, frontendFlag, style string
	var inPlace, force, includeDesktop bool

	// Legacy flags (backward compatibility)
	var apiOnly, includeExpo, mobileOnly, full bool

	cmd := &cobra.Command{
		Use:   "new <project-name|.>",
		Short: "Create a new Grit project",
		Long:  "Scaffold a new Grit project. Interactive by default — select architecture and frontend.\nUse flags to skip prompts: grit new my-app --single --vite\nUse `grit new .` to scaffold in the current directory.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectArg := strings.TrimSpace(args[0])
			projectName := projectArg

			if filepath.Clean(projectArg) == "." {
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("getting current directory: %w", err)
				}
				projectName = filepath.Base(cwd)
				inPlace = true
			}

			if err := scaffold.ValidateProjectName(projectName); err != nil {
				if filepath.Clean(projectArg) == "." {
					return fmt.Errorf("invalid current directory name %q for project generation: %w", projectName, err)
				}
				return err
			}

			opts := scaffold.Options{
				ProjectName:    projectName,
				Style:          style,
				InPlace:        inPlace,
				Force:          force,
				IncludeDesktop: includeDesktop,
				// Legacy flags
				APIOnly:     apiOnly,
				IncludeExpo: includeExpo,
				MobileOnly:  mobileOnly,
				Full:        full,
			}

			if !opts.InPlace {
				cwd, err := os.Getwd()
				if err == nil && filepath.Base(cwd) == projectName {
					opts.InPlace = true
				}
			}

			// Map architecture shorthand flags
			switch archFlag {
			case "single":
				opts.Architecture = scaffold.ArchSingle
			case "double":
				opts.Architecture = scaffold.ArchDouble
			case "triple":
				opts.Architecture = scaffold.ArchTriple
			case "api":
				opts.Architecture = scaffold.ArchAPI
			case "mobile":
				opts.Architecture = scaffold.ArchMobile
			case "":
				// Will be set by legacy flags or interactive prompt
			default:
				return fmt.Errorf("invalid architecture %q: must be single, double, triple, api, or mobile", archFlag)
			}

			// Map frontend shorthand flags
			switch frontendFlag {
			case "next":
				opts.Frontend = scaffold.FrontendNext
			case "vite", "tanstack":
				opts.Frontend = scaffold.FrontendTanStack
			case "":
				// Will be set by interactive prompt or default
			default:
				return fmt.Errorf("invalid frontend %q: must be next, vite, or tanstack", frontendFlag)
			}

			// Show interactive selector only when the user did not provide any
			// architecture/frontend shortcuts or explicit long-form flags.
			// We check the resolved flag values directly (set by PreRunE) rather
			// than cmd.Flags().Changed(), which is more reliable across flag sources.
			anyFlagSet := archFlag != "" || frontendFlag != "" ||
				apiOnly || mobileOnly || full || includeExpo

			if !anyFlagSet {
				printLogo()
				// Keep empty values so the prompt can collect architecture/frontend.
				opts.Architecture = ""
				opts.Frontend = ""
				if err := prompt.RunNewProjectPrompt(&opts); err != nil {
					return err
				}
			}

			// Final normalization (sets defaults for anything still empty)
			opts.Normalize()

			if err := opts.ValidateStyle(); err != nil {
				return err
			}

			// --desktop is incompatible with --single (single apps already bundle their own SPA).
			if opts.IncludeDesktop && opts.Architecture == scaffold.ArchSingle {
				return fmt.Errorf("--desktop is not supported with --single architecture (single apps already bundle their own frontend)")
			}

			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Creating new Grit project: %s\n", projectName)

			gray := color.New(color.FgHiBlack)
			gray.Printf("  Architecture: %s | Frontend: %s\n\n", opts.Architecture, opts.Frontend)

			if err := scaffold.Run(opts); err != nil {
				color.Red("\n  Error: %v\n", err)
				return err
			}

			printSuccess(projectName, opts)
			return nil
		},
	}

	// New flags
	cmd.Flags().StringVar(&archFlag, "arch", "", "Architecture: single, double, triple, api, mobile")
	cmd.Flags().StringVar(&frontendFlag, "frontend", "", "Frontend framework: next, vite (tanstack)")
	cmd.Flags().StringVar(&style, "style", "", "Admin panel style variant (default, modern, minimal, glass)")

	// Shorthand architecture flags
	cmd.Flags().BoolVar(&apiOnly, "api", false, "Shorthand for --arch=api")
	cmd.Flags().BoolVar(&mobileOnly, "mobile", false, "Shorthand for --arch=mobile")
	cmd.Flags().BoolVar(&full, "full", false, "Shorthand for --arch=triple with docs")
	cmd.Flags().BoolVar(&includeExpo, "expo", false, "Include Expo mobile app")
	cmd.Flags().BoolVar(&includeDesktop, "desktop", false, "Include a Wails desktop app that shares the monorepo API")

	// Shorthand frontend flags
	cmd.Flags().Bool("vite", false, "Shorthand for --frontend=vite (TanStack Router)")
	cmd.Flags().Bool("next", false, "Shorthand for --frontend=next (Next.js)")

	// Handle shorthand frontend flags
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if v, _ := cmd.Flags().GetBool("vite"); v {
			frontendFlag = "vite"
		}
		if v, _ := cmd.Flags().GetBool("next"); v {
			frontendFlag = "next"
		}
		// Shorthand single/double/triple
		if v, _ := cmd.Flags().GetBool("single"); v {
			archFlag = "single"
		}
		if v, _ := cmd.Flags().GetBool("double"); v {
			archFlag = "double"
		}
		if v, _ := cmd.Flags().GetBool("triple"); v {
			archFlag = "triple"
		}
		return nil
	}

	// Shorthand architecture flags
	cmd.Flags().Bool("single", false, "Shorthand for --arch=single")
	cmd.Flags().Bool("double", false, "Shorthand for --arch=double")
	cmd.Flags().Bool("triple", false, "Shorthand for --arch=triple")
	cmd.Flags().BoolVar(&inPlace, "here", false, "Scaffold into the current directory instead of creating a new folder")
	cmd.Flags().BoolVar(&force, "force", false, "Allow scaffolding into a non-empty directory (use with --here)")

	return cmd
}

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate code for your Grit project",
		Aliases: []string{"g"},
	}

	cmd.AddCommand(generateResourceCmd())
	cmd.AddCommand(generateSequenceCmd())

	return cmd
}

func generateSequenceCmd() *cobra.Command {
	var prefix string
	var reset string
	var width int

	cmd := &cobra.Command{
		Use:   "sequence <Name>",
		Short: "Generate a sequential numbering helper (e.g. INV-202605-0001)",
		Long: `Generate atomic sequential numbers for a resource.

The first invocation in a project creates internal/sequence/ — a generic
counter package backed by a database row. Every invocation also writes a
typed convenience wrapper at internal/services/<name>_sequence.go so
handlers call services.Next<Name>Number(db, t) without knowing the
prefix/reset/width.

Examples:
  grit generate sequence Invoice
  grit generate sequence Invoice --prefix INV --reset monthly --width 4
  grit generate sequence Order --prefix ORD --reset yearly --width 6
  grit generate sequence Receipt --reset never`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()
			return generate.GenerateSequence(generate.SequenceOptions{
				Name:   args[0],
				Prefix: prefix,
				Reset:  reset,
				Width:  width,
			})
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Alphabetic prefix (default: first 3 chars of name uppercased)")
	cmd.Flags().StringVar(&reset, "reset", "monthly", "When the counter resets: monthly, yearly, never")
	cmd.Flags().IntVar(&width, "width", 4, "Zero-padded width of the numeric portion")

	return cmd
}

func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove components from your Grit project",
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
	var forceRelease bool

	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"self-update"},
		Short:   "Update the Grit CLI to the latest version",
		Long: `Update Grit to the latest version.

Behaviour:
  1. Checks the latest version on GitHub.
  2. If you're already on latest, exits without doing anything.
  3. If Go is installed, runs 'go install ...@latest' (fast, idempotent).
  4. If Go is NOT installed, downloads the matching binary from the
     GitHub release for your OS / arch and atomically swaps it in.

Flags:
  --from-release   Skip the 'go install' path and always pull the
                   binary directly from GitHub releases.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()
			purple := color.New(color.FgHiMagenta, color.Bold)
			green := color.New(color.FgHiGreen, color.Bold)
			yellow := color.New(color.FgHiYellow)
			spinner := color.New(color.FgHiBlack)

			purple.Printf("\n  Grit self-update — current: v%s\n\n", version)

			// 1. Check latest first. Both strategies (go install / GitHub binary)
			//    do the same thing when already on latest — namely, nothing. Do
			//    a single cheap HTTP round-trip up front and bail if there's
			//    nothing to do. Saves a full 'go install' cycle (which would
			//    otherwise hit the module proxy and recompile every time the
			//    user runs `grit update`).
			spinner.Println("  → Checking GitHub for the latest release...")
			latest, upToDate, verr := selfupdate.LatestVersion(version)
			if verr != nil {
				// Network failure shouldn't be fatal — fall through to the
				// strategy that may have a cached copy.
				yellow.Printf("  ! Couldn't query GitHub: %v\n", verr)
				yellow.Println("  ! Proceeding with the update anyway.")
			} else if upToDate {
				fmt.Println()
				green.Printf("  ✓ Already on the latest version (v%s). Nothing to do.\n\n", version)
				return nil
			} else {
				spinner.Printf("  → New version available: v%s → v%s\n", version, latest)
			}

			// 2. Pick a strategy: GitHub binary swap or `go install`.
			useRelease := forceRelease
			if !useRelease {
				if _, err := exec.LookPath("go"); err != nil {
					spinner.Println("  → Go toolchain not on PATH — using GitHub binary.")
					useRelease = true
				}
			}

			if useRelease {
				if err := selfupdate.Run(version); err != nil {
					return fmt.Errorf("self-update: %w", err)
				}
				fmt.Println()
				return nil
			}

			// 3. `go install` strategy — overwrites the binary in $GOBIN.
			binPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("finding current binary: %w", err)
			}
			binPath, err = filepath.EvalSymlinks(binPath)
			if err != nil {
				return fmt.Errorf("resolving binary path: %w", err)
			}

			// Binary-replacement strategy by OS:
			//
			//   Linux / macOS: POSIX keeps the running process's inode alive
			//                  even if the file at the path is overwritten,
			//                  so 'go install' can write straight over the
			//                  current binary. We don't touch anything — go
			//                  install handles it.
			//
			//   Windows:       .exe files are locked while running, so the
			//                  binary must be moved aside before 'go install'
			//                  can write. We rename to .old, run go install,
			//                  then delete the .old. On failure we rename
			//                  back so the user is never stranded without a
			//                  working binary.
			var rollback func() // populated on Windows; nil elsewhere
			if runtime.GOOS == "windows" {
				oldPath := binPath + ".old"
				os.Remove(oldPath) // any leftover from a prior run
				spinner.Printf("  → Moving running binary aside: %s → .old\n", binPath)
				if err := os.Rename(binPath, oldPath); err != nil {
					return fmt.Errorf("renaming current binary: %w", err)
				}
				rollback = func() {
					// If go install didn't write a new binary at binPath,
					// put the original back.
					if _, err := os.Stat(binPath); os.IsNotExist(err) {
						_ = os.Rename(oldPath, binPath)
					}
				}
			}

			spinner.Println("  → Running: go install github.com/MUKE-coder/grit/v3/cmd/grit@latest")
			c := exec.Command("go", "install", "github.com/MUKE-coder/grit/v3/cmd/grit@latest")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				if rollback != nil {
					rollback()
					yellow.Println("  ! Restored the original binary.")
				}
				return fmt.Errorf("installing latest version: %w", err)
			}

			if runtime.GOOS == "windows" {
				// go install succeeded — drop the rename'd backup.
				_ = os.Remove(binPath + ".old")
			}

			fmt.Println()
			if latest != "" {
				green.Printf("  ✓ Updated to v%s\n", latest)
			} else {
				green.Println("  ✓ Grit CLI updated successfully!")
			}
			spinner.Println("  Run 'grit version' to verify.")
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().BoolVar(&forceRelease, "from-release", false,
		"Skip 'go install' and download the binary directly from the GitHub release")
	return cmd
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
		Long:  "Launch GORM Studio at http://localhost:8080/studio. For web projects, ensure your API server is running first. For desktop, starts the studio server and opens the browser.",
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
	if !opts.InPlace {
		cyan.Printf("    cd %s\n", name)
	}
	cyan.Println("    docker compose up -d")

	switch opts.Architecture {
	case scaffold.ArchAPI:
		cyan.Println("    cd apps/api && go run cmd/server/main.go")
	case scaffold.ArchSingle:
		cyan.Println("    cd frontend && pnpm install")
		cyan.Println("    go run main.go")
	default:
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
	if opts.ShouldIncludeSingleSPA() {
		gray.Printf("  Frontend:    http://localhost:5173\n")
	}
	if opts.ShouldIncludeExpo() {
		gray.Printf("  Expo:        exp://localhost:8081\n")
	}
	if opts.ShouldIncludeDesktop() {
		gray.Printf("  Desktop:     wails dev (from apps/desktop)\n")
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

// ── grit routes ──────────────────────────────────────────────────────────────

func routesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "routes",
		Short: "List all registered API routes",
		Long:  "Parse the routes.go file and display a table of all registered HTTP routes with their methods, paths, handlers, and middleware groups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting working directory: %w", err)
			}

			routesFile, err := routeparser.FindRoutesFile(cwd)
			if err != nil {
				return err
			}

			routes, err := routeparser.Parse(routesFile)
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  API Routes (%s)\n\n", routesFile)

			fmt.Println(routeparser.FormatTable(routes))
			return nil
		},
	}
}

// ── grit down / grit up ─────────────────────────────────────────────────────

func downCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Put the application in maintenance mode",
		Long:  "Creates a .maintenance file that triggers the maintenance middleware, returning 503 for all requests.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if maintenance.IsEnabled(cwd) {
				color.Yellow("\n  Application is already in maintenance mode.\n")
				return nil
			}

			if err := maintenance.Enable(cwd); err != nil {
				return err
			}

			color.New(color.FgHiYellow, color.Bold).Println("\n  Application is now in maintenance mode.")
			color.New(color.FgHiBlack).Println("  All requests will receive 503 Service Unavailable.")
			color.New(color.FgHiBlack).Println("  Run 'grit up' to bring it back online.")
			fmt.Println()
			return nil
		},
	}
}

func upCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Bring the application back online",
		Long:  "Removes the .maintenance file, allowing normal request handling to resume.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := maintenance.Disable(cwd); err != nil {
				return err
			}

			color.New(color.FgHiGreen, color.Bold).Println("\n  Application is back online!")
			color.New(color.FgHiBlack).Println("  Normal request handling has resumed.")
			fmt.Println()
			return nil
		},
	}
}

// ── grit deploy ──────────────────────────────────────────────────────────────

func deployCmd() *cobra.Command {
	var host, port, keyFile, domain, appPort string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy application to a remote server",
		Long:  "Build the application, upload via SSH, configure systemd service, and optionally set up Caddy reverse proxy with auto-TLS.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			// Try to detect app name from go.mod or grit.config
			appName := "grit-app"
			if data, err := os.ReadFile("go.mod"); err == nil {
				for _, line := range strings.Split(string(data), "\n") {
					if strings.HasPrefix(line, "module ") {
						parts := strings.Fields(line)
						if len(parts) >= 2 {
							appName = filepath.Base(parts[1])
						}
						break
					}
				}
			}

			// Fall back to env vars if flags not set
			if host == "" {
				host = os.Getenv("DEPLOY_HOST")
			}
			if keyFile == "" {
				keyFile = os.Getenv("DEPLOY_KEY_FILE")
			}
			if domain == "" {
				domain = os.Getenv("DEPLOY_DOMAIN")
			}

			cfg := deploy.Config{
				Host:    host,
				Port:    port,
				KeyFile: keyFile,
				AppName: appName,
				Domain:  domain,
				AppPort: appPort,
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Deploying %s to %s\n\n", appName, host)

			if err := deploy.Run(cfg); err != nil {
				color.Red("\n  Deploy failed: %v\n", err)
				return err
			}

			green := color.New(color.FgHiGreen, color.Bold)
			green.Println("\n  Deployment successful!")
			if domain != "" {
				color.New(color.FgHiBlack).Printf("  Live at: https://%s\n\n", domain)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "SSH host (e.g. user@server.com) or DEPLOY_HOST env var")
	cmd.Flags().StringVar(&port, "port", "22", "SSH port")
	cmd.Flags().StringVar(&keyFile, "key", "", "Path to SSH private key or DEPLOY_KEY_FILE env var")
	cmd.Flags().StringVar(&domain, "domain", "", "Domain for Caddy reverse proxy or DEPLOY_DOMAIN env var")
	cmd.Flags().StringVar(&appPort, "app-port", "8080", "Port the app runs on")

	return cmd
}