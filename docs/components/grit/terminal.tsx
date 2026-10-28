'use client'

import { useCallback, useState } from 'react'

/**
 * Terminal — the design system's `g-terminal`.
 *
 * Code stays dark in both themes on purpose, so a terminal looks identical
 * everywhere and a screenshot of one is stable.
 */

export interface TerminalLine {
  /** The command, without the prompt. */
  run?: string
  /** What it printed. */
  out?: string
  /** A success line, rendered with a tick. */
  ok?: string
}

/**
 * Copy, then confirm in place for 1.6s.
 *
 * In place rather than as a toast: the confirmation belongs where the action
 * happened, and a toast for "copied" is a notification nobody asked for.
 */
export function CopyButton({ text, label = 'Copy' }: { text: string; label?: string }) {
  const [copied, setCopied] = useState(false)

  const copy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(text)
    } catch {
      // Clipboard is blocked without a user gesture or over http. Falling
      // back keeps the button working rather than failing silently.
      const el = document.createElement('textarea')
      el.value = text
      el.setAttribute('readonly', '')
      el.style.position = 'fixed'
      el.style.opacity = '0'
      document.body.appendChild(el)
      el.select()
      try {
        document.execCommand('copy')
      } finally {
        el.remove()
      }
    }
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1600)
  }, [text])

  return (
    <button
      type="button"
      className="g-copy"
      data-copied={copied}
      onClick={copy}
      aria-label={copied ? 'Copied' : label}
    >
      {copied ? (
        <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
          <path d="M13 4.5 6.5 11 3 7.5" />
        </svg>
      ) : (
        <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
          <rect x="5.5" y="5.5" width="8" height="8" rx="1.5" />
          <path d="M10.5 5.5v-2a1.5 1.5 0 0 0-1.5-1.5H4a1.5 1.5 0 0 0-1.5 1.5v5A1.5 1.5 0 0 0 4 10h2" />
        </svg>
      )}
    </button>
  )
}

export function Terminal({
  title,
  lines,
  copy,
  className = '',
}: {
  title?: string
  lines: TerminalLine[]
  /** What the copy button puts on the clipboard. Defaults to the commands. */
  copy?: string
  className?: string
}) {
  const clipboard = copy ?? lines.map((l) => l.run).filter(Boolean).join('\n')

  return (
    <div className={`g-terminal ${className}`.trim()}>
      <div className="g-terminal__head">
        <div className="g-terminal__dots" aria-hidden="true">
          <i />
          <i />
          <i />
        </div>
        {title && <span className="g-terminal__title">{title}</span>}
        {clipboard && <CopyButton text={clipboard} label="Copy commands" />}
      </div>

      <div className="g-terminal__body">
        {lines.map((line, i) => (
          <div key={i}>
            {line.run !== undefined && (
              <div className="g-terminal__line">
                <span className="g-terminal__prompt" aria-hidden="true">
                  $
                </span>
                {line.run}
              </div>
            )}
            {line.out !== undefined && <div className="g-terminal__out">{line.out}</div>}
            {line.ok !== undefined && (
              <div className="g-terminal__out">
                <span className="g-terminal__ok" aria-hidden="true">
                  ✓{' '}
                </span>
                {line.ok}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
