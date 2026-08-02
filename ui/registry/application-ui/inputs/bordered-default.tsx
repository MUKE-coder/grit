'use client'

import * as React from 'react'

/* grit:slot input@1
 *
 * The stock Grit input. Published so `grit swap input bordered-default` can put
 * it back — a swap system with no way home is a one-way door.
 */

export type InputSize = 'sm' | 'md' | 'lg'

const SIZES: Record<InputSize, string> = {
  sm: 'h-8 rounded-lg px-3 text-[13px]',
  md: 'h-10 rounded-lg px-3.5 text-sm',
  lg: 'h-11 rounded-lg px-4 text-[15px]',
}

/* A textarea sizes itself by rows, so it gets padding and NO height. Two height
   utilities on one element resolve by Tailwind's internal ordering rather than
   by the order they appear in the string. */
const MULTILINE_SIZES: Record<InputSize, string> = {
  sm: 'rounded-lg px-3 py-2 text-[13px]',
  md: 'rounded-lg px-3.5 py-2.5 text-sm',
  lg: 'rounded-lg px-4 py-3 text-[15px]',
}

const BASE =
  'w-full border bg-bg-tertiary text-foreground transition-colors ' +
  'placeholder:text-text-muted ' +
  'focus:outline-none focus:ring-1 ' +
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
    ? 'border-danger focus:border-danger focus:ring-danger'
    : 'border-border focus:border-accent focus:ring-accent'
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
