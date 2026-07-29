"use client";

// command-item-01 — Command palette item (Linear Cmd+K-inspired)
export function CommandItem({
  label = "Create new issue",
  icon = "plus",
  shortcut,
  group,
  description,
  onSelect,
  isActive = false,
}: {
  label?: string; icon?: string; shortcut?: string[];
  group?: string; description?: string;
  onSelect?: () => void; isActive?: boolean;
}) {
  const iconPaths: Record<string, string> = {
    plus: "M12 4v16m8-8H4",
    search: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z",
    settings: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z",
    users: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z",
  };
  return (
    <button
      onClick={onSelect}
      className={"w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors " + (isActive ? "bg-accent/10 text-foreground" : "hover:bg-bg-hover text-text-secondary hover:text-foreground")}
    >
      <div className={"w-7 h-7 rounded-md flex items-center justify-center flex-shrink-0 " + (isActive ? "bg-accent/20 text-accent" : "bg-bg-hover text-text-muted")}>
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={iconPaths[icon] ?? iconPaths["plus"]} />
        </svg>
      </div>
      <div className="flex-1 min-w-0">
        <span className="text-sm font-medium block truncate">{label}</span>
        {description && <span className="text-xs text-text-muted truncate block">{description}</span>}
      </div>
      {shortcut && (
        <div className="flex items-center gap-1 flex-shrink-0">
          {shortcut.map((k) => (
            <kbd key={k} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{k}</kbd>
          ))}
        </div>
      )}
    </button>
  );
}

export default CommandItem;
