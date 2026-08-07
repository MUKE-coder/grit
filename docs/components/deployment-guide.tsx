import { Info } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'
import { GuideText } from '@/components/guide-text'
import type { DeploymentGuide, GuideBlock } from '@/config/deployment-guides'

/**
 * Renders a full deployment walkthrough from config/deployment-guides.ts.
 *
 * The guides are long on purpose — they are transcriptions of a deploy someone
 * actually did — so the job here is to keep them navigable: every `Part` is an
 * <h2> with an id, so the page's table of contents can link into them and a
 * reader can send a colleague to a step rather than to the page.
 *
 * Tables scroll inside their own container. A Compose-to-platform mapping is
 * three columns of long strings, and letting it set the page width means every
 * paragraph above it inherits a horizontal scrollbar on a phone.
 */

function slugify(text: string) {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

function Blocks({ blocks }: { blocks: GuideBlock[] }) {
  return (
    <>
      {blocks.map((block, i) => {
        switch (block.kind) {
          case 'p':
            return (
              <p key={i} className="text-sm text-muted-foreground leading-relaxed mb-4">
                <GuideText text={block.text} />
              </p>
            )

          case 'h3':
            return (
              <h3 key={i} className="text-lg font-semibold tracking-tight mt-8 mb-3">
                <GuideText text={block.text} />
              </h3>
            )

          case 'code':
            return (
              <div key={i} className="mb-4">
                <CodeBlock language={block.language} code={block.code} />
              </div>
            )

          case 'ul':
            return (
              <ul key={i} className="mb-4 space-y-2.5 pl-5 list-disc marker:text-muted-foreground/50">
                {block.items.map((item) => (
                  <li key={item} className="text-sm text-muted-foreground leading-relaxed">
                    <GuideText text={item} />
                  </li>
                ))}
              </ul>
            )

          case 'ol':
            return (
              <ol key={i} className="mb-5 space-y-4">
                {block.items.map((item, n) => (
                  <li key={item.text} className="relative pl-9">
                    <span
                      aria-hidden="true"
                      className="absolute left-0 top-0 flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-[11px] font-bold text-primary tabular-nums"
                    >
                      {n + 1}
                    </span>
                    <div className="text-sm text-muted-foreground leading-relaxed">
                      <GuideText text={item.text} />
                    </div>
                    {item.sub && (
                      <ul className="mt-2 space-y-1.5 pl-5 list-disc marker:text-muted-foreground/50">
                        {item.sub.map((s) => (
                          <li key={s} className="text-sm text-muted-foreground leading-relaxed">
                            <GuideText text={s} />
                          </li>
                        ))}
                      </ul>
                    )}
                    {item.code && (
                      <div className="mt-3">
                        <CodeBlock language={item.code.language} code={item.code.code} />
                      </div>
                    )}
                  </li>
                ))}
              </ol>
            )

          case 'note':
            return (
              <div
                key={i}
                className="mb-5 flex gap-3 rounded-xl border border-amber-500/25 bg-amber-500/[0.04] p-4"
              >
                <Info className="h-4 w-4 shrink-0 text-amber-500 mt-0.5" />
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <GuideText text={block.text} />
                </p>
              </div>
            )

          case 'table':
            return (
              /* Scrolls inside itself, so a wide mapping table does not give the
                 whole page a horizontal scrollbar on a phone. */
              <div
                key={i}
                className="mb-5 overflow-x-auto rounded-xl border border-border/50"
              >
                <table className="w-full text-left text-sm">
                  <thead className="bg-card/60">
                    <tr>
                      {block.headers.map((h) => (
                        <th
                          key={h}
                          scope="col"
                          className="whitespace-nowrap px-4 py-3 font-semibold"
                        >
                          <GuideText text={h} />
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/40">
                    {block.rows.map((row) => (
                      <tr key={row.join('|')}>
                        {row.map((cell, c) => (
                          <td
                            key={c}
                            className={`px-4 py-3 align-top text-muted-foreground ${
                              c === 0 ? 'font-medium text-foreground' : ''
                            }`}
                          >
                            <GuideText text={cell} />
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )
        }
      })}
    </>
  )
}

export function DeploymentGuideBody({ guide }: { guide: DeploymentGuide }) {
  return (
    <div>
      {guide.intro.map((p) => (
        <p key={p} className="text-sm text-muted-foreground leading-relaxed mb-4">
          <GuideText text={p} />
        </p>
      ))}

      {guide.sections.map((section) => (
        <section key={section.heading} className="mt-12">
          <h2
            id={slugify(section.heading)}
            className="scroll-mt-24 text-2xl font-bold tracking-tight mb-4"
          >
            <GuideText text={section.heading} />
          </h2>
          <Blocks blocks={section.blocks} />
        </section>
      ))}
    </div>
  )
}
