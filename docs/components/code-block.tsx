'use client'

import { Highlight, type PrismTheme } from 'prism-react-renderer'
import { CopyButton } from './copy-button'

const theme: PrismTheme = {
  plain: {
    color: '#e2e8f0',
    backgroundColor: 'transparent',
  },
  styles: [
    { types: ['comment', 'prolog', 'doctype', 'cdata'], style: { color: '#64748b', fontStyle: 'italic' as const } },
    { types: ['keyword', 'tag', 'operator', 'builtin'], style: { color: '#c792ea' } },
    { types: ['string', 'attr-value', 'template-string'], style: { color: '#c3e88d' } },
    { types: ['number', 'boolean'], style: { color: '#f78c6c' } },
    { types: ['function', 'method'], style: { color: '#82aaff' } },
    { types: ['class-name', 'type-name', 'maybe-class-name'], style: { color: '#ffcb6b' } },
    { types: ['punctuation'], style: { color: '#94a3b8' } },
    { types: ['variable', 'constant'], style: { color: '#f07178' } },
    { types: ['attr-name'], style: { color: '#ffcb6b' } },
    { types: ['selector', 'property'], style: { color: '#82aaff' } },
    { types: ['namespace'], style: { color: '#b2ccd6', opacity: 0.7 } },
    { types: ['deleted'], style: { color: '#ff6b6b' } },
    { types: ['inserted'], style: { color: '#c3e88d' } },
    { types: ['char', 'symbol'], style: { color: '#80cbc4' } },
    { types: ['regex'], style: { color: '#89ddff' } },
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

  return (
    <div className={`rounded-lg border border-slate-200 dark:border-border/40 bg-slate-900 overflow-hidden relative group ${className}`}>
      {/* Header */}
      {(filename || terminal) && (
        <div className="flex items-center justify-between px-4 py-2 border-b border-slate-700/50">
          <div className="flex items-center gap-2.5">
            {terminal && (
              <div className="flex items-center gap-1.5 mr-1">
                <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]/80" />
                <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]/80" />
                <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]/80" />
              </div>
            )}
            <span className="text-[12px] font-mono text-slate-500">
              {filename || 'Terminal'}
            </span>
          </div>
          <CopyButton
            text={terminal ? trimmedCode.replace(/^\$\s*/gm, '') : trimmedCode}
            className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7 text-slate-500 hover:text-slate-300"
          />
        </div>
      )}

      {/* Code content */}
      {terminal ? (
        <div className="p-4 font-mono text-[13px] leading-6 overflow-x-auto">
          {trimmedCode.split('\n').map((line, i) => {
            if (line.trim() === '') {
              return <div key={i} className="h-5" />
            }
            if (line.trimStart().startsWith('#')) {
              return (
                <div key={i}>
                  <span className="text-slate-600 text-xs">{line}</span>
                </div>
              )
            }
            return (
              <div key={i} className="flex">
                <span className="text-sky-400/60 select-none mr-2">$</span>
                <span className="text-slate-300">{line}</span>
              </div>
            )
          })}
        </div>
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
            className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7 text-slate-500 hover:text-slate-300"
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
