import type { ReactNode } from 'react'

interface StatCardProps {
  icon: ReactNode
  label: string
  value: string
  color?: string
  subtitle?: string
}

export function StatCard({ icon, label, value, color = 'bg-grit-blue', subtitle }: StatCardProps) {
  return (
    <div className="bg-surface rounded-xl border border-border p-5 shadow-sm">
      <div className="flex items-center gap-3 mb-3">
        <div className={`w-9 h-9 rounded-lg ${color} text-white flex items-center justify-center`}>
          {icon}
        </div>
        <p className="text-sm text-text-muted">{label}</p>
      </div>
      <p className="text-2xl font-bold text-text-primary">{value}</p>
      {subtitle && <p className="text-xs text-text-muted mt-1">{subtitle}</p>}
    </div>
  )
}
