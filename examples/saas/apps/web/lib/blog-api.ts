import { cache } from "react";
import type { Blog, PaginatedResponse } from "@repo/shared/types";
import { API_VERSION } from "@/lib/api";

// Server-side reads of the public blog API, for server components.
//
// In Docker the web container reaches the API by its service name, which the
// production compose file passes as API_INTERNAL_URL. Everywhere else the
// public URL works from the server as well.
const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

// The pages revalidate on the same interval, so a published edit shows within
// a minute.
const REVALIDATE_SECONDS = 60;

type BlogPage = {
  blogs: Blog[];
  meta: PaginatedResponse<Blog>["meta"] | undefined;
};

async function apiGet<T>(path: string): Promise<{ status: number; body: T | null }> {
  try {
    const res = await fetch(`${API_URL}/api/${API_VERSION}${path}`, {
      next: { revalidate: REVALIDATE_SECONDS },
    });
    if (!res.ok) return { status: res.status, body: null };
    return { status: res.status, body: (await res.json()) as T };
  } catch {
    // The API is unreachable, as it is while next build runs in CI. A list
    // renders without posts and fills in at the next revalidation.
    return { status: 0, body: null };
  }
}

// Published posts, newest first. cache() shares one request between the page
// and its metadata.
export const getPublishedBlogs = cache(
  async (page: number, pageSize: number): Promise<BlogPage> => {
    const { body } = await apiGet<PaginatedResponse<Blog>>(
      `/blogs?page=${page}&page_size=${pageSize}`
    );
    return { blogs: body?.data ?? [], meta: body?.meta };
  }
);

// One published post, or null when there is none by that slug. An API that
// fails is an error, not a missing post, so it is not cached as a 404.
export const getPublishedBlog = cache(async (slug: string): Promise<Blog | null> => {
  const { status, body } = await apiGet<{ data: Blog }>(`/blogs/${encodeURIComponent(slug)}`);
  if (status === 404) return null;
  if (!body) throw new Error(`the blog API answered ${status || "nothing"} for ${slug}`);
  return body.data;
});
