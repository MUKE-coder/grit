import type { ColumnDefinition } from "@/lib/resource";
import { ChevronUp, ChevronDown } from "@/lib/icons";

interface ColumnHeaderProps {
  column: ColumnDefinition;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  onSort?: (key: string) => void;
}

export function ColumnHeader({ column, sortBy, sortOrder, onSort }: ColumnHeaderProps) {
  const isSorted = sortBy === column.key;

  if (!column.sortable || !onSort) {
    return (
      <th
        className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider"
        style={column.width ? { width: column.width } : undefined}
      >
        {column.label}
      </th>
    );
  }

  return (
    <th
      className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider cursor-pointer hover:text-foreground transition-colors select-none"
      style={column.width ? { width: column.width } : undefined}
      onClick={() => onSort(column.key)}
    >
      <div className="flex items-center gap-1">
        <span>{column.label}</span>
        <div className="flex flex-col">
          <ChevronUp
            className={`h-3 w-3 -mb-0.5 ${
              isSorted && sortOrder === "asc" ? "text-accent" : "text-text-muted/30"
            }`}
          />
          <ChevronDown
            className={`h-3 w-3 -mt-0.5 ${
              isSorted && sortOrder === "desc" ? "text-accent" : "text-text-muted/30"
            }`}
          />
        </div>
      </div>
    </th>
  );
}
