import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/auth.store'

export const Route = createFileRoute('/invite/$token')({
  component: AcceptInvitePage,
})

function AcceptInvitePage() {
  const { token } = Route.useParams()
  const navigate = useNavigate()
  const login = useAuthStore((s) => s.login)
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({ name: '', password: '', confirmPassword: '' })

  const { data: invite, isLoading, error } = useQuery({
    queryKey: ['invite', token],
    queryFn: async () => { const r = await api.get(`/invite/${token}`); return r.data.data },
  })

  const handleAccept = async (e: React.FormEvent) => {
    e.preventDefault()
    if (form.password !== form.confirmPassword) {
      toast.error('Passwords do not match')
      return
    }
    setLoading(true)
    try {
      const res = await api.post(`/invite/${token}/accept`, {
        name: form.name,
        password: form.password,
      })
      const { user, tokens } = res.data.data
      login(user, tokens, [{ id: 0, name: invite.business_name, role: invite.role, branch_id: null }])
      toast.success('Welcome to ' + invite.business_name + '!')
      navigate({ to: '/' })
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to accept invitation')
    } finally {
      setLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <p className="text-text-muted">Loading invitation...</p>
      </div>
    )
  }

  if (error || !invite) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background px-4">
        <div className="text-center">
          <div className="w-16 h-16 rounded-2xl bg-error/10 text-error flex items-center justify-center mx-auto mb-4 text-2xl font-bold">!</div>
          <h1 className="text-xl font-bold text-text-primary mb-2">Invalid Invitation</h1>
          <p className="text-text-muted">This invitation may have expired or already been accepted.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-grit-blue text-white text-2xl font-bold shadow-lg mb-4">S</div>
          <h1 className="text-2xl font-bold text-text-primary">You're invited!</h1>
          <p className="text-text-muted mt-2">
            Join <strong>{invite.business_name}</strong> as a{' '}
            <span className="capitalize font-medium text-grit-blue">{invite.role?.replace('_', ' ')}</span>
          </p>
        </div>

        <form onSubmit={handleAccept} className="bg-surface rounded-xl shadow-sm border border-border p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-text-primary mb-1.5">Full Name</label>
            <input type="text" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
              placeholder="Your name" />
          </div>
          <div>
            <label className="block text-sm font-medium text-text-primary mb-1.5">Password</label>
            <input type="password" required minLength={6} value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
              placeholder="••••••••" />
          </div>
          <div>
            <label className="block text-sm font-medium text-text-primary mb-1.5">Confirm Password</label>
            <input type="password" required minLength={6} value={form.confirmPassword} onChange={(e) => setForm({ ...form, confirmPassword: e.target.value })}
              className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
              placeholder="••••••••" />
          </div>
          <button type="submit" disabled={loading}
            className="w-full py-2.5 rounded-lg bg-accent text-white font-semibold hover:bg-accent-hover transition disabled:opacity-50">
            {loading ? 'Accepting...' : 'Accept & Join'}
          </button>
        </form>
      </div>
    </div>
  )
}
