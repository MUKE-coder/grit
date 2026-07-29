"use client";

// alert-banner-01 — Alert/warning/success/error banner (common UI pattern)
export function AlertBanner({
  variant = "info",
  title,
  message = "Something needs your attention.",
  dismissible = true,
  onDismiss,
  action,
}: {
  variant?: "info" | "success" | "warning" | "error";
  title?: string; message?: string; dismissible?: boolean;
  onDismiss?: () => void;
  action?: { label: string; onClick: () => void };
}) {
  const config = {
    info: { bg: "bg-accent/10", border: "border-accent/30", icon: "text-accent", path: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" },
    success: { bg: "bg-success/10", border: "border-success/30", icon: "text-success", path: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" },
    warning: { bg: "bg-yellow-500/10", border: "border-yellow-500/30", icon: "text-yellow-400", path: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" },
    error: { bg: "bg-danger/10", border: "border-danger/30", icon: "text-danger", path: "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" },
  };
  const c = config[variant];
  return (
    <div className={"flex items-start gap-3 p-4 rounded-xl border " + c.bg + " " + c.border}>
      <svg className={"w-5 h-5 flex-shrink-0 mt-0.5 " + c.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={c.path} />
      </svg>
      <div className="flex-1 min-w-0">
        {title && <p className="text-sm font-semibold text-foreground mb-0.5">{title}</p>}
        <p className="text-sm text-text-secondary">{message}</p>
        {action && (
          <button onClick={action.onClick} className={"mt-2 text-sm font-semibold hover:underline " + c.icon}>{action.label}</button>
        )}
      </div>
      {dismissible && (
        <button onClick={onDismiss} className="p-1 rounded text-text-muted hover:text-foreground transition-colors flex-shrink-0">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      )}
    </div>
  );
}

export default AlertBanner;
