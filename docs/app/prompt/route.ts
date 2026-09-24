import { BUILD_WITH_AI_PROMPT } from '@/config/build-prompt'

/**
 * The build prompt as plain text, at https://gritframework.dev/prompt
 *
 * Three jobs. It is the fallback when a browser blocks the clipboard, it is
 * something an agent can fetch rather than be pasted, and it is a URL a person
 * can say out loud. Served from the same constant the button copies, so the
 * two cannot disagree.
 */
export const dynamic = 'force-static'

export function GET() {
  return new Response(BUILD_WITH_AI_PROMPT, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'Cache-Control': 'public, max-age=300, s-maxage=3600',
    },
  })
}
