"use client";

// toast-notification-01 — Toast notification (Sonner / React Hot Toast-inspired)
export function Toast({
  message = "Changes saved successfully",
  variant = "success",
  action,
  onDismiss,
}: {
  message?: string; variant?: "success" | "error" | "info" | "warning";
  action?: { label: string; onClick: () => void }; onDismiss?: () => void;
}) {
  const config = {
    success: { icon: "text-success bg-success/20", path: "M5 13l4 4L19 7" },
    error: { icon: "text-danger bg-danger/20", path: "M6 18L18 6M6 6l12 12" },
    info: { icon: "text-accent bg-accent/20", path: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" },
    warning: { icon: "text-yellow-400 bg-yellow-500/20", path: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" },
  };
  const c = config[variant];
  return (
    <div className="flex items-center gap-3 bg-bg-elevated border border-border rounded-xl px-4 py-3 shadow-lg min-w-64 max-w-sm">
      <div className={"w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0 " + c.icon}>
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={c.path} />
        </svg>
      </div>
      <span className="text-sm text-foreground flex-1">{message}</span>
      {action && (
        <button onClick={action.onClick} className="text-xs text-accent hover:underline font-medium flex-shrink-0">{action.label}</button>
      )}
      {onDismiss && (
        <button onClick={onDismiss} className="p-0.5 rounded text-text-muted hover:text-foreground transition-colors flex-shrink-0">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      )}
    </div>
  );
}

export default Toast;
