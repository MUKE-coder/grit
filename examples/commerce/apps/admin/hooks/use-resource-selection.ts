"use client";

import { useCallback, useMemo, useState } from "react";

/**
 * The rows ticked in a resource list, by id, and the rows behind them.
 *
 * selectedRows reads the page already loaded, so a custom bulk action like
 * "email the people I ticked" does not need a second round trip.
 */
export function useResourceSelection<T>(rows: T[]) {
  const [selection, setSelection] = useState<string[]>([]);

  const clearSelection = useCallback(() => setSelection([]), []);

  const selectedRows = useMemo(
    () => rows.filter((row) => selection.includes(String((row as Record<string, unknown>).id))),
    [rows, selection],
  );

  return { selection, setSelection, clearSelection, selectedRows };
}
