"use client";

// v3.31.40 -- per-user dashboard customisation.
//
// useDashboardLayout pulls the saved layout from /api/dashboard-layout.
// The response shape carries id, user_id, three kind-specific key
// arrays, and a date_preset. Empty id = "no row in DB yet" -- the
// dashboard renders the full default catalog in that case.
//
// useSaveDashboardLayout pushes the whole layout back as a PUT.
// Whole-resource replace, not patch, because the payload is tiny and
// the semantics are simpler -- whatever you send is the new layout.

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";
import type { SavedLayout } from "@/lib/dashboard-catalog";

export function useDashboardLayout() {
  return useQuery<SavedLayout>({
    queryKey: ["dashboard-layout"],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: SavedLayout }>("/api/dashboard-layout");
      return data.data;
    },
    staleTime: 5 * 60_000,
  });
}

export function useSaveDashboardLayout() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (body: {
      cards: string[];
      charts: string[];
      tables: string[];
      // v3.31.45 -- two new arrays. Both optional on the wire; if
      // omitted the API treats them as empty.
      resources?: string[];
      section_order?: string[];
      // v3.31.46 -- per-resource layout map. Keys are slugs; values
      // are "split" (default, can be omitted) or "tabs".
      resource_layouts?: Record<string, "split" | "tabs">;
      // v3.31.47 -- user-defined chart configs.
      custom_charts?: import("@/lib/dashboard-catalog").CustomChart[];
      date_preset: string;
    }) => {
      const { data } = await apiClient.put<{ data: SavedLayout }>(
        "/api/dashboard-layout",
        body,
      );
      return data.data;
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["dashboard-layout"] });
      toast.success("Dashboard preferences saved");
    },
    onError: () => {
      toast.error("Failed to save dashboard preferences");
    },
  });
}
