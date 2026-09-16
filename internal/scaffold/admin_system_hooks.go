package scaffold

import "strings"

// ─── M40: the system screens read through hooks ─────────────────────────────
//
// Four of the system screens talked to the API from inside the page: the form
// shares screen with eight apiClient calls and four useMutation blocks, SSO
// with five, access reviews with five, API keys with three. Each one carried
// its own query keys, its own invalidation and its own error handling, so the
// same endpoint was described differently in the page that read it and the
// page that wrote it, and nothing else could reuse either.
//
// These files are where those calls live now, beside use-system.ts. A page is
// left with its form state and its markup, which is all a page should have.

// adminUseAPIKeys is hooks/use-api-keys.ts.
func adminUseAPIKeys() string {
	src := `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-core";

export interface APIKey {
  id: string;
  name: string;
  kind: "publishable" | "secret";
  prefix: string;
  token?: string;
  endpoints?: string[];
  origins?: string[];
  rate_limit?: number;
  last_used_at?: string;
  expires_at?: string;
  revoked_at?: string;
  created_at: string;
}

export const apiKeyKeys = {
  all: ["api-keys"] as const,
};

export function useAPIKeys() {
  return useQuery<APIKey[]>({
    queryKey: apiKeyKeys.all,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/api-keys");
      return (data.data ?? []) as APIKey[];
    },
  });
}

// The page builds the body, because which fields it sends depends on which
// boxes the operator filled in. Everything after the request is here.
export function useCreateAPIKey() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/api-keys", body);
      return data.data as { token: string };
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: apiKeyKeys.all }),
    onError: (err: unknown) => toast.error(getApiErrorMessage(err, "Could not create the key.")),
  });
}

export function useRevokeAPIKey() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/api-keys/" + id);
    },
    onSuccess: () => {
      toast.success("Key revoked. Requests using it will now be refused.");
      queryClient.invalidateQueries({ queryKey: apiKeyKeys.all });
    },
    onError: (err: unknown) => toast.error(getApiErrorMessage(err, "Could not revoke the key.")),
  });
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// adminUseSSO is hooks/use-sso.ts.
func adminUseSSO() string {
	src := `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export const ssoKeys = {
  connections: ["sso-connections"] as const,
  test: (id: string) => ["sso-connection-test", id] as const,
};

// live is how many of them the API actually built at boot: for OIDC that means
// discovery against the customer's IdP succeeded. A connection that is saved
// but not live is the case the screen exists to make visible.
// Generic over the row: the screen owns the connection type it renders, so the
// hook does not keep a second copy of it to drift.
export function useSSOConnections<T = Record<string, unknown>>() {
  return useQuery<{ rows: T[]; live: number }>({
    queryKey: ssoKeys.connections,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/sso/connections");
      return { rows: (data.data ?? []) as T[], live: data.meta?.live ?? 0 };
    },
  });
}

// One hook for create and update, because the screen has one form and the only
// difference is whether an id came with it.
export function useSaveSSOConnection() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, body }: { id?: string; body: Record<string, unknown> }) => {
      if (id) {
        const { data } = await apiClient.put("/api/sso/connections/" + id, body);
        return data;
      }
      const { data } = await apiClient.post("/api/sso/connections", body);
      return data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ssoKeys.connections }),
  });
}

export function useDeleteSSOConnection() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/sso/connections/" + id);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ssoKeys.connections }),
  });
}

// Runs the connection's discovery against the customer's IdP and reports what
// came back. Not cached: the point of pressing Test is to ask again.
export function useTestSSOConnection() {
  return useMutation({
    mutationFn: async (id: string) => {
      const { data } = await apiClient.get("/api/sso/connections/" + id + "/test");
      return data.data as { ok: boolean; message?: string };
    },
  });
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// adminUseFormShares is hooks/use-form-shares.ts.
func adminUseFormShares() string {
	src := `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export interface PublicField {
  key: string;
  label: string;
  type: string;
  required: boolean;
}

export const formShareKeys = {
  shares: ["form-shares"] as const,
  submissions: (shareID?: string) => ["form-submissions", shareID ?? "all"] as const,
  resources: ["form-share-resources"] as const,
  fields: (resource: string) => ["form-share-fields", resource] as const,
};

export function useFormShares<T = unknown>() {
  return useQuery<{ data: T[] }>({
    queryKey: formShareKeys.shares,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/form-shares");
      return data;
    },
  });
}

export function useFormSubmissions<T = unknown>(shareID?: string) {
  return useQuery<{ data: T[] }>({
    queryKey: formShareKeys.submissions(shareID),
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/form-submissions", {
        params: shareID ? { share_id: shareID } : undefined,
      });
      return data;
    },
  });
}

// The resources a share can be opened on. The server decides: a resource is
// reachable through a public form only when it says so, which is the security
// boundary this screen sits in front of.
export function useFormShareResources() {
  return useQuery<string[]>({
    queryKey: formShareKeys.resources,
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: string[] }>("/api/admin/form-shares/resources");
      return data.data ?? [];
    },
    staleTime: 5 * 60_000,
  });
}

export function useFormShareFields(resource: string) {
  return useQuery<PublicField[]>({
    queryKey: formShareKeys.fields(resource),
    enabled: !!resource,
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: { fields: PublicField[] } }>(
        "/api/admin/form-shares/resources/" + resource + "/fields"
      );
      return data.data?.fields ?? [];
    },
    staleTime: 5 * 60_000,
  });
}

// One invalidation for all three writes: a share list, its submissions and the
// row being edited are all views of the same thing.
function refreshShares(queryClient: ReturnType<typeof useQueryClient>) {
  queryClient.invalidateQueries({ queryKey: formShareKeys.shares });
}

export function useCreateFormShare() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/admin/form-shares", body);
      return data.data;
    },
    onSuccess: () => refreshShares(queryClient),
  });
}

export function useUpdateFormShare() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, body }: { id: string; body: Record<string, unknown> }) => {
      const { data } = await apiClient.patch("/api/admin/form-shares/" + id, body);
      return data.data;
    },
    onSuccess: () => refreshShares(queryClient),
  });
}

export function useDeleteFormShare() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/admin/form-shares/" + id);
    },
    onSuccess: () => refreshShares(queryClient),
  });
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// adminUseAccessReviews is hooks/use-access-reviews.ts.
func adminUseAccessReviews() string {
	src := `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export const accessReviewKeys = {
  list: ["access-reviews"] as const,
  detail: (id: string | null) => ["access-review", id] as const,
};

export function useAccessReviews<T = unknown>() {
  return useQuery<T[]>({
    queryKey: accessReviewKeys.list,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/access-reviews");
      return (data.data ?? []) as T[];
    },
  });
}

export function useAccessReview<T = unknown>(id: string | null) {
  return useQuery<T>({
    queryKey: accessReviewKeys.detail(id),
    enabled: !!id,
    queryFn: async () => {
      const { data } = await apiClient.get("/api/access-reviews/" + id);
      return data.data as T;
    },
  });
}

// A campaign and its rows move together: opening one, deciding a row and
// completing it all change both the list's counts and the detail.
function refreshReviews(queryClient: ReturnType<typeof useQueryClient>, id: string | null) {
  queryClient.invalidateQueries({ queryKey: accessReviewKeys.list });
  if (id) queryClient.invalidateQueries({ queryKey: accessReviewKeys.detail(id) });
}

export function useOpenAccessReview<T = unknown>() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: { name: string; note?: string }) => {
      const { data } = await apiClient.post("/api/access-reviews", body);
      return data.data as T;
    },
    onSuccess: () => refreshReviews(queryClient, null),
  });
}

export function useDecideAccessReviewItem(reviewID: string | null) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (args: { itemId: string; decision: "approved" | "revoked" }) => {
      await apiClient.post(
        "/api/access-reviews/" + reviewID + "/items/" + args.itemId + "/decision",
        { decision: args.decision }
      );
    },
    onSuccess: () => refreshReviews(queryClient, reviewID),
  });
}

export function useCompleteAccessReview(reviewID: string | null) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      await apiClient.post("/api/access-reviews/" + reviewID + "/complete");
    },
    onSuccess: () => refreshReviews(queryClient, reviewID),
  });
}
`
	return strings.ReplaceAll(src, "~", "`")
}
