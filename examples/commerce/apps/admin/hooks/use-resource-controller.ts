"use client";

import { useCallback, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import type {
  BulkAction,
  ColumnDefinition,
  CustomBulkAction,
  ResourceDefinition,
  TableAction,
  TableTab,
} from "@/lib/resource";
import { useBulkResource, useDeleteResource, useResource } from "@/hooks/use-resource";
import { useResourceURLState } from "@/hooks/use-resource-url-state";
import { useResourceSelection } from "@/hooks/use-resource-selection";
import { useResourceDialogs } from "@/hooks/use-resource-dialogs";
import type { StatCard } from "@/components/chrome/StatCards";
import type { DateRange } from "@/components/tables/date-filter";

// The windows the default stat cards show beside the total. The list request
// asks for them (?counts=), so the cards cost no requests of their own.
const DEFAULT_STAT_COUNTS = ["created_7d", "created_30d", "updated_7d"];

export interface ResourceControllerOptions {
  /** Start on a page other than 1. */
  initialPage?: number;
  /** Override the resource's configured page size. */
  initialPageSize?: number;
}

export interface ResourceController<T = Record<string, unknown>> {
  resource: ResourceDefinition;

  // ── data ────────────────────────────────────────────────────────────
  rows: T[];
  meta:
    | { total: number; page: number; page_size: number; pages: number; counts?: Record<string, number> | null }
    | undefined;
  total: number;
  totalPages: number;
  isLoading: boolean;
  /** A refetch with rows already on screen: a new page, sort or search. */
  isFetching: boolean;

  // ── query state (all of it URL- or server-aware) ────────────────────
  page: number;
  pageSize: number;
  search: string;
  sortBy: string;
  sortOrder: "asc" | "desc";
  filters: Record<string, string>;
  dateRange: DateRange;
  setPage: (page: number) => void;
  setPageSize: (size: number) => void;
  setSearch: (value: string) => void;
  /** Toggles direction when the same key is passed twice. */
  setSort: (key: string) => void;
  setFilter: (key: string, value: string) => void;
  setDateRange: (range: DateRange) => void;

  // ── columns ─────────────────────────────────────────────────────────
  /** resource.table.columns minus hidden ones — what a table should render. */
  columns: ColumnDefinition[];
  allColumns: ColumnDefinition[];
  hiddenColumns: string[];
  toggleColumn: (key: string) => void;

  // ── selection ───────────────────────────────────────────────────────
  selection: string[];
  setSelection: (ids: string[]) => void;
  clearSelection: () => void;

  // ── actions ─────────────────────────────────────────────────────────
  actions: TableAction[];
  can: (action: TableAction) => boolean;
  create: () => void;
  /** Create with fields pre-filled, e.g. createWith({ parent_id: id }). */
  createWith: (defaults: Record<string, unknown>) => void;
  edit: (row: T) => void;
  view: (row: T) => void;
  /** Opens the confirm dialog; deletion happens on confirm. */
  remove: (id: string) => void;
  bulkRemove: () => void;
  isDeleting: boolean;
  isBulkDeleting: boolean;

  // ── bulk actions ────────────────────────────────────────────────────
  /** Built-ins the resource has switched on, minus any that make no sense
   *  in the current view (restore only appears while Archived is open). */
  bulkActions: BulkAction[];
  /** Custom ones from resources/<name>.custom.tsx, already filtered by visible(). */
  customBulkActions: CustomBulkAction<T>[];
  /** The rows behind the current selection, readable without a refetch. */
  selectedRows: T[];
  /** Opens the confirm dialog; archiving happens on confirm. */
  bulkArchive: () => void;
  /** Restores immediately: putting something back is not destructive. */
  bulkRestore: () => void;
  /** Opens the bulk edit dialog. */
  bulkEdit: () => void;
  /** Writes one field to every selected row and closes the dialog. */
  applyBulkEdit: (patch: Record<string, unknown>) => void;
  /** Runs a custom action, handing it the ids, the rows and the helpers. */
  runBulkAction: (action: CustomBulkAction<T>) => void;
  isBulkPending: boolean;
  /** Re-runs the list query. Handed to custom actions so they can refresh. */
  refresh: () => void;
  /** Speaks to the page's live region. Bulk changes never move focus. */
  announce: (message: string) => void;
  liveMessage: string;

  // ── tabs ────────────────────────────────────────────────────────────
  /** The resource's filter presets, or an empty array when it has none. */
  tabs: TableTab[];
  /** Key of the tab currently applied. "" when the resource has no tabs. */
  activeTab: string;
  setActiveTab: (key: string) => void;

  // ── archived view ───────────────────────────────────────────────────
  /** True while the Archived tab is open. */
  showArchived: boolean;
  setShowArchived: (value: boolean) => void;

  // ── dialog state, for anyone rendering their own ────────────────────
  form: {
    open: boolean;
    item: T | null;
    /** Pre-filled values for create mode, set by createWith(). */
    defaults?: Record<string, unknown>;
    close: () => void;
  };
  confirmDelete: { open: boolean; confirm: () => void; cancel: () => void };
  confirmBulkDelete: { open: boolean; confirm: () => void; cancel: () => void };
  confirmBulkArchive: { open: boolean; confirm: () => void; cancel: () => void };
  bulkEditor: { open: boolean; close: () => void };
  /** Set when a custom action asked to confirm first. */
  confirmCustom: {
    open: boolean;
    action: CustomBulkAction<T> | null;
    confirm: () => void;
    cancel: () => void;
  };
  importer: { open: boolean; setOpen: (open: boolean) => void };

  // ── odds and ends the default page needs ────────────────────────────
  /** Same query the table ran, for an export that matches what is on screen. */
  apiSearchParams: URLSearchParams;
  stats: StatCard[] | undefined;
  singularName: string;
  pluralName: string;
  isFormPage: boolean;
  isSteps: boolean;
}

/**
 * Everything a resource list page needs except the markup.
 *
 * const c = useResourceController(productsResource)
 * <MyTable rows={c.rows} onSort={c.setSort} onRowClick={c.edit} />
 *
 * Built from three hooks you can also use on their own: useResourceURLState
 * (what is being asked for), useResourceSelection (what is ticked) and
 * useResourceDialogs (what is open). This function joins them to the API.
 */
export function useResourceController<T = Record<string, unknown>>(
  resource: ResourceDefinition,
  options: ResourceControllerOptions = {},
): ResourceController<T> {
  const router = useRouter();
  const queryClient = useQueryClient();

  const isFormPage = resource.formView === "page" || resource.formView === "page-steps";
  const isSteps = resource.formView === "modal-steps" || resource.formView === "page-steps";

  const url = useResourceURLState(resource, options);
  const dialogs = useResourceDialogs<T>();
  const [liveMessage, setLiveMessage] = useState("");

  const statsConfig = resource.stats;
  const statsEnabled =
    statsConfig === undefined ||
    statsConfig === true ||
    (typeof statsConfig === "object" && statsConfig !== null && statsConfig.enabled !== false);
  const customStatCards =
    typeof statsConfig === "object" &&
    statsConfig !== null &&
    Array.isArray(statsConfig.cards) &&
    statsConfig.cards.length > 0;

  const { data, isLoading, isFetching } = useResource<T>(resource.endpoint, {
    page: url.page,
    pageSize: url.pageSize,
    search: url.debouncedSearch,
    sortBy: url.sortBy,
    sortOrder: url.sortOrder,
    // Tab filters, then the operator's own, then the archived flag. The
    // operator's win: picking "Unpaid" and then filtering by customer should
    // narrow the tab, not silently leave it.
    filters: {
      ...(url.tabs.find((t) => t.key === url.activeTab)?.filters ?? {}),
      ...url.filters,
      ...(url.showArchived ? { archived: "true" } : {}),
    },
    dateParams: url.dateParams,
    dateField: resource.table.dateFilter?.field,
    counts: statsEnabled && !customStatCards ? DEFAULT_STAT_COUNTS : undefined,
  });

  const rows = useMemo(() => data?.data ?? [], [data]);
  const { selection, setSelection, clearSelection, selectedRows } = useResourceSelection<T>(rows);

  // Switching tabs or the archived view changes which rows exist, so a
  // selection made in the other one is stale. Keeping it is how you archive
  // something you cannot see.
  const { setActiveTab: setURLActiveTab, setShowArchived: setURLShowArchived } = url;
  const setActiveTab = useCallback(
    (key: string) => {
      setURLActiveTab(key);
      clearSelection();
    },
    [setURLActiveTab, clearSelection],
  );
  const setShowArchived = useCallback(
    (value: boolean) => {
      setURLShowArchived(value);
      clearSelection();
    },
    [setURLShowArchived, clearSelection],
  );

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(
    resource.endpoint,
    singularName,
  );
  const { mutate: runBulk, isPending: isBulkPending } = useBulkResource(
    resource.endpoint,
    pluralName,
    singularName,
  );
  // Kept as its own name because the delete confirm dialog shows a spinner
  // for delete specifically, not for any bulk action in flight.
  const isBulkDeleting = isBulkPending;

  const { openForm, setDeletingId, setBulkDeleteOpen, setBulkArchiveOpen, setBulkEditOpen, setPendingCustom } = dialogs;

  const view = useCallback(
    (row: T) => {
      const id = String((row as Record<string, unknown>).id);
      router.push("/resources/" + resource.slug + "/" + id);
    },
    [router, resource.slug],
  );

  const edit = useCallback(
    (row: T) => {
      if (isFormPage) {
        const id = String((row as Record<string, unknown>).id);
        router.push("/resources/" + resource.slug + "?action=edit&edit=" + id);
      } else {
        openForm(row);
      }
    },
    [isFormPage, router, resource.slug, openForm],
  );

  const create = useCallback(() => {
    if (isFormPage) {
      router.push("/resources/" + resource.slug + "?action=create");
    } else {
      openForm(null);
    }
  }, [isFormPage, router, resource.slug, openForm]);

  /**
   * Create, with some fields already filled in.
   *
   * createWith({ parent_id: id }) is how the tree view adds a child to the row
   * you clicked. In page mode the values ride along as query params, because a
   * route change is the only state that survives the navigation.
   */
  const createWith = useCallback(
    (defaults: Record<string, unknown>) => {
      if (isFormPage) {
        const params = new URLSearchParams({ action: "create" });
        for (const [key, value] of Object.entries(defaults)) {
          if (value !== undefined && value !== null) params.set(key, String(value));
        }
        router.push("/resources/" + resource.slug + "?" + params.toString());
        return;
      }
      openForm(null, defaults);
    },
    [isFormPage, router, resource.slug, openForm],
  );

  const remove = useCallback((id: string) => setDeletingId(id), [setDeletingId]);

  const { deletingId } = dialogs;
  const doDelete = useCallback(() => {
    if (deletingId !== null) {
      deleteItem(deletingId, { onSuccess: () => setDeletingId(null) });
    }
  }, [deleteItem, deletingId, setDeletingId]);

  const bulkRemove = useCallback(() => {
    if (selection.length > 0) setBulkDeleteOpen(true);
  }, [selection, setBulkDeleteOpen]);

  const doBulkDelete = useCallback(() => {
    runBulk(
      { action: "delete", ids: selection },
      {
        onSuccess: () => {
          setBulkDeleteOpen(false);
          setSelection([]);
        },
      },
    );
  }, [runBulk, selection, setBulkDeleteOpen, setSelection]);

  // ── the rest of the bulk surface ──────────────────────────────────────

  const announce = useCallback((message: string) => {
    // Cleared first: setting the same string twice is not a change, and a
    // live region that has not changed says nothing. Two identical bulk
    // actions in a row would be announced once.
    setLiveMessage("");
    requestAnimationFrame(() => setLiveMessage(message));
  }, []);

  const refresh = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: [resource.endpoint] });
  }, [queryClient, resource.endpoint]);

  const bulkArchive = useCallback(() => {
    if (selection.length > 0) setBulkArchiveOpen(true);
  }, [selection, setBulkArchiveOpen]);

  const doBulkArchive = useCallback(() => {
    runBulk(
      { action: "archive", ids: selection },
      {
        onSuccess: () => {
          setBulkArchiveOpen(false);
          setSelection([]);
          announce(selection.length + " archived.");
        },
      },
    );
  }, [runBulk, selection, setBulkArchiveOpen, setSelection, announce]);

  // No confirm: putting something back is not destructive, and a dialog in
  // front of an undo is a dialog nobody reads.
  const bulkRestore = useCallback(() => {
    if (selection.length === 0) return;
    runBulk(
      { action: "restore", ids: selection },
      {
        onSuccess: () => {
          setSelection([]);
          announce(selection.length + " restored.");
        },
      },
    );
  }, [runBulk, selection, setSelection, announce]);

  const bulkEdit = useCallback(() => {
    if (selection.length > 0) setBulkEditOpen(true);
  }, [selection, setBulkEditOpen]);

  const applyBulkEdit = useCallback(
    (patch: Record<string, unknown>) => {
      runBulk(
        { action: "patch", ids: selection, patch },
        {
          onSuccess: () => {
            setBulkEditOpen(false);
            setSelection([]);
            announce(selection.length + " updated.");
          },
        },
      );
    },
    [runBulk, selection, setBulkEditOpen, setSelection, announce],
  );

  const runCustom = useCallback(
    (action: CustomBulkAction<T>) => {
      void action.onSelect(selection, selectedRows, {
        refresh,
        clearSelection,
        announce,
      });
    },
    [selection, selectedRows, refresh, clearSelection, announce],
  );

  const runBulkAction = useCallback(
    (action: CustomBulkAction<T>) => {
      if (action.confirm) {
        setPendingCustom(action);
        return;
      }
      runCustom(action);
    },
    [runCustom, setPendingCustom],
  );

  // Restore only makes sense on rows that are archived, and archive only on
  // rows that are not, so the two never appear together. Offering both is how
  // an operator ends up archiving what they meant to bring back.
  const { showArchived, dateParams } = url;
  const bulkActions = useMemo(() => {
    // ["edit", "export", "delete"] rather than ["delete"] alone: a resource
    // that predates bulkActions still gets the three that work against any
    // API. Archive and restore are opt-in because they need both the column
    // and the endpoint.
    const configured = resource.table.bulkActions ?? ["edit", "export", "delete"];
    return configured.filter((action) => {
      if (action === "restore") return showArchived;
      if (action === "archive") return !showArchived;
      return true;
    });
  }, [resource.table.bulkActions, showArchived]);

  const customBulkActions = useMemo(() => {
    const all = (resource.customBulkActions ?? []) as CustomBulkAction<T>[];
    return all.filter((action) => !action.visible || action.visible(selectedRows));
  }, [resource.customBulkActions, selectedRows]);

  const actions = resource.table.actions ?? ["create", "view", "edit", "delete"];
  const can = useCallback((action: TableAction) => actions.includes(action), [actions]);

  const stats: StatCard[] | undefined = useMemo(() => {
    if (!statsEnabled) return undefined;

    // Every stat endpoint gets whatever narrows the table, or "Total: 10,000"
    // sits above a table showing 142 matches. The archived view counts as
    // narrowing: without it the Archived tab reads "Total 9" over two rows.
    const applyViewParams = (cards: StatCard[]): StatCard[] => {
      const extra: Record<string, string> = { ...dateParams };
      if (showArchived) extra.archived = "true";
      if (Object.keys(extra).length === 0) return cards;
      return cards.map((card) => {
        if (!card.endpoint) return card;
        const sep = card.endpoint.includes("?") ? "&" : "?";
        const qs = new URLSearchParams(extra).toString();
        return { ...card, endpoint: card.endpoint + sep + qs };
      });
    };

    if (customStatCards && typeof statsConfig === "object" && statsConfig !== null && statsConfig.cards) {
      return applyViewParams(statsConfig.cards);
    }

    // The defaults read the list response: the total the table shows, and the
    // windows it asked for with ?counts=. So they describe the rows the table
    // matches, search and filters included, and cost no requests of their own.
    // A dash means the list endpoint does not answer ?counts=.
    const counts = data?.meta?.counts;
    const count = (name: string) => counts?.[name] ?? "—";
    const loading = isLoading;
    const defaults: StatCard[] = [
      { label: "Total", value: data?.meta?.total ?? "—", loading, icon: resource.icon || "Package" },
      { label: "This Week", value: count("created_7d"), loading, icon: "TrendingUp", color: "success" },
      { label: "This Month", value: count("created_30d"), loading, icon: "Calendar", color: "info" },
      { label: "Updated Recently", value: count("updated_7d"), loading, icon: "RefreshCw" },
    ];
    return defaults;
  }, [statsEnabled, customStatCards, statsConfig, resource.icon, dateParams, showArchived, data, isLoading]);

  return {
    resource,

    rows,
    meta: data?.meta,
    total: data?.meta?.total ?? 0,
    totalPages: data?.meta?.pages ?? 1,
    isLoading,
    isFetching: isFetching && !isLoading,

    page: url.page,
    pageSize: url.pageSize,
    search: url.search,
    sortBy: url.sortBy,
    sortOrder: url.sortOrder,
    filters: url.filters,
    dateRange: url.dateRange,
    setPage: url.setPage,
    setPageSize: url.setPageSize,
    setSearch: url.setSearch,
    setSort: url.setSort,
    setFilter: url.setFilter,
    setDateRange: url.setDateRange,

    columns: url.columns,
    allColumns: resource.table.columns,
    hiddenColumns: url.hiddenColumns,
    toggleColumn: url.toggleColumn,

    selection,
    setSelection,
    clearSelection,

    actions,
    can,
    create,
    createWith,
    edit,
    view,
    remove,
    bulkRemove,
    isDeleting,
    isBulkDeleting,

    bulkActions,
    customBulkActions,
    selectedRows,
    bulkArchive,
    bulkRestore,
    bulkEdit,
    applyBulkEdit,
    runBulkAction,
    isBulkPending,
    refresh,
    announce,
    liveMessage,

    tabs: url.tabs,
    activeTab: url.activeTab,
    setActiveTab,

    showArchived,
    setShowArchived,

    form: {
      open: dialogs.formOpen,
      item: dialogs.editingItem,
      defaults: dialogs.formDefaults,
      close: dialogs.closeForm,
    },
    confirmDelete: {
      open: deletingId !== null,
      confirm: doDelete,
      cancel: () => setDeletingId(null),
    },
    confirmBulkDelete: {
      open: dialogs.bulkDeleteOpen,
      confirm: doBulkDelete,
      cancel: () => setBulkDeleteOpen(false),
    },
    confirmBulkArchive: {
      open: dialogs.bulkArchiveOpen,
      confirm: doBulkArchive,
      cancel: () => setBulkArchiveOpen(false),
    },
    bulkEditor: {
      open: dialogs.bulkEditOpen,
      close: () => setBulkEditOpen(false),
    },
    confirmCustom: {
      open: dialogs.pendingCustom !== null,
      action: dialogs.pendingCustom,
      confirm: () => {
        if (dialogs.pendingCustom) runCustom(dialogs.pendingCustom);
        setPendingCustom(null);
      },
      cancel: () => setPendingCustom(null),
    },
    importer: { open: dialogs.importOpen, setOpen: dialogs.setImportOpen },

    apiSearchParams: url.apiSearchParams,
    stats,
    singularName,
    pluralName,
    isFormPage,
    isSteps,
  };
}
