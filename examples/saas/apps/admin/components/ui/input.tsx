"use client";

import * as React from "react";

/* grit:slot input@1
 *
 * SWAPPABLE. Replace this whole file with:
 *
 *     grit swap input <variant>
 *
 * Anything exported below is contract — see components/ui/button.tsx for the
 * rules a replacement has to keep.
 */

export type InputSize = "sm" | "md" | "lg";

const SIZES: Record<InputSize, string> = {
  sm: "h-8 rounded-lg px-3 text-[13px]",
  md: "h-10 rounded-lg px-3.5 text-sm",
  lg: "h-11 rounded-lg px-4 text-[15px]",
};

/* A textarea sizes itself by rows, so it gets vertical padding and NO height.
   Passing h-auto alongside h-10 would leave two height utilities fighting, and
   which one wins depends on Tailwind's internal ordering rather than on the
   order they appear in the string — the kind of thing that looks fine until a
   Tailwind upgrade silently reorders it. */
const MULTILINE_SIZES: Record<InputSize, string> = {
  sm: "rounded-lg px-3 py-2 text-[13px]",
  md: "rounded-lg px-3.5 py-2.5 text-sm",
  lg: "rounded-lg px-4 py-3 text-[15px]",
};

const BASE =
  "border bg-bg-tertiary text-foreground transition-colors " +
  "placeholder:text-text-muted " +
  "focus:outline-none focus:ring-1 " +
  "disabled:cursor-not-allowed disabled:opacity-60";

/**
 * The class string on its own — shared by <textarea> and <select> so the whole
 * form keeps one shape. Without this, swapping the input would restyle text
 * fields and leave every dropdown on the old look.
 */
export function inputClasses(opts?: {
  inputSize?: InputSize;
  invalid?: boolean;
  /** Drops the fixed height, for <textarea>. */
  multiline?: boolean;
  /**
   * Width. Full by default, which is what a field in a form column wants.
   *
   * A control in a toolbar wants the opposite: three full-width selects in a
   * wrapping row become three stacked rows, which is what the Users page's
   * filters looked like until this existed. Passing false sizes the control to
   * its content instead, within a sensible range so a long option cannot
   * stretch the row.
   *
   * It is an option rather than a className override because the two would be
   * w-full and w-auto in the same class attribute, and which of those wins is
   * decided by their order in the generated stylesheet rather than here.
   */
  fullWidth?: boolean;
  className?: string;
}): string {
  const table = opts?.multiline ? MULTILINE_SIZES : SIZES;
  const size = table[opts?.inputSize ?? "md"];
  const state = opts?.invalid
    ? "border-danger focus:border-danger focus:ring-danger"
    : "border-border focus:border-accent focus:ring-accent";
  const width = opts?.fullWidth === false ? "w-auto min-w-[7rem] max-w-[14rem]" : "w-full";
  return [BASE, width, size, state, opts?.className ?? ""].filter(Boolean).join(" ");
}

export interface InputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "size"> {
  /**
   * NOT "size". <input size> is a real HTML attribute taking a character count,
   * so a prop called size would both collide with it and silently render
   * size="md" into the DOM.
   */
  inputSize?: InputSize;
  /** Paints the error border. Pair it with aria-invalid for screen readers. */
  invalid?: boolean;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  function Input({ inputSize, invalid, className, ...rest }, ref) {
    return (
      <input
        ref={ref}
        aria-invalid={invalid || undefined}
        className={inputClasses({ inputSize, invalid, className })}
        {...rest}
      />
    );
  },
);
