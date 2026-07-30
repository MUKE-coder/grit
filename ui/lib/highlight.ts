import { codeToHtml } from 'shiki'

/**
 * Highlights block source with Shiki, at build time.
 *
 * This runs in a server component during prerender, so the browser receives
 * finished HTML and downloads no highlighter. Shiki's grammars are ~1MB —
 * shipping them to highlight a handful of static snippets would cost more than
 * every block on the page combined.
 *
 * Both themes are emitted in one pass: the light theme becomes inline styles
 * and the dark theme becomes CSS variables, which globals.css flips under
 * `.dark`. That means switching theme repaints instantly with no re-highlight
 * and no second copy of the markup.
 */
export async function highlight(code: string, lang = 'tsx'): Promise<string> {
  return codeToHtml(code, {
    lang,
    themes: { light: 'github-light', dark: 'github-dark-default' },
    defaultColor: false,
  })
}

/** Highlights many snippets concurrently. */
export async function highlightAll(
  items: { key: string; code: string; lang?: string }[],
): Promise<Record<string, string>> {
  const done = await Promise.all(
    items.map(async (i) => [i.key, await highlight(i.code, i.lang)] as const),
  )
  return Object.fromEntries(done)
}
