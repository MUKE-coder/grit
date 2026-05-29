'use client'

import { Highlight, type PrismTheme } from 'prism-react-renderer'
import { useTheme } from 'next-themes'
import { useEffect, useState } from 'react'
import { CopyButton } from './copy-button'

// GitHub Dark — tuned for our dark theme (matches gh.com syntax pal).
const githubDark: PrismTheme = {
  plain: { color: '#c9d1d9', backgroundColor: 'transparent' },
  styles: [
    { types: ['comment', 'prolog', 'doctype', 'cdata'], style: { color: '#8b949e', fontStyle: 'italic' as const } },
    { types: ['keyword', 'tag', 'operator', 'builtin'], style: { color: '#ff7b72' } },
    { types: ['string', 'attr-value', 'template-string', 'char'], style: { color: '#a5d6ff' } },
    { types: ['number', 'boolean'], style: { color: '#79c0ff' } },
    { types: ['function', 'method'], style: { color: '#d2a8ff' } },
    { types: ['class-name', 'type-name', 'maybe-class-name'], style: { color: '#ffa657' } },
    { types: ['punctuation'], style: { color: '#c9d1d9' } },
    { types: ['variable', 'constant'], style: { color: '#79c0ff' } },
    { types: ['attr-name'], style: { color: '#7ee787' } },
    { types: ['selector', 'property'], style: { color: '#79c0ff' } },
    { types: ['namespace'], style: { color: '#8b949e' } },
    { types: ['deleted'], style: { color: '#ffa198', backgroundColor: '#67060c' } },
    { types: ['inserted'], style: { color: '#aff5b4', backgroundColor: '#033a16' } },
    { types: ['regex'], style: { color: '#79c0ff' } },
  ],
}

// GitHub Light — the canonical gh.com light palette.
const githubLight: PrismTheme = {
  plain: { color: '#24292f', backgroundColor: 'transparent' },
  styles: [
    { types: ['comment', 'prolog', 'doctype', 'cdata'], style: { color: '#6e7781', fontStyle: 'italic' as const } },
    { types: ['keyword', 'tag', 'operator', 'builtin'], style: { color: '#cf222e' } },
    { types: ['string', 'attr-value', 'template-string', 'char'], style: { color: '#0a3069' } },
    { types: ['number', 'boolean'], style: { color: '#0550ae' } },
    { types: ['function', 'method'], style: { color: '#8250df' } },
    { types: ['class-name', 'type-name', 'maybe-class-name'], style: { color: '#953800' } },
    { types: ['punctuation'], style: { color: '#24292f' } },
    { types: ['variable', 'constant'], style: { color: '#0550ae' } },
    { types: ['attr-name'], style: { color: '#116329' } },
    { types: ['selector', 'property'], style: { color: '#0550ae' } },
    { types: ['namespace'], style: { color: '#6e7781' } },
    { types: ['deleted'], style: { color: '#82071e', backgroundColor: '#ffebe9' } },
    { types: ['inserted'], style: { color: '#116329', backgroundColor: '#dafbe1' } },
    { types: ['regex'], style: { color: '#0550ae' } },
  ],
}

interface CodeBlockProps {
  code: string
  language?: string
  filename?: string
  terminal?: boolean
  className?: string
  highlightLines?: number[]
}

export function CodeBlock({
  code,
  language = 'go',
  filename,
  terminal,
  className = 'mb-6',
  highlightLines = [],
}: CodeBlockProps) {
  const trimmedCode = code.trim()

  // Theme-aware Prism palette. Defer the hook reading until after mount
  // to avoid hydration mismatches (next-themes returns undefined on first
  // server render).
  const { resolvedTheme } = useTheme()
  const [mounted, setMounted] = useState(false)
  useEffect(() => setMounted(true), [])
  const theme = !mounted ? githubDark : resolvedTheme === 'light' ? githubLight : githubDark

  return (
    <div
      className={`rounded-lg overflow-hidden relative group ${className}
        border-2 border-slate-200/80 bg-white shadow-[0_1px_0_rgba(27,31,35,0.04),0_8px_24px_-12px_rgba(15,23,42,0.15)]
        dark:border-white/10 dark:bg-[#0d1117] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.04),0_20px_40px_-12px_rgba(0,0,0,0.6)]`}
    >
      {/* Header */}
      {(filename || terminal) && (
        <div className="flex items-center justify-between px-4 py-2.5 border-b border-slate-200 dark:border-white/[0.06] bg-slate-50/70 dark:bg-white/[0.02]">
          <div className="flex items-center gap-2.5">
            {terminal && (
              <div className="flex items-center gap-1.5 mr-1">
                <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
              </div>
            )}
            <span className="text-[12px] font-mono text-slate-600 dark:text-slate-400">
              {filename || 'Terminal'}
            </span>
          </div>
          <CopyButton
            text={terminal ? trimmedCode.replace(/^\$\s*/gm, '') : trimmedCode}
            className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7 text-slate-500 hover:text-slate-900 dark:hover:text-slate-100"
          />
        </div>
      )}

      {/* Code content — terminal mode adds a $ prompt gutter and forces
          language=bash so shell commands get real syntax highlighting
          instead of falling back to plain text. */}
      {terminal ? (
        <Highlight theme={theme} code={trimmedCode.replace(/^\$\s*/gm, '')} language="bash">
          {({ tokens, getLineProps, getTokenProps }) => (
            <pre className="p-4 text-[13px] leading-6 font-mono overflow-x-auto" style={{ background: 'transparent' }}>
              {tokens.map((line, i) => {
                const raw = trimmedCode.split('\n')[i] ?? ''
                if (raw.trim() === '') return <div key={i} className="h-5" />
                const isComment = raw.trimStart().startsWith('#')
                return (
                  <div key={i} {...getLineProps({ line })} className="flex">
                    {!isComment && (
                      <span className="text-sky-600/70 dark:text-sky-400/60 select-none mr-2 shrink-0">$</span>
                    )}
                    <span className={isComment ? 'text-slate-500 dark:text-slate-500 text-xs' : ''}>
                      {line.map((token, key) => (
                        <span key={key} {...getTokenProps({ token })} />
                      ))}
                    </span>
                  </div>
                )
              })}
            </pre>
          )}
        </Highlight>
      ) : (
        <Highlight theme={theme} code={trimmedCode} language={language}>
          {({ tokens, getLineProps, getTokenProps }) => (
            <pre className="p-4 text-[13px] leading-6 font-mono overflow-x-auto" style={{ background: 'transparent' }}>
              {tokens.map((line, i) => {
                const lineNum = i + 1
                const isHighlighted = highlightLines.includes(lineNum)
                return (
                  <div
                    key={i}
                    {...getLineProps({ line })}
                    className={isHighlighted ? 'bg-sky-500/[0.08] border-l-2 border-sky-400 -mx-4 px-4 pl-[14px]' : ''}
                  >
                    {line.map((token, key) => (
                      <span key={key} {...getTokenProps({ token })} />
                    ))}
                  </div>
                )
              })}
            </pre>
          )}
        </Highlight>
      )}

      {/* No header — show copy button in corner */}
      {!filename && !terminal && (
        <div className="absolute top-2 right-2">
          <CopyButton
            text={trimmedCode}
            className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7 text-slate-500 hover:text-slate-900 dark:hover:text-slate-100"
          />
        </div>
      )}
    </div>
  )
}

// Step component for Tailwind-style numbered installation steps
export function Step({ number, title, children }: { number: string; title: string; children: React.ReactNode }) {
  return (
    <div className="relative grid grid-cols-1 lg:grid-cols-[minmax(min-content,350px)_minmax(min-content,1fr)] gap-x-10 gap-y-4 pb-16">
      <div>
        <div className="flex items-center gap-3 mb-3">
          <span className="inline-flex items-center justify-center h-7 px-2 rounded border border-slate-200 dark:border-border/40 bg-slate-100 dark:bg-accent/20 text-[11px] font-mono font-medium text-slate-500 dark:text-muted-foreground">
            {number}
          </span>
          <h3 className="text-base font-semibold text-foreground">{title}</h3>
        </div>
        <div className="text-[15px] text-muted-foreground leading-relaxed">
          {children}
        </div>
      </div>
    </div>
  )
}

// StepWithCode component — text on left, code on right (like Tailwind docs)
export function StepWithCode({ number, title, description, code, filename, language, highlightLines }: {
  number: string
  title: string
  description: React.ReactNode
  code: string
  filename?: string
  language?: string
  highlightLines?: number[]
}) {
  return (
    <div className="relative grid grid-cols-1 lg:grid-cols-[minmax(min-content,380px)_minmax(min-content,1fr)] gap-x-10 gap-y-4 pb-14">
      <div>
        <div className="flex items-center gap-3 mb-3">
          <span className="inline-flex items-center justify-center h-7 px-2 rounded border border-slate-200 dark:border-border/40 bg-slate-100 dark:bg-accent/20 text-[11px] font-mono font-medium text-slate-500 dark:text-muted-foreground">
            {number}
          </span>
          <h3 className="text-base font-semibold text-foreground">{title}</h3>
        </div>
        <div className="text-[15px] text-muted-foreground leading-relaxed">
          {description}
        </div>
      </div>
      <div className="lg:mt-0">
        <CodeBlock
          code={code}
          filename={filename}
          language={language}
          highlightLines={highlightLines}
          className=""
        />
      </div>
    </div>
  )
}
