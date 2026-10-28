import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { DocsAutoToc } from '@/components/docs-auto-toc'

/**
 * The docs shell: three columns inside 1440px.
 *
 * 264 sidebar · 760 content · 232 on-this-page, from the design. The whole
 * row is centred rather than pinned to the viewport edge, which is the main
 * thing that changes on a wide screen: the old arrangement fixed the sidebar
 * to the left of the window, so at 2560px the navigation sat a metre away
 * from the text it belonged to.
 *
 * The TOC hides below 1180px and the sidebar becomes a drawer below 1024,
 * which is what the two `variant` props on those components switch to.
 *
 * ── Why this is a component and not app/docs/layout.tsx ──────────────────
 *
 * 232 docs pages render their own SiteHeader, DocsSidebar and <main> today.
 * Putting the shell in the layout would give every one of them a second
 * header until all 232 are converted, so pages opt in one at a time and the
 * layout keeps passing children through. When the last page is converted
 * this moves into the layout and the `variant` props go away.
 */
export function DocsShell({
  children,
  /** Rendered above the content, inside the 760px measure. */
  header,
}: {
  children: React.ReactNode
  header?: React.ReactNode
}) {
  return (
    <div className="min-h-screen bg-bg">
      {/* Tells the floating rail that app/docs/layout.tsx mounts to stand
          down: this shell renders its own in the third column. */}
      <span data-docs-shell hidden />
      <SiteHeader />

      <div className="mx-auto flex max-w-[1440px] flex-wrap items-start">
        <DocsSidebar variant="inline" />

        {/* flex-[999] so the content takes every pixel the two rails do not,
            and min-w-0 so a wide code block scrolls rather than widening the
            column and pushing the TOC off the screen. */}
        <main className="min-w-0 flex-[999_1_560px] px-6 pb-24 pt-12 md:px-10 lg:px-[4.5vw] xl:px-[72px]">
          {/* The measure, applied to the column rather than chased paragraph
              by paragraph. 760px of Geist at 15px is about 104 characters,
              and every paragraph that is not inside a .prose-grit block was
              running at that. Code blocks, tables and cards keep the full
              column; only the prose inside them narrows. */}
          <div className="mx-auto max-w-[760px] [&_p]:max-w-[var(--measure)]">
            {header}
            {children}
          </div>
        </main>

        <DocsAutoToc variant="inline" />
      </div>
    </div>
  )
}
