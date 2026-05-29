package scaffold

// This file holds the "wishlist" Vite/React components that every Grit
// single-app project ends up writing by hand: a real auth helper, a
// confirm dialog with a Promise-based hook, MoneyInput, Combobox,
// SessionExpiryMonitor, StatusBadge, StatsRow. They are wired into
// writeSingleFrontendFiles so a fresh `grit new --single --vite` ships
// with them already in place.
//
// Conventions:
//   - Zero external deps beyond what's already in the single-app
//     package.json (react 19, axios, clsx, tailwind-merge, lucide-react).
//   - Tailwind classes use the Grit design tokens
//     (bg-background, border-border, text-foreground, etc.).
//   - All identifiers are typed; no `any`.

// singleViteNavbar emits components/navbar.tsx for the Vite scaffold —
// uses TanStack Router's <Link> + useRouterState() in place of next/link
// and usePathname(). The Next.js variant lives in web_files.go::webNavbar
// and gets used by the monorepo --triple scaffold.
func singleViteNavbar(opts Options) string {
	return `import { useState } from "react"
import { Link, useRouterState } from "@tanstack/react-router"
import { Menu, X, Github } from "lucide-react"

const DOCS_URL = "https://grit-vert.vercel.app/docs"

const navLinks = [
  { href: "/", label: "Home" },
  { href: "/blog", label: "Blog" },
]

export function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false)
  const pathname = useRouterState({ select: (s) => s.location.pathname })

  return (
    <nav className="sticky top-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-lg">
      <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-6">
        <Link to="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent/15 border border-accent/20">
            <span className="text-accent font-mono font-bold text-sm">G</span>
          </div>
          <span className="text-lg font-bold tracking-tight">` + opts.ProjectName + `</span>
        </Link>

        <div className="hidden md:flex items-center gap-6">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              to={link.href}
              className={
                "text-sm transition-colors " +
                (pathname === link.href
                  ? "text-foreground font-medium"
                  : "text-text-secondary hover:text-foreground")
              }
            >
              {link.label}
            </Link>
          ))}
          <a
            href={DOCS_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-text-secondary hover:text-foreground transition-colors"
          >
            Docs
          </a>
          <a
            href="https://github.com/MUKE-coder/grit"
            target="_blank"
            rel="noopener noreferrer"
            className="text-text-secondary hover:text-foreground transition-colors"
          >
            <Github className="h-5 w-5" />
          </a>
        </div>

        <button
          onClick={() => setMobileOpen(!mobileOpen)}
          className="md:hidden p-2 text-text-secondary hover:text-foreground transition-colors"
          aria-label="Toggle menu"
        >
          {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>
      </div>

      {mobileOpen ? (
        <div className="md:hidden border-t border-border/50 bg-background/95 backdrop-blur-lg">
          <div className="mx-auto max-w-5xl px-6 py-4 flex flex-col gap-3">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                to={link.href}
                onClick={() => setMobileOpen(false)}
                className={
                  "text-sm py-2 transition-colors " +
                  (pathname === link.href
                    ? "text-foreground font-medium"
                    : "text-text-secondary hover:text-foreground")
                }
              >
                {link.label}
              </Link>
            ))}
            <a
              href={DOCS_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm py-2 text-text-secondary hover:text-foreground transition-colors"
            >
              Docs
            </a>
            <a
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm py-2 text-text-secondary hover:text-foreground transition-colors"
            >
              GitHub
            </a>
          </div>
        </div>
      ) : null}
    </nav>
  )
}
`
}

// singleViteFooter emits components/footer.tsx for the Vite scaffold.
func singleViteFooter(opts Options) string {
	return `const DOCS_URL = "https://grit-vert.vercel.app/docs"

export function Footer() {
  return (
    <footer className="border-t border-border/50 bg-background">
      <div className="mx-auto max-w-5xl px-6 py-8">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-text-muted">
            <span className="font-semibold text-text-secondary">` + opts.ProjectName + `</span>
            <span className="text-border">·</span>
            <span>Built with Grit</span>
          </div>
          <div className="flex items-center gap-6 text-sm text-text-muted">
            <a
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              GitHub
            </a>
            <a
              href={DOCS_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              Documentation
            </a>
          </div>
        </div>
        <div className="mt-6 pt-6 border-t border-border/30 text-center">
          <p className="text-xs text-text-muted">
            &copy; {new Date().getFullYear()} ` + opts.ProjectName + `. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}
`
}

// singleViteEnvTypes emits frontend/src/vite-env.d.ts. Without this,
// `import.meta.env.VITE_API_URL` errors under `tsc --noEmit` with
// "Property 'env' does not exist on type 'ImportMeta'". The triple-slash
// reference pulls in Vite's client types in addition to declaring our
// custom VITE_API_URL.
func singleViteEnvTypes() string {
	return `/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
`
}

// singleAuthLib emits frontend/src/lib/auth.ts — the missing helper layer
// the user spent a debug cycle re-deriving from scratch. Reads the
// {data:{user,tokens:{access_token,refresh_token,expires_at}}} envelope
// the Go API actually returns; persists tokens to localStorage;
// transparently swaps in the refresh token when the access token expires.
func singleAuthLib() string {
	return `import { api } from "./api"

export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  role: string
  avatar?: string
  job_title?: string
  bio?: string
  active: boolean
  created_at: string
  updated_at: string
}

export interface Tokens {
  access_token: string
  refresh_token: string
  expires_at: number
}

export interface LoginResult {
  user: User
  tokens: Tokens
}

interface Envelope<T> {
  data: T
  message?: string
}

interface TOTPChallenge {
  totp_required: true
  pending_token: string
}

const ACCESS_KEY = "grit.access_token"
const REFRESH_KEY = "grit.refresh_token"
const EXPIRES_KEY = "grit.token_expires_at"

export function getAccessToken(): string | null {
  return typeof window === "undefined" ? null : localStorage.getItem(ACCESS_KEY)
}

export function getRefreshToken(): string | null {
  return typeof window === "undefined" ? null : localStorage.getItem(REFRESH_KEY)
}

export function getTokenExpiresAt(): number | null {
  const raw = typeof window === "undefined" ? null : localStorage.getItem(EXPIRES_KEY)
  if (!raw) return null
  const parsed = Number(raw)
  return Number.isFinite(parsed) ? parsed : null
}

export function isAuthenticated(): boolean {
  return Boolean(getAccessToken())
}

function persistTokens(tokens: Tokens) {
  localStorage.setItem(ACCESS_KEY, tokens.access_token)
  localStorage.setItem(REFRESH_KEY, tokens.refresh_token)
  localStorage.setItem(EXPIRES_KEY, String(tokens.expires_at))
}

export function clearAuth() {
  localStorage.removeItem(ACCESS_KEY)
  localStorage.removeItem(REFRESH_KEY)
  localStorage.removeItem(EXPIRES_KEY)
}

/**
 * login() handles the {data:{user,tokens}} envelope that the Go API returns.
 * If the user has TOTP enabled the server replies with
 * {data:{totp_required:true, pending_token}} — surfaced via the totp field
 * on the result so callers can route to a verification screen.
 */
export async function login(email: string, password: string): Promise<LoginResult | { totp: TOTPChallenge }> {
  const res = await api.post<Envelope<LoginResult | TOTPChallenge>>("/api/auth/login", { email, password })
  const data = res.data.data
  if ("totp_required" in data) {
    return { totp: data }
  }
  persistTokens(data.tokens)
  return data
}

export async function register(payload: {
  first_name: string
  last_name: string
  email: string
  password: string
}): Promise<LoginResult> {
  const res = await api.post<Envelope<LoginResult>>("/api/auth/register", payload)
  persistTokens(res.data.data.tokens)
  return res.data.data
}

/** Fetches the current user via /api/auth/me. Throws if unauthenticated. */
export async function me(): Promise<User> {
  const res = await api.get<Envelope<User>>("/api/auth/me")
  return res.data.data
}

/**
 * refresh() exchanges the stored refresh token for a fresh access+refresh
 * pair. Clears auth on failure so the next request bounces to login
 * instead of looping forever on a dead token.
 */
export async function refresh(): Promise<Tokens | null> {
  const refreshToken = getRefreshToken()
  if (!refreshToken) return null
  try {
    const res = await api.post<Envelope<{ tokens: Tokens }>>("/api/auth/refresh", {
      refresh_token: refreshToken,
    })
    persistTokens(res.data.data.tokens)
    return res.data.data.tokens
  } catch {
    clearAuth()
    return null
  }
}

export async function logout(): Promise<void> {
  try {
    await api.post("/api/auth/logout")
  } catch {
    // Server-side logout is best-effort — we always clear local state.
  }
  clearAuth()
}
`
}

// viteAPIClientWithAuth replaces the bare viteAPIClient() that ships with
// the scaffold. Adds: Authorization header from getAccessToken(), a
// response interceptor that tries one transparent refresh on 401, and the
// idempotency-key behaviour preserved from the original.
func viteAPIClientWithAuth() string {
	return `import axios, { AxiosError, type AxiosRequestConfig } from "axios"

const API_URL = (import.meta.env.VITE_API_URL as string | undefined) ?? ""

const ACCESS_KEY = "grit.access_token"
const REFRESH_KEY = "grit.refresh_token"
const EXPIRES_KEY = "grit.token_expires_at"

export const api = axios.create({
  baseURL: API_URL,
  headers: { "Content-Type": "application/json" },
})

// Attach Authorization + Idempotency-Key on every outgoing request.
api.interceptors.request.use((config) => {
  const token = typeof window === "undefined" ? null : localStorage.getItem(ACCESS_KEY)
  if (token && config.headers) {
    config.headers.Authorization = "Bearer " + token
  }
  const method = (config.method || "get").toUpperCase()
  const unsafe = method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE"
  if (unsafe && config.headers && !config.headers["Idempotency-Key"]) {
    config.headers["Idempotency-Key"] = crypto.randomUUID()
  }
  return config
})

// Single-flight refresh: if many requests 401 at once, only one of them
// actually calls /api/auth/refresh; the rest await the same promise.
let inflightRefresh: Promise<string | null> | null = null

async function refreshAccessToken(): Promise<string | null> {
  if (inflightRefresh) return inflightRefresh
  inflightRefresh = (async () => {
    const rt = typeof window === "undefined" ? null : localStorage.getItem(REFRESH_KEY)
    if (!rt) return null
    try {
      const res = await axios.post(API_URL + "/api/auth/refresh", { refresh_token: rt })
      const tokens = res.data?.data?.tokens as
        | { access_token: string; refresh_token: string; expires_at: number }
        | undefined
      if (!tokens) return null
      localStorage.setItem(ACCESS_KEY, tokens.access_token)
      localStorage.setItem(REFRESH_KEY, tokens.refresh_token)
      localStorage.setItem(EXPIRES_KEY, String(tokens.expires_at))
      return tokens.access_token
    } catch {
      localStorage.removeItem(ACCESS_KEY)
      localStorage.removeItem(REFRESH_KEY)
      localStorage.removeItem(EXPIRES_KEY)
      return null
    } finally {
      inflightRefresh = null
    }
  })()
  return inflightRefresh
}

// Retry once on 401 after a refresh. Skip the refresh endpoint itself to
// avoid recursion if the refresh token is the one that died.
api.interceptors.response.use(
  (res) => res,
  async (err: AxiosError) => {
    const original = err.config as (AxiosRequestConfig & { _retried?: boolean }) | undefined
    const status = err.response?.status
    const url = original?.url ?? ""
    if (status !== 401 || !original || original._retried || url.includes("/api/auth/refresh")) {
      return Promise.reject(err)
    }
    const newToken = await refreshAccessToken()
    if (!newToken) return Promise.reject(err)
    original._retried = true
    original.headers = { ...(original.headers || {}), Authorization: "Bearer " + newToken }
    return api.request(original)
  },
)
`
}

// singleUseConfirmHook emits hooks/use-confirm.tsx — Promise-based
// confirmation dialog usable as `const ok = await confirm({...})`.
// Mount <ConfirmProvider> once at the app root.
func singleUseConfirmHook() string {
	return `import { createContext, useCallback, useContext, useEffect, useRef, useState } from "react"
import type { ReactNode } from "react"
import { cn } from "@/lib/utils"

export type ConfirmTone = "default" | "danger"

export interface ConfirmOptions {
  title: string
  message?: ReactNode
  confirmLabel?: string
  cancelLabel?: string
  tone?: ConfirmTone
}

type Resolver = (ok: boolean) => void

interface ConfirmContextValue {
  request: (opts: ConfirmOptions) => Promise<boolean>
}

const ConfirmContext = createContext<ConfirmContextValue | null>(null)

export function useConfirm() {
  const ctx = useContext(ConfirmContext)
  if (!ctx) {
    throw new Error("useConfirm must be used inside <ConfirmProvider>")
  }
  return ctx.request
}

interface PendingConfirm extends ConfirmOptions {
  resolve: Resolver
}

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [pending, setPending] = useState<PendingConfirm | null>(null)
  const confirmBtnRef = useRef<HTMLButtonElement | null>(null)

  const request = useCallback((opts: ConfirmOptions) => {
    return new Promise<boolean>((resolve) => {
      setPending({ ...opts, resolve })
    })
  }, [])

  const close = useCallback((ok: boolean) => {
    if (!pending) return
    pending.resolve(ok)
    setPending(null)
  }, [pending])

  useEffect(() => {
    if (!pending) return
    confirmBtnRef.current?.focus()
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") close(false)
      if (e.key === "Enter") close(true)
    }
    window.addEventListener("keydown", onKey)
    return () => window.removeEventListener("keydown", onKey)
  }, [pending, close])

  return (
    <ConfirmContext.Provider value={{ request }}>
      {children}
      {pending ? (
        <div
          role="dialog"
          aria-modal="true"
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
          onClick={() => close(false)}
        >
          <div
            className="relative w-full max-w-md mx-4 rounded-lg border border-border bg-bg-elevated p-6 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-lg font-semibold text-foreground">{pending.title}</h2>
            {pending.message ? (
              <div className="mt-2 text-sm text-text-secondary">{pending.message}</div>
            ) : null}
            <div className="mt-6 flex justify-end gap-2">
              <button
                type="button"
                onClick={() => close(false)}
                className="px-4 py-2 rounded-md border border-border text-foreground hover:bg-bg-hover transition-colors text-sm"
              >
                {pending.cancelLabel || "Cancel"}
              </button>
              <button
                type="button"
                ref={confirmBtnRef}
                onClick={() => close(true)}
                className={cn(
                  "px-4 py-2 rounded-md text-sm font-medium transition-colors",
                  pending.tone === "danger"
                    ? "bg-danger text-white hover:opacity-90"
                    : "bg-accent text-white hover:bg-accent-hover",
                )}
              >
                {pending.confirmLabel || "Confirm"}
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </ConfirmContext.Provider>
  )
}
`
}

// singleMoneyInput emits components/money-input.tsx — text input that
// formats with thousands separators while still exposing the underlying
// number via onValueChange. Empty input maps to null so the caller can
// distinguish "no value" from 0.
func singleMoneyInput() string {
	return `import { forwardRef, useMemo } from "react"
import type { InputHTMLAttributes } from "react"
import { cn } from "@/lib/utils"

export interface MoneyInputProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, "onChange" | "value" | "type"> {
  value: number | null | undefined
  onValueChange: (next: number | null) => void
  prefix?: string
  /** Locale for thousands separators. Defaults to the browser default. */
  locale?: string
}

const NUMERIC = /[^0-9.-]+/g

function formatDisplay(value: number | null | undefined, locale?: string) {
  if (value === null || value === undefined || Number.isNaN(value)) return ""
  return new Intl.NumberFormat(locale, { maximumFractionDigits: 2 }).format(value)
}

export const MoneyInput = forwardRef<HTMLInputElement, MoneyInputProps>(function MoneyInput(
  { value, onValueChange, prefix, locale, className, ...rest },
  ref,
) {
  const display = useMemo(() => formatDisplay(value, locale), [value, locale])

  return (
    <div className={cn("relative flex items-center", className)}>
      {prefix ? (
        <span className="absolute left-3 text-text-muted text-sm pointer-events-none select-none">
          {prefix}
        </span>
      ) : null}
      <input
        ref={ref}
        type="text"
        inputMode="decimal"
        autoComplete="off"
        value={display}
        onChange={(e) => {
          const raw = e.target.value.replace(NUMERIC, "")
          if (raw === "" || raw === "-") {
            onValueChange(null)
            return
          }
          const parsed = Number(raw)
          if (Number.isNaN(parsed)) return
          onValueChange(parsed)
        }}
        className={cn(
          "w-full rounded-md border border-border bg-bg-elevated px-3 py-2 text-foreground tabular-nums",
          "focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent",
          prefix ? "pl-10" : "",
        )}
        {...rest}
      />
    </div>
  )
})
`
}

// singleCombobox emits components/combobox.tsx — keyboard-friendly
// searchable select. Filter matches label OR optional sublabel.
func singleCombobox() string {
	return `import { useEffect, useId, useMemo, useRef, useState } from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"

export interface ComboboxOption {
  value: string
  label: string
  sublabel?: string
  disabled?: boolean
}

export interface ComboboxProps {
  options: ComboboxOption[]
  value: string | null
  onChange: (value: string | null) => void
  placeholder?: string
  emptyMessage?: string
  disabled?: boolean
  className?: string
}

export function Combobox({
  options,
  value,
  onChange,
  placeholder = "Select an option…",
  emptyMessage = "No matches",
  disabled,
  className,
}: ComboboxProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState("")
  const [activeIndex, setActiveIndex] = useState(0)
  const wrapperRef = useRef<HTMLDivElement | null>(null)
  const inputRef = useRef<HTMLInputElement | null>(null)
  const listboxId = useId()

  const selected = useMemo(
    () => options.find((o) => o.value === value) || null,
    [options, value],
  )

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return options
    return options.filter(
      (o) =>
        o.label.toLowerCase().includes(q) ||
        (o.sublabel ? o.sublabel.toLowerCase().includes(q) : false),
    )
  }, [options, query])

  useEffect(() => {
    if (!open) return
    function onClick(e: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener("mousedown", onClick)
    return () => document.removeEventListener("mousedown", onClick)
  }, [open])

  useEffect(() => {
    if (open) setActiveIndex(0)
  }, [open, query])

  function commit(opt: ComboboxOption | undefined) {
    if (!opt || opt.disabled) return
    onChange(opt.value === value ? null : opt.value)
    setOpen(false)
    setQuery("")
  }

  function handleKey(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "ArrowDown") {
      e.preventDefault()
      setActiveIndex((i) => Math.min(i + 1, filtered.length - 1))
    } else if (e.key === "ArrowUp") {
      e.preventDefault()
      setActiveIndex((i) => Math.max(i - 1, 0))
    } else if (e.key === "Enter") {
      e.preventDefault()
      commit(filtered[activeIndex])
    } else if (e.key === "Escape") {
      setOpen(false)
    }
  }

  return (
    <div ref={wrapperRef} className={cn("relative", className)}>
      <button
        type="button"
        disabled={disabled}
        onClick={() => {
          setOpen((o) => !o)
          setTimeout(() => inputRef.current?.focus(), 0)
        }}
        className={cn(
          "w-full flex items-center justify-between gap-2 rounded-md border border-border bg-bg-elevated px-3 py-2 text-sm text-foreground",
          "hover:bg-bg-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed",
          "focus:outline-none focus:ring-2 focus:ring-accent/50",
        )}
        aria-haspopup="listbox"
        aria-expanded={open}
      >
        <span className={cn("truncate text-left", !selected && "text-text-muted")}>
          {selected ? selected.label : placeholder}
        </span>
        <ChevronsUpDown className="h-4 w-4 text-text-muted shrink-0" />
      </button>

      {open ? (
        <div className="absolute z-50 mt-1 w-full rounded-md border border-border bg-bg-elevated shadow-lg">
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKey}
            placeholder="Search…"
            className="w-full bg-transparent border-b border-border px-3 py-2 text-sm text-foreground focus:outline-none"
          />
          <ul
            role="listbox"
            id={listboxId}
            className="max-h-64 overflow-y-auto py-1"
          >
            {filtered.length === 0 ? (
              <li className="px-3 py-2 text-sm text-text-muted">{emptyMessage}</li>
            ) : (
              filtered.map((opt, i) => {
                const isActive = i === activeIndex
                const isSelected = opt.value === value
                return (
                  <li
                    key={opt.value}
                    role="option"
                    aria-selected={isSelected}
                    onMouseEnter={() => setActiveIndex(i)}
                    onMouseDown={(e) => {
                      e.preventDefault()
                      commit(opt)
                    }}
                    className={cn(
                      "px-3 py-2 cursor-pointer flex items-start gap-2",
                      isActive ? "bg-bg-hover" : "",
                      opt.disabled ? "opacity-50 cursor-not-allowed" : "",
                    )}
                  >
                    <Check
                      className={cn(
                        "h-4 w-4 mt-0.5 shrink-0",
                        isSelected ? "text-accent" : "opacity-0",
                      )}
                    />
                    <div className="min-w-0 flex-1">
                      <div className="text-sm text-foreground truncate">{opt.label}</div>
                      {opt.sublabel ? (
                        <div className="text-xs text-text-muted truncate">{opt.sublabel}</div>
                      ) : null}
                    </div>
                  </li>
                )
              })
            )}
          </ul>
        </div>
      ) : null}
    </div>
  )
}
`
}

// singleSessionExpiryMonitor emits components/session-expiry-monitor.tsx.
// Decodes the JWT exp claim from localStorage, shows a modal 30s before
// expiry with Stay/Logout, and calls refresh() on Stay. Counter or refresh
// failure triggers a clean logout. Mount near the app root.
func singleSessionExpiryMonitor() string {
	return `import { useEffect, useRef, useState } from "react"
import { clearAuth, getAccessToken, refresh } from "@/lib/auth"

const WARNING_SECONDS = 30
const POLL_MS = 1000

function decodeExp(token: string): number | null {
  try {
    const payload = JSON.parse(atob(token.split(".")[1])) as { exp?: number }
    return typeof payload.exp === "number" ? payload.exp : null
  } catch {
    return null
  }
}

function secondsLeft(token: string | null): number | null {
  if (!token) return null
  const exp = decodeExp(token)
  if (!exp) return null
  return Math.floor(exp - Date.now() / 1000)
}

export interface SessionExpiryMonitorProps {
  /** Called when the session is no longer valid (refresh failed or user logged out). */
  onLoggedOut?: () => void
}

export function SessionExpiryMonitor({ onLoggedOut }: SessionExpiryMonitorProps) {
  const [showWarning, setShowWarning] = useState(false)
  const [countdown, setCountdown] = useState<number>(WARNING_SECONDS)
  const refreshingRef = useRef(false)

  useEffect(() => {
    const interval = window.setInterval(() => {
      const left = secondsLeft(getAccessToken())
      if (left === null) {
        setShowWarning(false)
        return
      }
      if (left <= 0) {
        setShowWarning(false)
        clearAuth()
        onLoggedOut?.()
        return
      }
      if (left <= WARNING_SECONDS) {
        setShowWarning(true)
        setCountdown(left)
      } else {
        setShowWarning(false)
      }
    }, POLL_MS)
    return () => window.clearInterval(interval)
  }, [onLoggedOut])

  async function stay() {
    if (refreshingRef.current) return
    refreshingRef.current = true
    const tokens = await refresh()
    refreshingRef.current = false
    if (!tokens) {
      setShowWarning(false)
      onLoggedOut?.()
      return
    }
    setShowWarning(false)
  }

  function logout() {
    clearAuth()
    setShowWarning(false)
    onLoggedOut?.()
  }

  if (!showWarning) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="w-full max-w-md mx-4 rounded-lg border border-border bg-bg-elevated p-6 shadow-2xl">
        <h2 className="text-lg font-semibold text-foreground">Session about to expire</h2>
        <p className="mt-2 text-sm text-text-secondary">
          You will be signed out in <span className="tabular-nums font-medium text-foreground">{countdown}</span> seconds.
        </p>
        <div className="mt-6 flex justify-end gap-2">
          <button
            type="button"
            onClick={logout}
            className="px-4 py-2 rounded-md border border-border text-foreground hover:bg-bg-hover transition-colors text-sm"
          >
            Log out
          </button>
          <button
            type="button"
            onClick={stay}
            className="px-4 py-2 rounded-md bg-accent text-white hover:bg-accent-hover transition-colors text-sm font-medium"
          >
            Stay signed in
          </button>
        </div>
      </div>
    </div>
  )
}
`
}

// singleStatusBadge emits components/status-badge.tsx — typed status →
// color mapping. The TStatus generic lets a calling app extend the enum
// (e.g. <StatusBadge<MyStatus> status="shipped" />) while keeping
// type-safety on the supplied label/tone maps.
func singleStatusBadge() string {
	return `import type { ReactNode } from "react"
import { cn } from "@/lib/utils"

export type StatusTone = "success" | "warning" | "danger" | "info" | "neutral"

const TONE_CLASSES: Record<StatusTone, string> = {
  success: "bg-success/15 text-success border-success/30",
  warning: "bg-warning/15 text-warning border-warning/30",
  danger: "bg-danger/15 text-danger border-danger/30",
  info: "bg-info/15 text-info border-info/30",
  neutral: "bg-bg-hover text-text-secondary border-border",
}

/** Default tone for common payment/order statuses. Override via toneFor prop. */
export const DEFAULT_TONE_MAP: Record<string, StatusTone> = {
  paid: "success",
  active: "success",
  approved: "success",
  completed: "success",
  pending: "warning",
  partially_paid: "warning",
  processing: "info",
  draft: "neutral",
  cancelled: "danger",
  failed: "danger",
  overdue: "danger",
  rejected: "danger",
  disabled: "neutral",
}

export interface StatusBadgeProps<TStatus extends string = string> {
  status: TStatus
  /** Optional explicit label override. Defaults to title-cased status. */
  label?: ReactNode
  toneFor?: (status: TStatus) => StatusTone
  className?: string
}

function defaultLabel(status: string) {
  return status.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase())
}

export function StatusBadge<TStatus extends string = string>({
  status,
  label,
  toneFor,
  className,
}: StatusBadgeProps<TStatus>) {
  const tone = toneFor ? toneFor(status) : DEFAULT_TONE_MAP[status] || "neutral"
  return (
    <span
      className={cn(
        "inline-flex items-center px-2 py-0.5 rounded-md border text-xs font-medium",
        TONE_CLASSES[tone],
        className,
      )}
    >
      {label || defaultLabel(status)}
    </span>
  )
}
`
}

// singleStatsRow emits components/stats-row.tsx — list of stat cards.
// Grid auto-fits between 3 and 5 columns; each card supports an optional
// icon and sub-label.
func singleStatsRow() string {
	return `import type { ComponentType, ReactNode } from "react"
import { cn } from "@/lib/utils"

export type StatTone = "default" | "success" | "warning" | "danger" | "info"

const TONE_CLASSES: Record<StatTone, string> = {
  default: "text-foreground",
  success: "text-success",
  warning: "text-warning",
  danger: "text-danger",
  info: "text-info",
}

export interface StatCard {
  label: string
  value: ReactNode
  sub?: ReactNode
  icon?: ComponentType<{ className?: string }>
  tone?: StatTone
}

export interface StatsRowProps {
  cards: StatCard[]
  className?: string
}

// gridColsClass picks a static Tailwind class so the JIT compiler can find
// it. Dynamic concatenation like "lg:grid-cols-" + n gets purged.
const GRID_COLS: Record<number, string> = {
  1: "lg:grid-cols-1",
  2: "lg:grid-cols-2",
  3: "lg:grid-cols-3",
  4: "lg:grid-cols-4",
  5: "lg:grid-cols-5",
}

export function StatsRow({ cards, className }: StatsRowProps) {
  if (cards.length === 0) return null
  const cols = GRID_COLS[Math.min(cards.length, 5)] || GRID_COLS[5]
  return (
    <div
      className={cn(
        "grid gap-3 grid-cols-1 sm:grid-cols-2",
        cols,
        className,
      )}
    >
      {cards.map((card) => {
        const Icon = card.icon
        const tone = card.tone || "default"
        return (
          <div
            key={card.label}
            className="rounded-lg border border-border bg-bg-elevated p-4"
          >
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium uppercase tracking-wide text-text-muted">
                {card.label}
              </span>
              {Icon ? <Icon className={cn("h-4 w-4", TONE_CLASSES[tone])} /> : null}
            </div>
            <div className={cn("mt-2 text-2xl font-semibold tabular-nums", TONE_CLASSES[tone])}>
              {card.value}
            </div>
            {card.sub ? (
              <div className="mt-1 text-xs text-text-secondary">{card.sub}</div>
            ) : null}
          </div>
        )
      })}
    </div>
  )
}
`
}
