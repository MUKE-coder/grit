"use client";

// floating-action-01 — Floating action button (Notion "New page"-inspired)
export function FloatingActionButton({
  label = "New",
  icon = "plus",
  onClick,
  position = "bottom-right",
}: {
  label?: string; icon?: "plus" | "edit" | "message";
  onClick?: () => void; position?: "bottom-right" | "bottom-left";
}) {
  const iconPaths: Record<string, string> = {
    plus: "M12 4v16m8-8H4",
    edit: "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
    message: "M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z",
  };
  const posClass = position === "bottom-right" ? "bottom-6 right-6" : "bottom-6 left-6";
  return (
    <button
      onClick={onClick}
      className={"fixed " + posClass + " flex items-center gap-2 px-4 py-3 bg-accent hover:bg-accent/90 text-white rounded-full shadow-lg hover:shadow-accent/30 transition-all font-semibold text-sm hover:scale-105 active:scale-95 z-40"}
    >
      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={iconPaths[icon] ?? iconPaths["plus"]} />
      </svg>
      {label}
    </button>
  );
}

export default FloatingActionButton;
