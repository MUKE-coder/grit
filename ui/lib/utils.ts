import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * The shadcn `cn` helper, present here for one reason: blocks that declare
 * shadcn primitives in `registryDependencies` import them, and those primitives
 * import this. Without it the previews for those blocks would not compile.
 *
 * Marketing blocks do not use it and should not start. They are authored with
 * stock Tailwind classes and no dependencies, so a copied block renders
 * correctly in any Tailwind project with nothing to install first.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
