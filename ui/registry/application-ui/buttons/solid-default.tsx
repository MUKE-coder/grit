'use client'

import * as React from 'react'
import { Loader2 } from 'lucide-react'

/* grit:slot button@1
 *
 * The stock Grit button — the one a new project ships with. It is published as
 * a variant so `grit swap button solid-default` can put it back after trying
 * something else. A swap system with no way home is a one-way door.
 */

export type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'
export type ButtonSize = 'sm' | 'md' | 'lg' | 'icon'

const VARIANTS: Record<ButtonVariant, string> = {
  primary: 'bg-accent text-accent-fg hover:bg-accent-hover focus-visible:ring-accent',
  secondary: 'bg-bg-tertiary text-foreground hover:bg-bg-hover focus-visible:ring-border',
  outline:
    'border border-border bg-transparent text-foreground hover:bg-bg-hover focus-visible:ring-border',
  ghost:
    'bg-transparent text-text-secondary hover:bg-bg-hover hover:text-foreground focus-visible:ring-border',
  danger: 'bg-danger text-danger-fg hover:opacity-90 focus-visible:ring-danger',
}

const SIZES: Record<ButtonSize, string> = {
  sm: 'h-8 gap-1.5 rounded-lg px-3 text-[13px]',
  md: 'h-10 gap-2 rounded-lg px-4 text-sm',
  lg: 'h-11 gap-2 rounded-lg px-5 text-[15px]',
  icon: 'h-10 w-10 gap-0 rounded-lg p-0 text-sm',
}

const BASE =
  'inline-flex shrink-0 items-center justify-center font-medium transition-colors ' +
  'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-offset-bg-secondary ' +
  'disabled:pointer-events-none disabled:opacity-50'

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

/* grit:preview-start
   Everything below is the demo the registry site renders. `grit swap` strips it,
   so the file that lands in your project contains only the contract above. */
export default function Preview() {
  return (
    <div className="flex min-h-full flex-col justify-center gap-8 bg-white p-10 dark:bg-gray-950">
      <div className="flex flex-wrap items-center gap-3">
        <Button variant="primary">Save changes</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="outline">Cancel</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="danger">Delete</Button>
      </div>
      <div className="flex flex-wrap items-center gap-3">
        <Button size="sm">Small</Button>
        <Button size="md">Medium</Button>
        <Button size="lg">Large</Button>
        <Button loading>Saving</Button>
        <Button disabled>Disabled</Button>
      </div>
    </div>
  )
}
/* grit:preview-end */
