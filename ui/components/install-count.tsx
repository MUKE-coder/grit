import { Download } from 'lucide-react'
import { formatCount, installLabel } from '@/lib/format-count'

/**
 * The install badge — "1.2k installs".
 *
 * Renders NOTHING when the count is absent or zero. Two reasons: a fresh block
 * showing "0 installs" reads as "nobody wanted this", which is a worse first
 * impression than no badge at all; and if the analytics read fails, every block
 * would claim zero rather than simply staying quiet.
 */
export function InstallCount({
  count,
  className = '',
}: {
  count?: number
  className?: string
}) {
  if (!count || count <= 0) return null

  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full bg-gray-500/[0.08] px-2.5 py-1 text-[11.5px] font-medium text-gray-600 dark:bg-white/[0.07] dark:text-gray-400 ${className}`}
      // The formatted value is approximate; expose the exact number to anyone
      // who hovers or uses a screen reader.
      title={`${count.toLocaleString()} ${installLabel(count)}`}
    >
      <Download aria-hidden="true" className="size-3.5 text-indigo-500 dark:text-indigo-400" />
      <span className="tabular-nums">{formatCount(count)}</span>
      <span className="text-gray-400 dark:text-gray-500">{installLabel(count)}</span>
    </span>
  )
}
