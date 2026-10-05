import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
