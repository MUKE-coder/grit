import { MousePointer2, Smartphone, SquarePen, Palette } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-7' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

const PROJECT_TYPES = [
  { icon: SquarePen, label: 'Web app', active: true },
  { icon: MousePointer2, label: 'UX/UI design' },
  { icon: Smartphone, label: 'Mobile app' },
  { icon: Palette, label: 'Branding & logo' },
]

export default function MinimalWithProductCard({
  title = 'AI-powered project briefs for developers',
  subtitle = 'Turn an idea into a complete technical brief in seconds. Let the assistant write the spec while you focus on building the thing.',
  ctaLabel = 'Get started for free',
  socialProof = 'Join 80,000+ developers',
  avatarCount = '1,234+',
}: {
  title?: string
  subtitle?: string
  ctaLabel?: string
  socialProof?: string
  avatarCount?: string
}) {
  return (
    <div className="relative min-h-screen overflow-hidden bg-gray-100 dark:bg-gray-950">
      {/* Corner blooms — the only colour in an otherwise neutral layout, so
          they carry the whole mood. */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 opacity-70 dark:opacity-30"
        style={{
          background:
            'radial-gradient(26rem 22rem at 4% 62%, rgba(251,191,36,0.30), transparent 60%), radial-gradient(26rem 22rem at 96% 60%, rgba(244,114,182,0.28), transparent 60%)',
        }}
      />

      {/* Nav */}
      <nav
        aria-label="Global"
        className="relative mx-auto flex max-w-6xl items-center justify-between px-6 py-6"
      >
        <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
          <GritMark className="size-6" />
          <span className="text-base font-semibold tracking-tight">grit</span>
        </a>
        <a
          href="#"
          className="rounded-full bg-white px-5 py-2 text-sm font-medium text-gray-900 shadow-sm transition-colors hover:bg-gray-50 dark:bg-white/10 dark:text-white dark:hover:bg-white/15"
        >
          Sign in
        </a>
      </nav>

      {/* Ruled canvas — hairlines that frame the content column */}
      <div className="relative mx-auto max-w-5xl px-6">
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-y-0 left-[12%] w-px bg-gray-900/[0.06] dark:bg-white/[0.06]"
        />
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-y-0 right-[12%] w-px bg-gray-900/[0.06] dark:bg-white/[0.06]"
        />

        <div className="border-y border-gray-900/[0.06] py-20 text-center dark:border-white/[0.06]">
          <h1 className="mx-auto max-w-3xl text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h1>
          <p className="mx-auto mt-6 max-w-xl text-base/7 text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>
          <a
            href="#"
            className="mt-10 inline-block rounded-full bg-gray-900 px-7 py-3.5 text-sm font-semibold text-white transition-colors hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {ctaLabel}
          </a>
        </div>

        {/* Product card */}
        <div className="relative pt-14 pb-16">
          {/* The product mock is a drawing of an app, not an app: fake rows,
          invented dates, no behaviour. aria-hidden so a screen reader does
          not read the chrome out as page content, and so its deliberately
          low-contrast styling stops being a text-contrast failure. Nothing
          inside is focusable, which is what makes hiding it safe. */}
          <div aria-hidden="true" className="mx-auto max-w-3xl rounded-2xl bg-white/70 p-3 shadow-[0_24px_70px_-20px_rgba(15,23,42,0.25)] backdrop-blur dark:bg-white/5">
            <div className="rounded-xl bg-white p-5 dark:bg-gray-900">
              <div className="mb-5 flex items-center justify-between">
                <p className="flex items-center gap-2 text-sm">
                  <span className="size-2 rounded-full bg-green-500" />
                  <span className="text-gray-400">My briefs</span>
                  <span className="text-gray-300 dark:text-gray-700">/</span>
                  <span className="font-semibold text-gray-900 dark:text-white">
                    Grit Studio 2026
                  </span>
                </p>
                <span className="size-8 rounded-full bg-gradient-to-br from-indigo-400 to-fuchsia-400" />
              </div>

              <div className="grid gap-4 sm:grid-cols-[1.3fr_1fr]">
                <div className="space-y-4">
                  <div className="rounded-xl bg-gray-50 p-4 dark:bg-white/5">
                    <p className="mb-2.5 text-sm font-medium text-gray-900 dark:text-white">
                      Introduction
                    </p>
                    {[100, 92, 78].map((w) => (
                      <div
                        key={w}
                        className="mb-1.5 h-1.5 rounded bg-gray-200 dark:bg-white/10"
                        style={{ width: `${w}%` }}
                      />
                    ))}
                  </div>
                  <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-gray-800">
                    <p className="mb-2.5 text-sm font-medium text-gray-900 dark:text-white">
                      Goal
                    </p>
                    {/* Selected text — the detail that makes this read as a
                        live editor rather than a static wireframe. */}
                    <div className="mb-1.5 h-2 w-[68%] rounded bg-blue-600" />
                    <div className="mb-1.5 h-1.5 w-full rounded bg-gray-200 dark:bg-white/10" />
                    <div className="h-1.5 w-4/5 rounded bg-gray-200 dark:bg-white/10" />
                  </div>
                </div>

                <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-gray-800">
                  <p className="mb-3 text-sm font-medium text-gray-900 dark:text-white">
                    What type of project?
                  </p>
                  <div className="grid grid-cols-2 gap-2">
                    {PROJECT_TYPES.map(({ icon: Icon, label, active }) => (
                      <div
                        key={label}
                        className={`relative rounded-lg border p-2.5 ${
                          active
                            ? 'border-gray-900 bg-white dark:border-white dark:bg-gray-900'
                            : 'border-gray-200 dark:border-white/10'
                        }`}
                      >
                        <Icon
                          aria-hidden="true"
                          className="mb-2 size-4 text-gray-700 dark:text-gray-300"
                        />
                        <p className="text-[11px] text-gray-700 dark:text-gray-300">
                          {label}
                        </p>
                        {active && (
                          <MousePointer2
                            aria-hidden="true"
                            className="absolute right-1.5 bottom-1.5 size-3.5 fill-gray-900 text-gray-900 dark:fill-white dark:text-white"
                          />
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Social proof */}
          <div className="mt-12 text-center">
            <p className="text-sm text-gray-600 dark:text-gray-400">{socialProof}</p>
            <div className="mt-3 flex items-center justify-center">
              <div className="flex -space-x-2">
                {[
                  'from-amber-300 to-orange-400',
                  'from-sky-300 to-blue-400',
                  'from-emerald-300 to-teal-400',
                  'from-fuchsia-300 to-purple-400',
                  'from-rose-300 to-pink-400',
                ].map((tone) => (
                  <span
                    key={tone}
                    className={`size-8 rounded-full bg-gradient-to-br ring-2 ring-gray-100 dark:ring-gray-950 ${tone}`}
                  />
                ))}
              </div>
              <span className="ml-3 rounded-full bg-white px-3 py-1 text-xs font-medium text-blue-600 shadow-sm dark:bg-white/10 dark:text-blue-400">
                {avatarCount}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
