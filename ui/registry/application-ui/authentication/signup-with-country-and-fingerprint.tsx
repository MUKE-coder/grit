'use client'

import { useEffect, useId, useMemo, useRef, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Check, ChevronsUpDown, Eye, EyeOff, Search, X } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  useFormField,
} from '@/components/ui/form'

/*
 * Sign-up with a searchable country picker and a signal bundle for
 * duplicate-account detection.
 *
 * ── Colour ──────────────────────────────────────────────────────────────────
 * The orange is orange-700 (#c2410c) wherever white text sits on it, not the
 * brighter orange-500/600 the palette suggests. Measured: white on orange-500
 * is 2.80:1 and on orange-600 is 3.56:1, both under the 4.5:1 WCAG AA needs
 * for a button label. orange-700 is 5.18:1. Links flip to orange-400 in dark
 * mode because orange-700 on a near-black panel drops to 3.43:1 and fails the
 * other way. A bright orange button with white text is the single most common
 * contrast failure in this palette.
 *
 * ── Country picker ──────────────────────────────────────────────────────────
 * APG combobox-with-listbox, not a styled <select>: 250 native options is a
 * scroll wheel. Focus stays on the input and aria-activedescendant carries the
 * highlight, so arrow keys and typing both keep working. The committed ISO
 * code lives in its own field; the visible box holds throwaway filter text.
 * One piece of state for both is how a form ends up submitting "Ger".
 *
 * The list is the full ISO 3166-1 set. A trimmed one is a signup form that
 * cannot be completed by people born in the wrong place.
 *
 * ── Duplicate-account signals ───────────────────────────────────────────────
 * The client collects coarse, non-identifying environment values and hashes
 * them into one opaque id. It deliberately does NOT canvas- or WebGL-
 * fingerprint: those raise entropy enough to track a person across unrelated
 * sites, which is a different and much bigger promise than "is this the same
 * browser that signed up twenty minutes ago".
 *
 * The IP is NOT collected here. A client cannot know its own public address,
 * and anything it claims about one is attacker-controlled. The server reads it
 * from the connection. See the server note below.
 *
 * This is an abuse signal, not an identity. Treat a match as a reason to
 * review, rate-limit or require verification — never to hard-block, because a
 * shared office NAT, a university, a family iPad and a VPN exit all produce
 * collisions between people who have done nothing wrong. Fingerprinting also
 * carries consent obligations under GDPR/ePrivacy in several jurisdictions;
 * check your basis before switching it on.
 */

/* ─── Countries ─────────────────────────────────────────────────────────── */

/** "code:Name" pairs, parsed once. Compact so the full set fits in one file. */
const ISO_3166 =
  'AF:Afghanistan|AX:Åland Islands|AL:Albania|DZ:Algeria|AS:American Samoa|AD:Andorra|AO:Angola|AI:Anguilla|AQ:Antarctica|AG:Antigua and Barbuda|AR:Argentina|AM:Armenia|AW:Aruba|AU:Australia|AT:Austria|AZ:Azerbaijan|BS:Bahamas|BH:Bahrain|BD:Bangladesh|BB:Barbados|BY:Belarus|BE:Belgium|BZ:Belize|BJ:Benin|BM:Bermuda|BT:Bhutan|BO:Bolivia|BQ:Bonaire, Sint Eustatius and Saba|BA:Bosnia and Herzegovina|BW:Botswana|BV:Bouvet Island|BR:Brazil|IO:British Indian Ocean Territory|BN:Brunei Darussalam|BG:Bulgaria|BF:Burkina Faso|BI:Burundi|CV:Cabo Verde|KH:Cambodia|CM:Cameroon|CA:Canada|KY:Cayman Islands|CF:Central African Republic|TD:Chad|CL:Chile|CN:China|CX:Christmas Island|CC:Cocos (Keeling) Islands|CO:Colombia|KM:Comoros|CG:Congo|CD:Congo, Democratic Republic of the|CK:Cook Islands|CR:Costa Rica|CI:Côte d’Ivoire|HR:Croatia|CU:Cuba|CW:Curaçao|CY:Cyprus|CZ:Czechia|DK:Denmark|DJ:Djibouti|DM:Dominica|DO:Dominican Republic|EC:Ecuador|EG:Egypt|SV:El Salvador|GQ:Equatorial Guinea|ER:Eritrea|EE:Estonia|SZ:Eswatini|ET:Ethiopia|FK:Falkland Islands|FO:Faroe Islands|FJ:Fiji|FI:Finland|FR:France|GF:French Guiana|PF:French Polynesia|TF:French Southern Territories|GA:Gabon|GM:Gambia|GE:Georgia|DE:Germany|GH:Ghana|GI:Gibraltar|GR:Greece|GL:Greenland|GD:Grenada|GP:Guadeloupe|GU:Guam|GT:Guatemala|GG:Guernsey|GN:Guinea|GW:Guinea-Bissau|GY:Guyana|HT:Haiti|HM:Heard Island and McDonald Islands|VA:Holy See|HN:Honduras|HK:Hong Kong|HU:Hungary|IS:Iceland|IN:India|ID:Indonesia|IR:Iran|IQ:Iraq|IE:Ireland|IM:Isle of Man|IL:Israel|IT:Italy|JM:Jamaica|JP:Japan|JE:Jersey|JO:Jordan|KZ:Kazakhstan|KE:Kenya|KI:Kiribati|KP:Korea, Democratic People’s Republic of|KR:Korea, Republic of|KW:Kuwait|KG:Kyrgyzstan|LA:Lao People’s Democratic Republic|LV:Latvia|LB:Lebanon|LS:Lesotho|LR:Liberia|LY:Libya|LI:Liechtenstein|LT:Lithuania|LU:Luxembourg|MO:Macao|MG:Madagascar|MW:Malawi|MY:Malaysia|MV:Maldives|ML:Mali|MT:Malta|MH:Marshall Islands|MQ:Martinique|MR:Mauritania|MU:Mauritius|YT:Mayotte|MX:Mexico|FM:Micronesia|MD:Moldova|MC:Monaco|MN:Mongolia|ME:Montenegro|MS:Montserrat|MA:Morocco|MZ:Mozambique|MM:Myanmar|NA:Namibia|NR:Nauru|NP:Nepal|NL:Netherlands|NC:New Caledonia|NZ:New Zealand|NI:Nicaragua|NE:Niger|NG:Nigeria|NU:Niue|NF:Norfolk Island|MK:North Macedonia|MP:Northern Mariana Islands|NO:Norway|OM:Oman|PK:Pakistan|PW:Palau|PS:Palestine, State of|PA:Panama|PG:Papua New Guinea|PY:Paraguay|PE:Peru|PH:Philippines|PN:Pitcairn|PL:Poland|PT:Portugal|PR:Puerto Rico|QA:Qatar|RE:Réunion|RO:Romania|RU:Russian Federation|RW:Rwanda|BL:Saint Barthélemy|SH:Saint Helena, Ascension and Tristan da Cunha|KN:Saint Kitts and Nevis|LC:Saint Lucia|MF:Saint Martin (French part)|PM:Saint Pierre and Miquelon|VC:Saint Vincent and the Grenadines|WS:Samoa|SM:San Marino|ST:Sao Tome and Principe|SA:Saudi Arabia|SN:Senegal|RS:Serbia|SC:Seychelles|SL:Sierra Leone|SG:Singapore|SX:Sint Maarten (Dutch part)|SK:Slovakia|SI:Slovenia|SB:Solomon Islands|SO:Somalia|ZA:South Africa|GS:South Georgia and the South Sandwich Islands|SS:South Sudan|ES:Spain|LK:Sri Lanka|SD:Sudan|SR:Suriname|SJ:Svalbard and Jan Mayen|SE:Sweden|CH:Switzerland|SY:Syrian Arab Republic|TW:Taiwan|TJ:Tajikistan|TZ:Tanzania|TH:Thailand|TL:Timor-Leste|TG:Togo|TK:Tokelau|TO:Tonga|TT:Trinidad and Tobago|TN:Tunisia|TR:Türkiye|TM:Turkmenistan|TC:Turks and Caicos Islands|TV:Tuvalu|UG:Uganda|UA:Ukraine|AE:United Arab Emirates|GB:United Kingdom|US:United States of America|UM:United States Minor Outlying Islands|UY:Uruguay|UZ:Uzbekistan|VU:Vanuatu|VE:Venezuela|VN:Viet Nam|VG:Virgin Islands (British)|VI:Virgin Islands (U.S.)|WF:Wallis and Futuna|EH:Western Sahara|YE:Yemen|ZM:Zambia|ZW:Zimbabwe'

export interface Country {
  code: string
  name: string
}

export const COUNTRIES: Country[] = ISO_3166.split('|').map((entry) => {
  const [code, name] = entry.split(':')
  return { code, name }
})

const COUNTRY_CODES = COUNTRIES.map((country) => country.code)

/* Diacritics folded, so "Cote" finds "Côte d'Ivoire" and "Turkiye" finds
   "Türkiye". Otherwise the hardest names to type are the ones you cannot
   search for. */
function fold(value: string) {
  return value
    .normalize('NFD')
    .replace(/\p{Diacritic}/gu, '')
    .toLowerCase()
}

/* ─── Device signals ────────────────────────────────────────────────────── */

/**
 * Coarse environment values, hashed to one opaque id.
 *
 * Every value here is something a site learns anyway from a normal request or
 * a media query. What is deliberately absent is canvas and WebGL rendering,
 * audio-stack timing and the font list: those push entropy high enough to
 * re-identify a person across unrelated sites, which is surveillance rather
 * than abuse control.
 *
 * Returns null when SubtleCrypto is unavailable, which is any page not served
 * over HTTPS or localhost. Signup must still work there, so callers treat null
 * as "no signal" rather than failing shut.
 */
export async function collectDeviceSignals(): Promise<string | null> {
  if (typeof window === 'undefined' || !window.crypto?.subtle) return null

  const nav = window.navigator as Navigator & { deviceMemory?: number }
  const signals = [
    // Screen shape survives a reload and most window resizing.
    `${window.screen.width}x${window.screen.height}x${window.screen.colorDepth}`,
    Intl.DateTimeFormat().resolvedOptions().timeZone ?? '',
    nav.language ?? '',
    nav.platform ?? '',
    String(nav.hardwareConcurrency ?? ''),
    String(nav.deviceMemory ?? ''),
    String(nav.maxTouchPoints ?? ''),
    // Reduced motion and colour scheme are stable per person, not per visit.
    window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'rm' : '',
    window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light',
  ].join('|')

  const digest = await window.crypto.subtle.digest('SHA-256', new TextEncoder().encode(signals))
  return Array.from(new Uint8Array(digest))
    .map((byte) => byte.toString(16).padStart(2, '0'))
    .join('')
}

/* ─── Schema ────────────────────────────────────────────────────────────── */

const schema = z
  .object({
    name: z.string().min(2, { message: 'Enter your full name' }),
    email: z.string().email({ message: 'Enter a valid email address' }),
    country: z
      .string()
      .min(1, { message: 'Choose your country' })
      .refine((code) => COUNTRY_CODES.includes(code), { message: 'Choose a country from the list' }),
    password: z
      .string()
      .min(8, { message: 'Password must be at least 8 characters' })
      .regex(/[a-z]/, { message: 'Password needs a lowercase letter' })
      .regex(/[A-Z]/, { message: 'Password needs an uppercase letter' })
      .regex(/\d/, { message: 'Password needs a number' }),
    confirmPassword: z.string(),
    /* Absent on http:// and in browsers without SubtleCrypto, so optional.
       A required signal is a signup form that fails shut for some people. */
    deviceId: z.string().optional(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Both passwords must match',
    path: ['confirmPassword'],
  })

export type SignUpValues = z.infer<typeof schema>

const RULES = [
  { id: 'length', label: 'At least 8 characters', test: (v: string) => v.length >= 8 },
  { id: 'lower', label: 'A lowercase letter', test: (v: string) => /[a-z]/.test(v) },
  { id: 'upper', label: 'An uppercase letter', test: (v: string) => /[A-Z]/.test(v) },
  { id: 'number', label: 'A number', test: (v: string) => /\d/.test(v) },
] as const

/* ─── Country combobox ──────────────────────────────────────────────────── */

function CountryCombobox({
  value,
  onChange,
  onBlur,
  invalid,
}: {
  value: string
  onChange: (code: string) => void
  onBlur?: () => void
  invalid?: boolean
}) {
  const [open, setOpen] = useState(false)
  const [filter, setFilter] = useState('')
  const [activeIndex, setActiveIndex] = useState(0)

  const base = useId()
  const listId = `${base}-listbox`
  const optionId = (index: number) => `${base}-option-${index}`

  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLUListElement>(null)
  const wrapperRef = useRef<HTMLDivElement>(null)

  const selected = useMemo(() => COUNTRIES.find((c) => c.code === value), [value])

  const matches = useMemo(() => {
    const needle = fold(filter.trim())
    if (!needle) return COUNTRIES
    // Prefix matches first: "ind" should offer India before the countries that
    // merely contain those letters.
    const starts: Country[] = []
    const contains: Country[] = []
    for (const country of COUNTRIES) {
      const name = fold(country.name)
      if (name.startsWith(needle) || fold(country.code) === needle) starts.push(country)
      else if (name.includes(needle)) contains.push(country)
    }
    return [...starts, ...contains]
  }, [filter])

  /* An activedescendant nobody can see is the same bug as no focus ring. */
  useEffect(() => {
    if (!open) return
    const node = listRef.current?.querySelector(`#${CSS.escape(optionId(activeIndex))}`)
    node?.scrollIntoView({ block: 'nearest' })
  }, [activeIndex, open])

  useEffect(() => {
    if (!open) return
    function onPointerDown(event: PointerEvent) {
      if (!wrapperRef.current?.contains(event.target as Node)) {
        setOpen(false)
        setFilter('')
      }
    }
    document.addEventListener('pointerdown', onPointerDown)
    return () => document.removeEventListener('pointerdown', onPointerDown)
  }, [open])

  function commit(index: number) {
    const country = matches[index]
    if (!country) return
    onChange(country.code)
    setOpen(false)
    setFilter('')
    inputRef.current?.focus()
  }

  function onKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    const last = matches.length - 1

    if (!open && (event.key === 'ArrowDown' || event.key === 'ArrowUp')) {
      event.preventDefault()
      setOpen(true)
      setActiveIndex(0)
      return
    }
    if (!open) return

    if (event.key === 'ArrowDown') {
      event.preventDefault()
      setActiveIndex((c) => (c >= last ? 0 : c + 1))
    } else if (event.key === 'ArrowUp') {
      event.preventDefault()
      setActiveIndex((c) => (c <= 0 ? last : c - 1))
    } else if (event.key === 'Home') {
      event.preventDefault()
      setActiveIndex(0)
    } else if (event.key === 'End') {
      event.preventDefault()
      setActiveIndex(last)
    } else if (event.key === 'Enter') {
      event.preventDefault()
      commit(activeIndex)
    } else if (event.key === 'Escape') {
      // Swallowed only when there is a popup to close, so Escape still reaches
      // a surrounding dialog when there is not.
      event.preventDefault()
      setOpen(false)
      setFilter('')
    } else if (event.key === 'Tab') {
      setOpen(false)
      setFilter('')
    }
  }

  return (
    <div ref={wrapperRef} className="relative">
      <div className="relative">
        <Search
          aria-hidden="true"
          className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-gray-400"
        />
        <input
          ref={inputRef}
          type="text"
          role="combobox"
          aria-expanded={open}
          aria-controls={listId}
          aria-autocomplete="list"
          aria-activedescendant={open && matches.length ? optionId(activeIndex) : undefined}
          aria-invalid={invalid || undefined}
          autoComplete="country-name"
          spellCheck={false}
          placeholder={selected ? selected.name : 'Start typing a country'}
          value={open ? filter : (selected?.name ?? '')}
          onChange={(event) => {
            setFilter(event.target.value)
            setActiveIndex(0)
            if (!open) setOpen(true)
          }}
          onFocus={() => setOpen(true)}
          onBlur={onBlur}
          onKeyDown={onKeyDown}
          className={`min-h-11 w-full rounded-md border bg-white pr-10 pl-9 text-sm text-gray-900 placeholder:text-gray-500 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:bg-gray-900 dark:text-white dark:placeholder:text-gray-400 ${
            invalid ? 'border-red-600 dark:border-red-500' : 'border-gray-300 dark:border-white/15'
          }`}
        />
        <span
          aria-hidden="true"
          className="pointer-events-none absolute inset-y-0 right-0 flex w-10 items-center justify-center text-gray-400"
        >
          <ChevronsUpDown className="size-4" />
        </span>
      </div>

      {open && (
        <ul
          ref={listRef}
          id={listId}
          role="listbox"
          aria-label="Countries"
          className="absolute z-50 mt-1 max-h-64 w-full overflow-y-auto rounded-md border border-gray-200 bg-white py-1 shadow-lg dark:border-white/15 dark:bg-gray-900"
        >
          {matches.length === 0 ? (
            /* Not role="option" — there is nothing here to choose. */
            <li className="px-3 py-2 text-sm text-gray-600 dark:text-gray-300">
              No country matches “{filter}”.
            </li>
          ) : (
            matches.map((country, index) => {
              const active = index === activeIndex
              const chosen = country.code === value
              return (
                <li
                  key={country.code}
                  id={optionId(index)}
                  role="option"
                  aria-selected={chosen}
                  /* pointerdown, not click: click lands after blur, which has
                     already closed the list out from under the pointer. */
                  onPointerDown={(event) => {
                    event.preventDefault()
                    commit(index)
                  }}
                  onMouseMove={() => setActiveIndex(index)}
                  className={`flex cursor-pointer items-center gap-2 px-3 py-2 text-sm ${
                    active
                      ? 'bg-orange-100 text-orange-950 dark:bg-orange-500/20 dark:text-orange-50'
                      : 'text-gray-900 dark:text-white'
                  }`}
                >
                  <span className="w-6 shrink-0 text-xs text-gray-600 tabular-nums dark:text-gray-300">
                    {country.code}
                  </span>
                  <span className="flex-1">{country.name}</span>
                  {chosen && (
                    <Check
                      aria-hidden="true"
                      className="size-4 text-orange-700 dark:text-orange-400"
                    />
                  )}
                </li>
              )
            })
          )}
        </ul>
      )}
    </div>
  )
}

/* ─── Password field ────────────────────────────────────────────────────── */

/* The relative wrapper lives here, not inside FormControl. FormControl is a
   Slot and puts the generated id on its direct child, so wrapping the input in
   a div sends the id to the div and the label points at a container. */
function PasswordField({
  field,
  label,
}: {
  field: React.ComponentProps<typeof Input>
  label: string
}) {
  const [shown, setShown] = useState(false)
  return (
    <div className="relative">
      <FormControl>
        <Input {...field} type={shown ? 'text' : 'password'} autoComplete="new-password" className="h-11 pr-11" />
      </FormControl>
      <button
        type="button"
        onClick={() => setShown((c) => !c)}
        aria-pressed={shown}
        className="absolute inset-y-0 right-0 inline-flex w-11 items-center justify-center text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:text-gray-300 dark:hover:text-white"
      >
        {shown ? <EyeOff aria-hidden="true" className="size-4" /> : <Eye aria-hidden="true" className="size-4" />}
        <span className="sr-only">Show {label}</span>
      </button>
    </div>
  )
}

/* Renders as the field's description so FormControl already links it, and the
   list updates silently rather than reading four rules on every keystroke. */
function PasswordRules({ value }: { value: string }) {
  const { formDescriptionId } = useFormField()
  return (
    <ul id={formDescriptionId} role="list" className="mt-3 grid gap-1.5 sm:grid-cols-2">
      {RULES.map((rule) => {
        const met = rule.test(value ?? '')
        return (
          <li
            key={rule.id}
            className={`flex items-center gap-2 text-sm ${
              met ? 'text-emerald-700 dark:text-emerald-400' : 'text-gray-600 dark:text-gray-300'
            }`}
          >
            {met ? (
              <Check aria-hidden="true" className="size-4 shrink-0" />
            ) : (
              <X aria-hidden="true" className="size-4 shrink-0" />
            )}
            {rule.label}
            <span className="sr-only">{met ? ': met' : ': not met yet'}</span>
          </li>
        )
      })}
    </ul>
  )
}

/* ─── Block ─────────────────────────────────────────────────────────────── */

export default function SignUpWithCountryAndFingerprint({
  onSubmit = async () => {},
  signInHref = '#',
  brand = 'Grit',
}: {
  onSubmit?: (values: SignUpValues) => Promise<void> | void
  signInHref?: string
  brand?: string
}) {
  const [submitting, setSubmitting] = useState(false)

  const form = useForm<SignUpValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      email: '',
      country: '',
      password: '',
      confirmPassword: '',
      deviceId: undefined,
    },
    mode: 'onTouched',
  })

  const password = form.watch('password')

  /* Collected on mount, well before submit, so a slow digest never delays the
     button. Failure is silent by design: no signal beats no signup. */
  useEffect(() => {
    let cancelled = false
    collectDeviceSignals().then((id) => {
      if (!cancelled && id) form.setValue('deviceId', id)
    })
    return () => {
      cancelled = true
    }
  }, [form])

  async function handleSubmit(values: SignUpValues) {
    setSubmitting(true)
    try {
      await onSubmit(values)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-orange-50/40 px-4 py-12 dark:bg-gray-950">
      <div className="w-full max-w-lg">
        <div className="mb-6 flex items-center gap-2.5">
          <span
            aria-hidden="true"
            className="flex size-9 items-center justify-center rounded-lg bg-orange-700 text-sm font-bold text-white"
          >
            {brand.slice(0, 1).toUpperCase()}
          </span>
          <span className="text-lg font-semibold text-gray-900 dark:text-white">{brand}</span>
        </div>

        <div className="rounded-2xl border border-gray-200 bg-white p-8 shadow-sm dark:border-white/10 dark:bg-gray-900">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
            Create your account
          </h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
            Already have one?{' '}
            <a
              href={signInHref}
              className="font-medium text-orange-700 underline-offset-4 hover:underline dark:text-orange-400"
            >
              Sign in
            </a>
            .
          </p>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="mt-8 space-y-5">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Full name</FormLabel>
                    <FormControl>
                      <Input {...field} autoComplete="name" className="h-11" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input {...field} type="email" autoComplete="email" className="h-11" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="country"
                render={({ field, fieldState }) => (
                  <FormItem>
                    <FormLabel>Country</FormLabel>
                    {/* No FormControl wrapper: the combobox owns several
                        elements and the generated id belongs on its input. */}
                    <CountryCombobox
                      value={field.value}
                      onChange={field.onChange}
                      onBlur={field.onBlur}
                      invalid={!!fieldState.error}
                    />
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <PasswordField field={field} label="password" />
                    <PasswordRules value={password ?? ''} />
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Confirm password</FormLabel>
                    <PasswordField field={field} label="confirmation" />
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* orange-700, not the brighter shades: white on orange-600 is
                  3.56:1 and fails AA for a label this size. */}
              <Button
                type="submit"
                disabled={submitting}
                className="h-11 w-full bg-orange-700 font-medium text-white hover:bg-orange-800 focus-visible:outline-orange-600"
              >
                {submitting ? 'Creating account...' : 'Create account'}
              </Button>
            </form>
          </Form>
        </div>

        <p className="mt-4 text-center text-xs text-gray-600 dark:text-gray-400">
          We record your country and a coarse device signature to stop duplicate
          accounts. No canvas or font fingerprinting.
        </p>
      </div>
    </div>
  )
}
