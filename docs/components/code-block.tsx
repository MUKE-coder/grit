'use client'

import { Highlight, type PrismTheme } from 'prism-react-renderer'
import { CopyButton } from './copy-button'

const theme: PrismTheme = {
  plain: {
    color: 'hsl(var(--foreground) / 0.8)',
    backgroundColor: 'transparent',
  },
  styles: [
    { types: ['comment', 'prolog', 'doctype', 'cdata'], style: { color: 'hsl(var(--muted-foreground) / 0.5)', fontStyle: 'italic' as const } },
    { types: ['keyword', 'tag', 'operator', 'builtin'], style: { color: '#c792ea' } },
    { types: ['string', 'attr-value', 'template-string'], style: { color: '#c3e88d' } },
    { types: ['number', 'boolean'], style: { color: '#f78c6c' } },
    { types: ['function', 'method'], style: { color: '#82aaff' } },
    { types: ['class-name', 'type-name', 'maybe-class-name'], style: { color: '#ffcb6b' } },
    { types: ['punctuation'], style: { color: 'hsl(var(--foreground) / 0.5)' } },
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
}

export function CodeBlock({ code, language = 'go', filename, terminal, className = 'mb-6' }: CodeBlockProps) {
  const trimmedCode = code.trim()

  return (
    <div className={`rounded-xl border border-border/40 bg-card/80 overflow-hidden relative group ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-2.5 border-b border-border/30 bg-accent/30">
        <div className="flex items-center gap-2">
          {terminal && (
            <div className="flex items-center gap-1.5">
              <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
              <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
              <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
            </div>
          )}
          {(filename || terminal) && (
            <span className={`text-[11px] font-mono text-muted-foreground/40 ${terminal ? 'ml-2' : ''}`}>
              {filename || 'terminal'}
            </span>
          )}
        </div>
        <CopyButton
          text={terminal ? trimmedCode.replace(/^\$\s*/gm, '') : trimmedCode}
          className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7"
        />
      </div>

      {/* Code content */}
      {terminal ? (
        <div className="p-5 font-mono text-sm overflow-x-auto space-y-1">
          {trimmedCode.split('\n').map((line, i) => {
            if (line.trim() === '') {
              return <div key={i} className="h-1" />
            }
            if (line.trimStart().startsWith('#')) {
              return (
                <div key={i}>
                  <span className="text-muted-foreground/40 text-xs">{line}</span>
                </div>
              )
            }
            return (
              <div key={i}>
                <span className="text-primary/50 select-none">$ </span>
                <span className="text-foreground/80">{line}</span>
              </div>
            )
          })}
        </div>
      ) : (
        <Highlight theme={theme} code={trimmedCode} language={language}>
          {({ tokens, getLineProps, getTokenProps }) => (
            <pre className="p-5 text-sm font-mono overflow-x-auto" style={{ background: 'transparent' }}>
              {tokens.map((line, i) => (
                <div key={i} {...getLineProps({ line })}>
                  {line.map((token, key) => (
                    <span key={key} {...getTokenProps({ token })} />
                  ))}
                </div>
              ))}
            </pre>
          )}
        </Highlight>
      )}
    </div>
  )
}
