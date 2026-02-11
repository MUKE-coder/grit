package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/internal/scaffold"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "grit",
		Short: "Grit — Go + React. Built with Grit.",
		Long:  "Grit is a full-stack meta-framework that fuses Go (Gin + GORM) with Next.js (React + TypeScript).",
	}

	rootCmd.AddCommand(newCmd())
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
	var apiOnly bool

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

			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Creating new Grit project: %s\n\n", projectName)

			opts := scaffold.Options{
				ProjectName: projectName,
				APIOnly:     apiOnly,
			}

			if err := scaffold.Run(opts); err != nil {
				color.Red("\n  Error: %v\n", err)
				return err
			}

			printSuccess(projectName, apiOnly)
			return nil
		},
	}

	cmd.Flags().BoolVar(&apiOnly, "api", false, "Scaffold only the Go API (no frontend apps)")

	return cmd
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
	gray.Println("  Go + React. Built with Grit.")
}

func printSuccess(name string, apiOnly bool) {
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

	if apiOnly {
		cyan.Println("    cd apps/api && go run cmd/server/main.go")
	} else {
		cyan.Println("    pnpm install")
		cyan.Println("    pnpm dev")
	}

	fmt.Println()
	gray.Println("  ─────────────────────────────────────")
	gray.Printf("  API:        http://localhost:8080\n")
	gray.Printf("  GORM Studio: http://localhost:8080/studio\n")

	if !apiOnly {
		gray.Printf("  Web App:    http://localhost:3000\n")
		gray.Printf("  Admin:      http://localhost:3001\n")
	}

	gray.Printf("  PostgreSQL: localhost:5432\n")
	gray.Printf("  Redis:      localhost:6379\n")
	gray.Printf("  MinIO:      http://localhost:9001\n")
	gray.Printf("  Mailhog:    http://localhost:8025\n")
	gray.Println("  ─────────────────────────────────────")
	fmt.Println()
}
