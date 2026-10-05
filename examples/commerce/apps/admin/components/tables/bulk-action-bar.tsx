"use client";

import type { CustomBulkAction } from "@/lib/resource";
import { getIcon, Archive, ArchiveRestore, Download, Pencil, Trash2, X } from "@/lib/icons";

/*
 * The bar that appears once rows are ticked.
 *
 * Fixed to the bottom of the viewport, centred. It sat in the flow at the foot
 * of the table first, on the reasoning that a floating bar covers the rows it
 * acts on. That reasoning only holds for a table that fits on screen: with
 * twenty rows you tick something near the top, the bar appears eight hundred
 * pixels below the fold, and as far as the operator can tell nothing happened.
 * A control that responds to a selection has to be where the selection is
 * being made.
 *
 * The original worry is answered rather than ignored. It is a centred pill
 * rather than a full-width bar, so the table is visible either side of it, and
 * the page reserves space underneath while it is shown, so the last rows can
 * still be scrolled clear of it.
 *
 * It is a labelled region so it turns up in a landmark list, and its arrival is
 * announced through the page's live region, because ticking a checkbox does not
 * move focus and a bar that silently appears is a bar a keyboard user never
 * learns about.
 *
 * Delete is the only red control. If everything is red then nothing is.
 */

export interface BulkActionBarProps {
  count: number;
  /** Built-ins the resource switched on, already filtered for the view. */
  actions: string[];
  custom: CustomBulkAction[];
  pending: boolean;
  singularName: string;
  pluralName: string;
  onEdit: () => void;
  onArchive: () => void;
  onRestore: () => void;
  onDelete: () => void;
  onExport: () => void;
  onCustom: (action: CustomBulkAction) => void;
  onClear: () => void;
}

export function BulkActionBar({
  count,
  actions,
  custom,
  pending,
  singularName,
  pluralName,
  onEdit,
  onArchive,
  onRestore,
  onDelete,
  onExport,
  onCustom,
  onClear,
}: BulkActionBarProps) {
  if (count === 0) return null;

  const noun = count === 1 ? singularName.toLowerCase() : pluralName.toLowerCase();
  const base =
    "inline-flex min-h-9 items-center gap-1.5 rounded-lg border px-3 text-sm font-medium transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent";
  const neutral = base + " border-border bg-bg-secondary text-foreground hover:bg-bg-hover";
  const danger = base + " border-danger/40 bg-transparent text-danger hover:bg-danger/10";

  return (
    <section
      aria-label="Bulk actions"
      className="fixed bottom-6 left-1/2 z-40 flex max-w-[calc(100vw-2rem)] -translate-x-1/2 flex-wrap items-center gap-3 rounded-xl border border-border bg-bg-elevated px-4 py-2.5 shadow-2xl shadow-black/25"
    >
      <p className="text-sm font-medium text-foreground">
        {count} {noun} selected
      </p>

      <div className="flex flex-wrap items-center gap-2">
        {actions.includes("edit") && (
          <button type="button" onClick={onEdit} disabled={pending} className={neutral}>
            <Pencil className="h-3.5 w-3.5" aria-hidden="true" />
            Edit
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {actions.includes("archive") && (
          <button type="button" onClick={onArchive} disabled={pending} className={neutral}>
            <Archive className="h-3.5 w-3.5" aria-hidden="true" />
            Archive
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {actions.includes("restore") && (
          <button type="button" onClick={onRestore} disabled={pending} className={neutral}>
            <ArchiveRestore className="h-3.5 w-3.5" aria-hidden="true" />
            Restore
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {actions.includes("export") && (
          <button type="button" onClick={onExport} disabled={pending} className={neutral}>
            <Download className="h-3.5 w-3.5" aria-hidden="true" />
            Export
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {custom.map((action) => {
          const Icon = action.icon ? getIcon(action.icon) : null;
          return (
            <button
              key={action.key}
              type="button"
              onClick={() => onCustom(action)}
              disabled={pending}
              className={action.variant === "danger" ? danger : neutral}
            >
              {Icon && <Icon className="h-3.5 w-3.5" aria-hidden="true" />}
              {action.label}
              <span className="sr-only"> for the {count} selected {noun}</span>
            </button>
          );
        })}

        {actions.includes("delete") && (
          <button type="button" onClick={onDelete} disabled={pending} className={danger}>
            <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
            Delete
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}
      </div>

      <button
        type="button"
        onClick={onClear}
        className="inline-flex min-h-9 items-center gap-1.5 rounded-lg px-3 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
      >
        <X className="h-3.5 w-3.5" aria-hidden="true" />
        Clear selection
      </button>
    </section>
  );
}
