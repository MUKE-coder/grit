import Link from 'next/link'

/**
 * Button — the design system's `g-btn`.
 *
 * Primary is INVERTED: a near-white fill with dark text, not the brand
 * violet. That is what lets violet mean "brand" everywhere else instead of
 * meaning "there is a button here". One primary per view.
 */
export type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'accent'
export type ButtonSize = 'sm' | 'md' | 'lg'

const VARIANT: Record<ButtonVariant, string> = {
  primary: 'g-btn--primary',
  secondary: 'g-btn--secondary',
  ghost: 'g-btn--ghost',
  accent: 'g-btn--accent',
}

const SIZE: Record<ButtonSize, string> = {
  sm: 'g-btn--sm',
  md: '',
  lg: 'g-btn--lg',
}

export function buttonClass(
  variant: ButtonVariant = 'secondary',
  size: ButtonSize = 'md',
  className = ''
) {
  return ['g-btn', VARIANT[variant], SIZE[size], className].filter(Boolean).join(' ')
}

type Common = {
  variant?: ButtonVariant
  size?: ButtonSize
  className?: string
  children: React.ReactNode
}

type AsLink = Common & { href: string; external?: boolean } & Omit<
    React.AnchorHTMLAttributes<HTMLAnchorElement>,
    'className' | 'href'
  >

type AsButton = Common & { href?: undefined } & Omit<
    React.ButtonHTMLAttributes<HTMLButtonElement>,
    'className'
  >

export function Button(props: AsLink | AsButton) {
  const { variant = 'secondary', size = 'md', className, children } = props

  if (props.href !== undefined) {
    const { href, external, variant: _v, size: _s, className: _c, ...rest } = props as AsLink
    const cls = buttonClass(variant, size, className)
    // An external link gets a real anchor; next/link is for routes.
    if (external || href.startsWith('http') || href.startsWith('mailto:')) {
      return (
        <a className={cls} href={href} target="_blank" rel="noreferrer" {...rest}>
          {children}
        </a>
      )
    }
    return (
      <Link className={cls} href={href} {...rest}>
        {children}
      </Link>
    )
  }

  const { variant: _v, size: _s, className: _c, ...rest } = props as AsButton
  return (
    <button type="button" className={buttonClass(variant, size, className)} {...rest}>
      {children}
    </button>
  )
}

/**
 * The arrow that nudges on hover. Separate because the animation is on a
 * child, so a button without an arrow pays nothing for it.
 */
export function ButtonArrow() {
  return (
    <svg
      className="g-btn__arrow"
      viewBox="0 0 16 16"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.5"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M3 8h10M9 4l4 4-4 4" />
    </svg>
  )
}
