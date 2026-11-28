package prompt

import (
	"github.com/charmbracelet/huh"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// Recommended is the project somebody gets by pressing Enter through every
// question, and the one /docs/start tells a newcomer to ask for by name:
//
//	grit new my-app --triple --next
//
// It is data rather than five accidents of list order because the picker and
// that page have to agree, and they did not. The architecture list opened on
// Full, so the documented recommendation and the thing a beginner got by
// pressing Enter were different projects, with nothing anywhere saying so
// (grit#240). TestEnterThroughThePromptGivesTheDocumentedProject reads both
// sides in one pass, because a test per side passes while the pair is broken.
var Recommended = struct {
	Architecture  string
	Frontend      string
	Theme         string
	DBProvider    string
	StorageDriver string
}{
	Architecture:  string(scaffold.ArchTriple),
	Frontend:      string(scaffold.FrontendNext),
	Theme:         scaffold.DefaultTheme,
	DBProvider:    scaffold.DefaultDBProvider,
	StorageDriver: scaffold.DefaultStorageDriver,
}

// recommendedSuffix marks the option Enter already lands on.
//
// Marking it is the whole point: a list of six architectures with nothing
// distinguishing them is a decision a newcomer has no basis for making, and
// the first question of their first command is the worst place to ask for one.
const recommendedSuffix = "  (recommended)"

// Hint is the line printed above the first question.
//
// Separate from printing it so a test can read it, and worded for the person
// who has no preference yet: the way out of five questions is Enter, and the
// project it produces is the one every guide assumes.
func Hint() string {
	return ui.Note(
		"New to Grit? Press Enter on every question.",
		"That gives you "+Recommended.Architecture+" + "+Recommended.Frontend+
			", the combination every guide here assumes.",
	)
}

// RunNewProjectPrompt shows an interactive prompt for project configuration.
// Returns the selected architecture and frontend. Skips prompts for fields
// that are already set (via CLI flags).
//
// Every picker is one huh.Form rather than a .Run() each, so the TUI renders
// once per step instead of once per question. On terminals with weak ANSI
// support (Git Bash / MINGW64), separate .Run() calls leave stale frames in
// the scrollback because the cursor-up codes do not fully apply.
//
// The same weakness is why ui.RunForm, and not form.Run, finishes the job: a
// form of several groups is drawn at its tallest group's height, and huh only
// learns that height from a message queued in Init, so the first frame is
// shorter than the rest. RunForm settles the size before the first frame.
func RunNewProjectPrompt(opts *scaffold.Options) error {
	arch := string(opts.Architecture)
	frontend := string(opts.Frontend)
	theme := opts.Theme
	db := opts.DBProvider
	storage := opts.StorageDriver

	// Open each picker on the recommended answer rather than on whatever
	// happens to be first in the list. huh writes the highlighted option on
	// submit, so this is what makes Enter mean the documented project.
	preselect(&arch, &frontend, &theme, &db, &storage)

	saas, orgs, billing := opts.SaaS, opts.SaaSOrganizations, opts.SaaSBilling

	form := newProjectForm(opts, &arch, &frontend, &theme, &db, &storage, &saas, &orgs, &billing)

	// One theme for every picker, from internal/ui: huh's default draws a
	// coloured border down the whole form and paints the options in its own
	// palette, which is a second design living beside the one the rest of the
	// output follows. RunForm also settles the size before the first frame,
	// which is what stops the list being drawn twice on Git Bash.
	if err := ui.RunForm(form); err != nil {
		return err
	}

	applyAnswers(opts, arch, frontend, theme, db, storage)
	opts.SaaS, opts.SaaSOrganizations, opts.SaaSBilling = saas, orgs, billing
	// Neither follow-up means anything on its own: a project that said no to
	// customers is not getting an organizations plugin because the cursor
	// happened to pass over the question.
	if !saas {
		opts.SaaSOrganizations, opts.SaaSBilling = false, false
	}
	return nil
}

// preselect fills in the recommended answer for anything still unset.
//
// Only the empty ones: a value already there came from a flag, and the group
// asking about it is hidden.
func preselect(arch, frontend, theme, db, storage *string) {
	for ptr, want := range map[*string]string{
		arch:     Recommended.Architecture,
		frontend: Recommended.Frontend,
		theme:    Recommended.Theme,
		db:       Recommended.DBProvider,
		storage:  Recommended.StorageDriver,
	} {
		if *ptr == "" {
			*ptr = want
		}
	}
}

// applyAnswers writes the picked values back onto the options.
//
// Separate from running the form so the test that drives the form as a
// bubbletea model ends up with the same Options a real `grit new` would, and
// not with its own idea of what the answers mean.
func applyAnswers(opts *scaffold.Options, arch, frontend, theme, db, storage string) {
	if arch == "full" {
		// Full is a shorthand, not an architecture. Normalize() turns opts.Full
		// into triple + web + admin + docs + expo + desktop.
		opts.Full = true
	} else {
		opts.Architecture = scaffold.Architecture(arch)
	}
	opts.Frontend = scaffold.Frontend(frontend)
	opts.Theme = theme
	opts.DBProvider = db
	opts.StorageDriver = storage
}

// newProjectForm builds the pickers.
//
// Separate from running them so a test can drive the form as the bubbletea
// model it is, which is the only way to see what a frame actually measures.
func newProjectForm(
	opts *scaffold.Options,
	arch, frontend, theme, db, storage *string,
	saas, orgs, billing *bool,
) *huh.Form {
	// "full" is a sentinel, not an Architecture enum: it maps to opts.Full,
	// which Normalize() expands to triple + web + admin + docs + expo + desktop.
	needsFrontend := func() bool {
		a := scaffold.Architecture(*arch)
		return *arch == "full" || a == scaffold.ArchSingle || a == scaffold.ArchDouble || a == scaffold.ArchTriple
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("arch").
				Title("Select architecture").
				Description("Press Enter to take the recommendation; every guide assumes it.").
				Options(
					huh.NewOption("Full — Web + Admin + API + Docs + Expo + Desktop (everything)", "full"),
					huh.NewOption(mark("Triple — Web + Admin + API (Turborepo)", string(scaffold.ArchTriple), Recommended.Architecture), string(scaffold.ArchTriple)),
					huh.NewOption("Double — Web + API (Turborepo)", string(scaffold.ArchDouble)),
					huh.NewOption("Single — Go API + embedded React SPA (one binary)", string(scaffold.ArchSingle)),
					huh.NewOption("API Only — Go API (no frontend)", string(scaffold.ArchAPI)),
					huh.NewOption("Mobile — API + Expo (React Native)", string(scaffold.ArchMobile)),
				).
				Value(arch),
		).WithHideFunc(func() bool { return opts.Architecture != "" }),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("frontend").
				Title("Select frontend framework").
				Options(
					huh.NewOption(mark("Next.js — SSR, SEO, App Router", string(scaffold.FrontendNext), Recommended.Frontend), string(scaffold.FrontendNext)),
					huh.NewOption(mark("TanStack Router — Vite, fast builds, small bundle (SPA)", string(scaffold.FrontendTanStack), Recommended.Frontend), string(scaffold.FrontendTanStack)),
				).
				Value(frontend),
		).WithHideFunc(func() bool {
			return opts.Frontend != "" || !needsFrontend()
		}),

		// Theme picker — runs for any architecture that includes a frontend.
		// The choice writes THEME=<name> to .env so the dashboard and auth
		// pages render with matching tokens, fonts and layout.
		//
		// Every theme in scaffold.ValidThemes belongs here. The list started
		// at three and grew to eight, and the five that arrived in v3.322.0
		// were reachable only with --theme: the picker is what most people
		// see, so for them the new layouts may as well not have shipped.
		// TestThePickerOffersEveryTheme keeps the two lists together.
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("theme").
				Title("Select visual theme").
				Description("Drives auth layout, dashboard tokens, fonts, and brand colors.").
				Options(
					huh.NewOption(mark("Atlas — split-screen, blue/white, team/organisation (Inter)", "atlas", Recommended.Theme), "atlas"),
					huh.NewOption("Aurora — Apple-inspired, monochrome black/white/grey (Geist)", "aurora"),
					huh.NewOption("Pulse — Cloudflare-inspired, premium blue, elevated cards (Onest)", "pulse"),
					huh.NewOption("Coral — a card over a blurred glimpse of the app, rose (Inter)", "coral"),
					huh.NewOption("Amber — a plain boxed form under a wordmark, amber (Inter)", "amber"),
					huh.NewOption("Sky — one bold heading, social sign-in first, crisp blue (Inter)", "sky"),
					huh.NewOption("Mono — a fine grid beside a panel of proof, black/white (Inter)", "mono"),
					huh.NewOption("Emerald — a form beside a customer quote, green (Inter)", "emerald"),
				).
				Value(theme),
		).WithHideFunc(func() bool {
			return opts.Theme != "" || !needsFrontend()
		}),

		// The database. Asked of every architecture, because the API is the one
		// thing every shape has.
		//
		// Postgres stays the default and the first option, but it needs a
		// server: without one, a new project migrates into nothing and the
		// first thing anybody sees is a connection refused. SQLite is the
		// shortest path to a running app and it was reachable only by reading
		// the help text for --db.
		//
		// Options come from scaffold.DBProviders, so this list and the flag
		// cannot drift. TestThePickerOffersEveryDatabase keeps them together.
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("db").
				Title("Select database").
				Description("Written to .env as DB_PROVIDER. Change it later by editing that line.").
				Options(databaseOptions()...).
				Value(db),
		).WithHideFunc(func() bool { return opts.DBProvider != "" }),

		// Where uploads go. Asked of every architecture for the same reason as
		// the database: the API is the one thing every shape has, and every
		// shape can take an upload.
		//
		// minio stays the default and the first option, because it is what
		// docker-compose.yml starts. local is second because it is the one
		// that needs nothing: no Docker, no bucket, no credentials. It worked
		// all along and was reachable only by reading a comment in .env after
		// the project already existed.
		//
		// Options come from scaffold.StorageDrivers, so this list and the flag
		// cannot drift. TestThePickerOffersEveryStorageDriver keeps them
		// together.
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("storage").
				Title("Select file storage").
				Description("Written to .env as STORAGE_DRIVER. s3, r2 and b2 need credentials before uploads work.").
				Options(storageOptions()...).
				Value(storage),
		).WithHideFunc(func() bool { return opts.StorageDriver != "" }),

		// Who the app is for.
		//
		// Asked last, because it is the only question whose answer is about
		// the product rather than the stack, and because the two that follow
		// it only make sense once it has been answered.
		//
		// Hidden for a shape with no browser to sign into. An API or a mobile
		// backend can have customers, but the customer area this turns on is
		// apps/web, and offering it where it cannot be written is a promise
		// the scaffold then breaks.
		huh.NewGroup(
			huh.NewConfirm().
				Key("saas").
				Title("Does this app have customers?").
				Description("Yes adds the signed-in customer area, and a resource owned by a user\nthen gets a customer screen as well as an admin one.").
				Affirmative("Yes, customers sign in").
				Negative("No, operators only").
				Value(saas),
		).WithHideFunc(func() bool { return !needsFrontend() }),

		huh.NewGroup(
			huh.NewConfirm().
				Key("saas-organizations").
				Title("Do customers belong to organizations?").
				Description("Teams, invitations, per-organization roles, and query scoping so one\norganization cannot read another's rows. The multitenant plugin.").
				Affirmative("Yes, organizations").
				Negative("No, one account each").
				Value(orgs),
		).WithHideFunc(func() bool { return !*saas || !needsFrontend() }),

		huh.NewGroup(
			huh.NewConfirm().
				Key("saas-billing").
				Title("Do you charge them?").
				Description("Stripe: server-priced checkout, the hosted billing portal, and webhooks\nthat settle it. Billing belongs to the organization when there is one.").
				Affirmative("Yes, take payments").
				Negative("Not yet").
				Value(billing),
		).WithHideFunc(func() bool { return !*saas || !needsFrontend() }),
	)
}

// mark appends the recommended suffix when this option is the recommended one.
//
// Taking the value and the recommendation rather than a bool so the label and
// the thing being recommended cannot come apart: move the recommendation and
// every label follows it.
func mark(label, value, recommended string) string {
	if value != recommended {
		return label
	}
	return label + recommendedSuffix
}

// storageOptions builds the picker's list from the providers the flag accepts.
//
// Labelled "<name> — <what it means>", with the description
// scaffold.StorageDrivers already holds, so the two say the same thing in both
// places.
func storageOptions() []huh.Option[string] {
	labels := map[string]string{
		"minio": "MinIO",
		"local": "Local disk",
		"s3":    "AWS S3",
		"r2":    "Cloudflare R2",
		"b2":    "Backblaze B2",
	}
	options := make([]huh.Option[string], 0, len(scaffold.StorageDriverOrder))
	for _, name := range scaffold.StorageDriverOrder {
		label := labels[name]
		if label == "" {
			label = name
		}
		options = append(options, huh.NewOption(
			mark(label+" — "+scaffold.StorageDrivers[name], name, Recommended.StorageDriver), name))
	}
	return options
}

// databaseOptions builds the picker's list from the engines the flag accepts.
//
// Labelled "<name> — <what it means>", with the description scaffold.DBProviders
// already holds, so the two say the same thing in both places.
func databaseOptions() []huh.Option[string] {
	labels := map[string]string{
		"postgres": "PostgreSQL",
		"mysql":    "MySQL / MariaDB",
		"sqlite":   "SQLite",
		"memory":   "In-memory",
	}
	options := make([]huh.Option[string], 0, len(scaffold.DBProviderOrder))
	for _, name := range scaffold.DBProviderOrder {
		label := labels[name]
		if label == "" {
			label = name
		}
		options = append(options, huh.NewOption(
			mark(label+" — "+scaffold.DBProviders[name], name, Recommended.DBProvider), name))
	}
	return options
}
