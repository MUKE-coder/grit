'use client'

import * as React from 'react'
import { Loader2 } from 'lucide-react'

/* grit:slot button@1
 *
 * Fully rounded, with a ring that expands out of the button on hover and settles.
 *
 * The ring is a ::after pseudo-element on a ring utility, NOT a box-shadow
 * transition — a shadow that animates its spread repaints the element's whole
 * bounding box every frame, and on a form with a dozen buttons that is visible.
 * A scaled pseudo-element is composited.
 *
 * Honours prefers-reduced-motion: the ring still appears, it just does not
 * travel. Removing the affordance entirely would leave those users without the
 * hover feedback everyone else gets.
 */

export type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'
export type ButtonSize = 'sm' | 'md' | 'lg' | 'icon'

const VARIANTS: Record<ButtonVariant, string> = {
  primary:
    'bg-accent text-accent-fg after:ring-accent/40 hover:bg-accent-hover focus-visible:ring-accent',
  secondary:
    'bg-bg-tertiary text-foreground after:ring-border hover:bg-bg-hover focus-visible:ring-border',
  outline:
    'border border-border bg-transparent text-foreground after:ring-border hover:bg-bg-hover focus-visible:ring-border',
  ghost:
    'bg-transparent text-text-secondary after:ring-transparent hover:bg-bg-hover hover:text-foreground focus-visible:ring-border',
  danger: 'bg-danger text-danger-fg after:ring-danger/40 hover:opacity-90 focus-visible:ring-danger',
}

const SIZES: Record<ButtonSize, string> = {
  sm: 'h-8 gap-1.5 rounded-full px-3.5 text-[13px]',
  md: 'h-10 gap-2 rounded-full px-5 text-sm',
  lg: 'h-11 gap-2 rounded-full px-6 text-[15px]',
  icon: 'h-10 w-10 gap-0 rounded-full p-0 text-sm',
}

const BASE =
  'relative isolate inline-flex shrink-0 items-center justify-center font-medium transition-colors ' +
  'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-offset-bg-secondary ' +
  'disabled:pointer-events-none disabled:opacity-50 ' +
  // The ring itself: a transparent overlay that scales up and fades on hover.
  "after:pointer-events-none after:absolute after:inset-0 after:-z-10 after:rounded-full after:ring-2 " +
  'after:opacity-0 after:transition after:duration-300 after:ease-out ' +
  'hover:after:scale-110 hover:after:opacity-100 ' +
  'motion-reduce:after:transition-none motion-reduce:hover:after:scale-100'

export function buttonClasses(opts?: {
  variant?: ButtonVariant
  size?: ButtonSize
  className?: string
}): string {
  const variant = VARIANTS[opts?.variant ?? 'primary']
  const size = SIZES[opts?.size ?? 'md']
  return [BASE, variant, size, opts?.className ?? ''].filter(Boolean).join(' ')
}

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
  loading?: boolean
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { variant, size, loading, disabled, className, children, type, ...rest },
  ref,
) {
  return (
    <button
      ref={ref}
      type={type ?? 'button'}
      disabled={disabled || loading}
      className={buttonClasses({ variant, size, className })}
      {...rest}
    >
      {loading && <Loader2 className="h-4 w-4 animate-spin" />}
      {children}
    </button>
  )
})

/* grit:preview-start */
export default function Preview() {
  return (
    <div className="flex min-h-full flex-col justify-center gap-8 bg-white p-10 dark:bg-gray-950">
      <div className="flex flex-wrap items-center gap-4">
        <Button variant="primary">Save changes</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="outline">Cancel</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="danger">Delete</Button>
      </div>
      <div className="flex flex-wrap items-center gap-4">
        <Button size="sm">Small</Button>
        <Button size="md">Medium</Button>
        <Button size="lg">Large</Button>
        <Button loading>Saving</Button>
        <Button disabled>Disabled</Button>
      </div>
      <p className="text-[13px] text-gray-500 dark:text-gray-400">Hover a button to see the ring.</p>
    </div>
  )
}
/* grit:preview-end */
