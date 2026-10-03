import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { llmsIndex, llmsCliReference } from '@/lib/llms'

export const dynamic = 'force-static'

// The working guide, read from the same file /skill.md serves and `grit init`
// writes into a project. One source: a guide that disagreed with itself
// depending on which URL an agent happened to fetch would be worse than none.
function workingGuide(): string {
  const raw = readFileSync(join(process.cwd(), 'public', 'skill.md'), 'utf8')
  // Strip the YAML front matter, which addresses a skill loader rather than a
  // reader, and leave the body.
  return raw.replace(/^---\n[\s\S]*?\n---\n/, '').trim()
}

export function GET() {
  const body = [llmsIndex(), '---', workingGuide(), '---', llmsCliReference()].join('\n\n')
  return new Response(body, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'Cache-Control': 'public, max-age=0, s-maxage=3600, stale-while-revalidate=86400',
    },
  })
}
