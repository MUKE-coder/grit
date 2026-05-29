import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { Save, Plus, Building2, Check } from 'lucide-react'
import toast from 'react-hot-toast'
import { useAuthStore } from '@/stores/auth.store'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Modal } from '@/components/ui/Modal'
import api from '@/lib/api'

export const Route = createFileRoute('/_app/settings')({
  component: SettingsPage,
})

function SettingsPage() {
  const queryClient = useQueryClient()
  const businesses = useAuthStore((s) => s.businesses)
  const currentBusinessID = useAuthStore((s) => s.currentBusinessID)
  const switchBusiness = useAuthStore((s) => s.switchBusiness)
  const login = useAuthStore((s) => s.login)
  const user = useAuthStore((s) => s.user)
  const tokens = useAuthStore((s) => s.tokens)
  const currentBiz = businesses.find((b) => b.id === currentBusinessID)

  const [name, setName] = useState(currentBiz?.name || '')
  const [saving, setSaving] = useState(false)
  const [showCreate, setShowCreate] = useState(false)

  // Fetch fresh businesses list
  const { data: freshBusinesses } = useQuery({
    queryKey: ['businesses'],
    queryFn: async () => {
      const res = await api.get('/businesses')
      return res.data.data as any[]
    },
  })

  const allBusinesses = freshBusinesses || businesses

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.put(`/businesses/${currentBusinessID}`, { name })
      queryClient.invalidateQueries({ queryKey: ['businesses'] })
      // Update store
      if (user && tokens) {
        const updated = businesses.map((b) =>
          b.id === currentBusinessID ? { ...b, name } : b
        )
        login(user, tokens, updated)
      }
      toast.success('Business name updated')
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to update')
    } finally {
      setSaving(false)
    }
  }

  const handleSwitch = (bizId: number) => {
    switchBusiness(bizId)
    const biz = allBusinesses.find((b: any) => b.id === bizId)
    setName(biz?.name || '')
    toast.success(`Switched to ${biz?.name}`)
    // Reload page to refresh all data for new business
    window.location.reload()
  }

  return (
    <div className="p-6 lg:p-8 max-w-2xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-text-primary">Settings</h1>
        <p className="text-text-muted mt-1">Business configuration</p>
      </div>

      {/* Current Business Name */}
      <form onSubmit={handleSave} className="bg-surface rounded-xl border border-border p-5 space-y-4 mb-6">
        <h3 className="text-sm font-semibold text-text-primary">Current Business</h3>
        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">Business Name</label>
          <input type="text" required value={name} onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue" />
        </div>
        <button type="submit" disabled={saving}
          className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-white font-semibold text-sm hover:bg-accent-hover transition disabled:opacity-50">
          <Save size={16} /> {saving ? 'Saving...' : 'Save Changes'}
        </button>
      </form>

      {/* All Businesses */}
      <div className="bg-surface rounded-xl border border-border shadow-sm">
        <div className="px-5 py-4 border-b border-border flex items-center justify-between">
          <div>
            <h3 className="text-sm font-semibold text-text-primary">Your Businesses</h3>
            <p className="text-xs text-text-muted mt-0.5">Switch between businesses or create a new one</p>
          </div>
          <button onClick={() => setShowCreate(true)}
            className="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg bg-accent text-white text-xs font-semibold hover:bg-accent-hover transition">
            <Plus size={14} /> New Business
          </button>
        </div>

        <div className="divide-y divide-border">
          {allBusinesses.map((biz: any) => {
            const isActive = biz.id === currentBusinessID
            return (
              <div key={biz.id} className={`px-5 py-4 flex items-center justify-between ${isActive ? 'bg-grit-blue-50/50' : ''}`}>
                <div className="flex items-center gap-3">
                  <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${isActive ? 'bg-grit-blue text-white' : 'bg-background text-text-muted'}`}>
                    <Building2 size={18} />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-text-primary">{biz.name}</p>
                    <p className="text-xs text-text-muted capitalize">
                      {biz.role || 'admin'}
                      {isActive && ' — Active'}
                    </p>
                  </div>
                </div>

                {isActive ? (
                  <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-grit-blue/10 text-grit-blue text-xs font-semibold">
                    <Check size={12} /> Active
                  </span>
                ) : (
                  <button onClick={() => handleSwitch(biz.id)}
                    className="px-3 py-1.5 rounded-lg border border-border text-xs font-semibold text-text-primary hover:bg-background transition">
                    Switch
                  </button>
                )}
              </div>
            )
          })}
        </div>
      </div>

      {/* Create Business Modal */}
      <CreateBusinessModal open={showCreate} onClose={() => setShowCreate(false)} />

    </div>
  )
}

function CreateBusinessModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const queryClient = useQueryClient()
  const user = useAuthStore((s) => s.user)
  const tokens = useAuthStore((s) => s.tokens)
  const businesses = useAuthStore((s) => s.businesses)
  const login = useAuthStore((s) => s.login)
  const switchBusiness = useAuthStore((s) => s.switchBusiness)

  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const res = await api.post('/businesses', { name })
      const newBiz = res.data.data.business
      queryClient.invalidateQueries({ queryKey: ['businesses'] })

      // Add to store and switch to it
      if (user && tokens) {
        const updated = [...businesses, { id: newBiz.id, name: newBiz.name, role: 'admin', branch_id: null }]
        login(user, tokens, updated)
        switchBusiness(newBiz.id)
      }

      toast.success(`"${name}" created! Switching to it...`)
      setName('')
      onClose()
      // Reload to load new business data
      setTimeout(() => window.location.reload(), 500)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to create business')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Create New Business">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-semibold text-text-primary mb-2">Business Name <span className="text-error">*</span></label>
          <input type="text" required value={name} onChange={(e) => setName(e.target.value)}
            className="w-full px-4 py-3 rounded-xl border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
            placeholder="e.g. Kampala Electronics" />
        </div>
        <p className="text-xs text-text-muted">
          A default "Main Branch" will be created automatically. You can add more branches later.
        </p>
        <button type="submit" disabled={loading}
          className="w-full py-3 rounded-xl bg-accent text-white font-semibold hover:bg-accent-hover transition disabled:opacity-50">
          {loading ? 'Creating...' : 'Create Business'}
        </button>
      </form>
    </Modal>
  )
}
