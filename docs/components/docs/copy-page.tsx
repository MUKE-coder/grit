'use client'

import { useCallback, useState } from 'react'
import { Check, Copy } from 'lucide-react'

/**
 * Copy the rendered page as Markdown.
 *
 * The reason this exists is that people paste documentation into a chat with
 * an agent, and what the clipboard gets from selecting a rendered page is
 * HTML-flavoured text with the navigation in it. This walks the content
 * column only, and emits the handful of constructs docs actually use:
 * headings, paragraphs, lists, code blocks, tables and links.
 *
 * It reads the DOM rather than a source file because there is no source
 * file: every docs page here is a React component, so the rendered output is
 * the only place the prose exists in one piece.
 */

/** Turns one element's inline content into Markdown, links included. */
function inline(node: Node): string {
  if (node.nodeType === Node.TEXT_NODE) return node.textContent ?? ''
  if (node.nodeType !== Node.ELEMENT_NODE) return ''

  const el = node as HTMLElement
  if (el.getAttribute('aria-hidden') === 'true') return ''
  if (el.classList.contains('sr-only')) return ''

  const inner = [...el.childNodes].map(inline).join('')

  switch (el.tagName) {
    case 'CODE':
      // A <code> inside a <pre> is handled by the block walker.
      return el.closest('pre') ? inner : `\`${inner}\``
    case 'STRONG':
    case 'B':
      return `**${inner}**`
    case 'EM':
    case 'I':
      return `*${inner}*`
    case 'A': {
      const href = el.getAttribute('href') ?? ''
      if (!href || href.startsWith('#')) return inner
      const abs = href.startsWith('http') ? href : new URL(href, location.origin).href
      return `[${inner}](${abs})`
    }
    case 'BR':
      return '\n'
    default:
      return inner
  }
}

const clean = (s: string) => s.replace(/[ \t]+/g, ' ').replace(/ ?\n ?/g, '\n').trim()

function blocks(root: HTMLElement): string[] {
  const out: string[] = []

  for (const el of [...root.children] as HTMLElement[]) {
    if (el.getAttribute('aria-hidden') === 'true') continue
    if (el.classList.contains('sr-only')) continue
    if (el.tagName === 'NAV' || el.tagName === 'SCRIPT' || el.tagName === 'STYLE') continue

    switch (el.tagName) {
      case 'H1':
      case 'H2':
      case 'H3':
      case 'H4': {
        const depth = Number(el.tagName[1])
        out.push(`${'#'.repeat(depth)} ${clean(inline(el))}`)
        break
      }
      case 'P': {
        const text = clean(inline(el))
        if (text) out.push(text)
        break
      }
      case 'UL':
      case 'OL': {
        const ordered = el.tagName === 'OL'
        const items = [...el.children]
          .filter((li) => li.tagName === 'LI')
          .map((li, i) => `${ordered ? `${i + 1}.` : '-'} ${clean(inline(li))}`)
          .filter((line) => line.length > 2)
        if (items.length) out.push(items.join('\n'))
        break
      }
      case 'PRE': {
        const code = el.textContent ?? ''
        // The language is on the <code> as `language-xxx` where a highlighter
        // put it, and absent where one did not. A fence with no language is
        // still a fence.
        const lang =
          [...(el.querySelector('code')?.classList ?? [])]
            .find((c) => c.startsWith('language-'))
            ?.slice('language-'.length) ?? ''
        out.push(`\`\`\`${lang}\n${code.replace(/\n+$/, '')}\n\`\`\``)
        break
      }
      case 'TABLE': {
        const rows = [...el.querySelectorAll('tr')].map((tr) =>
          [...tr.children].map((cell) => clean(inline(cell)).replace(/\|/g, '\\|'))
        )
        if (rows.length) {
          const [head, ...body] = rows
          out.push(
            [
              `| ${head.join(' | ')} |`,
              `| ${head.map(() => '---').join(' | ')} |`,
              ...body.map((r) => `| ${r.join(' | ')} |`),
            ].join('\n')
          )
        }
        break
      }
      case 'BLOCKQUOTE':
        out.push(
          clean(inline(el))
            .split('\n')
            .map((l) => `> ${l}`)
            .join('\n')
        )
        break
      case 'HR':
        out.push('---')
        break
      default:
        // A wrapper: recurse. Anything with no block children of its own and
        // some text is treated as a paragraph, which is how a card or a
        // callout survives the walk instead of vanishing.
        if (el.children.length) out.push(...blocks(el))
        else {
          const text = clean(inline(el))
          if (text) out.push(text)
        }
    }
  }

  return out
}

export function CopyPageButton({ className = '' }: { className?: string }) {
  const [copied, setCopied] = useState(false)

  const copy = useCallback(async () => {
    const main = document.querySelector('main')
    if (!main) return

    const markdown = [`${location.origin}${location.pathname}`, '', ...blocks(main)]
      .join('\n\n')
      .replace(/\n{3,}/g, '\n\n')
      .trim()

    try {
      await navigator.clipboard.writeText(markdown)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1600)
    } catch {
      // Blocked over http or without a gesture. The page is still readable.
    }
  }, [])

  return (
    <button
      type="button"
      onClick={() => void copy()}
      className={`inline-flex h-7 items-center gap-1.5 rounded-btn-sm border border-border px-2.5 text-[12px] font-medium text-fg-muted transition-colors hover:border-border-strong hover:text-fg ${className}`.trim()}
    >
      {copied ? (
        <Check className="h-3.5 w-3.5 text-success" aria-hidden="true" />
      ) : (
        <Copy className="h-3.5 w-3.5" aria-hidden="true" />
      )}
      {copied ? 'Copied as Markdown' : 'Copy page'}
    </button>
  )
}
