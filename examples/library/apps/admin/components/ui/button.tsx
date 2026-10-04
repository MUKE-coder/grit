"use client";

import * as React from "react";
// lucide-react directly, NOT @/lib/icons. A swapped-in variant is authored
// against the public registry where @/lib/icons does not exist, and @/lib/icons
// only re-exports a curated subset — a variant reaching for an icon that is in
// the iconMap but not the export block fails to compile. Importing the package
// keeps default and variant byte-identical in shape and needs no rewrite on swap.
import { Loader2 } from "lucide-react";

/* grit:slot button@1
 *
 * SWAPPABLE. Replace this whole file with:
 *
 *     grit swap button <variant>
 *
 * Anything exported below is contract. Call sites across the admin rely on it,
 * so a replacement must keep every variant and size working — including the
 * ones you personally would not use. Run `grit swap --check` to verify.
 */

export type ButtonVariant =
  | "primary"
  | "secondary"
  | "outline"
  | "ghost"
  | "danger";

export type ButtonSize = "sm" | "md" | "lg" | "icon";

const VARIANTS: Record<ButtonVariant, string> = {
  primary:
    "bg-accent text-accent-fg hover:bg-accent-hover focus-visible:ring-accent",
  secondary:
    "bg-bg-tertiary text-foreground hover:bg-bg-hover focus-visible:ring-border",
  outline:
    "border border-border bg-transparent text-foreground hover:bg-bg-hover focus-visible:ring-border",
  ghost:
    "bg-transparent text-text-secondary hover:bg-bg-hover hover:text-foreground focus-visible:ring-border",
  danger:
    "bg-danger text-white hover:opacity-90 focus-visible:ring-danger",
};

const SIZES: Record<ButtonSize, string> = {
  sm: "h-8 gap-1.5 rounded-lg px-3 text-[13px]",
  md: "h-10 gap-2 rounded-lg px-4 text-sm",
  lg: "h-11 gap-2 rounded-lg px-5 text-[15px]",
  // Square. Width is pinned so an icon-only button never collapses to its
  // glyph and becomes a 16px tap target.
  icon: "h-10 w-10 gap-0 rounded-lg p-0 text-sm",
};

const BASE =
  "inline-flex shrink-0 items-center justify-center font-medium transition-colors " +
  "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-offset-bg-secondary " +
  "disabled:pointer-events-none disabled:opacity-50";

/**
 * The class string on its own.
 *
 * Exported because a good third of the controls in the admin are <a> or <label>
 * that need to look like buttons. They call this instead of wrapping a Button,
 * which is what keeps a swap total rather than leaving links on the old style.
 */
export function buttonClasses(opts?: {
  variant?: ButtonVariant;
  size?: ButtonSize;
  className?: string;
}): string {
  const variant = VARIANTS[opts?.variant ?? "primary"];
  const size = SIZES[opts?.size ?? "md"];
  return [BASE, variant, size, opts?.className ?? ""].filter(Boolean).join(" ");
}

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  /** Shows a spinner and disables the button. */
  loading?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  function Button(
    { variant, size, loading, disabled, className, children, type, ...rest },
    ref,
  ) {
    return (
      <button
        ref={ref}
        // Defaulting to "button" is deliberate. The HTML default is "submit",
        // and an unlabelled button inside a form submits it — which is how a
        // Cancel next to a Save ends up saving.
        type={type ?? "button"}
        disabled={disabled || loading}
        className={buttonClasses({ variant, size, className })}
        {...rest}
      >
        {loading && <Loader2 className="h-4 w-4 animate-spin" />}
        {children}
      </button>
    );
  },
);
