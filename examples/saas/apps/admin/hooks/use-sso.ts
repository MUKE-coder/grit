import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
