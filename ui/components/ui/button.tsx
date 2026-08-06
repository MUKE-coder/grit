'use client'

import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'

import { cn } from '@/lib/utils'

/**
 * A preview stand-in for shadcn's Button.
 *
 * Blocks that declare `registryDependencies: ['button']` install shadcn's real
 * one into the user's project. This file exists only so those blocks render on
 * this site, and it is authored with stock Tailwind rather than shadcn's theme
 * tokens on purpose.
 *
 * The reason is a collision. This site already maps `accent`, `border` and
 * `foreground` to the admin's design tokens, because the swappable slot
 * previews under Application UI are authored against those names. shadcn uses
 * `accent` for a subtle hover surface; the admin uses it for a bright blue. Add
 * shadcn's token layer on top and the four swappable previews start rendering
 * blue hover states with unreadable text.
 *
 * So the trade is deliberate: previews are styled our way rather than being
 * left unstyled, and the installed article is shadcn's. The shapes, variants
 * and API match, which is what the surrounding block markup depends on.
 */
const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0',
  {
    variants: {
      variant: {
        default: 'bg-indigo-600 text-white shadow-sm hover:bg-indigo-500',
        destructive: 'bg-red-600 text-white shadow-sm hover:bg-red-500',
        outline:
          'border border-gray-300 bg-transparent text-gray-900 hover:bg-gray-50 dark:border-white/15 dark:text-white dark:hover:bg-white/5',
        secondary:
          'bg-gray-100 text-gray-900 hover:bg-gray-200 dark:bg-white/10 dark:text-white dark:hover:bg-white/15',
        ghost:
          'text-gray-900 hover:bg-gray-100 dark:text-white dark:hover:bg-white/10',
        link: 'text-indigo-600 underline-offset-4 hover:underline dark:text-indigo-400',
      },
      size: {
        /* min-h-11 rather than h-9: 44px is the smallest reliable tap target,
           and a button that is comfortable on a phone is not worse on a
           desktop. */
        default: 'min-h-11 px-4 py-2',
        sm: 'min-h-9 rounded-md px-3 text-xs',
        lg: 'min-h-12 rounded-md px-8',
        icon: 'size-11',
      },
    },
    defaultVariants: { variant: 'default', size: 'default' },
  },
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button'
    return (
      <Comp className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
    )
  },
)
Button.displayName = 'Button'

export { Button, buttonVariants }
