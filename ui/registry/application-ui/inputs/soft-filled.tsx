'use client'

import * as React from 'react'

/* grit:slot input@1
 *
 * Borderless until you touch it: a filled surface that grows a ring on focus.
 *
 * The border is transparent rather than absent. A field that only gains a border
 * on focus shifts its neighbours by two pixels the moment you click it, and on a
 * dense form the whole column jumps. Keeping the border and changing only its
 * colour costs nothing and holds the layout still.
 */

export type InputSize = 'sm' | 'md' | 'lg'

const SIZES: Record<InputSize, string> = {
  sm: 'h-8 rounded-xl px-3 text-[13px]',
  md: 'h-11 rounded-xl px-4 text-sm',
  lg: 'h-12 rounded-xl px-4 text-[15px]',
}

const MULTILINE_SIZES: Record<InputSize, string> = {
  sm: 'rounded-xl px-3 py-2 text-[13px]',
  md: 'rounded-xl px-4 py-3 text-sm',
  lg: 'rounded-xl px-4 py-3.5 text-[15px]',
}

const BASE =
  'w-full border border-transparent bg-bg-hover text-foreground transition-all ' +
  'placeholder:text-text-muted ' +
  'hover:bg-bg-tertiary ' +
  'focus:bg-bg-secondary focus:outline-none focus:ring-2 ' +
  'disabled:cursor-not-allowed disabled:opacity-60'

export function inputClasses(opts?: {
  inputSize?: InputSize
  invalid?: boolean
  multiline?: boolean
  className?: string
}): string {
  const table = opts?.multiline ? MULTILINE_SIZES : SIZES
  const size = table[opts?.inputSize ?? 'md']
  const state = opts?.invalid
    ? 'border-danger/60 focus:border-danger focus:ring-danger/25'
    : 'focus:border-accent/60 focus:ring-accent/25'
  return [BASE, size, state, opts?.className ?? ''].filter(Boolean).join(' ')
}

export interface InputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  /** NOT "size" — <input size> is a real HTML attribute taking a character count. */
  inputSize?: InputSize
  invalid?: boolean
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(function Input(
  { inputSize, invalid, className, ...rest },
  ref,
) {
  return (
    <input
      ref={ref}
      aria-invalid={invalid || undefined}
      className={inputClasses({ inputSize, invalid, className })}
      {...rest}
    />
  )
})

/* grit:preview-start */
export default function Preview() {
  return (
    <div className="mx-auto flex min-h-full max-w-md flex-col justify-center gap-4 bg-white p-10 dark:bg-gray-950">
      <Input placeholder="Small" inputSize="sm" />
      <Input placeholder="Medium (default)" />
      <Input placeholder="Large" inputSize="lg" />
      <Input placeholder="Invalid" invalid defaultValue="not-an-email" />
      <Input placeholder="Disabled" disabled />
      <textarea
        rows={3}
        placeholder="Textarea shares the same look"
        className={inputClasses({ multiline: true, className: 'resize-y' })}
      />
    </div>
  )
}
/* grit:preview-end */
