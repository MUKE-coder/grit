package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Mail drivers (issue 004, B1, B2, B4).
//
// The mailer was Resend over HTTP and nothing else, built only when
// RESEND_API_KEY was set. A fresh project therefore sent no mail at all: a
// password reset logged its link, and the Mailhog in docker compose received
// nothing. internal/mail is now a Transport interface with smtp, resend,
// mailgun, postmark, sendgrid, ses, log and failover drivers, picked by
// MAIL_MAILER through mail.FromConfig.
//
// The constants are shared by the templates and by repairMailDrivers, which
// brings an existing project the same text.

const mailerInitOld = `	// Email (Resend)
	var mailer *mail.Mailer
	if cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key" {
		mailer = mail.New(cfg.ResendAPIKey, cfg.MailFrom)
		log.Println("Email service configured")
	} else {
		log.Println("Warning: Resend API key not set (emails disabled)")
	}
`

const mailerInitOldSingle = `	// Email (Resend)
	var mailer *mail.Mailer
	if cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key" {
		mailer = mail.New(cfg.ResendAPIKey, cfg.MailFrom)
		log.Println("Email service configured")
	}
`

const mailerInitNew = `	// Email. MAIL_MAILER picks the transport: smtp, resend, mailgun, postmark,
	// sendgrid, ses, log or failover. Left empty, RESEND_API_KEY means Resend,
	// and development sends to Mailhog, falling back to the log.
	mailer, mailErr := mail.FromConfig(cfg)
	if mailErr != nil {
		log.Fatalf("Email: %v", mailErr)
	}
	if mailer != nil {
		log.Printf("Email: sending with %s", mailer.Driver())
	} else {
		log.Println("Warning: no mail driver configured (emails disabled). Set MAIL_MAILER, or RESEND_API_KEY for Resend.")
	}
`

const configMailFieldsOld = "\tResendAPIKey string\n\tMailFrom     string\n"

const configMailFieldsNew = configMailFieldsOld + `
	// Mail picks the transport (MAIL_MAILER) and holds each driver's
	// settings. internal/mail.FromConfig reads it.
	Mail MailConfig
`

const configMailLoadOld = "\t\tMailFrom:     getEnv(\"MAIL_FROM\", \"noreply@localhost\"),\n"

const configMailLoadNew = configMailLoadOld + "\t\tMail:         loadMailConfig(),\n"

const configMailFuncs = `// MailConfig configures internal/mail. MAIL_MAILER picks the driver, and each
// driver reads only its own keys; .env.example lists them.
type MailConfig struct {
	Mailer               string   // MAIL_MAILER: smtp, resend, mailgun, postmark, sendgrid, ses, log or failover
	FromName             string   // MAIL_FROM_NAME
	Failover             []string // MAIL_FAILOVER, tried in order
	AllowLogInProduction bool     // MAIL_ALLOW_LOG_IN_PRODUCTION
	LogPath              string   // MAIL_LOG_PATH, where the log driver saves messages

	SMTPHost       string
	SMTPPort       string
	SMTPUsername   string
	SMTPPassword   string
	SMTPEncryption string // tls, starttls or none; empty picks from the host and port

	MailgunDomain   string
	MailgunSecret   string
	MailgunEndpoint string // api.mailgun.net, or api.eu.mailgun.net

	PostmarkToken         string
	PostmarkMessageStream string

	SendGridAPIKey string

	SESRegion          string
	SESAccessKeyID     string
	SESSecretAccessKey string
	SESSessionToken    string
}

func loadMailConfig() MailConfig {
	return MailConfig{
		Mailer:               strings.ToLower(getEnv("MAIL_MAILER", "")),
		FromName:             getEnv("MAIL_FROM_NAME", ""),
		Failover:             splitCSV(getEnv("MAIL_FAILOVER", "")),
		AllowLogInProduction: getEnv("MAIL_ALLOW_LOG_IN_PRODUCTION", "false") == "true",
		LogPath:              getEnv("MAIL_LOG_PATH", "storage/mail"),

		SMTPHost: getEnv("SMTP_HOST", "localhost"),
		// Unset, the port follows MAILHOG_SMTP_PORT, so moving Mailhog's port
		// in .env moves the connection with it.
		SMTPPort:       getEnv("SMTP_PORT", getEnv("MAILHOG_SMTP_PORT", "1025")),
		SMTPUsername:   getEnv("SMTP_USERNAME", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SMTPEncryption: getEnv("SMTP_ENCRYPTION", ""),

		MailgunDomain:   getEnv("MAILGUN_DOMAIN", ""),
		MailgunSecret:   getEnv("MAILGUN_SECRET", ""),
		MailgunEndpoint: getEnv("MAILGUN_ENDPOINT", "api.mailgun.net"),

		PostmarkToken:         getEnv("POSTMARK_TOKEN", ""),
		PostmarkMessageStream: getEnv("POSTMARK_MESSAGE_STREAM", "outbound"),

		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),

		SESRegion:          getEnv("AWS_SES_REGION", getEnv("AWS_REGION", "us-east-1")),
		SESAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		SESSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		SESSessionToken:    getEnv("AWS_SESSION_TOKEN", ""),
	}
}
`

const healthDriverFieldOld = "\t\t\tConfigured bool   `json:\"configured,omitempty\"`\n"

const healthDriverFieldNew = healthDriverFieldOld + "\t\t\tDriver     string `json:\"driver,omitempty\"`\n"

const healthMailStatusOld = `		// Email is "configured" when Resend key is set + non-default. The
		// dashboard treats unconfigured as "` + "—" + `" not "down".
		mailStatus := compStatus{
			Configured: cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key",
			OK:         cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key",
		}
`

// healthMailStatusNew keeps the "// Email is" opening: the jobs health repair
// (health_repair.go) ends its block at that line.
const healthMailStatusNew = `		// Email is reported by the driver main.go picked: smtp, resend, log and
		// so on. No mailer is "not configured", which the dashboard shows as a
		// dash rather than as down.
		mailStatus := compStatus{Configured: svc.Mailer != nil, OK: svc.Mailer != nil}
		if svc.Mailer != nil {
			mailStatus.Driver = svc.Mailer.Driver()
		}
`

const envMailHeadOld = "# Email (Resend)\nRESEND_API_KEY=re_your_api_key\n"

const envMailHeadNew = `# Email
# MAIL_MAILER picks how mail is sent: smtp, resend, mailgun, postmark,
# sendgrid, ses, log or failover. Left empty, a real RESEND_API_KEY means
# Resend, and development sends to Mailhog (docker compose, inbox on
# MAILHOG_UI_PORT) and saves to storage/mail when Mailhog is not running.
# Production refuses MAIL_MAILER=log unless MAIL_ALLOW_LOG_IN_PRODUCTION=true.
# MAIL_MAILER=
RESEND_API_KEY=re_your_api_key
`

const envMailDriversNew = `MAIL_FROM_NAME=
# MAIL_FAILOVER=smtp,log

# SMTP (MAIL_MAILER=smtp). Unset, it is Mailhog on localhost:MAILHOG_SMTP_PORT.
# SMTP_ENCRYPTION is tls, starttls or none; empty means tls on port 465, none
# for localhost and starttls for anything else.
# SMTP_HOST=localhost
# SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_ENCRYPTION=

# Mailgun (MAIL_MAILER=mailgun). api.eu.mailgun.net for a domain in the EU.
MAILGUN_DOMAIN=
MAILGUN_SECRET=
MAILGUN_ENDPOINT=api.mailgun.net

# Postmark (MAIL_MAILER=postmark)
POSTMARK_TOKEN=
POSTMARK_MESSAGE_STREAM=outbound

# SendGrid (MAIL_MAILER=sendgrid)
SENDGRID_API_KEY=

# Amazon SES API v2 (MAIL_MAILER=ses). The keys are the usual AWS ones.
AWS_SES_REGION=us-east-1
# AWS_ACCESS_KEY_ID=
# AWS_SECRET_ACCESS_KEY=
`

// envSupportAnchor starts the block after the mail keys in .env.example.
const envSupportAnchor = "\n# Support inbox"

// repairMailDrivers gives an existing project the mail drivers: the MailConfig
// fields in config.go, the internal/mail package, mail.FromConfig in main.go,
// the driver name in /api/health and the keys in .env.example.
//
// Each step waits for the one before it, so a project whose files were edited
// past recognition keeps building on mail.New: the package is delivered only
// once config.go has MailConfig, and main.go switches only once the package
// has FromConfig.
func repairMailDrivers(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	mailDir := filepath.Join(apiRoot, "internal", "mail")
	mailer := filepath.Join(mailDir, "mailer.go")
	configPath := filepath.Join(apiRoot, "internal", "config", "config.go")
	if !fileExists(mailer) || !fileExists(configPath) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	// An edited mailer.go would not be overwritten, and the new drivers do not
	// compile against the old Mailer, so nothing is delivered around it.
	if !fileContains(mailer, "func NewWithTransport(") {
		if key, inside := manifest.Rel(root, mailer); inside && m.StatusOf(root, key) == manifest.Modified {
			fmt.Println("  ⚠ internal/mail/mailer.go has been edited, so the mail drivers were not added. Compare it with a new project's internal/mail to bring them in.")
			return nil
		}
	}

	if err := repairSourceFile(root, m, configPath, repairMailConfigSource); err != nil {
		return err
	}
	if !fileContains(configPath, "type MailConfig struct") {
		return nil
	}
	if err := writeMailFiles(root, opts); err != nil {
		return err
	}
	if !fileContains(mailer, "func NewWithTransport(") || !fileContains(filepath.Join(mailDir, "from_config.go"), "func FromConfig(") {
		return nil
	}

	for _, server := range []string{filepath.Join(apiRoot, "cmd", "server", "main.go"), filepath.Join(root, "main.go")} {
		if !fileExists(server) {
			continue
		}
		if err := repairSourceFile(root, m, server, repairMailerInitSource); err != nil {
			return err
		}
	}
	if routes := filepath.Join(apiRoot, "internal", "routes", "routes.go"); fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairMailHealthSource); err != nil {
			return err
		}
	}
	if env := filepath.Join(root, ".env.example"); fileExists(env) {
		if err := repairTextFile(root, m, env, repairEnvMailSource); err != nil {
			return err
		}
	}
	return nil
}

func repairMailConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "type MailConfig struct") {
		return src, nil, nil
	}
	if strings.Count(src, configMailFieldsOld) != 1 || strings.Count(src, configMailLoadOld) != 1 ||
		!strings.Contains(src, "func splitCSV(") || !strings.Contains(src, "\t\"strings\"\n") {
		return src, nil, []string{"config.go is not the file Grit wrote: add a Mail MailConfig field loaded from MAIL_MAILER and the driver keys, or mail stays on Resend alone"}
	}
	out := strings.Replace(src, configMailFieldsOld, configMailFieldsNew, 1)
	out = strings.Replace(out, configMailLoadOld, configMailLoadNew, 1)
	out = strings.TrimRight(out, "\n") + "\n\n" + configMailFuncs
	return out, []string{"reads MAIL_MAILER and each mail driver's settings"}, nil
}

func repairMailerInitSource(src string) (string, []string, []string) {
	if strings.Contains(src, "mail.FromConfig(") || !strings.Contains(src, "mail.New(") {
		return src, nil, nil
	}
	for _, old := range []string{mailerInitOld, mailerInitOldSingle} {
		if strings.Count(src, old) == 1 {
			return strings.Replace(src, old, mailerInitNew, 1),
				[]string{"the mailer comes from mail.FromConfig, so MAIL_MAILER picks the driver and development mail reaches Mailhog"}, nil
		}
	}
	return src, nil, []string{"the mailer is not built the way Grit wrote it: replace mail.New(cfg.ResendAPIKey, cfg.MailFrom) with mail.FromConfig(cfg) so MAIL_MAILER picks the driver"}
}

func repairMailHealthSource(src string) (string, []string, []string) {
	if strings.Contains(src, "svc.Mailer.Driver()") || !strings.Contains(src, "mailStatus := compStatus{") {
		return src, nil, nil
	}
	if strings.Count(src, healthMailStatusOld) != 1 || strings.Count(src, healthDriverFieldOld) != 1 {
		return src, nil, []string{"the health check's email status is not the one Grit wrote: report svc.Mailer.Driver() rather than whether RESEND_API_KEY is set"}
	}
	out := strings.Replace(src, healthDriverFieldOld, healthDriverFieldNew, 1)
	out = strings.Replace(out, healthMailStatusOld, healthMailStatusNew, 1)
	return out, []string{"/api/health reports the mail driver in use"}, nil
}

func repairEnvMailSource(src string) (string, []string, []string) {
	if strings.Contains(src, "MAIL_MAILER") || !strings.Contains(src, envMailHeadOld) {
		return src, nil, nil
	}
	if strings.Count(src, "\n"+envSupportAnchor) != 1 {
		return src, nil, []string{"the email block is not the one Grit wrote: see a new project's .env.example for MAIL_MAILER and the driver keys"}
	}
	out := strings.Replace(src, envMailHeadOld, envMailHeadNew, 1)
	out = strings.Replace(out, "\n"+envSupportAnchor, "\n"+envMailDriversNew+envSupportAnchor, 1)
	return out, []string{"documents MAIL_MAILER and the SMTP, Mailgun, Postmark, SendGrid and SES keys"}, nil
}
