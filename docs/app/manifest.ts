import type { MetadataRoute } from 'next'

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: 'Grit — Go + React Full-Stack Framework',
    short_name: 'Grit',
    description:
      'The full-stack meta-framework that combines Go, React, and a Filament-like admin panel.',
    start_url: '/',
    display: 'standalone',
    background_color: '#0a0a0f',
    theme_color: '#6c5ce7',
    icons: [
      {
        src: '/favicon.ico',
        sizes: 'any',
        type: 'image/x-icon',
      },
    ],
  }
}
