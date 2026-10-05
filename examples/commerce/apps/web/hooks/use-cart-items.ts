import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import type { CartItem } from "@repo/shared/types";

interface CartItemsResponse {
  data: CartItem[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

interface UseCartItemsParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
}

export function useCartItems({ page = 1, pageSize = 20, search = "", sortBy = "created_at", sortOrder = "desc" }: UseCartItemsParams = {}) {
  return useQuery<CartItemsResponse>({
    queryKey: ["cart_items", { page, pageSize, search, sortBy, sortOrder }],
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
      const { data } = await apiClient.get(`/api/cart_items?${params}`);
      return data;
    },
  });
}

export function useGetCartItem(id: string) {
  return useQuery<CartItem>({
    queryKey: ["cart_items", id],
    queryFn: async () => {
      const { data } = await apiClient.get(`/api/cart_items/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
}

export function useCreateCartItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/cart_items", input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cart_items"] });
    },
  });
}

export function useUpdateCartItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: { id: string } & Record<string, unknown>) => {
      const { data } = await apiClient.put(`/api/cart_items/${id}`, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cart_items"] });
    },
  });
}

export function useDeleteCartItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(`/api/cart_items/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cart_items"] });
    },
  });
}
