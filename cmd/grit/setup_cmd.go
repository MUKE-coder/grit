package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/setup"
)

// setupCmd is the installer: a starter, a shape, the services the project
// talks to, and the steps that make the result run.
//
// It is separate from grit new on purpose. grit new is the fast, scriptable
// scaffold and its behaviour is frozen; setup is the guided path that can ask
// questions, change the starter, and run the follow-up commands. Keeping them
// apart means a change here cannot surprise somebody whose build script runs
// `grit new my-app --triple`.
func setupCmd() *cobra.Command {
	var (
		arch, frontend, theme, db, starter string
		inPlace, force, yes                bool

		noInstall, noDocker, noMigrate, noSeed, noGit bool
		star, start, noStart                          bool
		withI18n                                      bool
		plugins                                       []string

		redisMode, redisURL string

		storageMode, storageEndpoint, storageAccessKey string
		storageSecretKey, storageBucket, storageRegion string
		storagePublicURL                               string

		mailMode, mailFrom, mailAPIKey, mailDomain string
		mailEndpoint, mailStream                   string
		smtpHost, smtpPort, smtpUsername           string
		smtpPassword, smtpEncryption               string
		sesRegion, awsAccessKey, awsSecretKey      string
	)

	cmd := &cobra.Command{
		Use:   "setup [project-name|.]",
		Short: "Create a new project with the interactive installer",
		Long: `Create a new project the way an installer does: choose a starter, a shape, a
database, and the services it talks to, then let setup install dependencies,
migrate, and make the first commit.

  grit setup my-app
  grit setup .
  grit setup my-app --yes --arch triple --starter blank --db sqlite --redis off --storage local --mail log --no-docker

Starters live in the project, not in a separate template repository. Blank is
the same scaffold with the built-in Blog removed; Blog keeps it exactly as
grit new writes it.

This command only creates new projects. Use grit upgrade for an existing one.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			projectName := ""
			if len(args) == 1 {
				projectName = strings.TrimSpace(args[0])
			}
			if projectName == "." {
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("getting current directory: %w", err)
				}
				projectName = filepath.Base(cwd)
				inPlace = true
			}

			if start && noStart {
				return fmt.Errorf("--start and --no-start cannot both be set")
			}

			opts := setup.Options{
				ProjectName:  projectName,
				InPlace:      inPlace,
				Force:        force,
				Version:      version,
				Architecture: arch,
				Frontend:     frontend,
				Theme:        theme,
				DBProvider:   db,
				Starter:      setup.Starter(starter),
				Redis: setup.RedisConfig{
					Mode: setup.RedisMode(redisMode),
					URL:  redisURL,
				},
				Storage: setup.StorageConfig{
					Mode:      setup.StorageMode(storageMode),
					Endpoint:  storageEndpoint,
					AccessKey: storageAccessKey,
					SecretKey: storageSecretKey,
					Bucket:    storageBucket,
					Region:    storageRegion,
					PublicURL: storagePublicURL,
				},
				Mail: setup.MailConfig{
					Mode:           setup.MailMode(mailMode),
					From:           mailFrom,
					APIKey:         mailAPIKey,
					Domain:         mailDomain,
					Endpoint:       mailEndpoint,
					Stream:         mailStream,
					SMTPHost:       smtpHost,
					SMTPPort:       smtpPort,
					SMTPUsername:   smtpUsername,
					SMTPPassword:   smtpPassword,
					SMTPEncryption: smtpEncryption,
					SESRegion:      sesRegion,
					AWSAccessKey:   awsAccessKey,
					AWSSecretKey:   awsSecretKey,
				},
				Plugins:     plugins,
				WithI18n:    withI18n,
				InstallDeps: !noInstall,
				StartDocker: !noDocker,
				Migrate:     !noMigrate,
				Seed:        !noSeed,
				Git:         !noGit,
				Star:        star,
				Yes:         yes,
			}
			if start {
				value := true
				opts.Start = &value
			}
			if noStart {
				value := false
				opts.Start = &value
			}

			if err := setup.Run(cmd.Context(), opts); err != nil {
				if errors.Is(err, setup.ErrCancelled) {
					fmt.Println("\n  Cancelled. Nothing was written.")
					return nil
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&arch, "arch", "", "Architecture: single, double, triple, desktop, api, mobile, full")
	cmd.Flags().StringVar(&frontend, "frontend", "", "Frontend: next (default) or tanstack (vite)")
	cmd.Flags().StringVar(&theme, "theme", "", "Theme: atlas (default), aurora, pulse, coral, amber, sky, mono, emerald")
	cmd.Flags().StringVar(&db, "db", "", "Database: postgres (default), mysql, sqlite, memory")
	cmd.Flags().StringVar(&starter, "starter", "", "Starter: blank (default) or blog")

	cmd.Flags().StringVar(&redisMode, "redis", "", "Redis: local (default), external or off")
	cmd.Flags().StringVar(&redisURL, "redis-url", "", "External Redis URL, redis:// or rediss://")

	cmd.Flags().StringVar(&storageMode, "storage", "", "Storage: minio (default), s3, r2, local or off")
	cmd.Flags().StringVar(&storageEndpoint, "storage-endpoint", "", "S3 or R2 endpoint")
	cmd.Flags().StringVar(&storageAccessKey, "storage-access-key", "", "S3 or R2 access key")
	cmd.Flags().StringVar(&storageSecretKey, "storage-secret-key", "", "S3 or R2 secret key")
	cmd.Flags().StringVar(&storageBucket, "storage-bucket", "", "S3 or R2 bucket")
	cmd.Flags().StringVar(&storageRegion, "storage-region", "", "S3 or R2 region")
	cmd.Flags().StringVar(&storagePublicURL, "storage-public-url", "", "Public CDN URL for stored files")

	cmd.Flags().StringVar(&mailMode, "mail", "", "Mail: mailhog (default), log, resend, mailgun, postmark, sendgrid, ses or smtp")
	cmd.Flags().StringVar(&mailFrom, "mail-from", "", "From address, default noreply@<project>.dev")
	cmd.Flags().StringVar(&mailAPIKey, "mail-api-key", "", "API key or secret for the chosen mail provider")
	cmd.Flags().StringVar(&mailDomain, "mail-domain", "", "Mailgun sending domain")
	cmd.Flags().StringVar(&mailEndpoint, "mail-endpoint", "", "Mailgun endpoint, default api.mailgun.net")
	cmd.Flags().StringVar(&mailStream, "mail-stream", "", "Postmark message stream, default outbound")
	cmd.Flags().StringVar(&smtpHost, "smtp-host", "", "SMTP host")
	cmd.Flags().StringVar(&smtpPort, "smtp-port", "", "SMTP port, default 587")
	cmd.Flags().StringVar(&smtpUsername, "smtp-username", "", "SMTP username")
	cmd.Flags().StringVar(&smtpPassword, "smtp-password", "", "SMTP password")
	cmd.Flags().StringVar(&smtpEncryption, "smtp-encryption", "", "SMTP encryption: tls, starttls or none")
	cmd.Flags().StringVar(&sesRegion, "ses-region", "", "Amazon SES region, default us-east-1")
	cmd.Flags().StringVar(&awsAccessKey, "aws-access-key", "", "AWS access key for SES")
	cmd.Flags().StringVar(&awsSecretKey, "aws-secret-key", "", "AWS secret key for SES")

	cmd.Flags().StringSliceVar(&plugins, "plugin", nil, "Built-in plugin to install; repeat or comma-separate")
	cmd.Flags().BoolVar(&withI18n, "i18n", false, "Add internationalisation (next-intl, translated API messages, en/fr/sw)")

	cmd.Flags().BoolVar(&yes, "yes", false, "Use the defaults and the flags, with no prompts. Needs a project name")
	cmd.Flags().BoolVar(&inPlace, "here", false, "Scaffold into the current directory")
	cmd.Flags().BoolVar(&force, "force", false, "Allow scaffolding into a non-empty directory (use with --here)")

	cmd.Flags().BoolVar(&noInstall, "no-install", false, "Skip pnpm install")
	cmd.Flags().BoolVar(&noDocker, "no-docker", false, "Do not start local services with docker compose")
	cmd.Flags().BoolVar(&noMigrate, "no-migrate", false, "Skip database migrations")
	cmd.Flags().BoolVar(&noSeed, "no-seed", false, "Skip sample data")
	cmd.Flags().BoolVar(&noGit, "no-git", false, "Skip git init and the first commit")

	cmd.Flags().BoolVar(&star, "star", false, "Star MUKE-coder/grit on GitHub without asking")
	cmd.Flags().BoolVar(&start, "start", false, "Run grit start after setup")
	cmd.Flags().BoolVar(&noStart, "no-start", false, "Do not run grit start after setup")

	return cmd
}
