// settings-section-01 — Settings card section (Vercel / Linear-inspired)
export function SettingsSection({
  title = "Workspace Name",
  description = "This is the name that will be displayed across your workspace and in emails.",
  children,
  footer,
  danger = false,
}: {
  title?: string; description?: string;
  children?: React.ReactNode; footer?: React.ReactNode; danger?: boolean;
}) {
  return (
    <div className={"rounded-xl border overflow-hidden " + (danger ? "border-danger/40" : "border-border")}>
      <div className="p-6">
        <h3 className={"text-base font-semibold mb-1 " + (danger ? "text-danger" : "text-foreground")}>{title}</h3>
        <p className="text-sm text-text-secondary mb-4">{description}</p>
        {children}
      </div>
      {footer && (
        <div className={"px-6 py-4 border-t flex items-center justify-between " + (danger ? "border-danger/20 bg-danger/5" : "border-border bg-background")}>
          {footer}
        </div>
      )}
    </div>
  );
}

export default SettingsSection;
