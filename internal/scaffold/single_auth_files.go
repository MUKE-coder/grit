package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Auth pages and the customer area for a single project's SPA.
//
// The SPA has shipped a complete auth library since it existed: login, register,
// refresh, logout, the TOTP challenge, a session-expiry monitor. It has never
// shipped a screen that calls any of it. `grit add web-auth` writes those screens
// for a Next.js web app; a single project is one binary with a Vite SPA, so the
// same command writes the TanStack Router versions here.
//
// The sections are layout routes, not a list of paths:
//
//	routes/_site.tsx      the public site: navbar and footer
//	routes/_auth.tsx      login and registration: full bleed, no chrome
//	routes/account/       the signed-in customer area: its own shell, guarded
//	routes/admin/         the admin panel, which has its own layout already
//
// _site and _auth start with an underscore, so they are pathless: they wrap
// pages without adding a segment to the URL, the way a route group does in Next.

// singleAuthFiles is everything `grit add web-auth` writes into a single project.
func singleAuthFiles(feRoot string, opts Options) []webAuthFile {
	src := filepath.Join(feRoot, "src")
	return []webAuthFile{
		{filepath.Join(src, "hooks", "use-auth.ts"), singleUseAuthHook()},
		{filepath.Join(src, "components", "user-menu.tsx"), singleUserMenu()},

		// The auth section: pathless, so /login is still /login.
		{filepath.Join(src, "routes", "_auth.tsx"), singleAuthLayoutRoute(opts)},
		{filepath.Join(src, "routes", "_auth", "login.tsx"), singleLoginRoute()},
		{filepath.Join(src, "routes", "_auth", "register.tsx"), singleRegisterRoute()},
		{filepath.Join(src, "routes", "_auth", "forgot-password.tsx"), singleForgotPasswordRoute()},
		{filepath.Join(src, "routes", "_auth", "reset-password.tsx"), singleResetPasswordRoute()},

		// The customer area.
		{filepath.Join(src, "routes", "account", "route.tsx"), singleAccountLayoutRoute(opts)},
		{filepath.Join(src, "routes", "account", "index.tsx"), singleAccountOverviewRoute()},
		{filepath.Join(src, "routes", "account", "profile.tsx"), singleAccountProfileRoute()},

		// REPLACEMENT: the navbar gains the user menu. Overwritten only with
		// --force, or when it is byte-for-byte the one the scaffold wrote.
		{filepath.Join(src, "components", "navbar.tsx"), singleViteNavbarWithAuth(opts)},
	}
}

// singleUseAuthHook wraps lib/auth.ts in React Query, so a screen can read the
// session without every page fetching it again.
func singleUseAuthHook() string {
	return `import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'

import { isAuthenticated, login, logout, me, register, type User } from '@/lib/auth'

// The session, cached. lib/auth.ts does the talking; this decides when.
//
// A missing token is not an error, it is the signed-out state: useMe returns
// null rather than throwing, so a page can branch on it without a try/catch.
export function useMe() {
  return useQuery<User | null>({
    queryKey: ['me'],
    queryFn: async () => {
      if (!isAuthenticated()) return null
      try {
        return await me()
      } catch {
        return null
      }
    },
    retry: false,
    staleTime: 10 * 60 * 1000,
  })
}

export function useAuth() {
  const { data: user, isLoading } = useMe()
  return { user: user ?? null, isAuthenticated: !!user, isLoading }
}

export function useLogin() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (credentials: { email: string; password: string }) =>
      login(credentials.email, credentials.password),
    onSuccess: (result) => {
      // Two-factor is on for this account: the server issued a pending token
      // rather than a session, and the screen has to ask for the code.
      if ('totp' in result) return
      queryClient.setQueryData(['me'], result.user)
      navigate({ to: '/account' })
    },
  })
}

export function useRegister() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: register,
    onSuccess: (result) => {
      queryClient.setQueryData(['me'], result.user)
      navigate({ to: '/account' })
    },
  })
}

export function useLogout() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: logout,
    onSettled: () => {
      queryClient.clear()
      // A full reload, not a navigate: it drops every cached page along with
      // the token, so nothing signed-in survives the sign-out.
      window.location.href = '/'
    },
  })
}
`
}

// singleUserMenu is the navbar's right-hand side: sign-in buttons, or the account.
func singleUserMenu() string {
	return `import { useEffect, useRef, useState } from 'react'
import { Link } from '@tanstack/react-router'
import { ChevronDown, LogOut, User as UserIcon } from 'lucide-react'

import { useMe, useLogout } from '@/hooks/use-auth'

// Signed out: Log in and Sign up. Signed in: the account menu. While useMe is
// in flight it renders a placeholder of the same size, so the navbar does not
// jump once the answer arrives.
export function UserMenu() {
  const { data: user, isLoading } = useMe()
  const logout = useLogout()
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    const close = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', close)
    return () => document.removeEventListener('mousedown', close)
  }, [open])

  if (isLoading) {
    return <div className="h-9 w-24 rounded-lg bg-bg-hover" aria-hidden />
  }

  if (!user) {
    return (
      <div className="flex items-center gap-2">
        <Link
          to="/login"
          className="rounded-lg px-3 py-2 text-sm text-text-secondary transition-colors hover:text-foreground"
        >
          Log in
        </Link>
        <Link
          to="/register"
          className="rounded-lg bg-accent px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-accent-hover"
        >
          Sign up
        </Link>
      </div>
    )
  }

  const initials =
    ((user.first_name?.[0] ?? '') + (user.last_name?.[0] ?? '')).toUpperCase() ||
    user.email.slice(0, 2).toUpperCase()

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-2 rounded-lg px-2 py-1.5 text-sm transition-colors hover:bg-bg-hover"
      >
        <span className="flex h-8 w-8 items-center justify-center rounded-full bg-accent/10 text-xs font-semibold text-accent">
          {initials}
        </span>
        <ChevronDown className="h-4 w-4 text-text-secondary" />
      </button>

      {open ? (
        <div className="absolute right-0 z-50 mt-2 w-56 rounded-xl border border-border bg-background p-1.5 shadow-lg">
          <div className="px-3 py-2">
            <p className="truncate text-sm font-medium text-foreground">
              {[user.first_name, user.last_name].filter(Boolean).join(' ') || user.email}
            </p>
            <p className="truncate text-xs text-text-secondary">{user.email}</p>
          </div>
          <Link
            to="/account"
            onClick={() => setOpen(false)}
            className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground"
          >
            <UserIcon className="h-4 w-4" />
            Account
          </Link>
          <button
            type="button"
            onClick={() => logout.mutate()}
            className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground"
          >
            <LogOut className="h-4 w-4" />
            Sign out
          </button>
        </div>
      ) : null}
    </div>
  )
}
`
}

// singleAuthLayoutRoute is the shell every auth screen renders inside.
func singleAuthLayoutRoute(opts Options) string {
	return fmt.Sprintf(`import { createFileRoute, Link, Outlet } from '@tanstack/react-router'

// The auth section. Pathless: this file adds no segment, so its children are at
// /login and /register. Full bleed, because a sign-in page with a marketing
// navbar above it is a sign-in page nobody finishes.
export const Route = createFileRoute('/_auth')({
  component: AuthLayout,
})

function AuthLayout() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-bg-secondary px-6 py-12">
      <div className="w-full max-w-md">
        <Link to="/" className="mb-8 flex items-center justify-center gap-2.5">
          <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-sm font-bold text-accent">
            %s
          </span>
          <span className="text-lg font-bold tracking-tight text-foreground">%s</span>
        </Link>
        <div className="rounded-2xl border border-border bg-background p-8">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
`, brandInitial(opts), opts.ProjectName)
}

// singleLoginRoute is the customer sign-in screen.
func singleLoginRoute() string {
	return `import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'

import { useLogin } from '@/hooks/use-auth'

export const Route = createFileRoute('/_auth/login')({
  component: LoginPage,
})

function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const login = useLogin()

  // Two-factor: the server answered with a pending token instead of a session.
  const totp = login.data && 'totp' in login.data ? login.data.totp : null

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold tracking-tight text-foreground">Welcome back</h1>
        <p className="mt-1 text-sm text-text-secondary">Sign in to your account</p>
      </div>

      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault()
          login.mutate({ email, password })
        }}
      >
        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Email</span>
          <input
            type="email"
            required
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent"
          />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Password</span>
          <input
            type="password"
            required
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent"
          />
        </label>

        {login.isError ? (
          <p className="text-sm text-danger">That email and password did not match.</p>
        ) : null}
        {totp ? (
          <p className="text-sm text-text-secondary">
            This account uses two-factor authentication. Finish signing in from the admin
            panel, or wire the verification step here with the pending token.
          </p>
        ) : null}

        <button
          type="submit"
          disabled={login.isPending}
          className="w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
        >
          {login.isPending ? 'Signing in...' : 'Sign in'}
        </button>
      </form>

      <div className="flex items-center justify-between text-sm">
        <Link to="/forgot-password" className="text-text-secondary hover:text-foreground">
          Forgot password?
        </Link>
        <Link to="/register" className="text-accent hover:text-accent-hover">
          Create an account
        </Link>
      </div>
    </div>
  )
}
`
}

// singleRegisterRoute is the sign-up screen.
func singleRegisterRoute() string {
	return `import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'

import { useRegister } from '@/hooks/use-auth'

export const Route = createFileRoute('/_auth/register')({
  component: RegisterPage,
})

function RegisterPage() {
  const [form, setForm] = useState({
    first_name: '',
    last_name: '',
    email: '',
    password: '',
  })
  const register = useRegister()

  const field = (key: keyof typeof form) => ({
    value: form[key],
    onChange: (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((f) => ({ ...f, [key]: e.target.value })),
    className:
      'w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent',
  })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold tracking-tight text-foreground">Create your account</h1>
        <p className="mt-1 text-sm text-text-secondary">It takes a moment</p>
      </div>

      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault()
          register.mutate(form)
        }}
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">First name</span>
            <input required autoComplete="given-name" {...field('first_name')} />
          </label>
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">Last name</span>
            <input required autoComplete="family-name" {...field('last_name')} />
          </label>
        </div>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Email</span>
          <input type="email" required autoComplete="email" {...field('email')} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Password</span>
          <input type="password" required autoComplete="new-password" {...field('password')} />
        </label>

        {register.isError ? (
          <p className="text-sm text-danger">
            That did not work. The email may already be registered.
          </p>
        ) : null}

        <button
          type="submit"
          disabled={register.isPending}
          className="w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
        >
          {register.isPending ? 'Creating...' : 'Create account'}
        </button>
      </form>

      <p className="text-center text-sm text-text-secondary">
        Already have an account?{' '}
        <Link to="/login" className="text-accent hover:text-accent-hover">
          Sign in
        </Link>
      </p>
    </div>
  )
}
`
}

// singleForgotPasswordRoute asks the API to send a reset link.
func singleForgotPasswordRoute() string {
	return `import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { useMutation } from '@tanstack/react-query'

import { api } from '@/lib/api'

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const [email, setEmail] = useState('')
  const send = useMutation({
    mutationFn: async () => api.post('/api/auth/forgot-password', { email }),
  })

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold tracking-tight text-foreground">Reset your password</h1>
        <p className="mt-1 text-sm text-text-secondary">
          We will email you a link to set a new one.
        </p>
      </div>

      {send.isSuccess ? (
        <p className="text-sm text-text-secondary">
          If that address has an account, the link is on its way. It expires in an hour.
        </p>
      ) : (
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault()
            send.mutate()
          }}
        >
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">Email</span>
            <input
              type="email"
              required
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent"
            />
          </label>

          <button
            type="submit"
            disabled={send.isPending}
            className="w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
          >
            {send.isPending ? 'Sending...' : 'Send the link'}
          </button>
        </form>
      )}

      <p className="text-center text-sm">
        <Link to="/login" className="text-text-secondary hover:text-foreground">
          Back to sign in
        </Link>
      </p>
    </div>
  )
}
`
}

// singleResetPasswordRoute completes the reset, with the token from the link.
func singleResetPasswordRoute() string {
	return `import { useState } from 'react'
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useMutation } from '@tanstack/react-query'

import { api } from '@/lib/api'

// The token arrives in the query string of the emailed link.
export const Route = createFileRoute('/_auth/reset-password')({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === 'string' ? search.token : '',
  }),
  component: ResetPasswordPage,
})

function ResetPasswordPage() {
  const { token } = Route.useSearch()
  const navigate = useNavigate()
  const [password, setPassword] = useState('')

  const reset = useMutation({
    mutationFn: async () =>
      api.post('/api/auth/reset-password', { token, password }),
    onSuccess: () => {
      navigate({ to: '/login' })
    },
  })

  if (!token) {
    return (
      <div className="space-y-4">
        <h1 className="text-xl font-bold tracking-tight text-foreground">Link incomplete</h1>
        <p className="text-sm text-text-secondary">
          This page needs the token from the email we sent. Open the link again, or ask
          for a new one.
        </p>
        <Link to="/forgot-password" className="text-sm text-accent hover:text-accent-hover">
          Send another link
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold tracking-tight text-foreground">Set a new password</h1>
        <p className="mt-1 text-sm text-text-secondary">Then sign in with it.</p>
      </div>

      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault()
          reset.mutate()
        }}
      >
        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">New password</span>
          <input
            type="password"
            required
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent"
          />
        </label>

        {reset.isError ? (
          <p className="text-sm text-danger">
            That link has expired or has already been used.
          </p>
        ) : null}

        <button
          type="submit"
          disabled={reset.isPending}
          className="w-full rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
        >
          {reset.isPending ? 'Saving...' : 'Save the password'}
        </button>
      </form>
    </div>
  )
}
`
}

// singleAccountLayoutRoute is the customer area's shell, and its guard.
func singleAccountLayoutRoute(opts Options) string {
	return fmt.Sprintf(`import { createFileRoute, Link, Outlet, redirect, useRouterState } from '@tanstack/react-router'
import { LayoutDashboard, LogOut, User as UserIcon } from 'lucide-react'

import { isAuthenticated } from '@/lib/auth'
import { useMe, useLogout } from '@/hooks/use-auth'

// The signed-in customer area. Everything under /account renders inside this
// shell, and beforeLoad turns away anybody without a token before a screen is
// ever rendered: no flash of a dashboard the visitor cannot see.
//
// Add a section by adding a line to NAV and a file under routes/account/.
export const Route = createFileRoute('/account')({
  beforeLoad: () => {
    if (!isAuthenticated()) {
      throw redirect({ to: '/login' })
    }
  },
  component: AccountLayout,
})

const NAV = [
  { to: '/account', label: 'Overview', icon: LayoutDashboard },
  { to: '/account/profile', label: 'Profile', icon: UserIcon },
]

function AccountLayout() {
  const pathname = useRouterState({ select: (s) => s.location.pathname })
  const { data: user } = useMe()
  const logout = useLogout()

  const links = NAV.map((item) => {
    const Icon = item.icon
    const active = item.to === '/account' ? pathname === '/account' : pathname.startsWith(item.to)
    return (
      <Link
        key={item.to}
        to={item.to}
        className={
          'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-colors ' +
          (active
            ? 'bg-accent/10 font-medium text-accent'
            : 'text-text-secondary hover:bg-bg-hover hover:text-foreground')
        }
      >
        <Icon className="h-4 w-4 shrink-0" />
        {item.label}
      </Link>
    )
  })

  const initials =
    ((user?.first_name?.[0] ?? '') + (user?.last_name?.[0] ?? '')).toUpperCase() ||
    (user?.email ?? '?').slice(0, 2).toUpperCase()

  return (
    <div className="flex min-h-screen bg-bg-secondary">
      <aside className="hidden w-64 shrink-0 flex-col border-r border-border bg-background md:flex">
        <Link to="/" className="flex h-20 items-center gap-3 px-6">
          <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-sm font-bold text-accent">
            %s
          </span>
          <span className="text-lg font-bold tracking-tight text-foreground">%s</span>
        </Link>

        <nav className="flex flex-1 flex-col gap-1 px-4">{links}</nav>

        <div className="border-t border-border p-4">
          <button
            type="button"
            onClick={() => logout.mutate()}
            className="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground"
          >
            <LogOut className="h-4 w-4 shrink-0" />
            Sign out
          </button>
        </div>
      </aside>

      <div className="flex min-h-screen flex-1 flex-col">
        <header className="flex h-20 items-center justify-between border-b border-border bg-background px-6">
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold text-foreground">
              Welcome, {user?.first_name || user?.email || 'there'}
            </p>
            <p className="truncate text-xs text-text-secondary">{user?.email}</p>
          </div>
          <div className="flex items-center gap-3">
            <Link
              to="/"
              className="hidden text-sm text-text-secondary transition-colors hover:text-foreground sm:block"
            >
              Back to site
            </Link>
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-accent/10 text-sm font-semibold text-accent">
              {initials}
            </span>
          </div>
        </header>

        <nav className="flex gap-2 overflow-x-auto border-b border-border bg-background px-4 py-2 md:hidden">
          {links}
        </nav>

        <main className="flex-1 px-6 py-8 lg:px-10">
          <div className="mx-auto max-w-4xl">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  )
}
`, brandInitial(opts), opts.ProjectName)
}

// singleAccountOverviewRoute is where a customer lands after signing in.
func singleAccountOverviewRoute() string {
	return `import { createFileRoute, Link } from '@tanstack/react-router'
import { CalendarDays, Mail, ShieldCheck } from 'lucide-react'

import { useMe } from '@/hooks/use-auth'

export const Route = createFileRoute('/account/')({
  component: AccountOverview,
})

function AccountOverview() {
  const { data: user, isLoading } = useMe()

  if (isLoading) {
    return <p className="text-sm text-text-secondary">Loading your account...</p>
  }

  const fullName = [user?.first_name, user?.last_name].filter(Boolean).join(' ') || '--'
  const joined = user?.created_at
    ? new Date(user.created_at).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : '--'

  const cards = [
    { icon: Mail, label: 'Email', value: user?.email ?? '--' },
    { icon: ShieldCheck, label: 'Role', value: user?.role ?? 'USER' },
    { icon: CalendarDays, label: 'Member since', value: joined },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <ShieldCheck className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Overview</h1>
          <p className="text-sm text-text-secondary">Your account at a glance</p>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        {cards.map((card) => {
          const Icon = card.icon
          return (
            <div key={card.label} className="rounded-xl border border-border bg-background p-5">
              <div className="flex items-center gap-2 text-text-secondary">
                <Icon className="h-4 w-4" />
                <p className="text-xs font-semibold uppercase tracking-wide">{card.label}</p>
              </div>
              <p className="mt-3 truncate text-sm text-foreground">{card.value}</p>
            </div>
          )
        })}
      </div>

      <div className="rounded-xl border border-border bg-background">
        <div className="border-b border-border px-6 py-4">
          <h2 className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            Profile
          </h2>
        </div>
        <dl className="divide-y divide-border">
          {[
            ['Full name', fullName],
            ['Email', user?.email ?? '--'],
            ['Job title', user?.job_title || '--'],
          ].map(([label, value]) => (
            <div key={label} className="flex items-center justify-between gap-4 px-6 py-4">
              <dt className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
                {label}
              </dt>
              <dd className="truncate text-sm text-foreground">{value}</dd>
            </div>
          ))}
        </dl>
        <div className="border-t border-border px-6 py-4">
          <Link
            to="/account/profile"
            className="inline-flex items-center rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-accent-hover"
          >
            Edit profile
          </Link>
        </div>
      </div>
    </div>
  )
}
`
}

// singleAccountProfileRoute lets a customer edit their own record.
func singleAccountProfileRoute() string {
	return `import { useEffect, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { UserCog } from 'lucide-react'

import { api } from '@/lib/api'
import { useMe } from '@/hooks/use-auth'

export const Route = createFileRoute('/account/profile')({
  component: AccountProfile,
})

// One endpoint does all of this: PUT /api/profile takes the name fields, the
// email, and a password when the customer wants to change it. An empty password
// is left out of the payload rather than sent as an empty string, which the API
// would hash.
function AccountProfile() {
  const { data: user } = useMe()
  const queryClient = useQueryClient()

  const [form, setForm] = useState({
    first_name: '',
    last_name: '',
    email: '',
    job_title: '',
    password: '',
  })
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    if (!user) return
    setForm({
      first_name: user.first_name ?? '',
      last_name: user.last_name ?? '',
      email: user.email ?? '',
      job_title: user.job_title ?? '',
      password: '',
    })
  }, [user])

  const save = useMutation({
    mutationFn: async () => {
      const payload: Record<string, string> = {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        job_title: form.job_title,
      }
      if (form.password) payload.password = form.password
      const { data } = await api.put('/api/profile', payload)
      return data
    },
    onSuccess: () => {
      setSaved(true)
      setForm((f) => ({ ...f, password: '' }))
      queryClient.invalidateQueries({ queryKey: ['me'] })
      window.setTimeout(() => setSaved(false), 4000)
    },
  })

  const field = (key: keyof typeof form) => ({
    value: form[key],
    onChange: (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((f) => ({ ...f, [key]: e.target.value })),
    className:
      'w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent',
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <UserCog className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Profile</h1>
          <p className="text-sm text-text-secondary">
            This is what the rest of the app knows about you
          </p>
        </div>
      </div>

      <form
        className="space-y-5 rounded-xl border border-border bg-background p-6"
        onSubmit={(e) => {
          e.preventDefault()
          save.mutate()
        }}
      >
        <div className="grid gap-5 sm:grid-cols-2">
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">First name</span>
            <input {...field('first_name')} />
          </label>
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">Last name</span>
            <input {...field('last_name')} />
          </label>
        </div>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Email</span>
          <input type="email" {...field('email')} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Job title</span>
          <input {...field('job_title')} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">New password</span>
          <input type="password" autoComplete="new-password" {...field('password')} />
          <span className="block text-xs text-text-secondary">
            Leave this empty to keep the password you have.
          </span>
        </label>

        <div className="flex items-center gap-3 pt-1">
          <button
            type="submit"
            disabled={save.isPending}
            className="inline-flex items-center rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
          >
            {save.isPending ? 'Saving...' : 'Save changes'}
          </button>
          {saved ? <span className="text-sm text-text-secondary">Saved.</span> : null}
          {save.isError ? (
            <span className="text-sm text-danger">
              That did not save. Check the fields and try again.
            </span>
          ) : null}
        </div>
      </form>
    </div>
  )
}
`
}

// brandInitial is the letter in the little square: the project's first, upper case.
func brandInitial(opts Options) string {
	if opts.ProjectName == "" {
		return "A"
	}
	return strings.ToUpper(string([]rune(opts.ProjectName)[0]))
}

// singleViteNavbarWithAuth is the SPA navbar with the account menu in it.
//
// Built from the base navbar rather than copied, so the two cannot drift: the
// only difference is the user menu, and a change to the logo or the links reaches
// both. The Next.js web app does the same thing in web_chrome_files.go.
func singleViteNavbarWithAuth(opts Options) string {
	nav := singleViteNavbar(opts)
	nav = strings.Replace(nav,
		`import { Menu, X, Github } from "lucide-react"`,
		`import { Menu, X, Github } from "lucide-react"

import { UserMenu } from "@/components/user-menu"`, 1)

	// The desktop row: the account menu sits after the GitHub icon.
	nav = strings.Replace(nav, `            <Github className="h-5 w-5" />
          </a>
        </div>`, `            <Github className="h-5 w-5" />
          </a>
          <UserMenu />
        </div>`, 1)

	// And the mobile drawer, where it is a pair of links rather than a dropdown.
	nav = strings.Replace(nav, `              GitHub
            </a>
          </div>`, `              GitHub
            </a>
            <div className="pt-2 border-t border-border/50">
              <UserMenu />
            </div>
          </div>`, 1)
	return nav
}
