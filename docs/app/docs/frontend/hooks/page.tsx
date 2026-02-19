import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function HooksPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Frontend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                React Query Hooks
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every resource generated with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code> produces
                a complete set of React Query hooks for CRUD operations with pagination, search, sorting,
                cache invalidation, and optimistic updates.
              </p>
            </div>

            <div className="prose-grit">
              {/* Auto-generated Hooks Pattern */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Auto-Generated Hooks Pattern
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource Post</code>, Grit
                  creates a hook file at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/use-posts.ts</code> in
                  both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin</code>.
                  Each file contains five hooks that map to the CRUD API endpoints.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Hook</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">HTTP Method</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">usePosts()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">GET /api/posts</td>
                        <td className="px-4 py-2.5">Paginated list with sorting and search</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">useGetPost(id)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">GET /api/posts/:id</td>
                        <td className="px-4 py-2.5">Fetch a single post by ID</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">useCreatePost()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">POST /api/posts</td>
                        <td className="px-4 py-2.5">Create a new post</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">useUpdatePost()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PUT /api/posts/:id</td>
                        <td className="px-4 py-2.5">Update an existing post</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">useDeletePost()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">DELETE /api/posts/:id</td>
                        <td className="px-4 py-2.5">Soft-delete a post</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* useList Hook */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  useList Hook
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The list hook accepts pagination, search, and sort parameters. It returns a React Query
                  result with the paginated data and metadata (total, page, pages).
                </p>

                <CodeBlock filename="hooks/use-posts.ts (list hook)" code={`interface UsePostsParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
}

export function usePosts({
  page = 1,
  pageSize = 20,
  search = "",
  sortBy = "created_at",
  sortOrder = "desc",
}: UsePostsParams = {}) {
  return useQuery<PostsResponse>({
    queryKey: ["posts", { page, pageSize, search, sortBy, sortOrder }],
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
      const { data } = await apiClient.get(\`/api/posts?\${params}\`);
      return data;
    },
  });
}`} />
              </div>

              {/* useGetById Hook */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  useGetById Hook
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Fetches a single resource by ID. The query is disabled when no valid ID is provided,
                  preventing unnecessary API calls. The query key includes the resource name and ID for
                  precise cache targeting.
                </p>

                <CodeBlock filename="hooks/use-posts.ts (getById hook)" code={`export function useGetPost(id: number) {
  return useQuery<Post>({
    queryKey: ["posts", id],
    queryFn: async () => {
      const { data } = await apiClient.get(\`/api/posts/\${id}\`);
      return data.data;
    },
    enabled: !!id,
  });
}`} />
              </div>

              {/* Mutation Hooks */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  useCreate, useUpdate, useDelete
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  All mutation hooks use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useMutation</code> from React Query.
                  On success, each mutation automatically invalidates the list query cache, causing the table
                  to refetch fresh data from the server.
                </p>

                <CodeBlock filename="hooks/use-posts.ts (mutations)" code={`export function useCreatePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/posts", input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}

export function useUpdatePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      ...input
    }: { id: number } & Record<string, unknown>) => {
      const { data } = await apiClient.put(\`/api/posts/\${id}\`, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}

export function useDeletePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      await apiClient.delete(\`/api/posts/\${id}\`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}`} />
              </div>

              {/* Query Key Patterns */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Query Key Patterns
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit follows a consistent query key structure that enables precise cache management.
                  The first element is always the resource name (plural), and optional subsequent elements
                  narrow the scope.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Query Key</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Scope</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">[&quot;posts&quot;]</td>
                        <td className="px-4 py-2.5">All post queries (invalidates everything)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">[&quot;posts&quot;, {'{...params}'}]</td>
                        <td className="px-4 py-2.5">Specific list page with filters</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">[&quot;posts&quot;, 42]</td>
                        <td className="px-4 py-2.5">Single post by ID</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  When a mutation calls <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">invalidateQueries({'{'}queryKey: [&quot;posts&quot;]{'}'})</code>,
                  React Query invalidates all queries whose key starts with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">[&quot;posts&quot;]</code> --
                  including all paginated list queries and individual item queries. This is why your table
                  automatically refreshes after creating, updating, or deleting a post.
                </p>
              </div>

              {/* Cache Invalidation */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Cache Invalidation on Mutations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit&apos;s generated hooks use the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">onSuccess</code> callback
                  to invalidate related queries after every mutation. This pattern ensures the UI always
                  reflects the latest server state without manual refetching.
                </p>

                <CodeBlock filename="invalidation-flow.txt" code={`User clicks "Delete" on row
  -> useDeletePost.mutate(42)
    -> DELETE /api/posts/42
      -> 200 OK
        -> onSuccess fires
          -> queryClient.invalidateQueries(["posts"])
            -> usePosts() refetches from server
              -> Table re-renders with updated data`} />
              </div>

              {/* Pagination + Search */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Pagination + Search Parameters
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The list hook accepts all pagination and filter parameters as a single options object.
                  Each combination of parameters creates a unique cache entry, so navigating between pages
                  is instant on revisit.
                </p>

                <CodeBlock language="tsx" filename="example-usage.tsx" code={`// Fetch page 2 with 10 items, search for "react", sort by title
const { data, isLoading, error } = usePosts({
  page: 2,
  pageSize: 10,
  search: "react",
  sortBy: "title",
  sortOrder: "asc",
});

// Access the paginated response
console.log(data?.data);          // Post[]
console.log(data?.meta.total);    // 42
console.log(data?.meta.pages);    // 5
console.log(data?.meta.page);     // 2`} />
              </div>

              {/* Full Component Example */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Using Hooks in Components
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is a complete example of a posts list page using the generated hooks,
                  including search, pagination, and delete functionality.
                </p>

                <CodeBlock language="tsx" filename="app/(dashboard)/posts/page.tsx" code={`"use client";

import { useState } from "react";
import { usePosts, useDeletePost } from "@/hooks/use-posts";
import { toast } from "sonner";

export default function PostsPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState("created_at");
  const [sortOrder, setSortOrder] = useState("desc");

  const { data, isLoading, error } = usePosts({
    page,
    pageSize: 20,
    search,
    sortBy,
    sortOrder,
  });

  const deletePost = useDeletePost();

  const handleDelete = (id: number) => {
    deletePost.mutate(id, {
      onSuccess: () => toast.success("Post deleted"),
      onError: () => toast.error("Failed to delete post"),
    });
  };

  const handleSort = (column: string) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortBy(column);
      setSortOrder("asc");
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading posts</div>;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Posts</h1>
        <input
          type="text"
          placeholder="Search posts..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1); // Reset to first page on search
          }}
          className="px-3 py-2 rounded-lg border border-border bg-bg-secondary"
        />
      </div>

      {/* Table */}
      <table className="w-full">
        <thead>
          <tr>
            <th onClick={() => handleSort("title")}>Title</th>
            <th onClick={() => handleSort("created_at")}>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.data.map((post) => (
            <tr key={post.id}>
              <td>{post.title}</td>
              <td>{new Date(post.created_at).toLocaleDateString()}</td>
              <td>
                <button onClick={() => handleDelete(post.id)}>
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Pagination */}
      <div className="flex items-center gap-2 mt-4">
        <button
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page === 1}
        >
          Previous
        </button>
        <span>
          Page {data?.meta.page} of {data?.meta.pages}
        </span>
        <button
          onClick={() => setPage((p) => p + 1)}
          disabled={page >= (data?.meta.pages ?? 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
}`} />
              </div>

              {/* Custom Hooks */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Custom Hooks Beyond Generated Ones
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The generated hooks cover standard CRUD operations. For custom endpoints, queries, or
                  business logic, create additional hooks following the same pattern. Use the shared
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apiClient</code> instance
                  to ensure JWT tokens are automatically injected.
                </p>

                <CodeBlock language="typescript" filename="hooks/use-post-stats.ts" code={`import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

interface PostStats {
  total: number;
  published: number;
  draft: number;
  viewsThisWeek: number;
}

export function usePostStats() {
  return useQuery<PostStats>({
    queryKey: ["posts", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/posts/stats");
      return data.data;
    },
    staleTime: 60 * 1000, // 1 minute
  });
}

// Bulk publish mutation
export function useBulkPublish() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (ids: number[]) => {
      const { data } = await apiClient.post("/api/posts/bulk-publish", {
        ids,
      });
      return data;
    },
    onSuccess: () => {
      // Invalidate both posts list and stats
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}`} />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80">
                    <strong>Convention:</strong> Always place hooks in the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/</code> directory,
                    name them with the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">use-</code> prefix
                    and kebab-case file names, and follow the query key pattern <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">[&quot;resource&quot;, ...specifics]</code>.
                  </p>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/frontend/web-app" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Web App
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/frontend/shared-package" className="gap-1.5">
                  Shared Package
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
