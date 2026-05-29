import { Search, Plus, Inbox, RefreshCw, ArrowLeft } from 'lucide-react'
import { useQueryClient } from '@tanstack/react-query'
import { createContext, useContext, useState } from 'react'
import { useRouter, useNavigate } from '@tanstack/react-router'
import { cn } from '@/lib/utils'

/**
 * Two-pane layout primitives — list on the left, detail on the right.
 * Mirrors WhatsApp Web / Outlook / rental-manager patterns. Used for every
 * CRUD-heavy resource: motorcycles, borrowers, loans, repayments, etc.
 *
 * Responsive behaviour:
 * - lg+ (≥1024px): both panes visible side-by-side as before.
 * - mobile: only ONE pane visible at a time. The page passes
 *   `mobileMode='list' | 'detail'` to TwoPane, which the panes consume via
 *   context to decide who's hidden. DetailPane shows a back button on
 *   mobile that calls `onMobileBack` to clear the page's selection state.
 */

interface TwoPaneCtx {
  /** Which pane is showing on mobile. lg+ ignores this and shows both. */
  mobileMode: 'list' | 'detail'
  /** Called when the user taps the back arrow in DetailPane on mobile. */
  onMobileBack?: () => void
}

const TwoPaneContext = createContext<TwoPaneCtx>({ mobileMode: 'list' })

/**
 * Outer wrapper. Place inside the app's main content slot.
 *
 * If `mobileMode` is omitted, TwoPane auto-detects it from the URL:
 * any page using the project-wide `?selected=ID` convention will switch
 * to detail mode on mobile when a row is selected, and the auto-back
 * button (in DetailPane) clears the selected param. This means every
 * existing list page becomes mobile-responsive without per-page edits.
 */
export function TwoPane({
  children,
  mobileMode: mobileModeProp,
  onMobileBack: onMobileBackProp,
}: {
  children: React.ReactNode
  mobileMode?: 'list' | 'detail'
  onMobileBack?: () => void
}) {
  const router = useRouter()
  const navigate = useNavigate()
  const search = (router.state.location.search ?? {}) as Record<string, unknown>
  const inferredMode: 'list' | 'detail' = mobileModeProp
    ?? ((search.selected !== undefined && search.selected !== null && search.selected !== '') ? 'detail' : 'list')
  const onMobileBack = onMobileBackProp ?? (() => {
    navigate({
      to: router.state.location.pathname as any,
      search: (prev: any) => ({ ...prev, selected: undefined }),
      replace: true,
    } as any)
  })
  return (
    <TwoPaneContext.Provider value={{ mobileMode: inferredMode, onMobileBack }}>
      <div className="flex h-full overflow-hidden">{children}</div>
    </TwoPaneContext.Provider>
  )
}

/** Sticky-header list column with title, search, optional filters, and a "+ New" button. */
export function ListPane({
  title,
  count,
  onNew,
  newLabel = 'New',
  search,
  onSearch,
  searchPlaceholder = 'Search...',
  filters,
  children,
  footer,
  refreshKey,
  extraAction,
}: {
  title: string
  count?: number
  onNew?: () => void
  newLabel?: string
  search: string
  onSearch: (v: string) => void
  searchPlaceholder?: string
  filters?: React.ReactNode
  children: React.ReactNode
  footer?: React.ReactNode
  refreshKey?: unknown[]
  extraAction?: React.ReactNode
}) {
  const { mobileMode } = useContext(TwoPaneContext)
  return (
    <aside
      className={cn(
        'shrink-0 border-r border-border bg-surface flex flex-col h-full',
        // Mobile: full width when in list mode, hidden when in detail mode.
        // Desktop: fixed list-pane width, always visible.
        'w-full lg:w-listpane',
        mobileMode === 'detail' && 'hidden lg:flex',
      )}
    >
      <div className="px-4 pt-4 pb-2 space-y-3 border-b border-border-subtle">
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-baseline gap-2 min-w-0">
            <h2 className="text-[15px] font-semibold text-foreground truncate">{title}</h2>
            {count !== undefined && (
              <span className="text-[12px] text-foreground-muted">{count}</span>
            )}
          </div>
          <div className="flex items-center gap-1.5">
            {extraAction}
            <RefreshButton queryKey={refreshKey} />
            {onNew && (
              <button
                type="button"
                onClick={onNew}
                className="h-8 px-2.5 rounded-lg bg-accent text-white text-[12px] font-medium hover:bg-accent-hover transition-colors inline-flex items-center gap-1"
              >
                <Plus className="h-3.5 w-3.5" />
                {newLabel}
              </button>
            )}
          </div>
        </div>

        <div className="relative">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-foreground-muted" />
          <input
            type="search"
            value={search}
            onChange={(e) => onSearch(e.target.value)}
            placeholder={searchPlaceholder}
            className="w-full h-8 pl-8 pr-2.5 rounded-lg border border-border bg-surface-2 text-[12.5px] placeholder:text-foreground-muted focus:border-accent focus:bg-surface focus:outline-none focus:ring-2 focus:ring-accent/15"
          />
        </div>

        {filters}
      </div>
      <div className="flex-1 overflow-y-auto">{children}</div>
      {footer && <div className="px-4 py-2 border-t border-border-subtle bg-surface-2 text-[12px] text-foreground-muted">{footer}</div>}
    </aside>
  )
}

/** RefreshButton — invalidate a list's queries and re-fetch on click. */
export function RefreshButton({ queryKey }: { queryKey?: unknown[] }) {
  const qc = useQueryClient()
  const [spinning, setSpinning] = useState(false)
  const onClick = async () => {
    if (!queryKey) return
    setSpinning(true)
    await qc.invalidateQueries({ queryKey })
    setTimeout(() => setSpinning(false), 400)
  }
  return (
    <button
      type="button"
      onClick={onClick}
      title="Refresh"
      className="h-8 w-8 rounded-lg border border-border bg-surface text-foreground-muted hover:bg-surface-hover hover:text-foreground transition flex items-center justify-center"
    >
      <RefreshCw className={cn('h-3.5 w-3.5', spinning && 'animate-spin')} />
    </button>
  )
}

/** Standard list row with avatar/icon, title, subtitle, and right-aligned meta. */
export function ListRow({
  onClick,
  selected,
  icon,
  title,
  subtitle,
  rightTop,
  rightBottom,
}: {
  onClick?: () => void
  selected?: boolean
  icon?: React.ReactNode
  title: React.ReactNode
  subtitle?: React.ReactNode
  rightTop?: React.ReactNode
  rightBottom?: React.ReactNode
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        'relative w-full flex items-start gap-3 py-3 px-4 text-left transition-colors border-b border-border-subtle',
        'hover:bg-surface-hover',
        selected && 'bg-accent-tint',
      )}
    >
      {selected && (
        <span aria-hidden className="absolute inset-y-0 left-0 w-[2px] bg-accent" />
      )}
      {icon && (
        <div className="h-9 w-9 rounded-full bg-surface-2 shrink-0 flex items-center justify-center text-[12px] font-semibold text-foreground-secondary overflow-hidden">
          {icon}
        </div>
      )}
      <div className="flex-1 min-w-0">
        <div className={cn('text-[13.5px] font-medium truncate', selected ? 'text-accent' : 'text-foreground')}>
          {title}
        </div>
        {subtitle && (
          <div className="text-[12px] text-foreground-muted truncate mt-0.5">
            {subtitle}
          </div>
        )}
      </div>
      {(rightTop || rightBottom) && (
        <div className="flex flex-col items-end gap-1 shrink-0">
          {rightTop && <div className="text-[11.5px] text-foreground-muted">{rightTop}</div>}
          {rightBottom}
        </div>
      )}
    </button>
  )
}

/** Detail column on the right. Pass `empty` to render the placeholder. */
export function DetailPane({
  header,
  children,
  empty,
  emptyTitle = 'Nothing selected',
  emptyHint = 'Pick an item from the list, or create a new one.',
}: {
  header?: React.ReactNode
  children?: React.ReactNode
  empty?: boolean
  emptyTitle?: string
  emptyHint?: string
}) {
  const { mobileMode, onMobileBack } = useContext(TwoPaneContext)
  // On mobile, the empty state lives in the list pane (the user simply
  // hasn't tapped a row yet) so we hide the empty DetailPane entirely.
  // On desktop, both panes always show, so we keep the empty state.
  if (empty) {
    return (
      <section className={cn(
        'flex-1 items-center justify-center bg-background',
        'hidden lg:flex',
      )}>
        <EmptyState title={emptyTitle} hint={emptyHint} />
      </section>
    )
  }
  return (
    <section
      className={cn(
        'flex-1 flex-col min-w-0 bg-background overflow-hidden',
        // Mobile: only show in detail mode
        mobileMode === 'list' ? 'hidden lg:flex' : 'flex',
      )}
    >
      {header && (
        <div className="border-b border-border bg-surface px-4 lg:px-6 py-3 lg:py-4 shrink-0">
          {/* Mobile back button — visible only when DetailPane is the active
              pane on small screens. Calls the page-level callback that
              clears the URL search-param selection. */}
          {onMobileBack && (
            <button
              type="button"
              onClick={onMobileBack}
              className="lg:hidden mb-2 -ml-1.5 inline-flex items-center gap-1 px-1.5 py-0.5 rounded-md text-[12.5px] text-foreground-muted hover:text-foreground hover:bg-surface-hover transition"
              aria-label="Back to list"
            >
              <ArrowLeft size={14} />
              Back
            </button>
          )}
          {header}
        </div>
      )}
      <div className="flex-1 overflow-y-auto p-4 lg:p-6">{children}</div>
    </section>
  )
}

export function EmptyState({
  title,
  hint,
  action,
}: {
  title: string
  hint?: string
  action?: React.ReactNode
}) {
  return (
    <div className="flex flex-col items-center text-center max-w-sm px-6">
      <div className="h-14 w-14 rounded-full bg-surface-2 flex items-center justify-center mb-4">
        <Inbox className="h-6 w-6 text-foreground-muted" />
      </div>
      <div className="text-[14px] font-semibold text-foreground">{title}</div>
      {hint && <div className="text-[13px] text-foreground-secondary mt-1">{hint}</div>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}

/** Section heading inside a DetailPane. */
export function DetailSection({
  title,
  action,
  children,
}: {
  title: string
  action?: React.ReactNode
  children: React.ReactNode
}) {
  return (
    <section className="mb-6">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
          {title}
        </h3>
        {action}
      </div>
      {children}
    </section>
  )
}

/** Labelled value row used heavily in detail panes. */
export function Field({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div>
      <div className="text-[11px] font-medium uppercase tracking-wider text-foreground-muted mb-0.5">
        {label}
      </div>
      <div className="text-[13.5px] text-foreground">
        {value || <span className="text-foreground-muted">—</span>}
      </div>
    </div>
  )
}

/** Grid of Field rows — typical detail-pane layout. */
export function FieldGrid({
  columns = 2,
  children,
}: {
  columns?: 1 | 2 | 3
  children: React.ReactNode
}) {
  const cls =
    columns === 1
      ? 'grid grid-cols-1 gap-4'
      : columns === 3
      ? 'grid grid-cols-1 sm:grid-cols-3 gap-4'
      : 'grid grid-cols-1 sm:grid-cols-2 gap-4'
  return <div className={cls}>{children}</div>
}
