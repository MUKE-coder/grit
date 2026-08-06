package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// addI18nCmd — `grit add i18n`.
//
// i18n is opt-in rather than scaffolded by default. A project that will only
// ever be English should not carry an unused next-intl dependency, three JSON
// catalogues and a language switcher in its navbar, and a half-translated admin
// is worse than an honestly monolingual one.
//
// `grit new --i18n` runs this same function, so there is one implementation and
// the two cannot drift.
func addI18nCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "i18n",
		Short: "Add internationalisation (next-intl + translated API messages)",
		Long: `Add internationalisation to an existing Grit project.

Installs three things:

  API        internal/i18n with embedded JSON catalogues, locale-resolving
             middleware, and response helpers so a translated error is the
             easy path rather than extra work.
  Next apps  next-intl wired the cookie way, with no locale prefix in the
             URL, so routes stay the same in every language.
  Catalogues English, French and Swahili, translated rather than stubbed.

The frontend cookie and the API both use grit_locale, so a page and the errors
its API calls return agree on the language.

Idempotent: existing files are skipped and every injection checks for itself,
so running it twice changes nothing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting current directory: %w", err)
			}

			res, err := scaffold.AddI18n(cwd, force)
			if err != nil {
				return err
			}

			printI18nResult(res)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite files that already exist")
	return cmd
}

func printI18nResult(res *scaffold.I18nResult) {
	fmt.Println()
	for _, f := range res.Written {
		fmt.Printf("  ✓ %s\n", f)
	}
	for _, f := range res.Skipped {
		fmt.Printf("  • %s (exists, left alone)\n", f)
	}
	for _, w := range res.Wired {
		fmt.Printf("  ✓ wired %s\n", w)
	}

	fmt.Println()
	fmt.Println("  Internationalisation added.")
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println("    1. pnpm install            (picks up next-intl)")
	fmt.Println("    2. Drop <LanguageSwitcher /> into your admin navbar")
	fmt.Println("    3. Use it: const t = useTranslations('nav'); t('dashboard')")
	fmt.Println()
	fmt.Println("  The API translates its own errors. Try:")
	fmt.Println("    curl -H 'Accept-Language: fr' localhost:8080/api/v1/products/nope")
	fmt.Println()
	fmt.Println("  Catalogues live in apps/api/internal/i18n/locales and")
	fmt.Println("  apps/<app>/messages. Add a locale by adding a file to both.")
	fmt.Println()
}
