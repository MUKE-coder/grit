package setup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// EnvUpdates are the .env lines the service choices own.
//
// Only the active provider's keys are written. A key belonging to a provider
// nobody chose stays empty and is harmless; a key that is written for the
// wrong provider is the kind of thing that makes a project send mail through
// the account somebody was testing last month.
func (o Options) EnvUpdates() map[string]string {
	updates := map[string]string{}

	switch o.Redis.Mode {
	case RedisExternal:
		updates["REDIS_URL"] = o.Redis.URL
	case RedisOff:
		updates["REDIS_URL"] = ""
	default:
		updates["REDIS_URL"] = "redis://localhost:6380"
	}

	switch o.Storage.Mode {
	case StorageS3:
		updates["STORAGE_DRIVER"] = "s3"
		updates["S3_ENDPOINT"] = o.Storage.Endpoint
		updates["S3_ACCESS_KEY"] = o.Storage.AccessKey
		updates["S3_SECRET_KEY"] = o.Storage.SecretKey
		updates["S3_BUCKET"] = o.Storage.Bucket
		updates["S3_REGION"] = o.Storage.Region
		updates["S3_PUBLIC_URL"] = o.Storage.PublicURL
	case StorageR2:
		updates["STORAGE_DRIVER"] = "r2"
		updates["R2_ENDPOINT"] = o.Storage.Endpoint
		updates["R2_ACCESS_KEY"] = o.Storage.AccessKey
		updates["R2_SECRET_KEY"] = o.Storage.SecretKey
		updates["R2_BUCKET"] = o.Storage.Bucket
		updates["R2_REGION"] = o.Storage.Region
		updates["R2_PUBLIC_URL"] = o.Storage.PublicURL
	case StorageLocal:
		updates["STORAGE_DRIVER"] = "local"
	case StorageOff:
		updates["STORAGE_DRIVER"] = ""
	default:
		updates["STORAGE_DRIVER"] = "minio"
	}

	updates["MAIL_FROM"] = o.Mail.From
	if o.Mail.FromName != "" {
		updates["MAIL_FROM_NAME"] = o.Mail.FromName
	}
	switch o.Mail.Mode {
	case MailLog:
		updates["MAIL_MAILER"] = "log"
		// The scaffold puts a placeholder Resend key in .env. Left there, the
		// mailer's "a real key means Resend" rule is one accidental paste away
		// from sending through an account nobody chose.
		updates["RESEND_API_KEY"] = ""
	case MailResend:
		updates["MAIL_MAILER"] = "resend"
		updates["RESEND_API_KEY"] = o.Mail.APIKey
	case MailMailgun:
		updates["MAIL_MAILER"] = "mailgun"
		updates["MAILGUN_DOMAIN"] = o.Mail.Domain
		updates["MAILGUN_SECRET"] = o.Mail.APIKey
		updates["MAILGUN_ENDPOINT"] = o.Mail.Endpoint
	case MailPostmark:
		updates["MAIL_MAILER"] = "postmark"
		updates["POSTMARK_TOKEN"] = o.Mail.APIKey
		updates["POSTMARK_MESSAGE_STREAM"] = o.Mail.Stream
	case MailSendGrid:
		updates["MAIL_MAILER"] = "sendgrid"
		updates["SENDGRID_API_KEY"] = o.Mail.APIKey
	case MailSES:
		updates["MAIL_MAILER"] = "ses"
		updates["AWS_SES_REGION"] = o.Mail.SESRegion
		updates["AWS_ACCESS_KEY_ID"] = o.Mail.AWSAccessKey
		updates["AWS_SECRET_ACCESS_KEY"] = o.Mail.AWSSecretKey
	case MailSMTP:
		updates["MAIL_MAILER"] = "smtp"
		updates["SMTP_HOST"] = o.Mail.SMTPHost
		updates["SMTP_PORT"] = o.Mail.SMTPPort
		updates["SMTP_USERNAME"] = o.Mail.SMTPUsername
		updates["SMTP_PASSWORD"] = o.Mail.SMTPPassword
		updates["SMTP_ENCRYPTION"] = o.Mail.SMTPEncryption
	default:
		updates["MAIL_MAILER"] = ""
		updates["RESEND_API_KEY"] = ""
	}

	return updates
}

// PatchEnv writes the service choices into an existing .env without touching
// anything else in it.
//
// A scaffolded .env holds generated secrets and long explanatory comments. It
// is not a config file to regenerate; it is the one file on the machine that
// must not lose a password because somebody picked Mailhog. Active keys and the
// single-comment form (`# MAIL_MAILER=`) are replaced. Documentation comments
// with indented examples (`#   DATABASE_URL=...`) are left alone: they are the
// explanation, not a value.
func PatchEnv(path string, updates map[string]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	lines := strings.Split(string(data), "\n")

	// Redis used to be off in a way the scaffold explained in two comment
	// lines. Choosing it back on and leaving those lines above the active URL
	// makes the file argue with itself.
	if updates["REDIS_URL"] != "" {
		lines = withoutLines(lines, "No Redis, because --db ", "Delete this line and start one")
	}

	keys := make([]string, 0, len(updates))
	for key := range updates {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := updates[key]
		replacement := key + "=" + value
		replaced := false
		for i, line := range lines {
			if got, ok := envKey(line); ok && got == key {
				lines[i] = replacement
				replaced = true
				break
			}
		}
		if !replaced {
			// A key the scaffold does not know yet, such as STORAGE_DRIVER.
			// Appended rather than inserted into a provider block that would
			// then no longer describe what follows it.
			lines = append(lines, replacement)
		}
	}

	out := strings.Join(lines, "\n")
	// The scaffold writes placeholders. Setup writes the credentials somebody
	// just typed, so the file is private from here on. os.WriteFile keeps the
	// mode of an existing file, which is why the Chmod is not optional.
	if err := os.WriteFile(path, []byte(out), 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("setting permissions on %s: %w", path, err)
	}
	return nil
}

// PatchServiceEnv applies the choices to the project's .env.
func PatchServiceEnv(root string, o Options) error {
	return PatchEnv(filepath.Join(root, ".env"), o.EnvUpdates())
}

// envKey returns the key of an active assignment or of the single-comment form
// the scaffold uses for optional variables.
func envKey(line string) (string, bool) {
	s := strings.TrimSuffix(line, "\r")
	switch {
	case strings.HasPrefix(s, "# "):
		s = s[2:]
		// "#   DATABASE_URL=..." is documentation. "# MAIL_MAILER=" is the
		// commented default the user is expected to uncomment.
		if strings.HasPrefix(s, " ") || strings.HasPrefix(s, "\t") {
			return "", false
		}
	case strings.HasPrefix(s, "#"):
		return "", false
	case strings.HasPrefix(s, " ") || strings.HasPrefix(s, "\t"):
		return "", false
	}

	eq := strings.IndexByte(s, '=')
	if eq <= 0 {
		return "", false
	}
	key := strings.TrimSpace(s[:eq])
	if key == "" {
		return "", false
	}
	for _, r := range key {
		if r == '_' || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			continue
		}
		return "", false
	}
	return key, true
}

// envValue returns the active value of key in an env file, or "" when the file
// or the key is missing. The scaffold writes optional keys in the single-comment
// form (`# KEY=`), so an uncommented value wins and a commented one is the
// fallback.
func envValue(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	var commented string
	for _, line := range strings.Split(string(data), "\n") {
		got, ok := envKey(line)
		if !ok || got != key {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		value := strings.TrimSpace(line[eq+1:])
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			if commented == "" {
				commented = value
			}
			continue
		}
		return value
	}
	return commented
}

func withoutLines(lines []string, starts ...string) []string {
	out := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		matched := false
		for _, start := range starts {
			if strings.Contains(line, start) {
				matched = true
				break
			}
		}
		if !matched {
			out = append(out, line)
			continue
		}
		// The scaffold writes its "no Redis" explanation as two adjacent
		// comment lines. Drop the pair, not just the first, or the sentence
		// that follows loses its subject.
		if i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), "#") {
			i++
		}
	}
	return out
}

// ReadEnv returns one value from a project's .env. It exists for the summary
// and for diagnostics; secrets are never printed by setup itself.
func ReadEnv(path, key string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		k, ok := envKey(scanner.Text())
		if !ok || k != key {
			continue
		}
		line := scanner.Text()
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			return ""
		}
		return strings.TrimSpace(line[eq+1:])
	}
	return ""
}
