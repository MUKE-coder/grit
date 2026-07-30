'use client'

import { useState } from 'react'
import { ArrowRight, Menu, X } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-8' }: { className?: string }) {
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


const NAV = [
  { name: 'Product', href: '#' },
  { name: 'Features', href: '#' },
  { name: 'Marketplace', href: '#' },
  { name: 'Company', href: '#' },
]

export default function SimpleCentered({
  announcement = 'Grit v3 is out. Read the release notes.',
  title = 'Ship your backend and frontend together',
  subtitle = 'One command scaffolds a Go API, an admin panel and a typed React frontend — wired, tested and ready to deploy. Spend your first week on the product instead of the plumbing.',
  primaryLabel = 'Get started',
  secondaryLabel = 'Read the docs',
}: {
  announcement?: string
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  const [menuOpen, setMenuOpen] = useState(false)

  return (
    <div className="bg-white dark:bg-gray-900">
      <header className="absolute inset-x-0 top-0 z-50">
        <nav
          aria-label="Global"
          className="flex items-center justify-between p-6 lg:px-8"
        >
          {/* lg:flex-1 on the outer div is what pushes the nav links to the
              centre; without it they butt straight up against the logo. */}
          <div className="flex lg:flex-1">
            <a href="#" className="-m-1.5 flex items-center gap-2 p-1.5">
              <span className="sr-only">Grit</span>
              <GritMark className="size-8 text-indigo-600 dark:text-indigo-400" />
            </a>
          </div>

          <div className="flex lg:hidden">
            <button
              type="button"
              onClick={() => setMenuOpen(true)}
              className="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700 dark:text-gray-300"
            >
              <span className="sr-only">Open main menu</span>
              <Menu aria-hidden="true" className="size-6" />
            </button>
          </div>

          <div className="hidden lg:flex lg:gap-x-12">
            {NAV.map((item) => (
              <a
                key={item.name}
                href={item.href}
                className="text-sm/6 font-semibold text-gray-900 dark:text-white"
              >
                {item.name}
              </a>
            ))}
          </div>

          <div className="hidden lg:flex lg:flex-1 lg:justify-end">
            <a
              href="#"
              className="text-sm/6 font-semibold text-gray-900 dark:text-white"
            >
              Log in <span aria-hidden="true">&rarr;</span>
            </a>
          </div>
        </nav>

        {/* Mobile menu */}
        {menuOpen && (
          <div className="lg:hidden" role="dialog" aria-modal="true">
            <div
              className="fixed inset-0 z-50 bg-gray-900/20"
              onClick={() => setMenuOpen(false)}
            />
            <div className="fixed inset-y-0 right-0 z-50 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10 dark:bg-gray-900 dark:sm:ring-white/10">
              <div className="flex items-center justify-between">
                <a href="#" className="-m-1.5 p-1.5">
                  <span className="sr-only">Grit</span>
                  <GritMark className="size-8 text-indigo-600 dark:text-indigo-400" />
                </a>
                <button
                  type="button"
                  onClick={() => setMenuOpen(false)}
                  className="-m-2.5 rounded-md p-2.5 text-gray-700 dark:text-gray-300"
                >
                  <span className="sr-only">Close menu</span>
                  <X aria-hidden="true" className="size-6" />
                </button>
              </div>
              <div className="mt-6 flow-root">
                <div className="-my-6 divide-y divide-gray-500/10">
                  <div className="space-y-2 py-6">
                    {NAV.map((item) => (
                      <a
                        key={item.name}
                        href={item.href}
                        className="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-gray-900 hover:bg-gray-50 dark:text-white dark:hover:bg-gray-800"
                      >
                        {item.name}
                      </a>
                    ))}
                  </div>
                  <div className="py-6">
                    <a
                      href="#"
                      className="-mx-3 block rounded-lg px-3 py-2.5 text-base/7 font-semibold text-gray-900 hover:bg-gray-50 dark:text-white dark:hover:bg-gray-800"
                    >
                      Log in
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </header>

      <div className="relative isolate px-6 pt-14 lg:px-8">
        {/* Decorative blurred gradient. aria-hidden so it is never announced. */}
        <div
          aria-hidden="true"
          className="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
        >
          <div
            style={{
              clipPath:
                'polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)',
            }}
            className="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem] dark:opacity-20"
          />
        </div>

        <div className="mx-auto max-w-2xl py-32 sm:py-48 lg:py-56">
          <div className="hidden sm:mb-8 sm:flex sm:justify-center">
            <div className="relative rounded-full px-3 py-1 text-sm/6 text-gray-600 ring-1 ring-gray-900/10 hover:ring-gray-900/20 dark:text-gray-400 dark:ring-white/10 dark:hover:ring-white/20">
              {announcement}{' '}
              <a href="#" className="font-semibold text-indigo-600 dark:text-indigo-400">
                <span aria-hidden="true" className="absolute inset-0" />
                Read more <span aria-hidden="true">&rarr;</span>
              </a>
            </div>
          </div>

          <div className="text-center">
            <h1 className="text-5xl font-semibold tracking-tight text-balance text-gray-900 sm:text-7xl dark:text-white">
              {title}
            </h1>
            <p className="mt-8 text-lg font-medium text-pretty text-gray-500 sm:text-xl/8 dark:text-gray-400">
              {subtitle}
            </p>
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <a
                href="#"
                className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                {primaryLabel}
              </a>
              <a
                href="#"
                className="group inline-flex items-center gap-1 text-sm/6 font-semibold text-gray-900 dark:text-white"
              >
                {secondaryLabel}
                <ArrowRight
                  aria-hidden="true"
                  className="size-4 transition-transform group-hover:translate-x-0.5"
                />
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
