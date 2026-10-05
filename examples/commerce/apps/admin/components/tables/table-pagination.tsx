import { useT } from "@/lib/i18n";

interface TablePaginationProps {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (size: number) => void;
}

export function TablePagination({
  page,
  pageSize,
  total,
  totalPages,
  onPageChange,
  onPageSizeChange,
}: TablePaginationProps) {
  const t = useT();
  if (total === 0) return null;

  const start = (page - 1) * pageSize + 1;
  const end = Math.min(page * pageSize, total);

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 border-t border-border p-4">
      <div className="flex items-center gap-3">
        <p className="text-sm text-text-muted">
          {t("table.showing", "Showing {start}–{end} of {total}", { start, end, total })}
        </p>
        {/* Toolbar chrome, not a form field — the input slot sets w-full,
            which would stretch this to fill the row. */}
        <select
          value={pageSize}
          onChange={(e) => onPageSizeChange(Number(e.target.value))}
          className="rounded-lg border border-border bg-bg-tertiary px-2 py-1 text-sm text-foreground focus:border-accent focus:outline-none"
        >
          {[10, 20, 50, 100].map((size) => (
            <option key={size} value={size}>
              {t("table.perPage", "{size} / page", { size })}
            </option>
          ))}
        </select>
      </div>

      <div className="flex items-center gap-1">
        <button
          onClick={() => onPageChange(1)}
          disabled={page <= 1}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {t("table.first", "First")}
        </button>
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {t("table.prev", "Prev")}
        </button>

        {generatePageNumbers(page, totalPages).map((p, i) =>
          p === -1 ? (
            <span key={`ellipsis-${i}`} className="px-1 text-text-muted">...</span>
          ) : (
            <button
              key={p}
              onClick={() => onPageChange(p)}
              className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                p === page
                  ? "bg-accent text-accent-fg"
                  : "border border-border bg-bg-tertiary text-text-secondary hover:bg-bg-hover"
              }`}
            >
              {p}
            </button>
          )
        )}

        <button
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {t("table.next", "Next")}
        </button>
        <button
          onClick={() => onPageChange(totalPages)}
          disabled={page >= totalPages}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {t("table.last", "Last")}
        </button>
      </div>
    </div>
  );
}

function generatePageNumbers(current: number, total: number): number[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);

  if (current <= 3) return [1, 2, 3, 4, -1, total];
  if (current >= total - 2) return [1, -1, total - 3, total - 2, total - 1, total];

  return [1, -1, current - 1, current, current + 1, -1, total];
}
