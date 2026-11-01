package setup

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Starter is the example content a new project begins with.
//
// The built-in Blog is the only example Grit writes, so "blank" is not a
// different set of files: it is the same scaffold with the Blog removed.
// Keeping one scaffold behind both starters is what stops the two from
// drifting.
type Starter string

const (
	StarterBlank Starter = "blank"
	StarterBlog  Starter = "blog"
)

// Title is what the picker and the summary call this starter.
func (s Starter) Title() string {
	switch s {
	case StarterBlog:
		return "Blog"
	default:
		return "Blank"
	}
}

// Detail is the one line under the title.
func (s Starter) Detail() string {
	switch s {
	case StarterBlog:
		return "the built-in Blog example, with its admin and public pages"
	default:
		return "a clean project with no example resource"
	}
}

// RedisMode is how the project reaches Redis.
type RedisMode string

const (
	RedisLocal    RedisMode = "local"
	RedisExternal RedisMode = "external"
	RedisOff      RedisMode = "off"
)

// StorageMode is where uploads go.
type StorageMode string

const (
	StorageMinIO StorageMode = "minio"
	StorageS3    StorageMode = "s3"
	StorageR2    StorageMode = "r2"
	StorageLocal StorageMode = "local"
	StorageOff   StorageMode = "off"
)

// MailMode is the mail transport.
type MailMode string

const (
	MailMailhog  MailMode = "mailhog"
	MailLog      MailMode = "log"
	MailResend   MailMode = "resend"
	MailMailgun  MailMode = "mailgun"
	MailPostmark MailMode = "postmark"
	MailSendGrid MailMode = "sendgrid"
	MailSES      MailMode = "ses"
	MailSMTP     MailMode = "smtp"
)

// RedisConfig is the Redis answer.
type RedisConfig struct {
	Mode RedisMode
	URL  string
}

// StorageConfig is the object storage answer. The keys are per provider and
// only the ones the chosen mode needs are written to .env.
type StorageConfig struct {
	Mode      StorageMode
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	PublicURL string
}

// MailConfig is the mail answer. APIKey is the secret for the API-based
// providers; Domain is Mailgun's sending domain.
type MailConfig struct {
	Mode           MailMode
	From           string
	FromName       string
	APIKey         string
	Domain         string
	Endpoint       string
	Stream         string
	SMTPHost       string
	SMTPPort       string
	SMTPUsername   string
	SMTPPassword   string
	SMTPEncryption string
	SESRegion      string
	AWSAccessKey   string
	AWSSecretKey   string
}

// Options is everything grit setup needs. The command fills it from flags and
// the wizard fills what is left.
type Options struct {
	ProjectName string
	InPlace     bool
	Force       bool
	Version     string

	Architecture string // single, double, triple, api, mobile, full, desktop
	Frontend     string // next, tanstack
	Theme        string
	DBProvider   string
	Starter      Starter

	Redis   RedisConfig
	Storage StorageConfig
	Mail    MailConfig

	Plugins  []string
	WithI18n bool

	InstallDeps bool
	StartDocker bool
	Migrate     bool
	Seed        bool

	Git   bool
	Star  bool
	Start *bool

	Yes bool
}

// Normalize fills every default and validates the answers that came from
// flags. The wizard calls it after the last answer, so a value the person
// typed is checked exactly like one a script passed.
func (o *Options) Normalize() error {
	if o.ProjectName == "" && o.InPlace {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}
		o.ProjectName = filepath.Base(cwd)
	}
	if o.ProjectName == "" {
		return fmt.Errorf("give the project a name: grit setup my-app")
	}
	if err := scaffold.ValidateProjectName(o.ProjectName); err != nil {
		return err
	}

	// The same courtesy `grit new` gives: running the command from the folder
	// that already has the project's name scaffolds into it instead of creating
	// a second folder inside it.
	if !o.InPlace {
		if cwd, err := os.Getwd(); err == nil && filepath.Base(cwd) == o.ProjectName {
			o.InPlace = true
		}
	}

	switch o.Architecture {
	case "":
		o.Architecture = "triple"
	case "single", "double", "triple", "api", "mobile", "full", "desktop":
	default:
		return fmt.Errorf("invalid architecture %q: use single, double, triple, desktop, api, mobile or full", o.Architecture)
	}
	if o.Architecture == "full" && o.Frontend == "" {
		o.Frontend = "next"
	}
	if o.Frontend == "" {
		o.Frontend = "next"
	}
	if o.Frontend != "next" && o.Frontend != "tanstack" {
		return fmt.Errorf("invalid frontend %q: use next or tanstack", o.Frontend)
	}

	theme := scaffold.Options{Theme: o.Theme}
	if err := theme.ValidateTheme(); err != nil {
		return err
	}
	o.Theme = theme.Theme

	db := scaffold.Options{DBProvider: o.DBProvider}
	if err := db.ValidateDBProvider(); err != nil {
		return err
	}
	o.DBProvider = db.DBProvider

	if o.Starter == "" {
		o.Starter = StarterBlank
	}
	if !o.SupportsStarters() {
		o.Starter = StarterBlank
	}
	if o.Starter != StarterBlank && o.Starter != StarterBlog {
		return fmt.Errorf("invalid starter %q: use blank or blog", o.Starter)
	}

	if o.Redis.Mode == "" {
		o.Redis.Mode = RedisLocal
	}
	switch o.Redis.Mode {
	case RedisLocal, RedisOff:
	case RedisExternal:
		if strings.TrimSpace(o.Redis.URL) == "" {
			return fmt.Errorf("external Redis needs a URL: --redis external --redis-url rediss://...")
		}
		if !strings.HasPrefix(o.Redis.URL, "redis://") && !strings.HasPrefix(o.Redis.URL, "rediss://") {
			return fmt.Errorf("Redis URL must start with redis:// or rediss://")
		}
	default:
		return fmt.Errorf("invalid Redis choice %q: use local, external or off", o.Redis.Mode)
	}

	if o.Storage.Mode == "" {
		o.Storage.Mode = StorageMinIO
	}
	switch o.Storage.Mode {
	case StorageMinIO, StorageLocal, StorageOff:
	case StorageS3:
		if o.Storage.Bucket == "" {
			return fmt.Errorf("S3 needs a bucket: --storage s3 --storage-bucket my-bucket")
		}
		if o.Storage.Region == "" {
			o.Storage.Region = "us-east-1"
		}
	case StorageR2:
		if o.Storage.Endpoint == "" || o.Storage.AccessKey == "" || o.Storage.SecretKey == "" || o.Storage.Bucket == "" {
			return fmt.Errorf("R2 needs an endpoint, access key, secret key and bucket")
		}
		if o.Storage.Region == "" {
			o.Storage.Region = "auto"
		}
	default:
		return fmt.Errorf("invalid storage choice %q: use minio, s3, r2, local or off", o.Storage.Mode)
	}

	if o.Mail.Mode == "" {
		o.Mail.Mode = MailMailhog
	}
	if o.Mail.From == "" {
		o.Mail.From = "noreply@" + o.ProjectName + ".dev"
	}
	if _, err := mail.ParseAddress(o.Mail.From); err != nil {
		return fmt.Errorf("MAIL_FROM %q is not an email address: %w", o.Mail.From, err)
	}
	switch o.Mail.Mode {
	case MailMailhog, MailLog:
	case MailResend:
		if o.Mail.APIKey == "" {
			return fmt.Errorf("Resend needs an API key: --mail resend --mail-api-key re_...")
		}
	case MailMailgun:
		if o.Mail.Domain == "" || o.Mail.APIKey == "" {
			return fmt.Errorf("Mailgun needs a domain and a secret: --mail mailgun --mail-domain example.com --mail-api-key ...")
		}
		if o.Mail.Endpoint == "" {
			o.Mail.Endpoint = "api.mailgun.net"
		}
	case MailPostmark:
		if o.Mail.APIKey == "" {
			return fmt.Errorf("Postmark needs a server token: --mail postmark --mail-api-key ...")
		}
		if o.Mail.Stream == "" {
			o.Mail.Stream = "outbound"
		}
	case MailSendGrid:
		if o.Mail.APIKey == "" {
			return fmt.Errorf("SendGrid needs an API key: --mail sendgrid --mail-api-key SG....")
		}
	case MailSES:
		if o.Mail.SESRegion == "" {
			o.Mail.SESRegion = "us-east-1"
		}
		if o.Mail.AWSAccessKey == "" || o.Mail.AWSSecretKey == "" {
			return fmt.Errorf("SES needs AWS credentials: --mail ses --aws-access-key ... --aws-secret-key ...")
		}
	case MailSMTP:
		if o.Mail.SMTPHost == "" {
			return fmt.Errorf("SMTP needs a host: --mail smtp --smtp-host smtp.example.com")
		}
		if o.Mail.SMTPPort == "" {
			o.Mail.SMTPPort = "587"
		}
		if o.Mail.SMTPEncryption == "" {
			o.Mail.SMTPEncryption = "starttls"
		}
	default:
		return fmt.Errorf("invalid mail choice %q: use mailhog, log, resend, mailgun, postmark, sendgrid, ses or smtp", o.Mail.Mode)
	}

	o.Plugins = normalizePlugins(o.Plugins)
	for _, name := range o.Plugins {
		if reason := o.PluginIncompatibility(name); reason != "" {
			return fmt.Errorf("plugin %q cannot be installed in this project: %s", name, reason)
		}
	}

	if o.Yes && o.Start == nil {
		start := false
		o.Start = &start
	}
	return nil
}

// ScaffoldOptions turns the setup answers into the scaffold's own options.
func (o Options) ScaffoldOptions() scaffold.Options {
	opts := scaffold.Options{
		ProjectName: o.ProjectName,
		Frontend:    scaffold.Frontend(o.Frontend),
		Theme:       o.Theme,
		DBProvider:  o.DBProvider,
		InPlace:     o.InPlace,
		Force:       o.Force,
		Version:     o.Version,
	}
	switch o.Architecture {
	case "full":
		// Full is a shorthand, not an Architecture: Normalize expands it to
		// triple + web + admin + docs + Expo + desktop.
		opts.Full = true
	case "desktop":
		opts.Architecture = scaffold.ArchTriple
		opts.IncludeDesktop = true
	default:
		opts.Architecture = scaffold.Architecture(o.Architecture)
	}
	return opts
}

// NeedsWebFrontend reports whether this shape has a browser frontend whose
// framework and theme are worth asking about. A mobile project has an Expo app
// and no Next/Vite choice, and an API project has no frontend at all.
func (o Options) NeedsWebFrontend() bool {
	switch o.Architecture {
	case "single", "double", "triple", "full", "desktop":
		return true
	default:
		return false
	}
}

// HasAdminPanel reports whether the Blog example has admin pages to remove or
// keep. It matches scaffold.Options.HasAdminPanel for the shapes setup offers.
func (o Options) HasAdminPanel() bool {
	return o.NeedsWebFrontend()
}

// HasSeparateAdmin reports whether the project has apps/admin, which is what
// the admin-only plugins actually write into. A double has an embedded admin
// and a single has one in its SPA, but those plugins only emit the panel for a
// triple.
func (o Options) HasSeparateAdmin() bool {
	return o.Architecture == "triple" || o.Architecture == "full" || o.Architecture == "desktop"
}

// SupportsStarters reports whether blank and blog are a real choice here. An
// API or mobile project has nowhere for the Blog's pages, so it starts blank.
func (o Options) SupportsStarters() bool {
	return o.HasAdminPanel()
}

// IncludesExpo reports whether the Expo app is part of this shape.
func (o Options) IncludesExpo() bool {
	return o.Architecture == "mobile" || o.Architecture == "full"
}

// IncludesDesktop reports whether the Wails app is part of this shape.
func (o Options) IncludesDesktop() bool {
	return o.Architecture == "desktop" || o.Architecture == "full"
}

// HasNextWeb reports whether the project has the Next.js web app a plugin can
// write a page into.
func (o Options) HasNextWeb() bool {
	if o.Frontend == "tanstack" {
		return false
	}
	switch o.Architecture {
	case "double", "triple", "full", "desktop":
		return true
	default:
		return false
	}
}

// PluginIncompatibility explains why a built-in plugin does not fit, or
// returns "" when it does. It mirrors the checks the plugins themselves make:
// asking for a plugin that would silently install half of itself is worse than
// not offering it.
func (o Options) PluginIncompatibility(name string) string {
	// A new single keeps its Go module in api/ (v3.380.0 and later), while the
	// built-in plugins still resolve a single app's API root to the project
	// root. Until that path handling changes, installing one here would write
	// files the compiler never reads. Excluding single from the picker is the
	// honest answer; changing plugin behaviour is not this command's job.
	if o.Architecture == "single" {
		return "the built-in plugins do not support the api/ layout of a new single app yet"
	}
	switch name {
	case "command-palette", "impersonate", "kanban", "saved-views", "device-pairing":
		if !o.HasSeparateAdmin() {
			return "it needs the separate admin app of a triple, desktop or full project"
		}
	case "push":
		if !o.IncludesExpo() {
			return "it needs the Expo app"
		}
	case "stripe":
		if !o.HasNextWeb() {
			return "it needs the Next.js web app"
		}
	}
	return ""
}

// CompatiblePlugins is the plugin list the picker offers, sorted by name.
func (o Options) CompatiblePlugins() []plugin.Plugin {
	var out []plugin.Plugin
	for _, p := range plugin.All() {
		if o.PluginIncompatibility(p.Name) == "" {
			out = append(out, p)
		}
	}
	return out
}

// LocalServices are the docker compose services this project actually needs.
// Starting the whole compose file starts MySQL's neighbours, MinIO and Mailhog
// for somebody who chose SQLite and external mail, which is both slower and a
// surprise.
func (o Options) LocalServices() []string {
	if !o.StartDocker {
		return nil
	}
	var services []string
	if o.DBProvider == "postgres" {
		services = append(services, "postgres")
	}
	if o.Redis.Mode == RedisLocal {
		services = append(services, "redis")
	}
	if o.Storage.Mode == StorageMinIO {
		services = append(services, "minio")
	}
	if o.Mail.Mode == MailMailhog {
		services = append(services, "mailhog")
	}
	return services
}

// SummaryItem is one row of the confirmation block.
type SummaryItem struct {
	Label  string
	Value  string
	Detail string
}

// Summary says what is about to happen in the words the rest of the command
// uses. It never prints a key or a URL that may contain one.
func (o Options) Summary() []SummaryItem {
	items := []SummaryItem{
		{"Project", o.ProjectName, ""},
		{"Starter", o.Starter.Title(), o.Starter.Detail()},
		{"Shape", architectureTitle(o.Architecture), frontendTitle(o.Frontend)},
		{"Database", o.DBProvider, ""},
		{"Redis", redisTitle(o.Redis), ""},
		{"Storage", storageTitle(o.Storage), ""},
		{"Mail", mailTitle(o.Mail), o.Mail.From},
	}
	if len(o.Plugins) == 0 {
		items = append(items, SummaryItem{"Plugins", "none", ""})
	} else {
		items = append(items, SummaryItem{"Plugins", strings.Join(o.Plugins, ", "), ""})
	}
	if o.WithI18n {
		items = append(items, SummaryItem{"Extras", "internationalisation", ""})
	}
	var actions []string
	if o.StartDocker && len(o.LocalServices()) > 0 {
		actions = append(actions, "local services")
	}
	if o.InstallDeps && o.ScaffoldOptions().ShouldIncludeFrontend() {
		actions = append(actions, "frontend dependencies")
	}
	if o.Migrate {
		actions = append(actions, "migrations")
	}
	if o.Seed {
		actions = append(actions, "sample data")
	}
	if o.Git {
		actions = append(actions, "initial commit")
	}
	if len(actions) == 0 {
		actions = append(actions, "none")
	}
	items = append(items, SummaryItem{"After", strings.Join(actions, ", "), ""})
	return items
}

// StarterRemovalNeeded reports whether the scaffolded project still has the
// Blog model the blank starter removes. A shape without the Blog (or a
// scaffold change that stops writing it) is already blank.
func (o Options) StarterRemovalNeeded(root string) bool {
	if o.Starter != StarterBlank {
		return false
	}
	model := filepath.Join(o.ScaffoldOptions().APIRoot(root), "internal", "models", "blog.go")
	_, err := os.Stat(model)
	return err == nil
}

func architectureTitle(name string) string {
	switch name {
	case "full":
		return "Full"
	case "desktop":
		return "Desktop"
	case "double":
		return "Double"
	case "single":
		return "Single"
	case "api":
		return "API only"
	case "mobile":
		return "Mobile"
	default:
		return "Triple"
	}
}

func frontendTitle(name string) string {
	if name == "tanstack" {
		return "TanStack Router"
	}
	return "Next.js"
}

func redisTitle(c RedisConfig) string {
	switch c.Mode {
	case RedisExternal:
		return "external"
	case RedisOff:
		return "off"
	default:
		return "local, in Docker"
	}
}

func storageTitle(c StorageConfig) string {
	switch c.Mode {
	case StorageS3:
		return "AWS S3"
	case StorageR2:
		return "Cloudflare R2"
	case StorageLocal:
		return "local disk"
	case StorageOff:
		return "off"
	default:
		return "MinIO, in Docker"
	}
}

func mailTitle(c MailConfig) string {
	switch c.Mode {
	case MailLog:
		return "log"
	case MailResend:
		return "Resend"
	case MailMailgun:
		return "Mailgun"
	case MailPostmark:
		return "Postmark"
	case MailSendGrid:
		return "SendGrid"
	case MailSES:
		return "Amazon SES"
	case MailSMTP:
		return "SMTP"
	default:
		return "Mailhog, in Docker"
	}
}

func normalizePlugins(names []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
