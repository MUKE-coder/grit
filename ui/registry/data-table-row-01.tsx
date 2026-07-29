"use client";

// data-table-row-01 — Data table row with actions (GitHub / Linear-inspired)
export function DataTableRow({
  cells = ["sarah@example.com", "Admin", "Mar 1, 2025", "Active"],
  statusIndex = 3,
  statusMap = { Active: "text-success", Pending: "text-yellow-400", Inactive: "text-text-muted", Suspended: "text-danger" },
  onEdit,
  onDelete,
  selected = false,
  onSelect,
}: {
  cells?: string[]; statusIndex?: number;
  statusMap?: Record<string, string>;
  onEdit?: () => void; onDelete?: () => void;
  selected?: boolean; onSelect?: () => void;
}) {
  return (
    <tr className={"border-b border-border transition-colors hover:bg-bg-hover " + (selected ? "bg-accent/5" : "")}>
      <td className="py-3 pl-4 pr-2">
        <input type="checkbox" checked={selected} onChange={onSelect} className="rounded border-border accent-accent" />
      </td>
      {cells.map((cell, i) => (
        <td key={i} className="py-3 px-3">
          {i === statusIndex ? (
            <span className={"text-sm font-medium " + (statusMap[cell] ?? "text-text-secondary")}>{cell}</span>
          ) : (
            <span className="text-sm text-text-secondary">{cell}</span>
          )}
        </td>
      ))}
      <td className="py-3 pr-4 pl-2">
        <div className="flex items-center gap-1 justify-end opacity-0 group-hover:opacity-100">
          <button onClick={onEdit} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
          </button>
          <button onClick={onDelete} className="p-1.5 rounded text-text-muted hover:text-danger transition-colors">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
          </button>
        </div>
      </td>
    </tr>
  );
}

export default DataTableRow;
