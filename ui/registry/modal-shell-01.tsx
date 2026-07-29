"use client";

// modal-shell-01 — Modal dialog shell (Headless UI / Radix-inspired pattern)
import { useEffect } from "react";
export function ModalShell({
  title = "Confirm action",
  description,
  children,
  footer,
  onClose,
  size = "md",
}: {
  title?: string; description?: string;
  children?: React.ReactNode; footer?: React.ReactNode;
  onClose?: () => void; size?: "sm" | "md" | "lg" | "xl";
}) {
  const sizeClasses = { sm: "max-w-sm", md: "max-w-md", lg: "max-w-lg", xl: "max-w-2xl" };
  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === "Escape") onClose?.(); };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-background/60 backdrop-blur-sm" onClick={onClose} />
      <div className={"relative w-full " + sizeClasses[size] + " bg-bg-elevated border border-border rounded-2xl shadow-2xl overflow-hidden"}>
        <div className="flex items-start justify-between p-6 border-b border-border">
          <div>
            <h2 className="text-base font-semibold text-foreground">{title}</h2>
            {description && <p className="text-sm text-text-secondary mt-1">{description}</p>}
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors -mr-1">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        {children && <div className="p-6">{children}</div>}
        {footer && <div className="px-6 py-4 border-t border-border flex items-center justify-end gap-3">{footer}</div>}
      </div>
    </div>
  );
}

export default ModalShell;
