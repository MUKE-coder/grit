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

import { useTransition } from 'react'
import { useRouter } from 'next/navigation'
import { useLocale, useTranslations } from 'next-intl'

import { LOCALES, isSupported, setLocaleCookie } from '@/lib/locale'

/*
 * The language switcher.
 *
 * A native <select>, on purpose. Neither app ships a dropdown menu primitive,
 * and before v3.222.0 this file imported one anyway, which failed the type
 * check of every project scaffolded with --i18n. A select needs no library and
 * is reachable by keyboard and announced by screen readers without any work.
 *
 * Writes the cookie, then router.refresh() so the server components re-render
 * with the new catalogue. The URL does not change, which is the point of the
 * cookie strategy: the page you are on stays the page you are on.
 */
export function LanguageSwitcher() {
  const router = useRouter()
  const active = useLocale()
  const t = useTranslations('common')
  const [pending, startTransition] = useTransition()

  return (
    <label className="inline-flex items-center">
      <span className="sr-only">{t('language')}</span>
      <select
        value={active}
        disabled={pending}
        onChange={(e) => {
          const code = e.target.value
          if (!isSupported(code) || code === active) return
          setLocaleCookie(code)
          startTransition(() => router.refresh())
        }}
        className="h-9 rounded-lg border border-border bg-bg-tertiary px-2 text-sm text-text-secondary transition-colors hover:bg-bg-hover focus:outline-none focus:ring-2 focus:ring-accent/40 disabled:opacity-60"
      >
        {LOCALES.map((l) => (
          <option key={l.code} value={l.code}>
            {l.label}
          </option>
        ))}
      </select>
    </label>
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
    "systemHub": "System Hub",
    "signOut": "Sign out"
  },
  "table": {
    "search": "Search",
    "filter": "Filter",
    "columns": "Toggle columns",
    "export": "Export",
    "import": "Import",
    "selected": "{count} selected",
    "empty": "No records found",
    "rowsPerPage": "Rows per page",
    "of": "of",
    "previous": "Previous",
    "next": "Next",
    "prev": "Prev",
    "first": "First",
    "last": "Last",
    "actions": "Actions",
    "view": "View",
    "new": "New {name}",
    "showing": "Showing {start}–{end} of {total}",
    "perPage": "{size} / page",
    "emptyHint": "Try adjusting your search or filters",
    "searchPlaceholder": "Search..."
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
    "optional": "Optional",
    "update": "Update",
    "editTitle": "Edit {name}",
    "createTitle": "Create {name}",
    "backTo": "Back to {name}"
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
    "systemHub": "Centre système",
    "signOut": "Se déconnecter"
  },
  "table": {
    "search": "Rechercher",
    "filter": "Filtrer",
    "columns": "Afficher les colonnes",
    "export": "Exporter",
    "import": "Importer",
    "selected": "{count} sélectionné(s)",
    "empty": "Aucun enregistrement",
    "rowsPerPage": "Lignes par page",
    "of": "sur",
    "previous": "Précédent",
    "next": "Suivant",
    "prev": "Préc.",
    "first": "Début",
    "last": "Fin",
    "actions": "Actions",
    "view": "Voir",
    "new": "Ajouter {name}",
    "showing": "{start} à {end} sur {total}",
    "perPage": "{size} par page",
    "emptyHint": "Essayez de modifier votre recherche ou vos filtres",
    "searchPlaceholder": "Rechercher..."
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
    "optional": "Facultatif",
    "update": "Mettre à jour",
    "editTitle": "Modifier {name}",
    "createTitle": "Créer {name}",
    "backTo": "Retour à {name}"
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
    "systemHub": "Kitovu cha mfumo",
    "signOut": "Toka"
  },
  "table": {
    "search": "Tafuta",
    "filter": "Chuja",
    "columns": "Onyesha safu",
    "export": "Hamisha",
    "import": "Ingiza",
    "selected": "{count} zimechaguliwa",
    "empty": "Hakuna rekodi",
    "rowsPerPage": "Safu kwa ukurasa",
    "of": "kati ya",
    "previous": "Iliyotangulia",
    "next": "Ifuatayo",
    "prev": "Nyuma",
    "first": "Mwanzo",
    "last": "Mwisho",
    "actions": "Vitendo",
    "view": "Tazama",
    "new": "Ongeza {name}",
    "showing": "{start}–{end} kati ya {total}",
    "perPage": "{size} kwa ukurasa",
    "emptyHint": "Jaribu kubadilisha utafutaji au vichujio",
    "searchPlaceholder": "Tafuta..."
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
    "optional": "Si lazima",
    "update": "Sasisha",
    "editTitle": "Hariri {name}",
    "createTitle": "Tengeneza {name}",
    "backTo": "Rudi kwa {name}"
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
