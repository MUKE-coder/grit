'use client'

import { NexoraChat as Chat } from '@nexora-ai/react'

export function NexoraChat() {
  return (
    <Chat
      apiUrl="https://nexoraapi.gritcms.com"
      agentId="6e8794f8-f318-430f-a73f-ba347fb11bb4"
    />
  )
}
