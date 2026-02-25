import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/design/theme')

export default function ThemePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Design System</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Theme &amp; Colors
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit ships with a premium dark theme inspired by Linear, Vercel Dashboard, and Raycast.
                Every component, from login pages to admin tables, is designed to look polished out of the box.
              </p>
            </div>

            <div className="prose-grit">
              {/* Design Philosophy */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Design Philosophy
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit&apos;s visual identity is built on three principles:
                </p>
                <div className="space-y-4">
                  {[
                    {
                      title: 'Dark by Default',
                      desc: 'The dark theme is not an afterthought — it is the primary design surface. Every color, shadow, and border is tuned for dark backgrounds. Light mode is available as an option, but dark mode is the identity of Grit.',
                    },
                    {
                      title: 'Premium CRM Aesthetic',
                      desc: 'Think Linear, Vercel, or Raycast — not generic Bootstrap or Material Design. Subtle gradients, muted borders, careful spacing, and purposeful use of color create a sense of quality.',
                    },
                    {
                      title: 'Functional Beauty',
                      desc: 'Every design choice serves usability. Color contrast meets WCAG guidelines. Spacing provides visual hierarchy. The UI looks good because it works well.',
                    },
                  ].map((item) => (
                    <div key={item.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Color System */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Color System
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses CSS custom properties (design tokens) for all colors. This makes the theme
                  fully customizable and enables dark/light mode switching.
                </p>

                {/* Background Colors */}
                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Background Colors
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Layered backgrounds create depth and visual hierarchy:
                </p>
                <div className="space-y-2 mb-6">
                  {[
                    { token: '--bg-primary', hex: '#0a0a0f', desc: 'Page background, base layer' },
                    { token: '--bg-secondary', hex: '#111118', desc: 'Card backgrounds, sidebar' },
                    { token: '--bg-tertiary', hex: '#1a1a24', desc: 'Elevated cards, modals' },
                    { token: '--bg-elevated', hex: '#22222e', desc: 'Elevated surfaces, dropdowns' },
                    { token: '--bg-hover', hex: '#2a2a38', desc: 'Hover states, active items' },
                    { token: '--border', hex: '#2a2a3a', desc: 'Borders, dividers, separators' },
                  ].map((color) => (
                    <div key={color.token} className="flex items-center gap-4 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div
                        className="h-10 w-10 rounded-lg border border-border/40 shrink-0"
                        style={{ backgroundColor: color.hex }}
                      />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-0.5">
                          <code className="text-xs font-mono text-primary/80">{color.token}</code>
                          <code className="text-[10px] font-mono text-muted-foreground/50">{color.hex}</code>
                        </div>
                        <p className="text-xs text-muted-foreground/60">{color.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Text Colors */}
                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Text Colors
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Three tiers of text prominence:
                </p>
                <div className="space-y-2 mb-6">
                  {[
                    { token: '--text-primary', hex: '#e8e8f0', desc: 'Headings, primary content, important text' },
                    { token: '--text-secondary', hex: '#9090a8', desc: 'Body text, descriptions, labels' },
                    { token: '--text-muted', hex: '#606078', desc: 'Timestamps, placeholders, subtle info' },
                  ].map((color) => (
                    <div key={color.token} className="flex items-center gap-4 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div
                        className="h-10 w-10 rounded-lg border border-border/40 shrink-0 flex items-center justify-center"
                        style={{ backgroundColor: '#0a0a0f' }}
                      >
                        <span className="text-xs font-bold" style={{ color: color.hex }}>Aa</span>
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-0.5">
                          <code className="text-xs font-mono text-primary/80">{color.token}</code>
                          <code className="text-[10px] font-mono text-muted-foreground/50">{color.hex}</code>
                        </div>
                        <p className="text-xs text-muted-foreground/60">{color.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Accent Color */}
                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Accent Color
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The accent color is used for interactive elements, links, buttons, and focus rings:
                </p>
                <div className="space-y-2 mb-6">
                  {[
                    { token: '--accent', hex: '#6c5ce7', desc: 'Primary accent — buttons, links, active states' },
                    { token: '--accent-hover', hex: '#7c6cf7', desc: 'Hover state for accent elements' },
                  ].map((color) => (
                    <div key={color.token} className="flex items-center gap-4 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div
                        className="h-10 w-10 rounded-lg border border-border/40 shrink-0"
                        style={{ backgroundColor: color.hex }}
                      />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-0.5">
                          <code className="text-xs font-mono text-primary/80">{color.token}</code>
                          <code className="text-[10px] font-mono text-muted-foreground/50">{color.hex}</code>
                        </div>
                        <p className="text-xs text-muted-foreground/60">{color.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Semantic Colors */}
                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Semantic Colors
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Used for status indicators, alerts, and feedback:
                </p>
                <div className="space-y-2 mb-6">
                  {[
                    { token: '--success', hex: '#00b894', desc: 'Success states, positive actions, active badges' },
                    { token: '--danger', hex: '#ff6b6b', desc: 'Errors, destructive actions, delete buttons' },
                    { token: '--warning', hex: '#fdcb6e', desc: 'Warnings, caution notices, pending states' },
                    { token: '--info', hex: '#74b9ff', desc: 'Informational messages, tips, neutral badges' },
                  ].map((color) => (
                    <div key={color.token} className="flex items-center gap-4 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div
                        className="h-10 w-10 rounded-lg border border-border/40 shrink-0"
                        style={{ backgroundColor: color.hex }}
                      />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-0.5">
                          <code className="text-xs font-mono text-primary/80">{color.token}</code>
                          <code className="text-[10px] font-mono text-muted-foreground/50">{color.hex}</code>
                        </div>
                        <p className="text-xs text-muted-foreground/60">{color.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Full palette visual */}
                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Full Palette
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden p-5">
                  <div className="grid grid-cols-6 gap-2 mb-4">
                    {[
                      { hex: '#0a0a0f', label: 'bg-primary' },
                      { hex: '#111118', label: 'bg-secondary' },
                      { hex: '#1a1a24', label: 'bg-tertiary' },
                      { hex: '#22222e', label: 'bg-elevated' },
                      { hex: '#2a2a38', label: 'bg-hover' },
                      { hex: '#2a2a3a', label: 'border' },
                    ].map((c) => (
                      <div key={c.label} className="text-center">
                        <div
                          className="h-12 rounded-lg border border-border/40 mb-1.5"
                          style={{ backgroundColor: c.hex }}
                        />
                        <span className="text-[9px] font-mono text-muted-foreground/50 leading-none">{c.label}</span>
                      </div>
                    ))}
                  </div>
                  <div className="grid grid-cols-6 gap-2">
                    {[
                      { hex: '#6c5ce7', label: 'accent' },
                      { hex: '#00b894', label: 'success' },
                      { hex: '#ff6b6b', label: 'danger' },
                      { hex: '#fdcb6e', label: 'warning' },
                      { hex: '#74b9ff', label: 'info' },
                      { hex: '#e8e8f0', label: 'text' },
                    ].map((c) => (
                      <div key={c.label} className="text-center">
                        <div
                          className="h-12 rounded-lg border border-border/40 mb-1.5"
                          style={{ backgroundColor: c.hex }}
                        />
                        <span className="text-[9px] font-mono text-muted-foreground/50 leading-none">{c.label}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Typography */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Typography
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses two fonts that complement each other:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Font</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Usage</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Weights</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium text-foreground/90">Onest</td>
                        <td className="px-4 py-2.5">UI elements: headings, body text, labels, buttons</td>
                        <td className="px-4 py-2.5 font-mono text-xs">400, 500, 600, 700</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium text-foreground/90">JetBrains Mono</td>
                        <td className="px-4 py-2.5">Code blocks, terminal output, monospace labels</td>
                        <td className="px-4 py-2.5 font-mono text-xs">400, 500, 600</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Font Loading
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Fonts are loaded via Google Fonts in the root layout. Next.js automatically optimizes
                  font loading with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next/font</code>:
                </p>
                <CodeBlock language="tsx" filename="apps/web/app/layout.tsx" code={`import { Onest, JetBrains_Mono } from 'next/font/google'

const onest = Onest({
  subsets: ['latin'],
  variable: '--font-onest',
  weight: ['400', '500', '600', '700'],
})

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  variable: '--font-mono',
  weight: ['400', '500', '600'],
})`} />
              </div>

              {/* Using the Theme in CSS */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Using the Theme
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The theme tokens are available as Tailwind CSS classes via shadcn/ui&apos;s configuration.
                  Use them like any Tailwind utility:
                </p>
                <CodeBlock filename="Tailwind usage examples" code={`{/* Background colors */}
<div className="bg-background">       {/* --bg-primary */}
<div className="bg-card">             {/* --bg-secondary */}
<div className="bg-accent">           {/* --bg-tertiary */}

{/* Text colors */}
<h1 className="text-foreground">      {/* --text-primary */}
<p className="text-muted-foreground"> {/* --text-secondary */}

{/* Accent / primary color */}
<button className="bg-primary text-primary-foreground">
<a className="text-primary hover:text-primary/80">

{/* Borders */}
<div className="border border-border">
<hr className="border-border/30">`} />
              </div>

              {/* Dark/Light Toggle */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Dark / Light Mode Toggle
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The admin panel includes a theme toggle powered by{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">next-themes</code>.
                  The toggle switches between dark and light modes by updating the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">class</code> attribute
                  on the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&lt;html&gt;</code> element.
                </p>
                <CodeBlock filename="Theme toggle component" code={`'use client'

import { useTheme } from 'next-themes'
import { Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
    >
      <Sun className="h-4 w-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <Moon className="absolute h-4 w-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
    </Button>
  )
}`} />
                <p className="text-sm text-muted-foreground/60">
                  The theme preference is stored in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localStorage</code> and
                  persists across sessions. The default is <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">dark</code>.
                </p>
              </div>

              {/* Customizing the Theme */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Customizing the Theme
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  All theme tokens are defined in your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">globals.css</code> file.
                  To customize the theme, update the CSS custom properties:
                </p>
                <CodeBlock language="css" filename="app/globals.css" code={`@layer base {
  :root {
    /* Light mode tokens */
    --background: 0 0% 100%;
    --foreground: 240 10% 3.9%;
    --card: 0 0% 97%;
    --primary: 262 52% 63%;      /* Purple accent */
    --primary-foreground: 0 0% 100%;
    /* ... */
  }

  .dark {
    /* Dark mode tokens (Grit defaults) */
    --background: 240 33% 4%;     /* #0a0a0f */
    --foreground: 240 14% 93%;    /* #e8e8f0 */
    --card: 240 18% 8%;           /* #111118 */
    --primary: 252 75% 63%;       /* #6c5ce7 */
    --primary-foreground: 0 0% 100%;
    --muted: 240 10% 42%;         /* #606078 */
    --muted-foreground: 248 12% 61%; /* #9090a8 */
    --border: 240 17% 20%;        /* #2a2a3a */
    --accent: 240 18% 12%;        /* #1a1a24 */
    /* ... */
  }
}`} />
                <p className="text-muted-foreground leading-relaxed mb-4">
                  To change the accent color from purple to blue, for example, update the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--primary</code> value
                  in both light and dark sections:
                </p>
                <CodeBlock filename="Custom blue accent" code={`.dark {
    --primary: 217 91% 60%;   /* #3b82f6 — Tailwind blue-500 */
}`} />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Tip:</strong> Use the{' '}
                    <Link href="https://ui.shadcn.com/themes" target="_blank" className="text-primary hover:underline">
                      shadcn/ui theme builder
                    </Link>{' '}
                    to visually design your custom theme, then paste the generated CSS into your{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">globals.css</code>.
                    Grit is fully compatible with any shadcn/ui theme.
                  </p>
                </div>
              </div>

              {/* Design Tokens Reference */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Quick Reference
                </h2>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Token</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Hex</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Tailwind Class</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      {[
                        { token: 'Background', hex: '#0a0a0f', tw: 'bg-background' },
                        { token: 'Card', hex: '#111118', tw: 'bg-card' },
                        { token: 'Accent surface', hex: '#1a1a24', tw: 'bg-accent' },
                        { token: 'Primary (purple)', hex: '#6c5ce7', tw: 'bg-primary / text-primary' },
                        { token: 'Foreground', hex: '#e8e8f0', tw: 'text-foreground' },
                        { token: 'Muted foreground', hex: '#9090a8', tw: 'text-muted-foreground' },
                        { token: 'Border', hex: '#2a2a3a', tw: 'border-border' },
                        { token: 'Destructive', hex: '#ff6b6b', tw: 'bg-destructive' },
                      ].map((item) => (
                        <tr key={item.token} className="border-b border-border/20 last:border-b-0">
                          <td className="px-4 py-2.5">{item.token}</td>
                          <td className="px-4 py-2.5">
                            <div className="flex items-center gap-2">
                              <div
                                className="h-3.5 w-3.5 rounded-sm border border-border/40"
                                style={{ backgroundColor: item.hex }}
                              />
                              <code className="font-mono text-xs">{item.hex}</code>
                            </div>
                          </td>
                          <td className="px-4 py-2.5 font-mono text-xs text-primary/70">{item.tw}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/infrastructure/deployment" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Deployment
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/tutorials/blog" className="gap-1.5">
                    Tutorials
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
