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
      { src: '/favicon.ico',                sizes: 'any',     type: 'image/x-icon' },
      { src: '/favicon-16x16.png',          sizes: '16x16',   type: 'image/png' },
      { src: '/favicon-32x32.png',          sizes: '32x32',   type: 'image/png' },
      { src: '/apple-touch-icon.png',       sizes: '180x180', type: 'image/png' },
      { src: '/android-chrome-192x192.png', sizes: '192x192', type: 'image/png', purpose: 'any' },
      { src: '/android-chrome-512x512.png', sizes: '512x512', type: 'image/png', purpose: 'any' },
    ],
  }
}
