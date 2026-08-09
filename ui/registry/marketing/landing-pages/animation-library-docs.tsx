/*
 * An open-source library page: hero, capability strip, a feature grid with code,
 * examples, a workflow section, a showcase, a changelog and sponsors.
 *
 * One file, because the registry installs one file per block.
 *
 * Identifiers are <code>, not spans that happen to be monospaced. On this kind
 * of page the difference matters more than usual: half the prose is API names,
 * and a screen reader that knows a run of text is code can announce it as code
 * instead of trying to pronounce useMotionValue.
 *
 * The changelog says MAJOR and MINOR in words as well as colour, and every date
 * is a <time> with a machine-readable datetime. A yellow badge and a red badge
 * are the same badge to a lot of readers.
 *
 * The showcase scrolls horizontally. Every card in it is a link, so tabbing
 * moves through the row and the browser scrolls each card into view — a scroll
 * region whose contents cannot be reached by keyboard is a common and quiet
 * failure, and the fix is usually just making sure the contents are focusable
 * rather than adding a pair of arrow buttons.
 *
 * Yellow is the brand and it is only ever used behind black text. White on
 * yellow-400 measures 1.5:1 and is unreadable at any size; black on the same
 * yellow is 13.9:1, which is why every yellow surface here flips the text
 * rather than the shade.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Docs', 'Examples', 'UI', 'AI kit']

const CAPABILITIES = [
  { n: '01', name: 'Free', body: 'Completely free to use, MIT licensed and open source.' },
  { n: '02', name: 'Production ready', body: 'Used across hundreds of thousands of sites.' },
  { n: '03', name: 'Hybrid engine', body: 'JavaScript and hardware-accelerated browser APIs in one library.' },
  { n: '04', name: 'Built for agents', body: 'Agent-compatible documentation, skills and API surface.' },
  { n: '05', name: 'Tiny footprint', body: 'Up to 90% smaller than the usual alternative.' },
]

const FEATURES = [
  {
    n: '01',
    title: 'Independent transforms',
    body: ['Animate ', 'x', ', ', 'y', ', ', 'rotate', ' and ', 'scale', ' on the same element, without wrappers.'],
    code: '{ rotate: 15, x: "50%" }',
    demo: 'square',
    yellow: false,
  },
  {
    n: '02',
    title: 'Scroll animation',
    body: ['Hardware-accelerated scroll-linked motion via ', 'scrollTimeline', '.'],
    code: 'scroll()',
    demo: 'orbit',
    yellow: true,
  },
  {
    n: '03',
    title: 'Native gestures',
    body: ['', 'hover', ', ', 'press', ' and ', 'drag', ' that feel native, not bolted on.'],
    code: 'drag = true',
    demo: 'gesture',
    yellow: false,
  },
  {
    n: '04',
    title: 'Layout animation',
    body: ['Animate between any two layouts with a single ', 'layout', ' prop.'],
    code: 'layout = true',
    demo: 'grid',
    yellow: false,
  },
  {
    n: '05',
    title: 'Spring physics',
    body: ['Real ', 'spring', ' maths for animations that react naturally to input.'],
    code: 'type: "spring"',
    demo: 'ball',
    yellow: false,
  },
  {
    n: '06',
    title: 'Exit animation',
    body: ['', 'AnimatePresence', ' keeps elements alive so they can animate as they leave the DOM.'],
    code: 'exit={{ ... }}',
    demo: 'stack',
    yellow: false,
  },
  {
    n: '07',
    title: 'Timeline sequences',
    body: ['', 'variants', ', ', 'stagger', ' and timelines orchestrate complex motion.'],
    code: 'stagger(0.04)',
    demo: 'lines',
    yellow: true,
  },
  {
    n: '08',
    title: 'Motion values',
    body: ['', 'useMotionValue', ' drives animations and derived state in real time.'],
    code: 'useTransform(v, x => x * 10)',
    demo: 'tilt',
    yellow: false,
  },
]

const EXAMPLES = [
  { name: 'Parallax', art: 'from-stone-700 via-stone-800 to-stone-950' },
  { name: 'Confetti', art: 'from-gray-800 via-gray-900 to-black' },
  { name: 'Typewriter', art: 'from-slate-800 via-slate-900 to-black' },
  { name: 'App folder', art: 'from-gray-100 via-gray-200 to-gray-300' },
  { name: 'Pointer animation', art: 'from-sky-100 via-white to-gray-200' },
  { name: 'Modal', art: 'from-rose-200 via-white to-rose-100' },
]

const SHOWCASE = [
  { name: '3D cylinder gallery', by: '@dimi', art: 'from-stone-200 via-stone-300 to-stone-400' },
  { name: 'Dia browser', by: '@buenasuerte', art: 'from-indigo-300 via-fuchsia-300 to-amber-200' },
  { name: 'Collins carousel', by: '@dimi', art: 'from-gray-100 via-gray-200 to-gray-300' },
  { name: 'Business switcher', by: '@egdiala', art: 'from-sky-200 via-white to-sky-100' },
  { name: 'Scroll grid', by: '@hev', art: 'from-gray-800 via-gray-900 to-black' },
  { name: 'Radial menu', by: '@sonya', art: 'from-violet-300 via-purple-300 to-indigo-400' },
]

const CHANGELOG = [
  {
    version: '13.0.0',
    kind: 'Major',
    kindClass: 'bg-yellow-400 text-gray-950',
    date: '2026-08-05',
    dateLabel: '5 August 2026',
    groups: [
      {
        label: 'Changed',
        items: [
          ['Removed the optional ', '@emotion/is-prop-valid', ' dependency in favour of an explicit ', '<MotionConfig isValidProp>', '.'],
        ],
      },
      {
        label: 'Fixed',
        items: [
          ['Hardware-accelerated SVG elements now apply their final style when the animation completes.'],
          ['', 'AnimatePresence', ' marks nodes safe to remove when rendering ', 'propagate', ' with no motion children.'],
        ],
      },
    ],
  },
  {
    version: '12.43.0',
    kind: 'Minor',
    kindClass: 'bg-white/15 text-gray-100',
    date: '2026-07-27',
    dateLabel: '27 July 2026',
    groups: [
      {
        label: 'Added',
        items: [
          ['Hardware acceleration for ', 'backgroundColor', ' in supported browsers.'],
          ['Hardware acceleration for SVG elements.'],
        ],
      },
    ],
  },
]

const SPONSORS = ['Framer', 'Figma', 'Vercel', 'Linear', 'Raycast', 'Supabase']

/* ── Demos ──────────────────────────────────────────────────────────────── */

function Demo({ kind, yellow }: { kind: string; yellow: boolean }) {
  const mark = yellow ? 'bg-gray-950' : 'bg-yellow-400'
  const line = yellow ? 'bg-gray-950/25' : 'bg-white/15'
  const ring = yellow ? 'border-gray-950/30' : 'border-white/20'

  return (
    <span aria-hidden="true" className="relative block h-28 select-none">
      {kind === 'square' && (
        <>
          <span className={`absolute inset-x-0 top-1/2 h-px ${line}`} />
          <span className={`absolute inset-y-0 left-1/3 w-px ${line}`} />
          <span className={`absolute top-6 left-4 size-7 rounded-sm ${mark}`} />
        </>
      )}
      {kind === 'orbit' && (
        <>
          <span className={`absolute top-1/2 left-1/2 size-24 -translate-x-1/2 -translate-y-1/2 rounded-full border ${ring}`} />
          <span className={`absolute top-1/2 left-1/2 size-16 -translate-x-1/2 -translate-y-1/2 rounded-full border ${ring}`} />
          <span className={`absolute top-[68%] left-[70%] size-2.5 rounded-full ${mark}`} />
        </>
      )}
      {kind === 'gesture' && (
        <>
          <span className={`absolute inset-x-8 inset-y-3 rounded border ${ring}`} />
          <span className={`absolute top-1/2 left-1/2 size-8 -translate-x-1/2 -translate-y-1/2 rounded-sm ${mark}`} />
        </>
      )}
      {kind === 'grid' && (
        <span className="absolute inset-x-10 inset-y-2 grid grid-cols-2 gap-2">
          {[0, 1, 2, 3].map((i) => (
            <span key={i} className={`rounded-sm ${yellow ? 'bg-gray-950' : 'bg-gray-100'}`} />
          ))}
        </span>
      )}
      {kind === 'ball' && (
        <>
          <span className={`absolute inset-x-4 bottom-8 h-px ${line}`} />
          <span className={`absolute bottom-6 left-1/2 size-6 -translate-x-1/2 rounded-full ${yellow ? 'bg-gray-950' : 'bg-white'}`} />
        </>
      )}
      {kind === 'stack' && (
        <span className="absolute inset-x-6 inset-y-3 space-y-2">
          {[0, 1, 2].map((i) => (
            <span key={i} className={`flex h-6 items-center gap-2 rounded ${yellow ? 'bg-gray-950/10' : 'bg-white/10'} px-2`}>
              <span className={`size-1.5 rounded-full ${mark}`} />
            </span>
          ))}
        </span>
      )}
      {kind === 'lines' && (
        <span className="absolute inset-x-6 inset-y-3 space-y-2">
          {[100, 82, 92, 70, 88].map((w, i) => (
            <span key={i} style={{ width: `${w}%` }} className={`block h-3 rounded-sm ${yellow ? 'bg-gray-950/20' : 'bg-white/15'}`} />
          ))}
        </span>
      )}
      {kind === 'tilt' && (
        <>
          <span className={`absolute top-4 left-1/2 size-8 -translate-x-1/2 rounded-sm ${mark}`} />
          <span className={`absolute bottom-4 left-1/2 size-8 -translate-x-1/2 rotate-12 rounded-sm border ${ring}`} />
        </>
      )}
    </span>
  )
}

/* Alternating plain text and identifiers. Odd indexes are code — the shape of
   the array is the markup, which keeps the copy readable in the data above. */
function Prose({ parts }: { parts: string[] }) {
  return (
    <>
      {parts.map((part, i) =>
        i % 2 === 1 ? (
          <code key={i} className="font-mono text-[0.9em]">
            {part}
          </code>
        ) : (
          <span key={i}>{part}</span>
        ),
      )}
    </>
  )
}

function SectionHead({
  n,
  label,
  heading,
  body,
  action,
  id,
}: {
  n: string
  label: string
  heading: string
  body?: string
  action?: ReactNode
  id: string
}) {
  return (
    <div className="grid gap-6 border-t border-white/15 pt-4 lg:grid-cols-[16rem_1fr_auto]">
      <p className="flex items-start justify-between gap-4 font-mono text-[11px] tracking-widest text-gray-400 uppercase lg:pr-8">
        <span aria-hidden="true">{n}</span>
        <span>{label}</span>
      </p>
      <div>
        <h2 id={id} className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
          {heading}
        </h2>
        {body && <p className="mt-4 max-w-sm text-sm text-pretty text-gray-300">{body}</p>}
      </div>
      {action && <div className="lg:text-right">{action}</div>}
    </div>
  )
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function AnimationLibraryDocs() {
  return (
    <div className="min-h-screen bg-gray-950 font-sans text-gray-50 antialiased">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:bg-yellow-400 focus:px-4 focus:py-2 focus:text-sm focus:text-gray-950"
      >
        Skip to content
      </a>

      {/* Hero */}
      <div className="relative isolate overflow-hidden">
        {/* The banner texture is drawn with gradients rather than shipped as an
            image. It is decoration, it has to survive any width, and a 400KB
            mesh render is a lot of bytes to say "this library is about
            motion". */}
        <span
          aria-hidden="true"
          className="absolute inset-0 -z-10 bg-yellow-400"
          style={{
            backgroundImage:
              'radial-gradient(60% 80% at 78% 30%, rgba(236,72,153,0.85), transparent 60%),' +
              'radial-gradient(50% 70% at 62% 70%, rgba(56,189,248,0.75), transparent 60%),' +
              'radial-gradient(45% 60% at 88% 62%, rgba(139,92,246,0.7), transparent 60%),' +
              'repeating-linear-gradient(115deg, rgba(255,255,255,0.22) 0 2px, transparent 2px 7px)',
          }}
        />

        <header className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
          <a href="#" className="flex items-center gap-2 text-sm font-semibold text-gray-950">
            <span aria-hidden="true">▰▰</span>
            Motion
          </a>
          <nav aria-label="Primary" className="hidden md:block">
            <ul role="list" className="flex items-center gap-7 font-mono text-[11px] tracking-widest text-gray-950 uppercase">
              {NAV.map((item) => (
                <li key={item}>
                  <a href="#" className="hover:underline">
                    {item}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
          <a href="#" className="bg-gray-950 px-4 py-2 font-mono text-[11px] tracking-widest text-yellow-400 uppercase hover:bg-gray-800">
            Motion+
          </a>
        </header>

        <div className="mx-auto max-w-6xl px-6 pt-10 pb-24">
          <div className="max-w-lg bg-gray-950 p-8">
            <p className="flex items-center justify-between gap-4 font-mono text-[11px] tracking-widest text-gray-400 uppercase">
              <span>Open source / MIT licence</span>
              <span>v13.0.0</span>
            </p>
            <h1 className="mt-6 text-4xl font-semibold tracking-tight text-balance">
              <span className="text-yellow-400">Motion.</span> Production-grade animation for the web.
            </h1>
            <a
              href="#library"
              className="mt-8 inline-block bg-yellow-400 px-6 py-3 font-mono text-[11px] tracking-widest text-gray-950 uppercase hover:bg-yellow-300"
            >
              Get started
            </a>
            <p className="mt-8 font-mono text-[11px] tracking-widest text-gray-400 uppercase">Available for</p>
            <ul role="list" className="mt-3 flex flex-wrap gap-x-6 gap-y-2 font-mono text-[11px] tracking-widest uppercase">
              {['React', 'JavaScript', 'Vue'].map((target) => (
                <li key={target}>{target}</li>
              ))}
            </ul>
          </div>
        </div>
      </div>

      <main id="main">
        {/* Capability strip */}
        <section aria-labelledby="capabilities-heading" className="border-y border-white/15">
          <h2 id="capabilities-heading" className="sr-only">
            Why use Motion
          </h2>
          <ul role="list" className="mx-auto grid max-w-6xl divide-y divide-white/15 px-6 sm:grid-cols-2 sm:divide-y-0 lg:grid-cols-5 lg:divide-x">
            {CAPABILITIES.map((item) => (
              <li key={item.n} className="p-5 lg:px-5">
                <p className="font-mono text-[11px] tracking-widest text-yellow-400 uppercase">
                  <span aria-hidden="true">{item.n} / </span>
                  {item.name}
                </p>
                <p className="mt-3 text-sm text-pretty text-gray-300">{item.body}</p>
              </li>
            ))}
          </ul>
        </section>

        {/* Feature grid */}
        <section id="library" aria-labelledby="library-heading" className="px-6 py-20">
          <div className="mx-auto max-w-6xl">
            <SectionHead
              id="library-heading"
              n="01"
              label="Animation library"
              heading="Animations that move."
              body="Create high-performance web animations with an easy API, from simple transforms to advanced interactive gestures."
            />

            <ul role="list" className="mt-12 grid gap-px border border-white/15 bg-white/15 sm:grid-cols-2 lg:grid-cols-4">
              {FEATURES.map((feature) => (
                <li
                  key={feature.n}
                  className={`flex flex-col ${feature.yellow ? 'bg-yellow-400 text-gray-950' : 'bg-gray-950'}`}
                >
                  <Demo kind={feature.demo} yellow={feature.yellow} />
                  <div className="flex flex-1 flex-col p-5">
                    <p className={`font-mono text-[11px] tracking-widest uppercase ${feature.yellow ? 'text-gray-950/70' : 'text-gray-400'}`}>
                      {feature.n}
                    </p>
                    <h3 className="mt-2 flex items-center justify-between gap-3 font-semibold">
                      {feature.title}
                      <span aria-hidden="true" className={feature.yellow ? 'text-gray-950/60' : 'text-gray-500'}>
                        →
                      </span>
                    </h3>
                    <p className={`mt-2 text-sm text-pretty ${feature.yellow ? 'text-gray-950/80' : 'text-gray-300'}`}>
                      <Prose parts={feature.body} />
                    </p>
                    <p
                      className={`mt-auto pt-6 font-mono text-[11px] ${
                        feature.yellow ? 'text-gray-950/70' : 'text-yellow-400'
                      }`}
                    >
                      <code>{feature.code}</code>
                    </p>
                  </div>
                </li>
              ))}
            </ul>

            <div className="flex flex-wrap items-center justify-between gap-4 border-x border-b border-white/15 px-5 py-4 font-mono text-[11px] tracking-widest uppercase">
              <p className="text-gray-400">Available for</p>
              <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
                {['React', 'JavaScript', 'Vue'].map((target) => (
                  <li key={target}>{target}</li>
                ))}
              </ul>
              <a href="#" className="text-yellow-400 hover:underline">
                All documentation →
              </a>
            </div>
          </div>
        </section>

        {/* Examples */}
        <section aria-labelledby="examples-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <SectionHead
              id="examples-heading"
              n="02"
              label="Examples"
              heading="Copy, paste, ship."
              body="Build previously complex effects, like magnetic cursors or infinite tickers, with 410+ copy-and-paste examples."
              action={
                <a href="#" className="font-mono text-[11px] tracking-widest text-yellow-400 uppercase hover:underline">
                  Browse all examples →
                </a>
              }
            />

            <div className="mt-12 border border-white/15">
              <div className="flex flex-wrap items-center justify-between gap-3 border-b border-white/15 px-4 py-2.5 font-mono text-[11px] tracking-widest uppercase">
                <p className="text-gray-400">
                  <span aria-hidden="true">&gt; </span>Live example / skeleton shimmer
                </p>
                <ul role="list" className="flex gap-2">
                  <li>
                    <a href="#" className="border border-white/20 px-3 py-1.5 hover:bg-white/10">
                      View source
                    </a>
                  </li>
                  <li>
                    <a href="#" className="border border-white/20 px-3 py-1.5 hover:bg-white/10">
                      Open in editor
                    </a>
                  </li>
                </ul>
              </div>

              <div className="grid place-items-center bg-gray-900/60 p-10">
                {/* A profile card, drawn. This is the thing the example renders,
                    so it is decoration for the surrounding copy. */}
                <span aria-hidden="true" className="block w-64 overflow-hidden rounded-xl bg-gray-950 select-none">
                  <span className="block h-20 bg-gradient-to-r from-blue-600 to-violet-600" />
                  <span className="block p-4">
                    <span className="-mt-10 block size-12 rounded-full border-4 border-gray-950 bg-yellow-400" />
                    <span className="mt-3 block text-sm font-semibold">Motion</span>
                    <span className="block text-xs text-gray-400">@motiondotdev</span>
                    <span className="mt-3 block text-xs text-gray-300">
                      Free and open source. Animation for React, JavaScript and Vue.
                    </span>
                    <span className="mt-4 grid grid-cols-3 gap-2 text-center">
                      {[['127', 'Posts'], ['11K', 'Followers'], ['5', 'Following']].map(([v, k]) => (
                        <span key={k} className="block rounded bg-white/5 py-2">
                          <span className="block text-xs font-semibold">{v}</span>
                          <span className="block text-[10px] text-gray-400">{k}</span>
                        </span>
                      ))}
                    </span>
                    <span className="mt-3 block rounded bg-blue-600 py-2 text-center text-xs font-semibold">Follow</span>
                  </span>
                </span>
              </div>

              <ul role="list" className="grid gap-px border-t border-white/15 bg-white/15 sm:grid-cols-3 lg:grid-cols-6">
                {EXAMPLES.map((example) => (
                  <li key={example.name} className="bg-gray-950">
                    <a href="#" className="block p-3 hover:bg-white/5">
                      <span aria-hidden="true" className={`block h-16 bg-gradient-to-br ${example.art}`} />
                      <span className="mt-2 block font-mono text-[11px] tracking-widest uppercase">{example.name}</span>
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </section>

        {/* Workflow */}
        <section aria-labelledby="workflow-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <SectionHead
              id="workflow-heading"
              n="03"
              label="Workflow"
              heading="Give every project a head start."
              body="Equip your agent with specialist context, then build from production-ready animated sections that inherit your design system."
            />

            <div className="mt-12 grid border border-white/15 lg:grid-cols-2">
              <div className="border-b border-white/15 p-8 lg:border-r lg:border-b-0">
                <p className="font-mono text-[11px] tracking-widest text-violet-300 uppercase">
                  Help your agent understand Motion
                </p>
                <h3 className="mt-8 text-3xl font-semibold tracking-tight text-balance">
                  Give your agent specialist judgement.
                </h3>
                <p className="mt-4 max-w-sm text-sm text-pretty text-gray-300">
                  Send the latest docs, 410+ example sources, performance audits and production-ready
                  springs straight to your agent.
                </p>
                {/* violet-500 with black text rather than white on violet-600.
                    White on violet-600 is 4.4:1 — close enough to look fine and
                    still short of the 4.5:1 a 12px label needs. */}
                <a href="#" className="mt-6 inline-block bg-violet-400 px-5 py-2.5 font-mono text-[11px] tracking-widest text-gray-950 uppercase hover:bg-violet-300">
                  Explore the AI kit
                </a>
              </div>

              <div className="p-8">
                <p className="font-mono text-[11px] tracking-widest text-sky-300 uppercase">
                  Use pre-built animated sections
                </p>
                <h3 className="mt-8 text-3xl font-semibold tracking-tight text-balance">
                  Begin with production-ready motion.
                </h3>
                <p className="mt-4 max-w-sm text-sm text-pretty text-gray-300">
                  Browse performance-rated animated sections, install the source and inherit the design
                  tokens your project already uses.
                </p>
                <a href="#" className="mt-6 inline-block bg-sky-300 px-5 py-2.5 font-mono text-[11px] tracking-widest text-gray-950 uppercase hover:bg-sky-200">
                  Browse the UI library
                </a>
              </div>

              <div className="border-t border-white/15 lg:col-span-2">
                <div className="grid lg:grid-cols-2">
                  <div className="border-b border-white/15 p-6 lg:border-r lg:border-b-0">
                    <span aria-hidden="true" className="block overflow-hidden rounded-lg border border-white/15 bg-gray-900 select-none">
                      <span className="flex gap-1.5 border-b border-white/10 px-3 py-2">
                        <span className="size-2.5 rounded-full bg-rose-500" />
                        <span className="size-2.5 rounded-full bg-yellow-400" />
                        <span className="size-2.5 rounded-full bg-emerald-500" />
                      </span>
                      <span className="block space-y-1.5 p-4 font-mono text-[11px]">
                        <span className="block text-gray-300">/motion create a photo carousel</span>
                        <span className="block text-yellow-400">Searching documentation for carousel</span>
                        <span className="block text-yellow-400">Searching examples for "photo carousel"</span>
                        <span className="block text-gray-300">Found 1 doc and 3 examples.</span>
                        <span className="block text-gray-300">Building from your design system.</span>
                      </span>
                    </span>
                  </div>
                  <div className="relative grid min-h-56 place-items-center bg-gray-100 p-6">
                    <span aria-hidden="true" className="block h-24 w-40 rounded-lg bg-gray-300" />
                    <p className="absolute right-4 bottom-4 bg-gray-950 px-3 py-1.5 font-mono text-[10px] tracking-widest uppercase">
                      Live section preview
                    </p>
                  </div>
                </div>
              </div>

              <div className="flex flex-wrap items-center justify-between gap-4 border-t border-white/15 bg-emerald-950 p-6 lg:col-span-2">
                <div>
                  <p className="font-mono text-[11px] tracking-widest text-emerald-300 uppercase">Motion+</p>
                  <h3 className="mt-2 text-xl font-semibold">Bring both into one workflow.</h3>
                  <p className="mt-1 max-w-md text-sm text-gray-300">
                    The AI kit and the UI library are both included, with premium APIs, examples and
                    lifetime updates.
                  </p>
                </div>
                <a href="#" className="bg-emerald-300 px-5 py-2.5 font-mono text-[11px] tracking-widest text-gray-950 uppercase hover:bg-emerald-200">
                  Get Motion+
                </a>
              </div>
            </div>
          </div>
        </section>

        {/* Showcase */}
        <section aria-labelledby="showcase-heading" className="pb-20">
          <div className="mx-auto max-w-6xl px-6">
            <SectionHead
              id="showcase-heading"
              n="04"
              label="Showcase"
              heading="Made with Motion."
              body="Everything is possible with Motion. Here are some of the best things the community has built."
              action={
                <a href="#" className="font-mono text-[11px] tracking-widest text-yellow-400 uppercase hover:underline">
                  Submit your work →
                </a>
              }
            />
          </div>

          {/* Scrolls sideways, and every card is a link, so tabbing walks the
              row and the browser brings each card into view. A scroller whose
              contents are not focusable is unreachable without a mouse. */}
          <ul
            role="list"
            className="mt-12 flex gap-px overflow-x-auto bg-white/15 pb-2 [scrollbar-width:thin]"
          >
            {SHOWCASE.map((item) => (
              <li key={item.name} className="w-64 shrink-0 bg-gray-950">
                <a href="#" className="block hover:opacity-90">
                  <span aria-hidden="true" className={`block h-36 bg-gradient-to-br ${item.art}`} />
                  <span className="flex items-center justify-between gap-3 p-3 font-mono text-[11px] tracking-widest uppercase">
                    <span>{item.name}</span>
                    <span className="text-gray-400">{item.by}</span>
                  </span>
                </a>
              </li>
            ))}
          </ul>
        </section>

        {/* Updates */}
        <section aria-labelledby="updates-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <SectionHead
              id="updates-heading"
              n="05"
              label="Updates"
              heading="Latest from Motion."
              body="Release notes and magazine stories from the team."
            />

            <div className="mt-12 grid gap-10 lg:grid-cols-2">
              <div>
                <h3 className="flex items-center justify-between border-b border-white/15 pb-3 font-mono text-[11px] tracking-widest text-gray-400 uppercase">
                  Changelog
                  <a href="#" className="text-yellow-400 hover:underline">
                    RSS
                  </a>
                </h3>

                <ol className="mt-6 space-y-10">
                  {CHANGELOG.map((release) => (
                    <li key={release.version}>
                      <p className="flex flex-wrap items-center gap-2">
                        {/* The release kind is a word. A yellow badge and a grey
                            badge are the same badge to anyone not distinguishing
                            them by colour. */}
                        <span className={`px-2 py-0.5 font-mono text-[10px] tracking-widest uppercase ${release.kindClass}`}>
                          {release.kind}
                        </span>
                      </p>
                      <h4 className="mt-3 border border-dashed border-white/25 px-3 py-2 font-mono text-2xl">
                        {release.version}
                      </h4>
                      <p className="mt-2 font-mono text-[11px] tracking-widest text-gray-400 uppercase">
                        <time dateTime={release.date}>{release.dateLabel}</time>
                      </p>

                      {release.groups.map((group) => (
                        <div key={group.label} className="mt-5">
                          <h5 className="font-mono text-[11px] tracking-widest text-yellow-400 uppercase">{group.label}</h5>
                          <ul role="list" className="mt-2 space-y-2 text-sm text-gray-300">
                            {group.items.map((item, i) => (
                              <li key={i} className="flex gap-2.5">
                                <span aria-hidden="true" className="mt-1.5 size-1 shrink-0 rounded-full bg-gray-500" />
                                <span className="text-pretty">
                                  <Prose parts={item} />
                                </span>
                              </li>
                            ))}
                          </ul>
                        </div>
                      ))}
                    </li>
                  ))}
                </ol>

                <a href="#" className="mt-8 inline-block border border-white/20 px-5 py-2.5 font-mono text-[11px] tracking-widest uppercase hover:bg-white/10">
                  Full changelog
                </a>
              </div>

              <div>
                <h3 className="flex items-center justify-between border-b border-white/15 pb-3 font-mono text-[11px] tracking-widest text-gray-400 uppercase">
                  Magazine
                  <a href="#" className="text-yellow-400 hover:underline">
                    RSS
                  </a>
                </h3>

                <article className="mt-6">
                  <span aria-hidden="true" className="block h-56 bg-gradient-to-br from-sky-300 via-blue-500 to-indigo-700" />
                  <p className="mt-4 font-mono text-[11px] tracking-widest text-yellow-400 uppercase">
                    Announcement — <time dateTime="2026-07-23">23 July 2026</time>
                  </p>
                  <h4 className="mt-2 text-xl font-semibold">
                    <a href="#" className="hover:underline">
                      Introducing the UI library
                    </a>
                  </h4>
                  <p className="mt-2 text-sm text-pretty text-gray-300">
                    Production-ready animated sections and components for React, performance-graded,
                    dropping into your design system through the registry or an agent.
                  </p>
                  <p className="mt-3 font-mono text-[11px] tracking-widest text-gray-400 uppercase">Matt Perry</p>
                </article>

                <a href="#" className="mt-8 inline-block border border-white/20 px-5 py-2.5 font-mono text-[11px] tracking-widest uppercase hover:bg-white/10">
                  All articles
                </a>
              </div>
            </div>
          </div>
        </section>

        {/* Sponsors */}
        <section aria-labelledby="sponsors-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <SectionHead
              id="sponsors-heading"
              n="06"
              label="Sponsors"
              heading="Trusted by the world's best teams."
              action={
                <a href="#" className="font-mono text-[11px] tracking-widest text-yellow-400 uppercase hover:underline">
                  Become a sponsor →
                </a>
              }
            />
            <ul role="list" className="mt-10 flex flex-wrap items-center gap-x-10 gap-y-4">
              {SPONSORS.map((name) => (
                <li key={name} className="text-lg font-semibold text-gray-400">
                  {name}
                </li>
              ))}
            </ul>
          </div>
        </section>
      </main>

      <footer className="border-t border-white/15 px-6 py-10">
        <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-6 font-mono text-[11px] tracking-widest uppercase">
          <p className="flex items-center gap-2 font-semibold">
            <span aria-hidden="true" className="text-yellow-400">▰▰</span>
            Motion
          </p>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-8 gap-y-2">
              {['Docs', 'Examples', 'Changelog', 'Licence', 'GitHub'].map((link) => (
                <li key={link}>
                  <a href="#" className="text-gray-400 hover:text-white">
                    {link}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
          <p className="text-gray-400">MIT licensed</p>
        </div>
      </footer>
    </div>
  )
}
