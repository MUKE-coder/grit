package setup

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// ErrCancelled is what the wizard returns when somebody presses Ctrl+C or
// answers "Cancel" on the summary. It is not a failure: the command exits 0
// and says nothing was written.
var ErrCancelled = errors.New("setup cancelled")

// Prompt asks for everything the flags did not already answer.
//
// One huh.Form holds every group, for the same reason grit new uses one: a
// separate .Run() per question leaves stale frames in terminals with weak
// cursor support. Hidden groups are skipped, so a flag that answers a question
// does not make the wizard ask it again.
func Prompt(o *Options) error {
	name := o.ProjectName
	arch := o.Architecture
	if arch == "" {
		arch = "triple"
	}
	frontend := o.Frontend
	if frontend == "" {
		frontend = "next"
	}
	theme := o.Theme
	if theme == "" {
		theme = "atlas"
	}
	starter := string(o.Starter)
	if starter == "" {
		starter = string(StarterBlank)
	}
	db := o.DBProvider
	if db == "" {
		db = "postgres"
	}

	redisMode := string(o.Redis.Mode)
	if redisMode == "" {
		redisMode = string(RedisLocal)
	}
	redisURL := o.Redis.URL

	storageMode := string(o.Storage.Mode)
	if storageMode == "" {
		storageMode = string(StorageMinIO)
	}
	storageEndpoint := o.Storage.Endpoint
	storageAccess := o.Storage.AccessKey
	storageSecret := o.Storage.SecretKey
	storageBucket := o.Storage.Bucket
	storageRegion := o.Storage.Region
	storagePublicURL := o.Storage.PublicURL

	mailMode := string(o.Mail.Mode)
	if mailMode == "" {
		mailMode = string(MailMailhog)
	}
	mailFrom := o.Mail.From
	mailAPIKey := o.Mail.APIKey
	mailDomain := o.Mail.Domain
	mailEndpoint := o.Mail.Endpoint
	mailStream := o.Mail.Stream
	smtpHost := o.Mail.SMTPHost
	smtpPort := o.Mail.SMTPPort
	smtpUsername := o.Mail.SMTPUsername
	smtpPassword := o.Mail.SMTPPassword
	smtpEncryption := o.Mail.SMTPEncryption
	sesRegion := o.Mail.SESRegion
	awsAccessKey := o.Mail.AWSAccessKey
	awsSecretKey := o.Mail.AWSSecretKey

	plugins := append([]string(nil), o.Plugins...)
	i18n := o.WithI18n

	actions := o.actionSelection()
	git := o.Git

	needsFrontend := func() bool {
		return arch == "single" || arch == "double" || arch == "triple" || arch == "full" || arch == "desktop"
	}
	hasAdmin := needsFrontend

	groups := []*huh.Group{
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Project name").
				Description("Lowercase letters, numbers and hyphens. It becomes the folder and the Go module.").
				Placeholder("my-app").
				Value(&name).
				Validate(scaffold.ValidateProjectName),
		).WithHideFunc(func() bool { return o.ProjectName != "" }),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("arch").
				Title("Select architecture").
				Options(
					huh.NewOption("Full — Web + Admin + API + Docs + Expo + Desktop", "full"),
					huh.NewOption("Triple — Web + Admin + API (Turborepo)", "triple"),
					huh.NewOption("Desktop — Triple + a Wails desktop app", "desktop"),
					huh.NewOption("Double — Web + API, admin inside the web app", "double"),
					huh.NewOption("Single — API + an embedded frontend, one binary", "single"),
					huh.NewOption("API only — Go API, no frontend", "api"),
					huh.NewOption("Mobile — API + Expo", "mobile"),
				).
				Value(&arch),
		).WithHideFunc(func() bool { return o.Architecture != "" }),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("frontend").
				Title("Select frontend").
				Options(
					huh.NewOption("Next.js — SSR, SEO, App Router", "next"),
					huh.NewOption("TanStack Router — Vite, fast builds, small bundle", "tanstack"),
				).
				Value(&frontend),
		).WithHideFunc(func() bool {
			return o.Frontend != "" || !needsFrontend()
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("starter").
				Title("Select starter").
				Description("Blank removes the built-in Blog example. Blog keeps it exactly as grit new writes it.").
				Options(
					huh.NewOption("Blank — a clean project with no example resource", string(StarterBlank)),
					huh.NewOption("Blog — the built-in Blog example, admin and public pages included", string(StarterBlog)),
				).
				Value(&starter),
		).WithHideFunc(func() bool {
			return o.Starter != "" || !hasAdmin()
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("theme").
				Title("Select visual theme").
				Description("Drives auth layout, dashboard tokens, fonts and brand colours.").
				Options(themeOptions()...).
				Value(&theme),
		).WithHideFunc(func() bool {
			return o.Theme != "" || !needsFrontend()
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("db").
				Title("Select database").
				Description("Written to .env as DB_PROVIDER. SQLite and memory need no server.").
				Options(databaseOptions()...).
				Value(&db),
		).WithHideFunc(func() bool { return o.DBProvider != "" }),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("redis").
				Title("Redis").
				Description("Cache, background jobs and cron. Local starts the compose service.").
				Options(
					huh.NewOption("Local — in Docker, on port 6380", string(RedisLocal)),
					huh.NewOption("External — a server you already run", string(RedisExternal)),
					huh.NewOption("Off — no cache or background jobs", string(RedisOff)),
				).
				Value(&redisMode),
		).WithHideFunc(func() bool { return o.Redis.Mode != "" }),

		huh.NewGroup(
			huh.NewInput().
				Key("redis-url").
				Title("Redis URL").
				Placeholder("rediss://default:password@host:6379").
				Value(&redisURL).
				Validate(func(s string) error {
					s = strings.TrimSpace(s)
					if s == "" {
						return errors.New("an external Redis needs a URL")
					}
					if !strings.HasPrefix(s, "redis://") && !strings.HasPrefix(s, "rediss://") {
						return errors.New("the URL must start with redis:// or rediss://")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return redisMode != string(RedisExternal) || o.Redis.URL != ""
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("storage").
				Title("File storage").
				Description("Uploads. MinIO is the local compose service; S3 and R2 are managed buckets.").
				Options(
					huh.NewOption("MinIO — local, in Docker", string(StorageMinIO)),
					huh.NewOption("AWS S3 — a bucket you own", string(StorageS3)),
					huh.NewOption("Cloudflare R2 — a bucket you own", string(StorageR2)),
					huh.NewOption("Local disk — no external service", string(StorageLocal)),
					huh.NewOption("Off — leave uploads unconfigured", string(StorageOff)),
				).
				Value(&storageMode),
		).WithHideFunc(func() bool { return o.Storage.Mode != "" }),

		huh.NewGroup(
			huh.NewInput().
				Title("S3 endpoint").
				Description("Leave empty for AWS. Set it for another S3-compatible provider.").
				Placeholder("https://s3.us-east-1.amazonaws.com").
				Value(&storageEndpoint),
			huh.NewInput().
				Title("S3 access key").
				Value(&storageAccess).
				Validate(required("access key")),
			huh.NewInput().
				Title("S3 secret key").
				EchoMode(huh.EchoModePassword).
				Value(&storageSecret).
				Validate(required("secret key")),
			huh.NewInput().
				Title("S3 bucket").
				Value(&storageBucket).
				Validate(required("bucket")),
			huh.NewInput().
				Title("S3 region").
				Placeholder("us-east-1").
				Value(&storageRegion),
			huh.NewInput().
				Title("Public URL (optional)").
				Description("A CDN domain, or the bucket's public origin. Needed for images to render.").
				Value(&storagePublicURL),
		).WithHideFunc(func() bool {
			return storageMode != string(StorageS3) || (storageAccess != "" && storageSecret != "" && storageBucket != "")
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("R2 endpoint").
				Placeholder("https://<account>.r2.cloudflarestorage.com").
				Value(&storageEndpoint).
				Validate(required("endpoint")),
			huh.NewInput().
				Title("R2 access key").
				Value(&storageAccess).
				Validate(required("access key")),
			huh.NewInput().
				Title("R2 secret key").
				EchoMode(huh.EchoModePassword).
				Value(&storageSecret).
				Validate(required("secret key")),
			huh.NewInput().
				Title("R2 bucket").
				Value(&storageBucket).
				Validate(required("bucket")),
			huh.NewInput().
				Title("R2 region").
				Placeholder("auto").
				Value(&storageRegion),
			huh.NewInput().
				Title("R2 public URL").
				Description("Required for images: the S3 endpoint rejects unsigned browser requests.").
				Placeholder("https://pub-xxxx.r2.dev").
				Value(&storagePublicURL),
		).WithHideFunc(func() bool {
			return storageMode != string(StorageR2) || (storageAccess != "" && storageSecret != "" && storageBucket != "" && storageEndpoint != "")
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("mail").
				Title("Mail").
				Description("How the app sends reset links and notifications.").
				Options(
					huh.NewOption("Mailhog — local, in Docker", string(MailMailhog)),
					huh.NewOption("Log — write messages to the log, development only", string(MailLog)),
					huh.NewOption("Resend", string(MailResend)),
					huh.NewOption("Mailgun", string(MailMailgun)),
					huh.NewOption("Postmark", string(MailPostmark)),
					huh.NewOption("SendGrid", string(MailSendGrid)),
					huh.NewOption("Amazon SES", string(MailSES)),
					huh.NewOption("SMTP — any provider", string(MailSMTP)),
				).
				Value(&mailMode),
		).WithHideFunc(func() bool { return o.Mail.Mode != "" }),

		huh.NewGroup(
			huh.NewInput().
				Key("mail-from").
				Title("Mail from").
				Description("Leave empty for noreply@<project>.dev.").
				Placeholder("noreply@example.com").
				Value(&mailFrom).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return nil
					}
					if _, err := mail.ParseAddress(s); err != nil {
						return errors.New("enter a valid email address")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return o.Mail.From != "" }),

		huh.NewGroup(
			huh.NewInput().
				Title("Resend API key").
				EchoMode(huh.EchoModePassword).
				Value(&mailAPIKey).
				Validate(required("Resend API key")),
		).WithHideFunc(func() bool {
			return mailMode != string(MailResend) || o.Mail.APIKey != ""
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Mailgun domain").
				Placeholder("mg.example.com").
				Value(&mailDomain).
				Validate(required("Mailgun domain")),
			huh.NewInput().
				Title("Mailgun secret").
				EchoMode(huh.EchoModePassword).
				Value(&mailAPIKey).
				Validate(required("Mailgun secret")),
			huh.NewInput().
				Title("Mailgun endpoint").
				Description("api.eu.mailgun.net for a domain in the EU.").
				Placeholder("api.mailgun.net").
				Value(&mailEndpoint),
		).WithHideFunc(func() bool {
			return mailMode != string(MailMailgun) || (o.Mail.Domain != "" && o.Mail.APIKey != "")
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Postmark server token").
				EchoMode(huh.EchoModePassword).
				Value(&mailAPIKey).
				Validate(required("Postmark server token")),
			huh.NewInput().
				Title("Message stream").
				Placeholder("outbound").
				Value(&mailStream),
		).WithHideFunc(func() bool {
			return mailMode != string(MailPostmark) || o.Mail.APIKey != ""
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("SendGrid API key").
				EchoMode(huh.EchoModePassword).
				Value(&mailAPIKey).
				Validate(required("SendGrid API key")),
		).WithHideFunc(func() bool {
			return mailMode != string(MailSendGrid) || o.Mail.APIKey != ""
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("SES region").
				Placeholder("us-east-1").
				Value(&sesRegion),
			huh.NewInput().
				Title("AWS access key").
				Value(&awsAccessKey).
				Validate(required("AWS access key")),
			huh.NewInput().
				Title("AWS secret key").
				EchoMode(huh.EchoModePassword).
				Value(&awsSecretKey).
				Validate(required("AWS secret key")),
		).WithHideFunc(func() bool {
			return mailMode != string(MailSES) || (o.Mail.AWSAccessKey != "" && o.Mail.AWSSecretKey != "")
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("SMTP host").
				Placeholder("smtp.example.com").
				Value(&smtpHost).
				Validate(required("SMTP host")),
			huh.NewInput().
				Title("SMTP port").
				Placeholder("587").
				Value(&smtpPort),
			huh.NewInput().
				Title("SMTP username").
				Value(&smtpUsername),
			huh.NewInput().
				Title("SMTP password").
				EchoMode(huh.EchoModePassword).
				Value(&smtpPassword),
			huh.NewInput().
				Title("SMTP encryption").
				Description("tls, starttls or none. Empty chooses from the port.").
				Placeholder("starttls").
				Value(&smtpEncryption),
		).WithHideFunc(func() bool {
			return mailMode != string(MailSMTP) || o.Mail.SMTPHost != ""
		}),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("plugins").
				Title("Plugins").
				Description("Space to select, Enter to continue. Only the plugins this shape can run are listed.").
				Options(pluginOptions(o, arch)...).
				Value(&plugins),
		).WithHideFunc(func() bool { return len(pluginOptions(o, arch)) == 0 }),

		huh.NewGroup(
			huh.NewConfirm().
				Key("i18n").
				Title("Add internationalisation?").
				Description("next-intl, translated API messages and en/fr/sw catalogues.").
				Affirmative("Yes").
				Negative("No").
				Value(&i18n),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("actions").
				Title("After scaffolding").
				Description("Space to keep or skip, Enter to continue.").
				Options(
					huh.NewOption("Start local services in Docker", "services"),
					huh.NewOption("Install frontend dependencies (pnpm install)", "deps"),
					huh.NewOption("Run database migrations", "migrate"),
					huh.NewOption("Seed sample data", "seed"),
				).
				Value(&actions),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Key("git").
				Title("Initialise git and make the first commit?").
				Description("Skipped when the folder is already inside a repository. .env stays ignored.").
				Affirmative("Yes").
				Negative("No").
				Value(&git),
		),
	}

	// RunForm, not Run: a form of several groups is drawn at its tallest
	// group's height, and huh only learns that height from a message that
	// arrives after the first frame. The repaint at the new height leaves the
	// first one on screen in Git Bash.
	if err := ui.RunForm(huh.NewForm(groups...)); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return ErrCancelled
		}
		return err
	}

	o.ProjectName = strings.TrimSpace(name)
	o.Architecture = arch
	o.Frontend = frontend
	o.Theme = theme
	o.Starter = Starter(starter)
	o.DBProvider = db
	o.Redis = RedisConfig{Mode: RedisMode(redisMode), URL: strings.TrimSpace(redisURL)}
	o.Storage = StorageConfig{
		Mode:      StorageMode(storageMode),
		Endpoint:  strings.TrimSpace(storageEndpoint),
		AccessKey: strings.TrimSpace(storageAccess),
		SecretKey: strings.TrimSpace(storageSecret),
		Bucket:    strings.TrimSpace(storageBucket),
		Region:    strings.TrimSpace(storageRegion),
		PublicURL: strings.TrimSpace(storagePublicURL),
	}
	o.Mail = MailConfig{
		Mode:           MailMode(mailMode),
		From:           strings.TrimSpace(mailFrom),
		APIKey:         strings.TrimSpace(mailAPIKey),
		Domain:         strings.TrimSpace(mailDomain),
		Endpoint:       strings.TrimSpace(mailEndpoint),
		Stream:         strings.TrimSpace(mailStream),
		SMTPHost:       strings.TrimSpace(smtpHost),
		SMTPPort:       strings.TrimSpace(smtpPort),
		SMTPUsername:   strings.TrimSpace(smtpUsername),
		SMTPPassword:   smtpPassword,
		SMTPEncryption: strings.TrimSpace(smtpEncryption),
		SESRegion:      strings.TrimSpace(sesRegion),
		AWSAccessKey:   strings.TrimSpace(awsAccessKey),
		AWSSecretKey:   awsSecretKey,
	}
	o.Plugins = plugins
	o.WithI18n = i18n
	o.StartDocker = contains(actions, "services")
	o.InstallDeps = contains(actions, "deps")
	o.Migrate = contains(actions, "migrate")
	o.Seed = contains(actions, "seed")
	o.Git = git
	return nil
}

// Confirm prints the summary and asks whether to create the project.
func Confirm(o Options) (bool, error) {
	printSummary(o)

	create := true
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create this project?").
				Affirmative("Create").
				Negative("Cancel").
				Value(&create),
		),
	)
	err := ui.RunForm(form)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, nil
		}
		return false, err
	}
	return create, nil
}

// AskStar asks about the GitHub star. It is deliberately a separate question
// after the project works: asking for a star before somebody has seen the
// project is asking for a favour, not an opinion.
func AskStar(o Options) (bool, error) {
	if o.Yes {
		return o.Star, nil
	}
	if o.Star {
		return true, nil
	}
	yes := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Star Grit on GitHub?").
				Description("It helps other people find the project. Either answer is fine.").
				Affirmative("Star it").
				Negative("Not now").
				Value(&yes),
		),
	)
	err := ui.RunForm(form)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, nil
		}
		return false, err
	}
	return yes, nil
}

// AskStart asks whether to run grit start now. The default is yes: the project
// was just installed and migrated, and the next useful thing to do is look at
// it.
func AskStart(o Options) (bool, error) {
	if o.Start != nil {
		return *o.Start, nil
	}
	yes := true
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Run grit start now?").
				Description("Starts the API and every frontend in the foreground. Ctrl+C stops them.").
				Affirmative("Start it").
				Negative("Not yet").
				Value(&yes),
		),
	)
	err := ui.RunForm(form)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, nil
		}
		return false, err
	}
	return yes, nil
}

func (o Options) actionSelection() []string {
	var actions []string
	if o.StartDocker {
		actions = append(actions, "services")
	}
	if o.InstallDeps {
		actions = append(actions, "deps")
	}
	if o.Migrate {
		actions = append(actions, "migrate")
	}
	if o.Seed {
		actions = append(actions, "seed")
	}
	return actions
}

func printSummary(o Options) {
	items := o.Summary()
	width := 0
	for _, item := range items {
		if n := len([]rune(item.Label)); n > width {
			width = n
		}
	}
	fmt.Println()
	fmt.Println(ui.Section("Summary"))
	fmt.Println()
	for _, item := range items {
		label := item.Label + strings.Repeat(" ", width-len([]rune(item.Label)))
		fmt.Println(ui.Answered(label, item.Value, item.Detail))
	}
	fmt.Println()
}

func required(what string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s is required", what)
		}
		return nil
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

// themeOptions is the picker's list. TestThemePickerOffersEveryTheme keeps it
// in step with scaffold.ValidThemes, the way the grit new picker is kept in
// step.
func themeOptions() []huh.Option[string] {
	labels := map[string]string{
		"atlas":   "Atlas — split-screen, blue/white, team/organisation",
		"aurora":  "Aurora — Apple-inspired, monochrome black/white/grey",
		"pulse":   "Pulse — Cloudflare-inspired, premium blue, elevated cards",
		"coral":   "Coral — a card over a blurred glimpse of the app, rose",
		"amber":   "Amber — a plain boxed form under a wordmark, amber",
		"sky":     "Sky — one bold heading, social sign-in first, crisp blue",
		"mono":    "Mono — a fine grid beside a panel of proof, black/white",
		"emerald": "Emerald — a form beside a customer quote, green",
	}
	options := make([]huh.Option[string], 0, len(scaffold.ValidThemes))
	for _, name := range scaffold.ValidThemes {
		label := labels[name]
		if label == "" {
			label = name
		}
		options = append(options, huh.NewOption(label, name))
	}
	return options
}

// databaseOptions is the picker's list, built from the same table the flag
// validates against.
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
		options = append(options, huh.NewOption(label+" — "+scaffold.DBProviders[name], name))
	}
	return options
}

func pluginOptions(o *Options, arch string) []huh.Option[string] {
	probe := *o
	probe.Architecture = arch
	plugins := probe.CompatiblePlugins()
	options := make([]huh.Option[string], 0, len(plugins))
	for _, p := range plugins {
		options = append(options, huh.NewOption(p.Name+" — "+p.Summary, p.Name))
	}
	return options
}
