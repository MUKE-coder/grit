package mail

import (
	"errors"
	"fmt"
	netmail "net/mail"
	"strings"

	"library/apps/api/internal/config"
)

// Drivers lists every value MAIL_MAILER accepts.
var Drivers = []string{"smtp", "resend", "mailgun", "postmark", "sendgrid", "ses", "log", "failover"}

// FromConfig builds the Mailer the configuration asks for.
//
// MAIL_MAILER names the driver. Left empty, a real RESEND_API_KEY means
// Resend, which is how every project sent mail before there were drivers, so
// one that set only that key keeps sending the same way. With neither,
// development sends to Mailhog from docker compose and falls back to the log
// when Mailhog is not running, so a reset link is never lost. Production gets
// no mailer (nil, nil), and the code that sends mail says so.
//
// A driver that is named but not configured is an error rather than a quiet
// fallback: an app that believes it sends through SES should not be writing
// its mail to a log.
func FromConfig(cfg *config.Config) (*Mailer, error) {
	driver := strings.ToLower(strings.TrimSpace(cfg.Mail.Mailer))
	var (
		t   Transport
		err error
	)
	switch {
	case driver != "":
		t, err = transportFor(cfg, driver)
	case resendKeySet(cfg.ResendAPIKey):
		t = NewResend(cfg.ResendAPIKey)
	case cfg.AppEnv == "production":
		return nil, nil
	default:
		t, err = failoverFor(cfg, []string{"smtp", "log"})
	}
	if err != nil {
		return nil, err
	}
	return NewWithTransport(t, fromAddress(cfg)), nil
}

func transportFor(cfg *config.Config, driver string) (Transport, error) {
	mc := cfg.Mail
	switch driver {
	case "resend":
		if !resendKeySet(cfg.ResendAPIKey) {
			return nil, errors.New("MAIL_MAILER=resend needs RESEND_API_KEY")
		}
		return NewResend(cfg.ResendAPIKey), nil
	case "smtp":
		return NewSMTP(mc.SMTPHost, mc.SMTPPort, mc.SMTPUsername, mc.SMTPPassword, mc.SMTPEncryption), nil
	case "mailgun":
		if mc.MailgunDomain == "" || mc.MailgunSecret == "" {
			return nil, errors.New("MAIL_MAILER=mailgun needs MAILGUN_DOMAIN and MAILGUN_SECRET")
		}
		return NewMailgun(mc.MailgunDomain, mc.MailgunSecret, mc.MailgunEndpoint), nil
	case "postmark":
		if mc.PostmarkToken == "" {
			return nil, errors.New("MAIL_MAILER=postmark needs POSTMARK_TOKEN")
		}
		return NewPostmark(mc.PostmarkToken, mc.PostmarkMessageStream), nil
	case "sendgrid":
		if mc.SendGridAPIKey == "" {
			return nil, errors.New("MAIL_MAILER=sendgrid needs SENDGRID_API_KEY")
		}
		return NewSendGrid(mc.SendGridAPIKey), nil
	case "ses":
		if mc.SESRegion == "" || mc.SESAccessKeyID == "" || mc.SESSecretAccessKey == "" {
			return nil, errors.New("MAIL_MAILER=ses needs AWS_SES_REGION, AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
		}
		return NewSES(mc.SESRegion, mc.SESAccessKeyID, mc.SESSecretAccessKey, mc.SESSessionToken), nil
	case "log":
		if cfg.AppEnv == "production" && !mc.AllowLogInProduction {
			return nil, errors.New("MAIL_MAILER=log writes every message, reset links included, to the log instead of sending it, so production refuses it. Set MAIL_ALLOW_LOG_IN_PRODUCTION=true if that is really what you want")
		}
		return NewLog(mc.LogPath), nil
	case "failover":
		if len(mc.Failover) == 0 {
			return nil, errors.New("MAIL_MAILER=failover needs MAIL_FAILOVER, a comma-separated list of drivers such as smtp,log")
		}
		return failoverFor(cfg, mc.Failover)
	default:
		return nil, fmt.Errorf("MAIL_MAILER=%q is not a mail driver; use one of %s", driver, strings.Join(Drivers, ", "))
	}
}

func failoverFor(cfg *config.Config, names []string) (Transport, error) {
	transports := make([]Transport, 0, len(names))
	for _, name := range names {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		if name == "failover" {
			return nil, errors.New("MAIL_FAILOVER cannot name failover itself")
		}
		t, err := transportFor(cfg, name)
		if err != nil {
			return nil, fmt.Errorf("MAIL_FAILOVER: %w", err)
		}
		transports = append(transports, t)
	}
	return NewFailover(transports...), nil
}

// resendKeySet is false for an empty key and for the placeholders the .env
// files ship with.
func resendKeySet(key string) bool {
	return key != "" && !strings.HasPrefix(key, "re_your_api_key")
}

// fromAddress is MAIL_FROM with MAIL_FROM_NAME as its display name.
func fromAddress(cfg *config.Config) string {
	from := strings.TrimSpace(cfg.MailFrom)
	if cfg.Mail.FromName == "" || strings.Contains(from, "<") {
		return from
	}
	return (&netmail.Address{Name: cfg.Mail.FromName, Address: from}).String()
}
