package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/deploy"
	"github.com/MUKE-coder/grit/v3/internal/expose"
	"github.com/MUKE-coder/grit/v3/internal/generate"
	"github.com/MUKE-coder/grit/v3/internal/maintenance"
	"github.com/MUKE-coder/grit/v3/internal/project"
	"github.com/MUKE-coder/grit/v3/internal/prompt"
	"github.com/MUKE-coder/grit/v3/internal/routeparser"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/selfupdate"
)

var version = "3.52.0"

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
	rootCmd.AddCommand(exposeCmd())
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(compileCmd())
	rootCmd.AddCommand(studioCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(backupCmd())
	rootCmd.AddCommand(restoreCmd())
	rootCmd.AddCommand(packageCmd())
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
	var archFlag, frontendFlag, style, theme string
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
				Theme:          theme,
				InPlace:        inPlace,
				Force:          force,
				IncludeDesktop: includeDesktop,
				Version:        version, // inject the current CLI version into scaffolded files
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

			if err := opts.ValidateTheme(); err != nil {
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
	cmd.Flags().StringVar(&theme, "theme", "", "Full theme: atlas (default), aurora, pulse — controls auth pages, dashboard, fonts, and brand colors. Can also be overridden at runtime via THEME=<name> in .env.")

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
	cmd.AddCommand(addWebAuthCmd())

	return cmd
}

// addWebAuthCmd — v3.31.22+ Phase 4 of PLAN_FORMS_AND_SHARING.
// Drops two files into apps/web/ that close the "how do I protect a
// customer-facing page" gap:
//   - middleware.ts          — SSR cookie check, redirects to /login
//   - components/ProtectedWebRoute.tsx — client wrapper using useMe()
//
// Idempotent: pre-existing files are skipped with a notice.
func addWebAuthCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "web-auth",
		Short: "Add page protection helpers to apps/web/",
		Long: "Scaffolds apps/web/middleware.ts (SSR cookie redirect) and\n" +
			"apps/web/components/ProtectedWebRoute.tsx (client guard).\n" +
			"Existing files are left alone unless --force is passed.\n\n" +
			"After running, declare which paths are protected in middleware.ts\n" +
			"by editing the matcher, or wrap a page with <ProtectedWebRoute>.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Println("\n  Adding web-auth helpers to apps/web/")

			return scaffold.AddWebAuth(root, force)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite middleware.ts / ProtectedWebRoute.tsx if they already exist")
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

// exposeCmd is the v3.31.21+ parent for "scaffold a public-facing
// Next.js page from an already-generated resource". Subcommands are
// form (a Create-styled page) and table (a paginated list page).
func exposeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Scaffold a web page that consumes an existing resource",
		Long: "Generates a single Next.js page in apps/web/ that uses the resource's shared Zod schema + the\n" +
			"generated React Query hook. Use to surface a Grit-managed resource on a public marketing page,\n" +
			"a customer dashboard, or anywhere outside the admin panel — without re-implementing the form/table.",
	}
	cmd.AddCommand(exposeFormCmd())
	cmd.AddCommand(exposeTableCmd())
	return cmd
}

func exposeFormCmd() *cobra.Command {
	var to string
	var force bool
	var publicShare bool
	var token string
	cmd := &cobra.Command{
		Use:   "form <Resource>",
		Short: "Scaffold a public form page for a resource",
		Long: "Emits a Next.js page that renders one input per resource field.\n\n" +
			"By default: submits via the authenticated useCreate<Resource>() hook —\n" +
			"the visitor must be signed in.\n\n" +
			"With --public-share: submits via /api/public/forms/<token>/submit — no\n" +
			"auth required. The token comes from a FormShare you created in the\n" +
			"admin (System → Public form sharing). Pass it via --token, or leave it\n" +
			"blank and the generated page reads NEXT_PUBLIC_FORM_TOKEN from .env.\n\n" +
			"Examples:\n" +
			"  grit expose form Contact --to apps/web/app/contact-us/page.tsx\n" +
			"  grit expose form Contact --to apps/web/app/contact-us/page.tsx --public-share --token 9CkLh7gJZ...",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			mode := "auth'd"
			if publicShare {
				mode = "public-share"
			}
			purple.Printf("\n  Exposing form for %s (%s) → %s\n\n", args[0], mode, to)

			if err := expose.Form(expose.Opts{
				Resource:    args[0],
				To:          to,
				Root:        root,
				Force:       force,
				PublicShare: publicShare,
				Token:       token,
			}); err != nil {
				return err
			}

			green := color.New(color.FgHiGreen)
			green.Printf("  ✓ Wrote %s\n\n", to)
			fmt.Println("  Next steps:")
			fmt.Println("    cd apps/web && pnpm dev")
			fmt.Printf("    open http://localhost:3000/%s\n", deriveURLFromPath(to))
			if publicShare && token == "" {
				fmt.Println()
				fmt.Println("  ⚠  --token was not provided. The generated page reads")
				fmt.Println("     NEXT_PUBLIC_FORM_TOKEN from your web app's .env — set it before")
				fmt.Println("     the page will work. Create the share in admin → System → Public")
				fmt.Println("     form sharing if you haven't already.")
			}
			fmt.Println()
			return nil
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Destination path (e.g. apps/web/app/contact-us/page.tsx) — required")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite the destination if it already exists")
	cmd.Flags().BoolVar(&publicShare, "public-share", false, "Submit via /api/public/forms/<token>/submit instead of the auth'd hook")
	cmd.Flags().StringVar(&token, "token", "", "FormShare token for the public submission. Optional — falls back to NEXT_PUBLIC_FORM_TOKEN at runtime when blank.")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func exposeTableCmd() *cobra.Command {
	var to string
	var force bool
	cmd := &cobra.Command{
		Use:   "table <Resource>",
		Short: "Scaffold a paginated list page for a resource",
		Long: "Emits a Next.js page that renders a paginated, searchable list of records using\n" +
			"use<Resources>(). Web-styled (plain Tailwind), not the heavy admin DataTable.\n\n" +
			"Example:\n  grit expose table Contact --to apps/web/app/contacts/page.tsx",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Exposing table for %s → %s\n\n", args[0], to)

			if err := expose.Table(expose.Opts{Resource: args[0], To: to, Root: root, Force: force}); err != nil {
				return err
			}

			green := color.New(color.FgHiGreen)
			green.Printf("  ✓ Wrote %s\n\n", to)
			fmt.Println("  Next steps:")
			fmt.Println("    cd apps/web && pnpm dev")
			fmt.Printf("    open http://localhost:3000/%s\n\n", deriveURLFromPath(to))
			return nil
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Destination path (e.g. apps/web/app/contacts/page.tsx) — required")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite the destination if it already exists")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

// deriveURLFromPath turns "apps/web/app/contact-us/page.tsx" into
// "contact-us" so the success message prints a usable URL. Best-effort —
// doesn't try to handle route groups or catch-all params.
func deriveURLFromPath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimSuffix(p, "/page.tsx")
	if i := strings.Index(p, "apps/web/app/"); i >= 0 {
		p = p[i+len("apps/web/app/"):]
	}
	return strings.TrimPrefix(p, "/")
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

// backupCmd shells out to the project's cmd/backup, which knows the database and
// object-storage config. Same pattern as `grit migrate`.
func backupCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Back up the entire database",
		Long: "Dump every registered model to a ZIP archive: one CSV per table, a dump.sql of\n" +
			"INSERTs in parent-to-child order, and a metadata.json manifest of row counts.\n\n" +
			"By default the archive is uploaded to object storage (R2 / S3 / MinIO) and indexed,\n" +
			"the same way the weekly cron does it. Pass --output to write a local file instead,\n" +
			"which needs no storage credentials.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}

			goArgs := []string{"run", "./cmd/backup"}
			if output != "" {
				abs, err := filepath.Abs(output)
				if err != nil {
					return err
				}
				goArgs = append(goArgs, "--output", abs)
			}

			c := exec.Command("go", goArgs...)
			c.Dir = apiDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			purple := color.New(color.FgHiMagenta, color.Bold)
			if output != "" {
				purple.Printf("\n  Backing up the database to %s...\n", output)
			} else {
				purple.Println("\n  Backing up the database to object storage...")
			}

			return c.Run()
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Write the archive to a local file instead of uploading it")

	return cmd
}

// restoreCmd replays an archive. The article's hardest-won lesson: a backup you
// have never restored is a rumour — so restore ships as a first-class command.
func restoreCmd() *cobra.Command {
	var noMigrate bool

	cmd := &cobra.Command{
		Use:   "restore <backup.zip>",
		Short: "Restore the database from a backup archive",
		Long: "Run migrations, then replay the archive's dump.sql inside a single transaction —\n" +
			"either every row lands or none does.\n\n" +
			"Point this at an EMPTY database: the archive carries data, not schema, and existing\n" +
			"rows will collide on their primary keys. Use --no-migrate if the schema already exists.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}

			// The child runs with cwd=apiDir, so a relative path would resolve
			// from the wrong directory.
			archive, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			if _, err := os.Stat(archive); err != nil {
				return fmt.Errorf("backup archive not found: %s", archive)
			}

			goArgs := []string{"run", "./cmd/restore"}
			if noMigrate {
				goArgs = append(goArgs, "--migrate=false")
			}
			goArgs = append(goArgs, archive)

			c := exec.Command("go", goArgs...)
			c.Dir = apiDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			color.New(color.FgHiMagenta, color.Bold).Printf("\n  Restoring the database from %s...\n", args[0])

			return c.Run()
		},
	}

	cmd.Flags().BoolVar(&noMigrate, "no-migrate", false, "Skip running migrations before restoring")

	return cmd
}

// packageCmd builds a distributable desktop installer. It's a thin, friendly
// wrapper over `wails build` — detects the desktop project, checks prerequisites,
// runs the build, and prints where the artifact landed. For a full versioned
// release (branded installers + GitHub release) use scripts/release-desktop.sh.
func packageCmd() *cobra.Command {
	var platform string
	var noInstaller, clean bool

	cmd := &cobra.Command{
		Use:   "package",
		Short: "Build a distributable desktop installer (.exe / .app / binary)",
		Long: "Compile the desktop app into a shippable artifact you can hand to a user.\n\n" +
			"On Windows this produces an NSIS installer (.exe) by default — the single file\n" +
			"someone double-clicks to install your app. On macOS/Linux it produces the\n" +
			"platform binary/app bundle. Run this from inside a `grit new-desktop` project.\n\n" +
			"It wraps `wails build`; for a full versioned release (branded installers +\n" +
			"GitHub release) use scripts/release-desktop.sh <version> instead.",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := project.DetectProject()
			if err != nil || info.Type != project.ProjectDesktop {
				return fmt.Errorf("not inside a desktop project — run this from a `grit new-desktop` app")
			}

			// wails is required; makensis is required for the Windows installer.
			if _, err := exec.LookPath("wails"); err != nil {
				return fmt.Errorf("wails is not installed or not on PATH — see https://wails.io/docs/gettingstarted/installation")
			}
			targetsWindows := strings.HasPrefix(platform, "windows") ||
				(platform == "" && runtime.GOOS == "windows")
			buildInstaller := targetsWindows && !noInstaller
			if buildInstaller {
				if _, err := exec.LookPath("makensis"); err != nil {
					return fmt.Errorf("makensis (NSIS) is not installed or not on PATH — needed to build the Windows installer.\n\n" +
						"  Install it:\n" +
						"    winget install NSIS.NSIS   (or)   choco install nsis   (or)   scoop install nsis\n" +
						"    Then add the NSIS folder (usually C:\\Program Files (x86)\\NSIS) to your PATH so 'makensis' resolves.\n" +
						"    Download instead: https://nsis.sourceforge.io\n\n" +
						"  Or run 'grit package --no-installer' to build just the raw .exe (no installer).")
				}
			}

			printLogo()
			purple := color.New(color.FgHiMagenta, color.Bold)

			wailsArgs := []string{"build"}
			if platform != "" {
				wailsArgs = append(wailsArgs, "-platform", platform)
			}
			if clean {
				wailsArgs = append(wailsArgs, "-clean")
			}
			if buildInstaller {
				wailsArgs = append(wailsArgs, "-nsis")
				purple.Println("\n  Building desktop installer (this can take a minute)...")
			} else {
				purple.Println("\n  Building desktop app...")
			}

			c := exec.Command("wails", wailsArgs...)
			c.Dir = info.Root
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				return fmt.Errorf("wails build failed: %w", err)
			}

			// Point at the artifacts. Wails writes everything to build/bin.
			binDir := filepath.Join(info.Root, "build", "bin")
			green := color.New(color.FgHiGreen, color.Bold)
			cyan := color.New(color.FgHiCyan)
			gray := color.New(color.FgHiBlack)
			green.Println("\n  ✓ Build complete")
			fmt.Println()
			if buildInstaller {
				cyan.Println("  Your installer is in build/bin/ (the *-installer.exe file).")
				gray.Println("  Hand that single file to anyone — double-click to install.")
			} else {
				cyan.Println("  Your app is in build/bin/.")
			}
			if entries, derr := os.ReadDir(binDir); derr == nil {
				for _, e := range entries {
					if e.IsDir() {
						continue
					}
					gray.Printf("    build/bin/%s\n", e.Name())
				}
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVar(&platform, "platform", "", "Target platform, e.g. windows/amd64, darwin/arm64 (default: host)")
	cmd.Flags().BoolVar(&noInstaller, "no-installer", false, "Build the raw binary only, skip the NSIS installer (Windows)")
	cmd.Flags().BoolVar(&clean, "clean", false, "Clean the build directory before building")

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
// airVersion pins the air hot-reload tool. Grit runs it via `go run` when it
// isn't already on PATH, so a Grit app hot-reloads out of the box — the user
// never has to `go install` anything. Bump this to upgrade the bundled air.
const airVersion = "v1.65.3"

// apiHotReloadArgv returns the command to run the Go API from apiDir with Go
// hot-reload, and a one-line note describing which path was chosen. Order:
//   1. a globally-installed `air` on PATH (fastest — respects the user's own)
//   2. air via `go run github.com/air-verse/air@<pinned>` — no install needed;
//      compiled once, then served from the build cache
// Both read the .air.toml that every Grit API ships with. There's no
// no-hot-reload fallback anymore: `go run` air always works when Go is present
// (the first run downloads air, which needs network — same as any first build).
func apiHotReloadArgv() (bin string, args []string, note string) {
	if p, err := exec.LookPath("air"); err == nil {
		return p, nil, "API hot-reload via air (.go files auto-rebuild)."
	}
	return "go",
		[]string{"run", "github.com/air-verse/air@" + airVersion},
		"API hot-reload via air (bundled — first run downloads it, then cached)."
}

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

			// v3.30.2: pin the install target to the resolved GitHub tag
			// instead of @latest. proxy.golang.org's view of @latest can lag
			// minutes behind a freshly-pushed tag, which produces the
			// surprising "self-update says X but binary reports Y" bug.
			// Falling back to @latest only when the GitHub lookup failed —
			// in that case the proxy is the only signal we have.
			target := "github.com/MUKE-coder/grit/v3/cmd/grit@latest"
			if latest != "" {
				target = "github.com/MUKE-coder/grit/v3/cmd/grit@v" + latest
			}

			// Force the install to land where grit ACTUALLY lives (e.g.
			// ~/.grit/bin from the install script), not the Go default GOBIN
			// (~/go/bin). Without this, `go install` writes the new binary to
			// GOBIN while we've just renamed the real one aside — leaving the
			// user's `grit` path empty. (Fixes the self-update bug where
			// `grit update` reported success but left no grit on PATH.)
			installDir := filepath.Dir(binPath)
			spinner.Printf("  → Running: go install %s  (GOBIN=%s)\n", target, installDir)
			c := exec.Command("go", "install", target)
			c.Env = append(os.Environ(), "GOBIN="+installDir)
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

			// Sanity check: ask the freshly-installed binary what version
			// it reports. If it doesn't match the tag we asked for, surface
			// it loudly — the user otherwise has no way to know they're
			// running a stale binary that lied about its version.
			if latest != "" {
				if got, err := installedVersion(binPath); err == nil && got != "" && got != latest {
					fmt.Println()
					yellow.Printf("  ! Expected v%s but the installed binary reports v%s.\n", latest, got)
					yellow.Println("    Your GOPATH/bin may have an older binary on PATH, or the Go module")
					yellow.Println("    proxy hasn't indexed the new tag yet. Try again in a minute, or run:")
					yellow.Printf("      go install github.com/MUKE-coder/grit/v3/cmd/grit@v%s\n", latest)
					fmt.Println()
					return nil
				}
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

// installedVersion runs `<binPath> version` and parses out the semver.
// Used by the self-update flow to verify that the freshly-installed
// binary actually matches the version we asked the Go proxy for. Returns
// the version without the leading 'v', or "" if the call fails or the
// output doesn't contain a recognisable version line.
//
// We intentionally exec a fresh process (rather than reading the in-memory
// `version` variable) — the running binary still has the old version
// baked in, so checking ourselves would always say "no change".
func installedVersion(binPath string) (string, error) {
	out, err := exec.Command(binPath, "version").Output()
	if err != nil {
		return "", err
	}
	// `grit version` prints something like:
	//   Grit CLI v3.30.2 ...
	// We scan for the first vX.Y.Z occurrence and strip the v.
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)
	if m := re.FindStringSubmatch(string(out)); len(m) >= 2 {
		return m[1], nil
	}
	return "", nil
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start development servers",
		Long: "With no args, starts every app in the project in parallel — the Go API, the " +
			"frontend apps (web/admin via pnpm dev), and, when apps/desktop exists, the Wails " +
			"desktop app too. Ctrl+C stops all of them.\n\n" +
			"Run a single app from anywhere in the project:\n" +
			"  grit start server    Go API only\n" +
			"  grit start web       web app (Next.js)\n" +
			"  grit start admin     admin panel (Next.js)\n" +
			"  grit start expo      Expo mobile app\n" +
			"  grit start desktop   Wails desktop app\n" +
			"  grit start client    all frontend apps via Turborepo\n\n" +
			"Inside a standalone desktop project (grit new-desktop), 'grit start' runs 'wails dev'.",
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
			// Web project — run server + client in parallel.
			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}
			return runDevPair(info.Root, apiDir)
		},
	}

	cmd.AddCommand(startClientCmd())
	cmd.AddCommand(startServerCmd())
	// Per-app starters, so `grit start <app>` from the root runs just one.
	cmd.AddCommand(startAppCmd("web", "Start the web app (Next.js)", "apps/web", "pnpm", []string{"dev"}, false))
	cmd.AddCommand(startAppCmd("admin", "Start the admin panel (Next.js)", "apps/admin", "pnpm", []string{"dev"}, false))
	cmd.AddCommand(startAppCmd("expo", "Start the Expo mobile app", "apps/expo", "pnpm", []string{"start"}, false))
	cmd.AddCommand(startAppCmd("desktop", "Start the Wails desktop app", "apps/desktop", "wails", []string{"dev"}, true))

	return cmd
}

// startAppCmd builds a `grit start <app>` subcommand that runs a single app's
// dev process from anywhere in the project. appSubdir is relative to the
// project root; requiresWails gates the desktop app on the Wails toolchain.
func startAppCmd(use, short, appSubdir, bin string, args []string, requiresWails bool) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			dir := filepath.Join(root, appSubdir)
			if _, err := os.Stat(dir); err != nil {
				return fmt.Errorf("this project has no %s app (expected %s/)", use, appSubdir)
			}
			if requiresWails {
				if _, err := exec.LookPath("wails"); err != nil {
					return fmt.Errorf("the Wails toolchain isn't on PATH — install it from https://wails.io, then run 'grit start %s'", use)
				}
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Starting %s...\n", use)

			c := exec.Command(bin, args...)
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			return c.Run()
		},
	}
}

// runDevPair launches the Go API and the pnpm dev pipeline in parallel,
// streams their stdout/stderr with a coloured prefix per source, and
// shuts both down cleanly on Ctrl+C. If either child exits first, the
// other is terminated so users don't end up with a zombie process.
//
// projectRoot: where `pnpm dev` runs (turbo picks up apps/web + admin).
// apiDir:      where `go run cmd/server/main.go` runs (apps/api).
func runDevPair(projectRoot, apiDir string) error {
	printLogo()
	purple := color.New(color.FgHiMagenta, color.Bold)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the API with Go hot-reload via air, bundled through `go run` so the
	// user never installs it. Uses a global `air` when present.
	apiBin, apiArgs, apiNote := apiHotReloadArgv()

	apiCmd := exec.CommandContext(ctx, apiBin, apiArgs...)
	apiCmd.Dir = apiDir
	apiCmd.Stdin = nil // child shouldn't fight for the terminal

	clientCmd := exec.CommandContext(ctx, "pnpm", "dev")
	clientCmd.Dir = projectRoot
	clientCmd.Stdin = nil

	// Prefix each line of output so a developer can tell which process said
	// what. ANSI colours mark the source even after copy-paste.
	procs := []devProc{
		{prefix: color.New(color.FgHiGreen).Sprint("[api] "), cmd: apiCmd},
		{prefix: color.New(color.FgHiCyan).Sprint("[web] "), cmd: clientCmd},
	}

	// If the monorepo includes the Wails desktop client, launch it too — so a
	// single `grit start` at the root brings up every app. Best-effort: skip
	// (with a note) when the Wails toolchain isn't installed, rather than
	// failing the whole command. Non-desktop projects are unaffected.
	desktopDir := filepath.Join(projectRoot, "apps", "desktop")
	desktopIncluded := false
	if _, err := os.Stat(filepath.Join(desktopDir, "wails.json")); err == nil {
		if _, err := exec.LookPath("wails"); err == nil {
			desktopCmd := exec.CommandContext(ctx, "wails", "dev")
			desktopCmd.Dir = desktopDir
			desktopCmd.Stdin = nil
			procs = append(procs, devProc{prefix: color.New(color.FgHiMagenta).Sprint("[desktop] "), cmd: desktopCmd})
			desktopIncluded = true
		} else {
			color.New(color.FgYellow).Println("  apps/desktop found, but 'wails' isn't on PATH — skipping the desktop app.")
			color.New(color.FgHiBlack).Println("  Install Wails (https://wails.io) to launch it with grit start.")
		}
	}

	if desktopIncluded {
		purple.Println("\n  Starting API + web + desktop in parallel...")
	} else {
		purple.Println("\n  Starting API + client apps in parallel...")
	}
	color.New(color.FgHiBlack).Println("  " + apiNote)
	color.New(color.FgHiBlack).Println("  Press Ctrl+C to stop everything.")
	fmt.Println()

	return runDevProcs(ctx, cancel, procs)
}

// devProc is one child in the parallel dev runner.
type devProc struct {
	prefix string
	cmd    *exec.Cmd
}

// runDevProcs starts every process, streams their prefixed output onto one
// terminal, forwards Ctrl+C / SIGTERM to all of them, and when any one exits
// (or the user interrupts) shuts the rest down so they leave together.
func runDevProcs(ctx context.Context, cancel context.CancelFunc, procs []devProc) error {
	var outWg sync.WaitGroup
	started := make([]*exec.Cmd, 0, len(procs))

	for _, p := range procs {
		out, _ := p.cmd.StdoutPipe()
		errPipe, _ := p.cmd.StderrPipe()
		if err := p.cmd.Start(); err != nil {
			for _, s := range started {
				_ = killProcess(s)
			}
			return fmt.Errorf("starting process: %w", err)
		}
		started = append(started, p.cmd)
		outWg.Add(2)
		go prefixCopy(&outWg, p.prefix, out, os.Stdout)
		go prefixCopy(&outWg, p.prefix, errPipe, os.Stderr)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// One waiter per process. doneCh is buffered to len(started) so no send
	// ever blocks (Wait is called exactly once per process).
	var procWg sync.WaitGroup
	doneCh := make(chan error, len(started))
	for _, c := range started {
		c := c
		procWg.Add(1)
		go func() { doneCh <- c.Wait(); procWg.Done() }()
	}

	select {
	case <-sigCh:
		color.New(color.FgHiBlack).Println("\n  Caught Ctrl+C — stopping all apps...")
	case err := <-doneCh:
		if err != nil {
			color.New(color.FgHiYellow).Printf("\n  A process exited: %v — stopping the rest.\n", err)
		} else {
			color.New(color.FgHiBlack).Println("\n  A process exited — stopping the rest.")
		}
	}

	// Polite Interrupt first, then cancel the context (kills on most platforms).
	for _, c := range started {
		if c.Process != nil {
			_ = c.Process.Signal(os.Interrupt)
		}
	}
	cancel()

	procWg.Wait()
	outWg.Wait()
	return nil
}

// prefixCopy streams r line-by-line into w, prefixing every line with
// prefix. Used so the API and the client share one terminal without the
// developer guessing whose output is whose.
func prefixCopy(wg *sync.WaitGroup, prefix string, r io.Reader, w io.Writer) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	// Default token cap is 64 KB which truncates fat lines from webpack /
	// turbo. Lift it so we don't lose error messages.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		fmt.Fprintln(w, prefix+scanner.Text())
	}
}

// killProcess sends SIGINT then falls back to Kill if the process is
// still alive — used in error paths where we never even got both
// children to Start.
func killProcess(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	_ = cmd.Process.Signal(os.Interrupt)
	return cmd.Process.Kill()
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
		Long:  "Runs the Go API server. Uses 'air' when it's on PATH (so .go file edits hot-reload), falls back to 'go run cmd/server/main.go' otherwise.",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiDir, err := findAPIDir()
			if err != nil {
				return err
			}

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Println("\n  Starting API server...")

			// Go hot-reload via air, bundled through `go run` (no install).
			bin, args2, note := apiHotReloadArgv()
			color.New(color.FgHiBlack).Println("  " + note)

			c := exec.Command(bin, args2...)
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

	gray.Printf("  PostgreSQL:  localhost:5434\n")
	gray.Printf("  Redis:       localhost:6380\n")
	gray.Printf("  MinIO:       http://localhost:9003\n")
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
	gray.Printf("  Desktop App:  http://localhost:34115  (Wails dev window)\n")
	gray.Printf("  Embedded API: http://127.0.0.1:34999  (curl / other clients)\n")
	gray.Printf("  Database:     SQLite (%s.db)\n", name)
	gray.Println("  ─────────────────────────────────────")
	fmt.Println()

	gray.Println("  To build a distributable installer:")
	cyan.Println("    grit package")
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
