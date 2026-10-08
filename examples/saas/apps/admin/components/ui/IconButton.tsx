"use client";

import { forwardRef } from "react";
import type { ButtonHTMLAttributes, ReactNode } from "react";

interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Icon node (e.g. <Plus className="h-4 w-4" />). Required. */
  icon: ReactNode;
  /** Text label shown on >=sm screens. The label always serves as the
   *  aria-label on mobile, so screen readers know what the button does. */
  label: string;
  /** Visual variant. Defaults to "primary". */
  variant?: "primary" | "secondary" | "ghost" | "danger";
}

const variantClass: Record<NonNullable<IconButtonProps["variant"]>, string> = {
  primary: "bg-accent text-accent-fg hover:bg-accent-hover",
  secondary: "border border-border bg-bg-elevated text-foreground hover:bg-bg-hover",
  ghost: "text-text-secondary hover:bg-bg-hover hover:text-foreground",
  danger: "bg-danger text-white hover:opacity-90",
};

/**
 * Auto-collapsing CTA. Stays text + icon on >=sm; collapses to icon-only
 * on mobile so table rows + page headers don't blow out of the viewport.
 * Pass label as the readable name; it doubles as the aria-label when the
 * text hides.
 */
export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(function IconButton(
  { icon, label, variant = "primary", className = "", ...rest },
  ref
) {
  return (
    <button
      {...rest}
      ref={ref}
      aria-label={label}
      className={
        "inline-flex items-center justify-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors disabled:opacity-50 " +
        // Mobile is icon-only at 36x36; >=sm reveals the label.
        "h-9 w-9 sm:h-9 sm:w-auto sm:px-3.5 " +
        variantClass[variant] + " " +
        className
      }
    >
      {icon}
      <span className="hidden sm:inline">{label}</span>
    </button>
  );
});
