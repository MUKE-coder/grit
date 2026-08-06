/*
 * A JSON response in a named file tab, scrolling inside a fixed height.
 *
 * The height is capped and the body scrolls, which is the point: an API
 * response is long and boring past the first ten lines, and letting it run to
 * its full length pushes everything below it off the screen. Capping it says
 * "this continues" without spending the page on it.
 *
 * The scroll container is `tabIndex={0}` with a label. A scrollable region that
 * cannot be focused cannot be scrolled from the keyboard at all — you can see
 * the content and have no way to reach it. This is the single most common
 * accessibility bug in a code block, and it is one attribute.
 *
 * Highlighting is a small hand-rolled tokeniser rather than Shiki or Prism.
 * For JSON the grammar is small enough that this is genuinely sufficient, and
 * the smallest real highlighter outweighs the rest of the page.
 *
 * Line numbers sit in an `aria-hidden` column outside the text, so copying the
 * response does not copy the numbers into the paste.
 */

const RESPONSE = `{
  "users": [
    {
      "name": "John Doe",
      "email": "john.doe@example.com",
      "age": 30,
      "verified": true,
      "cart": [
        { "id": 1, "name": "Product 1", "price": 10 },
        { "id": 2, "name": "Product 2", "price": 24 }
      ]
    },
    {
      "name": "Ada Byron",
      "email": "ada@example.com",
      "age": 36,
      "verified": false,
      "cart": []
    }
  ],
  "meta": { "total": 2, "page": 1, "pages": 1 }
}`

type Kind = 'key' | 'string' | 'number' | 'literal' | 'punct'

const TONE: Record<Kind, string> = {
  key: 'text-sky-600 dark:text-sky-400',
  string: 'text-emerald-600 dark:text-emerald-400',
  number: 'text-amber-600 dark:text-amber-400',
  literal: 'text-rose-600 dark:text-rose-400',
  punct: 'text-gray-400 dark:text-gray-500',
}

/* A key is a string followed by a colon; everything else is a value. */
function tokenize(line: string): { text: string; kind: Kind }[] {
  const pattern = /("(?:[^"\\]|\\.)*")(\s*:)?|(-?\d+(?:\.\d+)?)|\b(true|false|null)\b|([^\s"]+)|(\s+)/g
  const out: { text: string; kind: Kind }[] = []
  let match: RegExpExecArray | null
  while ((match = pattern.exec(line))) {
    const [text, str, colon, num, literal] = match
    if (str) {
      out.push({ text: str, kind: colon ? 'key' : 'string' })
      if (colon) out.push({ text: colon, kind: 'punct' })
    } else if (num) out.push({ text, kind: 'number' })
    else if (literal) out.push({ text, kind: 'literal' })
    else out.push({ text, kind: 'punct' })
  }
  return out
}

export default function CodeBlockJsonViewer({
  filename = 'response.json',
  json = RESPONSE,
}: {
  filename?: string
  json?: string
}) {
  const lines = json.split('\n')

  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-3xl px-6 lg:px-8">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-white/10 dark:bg-gray-900">
          <div className="flex items-center border-b border-gray-200 dark:border-white/10">
            <span className="flex items-center gap-2 border-r border-gray-200 px-4 py-3 font-mono text-xs text-gray-700 dark:border-white/10 dark:text-gray-300">
              <span aria-hidden="true" className="text-amber-500">
                {'{}'}
              </span>
              {filename}
            </span>
          </div>

          {/* tabIndex makes this reachable from the keyboard. Without it the
              content below the fold cannot be scrolled to at all. */}
          <div
            role="region"
            aria-label={`${filename} contents`}
            tabIndex={0}
            className="max-h-96 overflow-auto focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600"
          >
            <pre className="p-4 font-mono text-[13px]/6">
              <code>
                {lines.map((line, i) => (
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
      </div>
    </section>
  )
}
