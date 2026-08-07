import Link from 'next/link'
import { ExternalLink } from 'lucide-react'
import type { ReactNode } from 'react'

/**
 * Renders the small amount of inline markdown the deployment guides use:
 * `code`, **bold**, and [links](https://…).
 *
 * It builds React elements rather than setting innerHTML. The guides are
 * checked into this repo so the input is trusted today, but a component that
 * injects raw HTML is one careless edit away from being an XSS hole, and there
 * is no upside here — three patterns is a twenty-line tokeniser.
 *
 * Code spans are matched FIRST and consumed whole, so `**` inside a backticked
 * span stays literal. Getting that order wrong is how a shell snippet containing
 * asterisks silently turns half a command into bold text.
 *
 * Bold and link contents are re-rendered through this component, because the
 * guides nest the two constantly — "**Back up `/opt/orbita/.env`.**" is one
 * bold span containing a code span. Without the recursion the outer token
 * swallows the inner one and the backticks print literally on the page.
 * Recursion terminates because the delimiters are stripped before the inner
 * pass, so each level has strictly less to match.
 */

/* Order matters twice over. Code is first so `**` inside a backticked span
   stays literal, and bold is before italic so `**x**` is not read as an empty
   emphasis followed by `*x*`. */
const PATTERN = /(`[^`]+`)|(\*\*[^*]+\*\*)|(\[[^\]]+\]\([^)]+\))|(\*[^*\n]+\*)/g

export function GuideText({ text }: { text: string }) {
  const out: ReactNode[] = []
  let last = 0
  let match: RegExpExecArray | null
  const re = new RegExp(PATTERN)

  while ((match = re.exec(text))) {
    if (match.index > last) out.push(text.slice(last, match.index))
    const [token, code, bold, link, italic] = match

    if (code) {
      out.push(
        <code
          key={`${match.index}-c`}
          className="rounded bg-muted px-1.5 py-0.5 text-[0.85em] font-mono"
        >
          {code.slice(1, -1)}
        </code>,
      )
    } else if (bold) {
      out.push(
        <strong key={`${match.index}-b`} className="font-semibold text-foreground">
          <GuideText text={bold.slice(2, -2)} />
        </strong>,
      )
    } else if (link) {
      const label = link.slice(1, link.indexOf(']'))
      const href = link.slice(link.indexOf('(') + 1, -1)
      const external = href.startsWith('http')
      out.push(
        <Link
          key={`${match.index}-l`}
          href={href}
          target={external ? '_blank' : undefined}
          rel={external ? 'noreferrer' : undefined}
          className="inline-flex items-center gap-1 text-primary hover:underline"
        >
          <GuideText text={label} />
          {external && <ExternalLink className="h-3 w-3" />}
        </Link>,
      )
    } else if (italic) {
      out.push(
        <em key={`${match.index}-i`}>
          <GuideText text={italic.slice(1, -1)} />
        </em>,
      )
    }
    last = match.index + token.length
  }

  if (last < text.length) out.push(text.slice(last))
  return <>{out}</>
}
