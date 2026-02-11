package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/internal/generate"
	"github.com/MUKE-coder/grit/internal/scaffold"
)

var version = "0.5.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "grit",
		Short: "Grit — Go + React. Built with Grit.",
		Long:  "Grit is a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript).",
	}

	rootCmd.AddCommand(newCmd())
	rootCmd.AddCommand(generateCmd())
	rootCmd.AddCommand(syncCmd())
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

			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Creating new Grit project: %s\n\n", projectName)

			opts := scaffold.Options{
				ProjectName: projectName,
				APIOnly:     apiOnly,
				IncludeExpo: includeExpo,
				MobileOnly:  mobileOnly,
				Full:        full,
			}

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

func generateResourceCmd() *cobra.Command {
	var fromFile string
	var interactive bool
	var fields string

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

			gen, err := generate.NewGenerator(def)
			if err != nil {
				return err
			}

			return gen.Run()
		},
	}

	cmd.Flags().StringVar(&fromFile, "from", "", "YAML file defining the resource fields")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively define fields")
	cmd.Flags().StringVar(&fields, "fields", "", "Inline field definitions (e.g., \"title:string,content:text,published:bool\")")

	return cmd
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
	gray.Printf("  GORM Studio: http://localhost:8080/studio\n")

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
