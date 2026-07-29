'use client'

import { useState } from 'react'
import { Check, Copy } from 'lucide-react'

/**
 * Copy-to-clipboard with a confirmation state.
 *
 * navigator.clipboard is unavailable on insecure origins, so the failure path
 * is handled rather than swallowed — a button that looks like it copied but
 * did not is worse than one that visibly refuses.
 */
export function CopyButton({
  value,
  label,
  className = '',
}: {
  value: string
  label?: string
  className?: string
}) {
  const [state, setState] = useState<'idle' | 'copied' | 'failed'>('idle')

  async function copy() {
    try {
      await navigator.clipboard.writeText(value)
      setState('copied')
    } catch {
      setState('failed')
    }
    setTimeout(() => setState('idle'), 1800)
  }

  return (
    <button
      type="button"
      onClick={copy}
      aria-label={label ? `Copy ${label}` : 'Copy to clipboard'}
      className={`inline-flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs font-medium transition-colors ${
        state === 'failed'
          ? 'text-danger'
          : 'text-text-muted hover:text-foreground hover:bg-bg-hover'
      } ${className}`}
    >
      {state === 'copied' ? (
        <>
          <Check size={13} className="text-success" />
          Copied
        </>
      ) : state === 'failed' ? (
        <>
          <Copy size={13} />
          Press &#8984;C
        </>
      ) : (
        <>
          <Copy size={13} />
          {label ?? 'Copy'}
        </>
      )}
    </button>
  )
}
