package setup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/accessgen"
	"github.com/MUKE-coder/grit/v3/internal/plugin"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// Run is the whole command: ask, scaffold, configure, install, migrate,
// commit, and optionally start.
//
// Every step after the scaffold is a real command a person could run by hand.
// Keeping that shape is what makes the output useful when a step fails: setup
// says which command to run next, and the project is already on disk.
func Run(ctx context.Context, opts Options) error {
	if !opts.Yes {
		if os.Getenv("CI") != "" {
			return fmt.Errorf("grit setup is interactive, and CI is set; pass --yes with your answers as flags")
		}
		if !stdinIsTerminal() {
			return fmt.Errorf("grit setup needs an interactive terminal; for scripts, pass --yes with your answers as flags")
		}
		if err := Prompt(&opts); err != nil {
			return err
		}
	}

	if err := opts.Normalize(); err != nil {
		return err
	}

	if !opts.Yes {
		ok, err := Confirm(opts)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCancelled
		}
	}

	scaffoldOpts := opts.ScaffoldOptions()
	target := projectRoot(opts)
	if err := ensureNewProject(target); err != nil {
		return err
	}

	started := time.Now()
	if err := scaffold.Run(scaffoldOpts); err != nil {
		return err
	}

	root, err := filepath.Abs(projectRoot(opts))
	if err != nil {
		return fmt.Errorf("resolving project path: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.Section("Configuring"))

	// The blank starter is the scaffold minus its Blog example. It goes through
	// `grit remove resource`, which knows every file the generator wrote,
	// including the public pages and the sidebar entry.
	if opts.Starter == StarterBlank {
		if opts.StarterRemovalNeeded(root) {
			if err := RemoveBlog(ctx, root); err != nil {
				return fmt.Errorf("removing the Blog example: %w\n\nRetry with: cd %s && grit remove resource Blog --force", err, displayPath(opts, root))
			}
			// The single layouts put the panel and the public site in roots the
			// removal does not know yet; a monorepo still leaves the navbar
			// links behind. See starter.go.
			if err := cleanBlankStarter(root, opts.Version, opts); err != nil {
				return fmt.Errorf("finishing the blank starter: %w", err)
			}
			fmt.Println(ui.Step("Starter", "blank, Blog example removed"))
		} else {
			fmt.Println(ui.Step("Starter", "blank"))
		}
	} else {
		fmt.Println(ui.Step("Starter", "Blog, kept as grit new writes it"))
	}

	if opts.WithI18n {
		res, err := scaffold.AddI18n(root, false)
		if err != nil {
			return fmt.Errorf("adding internationalisation: %w\n\nRetry with: cd %s && grit add i18n", err, displayPath(opts, root))
		}
		fmt.Println(ui.Step("Internationalisation", fmt.Sprintf("%d files, %d wirings", len(res.Written), len(res.Wired))))
	}

	installed, err := InstallPlugins(ctx, root, opts)
	if err != nil {
		return err
	}
	for _, p := range installed {
		fmt.Println(ui.Step("Plugin", p.Name))
	}

	// Read the access table after the starter, i18n and plugins have finished
	// changing routes. A blank project must not keep the Blog's permissions,
	// and a plugin's routes must not be missing from the first admin session.
	if res, err := accessgen.Generate(root); err != nil {
		fmt.Println(ui.WarningBlock("Access table not written", err.Error(),
			"The project is fine. Run `grit doctor` inside it for the detail."))
	} else if res.Path != "" {
		fmt.Println(ui.Step("Access table", fmt.Sprintf("%d routes", res.Routes)))
	}

	if err := PatchServiceEnv(root, opts); err != nil {
		return fmt.Errorf("writing service settings to .env: %w", err)
	}
	fmt.Println(ui.Step("Environment", strings.Join([]string{
		redisTitle(opts.Redis),
		storageTitle(opts.Storage),
		mailTitle(opts.Mail),
	}, " · ")))

	services := opts.LocalServices()
	if len(services) > 0 {
		if err := StartLocalServices(ctx, root, services); err != nil {
			if toolMissing(err) {
				fmt.Println(ui.WarningBlock("Local services not started", err.Error(),
					"Install Docker, or start the services another way, then run:",
					"  docker compose up -d "+strings.Join(services, " ")))
			} else {
				fmt.Println(ui.WarningBlock("Local services did not start", err.Error(),
					"The project is complete. Start them with:",
					"  docker compose up -d "+strings.Join(services, " ")))
			}
		} else {
			fmt.Println(ui.Step("Local services", composeServicesLabel(services)))
		}
	}

	// Dependencies resolve after plugins, because a plugin's `go get` changes
	// go.mod after the scaffold's own tidy and a Node dependency is only in
	// package.json until pnpm reads it.
	apiDir := scaffoldOpts.APIRoot(root)
	needsTidy := scaffold.DependenciesPending() != ""
	if pending := scaffold.DependenciesPending(); pending != "" {
		if !filepath.IsAbs(pending) {
			apiDir = filepath.Join(root, pending)
		} else {
			apiDir = pending
		}
	}
	if needsTidy || len(installed) > 0 {
		if err := TidyGoDeps(ctx, apiDir); err != nil {
			return fmt.Errorf("resolving Go dependencies: %w\n\nRetry with: cd %s && go mod tidy", err, apiDir)
		}
		needsTidy = false
		fmt.Println(ui.Step("Go dependencies", "resolved"))
	}

	canStart := true
	needsFrontend := scaffoldOpts.ShouldIncludeFrontend()
	needsPnpm := needsFrontend && !opts.InstallDeps
	if needsFrontend {
		if !opts.InstallDeps {
			canStart = false
			fmt.Println(ui.WarningLine("Frontend dependencies skipped", "run pnpm install before grit start"))
		} else if err := InstallFrontendDeps(ctx, root, opts.Yes); err != nil {
			// A flaky network here must not cost the migration, the seed or the
			// first commit: none of them needs node. Report it, name the retry,
			// and carry on with what still works.
			canStart = false
			needsPnpm = true
			if errors.Is(err, ErrPnpmUnavailable) {
				fmt.Println(ui.WarningBlock("pnpm is not installed", err.Error(),
					"Install pnpm, then run:",
					"  cd "+displayPath(opts, root)+" && pnpm install"))
			} else {
				fmt.Println(ui.WarningBlock("Frontend dependencies did not install", err.Error(),
					"The project is complete. Retry with:",
					"  cd "+displayPath(opts, root)+" && pnpm install"))
			}
		} else {
			fmt.Println(ui.Step("Frontend dependencies", "installed"))
		}
	}

	needsMigrate := !opts.Migrate
	needsSeed := !opts.Seed
	if !opts.Migrate {
		canStart = false
	}
	if opts.Migrate {
		if err := RunMigrate(ctx, root); err != nil {
			needsMigrate = true
			fmt.Println()
			var cmdErr *CommandError
			if errors.As(err, &cmdErr) {
				if tail := tailOutput(cmdErr.Output, 30); strings.TrimSpace(tail) != "" {
					fmt.Println(ui.Section("Migration output"))
					fmt.Println()
					fmt.Println(tail)
					fmt.Println()
				}
			}
			fmt.Println(ui.FailureBlock("Migrations stopped", "The project was created, but the database did not finish migrating.",
				"Fix the cause above, then run: cd "+displayPath(opts, root)+" && grit migrate"))
			printNextSteps(opts, root, false, needsTidy, needsPnpm, needsMigrate, needsSeed, "")
			return errors.New("migrations did not finish")
		}
		needsMigrate = false
		fmt.Println(ui.Step("Database", "migrated"))
	}

	if opts.Seed {
		if err := RunSeed(ctx, root); err != nil {
			needsSeed = true
			var cmdErr *CommandError
			if errors.As(err, &cmdErr) {
				if tail := tailOutput(cmdErr.Output, 15); strings.TrimSpace(tail) != "" {
					fmt.Println(ui.WarningBlock("Sample data was not seeded", "The seeder stopped.",
						strings.Split(tail, "\n")...))
				} else {
					fmt.Println(ui.WarningLine("Sample data was not seeded", "run grit seed after the project starts"))
				}
			} else {
				fmt.Println(ui.WarningBlock("Sample data was not seeded", err.Error(),
					"The project and its tables are ready. Seed later with:",
					"  grit seed"))
			}
		} else {
			needsSeed = false
			fmt.Println(ui.Step("Sample data", "seeded"))
		}
	}

	if opts.Git {
		res, err := InitGit(ctx, root)
		if err != nil {
			fmt.Println(ui.WarningBlock("Git setup stopped", err.Error(),
				"The project is complete. Initialise and commit when you are ready."))
		} else if res.Committed {
			fmt.Println(ui.Step("Git", "initial commit"))
		} else if res.Skipped != "" {
			fmt.Println(ui.WarningLine("Git skipped", res.Skipped))
		}
	}

	if len(installed) > 0 {
		printPluginNotes(installed)
	}

	detail := []string{architectureTitle(opts.Architecture)}
	if opts.NeedsWebFrontend() {
		detail = append(detail, frontendTitle(opts.Frontend))
	}
	detail = append(detail, opts.DBProvider)
	if d := ui.Duration(time.Since(started)); d != "" {
		detail = append(detail, d)
	}
	fmt.Println()
	fmt.Println(ui.Result("Ready", opts.ProjectName, strings.Join(detail, " · ")))

	star, err := AskStar(opts)
	if err != nil {
		return err
	}
	if star {
		if err := StarGrit(ctx); err != nil {
			fmt.Println(ui.WarningBlock("Could not star Grit", err.Error(),
				"Star it at: https://github.com/MUKE-coder/grit"))
		} else {
			fmt.Println(ui.Step("GitHub", "Grit starred"))
		}
	}

	start, err := AskStart(opts)
	if err != nil {
		return err
	}
	if start && !canStart {
		fmt.Println(ui.WarningLine("Not starting yet", "frontend dependencies are not installed"))
		start = false
	}
	if start {
		fmt.Println()
		fmt.Println(ui.Note("Starting", "Ctrl+C stops the API and every frontend."))
		fmt.Println()
		return startForeground(root)
	}

	adminLogin := ""
	if opts.Seed && !needsSeed {
		adminLogin = adminCredentials(root)
	}
	printNextSteps(opts, root, canStart, needsTidy, needsPnpm, needsMigrate, needsSeed, adminLogin)
	return nil
}

// ensureNewProject is the one promise setup makes that grit new does not:
// this command only creates projects, so it must not be pointed at one that
// already exists.
func ensureNewProject(root string) error {
	if _, err := os.Stat(filepath.Join(root, "grit.json")); err == nil {
		return fmt.Errorf("%s is already a Grit project; grit setup only creates new projects", displayRoot(root))
	}
	return nil
}

func projectRoot(opts Options) string {
	if opts.InPlace {
		return "."
	}
	return opts.ProjectName
}

func displayRoot(root string) string {
	abs, err := filepath.Abs(root)
	if err != nil {
		return root
	}
	if base, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(base, abs); err == nil && rel != "." {
			return rel
		}
	}
	return abs
}

func displayPath(opts Options, root string) string {
	if opts.InPlace {
		return "."
	}
	return filepath.Base(root)
}

func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func printPluginNotes(plugins []plugin.Plugin) {
	hasNotes := false
	for _, p := range plugins {
		if len(p.NextSteps) > 0 {
			hasNotes = true
			break
		}
	}
	if !hasNotes {
		return
	}

	fmt.Println()
	fmt.Println(ui.Section("Plugin notes"))
	fmt.Println()
	for _, p := range plugins {
		if len(p.NextSteps) == 0 {
			continue
		}
		fmt.Printf("  %s\n", p.Name)
		for _, step := range p.NextSteps {
			fmt.Printf("    %s\n", step)
		}
	}
}

func printNextSteps(opts Options, root string, canStart, needsTidy, needsPnpm, needsMigrate, needsSeed bool, adminLogin string) {
	fmt.Println()
	fmt.Println(ui.Section("Next steps"))
	fmt.Println()

	type step struct{ cmd, note string }
	var steps []step
	if !opts.InPlace {
		steps = append(steps, step{"cd " + opts.ProjectName, ""})
	}
	if needsTidy {
		steps = append(steps, step{"go mod tidy", "in the API folder: this did not resolve"})
	}
	if needsPnpm {
		steps = append(steps, step{"pnpm install", "frontend deps, one time"})
	}
	if needsMigrate {
		steps = append(steps, step{"grit migrate", "create database tables"})
	}
	if needsSeed {
		steps = append(steps, step{"grit seed", "sample data, optional"})
	}
	if canStart {
		if opts.Architecture == "api" {
			steps = append(steps, step{"grit start server", "run the API"})
		} else {
			steps = append(steps, step{"grit start", "run everything"})
		}
		if opts.IncludesExpo() {
			steps = append(steps, step{"grit start expo", "run the Expo app"})
		}
		if opts.IncludesDesktop() {
			steps = append(steps, step{"grit start desktop", "run the Wails app"})
		}
	}

	column := 0
	for _, s := range steps {
		if n := len([]rune(s.cmd)); n > column {
			column = n
		}
	}
	for _, s := range steps {
		fmt.Println(ui.Command(s.cmd, s.note, column+4))
	}

	fmt.Println()
	fmt.Println(ui.Rule(52))

	studioLogin := localCredentials(root, "GORM_STUDIO")
	sentinelLogin := localCredentials(root, "SENTINEL")

	type link struct{ label, href, cred string }
	links := []link{
		{"API", "http://localhost:8080", ""},
		{"API docs", "http://localhost:8080/docs", ""},
		{"GORM Studio", "http://localhost:8080/studio", studioLogin},
		{"Sentinel", "http://localhost:8080/sentinel/ui", sentinelLogin},
	}
	if opts.Architecture == "double" || opts.Architecture == "triple" || opts.Architecture == "full" || opts.Architecture == "desktop" {
		links = append(links, link{"Web app", "http://localhost:3000", ""})
	}
	if opts.HasSeparateAdmin() {
		links = append(links, link{"Admin", "http://localhost:3001", adminLogin})
	} else if opts.Architecture == "double" {
		// A double embeds the panel in the web app at /admin. The URL table
		// already lists the web app, so this row is what makes the login
		// usable without opening the README.
		links = append(links, link{"Admin", "http://localhost:3000/admin", adminLogin})
	}
	if opts.Architecture == "single" {
		links = append(links, link{"Frontend", "http://localhost:5173", ""})
		links = append(links, link{"Admin", "http://localhost:5173/admin", adminLogin})
	}
	if opts.IncludesExpo() {
		links = append(links, link{"Expo", "exp://localhost:8081", ""})
	}
	if opts.IncludesDesktop() {
		links = append(links, link{"Desktop", "wails dev, from apps/desktop", ""})
	}
	if opts.Architecture == "full" {
		links = append(links, link{"Docs", "http://localhost:3002", ""})
	}
	if opts.DBProvider == "postgres" {
		links = append(links, link{"PostgreSQL", "localhost:5434", ""})
	}
	if opts.DBProvider == "mysql" {
		links = append(links, link{"MySQL", "localhost:3306", ""})
	}
	if opts.Redis.Mode == RedisLocal {
		links = append(links, link{"Redis", "localhost:6380", ""})
	}
	if opts.Storage.Mode == StorageMinIO {
		links = append(links, link{"MinIO", "http://localhost:9003", ""})
	}
	if opts.Mail.Mode == MailMailhog {
		links = append(links, link{"Mailhog", "http://localhost:8025", ""})
	}

	labelColumn, hrefColumn := 0, 0
	for _, l := range links {
		if n := len([]rune(l.label)); n > labelColumn {
			labelColumn = n
		}
		if n := len([]rune(l.href)); n > hrefColumn {
			hrefColumn = n
		}
	}
	for _, l := range links {
		line := ui.URL(l.label, l.href, labelColumn+3)
		if l.cred != "" {
			line += strings.Repeat(" ", hrefColumn-len([]rune(l.href))+2) + "| " + l.cred
		}
		fmt.Println(line)
	}

	// The seeder creates the admin account. Showing it here is the difference
	// between a project somebody can open and one they have to read the seeder
	// to sign in to. A password the person set is not echoed.
	if adminLogin == "" && !opts.Seed {
		fmt.Println()
		fmt.Println(ui.Note("Admin login", "Run grit seed to create "+adminCredentials(root)+"."))
	}
	if adminLogin != "" || studioLogin != "" || sentinelLogin != "" {
		fmt.Println()
		fmt.Println(ui.Gutter + ui.Muted.Render("Development credentials. Change them before production."))
	}
	fmt.Println()
}
