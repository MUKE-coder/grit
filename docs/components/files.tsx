'use client'

import { useState, type ReactNode } from 'react'
import { ChevronRight, FileCode } from 'lucide-react'
import { fileColor, folderEmoji } from '@/lib/file-meta'

/**
 * Composable file-tree, à la Fumadocs:
 *
 *   <Files title="myapp/">
 *     <Folder name="apps" defaultOpen>
 *       <File name="main.go" />
 *     </Folder>
 *     <File name="grit.json" />
 *   </Files>
 *
 * Vertical guide lines, emoji folders (📁 / 📂), and per-type file colours.
 */

export function Files({ title, children }: { title?: string; children: ReactNode }) {
  return (
    <div className="my-6 overflow-hidden rounded-xl border border-border/50 bg-card/40">
      {title && (
        <div className="border-b border-border/50 bg-muted/20 px-4 py-2 font-mono text-[12px] text-muted-foreground">
          {title}
        </div>
      )}
      <div className="p-3 text-[13px]">{children}</div>
    </div>
  )
}

export function File({
  name,
  icon,
  comment,
  highlight,
}: {
  name: string
  /** Override the default coloured glyph (e.g. an emoji). */
  icon?: ReactNode
  comment?: string
  /** Tint the row — e.g. for a generated file. */
  highlight?: boolean
}) {
  return (
    <div
      className={`flex items-center gap-2 rounded-md py-1 pl-1.5 pr-2 ${
        highlight ? 'bg-primary/[0.07]' : ''
      }`}
    >
      {icon ?? <FileCode className={`h-3.5 w-3.5 shrink-0 ${fileColor(name)}`} />}
      <span className="font-mono text-muted-foreground">{name}</span>
      {comment && <span className="ml-auto pl-3 text-[11px] text-muted-foreground/50">{comment}</span>}
    </div>
  )
}

export function Folder({
  name,
  defaultOpen = false,
  comment,
  children,
}: {
  name: string
  defaultOpen?: boolean
  comment?: string
  children?: ReactNode
}) {
  const [open, setOpen] = useState(defaultOpen)
  return (
    <div>
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center gap-1.5 rounded-md py-1 pl-1 pr-2 text-left hover:bg-accent/20"
      >
        <ChevronRight
          className={`h-3 w-3 shrink-0 text-muted-foreground/50 transition-transform ${open ? 'rotate-90' : ''}`}
        />
        <span className="text-[13px] leading-none">{folderEmoji(open)}</span>
        <span className="font-mono font-medium text-foreground/90">{name}</span>
        {comment && <span className="ml-auto pl-3 text-[11px] text-muted-foreground/50">{comment}</span>}
      </button>
      {children && (
        <div className={`ml-[9px] border-l border-border/50 pl-3 ${open ? 'block' : 'hidden'}`}>
          {children}
        </div>
      )}
    </div>
  )
}
