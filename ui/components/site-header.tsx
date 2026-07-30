'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { BookOpen, Github } from 'lucide-react'
import { CATALOG } from '@/registry/catalog'
import { ThemeToggle } from './theme-toggle'
import { GritUILogo } from './grit-ui-logo'

export function SiteHeader() {
  const pathname = usePathname() ?? ''

  return (
    <header className="glass-bar sticky top-0 z-50">
      <div className="mx-auto flex h-[52px] max-w-[100rem] items-center gap-7 px-6">
        <Link
          href="/"
          className="shrink-0 transition-transform duration-200 active:scale-[0.97]"
        >
          <GritUILogo size={22} />
          <span className="sr-only">Grit UI home</span>
        </Link>

        <nav className="hidden items-center gap-1 md:flex">
          {CATALOG.map((category) => {
            const active = pathname.startsWith(`/${category.slug}`)
            return (
              <Link
                key={category.slug}
                href={`/${category.slug}`}
                aria-current={active ? 'page' : undefined}
                className={`rounded-lg px-3 py-1.5 text-[13.5px] transition-colors duration-200 ${
                  active
                    ? 'bg-gray-500/[0.10] font-semibold text-gray-900 dark:bg-white/[0.10] dark:text-white'
                    : 'font-medium text-gray-600 hover:bg-gray-500/[0.07] hover:text-gray-900 dark:text-gray-400 dark:hover:bg-white/[0.06] dark:hover:text-white'
                }`}
              >
                {category.name}
              </Link>
            )
          })}
        </nav>

        <div className="ml-auto flex items-center gap-0.5">
          <IconLink
            href="https://gritframework.dev/docs/frontend/ui-components"
            label="Documentation"
            icon={BookOpen}
          />
          <IconLink
            href="https://github.com/MUKE-coder/grit"
            label="GitHub"
            icon={Github}
            external
          />
          <ThemeToggle />
        </div>
      </div>
    </header>
  )
}

function IconLink({
  href,
  label,
  icon: Icon,
  external,
}: {
  href: string
  label: string
  icon: typeof Github
  external?: boolean
}) {
  return (
    <Link
      href={href}
      target={external ? '_blank' : undefined}
      rel={external ? 'noreferrer' : undefined}
      aria-label={label}
      title={label}
      className="inline-flex size-8 items-center justify-center rounded-lg text-gray-500 transition-all duration-200 hover:bg-gray-500/[0.08] hover:text-gray-900 active:scale-[0.94] dark:text-gray-400 dark:hover:bg-white/[0.08] dark:hover:text-white"
    >
      <Icon aria-hidden="true" className="size-[17px]" />
    </Link>
  )
}
