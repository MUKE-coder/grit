"use client";

// sso-button-01 — SSO provider sign-in button (Okta / Google Workspace-inspired)
export function SSOButton({
  provider = "Google Workspace",
  orgName = "Acme Corp",
  orgDomain = "acme.com",
  logoChar = "G",
  logoColor = "bg-red-500",
  onClick,
}: {
  provider?: string; orgName?: string; orgDomain?: string;
  logoChar?: string; logoColor?: string; onClick?: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className="w-full flex items-center gap-4 p-4 bg-bg-elevated border border-border rounded-xl hover:border-accent/40 hover:bg-bg-hover transition-all group"
    >
      <div className={"w-10 h-10 rounded-lg " + logoColor + " flex items-center justify-center flex-shrink-0"}>
        <span className="text-white font-bold">{logoChar}</span>
      </div>
      <div className="flex-1 text-left">
        <p className="text-sm font-semibold text-foreground">Sign in with {provider}</p>
        <p className="text-xs text-text-muted">{orgName} · {orgDomain}</p>
      </div>
      <svg className="w-4 h-4 text-text-muted group-hover:text-foreground group-hover:translate-x-0.5 transition-all" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
    </button>
  );
}

export default SSOButton;
