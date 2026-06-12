import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        SEO is &quot;how do search engines find you&quot;; OG is &quot;how does your
        link preview when shared&quot;. Both are metadata. Next.js makes them
        one-liner additions. This lesson is the checklist.
      </p>

      <h2>The page metadata object</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(marketing)/page.tsx"
        code={`import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Acme — Run your retail business from one place',
  description: 'POS, inventory, reports for shops with 1 to 50 outlets.',
  openGraph: {
    title: 'Acme — Run your retail business from one place',
    description: 'POS, inventory, reports for shops with 1 to 50 outlets.',
    url: 'https://acme.com',
    siteName: 'Acme',
    images: [{ url: 'https://acme.com/og.png', width: 1200, height: 630 }],
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Acme — Run your retail business',
    description: 'POS, inventory, reports.',
    images: ['https://acme.com/og.png'],
  },
}

export default function HomePage() { /* ... */ }`}
      />
      <p>
        Next.js renders these into{' '}
        <code>&lt;meta&gt;</code> tags at build / request time. No work
        beyond declaring the object.
      </p>

      <h2>The shareable OG image</h2>
      <p>
        1200 × 630 PNG. Shows when someone shares your link in Slack /
        Twitter / WhatsApp. Build it once with{' '}
        <a href="https://www.figma.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Figma</a>{' '}
        + export PNG.
      </p>
      <p>
        Or generate dynamically with{' '}
        <code>next/og</code>:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/og/route.tsx"
        code={`import { ImageResponse } from 'next/og'

export async function GET() {
  return new ImageResponse(
    <div tw="flex w-full h-full bg-black text-white items-center justify-center">
      <div tw="text-7xl font-bold">Acme</div>
    </div>,
    { width: 1200, height: 630 },
  )
}`}
      />
      <p>
        Now <code>https://acme.com/og</code> renders a dynamic OG image —
        useful when you want page-specific previews (blog post titles,
        user profiles).
      </p>

      <h2>sitemap.xml</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/app/sitemap.ts"
        code={`import type { MetadataRoute } from 'next'

export default function sitemap(): MetadataRoute.Sitemap {
  return [
    { url: 'https://acme.com',         priority: 1.0, changeFrequency: 'weekly' },
    { url: 'https://acme.com/pricing', priority: 0.9, changeFrequency: 'monthly' },
    { url: 'https://acme.com/blog',    priority: 0.8, changeFrequency: 'weekly' },
    // For dynamic pages, fetch from the API and map them
  ]
}`}
      />
      <p>
        Available at <code>/sitemap.xml</code>. Google reads it to discover
        every public URL. For dynamic blog posts, fetch from your API and
        map each post to an entry.
      </p>

      <h2>robots.txt</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/app/robots.ts"
        code={`import type { MetadataRoute } from 'next'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: { userAgent: '*', allow: '/', disallow: '/dashboard/' },
    sitemap: 'https://acme.com/sitemap.xml',
  }
}`}
      />
      <p>
        Tells crawlers what to index. Always disallow{' '}
        <code>/dashboard</code> and similar logged-in surfaces — they&apos;d
        404 for the bot anyway.
      </p>

      <h2>For the admin app</h2>
      <p>
        <code>apps/admin</code> should NEVER be indexed. Add a global{' '}
        <code>robots</code> meta in its root layout:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/admin/app/layout.tsx"
        code={`export const metadata: Metadata = {
  robots: { index: false, follow: false },
}`}
      />

      <TipBox tone="success">
        <strong>Validate with the Facebook OG Debugger.</strong> Visit{' '}
        <a href="https://developers.facebook.com/tools/debug/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">developers.facebook.com/tools/debug</a>,{' '}
        paste your URL. You see exactly what the preview will look like.
        Same for{' '}
        <a href="https://cards-dev.twitter.com/validator" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Twitter Card Validator</a>.
      </TipBox>

      <KnowledgeCheck
        question="You ship your landing page with OG metadata. You share the link in Slack but the preview doesn't render. What's the most likely cause?"
        choices={[
          {
            label: 'Slack caches old previews — re-share or use the debugger to refresh',
            correct: true,
            feedback:
              "Right — Slack caches OG fetches aggressively. The first share with broken metadata sticks. Use the Facebook debugger to force a re-scrape, then re-share.",
          },
          {
            label: 'Slack doesn\'t support OG',
            feedback:
              "Wrong — Slack reads OG tags and renders previews. Cache is the usual culprit.",
          },
          {
            label: 'You forgot to add `og:image`',
            feedback:
              "Would explain a missing image but not a missing preview entirely. The metadata block in the lesson includes the image.",
          },
          {
            label: 'Slack requires HTTPS, you used HTTP',
            feedback:
              "Slack does prefer HTTPS but isn't strict. The cache issue is the more common cause.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              For chapter 2&apos;s assignment, complete the SEO + OG pass on
              your scaffolded web app:
            </p>
            <ol>
              <li>
                Add a metadata object to{' '}
                <code>apps/web/app/(marketing)/page.tsx</code>.
              </li>
              <li>
                Add <code>app/sitemap.ts</code> + <code>app/robots.ts</code>.
              </li>
              <li>
                Add a static OG image (or generate via <code>next/og</code>).
              </li>
              <li>
                Deploy to a public URL (Vercel, your VPS) and run the OG
                tester.
              </li>
              <li>
                Run Lighthouse — SEO score should be 100.
              </li>
            </ol>
            <p>Paste the OG preview + Lighthouse SEO in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            For local-only testing (no deploy yet), use{' '}
            <code>ngrok</code> to expose <code>localhost:3000</code> publicly.
            Then point the OG tester at the ngrok URL.
          </>
        }
        solution={
          <>
            <p>A healthy Slack preview looks like:</p>
            <ul>
              <li>Title</li>
              <li>Description</li>
              <li>1200 × 630 image</li>
              <li>Site name + URL</li>
            </ul>
            <p>
              If all four show, your metadata is right. Lighthouse SEO 100 is
              easy if your metadata is complete and you have a sitemap.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 3 — <strong>The User Dashboard</strong>. Signup, login, the
        logged-in surface. The customer&apos;s home base.
      </p>
    </>
  )
}
