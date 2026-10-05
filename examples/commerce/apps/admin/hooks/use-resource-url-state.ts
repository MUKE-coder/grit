"use client";

import { useCallback, useMemo, useState } from "react";
import {
  usePathname,
  useRouter,
  useSearchParams,
  type ReadonlyURLSearchParams,
} from "next/navigation";
import type { ColumnDefinition, ResourceDefinition, TableTab } from "@/lib/resource";
import { useDebouncedValue } from "@/hooks/use-resource";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";

// Read the date filter back out of the address bar so a refresh or a shared
// link rehydrates the same view.
function readDateRangeFromURL(sp: ReadonlyURLSearchParams | null): DateRange {
  if (!sp) return {};
  const preset = sp.get("date") as DateRange["preset"] | null;
  if (preset === "custom") {
    return {
      preset: "custom",
      from: sp.get("date_from") ?? undefined,
      to: sp.get("date_to") ?? undefined,
    };
  }
  if (preset === "today" || preset === "7d" || preset === "30d" || preset === "month") {
    return { preset };
  }
  return {};
}

// replace, not push: the back button should not collect one entry per filter
// tweak.
function writeDateRangeToURL(
  router: ReturnType<typeof useRouter>,
  pathname: string,
  current: ReadonlyURLSearchParams | null,
  range: DateRange,
) {
  const params = new URLSearchParams(current?.toString() ?? "");
  params.delete("date");
  params.delete("date_from");
  params.delete("date_to");
  if (range.preset) {
    params.set("date", range.preset);
    if (range.preset === "custom") {
      if (range.from) params.set("date_from", range.from);
      if (range.to) params.set("date_to", range.to);
    }
  }
  const qs = params.toString();
  router.replace(qs ? pathname + "?" + qs : pathname, { scroll: false });
}

export interface ResourceURLStateOptions {
  /** Start on a page other than 1. */
  initialPage?: number;
  /** Override the resource's configured page size. */
  initialPageSize?: number;
}

/**
 * What a resource list is asking the API for: page, size, search, sort,
 * filters, date range, tab and archived view, plus which columns are shown.
 *
 * Every setter that changes what is queried goes back to page 1, and the
 * date range round-trips through the address bar.
 */
export function useResourceURLState(resource: ResourceDefinition, options: ResourceURLStateOptions = {}) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const [page, setPage] = useState(options.initialPage ?? 1);
  const [pageSize, setPageSizeState] = useState(
    options.initialPageSize ?? resource.table.pageSize ?? 20,
  );
  const [search, setSearchState] = useState("");
  // The box shows every keystroke; the list asks once typing pauses.
  const debouncedSearch = useDebouncedValue(search, 300);
  const [sortBy, setSortBy] = useState(resource.table.defaultSort?.key ?? "");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">(
    resource.table.defaultSort?.direction ?? "desc",
  );
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [hiddenColumns, setHiddenColumns] = useState<string[]>([]);

  const [dateRange, setDateRangeState] = useState<DateRange>(() =>
    readDateRangeFromURL(searchParams),
  );
  const dateParams = useMemo(() => dateRangeToQueryParams(dateRange), [dateRange]);
  const setDateRange = useCallback(
    (next: DateRange) => {
      setDateRangeState(next);
      writeDateRangeToURL(router, pathname, searchParams, next);
      setPage(1);
    },
    [router, pathname, searchParams],
  );

  const tabs: TableTab[] = useMemo(() => resource.table.tabs ?? [], [resource.table.tabs]);
  // First tab on load. A tab strip where nothing is selected reads as broken,
  // and the first tab is conventionally the unfiltered one.
  const [activeTab, setActiveTabState] = useState(() => tabs[0]?.key ?? "");
  const [showArchived, setShowArchivedState] = useState(false);

  const setActiveTab = useCallback((key: string) => {
    setActiveTabState(key);
    setPage(1);
  }, []);

  const setShowArchived = useCallback((value: boolean) => {
    setShowArchivedState(value);
    setPage(1);
  }, []);

  // Mirrors the query useResource builds, so an export applies the same
  // filter and sort the operator is looking at.
  const apiSearchParams = useMemo(() => {
    const sp = new URLSearchParams();
    if (debouncedSearch) sp.set("search", debouncedSearch);
    if (sortBy) {
      sp.set("sort_by", sortBy);
      sp.set("sort_order", sortOrder);
    }
    Object.entries(filters).forEach(([k, v]) => {
      if (v) sp.set(k, v);
    });
    Object.entries(dateParams).forEach(([k, v]) => {
      if (v) sp.set(k, v);
    });
    const df = resource.table.dateFilter?.field;
    if (df && df !== "created_at") sp.set("date_field", df);
    return sp;
  }, [debouncedSearch, sortBy, sortOrder, filters, dateParams, resource.table.dateFilter?.field]);

  const columns: ColumnDefinition[] = useMemo(
    () => resource.table.columns.filter((col) => !col.hidden && !hiddenColumns.includes(col.key)),
    [resource.table.columns, hiddenColumns],
  );

  const toggleColumn = useCallback((key: string) => {
    setHiddenColumns((prev) =>
      prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key],
    );
  }, []);

  // Any change to what is being queried resets to page 1, otherwise a
  // search from page 7 lands on an empty page 7 of two results.
  const setSearch = useCallback((value: string) => {
    setSearchState(value);
    setPage(1);
  }, []);

  const setPageSize = useCallback((size: number) => {
    setPageSizeState(size);
    setPage(1);
  }, []);

  const setSort = useCallback(
    (key: string) => {
      if (sortBy === key) {
        setSortOrder((prev) => (prev === "asc" ? "desc" : "asc"));
      } else {
        setSortBy(key);
        setSortOrder("asc");
      }
      setPage(1);
    },
    [sortBy],
  );

  const setFilter = useCallback((key: string, value: string) => {
    setFilters((prev) => {
      if (!value) {
        const next = { ...prev };
        delete next[key];
        return next;
      }
      return { ...prev, [key]: value };
    });
    setPage(1);
  }, []);

  return {
    page,
    pageSize,
    search,
    debouncedSearch,
    sortBy,
    sortOrder,
    filters,
    dateRange,
    dateParams,
    setPage,
    setPageSize,
    setSearch,
    setSort,
    setFilter,
    setDateRange,
    columns,
    hiddenColumns,
    toggleColumn,
    tabs,
    activeTab,
    setActiveTab,
    showArchived,
    setShowArchived,
    apiSearchParams,
  };
}
