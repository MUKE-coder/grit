import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import type { Subscription } from "@repo/shared/types";

interface SubscriptionsResponse {
  data: Subscription[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

interface UseSubscriptionsParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
}

export function useSubscriptions({ page = 1, pageSize = 20, search = "", sortBy = "created_at", sortOrder = "desc" }: UseSubscriptionsParams = {}) {
  return useQuery<SubscriptionsResponse>({
    queryKey: ["subscriptions", { page, pageSize, search, sortBy, sortOrder }],
    queryFn: async () => {
      const params = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
        sort_by: sortBy,
        sort_order: sortOrder,
      });
      if (search) {
        params.set("search", search);
      }
      const { data } = await apiClient.get(`/api/subscriptions?${params}`);
      return data;
    },
  });
}

export function useGetSubscription(id: string) {
  return useQuery<Subscription>({
    queryKey: ["subscriptions", id],
    queryFn: async () => {
      const { data } = await apiClient.get(`/api/subscriptions/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
}

export function useCreateSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/subscriptions", input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["subscriptions"] });
    },
  });
}

export function useUpdateSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: { id: string } & Record<string, unknown>) => {
      const { data } = await apiClient.put(`/api/subscriptions/${id}`, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["subscriptions"] });
    },
  });
}

export function useDeleteSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(`/api/subscriptions/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["subscriptions"] });
    },
  });
}
