import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { Eye, EyeOff } from 'lucide-react'
import toast from 'react-hot-toast'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/auth.store'

export const Route = createFileRoute('/_auth/login')({
  component: LoginPage,
})

// Demo credentials — pre-fill so a first-time visitor can hit Enter and
// land in the app immediately. The seeder guarantees this user always
// exists; the demo-reset cron restores the state to canonical every 24 h.
const DEMO_EMAIL = 'admin@grit.demo'
const DEMO_PASSWORD = 'password123'

function LoginPage() {
  const navigate = useNavigate()
  const login = useAuthStore((s) => s.login)
  const [loading, setLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [form, setForm] = useState({ email: DEMO_EMAIL, password: DEMO_PASSWORD })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const res = await api.post('/auth/login', form)
      const { user, businesses, tokens } = res.data.data
      login(user, tokens, businesses)
      toast.success(`Welcome back, ${user.name}!`)
      navigate({ to: '/' })
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <div className="mb-6 rounded-xl border border-grit-blue/20 bg-grit-blue/5 p-3.5 text-[13px] leading-snug">
        <p className="font-semibold text-grit-blue mb-1">Welcome to the Grit Framework demo</p>
        <p className="text-foreground-muted">
          Credentials are pre-filled — just hit <kbd className="px-1.5 py-0.5 rounded border border-border bg-bg font-mono text-[11px]">Enter</kbd> to sign in.
          The database resets every 24h; please be respectful when editing data.
        </p>
      </div>

      <div className="mb-8">
        <h1 className="text-[28px] font-semibold text-foreground tracking-tight">Welcome back</h1>
        <p className="text-foreground-muted mt-1.5 text-[14px]">Please enter your details</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-[12.5px] font-medium text-foreground-secondary mb-1.5">
            Email
          </label>
          <input
            type="email"
            required
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            className="w-full h-11 px-3.5 rounded-lg border border-border bg-surface text-[14px] text-foreground placeholder:text-foreground-muted focus:outline-none focus:ring-2 focus:ring-accent/15 focus:border-accent transition"
            placeholder="Enter your email"
          />
        </div>

        <div>
          <label className="block text-[12.5px] font-medium text-foreground-secondary mb-1.5">
            Password
          </label>
          <div className="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              required
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full h-11 pl-3.5 pr-11 rounded-lg border border-border bg-surface text-[14px] text-foreground placeholder:text-foreground-muted focus:outline-none focus:ring-2 focus:ring-accent/15 focus:border-accent transition"
              placeholder="Enter your password"
            />
            <button
              type="button"
              onClick={() => setShowPassword((s) => !s)}
              className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 inline-flex items-center justify-center rounded-md text-foreground-muted hover:text-foreground hover:bg-surface-hover transition"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
              tabIndex={-1}
            >
              {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
            </button>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full h-11 mt-2 rounded-lg bg-accent text-white font-semibold text-[14px] hover:bg-accent-hover transition disabled:opacity-60 disabled:cursor-not-allowed shadow-xs"
        >
          {loading ? 'Signing in...' : 'Sign in'}
        </button>
      </form>
    </>
  )
}
