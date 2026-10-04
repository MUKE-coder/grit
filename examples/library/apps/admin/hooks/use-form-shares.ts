import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
