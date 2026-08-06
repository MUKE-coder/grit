package scaffold

// next-intl wiring for the generated Next.js apps.
//
// Follows the house pattern from jb.desishub.com/blog/nextjs-i18n-docs:
// cookie-based, no locale prefix in the URL. Routes stay the same in every
// language, so a link shared between two people works for both of them and
// analytics does not fragment by locale.
//
// The one thing added beyond that guide is the handshake with the Go API. The
// cookie name here and middleware.LocaleCookie on the backend have to agree,
// because the API sends its own error messages. Without that, a French admin
// reports English validation failures, which is the exact failure partial i18n
// is known for.

// i18nRequestTS emits i18n/request.ts.
func i18nRequestTS() string {
	return `import { getRequestConfig } from 'next-intl/server'
import { cookies } from 'next/headers'

import { DEFAULT_LOCALE, LOCALE_COOKIE, isSupported } from '@/lib/locale'

/*
 * Resolves the locale for a server render.
 *
 * Cookie only. There is no locale segment in the URL by design, so this is the
 * single source of truth on the server, and the same cookie travels to the Go
 * API on every fetch so both sides agree on the language.
 */
export default getRequestConfig(async () => {
  // Next 15 requires awaiting cookies().
  const store = await cookies()
  const fromCookie = store.get(LOCALE_COOKIE)?.value
  const locale = isSupported(fromCookie) ? fromCookie : DEFAULT_LOCALE

  return {
    locale,
    messages: (await import(` + "`../messages/${locale}.json`" + `)).default,
  }
})
`
}

// i18nLocaleLibTS emits lib/locale.ts — the shared constants both the server
// config and the switcher import, so the cookie name is written once.
func i18nLocaleLibTS() string {
	return `/*
 * Locale constants.
 *
 * LOCALE_COOKIE must match middleware.LocaleCookie in the Go API. The frontend
 * writes it, the browser sends it on every request including API calls, and the
 * backend reads it to translate its own error messages. One name, two runtimes.
 */

export const LOCALE_COOKIE = 'grit_locale'
export const DEFAULT_LOCALE = 'en'

/** Locales the app ships catalogues for. Keep in step with messages/ and with
 *  the API's internal/i18n/locales, or a switcher will offer a language one of
 *  the two cannot serve. */
export const LOCALES = [
  { code: 'en', label: 'English' },
  { code: 'fr', label: 'Français' },
  { code: 'sw', label: 'Kiswahili' },
] as const

export type Locale = (typeof LOCALES)[number]['code']

export function isSupported(value: string | undefined): value is Locale {
  return !!value && LOCALES.some((l) => l.code === value)
}

/** Writes the cookie the server config and the Go API both read. */
export function setLocaleCookie(locale: Locale) {
  // A year, path-wide, and Lax so it survives a normal navigation from an
  // external link. Not HttpOnly: the switcher is client-side and this carries
  // no security meaning.
  document.cookie = ` + "`${LOCALE_COOKIE}=${locale}; path=/; max-age=31536000; SameSite=Lax`" + `
}
`
}

// i18nSwitcherTSX emits components/language-switcher.tsx.
func i18nSwitcherTSX() string {
	return `'use client'

import { useState, useTransition } from 'react'
import { useRouter } from 'next/navigation'
import { useLocale } from 'next-intl'
import { Check, Languages } from 'lucide-react'

import {
  Button,
} from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { LOCALES, setLocaleCookie, type Locale } from '@/lib/locale'

/*
 * The language switcher.
 *
 * Writes the cookie, then router.refresh() so the server components re-render
 * with the new catalogue. The URL does not change, which is the point of the
 * cookie strategy: the page you are on stays the page you are on.
 *
 * useTransition keeps the control responsive while the server round-trips,
 * rather than leaving a menu that looks broken for a beat after the click.
 */
export function LanguageSwitcher() {
  const router = useRouter()
  const active = useLocale()
  const [pending, startTransition] = useTransition()
  const [open, setOpen] = useState(false)

  const choose = (code: Locale) => {
    if (code === active) return setOpen(false)
    setLocaleCookie(code)
    setOpen(false)
    startTransition(() => router.refresh())
  }

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          aria-label="Change language"
          disabled={pending}
        >
          <Languages className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="min-w-[10rem]">
        {LOCALES.map((l) => (
          <DropdownMenuItem
            key={l.code}
            onSelect={() => choose(l.code)}
            className="flex items-center justify-between gap-2"
          >
            {l.label}
            {l.code === active && <Check className="h-4 w-4" aria-hidden="true" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
`
}

// i18nMessagesEN emits messages/en.json for the Next.js apps.
//
// Nested rather than flat, which is what next-intl's useTranslations namespaces
// expect: useTranslations('nav') then t('dashboard').
func i18nMessagesEN() string {
	return `{
  "nav": {
    "dashboard": "Dashboard",
    "resources": "Resources",
    "users": "Users",
    "roles": "Roles",
    "settings": "Settings",
    "system": "System",
    "signOut": "Sign out"
  },
  "table": {
    "search": "Search",
    "filter": "Filter",
    "columns": "Toggle columns",
    "export": "Export",
    "import": "Import",
    "selected": "{count} selected",
    "empty": "Nothing here yet",
    "rowsPerPage": "Rows per page",
    "of": "of",
    "previous": "Previous",
    "next": "Next"
  },
  "form": {
    "save": "Save",
    "saving": "Saving",
    "cancel": "Cancel",
    "delete": "Delete",
    "create": "Create",
    "edit": "Edit",
    "confirmDelete": "Delete this permanently?",
    "required": "Required",
    "optional": "Optional"
  },
  "auth": {
    "signIn": "Sign in",
    "signUp": "Create account",
    "email": "Email",
    "password": "Password",
    "forgot": "Forgot your password?",
    "remember": "Remember this device"
  },
  "common": {
    "loading": "Loading",
    "error": "Something went wrong",
    "retry": "Try again",
    "language": "Language"
  }
}
`
}

func i18nMessagesFR() string {
	return `{
  "nav": {
    "dashboard": "Tableau de bord",
    "resources": "Ressources",
    "users": "Utilisateurs",
    "roles": "Rôles",
    "settings": "Paramètres",
    "system": "Système",
    "signOut": "Se déconnecter"
  },
  "table": {
    "search": "Rechercher",
    "filter": "Filtrer",
    "columns": "Afficher les colonnes",
    "export": "Exporter",
    "import": "Importer",
    "selected": "{count} sélectionné(s)",
    "empty": "Rien ici pour l'instant",
    "rowsPerPage": "Lignes par page",
    "of": "sur",
    "previous": "Précédent",
    "next": "Suivant"
  },
  "form": {
    "save": "Enregistrer",
    "saving": "Enregistrement",
    "cancel": "Annuler",
    "delete": "Supprimer",
    "create": "Créer",
    "edit": "Modifier",
    "confirmDelete": "Supprimer définitivement ?",
    "required": "Obligatoire",
    "optional": "Facultatif"
  },
  "auth": {
    "signIn": "Se connecter",
    "signUp": "Créer un compte",
    "email": "E-mail",
    "password": "Mot de passe",
    "forgot": "Mot de passe oublié ?",
    "remember": "Se souvenir de cet appareil"
  },
  "common": {
    "loading": "Chargement",
    "error": "Une erreur est survenue",
    "retry": "Réessayer",
    "language": "Langue"
  }
}
`
}

func i18nMessagesSW() string {
	return `{
  "nav": {
    "dashboard": "Dashibodi",
    "resources": "Rasilimali",
    "users": "Watumiaji",
    "roles": "Majukumu",
    "settings": "Mipangilio",
    "system": "Mfumo",
    "signOut": "Toka"
  },
  "table": {
    "search": "Tafuta",
    "filter": "Chuja",
    "columns": "Onyesha safu",
    "export": "Hamisha",
    "import": "Ingiza",
    "selected": "{count} zimechaguliwa",
    "empty": "Hakuna kitu bado",
    "rowsPerPage": "Safu kwa ukurasa",
    "of": "kati ya",
    "previous": "Iliyotangulia",
    "next": "Ifuatayo"
  },
  "form": {
    "save": "Hifadhi",
    "saving": "Inahifadhi",
    "cancel": "Ghairi",
    "delete": "Futa",
    "create": "Tengeneza",
    "edit": "Hariri",
    "confirmDelete": "Futa hii kabisa?",
    "required": "Inahitajika",
    "optional": "Si lazima"
  },
  "auth": {
    "signIn": "Ingia",
    "signUp": "Fungua akaunti",
    "email": "Barua pepe",
    "password": "Nenosiri",
    "forgot": "Umesahau nenosiri?",
    "remember": "Kumbuka kifaa hiki"
  },
  "common": {
    "loading": "Inapakia",
    "error": "Hitilafu imetokea",
    "retry": "Jaribu tena",
    "language": "Lugha"
  }
}
`
}
