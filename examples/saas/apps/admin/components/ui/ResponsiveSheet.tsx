"use client";

import { useEffect, useRef } from "react";
import type { ReactNode } from "react";
import { X } from "@/lib/icons";

interface ResponsiveSheetProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: ReactNode;
  /** Footer rendered at the bottom (typically Cancel + Submit). */
  footer?: ReactNode;
  /** Max width on desktop. Defaults to 'lg' (~36rem). */
  size?: "sm" | "md" | "lg" | "xl";
}

const sizeClass: Record<NonNullable<ResponsiveSheetProps["size"]>, string> = {
  sm: "md:max-w-sm",
  md: "md:max-w-md",
  lg: "md:max-w-lg",
  xl: "md:max-w-2xl",
};

/**
 * Adapts modal style to viewport. Desktop (>=md): a right-anchored sheet
 * that slides in from the right edge and spans full height — keeps the
 * dashboard context visible behind it (Walkie-Check style). Mobile: a
 * bottom-anchored sheet that slides up and stops at 90vh. Both lock body
 * scroll when open and close on backdrop click + Escape.
 */
export function ResponsiveSheet({
  open,
  onClose,
  title,
  description,
  children,
  footer,
  size = "lg",
}: ResponsiveSheetProps) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    document.addEventListener("keydown", onKey);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = "";
    };
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center md:items-stretch md:justify-end"
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />

      {/* Panel — mobile: bottom sheet (rounded top, capped 90vh).
          Desktop: right drawer (full height, rounded left, sized by prop). */}
      <div
        ref={ref}
        role="dialog"
        aria-modal="true"
        aria-labelledby="responsive-sheet-title"
        className={
          "relative z-10 flex w-full flex-col bg-bg-elevated text-foreground shadow-2xl " +
          "rounded-t-2xl md:rounded-none md:rounded-l-2xl " +
          "max-h-[90vh] md:max-h-none md:h-full " +
          ("md:w-full " + sizeClass[size])
        }
      >
        <header className="flex items-start justify-between border-b border-border px-5 py-4">
          <div className="min-w-0">
            <h2 id="responsive-sheet-title" className="text-lg font-semibold truncate">{title}</h2>
            {description && <p className="mt-0.5 text-sm text-text-secondary">{description}</p>}
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close"
            className="ml-3 inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted hover:bg-bg-hover hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </header>

        <div className="flex-1 overflow-y-auto px-5 py-4">{children}</div>

        {footer && (
          <footer className="flex items-center justify-end gap-2 border-t border-border px-5 py-3">
            {footer}
          </footer>
        )}
      </div>
    </div>
  );
}
