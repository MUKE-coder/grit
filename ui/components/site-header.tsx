import Link from 'next/link'
import { Github } from 'lucide-react'
import { CATALOG } from '@/registry/catalog'
import { ThemeToggle } from './theme-toggle'

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 border-b border-gray-200 bg-white/85 backdrop-blur dark:border-white/10 dark:bg-gray-950/85">
      <div className="mx-auto flex h-14 max-w-[100rem] items-center gap-6 px-6">
        <Link href="/" className="flex shrink-0 items-center gap-2">
          <span className="flex size-6 items-center justify-center rounded-md bg-gray-900 text-[11px] font-bold text-white dark:bg-white dark:text-gray-900">
            G
          </span>
          <span className="text-sm font-semibold tracking-tight">
            grit<span className="text-indigo-600 dark:text-indigo-400">UI</span>
          </span>
        </Link>

        <nav className="hidden items-center gap-5 md:flex">
          {CATALOG.map((category) => (
            <Link
              key={category.slug}
              href={`/${category.slug}`}
              className="text-sm text-gray-600 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
            >
              {category.name}
            </Link>
          ))}
        </nav>

        <div className="ml-auto flex items-center gap-1">
          <Link
            href="https://gritframework.dev/docs/frontend/ui-components"
            className="hidden px-2 text-sm text-gray-600 transition-colors hover:text-gray-900 sm:block dark:text-gray-400 dark:hover:text-white"
          >
            Docs
          </Link>
          <Link
            href="https://github.com/MUKE-coder/grit"
            target="_blank"
            rel="noreferrer"
            aria-label="GitHub"
            className="inline-flex size-8 items-center justify-center rounded-md text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-white/10 dark:hover:text-white"
          >
            <Github className="size-4" />
          </Link>
          <ThemeToggle />
        </div>
      </div>
    </header>
  )
}
