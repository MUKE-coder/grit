'use client'

/*
 * Not mounted.
 *
 * The floating "Ask AI" button was taken off the site: it covered the
 * bottom-right of every page, it was the only contrast failure the design
 * audit found on the home page (4.47:1 against a 4.5 minimum), and its
 * indigo glow is not a colour this design system has.
 *
 * The component is kept because bringing it back is one line in
 * app/layout.tsx. If it does come back, it needs a theme that reads from the
 * design tokens, and it belongs inside the command palette rather than
 * floating over the page.
 */

import { NexoraChat as Chat } from '@nexora-ai/react'

export function NexoraChat() {
  return (
    <Chat
      apiUrl="https://nexoraapi.gritcms.com"
      agentId="6e8794f8-f318-430f-a73f-ba347fb11bb4"
    />
  )
}
