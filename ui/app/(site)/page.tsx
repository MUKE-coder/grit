import Link from 'next/link'
import { ArrowRight, Copy, Layers, Palette, Sparkles, Zap } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { Reveal } from '@/components/reveal'
import { SubcategoryCard } from '@/components/subcategory-card'
import { CopyLine } from '@/components/copy-line'
import { CATALOG } from '@/registry/catalog'
import { countIn, baseUrl } from '@/lib/blocks'
import { BLOCK_COUNT } from '@/lib/block-map'

const PILLARS = [
  {
    icon: Copy,
    title: 'Copy, or one command',
    body: 'Every block is a shadcn registry item. Paste the source, or install it with a single line.',
  },
  {
    icon: Palette,
    title: 'Stock Tailwind only',
    body: 'No design tokens to merge, no config to patch. A block looks right in any Tailwind project.',
  },
  {
    icon: Layers,
    title: 'Light and dark, both correct',
    body: 'Not a dark theme with a light afterthought. Every block is checked in both.',
  },
  {
    icon: Zap,
    title: 'Yours the moment it lands',
    body: 'Blocks are files in your repo, not a dependency. Edit them, rename them, delete them.',
  },
]

export default function HomePage() {
  const base = baseUrl()
  const example = `npx shadcn@latest add ${base}/r/marketing-hero-sections-simple-centered.json`

  return (
    <div className="min-h-screen">
      <SiteHeader />

      {/* ── Hero ─────────────────────────────────────────────────────────── */}
      <section className="hairline relative overflow-hidden border-b">
        <div aria-hidden className="pointer-events-none absolute inset-0 grid-fade" />
        {/* A single wide, very low-opacity bloom. Apple uses light sparingly;
            two competing gradients read as a template. */}
        <div
          aria-hidden
          className="pointer-events-none absolute inset-x-0 -top-40 h-[42rem]"
          style={{
            background:
              'radial-gradient(46rem 26rem at 50% 30%, rgb(99 102 241 / 0.13), transparent 68%)',
          }}
        />

        {/* Centred, with the product directly beneath the copy. A left-aligned
            hero leaves half the viewport empty, which reads as unfinished at
            this width — and Apple always puts the thing itself on screen. */}
        <div className="relative mx-auto max-w-4xl px-6 pt-24 text-center lg:pt-32">
          <Reveal>
            <span className="glass inline-flex items-center gap-2 rounded-full px-3.5 py-1.5 text-[12.5px] font-medium text-gray-700 dark:text-gray-300">
              <Sparkles aria-hidden="true" className="size-3.5 text-indigo-500" />
              {BLOCK_COUNT} blocks and growing
            </span>
          </Reveal>

          <Reveal delay={0.06}>
            <h1 className="display-tight mt-7 text-[3.25rem] leading-[1.03] text-balance text-gray-900 sm:text-[4.5rem] dark:text-white">
              Beautifully crafted{' '}
              <span className="bg-gradient-to-br from-indigo-500 via-violet-500 to-sky-500 bg-clip-text text-transparent">
                React UI blocks
              </span>
            </h1>
          </Reveal>

          <Reveal delay={0.12}>
            <p className="mx-auto mt-7 max-w-xl text-[18px]/[1.6] text-pretty text-gray-500 dark:text-gray-400">
              Production-ready sections built with Tailwind CSS. Drop them into
              your project, then make them yours — the code is ordinary source you
              own from the moment it lands.
            </p>
          </Reveal>

          <Reveal delay={0.18}>
            <div className="mt-9 flex flex-wrap items-center justify-center gap-3">
              <Link
                href="/marketing/hero-sections"
                className="group sheen inline-flex items-center gap-2 rounded-full bg-gray-900 px-6 py-3 text-[14px] font-semibold text-white transition-all duration-200 hover:bg-gray-800 active:scale-[0.98] dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
              >
                Browse blocks
                <ArrowRight
                  aria-hidden="true"
                  className="size-4 transition-transform duration-200 group-hover:translate-x-0.5"
                />
              </Link>
              <Link
                href="https://gritframework.dev/docs/frontend/ui-components"
                className="glass inline-flex items-center rounded-full px-6 py-3 text-[14px] font-semibold text-gray-900 transition-transform duration-200 active:scale-[0.98] dark:text-white"
              >
                Documentation
              </Link>
            </div>
          </Reveal>

          <Reveal delay={0.24}>
            <div className="mx-auto mt-9 max-w-xl">
              <CopyLine value={example} />
            </div>
          </Reveal>
        </div>

        {/* Product window: a real block, rendered live in a frame that is
            clipped by the section so it appears to continue below the fold. */}
        <Reveal delay={0.3}>
          <div className="relative mx-auto mt-16 max-w-6xl px-6">
            <div className="hairline overflow-hidden rounded-t-2xl border border-b-0 bg-white lift-lg dark:bg-gray-900">
              <div className="hairline flex items-center gap-2 border-b px-4 py-2.5">
                <span className="size-2.5 rounded-full bg-[#ff5f57]" />
                <span className="size-2.5 rounded-full bg-[#febc2e]" />
                <span className="size-2.5 rounded-full bg-[#28c840]" />
                <span className="mx-auto font-mono text-[11px] text-gray-400">
                  marketing / hero-sections / simple-centered
                </span>
              </div>
              <iframe
                src="/preview/marketing-hero-sections-simple-centered?theme=light"
                title="A Grit UI block"
                loading="lazy"
                scrolling="no"
                className="h-[430px] w-full border-0 dark:hidden"
              />
              <iframe
                src="/preview/marketing-hero-sections-simple-centered?theme=dark"
                title="A Grit UI block"
                loading="lazy"
                scrolling="no"
                className="hidden h-[430px] w-full border-0 dark:block"
              />
            </div>
          </div>
        </Reveal>
      </section>

      {/* ── Pillars ──────────────────────────────────────────────────────── */}
      <section className="hairline border-b">
        <div className="mx-auto grid max-w-[100rem] grid-cols-1 gap-px bg-gray-500/[0.07] sm:grid-cols-2 lg:grid-cols-4 dark:bg-white/[0.07]">
          {PILLARS.map(({ icon: Icon, title, body }, i) => (
            <Reveal key={title} delay={i * 0.05}>
              <div className="group h-full bg-white px-7 py-9 transition-colors duration-300 hover:bg-gray-50/60 dark:bg-gray-950 dark:hover:bg-white/[0.02]">
                <span className="glass inline-flex size-10 items-center justify-center rounded-xl text-indigo-500 transition-transform duration-300 group-hover:-translate-y-0.5">
                  <Icon aria-hidden="true" className="size-[18px]" />
                </span>
                <h3 className="display mt-5 text-[14.5px] text-gray-900 dark:text-white">
                  {title}
                </h3>
                <p className="mt-2 text-[13.5px]/[1.6] text-gray-500 dark:text-gray-400">
                  {body}
                </p>
              </div>
            </Reveal>
          ))}
        </div>
      </section>

      {/* ── Catalogue ────────────────────────────────────────────────────── */}
      {CATALOG.map((category) => (
        <section key={category.slug} className="hairline border-b">
          <div className="mx-auto max-w-[100rem] px-6 py-16">
            <Reveal>
              <div className="max-w-2xl">
                <h2 className="display-tight text-[2rem] leading-tight text-gray-900 dark:text-white">
                  <Link
                    href={`/${category.slug}`}
                    className="transition-colors hover:text-indigo-600 dark:hover:text-indigo-400"
                  >
                    {category.name}
                  </Link>
                </h2>
                <p className="mt-3 text-[16px]/[1.6] text-gray-500 dark:text-gray-400">
                  {category.description}
                </p>
              </div>
            </Reveal>

            {category.groups.map((group) => (
              <div key={group.name} className="mt-12">
                <Reveal>
                  <p className="label-mono hairline border-b pb-3">{group.name}</p>
                </Reveal>
                <div className="mt-7 grid grid-cols-2 gap-x-6 gap-y-9 md:grid-cols-3 lg:grid-cols-4">
                  {group.subcategories.map((sub, i) => (
                    <Reveal key={sub.slug} delay={Math.min(i * 0.03, 0.2)}>
                      <SubcategoryCard
                        href={`/${category.slug}/${sub.slug}`}
                        name={sub.name}
                        count={countIn(category.slug, sub.slug)}
                      />
                    </Reveal>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>
      ))}

      <footer className="mx-auto max-w-[100rem] px-6 py-12">
        <p className="text-[13px] text-gray-400 dark:text-gray-600">
          Part of the{' '}
          <Link
            href="https://gritframework.dev"
            className="text-gray-600 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
          >
            Grit Framework
          </Link>
          . MIT licensed.
        </p>
      </footer>
    </div>
  )
}
