import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import type { UsageRecord } from "@repo/shared/types";

interface UsageRecordsResponse {
  data: UsageRecord[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

interface UseUsageRecordsParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
}

export function useUsageRecords({ page = 1, pageSize = 20, search = "", sortBy = "created_at", sortOrder = "desc" }: UseUsageRecordsParams = {}) {
  return useQuery<UsageRecordsResponse>({
    queryKey: ["usage_records", { page, pageSize, search, sortBy, sortOrder }],
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
      const { data } = await apiClient.get(`/api/usage_records?${params}`);
      return data;
    },
  });
}

export function useGetUsageRecord(id: string) {
  return useQuery<UsageRecord>({
    queryKey: ["usage_records", id],
    queryFn: async () => {
      const { data } = await apiClient.get(`/api/usage_records/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
}

export function useCreateUsageRecord() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/usage_records", input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["usage_records"] });
    },
  });
}

export function useUpdateUsageRecord() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: { id: string } & Record<string, unknown>) => {
      const { data } = await apiClient.put(`/api/usage_records/${id}`, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["usage_records"] });
    },
  });
}

export function useDeleteUsageRecord() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(`/api/usage_records/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["usage_records"] });
    },
  });
}
