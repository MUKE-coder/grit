'use client'

import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { CodeBlock } from '@/components/code-block'

// Renders a blog post's markdown body. Fenced code blocks are routed through the
// shared <CodeBlock> so they get Prism syntax highlighting + a copy button; the
// rest (headings, paragraphs, lists, tables, inline code) is styled by prose-grit.

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
          // <CodeBlock> renders its own <pre>, so unwrap react-markdown's.
          pre: ({ children }) => <>{children}</>,
          code({ className, children }) {
            const match = /language-(\w+)/.exec(className || '')
            const text = String(children).replace(/\n$/, '')
            const isBlock = Boolean(match) || text.includes('\n')

            if (isBlock) {
              const lang = match ? LANG_ALIAS[match[1]] || match[1] : 'bash'
              return <CodeBlock code={text} language={lang} className="my-5" />
            }
            // inline code — styled by prose-grit
            return <code>{children}</code>
          },
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  )
}
