package scaffold

// The web app's generic resource hooks.
//
// `grit generate resource` used to write ninety-two lines of React Query per
// resource into apps/web/hooks, the same file with the names changed, and it
// typed every mutation input as Record<string, unknown> although
// CreateProductInput was being generated two directories away in
// packages/shared. So the compiler could not tell you that you had sent
// `title` where the API wanted `name`, and a fix to the hook shape was one
// edit per resource per project (grit#146).
//
// The admin has had the good version of this for a long time. The web app got
// the earlier shape and kept it.
//
// Not shared with the admin's, which the issue suggests and which I looked at:
// apps/admin/hooks/use-resource.ts is 429 lines and is wired to the admin's own
// api-client, its resource registry, its optimistic-update conventions and its
// toast surface. Moving it to packages/shared would mean a package that
// imports from both apps or a third set of conventions neither uses. Two small
// generics that each know their own app is the cheaper truth.

// webResourceHooks is apps/web/hooks/use-resource.ts.
func webResourceHooks() string {
	return `import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryOptions,
} from "@tanstack/react-query";
import { apiClient } from "@/lib/api";

// One set of React Query hooks for any generated resource.
//
// ` + "`grit generate resource`" + ` writes a three-line file that calls this with
// the resource's own types, rather than ninety lines of the same code with the
// names changed. A fix here reaches every resource in the project at once.

/** The list envelope every Grit list endpoint returns. */
export interface ListResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

export interface ListParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  /** Anything else the endpoint filters on, passed through as query params. */
  filters?: Record<string, string | number | boolean | undefined>;
}

/**
 * The query keys for one resource.
 *
 * A factory rather than inline arrays, so an invalidation and a query cannot
 * disagree about the key: ["products", { page: 1 }] and ["products", "1"] are
 * different caches, and typing one of them by hand at each call site is how a
 * list stops refreshing after a create.
 */
export function resourceKeys(resource: string) {
  return {
    all: [resource] as const,
    lists: () => [resource, "list"] as const,
    list: (params: ListParams) => [resource, "list", params] as const,
    details: () => [resource, "detail"] as const,
    detail: (id: string) => [resource, "detail", id] as const,
  };
}

/**
 * Build the hooks for one resource.
 *
 * @param resource the cache key prefix, e.g. "products"
 * @param endpoint the API path, e.g. "/api/products"
 *
 * T is the row, C the create input and U the update input. All three come from
 * packages/shared, which is generated from the Go model, so a field renamed in
 * Go is a type error here rather than a 422 at runtime.
 */
/**
 * Whether an error is the API refusing a stale write.
 *
 * By code rather than by status: 409 is also what a unique-constraint
 * violation answers, and that one is the caller's to fix rather than a reason
 * to reload.
 */
export function isVersionConflict(error: unknown): boolean {
  const body = (error as { response?: { data?: { error?: { code?: string } } } })?.response?.data;
  return body?.error?.code === "VERSION_CONFLICT";
}

/** The version the record is at now, out of that 409. */
export function currentVersion(error: unknown): number | undefined {
  const details = (error as {
    response?: { data?: { error?: { details?: { current_version?: unknown } } } };
  })?.response?.data?.error?.details;
  const current = details?.current_version;
  return typeof current === "number" ? current : undefined;
}

export function resourceHooks<T, C = Partial<T>, U = Partial<C>>(
  resource: string,
  endpoint: string,
) {
  const keys = resourceKeys(resource);

  function useList(
    params: ListParams = {},
    options?: Omit<UseQueryOptions<ListResponse<T>>, "queryKey" | "queryFn">,
  ) {
    const {
      page = 1,
      pageSize = 20,
      search = "",
      sortBy = "created_at",
      sortOrder = "desc",
      filters,
    } = params;

    return useQuery<ListResponse<T>>({
      queryKey: keys.list({ page, pageSize, search, sortBy, sortOrder, filters }),
      queryFn: async () => {
        const query = new URLSearchParams({
          page: String(page),
          page_size: String(pageSize),
          sort_by: sortBy,
          sort_order: sortOrder,
        });
        if (search) query.set("search", search);
        for (const [key, value] of Object.entries(filters ?? {})) {
          // An undefined filter is "not filtering", not "filter by the string
          // undefined", which is what String(value) would have sent.
          if (value !== undefined && value !== "") query.set(key, String(value));
        }
        const { data } = await apiClient.get(endpoint + "?" + query.toString());
        return data;
      },
      ...options,
    });
  }

  function useOne(id: string, options?: Omit<UseQueryOptions<T>, "queryKey" | "queryFn">) {
    return useQuery<T>({
      queryKey: keys.detail(id),
      queryFn: async () => {
        const { data } = await apiClient.get(endpoint + "/" + id);
        return data.data;
      },
      enabled: !!id,
      ...options,
    });
  }

  function useCreate() {
    const queryClient = useQueryClient();
    return useMutation({
      mutationFn: async (input: C) => {
        const { data } = await apiClient.post(endpoint, input);
        return data.data as T;
      },
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: keys.all });
      },
    });
  }

  /**
   * Save an edit.
   *
   * version is the one the record was at when it was loaded, from useOne, and
   * it is sent as If-Match. The API applies the write only if the record is
   * still there and answers 409 VERSION_CONFLICT with the current version if it
   * is not, so two people editing the same row cannot silently overwrite each
   * other. Pass undefined to mean "overwrite whatever is there", and say so
   * where you do.
   *
   * Required rather than optional on purpose: a call site that does not pass it
   * does not compile. Optional, every one of them would have gone on sending no
   * header, which is the state this replaced (grit#182).
   */
  function useUpdate() {
    const queryClient = useQueryClient();
    return useMutation({
      mutationFn: async ({
        id,
        version,
        ...input
      }: { id: string; version: number | undefined } & U) => {
        const { data } = await apiClient.put(endpoint + "/" + id, input, {
          headers:
            typeof version === "number"
              ? { "If-Match": 'W/"' + version + '"' }
              : undefined,
        });
        return data.data as T;
      },
      onSuccess: (_result, variables) => {
        // The row AND the lists. Invalidating only the list leaves a detail
        // page showing what was there before the edit it just made.
        queryClient.invalidateQueries({ queryKey: keys.detail(variables.id) });
        queryClient.invalidateQueries({ queryKey: keys.lists() });
      },
      onError: (error, variables) => {
        // Somebody else saved first. Refetch, so whatever the page does with
        // the error is reading the record as it is rather than as it was.
        if (isVersionConflict(error)) {
          queryClient.invalidateQueries({ queryKey: keys.detail(variables.id) });
          queryClient.invalidateQueries({ queryKey: keys.lists() });
        }
      },
    });
  }

  function useRemove() {
    const queryClient = useQueryClient();
    return useMutation({
      mutationFn: async (id: string) => {
        await apiClient.delete(endpoint + "/" + id);
        return id;
      },
      onSuccess: (id) => {
        queryClient.removeQueries({ queryKey: keys.detail(id) });
        queryClient.invalidateQueries({ queryKey: keys.lists() });
      },
    });
  }

  return { keys, useList, useOne, useCreate, useUpdate, useRemove };
}
`
}
