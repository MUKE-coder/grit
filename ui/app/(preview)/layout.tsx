import type { Metadata } from 'next'
import '../globals.css'

export const metadata: Metadata = {
  title: 'Block preview',
  robots: { index: false, follow: false },
}

/**
 * A separate root layout for the preview frames, and the reason it exists is
 * worth stating: there is NO ThemeProvider here.
 *
 * Previously the previews shared the site's root layout, so next-themes ran
 * inside every iframe, read the visitor's stored preference from localStorage
 * (same origin, shared storage) and set `dark` on the frame's <html>. A preview
 * asked for light rendered dark for anyone browsing in dark mode, and no amount
 * of setting the class from the page could win — next-themes just set it back.
 *
 * With the provider gone, the frame's theme is decided by exactly one thing:
 * the `dark` class this layout is given from the search param. Nothing else can
 * touch it.
 */
export default function PreviewRootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="font-sans antialiased">{children}</body>
    </html>
  )
}
