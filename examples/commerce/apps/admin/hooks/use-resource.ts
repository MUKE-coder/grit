import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
  type QueryClient,
} from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-core";

// Every query about a resource starts with its endpoint, so invalidating
// [endpoint] still refreshes all of it. Each kind of read has a branch below
// that, so a save can refresh what it moved. Build keys here and nowhere else:
// two components asking for the same thing under two keys make two requests,
// and a key that does not start with the endpoint is never refreshed by a save.
export const resourceKeys = {
  all: (endpoint: string) => [endpoint] as const,
  lists: (endpoint: string) => [endpoint, "list"] as const,
  list: (endpoint: string, params: object) => [endpoint, "list", params] as const,
  detail: (endpoint: string, id: string) => [endpoint, "detail", id] as const,
  // PageHeader's stat cards key on [endpoint, "stat", ...].
  stats: (endpoint: string) => [endpoint, "stat"] as const,
  // The dashboard's stat card and latest table, which share one response.
  dashboardStats: (endpoint: string, params: object) => [endpoint, "dashboard-stats", params] as const,
  // Relationship pickers, single and multi, by what was typed.
  options: (endpoint: string, search = "") => [endpoint, "options", search] as const,
  // A table tab's count badge, over the filters the tab applies.
  tabCount: (endpoint: string, tab: string, filters: object) => [endpoint, "tab-count", tab, filters] as const,
  tree: (endpoint: string) => [endpoint, "tree"] as const,
};

// useDebouncedValue follows value once it has stopped changing for delay
// milliseconds. The search box uses it so typing "john" is one request.
export function useDebouncedValue<T>(value: T, delay: number): T {
  const [settled, setSettled] = useState(value);
  useEffect(() => {
    const timer = setTimeout(() => setSettled(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);
  return settled;
}

// refreshAfterSave refetches what a change to an existing row can move: the
// lists, the record and any other view of the resource. Stat cards are marked
// stale for their next mount instead: an edit rarely changes a count, and
// refetching them cost a request per card on every save. The default cards
// read their counts from the list, which does refetch.
function refreshAfterSave(queryClient: QueryClient, endpoint: string) {
  queryClient.invalidateQueries({
    queryKey: resourceKeys.all(endpoint),
    predicate: (query) => query.queryKey[1] !== "stat",
  });
  queryClient.invalidateQueries({ queryKey: resourceKeys.stats(endpoint), refetchType: "none" });
}

interface ResourceQueryParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  filters?: Record<string, string>;
  // v3.31.34 — date-window filter. dateParams comes from
  // dateRangeToQueryParams(); dateField overrides the server's
  // default "created_at" target column when set.
  dateParams?: Record<string, string>;
  dateField?: string;
  // Extra totals to return beside the page, such as created_7d: the default
  // stat cards read theirs from the list response instead of a request each.
  counts?: string[];
}

interface PaginatedResponse<T = Record<string, unknown>> {
  data: T[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
    // Answers ?counts=. Absent when the API does not support it.
    counts?: Record<string, number> | null;
  };
}

export function useResource<T = Record<string, unknown>>(
  endpoint: string,
  params: ResourceQueryParams = {}
) {
  const { page = 1, pageSize = 20, search, sortBy, sortOrder, filters, dateParams, dateField, counts } = params;

  return useQuery<PaginatedResponse<T>>({
    // v3.31.34: dateParams + dateField included in key so a date
    // filter change invalidates the cache and the list refetches.
    queryKey: resourceKeys.list(endpoint, { page, pageSize, search, sortBy, sortOrder, filters, dateParams, dateField, counts }),
    // The signal cancels a request the next keystroke or page has replaced, so
    // a slow answer to an old search cannot land on top of a newer one.
    queryFn: async ({ signal }) => {
      const searchParams = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
      });

      if (search) searchParams.set("search", search);
      if (sortBy) {
        searchParams.set("sort_by", sortBy);
        searchParams.set("sort_order", sortOrder ?? "desc");
      }
      if (filters) {
        Object.entries(filters).forEach(([key, value]) => {
          if (value) searchParams.set(key, value);
        });
      }
      if (dateParams) {
        Object.entries(dateParams).forEach(([key, value]) => {
          if (value) searchParams.set(key, value);
        });
      }
      if (dateField && dateField !== "created_at") {
        searchParams.set("date_field", dateField);
      }
      if (counts && counts.length > 0) {
        searchParams.set("counts", counts.join(","));
      }

      const { data } = await apiClient.get(`${endpoint}?${searchParams}`, { signal });
      return data;
    },
    // Keep the rows on screen while the next page, sort or search loads. The
    // table used to blank to a skeleton on every change.
    placeholderData: keepPreviousData,
  });
}

export function useResourceItem<T = Record<string, unknown>>(
  endpoint: string,
  id: string,
  options?: { enabled?: boolean }
) {
  return useQuery<{ data: T }>({
    queryKey: resourceKeys.detail(endpoint, id),
    queryFn: async () => {
      const { data } = await apiClient.get(`${endpoint}/${id}`);
      return data;
    },
    enabled: (options?.enabled ?? true) && !!id,
  });
}

export interface ResourceDashboardStats {
  resource: string;
  total: number;
  // Always 30 daily buckets, whatever the date range.
  series: { date: string; count: number }[];
  latest: Record<string, unknown>[];
}

// The most rows any dashboard widget shows. One number, so the stat card and
// the latest table ask for the same URL and share its answer.
export const DASHBOARD_LATEST_LIMIT = 10;

// useResourceDashboardStats reads a resource's dashboard stats. The stat card
// and the latest table used to fetch them separately under keys of their own,
// which was two requests per resource on every load and two more every minute.
// Both call this now and take their part with select.
export function useResourceDashboardStats<TSelected = ResourceDashboardStats>(
  endpoint: string,
  params: Record<string, string>,
  select?: (stats: ResourceDashboardStats) => TSelected,
) {
  return useQuery<ResourceDashboardStats, Error, TSelected>({
    queryKey: resourceKeys.dashboardStats(endpoint, params),
    queryFn: async ({ signal }) => {
      const search = new URLSearchParams({ ...params, limit: String(DASHBOARD_LATEST_LIMIT) });
      // The API registers stats under its own resource name, the last segment
      // of the endpoint (purchase_requests), not the admin slug
      // (purchase-requests), so every multi-word resource's card failed.
      const name = endpoint.split("/").filter(Boolean).pop();
      const { data } = await apiClient.get<{ data: ResourceDashboardStats }>(
        "/api/admin/dashboard/resource-stats/" + name + "?" + search.toString(),
        { signal },
      );
      return data.data;
    },
    select,
    // Keep stats around so coming back to the dashboard does not re-flash
    // the skeleton.
    staleTime: 30_000,
    refetchInterval: 60_000,
  });
}

// useRelationshipOptions lists the records a relationship picker offers: the
// first 100 of the related resource, narrowed on the server by what was typed,
// so a record past the first page can still be found. The single and the multi
// picker both use it, so a record created from either one refreshes both.
// Pass enabled: false until the options are needed, and a picker nobody opens
// asks for nothing.
export function useRelationshipOptions(
  endpoint: string,
  { search = "", enabled = true }: { search?: string; enabled?: boolean } = {},
) {
  const term = search.trim();
  return useQuery<Record<string, unknown>[]>({
    queryKey: resourceKeys.options(endpoint, term),
    queryFn: async ({ signal }) => {
      const params: Record<string, string> = { page_size: "100" };
      if (term) params.search = term;
      const { data } = await apiClient.get(endpoint, { params, signal });
      return Array.isArray(data?.data) ? data.data : Array.isArray(data) ? data : [];
    },
    enabled: enabled && !!endpoint,
    // The last answer stays up while the next search runs.
    placeholderData: keepPreviousData,
  });
}

// Every mutation hook takes an optional resource label (the singular, e.g.
// "Invoice") so toasts name what actually happened — "Invoice created
// successfully" rather than a bare "Created successfully". Omitting it keeps
// the old generic wording, so existing call sites still compile.
function said(label: string | undefined, verb: string) {
  return label ? label + " " + verb : verb.charAt(0).toUpperCase() + verb.slice(1);
}

function failed(label: string | undefined, verb: string) {
  return label ? "Failed to " + verb + " " + label.toLowerCase() : "Failed to " + verb;
}

// The four row mutations differ in three things: the HTTP method, the URL they
// build, and the word in the toast. They were four copies of the same eighteen
// lines, which is how usePatchResource ended up invalidating differently from
// useUpdateResource for no reason anybody could name.
//
// refresh says what a success should refetch: "all" for a create or delete,
// which can add or remove a row from any list, and "saved" for an edit, which
// cannot change a count and so leaves the stat cards alone.
type MutationSpec<TVars> = {
  method: "post" | "put" | "patch" | "delete";
  url: (endpoint: string, vars: TVars) => string;
  body?: (vars: TVars) => Record<string, unknown> | undefined;
  refresh: "all" | "saved";
  success: (label?: string) => string;
  failure: (label?: string) => string;
};

function makeResourceMutation<TVars>(spec: MutationSpec<TVars>) {
  return function useResourceMutation(endpoint: string, label?: string) {
    const queryClient = useQueryClient();

    return useMutation({
      mutationFn: async (vars: TVars) => {
        const { data } = await apiClient.request({
          method: spec.method,
          url: spec.url(endpoint, vars),
          data: spec.body?.(vars),
        });
        return data;
      },
      onSuccess: () => {
        if (spec.refresh === "all") {
          queryClient.invalidateQueries({ queryKey: resourceKeys.all(endpoint) });
        } else {
          refreshAfterSave(queryClient, endpoint);
        }
        toast.success(spec.success(label));
      },
      onError: (err: unknown) => {
        toast.error(getApiErrorMessage(err, spec.failure(label)));
      },
    });
  };
}

export const useCreateResource = makeResourceMutation<Record<string, unknown>>({
  method: "post",
  url: (endpoint) => endpoint,
  body: (body) => body,
  refresh: "all",
  success: (label) => said(label, "created successfully"),
  failure: (label) => failed(label, "create"),
});

export const useUpdateResource = makeResourceMutation<{ id: string; body: Record<string, unknown> }>({
  method: "put",
  url: (endpoint, { id }) => endpoint + "/" + id,
  body: ({ body }) => body,
  refresh: "saved",
  success: (label) => said(label, "updated successfully"),
  failure: (label) => failed(label, "update"),
});

// v3.31.18: partial updates for the grouped update view. Each group's Save
// button calls patch() with only the fields it owns. The Go-side Patch handler
// whitelists writable columns and silently drops anything else, so it is safe
// to send only a subset.
export const usePatchResource = makeResourceMutation<{ id: string; body: Record<string, unknown> }>({
  method: "patch",
  url: (endpoint, { id }) => endpoint + "/" + id,
  body: ({ body }) => body,
  refresh: "saved",
  success: (label) => (label ? label + " saved" : "Saved"),
  failure: (label) => failed(label, "save"),
});

export const useDeleteResource = makeResourceMutation<string>({
  method: "delete",
  url: (endpoint, id) => endpoint + "/" + id,
  refresh: "all",
  success: (label) => said(label, "deleted successfully"),
  failure: (label) => failed(label, "delete"),
});

export type BulkOperation = "delete" | "archive" | "restore" | "patch";

export interface BulkPayload {
  action: BulkOperation;
  ids: string[];
  /** Only read for "patch". */
  patch?: Record<string, unknown>;
}

const BULK_PAST: Record<BulkOperation, string> = {
  delete: "deleted",
  archive: "archived",
  restore: "restored",
  patch: "updated",
};

/**
 * One request for the whole selection, against POST <endpoint>/bulk.
 *
 * This used to be N parallel DELETEs, which is N transactions and N audit
 * entries, and leaves a half-applied result when the eleventh fails: the
 * operator is told it failed while ten rows are already gone. The server does
 * it in one transaction now, so the answer is all or nothing.
 *
 * Takes the PLURAL label ("Invoices") because the message counts rows, and
 * reports what the server actually did rather than what was asked: archiving
 * twelve rows of which three were already archived says nine.
 */
export function useBulkResource(endpoint: string, pluralLabel?: string, singularLabel?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    // ids are strings because Grit's models use UUID primary keys.
    mutationFn: async (payload: BulkPayload) => {
      try {
        const { data } = await apiClient.post(`${endpoint}/bulk`, payload);
        return { ...data, action: payload.action } as {
          data?: { affected: number; requested: number };
          action: BulkOperation;
        };
      } catch (err) {
        // No /bulk route on this endpoint. That is the normal state of an
        // upgraded project: grit upgrade replaces the admin but never
        // regenerates API handlers, so the browser gets the new code and the
        // server keeps the old routes. Falling back per row keeps the button
        // working instead of 404ing on every existing install.
        //
        // The fallback is genuinely worse: N requests, N transactions, and a
        // partial result if one fails. Run grit generate for the resource to
        // get the real endpoint.
        const status = (err as { response?: { status?: number } })?.response?.status;
        if (status !== 404) throw err;

        const results = await Promise.allSettled(
          payload.ids.map((id) =>
            payload.action === "delete"
              ? apiClient.delete(`${endpoint}/${id}`)
              : apiClient.patch(`${endpoint}/${id}`, payload.patch ?? {}),
          ),
        );
        const affected = results.filter((r) => r.status === "fulfilled").length;
        if (affected === 0) throw err;
        return {
          data: { affected, requested: payload.ids.length },
          action: payload.action,
        };
      }
    },
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      const affected = result?.data?.affected ?? 0;
      const requested = result?.data?.requested ?? affected;
      const noun = affected === 1 ? (singularLabel ?? "item") : (pluralLabel ?? "items");

      if (affected === 0) {
        // Not a success worth celebrating and not an error either. Saying
        // "0 archived" beats a green tick over a table that did not change.
        toast("Nothing to " + result.action + ": no matching rows");
        return;
      }
      const skipped = requested - affected;
      toast.success(
        affected + " " + noun + " " + BULK_PAST[result.action] +
          (skipped > 0 ? " (" + skipped + " already were)" : "")
      );
    },
    onError: (err: unknown) => {
      toast.error(getApiErrorMessage(err, "Bulk action failed. Nothing was changed."));
    },
  });
}

/**
 * Kept so existing call sites and hand-written pages keep working. Delegates
 * to the bulk endpoint rather than firing one request per row.
 *
 * @deprecated Use useBulkResource, which also archives, restores and patches.
 */
export function useBulkDeleteResource(endpoint: string, pluralLabel?: string) {
  const bulk = useBulkResource(endpoint, pluralLabel);
  return {
    ...bulk,
    mutate: (ids: string[], options?: Parameters<typeof bulk.mutate>[1]) =>
      bulk.mutate({ action: "delete", ids }, options),
  };
}
