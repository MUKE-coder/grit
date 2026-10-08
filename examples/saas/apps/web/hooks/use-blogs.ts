"use client";

import { useQuery } from "@tanstack/react-query";
import type { Blog, PaginatedResponse } from "@repo/shared/types";
import { api } from "@/lib/api";

// Types live in packages/shared — never inline them in a hook. Reasons:
//   - the same Blog shape is consumed by web, admin, and the Go API (via
//     grit sync's generated counterpart). One source of truth.
//   - editing the model in Go and running 'grit sync' updates ALL
//     consumers in one shot; inline duplicates silently drift.
type BlogMeta = PaginatedResponse<Blog>["meta"];

export function usePublicBlogs(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ["public-blogs", page, pageSize],
    queryFn: async (): Promise<{ blogs: Blog[]; meta: BlogMeta | undefined }> => {
      const { data } = await api.get<PaginatedResponse<Blog>>(
        `/api/blogs?page=${page}&page_size=${pageSize}`
      );
      return {
        blogs: data.data ?? [],
        meta: data.meta,
      };
    },
  });
}

export function usePublicBlog(slug: string) {
  return useQuery({
    queryKey: ["public-blog", slug],
    queryFn: async () => {
      const { data } = await api.get(`/api/blogs/${slug}`);
      return data.data as Blog;
    },
    enabled: !!slug,
  });
}

// Vite/TanStack blog routes import { useBlogs, useBlog } from this hook.
// Aliases so a single file works for both Next.js and Vite scaffolds.
export function useBlogs(page = 1, pageSize = 20) {
  const result = usePublicBlogs(page, pageSize);
  return { ...result, data: result.data?.blogs };
}

export const useBlog = usePublicBlog;
