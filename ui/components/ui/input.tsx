import * as React from 'react'

import { cn } from '@/lib/utils'

/**
 * A preview stand-in for shadcn's Input. See components/ui/button.tsx for why
 * these are authored with stock Tailwind rather than shadcn's theme tokens.
 *
 * `min-h-11` rather than `h-9`: a 44px target is the smallest that is reliably
 * tappable, and on a login form the input is the thing people are aiming at.
 */
const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<'input'>>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          'flex min-h-11 w-full rounded-md border border-gray-300 bg-transparent px-3 py-2 text-sm text-gray-900 shadow-sm transition-colors',
          'placeholder:text-gray-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600',
          'disabled:cursor-not-allowed disabled:opacity-50',
          'file:border-0 file:bg-transparent file:text-sm file:font-medium',
          'dark:border-white/15 dark:text-white dark:placeholder:text-gray-600',
          className,
        )}
        ref={ref}
        {...props}
      />
    )
  },
)
Input.displayName = 'Input'

export { Input }
