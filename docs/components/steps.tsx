import type { ReactNode } from 'react'

/**
 * Numbered steps, in the design's two-column gutter: a 32px rail carrying a
 * circled number and the line that joins one step to the next, then the
 * content.
 *
 * The numbers come from a CSS counter rather than a prop, so inserting a
 * step in the middle cannot leave the rest renumbered by hand and wrong.
 *
 *   <Steps>
 *     <Step title="Install the CLI">…</Step>
 *     <Step title="Create a project">…</Step>
 *   </Steps>
 *
 * The connector is drawn on each Step rather than on the container, so the
 * last one can stop at its own number instead of running on into the
 * whitespace below the final step.
 */

export function Steps({ children }: { children: ReactNode }) {
  return <div className="my-8 [counter-reset:step]">{children}</div>
}

export function Step({
  title,
  children,
  toc = true,
}: {
  title?: ReactNode
  children: ReactNode
  /**
   * Whether this step's title is a section of the page.
   *
   * True when the steps are the page: "Create a project" is somewhere a
   * reader navigates to, so it belongs in the table of contents and is set at
   * section size. False when the steps are a worked example inside a section
   * that already has its own heading, which is every plugin page: nine of
   * them listing "Install it" in the sidebar is noise, not navigation.
   *
   * It is one prop rather than two because it is one question. A step that is
   * not worth navigating to is not worth setting at 22px either.
   */
  toc?: boolean
}) {
  return (
    <div className="group grid grid-cols-[28px_minmax(0,1fr)] gap-x-5 [counter-increment:step]">
      <div aria-hidden="true" className="flex flex-col items-center">
        <span className="grid h-7 w-7 flex-none place-items-center rounded-full border border-border-strong bg-bg font-mono text-[12px] font-medium text-fg before:[content:counter(step)]" />
        {/* The rule grows to the height of the step beside it, and the last
            step has none: there is nothing below it to join to. group-last
            rather than last, because the rule is not the last child of
            anything, the step it belongs to is. */}
        <span className="mt-2 w-px flex-1 bg-border group-last:hidden" />
      </div>

      <div className={toc ? 'min-w-0 pb-10' : 'min-w-0 pb-8 group-last:pb-0'}>
        {title && (
          <h3
            {...(toc ? { 'data-toc': true } : {})}
            className={
              toc
                ? 'mt-0.5 text-[22px] font-semibold leading-[1.3] tracking-[-0.02em] text-fg'
                : 'mt-1 text-[17px] font-semibold leading-snug tracking-[-0.01em] text-fg'
            }
          >
            {title}
          </h3>
        )}
        <div className="mt-3 text-[15px] leading-[1.7] text-fg-muted [&>p]:max-w-[var(--measure)] [&_a:hover]:underline [&_a]:text-accent">
          {children}
        </div>
      </div>
    </div>
  )
}
