import type { Metadata } from 'next'
import { ThemeProvider } from '@/components/theme-provider'
import '../globals.css'

export const metadata: Metadata = {
  metadataBase: new URL('https://ui.gritframework.dev'),
  title: {
    default: 'Grit UI — React UI blocks built with Tailwind CSS',
    template: '%s · Grit UI',
  },
  description:
    'Professionally designed, fully responsive React components you can drop into your Tailwind projects. Marketing sections, application UI, and complete page examples. MIT licensed.',
  openGraph: {
    title: 'Grit UI — React UI blocks built with Tailwind CSS',
    description:
      'Marketing sections, application UI, and complete page examples. Copy the code or install with one command.',
    url: 'https://ui.gritframework.dev',
    siteName: 'Grit UI',
    type: 'website',
  },
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    // suppressHydrationWarning is required by next-themes: it sets the class on
    // <html> before React hydrates, which React would otherwise flag.
    <html lang="en" suppressHydrationWarning>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="bg-white font-sans text-gray-900 antialiased dark:bg-gray-950 dark:text-gray-100">
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          enableSystem={false}
          disableTransitionOnChange
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  )
}
