import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  metadataBase: new URL('https://ui.gritframework.dev'),
  title: {
    default: 'Grit UI — 100 React components, shadcn-compatible',
    template: '%s · Grit UI',
  },
  description:
    'A registry of 100 production-ready React components for marketing, SaaS, ecommerce, auth and app layout. Install any of them with npx shadcn add. MIT licensed.',
  openGraph: {
    title: 'Grit UI — 100 React components',
    description:
      'Marketing, SaaS, ecommerce, auth and layout components. Install with npx shadcn add. MIT licensed.',
    url: 'https://ui.gritframework.dev',
    siteName: 'Grit UI',
    type: 'website',
  },
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="dark">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Onest:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="font-sans antialiased">{children}</body>
    </html>
  )
}
