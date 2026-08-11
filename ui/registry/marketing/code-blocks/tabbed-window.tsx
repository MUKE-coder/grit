'use client'

import { useId, useState } from 'react'

/*
 * A code window with a tab per language, and line numbers.
 *
 * On highlighting: this tokenises with about thirty lines of regex rather than
 * pulling in Shiki or Prism. A marketing page showing eight lines of sample
 * code does not need a full grammar, and the smallest real highlighter is a
 * larger download than the rest of the page put together. The trade is honest
 * and worth stating: this will mis-colour something eventually. It is for
 * samples you control, not for user-submitted code.
 *
 * The tabs are real buttons in a `role="tablist"`, wired with `aria-selected`
 * and `aria-controls`, and the inactive panels are removed from the DOM rather
 * than hidden with CSS. Left and right arrows move between tabs, which is what
 * the pattern requires and what a row of divs with onClick does not do.
 *
 * The line numbers live in an `aria-hidden` column, not in the text. Put them
 * in the text and they are copied along with the code, which turns a
 * copy-paste into a syntax error.
 *
 * The <pre> scrolls on its own with `overflow-x-auto`, so a long line moves
 * inside the window instead of pushing the page sideways.
 */

type Kind = 'plain' | 'keyword' | 'string' | 'number' | 'comment' | 'fn' | 'punct'

const KEYWORDS = new Set([
  'const', 'let', 'var', 'function', 'return', 'await', 'async', 'import', 'from',
  'export', 'default', 'if', 'else', 'for', 'while', 'new', 'class', 'type',
  'interface', 'package', 'func', 'err', 'nil', 'range', 'true', 'false', 'null',
])

const TONE: Record<Kind, string> = {
  plain: 'text-gray-800 dark:text-gray-200',
  keyword: 'text-rose-700 dark:text-rose-400',
  string: 'text-emerald-700 dark:text-emerald-400',
  number: 'text-amber-700 dark:text-amber-400',
  comment: 'text-gray-500 dark:text-gray-400',
  fn: 'text-indigo-600 dark:text-indigo-400',
  punct: 'text-gray-500 dark:text-gray-400',
}

/* Deliberately small. Good enough for a fixed sample, not an editor. */
function tokenize(line: string): { text: string; kind: Kind }[] {
  const pattern =
    /(\/\/.*$)|('[^']*'|"[^"]*"|`[^`]*`)|(\b\d+(?:\.\d+)?\b)|([A-Za-z_$][\w$]*)(?=\s*\()|([A-Za-z_$][\w$]*)|([{}()[\].,;:=+\-*/<>!&|?]+)|(\s+)/g
  const out: { text: string; kind: Kind }[] = []
  let match: RegExpExecArray | null
  let last = 0
  while ((match = pattern.exec(line))) {
    if (match.index > last) out.push({ text: line.slice(last, match.index), kind: 'plain' })
    const [text, comment, str, num, fn, word, punct] = match
    let kind: Kind = 'plain'
    if (comment) kind = 'comment'
    else if (str) kind = 'string'
    else if (num) kind = 'number'
    else if (fn) kind = 'fn'
    else if (word) kind = KEYWORDS.has(word) ? 'keyword' : 'plain'
    else if (punct) kind = 'punct'
    out.push({ text, kind })
    last = match.index + text.length
  }
  if (last < line.length) out.push({ text: line.slice(last), kind: 'plain' })
  return out
}

export interface Tab {
  label: string
  code: string
}

const TABS: Tab[] = [
  {
    label: 'TypeScript',
    code: `import { createClient } from '@/lib/api'

const client = createClient({
  baseUrl: process.env.NEXT_PUBLIC_API_URL,
})

// Typed from the Go struct, so a renamed field fails the build
const { data } = await client.posts.list({ page: 1 })

console.log(data.items.length)`,
  },
  {
    label: 'Go',
    code: `package handlers

// The service holds the logic; the handler stays thin
func (h *PostHandler) List(c *gin.Context) {
  posts, meta, err := h.service.List(c.Request.Context())
  if err != nil {
    respond.Error(c, err)
    return
  }

  respond.Paginated(c, posts, meta)
}`,
  },
  {
    label: 'Shell',
    code: `# One command scaffolds the whole stack
grit new storefront --api

cd storefront
grit generate resource Post title:string body:text
grit dev`,
  },
]

export default function CodeBlockTabbedWindow({ tabs = TABS }: { tabs?: Tab[] }) {
  const [active, setActive] = useState(0)
  const id = useId()
  const lines = tabs[active].code.split('\n')

  /* Arrow keys move between tabs; this is required for the pattern. */
  function onKeyDown(event: React.KeyboardEvent) {
    if (event.key !== 'ArrowRight' && event.key !== 'ArrowLeft') return
    event.preventDefault()
    const next =
      event.key === 'ArrowRight'
        ? (active + 1) % tabs.length
        : (active - 1 + tabs.length) % tabs.length
    setActive(next)
    document.getElementById(`${id}-tab-${next}`)?.focus()
  }

  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-4xl px-6 lg:px-8">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-white/10 dark:bg-gray-900">
          <div className="flex items-center gap-1.5 px-4 pt-4">
            {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
              <span key={tone} aria-hidden="true" className={`size-2.5 rounded-full ${tone}`} />
            ))}
          </div>

          <div role="tablist" aria-label="Language" className="flex gap-1 px-4 pt-4">
            {tabs.map((tab, i) => (
              <button
                key={tab.label}
                id={`${id}-tab-${i}`}
                type="button"
                role="tab"
                aria-selected={i === active}
                aria-controls={`${id}-panel-${i}`}
                tabIndex={i === active ? 0 : -1}
                onClick={() => setActive(i)}
                onKeyDown={onKeyDown}
                className={`min-h-11 rounded-lg px-3 text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 ${
                  i === active
                    ? 'bg-gray-100 text-gray-900 dark:bg-white/10 dark:text-white'
                    : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>

          <div
            id={`${id}-panel-${active}`}
            role="tabpanel"
            aria-labelledby={`${id}-tab-${active}`}
            tabIndex={0}
            className="focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600"
          >
            {/* The pre scrolls, not the page. */}
            <pre className="overflow-x-auto p-4 font-mono text-[13px]/6">
              <code>
                {lines.map((line, i) => (
                  <span key={i} className="grid grid-cols-[2.5rem_1fr]">
                    {/* Numbers are outside the text, so copying the code does
                        not copy them into the paste. */}
                    <span
                      aria-hidden="true"
                      className="pr-4 text-right text-gray-300 select-none dark:text-gray-600"
                    >
                      {i + 1}
                    </span>
                    <span>
                      {tokenize(line).map((token, j) => (
                        <span key={j} className={TONE[token.kind]}>
                          {token.text}
                        </span>
                      ))}
                    </span>
                  </span>
                ))}
              </code>
            </pre>
          </div>
        </div>
      </div>
    </section>
  )
}
