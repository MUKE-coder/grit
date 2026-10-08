import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
