'use client'

import { useId, useState } from 'react'

/*
 * The source on the left, what it renders on the right.
 *
 * The preview is real markup, not a screenshot of one. That matters more than
 * it sounds: a screenshot of a login form is stale the moment you restyle the
 * button, and it is invisible to anyone who cannot see it. This preview
 * inherits the theme and is read out like any other content.
 *
 * The preview is marked `aria-hidden` all the same, and that is a deliberate
 * trade rather than an oversight. It is an illustration of the code beside it,
 * so its inputs and buttons would otherwise land in the tab order as
 * non-functional controls — a focus trap of things that do nothing. If you
 * wire it up for real, drop the attribute.
 *
 * Highlighting is a small hand-rolled tokeniser rather than Shiki or Prism:
 * eight lines of sample code do not justify a highlighter larger than the rest
 * of the page. It will mis-colour something eventually. Use it for samples you
 * control.
 */

type Kind = 'plain' | 'keyword' | 'string' | 'tag' | 'attr' | 'comment' | 'punct'

const KEYWORDS = new Set([
  'import', 'from', 'export', 'default', 'function', 'return', 'const', 'let',
  'async', 'await', 'if', 'else', 'type', 'interface',
])

const TONE: Record<Kind, string> = {
  plain: 'text-gray-800 dark:text-gray-200',
  keyword: 'text-rose-700 dark:text-rose-400',
  string: 'text-emerald-700 dark:text-emerald-400',
  tag: 'text-sky-700 dark:text-sky-400',
  attr: 'text-violet-700 dark:text-violet-400',
  comment: 'text-gray-500 dark:text-gray-400',
  punct: 'text-gray-500 dark:text-gray-400',
}

function tokenize(line: string): { text: string; kind: Kind }[] {
  const pattern =
    /(\/\/.*$|\{\/\*.*?\*\/\})|('[^']*'|"[^"]*"|`[^`]*`)|(<\/?[A-Za-z][\w.]*)|([A-Za-z-]+)(?==)|([A-Za-z_$][\w$]*)|([{}()[\].,;:=+\-*/<>!&|?]+)|(\s+)/g
  const out: { text: string; kind: Kind }[] = []
  let match: RegExpExecArray | null
  let last = 0
  while ((match = pattern.exec(line))) {
    if (match.index > last) out.push({ text: line.slice(last, match.index), kind: 'plain' })
    const [text, comment, str, tag, attr, word, punct] = match
    let kind: Kind = 'plain'
    if (comment) kind = 'comment'
    else if (str) kind = 'string'
    else if (tag) kind = 'tag'
    else if (attr) kind = 'attr'
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
    label: 'Next.js',
    code: `import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export default function LoginPage() {
  return (
    <form action={signIn} className="mx-auto max-w-sm">
      <h1>Welcome back</h1>
      <p>Sign in to continue</p>

      <Label htmlFor="email">Email</Label>
      <Input id="email" name="email" type="email" required />

      <Button type="submit">Continue</Button>
    </form>
  )
}`,
  },
  {
    label: 'Svelte',
    code: `<script lang="ts">
  import Button from '$lib/ui/button.svelte'
  import Input from '$lib/ui/input.svelte'

  let email = ''
</script>

<form method="POST" class="mx-auto max-w-sm">
  <h1>Welcome back</h1>
  <p>Sign in to continue</p>

  <label for="email">Email</label>
  <Input id="email" bind:value={email} type="email" required />

  <Button type="submit">Continue</Button>
</form>`,
  },
]

function Preview() {
  return (
    /* An illustration of the code, not a working form: its controls are kept
       out of the tab order so they are not focusable things that do nothing. */
    <div
      aria-hidden="true"
      className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-gray-900"
    >
      <span className="flex size-8 items-center justify-center rounded-lg bg-indigo-600 text-sm font-semibold text-white">
        G
      </span>
      <h3 className="mt-6 text-lg font-semibold text-gray-900 dark:text-white">Welcome back</h3>
      <p className="text-sm text-gray-600 dark:text-gray-400">Sign in to continue</p>

      <div className="mt-6 space-y-2">
        {['Continue with Google', 'Continue with GitHub'].map((label) => (
          <div
            key={label}
            className="flex min-h-11 items-center justify-center rounded-lg border border-gray-200 text-sm font-medium text-gray-700 dark:border-white/10 dark:text-gray-300"
          >
            {label}
          </div>
        ))}
      </div>

      <div className="my-6 flex items-center gap-3">
        <span className="h-px flex-1 bg-gray-200 dark:bg-white/10" />
        <span className="text-xs text-gray-400">or</span>
        <span className="h-px flex-1 bg-gray-200 dark:bg-white/10" />
      </div>

      <p className="text-sm font-medium text-gray-900 dark:text-white">Email</p>
      <div className="mt-1 flex min-h-11 items-center rounded-lg border border-gray-200 px-3 text-sm text-gray-400 dark:border-white/10 dark:text-gray-600">
        you@example.com
      </div>
      <div className="mt-4 flex min-h-11 items-center justify-center rounded-lg bg-indigo-600 text-sm font-semibold text-white">
        Continue
      </div>
    </div>
  )
}

export default function CodeBlockCodeWithPreview({
  title = 'The same component, in the framework you use',
  subtitle = 'Every block in the registry is plain source you own after installing it. There is no runtime to keep in step.',
  tabs = TABS,
}: {
  title?: string
  subtitle?: string
  tabs?: Tab[]
}) {
  const [active, setActive] = useState(0)
  const id = useId()

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
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <div className="mt-16 grid grid-cols-1 items-start gap-8 lg:grid-cols-[minmax(0,1.4fr)_minmax(0,1fr)]">
          <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-white/10 dark:bg-gray-900">
            <div
              role="tablist"
              aria-label="Framework"
              className="flex gap-1 border-b border-gray-200 px-3 dark:border-white/10"
            >
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
                  className={`-mb-px min-h-11 border-b-2 px-3 text-sm font-medium focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 ${
                    i === active
                      ? 'border-indigo-600 text-gray-900 dark:border-indigo-400 dark:text-white'
                      : 'border-transparent text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
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
              <pre className="overflow-x-auto p-4 font-mono text-[13px]/6">
                <code>
                  {tabs[active].code.split('\n').map((line, i) => (
                    <span key={i} className="grid grid-cols-[2.5rem_1fr]">
                      <span
                        aria-hidden="true"
                        className="pr-4 text-right text-gray-300 select-none dark:text-gray-600"
                      >
                        {i + 1}
                      </span>
                      <span className="whitespace-pre">
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

          <Preview />
        </div>
      </div>
    </section>
  )
}
