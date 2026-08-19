'use client'

import React from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { CodeBlock } from '@/components/code-block'

// Renders a blog post's markdown body. Fenced code blocks are routed through the
// shared <CodeBlock> so they get Prism syntax highlighting + a copy button; the
// rest (headings, paragraphs, lists, tables, inline code) is styled by prose-grit.

// Heading anchors.
//
// react-markdown emits headings with no id, so a "[Step 4f](#step-4f-...)"
// written inside a post scrolls nowhere. That is not hypothetical: the invoice
// guide has carried a dead one since it was published. rehype-slug would do
// this, and it is a dependency and a lockfile change for ten lines.
//
// Same rule as the docs sidebar's table of contents, so a heading gets the same
// id whichever surface renders it.
function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
}

// A heading's text, flattened out of whatever react-markdown handed us. The
// children of "## Step 4f: `Product` variants" are a string, an element and
// another string, and only the strings carry the words.
function headingText(node: React.ReactNode): string {
  if (typeof node === 'string' || typeof node === 'number') return String(node)
  if (Array.isArray(node)) return node.map(headingText).join('')
  if (React.isValidElement(node)) {
    return headingText((node.props as { children?: React.ReactNode }).children)
  }
  return ''
}

const LANG_ALIAS: Record<string, string> = {
  ts: 'typescript',
  js: 'javascript',
  sh: 'bash',
  shell: 'bash',
  yml: 'yaml',
}

export function BlogContent({ content }: { content: string }) {
  return (
    <div className="prose-grit">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          // Anchor targets, so a cross-reference inside a post lands on the
          // section it names.
          h2: ({ children }) => <h2 id={slugify(headingText(children))}>{children}</h2>,
          h3: ({ children }) => <h3 id={slugify(headingText(children))}>{children}</h3>,
          // <CodeBlock> renders its own <pre>, so unwrap react-markdown's.
          pre: ({ children }) => <>{children}</>,
          code({ className, children, node }) {
            const match = /language-(\w+)/.exec(className || '')
            const text = String(children).replace(/\n$/, '')
            const isBlock = Boolean(match) || text.includes('\n')

            if (isBlock) {
              const lang = match ? LANG_ALIAS[match[1]] || match[1] : 'bash'
              // Everything after the language on the fence line, so a block can
              // say where it lives:
              //
              //   ```tsx title="apps/admin/resources/products/products.custom.tsx"
              //
              // Worth the parsing. A snippet with no path is a snippet the
              // reader has to guess the home of, and in a post that touches
              // four files in a section, that guess is usually wrong.
              const meta = (node as { data?: { meta?: string } } | undefined)?.data?.meta ?? ''
              const filename = /title="([^"]+)"/.exec(meta)?.[1]
              return (
                <CodeBlock code={text} language={lang} filename={filename} className="my-5" />
              )
            }
            // inline code, styled by prose-grit
            return <code>{children}</code>
          },
          // Inline images are nearly always screenshots. A border and a radius
          // stop a light UI capture from bleeding into a light page, and the
          // caption comes free from the alt text you already had to write.
          //
          // Spans, not <figure>/<figcaption>: react-markdown wraps a lone image
          // in a <p>, and a <figure> inside a <p> is invalid, so the browser
          // hoists it out during parsing and React reports a hydration mismatch
          // against the server HTML it never had a chance to match.
          img: ({ src, alt }) => (
            <span className="my-8 block">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={typeof src === 'string' ? src : ''}
                alt={alt ?? ''}
                loading="lazy"
                className="block w-full rounded-xl border border-border/60"
              />
              {alt && (
                <span className="mt-2 block text-center text-[13px] text-muted-foreground/70">
                  {alt}
                </span>
              )}
            </span>
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  )
}
