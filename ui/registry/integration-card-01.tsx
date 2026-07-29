"use client";

// integration-card-01 — Integration connect card (Linear / Notion-inspired)
export function IntegrationCard({
  name = "GitHub",
  description = "Sync issues, link commits, and track pull requests directly from your workspace.",
  category = "Developer Tools",
  isConnected = false,
  logoChar = "G",
  logoColor = "bg-gray-800",
  onConnect,
  onDisconnect,
  onSettings,
}: {
  name?: string; description?: string; category?: string;
  isConnected?: boolean; logoChar?: string; logoColor?: string;
  onConnect?: () => void; onDisconnect?: () => void; onSettings?: () => void;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-4 hover:border-accent/30 transition-colors">
      <div className="flex items-start gap-3">
        <div className={"w-10 h-10 rounded-lg " + logoColor + " flex items-center justify-center flex-shrink-0"}>
          <span className="text-white font-bold text-lg">{logoChar}</span>
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="text-foreground font-semibold text-sm">{name}</h3>
            {isConnected && (
              <span className="text-xs bg-success/10 text-success border border-success/20 px-1.5 py-0.5 rounded-full">Connected</span>
            )}
          </div>
          <span className="text-xs text-text-muted">{category}</span>
        </div>
      </div>
      <p className="text-sm text-text-secondary leading-relaxed">{description}</p>
      <div className="flex items-center gap-2">
        {isConnected ? (
          <>
            <button onClick={onSettings} className="flex-1 text-sm py-2 px-3 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors text-center">
              Settings
            </button>
            <button onClick={onDisconnect} className="text-sm py-2 px-3 rounded-lg border border-danger/30 text-danger hover:bg-danger/10 transition-colors">
              Disconnect
            </button>
          </>
        ) : (
          <button onClick={onConnect} className="flex-1 text-sm py-2 px-3 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors text-center">
            Connect {name}
          </button>
        )}
      </div>
    </div>
  );
}

export default IntegrationCard;
